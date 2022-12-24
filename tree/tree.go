package tree

// Named is the interface for named nodes in a MapTree
type Named interface {
	GetName() string
}

// MapTree is a tree data structure
type MapTree[T Named] struct {
	Children map[string]*MapTree[T]
	Value    T
}

// New creates a new tree
func New[T Named](value T) *MapTree[T] {
	return &MapTree[T]{
		Children: make(map[string]*MapTree[T]),
		Value:    value,
	}
}

// AddTree adds a sub-tree without children
func (t *MapTree[T]) Add(child T) {
	t.Children[child.GetName()] = New(child)
}

// AddTree adds a sub-tree
func (t *MapTree[T]) AddTree(child *MapTree[T]) {
	t.Children[child.Value.GetName()] = child
}
