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

package sweep

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/zoomoid/waveman2/pkg/painter"
	"github.com/zoomoid/waveman2/pkg/plugin"
	"github.com/zoomoid/waveman2/pkg/utils"
)

var _ plugin.Plugin = &SweepPlugin{}

var Plugin plugin.Plugin = &SweepPlugin{
	data: newSweepData(),
}

type SweepPlugin struct {
	data    *sweepData
	painter *SweepPainter
}

func (l *SweepPlugin) Group() string {
	return group
}

func (l *SweepPlugin) Name() string {
	return "sweep"
}

func (l *SweepPlugin) Description() string {
	return description
}

func (l *SweepPlugin) Data() interface{} {
	return l.data
}

func (l *SweepPlugin) Validate() error {
	errs := l.data.validateLineOptions()
	errlist := utils.NewErrorList(errs)
	if errlist == nil {
		return nil
	}
	return errors.New(errlist.Error())
}

func (l *SweepPlugin) Flags(flags *pflag.FlagSet) error {
	data, ok := l.Data().(*sweepData)
	if !ok {
		return errors.New("sweep data struct is malformed")
	}
	flags.StringVar(&data.interpolation, "interpolation", string(DefaultInterpolation), "Interpolation mechanism to be used for smoothing the curve [none,fritsch-carlson,steffena,akima]")
	flags.StringVar(&data.fill, "fill-color", DefaultFillColor, "Color for the area enclosed by the line")
	flags.StringVar(&data.strokeColor, "stroke-color", DefaultStrokeColor, "Color of the line's stroke")
	flags.Float64Var(&data.strokeWidth, "stroke-width", DefaultStrokeWidth, "Width of the line's stroke")
	return nil
}

func (l *SweepPlugin) Completions(cmd *cobra.Command) {
	cmd.RegisterFlagCompletionFunc("interpolation", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return Interpolations, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.RegisterFlagCompletionFunc("fill-color", cobra.NoFileCompletions)
	cmd.RegisterFlagCompletionFunc("stroke-color", cobra.NoFileCompletions)
	cmd.RegisterFlagCompletionFunc("stroke-width", cobra.NoFileCompletions)
	cmd.RegisterFlagCompletionFunc("closed", cobra.NoFileCompletions)
	cmd.RegisterFlagCompletionFunc("inverted", cobra.NoFileCompletions)
}

func (l *SweepPlugin) Draw(options *painter.PainterOptions) []string {
	painter := NewPainter(options, l.data.toOptions(options.Width, options.Height))
	l.painter = painter
	return painter.Draw()
}

func (l *SweepPlugin) Painter() painter.Painter {
	return l.painter
}

func (l *sweepData) toOptions(width float64, height float64) *LineOptions {
	return &LineOptions{
		Interpolation: Interpolation(l.interpolation),
		Fill:          l.fill,
		Stroke: &Stroke{
			Color: l.strokeColor,
			Width: l.strokeWidth,
		},
		Spread:    width,
		Amplitude: height,
	}
}

func (l *sweepData) validateLineOptions() (errList []error) {
	if err := validateInterpolation(l.interpolation); err != nil {
		errList = append(errList, err)
	}
	if len(errList) == 0 {
		return nil
	}
	return errList
}

type sweepData struct {
	interpolation string
	fill          string
	strokeColor   string
	strokeWidth   float64
	spread        float64
	height        float64
}

func newSweepData() *sweepData {
	return &sweepData{
		interpolation: string(DefaultInterpolation),
		fill:          DefaultFillColor,
		strokeColor:   DefaultStrokeColor,
		strokeWidth:   DefaultStrokeWidth,
		spread:        DefaultSpread,
		height:        DefaultHeight,
	}
}
