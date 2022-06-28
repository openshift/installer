package agent

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"
)

// WaitForBootstrapComplete Wait for the bootstrap process to complete on
// cluster installations triggered by the agent installer.
func WaitForBootstrapComplete(assetDir string) (*Cluster, error) {

	ctx := context.Background()
	cluster, err := NewCluster(ctx, assetDir)
	if err != nil {
		logrus.Warn("unable to make cluster object to track installation")
		return nil, err
	}

	timeout := 40 * time.Minute
	waitContext, cancel := context.WithTimeout(cluster.Ctx, timeout)
	defer cancel()

	wait.Until(func() {
		bootstrap, err := cluster.IsBootstrapComplete()
		if bootstrap && err == nil {
			logrus.Info("cluster bootstrap is complete")
			cancel()
		}

		logrus.Info("Cluster is still installing")

	}, 5*time.Second, waitContext.Done())

	waitErr := waitContext.Err()
	if waitErr != nil && waitErr != context.Canceled {
		if err != nil {
			return cluster, err
		}
		return cluster, waitErr
	}

	return cluster, nil
}
