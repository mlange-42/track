package cli

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/mlange-42/track/util"
	"github.com/stretchr/testify/assert"
)

func TestExport(t *testing.T) {
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

	record := core.Record{
		Project: "test",
		Start:   util.DateTime(2001, 2, 3, 4, 5, 0),
		End:     util.DateTime(2001, 2, 3, 5, 5, 0),
		Note:    "Test note with +tag=1",
		Tags:    map[string]string{"tag": "1"},
		Pause: []core.Pause{
			{
				Start: util.DateTime(2001, 2, 3, 4, 10, 0),
				End:   util.DateTime(2001, 2, 3, 4, 15, 0),
				Note:  "Pause",
			},
		},
	}
	err = track.SaveRecord(&record, false)
	if err != nil {
		t.Fatal("error saving record")
	}

	cmd := RootCommand(track, "")
	cmd.SetArgs([]string{"export", "records"})

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
	expected := `start,end,project,total,work,pause,note,tags
2001-02-03 04:05,2001-02-03 05:05,test,01:00,00:55,00:05,"Test note with +tag=1",tag=1
`
	assert.Equal(t, expected, got, "unexpected CSV output")

	cmd = RootCommand(track, "")
	cmd.SetArgs([]string{"export", "records", "--json"})

	err = cmd.Execute()
	if err != nil {
		t.Fatal("error executing command")
	}

	cmd = RootCommand(track, "")
	cmd.SetArgs([]string{"export", "records", "--yaml"})

	err = cmd.Execute()
	if err != nil {
		t.Fatal("error executing command")
	}
}
