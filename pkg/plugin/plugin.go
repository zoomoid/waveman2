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

package plugin

import (
	"github.com/spf13/pflag"
	"github.com/zoomoid/waveman2/pkg/painter"
)

// FlagsFactory is a function type that, upon calling from the Plugin function, should register all flags the
// Plugin painter requires to the FlagSet of the command. Data is arbitrary data, that maps to the flag values.
type FlagsFactory func(flags *pflag.FlagSet, data interface{})

type Plugin interface {
	Validate() error
	Flags(flags *pflag.FlagSet) error
	Data() interface{}
	Name() string
	Description() string
	Enabled() *bool
	Draw(*painter.PainterOptions) []string
	Painter() painter.Painter
}

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
