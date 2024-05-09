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
}

// CloudProviderConfig generates the cloud provider config for the GCP platform.
func CloudProviderConfig(infraID, projectID, subnet, networkProjectID string) (string, error) {
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
{{ if ne .Global.NetworkProjectID "" }}network-project-id = {{.Global.NetworkProjectID}}{{end}}

`
