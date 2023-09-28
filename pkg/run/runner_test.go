/*
Copyright 2022-2023 zoomoid.

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
		Chunks:       200,
		Aggregator:   transform.AggregatorRootMeanSquare,
		Precision:    transform.PrecisionFull,
		Downsampling: transform.DownsamplingCenter,
		Clamping: &transform.Clamping{
			Min: 0.2,
			Max: 1,
		},
		Window: &transform.Window{
			Algorithm: transform.Tukey,
			P:         0.1,
		},
	}

	boxOptions := &box.BoxOptions{
		Color:     "black",
		Alignment: box.AlignmentCenter,
		BoxHeight: 80,
		BoxWidth:  6,
		Rounded:   6.0 / 2.0,
		Gap:       2,
	}

	svg, err := Box(fileFactory(), transformerOptions, boxOptions)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Create("./fixtures/TestBox.svg")
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.Write([]byte(svg))
	if err != nil {
		t.Fatal(err)
	}
}

func TestLine(t *testing.T) {
	transformerOptions := &transform.ReaderOptions{
		Chunks:       24,
		Aggregator:   transform.AggregatorRootMeanSquare,
		Downsampling: transform.DownsamplingCenter,
		Precision:    transform.PrecisionFull,
		Clamping: &transform.Clamping{
			Min: 0.2,
			Max: 0.9,
		},
		Window: &transform.Window{
			Algorithm: transform.Tukey,
			P:         0.05,
		},
	}

	lineOptions := &line.LineOptions{
		Interpolation: line.InterpolationSteffen,
		// Fill:          line.DefaultFillColor,
		Fill:      "black",
		Stroke:    nil,
		Closed:    false,
		Spread:    25,
		Amplitude: 100,
		Inverted:  true,
	}

	svg, err := Line(fileFactory(), transformerOptions, lineOptions)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Create("./fixtures/TestLine.svg")
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.Write([]byte(svg))
	if err != nil {
		t.Fatal(err)
	}
}

func TestLineMirrored(t *testing.T) {
	transformerOptions := &transform.ReaderOptions{
		Chunks:       24,
		Aggregator:   transform.AggregatorRootMeanSquare,
		Downsampling: transform.DownsamplingCenter,
		Precision:    transform.PrecisionFull,
		Clamping: &transform.Clamping{
			Min: 0.2,
			Max: 0.9,
		},
		Window: &transform.Window{
			Algorithm: transform.Tukey,
			P:         0.05,
		},
	}

	lineOptions := &line.LineOptions{
		Interpolation: line.InterpolationSteffen,
		// Fill:          line.DefaultFillColor,
		Fill:      "black",
		Stroke:    nil,
		Closed:    false,
		Spread:    25,
		Amplitude: 60,
		Mirrored:  true,
	}

	svg, err := Line(fileFactory(), transformerOptions, lineOptions)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Create("./fixtures/TestLineMirrored.svg")
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.Write([]byte(svg))
	if err != nil {
		t.Fatal(err)
	}
}
