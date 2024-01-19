package ibmcloud

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/sirupsen/logrus"

	configv1 "github.com/openshift/api/config/v1"
)

// https://github.com/kubernetes/kubernetes/blob/368ee4bb8ee7a0c18431cd87ee49f0c890aa53e5/staging/src/k8s.io/legacy-cloud-providers/gce/gce.go#L188
type config struct {
	Global     global     `gcfg:"global"`
	Kubernetes kubernetes `gcfg:"kubernetes"`
	Provider   provider   `gcfg:"provider"`
}

type global struct {
	Version string `gcfg:"version"`
}

type kubernetes struct {
	ConfigFile string `gcfg:"config-file"`
}

type provider struct {
	AccountID                string `gcfg:"accountID"`
	ClusterID                string `gcfg:"clusterID"`
	ClusterDefaultProvider   string `gcfg:"cluster-default-provider"`
	Region                   string `gcfg:"region"`
	G2CredentialsFilePath    string `gcfg:"g2Credentials"`
	G2ResourceGroupName      string `gcfg:"g2ResourceGroupName"`
	G2VPCName                string `gcfg:"g2VpcName"`
	G2WorkerServiceAccountID string `gcfg:"g2workerServiceAccountID"`
	G2VPCSubnetNames         string `gcfg:"g2VpcSubnetNames"`
	IAMEndpointOverride      string `gcfg:"iamEndpointOverride,omitempty"`
	VPCEndpointOverride      string `gcfg:"g2EndpointOverride,omitempty"`
	RMEndpointOverride       string `gcfg:"rmEndpointOverride,omitempty"`
}

// CloudProviderConfig generates the cloud provider config for the IBMCloud platform.
func CloudProviderConfig(infraID string, accountID string, region string, resourceGroupName string, vpcName string, subnets []string, controlPlaneZones []string, computeZones []string, serviceEndpoints []configv1.IBMCloudServiceEndpoint) (string, error) {
	if vpcName == "" {
		vpcName = fmt.Sprintf("%s-vpc", infraID)
	}

	var subnetNames string
	if len(subnets) > 0 {
		subnetNames = strings.Join(subnets, ",")
	} else {
		subnetNames = getVpcSubnetNames(infraID, controlPlaneZones, computeZones)
	}

	config := &config{
		Global: global{
			Version: "1.1.0",
		},
		Kubernetes: kubernetes{
			ConfigFile: "",
		},
		Provider: provider{
			AccountID:                accountID,
			ClusterID:                infraID,
			ClusterDefaultProvider:   "g2",
			Region:                   region,
			G2CredentialsFilePath:    "/etc/vpc/ibmcloud_api_key",
			G2ResourceGroupName:      resourceGroupName,
			G2VPCName:                vpcName,
			G2WorkerServiceAccountID: accountID,
			G2VPCSubnetNames:         subnetNames,
		},
	}

	// Add any IBM Cloud Service Endpoint overrides as necessary
	for _, endpoint := range serviceEndpoints {
		switch endpoint.Name {
		case configv1.IBMCloudServiceIAM:
			config.Provider.IAMEndpointOverride = endpoint.URL
		case configv1.IBMCloudServiceVPC:
			// Trim the version suffix ('/v1', '/v21', etc.) from the VPC URL since the CCM appends it
			suffixRegex := regexp.MustCompile(`(^.*)/v\d+[/]{0,1}$`)
			config.Provider.VPCEndpointOverride = suffixRegex.ReplaceAllString(endpoint.URL, "${1}")
		case configv1.IBMCloudServiceResourceManager:
			config.Provider.RMEndpointOverride = endpoint.URL
		default:
			logrus.Debugf("ignoring unnecessary endpoint override for cloud provider config: %s", endpoint.Name)
		}
	}

	buf := &bytes.Buffer{}
	template := template.Must(template.New("ibmcloud cloudproviderconfig").Parse(configTmpl))
	if err := template.Execute(buf, config); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// Generate a string of Subnet names for Control Plane and Compute based off the cluster name
func getVpcSubnetNames(infraID string, controlPlaneZones []string, computeZones []string) string {
	var subnetNames []string

	for cpIndex := range controlPlaneZones {
		// Add Control Plane subnet
		subnetNames = append(subnetNames, fmt.Sprintf("%s-subnet-control-plane-%s", infraID, controlPlaneZones[cpIndex]))
	}
	for comIndex := range computeZones {
		// Add Compute subnet
		subnetNames = append(subnetNames, fmt.Sprintf("%s-subnet-compute-%s", infraID, computeZones[comIndex]))
	}
	sort.Strings(subnetNames)
	return strings.Join(subnetNames, ",")
}

var configTmpl = `[global]
version = {{.Global.Version}}
[kubernetes]
config-file = {{ if ne .Kubernetes.ConfigFile "" }}{{ .Kubernetes.ConfigFile }}{{ else }}""{{ end }}
[provider]
accountID = {{.Provider.AccountID}}
clusterID = {{.Provider.ClusterID}}
cluster-default-provider = {{.Provider.ClusterDefaultProvider}}
region = {{.Provider.Region}}
g2Credentials = {{.Provider.G2CredentialsFilePath}}
g2ResourceGroupName = {{.Provider.G2ResourceGroupName}}
g2VpcName = {{.Provider.G2VPCName}}
g2workerServiceAccountID = {{.Provider.G2WorkerServiceAccountID}}
g2VpcSubnetNames = {{.Provider.G2VPCSubnetNames}}
{{ if ne .Provider.IAMEndpointOverride ""}}{{ printf "iamEndpointOverride = %s\n" .Provider.IAMEndpointOverride }}{{ end }}{{ if ne .Provider.VPCEndpointOverride ""}}{{ printf "g2EndpointOverride = %s\n" .Provider.VPCEndpointOverride }}{{ end }}{{ if ne .Provider.RMEndpointOverride ""}}{{ printf "rmEndpointOverride = %s\n" .Provider.RMEndpointOverride }}{{ end }}

`
