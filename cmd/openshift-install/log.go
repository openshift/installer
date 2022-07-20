package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/version"
)

type fileHook struct {
	file      io.Writer
	formatter logrus.Formatter
	level     logrus.Level
	config    logConfig

	truncateAtNewLine bool
}

func newFileHook(file io.Writer, level logrus.Level, formatter logrus.Formatter) *fileHook {
	return &fileHook{
		file:      file,
		formatter: formatter,
		level:     level,
	}
}

func newFileHookWithNewlineTruncate(file io.Writer, level logrus.Level, formatter logrus.Formatter) *fileHook {
	f := newFileHook(file, level, formatter)
	f.truncateAtNewLine = true
	return f
}

func (h fileHook) Levels() []logrus.Level {
	var levels []logrus.Level
	for _, level := range logrus.AllLevels {
		if level <= h.level {
			levels = append(levels, level)
		}
	}

	return levels
}

func (h *fileHook) Fire(entry *logrus.Entry) error {
	// logrus reuses the same entry for each invocation of hooks.
	// so we need to make sure we leave them message field as we received.
	orig := entry.Message
	defer func() { entry.Message = orig }()

	msgs := []string{orig}
	if h.truncateAtNewLine {
		msgs = strings.Split(orig, "\n")
	}

	level := logrus.InfoLevel
	if h.config.Level != nil {
		parsedLevel, err := logrus.ParseLevel(*h.config.Level)
		if err == nil {
			level = parsedLevel
		} else {
			logrus.Debugf("failed to parse level %s: %s", *h.config.Level, err.Error())
		}
	}

	if h.config.Fields != nil && level <= entry.Level {
		entrySet := entry.WithFields(*h.config.Fields)
		entrySet.Level = entry.Level
		entry = entrySet
	}

	for _, msg := range msgs {
		// this makes it easier to call format on entry
		// easy without creating a new one for each split message.
		entry.Message = msg
		line, err := h.formatter.Format(entry)
		if err != nil {
			return err
		}

		if _, err := h.file.Write(line); err != nil {
			return err
		}
	}

	return nil
}

func setupFileHook(baseDir string) func() {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to create base directory for logs"))
	}

	logfile, err := os.OpenFile(filepath.Join(baseDir, ".openshift_install.log"), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to open log file"))
	}
	originalHooks := logrus.LevelHooks{}
	for k, v := range logrus.StandardLogger().Hooks {
		originalHooks[k] = v
	}

	fileHook := newFileHook(logfile, logrus.TraceLevel, &logrus.TextFormatter{
		DisableColors:          true,
		DisableTimestamp:       false,
		FullTimestamp:          true,
		DisableLevelTruncation: false,
	})

	config, err := readLogConfigFile(baseDir)
	if err == nil {
		fileHook.config = config
	} else {
		logrus.Debugf("failed to parse log-config.yaml: %s", err.Error())
	}
	logrus.AddHook(fileHook)

	versionString, err := version.String()
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Debugf(versionString)
	if version.Commit != "" {
		logrus.Debugf("Built from commit %s", version.Commit)
	}

	return func() {
		logfile.Close()
		logrus.StandardLogger().ReplaceHooks(originalHooks)
	}
}
