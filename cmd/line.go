package cmd

import "github.com/zoomoid/waveman/v2/pkg/paint/line"

type lineData struct {
	interpolation string
	fill          string
	strokeColor   string
	strokeWidth   string
	spread        float64
	height        float64
	closed        bool
	inverted      bool
}

func newLineData() *lineData {
	return &lineData{
		interpolation: string(line.DefaultInterpolation),
		fill:          line.DefaultFillColor,
		strokeColor:   line.DefaultStrokeColor,
		strokeWidth:   line.DefaultStrokeWidth,
		spread:        line.DefaultSpread,
		height:        line.DefaultHeight,
		closed:        false,
		inverted:      false,
	}
}

const (
	LineMode string = "line"
)
