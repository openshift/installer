package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	configv1 "github.com/openshift/api/config/v1"
	configclient "github.com/openshift/client-go/config/clientset/versioned"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"

	"github.com/openshift/installer/pkg/asset/installconfig"
	assetstore "github.com/openshift/installer/pkg/asset/store"
	"github.com/openshift/installer/pkg/asset/tls"
	"github.com/openshift/installer/pkg/gather/ssh"
	"github.com/openshift/installer/pkg/terraform"
	gatheraws "github.com/openshift/installer/pkg/terraform/gather/aws"
	gatherazure "github.com/openshift/installer/pkg/terraform/gather/azure"
	gatherbaremetal "github.com/openshift/installer/pkg/terraform/gather/baremetal"
	gathergcp "github.com/openshift/installer/pkg/terraform/gather/gcp"
	gatherlibvirt "github.com/openshift/installer/pkg/terraform/gather/libvirt"
	gatheropenstack "github.com/openshift/installer/pkg/terraform/gather/openstack"
	gatherovirt "github.com/openshift/installer/pkg/terraform/gather/ovirt"
	gathervsphere "github.com/openshift/installer/pkg/terraform/gather/vsphere"
	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
	baremetaltypes "github.com/openshift/installer/pkg/types/baremetal"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
	libvirttypes "github.com/openshift/installer/pkg/types/libvirt"
	openstacktypes "github.com/openshift/installer/pkg/types/openstack"
	ovirttypes "github.com/openshift/installer/pkg/types/ovirt"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
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
		sshKeys   []string
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
	cmd.PersistentFlags().StringArrayVar(&gatherBootstrapOpts.sshKeys, "key", []string{}, "Path to SSH private keys that should be used for authentication. If no key was provided, SSH private keys from user's environment will be used")
	return cmd
}

func runGatherBootstrapCmd(directory string) error {
	assetStore, err := assetstore.NewStore(directory)
	if err != nil {
		return errors.Wrap(err, "failed to create asset store")
	}
	// add the default bootstrap key pair to the sshKeys list
	bootstrapSSHKeyPair := &tls.BootstrapSSHKeyPair{}
	if err := assetStore.Fetch(bootstrapSSHKeyPair); err != nil {
		return errors.Wrapf(err, "failed to fetch %s", bootstrapSSHKeyPair.Name())
	}
	tmpfile, err := ioutil.TempFile("", "bootstrap-ssh")
	if err != nil {
		return err
	}
	defer os.Remove(tmpfile.Name())
	if _, err := tmpfile.Write(bootstrapSSHKeyPair.Private()); err != nil {
		return err
	}
	if err := tmpfile.Close(); err != nil {
		return err
	}
	gatherBootstrapOpts.sshKeys = append(gatherBootstrapOpts.sshKeys, tmpfile.Name())

	tfStateFilePath := filepath.Join(directory, terraform.StateFileName)
	_, err = os.Stat(tfStateFilePath)
	if os.IsNotExist(err) {
		return unSupportedPlatformGather(directory)
	}
	if err != nil {
		return err
	}

	config := &installconfig.InstallConfig{}
	if err := assetStore.Fetch(config); err != nil {
		return errors.Wrapf(err, "failed to fetch %s", config.Name())
	}

	tfstate, err := terraform.ReadState(tfStateFilePath)
	if err != nil {
		return errors.Wrapf(err, "failed to read state from %q", tfStateFilePath)
	}
	bootstrap, port, masters, err := extractHostAddresses(config.Config, tfstate)
	if err != nil {
		if err2, ok := err.(errUnSupportedGatherPlatform); ok {
			logrus.Error(err2)
			return unSupportedPlatformGather(directory)
		}
		return errors.Wrapf(err, "failed to get bootstrap and control plane host addresses from %q", tfStateFilePath)
	}

	return logGatherBootstrap(bootstrap, port, masters, directory)
}

func logGatherBootstrap(bootstrap string, port int, masters []string, directory string) error {
	logrus.Info("Pulling debug logs from the bootstrap machine")
	client, err := ssh.NewClient("core", net.JoinHostPort(bootstrap, strconv.Itoa(port)), gatherBootstrapOpts.sshKeys)
	if err != nil && strings.Contains(err.Error(), "ssh: handshake failed: ssh: unable to authenticate") {
		return errors.Wrap(err, "failed to create SSH client, ensure the private key is added to your authentication agent (ssh-agent) or specified with the --key parameter")
	} else if err != nil {
		return errors.Wrap(err, "failed to create SSH client")
	}
	gatherID := time.Now().Format("20060102150405")
	if err := ssh.Run(client, fmt.Sprintf("/usr/local/bin/installer-gather.sh --id %s %s", gatherID, strings.Join(masters, " "))); err != nil {
		return errors.Wrap(err, "failed to run remote command")
	}
	file := filepath.Join(directory, fmt.Sprintf("log-bundle-%s.tar.gz", gatherID))
	if err := ssh.PullFileTo(client, fmt.Sprintf("/home/core/log-bundle-%s.tar.gz", gatherID), file); err != nil {
		return errors.Wrap(err, "failed to pull log file from remote")
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return errors.Wrap(err, "failed to stat log file")
	}
	logrus.Infof("Bootstrap gather logs captured here %q", path)
	return nil
}

func extractHostAddresses(config *types.InstallConfig, tfstate *terraform.State) (bootstrap string, port int, masters []string, err error) {
	port = 22
	switch config.Platform.Name() {
	case awstypes.Name:
		bootstrap, err = gatheraws.BootstrapIP(tfstate)
		if err != nil {
			return bootstrap, port, masters, err
		}
		masters, err = gatheraws.ControlPlaneIPs(tfstate)
		if err != nil {
			logrus.Error(err)
		}
	case azuretypes.Name:
		bootstrap, err = gatherazure.BootstrapIP(tfstate)
		if err != nil {
			return bootstrap, port, masters, err
		}
		masters, err = gatherazure.ControlPlaneIPs(tfstate)
		if err != nil {
			logrus.Error(err)
		}
	case baremetaltypes.Name:
		bootstrap = config.Platform.BareMetal.BootstrapProvisioningIP
		masters, err = gatherbaremetal.ControlPlaneIPs(config, tfstate)
		if err != nil {
			return bootstrap, port, masters, err
		}
	case gcptypes.Name:
		bootstrap, err = gathergcp.BootstrapIP(tfstate)
		if err != nil {
			return bootstrap, port, masters, err
		}
		masters, err = gathergcp.ControlPlaneIPs(tfstate)
		if err != nil {
			logrus.Error(err)
		}
	case libvirttypes.Name:
		bootstrap, err = gatherlibvirt.BootstrapIP(tfstate)
		if err != nil {
			return bootstrap, port, masters, err
		}
		masters, err = gatherlibvirt.ControlPlaneIPs(tfstate)
		if err != nil {
			logrus.Error(err)
		}
	case openstacktypes.Name:
		bootstrap, err = gatheropenstack.BootstrapIP(tfstate)
		if err != nil {
			return bootstrap, port, masters, err
		}
		masters, err = gatheropenstack.ControlPlaneIPs(tfstate)
		if err != nil {
			logrus.Error(err)
		}
	case ovirttypes.Name:
		bootstrap, err := gatherovirt.BootstrapIP(tfstate)
		if err != nil {
			return bootstrap, port, masters, err
		}
		masters, err = gatherovirt.ControlPlaneIPs(tfstate)
	case vspheretypes.Name:
		bootstrap, err = gathervsphere.BootstrapIP(config, tfstate)
		if err != nil {
			return bootstrap, port, masters, err
		}
		masters, err = gathervsphere.ControlPlaneIPs(config, tfstate)
		if err != nil {
			logrus.Error(err)
		}
	default:
		return "", port, nil, errUnSupportedGatherPlatform{Message: fmt.Sprintf("Cannot fetch the bootstrap and control plane host addresses from state file for %s platform", config.Platform.Name())}
	}
	return bootstrap, port, masters, nil
}

type errUnSupportedGatherPlatform struct {
	Message string
}

func (e errUnSupportedGatherPlatform) Error() string {
	return e.Message
}

func unSupportedPlatformGather(directory string) error {
	if gatherBootstrapOpts.bootstrap == "" || len(gatherBootstrapOpts.masters) == 0 {
		return errors.New("bootstrap host address and at least one control plane host address must be provided")
	}

	return logGatherBootstrap(gatherBootstrapOpts.bootstrap, 22, gatherBootstrapOpts.masters, directory)
}

func logClusterOperatorConditions(ctx context.Context, config *rest.Config) error {
	client, err := configclient.NewForConfig(config)
	if err != nil {
		return errors.Wrap(err, "creating a config client")
	}

	operators, err := client.ConfigV1().ClusterOperators().List(metav1.ListOptions{})
	if err != nil {
		return errors.Wrap(err, "listing ClusterOperator objects")
	}

	for _, operator := range operators.Items {
		for _, condition := range operator.Status.Conditions {
			if condition.Type == configv1.OperatorUpgradeable {
				continue
			} else if condition.Type == configv1.OperatorAvailable && condition.Status == configv1.ConditionTrue {
				continue
			} else if (condition.Type == configv1.OperatorDegraded || condition.Type == configv1.OperatorProgressing) && condition.Status == configv1.ConditionFalse {
				continue
			}
			if condition.Type == configv1.OperatorDegraded {
				logrus.Errorf("Cluster operator %s %s is %s with %s: %s", operator.ObjectMeta.Name, condition.Type, condition.Status, condition.Reason, condition.Message)
			} else {
				logrus.Infof("Cluster operator %s %s is %s with %s: %s", operator.ObjectMeta.Name, condition.Type, condition.Status, condition.Reason, condition.Message)
			}
		}
	}

	return nil
}
