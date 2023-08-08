package main

import (
	"fmt"

	"github.com/spgyip/olayc"
)

func main() {
	olayc.Load()
	fmt.Println(olayc.Int("foo.id", 99))
	fmt.Println(olayc.String("foo.name", "foo"))
	fmt.Println(olayc.String("foo.url", "http://www.default.com"))
}
