package interpolation

import (
	"fmt"
	"math"
	"strings"
)

type FritschCarlson struct {
	points []CubicCurvePoint

	// c piecewiseCubic
}

var _ Interpolator = &FritschCarlson{}

func (fc *FritschCarlson) Interpolate(samples [][2]float64) {
	xs := toSeries(samples, 0)
	ys := toSeries(samples, 1)

	n := len(xs)
	ms := make([]float64, n)

	dxs, _, ds := deltas(samples)
	m := len(ds)

	pd := ds[0]
	// See http://math.stackexchange.com/questions/45218/implementation-of-monotone-cubic-interpolation
	for i := 1; i < m; i++ {
		d := ds[i]
		if d == 0 || pd == 0 || (pd > 0) != (d > 0) {
			ms[i] = 0
		} else {
			// ms[i] = 3 * (xs[i+1] + xs[i-1]) / ((2*xs[i+1]-xs[i-1]-xs[i])/ds[i-1] + (xs[i+1]+xs[i]-2*xs[i-1])/ds[i])
			ms[i] = 3 * (dxs[i-1] + dxs[i]) / ((2*dxs[i]+dxs[i-1])/ds[i-1] + (dxs[i]+2*dxs[i-1])/ds[i])
			if math.IsInf(ms[i], 0) || math.IsNaN(ms[i]) {
				ms[i] = 0
			}
		}
		pd = d
	}

	ms[0] = fc.edges(xs, ds, true)
	ms[m] = fc.edges(xs, ds, false)

	// fc.c.fit(xs, ys, ms)

	fc.points = make([]CubicCurvePoint, m)
	for i := 0; i < m; i++ {
		x1 := xs[i] + dxs[i]/3.0 // current point + 1/3 the distance
		// y1 := fc.c.predict(x1)
		// instead of "predicting" the value by proper interpolation, we want to generate
		// control points that map better to bezier curves. Cubic interpolation will generate
		// points on a cubic that resembles the *output* of what a cubic bezier curve would yield.
		y1 := ys[i] + ms[i]*dxs[i]/3.0

		x2 := xs[i+1] - dxs[i]/3.0 // next point - 1/3 the distance
		// y2 := fc.c.predict(x2)
		y2 := ys[i+1] - ms[i+1]*dxs[i]/3.0

		fc.points[i] = CubicCurvePoint{
			Root: [2]float64{xs[i+1], ys[i+1]}, // the "next" point is the root of the cubic
			C1:   [2]float64{x1, y1},
			C2:   [2]float64{x2, y2},
		}
	}
}

func (fc *FritschCarlson) ToCurve(s [2]float64) string {
	segmentWriter := &strings.Builder{}
	fmt.Fprintf(segmentWriter, "M %f %f ", s[0], s[1])
	for _, p := range fc.points {
		segmentWriter.WriteString(p.ToSegment())
	}
	return segmentWriter.String()
}

func (fc *FritschCarlson) Points() []CubicCurvePoint {
	return fc.points
}

func (fc *FritschCarlson) edges(xs []float64, slopes []float64, leftEdge bool) float64 {
	n := len(xs)
	var dE, dI, h, hE, f float64
	if leftEdge {
		dE = slopes[0]
		dI = slopes[1]
		xE := xs[0]
		xM := xs[1]
		xI := xs[2]
		hE = xM - xE
		h = xI - xE
		f = xM + xI - 2*xE
	} else {
		dE = slopes[n-2]
		dI = slopes[n-3]
		xE := xs[n-1]
		xM := xs[n-2]
		xI := xs[n-3]
		hE = xE - xM
		h = xE - xI
		f = 2*xE - xI - xM
	}
	g := (f*dE - hE*dI) / h
	if g*dE <= 0 {
		return 0
	}
	if dE*dI <= 0 && math.Abs(g) > 3*math.Abs(dE) {
		return 3 * dE
	}
	return g
}
