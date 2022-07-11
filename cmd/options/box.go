package options

const (
	BoxFill         string = "fill-color"
	Alignment       string = "alignment"
	BoxHeight       string = "height"
	BoxHeightShort  string = "h"
	BoxWidth        string = "width"
	BoxWidthShort   string = "w"
	BoxRounded      string = "rounded"
	BoxRoundedShort string = "r"
	BoxGap          string = "gap"
)

const (
	BoxFillDescription    string = "Fill color of each box"
	AlignmentDescription  string = "Alignment of the shapes, chose one of 'top', 'center', or 'bottom'"
	BoxHeightDescription  string = "Height of the shape"
	BoxWidthDescription   string = "Width of each box"
	BoxRoundedDescription string = "Rounding factor of each box. Given in pixels. See SVG <rect> rx/ry attributes for details"
	BoxGapDescription     string = "Gap is the spacing left between each box. Boxes are centered horizonally, so half of gap is subtracted from the box's width"
)
