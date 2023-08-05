package main

import (
	"fmt"

	"github.com/spgyip/olayc"
)

func main() {
	// TODO: olayc.Require()
	olayc.Load()
	fmt.Println(olayc.Get("foo.id"))
	fmt.Println(olayc.Get("foo.name"))
	fmt.Println(olayc.Get("foo.labels"))
}
