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

func DCLNodePoolSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "ContainerAzure/NodePool",
			Description: "An Anthos node pool running on Azure.",
			StructName:  "NodePool",
			Reference: &dcl.Link{
				Text: "API reference",
				URL:  "https://cloud.google.com/kubernetes-engine/multi-cloud/docs/reference/rest/v1/projects.locations.azureClusters.azureNodePools",
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
				Description: "The function used to get information about a NodePool",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "nodePool",
						Required:    true,
						Description: "A full instance of a NodePool",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a NodePool",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "nodePool",
						Required:    true,
						Description: "A full instance of a NodePool",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a NodePool",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "nodePool",
						Required:    true,
						Description: "A full instance of a NodePool",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all NodePool",
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
					dcl.PathParameters{
						Name:     "cluster",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
				},
			},
			List: &dcl.Path{
				Description: "The function used to list information about many NodePool",
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
					dcl.PathParameters{
						Name:     "cluster",
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
				"NodePool": &dcl.Component{
					Title:           "NodePool",
					ID:              "projects/{{project}}/locations/{{location}}/azureClusters/{{cluster}}/azureNodePools/{{name}}",
					ParentContainer: "project",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"name",
							"version",
							"config",
							"subnetId",
							"autoscaling",
							"maxPodsConstraint",
							"project",
							"location",
							"cluster",
						},
						Properties: map[string]*dcl.Property{
							"annotations": &dcl.Property{
								Type: "object",
								AdditionalProperties: &dcl.Property{
									Type: "string",
								},
								GoName:      "Annotations",
								Description: "Optional. Annotations on the node pool. This field has the same restrictions as Kubernetes annotations. The total size of all keys and values combined is limited to 256k. Keys can have 2 segments: prefix (optional) and name (required), separated by a slash (/). Prefix must be a DNS subdomain. Name must be 63 characters or less, begin and end with alphanumerics, with dashes (-), underscores (_), dots (.), and alphanumerics between.",
							},
							"autoscaling": &dcl.Property{
								Type:        "object",
								GoName:      "Autoscaling",
								GoType:      "NodePoolAutoscaling",
								Description: "Autoscaler configuration for this node pool.",
								Required: []string{
									"minNodeCount",
									"maxNodeCount",
								},
								Properties: map[string]*dcl.Property{
									"maxNodeCount": &dcl.Property{
										Type:        "integer",
										Format:      "int64",
										GoName:      "MaxNodeCount",
										Description: "Maximum number of nodes in the node pool. Must be >= min_node_count.",
									},
									"minNodeCount": &dcl.Property{
										Type:        "integer",
										Format:      "int64",
										GoName:      "MinNodeCount",
										Description: "Minimum number of nodes in the node pool. Must be >= 1 and <= max_node_count.",
									},
								},
							},
							"azureAvailabilityZone": &dcl.Property{
								Type:          "string",
								GoName:        "AzureAvailabilityZone",
								Description:   "Optional. The Azure availability zone of the nodes in this nodepool. When unspecified, it defaults to `1`.",
								Immutable:     true,
								ServerDefault: true,
							},
							"cluster": &dcl.Property{
								Type:        "string",
								GoName:      "Cluster",
								Description: "The azureCluster for the resource",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Gkemulticloud/Cluster",
										Field:    "name",
										Parent:   true,
									},
								},
								Parameter: true,
							},
							"config": &dcl.Property{
								Type:        "object",
								GoName:      "Config",
								GoType:      "NodePoolConfig",
								Description: "The node configuration of the node pool.",
								Required: []string{
									"sshConfig",
								},
								Properties: map[string]*dcl.Property{
									"labels": &dcl.Property{
										Type: "object",
										AdditionalProperties: &dcl.Property{
											Type: "string",
										},
										GoName:      "Labels",
										Description: "Optional. The initial labels assigned to nodes of this node pool. An object containing a list of \"key\": value pairs. Example: { \"name\": \"wrench\", \"mass\": \"1.3kg\", \"count\": \"3\" }.",
										Immutable:   true,
									},
									"proxyConfig": &dcl.Property{
										Type:        "object",
										GoName:      "ProxyConfig",
										GoType:      "NodePoolConfigProxyConfig",
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
									"rootVolume": &dcl.Property{
										Type:          "object",
										GoName:        "RootVolume",
										GoType:        "NodePoolConfigRootVolume",
										Description:   "Optional. Configuration related to the root volume provisioned for each node pool machine. When unspecified, it defaults to a 32-GiB Azure Disk.",
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
										GoType:      "NodePoolConfigSshConfig",
										Description: "SSH configuration for how to access the node pool machines.",
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
									"tags": &dcl.Property{
										Type: "object",
										AdditionalProperties: &dcl.Property{
											Type: "string",
										},
										GoName:      "Tags",
										Description: "Optional. A set of tags to apply to all underlying Azure resources for this node pool. This currently only includes Virtual Machine Scale Sets. Specify at most 50 pairs containing alphanumerics, spaces, and symbols (.+-=_:@/). Keys can be up to 127 Unicode characters. Values can be up to 255 Unicode characters.",
										Immutable:   true,
									},
									"vmSize": &dcl.Property{
										Type:          "string",
										GoName:        "VmSize",
										Description:   "Optional. The Azure VM size name. Example: `Standard_DS2_v2`. See (/anthos/clusters/docs/azure/reference/supported-vms) for options. When unspecified, it defaults to `Standard_DS2_v2`.",
										Immutable:     true,
										ServerDefault: true,
									},
								},
							},
							"createTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "CreateTime",
								ReadOnly:    true,
								Description: "Output only. The time at which this node pool was created.",
								Immutable:   true,
							},
							"etag": &dcl.Property{
								Type:        "string",
								GoName:      "Etag",
								ReadOnly:    true,
								Description: "Allows clients to perform consistent read-modify-writes through optimistic concurrency control. May be sent on update and delete requests to ensure the client has an up-to-date value before proceeding.",
								Immutable:   true,
							},
							"location": &dcl.Property{
								Type:        "string",
								GoName:      "Location",
								Description: "The location for the resource",
								Immutable:   true,
								Parameter:   true,
							},
							"management": &dcl.Property{
								Type:        "object",
								GoName:      "Management",
								GoType:      "NodePoolManagement",
								Description: "The Management configuration for this node pool.",
								Properties: map[string]*dcl.Property{
									"autoRepair": &dcl.Property{
										Type:        "boolean",
										GoName:      "AutoRepair",
										Description: "Optional. Whether or not the nodes will be automatically repaired.",
									},
								},
							},
							"maxPodsConstraint": &dcl.Property{
								Type:        "object",
								GoName:      "MaxPodsConstraint",
								GoType:      "NodePoolMaxPodsConstraint",
								Description: "The constraint on the maximum number of pods that can be run simultaneously on a node in the node pool.",
								Immutable:   true,
								Required: []string{
									"maxPodsPerNode",
								},
								Properties: map[string]*dcl.Property{
									"maxPodsPerNode": &dcl.Property{
										Type:        "integer",
										Format:      "int64",
										GoName:      "MaxPodsPerNode",
										Description: "The maximum number of pods to schedule on a single node.",
										Immutable:   true,
									},
								},
							},
							"name": &dcl.Property{
								Type:        "string",
								GoName:      "Name",
								Description: "The name of this resource.",
								Immutable:   true,
								HasLongForm: true,
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
								Description: "Output only. If set, there are currently pending changes to the node pool.",
								Immutable:   true,
							},
							"state": &dcl.Property{
								Type:        "string",
								GoName:      "State",
								GoType:      "NodePoolStateEnum",
								ReadOnly:    true,
								Description: "Output only. The current state of the node pool. Possible values: STATE_UNSPECIFIED, PROVISIONING, RUNNING, RECONCILING, STOPPING, ERROR, DEGRADED",
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
							"subnetId": &dcl.Property{
								Type:        "string",
								GoName:      "SubnetId",
								Description: "The ARM ID of the subnet where the node pool VMs run. Make sure it's a subnet under the virtual network in the cluster configuration.",
								Immutable:   true,
							},
							"uid": &dcl.Property{
								Type:        "string",
								GoName:      "Uid",
								ReadOnly:    true,
								Description: "Output only. A globally unique identifier for the node pool.",
								Immutable:   true,
							},
							"updateTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "UpdateTime",
								ReadOnly:    true,
								Description: "Output only. The time at which this node pool was last updated.",
								Immutable:   true,
							},
							"version": &dcl.Property{
								Type:        "string",
								GoName:      "Version",
								Description: "The Kubernetes version (e.g. `1.19.10-gke.1000`) running on this node pool.",
							},
						},
					},
				},
			},
		},
	}
}
