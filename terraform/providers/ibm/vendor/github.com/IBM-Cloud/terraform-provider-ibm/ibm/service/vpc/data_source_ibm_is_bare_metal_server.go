// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMIsBareMetalServer() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMISBareMetalServerRead,

		Schema: map[string]*schema.Schema{
			"identifier": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				AtLeastOneOf:  []string{isBareMetalServerName, "identifier"},
				ConflictsWith: []string{isBareMetalServerName},
				ValidateFunc:  validate.InvokeDataSourceValidator("ibm_is_bare_metal_server", "identifier"),
			},
			isBareMetalServerName: {
				Type:          schema.TypeString,
				Optional:      true,
				AtLeastOneOf:  []string{isBareMetalServerName, "identifier"},
				Computed:      true,
				ConflictsWith: []string{"identifier"},
				ValidateFunc:  validate.InvokeValidator("ibm_is_bare_metal_server", isBareMetalServerName),
				Description:   "Bare metal server name",
			},
			isBareMetalServerBandwidth: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total bandwidth (in megabits per second)",
			},
			isBareMetalServerEnableSecureBoot: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether secure boot is enabled. If enabled, the image must support secure boot or the server will fail to boot.",
			},

			isBareMetalServerTrustedPlatformModule: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isBareMetalServerTrustedPlatformModuleMode: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The trusted platform module mode to use. The specified value must be listed in the bare metal server profile's supported_trusted_platform_module_modes",
						},
						isBareMetalServerTrustedPlatformModuleEnabled: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether the trusted platform module is enabled.",
						},
						isBareMetalServerTrustedPlatformModuleSupportedModes: {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         flex.ResourceIBMVPCHash,
							Computed:    true,
							Description: "The trusted platform module (TPM) mode:: disabled: No TPM functionality, tpm_2: TPM 2.0. The enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected property value was encountered. Enum: [ disabled, tpm_2 ]",
						},
					},
				},
			},

			isBareMetalServerBootTarget: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for this bare metal server disk",
			},
			isBareMetalServerCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the bare metal server was created",
			},
			isBareMetalServerCPU: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The bare metal server CPU configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isBareMetalServerCPUArchitecture: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CPU architecture",
						},
						isBareMetalServerCPUCoreCount: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total number of cores",
						},
						isBareMetalServerCpuSocketCount: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total number of CPU sockets",
						},
						isBareMetalServerCpuThreadPerCore: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total number of hardware threads per core",
						},
					},
				},
			},
			isBareMetalServerCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this bare metal server",
			},
			isBareMetalServerDisks: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The disks for this bare metal server, including any disks that are associated with the boot_target.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isBareMetalServerDiskHref: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this bare metal server disk",
						},
						isBareMetalServerDiskID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this bare metal server disk",
						},
						isBareMetalServerDiskInterfaceType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The disk interface used for attaching the disk. Supported values are [ nvme, sata ]",
						},
						isBareMetalServerDiskName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this disk",
						},
						isBareMetalServerDiskResourceType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type",
						},
						isBareMetalServerDiskSize: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The size of the disk in GB (gigabytes)",
						},
					},
				},
			},
			isBareMetalServerHref: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this bare metal server",
			},
			isBareMetalServerMemory: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The amount of memory, truncated to whole gibibytes",
			},

			isBareMetalServerPrimaryNetworkInterface: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Primary Network interface info",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						isBareMetalServerNicAllowIPSpoofing: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether IP spoofing is allowed on this interface.",
						},
						isBareMetalServerNicName: {
							Type:     schema.TypeString,
							Computed: true,
						},
						isBareMetalServerNicPortSpeed: {
							Type:     schema.TypeInt,
							Computed: true,
						},
						isBareMetalServerNicHref: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "This URL of the interface",
						},

						isBareMetalServerNicSecurityGroups: {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						isBareMetalServerNicSubnet: {
							Type:     schema.TypeString,
							Computed: true,
						},
						isBareMetalServerNicPrimaryIP: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "IPv4, The IP address. ",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isBareMetalServerNicIpAddress: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The globally unique IP address",
									},
									isBareMetalServerNicIpHref: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this reserved IP",
									},
									isBareMetalServerNicIpName: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this reserved IP. If unspecified, the name will be a hyphenated list of randomly-selected words. Names must be unique within the subnet the reserved IP resides in. ",
									},
									isBareMetalServerNicIpID: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Identifies a reserved IP by a unique property.",
									},
									isBareMetalServerNicResourceType: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type",
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
				Description: "The primary network attachment.",
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
							Description: "The virtual network interface for this bare metal server network attachment.",
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

			isBareMetalServerNetworkInterfaces: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						isBareMetalServerNicHref: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this network interface",
						},
						isBareMetalServerNicAllowIPSpoofing: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether IP spoofing is allowed on this interface.",
						},
						isBareMetalServerNicName: {
							Type:     schema.TypeString,
							Computed: true,
						},
						isBareMetalServerNicSecurityGroups: {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						isBareMetalServerNicSubnet: {
							Type:     schema.TypeString,
							Computed: true,
						},
						isBareMetalServerNicPrimaryIP: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "IPv4, The IP address. ",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isBareMetalServerNicIpAddress: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The globally unique IP address",
									},
									isBareMetalServerNicIpHref: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this reserved IP",
									},
									isBareMetalServerNicIpName: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this reserved IP. If unspecified, the name will be a hyphenated list of randomly-selected words. Names must be unique within the subnet the reserved IP resides in. ",
									},
									isBareMetalServerNicIpID: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Identifies a reserved IP by a unique property.",
									},
									isBareMetalServerNicResourceType: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type",
									},
								},
							},
						},
					},
				},
			},

			"network_attachments": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The network attachments for this bare metal server, including the primary network attachment.",
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
							Description: "The virtual network interface for this bare metal server network attachment.",
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

			isBareMetalServerKeys: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "SSH key Ids for the bare metal server",
			},

			isBareMetalServerImage: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "image name",
			},
			isBareMetalServerFirmwareUpdateTypeAvailable: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of firmware update available",
			},
			isBareMetalServerProfile: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "profil name",
			},

			isBareMetalServerZone: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Zone name",
			},

			isBareMetalServerVPC: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The VPC the bare metal server is to be a part of",
			},

			isBareMetalServerResourceGroup: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource group name",
			},
			isBareMetalServerResourceType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource type name",
			},

			isBareMetalServerStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Bare metal server status",
			},

			isBareMetalServerStatusReasons: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isBareMetalServerStatusReasonsCode: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A snake case string succinctly identifying the status reason",
						},

						isBareMetalServerStatusReasonsMessage: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An explanation of the status reason",
						},

						isBareMetalServerStatusReasonsMoreInfo: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Link to documentation about this status reason",
						},
					},
				},
			},
			isBareMetalServerTags: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_bare_metal_server", "tags")},
				Set:         flex.ResourceIBMVPCHash,
				Description: "Tags for the Bare metal server",
			},
			isBareMetalServerAccessTags: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of access tags",
			},
		},
	}
}

func DataSourceIBMIsBareMetalServerValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "identifier",
			ValidateFunctionIdentifier: validate.ValidateNoZeroValues,
			Type:                       validate.TypeString})

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isBareMetalServerName,
			ValidateFunctionIdentifier: validate.ValidateNoZeroValues,
			Type:                       validate.TypeString})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "tags",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9:_ .-]+$`,
			MinValueLength:             1,
			MaxValueLength:             128})

	ibmISBMSDataSourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_bare_metal_server", Schema: validateSchema}
	return &ibmISBMSDataSourceValidator
}

func dataSourceIBMISBareMetalServerRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := d.Get("identifier").(string)
	name := d.Get(isBareMetalServerName).(string)

	sess, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}
	var bms *vpcv1.BareMetalServer
	if id != "" {
		options := &vpcv1.GetBareMetalServerOptions{}
		options.ID = &id
		server, response, err := sess.GetBareMetalServerWithContext(context, options)
		if err != nil {
			log.Printf("[DEBUG] GetBareMetalServerWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("[ERROR] Error Getting Bare Metal Server (%s): %s\n%s", id, err, response))
		}
		bms = server
	} else if name != "" {
		options := &vpcv1.ListBareMetalServersOptions{}
		options.Name = &name
		bmservers, response, err := sess.ListBareMetalServersWithContext(context, options)
		if err != nil {
			log.Printf("[DEBUG] ListBareMetalServersWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("[ERROR] Error Listing Bare Metal Server (%s): %s\n%s", name, err, response))
		}
		if len(bmservers.BareMetalServers) == 0 {
			return diag.FromErr(fmt.Errorf("[ERROR] No bare metal servers found with name %s", name))
		}
		bms = &bmservers.BareMetalServers[0]
	}

	d.SetId(*bms.ID)
	d.Set(isBareMetalServerBandwidth, bms.Bandwidth)
	if bms.BootTarget != nil {
		bmsBootTargetIntf := bms.BootTarget.(*vpcv1.BareMetalServerBootTarget)
		bmsBootTarget := bmsBootTargetIntf.ID
		d.Set(isBareMetalServerBootTarget, bmsBootTarget)
	}

	// set keys and image using initialization

	optionsInitialization := &vpcv1.GetBareMetalServerInitializationOptions{
		ID: bms.ID,
	}

	initialization, response, err := sess.GetBareMetalServerInitialization(optionsInitialization)
	if err != nil || initialization == nil {
		return diag.FromErr(fmt.Errorf("[Error] Error getting Bare Metal Server (%s) initialization : %s\n%s", *bms.ID, err, response))
	}

	d.Set(isBareMetalServerImage, initialization.Image.ID)

	keyListList := []string{}
	for i := 0; i < len(initialization.Keys); i++ {
		keyListList = append(keyListList, string(*(initialization.Keys[i].ID)))
	}
	d.Set(isBareMetalServerKeys, keyListList)

	cpuList := make([]map[string]interface{}, 0)
	if bms.Cpu != nil {
		currentCPU := map[string]interface{}{}
		currentCPU[isBareMetalServerCPUArchitecture] = *bms.Cpu.Architecture
		currentCPU[isBareMetalServerCPUCoreCount] = *bms.Cpu.CoreCount
		currentCPU[isBareMetalServerCpuSocketCount] = *bms.Cpu.SocketCount
		currentCPU[isBareMetalServerCpuThreadPerCore] = *bms.Cpu.ThreadsPerCore
		cpuList = append(cpuList, currentCPU)
	}
	d.Set(isBareMetalServerCPU, cpuList)
	if err = d.Set(isBareMetalServerCreatedAt, bms.CreatedAt.String()); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_at: %s", err))
	}
	if err = d.Set(isBareMetalServerCRN, bms.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting crn: %s", err))
	}

	diskList := make([]map[string]interface{}, 0)
	if bms.Disks != nil {
		for _, disk := range bms.Disks {
			currentDisk := map[string]interface{}{
				isBareMetalServerDiskHref:          disk.Href,
				isBareMetalServerDiskID:            disk.ID,
				isBareMetalServerDiskInterfaceType: disk.InterfaceType,
				isBareMetalServerDiskName:          disk.Name,
				isBareMetalServerDiskResourceType:  disk.ResourceType,
				isBareMetalServerDiskSize:          disk.Size,
			}
			diskList = append(diskList, currentDisk)
		}
	}
	d.Set(isBareMetalServerDisks, diskList)
	if err = d.Set(isBareMetalServerHref, bms.Href); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting href: %s", err))
	}
	if err = d.Set(isBareMetalServerMemory, bms.Memory); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting memory: %s", err))
	}
	if err = d.Set(isBareMetalServerName, *bms.Name); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
	}
	if err = d.Set("identifier", *bms.ID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting identifier: %s", err))
	}
	if bms.Firmware != nil && bms.Firmware.Update != nil {
		if err = d.Set(isBareMetalServerFirmwareUpdateTypeAvailable, *bms.Firmware.Update); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting availble firmware update type: %s", err))
		}
	}

	//enable secure boot
	if err = d.Set(isBareMetalServerEnableSecureBoot, bms.EnableSecureBoot); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting enable_secure_boot: %s", err))
	}

	// tpm
	if bms.TrustedPlatformModule != nil {
		trustedPlatformModuleMap, err := resourceIBMIsBareMetalServerBareMetalServerTrustedPlatformModulePrototypeToMap(bms.TrustedPlatformModule)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set(isBareMetalServerTrustedPlatformModule, []map[string]interface{}{trustedPlatformModuleMap}); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting trusted_platform_module: %s", err))
		}
	}

	//pni

	if bms.PrimaryNetworkInterface != nil {
		primaryNicList := make([]map[string]interface{}, 0)
		currentPrimNic := map[string]interface{}{}
		currentPrimNic["id"] = *bms.PrimaryNetworkInterface.ID
		currentPrimNic[isBareMetalServerNicHref] = *bms.PrimaryNetworkInterface.Href
		currentPrimNic[isBareMetalServerNicName] = *bms.PrimaryNetworkInterface.Name
		currentPrimNic[isBareMetalServerNicHref] = *bms.PrimaryNetworkInterface.Href
		currentPrimNic[isBareMetalServerNicSubnet] = *bms.PrimaryNetworkInterface.Subnet.ID
		if bms.PrimaryNetworkInterface.PrimaryIP != nil {
			primaryIpList := make([]map[string]interface{}, 0)
			currentIP := map[string]interface{}{}
			if bms.PrimaryNetworkInterface.PrimaryIP.Href != nil {
				currentIP[isBareMetalServerNicIpAddress] = *bms.PrimaryNetworkInterface.PrimaryIP.Address
			}
			if bms.PrimaryNetworkInterface.PrimaryIP.Href != nil {
				currentIP[isBareMetalServerNicIpHref] = *bms.PrimaryNetworkInterface.PrimaryIP.Href
			}
			if bms.PrimaryNetworkInterface.PrimaryIP.Name != nil {
				currentIP[isBareMetalServerNicIpName] = *bms.PrimaryNetworkInterface.PrimaryIP.Name
			}
			if bms.PrimaryNetworkInterface.PrimaryIP.ID != nil {
				currentIP[isBareMetalServerNicIpID] = *bms.PrimaryNetworkInterface.PrimaryIP.ID
			}
			if bms.PrimaryNetworkInterface.PrimaryIP.ResourceType != nil {
				currentIP[isBareMetalServerNicResourceType] = *bms.PrimaryNetworkInterface.PrimaryIP.ResourceType
			}
			primaryIpList = append(primaryIpList, currentIP)
			currentPrimNic[isBareMetalServerNicPrimaryIP] = primaryIpList
		}
		getnicoptions := &vpcv1.GetBareMetalServerNetworkInterfaceOptions{
			BareMetalServerID: bms.ID,
			ID:                bms.PrimaryNetworkInterface.ID,
		}
		bmsnic, response, err := sess.GetBareMetalServerNetworkInterface(getnicoptions)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error getting network interfaces attached to the bare metal server %s\n%s", err, response))
		}

		switch reflect.TypeOf(bmsnic).String() {
		case "*vpcv1.BareMetalServerNetworkInterfaceByPci":
			{
				primNic := bmsnic.(*vpcv1.BareMetalServerNetworkInterfaceByPci)
				currentPrimNic[isInstanceNicAllowIPSpoofing] = *primNic.AllowIPSpoofing
				currentPrimNic[isBareMetalServerNicPortSpeed] = *primNic.PortSpeed
				if len(primNic.SecurityGroups) != 0 {
					secgrpList := []string{}
					for i := 0; i < len(primNic.SecurityGroups); i++ {
						secgrpList = append(secgrpList, string(*(primNic.SecurityGroups[i].ID)))
					}
					currentPrimNic[isInstanceNicSecurityGroups] = flex.NewStringSet(schema.HashString, secgrpList)
				}
			}
		case "*vpcv1.BareMetalServerNetworkInterfaceByVlan":
			{
				primNic := bmsnic.(*vpcv1.BareMetalServerNetworkInterfaceByVlan)
				currentPrimNic[isInstanceNicAllowIPSpoofing] = *primNic.AllowIPSpoofing
				currentPrimNic[isBareMetalServerNicPortSpeed] = *primNic.PortSpeed

				if len(primNic.SecurityGroups) != 0 {
					secgrpList := []string{}
					for i := 0; i < len(primNic.SecurityGroups); i++ {
						secgrpList = append(secgrpList, string(*(primNic.SecurityGroups[i].ID)))
					}
					currentPrimNic[isInstanceNicSecurityGroups] = flex.NewStringSet(schema.HashString, secgrpList)
				}
			}
		case "*vpcv1.BareMetalServerNetworkInterfaceByHiperSocket":
			{
				primNic := bmsnic.(*vpcv1.BareMetalServerNetworkInterfaceByHiperSocket)
				currentPrimNic[isInstanceNicAllowIPSpoofing] = *primNic.AllowIPSpoofing

				if len(primNic.SecurityGroups) != 0 {
					secgrpList := []string{}
					for i := 0; i < len(primNic.SecurityGroups); i++ {
						secgrpList = append(secgrpList, string(*(primNic.SecurityGroups[i].ID)))
					}
					currentPrimNic[isInstanceNicSecurityGroups] = flex.NewStringSet(schema.HashString, secgrpList)
				}
			}
		}

		primaryNicList = append(primaryNicList, currentPrimNic)
		d.Set(isBareMetalServerPrimaryNetworkInterface, primaryNicList)
	}

	primaryNetworkAttachment := []map[string]interface{}{}
	if bms.PrimaryNetworkAttachment != nil {
		modelMap, err := dataSourceIBMIsBareMetalServerBareMetalServerNetworkAttachmentReferenceToMap(bms.PrimaryNetworkAttachment)
		if err != nil {
			return diag.FromErr(err)
		}
		primaryNetworkAttachment = append(primaryNetworkAttachment, modelMap)
	}
	if err = d.Set("primary_network_attachment", primaryNetworkAttachment); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting primary_network_attachment %s", err))
	}

	//ni

	interfacesList := make([]map[string]interface{}, 0)
	for _, intfc := range bms.NetworkInterfaces {
		if intfc.ID != nil && *intfc.ID != *bms.PrimaryNetworkInterface.ID {
			currentNic := map[string]interface{}{}
			currentNic["id"] = *intfc.ID
			currentNic[isBareMetalServerNicHref] = *intfc.Href
			currentNic[isBareMetalServerNicName] = *intfc.Name
			if intfc.PrimaryIP != nil {
				primaryIpList := make([]map[string]interface{}, 0)
				currentIP := map[string]interface{}{}
				if intfc.PrimaryIP.Href != nil {
					currentIP[isBareMetalServerNicIpAddress] = *intfc.PrimaryIP.Address
				}
				if intfc.PrimaryIP.Href != nil {
					currentIP[isBareMetalServerNicIpHref] = *intfc.PrimaryIP.Href
				}
				if intfc.PrimaryIP.Name != nil {
					currentIP[isBareMetalServerNicIpName] = *intfc.PrimaryIP.Name
				}
				if intfc.PrimaryIP.ID != nil {
					currentIP[isBareMetalServerNicIpID] = *intfc.PrimaryIP.ID
				}
				if intfc.PrimaryIP.ResourceType != nil {
					currentIP[isBareMetalServerNicResourceType] = *intfc.PrimaryIP.ResourceType
				}
				primaryIpList = append(primaryIpList, currentIP)
				currentNic[isBareMetalServerNicPrimaryIP] = primaryIpList
			}
			getnicoptions := &vpcv1.GetBareMetalServerNetworkInterfaceOptions{
				BareMetalServerID: bms.ID,
				ID:                intfc.ID,
			}
			bmsnicintf, response, err := sess.GetBareMetalServerNetworkInterface(getnicoptions)
			if err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error getting network interfaces attached to the bare metal server %s\n%s", err, response))
			}

			switch reflect.TypeOf(bmsnicintf).String() {
			case "*vpcv1.BareMetalServerNetworkInterfaceByPci":
				{
					bmsnic := bmsnicintf.(*vpcv1.BareMetalServerNetworkInterfaceByPci)
					currentNic[isBareMetalServerNicAllowIPSpoofing] = *bmsnic.AllowIPSpoofing
					currentNic[isBareMetalServerNicSubnet] = *bmsnic.Subnet.ID
					if len(bmsnic.SecurityGroups) != 0 {
						secgrpList := []string{}
						for i := 0; i < len(bmsnic.SecurityGroups); i++ {
							secgrpList = append(secgrpList, string(*(bmsnic.SecurityGroups[i].ID)))
						}
						currentNic[isBareMetalServerNicSecurityGroups] = flex.NewStringSet(schema.HashString, secgrpList)
					}
				}
			case "*vpcv1.BareMetalServerNetworkInterfaceByVlan":
				{
					bmsnic := bmsnicintf.(*vpcv1.BareMetalServerNetworkInterfaceByVlan)
					currentNic[isBareMetalServerNicAllowIPSpoofing] = *bmsnic.AllowIPSpoofing
					currentNic[isBareMetalServerNicSubnet] = *bmsnic.Subnet.ID
					if len(bmsnic.SecurityGroups) != 0 {
						secgrpList := []string{}
						for i := 0; i < len(bmsnic.SecurityGroups); i++ {
							secgrpList = append(secgrpList, string(*(bmsnic.SecurityGroups[i].ID)))
						}
						currentNic[isBareMetalServerNicSecurityGroups] = flex.NewStringSet(schema.HashString, secgrpList)
					}
				}
			case "*vpcv1.BareMetalServerNetworkInterfaceByHiperSocket":
				{
					bmsnic := bmsnicintf.(*vpcv1.BareMetalServerNetworkInterfaceByHiperSocket)
					currentNic[isBareMetalServerNicAllowIPSpoofing] = *bmsnic.AllowIPSpoofing
					currentNic[isBareMetalServerNicSubnet] = *bmsnic.Subnet.ID
					if len(bmsnic.SecurityGroups) != 0 {
						secgrpList := []string{}
						for i := 0; i < len(bmsnic.SecurityGroups); i++ {
							secgrpList = append(secgrpList, string(*(bmsnic.SecurityGroups[i].ID)))
						}
						currentNic[isBareMetalServerNicSecurityGroups] = flex.NewStringSet(schema.HashString, secgrpList)
					}
				}
			}
			interfacesList = append(interfacesList, currentNic)
		}
	}
	d.Set(isBareMetalServerNetworkInterfaces, interfacesList)

	networkAttachments := []map[string]interface{}{}
	if bms.NetworkAttachments != nil {
		for _, modelItem := range bms.NetworkAttachments {
			if *bms.PrimaryNetworkAttachment.ID != *modelItem.ID {
				modelMap, err := dataSourceIBMIsBareMetalServerBareMetalServerNetworkAttachmentReferenceToMap(&modelItem)
				if err != nil {
					return diag.FromErr(err)
				}
				networkAttachments = append(networkAttachments, modelMap)
			}
		}
	}
	if err = d.Set("network_attachments", networkAttachments); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting network_attachments %s", err))
	}

	if err = d.Set(isBareMetalServerProfile, *bms.Profile.Name); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting profile: %s", err))
	}
	if bms.ResourceGroup != nil {
		d.Set(isBareMetalServerResourceGroup, *bms.ResourceGroup.ID)
	}
	if err = d.Set(isBareMetalServerResourceType, bms.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting resource_type: %s", err))
	}
	if err = d.Set(isBareMetalServerStatus, *bms.Status); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting status: %s", err))
	}
	statusReasonsList := make([]map[string]interface{}, 0)
	if bms.StatusReasons != nil {
		for _, sr := range bms.StatusReasons {
			currentSR := map[string]interface{}{}
			if sr.Code != nil && sr.Message != nil {
				currentSR[isBareMetalServerStatusReasonsCode] = *sr.Code
				currentSR[isBareMetalServerStatusReasonsMessage] = *sr.Message
				if sr.MoreInfo != nil {
					currentSR[isBareMetalServerStatusReasonsMoreInfo] = *sr.MoreInfo
				}
				statusReasonsList = append(statusReasonsList, currentSR)
			}
		}
	}
	d.Set(isBareMetalServerStatusReasons, statusReasonsList)

	if err = d.Set(isBareMetalServerVPC, *bms.VPC.ID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting vpc: %s", err))
	}
	if err = d.Set(isBareMetalServerZone, *bms.Zone.Name); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting zone: %s", err))
	}

	tags, err := flex.GetGlobalTagsUsingCRN(meta, *bms.CRN, "", isBareMetalServerAccessTagType)
	if err != nil {
		log.Printf(
			"[ERROR] Error on get of resource bare metal server (%s) tags: %s", d.Id(), err)
	}
	d.Set(isBareMetalServerTags, tags)

	accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *bms.CRN, "", isBareMetalServerAccessTagType)
	if err != nil {
		log.Printf(
			"[ERROR] Error on get of resource bare metal server (%s) tags: %s", d.Id(), err)
	}
	d.Set(isBareMetalServerAccessTags, accesstags)

	return nil
}

func dataSourceIBMIsBareMetalServerBareMetalServerNetworkAttachmentReferenceToMap(model *vpcv1.BareMetalServerNetworkAttachmentReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIsBareMetalServerBareMetalServerNetworkAttachmentReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	primaryIPMap, err := dataSourceIBMIsBareMetalServerReservedIPReferenceToMap(model.PrimaryIP)
	if err != nil {
		return modelMap, err
	}
	modelMap["primary_ip"] = []map[string]interface{}{primaryIPMap}
	modelMap["resource_type"] = model.ResourceType
	subnetMap, err := dataSourceIBMIsBareMetalServerSubnetReferenceToMap(model.Subnet)
	if err != nil {
		return modelMap, err
	}
	modelMap["subnet"] = []map[string]interface{}{subnetMap}
	virtualNetworkInterfaceMap, err := dataSourceIBMIsBareMetalServerVirtualNetworkInterfaceReferenceAttachmentContextToMap(model.VirtualNetworkInterface)
	if err != nil {
		return modelMap, err
	}
	modelMap["virtual_network_interface"] = []map[string]interface{}{virtualNetworkInterfaceMap}
	return modelMap, nil
}

func dataSourceIBMIsBareMetalServerVirtualNetworkInterfaceReferenceAttachmentContextToMap(model *vpcv1.VirtualNetworkInterfaceReferenceAttachmentContext) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.CRN
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}

func dataSourceIBMIsBareMetalServerBareMetalServerNetworkAttachmentReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func dataSourceIBMIsBareMetalServerReservedIPReferenceToMap(model *vpcv1.ReservedIPReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["address"] = model.Address
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIsBareMetalServerReservedIPReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}

func dataSourceIBMIsBareMetalServerReservedIPReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func dataSourceIBMIsBareMetalServerSubnetReferenceToMap(model *vpcv1.SubnetReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.CRN
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIsBareMetalServerSubnetReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}

func dataSourceIBMIsBareMetalServerSubnetReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}
