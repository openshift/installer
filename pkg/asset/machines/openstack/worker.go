// Package openstack generates Machine objects for openstack.
package openstack

import (
	"text/template"

	"github.com/openshift/installer/pkg/types/openstack"
)

// Config is used to generate the machine.
type Config struct {
	ClusterName string
	Replicas    int64
	Image       string
	Tags        map[string]string
	Region      string
	Machine     openstack.MachinePool
	Trunk       bool
}

// ComputeMachineSetTmpl is template for compute machineset.
var ComputeMachineSetTmpl = template.Must(template.New("openstack-compute-machineset").Parse(`
apiVersion: cluster.k8s.io/v1alpha1
kind: MachineSet
metadata:
  name: {{.ClusterName}}-compute-0
  namespace: openshift-cluster-api
  labels:
    sigs.k8s.io/cluster-api-cluster: {{.ClusterName}}
    sigs.k8s.io/cluster-api-machine-role: compute
    sigs.k8s.io/cluster-api-machine-type: compute
spec:
  replicas: {{.Replicas}}
  selector:
    matchLabels:
      sigs.k8s.io/cluster-api-machineset: {{.ClusterName}}-compute-0
      sigs.k8s.io/cluster-api-cluster: {{.ClusterName}}
  template:
    metadata:
      labels:
        sigs.k8s.io/cluster-api-machineset: {{.ClusterName}}-compute-0
        sigs.k8s.io/cluster-api-cluster: {{.ClusterName}}
        sigs.k8s.io/cluster-api-machine-role: compute
        sigs.k8s.io/cluster-api-machine-type: compute
    spec:
      providerSpec:
        value:
          apiVersion: openstack.cluster.k8s.io/v1alpha1
          kind: OpenStackMachineProviderConfig
          image:
            id: {{.Image}}
          flavor: {{.Machine.FlavorName}}
          placement:
            region: {{.Region}}
          subnet:
            filters:
            - name: "tag:Name"
              values:
              - "{{.ClusterName}}-compute-*"
          tags:
{{- range $key,$value := .Tags}}
            - name: "{{$key}}"
              value: "{{$value}}"
{{- end}}
          securityGroups:
            - filters:
              - name: "tag:Name"
                values:
                - "{{.ClusterName}}_compute_sg"
          userDataSecret:
            name: compute-user-data
          trunk: {{.Trunk}}
      versions:
        kubelet: ""
        controlPlane: ""
`))
