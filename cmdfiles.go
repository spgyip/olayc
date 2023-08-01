package olayc

import (
	"flag"

	"github.com/pkg/errors"
)

const (
	internalFlagFileYaml = "olayc.file.yaml"
)

// cmdFiles is type of `flags.Value`, which is parsed from commandline arguments
type cmdFiles []string

// Implement `flags.Value`'s `String` method
func (l *cmdFiles) String() string {
	return "OlayConfig cmdfiles"
}

// Implement `flags.Value`'s `Set`method
func (l *cmdFiles) Set(v string) error {
	*l = append(*l, v)
	return nil
}

// Parse cmdFiles from args
func parseCmdFilesFromArgs(args []string) (cmdFiles, error) {
	var files cmdFiles
	fs := flag.NewFlagSet("none", flag.ContinueOnError)
	fs.Var(&files, internalFlagFileYaml, "OlayConfig yaml file")
	if err := fs.Parse(args); err != nil {
		errors.Wrap(err, "Parse cmdfiles from args error")
	}
	return files, nil
}
