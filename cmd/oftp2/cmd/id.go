package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thomsmits/oftp2-client/internal/liboftp2/client"
)

var idCommand = &cobra.Command{
	Use:     "id",
	Short:   "Try to determine the server's ODETTE id",
	Long:    `Connects to the server with the id LOCAL (or the id set via -i flag) and tries to find the server's id.`,
	Example: `oftp2 id`,
	Run: func(cmd *cobra.Command, args []string) {
		determineId()
	},
}

func determineId() {

	r := client.OFTP2Client{
		ServerHost: activeOptions.Server,
		ServerPort: activeOptions.Port,
		OdetteId:   activeOptions.OdetteId,
		Verbose:    activeOptions.Verbose,
	}

	ssid, err := r.QueryServerCapabilities()
	if err != nil {
		print(err.Error() + "\n")
		os.Exit(1)
	}

	r.Close()

	fmt.Printf("Server's id is: '%s'\n", ssid.Id)
}
