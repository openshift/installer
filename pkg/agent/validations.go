package agent

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"

	"github.com/openshift/assisted-service/api/common"
	"github.com/openshift/assisted-service/models"
	"github.com/sirupsen/logrus"
)

const (
	validationFailure string = "failure"
	validationError   string = "error"
)

// Re-using Assisted UI host validation labels (see https://github.com/openshift-assisted/assisted-ui-lib)
// for logging human-friendly messages in case of validation failures
var hostValidationLabels = map[string]string{
	"odf-requirements-satisfied":                      "ODF requirements",
	"disk-encryption-requirements-satisfied":          "Disk encryption requirements",
	"compatible-with-cluster-platform":                "",
	"has-default-route":                               "Default route to host",
	"sufficient-network-latency-requirement-for-role": "Network latency",
	"sufficient-packet-loss-requirement-for-role":     "Packet loss",
	"has-inventory":                                   "Hardware information",
	"has-min-cpu-cores":                               "Minimum CPU cores",
	"has-min-memory":                                  "Minimum Memory",
	"has-min-valid-disks":                             "Minimum disks of required size",
	"has-cpu-cores-for-role":                          "Minimum CPU cores for selected role",
	"has-memory-for-role":                             "Minimum memory for selected role",
	"hostname-unique":                                 "Unique hostname",
	"hostname-valid":                                  "Valid hostname",
	"connected":                                       "Connected",
	"media-connected":                                 "Media Connected",
	"machine-cidr-defined":                            "Machine CIDR",
	"belongs-to-machine-cidr":                         "Belongs to machine CIDR",
	"ignition-downloadable":                           "Ignition file downloadable",
	"belongs-to-majority-group":                       "Belongs to majority connected group",
	"valid-platform-network-settings":                 "Platform network settings",
	"ntp-synced":                                      "NTP synchronization",
	"container-images-available":                      "Container images availability",
	"lso-requirements-satisfied":                      "LSO requirements",
	"ocs-requirements-satisfied":                      "OCS requirements",
	"sufficient-installation-disk-speed":              "Installation disk speed",
	"cnv-requirements-satisfied":                      "CNV requirements",
	"api-domain-name-resolved-correctly":              "API domain name resolution",
	"api-int-domain-name-resolved-correctly":          "API internal domain name resolution",
	"apps-domain-name-resolved-correctly":             "Application ingress domain name resolution",
	"dns-wildcard-not-configured":                     "DNS wildcard not configured",
	"non-overlapping-subnets":                         "Non overlapping subnets",
	"vsphere-disk-uuid-enabled":                       "Vsphere disk uuidenabled",
}

type validationTrace struct {
	header   string
	category string
	label    string
	message  string
}

var previousValidations []validationTrace

func logValidationsStatus(errorMsg string, validations string, log *logrus.Logger) []validationTrace {

	traces := []validationTrace{}
	if validations == "" {
		return traces
	}

	validationsInfo := common.ValidationsStatus{}
	err := json.Unmarshal([]byte(validations), &validationsInfo)
	if err != nil {
		return []validationTrace{{header: errorMsg, message: "unable to verify validations"}}
	}

	for category, validationResults := range validationsInfo {
		for _, r := range validationResults {
			switch r.Status {
			case validationFailure, validationError:
				label := r.ID
				if v, ok := hostValidationLabels[r.ID]; ok {
					label = v
				}

				traces = append(traces, validationTrace{
					header:   errorMsg,
					category: category,
					label:    label,
					message:  r.Message,
				})
			}
		}
	}

	return traces
}

func checkHostsValidations(cluster *models.Cluster, log *logrus.Logger) bool {

	var currentValidations []validationTrace

	currentValidations = append(currentValidations, logValidationsStatus("Validation failure found for cluster", cluster.ValidationsInfo, log)...)
	for _, h := range cluster.Hosts {
		currentValidations = append(currentValidations, logValidationsStatus(fmt.Sprintf("Validation failure found for %s", h.RequestedHostname), h.ValidationsInfo, log)...)
	}

	sort.Slice(currentValidations, func(i, j int) bool {
		if currentValidations[i].header != currentValidations[j].header {
			return currentValidations[i].header < currentValidations[j].header
		}
		if currentValidations[i].category != currentValidations[j].category {
			return currentValidations[i].category < currentValidations[j].category
		}
		return currentValidations[i].label < currentValidations[j].label
	})

	if !reflect.DeepEqual(currentValidations, previousValidations) {
		previousValidations = currentValidations

		if len(previousValidations) == 0 {
			log.Info("Pre-installation validations are OK")
			return true
		}

		log.Info("Checking for validation failures ----------------------------------------------")
		for _, v := range previousValidations {
			log.WithFields(logrus.Fields{
				"category": v.category,
				"label":    v.label,
				"message":  v.message,
			}).Error(v.header)
		}
	}

	return len(previousValidations) == 0
}
