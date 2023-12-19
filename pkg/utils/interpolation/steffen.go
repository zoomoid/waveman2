package interpolation

import (
	"fmt"
	"math"
	"strings"
)

type Steffen struct {
	points []CubicCurvePoint

	// c piecewiseCubic
}

var _ Interpolator = &Steffen{}

func (st *Steffen) Interpolate(samples [][2]float64) {
	xs := toSeries(samples, 0)
	ys := toSeries(samples, 1)

	dxs, _, ds := deltas(samples)

	ms := make([]float64, len(ds))
	n := len(ms)

	// See http://math.stackexchange.com/questions/45218/implementation-of-monotone-cubic-interpolation
	for i := 1; i < n-1; i++ {
		h := 0.5 * ((dxs[i]*ds[i-1] + dxs[i-1]*ds[i]) / (dxs[i-1] + dxs[i]))
		ms[i] = (sign(dxs[i-1]) + sign(dxs[i])) * math.Min(math.Min(math.Abs(dxs[i-1]), math.Abs(dxs[i])), h)
	}

	m := n - 1

	ms[0] = ds[0]
	ms[m] = ds[m-1]

	// st.c.fit(xs, ys, ms)

	st.points = make([]CubicCurvePoint, m)
	for i := 0; i < m; i++ {
		x1 := xs[i] + dxs[i]/3 // current point + 1/3 the distance
		// see fritsch_carlson.go for why we don't use predict here
		// y1 := st.c.predict(x1)
		y1 := ys[i] + ms[i]*dxs[i]/3.0

		x2 := xs[i+1] - dxs[i]/3 // next point - 1/3 the distance
		// y2 := st.c.predict(x2)
		y2 := ys[i+1] - ms[i+1]*dxs[i]/3.0

		st.points[i] = CubicCurvePoint{
			Root: [2]float64{xs[i+1], ys[i+1]}, // the "next" point is the root of the cubic
			C1:   [2]float64{x1, y1},
			C2:   [2]float64{x2, y2},
		}
	}
}

func (st *Steffen) ToCurve(s [2]float64) string {
	segmentWriter := &strings.Builder{}
	fmt.Fprintf(segmentWriter, "M %f %f ", s[0], s[1])
	for _, p := range st.points {
		segmentWriter.WriteString(p.ToSegment())
	}
	return segmentWriter.String()
}

func (st *Steffen) Points() []CubicCurvePoint {
	return st.points
}
