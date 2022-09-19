/*
Copyright 2022 zoomoid.

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

package validation

import (
	"fmt"

	"github.com/zoomoid/waveman2/pkg/transform"
)

func ValidateDownsamplingMode(mode string) error {
	m := transform.DownsamplingMode(mode)
	switch m {
	case transform.DownsamplingCenter,
		transform.DownsamplingHead,
		transform.DownsamplingTail,
		transform.DownsamplingNone,
		transform.DownsamplingEmpty:
		return nil
	}
	return fmt.Errorf("downsampling mode %s is not supported", mode)
}

func ValidateDownsamplingFactor(factor int) error {
	f := transform.Precision(factor)
	switch f {
	case transform.Precision128,
		transform.Precision64,
		transform.Precision32,
		transform.Precision16,
		transform.Precision8,
		transform.Precision4,
		transform.Precision2,
		transform.PrecisionFull:
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
	case transform.AggregatorAverage,
		transform.AggregatorMax,
		transform.AggregatorMeanSquare,
		transform.AggregatorRootMeanSquare,
		transform.AggregatorRoundedAverage,
		transform.AggregatorEmpty:
		return nil
	}
	return fmt.Errorf("aggregator %s is not supported", aggregator)
}
