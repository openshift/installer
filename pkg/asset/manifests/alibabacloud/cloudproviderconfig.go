package alibabacloud

import (
	"bytes"
	"text/template"
)

// ref: https://github.com/kubernetes/kubernetes/blob/368ee4bb8ee7a0c18431cd87ee49f0c890aa53e5/staging/src/k8s.io/legacy-cloud-providers/gce/gce.go#L188
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
	AccessKeyID     string `gcfg:"accessKeyID"`
	AccessKeySecret string `gcfg:"accessKeySecret"`
	ClusterID       string `gcfg:"clusterID"`
}

// CloudProviderConfig generates the cloud provider config for the Alibaba Cloud platform.
func CloudProviderConfig(infraID string, accessKeyID string, accessKeySecret string) (string, error) {
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
			VLANIPConfigMap: "alibaba-cloud-provider-vlan-ip-config",
		},
		Provider: provider{
			AccessKeyID:     accessKeyID,
			AccessKeySecret: accessKeySecret,
			ClusterID:       infraID,
		},
	}
	buf := &bytes.Buffer{}
	template := template.Must(template.New("alibaba cloudproviderconfig").Parse(configTmpl))
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
clusterID = {{.Provider.ClusterID}}

`
