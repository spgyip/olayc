package olayc

import (
	"encoding/json"
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

// KV is composition of key and value.
type KV struct {
	key   string
	value any
}

// Print OlayConfig usage.
func usage() {
	fmt.Println("Usage of olayc:")
	for _, fn := range internalFlags {
		fmt.Printf("  -%v | -%v\n", fn.full, fn.short)
		fmt.Printf("         %v\n", fn.help)
	}
}

// LoadOptions set options for `Load()`.
type loadOptionFunc func(*loadOptions)
type loadOptions struct {
	filesRequired []string
}

// WithFileRequire returns a loadOptionFunc appends a required file.
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

// Load yaml config from file.
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
		return errors.Wrap(err, "LoadYaml error")
	}
	copyMap(c.merged, m)
	return nil
}

// Load json config from file.
func (c *OlayConfig) LoadJsonFile(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return errors.Wrap(err, "LoadJsonFile error")
	}
	return c.LoadJson(data)
}

// Load json from bytes.
// To unmarshal a JSON object into a map using the standard library "encoding/json",
// the map's key type must either be any string type, an integer.
// Thus, the unmarshal map type is map[string]any(and all sub-maps), it not compatible with `copyMap()` which is accepting type map[any]any.
// We must convert `map[string]any` to `map[any]any`, this is simplily done by marshal/unmarshal with "gopkg.in/yaml.v2".
func (c *OlayConfig) LoadJson(data []byte) error {
	var m = make(map[string]any)
	err := json.Unmarshal(data, &m)
	if err != nil {
		return errors.Wrap(err, "LoadJson error")
	}

	var m1 map[any]any
	m1, err = convertMap(m)
	if err != nil {
		return errors.Wrap(err, "LoadJson error")
	}

	copyMap(c.merged, m1)
	return nil
}

// Load from arguments. Return numbers of kvs loaded.
// The internal olayc flags which prefix with `-oc.|--oc.` are ignored.
//
// If there are overlap keys, refer to 'LoadKVs()'.
func (c *OlayConfig) LoadArgs(args []string) (int, error) {
	var kvs []KV
	psr := &flagParser{}
	psr.parse(args)
	for _, kv := range psr.kvs {
		if strings.HasPrefix(kv.key, internalFlagPrefix) {
			continue
		}
		kvs = append(kvs, kv)
	}
	return c.LoadKVs(kvs)
}

// Load from environments. Return numbers of kvs loaded.
//
// The key is converted to lower case and the seperator '_' is replaced by '.'.
// E.g. 'LC_CTYPE=UTF-8', is converted to 'lc.ctype=UTF-8'.
// The anterior '_' in keys are trimed, e.g. '_P9K_SSH_TTY' is converted to `p9k.ssh.tty`.
//
// If there are overlap envs, e.g. 'TERM=tmux' 'TERM_PROGRAM=tmux', refer to 'LoadKVs()'.
func (c *OlayConfig) LoadEnvs(envs []string) (int, error) {
	psr := &envParser{}
	psr.parse(envs)
	return c.LoadKVs(psr.kvs)
}

// Load from key-value pairs. Return number of kvs loaded.
//
// If there are overlap keys, e.g. key1 'foo.redis=redis.cluster' and key2 'foo.redis.host=redis.cluster'.
// As seen, key2 contains the key1. If get with key 'foo.redis', only one value will be returned, either of 'redis.cluster' and '{"host": "redis.cluster"}'.
// The previously loaded key is more prior than the latter ones, so the latter one is ignored.
// For example, if 'foo.redis' is loaded previously, the return value is 'redis.cluster',
// or if the 'foo.redis.host' is loaded previously, the return value is '{"host": "redis.cluster"}'.
func (c *OlayConfig) LoadKVs(kvs []KV) (int, error) {
	var m = make(map[any]any)
	for _, kv := range kvs {
		var cur any = m
		sps := strings.Split(kv.key, ".")
		for j, sp := range sps {
			var curM map[any]any
			var ok bool
			// Current node is scalar value
			if curM, ok = cur.(map[any]any); !ok {
				break
			}

			// Add subtree or value if empty
			if _, ok = curM[sp]; !ok {
				if j == len(sps)-1 {
					curM[sp] = kv.value
				} else {
					curM[sp] = make(map[any]any)
				}
			}
			cur = curM[sp]
		}
	}
	copyMap(c.merged, m)
	return len(kvs), nil
}

// Get value with the given key, return nil if doesn't exist.
// The key is splitted by seperator '.', e.g. 'foo.name'.
// The key is case sensitive, thus, 'foo.Name' is different from 'foo.name'.
// Use `Root` key to get the whole configure.
// If it doesn't exist, 'Value.IsNil()' is true.
// TODO: How if the value is set to nil, should tell the difference between not-exist and nil value?
func (c *OlayConfig) Get(key string) Value {
	var cur any = c.merged
	if key == Root {
		return Value{v: cur}
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
	return Value{v: cur}
}

// Get string value, return defaultValue if it doesn't exisit.
func (c *OlayConfig) String(key string, defaultValue string) string {
	v := c.Get(key)
	if v.IsNil() {
		return defaultValue
	}
	return v.String()
}

// Get int value, return defaultValue if it doesn't exisit.
func (c *OlayConfig) Int(key string, defaultValue int) int {
	v := c.Get(key)
	if v.IsNil() {
		return defaultValue
	}
	return v.Int()
}

// Get uint value, return defaultValue if it doesn't exisit.
func (c *OlayConfig) Uint(key string, defaultValue uint) uint {
	v := c.Get(key)
	if v.IsNil() {
		return defaultValue
	}
	return v.Uint()
}

// Get int64 value, return defaultValue if it doesn't exisit.
func (c *OlayConfig) Int64(key string, defaultValue int64) int64 {
	v := c.Get(key)
	if v.IsNil() {
		return defaultValue
	}
	return v.Int64()
}

// Get uint64 value, return defaultValue if it doesn't exisit.
func (c *OlayConfig) Uint64(key string, defaultValue uint64) uint64 {
	v := c.Get(key)
	if v.IsNil() {
		return defaultValue
	}
	return v.Uint64()
}

// Get float64 value, return defaultValue if it doesn't exisit.
func (c *OlayConfig) Float64(key string, defaultValue float64) float64 {
	v := c.Get(key)
	if v.IsNil() {
		return defaultValue
	}
	return v.Float64()
}

// Get bool value, return defaultValue if it doesn't exisit.
func (c *OlayConfig) Bool(key string, defaultValue bool) bool {
	v := c.Get(key)
	if v.IsNil() {
		return defaultValue
	}
	return v.Bool()
}

// Unmarshal out, return error if it doesn't exist.
func (c *OlayConfig) Unmarshal(key string, out any) error {
	v := c.Get(key)
	if v.IsNil() {
		return errors.Errorf("key doesn't exists: %v", key)
	}
	return v.Unmarshal(out)
}

// Return Yaml bytes.
func (c *OlayConfig) ToYaml() (string, error) {
	v := c.Get(Root)
	data, err := v.MarshalToYaml()
	if err != nil {
		return "", err
	}
	return string(data), nil
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
	type inputFileType int
	const (
		Yaml inputFileType = iota
		Json
	)
	type inputFile struct {
		name string
		typ  inputFileType
	}

	var help = false
	var silent = false
	var dryrun = false
	var files []inputFile

	var opt loadOptions
	for _, of := range opts {
		of(&opt)
	}

	psr := &flagParser{}
	psr.parse(os.Args[1:])
	for _, kv := range psr.kvs {
		if internalFlags["silent"].is(kv.key) {
			if kv.value == true {
				silent = true
			} else {
				silent = false
			}
		} else if internalFlags["help"].is(kv.key) {
			help = true
		} else if internalFlags["dryrun"].is(kv.key) {
			dryrun = true
		} else if internalFlags["file.yaml"].is(kv.key) {
			files = append(files, inputFile{kv.value.(string), Yaml})
		} else if internalFlags["file.json"].is(kv.key) {
			files = append(files, inputFile{kv.value.(string), Json})
		} else if strings.HasPrefix(kv.key, internalFlagPrefix) {
			fmt.Printf("[OlayConfig][Error] Unknown oc flag: %v\n", kv.key)
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
		fmt.Printf("[OlayConfig] Dry run mode is %v.\n", dryrun)
	}

	if len(opt.filesRequired) > 0 && !silent {
		fmt.Println("[OlayConfig] Files required:")
		for _, name := range opt.filesRequired {
			fmt.Printf("  - %v\n", name)
		}
	}

	// Check required files
	var fileCheck = true
	for _, fr := range opt.filesRequired {
		var ok = false
		for _, fy := range files {
			// Check if fy has suffix of fr.
			n := len(fy.name) - len(fr)
			if n >= 0 && fy.name[n:] == fr {
				ok = true
				break
			}
		}
		if !ok {
			fmt.Printf("[OlayConfig][Error] File %v is required.\n", fr)
			fileCheck = false
		}
	}
	if !fileCheck {
		fmt.Println("[OlayConfig] Use '-oc.f.(y|j)=....'")
		os.Exit(1)
	}

	// Load commandline arguments
	n, err := defaultC.LoadArgs(os.Args[1:])
	if err != nil {
		fmt.Printf("[OlayConfig][Error] Load arguments fail, error: %v]\n", err)
		os.Exit(1)
	}
	if !silent {
		fmt.Printf("[OlayConfig] Commandlines loaded, totally %v KVs.\n", n)
	}

	// Load ENVs
	n, err = defaultC.LoadEnvs(os.Environ())
	if err != nil {
		if err != nil {
			fmt.Printf("[OlayConfig][Error] Load environments fail, error: %v]\n", err)
			os.Exit(1)
		}
	}
	if !silent {
		fmt.Printf("[OlayConfig] Environments loaded, totally %v KVs.\n", n)
	}

	// Load yaml/json files
	for _, f := range files {
		var err error
		if f.typ == Yaml {
			err = defaultC.LoadYamlFile(f.name)

		} else if f.typ == Json {
			err = defaultC.LoadJsonFile(f.name)
		}
		if err != nil {
			fmt.Printf("[OlayConfig][Error] Load fail, error: %v\n", err)
			os.Exit(1)
		}
		if !silent {
			fmt.Printf("[OlayConfig] File loaded: %v.\n", f.name)
		}
	}

	if dryrun {
		fmt.Println("[OlayConfig] Dry run mode is on, program will exit after yaml printed.")
		fmt.Println(defaultC.ToYaml())
		os.Exit(0)
	}
}

// Get value with default OlaycConfig.
func Get(key string) Value {
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

// ToYaml with default OlayConfig.
func ToYaml() (string, error) {
	return defaultC.ToYaml()
}
