package main

import (
	"container/heap"
	"context"
	"fmt"
	"os"
	"path/filepath"
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
	"github.com/openshift/installer/pkg/gather/ssh"
	"github.com/openshift/installer/pkg/terraform"
	gatheraws "github.com/openshift/installer/pkg/terraform/gather/aws"
	gatherazure "github.com/openshift/installer/pkg/terraform/gather/azure"
	gathergcp "github.com/openshift/installer/pkg/terraform/gather/gcp"
	gatherlibvirt "github.com/openshift/installer/pkg/terraform/gather/libvirt"
	gatheropenstack "github.com/openshift/installer/pkg/terraform/gather/openstack"
	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
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
	client, err := ssh.NewClient("core", fmt.Sprintf("%s:%d", bootstrap, port), gatherBootstrapOpts.sshKeys)
	if err != nil && len(gatherBootstrapOpts.sshKeys) == 0 {
		return errors.Wrap(err, "failed to create SSH client, ensure the proper ssh key is in your keyring or specify with --key")
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
	logrus.Infof("Bootstrap gather logs captured here %q", file)
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

type operatorConditionType int

const (
	conditionUnknown     operatorConditionType = -1
	conditionAvailable   operatorConditionType = 0
	conditionDegraded    operatorConditionType = 1
	conditionProgressing operatorConditionType = 2
	conditionUpgradeable operatorConditionType = 3
)

type OperatorCondition struct {
	priority         operatorConditionType
	operatorName     string
	conditionType    configv1.ClusterStatusConditionType
	conditionStatus  configv1.ConditionStatus
	conditionReason  string
	conditionMessage string
}

type OperatorConditions []OperatorCondition

func (o OperatorConditions) Len() int                 { return len(o) }
func (o OperatorConditions) Less(i, j int) bool       { return o[i].priority > o[j].priority }
func (o OperatorConditions) Swap(i, j int)            { o[i], o[j] = o[j], o[i] }
func (o *OperatorConditions) Push(pushed interface{}) { *o = append(*o, pushed.(OperatorCondition)) }
func (o *OperatorConditions) Pop() (popped interface{}) {
	popped = (*o)[len(*o)-1]
	*o = (*o)[:len(*o)-1]
	return
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

	m := make(map[string]*OperatorConditions)
	for _, operator := range operators.Items {
		for _, condition := range operator.Status.Conditions {
			if condition.Reason == "" {
				condition.Reason = "noreason"
			}
			if condition.Message == "" {
				condition.Message = "nomessage"
			}

			switch condition.Type {
			case configv1.OperatorAvailable:
				if condition.Status == configv1.ConditionTrue {
					continue
				}
			case configv1.OperatorDegraded:
				if condition.Status == configv1.ConditionFalse {
					continue
				}
				oc := OperatorCondition{conditionDegraded, operator.ObjectMeta.Name,
					condition.Type, condition.Status, condition.Reason, condition.Message}
				if m[operator.ObjectMeta.Name] == nil {
					heap.Init(m[operator.ObjectMeta.Name])
				}
				heap.Push(m[operator.ObjectMeta.Name], oc)
			case configv1.OperatorProgressing:
				if condition.Status == configv1.ConditionFalse {
					continue
				}
				oc := OperatorCondition{conditionProgressing, operator.ObjectMeta.Name,
					condition.Type, condition.Status, condition.Reason, condition.Message}
				if m[operator.ObjectMeta.Name] == nil {
					heap.Init(m[operator.ObjectMeta.Name])
				}
				heap.Push(m[operator.ObjectMeta.Name], oc)
			default:
				oc := OperatorCondition{conditionUnknown, operator.ObjectMeta.Name,
					condition.Type, condition.Status, condition.Reason, condition.Message}
				if m[operator.ObjectMeta.Name] == nil {
					heap.Init(m[operator.ObjectMeta.Name])
				}
				heap.Push(m[operator.ObjectMeta.Name], oc)
				continue
			}
			oc := OperatorCondition{conditionUnknown, operator.ObjectMeta.Name,
				condition.Type, condition.Status, condition.Reason, condition.Message}
			if m[operator.ObjectMeta.Name] == nil {
				heap.Init(m[operator.ObjectMeta.Name])
			}
			heap.Push(m[operator.ObjectMeta.Name], oc)
		}
	}

	for key, _ := range m {
		logfunc := logrus.Infof
		err := false

		for m[key].Len() != 0 {
			oc := heap.Pop(m[key]).(OperatorCondition)
			switch oc.priority {
			case conditionDegraded:
				logfunc = logrus.Errorf
				err = true
			case conditionProgressing:
				logfunc = logrus.Debugf
			case conditionUnknown:
				logfunc = logrus.Infof
				if err {
					logfunc = logrus.Debugf
				}
			}

			logfunc("Cluster-operator %s has %s=%s with %s: %s", oc.operatorName,
				oc.conditionType, oc.conditionStatus, oc.conditionReason, oc.conditionMessage)
		}
	}

	return nil
}
