package olayc

import (
	"flag"

	"github.com/pkg/errors"
)

// cmdFiles is type of `flags.Value`, which is parsed from commandline arguments.
type cmdFiles []string

// newCmdFiles allocates and returns a new cmdFiles.
func newCmdFiles() *cmdFiles {
	return &cmdFiles{}
}

// Parse files from arguments
func (l *cmdFiles) parseFromArgs(args []string, name string, usage string) error {
	fs := flag.NewFlagSet("olayc", flag.ContinueOnError)
	fs.Var(l, name, usage)
	if err := fs.Parse(args); err != nil {
		errors.Wrap(err, "Parse cmdfiles from args error")
	}
	return nil
}

// Implement `flags.Value`'s `String` method.
func (l *cmdFiles) String() string {
	return "OlayConfig cmdfiles"
}

// Implement `flags.Value`'s `Set`method.
func (l *cmdFiles) Set(v string) error {
	*l = append(*l, v)
	return nil
}

func parseArguments(args []string) {

}
