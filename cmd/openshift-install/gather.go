package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/openshift/installer/pkg/asset/installconfig"
	assetstore "github.com/openshift/installer/pkg/asset/store"
	"github.com/openshift/installer/pkg/terraform"
	gatheraws "github.com/openshift/installer/pkg/terraform/gather/aws"
	gatherlibvirt "github.com/openshift/installer/pkg/terraform/gather/libvirt"
	gatheropenstack "github.com/openshift/installer/pkg/terraform/gather/openstack"
	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	libvirttypes "github.com/openshift/installer/pkg/types/libvirt"
	openstacktypes "github.com/openshift/installer/pkg/types/openstack"
)

func newGatherCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gather",
		Short: "Gather debugging data for a given installation failure",
		Long: `Gather debugging data for a given installation failure.

When installation for Openshift cluster fails, gathering all the data useful for debugging can
become a difficult task. This command helps users to collect the most relevant information that can be used
to debug the installation failures`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(newGatherBootstrapCmd())
	return cmd
}

var (
	gatherBootstrapOpts struct {
		bootstrap string
		masters   []string
	}
)

func newGatherBootstrapCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bootstrap",
		Short: "Gather debugging data for a failing-to-bootstrap control plane",
		Args:  cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			cleanup := setupFileHook(rootOpts.dir)
			defer cleanup()
			err := runGatherBootstrapCmd(rootOpts.dir)
			if err != nil {
				logrus.Fatal(err)
			}
		},
	}
	cmd.PersistentFlags().StringVar(&gatherBootstrapOpts.bootstrap, "bootstrap", "", "Hostname or IP of the bootstrap host")
	cmd.PersistentFlags().StringArrayVar(&gatherBootstrapOpts.masters, "master", []string{}, "Hostnames or IPs of all control plane hosts")
	return cmd
}

func runGatherBootstrapCmd(directory string) error {
	tfStateFilePath := filepath.Join(directory, terraform.StateFileName)
	_, err := os.Stat(tfStateFilePath)
	if os.IsNotExist(err) {
		return unSupportedPlatformGather()
	}
	if err != nil {
		return err
	}

	assetStore, err := assetstore.NewStore(directory)
	if err != nil {
		return errors.Wrap(err, "failed to create asset store")
	}

	config := &installconfig.InstallConfig{}
	if err := assetStore.Fetch(config); err != nil {
		return errors.Wrapf(err, "failed to fetch %s", config.Name())
	}

	tfstate, err := terraform.ReadState(tfStateFilePath)
	if err != nil {
		return errors.Wrapf(err, "failed to read state from %q", tfStateFilePath)
	}
	bootstrap, masters, err := extractHostAddresses(config.Config, tfstate)
	if err != nil {
		if err2, ok := err.(errUnSupportedGatherPlatform); ok {
			logrus.Error(err2)
			return unSupportedPlatformGather()
		}
		return errors.Wrapf(err, "failed to get bootstrap and control plane host addresses from %q", tfStateFilePath)
	}

	logGatherBootstrap(bootstrap, masters)
	return nil
}

func logGatherBootstrap(bootstrap string, masters []string) {
	if s, ok := os.LookupEnv("SSH_AUTH_SOCK"); !ok || s == "" {
		logrus.Info("Make sure ssh-agent is running, env SSH_AUTH_SOCK is set to the ssh-agent's UNIX socket and your private key is added to the agent.")
	}
	logrus.Info("Use the following commands to gather logs from the cluster")
	logrus.Infof("ssh -A core@%s '/usr/local/bin/installer-gather.sh %s'", bootstrap, strings.Join(masters, " "))
	logrus.Infof("scp core@%s:~/log-bundle.tar.gz .", bootstrap)
}

func extractHostAddresses(config *types.InstallConfig, tfstate *terraform.State) (bootstrap string, masters []string, err error) {
	switch config.Platform.Name() {
	case awstypes.Name:
		bootstrap, err = gatheraws.BootstrapIP(tfstate)
		if err != nil {
			return bootstrap, masters, err
		}
		masters, err = gatheraws.ControlPlaneIPs(tfstate)
		if err != nil {
			logrus.Error(err)
		}
	case libvirttypes.Name:
		bootstrap, err = gatherlibvirt.BootstrapIP(tfstate)
		if err != nil {
			return bootstrap, masters, err
		}
		masters, err = gatherlibvirt.ControlPlaneIPs(tfstate)
		if err != nil {
			logrus.Error(err)
		}
	case openstacktypes.Name:
		bootstrap, err = gatheropenstack.BootstrapIP(tfstate)
		if err != nil {
			return bootstrap, masters, err
		}
		masters, err = gatheropenstack.ControlPlaneIPs(tfstate)
		if err != nil {
			logrus.Error(err)
		}
	default:
		return "", nil, errUnSupportedGatherPlatform{Message: fmt.Sprintf("Cannot fetch the bootstrap and control plane host addresses from state file for %s platform", config.Platform.Name())}
	}
	return bootstrap, masters, nil
}

type errUnSupportedGatherPlatform struct {
	Message string
}

func (e errUnSupportedGatherPlatform) Error() string {
	return e.Message
}

func unSupportedPlatformGather() error {
	if gatherBootstrapOpts.bootstrap == "" || len(gatherBootstrapOpts.masters) == 0 {
		return errors.New("boostrap host address and at least one control plane host address must be provided")
	}

	logGatherBootstrap(gatherBootstrapOpts.bootstrap, gatherBootstrapOpts.masters)
	return nil
}
