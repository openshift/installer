//go:generate go run plugin_versions_generate.go ../../../../go.mod plugin_versions.go

// Package plugins is collection of all the terraform plugins that are used/required by installer.
package plugins

type tfpluginFunc func()

// TFPlugin contains terraform exec function and plugin version
type TFPlugin struct {
	Name    string
	Exec    tfpluginFunc
	Version string
}

// KnownPlugins is a map of all the known plugin names to their exec functions.
var KnownPlugins = make(map[string]*TFPlugin)
