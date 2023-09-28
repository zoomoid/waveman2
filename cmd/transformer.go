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

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/zoomoid/waveman2/cmd/options"
	"github.com/zoomoid/waveman2/cmd/validation"
	"github.com/zoomoid/waveman2/pkg/transform"
	"github.com/zoomoid/waveman2/pkg/utils"
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

func addTransformerFlagCompletion(cmd *cobra.Command) {
	cmd.RegisterFlagCompletionFunc(options.DownsamplingMode, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return transform.DownsamplingModes, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.RegisterFlagCompletionFunc(options.DownsamplingFactor, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return transform.DownsamplingPrecisions, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.RegisterFlagCompletionFunc(options.Aggregator, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return transform.Aggregators, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.RegisterFlagCompletionFunc(options.Chunks, cobra.NoFileCompletions)
}

func (t *transformerData) validateTransformerOptions() utils.ErrorList {
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
	return utils.NewErrorList(errList)
}

func (t *transformerData) toOptions() *transform.ReaderOptions {
	return &transform.ReaderOptions{
		Chunks:       t.chunks,
		Aggregator:   transform.Aggregator(t.aggregator),
		Precision:    transform.Precision(t.downsamplingFactor),
		Downsampling: transform.DownsamplingMode(t.downsamplingMode),
	}
}
