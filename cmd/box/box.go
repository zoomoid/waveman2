package box

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/zoomoid/waveman/v2/cmd/options"
	b "github.com/zoomoid/waveman/v2/pkg/paint/box"
)

type boxData struct {
	color     string
	alignment string
	height    float64
	width     float64
	rounded   float64
	gap       float64
}

func newBoxData() *boxData {
	return &boxData{
		color:     b.DefaultColor,
		alignment: string(b.DefaultAlignment),
		height:    b.DefaultHeight,
		width:     b.DefaultWidth,
		rounded:   b.DefaultRounded,
		gap:       b.DefaultRounded,
	}
}

func NewCommand(data *boxData) *cobra.Command {
	if data == nil {
		data = newBoxData()
	}

	cmd := &cobra.Command{}

	addLineOptions(cmd.Flags(), data)

	return cmd
}

func addLineOptions(flags *pflag.FlagSet, data *boxData) {
	flags.StringVar(&data.color, options.BoxFill, b.DefaultColor, options.BoxFillDescription)
	flags.StringVar(&data.alignment, options.Alignment, string(b.DefaultAlignment), options.AlignmentDescription)
	flags.Float64VarP(&data.height, options.BoxHeight, options.BoxHeightShort, b.DefaultHeight, options.BoxHeightDescription)
	flags.Float64VarP(&data.width, options.BoxWidth, options.BoxWidthShort, b.DefaultWidth, options.BoxWidthDescription)
	flags.Float64VarP(&data.rounded, options.BoxRounded, options.BoxRoundedShort, b.DefaultRounded, options.BoxRoundedDescription)
	flags.Float64Var(&data.gap, options.BoxGap, b.DefaultGap, options.BoxGapDescription)
}
