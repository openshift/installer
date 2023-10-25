/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	iamv1 "sigs.k8s.io/cluster-api-provider-aws/v2/iam/api/v1beta1"
)

const (
	// DefaultBootstrapUserName is the default bootstrap user name.
	DefaultBootstrapUserName = "bootstrapper.cluster-api-provider-aws.sigs.k8s.io"
	// DefaultBootstrapGroupName is the default bootstrap user name.
	DefaultBootstrapGroupName = "bootstrapper.cluster-api-provider-aws.sigs.k8s.io"
	// DefaultStackName is the default CloudFormation stack name.
	DefaultStackName = "cluster-api-provider-aws-sigs-k8s-io"
	// DefaultPartitionName is the default security partition for AWS ARNs.
	DefaultPartitionName = "aws"
	// PartitionNameUSGov is the default security partition for AWS ARNs.
	PartitionNameUSGov = "aws-us-gov"
	// DefaultKMSAliasPattern is the default KMS alias.
	DefaultKMSAliasPattern = "cluster-api-provider-aws-*"
	// DefaultS3BucketPrefix is the default S3 bucket prefix.
	DefaultS3BucketPrefix = "cluster-api-provider-aws-"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

// SetDefaults_BootstrapUser is used by defaulter-gen.
func SetDefaults_BootstrapUser(obj *BootstrapUser) { //nolint:golint,stylecheck
	if obj != nil {
		if obj.UserName == "" {
			obj.UserName = DefaultBootstrapUserName
		}
		if obj.GroupName == "" {
			obj.GroupName = DefaultBootstrapGroupName
		}
	}
}

// SetDefaults_AWSIAMConfigurationSpec is used by defaulter-gen.
func SetDefaults_AWSIAMConfigurationSpec(obj *AWSIAMConfigurationSpec) { //nolint:golint,stylecheck
	if obj.NameSuffix == nil {
		obj.NameSuffix = pointer.String(iamv1.DefaultNameSuffix)
	}
	if obj.Partition == "" {
		obj.Partition = DefaultPartitionName
	}
	if obj.StackName == "" {
		obj.StackName = DefaultStackName
	}
	if obj.EKS == nil {
		obj.EKS = &EKSConfig{
			Disable:              false,
			AllowIAMRoleCreation: false,
			DefaultControlPlaneRole: AWSIAMRoleSpec{
				Disable: false,
			},
		}
	} else if !obj.EKS.Disable {
		obj.Nodes.EC2ContainerRegistryReadOnly = true
	}
	if obj.EventBridge == nil {
		obj.EventBridge = &EventBridgeConfig{
			Enable: false,
		}
	}
	if obj.EKS.ManagedMachinePool == nil {
		obj.EKS.ManagedMachinePool = &AWSIAMRoleSpec{
			Disable: true,
		}
	}
	if obj.EKS.Fargate == nil {
		obj.EKS.Fargate = &AWSIAMRoleSpec{
			Disable: true,
		}
	}
	if len(obj.SecureSecretsBackends) == 0 {
		obj.SecureSecretsBackends = []infrav1.SecretBackend{
			infrav1.SecretBackendSecretsManager,
		}
	}
	if len(obj.EKS.KMSAliasPrefix) == 0 {
		obj.EKS.KMSAliasPrefix = DefaultKMSAliasPattern
	}

	if obj.S3Buckets.NamePrefix == "" {
		obj.S3Buckets.NamePrefix = DefaultS3BucketPrefix
	}
}

// SetDefaults_AWSIAMConfiguration is used by defaulter-gen.
func SetDefaults_AWSIAMConfiguration(obj *AWSIAMConfiguration) { //nolint:golint,stylecheck
	obj.APIVersion = SchemeGroupVersion.String()
	obj.Kind = "AWSIAMConfiguration"
	if obj.Spec.NameSuffix == nil {
		obj.Spec.NameSuffix = pointer.String(iamv1.DefaultNameSuffix)
	}
	if obj.Spec.StackName == "" {
		obj.Spec.StackName = DefaultStackName
	}
}
