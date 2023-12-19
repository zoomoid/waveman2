package interpolation

import (
	"fmt"
)

// type piecewiseCubic struct {
// 	coeffs [][4]float64

// 	xs []float64

// 	lastY float64

// 	lastDxDy float64
// }

// func (c *piecewiseCubic) fit(xs, ys, dydxs []float64) {
// 	n := len(xs)
// 	if len(ys) != n {
// 		panic("xs and ys are not equal in length")
// 	}
// 	if len(dydxs) != n {
// 		panic("xs and dydxs are not equal in length")
// 	}
// 	if n < 2 {
// 		panic("too few points")
// 	}

// 	m := n - 1

// 	c.coeffs = make([][4]float64, m)
// 	for i := 0; i < m; i++ {
// 		dx := xs[i+1] - xs[i]
// 		dy := ys[i+1] - ys[i]

// 		// a_0
// 		a0 := ys[i]
// 		// a_1 - slopes previously computed from the interpolators
// 		a1 := dydxs[i]
// 		// a_2
// 		a2 := (3*dy - (2*dydxs[i]+dydxs[i+1])*dx) / dx / dx
// 		// a_3
// 		a3 := (-2*dy + (dydxs[i]+dydxs[i+1])*dx) / dx / dx / dx

// 		c.coeffs[i] = [4]float64{a0, a1, a2, a3}
// 	}

// 	c.xs = append(c.xs[:0], xs...)
// 	c.lastY = ys[m]
// 	c.lastDxDy = dydxs[m]
// }

// func (c *piecewiseCubic) fitWithSecondDerivative(xs, ys, d2ydx2s []float64) {
// 	n := len(xs)
// 	if n < 2 {
// 		panic("too few points")
// 	}
// 	m := n - 1
// 	c.coeffs = make([][4]float64, m)
// 	for i := 0; i < m; i++ {
// 		dx := xs[i+1] - xs[i]
// 		dy := ys[i+1] - ys[i]

// 		dm := d2ydx2s[i+1] - d2ydx2s[i]

// 		a0 := ys[i]
// 		a1 := (dy - (d2ydx2s[i]+dm/3)*dx*dx/2) / dx
// 		a2 := d2ydx2s[i] / 2
// 		a3 := dm / 6 / dx

// 		c.coeffs[i] = [4]float64{a0, a1, a2, a3}
// 	}
// 	c.xs = append(c.xs[:0], xs...)
// 	c.lastY = ys[m]

// 	lastDx := xs[m] - xs[m-1]
// 	a1 := c.coeffs[m-1][1]
// 	a2 := c.coeffs[m-1][2]
// 	a3 := c.coeffs[m-1][3]
// 	c.lastDxDy = a1 + 2*a2*lastDx + 3*a3*lastDx*lastDx
// }

// func (c *piecewiseCubic) predict(x float64) float64 {
// 	i := findSegment(c.xs, x)
// 	if i < 0 {
// 		return c.coeffs[0][0]
// 	}
// 	m := len(c.xs) - 1
// 	if x == c.xs[i] {
// 		if i < m {
// 			return c.coeffs[i][0]
// 		}
// 		return c.lastY
// 	}
// 	if i == m {
// 		return c.lastY
// 	}
// 	dx := x - c.xs[i]
// 	a := c.coeffs[i]

// 	return ((a[3]*dx+a[2])*dx+a[1])*dx + a[0]
// }

// func (c *piecewiseCubic) predictDerivative(x float64) float64 {
// 	i := findSegment(c.xs, x)
// 	if i < 0 {
// 		return c.coeffs[0][0]
// 	}
// 	m := len(c.xs) - 1
// 	if x == c.xs[i] {
// 		if i < m {
// 			return c.coeffs[i][0]
// 		}
// 		return c.lastDxDy
// 	}
// 	if i == m {
// 		return c.lastDxDy
// 	}
// 	dx := x - c.xs[i]
// 	a := c.coeffs[i]

// 	return (3*a[3]*dx+2*a[2])*dx + a[1]
// }

type CubicCurvePoint struct {
	Root [2]float64
	C1   [2]float64
	C2   [2]float64
}

func (p CubicCurvePoint) Clone() CubicCurvePoint {
	return CubicCurvePoint{
		Root: [2]float64{p.Root[0], p.Root[1]},
		C1:   [2]float64{p.C1[0], p.C1[1]},
		C2:   [2]float64{p.C2[0], p.C2[1]},
	}
}

func (p CubicCurvePoint) ToSegment() string {
	return fmt.Sprintf("C %g %g, %g %g, %g %g ", p.C1[0], p.C1[1], p.C2[0], p.C2[1], p.Root[0], p.Root[1])
}

func (CubicCurvePoint) ZeroSegment() string {
	return "c 0 0, 0 0, 0 0"
}

// findSegment returns 0 <= i < len(xs) such that xs[i] <= x < xs[i + 1], where xs[len(xs)]
// is assumed to be +Inf. If no such i is found, it returns -1. It assumes that len(xs) >= 2
// without checking.
// func findSegment(xs []float64, x float64) int {
// 	return sort.Search(len(xs), func(i int) bool { return xs[i] > x }) - 1
// }
