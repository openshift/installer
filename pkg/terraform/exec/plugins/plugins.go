// Package plugins is collection of all the terraform plugins that are used/required by installer.
package plugins

// TFpluginFunc is a callback function that runs the terraform plugin
type TFpluginFunc func()

// TFPlugin contains terraform plugin name, exec function, and plugin resources
type TFPlugin struct {
	Name string
	Exec TFpluginFunc
}

// KnownPlugins is a map of all the known plugin names to their exec functions.
var KnownPlugins = make(map[string]*TFPlugin)
