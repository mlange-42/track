package tree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Name string
}

func (t testStruct) GetName() string {
	return t.Name
}

type testTree = MapTree[testStruct]
type testNode = MapNode[testStruct]

func TestAdd(t *testing.T) {
	node := NewNode(testStruct{Name: "root"})
	node.Add(testStruct{Name: "child"})

	assert.Equal(t, 1, len(node.Children))

	ch, ok := node.Children["child"]
	assert.Equal(t, true, ok)
	assert.Equal(t, testStruct{Name: "child"}, ch.Value)

	assert.Equal(t, node, ch.Parent)
}

func TestNodeAdd(t *testing.T) {
	node := NewNode(testStruct{Name: "root"})
	node.AddNode(NewNode(testStruct{Name: "child"}))

	assert.Equal(t, 1, len(node.Children))

	ch, ok := node.Children["child"]
	assert.Equal(t, true, ok)
	assert.Equal(t, testStruct{Name: "child"}, ch.Value)

	assert.Equal(t, node, ch.Parent)
}

func TestAncestorDescendants(t *testing.T) {
	root := NewNode(testStruct{Name: "root"})
	a := NewNode(testStruct{Name: "a"})
	a1 := NewNode(testStruct{Name: "a1"})
	a2 := NewNode(testStruct{Name: "a2"})
	b := NewNode(testStruct{Name: "b"})
	b1 := NewNode(testStruct{Name: "b1"})
	b11 := NewNode(testStruct{Name: "b11"})

	root.AddNode(a)
	a.AddNode(a1)
	a.AddNode(a2)

	root.AddNode(b)
	b.AddNode(b1)
	b1.AddNode(b11)

	tree := MapTree[testStruct]{
		Root: root,
		Nodes: map[string]*testNode{
			"root": root,
			"a":    a,
			"a1":   a1,
			"a2":   a2,
			"b":    b,
			"b1":   b1,
			"b11":  b11,
		},
	}

	tt := []struct {
		title          string
		node           *testNode
		expAncestors   []*testNode
		expDescendants []*testNode
	}{
		{
			title:          "Node root",
			node:           root,
			expAncestors:   []*testNode{},
			expDescendants: []*testNode{a, a1, a2, b, b1, b11},
		},
		{
			title:          "Node a",
			node:           a,
			expAncestors:   []*testNode{root},
			expDescendants: []*testNode{a1, a2},
		},
		{
			title:          "Node a1",
			node:           a1,
			expAncestors:   []*testNode{a, root},
			expDescendants: []*testNode{},
		},
		{
			title:          "Node b",
			node:           b,
			expAncestors:   []*testNode{root},
			expDescendants: []*testNode{b1, b11},
		},
		{
			title:          "Node b1",
			node:           b1,
			expAncestors:   []*testNode{b, root},
			expDescendants: []*testNode{b11},
		},
		{
			title:          "Node b11",
			node:           b11,
			expAncestors:   []*testNode{b1, b, root},
			expDescendants: []*testNode{},
		},
	}

	for _, test := range tt {
		anc, ok := tree.Ancestors(test.node.Value.Name)
		if !ok {
			t.Fatalf("Should be able to determine ancestors")
		}
		des, ok := tree.Descendants(test.node.Value.Name)
		if !ok {
			t.Fatalf("Should be able to determine descendants")
		}
		assert.Equal(t, test.expAncestors, anc, "Ancestors don't match in %s", test.title)
		assert.Equal(t, test.expDescendants, des, "Descendants don't match in %s", test.title)
	}
}
