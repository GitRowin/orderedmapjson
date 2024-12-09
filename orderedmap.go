package orderedmapjson

import (
	"fmt"
	"github.com/elliotchance/orderedmap/v3"
	"strings"
)

type OrderedMap[K comparable, V any] struct {
	*orderedmap.OrderedMap[K, V]
	escapeHTML bool
}

func newOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		OrderedMap: orderedmap.NewOrderedMap[K, V](),
		escapeHTML: true, // Default to true for consistency with encoding/json
	}
}

func (m *OrderedMap[K, V]) SetEscapeHTML(on bool) {
	m.escapeHTML = on
}

func (m *OrderedMap[K, V]) String() string {
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
