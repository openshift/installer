package agent

import (
	"fmt"
	"net"
	"path"
	"path/filepath"
	"strconv"

	"github.com/sirupsen/logrus"

	gatherssh "github.com/openshift/installer/pkg/gather/ssh"
)

// PullAgentGatherArchive SSHs to the rendezvous host and runs the
// agent-gather script, pulling the resulting tar.xz archive to the
// local directory.
func PullAgentGatherArchive(rendezvousIP string, sshKeys []string, directory, gatherID string) (string, error) {
	logrus.Info("Pulling agent-gather data from the rendezvous host")

	address := net.JoinHostPort(rendezvousIP, strconv.Itoa(22))
	client, err := gatherssh.NewClient("core", address, sshKeys)
	if err != nil {
		return "", fmt.Errorf("failed to create SSH client for rendezvous host %s: %w", rendezvousIP, err)
	}

	// Run agent-gather with -i so it writes to a predictable path
	cmd := fmt.Sprintf("sudo /usr/local/bin/agent-gather -i %s", gatherID)
	if err := gatherssh.Run(client, cmd); err != nil {
		return "", fmt.Errorf("failed to run agent-gather on rendezvous host %s: %w", rendezvousIP, err)
	}

	archiveName := fmt.Sprintf("agent-gather-%s.tar.xz", gatherID)
	remoteFile := path.Join("/home/core", archiveName)
	localFile := filepath.Join(directory, archiveName)
	if err := gatherssh.PullFileTo(client, remoteFile, localFile); err != nil {
		return "", fmt.Errorf("failed to pull agent-gather archive: %w", err)
	}

	absPath, err := filepath.Abs(localFile)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	logrus.Info("Successfully pulled agent-gather data")
	return absPath, nil
}
