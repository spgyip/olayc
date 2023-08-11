package olayc

import (
	"reflect"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// Copy map values from src to dst, keep the src value if key is conflicted.
//
// There are some cases(see map_test.go), for a specified key:
// - If it doesn't exist in dst map, copy src map value to dst map.
// For example, dst's `foo.name=foo0` will be copied to result.
// src:
// `
// foo:
//   id: 123
// `
// dst:
// `
// foo:
//   name: foo0
// `
// result:
// `
// foo:
//   id: 123
//   name: foo0
// `
//
// - If it's a leaf node in dst map, the src map value will be ignored.
// For example, dst's `foo.id=456` will be ignored, the result is using src's `foo.id=123`
// src:
// `
// foo:
//   id: 123
// `
// dst:
// `
// foo:
//   id: 456
//   name: foo0
// `
// result:
// `
// foo:
//  id: 123
//  name: foo0
// `
//
// - If it's not a leaf node in dst map, but it's leaf node in src map, the src map value will be ignored.
// For example, the dst's `foo.redis="redis.cluster"` will be ignored, because the src map has a key path `foo.redis.*`.
// src:
// `
// foo:
//   id: 123
//   redis:
//     host: redis.cluster
//     port: 6380
// `
// dst:
// `
// foo:
//   name: foo0
//   redis: redis.cluster
// `
// result:
// `
// foo:
//   id: 123
//   foo: foo0
//   redis:
//     host: redis.cluster
//     port: 6380
// `
func copyMap(dst map[any]any, src map[any]any) {
	copyMapDFS(dst, src)
}

// Deep first search copy.
func copyMapDFS(dst map[any]any, src map[any]any) {
	for k, valSrc := range src {
		// Key doesn't exisit in dst, copy it to dst.
		valDst, ok := dst[k]
		if !ok {
			dst[k] = valSrc
			continue
		}

		// If dst value type is scalar type or array/slice,
		// which means it's a leaf node, keep the dst value and ignore the src value.
		typDst := reflect.TypeOf(valDst)
		switch typDst.Kind() {
		case reflect.Int, reflect.Int8,
			reflect.Int32, reflect.Int64, reflect.Uint,
			reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64,
			reflect.String, reflect.Bool,
			reflect.Array, reflect.Slice:
			// Do nothing and continue next key.
			continue
		}

		// If dst and src are both map type, recursively copy with DFS.
		nextDst, isDstMapType := dst[k].(map[any]any)
		nextSrc, isSrcMapType := src[k].(map[any]any)
		if !isDstMapType || !isSrcMapType {
			continue
		}
		copyMapDFS(nextDst, nextSrc)
	}
}

// Convert map[string]any to map[any]any.
// By using yaml marshal and unmarshal.
func convertMap(m map[string]any) (map[any]any, error) {
	data, err := yaml.Marshal(m)
	if err != nil {
		return nil, errors.Wrap(err, "Convert map fail")
	}

	var out = make(map[any]any)
	err = yaml.Unmarshal(data, &out)
	if err != nil {
		return nil, errors.Wrap(err, "Convert map fail")
	}
	return out, nil
}
