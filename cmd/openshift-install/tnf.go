package main

import (
	"context"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	configv1 "github.com/openshift/api/config/v1"
	configclient "github.com/openshift/client-go/config/clientset/versioned"
	"github.com/openshift/installer/cmd/openshift-install/command"
	"github.com/openshift/installer/pkg/fencing"
	"github.com/openshift/installer/pkg/metrics/timer"
)

var tnfOpts struct {
	sshKeys    []string
	restConfig *rest.Config
}

func newTNFCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tnf",
		Short: "Commands for Two Node with Fencing clusters",
		Long: `Utilities for managing and validating Two Node with Fencing clusters.

These commands require a deployed two node cluster with fencing and SSH access
to both control plane nodes.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			runRootCmd(cmd, args)

			kubeconfigPath := filepath.Join(command.RootOpts.Dir, "auth", "kubeconfig")
			config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
			if err != nil {
				logrus.Fatalf("Failed to load kubeconfig from %s: %v", kubeconfigPath, err)
			}
			tnfOpts.restConfig = config

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
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(newTNFValidateFencingCmd(ctx))
	return cmd
}

func newTNFValidateFencingCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate-fencing",
		Short: "Validate fencing configuration and fence both nodes sequentially",
		Long: `Validate fencing on a Two Node with Fencing cluster.

This command connects to both control plane nodes via SSH, runs pre-flight
checks (STONITH, Pacemaker, etcd), then fences each node sequentially and
verifies recovery. This is a DISRUPTIVE operation — nodes will be power-cycled.

Requires SSH access to both nodes as user "core" and a cluster-admin kubeconfig.`,
		Args: cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			timer.StartTimer(timer.TotalTimeElapsed)
			cleanup := command.SetupFileHook(command.RootOpts.Dir)
			defer cleanup()

			kubeClient, err := kubernetes.NewForConfig(tnfOpts.restConfig)
			if err != nil {
				logrus.Fatalf("Failed to create Kubernetes client: %v", err)
			}

			if err := fencing.Run(ctx, fencing.Config{
				KubeClient: kubeClient,
				SSHUser:    "core",
				SSHKeys:    tnfOpts.sshKeys,
			}); err != nil {
				logrus.Fatalf("Fencing validation failed: %v", err)
			}

			timer.StopTimer(timer.TotalTimeElapsed)
			timer.LogSummary()
		},
	}
	cmd.PersistentFlags().StringArrayVar(&tnfOpts.sshKeys, "key", nil,
		"Path to SSH private keys for node access. If not provided, SSH agent or default keys are used.")
	return cmd
}
