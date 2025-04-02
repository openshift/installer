// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isInstances         = "instances"
	isInstanceGroupName = "instance_group_name"
)

func DataSourceIBMISInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISInstancesRead,

		Schema: map[string]*schema.Schema{
			isInstanceGroup: {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"vpc", "vpc_crn", "vpc_name", isInstanceGroupName},
				Description:   "Instance group ID to filter the instances attached to it",
			},
			isInstanceGroupName: {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"vpc", "vpc_crn", "vpc_name", isInstanceGroup},
				Description:   "Instance group name to filter the instances attached to it",
			},
			"vpc_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"vpc", "vpc_crn", "instance_group"},
				Description:   "Name of the vpc to filter the instances attached to it",
			},

			"vpc": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"vpc_name", "vpc_crn", "instance_group"},
				Description:   "VPC ID to filter the instances attached to it",
			},

			"vpc_crn": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"vpc_name", "vpc", "instance_group"},
				Description:   "VPC CRN to filter the instances attached to it",
			},

			"resource_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Instance resource group",
			},

			// cluster changes
			"cluster_network_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters the collection to instances with a `cluster_network.id` property matching the specified identifier.",
			},
			"cluster_network_crn": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters the collection to instances with a `cluster_network.crn` property matching the specified CRN.",
			},
			"cluster_network_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters the collection to resources with a `cluster_network.name` property matching the exact specified name.",
			},

			"dedicated_host_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"dedicated_host"},
				Description:   "Name of the dedicated host to filter the instances attached to it",
			},

			"dedicated_host": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"dedicated_host_name"},
				Description:   "ID of the dedicated host to filter the instances attached to it",
			},

			"placement_group_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"placement_group"},
				Description:   "Name of the placement group to filter the instances attached to it",
			},

			"placement_group": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"placement_group_name"},
				Description:   "ID of the placement group to filter the instances attached to it",
			},

			isInstances: {
				Type:        schema.TypeList,
				Description: "List of instances",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance id",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance name",
						},
						"crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The crn for this Instance",
						},
						// cluster changes
						"cluster_network": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, the cluster network that this virtual server instance resides in.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this cluster network.",
									},
									"deleted": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"more_info": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
									},
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this cluster network.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this cluster network.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this cluster network. The name must not be used by another cluster network in the region.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
						"cluster_network_attachments": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The cluster network attachments for this virtual server instance.The cluster network attachments are ordered for consistent instance configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this instance cluster network attachment.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this instance cluster network attachment.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this instance cluster network attachment. The name is unique across all network attachments for the instance.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
						"confidential_compute_mode": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The confidential compute mode to use for this virtual server instance.If unspecified, the default confidential compute mode from the profile will be used.",
						},
						"enable_secure_boot": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether secure boot is enabled for this virtual server instance.If unspecified, the default secure boot mode from the profile will be used.",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance memory",
						},
						"numa_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of NUMA nodes this virtual server instance is provisioned on. This property may be absent if the instance's `status` is not `running`.",
						},
						isInstanceMetadataServiceEnabled: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether the metadata service endpoint is available to the virtual server instance",
						},
						isInstanceMetadataService: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The metadata service configuration",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isInstanceMetadataServiceEnabled1: {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicates whether the metadata service endpoint will be available to the virtual server instance",
									},

									isInstanceMetadataServiceProtocol: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The communication protocol to use for the metadata service endpoint. Applies only when the metadata service is enabled.",
									},

									isInstanceMetadataServiceRespHopLimit: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The hop limit (IP time to live) for IP response packets from the metadata service",
									},
								},
							},
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance status",
						},
						isInstanceAvailablePolicyHostFailure: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The availability policy to use for this virtual server instance. The action to perform if the compute host experiences a failure.",
						},

						isInstanceLifecycleState: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The lifecycle state of the virtual server instance.",
						},
						isInstanceLifecycleReasons: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The reasons for the current lifecycle_state (if any).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isInstanceLifecycleReasonsCode: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A snake case string succinctly identifying the reason for this lifecycle state.",
									},

									isInstanceLifecycleReasonsMessage: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "An explanation of the reason for this lifecycle state.",
									},

									isInstanceLifecycleReasonsMoreInfo: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about the reason for this lifecycle state.",
									},
								},
							},
						},

						isInstanceStatusReasons: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The reasons for the current status (if any).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isInstanceStatusReasonsCode: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A snake case string succinctly identifying the status reason",
									},

									isInstanceStatusReasonsMessage: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "An explanation of the status reason",
									},

									isInstanceStatusReasonsMoreInfo: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about this status reason",
									},
								},
							},
						},

						"resource_group": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance resource group",
						},
						"vpc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "vpc attached to the instance",
						},

						isInstanceCatalogOffering: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The catalog offering or offering version to use when provisioning this virtual server instance. If an offering is specified, the latest version of that offering will be used. The specified offering or offering version may be in a different account in the same enterprise, subject to IAM policies.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isInstanceCatalogOfferingOfferingCrn: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Identifies a catalog offering by a unique CRN property",
									},
									isInstanceCatalogOfferingVersionCrn: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Identifies a version of a catalog offering by a unique CRN property",
									},
									isInstanceCatalogOfferingPlanCrn: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this catalog offering version's billing plan",
									},
									"deleted": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted and provides some supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"more_info": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
									},
								},
							},
						},
						isInstanceTags: {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         flex.ResourceIBMVPCHash,
							Description: "list of tags for the instance",
						},

						isInstanceAccessTags: {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         flex.ResourceIBMVPCHash,
							Description: "list of access tags for the instance",
						},
						"boot_volume": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Instance Boot Volume",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance Boot volume id",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance Boot volume name",
									},
									"device": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance Boot volume device",
									},
									"volume_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance Boot volume's volume id",
									},
									"volume_crn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance Boot volume's volume CRN",
									},
								},
							},
						},

						"volume_attachments": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Instance Volume Attachments",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance volume Attachment id",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance volume Attachment name",
									},
									"volume_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance volume Attachment's volume id",
									},
									"volume_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance volume Attachment's volume name",
									},
									"volume_crn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance volume Attachment's volume CRN",
									},
								},
							},
						},

						"primary_network_interface": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Instance Primary Network Interface",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance Primary Network interface id",
									},
									isInstanceNicName: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance Primary Network interface name",
									},
									isInstanceNicPrimaryIpv4Address: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance Primary Network interface IPV4 Address",
									},
									isInstanceNicPrimaryIP: {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The primary IP address to bind to the network interface. This can be specified using an existing reserved IP, or a prototype object for a new reserved IP.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												isInstanceNicReservedIpAddress: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The IP address to reserve, which must not already be reserved on the subnet.",
												},
												isInstanceNicReservedIpHref: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The URL for this reserved IP",
												},
												isInstanceNicReservedIpName: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The user-defined name for this reserved IP. If unspecified, the name will be a hyphenated list of randomly-selected words. Names must be unique within the subnet the reserved IP resides in. ",
												},
												isInstanceNicReservedIpId: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Identifies a reserved IP by a unique property.",
												},
												isInstanceNicReservedIpResourceType: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The resource type",
												},
											},
										},
									},
									isInstanceNicSecurityGroups: {
										Type:        schema.TypeSet,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Set:         schema.HashString,
										Description: "Instance Primary Network interface security groups",
									},
									isInstanceNicSubnet: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance Primary Network interface subnet",
									},
								},
							},
						},
						"primary_network_attachment": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The primary network attachment for this virtual server instance.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"deleted": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"more_info": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
									},
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this network attachment.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this network attachment.",
									},
									"name": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"primary_ip": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The primary IP address of the virtual network interface for the network attachment.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"address": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The IP address.If the address has not yet been selected, the value will be `0.0.0.0`.This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.",
												},
												"deleted": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"more_info": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Link to documentation about deleted resources.",
															},
														},
													},
												},
												"href": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The URL for this reserved IP.",
												},
												"id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The unique identifier for this reserved IP.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name for this reserved IP. The name is unique across all reserved IPs in a subnet.",
												},
												"resource_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The resource type.",
												},
											},
										},
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
									"subnet": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The subnet of the virtual network interface for the network attachment.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"crn": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The CRN for this subnet.",
												},
												"deleted": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"more_info": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Link to documentation about deleted resources.",
															},
														},
													},
												},
												"href": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The URL for this subnet.",
												},
												"id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The unique identifier for this subnet.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name for this subnet. The name is unique across all subnets in the VPC.",
												},
												"resource_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The resource type.",
												},
											},
										},
									},
									"virtual_network_interface": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The virtual network interface for this instance network attachment.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"crn": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The CRN for this virtual network interface.",
												},
												"href": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The URL for this virtual network interface.",
												},
												"id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The unique identifier for this virtual network interface.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name for this virtual network interface. The name is unique across all virtual network interfaces in the VPC.",
												},
												"resource_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The resource type.",
												},
											},
										},
									},
								},
							},
						},

						"placement_target": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The placement restrictions for the virtual server instance.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this placement target resource.",
									},
									"deleted": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"more_info": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
									},
									"href": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this placement target resource.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this placement target resource.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique user-defined name for this placement target resource. If unspecified, the name will be a hyphenated list of randomly-selected words.",
									},
									"resource_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of resource referenced.",
									},
								},
							},
						},
						"network_interfaces": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Instance Network Interfaces",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance Network interface id",
									},
									isInstanceNicName: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance Network interface name",
									},
									isInstanceNicPrimaryIpv4Address: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance Network interface IPV4 Address",
									},
									isInstanceNicPrimaryIP: {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The primary IP address to bind to the network interface. This can be specified using an existing reserved IP, or a prototype object for a new reserved IP.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												isInstanceNicReservedIpAddress: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The IP address to reserve, which must not already be reserved on the subnet.",
												},
												isInstanceNicReservedIpHref: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The URL for this reserved IP",
												},
												isInstanceNicReservedIpName: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The user-defined name for this reserved IP. If unspecified, the name will be a hyphenated list of randomly-selected words. Names must be unique within the subnet the reserved IP resides in. ",
												},
												isInstanceNicReservedIpId: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Identifies a reserved IP by a unique property.",
												},
												isInstanceNicReservedIpResourceType: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The resource type",
												},
											},
										},
									},
									isInstanceNicSecurityGroups: {
										Type:        schema.TypeSet,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Set:         schema.HashString,
										Description: "Instance Network interface security groups",
									},
									isInstanceNicSubnet: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance Network interface subnet",
									},
								},
							},
						},

						"network_attachments": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The network attachments for this virtual server instance, including the primary network attachment.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"deleted": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"more_info": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
									},
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this network attachment.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this network attachment.",
									},
									"name": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"primary_ip": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The primary IP address of the virtual network interface for the network attachment.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"address": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The IP address.If the address has not yet been selected, the value will be `0.0.0.0`.This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.",
												},
												"deleted": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"more_info": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Link to documentation about deleted resources.",
															},
														},
													},
												},
												"href": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The URL for this reserved IP.",
												},
												"id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The unique identifier for this reserved IP.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name for this reserved IP. The name is unique across all reserved IPs in a subnet.",
												},
												"resource_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The resource type.",
												},
											},
										},
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
									"subnet": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The subnet of the virtual network interface for the network attachment.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"crn": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The CRN for this subnet.",
												},
												"deleted": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"more_info": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Link to documentation about deleted resources.",
															},
														},
													},
												},
												"href": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The URL for this subnet.",
												},
												"id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The unique identifier for this subnet.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name for this subnet. The name is unique across all subnets in the VPC.",
												},
												"resource_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The resource type.",
												},
											},
										},
									},
									"virtual_network_interface": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The virtual network interface for this instance network attachment.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"crn": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The CRN for this virtual network interface.",
												},
												"href": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The URL for this virtual network interface.",
												},
												"id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The unique identifier for this virtual network interface.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name for this virtual network interface. The name is unique across all virtual network interfaces in the VPC.",
												},
												"resource_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The resource type.",
												},
											},
										},
									},
								},
							},
						},

						"profile": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Profile",
						},
						isInstanceTotalVolumeBandwidth: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The amount of bandwidth (in megabits per second) allocated exclusively to instance storage volumes",
						},

						isInstanceBandwidth: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total bandwidth (in megabits per second) shared across the instance's network interfaces and storage volumes",
						},

						isInstanceTotalNetworkBandwidth: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The amount of bandwidth (in megabits per second) allocated exclusively to instance network interfaces.",
						},
						"vcpu": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Instance vcpu",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"architecture": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance vcpu architecture",
									},
									"count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Instance vcpu count",
									},
									"manufacturer": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance vcpu manufacturer",
									},
								},
							},
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance zone",
						},
						"image": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Instance Image",
						},

						isInstanceGpu: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Instance GPU",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isInstanceGpuCount: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Instance GPU Count",
									},
									isInstanceGpuMemory: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Instance GPU Memory",
									},
									isInstanceGpuManufacturer: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance GPU Manufacturer",
									},
									isInstanceGpuModel: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance GPU Model",
									},
								},
							},
						},

						isInstanceDisks: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Collection of the instance's disks.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"created_at": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The date and time that the disk was created.",
									},
									"href": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this instance disk.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this instance disk.",
									},
									"interface_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The disk interface used for attaching the disk.The enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected property value was encountered.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this disk.",
									},
									"resource_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
									"size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The size of the disk in GB (gigabytes).",
									},
								},
							},
						},
						"health_reasons": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The reasons for the current health_state (if any).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A snake case string succinctly identifying the reason for this health state.",
									},
									"message": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "An explanation of the reason for this health state.",
									},
									"more_info": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about the reason for this health state.",
									},
								},
							},
						},
						"health_state": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The health of this resource",
						},
						isInstanceReservation: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The reservation used by this virtual server instance",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isReservationId: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this reservation.",
									},
									isReservationCrn: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this reservation.",
									},
									isReservationName: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this reservation. The name is unique across all reservations in the region.",
									},
									isReservationHref: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this reservation.",
									},
									isReservationResourceType: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
									isReservationDeleted: &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												isReservationDeletedMoreInfo: &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
									},
								},
							},
						},
						isReservationAffinity: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isReservationAffinityPolicyResp: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The reservation affinity policy to use for this virtual server instance.",
									},
									isReservationAffinityPool: &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The pool of reservations available for use by this virtual server instance.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												isReservationId: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The unique identifier for this reservation.",
												},
												isReservationCrn: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The CRN for this reservation.",
												},
												isReservationName: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name for this reservation. The name is unique across all reservations in the region.",
												},
												isReservationHref: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The URL for this reservation.",
												},
												isReservationResourceType: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The resource type.",
												},
												isReservationDeleted: &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															isReservationDeletedMoreInfo: &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Link to documentation about deleted resources.",
															},
														},
													},
												},
											},
										},
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

func dataSourceIBMISInstancesRead(d *schema.ResourceData, meta interface{}) error {

	err := instancesList(d, meta)
	if err != nil {
		return err
	}
	return nil
}

func instancesList(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	var vpcName, vpcID, vpcCrn, resourceGroup, insGrp, dHostNameStr, dHostIdStr, placementGrpNameStr, placementGrpIdStr string

	if vpc, ok := d.GetOk("vpc_name"); ok {
		vpcName = vpc.(string)
	}

	if vpc, ok := d.GetOk("vpc"); ok {
		vpcID = vpc.(string)
	}

	if vpccrn, ok := d.GetOk("vpc_crn"); ok {
		vpcCrn = vpccrn.(string)
	}

	if rg, ok := d.GetOk("resource_group"); ok {
		resourceGroup = rg.(string)
	}

	if dHostNameIntf, ok := d.GetOk("dedicated_host_name"); ok {
		dHostNameStr = dHostNameIntf.(string)
	}

	if dHostIdIntf, ok := d.GetOk("dedicated_host"); ok {
		dHostIdStr = dHostIdIntf.(string)
	}

	if placementGrpNameIntf, ok := d.GetOk("placement_group_name"); ok {
		placementGrpNameStr = placementGrpNameIntf.(string)
	}

	if placementGrpIdIntf, ok := d.GetOk("placement_group"); ok {
		placementGrpIdStr = placementGrpIdIntf.(string)
	}

	if insGrpInf, ok := d.GetOk(isInstanceGroup); ok {
		insGrp = insGrpInf.(string)
	} else if insGrpNameInf, ok := d.GetOk(isInstanceGroupName); ok {
		insGrpName := insGrpNameInf.(string)
		start := ""
		allrecs := []vpcv1.InstanceGroup{}
		for {
			listInstanceGroupOptions := vpcv1.ListInstanceGroupsOptions{}
			if start != "" {
				listInstanceGroupOptions.Start = &start
			}
			instanceGroupsCollection, response, err := sess.ListInstanceGroups(&listInstanceGroupOptions)
			if err != nil {
				return fmt.Errorf("[ERROR] Error Fetching InstanceGroups %s\n%s", err, response)
			}
			start = flex.GetNext(instanceGroupsCollection.Next)
			allrecs = append(allrecs, instanceGroupsCollection.InstanceGroups...)

			if start == "" {
				break
			}

		}

		for _, instanceGroup := range allrecs {
			if *instanceGroup.Name == insGrpName {
				insGrp = *instanceGroup.ID
				break
			}
		}
	}

	listInstancesOptions := &vpcv1.ListInstancesOptions{}

	// cluster changes
	if _, ok := d.GetOk("cluster_network_id"); ok {
		listInstancesOptions.SetClusterNetworkID(d.Get("cluster_network_id").(string))
	}
	if _, ok := d.GetOk("cluster_network_crn"); ok {
		listInstancesOptions.SetClusterNetworkCRN(d.Get("cluster_network_crn").(string))
	}
	if _, ok := d.GetOk("cluster_network_name"); ok {
		listInstancesOptions.SetClusterNetworkName(d.Get("cluster_network_name").(string))
	}

	if vpcName != "" {
		listInstancesOptions.VPCName = &vpcName
	}
	if vpcID != "" {
		listInstancesOptions.VPCID = &vpcID
	}
	if resourceGroup != "" {
		listInstancesOptions.ResourceGroupID = &resourceGroup
	}
	if vpcCrn != "" {
		listInstancesOptions.VPCCRN = &vpcCrn
	}

	if dHostNameStr != "" {
		listInstancesOptions.DedicatedHostName = &dHostNameStr
	}

	if dHostIdStr != "" {
		listInstancesOptions.DedicatedHostID = &dHostIdStr
	}

	if placementGrpNameStr != "" {
		listInstancesOptions.PlacementGroupName = &placementGrpNameStr
	}

	if placementGrpIdStr != "" {
		listInstancesOptions.PlacementGroupID = &placementGrpIdStr
	}

	start := ""
	allrecs := []vpcv1.Instance{}
	for {

		if start != "" {
			listInstancesOptions.Start = &start
		}

		instances, response, err := sess.ListInstances(listInstancesOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Fetching Instances %s\n%s", err, response)
		}
		start = flex.GetNext(instances.Next)
		allrecs = append(allrecs, instances.Instances...)
		if start == "" {
			break
		}
	}

	if insGrp != "" {
		membershipMap := map[string]bool{}
		start := ""
		for {
			listInstanceGroupMembershipsOptions := vpcv1.ListInstanceGroupMembershipsOptions{
				InstanceGroupID: &insGrp,
			}
			if start != "" {
				listInstanceGroupMembershipsOptions.Start = &start
			}
			instanceGroupMembershipCollection, response, err := sess.ListInstanceGroupMemberships(&listInstanceGroupMembershipsOptions)
			if err != nil {
				return fmt.Errorf("[ERROR] Error Getting InstanceGroup Membership Collection %s\n%s", err, response)
			}

			start = flex.GetNext(instanceGroupMembershipCollection.Next)
			for _, membershipItem := range instanceGroupMembershipCollection.Memberships {
				membershipMap[*membershipItem.Instance.ID] = true
			}

			if start == "" {
				break
			}

		}

		//Filtering instance allrecs to contain instance group members only
		i := 0
		for _, ins := range allrecs {
			if membershipMap[*ins.ID] {
				allrecs[i] = ins
				i++
			}
		}
		allrecs = allrecs[:i]
	}

	instancesInfo := make([]map[string]interface{}, 0)
	for _, instance := range allrecs {
		id := *instance.ID
		l := map[string]interface{}{}
		l["id"] = id
		l["crn"] = *instance.CRN
		l["name"] = *instance.Name
		l["memory"] = *instance.Memory
		if instance.NumaCount != nil {
			l["numa_count"] = *instance.NumaCount
		}
		l["confidential_compute_mode"] = instance.ConfidentialComputeMode

		l["enable_secure_boot"] = instance.EnableSecureBoot
		if instance.MetadataService != nil {
			l[isInstanceMetadataServiceEnabled] = *instance.MetadataService.Enabled
			metadataService := []map[string]interface{}{}
			metadataServiceMap := map[string]interface{}{}

			metadataServiceMap[isInstanceMetadataServiceEnabled1] = instance.MetadataService.Enabled
			if instance.MetadataService.Protocol != nil {
				metadataServiceMap[isInstanceMetadataServiceProtocol] = instance.MetadataService.Protocol
			}
			if instance.MetadataService.ResponseHopLimit != nil {
				metadataServiceMap[isInstanceMetadataServiceRespHopLimit] = instance.MetadataService.ResponseHopLimit
			}
			metadataService = append(metadataService, metadataServiceMap)
			l[isInstanceMetadataService] = metadataService
		}
		l["status"] = *instance.Status
		l["resource_group"] = *instance.ResourceGroup.ID
		l["vpc"] = *instance.VPC.ID

		if instance.AvailabilityPolicy != nil && instance.AvailabilityPolicy.HostFailure != nil {
			l[isInstanceAvailablePolicyHostFailure] = *instance.AvailabilityPolicy.HostFailure
		}

		if instance.PlacementTarget != nil {
			placementTargetMap := resourceIbmIsInstanceInstancePlacementToMap(*instance.PlacementTarget.(*vpcv1.InstancePlacementTarget))
			l["placement_target"] = []map[string]interface{}{placementTargetMap}
		}
		if instance.Bandwidth != nil {
			l[isInstanceBandwidth] = int(*instance.Bandwidth)
		}

		if instance.TotalNetworkBandwidth != nil {
			l[isInstanceTotalNetworkBandwidth] = int(*instance.TotalNetworkBandwidth)
		}

		if instance.TotalVolumeBandwidth != nil {
			l[isInstanceTotalVolumeBandwidth] = int(*instance.TotalVolumeBandwidth)
		}

		// catalog
		if instance.CatalogOffering != nil {
			versionCrn := *instance.CatalogOffering.Version.CRN
			catalogList := make([]map[string]interface{}, 0)
			catalogMap := map[string]interface{}{}
			catalogMap[isInstanceCatalogOfferingVersionCrn] = versionCrn
			if instance.CatalogOffering.Plan != nil {
				if instance.CatalogOffering.Plan.CRN != nil && *instance.CatalogOffering.Plan.CRN != "" {
					catalogMap[isInstanceCatalogOfferingPlanCrn] = *instance.CatalogOffering.Plan.CRN
				}
				if instance.CatalogOffering.Plan.Deleted != nil {
					deletedMap := resourceIbmIsInstanceCatalogOfferingVersionPlanReferenceDeletedToMap(*instance.CatalogOffering.Plan.Deleted)
					catalogMap["deleted"] = []map[string]interface{}{deletedMap}
				}
			}
			catalogList = append(catalogList, catalogMap)
			l[isInstanceCatalogOffering] = catalogList
		}
		clusterNetwork := []map[string]interface{}{}
		if instance.ClusterNetwork != nil {
			clusterNetworkMap, err := DataSourceIBMIsInstancesClusterNetworkReferenceToMap(instance.ClusterNetwork)
			if err != nil {
				return err
			}
			clusterNetwork = append(clusterNetwork, clusterNetworkMap)
		}
		l["cluster_network"] = clusterNetwork
		clusterNetworkAttachments := []map[string]interface{}{}
		for _, clusterNetworkAttachmentsItem := range instance.ClusterNetworkAttachments {
			clusterNetworkAttachmentsItemMap, err := DataSourceIBMIsInstancesInstanceClusterNetworkAttachmentReferenceToMap(&clusterNetworkAttachmentsItem) // #nosec G601
			if err != nil {
				return err
			}
			clusterNetworkAttachments = append(clusterNetworkAttachments, clusterNetworkAttachmentsItemMap)
		}
		l["cluster_network_attachments"] = clusterNetworkAttachments

		if instance.BootVolumeAttachment != nil {
			bootVolList := make([]map[string]interface{}, 0)
			bootVol := map[string]interface{}{}
			bootVol["id"] = *instance.BootVolumeAttachment.ID
			bootVol["name"] = *instance.BootVolumeAttachment.Name
			if instance.BootVolumeAttachment.Device != nil {
				bootVol["device"] = *instance.BootVolumeAttachment.Device.ID
			}
			if instance.BootVolumeAttachment.Volume != nil {
				bootVol["volume_id"] = *instance.BootVolumeAttachment.Volume.ID
				bootVol["volume_crn"] = *instance.BootVolumeAttachment.Volume.CRN
			}
			bootVolList = append(bootVolList, bootVol)
			l["boot_volume"] = bootVolList
		}
		tags, err := flex.GetGlobalTagsUsingCRN(meta, *instance.CRN, "", isInstanceUserTagType)
		if err != nil {
			log.Printf(
				"Error on get of resource vpc Instance (%s) tags: %s", d.Id(), err)
		}
		l[isInstanceTags] = tags

		accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *instance.CRN, "", isInstanceAccessTagType)
		if err != nil {
			log.Printf(
				"Error on get of resource vpc Instance (%s) access tags: %s", d.Id(), err)
		}
		l[isInstanceAccessTags] = accesstags
		//set the status reasons
		statusReasonsList := make([]map[string]interface{}, 0)
		if instance.StatusReasons != nil {
			for _, sr := range instance.StatusReasons {
				currentSR := map[string]interface{}{}
				if sr.Code != nil && sr.Message != nil {
					currentSR[isInstanceStatusReasonsCode] = *sr.Code
					currentSR[isInstanceStatusReasonsMessage] = *sr.Message
					if sr.MoreInfo != nil {
						currentSR[isInstanceStatusReasonsMoreInfo] = *sr.MoreInfo
					}
					statusReasonsList = append(statusReasonsList, currentSR)
				}
			}
		}
		l[isInstanceStatusReasons] = statusReasonsList

		if instance.VolumeAttachments != nil {
			volList := make([]map[string]interface{}, 0)
			for _, volume := range instance.VolumeAttachments {
				vol := map[string]interface{}{}
				if volume.Volume != nil {
					vol["id"] = *volume.ID
					vol["volume_id"] = *volume.Volume.ID
					vol["name"] = *volume.Name
					vol["volume_name"] = *volume.Volume.Name
					vol["volume_crn"] = *volume.Volume.CRN
					volList = append(volList, vol)
				}
			}
			l["volume_attachments"] = volList
		}

		if instance.PrimaryNetworkInterface != nil {
			primaryNicList := make([]map[string]interface{}, 0)
			currentPrimNic := map[string]interface{}{}
			currentPrimNic["id"] = *instance.PrimaryNetworkInterface.ID
			currentPrimNic[isInstanceNicName] = *instance.PrimaryNetworkInterface.Name

			// reserved ip changes
			primaryIpList := make([]map[string]interface{}, 0)
			currentPrimIp := map[string]interface{}{}

			if instance.PrimaryNetworkInterface.PrimaryIP.Address != nil {
				currentPrimNic[isInstanceNicPrimaryIpv4Address] = *instance.PrimaryNetworkInterface.PrimaryIP.Address
				currentPrimIp[isInstanceNicReservedIpAddress] = *instance.PrimaryNetworkInterface.PrimaryIP.Address
			}
			if instance.PrimaryNetworkInterface.PrimaryIP.Href != nil {
				currentPrimIp[isInstanceNicReservedIpHref] = *instance.PrimaryNetworkInterface.PrimaryIP.Href
			}
			if instance.PrimaryNetworkInterface.PrimaryIP.Name != nil {
				currentPrimIp[isInstanceNicReservedIpName] = *instance.PrimaryNetworkInterface.PrimaryIP.Name
			}
			if instance.PrimaryNetworkInterface.PrimaryIP.ID != nil {
				currentPrimIp[isInstanceNicReservedIpId] = *instance.PrimaryNetworkInterface.PrimaryIP.ID
			}
			if instance.PrimaryNetworkInterface.PrimaryIP.ResourceType != nil {
				currentPrimIp[isInstanceNicReservedIpResourceType] = *instance.PrimaryNetworkInterface.PrimaryIP.ResourceType
			}
			primaryIpList = append(primaryIpList, currentPrimIp)
			currentPrimNic[isInstanceNicPrimaryIP] = primaryIpList

			getnicoptions := &vpcv1.GetInstanceNetworkInterfaceOptions{
				InstanceID: &id,
				ID:         instance.PrimaryNetworkInterface.ID,
			}
			insnic, response, err := sess.GetInstanceNetworkInterface(getnicoptions)
			if err != nil {
				return fmt.Errorf("[ERROR] Error getting network interfaces attached to the instance %s\n%s", err, response)
			}
			currentPrimNic[isInstanceNicSubnet] = *insnic.Subnet.ID
			if len(insnic.SecurityGroups) != 0 {
				secgrpList := []string{}
				for i := 0; i < len(insnic.SecurityGroups); i++ {
					secgrpList = append(secgrpList, string(*(insnic.SecurityGroups[i].ID)))
				}
				currentPrimNic[isInstanceNicSecurityGroups] = flex.NewStringSet(schema.HashString, secgrpList)
			}

			primaryNicList = append(primaryNicList, currentPrimNic)
			l["primary_network_interface"] = primaryNicList
		}

		primaryNetworkAttachment := []map[string]interface{}{}
		if instance.PrimaryNetworkAttachment != nil {
			modelMap, err := dataSourceIBMIsInstanceInstanceNetworkAttachmentReferenceToMap(instance.PrimaryNetworkAttachment)
			if err != nil {
				return err
			}
			primaryNetworkAttachment = append(primaryNetworkAttachment, modelMap)
		}
		l["primary_network_attachment"] = primaryNetworkAttachment

		if instance.NetworkInterfaces != nil {
			interfacesList := make([]map[string]interface{}, 0)
			for _, intfc := range instance.NetworkInterfaces {
				if *intfc.ID != *instance.PrimaryNetworkInterface.ID {
					currentNic := map[string]interface{}{}
					currentNic["id"] = *intfc.ID
					currentNic[isInstanceNicName] = *intfc.Name

					// reserved ip changes
					primaryIpList := make([]map[string]interface{}, 0)
					currentPrimIp := map[string]interface{}{}
					if intfc.PrimaryIP.Address != nil {
						currentPrimIp[isInstanceNicReservedIpAddress] = *intfc.PrimaryIP.Address
						currentNic[isInstanceNicPrimaryIpv4Address] = *intfc.PrimaryIP.Address
					}
					if intfc.PrimaryIP.Href != nil {
						currentPrimIp[isInstanceNicReservedIpHref] = *intfc.PrimaryIP.Href
					}
					if intfc.PrimaryIP.Name != nil {
						currentPrimIp[isInstanceNicReservedIpName] = *intfc.PrimaryIP.Name
					}
					if intfc.PrimaryIP.ID != nil {
						currentPrimIp[isInstanceNicReservedIpId] = *intfc.PrimaryIP.ID
					}
					if intfc.PrimaryIP.ResourceType != nil {
						currentPrimIp[isInstanceNicReservedIpResourceType] = *intfc.PrimaryIP.ResourceType
					}
					primaryIpList = append(primaryIpList, currentPrimIp)
					currentNic[isInstanceNicPrimaryIP] = primaryIpList

					getnicoptions := &vpcv1.GetInstanceNetworkInterfaceOptions{
						InstanceID: &id,
						ID:         intfc.ID,
					}
					insnic, response, err := sess.GetInstanceNetworkInterface(getnicoptions)
					if err != nil {
						return fmt.Errorf("[ERROR] Error getting network interfaces attached to the instance %s\n%s", err, response)
					}
					currentNic[isInstanceNicSubnet] = *insnic.Subnet.ID
					if len(insnic.SecurityGroups) != 0 {
						secgrpList := []string{}
						for i := 0; i < len(insnic.SecurityGroups); i++ {
							secgrpList = append(secgrpList, string(*(insnic.SecurityGroups[i].ID)))
						}
						currentNic[isInstanceNicSecurityGroups] = flex.NewStringSet(schema.HashString, secgrpList)
					}
					interfacesList = append(interfacesList, currentNic)
				}
			}
			l["network_interfaces"] = interfacesList
		}

		networkAttachments := []map[string]interface{}{}
		if instance.NetworkAttachments != nil {
			for _, modelItem := range instance.NetworkAttachments {
				modelMap, err := dataSourceIBMIsInstanceInstanceNetworkAttachmentReferenceToMap(&modelItem)
				if err != nil {
					return err
				}
				networkAttachments = append(networkAttachments, modelMap)
			}
		}
		l["network_attachments"] = networkAttachments

		l["profile"] = *instance.Profile.Name

		cpuList := make([]map[string]interface{}, 0)
		if instance.Vcpu != nil {
			currentCPU := map[string]interface{}{}
			currentCPU["architecture"] = *instance.Vcpu.Architecture
			currentCPU["count"] = *instance.Vcpu.Count
			currentCPU["manufacturer"] = *instance.Vcpu.Manufacturer
			cpuList = append(cpuList, currentCPU)
		}
		l["vcpu"] = cpuList

		gpuList := make([]map[string]interface{}, 0)
		if instance.Gpu != nil {
			currentGpu := map[string]interface{}{}
			currentGpu[isInstanceGpuManufacturer] = instance.Gpu.Manufacturer
			currentGpu[isInstanceGpuModel] = instance.Gpu.Model
			currentGpu[isInstanceGpuCount] = instance.Gpu.Count
			currentGpu[isInstanceGpuMemory] = instance.Gpu.Memory
			gpuList = append(gpuList, currentGpu)
			l[isInstanceGpu] = gpuList
		}

		//set the lifecycle status, reasons
		if instance.LifecycleState != nil {
			l[isInstanceLifecycleState] = *instance.LifecycleState
		}
		if instance.LifecycleReasons != nil {
			l[isInstanceLifecycleReasons] = dataSourceInstanceFlattenLifecycleReasons(instance.LifecycleReasons)
		}

		l["zone"] = *instance.Zone.Name
		if instance.Image != nil {
			l["image"] = *instance.Image.ID
		}

		if instance.Disks != nil {
			l[isInstanceDisks] = dataSourceInstanceFlattenDisks(instance.Disks)
		}
		if instance.HealthReasons != nil {
			healthReasonsList := []map[string]interface{}{}
			for _, healthReasonsItem := range instance.HealthReasons {
				healthReasonsList = append(healthReasonsList, dataSourceInstancesCollectionHealthReasonsToMap(healthReasonsItem))
			}
			l["health_reasons"] = healthReasonsList
		}
		if instance.HealthState != nil {
			l["health_state"] = instance.HealthState
		}
		if instance.ReservationAffinity != nil {
			reservationAffinity := []map[string]interface{}{}
			reservationAffinityMap := map[string]interface{}{}

			reservationAffinityMap[isReservationAffinityPolicyResp] = instance.ReservationAffinity.Policy
			if instance.ReservationAffinity.Pool != nil {
				poolList := make([]map[string]interface{}, 0)
				for _, pool := range instance.ReservationAffinity.Pool {
					res := map[string]interface{}{}

					res[isReservationId] = *pool.ID
					res[isReservationHref] = *pool.Href
					res[isReservationName] = *pool.Name
					res[isReservationCrn] = *pool.CRN
					res[isReservationResourceType] = *pool.ResourceType
					if pool.Deleted != nil {
						deletedList := []map[string]interface{}{}
						deletedMap := dataSourceReservationDeletedToMap(*pool.Deleted)
						deletedList = append(deletedList, deletedMap)
						res[isReservationDeleted] = deletedList
					}
					poolList = append(poolList, res)
				}
				reservationAffinityMap[isReservationAffinityPool] = poolList
			}
			reservationAffinity = append(reservationAffinity, reservationAffinityMap)
			l[isReservationAffinity] = reservationAffinity
		}
		if instance.Reservation != nil {
			resList := make([]map[string]interface{}, 0)
			res := map[string]interface{}{}

			res[isReservationId] = *instance.Reservation.ID
			res[isReservationHref] = *instance.Reservation.Href
			res[isReservationName] = *instance.Reservation.Name
			res[isReservationCrn] = *instance.Reservation.CRN
			res[isReservationResourceType] = *instance.Reservation.ResourceType
			if instance.Reservation.Deleted != nil {
				deletedList := []map[string]interface{}{}
				deletedMap := dataSourceReservationDeletedToMap(*instance.Reservation.Deleted)
				deletedList = append(deletedList, deletedMap)
				res[isReservationDeleted] = deletedList
			}
			resList = append(resList, res)
			l[isInstanceReservation] = resList
		}

		instancesInfo = append(instancesInfo, l)
	}
	d.SetId(dataSourceIBMISInstancesID(d))
	d.Set(isInstances, instancesInfo)
	return nil
}

// dataSourceIBMISInstancesID returns a reasonable ID for a Instance list.
func dataSourceIBMISInstancesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIBMIsInstancesClusterNetworkReferenceToMap(model *vpcv1.ClusterNetworkReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsInstancesDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func DataSourceIBMIsInstancesInstanceClusterNetworkAttachmentReferenceToMap(model *vpcv1.InstanceClusterNetworkAttachmentReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func DataSourceIBMIsInstancesDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = *model.MoreInfo
	return modelMap, nil
}

func dataSourceInstancesCollectionHealthReasonsToMap(statusReasonsItem vpcv1.InstanceHealthReason) (healthReasonsMap map[string]interface{}) {
	healthReasonsMap = map[string]interface{}{}

	if statusReasonsItem.Code != nil {
		healthReasonsMap["code"] = statusReasonsItem.Code
	}
	if statusReasonsItem.Message != nil {
		healthReasonsMap["message"] = statusReasonsItem.Message
	}
	if statusReasonsItem.MoreInfo != nil {
		healthReasonsMap["more_info"] = statusReasonsItem.MoreInfo
	}

	return healthReasonsMap
}
