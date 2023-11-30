package powervs

import (
	"bytes"
	"strings"
	"text/template"
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
	PowerVSCloudInstanceID   string `gcfg:"powerVSCloudInstanceID"`
	PowerVSCloudInstanceName string `gcfg:"powerVSCloudInstanceName"`
	PowerVSRegion            string `gcfg:"powerVSRegion"`
	PowerVSZone              string `gcfg:"powerVSZone"`
}

// CloudProviderConfig generates the cloud provider config for the IBM Power VS platform.
func CloudProviderConfig(infraID string, accountID string, vpcName string, region string, resourceGroupName string, subnets []string, cloudInstGUID string, cloudInstName string, pvsRegion string, pvsZone string) (string, error) {
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
			G2VPCSubnetNames:         strings.Join(subnets, ","),
			PowerVSCloudInstanceID:   cloudInstGUID,
			PowerVSCloudInstanceName: cloudInstName,
			PowerVSRegion:            pvsRegion,
			PowerVSZone:              pvsZone,
		},
	}
	buf := &bytes.Buffer{}
	template := template.Must(template.New("powervs cloudproviderconfig").Parse(configTmpl))
	if err := template.Execute(buf, config); err != nil {
		return "", err
	}
	return buf.String(), nil
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
powerVSCloudInstanceID = {{.Provider.PowerVSCloudInstanceID}}
powerVSCloudInstanceName = {{.Provider.PowerVSCloudInstanceName}}
powerVSRegion = {{.Provider.PowerVSRegion}}
powerVSZone = {{.Provider.PowerVSZone}}
`
