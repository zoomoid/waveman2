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

package transform

import (
	"math"
)

// toMono sums a slice of stereo samples to mono *only in a visual way*, sort of
// like a peak meter would.
// Summing is done before taking the absolute value, such that both signals can
// interfere and inverse polarity can cancel out the signal.
func toMono(samples [][2]float64) []float64 {
	sc := make([]float64, len(samples))
	for i := 0; i < len(samples); i++ {
		// we sum channels first (allowing them to null) and then take the absolute value
		// to implement visual peak metering
		sc[i] = math.Abs(samples[i][0]+samples[i][1]) / 2
	}
	return sc
}

// normalize implements regular feature scaling to [0,1] on a slice of float64
func normalize(samples []float64) []float64 {
	localMax := max(samples)
	localMin := min(samples)
	scalar := localMax - localMin
	sc := make([]float64, len(samples))
	for i, sample := range samples {
		sc[i] = (sample - localMin) / scalar
	}
	return sc
}

// sum takes a slice of samples and sums them up.
// This function is used in various, more advanced processing methods
// downstream
func sum(samples []float64) float64 {
	sum := float64(0)
	n := len(samples)
	for i := 0; i < n; i += 1 {
		sum += samples[i]
	}
	return sum
}

// max determines the maximum value in a slice of float64 samples
func max(samples []float64) float64 {
	max := float64(math.Inf(-1))
	n := len(samples)
	for i := 0; i < n; i += 1 {
		sample := samples[i]
		if sample > max {
			max = sample
		}
	}
	return max
}

// min determines the minimum value in a slice of float64 samples
func min(samples []float64) float64 {
	min := float64(math.Inf(1))
	n := len(samples)
	for i := 0; i < n; i += 1 {
		sample := samples[i]
		if sample < min {
			min = sample
		}
	}
	return min
}

// mean determines the mean value in a slice of float64 samples
func mean(samples []float64) float64 {
	o := sum(samples) / float64(len(samples))
	return o
}

// roundedMean determines the mean value in a slice of float64 samples
// and rounds it to a given amount of digits of precision
func roundedMean(samples []float64, digits uint) float64 {
	o := mean(samples)
	roundingPrecision := math.Pow(10, float64(digits))
	o = math.Round(o * roundingPrecision)
	o = o / roundingPrecision
	return o
}

// meanSquare calculates the squared mean of a slice of float64 samples
func meanSquare(samples []float64) float64 {
	// mutate slice of samples to their squares before using sum
	o := float64(0)
	n := len(samples)
	for i := 0; i < n; i += 1 {
		sample := samples[i]
		o += sample * sample
	}
	// log.Print(o, n)
	o = o / float64(len(samples))
	return o
}

// rootMeanSquare determines the root-mean-square measure of a slice of float64 samples
func rootMeanSquare(samples []float64) float64 {
	return math.Sqrt(meanSquare(samples))
}

// clamp clips a value between a lower (min) and an upper bound (max)
func clamp(value, min, max float64) float64 {
	return math.Max(min, math.Min(max, value))
}

// tukey implements the Tukey Window function for a discrete
// set of samples given a window parameter alpha.
//
// alpha is expected to be in [0,1], but for safety reasons,
// additionally clamped to that interval.
func tukey(samples []float64, alpha float64) []float64 {
	// clamp alpha to [0,1]
	alpha = clamp(alpha, 0, 1)
	w := make([]float64, len(samples))
	N := float64(len(samples) - 1)
	for n := range w {
		if 0 <= n && n < int((alpha*N)/2) {
			s1 := (2 * math.Pi * float64(n)) / (alpha * N)
			w[n] = 0.5 * (1.0 - math.Cos(s1))
		} else if int((alpha*N)/2) <= n && n < int(N*(1-alpha/2)) {
			w[n] = 1
		} else if int(N*(1-alpha/2)) < n && n <= int(N) {
			w[n] = w[int(N)-n]
		} else {
			w[n] = 0
		}
	}
	// fold vectors
	s := make([]float64, len(samples))
	for idx, x := range samples {
		s[idx] = w[idx] * x
	}
	return s
}

// hann implements the Hann Window function for a discrete
// set of samples.
//
// see https://en.wikipedia.org/wiki/Window_function#Hann_and_Hamming_windows
func hann(samples []float64, _ float64) []float64 {
	w := make([]float64, len(samples))
	N := float64(len(samples) - 1)
	for n := range w {
		if 0 <= n && n <= int(N) {
			s1 := (2 * math.Pi * float64(n)) / N
			w[n] = 0.5 * (1.0 - math.Cos(s1))
		} else {
			w[n] = 0
		}
	}
	// fold vectors
	s := make([]float64, len(samples))
	for idx, x := range samples {
		s[idx] = w[idx] * x
	}
	return s
}

// planck_taper implements the Planck-taper Window function for a discrete
// set of samples given a window parameter eps.
// Eps is technically not bounded, but values larger than 1 will
// violate symmetry properties, and are thus of less use.
//
// see https://en.wikipedia.org/wiki/Window_function#Planck-taper_window
func planck_taper(samples []float64, eps float64) []float64 {
	w := make([]float64, len(samples))
	N := float64(len(samples) - 1)
	for n := range w {
		if 0 <= n && n < int(eps*N) {
			s1 := math.Exp((eps*N)/float64(n) - (eps*N)/(eps*N-float64(n)))
			w[n] = 1 / (1 + s1)
		} else if int(eps*N) <= n && n < int(N/2) {
			w[n] = 1
		} else if int(N/2) <= n && n <= int(N) {
			w[n] = w[int(N)-n]
		} else {
			w[n] = 0
		}
	}
	s := make([]float64, len(samples))
	for idx, x := range samples {
		s[idx] = w[idx] * x
	}
	return s
}
