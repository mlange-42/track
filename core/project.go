package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gookit/color"
	"github.com/mlange-42/track/util"
	"gopkg.in/yaml.v3"
)

// RootPattern is the pattern to display the workspace in the project tree
const RootPattern = "<%s>"

// ProjectTree is a tree of projects
type ProjectTree = util.MapTree[Project]

// NewTree creates a new project tree
func NewTree(project Project) *ProjectTree {
	return util.NewTree(project)
}

// ProjectNode is a tree of projects
type ProjectNode = util.MapNode[Project]

// Project holds and manipulates data for a project
type Project struct {
	Name         string
	Parent       string
	RequiredTags []string `yaml:"requiredTags"`
	Color        uint8
	FgColor      uint8          `yaml:"fgColor"`
	Render       color.Style256 `yaml:"-"`
	Symbol       string
	Archived     bool
}

// NewProject creates a new project
func NewProject(name string, parent string, symbol string, requiredTags []string, fgColor, color uint8) Project {
	p := Project{
		Name:         name,
		Parent:       parent,
		RequiredTags: requiredTags,
		Symbol:       symbol,
		Color:        color,
		FgColor:      fgColor,
		Archived:     false,
	}
	p.SetColors(fgColor, color)
	return p
}

type tempProject struct {
	Name         string
	Parent       string
	RequiredTags []string `yaml:"requiredTags"`
	Color        uint8
	FgColor      uint8 `yaml:"fgColor"`
	Symbol       string
	Archived     bool
}

// GetName implements the Named interface required for the MapTree
func (p Project) GetName() string {
	return p.Name
}

// UnmarshalYAML un-marshals a project
func (p *Project) UnmarshalYAML(value *yaml.Node) error {
	var tmp tempProject
	err := value.Decode(&tmp)
	if err != nil {
		return err
	}
	p.Name = tmp.Name
	p.Parent = tmp.Parent
	p.RequiredTags = tmp.RequiredTags
	p.Symbol = tmp.Symbol
	p.Archived = tmp.Archived

	p.SetColors(tmp.FgColor, tmp.Color)

	return nil
}

// SetColors sets project colors
func (p *Project) SetColors(fgCol, col uint8) {
	p.Color = col
	p.FgColor = fgCol
	p.Render = *color.S256(fgCol, col)
}

// ProjectExists checks if a project exists on disk
func (t *Track) ProjectExists(name string) bool {
	return util.FileExists(t.ProjectPath(name))
}

// SaveProject saves a project to disk.
// Argument `force` allows to overwrite an existing file.
func (t *Track) SaveProject(project Project, force bool) error {
	path := t.ProjectPath(project.Name)

	if !force && util.FileExists(path) {
		return fmt.Errorf("Project '%s' already exists", project.Name)
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	bytes, err := yaml.Marshal(&project)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(file, "%s Project %s\n\n", YamlCommentPrefix, project.Name)
	if err != nil {
		return err
	}

	_, err = file.Write(bytes)

	return err
}

// LoadProject loads a project by it's name
func (t *Track) LoadProject(name string) (Project, error) {
	path := t.ProjectPath(name)
	return t.loadProjectFromFile(path)
}

// loadProjectFromFile loads a project from the given path
func (t *Track) loadProjectFromFile(path string) (Project, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return Project{}, err
	}

	var project Project

	if err := yaml.Unmarshal(file, &project); err != nil {
		return Project{}, err
	}

	return project, nil
}

// LoadAllProjects loads all projects in the current workspace
func (t *Track) LoadAllProjects() (map[string]Project, error) {
	path := t.ProjectsDir()

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	projects := make(map[string]Project)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		project, err := t.loadProjectFromFile(filepath.Join(path, file.Name()))
		if err != nil {
			return nil, err
		}
		projects[project.Name] = project
	}

	return projects, nil
}

// DeleteProject deletes a project and potentially all associated records
//
// Argument `deleteRecords` determines whether records should be deleted
// Argument `dryRun` can be used to dry-run deleting
func (t *Track) DeleteProject(project *Project, deleteRecords bool, dryRun bool) (int, error) {
	counter := 0

	if deleteRecords {
		// TODO: make a backup
		filters := NewFilter(
			[]func(r *Record) bool{
				FilterByProjects([]string{project.Name}),
			}, util.NoTime, util.NoTime,
		)
		fn, results, _ := t.AllRecordsFiltered(filters, false)

		go fn()

		for res := range results {
			if res.Err != nil {
				return counter, res.Err
			}
			if !dryRun {
				t.DeleteRecord(&res.Record)
			}
			counter++
		}
	}

	if !dryRun {
		err := os.Remove(t.ProjectPath(project.Name))
		if err != nil {
			return counter, err
		}
	}

	return counter, nil
}

// ToProjectTree creates a MapTree of the given projects
func (t *Track) ToProjectTree(projects map[string]Project) (*ProjectTree, error) {
	pTree := NewTree(
		Project{
			Name:   t.WorkspaceLabel(),
			Symbol: " ",
		},
	)

	nodes := map[string]*ProjectNode{pTree.Root.Value.Name: pTree.Root}

	for name, project := range projects {
		nodes[name] = util.NewNode(project)
	}

	for _, tree := range nodes {
		if tree == pTree.Root {
			continue
		}
		var err error
		if tree.Value.Parent == "" {
			err = pTree.AddNode(pTree.Root, tree)
		} else {
			if tt, ok := nodes[tree.Value.Parent]; ok {
				err = pTree.AddNode(tt, tree)
			} else {
				err = pTree.AddNode(pTree.Root, tree)
			}
		}
		if err != nil {
			return nil, err
		}
	}

	return pTree, nil
}

// CheckParents tests if the parent project is valid and that there are no circular relations
func (t *Track) CheckParents(p Project) error {
	projects, err := t.LoadAllProjects()
	if err != nil {
		return err
	}
	return t.checkParentsRecursive(p, p, projects)
}

func (t *Track) checkParentsRecursive(p Project, start Project, projects map[string]Project) error {
	if p.Parent == "" {
		return nil
	}
	if p.Parent == p.Name {
		return fmt.Errorf("can't make project '%s' a parent of itself", p.Parent)
	}
	parent, ok := projects[p.Parent]
	if !ok {
		return fmt.Errorf("project '%s' does not exist", p.Parent)
	}
	if parent.Name == start.Name {
		return fmt.Errorf("circular parent relationship")
	}
	return t.checkParentsRecursive(parent, start, projects)
}
