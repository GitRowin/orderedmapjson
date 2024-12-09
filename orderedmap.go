package orderedmapjson

import (
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

func (m *orderedMap[V]) SetEscapeHTML(on bool) {
	m.escapeHTML = on
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
