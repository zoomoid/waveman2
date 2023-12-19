/*
Copyright 2022-2023 zoomoid.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
