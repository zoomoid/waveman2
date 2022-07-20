package run

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/zoomoid/waveman2/pkg/painter/box"
	"github.com/zoomoid/waveman2/pkg/painter/line"
	"github.com/zoomoid/waveman2/pkg/transform"
)

const (
	TestFile = "../../hack/Morgendämmerung.mp3"
)

func fileFactory() io.Reader {
	fn, err := filepath.Abs(TestFile)
	if err != nil {
		log.Fatal(errors.New("failed to construct absolute path"))
	}
	f, err := os.Open(fn)
	if err != nil {
		log.Fatal(errors.New("failed to construct absolute path"))
	}
	return f
}

func TestBox(t *testing.T) {
	transformerOptions := &transform.ReaderOptions{
		Chunks:       50,
		Aggregator:   transform.AggregatorRootMeanSquare,
		Precision:    transform.Precision8,
		Downsampling: transform.DownsamplingCenter,
	}

	boxOptions := &box.BoxOptions{
		Color:     "black",
		Alignment: box.AlignmentCenter,
		BoxHeight: 200,
		BoxWidth:  10,
		Rounded:   5,
		Gap:       2,
	}

	svg, err := Box(fileFactory(), transformerOptions, boxOptions)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Print(svg)
}

func TestLine(t *testing.T) {
	transformerOptions := &transform.ReaderOptions{
		Chunks:       64,
		Aggregator:   transform.AggregatorRootMeanSquare,
		Downsampling: transform.DownsamplingCenter,
		Precision:    transform.Precision4,
	}

	lineOptions := &line.LineOptions{
		Interpolation: line.DefaultInterpolation,
		Fill:          line.DefaultFillColor,
		Stroke: &line.Stroke{
			Color: line.DefaultStrokeColor,
			Width: "2px",
		},
		Closed:    true,
		Spread:    10,
		Amplitude: 50,
	}

	svg, err := Line(fileFactory(), transformerOptions, lineOptions)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Print(svg)
}
