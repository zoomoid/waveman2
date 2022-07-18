package painter

type PainterOptions struct {
	// Data contains all sample points to use in a drawing context
	Data   []float64
	Height float64
	Width  float64
}

type Painter interface {
	// Height is the interface function for getting the painter canvas's total height
	Height() float64
	// TotalHeight is the interface function for getting the painter canvas's total width
	Width() float64
	// Draw is the interface function for converting a slice of samples into SVG elements
	Draw() []string
}

const (
	DefaultWidth  float64 = 10
	DefaultHeight float64 = 200
)
