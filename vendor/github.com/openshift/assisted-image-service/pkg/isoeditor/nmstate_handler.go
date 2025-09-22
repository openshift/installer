package isoeditor

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

//go:generate mockgen -package=isoeditor -destination=mock_nmstate_handler.go . NmstateHandler
type NmstateHandler interface {
	BuildNmstateCpioArchive(rootfsPath string) ([]byte, error)
}

type nmstateHandler struct {
	workDir  string
	executer Executer
}

func NewNmstateHandler(workDir string, executer Executer) NmstateHandler {
	return &nmstateHandler{
		workDir:  workDir,
		executer: executer,
	}
}

func (n *nmstateHandler) BuildNmstateCpioArchive(rootfsPath string) ([]byte, error) {
	// Extract nmstatectl binary
	var err error
	nmstateDir := filepath.Join(n.workDir, "nmstate")
	err = os.MkdirAll(nmstateDir, os.ModePerm)
	if err != nil {
		return nil, err
	}
	// Remove temp dir
	defer func() {
		removeErr := os.RemoveAll(nmstateDir)
		if removeErr != nil {
			log.WithError(removeErr).Error("failed to remove nmstate temp dir")
		}
	}()

	binaryPath, err := n.extractNmstatectl(rootfsPath, nmstateDir)
	if err != nil {
		return nil, err
	}
	nmstatectlPath := filepath.Join(nmstateDir, "squashfs-root", binaryPath)

	// Check if nmstatectl binary file exists
	if _, err = os.Stat(nmstatectlPath); os.IsNotExist(err) {
		return nil, err
	}

	// Read binary
	nmstateBinContent, err := os.ReadFile(nmstatectlPath)
	if err != nil {
		return nil, err
	}

	// Create a compressed RAM disk image with the nmstatectl binary
	compressedCpio, err := generateCompressedCPIO([]fileEntry{
		{
			Content: nmstateBinContent,
			Path:    NmstatectlPathInRamdisk,
			Mode:    0o100_755,
		},
	})
	if err != nil {
		return nil, err
	}

	return compressedCpio, err
}

// TODO: Update the code to utilize go-diskfs's squashfs instead of unsquashfs once go-diskfs supports the zstd compression format used by CoreOS - MGMT-19227
func (n *nmstateHandler) extractNmstatectl(rootfsPath, nmstateDir string) (string, error) {
	_, err := n.executer.Execute(fmt.Sprintf("cat %s | cpio -i", rootfsPath), nmstateDir)
	if err != nil {
		log.Errorf("failed to extract rootfs.img using cpio command: %v", err.Error())
		return "", err
	}
	// limiting files is needed on el<=9 due to https://github.com/plougher/squashfs-tools/issues/125
	ulimit := "ulimit -n 1024"

	// Listing the filesystem concisely, displaying only files (using `-lc` option).
	// Each file in the output won't include any prefix before `/ostree` (by using `-dest ''` option),
	// which is useful when invoking `-extract-file` (after finding the `nmstatectl` binary path).
	list, err := n.executer.Execute(fmt.Sprintf("%s ; unsquashfs -dest '' -lc %s", ulimit, "root.squashfs"), nmstateDir)
	if err != nil {
		log.Errorf("failed to unsquashfs root.squashfs: %v", err.Error())
		return "", err
	}

	r, err := regexp.Compile(".*nmstatectl")
	if err != nil {
		log.Errorf("failed to compile regexp: %v", err.Error())
		return "", err
	}
	binaryPath := r.FindString(list)

	_, err = n.executer.Execute(fmt.Sprintf("%s ; unsquashfs -no-xattrs %s -extract-file %s", ulimit, "root.squashfs", binaryPath), nmstateDir)
	if err != nil {
		log.Errorf("failed to unsquashfs root.squashfs: %v", err.Error())
		return "", err
	}
	return binaryPath, nil
}

//go:generate mockgen -package=isoeditor -destination=mock_executer.go . Executer
type Executer interface {
	Execute(command, workDir string) (string, error)
}

type CommonExecuter struct{}

func (e *CommonExecuter) Execute(command, workDir string) (string, error) {
	var stdoutBytes, stderrBytes bytes.Buffer
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = &stdoutBytes
	cmd.Stderr = &stderrBytes
	log.Infof(fmt.Sprintf("Running cmd: %s\n", command))
	cmd.Dir = workDir
	err := cmd.Run()
	if err != nil {
		return "", errors.Wrapf(err, "Failed to execute cmd (%s): %s", cmd, stderrBytes.String())
	}

	return strings.TrimSuffix(stdoutBytes.String(), "\n"), nil
}
