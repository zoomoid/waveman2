package box

type Position struct {
	// Horizontal position value of the uppper-left corner of the rectangle
	x float64
	// Vertical position value of the upper-left corner of the rectangle
	y float64
}

// X position of the rectangle in 2D space
func (p *Position) X() float64 {
	return p.x
}

// Y position of the rectangle in 2D space
func (p *Position) Y() float64 {
	return p.y
}

type Dimensions struct {
	// height of the rectangle shape
	height float64
	// width of the rectangle shape
	width float64
}

// Height returns the rectangle's height
func (d *Dimensions) Height() float64 {
	return d.height
}

// Width returns the rectangle's width
func (d *Dimensions) Width() float64 {
	return d.width
}

type Rectangle struct {
	// Position in 2D space
	Position
	// Height and width of the rectangle
	Dimensions
	// Rounded is the rectangle edge rounding
	Rounded float64
	// Color is the fill color of the rectangle
	Color string
}
