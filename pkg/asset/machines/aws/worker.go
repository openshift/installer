// Package aws generates Machine objects for aws.
package aws

import (
	"fmt"
	"text/template"

        awsprovider "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsproviderconfig/v1alpha1"
        clusterapi "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

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

func createWorkerMachineSets(config *WorkerConfig) metav1.List {
	list := metav1.List{}
	list.Kind = "List"
	items := make([]runtime.RawExtension, 0)
	for i, instance := range config.Instances {
		ms := clusterapi.MachineSet{}
		ms.Kind = "MachineSet"
		ms.Name = fmt.Sprintf("%s-worker-%d", config.ClusterName, i)
		ms.Namespace = "openshift-cluster-api"
		ms.Labels = map[string]string{
			"sigs.k8s.io/cluster-api-cluster":      config.ClusterName,
			"sigs.k8s.io/cluster-api-machine-role": "worker",
			"sigs.k8s.io/cluster-api-machine-type": "worker",
		}
		ms.Spec = clusterapi.MachineSetSpec{}
		ms.Spec.Replicas = new(int32)
		*ms.Spec.Replicas = int32(instance.Replicas)
		ms.Spec.Selector = metav1.LabelSelector{}
		ms.Spec.Selector.MatchLabels = map[string]string{
			"sigs.k8s.io/cluster-api-machineset": "worker",
			"sigs.k8s.io/cluster-api-cluster":    config.ClusterName,
		}

		template := clusterapi.MachineTemplateSpec{}
		template.Labels = map[string]string{
			"sigs.k8s.io/cluster-api-machineset":   "worker",
			"sigs.k8s.io/cluster-api-cluster":      config.ClusterName,
			"sigs.k8s.io/cluster-api-machine-role": "worker",
			"sigs.k8s.io/cluster-api-machine-type": "worker",
		}
		template.Spec = clusterapi.MachineSpec{}
		provider := awsprovider.AWSMachineProviderConfig{}
		provider.APIVersion = "aws.cluster.k8s.io/v1alpha1"
		provider.InstanceType = config.Machine.InstanceType
		provider.AMI = awsprovider.AWSResourceReference{ID: new(string)}
		*provider.AMI.ID = config.AMIID
		provider.Subnet = awsprovider.AWSResourceReference{
			Filters: []awsprovider.Filter{
				{
					Name:   "tag:Name",
					Values: []string{fmt.Sprintf("%s-worker-%s", config.ClusterName, instance.AvailabilityZone)},
				},
			},
		}
		provider.Placement = awsprovider.Placement{Region: config.Region, AvailabilityZone: instance.AvailabilityZone}
		provider.IAMInstanceProfile = &awsprovider.AWSResourceReference{ID: new(string)}
		*provider.IAMInstanceProfile.ID = fmt.Sprintf("%s-worker-profile", config.ClusterName)
		provider.Tags = make([]awsprovider.TagSpecification, 0)
		for tagName, tagValue := range config.Tags {
			provider.Tags = append(provider.Tags, awsprovider.TagSpecification{Name: tagName, Value: tagValue})
		}
		provider.SecurityGroups = []awsprovider.AWSResourceReference{
			{
				Filters: []awsprovider.Filter{
					{
						Name:   "tag:Name",
						Values: []string{fmt.Sprintf("%s_worker_sg", config.ClusterName)},
					},
				},
			},
		}
		provider.UserDataSecret = &corev1.LocalObjectReference{Name: "worker-user-data"}

		template.Spec.ProviderConfig = clusterapi.ProviderConfig{Value: &runtime.RawExtension{Object: &provider}}
		template.Spec.Versions = clusterapi.MachineVersionInfo{
			Kubelet:      "",
			ControlPlane: "",
		}
		ms.Spec.Template = template

		// add the machineset to items
		items = append(items, runtime.RawExtension{Object: &ms})
	}
	list.Items = items
	return list
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
        sigs.k8s.io/cluster-api-machineset: {{$c.ClusterName}}-worker-{{$index}}
        sigs.k8s.io/cluster-api-cluster: {{$c.ClusterName}}
    template:
      metadata:
        labels:
          sigs.k8s.io/cluster-api-machineset: {{$c.ClusterName}}-worker-{{$index}}
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
