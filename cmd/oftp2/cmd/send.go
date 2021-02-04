package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thomsmits/oftp2-client/internal/liboftp2/client"
)

var sendCommand = &cobra.Command{
	Use:     "send ODETTEID FILEPATH DATASETNAME",
	Short:   "Send file to server",
	Long:    `Sends a file to the server`,
	Example: `oftp2 send O20222CUSTOMER /tmp/data DATA22`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 3 {
			return errors.New("please provide an receiving ODETTE ID, a file path and a dataset name\n")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		sendFile(args[0], args[1], args[2])
	},
}

func sendFile(odetteId, filePath, datasetName string) {

	s := client.OFTP2Client{
		ServerHost: activeOptions.Server,
		ServerPort: activeOptions.Port,
		Verbose:    activeOptions.Verbose,
		OdetteId:   activeOptions.OdetteId,
	}

	err := s.Connect()
	if err != nil {
		panic(err)
	}

	err = s.StartSession("", false, false, false)
	if err != nil {
		fmt.Printf("start session failed: %v\n", err)
		os.Exit(1)
	}

	err = s.SendFile(datasetName,
		filePath,
		client.FileFormatUnstructured,
		//"O2010CUSTOMER",
		odetteId,
		client.SecurityLevelNone,
		false,
		false,
		false,
		false)
	if err != nil {
		fmt.Printf("send file failed: %v\n", err)
		os.Exit(1)
	}

	err = s.EndSession()
	if err != nil {
		fmt.Printf("end session failed: %v\n", err)
		os.Exit(1)
	}

	s.Close()
}
