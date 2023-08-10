package olayc

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

const (
	// Root key
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

type loadOptionFunc func(*loadOptions)
type loadOptions struct {
	filesRequired []string
}

func WithFileRequire(name string) loadOptionFunc {
	return func(opt *loadOptions) {
		opt.filesRequired = append(opt.filesRequired, name)
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
	return c.LoadYaml(data)
}

// Load yaml from bytes.
func (c *OlayConfig) LoadYaml(data []byte) error {
	var m = make(map[any]any)
	err := yaml.Unmarshal(data, &m)
	if err != nil {
		return errors.Wrap(err, "LoadYamlBytes error")
	}
	copyMap(c.merged, m)
	return nil
}

// Load from arguments.
// Return numbers of kvs loaded.
// The internal olayc flags which prefix with `-oc.|--oc.` are ignored.
func (c *OlayConfig) LoadArgs(args []string) (int, error) {
	var kvs []KV
	fp := &flagParser{}
	fp.parse(args)
	for _, kv := range fp.kvs {
		if strings.HasPrefix(kv.key, internalFlagPrefix) {
			continue
		}
		kvs = append(kvs, kv)
	}
	return c.LoadKVs(kvs)
}

// Load from key-value pairs.
// Return number of kvs loaded.
func (c *OlayConfig) LoadKVs(kvs []KV) (int, error) {
	var m = make(map[any]any)
	for _, kv := range kvs {
		var cur any = m
		sps := strings.Split(kv.key, ".")
		for j, sp := range sps {
			curM := cur.(map[any]any)
			if j == len(sps)-1 {
				curM[sp] = kv.value
			} else {
				if _, ok := curM[sp]; !ok {
					curM[sp] = make(map[any]any)
				}
				cur = curM[sp]
			}
		}
	}
	copyMap(c.merged, m)
	return len(kvs), nil
}

// Get value with the given key, return nil if doesn't exist.
// The key is splitted by seperator '.', e.g. 'foo.name'.
// The key is case sensitive, thus, 'foo.Name' is different from 'foo.name'.
// Use `Root` key to get the whole configure.
// Return nil if it doesn't exist.
// TODO: How if the value is set to nil, should tell the difference between not-exist and nil value?
func (c *OlayConfig) Get(key string) *Value {
	var cur any = c.merged
	if key == Root {
		return &Value{v: cur}
	}
	sps := strings.Split(key, ".")
	for _, sp := range sps {
		var ok bool
		var curM map[any]any
		if curM, ok = cur.(map[any]any); !ok {
			cur = nil
			break
		}
		if cur, ok = curM[sp]; !ok {
			cur = nil
			break
		}
	}
	if cur == nil {
		return nil
	}
	return &Value{v: cur}
}

// Get string value, return defaultValue if it doesn't exisit.
func (c *OlayConfig) String(key string, defaultValue string) string {
	v := c.Get(key)
	if v == nil {
		return defaultValue
	}
	return v.String()
}

// Get int value, return defaultValue if it doesn't exisit.
func (c *OlayConfig) Int(key string, defaultValue int) int {
	v := c.Get(key)
	if v == nil {
		return defaultValue
	}
	return v.Int()
}

// Get uint value, return defaultValue if it doesn't exisit.
func (c *OlayConfig) Uint(key string, defaultValue uint) uint {
	v := c.Get(key)
	if v == nil {
		return defaultValue
	}
	return v.Uint()
}

// Get int64 value, return defaultValue if it doesn't exisit.
func (c *OlayConfig) Int64(key string, defaultValue int64) int64 {
	v := c.Get(key)
	if v == nil {
		return defaultValue
	}
	return v.Int64()
}

// Get uint64 value, return defaultValue if it doesn't exisit.
func (c *OlayConfig) Uint64(key string, defaultValue uint64) uint64 {
	v := c.Get(key)
	if v == nil {
		return defaultValue
	}
	return v.Uint64()
}

// Get float64 value, return defaultValue if it doesn't exisit.
func (c *OlayConfig) Float64(key string, defaultValue float64) float64 {
	v := c.Get(key)
	if v == nil {
		return defaultValue
	}
	return v.Float64()
}

// Get bool value, return defaultValue if it doesn't exisit.
func (c *OlayConfig) Bool(key string, defaultValue bool) bool {
	v := c.Get(key)
	if v == nil {
		return defaultValue
	}
	return v.Bool()
}

// Unmarshal out, return error if it doesn't exist.
func (c *OlayConfig) Unmarshal(key string, out any) error {
	v := c.Get(key)
	if v == nil {
		return errors.Errorf("key doesn't exists: %v", key)
	}
	return v.Unmarshal(out)
}

// `defaultC` is the default OlayConfig.
var defaultC = New()

// Load the default OlayConfig from configure sources:
// - Commandline arguments, e.g. -foo.name=foo
// - Enviroments, e.g. FOO_NAME=hello
// - Yaml files, e.g. `-oc.f.y=foo.yaml`
// - Json files, e.g. `-oc.f.j=foo.json`
//
// If errors happen, e.g. load file fail, error message will be printed and call os.Exit(1).
func Load(opts ...loadOptionFunc) {
	var yamlFiles []string
	var help = false
	var silent = false

	var opt loadOptions
	for _, of := range opts {
		of(&opt)
	}

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
		} else if strings.HasPrefix(kv.key, internalFlagPrefix) {
			fmt.Printf("[OlayConfig] Unknown oc flag: %v\n", kv.key)
			usage()
			os.Exit(1)
		}
	}

	if help {
		usage()
		os.Exit(0)
	}

	// Check required files
	var fileCheck = true
	for _, fr := range opt.filesRequired {
		var ok = false
		for _, fy := range yamlFiles {
			// Check if fy has suffix of fr.
			n := len(fy) - len(fr)
			if n >= 0 && fy[n:] == fr {
				ok = true
				break
			}
		}
		if !ok {
			fmt.Printf("[OlayConfig] File %v is required.\n", fr)
			fileCheck = false
		}
	}
	if !fileCheck {
		fmt.Println("[OlayConfig] Use '-oc.f.(y|j)=....'")
		os.Exit(1)
	}

	if !silent {
		fmt.Println("[OlayConfig] Silent mode is off, verbose messages will be printed.")
		fmt.Printf("[OlayConfig] Silent mode can be turned on with '-%v'.\n", internalFlags["silent"].short)
	}

	// Priority
	//  - Commandline arguments
	//  - Yaml/Json files
	n, err := defaultC.LoadArgs(os.Args[1:])
	if err != nil {
		fmt.Printf("[OlayConfig] Load arguments fail, error: %v]\n", err)
		os.Exit(1)
	}
	if !silent {
		fmt.Printf("[OlayConfig] Commandlines loaded, totally %v KVs.\n", n)
	}

	for _, file := range yamlFiles {
		err := defaultC.LoadYamlFile(file)
		if err != nil {
			fmt.Printf("[OlayConfig] Load fail, error: %v\n", err)
			os.Exit(1)
		}
		if !silent {
			fmt.Printf("[OlayConfig] Yaml file loaded: %v.\n", file)
		}
	}
}

// Get value with default OlaycConfig.
func Get(key string) *Value {
	return defaultC.Get(key)
}

// Get string with default OlayConfig.
func String(key string, defaultValue string) string {
	return defaultC.String(key, defaultValue)
}

// Get int with default OlayConfig.
func Int(key string, defaultValue int) int {
	return defaultC.Int(key, defaultValue)
}

// Get uint with default OlayConfig.
func Uint(key string, defaultValue uint) uint {
	return defaultC.Uint(key, defaultValue)
}

// Get int64 with default OlayConfig.
func Int64(key string, defaultValue int64) int64 {
	return defaultC.Int64(key, defaultValue)
}

// Get uint64 with default OlayConfig.
func Uint64(key string, defaultValue uint64) uint64 {
	return defaultC.Uint64(key, defaultValue)
}

// Get float64 with default OlayConfig.
func Float64(key string, defaultValue float64) float64 {
	return defaultC.Float64(key, defaultValue)
}

// Get bool with default OlayConfig.
func Bool(key string, defaultValue bool) bool {
	return defaultC.Bool(key, defaultValue)
}

// Unmarshal with default OlayConfig.
func Unmarshal(key string, out any) error {
	return defaultC.Unmarshal(key, out)
}
