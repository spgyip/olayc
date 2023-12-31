package main

import (
	"fmt"
	"reflect"

	"github.com/spgyip/olayc"
)

func main() {
	olayc.Load(
		olayc.WithFileRequire("test1.yaml"),
		olayc.WithFileRequire("test2.yaml"),
		olayc.WithUsage("foo.id", reflect.Int, 99, "Set foo ID"),
	)
	fmt.Println("foo.id:", olayc.Int("foo.id", 99))
	fmt.Println("foo.name:", olayc.String("foo.name", "foo"))
	fmt.Println("foo.url:", olayc.String("foo.url", "http://www.default.com"))
	fmt.Println("foo.redis.host:", olayc.String("foo.redis.host", "localhost"))
	fmt.Println("foo.redis.port:", olayc.Int("foo.redis.port", 0))
}
