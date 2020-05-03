// Package exec is glue between the vendored terraform codebase and installer.
package exec

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/logutils"
	"github.com/hashicorp/terraform-plugin-sdk/helper/logging"
	"github.com/hashicorp/terraform/command"
	"github.com/mitchellh/cli"
)

type cmdFunc func(command.Meta) cli.Command

var commands = map[string]cmdFunc{
	"apply": func(meta command.Meta) cli.Command {
		return &command.ApplyCommand{Meta: meta}
	},
	"destroy": func(meta command.Meta) cli.Command {
		return &command.ApplyCommand{Meta: meta, Destroy: true}
	},
	"init": func(meta command.Meta) cli.Command {
		return &command.InitCommand{Meta: meta}
	},
}

func runner(cmd string, dir string, args []string, stdout, stderr io.Writer) int {
	lf := ioutil.Discard
	if level := logging.LogLevel(); level != "" {
		lf = &logutils.LevelFilter{
			Levels:   logging.ValidLevels,
			MinLevel: logutils.LogLevel(level),
			Writer:   stdout,
		}
	}
	log.SetOutput(lf)
	defer log.SetOutput(os.Stderr)

	// Make sure we clean up any managed plugins at the end of this
	defer plugin.CleanupClients()

	sdCh, cancel := makeShutdownCh()
	defer cancel()

	pluginDirs, err := globalPluginDirs(dir)
	if err != nil {
		fmt.Fprintf(stderr, "Error discovering plugin directories for Terraform: %v", err)
		return 1
	}
	meta := command.Meta{
		Color:            false,
		GlobalPluginDirs: pluginDirs,
		Ui: &supressedUI{
			Ui: &cli.BasicUi{
				Writer:      stdout,
				ErrorWriter: stderr,
			},
		},

		OverrideDataDir: filepath.Join(dir, ".tf"),

		ShutdownCh: sdCh,
	}

	f := commands[cmd]

	oldStderr := os.Stderr
	outR, outW, err := os.Pipe()
	if err != nil {
		fmt.Fprintf(stderr, "error creating Pipe: %v", err)
		return 1
	}
	os.Stderr = outW
	go func() {
		scanner := bufio.NewScanner(outR)
		for scanner.Scan() {
			fmt.Fprintf(lf, "%s\n", scanner.Bytes())
		}
	}()
	defer func() {
		outW.Close()
		os.Stderr = oldStderr
	}()
	return f(meta).Run(args)
}

// Apply is wrapper around `terraform apply` subcommand.
func Apply(datadir string, args []string, stdout, stderr io.Writer) int {
	return runner("apply", datadir, args, stdout, stderr)
}

// Destroy is wrapper around `terraform destroy` subcommand.
func Destroy(datadir string, args []string, stdout, stderr io.Writer) int {
	return runner("destroy", datadir, args, stdout, stderr)
}

// Init is wrapper around `terraform init` subcommand.
func Init(datadir string, args []string, stdout, stderr io.Writer) int {
	return runner("init", datadir, args, stdout, stderr)
}

// makeShutdownCh creates an interrupt listener and returns a channel.
// A message will be sent on the channel for every interrupt received.
func makeShutdownCh() (<-chan struct{}, func()) {
	resultCh := make(chan struct{})
	signalCh := make(chan os.Signal, 4)

	handle := []os.Signal{}
	handle = append(handle, ignoreSignals...)
	handle = append(handle, forwardSignals...)

	signal.Notify(signalCh, handle...)
	go func() {
		for {
			<-signalCh
			resultCh <- struct{}{}
		}
	}()

	return resultCh, func() { signal.Reset(handle...) }
}

// suppressedUI suppresses the Ui's warnings from error to
// info.
type supressedUI struct {
	cli.Ui
}

func (sui *supressedUI) Warn(msg string) {
	sui.Ui.Info(msg)
}
