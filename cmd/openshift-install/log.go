package main

import (
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/version"
)

type fileHook struct {
	file      io.Writer
	formatter logrus.Formatter
	level     logrus.Level
}

func newFileHook(file io.Writer, level logrus.Level, formatter logrus.Formatter) *fileHook {
	return &fileHook{
		file:      file,
		formatter: formatter,
		level:     level,
	}
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
	line, err := h.formatter.Format(entry)
	if err != nil {
		return err
	}

	_, err = h.file.Write(line)
	return err
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
	logrus.AddHook(newFileHook(logfile, logrus.TraceLevel, &logrus.TextFormatter{
		DisableColors:          true,
		DisableTimestamp:       false,
		FullTimestamp:          true,
		DisableLevelTruncation: false,
	}))

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
