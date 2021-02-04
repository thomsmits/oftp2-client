package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

const rootDoc = `
An OFTP2 client implemented in Go.`

var rootCmd = &cobra.Command{
	Use:   "oftp2",
	Short: "oftp2 is an OFTP2 client",
	Long:  rootDoc,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&activeOptions.OdetteId, "odetteId", "i", "LOCAL", "Odette ID of this client")
	rootCmd.PersistentFlags().IntVarP(&activeOptions.Port, "port", "p", 3305, "Port of the Odette server")
	rootCmd.PersistentFlags().StringVarP(&activeOptions.Server, "host", "s", "localhost", "host to connect to")
	rootCmd.PersistentFlags().BoolVarP(&activeOptions.Verbose, "verbose", "v", false, "verbose output")

	rootCmd.AddCommand(queryCommand)
	rootCmd.AddCommand(sendCommand)
	rootCmd.AddCommand(idCommand)
}

// Execute the command.
func Execute() {
	err := rootCmd.Execute()

	if err != nil {
		print(err.Error() + "\n")
		os.Exit(1)
	}
}
