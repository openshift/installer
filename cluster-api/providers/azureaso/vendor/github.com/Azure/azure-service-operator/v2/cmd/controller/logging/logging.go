/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package logging

import (
	"flag"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"k8s.io/klog/v2/textlogger"

	"github.com/Azure/azure-service-operator/v2/internal/logging"
)

type Config struct {
	// verbosity indicates the level of logging.
	// Higher values indicate more logging.
	verbosity int

	// useJSON indicates whether we should output logs in JSON format.
	// Default is no
	useJSON bool
}

// Create returns a new logger, ready for use.
// This can be called multiple times if required.
func Create(cfg *Config) logr.Logger {
	if cfg != nil && cfg.useJSON {
		log, err := createJSONLogger(cfg)
		if err != nil {
			log = createTextLogger(cfg)
			log.Error(err, "failed to create JSON logger, falling back to text")
		}

		return log
	}

	return createTextLogger(cfg)
}

func createTextLogger(cfg *Config) logr.Logger {
	opts := []textlogger.ConfigOption{}
	if cfg != nil {
		opts = append(opts, textlogger.Verbosity(cfg.verbosity))
	}

	c := textlogger.NewConfig(opts...)
	return textlogger.NewLogger(c)
}

func createJSONLogger(cfg *Config) (logr.Logger, error) {
	level := zap.InfoLevel
	if cfg != nil {
		switch cfg.verbosity {
		case 0:
			level = zap.ErrorLevel
		case 1:
			level = zap.WarnLevel
		case 2:
			level = zap.InfoLevel
		default: // 3 or above
			level = zap.DebugLevel
		}
	}

	encoder := zap.NewProductionEncoderConfig()
	encoder.EncodeTime = zapcore.ISO8601TimeEncoder

	c := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    encoder,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := c.Build()
	if err != nil {
		return logr.Logger{}, err
	}

	return zapr.NewLogger(logger), nil
}

// InitFlags initializes the flags for the logging package
func InitFlags(fs *flag.FlagSet) *Config {
	result := &Config{}

	fs.IntVar(&result.verbosity, "verbose", logging.Verbose, "Enable verbose logging")
	fs.IntVar(&result.verbosity, "v", logging.Verbose, "Enable verbose logging")

	fs.BoolVar(&result.useJSON, "json-logging", false, "Enable JSON logging")

	return result
}
