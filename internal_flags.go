package olayc

type internalFlag struct {
	full  string
	short string
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

var internalFlags = map[string]internalFlag{
	"help": internalFlag{
		"oc.help",
		"oc.h",
		"Print this help message.",
	},
	"silent": internalFlag{
		"oc.silent",
		"oc.s",
		"Set silent mode, default is false.",
	},
	"file.yaml": internalFlag{
		"oc.file.yaml",
		"oc.f.y",
		"Load yaml file.",
	},
	"file.json": internalFlag{
		"oc.file.json",
		"oc.f.j",
		"Load json file.",
	},
}
