package box

import "github.com/lithammer/dedent"

var (
	group string = "core/v1"

	description string = dedent.Dedent(`
		The box painter draws a simple box for each data point.

		The box color can be set with --color.
		
		The alignment axis can be either "top", "center", or "bottom",
		and set with --alignment. 
		
		--height (or -h) sets the height of highest box, thus also the
		height of the entire canvas. 
		
		--width (or -w) sets the width of each box's bounding box.
		
		--gap sets the space left between each box. Boxes are painted centered inside 
		their bounding box:

		|-------------------------------------------|
		|<- 0.5 * gap ->|----BOX----|<- 0.5 * gap ->|
		|<----------------- width ----------------->|

		--rounded (or -r) parameter controls the rounding of the rectangles.
		Notably, rounding requires the boxes to have a minimum height, namely at least
		the width of the box, to look aesthetically pleasing. When using --rounded,
		each box's height will have its width as a lower bound.
	`)
)
