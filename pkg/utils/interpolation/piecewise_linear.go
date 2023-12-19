package interpolation

import (
	"fmt"
	"strings"
)

// For usage with SVG, we don't actually need to do anything, just use "L" in ToCurve
type PiecewiseLinear struct {
	points [][2]float64
}

func (pl *PiecewiseLinear) Interpolate(samples [][2]float64) {
	pl.points = samples
}

func (pl *PiecewiseLinear) ToCurve(startPoint [2]float64) string {
	n := len(pl.points)
	segmentWriter := &strings.Builder{}
	fmt.Fprintf(segmentWriter, "M %f %f ", startPoint[0], startPoint[1])

	for i := 0; i < n; i++ {
		segmentWriter.WriteString(LinePoint(pl.points[i]).ToSegment())
	}
	return segmentWriter.String()
}

type LinePoint [2]float64

func (pl *PiecewiseLinear) Points() []LinePoint {
	ps := make([]LinePoint, len(pl.points))
	for i, p := range pl.points {
		ps[i] = LinePoint(p)
	}
	return ps
}

func (p LinePoint) Clone() LinePoint {
	return LinePoint([2]float64{p[0], p[1]})
}

func (p LinePoint) ToSegment() string {
	return fmt.Sprintf("L %g %g ", p[0], p[1])
}

func (LinePoint) ZeroSegment() string {
	return "l 0 0"
}
