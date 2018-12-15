// Package plugins is collection of all the terraform plugins that are used/required by installer.
package plugins

// KnownPlugins is a map of all the known plugin names to their exec functions.
var KnownPlugins = map[string]func(){}
