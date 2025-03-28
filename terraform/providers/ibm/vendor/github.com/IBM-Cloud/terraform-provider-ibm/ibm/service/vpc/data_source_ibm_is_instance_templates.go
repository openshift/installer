// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"
	"reflect"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isInstanceTemplates                     = "templates"
	isInstanceTemplatesFirst                = "first"
	isInstanceTemplatesHref                 = "href"
	isInstanceTemplatesCrn                  = "crn"
	isInstanceTemplatesLimit                = "limit"
	isInstanceTemplatesNext                 = "next"
	isInstanceTemplatesTotalCount           = "total_count"
	isInstanceTemplatesName                 = "name"
	isInstanceTemplatesPortSpeed            = "port_speed"
	isInstanceTemplatesPortType             = "type"
	isInstanceTemplatesPortValue            = "value"
	isInstanceTemplatesDeleteVol            = "delete_volume_on_instance_delete"
	isInstanceTemplatesVol                  = "volume"
	isInstanceTemplatesMemory               = "memory"
	isInstanceTemplatesMemoryValue          = "value"
	isInstanceTemplatesMemoryType           = "type"
	isInstanceTemplatesMemoryValues         = "values"
	isInstanceTemplatesMemoryDefault        = "default"
	isInstanceTemplatesMemoryMin            = "min"
	isInstanceTemplatesMemoryMax            = "max"
	isInstanceTemplatesMemoryStep           = "step"
	isInstanceTemplatesSocketCount          = "socket_count"
	isInstanceTemplatesSocketValue          = "value"
	isInstanceTemplatesSocketType           = "type"
	isInstanceTemplatesSocketValues         = "values"
	isInstanceTemplatesSocketDefault        = "default"
	isInstanceTemplatesSocketMin            = "min"
	isInstanceTemplatesSocketMax            = "max"
	isInstanceTemplatesSocketStep           = "step"
	isInstanceTemplatesVcpuArch             = "vcpu_architecture"
	isInstanceTemplatesVcpuArchType         = "type"
	isInstanceTemplatesVcpuArchValue        = "value"
	isInstanceTemplatesVcpuCount            = "vcpu_count"
	isInstanceTemplatesVcpuCountValue       = "value"
	isInstanceTemplatesVcpuCountType        = "type"
	isInstanceTemplatesVcpuCountValues      = "values"
	isInstanceTemplatesVcpuCountDefault     = "default"
	isInstanceTemplatesVcpuCountMin         = "min"
	isInstanceTemplatesVcpuCountMax         = "max"
	isInstanceTemplatesVcpuCountStep        = "step"
	isInstanceTemplatesStart                = "start"
	isInstanceTemplatesVersion              = "version"
	isInstanceTemplatesGeneration           = "generation"
	isInstanceTemplatesBootVolumeAttachment = "boot_volume_attachment"

	isInstanceTemplateVPC                     = "vpc"
	isInstanceTemplateZone                    = "zone"
	isInstanceTemplateProfile                 = "profile"
	isInstanceTemplateKeys                    = "keys"
	isInstanceTemplateVolumeAttachments       = "volume_attachments"
	isInstanceTemplateNetworkInterfaces       = "network_interfaces"
	isInstanceTemplatePrimaryNetworkInterface = "primary_network_interface"
	isInstanceTemplateNicName                 = "name"
	isInstanceTemplateNicPortSpeed            = "port_speed"
	isInstanceTemplateNicAllowIPSpoofing      = "allow_ip_spoofing"
	isInstanceTemplateNicPrimaryIpv4Address   = "primary_ipv4_address"
	isInstanceTemplateNicSecondaryAddress     = "secondary_addresses"
	isInstanceTemplateNicSecurityGroups       = "security_groups"
	isInstanceTemplateNicSubnet               = "subnet"
	isInstanceTemplateNicFloatingIPs          = "floating_ips"
	isInstanceTemplateUserData                = "user_data"
	isInstanceTemplateGeneration              = "generation"
	isInstanceTemplateImage                   = "image"
	isInstanceTemplateResourceGroup           = "resource_group"
	isInstanceTemplateName                    = "name"
	isInstanceTemplateDeleteVolume            = "delete_volume_on_instance_delete"
	isInstanceTemplateVolAttName              = "name"
	isInstanceTemplateVolAttVolume            = "volume"
)

func DataSourceIBMISInstanceTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISInstanceTemplatesRead,
		Schema: map[string]*schema.Schema{
			isInstanceTemplates: {
				Type:        schema.TypeList,
				Description: "Collection of instance templates",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						// cluster changes
						"cluster_network_attachments": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The cluster network attachments to create for this virtual server instance. A cluster network attachment represents a device that is connected to a cluster network. The number of network attachments must match one of the values from the instance profile's `cluster_network_attachment_count` before the instance can be started.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cluster_network_interface": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "A cluster network interface for the instance cluster network attachment. This can bespecified using an existing cluster network interface that does not already have a `target`,or a prototype object for a new cluster network interface.This instance must reside in the same VPC as the specified cluster network interface. Thecluster network interface must reside in the same cluster network as the`cluster_network_interface` of any other `cluster_network_attachments` for this instance.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"auto_delete": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Indicates whether this cluster network interface will be automatically deleted when `target` is deleted.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name for this cluster network interface. The name must not be used by another interface in the cluster network. Names beginning with `ibm-` are reserved for provider-owned resources, and are not allowed. If unspecified, the name will be a hyphenated list of randomly-selected words.",
												},
												"primary_ip": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The primary IP address to bind to the cluster network interface. May be eithera cluster network subnet reserved IP identity, or a cluster network subnet reserved IPprototype object which will be used to create a new cluster network subnet reserved IP.If a cluster network subnet reserved IP identity is provided, the specified clusternetwork subnet reserved IP must be unbound.If a cluster network subnet reserved IP prototype object with an address is provided,the address must be available on the cluster network interface's cluster networksubnet. If no address is specified, an available address on the cluster network subnetwill be automatically selected and reserved.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The unique identifier for this cluster network subnet reserved IP.",
															},
															"href": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The URL for this cluster network subnet reserved IP.",
															},
															"address": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The IP address to reserve, which must not already be reserved on the subnet.If unspecified, an available address on the subnet will automatically be selected.",
															},
															"auto_delete": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Indicates whether this cluster network subnet reserved IP member will be automatically deleted when either `target` is deleted, or the cluster network subnet reserved IP is unbound.",
															},
															"name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The name for this cluster network subnet reserved IP. The name must not be used by another reserved IP in the cluster network subnet. Names starting with `ibm-` are reserved for provider-owned resources, and are not allowed. If unspecified, the name will be a hyphenated list of randomly-selected words.",
															},
														},
													},
												},
												"subnet": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The associated cluster network subnet. Required if `primary_ip` does not specify acluster network subnet reserved IP identity.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The unique identifier for this cluster network subnet.",
															},
															"href": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The URL for this cluster network subnet.",
															},
														},
													},
												},
												"id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The unique identifier for this cluster network interface.",
												},
												"href": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The URL for this cluster network interface.",
												},
											},
										},
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this cluster network attachment. Names must be unique within the instance the cluster network attachment resides in. If unspecified, the name will be a hyphenated list of randomly-selected words. Names starting with `ibm-` are reserved for provider-owned resources, and are not allowed.",
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
						isInstanceAvailablePolicyHostFailure: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The availability policy to use for this virtual server instance. The action to perform if the compute host experiences a failure.",
						},
						isInstanceTemplatesName: {
							Type:     schema.TypeString,
							Computed: true,
						},
						isInstanceTemplateMetadataServiceEnabled: {
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
						isInstanceTemplatesHref: {
							Type:     schema.TypeString,
							Computed: true,
						},
						isInstanceTemplatesCrn: {
							Type:     schema.TypeString,
							Computed: true,
						},
						isInstanceTemplateVPC: {
							Type:     schema.TypeString,
							Computed: true,
						},
						isInstanceTemplateZone: {
							Type:     schema.TypeString,
							Computed: true,
						},
						isInstanceTemplateProfile: {
							Type:     schema.TypeString,
							Computed: true,
						},
						isInstanceTemplateKeys: {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						isInstanceTotalVolumeBandwidth: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The amount of bandwidth (in megabits per second) allocated exclusively to instance storage volumes",
						},
						isInstanceDefaultTrustedProfileAutoLink: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "If set to `true`, the system will create a link to the specified `target` trusted profile during instance creation. Regardless of whether a link is created by the system or manually using the IAM Identity service, it will be automatically deleted when the instance is deleted.",
						},
						isInstanceDefaultTrustedProfileTarget: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier or CRN of the default IAM trusted profile to use for this virtual server instance.",
						},

						isInstanceTemplateVolumeAttachments: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isInstanceTemplatesDeleteVol: {
										Type:     schema.TypeBool,
										Computed: true,
									},
									isInstanceTemplatesName: {
										Type:     schema.TypeString,
										Computed: true,
									},
									isInstanceTemplatesVol: {
										Type:     schema.TypeString,
										Computed: true,
									},

									isInstanceTemplateVolAttVolPrototype: {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												isInstanceTemplateVolAttVolIops: {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The maximum I/O operations per second (IOPS) for the volume.",
												},
												isInstanceTemplateVolAttVolProfile: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The  globally unique name for the volume profile to use for this volume.",
												},
												isInstanceTemplateVolAttVolCapacity: {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The capacity of the volume in gigabytes. The specified minimum and maximum capacity values for creating or updating volumes may expand in the future.",
												},
												isInstanceTemplateVolAttVolEncryptionKey: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The CRN of the [Key Protect Root Key](https://cloud.ibm.com/docs/key-protect?topic=key-protect-getting-started-tutorial) or [Hyper Protect Crypto Service Root Key](https://cloud.ibm.com/docs/hs-crypto?topic=hs-crypto-get-started) for this resource.",
												},
												isInstanceTemplateVolAttTags: {
													Type:        schema.TypeSet,
													Computed:    true,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Set:         flex.ResourceIBMVPCHash,
													Description: "The user tags associated with this volume.",
												},
											},
										},
									},
								},
							},
						},

						isInstanceTemplateCatalogOffering: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The catalog offering or offering version to use when provisioning this virtual server instance template. If an offering is specified, the latest version of that offering will be used. The specified offering or offering version may be in a different account in the same enterprise, subject to IAM policies.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isInstanceTemplateCatalogOfferingOfferingCrn: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Identifies a catalog offering by a unique CRN property",
									},
									isInstanceTemplateCatalogOfferingVersionCrn: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Identifies a version of a catalog offering by a unique CRN property",
									},
									isInstanceTemplateCatalogOfferingPlanCrn: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this catalog offering version's billing plan",
									},
								},
							},
						},

						"placement_target": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The placement restrictions for the virtual server instance. For the target tobe changed, the instance `status` must be `stopping` or `stopped`.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this dedicated host.",
									},
									"crn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this dedicated host.",
									},
									"href": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this dedicated host.",
									},
								},
							},
						},

						"network_attachments": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The additional network attachments to create for the virtual server instance.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this network attachment. Names must be unique within the instance the network attachment resides in. If unspecified, the name will be a hyphenated list of randomly-selected words.",
									},
									"virtual_network_interface": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "A virtual network interface for the instance network attachment. This can be specifiedusing an existing virtual network interface, or a prototype object for a new virtualnetwork interface.If an existing virtual network interface is specified, `enable_infrastructure_nat` must be`false`.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"allow_ip_spoofing": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Indicates whether source IP spoofing is allowed on this interface. If `false`, source IP spoofing is prevented on this interface. If `true`, source IP spoofing is allowed on this interface.",
												},
												"auto_delete": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Indicates whether this virtual network interface will be automatically deleted when`target` is deleted.",
												},
												"enable_infrastructure_nat": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "If `true`:- The VPC infrastructure performs any needed NAT operations.- `floating_ips` must not have more than one floating IP.If `false`:- Packets are passed unchanged to/from the virtual network interface,  allowing the workload to perform any needed NAT operations.- `allow_ip_spoofing` must be `false`.- If the virtual network interface is attached:  - The target `resource_type` must be `bare_metal_server_network_attachment`.  - The target `interface_type` must not be `hipersocket`.",
												},
												"ips": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Additional IP addresses to bind to the virtual network interface. Each item may be either a reserved IP identity, or as a reserved IP prototype object which will be used to create a new reserved IP. All IP addresses must be in the same subnet as the primary IP.If reserved IP identities are provided, the specified reserved IPs must be unbound.If reserved IP prototype objects with addresses are provided, the addresses must be available on the virtual network interface's subnet. For any prototype objects that do not specify an address, an available address on the subnet will be automatically selected and reserved.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The unique identifier for this reserved IP.",
															},
															"href": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The URL for this reserved IP.",
															},
															"address": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The IP address to reserve, which must not already be reserved on the subnet.If unspecified, an available address on the subnet will automatically be selected.",
															},
															"auto_delete": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Indicates whether this reserved IP member will be automatically deleted when either`target` is deleted, or the reserved IP is unbound.",
															},
															"name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The name for this reserved IP. The name must not be used by another reserved IP in the subnet. Names starting with `ibm-` are reserved for provider-owned resources, and are not allowed. If unspecified, the name will be a hyphenated list of randomly-selected words.",
															},
														},
													},
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name for this virtual network interface. The name must not be used by another virtual network interface in the VPC. If unspecified, the name will be a hyphenated list of randomly-selected words. Names beginning with `ibm-` are reserved for provider-owned resources, and are not allowed.",
												},
												"primary_ip": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The primary IP address to bind to the virtual network interface. May be either areserved IP identity, or a reserved IP prototype object which will be used to create anew reserved IP.If a reserved IP identity is provided, the specified reserved IP must be unbound.If a reserved IP prototype object with an address is provided, the address must beavailable on the virtual network interface's subnet. If no address is specified,an available address on the subnet will be automatically selected and reserved.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The unique identifier for this reserved IP.",
															},
															"href": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The URL for this reserved IP.",
															},
															"address": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The IP address to reserve, which must not already be reserved on the subnet.If unspecified, an available address on the subnet will automatically be selected.",
															},
															"auto_delete": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Indicates whether this reserved IP member will be automatically deleted when either`target` is deleted, or the reserved IP is unbound.",
															},
															"name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The name for this reserved IP. The name must not be used by another reserved IP in the subnet. Names starting with `ibm-` are reserved for provider-owned resources, and are not allowed. If unspecified, the name will be a hyphenated list of randomly-selected words.",
															},
														},
													},
												},
												"protocol_state_filtering_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The protocol state filtering mode used for this virtual network interface.",
												},
												"resource_group": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The resource group to use for this virtual network interface. If unspecified, thevirtual server instance's resource group will be used.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The unique identifier for this resource group.",
															},
														},
													},
												},
												"security_groups": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The security groups to use for this virtual network interface. If unspecified, the default security group of the VPC for the subnet is used.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The unique identifier for this security group.",
															},
															"crn": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The security group's CRN.",
															},
															"href": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The security group's canonical URL.",
															},
														},
													},
												},
												"subnet": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The associated subnet. Required if `primary_ip` does not specify a reserved IP.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The unique identifier for this subnet.",
															},
															"crn": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The CRN for this subnet.",
															},
															"href": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The URL for this subnet.",
															},
														},
													},
												},
												"id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The unique identifier for this virtual network interface.",
												},
												"href": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The URL for this virtual network interface.",
												},
												"crn": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The CRN for this virtual network interface.",
												},
											},
										},
									},
								},
							},
						},
						"primary_network_attachment": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The primary network attachment to create for the virtual server instance.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this network attachment. Names must be unique within the instance the network attachment resides in. If unspecified, the name will be a hyphenated list of randomly-selected words.",
									},
									"virtual_network_interface": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "A virtual network interface for the instance network attachment. This can be specifiedusing an existing virtual network interface, or a prototype object for a new virtualnetwork interface.If an existing virtual network interface is specified, `enable_infrastructure_nat` must be`false`.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"allow_ip_spoofing": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Indicates whether source IP spoofing is allowed on this interface. If `false`, source IP spoofing is prevented on this interface. If `true`, source IP spoofing is allowed on this interface.",
												},
												"auto_delete": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Indicates whether this virtual network interface will be automatically deleted when`target` is deleted.",
												},
												"enable_infrastructure_nat": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "If `true`:- The VPC infrastructure performs any needed NAT operations.- `floating_ips` must not have more than one floating IP.If `false`:- Packets are passed unchanged to/from the virtual network interface,  allowing the workload to perform any needed NAT operations.- `allow_ip_spoofing` must be `false`.- If the virtual network interface is attached:  - The target `resource_type` must be `bare_metal_server_network_attachment`.  - The target `interface_type` must not be `hipersocket`.",
												},
												"ips": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Additional IP addresses to bind to the virtual network interface. Each item may be either a reserved IP identity, or as a reserved IP prototype object which will be used to create a new reserved IP. All IP addresses must be in the same subnet as the primary IP.If reserved IP identities are provided, the specified reserved IPs must be unbound.If reserved IP prototype objects with addresses are provided, the addresses must be available on the virtual network interface's subnet. For any prototype objects that do not specify an address, an available address on the subnet will be automatically selected and reserved.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The unique identifier for this reserved IP.",
															},
															"href": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The URL for this reserved IP.",
															},
															"address": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The IP address to reserve, which must not already be reserved on the subnet.If unspecified, an available address on the subnet will automatically be selected.",
															},
															"auto_delete": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Indicates whether this reserved IP member will be automatically deleted when either`target` is deleted, or the reserved IP is unbound.",
															},
															"name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The name for this reserved IP. The name must not be used by another reserved IP in the subnet. Names starting with `ibm-` are reserved for provider-owned resources, and are not allowed. If unspecified, the name will be a hyphenated list of randomly-selected words.",
															},
														},
													},
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name for this virtual network interface. The name must not be used by another virtual network interface in the VPC. If unspecified, the name will be a hyphenated list of randomly-selected words. Names beginning with `ibm-` are reserved for provider-owned resources, and are not allowed.",
												},
												"primary_ip": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The primary IP address to bind to the virtual network interface. May be either areserved IP identity, or a reserved IP prototype object which will be used to create anew reserved IP.If a reserved IP identity is provided, the specified reserved IP must be unbound.If a reserved IP prototype object with an address is provided, the address must beavailable on the virtual network interface's subnet. If no address is specified,an available address on the subnet will be automatically selected and reserved.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The unique identifier for this reserved IP.",
															},
															"href": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The URL for this reserved IP.",
															},
															"address": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The IP address to reserve, which must not already be reserved on the subnet.If unspecified, an available address on the subnet will automatically be selected.",
															},
															"auto_delete": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Indicates whether this reserved IP member will be automatically deleted when either`target` is deleted, or the reserved IP is unbound.",
															},
															"name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The name for this reserved IP. The name must not be used by another reserved IP in the subnet. Names starting with `ibm-` are reserved for provider-owned resources, and are not allowed. If unspecified, the name will be a hyphenated list of randomly-selected words.",
															},
														},
													},
												},
												"protocol_state_filtering_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The protocol state filtering mode used for this virtual network interface.",
												},
												"resource_group": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The resource group to use for this virtual network interface. If unspecified, thevirtual server instance's resource group will be used.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The unique identifier for this resource group.",
															},
														},
													},
												},
												"security_groups": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The security groups to use for this virtual network interface. If unspecified, the default security group of the VPC for the subnet is used.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The unique identifier for this security group.",
															},
															"crn": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The security group's CRN.",
															},
															"href": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The security group's canonical URL.",
															},
														},
													},
												},
												"subnet": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The associated subnet. Required if `primary_ip` does not specify a reserved IP.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The unique identifier for this subnet.",
															},
															"crn": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The CRN for this subnet.",
															},
															"href": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The URL for this subnet.",
															},
														},
													},
												},
												"id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The unique identifier for this virtual network interface.",
												},
												"href": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The URL for this virtual network interface.",
												},
												"crn": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The CRN for this virtual network interface.",
												},
											},
										},
									},
								},
							},
						},

						isInstanceTemplatePrimaryNetworkInterface: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isInstanceTemplateNicName: {
										Type:     schema.TypeString,
										Computed: true,
									},
									isInstanceTemplateNicPrimaryIpv4Address: {
										Type:     schema.TypeString,
										Computed: true,
									},
									isInstanceTemplateNicPrimaryIP: {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The primary IP address to bind to the network interface. This can be specified using an existing reserved IP, or a prototype object for a new reserved IP.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												isInstanceTemplateNicReservedIpAddress: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The IP address to reserve, which must not already be reserved on the subnet.",
												},
												isInstanceTemplateNicReservedIpName: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The user-defined name for this reserved IP. If unspecified, the name will be a hyphenated list of randomly-selected words. Names must be unique within the subnet the reserved IP resides in. ",
												},
												isInstanceTemplateNicReservedIpId: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Identifies a reserved IP by a unique property.",
												},
											},
										},
									},
									isInstanceTemplateNicSecurityGroups: {
										Type:     schema.TypeSet,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Set:      schema.HashString,
									},
									isInstanceTemplateNicSubnet: {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},

						isInstanceTemplateNetworkInterfaces: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isInstanceTemplateNicName: {
										Type:     schema.TypeString,
										Computed: true,
									},
									isInstanceTemplateNicPrimaryIpv4Address: {
										Type:     schema.TypeString,
										Computed: true,
									},
									isInstanceTemplateNicPrimaryIP: {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The primary IP address to bind to the network interface. This can be specified using an existing reserved IP, or a prototype object for a new reserved IP.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												isInstanceTemplateNicReservedIpAddress: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The IP address to reserve, which must not already be reserved on the subnet.",
												},
												isInstanceTemplateNicReservedIpName: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The user-defined name for this reserved IP. If unspecified, the name will be a hyphenated list of randomly-selected words. Names must be unique within the subnet the reserved IP resides in. ",
												},
												isInstanceTemplateNicReservedIpId: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Identifies a reserved IP by a unique property.",
												},
											},
										},
									},
									isInstanceTemplateNicSecurityGroups: {
										Type:     schema.TypeSet,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Set:      schema.HashString,
									},
									isInstanceTemplateNicSubnet: {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						isInstanceTemplateUserData: {
							Type:     schema.TypeString,
							Computed: true,
						},
						isInstanceTemplateImage: {
							Type:     schema.TypeString,
							Computed: true,
						},
						isInstanceTemplatesBootVolumeAttachment: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isInstanceTemplatesDeleteVol: {
										Type:     schema.TypeBool,
										Computed: true,
									},
									isInstanceTemplatesName: {
										Type:     schema.TypeString,
										Computed: true,
									},
									isInstanceTemplatesVol: {
										Type:     schema.TypeString,
										Computed: true,
									},
									isInstanceTemplateBootSize: {
										Type:     schema.TypeInt,
										Computed: true,
									},
									isInstanceTemplateBootProfile: {
										Type:     schema.TypeString,
										Computed: true,
									},
									isInstanceTemplateBootVolumeTags: {
										Type:        schema.TypeSet,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Set:         flex.ResourceIBMVPCHash,
										Description: "The user tags associated with this volume.",
									},
								},
							},
						},
						isInstanceTemplateResourceGroup: {
							Type:     schema.TypeString,
							Computed: true,
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
									isReservationAffinityPool: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The reservation associated with this template.",
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

func dataSourceIBMISInstanceTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	instanceC, err := vpcClient(meta)
	if err != nil {
		return err
	}
	listInstanceTemplatesOptions := &vpcv1.ListInstanceTemplatesOptions{}
	availableTemplates, _, err := instanceC.ListInstanceTemplates(listInstanceTemplatesOptions)
	if err != nil {
		return err
	}
	templates := make([]map[string]interface{}, 0)
	for _, instTempl := range availableTemplates.Templates {
		template := map[string]interface{}{}
		instance := instTempl.(*vpcv1.InstanceTemplate)
		template["id"] = instance.ID
		template[isInstanceTemplatesHref] = instance.Href
		template[isInstanceTemplatesCrn] = instance.CRN
		template[isInstanceTemplateName] = instance.Name
		template[isInstanceTemplateUserData] = instance.UserData

		if instance.DefaultTrustedProfile != nil {
			if instance.DefaultTrustedProfile.AutoLink != nil {
				template[isInstanceDefaultTrustedProfileAutoLink] = instance.DefaultTrustedProfile.AutoLink
			}
			if instance.DefaultTrustedProfile.Target != nil {
				switch reflect.TypeOf(instance.DefaultTrustedProfile.Target).String() {
				case "*vpcv1.TrustedProfileIdentityTrustedProfileByID":
					{
						target := instance.DefaultTrustedProfile.Target.(*vpcv1.TrustedProfileIdentityByID)
						template[isInstanceDefaultTrustedProfileTarget] = target.ID
					}
				case "*vpcv1.TrustedProfileIdentityTrustedProfileByCRN":
					{
						target := instance.DefaultTrustedProfile.Target.(*vpcv1.TrustedProfileIdentityByCRN)
						template[isInstanceDefaultTrustedProfileTarget] = target.CRN
					}
				}
			}
		}

		if instance.PlacementTarget != nil {
			placementTargetMap := resourceIbmIsInstanceTemplateInstancePlacementTargetPrototypeToMap(*instance.PlacementTarget.(*vpcv1.InstancePlacementTargetPrototype))
			template["placement_target"] = []map[string]interface{}{placementTargetMap}
		}

		// catalog offering if any
		if instance.CatalogOffering != nil {
			catOfferingList := make([]map[string]interface{}, 0)
			insTempCatalogOffering := instance.CatalogOffering.(*vpcv1.InstanceCatalogOfferingPrototype)

			currentOffering := map[string]interface{}{}
			if insTempCatalogOffering.Offering != nil {
				offering := insTempCatalogOffering.Offering.(*vpcv1.CatalogOfferingIdentity)
				currentOffering[isInstanceTemplateCatalogOfferingOfferingCrn] = *offering.CRN
			}
			if insTempCatalogOffering.Version != nil {
				version := insTempCatalogOffering.Version.(*vpcv1.CatalogOfferingVersionIdentity)
				currentOffering[isInstanceTemplateCatalogOfferingVersionCrn] = *version.CRN
			}
			if insTempCatalogOffering.Plan != nil {
				plan := insTempCatalogOffering.Plan.(*vpcv1.CatalogOfferingVersionPlanIdentity)
				if plan.CRN != nil && *plan.CRN != "" {
					currentOffering[isInstanceTemplateCatalogOfferingPlanCrn] = *plan.CRN
				}
			}
			catOfferingList = append(catOfferingList, currentOffering)
			template[isInstanceTemplateCatalogOffering] = catOfferingList
		}

		template["confidential_compute_mode"] = instance.ConfidentialComputeMode
		if instance.ClusterNetworkAttachments != nil {
			clusterNetworkAttachments := []map[string]interface{}{}
			for _, clusterNetworkAttachmentsItem := range instance.ClusterNetworkAttachments {
				clusterNetworkAttachmentsItemMap, err := DataSourceIBMIsInstanceTemplatesInstanceClusterNetworkAttachmentPrototypeInstanceContextToMap(&clusterNetworkAttachmentsItem) // #nosec G601
				if err != nil {
					return err
				}
				clusterNetworkAttachments = append(clusterNetworkAttachments, clusterNetworkAttachmentsItemMap)
			}
			template["cluster_network_attachments"] = clusterNetworkAttachments
		}
		template["enable_secure_boot"] = instance.EnableSecureBoot

		if instance.MetadataService != nil {
			template[isInstanceTemplateMetadataServiceEnabled] = *instance.MetadataService.Enabled

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
			template[isInstanceMetadataService] = metadataService
		}

		if instance.AvailabilityPolicy != nil && instance.AvailabilityPolicy.HostFailure != nil {
			template[isInstanceTemplateAvailablePolicyHostFailure] = *instance.AvailabilityPolicy.HostFailure
		}
		if instance.Keys != nil {
			keys := []string{}
			for _, intfc := range instance.Keys {
				instanceKeyIntf := intfc.(*vpcv1.KeyIdentity)
				keys = append(keys, *instanceKeyIntf.ID)
			}
			template[isInstanceTemplateKeys] = keys
		}
		if instance.Profile != nil {
			instanceProfileIntf := instance.Profile
			identity := instanceProfileIntf.(*vpcv1.InstanceProfileIdentity)
			template[isInstanceTemplateProfile] = identity.Name
		}

		if instance.PlacementTarget != nil {
			placementTargetList := []map[string]interface{}{}
			placementTargetMap := dataSourceInstanceTemplateCollectionTemplatesPlacementTargetToMap(*instance.PlacementTarget.(*vpcv1.InstancePlacementTargetPrototype))
			placementTargetList = append(placementTargetList, placementTargetMap)
			template["placement_target"] = placementTargetList
		}

		if instance.TotalVolumeBandwidth != nil {
			template[isInstanceTotalVolumeBandwidth] = int(*instance.TotalVolumeBandwidth)
		}

		// vni

		networkAttachments := []map[string]interface{}{}
		if instance.NetworkAttachments != nil {
			for _, modelItem := range instance.NetworkAttachments {
				modelMap, err := dataSourceIBMIsInstanceTemplateInstanceNetworkAttachmentPrototypeToMap(&modelItem)
				if err != nil {
					return err
				}
				networkAttachments = append(networkAttachments, modelMap)
			}
		}
		template["network_attachments"] = networkAttachments

		primaryNetworkAttachment := []map[string]interface{}{}
		if instance.PrimaryNetworkAttachment != nil {
			modelMap, err := dataSourceIBMIsInstanceTemplateInstanceNetworkAttachmentPrototypeToMap(instance.PrimaryNetworkAttachment)
			if err != nil {
				return err
			}
			primaryNetworkAttachment = append(primaryNetworkAttachment, modelMap)
		}
		template["primary_network_attachment"] = primaryNetworkAttachment

		if instance.PrimaryNetworkInterface != nil {
			interfaceList := make([]map[string]interface{}, 0)
			currentPrimNic := map[string]interface{}{}
			currentPrimNic[isInstanceTemplateNicName] = *instance.PrimaryNetworkInterface.Name

			if instance.PrimaryNetworkInterface.PrimaryIP != nil {
				primaryipIntf := instance.PrimaryNetworkInterface.PrimaryIP
				primaryIpList := make([]map[string]interface{}, 0)
				currentPrimIp := map[string]interface{}{}
				switch reflect.TypeOf(primaryipIntf).String() {
				case "*vpcv1.NetworkInterfaceIPPrototype":
					{
						primaryip := primaryipIntf.(*vpcv1.NetworkInterfaceIPPrototype)
						currentPrimNic[isInstanceTemplateNicPrimaryIpv4Address] = primaryip.Address
						currentPrimIp[isInstanceTemplateNicReservedIpAddress] = primaryip.Address
						currentPrimIp[isInstanceTemplateNicReservedIpId] = primaryip.ID
					}
				case "*vpcv1.NetworkInterfaceIPPrototypeReservedIPPrototypeNetworkInterfaceContext":
					{
						primaryip := primaryipIntf.(*vpcv1.NetworkInterfaceIPPrototypeReservedIPPrototypeNetworkInterfaceContext)
						currentPrimNic[isInstanceTemplateNicPrimaryIpv4Address] = primaryip.Address
						currentPrimIp[isInstanceTemplateNicReservedIpAddress] = primaryip.Address
					}
				case "*vpcv1.NetworkInterfaceIPPrototypeReservedIPIdentity":
					{
						primaryip := primaryipIntf.(*vpcv1.NetworkInterfaceIPPrototypeReservedIPIdentity)
						currentPrimIp[isInstanceTemplateNicReservedIpId] = primaryip.ID
					}
				}
				primaryIpList = append(primaryIpList, currentPrimIp)
				currentPrimNic[isInstanceTemplateNicPrimaryIP] = primaryIpList
			}
			subInf := instance.PrimaryNetworkInterface.Subnet
			subnetIdentity := subInf.(*vpcv1.SubnetIdentity)
			currentPrimNic[isInstanceTemplateNicSubnet] = *subnetIdentity.ID

			if len(instance.PrimaryNetworkInterface.SecurityGroups) != 0 {
				secgrpList := []string{}
				for i := 0; i < len(instance.PrimaryNetworkInterface.SecurityGroups); i++ {
					secGrpInf := instance.PrimaryNetworkInterface.SecurityGroups[i]
					secGrpIdentity := secGrpInf.(*vpcv1.SecurityGroupIdentity)
					secgrpList = append(secgrpList, string(*secGrpIdentity.ID))
				}
				currentPrimNic[isInstanceTemplateNicSecurityGroups] = flex.NewStringSet(schema.HashString, secgrpList)
			}
			interfaceList = append(interfaceList, currentPrimNic)
			template[isInstanceTemplatePrimaryNetworkInterface] = interfaceList
		}

		if instance.NetworkInterfaces != nil {
			interfacesList := make([]map[string]interface{}, 0)
			for _, intfc := range instance.NetworkInterfaces {
				currentNic := map[string]interface{}{}
				currentNic[isInstanceTemplateNicName] = *intfc.Name
				if intfc.PrimaryIP != nil {
					primaryipIntf := intfc.PrimaryIP
					primaryIpList := make([]map[string]interface{}, 0)
					currentPrimIp := map[string]interface{}{}
					switch reflect.TypeOf(primaryipIntf).String() {
					case "*vpcv1.NetworkInterfaceIPPrototype":
						{
							primaryip := primaryipIntf.(*vpcv1.NetworkInterfaceIPPrototype)
							currentNic[isInstanceTemplateNicPrimaryIpv4Address] = primaryip.Address
							currentPrimIp[isInstanceTemplateNicReservedIpAddress] = primaryip.Address
							currentPrimIp[isInstanceTemplateNicReservedIpId] = primaryip.ID
						}
					case "*vpcv1.NetworkInterfaceIPPrototypeReservedIPPrototypeNetworkInterfaceContext":
						{
							primaryip := primaryipIntf.(*vpcv1.NetworkInterfaceIPPrototypeReservedIPPrototypeNetworkInterfaceContext)
							currentNic[isInstanceTemplateNicPrimaryIpv4Address] = primaryip.Address
							currentPrimIp[isInstanceTemplateNicReservedIpAddress] = primaryip.Address
						}
					case "*vpcv1.NetworkInterfaceIPPrototypeReservedIPIdentity":
						{
							primaryip := primaryipIntf.(*vpcv1.NetworkInterfaceIPPrototypeReservedIPIdentity)
							currentPrimIp[isInstanceTemplateNicReservedIpId] = primaryip.ID
						}
					}
					primaryIpList = append(primaryIpList, currentPrimIp)
					currentNic[isInstanceTemplateNicPrimaryIP] = primaryIpList
				}
				//currentNic[isInstanceTemplateNicAllowIpSpoofing] = intfc.AllowIpSpoofing
				subInf := intfc.Subnet
				subnetIdentity := subInf.(*vpcv1.SubnetIdentity)
				currentNic[isInstanceTemplateNicSubnet] = *subnetIdentity.ID
				if len(intfc.SecurityGroups) != 0 {
					secgrpList := []string{}
					for i := 0; i < len(intfc.SecurityGroups); i++ {
						secGrpInf := intfc.SecurityGroups[i]
						secGrpIdentity := secGrpInf.(*vpcv1.SecurityGroupIdentity)
						secgrpList = append(secgrpList, string(*secGrpIdentity.ID))
					}
					currentNic[isInstanceTemplateNicSecurityGroups] = flex.NewStringSet(schema.HashString, secgrpList)
				}

				interfacesList = append(interfacesList, currentNic)
			}
			template[isInstanceTemplateNetworkInterfaces] = interfacesList
		}

		if instance.Image != nil {
			imageInf := instance.Image
			imageIdentity := imageInf.(*vpcv1.ImageIdentity)
			template[isInstanceTemplateImage] = imageIdentity.ID
		}

		if instance.VPC != nil {
			vpcInf := instance.VPC
			vpcRef := vpcInf.(*vpcv1.VPCIdentity)
			template[isInstanceTemplateVPC] = vpcRef.ID
		}

		if instance.Zone != nil {
			zoneInf := instance.Zone
			zone := zoneInf.(*vpcv1.ZoneIdentity)
			template[isInstanceTemplateZone] = zone.Name
		}

		interfacesList := make([]map[string]interface{}, 0)
		if instance.VolumeAttachments != nil {
			for _, volume := range instance.VolumeAttachments {
				volumeAttach := map[string]interface{}{}
				volumeAttach[isInstanceTemplateVolAttName] = *volume.Name
				volumeAttach[isInstanceTemplateDeleteVolume] = *volume.DeleteVolumeOnInstanceDelete
				volumeIntf := volume.Volume
				volumeInst := volumeIntf.(*vpcv1.VolumeAttachmentPrototypeVolume)
				newVolumeArr := []map[string]interface{}{}
				newVolume := map[string]interface{}{}

				if volumeInst.ID != nil {
					volumeAttach[isInstanceTemplateVolAttVolume] = *volumeInst.ID
				}

				if volumeInst.Capacity != nil {
					newVolume[isInstanceTemplateVolAttVolCapacity] = *volumeInst.Capacity
				}
				if volumeInst.Profile != nil {
					profile := volumeInst.Profile.(*vpcv1.VolumeProfileIdentity)
					newVolume[isInstanceTemplateVolAttVolProfile] = profile.Name
				}

				if volumeInst.Iops != nil {
					newVolume[isInstanceTemplateVolAttVolIops] = *volumeInst.Iops
				}
				if volumeInst.EncryptionKey != nil {
					encryptionKey := volumeInst.EncryptionKey.(*vpcv1.EncryptionKeyIdentity)
					newVolume[isInstanceTemplateVolAttVolEncryptionKey] = *encryptionKey.CRN
				}
				if volumeInst.UserTags != nil {
					newVolume[isInstanceTemplateVolAttTags] = instance.BootVolumeAttachment.Volume.UserTags
				}
				newVolumeArr = append(newVolumeArr, newVolume)
				volumeAttach[isInstanceTemplateVolAttVolPrototype] = newVolumeArr

				interfacesList = append(interfacesList, volumeAttach)
			}
			template[isInstanceTemplateVolumeAttachments] = interfacesList
		}

		if instance.BootVolumeAttachment != nil {
			bootVolList := make([]map[string]interface{}, 0)
			bootVol := map[string]interface{}{}

			bootVol[isInstanceTemplatesDeleteVol] = *instance.BootVolumeAttachment.DeleteVolumeOnInstanceDelete
			if instance.BootVolumeAttachment.Volume != nil {
				volumeIntf := instance.BootVolumeAttachment.Volume
				bootVol[isInstanceTemplatesName] = volumeIntf.Name
				bootVol[isInstanceTemplatesVol] = volumeIntf.Name
				bootVol[isInstanceTemplateBootSize] = volumeIntf.Capacity
				if instance.BootVolumeAttachment.Volume.Profile != nil {
					volProfIntf := instance.BootVolumeAttachment.Volume.Profile
					volProfInst := volProfIntf.(*vpcv1.VolumeProfileIdentity)
					bootVol[isInstanceTemplateBootProfile] = volProfInst.Name
				}
				if instance.BootVolumeAttachment.Volume.UserTags != nil {
					bootVol[isInstanceTemplateBootVolumeTags] = instance.BootVolumeAttachment.Volume.UserTags
				}
			}
			bootVolList = append(bootVolList, bootVol)
			template[isInstanceTemplatesBootVolumeAttachment] = bootVolList
		}
		if instance.ReservationAffinity != nil {
			reservationAffinity := []map[string]interface{}{}
			reservationAffinityMap := map[string]interface{}{}

			reservationAffinityMap[isReservationAffinityPolicyResp] = instance.ReservationAffinity.Policy
			if instance.ReservationAffinity.Pool != nil && len(instance.ReservationAffinity.Pool) > 0 {
				pool := instance.ReservationAffinity.Pool[0]
				res := ""
				if idPool, ok := pool.(*vpcv1.ReservationIdentityByID); ok {
					res = *idPool.ID
				} else if crnPool, ok := pool.(*vpcv1.ReservationIdentityByCRN); ok {
					res = *crnPool.CRN
				} else if hrefPool, ok := pool.(*vpcv1.ReservationIdentityByHref); ok {
					res = *hrefPool.Href
				}
				reservationAffinityMap[isReservationAffinityPool] = res
			}
			reservationAffinity = append(reservationAffinity, reservationAffinityMap)
			template[isReservationAffinity] = reservationAffinity
		}

		if instance.ResourceGroup != nil {
			rg := instance.ResourceGroup
			template[isInstanceTemplateResourceGroup] = rg.ID
		}

		templates = append(templates, template)
	}
	d.SetId(dataSourceIBMISInstanceTemplatesID(d))
	d.Set(isInstanceTemplates, templates)
	return nil
}

// dataSourceIBMISInstanceTemplatesID returns a reasonable ID for a instance templates list.
func dataSourceIBMISInstanceTemplatesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceInstanceTemplateCollectionTemplatesPlacementTargetToMap(placementTargetItem vpcv1.InstancePlacementTargetPrototype) (placementTargetMap map[string]interface{}) {
	placementTargetMap = map[string]interface{}{}

	if placementTargetItem.ID != nil {
		placementTargetMap["id"] = placementTargetItem.ID
	}
	if placementTargetItem.CRN != nil {
		placementTargetMap["crn"] = placementTargetItem.CRN
	}
	if placementTargetItem.Href != nil {
		placementTargetMap["href"] = placementTargetItem.Href
	}

	return placementTargetMap
}

func DataSourceIBMIsInstanceTemplatesInstanceClusterNetworkAttachmentPrototypeInstanceContextToMap(model *vpcv1.InstanceClusterNetworkAttachmentPrototypeInstanceContext) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	clusterNetworkInterfaceMap, err := DataSourceIBMIsInstanceTemplatesInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceToMap(model.ClusterNetworkInterface)
	if err != nil {
		return modelMap, err
	}
	modelMap["cluster_network_interface"] = []map[string]interface{}{clusterNetworkInterfaceMap}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	return modelMap, nil
}

func DataSourceIBMIsInstanceTemplatesInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceToMap(model vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceInstanceClusterNetworkInterfacePrototypeInstanceClusterNetworkAttachment); ok {
		return DataSourceIBMIsInstanceTemplatesInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceInstanceClusterNetworkInterfacePrototypeInstanceClusterNetworkAttachmentToMap(model.(*vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceInstanceClusterNetworkInterfacePrototypeInstanceClusterNetworkAttachment))
	} else if _, ok := model.(*vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentity); ok {
		return DataSourceIBMIsInstanceTemplatesInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityToMap(model.(*vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentity))
	} else if _, ok := model.(*vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterface); ok {
		modelMap := make(map[string]interface{})
		model := model.(*vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterface)
		if model.AutoDelete != nil {
			modelMap["auto_delete"] = *model.AutoDelete
		}
		if model.Name != nil {
			modelMap["name"] = *model.Name
		}
		if model.PrimaryIP != nil {
			primaryIPMap, err := DataSourceIBMIsInstanceTemplatesClusterNetworkInterfacePrimaryIPPrototypeToMap(model.PrimaryIP)
			if err != nil {
				return modelMap, err
			}
			modelMap["primary_ip"] = []map[string]interface{}{primaryIPMap}
		}
		if model.Subnet != nil {
			subnetMap, err := DataSourceIBMIsInstanceTemplatesClusterNetworkSubnetIdentityToMap(model.Subnet)
			if err != nil {
				return modelMap, err
			}
			modelMap["subnet"] = []map[string]interface{}{subnetMap}
		}
		if model.ID != nil {
			modelMap["id"] = *model.ID
		}
		if model.Href != nil {
			modelMap["href"] = *model.Href
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceIntf subtype encountered")
	}
}

func DataSourceIBMIsInstanceTemplatesClusterNetworkInterfacePrimaryIPPrototypeToMap(model vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContext); ok {
		return DataSourceIBMIsInstanceTemplatesClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextToMap(model.(*vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContext))
	} else if _, ok := model.(*vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPPrototypeClusterNetworkInterfacePrimaryIPContext); ok {
		return DataSourceIBMIsInstanceTemplatesClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPPrototypeClusterNetworkInterfacePrimaryIPContextToMap(model.(*vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPPrototypeClusterNetworkInterfacePrimaryIPContext))
	} else if _, ok := model.(*vpcv1.ClusterNetworkInterfacePrimaryIPPrototype); ok {
		modelMap := make(map[string]interface{})
		model := model.(*vpcv1.ClusterNetworkInterfacePrimaryIPPrototype)
		if model.ID != nil {
			modelMap["id"] = *model.ID
		}
		if model.Href != nil {
			modelMap["href"] = *model.Href
		}
		if model.Address != nil {
			modelMap["address"] = *model.Address
		}
		if model.AutoDelete != nil {
			modelMap["auto_delete"] = *model.AutoDelete
		}
		if model.Name != nil {
			modelMap["name"] = *model.Name
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeIntf subtype encountered")
	}
}

func DataSourceIBMIsInstanceTemplatesClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextToMap(model vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByID); ok {
		return DataSourceIBMIsInstanceTemplatesClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByIDToMap(model.(*vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByID))
	} else if _, ok := model.(*vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByHref); ok {
		return DataSourceIBMIsInstanceTemplatesClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByHrefToMap(model.(*vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByHref))
	} else if _, ok := model.(*vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContext); ok {
		modelMap := make(map[string]interface{})
		model := model.(*vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContext)
		if model.ID != nil {
			modelMap["id"] = *model.ID
		}
		if model.Href != nil {
			modelMap["href"] = *model.Href
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextIntf subtype encountered")
	}
}

func DataSourceIBMIsInstanceTemplatesClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByIDToMap(model *vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByID) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	return modelMap, nil
}

func DataSourceIBMIsInstanceTemplatesClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByHrefToMap(model *vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByHref) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	return modelMap, nil
}

func DataSourceIBMIsInstanceTemplatesClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPPrototypeClusterNetworkInterfacePrimaryIPContextToMap(model *vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPPrototypeClusterNetworkInterfacePrimaryIPContext) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Address != nil {
		modelMap["address"] = *model.Address
	}
	if model.AutoDelete != nil {
		modelMap["auto_delete"] = *model.AutoDelete
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	return modelMap, nil
}

func DataSourceIBMIsInstanceTemplatesClusterNetworkSubnetIdentityToMap(model vpcv1.ClusterNetworkSubnetIdentityIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.ClusterNetworkSubnetIdentityByID); ok {
		return DataSourceIBMIsInstanceTemplatesClusterNetworkSubnetIdentityByIDToMap(model.(*vpcv1.ClusterNetworkSubnetIdentityByID))
	} else if _, ok := model.(*vpcv1.ClusterNetworkSubnetIdentityByHref); ok {
		return DataSourceIBMIsInstanceTemplatesClusterNetworkSubnetIdentityByHrefToMap(model.(*vpcv1.ClusterNetworkSubnetIdentityByHref))
	} else if _, ok := model.(*vpcv1.ClusterNetworkSubnetIdentity); ok {
		modelMap := make(map[string]interface{})
		model := model.(*vpcv1.ClusterNetworkSubnetIdentity)
		if model.ID != nil {
			modelMap["id"] = *model.ID
		}
		if model.Href != nil {
			modelMap["href"] = *model.Href
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized vpcv1.ClusterNetworkSubnetIdentityIntf subtype encountered")
	}
}

func DataSourceIBMIsInstanceTemplatesClusterNetworkSubnetIdentityByIDToMap(model *vpcv1.ClusterNetworkSubnetIdentityByID) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	return modelMap, nil
}

func DataSourceIBMIsInstanceTemplatesClusterNetworkSubnetIdentityByHrefToMap(model *vpcv1.ClusterNetworkSubnetIdentityByHref) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	return modelMap, nil
}

func DataSourceIBMIsInstanceTemplatesInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceInstanceClusterNetworkInterfacePrototypeInstanceClusterNetworkAttachmentToMap(model *vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceInstanceClusterNetworkInterfacePrototypeInstanceClusterNetworkAttachment) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AutoDelete != nil {
		modelMap["auto_delete"] = *model.AutoDelete
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.PrimaryIP != nil {
		primaryIPMap, err := DataSourceIBMIsInstanceTemplatesClusterNetworkInterfacePrimaryIPPrototypeToMap(model.PrimaryIP)
		if err != nil {
			return modelMap, err
		}
		modelMap["primary_ip"] = []map[string]interface{}{primaryIPMap}
	}
	if model.Subnet != nil {
		subnetMap, err := DataSourceIBMIsInstanceTemplatesClusterNetworkSubnetIdentityToMap(model.Subnet)
		if err != nil {
			return modelMap, err
		}
		modelMap["subnet"] = []map[string]interface{}{subnetMap}
	}
	return modelMap, nil
}

func DataSourceIBMIsInstanceTemplatesInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityToMap(model vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByID); ok {
		return DataSourceIBMIsInstanceTemplatesInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByIDToMap(model.(*vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByID))
	} else if _, ok := model.(*vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByHref); ok {
		return DataSourceIBMIsInstanceTemplatesInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByHrefToMap(model.(*vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByHref))
	} else if _, ok := model.(*vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentity); ok {
		modelMap := make(map[string]interface{})
		model := model.(*vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentity)
		if model.ID != nil {
			modelMap["id"] = *model.ID
		}
		if model.Href != nil {
			modelMap["href"] = *model.Href
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityIntf subtype encountered")
	}
}

func DataSourceIBMIsInstanceTemplatesInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByIDToMap(model *vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByID) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	return modelMap, nil
}

func DataSourceIBMIsInstanceTemplatesInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByHrefToMap(model *vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByHref) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	return modelMap, nil
}
