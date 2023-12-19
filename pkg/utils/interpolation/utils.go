package interpolation

func toSeries(samples [][2]float64, i int) []float64 {
	if i < 0 || i > 1 {
		panic("invalid index for samples slice")
	}

	s := make([]float64, len(samples))
	for idx, sample := range samples {
		s[idx] = sample[i]
	}
	return s
}

// func toPoints(xs []float64, ys []float64) [][2]float64 {
// 	if len(xs) != len(ys) {
// 		panic("series slices are not equal in length")
// 	}
// 	ps := make([][2]float64, len(xs))
// 	for i := 0; i < len(xs); i++ {
// 		ps[i] = [2]float64{xs[i], ys[i]}
// 	}
// 	return ps
// }

// sign is the default mathematical sign function, returning
// -1 if d < 0, 0 if d == 0, and 1 if d > 0.
func sign(d float64) float64 {
	if d < 0 {
		return -1
	}
	if d > 0 {
		return 1
	}
	return 0
}

// deltas calculates the difference between a point and its next
// neighbour and also the derivative, which is required for determining the
// slope of the interpolated cubic curve in MonotonicCube
func deltas(samples [][2]float64) (dxs []float64, dys []float64, ds []float64) {
	n := len(samples)
	dxs = make([]float64, n)
	dys = make([]float64, n)
	ds = make([]float64, n-1)
	for i := 0; i < n-1; i++ {
		dxs[i] = samples[i+1][0] - samples[i][0]
		dys[i] = samples[i+1][1] - samples[i][1]
		ds[i] = dys[i] / dxs[i]
	}
	return dxs, dys, ds
}
