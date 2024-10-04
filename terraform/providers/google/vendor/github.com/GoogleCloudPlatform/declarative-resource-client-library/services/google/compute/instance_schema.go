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
package compute

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func DCLInstanceSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Compute/Instance",
			Description: "The Compute Instance resource",
			StructName:  "Instance",
			HasIAM:      true,
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a Instance",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "instance",
						Required:    true,
						Description: "A full instance of a Instance",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a Instance",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "instance",
						Required:    true,
						Description: "A full instance of a Instance",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a Instance",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "instance",
						Required:    true,
						Description: "A full instance of a Instance",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all Instance",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:     "project",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
					dcl.PathParameters{
						Name:     "zone",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
				},
			},
			List: &dcl.Path{
				Description: "The function used to list information about many Instance",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:     "project",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
					dcl.PathParameters{
						Name:     "zone",
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
				"Instance": &dcl.Component{
					Title: "Instance",
					ID:    "projects/{{project}}/zones/{{zone}}/instances/{{name}}",
					Locations: []string{
						"zone",
					},
					UsesStateHint:   true,
					ParentContainer: "project",
					LabelsField:     "labels",
					HasCreate:       true,
					HasIAM:          true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"zone",
							"project",
						},
						Properties: map[string]*dcl.Property{
							"canIPForward": &dcl.Property{
								Type:        "boolean",
								GoName:      "CanIPForward",
								Description: "Allows this instance to send and receive packets with non-matching destination or source IPs. This is required if you plan to use this instance to forward routes.",
								Immutable:   true,
							},
							"cpuPlatform": &dcl.Property{
								Type:        "string",
								GoName:      "CpuPlatform",
								ReadOnly:    true,
								Description: "The CPU platform used by this instance.",
								Immutable:   true,
							},
							"creationTimestamp": &dcl.Property{
								Type:        "string",
								GoName:      "CreationTimestamp",
								ReadOnly:    true,
								Description: "Creation timestamp in RFC3339 text format.",
								Immutable:   true,
							},
							"deletionProtection": &dcl.Property{
								Type:        "boolean",
								GoName:      "DeletionProtection",
								Description: "Whether the resource should be protected against deletion.",
							},
							"description": &dcl.Property{
								Type:        "string",
								GoName:      "Description",
								Description: "An optional description of this resource.",
								Immutable:   true,
							},
							"disks": &dcl.Property{
								Type:        "array",
								GoName:      "Disks",
								Description: "An array of disks that are associated with the instances that are created from this template.",
								Immutable:   true,
								SendEmpty:   true,
								ListType:    "list",
								Items: &dcl.Property{
									Type:   "object",
									GoType: "InstanceDisks",
									Properties: map[string]*dcl.Property{
										"autoDelete": &dcl.Property{
											Type:        "boolean",
											GoName:      "AutoDelete",
											Description: "Specifies whether the disk will be auto-deleted when the instance is deleted (but not when the disk is detached from the instance).  Tip: Disks should be set to autoDelete=true so that leftover disks are not left behind on machine deletion.",
											Immutable:   true,
										},
										"boot": &dcl.Property{
											Type:        "boolean",
											GoName:      "Boot",
											Description: "Indicates that this is a boot disk. The virtual machine will use the first partition of the disk for its root filesystem.",
											Immutable:   true,
										},
										"deviceName": &dcl.Property{
											Type:        "string",
											GoName:      "DeviceName",
											Description: "Specifies a unique device name of your choice that is reflected into the /dev/disk/by-id/google-* tree of a Linux operating system running within the instance. This name can be used to reference the device for mounting, resizing, and so on, from within the instance.",
											Immutable:   true,
										},
										"diskEncryptionKey": &dcl.Property{
											Type:        "object",
											GoName:      "DiskEncryptionKey",
											GoType:      "InstanceDisksDiskEncryptionKey",
											Description: "Encrypts or decrypts a disk using a customer-supplied encryption key.",
											Immutable:   true,
											Properties: map[string]*dcl.Property{
												"rawKey": &dcl.Property{
													Type:        "string",
													GoName:      "RawKey",
													Description: "Specifies a 256-bit customer-supplied encryption key, encoded in RFC 4648 base64 to either encrypt or decrypt this resource.",
													Immutable:   true,
												},
												"rsaEncryptedKey": &dcl.Property{
													Type:        "string",
													GoName:      "RsaEncryptedKey",
													Description: "Specifies an RFC 4648 base64 encoded, RSA-wrapped 2048-bit customer-supplied encryption key to either encrypt or decrypt this resource.",
													Immutable:   true,
												},
												"sha256": &dcl.Property{
													Type:        "string",
													GoName:      "Sha256",
													ReadOnly:    true,
													Description: "The RFC 4648 base64 encoded SHA-256 hash of the customer-supplied encryption key that protects this resource.",
													Immutable:   true,
												},
											},
										},
										"index": &dcl.Property{
											Type:        "integer",
											Format:      "int64",
											GoName:      "Index",
											Description: "Assigns a zero-based index to this disk, where 0 is reserved for the boot disk. For example, if you have many disks attached to an instance, each disk would have a unique index number. If not specified, the server will choose an appropriate value.",
											Immutable:   true,
										},
										"initializeParams": &dcl.Property{
											Type:        "object",
											GoName:      "InitializeParams",
											GoType:      "InstanceDisksInitializeParams",
											Description: "Specifies the parameters for a new disk that will be created alongside the new instance. Use initialization parameters to create boot disks or local SSDs attached to the new instance.",
											Immutable:   true,
											Unreadable:  true,
											Properties: map[string]*dcl.Property{
												"diskName": &dcl.Property{
													Type:        "string",
													GoName:      "DiskName",
													Description: "Specifies the disk name. If not specified, the default is to use the name of the instance.",
													Immutable:   true,
												},
												"diskSizeGb": &dcl.Property{
													Type:        "integer",
													Format:      "int64",
													GoName:      "DiskSizeGb",
													Description: "Specifies the size of the disk in base-2 GB.",
													Immutable:   true,
												},
												"diskType": &dcl.Property{
													Type:        "string",
													GoName:      "DiskType",
													Description: "Reference to a disk type. Specifies the disk type to use to create the instance. If not specified, the default is pd-standard.",
													Immutable:   true,
													ResourceReferences: []*dcl.PropertyResourceReference{
														&dcl.PropertyResourceReference{
															Resource: "Compute/DiskType",
															Field:    "name",
														},
													},
												},
												"sourceImage": &dcl.Property{
													Type:        "string",
													GoName:      "SourceImage",
													Description: "The source image to create this disk. When creating a new instance, one of initializeParams.sourceImage or disks.source is required.  To create a disk with one of the public operating system images, specify the image by its family name.",
													Immutable:   true,
												},
												"sourceImageEncryptionKey": &dcl.Property{
													Type:        "object",
													GoName:      "SourceImageEncryptionKey",
													GoType:      "InstanceDisksInitializeParamsSourceImageEncryptionKey",
													Description: "The customer-supplied encryption key of the source image. Required if the source image is protected by a customer-supplied encryption key.  Instance templates do not store customer-supplied encryption keys, so you cannot create disks for instances in a managed instance group if the source images are encrypted with your own keys.",
													Immutable:   true,
													Properties: map[string]*dcl.Property{
														"rawKey": &dcl.Property{
															Type:        "string",
															GoName:      "RawKey",
															Description: "Specifies a 256-bit customer-supplied encryption key, encoded in RFC 4648 base64 to either encrypt or decrypt this resource.",
															Immutable:   true,
														},
														"sha256": &dcl.Property{
															Type:        "string",
															GoName:      "Sha256",
															ReadOnly:    true,
															Description: "The RFC 4648 base64 encoded SHA-256 hash of the customer-supplied encryption key that protects this resource.",
															Immutable:   true,
														},
													},
												},
											},
										},
										"interface": &dcl.Property{
											Type:        "string",
											GoName:      "Interface",
											GoType:      "InstanceDisksInterfaceEnum",
											Description: "Specifies the disk interface to use for attaching this disk, which is either SCSI or NVME. The default is SCSI. Persistent disks must always use SCSI and the request will fail if you attempt to attach a persistent disk in any other format than SCSI.",
											Immutable:   true,
											Enum: []string{
												"SCSI",
												"NVME",
											},
										},
										"mode": &dcl.Property{
											Type:        "string",
											GoName:      "Mode",
											GoType:      "InstanceDisksModeEnum",
											Description: "The mode in which to attach this disk, either READ_WRITE or READ_ONLY. If not specified, the default is to attach the disk in READ_WRITE mode.",
											Immutable:   true,
											Enum: []string{
												"READ_WRITE",
												"READ_ONLY",
											},
										},
										"source": &dcl.Property{
											Type:        "string",
											GoName:      "Source",
											Description: "Reference to a disk. When creating a new instance, one of initializeParams.sourceImage or disks.source is required.  If desired, you can also attach existing non-root persistent disks using this property. This field is only applicable for persistent disks.",
											Immutable:   true,
											ResourceReferences: []*dcl.PropertyResourceReference{
												&dcl.PropertyResourceReference{
													Resource: "Compute/Disk",
													Field:    "selfLink",
												},
											},
										},
										"type": &dcl.Property{
											Type:        "string",
											GoName:      "Type",
											GoType:      "InstanceDisksTypeEnum",
											Description: "Specifies the type of the disk, either SCRATCH or PERSISTENT. If not specified, the default is PERSISTENT.",
											Immutable:   true,
											Enum: []string{
												"SCRATCH",
												"PERSISTENT",
											},
										},
									},
								},
								Unreadable: true,
							},
							"guestAccelerators": &dcl.Property{
								Type:        "array",
								GoName:      "GuestAccelerators",
								Description: "List of the type and count of accelerator cards attached to the instance",
								Immutable:   true,
								SendEmpty:   true,
								ListType:    "list",
								Items: &dcl.Property{
									Type:   "object",
									GoType: "InstanceGuestAccelerators",
									Properties: map[string]*dcl.Property{
										"acceleratorCount": &dcl.Property{
											Type:        "integer",
											Format:      "int64",
											GoName:      "AcceleratorCount",
											Description: "The number of the guest accelerator cards exposed to this instance.",
											Immutable:   true,
										},
										"acceleratorType": &dcl.Property{
											Type:        "string",
											GoName:      "AcceleratorType",
											Description: "Full or partial URL of the accelerator type resource to expose to this instance.",
											Immutable:   true,
										},
									},
								},
							},
							"hostname": &dcl.Property{
								Type:        "string",
								GoName:      "Hostname",
								Description: "The hostname of the instance to be created. The specified hostname must be RFC1035 compliant. If hostname is not specified, the default hostname is [INSTANCE_NAME].c.[PROJECT_ID].internal when using the global DNS, and [INSTANCE_NAME].[ZONE].c.[PROJECT_ID].internal when using zonal DNS.",
								Immutable:   true,
							},
							"id": &dcl.Property{
								Type:        "string",
								GoName:      "Id",
								ReadOnly:    true,
								Description: "The unique identifier for the resource. This identifier is defined by the server.",
								Immutable:   true,
							},
							"labels": &dcl.Property{
								Type: "object",
								AdditionalProperties: &dcl.Property{
									Type: "string",
								},
								GoName:      "Labels",
								Description: "Labels to apply to this instance.  A list of key->value pairs.",
							},
							"machineType": &dcl.Property{
								Type:                "string",
								GoName:              "MachineType",
								Description:         "A reference to a machine type which defines VM kind.",
								ForwardSlashAllowed: true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Compute/MachineType",
										Field:    "name",
									},
								},
								HasLongForm: true,
							},
							"metadata": &dcl.Property{
								Type: "object",
								AdditionalProperties: &dcl.Property{
									Type: "string",
								},
								GoName:      "Metadata",
								Description: "The metadata key/value pairs to assign to instances that are created from this template. These pairs can consist of custom metadata or predefined keys.",
							},
							"minCpuPlatform": &dcl.Property{
								Type:        "string",
								GoName:      "MinCpuPlatform",
								Description: "Specifies a minimum CPU platform for the VM instance. Applicable values are the friendly names of CPU platforms",
								Immutable:   true,
							},
							"name": &dcl.Property{
								Type:        "string",
								GoName:      "Name",
								Description: "The name of the resource, provided by the client when initially creating the resource. The resource name must be 1-63 characters long, and comply with RFC1035. Specifically, the name must be 1-63 characters long and match the regular expression `[a-z]([-a-z0-9]*[a-z0-9])?` which means the first character must be a lowercase letter, and all following characters must be a dash, lowercase letter, or digit, except the last character, which cannot be a dash.",
								Immutable:   true,
							},
							"networkInterfaces": &dcl.Property{
								Type:        "array",
								GoName:      "NetworkInterfaces",
								Description: "An array of configurations for this interface. This specifies how this interface is configured to interact with other network services, such as connecting to the internet. Only one network interface is supported per instance.",
								Immutable:   true,
								SendEmpty:   true,
								ListType:    "list",
								Items: &dcl.Property{
									Type:   "object",
									GoType: "InstanceNetworkInterfaces",
									Properties: map[string]*dcl.Property{
										"accessConfigs": &dcl.Property{
											Type:        "array",
											GoName:      "AccessConfigs",
											Description: "An array of configurations for this interface. Currently, only one access config, ONE_TO_ONE_NAT, is supported. If there are no accessConfigs specified, then this instance will have no external internet access.",
											Immutable:   true,
											SendEmpty:   true,
											ListType:    "list",
											Items: &dcl.Property{
												Type:   "object",
												GoType: "InstanceNetworkInterfacesAccessConfigs",
												Required: []string{
													"name",
													"type",
												},
												Properties: map[string]*dcl.Property{
													"externalIPv6": &dcl.Property{
														Type:        "string",
														GoName:      "ExternalIPv6",
														ReadOnly:    true,
														Description: "The first IPv6 address of the external IPv6 range associated with this instance, prefix length is stored in externalIpv6PrefixLength in ipv6AccessConfig. The field is output only, an IPv6 address from a subnetwork associated with the instance will be allocated dynamically.",
														Immutable:   true,
													},
													"externalIPv6PrefixLength": &dcl.Property{
														Type:        "string",
														GoName:      "ExternalIPv6PrefixLength",
														ReadOnly:    true,
														Description: "The prefix length of the external IPv6 range.",
														Immutable:   true,
													},
													"name": &dcl.Property{
														Type:        "string",
														GoName:      "Name",
														Description: "The name of this access configuration. The default and recommended name is External NAT but you can use any arbitrary string you would like. For example, My external IP or Network Access.",
														Immutable:   true,
													},
													"natIP": &dcl.Property{
														Type:          "string",
														GoName:        "NatIP",
														Description:   "Reference to an address. An external IP address associated with this instance. Specify an unused static external IP address available to the project or leave this field undefined to use an IP from a shared ephemeral IP address pool. If you specify a static external IP address, it must live in the same region as the zone of the instance.",
														Immutable:     true,
														ServerDefault: true,
														ResourceReferences: []*dcl.PropertyResourceReference{
															&dcl.PropertyResourceReference{
																Resource: "Compute/Address",
																Field:    "selfLink",
															},
														},
													},
													"networkTier": &dcl.Property{
														Type:          "string",
														GoName:        "NetworkTier",
														GoType:        "InstanceNetworkInterfacesAccessConfigsNetworkTierEnum",
														Description:   "This signifies the networking tier used for configuring this access configuration and can only take the following values: PREMIUM, STANDARD. If an AccessConfig is specified without a valid external IP address, an ephemeral IP will be created with this networkTier. If an AccessConfig with a valid external IP address is specified, it must match that of the networkTier associated with the Address resource owning that IP.",
														Immutable:     true,
														ServerDefault: true,
														Enum: []string{
															"PREMIUM",
															"STANDARD",
														},
													},
													"publicPtrDomainName": &dcl.Property{
														Type:        "string",
														GoName:      "PublicPtrDomainName",
														Description: "The DNS domain name for the public PTR record. You can set this field only if the setPublicPtr field is enabled.",
														Immutable:   true,
													},
													"setPublicPtr": &dcl.Property{
														Type:        "boolean",
														GoName:      "SetPublicPtr",
														Description: "Specifies whether a public DNS 'PTR' record should be created to map the external IP address of the instance to a DNS domain name.",
														Immutable:   true,
													},
													"type": &dcl.Property{
														Type:        "string",
														GoName:      "Type",
														GoType:      "InstanceNetworkInterfacesAccessConfigsTypeEnum",
														Description: "The type of configuration. The default and only option is ONE_TO_ONE_NAT.",
														Immutable:   true,
														Enum: []string{
															"ONE_TO_ONE_NAT",
														},
													},
												},
											},
										},
										"aliasIPRanges": &dcl.Property{
											Type:        "array",
											GoName:      "AliasIPRanges",
											Description: "An array of alias IP ranges for this network interface. Can only be specified for network interfaces on subnet-mode networks.",
											Immutable:   true,
											SendEmpty:   true,
											ListType:    "list",
											Items: &dcl.Property{
												Type:   "object",
												GoType: "InstanceNetworkInterfacesAliasIPRanges",
												Properties: map[string]*dcl.Property{
													"ipCidrRange": &dcl.Property{
														Type:        "string",
														GoName:      "IPCidrRange",
														Description: "The IP CIDR range represented by this alias IP range. This IP CIDR range must belong to the specified subnetwork and cannot contain IP addresses reserved by system or used by other network interfaces. This range may be a single IP address (e.g. 10.2.3.4), a netmask (e.g. /24) or a CIDR format string (e.g. 10.1.2.0/24).",
														Immutable:   true,
													},
													"subnetworkRangeName": &dcl.Property{
														Type:        "string",
														GoName:      "SubnetworkRangeName",
														Description: "Optional subnetwork secondary range name specifying the secondary range from which to allocate the IP CIDR range for this alias IP range. If left unspecified, the primary range of the subnetwork will be used.",
														Immutable:   true,
													},
												},
											},
										},
										"ipv6AccessConfigs": &dcl.Property{
											Type:        "array",
											GoName:      "IPv6AccessConfigs",
											Description: "An array of IPv6 access configurations for this interface. Currently, only one IPv6 access config, DIRECT_IPV6, is supported. If there is no ipv6AccessConfig specified, then this instance will have no external IPv6 Internet access.",
											Immutable:   true,
											SendEmpty:   true,
											ListType:    "list",
											Items: &dcl.Property{
												Type:   "object",
												GoType: "InstanceNetworkInterfacesIPv6AccessConfigs",
												Required: []string{
													"name",
													"type",
												},
												Properties: map[string]*dcl.Property{
													"externalIPv6": &dcl.Property{
														Type:        "string",
														GoName:      "ExternalIPv6",
														ReadOnly:    true,
														Description: "The first IPv6 address of the external IPv6 range associated with this instance, prefix length is stored in externalIpv6PrefixLength in ipv6AccessConfig. The field is output only, an IPv6 address from a subnetwork associated with the instance will be allocated dynamically.",
														Immutable:   true,
													},
													"externalIPv6PrefixLength": &dcl.Property{
														Type:        "string",
														GoName:      "ExternalIPv6PrefixLength",
														ReadOnly:    true,
														Description: "The prefix length of the external IPv6 range.",
														Immutable:   true,
													},
													"name": &dcl.Property{
														Type:        "string",
														GoName:      "Name",
														Description: "The name of this access configuration. The default and recommended name is External NAT but you can use any arbitrary string you would like. For example, My external IP or Network Access.",
														Immutable:   true,
													},
													"natIP": &dcl.Property{
														Type:          "string",
														GoName:        "NatIP",
														Description:   "Reference to an address. An external IP address associated with this instance. Specify an unused static external IP address available to the project or leave this field undefined to use an IP from a shared ephemeral IP address pool. If you specify a static external IP address, it must live in the same region as the zone of the instance.",
														Immutable:     true,
														ServerDefault: true,
														ResourceReferences: []*dcl.PropertyResourceReference{
															&dcl.PropertyResourceReference{
																Resource: "Compute/Address",
																Field:    "selfLink",
															},
														},
													},
													"networkTier": &dcl.Property{
														Type:        "string",
														GoName:      "NetworkTier",
														GoType:      "InstanceNetworkInterfacesIPv6AccessConfigsNetworkTierEnum",
														Description: "This signifies the networking tier used for configuring this access configuration and can only take the following values: PREMIUM, STANDARD. If an AccessConfig is specified without a valid external IP address, an ephemeral IP will be created with this networkTier. If an AccessConfig with a valid external IP address is specified, it must match that of the networkTier associated with the Address resource owning that IP.",
														Immutable:   true,
														Enum: []string{
															"PREMIUM",
															"STANDARD",
														},
													},
													"publicPtrDomainName": &dcl.Property{
														Type:        "string",
														GoName:      "PublicPtrDomainName",
														Description: "The DNS domain name for the public PTR record. You can set this field only if the setPublicPtr field is enabled.",
														Immutable:   true,
													},
													"setPublicPtr": &dcl.Property{
														Type:        "boolean",
														GoName:      "SetPublicPtr",
														Description: "Specifies whether a public DNS 'PTR' record should be created to map the external IP address of the instance to a DNS domain name.",
														Immutable:   true,
													},
													"type": &dcl.Property{
														Type:        "string",
														GoName:      "Type",
														GoType:      "InstanceNetworkInterfacesIPv6AccessConfigsTypeEnum",
														Description: "The type of configuration. The default and only option is ONE_TO_ONE_NAT.",
														Immutable:   true,
														Enum: []string{
															"ONE_TO_ONE_NAT",
														},
													},
												},
											},
										},
										"name": &dcl.Property{
											Type:        "string",
											GoName:      "Name",
											ReadOnly:    true,
											Description: "The name of the network interface, generated by the server. For network devices, these are eth0, eth1, etc",
											Immutable:   true,
										},
										"network": &dcl.Property{
											Type:          "string",
											GoName:        "Network",
											Description:   "Specifies the title of an existing network.  When creating an instance, if neither the network nor the subnetwork is specified, the default network global/networks/default is used; if the network is not specified but the subnetwork is specified, the network is inferred.",
											Immutable:     true,
											ServerDefault: true,
											ResourceReferences: []*dcl.PropertyResourceReference{
												&dcl.PropertyResourceReference{
													Resource: "Compute/Network",
													Field:    "name",
												},
											},
										},
										"networkIP": &dcl.Property{
											Type:          "string",
											GoName:        "NetworkIP",
											Description:   "An IPv4 internal network address to assign to the instance for this network interface. If not specified by the user, an unused internal IP is assigned by the system.",
											Immutable:     true,
											ServerDefault: true,
										},
										"subnetwork": &dcl.Property{
											Type:          "string",
											GoName:        "Subnetwork",
											Description:   "Reference to a VPC network. If the network resource is in legacy mode, do not provide this property.  If the network is in auto subnet mode, providing the subnetwork is optional. If the network is in custom subnet mode, then this field should be specified.",
											Immutable:     true,
											ServerDefault: true,
											ResourceReferences: []*dcl.PropertyResourceReference{
												&dcl.PropertyResourceReference{
													Resource: "Compute/Subnetwork",
													Field:    "name",
												},
											},
										},
									},
								},
							},
							"project": &dcl.Property{
								Type:        "string",
								GoName:      "Project",
								Description: "The project id of the resource.",
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
							"scheduling": &dcl.Property{
								Type:          "object",
								GoName:        "Scheduling",
								GoType:        "InstanceScheduling",
								Description:   "Sets the scheduling options for this instance.",
								Immutable:     true,
								ServerDefault: true,
								Properties: map[string]*dcl.Property{
									"automaticRestart": &dcl.Property{
										Type:        "boolean",
										GoName:      "AutomaticRestart",
										Description: "Specifies whether the instance should be automatically restarted if it is terminated by Compute Engine (not terminated by a user). You can only set the automatic restart option for standard instances. Preemptible instances cannot be automatically restarted.",
										Immutable:   true,
									},
									"onHostMaintenance": &dcl.Property{
										Type:        "string",
										GoName:      "OnHostMaintenance",
										Description: "Defines the maintenance behavior for this instance. For standard instances, the default behavior is MIGRATE. For preemptible instances, the default and only possible behavior is TERMINATE. For more information, see Setting Instance Scheduling Options.",
										Immutable:   true,
									},
									"preemptible": &dcl.Property{
										Type:        "boolean",
										GoName:      "Preemptible",
										Description: "Defines whether the instance is preemptible. This can only be set during instance creation, it cannot be set or changed after the instance has been created.",
										Immutable:   true,
									},
								},
							},
							"selfLink": &dcl.Property{
								Type:        "string",
								GoName:      "SelfLink",
								ReadOnly:    true,
								Description: "The self link of the instance",
								Immutable:   true,
							},
							"serviceAccounts": &dcl.Property{
								Type:        "array",
								GoName:      "ServiceAccounts",
								Description: "A list of service accounts, with their specified scopes, authorized for this instance. Only one service account per VM instance is supported.",
								Immutable:   true,
								SendEmpty:   true,
								ListType:    "list",
								Items: &dcl.Property{
									Type:   "object",
									GoType: "InstanceServiceAccounts",
									Properties: map[string]*dcl.Property{
										"email": &dcl.Property{
											Type:        "string",
											GoName:      "Email",
											Description: "Email address of the service account.",
											Immutable:   true,
										},
										"scopes": &dcl.Property{
											Type:        "array",
											GoName:      "Scopes",
											Description: "The list of scopes to be made available for this service account.",
											Immutable:   true,
											SendEmpty:   true,
											ListType:    "list",
											Items: &dcl.Property{
												Type:   "string",
												GoType: "string",
											},
										},
									},
								},
							},
							"shieldedInstanceConfig": &dcl.Property{
								Type:          "object",
								GoName:        "ShieldedInstanceConfig",
								GoType:        "InstanceShieldedInstanceConfig",
								Description:   "Configuration for various parameters related to shielded instances.",
								ServerDefault: true,
								Properties: map[string]*dcl.Property{
									"enableIntegrityMonitoring": &dcl.Property{
										Type:        "boolean",
										GoName:      "EnableIntegrityMonitoring",
										Description: "Defines whether the instance has integrity monitoring enabled.",
									},
									"enableSecureBoot": &dcl.Property{
										Type:        "boolean",
										GoName:      "EnableSecureBoot",
										Description: "Defines whether the instance has Secure Boot enabled.",
									},
									"enableVtpm": &dcl.Property{
										Type:        "boolean",
										GoName:      "EnableVtpm",
										Description: "Defines whether the instance has the vTPM enabled",
									},
								},
							},
							"status": &dcl.Property{
								Type:          "string",
								GoName:        "Status",
								GoType:        "InstanceStatusEnum",
								Description:   "The status of the instance. One of the following values: PROVISIONING, STAGING, RUNNING, STOPPING, SUSPENDING, SUSPENDED, and TERMINATED.  As a user, use RUNNING to keep a machine \"on\" and TERMINATED to turn a machine off",
								Immutable:     true,
								ServerDefault: true,
								Enum: []string{
									"PROVISIONING",
									"STAGING",
									"RUNNING",
									"STOPPING",
									"SUSPENDING",
									"SUSPENDED",
									"TERMINATED",
								},
							},
							"statusMessage": &dcl.Property{
								Type:        "string",
								GoName:      "StatusMessage",
								ReadOnly:    true,
								Description: "An optional, human-readable explanation of the status.",
								Immutable:   true,
							},
							"tags": &dcl.Property{
								Type:        "array",
								GoName:      "Tags",
								Description: "A list of tags to apply to this instance. Tags are used to identify valid sources or targets for network firewalls and are specified by the client during instance creation. Each tag within the list must comply with RFC1035.",
								SendEmpty:   true,
								ListType:    "list",
								Items: &dcl.Property{
									Type:   "string",
									GoType: "string",
								},
							},
							"zone": &dcl.Property{
								Type:        "string",
								GoName:      "Zone",
								Description: "A reference to the zone where the machine resides.",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Compute/Zone",
										Field:    "name",
										Parent:   true,
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
