// Copyright 2023 Google LLC. All Rights Reserved.
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
//     http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package containeraws

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func DCLClusterSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "ContainerAws/Cluster",
			Description: "An Anthos cluster running on AWS.",
			StructName:  "Cluster",
			Reference: &dcl.Link{
				Text: "API reference",
				URL:  "https://cloud.google.com/anthos/clusters/docs/multi-cloud/reference/rest/v1/projects.locations.awsClusters",
			},
			Guides: []*dcl.Link{
				&dcl.Link{
					Text: "Multicloud overview",
					URL:  "https://cloud.google.com/anthos/clusters/docs/multi-cloud",
				},
			},
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a Cluster",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "cluster",
						Required:    true,
						Description: "A full instance of a Cluster",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a Cluster",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "cluster",
						Required:    true,
						Description: "A full instance of a Cluster",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a Cluster",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "cluster",
						Required:    true,
						Description: "A full instance of a Cluster",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all Cluster",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:     "project",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
					dcl.PathParameters{
						Name:     "location",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
				},
			},
			List: &dcl.Path{
				Description: "The function used to list information about many Cluster",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:     "project",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
					dcl.PathParameters{
						Name:     "location",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
				},
			},
		},
		Components: &dcl.Components{
			Schemas: map[string]*dcl.Component{
				"Cluster": &dcl.Component{
					Title:           "Cluster",
					ID:              "projects/{{project}}/locations/{{location}}/awsClusters/{{name}}",
					ParentContainer: "project",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"name",
							"networking",
							"awsRegion",
							"controlPlane",
							"authorization",
							"project",
							"location",
							"fleet",
						},
						Properties: map[string]*dcl.Property{
							"annotations": &dcl.Property{
								Type: "object",
								AdditionalProperties: &dcl.Property{
									Type: "string",
								},
								GoName:      "Annotations",
								Description: "Optional. Annotations on the cluster. This field has the same restrictions as Kubernetes annotations. The total size of all keys and values combined is limited to 256k. Key can have 2 segments: prefix (optional) and name (required), separated by a slash (/). Prefix must be a DNS subdomain. Name must be 63 characters or less, begin and end with alphanumerics, with dashes (-), underscores (_), dots (.), and alphanumerics between.",
							},
							"authorization": &dcl.Property{
								Type:        "object",
								GoName:      "Authorization",
								GoType:      "ClusterAuthorization",
								Description: "Configuration related to the cluster RBAC settings.",
								Required: []string{
									"adminUsers",
								},
								Properties: map[string]*dcl.Property{
									"adminUsers": &dcl.Property{
										Type:        "array",
										GoName:      "AdminUsers",
										Description: "Users to perform operations as a cluster admin. A managed ClusterRoleBinding will be created to grant the `cluster-admin` ClusterRole to the users. Up to ten admin users can be provided. For more info on RBAC, see https://kubernetes.io/docs/reference/access-authn-authz/rbac/#user-facing-roles",
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "object",
											GoType: "ClusterAuthorizationAdminUsers",
											Required: []string{
												"username",
											},
											Properties: map[string]*dcl.Property{
												"username": &dcl.Property{
													Type:        "string",
													GoName:      "Username",
													Description: "The name of the user, e.g. `my-gcp-id@gmail.com`.",
												},
											},
										},
									},
								},
							},
							"awsRegion": &dcl.Property{
								Type:        "string",
								GoName:      "AwsRegion",
								Description: "The AWS region where the cluster runs. Each Google Cloud region supports a subset of nearby AWS regions. You can call to list all supported AWS regions within a given Google Cloud region.",
								Immutable:   true,
							},
							"controlPlane": &dcl.Property{
								Type:        "object",
								GoName:      "ControlPlane",
								GoType:      "ClusterControlPlane",
								Description: "Configuration related to the cluster control plane.",
								Required: []string{
									"version",
									"subnetIds",
									"configEncryption",
									"iamInstanceProfile",
									"databaseEncryption",
									"awsServicesAuthentication",
								},
								Properties: map[string]*dcl.Property{
									"awsServicesAuthentication": &dcl.Property{
										Type:        "object",
										GoName:      "AwsServicesAuthentication",
										GoType:      "ClusterControlPlaneAwsServicesAuthentication",
										Description: "Authentication configuration for management of AWS resources.",
										Required: []string{
											"roleArn",
										},
										Properties: map[string]*dcl.Property{
											"roleArn": &dcl.Property{
												Type:        "string",
												GoName:      "RoleArn",
												Description: "The Amazon Resource Name (ARN) of the role that the Anthos Multi-Cloud API will assume when managing AWS resources on your account.",
											},
											"roleSessionName": &dcl.Property{
												Type:          "string",
												GoName:        "RoleSessionName",
												Description:   "Optional. An identifier for the assumed role session. When unspecified, it defaults to `multicloud-service-agent`.",
												ServerDefault: true,
											},
										},
									},
									"configEncryption": &dcl.Property{
										Type:        "object",
										GoName:      "ConfigEncryption",
										GoType:      "ClusterControlPlaneConfigEncryption",
										Description: "The ARN of the AWS KMS key used to encrypt cluster configuration.",
										Required: []string{
											"kmsKeyArn",
										},
										Properties: map[string]*dcl.Property{
											"kmsKeyArn": &dcl.Property{
												Type:        "string",
												GoName:      "KmsKeyArn",
												Description: "The ARN of the AWS KMS key used to encrypt cluster configuration.",
											},
										},
									},
									"databaseEncryption": &dcl.Property{
										Type:        "object",
										GoName:      "DatabaseEncryption",
										GoType:      "ClusterControlPlaneDatabaseEncryption",
										Description: "The ARN of the AWS KMS key used to encrypt cluster secrets.",
										Immutable:   true,
										Required: []string{
											"kmsKeyArn",
										},
										Properties: map[string]*dcl.Property{
											"kmsKeyArn": &dcl.Property{
												Type:        "string",
												GoName:      "KmsKeyArn",
												Description: "The ARN of the AWS KMS key used to encrypt cluster secrets.",
												Immutable:   true,
											},
										},
									},
									"iamInstanceProfile": &dcl.Property{
										Type:        "string",
										GoName:      "IamInstanceProfile",
										Description: "The name of the AWS IAM instance pofile to assign to each control plane replica.",
									},
									"instanceType": &dcl.Property{
										Type:          "string",
										GoName:        "InstanceType",
										Description:   "Optional. The AWS instance type. When unspecified, it defaults to `m5.large`.",
										ServerDefault: true,
									},
									"mainVolume": &dcl.Property{
										Type:          "object",
										GoName:        "MainVolume",
										GoType:        "ClusterControlPlaneMainVolume",
										Description:   "Optional. Configuration related to the main volume provisioned for each control plane replica. The main volume is in charge of storing all of the cluster's etcd state. Volumes will be provisioned in the availability zone associated with the corresponding subnet. When unspecified, it defaults to 8 GiB with the GP2 volume type.",
										Immutable:     true,
										ServerDefault: true,
										Properties: map[string]*dcl.Property{
											"iops": &dcl.Property{
												Type:          "integer",
												Format:        "int64",
												GoName:        "Iops",
												Description:   "Optional. The number of I/O operations per second (IOPS) to provision for GP3 volume.",
												Immutable:     true,
												ServerDefault: true,
											},
											"kmsKeyArn": &dcl.Property{
												Type:        "string",
												GoName:      "KmsKeyArn",
												Description: "Optional. The Amazon Resource Name (ARN) of the Customer Managed Key (CMK) used to encrypt AWS EBS volumes. If not specified, the default Amazon managed key associated to the AWS region where this cluster runs will be used.",
												Immutable:   true,
											},
											"sizeGib": &dcl.Property{
												Type:          "integer",
												Format:        "int64",
												GoName:        "SizeGib",
												Description:   "Optional. The size of the volume, in GiBs. When unspecified, a default value is provided. See the specific reference in the parent resource.",
												Immutable:     true,
												ServerDefault: true,
											},
											"throughput": &dcl.Property{
												Type:          "integer",
												Format:        "int64",
												GoName:        "Throughput",
												Description:   "Optional. The throughput to provision for the volume, in MiB/s. Only valid if the volume type is GP3.",
												Immutable:     true,
												ServerDefault: true,
											},
											"volumeType": &dcl.Property{
												Type:          "string",
												GoName:        "VolumeType",
												GoType:        "ClusterControlPlaneMainVolumeVolumeTypeEnum",
												Description:   "Optional. Type of the EBS volume. When unspecified, it defaults to GP2 volume. Possible values: VOLUME_TYPE_UNSPECIFIED, GP2, GP3",
												Immutable:     true,
												ServerDefault: true,
												Enum: []string{
													"VOLUME_TYPE_UNSPECIFIED",
													"GP2",
													"GP3",
												},
											},
										},
									},
									"proxyConfig": &dcl.Property{
										Type:        "object",
										GoName:      "ProxyConfig",
										GoType:      "ClusterControlPlaneProxyConfig",
										Description: "Proxy configuration for outbound HTTP(S) traffic.",
										Required: []string{
											"secretArn",
											"secretVersion",
										},
										Properties: map[string]*dcl.Property{
											"secretArn": &dcl.Property{
												Type:        "string",
												GoName:      "SecretArn",
												Description: "The ARN of the AWS Secret Manager secret that contains the HTTP(S) proxy configuration.",
											},
											"secretVersion": &dcl.Property{
												Type:        "string",
												GoName:      "SecretVersion",
												Description: "The version string of the AWS Secret Manager secret that contains the HTTP(S) proxy configuration.",
											},
										},
									},
									"rootVolume": &dcl.Property{
										Type:          "object",
										GoName:        "RootVolume",
										GoType:        "ClusterControlPlaneRootVolume",
										Description:   "Optional. Configuration related to the root volume provisioned for each control plane replica. Volumes will be provisioned in the availability zone associated with the corresponding subnet. When unspecified, it defaults to 32 GiB with the GP2 volume type.",
										ServerDefault: true,
										Properties: map[string]*dcl.Property{
											"iops": &dcl.Property{
												Type:          "integer",
												Format:        "int64",
												GoName:        "Iops",
												Description:   "Optional. The number of I/O operations per second (IOPS) to provision for GP3 volume.",
												ServerDefault: true,
											},
											"kmsKeyArn": &dcl.Property{
												Type:        "string",
												GoName:      "KmsKeyArn",
												Description: "Optional. The Amazon Resource Name (ARN) of the Customer Managed Key (CMK) used to encrypt AWS EBS volumes. If not specified, the default Amazon managed key associated to the AWS region where this cluster runs will be used.",
											},
											"sizeGib": &dcl.Property{
												Type:          "integer",
												Format:        "int64",
												GoName:        "SizeGib",
												Description:   "Optional. The size of the volume, in GiBs. When unspecified, a default value is provided. See the specific reference in the parent resource.",
												ServerDefault: true,
											},
											"throughput": &dcl.Property{
												Type:          "integer",
												Format:        "int64",
												GoName:        "Throughput",
												Description:   "Optional. The throughput to provision for the volume, in MiB/s. Only valid if the volume type is GP3.",
												ServerDefault: true,
											},
											"volumeType": &dcl.Property{
												Type:          "string",
												GoName:        "VolumeType",
												GoType:        "ClusterControlPlaneRootVolumeVolumeTypeEnum",
												Description:   "Optional. Type of the EBS volume. When unspecified, it defaults to GP2 volume. Possible values: VOLUME_TYPE_UNSPECIFIED, GP2, GP3",
												ServerDefault: true,
												Enum: []string{
													"VOLUME_TYPE_UNSPECIFIED",
													"GP2",
													"GP3",
												},
											},
										},
									},
									"securityGroupIds": &dcl.Property{
										Type:        "array",
										GoName:      "SecurityGroupIds",
										Description: "Optional. The IDs of additional security groups to add to control plane replicas. The Anthos Multi-Cloud API will automatically create and manage security groups with the minimum rules needed for a functioning cluster.",
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "string",
											GoType: "string",
										},
									},
									"sshConfig": &dcl.Property{
										Type:        "object",
										GoName:      "SshConfig",
										GoType:      "ClusterControlPlaneSshConfig",
										Description: "Optional. SSH configuration for how to access the underlying control plane machines.",
										Required: []string{
											"ec2KeyPair",
										},
										Properties: map[string]*dcl.Property{
											"ec2KeyPair": &dcl.Property{
												Type:        "string",
												GoName:      "Ec2KeyPair",
												Description: "The name of the EC2 key pair used to login into cluster machines.",
											},
										},
									},
									"subnetIds": &dcl.Property{
										Type:        "array",
										GoName:      "SubnetIds",
										Description: "The list of subnets where control plane replicas will run. A replica will be provisioned on each subnet and up to three values can be provided. Each subnet must be in a different AWS Availability Zone (AZ).",
										Immutable:   true,
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "string",
											GoType: "string",
										},
									},
									"tags": &dcl.Property{
										Type: "object",
										AdditionalProperties: &dcl.Property{
											Type: "string",
										},
										GoName:      "Tags",
										Description: "Optional. A set of AWS resource tags to propagate to all underlying managed AWS resources. Specify at most 50 pairs containing alphanumerics, spaces, and symbols (.+-=_:@/). Keys can be up to 127 Unicode characters. Values can be up to 255 Unicode characters.",
									},
									"version": &dcl.Property{
										Type:        "string",
										GoName:      "Version",
										Description: "The Kubernetes version to run on control plane replicas (e.g. `1.19.10-gke.1000`). You can list all supported versions on a given Google Cloud region by calling .",
									},
								},
							},
							"createTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "CreateTime",
								ReadOnly:    true,
								Description: "Output only. The time at which this cluster was created.",
								Immutable:   true,
							},
							"description": &dcl.Property{
								Type:        "string",
								GoName:      "Description",
								Description: "Optional. A human readable description of this cluster. Cannot be longer than 255 UTF-8 encoded bytes.",
							},
							"endpoint": &dcl.Property{
								Type:        "string",
								GoName:      "Endpoint",
								ReadOnly:    true,
								Description: "Output only. The endpoint of the cluster's API server.",
								Immutable:   true,
							},
							"etag": &dcl.Property{
								Type:        "string",
								GoName:      "Etag",
								ReadOnly:    true,
								Description: "Allows clients to perform consistent read-modify-writes through optimistic concurrency control. May be sent on update and delete requests to ensure the client has an up-to-date value before proceeding.",
								Immutable:   true,
							},
							"fleet": &dcl.Property{
								Type:        "object",
								GoName:      "Fleet",
								GoType:      "ClusterFleet",
								Description: "Fleet configuration.",
								Immutable:   true,
								Required: []string{
									"project",
								},
								Properties: map[string]*dcl.Property{
									"membership": &dcl.Property{
										Type:        "string",
										GoName:      "Membership",
										ReadOnly:    true,
										Description: "The name of the managed Hub Membership resource associated to this cluster. Membership names are formatted as projects/<project-number>/locations/global/membership/<cluster-id>.",
										Immutable:   true,
									},
									"project": &dcl.Property{
										Type:        "string",
										GoName:      "Project",
										Description: "The number of the Fleet host project where this cluster will be registered.",
										Immutable:   true,
										ResourceReferences: []*dcl.PropertyResourceReference{
											&dcl.PropertyResourceReference{
												Resource: "Cloudresourcemanager/Project",
												Field:    "name",
												Parent:   true,
											},
										},
									},
								},
							},
							"location": &dcl.Property{
								Type:        "string",
								GoName:      "Location",
								Description: "The location for the resource",
								Immutable:   true,
							},
							"name": &dcl.Property{
								Type:        "string",
								GoName:      "Name",
								Description: "The name of this resource.",
								Immutable:   true,
							},
							"networking": &dcl.Property{
								Type:        "object",
								GoName:      "Networking",
								GoType:      "ClusterNetworking",
								Description: "Cluster-wide networking configuration.",
								Required: []string{
									"vpcId",
									"podAddressCidrBlocks",
									"serviceAddressCidrBlocks",
								},
								Properties: map[string]*dcl.Property{
									"perNodePoolSgRulesDisabled": &dcl.Property{
										Type:        "boolean",
										GoName:      "PerNodePoolSgRulesDisabled",
										Description: "Disable the per node pool subnet security group rules on the control plane security group. When set to true, you must also provide one or more security groups that ensure node pools are able to send requests to the control plane on TCP/443 and TCP/8132. Failure to do so may result in unavailable node pools.",
									},
									"podAddressCidrBlocks": &dcl.Property{
										Type:        "array",
										GoName:      "PodAddressCidrBlocks",
										Description: "All pods in the cluster are assigned an RFC1918 IPv4 address from these ranges. Only a single range is supported. This field cannot be changed after creation.",
										Immutable:   true,
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "string",
											GoType: "string",
										},
									},
									"serviceAddressCidrBlocks": &dcl.Property{
										Type:        "array",
										GoName:      "ServiceAddressCidrBlocks",
										Description: "All services in the cluster are assigned an RFC1918 IPv4 address from these ranges. Only a single range is supported. This field cannot be changed after creation.",
										Immutable:   true,
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "string",
											GoType: "string",
										},
									},
									"vpcId": &dcl.Property{
										Type:        "string",
										GoName:      "VPCId",
										Description: "The VPC associated with the cluster. All component clusters (i.e. control plane and node pools) run on a single VPC. This field cannot be changed after creation.",
										Immutable:   true,
									},
								},
							},
							"project": &dcl.Property{
								Type:        "string",
								GoName:      "Project",
								Description: "The project for the resource",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Cloudresourcemanager/Project",
										Field:    "name",
										Parent:   true,
									},
								},
							},
							"reconciling": &dcl.Property{
								Type:        "boolean",
								GoName:      "Reconciling",
								ReadOnly:    true,
								Description: "Output only. If set, there are currently changes in flight to the cluster.",
								Immutable:   true,
							},
							"state": &dcl.Property{
								Type:        "string",
								GoName:      "State",
								GoType:      "ClusterStateEnum",
								ReadOnly:    true,
								Description: "Output only. The current state of the cluster. Possible values: STATE_UNSPECIFIED, PROVISIONING, RUNNING, RECONCILING, STOPPING, ERROR, DEGRADED",
								Immutable:   true,
								Enum: []string{
									"STATE_UNSPECIFIED",
									"PROVISIONING",
									"RUNNING",
									"RECONCILING",
									"STOPPING",
									"ERROR",
									"DEGRADED",
								},
							},
							"uid": &dcl.Property{
								Type:        "string",
								GoName:      "Uid",
								ReadOnly:    true,
								Description: "Output only. A globally unique identifier for the cluster.",
								Immutable:   true,
							},
							"updateTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "UpdateTime",
								ReadOnly:    true,
								Description: "Output only. The time at which this cluster was last updated.",
								Immutable:   true,
							},
							"workloadIdentityConfig": &dcl.Property{
								Type:        "object",
								GoName:      "WorkloadIdentityConfig",
								GoType:      "ClusterWorkloadIdentityConfig",
								ReadOnly:    true,
								Description: "Output only. Workload Identity settings.",
								Immutable:   true,
								Properties: map[string]*dcl.Property{
									"identityProvider": &dcl.Property{
										Type:        "string",
										GoName:      "IdentityProvider",
										Description: "The ID of the OIDC Identity Provider (IdP) associated to the Workload Identity Pool.",
										Immutable:   true,
									},
									"issuerUri": &dcl.Property{
										Type:        "string",
										GoName:      "IssuerUri",
										Description: "The OIDC issuer URL for this cluster.",
										Immutable:   true,
									},
									"workloadPool": &dcl.Property{
										Type:        "string",
										GoName:      "WorkloadPool",
										Description: "The Workload Identity Pool associated to the cluster.",
										Immutable:   true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
