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
	"github.com/zoomoid/waveman2/pkg/painter/box"
	"github.com/zoomoid/waveman2/pkg/plugin"
	options "github.com/zoomoid/waveman2/pkg/reference/options/box"
	"github.com/zoomoid/waveman2/pkg/utils"
)

var _ plugin.Plugin = &BoxPainter{}

var BoxPainterPlugin plugin.Plugin = &BoxPainter{
	data:    newBoxData(),
	enabled: false,
}

const (
	BoxMode        string = "box"
	BoxDescription string = "Create a box waveform"
)

type BoxPainter struct {
	data    *boxData
	enabled bool
	painter *box.BoxPainter
}

func (b *BoxPainter) Name() string {
	return BoxMode
}

func (b *BoxPainter) Description() string {
	return BoxDescription
}

func (b *BoxPainter) Enabled() *bool {
	return &b.enabled
}

func (b *BoxPainter) Data() interface{} {
	return b.data
}

func (b *BoxPainter) Validate() error {
	errs := b.data.validateBoxOptions()
	errlist := utils.NewErrorList(errs)
	if errlist == nil {
		return nil
	}
	return errors.New(errlist.Error())
}

func (b *BoxPainter) Flags(flags *pflag.FlagSet) error {
	data, ok := b.Data().(*boxData)
	if !ok {
		return errors.New("box data struct is malformed")
	}
	flags.BoolVar(b.Enabled(), BoxMode, false, BoxDescription)
	flags.StringVar(&data.color, options.BoxFill, box.DefaultColor, options.BoxFillDescription)
	flags.StringVar(&data.alignment, options.Alignment, string(box.DefaultAlignment), options.AlignmentDescription)
	flags.Float64Var(&data.rounded, options.BoxRounded, box.DefaultRounded, options.BoxRoundedDescription)
	flags.Float64Var(&data.gap, options.BoxGap, box.DefaultGap, options.BoxGapDescription)
	return nil
}

func (b *BoxPainter) Completions(cmd *cobra.Command) {
	cmd.RegisterFlagCompletionFunc(options.BoxFill, cobra.NoFileCompletions)
	cmd.RegisterFlagCompletionFunc(options.Alignment, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return box.Alignments, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.RegisterFlagCompletionFunc(options.BoxRounded, cobra.NoFileCompletions)
	cmd.RegisterFlagCompletionFunc(options.BoxGap, cobra.NoFileCompletions)
}

func (b *BoxPainter) Draw(options *painter.PainterOptions) []string {
	painter := box.New(options, b.data.toOptions(options.Width, options.Height))
	b.painter = painter
	return painter.Draw()
}

func (b *BoxPainter) Painter() painter.Painter {
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
	if err := validation.ValidateAlignment(b.alignment); err != nil {
		errList = append(errList, err)
	}
	if err := validation.ValidateGap(b.gap, b.width); err != nil {
		errList = append(errList, err)
	}
	if len(errList) == 0 {
		return nil
	}
	return errList
}

func (b *boxData) toOptions(width float64, height float64) *box.BoxOptions {
	p := &box.BoxOptions{
		Alignment: box.Alignment(b.alignment),
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
		color:     box.DefaultColor,
		alignment: string(box.DefaultAlignment),
		height:    box.DefaultHeight,
		width:     box.DefaultWidth,
		rounded:   box.DefaultRounded,
		gap:       box.DefaultRounded,
	}
}
