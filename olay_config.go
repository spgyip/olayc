package olayc

import (
	"log"
	"os"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// OlayConfig is composition of multiple configure sources, each source is overlayed from bottom to top.
// The top layer is visible if there is key conflicted among layers.
// The configure sources can be configure files, environments and commandline arguments.
type OlayConfig struct {
	layers []ConfigSource
	merged map[any]any
}

// NewOlayConfig allocates and returns a new OlayConfig.
func NewOlayConfig() *OlayConfig {
	return &OlayConfig{
		merged: make(map[any]any),
	}
}

// Load yaml config from file, stack on top layer
func (c *OlayConfig) LoadYamlFile(filepath string) error {
	log.Println("[OlayConfig] Load yaml file: ", filepath)
	data, err := os.ReadFile(filepath)
	if err != nil {
		return errors.Wrap(err, "LoadYamlFile error")
	}

	var m = make(map[any]any)
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

// Load the default OlayConfig from configure sources:
//   Commandline arguments, e.g. --foo.name=foo
//   Enviroments, e.g. FOO_NAME=hello
//   Yaml files, specified by commandline `--olay.file.yaml=foo.yaml`
//   Json files, specified by commandline `--olay.file.json=foo.json`
func Load() error {
	files, err := parseCmdFilesFromArgs(os.Args[1:])
	if err != nil {
		return errors.Wrap(err, "Load error")
	}
	for _, file := range files {
		defaultC.LoadYamlFile(file)
	}
	return nil
}
