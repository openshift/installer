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
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Description: "Primary Network interface info",
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

			isBareMetalServerNetworkInterfaces: {
				Type:     schema.TypeSet,
				Optional: true,
				Set:      resourceIBMBMSNicSet,
				Computed: true,
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

			isBareMetalServerKeys: {
				Type:             schema.TypeSet,
				Required:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				Set:              schema.HashString,
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "SSH key Ids for the bare metal server",
			},

			isBareMetalServerImage: {
				Type:        schema.TypeString,
				ForceNew:    true,
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
				ForceNew:    true,
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
	options := &vpcv1.CreateBareMetalServerOptions{}
	var imageStr string
	if image, ok := d.GetOk(isBareMetalServerImage); ok {
		imageStr = image.(string)
	}

	// enable secure boot

	if _, ok := d.GetOkExists(isBareMetalServerEnableSecureBoot); ok {
		options.SetEnableSecureBoot(d.Get(isBareMetalServerEnableSecureBoot).(bool))
	}

	// trusted_platform_module

	if _, ok := d.GetOk(isBareMetalServerTrustedPlatformModule); ok {
		trustedPlatformModuleModel, err := resourceIBMIsBareMetalServerMapToBareMetalServerTrustedPlatformModulePrototype(d.Get("trusted_platform_module.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		options.SetTrustedPlatformModule(trustedPlatformModuleModel)
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

	bms, response, err := sess.CreateBareMetalServerWithContext(context, options)
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
	bmsBootTargetIntf := bms.BootTarget.(*vpcv1.BareMetalServerBootTarget)
	bmsBootTarget := bmsBootTargetIntf.ID
	d.Set(isBareMetalServerBootTarget, bmsBootTarget)
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
	isServerStopped := false
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
