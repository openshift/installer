package ibmcloud

import (
	"bytes"
	"text/template"

	"github.com/IBM-Cloud/bluemix-go/crn"
	"github.com/openshift/installer/pkg/types"
)

// https://github.com/kubernetes/kubernetes/blob/368ee4bb8ee7a0c18431cd87ee49f0c890aa53e5/staging/src/k8s.io/legacy-cloud-providers/gce/gce.go#L188
type config struct {
	Global                 global                 `gcfg:"global"`
	Kubernetes             kubernetes             `gcfg:"kubernetes"`
	LoadBalancerDeployment loadBalancerDeployment `gcfg:"load-balancer-deployment"`
	Provider               provider               `gcfg:"provider"`
}

type global struct {
	Version string `gcfg:"version"`
}

type kubernetes struct {
	ConfigFile string `gcfg:"config-file"`
}

type loadBalancerDeployment struct {
	Image           string `gcfg:"image"`
	Application     string `gcfg:"application"`
	VLANIPConfigMap string `gcfg:"vlan-ip-config-map"`
}

type provider struct {
	AccountID string `gcfg:"accountID"`
	ClusterID string `gcfg:"clusterID"`
}

// CloudProviderConfig generates the cloud provider config for the IBMCloud platform.
func CloudProviderConfig(infraID string, installConfig types.InstallConfig) (string, error) {
	CRN, err := crn.Parse(installConfig.Platform.IBMCloud.CISInstanceCRN)
	if err != nil {
		return "", err
	}

	config := &config{
		Global: global{
			Version: "1.1.0",
		},
		Kubernetes: kubernetes{
			ConfigFile: "/mnt/etc/kubernetes/controller-manager-kubeconfig",
		},
		LoadBalancerDeployment: loadBalancerDeployment{
			Image:           "[REGISTRY]/[NAMESPACE]/keepalived:[TAG]",
			Application:     "keepalived",
			VLANIPConfigMap: "ibm-cloud-provider-vlan-ip-config",
		},
		Provider: provider{
			AccountID: CRN.Scope,
			ClusterID: infraID,
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
config-file = {{.Kubernetes.ConfigFile}}
[load-balancer-deployment]
image = {{.LoadBalancerDeployment.Image}}
application = {{.LoadBalancerDeployment.Application}}
vlan-ip-config-map = {{.LoadBalancerDeployment.VLANIPConfigMap}}
[provider]
accountID = {{.Provider.AccountID}}
clusterID = {{.Provider.ClusterID}}

`
