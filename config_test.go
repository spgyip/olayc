package olayc

/*
func TestDefaultConfig(t *testing.T) {
	var expect = map[string]any{
		"foo": map[string]any{
			"id":   123,
			"name": "foo1",
			"url":  "http://www.example.com",
		},
	}

	got, err := loadDefaultFromArgs([]string{
		"-olayc.file.yaml", "testdata/test1.yaml",
		"-olayc.file.yaml", "testdata/test2.yaml",
	})
	if err != nil {
		t.Fatal(err)
	}
	equal := reflect.DeepEqual(expect, got)
	t.Log(reflect.TypeOf(expect["foo"]))
	t.Log(reflect.TypeOf(got["foo"]))

	if !equal {
		t.Fatalf("got(%v)!=expect(%v)", got, expect)
	}
}

func TestYamlConfigGetScalar(t *testing.T) {
	c := newYamlConfig()
	err := c.loadFromFile("testdata/test1.yaml")
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		key    string
		expect any
	}{
		{"foo.id", 123},
		{"foo.name", "foo1"},
		{"foo.labels.app", "foo"},
		{"foo.labels.zone", "sz"},
		{"foo.notexist", nil},
	} {
		got := c.Get(tc.key)
		if got != tc.expect {
			t.Errorf("got(%v)!=expect(%v) with key %v", got, tc.expect, tc.key)
		}
	}
}*/
