package agent

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/assisted-service/client/installer"
	"github.com/openshift/assisted-service/models"
)

type AgentCluster struct {
	Ctx                    context.Context
	API                    *clientSet
	agentClusterID         *strfmt.UUID
	agentClusterInfraEnvID *strfmt.UUID
	installHistory         *clusterInstallStatusHistory
}

type clientSet struct {
	Kube *agentClusterKubeAPIClient
	Rest *nodeZeroRestClient
}

type clusterInstallStatusHistory struct {
	AgentRestAPISeen                                  bool
	AgentClusterStatusAddingHostsSeen                 bool
	AgentClusterStatusCancelledSeen                   bool
	AgentClusterStatusInstallingSeen                  bool
	AgentClusterStatusInstallingPendingUserActionSeen bool
	AgentClusterStatusInsufficientSeen                bool
	AgentClusterStatusFinalizingSeen                  bool
	AgentClusterStatusErrorSeen                       bool
	AgentClusterStatusPendingForInputSeen             bool
	AgentClusterStatusPreparingForInstallationSeen    bool
	AgentClusterStatusReadySeen                       bool
	AgentCurrentClusterStatus                         string
	AgentPreviousClusterStatus                        string
	AgentHostValidationsPassed                        bool
	ClusterKubeAPISeen                                bool
}

// NewAgentCluster initializes a AgentCluster object
func NewAgentCluster(ctx context.Context, assetDir string) (*AgentCluster, error) {

	czero := &AgentCluster{}
	capi := &clientSet{}

	restclient, err := NewNodeZeroRestClient(ctx, assetDir)
	if err != nil {
		logrus.Fatal(err)
	}
	kubeclient, err := NewAgentClusterKubeAPIClient(ctx, assetDir)
	if err != nil {
		logrus.Fatal(err)
	}

	capi.Rest = restclient
	capi.Kube = kubeclient

	agentClusterID, err := restclient.getAgentClusterClusterID()
	if err != nil {
		return nil, err
	}
	agentClusterInfraEnvID, err := restclient.getAgentClusterInfraEnvID()
	if err != nil {
		return nil, err
	}

	cinstallstatushistory := &clusterInstallStatusHistory{
		AgentRestAPISeen:                                  false,
		AgentClusterStatusAddingHostsSeen:                 false,
		AgentClusterStatusCancelledSeen:                   false,
		AgentClusterStatusInstallingSeen:                  false,
		AgentClusterStatusInstallingPendingUserActionSeen: false,
		AgentClusterStatusInsufficientSeen:                false,
		AgentClusterStatusFinalizingSeen:                  false,
		AgentClusterStatusErrorSeen:                       false,
		AgentClusterStatusPendingForInputSeen:             false,
		AgentClusterStatusPreparingForInstallationSeen:    false,
		AgentClusterStatusReadySeen:                       false,
		AgentCurrentClusterStatus:                         "",
		AgentPreviousClusterStatus:                        "",
		AgentHostValidationsPassed:                        false,
		ClusterKubeAPISeen:                                false,
	}

	czero.Ctx = ctx
	czero.API = capi
	czero.agentClusterID = agentClusterID
	czero.agentClusterInfraEnvID = agentClusterInfraEnvID
	czero.installHistory = cinstallstatushistory
	return czero, nil
}

func (czero *AgentCluster) Get() (*models.Cluster, error) {
	// GET /v2/clusters/{cluster_zero_id}
	getClusterParams := &installer.V2GetClusterParams{ClusterID: *czero.agentClusterID}
	result, err := czero.API.Rest.Client.Installer.V2GetCluster(czero.Ctx, getClusterParams)
	if err != nil {
		return nil, err
	}
	agentCluster := result.Payload
	return agentCluster, nil
}

func (czero *AgentCluster) IsBootstrapComplete() (bool, error) {

	agentRestAPILive, agentRestAPIErr := czero.API.Rest.IsAgentAPILive()
	if agentRestAPIErr != nil {
		logrus.Debug("Node Zero Agent API is not available.")
		logrus.Debug(agentRestAPIErr)
	}

	clusterKubeAPILive, clusterKubeAPIErr := czero.API.Kube.IsKubeAPILive()
	if clusterKubeAPIErr != nil {
		logrus.Debug("Cluster Kube API is not available.")
		logrus.Debug(clusterKubeAPIErr)
	}

	if clusterKubeAPILive {
		// First time we see the cluster Kube API
		if !czero.installHistory.ClusterKubeAPISeen {
			logrus.Info("Cluster Kube API Initialized")
			czero.installHistory.ClusterKubeAPISeen = true
		}

		configmap, _ := czero.API.Kube.IsBootstrapConfigMapComplete()
		if configmap {
			logrus.Info("Bootstrap configMap status is complete.")
			return true, nil
		}
	}

	if agentRestAPILive {
		// First time we see the agent Rest API
		if !czero.installHistory.AgentRestAPISeen {
			logrus.Info("Node Zero Agent API Initialized")
			czero.installHistory.AgentRestAPISeen = true
		}
		logrus.Trace("Getting cluster info from Node Zero Agent API")
		clusterState, _ := czero.Get()

		// TODO(lranjbar)[AGENT-172]: Add CheckHostValidations
		// if !czero.installHistory.AgentHostValidationsPassed {
		// 	validations := CheckHostValidations()
		// 	if validations {
		// 		czero.installHistory.AgentHostValidationsPassed = true
		// 	}
		// }

		czero.PrintInstallStatus(clusterState)
		czero.installHistory.AgentCurrentClusterStatus = *clusterState.Status

		// Update Install History object when we see these states
		czero.updateInstallHistoryClusterStatus(clusterState)

		stopped, _ := czero.HasStoppedInstalling(*clusterState.Status)
		if stopped {
			logrus.Error("Cluster has stopped installing.")
			errored, _ := czero.HasErrored(*clusterState.Status)
			if errored {
				return false, errors.New("Cluster installation has stopped due to errors.")
			}
			return false, errors.New("Cluster has stopped installing and/or insallation was cancelled.")
		}

	}

	// both API's are not available
	if !agentRestAPILive && !clusterKubeAPILive {
		logrus.Debug("Current API Status: Node Zero Agent API: down, Cluster Kube API: down")
		if !czero.installHistory.AgentRestAPISeen && !czero.installHistory.ClusterKubeAPISeen {
			logrus.Debug("Nero Zero Agent API never initialized. Cluster API never initialized.")
			logrus.Warn("Unable to detect installation. Cluster install has either not initalized or was not started.")
			return false, nil
		}

		if czero.installHistory.AgentRestAPISeen && !czero.installHistory.ClusterKubeAPISeen {
			logrus.Debug("Cluster API never initialized.")
			logrus.Debug("Cluster install status last seen was: %s", czero.installHistory.AgentCurrentClusterStatus)
			return false, errors.New("Cluster installation did not complete.")
		}
	}

	logrus.Debug("Bootstrap is not complete. Sleeping...")
	return false, nil
}

func (czero *AgentCluster) IsInstallComplete() (bool, error) {
	return true, nil
}

func (czero *AgentCluster) IsInstalling(status string) (bool, string) {
	clusterInstallingStates := map[string]bool{
		models.ClusterStatusAddingHosts:                 true,
		models.ClusterStatusCancelled:                   false,
		models.ClusterStatusInstalling:                  true,
		models.ClusterStatusInstallingPendingUserAction: false,
		models.ClusterStatusInsufficient:                false,
		models.ClusterStatusError:                       false,
		models.ClusterStatusFinalizing:                  true,
		models.ClusterStatusPendingForInput:             true,
		models.ClusterStatusPreparingForInstallation:    true,
		models.ClusterStatusReady:                       true,
	}
	return clusterInstallingStates[status], status
}

func (czero *AgentCluster) HasErrored(status string) (bool, string) {
	clusterErrorStates := map[string]bool{
		models.ClusterStatusAddingHosts:                 false,
		models.ClusterStatusCancelled:                   false,
		models.ClusterStatusInstalling:                  false,
		models.ClusterStatusInstallingPendingUserAction: true,
		models.ClusterStatusInsufficient:                true,
		models.ClusterStatusError:                       true,
		models.ClusterStatusFinalizing:                  false,
		models.ClusterStatusPendingForInput:             false,
		models.ClusterStatusPreparingForInstallation:    false,
		models.ClusterStatusReady:                       false,
	}
	return clusterErrorStates[status], status
}

func (czero *AgentCluster) HasStoppedInstalling(status string) (bool, string) {
	clusterStoppedInstallingStates := map[string]bool{
		models.ClusterStatusAddingHosts:                 false,
		models.ClusterStatusCancelled:                   true,
		models.ClusterStatusFinalizing:                  false,
		models.ClusterStatusInstalling:                  false,
		models.ClusterStatusInstallingPendingUserAction: true,
		models.ClusterStatusInsufficient:                true,
		models.ClusterStatusError:                       true,
		models.ClusterStatusPendingForInput:             false,
		models.ClusterStatusPreparingForInstallation:    false,
		models.ClusterStatusReady:                       false,
	}
	return clusterStoppedInstallingStates[status], status
}

// TODO(lranjbar): Print install status from the Cluster object
func (czero *AgentCluster) PrintInstallStatus(cluster *models.Cluster) error {

	// Don't print the same status message back to back
	if *cluster.Status != czero.installHistory.AgentCurrentClusterStatus {
		friendlyStatus := humanFriendlyClusterInstallStatus(*cluster.Status)
		logrus.Info(friendlyStatus)
	}

	return nil
}

func humanFriendlyClusterInstallStatus(status string) string {
	clusterStoppedInstallingStates := map[string]string{
		models.ClusterStatusAddingHosts:                 "Cluster is adding hosts.",
		models.ClusterStatusCancelled:                   "Cluster installation cancelled.",
		models.ClusterStatusError:                       "Cluster has hosts in error.",
		models.ClusterStatusFinalizing:                  "Finalizing cluster installation.",
		models.ClusterStatusInstalling:                  "Cluster installation in progress.",
		models.ClusterStatusInstallingPendingUserAction: "Cluster has hosts requiring user input.",
		models.ClusterStatusInsufficient:                "Cluster is not ready for install. Check hardware settings.",
		models.ClusterStatusPendingForInput:             "User input is required to continue cluster installation.",
		models.ClusterStatusPreparingForInstallation:    "Preparing cluster for installation.",
		models.ClusterStatusReady:                       "Cluster is ready for install.",
	}
	return clusterStoppedInstallingStates[status]

}

func (czero *AgentCluster) updateInstallHistoryClusterStatus(cluster *models.Cluster) error {

	switch *cluster.Status {
	case models.ClusterStatusAddingHosts:
		if !czero.installHistory.AgentClusterStatusAddingHostsSeen {
			czero.installHistory.AgentClusterStatusAddingHostsSeen = true
		}
	case models.ClusterStatusCancelled:
		if !czero.installHistory.AgentClusterStatusCancelledSeen {
			czero.installHistory.AgentClusterStatusCancelledSeen = true
		}
	case models.ClusterStatusError:
		if !czero.installHistory.AgentClusterStatusErrorSeen {
			czero.installHistory.AgentClusterStatusErrorSeen = true
		}
	case models.ClusterStatusFinalizing:
		if !czero.installHistory.AgentClusterStatusFinalizingSeen {
			czero.installHistory.AgentClusterStatusFinalizingSeen = true
		}
	case models.ClusterStatusInsufficient:
		if !czero.installHistory.AgentClusterStatusInsufficientSeen {
			czero.installHistory.AgentClusterStatusInsufficientSeen = true
		}
	case models.ClusterStatusInstalling:
		if !czero.installHistory.AgentClusterStatusInstallingSeen {
			czero.installHistory.AgentClusterStatusInstallingSeen = true
		}
	case models.ClusterStatusInstallingPendingUserAction:
		if !czero.installHistory.AgentClusterStatusInstallingPendingUserActionSeen {
			czero.installHistory.AgentClusterStatusInstallingPendingUserActionSeen = true
		}
	case models.ClusterStatusPendingForInput:
		if !czero.installHistory.AgentClusterStatusPendingForInputSeen {
			czero.installHistory.AgentClusterStatusPendingForInputSeen = true
		}
	case models.ClusterStatusPreparingForInstallation:
		if !czero.installHistory.AgentClusterStatusPreparingForInstallationSeen {
			czero.installHistory.AgentClusterStatusPreparingForInstallationSeen = true
		}
	case models.ClusterStatusReady:
		if !czero.installHistory.AgentClusterStatusReadySeen {
			czero.installHistory.AgentClusterStatusReadySeen = true
		}
	}

	return nil
}
