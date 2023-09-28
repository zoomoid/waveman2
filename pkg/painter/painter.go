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

package painter

type PainterOptions struct {
	// Data contains all sample points to use in a drawing context
	Data   []float64
	Height float64
	Width  float64
}

type Painter interface {
	// Height is the interface function for getting the painter canvas's total height
	Height() float64
	// TotalHeight is the interface function for getting the painter canvas's total width
	Width() float64
	// Draw is the interface function for converting a slice of samples into SVG elements
	Draw() []string
	// Viewbox templates the SVG viewbox for the canvas
	Viewbox() string
}

const (
	DefaultWidth  float64 = 10
	DefaultHeight float64 = 200
)
