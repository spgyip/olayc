package olayc

import "reflect"

// Copy map key-values from src to dst, dst can't override the src values.
// There are some cases, for a specified key:
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
//     port: 3096
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
//     port: 3096
// `
func CopyMap(dst map[any]any, src map[any]any) {
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
		case reflect.Int, reflect.Int32,
			reflect.Int64, reflect.Uint,
			reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64,
			reflect.String, reflect.Bool,
			reflect.Array, reflect.Slice:
			// Do nothing and continue next key.
			continue
		}

		// If dst and src are both map type, recursively copy.
		nextDst, isDstMapType := dst[k].(map[any]any)
		nextSrc, isSrcMapType := src[k].(map[any]any)
		if !isDstMapType || !isSrcMapType {
			continue
		}
		copyMapDFS(nextDst, nextSrc)
	}
}
