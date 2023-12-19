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

package interpolation

// Determine desired slope (m) at each point using Fritsch-Carlson method
//
// See http://math.stackexchange.com/questions/45218/implementation-of-monotone-cubic-interpolation
// func fritschCarlson(n int, ds []float64, dxs []float64) []float64 {
// 	ms := make([]float64, n)

// 	ms[0] = ds[0]
// 	ms[n-1] = ds[n-2]

// 	for i := 1; i < n-1; i++ {
// 		if ds[i] == 0 || ds[i-1] == 0 || (ds[i-1] > 0) != (ds[i] > 0) {
// 			ms[i] = 0
// 		} else {
// 			ms[i] = 3 * (dxs[i-1] + dxs[i]) / ((2*dxs[i]+dxs[i-1])/ds[i-1] + (dxs[i]+2*dxs[i-1])/ds[i])
// 			if math.IsInf(ms[i], 0) || math.IsNaN(ms[i]) {
// 				ms[i] = 0
// 			}
// 		}
// 	}
// 	return ms
// }

// Determine desired slope (m) at each point using Steffen method
//
// See http://math.stackexchange.com/questions/45218/implementation-of-monotone-cubic-interpolation
// func steffen(n int, ds []float64, dxs []float64) []float64 {
// 	ms := make([]float64, n)

// 	ms[0] = ds[0]
// 	ms[n-1] = ds[n-2]
// 	for i := 1; i < n-1; i++ {
// 		h := 0.5 * ((dxs[i]*ds[i-1] + dxs[i-1]*ds[i]) / (dxs[i-1] + dxs[i]))
// 		ms[i] = (sign(dxs[i-1]) + sign(dxs[i])) * math.Min(math.Min(math.Abs(dxs[i-1]), math.Abs(dxs[i])), h)
// 	}
// 	return ms
// }
