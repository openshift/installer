package process

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"path"
	"regexp"
	"sync"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

// ListenAddr represents some listening address and port.
type ListenAddr struct {
	Address string
	Port    string
}

// URL returns a URL for this address with the given scheme and subpath.
func (l *ListenAddr) URL(scheme string, path string) *url.URL {
	return &url.URL{
		Scheme: scheme,
		Host:   l.HostPort(),
		Path:   path,
	}
}

// HostPort returns the joined host-port pair for this address.
func (l *ListenAddr) HostPort() string {
	return net.JoinHostPort(l.Address, l.Port)
}

// HealthCheck describes the information needed to health-check a process via
// some health-check URL.
type HealthCheck struct {
	url.URL

	// HealthCheckPollInterval is the interval which will be used for polling the
	// endpoint described by Host, Port, and Path.
	//
	// If left empty it will default to 100 Milliseconds.
	PollInterval time.Duration
}

// State define the state of the process.
type State struct {
	Cmd *exec.Cmd

	// HealthCheck describes how to check if this process is up.  If we get an http.StatusOK,
	// we assume the process is ready to operate.
	//
	// For example, the /healthz endpoint of the k8s API server, or the /health endpoint of etcd.
	HealthCheck *HealthCheck

	Dir  string
	Args []string
	Env  []string

	StopTimeout  time.Duration
	StartTimeout time.Duration

	Path string

	// ready holds whether the process is currently in ready state (hit the ready condition) or not.
	// It will be set to true on a successful `Start()` and set to false on a successful `Stop()`
	ready bool

	// waitDone is closed when our call to wait finishes up, and indicates that
	// our process has terminated.
	waitDone chan struct{}
	errMu    sync.Mutex
	exitErr  error
	exited   bool
}

// Init sets up this process, configuring binary paths if missing, initializing
// temporary directories, etc.
//
// This defaults all defaultable fields.
func (ps *State) Init(name string) error {
	if ps.Path == "" {
		if name == "" {
			return fmt.Errorf("must have at least one of name or path")
		}
	}

	if ps.StartTimeout == 0 {
		ps.StartTimeout = 20 * time.Second
	}

	if ps.StopTimeout == 0 {
		ps.StopTimeout = 20 * time.Second
	}
	return nil
}

type stopChannel chan struct{}

// CheckFlag checks the help output of this command for the presence of the given flag, specified
// without the leading `--` (e.g. `CheckFlag("insecure-port")` checks for `--insecure-port`),
// returning true if the flag is present.
func (ps *State) CheckFlag(flag string) (bool, error) {
	cmd := exec.Command(ps.Path, "--help") //nolint:gosec
	outContents, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("unable to run command %q to check for flag %q: %w", ps.Path, flag, err)
	}
	pat := `(?m)^\s*--` + flag + `\b` // (m --> multi-line --> ^ matches start of line)
	matched, err := regexp.Match(pat, outContents)
	if err != nil {
		return false, fmt.Errorf("unable to check command %q for flag %q in help output: %w", ps.Path, flag, err)
	}
	return matched, nil
}

// Start starts the apiserver, waits for it to come up, and returns an error,
// if occurred.
func (ps *State) Start(ctx context.Context, stdout io.Writer, stderr io.Writer) (err error) {
	if ps.ready {
		return nil
	}

	ps.Cmd = exec.CommandContext(ctx, ps.Path, ps.Args...) //nolint:gosec
	ps.Cmd.Env = append(ps.Cmd.Environ(), ps.Env...)
	ps.Cmd.Stdout = stdout
	ps.Cmd.Stderr = stderr
	ps.Cmd.Dir = ps.Dir
	ps.Cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	ready := make(chan bool)
	timedOut := time.After(ps.StartTimeout)

	pollerStopCh := make(stopChannel)
	if ps.HealthCheck != nil {
		go pollURLUntilOK(ps.HealthCheck.URL, ps.HealthCheck.PollInterval, ready, pollerStopCh)
	} else {
		// Assume that if we're not health-checking, we're ready to go.
		close(ready)
	}

	ps.waitDone = make(chan struct{})
	if err := ps.Cmd.Start(); err != nil {
		ps.errMu.Lock()
		defer ps.errMu.Unlock()
		ps.exited = true
		return err
	}
	go func() {
		defer close(ps.waitDone)
		err := ps.Cmd.Wait()

		ps.errMu.Lock()
		defer ps.errMu.Unlock()
		ps.exitErr = err
		ps.exited = true
	}()

	select {
	case <-ready:
		ps.ready = true
		return nil
	case <-ps.waitDone:
		close(pollerStopCh)
		return fmt.Errorf("timeout waiting for process %s to start successfully "+
			"(it may have failed to start, or stopped unexpectedly before becoming ready)",
			path.Base(ps.Path))
	case <-timedOut:
		close(pollerStopCh)
		if ps.Cmd != nil {
			// intentionally ignore this -- we might've crashed, failed to start, etc
			ps.Cmd.Process.Signal(syscall.SIGTERM) //nolint:errcheck
		}
		return fmt.Errorf("timeout waiting for process %s to start", path.Base(ps.Path))
	}
}

// Exited returns true if the process exited, and may also
// return an error (as per Cmd.Wait) if the process did not
// exit with error code 0.
func (ps *State) Exited() (bool, error) {
	ps.errMu.Lock()
	defer ps.errMu.Unlock()
	return ps.exited, ps.exitErr
}

func pollURLUntilOK(url url.URL, interval time.Duration, ready chan bool, stopCh stopChannel) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				// there's probably certs *somewhere*,
				// but it's fine to just skip validating
				// them for health checks during testing
				InsecureSkipVerify: true, //nolint:gosec
			},
		},
	}
	if interval <= 0 {
		interval = 100 * time.Millisecond
	}
	for {
		res, err := client.Get(url.String())
		if err == nil {
			res.Body.Close()
			if res.StatusCode == http.StatusOK {
				ready <- true
				return
			}
		}

		select {
		case <-stopCh:
			return
		default:
			time.Sleep(interval)
		}
	}
}

// Stop stops this process gracefully, waits for its termination.
func (ps *State) Stop() error {
	if ps.Cmd == nil {
		return nil
	}
	if done, err := ps.Exited(); done {
		if err != nil {
			logrus.Warnf("process %s exited with error: %v", path.Base(ps.Path), err)
		}
		return nil
	}
	if err := ps.Cmd.Process.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("unable to signal for process %s to stop: %w", ps.Path, err)
	}

	timedOut := time.After(ps.StopTimeout)
	select {
	case <-ps.waitDone:
		break
	case <-timedOut:
		if err := ps.Cmd.Process.Signal(syscall.SIGKILL); err != nil {
			return fmt.Errorf("unable to signal for process %s to stop: %w", ps.Path, err)
		}
		return fmt.Errorf("timeout waiting for process %s to stop, sent SIGKILL", path.Base(ps.Path))
	}
	ps.ready = false
	return nil
}
