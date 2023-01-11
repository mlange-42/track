package treemap

import (
	"bytes"
	"image/color"
	"io"

	"github.com/mlange-42/track/core"
	"github.com/nikolaydubina/treemap"
	"github.com/nikolaydubina/treemap/parser"
	"github.com/nikolaydubina/treemap/render"
)

// SvgOptions are options for SVG treemap rendering
type SvgOptions struct {
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

// SvgRenderer renders a tree in treemap SVG format
type SvgRenderer struct {
	Track    *core.Track
	Reporter core.Reporter
	Options  SvgOptions
}

// Render renders a FileTree
func (r SvgRenderer) Render(w io.Writer) error {
	csvRenderer := CsvRenderer{
		Track:    r.Track,
		Reporter: r.Reporter,
	}

	buffer := bytes.Buffer{}
	err := csvRenderer.Render(&buffer)
	if err != nil {
		return err
	}

	svgBytes, err := toSvg(buffer.String(), &r.Options)
	if err != nil {
		return err
	}
	w.Write(svgBytes)
	return nil
}

var grey = color.RGBA{128, 128, 128, 255}

func toSvg(s string, flags *SvgOptions) ([]byte, error) {
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
