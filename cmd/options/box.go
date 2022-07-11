package options

const (
	BoxFill         string = "color"
	Alignment       string = "alignment"
	BoxRounded      string = "rounded"
	BoxRoundedShort string = "r"
	BoxGap          string = "gap"
)

const (
	BoxFillDescription    string = "Fill color of each box"
	AlignmentDescription  string = "Alignment of the shapes, chose one of 'top', 'center', or 'bottom'"
	BoxRoundedDescription string = "Rounding factor of each box. Given in pixels. See SVG <rect> rx/ry attributes for details"
	BoxGapDescription     string = "Gap is the spacing left between each box. Boxes are centered horizonally, so half of gap is subtracted from the box's width"
)
