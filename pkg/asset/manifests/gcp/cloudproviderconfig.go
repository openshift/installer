package gcp

import (
	"bytes"
	"fmt"
	"text/template"
)

// https://github.com/kubernetes/kubernetes/blob/368ee4bb8ee7a0c18431cd87ee49f0c890aa53e5/staging/src/k8s.io/legacy-cloud-providers/gce/gce.go#L188
type config struct {
	Global global `gcfg:"global"`
}

type global struct {
	ProjectID string `gcfg:"project-id"`

	Regional  bool `gcfg:"regional"`
	Multizone bool `gcfg:"multizone"`

	NodeTags                     []string `gcfg:"node-tags"`
	NodeInstancePrefix           string   `gcfg:"node-instance-prefix"`
	ExternalInstanceGroupsPrefix string   `gcfg:"external-instance-groups-prefix"`

	SubnetworkName string `gcfg:"subnetwork-name"`

	NetworkProjectID string `gcfg:"network-project-id"`

	// APIEndpoint is the compute API endpoint to use. If this is blank,
	// then the default endpoint is used.
	APIEndpoint string `gcfg:"api-endpoint"`
	// ContainerAPIEndpoint is the container API endpoint to use. If this is blank,
	// then the default endpoint is used.
	ContainerAPIEndpoint string `gcfg:"container-api-endpoint"`

	FirewallManagement string `gcfg:"firewall-rules-management"`
}

// CloudProviderConfig generates the cloud provider config for the GCP platform.
func CloudProviderConfig(infraID, projectID, subnet, networkProjectID, apiEndpoint, containerAPIEndpoint, firewallManagement string) (string, error) {
	config := &config{
		Global: global{
			ProjectID: projectID,

			// To make sure k8s cloud provider is looking for instances in all zones.
			Regional:  true,
			Multizone: true,

			// To make sure k8s cloud provider has tags for firewall for load balancer.
			// The CAPI gcp provider uses the node tag "control-plane" for master nodes.
			NodeTags:                     []string{fmt.Sprintf("%s-master", infraID), fmt.Sprintf("%s-control-plane", infraID), fmt.Sprintf("%s-worker", infraID)},
			NodeInstancePrefix:           infraID,
			ExternalInstanceGroupsPrefix: infraID,

			// Used for internal load balancers
			SubnetworkName: subnet,

			// Used for shared vpc installations,
			NetworkProjectID: networkProjectID,

			// Used for api endpoint overrides in the cloud provider.
			APIEndpoint:          apiEndpoint,
			ContainerAPIEndpoint: containerAPIEndpoint,

			FirewallManagement: firewallManagement,
		},
	}

	buf := &bytes.Buffer{}
	template := template.Must(template.New("gce cloudproviderconfig").Parse(configTmpl))
	if err := template.Execute(buf, config); err != nil {
		return "", err
	}
	return buf.String(), nil
}

var configTmpl = `[global]
project-id      = {{.Global.ProjectID}}
regional        = {{.Global.Regional}}
multizone       = {{.Global.Multizone}}
{{range $idx, $tag := .Global.NodeTags -}}
node-tags       = {{$tag}}
{{end -}}
node-instance-prefix = {{.Global.NodeInstancePrefix}}
external-instance-groups-prefix = {{.Global.ExternalInstanceGroupsPrefix}}
subnetwork-name = {{.Global.SubnetworkName}}
{{- if ne .Global.NetworkProjectID "" }}{{"\n"}}network-project-id = {{.Global.NetworkProjectID}}{{ end }}
{{- if ne .Global.APIEndpoint "" }}{{"\n"}}api-endpoint = {{.Global.APIEndpoint}}{{ end }}
{{- if ne .Global.ContainerAPIEndpoint "" }}{{"\n"}}container-api-endpoint = {{.Global.ContainerAPIEndpoint}}{{ end }}
{{- if ne .Global.FirewallManagement "" }}{{"\n"}}firewall-rules-management = {{.Global.FirewallManagement}}{{ end }}

`
