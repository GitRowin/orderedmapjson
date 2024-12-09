package orderedmapjson

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type TypedOrderedMap[V any] struct {
	*orderedMap[V]
}

func NewTypedOrderedMap[V any]() *TypedOrderedMap[V] {
	return &TypedOrderedMap[V]{
		orderedMap: newOrderedMap[V](),
	}
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
