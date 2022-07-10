package cmd

import "github.com/zoomoid/waveman/v2/pkg/transform"

// transfomerData captures all properties defineable by flags
// at calling the command
type transfomerData struct {
	downsamplingMode   string
	downsamplingFactor int
	aggregator         string
	filename           string
	chunks             int
	output             string
}

func newTransformerData() *transfomerData {
	return &transfomerData{
		downsamplingMode:   string(transform.DefaultDownsamplingMode),
		downsamplingFactor: int(transform.DefaultPrecision),
		aggregator:         string(transform.DefaultAggregator),
		filename:           "",
		chunks:             transform.DefaultChunks,
	}
}
