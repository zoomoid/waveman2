package validation

import (
	"fmt"

	"github.com/zoomoid/waveman/v2/pkg/transform"
)

func ValidateDownsamplingMode(mode string) error {
	m := transform.DownsamplingMode(mode)
	switch m {
	case transform.DownsamplingCenter:
	case transform.DownsamplingHead:
	case transform.DownsamplingTail:
	case transform.DownsamplingNone:
	case transform.DownsamplingEmpty:
		return nil
	}
	return fmt.Errorf("downsampling mode %s is not supported", mode)
}

func ValidateDownsamplingFactor(factor int) error {
	f := transform.Precision(factor)
	switch f {
	case transform.Precision128:
	case transform.Precision64:
	case transform.Precision32:
	case transform.Precision16:
	case transform.Precision8:
	case transform.Precision4:
	case transform.Precision2:
	case transform.PrecisionFull:
		return nil
	}
	return fmt.Errorf("downsampling factor %d is not supported", factor)
}

func ValidateChunks(chunks int) error {
	if chunks > 0 {
		return nil
	}
	return fmt.Errorf("downsampling factor must be strictly positve")
}

func ValidateAggregator(aggregator string) error {
	a := transform.Aggregator(aggregator)
	switch a {
	case transform.AggregatorAverage:
	case transform.AggregatorMax:
	case transform.AggregatorMeanSquare:
	case transform.AggregatorRootMeanSquare:
	case transform.AggregatorRoundedAverage:
	case transform.AggregatorEmpty:
		return nil
	}
	return fmt.Errorf("aggregator %s is not supported", aggregator)
}
