package zero

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/assisted-service/client/installer"
	"github.com/openshift/assisted-service/models"
)

type ClusterZero struct {
	Ctx                   context.Context
	Api                   *zeroClient
	clusterZeroID         *strfmt.UUID
	clusterZeroInfraEnvID *strfmt.UUID
	installHistory        *clusterInstallStatusHistory
}

type zeroClient struct {
	Kube *clusterZeroKubeAPIClient
	Rest *nodeZeroRestClient
}

type clusterInstallStatusHistory struct {
	AgentRestApiSeen                                  bool
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
	ClusterKubeApiSeen                                bool
}

func NewClusterZero(ctx context.Context, assetDir string) (*ClusterZero, error) {

	czero := &ClusterZero{}
	capi := &zeroClient{}

	restclient, err := NewNodeZeroRestClient(ctx, assetDir)
	if err != nil {
		logrus.Fatal(err)
	}
	kubeclient, err := NewClusterZeroKubeAPIClient(ctx, assetDir)
	if err != nil {
		logrus.Fatal(err)
	}

	capi.Rest = restclient
	capi.Kube = kubeclient

	clusterZeroID, err := restclient.getClusterZeroClusterID()
	if err != nil {
		return nil, err
	}
	clusterZeroInfraEnvID, err := restclient.getClusterZeroInfraEnvID()
	if err != nil {
		return nil, err
	}

	cinstallstatushistory := &clusterInstallStatusHistory{
		AgentRestApiSeen:                                  false,
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
		ClusterKubeApiSeen:                                false,
	}

	czero.Ctx = ctx
	czero.Api = capi
	czero.clusterZeroID = clusterZeroID
	czero.clusterZeroInfraEnvID = clusterZeroInfraEnvID
	czero.installHistory = cinstallstatushistory
	return czero, nil
}

func (czero *ClusterZero) Get() (*models.Cluster, error) {
	// GET /v2/clusters/{cluster_zero_id}
	getClusterParams := &installer.V2GetClusterParams{ClusterID: *czero.clusterZeroID}
	result, err := czero.Api.Rest.Client.Installer.V2GetCluster(czero.Ctx, getClusterParams)
	if err != nil {
		return nil, err
	}
	clusterZero := result.Payload
	return clusterZero, nil
}

func (czero *ClusterZero) IsBootstrapComplete() (bool, error) {

	agentRestApiLive, agentRestApiErr := czero.Api.Rest.IsAgentAPILive()
	if agentRestApiErr != nil {
		logrus.Debug("Node Zero Agent API is not available.")
		logrus.Debug(agentRestApiErr)
	}

	clusterKubeApiLive, clusterKubeApiErr := czero.Api.Kube.IsKubeAPILive()
	if clusterKubeApiErr != nil {
		logrus.Debug("Cluster Kube API is not available.")
		logrus.Debug(clusterKubeApiErr)
	}

	if clusterKubeApiLive {
		// First time we see the cluster Kube Api
		if !czero.installHistory.ClusterKubeApiSeen {
			logrus.Info("Cluster Kube API Initialized")
			czero.installHistory.ClusterKubeApiSeen = true
		}

		configmap, _ := czero.Api.Kube.IsBootstrapConfigMapComplete()
		if configmap {
			logrus.Info("Bootstrap configMap status is complete.")
			return true, nil
		}
	}

	if agentRestApiLive {
		// First time we see the agent Rest Api
		if !czero.installHistory.AgentRestApiSeen {
			logrus.Info("Node Zero Agent API Initialized")
			czero.installHistory.AgentRestApiSeen = true
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

	// both Api's are not available
	if !agentRestApiLive && !clusterKubeApiLive {
		logrus.Debug("Current API Status: Node Zero Agent API: down, Cluster Kube API: down")
		if !czero.installHistory.AgentRestApiSeen && !czero.installHistory.ClusterKubeApiSeen {
			logrus.Debug("Nero Zero Agent API never initialized. Cluster API never initialized.")
			logrus.Warn("Unable to detect installation. Cluster install has either not initalized or was not started.")
			return false, nil
		}

		if czero.installHistory.AgentRestApiSeen && !czero.installHistory.ClusterKubeApiSeen {
			logrus.Debug("Cluster API never initialized.")
			logrus.Debug("Cluster install status last seen was: %s", czero.installHistory.AgentCurrentClusterStatus)
			return false, errors.New("Cluster installation did not complete.")
		}
	}

	logrus.Debug("Bootstrap is not complete. Sleeping...")
	return false, nil
}

func (czero *ClusterZero) IsInstallComplete() (bool, error) {
	return true, nil
}

func (czero *ClusterZero) IsInstalling(status string) (bool, string) {
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

func (czero *ClusterZero) HasErrored(status string) (bool, string) {
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

func (czero *ClusterZero) HasStoppedInstalling(status string) (bool, string) {
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
func (czero *ClusterZero) PrintInstallStatus(cluster *models.Cluster) error {

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

func (czero *ClusterZero) updateInstallHistoryClusterStatus(cluster *models.Cluster) error {

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
