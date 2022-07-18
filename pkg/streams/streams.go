package streams

import (
	"io"
	"os"
)

type IO struct {
	In     io.Reader
	Out    io.Writer
	ErrOut io.Writer
}

var (
	// Default I/O stream binding
	DefaultStreams *IO = &IO{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}
)
