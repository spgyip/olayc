package olayc

import (
	"reflect"
	"testing"
)

func TestMergeMapAppendEmpty(t *testing.T) {
	var m0 = map[any]any{}
	var m1 = map[any]any{
		"foo": map[any]any{
			"id": 123,
		},
	}

	var expect = map[any]any{
		"foo": map[any]any{
			"id": 123,
		},
	}

	copyMap(m0, m1)
	if !reflect.DeepEqual(m0, expect) {
		t.Log("m0:", m0)
		t.Log("expect:", expect)
		t.Fatal("Merged result is not expected.")
	}
}

func TestMergeMapAppendNew(t *testing.T) {
	var m0 = map[any]any{
		"foo": map[any]any{
			"id": 123,
		},
	}
	var m1 = map[any]any{
		"foo": map[any]any{
			"name": "foo1",
		},
	}

	var expect = map[any]any{
		"foo": map[any]any{
			"id":   123,
			"name": "foo1",
		},
	}

	copyMap(m0, m1)
	if !reflect.DeepEqual(m0, expect) {
		t.Log("m0:", m0)
		t.Log("expect:", expect)
		t.Fatal("Merged result is not expected.")
	}
}

func TestMergeMapIgnoreNode(t *testing.T) {
	var m0 = map[any]any{
		"foo": map[any]any{
			"id": 123,
		},
	}
	var m1 = map[any]any{
		"foo": map[any]any{
			"id":   456,
			"name": "foo1",
		},
	}

	var expect = map[any]any{
		"foo": map[any]any{
			"id":   123,
			"name": "foo1",
		},
	}

	copyMap(m0, m1)
	if !reflect.DeepEqual(m0, expect) {
		t.Log("m0:", m0)
		t.Log("expect:", expect)
		t.Fatal("Merged result is not expected.")
	}
}

func TestMergeMapIgnoreSubTree(t *testing.T) {
	var m0 = map[any]any{
		"foo": map[any]any{
			"id": 123,
			"redis": map[any]any{
				"host": "redis.cluster",
				"ip":   6380,
			},
		},
	}
	var m1 = map[any]any{
		"foo": map[any]any{
			"name":  "foo1",
			"redis": "redis.cluster",
		},
	}

	var expect = map[any]any{
		"foo": map[any]any{
			"id":   123,
			"name": "foo1",
			"redis": map[any]any{
				"host": "redis.cluster",
				"ip":   6380,
			},
		},
	}

	copyMap(m0, m1)
	if !reflect.DeepEqual(m0, expect) {
		t.Log("m0:", m0)
		t.Log("expect:", expect)
		t.Fatal("Merged result is not expected.")
	}
}
