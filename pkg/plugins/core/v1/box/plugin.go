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

package box

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/zoomoid/waveman2/pkg/painter"
	"github.com/zoomoid/waveman2/pkg/plugin"
	"github.com/zoomoid/waveman2/pkg/utils"
)

var _ plugin.Plugin = &BoxPlugin{}

var Plugin plugin.Plugin = &BoxPlugin{
	data: newBoxData(),
}

type BoxPlugin struct {
	data    *boxData
	painter *BoxPainter
}

func (b *BoxPlugin) Group() string {
	return group
}

func (b *BoxPlugin) Name() string {
	return "box"
}

func (b *BoxPlugin) Description() string {
	return description
}

func (b *BoxPlugin) Data() interface{} {
	return b.data
}

func (b *BoxPlugin) Validate() error {
	errs := b.data.validateBoxOptions()
	errlist := utils.NewErrorList(errs)
	if errlist == nil {
		return nil
	}
	return errors.New(errlist.Error())
}

func (b *BoxPlugin) Flags(flags *pflag.FlagSet) error {
	data, ok := b.Data().(*boxData)
	if !ok {
		return errors.New("box data struct is malformed")
	}
	flags.StringVar(&data.color, "color", DefaultColor, "Fill color of each box")
	flags.StringVar(&data.alignment, "alignment", string(DefaultAlignment), "Alignment of the shapes, chose one of 'top', 'center', or 'bottom'")
	flags.Float64Var(&data.rounded, "rounded", DefaultRounded, "Rounding factor of each box. Given in pixels. See SVG <rect> rx/ry attributes for details")
	flags.Float64Var(&data.gap, "gap", DefaultGap, "Gap is the spacing left between each box. Boxes are centered horizonally, so half of gap is subtracted from the box's width")
	return nil
}

func (b *BoxPlugin) Completions(cmd *cobra.Command) {
	cmd.RegisterFlagCompletionFunc("color", cobra.NoFileCompletions)
	cmd.RegisterFlagCompletionFunc("alignment", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return Alignments, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.RegisterFlagCompletionFunc("rounded", cobra.NoFileCompletions)
	cmd.RegisterFlagCompletionFunc("gap", cobra.NoFileCompletions)
}

func (b *BoxPlugin) Draw(options *painter.PainterOptions) []string {
	painter := NewPainter(options, b.data.toOptions(options.Width, options.Height))
	b.painter = painter
	return painter.Draw()
}

func (b *BoxPlugin) Painter() painter.Painter {
	return b.painter
}

type boxData struct {
	color     string
	alignment string
	height    float64
	width     float64
	rounded   float64
	gap       float64
}

func (b *boxData) validateBoxOptions() (errList []error) {
	if err := validateAlignment(b.alignment); err != nil {
		errList = append(errList, err)
	}
	if err := validateGap(b.gap, b.width); err != nil {
		errList = append(errList, err)
	}
	if len(errList) == 0 {
		return nil
	}
	return errList
}

func (b *boxData) toOptions(width float64, height float64) *BoxOptions {
	p := &BoxOptions{
		Alignment: Alignment(b.alignment),
		Color:     b.color,
		BoxHeight: height,
		BoxWidth:  width,
		Rounded:   b.rounded,
		Gap:       b.gap,
	}
	return p
}

func newBoxData() *boxData {
	return &boxData{
		color:     DefaultColor,
		alignment: string(DefaultAlignment),
		height:    DefaultHeight,
		width:     DefaultWidth,
		rounded:   DefaultRounded,
		gap:       DefaultRounded,
	}
}
