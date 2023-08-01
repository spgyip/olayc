package olayc

import (
	"os"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// yamlConfig implements a ConfigSoure, which loads configures from yaml file.
type yamlConfig struct {
	m map[any]any
}

// newYamlConfig allocates and returns a new yamlConfig.
func newYamlConfig() *yamlConfig {
	return &yamlConfig{
		m: make(map[any]any),
	}
}

// Load yaml data from file and unmarshal it
func (c *yamlConfig) loadFromFile(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return errors.Wrap(err, "loadFromFile error")
	}
	err = yaml.Unmarshal(data, &c.m)
	if err != nil {
		return errors.Wrap(err, "loadFromFile error")
	}
	return nil
}

// Implement ConfigSource's method.
func (c *yamlConfig) Get(key string) Value {
	var cur any = c.m
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
