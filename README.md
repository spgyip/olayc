Overlay configuration
===================================

Overlay configuration is composition of multiple configure sources, each source is overlayed from bottom to top.
The top layer is visible if there is key conflicted among layers.

![layers](readme/images/layers.png)

# Examples

See `examples/`. Build with `make all`, binaries are built in `bin/`.

```
./bin/simple -oc.f.y=./testdata/test1.yaml \
             -oc.f.y=./testdata/test2.yaml
```

There are default verbose logs, silent mode can be turned on with `-oc.s`:

```
./bin/simple -oc.s \
             -oc.f.y=./testdata/test1.yaml \
             -oc.f.y=./testdata/test2.yaml
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
	Id   int    `yaml:'id'`
	Name string `yaml: 'name'`
	Url  string `yaml: 'url'`
}

olayc.Load()
olayc.Unmarshal("foo", &cfg)
```

