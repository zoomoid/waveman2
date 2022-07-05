package run

import (
	"fmt"
	"testing"

	"github.com/zoomoid/waveman/v2/pkg/paint/box"
	"github.com/zoomoid/waveman/v2/pkg/paint/line"
	"github.com/zoomoid/waveman/v2/pkg/transform"
)

func TestBox(t *testing.T) {
	transformerOptions := &transform.ReaderOptions{
		Chunks:    64,
		Filename:  "../../hack/Morgendämmerung.mp3",
		Mode:      transform.TransformerModeRootMeanSquare,
		Precision: transform.Precision2,
	}

	boxOptions := &box.BoxOptions{
		Color:     "black",
		Alignment: box.AlignmentCenter,
		Height:    200,
		Width:     10,
		Rounded:   5,
		Gap:       2,
	}

	svg, err := Box(transformerOptions, boxOptions)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Print(svg)
}

func TestLine(t *testing.T) {
	transformerOptions := &transform.ReaderOptions{
		Chunks:   64,
		Filename: "../../hack/Morgendämmerung.mp3",
		Mode:     transform.TransformerModeRootMeanSquare,
	}

	lineOptions := &line.LineOptions{
		Interpolation: line.DefaultInterpolation,
		Fill:          line.DefaultFillColor,
		Stroke: &line.Stroke{
			Color: line.DefaultStrokeColor,
			Width: "2px",
		},
		Closed: true,
		Spread: 10,
		Height: 50,
	}

	svg, err := Line(transformerOptions, lineOptions)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Print(svg)
}
