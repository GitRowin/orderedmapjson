package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type TypedOrderedMap[K comparable, V any] struct {
	*OrderedMap[K, V]
}

func NewTypedOrderedMap[K comparable, V any]() *TypedOrderedMap[K, V] {
	return &TypedOrderedMap[K, V]{
		OrderedMap: newOrderedMap[K, V](),
	}
}

func (m *TypedOrderedMap[K, V]) SetEscapeHTML(on bool) {
	m.escapeHTML = on
}

func (m *TypedOrderedMap[K, V]) String() string {
	builder := strings.Builder{}

	builder.WriteString("{")

	index := 0
	for k, v := range m.AllFromFront() {
		if index > 0 {
			builder.WriteString(",")
		}

		builder.WriteString(fmt.Sprintf("%v:%v", k, v))
		index++
	}

	builder.WriteString("}")
	return builder.String()
}

func (m *TypedOrderedMap[K, V]) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(m.escapeHTML)

	index := 0
	for k, v := range m.AllFromFront() {
		if index > 0 {
			buf.WriteByte(',')
		}

		if err := encoder.Encode(k); err != nil {
			return nil, err
		}

		buf.WriteByte(':')

		if err := encoder.Encode(v); err != nil {
			return nil, err
		}

		index++
	}

	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func (m *TypedOrderedMap[K, V]) UnmarshalJSON(b []byte) error {
	if m.OrderedMap == nil {
		m.OrderedMap = newOrderedMap[K, V]()
	}

	decoder := json.NewDecoder(bytes.NewReader(b))

	// Skip '{'
	token, err := decoder.Token()

	if err != nil {
		return err
	}

	if token != json.Delim('{') {
		return fmt.Errorf("expected '{' but got %v", token)
	}

	for {
		token, err := decoder.Token()

		if err != nil {
			return err
		}

		// Reached end of map
		if token == json.Delim('}') {
			return nil
		}

		key, ok := token.(K)

		if !ok {
			return fmt.Errorf("unexpected key type: %T", token)
		}

		var value V
		if err := decoder.Decode(&value); err != nil {
			return err
		}

		m.Set(key, value)
	}
}
