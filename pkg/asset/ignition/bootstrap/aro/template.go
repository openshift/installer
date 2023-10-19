package aro

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/types/azure"
)

// TemplateData holds data specific to templates used for the ARO platform.
type TemplateData struct {
	ClusterDomain            string   `json:"cluster_domain"`
	APIIntIP                 string   `json:"api_int_ip"`
	IngressIP                string   `json:"ingress_ip"`
	GatewayDomains           []string `json:"gateway_domains"`
	GatewayPrivateEndpointIP string   `json:"gateway_private_endpoint_ip"`
}

// GetTemplateData returns platform-specific data for bootstrap templates.
func GetTemplateData(config *azure.Platform) *TemplateData {
	// XXX Hack City - *for now*, look at better way later.
	aroInfraFile := os.Getenv("ARO_INFRASTRUCTURE_FILE")
	if aroInfraFile == "" {
		logrus.Warnf("ARO_INFRASTRUCTURE_FILE must be set for ARO")
		return nil
	}

	file, err := ioutil.ReadFile(aroInfraFile)
	if err != nil {
		logrus.Warnf("error reading %s: %v", aroInfraFile, err)
		return nil
	}

	var templateData TemplateData
	err = json.Unmarshal([]byte(file), &templateData)
	if err != nil {
		logrus.Warnf("error parsing ARO JSONs: %v", err)
		return nil
	}

	return &templateData
}
