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

package plugin

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/zoomoid/waveman2/pkg/painter"
)

type Plugin interface {
	// Validate is a hook for plugins that ensures all parameters are correct
	Validate() error
	// Flags appends the plugins flags to its command
	Flags(flags *pflag.FlagSet) error
	// Completions hooks into the completions mechanism of cobra
	Completions(cmd *cobra.Command)
	// Data returns the internal struct the plugin uses to manage its configuration.
	// Needs type assertion to be worked with
	Data() interface{}
	// Group returns the plugin's group it belongs to. For core plugins, this is "core/v1",
	// other plugins may define the group at their own will
	Group() string
	// Name returns the plugin's name
	Name() string
	// Description returns the usage description as it should be printed by `waveman <plugin> --help`
	Description() string
	// Draw takes options configured from flags and the transformer stage, and converts the chunks of
	// samples to SVG elements. This may be more than one, hence a string slice. The core plugins wrap
	// their painter's output in SVG group elements, thus their slice only ever contains one element.
	Draw(*painter.PainterOptions) []string
	// Painter instantiates the plugin's painter backend and returns it
	Painter() painter.Painter
}

// FlagsFactory is a function type that, upon calling from the Plugin function, should register all flags the
// Plugin painter requires to the FlagSet of the command. Data is arbitrary data, that maps to the flag values.
type FlagsFactory func(flags *pflag.FlagSet, data interface{})

// NoCompletions is the baseline completions hook to be used by plugins that don't support completions
func NoCompletions(cmd *cobra.Command) {}

// Plugins wraps a map of plugins to provide the visit function
type Plugins map[string]Plugin

// visit implements the visitor pattern for the map of plugins contained in the Plugins type
func (p Plugins) Visit(f func(plugin Plugin) bool) bool {
	for _, plugin := range p {
		if match := f(plugin); match {
			return match
		}
	}
	return false
}
