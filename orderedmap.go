package orderedmapjson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/elliotchance/orderedmap/v3"
)

type orderedMap[V any] struct {
	*orderedmap.OrderedMap[string, V]
	escapeHTML bool
	useNumber  bool
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

// SetEscapeHTML sets whether problematic HTML characters (<, >, and &) are escaped by MarshalJSON.
// Note that json.Marshal HTML-escapes the output of MarshalJSON regardless of this setting.
// To produce unescaped output, use a json.Encoder with SetEscapeHTML(false) or call MarshalJSON directly.
func (m *orderedMap[V]) SetEscapeHTML(on bool) {
	m.escapeHTML = on
}

// SetUseNumber sets whether UnmarshalJSON decodes numbers as json.Number instead of float64,
// like json.Decoder's UseNumber.
func (m *orderedMap[V]) SetUseNumber(on bool) {
	m.useNumber = on
}

// clear removes all entries from the map, keeping its settings.
func (m *orderedMap[V]) clear() {
	m.OrderedMap = orderedmap.NewOrderedMap[string, V]()
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

		// Remove the newline added by Encode
		buf.Truncate(buf.Len() - 1)

		buf.WriteByte(':')

		if err := encoder.Encode(v); err != nil {
			return nil, err
		}

		// Remove the newline added by Encode
		buf.Truncate(buf.Len() - 1)

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

		builder.WriteString(fmt.Sprintf("%v:%s", k, formatValue(v)))
		index++
	}

	builder.WriteString("}")
	return builder.String()
}

// formatValue formats a value for String, formatting nested []any values
// consistently with the map and slice String methods.
func formatValue(value any) string {
	values, ok := value.([]any)

	if !ok {
		return fmt.Sprintf("%v", value)
	}

	builder := strings.Builder{}

	builder.WriteString("[")

	for i, v := range values {
		if i > 0 {
			builder.WriteString(",")
		}

		builder.WriteString(formatValue(v))
	}

	builder.WriteString("]")
	return builder.String()
}

// Copy returns a shallow copy of the map: the values are copied as-is,
// so nested maps and slices are shared between the original and the copy.
func (m *orderedMap[V]) Copy() *orderedMap[V] {
	return &orderedMap[V]{
		OrderedMap: m.OrderedMap.Copy(),
		escapeHTML: m.escapeHTML,
		useNumber:  m.useNumber,
	}
}

var jsonNull = []byte("null")

// ensureNoTrailingData returns an error if the decoder has any tokens left.
func ensureNoTrailingData(decoder *json.Decoder) error {
	token, err := decoder.Token()

	if err == io.EOF {
		return nil
	}

	if err != nil {
		return err
	}

	return fmt.Errorf("unexpected data after top-level value: %v", token)
}
