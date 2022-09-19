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

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/zoomoid/waveman2/cmd/options"
	"github.com/zoomoid/waveman2/pkg/painter"
)

type filenameOptions struct {
	filenames []string
	recursive bool
	output    string
}

type sharedPainterOptions struct {
	height float64
	width  float64
}

func addDimensionFlags(flags *pflag.FlagSet, data *sharedPainterOptions) {
	flags.Float64VarP(&data.width, options.Width, options.WidthShort, painter.DefaultWidth, options.WidthDescription)
	flags.Float64VarP(&data.height, options.Height, options.HeightShort, painter.DefaultHeight, options.HeightDescription)
}

func addIOFlags(flags *pflag.FlagSet, data *filenameOptions) {
	flags.StringSliceVarP(&data.filenames, options.Filename, options.FilenameShort, nil, options.FilenameDescription)
	flags.BoolVarP(&data.recursive, options.Recursive, options.RecursiveShort, false, options.RecursiveDescription)
	flags.StringVarP(&data.output, options.Output, options.OutputShort, "", options.OutputDescription)
}

func addDimensionFlagsCompletion(cmd *cobra.Command) {
	cmd.RegisterFlagCompletionFunc(options.Width, cobra.NoFileCompletions)
	cmd.RegisterFlagCompletionFunc(options.Height, cobra.NoFileCompletions)
}

func addIOFlagsCompletion(cmd *cobra.Command) {
	cmd.RegisterFlagCompletionFunc(options.Filename, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"mp3"}, cobra.ShellCompDirectiveFilterFileExt
	})
	cmd.RegisterFlagCompletionFunc(options.Output, cobra.NoFileCompletions)
}
