package cli

import (
	"fmt"
	"image/color"
	"os"
	"strings"

	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/mlange-42/track/util"
	"github.com/nikolaydubina/treemap"
	"github.com/nikolaydubina/treemap/parser"
	"github.com/nikolaydubina/treemap/render"
	"github.com/spf13/cobra"
)

type svgFlags struct {
	W             float64
	H             float64
	MarginBox     float64
	PaddingBox    float64
	Padding       float64
	ColorScheme   string
	ColorBorder   string
	ImputeHeat    bool
	KeepLongPaths bool
}

func treemapReportCommand(t *core.Track, options *filterOptions) *cobra.Command {
	var svg svgFlags
	var csv bool

	treemap := &cobra.Command{
		Use:     "treemap",
		Short:   "Generates a treemap of time tracking in SVG format",
		Aliases: []string{"m"},
		Args:    util.WrappedArgs(cobra.NoArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			projects, err := t.LoadAllProjects()
			if err != nil {
				return fmt.Errorf("failed to generate report: %s", err.Error())
			}

			filters, err := createFilters(options, projects, false)
			if err != nil {
				return fmt.Errorf("failed to generate report: %s", err.Error())
			}

			startTime, endTime, err := parseStartEnd(options)
			if err != nil {
				return fmt.Errorf("failed to generate report: %s", err.Error())
			}
			reporter, err := core.NewReporter(
				t, options.projects, filters,
				options.includeArchived, startTime, endTime,
			)
			if err != nil {
				return fmt.Errorf("failed to generate report: %s", err.Error())
			}

			tree, err := t.ToProjectTree(reporter.Projects)
			if err != nil {
				return fmt.Errorf("failed to generate report: %s", err.Error())
			}

			formatter := TreemapPrinter{*reporter}
			str := formatter.Print(tree)
			if csv {
				out.Print(str)
				return nil
			}

			svgBytes, err := toSvg(str, &svg)
			if err != nil {
				out.Err("failed to generate report: %s", err.Error())
				return nil
			}
			os.Stdout.Write(svgBytes)
			return nil
		},
	}
	treemap.Flags().StringVarP(&options.start, "start", "s", "", "Start date (start at 00:00)")
	treemap.Flags().StringVarP(&options.end, "end", "e", "", "End date (inclusive: end at 24:00)")

	treemap.Flags().BoolVar(&csv, "csv", false, "Generate raw CSV output for github.com/nikolaydubina/treemap")

	treemap.Flags().Float64Var(&svg.W, "w", 1028, "width of output")
	treemap.Flags().Float64Var(&svg.H, "h", 640, "height of output")
	treemap.Flags().Float64Var(&svg.MarginBox, "margin-box", 4, "margin between boxes")
	treemap.Flags().Float64Var(&svg.PaddingBox, "padding-box", 4, "padding between box border and content")
	treemap.Flags().Float64Var(&svg.Padding, "padding", 32, "padding around root content")
	treemap.Flags().StringVar(&svg.ColorScheme, "color", "balance", "color scheme (RdBu, balance, RdYlGn, none)")
	treemap.Flags().StringVar(&svg.ColorBorder, "color-border", "auto", "color of borders (light, dark, auto)")
	treemap.Flags().BoolVar(&svg.ImputeHeat, "impute-heat", false, "impute heat for parents(weighted sum) and leafs(0.5)")
	treemap.Flags().BoolVar(&svg.KeepLongPaths, "long-paths", false, "keep long paths when paren has single child")

	return treemap
}

// TreemapPrinter prints a tree in treemap CSV format
type TreemapPrinter struct {
	Reporter core.Reporter
}

// Print prints a FileTree
func (p TreemapPrinter) Print(t *core.ProjectTree) string {
	sb := strings.Builder{}
	p.print(t.Root, &sb, "")
	return sb.String()
}

func (p TreemapPrinter) print(t *core.ProjectNode, sb *strings.Builder, path string) {
	total := p.Reporter.TotalTime[t.Value.Name]
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
		sb,
		"%s,%f,%f\n",
		strings.Replace(path, ",", "-", -1),
		totalHours,
		0.0,
	)
	for _, child := range t.Children {
		p.print(child, sb, path)
	}
}

var grey = color.RGBA{128, 128, 128, 255}

func toSvg(s string, flags *svgFlags) ([]byte, error) {
	parser := parser.CSVTreeParser{}
	tree, err := parser.ParseString(s)
	if err != nil || tree == nil {
		return []byte{}, err
	}

	treemap.SetNamesFromPaths(tree)
	if !flags.KeepLongPaths {
		treemap.CollapseLongPaths(tree)
	}

	sizeImputer := treemap.SumSizeImputer{EmptyLeafSize: 1}
	sizeImputer.ImputeSize(*tree)

	if flags.ImputeHeat {
		heatImputer := treemap.WeightedHeatImputer{EmptyLeafHeat: 0.5}
		heatImputer.ImputeHeat(*tree)
	}

	tree.NormalizeHeat()

	var colorer render.Colorer

	palette, hasPalette := render.GetPalette(flags.ColorScheme)
	treeHueColorer := render.TreeHueColorer{
		Offset: 0,
		Hues:   map[string]float64{},
		C:      0.5,
		L:      0.5,
		DeltaH: 10,
		DeltaC: 0.3,
		DeltaL: 0.1,
	}

	var borderColor color.Color
	borderColor = color.White
	colorer = treeHueColorer
	borderColor = color.White

	colorer = treeHueColorer
	borderColor = color.White

	switch {
	case flags.ColorScheme == "none":
		colorer = render.NoneColorer{}
		borderColor = grey
	case flags.ColorScheme == "balanced":
		colorer = treeHueColorer
		borderColor = color.White
	case hasPalette && tree.HasHeat():
		colorer = render.HeatColorer{Palette: palette}
	case tree.HasHeat():
		palette, _ := render.GetPalette("RdBu")
		colorer = render.HeatColorer{Palette: palette}
	default:
		colorer = treeHueColorer
	}

	switch {
	case flags.ColorBorder == "light":
		borderColor = color.White
	case flags.ColorBorder == "dark":
		borderColor = grey
	}

	uiBuilder := render.UITreeMapBuilder{
		Colorer:     colorer,
		BorderColor: borderColor,
	}
	spec := uiBuilder.NewUITreeMap(*tree, flags.W, flags.H, flags.MarginBox, flags.PaddingBox, flags.Padding)
	renderer := render.SVGRenderer{}

	return renderer.Render(spec, flags.W, flags.H), nil
}
