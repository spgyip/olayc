package olayc

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

type ErrorHandling int

const (
	internalFlagNameYamlFile = "olayc.file.yaml"
	internalFlagNameJsonFile = "olayc.file.json"
)

// Print OlayConfig usage.
func usage() {
	fmt.Println("Usage of olayc:")
	fmt.Println("  -olayc.help int")
	fmt.Println("       Print this help message and call os.Exit(0).")
	fmt.Println("  -olayc.silent bool")
	fmt.Println("       Turn silent mode on/off. Default is false.")
	fmt.Println("  -olayc.file.yaml string")
	fmt.Println("       Load yaml file.")
	fmt.Println("  -olayc.json.yaml string")
	fmt.Println("       Load json file.")
}

// OlayConfig is composition of multiple configure sources, each source is overlayed from bottom to top.
// The top layer is visible if there is key conflicted among layers.
// The configure sources can be configure files, environments and commandline arguments.
type OlayConfig struct {
	merged map[any]any

	silent bool // If silent is true, there will be no verbose logs.
}

// NewOlayConfig allocates and returns a new OlayConfig.
func NewOlayConfig() *OlayConfig {
	return &OlayConfig{
		merged: make(map[any]any),
		silent: false,
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
// - Commandline arguments, e.g. --foo.name=foo
// - Enviroments, e.g. FOO_NAME=hello
// - Yaml files, specified by commandline `--olay.file.yaml=foo.yaml`
// - Json files, specified by commandline `--olay.file.json=foo.json`
//
// If encounter errors, e.g. load file fail, error message will be printed and call os.Exit(1).
func Load() {
	flgs := &flags{}
	flgs.parse(os.Args[1:])

	var yamlFiles []string
	var help = false
	for _, kv := range flgs.kvs {
		if kv.key == "olayc.silent" {
			if kv.value == true {
				defaultC.silent = true
			} else {
				defaultC.silent = false
			}
		} else if kv.key == "olayc.help" {
			help = true
		} else if kv.key == internalFlagNameYamlFile {
			yamlFiles = append(yamlFiles, kv.value.(string))
		}
	}

	if help {
		usage()
		os.Exit(0)
	}

	if !defaultC.silent {
		fmt.Println("[OlayConfig] Silent mode is off, verbose messages will be printed.")
		fmt.Println("[OlayConfig] Silent mode can be turned on with '-olayc.silent'.")
	}

	for _, file := range yamlFiles {
		err := defaultC.LoadYamlFile(file)
		if err != nil {
			fmt.Printf("[OlayConfig] Load fail, error: %v\n", err)
			os.Exit(1)
		}
		if !defaultC.silent {
			fmt.Printf("[OlayConfig] Loaded yaml file %v.\n", file)
		}
	}
}

// Get from defaultC.
func Get(key string) Value {
	return defaultC.Get(key)
}
