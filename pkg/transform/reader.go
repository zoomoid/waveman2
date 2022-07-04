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
	Steps    int
	Filename string
	Mode     TransformerMode
}

type ReaderContext struct {
	steps     int
	filename  string
	mode      TransformerMode
	reader    *os.File
	decoder   beep.StreamSeekCloser
	blocks    []float64
	chunkSize int
}

func New(options *ReaderOptions) (*ReaderContext, error) {
	if options.Steps == 0 {
		options.Steps = 64
	}
	if options.Mode == "" {
		options.Mode = DefaultTransformerMode
	}
	if options.Filename == "" {
		return nil, ErrNoFile
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

	chunkSize := stream.Len() / options.Steps
	blocks := make([]float64, chunkSize)

	ctx := &ReaderContext{
		steps:     options.Steps,
		filename:  options.Filename,
		mode:      options.Mode,
		reader:    f,
		decoder:   stream,
		blocks:    blocks,
		chunkSize: chunkSize,
	}

	err = ctx.read()
	if err != nil {
		return nil, errors.New("failed to decode file")
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
	log.Debug().Msg(fmt.Sprintf("%d chunks, %d samples per chunk, %d block", r.steps, r.chunkSize, len(r.blocks)))
	for i := range r.blocks {
		b := make([][2]float64, r.chunkSize)
		_, ok := r.decoder.Stream(b)
		if !ok {
			return errors.New("failed to decode block")
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
