package transform

import "math"

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
	for _, sample := range samples {
		sum += sample
	}
	return sum
}

// max determines the maximum value in a slice of float64 samples
func max(samples []float64) float64 {
	max := float64(math.Inf(-1))
	for _, sample := range samples {
		if sample > max {
			max = sample
		}
	}
	return max
}

// min determines the minimum value in a slice of float64 samples
func min(samples []float64) float64 {
	min := float64(math.Inf(1))
	for _, sample := range samples {
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
	precision := math.Pow(10, float64(digits))
	o = math.Round(o * precision)
	o = o / precision
	return o
}

// meanSquare calculates the squared mean of a slice of float64 samples
func meanSquare(samples []float64) float64 {
	// mutate slice of samples to their squares before using sum
	sc := make([]float64, len(samples))
	for i, sample := range samples {
		sc[i] = sample * sample
	}
	o := sum(sc) / float64(len(sc))
	return o
}

// rootMeanSquare determines the root-mean-square measure of a slice of float64 samples
func rootMeanSquare(samples []float64) float64 {
	return math.Sqrt(meanSquare(samples))
}
