// Package plugins is collection of all the terraform plugins that are used/required by installer.
package plugins

import (
	"github.com/pkg/errors"
)

// TFpluginFunc is a callback function that runs the terraform plugin
type TFpluginFunc func()

// TFPlugin contains terraform plugin name, exec function, and plugin resources
type TFPlugin struct {
	Name      string
	Exec      TFpluginFunc
	Resources []string
}

// KnownPlugins is a map of all the known plugin names to their exec functions.
var KnownPlugins = make(map[string]*TFPlugin)

// GetPluginName returns the plugin name of a terraform provider
func GetPluginName(platformName string) (string, error) {
	for pluginName, plugin := range KnownPlugins {
		if plugin.Name == platformName {
			return pluginName, nil
		}
	}

	return "", errors.Errorf("unable to determine plugin name for platform %s", platformName)
}

// GetPlugin returns the terraform plugin
func GetPlugin(platformName string) (*TFPlugin, error) {
	pluginName, err := GetPluginName(platformName)
	if err != nil {
		return nil, err
	}

	plugin, ok := KnownPlugins[pluginName]
	if !ok {
		return nil, errors.Errorf("%s: no such plugin", pluginName)
	}

	return plugin, nil
}
