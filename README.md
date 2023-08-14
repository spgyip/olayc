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
- [X] Commandline arguments
- [X] Yaml file
- [X] Json file
- [X] Environments
- [ ] Etcd KV

# Overlay

Every configure sources are overlayed from bottom to top, the upper layer has more priority than the lower ones, which means if there are keys conflited the value in the upper layer will be got.

![layers](readme/images/layers.png)

# Examples

See `examples/`. Build with `make all`, binaries are built in `bin/`.

## Load yaml files

Use `-of.f.y=...` to add yaml file.

```shell
./bin/simple -oc.f.y=./testdata/test1.yaml \
             -oc.f.y=./testdata/test2.yaml
```

> The left yaml file is prior to right file, thus, if there are same keys in `test1.yaml` and `test2.yaml`, the value in `test1.yaml` will be got.

## Load json files

Use `-of.f.j=...` to add json file.

```shell
./bin/simple -oc.f.y=./testdata/test1.yaml \
             -oc.f.y=./testdata/test2.yaml
             -oc.f.j=./testdata/test1.json
```

## Load from commandline arguments

Use commandline argument seperated by `.`.

```shell
./bin/simple -oc.f.y=./testdata/test1.yaml \
             -oc.f.y=./testdata/test2.yaml \
             -oc.f.j=./testdata/test1.json
             -foo.redis.host=redis.othercluster \
             -foo.redis.port=999
```

## Load from environment variables

Environment variables are not loaded on default, can be turned on with `-oc.env | -oc.e`. The environment format "FOO\_NAME" is converted to "foo.name".

```shell
FOO_NAME=foo-env ./bin/simple -oc.e \
                              -oc.f.y=./testdata/test1.yaml \
                              -oc.f.y=./testdata/test2.yaml \
                              -oc.f.j=./testdata/test1.json
                              -foo.redis.host=redis.othercluster \
                              -foo.redis.port=999
```

## Silent mode

There are default verbose logs, silent mode can be turned on with `-oc.s`:

```shell
./bin/simple -oc.s \
             -oc.e \
             -oc.f.y=./testdata/test1.yaml \
             -oc.f.y=./testdata/test2.yaml \
             -oc.f.j=./testdata/test1.json \
             -foo.redis.host=redis.othercluster \
             -foo.redis.port=999
```

## Dry run mode

Use `-oc.dr` to turn on dry run mode, olayc loads and prints out the merged configurations with yaml format then exits the program. It's convenient for pre-checking what configurations are loaded.

```shell
./bin/simple -oc.dr \
             -oc.s \
             -oc.e \
             -oc.f.y=./testdata/test1.yaml \
             -oc.f.y=./testdata/test2.yaml \
             -oc.f.j=./testdata/test1.json \
             -foo.redis.host=redis.othercluster \
             -foo.redis.port=999

[OlayConfig] Dry run mode is on, program will exit after yaml printed.
foo:
  id: 123
  labels:
    app: foo
    zone: sz
  name: foo1
  redis:
    host: redis.othercluster
    port: 999
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

```go
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

- Commandline arguments
- Environment variables
- Yaml/Json Files

# Internal olayc flags

There are internal olayc flags which are prefix with '-oc.|--oc.', use `-oc.h` to see help message.

```shell
Usage of olayc:
  -oc.help | -oc.h
         Print this help message.
  -oc.silent | -oc.s
         Set silent mode, default is false.
  -oc.file.yaml | -oc.f.y
         Load yaml file.
  -oc.file.json | -oc.f.j
         Load json file.
  -oc.env | -oc.e
         Load from environments.
  -oc.dryrun | -oc.dr
         Dry run, load and print Yaml then exit.
```

# Overlap keys

When use commandline arguments or environment variables, keys may be overlap, for examples

```shell
-foo.redis=cluster1
-foo.redis.name=cluster2
```

```shell
FOO_REDIS=cluster1
FOO_REDIS_NAME=cluster2
```

There is only one value can be got when get with key 'foo.redis'. This depends on the order of loading. The previously loaded key is prior to the latter ones, the latter ones will be ignored.

If the "-foo.redis=cluster1" is loaded previously, the "-foo.redis.name=cluster2" will be ignored. This is resulting to:

```go
Get("foo.redis")         => "cluster1"
Get("foo.redis.name")    => nil(NOT EXIST)
```

If the "-foo.redis.name=cluster2" is loaded previously, the "-foo.redis=cluster1" will be ignored. This is resulting to:

```go
Get("foo.redis")         => {"name": "cluster2"}
Get("foo.redis.name")    => "cluster2"
```

> Should notice that `Get("foo.redis")` will return a subnode, but not nil, in the case 2.

# Tests

Make sure run tests with `make test` if you fork this repo and try to add some features.
