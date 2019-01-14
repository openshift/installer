// Package openstack generates Machine objects for openstack.
package openstack

import (
	"text/template"

	"github.com/openshift/installer/pkg/types/openstack"
)

// MasterConfig is used to generate the machine.
type MasterConfig struct {
	CloudName   string
	ClusterName string
	Instances   []string
	Image       string
	Tags        map[string]string
	Region      string
	Machine     openstack.MachinePool
	Trunk       bool
}

// MasterMachinesTmpl is the template for master machines.
var MasterMachinesTmpl = template.Must(template.New("openstack-master-machines").Parse(`
{{- $c := . -}}
kind: List
apiVersion: v1
metadata:
  resourceVersion: ""
  selfLink: ""
items:
{{- range $index,$instance := .Instances}}
- apiVersion: cluster.k8s.io/v1alpha1
  kind: Machine
  metadata:
    name: {{$c.ClusterName}}-master-{{$index}}
    namespace: openshift-cluster-api
    labels:
      sigs.k8s.io/cluster-api-cluster: {{$c.ClusterName}}
      sigs.k8s.io/cluster-api-machine-role: master
      sigs.k8s.io/cluster-api-machine-type: master
  spec:
    providerSpec:
      value:
        apiVersion: openstack.cluster.k8s.io/v1alpha1
        kind: OpenStackMachineProviderConfig
        cloudName: {{$c.CloudName}}
        cloudsSecret: "openstack-credentials"
        image: {{$c.Image}}
        flavor: {{$c.Machine.FlavorName}}
        placement:
          region: {{$c.Region}}
        networks:
{{- range $key,$value := $c.Tags}}
        - filter:
            tags: "{{$key}}={{$value}}"
{{- end}}
        securityGroups:
          - master
        userDataSecret:
          name: master-user-data
        trunk: {{$c.Trunk}}
    versions:
      kubelet: "v1.11.0"
      controlPlane: ""
{{- end -}}
`))
