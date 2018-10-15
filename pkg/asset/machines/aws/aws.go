// Package aws generates Machine objects for aws.
package aws

import (
	"text/template"

	"github.com/openshift/installer/pkg/types"
)

// Config is used to generate the machine.
type Config struct {
	ClusterName string
	Replicas    int64
	AMIID       string
	Tags        map[string]string
	Region      string
	Machine     types.AWSMachinePoolPlatform
}

// WorkerMachineSetTmpl is template for worker machineset.
var WorkerMachineSetTmpl = template.Must(template.New("aws-worker-machineset").Parse(`
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
      sigs.k8s.io/cluster-api-machineset: worker
      sigs.k8s.io/cluster-api-cluster: {{.ClusterName}}
  template:
    metadata:
      labels:
        sigs.k8s.io/cluster-api-machineset: worker
        sigs.k8s.io/cluster-api-cluster: {{.ClusterName}}
        sigs.k8s.io/cluster-api-machine-role: worker
        sigs.k8s.io/cluster-api-machine-type: worker
    spec:
      providerConfig:
        value:
          apiVersion: aws.cluster.k8s.io/v1alpha1
          kind: AWSMachineProviderConfig
          ami:
            id: {{.AMIID}}
          instanceType: {{.Machine.InstanceType}}
          placement:
            region: {{.Region}}
          subnet:
            filters:
            - name: "tag:Name"
              values:
              - "{{.ClusterName}}-worker-*"
          iamInstanceProfile:
            id: "{{.ClusterName}}-worker-profile"
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
