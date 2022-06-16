package agent

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/openshift/installer/pkg/agent/zero"
)

// TODO(lranjbar)[AGENT-172]: Add wait for cluster validation
func WaitForClusterValidationSuccess(assetDir string) error {
	logrus.Info("agentWaitForValidationSuccess")

	ctx := context.Background()
	clusterZero, err := zero.NewClusterZero(ctx, assetDir)
	if err != nil {
		return err
	}

	// // Wait to see assisted service API for the first time
	WaitForNodeZeroAgentRestAPIInit(clusterZero, 5)

	// // Wait for cluster validations to succeed
	// WaitForClusterZeroManifestsToValidate(clusterZero, 5)

	return nil
}

// Wait for the bootstrap complete triggered by the agent installer.
func WaitForBootstrapComplete(assetDir string) error {
	logrus.Info("WaitForBootstrapComplete")

	ctx := context.Background()
	clusterZero, err := zero.NewClusterZero(ctx, assetDir)
	if err != nil {
		return err
	}

	// Wait to see assisted service API for the first time
	WaitForNodeZeroAgentRestAPIInit(clusterZero, 5)

	// TODO(lranjbar)[AGENT-172]: Add wait for cluster validation
	// Wait for cluster validations to succeed
	// WaitForClusterZeroManifestsToValidate(clusterZero, 5)

	// Wait for cluster Kube API to come up and kubeconfig to be created
	WaitForClusterZeroKubeConfigToExist(clusterZero, 5)
	WaitForClusterZeroKubeAPILive(clusterZero, 5)

	// Wait for bootstrap node to reboot
	WaitForNodeZeroReboot(clusterZero, 30)

	// DEV_NOTES(lranjbar): Maybe this is WaitFor_BootstrapConfigMap_OR_NodeZeroReboot()
	// Research notes: In installer main package create.go:
	// waitForBootstrapComplete(), waitForBootstrapConfigMap()

	// Wait for bootstrap configmap
	// WaitForBootstrapConfigMap(clusterZero, 30)

	return nil
}

// TODO(lranjbar)[AGENT-173]: Add wait for install complete in AGENT-173
// Wait for the installation complete triggered by the agent installer.
func WaitForInstallComplete(assetDir string) error {
	logrus.Info("WaitForInstallComplete")

	// Research notes: In installer main package create.go:
	// waitForInitializedCluster()

	// WaitForBootstrapComplete()

	return nil
}

func WaitForNodeZeroAgentRestAPIInit(clusterZero *zero.ClusterZero, timeoutMins int) error {
	logrus.Info("WaitForNodeZeroAgentRestAPIInit")

	timeout := time.Duration(timeoutMins) * time.Minute
	waitContext, cancel := context.WithTimeout(clusterZero.Ctx, timeout)
	defer cancel()

	wait.Until(func() {
		live, err := clusterZero.Api.Rest.IsAgentAPILive()
		if live && err == nil {
			logrus.Info("Node Zero Agent API Initialized")
			cancel()
		}
	}, 5*time.Second, waitContext.Done())

	return nil
}

func WaitForClusterZeroManifestsToValidate(clusterZero *zero.ClusterZero, timeoutMins int) error {
	logrus.Info("WaitForClusterZeroManifestsToValidate")

	// timeout := time.Duration(timeoutMins) * time.Minute
	// waitContext, cancel := context.WithTimeout(clusterZero.Ctx, timeout)
	// defer cancel()

	// wait.Until(func() {
	// 	clusterState, _ := clusterZero.Get()
	// 	validate, err := clusterZero.Api.Kube.ParseValidationInfo(clusterState)
	// 	if validate && err == nil {
	// 		cancel()
	// 	}
	// }, 5*time.Second, waitContext.Done())

	return nil
}

// DEV_NOTES(lranjbar): Potentially redundant? We will fail when making the client if kubeconfig is not around
func WaitForClusterZeroKubeConfigToExist(clusterZero *zero.ClusterZero, timeoutMins int) error {

	timeout := time.Duration(timeoutMins) * time.Minute
	waitContext, cancel := context.WithTimeout(clusterZero.Ctx, timeout)
	defer cancel()

	wait.Until(func() {
		exist, err := clusterZero.Api.Kube.DoesKubeConfigExist()
		if exist && err == nil {
			logrus.Info("Found kubeconfig")
			cancel()
		}
	}, 5*time.Second, waitContext.Done())

	return nil
}

func WaitForClusterZeroKubeAPILive(clusterZero *zero.ClusterZero, timeoutMins int) error {

	timeout := time.Duration(timeoutMins) * time.Minute
	waitContext, cancel := context.WithTimeout(clusterZero.Ctx, timeout)
	defer cancel()

	wait.Until(func() {
		live, err := clusterZero.Api.Kube.IsKubeAPILive()
		if live && err == nil {
			logrus.Info("Cluster API Initialized")
			cancel()
		}
	}, 5*time.Second, waitContext.Done())

	return nil
}

// DEV_NOTES(lranjbar): Probably rename this function... I think it is safe to assume this logic.
// If the Assisted REST API on node zero is gone then it is likely that node zero has rebooted
func WaitForNodeZeroReboot(clusterZero *zero.ClusterZero, timeoutMins int) error {

	timeout := time.Duration(timeoutMins) * time.Minute
	waitContext, cancel := context.WithTimeout(clusterZero.Ctx, timeout)
	defer cancel()

	wait.Until(func() {
		clusterState, err := clusterZero.Get()
		if err != nil {
			logrus.Info("Node Zero Agent API Gone")
			cancel()
		} else {
			installing, _ := clusterZero.IsInstalling(clusterState)
			if !installing {
				logrus.Fatal("Cluster has stopped installing")
				cancel()
			}
		}
	}, 5*time.Second, waitContext.Done())

	return nil
}

// TODO(lranjbar): Look at waitForBootStrapConfigMap in the main function
func WaitForBootstrapConfigMap(clusterZero *zero.ClusterZero, timeoutMins int) error {

	return nil
}
