# orderedmapjson

orderedmapjson is a package that extends the [orderedmap](https://github.com/elliotchance/orderedmap) package by adding support for JSON marshalling and unmarshalling. Keys are required to be strings, as per the JSON specification.

# Installation

```
go get github.com/GitRowin/orderedmapjson
```

# Usage

The library provides three types: `TypedOrderedMap`, `AnyOrderedMap` and `AnyOrderedMapSlice`.

## TypedOrderedMap

TypedOrderedMap allows you to specify the value type of the map.

```go
m := orderedmapjson.NewTypedOrderedMap[int]()

m.Set("q", 1)
m.Set("w", 2)
m.Set("e", 3)
m.Set("r", 4)
m.Set("t", 5)
m.Set("y", 6)

b, err := json.Marshal(m)

if err != nil {
    panic(err)
}

fmt.Println(string(b)) // Output: {"q":1,"w":2,"e":3,"r":4,"t":5,"y":6}
```

## AnyOrderedMap

AnyOrderedMap is similar to `TypedOrderedMap[any]`, with one key difference: it unmarshals JSON objects at any level as `*AnyOrderedMap` rather than `map[string]any`. This ensures that the order of nested objects is preserved too.

```go
const input = `{"foo":"bar","obj":{"2":2,"3":"3","1":1},"baz":[1, "2", 3, null]}`

var m *orderedmapjson.AnyOrderedMap
if err := json.Unmarshal([]byte(input), &m); err != nil {
    panic(err)
}

if obj, ok := m.Get("obj"); ok {
    fmt.Printf("%T %v\n", obj, obj) // Output: *orderedmapjson.AnyOrderedMap {2:2,3:3,1:1}
}

if baz, ok := m.Get("baz"); ok {
    fmt.Println(baz.([]any)[2]) // Output: 3
}
```

## AnyOrderedMapSlice

For unmarshalling JSON arrays. Like AnyOrderedMap, it unmarshals JSON objects at any level as `*AnyOrderedMap` rather than `map[string]any`. The decoded elements are available through the `Values` field. The zero value is ready to use.

```go
const input = `["foo",1,{"3":"3","1":"1","2":"2"}]`

var values orderedmapjson.AnyOrderedMapSlice
if err := json.Unmarshal([]byte(input), &values); err != nil {
    panic(err)
}

fmt.Println(values.Values[2]) // Output: {3:3,1:1,2:2}
```

# Options

## SetUseNumber

By default, numbers are unmarshalled as `float64`, consistent with `encoding/json`. This can silently lose precision for large numbers. Call `SetUseNumber(true)` to unmarshal numbers as `json.Number` instead, like `json.Decoder`'s `UseNumber`. This also preserves the exact number representation when marshalling the map back to JSON.

```go
m := orderedmapjson.NewAnyOrderedMap()
m.SetUseNumber(true)

if err := json.Unmarshal([]byte(`{"foo":10000000000000000001}`), &m); err != nil {
    panic(err)
}

foo, _ := m.Get("foo")
fmt.Printf("%T %v\n", foo, foo) // Output: json.Number 10000000000000000001
```

Note that the setting must be enabled before unmarshalling, so unmarshalling into a nil map (`var m *orderedmapjson.AnyOrderedMap`) always uses `float64`. `AnyOrderedMapSlice` does not have this limitation, because its zero value is usable directly:

```go
var values orderedmapjson.AnyOrderedMapSlice
values.SetUseNumber(true)

if err := json.Unmarshal([]byte(`[10000000000000000001]`), &values); err != nil {
    panic(err)
}

fmt.Println(values.Values[0]) // Output: 10000000000000000001
```

## SetEscapeHTML

By default, problematic HTML characters (`<`, `>`, and `&`) are escaped when marshalling, consistent with `encoding/json`. Call `SetEscapeHTML(false)` to disable this. Note that `json.Marshal` HTML-escapes the output of custom marshallers regardless of this setting, so a `json.Encoder` must be used to produce unescaped output:

```go
m := orderedmapjson.NewAnyOrderedMap()
m.SetEscapeHTML(false)
m.Set("foo", "<bar>")

var buf bytes.Buffer
encoder := json.NewEncoder(&buf)
encoder.SetEscapeHTML(false)

if err := encoder.Encode(m); err != nil {
    panic(err)
}

fmt.Print(buf.String()) // Output: {"foo":"<bar>"}
```

`AnyOrderedMap`'s and `AnyOrderedMapSlice`'s `SetEscapeHTML` also apply the setting to all nested `AnyOrderedMap` values. `TypedOrderedMap`'s only applies to the map itself: nested maps have their own setting.

# Notes

- The map types must be created using their constructors: their zero value is only usable as an unmarshalling target. `AnyOrderedMapSlice`'s zero value is fully usable.
- The settings of `AnyOrderedMap` and `AnyOrderedMapSlice` are inherited by nested maps created during unmarshalling. `TypedOrderedMap`'s settings only apply to the map itself: nested maps have their own settings.
- `Copy()` returns a shallow copy: the values are copied as-is, so nested maps and slices are shared between the original and the copy.
- Like `encoding/json`, unmarshalling into a non-empty map does not clear it first: existing entries are kept unless overwritten.
- Like `encoding/json`, duplicate keys are permitted and the last value wins. The key keeps the position of its first occurrence.
- Like `encoding/json`, a nil `Values` field is marshalled as `null` rather than `[]`.
- Unmarshalling `null` clears the map (or sets `AnyOrderedMapSlice`'s `Values` to nil), consistent with how `encoding/json` unmarshals `null` into maps and slices. Settings are kept.
