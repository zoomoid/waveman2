package reference

import (
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/zoomoid/waveman2/cmd/options"
	cmdutils "github.com/zoomoid/waveman2/cmd/utils"
	"github.com/zoomoid/waveman2/cmd/validation"
	"github.com/zoomoid/waveman2/pkg/painter"
	"github.com/zoomoid/waveman2/pkg/painter/line"
	"github.com/zoomoid/waveman2/pkg/plugin"
)

var _ plugin.Plugin = &LinePainter{}

var LinePainterPlugin plugin.Plugin = &LinePainter{
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
	painter *line.LinePainter
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
	errs := l.data.validateLineOptions()
	errlist := cmdutils.NewErrorList(errs)
	if errlist == nil {
		return nil
	}
	return errors.New(errlist.Error())
}

func (l *LinePainter) Flags(flags *pflag.FlagSet) error {
	data, ok := l.Data().(*lineData)
	if !ok {
		return errors.New("line data struct is malformed")
	}
	flags.BoolVar(l.Enabled(), LineMode, false, LineDescription)
	flags.StringVar(&data.interpolation, options.Interpolation, string(line.DefaultInterpolation), options.InterpolationDescription)
	flags.StringVar(&data.fill, options.LineFill, line.DefaultFillColor, options.LineFillDescription)
	flags.StringVar(&data.strokeColor, options.StrokeColor, line.DefaultStrokeColor, options.StrokeColorDescription)
	flags.StringVar(&data.strokeWidth, options.StrokeWidth, line.DefaultStrokeWidth, options.StrokeWidthDescription)
	flags.BoolVarP(&data.closed, options.Closed, options.ClosedShort, false, options.ClosedDescription)
	flags.BoolVarP(&data.inverted, options.Inverted, options.InvertedShort, false, options.InvertedDescription)
	return nil
}

func (l *LinePainter) Draw(options *painter.PainterOptions) []string {
	painter := line.New(options, l.data.toOptions(options.Width, options.Height))
	l.painter = painter
	return painter.Draw()
}

func (l *LinePainter) Painter() painter.Painter {
	return l.painter
}

func (l *lineData) toOptions(width float64, height float64) *line.LineOptions {
	return &line.LineOptions{
		Interpolation: line.Interpolation(l.interpolation),
		Fill:          l.fill,
		Stroke: &line.Stroke{
			Color: l.strokeColor,
			Width: l.strokeWidth,
		},
		Closed:    l.closed,
		Spread:    width,
		Amplitude: height,
		Inverted:  l.inverted,
	}
}

func (l *lineData) validateLineOptions() (errList []error) {
	if err := validation.ValidateInterpolation(l.interpolation); err != nil {
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
