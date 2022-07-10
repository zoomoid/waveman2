package cmd

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/zoomoid/waveman/v2/cmd/options"
	"github.com/zoomoid/waveman/v2/cmd/validation"
	"github.com/zoomoid/waveman/v2/pkg/paint/box"
	"github.com/zoomoid/waveman/v2/pkg/paint/line"
	"github.com/zoomoid/waveman/v2/pkg/transform"
)

func Execute() {
	rootCmd, err := NewWavemanCommand(nil)
	if err != nil {
		log.Fatal().Err(err)
	}
	if err := rootCmd.command.Execute(); err != nil {
		log.Fatal().Err(err)
	}
}

// newContext constructs a data struct to be used in the closure of the NewWavemanCommand
// constructor function when no struct is given as a parameter. This is the default case,
// the ability to pass data as a parameter is present to make unit testing commands
// possible
func newContext() *context {
	return &context{
		transfomerData: newTransformerData(),
		boxData:        newBoxData(),
		lineData:       newLineData(),
	}
}

type context struct {
	*transfomerData
	*boxData
	*lineData
}

type WavemanCommand struct {
	command           *cobra.Command
	context           *context
	availablePainters []string
	modes             map[string]*bool
}

var referencePainters = []string{options.Box, options.Line}

// NewWaveman creates a new cobra command and adds the relevant flags to the root command.
// It also creates the link to the subcommands
func NewWavemanCommand(data *context) (*WavemanCommand, error) {
	if data == nil {
		data = newContext()
	}

	cmd := &cobra.Command{
		Use: "waveman",
		// TODO: add short and long description
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	// add transformer flags
	addTransformerOptions(cmd.Flags(), data.transfomerData)

	// add reference painter implementation flags

	modes := make(map[string]*bool)

	cmd.Flags().BoolVar(modes[options.Box], options.Box, false, options.BoxDescription)
	addBoxOptions(cmd.Flags(), data.boxData)
	cmd.Flags().BoolVar(modes[options.Line], options.Line, false, options.LineDescription)
	addLineOptions(cmd.Flags(), data.lineData)

	if err := validation.ValidatePainterModes(cmd.Flags(), referencePainters); err != nil {
		return nil, err
	}

	return &WavemanCommand{
		command:           cmd,
		context:           data,
		availablePainters: referencePainters,
	}, nil
}

// FlagsFactory is a function type that, upon calling from the Plugin function, should register all flags the
// Plugin painter requires to the FlagSet of the command. Data is arbitrary data, that maps to the flag values.
type FlagsFactory func(flags *pflag.FlagSet, data interface{})

// Plugin allows a user to patch in additional painters and register their flags to the waveman command.
func (w *WavemanCommand) Plugin(painterName string, flagValue *bool, flagDescription string, flags FlagsFactory, data *interface{}) (*WavemanCommand, error) {
	if _, ok := w.modes[painterName]; ok {
		return w, fmt.Errorf("painter %s is already registered in the chain", painterName)
	}
	w.modes[painterName] = flagValue
	w.availablePainters = append(w.availablePainters, painterName)
	w.command.Flags().BoolVar(flagValue, painterName, false, flagDescription)
	flags(w.command.Flags(), data)
	return w, nil
}

func addTransformerOptions(flags *pflag.FlagSet, data *transfomerData) {
	flags.StringVar(&data.downsamplingMode, options.DownsamplingMode, "", options.DownsamplingModeDescription)
	flags.IntVar(&data.downsamplingFactor, options.DownsamplingFactor, 1, options.DownsamplingFactorDescription)
	flags.StringVar(&data.aggregator, options.Aggregator, string(transform.DefaultAggregator), options.AggregatorDescription)
	flags.StringVarP(&data.filename, options.Filename, options.FilenameShort, "", options.FilenameDescription)
	flags.IntVarP(&data.chunks, options.Chunks, options.ChunksShort, transform.DefaultChunks, options.ChunksDescription)
	flags.StringVarP(&data.output, options.Output, options.OutputShort, "", options.OutputDescription)
}

func addBoxOptions(flags *pflag.FlagSet, data *boxData) {
	flags.StringVar(&data.color, options.BoxFill, box.DefaultColor, options.BoxFillDescription)
	flags.StringVar(&data.alignment, options.Alignment, string(box.DefaultAlignment), options.AlignmentDescription)
	flags.Float64VarP(&data.height, options.BoxHeight, options.BoxHeightShort, box.DefaultHeight, options.BoxHeightDescription)
	flags.Float64VarP(&data.width, options.BoxWidth, options.BoxWidthShort, box.DefaultWidth, options.BoxWidthDescription)
	flags.Float64VarP(&data.rounded, options.BoxRounded, options.BoxRoundedShort, box.DefaultRounded, options.BoxRoundedDescription)
	flags.Float64Var(&data.gap, options.BoxGap, box.DefaultGap, options.BoxGapDescription)
}

func addLineOptions(flags *pflag.FlagSet, data *lineData) {
	flags.StringVar(&data.interpolation, options.Interpolation, string(line.DefaultInterpolation), options.InterpolationDescription)
	flags.StringVar(&data.fill, options.LineFill, line.DefaultFillColor, options.LineFillDescription)
	flags.StringVar(&data.strokeColor, options.StrokeColor, line.DefaultStrokeColor, options.StrokeColorDescription)
	flags.StringVar(&data.strokeWidth, options.StrokeWidth, line.DefaultStrokeWidth, options.StrokeWidthDescription)
	flags.Float64VarP(&data.spread, options.LineSpread, options.LineSpreadShort, line.DefaultSpread, options.LineSpreadDescription)
	flags.Float64VarP(&data.height, options.LineHeight, options.LineHeightShort, line.DefaultHeight, options.LineHeightDescription)
	flags.BoolVarP(&data.closed, options.Closed, options.ClosedShort, false, options.ClosedDescription)
	flags.BoolVarP(&data.inverted, options.Inverted, options.InvertedShort, false, options.InvertedDescription)
}
