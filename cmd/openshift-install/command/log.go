package command

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/version"
)

var (
	// RootOpts holds the log directory and log level configuration.
	RootOpts struct {
		Dir      string
		LogLevel string
	}

	// logCmdOnce logs the invoked openshift-install command details once.
	logCmdOnce sync.Once
)

type fileHook struct {
	file      io.Writer
	formatter logrus.Formatter
	level     logrus.Level

	truncateAtNewLine bool
}

func newFileHook(file io.Writer, level logrus.Level, formatter logrus.Formatter) *fileHook {
	return &fileHook{
		file:      file,
		formatter: formatter,
		level:     level,
	}
}

// NewFileHookWithNewlineTruncate returns a new FileHook with truncated new lines.
func NewFileHookWithNewlineTruncate(file io.Writer, level logrus.Level, formatter logrus.Formatter) *fileHook {
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

// redactCommandArgs returns a copy of args with sensitive flag values redacted.
// It handles both "--flag value" and "--flag=value" forms.
func redactCommandArgs(args []string) []string {
	// Flags whose values should be redacted
	sensitiveFlags := map[string]bool{
		"--key":    true,
		"--master": true,
	}

	redacted := make([]string, len(args))
	copy(redacted, args)

	for i := 0; i < len(redacted); i++ {
		arg := redacted[i]

		// Handle --flag=value form
		if strings.Contains(arg, "=") {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) == 2 && sensitiveFlags[parts[0]] {
				redacted[i] = parts[0] + "=<redacted>"
			}
			continue
		}

		// Handle --flag value form
		if sensitiveFlags[arg] && i+1 < len(redacted) {
			redacted[i+1] = "<redacted>"
			i++ // skip the next argument since we just redacted it
		}
	}

	return redacted
}

// SetupFileHook creates the base log directory and configures logrus options.
func SetupFileHook(baseDir string) func() {
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

	logCmdOnce.Do(func() {
		logrus.Debugf("Running: %s", strings.Join(redactCommandArgs(os.Args), " "))
		versionString, err := version.String()
		if err != nil {
			logrus.Fatal(err)
		}
		logrus.Debugf("%s", versionString)
		if version.Commit != "" {
			logrus.Debugf("Built from commit %s", version.Commit)
		}
	})

	return func() {
		logfile.Close()
		logrus.StandardLogger().ReplaceHooks(originalHooks)
	}
}
