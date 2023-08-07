Overlay configuration
===================================

![layers](readme/images/layers.png)

# Examples

See `examples/`. Build with `make all`, binaries are built in `bin/`.

# Usage

Load default olayc and get value with key

```
// simple/main.go
olayc.Load()
olayc.Get("foo.id")
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





