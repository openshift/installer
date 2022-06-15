package agent

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"
)

// func WaitFor() error {
// 	logrus.Info("WaitFor command")

// 	return nil
// }

// TODO(lranjbar)[AGENT-172]: Add wait for cluster validation
func WaitForClusterValidationSuccess() error {
	logrus.Info("agentWaitForValidationSuccess")

	// zeroClient, err := NewNodeZeroClient()
	// if err != nil {
	// 	return err
	// }

	// // Wait to see assisted service API for the first time
	// WaitForNodeZeroAgentRestAPIInit(zeroClient, 5)

	// clusterZero, err := NewClusterZero(zeroClient)
	// if err != nil {
	// 	return err
	// }

	// // Wait for cluster validations to succeed
	// WaitForClusterZeroManifestsToValidate(clusterZero, 5)

	return nil
}

// Wait for the bootstrap complete triggered by the agent installer.
func WaitForBootstrapComplete() error {
	logrus.Info("WaitForBootstrapComplete")

	zeroClient, err := NewNodeZeroClient()
	if err != nil {
		return err
	}

	// Wait to see assisted service API for the first time
	WaitForNodeZeroAgentRestAPIInit(zeroClient, 5)

	// Research notes: In installer main package create.go:
	// waitForBootstrapComplete(), waitForBootstrapConfigMap()

	// TODO(lranjbar)[AGENT-172]: Add wait for cluster validation
	// clusterZero, err := NewClusterZero(zeroClient)
	// if err != nil {
	// 	return err
	// }
	// Wait for cluster validations to succeed
	// WaitForClusterZeroManifestsToValidate(clusterZero, 5)

	// Wait for cluster Kube API to come up and kubeconfig to be created
	// WaitForClusterZeroKubeAPILive()

	// Wait for bootstrap configmap

	// Wait for bootstrap node to reboot
	// WaitForNodeZeroReboot()

	return nil
}

// Wait for the installation complete triggered by the agent installer.
func WaitForInstallComplete() error {
	logrus.Info("WaitForInstallComplete")

	// Research notes: In installer main package create.go:
	// waitForInitializedCluster()

	// TODO(lranjbar): Add wait for install complete in AGENT-173
	// WaitForBootstrapComplete()

	return nil
}

func WaitForNodeZeroAgentRestAPIInit(zeroClient *nodeZeroClient, timeoutMins int) error {
	logrus.Info("WaitForNodeZeroAgentRestAPIInit")

	timeout := time.Duration(timeoutMins) * time.Minute
	serviceContext, cancel := context.WithTimeout(zeroClient.ctx, timeout)
	defer cancel()

	wait.Until(func() {
		live, err := zeroClient.isAgentAPILive()
		if live && err == nil {
			logrus.Info("Node Zero Agent API Initialized")
			cancel()
		}
	}, 5*time.Second, serviceContext.Done())

	return nil
}

func WaitForClusterZeroManifestsToValidate(clusterZero *clusterZero, timeoutMins int) error {
	logrus.Info("WaitForClusterZeroManifestsToValidate")

	timeout := time.Duration(timeoutMins) * time.Minute
	serviceContext, cancel := context.WithTimeout(clusterZero.zeroClient.ctx, timeout)
	defer cancel()

	wait.Until(func() {
		clusterProgress, _ := clusterZero.get()
		validate, err := clusterZero.parseValidationInfo(clusterProgress)
		if validate && err == nil {
			cancel()
		}
	}, 5*time.Second, serviceContext.Done())

	return nil
}

func WaitForClusterZeroKubeAPILive(timeoutMins int) error {

	return nil
}

// TODO(lranjbar): How to detect?
func WaitForNodeZeroReboot(timeoutMins int) error {

	return nil
}
