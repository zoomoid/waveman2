package cmd

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"github.com/zoomoid/waveman/v2/cmd/options"
	cmdutils "github.com/zoomoid/waveman/v2/cmd/utils"
	"github.com/zoomoid/waveman/v2/cmd/validation"
	"github.com/zoomoid/waveman/v2/pkg/painter/line"
)

var _ Plugin = &LinePainter{}

var LinePainterPlugin Plugin = &LinePainter{
	data:    newLineData(),
	enabled: false,
}

const (
	LineMode        string = "line"
	LineDescription string = "Create a line waveform"
)

type LinePainter struct {
	data    *lineData
	enabled bool
}

func (l *LinePainter) Name() string {
	return LineMode
}

func (l *LinePainter) Description() string {
	return LineDescription
}

func (l *LinePainter) Enabled() *bool {
	return &l.enabled
}

func (l *LinePainter) Data() interface{} {
	return l.data
}

func (l *LinePainter) Validate() error {
	errs := l.data.ValidateLineOptions()
	errlist := cmdutils.NewErrorList(errs)
	if errlist == nil {
		return nil
	}
	return errors.New(errlist.Error())
}

func (l *LinePainter) Flags(flags *pflag.FlagSet) {
	data, ok := l.Data().(lineData)
	if !ok {
		log.Fatal().Msg("line data struct is malformed")
	}
	flags.BoolVar(l.Enabled(), LineMode, false, LineDescription)
	flags.StringVar(&data.interpolation, options.Interpolation, string(line.DefaultInterpolation), options.InterpolationDescription)
	flags.StringVar(&data.fill, options.LineFill, line.DefaultFillColor, options.LineFillDescription)
	flags.StringVar(&data.strokeColor, options.StrokeColor, line.DefaultStrokeColor, options.StrokeColorDescription)
	flags.StringVar(&data.strokeWidth, options.StrokeWidth, line.DefaultStrokeWidth, options.StrokeWidthDescription)
	flags.Float64VarP(&data.spread, options.LineSpread, options.LineSpreadShort, line.DefaultSpread, options.LineSpreadDescription)
	flags.Float64VarP(&data.height, options.LineHeight, options.LineHeightShort, line.DefaultHeight, options.LineHeightDescription)
	flags.BoolVarP(&data.closed, options.Closed, options.ClosedShort, false, options.ClosedDescription)
	flags.BoolVarP(&data.inverted, options.Inverted, options.InvertedShort, false, options.InvertedDescription)
}

func (l *lineData) ValidateLineOptions() (errList []error) {
	if err := validation.ValidateInterpolation(l.interpolation); err != nil {
		errList = append(errList, err)
	}
	if err := validation.ValidateSpread(l.spread); err != nil {
		errList = append(errList, err)
	}
	if err := validation.ValidateLineHeight(l.height); err != nil {
		errList = append(errList, err)
	}
	if len(errList) == 0 {
		return nil
	}
	return errList
}

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
		interpolation: string(line.DefaultInterpolation),
		fill:          line.DefaultFillColor,
		strokeColor:   line.DefaultStrokeColor,
		strokeWidth:   line.DefaultStrokeWidth,
		spread:        line.DefaultSpread,
		height:        line.DefaultHeight,
		closed:        false,
		inverted:      false,
	}
}
