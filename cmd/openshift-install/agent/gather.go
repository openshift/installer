package agent

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/openshift/installer/cmd/openshift-install/command"
	agentpkg "github.com/openshift/installer/pkg/agent"
	assetstore "github.com/openshift/installer/pkg/asset/store"
	"github.com/openshift/installer/pkg/asset/tls"
)

var agentGatherOpts struct {
	sshKeys []string
}

// NewGatherCmd creates the commands for gathering debug data from an agent-based installation.
func NewGatherCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gather",
		Short: "Gather debugging data for a failed agent-based installation",
		Long: `Gather debugging data for a failed agent-based installation.

When an agent-based installation fails, this command collects debugging
data from the rendezvous host to help diagnose the issue.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(newAgentGatherCmd())
	return cmd
}

func newAgentGatherCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bootstrap",
		Short: "Gather debugging data from the rendezvous host",
		Args:  cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			cleanup := command.SetupFileHook(command.RootOpts.Dir)
			defer cleanup()

			bundlePath, err := runAgentGatherCmd(command.RootOpts.Dir)
			if err != nil {
				logrus.Fatal(err)
			}
			logrus.Infof("Agent gather logs captured here %q", bundlePath)
		},
	}

	cmd.PersistentFlags().StringArrayVar(&agentGatherOpts.sshKeys, "key", []string{},
		"Path to SSH private keys that should be used for authentication. "+
			"If no key was provided, SSH private keys from user's environment will be used")
	return cmd
}

func runAgentGatherCmd(directory string) (string, error) {
	ctx := context.TODO()

	store, err := assetstore.NewStore(directory)
	if err != nil {
		return "", fmt.Errorf("failed to create asset store: %w", err)
	}

	rendezvousIP, err := agentpkg.FindRendezvousIPFromAssetStore(store)
	if err != nil {
		return "", fmt.Errorf("failed to determine rendezvous host: %w", err)
	}
	logrus.Infof("Rendezvous host IP: %s", rendezvousIP)

	// add the bootstrap SSH key pair to the sshKeys list automatically
	bootstrapSSHKeyPair := &tls.BootstrapSSHKeyPair{}
	if err := store.Fetch(ctx, bootstrapSSHKeyPair); err != nil {
		logrus.Debugf("Failed to fetch bootstrap SSH key pair: %v", err)
	} else {
		tmpfile, err := os.CreateTemp("", "bootstrap-ssh")
		if err != nil {
			return "", err
		}
		defer os.Remove(tmpfile.Name())
		if _, err := tmpfile.Write(bootstrapSSHKeyPair.Private()); err != nil {
			return "", err
		}
		if err := tmpfile.Close(); err != nil {
			return "", err
		}
		agentGatherOpts.sshKeys = append(agentGatherOpts.sshKeys, tmpfile.Name())
	}

	gatherID := time.Now().Format("20060102150405")

	bundlePath, err := agentpkg.PullAgentGatherArchive(rendezvousIP, agentGatherOpts.sshKeys, directory, gatherID)
	if err != nil {
		return "", fmt.Errorf("failed to gather data from rendezvous host: %w", err)
	}

	return bundlePath, nil
}
