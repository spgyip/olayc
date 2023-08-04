Overlay configuration
===================================

# Examples

See `examples/`. Build with `make all`, binarieds are built in `bin/`.

# Usage

Load default olayc and get value with key

```
// simple/main.go
olayc.Load()
olayc.Get("foo.id")
```

Load yaml files with commandline arguments

```
./bin/simple -olayc.file.yaml=./testdata/test1.yaml \
         -olayc.file.yaml=./testdata/test2.yaml
```

Turn on silent mode with commandline argument

```
./bin/simple -olayc.silent=true
```





