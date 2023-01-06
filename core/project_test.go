package core

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/mlange-42/track/fs"
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

func TestSaveLoadProject(t *testing.T) {
	dir, err := ioutil.TempDir("", "track-test")
	assert.Nil(t, err, "Error creating temporary directory")
	defer os.Remove(dir)

	track, err := NewTrack(&dir)
	assert.Nil(t, err, "Error creating Track instance")

	assert.Equal(t, dir, track.RootDir, "Wrong root directory")

	assert.False(t, fs.FileExists(track.ProjectPath("test")), "File must not exist")
	assert.False(t, track.ProjectExists("test"), "Project must not exist")

	project := NewProject("test", "", "T", 0, 15)
	err = track.SaveProject(project, false)
	assert.Nil(t, err, "Error saving project")

	// Test overwriting
	err = track.SaveProject(project, false)
	assert.True(t, err != nil, "Expect error saving project: should not overwrite")
	err = track.SaveProject(project, true)
	assert.Nil(t, err, "Error saving project with force")

	assert.True(t, fs.FileExists(track.ProjectPath("test")), "File must exist")
	assert.True(t, track.ProjectExists("test"), "Project must exist")

	allProjects, err := track.LoadAllProjects()
	assert.Nil(t, err, "Error loading projects")

	assert.Equal(t, map[string]Project{"test": project}, allProjects, "Loaded project not equal to saved project")

	newProject, err := track.LoadProjectByName("test")
	assert.Nil(t, err, "Error loading project")
	assert.Equal(t, project, newProject, "Loaded project not equal to saved project")

	_, err = track.DeleteProject(&project, true, false)
	assert.Nil(t, err, "Error deleting project")

	assert.False(t, fs.FileExists(track.ProjectPath("test")), "File must not exist")
	assert.False(t, track.ProjectExists("test"), "Project must not exist")
}

func TestCheckParents(t *testing.T) {
	dir, err := ioutil.TempDir("", "track-test")
	assert.Nil(t, err, "Error creating temporary directory")
	defer os.Remove(dir)

	track, err := NewTrack(&dir)
	assert.Nil(t, err, "Error creating Track instance")

	p1 := NewProject("p1", "", "T", 0, 15)
	p2 := NewProject("p2", "", "T", 0, 15)
	p3 := NewProject("p3", "", "T", 0, 15)

	err = track.SaveProject(p1, false)
	assert.Nil(t, err, "Error saving project")
	err = track.SaveProject(p2, false)
	assert.Nil(t, err, "Error saving project")
	err = track.SaveProject(p3, false)
	assert.Nil(t, err, "Error saving project")

	assert.True(t, track.CheckParents(p1) == nil, "Unexpected error in parent check")

	p1.Parent = "p1"
	err = track.SaveProject(p1, true)
	assert.Nil(t, err, "Error saving project")

	assert.True(t, track.CheckParents(p1) != nil, "Expected error in parent check")

	p2.Parent = "p3"
	err = track.SaveProject(p2, true)
	assert.Nil(t, err, "Error saving project")
	p3.Parent = "p2"
	err = track.SaveProject(p3, true)
	assert.Nil(t, err, "Error saving project")

	assert.True(t, track.CheckParents(p2) != nil, "Expected error in parent check")
	assert.True(t, track.CheckParents(p3) != nil, "Expected error in parent check")
}
