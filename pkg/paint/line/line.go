package line

import (
	"bytes"
	"text/template"

	"github.com/zoomoid/waveman/v2/pkg/paint"
)

type Interpolation string

const (
	// InterpolationNone means no interpolated points and a piecewise linear curve as the result
	InterpolationNone Interpolation = "none"
	// InterpolationFritschCarlson applies the Fritsch-Carlson method for
	// interpolating hermitic cubic splines to fit the data points
	// See: http://math.stackexchange.com/questions/45218/implementation-of-monotone-cubic-interpolation
	InterpolationFritschCarlson Interpolation = "fritsch-carlson"
	// InterpolationSteffen applies the Steffen method for interpolating
	// hermetic cubic splices to fit the data points
	// See: http://math.stackexchange.com/questions/45218/implementation-of-monotone-cubic-interpolation
	InterpolationSteffen Interpolation = "steffen"
)

const (
	DefaultPathTemplate = `<path d="{{.Points}}" fill="{{.Fill}}" stroke="{{.Stroke.Color}}" stroke-width="{{.Stroke.Width}}" />`
)

const (
	// Default interpolation mode
	DefaultInterpolation = InterpolationFritschCarlson
	// Default stroke width
	DefaultStrokeWidth = "5px"
	// Default stroke color
	DefaultStrokeColor = "black"
	// Default fill color
	DefaultFillColor = "rgba(0 0 0 / 0.5)"
	// Default horizontal spread of data points
	DefaultSpread = float64(10)
	// DefaultHeight of a canvas is 200px
	DefaultHeight = float64(200)
)

type Stroke struct {
	// Color is the CSS-compliant stroke color used for the path
	Color string
	// Width is the stroke width used for the path
	Width string
}

type LineOptions struct {
	// Interpolation choses the point interpolation mode
	Interpolation Interpolation
	// Fill is the CSS-compliant color value for the fill
	Fill string
	// Stroke is a struct defining properties of the stroke, namely color and
	// width
	Stroke *Stroke
	// Closed makes the bezier curve a closed one by appending the "Z" parameter
	// to the path
	Closed bool
	// Spread is the horizontal scaling factor by which all indices are scaled
	Spread float64
	// Height is the vertical scaling factor by which all sample points are scaled
	// up. Height also determines the total canvas height when normalized samples are used
	Height float64
	// Inverted transforms the SVG group to be horizontically flipped
	Inverted bool
}

// Compile-time type checking for LinePainter to implement all functions required
// by the Painter interface
var _ paint.Painter = &LinePainter{}

type LinePainter struct {
	// Embed all painter options, i.e., data points
	*paint.PainterOptions
	// Embed all options for the line painter
	*LineOptions
}

// New constructs a new Line painter with the passed options and fills in defaults
// for missing fields
func New(painter *paint.PainterOptions, options *LineOptions) *LinePainter {
	if options.Interpolation == "" {
		options.Interpolation = DefaultInterpolation
	}
	if options.Fill == "" {
		options.Fill = DefaultFillColor
	}
	if options.Stroke == nil {
		options.Stroke = &Stroke{
			Color: DefaultStrokeColor,
			Width: DefaultStrokeWidth,
		}
	}
	if options.Height == 0 {
		options.Height = DefaultHeight
	}
	if options.Spread == 0 {
		options.Spread = DefaultSpread
	}

	return &LinePainter{
		PainterOptions: painter,
		LineOptions:    options,
	}
}

// TotalHeight is the maximum height of any point in the curve
// since the interpolation does not overshoot, the maximum height, thus the
// total height, is the same as the Height scaling, as long as normalized data
// is passed to the painter.
func (l *LinePainter) TotalHeight() float64 {
	return l.Height
}

// TotalWidth is the width of the canvas, that is, the width the path spans horizontally.
func (l *LinePainter) TotalWidth() float64 {
	return float64(len(l.PainterOptions.Data)-1)*l.Spread + 2*l.Spread
}

// Draw implements line drawing with optional interpolation to smooth out the curve.
// Curves are, by default, anchored to the top-left corner. If you want to change
// this, consider transforming the entire SVG in-post using CSS transforms.
func (l *LinePainter) Draw() []string {
	output := make([]string, 3)

	pathTemplate := template.New("path")
	pathTemplate.Parse(DefaultPathTemplate)

	// make a slice of pairs that have the spread x values and their y values
	// paired
	samples := make([][2]float64, 0)
	samples = append(samples, [2]float64{0, 0})
	for i, sample := range l.Data {
		// offset samples in X direction by one unit of spread to account for start points
		samples = append(samples, [2]float64{float64(i)*l.Spread + l.Spread, sample * l.Height})
	}

	samples = append(samples, [2]float64{l.TotalWidth(), 0})

	line := ""
	switch l.Interpolation {
	case InterpolationSteffen:
		line = MonotonicCube(samples, steffen)
	case InterpolationNone:
		line = None(samples)
	case InterpolationFritschCarlson:
		line = MonotonicCube(samples, fritschCarlson)
	}

	if l.Closed {
		line += " Z\n"
	}

	bindings := struct {
		Fill   string
		Points string
		Stroke *Stroke
	}{
		Fill:   l.Fill,
		Points: line,
		Stroke: l.Stroke,
	}

	templateBuf := &bytes.Buffer{}
	pathTemplate.Execute(templateBuf, bindings)

	if l.Inverted {
		output = append(output, `<g style="transform: scaleY(-1)">`)
	} else {
		output = append(output, "<g>")
	}
	output = append(output, templateBuf.String())
	output = append(output, "</g>")
	return output
}
