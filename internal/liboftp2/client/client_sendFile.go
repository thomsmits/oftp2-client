package client

import (
	"errors"
	"fmt"
	"os"

	"github.com/thomsmits/oftp2-client/internal/liboftp2/wire/endfile"
	"github.com/thomsmits/oftp2-client/internal/liboftp2/wire/startfile"
	"github.com/thomsmits/oftp2-client/internal/liboftp2/wire/transfer"
)

// SendFile sends a file to the OFTP2 server
func (s *OFTP2Client) SendFile(datasetName string, filePath string, format OFTP2FileFormat, destination string, securityLevel OFTP2SecurityLevel, cipher, compression, envelope, signed bool) error {

	var cipherSuite, compressionIndicator, envelopeIndicator int

	if cipher {
		cipherSuite = 1
	} else {
		cipherSuite = 0
	}

	if compression {
		compressionIndicator = 1
	} else {
		compressionIndicator = 0
	}

	if envelope {
		envelopeIndicator = 1
	} else {
		envelopeIndicator = 0
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	fileSize := fileInfo.Size()

	var maxRecordSize int

	if format == FileFormatText || format == FileFormatUnstructured {
		maxRecordSize = 0
	} else {
		// TODO: Determine correct value
		maxRecordSize = 0
	}

	// Report at least 1 kB of file size, even if file is smaller
	fileSizeInK := uint64(fileSize / 1024)
	if fileSizeInK == 0 && fileSize != 0 {
		fileSizeInK = 1
	}

	sfid := startfile.SFID{
		DatasetName:            datasetName,
		FileDateTime:           fileInfo.ModTime(),
		UserData:               "",
		Destination:            destination,
		Originator:             s.OdetteId,
		FileFormat:             string(format),
		MaxRecordSize:          maxRecordSize,
		FileSizeInK:            fileSizeInK,
		OriginalFileSizeInK:    fileSizeInK,
		RestartPosition:        0,
		SecurityLevel:          int(securityLevel),
		CipherSuite:            cipherSuite,
		Compression:            compressionIndicator,
		Envelope:               envelopeIndicator,
		SigningRequired:        signed,
		VirtualFileDescription: "",
	}

	err = s.write(sfid.Marshal())
	if err != nil {
		return err
	}

	// Read answer from server
	// Read answer from communication partner
	buffer, err := s.read()
	answer, t, err := DetermineMessageType(buffer)
	if err != nil {
		return err
	}

	if t == "SFNA" {
		return errors.New(fmt.Sprintf("partner does not accept file: %v", answer))
	} else if t != "SFPA" {
		return errors.New(fmt.Sprintf("unknown answer. Expected SFPA or SFNA, got %s", t))
	}

	// get server answer
	sfpa := startfile.SFPA{}
	err = sfpa.Parse(buffer)
	if err != nil {
		return err
	}

	if sfpa.AnswerCount > sfid.RestartPosition {
		return errors.New(fmt.Sprintf("restart positions do not fit we: %d, server: %d", sfid.RestartPosition, sfpa.AnswerCount))
	}

	// Server is ready to receive our data, so send it to it
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var bytesTransmitted uint64

	var credits = s.serverCredit
	var maxReadBufferSize = int(s.maxReadBufferSize())

	for true {
		fileBuffer := make([]byte, maxReadBufferSize)
		n, err := file.Read(fileBuffer)
		if n != 0 && err != nil {
			return err
		}

		if n == 0 {
			// EOF
			break
		}

		sendBuffer := s.splitBufferIntoSubRecords(fileBuffer[0:n], n < maxReadBufferSize)

		data := transfer.DATA{
			Length: uint64(n),
			Buffer: sendBuffer,
		}

		err = s.write(data.Marshal())
		if err != nil {
			return err
		}

		bytesTransmitted += uint64(n)

		credits--

		if credits <= 0 {
			// we exceeded our credits, wait for the server to send us a CDT command
			// get the credit command
			cdt := transfer.CDT{}

			buffer, err := s.read()
			err = cdt.Parse(buffer)
			if err != nil {
				return err
			}

			credits = s.serverCredit
		}

	}

	var recordCount uint64 = 0

	if format == FileFormatText || format == FileFormatUnstructured {
		recordCount = 0
	} else {
		// TODO: Determine correct value
		recordCount = 0
	}

	// End file
	efid := endfile.EFID{
		RecordCount: recordCount,
		UnitCount:   bytesTransmitted,
	}

	err = s.write(efid.Marshal())
	if err != nil {
		return err
	}

	return nil
}
