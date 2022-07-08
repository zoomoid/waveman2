package line

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/zoomoid/waveman/v2/cmd/options"
	l "github.com/zoomoid/waveman/v2/pkg/paint/line"
)

type lineData struct {
	interpolation string
	fill          string
	strokeColor   string
	strokeWidth   string
	spread        float64
	height        float64
	closed        bool
	inverted      bool
}

func newLineData() *lineData {
	return &lineData{
		interpolation: string(l.DefaultInterpolation),
		fill:          l.DefaultFillColor,
		strokeColor:   l.DefaultStrokeColor,
		strokeWidth:   l.DefaultStrokeWidth,
		spread:        l.DefaultSpread,
		height:        l.DefaultHeight,
		closed:        false,
		inverted:      false,
	}
}

func NewCommand(data *lineData) *cobra.Command {
	if data == nil {
		data = newLineData()
	}

	cmd := &cobra.Command{
		Use: "line",
		// TODO: add short and long description
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	addLineOptions(cmd.Flags(), data)

	return cmd
}

func addLineOptions(flags *pflag.FlagSet, data *lineData) {
	flags.StringVar(&data.interpolation, options.Interpolation, string(l.DefaultInterpolation), options.InterpolationDescription)
	flags.StringVar(&data.fill, options.LineFill, l.DefaultFillColor, options.LineFillDescription)
	flags.StringVar(&data.strokeColor, options.StrokeColor, l.DefaultStrokeColor, options.StrokeColorDescription)
	flags.StringVar(&data.strokeWidth, options.StrokeWidth, l.DefaultStrokeWidth, options.StrokeWidthDescription)
	flags.Float64VarP(&data.spread, options.LineSpread, options.LineSpreadShort, l.DefaultSpread, options.LineSpreadDescription)
	flags.Float64VarP(&data.height, options.LineHeight, options.LineHeightShort, l.DefaultHeight, options.LineHeightDescription)
	flags.BoolVarP(&data.closed, options.Closed, options.ClosedShort, false, options.ClosedDescription)
	flags.BoolVarP(&data.inverted, options.Inverted, options.InvertedShort, false, options.InvertedDescription)
}
