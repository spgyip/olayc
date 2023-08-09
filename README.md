Overlay configuration
===================================

Overlay configuration is composition of multiple configure sources, each source is overlayed from bottom to top.
The top layer is visible if there is key conflicted among layers.

![layers](readme/images/layers.png)

# Examples

See `examples/`. Build with `make all`, binaries are built in `bin/`.

## Load yaml files

```
./bin/simple -oc.f.y=./testdata/test1.yaml \
             -oc.f.y=./testdata/test2.yaml
```

> The left yaml file is more prior to right file, thus, if there are same keys in `test1.yaml` and `test2.yaml`, the value in `test1.yaml` will be got.

## Load from commandline arguments

```
./bin/simple -oc.f.y=./testdata/test1.yaml \
             -oc.f.y=./testdata/test2.yaml \
             -foo.id=999 \
             -foo.name=hello
```

> Commandline arguments is more prior to yaml files, thus, `foo.id` will be got with value `999` which is from commandline argument.

## Silent mode

There are default verbose logs, silent mode can be turned on with `-oc.s`:

```
./bin/simple -oc.s \
             -oc.f.y=./testdata/test1.yaml \
             -oc.f.y=./testdata/test2.yaml
             -foo.id=999 \
             -foo.name=hello
```

# Usage

## Single value

```go
olayc.Load()
id := olayc.Int("foo.id", 99))
name: = olayc.String("foo.name", "foo")
url := olayc.String("foo.url", "http://www.default.com"))
```

## Unmarshal

Unmarshal is using yaml field tags.

```go

var cfg struct {
	Foo struct {
		Id   int    `yaml:'id'`
		Name string `yaml: 'name'`
		Url  string `yaml: 'url'`
	} `yaml: 'foo'`
}

olayc.Load()
olayc.Unmarshal(olayc.Root, &cfg)
```

# Priority

The default olayc has default priority when multiple configure sources are loaded:

- Commandline arguments, left prior
- Environment variables
- Yaml/Json Files, left prior

> In fact, the source has more priority is placed on upper layer.

# Tests

Make sure run tests with `make test` if you fork this repo and try to add some features.
