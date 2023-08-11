Overlay configuration
===================================

Configurations are very common scenerios controlling program behaviors at runtime. There are different ways of configuration, e.g. from configure files, commandline arguments, environment variables. Overlay configuration is composition of multiple configure sources, and provides a unified interface to load from different configure sources. 

```
                                [### Configure sources ###]

                                -------------------------
                       --------- | commandline arguments |
                       |         -------------------------
    ---------------    |         --------------
    |    olayc    | <----------- | Yaml files |
    ---------------    |         --------------
                       |         --------------
                       --------- |    ENVs    |
                                 --------------
```

Currently supported: 
- [X] Commandline argument
- [X] Yaml file
- [X] Json file
- [ ] Environment
- [ ] Etcd KV

# Overlay

Every configure sources are overlayed from bottom to top, the upper layer has more priority than the lower ones, which means if there are keys conflited the value in the upper layer will be got.

![layers](readme/images/layers.png)

# Examples

See `examples/`. Build with `make all`, binaries are built in `bin/`.

## Load yaml files

Use `-of.f.y=...` to add yaml file.

```
./bin/simple -oc.f.y=./testdata/test1.yaml \
             -oc.f.y=./testdata/test2.yaml
```

> The left yaml file is prior to right file, thus, if there are same keys in `test1.yaml` and `test2.yaml`, the value in `test1.yaml` will be got.

## Load json files

Use `-of.f.j=...` to add json file.

```
./bin/simple -oc.f.y=./testdata/test1.yaml \
             -oc.f.y=./testdata/test2.yaml
             -oc.f.j=./testdata/test1.yaml
```

## Load from commandline arguments

Use commandline argument seperated by `.`.

```
./bin/simple -oc.f.y=./testdata/test1.yaml \
             -oc.f.y=./testdata/test2.yaml \
             -oc.f.j=./testdata/test1.yaml
             -foo.redis.host=redis.othercluster \
             -foo.redis.port=999
```

> Commandline arguments is prior to yaml files, thus, `foo.id` will be got with value `999` which is from commandline argument.

## Silent mode

There are default verbose logs, silent mode can be turned on with `-oc.s`:

```go
./bin/simple -oc.s \
             -oc.f.y=./testdata/test1.yaml \
             -oc.f.y=./testdata/test2.yaml
             -foo.id=999 \
             -foo.name=hello
```

## Dry run mode

Use `-oc.dr` to turn on dry run mode, olayc loads and prints out the merged configurations with yaml format then exits the program. It's convenient for pre-checking.

```go
./bin/simple -oc.dr \
             -oc.s \
             -oc.f.y=./testdata/test1.yaml \
             -oc.f.y=./testdata/test2.yaml
             -foo.id=999 \
             -foo.name=hello

[OlayConfig] Dry run mode is on, program will exit after yaml printed.
foo:
  id: 123
  labels:
    app: foo
    zone: sz
  name: foo1
  redis:
    host: redis.cluster
    port: 8306
  url: http://www.example.com

```

# Usage

## Load

It's easy to initialize the default olayc with `Load()`.

```go
olayc.Load()
```

It's very common there are configure files must be loaded before program startup. 
Use `WithFileRequire()` on `Load()`, program will terminates if the required file is not loaded.

```
olayc.Load(
    olayc.WithFileRequire("test1.yaml"),
)
```

## Get scalar value

```go
olayc.Load()
id := olayc.Int("foo.id", 99))
name: = olayc.String("foo.name", "foo")
url := olayc.String("foo.url", "http://www.default.com"))
```

## Unmarshal to struct

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

The default olayc has default priority when multiple configure sources are loaded, which are:

- Commandline arguments, left prior
- Environment variables
- Yaml/Json Files, left prior

# Internal olayc flags

There are internal olayc flags which are prefix with '-oc.|--oc.', use `-oc.h` to see help message.

```
Usage of olayc:
  oc.help | oc.h
         Print this help message.
  oc.silent | oc.s
         Set silent mode, default is false.
  oc.file.yaml | oc.f.y
         Load yaml file.
  oc.file.json | oc.f.j
         Load json file.
  oc.dryrun | oc.dr
         Dry run, load and print Yaml then exit.
```

# Tests

Make sure run tests with `make test` if you fork this repo and try to add some features.
