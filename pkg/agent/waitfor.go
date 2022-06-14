package agent

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"k8s.io/apimachinery/pkg/util/wait"
)

// WaitFor wait for the installation complete triggered by the agent installer.
func WaitFor() error {
	logrus.Info("WaitFor command")

	// ctx := context.Background()
	// zeroClient, err := NewNodeZeroClient()
	// if err != nil {
	// 	return err
	// }

	return nil
}

func WaitForNodeZeroAgentRestAPIInit(ctx context.Context, zeroClient *nodeZeroClient, timeoutMins int) error {
	logrus.Info("WaitForNodeZeroAgentRestAPIInit")

	timeout := time.Duration(timeoutMins) * time.Minute
	serviceContext, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	wait.Until(func() {
		live, err := isAgentAPILive(ctx, zeroClient)
		if live && err == nil {
			logrus.Info("Node Zero Agent API Initialized")
			cancel()
		}
	}, 5*time.Second, serviceContext.Done())

	return nil
}

func WaitForClusterValidationSuccess(ctx context.Context, zeroClient *nodeZeroClient, timeoutMins int) error {
	logrus.Info("agentWaitForValidationSuccess")

	// timeout := time.Duration(timeoutMins) * time.Minute
	// serviceContext, cancel := context.WithTimeout(ctx, timeout)
	// defer cancel()

	// TODO(lranjbar): Update for validations info checking
	// wait.Until(func() {
	// 	clusterZero, err := getClusterZero(zeroClient, ctx)
	// 	if clusterZero.ValidationInfo && err == nil {
	// 		cancel()
	// 	}
	// }, 5*time.Second, serviceContext.Done())

	return nil
}

func WaitForKubeAPI(ctx context.Context, zeroClient *nodeZeroClient, timeoutMins int) error {

	return nil
}

func WaitForNodeZeroReboot(ctx context.Context, zeroClient *nodeZeroClient, timeoutMins int) error {

	return nil
}

func WaitForBootstrapComplete(ctx context.Context, zeroClient *nodeZeroClient, timeoutMins int) error {
	logrus.Info("WaitForBootstrapComplete")

	// Wait to see assisted service API for the first time
	WaitForNodeZeroAgentRestAPIInit(ctx, zeroClient, 5)

	// Wait for cluster validations to succeed
	// WaitForClusterValidationSuccess(ctx, zeroClient, 1)

	// Wait for kubeAPI to come up and kubeconfig to be created
	// WaitForKubeAPI(ctx, zeroClient, 20)

	// TODO(lranjbar): How to detect?
	// Wait for bootstrap node to reboot
	// WaitForNodeZeroReboot(ctx, zeroClient, 20)

	return nil
}

func WaitForInstallComplete() error {
	logrus.Info("WaitForInstallComplete")

	// WaitForBootstrapComplete()

	return nil
}

func isKubeAPILive() (bool, error) {

	return false, nil
}

func doesKubeConfigExist() (bool, error) {

	return false, nil
}

func printInstallStatus() error {

	return nil
}
