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

package line

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/zoomoid/waveman2/pkg/painter"
	"github.com/zoomoid/waveman2/pkg/plugin"
	"github.com/zoomoid/waveman2/pkg/utils"
)

var _ plugin.Plugin = &LinePlugin{}

var Plugin plugin.Plugin = &LinePlugin{
	data: newLineData(),
}

type LinePlugin struct {
	data    *lineData
	painter *LinePainter
}

func (l *LinePlugin) Group() string {
	return group
}

func (l *LinePlugin) Name() string {
	return "line"
}

func (l *LinePlugin) Description() string {
	return description
}

func (l *LinePlugin) Data() interface{} {
	return l.data
}

func (l *LinePlugin) Validate() error {
	errs := l.data.validateLineOptions()
	errlist := utils.NewErrorList(errs)
	if errlist == nil {
		return nil
	}
	return errors.New(errlist.Error())
}

func (l *LinePlugin) Flags(flags *pflag.FlagSet) error {
	data, ok := l.Data().(*lineData)
	if !ok {
		return errors.New("line data struct is malformed")
	}
	flags.StringVar(&data.interpolation, "interpolation", string(DefaultInterpolation), "Interpolation mechanism to be used for smoothing the curve [none,fritsch-carlson,steffen]")
	flags.StringVar(&data.fill, "fill-color", DefaultFillColor, "Color for the area enclosed by the line")
	flags.StringVar(&data.strokeColor, "stroke-color", DefaultStrokeColor, "Color of the line's stroke")
	flags.Float64Var(&data.strokeWidth, "stroke-width", DefaultStrokeWidth, "Width of the line's stroke")
	flags.BoolVarP(&data.closed, "closed", "c", false, "Whether the SVG path should be closed or left open")
	flags.BoolVarP(&data.inverted, "inverted", "i", false, "Whether the shape should be inverted horizontally, i.e., switch the vertical alignment from top to bottom")
	return nil
}

func (l *LinePlugin) Completions(cmd *cobra.Command) {
	cmd.RegisterFlagCompletionFunc("interpolation", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return Interpolations, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.RegisterFlagCompletionFunc("fill-color", cobra.NoFileCompletions)
	cmd.RegisterFlagCompletionFunc("stroke-color", cobra.NoFileCompletions)
	cmd.RegisterFlagCompletionFunc("stroke-width", cobra.NoFileCompletions)
	cmd.RegisterFlagCompletionFunc("closed", cobra.NoFileCompletions)
	cmd.RegisterFlagCompletionFunc("inverted", cobra.NoFileCompletions)
}

func (l *LinePlugin) Draw(options *painter.PainterOptions) []string {
	painter := NewPainter(options, l.data.toOptions(options.Width, options.Height))
	l.painter = painter
	return painter.Draw()
}

func (l *LinePlugin) Painter() painter.Painter {
	return l.painter
}

func (l *lineData) toOptions(width float64, height float64) *LineOptions {
	return &LineOptions{
		Interpolation: Interpolation(l.interpolation),
		Fill:          l.fill,
		Stroke: &Stroke{
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
	if err := validateInterpolation(l.interpolation); err != nil {
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
	mirrored      bool
}

func newLineData() *lineData {
	return &lineData{
		interpolation: string(DefaultInterpolation),
		fill:          DefaultFillColor,
		strokeColor:   DefaultStrokeColor,
		strokeWidth:   DefaultStrokeWidth,
		spread:        DefaultSpread,
		height:        DefaultHeight,
		closed:        false,
		inverted:      false,
	}
}
