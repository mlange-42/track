package util

import "fmt"

// Named is an interface for stuff that has a name
type Named interface {
	GetName() string
}

// MapNode is a node in the tree data structure
type MapNode[T Named] struct {
	Parent   *MapNode[T]
	Children map[string]*MapNode[T]
	Value    T
}

// NewNode creates a new tree node
func NewNode[T Named](value T) *MapNode[T] {
	return &MapNode[T]{
		Children: make(map[string]*MapNode[T]),
		Value:    value,
	}
}

// MapTree is a tree data structure
type MapTree[T Named] struct {
	Root  *MapNode[T]
	Nodes map[string]*MapNode[T]
}

// NewTree creates a new tree node
func NewTree[T Named](value T) *MapTree[T] {
	root := NewNode(value)
	return &MapTree[T]{
		Root:  root,
		Nodes: map[string]*MapNode[T]{value.GetName(): root},
	}
}

// Ancestors returns a slice of all ancestors (i.e. recursive parents),
// and an ok bool whether the requested node was found.
//
// The first ancestor is the direct parent, while the last ancestor is the root node.
func (t *MapTree[T]) Ancestors(name string) ([]*MapNode[T], bool) {
	res := []*MapNode[T]{}
	curr, ok := t.Nodes[name]
	if !ok {
		return res, false
	}
	for curr.Parent != nil {
		res = append(res, curr.Parent)
		curr = curr.Parent
	}
	return res, true
}

// Descendants returns a slice of all descendants (i.e. recursive children),
// and an ok bool whether the requested node was found.
//
// Descendants in the returned slice have undefined order.
func (t *MapTree[T]) Descendants(name string) ([]*MapNode[T], bool) {
	res := []*MapNode[T]{}
	curr, ok := t.Nodes[name]
	if !ok {
		return res, false
	}
	desc := t.descendants(curr, res)
	return desc, true
}

func (t *MapTree[T]) descendants(n *MapNode[T], res []*MapNode[T]) []*MapNode[T] {
	for _, child := range n.Children {
		res = append(res, child)
		res = t.descendants(child, res)
	}
	return res
}

// AddTree adds a sub-tree without children
func (t *MapTree[T]) Add(parent *MapNode[T], child T) (*MapNode[T], error) {
	name := child.GetName()
	if _, ok := t.Nodes[name]; ok {
		return nil, fmt.Errorf("duplicate key '%s'", name)
	}

	node := NewNode(child)
	node.Parent = parent
	parent.Children[name] = node
	t.Nodes[name] = node

	return node, nil
}

// AddNode adds a sub-tree
func (t *MapTree[T]) AddNode(parent *MapNode[T], child *MapNode[T]) error {
	name := child.Value.GetName()
	if _, ok := t.Nodes[name]; ok {
		return fmt.Errorf("duplicate key '%s'", name)
	}

	child.Parent = parent
	parent.Children[name] = child
	t.Nodes[name] = child

	return nil
}

// Aggregate aggregates values over the tree
func Aggregate[T Named, V any](t *MapTree[T], values map[string]V, zero V, fn func(a, b V) V) {
	aggregate(t.Root, values, zero, fn)
}

func aggregate[T Named, V any](nd *MapNode[T], values map[string]V, zero V, fn func(a, b V) V) V {
	agg, ok := values[nd.Value.GetName()]
	if !ok {
		agg = zero
	}
	for _, child := range nd.Children {
		v := aggregate(child, values, zero, fn)
		agg = fn(agg, v)
	}
	values[nd.Value.GetName()] = agg
	return agg
}
