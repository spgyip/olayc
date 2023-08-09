package olayc

import (
	"reflect"
	"testing"
)

// TODO: Add tests for success Get.
func TestConfigGetValueNotExisit(t *testing.T) {
	var testdata = []byte(`
foo:
  name: foo1
  id: 123
`)

	var c = New()
	err := c.LoadYamlBytes(testdata)
	if err != nil {
		t.Fatal(err)
	}

	for i, test := range []struct {
		key string
	}{
		{"fooa"},
		{"foo.namea"},
		{"foo.name.a"},
	} {
		got := c.Get(test.key)
		if got != nil {
			t.Errorf("[%v] got not nil(%v), expect nil\n", i, got)
		}
	}
}

func TestConfigYamlSingle(t *testing.T) {
	var testdata = []byte(`
foo:
  name: test-foo
  id: 123
  pi: 3.1415926
  temp: -50
  onoff: True
`)
	var c = New()
	err := c.LoadYamlBytes(testdata)
	if err != nil {
		t.Fatal(err)
	}

	for i, test := range []struct {
		key          string
		defaultValue any
		knd          reflect.Kind
		expect       any
	}{
		{"foo.name", "", reflect.String, "test-foo"},
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

func TestConfigYamlOverlay(t *testing.T) {
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
	err := c.LoadYamlBytes(testdata1)
	if err != nil {
		t.Fatal(err)
	}
	err = c.LoadYamlBytes(testdata2)
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
	err := c.LoadYamlBytes(testdata)
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
	err := c.LoadYamlBytes(testdata)
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
