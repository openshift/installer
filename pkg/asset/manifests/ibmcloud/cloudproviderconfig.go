package ibmcloud

import (
	"bytes"
	"fmt"
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
	G2VPCName                string `gcfg:"g2VpcName"`
	G2WorkerServiceAccountID string `gcfg:"g2workerServiceAccountID"`
}

// CloudProviderConfig generates the cloud provider config for the IBMCloud platform.
func CloudProviderConfig(infraID string, accountID string, region string) (string, error) {
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
			G2VPCName:                fmt.Sprintf("%s-vpc", infraID),
			G2WorkerServiceAccountID: accountID,
		},
	}
	buf := &bytes.Buffer{}
	template := template.Must(template.New("ibmcloud cloudproviderconfig").Parse(configTmpl))
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
g2VpcName = {{.Provider.G2VPCName}}
g2workerServiceAccountID = {{.Provider.G2WorkerServiceAccountID}}

`
