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

package transform

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"
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

func transformerFactory(precision Precision) (*ReaderOptions, io.Reader) {
	f := fileFactory()
	ro := &ReaderOptions{
		Chunks:       32,
		Aggregator:   AggregatorRootMeanSquare,
		Precision:    precision,
		Downsampling: DownsamplingCenter,
	}
	return ro, f
}

func BenchmarkReaderFull(b *testing.B) {
	t, f := transformerFactory(PrecisionFull)
	New(t, f)
}

func BenchmarkReader2(b *testing.B) {
	t, f := transformerFactory(Precision2)
	New(t, f)
}

func BenchmarkReader4(b *testing.B) {
	t, f := transformerFactory(Precision4)
	New(t, f)
}

func BenchmarkReader8(b *testing.B) {
	t, f := transformerFactory(Precision8)
	New(t, f)
}

func BenchmarkReader16(b *testing.B) {
	t, f := transformerFactory(Precision16)
	New(t, f)
}

func BenchmarkReader32(b *testing.B) {
	t, f := transformerFactory(Precision32)
	New(t, f)
}

func BenchmarkReader64(b *testing.B) {
	t, f := transformerFactory(Precision64)
	New(t, f)
}

func BenchmarkReader128(b *testing.B) {
	t, f := transformerFactory(Precision128)
	New(t, f)
}

func TestNew(t *testing.T) {
	options := &ReaderOptions{
		Chunks:     64,
		Aggregator: AggregatorRootMeanSquare,
	}

	f := fileFactory()

	ctx, err := New(options, f)
	if err != nil {
		t.Fatal(err)
	}

	blocks := ctx.Blocks()
	if len(blocks) != 64 {
		t.Fatalf("wrong number of chunks, expected %d, found %d", 64, len(blocks))
	}

	for _, sample := range blocks {
		if sample != 0 {
			return
		}
	}
	t.Fatalf("block slice only contains 0 entries, expected at least one non-null sample")
}
