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

package reference

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/zoomoid/waveman2/cmd/validation"
	"github.com/zoomoid/waveman2/pkg/painter"
	"github.com/zoomoid/waveman2/pkg/painter/line"
	"github.com/zoomoid/waveman2/pkg/plugin"
	options "github.com/zoomoid/waveman2/pkg/reference/options/line"
	"github.com/zoomoid/waveman2/pkg/utils"
)

var _ plugin.Plugin = &LinePainter{}

var LinePainterPlugin plugin.Plugin = &LinePainter{
	data:    newLineData(),
	enabled: false,
}

const (
	LineMode        string = "line"
	LineDescription string = "Create a line waveform"
)

type LinePainter struct {
	data    *lineData
	enabled bool
	painter *line.LinePainter
}

func (l *LinePainter) Name() string {
	return LineMode
}

func (l *LinePainter) Description() string {
	return LineDescription
}

func (l *LinePainter) Enabled() *bool {
	return &l.enabled
}

func (l *LinePainter) Data() interface{} {
	return l.data
}

func (l *LinePainter) Validate() error {
	errs := l.data.validateLineOptions()
	errlist := utils.NewErrorList(errs)
	if errlist == nil {
		return nil
	}
	return errors.New(errlist.Error())
}

func (l *LinePainter) Flags(flags *pflag.FlagSet) error {
	data, ok := l.Data().(*lineData)
	if !ok {
		return errors.New("line data struct is malformed")
	}
	flags.BoolVar(l.Enabled(), LineMode, false, LineDescription)
	flags.StringVar(&data.interpolation, options.Interpolation, string(line.DefaultInterpolation), options.InterpolationDescription)
	flags.StringVar(&data.fill, options.LineFill, line.DefaultFillColor, options.LineFillDescription)
	flags.StringVar(&data.strokeColor, options.StrokeColor, line.DefaultStrokeColor, options.StrokeColorDescription)
	flags.Float64Var(&data.strokeWidth, options.StrokeWidth, line.DefaultStrokeWidth, options.StrokeWidthDescription)
	flags.BoolVarP(&data.closed, options.Closed, options.ClosedShort, false, options.ClosedDescription)
	flags.BoolVarP(&data.inverted, options.Inverted, options.InvertedShort, false, options.InvertedDescription)
	return nil
}

func (l *LinePainter) Completions(cmd *cobra.Command) {
	cmd.RegisterFlagCompletionFunc(options.Interpolation, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return line.Interpolations, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.RegisterFlagCompletionFunc(options.LineFill, cobra.NoFileCompletions)
	cmd.RegisterFlagCompletionFunc(options.StrokeColor, cobra.NoFileCompletions)
	cmd.RegisterFlagCompletionFunc(options.StrokeWidth, cobra.NoFileCompletions)
	cmd.RegisterFlagCompletionFunc(options.Closed, cobra.NoFileCompletions)
	cmd.RegisterFlagCompletionFunc(options.Inverted, cobra.NoFileCompletions)
}

func (l *LinePainter) Draw(options *painter.PainterOptions) []string {
	painter := line.New(options, l.data.toOptions(options.Width, options.Height))
	l.painter = painter
	return painter.Draw()
}

func (l *LinePainter) Painter() painter.Painter {
	return l.painter
}

func (l *lineData) toOptions(width float64, height float64) *line.LineOptions {
	return &line.LineOptions{
		Interpolation: line.Interpolation(l.interpolation),
		Fill:          l.fill,
		Stroke: &line.Stroke{
			Color: l.strokeColor,
			Width: l.strokeWidth,
		},
		Closed:    l.closed,
		Spread:    width,
		Amplitude: height,
		Inverted:  l.inverted,
	}
}

func (l *lineData) validateLineOptions() (errList []error) {
	if err := validation.ValidateInterpolation(l.interpolation); err != nil {
		errList = append(errList, err)
	}
	if len(errList) == 0 {
		return nil
	}
	return errList
}

type lineData struct {
	interpolation string
	fill          string
	strokeColor   string
	strokeWidth   float64
	spread        float64
	height        float64
	closed        bool
	inverted      bool
}

func newLineData() *lineData {
	return &lineData{
		interpolation: string(line.DefaultInterpolation),
		fill:          line.DefaultFillColor,
		strokeColor:   line.DefaultStrokeColor,
		strokeWidth:   line.DefaultStrokeWidth,
		spread:        line.DefaultSpread,
		height:        line.DefaultHeight,
		closed:        false,
		inverted:      false,
	}
}
