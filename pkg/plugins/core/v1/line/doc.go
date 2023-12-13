package line

import "github.com/lithammer/dedent"

var (
	group string = "core/v1"

	description string = dedent.Dedent(`
		The box painter draws a line that traces each data point.

		A line's path can be closed by setting the --closed (or -c) flag.
		This will close the <path> by appending "Z" at the end of the data points.

		When the path is closed, the color of the enclosed shape can be set with 
		--fill-color. 
		
		The color of the line is set with --stroke-color, and the width of the line 
		with --stroke-width. All those require SVG/CSS-compliant values for the
		attributes.

		The shape can be horizontally mirrored by setting --inverted (-i).

		To create a symmetric shape, similar to Box with alignment = center, but with a 
		continuous line, use --mirrored.

		Similarly to the box painter, the --height (or -h) flag controls the shape's overall
		height. 
		
		--spread (or -s) controls the horizontal spacing between each of the sample
		points. 

		To make the line appear smoothly from a discrete set of points, we interpolate 
		control points for each sample point using cubic hermetic interpolation to fit 
		cubic polynomials. Namely, we implement 2 interpolation schemes: Fritsch-Carlson
		and Steffen. Details can be seen here: 
		http://math.stackexchange.com/questions/45218/implementation-of-monotone-cubic-interpolation
		This way, the shape appears smooth. Interpolation can also be controlled with
		flags: By default, the Frisch-Carlson scheme is used; setting "--interpolation steffen"
		uses the Steffen scheme. If you want to disable interpolation entirely, set
		"--interpolation none".
	`)
)
