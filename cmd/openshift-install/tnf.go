package main

import (
	"context"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	configv1 "github.com/openshift/api/config/v1"
	configclient "github.com/openshift/client-go/config/clientset/versioned"
	"github.com/openshift/installer/cmd/openshift-install/command"
	"github.com/openshift/installer/pkg/fencing"
	"github.com/openshift/installer/pkg/metrics/timer"
)

var tnfValidateOpts struct {
	sshKeys []string
}

func newTNFValidateFencingCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tnf-validate-fencing",
		Short: "Validate fencing by power-cycling both nodes sequentially (DISRUPTIVE)",
		Long: `Validate fencing on a Two Node with Fencing cluster.

This command connects to both control plane nodes via SSH, runs pre-flight
checks (STONITH, Pacemaker, etcd), then fences each node sequentially and
verifies recovery.

WARNING: This is a DISRUPTIVE operation — nodes will be forcibly powered off
via STONITH and must recover automatically.

Requires SSH access to both nodes as user "core" and a cluster-admin kubeconfig.`,
		Args: cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			runRootCmd(cmd, args)
			timer.StartTimer(timer.TotalTimeElapsed)
			cleanup := command.SetupFileHook(command.RootOpts.Dir)
			defer cleanup()

			kubeconfigPath := filepath.Join(command.RootOpts.Dir, "auth", "kubeconfig")
			config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
			if err != nil {
				logrus.Fatalf("Failed to load kubeconfig from %s: %v", kubeconfigPath, err)
			}

			cfgClient, err := configclient.NewForConfig(config)
			if err != nil {
				logrus.Fatalf("Failed to create config client: %v", err)
			}
			infra, err := cfgClient.ConfigV1().Infrastructures().Get(ctx, "cluster", metav1.GetOptions{})
			if err != nil {
				logrus.Fatalf("Failed to read Infrastructure CR: %v", err)
			}
			if infra.Status.ControlPlaneTopology != configv1.DualReplicaTopologyMode {
				logrus.Fatalf("This command requires a Two Node with Fencing (DualReplica) cluster, found %q", infra.Status.ControlPlaneTopology)
			}

			kubeClient, err := kubernetes.NewForConfig(config)
			if err != nil {
				logrus.Fatalf("Failed to create Kubernetes client: %v", err)
			}

			if err := fencing.Run(ctx, fencing.Config{
				KubeClient: kubeClient,
				SSHUser:    "core",
				SSHKeys:    tnfValidateOpts.sshKeys,
			}); err != nil {
				logrus.Fatalf("Fencing validation failed: %v", err)
			}

			timer.StopTimer(timer.TotalTimeElapsed)
			timer.LogSummary()
		},
	}
	cmd.PersistentFlags().StringArrayVar(&tnfValidateOpts.sshKeys, "key", nil,
		"Path to SSH private keys for node access. If not provided, SSH agent or default keys are used.")
	return cmd
}
