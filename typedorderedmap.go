package orderedmapjson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type TypedOrderedMap[V any] struct {
	*orderedMap[V]
}

func NewTypedOrderedMap[V any]() *TypedOrderedMap[V] {
	return &TypedOrderedMap[V]{
		orderedMap: newOrderedMap[V](),
	}
}

func (m *TypedOrderedMap[V]) SetEscapeHTML(on bool) {
	m.escapeHTML = on
}

func (m *TypedOrderedMap[V]) String() string {
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

func (m *TypedOrderedMap[V]) MarshalJSON() ([]byte, error) {
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

func (m *TypedOrderedMap[V]) UnmarshalJSON(b []byte) error {
	if m.orderedMap == nil {
		m.orderedMap = newOrderedMap[V]()
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

		key, ok := token.(string)

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
