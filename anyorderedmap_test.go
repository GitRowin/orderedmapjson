package orderedmapjson

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestAnyOrderedMap(t *testing.T) {
	const input = `{"foo":"bar","123":true,"abc":null,"_obj":{"foo":"bar"},"_array":[{"b":"b","a":"a","0":null,"123":{"q":"","w":"","e":"","r":"","t":"","y":""}}],"q":"","w":"","e":"","r":"","t":"","y":""}`

	var m *AnyOrderedMap
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := json.Marshal(m)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := string(data)

	if output != input {
		t.Fatalf("expected %s, got %s", input, output)
	}
}

func TestAnyOrderedMapSlice(t *testing.T) {
	const input = `[{"c":"c","a":"a","b":"b"},"test",123,{"3":"3","2":"2","1":"1"}]`

	var values AnyOrderedMapSlice
	if err := json.Unmarshal([]byte(input), &values); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := json.Marshal(values)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := string(data)

	if output != input {
		t.Fatalf("expected %s, got %s", input, output)
	}
}

func TestAnyOrderedMapEmpty(t *testing.T) {
	const input = `{"obj":{},"array":[],"nested":[[],[{}]]}`

	var m *AnyOrderedMap
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := json.Marshal(m)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := string(data)

	if output != input {
		t.Fatalf("expected %s, got %s", input, output)
	}
}

func TestAnyOrderedMapSliceEmpty(t *testing.T) {
	const input = `[]`

	var values AnyOrderedMapSlice
	if err := json.Unmarshal([]byte(input), &values); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := json.Marshal(values)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := string(data)

	if output != input {
		t.Fatalf("expected %s, got %s", input, output)
	}
}

func TestAnyOrderedMapSliceUseNumber(t *testing.T) {
	const input = `[10000000000000000001,{"foo":1e2},[0.1000]]`

	var values AnyOrderedMapSlice
	values.SetUseNumber(true)

	if err := json.Unmarshal([]byte(input), &values); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := values.Values[0].(json.Number); !ok {
		t.Fatalf("expected json.Number, got %T", values.Values[0])
	}

	data, err := json.Marshal(values)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := string(data)

	if output != input {
		t.Fatalf("expected %s, got %s", input, output)
	}
}

func TestAnyOrderedMapSliceSetEscapeHTML(t *testing.T) {
	const input = `["<foo>",{"bar":"&"},[">"]]`

	var values AnyOrderedMapSlice
	if err := json.Unmarshal([]byte(input), &values); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Applies recursively to nested maps
	values.SetEscapeHTML(false)

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(values); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := strings.TrimSuffix(buf.String(), "\n")

	if output != input {
		t.Fatalf("expected %s, got %s", input, output)
	}
}

func TestAnyOrderedMapUseNumber(t *testing.T) {
	const input = `{"foo":10000000000000000001,"bar":1e2,"array":[0.1000]}`

	m := NewAnyOrderedMap()
	m.SetUseNumber(true)

	if err := json.Unmarshal([]byte(input), &m); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	foo, _ := m.Get("foo")

	if _, ok := foo.(json.Number); !ok {
		t.Fatalf("expected json.Number, got %T", foo)
	}

	data, err := json.Marshal(m)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := string(data)

	if output != input {
		t.Fatalf("expected %s, got %s", input, output)
	}
}

func TestAnyOrderedMapTrailingData(t *testing.T) {
	m := NewAnyOrderedMap()

	if err := m.UnmarshalJSON([]byte(`{"foo":1} {"bar":2}`)); err == nil {
		t.Fatal("expected error")
	}

	var values AnyOrderedMapSlice

	if err := values.UnmarshalJSON([]byte(`[1] [2]`)); err == nil {
		t.Fatal("expected error")
	}
}

func TestAnyOrderedMapInvalidInput(t *testing.T) {
	m := NewAnyOrderedMap()

	for _, input := range []string{``, `[1]`, `123`, `{`, `{"a":`, `{"a":1`, `{"a":1,`, `{"a":{`, `{"a":[`, `{"a":1}x`} {
		if err := m.UnmarshalJSON([]byte(input)); err == nil {
			t.Fatalf("expected error for input %q", input)
		}
	}
}

func TestAnyOrderedMapSliceInvalidInput(t *testing.T) {
	var values AnyOrderedMapSlice

	for _, input := range []string{``, `{"a":1}`, `123`, `[`, `[1`, `[1,`, `[{`, `[[`, `[1]x`} {
		if err := values.UnmarshalJSON([]byte(input)); err == nil {
			t.Fatalf("expected error for input %q", input)
		}
	}
}

func TestAnyOrderedMapMarshalUnsupportedType(t *testing.T) {
	m := NewAnyOrderedMap()
	m.Set("ch", make(chan int))

	if _, err := json.Marshal(m); err == nil {
		t.Fatal("expected error")
	}

	var values AnyOrderedMapSlice
	values.Values = []any{make(chan int)}

	if _, err := json.Marshal(values); err == nil {
		t.Fatal("expected error")
	}
}

func TestAnyOrderedMapMarshalJSONDirect(t *testing.T) {
	m := NewAnyOrderedMap()
	m.Set("foo", "bar")
	m.Set("123", 123)

	data, err := m.MarshalJSON()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	const expected = `{"foo":"bar","123":123}`

	if string(data) != expected {
		t.Fatalf("expected %s, got %q", expected, string(data))
	}
}

func TestAnyOrderedMapSetEscapeHTML(t *testing.T) {
	const input = `{"foo":"<bar>","obj":{"baz":"&"},"array":[{"qux":">"}]}`

	var m *AnyOrderedMap
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Applies recursively to nested maps
	m.SetEscapeHTML(false)

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(m); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := strings.TrimSuffix(buf.String(), "\n")

	if output != input {
		t.Fatalf("expected %s, got %s", input, output)
	}
}

func TestAnyOrderedMapUnmarshalNull(t *testing.T) {
	m := NewAnyOrderedMap()
	m.SetUseNumber(true)
	m.Set("foo", "bar")

	// Like encoding/json does for maps, unmarshalling null clears the map
	if err := m.UnmarshalJSON([]byte(`null`)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if m.Len() != 0 {
		t.Fatalf("expected map to be cleared, got %s", m.String())
	}

	if !m.useNumber {
		t.Fatal("expected settings to be kept")
	}

	// The map must remain usable
	m.Set("foo", "bar")

	if m.Len() != 1 {
		t.Fatalf("expected map to be usable, got %s", m.String())
	}
}

func TestAnyOrderedMapSliceUnmarshalNull(t *testing.T) {
	values := AnyOrderedMapSlice{Values: []any{"foo"}}

	// Like encoding/json does for slices, unmarshalling null sets the slice to nil
	if err := json.Unmarshal([]byte(`null`), &values); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if values.Values != nil {
		t.Fatalf("expected slice to be nil, got %s", values.String())
	}
}

func TestAnyOrderedMapString(t *testing.T) {
	const input = `{"a":[1,[2,3],{"b":"c"}]}`

	var m *AnyOrderedMap
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	const expected = `{a:[1,[2,3],{b:c}]}`

	if m.String() != expected {
		t.Fatalf("expected %s, got %s", expected, m.String())
	}
}

func TestAnyOrderedMapSliceString(t *testing.T) {
	const input = `[1,[2,3],{"b":"c"}]`

	var values AnyOrderedMapSlice
	if err := json.Unmarshal([]byte(input), &values); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	const expected = `[1,[2,3],{b:c}]`

	if values.String() != expected {
		t.Fatalf("expected %s, got %s", expected, values.String())
	}
}

func TestAnyOrderedMapSettingsInheritance(t *testing.T) {
	const input = `{"obj":{"foo":"<bar>"},"array":[{"baz":"&"}]}`

	m := NewAnyOrderedMap()
	m.SetEscapeHTML(false)
	m.SetUseNumber(true)

	if err := json.Unmarshal([]byte(input), &m); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	obj, _ := m.Get("obj")
	nested := obj.(*AnyOrderedMap)

	// Nested maps inherit the parent's settings at unmarshal time
	if nested.escapeHTML || !nested.useNumber {
		t.Fatalf("expected nested map to inherit settings, got escapeHTML=%v useNumber=%v", nested.escapeHTML, nested.useNumber)
	}

	data, err := nested.MarshalJSON()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	const expected = `{"foo":"<bar>"}`

	if string(data) != expected {
		t.Fatalf("expected %s, got %s", expected, string(data))
	}

	array, _ := m.Get("array")
	inArray := array.([]any)[0].(*AnyOrderedMap)

	if inArray.escapeHTML || !inArray.useNumber {
		t.Fatalf("expected map in array to inherit settings, got escapeHTML=%v useNumber=%v", inArray.escapeHTML, inArray.useNumber)
	}
}

func TestAnyOrderedMapSliceMarshalNil(t *testing.T) {
	// Like encoding/json does for nil slices, nil Values marshal as null
	var values AnyOrderedMapSlice

	data, err := json.Marshal(values)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(data) != `null` {
		t.Fatalf("expected null, got %s", string(data))
	}

	// Empty but non-nil Values marshal as []
	values.Values = []any{}

	data, err = json.Marshal(values)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(data) != `[]` {
		t.Fatalf("expected [], got %s", string(data))
	}
}

func TestAnyOrderedMapMerge(t *testing.T) {
	m := NewAnyOrderedMapWithCapacity(4)
	m.Set("a", 1)
	m.Set("b", 2)

	// Like encoding/json, unmarshalling into a non-empty map keeps existing entries
	if err := json.Unmarshal([]byte(`{"b":3,"c":4}`), &m); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	const expected = `{a:1,b:3,c:4}`

	if m.String() != expected {
		t.Fatalf("expected %s, got %s", expected, m.String())
	}
}

func TestAnyOrderedMapDuplicateKeys(t *testing.T) {
	const input = `{"a":1,"b":2,"a":3}`

	var m *AnyOrderedMap
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Like encoding/json, the last value wins. The key keeps its original position.
	const expected = `{"a":3,"b":2}`

	data, err := json.Marshal(m)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(data) != expected {
		t.Fatalf("expected %s, got %s", expected, string(data))
	}
}

func TestAnyOrderedMapSliceMarshalJSONDirect(t *testing.T) {
	values := AnyOrderedMapSlice{Values: []any{"foo", 123}}

	data, err := values.MarshalJSON()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	const expected = `["foo",123]`

	if string(data) != expected {
		t.Fatalf("expected %s, got %q", expected, string(data))
	}
}

func TestAnyOrderedMapCopySettings(t *testing.T) {
	m := NewAnyOrderedMap()
	m.SetEscapeHTML(false)
	m.SetUseNumber(true)

	mm := m.Copy()

	if mm.escapeHTML || !mm.useNumber {
		t.Fatalf("expected copy to preserve settings, got escapeHTML=%v useNumber=%v", mm.escapeHTML, mm.useNumber)
	}
}

func TestAnyOrderedMapCopy(t *testing.T) {
	m := NewAnyOrderedMap()
	m.Set("foo", "bar")
	m.Set("123", true)
	m.Set("abc", nil)

	mm := m.Copy()

	if mm.String() != m.String() {
		t.Fatalf("expected %s, got %s", m.String(), mm.String())
	}
}
