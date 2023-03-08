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
func WaitForBootstrapComplete(cluster *Cluster) error {
	start := time.Now()
	previous := time.Now()
	timeout := 60 * time.Minute
	waitContext, cancel := context.WithTimeout(cluster.Ctx, timeout)
	defer cancel()

	var lastErr, lastErrOnExit error
	wait.Until(func() {
		bootstrap, exitOnErr, err := cluster.IsBootstrapComplete()
		if bootstrap && err == nil {
			logrus.Info("cluster bootstrap is complete")
			cancel()
		}

		if err != nil {
			if exitOnErr {
				lastErrOnExit = err
				cancel()
			} else {
				if errors.Is(err, lastErr) {
					logrus.Info(err)
					lastErr = err
				}
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
		if errors.Is(waitErr, context.Canceled) && lastErrOnExit != nil {
			return errors.Wrap(lastErrOnExit, "bootstrap process returned error")
		}
		if waitErr == context.DeadlineExceeded {
			return errors.Wrap(waitErr, "bootstrap process timed out")
		}
	}

	return nil
}

// WaitForInstallComplete Waits for the cluster installation triggered by the
// agent installer to be complete.
func WaitForInstallComplete(cluster *Cluster) error {
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
		return errors.Wrap(waitErr, "Cluster installation timed out")
	}
	return nil
}
