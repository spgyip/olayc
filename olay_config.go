package olayc

import (
	"log"
	"os"

	"github.com/pkg/errors"
)

// OlayConfig is combination of different configure sources,
// e.g. configure files, environments and commandline arguments.
// Each configure source is overlayed bottom-to-top
// , when there is key conflicted, the upper one will be picked.
type OlayConfig struct {
	// Each layer is appended in order, so the last layer is the top layer, the first layer is the bottom layer.
	layers []ConfigSource
}

// NewOlayConfig allocates and returns a new OlayConfig.
func NewOlayConfig() *OlayConfig {
	return &OlayConfig{
		layers: make([]ConfigSource, 0),
	}
}

// Load yaml config from file, stack on top layer
func (c *OlayConfig) LoadYamlFile(filepath string) error {
	log.Println("[OlayConfig] Load yaml file: ", filepath)
	yc := newYamlConfig()
	err := yc.loadFromFile(filepath)
	if err != nil {
		return errors.Wrap(err, "loadYamlFile error")
	}
	c.layers = append(c.layers, yc)
	return nil
}

// Implement `ConfigSource`'s method.
// Get value from from layers, from upper to layer, return the first found value.
// Return nil if key doesn't exist.
func (c *OlayConfig) Get(key string) Value {
	// The last layer is top layer, so iterate in reverse order.
	for i := len(c.layers) - 1; i >= 0; i-- {
		layer := c.layers[i]
		v := layer.Get(key)
		if v == nil {
			continue
		}
		return v
	}
	return nil
}

// `defaultC` is the default OlayConfig
var defaultC = NewOlayConfig()

// Load the default OlayConfig from the following configure sources:
//   Commandline arguments, e.g. --foo.name=foo
//   Enviroments, e.g. FOO_NAME=hello
//   Yaml files, from commandline `--olay.file.yaml=foo.yaml`
//   Json files, from commandline `--olay.file.json=foo.json`
func LoadDefault() error {
	files, err := parseCmdFilesFromArgs(os.Args[1:])
	if err != nil {
		return errors.Wrap(err, "LoadDefault error")
	}
	for _, file := range files {
		defaultC.LoadYamlFile(file)
	}
	return nil
}
