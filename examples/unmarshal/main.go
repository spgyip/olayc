package main

import (
	"fmt"
	"os"

	"github.com/spgyip/olayc"
)

type config struct {
	Id   int    `yaml:'id'`
	Name string `yaml: 'name'`
	Url  string `yaml: 'url'`
}

func main() {
	var cfg config
	olayc.Load()
	err := olayc.Unmarshal("foo", &cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	fmt.Println(cfg)
}
