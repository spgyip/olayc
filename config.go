package olayc

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

const (
	// Root value
	Root = ""
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

// New allocates and returns a new OlayConfig.
func New() *OlayConfig {
	return &OlayConfig{
		merged: make(map[any]any),
	}
}

// Load yaml config from file, stack on top layer
func (c *OlayConfig) LoadYamlFile(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return errors.Wrap(err, "LoadYamlFile error")
	}
	return c.LoadYamlBytes(data)
}

// Load yaml from bytes.
func (c *OlayConfig) LoadYamlBytes(data []byte) error {
	var m = make(map[any]any)
	err := yaml.Unmarshal(data, &m)
	if err != nil {
		return errors.Wrap(err, "LoadYamlBytes error")
	}
	copyMap(c.merged, m)
	return nil
}

// Get value with the given key, return nil if doesn't exist.
// The key is splitted by seperator '.', e.g. 'foo.name'.
// The key is case sensitive, thus, 'foo.Name' is different from 'foo.name'.
// Return nil if it doesn't exist.
func (c *OlayConfig) Get(key string) Value {
	var cur any = c.merged
	if key == Root {
		return cur
	}
	sps := strings.Split(key, ".")
	for _, sp := range sps {
		next, ok := cur.(map[any]any)[sp]
		if !ok {
			return nil
		}
		cur = next
	}
	return cur
}

// Get string value, return defaultValue if it doesn't exist.
func (c *OlayConfig) String(key string, defaultValue string) string {
	v := c.Get(key)
	if v == nil {
		return defaultValue
	}

	var s = defaultValue
	switch x := v.(type) {
	case string:
		s = x
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64,
		bool:
		s = fmt.Sprintf("%v", x)
	}
	return s
}

// Get int value, return defaultValue if it doesn't exist.
func (c *OlayConfig) Int(key string, defaultValue int) int {
	v := c.Get(key)
	if v == nil {
		return defaultValue
	}

	var i int
	switch x := v.(type) {
	case int:
		i = int(x)
	case int8:
		i = int(x)
	case int16:
		i = int(x)
	case int32:
		i = int(x)
	case int64:
		i = int(x)
	case uint:
		i = int(x)
	case uint8:
		i = int(x)
	case uint16:
		i = int(x)
	case uint32:
		i = int(x)
	case uint64:
		i = int(x)
	default:
		return defaultValue
	}
	return i
}

// Get uint value, return defaultValue if it doesn't exist.
func (c *OlayConfig) Uint(key string, defaultValue uint) uint {
	v := c.Get(key)
	if v == nil {
		return defaultValue
	}

	var i = defaultValue
	switch x := v.(type) {
	case int:
		i = uint(x)
	case int8:
		i = uint(x)
	case int16:
		i = uint(x)
	case int32:
		i = uint(x)
	case int64:
		i = uint(x)
	case uint:
		i = uint(x)
	case uint8:
		i = uint(x)
	case uint16:
		i = uint(x)
	case uint32:
		i = uint(x)
	case uint64:
		i = uint(x)
	}
	return i
}

// Get int64 value, return defaultValue if it doesn't exist.
func (c *OlayConfig) Int64(key string, defaultValue int64) int64 {
	v := c.Get(key)
	if v == nil {
		return defaultValue
	}

	var i = defaultValue
	switch x := v.(type) {
	case int:
		i = int64(x)
	case int8:
		i = int64(x)
	case int16:
		i = int64(x)
	case int32:
		i = int64(x)
	case int64:
		i = int64(x)
	case uint:
		i = int64(x)
	case uint8:
		i = int64(x)
	case uint16:
		i = int64(x)
	case uint32:
		i = int64(x)
	case uint64:
		i = int64(x)
	}
	return i
}

// Get uint64 value, return defaultValue if it doesn't exist.
func (c *OlayConfig) Uint64(key string, defaultValue uint64) uint64 {
	v := c.Get(key)
	if v == nil {
		return defaultValue
	}

	var i = defaultValue
	switch x := v.(type) {
	case int:
		i = uint64(x)
	case int8:
		i = uint64(x)
	case int16:
		i = uint64(x)
	case int32:
		i = uint64(x)
	case int64:
		i = uint64(x)
	case uint:
		i = uint64(x)
	case uint8:
		i = uint64(x)
	case uint16:
		i = uint64(x)
	case uint32:
		i = uint64(x)
	case uint64:
		i = uint64(x)
	}
	return i
}

// Get float64 value, return defaultValue if it doesn't exist.
func (c *OlayConfig) Float64(key string, defaultValue float64) float64 {
	v := c.Get(key)
	if v == nil {
		return defaultValue
	}

	var i = defaultValue
	switch x := v.(type) {
	case float32:
		i = float64(x)
	case float64:
		i = float64(x)
	}
	return i
}

// Get bool value, return defaultValue if it doesn't exist.
func (c *OlayConfig) Bool(key string, defaultValue bool) bool {
	v := c.Get(key)
	if v == nil {
		return defaultValue
	}

	var i = defaultValue
	switch x := v.(type) {
	case bool:
		i = bool(x)
	}
	return i
}

// Unmarshal is implemented by using yaml utility,
// value is firstly marshalled to yaml bytes,
// then the yaml bytes is unmarshal to target out.
// Thus, if 'out' is a struct, you must use the yaml struct tag.
func (c *OlayConfig) Unmarshal(key string, out any) error {
	v := c.Get(key)
	if v == nil {
		return errors.New("key doesn't exist")
	}
	return UnmarshalValue(v, out)
}

// `defaultC` is the default OlayConfig.
var defaultC = New()

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

// Get value with default OlaycConfig .
func Get(key string) Value {
	return defaultC.Get(key)
}

// Get string with default OlayConfig
func String(key string, defaultValue string) string {
	return defaultC.String(key, defaultValue)
}

// Get int with default OlayConfig
func Int(key string, defaultValue int) int {
	return defaultC.Int(key, defaultValue)
}

// Get uint with default OlayConfig
func Uint(key string, defaultValue uint) uint {
	return defaultC.Uint(key, defaultValue)
}

// Get int64 with default OlayConfig
func Int64(key string, defaultValue int64) int64 {
	return defaultC.Int64(key, defaultValue)
}

// Get uint64 with default OlayConfig
func Uint64(key string, defaultValue uint64) uint64 {
	return defaultC.Uint64(key, defaultValue)
}

// Get float64 with default OlayConfig
func Float64(key string, defaultValue float64) float64 {
	return defaultC.Float64(key, defaultValue)
}

// Get bool with default OlayConfig
func Bool(key string, defaultValue bool) bool {
	return defaultC.Bool(key, defaultValue)
}

// Unmarshal with default OlayConfig
func Unmarshal(key string, out any) error {
	return defaultC.Unmarshal(key, out)
}
