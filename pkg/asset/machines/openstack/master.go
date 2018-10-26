// Package openstack generates Machine objects for openstack.
package openstack

import (
	"text/template"

	"github.com/openshift/installer/pkg/types"
)

// MasterConfig is used to generate the machine.
type MasterConfig struct {
	ClusterName string
	Instances   []string
	Image       string
	Tags        map[string]string
	Region      string
	Machine     types.OpenStackMachinePoolPlatform
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
    providerConfig:
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
            - "{{$c.ClusterName}}-master-*"
        tags:
{{- range $key,$value := $c.Tags}}
          - name: "{{$key}}"
            value: "{{$value}}"
{{- end}}
        securityGroups:
          - filters:
            - name: "tag:Name"
              values:
              - "{{$c.ClusterName}}_master_sg"
        userDataSecret:
          name: master-user-data
    versions:
      kubelet: ""
      controlPlane: ""
{{- end -}}
`))
