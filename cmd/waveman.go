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
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/zoomoid/waveman2/cmd/options"
	"github.com/zoomoid/waveman2/cmd/validation"
	"github.com/zoomoid/waveman2/pkg/painter"
	"github.com/zoomoid/waveman2/pkg/plugin"
	r "github.com/zoomoid/waveman2/pkg/reference"
	"github.com/zoomoid/waveman2/pkg/streams"
	"github.com/zoomoid/waveman2/pkg/svg"
	"github.com/zoomoid/waveman2/pkg/transform"
	"github.com/zoomoid/waveman2/pkg/utils"
	"github.com/zoomoid/waveman2/pkg/visitor"
)

type Waveman struct {
	cmd     *cobra.Command
	options *WavemanOptions
	painter *plugin.Plugin
	io      *streams.IO
}

var Version = "0.0.0-dev.0"

// NewWaveman creates a new cobra command and adds the relevant flags to the root command.
// It also creates the link to the subcommands
func NewWaveman(data *WavemanOptions, streams *streams.IO) *Waveman {
	if data == nil {
		data = NewWavemanOptions()
	}

	waveman := &Waveman{
		options: data,
		io:      streams,
	}

	cmd := &cobra.Command{
		Use:     "waveman",
		Short:   WavemanShort,
		Long:    WavemanLong,
		Example: WavemanExamples,
		Version: Version,
	}

	// add transformer flags
	addTranformerFlags(cmd.PersistentFlags(), data.transformerData)
	addTransformerFlagCompletion(cmd)
	// add shared painter flags, like height and width
	addDimensionFlags(cmd.PersistentFlags(), data.sharedPainterOptions)
	addDimensionFlagsCompletion(cmd)
	// add -f/-o/-r flags
	addIOFlags(cmd.PersistentFlags(), data.filenameOptions)
	addIOFlagsCompletion(cmd)

	waveman.cmd = cmd

	return waveman
}

func (w *Waveman) V(version string) *Waveman {
	// log.Debug().Msgf("Using waveman version %s", version)
	Version = version
	w.cmd.Version = version
	return w
}

// WavemanOptions contains all data passed into waveman as flags
type WavemanOptions struct {
	*transformerData
	*filenameOptions
	*sharedPainterOptions
	plugins plugin.Plugins
}

// Plugin allows a user to patch in additional painters and register their flags to the waveman command.
func (w *Waveman) Plugin(plugin plugin.Plugin) *Waveman {
	painterName := plugin.Name()
	if _, ok := w.options.plugins[painterName]; ok {
		log.Warn().
			Msgf("painter %s is already registered, not replacing existing painter", painterName)
		return w
	}
	w.options.plugins[painterName] = plugin
	// add plugin flags
	err := plugin.Flags(w.cmd.PersistentFlags())
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("")
	}
	return w
}

// Complete finalizes the Waveman configuration and creates a runner
func (w *Waveman) Complete() *cobra.Command {
	w.cmd.RunE = func(_ *cobra.Command, _ []string) error {
		err := w.options.Validate() // run all data validations
		if err != nil {
			return err
		}

		// we finished all plugin registrations, so find the selected painter
		// NOTE: this does not perform the mutual exclusivity constraint check
		// for modes, this is done in Validate() at runtime of the cobra command
		var selected *plugin.Plugin
		w.options.plugins.Visit(func(plugin plugin.Plugin) bool {
			if *plugin.Enabled() {
				selected = &plugin
				return true
			}
			return false
		})

		// handle the case where no mode flag was set, fall back to box painter
		if selected == nil {
			log.Warn().Msgf("No painter selected, falling back to BoxPainter")
			plugin, ok := w.options.plugins[r.BoxMode]
			if !ok {
				log.Fatal().Msgf("Default box painter is not instantiated")
			}
			selected = &plugin
		}

		w.painter = selected

		o := w.options

		// use stdout only if a singleton file is given and the -o flag did not specify elsewise
		// when -o is not specified, assume output to Stdout.
		// Validation during expansion of -f flags (--recursive included) will fail if more
		// than one input is present
		// if more than one file is passed with -f flags, we can skip this, because then we
		// will have to create parallel output files
		useStdout := options.OutputType(o.output) == options.OutputTypeEmpty && len(o.filenames) <= 1
		filenames := o.filenames
		recursive := o.recursive

		// expand all paths given to the CLI into visitors
		visitors, errs := visitor.ExpandPaths(filenames, recursive, useStdout, w.io)

		// catch any errors encountered in the process
		if el := utils.NewErrorList(errs); errs != nil {
			log.Fatal().Msg(el.Error())
		}

		err = visitors.
			ContinueOnError().
			UseStdout(useStdout).
			Visit(func(f *visitor.File) error {
				transformer, err := transform.New(o.transformerData.toOptions(), f.Reader())
				if err != nil {
					return err
				}
				samples := transformer.Blocks()
				p := *w.painter
				if p == nil {
					return fmt.Errorf("painter is nil")
				}
				elements := p.Draw(&painter.PainterOptions{
					Data:   samples,
					Height: o.height,
					Width:  o.width,
				})
				out, err := svg.Template(elements, true, p.Painter().Viewbox())
				if err != nil {
					return err
				}
				f.Print(out)

				return nil
			})

		return err
	}

	addShellCompletionSubcommand(w.cmd)

	return w.cmd
}

// NewWavemanOptions constructs a data struct to be used in the closure of the NewWavemanCommand
// constructor function when no struct is given as a parameter. This is the default case,
// the ability to pass data as a parameter is present to make unit testing commands
// possible
func NewWavemanOptions() *WavemanOptions {
	return &WavemanOptions{
		transformerData:      newTransformerData(),
		plugins:              make(map[string]plugin.Plugin),
		sharedPainterOptions: newSharedPainterData(),
		filenameOptions:      newFilenameData(),
	}
}

func newSharedPainterData() *sharedPainterOptions {
	return &sharedPainterOptions{
		height: painter.DefaultHeight,
		width:  painter.DefaultWidth,
	}
}

func newFilenameData() *filenameOptions {
	return &filenameOptions{
		filenames: []string{},
		recursive: false,
	}
}

// Validate checks all flags and data bindings to fulfill their defined
// validation functions and returns early when a value does not satisfy the conditions
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
	if err := validation.ValidateOutput(o.output); err != nil {
		return err
	}
	// if err := validation.ValidateFilenames(o.filenames); err != nil {
	// 	return err
	// }

	// mutually exclude plugins
	var mode string
	var collision string
	foundCollision := o.plugins.Visit(func(p plugin.Plugin) bool {
		if *p.Enabled() {
			if mode == "" {
				mode = p.Name()
			} else {
				collision = p.Name()
				return true // break on first collision
			}
		}
		return false
	})
	if foundCollision {
		return fmt.Errorf("cannnot use --%s and --%s simultaneously", mode, collision)
	}

	var errs []error
	// Run each plugin's validation hook
	o.plugins.Visit(func(p plugin.Plugin) bool {
		if err := p.Validate(); err != nil {
			errs = append(errs, err)
		}
		return false
	})

	if err := utils.NewErrorList(errs); err != nil {
		return errors.New(err.Error())
	}
	return nil
}
