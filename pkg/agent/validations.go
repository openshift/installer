package agent

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/assisted-service/api/common"
	"github.com/openshift/assisted-service/models"
)

// These validation status strings are defined to match the ones from assisted-service/internal.
// These statuses are the statuses that mattered to us for this implementation.
// Cluster: https://github.com/openshift/assisted-service/blob/master/internal/cluster/validator.go
// Host: https://github.com/openshift/assisted-service/blob/master/internal/host/validator.go
const (
	validationFailure string = "failure"
	validationError   string = "error"
	validationSuccess string = "success"
)

type validationResults struct {
	ClusterValidationHistory map[string]*validationResultHistory
	HostValidationHistory    map[string]map[string]*validationResultHistory
}

type validationResultHistory struct {
	numFailures     int
	seen            bool
	currentStatus   string
	currentMessage  string
	previousStatus  string
	previousMessage string
}

func checkValidations(cluster *models.Cluster, validationResults *validationResults, log *logrus.Logger, ch chan logEntry) error {
	clusterLogPrefix := "Cluster validation: "
	updatedClusterValidationHistory, err := updateValidationResultHistory(clusterLogPrefix, cluster.ValidationsInfo, validationResults.ClusterValidationHistory, log, ch)
	if err != nil {
		return err
	}
	validationResults.ClusterValidationHistory = updatedClusterValidationHistory

	for _, h := range cluster.Hosts {
		hostLogPrefix := "Host " + h.RequestedHostname + " validation: "
		if _, ok := validationResults.HostValidationHistory[h.RequestedHostname]; !ok {
			validationResults.HostValidationHistory[h.RequestedHostname] = make(map[string]*validationResultHistory)
		}
		updatedHostValidationHistory, err := updateValidationResultHistory(hostLogPrefix, h.ValidationsInfo, validationResults.HostValidationHistory[h.RequestedHostname], log, ch)
		if err != nil {
			return err
		}
		validationResults.HostValidationHistory[h.RequestedHostname] = updatedHostValidationHistory
	}
	return nil
}

func updateValidationResultHistory(logPrefix string, validationsInfoString string, validationHistory map[string]*validationResultHistory, log *logrus.Logger, ch chan logEntry) (map[string]*validationResultHistory, error) {
	if validationsInfoString == "" {
		return validationHistory, nil
	}

	validationsInfo := common.ValidationsStatus{}
	err := json.Unmarshal([]byte(validationsInfoString), &validationsInfo)
	if err != nil {
		return nil, errors.Wrap(err, "unable to verify validations")
	}

	for _, validationResults := range validationsInfo {
		for _, r := range validationResults {
			// If validation ID does not exist create it
			if _, ok := validationHistory[r.ID]; !ok {
				validationHistory[r.ID] = &validationResultHistory{}
			}
			validationHistory[r.ID].previousMessage = validationHistory[r.ID].currentMessage
			validationHistory[r.ID].previousStatus = validationHistory[r.ID].currentStatus
			validationHistory[r.ID].currentMessage = r.Message
			validationHistory[r.ID].currentStatus = r.Status
			switch r.Status {
			case validationFailure, validationError:
				validationHistory[r.ID].numFailures++
			}
			logValidationHistory(logPrefix, validationHistory[r.ID], log, ch)
		}
	}
	return validationHistory, nil
}

func logValidationHistory(logPrefix string, history *validationResultHistory, logger *logrus.Logger, ch chan logEntry) {
	// First time we print something
	if !history.seen {
		history.seen = true
		switch history.currentStatus {
		case validationSuccess:
			log(Debug, logPrefix+history.currentMessage, logger, ch)
		case validationFailure, validationError:
			log(Warning, logPrefix+history.currentMessage, logger, ch)
		default:
			log(Trace, logPrefix+history.currentMessage, logger, ch)
		}
		return
	}
	// We have already printed something
	if history.currentMessage != history.previousMessage {
		switch history.currentStatus {
		case validationSuccess:
			if history.previousStatus == validationError || history.previousStatus == validationFailure {
				log(Info, logPrefix+history.currentMessage, logger, ch)
			}
		case validationFailure, validationError:
			log(Warning, logPrefix+history.currentMessage, logger, ch)
		default:
			log(Trace, logPrefix+history.currentMessage, logger, ch)
		}
	}
}
