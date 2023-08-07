Overlay configuration
===================================

OlayConfig is composition of multiple configure sources, each source is overlayed from bottom to top.
The top layer is visible if there is key conflicted among layers.

![layers](readme/images/layers.png)

# Examples

See `examples/`. Build with `make all`, binaries are built in `bin/`.

# Usage

Load default olayc and get value with key

```
// simple/main.go
olayc.Load()
olayc.Get("foo.id")

olayc.String("foo.name", "default")
olayc.Int("foo.id", 0)
olayc.Uint("foo.id", 0)
olayc.Int64("foo.id", 0)
olayc.Uint64("foo.id", 0)
olayc.Float64("foo.id", 0)
olayc.Bool("foo.onof", false)

var redis = struct {
    Host string `yaml:"host"`
    Port int `yaml:"port"`
}
olayc.Unmarshal("foo.redis", &redis)
```

Load yaml files with commandline arguments

```
./bin/simple -oc.f.y=./testdata/test1.yaml \
             -oc.f.y=./testdata/test2.yaml
```

Turn on silent mode with commandline argument

```
./bin/simple -oc.s=true
```





