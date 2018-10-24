// Package openstack generates Machine objects for openstack.
package openstack

import (
	"text/template"

	"github.com/openshift/installer/pkg/types"
)

// Config is used to generate the machine.
type Config struct {
	ClusterName string
	Replicas    int64
	Image       string
	Tags        map[string]string
	Region      string
	Machine     types.OpenStackMachinePoolPlatform
}

// WorkerMachineSetTmpl is template for worker machineset.
var WorkerMachineSetTmpl = template.Must(template.New("openstack-worker-machineset").Parse(`
apiVersion: cluster.k8s.io/v1alpha1
kind: MachineSet
metadata:
  name: {{.ClusterName}}-worker-0
  namespace: openshift-cluster-api
  labels:
    sigs.k8s.io/cluster-api-cluster: {{.ClusterName}}
    sigs.k8s.io/cluster-api-machine-role: worker
    sigs.k8s.io/cluster-api-machine-type: worker
spec:
  replicas: {{.Replicas}}
  selector:
    matchLabels:
      sigs.k8s.io/cluster-api-machineset: {{.ClusterName}}-worker-0
      sigs.k8s.io/cluster-api-cluster: {{.ClusterName}}
  template:
    metadata:
      labels:
        sigs.k8s.io/cluster-api-machineset: {{.ClusterName}}-worker-0
        sigs.k8s.io/cluster-api-cluster: {{.ClusterName}}
        sigs.k8s.io/cluster-api-machine-role: worker
        sigs.k8s.io/cluster-api-machine-type: worker
    spec:
      providerConfig:
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
              - "{{.ClusterName}}-worker-*"
          tags:
{{- range $key,$value := .Tags}}
            - name: "{{$key}}"
              value: "{{$value}}"
{{- end}}
          securityGroups:
            - filters:
              - name: "tag:Name"
                values:
                - "{{.ClusterName}}_worker_sg"
          userDataSecret:
            name: worker-user-data
      versions:
        kubelet: ""
        controlPlane: ""
`))
