package cmd

import (
	"fmt"

	"github.com/lithammer/dedent"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/zoomoid/waveman/v2/cmd/options"
	cmdutils "github.com/zoomoid/waveman/v2/cmd/utils"
	"github.com/zoomoid/waveman/v2/cmd/validation"
	"github.com/zoomoid/waveman/v2/pkg/painter"
)

var (
	WavemanShort string = "waveman generates stylized visual waveforms from mp3 files. Comes with a box painter and a line painter, but can be extended to with other painters easily."

	WavemanLong string = dedent.Dedent(`
		Generate SVG waveforms for one or more mp3 files.

		Prints SVG to stdout when not --output is not specified. When passing in a
		directory, will create SVG files named by the mp3 source files. When the
		--recursive flag is used, *all* mp3 files below the path are used and SVG files
		are colocated with the source mp3 files.
		
		You can configure the sample decoder/transformer in various ways: The number of
		chunks to be passed down to the painter can be set with --chunks (or -n). The
		number must be non-negative. The aggregation function by default uses
		Root-Mean-Square ("rms") for the samples in each chunk. This mimmicks the way
		metering in most DAWs would. Instead, you can also choose "avg", "max",
		"mean-square", or "rounded-avg". The last mode is particularly nice if you don't
		like large floating point numbers in you SVG code, rounding to 3 digits by
		default.
		
		You can improve performance of the waveman by aggressively downsampling the
		audio file. We tested this out and found that using full resolution for the
		aggregation of samples yields minimum visual changes to the audio file, compared
		to the use of high downsampling ratios. We downsample evenly, reducing the
		window of samples per chunks in powers of two. This means that the downsampling
		factor is given as 1/2, 1/4, 1/8, 1/16, etc., up to 1/128.

		To set this factor with flags, use the inverse in --downsampling-factor, e.g.
		"--downsampling-factor 16" for a downsampling factor of 1/16. 
		
		Due to I/O bottlenecks, *this is not done evenly* throughout the file. Instead,
		the downsampling window of a chunk is located either at the start, the middle,
		or the end of a chunk. This behaviour can be set with --downsampling-mode, either
		"head", "center", or "tail".

		--------     ---------------     -----------                  -------
		| File | --> | Transformer | --> | Painter | --> Elements --> | SVG |
		--------     ---------------     -----------                  -------

		The reference implementation brings painters for boxes and lines. Both have 
		multiple configuration options.

		The box painter is used by setting the --box flag. The box color can be set
		with --color. The alignment axis can be either "top", "center", or "bottom",
		and set with --alignment. --height (or -h) sets the height of highest box, thus also the
		height of the entire canvas. --width (or -w) sets the width of each box's bounding box.
		--gap sets the space left between each box. Boxes are painted centered inside 
		their bounding box: 

		|-------------------------------------------|
		|<- 0.5 * gap ->|----BOX----|<- 0.5 * gap ->|
		|<----------------- width ----------------->|

		Lastly, the --rounded (or -r) parameter controls the rounding of the rectangles.
		Notably, rounding requires the boxes to have a minimum height, namely at least
		the width of the box, to look aesthetically pleasing. When using --rounded,
		each box's height will have its width as a lower bound.

		The line painter is used by setting the --line flag. A line's path can be closed
		by setting the --closed (or -c) flag. This will close the <path> by appending "Z" 
		at the end of the data points. The resulting shape can be horizontally mirrored by 
		setting --inverted (-i). This uses CSS transforms as linear transformation, rather 
		than computing the data points with offset.
		
		When the path is closed, the color of the enclosed shape can be set with 
		--fill-color. The color of the line is set with --stroke-color, and the width of
		the line with --stroke-width. All those require SVG/CSS-compliant values for the
		attributes.

		Similarly to the box painter, the --height (or -h) flag controls the shape's overall
		height. --spread (or -s) controls the horizontal spacing between each of the sample
		points. 

		To make the line appear smoothly from a discrete set of points, we interpolate 
		control points for each sample point using cubic hermetic interpolation to fit 
		cubic polynomials. Namely, we implement 2 interpolation schemes: Fritsch-Carlson
		and Steffen. Details can be seen here: 
		http://math.stackexchange.com/questions/45218/implementation-of-monotone-cubic-interpolation
		This way, the shape appears smooth. Interpolation can also be controlled with
		flags: By default, the Frisch-Carlson scheme is used; setting "--interpolation steffen"
		uses the Steffen scheme. If you want to disable interpolation entirely, set
		"--interpolation none".
	`)

	WavemanExamples string = dedent.Dedent(`
		# Create a black box waveform with 50 blocks for a single mp3
		waveman --box --chunks 50 -f audio.mp3

		# Create a line waveform with 32 sample points for a single mp3
		waveman --line --chunks 32 -f audio.mp3

		# Create a red box waveform with 50 blocks at 1/8 downsampling factor
		waveman --box --chunks 50 --fill-color red --downsampling-factor 8 -f audio.mp3

		# Create a green box waveform for *all* mp3 files in the directory at 
		# 1/4 downsampling
		waveman --box --fill-color green --downsampling-factor 4 -f ./

		# Create a closed line waveform with 128 sample points and 1/64 downsampling 
		# from the start of each chunk, spread apart 50 pixels, with a thicker yellow 
		# line and flip the shape horizontally
		waveman --line --stroke-color yellow --stroke-width 5px \ 
			--closed --inverted --spread 50 \
			--downsampling-factor 64 --downsampling-mode head \
			-f audio.mp3
	`)
)

// NewWaveman creates a new cobra command and adds the relevant flags to the root command.
// It also creates the link to the subcommands
func NewWaveman(data *WavemanOptions, streams *IOStreams) *WavemanCommand {
	if data == nil {
		data = NewWavemanOptions()
	}

	waveman := &WavemanCommand{
		c: data,
		o: streams,
	}

	cmd := &cobra.Command{
		Use:     "waveman",
		Short:   WavemanShort,
		Long:    WavemanLong,
		Example: WavemanExamples,
	}

	// add transformer flags
	addTransformerOptions(cmd.Flags(), data.transformerData)
	addDimensionOptions(cmd.Flags(), data)

	waveman.cmd = cmd

	return waveman
}

type WavemanCommand struct {
	cmd     *cobra.Command
	c       *WavemanOptions
	painter *Plugin
	o       *IOStreams
}

// Plugin allows a user to patch in additional painters and register their flags to the waveman command.
func (w *WavemanCommand) Plugin(plugin Plugin) *WavemanCommand {
	painterName := plugin.Name()
	if _, ok := w.c.plugins[painterName]; ok {
		log.Warn().Msgf("painter %s is already registered, not replacing existing painter", painterName)
		return w
	}
	w.c.plugins[painterName] = plugin
	// add plugin flags
	plugin.Flags(w.cmd.Flags())
	return w
}

// Complete finalizes the Waveman configuration and creates a runner
func (w *WavemanCommand) Complete() *cobra.Command {
	// we finished all plugin registrations, so find the selected painter
	// NOTE: this does not perform the mutual exclusivity constraint check
	// for modes, this is done in Validate() at runtime of the cobra command

	w.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		err := w.c.Validate() // run all data validations
		if err != nil {
			return err
		}
		var selected *Plugin
		w.c.plugins.visit(func(plugin Plugin) bool {
			if *plugin.Enabled() {
				selected = &plugin
				return true
			}
			return false
		})

		// handle the case where no mode flag was set, fall back to box painter
		if selected == nil {
			log.Warn().Msgf("No painter selected, falling back to BoxPainter")
			plugin, ok := w.c.plugins[BoxMode]
			if !ok {
				log.Fatal().Msgf("Default box painter is not instantiated")
			}
			w.painter = &plugin
		}

		w.painter = selected

		// if we got here, we also have a selected painter available

		// TODO: implement me!
		return nil
	}

	return w.cmd
}

// NewWavemanOptions constructs a data struct to be used in the closure of the NewWavemanCommand
// constructor function when no struct is given as a parameter. This is the default case,
// the ability to pass data as a parameter is present to make unit testing commands
// possible
func NewWavemanOptions() *WavemanOptions {
	return &WavemanOptions{
		transformerData: newTransformerData(),
		plugins:         make(map[string]Plugin),
	}
}

type WavemanOptions struct {
	*transformerData
	height  float64
	width   float64
	plugins Plugins
}

// Validate first checks all properties of the transformer
func (o *WavemanOptions) Validate() error {
	transformerErrors := o.validateTransformerOptions()
	if transformerErrors != nil {
		return errors.New(transformerErrors.Error())
	}
	if err := validation.ValidateHeight(o.height); err != nil {
		return err
	}
	if err := validation.ValidateWidth(o.width); err != nil {
		return err
	}

	// mutually exclude plugins
	var mode string
	var collision string
	foundCollision := o.plugins.visit(func(plugin Plugin) bool {
		if *plugin.Enabled() {
			if mode == "" {
				mode = plugin.Name()
			} else {
				collision = plugin.Name()
				return true // break on first collision
			}
		}
		return false
	})
	if foundCollision {
		return fmt.Errorf("cannnot use --%s and --%s simultaneously", mode, collision)
	}
	var errs []error

	o.plugins.visit(func(plugin Plugin) bool {
		if err := plugin.Validate(); err != nil {
			errs = append(errs, err)
		}
		return false
	})

	if err := cmdutils.NewErrorList(errs); err != nil {
		return errors.New(err.Error())
	}
	return nil
}

func addDimensionOptions(flags *pflag.FlagSet, data *WavemanOptions) {
	flags.Float64VarP(&data.width, options.Width, options.WidthShort, painter.DefaultWidth, options.WidthDescription)
	flags.Float64VarP(&data.height, options.Height, options.HeightShort, painter.DefaultHeight, options.HeightDescription)
}
