// Package openstack generates Machine objects for openstack.
package openstack

import (
	"text/template"

	"github.com/openshift/installer/pkg/types/openstack"
)

// Config is used to generate the machine.
type Config struct {
	CloudName   string
	ClusterName string
	Replicas    int64
	Image       string
	Tags        map[string]string
	Region      string
	Machine     openstack.MachinePool
	Trunk       bool
}

// WorkerMachineSetTmpl is template for worker machineset.
var WorkerMachineSetTmpl = template.Must(template.New("openstack-worker-machineset").Parse(`
apiVersion: cluster.k8s.io/v1alpha1
kind: MachineSet
metadata:
  name: {{.ClusterName}}-worker
  namespace: openshift-cluster-api
  labels:
    sigs.k8s.io/cluster-api-cluster: {{.ClusterName}}
    sigs.k8s.io/cluster-api-machine-role: worker
    sigs.k8s.io/cluster-api-machine-type: worker
spec:
  replicas: {{.Replicas}}
  selector:
    matchLabels:
      sigs.k8s.io/cluster-api-machineset: {{.ClusterName}}-worker
      sigs.k8s.io/cluster-api-cluster: {{.ClusterName}}
  template:
    metadata:
      labels:
        sigs.k8s.io/cluster-api-machineset: {{.ClusterName}}-worker
        sigs.k8s.io/cluster-api-cluster: {{.ClusterName}}
        sigs.k8s.io/cluster-api-machine-role: worker
        sigs.k8s.io/cluster-api-machine-type: worker
    spec:
      providerSpec:
        value:
          apiVersion: openstack.cluster.k8s.io/v1alpha1
          kind: OpenStackMachineProviderConfig
          cloudName: {{.CloudName}}
          cloudsSecret: "openstack-credentials"
          image: {{.Image}}
          flavor: {{.Machine.FlavorName}}
          placement:
            region: {{.Region}}
          networks:
{{- range $key,$value := .Tags}}
          - filter:
              tags: "{{$key}}={{$value}}"
{{- end}}
          securityGroups:
            - worker
          userDataSecret:
            name: worker-user-data
          trunk: {{.Trunk}}
      versions:
        kubelet: "v1.11.0"
        controlPlane: ""
`))
