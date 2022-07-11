package cmd

import "github.com/spf13/pflag"

// FlagsFactory is a function type that, upon calling from the Plugin function, should register all flags the
// Plugin painter requires to the FlagSet of the command. Data is arbitrary data, that maps to the flag values.
type FlagsFactory func(flags *pflag.FlagSet, data interface{})

type Plugin interface {
	Validate() error
	Flags(flags *pflag.FlagSet)
	Data() interface{}
	Name() string
	Description() string
	Enabled() *bool
}

// Plugins wraps a map of plugins to provide the visit function
type Plugins map[string]Plugin

// visit implements the visitor pattern for the map of plugins contained in the Plugins type
func (p Plugins) visit(f func(plugin Plugin) bool) bool {
	for _, plugin := range p {
		if match := f(plugin); match {
			return match
		}
	}
	return false
}
