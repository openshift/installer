// Package aws generates Machine objects for aws.
package aws

import (
	"text/template"
)

// MasterConfig is used to generate master machines
type MasterConfig struct {
	MachineConfig
	Instances []MasterInstance
}

// MasterInstance contains information specific to each
// master machine instance to create.
type MasterInstance struct {
	AvailabilityZone string
}

// MasterMachineTmpl is a template for a list of master machines.
var MasterMachineTmpl = template.Must(template.New("aws-master-machine").Parse(`
{{- $c := . -}}
kind: List
apiVersion: v1
metadata:
  resourceVersion: ""
  selfLink: ""
items:
{{- range $index,$instance := $c.Instances}}
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
            - "{{$c.ClusterName}}-master-{{$instance.AvailabilityZone}}"
        iamInstanceProfile:
          id: "{{$c.ClusterName}}-master-profile"
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
          name: "master-user-data-{{$index}}"
    versions:
      kubelet: ""
      controlPlane: ""
{{- end}}
`))
