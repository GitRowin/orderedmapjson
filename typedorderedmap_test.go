package orderedmapjson

import (
	"encoding/json"
	"testing"
)

func TestTypedOrderedMapInt(t *testing.T) {
	const input = `{"nine":9,"eight":8,"seven":7,"six":6,"five":5,"four":4,"three":3,"two":2,"one":1,"zero":0}`

	m := NewTypedOrderedMap[int]()

	if err := json.Unmarshal([]byte(input), m); err != nil {
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

	m := NewTypedOrderedMap[*TypedOrderedMap[string]]()

	if err := json.Unmarshal([]byte(input), m); err != nil {
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

	m := NewTypedOrderedMap[*TypedOrderedMap[int]]()

	if err := json.Unmarshal([]byte(input), m); err != nil {
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

	m := NewTypedOrderedMap[int]()

	if err := json.Unmarshal([]byte(input), m); err == nil {
		t.Fatal("expected error")
	}
}

func TestTypedOrderedMapCopy(t *testing.T) {
	m := NewTypedOrderedMap[string]()
	m.Set("foo", "bar")
	m.Set("123", "true")
	m.Set("abc", "nil")

	mm := m.Copy()

	if mm.String() != m.String() {
		t.Fatalf("expected %s, got %s", mm.String(), m.String())
	}
}
