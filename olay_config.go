package olayc

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// Print OlayConfig usage.
func usage() {
	fmt.Println("Usage of olayc:")
	for _, fn := range internalFlags {
		fmt.Printf("  %v | %v\n", fn.full, fn.short)
		fmt.Printf("         %v\n", fn.help)
	}
}

// OlayConfig is composition of multiple configure sources, each source is overlayed from bottom to top.
// The top layer is visible if there is key conflicted among layers.
// The configure sources can be configure files, environments and commandline arguments.
type OlayConfig struct {
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
	var err error
	var data []byte
	var m = make(map[any]any)

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

// Get value with the given key, return nil if doesn't exist.
// The key is splitted by seperator '.', e.g. 'foo.name'.
// The key is case sensitive, thus, 'foo.Name' is different from 'foo.name'.
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
// - Commandline arguments, e.g. -foo.name=foo
// - Enviroments, e.g. FOO_NAME=hello
// - Yaml files, e.g. `-oc.f.y=foo.yaml`
// - Json files, e.g. `-oc.f.j=foo.json`
//
// If encounter errors, e.g. load file fail, error message will be printed and call os.Exit(1).
func Load() {
	var yamlFiles []string
	var help = false
	var silent = false

	fp := &flagParser{}
	fp.parse(os.Args[1:])
	for _, kv := range fp.kvs {
		if internalFlags["silent"].is(kv.key) {
			if kv.value == true {
				silent = true
			} else {
				silent = false
			}
		} else if internalFlags["help"].is(kv.key) {
			help = true
		} else if internalFlags["file.yaml"].is(kv.key) {
			yamlFiles = append(yamlFiles, kv.value.(string))
		} else if kv.key[:3] == "oc." {
			fmt.Printf("[OlayConfig] Unknown oc flag: %v\n", kv.key)
			usage()
			os.Exit(1)
		}
	}

	if help {
		usage()
		os.Exit(0)
	}

	if !silent {
		fmt.Println("[OlayConfig] Silent mode is off, verbose messages will be printed.")
		fmt.Printf("[OlayConfig] Silent mode can be turned on with '-%v'.\n", internalFlags["silent"].short)
	}

	for _, file := range yamlFiles {
		err := defaultC.LoadYamlFile(file)
		if err != nil {
			fmt.Printf("[OlayConfig] Load fail, error: %v\n", err)
			os.Exit(1)
		}
		if !silent {
			fmt.Printf("[OlayConfig] Loaded yaml file %v.\n", file)
		}
	}
}

// Get value with the default OlaycConfig .
func Get(key string) Value {
	return defaultC.Get(key)
}
