package orderedmapjson

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// TypedOrderedMap is a JSON-marshallable ordered map with string keys and values of type V.
// The zero value is not usable, except as an unmarshalling target:
// use NewTypedOrderedMap or NewTypedOrderedMapWithCapacity.
type TypedOrderedMap[V any] struct {
	*orderedMap[V]
}

// NewTypedOrderedMap returns an initialized TypedOrderedMap.
func NewTypedOrderedMap[V any]() *TypedOrderedMap[V] {
	return &TypedOrderedMap[V]{
		orderedMap: newOrderedMap[V](),
	}
}

// NewTypedOrderedMapWithCapacity returns an initialized TypedOrderedMap with the given capacity.
func NewTypedOrderedMapWithCapacity[V any](capacity int) *TypedOrderedMap[V] {
	return &TypedOrderedMap[V]{
		orderedMap: newOrderedMapWithCapacity[V](capacity),
	}
}

func (m *TypedOrderedMap[V]) UnmarshalJSON(b []byte) error {
	// Decode null as an empty map, so null in a non-pointer field doesn't error.
	if bytes.Equal(b, jsonNull) {
		if m.orderedMap != nil {
			m.clear()
		}

		return nil
	}

	if m.orderedMap == nil {
		m.orderedMap = newOrderedMap[V]()
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

	for {
		token, err := decoder.Token()

		if err != nil {
			return err
		}

		// Reached end of map.
		if token == json.Delim('}') {
			break
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

	return ensureNoTrailingData(decoder)
}

// Copy returns a shallow copy of the map: the values are copied as-is,
// so nested maps and slices are shared between the original and the copy.
func (m *TypedOrderedMap[V]) Copy() *TypedOrderedMap[V] {
	return &TypedOrderedMap[V]{
		orderedMap: m.orderedMap.Copy(),
	}
}
