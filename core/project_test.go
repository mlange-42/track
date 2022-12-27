package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToTree(t *testing.T) {
	track := Track{
		Config: Config{
			Workspace: "default",
		},
	}

	projects := map[string]Project{
		"p1": {
			Name: "p1",
		},
		"p1a": {
			Name:   "p1a",
			Parent: "p1",
		},
		"p1b": {
			Name:   "p1b",
			Parent: "p1",
		},
		"p2": {
			Name: "p2",
		},
		"p2a": {
			Name:   "p2a",
			Parent: "p2",
		},
	}

	pTree, err := track.ToProjectTree(projects)
	if err != nil {
		t.Fatal(err)
	}

	tt := []struct {
		title          string
		node           *ProjectNode
		expAncestors   []*ProjectNode
		expDescendants []*ProjectNode
	}{
		{
			title: "Node p1",
			node:  pTree.Nodes["p1"],
			expAncestors: []*ProjectNode{
				pTree.Root,
			},
			expDescendants: []*ProjectNode{
				pTree.Nodes["p1a"],
				pTree.Nodes["p1b"],
			},
		},
		{
			title: "Node p1a",
			node:  pTree.Nodes["p1a"],
			expAncestors: []*ProjectNode{
				pTree.Nodes["p1"],
				pTree.Root,
			},
			expDescendants: []*ProjectNode{},
		},
		{
			title: "Node p1b",
			node:  pTree.Nodes["p1b"],
			expAncestors: []*ProjectNode{
				pTree.Nodes["p1"],
				pTree.Root,
			},
			expDescendants: []*ProjectNode{},
		},
		{
			title: "Node p2",
			node:  pTree.Nodes["p2"],
			expAncestors: []*ProjectNode{
				pTree.Root,
			},
			expDescendants: []*ProjectNode{
				pTree.Nodes["p2a"],
			},
		},
		{
			title: "Node p2a",
			node:  pTree.Nodes["p2a"],
			expAncestors: []*ProjectNode{
				pTree.Nodes["p2"],
				pTree.Root,
			},
			expDescendants: []*ProjectNode{},
		},
	}

	for _, test := range tt {
		anc, ok := pTree.Ancestors(test.node.Value.Name)
		if !ok {
			t.Fatalf("Should be able to determine ancestors")
		}
		des, ok := pTree.Descendants(test.node.Value.Name)
		if !ok {
			t.Fatalf("Should be able to determine descendants")
		}
		assert.Equal(t, test.expAncestors, anc, "Ancestors don't match in %s", test.title)
		assert.ElementsMatch(t, test.expDescendants, des, "Descendants don't match in %s", test.title)
	}
}
