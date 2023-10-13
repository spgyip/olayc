package olayc

import (
	"reflect"
	"testing"
)

func TestMergeMapAppendEmpty(t *testing.T) {
	var dst = map[any]any{}
	var src = map[any]any{
		"foo": map[any]any{
			"id": 123,
		},
	}

	var expect = map[any]any{
		"foo": map[any]any{
			"id": 123,
		},
	}

	copyMap(dst, src)
	if !reflect.DeepEqual(dst, expect) {
		t.Log("m0:", dst)
		t.Log("expect:", expect)
		t.Fatal("Merged result is not expected.")
	}
}

func TestMergeMapAppendNew(t *testing.T) {
	var dst = map[any]any{
		"foo": map[any]any{
			"id": 123,
		},
	}
	var src = map[any]any{
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

	copyMap(dst, src)
	if !reflect.DeepEqual(dst, expect) {
		t.Log("m0:", dst)
		t.Log("expect:", expect)
		t.Fatal("Merged result is not expected.")
	}
}

func TestMergeMapIgnoreNode(t *testing.T) {
	var dst = map[any]any{
		"foo": map[any]any{
			"id": 123,
		},
	}
	var src = map[any]any{
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

	copyMap(dst, src)
	if !reflect.DeepEqual(dst, expect) {
		t.Log("m0:", dst)
		t.Log("expect:", expect)
		t.Fatal("Merged result is not expected.")
	}
}

func TestMergeMapIgnoreSubTree(t *testing.T) {
	var dst = map[any]any{
		"foo": map[any]any{
			"id": 123,
			"redis": map[any]any{
				"host": "redis.cluster",
				"ip":   6380,
			},
		},
	}
	var src = map[any]any{
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

	copyMap(dst, src)
	if !reflect.DeepEqual(dst, expect) {
		t.Log("m0:", dst)
		t.Log("expect:", expect)
		t.Fatal("Merged result is not expected.")
	}
}
