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

	NodeTags []string `gcfg:"node-tags"`

	SubnetworkName string `gcfg:"subnetwork-name"`
}

// CloudProviderConfig generates the cloud provider config for the GCP platform.
func CloudProviderConfig(infraID, projectID, subnet string) (string, error) {
	config := &config{
		Global: global{
			ProjectID: projectID,

			// To make sure k8s cloud provider is looking for instances in all zones.
			Regional:  true,
			Multizone: true,

			// To make sure k8s cloud provide has tags for firewal for load balancer.
			NodeTags: []string{fmt.Sprintf("%s-master", infraID), fmt.Sprintf("%s-worker", infraID)},

			// Used for internal load balancers
			SubnetworkName: subnet,
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
subnetwork-name = {{.Global.SubnetworkName}}

`
