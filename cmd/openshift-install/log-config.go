package main

import (
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type logConfig struct {
	// Fields contains user-provided fields to be added to all log entries.
	// +optional
	Fields *logrus.Fields `json:"fields,omitempty"`

	// Level sets the level of logging to standard out.
	// Valid values: panic, fatal, error, warning, info, debug, trace
	// +optional
	Level *string `json:"level,omitempty"`
}

func readLogConfigFile(directory string) (logConfig logConfig, err error) {
	file, err := os.ReadFile(path.Join(directory, "log-config.yaml"))
	if err != nil {
		return
	}
	err = yaml.Unmarshal(file, &logConfig)
	return
}
