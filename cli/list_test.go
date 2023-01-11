package cli

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/mlange-42/track/util"
	"github.com/stretchr/testify/assert"
)

func TestListWorkspaces(t *testing.T) {
	track, err := setupTestCommand()
	if err != nil {
		t.Fatal("error setting up test")
	}
	defer os.Remove(track.RootDir)

	cmd := RootCommand(track, "")
	cmd.SetArgs([]string{"list", "workspaces"})

	buffer := bytes.NewBufferString("")
	out.StdOut = buffer
	err = cmd.Execute()
	if err != nil {
		t.Fatal("error executing command")
	}

	outStr, err := io.ReadAll(buffer)
	if err != nil {
		t.Fatal("error reading output")
	}

	got := string(outStr)

	assert.Contains(t, got, "default", "First line should contain workspace")
}

func TestListProjects(t *testing.T) {
	track, err := setupTestCommand()
	if err != nil {
		t.Fatal("error setting up test")
	}
	defer os.Remove(track.RootDir)

	project := core.NewProject("test", "", "t", []string{}, 15, 0)
	err = track.SaveProject(project, false)
	if err != nil {
		t.Fatal("error saving project")
	}
	child := core.NewProject("child", "test", "t", []string{}, 15, 0)
	err = track.SaveProject(child, false)
	if err != nil {
		t.Fatal("error saving project")
	}

	cmd := RootCommand(track, "")
	cmd.SetArgs([]string{"list", "projects"})

	buffer := bytes.NewBufferString("")
	out.StdOut = buffer
	err = cmd.Execute()
	if err != nil {
		t.Fatal("error executing command")
	}

	outStr, err := io.ReadAll(buffer)
	if err != nil {
		t.Fatal("error reading output")
	}

	got := strings.Split(string(outStr), "\n")

	assert.Contains(t, got[0], "<default>", "First line should contain workspace")
	assert.Contains(t, got[1], "test", "Second line should parent project")
	assert.Contains(t, got[2], "child", "Third line should child project")
}

func TestListTags(t *testing.T) {
	track, err := setupTestCommand()
	if err != nil {
		t.Fatal("error setting up test")
	}
	defer os.Remove(track.RootDir)

	project := core.NewProject("test", "", "t", []string{}, 15, 0)
	err = track.SaveProject(project, false)
	if err != nil {
		t.Fatal("error saving project")
	}

	record1 := core.Record{
		Project: "test",
		Start:   util.DateTime(2001, 2, 3, 4, 5, 0),
		End:     util.DateTime(2001, 2, 3, 5, 5, 0),
		Note:    "Test note with +key=value and +tag and +foo=bar",
	}
	record2 := core.Record{
		Project: "test",
		Start:   util.DateTime(2001, 2, 3, 6, 5, 0),
		End:     util.DateTime(2001, 2, 3, 7, 5, 0),
		Note:    "Test note with +tag and +foo=baz",
	}
	err = track.SaveRecord(&record1, false)
	if err != nil {
		t.Fatal("error saving record")
	}
	err = track.SaveRecord(&record2, false)
	if err != nil {
		t.Fatal("error saving record")
	}

	cmd := RootCommand(track, "")
	cmd.SetArgs([]string{"list", "tags"})

	buffer := bytes.NewBufferString("")
	out.StdOut = buffer
	err = cmd.Execute()
	if err != nil {
		t.Fatal("error executing command")
	}

	outStr, err := io.ReadAll(buffer)
	if err != nil {
		t.Fatal("error reading output")
	}

	got := strings.Split(string(outStr), "\n")

	assert.Contains(t, got[0], "foo", "Wrong tag")
	assert.Contains(t, got[0], "[bar baz]", "Wrong value")
	assert.Contains(t, got[0], "2", "Wrong count")

	assert.Contains(t, got[1], "key", "Wrong tag")
	assert.Contains(t, got[1], "[value]", "Wrong value")
	assert.Contains(t, got[1], "1", "Wrong count")

	assert.Contains(t, got[2], "tag", "Wrong tag")
	assert.Contains(t, got[2], "2", "Wrong count")
}

func TestListColors(t *testing.T) {
	track, err := setupTestCommand()
	if err != nil {
		t.Fatal("error setting up test")
	}
	defer os.Remove(track.RootDir)

	cmd := RootCommand(track, "")
	cmd.SetArgs([]string{"list", "colors"})

	buffer := bytes.NewBufferString("")
	out.StdOut = buffer
	err = cmd.Execute()
	if err != nil {
		t.Fatal("error executing command")
	}

	outStr, err := io.ReadAll(buffer)
	if err != nil {
		t.Fatal("error reading output")
	}

	got := strings.Split(string(outStr), "\n")

	assert.Equal(t, 18, len(got), "Wrong number of lines printed")
}

func TestListRecords(t *testing.T) {
	track, err := setupTestCommand()
	if err != nil {
		t.Fatal("error setting up test")
	}
	defer os.Remove(track.RootDir)

	project := core.NewProject("test", "", "t", []string{}, 15, 0)
	err = track.SaveProject(project, false)
	if err != nil {
		t.Fatal("error saving project")
	}

	record1 := core.Record{
		Project: "test",
		Start:   util.DateTime(2001, 2, 3, 4, 5, 0),
		End:     util.DateTime(2001, 2, 3, 5, 5, 0),
		Note:    "Test note with +key=value and +tag and +foo=bar",
	}
	record2 := core.Record{
		Project: "test",
		Start:   util.DateTime(2001, 2, 3, 6, 5, 0),
		End:     util.DateTime(2001, 2, 3, 7, 5, 0),
		Note:    "Test note with +tag and +foo=baz",
	}
	err = track.SaveRecord(&record1, false)
	if err != nil {
		t.Fatal("error saving record")
	}
	err = track.SaveRecord(&record2, false)
	if err != nil {
		t.Fatal("error saving record")
	}

	cmd := RootCommand(track, "")
	cmd.SetArgs([]string{"list", "records", "2001-02-03"})

	buffer := bytes.NewBufferString("")
	out.StdOut = buffer
	err = cmd.Execute()
	if err != nil {
		t.Fatal("error executing command")
	}

	outStr, err := io.ReadAll(buffer)
	if err != nil {
		t.Fatal("error reading output")
	}

	got := strings.Split(string(outStr), "\n")
	fmt.Println(string(outStr))

	assert.Contains(t, got[0], "test", "Wrong project name")
	assert.Contains(t, got[0], "2001-02-03 04:05 - 05:05", "Wrong time range")
	assert.Contains(t, got[0], "Test note with +key=value and +tag and +foo=bar", "Wrong note")

	assert.Contains(t, got[1], "test", "Wrong project name")
	assert.Contains(t, got[1], "2001-02-03 06:05 - 07:05", "Wrong time range")
	assert.Contains(t, got[1], "Test note with +tag and +foo=baz", "Wrong note")
}
