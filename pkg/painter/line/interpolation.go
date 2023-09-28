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
	"fmt"
	"math"
	"strings"
)

// MonotonicCube calculates the data points for an SVG path by monotonic cubic
// hermetic interpolation such that a path appears smooth without having to
// specify control points for the bezier curves manually. It implements the
// Fritsch-Carlson method and follows roughly the structure of
// https://gionkunz.github.io/chartist-js/
func MonotonicCube(samples [][2]float64, slope func(int, []float64, []float64) []float64, startPoint [2]float64) string {
	n := len(samples)

	dxs, _, ds := calculateDeltas(samples)

	ms := slope(n, ds, dxs)

	segments := make([]string, 0)
	start := fmt.Sprintf("M %f %f", startPoint[0], startPoint[1])
	segments = append(segments, start)

	for i := 0; i < n-1; i++ {
		// first control point
		x1 := samples[i][0] + dxs[i]/3
		y1 := samples[i][1] + ms[i]*dxs[i]/3
		// second control point
		x2 := samples[i+1][0] - dxs[i]/3
		y2 := samples[i+1][1] - ms[i+1]*dxs[i]/3
		// endpoints
		x := samples[i+1][0]
		y := samples[i+1][1]
		segment := fmt.Sprintf("C %g %g, %g %g, %g %g\n", x1, y1, x2, y2, x, y)
		segments = append(segments, segment)
	}
	return strings.Join(segments, " ")
}

// Determine desired slope (m) at each point using Fritsch-Carlson method See:
// http://math.stackexchange.com/questions/45218/implementation-of-monotone-cubic-interpolation
func fritschCarlson(n int, ds []float64, dxs []float64) []float64 {
	ms := make([]float64, n)

	ms[0] = ds[0]
	ms[n-1] = ds[n-2]

	for i := 1; i < n-1; i++ {
		if ds[i] == 0 || ds[i-1] == 0 || (ds[i-1] > 0) != (ds[i] > 0) {
			ms[i] = 0
		} else {
			ms[i] = 3 * (dxs[i-1] + dxs[i]) / ((2*dxs[i]+dxs[i-1])/ds[i-1] + (dxs[i]+2*dxs[i-1])/ds[i])
			if math.IsInf(ms[i], 0) || math.IsNaN(ms[i]) {
				ms[i] = 0
			}
		}
	}
	return ms
}

// Determine desired slope (m) at each point using Steffen method See:
// http://math.stackexchange.com/questions/45218/implementation-of-monotone-cubic-interpolation
func steffen(n int, ds []float64, dxs []float64) []float64 {
	ms := make([]float64, n)

	ms[0] = ds[0]
	ms[n-1] = ds[n-2]
	for i := 1; i < n-1; i++ {
		h := 0.5 * ((dxs[i]*ds[i-1] + dxs[i-1]*ds[i]) / (dxs[i-1] + dxs[i]))
		ms[i] = (sign(dxs[i-1]) + sign(dxs[i])) * math.Min(math.Min(math.Abs(dxs[i-1]), math.Abs(dxs[i])), h)
	}
	return ms
}

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

// calculateDeltas calculates the difference between a point and its next
// neighbour and also the derivative, which is required for determining the
// slope of the interpolated cubic curve in MonotonicCube
func calculateDeltas(samples [][2]float64) (dxs []float64, dys []float64, ds []float64) {
	dxs = make([]float64, len(samples)-1)
	dys = make([]float64, len(samples)-1)
	ds = make([]float64, len(samples)-1)
	for i := 0; i < len(samples)-1; i++ {
		dxs[i] = samples[i+1][0] - samples[i][0]
		dys[i] = samples[i+1][1] - samples[i][1]
		ds[i] = dys[i] / dxs[i]
	}
	return dxs, dys, ds
}

// None calculates a line from a given slice of samples without any
// interpolation
func None(samples [][2]float64) string {
	n := len(samples)
	segments := make([]string, 0)
	for i := 0; i < n; i++ {
		x := samples[i][0]
		y := samples[i][1]
		segment := fmt.Sprintf("L %g %g \n", x, y)
		segments = append(segments, segment)
	}
	return strings.Join(segments, " ")
}
