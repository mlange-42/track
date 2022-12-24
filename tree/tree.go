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

// MapNode is a node in the tree data structure
type MapNode[T Named] struct {
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
	t.Children[child.GetName()] = NewNode(child)
}

// AddTree adds a sub-tree
func (t *MapNode[T]) AddTree(child *MapNode[T]) {
	t.Children[child.Value.GetName()] = child
}

// AddTree adds a sub-tree
func (t *MapNode[T]) Find(name string) (*MapNode[T], bool) {
	if t.Value.GetName() == name {
		return t, true
	}
	if tr, ok := t.Children[name]; ok {
		return tr, true
	}
	for _, child := range t.Children {
		if tr, ok := child.Find(name); ok {
			return tr, ok
		}
	}
	return nil, false
}
