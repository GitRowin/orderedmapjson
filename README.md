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

AnyOrderedMap is similar to `TypedOrderedMap[any]`, with one key difference: it unmarshals JSON objects at any level as `AnyOrderedMap[any]` rather than `map[string]any`. This ensures that the order of nested objects is preserved too.

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

For unmarshalling JSON arrays. Like AnyOrderedMap, it unmarshals JSON objects at any level as `AnyOrderedMap[any]` rather than `map[string]any`.

```go
const input = `["foo",1,{"3":"3","1":"1","2":"2"}]`

var values AnyOrderedMapSlice
if err := json.Unmarshal([]byte(input), &values); err != nil {
    panic(err)
}

fmt.Println(values[2]) // Output: {3:3,1:1,2:2}
```
