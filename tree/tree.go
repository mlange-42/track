package tree

// Named is the interface for named nodes in a MapTree
type Named interface {
	GetName() string
}

// MapTree is a tree data structure
type MapTree[T Named] struct {
	Root  *MapNode[T]
	Nodes map[string]*MapNode[T]
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

// AddTree adds a sub-tree without children
func (t *MapNode[T]) Add(child T) {
	node := NewNode(child)
	node.Parent = t
	t.Children[child.GetName()] = node
}

// AddNode adds a sub-tree
func (t *MapNode[T]) AddNode(child *MapNode[T]) {
	child.Parent = t
	t.Children[child.Value.GetName()] = child
}
