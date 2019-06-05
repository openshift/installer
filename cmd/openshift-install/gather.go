package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/openshift/installer/pkg/asset/installconfig"
	assetstore "github.com/openshift/installer/pkg/asset/store"
	"github.com/openshift/installer/pkg/gather/ssh"
	"github.com/openshift/installer/pkg/terraform"
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
		return unSupportedPlatformGather(directory)
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

	sfRaw, err := ioutil.ReadFile(tfStateFilePath)
	if err != nil {
		return errors.Wrapf(err, "failed to read %q", tfStateFilePath)
	}

	var tfstate terraformState
	if err := json.Unmarshal(sfRaw, &tfstate); err != nil {
		return errors.Wrapf(err, "failed to unmarshal %q", tfStateFilePath)
	}

	bootstrap, masters, err := extractHostAddresses(config.Config, tfstate)
	if err != nil {
		if err2, ok := err.(errUnSupportedGatherPlatform); ok {
			logrus.Error(err2)
			return unSupportedPlatformGather(directory)
		}
		return errors.Wrapf(err, "failed to get bootstrap and control plane host addresses from %q", tfStateFilePath)
	}

	return logGatherBootstrap(bootstrap, port, masters, directory)
}

func logGatherBootstrap(bootstrap string, port int, masters []string, directory string) {
	logrus.Info("Pulling logs from bootstrap for debugging")

	client, err := ssh.NewClient("core", fmt.Sprintf("%s:%d", bootstrap, port), nil)
	if err != nil {
		return errors.Wrap(err, "failed to create SSH client")
	}
	if err := ssh.Run(client, fmt.Sprintf("/usr/local/bin/installer-gather.sh %s", strings.Join(masters, " "))); err != nil {
		return errors.Wrap(err, "failed to run remote command")
	}
	file := filepath.Join(directory, fmt.Sprintf("log-bundle-%s.tar.gz", time.Now().Format("20060102150405")))
	if err := ssh.PullFileTo(client, "/home/core/log-bundle.tar.gz", file); err != nil {
		return errors.Wrap(err, "failed to pull log file from remote")
	}
	logrus.Infof("Bootstrap gather logs captured here %q", file)
	return nil
}

func extractHostAddresses(config *types.InstallConfig, tfstate terraformState) (bootstrap string, masters []string, err error) {
	mcount := *config.ControlPlane.Replicas
	switch config.Platform.Name() {
	case awstypes.Name:
		bm := tfstate.Modules["root/bootstrap"]
		bootstrap, _, err = unstructured.NestedString(bm.Resources["aws_instance.bootstrap"], "primary", "attributes", "public_ip")
		if err != nil {
			return bootstrap, masters, errors.Wrapf(err, "failed to get bootstrap host addresses")
		}

		mm := tfstate.Modules["root/masters"]
		for idx := int64(0); idx < mcount; idx++ {
			r := fmt.Sprintf("aws_instance.master.%d", idx)
			if mcount == 1 {
				r = "aws_instance.master"
			}
			var master string
			master, _, err = unstructured.NestedString(mm.Resources[r], "primary", "attributes", "private_ip")
			if err != nil {
				return bootstrap, masters, errors.Wrapf(err, "failed to get master host addresses")
			}
			masters = append(masters, master)
		}
	case libvirttypes.Name:
		bm := tfstate.Modules["root/bootstrap"]
		bootstrap, _, err = unstructured.NestedString(bm.Resources["libvirt_domain.bootstrap"], "primary", "attributes", "network_interface.0.hostname")
		if err != nil {
			return bootstrap, masters, errors.Wrapf(err, "failed to get bootstrap host addresses")
		}

		rm := tfstate.Modules["root"]
		for idx := int64(0); idx < mcount; idx++ {
			r := fmt.Sprintf("libvirt_domain.master.%d", idx)
			if mcount == 1 {
				r = "libvirt_domain.master"
			}
			var master string
			master, _, err = unstructured.NestedString(rm.Resources[r], "primary", "attributes", "network_interface.0.hostname")
			if err != nil {
				return bootstrap, masters, errors.Wrapf(err, "failed to get master host addresses")
			}
			masters = append(masters, master)
		}
	case openstacktypes.Name:
		bm := tfstate.Modules["root/bootstrap"]
		bootstrap, _, err = unstructured.NestedString(bm.Resources["openstack_compute_instance_v2.bootstrap"], "primary", "attributes", "access_ip_v4")
		if err != nil {
			return bootstrap, masters, errors.Wrapf(err, "failed to get bootstrap host addresses")
		}

		mm := tfstate.Modules["root/masters"]
		for idx := int64(0); idx < mcount; idx++ {
			r := fmt.Sprintf("openstack_compute_instance_v2.master_conf.%d", idx)
			if mcount == 1 {
				r = "openstack_compute_instance_v2.master_conf"
			}
			var master string
			master, _, err = unstructured.NestedString(mm.Resources[r], "primary", "attributes", "access_ip_v4")
			if err != nil {
				return bootstrap, masters, errors.Wrapf(err, "failed to get master host addresses")
			}
			masters = append(masters, master)
		}
	default:
		return "", nil, errUnSupportedGatherPlatform{Message: fmt.Sprintf("Cannot fetch the bootstrap and control plane host addresses from state file for %s platform", config.Platform.Name())}
	}
	return bootstrap, masters, nil
}

type terraformState struct {
	Modules map[string]terraformStateModule
}

type terraformStateModule struct {
	Resources map[string]map[string]interface{} `json:"resources"`
}

func (tfs *terraformState) UnmarshalJSON(raw []byte) error {
	var transform struct {
		Modules []struct {
			Path []string `json:"path"`
			terraformStateModule
		} `json:"modules"`
	}
	if err := json.Unmarshal(raw, &transform); err != nil {
		return err
	}
	if tfs == nil {
		tfs = &terraformState{}
	}
	if tfs.Modules == nil {
		tfs.Modules = make(map[string]terraformStateModule)
	}
	for _, m := range transform.Modules {
		tfs.Modules[strings.Join(m.Path, "/")] = terraformStateModule{Resources: m.Resources}
	}
	return nil
}

type errUnSupportedGatherPlatform struct {
	Message string
}

func (e errUnSupportedGatherPlatform) Error() string {
	return e.Message
}

func unSupportedPlatformGather(directory string) error {
	if gatherBootstrapOpts.bootstrap == "" || len(gatherBootstrapOpts.masters) == 0 {
		return errors.New("boostrap host address and at least one control plane host address must be provided")
	}

	return logGatherBootstrap(gatherBootstrapOpts.bootstrap, 22, gatherBootstrapOpts.masters, directory)
}
