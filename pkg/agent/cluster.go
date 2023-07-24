package agent

import (
	"context"
	"os"
	"path/filepath"
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
	Ctx                    context.Context
	API                    *clientSet
	assetDir               string
	clusterConsoleRouteURL string
	clusterID              *strfmt.UUID
	clusterInfraEnvID      *strfmt.UUID
	installHistory         *clusterInstallStatusHistory
}

type clientSet struct {
	Kube      *ClusterKubeAPIClient
	OpenShift *ClusterOpenShiftAPIClient
	Rest      *NodeZeroRestClient
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
	RestAPIInfraEnvEventList                            models.EventList
	RestAPIPreviousClusterStatus                        string
	RestAPIPreviousEventMessage                         string
	RestAPIHostValidationsPassed                        bool
	ClusterKubeAPISeen                                  bool
	ClusterBootstrapComplete                            bool
	ClusterOperatorsInitialized                         bool
	ClusterConsoleRouteCreated                          bool
	ClusterConsoleRouteURLCreated                       bool
	ClusterInstallComplete                              bool
	NotReadyTime                                        time.Time
	ValidationResults                                   *validationResults
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

	ocpclient, err := NewClusterOpenShiftAPIClient(ctx, assetDir)
	if err != nil {
		logrus.Fatal(err)
	}

	capi.Rest = restclient
	capi.Kube = kubeclient
	capi.OpenShift = ocpclient

	cinstallstatushistory := &clusterInstallStatusHistory{
		RestAPISeen:                   false,
		RestAPIInfraEnvEventList:      nil,
		RestAPIPreviousClusterStatus:  "",
		RestAPIPreviousEventMessage:   "",
		RestAPIHostValidationsPassed:  false,
		ClusterKubeAPISeen:            false,
		ClusterBootstrapComplete:      false,
		ClusterOperatorsInitialized:   false,
		ClusterConsoleRouteCreated:    false,
		ClusterConsoleRouteURLCreated: false,
		ClusterInstallComplete:        false,
	}

	cvalidationresults := &validationResults{
		ClusterValidationHistory: make(map[string]*validationResultHistory),
		HostValidationHistory:    make(map[string]map[string]*validationResultHistory),
	}

	czero.Ctx = ctx
	czero.API = capi
	czero.clusterID = nil
	czero.clusterInfraEnvID = nil
	czero.assetDir = assetDir
	czero.clusterConsoleRouteURL = ""
	czero.installHistory = cinstallstatushistory
	czero.installHistory.ValidationResults = cvalidationresults
	return czero, nil
}

// IsBootstrapComplete (is-bootstrap-complete, exit-on-error, returned-error)
// IsBootstrapComplete Determine if the cluster has completed the bootstrap process.
func (czero *Cluster) IsBootstrapComplete() (bool, bool, error) {

	if czero.installHistory.ClusterBootstrapComplete {
		logrus.Info("Bootstrap is complete")
		return true, false, nil
	}

	clusterKubeAPILive, clusterKubeAPIErr := czero.API.Kube.IsKubeAPILive()
	if clusterKubeAPIErr != nil {
		logrus.Trace(errors.Wrap(clusterKubeAPIErr, "Bootstrap Kube API is not available"))
	}

	agentRestAPILive, agentRestAPIErr := czero.API.Rest.IsRestAPILive()
	if agentRestAPIErr != nil {
		logrus.Trace(errors.Wrap(agentRestAPIErr, "Agent Rest API is not available"))
	}

	// Both API's are not available
	if !agentRestAPILive && !clusterKubeAPILive {
		logrus.Trace("Current API Status: Agent Rest API: down, Bootstrap Kube API: down")
		if !czero.installHistory.RestAPISeen && !czero.installHistory.ClusterKubeAPISeen {
			logrus.Debug("Agent Rest API never initialized. Bootstrap Kube API never initialized")
			logrus.Info("Waiting for cluster install to initialize. Sleeping for 30 seconds")
			time.Sleep(30 * time.Second)
			return false, false, nil
		}

		if czero.installHistory.RestAPISeen && !czero.installHistory.ClusterKubeAPISeen {
			logrus.Debug("Bootstrap Kube API never initialized")
			logrus.Debugf("Cluster install status from Agent Rest API last seen was: %s", czero.installHistory.RestAPIPreviousClusterStatus)
			return false, false, errors.New("cluster bootstrap did not complete")
		}
	}

	// Kube API is available
	if clusterKubeAPILive {

		// First time we see the cluster Kube API
		if !czero.installHistory.ClusterKubeAPISeen {
			logrus.Info("Bootstrap Kube API Initialized")
			czero.installHistory.ClusterKubeAPISeen = true
		}

		configmap, err := czero.API.Kube.IsBootstrapConfigMapComplete()
		if configmap {
			logrus.Info("Bootstrap configMap status is complete")
			czero.installHistory.ClusterBootstrapComplete = true
			return true, false, nil
		}
		if err != nil {
			logrus.Debug(err)
		}
	}

	// Agent Rest API is available
	if agentRestAPILive {

		// First time we see the agent Rest API
		if !czero.installHistory.RestAPISeen {
			logrus.Debug("Agent Rest API Initialized")
			czero.installHistory.RestAPISeen = true
			czero.installHistory.NotReadyTime = time.Now()
		}

		// Lazy loading of the clusterID and clusterInfraEnvID
		if czero.clusterID == nil {
			clusterID, err := czero.API.Rest.getClusterID()
			if err != nil {
				return false, false, errors.Wrap(err, "Unable to retrieve clusterID from Agent Rest API")
			}
			czero.clusterID = clusterID
		}

		if czero.clusterInfraEnvID == nil {
			clusterInfraEnvID, err := czero.API.Rest.getClusterInfraEnvID()
			if err != nil {
				return false, false, errors.Wrap(err, "Unable to retrieve clusterInfraEnvID from Agent Rest API")
			}
			czero.clusterInfraEnvID = clusterInfraEnvID
		}

		logrus.Trace("Getting cluster metadata from Agent Rest API")
		clusterMetadata, err := czero.GetClusterRestAPIMetadata()
		if err != nil {
			return false, false, errors.Wrap(err, "Unable to retrieve cluster metadata from Agent Rest API")
		}

		if clusterMetadata == nil {
			return false, false, errors.New("cluster metadata returned nil from Agent Rest API")
		}

		czero.PrintInstallStatus(clusterMetadata)

		// If status indicates pending action, log host info to help pinpoint what is missing
		if (*clusterMetadata.Status != czero.installHistory.RestAPIPreviousClusterStatus) &&
			(*clusterMetadata.Status == models.ClusterStatusInstallingPendingUserAction) {
			for _, host := range clusterMetadata.Hosts {
				if *host.Status == models.ClusterStatusInstallingPendingUserAction {
					logrus.Warningf("Host %s %s", host.RequestedHostname, *host.StatusInfo)
				}
			}
		}

		if *clusterMetadata.Status == models.ClusterStatusReady {
			stuck, err := czero.IsClusterStuckInReady()
			if err != nil {
				return false, stuck, err
			}
		} else {
			czero.installHistory.NotReadyTime = time.Now()
		}

		czero.installHistory.RestAPIPreviousClusterStatus = *clusterMetadata.Status

		installing, _ := czero.IsInstalling(*clusterMetadata.Status)
		if !installing {
			errored, _ := czero.HasErrored(*clusterMetadata.Status)
			if errored {
				return false, false, errors.New("cluster has stopped installing... working to recover installation")
			} else if *clusterMetadata.Status == models.ClusterStatusCancelled {
				return false, true, errors.New("cluster installation was cancelled")
			}
		}

		validationsErr := checkValidations(clusterMetadata, czero.installHistory.ValidationResults, logrus.StandardLogger())
		if validationsErr != nil {
			return false, false, errors.Wrap(validationsErr, "cluster host validations failed")

		}

		// Print most recent event associated with the clusterInfraEnvID
		eventList, err := czero.API.Rest.GetInfraEnvEvents(czero.clusterInfraEnvID)
		if err != nil {
			return false, false, errors.Wrap(err, "Unable to retrieve events about the cluster from the Agent Rest API")
		}
		if len(eventList) == 0 {
			logrus.Trace("No cluster events detected from the Agent Rest API")
		} else {
			mostRecentEvent := eventList[len(eventList)-1]
			// Don't print the same status message back to back
			if *mostRecentEvent.Message != czero.installHistory.RestAPIPreviousEventMessage {
				if *mostRecentEvent.Severity == models.EventSeverityInfo {
					logrus.Info(*mostRecentEvent.Message)
				} else {
					logrus.Warn(*mostRecentEvent.Message)
				}
			}
			czero.installHistory.RestAPIPreviousEventMessage = *mostRecentEvent.Message
			czero.installHistory.RestAPIInfraEnvEventList = eventList
		}

	}

	logrus.Trace("cluster bootstrap is not complete")
	return false, false, nil
}

// IsInstallComplete Determine if the cluster has completed installation.
func (czero *Cluster) IsInstallComplete() (bool, error) {

	if czero.installHistory.ClusterInstallComplete {
		logrus.Info("Cluster installation is complete")
		return true, nil
	}

	if !czero.installHistory.ClusterOperatorsInitialized {
		initialized, err := czero.API.OpenShift.AreClusterOperatorsInitialized()
		if initialized && err == nil {
			czero.installHistory.ClusterOperatorsInitialized = true
		}
		if err != nil {
			return false, errors.Wrap(err, "Error while initializing cluster operators")
		}

	}

	if !czero.installHistory.ClusterConsoleRouteCreated {
		route, err := czero.API.OpenShift.IsConsoleRouteAvailable()
		if route && err == nil {
			czero.installHistory.ClusterConsoleRouteCreated = true
		}
		if err != nil {
			return false, errors.Wrap(err, "Error while waiting for console route")
		}

	}

	if !czero.installHistory.ClusterConsoleRouteURLCreated {
		available, url, err := czero.API.OpenShift.IsConsoleRouteURLAvailable()
		if available && url != "" && err == nil {
			czero.clusterConsoleRouteURL = url
			czero.installHistory.ClusterConsoleRouteURLCreated = true
		}
		if err != nil {
			return false, errors.Wrap(err, "Error while waiting for console route URL")
		}
	}

	if czero.installHistory.ClusterOperatorsInitialized &&
		czero.installHistory.ClusterConsoleRouteCreated &&
		czero.installHistory.ClusterConsoleRouteURLCreated {
		czero.installHistory.ClusterInstallComplete = true
		return true, nil
	}

	return false, nil
}

// IsClusterStuckInReady Determine if the cluster has stopped transitioning out of the Ready state
func (czero *Cluster) IsClusterStuckInReady() (bool, error) {

	// If the status changes back to Ready from Installing it indicates an error. This condition
	// will be retried
	if czero.installHistory.RestAPIPreviousClusterStatus == models.ClusterStatusPreparingForInstallation {
		return false, errors.New("failed to prepare cluster installation, retrying")
	}

	// Check if stuck in Ready state
	if czero.installHistory.RestAPIPreviousClusterStatus == models.ClusterStatusReady {
		current := time.Now()
		elapsed := current.Sub(czero.installHistory.NotReadyTime)
		if elapsed > 1*time.Minute {
			return true, errors.New("failed to progress after all hosts available")
		}
	}

	return false, nil
}

// GetClusterRestAPIMetadata Retrieve the current cluster metadata from the Agent Rest API
func (czero *Cluster) GetClusterRestAPIMetadata() (*models.Cluster, error) {
	// GET /v2/clusters/{cluster_zero_id}
	if czero.clusterID != nil {
		getClusterParams := &installer.V2GetClusterParams{ClusterID: *czero.clusterID}
		result, err := czero.API.Rest.Client.Installer.V2GetCluster(czero.Ctx, getClusterParams)
		if err != nil {
			return nil, err
		}
		return result.Payload, nil
	}
	return nil, errors.New("no clusterID known for the cluster")
}

// HasErrored Determine if the cluster installation has errored using the models from the Agent Rest API.
func (czero *Cluster) HasErrored(status string) (bool, string) {
	clusterErrorStates := map[string]bool{
		models.ClusterStatusAddingHosts:                 false,
		models.ClusterStatusCancelled:                   false,
		models.ClusterStatusInstalling:                  false,
		models.ClusterStatusInstallingPendingUserAction: true,
		models.ClusterStatusInsufficient:                false,
		models.ClusterStatusError:                       true,
		models.ClusterStatusFinalizing:                  false,
		models.ClusterStatusPendingForInput:             false,
		models.ClusterStatusPreparingForInstallation:    false,
		models.ClusterStatusReady:                       false,
	}
	return clusterErrorStates[status], status
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

// PrintInfraEnvRestAPIEventList Prints the whole event list for debugging
func (czero *Cluster) PrintInfraEnvRestAPIEventList() {
	if czero.installHistory.RestAPIInfraEnvEventList != nil {
		for i := 0; i < len(czero.installHistory.RestAPIInfraEnvEventList); i++ {
			logrus.Debug(*czero.installHistory.RestAPIInfraEnvEventList[i].Message)
		}
	} else {
		logrus.Debug("No events logged from the Agent Rest API")
	}
}

// PrintInstallationComplete Prints the installation complete information
func (czero *Cluster) PrintInstallationComplete() error {
	absDir, err := filepath.Abs(czero.assetDir)
	if err != nil {
		return err
	}
	kubeconfig := filepath.Join(absDir, "auth", "kubeconfig")
	pwFile := filepath.Join(absDir, "auth", "kubeadmin-password")
	pw, err := os.ReadFile(pwFile)
	if err != nil {
		return err
	}
	logrus.Info("Install complete!")
	logrus.Infof("To access the cluster as the system:admin user when using 'oc', run\n    export KUBECONFIG=%s", kubeconfig)
	logrus.Infof("Access the OpenShift web-console here: %s", czero.clusterConsoleRouteURL)
	logrus.Infof("Login to the console with user: %q, and password: %q", "kubeadmin", pw)
	return nil

}

// PrintInstallStatus Print a human friendly message using the models from the Agent Rest API.
func (czero *Cluster) PrintInstallStatus(cluster *models.Cluster) error {

	friendlyStatus := humanFriendlyClusterInstallStatus(*cluster.Status)
	logrus.Trace(friendlyStatus)
	// Don't print the same status message back to back
	if *cluster.Status != czero.installHistory.RestAPIPreviousClusterStatus {
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
		models.ClusterStatusInsufficient:                "Cluster is not ready for install. Check validations",
		models.ClusterStatusPendingForInput:             "User input is required to continue cluster installation",
		models.ClusterStatusPreparingForInstallation:    "Preparing cluster for installation",
		models.ClusterStatusReady:                       "Cluster is ready for install",
	}
	return clusterStoppedInstallingStates[status]

}
