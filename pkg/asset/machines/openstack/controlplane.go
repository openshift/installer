// Package openstack generates Machine objects for openstack.
package openstack

import (
	"text/template"

	"github.com/openshift/installer/pkg/types/openstack"
)

// ControlPlaneConfig is used to generate the machine.
type ControlPlaneConfig struct {
	ClusterName string
	Instances   []string
	Image       string
	Tags        map[string]string
	Region      string
	Machine     openstack.MachinePool
	Trunk       bool
}

// ControlPlaneMachinesTmpl is the template for control plane machines.
var ControlPlaneMachinesTmpl = template.Must(template.New("openstack-controlplane-machines").Parse(`
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
    name: {{$c.ClusterName}}-controlplane-{{$index}}
    namespace: openshift-cluster-api
    labels:
      sigs.k8s.io/cluster-api-cluster: {{$c.ClusterName}}
      sigs.k8s.io/cluster-api-machine-role: controlplane
      sigs.k8s.io/cluster-api-machine-type: controlplane
  spec:
    providerSpec:
      value:
        apiVersion: openstack.cluster.k8s.io/v1alpha1
        kind: OpenStackMachineProviderConfig
        image:
          id: {{$c.Image}}
        flavor: {{$c.Machine.FlavorName}}
        placement:
          region: {{$c.Region}}
        subnet:
          filters:
          - name: "tag:Name"
            values:
            - "{{$c.ClusterName}}-controlplane-*"
        tags:
{{- range $key,$value := $c.Tags}}
          - name: "{{$key}}"
            value: "{{$value}}"
{{- end}}
        securityGroups:
          - filters:
            - name: "tag:Name"
              values:
              - "{{$c.ClusterName}}_controlplane_sg"
        userDataSecret:
          name: controlplane-user-data
        trunk: {{$c.Trunk}}
    versions:
      kubelet: ""
      controlPlane: ""
{{- end -}}
`))
