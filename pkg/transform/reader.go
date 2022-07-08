package transform

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

type DownsamplingMode string

const (
	DownsamplingNone   DownsamplingMode = "none"
	DownsamplingHead   DownsamplingMode = "head"
	DownsamplingCenter DownsamplingMode = "center"
	DownsamplingTail   DownsamplingMode = "tail"
)

type Aggregator string

const (
	TransformerModeAverage        Aggregator = "avg"
	TransformerModeRoundedAverage Aggregator = "rounded-avg"
	TransformerModeMax            Aggregator = "max"
	TransformerModeMeanSquare     Aggregator = "mean-square"
	TransformerModeRootMeanSquare Aggregator = "rms"
)

var (
	ErrNoFile                error            = errors.New("no file given")
	DefaultTransformerMode   Aggregator       = TransformerModeRootMeanSquare
	DefaultRoundingPrecision uint             = 3
	DefaultDownsamplingMode  DownsamplingMode = DownsamplingNone
	DefaultPrecision         Precision        = PrecisionFull
	DefaultChunks            int              = 64
)

type ReaderOptions struct {
	Chunks       int
	Filename     string
	Aggregator   Aggregator
	Precision    Precision
	Downsampling DownsamplingMode
}

type Precision int

const (
	// Lowest level of precision. Only 1/128 samples of the full chunk are used
	Precision128 Precision = 128
	// Only 1/64 samples of the full chunk are used
	Precision64 Precision = 64
	// Only 1/32 samples of the full chunk are used
	Precision32 Precision = 32
	// Only 1/16 samples of the full chunk are used
	Precision16 Precision = 16
	// Only 1/8 samples of the full chunk are used
	Precision8 Precision = 8
	// Only 1/4 samples of the full chunk are used
	Precision4 Precision = 4
	// Only 1/2 samples of the full chunk are used
	Precision2 Precision = 2
	// Highest level of precision, samples are used as-is
	PrecisionFull Precision = 1
)

type ReaderContext struct {
	chunks             int
	filename           string
	mode               Aggregator
	file               *os.File
	decoder            *Mp3Decoder
	blocks             []float64
	chunkSize          int
	precision          Precision
	samplesPerChunk    int
	singleSampleBuffer [][2]float64
	downsampling       DownsamplingMode
}

func New(options *ReaderOptions) (*ReaderContext, error) {
	if options.Chunks == 0 {
		options.Chunks = DefaultChunks
	}
	if options.Aggregator == "" {
		options.Aggregator = DefaultTransformerMode
	}
	if options.Filename == "" {
		return nil, ErrNoFile
	}
	if options.Precision == 0 {
		options.Precision = DefaultPrecision
	}
	if options.Downsampling == "" {
		options.Downsampling = DefaultDownsamplingMode
	}

	fn, err := filepath.Abs(options.Filename)
	if err != nil {
		return nil, errors.New("failed to construct absolute path")
	}
	f, err := os.Open(fn)
	if err != nil {
		return nil, errors.New("failed to open file")
	}

	d, err := newDecoder(f)
	if err != nil {
		return nil, err
	}

	chunkSize := int(d.length()) / options.Chunks
	blocks := make([]float64, options.Chunks)
	samplesPerChunk := (chunkSize / DefaultGoMp3FrameWidth) / int(options.Precision)
	singleSampleBuffer := make([][2]float64, 1)
	ctx := &ReaderContext{
		chunks:             options.Chunks,
		filename:           options.Filename,
		mode:               options.Aggregator,
		file:               f,
		decoder:            d,
		blocks:             blocks,
		chunkSize:          chunkSize,
		precision:          options.Precision,
		samplesPerChunk:    samplesPerChunk,
		singleSampleBuffer: singleSampleBuffer,
		downsampling:       options.Downsampling,
	}

	err = ctx.process()
	if err != nil {
		return nil, err
	}

	return ctx, nil
}

func (r *ReaderContext) Close() {
	defer r.file.Close()
}

func (r *ReaderContext) Blocks() []float64 {
	return r.blocks
}

func (r *ReaderContext) process() error {
	log.Debug().
		Int("chunks", r.chunks).
		Int("raw samples per chunk", r.chunkSize).
		Int("samples per resampled chunk", r.samplesPerChunk).
		Int("total samples", int(r.decoder.length())).
		Send()

	blockBuffer := make([][2]float64, r.samplesPerChunk)
	for i := range r.blocks {
		var err error
		switch r.downsampling {
		case DownsamplingHead:
			_, err = r.downsampleHead(blockBuffer)
		case DownsamplingCenter:
			_, err = r.downsampleCenter(blockBuffer, i)
		case DownsamplingTail:
			_, err = r.downsampleTail(blockBuffer)
		case DownsamplingNone:
			_, err = r.decoder.read(blockBuffer)
		default:
			return fmt.Errorf("downsampling mode %s is not supported", r.downsampling)
		}
		if err != nil {
			return err
		}

		monoSignal := toMono(blockBuffer)

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

func (r *ReaderContext) downsampleHead(block [][2]float64) (int, error) {
	n, err := r.decoder.read(block)
	if err != nil {
		return n, err
	}
	seekSize := r.chunkSize - (n * DefaultGoMp3FrameWidth)
	_, err = r.decoder.seek(int64(seekSize), io.SeekCurrent)
	if errors.Is(err, io.EOF) {
		return n, nil
	}
	return n, err
}

func (r *ReaderContext) downsampleTail(block [][2]float64) (int, error) {
	n := len(block)
	seekSize := r.chunkSize - (n * DefaultGoMp3FrameWidth)
	sb, err := r.decoder.seek(int64(seekSize), io.SeekCurrent)
	if errors.Is(err, io.EOF) {
		return int(sb) / r.decoder.width, nil
	}
	if err != nil {
		return 0, err
	}
	rb, err := r.decoder.read(block)
	if errors.Is(err, io.EOF) {
		return rb, nil
	}
	if err != nil {
		return 0, err
	}
	return r.samplesPerChunk, nil
}

func (r *ReaderContext) downsampleCenter(block [][2]float64, chunk int) (int, error) {
	n := r.samplesPerChunk * r.decoder.width
	lq := (r.chunkSize / 2) - (n / 2)
	seekTo := (int64(r.chunkSize*(chunk) + lq))
	sb, err := r.decoder.seek(seekTo, io.SeekStart)
	if errors.Is(err, io.EOF) {
		return int(sb) / r.decoder.width, nil
	}
	if err != nil {
		return 0, err
	}
	rb, err := r.decoder.read(block)
	if errors.Is(err, io.EOF) {
		return rb, nil
	}
	if err != nil {
		return 0, err
	}
	seekEnd := int64((chunk + 1) * r.chunkSize)
	sb, err = r.decoder.seek(seekEnd, io.SeekStart)
	if errors.Is(err, io.EOF) {
		return int(sb) / r.decoder.width, nil
	}
	if err != nil {
		return 0, err
	}
	return r.samplesPerChunk, nil
}

// func (r *ReaderContext) downsampleEvenly(block [][2]float64) (n int, err error) {
// 	t := time.Now()
// 	for i := 0; i < r.samplesPerChunk; i++ {
// 		_, err := r.decoder.read(r.singleSampleBuffer)
// 		if err != nil {
// 			return 0, err
// 		}
// 		block[i] = r.singleSampleBuffer[0]
// 		seekSize := (int(r.precision) - 1) * DefaultGoMp3FrameWidth
// 		_, err = r.decoder.seek(int64(seekSize), io.SeekCurrent)
// 		if err != nil {
// 			return 0, err
// 		}
// 	}
// 	fmt.Printf("duration: %v", time.Since(t))
// 	return len(block), nil
// }
