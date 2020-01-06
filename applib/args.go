package applib

import (
	"flag"
)

// ProgramArgs contains console arguments
type ProgramArgs struct {

	Port    int				// Port to listen
	Bind    bool			// IP address to bind

}

// getArguments parse command line arguments
// enables setting custom listening port 
func getArguments() *ProgramArgs {

	args := new(ProgramArgs)
	portPtr := flag.Int("port", 5580, "Service port\n")
	bindPtr := flag.Bool("bind", false, "Bind all interfaces\n")

	flag.Parse()

	args.Port = *portPtr
	args.Bind = *bindPtr

	return args
}
