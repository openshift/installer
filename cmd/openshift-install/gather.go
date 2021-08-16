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
	"syscall"
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
	"github.com/openshift/installer/pkg/gather/service"
	"github.com/openshift/installer/pkg/gather/ssh"
	"github.com/openshift/installer/pkg/metrics/gatherer"
	timer "github.com/openshift/installer/pkg/metrics/timer"
	platformstages "github.com/openshift/installer/pkg/terraform/stages/platform"
)

func newGatherCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gather",
		Short: "Gather debugging data for a given installation failure",
		Long: `Gather debugging data for a given installation failure.

When installation for OpenShift cluster fails, gathering all the data useful for debugging can
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
		bootstrap    string
		masters      []string
		sshKeys      []string
		skipAnalysis bool
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
			gatherer.InitializeInvocationMetrics(gatherer.GatherMetricName)
			timer.StartTimer(timer.TotalTimeElapsed)
			bundlePath, err := runGatherBootstrapCmd(rootOpts.dir)
			if err != nil {
				gatherer.LogError("failed", gatherer.CurrentInvocationContext)
				logrus.Fatal(err)
			}
			if !gatherBootstrapOpts.skipAnalysis {
				if err := service.AnalyzeGatherBundle(bundlePath); err != nil {
					logrus.Fatal(err)
				}
			}
			timer.StopTimer(timer.TotalTimeElapsed)
		},
	}
	cmd.PersistentFlags().StringVar(&gatherBootstrapOpts.bootstrap, "bootstrap", "", "Hostname or IP of the bootstrap host")
	cmd.PersistentFlags().StringArrayVar(&gatherBootstrapOpts.masters, "master", []string{}, "Hostnames or IPs of all control plane hosts")
	cmd.PersistentFlags().StringArrayVar(&gatherBootstrapOpts.sshKeys, "key", []string{}, "Path to SSH private keys that should be used for authentication. If no key was provided, SSH private keys from user's environment will be used")
	cmd.PersistentFlags().BoolVar(&gatherBootstrapOpts.skipAnalysis, "skipAnalysis", false, "Skip analysis of the gathered data")
	return cmd
}

func runGatherBootstrapCmd(directory string) (string, error) {
	assetStore, err := assetstore.NewStore(directory)
	if err != nil {
		return "", errors.Wrap(err, "failed to create asset store")
	}
	// add the default bootstrap key pair to the sshKeys list
	bootstrapSSHKeyPair := &tls.BootstrapSSHKeyPair{}
	if err := assetStore.Fetch(bootstrapSSHKeyPair); err != nil {
		return "", errors.Wrapf(err, "failed to fetch %s", bootstrapSSHKeyPair.Name())
	}
	tmpfile, err := ioutil.TempFile("", "bootstrap-ssh")
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
	gatherBootstrapOpts.sshKeys = append(gatherBootstrapOpts.sshKeys, tmpfile.Name())

	bootstrap := gatherBootstrapOpts.bootstrap
	port := 22
	masters := gatherBootstrapOpts.masters
	if bootstrap == "" && len(masters) == 0 {
		config := &installconfig.InstallConfig{}
		if err := assetStore.Fetch(config); err != nil {
			return "", errors.Wrapf(err, "failed to fetch %s", config.Name())
		}

		for _, stage := range platformstages.StagesForPlatform(config.Config.Platform.Name()) {
			stageBootstrap, stagePort, stageMasters, err := stage.ExtractHostAddresses(directory, config.Config)
			if err != nil {
				return "", err
			}
			if stageBootstrap != "" {
				bootstrap = stageBootstrap
			}
			if stagePort != 0 {
				port = stagePort
			}
			if len(stageMasters) > 0 {
				masters = stageMasters
			}
		}

		if bootstrap == "" || len(masters) == 0 {
			return "", errors.New("bootstrap host address and at least one control plane host address must be provided")
		}
	} else if bootstrap == "" || len(masters) == 0 {
		return "", errors.New("must provide both bootstrap host address and at least one control plane host address when providing one")
	}

	return logGatherBootstrap(bootstrap, port, masters, directory)
}

func logGatherBootstrap(bootstrap string, port int, masters []string, directory string) (string, error) {
	logrus.Info("Pulling debug logs from the bootstrap machine")
	client, err := ssh.NewClient("core", net.JoinHostPort(bootstrap, strconv.Itoa(port)), gatherBootstrapOpts.sshKeys)
	if err != nil {
		if errors.Is(err, syscall.ECONNREFUSED) {
			return "", errors.Wrap(err, "failed to connect to the bootstrap machine")
		}
		return "", errors.Wrap(err, "failed to create SSH client")
	}

	gatherID := time.Now().Format("20060102150405")
	if err := ssh.Run(client, fmt.Sprintf("/usr/local/bin/installer-gather.sh --id %s %s", gatherID, strings.Join(masters, " "))); err != nil {
		return "", errors.Wrap(err, "failed to run remote command")
	}
	file := filepath.Join(directory, fmt.Sprintf("log-bundle-%s.tar.gz", gatherID))
	if err := ssh.PullFileTo(client, fmt.Sprintf("/home/core/log-bundle-%s.tar.gz", gatherID), file); err != nil {
		return "", errors.Wrap(err, "failed to pull log file from remote")
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", errors.Wrap(err, "failed to stat log file")
	}
	logrus.Infof("Bootstrap gather logs captured here %q", path)
	return path, nil
}

func logClusterOperatorConditions(ctx context.Context, config *rest.Config) error {
	client, err := configclient.NewForConfig(config)
	if err != nil {
		return errors.Wrap(err, "creating a config client")
	}

	operators, err := client.ConfigV1().ClusterOperators().List(ctx, metav1.ListOptions{})
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
