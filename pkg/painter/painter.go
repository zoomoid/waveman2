package painter

type PainterOptions struct {
	// Data contains all sample points to use in a drawing context
	Data []float64
}

type Painter interface {
	// TotalHeight is the interface function for getting the painter canvas's total height
	TotalHeight() float64
	// TotalHeight is the interface function for getting the painter canvas's total width
	TotalWidth() float64
	// Draw is the interface function for converting a slice of samples into SVG elements
	Draw() []string
}
