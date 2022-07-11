package options

const (
	Interpolation string = "interpolation"
	LineFill      string = "fill-color"
	StrokeColor   string = "stroke-color"
	StrokeWidth   string = "stroke-width"
	Closed        string = "closed"
	ClosedShort   string = "c"
	Inverted      string = "inverted"
	InvertedShort string = "i"
)

const (
	InterpolationDescription string = "Interpolation mechanism to be used for smoothing the curve. Choose one of 'none', 'fritsch-carlson', or 'steffen'"
	LineFillDescription      string = "Color for the area enclosed by the line"
	StrokeColorDescription   string = "Color of the line's stroke"
	StrokeWidthDescription   string = "Width of the line's stroke"
	ClosedDescription        string = "Whether the SVG path should be closed or left open"
	InvertedDescription      string = "Whether the shape should be inverted horizontally, i.e., switch the vertical alignment from top to bottom"
)
