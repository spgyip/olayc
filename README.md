Overlay configuration
===================================

```
olayc.LoadDefault()
olayc.Get("foo.name").String("")
olayc.Get("foo.id").Int(0)
olayc.Get("foo.id").Int64(0)
olayc.Get("foo.id").Uint(0)
olayc.Get("foo.id").Uint64(0)
olayc.Get("foo.id").Bool(false)
olayc.Get("temperate.degree").Float64(0.0)

var t foo struct {
    Name string
    Id int
}

olayc.Unmarshal("foo", &t)
olayc.Unmarshal(olayc.ROOT, &t)
```
