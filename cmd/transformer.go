package cmd

import (
	"github.com/spf13/pflag"
	"github.com/zoomoid/waveman/v2/cmd/options"
	cmdutils "github.com/zoomoid/waveman/v2/cmd/utils"
	"github.com/zoomoid/waveman/v2/cmd/validation"
	"github.com/zoomoid/waveman/v2/pkg/transform"
)

// transformerData captures all properties defineable by flags
// at calling the command
type transformerData struct {
	downsamplingMode   string
	downsamplingFactor int
	aggregator         string
	chunks             int
}

func newTransformerData() *transformerData {
	return &transformerData{
		downsamplingMode:   string(transform.DefaultDownsamplingMode),
		downsamplingFactor: int(transform.DefaultPrecision),
		aggregator:         string(transform.DefaultAggregator),
		chunks:             transform.DefaultChunks,
	}
}

func addTranformerFlags(flags *pflag.FlagSet, data *transformerData) {
	flags.StringVar(&data.downsamplingMode, options.DownsamplingMode, "", options.DownsamplingModeDescription)
	flags.IntVar(&data.downsamplingFactor, options.DownsamplingFactor, 1, options.DownsamplingFactorDescription)
	flags.StringVar(&data.aggregator, options.Aggregator, string(transform.DefaultAggregator), options.AggregatorDescription)
	flags.IntVarP(&data.chunks, options.Chunks, options.ChunksShort, transform.DefaultChunks, options.ChunksDescription)
}

func (t *transformerData) validateTransformerOptions() cmdutils.ErrorList {
	errList := []error{}
	if err := validation.ValidateDownsamplingFactor(t.downsamplingFactor); err != nil {
		errList = append(errList, err)
	}
	if err := validation.ValidateDownsamplingMode(t.downsamplingMode); err != nil {
		errList = append(errList, err)
	}
	if err := validation.ValidateChunks(t.chunks); err != nil {
		errList = append(errList, err)
	}
	if err := validation.ValidateAggregator(t.aggregator); err != nil {
		errList = append(errList, err)
	}
	return cmdutils.NewErrorList(errList)
}

func (t *transformerData) toOptions() *transform.ReaderOptions {
	return &transform.ReaderOptions{
		Chunks:       t.chunks,
		Aggregator:   transform.Aggregator(t.aggregator),
		Precision:    transform.Precision(t.downsamplingFactor),
		Downsampling: transform.DownsamplingMode(t.downsamplingMode),
	}
}
