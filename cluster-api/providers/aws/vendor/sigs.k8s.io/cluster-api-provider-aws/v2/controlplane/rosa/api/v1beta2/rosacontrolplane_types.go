/*
Copyright 2023 The Kubernetes Authors.

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

package v1beta2

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// RosaEndpointAccessType specifies the publishing scope of cluster endpoints.
type RosaEndpointAccessType string

const (
	// Public endpoint access allows public API server access and
	// private node communication with the control plane.
	Public RosaEndpointAccessType = "Public"

	// Private endpoint access allows only private API server access and private
	// node communication with the control plane.
	Private RosaEndpointAccessType = "Private"
)

// VersionGateAckType specifies the version gate acknowledgment.
type VersionGateAckType string

const (
	// Acknowledge if acknowledgment is required and proceed with the upgrade.
	Acknowledge VersionGateAckType = "Acknowledge"

	// WaitForAcknowledge if acknowledgment is required, wait not to proceed with the upgrade.
	WaitForAcknowledge VersionGateAckType = "WaitForAcknowledge"

	// AlwaysAcknowledge always acknowledg if required and proceed with the upgrade.
	AlwaysAcknowledge VersionGateAckType = "AlwaysAcknowledge"
)

// ChannelGroupType specifies the OpenShift version channel group.
type ChannelGroupType string

const (
	// Stable channel group is the default channel group for stable releases.
	Stable ChannelGroupType = "stable"

	// Eus channel group is for eus channel releases.
	Eus ChannelGroupType = "eus"

	// Fast channel group is for fast channel releases.
	Fast ChannelGroupType = "fast"

	// Candidate channel group is for testing candidate builds.
	Candidate ChannelGroupType = "candidate"

	// Nightly channel group is for testing nigtly builds.
	Nightly ChannelGroupType = "nightly"
)

// AutoNodeMode specifies the AutoNode mode for the ROSA Control Plane.
type AutoNodeMode string

const (
	// AutoNodeModeEnabled enable AutoNode
	AutoNodeModeEnabled AutoNodeMode = "Enabled"

	// AutoNodeModeDisabled Disabled AutoNode
	AutoNodeModeDisabled AutoNodeMode = "Disabled"
)

// RosaControlPlaneSpec defines the desired state of ROSAControlPlane.
type RosaControlPlaneSpec struct { //nolint: maligned
	// Cluster name must be valid DNS-1035 label, so it must consist of lower case alphanumeric
	// characters or '-', start with an alphabetic character, end with an alphanumeric character
	// and have a max length of 54 characters.
	//
	// +immutable
	// +kubebuilder:validation:XValidation:rule="self == oldSelf", message="rosaClusterName is immutable"
	// +kubebuilder:validation:MaxLength:=54
	// +kubebuilder:validation:Pattern:=`^[a-z]([-a-z0-9]*[a-z0-9])?$`
	RosaClusterName string `json:"rosaClusterName"`

	// DomainPrefix is an optional prefix added to the cluster's domain name. It will be used
	// when generating a sub-domain for the cluster on openshiftapps domain. It must be valid DNS-1035 label
	// consisting of lower case alphanumeric characters or '-', start with an alphabetic character
	// end with an alphanumeric character and have a max length of 15 characters.
	//
	// +immutable
	// +kubebuilder:validation:XValidation:rule="self == oldSelf", message="domainPrefix is immutable"
	// +kubebuilder:validation:MaxLength:=15
	// +kubebuilder:validation:Pattern:=`^[a-z]([-a-z0-9]*[a-z0-9])?$`
	// +optional
	DomainPrefix string `json:"domainPrefix,omitempty"`

	// The Subnet IDs to use when installing the cluster.
	// SubnetIDs should come in pairs; two per availability zone, one private and one public.
	// +optional
	Subnets []string `json:"subnets,omitempty"`

	// AvailabilityZones describe AWS AvailabilityZones of the worker nodes.
	// should match the AvailabilityZones of the provided Subnets.
	// a machinepool will be created for each availabilityZone.
	// +optional
	AvailabilityZones []string `json:"availabilityZones,omitempty"`

	// The AWS Region the cluster lives in.
	Region string `json:"region"`

	// OpenShift semantic version, for example "4.14.5".
	Version string `json:"version"`

	// OpenShift version channel group, default is stable.
	//
	// +kubebuilder:validation:Enum=stable;eus;fast;candidate;nightly
	// +kubebuilder:default=stable
	ChannelGroup ChannelGroupType `json:"channelGroup"`

	// VersionGate requires acknowledgment when upgrading ROSA-HCP y-stream versions (e.g., from 4.15 to 4.16).
	// Default is WaitForAcknowledge.
	// WaitForAcknowledge: If acknowledgment is required, the upgrade will not proceed until VersionGate is set to Acknowledge or AlwaysAcknowledge.
	// Acknowledge: If acknowledgment is required, apply it for the upgrade. After upgrade is done set the version gate to WaitForAcknowledge.
	// AlwaysAcknowledge: If acknowledgment is required, apply it and proceed with the upgrade.
	//
	// +kubebuilder:validation:Enum=Acknowledge;WaitForAcknowledge;AlwaysAcknowledge
	// +kubebuilder:default=WaitForAcknowledge
	VersionGate VersionGateAckType `json:"versionGate"`

	// RosaRoleConfigRef is a reference to a RosaRoleConfig resource that contains account roles, operator roles and OIDC configuration.
	// RosaRoleConfigRef and role fields such as installerRoleARN, supportRoleARN, workerRoleARN, rolesRef and oidcID are mutually exclusive.
	//
	// +optional
	RosaRoleConfigRef *corev1.LocalObjectReference `json:"rosaRoleConfigRef,omitempty"`

	// AWS IAM roles used to perform credential requests by the openshift operators.
	// Required if RosaRoleConfigRef is not specified.
	// +optional
	RolesRef AWSRolesRef `json:"rolesRef,omitempty"`

	// The ID of the internal OpenID Connect Provider.
	// Required if RosaRoleConfigRef is not specified.
	//
	// +kubebuilder:validation:XValidation:rule="self == oldSelf", message="oidcID is immutable"
	// +optional
	OIDCID string `json:"oidcID,omitempty"`

	// EnableExternalAuthProviders enables external authentication configuration for the cluster.
	//
	// +kubebuilder:default=false
	// +kubebuilder:validation:XValidation:rule="self == oldSelf", message="enableExternalAuthProviders is immutable"
	// +optional
	EnableExternalAuthProviders bool `json:"enableExternalAuthProviders,omitempty"`

	// ExternalAuthProviders are external OIDC identity providers that can issue tokens for this cluster.
	// Can only be set if "enableExternalAuthProviders" is set to "True".
	//
	// At most one provider can be configured.
	//
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MaxItems=1
	ExternalAuthProviders []ExternalAuthProvider `json:"externalAuthProviders,omitempty"`

	// InstallerRoleARN is an AWS IAM role that OpenShift Cluster Manager will assume to create the cluster.
	// Required if RosaRoleConfigRef is not specified.
	// +optional
	InstallerRoleARN string `json:"installerRoleARN,omitempty"`
	// SupportRoleARN is an AWS IAM role used by Red Hat SREs to enable
	// access to the cluster account in order to provide support.
	// Required if RosaRoleConfigRef is not specified.
	// +optional
	SupportRoleARN string `json:"supportRoleARN,omitempty"`
	// WorkerRoleARN is an AWS IAM role that will be attached to worker instances.
	// Required if RosaRoleConfigRef is not specified.
	// +optional
	WorkerRoleARN string `json:"workerRoleARN,omitempty"`

	// BillingAccount is an optional AWS account to use for billing the subscription fees for ROSA HCP clusters.
	// The cost of running each ROSA HCP cluster will be billed to the infrastructure account in which the cluster
	// is running.
	//
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf", message="billingAccount is immutable"
	// +kubebuilder:validation:XValidation:rule="self.matches('^[0-9]{12}$')", message="billingAccount must be a valid AWS account ID"
	// +immutable
	// +optional
	BillingAccount string `json:"billingAccount,omitempty"`

	// DefaultMachinePoolSpec defines the configuration for the default machinepool(s) provisioned as part of the cluster creation.
	// One MachinePool will be created with this configuration per AvailabilityZone. Those default machinepools are required for openshift cluster operators
	// to work properly.
	// As these machinepool not created using ROSAMachinePool CR, they will not be visible/managed by ROSA CAPI provider.
	// `rosa list machinepools -c <rosaClusterName>` can be used to view those machinepools.
	//
	// This field will be removed in the future once the current limitation is resolved.
	//
	// +optional
	DefaultMachinePoolSpec DefaultMachinePoolSpec `json:"defaultMachinePoolSpec,omitempty"`

	// Network config for the ROSA HCP cluster.
	// +optional
	Network *NetworkSpec `json:"network,omitempty"`

	// EndpointAccess specifies the publishing scope of cluster endpoints. The
	// default is Public.
	//
	// +kubebuilder:validation:Enum=Public;Private
	// +kubebuilder:default=Public
	// +optional
	EndpointAccess RosaEndpointAccessType `json:"endpointAccess,omitempty"`

	// AdditionalTags are user-defined tags to be added on the AWS resources associated with the control plane.
	// +optional
	AdditionalTags infrav1.Tags `json:"additionalTags,omitempty"`

	// EtcdEncryptionKMSARN is the ARN of the KMS key used to encrypt etcd. The key itself needs to be
	// created out-of-band by the user and tagged with `red-hat:true`.
	// +optional
	EtcdEncryptionKMSARN string `json:"etcdEncryptionKMSARN,omitempty"`

	// AuditLogRoleARN defines the role that is used to forward audit logs to AWS CloudWatch.
	// If not set, audit log forwarding is disabled.
	// +optional
	AuditLogRoleARN string `json:"auditLogRoleARN,omitempty"`

	// ProvisionShardID defines the shard where ROSA hosted control plane components will be hosted.
	//
	// +kubebuilder:validation:XValidation:rule="self == oldSelf", message="provisionShardID is immutable"
	// +optional
	ProvisionShardID string `json:"provisionShardID,omitempty"`

	// CredentialsSecretRef references a secret with necessary credentials to connect to the OCM API.
	// The secret should contain the following data keys:
	// - ocmToken: eyJhbGciOiJIUzI1NiIsI....
	// - ocmApiUrl: Optional, defaults to 'https://api.openshift.com'
	// +optional
	CredentialsSecretRef *corev1.LocalObjectReference `json:"credentialsSecretRef,omitempty"`

	// IdentityRef is a reference to an identity to be used when reconciling the managed control plane.
	// If no identity is specified, the default identity for this controller will be used.
	//
	// +optional
	IdentityRef *infrav1.AWSIdentityReference `json:"identityRef,omitempty"`

	// ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// +optional
	ControlPlaneEndpoint clusterv1.APIEndpoint `json:"controlPlaneEndpoint"`

	// ClusterRegistryConfig represents registry config used with the cluster.
	// +optional
	ClusterRegistryConfig *RegistryConfig `json:"clusterRegistryConfig,omitempty"`

	// autoNode set the autoNode mode and roleARN.
	// +optional
	AutoNode *AutoNode `json:"autoNode,omitempty"`

	// ROSANetworkRef references ROSANetwork custom resource that contains the networking infrastructure
	// for the ROSA HCP cluster.
	// +optional
	ROSANetworkRef *corev1.LocalObjectReference `json:"rosaNetworkRef,omitempty"`
}

// AutoNode set the AutoNode mode and AutoNode role ARN.
type AutoNode struct {
	// mode specifies the mode for the AutoNode. Setting Enable/Disable mode will allows/disallow karpenter AutoNode scaling.
	// +kubebuilder:validation:Enum=Enabled;Disabled
	// +kubebuilder:default=Disabled
	// +optional
	Mode AutoNodeMode `json:"mode,omitempty"`

	// roleARN sets the autoNode role ARN, which includes the IAM policy and cluster-specific role that grant the necessary permissions to the Karpenter controller.
	// The role must be attached with the same OIDC-ID that is used with the ROSA-HCP cluster.
	// +kubebuilder:validation:MaxLength:=2048
	// +optional
	RoleARN string `json:"roleARN,omitempty"`
}

// RegistryConfig for ROSA-HCP cluster
type RegistryConfig struct {
	// AdditionalTrustedCAs containing the registry hostname as the key, and the PEM-encoded certificate as the value,
	// for each additional registry CA to trust.
	// +optional
	AdditionalTrustedCAs map[string]string `json:"additionalTrustedCAs,omitempty"`

	// AllowedRegistriesForImport limits the container image registries that normal users may import
	// images from. Set this list to the registries that you trust to contain valid Docker
	// images and that you want applications to be able to import from.
	// +optional
	AllowedRegistriesForImport []RegistryLocation `json:"allowedRegistriesForImport,omitempty"`

	// RegistrySources contains configuration that determines how the container runtime
	// should treat individual registries when accessing images. It does not contain configuration
	// for the internal cluster registry. AllowedRegistries, BlockedRegistries are mutually exclusive.
	// +optional
	RegistrySources *RegistrySources `json:"registrySources,omitempty"`
}

// RegistryLocation contains a location of the registry specified by the registry domain name.
type RegistryLocation struct {
	// domainName specifies a domain name for the registry. The domain name might include wildcards, like '*' or '??'.
	// In case the registry use non-standard (80 or 443) port, the port should be included in the domain name as well.
	// +optional
	DomainName string `json:"domainName,omitempty"`

	// insecure indicates whether the registry is secure (https) or insecure (http), default is secured.
	// +kubebuilder:default=false
	// +optional
	Insecure bool `json:"insecure,omitempty"`
}

// RegistrySources contains registries configuration.
type RegistrySources struct {
	// AllowedRegistries are the registries for which image pull and push actions are allowed.
	// To specify all subdomains, add the asterisk (*) wildcard character as a prefix to the domain name,
	// For example, *.example.com.
	// You can specify an individual repository within a registry, For example: reg1.io/myrepo/myapp:latest.
	// All other registries are blocked.
	// +optional
	AllowedRegistries []string `json:"allowedRegistries,omitempty"`

	// BlockedRegistries are the registries for which image pull and push actions are denied.
	// To specify all subdomains, add the asterisk (*) wildcard character as a prefix to the domain name,
	// For example, *.example.com.
	// You can specify an individual repository within a registry, For example: reg1.io/myrepo/myapp:latest.
	// All other registries are allowed.
	// +optional
	BlockedRegistries []string `json:"blockedRegistries,omitempty"`

	// InsecureRegistries are registries which do not have a valid TLS certificate or only support HTTP connections.
	// To specify all subdomains, add the asterisk (*) wildcard character as a prefix to the domain name,
	// For example, *.example.com.
	// You can specify an individual repository within a registry, For example: reg1.io/myrepo/myapp:latest.
	// +optional
	InsecureRegistries []string `json:"insecureRegistries,omitempty"`
}

// NetworkSpec for ROSA-HCP.
type NetworkSpec struct {
	// IP addresses block used by OpenShift while installing the cluster, for example "10.0.0.0/16".
	// +kubebuilder:validation:Format=cidr
	// +optional
	MachineCIDR string `json:"machineCIDR,omitempty"`

	// IP address block from which to assign pod IP addresses, for example `10.128.0.0/14`.
	// +kubebuilder:validation:Format=cidr
	// +optional
	PodCIDR string `json:"podCIDR,omitempty"`

	// IP address block from which to assign service IP addresses, for example `172.30.0.0/16`.
	// +kubebuilder:validation:Format=cidr
	// +optional
	ServiceCIDR string `json:"serviceCIDR,omitempty"`

	// Network host prefix which is defaulted to `23` if not specified.
	// +kubebuilder:default=23
	// +optional
	HostPrefix int `json:"hostPrefix,omitempty"`

	// The CNI network type default is OVNKubernetes.
	// +kubebuilder:validation:Enum=OVNKubernetes;Other
	// +kubebuilder:default=OVNKubernetes
	// +optional
	NetworkType string `json:"networkType,omitempty"`
}

// DefaultMachinePoolSpec defines the configuration for the required worker nodes provisioned as part of the cluster creation.
type DefaultMachinePoolSpec struct {
	// The instance type to use, for example `r5.xlarge`. Instance type ref; https://aws.amazon.com/ec2/instance-types/
	// +optional
	InstanceType string `json:"instanceType,omitempty"`

	// Autoscaling specifies auto scaling behaviour for the default MachinePool. Autoscaling min/max value
	// must be equal or multiple of the availability zones count.
	// +optional
	Autoscaling *AutoScaling `json:"autoscaling,omitempty"`

	// VolumeSize set the disk volume size for the default workers machine pool in Gib. The default is 300 GiB.
	// +kubebuilder:validation:Minimum=75
	// +kubebuilder:validation:Maximum=16384
	// +immutable
	// +optional
	VolumeSize int `json:"volumeSize,omitempty"`
}

// AutoScaling specifies scaling options.
type AutoScaling struct {
	// +kubebuilder:validation:Minimum=1
	MinReplicas int `json:"minReplicas,omitempty"`
	// +kubebuilder:validation:Minimum=1
	MaxReplicas int `json:"maxReplicas,omitempty"`
}

// AWSRolesRef contains references to various AWS IAM roles required for operators to make calls against the AWS API.
type AWSRolesRef struct {
	// The referenced role must have a trust relationship that allows it to be assumed via web identity.
	// https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_providers_oidc.html.
	// Example:
	// {
	//		"Version": "2012-10-17",
	//		"Statement": [
	//			{
	//				"Effect": "Allow",
	//				"Principal": {
	//					"Federated": "{{ .ProviderARN }}"
	//				},
	//					"Action": "sts:AssumeRoleWithWebIdentity",
	//				"Condition": {
	//					"StringEquals": {
	//						"{{ .ProviderName }}:sub": {{ .ServiceAccounts }}
	//					}
	//				}
	//			}
	//		]
	//	}
	//
	// IngressARN is an ARN value referencing a role appropriate for the Ingress Operator.
	//
	// The following is an example of a valid policy document:
	//
	// {
	//	"Version": "2012-10-17",
	//	"Statement": [
	//		{
	//			"Effect": "Allow",
	//			"Action": [
	//				"elasticloadbalancing:DescribeLoadBalancers",
	//				"tag:GetResources",
	//				"route53:ListHostedZones"
	//			],
	//			"Resource": "*"
	//		},
	//		{
	//			"Effect": "Allow",
	//			"Action": [
	//				"route53:ChangeResourceRecordSets"
	//			],
	//			"Resource": [
	//				"arn:aws:route53:::PUBLIC_ZONE_ID",
	//				"arn:aws:route53:::PRIVATE_ZONE_ID"
	//			]
	//		}
	//	]
	// }
	IngressARN string `json:"ingressARN"`

	// ImageRegistryARN is an ARN value referencing a role appropriate for the Image Registry Operator.
	//
	// The following is an example of a valid policy document:
	//
	// {
	//	"Version": "2012-10-17",
	//	"Statement": [
	//		{
	//			"Effect": "Allow",
	//			"Action": [
	//				"s3:CreateBucket",
	//				"s3:DeleteBucket",
	//				"s3:PutBucketTagging",
	//				"s3:GetBucketTagging",
	//				"s3:PutBucketPublicAccessBlock",
	//				"s3:GetBucketPublicAccessBlock",
	//				"s3:PutEncryptionConfiguration",
	//				"s3:GetEncryptionConfiguration",
	//				"s3:PutLifecycleConfiguration",
	//				"s3:GetLifecycleConfiguration",
	//				"s3:GetBucketLocation",
	//				"s3:ListBucket",
	//				"s3:GetObject",
	//				"s3:PutObject",
	//				"s3:DeleteObject",
	//				"s3:ListBucketMultipartUploads",
	//				"s3:AbortMultipartUpload",
	//				"s3:ListMultipartUploadParts"
	//			],
	//			"Resource": "*"
	//		}
	//	]
	// }
	ImageRegistryARN string `json:"imageRegistryARN"`

	// StorageARN is an ARN value referencing a role appropriate for the Storage Operator.
	//
	// The following is an example of a valid policy document:
	//
	// {
	//	"Version": "2012-10-17",
	//	"Statement": [
	//		{
	//			"Effect": "Allow",
	//			"Action": [
	//				"ec2:AttachVolume",
	//				"ec2:CreateSnapshot",
	//				"ec2:CreateTags",
	//				"ec2:CreateVolume",
	//				"ec2:DeleteSnapshot",
	//				"ec2:DeleteTags",
	//				"ec2:DeleteVolume",
	//				"ec2:DescribeInstances",
	//				"ec2:DescribeSnapshots",
	//				"ec2:DescribeTags",
	//				"ec2:DescribeVolumes",
	//				"ec2:DescribeVolumesModifications",
	//				"ec2:DetachVolume",
	//				"ec2:ModifyVolume"
	//			],
	//			"Resource": "*"
	//		}
	//	]
	// }
	StorageARN string `json:"storageARN"`

	// NetworkARN is an ARN value referencing a role appropriate for the Network Operator.
	//
	// The following is an example of a valid policy document:
	//
	// {
	//	"Version": "2012-10-17",
	//	"Statement": [
	//		{
	//			"Effect": "Allow",
	//			"Action": [
	//				"ec2:DescribeInstances",
	//        "ec2:DescribeInstanceStatus",
	//        "ec2:DescribeInstanceTypes",
	//        "ec2:UnassignPrivateIpAddresses",
	//        "ec2:AssignPrivateIpAddresses",
	//        "ec2:UnassignIpv6Addresses",
	//        "ec2:AssignIpv6Addresses",
	//        "ec2:DescribeSubnets",
	//        "ec2:DescribeNetworkInterfaces"
	//			],
	//			"Resource": "*"
	//		}
	//	]
	// }
	NetworkARN string `json:"networkARN"`

	// KubeCloudControllerARN is an ARN value referencing a role appropriate for the KCM/KCC.
	// Source: https://cloud-provider-aws.sigs.k8s.io/prerequisites/#iam-policies
	//
	// The following is an example of a valid policy document:
	//
	//  {
	//  "Version": "2012-10-17",
	//  "Statement": [
	//    {
	//      "Action": [
	//        "autoscaling:DescribeAutoScalingGroups",
	//        "autoscaling:DescribeLaunchConfigurations",
	//        "autoscaling:DescribeTags",
	//        "ec2:DescribeAvailabilityZones",
	//        "ec2:DescribeInstances",
	//        "ec2:DescribeImages",
	//        "ec2:DescribeRegions",
	//        "ec2:DescribeRouteTables",
	//        "ec2:DescribeSecurityGroups",
	//        "ec2:DescribeSubnets",
	//        "ec2:DescribeVolumes",
	//        "ec2:CreateSecurityGroup",
	//        "ec2:CreateTags",
	//        "ec2:CreateVolume",
	//        "ec2:ModifyInstanceAttribute",
	//        "ec2:ModifyVolume",
	//        "ec2:AttachVolume",
	//        "ec2:AuthorizeSecurityGroupIngress",
	//        "ec2:CreateRoute",
	//        "ec2:DeleteRoute",
	//        "ec2:DeleteSecurityGroup",
	//        "ec2:DeleteVolume",
	//        "ec2:DetachVolume",
	//        "ec2:RevokeSecurityGroupIngress",
	//        "ec2:DescribeVpcs",
	//        "elasticloadbalancing:AddTags",
	//        "elasticloadbalancing:AttachLoadBalancerToSubnets",
	//        "elasticloadbalancing:ApplySecurityGroupsToLoadBalancer",
	//        "elasticloadbalancing:CreateLoadBalancer",
	//        "elasticloadbalancing:CreateLoadBalancerPolicy",
	//        "elasticloadbalancing:CreateLoadBalancerListeners",
	//        "elasticloadbalancing:ConfigureHealthCheck",
	//        "elasticloadbalancing:DeleteLoadBalancer",
	//        "elasticloadbalancing:DeleteLoadBalancerListeners",
	//        "elasticloadbalancing:DescribeLoadBalancers",
	//        "elasticloadbalancing:DescribeLoadBalancerAttributes",
	//        "elasticloadbalancing:DetachLoadBalancerFromSubnets",
	//        "elasticloadbalancing:DeregisterInstancesFromLoadBalancer",
	//        "elasticloadbalancing:ModifyLoadBalancerAttributes",
	//        "elasticloadbalancing:RegisterInstancesWithLoadBalancer",
	//        "elasticloadbalancing:SetLoadBalancerPoliciesForBackendServer",
	//        "elasticloadbalancing:AddTags",
	//        "elasticloadbalancing:CreateListener",
	//        "elasticloadbalancing:CreateTargetGroup",
	//        "elasticloadbalancing:DeleteListener",
	//        "elasticloadbalancing:DeleteTargetGroup",
	//        "elasticloadbalancing:DeregisterTargets",
	//        "elasticloadbalancing:DescribeListeners",
	//        "elasticloadbalancing:DescribeLoadBalancerPolicies",
	//        "elasticloadbalancing:DescribeTargetGroups",
	//        "elasticloadbalancing:DescribeTargetHealth",
	//        "elasticloadbalancing:ModifyListener",
	//        "elasticloadbalancing:ModifyTargetGroup",
	//        "elasticloadbalancing:RegisterTargets",
	//        "elasticloadbalancing:SetLoadBalancerPoliciesOfListener",
	//        "iam:CreateServiceLinkedRole",
	//        "kms:DescribeKey"
	//      ],
	//      "Resource": [
	//        "*"
	//      ],
	//      "Effect": "Allow"
	//    }
	//  ]
	// }
	// +immutable
	KubeCloudControllerARN string `json:"kubeCloudControllerARN"`

	// NodePoolManagementARN is an ARN value referencing a role appropriate for the CAPI Controller.
	//
	// The following is an example of a valid policy document:
	//
	// {
	//   "Version": "2012-10-17",
	//  "Statement": [
	//    {
	//      "Action": [
	//        "ec2:AssociateRouteTable",
	//        "ec2:AttachInternetGateway",
	//        "ec2:AuthorizeSecurityGroupIngress",
	//        "ec2:CreateInternetGateway",
	//        "ec2:CreateNatGateway",
	//        "ec2:CreateRoute",
	//        "ec2:CreateRouteTable",
	//        "ec2:CreateSecurityGroup",
	//        "ec2:CreateSubnet",
	//        "ec2:CreateTags",
	//        "ec2:DeleteInternetGateway",
	//        "ec2:DeleteNatGateway",
	//        "ec2:DeleteRouteTable",
	//        "ec2:DeleteSecurityGroup",
	//        "ec2:DeleteSubnet",
	//        "ec2:DeleteTags",
	//        "ec2:DescribeAccountAttributes",
	//        "ec2:DescribeAddresses",
	//        "ec2:DescribeAvailabilityZones",
	//        "ec2:DescribeImages",
	//        "ec2:DescribeInstances",
	//        "ec2:DescribeInternetGateways",
	//        "ec2:DescribeNatGateways",
	//        "ec2:DescribeNetworkInterfaces",
	//        "ec2:DescribeNetworkInterfaceAttribute",
	//        "ec2:DescribeRouteTables",
	//        "ec2:DescribeSecurityGroups",
	//        "ec2:DescribeSubnets",
	//        "ec2:DescribeVpcs",
	//        "ec2:DescribeVpcAttribute",
	//        "ec2:DescribeVolumes",
	//        "ec2:DetachInternetGateway",
	//        "ec2:DisassociateRouteTable",
	//        "ec2:DisassociateAddress",
	//        "ec2:ModifyInstanceAttribute",
	//        "ec2:ModifyNetworkInterfaceAttribute",
	//        "ec2:ModifySubnetAttribute",
	//        "ec2:RevokeSecurityGroupIngress",
	//        "ec2:RunInstances",
	//        "ec2:TerminateInstances",
	//        "tag:GetResources",
	//        "ec2:CreateLaunchTemplate",
	//        "ec2:CreateLaunchTemplateVersion",
	//        "ec2:DescribeLaunchTemplates",
	//        "ec2:DescribeLaunchTemplateVersions",
	//        "ec2:DeleteLaunchTemplate",
	//        "ec2:DeleteLaunchTemplateVersions"
	//      ],
	//      "Resource": [
	//        "*"
	//      ],
	//      "Effect": "Allow"
	//    },
	//    {
	//      "Condition": {
	//        "StringLike": {
	//          "iam:AWSServiceName": "elasticloadbalancing.amazonaws.com"
	//        }
	//      },
	//      "Action": [
	//        "iam:CreateServiceLinkedRole"
	//      ],
	//      "Resource": [
	//        "arn:*:iam::*:role/aws-service-role/elasticloadbalancing.amazonaws.com/AWSServiceRoleForElasticLoadBalancing"
	//      ],
	//      "Effect": "Allow"
	//    },
	//    {
	//      "Action": [
	//        "iam:PassRole"
	//      ],
	//      "Resource": [
	//        "arn:*:iam::*:role/*-worker-role"
	//      ],
	//      "Effect": "Allow"
	//    },
	// 	  {
	// 	  	"Effect": "Allow",
	// 	  	"Action": [
	// 	  		"kms:Decrypt",
	// 	  		"kms:ReEncrypt",
	// 	  		"kms:GenerateDataKeyWithoutPlainText",
	// 	  		"kms:DescribeKey"
	// 	  	],
	// 	  	"Resource": "*"
	// 	  },
	// 	  {
	// 	  	"Effect": "Allow",
	// 	  	"Action": [
	// 	  		"kms:CreateGrant"
	// 	  	],
	// 	  	"Resource": "*",
	// 	  	"Condition": {
	// 	  		"Bool": {
	// 	  			"kms:GrantIsForAWSResource": true
	// 	  		}
	// 	  	}
	// 	  }
	//  ]
	// }
	//
	// +immutable
	NodePoolManagementARN string `json:"nodePoolManagementARN"`

	// ControlPlaneOperatorARN  is an ARN value referencing a role appropriate for the Control Plane Operator.
	//
	// The following is an example of a valid policy document:
	//
	// {
	//	"Version": "2012-10-17",
	//	"Statement": [
	//		{
	//			"Effect": "Allow",
	//			"Action": [
	//				"ec2:CreateVpcEndpoint",
	//				"ec2:DescribeVpcEndpoints",
	//				"ec2:ModifyVpcEndpoint",
	//				"ec2:DeleteVpcEndpoints",
	//				"ec2:CreateTags",
	//				"route53:ListHostedZones",
	//				"ec2:CreateSecurityGroup",
	//				"ec2:AuthorizeSecurityGroupIngress",
	//				"ec2:AuthorizeSecurityGroupEgress",
	//				"ec2:DeleteSecurityGroup",
	//				"ec2:RevokeSecurityGroupIngress",
	//				"ec2:RevokeSecurityGroupEgress",
	//				"ec2:DescribeSecurityGroups",
	//				"ec2:DescribeVpcs",
	//			],
	//			"Resource": "*"
	//		},
	//		{
	//			"Effect": "Allow",
	//			"Action": [
	//				"route53:ChangeResourceRecordSets",
	//				"route53:ListResourceRecordSets"
	//			],
	//			"Resource": "arn:aws:route53:::%s"
	//		}
	//	]
	// }
	// +immutable
	ControlPlaneOperatorARN string `json:"controlPlaneOperatorARN"`
	KMSProviderARN          string `json:"kmsProviderARN"`
}

// RosaControlPlaneStatus defines the observed state of ROSAControlPlane.
type RosaControlPlaneStatus struct {
	// ExternalManagedControlPlane indicates to cluster-api that the control plane
	// is managed by an external service such as AKS, EKS, GKE, etc.
	// +kubebuilder:default=true
	ExternalManagedControlPlane *bool `json:"externalManagedControlPlane,omitempty"`
	// Initialized denotes whether or not the control plane has the
	// uploaded kubernetes config-map.
	// +optional
	Initialized bool `json:"initialized"`
	// Ready denotes that the ROSAControlPlane API Server is ready to receive requests.
	// +kubebuilder:default=false
	Ready bool `json:"ready"`
	// FailureMessage will be set in the event that there is a terminal problem
	// reconciling the state and will be set to a descriptive error message.
	//
	// This field should not be set for transitive errors that a controller
	// faces that are expected to be fixed automatically over
	// time (like service outages), but instead indicate that something is
	// fundamentally wrong with the spec or the configuration of
	// the controller, and that manual intervention is required.
	//
	// +optional
	FailureMessage *string `json:"failureMessage,omitempty"`
	// Conditions specifies the conditions for the managed control plane
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`

	// ID is the cluster ID given by ROSA.
	ID string `json:"id,omitempty"`
	// ConsoleURL is the url for the openshift console.
	ConsoleURL string `json:"consoleURL,omitempty"`
	// OIDCEndpointURL is the endpoint url for the managed OIDC provider.
	OIDCEndpointURL string `json:"oidcEndpointURL,omitempty"`

	// OpenShift semantic version, for example "4.14.5".
	// +optional
	Version string `json:"version"`

	// Available upgrades for the ROSA hosted control plane.
	AvailableUpgrades []string `json:"availableUpgrades,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=rosacontrolplanes,shortName=rosacp,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this RosaControl belongs"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Control plane infrastructure is ready for worker nodes"
// +k8s:defaulter-gen=true

// ROSAControlPlane is the Schema for the ROSAControlPlanes API.
type ROSAControlPlane struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RosaControlPlaneSpec   `json:"spec,omitempty"`
	Status RosaControlPlaneStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ROSAControlPlaneList contains a list of ROSAControlPlane.
type ROSAControlPlaneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ROSAControlPlane `json:"items"`
}

// GetConditions returns the control planes conditions.
func (r *ROSAControlPlane) GetConditions() clusterv1.Conditions {
	return r.Status.Conditions
}

// SetConditions sets the status conditions for the AWSManagedControlPlane.
func (r *ROSAControlPlane) SetConditions(conditions clusterv1.Conditions) {
	r.Status.Conditions = conditions
}

func init() {
	SchemeBuilder.Register(&ROSAControlPlane{}, &ROSAControlPlaneList{})
}
