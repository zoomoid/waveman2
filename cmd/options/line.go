package options

const (
	Line            string = "line"
	Interpolation   string = "interpolation"
	LineFill        string = "fill-color"
	StrokeColor     string = "stroke-color"
	StrokeWidth     string = "stroke-width"
	LineSpread      string = "spread"
	LineSpreadShort string = "s"
	LineHeight      string = "height"
	LineHeightShort string = "h"
	Closed          string = "closed"
	ClosedShort     string = "c"
	Inverted        string = "inverted"
	InvertedShort   string = "i"
)

const (
	LineDescription          string = "Create a line waveform"
	InterpolationDescription string = "Interpolation mechanism to be used for smoothing the curve. Choose one of 'none', 'fritsch-carlson', or 'steffen'"
	LineFillDescription      string = "Color for the area enclosed by the line"
	StrokeColorDescription   string = "Color of the line's stroke"
	StrokeWidthDescription   string = "Width of the line's stroke"
	LineSpreadDescription    string = "Per-point distance by which each sample point gets distributed horizontally"
	LineHeightDescription    string = "Height of the line shape"
	ClosedDescription        string = "Whether the SVG path should be closed or left open"
	InvertedDescription      string = "Whether the shape should be inverted horizontally, i.e., switch the vertical alignment from top to bottom"
)
