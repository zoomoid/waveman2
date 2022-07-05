package transform

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/rs/zerolog/log"
)

type TransformerMode string

const (
	TransformerModeAverage        TransformerMode = "avg"
	TransformerModeRoundedAverage TransformerMode = "rounded_avg"
	TransformerModeMax            TransformerMode = "max"
	TransformerModeMeanSquare     TransformerMode = "mean_square"
	TransformerModeRootMeanSquare TransformerMode = "rms"
)

var (
	ErrNoFile                error           = errors.New("no file given")
	DefaultTransformerMode   TransformerMode = TransformerModeRootMeanSquare
	DefaultRoundingPrecision uint            = 3
)

type ReaderOptions struct {
	Chunks    int
	Filename  string
	Mode      TransformerMode
	Precision Precision
}

type Precision uint8

const (

	// Lowest level of precision. Only every 128th sample is used
	Precision128 Precision = 128
	// Only every 64th sample is used
	Precision64 Precision = 64
	// Only every 32nd sample is used
	Precision32 Precision = 32
	// Only every 16th sample is used
	Precision16 Precision = 16
	// Only every 8th sample is used
	Precision8 Precision = 8
	// Only every 4th sample is used
	Precision4 Precision = 4
	// Only every 2nd sample is used
	Precision2 Precision = 2
	// Highest level of precision, samples are used as-is
	PrecisionFull Precision = 1
)

type ReaderContext struct {
	chunks    int
	filename  string
	mode      TransformerMode
	reader    *os.File
	decoder   beep.StreamSeekCloser
	blocks    []float64
	chunkSize int
	precision Precision
}

func New(options *ReaderOptions) (*ReaderContext, error) {
	if options.Chunks == 0 {
		options.Chunks = 64
	}
	if options.Mode == "" {
		options.Mode = DefaultTransformerMode
	}
	if options.Filename == "" {
		return nil, ErrNoFile
	}
	if options.Precision == 0 {
		options.Precision = PrecisionFull
	}

	fn, err := filepath.Abs(options.Filename)
	if err != nil {
		return nil, errors.New("failed to construct absolute path")
	}
	f, err := os.Open(fn)
	if err != nil {
		return nil, errors.New("failed to open file")
	}

	stream, _, err := mp3.Decode(f)
	if err != nil {
		return nil, errors.New("failed to construct decoder")
	}

	chunkSize := stream.Len() / options.Chunks
	blocks := make([]float64, options.Chunks)

	ctx := &ReaderContext{
		chunks:    options.Chunks,
		filename:  options.Filename,
		mode:      options.Mode,
		reader:    f,
		decoder:   stream,
		blocks:    blocks,
		chunkSize: chunkSize,
		precision: options.Precision,
	}

	err = ctx.read()
	if err != nil {
		return nil, err
	}

	return ctx, nil
}

func (r *ReaderContext) Close() {
	defer r.reader.Close()
}

func (r *ReaderContext) Blocks() []float64 {
	return r.blocks
}

func (r *ReaderContext) read() error {
	samplesPerChunk := r.chunkSize / int(r.precision)
	log.Debug().
		Int("chunks", r.chunks).
		Int("raw samples per chunk", r.chunkSize).
		Int("samples per resampled chunk", samplesPerChunk).
		Int("total samples", r.decoder.Len()).
		Send()

	for i := range r.blocks {
		b := make([][2]float64, samplesPerChunk)
		n, ok := r.decoder.Stream(b)
		err := r.decoder.Err()
		if !ok && err != nil {
			return err
		}
		seeking := r.chunkSize - n
		err = r.decoder.Seek(i*r.chunkSize + seeking)
		if err != nil {
			return err
		}

		monoSignal := toMono(b)

		block := float64(0)
		switch r.mode {
		case TransformerModeMax:
			block = max(monoSignal)
		case TransformerModeAverage:
			block = mean(monoSignal)
		case TransformerModeRoundedAverage:
			block = roundedMean(monoSignal, DefaultRoundingPrecision)
		case TransformerModeMeanSquare:
			block = meanSquare(monoSignal)
		case TransformerModeRootMeanSquare:
			block = rootMeanSquare(monoSignal)
		default:
			return fmt.Errorf("mode %s is not implemented", r.mode)
		}
		r.blocks[i] = block
	}

	// last step is to normalize the block range to [0,1]
	r.blocks = normalize(r.blocks)

	return nil
}

// func (r *ReaderContext) resample(samples [][2]float64, precision uint) [][2]float64 {
// 	if Precision(precision) == PrecisionFull {
// 		return samples
// 	}
// 	s := make([][2]float64, 0)
// 	for i := 0; i < len(samples); i += int(precision) {
// 		s = append(s, samples[i])
// 	}
// 	return s
// }
