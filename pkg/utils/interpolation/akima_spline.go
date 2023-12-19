package interpolation

import (
	"fmt"
	"math"
	"strings"
)

type AkimaSpline struct {
	points []CubicCurvePoint

	// c piecewiseCubic
}

var _ Interpolator = &AkimaSpline{}

func (as *AkimaSpline) Interpolate(samples [][2]float64) {
	n := len(samples)
	xs := toSeries(samples, 0)
	ys := toSeries(samples, 1)

	dydxs := make([]float64, n)

	if n == 2 {
		dx := xs[1] - xs[0]
		slope := (ys[1] - ys[0]) / dx
		dydxs[0] = slope
		dydxs[1] = slope

		// as.c.fit(xs, ys, dydxs)
		return
	}

	slopes := as.slopes(xs, ys)

	for i := 0; i < n; i++ {
		wL, wR := as.weights(slopes, i)
		dydxs[i] = as.weightedAverage(slopes[i+1], slopes[i+2], wL, wR)
	}

	// as.c.fit(xs, ys, dydxs)

	m := n - 1
	as.points = make([]CubicCurvePoint, m)
	for i := 0; i < m; i++ {
		dx := (xs[i+1] - xs[i]) / 3.0
		x1 := xs[i] + dx // current point + 1/3 the distance
		// y1 := as.c.predict(x1)
		y1 := ys[i] + dydxs[i]*(xs[i+1]-xs[i])/3.0

		x2 := xs[i+1] - dx // next point - 1/3 the distance
		// y2 := as.c.predict(x2)
		y2 := ys[i+1] - dydxs[i+1]*(xs[i+1]-xs[i])/3

		as.points[i] = CubicCurvePoint{
			Root: [2]float64{xs[i+1], ys[i+1]}, // the "next" point is the root of the cubic
			C1:   [2]float64{x1, y1},
			C2:   [2]float64{x2, y2},
		}
	}
}

func (as *AkimaSpline) ToCurve(s [2]float64) string {
	segmentWriter := &strings.Builder{}
	fmt.Fprintf(segmentWriter, "M %f %f ", s[0], s[1])
	for _, p := range as.points {
		segmentWriter.WriteString(p.ToSegment())
	}
	return segmentWriter.String()
}

func (as *AkimaSpline) Points() []CubicCurvePoint {
	return as.points
}

func (*AkimaSpline) slopes(xs, ys []float64) []float64 {
	n := len(xs)
	if n <= 2 {
		panic("akima: too few points")
	}

	m := n + 3
	slopes := make([]float64, m)
	for i := 2; i < m-2; i++ {
		dx := xs[i-1] - xs[i-2]
		slopes[i] = (ys[i-1] - ys[i-2]) / dx
	}
	slopes[0] = 3*slopes[2] - 2*slopes[3]
	slopes[1] = 2*slopes[2] - slopes[3]

	slopes[m-2] = 2*slopes[m-3] - slopes[m-4]
	slopes[m-1] = 3*slopes[m-3] - 2*slopes[m-4]
	return slopes
}

func (*AkimaSpline) weightedAverage(v1, v2, w1, w2 float64) float64 {
	w := w1 + w2
	if w > 0 {
		return (v1*w1 + v2*w2) / w
	}
	return 0.5*v1 + 0.5*v2
}

func (*AkimaSpline) weights(slopes []float64, i int) (float64, float64) {
	wL := math.Abs(slopes[i+2] - slopes[i+3])
	wR := math.Abs(slopes[i+1] - slopes[i])
	return wL, wR
}
