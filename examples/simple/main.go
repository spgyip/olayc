package main

import (
	"fmt"

	"github.com/spgyip/olayc"
)

func main() {
	olayc.Load(
		olayc.WithFileRequire("test1.yaml"),
		olayc.WithFileRequire("test2.yaml"),
	)
	fmt.Println(olayc.Int("foo.id", 99))
	fmt.Println(olayc.String("foo.name", "foo"))
	fmt.Println(olayc.String("foo.url", "http://www.default.com"))
}
