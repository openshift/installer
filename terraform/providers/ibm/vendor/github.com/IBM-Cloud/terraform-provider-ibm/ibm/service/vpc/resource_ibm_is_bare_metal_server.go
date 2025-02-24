// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isBareMetalServerAction                              = "action"
	isBareMetalServerEnableSecureBoot                    = "enable_secure_boot"
	isBareMetalServerTrustedPlatformModule               = "trusted_platform_module"
	isBareMetalServerTrustedPlatformModuleMode           = "mode"
	isBareMetalServerTrustedPlatformModuleEnabled        = "enabled"
	isBareMetalServerTrustedPlatformModuleSupportedModes = "supported_modes"
	isBareMetalServerBandwidth                           = "bandwidth"
	isBareMetalServerBootTarget                          = "boot_target"
	isBareMetalServerCreatedAt                           = "created_at"
	isBareMetalServerCPU                                 = "cpu"
	isBareMetalServerCPUArchitecture                     = "architecture"
	isBareMetalServerCPUCoreCount                        = "core_count"
	isBareMetalServerCpuSocketCount                      = "socket_count"
	isBareMetalServerCpuThreadPerCore                    = "threads_per_core"
	isBareMetalServerCRN                                 = "crn"
	isBareMetalServerDisks                               = "disks"
	isBareMetalServerDiskID                              = "id"
	isBareMetalServerDiskSize                            = "size"
	isBareMetalServerDiskName                            = "name"
	isBareMetalServerDiskInterfaceType                   = "interface_type"
	isBareMetalServerHref                                = "href"
	isBareMetalServerMemory                              = "memory"
	isBareMetalServerTags                                = "tags"
	isBareMetalServerName                                = "name"
	isBareMetalServerNetworkInterfaces                   = "network_interfaces"
	isBareMetalServerPrimaryNetworkInterface             = "primary_network_interface"
	isBareMetalServerProfile                             = "profile"
	isBareMetalServerResourceGroup                       = "resource_group"
	isBareMetalServerResourceType                        = "resource_type"
	isBareMetalServerStatus                              = "status"
	isBareMetalServerStatusReasons                       = "status_reasons"
	isBareMetalServerVPC                                 = "vpc"
	isBareMetalServerZone                                = "zone"
	isBareMetalServerStatusReasonsCode                   = "code"
	isBareMetalServerStatusReasonsMessage                = "message"
	isBareMetalServerStatusReasonsMoreInfo               = "more_info"
	isBareMetalServerDeleteType                          = "delete_type"
	isBareMetalServerImage                               = "image"
	isBareMetalServerFirmwareUpdateTypeAvailable         = "firmware_update_type_available"
	isBareMetalServerKeys                                = "keys"
	isBareMetalServerUserData                            = "user_data"
	isBareMetalServerNicName                             = "name"
	isBareMetalServerNicPortSpeed                        = "port_speed"
	isBareMetalServerNicAllowIPSpoofing                  = "allow_ip_spoofing"
	isBareMetalServerNicSecurityGroups                   = "security_groups"
	isBareMetalServerNicSubnet                           = "subnet"
	isBareMetalServerUserAccounts                        = "user_accounts"
	isBareMetalServerActionDeleting                      = "deleting"
	isBareMetalServerActionDeleted                       = "deleted"
	isBareMetalServerActionStatusStopping                = "stopping"
	isBareMetalServerActionStatusStopped                 = "stopped"
	isBareMetalServerActionStatusStarting                = "starting"
	isBareMetalServerStatusRunning                       = "running"
	isBareMetalServerStatusPending                       = "pending"
	isBareMetalServerStatusRestarting                    = "restarting"
	isBareMetalServerStatusFailed                        = "failed"
	isBareMetalServerAccessTags                          = "access_tags"
	isBareMetalServerUserTagType                         = "user"
	isBareMetalServerAccessTagType                       = "access"
)

func ResourceIBMIsBareMetalServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMISBareMetalServerCreate,
		ReadContext:   resourceIBMISBareMetalServerRead,
		UpdateContext: resourceIBMISBareMetalServerUpdate,
		DeleteContext: resourceIBMISBareMetalServerDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) (result []*schema.ResourceData, err error) {
				log.Printf("[INFO] Bare metal server (%s) importing", d.Id())
				d.Set(isBareMetalServerDeleteType, "hard")
				return []*schema.ResourceData{d}, nil
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			customdiff.Sequence(
				func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
					return flex.ResourceTagsCustomizeDiff(diff)
				},
			),
			customdiff.Sequence(
				func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
					return flex.ResourceValidateAccessTags(diff, v)
				},
			),
		),

		Schema: map[string]*schema.Schema{

			isBareMetalServerName: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_bare_metal_server", isBareMetalServerName),
				Description:  "Bare metal server name",
			},

			isBareMetalServerEnableSecureBoot: {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates whether secure boot is enabled. If enabled, the image must support secure boot or the server will fail to boot.",
			},

			isBareMetalServerTrustedPlatformModule: {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isBareMetalServerTrustedPlatformModuleMode: {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validate.InvokeValidator("ibm_is_bare_metal_server", isBareMetalServerTrustedPlatformModuleMode),
							Description:  "The trusted platform module mode to use. The specified value must be listed in the bare metal server profile's supported_trusted_platform_module_modes",
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

			isBareMetalServerAction: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_bare_metal_server", isBareMetalServerAction),
				Description:  "This restart/start/stops a bare metal server.",
			},
			isBareMetalServerBandwidth: {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The total bandwidth (in megabits per second)",
			},
			isBareMetalServerBootTarget: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for this bare metal server disk",
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
			isBareMetalServerFirmwareUpdateTypeAvailable: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of firmware update available",
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
			isBareMetalServerDeleteType: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "hard",
				Description: "Enables stopping type of the bare metal server before deleting",
			},
			isBareMetalServerPrimaryNetworkInterface: {
				Type:          schema.TypeList,
				MinItems:      1,
				MaxItems:      1,
				Optional:      true,
				Computed:      true,
				ExactlyOneOf:  []string{"primary_network_attachment", "primary_network_interface"},
				ConflictsWith: []string{"primary_network_attachment", "network_attachments"},
				Description:   "Primary Network interface info",
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
						isBareMetalServerNicEnableInfraNAT: {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "If true, the VPC infrastructure performs any needed NAT operations. If false, the packet is passed unmodified to/from the network interface, allowing the workload to perform any needed NAT operations.",
						},
						isBareMetalServerNicInterfaceType: {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validate.InvokeValidator("ibm_is_bare_metal_server", isBareMetalServerNicInterfaceType),
							Description:  "The network interface type: [ pci, hipersocket ]",
						},
						isBareMetalServerNicPrimaryIP: {
							Type:        schema.TypeList,
							Optional:    true,
							MinItems:    0,
							MaxItems:    1,
							Computed:    true,
							Description: "title: IPv4, The IP address. ",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isBareMetalServerNicIpAddress: {
										Type:          schema.TypeString,
										Optional:      true,
										Computed:      true,
										ConflictsWith: []string{"primary_network_interface.0.primary_ip.0.reserved_ip"},
										Description:   "The globally unique IP address",
									},
									isBareMetalServerNicIpHref: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this reserved IP",
									},
									isBareMetalServerNicIpAutoDelete: {
										Type:          schema.TypeBool,
										Optional:      true,
										Computed:      true,
										ConflictsWith: []string{"primary_network_interface.0.primary_ip.0.reserved_ip"},
										Description:   "Indicates whether this reserved IP member will be automatically deleted when either target is deleted, or the reserved IP is unbound.",
									},
									isBareMetalServerNicIpName: {
										Type:          schema.TypeString,
										Optional:      true,
										Computed:      true,
										ConflictsWith: []string{"primary_network_interface.0.primary_ip.0.reserved_ip"},
										Description:   "The user-defined name for this reserved IP. If unspecified, the name will be a hyphenated list of randomly-selected words. Names must be unique within the subnet the reserved IP resides in. ",
									},
									isBareMetalServerNicIpID: {
										Type:          schema.TypeString,
										Optional:      true,
										Computed:      true,
										ConflictsWith: []string{"primary_network_interface.0.primary_ip.0.address", "primary_network_interface.0.primary_ip.0.auto_delete", "primary_network_interface.0.primary_ip.0.name"},
										Description:   "Identifies a reserved IP by a unique property.",
									},
									isBareMetalServerNicResourceType: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type",
									},
								},
							},
						},
						isBareMetalServerNicAllowIPSpoofing: {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Indicates whether IP spoofing is allowed on this interface.",
						},
						isBareMetalServerNicName: {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						isBareMetalServerNicPortSpeed: {
							Type:     schema.TypeInt,
							Computed: true,
						},

						isBareMetalServerNicSecurityGroups: {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						isBareMetalServerNicSubnet: {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						isBareMetalServerNicAllowedVlans: {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Set:         schema.HashInt,
							Description: "Indicates what VLAN IDs (for VLAN type only) can use this physical (PCI type) interface. A given VLAN can only be in the allowed_vlans array for one PCI type adapter per bare metal server.",
						},
					},
				},
			},

			"primary_network_attachment": &schema.Schema{
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				Description:   "The primary network attachment.",
				ExactlyOneOf:  []string{"primary_network_attachment", "primary_network_interface"},
				ConflictsWith: []string{"primary_network_interface", "network_interfaces"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						// pna can accept either vni id or prototype
						isBareMetalServerNicAllowedVlans: {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Set:         schema.HashInt,
							Description: "Indicates what VLAN IDs (for VLAN type only) can use this physical (PCI type) interface. A given VLAN can only be in the allowed_vlans array for one PCI type adapter per bare metal server.",
						},

						isBareMetalServerNicAllowInterfaceToFloat: {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Indicates if the interface can float to any other server within the same resource_group. The interface will float automatically if the network detects a GARP or RARP on another bare metal server in the resource group. Applies only to vlan type interfaces.",
						},

						isBareMetalServerNicVlan: {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Indicates the 802.1Q VLAN ID tag that must be used for all traffic on this interface",
						},

						isBareMetalServerNicInterfaceType: {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validate.InvokeValidator("ibm_is_bare_metal_server", isBareMetalServerNicInterfaceType),
							Description:  "The network interface type: [ pci, vlan, hipersocket ]",
						},
						"name": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validate.InvokeValidator("ibm_is_bare_metal_server", "name"),
							Description:  "The name for this bare metal server network attachment. The name is unique across all network attachments for the bare metal server.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this bare metal server network attachment.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this bare metal server network attachment.",
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

						"virtual_network_interface": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Computed:    true,
							Description: "A virtual network interface for the bare metal server network attachment. This can be specified using an existing virtual network interface, or a prototype object for a new virtual network interface.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The virtual network interface id for this bare metal server network attachment.",
									},
									"allow_ip_spoofing": &schema.Schema{
										Type:          schema.TypeBool,
										Optional:      true,
										ConflictsWith: []string{"primary_network_attachment.0.virtual_network_interface.0.id"},
										Computed:      true,
										Description:   "Indicates whether source IP spoofing is allowed on this interface. If `false`, source IP spoofing is prevented on this interface. If `true`, source IP spoofing is allowed on this interface.",
									},
									"auto_delete": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Computed:    true,
										Description: "Indicates whether this virtual network interface will be automatically deleted when`target` is deleted.",
									},
									"enable_infrastructure_nat": &schema.Schema{
										Type:          schema.TypeBool,
										Optional:      true,
										ConflictsWith: []string{"primary_network_attachment.0.virtual_network_interface.0.id"},
										Computed:      true,
										Description:   "If `true`:- The VPC infrastructure performs any needed NAT operations.- `floating_ips` must not have more than one floating IP.If `false`:- Packets are passed unchanged to/from the network interface,  allowing the workload to perform any needed NAT operations.- `allow_ip_spoofing` must be `false`.- If the virtual network interface is attached:  - The target `resource_type` must be `bare_metal_server_network_attachment`.  - The target `interface_type` must not be `hipersocket`.",
									},
									"ips": &schema.Schema{
										Type:          schema.TypeSet,
										Optional:      true,
										Set:           hashIpsList,
										ConflictsWith: []string{"primary_network_attachment.0.virtual_network_interface.0.id"},
										Computed:      true,
										Description:   "The reserved IPs bound to this virtual network interface.May be empty when `lifecycle_state` is `pending`.",
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
												"auto_delete": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Indicates whether this reserved IP member will be automatically deleted when either target is deleted, or the reserved IP is unbound.",
												},
												"href": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The URL for this reserved IP.",
												},
												"reserved_ip": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
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
									"name": &schema.Schema{
										Type:          schema.TypeString,
										Optional:      true,
										ConflictsWith: []string{"primary_network_attachment.0.virtual_network_interface.0.id"},
										Computed:      true,
										ValidateFunc:  validate.InvokeValidator("ibm_is_virtual_network_interface", "vni_name"),
										Description:   "The name for this virtual network interface. The name is unique across all virtual network interfaces in the VPC.",
									},
									"protocol_state_filtering_mode": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validate.InvokeValidator("ibm_is_virtual_network_interface", "protocol_state_filtering_mode"),
										Description:  "The protocol state filtering mode used for this virtual network interface.",
									},
									"primary_ip": &schema.Schema{
										Type:          schema.TypeList,
										Optional:      true,
										MaxItems:      1,
										ConflictsWith: []string{"primary_network_attachment.0.virtual_network_interface.0.id"},
										Computed:      true,
										Description:   "The primary IP address of the virtual network interface for the bare metal server networkattachment.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"address": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
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
												"reserved_ip": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "The unique identifier for this reserved IP.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
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
									"resource_group": &schema.Schema{
										Type:          schema.TypeString,
										Optional:      true,
										ConflictsWith: []string{"primary_network_attachment.0.virtual_network_interface.0.id"},
										Computed:      true,
										Description:   "The resource group id for this virtual network interface.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
									"security_groups": {
										Type:          schema.TypeSet,
										Optional:      true,
										Computed:      true,
										ConflictsWith: []string{"primary_network_attachment.0.virtual_network_interface.0.id"},
										ForceNew:      true,
										Elem:          &schema.Schema{Type: schema.TypeString},
										Set:           schema.HashString,
										Description:   "The security groups for this virtual network interface.",
									},
									"subnet": &schema.Schema{
										Type:          schema.TypeString,
										Optional:      true,
										ConflictsWith: []string{"primary_network_attachment.0.virtual_network_interface.0.id"},
										Computed:      true,
										ForceNew:      true,
										Description:   "The associated subnet id.",
									},
								},
							},
						},
					},
				},
			},

			isBareMetalServerNetworkInterfaces: {
				Type:          schema.TypeSet,
				Optional:      true,
				Set:           resourceIBMBMSNicSet,
				ConflictsWith: []string{"primary_network_attachment", "network_attachments"},
				Computed:      true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						isBareMetalServerNicAllowIPSpoofing: {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Indicates whether IP spoofing is allowed on this interface.",
						},
						isBareMetalServerNicName: {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: flex.ApplyOnce,
							Description:      "The user-defined name for this network interface. If unspecified, the name will be a hyphenated list of randomly-selected words",
						},
						isBareMetalServerNicPortSpeed: {
							Type:     schema.TypeInt,
							Computed: true,
						},
						isBareMetalServerNicHref: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this network interface",
						},
						isBareMetalServerNicEnableInfraNAT: {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "If true, the VPC infrastructure performs any needed NAT operations. If false, the packet is passed unmodified to/from the network interface, allowing the workload to perform any needed NAT operations.",
						},
						isBareMetalServerNicInterfaceType: {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validate.InvokeValidator("ibm_is_bare_metal_server", isBareMetalServerNicInterfaceType),
							Description:  "The network interface type: [ pci, vlan, hipersocket ]",
						},
						isBareMetalServerNicPrimaryIP: {
							Type:        schema.TypeList,
							MinItems:    0,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "title: IPv4, The IP address. ",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isBareMetalServerNicIpAddress: {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The globally unique IP address",
									},
									isBareMetalServerNicIpHref: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this reserved IP",
									},
									isBareMetalServerNicIpAutoDelete: {
										Type:        schema.TypeBool,
										Optional:    true,
										Computed:    true,
										Description: "Indicates whether this reserved IP member will be automatically deleted when either target is deleted, or the reserved IP is unbound.",
									},
									isBareMetalServerNicIpName: {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The user-defined name for this reserved IP. If unspecified, the name will be a hyphenated list of randomly-selected words. Names must be unique within the subnet the reserved IP resides in. ",
									},
									isBareMetalServerNicIpID: {
										Type:        schema.TypeString,
										Optional:    true,
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
						isBareMetalServerNicSecurityGroups: {
							Type:        schema.TypeSet,
							Optional:    true,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         schema.HashString,
							Description: "Collection of security group ids",
						},
						isBareMetalServerNicSubnet: {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    false,
							Description: "The associated subnet",
						},
						isBareMetalServerNicAllowedVlans: {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Set:         schema.HashInt,
							Description: "Indicates what VLAN IDs (for VLAN type only) can use this physical (PCI type) interface. A given VLAN can only be in the allowed_vlans array for one PCI type adapter per bare metal server.",
						},

						isBareMetalServerNicAllowInterfaceToFloat: {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Indicates if the interface can float to any other server within the same resource_group. The interface will float automatically if the network detects a GARP or RARP on another bare metal server in the resource group. Applies only to vlan type interfaces.",
						},

						isBareMetalServerNicVlan: {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Indicates the 802.1Q VLAN ID tag that must be used for all traffic on this interface",
						},
					},
				},
			},

			"network_attachments": &schema.Schema{
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"primary_network_interface", "network_interfaces"},
				Description:   "The network attachments for this bare metal server, including the primary network attachment.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						// sna can accept either vni id or prototype
						isBareMetalServerNicAllowedVlans: {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Set:         schema.HashInt,
							Description: "Indicates what VLAN IDs (for VLAN type only) can use this physical (PCI type) interface. A given VLAN can only be in the allowed_vlans array for one PCI type adapter per bare metal server.",
						},

						"allow_to_float": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Indicates if the interface can float to any other server within the same resource_group. The interface will float automatically if the network detects a GARP or RARP on another bare metal server in the resource group. Applies only to vlan type interfaces.",
						},

						isBareMetalServerNicVlan: {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Indicates the 802.1Q VLAN ID tag that must be used for all traffic on this interface",
						},

						isBareMetalServerNicInterfaceType: {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validate.InvokeValidator("ibm_is_bare_metal_server", isBareMetalServerNicInterfaceType),
							Description:  "The network interface type: [ pci, vlan, hipersocket ]",
						},
						"name": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.InvokeValidator("ibm_is_bare_metal_server", "name"),
							Description:  "The name for this bare metal server network attachment. The name is unique across all network attachments for the bare metal server.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this bare metal server network attachment.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this bare metal server network attachment.",
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

						"virtual_network_interface": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "A virtual network interface for the bare metal server network attachment. This can be specified using an existing virtual network interface, or a prototype object for a new virtual network interface.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The virtual network interface id for this bare metal server network attachment.",
									},
									"allow_ip_spoofing": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Computed:    true,
										Description: "Indicates whether source IP spoofing is allowed on this interface. If `false`, source IP spoofing is prevented on this interface. If `true`, source IP spoofing is allowed on this interface.",
									},
									"auto_delete": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Computed:    true,
										Description: "Indicates whether this virtual network interface will be automatically deleted when`target` is deleted.",
									},
									"enable_infrastructure_nat": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Computed:    true,
										Description: "If `true`:- The VPC infrastructure performs any needed NAT operations.- `floating_ips` must not have more than one floating IP.If `false`:- Packets are passed unchanged to/from the network interface,  allowing the workload to perform any needed NAT operations.- `allow_ip_spoofing` must be `false`.- If the virtual network interface is attached:  - The target `resource_type` must be `bare_metal_server_network_attachment`.  - The target `interface_type` must not be `hipersocket`.",
									},
									"ips": &schema.Schema{
										Type:        schema.TypeSet,
										Optional:    true,
										Set:         hashIpsList,
										Computed:    true,
										Description: "The reserved IPs bound to this virtual network interface.May be empty when `lifecycle_state` is `pending`.",
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
												"auto_delete": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Indicates whether this reserved IP member will be automatically deleted when either target is deleted, or the reserved IP is unbound.",
												},
												"href": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The URL for this reserved IP.",
												},
												"reserved_ip": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
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
									"name": &schema.Schema{
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validate.InvokeValidator("ibm_is_virtual_network_interface", "vni_name"),
										Description:  "The name for this virtual network interface. The name is unique across all virtual network interfaces in the VPC.",
									},
									"protocol_state_filtering_mode": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validate.InvokeValidator("ibm_is_virtual_network_interface", "protocol_state_filtering_mode"),
										Description:  "The protocol state filtering mode used for this virtual network interface.",
									},
									"primary_ip": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "The primary IP address of the virtual network interface for the bare metal server networkattachment.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"address": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
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
												"reserved_ip": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "The unique identifier for this reserved IP.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
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
									"resource_group": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The resource group id for this virtual network interface.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
									"security_groups": {
										Type:        schema.TypeSet,
										Optional:    true,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Set:         schema.HashString,
										Description: "The security groups for this virtual network interface.",
									},
									"subnet": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The associated subnet id.",
									},
								},
							},
						},
					},
				},
			},

			isBareMetalServerKeys: {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "SSH key Ids for the bare metal server",
			},

			isBareMetalServerImage: {
				Type:        schema.TypeString,
				ForceNew:    false,
				Required:    true,
				Description: "image id",
			},
			isBareMetalServerProfile: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "profile name",
			},

			isBareMetalServerUserData: {
				Type:        schema.TypeString,
				ForceNew:    false,
				Optional:    true,
				Description: "User data given for the bare metal server",
			},

			isBareMetalServerZone: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone name",
			},

			isBareMetalServerVPC: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The VPC the bare metal server is to be a part of",
			},

			isBareMetalServerResourceGroup: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
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
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_bare_metal_server", "tags")},
				Set:         flex.ResourceIBMVPCHash,
				Description: "Tags for the Bare metal server",
			},

			isBareMetalServerAccessTags: {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_bare_metal_server", "accesstag")},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of access management tags",
			},
		},
	}
}

func ResourceIBMIsBareMetalServerValidator() *validate.ResourceValidator {
	bareMetalServerActions := "start, restart, stop"
	tpmModes := "disabled, tpm_2"
	interface_types := "pci, hipersocket"
	validateSchema := make([]validate.ValidateSchema, 1)

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isBareMetalServerNicInterfaceType,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			Default:                    "pci",
			AllowedValues:              interface_types})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isBareMetalServerName,
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "tags",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9:_ .-]+$`,
			MinValueLength:             1,
			MaxValueLength:             128})

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isBareMetalServerAction,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              bareMetalServerActions})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isBareMetalServerTrustedPlatformModuleMode,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              tpmModes})

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "accesstag",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([A-Za-z0-9_.-]|[A-Za-z0-9_.-][A-Za-z0-9_ .-]*[A-Za-z0-9_.-]):([A-Za-z0-9_.-]|[A-Za-z0-9_.-][A-Za-z0-9_ .-]*[A-Za-z0-9_.-])$`,
			MinValueLength:             1,
			MaxValueLength:             128})
	ibmISBareMetalServerResourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_bare_metal_server", Schema: validateSchema}
	return &ibmISBareMetalServerResourceValidator
}

func resourceIBMISBareMetalServerCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	createbmsoptions := &vpcv1.CreateBareMetalServerOptions{}
	options := &vpcv1.BareMetalServerPrototype{}
	var imageStr string
	if image, ok := d.GetOk(isBareMetalServerImage); ok {
		imageStr = image.(string)
	}

	if bandwidthIntf, ok := d.GetOk(isBareMetalServerBandwidth); ok {
		bandwidth := int64(bandwidthIntf.(int))
		options.Bandwidth = &bandwidth
	}
	// enable secure boot

	if _, ok := d.GetOkExists(isBareMetalServerEnableSecureBoot); ok {
		options.EnableSecureBoot = core.BoolPtr(d.Get(isBareMetalServerEnableSecureBoot).(bool))
	}

	// trusted_platform_module

	if _, ok := d.GetOk(isBareMetalServerTrustedPlatformModule); ok {
		trustedPlatformModuleModel, err := resourceIBMIsBareMetalServerMapToBareMetalServerTrustedPlatformModulePrototype(d.Get("trusted_platform_module.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		options.TrustedPlatformModule = trustedPlatformModuleModel
	}

	keySet := d.Get(isBareMetalServerKeys).(*schema.Set)
	if keySet.Len() != 0 {
		keyobjs := make([]vpcv1.KeyIdentityIntf, keySet.Len())
		for i, key := range keySet.List() {
			keystr := key.(string)
			keyobjs[i] = &vpcv1.KeyIdentity{
				ID: &keystr,
			}
		}
		options.Initialization = &vpcv1.BareMetalServerInitializationPrototype{
			Image: &vpcv1.ImageIdentity{
				ID: &imageStr,
			},
			Keys: keyobjs,
		}
		if userdata, ok := d.GetOk(isBareMetalServerUserData); ok {
			userdatastr := userdata.(string)
			options.Initialization.UserData = &userdatastr
		}
	}

	if name, ok := d.GetOk(isBareMetalServerName); ok {
		nameStr := name.(string)
		options.Name = &nameStr
	}

	if primnicintf, ok := d.GetOk(isBareMetalServerPrimaryNetworkInterface); ok && len(primnicintf.([]interface{})) > 0 {
		primnic := primnicintf.([]interface{})[0].(map[string]interface{})
		subnetintf, _ := primnic[isBareMetalServerNicSubnet]
		subnetintfstr := subnetintf.(string)
		var primnicobj = &vpcv1.BareMetalServerPrimaryNetworkInterfacePrototype{}
		primnicobj.Subnet = &vpcv1.SubnetIdentity{
			ID: &subnetintfstr,
		}
		name, _ := primnic[isBareMetalServerNicName]
		namestr := name.(string)
		if namestr != "" {
			primnicobj.Name = &namestr
		}

		if primaryIpIntf, ok := primnic[isBareMetalServerNicPrimaryIP]; ok && len(primaryIpIntf.([]interface{})) > 0 {
			primaryIp := primaryIpIntf.([]interface{})[0].(map[string]interface{})

			reservedIpIdOk, ok := primaryIp[isBareMetalServerNicIpID]
			if ok && reservedIpIdOk.(string) != "" {
				ipid := reservedIpIdOk.(string)
				primnicobj.PrimaryIP = &vpcv1.NetworkInterfaceIPPrototypeReservedIPIdentity{
					ID: &ipid,
				}
			} else {
				primaryip := &vpcv1.NetworkInterfaceIPPrototypeReservedIPPrototypeNetworkInterfaceContext{}

				reservedIpAddressOk, okAdd := primaryIp[isBareMetalServerNicIpAddress]
				if okAdd && reservedIpAddressOk.(string) != "" {
					reservedIpAddress := reservedIpAddressOk.(string)
					primaryip.Address = &reservedIpAddress
				}

				reservedIpNameOk, okName := primaryIp[isBareMetalServerNicIpName]
				if okName && reservedIpNameOk.(string) != "" {
					reservedIpName := reservedIpNameOk.(string)
					primaryip.Name = &reservedIpName
				}
				reservedIpAutoOk, okAuto := primaryIp[isBareMetalServerNicIpAutoDelete]
				if okAuto {
					reservedIpAuto := reservedIpAutoOk.(bool)
					primaryip.AutoDelete = &reservedIpAuto
				}
				if okAdd || okName || okAuto {
					primnicobj.PrimaryIP = primaryip
				}
			}
		}

		allowIPSpoofing, ok := primnic[isBareMetalServerNicAllowIPSpoofing]

		if ok && allowIPSpoofing != nil {
			allowIPSpoofingbool := allowIPSpoofing.(bool)
			if allowIPSpoofingbool {
				primnicobj.AllowIPSpoofing = &allowIPSpoofingbool
			}
		}
		enableInfraNATbool := true
		enableInfraNAT, ok := primnic[isBareMetalServerNicEnableInfraNAT]
		if ok && enableInfraNAT != nil {
			enableInfraNATbool = enableInfraNAT.(bool)
			primnicobj.EnableInfrastructureNat = &enableInfraNATbool
		}

		secgrpintf, ok := primnic[isBareMetalServerNicSecurityGroups]
		if ok {
			secgrpSet := secgrpintf.(*schema.Set)
			if secgrpSet.Len() != 0 {
				var secgrpobjs = make([]vpcv1.SecurityGroupIdentityIntf, secgrpSet.Len())
				for i, secgrpIntf := range secgrpSet.List() {
					secgrpIntfstr := secgrpIntf.(string)
					secgrpobjs[i] = &vpcv1.SecurityGroupIdentity{
						ID: &secgrpIntfstr,
					}
				}
				primnicobj.SecurityGroups = secgrpobjs
			}
		}

		if interfaceTypeOk, ok := primnic[isBareMetalServerNicInterfaceType]; ok && interfaceTypeOk.(string) != "" {
			interfaceType := interfaceTypeOk.(string)
			primnicobj.InterfaceType = &interfaceType
		} else if allowedVlansOk, ok := primnic[isBareMetalServerNicAllowedVlans]; ok {
			allowedVlansList := allowedVlansOk.(*schema.Set).List()

			allowedVlans := make([]int64, 0, len(allowedVlansList))
			for _, k := range allowedVlansList {
				allowedVlans = append(allowedVlans, int64(k.(int)))
			}
			primnicobj.AllowedVlans = allowedVlans
			interfaceType := "pci"
			primnicobj.InterfaceType = &interfaceType
		}
		options.PrimaryNetworkInterface = primnicobj
	}

	if _, ok := d.GetOk("primary_network_attachment"); ok {
		primarynetworkAttachmentsIntf := d.Get("primary_network_attachment")
		i := 0
		allowipspoofing := fmt.Sprintf("primary_network_attachment.0.virtual_network_interface.%d.allow_ip_spoofing", i)
		autodelete := fmt.Sprintf("primary_network_attachment.0.virtual_network_interface.%d.autodelete", i)
		enablenat := fmt.Sprintf("primary_network_attachment.0.virtual_network_interface.%d.enable_infrastructure_nat", i)
		primaryNetworkAttachmentModel, err := resourceIBMIsBareMetalServerMapToBareMetalServerPrimaryNetworkAttachmentPrototype(allowipspoofing, autodelete, enablenat, d, primarynetworkAttachmentsIntf.([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		options.PrimaryNetworkAttachment = primaryNetworkAttachmentModel
	}
	if i, ok := d.GetOk("network_attachments"); ok {
		allowipspoofing := fmt.Sprintf("network_attachments.%d.virtual_network_interface.0.allow_ip_spoofing", i)
		autodelete := fmt.Sprintf("network_attachments.%d.virtual_network_interface.0.autodelete", i)
		enablenat := fmt.Sprintf("network_attachments.%d.virtual_network_interface.0.enable_infrastructure_nat", i)
		allowfloat := fmt.Sprintf("network_attachments.%d.virtual_network_interface.0.allow_to_float", i)
		networkAttachmentsIntf := d.Get("network_attachments")
		networkAttachments := []vpcv1.BareMetalServerNetworkAttachmentPrototypeIntf{}
		for _, networkAttachmentsItem := range networkAttachmentsIntf.([]interface{}) {
			networkAttachmentsItemModel, err := resourceIBMIsBareMetalServerMapToBareMetalServerNetworkAttachmentPrototype(allowipspoofing, allowfloat, autodelete, enablenat, d, networkAttachmentsItem.(map[string]interface{}))
			if err != nil {
				return diag.FromErr(err)
			}
			networkAttachments = append(networkAttachments, networkAttachmentsItemModel)
		}
		options.NetworkAttachments = networkAttachments
	}

	if nicsintf, ok := d.GetOk(isBareMetalServerNetworkInterfaces); ok {

		nics := nicsintf.(*schema.Set).List()
		inlinenicobj := make([]vpcv1.BareMetalServerNetworkInterfacePrototypeIntf, 0)
		for _, resource := range nics {
			nic := resource.(map[string]interface{})
			interfaceType := ""

			if allowedVlansOk, ok := nic[isBareMetalServerNicAllowedVlans]; ok {
				interfaceType = "pci"
				var nicobj = &vpcv1.BareMetalServerNetworkInterfacePrototypeBareMetalServerNetworkInterfaceByPciPrototype{}
				nicobj.InterfaceType = &interfaceType

				allowedVlansList := allowedVlansOk.(*schema.Set).List()

				if len(allowedVlansList) > 0 {
					allowedVlans := make([]int64, 0, len(allowedVlansList))
					for _, k := range allowedVlansList {
						allowedVlans = append(allowedVlans, int64(k.(int)))
					}
					nicobj.AllowedVlans = allowedVlans

					subnetintf, _ := nic[isBareMetalServerNicSubnet]
					subnetintfstr := subnetintf.(string)
					nicobj.Subnet = &vpcv1.SubnetIdentity{
						ID: &subnetintfstr,
					}
					name, _ := nic[isBareMetalServerNicName]
					namestr := name.(string)
					if namestr != "" {
						nicobj.Name = &namestr
					}

					enableInfraNAT, ok := nic[isBareMetalServerNicEnableInfraNAT]
					enableInfraNATbool := enableInfraNAT.(bool)
					if ok {
						nicobj.EnableInfrastructureNat = &enableInfraNATbool
					}

					if primaryIpIntf, ok := nic[isBareMetalServerNicPrimaryIP]; ok && len(primaryIpIntf.([]interface{})) > 0 {
						primaryIp := primaryIpIntf.([]interface{})[0].(map[string]interface{})

						reservedIpIdOk, ok := primaryIp[isBareMetalServerNicIpID]
						if ok && reservedIpIdOk.(string) != "" {
							ipid := reservedIpIdOk.(string)
							nicobj.PrimaryIP = &vpcv1.NetworkInterfaceIPPrototypeReservedIPIdentity{
								ID: &ipid,
							}
						} else {
							primaryip := &vpcv1.NetworkInterfaceIPPrototypeReservedIPPrototypeNetworkInterfaceContext{}
							reservedIpAddressOk, okAdd := primaryIp[isBareMetalServerNicIpAddress]
							if okAdd && reservedIpAddressOk.(string) != "" {
								reservedIpAddress := reservedIpAddressOk.(string)
								primaryip.Address = &reservedIpAddress
							}
							reservedIpNameOk, okName := primaryIp[isBareMetalServerNicIpName]
							if okName && reservedIpNameOk.(string) != "" {
								reservedIpName := reservedIpNameOk.(string)
								primaryip.Name = &reservedIpName
							}
							reservedIpAutoOk, okAuto := primaryIp[isBareMetalServerNicIpAutoDelete]
							if okAuto {
								reservedIpAuto := reservedIpAutoOk.(bool)
								primaryip.AutoDelete = &reservedIpAuto
							}
							if okAdd || okName || okAuto {
								nicobj.PrimaryIP = primaryip
							}
						}

					}

					allowIPSpoofing, ok := nic[isBareMetalServerNicAllowIPSpoofing]
					allowIPSpoofingbool := allowIPSpoofing.(bool)
					if ok && allowIPSpoofingbool {
						nicobj.AllowIPSpoofing = &allowIPSpoofingbool
					}
					secgrpintf, ok := nic[isBareMetalServerNicSecurityGroups]
					if ok {
						secgrpSet := secgrpintf.(*schema.Set)
						if secgrpSet.Len() != 0 {
							var secgrpobjs = make([]vpcv1.SecurityGroupIdentityIntf, secgrpSet.Len())
							for i, secgrpIntf := range secgrpSet.List() {
								secgrpIntfstr := secgrpIntf.(string)
								secgrpobjs[i] = &vpcv1.SecurityGroupIdentity{
									ID: &secgrpIntfstr,
								}
							}
							nicobj.SecurityGroups = secgrpobjs
						}
					}
					inlinenicobj = append(inlinenicobj, nicobj)
				} else if interfaceTypeOk, ok := nic[isBareMetalServerNicInterfaceType]; ok && interfaceTypeOk.(string) == "hipersocket" {
					interfaceType = interfaceTypeOk.(string)
					var nicobj = &vpcv1.BareMetalServerNetworkInterfacePrototypeBareMetalServerNetworkInterfaceByHiperSocketPrototype{}
					nicobj.InterfaceType = &interfaceType

					subnetintf, _ := nic[isBareMetalServerNicSubnet]
					subnetintfstr := subnetintf.(string)
					nicobj.Subnet = &vpcv1.SubnetIdentity{
						ID: &subnetintfstr,
					}
					name, _ := nic[isBareMetalServerNicName]
					namestr := name.(string)
					if namestr != "" {
						nicobj.Name = &namestr
					}

					enableInfraNAT, ok := nic[isBareMetalServerNicEnableInfraNAT]
					enableInfraNATbool := enableInfraNAT.(bool)
					if ok {
						nicobj.EnableInfrastructureNat = &enableInfraNATbool
					}

					if primaryIpIntf, ok := nic[isBareMetalServerNicPrimaryIP]; ok && len(primaryIpIntf.([]interface{})) > 0 {
						// primaryIp := primaryIpIntf.([]interface{})[0].(map[string]interface{})
						// reservedIpAddressOk, ok := primaryIp[isBareMetalServerNicIpAddress]
						// if ok && reservedIpAddressOk.(string) != "" {
						// 	reservedIpAddress := reservedIpAddressOk.(string)
						// 	nicobj.PrimaryIpv4Address = &reservedIpAddress
						// }

						primaryIp := primaryIpIntf.([]interface{})[0].(map[string]interface{})
						reservedIpIdOk, ok := primaryIp[isBareMetalServerNicIpID]
						if ok && reservedIpIdOk.(string) != "" {
							ipid := reservedIpIdOk.(string)
							nicobj.PrimaryIP = &vpcv1.NetworkInterfaceIPPrototypeReservedIPIdentity{
								ID: &ipid,
							}
						} else {
							primaryip := &vpcv1.NetworkInterfaceIPPrototypeReservedIPPrototypeNetworkInterfaceContext{}

							reservedIpAddressOk, okAdd := primaryIp[isBareMetalServerNicIpAddress]
							if okAdd && reservedIpAddressOk.(string) != "" {
								reservedIpAddress := reservedIpAddressOk.(string)
								primaryip.Address = &reservedIpAddress
							}

							reservedIpNameOk, okName := primaryIp[isBareMetalServerNicIpName]
							if okName && reservedIpNameOk.(string) != "" {
								reservedIpName := reservedIpNameOk.(string)
								primaryip.Name = &reservedIpName
							}

							reservedIpAutoOk, okAuto := primaryIp[isBareMetalServerNicIpAutoDelete]
							if okAuto {
								reservedIpAuto := reservedIpAutoOk.(bool)
								primaryip.AutoDelete = &reservedIpAuto
							}
							if okAdd || okName || okAuto {
								nicobj.PrimaryIP = primaryip
							}
						}
					}

					allowIPSpoofing, ok := nic[isBareMetalServerNicAllowIPSpoofing]
					allowIPSpoofingbool := allowIPSpoofing.(bool)
					if ok && allowIPSpoofingbool {
						nicobj.AllowIPSpoofing = &allowIPSpoofingbool
					}
					secgrpintf, ok := nic[isBareMetalServerNicSecurityGroups]
					if ok {
						secgrpSet := secgrpintf.(*schema.Set)
						if secgrpSet.Len() != 0 {
							var secgrpobjs = make([]vpcv1.SecurityGroupIdentityIntf, secgrpSet.Len())
							for i, secgrpIntf := range secgrpSet.List() {
								secgrpIntfstr := secgrpIntf.(string)
								secgrpobjs[i] = &vpcv1.SecurityGroupIdentity{
									ID: &secgrpIntfstr,
								}
							}
							nicobj.SecurityGroups = secgrpobjs
						}
					}
					inlinenicobj = append(inlinenicobj, nicobj)
				} else {
					interfaceType = "vlan"
					var nicobj = &vpcv1.BareMetalServerNetworkInterfacePrototypeBareMetalServerNetworkInterfaceByVlanPrototype{}
					nicobj.InterfaceType = &interfaceType

					if aitf, ok := nic[isBareMetalServerNicAllowInterfaceToFloat]; ok {
						allowInterfaceToFloat := aitf.(bool)
						nicobj.AllowInterfaceToFloat = &allowInterfaceToFloat
					}
					if vlan, ok := nic[isBareMetalServerNicVlan]; ok {
						vlanInt := int64(vlan.(int))
						nicobj.Vlan = &vlanInt
					}

					subnetintf, _ := nic[isBareMetalServerNicSubnet]
					subnetintfstr := subnetintf.(string)
					nicobj.Subnet = &vpcv1.SubnetIdentity{
						ID: &subnetintfstr,
					}
					name, _ := nic[isBareMetalServerNicName]
					namestr := name.(string)
					if namestr != "" {
						nicobj.Name = &namestr
					}

					enableInfraNAT, ok := nic[isBareMetalServerNicEnableInfraNAT]
					enableInfraNATbool := enableInfraNAT.(bool)
					if ok {
						nicobj.EnableInfrastructureNat = &enableInfraNATbool
					}

					if primaryIpIntf, ok := nic[isBareMetalServerNicPrimaryIP]; ok && len(primaryIpIntf.([]interface{})) > 0 {
						primaryIp := primaryIpIntf.([]interface{})[0].(map[string]interface{})
						reservedIpIdOk, ok := primaryIp[isBareMetalServerNicIpID]
						if ok && reservedIpIdOk.(string) != "" {
							ipid := reservedIpIdOk.(string)
							nicobj.PrimaryIP = &vpcv1.NetworkInterfaceIPPrototypeReservedIPIdentity{
								ID: &ipid,
							}
						} else {
							primaryip := &vpcv1.NetworkInterfaceIPPrototypeReservedIPPrototypeNetworkInterfaceContext{}

							reservedIpAddressOk, okAdd := primaryIp[isBareMetalServerNicIpAddress]
							if okAdd && reservedIpAddressOk.(string) != "" {
								reservedIpAddress := reservedIpAddressOk.(string)
								primaryip.Address = &reservedIpAddress
							}

							reservedIpNameOk, okName := primaryIp[isBareMetalServerNicIpName]
							if okName && reservedIpNameOk.(string) != "" {
								reservedIpName := reservedIpNameOk.(string)
								primaryip.Name = &reservedIpName
							}

							reservedIpAutoOk, okAuto := primaryIp[isBareMetalServerNicIpAutoDelete]
							if okAuto {
								reservedIpAuto := reservedIpAutoOk.(bool)
								primaryip.AutoDelete = &reservedIpAuto
							}
							if okAdd || okName || okAuto {
								nicobj.PrimaryIP = primaryip
							}
						}
					}

					allowIPSpoofing, ok := nic[isBareMetalServerNicAllowIPSpoofing]
					allowIPSpoofingbool := allowIPSpoofing.(bool)
					if ok && allowIPSpoofingbool {
						nicobj.AllowIPSpoofing = &allowIPSpoofingbool
					}
					secgrpintf, ok := nic[isBareMetalServerNicSecurityGroups]
					if ok {
						secgrpSet := secgrpintf.(*schema.Set)
						if secgrpSet.Len() != 0 {
							var secgrpobjs = make([]vpcv1.SecurityGroupIdentityIntf, secgrpSet.Len())
							for i, secgrpIntf := range secgrpSet.List() {
								secgrpIntfstr := secgrpIntf.(string)
								secgrpobjs[i] = &vpcv1.SecurityGroupIdentity{
									ID: &secgrpIntfstr,
								}
							}
							nicobj.SecurityGroups = secgrpobjs
						}
					}
					inlinenicobj = append(inlinenicobj, nicobj)
				}
			} else if interfaceTypeOk, ok := nic[isBareMetalServerNicInterfaceType]; ok && interfaceTypeOk.(string) == "hipersocket" {
				interfaceType = interfaceTypeOk.(string)
				var nicobj = &vpcv1.BareMetalServerNetworkInterfacePrototypeBareMetalServerNetworkInterfaceByHiperSocketPrototype{}
				nicobj.InterfaceType = &interfaceType

				subnetintf, _ := nic[isBareMetalServerNicSubnet]
				subnetintfstr := subnetintf.(string)
				nicobj.Subnet = &vpcv1.SubnetIdentity{
					ID: &subnetintfstr,
				}
				name, _ := nic[isBareMetalServerNicName]
				namestr := name.(string)
				if namestr != "" {
					nicobj.Name = &namestr
				}

				enableInfraNAT, ok := nic[isBareMetalServerNicEnableInfraNAT]
				enableInfraNATbool := enableInfraNAT.(bool)
				if ok {
					nicobj.EnableInfrastructureNat = &enableInfraNATbool
				}

				if primaryIpIntf, ok := nic[isBareMetalServerNicPrimaryIP]; ok && len(primaryIpIntf.([]interface{})) > 0 {
					// primaryIp := primaryIpIntf.([]interface{})[0].(map[string]interface{})
					// reservedIpAddressOk, ok := primaryIp[isBareMetalServerNicIpAddress]
					// if ok && reservedIpAddressOk.(string) != "" {
					// 	reservedIpAddress := reservedIpAddressOk.(string)
					// 	nicobj.PrimaryIpv4Address = &reservedIpAddress
					// }

					primaryIp := primaryIpIntf.([]interface{})[0].(map[string]interface{})
					reservedIpIdOk, ok := primaryIp[isBareMetalServerNicIpID]
					if ok && reservedIpIdOk.(string) != "" {
						ipid := reservedIpIdOk.(string)
						nicobj.PrimaryIP = &vpcv1.NetworkInterfaceIPPrototypeReservedIPIdentity{
							ID: &ipid,
						}
					} else {
						primaryip := &vpcv1.NetworkInterfaceIPPrototypeReservedIPPrototypeNetworkInterfaceContext{}

						reservedIpAddressOk, okAdd := primaryIp[isBareMetalServerNicIpAddress]
						if okAdd && reservedIpAddressOk.(string) != "" {
							reservedIpAddress := reservedIpAddressOk.(string)
							primaryip.Address = &reservedIpAddress
						}

						reservedIpNameOk, okName := primaryIp[isBareMetalServerNicIpName]
						if okName && reservedIpNameOk.(string) != "" {
							reservedIpName := reservedIpNameOk.(string)
							primaryip.Name = &reservedIpName
						}

						reservedIpAutoOk, okAuto := primaryIp[isBareMetalServerNicIpAutoDelete]
						if okAuto {
							reservedIpAuto := reservedIpAutoOk.(bool)
							primaryip.AutoDelete = &reservedIpAuto
						}
						if okAdd || okName || okAuto {
							nicobj.PrimaryIP = primaryip
						}
					}
				}

				allowIPSpoofing, ok := nic[isBareMetalServerNicAllowIPSpoofing]
				allowIPSpoofingbool := allowIPSpoofing.(bool)
				if ok && allowIPSpoofingbool {
					nicobj.AllowIPSpoofing = &allowIPSpoofingbool
				}
				secgrpintf, ok := nic[isBareMetalServerNicSecurityGroups]
				if ok {
					secgrpSet := secgrpintf.(*schema.Set)
					if secgrpSet.Len() != 0 {
						var secgrpobjs = make([]vpcv1.SecurityGroupIdentityIntf, secgrpSet.Len())
						for i, secgrpIntf := range secgrpSet.List() {
							secgrpIntfstr := secgrpIntf.(string)
							secgrpobjs[i] = &vpcv1.SecurityGroupIdentity{
								ID: &secgrpIntfstr,
							}
						}
						nicobj.SecurityGroups = secgrpobjs
					}
				}
				inlinenicobj = append(inlinenicobj, nicobj)
			} else {
				interfaceType = "vlan"
				var nicobj = &vpcv1.BareMetalServerNetworkInterfacePrototypeBareMetalServerNetworkInterfaceByVlanPrototype{}
				nicobj.InterfaceType = &interfaceType

				if aitf, ok := nic[isBareMetalServerNicAllowInterfaceToFloat]; ok {
					allowInterfaceToFloat := aitf.(bool)
					nicobj.AllowInterfaceToFloat = &allowInterfaceToFloat
				}
				if vlan, ok := nic[isBareMetalServerNicVlan]; ok {
					vlanInt := int64(vlan.(int))
					nicobj.Vlan = &vlanInt
				}

				subnetintf, _ := nic[isBareMetalServerNicSubnet]
				subnetintfstr := subnetintf.(string)
				nicobj.Subnet = &vpcv1.SubnetIdentity{
					ID: &subnetintfstr,
				}
				name, _ := nic[isBareMetalServerNicName]
				namestr := name.(string)
				if namestr != "" {
					nicobj.Name = &namestr
				}

				enableInfraNAT, ok := nic[isBareMetalServerNicEnableInfraNAT]
				enableInfraNATbool := enableInfraNAT.(bool)
				if ok {
					nicobj.EnableInfrastructureNat = &enableInfraNATbool
				}

				if primaryIpIntf, ok := nic[isBareMetalServerNicPrimaryIP]; ok && len(primaryIpIntf.([]interface{})) > 0 {
					primaryIp := primaryIpIntf.([]interface{})[0].(map[string]interface{})
					reservedIpIdOk, ok := primaryIp[isBareMetalServerNicIpID]
					if ok && reservedIpIdOk.(string) != "" {
						ipid := reservedIpIdOk.(string)
						nicobj.PrimaryIP = &vpcv1.NetworkInterfaceIPPrototypeReservedIPIdentity{
							ID: &ipid,
						}
					} else {
						primaryip := &vpcv1.NetworkInterfaceIPPrototypeReservedIPPrototypeNetworkInterfaceContext{}

						reservedIpAddressOk, okAdd := primaryIp[isBareMetalServerNicIpAddress]
						if okAdd && reservedIpAddressOk.(string) != "" {
							reservedIpAddress := reservedIpAddressOk.(string)
							primaryip.Address = &reservedIpAddress
						}

						reservedIpNameOk, okName := primaryIp[isBareMetalServerNicIpName]
						if okName && reservedIpNameOk.(string) != "" {
							reservedIpName := reservedIpNameOk.(string)
							primaryip.Name = &reservedIpName
						}

						reservedIpAutoOk, okAuto := primaryIp[isBareMetalServerNicIpAutoDelete]
						if okAuto {
							reservedIpAuto := reservedIpAutoOk.(bool)
							primaryip.AutoDelete = &reservedIpAuto
						}
						if okAdd || okName || okAuto {
							nicobj.PrimaryIP = primaryip
						}
					}
				}

				allowIPSpoofing, ok := nic[isBareMetalServerNicAllowIPSpoofing]
				allowIPSpoofingbool := allowIPSpoofing.(bool)
				if ok && allowIPSpoofingbool {
					nicobj.AllowIPSpoofing = &allowIPSpoofingbool
				}
				secgrpintf, ok := nic[isBareMetalServerNicSecurityGroups]
				if ok {
					secgrpSet := secgrpintf.(*schema.Set)
					if secgrpSet.Len() != 0 {
						var secgrpobjs = make([]vpcv1.SecurityGroupIdentityIntf, secgrpSet.Len())
						for i, secgrpIntf := range secgrpSet.List() {
							secgrpIntfstr := secgrpIntf.(string)
							secgrpobjs[i] = &vpcv1.SecurityGroupIdentity{
								ID: &secgrpIntfstr,
							}
						}
						nicobj.SecurityGroups = secgrpobjs
					}
				}
				inlinenicobj = append(inlinenicobj, nicobj)
			}
		}
		options.NetworkInterfaces = inlinenicobj
	}

	if rgrp, ok := d.GetOk(isBareMetalServerResourceGroup); ok {
		rg := rgrp.(string)
		options.ResourceGroup = &vpcv1.ResourceGroupIdentity{
			ID: &rg,
		}
	}

	if p, ok := d.GetOk(isBareMetalServerProfile); ok {
		profile := p.(string)
		options.Profile = &vpcv1.BareMetalServerProfileIdentity{
			Name: &profile,
		}
	}

	if z, ok := d.GetOk(isBareMetalServerZone); ok {
		zone := z.(string)
		options.Zone = &vpcv1.ZoneIdentity{
			Name: &zone,
		}
	}

	if v, ok := d.GetOk(isBareMetalServerVPC); ok {
		vpc := v.(string)
		options.VPC = &vpcv1.VPCIdentity{
			ID: &vpc,
		}
	}
	createbmsoptions.BareMetalServerPrototype = options
	bms, response, err := sess.CreateBareMetalServerWithContext(context, createbmsoptions)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[DEBUG] Create bare metal server err %s\n%s", err, response))
	}
	d.SetId(*bms.ID)
	log.Printf("[INFO] Bare Metal Server : %s", *bms.ID)
	_, err = isWaitForBareMetalServerAvailable(sess, d.Id(), d.Timeout(schema.TimeoutCreate), d)
	if err != nil {
		return diag.FromErr(err)
	}
	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk(isBareMetalServerTags); ok || v != "" {
		oldList, newList := d.GetChange(isBareMetalServerTags)
		err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *bms.CRN, "", isBareMetalServerUserTagType)
		if err != nil {
			log.Printf(
				"[ERROR] Error on create of resource bare metal server (%s) tags: %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk(isBareMetalServerAccessTags); ok {
		oldList, newList := d.GetChange(isBareMetalServerAccessTags)
		err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *bms.CRN, "", isBareMetalServerAccessTagType)
		if err != nil {
			log.Printf(
				"[ERROR] Error on create of resource bare metal server (%s) access tags: %s", d.Id(), err)
		}
	}

	return resourceIBMISBareMetalServerRead(context, d, meta)
}

func resourceIBMISBareMetalServerRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := d.Id()
	err := bareMetalServerGet(context, d, meta, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func bareMetalServerGet(context context.Context, d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	options := &vpcv1.GetBareMetalServerOptions{
		ID: &id,
	}
	bms, response, err := sess.GetBareMetalServerWithContext(context, options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error getting Bare Metal Server (%s): %s\n%s", id, err, response)
	}
	d.SetId(*bms.ID)
	d.Set(isBareMetalServerBandwidth, bms.Bandwidth)
	if bms.BootTarget != nil {
		bmsBootTargetIntf := bms.BootTarget.(*vpcv1.BareMetalServerBootTarget)
		bmsBootTarget := bmsBootTargetIntf.ID
		d.Set(isBareMetalServerBootTarget, bmsBootTarget)
	}
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
	d.Set(isBareMetalServerCRN, *bms.CRN)
	if bms.Firmware != nil && bms.Firmware.Update != nil {
		d.Set(isBareMetalServerFirmwareUpdateTypeAvailable, *bms.Firmware.Update)
	}

	//enable secure boot
	if err = d.Set(isBareMetalServerEnableSecureBoot, bms.EnableSecureBoot); err != nil {
		return fmt.Errorf("[ERROR] Error setting enable_secure_boot: %s", err)
	}

	// tpm
	if bms.TrustedPlatformModule != nil {
		trustedPlatformModuleMap, err := resourceIBMIsBareMetalServerBareMetalServerTrustedPlatformModulePrototypeToMap(bms.TrustedPlatformModule)
		if err != nil {
			return (err)
		}
		if err = d.Set(isBareMetalServerTrustedPlatformModule, []map[string]interface{}{trustedPlatformModuleMap}); err != nil {
			return fmt.Errorf("[ERROR] Error setting trusted_platform_module: %s", err)
		}
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
	d.Set(isBareMetalServerHref, *bms.Href)
	d.Set(isBareMetalServerMemory, *bms.Memory)
	d.Set(isBareMetalServerName, *bms.Name)

	// get initialization
	getBmsInitialization := &vpcv1.GetBareMetalServerInitializationOptions{
		ID: bms.ID,
	}
	bmsinitialization, response, err := sess.GetBareMetalServerInitializationWithContext(context, getBmsInitialization)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error getting Bare Metal Server (%s) initialization: %s\n%s", id, err, response)
	}
	if bmsinitialization != nil && bmsinitialization.Image.ID != nil {
		d.Set(isBareMetalServerImage, *bmsinitialization.Image.ID)
	}
	if bmsinitialization != nil && bmsinitialization.Keys != nil {
		keyList := []string{}
		if len(bmsinitialization.Keys) != 0 {
			for i := 0; i < len(bmsinitialization.Keys); i++ {
				keyList = append(keyList, string(*(bmsinitialization.Keys[i].ID)))
			}
		}
		d.Set(isBareMetalServerKeys, keyList)
	}

	//pni

	if bms.PrimaryNetworkInterface != nil {
		primaryNicList := make([]map[string]interface{}, 0)
		currentPrimNic := map[string]interface{}{}
		currentPrimNic["id"] = *bms.PrimaryNetworkInterface.ID
		currentPrimNic[isBareMetalServerNicName] = *bms.PrimaryNetworkInterface.Name
		currentPrimNic[isBareMetalServerNicHref] = *bms.PrimaryNetworkInterface.Href
		currentPrimNic[isBareMetalServerNicSubnet] = *bms.PrimaryNetworkInterface.Subnet.ID
		getnicoptions := &vpcv1.GetBareMetalServerNetworkInterfaceOptions{
			BareMetalServerID: &id,
			ID:                bms.PrimaryNetworkInterface.ID,
		}
		bmsnic, response, err := sess.GetBareMetalServerNetworkInterfaceWithContext(context, getnicoptions)

		if err != nil {
			return fmt.Errorf("[ERROR] Error getting primary network interface attached to the bare metal server %s\n%s", err, response)
		}

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

			getripoptions := &vpcv1.GetSubnetReservedIPOptions{
				SubnetID: bms.PrimaryNetworkInterface.Subnet.ID,
				ID:       bms.PrimaryNetworkInterface.PrimaryIP.ID,
			}
			bmsRip, response, err := sess.GetSubnetReservedIP(getripoptions)
			if err != nil {
				return fmt.Errorf("[ERROR] Error getting network interface reserved ip(%s) attached to the bare metal server primary network interface(%s): %s\n%s", *bms.PrimaryNetworkInterface.PrimaryIP.ID, *bms.PrimaryNetworkInterface.ID, err, response)
			}
			currentIP[isBareMetalServerNicIpAutoDelete] = bmsRip.AutoDelete
		}
		switch reflect.TypeOf(bmsnic).String() {
		case "*vpcv1.BareMetalServerNetworkInterfaceByPci":
			{
				primNic := bmsnic.(*vpcv1.BareMetalServerNetworkInterfaceByPci)
				currentPrimNic[isBareMetalServerNicAllowIPSpoofing] = *primNic.AllowIPSpoofing
				currentPrimNic[isBareMetalServerNicEnableInfraNAT] = *primNic.EnableInfrastructureNat
				currentPrimNic[isBareMetalServerNicPortSpeed] = *primNic.PortSpeed
				currentPrimNic[isBareMetalServerNicInterfaceType] = *primNic.InterfaceType

				if len(primNic.SecurityGroups) != 0 {
					secgrpList := []string{}
					for i := 0; i < len(primNic.SecurityGroups); i++ {
						secgrpList = append(secgrpList, string(*(primNic.SecurityGroups[i].ID)))
					}
					currentPrimNic[isBareMetalServerNicSecurityGroups] = flex.NewStringSet(schema.HashString, secgrpList)
				}

				if primNic.AllowedVlans != nil {
					var out = make([]interface{}, len(primNic.AllowedVlans), len(primNic.AllowedVlans))
					for i, v := range primNic.AllowedVlans {
						out[i] = int(v)
					}
					currentPrimNic[isBareMetalServerNicAllowedVlans] = schema.NewSet(schema.HashInt, out)
				}
			}
		case "*vpcv1.BareMetalServerNetworkInterfaceByVlan":
			{
				primNic := bmsnic.(*vpcv1.BareMetalServerNetworkInterfaceByVlan)
				currentPrimNic[isBareMetalServerNicAllowIPSpoofing] = *primNic.AllowIPSpoofing
				currentPrimNic[isBareMetalServerNicEnableInfraNAT] = *primNic.EnableInfrastructureNat
				currentPrimNic[isBareMetalServerNicInterfaceType] = *primNic.InterfaceType
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
				currentPrimNic[isBareMetalServerNicEnableInfraNAT] = *primNic.EnableInfrastructureNat
				currentPrimNic[isBareMetalServerNicInterfaceType] = *primNic.InterfaceType
				currentPrimNic[isBareMetalServerNicPortSpeed] = *primNic.PortSpeed

				if len(primNic.SecurityGroups) != 0 {
					secgrpList := []string{}
					for i := 0; i < len(primNic.SecurityGroups); i++ {
						secgrpList = append(secgrpList, string(*(primNic.SecurityGroups[i].ID)))
					}
					currentPrimNic[isBareMetalServerNicSecurityGroups] = flex.NewStringSet(schema.HashString, secgrpList)
				}
			}
		}

		primaryNicList = append(primaryNicList, currentPrimNic)
		d.Set(isBareMetalServerPrimaryNetworkInterface, primaryNicList)
	}

	if !core.IsNil(bms.PrimaryNetworkAttachment) {
		pnaId := *bms.PrimaryNetworkAttachment.ID
		getBareMetalServerNetworkAttachment := &vpcv1.GetBareMetalServerNetworkAttachmentOptions{
			BareMetalServerID: &id,
			ID:                &pnaId,
		}
		pna, response, err := sess.GetBareMetalServerNetworkAttachment(getBareMetalServerNetworkAttachment)
		if err != nil {
			return fmt.Errorf("[ERROR] Error on GetBareMetalServerNetworkAttachment in bms : %s\n%s", err, response)
		}
		primaryNetworkAttachmentMap, err := resourceIBMIsBareMetalServerBareMetalServerNetworkAttachmentReferenceToMap(bms.PrimaryNetworkAttachment, pna, sess)
		if err != nil {
			return err
		}
		if err = d.Set("primary_network_attachment", []map[string]interface{}{primaryNetworkAttachmentMap}); err != nil {
			return fmt.Errorf("[ERROR] Error setting primary_network_attachment: %s", err)
		}
	}

	//ni
	if bms.NetworkInterfaces != nil {
		interfacesList := make([]map[string]interface{}, 0)
		for _, intfc := range bms.NetworkInterfaces {
			flagAllowFloat := false
			if *intfc.ID != *bms.PrimaryNetworkInterface.ID {
				currentNic := map[string]interface{}{}
				subnetId := *intfc.Subnet.ID
				ripId := ""
				nicId := *intfc.ID
				currentNic["id"] = nicId
				currentNic[isBareMetalServerNicName] = *intfc.Name
				currentNic[isBareMetalServerNicHref] = *intfc.Href
				currentNic[isBareMetalServerNicSubnet] = subnetId

				getnicoptions := &vpcv1.GetBareMetalServerNetworkInterfaceOptions{
					BareMetalServerID: &id,
					ID:                &nicId,
				}
				bmsnicintf, response, err := sess.GetBareMetalServerNetworkInterfaceWithContext(context, getnicoptions)
				if err != nil {
					return fmt.Errorf("[ERROR] Error getting network interface(%s) attached to the bare metal server(%s) %s\n%s", nicId, id, err, response)
				}
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
						ripId = *intfc.PrimaryIP.ID
						currentIP[isBareMetalServerNicIpID] = ripId
						getripoptions := &vpcv1.GetSubnetReservedIPOptions{
							SubnetID: &subnetId,
							ID:       &ripId,
						}
						bmsRip, response, err := sess.GetSubnetReservedIP(getripoptions)
						if err != nil {
							return fmt.Errorf("[ERROR] Error getting network interface reserved ip(%s) attached to the bare metal server network interface(%s): %s\n%s", ripId, nicId, err, response)
						}
						if bmsRip.AutoDelete != nil {
							currentIP[isBareMetalServerNicIpAutoDelete] = *bmsRip.AutoDelete
						}
					}
					if intfc.PrimaryIP.ResourceType != nil {
						currentIP[isBareMetalServerNicResourceType] = *intfc.PrimaryIP.ResourceType
					}
					primaryIpList = append(primaryIpList, currentIP)
					currentNic[isBareMetalServerNicPrimaryIP] = primaryIpList
				}

				switch reflect.TypeOf(bmsnicintf).String() {
				case "*vpcv1.BareMetalServerNetworkInterfaceByPci":
					{
						bmsnic := bmsnicintf.(*vpcv1.BareMetalServerNetworkInterfaceByPci)
						currentNic[isBareMetalServerNicAllowIPSpoofing] = *bmsnic.AllowIPSpoofing
						currentNic[isBareMetalServerNicEnableInfraNAT] = *bmsnic.EnableInfrastructureNat
						currentNic[isBareMetalServerNicPortSpeed] = *bmsnic.PortSpeed
						currentNic[isBareMetalServerNicInterfaceType] = *bmsnic.InterfaceType
						if len(bmsnic.SecurityGroups) != 0 {
							secgrpList := []string{}
							for i := 0; i < len(bmsnic.SecurityGroups); i++ {
								secgrpList = append(secgrpList, string(*(bmsnic.SecurityGroups[i].ID)))
							}
							currentNic[isBareMetalServerNicSecurityGroups] = flex.NewStringSet(schema.HashString, secgrpList)
						}
						if bmsnic.AllowedVlans != nil {
							var out = make([]interface{}, len(bmsnic.AllowedVlans), len(bmsnic.AllowedVlans))
							for i, v := range bmsnic.AllowedVlans {
								out[i] = int(v)
							}
							currentNic[isBareMetalServerNicAllowedVlans] = schema.NewSet(schema.HashInt, out)
						}
					}
				case "*vpcv1.BareMetalServerNetworkInterfaceByVlan":
					{
						bmsnic := bmsnicintf.(*vpcv1.BareMetalServerNetworkInterfaceByVlan)
						if bmsnic.AllowInterfaceToFloat != nil {
							flagAllowFloat = *bmsnic.AllowInterfaceToFloat
						}
						currentNic[isBareMetalServerNicAllowIPSpoofing] = *bmsnic.AllowIPSpoofing
						currentNic[isBareMetalServerNicEnableInfraNAT] = *bmsnic.EnableInfrastructureNat
						currentNic[isBareMetalServerNicPortSpeed] = *bmsnic.PortSpeed
						currentNic[isBareMetalServerNicInterfaceType] = *bmsnic.InterfaceType
						if bmsnic.Vlan != nil {
							currentNic[isBareMetalServerNicVlan] = *bmsnic.Vlan
						}
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
						currentNic[isBareMetalServerNicEnableInfraNAT] = *bmsnic.EnableInfrastructureNat
						currentNic[isBareMetalServerNicPortSpeed] = *bmsnic.PortSpeed
						currentNic[isBareMetalServerNicInterfaceType] = *bmsnic.InterfaceType
						if len(bmsnic.SecurityGroups) != 0 {
							secgrpList := []string{}
							for i := 0; i < len(bmsnic.SecurityGroups); i++ {
								secgrpList = append(secgrpList, string(*(bmsnic.SecurityGroups[i].ID)))
							}
							currentNic[isBareMetalServerNicSecurityGroups] = flex.NewStringSet(schema.HashString, secgrpList)
						}
					}
				}
				if !flagAllowFloat {
					interfacesList = append(interfacesList, currentNic)
				}
			}
		}
		d.Set(isBareMetalServerNetworkInterfaces, interfacesList)
	}

	if !core.IsNil(bms.NetworkAttachments) {
		networkAttachments := []map[string]interface{}{}
		for _, networkAttachmentsItem := range bms.NetworkAttachments {
			naId := *networkAttachmentsItem.ID
			if *bms.PrimaryNetworkAttachment.ID != naId {
				getBareMetalServerNetworkAttachment := &vpcv1.GetBareMetalServerNetworkAttachmentOptions{
					BareMetalServerID: &id,
					ID:                &naId,
				}
				na, response, err := sess.GetBareMetalServerNetworkAttachment(getBareMetalServerNetworkAttachment)
				if err != nil {
					return fmt.Errorf("[ERROR] Error on GetBareMetalServerNetworkAttachment in baremetal server : %s\n%s", err, response)
				}
				networkAttachmentsItemMap, err := resourceIBMIsBareMetalServerBareMetalServerNetworkAttachmentReferenceToMap(&networkAttachmentsItem, na, sess)
				if err != nil {
					return err
				}
				networkAttachments = append(networkAttachments, networkAttachmentsItemMap)
			}
		}
		if err = d.Set("network_attachments", networkAttachments); err != nil {
			return fmt.Errorf("[ERROR] Error setting network_attachments: %s", err)
		}
	}

	d.Set(isBareMetalServerProfile, *bms.Profile.Name)
	if bms.ResourceGroup != nil {
		d.Set(isBareMetalServerResourceGroup, *bms.ResourceGroup.ID)
	}
	d.Set(isBareMetalServerResourceType, bms.ResourceType)
	d.Set(isBareMetalServerStatus, *bms.Status)
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
	d.Set(isBareMetalServerVPC, *bms.VPC.ID)
	d.Set(isBareMetalServerZone, *bms.Zone.Name)

	tags, err := flex.GetGlobalTagsUsingCRN(meta, *bms.CRN, "", isBareMetalServerUserTagType)
	if err != nil {
		log.Printf(
			"[ERROR] Error on get of resource bare metal server (%s) tags: %s", d.Id(), err)
	}
	d.Set(isBareMetalServerTags, tags)

	accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *bms.CRN, "", isBareMetalServerAccessTagType)
	if err != nil {
		log.Printf(
			"[ERROR] Error on get of resource bare metal server (%s) access tags: %s", d.Id(), err)
	}
	d.Set(isBareMetalServerAccessTags, accesstags)

	return nil
}

func resourceIBMISBareMetalServerUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	id := d.Id()

	err := bareMetalServerUpdate(context, d, meta, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceIBMISBareMetalServerRead(context, d, meta)
}

func bareMetalServerUpdate(context context.Context, d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	if d.HasChange("image") || d.HasChange("keys") || d.HasChange("user_data") {
		stopServerIfStartingForInitialization := false
		newImageId := d.Get("image").(string)
		initializationPatch := &vpcv1.ReplaceBareMetalServerInitializationOptions{
			ID: &id,
			Image: &vpcv1.ImageIdentityByID{
				ID: &newImageId,
			},
		}
		// apply the user data file, if its not updated use the existing
		newUserData := d.Get("user_data").(string)
		initializationPatch.UserData = &newUserData
		// apply the keys, if its not updated use the existing
		keySet := d.Get(isBareMetalServerKeys).(*schema.Set)
		if keySet.Len() != 0 {
			keyobjs := make([]vpcv1.KeyIdentityIntf, keySet.Len())
			for i, key := range keySet.List() {
				keystr := key.(string)
				keyobjs[i] = &vpcv1.KeyIdentity{
					ID: &keystr,
				}
			}
			initializationPatch.Keys = keyobjs
		}

		stopServerIfStartingForInitialization, err = resourceStopServerIfRunning(id, "hard", d, context, sess, stopServerIfStartingForInitialization)
		if err != nil {
			return err
		}
		_, res, err := sess.ReplaceBareMetalServerInitialization(initializationPatch)
		if err != nil {
			return fmt.Errorf("ReplaceBareMetalServerInitialization failed %s\n%s", err, res)
		}
		_, err = isWaitForBareMetalServerStoppedOnReload(sess, d.Id(), d.Timeout(schema.TimeoutUpdate), d)
		if err != nil {
			return err
		}
		if stopServerIfStartingForInitialization {
			_, err = resourceStartServerIfStopped(id, "hard", d, context, sess, stopServerIfStartingForInitialization)
			if err != nil {
				return err
			}
		}
	}
	isServerStopped := false

	// network attachments

	// primary network attachment

	if d.HasChange("primary_network_attachment") && !d.IsNewResource() {
		nameChanged := d.HasChange("primary_network_attachment.0.name")
		avChanged := d.HasChange("primary_network_attachment.0.allowed_vlans")
		vniChanged := d.HasChange("primary_network_attachment.0.virtual_network_interface")
		if nameChanged || avChanged {
			pnacId := d.Get("primary_network_attachment.0.id").(string)
			updateBareMetalServerNetworkAttachmentOptions := &vpcv1.UpdateBareMetalServerNetworkAttachmentOptions{}
			updateBareMetalServerNetworkAttachmentOptions.SetBareMetalServerID(d.Id())
			updateBareMetalServerNetworkAttachmentOptions.SetID(pnacId)
			patchVals := &vpcv1.BareMetalServerNetworkAttachmentPatch{}
			if avChanged {
				var allowedVlans []int64
				for _, v := range d.Get("primary_network_attachment.0.allowed_vlans").(*schema.Set).List() {
					allowedVlansItem := int64(v.(int))
					allowedVlans = append(allowedVlans, allowedVlansItem)
				}
				patchVals.AllowedVlans = allowedVlans
			}
			if nameChanged {
				newName := d.Get("primary_network_attachment.0.name").(string)
				patchVals.Name = &newName
			}
			updateBareMetalServerNetworkAttachmentOptions.BareMetalServerNetworkAttachmentPatch, _ = patchVals.AsPatch()
			_, response, err := sess.UpdateBareMetalServerNetworkAttachmentWithContext(context, updateBareMetalServerNetworkAttachmentOptions)
			if err != nil {
				log.Printf("[DEBUG] UpdateBareMetalServerNetworkAttachmentWithContext failed %s\n%s", err, response)
				return fmt.Errorf("UpdateBareMetalServerNetworkAttachmentWithContext failed %s\n%s", err, response)
			}
		}
		if vniChanged {
			vniId := d.Get("primary_network_attachment.0.virtual_network_interface.0.id").(string)
			updateVirtualNetworkInterfaceOptions := &vpcv1.UpdateVirtualNetworkInterfaceOptions{
				ID: &vniId,
			}
			virtualNetworkInterfacePatch := &vpcv1.VirtualNetworkInterfacePatch{}
			if d.HasChange("primary_network_attachment.0.virtual_network_interface.0.auto_delete") {
				autodelete := d.Get("primary_network_attachment.0.virtual_network_interface.0.auto_delete").(bool)
				virtualNetworkInterfacePatch.AutoDelete = &autodelete
			}
			if d.HasChange("primary_network_attachment.0.virtual_network_interface.0.name") {
				name := d.Get("primary_network_attachment.0.virtual_network_interface.0.name").(string)
				virtualNetworkInterfacePatch.Name = &name
			}
			if d.HasChange("primary_network_attachment.0.virtual_network_interface.0.enable_infrastructure_nat") {
				enableNat := d.Get("primary_network_attachment.0.virtual_network_interface.0.enable_infrastructure_nat").(bool)
				virtualNetworkInterfacePatch.EnableInfrastructureNat = &enableNat
			}
			if d.HasChange("primary_network_attachment.0.virtual_network_interface.0.allow_ip_spoofing") {
				allIpSpoofing := d.Get("primary_network_attachment.0.virtual_network_interface.0.allow_ip_spoofing").(bool)
				virtualNetworkInterfacePatch.AllowIPSpoofing = &allIpSpoofing
			}
			if d.HasChange("primary_network_attachment.0.virtual_network_interface.0.protocol_state_filtering_mode") {
				psfMode := d.Get("primary_network_attachment.0.virtual_network_interface.0.protocol_state_filtering_mode").(string)
				virtualNetworkInterfacePatch.ProtocolStateFilteringMode = &psfMode
			}
			virtualNetworkInterfacePatchAsPatch, err := virtualNetworkInterfacePatch.AsPatch()
			if err != nil {
				return fmt.Errorf("[ERROR] Error encountered while apply as patch for virtualNetworkInterfacePatch of BareMetalServer(%s) vni (%s) %s", d.Id(), vniId, err)
			}
			updateVirtualNetworkInterfaceOptions.VirtualNetworkInterfacePatch = virtualNetworkInterfacePatchAsPatch
			_, response, err := sess.UpdateVirtualNetworkInterfaceWithContext(context, updateVirtualNetworkInterfaceOptions)
			if err != nil {
				log.Printf("[DEBUG] UpdateVirtualNetworkInterfaceWithContext failed %s\n%s", err, response)
				return fmt.Errorf("UpdateVirtualNetworkInterfaceWithContext failed during BareMetalServer(%s) network attachment patch %s\n%s", d.Id(), err, response)
			}

			if d.HasChange("primary_network_attachment.0.virtual_network_interface.0.ips") {
				oldips, newips := d.GetChange("primary_network_attachment.0.virtual_network_interface.0.ips")
				os := oldips.(*schema.Set)
				ns := newips.(*schema.Set)
				var oldset, newset *schema.Set

				var out = make([]interface{}, ns.Len(), ns.Len())
				for i, nA := range ns.List() {
					newPack := nA.(map[string]interface{})
					out[i] = newPack["reserved_ip"].(string)
				}
				newset = schema.NewSet(schema.HashString, out)

				out = make([]interface{}, os.Len(), os.Len())
				for i, oA := range os.List() {
					oldPack := oA.(map[string]interface{})
					out[i] = oldPack["reserved_ip"].(string)
				}
				oldset = schema.NewSet(schema.HashString, out)

				remove := flex.ExpandStringList(oldset.Difference(newset).List())
				add := flex.ExpandStringList(newset.Difference(oldset).List())

				if add != nil && len(add) > 0 {
					for _, ipItem := range add {
						if ipItem != "" {

							addVirtualNetworkInterfaceIPOptions := &vpcv1.AddVirtualNetworkInterfaceIPOptions{}
							addVirtualNetworkInterfaceIPOptions.SetVirtualNetworkInterfaceID(vniId)
							addVirtualNetworkInterfaceIPOptions.SetID(ipItem)
							_, response, err := sess.AddVirtualNetworkInterfaceIPWithContext(context, addVirtualNetworkInterfaceIPOptions)
							if err != nil {
								log.Printf("[DEBUG] AddVirtualNetworkInterfaceIPWithContext failed in VirtualNetworkInterface patch during BareMetalServer nac patch %s\n%s", err, response)
								return fmt.Errorf("AddVirtualNetworkInterfaceIPWithContext failed in VirtualNetworkInterface patch during BareMetalServer nac patch %s\n%s", err, response)
							}
						}
					}
				}
				if remove != nil && len(remove) > 0 {
					for _, ipItem := range remove {
						if ipItem != "" {

							removeVirtualNetworkInterfaceIPOptions := &vpcv1.RemoveVirtualNetworkInterfaceIPOptions{}
							removeVirtualNetworkInterfaceIPOptions.SetVirtualNetworkInterfaceID(vniId)
							removeVirtualNetworkInterfaceIPOptions.SetID(ipItem)
							response, err := sess.RemoveVirtualNetworkInterfaceIPWithContext(context, removeVirtualNetworkInterfaceIPOptions)
							if err != nil {
								log.Printf("[DEBUG] RemoveVirtualNetworkInterfaceIPWithContext failed in VirtualNetworkInterface patch during BareMetalServer nac patch %s\n%s", err, response)
								return fmt.Errorf("RemoveVirtualNetworkInterfaceIPWithContext failed in VirtualNetworkInterface patch during BareMetalServer nac patch %s\n%s", err, response)
							}
						}
					}
				}

			}
			if d.HasChange("primary_network_attachment.0.virtual_network_interface.0.primary_ip") {
				subnetId := d.Get("primary_network_attachment.0.virtual_network_interface.0.subnet").(string)
				ripId := d.Get("primary_network_attachment.0.virtual_network_interface.0.primary_ip.0.reserved_ip").(string)
				updateripoptions := &vpcv1.UpdateSubnetReservedIPOptions{
					SubnetID: &subnetId,
					ID:       &ripId,
				}
				reservedIpPath := &vpcv1.ReservedIPPatch{}
				if d.HasChange("primary_network_attachment.0.virtual_network_interface.0.primary_ip.0.name") {
					name := d.Get("primary_network_attachment.0.virtual_network_interface.0.primary_ip.0.name").(string)
					reservedIpPath.Name = &name
				}
				if d.HasChange("primary_network_attachment.0.virtual_network_interface.0.primary_ip.0.auto_delete") {
					auto := d.Get("primary_network_attachment.0.virtual_network_interface.0.primary_ip.0.auto_delete").(bool)
					reservedIpPath.AutoDelete = &auto
				}
				reservedIpPathAsPatch, err := reservedIpPath.AsPatch()
				if err != nil {
					return fmt.Errorf("[ERROR] Error calling reserved ip as patch on vni patch \n%s", err)
				}
				updateripoptions.ReservedIPPatch = reservedIpPathAsPatch
				_, response, err := sess.UpdateSubnetReservedIP(updateripoptions)
				if err != nil {
					return fmt.Errorf("[ERROR] Error updating vni reserved ip(%s): %s\n%s", ripId, err, response)
				}
			}
			if d.HasChange("primary_network_attachment.0.virtual_network_interface.0.security_groups") {
				ovs, nvs := d.GetChange("primary_network_attachment.0.virtual_network_interface.0.security_groups")
				ov := ovs.(*schema.Set)
				nv := nvs.(*schema.Set)
				remove := flex.ExpandStringList(ov.Difference(nv).List())
				add := flex.ExpandStringList(nv.Difference(ov).List())
				if len(add) > 0 {
					for i := range add {
						createsgnicoptions := &vpcv1.CreateSecurityGroupTargetBindingOptions{
							SecurityGroupID: &add[i],
							ID:              &vniId,
						}
						_, response, err := sess.CreateSecurityGroupTargetBinding(createsgnicoptions)
						if err != nil {
							return (fmt.Errorf("[ERROR] Error while creating security group %q for virtual network interface %s\n%s: %q", add[i], d.Id(), err, response))
						}
						_, err = isWaitForVirtualNetworkInterfaceAvailable(sess, vniId, d.Timeout(schema.TimeoutUpdate))
						if err != nil {
							return (err)
						}
					}

				}
				if len(remove) > 0 {
					for i := range remove {
						deletesgnicoptions := &vpcv1.DeleteSecurityGroupTargetBindingOptions{
							SecurityGroupID: &remove[i],
							ID:              &vniId,
						}
						response, err := sess.DeleteSecurityGroupTargetBinding(deletesgnicoptions)
						if err != nil {
							return (fmt.Errorf("[ERROR] Error while removing security group %q for virtual network interface %s\n%s: %q", remove[i], d.Id(), err, response))
						}
						_, err = isWaitForVirtualNetworkInterfaceAvailable(sess, vniId, d.Timeout(schema.TimeoutUpdate))
						if err != nil {
							return (err)
						}
					}
				}
			}

		}
	}

	// network attachments

	if d.HasChange("network_attachments") && !d.IsNewResource() {
		// nacs := d.Get("network_attachments").([]interface{})
		ots, nts := d.GetChange("network_attachments")
		otsIntf := ots.([]interface{})
		ntsIntf := nts.([]interface{})

		// out := make([]string, len(otsIntf))
		listToRemove, listToAdd, serverToStop, listToUpdate := findNetworkAttachmentDifferences(otsIntf, ntsIntf, d.Id(), sess, d)

		if listToUpdate != nil {
			return fmt.Errorf("[ERROR] Error while updating network attachment BareMetalServer(%s) \n%s", d.Id(), err)
		}
		serverStopped := false
		if serverToStop {
			// stop the server
			serverStopped = true
			isServerStopped, err = resourceStopServerIfRunning(id, "hard", d, context, sess, isServerStopped)
			if err != nil {
				return err
			}
		}
		for _, removeItem := range listToRemove {
			res, err := sess.DeleteBareMetalServerNetworkAttachment(&removeItem)
			if err != nil {
				return fmt.Errorf("[ERROR] Error while removing network attachment(%s) of BareMetalServer(%s) \n%s: %q", *removeItem.ID, d.Id(), err, res)
			}
		}
		for _, addItem := range listToAdd {
			_, res, err := sess.CreateBareMetalServerNetworkAttachment(&addItem)
			if err != nil {
				return fmt.Errorf("[ERROR] Error while adding network attachment(%s) of BareMetalServer(%s) \n%s: %q", *addItem.BareMetalServerID, d.Id(), err, res)
			}
		}
		if serverStopped && isServerStopped {
			// retstart ther server
			isServerStopped, err = resourceStartServerIfStopped(id, "hard", d, context, sess, isServerStopped)
			if err != nil {
				return err
			}
		}
		// j := 0
		// for _, currOtsG := range otsIntf {
		// 	currOts := currOtsG.(map[string]interface{})
		// 	flag := false
		// 	for _, currNtsG := range ntsIntf {
		// 		currNts := currNtsG.(map[string]interface{})
		// 		if currOts["id"].(string) == currNts["id"].(string) {
		// 			flag = true
		// 		}
		// 	}
		// 	if !flag {
		// 		log.Printf("[INFO] Nac with name (%s) will be deleted", currOts["name"].(string))
		// 		nacId := currOts["id"]
		// 		if nacId != nil && nacId.(string) != "" {
		// 			nacIdStr := nacId.(string)
		// 			if !containsNacId(out, nacIdStr) {
		// 				out[j] = nacIdStr
		// 				j = j + 1
		// 				deleteBareMetalServerNetworkAttachmentOptions := &vpcv1.DeleteBareMetalServerNetworkAttachmentOptions{
		// 					BareMetalServerID: &id,
		// 					ID:                &nacIdStr,
		// 				}
		// 				res, err := sess.DeleteBareMetalServerNetworkAttachment(deleteBareMetalServerNetworkAttachmentOptions)
		// 				if err != nil {
		// 					return fmt.Errorf("[ERROR] Error while deleting network attachment(%s) of BareMetalServer(%s) \n%s: %q", nacIdStr, d.Id(), err, res)
		// 				}
		// 			}
		// 		}
		// 	}
		// }

		// for i, nac := range nacs {
		// 	nacIdKey := fmt.Sprintf("network_attachments.%d.id", i)
		// 	nacId := d.Get(nacIdKey).(string)
		// 	// if nacId is empty, then create
		// 	// if nacId == "" || containsNacId(out, nacId) {

		// 	if nacId == "" {
		// 		log.Printf("[DEBUG] nacId is empty")
		// 		allowipspoofing := fmt.Sprintf("network_attachments.%d.virtual_network_interface.0.allow_ip_spoofing", i)
		// 		autodelete := fmt.Sprintf("network_attachments.%d.virtual_network_interface.0.auto_delete", i)
		// 		enablenat := fmt.Sprintf("network_attachments.%d.virtual_network_interface.0.enable_infrastructure_nat", i)
		// 		nacMap := nac.(map[string]interface{})
		// 		VirtualNetworkInterfaceModel, err := resourceIBMIsBareMetalServerMapToVirtualNetworkInterfacePrototypeAttachmentContext(allowipspoofing, autodelete, enablenat, d, nacMap["virtual_network_interface"].([]interface{})[0].(map[string]interface{}))
		// 		if err != nil {
		// 			return err
		// 		}
		// 		nacNameStr := nacMap["name"].(string)
		// 		createBareMetalServerNetworkAttachmentOptions := &vpcv1.CreateBareMetalServerNetworkAttachmentOptions{
		// 			BareMetalServerID: &id,
		// 		}
		// 		bareMetalServerNetworkAttachmentPrototype := &vpcv1.BareMetalServerNetworkAttachmentPrototype{
		// 			Name:                    &nacNameStr,
		// 			VirtualNetworkInterface: VirtualNetworkInterfaceModel,
		// 		}
		// 		createBareMetalServerNetworkAttachmentOptions.BareMetalServerNetworkAttachmentPrototype = bareMetalServerNetworkAttachmentPrototype
		// 		_, res, err := sess.CreateBareMetalServerNetworkAttachment(createBareMetalServerNetworkAttachmentOptions)
		// 		if err != nil {
		// 			return fmt.Errorf("[ERROR] Error while creating network attachment(%s) of BareMetalServer(%s) \n%s: %q", nacNameStr, d.Id(), err, res)
		// 		}
		// 	} else {
		// 		log.Printf("[DEBUG] nacId is not empty")
		// 		nacName := fmt.Sprintf("network_attachments.%d.name", i)
		// 		nacVniName := fmt.Sprintf("network_attachments.%d.virtual_network_interface", i)
		// 		primaryipName := fmt.Sprintf("%s.%s", nacVniName, "0.primary_ip")
		// 		sgName := fmt.Sprintf("%s.%s", nacVniName, "0.security_groups")
		// 		if d.HasChange(nacName) {
		// 			networkName := d.Get(nacName).(string)
		// 			updateBareMetalServerNetworkAttachmentOptions := &vpcv1.UpdateBareMetalServerNetworkAttachmentOptions{
		// 				BareMetalServerID: &id,
		// 				ID:                &nacId,
		// 			}
		// 			bareMetalServerNetworkAttachmentPatch := &vpcv1.InstanceNetworkAttachmentPatch{
		// 				Name: &networkName,
		// 			}
		// 			bareMetalServerNetworkAttachmentPatchAsPatch, err := bareMetalServerNetworkAttachmentPatch.AsPatch()
		// 			if err != nil {
		// 				return (fmt.Errorf("[ERROR] Error encountered while apply as patch for BareMetalServerNetworkAttachmentPatchAsPatch of network attachment(%s) of instance(%s) %s", nacId, id, err))
		// 			}
		// 			updateBareMetalServerNetworkAttachmentOptions.BareMetalServerNetworkAttachmentPatch = bareMetalServerNetworkAttachmentPatchAsPatch
		// 			_, res, err := sess.UpdateBareMetalServerNetworkAttachment(updateBareMetalServerNetworkAttachmentOptions)
		// 			if err != nil {
		// 				return (fmt.Errorf("[ERROR] Error encountered while updating network attachment(%s) name of BareMetalServer(%s) %s/n%s", nacId, id, err, res))
		// 			}
		// 			// output, err := json.MarshalIndent(updateInstanceNetworkAttachmentOptions, "", "    ")
		// 			// if err == nil {
		// 			// 	log.Printf("%+v\n", string(output))
		// 			// } else {
		// 			// 	log.Printf("Error : %#v", updateInstanceNetworkAttachmentOptions)
		// 			// }
		// 		}
		// 		if d.HasChange(nacVniName) {
		// 			vniId := d.Get(fmt.Sprintf("%s.%s", nacVniName, "0.id")).(string)
		// 			updateVirtualNetworkInterfaceOptions := &vpcv1.UpdateVirtualNetworkInterfaceOptions{
		// 				ID: &vniId,
		// 			}
		// 			virtualNetworkInterfacePatch := &vpcv1.VirtualNetworkInterfacePatch{}
		// 			autoDeleteName := fmt.Sprintf("%s.%s", nacVniName, "0.auto_delete")
		// 			nameName := fmt.Sprintf("%s.%s", nacVniName, "0.name")
		// 			ipsName := fmt.Sprintf("%s.%s", nacVniName, "0.ips")
		// 			enableNatName := fmt.Sprintf("%s.%s", nacVniName, "0.enable_infrastructure_nat")
		// 			allowIpSpoofingName := fmt.Sprintf("%s.%s", nacVniName, "0.allow_ip_spoofing")
		// 			if d.HasChange(autoDeleteName) {
		// 				autodelete := d.Get(autoDeleteName).(bool)
		// 				virtualNetworkInterfacePatch.AutoDelete = &autodelete
		// 			}
		// 			if d.HasChange(nameName) {
		// 				name := d.Get(nameName).(string)
		// 				virtualNetworkInterfacePatch.Name = &name
		// 			}
		// 			if d.HasChange(enableNatName) {
		// 				enableNat := d.Get(enableNatName).(bool)
		// 				virtualNetworkInterfacePatch.EnableInfrastructureNat = &enableNat
		// 			}
		// 			if d.HasChange(allowIpSpoofingName) {
		// 				allIpSpoofing := d.Get(allowIpSpoofingName).(bool)
		// 				virtualNetworkInterfacePatch.AllowIPSpoofing = &allIpSpoofing
		// 			}
		// 			virtualNetworkInterfacePatchAsPatch, err := virtualNetworkInterfacePatch.AsPatch()
		// 			if err != nil {
		// 				return fmt.Errorf("[ERROR] Error encountered while apply as patch for virtualNetworkInterfacePatch of instance(%s) vni (%s) %s", d.Id(), vniId, err)
		// 			}
		// 			updateVirtualNetworkInterfaceOptions.VirtualNetworkInterfacePatch = virtualNetworkInterfacePatchAsPatch
		// 			_, response, err := sess.UpdateVirtualNetworkInterface(updateVirtualNetworkInterfaceOptions)
		// 			if err != nil {
		// 				log.Printf("[DEBUG] UpdateVirtualNetworkInterfaceWithContext failed %s\n%s", err, response)
		// 				return fmt.Errorf("UpdateVirtualNetworkInterfaceWithContext failed during instance(%s) network attachment patch %s\n%s", d.Id(), err, response)
		// 			}

		// 			if d.HasChange(ipsName) {
		// 				oldips, newips := d.GetChange(ipsName)
		// 				os := oldips.(*schema.Set)
		// 				ns := newips.(*schema.Set)
		// 				var oldset, newset *schema.Set

		// 				var out = make([]interface{}, ns.Len(), ns.Len())
		// 				for i, nA := range ns.List() {
		// 					newPack := nA.(map[string]interface{})
		// 					out[i] = newPack["reserved_ip"].(string)
		// 				}
		// 				newset = schema.NewSet(schema.HashString, out)

		// 				out = make([]interface{}, os.Len(), os.Len())
		// 				for i, oA := range os.List() {
		// 					oldPack := oA.(map[string]interface{})
		// 					out[i] = oldPack["reserved_ip"].(string)
		// 				}
		// 				oldset = schema.NewSet(schema.HashString, out)

		// 				remove := flex.ExpandStringList(oldset.Difference(newset).List())
		// 				add := flex.ExpandStringList(newset.Difference(oldset).List())

		// 				if add != nil && len(add) > 0 {
		// 					for _, ipItem := range add {
		// 						if ipItem != "" {
		// 							addVirtualNetworkInterfaceIPOptions := &vpcv1.AddVirtualNetworkInterfaceIPOptions{}
		// 							addVirtualNetworkInterfaceIPOptions.SetVirtualNetworkInterfaceID(vniId)
		// 							addVirtualNetworkInterfaceIPOptions.SetID(ipItem)
		// 							_, response, err := sess.AddVirtualNetworkInterfaceIP(addVirtualNetworkInterfaceIPOptions)
		// 							if err != nil {
		// 								log.Printf("[DEBUG] AddVirtualNetworkInterfaceIPWithContext failed in VirtualNetworkInterface patch during instance nac patch %s\n%s", err, response)
		// 								return fmt.Errorf("AddVirtualNetworkInterfaceIPWithContext failed in VirtualNetworkInterface patch during instance nac patch %s\n%s", err, response)
		// 							}
		// 						}
		// 					}
		// 				}
		// 				if remove != nil && len(remove) > 0 {
		// 					for _, ipItem := range remove {
		// 						if ipItem != "" {
		// 							removeVirtualNetworkInterfaceIPOptions := &vpcv1.RemoveVirtualNetworkInterfaceIPOptions{}
		// 							removeVirtualNetworkInterfaceIPOptions.SetVirtualNetworkInterfaceID(vniId)
		// 							removeVirtualNetworkInterfaceIPOptions.SetID(ipItem)
		// 							response, err := sess.RemoveVirtualNetworkInterfaceIP(removeVirtualNetworkInterfaceIPOptions)
		// 							if err != nil {
		// 								log.Printf("[DEBUG] RemoveVirtualNetworkInterfaceIPWithContext failed in VirtualNetworkInterface patch during instance nac patch %s\n%s", err, response)
		// 								return fmt.Errorf("RemoveVirtualNetworkInterfaceIPWithContext failed in VirtualNetworkInterface patch during instance nac patch %s\n%s", err, response)
		// 							}
		// 						}
		// 					}
		// 				}
		// 			}

		// 			if d.HasChange(primaryipName) {
		// 				subnetIdName := fmt.Sprintf("%s.%s", nacVniName, "0.subnet")
		// 				ripIdName := fmt.Sprintf("%s.%s", primaryipName, "0.reserved_ip")
		// 				subnetId := d.Get(subnetIdName).(string)
		// 				primaryipNameName := fmt.Sprintf("%s.%s", primaryipName, "0.name")
		// 				primaryipAutoDeleteName := fmt.Sprintf("%s.%s", primaryipName, "0.name")
		// 				ripId := d.Get(ripIdName).(string)
		// 				updateripoptions := &vpcv1.UpdateSubnetReservedIPOptions{
		// 					SubnetID: &subnetId,
		// 					ID:       &ripId,
		// 				}
		// 				reservedIpPath := &vpcv1.ReservedIPPatch{}
		// 				if d.HasChange(primaryipNameName) {
		// 					name := d.Get(primaryipNameName).(string)
		// 					reservedIpPath.Name = &name
		// 				}
		// 				if d.HasChange(primaryipAutoDeleteName) {
		// 					auto := d.Get(primaryipAutoDeleteName).(bool)
		// 					reservedIpPath.AutoDelete = &auto
		// 				}
		// 				reservedIpPathAsPatch, err := reservedIpPath.AsPatch()
		// 				if err != nil {
		// 					return fmt.Errorf("[ERROR] Error calling reserved ip as patch on vni patch \n%s", err)
		// 				}
		// 				updateripoptions.ReservedIPPatch = reservedIpPathAsPatch
		// 				_, response, err := sess.UpdateSubnetReservedIP(updateripoptions)
		// 				if err != nil {
		// 					return fmt.Errorf("[ERROR] Error updating vni reserved ip(%s): %s\n%s", ripId, err, response)
		// 				}
		// 			}
		// 			if d.HasChange(sgName) {
		// 				ovs, nvs := d.GetChange(sgName)
		// 				ov := ovs.(*schema.Set)
		// 				nv := nvs.(*schema.Set)
		// 				remove := flex.ExpandStringList(ov.Difference(nv).List())
		// 				add := flex.ExpandStringList(nv.Difference(ov).List())
		// 				if len(add) > 0 {
		// 					for i := range add {
		// 						createsgnicoptions := &vpcv1.CreateSecurityGroupTargetBindingOptions{
		// 							SecurityGroupID: &add[i],
		// 							ID:              &vniId,
		// 						}
		// 						_, response, err := sess.CreateSecurityGroupTargetBinding(createsgnicoptions)
		// 						if err != nil {
		// 							return fmt.Errorf("[ERROR] Error while creating security group %q for virtual network interface %s\n%s: %q", add[i], vniId, err, response)
		// 						}
		// 						_, err = isWaitForVirtualNetworkInterfaceAvailable(sess, vniId, d.Timeout(schema.TimeoutUpdate))
		// 						if err != nil {
		// 							return err
		// 						}
		// 					}

		// 				}
		// 				if len(remove) > 0 {
		// 					for i := range remove {
		// 						deletesgnicoptions := &vpcv1.DeleteSecurityGroupTargetBindingOptions{
		// 							SecurityGroupID: &remove[i],
		// 							ID:              &vniId,
		// 						}
		// 						response, err := sess.DeleteSecurityGroupTargetBinding(deletesgnicoptions)
		// 						if err != nil {
		// 							return fmt.Errorf("[ERROR] Error while removing security group %q for virtual network interface %s\n%s: %q", remove[i], d.Id(), err, response)
		// 						}
		// 						_, err = isWaitForVirtualNetworkInterfaceAvailable(sess, vniId, d.Timeout(schema.TimeoutUpdate))
		// 						if err != nil {
		// 							return err
		// 						}
		// 					}
		// 				}
		// 			}

		// 		}

		// 	}
		// 	// }
		// }

	}

	if d.HasChange(isBareMetalServerTags) || d.HasChange(isBareMetalServerAccessTags) {
		bmscrn := d.Get(isBareMetalServerCRN).(string)
		if bmscrn == "" {
			options := &vpcv1.GetBareMetalServerOptions{
				ID: &id,
			}
			bms, response, err := sess.GetBareMetalServerWithContext(context, options)
			if err != nil {
				if response != nil && response.StatusCode == 404 {
					d.SetId("")
					return nil
				}
				return fmt.Errorf("[ERROR] Error getting Bare Metal Server (%s): %s\n%s", id, err, response)
			}
			bmscrn = *bms.CRN
		}
		if d.HasChange(isBareMetalServerTags) {
			oldList, newList := d.GetChange(isBareMetalServerTags)
			err = flex.UpdateTagsUsingCRN(oldList, newList, meta, bmscrn)
			if err != nil {
				log.Printf(
					"[ERROR] Error on update of vpc Bare metal server (%s) tags: %s", id, err)
			}
		}
		if d.HasChange(isBareMetalServerAccessTags) {
			oldList, newList := d.GetChange(isBareMetalServerAccessTags)
			err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, bmscrn, "", isBareMetalServerAccessTagType)
			if err != nil {
				log.Printf(
					"[ERROR] Error on update of resource vpc Bare metal server (%s) access tags: %s", id, err)
			}
		}
	}

	if d.HasChange(isBareMetalServerNetworkInterfaces) {
		oldList, newList := d.GetChange(isBareMetalServerNetworkInterfaces)
		if oldList == nil {
			oldList = new(schema.Set)
		}
		if newList == nil {
			newList = new(schema.Set)
		}
		os := oldList.(*schema.Set)
		ns := newList.(*schema.Set)
		for _, nA := range ns.List() {
			newPack := nA.(map[string]interface{})
			for _, oA := range os.List() {
				oldPack := oA.(map[string]interface{})
				if strings.Compare(newPack["name"].(string), oldPack["name"].(string)) == 0 {
					networkId := newPack["id"].(string)
					newAllowedVlans := newPack[isBareMetalServerNicAllowedVlans].(*schema.Set)
					newNicName := newPack[isBareMetalServerNicName].(string)
					newIpSpoofing := newPack[isBareMetalServerNicAllowIPSpoofing].(bool)
					newInfraNat := newPack[isBareMetalServerNicEnableInfraNAT].(bool)

					oldAllowedVlans := oldPack[isBareMetalServerNicAllowedVlans].(*schema.Set)
					oldNicName := oldPack[isBareMetalServerNicName].(string)
					oldIpSpoofing := oldPack[isBareMetalServerNicAllowIPSpoofing].(bool)
					oldInfraNat := oldPack[isBareMetalServerNicEnableInfraNAT].(bool)

					if oldAllowedVlans.Difference(newAllowedVlans).Len() > 0 || newAllowedVlans.Difference(oldAllowedVlans).Len() > 0 || newInfraNat != oldInfraNat || newIpSpoofing != oldIpSpoofing {

						updatepnicfoptions := &vpcv1.UpdateBareMetalServerNetworkInterfaceOptions{
							BareMetalServerID: &id,
							ID:                &networkId,
						}

						bmsPatchModel := &vpcv1.BareMetalServerNetworkInterfacePatch{}
						if strings.Compare(newNicName, oldNicName) != 0 {
							bmsPatchModel.Name = &newNicName
						}

						if oldAllowedVlans.Difference(newAllowedVlans).Len() > 0 || newAllowedVlans.Difference(oldAllowedVlans).Len() > 0 {
							allowedVlansList := newPack[isBareMetalServerNicAllowedVlans].(*schema.Set).List()
							allowedVlans := make([]int64, 0, len(allowedVlansList))
							for _, k := range allowedVlansList {
								allowedVlans = append(allowedVlans, int64(k.(int)))
							}
							bmsPatchModel.AllowedVlans = allowedVlans
						}

						if newIpSpoofing != oldIpSpoofing {
							bmsPatchModel.AllowIPSpoofing = &newIpSpoofing
						}
						if newInfraNat != oldInfraNat {
							bmsPatchModel.EnableInfrastructureNat = &newInfraNat
						}
						networkInterfacePatch, err := bmsPatchModel.AsPatch()
						if err != nil {
							return fmt.Errorf("[ERROR] Error calling asPatch for BareMetalServerNetworkInterfacePatch: %s", err)
						}
						updatepnicfoptions.BareMetalServerNetworkInterfacePatch = networkInterfacePatch
						_, response, err := sess.UpdateBareMetalServerNetworkInterface(updatepnicfoptions)
						if err != nil {
							return fmt.Errorf("[ERROR] Error while updating network interface(%s) of bar emetal server(%s) \n%s: %q", networkId, d.Id(), err, response)
						}
						ns.Remove(nA)
						os.Remove(oA)
					}
				}
			}
		}
		remove := os.Difference(ns).List()

		if len(remove) > 0 {
			// check if any removing nic is of pci type
			flag := false
			for _, rem := range remove {
				oldNic := rem.(map[string]interface{})
				interfaceType := oldNic["interface_type"].(string)
				if interfaceType == "pci" {
					flag = true
				}
			}
			if flag {
				isServerStopped, err = resourceStopServerIfRunning(id, "hard", d, context, sess, isServerStopped)
				if err != nil {
					return err
				}
			}
			for _, rem := range remove {
				oldNic := rem.(map[string]interface{})
				networkId := oldNic["id"].(string)
				removeBMSNic := &vpcv1.DeleteBareMetalServerNetworkInterfaceOptions{
					BareMetalServerID: &id,
					ID:                &networkId,
				}
				_, err = sess.DeleteBareMetalServerNetworkInterface(removeBMSNic)
				if err != nil {
					return err
				}
			}

		}
		add := ns.Difference(os).List()
		if len(add) > 0 {
			// check if any adding nic is of pci type
			flag := false
			for _, a := range add {
				oldNic := a.(map[string]interface{})
				allowedVlansOk, ok := oldNic[isBareMetalServerNicAllowedVlans]
				if ok && len(allowedVlansOk.(*schema.Set).List()) > 0 {
					flag = true
				}
			}
			if flag {
				isServerStopped, err = resourceStopServerIfRunning(id, "hard", d, context, sess, isServerStopped)
				if err != nil {
					return err
				}
			}
			for _, a := range add {
				nic := a.(map[string]interface{})
				addNicOptions := &vpcv1.CreateBareMetalServerNetworkInterfaceOptions{
					BareMetalServerID: &id,
				}
				interfaceType := ""
				if interfaceTypeOk, ok := nic[isBareMetalServerNicInterfaceType]; ok && interfaceTypeOk.(string) == "hipersocket" {
					interfaceType = "hipersocket"
					var nicobj = &vpcv1.BareMetalServerNetworkInterfacePrototypeBareMetalServerNetworkInterfaceByHiperSocketPrototype{}
					nicobj.InterfaceType = &interfaceType

					subnetintf, _ := nic[isBareMetalServerNicSubnet]
					subnetintfstr := subnetintf.(string)
					nicobj.Subnet = &vpcv1.SubnetIdentity{
						ID: &subnetintfstr,
					}
					name, _ := nic[isBareMetalServerNicName]
					namestr := name.(string)
					if namestr != "" {
						nicobj.Name = &namestr
					}

					enableInfraNAT, ok := nic[isBareMetalServerNicEnableInfraNAT]
					enableInfraNATbool := enableInfraNAT.(bool)
					if ok {
						nicobj.EnableInfrastructureNat = &enableInfraNATbool
					}

					if primaryIpIntf, ok := nic[isBareMetalServerNicPrimaryIP]; ok && len(primaryIpIntf.([]interface{})) > 0 {
						primaryIp := primaryIpIntf.([]interface{})[0].(map[string]interface{})
						reservedIpIdOk, ok := primaryIp[isBareMetalServerNicIpID]
						if ok && reservedIpIdOk.(string) != "" {
							ipid := reservedIpIdOk.(string)
							nicobj.PrimaryIP = &vpcv1.NetworkInterfaceIPPrototypeReservedIPIdentity{
								ID: &ipid,
							}
						} else {
							primaryip := &vpcv1.NetworkInterfaceIPPrototypeReservedIPPrototypeNetworkInterfaceContext{}

							reservedIpAddressOk, okAdd := primaryIp[isBareMetalServerNicIpAddress]
							if okAdd && reservedIpAddressOk.(string) != "" {
								reservedIpAddress := reservedIpAddressOk.(string)
								primaryip.Address = &reservedIpAddress
							}

							reservedIpNameOk, okName := primaryIp[isBareMetalServerNicIpName]
							if okName && reservedIpNameOk.(string) != "" {
								reservedIpName := reservedIpNameOk.(string)
								primaryip.Name = &reservedIpName
							}

							reservedIpAutoOk, okAuto := primaryIp[isBareMetalServerNicIpAutoDelete]
							if okAuto {
								reservedIpAuto := reservedIpAutoOk.(bool)
								primaryip.AutoDelete = &reservedIpAuto
							}
							if okAdd || okName || okAuto {
								nicobj.PrimaryIP = primaryip
							}
						}
					}

					allowIPSpoofing, ok := nic[isBareMetalServerNicAllowIPSpoofing]
					allowIPSpoofingbool := allowIPSpoofing.(bool)
					if ok && allowIPSpoofingbool {
						nicobj.AllowIPSpoofing = &allowIPSpoofingbool
					}
					secgrpintf, ok := nic[isBareMetalServerNicSecurityGroups]
					if ok {
						secgrpSet := secgrpintf.(*schema.Set)
						if secgrpSet.Len() != 0 {
							var secgrpobjs = make([]vpcv1.SecurityGroupIdentityIntf, secgrpSet.Len())
							for i, secgrpIntf := range secgrpSet.List() {
								secgrpIntfstr := secgrpIntf.(string)
								secgrpobjs[i] = &vpcv1.SecurityGroupIdentity{
									ID: &secgrpIntfstr,
								}
							}
							nicobj.SecurityGroups = secgrpobjs
						}
					}
					addNicOptions.BareMetalServerNetworkInterfacePrototype = nicobj
				} else if allowedVlansOk, ok := nic[isBareMetalServerNicAllowedVlans]; ok {
					interfaceType = "pci"
					var nicobj = &vpcv1.BareMetalServerNetworkInterfacePrototypeBareMetalServerNetworkInterfaceByPciPrototype{}
					nicobj.InterfaceType = &interfaceType

					allowedVlansList := allowedVlansOk.(*schema.Set).List()

					if len(allowedVlansList) > 0 {
						allowedVlans := make([]int64, 0, len(allowedVlansList))
						for _, k := range allowedVlansList {
							allowedVlans = append(allowedVlans, int64(k.(int)))
						}
						nicobj.AllowedVlans = allowedVlans

						subnetintf, _ := nic[isBareMetalServerNicSubnet]
						subnetintfstr := subnetintf.(string)
						nicobj.Subnet = &vpcv1.SubnetIdentity{
							ID: &subnetintfstr,
						}
						name, _ := nic[isBareMetalServerNicName]
						namestr := name.(string)
						if namestr != "" {
							nicobj.Name = &namestr
						}

						enableInfraNAT, ok := nic[isBareMetalServerNicEnableInfraNAT]
						enableInfraNATbool := enableInfraNAT.(bool)
						if ok {
							nicobj.EnableInfrastructureNat = &enableInfraNATbool
						}

						if primaryIpIntf, ok := nic[isBareMetalServerNicPrimaryIP]; ok && len(primaryIpIntf.([]interface{})) > 0 {
							primaryIp := primaryIpIntf.([]interface{})[0].(map[string]interface{})

							reservedIpIdOk, ok := primaryIp[isBareMetalServerNicIpID]
							if ok && reservedIpIdOk.(string) != "" {
								ipid := reservedIpIdOk.(string)
								nicobj.PrimaryIP = &vpcv1.NetworkInterfaceIPPrototypeReservedIPIdentity{
									ID: &ipid,
								}
							} else {
								primaryip := &vpcv1.NetworkInterfaceIPPrototypeReservedIPPrototypeNetworkInterfaceContext{}
								reservedIpAddressOk, okAdd := primaryIp[isBareMetalServerNicIpAddress]
								if okAdd && reservedIpAddressOk.(string) != "" {
									reservedIpAddress := reservedIpAddressOk.(string)
									primaryip.Address = &reservedIpAddress
								}
								reservedIpNameOk, okName := primaryIp[isBareMetalServerNicIpName]
								if okName && reservedIpNameOk.(string) != "" {
									reservedIpName := reservedIpNameOk.(string)
									primaryip.Name = &reservedIpName
								}
								reservedIpAutoOk, okAuto := primaryIp[isBareMetalServerNicIpAutoDelete]
								if okAuto {
									reservedIpAuto := reservedIpAutoOk.(bool)
									primaryip.AutoDelete = &reservedIpAuto
								}
								if okAdd || okName || okAuto {
									nicobj.PrimaryIP = primaryip
								}
							}

						}

						allowIPSpoofing, ok := nic[isBareMetalServerNicAllowIPSpoofing]
						allowIPSpoofingbool := allowIPSpoofing.(bool)
						if ok && allowIPSpoofingbool {
							nicobj.AllowIPSpoofing = &allowIPSpoofingbool
						}
						secgrpintf, ok := nic[isBareMetalServerNicSecurityGroups]
						if ok {
							secgrpSet := secgrpintf.(*schema.Set)
							if secgrpSet.Len() != 0 {
								var secgrpobjs = make([]vpcv1.SecurityGroupIdentityIntf, secgrpSet.Len())
								for i, secgrpIntf := range secgrpSet.List() {
									secgrpIntfstr := secgrpIntf.(string)
									secgrpobjs[i] = &vpcv1.SecurityGroupIdentity{
										ID: &secgrpIntfstr,
									}
								}
								nicobj.SecurityGroups = secgrpobjs
							}
						}
						addNicOptions.BareMetalServerNetworkInterfacePrototype = nicobj
					} else {
						interfaceType = "vlan"
						var nicobj = &vpcv1.BareMetalServerNetworkInterfacePrototypeBareMetalServerNetworkInterfaceByVlanPrototype{}
						nicobj.InterfaceType = &interfaceType

						if aitf, ok := nic[isBareMetalServerNicAllowInterfaceToFloat]; ok {
							allowInterfaceToFloat := aitf.(bool)
							nicobj.AllowInterfaceToFloat = &allowInterfaceToFloat
						}
						if vlan, ok := nic[isBareMetalServerNicVlan]; ok {
							vlanInt := int64(vlan.(int))
							nicobj.Vlan = &vlanInt
						}

						subnetintf, _ := nic[isBareMetalServerNicSubnet]
						subnetintfstr := subnetintf.(string)
						nicobj.Subnet = &vpcv1.SubnetIdentity{
							ID: &subnetintfstr,
						}
						name, _ := nic[isBareMetalServerNicName]
						namestr := name.(string)
						if namestr != "" {
							nicobj.Name = &namestr
						}

						enableInfraNAT, ok := nic[isBareMetalServerNicEnableInfraNAT]
						enableInfraNATbool := enableInfraNAT.(bool)
						if ok {
							nicobj.EnableInfrastructureNat = &enableInfraNATbool
						}

						if primaryIpIntf, ok := nic[isBareMetalServerNicPrimaryIP]; ok && len(primaryIpIntf.([]interface{})) > 0 {
							primaryIp := primaryIpIntf.([]interface{})[0].(map[string]interface{})
							reservedIpIdOk, ok := primaryIp[isBareMetalServerNicIpID]
							if ok && reservedIpIdOk.(string) != "" {
								ipid := reservedIpIdOk.(string)
								nicobj.PrimaryIP = &vpcv1.NetworkInterfaceIPPrototypeReservedIPIdentity{
									ID: &ipid,
								}
							} else {
								primaryip := &vpcv1.NetworkInterfaceIPPrototypeReservedIPPrototypeNetworkInterfaceContext{}

								reservedIpAddressOk, okAdd := primaryIp[isBareMetalServerNicIpAddress]
								if okAdd && reservedIpAddressOk.(string) != "" {
									reservedIpAddress := reservedIpAddressOk.(string)
									primaryip.Address = &reservedIpAddress
								}

								reservedIpNameOk, okName := primaryIp[isBareMetalServerNicIpName]
								if okName && reservedIpNameOk.(string) != "" {
									reservedIpName := reservedIpNameOk.(string)
									primaryip.Name = &reservedIpName
								}

								reservedIpAutoOk, okAuto := primaryIp[isBareMetalServerNicIpAutoDelete]
								if okAuto {
									reservedIpAuto := reservedIpAutoOk.(bool)
									primaryip.AutoDelete = &reservedIpAuto
								}
								if okAdd || okName || okAuto {
									nicobj.PrimaryIP = primaryip
								}
							}
						}

						allowIPSpoofing, ok := nic[isBareMetalServerNicAllowIPSpoofing]
						allowIPSpoofingbool := allowIPSpoofing.(bool)
						if ok && allowIPSpoofingbool {
							nicobj.AllowIPSpoofing = &allowIPSpoofingbool
						}
						secgrpintf, ok := nic[isBareMetalServerNicSecurityGroups]
						if ok {
							secgrpSet := secgrpintf.(*schema.Set)
							if secgrpSet.Len() != 0 {
								var secgrpobjs = make([]vpcv1.SecurityGroupIdentityIntf, secgrpSet.Len())
								for i, secgrpIntf := range secgrpSet.List() {
									secgrpIntfstr := secgrpIntf.(string)
									secgrpobjs[i] = &vpcv1.SecurityGroupIdentity{
										ID: &secgrpIntfstr,
									}
								}
								nicobj.SecurityGroups = secgrpobjs
							}
						}
						addNicOptions.BareMetalServerNetworkInterfacePrototype = nicobj
					}
				} else {
					interfaceType = "vlan"
					var nicobj = &vpcv1.BareMetalServerNetworkInterfacePrototypeBareMetalServerNetworkInterfaceByVlanPrototype{}
					nicobj.InterfaceType = &interfaceType

					if aitf, ok := nic[isBareMetalServerNicAllowInterfaceToFloat]; ok {
						allowInterfaceToFloat := aitf.(bool)
						nicobj.AllowInterfaceToFloat = &allowInterfaceToFloat
					}
					if vlan, ok := nic[isBareMetalServerNicVlan]; ok {
						vlanInt := int64(vlan.(int))
						nicobj.Vlan = &vlanInt
					}

					subnetintf, _ := nic[isBareMetalServerNicSubnet]
					subnetintfstr := subnetintf.(string)
					nicobj.Subnet = &vpcv1.SubnetIdentity{
						ID: &subnetintfstr,
					}
					name, _ := nic[isBareMetalServerNicName]
					namestr := name.(string)
					if namestr != "" {
						nicobj.Name = &namestr
					}

					enableInfraNAT, ok := nic[isBareMetalServerNicEnableInfraNAT]
					enableInfraNATbool := enableInfraNAT.(bool)
					if ok {
						nicobj.EnableInfrastructureNat = &enableInfraNATbool
					}

					if primaryIpIntf, ok := nic[isBareMetalServerNicPrimaryIP]; ok && len(primaryIpIntf.([]interface{})) > 0 {
						primaryIp := primaryIpIntf.([]interface{})[0].(map[string]interface{})
						reservedIpIdOk, ok := primaryIp[isBareMetalServerNicIpID]
						if ok && reservedIpIdOk.(string) != "" {
							ipid := reservedIpIdOk.(string)
							nicobj.PrimaryIP = &vpcv1.NetworkInterfaceIPPrototypeReservedIPIdentity{
								ID: &ipid,
							}
						} else {
							primaryip := &vpcv1.NetworkInterfaceIPPrototypeReservedIPPrototypeNetworkInterfaceContext{}

							reservedIpAddressOk, okAdd := primaryIp[isBareMetalServerNicIpAddress]
							if okAdd && reservedIpAddressOk.(string) != "" {
								reservedIpAddress := reservedIpAddressOk.(string)
								primaryip.Address = &reservedIpAddress
							}

							reservedIpNameOk, okName := primaryIp[isBareMetalServerNicIpName]
							if okName && reservedIpNameOk.(string) != "" {
								reservedIpName := reservedIpNameOk.(string)
								primaryip.Name = &reservedIpName
							}

							reservedIpAutoOk, okAuto := primaryIp[isBareMetalServerNicIpAutoDelete]
							if okAuto {
								reservedIpAuto := reservedIpAutoOk.(bool)
								primaryip.AutoDelete = &reservedIpAuto
							}
							if okAdd || okName || okAuto {
								nicobj.PrimaryIP = primaryip
							}
						}
					}

					allowIPSpoofing, ok := nic[isBareMetalServerNicAllowIPSpoofing]
					allowIPSpoofingbool := allowIPSpoofing.(bool)
					if ok && allowIPSpoofingbool {
						nicobj.AllowIPSpoofing = &allowIPSpoofingbool
					}
					secgrpintf, ok := nic[isBareMetalServerNicSecurityGroups]
					if ok {
						secgrpSet := secgrpintf.(*schema.Set)
						if secgrpSet.Len() != 0 {
							var secgrpobjs = make([]vpcv1.SecurityGroupIdentityIntf, secgrpSet.Len())
							for i, secgrpIntf := range secgrpSet.List() {
								secgrpIntfstr := secgrpIntf.(string)
								secgrpobjs[i] = &vpcv1.SecurityGroupIdentity{
									ID: &secgrpIntfstr,
								}
							}
							nicobj.SecurityGroups = secgrpobjs
						}
					}
					addNicOptions.BareMetalServerNetworkInterfacePrototype = nicobj
				}
				_, _, err := sess.CreateBareMetalServerNetworkInterface(addNicOptions)
				if err != nil {
					return err
				}
			}

		}

		// if len(add) > 0 {
		// }
		// if len(remove) > 0 {
		// }
		// nics := d.Get(isBareMetalServerNetworkInterfaces).(*schema.Set).List()
		// for i := range nics {
		// 	securitygrpKey := fmt.Sprintf("network_interfaces.%d.security_groups", i)
		// 	networkNameKey := fmt.Sprintf("network_interfaces.%d.name", i)
		// 	ipSpoofingKey := fmt.Sprintf("network_interfaces.%d.allow_ip_spoofing", i)
		// 	infraNatKey := fmt.Sprintf("network_interfaces.%d.enable_infrastructure_nat", i)
		// 	allowedVlans := fmt.Sprintf("network_interfaces.%d.allowed_vlans", i)
		// 	primaryipname := fmt.Sprintf("network_interfaces.%d.primary_ip.0.name", i)
		// 	primaryipauto := fmt.Sprintf("network_interfaces.%d.primary_ip.0.auto_delete", i)
		// 	primaryiprip := fmt.Sprintf("network_interfaces.%d.primary_ip.0.reserved_ip", i)
		// 	if d.HasChange(primaryipname) || d.HasChange(primaryipauto) {
		// 		subnetId := d.Get(isBareMetalServerNicSubnet).(string)
		// 		ripId := d.Get(primaryiprip).(string)
		// 		updateripoptions := &vpcv1.UpdateSubnetReservedIPOptions{
		// 			SubnetID: &subnetId,
		// 			ID:       &ripId,
		// 		}
		// 		reservedIpPath := &vpcv1.ReservedIPPatch{}
		// 		if d.HasChange(primaryipname) {
		// 			name := d.Get(primaryipname).(string)
		// 			reservedIpPath.Name = &name
		// 		}
		// 		if d.HasChange(primaryipauto) {
		// 			auto := d.Get(primaryipauto).(bool)
		// 			reservedIpPath.AutoDelete = &auto
		// 		}
		// 		reservedIpPathAsPatch, err := reservedIpPath.AsPatch()
		// 		if err != nil {
		// 			return fmt.Errorf("[ERROR] Error calling reserved ip as patch \n%s", err)
		// 		}
		// 		updateripoptions.ReservedIPPatch = reservedIpPathAsPatch
		// 		_, response, err := sess.UpdateSubnetReservedIP(updateripoptions)
		// 		if err != nil {
		// 			return fmt.Errorf("[ERROR] Error updating bare metal server network interface reserved ip(%s): %s\n%s", ripId, err, response)
		// 		}
		// 	}

		// 	if d.HasChange(securitygrpKey) {
		// 		ovs, nvs := d.GetChange(securitygrpKey)
		// 		ov := ovs.(*schema.Set)
		// 		nv := nvs.(*schema.Set)
		// 		remove := flex.ExpandStringList(ov.Difference(nv).List())
		// 		add := flex.ExpandStringList(nv.Difference(ov).List())
		// 		if len(add) > 0 {
		// 			networkIDKey := fmt.Sprintf("network_interfaces.%d.id", i)
		// 			networkID := d.Get(networkIDKey).(string)
		// 			for i := range add {
		// 				createsgnicoptions := &vpcv1.CreateSecurityGroupTargetBindingOptions{
		// 					SecurityGroupID: &add[i],
		// 					ID:              &networkID,
		// 				}
		// 				_, response, err := sess.CreateSecurityGroupTargetBinding(createsgnicoptions)
		// 				if err != nil {
		// 					return fmt.Errorf("[ERROR] Error while creating security group %q for network interface of bare metal server %s\n%s: %q", add[i], d.Id(), err, response)
		// 				}
		// 			}

		// 		}
		// 		if len(remove) > 0 {
		// 			networkIDKey := fmt.Sprintf("network_interfaces.%d.id", i)
		// 			networkID := d.Get(networkIDKey).(string)
		// 			for i := range remove {
		// 				deletesgnicoptions := &vpcv1.DeleteSecurityGroupTargetBindingOptions{
		// 					SecurityGroupID: &remove[i],
		// 					ID:              &networkID,
		// 				}
		// 				response, err := sess.DeleteSecurityGroupTargetBinding(deletesgnicoptions)
		// 				if err != nil {
		// 					return fmt.Errorf("[ERROR] Error while removing security group %q for network interface of instance %s\n%s: %q", remove[i], d.Id(), err, response)
		// 				}
		// 			}
		// 		}

		// 	}

	}
	options := &vpcv1.UpdateBareMetalServerOptions{
		ID: &id,
	}
	bmsPatchModel := &vpcv1.BareMetalServerPatch{}
	flag := false

	if d.HasChange(isBareMetalServerBandwidth) && !d.IsNewResource() {
		bandwidth := int64(d.Get(isBareMetalServerBandwidth).(int))
		bmsPatchModel.Bandwidth = &bandwidth
		flag = true
	}

	if d.HasChange(isBareMetalServerEnableSecureBoot) {
		newEnableSecureBoot := d.Get(isBareMetalServerEnableSecureBoot).(bool)
		bmsPatchModel.EnableSecureBoot = &newEnableSecureBoot
		flag = true
		isServerStopped, err = resourceStopServerIfRunning(id, "hard", d, context, sess, isServerStopped)
		if err != nil {
			return err
		}
	}

	// tpm
	if d.HasChange("trusted_platform_module") && d.HasChange("trusted_platform_module.0.mode") {
		bareMetalServerTrustedPlatformModulePatch := &vpcv1.BareMetalServerTrustedPlatformModulePatch{}
		newModeTPM := d.Get("trusted_platform_module.0.mode").(string)
		bareMetalServerTrustedPlatformModulePatch.Mode = &newModeTPM
		bmsPatchModel.TrustedPlatformModule = bareMetalServerTrustedPlatformModulePatch
		flag = true
		isServerStopped, err = resourceStopServerIfRunning(id, "hard", d, context, sess, isServerStopped)
		if err != nil {
			return err
		}
	}

	if d.HasChange(isBareMetalServerPrimaryNetworkInterface) {
		nicId := d.Get("primary_network_interface.0.id").(string)
		nicflag := false
		if d.HasChange("primary_network_interface.0.primary_ip.0.name") || d.HasChange("primary_network_interface.0.primary_ip.0.auto_delete") {
			subnetId := d.Get("primary_network_interface.0.subnet").(string)
			ripId := d.Get("primary_network_interface.0.primary_ip.0.reserved_ip").(string)
			updateripoptions := &vpcv1.UpdateSubnetReservedIPOptions{
				SubnetID: &subnetId,
				ID:       &ripId,
			}
			reservedIpPath := &vpcv1.ReservedIPPatch{}
			if d.HasChange("primary_network_interface.0.primary_ip.0.name") {
				name := d.Get("primary_network_interface.0.primary_ip.0.name").(string)
				reservedIpPath.Name = &name
			}
			if d.HasChange("primary_network_interface.0.primary_ip.0.auto_delete") {
				auto := d.Get("primary_network_interface.0.primary_ip.0.auto_delete").(bool)
				reservedIpPath.AutoDelete = &auto
			}
			reservedIpPathAsPatch, err := reservedIpPath.AsPatch()
			if err != nil {
				return fmt.Errorf("[ERROR] Error calling reserved ip as patch \n%s", err)
			}
			updateripoptions.ReservedIPPatch = reservedIpPathAsPatch
			_, response, err := sess.UpdateSubnetReservedIP(updateripoptions)
			if err != nil {
				return fmt.Errorf("[ERROR] Error updating bare metal server primary network interface reserved ip(%s): %s\n%s", ripId, err, response)
			}
		}
		bmsNicUpdateOptions := &vpcv1.UpdateBareMetalServerNetworkInterfaceOptions{
			BareMetalServerID: &id,
			ID:                &nicId,
		}
		bmsNicPatchModel := &vpcv1.BareMetalServerNetworkInterfacePatch{}
		if d.HasChange("primary_network_interface.0.allowed_vlans") {
			if allowedVlansOk, ok := d.GetOk("primary_network_interface.0.allowed_vlans"); ok {
				allowedVlansList := allowedVlansOk.(*schema.Set).List()
				allowedVlans := make([]int64, 0, len(allowedVlansList))
				for _, k := range allowedVlansList {
					allowedVlans = append(allowedVlans, int64(k.(int)))
				}
				bmsNicPatchModel.AllowedVlans = allowedVlans
			}
			nicflag = true
		}
		if d.HasChange("primary_network_interface.0.allow_ip_spoofing") {

			if allowIpSpoofingOk, ok := d.GetOk("primary_network_interface.0.allow_ip_spoofing"); ok {
				allowIpSpoofing := allowIpSpoofingOk.(bool)
				if allowIpSpoofing {
					bmsNicPatchModel.AllowIPSpoofing = &allowIpSpoofing
				}
				nicflag = true
			}
		}
		if d.HasChange("primary_network_interface.0.enable_infrastructure_nat") {
			if enableNatOk, ok := d.GetOk("primary_network_interface.0.enable_infrastructure_nat"); ok {
				enableNat := enableNatOk.(bool)
				bmsNicPatchModel.EnableInfrastructureNat = &enableNat
				nicflag = true
			}
		}
		if d.HasChange("primary_network_interface.0.security_groups") && !d.IsNewResource() {
			ovs, nvs := d.GetChange("primary_network_interface.0.security_groups")
			ov := ovs.(*schema.Set)
			nv := nvs.(*schema.Set)
			remove := flex.ExpandStringList(ov.Difference(nv).List())
			add := flex.ExpandStringList(nv.Difference(ov).List())
			if len(add) > 0 {
				networkID := d.Get("primary_network_interface.0.id").(string)
				for i := range add {
					createsgnicoptions := &vpcv1.CreateSecurityGroupTargetBindingOptions{
						SecurityGroupID: &add[i],
						ID:              &networkID,
					}
					_, response, err := sess.CreateSecurityGroupTargetBinding(createsgnicoptions)
					if err != nil {
						return fmt.Errorf("[ERROR] Error while creating security group %q for primary network interface of bare metal server %s\n%s: %q", add[i], d.Id(), err, response)
					}
					_, err = isWaitForBareMetalServerAvailable(sess, id, d.Timeout(schema.TimeoutUpdate), d)
					if err != nil {
						return err
					}
				}

			}
			if len(remove) > 0 {
				networkID := d.Get("primary_network_interface.0.id").(string)
				for i := range remove {
					deletesgnicoptions := &vpcv1.DeleteSecurityGroupTargetBindingOptions{
						SecurityGroupID: &remove[i],
						ID:              &networkID,
					}
					response, err := sess.DeleteSecurityGroupTargetBinding(deletesgnicoptions)
					if err != nil {
						return fmt.Errorf("[ERROR] Error while removing security group %q for primary network interface of bare metal server %s\n%s: %q", remove[i], d.Id(), err, response)
					}
					_, err = isWaitForBareMetalServerAvailable(sess, id, d.Timeout(schema.TimeoutUpdate), d)
					if err != nil {
						return err
					}
				}
			}
		}
		if d.HasChange("primary_network_interface.0.name") {
			if nameOk, ok := d.GetOk("primary_network_interface.0.name"); ok {
				name := nameOk.(string)
				bmsNicPatchModel.Name = &name
				nicflag = true
			}
		}
		if nicflag {
			bmsNicPatch, err := bmsNicPatchModel.AsPatch()
			if err != nil {
				return err
			}
			bmsNicUpdateOptions.BareMetalServerNetworkInterfacePatch = bmsNicPatch
			_, _, err = sess.UpdateBareMetalServerNetworkInterfaceWithContext(context, bmsNicUpdateOptions)
			if err != nil {
				return err
			}
			_, err = isWaitForBareMetalServerAvailable(sess, id, d.Timeout(schema.TimeoutUpdate), d)
			if err != nil {
				return err
			}
		}
	}
	if d.HasChange(isBareMetalServerName) {
		flag = true
		nameStr := ""
		if name, ok := d.GetOk(isBareMetalServerName); ok {
			nameStr = name.(string)
		}
		bmsPatchModel.Name = &nameStr
	}
	if flag {
		bmsPatch, err := bmsPatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("[ERROR] Error calling asPatch for BareMetalServerPatch: %s", err)
		}
		options.BareMetalServerPatch = bmsPatch
		_, response, err := sess.UpdateBareMetalServerWithContext(context, options)
		if err != nil {
			return fmt.Errorf("[ERROR] Error updating Bare Metal Server: %s\n%s", err, response)
		}
	}

	if d.HasChange(isBareMetalServerAction) {
		action := ""
		if actionOk, ok := d.GetOk(isBareMetalServerAction); ok {
			action = actionOk.(string)
		}
		if action == "start" {
			isBareMetalServerStart(sess, d.Id(), d, 10)
		} else if action == "stop" {
			isBareMetalServerStop(sess, d.Id(), d, 10)
		} else if action == "restart" {
			isBareMetalServerRestart(sess, d.Id(), d, 10)
		}
	}

	if flag || isServerStopped {
		isServerStopped, err = resourceStartServerIfStopped(id, "hard", d, context, sess, isServerStopped)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceIBMISBareMetalServerDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := d.Id()
	deleteType := "hard"
	if dt, ok := d.GetOk(isBareMetalServerDeleteType); ok {
		deleteType = dt.(string)
	}
	err := bareMetalServerDelete(context, d, meta, id, deleteType)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func bareMetalServerDelete(context context.Context, d *schema.ResourceData, meta interface{}, id, deleteType string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	getBmsOptions := &vpcv1.GetBareMetalServerOptions{
		ID: &id,
	}
	bms, response, err := sess.GetBareMetalServerWithContext(context, getBmsOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil
		}
		return fmt.Errorf("[ERROR] Error Getting Bare Metal Server (%s): %s\n%s", id, err, response)
	}
	if *bms.Status == "running" {

		options := &vpcv1.StopBareMetalServerOptions{
			ID:   bms.ID,
			Type: &deleteType,
		}

		response, err := sess.StopBareMetalServerWithContext(context, options)
		if err != nil && response != nil && response.StatusCode != 204 {
			return fmt.Errorf("[ERROR] Error stopping Bare Metal Server (%s): %s\n%s", id, err, response)
		}
		isWaitForBareMetalServerActionStop(sess, d.Timeout(schema.TimeoutDelete), id, d)

	}
	options := &vpcv1.DeleteBareMetalServerOptions{
		ID: &id,
	}
	response, err = sess.DeleteBareMetalServerWithContext(context, options)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Deleting Bare Metal Server : %s\n%s", err, response)
	}
	_, err = isWaitForBareMetalServerDeleted(sess, id, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func isWaitForBareMetalServerDeleted(bmsC *vpcv1.VpcV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for  (%s) to be deleted.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isBareMetalServerActionDeleting},
		Target:     []string{"done", "", isBareMetalServerActionDeleted, isBareMetalServerStatusFailed},
		Refresh:    isBareMetalServerDeleteRefreshFunc(bmsC, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isBareMetalServerDeleteRefreshFunc(bmsC *vpcv1.VpcV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		bmsgetoptions := &vpcv1.GetBareMetalServerOptions{
			ID: &id,
		}
		bms, response, err := bmsC.GetBareMetalServer(bmsgetoptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return bms, isBareMetalServerActionDeleted, nil
			}
			return bms, "", fmt.Errorf("[ERROR] Error Getting Bare Metal Server: %s\n%s", err, response)
		}
		if *bms.Status == isBareMetalServerStatusFailed {
			return bms, *bms.Status, fmt.Errorf("[ERROR] The Bare Metal Server (%s) failed to delete: %v", *bms.ID, err)
		}
		return bms, isBareMetalServerActionDeleting, err
	}
}

func isWaitForBareMetalServerAvailable(client *vpcv1.VpcV1, id string, timeout time.Duration, d *schema.ResourceData) (interface{}, error) {
	log.Printf("Waiting for Bare Metal Server (%s) to be available.", id)
	communicator := make(chan interface{})
	stateConf := &resource.StateChangeConf{
		Pending:    []string{isBareMetalServerStatusPending, isBareMetalServerActionStatusStarting},
		Target:     []string{isBareMetalServerStatusRunning, isBareMetalServerStatusFailed},
		Refresh:    isBareMetalServerRefreshFunc(client, id, d, communicator),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	return stateConf.WaitForState()
}

func isBareMetalServerRefreshFunc(client *vpcv1.VpcV1, id string, d *schema.ResourceData, communicator chan interface{}) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		bmsgetoptions := &vpcv1.GetBareMetalServerOptions{
			ID: &id,
		}
		bms, response, err := client.GetBareMetalServer(bmsgetoptions)
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] Error getting Bare Metal Server: %s\n%s", err, response)
		}
		d.Set(isBareMetalServerStatus, *bms.Status)

		select {
		case data := <-communicator:
			return nil, "", data.(error)
		default:
			fmt.Println("no message sent")
		}

		if *bms.Status == "running" || *bms.Status == "failed" {
			// let know the isRestartStartAction() to stop
			close(communicator)
			if *bms.Status == "failed" {
				bmsStatusReason := bms.StatusReasons

				//set the status reasons
				if bms.StatusReasons != nil {
					statusReasonsList := make([]map[string]interface{}, 0)
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
					d.Set(isBareMetalServerStatusReasons, statusReasonsList)
				}

				out, err := json.MarshalIndent(bmsStatusReason, "", "    ")
				if err != nil {
					return bms, *bms.Status, fmt.Errorf("[ERROR] The Bare Metal Server (%s) went into failed state during the operation \n [WARNING] Running terraform apply again will remove the tainted bare metal server and attempt to create the bare metal server again replacing the previous configuration", *bms.ID)
				}
				return bms, *bms.Status, fmt.Errorf("[ERROR] Bare Metal Server (%s) went into failed state during the operation \n (%+v) \n [WARNING] Running terraform apply again will remove the tainted Bare Metal Server and attempt to create the Bare Metal Server again replacing the previous configuration", *bms.ID, string(out))
			}
			return bms, *bms.Status, nil

		}
		return bms, isBareMetalServerStatusPending, nil
	}
}
func isWaitForBareMetalServerStoppedOnReload(client *vpcv1.VpcV1, id string, timeout time.Duration, d *schema.ResourceData) (interface{}, error) {
	log.Printf("Waiting for Bare Metal Server (%s) to be stopped for reload success.", id)
	communicator := make(chan interface{})
	stateConf := &resource.StateChangeConf{
		Pending:    []string{isBareMetalServerStatusPending, isBareMetalServerActionStatusStarting, "reinitializing"},
		Target:     []string{isBareMetalServerStatusRunning, isBareMetalServerStatusFailed, "stopped"},
		Refresh:    isBareMetalServerRefreshFuncForReload(client, id, d, communicator),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	return stateConf.WaitForState()
}

func isBareMetalServerRefreshFuncForReload(client *vpcv1.VpcV1, id string, d *schema.ResourceData, communicator chan interface{}) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		bmsgetoptions := &vpcv1.GetBareMetalServerOptions{
			ID: &id,
		}
		bms, response, err := client.GetBareMetalServer(bmsgetoptions)
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] Error getting Bare Metal Server: %s\n%s", err, response)
		}
		d.Set(isBareMetalServerStatus, *bms.Status)

		select {
		case data := <-communicator:
			return nil, "", data.(error)
		default:
			fmt.Println("no message sent")
		}

		if *bms.Status == "running" || *bms.Status == "failed" {
			// let know the isRestartStartAction() to stop
			close(communicator)
			if *bms.Status == "failed" {
				bmsStatusReason := bms.StatusReasons

				//set the status reasons
				if bms.StatusReasons != nil {
					statusReasonsList := make([]map[string]interface{}, 0)
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
					d.Set(isBareMetalServerStatusReasons, statusReasonsList)
				}

				out, err := json.MarshalIndent(bmsStatusReason, "", "    ")
				if err != nil {
					return bms, *bms.Status, fmt.Errorf("[ERROR] The Bare Metal Server (%s) went into failed state during the operation \n [WARNING] Running terraform apply again will remove the tainted bare metal server and attempt to create the bare metal server again replacing the previous configuration", *bms.ID)
				}
				return bms, *bms.Status, fmt.Errorf("[ERROR] Bare Metal Server (%s) went into failed state during the operation \n (%+v) \n [WARNING] Running terraform apply again will remove the tainted Bare Metal Server and attempt to create the Bare Metal Server again replacing the previous configuration", *bms.ID, string(out))
			}
			return bms, *bms.Status, nil

		}
		return bms, *bms.Status, nil
	}
}

func isWaitForBareMetalServerActionStop(bmsC *vpcv1.VpcV1, timeout time.Duration, id string, d *schema.ResourceData) (interface{}, error) {
	communicator := make(chan interface{})
	stateConf := &resource.StateChangeConf{
		Pending: []string{isBareMetalServerStatusRunning, isBareMetalServerStatusPending, isBareMetalServerActionStatusStopping},
		Target:  []string{isBareMetalServerActionStatusStopped, isBareMetalServerStatusFailed, ""},
		Refresh: func() (interface{}, string, error) {
			getbmsoptions := &vpcv1.GetBareMetalServerOptions{
				ID: &id,
			}
			bms, response, err := bmsC.GetBareMetalServer(getbmsoptions)
			if err != nil {
				return nil, "", fmt.Errorf("[ERROR] Error Getting Bare Metal Server: %s\n%s", err, response)
			}
			select {
			case data := <-communicator:
				return nil, "", data.(error)
			default:
				fmt.Println("no message sent")
			}
			if *bms.Status == isBareMetalServerStatusFailed {
				// let know the isRestartStopAction() to stop
				close(communicator)
				return bms, *bms.Status, fmt.Errorf("[ERROR] The  Bare Metal Server %s failed to stop: %v", id, err)
			}
			return bms, *bms.Status, nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isBareMetalServerRestartStopAction(bmsC *vpcv1.VpcV1, id string, d *schema.ResourceData, forceTimeout int, communicator chan interface{}) {
	subticker := time.NewTicker(time.Duration(forceTimeout) * time.Minute)
	for {
		select {

		case <-subticker.C:
			log.Println("Bare Metal Server is still in stopping state, retrying to stop with -force")
			actiontype := "hard"
			createbmssactoptions := &vpcv1.StopBareMetalServerOptions{
				ID:   &id,
				Type: &actiontype,
			}
			response, err := bmsC.StopBareMetalServer(createbmssactoptions)
			if err != nil {
				communicator <- fmt.Errorf("[ERROR] Error retrying Bare Metal Server action stop: %s\n%s", err, response)
				return
			}
		case <-communicator:
			// indicates refresh func is reached target and not proceed with the thread)
			subticker.Stop()
			return

		}
	}
}

func isBareMetalServerStart(bmsC *vpcv1.VpcV1, id string, d *schema.ResourceData, forceTimeout int) (interface{}, error) {
	createbmsactoptions := &vpcv1.StartBareMetalServerOptions{
		ID: &id,
	}
	response, err := bmsC.StartBareMetalServer(createbmsactoptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil, nil
		}
		return nil, fmt.Errorf("[ERROR] Error creating Bare Metal Server action start : %s\n%s", err, response)
	}
	_, err = isWaitForBareMetalServerAvailable(bmsC, d.Id(), d.Timeout(schema.TimeoutUpdate), d)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
func isBareMetalServerStop(bmsC *vpcv1.VpcV1, id string, d *schema.ResourceData, forceTimeout int) (interface{}, error) {
	stoppingType := "soft"
	createbmsactoptions := &vpcv1.StopBareMetalServerOptions{
		ID:   &id,
		Type: &stoppingType,
	}
	response, err := bmsC.StopBareMetalServer(createbmsactoptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil, nil
		}
		return nil, fmt.Errorf("[ERROR] Error creating Bare Metal Server Action stop: %s\n%s", err, response)
	}
	_, err = isWaitForBareMetalServerActionStop(bmsC, d.Timeout(schema.TimeoutUpdate), d.Id(), d)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
func isBareMetalServerRestart(bmsC *vpcv1.VpcV1, id string, d *schema.ResourceData, forceTimeout int) (interface{}, error) {
	createbmsactoptions := &vpcv1.RestartBareMetalServerOptions{
		ID: &id,
	}
	response, err := bmsC.RestartBareMetalServer(createbmsactoptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil, nil
		}
		return nil, fmt.Errorf("[ERROR] Error creating Bare Metal Server action restart: %s\n%s", err, response)
	}
	_, err = isWaitForBareMetalServerAvailable(bmsC, d.Id(), d.Timeout(schema.TimeoutUpdate), d)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func resourceIBMBMSNicSet(v interface{}) int {
	var buf bytes.Buffer
	a := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", a["subnet"].(string)))
	// buf.WriteString(fmt.Sprintf("%s-", a["name"].(string)))
	buf.WriteString(fmt.Sprintf("%d-", a["vlan"].(int)))
	buf.WriteString(fmt.Sprintf("%v-", a["allowed_vlans"].(*schema.Set)))
	return conns.String(buf.String())
}

func resourceStopServerIfRunning(id, stoppingType string, d *schema.ResourceData, context context.Context, sess *vpcv1.VpcV1, isServerStopped bool) (bool, error) {
	getBmsOptions := &vpcv1.GetBareMetalServerOptions{
		ID: &id,
	}
	bms, response, err := sess.GetBareMetalServerWithContext(context, getBmsOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return isServerStopped, nil
		}
		return isServerStopped, fmt.Errorf("[ERROR] Error Getting Bare Metal Server (%s): %s\n%s", id, err, response)
	}
	if *bms.Status == "running" {

		options := &vpcv1.StopBareMetalServerOptions{
			ID:   bms.ID,
			Type: &stoppingType,
		}

		response, err := sess.StopBareMetalServerWithContext(context, options)
		if err != nil && response != nil && response.StatusCode != 204 {
			return isServerStopped, fmt.Errorf("[ERROR] Error stopping Bare Metal Server (%s): %s\n%s", id, err, response)
		}
		isServerStopped = true
		isWaitForBareMetalServerActionStop(sess, d.Timeout(schema.TimeoutDelete), id, d)
	}
	return isServerStopped, nil
}

func resourceStartServerIfStopped(id, stoppingType string, d *schema.ResourceData, context context.Context, sess *vpcv1.VpcV1, isServerStopped bool) (bool, error) {
	getBmsOptions := &vpcv1.GetBareMetalServerOptions{
		ID: &id,
	}
	bms, response, err := sess.GetBareMetalServerWithContext(context, getBmsOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return isServerStopped, nil
		}
		return isServerStopped, fmt.Errorf("[ERROR] Error Getting Bare Metal Server (%s): %s\n%s", id, err, response)
	}
	if *bms.Status == "stopped" {

		createbmsactoptions := &vpcv1.StartBareMetalServerOptions{
			ID: &id,
		}
		response, err := sess.StartBareMetalServer(createbmsactoptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return isServerStopped, nil
			}
			return isServerStopped, fmt.Errorf("[ERROR] Error creating Bare Metal Server action start : %s\n%s", err, response)
		}
		isServerStopped = true
		_, err = isWaitForBareMetalServerAvailable(sess, d.Id(), d.Timeout(schema.TimeoutUpdate), d)
		if err != nil {
			return isServerStopped, err
		}
	}
	return isServerStopped, nil
}

func resourceIBMIsBareMetalServerMapToBareMetalServerTrustedPlatformModulePrototype(modelMap map[string]interface{}) (*vpcv1.BareMetalServerTrustedPlatformModulePrototype, error) {
	model := &vpcv1.BareMetalServerTrustedPlatformModulePrototype{}
	// if modelMap[isBareMetalServerTrustedPlatformModuleEnabled] != nil {
	// 	model.Enabled = core.BoolPtr(modelMap[isBareMetalServerTrustedPlatformModuleEnabled].(bool))
	// }
	if modelMap[isBareMetalServerTrustedPlatformModuleMode] != nil && modelMap[isBareMetalServerTrustedPlatformModuleMode].(string) != "" {
		model.Mode = core.StringPtr(modelMap[isBareMetalServerTrustedPlatformModuleMode].(string))
	}
	return model, nil
}

func resourceIBMIsBareMetalServerBareMetalServerTrustedPlatformModulePrototypeToMap(model *vpcv1.BareMetalServerTrustedPlatformModule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Enabled != nil {
		modelMap[isBareMetalServerTrustedPlatformModuleEnabled] = model.Enabled
	}
	if model.Mode != nil {
		modelMap[isBareMetalServerTrustedPlatformModuleMode] = model.Mode
	}
	if model.SupportedModes != nil {
		modelMap[isBareMetalServerTrustedPlatformModuleSupportedModes] = model.SupportedModes
	}
	return modelMap, nil
}

func resourceIBMIsBareMetalServerBareMetalServerNetworkAttachmentReferenceToMap(model *vpcv1.BareMetalServerNetworkAttachmentReference, na vpcv1.BareMetalServerNetworkAttachmentIntf, instanceC *vpcv1.VpcV1) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Deleted != nil {
		deletedMap, err := resourceIBMIsBareMetalServerBareMetalServerNetworkAttachmentReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	vniMap := make(map[string]interface{})
	vniid := ""
	switch reflect.TypeOf(na).String() {
	case "*vpcv1.BareMetalServerNetworkAttachmentByPci":
		{
			vna := na.(*vpcv1.BareMetalServerNetworkAttachmentByPci)
			if vna.AllowedVlans != nil {
				var out = make([]interface{}, len(vna.AllowedVlans))
				for i, v := range vna.AllowedVlans {
					out[i] = int(v)
				}
				modelMap["allowed_vlans"] = schema.NewSet(schema.HashInt, out)
			}
			if vna.VirtualNetworkInterface != nil {
				vniid = *vna.VirtualNetworkInterface.ID
				vniMap["id"] = vniid
				vniMap["name"] = vna.VirtualNetworkInterface.Name
				vniMap["resource_type"] = vna.VirtualNetworkInterface.ResourceType
			}
		}
	case "*vpcv1.BareMetalServerNetworkAttachmentByVlan":
		{
			vna := na.(*vpcv1.BareMetalServerNetworkAttachmentByVlan)
			if vna.Vlan != nil {
				modelMap["vlan"] = *vna.Vlan
			}
			if vna.AllowToFloat != nil {
				modelMap["allow_to_float"] = *vna.AllowToFloat
			}
			if vna.VirtualNetworkInterface != nil {
				vniid = *vna.VirtualNetworkInterface.ID
				vniMap["id"] = vniid
				vniMap["name"] = vna.VirtualNetworkInterface.Name
				vniMap["resource_type"] = vna.VirtualNetworkInterface.ResourceType
			}
		}
	default:
		{
			vna := na.(*vpcv1.BareMetalServerNetworkAttachment)
			if vna.VirtualNetworkInterface != nil {
				vniid = *vna.VirtualNetworkInterface.ID
				vniMap["id"] = vniid
				vniMap["name"] = vna.VirtualNetworkInterface.Name
				vniMap["resource_type"] = vna.VirtualNetworkInterface.ResourceType
			}
			if vna.AllowedVlans != nil {
				var out = make([]interface{}, len(vna.AllowedVlans))
				for i, v := range vna.AllowedVlans {
					out[i] = int(v)
				}
				modelMap["allowed_vlans"] = schema.NewSet(schema.HashInt, out)
			}
			if vna.Vlan != nil {
				modelMap["vlan"] = *vna.Vlan
			}
			if vna.AllowToFloat != nil {
				modelMap["allow_to_float"] = *vna.AllowToFloat
			}

		}
	}

	getVirtualNetworkInterfaceOptions := &vpcv1.GetVirtualNetworkInterfaceOptions{
		ID: &vniid,
	}
	vniDetails, response, err := instanceC.GetVirtualNetworkInterface(getVirtualNetworkInterfaceOptions)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error on GetInstanceNetworkAttachment in instance : %s\n%s", err, response)
	}
	vniMap["allow_ip_spoofing"] = vniDetails.AllowIPSpoofing
	vniMap["auto_delete"] = vniDetails.AutoDelete
	vniMap["enable_infrastructure_nat"] = vniDetails.EnableInfrastructureNat
	vniMap["resource_group"] = vniDetails.ResourceGroup.ID
	vniMap["protocol_state_filtering_mode"] = vniDetails.ProtocolStateFilteringMode
	primaryipId := *vniDetails.PrimaryIP.ID
	if !core.IsNil(vniDetails.Ips) {
		ips := []map[string]interface{}{}
		for _, ipsItem := range vniDetails.Ips {
			if *ipsItem.ID != primaryipId {
				ipsItemMap, err := resourceIBMIsVirtualNetworkInterfaceReservedIPReferenceToMap(&ipsItem, true)
				if err != nil {
					return nil, err
				}
				ips = append(ips, ipsItemMap)
			}
		}
		vniMap["ips"] = ips
	}
	primaryIPMap, err := resourceIBMIsBareMetalServerReservedIPReferenceToMap(model.PrimaryIP)
	if err != nil {
		return modelMap, err
	}
	vniMap["primary_ip"] = []map[string]interface{}{primaryIPMap}

	if !core.IsNil(vniDetails.SecurityGroups) {
		securityGroups := make([]string, 0)
		for _, securityGroupsItem := range vniDetails.SecurityGroups {
			if securityGroupsItem.ID != nil {
				securityGroups = append(securityGroups, *securityGroupsItem.ID)
			}
		}
		vniMap["security_groups"] = securityGroups
	}
	if model.ResourceType != nil {
		modelMap["resource_type"] = *model.ResourceType
	}
	if model.Subnet != nil {
		vniMap["subnet"] = *model.Subnet.ID
	}
	modelMap["virtual_network_interface"] = []map[string]interface{}{vniMap}
	return modelMap, nil
}

func resourceIBMIsBareMetalServerBareMetalServerNetworkAttachmentReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func resourceIBMIsBareMetalServerReservedIPReferenceToMap(model *vpcv1.ReservedIPReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["address"] = model.Address
	if model.Deleted != nil {
		deletedMap, err := resourceIBMIsBareMetalServerReservedIPReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = model.Href
	modelMap["reserved_ip"] = model.ID
	modelMap["name"] = model.Name
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}

func resourceIBMIsBareMetalServerReservedIPReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func resourceIBMIsBareMetalServerMapToBareMetalServerPrimaryNetworkAttachmentPrototype(allowipspoofing, autodelete, enablenat string, d *schema.ResourceData, modelMap map[string]interface{}) (*vpcv1.BareMetalServerPrimaryNetworkAttachmentPrototype, error) {
	model := &vpcv1.BareMetalServerPrimaryNetworkAttachmentPrototype{}
	interface_type := "pci"
	if modelMap["allowed_vlans"] != nil && modelMap["allowed_vlans"].(*schema.Set).Len() > 0 {
		allowedVlans := []int64{}
		for _, allowedVlansItem := range modelMap["allowed_vlans"].(*schema.Set).List() {
			allowedVlans = append(allowedVlans, int64(allowedVlansItem.(int)))
		}
		model.AllowedVlans = allowedVlans
		interface_type = "pci"
	}
	if modelMap["interface_type"].(string) != "" {
		interface_type = modelMap["interface_type"].(string)
	}
	model.InterfaceType = &interface_type
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	VirtualNetworkInterfaceModel, err := resourceIBMIsBareMetalServerMapToVirtualNetworkInterfacePrototypeAttachmentContext(allowipspoofing, autodelete, enablenat, d, modelMap["virtual_network_interface"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.VirtualNetworkInterface = VirtualNetworkInterfaceModel
	return model, nil
}
func resourceIBMIsBareMetalServerMapToVirtualNetworkInterfaceIPsReservedIPPrototype(modelMap map[string]interface{}) (vpcv1.VirtualNetworkInterfaceIPPrototypeIntf, error) {
	model := &vpcv1.VirtualNetworkInterfaceIPPrototype{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	if modelMap["address"] != nil && modelMap["address"].(string) != "" {
		model.Address = core.StringPtr(modelMap["address"].(string))
	}
	if modelMap["auto_delete"] != nil {
		model.AutoDelete = core.BoolPtr(modelMap["auto_delete"].(bool))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	return model, nil
}
func resourceIBMIsBareMetalServerMapToVirtualNetworkInterfacePrimaryIPReservedIPPrototype(modelMap map[string]interface{}) (vpcv1.VirtualNetworkInterfacePrimaryIPPrototypeIntf, error) {
	model := &vpcv1.VirtualNetworkInterfacePrimaryIPPrototype{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	if modelMap["address"] != nil && modelMap["address"].(string) != "" {
		model.Address = core.StringPtr(modelMap["address"].(string))
	}
	if modelMap["auto_delete"] != nil {
		model.AutoDelete = core.BoolPtr(modelMap["auto_delete"].(bool))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	return model, nil
}
func resourceIBMIsBareMetalServerMapToVirtualNetworkInterfacePrototypeAttachmentContext(allowipspoofing, autodelete, enablenat string, d *schema.ResourceData, modelMap map[string]interface{}) (vpcv1.BareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterfaceIntf, error) {
	model := &vpcv1.BareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterface{}
	if _, ok := d.GetOkExists(allowipspoofing); ok && modelMap["allow_ip_spoofing"] != nil {
		model.AllowIPSpoofing = core.BoolPtr(modelMap["allow_ip_spoofing"].(bool))
	}
	if _, ok := d.GetOkExists(autodelete); ok && modelMap["auto_delete"] != nil {
		model.AutoDelete = core.BoolPtr(modelMap["auto_delete"].(bool))
	}
	if _, ok := d.GetOkExists(enablenat); ok && modelMap["enable_infrastructure_nat"] != nil {
		model.EnableInfrastructureNat = core.BoolPtr(modelMap["enable_infrastructure_nat"].(bool))
	}
	if modelMap["protocol_state_filtering_mode"] != nil {
		if pStateFilteringInt, ok := modelMap["protocol_state_filtering_mode"]; ok {
			protocolStateFilteringMode := pStateFilteringInt.(string)
			if protocolStateFilteringMode != "" {
				model.ProtocolStateFilteringMode = core.StringPtr(protocolStateFilteringMode)
			}
		}
	}
	if modelMap["ips"] != nil && modelMap["ips"].(*schema.Set).Len() > 0 {
		ips := []vpcv1.VirtualNetworkInterfaceIPPrototypeIntf{}
		for _, ipsItem := range modelMap["ips"].(*schema.Set).List() {
			ipsItemModel, err := resourceIBMIsBareMetalServerMapToVirtualNetworkInterfaceIPsReservedIPPrototype(ipsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			ips = append(ips, ipsItemModel)
		}
		model.Ips = ips
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["primary_ip"] != nil && len(modelMap["primary_ip"].([]interface{})) > 0 {
		PrimaryIPModel, err := resourceIBMIsBareMetalServerMapToVirtualNetworkInterfacePrimaryIPReservedIPPrototype(modelMap["primary_ip"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.PrimaryIP = PrimaryIPModel
	}
	if modelMap["resource_group"] != nil && modelMap["resource_group"].(string) != "" {
		resourceGroupId := modelMap["resource_group"].(string)
		model.ResourceGroup = &vpcv1.ResourceGroupIdentity{
			ID: &resourceGroupId,
		}
	}
	if modelMap["security_groups"] != nil {
		securityGroups := []vpcv1.SecurityGroupIdentityIntf{}
		sg := modelMap["security_groups"].(*schema.Set)
		for _, v := range sg.List() {
			value := v.(string)
			securityGroupsItem := &vpcv1.SecurityGroupIdentity{
				ID: &value,
			}
			securityGroups = append(securityGroups, securityGroupsItem)
		}
		model.SecurityGroups = securityGroups
	}
	if modelMap["subnet"] != nil && modelMap["subnet"].(string) != "" {
		subnetId := modelMap["subnet"].(string)
		model.Subnet = &vpcv1.SubnetIdentity{
			ID: &subnetId,
		}
	}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	if modelMap["crn"] != nil && modelMap["crn"].(string) != "" {
		model.CRN = core.StringPtr(modelMap["crn"].(string))
	}
	return model, nil
}
func resourceIBMIsBareMetalServerMapToSubnetIdentity(modelMap map[string]interface{}) (vpcv1.SubnetIdentityIntf, error) {
	model := &vpcv1.SubnetIdentity{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["crn"] != nil && modelMap["crn"].(string) != "" {
		model.CRN = core.StringPtr(modelMap["crn"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	return model, nil
}
func resourceIBMIsBareMetalServerMapToSecurityGroupIdentity(modelMap map[string]interface{}) (vpcv1.SecurityGroupIdentityIntf, error) {
	model := &vpcv1.SecurityGroupIdentity{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["crn"] != nil && modelMap["crn"].(string) != "" {
		model.CRN = core.StringPtr(modelMap["crn"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	return model, nil
}

func resourceIBMIsBareMetalServerMapToBareMetalServerNetworkAttachmentPrototype(allowipspoofing, allowfloat, autodelete, enablenat string, d *schema.ResourceData, modelMap map[string]interface{}) (vpcv1.BareMetalServerNetworkAttachmentPrototypeIntf, error) {
	if modelMap["vlan"] != nil && int64(modelMap["vlan"].(int)) != 0 {
		return resourceIBMIsBareMetalServerMapToBareMetalServerNetworkAttachmentPrototypeBareMetalServerNetworkAttachmentByVlanPrototype(allowipspoofing, allowfloat, autodelete, enablenat, d, modelMap)
	} else {
		return resourceIBMIsBareMetalServerMapToBareMetalServerNetworkAttachmentPrototypeBareMetalServerNetworkAttachmentByPciPrototype(allowipspoofing, autodelete, enablenat, d, modelMap)
	}
}

func resourceIBMIsBareMetalServerMapToBareMetalServerNetworkAttachmentPrototypeBareMetalServerNetworkAttachmentByVlanPrototype(allowipspoofing, allowfloat, autodelete, enablenat string, d *schema.ResourceData, modelMap map[string]interface{}) (*vpcv1.BareMetalServerNetworkAttachmentPrototypeBareMetalServerNetworkAttachmentByVlanPrototype, error) {
	model := &vpcv1.BareMetalServerNetworkAttachmentPrototypeBareMetalServerNetworkAttachmentByVlanPrototype{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	VirtualNetworkInterfaceModel, err := resourceIBMIsBareMetalServerMapToVirtualNetworkInterfacePrototypeAttachmentContext(allowipspoofing, autodelete, enablenat, d, modelMap["virtual_network_interface"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}

	model.VirtualNetworkInterface = VirtualNetworkInterfaceModel
	if _, ok := d.GetOkExists(allowfloat); ok && modelMap["allow_to_float"] != nil {
		model.AllowToFloat = core.BoolPtr(modelMap["allow_to_float"].(bool))
	}
	model.InterfaceType = core.StringPtr("vlan")
	model.Vlan = core.Int64Ptr(int64(modelMap["vlan"].(int)))
	return model, nil
}

func resourceIBMIsBareMetalServerMapToBareMetalServerNetworkAttachmentPrototypeBareMetalServerNetworkAttachmentByPciPrototype(allowipspoofing, autodelete, enablenat string, d *schema.ResourceData, modelMap map[string]interface{}) (*vpcv1.BareMetalServerNetworkAttachmentPrototypeBareMetalServerNetworkAttachmentByPciPrototype, error) {
	model := &vpcv1.BareMetalServerNetworkAttachmentPrototypeBareMetalServerNetworkAttachmentByPciPrototype{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	VirtualNetworkInterfaceModel, err := resourceIBMIsBareMetalServerMapToVirtualNetworkInterfacePrototypeAttachmentContext(allowipspoofing, autodelete, enablenat, d, modelMap["virtual_network_interface"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.VirtualNetworkInterface = VirtualNetworkInterfaceModel
	if modelMap["allowed_vlans"] != nil {
		allowedVlans := []int64{}
		for _, allowedVlansItem := range modelMap["allowed_vlans"].(*schema.Set).List() {
			allowedVlans = append(allowedVlans, int64(allowedVlansItem.(int)))
		}
		model.AllowedVlans = allowedVlans
	}
	model.InterfaceType = core.StringPtr("pci")
	return model, nil
}

func findNetworkAttachmentDifferences(oldList, newList []interface{}, bareMetalServerId string, sess *vpcv1.VpcV1, d *schema.ResourceData) ([]vpcv1.DeleteBareMetalServerNetworkAttachmentOptions, []vpcv1.CreateBareMetalServerNetworkAttachmentOptions, bool, error) {
	var wg sync.WaitGroup
	wg.Add(3)

	var listToDelete []vpcv1.DeleteBareMetalServerNetworkAttachmentOptions
	var listToAdd []vpcv1.CreateBareMetalServerNetworkAttachmentOptions
	var err error
	var serverRestartNeeded bool

	go func() {
		listToDelete, serverRestartNeeded = compareRemovedNacs(oldList, newList, bareMetalServerId)
		wg.Done()
	}()

	go func() {
		listToAdd, serverRestartNeeded = compareAddedNacs(oldList, newList, bareMetalServerId)
		wg.Done()
	}()

	go func() {
		err = compareModifiedNacs(oldList, newList, bareMetalServerId, sess, d)
		wg.Done()
	}()

	wg.Wait()
	return listToDelete, listToAdd, serverRestartNeeded, err
}
func compareRemovedNacs(oldList, newList []interface{}, bareMetalServerId string) ([]vpcv1.DeleteBareMetalServerNetworkAttachmentOptions, bool) {
	var removed []vpcv1.DeleteBareMetalServerNetworkAttachmentOptions
	newListMap := make(map[string]struct{})
	restartNeeded := false
	for _, newListitem := range newList {
		newListitemMap := newListitem.(map[string]interface{})
		// list of ids in new list
		if newListitemMap["id"] != nil && newListitemMap["id"].(string) != "" {
			newListMap[newListitemMap["id"].(string)] = struct{}{}
		}
	}
	// find the ids missing in the oldList, difference are the ones to be removed
	for _, oldListitem1 := range oldList {
		oldListitemMap := oldListitem1.(map[string]interface{})
		if _, exists := newListMap[oldListitemMap["id"].(string)]; !exists {
			deleteNac := &vpcv1.DeleteBareMetalServerNetworkAttachmentOptions{
				BareMetalServerID: &bareMetalServerId,
				ID:                core.StringPtr(oldListitemMap["id"].(string)),
			}
			if oldListitemMap["interface_type"].(string) == "pci" {
				restartNeeded = true
			}
			removed = append(removed, *deleteNac)
		}
	}
	return removed, restartNeeded
}

func compareAddedNacs(oldList, newList []interface{}, bareMetalServerId string) ([]vpcv1.CreateBareMetalServerNetworkAttachmentOptions, bool) {
	var added []vpcv1.CreateBareMetalServerNetworkAttachmentOptions

	restartNeeded := false
	// the nac(s) which dont have the id are to be created
	for _, newListitem := range newList {
		newListitemmMap := newListitem.(map[string]interface{})
		// if _, exists := oldListMap[newListitemmMap["name"].(string)]; !exists {
		if (newListitemmMap["id"] == nil) || (newListitemmMap["id"] != nil && newListitemmMap["id"].(string) == "") {
			addNac := &vpcv1.CreateBareMetalServerNetworkAttachmentOptions{
				BareMetalServerID: &bareMetalServerId,
			}
			if newListitemmMap["vlan"] != nil && newListitemmMap["vlan"].(int) != 0 {
				vlanId := int64(newListitemmMap["vlan"].(int))
				interfaceType := "vlan"
				nacAttPrototype := &vpcv1.BareMetalServerNetworkAttachmentPrototypeBareMetalServerNetworkAttachmentByVlanPrototype{
					Vlan:          core.Int64Ptr(vlanId),
					InterfaceType: &interfaceType,
				}
				name := newListitemmMap["name"].(string)
				if name != "" {
					nacAttPrototype.Name = &name
				}
				if newListitemmMap["virtual_network_interface"] != nil {
					newListItemVniMap := newListitemmMap["virtual_network_interface"].([]interface{})[0].(map[string]interface{})
					virtualNetworkInterface := &vpcv1.BareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterface{}
					if newListItemVniMap["allow_ip_spoofing"] != nil {
						virtualNetworkInterface.AllowIPSpoofing = core.BoolPtr(newListItemVniMap["allow_ip_spoofing"].(bool))
					}
					if newListItemVniMap["auto_delete"] != nil {
						virtualNetworkInterface.AutoDelete = core.BoolPtr(newListItemVniMap["auto_delete"].(bool))
					}
					if newListItemVniMap["enable_infrastructure_nat"] != nil {
						virtualNetworkInterface.EnableInfrastructureNat = core.BoolPtr(newListItemVniMap["enable_infrastructure_nat"].(bool))
					}
					if newListItemVniMap["ips"] != nil && newListItemVniMap["ips"].(*schema.Set).Len() > 0 {
						ips := []vpcv1.VirtualNetworkInterfaceIPPrototypeIntf{}
						for _, ipsItem := range newListItemVniMap["ips"].(*schema.Set).List() {
							ipsItemModelMap := ipsItem.(map[string]interface{})
							idIntf := ipsItemModelMap["id"]
							if idIntf != nil && idIntf.(string) != "" {
								ipsItemModel := &vpcv1.VirtualNetworkInterfaceIPPrototype{
									ID: core.StringPtr(idIntf.(string)),
								}
								ips = append(ips, ipsItemModel)
							}
						}
						virtualNetworkInterface.Ips = ips
					}
					if newListItemVniMap["name"] != nil && newListItemVniMap["name"].(string) != "" {
						virtualNetworkInterface.Name = core.StringPtr(newListItemVniMap["name"].(string))
					}
					if newListItemVniMap["primary_ip"] != nil && len(newListItemVniMap["primary_ip"].([]interface{})) > 0 {
						primaryIPMapModel := newListItemVniMap["primary_ip"].([]interface{})[0].(map[string]interface{})
						primaryIPModel := &vpcv1.VirtualNetworkInterfacePrimaryIPPrototype{}
						if primaryIPMapModel["id"] != nil && primaryIPMapModel["id"].(string) != "" {
							primaryIPModel.ID = core.StringPtr(primaryIPMapModel["id"].(string))
						}
						if primaryIPMapModel["href"] != nil && primaryIPMapModel["href"].(string) != "" {
							primaryIPModel.Href = core.StringPtr(primaryIPMapModel["href"].(string))
						}
						if primaryIPMapModel["address"] != nil && primaryIPMapModel["address"].(string) != "" {
							primaryIPModel.Address = core.StringPtr(primaryIPMapModel["address"].(string))
						}
						if primaryIPMapModel["auto_delete"] != nil {
							primaryIPModel.AutoDelete = core.BoolPtr(primaryIPMapModel["auto_delete"].(bool))
						}
						if primaryIPMapModel["name"] != nil && primaryIPMapModel["name"].(string) != "" {
							primaryIPModel.Name = core.StringPtr(primaryIPMapModel["name"].(string))
						}
						virtualNetworkInterface.PrimaryIP = primaryIPModel
					}
					if newListItemVniMap["protocol_state_filtering_mode"] != nil && newListItemVniMap["protocol_state_filtering_mode"].(string) != "" {
						virtualNetworkInterface.ProtocolStateFilteringMode = core.StringPtr(newListItemVniMap["protocol_state_filtering_mode"].(string))
					}
					if newListItemVniMap["resource_group"] != nil && newListItemVniMap["resource_group"].(string) != "" {

						virtualNetworkInterface.ResourceGroup = &vpcv1.ResourceGroupIdentity{
							ID: core.StringPtr(newListItemVniMap["resource_group"].(string)),
						}
					}
					if newListItemVniMap["security_groups"] != nil && newListItemVniMap["security_groups"].(*schema.Set).Len() > 0 {
						securityGroups := []vpcv1.SecurityGroupIdentityIntf{}
						for _, securityGroupsItem := range newListItemVniMap["security_groups"].(*schema.Set).List() {
							securityGroupsItemModel := &vpcv1.SecurityGroupIdentity{
								ID: core.StringPtr(securityGroupsItem.(string)),
							}
							securityGroups = append(securityGroups, securityGroupsItemModel)
						}
						virtualNetworkInterface.SecurityGroups = securityGroups
					}
					if newListItemVniMap["subnet"] != nil && newListItemVniMap["subnet"].(string) != "" {
						virtualNetworkInterface.Subnet = &vpcv1.SubnetIdentity{
							ID: core.StringPtr(newListItemVniMap["subnet"].(string)),
						}
					}
					if newListItemVniMap["id"] != nil && newListItemVniMap["id"].(string) != "" {
						virtualNetworkInterface.ID = core.StringPtr(newListItemVniMap["id"].(string))
					}
					nacAttPrototype.VirtualNetworkInterface = virtualNetworkInterface
				}
				addNac.BareMetalServerNetworkAttachmentPrototype = nacAttPrototype
			} else {
				restartNeeded = true
				name := newListitemmMap["name"].(string)
				interfaceType := "pci"
				nacAttPrototype := &vpcv1.BareMetalServerNetworkAttachmentPrototypeBareMetalServerNetworkAttachmentByPciPrototype{
					InterfaceType: &interfaceType,
				}
				if name != "" {
					nacAttPrototype.Name = &name
				}
				if newListitemmMap["virtual_network_interface"] != nil {
					newListVniitemmMap := newListitemmMap["virtual_network_interface"].([]interface{})[0].(map[string]interface{})
					virtualNetworkInterface := &vpcv1.BareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterface{}
					if newListVniitemmMap["allow_ip_spoofing"] != nil {
						virtualNetworkInterface.AllowIPSpoofing = core.BoolPtr(newListVniitemmMap["allow_ip_spoofing"].(bool))
					}
					if newListVniitemmMap["auto_delete"] != nil {
						virtualNetworkInterface.AutoDelete = core.BoolPtr(newListVniitemmMap["auto_delete"].(bool))
					}
					if newListVniitemmMap["enable_infrastructure_nat"] != nil {
						virtualNetworkInterface.EnableInfrastructureNat = core.BoolPtr(newListVniitemmMap["enable_infrastructure_nat"].(bool))
					}
					if newListVniitemmMap["ips"] != nil && newListVniitemmMap["ips"].(*schema.Set).Len() > 0 {
						ips := []vpcv1.VirtualNetworkInterfaceIPPrototypeIntf{}
						for _, ipsItem := range newListVniitemmMap["ips"].(*schema.Set).List() {
							ipsItemModelMap := ipsItem.(map[string]interface{})
							idIntf := ipsItemModelMap["id"]
							if idIntf != nil && idIntf.(string) != "" {
								ipsItemModel := &vpcv1.VirtualNetworkInterfaceIPPrototype{
									ID: core.StringPtr(idIntf.(string)),
								}
								ips = append(ips, ipsItemModel)
							}
						}
						virtualNetworkInterface.Ips = ips
					}
					if newListVniitemmMap["name"] != nil && newListVniitemmMap["name"].(string) != "" {
						virtualNetworkInterface.Name = core.StringPtr(newListVniitemmMap["name"].(string))
					}
					if newListVniitemmMap["primary_ip"] != nil && len(newListVniitemmMap["primary_ip"].([]interface{})) > 0 {
						primaryIPMapModel := newListVniitemmMap["primary_ip"].([]interface{})[0].(map[string]interface{})
						primaryIPModel := &vpcv1.VirtualNetworkInterfacePrimaryIPPrototype{}
						if primaryIPMapModel["id"] != nil && primaryIPMapModel["id"].(string) != "" {
							primaryIPModel.ID = core.StringPtr(primaryIPMapModel["id"].(string))
						}
						if primaryIPMapModel["href"] != nil && primaryIPMapModel["href"].(string) != "" {
							primaryIPModel.Href = core.StringPtr(primaryIPMapModel["href"].(string))
						}
						if primaryIPMapModel["address"] != nil && primaryIPMapModel["address"].(string) != "" {
							primaryIPModel.Address = core.StringPtr(primaryIPMapModel["address"].(string))
						}
						if primaryIPMapModel["auto_delete"] != nil {
							primaryIPModel.AutoDelete = core.BoolPtr(primaryIPMapModel["auto_delete"].(bool))
						}
						if primaryIPMapModel["name"] != nil && primaryIPMapModel["name"].(string) != "" {
							primaryIPModel.Name = core.StringPtr(primaryIPMapModel["name"].(string))
						}
						virtualNetworkInterface.PrimaryIP = primaryIPModel
					}
					if newListVniitemmMap["protocol_state_filtering_mode"] != nil && newListVniitemmMap["protocol_state_filtering_mode"].(string) != "" {
						virtualNetworkInterface.ProtocolStateFilteringMode = core.StringPtr(newListVniitemmMap["protocol_state_filtering_mode"].(string))
					}
					if newListVniitemmMap["resource_group"] != nil && newListVniitemmMap["resource_group"].(string) != "" {

						virtualNetworkInterface.ResourceGroup = &vpcv1.ResourceGroupIdentity{
							ID: core.StringPtr(newListVniitemmMap["resource_group"].(string)),
						}
					}
					if newListVniitemmMap["security_groups"] != nil && newListVniitemmMap["security_groups"].(*schema.Set).Len() > 0 {
						securityGroups := []vpcv1.SecurityGroupIdentityIntf{}
						for _, securityGroupsItem := range newListVniitemmMap["security_groups"].(*schema.Set).List() {
							securityGroupsItemModel := &vpcv1.SecurityGroupIdentity{
								ID: core.StringPtr(securityGroupsItem.(string)),
							}
							securityGroups = append(securityGroups, securityGroupsItemModel)
						}
						virtualNetworkInterface.SecurityGroups = securityGroups
					}
					if newListVniitemmMap["subnet"] != nil && newListVniitemmMap["subnet"].(string) != "" {
						virtualNetworkInterface.Subnet = &vpcv1.SubnetIdentity{
							ID: core.StringPtr(newListVniitemmMap["subnet"].(string)),
						}
					}
					if newListVniitemmMap["id"] != nil && newListVniitemmMap["id"].(string) != "" {
						virtualNetworkInterface.ID = core.StringPtr(newListVniitemmMap["id"].(string))
					}
					nacAttPrototype.VirtualNetworkInterface = virtualNetworkInterface
				}
				allowedVlansIntf := newListitemmMap["allowed_vlans"]
				if allowedVlansIntf != nil && allowedVlansIntf.(*schema.Set).Len() > 0 {
					allowedVlans := []int64{}
					for _, allowedVlansItem := range allowedVlansIntf.(*schema.Set).List() {
						allowedVlans = append(allowedVlans, int64(allowedVlansItem.(int)))
					}
					nacAttPrototype.AllowedVlans = allowedVlans
				}
				addNac.BareMetalServerNetworkAttachmentPrototype = nacAttPrototype
			}
			added = append(added, *addNac)
		}
	}

	return added, restartNeeded
}

func compareModifiedNacs(oldList, newList []interface{}, bareMetalServerId string, sess *vpcv1.VpcV1, d *schema.ResourceData) error {
	list2Map := make(map[string]interface{})

	for _, newListitem := range newList {
		newListitemMap := newListitem.(map[string]interface{})
		if newListitemMap["id"] != nil && newListitemMap["id"] != "" {
			list2Map[newListitemMap["id"].(string)] = newListitem
		}
	}

	for _, oldListitem1 := range oldList {
		oldListitemMap := oldListitem1.(map[string]interface{})
		if oldListitemMap["id"] != nil {
			if oldListitem2, exists := list2Map[oldListitemMap["id"].(string)]; exists {
				s1 := oldListitem1.(map[string]interface{})
				s2 := oldListitem2.(map[string]interface{})
				id := s1["id"].(string)
				modilfiedNac := &vpcv1.UpdateBareMetalServerNetworkAttachmentOptions{
					BareMetalServerID: &bareMetalServerId,
					ID:                &id,
				}
				s1AllowedVlan := s1["allowed_vlans"]
				s2AllowedVlan := s2["allowed_vlans"]
				s1Name := s1["name"]
				s2Name := s2["name"]
				s1Vni := s1["virtual_network_interface"]
				s2Vni := s2["virtual_network_interface"]
				bmsNacPatchModel := &vpcv1.BareMetalServerNetworkAttachmentPatch{}
				hasChanged := false
				if s1AllowedVlan != nil && s2AllowedVlan != nil {
					if !s1AllowedVlan.(*schema.Set).Equal(s2AllowedVlan.(*schema.Set)) {
						hasChanged = true
						allowedVlansList := s2AllowedVlan.(*schema.Set).List()
						allowedVlans := make([]int64, 0, len(allowedVlansList))
						for _, k := range allowedVlansList {
							allowedVlans = append(allowedVlans, int64(k.(int)))
						}
						bmsNacPatchModel.AllowedVlans = allowedVlans
					}
				}
				if s1Name != nil && s2Name != nil {
					if s1Name.(string) != s2Name.(string) {
						hasChanged = true
						bmsNacPatchModel.Name = core.StringPtr(s2Name.(string))
					}
				}
				if hasChanged {
					bmsNacPatch, _ := bmsNacPatchModel.AsPatch()
					modilfiedNac.BareMetalServerNetworkAttachmentPatch = bmsNacPatch
					_, res, err := sess.UpdateBareMetalServerNetworkAttachment(modilfiedNac)
					if err != nil {
						return fmt.Errorf("%s/n%v", err, res)
					}
				}
				if s1Vni != nil && s2Vni != nil {
					s1VniMap := s1Vni.([]interface{})[0].(map[string]interface{})
					s2VniMap := s2Vni.([]interface{})[0].(map[string]interface{})
					vniId := s1VniMap["id"].(string)
					s1vniMapAIS := s1VniMap["allow_ip_spoofing"]
					s1vniMapAD := s1VniMap["auto_delete"]
					s1vniMapEIN := s1VniMap["enable_infrastructure_nat"]
					s1vniMapIPS := s1VniMap["ips"]
					s1vniMapName := s1VniMap["name"]
					s1vniMapSG := s1VniMap["security_groups"]
					s2vniMapAIS := s2VniMap["allow_ip_spoofing"]
					s2vniMapAD := s2VniMap["auto_delete"]
					s2vniMapEIN := s2VniMap["enable_infrastructure_nat"]
					s2vniMapIPS := s2VniMap["ips"]
					s2vniMapName := s2VniMap["name"]
					s2vniMapSG := s2VniMap["security_groups"]
					s1vniPSFM := s1VniMap["protocol_state_filtering_mode"]
					s2vniPSFM := s2VniMap["protocol_state_filtering_mode"]
					vniUpdateOptions := &vpcv1.UpdateVirtualNetworkInterfaceOptions{
						ID: &vniId,
					}
					hasChanged := false
					vniPatch := &vpcv1.VirtualNetworkInterfacePatch{}
					if s1vniMapAIS != nil && s2vniMapAIS != nil && s1vniMapAIS.(bool) != s2vniMapAIS.(bool) {
						vniPatch.AllowIPSpoofing = core.BoolPtr(s2vniMapAIS.(bool))
						hasChanged = true
					}
					if s1vniMapAD != nil && s2vniMapAD != nil && s1vniMapAD.(bool) != s2vniMapAD.(bool) {
						vniPatch.AutoDelete = core.BoolPtr(s2vniMapAD.(bool))
						hasChanged = true
					}
					if s1vniMapEIN != nil && s2vniMapEIN != nil && s1vniMapEIN.(bool) != s2vniMapEIN.(bool) {
						vniPatch.EnableInfrastructureNat = core.BoolPtr(s2vniMapEIN.(bool))
						hasChanged = true
					}
					if s1vniMapName != nil && s2vniMapName != nil && s1vniMapName.(string) != s2vniMapName.(string) {
						vniPatch.Name = core.StringPtr(s2vniMapName.(string))
						hasChanged = true
					}
					if s1vniPSFM != nil && s2vniPSFM != nil && s1vniPSFM.(string) != s2vniPSFM.(string) {
						vniPatch.ProtocolStateFilteringMode = core.StringPtr(s2vniPSFM.(string))
						hasChanged = true
					}
					if hasChanged {
						virtualNetworkInterfacePatchAsPatch, err := vniPatch.AsPatch()
						vniUpdateOptions.VirtualNetworkInterfacePatch = virtualNetworkInterfacePatchAsPatch
						_, res, err := sess.UpdateVirtualNetworkInterface(vniUpdateOptions)
						if err != nil {
							return fmt.Errorf("%s/n%v", err, res)
						}
					}
					if s1vniMapIPS != nil && s2vniMapIPS != nil && !s1vniMapIPS.(*schema.Set).Equal(s2vniMapIPS.(*schema.Set)) {

						os := s1vniMapIPS.(*schema.Set)
						ns := s2vniMapIPS.(*schema.Set)
						var oldset, newset *schema.Set

						var out = make([]interface{}, ns.Len(), ns.Len())
						for i, nA := range ns.List() {
							newPack := nA.(map[string]interface{})
							out[i] = newPack["reserved_ip"].(string)
						}
						newset = schema.NewSet(schema.HashString, out)

						out = make([]interface{}, os.Len(), os.Len())
						for i, oA := range os.List() {
							oldPack := oA.(map[string]interface{})
							out[i] = oldPack["reserved_ip"].(string)
						}
						oldset = schema.NewSet(schema.HashString, out)

						remove := flex.ExpandStringList(oldset.Difference(newset).List())
						add := flex.ExpandStringList(newset.Difference(oldset).List())

						if add != nil && len(add) > 0 {
							for _, ipItem := range add {
								if ipItem != "" {

									addVirtualNetworkInterfaceIPOptions := &vpcv1.AddVirtualNetworkInterfaceIPOptions{}
									addVirtualNetworkInterfaceIPOptions.SetVirtualNetworkInterfaceID(vniId)
									addVirtualNetworkInterfaceIPOptions.SetID(ipItem)
									_, response, err := sess.AddVirtualNetworkInterfaceIP(addVirtualNetworkInterfaceIPOptions)
									if err != nil {
										log.Printf("[DEBUG] AddVirtualNetworkInterfaceIP failed in VirtualNetworkInterface patch during BareMetalServer nac patch %s\n%s", err, response)
										return fmt.Errorf("AddVirtualNetworkInterfaceIP failed in VirtualNetworkInterface patch during BareMetalServer nac patch %s\n%s", err, response)
									}
								}
							}
						}
						if remove != nil && len(remove) > 0 {
							for _, ipItem := range remove {
								if ipItem != "" {

									removeVirtualNetworkInterfaceIPOptions := &vpcv1.RemoveVirtualNetworkInterfaceIPOptions{}
									removeVirtualNetworkInterfaceIPOptions.SetVirtualNetworkInterfaceID(vniId)
									removeVirtualNetworkInterfaceIPOptions.SetID(ipItem)
									response, err := sess.RemoveVirtualNetworkInterfaceIP(removeVirtualNetworkInterfaceIPOptions)
									if err != nil {
										log.Printf("[DEBUG] RemoveVirtualNetworkInterfaceIP failed in VirtualNetworkInterface patch during BareMetalServer nac patch %s\n%s", err, response)
										return fmt.Errorf("RemoveVirtualNetworkInterfaceIP failed in VirtualNetworkInterface patch during BareMetalServer nac patch %s\n%s", err, response)
									}
								}
							}
						}

					}
					if s1vniMapSG != nil && s2vniMapSG != nil && !s1vniMapSG.(*schema.Set).Equal(s2vniMapSG.(*schema.Set)) {

						ov := s1vniMapSG.(*schema.Set)
						nv := s2vniMapSG.(*schema.Set)
						remove := flex.ExpandStringList(ov.Difference(nv).List())
						add := flex.ExpandStringList(nv.Difference(ov).List())
						if len(add) > 0 {
							for i := range add {
								createsgnicoptions := &vpcv1.CreateSecurityGroupTargetBindingOptions{
									SecurityGroupID: &add[i],
									ID:              &vniId,
								}
								_, response, err := sess.CreateSecurityGroupTargetBinding(createsgnicoptions)
								if err != nil {
									return (fmt.Errorf("[ERROR] Error while creating security group %q for virtual network interface %s\n%s: %q", add[i], d.Id(), err, response))
								}
								_, err = isWaitForVirtualNetworkInterfaceAvailable(sess, vniId, d.Timeout(schema.TimeoutUpdate))
								if err != nil {
									return (err)
								}
							}

						}
						if len(remove) > 0 {
							for i := range remove {
								deletesgnicoptions := &vpcv1.DeleteSecurityGroupTargetBindingOptions{
									SecurityGroupID: &remove[i],
									ID:              &vniId,
								}
								response, err := sess.DeleteSecurityGroupTargetBinding(deletesgnicoptions)
								if err != nil {
									return (fmt.Errorf("[ERROR] Error while removing security group %q for virtual network interface %s\n%s: %q", remove[i], d.Id(), err, response))
								}
								_, err = isWaitForVirtualNetworkInterfaceAvailable(sess, vniId, d.Timeout(schema.TimeoutUpdate))
								if err != nil {
									return (err)
								}
							}
						}

					}
				}
			}
		}
	}
	return nil
}
