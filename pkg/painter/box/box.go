package box

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/zoomoid/waveman2/pkg/painter"
)

// Alignment is the categorical type for determining the box alignment axis
type Alignment string

const (
	// AlignmentTop pulls all boxes to the canvas's upper boundary
	AlignmentTop Alignment = "top"
	// AlignmentCenter centers all boxes in the canvas's horizontal center axis
	AlignmentCenter Alignment = "center"
	// AlignmentBottom pulls all boxes to the canvas's lower boundary
	AlignmentBottom Alignment = "bottom"
	// AlignmentEmpty is used for catching unitialized alignment
	AlignmentEmpty Alignment = ""
)

const (
	DefaultRectangleTemplate = `<rect width="{{.Width}}" height="{{.Height}}" x="{{.X}}" y="{{.Y}}" rx="{{.Rounded}}" ry="{{.Rounded}}" fill="{{.Color}}" />`
)

const (
	// DefaultColor for boxes is black
	DefaultColor = "black"
	// DefaultAlignment for boxes is center. If an alignment type outside of the
	// supported ones is used, this default alignment is automatically chosen as a
	// fallback
	DefaultAlignment = AlignmentCenter
	// DefaultHeight of a canvas is 200px
	DefaultHeight = float64(200)
	// DefaultWidth of each box is 20px
	DefaultWidth = float64(20)
	// DefaultGap of each box is 10px
	DefaultGap = float64(5)
	// DefaultRounded rounding ratio is 10px
	DefaultRounded = float64(10)
)

// Compile-time type checking for BoxPainter to implement all functions required
// by the Painter interface
var _ painter.Painter = &BoxPainter{}

type BoxOptions struct {
	// Color for each rectangle, in a CSS-compliant format
	Color string
	// Alignment of the boxes, either top, center, or bottom
	Alignment Alignment
	// BoxHeight is the factor by which each sample value gets scaled upwards. Since
	// samples are expected to be normalized, this means that height is also the
	// maximum, thus total height of the graphic
	BoxHeight float64
	// BoxWidth is the absolute width of each bounding box of a box, including the
	// gap. Internally, the width is reduced by the gap
	BoxWidth float64
	// Rounded is the rounding value for all boxes. This is by default
	// symmetrical, there is currently no way to set this for x and y
	// indepdendently
	Rounded float64
	// Gap is the spacing between each box. Boxes are placed horizontally centered
	// in the bounding box of the height and the width, with their inner width
	// being reduced by the gap.
	Gap float64
	// totalWidth is the canvas's width that results from adding up each box's
	// width
	totalWidth float64
	// totalHeight is the canvas's total height (which, as long as samples are
	// normalized beforehand, will always be equal to BoxOptions.Height)
	totalHeight float64
}

// BoxPainter is the struct containing context for drawing a waveform as SVG
// rectangles
type BoxPainter struct {
	// Embed all painter options, i.e., data points
	*painter.PainterOptions
	// Embed all options for the box drawer
	*BoxOptions
}

// Height returns the canvas's total height. When normalized samples are
// used, this is equal to the height scaling factor
func (o *BoxPainter) Height() float64 {
	return o.totalHeight
}

// Width returns the canvas's total width. This is equal to the number of
// samples times the width of each box.
func (o *BoxPainter) Width() float64 {
	return o.totalWidth
}

// New constructs a new Box painter with the passed options and fills in
// defaults for missing fields
func New(painter *painter.PainterOptions, options *BoxOptions) *BoxPainter {
	if options.Color == "" {
		options.Color = DefaultColor
	}
	if options.Alignment == AlignmentEmpty {
		options.Alignment = AlignmentCenter
	}
	if options.BoxHeight == 0 {
		options.BoxHeight = DefaultHeight
	}
	if options.BoxWidth == 0 {
		options.BoxWidth = DefaultWidth
	}

	options.totalHeight = options.BoxHeight
	options.totalWidth = options.BoxWidth * float64(len(painter.Data))
	return &BoxPainter{
		PainterOptions: painter,
		BoxOptions:     options,
	}
}

// Draw implements the Painter interface's required Draw() function. For each
// sample, an SVG rectangle is created, and all of them are wrapped inside an
// SVG group element.
func (o *BoxPainter) Draw() []string {
	output := make([]string, 0)

	rectTemplate := template.New("rect")
	rectTemplate.Parse(DefaultRectangleTemplate)

	output = append(output, "<g>")
	for index, sample := range o.Data {
		buf := &bytes.Buffer{}
		if sample*o.BoxHeight < o.BoxWidth {
			sample = (o.BoxWidth - o.Gap) / o.BoxHeight
		}
		rect := o.perSample(index, sample)
		rectTemplate.Execute(buf, rect)
		output = append(output, buf.String())
	}
	output = append(output, "</g>")
	return output
}

// perSample is the handler that creates a Rectangle struct for each sample and
// its index.
func (o *BoxPainter) perSample(index int, sample float64) *Rectangle {
	rect := &Rectangle{}
	switch o.Alignment {
	case AlignmentBottom:
		rect = o.alignBottom(index, sample)
	case AlignmentTop:
		rect = o.alignTop(index, sample)
	case AlignmentCenter:
		rect = o.alignCenter(index, sample)
	}
	rect.Color = o.Color
	rect.Rounded = o.Rounded

	return rect
}

// alignTop implements the AlignmentTop alignment mode. It returns a minimal
// rectangle.
func (o *BoxPainter) alignTop(index int, sample float64) *Rectangle {
	pos := Position{
		x: float64(index)*o.BoxWidth + (0.5 * o.Gap),
		y: o.BoxHeight,
	}
	size := Dimensions{
		width:  o.BoxWidth - o.Gap,
		height: sample * o.BoxHeight,
	}
	return &Rectangle{
		Position:   pos,
		Dimensions: size,
	}
}

// alignCenter implement the AlignmentCenter alignment mode. It returns a
// minimal rectangle.
func (o *BoxPainter) alignCenter(index int, sample float64) *Rectangle {
	pos := Position{
		x: float64(index)*o.BoxWidth + (0.5 * o.Gap),
		y: (0.5 * o.BoxHeight) - (0.5 * sample * o.BoxHeight),
	}
	size := Dimensions{
		width:  o.BoxWidth - o.Gap,
		height: sample * o.BoxHeight,
	}
	return &Rectangle{
		Position:   pos,
		Dimensions: size,
	}
}

// alignBottom implement the AlignmentBottom alignment mode. It returns a
// minimal rectangle.
func (o *BoxPainter) alignBottom(index int, sample float64) *Rectangle {
	pos := Position{
		x: float64(index)*o.BoxWidth + (0.5 * o.Gap),
		y: (1 - sample) * o.BoxHeight,
	}
	size := Dimensions{
		width:  o.BoxWidth - o.Gap,
		height: sample * o.BoxHeight,
	}
	return &Rectangle{
		Position:   pos,
		Dimensions: size,
	}
}

func (o *BoxPainter) Viewbox() string {
	return fmt.Sprintf("0 0 %f %f", o.totalWidth, o.totalHeight)
}
