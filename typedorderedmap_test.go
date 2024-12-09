package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestTypedOrderedMapStringInt(t *testing.T) {
	const input = `{"nine":9,"eight":8,"seven":7,"six":6,"five":5,"four":4,"three":3,"two":2,"one":1,"zero":0}`

	m := NewTypedOrderedMap[string, int]()

	if err := json.Unmarshal([]byte(input), m); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fmt.Println(m)

	data, err := json.Marshal(m)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := string(data)

	if output != input {
		t.Fatalf("expected %s, got %s", input, output)
	}
}

func TestTypedOrderedMapStringOrderedMapStringString(t *testing.T) {
	const input = `{"one":{"foo":"bar"},"two":{"foo":"bar"},"three":{"foo":"bar"}}`

	m := NewTypedOrderedMap[string, *TypedOrderedMap[string, string]]()

	if err := json.Unmarshal([]byte(input), m); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fmt.Println(m)

	data, err := json.Marshal(m)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := string(data)

	if output != input {
		t.Fatalf("expected %s, got %s", input, output)
	}
}
