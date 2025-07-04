package agent

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/assisted-service/client/installer"
	"github.com/openshift/assisted-service/models"
	"github.com/openshift/installer/pkg/asset/agent/gencrypto"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
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
	workflow               workflow.AgentWorkflowType
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
	NotReadyTime                                        time.Time
	ValidationResults                                   *validationResults
	ClusterInitTime                                     time.Time
}

// NewCluster initializes a Cluster object
func NewCluster(ctx context.Context, assetDir, rendezvousIP, kubeconfigPath, sshKey string, workflowType workflow.AgentWorkflowType) (*Cluster, error) {
	czero := &Cluster{}
	capi := &clientSet{}

	var watcherAuthToken string
	var err error

	switch workflowType {
	case workflow.AgentWorkflowTypeInstall:
		watcherAuthToken, err = FindAuthTokenFromAssetStore(assetDir)
		if err != nil {
			return nil, err
		}
	case workflow.AgentWorkflowTypeAddNodes:
		watcherAuthToken, err = gencrypto.GetWatcherAuthTokenFromCluster(ctx, kubeconfigPath)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("AgentWorkflowType value not supported: %s", workflowType)
	}

	restclient := NewNodeZeroRestClient(ctx, rendezvousIP, sshKey, watcherAuthToken)

	kubeclient, err := NewClusterKubeAPIClient(ctx, kubeconfigPath)
	if err != nil {
		logrus.Fatal(err)
	}

	ocpclient, err := NewClusterOpenShiftAPIClient(ctx, kubeconfigPath)
	if err != nil {
		logrus.Fatal(err)
	}

	capi.Rest = restclient
	capi.Kube = kubeclient
	capi.OpenShift = ocpclient

	cinstallstatushistory := &clusterInstallStatusHistory{
		RestAPISeen:                  false,
		RestAPIInfraEnvEventList:     nil,
		RestAPIPreviousClusterStatus: "",
		RestAPIPreviousEventMessage:  "",
		RestAPIHostValidationsPassed: false,
		ClusterKubeAPISeen:           false,
		ClusterBootstrapComplete:     false,
		ClusterInitTime:              time.Now(),
	}

	cvalidationresults := &validationResults{
		ClusterValidationHistory: make(map[string]*validationResultHistory),
		HostValidationHistory:    make(map[string]map[string]*validationResultHistory),
	}

	czero.Ctx = ctx
	czero.API = capi
	czero.workflow = workflowType
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

	clusterKubeAPILive := czero.API.Kube.IsKubeAPILive()

	agentRestAPILive := czero.API.Rest.IsRestAPILive()

	// Both API's are not available
	if !agentRestAPILive && !clusterKubeAPILive {
		// Current API Status: Agent Rest API: down, Bootstrap Kube API: down
		if !czero.installHistory.RestAPISeen && !czero.installHistory.ClusterKubeAPISeen {
			logrus.Debug("Agent Rest API never initialized. Bootstrap Kube API never initialized")
			elapsedSinceInit := time.Since(czero.installHistory.ClusterInitTime)
			// After allowing time for the interface to come up, check if Node0 can be accessed
			if elapsedSinceInit > 2*time.Minute && !czero.CanConnectToNodeZero() {
				logrus.Info("Cannot access Rendezvous Host. There may be a network configuration problem, check console for additional info")
			} else {
				logrus.Info("Waiting for cluster install to initialize. Sleeping for 30 seconds")
			}

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
		}
		if err != nil {
			logrus.Debug(err)
		}
	}

	// Agent Rest API is available
	if agentRestAPILive {
		exitOnErr, err := czero.MonitorStatusFromAssistedService(nil)
		if err != nil {
			return false, exitOnErr, err
		}
	}

	// cluster bootstrap is not complete
	return false, false, nil
}

// MonitorStatusFromAssistedService (exit-on-error, returned-error)
// checks if the Assisted Service API is up, and both cluster and
// infraenv have been registered.
//
// After those preconditions are met,
// it then reports on the host validation status and overall cluster
// status and updates the cluster's install history.
//
// After cluster or host installation has started, new events from
// the Assisted Service API are also logged and updated to the cluster's
// install history.
func (czero *Cluster) MonitorStatusFromAssistedService(ch chan logEntry) (bool, error) {
	logger := logrus.StandardLogger()
	resource := "cluster"
	if czero.workflow == workflow.AgentWorkflowTypeAddNodes {
		resource = "host"
	}

	// First time we see the agent Rest API
	if !czero.installHistory.RestAPISeen {
		log(Debug, "Agent Rest API Initialized", logger, ch)
		czero.installHistory.RestAPISeen = true
		czero.installHistory.NotReadyTime = time.Now()
	}

	// Lazy loading of the clusterID and clusterInfraEnvID
	if czero.clusterID == nil {
		clusterID, err := czero.API.Rest.getClusterID()
		if err != nil {
			return false, errors.Wrap(err, "Unable to retrieve clusterID from Agent Rest API")
		}
		czero.clusterID = clusterID
	}

	if czero.clusterInfraEnvID == nil {
		clusterInfraEnvID, err := czero.API.Rest.getClusterInfraEnvID()
		if err != nil {
			return false, errors.Wrap(err, "Unable to retrieve clusterInfraEnvID from Agent Rest API")
		}
		czero.clusterInfraEnvID = clusterInfraEnvID
	}

	// Getting cluster metadata from Agent Rest API
	clusterMetadata, err := czero.GetClusterRestAPIMetadata()
	if err != nil {
		return false, errors.Wrap(err, "Unable to retrieve cluster metadata from Agent Rest API")
	}

	if clusterMetadata == nil {
		return false, errors.New("cluster metadata returned nil from Agent Rest API")
	}

	czero.PrintInstallStatus(clusterMetadata, ch)

	// If status indicates pending action, log host info to help pinpoint what is missing
	if (*clusterMetadata.Status != czero.installHistory.RestAPIPreviousClusterStatus) &&
		(*clusterMetadata.Status == models.ClusterStatusInstallingPendingUserAction) {
		for _, host := range clusterMetadata.Hosts {
			if *host.Status == models.ClusterStatusInstallingPendingUserAction {
				log(Warning, fmt.Sprintf("Host %s %s", host.RequestedHostname, *host.StatusInfo), logger, ch)
			}
		}
	}

	if *clusterMetadata.Status == models.ClusterStatusReady {
		stuck, err := czero.IsClusterStuckInReady()
		if err != nil {
			return stuck, err
		}
	} else {
		czero.installHistory.NotReadyTime = time.Now()
	}

	czero.installHistory.RestAPIPreviousClusterStatus = *clusterMetadata.Status

	installing, _ := czero.IsInstalling(*clusterMetadata.Status)
	if !installing {
		errored, _ := czero.HasErrored(*clusterMetadata.Status)
		if errored {
			return false, fmt.Errorf("%s has stopped installing... working to recover installation", resource)
		} else if *clusterMetadata.Status == models.ClusterStatusCancelled {
			return true, fmt.Errorf("%s installation was cancelled", resource)
		}
	}

	validationsErr := checkValidations(clusterMetadata, czero.installHistory.ValidationResults, logger, ch)
	if validationsErr != nil {
		return false, errors.Wrap(validationsErr, "host validations failed")
	}

	// Print most recent event associated with the clusterInfraEnvID
	eventList, err := czero.API.Rest.GetInfraEnvEvents(czero.clusterInfraEnvID)
	if err != nil {
		return false, errors.Wrap(err, fmt.Sprintf("Unable to retrieve events about the %s from the Agent Rest API", resource))
	}
	if len(eventList) == 0 {
		// No cluster events detected from the Agent Rest API
	} else {
		mostRecentEvent := eventList[len(eventList)-1]
		// Don't print the same status message back to back
		if *mostRecentEvent.Message != czero.installHistory.RestAPIPreviousEventMessage {
			if *mostRecentEvent.Severity == models.EventSeverityInfo {
				log(Info, *mostRecentEvent.Message, logger, ch)
			} else {
				log(Warning, *mostRecentEvent.Message, logger, ch)
			}
		}
		czero.installHistory.RestAPIPreviousEventMessage = *mostRecentEvent.Message
		czero.installHistory.RestAPIInfraEnvEventList = eventList
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

// PrintInstallStatus Print a human friendly message using the models from the Agent Rest API.
func (czero *Cluster) PrintInstallStatus(cluster *models.Cluster, ch chan logEntry) {
	friendlyStatus := czero.humanFriendlyClusterInstallStatus(*cluster.Status)
	// Don't print the same status message back to back
	if *cluster.Status != czero.installHistory.RestAPIPreviousClusterStatus {
		log(Info, friendlyStatus, logrus.StandardLogger(), ch)
	}
}

// CanConnectToNodeZero Checks if NodeZero is reachable.
func (czero *Cluster) CanConnectToNodeZero() bool {
	ip := czero.API.Rest.NodeZeroIP
	port := "22"

	conn, err := net.Dial("tcp", net.JoinHostPort(ip, port))
	if err != nil {
		logrus.Debugf("Failed to connect to the Rendezvous Host on port %s: %s", port, err)
		return false
	}
	conn.Close()
	return true
}

// Human friendly install status strings mapped to the Agent Rest API cluster statuses
func (czero *Cluster) humanFriendlyClusterInstallStatus(status string) string {
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
	switch czero.workflow {
	case workflow.AgentWorkflowTypeAddNodes:
		return fmt.Sprintf("Node %s: %s", czero.API.Rest.NodeZeroIP, clusterStoppedInstallingStates[status])
	default:
		return clusterStoppedInstallingStates[status]
	}
}
