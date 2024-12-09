package orderedmapjson

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type AnyOrderedMap struct {
	*orderedMap[any]
}

func NewAnyOrderedMap() *AnyOrderedMap {
	return &AnyOrderedMap{
		orderedMap: newOrderedMap[any](),
	}
}

func (m *AnyOrderedMap) UnmarshalJSON(b []byte) error {
	if m.orderedMap == nil {
		m.orderedMap = newOrderedMap[any]()
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

	return unmarshalAnyOrderedMap(decoder, m)
}

func unmarshalAnyOrderedMap(decoder *json.Decoder, m *AnyOrderedMap) error {
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

		token, err = decoder.Token()

		if err != nil {
			return err
		}

		switch token {
		case json.Delim('{'):
			mm := NewAnyOrderedMap()

			if err := unmarshalAnyOrderedMap(decoder, mm); err != nil {
				return err
			}

			m.Set(key, mm)
		case json.Delim('['):
			values, err := unmarshalAnyOrderedMapArray(decoder)

			if err != nil {
				return err
			}

			m.Set(key, values)
		default:
			m.Set(key, token)
		}
	}
}

func unmarshalAnyOrderedMapArray(decoder *json.Decoder) ([]any, error) {
	var values []any
	for {
		token, err := decoder.Token()

		if err != nil {
			return values, err
		}

		switch token {
		// Reached end of array
		case json.Delim(']'):
			return values, nil
		case json.Delim('{'):
			mm := NewAnyOrderedMap()

			if err := unmarshalAnyOrderedMap(decoder, mm); err != nil {
				return nil, err
			}

			values = append(values, mm)
		case json.Delim('['):
			vv, err := unmarshalAnyOrderedMapArray(decoder)

			if err != nil {
				return values, err
			}

			values = append(values, vv)
		default:
			values = append(values, token)
		}
	}
}
