package main

import (
	"github.com/thomsmits/oftp2-client/cmd/oftp2/cmd"
)

// Main function of the tool.
func main() {

	// As we are using the cobra framework for command line tools, this method is
	// more or less empty and we delegate the whole option and command line handling
	// to the framework.
	//
	// To add a new command, go to the cmd folder, create a new file and instantiate
	// a cobra.Command structure. Then add this structure to the available commands
	// in the cmd/root.go file (init method).
	cmd.Execute()
}
