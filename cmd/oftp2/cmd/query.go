package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thomsmits/oftp2-client/internal/liboftp2/client"
)

var queryCommand = &cobra.Command{
	Use:     "query",
	Short:   "Query the server",
	Long:    `Queries the server to determine its capabilities.`,
	Example: `oftp2 query -i LOCAL`,
	Run: func(cmd *cobra.Command, args []string) {
		queryClient()
	},
}

func queryClient() {

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

	fmt.Printf("\nData received from remote system:\n")
	fmt.Printf("------------------------------------------\n")
	fmt.Printf("%s", ssid.String())
	fmt.Printf("------------------------------------------\n")
	r.Close()

	os.Exit(0)
}
