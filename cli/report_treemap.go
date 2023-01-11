package cli

import (
	"fmt"

	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/mlange-42/track/render"
	"github.com/mlange-42/track/render/treemap"
	"github.com/mlange-42/track/util"
	"github.com/spf13/cobra"
)

func treemapReportCommand(t *core.Track, options *filterOptions) *cobra.Command {
	var svg treemap.SvgOptions
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

			var renderer render.Renderer
			if csv {
				renderer = treemap.CsvRenderer{
					Track:    t,
					Reporter: *reporter,
				}
			} else {
				renderer = treemap.SvgRenderer{
					Track:    t,
					Reporter: *reporter,
					Options:  svg,
				}
			}
			renderer.Render(out.StdOut)
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
