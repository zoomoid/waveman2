package cmd

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"github.com/zoomoid/waveman/v2/cmd/options"
	cmdutils "github.com/zoomoid/waveman/v2/cmd/utils"
	"github.com/zoomoid/waveman/v2/cmd/validation"
	"github.com/zoomoid/waveman/v2/pkg/painter/box"
)

var _ Plugin = &BoxPainter{}

var BoxPainterPlugin Plugin = &BoxPainter{
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
	return *b.data
}

func (b *BoxPainter) Validate() error {
	errs := b.data.validateBoxOptions()
	errlist := cmdutils.NewErrorList(errs)
	if errlist == nil {
		return nil
	}
	return errors.New(errlist.Error())
}

func (b *BoxPainter) Flags(flags *pflag.FlagSet) {
	data, ok := b.Data().(boxData)
	if !ok {
		log.Fatal().Msg("box data struct is malformed")
	}
	flags.BoolVar(b.Enabled(), BoxMode, false, BoxDescription)
	flags.StringVar(&data.color, options.BoxFill, box.DefaultColor, options.BoxFillDescription)
	flags.StringVar(&data.alignment, options.Alignment, string(box.DefaultAlignment), options.AlignmentDescription)
	flags.Float64VarP(&data.rounded, options.BoxRounded, options.BoxRoundedShort, box.DefaultRounded, options.BoxRoundedDescription)
	flags.Float64Var(&data.gap, options.BoxGap, box.DefaultGap, options.BoxGapDescription)
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

type boxData struct {
	color     string
	alignment string
	height    float64
	width     float64
	rounded   float64
	gap       float64
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
