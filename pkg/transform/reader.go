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
	"fmt"
	"io"
	"strings"
)

type DownsamplingMode string

const (
	DownsamplingNone   DownsamplingMode = "none"
	DownsamplingHead   DownsamplingMode = "head"
	DownsamplingCenter DownsamplingMode = "center"
	DownsamplingTail   DownsamplingMode = "tail"
	DownsamplingEmpty  DownsamplingMode = ""
)

var DownsamplingModes = []string{"center", "head", "tail", "none"}

type Aggregator string

const (
	AggregatorAverage        Aggregator = "avg"
	AggregatorRoundedAverage Aggregator = "rounded-avg"
	AggregatorMax            Aggregator = "max"
	AggregatorMeanSquare     Aggregator = "mean-square"
	AggregatorRootMeanSquare Aggregator = "rms"
	AggregatorEmpty          Aggregator = ""
)

type WindowAlgorithm int

var WindowAlgorithms = []string{Rectangular.String(), Hann.String(), Tukey.String(), PlanckTaper.String()}

const (
	Rectangular WindowAlgorithm = iota
	Hann
	Tukey
	PlanckTaper
)

func (a WindowAlgorithm) String() string {
	switch a {
	case Rectangular:
		return "rectangular"
	case Hann:
		return "hann"
	case Tukey:
		return "tukey"
	case PlanckTaper:
		return "plank-taper"
	default:
		return ""
	}
}

func WindowAlgorithmFromString(a string) WindowAlgorithm {
	a = strings.ToLower(a)
	switch a {
	case "rect", "rectangle", "square":
		return Rectangular
	case "hann":
		return Hann
	case "tukey":
		return Tukey
	case "plank-taper", "plank_taper", "planktaper":
		return PlanckTaper
	default:
		return Rectangular
	}
}

var Aggregators = []string{"rms", "mean-square", "rounded-avg", "avg", "max"}

var (
	ErrNoFile error = errors.New("no file given")

	DefaultAggregator        Aggregator       = AggregatorRootMeanSquare
	DefaultRoundingPrecision uint             = 3
	DefaultDownsamplingMode  DownsamplingMode = DownsamplingCenter
	DefaultPrecision         Precision        = PrecisionFull
	DefaultChunks            int              = 64
	DefaultWindowParameter   float64          = 0
	DefaultWindowAlgorithm   WindowAlgorithm  = Rectangular
	DefaultClamping          *Clamping        = &Clamping{
		Min: 0,
		Max: 1,
	}
	DefaultWindow *Window = &Window{
		Algorithm: DefaultWindowAlgorithm,
		P:         0,
	}
)

type ReaderOptions struct {
	Chunks       int
	Aggregator   Aggregator
	Precision    Precision
	Downsampling DownsamplingMode

	Normalize bool
	Window    *Window
	Clamping  *Clamping
}

type Window struct {
	Algorithm WindowAlgorithm
	P         float64
}

type Clamping struct {
	Min float64
	Max float64
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

// DownsamplingModes contains all currently supported levels as integers for Cobra flag autocompletion
var DownsamplingPrecisions = []string{"1", "2", "4", "8", "16", "32", "64", "128"}

const (
	// Maximum precision alias for range checking
	MaximumPrecision = PrecisionFull
	// Minimum precision alias for range checking
	MinimumPrecision = Precision128
)

type Transformer ReaderContext

type ReaderContext struct {
	chunks             int
	mode               Aggregator
	reader             io.Reader
	decoder            *Mp3Decoder
	blocks             []float64
	chunkSize          int
	precision          Precision
	samplesPerChunk    int
	singleSampleBuffer [][2]float64
	downsampling       DownsamplingMode
	clipping           *Clamping
	windowParam        float64
	windowAlgo         WindowAlgorithm
	normalize          bool
}

func New(options *ReaderOptions, reader io.Reader) (*ReaderContext, error) {
	if options.Chunks == 0 {
		options.Chunks = DefaultChunks
	}
	if options.Aggregator == AggregatorEmpty {
		options.Aggregator = DefaultAggregator
	}
	if options.Clamping == nil {
		options.Clamping = DefaultClamping
	}
	if options.Window == nil {
		options.Window = DefaultWindow
	}
	if options.Precision == 0 {
		options.Precision = DefaultPrecision
	}
	if options.Downsampling == DownsamplingEmpty {
		options.Downsampling = DefaultDownsamplingMode
	}

	d, err := newDecoder(reader)
	if err != nil {
		return nil, err
	}

	chunkSize := int(d.length()) / options.Chunks
	blocks := make([]float64, options.Chunks)
	samplesPerChunk := (chunkSize / DefaultGoMp3FrameWidth) / int(options.Precision)
	singleSampleBuffer := make([][2]float64, 1)

	ctx := &ReaderContext{
		chunks:             options.Chunks,
		mode:               options.Aggregator,
		reader:             reader,
		decoder:            d,
		blocks:             blocks,
		chunkSize:          chunkSize,
		precision:          options.Precision,
		samplesPerChunk:    samplesPerChunk,
		singleSampleBuffer: singleSampleBuffer,
		downsampling:       options.Downsampling,
		windowParam:        options.Window.P,
		windowAlgo:         options.Window.Algorithm,
		clipping:           options.Clamping,
		normalize:          options.Normalize,
	}

	err = ctx.process()
	if err != nil {
		return nil, err
	}

	return ctx, nil
}

func (r *ReaderContext) Close() {
	// defer r.reader.Close()
}

func (r *ReaderContext) Blocks() []float64 {
	return r.blocks
}

func (r *ReaderContext) process() error {
	// log.Debug().
	// 	Int("chunks", r.chunks).
	// 	Int("raw samples per chunk", r.chunkSize).
	// 	Int("samples per resampled chunk", r.samplesPerChunk).
	// 	Int("total samples", int(r.decoder.length())).
	// 	Send()

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
		case AggregatorMax:
			block = max(monoSignal)
		case AggregatorAverage:
			block = mean(monoSignal)
		case AggregatorRoundedAverage:
			block = roundedMean(monoSignal, DefaultRoundingPrecision)
		case AggregatorMeanSquare:
			block = meanSquare(monoSignal)
		case AggregatorRootMeanSquare:
			block = rootMeanSquare(monoSignal)
		default:
			return fmt.Errorf("mode %s is not implemented", r.mode)
		}
		r.blocks[i] = block
	}

	// last step is to normalize the block range to [0,1]
	if r.normalize {
		r.blocks = normalize(r.blocks)
	}
	for idx, sample := range r.blocks {
		r.blocks[idx] = clamp(sample, r.clipping.Min, r.clipping.Max)
	}

	switch r.windowAlgo {
	case Hann:
		r.blocks = hann(r.blocks, r.windowParam)
	case Tukey:
		r.blocks = tukey(r.blocks, r.windowParam)
	case PlanckTaper:
		r.blocks = planck_taper(r.blocks, r.windowParam)
	case Rectangular:
		// rectangular window over the entire range equals no window
	default:
	}

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
