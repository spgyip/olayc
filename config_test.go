package olayc

import (
	"reflect"
	"testing"
)

func TestConfigGetValueLoadYaml(t *testing.T) {
	var testdata = []byte(`
foo:
  name: foo1
  id: 123
  pi: 3.1415926
  temp: -50
  onoff: True
`)

	var c = New()
	err := c.LoadYaml(testdata)
	if err != nil {
		t.Fatal(err)
	}

	for i, test := range []struct {
		key    string
		expect any
	}{
		{"foo-not-exisit", Value{v: nil}},
		{"foo.name-not-exisit", Value{v: nil}},
		{"foo.name.not-exisit", Value{v: nil}},
		{"foo.name", Value{string("foo1")}},
		{"foo.id", Value{int(123)}},
		{"foo.pi", Value{float64(3.1415926)}},
		{"foo.onoff", Value{bool(true)}},
	} {
		got := c.Get(test.key)
		if got != test.expect {
			t.Errorf("[%v] got(%v)!=expect(%v)\n", i, got, test.expect)
		}
	}
}

func TestConfigGetValueLoadJson(t *testing.T) {
	var testdata = []byte(`
{
  "foo": {
    "name": "foo1",
    "id": 123,
	"pi": 3.1415926,
	"temp": -50,
	"onoff": true
  }
}
`)

	var c = New()
	err := c.LoadYaml(testdata)
	if err != nil {
		t.Fatal(err)
	}

	for i, test := range []struct {
		key    string
		expect any
	}{
		{"foo-not-exisit", Value{v: nil}},
		{"foo.name-not-exisit", Value{v: nil}},
		{"foo.name.not-exisit", Value{v: nil}},
		{"foo.name", Value{string("foo1")}},
		{"foo.id", Value{int(123)}},
		{"foo.pi", Value{float64(3.1415926)}},
		{"foo.onoff", Value{bool(true)}},
	} {
		got := c.Get(test.key)
		if got != test.expect {
			t.Errorf("[%v] got(%v)!=expect(%v)\n", i, got, test.expect)
		}
	}
}

func TestConfigGetValueLoadKVs(t *testing.T) {
	var c = New()
	var kvs = []KV{
		{"foo.name", "foo1"},
		{"foo.id", 123},
		{"foo.temp", -50},
		{"foo.pi", 3.1415926},
		{"foo.onoff", true},

		// Overlap keys
		{"foo.redis", "redis.cluster"},
		{"foo.redis.host", "redis.cluster"}, // Should be ignored

		// Overlap keys
		{"foo.term.program", "tmux"},
		{"foo.term", "tmux"}, // Should be ignored
	}
	_, err := c.LoadKVs(kvs)
	if err != nil {
		t.Fatal(err)
	}

	for i, test := range []struct {
		key    string
		expect any
	}{
		{"foo.name", Value{string("foo1")}},
		{"foo.id", Value{int(123)}},
		{"foo.temp", Value{int(-50)}},
		{"foo.pi", Value{float64(3.1415926)}},
		{"foo.onoff", Value{bool(true)}},
		{"foo.redis", Value{string("redis.cluster")}},
		{"foo.redis.host", Value{v: nil}},
		{"foo.term.program", Value{string("tmux")}},
		//{"foo.term", Value{v: nil}},
	} {
		got := c.Get(test.key)
		if got != test.expect {
			t.Errorf("[%v] key=%v, got(%v)!=expect(%v)\n", i, test.key, got, test.expect)
		}
	}
}

func TestConfigGetScalarLoadYaml(t *testing.T) {
	var testdata = []byte(`
foo:
  name: foo1
  id: 123
  pi: 3.1415926
  temp: -50
  onoff: True
`)
	var c = New()
	err := c.LoadYaml(testdata)
	if err != nil {
		t.Fatal(err)
	}

	for i, test := range []struct {
		key          string
		defaultValue any
		knd          reflect.Kind
		expect       any
	}{
		{"foo.name", "", reflect.String, "foo1"},
		{"foo.id", "", reflect.String, "123"},
		{"foo.pi", "", reflect.String, "3.1415926"},
		{"foo.onoff", "", reflect.String, "true"},
		{"foo.name.not-exist", "default-name", reflect.String, "default-name"},
		{"foo.id", int(0), reflect.Int, int(123)},
		{"foo.id", uint(0), reflect.Uint, uint(123)},
		{"foo.id", int64(0), reflect.Int64, int64(123)},
		{"foo.id", uint64(0), reflect.Uint64, uint64(123)},
		{"foo.temp", int(0), reflect.Int, int(-50)},
		{"foo.temp", int64(0), reflect.Int64, int64(-50)},
		{"foo.pi", float64(0), reflect.Float64, float64(3.1415926)},
		{"foo.onoff", bool(false), reflect.Bool, bool(true)},
	} {
		var got any
		switch test.knd {
		case reflect.String:
			got = c.String(test.key, test.defaultValue.(string))
		case reflect.Int:
			got = c.Int(test.key, test.defaultValue.(int))
		case reflect.Uint:
			got = c.Uint(test.key, test.defaultValue.(uint))
		case reflect.Int64:
			got = c.Int64(test.key, test.defaultValue.(int64))
		case reflect.Uint64:
			got = c.Uint64(test.key, test.defaultValue.(uint64))
		case reflect.Float64:
			got = c.Float64(test.key, test.defaultValue.(float64))
		case reflect.Bool:
			got = c.Bool(test.key, test.defaultValue.(bool))
		}

		if got != test.expect {
			t.Errorf("[%v] got(\"%v\")!=expect(\"%v\")\n", i, got, test.expect)
		}
	}
}

func TestConfigGetScalarLoadJson(t *testing.T) {
	var testdata = []byte(`
{
  "foo": {
	"name": "foo1",
    "id": 123,
    "pi": 3.1415926,
    "temp": -50,
    "onoff": true
  }
}
`)
	var c = New()
	err := c.LoadJson(testdata)
	if err != nil {
		t.Fatal(err)
	}

	for i, test := range []struct {
		key          string
		defaultValue any
		knd          reflect.Kind
		expect       any
	}{
		{"foo.name", "", reflect.String, "foo1"},
		{"foo.name.not-exist", "default-name", reflect.String, "default-name"},
		{"foo.id", int(0), reflect.Int, int(123)},
		{"foo.id", uint(0), reflect.Uint, uint(123)},
		{"foo.id", int64(0), reflect.Int64, int64(123)},
		{"foo.id", uint64(0), reflect.Uint64, uint64(123)},
		{"foo.temp", int(0), reflect.Int, int(-50)},
		{"foo.temp", int64(0), reflect.Int64, int64(-50)},
		{"foo.pi", float64(0), reflect.Float64, float64(3.1415926)},
		{"foo.onoff", bool(false), reflect.Bool, bool(true)},
	} {
		var got any
		switch test.knd {
		case reflect.String:
			got = c.String(test.key, test.defaultValue.(string))
		case reflect.Int:
			got = c.Int(test.key, test.defaultValue.(int))
		case reflect.Uint:
			got = c.Uint(test.key, test.defaultValue.(uint))
		case reflect.Int64:
			got = c.Int64(test.key, test.defaultValue.(int64))
		case reflect.Uint64:
			got = c.Uint64(test.key, test.defaultValue.(uint64))
		case reflect.Float64:
			got = c.Float64(test.key, test.defaultValue.(float64))
		case reflect.Bool:
			got = c.Bool(test.key, test.defaultValue.(bool))
		}

		if got != test.expect {
			t.Errorf("[%v] got(\"%v\")!=expect(\"%v\")\n", i, got, test.expect)
		}
	}
}

func TestConfigGetScalarLoadKVs(t *testing.T) {
	var c = New()
	var kvs = []KV{
		{"foo.name", "foo1"},
		{"foo.id", 123},
		{"foo.temp", -50},
		{"foo.pi", 3.1415926},
		{"foo.onoff", true},

		// Overlay keys
		{"foo.redis", "redis.cluster"},
		{"foo.redis.host", "redis.cluster"}, // Should be ignored
	}
	_, err := c.LoadKVs(kvs)
	if err != nil {
		t.Fatal(err)
	}

	for i, test := range []struct {
		key          string
		defaultValue any
		knd          reflect.Kind
		expect       any
	}{
		{"foo.name", "", reflect.String, "foo1"},
		{"foo.id", "", reflect.String, "123"},
		{"foo.pi", "", reflect.String, "3.1415926"},
		{"foo.onoff", "", reflect.String, "true"},
		{"foo.name.not-exist", "default-name", reflect.String, "default-name"},
		{"foo.id", int(0), reflect.Int, int(123)},
		{"foo.id", uint(0), reflect.Uint, uint(123)},
		{"foo.id", int64(0), reflect.Int64, int64(123)},
		{"foo.id", uint64(0), reflect.Uint64, uint64(123)},
		{"foo.temp", int(0), reflect.Int, int(-50)},
		{"foo.temp", int64(0), reflect.Int64, int64(-50)},
		{"foo.pi", float64(0), reflect.Float64, float64(3.1415926)},
		{"foo.onoff", bool(false), reflect.Bool, bool(true)},
		{"foo.redis", "", reflect.String, "redis.cluster"},
		{"foo.redis.host", "", reflect.String, ""},
	} {
		var got any
		switch test.knd {
		case reflect.String:
			got = c.String(test.key, test.defaultValue.(string))
		case reflect.Int:
			got = c.Int(test.key, test.defaultValue.(int))
		case reflect.Uint:
			got = c.Uint(test.key, test.defaultValue.(uint))
		case reflect.Int64:
			got = c.Int64(test.key, test.defaultValue.(int64))
		case reflect.Uint64:
			got = c.Uint64(test.key, test.defaultValue.(uint64))
		case reflect.Float64:
			got = c.Float64(test.key, test.defaultValue.(float64))
		case reflect.Bool:
			got = c.Bool(test.key, test.defaultValue.(bool))
		}

		if got != test.expect {
			t.Errorf("[%v] key=%v, got(\"%v\")!=expect(\"%v\")\n", i, test.key, got, test.expect)
		}
	}
}

func TestConfigGetScalarLoadEnvs(t *testing.T) {
	var c = New()
	var envs = []string{
		"SHELL=/bin/zsh",
		"SHLVL=2",
		"SSH_AUTH_SOCK=/private/tmp/com.apple.launchd.guACPbxyU5/Listeners",
		"TERM=xterm-256color",
		"TERM_PROGRAM=tmux",                                           // Overlap with "TERM"
		"TERM_PROGRAM_VERSION=3.2",                                    // Overlap with "TERM"
		"TERM_SESSION_ID=w0t0p0:53538359-4F0F-4551-9088-091F2DDE2173", // Overlap with "TERM"
		"TMPDIR=/var/folders/z_/kk3ybcfj4p5gl21bpwb5by8h0000gn/T/",
		"TMUX=/private/tmp/tmux-501/default,48207,0",
		"TMUX_PANE=%1", // Overlap with "TMUX"
		"XPC_FLAGS=0x0",
		"XPC_SERVICE_NAME=0",
		"_P9K_SSH_TTY=/dev/ttys002",
		"_P9K_TTY=/dev/ttys002",
		"__CFBundleIdentifier=com.googlecode.iterm2",
		"__CF_USER_TEXT_ENCODING=0x0:25:52",
		"P9K_SSH=0",
		"_=/usr/bin/env",
	}
	_, err := c.LoadEnvs(envs)
	if err != nil {
		t.Fatal(err)
	}

	for i, test := range []struct {
		key          string
		defaultValue any
		knd          reflect.Kind
		expect       any
	}{
		{"shell", "", reflect.String, "/bin/zsh"},
		{"shlvl", uint64(0), reflect.Uint64, uint64(2)},
		{"ssh.auth.sock", "", reflect.String, "/private/tmp/com.apple.launchd.guACPbxyU5/Listeners"},
		{"term", "", reflect.String, "xterm-256color"},
		{"term.program", "", reflect.String, ""},
		{"term.program.version", 0.0, reflect.Float64, 0.0},
		{"term.session.id", "", reflect.String, ""},
		{"tmpdir", "", reflect.String, "/var/folders/z_/kk3ybcfj4p5gl21bpwb5by8h0000gn/T/"},
		{"tmux", "", reflect.String, "/private/tmp/tmux-501/default,48207,0"},
		{"tmux.pane", "", reflect.String, ""},
		{"xpc.flags", "", reflect.String, "0x0"},
		{"xpc.service.name", uint64(0), reflect.Uint64, uint64(0)},
		{"p9k.ssh.tty", "", reflect.String, "/dev/ttys002"},
		{"p9k.tty", "", reflect.String, "/dev/ttys002"},
		{"cfbundleidentifier", "", reflect.String, "com.googlecode.iterm2"},
		{"cf.user.text.encoding", "", reflect.String, "0x0:25:52"},
		{"p9k.ssh", uint64(0), reflect.Uint64, uint64(0)},
	} {
		var got any
		switch test.knd {
		case reflect.String:
			got = c.String(test.key, test.defaultValue.(string))
		case reflect.Int:
			got = c.Int(test.key, test.defaultValue.(int))
		case reflect.Uint:
			got = c.Uint(test.key, test.defaultValue.(uint))
		case reflect.Int64:
			got = c.Int64(test.key, test.defaultValue.(int64))
		case reflect.Uint64:
			got = c.Uint64(test.key, test.defaultValue.(uint64))
		case reflect.Float64:
			got = c.Float64(test.key, test.defaultValue.(float64))
		case reflect.Bool:
			got = c.Bool(test.key, test.defaultValue.(bool))
		}

		if got != test.expect {
			t.Errorf("[%v] key=%v, got(\"%v\")!=expect(\"%v\")\n", i, test.key, got, test.expect)
		}
	}
}

func TestConfigGetScalarLoadArgs(t *testing.T) {
	var c = New()
	var args = []string{
		"-foo.name", "foo1",
		"-foo.id=123",
		"-foo.temp=-50",
		"-foo.pi=3.1415926",
		"-foo.onoff",

		// Overlap keys
		"-foo.redis=redis.cluster",
		"-foo.redis.host=redis.cluster", // Should be ignored
	}
	_, err := c.LoadArgs(args)
	if err != nil {
		t.Fatal(err)
	}

	for i, test := range []struct {
		key          string
		defaultValue any
		knd          reflect.Kind
		expect       any
	}{
		{"foo.name", "", reflect.String, "foo1"},
		{"foo.id", "", reflect.String, "123"},
		{"foo.pi", "", reflect.String, "3.1415926"},
		{"foo.onoff", "", reflect.String, "true"},
		{"foo.name.not-exist", "default-name", reflect.String, "default-name"},
		{"foo.id", int(0), reflect.Int, int(123)},
		{"foo.id", uint(0), reflect.Uint, uint(123)},
		{"foo.id", int64(0), reflect.Int64, int64(123)},
		{"foo.id", uint64(0), reflect.Uint64, uint64(123)},
		{"foo.temp", int(0), reflect.Int, int(-50)},
		{"foo.temp", int64(0), reflect.Int64, int64(-50)},
		{"foo.pi", float64(0), reflect.Float64, float64(3.1415926)},
		{"foo.onoff", bool(false), reflect.Bool, bool(true)},
		{"foo.redis", "", reflect.String, "redis.cluster"},
		{"foo.redis.host", "", reflect.String, ""},
	} {
		var got any
		switch test.knd {
		case reflect.String:
			got = c.String(test.key, test.defaultValue.(string))
		case reflect.Int:
			got = c.Int(test.key, test.defaultValue.(int))
		case reflect.Uint:
			got = c.Uint(test.key, test.defaultValue.(uint))
		case reflect.Int64:
			got = c.Int64(test.key, test.defaultValue.(int64))
		case reflect.Uint64:
			got = c.Uint64(test.key, test.defaultValue.(uint64))
		case reflect.Float64:
			got = c.Float64(test.key, test.defaultValue.(float64))
		case reflect.Bool:
			got = c.Bool(test.key, test.defaultValue.(bool))
		}

		if got != test.expect {
			t.Errorf("[%v] key=%v, got(\"%v\")!=expect(\"%v\")\n", i, test.key, got, test.expect)
		}
	}
}

func TestConfigGetScalarOverlayWithYaml(t *testing.T) {
	var testdata1 = []byte(`
foo:
  name: foo1
`)

	var testdata2 = []byte(`
foo:
  name: foo2
  id: 123
  pi: 3.1415926
`)

	var c = New()
	err := c.LoadYaml(testdata1)
	if err != nil {
		t.Fatal(err)
	}
	err = c.LoadYaml(testdata2)
	if err != nil {
		t.Fatal(err)
	}

	for i, test := range []struct {
		key          string
		defaultValue any
		knd          reflect.Kind
		expect       any
	}{
		{"foo.name", "", reflect.String, "foo1"},
		{"foo.id", int(0), reflect.Int, int(123)},
		{"foo.pi", float64(0), reflect.Float64, float64(3.1415926)},
	} {
		var got any
		switch test.knd {
		case reflect.String:
			got = c.String(test.key, test.defaultValue.(string))
		case reflect.Int:
			got = c.Int(test.key, test.defaultValue.(int))
		case reflect.Uint:
			got = c.Uint(test.key, test.defaultValue.(uint))
		case reflect.Int64:
			got = c.Int64(test.key, test.defaultValue.(int64))
		case reflect.Uint64:
			got = c.Uint64(test.key, test.defaultValue.(uint64))
		case reflect.Float64:
			got = c.Float64(test.key, test.defaultValue.(float64))
		case reflect.Bool:
			got = c.Bool(test.key, test.defaultValue.(bool))
		}

		if got != test.expect {
			t.Errorf("[%v] got(\"%v\")!=expect(\"%v\")\n", i, got, test.expect)
		}
	}
}

func TestConfigUnmarshalRoot(t *testing.T) {
	var testdata = []byte(`
foo:
  id: 123
  name: foo1
`)

	type testConfig struct {
		Foo struct {
			Id   int    `yaml:'id'`
			Name string `yaml:'name'`
		} `yaml: 'foo'`
	}

	var got testConfig
	var expect = testConfig{
		Foo: struct {
			Id   int    `yaml:'id'`
			Name string `yaml:'name'`
		}{Id: 123, Name: "foo1"},
	}

	var c = New()
	err := c.LoadYaml(testdata)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Unmarshal(Root, &got)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expect, got) {
		t.Fatalf("expect(%v)!=got(%v)\n", expect, got)
	}
}

func TestConfigUnmarshalSubTrue(t *testing.T) {
	var testdata = []byte(`
foo:
  id: 123
  name: foo1
`)

	type testConfig struct {
		Id   int    `yaml:'id'`
		Name string `yaml:'name'`
	}

	var got testConfig
	var expect = testConfig{Id: 123, Name: "foo1"}

	var c = New()
	err := c.LoadYaml(testdata)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Unmarshal("foo", &got)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expect, got) {
		t.Fatalf("expect(%v)!=got(%v)\n", expect, got)
	}
}

func TestConfigGetScalarOverlayWithAllSources(t *testing.T) {
	var args = []string{
		"-foo.name", "foo1",
		"-foo.id=123",
	}

	var envs = []string{
		"FOO_NAME=foo2",
		"FOO_RUNTIME_MEM=100m",
		"FOO_RUNTIME_CPU=100m",
		"FOO_REDIS_PORT=16380",
	}
	var testyaml = []byte(`
foo:
  redis:
    host: redis.cluster
    port: 6380
`)
	var testjson = []byte(`
{
  "foo": {
    "mysql": {
      "dsn": "root:123456@tcp(localhost:5555)/testdb"
	}
  }
}
`)

	var err error
	var c = New()
	_, err = c.LoadArgs(args)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.LoadEnvs(envs)
	if err != nil {
		t.Fatal(err)
	}
	err = c.LoadYaml(testyaml)
	if err != nil {
		t.Fatal(err)
	}
	err = c.LoadJson(testjson)
	if err != nil {
		t.Fatal(err)
	}

	for i, test := range []struct {
		key          string
		defaultValue any
		knd          reflect.Kind
		expect       any
	}{
		{"foo.name", "", reflect.String, "foo1"},
		{"foo.id", int(0), reflect.Int, int(123)},
		{"foo.runtime.mem", "", reflect.String, "100m"},
		{"foo.runtime.cpu", "", reflect.String, "100m"},
		{"foo.redis.host", "", reflect.String, "redis.cluster"},
		{"foo.redis.port", int(0), reflect.Int, int(16380)},
		{"foo.mysql.dsn", "", reflect.String, "root:123456@tcp(localhost:5555)/testdb"},
	} {
		var got any
		switch test.knd {
		case reflect.String:
			got = c.String(test.key, test.defaultValue.(string))
		case reflect.Int:
			got = c.Int(test.key, test.defaultValue.(int))
		case reflect.Uint:
			got = c.Uint(test.key, test.defaultValue.(uint))
		case reflect.Int64:
			got = c.Int64(test.key, test.defaultValue.(int64))
		case reflect.Uint64:
			got = c.Uint64(test.key, test.defaultValue.(uint64))
		case reflect.Float64:
			got = c.Float64(test.key, test.defaultValue.(float64))
		case reflect.Bool:
			got = c.Bool(test.key, test.defaultValue.(bool))
		}

		if got != test.expect {
			t.Errorf("[%v] got(\"%v\")!=expect(\"%v\")\n", i, got, test.expect)
		}
	}
}
