package main

import (
	"encoding/json"
	"testing"
)

func TestAnyOrderedMap(t *testing.T) {
	const input = `{"foo":"bar","123":true,"abc":null,"_obj":{"foo":"bar"},"_array":[{"b":"b","a":"a","0":null,"123":{"q":"","w":"","e":"","r":"","t":"","y":""}}],"q":"","w":"","e":"","r":"","t":"","y":""}`

	m := NewAnyOrderedMap()

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
