package main

import (
	"fmt"
	"os"

	"github.com/spgyip/olayc"
)

type config struct {
	Foo struct {
		Id   int    `yaml:'id'`
		Name string `yaml: 'name'`
		Url  string `yaml: 'url'`
	} `yaml: 'foo'`
}

func main() {
	var cfg config

	olayc.Load(
		olayc.WithFileRequire("test1.yaml"),
		olayc.WithFileRequire("test2.yaml"),
	)
	err := olayc.Unmarshal(olayc.Root, &cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	fmt.Println(cfg)
}
