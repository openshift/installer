// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testscript

import (
	cryptorand "crypto/rand"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// TestingM is implemented by *testing.M. It's defined as an interface
// to allow testscript to co-exist with other testing frameworks
// that might also wish to call M.Run.
type TestingM interface {
	Run() int
}

var ignoreMissedCoverage = false

// IgnoreMissedCoverage causes any missed coverage information
// (for example when a function passed to RunMain
// calls os.Exit, for example) to be ignored.
// This function should be called before calling RunMain.
func IgnoreMissedCoverage() {
	ignoreMissedCoverage = true
}

// RunMain should be called within a TestMain function to allow
// subcommands to be run in the testscript context.
//
// The commands map holds the set of command names, each
// with an associated run function which should return the
// code to pass to os.Exit. It's OK for a command function to
// exit itself, but this may result in loss of coverage information.
//
// When Run is called, these commands are installed as regular commands in the shell
// path, so can be invoked with "exec" or via any other command (for example a shell script).
//
// For backwards compatibility, the commands declared in the map can be run
// without "exec" - that is, "foo" will behave like "exec foo".
// This can be disabled with Params.RequireExplicitExec to keep consistency
// across test scripts, and to keep separate process executions explicit.
//
// This function returns an exit code to pass to os.Exit, after calling m.Run.
func RunMain(m TestingM, commands map[string]func() int) (exitCode int) {
	// Depending on os.Args[0], this is either the top-level execution of
	// the test binary by "go test", or the execution of one of the provided
	// commands via "foo" or "exec foo".

	cmdName := filepath.Base(os.Args[0])
	if runtime.GOOS == "windows" {
		cmdName = strings.TrimSuffix(cmdName, ".exe")
	}
	mainf := commands[cmdName]
	if mainf == nil {
		// Unknown command; this is just the top-level execution of the
		// test binary by "go test".

		// Set up all commands in a directory, added in $PATH.
		tmpdir, err := ioutil.TempDir("", "testscript-main")
		if err != nil {
			log.Printf("could not set up temporary directory: %v", err)
			return 2
		}
		defer func() {
			if err := os.RemoveAll(tmpdir); err != nil {
				log.Printf("cannot delete temporary directory: %v", err)
				exitCode = 2
			}
		}()
		bindir := filepath.Join(tmpdir, "bin")
		if err := os.MkdirAll(bindir, 0o777); err != nil {
			log.Printf("could not set up PATH binary directory: %v", err)
			return 2
		}
		os.Setenv("PATH", bindir+string(filepath.ListSeparator)+os.Getenv("PATH"))

		flag.Parse()
		// If we are collecting a coverage profile, set up a shared
		// directory for all executed test binary sub-processes to write
		// their profiles to. Before finishing, we'll merge all of those
		// profiles into the main profile.
		if coverProfile() != "" {
			coverdir := filepath.Join(tmpdir, "cover")
			if err := os.MkdirAll(coverdir, 0o777); err != nil {
				log.Printf("could not set up cover directory: %v", err)
				return 2
			}
			os.Setenv("TESTSCRIPT_COVER_DIR", coverdir)
			defer func() {
				if err := finalizeCoverProfile(coverdir); err != nil {
					log.Printf("cannot merge cover profiles: %v", err)
					exitCode = 2
				}
			}()
		}

		// We're not in a subcommand.
		for name := range commands {
			name := name
			// Set up this command in the directory we added to $PATH.
			binfile := filepath.Join(bindir, name)
			if runtime.GOOS == "windows" {
				binfile += ".exe"
			}
			binpath, err := os.Executable()
			if err == nil {
				err = copyBinary(binpath, binfile)
			}
			if err != nil {
				log.Printf("could not set up %s in $PATH: %v", name, err)
				return 2
			}
			scriptCmds[name] = func(ts *TestScript, neg bool, args []string) {
				if ts.params.RequireExplicitExec {
					ts.Fatalf("use 'exec %s' rather than '%s' (because RequireExplicitExec is enabled)", name, name)
				}
				ts.cmdExec(neg, append([]string{name}, args...))
			}
		}
		return m.Run()
	}
	// The command being registered is being invoked, so run it, then exit.
	os.Args[0] = cmdName
	coverdir := os.Getenv("TESTSCRIPT_COVER_DIR")
	if coverdir == "" {
		// No coverage, act as normal.
		return mainf()
	}

	// For a command "foo", write ${TESTSCRIPT_COVER_DIR}/foo-${RANDOM}.
	// Note that we do not use ioutil.TempFile as that creates the file.
	// In this case, we want to leave it to -test.coverprofile to create the
	// file, as otherwise we could end up with an empty file.
	// Later, when merging profiles, trying to merge an empty file would
	// result in a confusing error.
	rnd, err := nextRandom()
	if err != nil {
		log.Printf("could not obtain random number: %v", err)
		return 2
	}
	cprof := filepath.Join(coverdir, fmt.Sprintf("%s-%x", cmdName, rnd))
	return runCoverSubcommand(cprof, mainf)
}

func nextRandom() ([]byte, error) {
	p := make([]byte, 6)
	_, err := cryptorand.Read(p)
	return p, err
}

// copyBinary makes a copy of a binary to a new location. It is used as part of
// setting up top-level commands in $PATH.
//
// It does not attempt to use symlinks for two reasons:
//
// First, some tools like cmd/go's -toolexec will be clever enough to realise
// when they're given a symlink, and they will use the symlink target for
// executing the program. This breaks testscript, as we depend on os.Args[0] to
// know what command to run.
//
// Second, symlinks might not be available on some environments, so we have to
// implement a "full copy" fallback anyway.
//
// However, we do try to use a hard link, since that will probably work on most
// unix-like setups. Note that "go test" also places test binaries in the
// system's temporary directory, like we do. We don't use hard links on Windows,
// as that can lead to "access denied" errors when removing.
func copyBinary(from, to string) error {
	if runtime.GOOS != "windows" {
		if err := os.Link(from, to); err == nil {
			return nil
		}
	}
	writer, err := os.OpenFile(to, os.O_WRONLY|os.O_CREATE, 0o777)
	if err != nil {
		return err
	}
	defer writer.Close()

	reader, err := os.Open(from)
	if err != nil {
		return err
	}
	defer reader.Close()

	_, err = io.Copy(writer, reader)
	return err
}

// runCoverSubcommand runs the given function, then writes any generated
// coverage information to the cprof file.
// This is called inside a separately run executable.
func runCoverSubcommand(cprof string, mainf func() int) (exitCode int) {
	// Change the error handling mode to PanicOnError
	// so that in the common case of calling flag.Parse in main we'll
	// be able to catch the panic instead of just exiting.
	flag.CommandLine.Init(flag.CommandLine.Name(), flag.PanicOnError)
	defer func() {
		panicErr := recover()
		if err, ok := panicErr.(error); ok {
			// The flag package will already have printed this error, assuming,
			// that is, that the error was created in the flag package.
			// TODO check the stack to be sure it was actually raised by the flag package.
			exitCode = 2
			if err == flag.ErrHelp {
				exitCode = 0
			}
			panicErr = nil
		}
		// Set os.Args so that flag.Parse will tell testing the correct
		// coverprofile setting. Unfortunately this isn't sufficient because
		// the testing oackage explicitly avoids calling flag.Parse again
		// if flag.Parsed returns true, so we the coverprofile value directly
		// too.
		os.Args = []string{os.Args[0], "-test.coverprofile=" + cprof}
		setCoverProfile(cprof)

		// Suppress the chatty coverage and test report.
		devNull, err := os.Open(os.DevNull)
		if err != nil {
			panic(err)
		}
		os.Stdout = devNull
		os.Stderr = devNull

		// Run MainStart (recursively, but it we should be ok) with no tests
		// so that it writes the coverage profile.
		m := mainStart()
		if code := m.Run(); code != 0 && exitCode == 0 {
			exitCode = code
		}
		if _, err := os.Stat(cprof); err != nil {
			log.Printf("failed to write coverage profile %q", cprof)
		}
		if panicErr != nil {
			// The error didn't originate from the flag package (we know that
			// flag.PanicOnError causes an error value that implements error),
			// so carry on panicking.
			panic(panicErr)
		}
	}()
	return mainf()
}

func coverProfileFlag() flag.Getter {
	f := flag.CommandLine.Lookup("test.coverprofile")
	if f == nil {
		// We've imported testing so it definitely should be there.
		panic("cannot find test.coverprofile flag")
	}
	return f.Value.(flag.Getter)
}

func coverProfile() string {
	return coverProfileFlag().Get().(string)
}

func setCoverProfile(cprof string) {
	coverProfileFlag().Set(cprof)
}

type nopTestDeps struct{}

func (nopTestDeps) MatchString(pat, str string) (result bool, err error) {
	return false, nil
}

func (nopTestDeps) StartCPUProfile(w io.Writer) error {
	return nil
}

func (nopTestDeps) StopCPUProfile() {}

func (nopTestDeps) WriteProfileTo(name string, w io.Writer, debug int) error {
	return nil
}

func (nopTestDeps) ImportPath() string {
	return ""
}
func (nopTestDeps) StartTestLog(w io.Writer) {}

func (nopTestDeps) StopTestLog() error {
	return nil
}

// Note: WriteHeapProfile is needed for Go 1.10 but not Go 1.11.
func (nopTestDeps) WriteHeapProfile(io.Writer) error {
	// Not needed for Go 1.10.
	return nil
}

// Note: SetPanicOnExit0 was added in Go 1.16.
func (nopTestDeps) SetPanicOnExit0(bool) {}
