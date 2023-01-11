package treemap

import (
	"fmt"
	"io"
	"strings"

	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/util"
)

// CsvRenderer renders a tree in treemap CSV format
type CsvRenderer struct {
	Track    *core.Track
	Reporter core.Reporter
}

// Render renders a FileTree
func (r CsvRenderer) Render(w io.Writer) error {
	tree, err := r.Track.ToProjectTree(r.Reporter.Projects)
	if err != nil {
		return err
	}

	return r.render(tree.Root, w, "")
}

func (r CsvRenderer) render(t *core.ProjectNode, w io.Writer, path string) error {
	total := r.Reporter.TotalTime[t.Value.Name]
	if len(path) == 0 {
		path = fmt.Sprintf("%s (%s)", t.Value.Name, util.FormatDuration(total))
	} else {
		path = fmt.Sprintf("%s/%s (%s)", path, t.Value.Name, util.FormatDuration(total))
	}

	totalHours := total.Hours()
	if totalHours == 0 {
		totalHours = 0.001
	}

	fmt.Fprintf(
		w,
		"%s,%f,%f\n",
		strings.Replace(path, ",", "-", -1),
		totalHours,
		0.0,
	)
	for _, child := range t.Children {
		r.render(child, w, path)
	}
	return nil
}
