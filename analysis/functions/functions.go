package functions

import (
	"go/ast"

	"github.com/saintedlama/goarch/analysis/common"
)

// Item represents a function or method declaration entry.
type Item struct {
	Ref      common.Ref
	Name     string
	Receiver string
	Node     *ast.FuncDecl
}

// MatchFunc is a function type that matches function entries.
type MatchFunc func(Item) bool

// Collection stores function entries and provides convenience query APIs.
type Collection struct {
	items []Item
}

// All returns all function entries.
func (c Collection) All() []Item {
	return c.items
}

// Len returns number of function entries.
func (c Collection) Len() int {
	return len(c.items)
}

// Match applies matcher to all function entries and converts matches into code refs.
func (c Collection) Match(matcher MatchFunc) []common.Ref {
	if matcher == nil {
		return nil
	}

	var refs []common.Ref
	for _, item := range c.items {
		if !matcher(item) {
			continue
		}
		refs = append(refs, item.Ref)
	}

	return refs
}

// Add appends an entry to the collection.
func (c *Collection) Add(item Item) {
	c.items = append(c.items, item)
}
