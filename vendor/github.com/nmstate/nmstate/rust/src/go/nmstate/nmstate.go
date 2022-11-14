package nmstate

// #cgo CFLAGS: -g -Wall
// #cgo LDFLAGS: -lnmstate
// #include <nmstate.h>
// #include <stdlib.h>
import "C"
import (
	"fmt"
	"io"
	"time"
)

type Nmstate struct {
	timeout    uint
	logsWriter io.Writer
	flags      byte
}

const (
	kernelOnly = 2 << iota
	noVerify
	includeStatusData
	includeSecrets
	noCommit
	memoryOnly
	runningConfigOnly
)

func New(options ...func(*Nmstate)) *Nmstate {
	nms := &Nmstate{}
	for _, option := range options {
		option(nms)
	}
	return nms
}

func WithTimeout(timeout time.Duration) func(*Nmstate) {
	return func(n *Nmstate) {
		n.timeout = uint(timeout.Seconds())
	}
}

func WithLogsWritter(log_writter io.Writer) func(*Nmstate) {
	return func(n *Nmstate) {
		n.logsWriter = log_writter
	}
}

func WithKernelOnly() func(*Nmstate) {
	return func(n *Nmstate) {
		n.flags = n.flags | kernelOnly
	}
}

func WithNoVerify() func(*Nmstate) {
	return func(n *Nmstate) {
		n.flags = n.flags | noVerify
	}
}

func WithIncludeStatusData() func(*Nmstate) {
	return func(n *Nmstate) {
		n.flags = n.flags | includeStatusData
	}
}

func WithIncludeSecrets() func(*Nmstate) {
	return func(n *Nmstate) {
		n.flags = n.flags | includeSecrets
	}
}

func WithNoCommit() func(*Nmstate) {
	return func(n *Nmstate) {
		n.flags = n.flags | noCommit
	}
}

func WithMemoryOnly() func(*Nmstate) {
	return func(n *Nmstate) {
		n.flags = n.flags | memoryOnly
	}
}

func WithRunningConfigOnly() func(*Nmstate) {
	return func(n *Nmstate) {
		n.flags = n.flags | runningConfigOnly
	}
}

// Retrieve the network state in json format. This function returns the current
// network state or an error.
func (n *Nmstate) RetrieveNetState() (string, error) {
	var (
		state    *C.char
		log      *C.char
		err_kind *C.char
		err_msg  *C.char
	)
	rc := C.nmstate_net_state_retrieve(C.uint(n.flags), &state, &log, &err_kind, &err_msg)
	defer func() {
		C.nmstate_cstring_free(state)
		C.nmstate_cstring_free(err_msg)
		C.nmstate_cstring_free(err_kind)
		C.nmstate_cstring_free(log)
	}()
	if rc != 0 {
		return "", fmt.Errorf("failed retrieving nmstate net state with rc: %d, err_msg: %s, err_kind: %s", rc, C.GoString(err_msg), C.GoString(err_kind))
	}
	if err := n.writeLog(log); err != nil {
		return "", fmt.Errorf("failed when retrieving state: %v", err)
	}
	return C.GoString(state), nil
}

// Apply the network state in json format. This function returns the applied
// network state or an error.
func (n *Nmstate) ApplyNetState(state string) (string, error) {
	var (
		c_state  *C.char
		log      *C.char
		err_kind *C.char
		err_msg  *C.char
	)
	c_state = C.CString(state)
	rc := C.nmstate_net_state_apply(C.uint(n.flags), c_state, C.uint(n.timeout), &log, &err_kind, &err_msg)

	defer func() {
		C.nmstate_cstring_free(c_state)
		C.nmstate_cstring_free(err_msg)
		C.nmstate_cstring_free(err_kind)
		C.nmstate_cstring_free(log)
	}()
	if rc != 0 {
		return "", fmt.Errorf("failed applying nmstate net state %s with rc: %d, err_msg: %s, err_kind: %s", state, rc, C.GoString(err_msg), C.GoString(err_kind))
	}
	if err := n.writeLog(log); err != nil {
		return "", fmt.Errorf("failed when applying state: %v", err)
	}
	return state, nil
}

// Commit the checkpoint path provided. This function returns the committed
// checkpoint path or an error.
func (n *Nmstate) CommitCheckpoint(checkpoint string) (string, error) {
	var (
		c_checkpoint *C.char
		log          *C.char
		err_kind     *C.char
		err_msg      *C.char
	)
	c_checkpoint = C.CString(checkpoint)
	rc := C.nmstate_checkpoint_commit(c_checkpoint, &log, &err_kind, &err_msg)

	defer func() {
		C.nmstate_cstring_free(c_checkpoint)
		C.nmstate_cstring_free(err_msg)
		C.nmstate_cstring_free(err_kind)
		C.nmstate_cstring_free(log)
	}()
	if rc != 0 {
		return "", fmt.Errorf("failed commiting checkpoint %s with rc: %d, err_msg: %s, err_kind: %s", checkpoint, rc, C.GoString(err_msg), C.GoString(err_kind))
	}
	if err := n.writeLog(log); err != nil {
		return "", fmt.Errorf("failed when commiting: %v", err)
	}
	return checkpoint, nil
}

// Rollback to the checkpoint provided. This function returns the checkpoint
// path used for rollback or an error.
func (n *Nmstate) RollbackCheckpoint(checkpoint string) (string, error) {
	var (
		c_checkpoint *C.char
		log          *C.char
		err_kind     *C.char
		err_msg      *C.char
	)
	c_checkpoint = C.CString(checkpoint)
	rc := C.nmstate_checkpoint_rollback(c_checkpoint, &log, &err_kind, &err_msg)

	defer func() {
		C.nmstate_cstring_free(c_checkpoint)
		C.nmstate_cstring_free(err_msg)
		C.nmstate_cstring_free(err_kind)
		C.nmstate_cstring_free(log)
	}()
	if rc != 0 {
		return "", fmt.Errorf("failed when doing rollback checkpoint %s with rc: %d, err_msg: %s, err_kind: %s", checkpoint, rc, C.GoString(err_msg), C.GoString(err_kind))
	}
	if err := n.writeLog(log); err != nil {
		return "", fmt.Errorf("failed when doing rollback: %v", err)
	}
	return checkpoint, nil
}

func (n *Nmstate) writeLog(log *C.char) error {
	if n.logsWriter == nil {
		return nil
	}
	_, err := io.WriteString(n.logsWriter, C.GoString(log))
	if err != nil {
		return fmt.Errorf("failed writting logs: %v", err)
	}
	return nil
}

// GenerateConfiguration generates the configuration for the state provided.
// This function returns the configuration files for the state provided.
func (n *Nmstate) GenerateConfiguration(state string) (string, error) {
	var (
		c_state  *C.char
		config   *C.char
		log      *C.char
		err_kind *C.char
		err_msg  *C.char
	)
	c_state = C.CString(state)
	rc := C.nmstate_generate_configurations(c_state, &config, &log, &err_kind, &err_msg)

	defer func() {
		C.nmstate_cstring_free(c_state)
		C.nmstate_cstring_free(config)
		C.nmstate_cstring_free(err_msg)
		C.nmstate_cstring_free(err_kind)
		C.nmstate_cstring_free(log)
	}()
	if rc != 0 {
		return "", fmt.Errorf("failed when generating the configuration %s with rc: %d, err_msg: %s, err_kind: %s", state, rc, C.GoString(err_msg), C.GoString(err_kind))
	}
	if err := n.writeLog(log); err != nil {
		return "", fmt.Errorf("failed when generating the configuration: %v", err)
	}
	return C.GoString(config), nil
}
