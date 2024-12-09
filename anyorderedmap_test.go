package orderedmapjson

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

func TestUnmarshalArrayWithAnyOrderedMap(t *testing.T) {
	const input = `[{"c":"c","a":"a","b":"b"},"test",123,{"3":"3","2":"2","1":"1"}]`

	values, err := UnmarshalArrayWithAnyOrderedMap([]byte(input))

	data, err := json.Marshal(values)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := string(data)

	if output != input {
		t.Fatalf("expected %s, got %s", input, output)
	}
}
