package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"

	configv1 "github.com/openshift/api/config/v1"
	configclient "github.com/openshift/client-go/config/clientset/versioned"
	"github.com/openshift/installer/cmd/openshift-install/command"
	"github.com/openshift/installer/pkg/asset/installconfig"
	assetstore "github.com/openshift/installer/pkg/asset/store"
	"github.com/openshift/installer/pkg/asset/tls"
	"github.com/openshift/installer/pkg/clusterapi"
	serialgather "github.com/openshift/installer/pkg/gather"
	"github.com/openshift/installer/pkg/gather/service"
	"github.com/openshift/installer/pkg/gather/ssh"
	"github.com/openshift/installer/pkg/infrastructure"
	infra "github.com/openshift/installer/pkg/infrastructure/platform"

	_ "github.com/openshift/installer/pkg/gather/aws"
	_ "github.com/openshift/installer/pkg/gather/azure"
	_ "github.com/openshift/installer/pkg/gather/gcp"
)

func newGatherCmd(ctx context.Context) *cobra.Command {
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
	cmd.AddCommand(newGatherBootstrapCmd(ctx))
	return cmd
}

var gatherBootstrapOpts struct {
	bootstrap    string
	masters      []string
	sshKeys      []string
	skipAnalysis bool
}

func newGatherBootstrapCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bootstrap",
		Short: "Gather debugging data for a failing-to-bootstrap control plane",
		Args:  cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			cleanup := command.SetupFileHook(command.RootOpts.Dir)
			defer cleanup()
			bundlePath, err := runGatherBootstrapCmd(ctx, command.RootOpts.Dir)
			if err != nil {
				logrus.Fatal(err)
			}

			if !gatherBootstrapOpts.skipAnalysis {
				if err := service.AnalyzeGatherBundle(bundlePath); err != nil {
					logrus.Fatal(err)
				}
			}

			logrus.Infof("Bootstrap gather logs captured here %q", bundlePath)
		},
	}
	cmd.PersistentFlags().StringVar(&gatherBootstrapOpts.bootstrap, "bootstrap", "", "Hostname or IP of the bootstrap host")
	cmd.PersistentFlags().StringArrayVar(&gatherBootstrapOpts.masters, "master", []string{}, "Hostnames or IPs of all control plane hosts")
	cmd.PersistentFlags().StringArrayVar(&gatherBootstrapOpts.sshKeys, "key", []string{}, "Path to SSH private keys that should be used for authentication. If no key was provided, SSH private keys from user's environment will be used")
	cmd.PersistentFlags().BoolVar(&gatherBootstrapOpts.skipAnalysis, "skipAnalysis", false, "Skip analysis of the gathered data")
	return cmd
}

func runGatherBootstrapCmd(ctx context.Context, directory string) (string, error) {
	assetStore, err := assetstore.NewStore(directory)
	if err != nil {
		return "", fmt.Errorf("failed to create asset store: %w", err)
	}
	// add the default bootstrap key pair to the sshKeys list
	bootstrapSSHKeyPair := &tls.BootstrapSSHKeyPair{}
	if err := assetStore.Fetch(ctx, bootstrapSSHKeyPair); err != nil {
		return "", fmt.Errorf("failed to fetch %s: %w", bootstrapSSHKeyPair.Name(), err)
	}
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
	gatherBootstrapOpts.sshKeys = append(gatherBootstrapOpts.sshKeys, tmpfile.Name())

	ha := &infrastructure.HostAddresses{
		Bootstrap: gatherBootstrapOpts.bootstrap,
		Port:      22,
		Masters:   gatherBootstrapOpts.masters,
	}

	if ha.Bootstrap == "" && len(ha.Masters) == 0 {
		config := &installconfig.InstallConfig{}
		if err := assetStore.Fetch(ctx, config); err != nil {
			return "", fmt.Errorf("failed to fetch %s: %w", config.Name(), err)
		}

		provider, err := infra.ProviderForPlatform(config.Config.Platform.Name(), config.Config.EnabledFeatureGates())
		if err != nil {
			return "", fmt.Errorf("error getting infrastructure provider: %w", err)
		}
		if err = provider.ExtractHostAddresses(directory, config.Config, ha); err != nil {
			logrus.Warnf("Failed to extract host addresses: %s", err.Error())
		}
	}

	if ha.Bootstrap == "" {
		return "", errors.New("must provide bootstrap host address")
	}

	return gatherBootstrap(ha.Bootstrap, ha.Port, ha.Masters, directory)
}

func gatherBootstrap(bootstrap string, port int, masters []string, directory string) (string, error) {
	gatherID := time.Now().Format("20060102150405")
	archives := map[string]string{}

	if capiManifestsBundlePath, err := gatherCAPIArtifacts(directory, gatherID); err != nil {
		// Do not fail the whole gather if we can't find capi manifests (we can be running terraform)
		logrus.Infof("Failed to gather Cluster API manifests: %s", err.Error())
	} else {
		archives[capiManifestsBundlePath] = "clusterapi"
	}

	serialLogBundle := filepath.Join(directory, fmt.Sprintf("serial-log-bundle-%s.tar.gz", gatherID))
	serialLogBundlePath, err := filepath.Abs(serialLogBundle)
	if err != nil {
		return "", fmt.Errorf("failed to stat log file: %w", err)
	}

	consoleGather, err := serialgather.New(logrus.StandardLogger(), serialLogBundlePath, bootstrap, masters, directory)
	if err != nil {
		logrus.Infof("Skipping VM console logs gather: %s", err.Error())
	} else {
		logrus.Info("Pulling VM console logs")
		if err := consoleGather.Run(); err != nil {
			logrus.Infof("Failed to gather VM console logs: %s", err.Error())
		} else {
			archives[serialLogBundlePath] = "serial"
		}
	}

	clusterLogBundlePath, err := pullLogsFromBootstrap(gatherID, bootstrap, port, masters, directory)
	if err != nil {
		logrus.Infof("Failed to gather bootstrap logs: %s", err.Error())
	} else {
		archives[clusterLogBundlePath] = ""
	}

	if len(archives) == 0 {
		return "", fmt.Errorf("failed to gather VM console and bootstrap logs")
	}

	logBundlePath := filepath.Join(directory, fmt.Sprintf("log-bundle-%s.tar.gz", gatherID))
	err = serialgather.CombineArchives(logBundlePath, archives)
	if err != nil {
		return "", fmt.Errorf("failed to combine archives: %w", err)
	}

	return logBundlePath, nil
}

func pullLogsFromBootstrap(gatherID string, bootstrap string, port int, masters []string, directory string) (string, error) {
	logrus.Info("Pulling debug logs from the bootstrap machine")
	client, err := ssh.NewClient("core", net.JoinHostPort(bootstrap, strconv.Itoa(port)), gatherBootstrapOpts.sshKeys)
	if err != nil {
		if errors.Is(err, syscall.ECONNREFUSED) || errors.Is(err, syscall.ETIMEDOUT) {
			return "", fmt.Errorf("failed to connect to the bootstrap machine: %w", err)
		}
		return "", fmt.Errorf("failed to create SSH client: %w", err)
	}

	if err := ssh.Run(client, fmt.Sprintf("/usr/local/bin/installer-gather.sh --id %s %s", gatherID, strings.Join(masters, " "))); err != nil {
		return "", fmt.Errorf("failed to run remote command: %w", err)
	}

	file := filepath.Join(directory, fmt.Sprintf("cluster-log-bundle-%s.tar.gz", gatherID))
	if err := ssh.PullFileTo(client, fmt.Sprintf("/home/core/log-bundle-%s.tar.gz", gatherID), file); err != nil {
		return "", fmt.Errorf("failed to pull log file from remote: %w", err)
	}

	clusterLogBundlePath, err := filepath.Abs(file)
	if err != nil {
		return "", fmt.Errorf("failed to stat log file: %w", err)
	}

	return clusterLogBundlePath, nil
}

func logClusterOperatorConditions(ctx context.Context, config *rest.Config) error {
	client, err := configclient.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("creating a config client: %w", err)
	}

	operators, err := client.ConfigV1().ClusterOperators().List(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("listing ClusterOperator objects: %w", err)
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
			if condition.Type == configv1.OperatorAvailable || condition.Type == configv1.OperatorDegraded {
				logrus.Errorf("Cluster operator %s %s is %s with %s: %s", operator.ObjectMeta.Name, condition.Type, condition.Status, condition.Reason, condition.Message)
			} else {
				logrus.Infof("Cluster operator %s %s is %s with %s: %s", operator.ObjectMeta.Name, condition.Type, condition.Status, condition.Reason, condition.Message)
			}
		}
	}

	return nil
}

func gatherCAPIArtifacts(directory, gatherID string) (string, error) {
	logrus.Infoln("Pulling Cluster API artifacts")
	dir, err := filepath.Abs(directory)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for %s: %w", directory, err)
	}

	capiDir := filepath.Join(dir, clusterapi.ArtifactsDir)
	if _, err := os.Stat(capiDir); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return "", fmt.Errorf("either Cluster API manifests not generated or terraform provision")
		}
		return "", fmt.Errorf("failed to stat Cluster API output directory: %w", err)
	}

	bundleDir := filepath.Join(dir, fmt.Sprintf("capi-artifacts-bundle-%s", gatherID))
	// Symlink the hidden directory so the artifacts are not hidden in the archive
	if err := os.Symlink(capiDir, bundleDir); err != nil {
		return "", fmt.Errorf("failed to copy Cluster API artifacts: %w", err)
	}
	defer os.Remove(bundleDir)

	var capiArtifacts []string
	manifests, err := filepath.Glob(filepath.Join(bundleDir, "*.yaml"))
	if err != nil {
		return "", fmt.Errorf("failed to gather Cluster API manifests: %w", err)
	}
	capiArtifacts = append(capiArtifacts, manifests...)

	logs, err := filepath.Glob(filepath.Join(bundleDir, "*.log"))
	if err != nil {
		return "", fmt.Errorf("failed to gather Cluster API control plane logs: %w", err)
	}
	capiArtifacts = append(capiArtifacts, logs...)

	capiArtifactsBundlePath := fmt.Sprintf("%s.tar.gz", bundleDir)
	if err := serialgather.CreateArchive(capiArtifacts, capiArtifactsBundlePath); err != nil {
		return "", fmt.Errorf("failed to create clusterapi bundle file: %w", err)
	}
	return capiArtifactsBundlePath, nil
}
