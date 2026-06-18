package orderedmapjson

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// AnyOrderedMap is a JSON-marshallable ordered map with string keys and values of any type.
// Unlike TypedOrderedMap[any], it unmarshals JSON objects at any level as *AnyOrderedMap
// rather than map[string]any, preserving the order of nested objects too.
// The zero value is not usable, except as an unmarshalling target:
// use NewAnyOrderedMap or NewAnyOrderedMapWithCapacity.
type AnyOrderedMap struct {
	*orderedMap[any]
}

// NewAnyOrderedMap returns an initialized AnyOrderedMap.
func NewAnyOrderedMap() *AnyOrderedMap {
	return &AnyOrderedMap{
		orderedMap: newOrderedMap[any](),
	}
}

// NewAnyOrderedMapWithCapacity returns an initialized AnyOrderedMap with the given capacity.
func NewAnyOrderedMapWithCapacity(capacity int) *AnyOrderedMap {
	return &AnyOrderedMap{
		orderedMap: newOrderedMapWithCapacity[any](capacity),
	}
}

// newNestedAnyOrderedMap returns a new AnyOrderedMap for a nested JSON object,
// inheriting the parent's settings.
func newNestedAnyOrderedMap(escapeHTML, useNumber bool) *AnyOrderedMap {
	m := NewAnyOrderedMap()
	m.escapeHTML = escapeHTML
	m.useNumber = useNumber
	return m
}

// SetEscapeHTML sets whether problematic HTML characters (<, >, and &) are escaped by MarshalJSON.
// The setting is also applied to all nested AnyOrderedMap values.
// Note that json.Marshal HTML-escapes the output of MarshalJSON regardless of this setting.
// To produce unescaped output, use a json.Encoder with SetEscapeHTML(false) or call MarshalJSON directly.
func (m *AnyOrderedMap) SetEscapeHTML(on bool) {
	m.orderedMap.SetEscapeHTML(on)

	for _, v := range m.AllFromFront() {
		setEscapeHTML(v, on)
	}
}

func setEscapeHTML(value any, on bool) {
	switch value := value.(type) {
	case *AnyOrderedMap:
		value.SetEscapeHTML(on)
	case []any:
		for _, v := range value {
			setEscapeHTML(v, on)
		}
	}
}

func (m *AnyOrderedMap) UnmarshalJSON(b []byte) error {
	// Decode null as an empty map, so null in a non-pointer field doesn't error.
	if bytes.Equal(b, jsonNull) {
		if m.orderedMap != nil {
			m.clear()
		}

		return nil
	}

	if m.orderedMap == nil {
		m.orderedMap = newOrderedMap[any]()
	}

	decoder := json.NewDecoder(bytes.NewReader(b))

	if m.useNumber {
		decoder.UseNumber()
	}

	// Skip '{'.
	token, err := decoder.Token()

	if err != nil {
		return err
	}

	if token != json.Delim('{') {
		return fmt.Errorf("expected '{' but got %v", token)
	}

	if err := unmarshalAnyOrderedMap(decoder, m); err != nil {
		return err
	}

	return ensureNoTrailingData(decoder)
}

// Copy returns a shallow copy of the map: the values are copied as-is,
// so nested maps and slices are shared between the original and the copy.
func (m *AnyOrderedMap) Copy() *AnyOrderedMap {
	return &AnyOrderedMap{
		orderedMap: m.orderedMap.Copy(),
	}
}

func unmarshalAnyOrderedMap(decoder *json.Decoder, m *AnyOrderedMap) error {
	for {
		token, err := decoder.Token()

		if err != nil {
			return err
		}

		// Reached end of map.
		if token == json.Delim('}') {
			return nil
		}

		key, ok := token.(string)

		if !ok {
			return fmt.Errorf("unexpected key type: %T", token)
		}

		token, err = decoder.Token()

		if err != nil {
			return err
		}

		switch token {
		case json.Delim('{'):
			mm := newNestedAnyOrderedMap(m.escapeHTML, m.useNumber)

			if err := unmarshalAnyOrderedMap(decoder, mm); err != nil {
				return err
			}

			m.Set(key, mm)
		case json.Delim('['):
			values, err := unmarshalAnyOrderedMapArray(decoder, m.escapeHTML, m.useNumber)

			if err != nil {
				return err
			}

			m.Set(key, values)
		default:
			m.Set(key, token)
		}
	}
}

func unmarshalAnyOrderedMapArray(decoder *json.Decoder, escapeHTML, useNumber bool) ([]any, error) {
	values := []any{}
	for {
		token, err := decoder.Token()

		if err != nil {
			return nil, err
		}

		switch token {
		// Reached end of array.
		case json.Delim(']'):
			return values, nil
		case json.Delim('{'):
			mm := newNestedAnyOrderedMap(escapeHTML, useNumber)

			if err := unmarshalAnyOrderedMap(decoder, mm); err != nil {
				return nil, err
			}

			values = append(values, mm)
		case json.Delim('['):
			vv, err := unmarshalAnyOrderedMapArray(decoder, escapeHTML, useNumber)

			if err != nil {
				return nil, err
			}

			values = append(values, vv)
		default:
			values = append(values, token)
		}
	}
}

// AnyOrderedMapSlice unmarshals JSON arrays, decoding JSON objects at any level
// as *AnyOrderedMap rather than map[string]any. The zero value is ready to use.
type AnyOrderedMapSlice struct {
	Values []any

	noEscapeHTML bool // Inverted so that the zero value matches encoding/json's default.
	useNumber    bool
}

// SetEscapeHTML sets whether problematic HTML characters (<, >, and &) are escaped by MarshalJSON.
// The setting is also applied to all nested AnyOrderedMap values.
// Note that json.Marshal HTML-escapes the output of MarshalJSON regardless of this setting.
// To produce unescaped output, use a json.Encoder with SetEscapeHTML(false) or call MarshalJSON directly.
func (s *AnyOrderedMapSlice) SetEscapeHTML(on bool) {
	s.noEscapeHTML = !on

	for _, v := range s.Values {
		setEscapeHTML(v, on)
	}
}

// SetUseNumber sets whether UnmarshalJSON decodes numbers as json.Number instead of float64,
// like json.Decoder's UseNumber.
func (s *AnyOrderedMapSlice) SetUseNumber(on bool) {
	s.useNumber = on
}

func (s AnyOrderedMapSlice) MarshalJSON() ([]byte, error) {
	// Like encoding/json does for nil slices, marshal nil Values as null.
	if s.Values == nil {
		return []byte("null"), nil
	}

	var buf bytes.Buffer
	buf.WriteByte('[')
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(!s.noEscapeHTML)

	for i, v := range s.Values {
		if i > 0 {
			buf.WriteByte(',')
		}

		if err := encoder.Encode(v); err != nil {
			return nil, err
		}

		// Remove the newline added by Encode.
		buf.Truncate(buf.Len() - 1)
	}

	buf.WriteByte(']')
	return buf.Bytes(), nil
}

func (s *AnyOrderedMapSlice) UnmarshalJSON(b []byte) error {
	// Like encoding/json does for slices, unmarshalling null sets the slice to nil.
	if bytes.Equal(b, jsonNull) {
		s.Values = nil
		return nil
	}

	decoder := json.NewDecoder(bytes.NewReader(b))

	if s.useNumber {
		decoder.UseNumber()
	}

	// Skip '['.
	token, err := decoder.Token()

	if err != nil {
		return err
	}

	if token != json.Delim('[') {
		return fmt.Errorf("expected '[' but got %v", token)
	}

	values, err := unmarshalAnyOrderedMapArray(decoder, !s.noEscapeHTML, s.useNumber)

	if err != nil {
		return err
	}

	if err := ensureNoTrailingData(decoder); err != nil {
		return err
	}

	s.Values = values
	return nil
}

func (s AnyOrderedMapSlice) String() string {
	return formatValue(s.Values)
}
