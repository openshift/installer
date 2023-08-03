// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isInstanceName                    = "name"
	IsInstanceCRN                     = "crn"
	isInstanceKeys                    = "keys"
	isInstanceTags                    = "tags"
	isInstanceBootVolumeTags          = "tags"
	isInstanceNetworkInterfaces       = "network_interfaces"
	isInstancePrimaryNetworkInterface = "primary_network_interface"
	isInstanceNicName                 = "name"
	isInstanceProfile                 = "profile"
	isInstanceNicPortSpeed            = "port_speed"
	isInstanceNicAllowIPSpoofing      = "allow_ip_spoofing"
	isInstanceNicPrimaryIpv4Address   = "primary_ipv4_address"
	isInstanceNicSecondaryAddress     = "secondary_addresses"
	isInstanceNicSecurityGroups       = "security_groups"
	isInstanceNicSubnet               = "subnet"
	isInstanceNicFloatingIP           = "floating_ip"
	isInstanceNicFloatingIPs          = "floating_ips"
	isInstanceUserData                = "user_data"
	isInstanceVolumes                 = "volumes"
	isInstanceVPC                     = "vpc"
	isInstanceZone                    = "zone"
	isInstanceBootVolume              = "boot_volume"
	isInstanceVolumeSnapshot          = "snapshot"
	isInstanceSourceTemplate          = "instance_template"
	isInstanceBandwidth               = "bandwidth"
	isInstanceTotalVolumeBandwidth    = "total_volume_bandwidth"
	isInstanceTotalNetworkBandwidth   = "total_network_bandwidth"
	isInstanceVolAttVolAutoDelete     = "auto_delete_volume"
	isInstanceVolAttVolBillingTerm    = "billing_term"
	isInstanceImage                   = "image"
	isInstanceCPU                     = "vcpu"
	isInstanceCPUArch                 = "architecture"
	isInstanceCPUCores                = "cores"
	isInstanceCPUCount                = "count"
	isInstanceCPUManufacturer         = "manufacturer"
	isInstanceGpu                     = "gpu"
	isInstanceGpuCores                = "cores"
	isInstanceGpuCount                = "count"
	isInstanceGpuManufacturer         = "manufacturer"
	isInstanceGpuMemory               = "memory"
	isInstanceGpuModel                = "model"
	isInstanceMemory                  = "memory"
	isInstanceDisks                   = "disks"
	isInstanceDedicatedHost           = "dedicated_host"
	isInstanceStatus                  = "status"
	isInstanceStatusReasons           = "status_reasons"
	isInstanceStatusReasonsCode       = "code"
	isInstanceStatusReasonsMessage    = "message"
	isInstanceStatusReasonsMoreInfo   = "more_info"
	isEnableCleanDelete               = "wait_before_delete"
	isInstanceProvisioning            = "provisioning"
	isInstanceProvisioningDone        = "done"
	isInstanceAvailable               = "available"
	isInstanceDeleting                = "deleting"
	isInstanceDeleteDone              = "done"
	isInstanceFailed                  = "failed"

	isInstanceStatusRestarting           = "restarting"
	isInstanceStatusStarting             = "starting"
	isInstanceActionStatusStopping       = "stopping"
	isInstanceActionStatusStopped        = "stopped"
	isInstanceStatusPending              = "pending"
	isInstanceStatusRunning              = "running"
	isInstanceStatusFailed               = "failed"
	isInstanceAvailablePolicyHostFailure = "availability_policy_host_failure"

	isInstanceBootAttachmentName       = "name"
	isInstanceBootVolumeId             = "volume_id"
	isInstanceBootSize                 = "size"
	isInstanceBootIOPS                 = "iops"
	isInstanceBootEncryption           = "encryption"
	isInstanceBootProfile              = "profile"
	isInstanceAction                   = "action"
	isInstanceVolumeAttachments        = "volume_attachments"
	isInstanceVolumeAttaching          = "attaching"
	isInstanceVolumeAttached           = "attached"
	isInstanceVolumeDetaching          = "detaching"
	isInstanceResourceGroup            = "resource_group"
	isInstanceLifecycleReasons         = "lifecycle_reasons"
	isInstanceLifecycleState           = "lifecycle_state"
	isInstanceLifecycleReasonsCode     = "code"
	isInstanceLifecycleReasonsMessage  = "message"
	isInstanceLifecycleReasonsMoreInfo = "more_info"

	isInstanceCatalogOffering            = "catalog_offering"
	isInstanceCatalogOfferingOfferingCrn = "offering_crn"
	isInstanceCatalogOfferingVersionCrn  = "version_crn"

	isPlacementTargetDedicatedHost      = "dedicated_host"
	isPlacementTargetDedicatedHostGroup = "dedicated_host_group"
	isInstancePlacementTarget           = "placement_target"
	isPlacementTargetPlacementGroup     = "placement_group"

	isInstanceDefaultTrustedProfileAutoLink = "default_trusted_profile_auto_link"
	isInstanceDefaultTrustedProfileTarget   = "default_trusted_profile_target"
	isInstanceMetadataServiceEnabled        = "metadata_service_enabled"

	isInstanceAccessTags                  = "access_tags"
	isInstanceUserTagType                 = "user"
	isInstanceAccessTagType               = "access"
	isInstanceMetadataService             = "metadata_service"
	isInstanceMetadataServiceEnabled1     = "enabled"
	isInstanceMetadataServiceProtocol     = "protocol"
	isInstanceMetadataServiceRespHopLimit = "response_hop_limit"
)

func ResourceIBMISInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceIBMisInstanceCreate,
		Read:   resourceIBMisInstanceRead,
		Update: resourceIBMisInstanceUpdate,
		Delete: resourceIBMisInstanceDelete,
		Exists: resourceIBMisInstanceExists,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) (result []*schema.ResourceData, err error) {
				log.Printf("[INFO] Instance (%s) importing", d.Id())
				id := d.Id()
				instanceC, err := vpcClient(meta)
				if err != nil {
					return nil, err
				}
				getinsOptions := &vpcv1.GetInstanceOptions{
					ID: &id,
				}
				instance, response, err := instanceC.GetInstance(getinsOptions)
				if err != nil {
					if response != nil && response.StatusCode == 404 {
						d.SetId("")
						return nil, nil
					}
					return nil, fmt.Errorf("[ERROR] Error Getting Instance: %s\n%s", err, response)
				}
				var volumes []string
				volumes = make([]string, 0)
				if instance.VolumeAttachments != nil {
					for _, volume := range instance.VolumeAttachments {
						if volume.Volume != nil && *volume.Volume.ID != *instance.BootVolumeAttachment.Volume.ID {
							volumes = append(volumes, *volume.Volume.ID)
						}
					}
				}
				d.Set(isInstanceVolumes, flex.NewStringSet(schema.HashString, volumes))
				return []*schema.ResourceData{d}, nil
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			customdiff.Sequence(
				func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
					return flex.InstanceProfileValidate(diff)
				}),
			customdiff.Sequence(
				func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
					return flex.ResourceTagsCustomizeDiff(diff)
				}),
			customdiff.Sequence(
				func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
					return flex.ResourceValidateAccessTags(diff, v)
				}),
		),

		Schema: map[string]*schema.Schema{
			isInstanceAvailablePolicyHostFailure: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The availability policy to use for this virtual server instance",
			},

			isInstanceName: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validate.InvokeValidator("ibm_is_instance", isInstanceName),
				Description:  "Instance name",
			},
			isInstanceVPC: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: "VPC id",
			},
			IsInstanceCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Crn for this Instance",
			},

			isInstanceSourceTemplate: {
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				AtLeastOneOf:  []string{isInstanceImage, isInstanceSourceTemplate, "boot_volume.0.snapshot", "catalog_offering.0.offering_crn", "catalog_offering.0.version_crn", "boot_volume.0.volume_id"},
				ConflictsWith: []string{"boot_volume.0.snapshot", "boot_volume.0.volume_id"},
				Description:   "Id of the instance template",
			},
			isInstanceZone: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Computed:    true,
				Optional:    true,
				Description: "Zone name",
			},

			isInstanceProfile: {
				Type:        schema.TypeString,
				ForceNew:    false,
				Computed:    true,
				Optional:    true,
				Description: "Profile info",
			},
			isInstanceDefaultTrustedProfileAutoLink: {
				Type:         schema.TypeBool,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				RequiredWith: []string{isInstanceDefaultTrustedProfileTarget},
				Description:  "If set to `true`, the system will create a link to the specified `target` trusted profile during instance creation. Regardless of whether a link is created by the system or manually using the IAM Identity service, it will be automatically deleted when the instance is deleted.",
			},
			isInstanceDefaultTrustedProfileTarget: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The unique identifier or CRN of the default IAM trusted profile to use for this virtual server instance.",
			},
			isPlacementTargetDedicatedHost: {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{isPlacementTargetDedicatedHostGroup, isPlacementTargetPlacementGroup},
				Description:   "Unique Identifier of the Dedicated Host where the instance will be placed",
			},

			isPlacementTargetDedicatedHostGroup: {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{isPlacementTargetDedicatedHost, isPlacementTargetPlacementGroup},
				Description:   "Unique Identifier of the Dedicated Host Group where the instance will be placed",
			},

			isPlacementTargetPlacementGroup: {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{isPlacementTargetDedicatedHost, isPlacementTargetDedicatedHostGroup},
				Description:   "Unique Identifier of the Placement Group for restricting the placement of the instance",
			},

			isInstanceTotalVolumeBandwidth: {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_instance", isInstanceTotalVolumeBandwidth),
				Description:  "The amount of bandwidth (in megabits per second) allocated exclusively to instance storage volumes",
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

			isInstanceKeys: {
				Type:             schema.TypeSet,
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				Set:              schema.HashString,
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "SSH key Ids for the instance",
			},

			isInstanceTags: {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_instance", "tags")},
				Set:         flex.ResourceIBMVPCHash,
				Description: "list of tags for the instance",
			},

			isInstanceAccessTags: {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_instance", "accesstag")},
				Set:         flex.ResourceIBMVPCHash,
				Description: "list of access tags for the instance",
			},

			isEnableCleanDelete: {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          true,
				DiffSuppressFunc: suppressEnableCleanDelete,
				Description:      "Enables stopping of instance before deleting and waits till deletion is complete",
			},

			isInstanceAction: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_instance", isInstanceAction),
				Description:  "Enables stopping of instance before deleting and waits till deletion is complete",
			},

			isInstanceActionForce: {
				Type:         schema.TypeBool,
				Optional:     true,
				RequiredWith: []string{isInstanceAction},
				Default:      false,
				Description:  "If set to true, the action will be forced immediately, and all queued actions deleted. Ignored for the start action.",
			},

			isInstanceVolumeAttachments: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume_crn": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			isInstanceCatalogOffering: {
				Type:        schema.TypeList,
				MinItems:    0,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Description: "The catalog offering or offering version to use when provisioning this virtual server instance. If an offering is specified, the latest version of that offering will be used. The specified offering or offering version may be in a different account in the same enterprise, subject to IAM policies.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isInstanceCatalogOfferingOfferingCrn: {
							Type:          schema.TypeString,
							Optional:      true,
							ForceNew:      true,
							ConflictsWith: []string{"catalog_offering.0.version_crn"},
							RequiredWith:  []string{isInstanceZone, isInstancePrimaryNetworkInterface, isInstanceKeys, isInstanceVPC, isInstanceProfile},
							Description:   "Identifies a catalog offering by a unique CRN property",
						},
						isInstanceCatalogOfferingVersionCrn: {
							Type:          schema.TypeString,
							Optional:      true,
							ForceNew:      true,
							ConflictsWith: []string{"catalog_offering.0.offering_crn"},
							RequiredWith:  []string{isInstanceZone, isInstancePrimaryNetworkInterface, isInstanceKeys, isInstanceVPC, isInstanceProfile},
							Description:   "Identifies a version of a catalog offering by a unique CRN property",
						},
					},
				},
			},

			isInstancePrimaryNetworkInterface: {
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "Primary Network interface info",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						isInstanceNicAllowIPSpoofing: {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Indicates whether IP spoofing is allowed on this interface.",
						},
						isInstanceNicName: {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						isInstanceNicPortSpeed: {
							Type:             schema.TypeInt,
							Optional:         true,
							Computed:         true,
							DiffSuppressFunc: flex.ApplyOnce,
							Deprecated:       "This field is deprected",
						},
						isInstanceNicPrimaryIpv4Address: {
							Type:          schema.TypeString,
							ForceNew:      true,
							Optional:      true,
							Computed:      true,
							ConflictsWith: []string{"primary_network_interface.0.primary_ip.0.address"},
							Deprecated:    "primary_ipv4_address is deprecated and support will be removed. Use primary_ip instead",
						},
						isInstanceNicPrimaryIP: {
							Type:        schema.TypeList,
							MinItems:    0,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "The primary IP address to bind to the network interface. This can be specified using an existing reserved IP, or a prototype object for a new reserved IP.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isInstanceNicReservedIpAddress: {
										Type:          schema.TypeString,
										Computed:      true,
										ForceNew:      true,
										Optional:      true,
										ConflictsWith: []string{"primary_network_interface.0.primary_ipv4_address"},
										Description:   "The IP address to reserve, which must not already be reserved on the subnet.",
									},
									isInstanceNicReservedIpHref: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this reserved IP",
									},
									isInstanceNicReservedIpAutoDelete: {
										Type:        schema.TypeBool,
										Optional:    true,
										Computed:    true,
										Description: "Indicates whether this reserved IP member will be automatically deleted when either target is deleted, or the reserved IP is unbound.",
									},
									isInstanceNicReservedIpName: {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The user-defined name for this reserved IP. If unspecified, the name will be a hyphenated list of randomly-selected words. Names must be unique within the subnet the reserved IP resides in. ",
									},
									isInstanceNicReservedIpId: {
										Type:          schema.TypeString,
										Optional:      true,
										ForceNew:      true,
										ConflictsWith: []string{"primary_network_interface.0.primary_ipv4_address", "primary_network_interface.0.primary_ip.0.address"},
										Computed:      true,
										Description:   "Identifies a reserved IP by a unique property.",
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
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						isInstanceNicSubnet: {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			isInstanceNetworkInterfaces: {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						isInstanceNicAllowIPSpoofing: {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Indicates whether IP spoofing is allowed on this interface.",
						},
						isInstanceNicName: {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						isInstanceNicPrimaryIpv4Address: {
							Type:       schema.TypeString,
							ForceNew:   true,
							Optional:   true,
							Deprecated: "primary_ipv4_address is deprecated and support will be removed. Use primary_ip instead",
							Computed:   true,
						},
						isInstanceNicPrimaryIP: {
							Type:        schema.TypeList,
							MinItems:    0,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "The primary IP address to bind to the network interface. This can be specified using an existing reserved IP, or a prototype object for a new reserved IP.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isInstanceNicReservedIpAddress: {
										Type:        schema.TypeString,
										Computed:    true,
										ForceNew:    true,
										Optional:    true,
										Description: "The IP address to reserve, which must not already be reserved on the subnet.",
									},
									isInstanceNicReservedIpAutoDelete: {
										Type:        schema.TypeBool,
										Optional:    true,
										Computed:    true,
										Description: "Indicates whether this reserved IP member will be automatically deleted when either target is deleted, or the reserved IP is unbound.",
									},
									isInstanceNicReservedIpHref: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this reserved IP",
									},
									isInstanceNicReservedIpName: {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The user-defined name for this reserved IP. If unspecified, the name will be a hyphenated list of randomly-selected words. Names must be unique within the subnet the reserved IP resides in. ",
									},
									isInstanceNicReservedIpId: {
										Type:        schema.TypeString,
										Optional:    true,
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
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						isInstanceNicSubnet: {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			isInstanceUserData: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "User data given for the instance",
			},

			isInstanceImage: {
				Type:          schema.TypeString,
				ForceNew:      true,
				Computed:      true,
				Optional:      true,
				ConflictsWith: []string{"boot_volume.0.snapshot", "catalog_offering.0.offering_crn", "catalog_offering.0.version_crn", "boot_volume.0.volume_id"},
				AtLeastOneOf:  []string{isInstanceImage, isInstanceSourceTemplate, "boot_volume.0.snapshot", "catalog_offering.0.offering_crn", "catalog_offering.0.version_crn", "boot_volume.0.volume_id"},
				RequiredWith:  []string{isInstanceZone, isInstancePrimaryNetworkInterface, isInstanceKeys, isInstanceVPC, isInstanceProfile},
				Description:   "image id",
			},

			isInstanceBootVolume: {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isInstanceBootVolumeId: {
							Type:          schema.TypeString,
							Optional:      true,
							ForceNew:      true,
							Computed:      true,
							RequiredWith:  []string{isInstanceZone, isInstancePrimaryNetworkInterface, isInstanceProfile, isInstanceKeys, isInstanceVPC},
							AtLeastOneOf:  []string{isInstanceImage, isInstanceSourceTemplate, "boot_volume.0.volume_id", "boot_volume.0.snapshot"},
							ConflictsWith: []string{isInstanceImage, isInstanceSourceTemplate, "boot_volume.0.snapshot", "boot_volume.0.name", "boot_volume.0.encryption", "catalog_offering.0.offering_crn", "catalog_offering.0.version_crn"},
							Description:   "The unique identifier for this volume",
						},
						isInstanceVolAttVolAutoDelete: {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Auto delete boot volume along with instance",
						},
						isInstanceBootAttachmentName: {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validate.InvokeValidator("ibm_is_instance", isInstanceBootAttachmentName),
						},

						isInstanceVolumeSnapshot: {
							Type:          schema.TypeString,
							RequiredWith:  []string{isInstanceZone, isInstancePrimaryNetworkInterface, isInstanceProfile, isInstanceKeys, isInstanceVPC},
							AtLeastOneOf:  []string{isInstanceImage, isInstanceSourceTemplate, "boot_volume.0.snapshot", "catalog_offering.0.offering_crn", "catalog_offering.0.version_crn", "boot_volume.0.volume_id"},
							ConflictsWith: []string{isInstanceImage, isInstanceSourceTemplate, "catalog_offering.0.offering_crn", "catalog_offering.0.version_crn", "boot_volume.0.volume_id"},
							Optional:      true,
							ForceNew:      true,
							Computed:      true,
						},
						isInstanceBootEncryption: {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: flex.ApplyOnce,
							Computed:         true,
						},
						isInstanceBootSize: {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validate.InvokeValidator("ibm_is_instance", isInstanceBootSize),
						},
						isInstanceBootIOPS: {
							Type:     schema.TypeInt,
							Computed: true,
						},
						isInstanceBootProfile: {
							Type:     schema.TypeString,
							Computed: true,
						},
						isInstanceBootVolumeTags: {
							Type:        schema.TypeSet,
							Optional:    true,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_instance", "tags")},
							Set:         flex.ResourceIBMVPCHash,
							Description: "UserTags for the volume instance",
						},
					},
				},
			},

			isInstanceVolumes: {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of volumes",
			},
			isInstanceVolAttVolAutoDelete: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Auto delete volume along with instance",
			},

			isInstanceResourceGroup: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: "Instance resource group",
			},

			isInstanceCPU: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isInstanceCPUArch: {
							Type:     schema.TypeString,
							Computed: true,
						},
						isInstanceCPUCount: {
							Type:     schema.TypeInt,
							Computed: true,
						},
						isInstanceCPUManufacturer: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VCPU manufacturer",
						},
					},
				},
			},

			isInstanceGpu: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The virtual server instance GPU configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isInstanceGpuCount: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of GPUs assigned to the instance",
						},
						isInstanceGpuMemory: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The overall amount of GPU memory in GiB (gibibytes)",
						},
						isInstanceGpuManufacturer: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The GPU manufacturer",
						},
						isInstanceGpuModel: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The GPU model",
						},
					},
				},
			},

			isInstanceMemory: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Instance memory",
			},

			isInstanceStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "instance status",
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
			isInstanceMetadataServiceEnabled: {
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{isInstanceMetadataService},
				Deprecated:    "Use metadata_service instead",
				Description:   "Indicates whether the metadata service endpoint is available to the virtual server instance",
			},

			isInstanceMetadataService: {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				MinItems:      1,
				MaxItems:      1,
				ConflictsWith: []string{isInstanceMetadataServiceEnabled},
				Description:   "The metadata service configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isInstanceMetadataServiceEnabled1: {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Indicates whether the metadata service endpoint will be available to the virtual server instance",
						},

						isInstanceMetadataServiceProtocol: {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							Description:  "The communication protocol to use for the metadata service endpoint. Applies only when the metadata service is enabled.",
							ValidateFunc: validate.InvokeValidator("ibm_is_instance", isInstanceMetadataServiceProtocol),
						},

						isInstanceMetadataServiceRespHopLimit: {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							Description:  "The hop limit (IP time to live) for IP response packets from the metadata service",
							ValidateFunc: validate.InvokeValidator("ibm_is_instance", isInstanceMetadataServiceRespHopLimit),
						},
					},
				},
			},
			flex.ResourceControllerURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the IBM Cloud dashboard that can be used to explore and view details about this instance",
			},

			flex.ResourceName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the resource",
			},

			flex.ResourceCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},

			flex.ResourceStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the resource",
			},

			flex.ResourceGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group name in which resource is provisioned",
			},

			"force_recovery_time": {
				Description: "Define timeout to force the instances to start/stop in minutes.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			isInstanceDisks: &schema.Schema{
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
			isInstancePlacementTarget: &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The placement restrictions for the virtual server instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this placement target.",
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
							Description: "The URL for this placement target.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this placement target.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique user-defined name for this placement target.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of resource referenced.",
						},
					},
				},
			},
		},
	}
}

func ResourceIBMISInstanceValidator() *validate.ResourceValidator {
	actions := "stop, start, reboot"
	host_failure := "restart, stop"
	metadataServiceProtocol := "https, http"
	validateSchema := make([]validate.ValidateSchema, 0)

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isInstanceMetadataServiceRespHopLimit,
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			Optional:                   true,
			MinValue:                   "1",
			MaxValue:                   "64"})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isInstanceMetadataServiceProtocol,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              metadataServiceProtocol})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isInstanceName,
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
			Identifier:                 isInstanceTotalVolumeBandwidth,
			ValidateFunctionIdentifier: validate.IntAtLeast,
			Type:                       validate.TypeInt,
			Optional:                   true,
			MinValue:                   "500"})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isInstanceBootSize,
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			Optional:                   true,
			MinValue:                   "1",
			MaxValue:                   "250"})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isInstanceAction,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              actions})

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isInstanceBootAttachmentName,
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isInstanceAvailablePolicyHostFailure,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              host_failure})

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "accesstag",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([A-Za-z0-9_.-]|[A-Za-z0-9_.-][A-Za-z0-9_ .-]*[A-Za-z0-9_.-]):([A-Za-z0-9_.-]|[A-Za-z0-9_.-][A-Za-z0-9_ .-]*[A-Za-z0-9_.-])$`,
			MinValueLength:             1,
			MaxValueLength:             128})

	ibmISInstanceValidator := validate.ResourceValidator{ResourceName: "ibm_is_instance", Schema: validateSchema}
	return &ibmISInstanceValidator
}

func instanceCreateByImage(d *schema.ResourceData, meta interface{}, profile, name, vpcID, zone, image string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	instanceproto := &vpcv1.InstancePrototype{
		Image: &vpcv1.ImageIdentity{
			ID: &image,
		},
		Zone: &vpcv1.ZoneIdentity{
			Name: &zone,
		},
		Profile: &vpcv1.InstanceProfileIdentity{
			Name: &profile,
		},
		Name: &name,
		VPC: &vpcv1.VPCIdentity{
			ID: &vpcID,
		},
	}
	if defaultTrustedProfileTargetIntf, ok := d.GetOk(isInstanceDefaultTrustedProfileTarget); ok {
		defaultTrustedProfiletarget := defaultTrustedProfileTargetIntf.(string)

		target := &vpcv1.TrustedProfileIdentity{}
		if strings.HasPrefix(defaultTrustedProfiletarget, "crn") {
			target.CRN = &defaultTrustedProfiletarget
		} else {
			target.ID = &defaultTrustedProfiletarget
		}
		instanceproto.DefaultTrustedProfile = &vpcv1.InstanceDefaultTrustedProfilePrototype{
			Target: target,
		}

		if defaultTrustedProfileAutoLinkIntf, ok := d.GetOkExists(isInstanceDefaultTrustedProfileAutoLink); ok {
			defaultTrustedProfileAutoLink := defaultTrustedProfileAutoLinkIntf.(bool)
			instanceproto.DefaultTrustedProfile.AutoLink = &defaultTrustedProfileAutoLink
		}
	}
	if availablePolicyItem, ok := d.GetOk(isInstanceAvailablePolicyHostFailure); ok {
		hostFailure := availablePolicyItem.(string)
		instanceproto.AvailabilityPolicy = &vpcv1.InstanceAvailabilityPolicyPrototype{
			HostFailure: &hostFailure,
		}
	}

	if totalVolBandwidthIntf, ok := d.GetOk(isInstanceTotalVolumeBandwidth); ok {
		totalVolBandwidthStr := int64(totalVolBandwidthIntf.(int))
		instanceproto.TotalVolumeBandwidth = &totalVolBandwidthStr
	}
	if dHostIdInf, ok := d.GetOk(isPlacementTargetDedicatedHost); ok {
		dHostIdStr := dHostIdInf.(string)
		dHostPlaementTarget := &vpcv1.InstancePlacementTargetPrototypeDedicatedHostIdentity{
			ID: &dHostIdStr,
		}
		instanceproto.PlacementTarget = dHostPlaementTarget
	} else if dHostGrpIdInf, ok := d.GetOk(isPlacementTargetDedicatedHostGroup); ok {
		dHostGrpIdStr := dHostGrpIdInf.(string)
		dHostGrpPlaementTarget := &vpcv1.InstancePlacementTargetPrototypeDedicatedHostGroupIdentity{
			ID: &dHostGrpIdStr,
		}
		instanceproto.PlacementTarget = dHostGrpPlaementTarget
	} else if placementGroupInf, ok := d.GetOk(isPlacementTargetPlacementGroup); ok {
		placementGrpStr := placementGroupInf.(string)
		placementGrp := &vpcv1.InstancePlacementTargetPrototypePlacementGroupIdentity{
			ID: &placementGrpStr,
		}
		instanceproto.PlacementTarget = placementGrp
	}

	if boot, ok := d.GetOk(isInstanceBootVolume); ok {
		bootvol := boot.([]interface{})[0].(map[string]interface{})
		var volTemplate = &vpcv1.VolumePrototypeInstanceByImageContext{}
		name, ok := bootvol[isInstanceBootAttachmentName]
		namestr := name.(string)
		if namestr != "" && ok {
			volTemplate.Name = &namestr
		}
		sizeOk, ok := bootvol[isInstanceBootSize]
		size := sizeOk.(int)
		if size != 0 && ok {
			sizeInt64 := int64(size)
			volTemplate.Capacity = &sizeInt64
		}
		enc, ok := bootvol[isInstanceBootEncryption]
		encstr := enc.(string)
		if ok && encstr != "" {
			volTemplate.EncryptionKey = &vpcv1.EncryptionKeyIdentity{
				CRN: &encstr,
			}
		}

		volprof := "general-purpose"
		volTemplate.Profile = &vpcv1.VolumeProfileIdentity{
			Name: &volprof,
		}
		var userTags *schema.Set
		if v, ok := bootvol[isInstanceBootVolumeTags]; ok {
			userTags = v.(*schema.Set)
			if userTags != nil && userTags.Len() != 0 {
				userTagsArray := make([]string, userTags.Len())
				for i, userTag := range userTags.List() {
					userTagStr := userTag.(string)
					userTagsArray[i] = userTagStr
				}
				volTemplate.UserTags = userTagsArray
			}
		}
		deleteboolIntf := bootvol[isInstanceVolAttVolAutoDelete]
		deletebool := deleteboolIntf.(bool)
		instanceproto.BootVolumeAttachment = &vpcv1.VolumeAttachmentPrototypeInstanceByImageContext{
			DeleteVolumeOnInstanceDelete: &deletebool,
			Volume:                       volTemplate,
		}

	}

	if primnicintf, ok := d.GetOk(isInstancePrimaryNetworkInterface); ok {
		primnic := primnicintf.([]interface{})[0].(map[string]interface{})
		subnetintf, _ := primnic[isInstanceNicSubnet]
		subnetintfstr := subnetintf.(string)
		var primnicobj = &vpcv1.NetworkInterfacePrototype{}
		primnicobj.Subnet = &vpcv1.SubnetIdentity{
			ID: &subnetintfstr,
		}
		name, _ := primnic[isInstanceNicName]
		namestr := name.(string)
		if namestr != "" {
			primnicobj.Name = &namestr
		}

		// reserved ip changes

		var ipv4str, reservedIp, reservedipv4, reservedipname string
		var autodelete, okAuto bool
		ipv4, _ := primnic[isInstanceNicPrimaryIpv4Address]
		ipv4str = ipv4.(string)

		primaryIpOk, ok := primnic[isInstanceNicPrimaryIP]
		if ok && len(primaryIpOk.([]interface{})) > 0 {
			primip := primaryIpOk.([]interface{})[0].(map[string]interface{})

			reservedipok, _ := primip[isInstanceNicReservedIpId]
			reservedIp = reservedipok.(string)

			reservedipv4Ok, _ := primip[isInstanceNicReservedIpAddress]
			reservedipv4 = reservedipv4Ok.(string)

			reservedipnameOk, _ := primip[isInstanceNicReservedIpName]
			reservedipname = reservedipnameOk.(string)
			var reservedipautodeleteok interface{}
			reservedipautodeleteok, okAuto = primip[isInstanceNicReservedIpAutoDelete]
			autodelete = reservedipautodeleteok.(bool)
		}
		if ipv4str != "" && reservedipv4 != "" && ipv4str != reservedipv4 {
			return fmt.Errorf("[ERROR] Error creating instance, primary_network_interface error, use either primary_ipv4_address(%s) or primary_ip.0.address(%s)", ipv4str, reservedipv4)
		}
		if reservedIp != "" && (ipv4str != "" || reservedipv4 != "" || reservedipname != "") {
			return fmt.Errorf("[ERROR] Error creating instance, primary_network_interface error, reserved_ip(%s) is mutually exclusive with other primary_ip attributes", reservedIp)
		}
		if reservedIp != "" {
			primnicobj.PrimaryIP = &vpcv1.NetworkInterfaceIPPrototypeReservedIPIdentity{
				ID: &reservedIp,
			}
		} else {
			if ipv4str != "" || reservedipv4 != "" || reservedipname != "" || okAuto {
				primaryipobj := &vpcv1.NetworkInterfaceIPPrototypeReservedIPPrototypeNetworkInterfaceContext{}
				if ipv4str != "" {
					primaryipobj.Address = &ipv4str
				}
				if reservedipv4 != "" {
					primaryipobj.Address = &reservedipv4
				}
				if reservedipname != "" {
					primaryipobj.Name = &reservedipname
				}
				if okAuto {
					primaryipobj.AutoDelete = &autodelete
				}
				primnicobj.PrimaryIP = primaryipobj
			}
		}

		allowIPSpoofing, ok := primnic[isInstanceNicAllowIPSpoofing]
		allowIPSpoofingbool := allowIPSpoofing.(bool)
		if ok {
			primnicobj.AllowIPSpoofing = &allowIPSpoofingbool
		}
		secgrpintf, ok := primnic[isInstanceNicSecurityGroups]
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
		instanceproto.PrimaryNetworkInterface = primnicobj
	}

	if nicsintf, ok := d.GetOk(isInstanceNetworkInterfaces); ok {
		nics := nicsintf.([]interface{})
		var intfs []vpcv1.NetworkInterfacePrototype
		for _, resource := range nics {
			nic := resource.(map[string]interface{})
			nwInterface := &vpcv1.NetworkInterfacePrototype{}
			subnetintf, _ := nic[isInstanceNicSubnet]
			subnetintfstr := subnetintf.(string)
			nwInterface.Subnet = &vpcv1.SubnetIdentity{
				ID: &subnetintfstr,
			}
			name, ok := nic[isInstanceNicName]
			namestr := name.(string)
			if ok && namestr != "" {
				nwInterface.Name = &namestr
			}

			// reserved ip changes

			var ipv4str, reservedIp, reservedipv4, reservedipname string
			var autodelete, okAuto bool
			ipv4, _ := nic[isInstanceNicPrimaryIpv4Address]
			ipv4str = ipv4.(string)

			primaryIpOk, ok := nic[isInstanceNicPrimaryIP]
			if ok && len(primaryIpOk.([]interface{})) > 0 {
				primip := primaryIpOk.([]interface{})[0].(map[string]interface{})

				reservedipok, _ := primip[isInstanceNicReservedIpId]
				reservedIp = reservedipok.(string)

				reservedipv4Ok, _ := primip[isInstanceNicReservedIpAddress]
				reservedipv4 = reservedipv4Ok.(string)

				reservedipnameOk, _ := primip[isInstanceNicReservedIpName]
				reservedipname = reservedipnameOk.(string)
				var reservedipautodeleteok interface{}
				reservedipautodeleteok, okAuto = primip[isInstanceNicReservedIpAutoDelete]
				autodelete = reservedipautodeleteok.(bool)
			}
			if ipv4str != "" && reservedipv4 != "" && ipv4str != reservedipv4 {
				return fmt.Errorf("[ERROR] Error creating instance, network_interfaces error, use either primary_ipv4_address(%s) or primary_ip.0.address(%s)", ipv4str, reservedipv4)
			}
			if reservedIp != "" && (ipv4str != "" || reservedipv4 != "" || reservedipname != "") {
				return fmt.Errorf("[ERROR] Error creating instance, network_interfaces error, reserved_ip(%s) is mutually exclusive with other primary_ip attributes", reservedIp)
			}
			if reservedIp != "" {
				nwInterface.PrimaryIP = &vpcv1.NetworkInterfaceIPPrototypeReservedIPIdentity{
					ID: &reservedIp,
				}
			} else {
				if ipv4str != "" || reservedipv4 != "" || reservedipname != "" || okAuto {
					primaryipobj := &vpcv1.NetworkInterfaceIPPrototypeReservedIPPrototypeNetworkInterfaceContext{}
					if ipv4str != "" {
						primaryipobj.Address = &ipv4str
					}
					if reservedipv4 != "" {
						primaryipobj.Address = &reservedipv4
					}
					if reservedipname != "" {
						primaryipobj.Name = &reservedipname
					}
					if okAuto {
						primaryipobj.AutoDelete = &autodelete
					}
					nwInterface.PrimaryIP = primaryipobj
				}
			}
			allowIPSpoofing, ok := nic[isInstanceNicAllowIPSpoofing]
			allowIPSpoofingbool := allowIPSpoofing.(bool)
			if ok {
				nwInterface.AllowIPSpoofing = &allowIPSpoofingbool
			}
			secgrpintf, ok := nic[isInstanceNicSecurityGroups]
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
					nwInterface.SecurityGroups = secgrpobjs
				}
			}
			intfs = append(intfs, *nwInterface)
		}
		instanceproto.NetworkInterfaces = intfs
	}

	keySet := d.Get(isInstanceKeys).(*schema.Set)
	if keySet.Len() != 0 {
		keyobjs := make([]vpcv1.KeyIdentityIntf, keySet.Len())
		for i, key := range keySet.List() {
			keystr := key.(string)
			keyobjs[i] = &vpcv1.KeyIdentity{
				ID: &keystr,
			}
		}
		instanceproto.Keys = keyobjs
	}

	if userdata, ok := d.GetOk(isInstanceUserData); ok {
		userdatastr := userdata.(string)
		instanceproto.UserData = &userdatastr
	}

	if grp, ok := d.GetOk(isInstanceResourceGroup); ok {
		grpstr := grp.(string)
		instanceproto.ResourceGroup = &vpcv1.ResourceGroupIdentity{
			ID: &grpstr,
		}

	}

	metadataServiceEnabled := d.Get(isInstanceMetadataServiceEnabled).(bool)
	if metadataServiceEnabled {
		instanceproto.MetadataService = &vpcv1.InstanceMetadataServicePrototype{
			Enabled: &metadataServiceEnabled,
		}
	}

	if metadataService := GetInstanceMetadataServiceOptions(d); metadataService != nil {
		instanceproto.MetadataService = metadataService
	}

	options := &vpcv1.CreateInstanceOptions{
		InstancePrototype: instanceproto,
	}

	instance, response, err := sess.CreateInstance(options)
	if err != nil {
		log.Printf("[DEBUG] Instance err %s\n%s", err, response)
		return err
	}
	d.SetId(*instance.ID)

	log.Printf("[INFO] Instance : %s", *instance.ID)
	d.Set(isInstanceStatus, instance.Status)

	_, err = isWaitForInstanceAvailable(sess, d.Id(), d.Timeout(schema.TimeoutCreate), d)
	if err != nil {
		return err
	}

	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk(isInstanceTags); ok || v != "" {
		oldList, newList := d.GetChange(isInstanceTags)
		err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *instance.CRN, "", isInstanceUserTagType)
		if err != nil {
			log.Printf(
				"Error on create of resource instance (%s) tags: %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk(isInstanceAccessTags); ok {
		oldList, newList := d.GetChange(isInstanceAccessTags)
		err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *instance.CRN, "", isInstanceAccessTagType)
		if err != nil {
			log.Printf(
				"Error on create of resource instance (%s) access tags: %s", d.Id(), err)
		}
	}
	return nil
}
func instanceCreateByCatalogOffering(d *schema.ResourceData, meta interface{}, profile, name, vpcID, zone, image, offerringCrn, versionCrn string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	instanceproto := &vpcv1.InstancePrototypeInstanceByCatalogOffering{
		Zone: &vpcv1.ZoneIdentity{
			Name: &zone,
		},
		Profile: &vpcv1.InstanceProfileIdentity{
			Name: &profile,
		},
		Name: &name,
		VPC: &vpcv1.VPCIdentity{
			ID: &vpcID,
		},
	}
	if offerringCrn != "" {
		catalogOffering := &vpcv1.CatalogOfferingIdentityCatalogOfferingByCRN{
			CRN: &offerringCrn,
		}
		offeringPrototype := &vpcv1.InstanceCatalogOfferingPrototypeCatalogOfferingByOffering{
			Offering: catalogOffering,
		}
		instanceproto.CatalogOffering = offeringPrototype
	}
	if versionCrn != "" {
		versionOffering := &vpcv1.CatalogOfferingVersionIdentityCatalogOfferingVersionByCRN{
			CRN: &versionCrn,
		}
		versionPrototype := &vpcv1.InstanceCatalogOfferingPrototypeCatalogOfferingByVersion{
			Version: versionOffering,
		}
		instanceproto.CatalogOffering = versionPrototype
	}
	if defaultTrustedProfileTargetIntf, ok := d.GetOk(isInstanceDefaultTrustedProfileTarget); ok {
		defaultTrustedProfiletarget := defaultTrustedProfileTargetIntf.(string)

		target := &vpcv1.TrustedProfileIdentity{}
		if strings.HasPrefix(defaultTrustedProfiletarget, "crn") {
			target.CRN = &defaultTrustedProfiletarget
		} else {
			target.ID = &defaultTrustedProfiletarget
		}
		instanceproto.DefaultTrustedProfile = &vpcv1.InstanceDefaultTrustedProfilePrototype{
			Target: target,
		}

		if defaultTrustedProfileAutoLinkIntf, ok := d.GetOkExists(isInstanceDefaultTrustedProfileAutoLink); ok {
			defaultTrustedProfileAutoLink := defaultTrustedProfileAutoLinkIntf.(bool)
			instanceproto.DefaultTrustedProfile.AutoLink = &defaultTrustedProfileAutoLink
		}
	}
	if availablePolicyItem, ok := d.GetOk(isInstanceAvailablePolicyHostFailure); ok {
		hostFailure := availablePolicyItem.(string)
		instanceproto.AvailabilityPolicy = &vpcv1.InstanceAvailabilityPolicyPrototype{
			HostFailure: &hostFailure,
		}
	}

	if totalVolBandwidthIntf, ok := d.GetOk(isInstanceTotalVolumeBandwidth); ok {
		totalVolBandwidthStr := int64(totalVolBandwidthIntf.(int))
		instanceproto.TotalVolumeBandwidth = &totalVolBandwidthStr
	}
	if dHostIdInf, ok := d.GetOk(isPlacementTargetDedicatedHost); ok {
		dHostIdStr := dHostIdInf.(string)
		dHostPlaementTarget := &vpcv1.InstancePlacementTargetPrototypeDedicatedHostIdentity{
			ID: &dHostIdStr,
		}
		instanceproto.PlacementTarget = dHostPlaementTarget
	} else if dHostGrpIdInf, ok := d.GetOk(isPlacementTargetDedicatedHostGroup); ok {
		dHostGrpIdStr := dHostGrpIdInf.(string)
		dHostGrpPlaementTarget := &vpcv1.InstancePlacementTargetPrototypeDedicatedHostGroupIdentity{
			ID: &dHostGrpIdStr,
		}
		instanceproto.PlacementTarget = dHostGrpPlaementTarget
	} else if placementGroupInf, ok := d.GetOk(isPlacementTargetPlacementGroup); ok {
		placementGrpStr := placementGroupInf.(string)
		placementGrp := &vpcv1.InstancePlacementTargetPrototypePlacementGroupIdentity{
			ID: &placementGrpStr,
		}
		instanceproto.PlacementTarget = placementGrp
	}

	if boot, ok := d.GetOk(isInstanceBootVolume); ok {
		bootvol := boot.([]interface{})[0].(map[string]interface{})
		var volTemplate = &vpcv1.VolumePrototypeInstanceByImageContext{}
		name, ok := bootvol[isInstanceBootAttachmentName]
		namestr := name.(string)
		if namestr != "" && ok {
			volTemplate.Name = &namestr
		}
		sizeOk, ok := bootvol[isInstanceBootSize]
		size := sizeOk.(int)
		if size != 0 && ok {
			sizeInt64 := int64(size)
			volTemplate.Capacity = &sizeInt64
		}
		enc, ok := bootvol[isInstanceBootEncryption]
		encstr := enc.(string)
		if ok && encstr != "" {
			volTemplate.EncryptionKey = &vpcv1.EncryptionKeyIdentity{
				CRN: &encstr,
			}
		}

		volprof := "general-purpose"
		volTemplate.Profile = &vpcv1.VolumeProfileIdentity{
			Name: &volprof,
		}
		deleteboolIntf := bootvol[isInstanceVolAttVolAutoDelete]
		deletebool := deleteboolIntf.(bool)
		instanceproto.BootVolumeAttachment = &vpcv1.VolumeAttachmentPrototypeInstanceByImageContext{
			DeleteVolumeOnInstanceDelete: &deletebool,
			Volume:                       volTemplate,
		}

	}

	if primnicintf, ok := d.GetOk(isInstancePrimaryNetworkInterface); ok {
		primnic := primnicintf.([]interface{})[0].(map[string]interface{})
		subnetintf, _ := primnic[isInstanceNicSubnet]
		subnetintfstr := subnetintf.(string)
		var primnicobj = &vpcv1.NetworkInterfacePrototype{}
		primnicobj.Subnet = &vpcv1.SubnetIdentity{
			ID: &subnetintfstr,
		}
		name, _ := primnic[isInstanceNicName]
		namestr := name.(string)
		if namestr != "" {
			primnicobj.Name = &namestr
		}

		// reserved ip changes

		var ipv4str, reservedIp, reservedipv4, reservedipname string
		var autodelete, okAuto bool
		ipv4, _ := primnic[isInstanceNicPrimaryIpv4Address]
		ipv4str = ipv4.(string)

		primaryIpOk, ok := primnic[isInstanceNicPrimaryIP]
		if ok && len(primaryIpOk.([]interface{})) > 0 {
			primip := primaryIpOk.([]interface{})[0].(map[string]interface{})

			reservedipok, _ := primip[isInstanceNicReservedIpId]
			reservedIp = reservedipok.(string)

			reservedipv4Ok, _ := primip[isInstanceNicReservedIpAddress]
			reservedipv4 = reservedipv4Ok.(string)

			reservedipnameOk, _ := primip[isInstanceNicReservedIpName]
			reservedipname = reservedipnameOk.(string)
			var reservedipautodeleteok interface{}
			reservedipautodeleteok, okAuto = primip[isInstanceNicReservedIpAutoDelete]
			autodelete = reservedipautodeleteok.(bool)
		}
		if ipv4str != "" && reservedipv4 != "" && ipv4str != reservedipv4 {
			return fmt.Errorf("[ERROR] Error creating instance, primary_network_interface error, use either primary_ipv4_address(%s) or primary_ip.0.address(%s)", ipv4str, reservedipv4)
		}
		if reservedIp != "" && (ipv4str != "" || reservedipv4 != "" || reservedipname != "") {
			return fmt.Errorf("[ERROR] Error creating instance, primary_network_interface error, reserved_ip(%s) is mutually exclusive with other primary_ip attributes", reservedIp)
		}
		if reservedIp != "" {
			primnicobj.PrimaryIP = &vpcv1.NetworkInterfaceIPPrototypeReservedIPIdentity{
				ID: &reservedIp,
			}
		} else {
			if ipv4str != "" || reservedipv4 != "" || reservedipname != "" || okAuto {
				primaryipobj := &vpcv1.NetworkInterfaceIPPrototypeReservedIPPrototypeNetworkInterfaceContext{}
				if ipv4str != "" {
					primaryipobj.Address = &ipv4str
				}
				if reservedipv4 != "" {
					primaryipobj.Address = &reservedipv4
				}
				if reservedipname != "" {
					primaryipobj.Name = &reservedipname
				}
				if okAuto {
					primaryipobj.AutoDelete = &autodelete
				}
				primnicobj.PrimaryIP = primaryipobj
			}
		}

		allowIPSpoofing, ok := primnic[isInstanceNicAllowIPSpoofing]
		allowIPSpoofingbool := allowIPSpoofing.(bool)
		if ok {
			primnicobj.AllowIPSpoofing = &allowIPSpoofingbool
		}
		secgrpintf, ok := primnic[isInstanceNicSecurityGroups]
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
		instanceproto.PrimaryNetworkInterface = primnicobj
	}

	if nicsintf, ok := d.GetOk(isInstanceNetworkInterfaces); ok {
		nics := nicsintf.([]interface{})
		var intfs []vpcv1.NetworkInterfacePrototype
		for _, resource := range nics {
			nic := resource.(map[string]interface{})
			nwInterface := &vpcv1.NetworkInterfacePrototype{}
			subnetintf, _ := nic[isInstanceNicSubnet]
			subnetintfstr := subnetintf.(string)
			nwInterface.Subnet = &vpcv1.SubnetIdentity{
				ID: &subnetintfstr,
			}
			name, ok := nic[isInstanceNicName]
			namestr := name.(string)
			if ok && namestr != "" {
				nwInterface.Name = &namestr
			}

			// reserved ip changes

			var ipv4str, reservedIp, reservedipv4, reservedipname string
			var autodelete, okAuto bool
			ipv4, _ := nic[isInstanceNicPrimaryIpv4Address]
			ipv4str = ipv4.(string)

			primaryIpOk, ok := nic[isInstanceNicPrimaryIP]
			if ok && len(primaryIpOk.([]interface{})) > 0 {
				primip := primaryIpOk.([]interface{})[0].(map[string]interface{})

				reservedipok, _ := primip[isInstanceNicReservedIpId]
				reservedIp = reservedipok.(string)

				reservedipv4Ok, _ := primip[isInstanceNicReservedIpAddress]
				reservedipv4 = reservedipv4Ok.(string)

				reservedipnameOk, _ := primip[isInstanceNicReservedIpName]
				reservedipname = reservedipnameOk.(string)
				var reservedipautodeleteok interface{}
				reservedipautodeleteok, okAuto = primip[isInstanceNicReservedIpAutoDelete]
				autodelete = reservedipautodeleteok.(bool)
			}
			if ipv4str != "" && reservedipv4 != "" && ipv4str != reservedipv4 {
				return fmt.Errorf("[ERROR] Error creating instance, network_interfaces error, use either primary_ipv4_address(%s) or primary_ip.0.address(%s)", ipv4str, reservedipv4)
			}
			if reservedIp != "" && (ipv4str != "" || reservedipv4 != "" || reservedipname != "") {
				return fmt.Errorf("[ERROR] Error creating instance, network_interfaces error, reserved_ip(%s) is mutually exclusive with other primary_ip attributes", reservedIp)
			}
			if reservedIp != "" {
				nwInterface.PrimaryIP = &vpcv1.NetworkInterfaceIPPrototypeReservedIPIdentity{
					ID: &reservedIp,
				}
			} else {
				if ipv4str != "" || reservedipv4 != "" || reservedipname != "" || okAuto {
					primaryipobj := &vpcv1.NetworkInterfaceIPPrototypeReservedIPPrototypeNetworkInterfaceContext{}
					if ipv4str != "" {
						primaryipobj.Address = &ipv4str
					}
					if reservedipv4 != "" {
						primaryipobj.Address = &reservedipv4
					}
					if reservedipname != "" {
						primaryipobj.Name = &reservedipname
					}
					if okAuto {
						primaryipobj.AutoDelete = &autodelete
					}
					nwInterface.PrimaryIP = primaryipobj
				}
			}
			allowIPSpoofing, ok := nic[isInstanceNicAllowIPSpoofing]
			allowIPSpoofingbool := allowIPSpoofing.(bool)
			if ok {
				nwInterface.AllowIPSpoofing = &allowIPSpoofingbool
			}
			secgrpintf, ok := nic[isInstanceNicSecurityGroups]
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
					nwInterface.SecurityGroups = secgrpobjs
				}
			}
			intfs = append(intfs, *nwInterface)
		}
		instanceproto.NetworkInterfaces = intfs
	}

	keySet := d.Get(isInstanceKeys).(*schema.Set)
	if keySet.Len() != 0 {
		keyobjs := make([]vpcv1.KeyIdentityIntf, keySet.Len())
		for i, key := range keySet.List() {
			keystr := key.(string)
			keyobjs[i] = &vpcv1.KeyIdentity{
				ID: &keystr,
			}
		}
		instanceproto.Keys = keyobjs
	}

	if userdata, ok := d.GetOk(isInstanceUserData); ok {
		userdatastr := userdata.(string)
		instanceproto.UserData = &userdatastr
	}

	if grp, ok := d.GetOk(isInstanceResourceGroup); ok {
		grpstr := grp.(string)
		instanceproto.ResourceGroup = &vpcv1.ResourceGroupIdentity{
			ID: &grpstr,
		}

	}

	metadataServiceEnabled := d.Get(isInstanceMetadataServiceEnabled).(bool)
	if metadataServiceEnabled {
		instanceproto.MetadataService = &vpcv1.InstanceMetadataServicePrototype{
			Enabled: &metadataServiceEnabled,
		}
	}

	if metadataService := GetInstanceMetadataServiceOptions(d); metadataService != nil {
		instanceproto.MetadataService = metadataService
	}

	options := &vpcv1.CreateInstanceOptions{
		InstancePrototype: instanceproto,
	}

	instance, response, err := sess.CreateInstance(options)
	if err != nil {
		log.Printf("[DEBUG] Instance err %s\n%s", err, response)
		return err
	}
	d.SetId(*instance.ID)

	log.Printf("[INFO] Instance : %s", *instance.ID)
	d.Set(isInstanceStatus, instance.Status)

	_, err = isWaitForInstanceAvailable(sess, d.Id(), d.Timeout(schema.TimeoutCreate), d)
	if err != nil {
		return err
	}

	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk(isInstanceTags); ok || v != "" {
		oldList, newList := d.GetChange(isInstanceTags)
		err = flex.UpdateTagsUsingCRN(oldList, newList, meta, *instance.CRN)
		if err != nil {
			log.Printf(
				"Error on create of resource instance (%s) tags: %s", d.Id(), err)
		}
	}
	return nil
}

func instanceCreateByTemplate(d *schema.ResourceData, meta interface{}, profile, name, vpcID, zone, image, template string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	instanceproto := &vpcv1.InstancePrototypeInstanceBySourceTemplate{
		SourceTemplate: &vpcv1.InstanceTemplateIdentity{
			ID: &template,
		},
		Name: &name,
	}

	if defaultTrustedProfileTargetIntf, ok := d.GetOk(isInstanceDefaultTrustedProfileTarget); ok {
		defaultTrustedProfiletarget := defaultTrustedProfileTargetIntf.(string)

		target := &vpcv1.TrustedProfileIdentity{}
		if strings.HasPrefix(defaultTrustedProfiletarget, "crn") {
			target.CRN = &defaultTrustedProfiletarget
		} else {
			target.ID = &defaultTrustedProfiletarget
		}
		instanceproto.DefaultTrustedProfile = &vpcv1.InstanceDefaultTrustedProfilePrototype{
			Target: target,
		}

		if defaultTrustedProfileAutoLinkIntf, ok := d.GetOkExists(isInstanceDefaultTrustedProfileAutoLink); ok {
			defaultTrustedProfileAutoLink := defaultTrustedProfileAutoLinkIntf.(bool)
			instanceproto.DefaultTrustedProfile.AutoLink = &defaultTrustedProfileAutoLink
		}
	}
	if profile != "" {
		instanceproto.Profile = &vpcv1.InstanceProfileIdentity{
			Name: &profile,
		}
	}
	if totalVolBandwidthIntf, ok := d.GetOk(isInstanceTotalVolumeBandwidth); ok {
		totalVolBandwidthStr := int64(totalVolBandwidthIntf.(int))
		instanceproto.TotalVolumeBandwidth = &totalVolBandwidthStr
	}

	if vpcID != "" {
		instanceproto.VPC = &vpcv1.VPCIdentity{
			ID: &vpcID,
		}
	}
	if zone != "" {
		instanceproto.Zone = &vpcv1.ZoneIdentity{
			Name: &zone,
		}
	}

	if dHostIdInf, ok := d.GetOk(isPlacementTargetDedicatedHost); ok {
		dHostIdStr := dHostIdInf.(string)
		dHostPlaementTarget := &vpcv1.InstancePlacementTargetPrototypeDedicatedHostIdentity{
			ID: &dHostIdStr,
		}
		instanceproto.PlacementTarget = dHostPlaementTarget
	} else if dHostGrpIdInf, ok := d.GetOk(isPlacementTargetDedicatedHostGroup); ok {
		dHostGrpIdStr := dHostGrpIdInf.(string)
		dHostGrpPlaementTarget := &vpcv1.InstancePlacementTargetPrototypeDedicatedHostGroupIdentity{
			ID: &dHostGrpIdStr,
		}
		instanceproto.PlacementTarget = dHostGrpPlaementTarget
	} else if placementGroupInf, ok := d.GetOk(isPlacementTargetPlacementGroup); ok {
		placementGrpStr := placementGroupInf.(string)
		placementGrp := &vpcv1.InstancePlacementTargetPrototypePlacementGroupIdentity{
			ID: &placementGrpStr,
		}
		instanceproto.PlacementTarget = placementGrp
	}
	if availablePolicyItem, ok := d.GetOk(isInstanceAvailablePolicyHostFailure); ok {
		hostFailure := availablePolicyItem.(string)
		instanceproto.AvailabilityPolicy = &vpcv1.InstanceAvailabilityPolicyPrototype{
			HostFailure: &hostFailure,
		}
	}
	if boot, ok := d.GetOk(isInstanceBootVolume); ok {
		bootvol := boot.([]interface{})[0].(map[string]interface{})
		var volTemplate = &vpcv1.VolumePrototypeInstanceByImageContext{}
		name, ok := bootvol[isInstanceBootAttachmentName]
		namestr := name.(string)
		if namestr != "" && ok {
			volTemplate.Name = &namestr
		}
		sizeOk, ok := bootvol[isInstanceBootSize]
		size := sizeOk.(int)
		if size != 0 && ok {
			sizeInt64 := int64(size)
			volTemplate.Capacity = &sizeInt64
		}
		enc, ok := bootvol[isInstanceBootEncryption]
		encstr := enc.(string)
		if ok && encstr != "" {
			volTemplate.EncryptionKey = &vpcv1.EncryptionKeyIdentity{
				CRN: &encstr,
			}
		}

		volprof := "general-purpose"

		volTemplate.Profile = &vpcv1.VolumeProfileIdentity{
			Name: &volprof,
		}
		var userTags *schema.Set
		if v, ok := bootvol[isInstanceBootVolumeTags]; ok {
			userTags = v.(*schema.Set)
			if userTags != nil && userTags.Len() != 0 {
				userTagsArray := make([]string, userTags.Len())
				for i, userTag := range userTags.List() {
					userTagStr := userTag.(string)
					userTagsArray[i] = userTagStr
				}
				volTemplate.UserTags = userTagsArray
			}
		}
		deleteboolIntf := bootvol[isInstanceVolAttVolAutoDelete]
		deletebool := deleteboolIntf.(bool)

		instanceproto.BootVolumeAttachment = &vpcv1.VolumeAttachmentPrototypeInstanceByImageContext{
			DeleteVolumeOnInstanceDelete: &deletebool,
			Volume:                       volTemplate,
		}
	}

	if primnicintf, ok := d.GetOk(isInstancePrimaryNetworkInterface); ok {
		primnic := primnicintf.([]interface{})[0].(map[string]interface{})
		subnetintf, _ := primnic[isInstanceNicSubnet]
		subnetintfstr := subnetintf.(string)
		var primnicobj = &vpcv1.NetworkInterfacePrototype{}
		primnicobj.Subnet = &vpcv1.SubnetIdentity{
			ID: &subnetintfstr,
		}
		name, _ := primnic[isInstanceNicName]
		namestr := name.(string)
		if namestr != "" {
			primnicobj.Name = &namestr
		}

		// reserved ip changes

		var ipv4str, reservedIp, reservedipv4, reservedipname string
		var autodelete, okAuto bool
		ipv4, _ := primnic[isInstanceNicPrimaryIpv4Address]
		ipv4str = ipv4.(string)

		primaryIpOk, ok := primnic[isInstanceNicPrimaryIP]
		if ok && len(primaryIpOk.([]interface{})) > 0 {
			primip := primaryIpOk.([]interface{})[0].(map[string]interface{})

			reservedipok, _ := primip[isInstanceNicReservedIpId]
			reservedIp = reservedipok.(string)

			reservedipv4Ok, _ := primip[isInstanceNicReservedIpAddress]
			reservedipv4 = reservedipv4Ok.(string)

			reservedipnameOk, _ := primip[isInstanceNicReservedIpName]
			reservedipname = reservedipnameOk.(string)
			var reservedipautodeleteok interface{}
			reservedipautodeleteok, okAuto = primip[isInstanceNicReservedIpAutoDelete]
			autodelete = reservedipautodeleteok.(bool)
		}
		if ipv4str != "" && reservedipv4 != "" && ipv4str != reservedipv4 {
			return fmt.Errorf("[ERROR] Error creating instance, primary_network_interface error, use either primary_ipv4_address(%s) or primary_ip.0.address(%s)", ipv4str, reservedipv4)
		}
		if reservedIp != "" && (ipv4str != "" || reservedipv4 != "" || reservedipname != "") {
			return fmt.Errorf("[ERROR] Error creating instance, primary_network_interface error, reserved_ip(%s) is mutually exclusive with other primary_ip attributes", reservedIp)
		}
		if reservedIp != "" {
			primnicobj.PrimaryIP = &vpcv1.NetworkInterfaceIPPrototypeReservedIPIdentity{
				ID: &reservedIp,
			}
		} else {
			if ipv4str != "" || reservedipv4 != "" || reservedipname != "" || okAuto {
				primaryipobj := &vpcv1.NetworkInterfaceIPPrototypeReservedIPPrototypeNetworkInterfaceContext{}
				if ipv4str != "" {
					primaryipobj.Address = &ipv4str
				}
				if reservedipv4 != "" {
					primaryipobj.Address = &reservedipv4
				}
				if reservedipname != "" {
					primaryipobj.Name = &reservedipname
				}
				if okAuto {
					primaryipobj.AutoDelete = &autodelete
				}
				primnicobj.PrimaryIP = primaryipobj
			}
		}
		allowIPSpoofing, ok := primnic[isInstanceNicAllowIPSpoofing]
		allowIPSpoofingbool := allowIPSpoofing.(bool)
		if ok {
			primnicobj.AllowIPSpoofing = &allowIPSpoofingbool
		}
		secgrpintf, ok := primnic[isInstanceNicSecurityGroups]
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
		instanceproto.PrimaryNetworkInterface = primnicobj
	}

	if nicsintf, ok := d.GetOk(isInstanceNetworkInterfaces); ok {
		nics := nicsintf.([]interface{})
		var intfs []vpcv1.NetworkInterfacePrototype
		for _, resource := range nics {
			nic := resource.(map[string]interface{})
			nwInterface := &vpcv1.NetworkInterfacePrototype{}
			subnetintf, _ := nic[isInstanceNicSubnet]
			subnetintfstr := subnetintf.(string)
			nwInterface.Subnet = &vpcv1.SubnetIdentity{
				ID: &subnetintfstr,
			}
			name, ok := nic[isInstanceNicName]
			namestr := name.(string)
			if ok && namestr != "" {
				nwInterface.Name = &namestr
			}

			// reserved ip changes

			var ipv4str, reservedIp, reservedipv4, reservedipname string
			var autodelete, okAuto bool
			ipv4, _ := nic[isInstanceNicPrimaryIpv4Address]
			ipv4str = ipv4.(string)

			primaryIpOk, ok := nic[isInstanceNicPrimaryIP]
			if ok && len(primaryIpOk.([]interface{})) > 0 {
				primip := primaryIpOk.([]interface{})[0].(map[string]interface{})

				reservedipok, _ := primip[isInstanceNicReservedIpId]
				reservedIp = reservedipok.(string)

				reservedipv4Ok, _ := primip[isInstanceNicReservedIpAddress]
				reservedipv4 = reservedipv4Ok.(string)

				reservedipnameOk, _ := primip[isInstanceNicReservedIpName]
				reservedipname = reservedipnameOk.(string)
				var reservedipautodeleteok interface{}
				reservedipautodeleteok, okAuto = primip[isInstanceNicReservedIpAutoDelete]
				autodelete = reservedipautodeleteok.(bool)
			}
			if ipv4str != "" && reservedipv4 != "" && ipv4str != reservedipv4 {
				return fmt.Errorf("[ERROR] Error creating instance, network_interfaces error, use either primary_ipv4_address(%s) or primary_ip.0.address(%s)", ipv4str, reservedipv4)
			}
			if reservedIp != "" && (ipv4str != "" || reservedipv4 != "" || reservedipname != "") {
				return fmt.Errorf("[ERROR] Error creating instance, network_interfaces error, reserved_ip(%s) is mutually exclusive with other primary_ip attributes", reservedIp)
			}
			if reservedIp != "" {
				nwInterface.PrimaryIP = &vpcv1.NetworkInterfaceIPPrototypeReservedIPIdentity{
					ID: &reservedIp,
				}
			} else {
				if ipv4str != "" || reservedipv4 != "" || reservedipname != "" || okAuto {
					primaryipobj := &vpcv1.NetworkInterfaceIPPrototypeReservedIPPrototypeNetworkInterfaceContext{}
					if ipv4str != "" {
						primaryipobj.Address = &ipv4str
					}
					if reservedipv4 != "" {
						primaryipobj.Address = &reservedipv4
					}
					if reservedipname != "" {
						primaryipobj.Name = &reservedipname
					}
					if okAuto {
						primaryipobj.AutoDelete = &autodelete
					}
					nwInterface.PrimaryIP = primaryipobj
				}
			}
			allowIPSpoofing, ok := nic[isInstanceNicAllowIPSpoofing]
			allowIPSpoofingbool := allowIPSpoofing.(bool)
			if ok {
				nwInterface.AllowIPSpoofing = &allowIPSpoofingbool
			}
			secgrpintf, ok := nic[isInstanceNicSecurityGroups]
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
					nwInterface.SecurityGroups = secgrpobjs
				}
			}
			intfs = append(intfs, *nwInterface)
		}
		instanceproto.NetworkInterfaces = intfs
	}

	keySet := d.Get(isInstanceKeys).(*schema.Set)
	if keySet.Len() != 0 {
		keyobjs := make([]vpcv1.KeyIdentityIntf, keySet.Len())
		for i, key := range keySet.List() {
			keystr := key.(string)
			keyobjs[i] = &vpcv1.KeyIdentity{
				ID: &keystr,
			}
		}
		instanceproto.Keys = keyobjs
	}

	if userdata, ok := d.GetOk(isInstanceUserData); ok {
		userdatastr := userdata.(string)
		instanceproto.UserData = &userdatastr
	}

	if grp, ok := d.GetOk(isInstanceResourceGroup); ok {
		grpstr := grp.(string)
		instanceproto.ResourceGroup = &vpcv1.ResourceGroupIdentity{
			ID: &grpstr,
		}

	}

	if metadataServiceEnabled, ok := d.GetOkExists(isInstanceMetadataServiceEnabled); ok {
		metadataServiceEnabledBool := metadataServiceEnabled.(bool)
		instanceproto.MetadataService = &vpcv1.InstanceMetadataServicePrototype{
			Enabled: &metadataServiceEnabledBool,
		}
	}

	if metadataService := GetInstanceMetadataServiceOptions(d); metadataService != nil {
		instanceproto.MetadataService = metadataService
	}

	options := &vpcv1.CreateInstanceOptions{
		InstancePrototype: instanceproto,
	}

	instance, response, err := sess.CreateInstance(options)
	if err != nil {
		log.Printf("[DEBUG] Instance err %s\n%s", err, response)
		return err
	}
	d.SetId(*instance.ID)

	log.Printf("[INFO] Instance : %s", *instance.ID)
	d.Set(isInstanceStatus, instance.Status)

	_, err = isWaitForInstanceAvailable(sess, d.Id(), d.Timeout(schema.TimeoutCreate), d)
	if err != nil {
		return err
	}

	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk(isInstanceTags); ok || v != "" {
		oldList, newList := d.GetChange(isInstanceTags)
		err = flex.UpdateTagsUsingCRN(oldList, newList, meta, *instance.CRN)
		if err != nil {
			log.Printf(
				"Error on create of resource instance (%s) tags: %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk(isInstanceAccessTags); ok {
		oldList, newList := d.GetChange(isInstanceAccessTags)
		err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *instance.CRN, "", isInstanceAccessTagType)
		if err != nil {
			log.Printf(
				"Error on create of resource instance (%s) access tags: %s", d.Id(), err)
		}
	}
	return nil
}

func instanceCreateBySnapshot(d *schema.ResourceData, meta interface{}, profile, name, vpcID, zone string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	instanceproto := &vpcv1.InstancePrototypeInstanceBySourceSnapshot{
		Zone: &vpcv1.ZoneIdentity{
			Name: &zone,
		},
		Profile: &vpcv1.InstanceProfileIdentity{
			Name: &profile,
		},
		Name: &name,
		VPC: &vpcv1.VPCIdentity{
			ID: &vpcID,
		},
	}

	if defaultTrustedProfileTargetIntf, ok := d.GetOk(isInstanceDefaultTrustedProfileTarget); ok {
		defaultTrustedProfiletarget := defaultTrustedProfileTargetIntf.(string)

		target := &vpcv1.TrustedProfileIdentity{}
		if strings.HasPrefix(defaultTrustedProfiletarget, "crn") {
			target.CRN = &defaultTrustedProfiletarget
		} else {
			target.ID = &defaultTrustedProfiletarget
		}
		instanceproto.DefaultTrustedProfile = &vpcv1.InstanceDefaultTrustedProfilePrototype{
			Target: target,
		}

		if defaultTrustedProfileAutoLinkIntf, ok := d.GetOkExists(isInstanceDefaultTrustedProfileAutoLink); ok {
			defaultTrustedProfileAutoLink := defaultTrustedProfileAutoLinkIntf.(bool)
			instanceproto.DefaultTrustedProfile.AutoLink = &defaultTrustedProfileAutoLink
		}
	}

	if dHostIdInf, ok := d.GetOk(isPlacementTargetDedicatedHost); ok {
		dHostIdStr := dHostIdInf.(string)
		dHostPlaementTarget := &vpcv1.InstancePlacementTargetPrototypeDedicatedHostIdentity{
			ID: &dHostIdStr,
		}
		instanceproto.PlacementTarget = dHostPlaementTarget
	} else if dHostGrpIdInf, ok := d.GetOk(isPlacementTargetDedicatedHostGroup); ok {
		dHostGrpIdStr := dHostGrpIdInf.(string)
		dHostGrpPlaementTarget := &vpcv1.InstancePlacementTargetPrototypeDedicatedHostGroupIdentity{
			ID: &dHostGrpIdStr,
		}
		instanceproto.PlacementTarget = dHostGrpPlaementTarget
	} else if placementGroupInf, ok := d.GetOk(isPlacementTargetPlacementGroup); ok {
		placementGrpStr := placementGroupInf.(string)
		placementGrp := &vpcv1.InstancePlacementTargetPrototypePlacementGroupIdentity{
			ID: &placementGrpStr,
		}
		instanceproto.PlacementTarget = placementGrp
	}

	if boot, ok := d.GetOk(isInstanceBootVolume); ok {
		bootvol := boot.([]interface{})[0].(map[string]interface{})
		var volTemplate = &vpcv1.VolumePrototypeInstanceBySourceSnapshotContext{}

		name, ok := bootvol[isInstanceBootAttachmentName]
		namestr := name.(string)
		if namestr != "" && ok {
			volTemplate.Name = &namestr
		}
		sizeOk, ok := bootvol[isInstanceBootSize]
		size := sizeOk.(int)
		if size != 0 && ok {
			sizeInt64 := int64(size)
			volTemplate.Capacity = &sizeInt64
		}
		enc, ok := bootvol[isInstanceBootEncryption]
		encstr := enc.(string)
		if ok && encstr != "" {
			volTemplate.EncryptionKey = &vpcv1.EncryptionKeyIdentity{
				CRN: &encstr,
			}
		}
		volprof := "general-purpose"
		volTemplate.Profile = &vpcv1.VolumeProfileIdentity{
			Name: &volprof,
		}
		var userTags *schema.Set
		if v, ok := bootvol[isInstanceBootVolumeTags]; ok {
			userTags = v.(*schema.Set)
			if userTags != nil && userTags.Len() != 0 {
				userTagsArray := make([]string, userTags.Len())
				for i, userTag := range userTags.List() {
					userTagStr := userTag.(string)
					userTagsArray[i] = userTagStr
				}
				volTemplate.UserTags = userTagsArray
			}
		}
		snapshotId, ok := bootvol[isInstanceVolumeSnapshot]
		snapshotIdStr := snapshotId.(string)
		if snapshotIdStr != "" && ok {
			volTemplate.SourceSnapshot = &vpcv1.SnapshotIdentity{
				ID: &snapshotIdStr,
			}
		}
		deleteboolIntf := bootvol[isInstanceVolAttVolAutoDelete]
		deletebool := deleteboolIntf.(bool)
		instanceproto.BootVolumeAttachment = &vpcv1.VolumeAttachmentPrototypeInstanceBySourceSnapshotContext{
			DeleteVolumeOnInstanceDelete: &deletebool,
			Volume:                       volTemplate,
		}
	}
	if totalVolBandwidthIntf, ok := d.GetOk(isInstanceTotalVolumeBandwidth); ok {
		totalVolBandwidthStr := int64(totalVolBandwidthIntf.(int))
		instanceproto.TotalVolumeBandwidth = &totalVolBandwidthStr
	}

	if primnicintf, ok := d.GetOk(isInstancePrimaryNetworkInterface); ok {
		primnic := primnicintf.([]interface{})[0].(map[string]interface{})
		subnetintf, _ := primnic[isInstanceNicSubnet]
		subnetintfstr := subnetintf.(string)
		var primnicobj = &vpcv1.NetworkInterfacePrototype{}
		primnicobj.Subnet = &vpcv1.SubnetIdentity{
			ID: &subnetintfstr,
		}
		name, _ := primnic[isInstanceNicName]
		namestr := name.(string)
		if namestr != "" {
			primnicobj.Name = &namestr
		}

		// reserved ip changes

		var ipv4str, reservedIp, reservedipv4, reservedipname string
		var autodelete, okAuto bool
		ipv4, _ := primnic[isInstanceNicPrimaryIpv4Address]
		ipv4str = ipv4.(string)

		primaryIpOk, ok := primnic[isInstanceNicPrimaryIP]
		if ok && len(primaryIpOk.([]interface{})) > 0 {
			primip := primaryIpOk.([]interface{})[0].(map[string]interface{})

			reservedipok, _ := primip[isInstanceNicReservedIpId]
			reservedIp = reservedipok.(string)

			reservedipv4Ok, _ := primip[isInstanceNicReservedIpAddress]
			reservedipv4 = reservedipv4Ok.(string)

			reservedipnameOk, _ := primip[isInstanceNicReservedIpName]
			reservedipname = reservedipnameOk.(string)
			var reservedipautodeleteok interface{}
			reservedipautodeleteok, okAuto = primip[isInstanceNicReservedIpAutoDelete]
			autodelete = reservedipautodeleteok.(bool)
		}
		if ipv4str != "" && reservedipv4 != "" && ipv4str != reservedipv4 {
			return fmt.Errorf("[ERROR] Error creating instance, primary_network_interface error, use either primary_ipv4_address(%s) or primary_ip.0.address(%s)", ipv4str, reservedipv4)
		}
		if reservedIp != "" && (ipv4str != "" || reservedipv4 != "" || reservedipname != "") {
			return fmt.Errorf("[ERROR] Error creating instance, primary_network_interface error, reserved_ip(%s) is mutually exclusive with other primary_ip attributes", reservedIp)
		}
		if reservedIp != "" {
			primnicobj.PrimaryIP = &vpcv1.NetworkInterfaceIPPrototypeReservedIPIdentity{
				ID: &reservedIp,
			}
		} else {
			if ipv4str != "" || reservedipv4 != "" || reservedipname != "" || okAuto {
				primaryipobj := &vpcv1.NetworkInterfaceIPPrototypeReservedIPPrototypeNetworkInterfaceContext{}
				if ipv4str != "" {
					primaryipobj.Address = &ipv4str
				}
				if reservedipv4 != "" {
					primaryipobj.Address = &reservedipv4
				}
				if reservedipname != "" {
					primaryipobj.Name = &reservedipname
				}
				if okAuto {
					primaryipobj.AutoDelete = &autodelete
				}
				primnicobj.PrimaryIP = primaryipobj
			}
		}

		allowIPSpoofing, ok := primnic[isInstanceNicAllowIPSpoofing]
		allowIPSpoofingbool := allowIPSpoofing.(bool)
		if ok {
			primnicobj.AllowIPSpoofing = &allowIPSpoofingbool
		}
		secgrpintf, ok := primnic[isInstanceNicSecurityGroups]
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
		instanceproto.PrimaryNetworkInterface = primnicobj
	}

	if nicsintf, ok := d.GetOk(isInstanceNetworkInterfaces); ok {
		nics := nicsintf.([]interface{})
		var intfs []vpcv1.NetworkInterfacePrototype
		for _, resource := range nics {
			nic := resource.(map[string]interface{})
			nwInterface := &vpcv1.NetworkInterfacePrototype{}
			subnetintf, _ := nic[isInstanceNicSubnet]
			subnetintfstr := subnetintf.(string)
			nwInterface.Subnet = &vpcv1.SubnetIdentity{
				ID: &subnetintfstr,
			}
			name, ok := nic[isInstanceNicName]
			namestr := name.(string)
			if ok && namestr != "" {
				nwInterface.Name = &namestr
			}
			// reserved ip changes

			var ipv4str, reservedIp, reservedipv4, reservedipname string
			var autodelete, okAuto bool
			ipv4, _ := nic[isInstanceNicPrimaryIpv4Address]
			ipv4str = ipv4.(string)

			primaryIpOk, ok := nic[isInstanceNicPrimaryIP]
			if ok && len(primaryIpOk.([]interface{})) > 0 {
				primip := primaryIpOk.([]interface{})[0].(map[string]interface{})

				reservedipok, _ := primip[isInstanceNicReservedIpId]
				reservedIp = reservedipok.(string)

				reservedipv4Ok, _ := primip[isInstanceNicReservedIpAddress]
				reservedipv4 = reservedipv4Ok.(string)

				reservedipnameOk, _ := primip[isInstanceNicReservedIpName]
				reservedipname = reservedipnameOk.(string)
				var reservedipautodeleteok interface{}
				reservedipautodeleteok, okAuto = primip[isInstanceNicReservedIpAutoDelete]
				autodelete = reservedipautodeleteok.(bool)
			}
			if ipv4str != "" && reservedipv4 != "" && ipv4str != reservedipv4 {
				return fmt.Errorf("[ERROR] Error creating instance, network_interfaces error, use either primary_ipv4_address(%s) or primary_ip.0.address(%s)", ipv4str, reservedipv4)
			}
			if reservedIp != "" && (ipv4str != "" || reservedipv4 != "" || reservedipname != "") {
				return fmt.Errorf("[ERROR] Error creating instance, network_interfaces error, reserved_ip(%s) is mutually exclusive with other primary_ip attributes", reservedIp)
			}
			if reservedIp != "" {
				nwInterface.PrimaryIP = &vpcv1.NetworkInterfaceIPPrototypeReservedIPIdentity{
					ID: &reservedIp,
				}
			} else {
				if ipv4str != "" || reservedipv4 != "" || reservedipname != "" || okAuto {
					primaryipobj := &vpcv1.NetworkInterfaceIPPrototypeReservedIPPrototypeNetworkInterfaceContext{}
					if ipv4str != "" {
						primaryipobj.Address = &ipv4str
					}
					if reservedipv4 != "" {
						primaryipobj.Address = &reservedipv4
					}
					if reservedipname != "" {
						primaryipobj.Name = &reservedipname
					}
					if okAuto {
						primaryipobj.AutoDelete = &autodelete
					}
					nwInterface.PrimaryIP = primaryipobj
				}
			}
			allowIPSpoofing, ok := nic[isInstanceNicAllowIPSpoofing]
			allowIPSpoofingbool := allowIPSpoofing.(bool)
			if ok {
				nwInterface.AllowIPSpoofing = &allowIPSpoofingbool
			}
			secgrpintf, ok := nic[isInstanceNicSecurityGroups]
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
					nwInterface.SecurityGroups = secgrpobjs
				}
			}
			intfs = append(intfs, *nwInterface)
		}
		instanceproto.NetworkInterfaces = intfs
	}

	keySet := d.Get(isInstanceKeys).(*schema.Set)
	if keySet.Len() != 0 {
		keyobjs := make([]vpcv1.KeyIdentityIntf, keySet.Len())
		for i, key := range keySet.List() {
			keystr := key.(string)
			keyobjs[i] = &vpcv1.KeyIdentity{
				ID: &keystr,
			}
		}
		instanceproto.Keys = keyobjs
	}

	if userdata, ok := d.GetOk(isInstanceUserData); ok {
		userdatastr := userdata.(string)
		instanceproto.UserData = &userdatastr
	}

	if grp, ok := d.GetOk(isInstanceResourceGroup); ok {
		grpstr := grp.(string)
		instanceproto.ResourceGroup = &vpcv1.ResourceGroupIdentity{
			ID: &grpstr,
		}

	}
	if availablePolicyItem, ok := d.GetOk(isInstanceAvailablePolicyHostFailure); ok {
		hostFailure := availablePolicyItem.(string)
		instanceproto.AvailabilityPolicy = &vpcv1.InstanceAvailabilityPolicyPrototype{
			HostFailure: &hostFailure,
		}
	}
	metadataServiceEnabled := d.Get(isInstanceMetadataServiceEnabled).(bool)
	if metadataServiceEnabled {
		instanceproto.MetadataService = &vpcv1.InstanceMetadataServicePrototype{
			Enabled: &metadataServiceEnabled,
		}
	}

	if metadataService := GetInstanceMetadataServiceOptions(d); metadataService != nil {
		instanceproto.MetadataService = metadataService
	}

	options := &vpcv1.CreateInstanceOptions{
		InstancePrototype: instanceproto,
	}

	instance, response, err := sess.CreateInstance(options)
	if err != nil {
		log.Printf("[DEBUG] Instance err %s\n%s", err, response)
		return err
	}
	d.SetId(*instance.ID)

	log.Printf("[INFO] Instance : %s", *instance.ID)
	d.Set(isInstanceStatus, instance.Status)

	_, err = isWaitForInstanceAvailable(sess, d.Id(), d.Timeout(schema.TimeoutCreate), d)
	if err != nil {
		return err
	}

	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk(isInstanceTags); ok || v != "" {
		oldList, newList := d.GetChange(isInstanceTags)
		err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *instance.CRN, "", isInstanceUserTagType)
		if err != nil {
			log.Printf(
				"Error on create of resource instance (%s) tags: %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk(isInstanceAccessTags); ok {
		oldList, newList := d.GetChange(isInstanceAccessTags)
		err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *instance.CRN, "", isInstanceAccessTagType)
		if err != nil {
			log.Printf(
				"Error on create of resource instance (%s) access tags: %s", d.Id(), err)
		}
	}
	return nil
}

func instanceCreateByVolume(d *schema.ResourceData, meta interface{}, profile, name, vpcID, zone string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	instanceproto := &vpcv1.InstancePrototypeInstanceByVolume{
		Zone: &vpcv1.ZoneIdentity{
			Name: &zone,
		},
		Profile: &vpcv1.InstanceProfileIdentity{
			Name: &profile,
		},
		Name: &name,
		VPC: &vpcv1.VPCIdentity{
			ID: &vpcID,
		},
	}

	if defaultTrustedProfileTargetIntf, ok := d.GetOk(isInstanceDefaultTrustedProfileTarget); ok {
		defaultTrustedProfiletarget := defaultTrustedProfileTargetIntf.(string)

		target := &vpcv1.TrustedProfileIdentity{}
		if strings.HasPrefix(defaultTrustedProfiletarget, "crn") {
			target.CRN = &defaultTrustedProfiletarget
		} else {
			target.ID = &defaultTrustedProfiletarget
		}
		instanceproto.DefaultTrustedProfile = &vpcv1.InstanceDefaultTrustedProfilePrototype{
			Target: target,
		}

		if defaultTrustedProfileAutoLinkIntf, ok := d.GetOkExists(isInstanceDefaultTrustedProfileAutoLink); ok {
			defaultTrustedProfileAutoLink := defaultTrustedProfileAutoLinkIntf.(bool)
			instanceproto.DefaultTrustedProfile.AutoLink = &defaultTrustedProfileAutoLink
		}
	}

	if dHostIdInf, ok := d.GetOk(isPlacementTargetDedicatedHost); ok {
		dHostIdStr := dHostIdInf.(string)
		dHostPlaementTarget := &vpcv1.InstancePlacementTargetPrototypeDedicatedHostIdentity{
			ID: &dHostIdStr,
		}
		instanceproto.PlacementTarget = dHostPlaementTarget
	} else if dHostGrpIdInf, ok := d.GetOk(isPlacementTargetDedicatedHostGroup); ok {
		dHostGrpIdStr := dHostGrpIdInf.(string)
		dHostGrpPlaementTarget := &vpcv1.InstancePlacementTargetPrototypeDedicatedHostGroupIdentity{
			ID: &dHostGrpIdStr,
		}
		instanceproto.PlacementTarget = dHostGrpPlaementTarget
	} else if placementGroupInf, ok := d.GetOk(isPlacementTargetPlacementGroup); ok {
		placementGrpStr := placementGroupInf.(string)
		placementGrp := &vpcv1.InstancePlacementTargetPrototypePlacementGroupIdentity{
			ID: &placementGrpStr,
		}
		instanceproto.PlacementTarget = placementGrp
	}

	if boot, ok := d.GetOk(isInstanceBootVolume); ok {
		bootvol := boot.([]interface{})[0].(map[string]interface{})
		volumeId, ok := bootvol[isInstanceBootVolumeId]

		volumeIdStr := volumeId.(string)
		bootVolAttachment := &vpcv1.VolumeAttachmentPrototypeInstanceByVolumeContext{}
		if ok && volumeIdStr != "" {
			volumeIdentity := &vpcv1.VolumeIdentity{
				ID: &volumeIdStr,
			}
			bootVolAttachment.Volume = volumeIdentity
		}
		if autoDeleteIntf, ok := d.GetOk("boot_volume.0.auto_delete_volume"); ok {
			autoDelete := autoDeleteIntf.(bool)
			bootVolAttachment.DeleteVolumeOnInstanceDelete = &autoDelete
		}
		instanceproto.BootVolumeAttachment = bootVolAttachment
	}
	if totalVolBandwidthIntf, ok := d.GetOk(isInstanceTotalVolumeBandwidth); ok {
		totalVolBandwidthStr := int64(totalVolBandwidthIntf.(int))
		instanceproto.TotalVolumeBandwidth = &totalVolBandwidthStr
	}

	if primnicintf, ok := d.GetOk(isInstancePrimaryNetworkInterface); ok {
		primnic := primnicintf.([]interface{})[0].(map[string]interface{})
		subnetintf, _ := primnic[isInstanceNicSubnet]
		subnetintfstr := subnetintf.(string)
		var primnicobj = &vpcv1.NetworkInterfacePrototype{}
		primnicobj.Subnet = &vpcv1.SubnetIdentity{
			ID: &subnetintfstr,
		}
		name, _ := primnic[isInstanceNicName]
		namestr := name.(string)
		if namestr != "" {
			primnicobj.Name = &namestr
		}

		// reserved ip changes

		var ipv4str, reservedIp, reservedipv4, reservedipname string
		var autodelete, okAuto bool
		ipv4, _ := primnic[isInstanceNicPrimaryIpv4Address]
		ipv4str = ipv4.(string)

		primaryIpOk, ok := primnic[isInstanceNicPrimaryIP]
		if ok && len(primaryIpOk.([]interface{})) > 0 {
			primip := primaryIpOk.([]interface{})[0].(map[string]interface{})

			reservedipok, _ := primip[isInstanceNicReservedIpId]
			reservedIp = reservedipok.(string)

			reservedipv4Ok, _ := primip[isInstanceNicReservedIpAddress]
			reservedipv4 = reservedipv4Ok.(string)

			reservedipnameOk, _ := primip[isInstanceNicReservedIpName]
			reservedipname = reservedipnameOk.(string)
			var reservedipautodeleteok interface{}
			reservedipautodeleteok, okAuto = primip[isInstanceNicReservedIpAutoDelete]
			autodelete = reservedipautodeleteok.(bool)
		}
		if ipv4str != "" && reservedipv4 != "" && ipv4str != reservedipv4 {
			return fmt.Errorf("[ERROR] Error creating instance, primary_network_interface error, use either primary_ipv4_address(%s) or primary_ip.0.address(%s)", ipv4str, reservedipv4)
		}
		if reservedIp != "" && (ipv4str != "" || reservedipv4 != "" || reservedipname != "") {
			return fmt.Errorf("[ERROR] Error creating instance, primary_network_interface error, reserved_ip(%s) is mutually exclusive with other primary_ip attributes", reservedIp)
		}
		if reservedIp != "" {
			primnicobj.PrimaryIP = &vpcv1.NetworkInterfaceIPPrototypeReservedIPIdentity{
				ID: &reservedIp,
			}
		} else {
			if ipv4str != "" || reservedipv4 != "" || reservedipname != "" || okAuto {
				primaryipobj := &vpcv1.NetworkInterfaceIPPrototypeReservedIPPrototypeNetworkInterfaceContext{}
				if ipv4str != "" {
					primaryipobj.Address = &ipv4str
				}
				if reservedipv4 != "" {
					primaryipobj.Address = &reservedipv4
				}
				if reservedipname != "" {
					primaryipobj.Name = &reservedipname
				}
				if okAuto {
					primaryipobj.AutoDelete = &autodelete
				}
				primnicobj.PrimaryIP = primaryipobj
			}
		}

		allowIPSpoofing, ok := primnic[isInstanceNicAllowIPSpoofing]
		allowIPSpoofingbool := allowIPSpoofing.(bool)
		if ok {
			primnicobj.AllowIPSpoofing = &allowIPSpoofingbool
		}
		secgrpintf, ok := primnic[isInstanceNicSecurityGroups]
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
		instanceproto.PrimaryNetworkInterface = primnicobj
	}

	if nicsintf, ok := d.GetOk(isInstanceNetworkInterfaces); ok {
		nics := nicsintf.([]interface{})
		var intfs []vpcv1.NetworkInterfacePrototype
		for _, resource := range nics {
			nic := resource.(map[string]interface{})
			nwInterface := &vpcv1.NetworkInterfacePrototype{}
			subnetintf, _ := nic[isInstanceNicSubnet]
			subnetintfstr := subnetintf.(string)
			nwInterface.Subnet = &vpcv1.SubnetIdentity{
				ID: &subnetintfstr,
			}
			name, ok := nic[isInstanceNicName]
			namestr := name.(string)
			if ok && namestr != "" {
				nwInterface.Name = &namestr
			}
			// reserved ip changes

			var ipv4str, reservedIp, reservedipv4, reservedipname string
			var autodelete, okAuto bool
			ipv4, _ := nic[isInstanceNicPrimaryIpv4Address]
			ipv4str = ipv4.(string)

			primaryIpOk, ok := nic[isInstanceNicPrimaryIP]
			if ok && len(primaryIpOk.([]interface{})) > 0 {
				primip := primaryIpOk.([]interface{})[0].(map[string]interface{})

				reservedipok, _ := primip[isInstanceNicReservedIpId]
				reservedIp = reservedipok.(string)

				reservedipv4Ok, _ := primip[isInstanceNicReservedIpAddress]
				reservedipv4 = reservedipv4Ok.(string)

				reservedipnameOk, _ := primip[isInstanceNicReservedIpName]
				reservedipname = reservedipnameOk.(string)
				var reservedipautodeleteok interface{}
				reservedipautodeleteok, okAuto = primip[isInstanceNicReservedIpAutoDelete]
				autodelete = reservedipautodeleteok.(bool)
			}
			if ipv4str != "" && reservedipv4 != "" && ipv4str != reservedipv4 {
				return fmt.Errorf("[ERROR] Error creating instance, network_interfaces error, use either primary_ipv4_address(%s) or primary_ip.0.address(%s)", ipv4str, reservedipv4)
			}
			if reservedIp != "" && (ipv4str != "" || reservedipv4 != "" || reservedipname != "") {
				return fmt.Errorf("[ERROR] Error creating instance, network_interfaces error, reserved_ip(%s) is mutually exclusive with other primary_ip attributes", reservedIp)
			}
			if reservedIp != "" {
				nwInterface.PrimaryIP = &vpcv1.NetworkInterfaceIPPrototypeReservedIPIdentity{
					ID: &reservedIp,
				}
			} else {
				if ipv4str != "" || reservedipv4 != "" || reservedipname != "" || okAuto {
					primaryipobj := &vpcv1.NetworkInterfaceIPPrototypeReservedIPPrototypeNetworkInterfaceContext{}
					if ipv4str != "" {
						primaryipobj.Address = &ipv4str
					}
					if reservedipv4 != "" {
						primaryipobj.Address = &reservedipv4
					}
					if reservedipname != "" {
						primaryipobj.Name = &reservedipname
					}
					if okAuto {
						primaryipobj.AutoDelete = &autodelete
					}
					nwInterface.PrimaryIP = primaryipobj
				}
			}
			allowIPSpoofing, ok := nic[isInstanceNicAllowIPSpoofing]
			allowIPSpoofingbool := allowIPSpoofing.(bool)
			if ok {
				nwInterface.AllowIPSpoofing = &allowIPSpoofingbool
			}
			secgrpintf, ok := nic[isInstanceNicSecurityGroups]
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
					nwInterface.SecurityGroups = secgrpobjs
				}
			}
			intfs = append(intfs, *nwInterface)
		}
		instanceproto.NetworkInterfaces = intfs
	}

	keySet := d.Get(isInstanceKeys).(*schema.Set)
	if keySet.Len() != 0 {
		keyobjs := make([]vpcv1.KeyIdentityIntf, keySet.Len())
		for i, key := range keySet.List() {
			keystr := key.(string)
			keyobjs[i] = &vpcv1.KeyIdentity{
				ID: &keystr,
			}
		}
		instanceproto.Keys = keyobjs
	}

	if userdata, ok := d.GetOk(isInstanceUserData); ok {
		userdatastr := userdata.(string)
		instanceproto.UserData = &userdatastr
	}

	if grp, ok := d.GetOk(isInstanceResourceGroup); ok {
		grpstr := grp.(string)
		instanceproto.ResourceGroup = &vpcv1.ResourceGroupIdentity{
			ID: &grpstr,
		}

	}
	if availablePolicyItem, ok := d.GetOk(isInstanceAvailablePolicyHostFailure); ok {
		hostFailure := availablePolicyItem.(string)
		instanceproto.AvailabilityPolicy = &vpcv1.InstanceAvailabilityPolicyPrototype{
			HostFailure: &hostFailure,
		}
	}
	metadataServiceEnabled := d.Get(isInstanceMetadataServiceEnabled).(bool)
	if metadataServiceEnabled {
		instanceproto.MetadataService = &vpcv1.InstanceMetadataServicePrototype{
			Enabled: &metadataServiceEnabled,
		}
	}
	if metadataService := GetInstanceMetadataServiceOptions(d); metadataService != nil {
		instanceproto.MetadataService = metadataService
	}

	options := &vpcv1.CreateInstanceOptions{
		InstancePrototype: instanceproto,
	}

	instance, response, err := sess.CreateInstance(options)
	if err != nil {
		log.Printf("[DEBUG] Instance err %s\n%s", err, response)
		return err
	}
	d.SetId(*instance.ID)

	log.Printf("[INFO] Instance : %s", *instance.ID)
	d.Set(isInstanceStatus, instance.Status)

	_, err = isWaitForInstanceAvailable(sess, d.Id(), d.Timeout(schema.TimeoutCreate), d)
	if err != nil {
		return err
	}

	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk(isInstanceTags); ok || v != "" {
		oldList, newList := d.GetChange(isInstanceTags)
		err = flex.UpdateTagsUsingCRN(oldList, newList, meta, *instance.CRN)
		if err != nil {
			log.Printf(
				"Error on create of resource instance (%s) tags: %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk(isInstanceAccessTags); ok {
		oldList, newList := d.GetChange(isInstanceAccessTags)
		err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *instance.CRN, "", isInstanceAccessTagType)
		if err != nil {
			log.Printf(
				"Error on create of resource instance (%s) access tags: %s", d.Id(), err)
		}
	}
	return nil
}

func resourceIBMisInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	profile := d.Get(isInstanceProfile).(string)
	name := d.Get(isInstanceName).(string)
	vpcID := d.Get(isInstanceVPC).(string)
	zone := d.Get(isInstanceZone).(string)
	image := d.Get(isInstanceImage).(string)
	snapshot := d.Get("boot_volume.0.snapshot").(string)
	volume := d.Get("boot_volume.0.volume_id").(string)
	template := d.Get(isInstanceSourceTemplate).(string)
	if catalogOfferingOk, ok := d.GetOk(isInstanceCatalogOffering); ok {
		catalogOffering := catalogOfferingOk.([]interface{})[0].(map[string]interface{})
		offeringCrn, _ := catalogOffering[isInstanceCatalogOfferingOfferingCrn].(string)
		versionCrn, _ := catalogOffering[isInstanceCatalogOfferingVersionCrn].(string)
		err := instanceCreateByCatalogOffering(d, meta, profile, name, vpcID, zone, image, offeringCrn, versionCrn)
		if err != nil {
			return err
		}

	} else if volume != "" {
		err := instanceCreateByVolume(d, meta, profile, name, vpcID, zone)
		if err != nil {
			return err
		}
	} else if snapshot != "" {
		err := instanceCreateBySnapshot(d, meta, profile, name, vpcID, zone)
		if err != nil {
			return err
		}
	} else if template != "" {
		err := instanceCreateByTemplate(d, meta, profile, name, vpcID, zone, image, template)
		if err != nil {
			return err
		}
	} else {
		err := instanceCreateByImage(d, meta, profile, name, vpcID, zone, image)
		if err != nil {
			return err
		}
	}

	return resourceIBMisInstanceUpdate(d, meta)
}

func isWaitForInstanceAvailable(instanceC *vpcv1.VpcV1, id string, timeout time.Duration, d *schema.ResourceData) (interface{}, error) {
	log.Printf("Waiting for instance (%s) to be available.", id)

	communicator := make(chan interface{})

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isInstanceProvisioning},
		Target:     []string{isInstanceStatusRunning, "available", "failed", ""},
		Refresh:    isInstanceRefreshFunc(instanceC, id, d, communicator),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	if v, ok := d.GetOk("force_recovery_time"); ok {
		forceTimeout := v.(int)
		go isRestartStartAction(instanceC, id, d, forceTimeout, communicator)
	}

	return stateConf.WaitForState()
}

func isInstanceRefreshFunc(instanceC *vpcv1.VpcV1, id string, d *schema.ResourceData, communicator chan interface{}) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getinsOptions := &vpcv1.GetInstanceOptions{
			ID: &id,
		}
		instance, response, err := instanceC.GetInstance(getinsOptions)
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] Error Getting Instance: %s\n%s", err, response)
		}
		d.Set(isInstanceStatus, *instance.Status)

		select {
		case data := <-communicator:
			return nil, "", data.(error)
		default:
			fmt.Println("no message sent")
		}

		if *instance.Status == "available" || *instance.Status == "failed" || *instance.Status == "running" {
			// let know the isRestartStartAction() to stop
			close(communicator)
			// taint the instance if status is failed
			if *instance.Status == "failed" {
				instanceStatusReason := instance.StatusReasons

				//set the status reasons
				if instance.StatusReasons != nil {
					statusReasonsList := make([]map[string]interface{}, 0)
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
					d.Set(isInstanceStatusReasons, statusReasonsList)
				}

				out, err := json.MarshalIndent(instanceStatusReason, "", "    ")
				if err != nil {
					return instance, *instance.Status, fmt.Errorf("[ERROR] Instance (%s) went into failed state during the operation \n [WARNING] Running terraform apply again will remove the tainted instance and attempt to create the instance again replacing the previous configuration", *instance.ID)
				}
				return instance, *instance.Status, fmt.Errorf("[ERROR] Instance (%s) went into failed state during the operation \n (%+v) \n [WARNING] Running terraform apply again will remove the tainted instance and attempt to create the instance again replacing the previous configuration", *instance.ID, string(out))
			}
			return instance, *instance.Status, nil

		}
		return instance, isInstanceProvisioning, nil
	}
}

func isRestartStartAction(instanceC *vpcv1.VpcV1, id string, d *schema.ResourceData, forceTimeout int, communicator chan interface{}) {
	subticker := time.NewTicker(time.Duration(forceTimeout) * time.Minute)
	//subticker := time.NewTicker(time.Duration(forceTimeout) * time.Second)
	for {
		select {

		case <-subticker.C:
			log.Println("Instance is still in starting state, force retry by restarting the instance.")
			actiontype := "stop"
			createinsactoptions := &vpcv1.CreateInstanceActionOptions{
				InstanceID: &id,
				Type:       &actiontype,
			}
			_, response, err := instanceC.CreateInstanceAction(createinsactoptions)
			if err != nil {
				communicator <- fmt.Errorf("[ERROR] Error retrying instance action start: %s\n%s", err, response)
				return
			}
			waitTimeout := time.Duration(1) * time.Minute
			_, _ = isWaitForInstanceActionStop(instanceC, waitTimeout, id, d)
			actiontype = "start"
			createinsactoptions = &vpcv1.CreateInstanceActionOptions{
				InstanceID: &id,
				Type:       &actiontype,
			}
			_, response, err = instanceC.CreateInstanceAction(createinsactoptions)
			if err != nil {
				communicator <- fmt.Errorf("[ERROR] Error retrying instance action start: %s\n%s", err, response)
				return
			}
		case <-communicator:
			// indicates refresh func is reached target and not proceed with the thread
			subticker.Stop()
			return

		}
	}
}
func resourceIBMisInstanceRead(d *schema.ResourceData, meta interface{}) error {

	ID := d.Id()

	err := instanceGet(d, meta, ID)
	if err != nil {
		return err
	}
	return nil
}

func instanceGet(d *schema.ResourceData, meta interface{}, id string) error {
	instanceC, err := vpcClient(meta)
	if err != nil {
		return err
	}
	getinsOptions := &vpcv1.GetInstanceOptions{
		ID: &id,
	}
	getinsIniOptions := &vpcv1.GetInstanceInitializationOptions{
		ID: &id,
	}
	instance, response, err := instanceC.GetInstance(getinsOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error getting Instance: %s\n%s", err, response)
	}
	instanceInitialization, response, err := instanceC.GetInstanceInitialization(getinsIniOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error getting Instance initialization details: %s\n%s", err, response)
	}
	if instanceInitialization.DefaultTrustedProfile != nil && instanceInitialization.DefaultTrustedProfile.AutoLink != nil {
		d.Set(isInstanceDefaultTrustedProfileAutoLink, *instanceInitialization.DefaultTrustedProfile.AutoLink)
	}
	if instanceInitialization.DefaultTrustedProfile != nil && instanceInitialization.DefaultTrustedProfile.Target != nil {
		d.Set(isInstanceDefaultTrustedProfileTarget, *instanceInitialization.DefaultTrustedProfile.Target.ID)
	}

	if instance.AvailabilityPolicy != nil && instance.AvailabilityPolicy.HostFailure != nil {
		d.Set(isInstanceAvailablePolicyHostFailure, *instance.AvailabilityPolicy.HostFailure)
	}

	// catalog
	if instance.CatalogOffering != nil {
		versionCrn := *instance.CatalogOffering.Version.CRN
		catalogList := make([]map[string]interface{}, 0)
		catalogMap := map[string]interface{}{}
		catalogMap[isInstanceCatalogOfferingVersionCrn] = versionCrn
		catalogList = append(catalogList, catalogMap)
		d.Set(isInstanceCatalogOffering, catalogList)
	}
	d.Set(isInstanceName, *instance.Name)
	if instance.Profile != nil {
		d.Set(isInstanceProfile, *instance.Profile.Name)
	}
	cpuList := make([]map[string]interface{}, 0)
	if instance.Vcpu != nil {
		currentCPU := map[string]interface{}{}
		currentCPU[isInstanceCPUArch] = *instance.Vcpu.Architecture
		currentCPU[isInstanceCPUCount] = *instance.Vcpu.Count
		currentCPU[isInstanceCPUManufacturer] = instance.Vcpu.Manufacturer
		cpuList = append(cpuList, currentCPU)
	}
	d.Set(isInstanceCPU, cpuList)

	if instance.Bandwidth != nil {
		d.Set(isInstanceBandwidth, int(*instance.Bandwidth))
	}

	if instance.TotalNetworkBandwidth != nil {
		d.Set(isInstanceTotalNetworkBandwidth, int(*instance.TotalNetworkBandwidth))
	}

	if instance.TotalVolumeBandwidth != nil {
		d.Set(isInstanceTotalVolumeBandwidth, int(*instance.TotalVolumeBandwidth))
	}

	d.Set(isInstanceMemory, *instance.Memory)
	gpuList := make([]map[string]interface{}, 0)
	if instance.Gpu != nil {
		currentGpu := map[string]interface{}{}
		currentGpu[isInstanceGpuManufacturer] = instance.Gpu.Manufacturer
		currentGpu[isInstanceGpuModel] = instance.Gpu.Model
		currentGpu[isInstanceGpuCount] = instance.Gpu.Count
		currentGpu[isInstanceGpuMemory] = instance.Gpu.Memory
		gpuList = append(gpuList, currentGpu)
	}
	d.Set(isInstanceGpu, gpuList)

	if instance.PrimaryNetworkInterface != nil {
		primaryNicList := make([]map[string]interface{}, 0)
		currentPrimNic := map[string]interface{}{}
		currentPrimNic["id"] = *instance.PrimaryNetworkInterface.ID
		currentPrimNic[isInstanceNicName] = *instance.PrimaryNetworkInterface.Name

		//reserved ip changes
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
		getripoptions := &vpcv1.GetSubnetReservedIPOptions{
			SubnetID: instance.PrimaryNetworkInterface.Subnet.ID,
			ID:       instance.PrimaryNetworkInterface.PrimaryIP.ID,
		}
		insRip, response, err := instanceC.GetSubnetReservedIP(getripoptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error getting network interface reserved ip(%s) attached to the instance network interface(%s): %s\n%s", *instance.PrimaryNetworkInterface.PrimaryIP.ID, *instance.PrimaryNetworkInterface.ID, err, response)
		}
		currentPrimIp[isInstanceNicReservedIpAutoDelete] = insRip.AutoDelete

		primaryIpList = append(primaryIpList, currentPrimIp)
		currentPrimNic[isInstanceNicPrimaryIP] = primaryIpList

		getnicoptions := &vpcv1.GetInstanceNetworkInterfaceOptions{
			InstanceID: &id,
			ID:         instance.PrimaryNetworkInterface.ID,
		}
		insnic, response, err := instanceC.GetInstanceNetworkInterface(getnicoptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error getting network interfaces attached to the instance %s\n%s", err, response)
		}
		currentPrimNic[isInstanceNicAllowIPSpoofing] = *insnic.AllowIPSpoofing
		if insnic.PortSpeed != nil {
			currentPrimNic[isInstanceNicPortSpeed] = *insnic.PortSpeed
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
		d.Set(isInstancePrimaryNetworkInterface, primaryNicList)
	}

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

				getripoptions := &vpcv1.GetSubnetReservedIPOptions{
					SubnetID: intfc.Subnet.ID,
					ID:       intfc.PrimaryIP.ID,
				}
				insRip, response, err := instanceC.GetSubnetReservedIP(getripoptions)
				if err != nil {
					return fmt.Errorf("[ERROR] Error getting network interface reserved ip(%s) attached to the instance network interface(%s): %s\n%s", *intfc.PrimaryIP.ID, *intfc.ID, err, response)
				}
				currentPrimIp[isInstanceNicReservedIpAutoDelete] = insRip.AutoDelete

				primaryIpList = append(primaryIpList, currentPrimIp)
				currentNic[isInstanceNicPrimaryIP] = primaryIpList

				getnicoptions := &vpcv1.GetInstanceNetworkInterfaceOptions{
					InstanceID: &id,
					ID:         intfc.ID,
				}
				insnic, response, err := instanceC.GetInstanceNetworkInterface(getnicoptions)
				if err != nil {
					return fmt.Errorf("[ERROR] Error getting network interfaces attached to the instance %s\n%s", err, response)
				}
				currentNic[isInstanceNicAllowIPSpoofing] = *insnic.AllowIPSpoofing
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

		d.Set(isInstanceNetworkInterfaces, interfacesList)
	}

	if instance.Image != nil {
		d.Set(isInstanceImage, *instance.Image.ID)
	}

	d.Set(isInstanceStatus, *instance.Status)

	//set the status reasons
	if instance.StatusReasons != nil {
		statusReasonsList := make([]map[string]interface{}, 0)
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
		d.Set(isInstanceStatusReasons, statusReasonsList)
	}

	//set the lifecycle status, reasons
	if instance.LifecycleState != nil {
		d.Set(isInstanceLifecycleState, *instance.LifecycleState)
	}
	if instance.LifecycleReasons != nil {
		d.Set(isInstanceLifecycleReasons, dataSourceInstanceFlattenLifecycleReasons(instance.LifecycleReasons))
	}

	d.Set(isInstanceVPC, *instance.VPC.ID)
	d.Set(isInstanceZone, *instance.Zone.Name)

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
		d.Set(isInstanceVolumeAttachments, volList)
	}

	if instance.BootVolumeAttachment != nil {
		bootVolList := make([]map[string]interface{}, 0)
		bootVol := map[string]interface{}{}
		if instance.BootVolumeAttachment.Volume != nil {
			bootVol[isInstanceBootAttachmentName] = *instance.BootVolumeAttachment.Volume.Name
			bootVol[isInstanceBootVolumeId] = *instance.BootVolumeAttachment.Volume.ID

			instanceId := *instance.ID
			bootVolID := *instance.BootVolumeAttachment.ID
			getinsVolAttOptions := &vpcv1.GetInstanceVolumeAttachmentOptions{
				InstanceID: &instanceId,
				ID:         &bootVolID,
			}
			bootVolumeAtt, response, err := instanceC.GetInstanceVolumeAttachment(getinsVolAttOptions)
			if err != nil {
				return fmt.Errorf("[ERROR] Error getting Instance boot volume attachment : %s\n%s", err, response)
			}

			bootVol[isInstanceVolAttVolAutoDelete] = *bootVolumeAtt.DeleteVolumeOnInstanceDelete
			options := &vpcv1.GetVolumeOptions{
				ID: instance.BootVolumeAttachment.Volume.ID,
			}
			vol, response, err := instanceC.GetVolume(options)
			if err != nil {
				log.Printf("Error Getting Boot Volume (%s): %s\n%s", id, err, response)
			}
			if vol != nil {
				bootVol[isInstanceBootSize] = *vol.Capacity
				bootVol[isInstanceBootIOPS] = *vol.Iops
				bootVol[isInstanceBootProfile] = *vol.Profile.Name
				if vol.EncryptionKey != nil {
					bootVol[isInstanceBootEncryption] = *vol.EncryptionKey.CRN
				}
				if vol.SourceSnapshot != nil {
					bootVol[isInstanceVolumeSnapshot] = vol.SourceSnapshot.ID
				}
				if vol.UserTags != nil {
					bootVol[isInstanceBootVolumeTags] = vol.UserTags
				}
			}
		}
		bootVolList = append(bootVolList, bootVol)
		d.Set(isInstanceBootVolume, bootVolList)
	}
	tags, err := flex.GetGlobalTagsUsingCRN(meta, *instance.CRN, "", isInstanceUserTagType)
	if err != nil {
		log.Printf(
			"Error on get of resource Instance (%s) tags: %s", d.Id(), err)
	}
	d.Set(isInstanceTags, tags)
	accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *instance.CRN, "", isInstanceAccessTagType)
	if err != nil {
		log.Printf(
			"Error on get of resource Instance (%s) access tags: %s", d.Id(), err)
	}
	d.Set(isInstanceAccessTags, accesstags)
	controller, err := flex.GetBaseController(meta)
	if err != nil {
		return err
	}
	d.Set(flex.ResourceControllerURL, controller+"/vpc-ext/compute/vs")
	d.Set(flex.ResourceName, *instance.Name)
	d.Set(flex.ResourceCRN, *instance.CRN)
	d.Set(IsInstanceCRN, *instance.CRN)
	d.Set(flex.ResourceStatus, *instance.Status)
	if instance.ResourceGroup != nil {
		d.Set(isInstanceResourceGroup, *instance.ResourceGroup.ID)
		d.Set(flex.ResourceGroupName, *instance.ResourceGroup.Name)
	}
	if instance.MetadataService != nil {
		d.Set(isInstanceMetadataServiceEnabled, instance.MetadataService.Enabled)
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
		d.Set(isInstanceMetadataService, metadataService)
	}

	if instance.Disks != nil {
		disks := []map[string]interface{}{}
		for _, disksItem := range instance.Disks {
			disksItemMap := resourceIbmIsInstanceInstanceDiskToMap(disksItem)
			disks = append(disks, disksItemMap)
		}
		if err = d.Set(isInstanceDisks, disks); err != nil {
			return fmt.Errorf("[ERROR] Error setting disks: %s", err)
		}
	}

	placementTarget := []map[string]interface{}{}
	if instance.PlacementTarget != nil {
		placementTargetMap := resourceIbmIsInstanceInstancePlacementToMap(*instance.PlacementTarget.(*vpcv1.InstancePlacementTarget))
		placementTarget = append(placementTarget, placementTargetMap)
	}
	if err = d.Set(isInstancePlacementTarget, placementTarget); err != nil {
		return fmt.Errorf("[ERROR] Error setting placement_target: %s", err)
	}
	return nil
}

func instanceUpdate(d *schema.ResourceData, meta interface{}) error {
	instanceC, err := vpcClient(meta)
	if err != nil {
		return err
	}
	id := d.Id()

	bootVolSize := "boot_volume.0.size"
	if d.HasChange(bootVolSize) && !d.IsNewResource() {
		old, new := d.GetChange(bootVolSize)
		if new.(int) < old.(int) {
			return fmt.Errorf("[ERROR] Error while updating boot volume size of the instance, only expansion is possible")
		}
		bootVol := int64(new.(int))
		volId := d.Get("boot_volume.0.volume_id").(string)
		updateVolumeOptions := &vpcv1.UpdateVolumeOptions{
			ID: &volId,
		}
		volPatchModel := &vpcv1.VolumePatch{
			Capacity: &bootVol,
		}
		volPatchModelAsPatch, err := volPatchModel.AsPatch()

		if err != nil {
			return (fmt.Errorf("[ERROR] Error encountered while apply as patch for boot volume of instance %s", err))
		}

		updateVolumeOptions.VolumePatch = volPatchModelAsPatch

		vol, res, err := instanceC.UpdateVolume(updateVolumeOptions)

		if vol == nil || err != nil {
			return (fmt.Errorf("[ERROR] Error encountered while expanding boot volume of instance %s/n%s", err, res))
		}

		_, err = isWaitForVolumeAvailable(instanceC, volId, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
	}
	bootVolTags := "boot_volume.0.tags"
	if d.HasChange(bootVolTags) && !d.IsNewResource() {
		var userTags *schema.Set
		if v, ok := d.GetOk("boot_volume.0.tags"); ok {
			volId := d.Get("boot_volume.0.volume_id").(string)
			updateVolumeOptions := &vpcv1.UpdateVolumeOptions{
				ID: &volId,
			}
			userTags = v.(*schema.Set)
			if userTags != nil && userTags.Len() != 0 {
				userTagsArray := make([]string, userTags.Len())
				for i, userTag := range userTags.List() {
					userTagStr := userTag.(string)
					userTagsArray[i] = userTagStr
				}
				volumePatchModel := &vpcv1.VolumePatch{}
				volumePatchModel.UserTags = userTagsArray
				volumePatch, err := volumePatchModel.AsPatch()
				if err != nil {
					return fmt.Errorf("[ERROR] Error encountered while apply as patch for boot volume of instance %s", err)
				}
				optionsget := &vpcv1.GetVolumeOptions{
					ID: &volId,
				}
				_, response, err := instanceC.GetVolume(optionsget)
				if err != nil {
					return fmt.Errorf("[ERROR] Error getting Boot Volume (%s): %s\n%s", id, err, response)
				}
				eTag := response.Headers.Get("ETag")
				updateVolumeOptions.IfMatch = &eTag
				updateVolumeOptions.VolumePatch = volumePatch
				vol, res, err := instanceC.UpdateVolume(updateVolumeOptions)
				if vol == nil || err != nil {
					return (fmt.Errorf("[ERROR] Error encountered while applying tags for boot volume of instance %s/n%s", err, res))
				}
				_, err = isWaitForVolumeAvailable(instanceC, volId, d.Timeout(schema.TimeoutCreate))
				if err != nil {
					return err
				}
			}
		}
	}
	bootVolAutoDel := "boot_volume.0.auto_delete_volume"
	if d.HasChange(bootVolAutoDel) && !d.IsNewResource() {
		listvolattoptions := &vpcv1.ListInstanceVolumeAttachmentsOptions{
			InstanceID: &id,
		}
		vols, _, err := instanceC.ListInstanceVolumeAttachments(listvolattoptions)
		if err != nil {
			return err
		}

		auto_delete := d.Get(bootVolAutoDel).(bool)
		for _, vol := range vols.VolumeAttachments {
			if *vol.Type == "boot" {
				volAttachmentID := *vol.ID
				updateInstanceVolAttOptions := &vpcv1.UpdateInstanceVolumeAttachmentOptions{
					InstanceID: &id,
					ID:         &volAttachmentID,
				}
				volAttNamePatchModel := &vpcv1.VolumeAttachmentPatch{
					DeleteVolumeOnInstanceDelete: &auto_delete,
				}
				volAttNamePatchModelAsPatch, err := volAttNamePatchModel.AsPatch()
				if err != nil || volAttNamePatchModelAsPatch == nil {
					return fmt.Errorf("[ERROR] Error Instance volume attachment (%s) as patch : %s", id, err)
				}
				updateInstanceVolAttOptions.VolumeAttachmentPatch = volAttNamePatchModelAsPatch

				instanceVolAttUpdate, response, err := instanceC.UpdateInstanceVolumeAttachment(updateInstanceVolAttOptions)
				if err != nil || instanceVolAttUpdate == nil {
					log.Printf("[DEBUG] Instance volume attachment updation err %s\n%s", err, response)
					return err
				}
			}
		}
	}
	if d.HasChange(isPlacementTargetDedicatedHost) || d.HasChange(isPlacementTargetDedicatedHostGroup) && !d.IsNewResource() {
		dedicatedHost := d.Get(isPlacementTargetDedicatedHost).(string)
		dedicatedHostGroup := d.Get(isPlacementTargetDedicatedHostGroup).(string)
		actiontype := "stop"

		if dedicatedHost == "" && dedicatedHostGroup == "" {
			return fmt.Errorf("[ERROR] Error: Instances cannot be moved from private to public hosts")
		}

		createinsactoptions := &vpcv1.CreateInstanceActionOptions{
			InstanceID: &id,
			Type:       &actiontype,
		}
		_, response, err := instanceC.CreateInstanceAction(createinsactoptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return nil
			}
			return fmt.Errorf("[ERROR] Error Creating Instance Action: %s\n%s", err, response)
		}
		_, err = isWaitForInstanceActionStop(instanceC, d.Timeout(schema.TimeoutUpdate), id, d)
		if err != nil {
			return err
		}

		updateOptions := &vpcv1.UpdateInstanceOptions{
			ID: &id,
		}

		instancePatchModel := &vpcv1.InstancePatch{}

		if dedicatedHost != "" {
			placementTarget := &vpcv1.InstancePlacementTargetPatch{
				ID: &dedicatedHost,
			}
			instancePatchModel.PlacementTarget = placementTarget
		} else if dedicatedHostGroup != "" {
			placementTarget := &vpcv1.InstancePlacementTargetPatch{
				ID: &dedicatedHostGroup,
			}
			instancePatchModel.PlacementTarget = placementTarget
		}

		instancePatch, err := instancePatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("[ERROR] Error calling asPatch with total volume bandwidth for InstancePatch: %s", err)
		}

		updateOptions.InstancePatch = instancePatch

		_, _, err = instanceC.UpdateInstance(updateOptions)
		if err != nil {
			return err
		}

		actiontype = "start"
		createinsactoptions = &vpcv1.CreateInstanceActionOptions{
			InstanceID: &id,
			Type:       &actiontype,
		}
		_, response, err = instanceC.CreateInstanceAction(createinsactoptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return nil
			}
			return fmt.Errorf("[ERROR] Error Creating Instance Action: %s\n%s", err, response)
		}
		_, err = isWaitForInstanceActionStart(instanceC, d.Timeout(schema.TimeoutUpdate), id, d)
		if err != nil {
			return err
		}
	}

	if d.HasChange(isInstanceAction) && !d.IsNewResource() {

		actiontype := d.Get(isInstanceAction).(string)
		if actiontype != "" {
			getinsOptions := &vpcv1.GetInstanceOptions{
				ID: &id,
			}
			instance, response, err := instanceC.GetInstance(getinsOptions)
			if err != nil {
				return fmt.Errorf("[ERROR] Error Getting Instance (%s): %s\n%s", id, err, response)
			}
			if (actiontype == "stop" || actiontype == "reboot") && *instance.Status != isInstanceStatusRunning {
				d.Set(isInstanceAction, nil)
				return fmt.Errorf("[ERROR] Error with stop/reboot action: Cannot invoke stop/reboot action while instance is not in running state")
			} else if actiontype == "start" && *instance.Status != isInstanceActionStatusStopped {
				d.Set(isInstanceAction, nil)
				return fmt.Errorf("[ERROR] Error with start action: Cannot invoke start action while instance is not in stopped state")
			}
			createinsactoptions := &vpcv1.CreateInstanceActionOptions{
				InstanceID: &id,
				Type:       &actiontype,
			}
			if instanceActionForceIntf, ok := d.GetOk(isInstanceActionForce); ok {
				force := instanceActionForceIntf.(bool)
				createinsactoptions.Force = &force
			}
			_, response, err = instanceC.CreateInstanceAction(createinsactoptions)
			if err != nil {
				return fmt.Errorf("[ERROR] Error Creating Instance Action: %s\n%s", err, response)
			}
			if actiontype == "stop" {
				_, err = isWaitForInstanceActionStop(instanceC, d.Timeout(schema.TimeoutUpdate), id, d)
				if err != nil {
					return err
				}
			} else if actiontype == "start" || actiontype == "reboot" {
				_, err = isWaitForInstanceActionStart(instanceC, d.Timeout(schema.TimeoutUpdate), id, d)
				if err != nil {
					return err
				}
			}

		}
	}
	if d.HasChange(isInstanceVolumes) {
		old, new := d.GetChange(isInstanceVolumes)
		oldaddons := old.([]interface{})
		newaddons := new.([]interface{})
		var oldaddon, newaddon, add []string
		for _, v := range oldaddons {
			oldaddon = append(oldaddon, v.(string))
		}
		for _, v := range newaddons {
			newaddon = append(newaddon, v.(string))
		}
		// 1. Remove old addons no longer appearing in the new set
		// 2. Add new addons not already provisioned
		remove := flex.Listdifference(oldaddon, newaddon)
		add = flex.Listdifference(newaddon, oldaddon)
		var volautoDelete bool
		if volumeautodeleteIntf, ok := d.GetOk(isInstanceVolAttVolAutoDelete); ok && volumeautodeleteIntf != nil {
			volautoDelete = volumeautodeleteIntf.(bool)
		}

		if len(add) > 0 {
			for i := range add {
				createvolattoptions := &vpcv1.CreateInstanceVolumeAttachmentOptions{
					InstanceID: &id,
					Volume: &vpcv1.VolumeAttachmentPrototypeVolume{
						ID: &add[i],
					},
					DeleteVolumeOnInstanceDelete: &volautoDelete,
				}
				vol, _, err := instanceC.CreateInstanceVolumeAttachment(createvolattoptions)
				if err != nil {
					return fmt.Errorf("[ERROR] Error while attaching volume %q for instance %s: %q", add[i], d.Id(), err)
				}
				_, err = isWaitForInstanceVolumeAttached(instanceC, d, id, *vol.ID)
				if err != nil {
					return err
				}
			}

		}
		if len(remove) > 0 {
			for i := range remove {
				listvolattoptions := &vpcv1.ListInstanceVolumeAttachmentsOptions{
					InstanceID: &id,
				}
				vols, _, err := instanceC.ListInstanceVolumeAttachments(listvolattoptions)
				if err != nil {
					return err
				}
				for _, vol := range vols.VolumeAttachments {
					if *vol.Volume.ID == remove[i] {
						delvolattoptions := &vpcv1.DeleteInstanceVolumeAttachmentOptions{
							InstanceID: &id,
							ID:         vol.ID,
						}
						_, err := instanceC.DeleteInstanceVolumeAttachment(delvolattoptions)
						if err != nil {
							return fmt.Errorf("[ERROR] Error while removing volume %q for instance %s: %q", remove[i], d.Id(), err)
						}
						_, err = isWaitForInstanceVolumeDetached(instanceC, d, d.Id(), *vol.ID)
						if err != nil {
							return err
						}
						break
					}
				}
			}
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
				_, response, err := instanceC.CreateSecurityGroupTargetBinding(createsgnicoptions)
				if err != nil {
					return fmt.Errorf("[ERROR] Error while creating security group %q for primary network interface of instance %s\n%s: %q", add[i], d.Id(), err, response)
				}
				_, err = isWaitForInstanceAvailable(instanceC, d.Id(), d.Timeout(schema.TimeoutUpdate), d)
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
				response, err := instanceC.DeleteSecurityGroupTargetBinding(deletesgnicoptions)
				if err != nil {
					return fmt.Errorf("[ERROR] Error while removing security group %q for primary network interface of instance %s\n%s: %q", remove[i], d.Id(), err, response)
				}
				_, err = isWaitForInstanceAvailable(instanceC, d.Id(), d.Timeout(schema.TimeoutUpdate), d)
				if err != nil {
					return err
				}
			}
		}
	}

	if !d.IsNewResource() && (d.HasChange("primary_network_interface.0.primary_ip.0.name") || d.HasChange("primary_network_interface.0.primary_ip.0.auto_delete")) {
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
		_, response, err := instanceC.UpdateSubnetReservedIP(updateripoptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error updating instance network interface reserved ip(%s): %s\n%s", ripId, err, response)
		}
	}

	if (d.HasChange("primary_network_interface.0.allow_ip_spoofing") || d.HasChange("primary_network_interface.0.name")) && !d.IsNewResource() {
		newName := d.Get("primary_network_interface.0.name").(string)
		networkID := d.Get("primary_network_interface.0.id").(string)
		allowIPSpoofing := d.Get("primary_network_interface.0.allow_ip_spoofing").(bool)
		updatepnicfoptions := &vpcv1.UpdateInstanceNetworkInterfaceOptions{
			InstanceID: &id,
			ID:         &networkID,
		}

		networkInterfacePatchModel := &vpcv1.NetworkInterfacePatch{
			Name:            &newName,
			AllowIPSpoofing: &allowIPSpoofing,
		}
		networkInterfacePatch, err := networkInterfacePatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("[ERROR] Error calling asPatch for NetworkInterfacePatch: %s", err)
		}
		updatepnicfoptions.NetworkInterfacePatch = networkInterfacePatch

		_, response, err := instanceC.UpdateInstanceNetworkInterface(updatepnicfoptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error while updating name %s for primary network interface of instance %s\n%s: %q", newName, d.Id(), err, response)
		}
		_, err = isWaitForInstanceAvailable(instanceC, d.Id(), d.Timeout(schema.TimeoutUpdate), d)
		if err != nil {
			return err
		}
	}

	if d.HasChange(isInstanceNetworkInterfaces) && !d.IsNewResource() {
		nics := d.Get(isInstanceNetworkInterfaces).([]interface{})
		for i := range nics {
			securitygrpKey := fmt.Sprintf("network_interfaces.%d.security_groups", i)
			networkNameKey := fmt.Sprintf("network_interfaces.%d.name", i)
			subnetKey := fmt.Sprintf("network_interfaces.%d.subnet", i)
			ipSpoofingKey := fmt.Sprintf("network_interfaces.%d.allow_ip_spoofing", i)
			primaryipname := fmt.Sprintf("network_interfaces.%d.primary_ip.0.name", i)
			primaryipauto := fmt.Sprintf("network_interfaces.%d.primary_ip.0.auto_delete", i)
			primaryiprip := fmt.Sprintf("network_interfaces.%d.primary_ip.0.reserved_ip", i)
			if d.HasChange(primaryipname) || d.HasChange(primaryipauto) {
				subnetId := d.Get(subnetKey).(string)
				ripId := d.Get(primaryiprip).(string)
				updateripoptions := &vpcv1.UpdateSubnetReservedIPOptions{
					SubnetID: &subnetId,
					ID:       &ripId,
				}
				reservedIpPath := &vpcv1.ReservedIPPatch{}
				if d.HasChange(primaryipname) {
					name := d.Get(primaryipname).(string)
					reservedIpPath.Name = &name
				}
				if d.HasChange(primaryipauto) {
					auto := d.Get(primaryipauto).(bool)
					reservedIpPath.AutoDelete = &auto
				}
				reservedIpPathAsPatch, err := reservedIpPath.AsPatch()
				if err != nil {
					return fmt.Errorf("[ERROR] Error calling reserved ip as patch \n%s", err)
				}
				updateripoptions.ReservedIPPatch = reservedIpPathAsPatch
				_, response, err := instanceC.UpdateSubnetReservedIP(updateripoptions)
				if err != nil {
					return fmt.Errorf("[ERROR] Error updating instance network interface reserved ip(%s): %s\n%s", ripId, err, response)
				}
			}

			if d.HasChange(securitygrpKey) {
				ovs, nvs := d.GetChange(securitygrpKey)
				ov := ovs.(*schema.Set)
				nv := nvs.(*schema.Set)
				remove := flex.ExpandStringList(ov.Difference(nv).List())
				add := flex.ExpandStringList(nv.Difference(ov).List())
				if len(add) > 0 {
					networkIDKey := fmt.Sprintf("network_interfaces.%d.id", i)
					networkID := d.Get(networkIDKey).(string)
					for i := range add {
						createsgnicoptions := &vpcv1.CreateSecurityGroupTargetBindingOptions{
							SecurityGroupID: &add[i],
							ID:              &networkID,
						}
						_, response, err := instanceC.CreateSecurityGroupTargetBinding(createsgnicoptions)
						if err != nil {
							return fmt.Errorf("[ERROR] Error while creating security group %q for network interface of instance %s\n%s: %q", add[i], d.Id(), err, response)
						}
						_, err = isWaitForInstanceAvailable(instanceC, d.Id(), d.Timeout(schema.TimeoutUpdate), d)
						if err != nil {
							return err
						}
					}

				}
				if len(remove) > 0 {
					networkIDKey := fmt.Sprintf("network_interfaces.%d.id", i)
					networkID := d.Get(networkIDKey).(string)
					for i := range remove {
						deletesgnicoptions := &vpcv1.DeleteSecurityGroupTargetBindingOptions{
							SecurityGroupID: &remove[i],
							ID:              &networkID,
						}
						response, err := instanceC.DeleteSecurityGroupTargetBinding(deletesgnicoptions)
						if err != nil {
							return fmt.Errorf("[ERROR] Error while removing security group %q for network interface of instance %s\n%s: %q", remove[i], d.Id(), err, response)
						}
						_, err = isWaitForInstanceAvailable(instanceC, d.Id(), d.Timeout(schema.TimeoutUpdate), d)
						if err != nil {
							return err
						}
					}
				}

			}

			if d.HasChange(networkNameKey) || d.HasChange(ipSpoofingKey) {
				newName := d.Get(networkNameKey).(string)
				networkIDKey := fmt.Sprintf("network_interfaces.%d.id", i)
				networkID := d.Get(networkIDKey).(string)
				ipSpoofing := d.Get(ipSpoofingKey).(bool)
				updatepnicfoptions := &vpcv1.UpdateInstanceNetworkInterfaceOptions{
					InstanceID: &id,
					ID:         &networkID,
				}

				instancePatchModel := &vpcv1.NetworkInterfacePatch{
					Name:            &newName,
					AllowIPSpoofing: &ipSpoofing,
				}
				networkInterfacePatch, err := instancePatchModel.AsPatch()
				if err != nil {
					return fmt.Errorf("[ERROR] Error calling asPatch for NetworkInterfacePatch: %s", err)
				}
				updatepnicfoptions.NetworkInterfacePatch = networkInterfacePatch

				_, response, err := instanceC.UpdateInstanceNetworkInterface(updatepnicfoptions)
				if err != nil {
					return fmt.Errorf("[ERROR] Error while updating name %s for network interface of instance %s\n%s: %q", newName, d.Id(), err, response)
				}
				if err != nil {
					return err
				}
			}
		}

	}

	if d.HasChange(isInstanceTotalVolumeBandwidth) && !d.IsNewResource() {
		totalVolBandwidth := int64(d.Get(isInstanceTotalVolumeBandwidth).(int))
		updnetoptions := &vpcv1.UpdateInstanceOptions{
			ID: &id,
		}

		instancePatchModel := &vpcv1.InstancePatch{
			TotalVolumeBandwidth: &totalVolBandwidth,
		}
		instancePatch, err := instancePatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("[ERROR] Error calling asPatch with total volume bandwidth for InstancePatch: %s", err)
		}
		updnetoptions.InstancePatch = instancePatch

		_, _, err = instanceC.UpdateInstance(updnetoptions)
		if err != nil {
			return err
		}
	}

	if d.HasChange(isInstanceName) && !d.IsNewResource() {
		name := d.Get(isInstanceName).(string)
		updnetoptions := &vpcv1.UpdateInstanceOptions{
			ID: &id,
		}

		instancePatchModel := &vpcv1.InstancePatch{
			Name: &name,
		}
		instancePatch, err := instancePatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("[ERROR] Error calling asPatch for InstancePatch: %s", err)
		}
		updnetoptions.InstancePatch = instancePatch

		_, _, err = instanceC.UpdateInstance(updnetoptions)
		if err != nil {
			return err
		}
	}

	if d.HasChange(isInstanceMetadataServiceEnabled) && !d.IsNewResource() {
		enabled := d.Get(isInstanceMetadataServiceEnabled).(bool)
		updatedoptions := &vpcv1.UpdateInstanceOptions{
			ID: &id,
		}
		instancePatchModel := &vpcv1.InstancePatch{
			MetadataService: &vpcv1.InstanceMetadataServicePatch{
				Enabled: &enabled,
			},
		}
		instancePatch, err := instancePatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("Error calling asPatch for InstancePatch: %s", err)
		}
		updatedoptions.InstancePatch = instancePatch

		_, _, err = instanceC.UpdateInstance(updatedoptions)
		if err != nil {
			return err
		}
	}

	if d.HasChange(isInstanceMetadataService) && !d.IsNewResource() {
		metadataServiceIntf := d.Get(isInstanceMetadataService)
		updatedoptions := &vpcv1.UpdateInstanceOptions{
			ID: &id,
		}
		metadataServicePatchModel := &vpcv1.InstanceMetadataServicePatch{}
		instancePatchModel := &vpcv1.InstancePatch{}
		metadataServiceMap := metadataServiceIntf.([]interface{})[0].(map[string]interface{})
		if d.HasChange(isInstanceMetadataService + ".0." + isInstanceMetadataServiceEnabled1) {
			enabledIntf, ok := metadataServiceMap[isInstanceMetadataServiceEnabled1]
			if ok {
				enabled := enabledIntf.(bool)
				metadataServicePatchModel.Enabled = &enabled
			}
		}
		if d.HasChange(isInstanceMetadataService + ".0." + isInstanceMetadataServiceProtocol) {
			protocolIntf, ok := metadataServiceMap[isInstanceMetadataServiceProtocol]
			if ok {
				protocol := protocolIntf.(string)
				metadataServicePatchModel.Protocol = &protocol
			}
		}
		if d.HasChange(isInstanceMetadataService + ".0." + isInstanceMetadataServiceRespHopLimit) {
			respHopLimitIntf, ok := metadataServiceMap[isInstanceMetadataServiceRespHopLimit]
			if ok {
				respHopLimit := int64(respHopLimitIntf.(int))
				metadataServicePatchModel.ResponseHopLimit = &respHopLimit
			}
		}
		instancePatchModel.MetadataService = metadataServicePatchModel

		instancePatch, err := instancePatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("Error calling asPatch for InstancePatch: %s", err)
		}
		updatedoptions.InstancePatch = instancePatch

		_, _, err = instanceC.UpdateInstance(updatedoptions)
		if err != nil {
			return err
		}
	}

	if d.HasChange(isInstanceAvailablePolicyHostFailure) && !d.IsNewResource() {

		updatedoptions := &vpcv1.UpdateInstanceOptions{
			ID: &id,
		}
		availablePolicyHostFailure := d.Get(isInstanceAvailablePolicyHostFailure).(string)
		instancePatchModel := &vpcv1.InstancePatch{
			AvailabilityPolicy: &vpcv1.InstanceAvailabilityPolicyPatch{
				HostFailure: &availablePolicyHostFailure,
			},
		}
		instancePatch, err := instancePatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("Error calling asPatch for InstancePatch: %s", err)
		}
		updatedoptions.InstancePatch = instancePatch

		_, _, err = instanceC.UpdateInstance(updatedoptions)
		if err != nil {
			return err
		}
	}

	if d.HasChange(isInstanceProfile) && !d.IsNewResource() {

		getinsOptions := &vpcv1.GetInstanceOptions{
			ID: &id,
		}
		instance, response, err := instanceC.GetInstance(getinsOptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("[ERROR] Error Getting Instance (%s): %s\n%s", id, err, response)
		}

		if instance != nil && *instance.Status == "running" {
			actiontype := "stop"
			createinsactoptions := &vpcv1.CreateInstanceActionOptions{
				InstanceID: &id,
				Type:       &actiontype,
			}
			_, response, err = instanceC.CreateInstanceAction(createinsactoptions)
			if err != nil {
				if response != nil && response.StatusCode == 404 {
					return nil
				}
				return fmt.Errorf("[ERROR] Error Creating Instance Action: %s\n%s", err, response)
			}
			_, err = isWaitForInstanceActionStop(instanceC, d.Timeout(schema.TimeoutUpdate), id, d)
			if err != nil {
				return err
			}
		}

		updnetoptions := &vpcv1.UpdateInstanceOptions{
			ID: &id,
		}

		instanceProfile := d.Get(isInstanceProfile).(string)
		profile := &vpcv1.InstancePatchProfile{
			Name: &instanceProfile,
		}
		instancePatchModel := &vpcv1.InstancePatch{
			Profile: profile,
		}
		instancePatch, err := instancePatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("[ERROR] Error calling asPatch for InstancePatch: %s", err)
		}
		updnetoptions.InstancePatch = instancePatch

		_, response, err = instanceC.UpdateInstance(updnetoptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error in UpdateInstancePatch: %s\n%s", err, response)
		}

		actiontype := "start"
		createinsactoptions := &vpcv1.CreateInstanceActionOptions{
			InstanceID: &id,
			Type:       &actiontype,
		}
		_, response, err = instanceC.CreateInstanceAction(createinsactoptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return nil
			}
			return fmt.Errorf("[ERROR] Error Creating Instance Action: %s\n%s", err, response)
		}
		_, err = isWaitForInstanceAvailable(instanceC, d.Id(), d.Timeout(schema.TimeoutUpdate), d)
		if err != nil {
			return err
		}

	}

	getinsOptions := &vpcv1.GetInstanceOptions{
		ID: &id,
	}
	instance, response, err := instanceC.GetInstance(getinsOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Getting Instance: %s\n%s", err, response)
	}
	if d.HasChange(isInstanceTags) {
		oldList, newList := d.GetChange(isInstanceTags)
		err = flex.UpdateTagsUsingCRN(oldList, newList, meta, *instance.CRN)
		if err != nil {
			log.Printf(
				"Error on update of resource Instance (%s) tags: %s", d.Id(), err)
		}
	}
	if d.HasChange(isInstanceAccessTags) {
		oldList, newList := d.GetChange(isInstanceAccessTags)
		err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *instance.CRN, "", isInstanceAccessTagType)
		if err != nil {
			log.Printf(
				"Error on update of resource Instance (%s) access tags: %s", d.Id(), err)
		}
	}
	return nil
}

func resourceIBMisInstanceUpdate(d *schema.ResourceData, meta interface{}) error {

	err := instanceUpdate(d, meta)
	if err != nil {
		return err
	}

	return resourceIBMisInstanceRead(d, meta)
}

func instanceDelete(d *schema.ResourceData, meta interface{}, id string) error {
	instanceC, err := vpcClient(meta)
	if err != nil {
		return err
	}

	cleanDelete := d.Get(isEnableCleanDelete).(bool)
	getinsOptions := &vpcv1.GetInstanceOptions{
		ID: &id,
	}
	_, response, err := instanceC.GetInstance(getinsOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error Getting Instance (%s): %s\n%s", id, err, response)
	}

	bootvolid := ""

	if cleanDelete {
		actiontype := "stop"
		createinsactoptions := &vpcv1.CreateInstanceActionOptions{
			InstanceID: &id,
			Type:       &actiontype,
		}
		_, response, err = instanceC.CreateInstanceAction(createinsactoptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return nil
			}
			return fmt.Errorf("[ERROR] Error Creating Instance Action: %s\n%s", err, response)
		}
		_, err = isWaitForInstanceActionStop(instanceC, d.Timeout(schema.TimeoutDelete), id, d)
		if err != nil {
			return err
		}
		listvolattoptions := &vpcv1.ListInstanceVolumeAttachmentsOptions{
			InstanceID: &id,
		}
		vols, response, err := instanceC.ListInstanceVolumeAttachments(listvolattoptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Listing volume attachments to the instance: %s\n%s", err, response)
		}
		for _, vol := range vols.VolumeAttachments {
			if *vol.Type == "data" {
				delvolattoptions := &vpcv1.DeleteInstanceVolumeAttachmentOptions{
					InstanceID: &id,
					ID:         vol.ID,
				}
				_, err := instanceC.DeleteInstanceVolumeAttachment(delvolattoptions)
				if err != nil {
					return fmt.Errorf("[ERROR] Error while removing volume Attachment %q for instance %s: %q", *vol.ID, d.Id(), err)
				}
				_, err = isWaitForInstanceVolumeDetached(instanceC, d, d.Id(), *vol.ID)
				if err != nil {
					return err
				}
			}
			if *vol.Type == "boot" {
				bootvolid = *vol.Volume.ID
			}
		}
	}
	deleteinstanceOptions := &vpcv1.DeleteInstanceOptions{
		ID: &id,
	}
	_, err = instanceC.DeleteInstance(deleteinstanceOptions)
	if err != nil {
		return err
	}
	if cleanDelete {
		_, err = isWaitForInstanceDelete(instanceC, d, d.Id())
		if err != nil {
			return err
		}
		if _, ok := d.GetOk(isInstanceBootVolume); ok {
			autoDel := d.Get("boot_volume.0.auto_delete_volume").(bool)
			if autoDel {
				_, err = isWaitForVolumeDeleted(instanceC, bootvolid, d.Timeout(schema.TimeoutDelete))
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func resourceIBMisInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	id := d.Id()
	err := instanceDelete(d, meta, id)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func instanceExists(d *schema.ResourceData, meta interface{}, id string) (bool, error) {
	instanceC, err := vpcClient(meta)
	if err != nil {
		return false, err
	}
	getinsOptions := &vpcv1.GetInstanceOptions{
		ID: &id,
	}
	_, response, err := instanceC.GetInstance(getinsOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("[ERROR] Error Getting Instance: %s\n%s", err, response)
	}
	return true, nil
}

func resourceIBMisInstanceExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	id := d.Id()

	exists, err := instanceExists(d, meta, id)
	return exists, err

}

func isWaitForInstanceDelete(instanceC *vpcv1.VpcV1, d *schema.ResourceData, id string) (interface{}, error) {

	stateConf := &resource.StateChangeConf{
		Pending: []string{isInstanceDeleting, isInstanceAvailable},
		Target:  []string{isInstanceDeleteDone, ""},
		Refresh: func() (interface{}, string, error) {
			getinsoptions := &vpcv1.GetInstanceOptions{
				ID: &id,
			}
			instance, response, err := instanceC.GetInstance(getinsoptions)
			if err != nil {
				if response != nil && response.StatusCode == 404 {
					return instance, isInstanceDeleteDone, nil
				}
				return nil, "", fmt.Errorf("[ERROR] Error Getting Instance: %s\n%s", err, response)
			}
			if *instance.Status == isInstanceFailed {
				return instance, *instance.Status, fmt.Errorf("[ERROR] The  instance %s failed to delete: %v", d.Id(), err)
			}
			return instance, isInstanceDeleting, nil
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isWaitForInstanceActionStop(instanceC *vpcv1.VpcV1, timeout time.Duration, id string, d *schema.ResourceData) (interface{}, error) {
	communicator := make(chan interface{})
	stateConf := &resource.StateChangeConf{
		Pending: []string{isInstanceStatusRunning, isInstanceStatusPending, isInstanceActionStatusStopping},
		Target:  []string{isInstanceActionStatusStopped, isInstanceStatusFailed, ""},
		Refresh: func() (interface{}, string, error) {
			getinsoptions := &vpcv1.GetInstanceOptions{
				ID: &id,
			}
			instance, response, err := instanceC.GetInstance(getinsoptions)
			if err != nil {
				return nil, "", fmt.Errorf("[ERROR] Error Getting Instance: %s\n%s", err, response)
			}
			select {
			case data := <-communicator:
				return nil, "", data.(error)
			default:
				fmt.Println("no message sent")
			}
			if *instance.Status == isInstanceStatusFailed {
				// let know the isRestartStopAction() to stop
				close(communicator)
				return instance, *instance.Status, fmt.Errorf("[ERROR] The  instance %s failed to stop: %v", id, err)
			}
			return instance, *instance.Status, nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	if v, ok := d.GetOk("force_recovery_time"); ok {
		forceTimeout := v.(int)
		go isRestartStopAction(instanceC, id, d, forceTimeout, communicator)
	}

	return stateConf.WaitForState()
}

func isWaitForInstanceActionStart(instanceC *vpcv1.VpcV1, timeout time.Duration, id string, d *schema.ResourceData) (interface{}, error) {
	communicator := make(chan interface{})
	stateConf := &resource.StateChangeConf{
		Pending: []string{isInstanceActionStatusStopped, isInstanceStatusPending, isInstanceActionStatusStopping, isInstanceStatusStarting, isInstanceStatusRestarting},
		Target:  []string{isInstanceStatusRunning, isInstanceStatusFailed, ""},
		Refresh: func() (interface{}, string, error) {
			getinsoptions := &vpcv1.GetInstanceOptions{
				ID: &id,
			}
			instance, response, err := instanceC.GetInstance(getinsoptions)
			if err != nil {
				return nil, "", fmt.Errorf("[ERROR] Error Getting Instance: %s\n%s", err, response)
			}
			select {
			case data := <-communicator:
				return nil, "", data.(error)
			default:
				fmt.Println("no message sent")
			}
			if *instance.Status == isInstanceStatusFailed {
				// let know the isRestartStopAction() to stop
				close(communicator)
				return instance, *instance.Status, fmt.Errorf("[ERROR] The  instance %s failed to start: %v", id, err)
			}
			return instance, *instance.Status, nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	if v, ok := d.GetOk("force_recovery_time"); ok {
		forceTimeout := v.(int)
		go isRestartStopAction(instanceC, id, d, forceTimeout, communicator)
	}

	return stateConf.WaitForState()
}

func isRestartStopAction(instanceC *vpcv1.VpcV1, id string, d *schema.ResourceData, forceTimeout int, communicator chan interface{}) {
	subticker := time.NewTicker(time.Duration(forceTimeout) * time.Minute)
	//subticker := time.NewTicker(time.Duration(forceTimeout) * time.Second)
	for {
		select {

		case <-subticker.C:
			log.Println("Instance is still in stopping state, retrying to stop with -force")
			actiontype := "stop"
			createinsactoptions := &vpcv1.CreateInstanceActionOptions{
				InstanceID: &id,
				Type:       &actiontype,
			}
			_, response, err := instanceC.CreateInstanceAction(createinsactoptions)
			if err != nil {
				communicator <- fmt.Errorf("[ERROR] Error retrying instance action stop: %s\n%s", err, response)
				return
			}
		case <-communicator:
			// indicates refresh func is reached target and not proceed with the thread)
			subticker.Stop()
			return

		}
	}
}

func isWaitForInstanceVolumeAttached(instanceC *vpcv1.VpcV1, d *schema.ResourceData, id, volID string) (interface{}, error) {
	log.Printf("Waiting for instance (%s) volume (%s) to be attached.", id, volID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{isInstanceVolumeAttaching},
		Target:     []string{isInstanceVolumeAttached, ""},
		Refresh:    isInstanceVolumeRefreshFunc(instanceC, id, volID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isInstanceVolumeRefreshFunc(instanceC *vpcv1.VpcV1, id, volID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getvolattoptions := &vpcv1.GetInstanceVolumeAttachmentOptions{
			InstanceID: &id,
			ID:         &volID,
		}
		vol, response, err := instanceC.GetInstanceVolumeAttachment(getvolattoptions)
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] Error Attaching volume: %s\n%s", err, response)
		}

		if *vol.Status == isInstanceVolumeAttached {
			return vol, isInstanceVolumeAttached, nil
		}

		return vol, isInstanceVolumeAttaching, nil
	}
}

func isWaitForInstanceVolumeDetached(instanceC *vpcv1.VpcV1, d *schema.ResourceData, id, volID string) (interface{}, error) {

	stateConf := &resource.StateChangeConf{
		Pending: []string{isInstanceVolumeAttached, isInstanceVolumeDetaching},
		Target:  []string{isInstanceDeleteDone, ""},
		Refresh: func() (interface{}, string, error) {
			getvolattoptions := &vpcv1.GetInstanceVolumeAttachmentOptions{
				InstanceID: &id,
				ID:         &volID,
			}
			vol, response, err := instanceC.GetInstanceVolumeAttachment(getvolattoptions)
			if err != nil {
				if response != nil && response.StatusCode == 404 {
					return vol, isInstanceDeleteDone, nil
				}
				return nil, "", fmt.Errorf("[ERROR] Error Detaching: %s\n%s", err, response)
			}
			if *vol.Status == isInstanceFailed {
				return vol, *vol.Status, fmt.Errorf("[ERROR] The instance %s failed to detach volume %s: %v", d.Id(), volID, err)
			}
			return vol, isInstanceVolumeDetaching, nil
		},
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func resourceIbmIsInstanceInstanceDiskToMap(instanceDisk vpcv1.InstanceDisk) map[string]interface{} {
	instanceDiskMap := map[string]interface{}{}

	instanceDiskMap["created_at"] = instanceDisk.CreatedAt.String()
	instanceDiskMap["href"] = instanceDisk.Href
	instanceDiskMap["id"] = instanceDisk.ID
	instanceDiskMap["interface_type"] = instanceDisk.InterfaceType
	instanceDiskMap["name"] = instanceDisk.Name
	instanceDiskMap["resource_type"] = instanceDisk.ResourceType
	instanceDiskMap["size"] = flex.IntValue(instanceDisk.Size)

	return instanceDiskMap
}

func suppressEnableCleanDelete(k, old, new string, d *schema.ResourceData) bool {
	// During import
	if old == "" && d.Id() != "" {
		return true
	}
	return false
}

func resourceIbmIsInstanceInstancePlacementToMap(instancePlacement vpcv1.InstancePlacementTarget) map[string]interface{} {
	instancePlacementMap := map[string]interface{}{}

	instancePlacementMap["crn"] = instancePlacement.CRN
	if instancePlacement.Deleted != nil {
		DeletedMap := resourceIbmIsInstanceDedicatedHostGroupReferenceDeletedToMap(*instancePlacement.Deleted)
		instancePlacementMap["deleted"] = []map[string]interface{}{DeletedMap}
	}
	instancePlacementMap["href"] = instancePlacement.Href
	instancePlacementMap["id"] = instancePlacement.ID
	instancePlacementMap["name"] = instancePlacement.Name
	instancePlacementMap["resource_type"] = instancePlacement.ResourceType

	return instancePlacementMap
}

func resourceIbmIsInstanceDedicatedHostGroupReferenceDeletedToMap(dedicatedHostGroupReferenceDeleted vpcv1.DedicatedHostGroupReferenceDeleted) map[string]interface{} {
	dedicatedHostGroupReferenceDeletedMap := map[string]interface{}{}

	dedicatedHostGroupReferenceDeletedMap["more_info"] = dedicatedHostGroupReferenceDeleted.MoreInfo

	return dedicatedHostGroupReferenceDeletedMap
}

func GetInstanceMetadataServiceOptions(d *schema.ResourceData) (metadataService *vpcv1.InstanceMetadataServicePrototype) {

	if metadataServiceIntf, ok := d.GetOk(isInstanceMetadataService); ok {
		metadataService = &vpcv1.InstanceMetadataServicePrototype{}
		metadataServiceMap := metadataServiceIntf.([]interface{})[0].(map[string]interface{})
		enabledIntf, ok := metadataServiceMap[isInstanceMetadataServiceEnabled1]

		if ok {
			enabled := enabledIntf.(bool)
			metadataService.Enabled = &enabled
		}
		protocolIntf, ok := metadataServiceMap[isInstanceMetadataServiceProtocol]
		if ok && protocolIntf.(string) != "" {
			protocol := protocolIntf.(string)
			metadataService.Protocol = &protocol
		}
		respHopLimitIntf, ok := metadataServiceMap[isInstanceMetadataServiceRespHopLimit]
		if ok && int64(respHopLimitIntf.(int)) != 0 {
			respHopLimit := int64(respHopLimitIntf.(int))
			metadataService.ResponseHopLimit = &respHopLimit
		}
		return
	}
	return nil
}
