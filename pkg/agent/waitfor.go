package agent

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"
)

// WaitForBootstrapComplete Wait for the bootstrap process to complete on
// cluster installations triggered by the agent installer.
func WaitForBootstrapComplete(assetDir string) (*Cluster, error) {

	ctx := context.Background()
	cluster, err := NewCluster(ctx, assetDir)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	previous := time.Now()
	timeout := 60 * time.Minute
	waitContext, cancel := context.WithTimeout(cluster.Ctx, timeout)
	defer cancel()

	var lastErr error
	wait.Until(func() {
		bootstrap, exitOnErr, err := cluster.IsBootstrapComplete()
		if bootstrap && err == nil {
			logrus.Info("cluster bootstrap is complete")
			cancel()
		}

		if err != nil {
			if exitOnErr {
				lastErr = err
				cancel()
			} else {
				logrus.Info(err)
			}
		}

		current := time.Now()
		elapsed := current.Sub(previous)
		elapsedTotal := current.Sub(start)
		if elapsed >= 1*time.Minute {
			logrus.Tracef("elapsed: %s, elapsedTotal: %s", elapsed.String(), elapsedTotal.String())
			previous = current
		}

	}, 2*time.Second, waitContext.Done())

	waitErr := waitContext.Err()
	if waitErr != nil {
		if waitErr == context.Canceled && lastErr != nil {
			return cluster, errors.Wrap(lastErr, "bootstrap process returned error")
		}
		if waitErr == context.DeadlineExceeded {
			return cluster, errors.Wrap(waitErr, "bootstrap process timed out")
		}
	}

	return cluster, nil
}

// WaitForInstallComplete Waits for the cluster installation triggered by the
// agent installer to be complete.
func WaitForInstallComplete(assetDir string) (*Cluster, error) {

	cluster, err := WaitForBootstrapComplete(assetDir)

	if err != nil {
		return cluster, errors.Wrap(err, "error occured during bootstrap process")
	}

	timeout := 90 * time.Minute
	waitContext, cancel := context.WithTimeout(cluster.Ctx, timeout)
	defer cancel()

	wait.Until(func() {
		installed, err := cluster.IsInstallComplete()
		if installed && err == nil {
			logrus.Info("Cluster is installed")
			cancel()
		}

	}, 2*time.Second, waitContext.Done())

	waitErr := waitContext.Err()
	if waitErr != nil && waitErr != context.Canceled {
		if err != nil {
			return cluster, errors.Wrap(err, "Error occurred during installation")
		}
		return cluster, errors.Wrap(waitErr, "Cluster installation timed out")
	}
	return cluster, nil
}
