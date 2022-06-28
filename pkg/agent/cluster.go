package agent

import (
	"context"
	"time"

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
	czero.clusterID = nil
	czero.clusterInfraEnvID = nil
	czero.installHistory = cinstallstatushistory
	return czero, nil
}

// GetClusterRestAPIMetadata Retrieve the current cluster metadata from the Agent Rest API
func (czero *Cluster) GetClusterRestAPIMetadata() (*models.Cluster, error) {
	// GET /v2/clusters/{cluster_zero_id}
	getClusterParams := &installer.V2GetClusterParams{ClusterID: *czero.clusterID}
	result, err := czero.API.Rest.Client.Installer.V2GetCluster(czero.Ctx, getClusterParams)
	if err != nil {
		return nil, err
	}
	cluster := result.Payload
	return cluster, nil
}

// IsBootstrapComplete Determine if the cluster that agent installer
// is installing has completed the bootstrap process.
func (czero *Cluster) IsBootstrapComplete() (bool, error) {

	clusterKubeAPILive, clusterKubeAPIErr := czero.API.Kube.IsKubeAPILive()
	if clusterKubeAPIErr != nil {
		logrus.Debug(errors.Wrap(clusterKubeAPIErr, "cluster Kube API is not available"))
	}

	if clusterKubeAPILive {

		// First time we see the cluster Kube API
		if !czero.installHistory.ClusterKubeAPISeen {
			logrus.Info("cluster Kube API Initialized")
			czero.installHistory.ClusterKubeAPISeen = true
		}

		configmap, _ := czero.API.Kube.IsBootstrapConfigMapComplete()
		if configmap {
			logrus.Info("bootstrap configMap status is complete")
			return true, nil
		}
	}

	agentRestAPILive, agentRestAPIErr := czero.API.Rest.IsRestAPILive()
	if agentRestAPIErr != nil {
		logrus.Debug(errors.Wrap(agentRestAPIErr, "node zero Agent API is not available"))
	}

	if agentRestAPILive {

		// First time we see the agent Rest API
		if !czero.installHistory.RestAPISeen {
			logrus.Debug("node zero Agent API Initialized")
			czero.installHistory.RestAPISeen = true
		}

		// Lazy loading of the clusterID and clusterInfraEnvID
		if czero.clusterID == nil {
			clusterID, err := czero.API.Rest.getClusterID()
			if err != nil {
				return false, err
			}
			czero.clusterID = clusterID

		}

		if czero.clusterInfraEnvID == nil {
			clusterInfraEnvID, err := czero.API.Rest.getClusterInfraEnvID()
			if err != nil {
				return false, err
			}
			czero.clusterInfraEnvID = clusterInfraEnvID
		}

		logrus.Trace("getting cluster metadata from Node Zero Agent API")
		clusterMetadata, err := czero.GetClusterRestAPIMetadata()
		if err != nil {
			return false, err
		}

		// TODO[AGENT-172]: Add CheckHostValidations

		czero.PrintInstallStatus(clusterMetadata)
		czero.installHistory.RestAPICurrentClusterStatus = *clusterMetadata.Status

		// Update Install History object when we see these states
		czero.updateInstallHistoryClusterStatus(clusterMetadata)

		installing, _ := czero.IsInstalling(*clusterMetadata.Status)
		if !installing {
			logrus.Warn("Cluster has stopped installing... working to recover installation")
			errored, _ := czero.HasErrored(*clusterMetadata.Status)
			if errored {
				return false, errors.New("cluster installation has stopped due to errors")
			} else if *clusterMetadata.Status == models.ClusterStatusCancelled {
				return false, errors.New("cluster insallation was cancelled")
			}
		}

	}

	// both API's are not available
	if !agentRestAPILive && !clusterKubeAPILive {
		logrus.Debug("current API Status: Node Zero Agent API: down, Cluster Kube API: down")
		if !czero.installHistory.RestAPISeen && !czero.installHistory.ClusterKubeAPISeen {
			logrus.Debug("node zero Agent API never initialized. Cluster API never initialized")
			logrus.Info("Waiting for cluster install to intialize. Sleeping for 30 seconds")
			time.Sleep(30 * time.Second)
			return false, nil
		}

		if czero.installHistory.RestAPISeen && !czero.installHistory.ClusterKubeAPISeen {
			logrus.Debug("cluster API never initialized")
			logrus.Debugf("cluster install status last seen was: %s", czero.installHistory.RestAPICurrentClusterStatus)
			return false, errors.New("cluster installation did not complete")
		}
	}

	logrus.Debug("bootstrap is not complete. sleeping... ")
	return false, nil
}

// IsInstalling Determine if the cluster is still installing using the models from the Agent Rest API.
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

// HasErrored Determine if the cluster installation has errored using the models from the Agent Rest API.
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

// PrintInstallStatus Print a human friendly message using the models from the Agent Rest API.
func (czero *Cluster) PrintInstallStatus(cluster *models.Cluster) error {

	friendlyStatus := humanFriendlyClusterInstallStatus(*cluster.Status)
	logrus.Debug(friendlyStatus)
	// Don't print the same status message back to back
	if *cluster.Status != czero.installHistory.RestAPICurrentClusterStatus {
		logrus.Info(friendlyStatus)
	}

	return nil
}

// Human friendly install status strings mapped to the Agent Rest API cluster statuses
func humanFriendlyClusterInstallStatus(status string) string {
	clusterStoppedInstallingStates := map[string]string{
		models.ClusterStatusAddingHosts:                 "Cluster is adding hosts",
		models.ClusterStatusCancelled:                   "Cluster installation cancelled",
		models.ClusterStatusError:                       "Cluster has hosts in error",
		models.ClusterStatusFinalizing:                  "Finalizing cluster installation",
		models.ClusterStatusInstalling:                  "Cluster installation in progress",
		models.ClusterStatusInstallingPendingUserAction: "Cluster has hosts requiring user input",
		models.ClusterStatusInsufficient:                "Cluster is not ready for install. Check host validations",
		models.ClusterStatusPendingForInput:             "User input is required to continue cluster installation",
		models.ClusterStatusPreparingForInstallation:    "Preparing cluster for installation",
		models.ClusterStatusReady:                       "Cluster is ready for install",
	}
	return clusterStoppedInstallingStates[status]

}

// Update the install history struct when we see the status from the Agent Rest API
func (czero *Cluster) updateInstallHistoryClusterStatus(cluster *models.Cluster) {
	switch *cluster.Status {
	case models.ClusterStatusAddingHosts:
		czero.installHistory.RestAPIClusterStatusAddingHostsSeen = true
	case models.ClusterStatusCancelled:
		czero.installHistory.RestAPIClusterStatusCancelledSeen = true
	case models.ClusterStatusError:
		czero.installHistory.RestAPIClusterStatusErrorSeen = true
	case models.ClusterStatusFinalizing:
		czero.installHistory.RestAPIClusterStatusFinalizingSeen = true
	case models.ClusterStatusInsufficient:
		czero.installHistory.RestAPIClusterStatusInsufficientSeen = true
	case models.ClusterStatusInstalling:
		czero.installHistory.RestAPIClusterStatusInstallingSeen = true
	case models.ClusterStatusInstallingPendingUserAction:
		czero.installHistory.RestAPIClusterStatusInstallingPendingUserActionSeen = true
	case models.ClusterStatusPendingForInput:
		czero.installHistory.RestAPIClusterStatusPendingForInputSeen = true
	case models.ClusterStatusPreparingForInstallation:
		czero.installHistory.RestAPIClusterStatusPreparingForInstallationSeen = true
	case models.ClusterStatusReady:
		czero.installHistory.RestAPIClusterStatusReadySeen = true
	}
}
