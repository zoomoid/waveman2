package transform

import "testing"

func TestSum(t *testing.T) {
	s := sum([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	if s != ((10 * 11) / 2) {
		t.Errorf("Expected %d, found %g", ((10 * 11) / 2), s)
	}
}
