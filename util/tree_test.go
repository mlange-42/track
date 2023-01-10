package util

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
	tr := NewTree(
		testStruct{Name: "root"},
	)
	tr.Add(tr.Root, testStruct{Name: "child"})

	assert.Equal(t, 1, len(tr.Root.Children))

	ch, ok := tr.Root.Children["child"]
	assert.Equal(t, true, ok)
	assert.Equal(t, testStruct{Name: "child"}, ch.Value)

	assert.Equal(t, tr.Root, ch.Parent)
}

func TestNodeAdd(t *testing.T) {
	tr := NewTree(
		testStruct{Name: "root"},
	)
	tr.AddNode(tr.Root, NewNode(testStruct{Name: "child"}))

	assert.Equal(t, 1, len(tr.Root.Children))

	ch, ok := tr.Root.Children["child"]
	assert.Equal(t, true, ok)
	assert.Equal(t, testStruct{Name: "child"}, ch.Value)

	assert.Equal(t, tr.Root, ch.Parent)
}

func TestAncestorDescendants(t *testing.T) {
	tr := NewTree(
		testStruct{Name: "root"},
	)
	root := tr.Root
	a := NewNode(testStruct{Name: "a"})
	a1 := NewNode(testStruct{Name: "a1"})
	a2 := NewNode(testStruct{Name: "a2"})
	b := NewNode(testStruct{Name: "b"})
	b1 := NewNode(testStruct{Name: "b1"})
	b11 := NewNode(testStruct{Name: "b11"})

	tr.AddNode(root, a)
	tr.AddNode(a, a1)
	tr.AddNode(a, a2)

	tr.AddNode(root, b)
	tr.AddNode(b, b1)
	tr.AddNode(b1, b11)

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
		expAncestors   []string
		expDescendants []string
	}{
		{
			title:          "Node root",
			node:           root,
			expAncestors:   []string{},
			expDescendants: []string{"a", "a1", "a2", "b", "b1", "b11"},
		},
		{
			title:          "Node a",
			node:           a,
			expAncestors:   []string{"root"},
			expDescendants: []string{"a1", "a2"},
		},
		{
			title:          "Node a1",
			node:           a1,
			expAncestors:   []string{"a", "root"},
			expDescendants: []string{},
		},
		{
			title:          "Node b",
			node:           b,
			expAncestors:   []string{"root"},
			expDescendants: []string{"b1", "b11"},
		},
		{
			title:          "Node b1",
			node:           b1,
			expAncestors:   []string{"b", "root"},
			expDescendants: []string{"b11"},
		},
		{
			title:          "Node b11",
			node:           b11,
			expAncestors:   []string{"b1", "b", "root"},
			expDescendants: []string{},
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
		ancStr := make([]string, len(anc), len(anc))
		desStr := make([]string, len(des), len(des))
		for i, a := range anc {
			ancStr[i] = a.Value.Name
		}
		for i, a := range des {
			desStr[i] = a.Value.Name
		}
		assert.Equal(t, test.expAncestors, ancStr, "Ancestors don't match in %s", test.title)
		assert.ElementsMatch(t, test.expDescendants, desStr, "Descendants don't match in %s", test.title)
	}
}
