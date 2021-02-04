# Odette FTP 2 (OFTP2) Library

## Purpose

Odette FTP 2 (OFTP2) is a protocol specified in [RFC 5024](https://tools.ietf.org/html/rfc5024). It is a file transfer protocol with roots in the ancient B2B world of the 1980ies. In contrast to FTP it has no notion of directories but just transfers files between locations. Locations are identified using a so called _ODETTE ID_, which is formatted according to ISO-6523.

This project is a clean room implementation of the OFTP2 protocol, just using the specification found in the RFC. No other implementation was consulted or reviewed. Therefore, it can be expected to be free of external IP. See the [license](LICENSE).

The purpose is not to implement the protocol to its full extend but to provide the means to communication with an OFTP2 server as a client sending files to the server. For example, this project currently does not support the server side of the protocol. Nevertheless, regarding the protocol messages, the complete OFTP2 specification is implemented, making a server implementation at least feasible.

From this project, you get

  * an OFTP2 library for sending files to a server
  * a command-line utility using this library as a proof of concept

The library has more functionality than the command-line utility exposes. Therefore, if you want to do more fancy stuff, have a look at the `OFTP2Client` struct and its methods.

## Structure of project

  * [cmd/oft2/main.go](cmd/oft2/main.go) command line tool, to interact with OFTP2 servers based on the [cobra](https://github.com/spf13/cobra) framework
  * [internal/liboftp2](internal/liboftp2) Library to talk to OFTP2 servers
    * [client](internal/liboftp2/client) The client functions
    * [wire](internal/liboftp2/wire) Implementation of the different messages, clients and servers can exchange

## Talking to a server

### Odette ID

To talk to an OFTP2 server, you need a valid Odette ID that is configured as a communication partner on the server. If you know such an id, you can ask the server for its own id, with the `id` subcommand:

```
$ oftp2 id -s 10.0.0.1 -i MYODETTEID
```

You can try `LOCAL` and `O1999MENDELSONDE` which sometimes work out-of-the-box due to the server implementation.

The format of the ID is: `Onnnnccccccccccccccssssss`

  * `O`: the first character is always `O`
  * `nnnn`: a four-digit number with leading zeros, called the _International Code Designator_ (_ICD_). Only the numbers 0002-0189 are globally registered (see [ICD list](https://www.cyber-identity.com/download/ICD-list.pdf)). In Germany, you can expect to see `0013` and `0177` quite often because they are registered to the [VDA](https://www.vda.de) and [Odette consortium}(https://www.odette.org).
  * `cccccccccccccc`: A 14 character _Organisation Code_ that normally consists of the character set `[A-Z\-0-9]`. Which characters are allowed, depends on the [ICD](https://www.cyber-identity.com/download/ICD-list.pdf).
  * `ssssss`: a 6 character computer sub address that normally consists of the character set `[A-Z\-0-9]`.

The ID has not to take up the fill 25 characters but can be shorter. In this case, it is padded with blanks to the right.

If the Odette id starts with `O013` the next six characters are numeric, e.g. `O0013000243AEG`.

### Server Capabilities

Using a valid Odette ID for the sender and knowing the ID of the server (e.g. determined via the ID command), you can query the server for its capabilities using the `query` sub command of the tool. 

```
$ oftp2 query -s 10.0.0.1 -i MYODETTEID
```

A typical answer may be:

```
Data received from remote system:
------------------------------------------
OdetteId                : SERVERODETTEID
Max buffer size         : 10240
Capability              : Location can only receive files.
Compression supported   : false
Restart supported       : false
Special logic supported : false
Authentication supported: false
User data               :
------------------------------------------
```

The server can only receive files, does not support compression, authentication, or resuming a file transfer session. "Special logic" refers to additional processing steps performed with data received. "User data" is additional data, configured for the Odette ID you provided as the sender.

### Sending a file

If the server does not require authentication - which is common - you can simply send data to the server if you have a valid Odette ID. This is done with the `send` sub command, for example:

```
$ oftp2 send -s 10.0.0.1 -i MYODETTEID SERVERODETTEID /tmp/filetosend DATASETNAME
```

The `SERVERODETTEID` and be the ID determined by the `id` sub command or another ID, the server is configured to receive files for. `DATASETNAME` is an arbitrary name chosen for the data. The command above send ths file `/tmp/filetosend` to the sever `10.0.0.1`. 

## Fuzzing

This project was created during a security assessment to find vulnerabilities in an OFTP2 installation. Therefore, the `OFTP2Client` structure contains a field `Fuzzer` for a callback function that allows you to mess around with the data send over the wire to the server. To use it, implement a new command in [cmd/oft2/cmd](cmd/oft2/cmd/) and provide a callback function for fuzzing the data. Please keep in mind that the fuzzer gets the raw data which is afterwards directly written to the TCP socket, so it is up to you to maintain the right length and header information if you change the data.

## Limitations

  * The tool provides only client functionality, i.e. sending files. Receiving files is not implemented yet, but the messages are already understood by the tool.
  * The command line tool does note provide commands for authentication but the library already supports the necessary protocol steps.
