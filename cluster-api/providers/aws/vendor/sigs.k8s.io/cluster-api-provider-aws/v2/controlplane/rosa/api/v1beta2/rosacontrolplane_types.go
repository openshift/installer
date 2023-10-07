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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

type RosaControlPlaneSpec struct { //nolint: maligned
	// The Subnet IDs to use when installing the cluster.
	// SubnetIDs should come in pairs; two per availability zone, one private and one public.
	Subnets []string `json:"subnets"`

	// AWS AvailabilityZones of the worker nodes
	// should match the AvailabilityZones of the Subnets.
	AvailabilityZones []string `json:"availabilityZones"`

	// Block of IP addresses used by OpenShift while installing the cluster, for example "10.0.0.0/16".
	MachineCIDR *string `json:"machineCIDR"`

	// The AWS Region the cluster lives in.
	Region *string `json:"region"`

	// Openshift version, for example "openshift-v4.12.15".
	Version *string `json:"version"`

	// ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// +optional
	ControlPlaneEndpoint clusterv1.APIEndpoint `json:"controlPlaneEndpoint"`

	// AWS IAM roles used to perform credential requests by the openshift operators.
	RolesRef AWSRolesRef `json:"rolesRef"`

	// The ID of the OpenID Connect Provider.
	OIDCID *string `json:"oidcID"`

	// TODO: these are to satisfy ocm sdk. Explore how to drop them.
	AccountID        *string `json:"accountID"`
	CreatorARN       *string `json:"creatorARN"`
	InstallerRoleARN *string `json:"installerRoleARN"`
	SupportRoleARN   *string `json:"supportRoleARN"`
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

type RosaControlPlaneStatus struct {
	// ExternalManagedControlPlane indicates to cluster-api that the control plane
	// is managed by an external service such as AKS, EKS, GKE, etc.
	// +kubebuilder:default=true
	ExternalManagedControlPlane *bool `json:"externalManagedControlPlane,omitempty"`
	// Initialized denotes whether or not the control plane has the
	// uploaded kubernetes config-map.
	// +optional
	Initialized bool `json:"initialized"`
	// Ready denotes that the AWSManagedControlPlane API Server is ready to
	// receive requests and that the VPC infra is ready.
	// +kubebuilder:default=false
	Ready bool `json:"ready"`
	// ErrorMessage indicates that there is a terminal problem reconciling the
	// state, and will be set to a descriptive error message.
	// +optional
	FailureMessage *string `json:"failureMessage,omitempty"`
	// Conditions specifies the cpnditions for the managed control plane
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=rosacontrolplanes,shortName=rosacp,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this RosaControl belongs"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Control plane infrastructure is ready for worker nodes"

type ROSAControlPlane struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RosaControlPlaneSpec   `json:"spec,omitempty"`
	Status RosaControlPlaneStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

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
