package agent

import (
	"context"
	"os"

	"github.com/openshift/assisted-service/client/installer"
	"github.com/openshift/assisted-service/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// PrintMostRecentEvent Print most recent event associated with the clusterInfraEnvID
func PrintMostRecentEvent(cluster *Cluster) {
	eventList, err := cluster.API.Rest.GetInfraEnvEvents(cluster.clusterInfraEnvID)
	if err != nil {
		errors.Wrap(err, "Unable to retrieve events about the cluster from the Agent Rest API")
	}
	if len(eventList) == 0 {
		logrus.Trace("No cluster installation events detected from the Agent Rest API.")
	} else {
		mostRecentEvent := eventList[len(eventList)-1]
		// Don't print the same status message back to back
		if *mostRecentEvent.Message != cluster.installHistory.RestAPIPreviousEventMessage {
			if *mostRecentEvent.Severity == models.EventSeverityInfo {
				logrus.Info(*mostRecentEvent.Message)
			} else {
				logrus.Warn(*mostRecentEvent.Message)
			}
		}
		cluster.installHistory.RestAPIPreviousEventMessage = *mostRecentEvent.Message
		cluster.installHistory.RestAPIInfraEnvEventList = eventList
	}
}

// RetrieveEventsLog Retrieve the event log from the Agent API
func RetrieveEventsLog(assetDir string) {
	ctx := context.Background()
	cluster, err := NewCluster(ctx, assetDir)
	if err != nil {
		logrus.Warn("unable to make cluster object to track installation")
	}
	eventList, err := cluster.API.Rest.GetInfraEnvEvents(cluster.clusterInfraEnvID)
	if err != nil {
		errors.Wrap(err, "Agent Rest API is not available. Unable to retrieve installation events.")
	}
	PrintFullEventsLog(eventList)
}

// PrintFullEventsLog Print the entire event log from the Agent API
func PrintFullEventsLog(eventList models.EventList) {
	if len(eventList) == 0 {
		logrus.Trace("No cluster installation events detected from the Agent Rest API.")
	} else {
		for _, event := range eventList {
			logrus.Info(event.Message)
		}
	}
}

// DownloadClusterLogs Download the cluster logs from the Agent API
func DownloadClusterLogs(assetDir string) error {
	ctx := context.Background()
	cluster, err := NewCluster(ctx, assetDir)
	if err != nil {
		logrus.Warn("unable to make cluster object to track installation")
	}

	filename := "logs-" + string(*cluster.clusterID)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer func() {
		file.Close()
	}()

	logrus.Infof("Downloading cluster logs to %s", filename)
	_, err = cluster.API.Rest.Client.Installer.V2DownloadClusterLogs(
		ctx, &installer.V2DownloadClusterLogsParams{ClusterID: *cluster.clusterID}, file)
	if err != nil {
		return err
	}

	return nil
}
