/*
Copyright 2022 zoomoid.

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

package options

const (
	BoxFill    string = "color"
	Alignment  string = "alignment"
	BoxRounded string = "rounded"
	BoxGap     string = "gap"
)

const (
	BoxFillDescription    string = "Fill color of each box"
	AlignmentDescription  string = "Alignment of the shapes, chose one of 'top', 'center', or 'bottom'"
	BoxRoundedDescription string = "Rounding factor of each box. Given in pixels. See SVG <rect> rx/ry attributes for details"
	BoxGapDescription     string = "Gap is the spacing left between each box. Boxes are centered horizonally, so half of gap is subtracted from the box's width"
)
