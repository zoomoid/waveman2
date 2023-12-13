package decoder

import "io"

const (
	Channels   int = 2
	Precision  int = 2
	FrameWidth int = Channels * Precision
)

type Decoder interface {
	Length() int

	Read(samples [][2]float64) (n int, err error)

	Decode(p []byte) (sample [2]float64, n int)

	io.Seeker
	io.Closer
}
