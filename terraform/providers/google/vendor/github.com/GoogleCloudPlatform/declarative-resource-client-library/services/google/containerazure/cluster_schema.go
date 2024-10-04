// Copyright 2024 Google LLC. All Rights Reserved.
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
package containerazure

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func DCLClusterSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "ContainerAzure/Cluster",
			Description: "An Anthos cluster running on Azure.",
			StructName:  "Cluster",
			Reference: &dcl.Link{
				Text: "API reference",
				URL:  "https://cloud.google.com/kubernetes-engine/multi-cloud/docs/reference/rest/v1/projects.locations.azureClusters",
			},
			Guides: []*dcl.Link{
				&dcl.Link{
					Text: "Multicloud overview",
					URL:  "https://cloud.google.com/kubernetes-engine/multi-cloud/docs",
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
					ID:              "projects/{{project}}/locations/{{location}}/azureClusters/{{name}}",
					ParentContainer: "project",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"name",
							"azureRegion",
							"resourceGroupId",
							"networking",
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
								Description: "Optional. Annotations on the cluster. This field has the same restrictions as Kubernetes annotations. The total size of all keys and values combined is limited to 256k. Keys can have 2 segments: prefix (optional) and name (required), separated by a slash (/). Prefix must be a DNS subdomain. Name must be 63 characters or less, begin and end with alphanumerics, with dashes (-), underscores (_), dots (.), and alphanumerics between.",
								Immutable:   true,
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
									"adminGroups": &dcl.Property{
										Type:        "array",
										GoName:      "AdminGroups",
										Description: "Groups of users that can perform operations as a cluster admin. A managed ClusterRoleBinding will be created to grant the `cluster-admin` ClusterRole to the groups. Up to ten admin groups can be provided. For more info on RBAC, see https://kubernetes.io/docs/reference/access-authn-authz/rbac/#user-facing-roles",
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "object",
											GoType: "ClusterAuthorizationAdminGroups",
											Required: []string{
												"group",
											},
											Properties: map[string]*dcl.Property{
												"group": &dcl.Property{
													Type:        "string",
													GoName:      "Group",
													Description: "The name of the group, e.g. `my-group@domain.com`.",
												},
											},
										},
									},
									"adminUsers": &dcl.Property{
										Type:        "array",
										GoName:      "AdminUsers",
										Description: "Users that can perform operations as a cluster admin. A new ClusterRoleBinding will be created to grant the cluster-admin ClusterRole to the users. Up to ten admin users can be provided. For more info on RBAC, see https://kubernetes.io/docs/reference/access-authn-authz/rbac/#user-facing-roles",
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
							"azureRegion": &dcl.Property{
								Type:        "string",
								GoName:      "AzureRegion",
								Description: "The Azure region where the cluster runs. Each Google Cloud region supports a subset of nearby Azure regions. You can call to list all supported Azure regions within a given Google Cloud region.",
								Immutable:   true,
							},
							"azureServicesAuthentication": &dcl.Property{
								Type:        "object",
								GoName:      "AzureServicesAuthentication",
								GoType:      "ClusterAzureServicesAuthentication",
								Description: "Azure authentication configuration for management of Azure resources",
								Conflicts: []string{
									"client",
								},
								Required: []string{
									"tenantId",
									"applicationId",
								},
								Properties: map[string]*dcl.Property{
									"applicationId": &dcl.Property{
										Type:        "string",
										GoName:      "ApplicationId",
										Description: "The Azure Active Directory Application ID for Authentication configuration.",
									},
									"tenantId": &dcl.Property{
										Type:        "string",
										GoName:      "TenantId",
										Description: "The Azure Active Directory Tenant ID for Authentication configuration.",
									},
								},
							},
							"client": &dcl.Property{
								Type:        "string",
								GoName:      "Client",
								Description: "Name of the AzureClient. The `AzureClient` resource must reside on the same GCP project and region as the `AzureCluster`. `AzureClient` names are formatted as `projects/<project-number>/locations/<region>/azureClients/<client-id>`. See Resource Names (https:cloud.google.com/apis/design/resource_names) for more details on Google Cloud resource names.",
								Conflicts: []string{
									"azureServicesAuthentication",
								},
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "ContainerAzure/AzureClient",
										Field:    "name",
									},
								},
							},
							"controlPlane": &dcl.Property{
								Type:        "object",
								GoName:      "ControlPlane",
								GoType:      "ClusterControlPlane",
								Description: "Configuration related to the cluster control plane.",
								Required: []string{
									"version",
									"subnetId",
									"sshConfig",
								},
								Properties: map[string]*dcl.Property{
									"databaseEncryption": &dcl.Property{
										Type:        "object",
										GoName:      "DatabaseEncryption",
										GoType:      "ClusterControlPlaneDatabaseEncryption",
										Description: "Optional. Configuration related to application-layer secrets encryption.",
										Immutable:   true,
										Required: []string{
											"keyId",
										},
										Properties: map[string]*dcl.Property{
											"keyId": &dcl.Property{
												Type:        "string",
												GoName:      "KeyId",
												Description: "The ARM ID of the Azure Key Vault key to encrypt / decrypt data. For example: `/subscriptions/<subscription-id>/resourceGroups/<resource-group-id>/providers/Microsoft.KeyVault/vaults/<key-vault-id>/keys/<key-name>` Encryption will always take the latest version of the key and hence specific version is not supported.",
												Immutable:   true,
											},
										},
									},
									"mainVolume": &dcl.Property{
										Type:          "object",
										GoName:        "MainVolume",
										GoType:        "ClusterControlPlaneMainVolume",
										Description:   "Optional. Configuration related to the main volume provisioned for each control plane replica. The main volume is in charge of storing all of the cluster's etcd state. When unspecified, it defaults to a 8-GiB Azure Disk.",
										Immutable:     true,
										ServerDefault: true,
										Properties: map[string]*dcl.Property{
											"sizeGib": &dcl.Property{
												Type:          "integer",
												Format:        "int64",
												GoName:        "SizeGib",
												Description:   "Optional. The size of the disk, in GiBs. When unspecified, a default value is provided. See the specific reference in the parent resource.",
												Immutable:     true,
												ServerDefault: true,
											},
										},
									},
									"proxyConfig": &dcl.Property{
										Type:        "object",
										GoName:      "ProxyConfig",
										GoType:      "ClusterControlPlaneProxyConfig",
										Description: "Proxy configuration for outbound HTTP(S) traffic.",
										Immutable:   true,
										Required: []string{
											"resourceGroupId",
											"secretId",
										},
										Properties: map[string]*dcl.Property{
											"resourceGroupId": &dcl.Property{
												Type:        "string",
												GoName:      "ResourceGroupId",
												Description: "The ARM ID the of the resource group containing proxy keyvault. Resource group ids are formatted as `/subscriptions/<subscription-id>/resourceGroups/<resource-group-name>`",
												Immutable:   true,
											},
											"secretId": &dcl.Property{
												Type:        "string",
												GoName:      "SecretId",
												Description: "The URL the of the proxy setting secret with its version. Secret ids are formatted as `https:<key-vault-name>.vault.azure.net/secrets/<secret-name>/<secret-version>`.",
												Immutable:   true,
											},
										},
									},
									"replicaPlacements": &dcl.Property{
										Type:        "array",
										GoName:      "ReplicaPlacements",
										Description: "Configuration for where to place the control plane replicas. Up to three replica placement instances can be specified. If replica_placements is set, the replica placement instances will be applied to the three control plane replicas as evenly as possible.",
										Immutable:   true,
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "object",
											GoType: "ClusterControlPlaneReplicaPlacements",
											Required: []string{
												"subnetId",
												"azureAvailabilityZone",
											},
											Properties: map[string]*dcl.Property{
												"azureAvailabilityZone": &dcl.Property{
													Type:        "string",
													GoName:      "AzureAvailabilityZone",
													Description: "For a given replica, the Azure availability zone where to provision the control plane VM and the ETCD disk.",
													Immutable:   true,
												},
												"subnetId": &dcl.Property{
													Type:        "string",
													GoName:      "SubnetId",
													Description: "For a given replica, the ARM ID of the subnet where the control plane VM is deployed. Make sure it's a subnet under the virtual network in the cluster configuration.",
													Immutable:   true,
												},
											},
										},
									},
									"rootVolume": &dcl.Property{
										Type:          "object",
										GoName:        "RootVolume",
										GoType:        "ClusterControlPlaneRootVolume",
										Description:   "Optional. Configuration related to the root volume provisioned for each control plane replica. When unspecified, it defaults to 32-GiB Azure Disk.",
										Immutable:     true,
										ServerDefault: true,
										Properties: map[string]*dcl.Property{
											"sizeGib": &dcl.Property{
												Type:          "integer",
												Format:        "int64",
												GoName:        "SizeGib",
												Description:   "Optional. The size of the disk, in GiBs. When unspecified, a default value is provided. See the specific reference in the parent resource.",
												Immutable:     true,
												ServerDefault: true,
											},
										},
									},
									"sshConfig": &dcl.Property{
										Type:        "object",
										GoName:      "SshConfig",
										GoType:      "ClusterControlPlaneSshConfig",
										Description: "SSH configuration for how to access the underlying control plane machines.",
										Required: []string{
											"authorizedKey",
										},
										Properties: map[string]*dcl.Property{
											"authorizedKey": &dcl.Property{
												Type:        "string",
												GoName:      "AuthorizedKey",
												Description: "The SSH public key data for VMs managed by Anthos. This accepts the authorized_keys file format used in OpenSSH according to the sshd(8) manual page.",
											},
										},
									},
									"subnetId": &dcl.Property{
										Type:        "string",
										GoName:      "SubnetId",
										Description: "The ARM ID of the subnet where the control plane VMs are deployed. Example: `/subscriptions//resourceGroups//providers/Microsoft.Network/virtualNetworks//subnets/default`.",
										Immutable:   true,
									},
									"tags": &dcl.Property{
										Type: "object",
										AdditionalProperties: &dcl.Property{
											Type: "string",
										},
										GoName:      "Tags",
										Description: "Optional. A set of tags to apply to all underlying control plane Azure resources.",
										Immutable:   true,
									},
									"version": &dcl.Property{
										Type:        "string",
										GoName:      "Version",
										Description: "The Kubernetes version to run on control plane replicas (e.g. `1.19.10-gke.1000`). You can list all supported versions on a given Google Cloud region by calling GetAzureServerConfig.",
									},
									"vmSize": &dcl.Property{
										Type:          "string",
										GoName:        "VmSize",
										Description:   "Optional. The Azure VM size name. Example: `Standard_DS2_v2`. For available VM sizes, see https://docs.microsoft.com/en-us/azure/virtual-machines/vm-naming-conventions. When unspecified, it defaults to `Standard_DS2_v2`.",
										ServerDefault: true,
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
										HasLongForm: true,
									},
								},
							},
							"location": &dcl.Property{
								Type:        "string",
								GoName:      "Location",
								Description: "The location for the resource",
								Immutable:   true,
								Parameter:   true,
							},
							"name": &dcl.Property{
								Type:        "string",
								GoName:      "Name",
								Description: "The name of this resource.",
								Immutable:   true,
								HasLongForm: true,
							},
							"networking": &dcl.Property{
								Type:        "object",
								GoName:      "Networking",
								GoType:      "ClusterNetworking",
								Description: "Cluster-wide networking configuration.",
								Immutable:   true,
								Required: []string{
									"virtualNetworkId",
									"podAddressCidrBlocks",
									"serviceAddressCidrBlocks",
								},
								Properties: map[string]*dcl.Property{
									"podAddressCidrBlocks": &dcl.Property{
										Type:        "array",
										GoName:      "PodAddressCidrBlocks",
										Description: "The IP address range of the pods in this cluster, in CIDR notation (e.g. `10.96.0.0/14`). All pods in the cluster get assigned a unique RFC1918 IPv4 address from these ranges. Only a single range is supported. This field cannot be changed after creation.",
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
										Description: "The IP address range for services in this cluster, in CIDR notation (e.g. `10.96.0.0/14`). All services in the cluster get assigned a unique RFC1918 IPv4 address from these ranges. Only a single range is supported. This field cannot be changed after creating a cluster.",
										Immutable:   true,
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "string",
											GoType: "string",
										},
									},
									"virtualNetworkId": &dcl.Property{
										Type:        "string",
										GoName:      "VirtualNetworkId",
										Description: "The Azure Resource Manager (ARM) ID of the VNet associated with your cluster. All components in the cluster (i.e. control plane and node pools) run on a single VNet. Example: `/subscriptions/*/resourceGroups/*/providers/Microsoft.Network/virtualNetworks/*` This field cannot be changed after creation.",
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
								Parameter: true,
							},
							"reconciling": &dcl.Property{
								Type:        "boolean",
								GoName:      "Reconciling",
								ReadOnly:    true,
								Description: "Output only. If set, there are currently changes in flight to the cluster.",
								Immutable:   true,
							},
							"resourceGroupId": &dcl.Property{
								Type:        "string",
								GoName:      "ResourceGroupId",
								Description: "The ARM ID of the resource group where the cluster resources are deployed. For example: `/subscriptions/*/resourceGroups/*`",
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
