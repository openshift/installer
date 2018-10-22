// Package aws generates Machine objects for aws.
package aws

import (
	"text/template"

	"github.com/openshift/installer/pkg/types"
)

// WorkerConfig is used to generate the worker machinesets.
type WorkerConfig struct {
	Instances []WorkerMachinesetInstance
	MachineConfig
}

// WorkerMachinesetInstance constains information specific to
// one worker machineset.
type WorkerMachinesetInstance struct {
	Replicas         int64
	AvailabilityZone string
}

// MachineConfig contains fields common to worker and master
// machine configurations
type MachineConfig struct {
	ClusterName string
	AMIID       string
	Tags        map[string]string
	Region      string
	Machine     types.AWSMachinePoolPlatform
}

// WorkerMachineSetsTmpl is template for worker machinesets.
var WorkerMachineSetsTmpl = template.Must(template.New("aws-worker-machinesets").Parse(`
{{- $c := . -}}
kind: List
apiVersion: v1
metadata:
  resourceVersion: ""
  selfLink: ""
items:
{{- range $index,$instance := $c.Instances}}
- apiVersion: cluster.k8s.io/v1alpha1
  kind: MachineSet
  metadata:
    name: {{$c.ClusterName}}-worker-{{$index}}
    namespace: openshift-cluster-api
    labels:
      sigs.k8s.io/cluster-api-cluster: {{$c.ClusterName}}
      sigs.k8s.io/cluster-api-machine-role: worker
      sigs.k8s.io/cluster-api-machine-type: worker
  spec:
    replicas: {{$instance.Replicas}}
    selector:
      matchLabels:
        sigs.k8s.io/cluster-api-machineset: worker
        sigs.k8s.io/cluster-api-cluster: {{$c.ClusterName}}
    template:
      metadata:
        labels:
          sigs.k8s.io/cluster-api-machineset: worker
          sigs.k8s.io/cluster-api-cluster: {{$c.ClusterName}}
          sigs.k8s.io/cluster-api-machine-role: worker
          sigs.k8s.io/cluster-api-machine-type: worker
      spec:
        providerConfig:
          value:
            apiVersion: aws.cluster.k8s.io/v1alpha1
            kind: AWSMachineProviderConfig
            ami:
              id: {{$c.AMIID}}
            instanceType: {{$c.Machine.InstanceType}}
            placement:
              region: {{$c.Region}}
              availabilityZone: {{$instance.AvailabilityZone}}
            subnet:
              filters:
              - name: "tag:Name"
                values:
                - "{{$c.ClusterName}}-worker-{{$instance.AvailabilityZone}}"
            iamInstanceProfile:
              id: "{{$c.ClusterName}}-worker-profile"
            tags:
              {{- range $key,$value := $c.Tags}}
              - name: "{{$key}}"
                value: "{{$value}}"
              {{- end}}
            securityGroups:
              - filters:
                - name: "tag:Name"
                  values:
                  - "{{$c.ClusterName}}_worker_sg"
            userDataSecret:
              name: worker-user-data
        versions:
          kubelet: ""
          controlPlane: ""
{{- end}}
`))
