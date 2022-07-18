package cmd

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/zoomoid/waveman/v2/cmd/options"
	"github.com/zoomoid/waveman/v2/pkg/painter"
)

type filenameOptions struct {
	filenames []string
	recursive bool
	output    string
}

func (o *filenameOptions) validate(useStdout bool) []error {
	var errors []error
	if !useStdout && len(o.filenames) > 1 {
		errors = append(errors, fmt.Errorf("cannot use multiple files with stdout target, use --output file"))
	}
	return errors
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
