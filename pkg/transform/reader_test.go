package transform

import (
	"testing"
)

func transformerFactory(precision Precision) *ReaderOptions {
	return &ReaderOptions{
		Chunks:       32,
		Filename:     "../../hack/Morgendämmerung.mp3",
		Aggregator:   AggregatorRootMeanSquare,
		Precision:    precision,
		Downsampling: DownsamplingCenter,
	}
}

func BenchmarkReaderFull(b *testing.B) {
	t := transformerFactory(PrecisionFull)
	New(t)
}

func BenchmarkReader2(b *testing.B) {
	t := transformerFactory(Precision2)
	New(t)
}

func BenchmarkReader4(b *testing.B) {
	t := transformerFactory(Precision4)
	New(t)
}

func BenchmarkReader8(b *testing.B) {
	t := transformerFactory(Precision8)
	New(t)
}

func BenchmarkReader16(b *testing.B) {
	t := transformerFactory(Precision16)
	New(t)
}

func BenchmarkReader32(b *testing.B) {
	t := transformerFactory(Precision32)
	New(t)
}

func BenchmarkReader64(b *testing.B) {
	t := transformerFactory(Precision64)
	New(t)
}

func BenchmarkReader128(b *testing.B) {
	t := transformerFactory(Precision128)
	New(t)
}

func TestNew(t *testing.T) {
	options := &ReaderOptions{
		Chunks:     64,
		Filename:   "../../hack/Morgendämmerung.mp3",
		Aggregator: AggregatorRootMeanSquare,
	}

	ctx, err := New(options)
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
