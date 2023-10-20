package aro

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/openshift/installer/pkg/types/azure"
)

// TemplateData holds data specific to templates used for the ARO platform.
type TemplateData struct {

	// DNSmasq
	ClusterDomain            string   `json:"clusterDomain" yaml:"clusterDomain"`
	APIIntIP                 string   `json:"apiIntIP" yaml:"apiIntIP"`
	IngressIP                string   `json:"ingressIP" yaml:"ingressIP"`
	GatewayDomains           []string `json:"gatewayDomains" yaml:"gatewayDomains"`
	GatewayPrivateEndpointIP string   `json:"gatewayPrivateEndpointIP" yaml:"gatewayPrivateEndpointIP"`

	// Logging
	Account           string `json:"account" yaml:"account"`
	Certificate       string `json:"ceriticate" yaml:"certificate"`
	ConfigVersion     string `json:"configVersion" yaml:"configVersion"`
	Environment       string `json:"environment" yaml:"environment"`
	FluentbitImage    string `json:"fluentbitImage" yaml:"fluentbitImage"`
	Key               string `json:"key" yaml:"key"`
	MdsdImage         string `json:"mdsdImage" yaml:"mdsdImage"`
	Namespace         string `json:"namespace" yaml:"namespace"`
	Region            string `json:"region": yaml:"region"`
	ResourceGroupName string `json:"resourceGroupName" yaml:"resourceGroupName"`
	ResourceID        string `json:"resourceID" yaml:"resourceID"`
	ResourceName      string `json:"resourceName" yaml:"resourceName"`
	SubscriptionID    string `json:"subscriptionID" yaml:"subscriptionID"`
}

// GetTemplateData returns platform-specific data for bootstrap templates.
func GetTemplateData(config *azure.Platform) *TemplateData {
	// XXX Hack City - *for now*, look at better way later.
	aroInfraFile := strings.ToLower(os.Getenv("ARO_INFRASTRUCTURE_FILE"))
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
	if strings.HasSuffix(aroInfraFile, "json") {
		err = json.Unmarshal([]byte(file), &templateData)
		if err != nil {
			logrus.Warnf("error parsing ARO JSON config: %v", err)
			return nil
		}

	} else if strings.HasSuffix(aroInfraFile, "yaml") {
		err = yaml.Unmarshal([]byte(file), &templateData)
		if err != nil {
			logrus.Warnf("error parsing ARO YAML config: %v", err)
			return nil
		}
	}

	return &templateData
}
