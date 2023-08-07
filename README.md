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
id := olayc.Int("foo.id", 99))
name: = olayc.String("foo.name", "foo")
url := olayc.String("foo.url", "http://www.default.com"))
```

Load yaml files with commandline arguments

```
./bin/simple -oc.f.y=./testdata/test1.yaml \
             -oc.f.y=./testdata/test2.yaml
```

Turn on silent mode with commandline argument

```
./bin/simple -oc.s \
             -oc.f.y=./testdata/test1.yaml \
             -oc.f.y=./testdata/test2.yaml
```





