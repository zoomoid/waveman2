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

package line

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/rs/zerolog/log"
	"github.com/zoomoid/waveman2/pkg/painter"
)

type Interpolation string

const (
	// InterpolationNone means no interpolated points and a piecewise linear curve as the result
	InterpolationNone Interpolation = "none"
	// InterpolationFritschCarlson applies the Fritsch-Carlson method for
	// interpolating hermitic cubic splines to fit the data points
	//
	// See: http://math.stackexchange.com/questions/45218/implementation-of-monotone-cubic-interpolation
	InterpolationFritschCarlson Interpolation = "fritsch-carlson"
	// InterpolationSteffen applies the Steffen method for interpolating
	// hermetic cubic splices to fit the data points.
	//
	// See: http://math.stackexchange.com/questions/45218/implementation-of-monotone-cubic-interpolation
	InterpolationSteffen Interpolation = "steffen"
	// InterpolationEmpty is for catching uninitialized interpolation modes
	InterpolationEmpty Interpolation = ""
)

var Interpolations = []string{"fritsch-carlson", "none", "steffen"}

const (
	DefaultPathTemplate string = `<path d="{{.Points}}" fill="{{.Fill}}" stroke="{{.Stroke.Color}}" stroke-width="{{.Stroke.Width}}" />`
)

type Stroke struct {
	// Color is the CSS-compliant stroke color used for the path
	Color string
	// Width is the stroke width used for the path
	Width float64
}

type LineOptions struct {
	// Interpolation choses the point interpolation mode
	Interpolation Interpolation
	// Fill is the CSS-compliant color value for the fill
	Fill string
	// Stroke is a struct defining properties of the stroke, namely color and
	// width
	Stroke *Stroke
	// Closed makes the bezier curve a closed one by appending the "Z" parameter
	// to the path
	Closed bool
	// Spread is the horizontal scaling factor by which all indices are scaled
	Spread float64
	// Amplitude is the vertical scaling factor by which all sample points are scaled
	// up. Amplitude also determines the total canvas height when normalized samples are used
	Amplitude float64
	// Inverted transforms the SVG group to be horizontically flipped
	Inverted bool
	// Mirrored mirros the line along the x axis
	Mirrored bool
}

const (
	// Default interpolation mode
	DefaultInterpolation Interpolation = InterpolationFritschCarlson
	// Default stroke width
	DefaultStrokeWidth float64 = 0
	// Default stroke color
	DefaultStrokeColor string = "none"
	// Default fill color
	DefaultFillColor string = "rgba(0 0 0 / 0.5)"
	// Default horizontal spread of data points
	DefaultSpread = float64(10)
	// DefaultHeight of a canvas is 200px
	DefaultHeight = float64(200)
)

var (
	DefaultStroke Stroke = Stroke{
		Color: DefaultStrokeColor,
		Width: DefaultStrokeWidth,
	}
)

// Compile-time type checking for LinePainter to implement all functions required
// by the Painter interface
var _ painter.Painter = &LinePainter{}

type LinePainter struct {
	// Embed all painter options, i.e., data points
	*painter.PainterOptions
	// Embed all options for the line painter
	*LineOptions
}

// NewPainter constructs a new Line painter with the passed options and fills in defaults
// for missing fields
func NewPainter(painter *painter.PainterOptions, options *LineOptions) *LinePainter {
	if options.Interpolation == InterpolationEmpty {
		options.Interpolation = DefaultInterpolation
	}
	if options.Fill == "" {
		options.Fill = DefaultFillColor
	}
	if options.Stroke == nil {
		options.Stroke = &Stroke{
			Color: DefaultStrokeColor,
			Width: DefaultStrokeWidth,
		}
	}
	if options.Stroke.Color == "" {
		options.Stroke.Color = DefaultStrokeColor
	}
	if options.Stroke.Width < 0 {
		options.Stroke.Width = DefaultStrokeWidth
	}
	if options.Amplitude == 0 {
		options.Amplitude = DefaultHeight
	}
	if options.Spread == 0 {
		options.Spread = DefaultSpread
	}

	return &LinePainter{
		PainterOptions: painter,
		LineOptions:    options,
	}
}

// Height is the maximum height of any point in the curve
// since the interpolation does not overshoot, the maximum height, thus the
// total height, is the same as the Height scaling, as long as normalized data
// is passed to the painter.
func (l *LinePainter) Height() float64 {
	return l.Amplitude
}

// Width is the width of the canvas, that is, the width the path spans horizontally.
func (l *LinePainter) Width() float64 {
	return float64(len(l.PainterOptions.Data)-1)*l.Spread + 2*l.Spread
}

// Draw implements line drawing with optional interpolation to smooth out the curve.
// Curves are, by default, anchored to the top-left corner. If you want to change
// this, consider transforming the entire SVG in-post using CSS transforms.
func (l *LinePainter) Draw() []string {
	output := make([]string, 3)

	pathTemplate := template.New("path")
	pathTemplate.Parse(DefaultPathTemplate)

	// make a slice of pairs that have the spread x values and their y values paired
	samples := make([][2]float64, 0, len(l.Data)+2)

	var direction float64 = -1
	var offset float64 = l.Amplitude
	if l.Inverted && l.Mirrored {
		log.Warn().Msg("LinePainter.Inverted and LinePainter.Mirrored are mutually exclusive, deactiving LinePainter.Inverted...")
		l.Inverted = false
	}

	if l.Inverted {
		direction = 1
		offset = 0
	}

	samples = append(samples, [2]float64{0, offset})
	for i, sample := range l.Data {
		// offset samples in X direction by one unit of spread to account for start points
		samples = append(samples, [2]float64{float64(i)*l.Spread + l.Spread, direction*l.Amplitude*sample + offset})
	}
	samples = append(samples, [2]float64{l.Width(), offset})

	line := ""
	switch l.Interpolation {
	case InterpolationSteffen:
		line = MonotonicCube(samples, steffen, samples[0])
	case InterpolationNone:
		line = None(samples)
	case InterpolationFritschCarlson:
		line = MonotonicCube(samples, fritschCarlson, samples[0])
	}

	if l.Closed {
		line += " Z\n"
	}

	bindings := struct {
		Fill   string
		Points string
		Stroke *Stroke
	}{
		Fill:   l.Fill,
		Points: line,
		Stroke: l.Stroke,
	}

	templateBuf := &bytes.Buffer{}
	pathTemplate.Execute(templateBuf, bindings)

	output = append(output, `<g style="transform-origin: center center;">`)
	output = append(output, templateBuf.String())
	output = append(output, "</g>")

	if l.Mirrored {
		// add the samples again, but in the other direction for line continuity
		mirrorPoints := make([][2]float64, 0, len(l.Data)+2)
		mirrorPoints = append(mirrorPoints, [2]float64{0, offset})
		for i, x := range l.Data {
			mirrorPoints = append(mirrorPoints, [2]float64{float64(i)*l.Spread + l.Spread, 2*offset - (1-x)*l.Amplitude})
		}
		mirrorPoints = append(mirrorPoints, [2]float64{l.Width(), offset})

		line := ""
		switch l.Interpolation {
		case InterpolationSteffen:
			line = MonotonicCube(mirrorPoints, steffen, mirrorPoints[0])
		case InterpolationNone:
			line = None(mirrorPoints)
		case InterpolationFritschCarlson:
			line = MonotonicCube(mirrorPoints, fritschCarlson, mirrorPoints[0])
		}

		if l.Closed {
			line += " Z\n"
		}

		bindings := struct {
			Fill   string
			Points string
			Stroke *Stroke
		}{
			Fill:   l.Fill,
			Points: line,
			Stroke: l.Stroke,
		}

		pathTemplate.Execute(templateBuf, bindings)
		output = append(output, `<g style="transform-origin: center center;">`)
		output = append(output, templateBuf.String())
		output = append(output, "</g>")
	}

	return output
}

func (l *LinePainter) Viewbox() string {
	// calculate the viewBox: we need to offset the viewbox by the stroke width in all directions to not clip it
	offset := l.Stroke.Width
	mirrorMul := 1
	if l.Mirrored {
		mirrorMul = 2
	}
	return fmt.Sprintf("%f %f %f %f", (-1 * offset), offset, l.Width()+offset, float64(mirrorMul)*l.Height()+2*offset)
}