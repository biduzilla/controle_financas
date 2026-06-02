package cache

import (
	"fmt"
	"strings"
)

type KeyBuilder struct {
	prefix string
}

func NewKeyBuilder(prefix string) *KeyBuilder {
	return &KeyBuilder{prefix: prefix}
}

func (kb *KeyBuilder) Item(id string) string {
	return fmt.Sprintf("%s:%s", kb.prefix, id)
}

func (kb *KeyBuilder) List(params ...any) string {
	var base strings.Builder
	fmt.Fprintf(&base, "%s:list", kb.prefix)
	for _, p := range params {
		fmt.Fprintf(&base, ":%v", p)
	}
	return base.String()
}
