package agent

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/openshift/installer/pkg/agent/zero"
)

// Wait for the bootstrap complete triggered by the agent installer.
func WaitForBootstrapComplete(assetDir string) (*zero.ClusterZero, error) {
	logrus.Info("WaitForBootstrapComplete")

	ctx := context.Background()
	clusterZero, err := zero.NewClusterZero(ctx, assetDir)
	if err != nil {
		logrus.Warn("Unable to make cluster zero object")
		return nil, err
	}

	timeout := 40 * time.Minute
	waitContext, cancel := context.WithTimeout(clusterZero.Ctx, timeout)
	defer cancel()

	wait.Until(func() {
		bootstrap, err := clusterZero.IsBootstrapComplete()
		if bootstrap && err == nil {
			logrus.Info("Cluster bootstrap is complete")
			cancel()
		}

		if err != nil {
			cancel()
		}
	}, 5*time.Second, waitContext.Done())

	if err != nil {
		return clusterZero, err
	}

	return clusterZero, nil
}

// TODO(lranjbar)[AGENT-173]: Add wait for install complete in AGENT-173
// Wait for the installation complete triggered by the agent installer.
func WaitForInstallComplete(assetDir string) error {
	logrus.Info("WaitForInstallComplete")

	// Research notes: In installer main package create.go:
	// waitForInitializedCluster()

	cluster, err := WaitForBootstrapComplete(assetDir)

	if err != nil {
		return err
	}

	timeout := 40 * time.Minute
	waitContext, cancel := context.WithTimeout(cluster.Ctx, timeout)
	defer cancel()

	wait.Until(func() {
		installed, err := cluster.IsInstallComplete()
		if installed && err == nil {
			logrus.Info("Cluster is installed.")
			cancel()
		}

	}, 5*time.Second, waitContext.Done())

	return nil
}
