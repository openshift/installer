package agent

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/assisted-service/client/installer"
	"github.com/openshift/assisted-service/models"
)

// Cluster is a struct designed to help interact with the cluster that is
// currently being installed by agent installer.
type Cluster struct {
	Ctx               context.Context
	API               *clientSet
	clusterID         *strfmt.UUID
	clusterInfraEnvID *strfmt.UUID
	installHistory    *clusterInstallStatusHistory
}

type clientSet struct {
	Kube *ClusterKubeAPIClient
	Rest *NodeZeroRestClient
}

type clusterInstallStatusHistory struct {
	RestAPISeen                                         bool
	RestAPIClusterStatusAddingHostsSeen                 bool
	RestAPIClusterStatusCancelledSeen                   bool
	RestAPIClusterStatusInstallingSeen                  bool
	RestAPIClusterStatusInstallingPendingUserActionSeen bool
	RestAPIClusterStatusInsufficientSeen                bool
	RestAPIClusterStatusFinalizingSeen                  bool
	RestAPIClusterStatusErrorSeen                       bool
	RestAPIClusterStatusPendingForInputSeen             bool
	RestAPIClusterStatusPreparingForInstallationSeen    bool
	RestAPIClusterStatusReadySeen                       bool
	RestAPICurrentClusterStatus                         string
	RestAPIPreviousClusterStatus                        string
	RestAPIHostValidationsPassed                        bool
	ClusterKubeAPISeen                                  bool
}

// NewCluster initializes a Cluster object
func NewCluster(ctx context.Context, assetDir string) (*Cluster, error) {

	czero := &Cluster{}
	capi := &clientSet{}

	restclient, err := NewNodeZeroRestClient(ctx, assetDir)
	if err != nil {
		logrus.Fatal(err)
	}
	kubeclient, err := NewClusterKubeAPIClient(ctx, assetDir)
	if err != nil {
		logrus.Fatal(err)
	}

	capi.Rest = restclient
	capi.Kube = kubeclient

	clusterID, err := restclient.getClusterID()
	if err != nil {
		return nil, err
	}
	clusterInfraEnvID, err := restclient.getClusterInfraEnvID()
	if err != nil {
		return nil, err
	}

	cinstallstatushistory := &clusterInstallStatusHistory{
		RestAPISeen:                                         false,
		RestAPIClusterStatusAddingHostsSeen:                 false,
		RestAPIClusterStatusCancelledSeen:                   false,
		RestAPIClusterStatusInstallingSeen:                  false,
		RestAPIClusterStatusInstallingPendingUserActionSeen: false,
		RestAPIClusterStatusInsufficientSeen:                false,
		RestAPIClusterStatusFinalizingSeen:                  false,
		RestAPIClusterStatusErrorSeen:                       false,
		RestAPIClusterStatusPendingForInputSeen:             false,
		RestAPIClusterStatusPreparingForInstallationSeen:    false,
		RestAPIClusterStatusReadySeen:                       false,
		RestAPICurrentClusterStatus:                         "",
		RestAPIPreviousClusterStatus:                        "",
		RestAPIHostValidationsPassed:                        false,
		ClusterKubeAPISeen:                                  false,
	}

	czero.Ctx = ctx
	czero.API = capi
	czero.clusterID = clusterID
	czero.clusterInfraEnvID = clusterInfraEnvID
	czero.installHistory = cinstallstatushistory
	return czero, nil
}

// Cluster.Get Retrieve the current cluster metadata from the Agent Rest API
func (czero *Cluster) Get() (*models.Cluster, error) {
	// GET /v2/clusters/{cluster_zero_id}
	getClusterParams := &installer.V2GetClusterParams{ClusterID: *czero.clusterID}
	result, err := czero.API.Rest.Client.Installer.V2GetCluster(czero.Ctx, getClusterParams)
	if err != nil {
		return nil, err
	}
	cluster := result.Payload
	return cluster, nil
}

// Cluster.IsBootstrapComplete Determine if the cluster that agent installer
// is installing has completed the bootstrap process.
func (czero *Cluster) IsBootstrapComplete() (bool, error) {

	agentRestAPILive, agentRestAPIErr := czero.API.Rest.IsRestAPILive()
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
		if !czero.installHistory.RestAPISeen {
			logrus.Info("Node Zero Agent API Initialized")
			czero.installHistory.RestAPISeen = true
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
		czero.installHistory.RestAPICurrentClusterStatus = *clusterState.Status

		// Update Install History object when we see these states
		czero.updateInstallHistoryClusterStatus(clusterState)

		stopped, _ := czero.HasStoppedInstalling(*clusterState.Status)
		if stopped {
			logrus.Error("Cluster has stopped installing")
			errored, _ := czero.HasErrored(*clusterState.Status)
			if errored {
				return false, errors.New("Cluster installation has stopped due to errors")
			}
			return false, errors.New("Cluster has stopped installing and/or insallation was cancelled")
		}

	}

	// both API's are not available
	if !agentRestAPILive && !clusterKubeAPILive {
		logrus.Debug("Current API Status: Node Zero Agent API: down, Cluster Kube API: down")
		if !czero.installHistory.RestAPISeen && !czero.installHistory.ClusterKubeAPISeen {
			logrus.Debug("Nero Zero Agent API never initialized. Cluster API never initialized.")
			logrus.Warn("Unable to detect installation. Cluster install has either not initalized or was not started.")
			return false, nil
		}

		if czero.installHistory.RestAPISeen && !czero.installHistory.ClusterKubeAPISeen {
			logrus.Debug("Cluster API never initialized.")
			logrus.Debug("Cluster install status last seen was: %s", czero.installHistory.RestAPICurrentClusterStatus)
			return false, errors.New("Cluster installation did not complete.")
		}
	}

	logrus.Debug("Bootstrap is not complete. Sleeping...")
	return false, nil
}

// Cluster.IsInstallComplete Determine if the cluster has completed installation.
func (czero *Cluster) IsInstallComplete() (bool, error) {
	return true, nil
}

// Cluster.IsInstalling Determine if the cluster is still installing using
// the models from the Agent Rest API.
func (czero *Cluster) IsInstalling(status string) (bool, string) {
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

// Cluster.HasErrored Determine if the cluster installation has errored using
// the models from the Agent Rest API.
func (czero *Cluster) HasErrored(status string) (bool, string) {
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

// Cluster.IsInstalling Determine if the cluster has stopped installing using
// the models from the Agent Rest API.
func (czero *Cluster) HasStoppedInstalling(status string) (bool, string) {
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

// PrintInstallStatus Print a human friendly message using the models from the
// Agent Rest API.
func (czero *Cluster) PrintInstallStatus(cluster *models.Cluster) error {

	// Don't print the same status message back to back
	if *cluster.Status != czero.installHistory.RestAPICurrentClusterStatus {
		friendlyStatus := humanFriendlyClusterInstallStatus(*cluster.Status)
		logrus.Info(friendlyStatus)
	}

	return nil
}

// Human friendly install status strings mapped to the Agent Rest API cluster statuses
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

// Update the install history struct when we see the status from the Agent Rest API
func (czero *Cluster) updateInstallHistoryClusterStatus(cluster *models.Cluster) error {

	switch *cluster.Status {
	case models.ClusterStatusAddingHosts:
		if !czero.installHistory.RestAPIClusterStatusAddingHostsSeen {
			czero.installHistory.RestAPIClusterStatusAddingHostsSeen = true
		}
	case models.ClusterStatusCancelled:
		if !czero.installHistory.RestAPIClusterStatusCancelledSeen {
			czero.installHistory.RestAPIClusterStatusCancelledSeen = true
		}
	case models.ClusterStatusError:
		if !czero.installHistory.RestAPIClusterStatusErrorSeen {
			czero.installHistory.RestAPIClusterStatusErrorSeen = true
		}
	case models.ClusterStatusFinalizing:
		if !czero.installHistory.RestAPIClusterStatusFinalizingSeen {
			czero.installHistory.RestAPIClusterStatusFinalizingSeen = true
		}
	case models.ClusterStatusInsufficient:
		if !czero.installHistory.RestAPIClusterStatusInsufficientSeen {
			czero.installHistory.RestAPIClusterStatusInsufficientSeen = true
		}
	case models.ClusterStatusInstalling:
		if !czero.installHistory.RestAPIClusterStatusInstallingSeen {
			czero.installHistory.RestAPIClusterStatusInstallingSeen = true
		}
	case models.ClusterStatusInstallingPendingUserAction:
		if !czero.installHistory.RestAPIClusterStatusInstallingPendingUserActionSeen {
			czero.installHistory.RestAPIClusterStatusInstallingPendingUserActionSeen = true
		}
	case models.ClusterStatusPendingForInput:
		if !czero.installHistory.RestAPIClusterStatusPendingForInputSeen {
			czero.installHistory.RestAPIClusterStatusPendingForInputSeen = true
		}
	case models.ClusterStatusPreparingForInstallation:
		if !czero.installHistory.RestAPIClusterStatusPreparingForInstallationSeen {
			czero.installHistory.RestAPIClusterStatusPreparingForInstallationSeen = true
		}
	case models.ClusterStatusReady:
		if !czero.installHistory.RestAPIClusterStatusReadySeen {
			czero.installHistory.RestAPIClusterStatusReadySeen = true
		}
	}

	return nil
}
