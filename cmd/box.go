package cmd

import "github.com/zoomoid/waveman/v2/pkg/paint/box"

type boxData struct {
	color     string
	alignment string
	height    float64
	width     float64
	rounded   float64
	gap       float64
}

func newBoxData() *boxData {
	return &boxData{
		color:     box.DefaultColor,
		alignment: string(box.DefaultAlignment),
		height:    box.DefaultHeight,
		width:     box.DefaultWidth,
		rounded:   box.DefaultRounded,
		gap:       box.DefaultRounded,
	}
}

const (
	BoxMode string = "box"
)
