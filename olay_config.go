package olayc

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

const (
	internalFlagNameYamlFile = "olayc.file.yaml"
	internalFlagNameJsonFile = "olayc.file.json"
)

// OlayConfig is composition of multiple configure sources, each source is overlayed from bottom to top.
// The top layer is visible if there is key conflicted among layers.
// The configure sources can be configure files, environments and commandline arguments.
type OlayConfig struct {
	merged map[any]any

	errorHandling int
}

// NewOlayConfig allocates and returns a new OlayConfig.
func NewOlayConfig() *OlayConfig {
	return &OlayConfig{
		merged: make(map[any]any),
	}
}

// Load yaml config from file, stack on top layer
func (c *OlayConfig) LoadYamlFile(filepath string) error {
	var err error
	var data []byte
	var m = make(map[any]any)

	defer func() {
		if err != nil {
			fmt.Printf("[OlayConfig] %v loaded error: %v\n", filepath, err)
		} else {
			fmt.Printf("[OlayConfig] %v loaded.\n", filepath)
		}
	}()

	data, err = os.ReadFile(filepath)
	if err != nil {
		return errors.Wrap(err, "LoadYamlFile error")
	}

	err = yaml.Unmarshal(data, &m)
	if err != nil {
		return errors.Wrap(err, "LoadYamlFile error")
	}

	copyMap(c.merged, m)
	return nil
}

// Implement `ConfigSource`'s method.
// Return nil if key doesn't exist.
func (c *OlayConfig) Get(key string) Value {
	var cur any = c.merged
	sps := strings.Split(key, ".")
	for _, sp := range sps {
		if reflect.TypeOf(cur).Kind() != reflect.Map {
			return nil
		}
		next, ok := cur.(map[any]any)[sp]
		if !ok {
			return nil
		}
		cur = next
	}
	return cur
}

// `defaultC` is the default OlayConfig.
var defaultC = NewOlayConfig()

// Load the default OlayConfig from configure sources:
// - Commandline arguments, e.g. --foo.name=foo
// - Enviroments, e.g. FOO_NAME=hello
// - Yaml files, specified by commandline `--olay.file.yaml=foo.yaml`
// - Json files, specified by commandline `--olay.file.json=foo.json`
func Load() error {
	files := newCmdFiles()
	err := files.parseFromArgs(os.Args[1:], internalFlagNameYamlFile, "OlayConfig internal, load yaml file")
	if err != nil {
		return errors.Wrap(err, "Load error")
	}
	for _, file := range *files {
		defaultC.LoadYamlFile(file)
	}
	return nil
}

// Get from defaultC.
func Get(key string) Value {
	return defaultC.Get(key)
}
