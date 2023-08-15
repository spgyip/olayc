package olayc

import (
	"fmt"
	"reflect"
)

type internalFlag struct {
	full  string
	short string
	knd   reflect.Kind
	help  string
}

const (
	internalFlagPrefix = "oc."
)

func (fl internalFlag) is(key string) bool {
	if key == fl.full || key == fl.short {
		return true
	}
	return false
}

// internalFlags defines OlayConfig preserved internal flags, all prefixed with internalFlagPrefix.
var internalFlags = map[string]internalFlag{
	"help": internalFlag{
		"oc.help",
		"oc.h",
		reflect.Bool,
		"Print this help message.",
	},
	"silent": internalFlag{
		"oc.silent",
		"oc.s",
		reflect.Bool,
		"Set silent mode, default is false.",
	},
	"file.yaml": internalFlag{
		"oc.file.yaml",
		"oc.f.y",
		reflect.String,
		"Load yaml file.",
	},
	"file.json": internalFlag{
		"oc.file.json",
		"oc.f.j",
		reflect.String,
		"Load json file.",
	},
	"env": internalFlag{
		"oc.env",
		"oc.e",
		reflect.Bool,
		"Load from environments.",
	},
	"dryrun": internalFlag{
		"oc.dryrun",
		"oc.dr",
		reflect.Bool,
		"Dry run, load and print Yaml then exit.",
	},
}

// Print OlayConfig usage message.
func usageOlayc() {
	fmt.Println("Usage of olayc:")
	for _, fn := range internalFlags {
		fmt.Printf("  -%v|-%v %v\n", fn.full, fn.short, fn.knd)
		fmt.Printf("         %v\n", fn.help)
	}
}
