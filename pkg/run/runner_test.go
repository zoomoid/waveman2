/*
Copyright 2022 zoomoid.

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
			Width: 2,
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
