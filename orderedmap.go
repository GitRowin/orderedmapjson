package orderedmapjson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elliotchance/orderedmap/v3"
	"strings"
)

type orderedMap[V any] struct {
	*orderedmap.OrderedMap[string, V]
	escapeHTML bool
}

func newOrderedMap[V any]() *orderedMap[V] {
	return &orderedMap[V]{
		OrderedMap: orderedmap.NewOrderedMap[string, V](),
		escapeHTML: true, // Default to true for consistency with encoding/json
	}
}

func newOrderedMapWithCapacity[V any](capacity int) *orderedMap[V] {
	return &orderedMap[V]{
		OrderedMap: orderedmap.NewOrderedMapWithCapacity[string, V](capacity),
		escapeHTML: true, // Default to true for consistency with encoding/json
	}
}

func (m *orderedMap[V]) SetEscapeHTML(on bool) {
	m.escapeHTML = on
}

func (m *orderedMap[V]) MarshalJSON() ([]byte, error) {
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

func (m *orderedMap[V]) String() string {
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
