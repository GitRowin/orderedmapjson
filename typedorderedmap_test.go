package orderedmapjson

import (
	"encoding/json"
	"testing"
)

func TestTypedOrderedMapInt(t *testing.T) {
	const input = `{"nine":9,"eight":8,"seven":7,"six":6,"five":5,"four":4,"three":3,"two":2,"one":1,"zero":0}`

	var m *TypedOrderedMap[int]
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

func TestTypedOrderedMapOrderedMapString(t *testing.T) {
	const input = `{"one":{"foo":"bar"},"two":{"foo":"bar"},"three":{"foo":"bar"}}`

	var m *TypedOrderedMap[*TypedOrderedMap[string]]
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

func TestTypedOrderedMapOrderedMapInt(t *testing.T) {
	const input = `{"one":{"foo":1},"two":{"foo":2},"three":{"foo":3}}`

	var m *TypedOrderedMap[*TypedOrderedMap[int]]
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

func TestTypedOrderedMapIntWrongType(t *testing.T) {
	const input = `{"one":1,"two":2,"three":"3"}`

	var m *TypedOrderedMap[int]
	if err := json.Unmarshal([]byte(input), &m); err == nil {
		t.Fatal("expected error")
	}
}

func TestTypedOrderedMapUseNumber(t *testing.T) {
	const input = `{"foo":10000000000000000001}`

	m := NewTypedOrderedMap[any]()
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

func TestTypedOrderedMapTrailingData(t *testing.T) {
	m := NewTypedOrderedMap[int]()

	if err := m.UnmarshalJSON([]byte(`{"foo":1} {"bar":2}`)); err == nil {
		t.Fatal("expected error")
	}
}

func TestTypedOrderedMapUnmarshalNull(t *testing.T) {
	m := NewTypedOrderedMap[int]()
	m.Set("foo", 1)

	// Like encoding/json does for maps, unmarshalling null clears the map
	if err := m.UnmarshalJSON([]byte(`null`)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if m.Len() != 0 {
		t.Fatalf("expected map to be cleared, got %s", m.String())
	}

	// The map must remain usable
	m.Set("foo", 1)

	if m.Len() != 1 {
		t.Fatalf("expected map to be usable, got %s", m.String())
	}
}

func TestTypedOrderedMapNullField(t *testing.T) {
	// Non-pointer fields receive "null" in their UnmarshalJSON
	var holder struct {
		M TypedOrderedMap[int]
		S AnyOrderedMapSlice
	}

	if err := json.Unmarshal([]byte(`{"M":null,"S":null}`), &holder); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTypedOrderedMapInvalidInput(t *testing.T) {
	m := NewTypedOrderedMap[int]()

	for _, input := range []string{``, `[1]`, `123`, `{`, `{"a":`, `{"a":1`, `{"a":1,`, `{"a":1}x`} {
		if err := m.UnmarshalJSON([]byte(input)); err == nil {
			t.Fatalf("expected error for input %q", input)
		}
	}
}

func TestTypedOrderedMapMerge(t *testing.T) {
	m := NewTypedOrderedMapWithCapacity[int](4)
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

func TestTypedOrderedMapCopy(t *testing.T) {
	m := NewTypedOrderedMap[string]()
	m.Set("foo", "bar")
	m.Set("123", "true")
	m.Set("abc", "nil")

	mm := m.Copy()

	if mm.String() != m.String() {
		t.Fatalf("expected %s, got %s", m.String(), mm.String())
	}
}
