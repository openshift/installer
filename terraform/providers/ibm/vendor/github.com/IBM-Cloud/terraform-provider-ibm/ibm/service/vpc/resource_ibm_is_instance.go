// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
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
	isInstanceVolumeSnapshotCrn       = "snapshot_crn"
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
	isInstanceCatalogOfferingPlanCrn     = "plan_crn"

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

			// cluster changes
			"cluster_network_attachments": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				// DiffSuppressFunc: diffSuppressClusterNetworkAttachment,
				Description: "The cluster network attachments for this virtual server instance.The cluster network attachments are ordered for consistent instance configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_network_interface": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Required:    true,
							Description: "A cluster network interface for the instance cluster network attachment. This can bespecified using an existing cluster network interface that does not already have a `target`,or a prototype object for a new cluster network interface.This instance must reside in the same VPC as the specified cluster network interface. Thecluster network interface must reside in the same cluster network as the`cluster_network_interface` of any other `cluster_network_attachments` for this instance.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"auto_delete": &schema.Schema{
										Type:                  schema.TypeBool,
										Optional:              true,
										Computed:              true,
										DiffSuppressOnRefresh: true,
										DiffSuppressFunc:      flex.ApplyOnlyOnce,
										Description:           "Indicates whether this cluster network interface will be automatically deleted when `target` is deleted.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The name for this cluster network interface. The name must not be used by another interface in the cluster network. Names beginning with `ibm-` are reserved for provider-owned resources, and are not allowed. If unspecified, the name will be a hyphenated list of randomly-selected words.",
									},
									"primary_ip": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Computed:    true,
										Description: "The primary IP address to bind to the cluster network interface. May be eithera cluster network subnet reserved IP identity, or a cluster network subnet reserved IPprototype object which will be used to create a new cluster network subnet reserved IP.If a cluster network subnet reserved IP identity is provided, the specified clusternetwork subnet reserved IP must be unbound.If a cluster network subnet reserved IP prototype object with an address is provided,the address must be available on the cluster network interface's cluster networksubnet. If no address is specified, an available address on the cluster network subnetwill be automatically selected and reserved.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
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
													Optional:    true,
													Computed:    true,
													Description: "The IP address to reserve, which must not already be reserved on the subnet.If unspecified, an available address on the subnet will automatically be selected.",
												},
												"auto_delete": &schema.Schema{
													Type:             schema.TypeBool,
													Optional:         true,
													Computed:         true,
													DiffSuppressFunc: flex.ApplyOnlyOnce,
													Description:      "Indicates whether this cluster network subnet reserved IP member will be automatically deleted when either `target` is deleted, or the cluster network subnet reserved IP is unbound.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "The name for this cluster network subnet reserved IP. The name must not be used by another reserved IP in the cluster network subnet. Names starting with `ibm-` are reserved for provider-owned resources, and are not allowed. If unspecified, the name will be a hyphenated list of randomly-selected words.",
												},
											},
										},
									},
									"subnet": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Computed:    true,
										Description: "The associated cluster network subnet. Required if `primary_ip` does not specify acluster network subnet reserved IP identity.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "The unique identifier for this cluster network subnet.",
												},
												"href": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "The URL for this cluster network subnet.",
												},
											},
										},
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
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
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: flex.ApplyOnlyOnce,
							Description:      "The name for this cluster network attachment. Names must be unique within the instance the cluster network attachment resides in. If unspecified, the name will be a hyphenated list of randomly-selected words. Names starting with `ibm-` are reserved for provider-owned resources, and are not allowed.",
						},
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
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
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
							Optional:    true,
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

			"confidential_compute_mode": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_instance", "confidential_compute_mode"),
				Description:  "The confidential compute mode to use for this virtual server instance.If unspecified, the default confidential compute mode from the profile will be used.",
			},
			"enable_secure_boot": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates whether secure boot is enabled for this virtual server instance.If unspecified, the default secure boot mode from the profile will be used.",
			},

			isInstanceSourceTemplate: {
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				AtLeastOneOf:  []string{isInstanceImage, isInstanceSourceTemplate, "boot_volume.0.snapshot", "boot_volume.0.snapshot_crn", "catalog_offering.0.offering_crn", "catalog_offering.0.version_crn", "boot_volume.0.volume_id"},
				ConflictsWith: []string{"boot_volume.0.snapshot", "boot_volume.0.snapshot_crn", "boot_volume.0.volume_id"},
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
							RequiredWith:  []string{isInstanceZone, isInstanceVPC, isInstanceProfile},
							Description:   "Identifies a catalog offering by a unique CRN property",
						},
						isInstanceCatalogOfferingVersionCrn: {
							Type:          schema.TypeString,
							Optional:      true,
							ForceNew:      true,
							ConflictsWith: []string{"catalog_offering.0.offering_crn"},
							RequiredWith:  []string{isInstanceZone, isInstanceVPC, isInstanceProfile},
							Description:   "Identifies a version of a catalog offering by a unique CRN property",
						},
						isInstanceCatalogOfferingPlanCrn: {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							ForceNew:    true,
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

			isInstancePrimaryNetworkInterface: {
				Type:          schema.TypeList,
				MinItems:      1,
				MaxItems:      1,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"primary_network_attachment", "network_attachments"},
				Description:   "Primary Network interface info",
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

			// volume_prototypes
			"volume_prototypes": {
				Type:             schema.TypeList,
				Optional:         true,
				DiffSuppressFunc: diffSuppressVolumePrototypes,
				ConflictsWith:    []string{isInstanceVolumes},
				Computed:         true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: flex.ApplyOnce,
						},
						"delete_volume_on_instance_delete": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"volume_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"volume_crn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume_resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume_iops": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The maximum I/O operations per second (IOPS) for the volume.",
						},
						"volume_profile": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The  globally unique name for the volume profile to use for this volume.",
						},
						"volume_capacity": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The capacity of the volume in gigabytes. The specified minimum and maximum capacity values for creating or updating volumes may expand in the future.",
						},
						"volume_source_snapshot": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The snapshot from which to clone the volume",
						},
						"volume_encryption_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The CRN of the [Key Protect Root Key](https://cloud.ibm.com/docs/key-protect?topic=key-protect-getting-started-tutorial) or [Hyper Protect Crypto Service Root Key](https://cloud.ibm.com/docs/hs-crypto?topic=hs-crypto-get-started) for this resource.",
						},
						"volume_tags": {
							Type:        schema.TypeSet,
							Optional:    true,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_instance_template", "tags")},
							Set:         flex.ResourceIBMVPCHash,
							Description: "UserTags for the volume instance",
						},
					},
				},
			},

			"primary_network_attachment": &schema.Schema{
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				Description:   "The primary network attachment for this virtual server instance.",
				ExactlyOneOf:  []string{"primary_network_attachment", "primary_network_interface"},
				ConflictsWith: []string{"primary_network_interface", "network_interfaces"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						// pna can accept either vni id or prototype
						"name": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validate.InvokeValidator("ibm_is_instance_network_attachment", "name"),
							Description:  "The name for this instance network attachment. The name is unique across all network attachments for the instance.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this instance network attachment.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
						// primary_ip for consistency
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
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this instance network attachment.",
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

						// vni properties
						"virtual_network_interface": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Computed:    true,
							Description: "A virtual network interface for the instance network attachment. This can be specified using an existing virtual network interface, or a prototype object for a new virtual network interface.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The virtual network interface id for this instance network attachment.",
									},
									"allow_ip_spoofing": &schema.Schema{
										Type:          schema.TypeBool,
										Optional:      true,
										ConflictsWith: []string{"primary_network_attachment.0.virtual_network_interface.0.id"},
										Computed:      true,
										Description:   "Indicates whether source IP spoofing is allowed on this interface. If `false`, source IP spoofing is prevented on this interface. If `true`, source IP spoofing is allowed on this interface.",
									},
									"auto_delete": &schema.Schema{
										Type:          schema.TypeBool,
										Optional:      true,
										ConflictsWith: []string{"primary_network_attachment.0.virtual_network_interface.0.id"},
										Computed:      true,
										Description:   "Indicates whether this virtual network interface will be automatically deleted when`target` is deleted.",
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
										ConflictsWith: []string{"primary_network_attachment.0.virtual_network_interface.0.id"},
										Computed:      true,
										Set:           hashIpsList,
										Description:   "The reserved IPs bound to this virtual network interface.May be empty when `lifecycle_state` is `pending`.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"address": &schema.Schema{
													Type: schema.TypeString,
													// Optional:    true,
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
													Type: schema.TypeBool,
													// Optional:    true,
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
													Type: schema.TypeString,
													// Optional:    true,
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
										ConflictsWith: []string{"primary_network_attachment.0.virtual_network_interface.0.id"},
										Computed:      true,
										Description:   "The primary IP address of the virtual network interface for the instance networkattachment.",
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
												"auto_delete": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Default:     true,
													Description: "Indicates whether this reserved ip will be automatically deleted when `target` is deleted.",
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
										ConflictsWith: []string{"primary_network_attachment.0.virtual_network_interface.0.id"},
										Computed:      true,
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

			isInstanceNetworkInterfaces: {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"primary_network_attachment", "network_attachments"},
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

			"network_attachments": &schema.Schema{
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"primary_network_interface", "network_interfaces"},
				Description:   "The network attachments for this virtual server instance, including the primary network attachment.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						// pna can accept either vni id or prototype
						"name": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validate.InvokeValidator("ibm_is_instance_network_attachment", "name"),
							Description:  "The name for this instance network attachment. The name is unique across all network attachments for the instance.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this instance network attachment.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this instance network attachment.",
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
						// primary_ip for consistency
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
						"virtual_network_interface": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "A virtual network interface for the instance network attachment. This can be specified using an existing virtual network interface, or a prototype object for a new virtual network interface.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The virtual network interface id for this instance network attachment.",
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
										Computed:    true,
										Set:         hashIpsList,
										Description: "The reserved IPs bound to this virtual network interface.May be empty when `lifecycle_state` is `pending`.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"address": &schema.Schema{
													Type: schema.TypeString,
													// Optional:    true,
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
													Type: schema.TypeBool,
													// Optional:    true,
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
													Type: schema.TypeString,
													// Optional:    true,
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
										Description: "The primary IP address of the virtual network interface for the instance networkattachment.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"address": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "The IP address.If the address has not yet been selected, the value will be `0.0.0.0`.This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.",
												},
												"auto_delete": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Default:     true,
													Description: "Indicates whether this reserved IP member will be automatically deleted when either target is deleted, or the reserved IP is unbound.",
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
				ConflictsWith: []string{"boot_volume.0.snapshot", "boot_volume.0.snapshot_crn", "catalog_offering.0.offering_crn", "catalog_offering.0.version_crn", "boot_volume.0.volume_id"},
				AtLeastOneOf:  []string{isInstanceImage, isInstanceSourceTemplate, "boot_volume.0.snapshot", "boot_volume.0.snapshot_crn", "catalog_offering.0.offering_crn", "catalog_offering.0.version_crn", "boot_volume.0.volume_id"},
				RequiredWith:  []string{isInstanceZone, isInstanceVPC, isInstanceProfile},
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
							RequiredWith:  []string{isInstanceZone, isInstanceProfile, isInstanceVPC},
							AtLeastOneOf:  []string{isInstanceImage, isInstanceSourceTemplate, "boot_volume.0.volume_id", "boot_volume.0.snapshot", "boot_volume.0.snapshot_crn", "catalog_offering.0.offering_crn", "catalog_offering.0.version_crn"},
							ConflictsWith: []string{isInstanceImage, isInstanceSourceTemplate, "boot_volume.0.snapshot", "boot_volume.0.snapshot_crn", "boot_volume.0.name", "boot_volume.0.encryption", "catalog_offering.0.offering_crn", "catalog_offering.0.version_crn"},
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
							RequiredWith:  []string{isInstanceZone, isInstanceProfile, isInstanceVPC},
							AtLeastOneOf:  []string{isInstanceImage, isInstanceSourceTemplate, "boot_volume.0.snapshot", "boot_volume.0.snapshot_crn", "catalog_offering.0.offering_crn", "catalog_offering.0.version_crn", "boot_volume.0.volume_id"},
							ConflictsWith: []string{isInstanceImage, isInstanceSourceTemplate, "catalog_offering.0.offering_crn", "catalog_offering.0.version_crn", "boot_volume.0.volume_id", "boot_volume.0.snapshot_crn"},
							Optional:      true,
							ForceNew:      true,
							Computed:      true,
						},
						isInstanceVolumeSnapshotCrn: {
							Type:          schema.TypeString,
							RequiredWith:  []string{isInstanceZone, isInstanceProfile, isInstanceVPC},
							AtLeastOneOf:  []string{isInstanceImage, isInstanceSourceTemplate, "boot_volume.0.snapshot", "boot_volume.0.snapshot_crn", "catalog_offering.0.offering_crn", "catalog_offering.0.version_crn", "boot_volume.0.volume_id"},
							ConflictsWith: []string{isInstanceImage, isInstanceSourceTemplate, "catalog_offering.0.offering_crn", "catalog_offering.0.version_crn", "boot_volume.0.volume_id", "boot_volume.0.snapshot"},
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
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							// ValidateFunc: validate.InvokeValidator("ibm_is_instance", isInstanceBootSize),
						},
						isInstanceBootIOPS: {
							Type:     schema.TypeInt,
							Computed: true,
							Optional: true,
						},
						isInstanceBootProfile: {
							Type:     schema.TypeString,
							Optional: true,
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
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"volume_prototypes"},
				Elem:          &schema.Schema{Type: schema.TypeString},
				Description:   "List of volumes",
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

			"numa_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of NUMA nodes this virtual server instance is provisioned on. This property may be absent if the instance's `status` is not `running`.",
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
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The reservation affinity policy to use for this virtual server instance.",
						},
						isReservationAffinityPool: &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "The pool of reservations available for use by this virtual server instance.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isReservationId: {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The unique identifier for this reservation.",
									},
									isReservationCrn: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this reservation.",
									},
									isReservationHref: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this reservation.",
									},
									isReservationName: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this reservation. The name is unique across all reservations in the region.",
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
	}
}

func ResourceIBMISInstanceValidator() *validate.ResourceValidator {
	actions := "stop, start, reboot"
	host_failure := "restart, stop"
	metadataServiceProtocol := "https, http"
	validateSchema := make([]validate.ValidateSchema, 0)

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "confidential_compute_mode",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "disabled, sgx, tdx",
			Regexp:                     `^[a-z][a-z0-9]*(_[a-z0-9]+)*$`,
			MinValueLength:             1,
			MaxValueLength:             128,
		},
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

func instanceCreateByImage(d *schema.ResourceData, meta interface{}, profile, name, vpcID, zone, image, bootProfile string) error {
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

	// cluster changes
	if clusterNetworkAttachmentOk, ok := d.GetOk("cluster_network_attachments"); ok {
		clusterNetworkAttachmentList := clusterNetworkAttachmentOk.([]interface{})
		if len(clusterNetworkAttachmentList) > 0 {
			clusterNetworkAttachments := []vpcv1.InstanceClusterNetworkAttachmentPrototypeInstanceContext{}
			for _, clusterNetworkAttachmentsItem := range clusterNetworkAttachmentList {
				clusterNetworkAttachmentsItemModel, err := ResourceIBMIsInstanceMapToInstanceClusterNetworkAttachmentPrototypeInstanceContext(clusterNetworkAttachmentsItem.(map[string]interface{}))
				if err != nil {
					return err
				}
				clusterNetworkAttachments = append(clusterNetworkAttachments, *clusterNetworkAttachmentsItemModel)
			}
			instanceproto.ClusterNetworkAttachments = clusterNetworkAttachments
		}
	}

	// volume_prototypes
	if volumeattintf, ok := d.GetOk("volume_prototypes"); ok {
		volumeatt := []vpcv1.VolumeAttachmentPrototype{}
		for i, _ := range volumeattintf.([]interface{}) {
			volumeattItemModel := &vpcv1.VolumeAttachmentPrototype{}
			volumeattItemPrototypeModel := &vpcv1.VolumeAttachmentPrototypeVolume{}
			if attNameOk, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.name", i)); ok {
				attName := attNameOk.(string)
				if attName != "" {
					volumeattItemModel.Name = &attName
				}
			}
			if vname, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.volume_name", i)); ok {
				volName := vname.(string)
				if volName != "" {
					volumeattItemPrototypeModel.Name = &volName
				}
			}
			if volAutoDelete, ok := d.GetOkExists(fmt.Sprintf("volume_prototypes.%d.delete_volume_on_instance_delete", i)); ok {
				volumeattItemModel.DeleteVolumeOnInstanceDelete = core.BoolPtr(volAutoDelete.(bool))
			}
			if volIops, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.volume_iops", i)); ok {
				if volIops.(int) != 0 {
					volumeattItemPrototypeModel.Iops = core.Int64Ptr(int64(volIops.(int)))
				}
			}
			if volCapacity, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.volume_capacity", i)); ok {
				if volCapacity != 0 {
					volumeattItemPrototypeModel.Capacity = core.Int64Ptr(int64(volCapacity.(int)))
				}
			}
			if volEncKeyOk, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.volume_encryption_key", i)); ok {
				volEncKey := volEncKeyOk.(string)
				if volEncKey != "" {
					volumeattItemPrototypeModel.EncryptionKey = &vpcv1.EncryptionKeyIdentity{
						CRN: &volEncKey,
					}
				}
			}
			if volProfileOk, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.volume_profile", i)); ok {
				volProfile := volProfileOk.(string)
				if volProfile != "" {
					volumeattItemPrototypeModel.Profile = &vpcv1.VolumeProfileIdentity{
						Name: &volProfile,
					}
				}
			}
			if volRgOk, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.volume_resource_group", i)); ok {
				volRg := volRgOk.(string)
				if volRg != "" {
					volumeattItemPrototypeModel.ResourceGroup = &vpcv1.ResourceGroupIdentity{
						ID: &volRg,
					}
				}
			}
			if volSnapshotok, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.volume_source_snapshot", i)); ok {
				volSnapshot := volSnapshotok.(string)
				if volSnapshot != "" {
					volumeattItemPrototypeModel.SourceSnapshot = &vpcv1.SnapshotIdentity{
						ID: &volSnapshot,
					}
				}
			}
			volTags := d.Get(fmt.Sprintf("volume_prototypes.%d.volume_tags", i)).(*schema.Set)
			if volTags != nil && volTags.Len() != 0 {
				userTagsArray := make([]string, volTags.Len())
				for i, userTag := range volTags.List() {
					userTagStr := userTag.(string)
					userTagsArray[i] = userTagStr
				}
				volumeattItemPrototypeModel.UserTags = userTagsArray
			}

			volumeattItemModel.Volume = volumeattItemPrototypeModel

			volumeatt = append(volumeatt, *volumeattItemModel)
		}
		instanceproto.VolumeAttachments = volumeatt
	}
	if _, ok := d.GetOk("confidential_compute_mode"); ok {
		instanceproto.ConfidentialComputeMode = core.StringPtr(d.Get("confidential_compute_mode").(string))
	}
	if _, ok := d.GetOkExists("enable_secure_boot"); ok {
		instanceproto.EnableSecureBoot = core.BoolPtr(d.Get("enable_secure_boot").(bool))
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
		iopsOk, ok := bootvol[isInstanceBootIOPS]
		iops := iopsOk.(int)
		if iops != 0 && ok {
			iopsInt64 := int64(iops)
			volTemplate.Iops = &iopsInt64
		}
		enc, ok := bootvol[isInstanceBootEncryption]
		encstr := enc.(string)
		if ok && encstr != "" {
			volTemplate.EncryptionKey = &vpcv1.EncryptionKeyIdentity{
				CRN: &encstr,
			}
		}
		if bootProfile == "" {
			bootProfile = "general-purpose"
		}
		volTemplate.Profile = &vpcv1.VolumeProfileIdentity{
			Name: &bootProfile,
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

	if networkattachmentsintf, ok := d.GetOk("network_attachments"); ok {
		networkAttachments := []vpcv1.InstanceNetworkAttachmentPrototype{}
		for i, networkAttachmentsItem := range networkattachmentsintf.([]interface{}) {
			allowipspoofing := fmt.Sprintf("network_attachments.%d.virtual_network_interface.0.allow_ip_spoofing", i)
			autodelete := fmt.Sprintf("network_attachments.%d.virtual_network_interface.0.auto_delete", i)
			enablenat := fmt.Sprintf("network_attachments.%d.virtual_network_interface.0.enable_infrastructure_nat", i)
			// allowipspoofing := "primary_network_attachment.0.allow_ip_spoofing"
			// autodelete := "primary_network_attachment.0.autodelete"
			// enablenat := "primary_network_attachment.0.enable_infrastructure_nat"
			networkAttachmentsItemModel, err := resourceIBMIsInstanceMapToInstanceNetworkAttachmentPrototype(allowipspoofing, autodelete, enablenat, d, networkAttachmentsItem.(map[string]interface{}))
			if err != nil {
				return err
			}
			networkAttachments = append(networkAttachments, *networkAttachmentsItemModel)
		}
		instanceproto.NetworkAttachments = networkAttachments
	}
	if primnetworkattachmentintf, ok := d.GetOk("primary_network_attachment"); ok && len(primnetworkattachmentintf.([]interface{})) > 0 {
		i := 0
		allowipspoofing := fmt.Sprintf("primary_network_attachment.%d.virtual_network_interface.0.allow_ip_spoofing", i)
		autodelete := fmt.Sprintf("primary_network_attachment.%d.virtual_network_interface.0.auto_delete", i)
		enablenat := fmt.Sprintf("primary_network_attachment.%d.virtual_network_interface.0.enable_infrastructure_nat", i)
		primaryNetworkAttachmentModel, err := resourceIBMIsInstanceMapToInstanceNetworkAttachmentPrototype(allowipspoofing, autodelete, enablenat, d, primnetworkattachmentintf.([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return err
		}
		instanceproto.PrimaryNetworkAttachment = primaryNetworkAttachmentModel
	}
	if resAffinity, ok := d.GetOk(isReservationAffinity); ok {
		resAff := resAffinity.([]interface{})[0].(map[string]interface{})
		var resAffinity = &vpcv1.InstanceReservationAffinityPrototype{}
		policy, ok := resAff["policy"]
		policyStr := policy.(string)
		if policyStr != "" && ok {
			resAffinity.Policy = &policyStr
		}
		poolIntf, okPool := resAff[isReservationAffinityPool]
		if okPool && poolIntf != nil && poolIntf.([]interface{}) != nil && len(poolIntf.([]interface{})) > 0 {
			pool := poolIntf.([]interface{})[0].(map[string]interface{})
			id, okId := pool["id"]
			if okId {
				idStr, ok := id.(string)
				if idStr != "" && ok {
					var resAffPool = make([]vpcv1.ReservationIdentityIntf, 1)
					resAffPool[0] = &vpcv1.ReservationIdentity{
						ID: &idStr,
					}
					resAffinity.Pool = resAffPool
				}
			}
		}
		instanceproto.ReservationAffinity = resAffinity
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

	if keySetIntf, ok := d.GetOk(isInstanceKeys); ok {
		keySet := keySetIntf.(*schema.Set)
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
				"[ERROR] Error on create of resource instance (%s) tags: %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk(isInstanceAccessTags); ok {
		oldList, newList := d.GetChange(isInstanceAccessTags)
		err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *instance.CRN, "", isInstanceAccessTagType)
		if err != nil {
			log.Printf(
				"[ERROR] Error on create of resource instance (%s) access tags: %s", d.Id(), err)
		}
	}
	return nil
}
func instanceCreateByCatalogOffering(d *schema.ResourceData, meta interface{}, profile, name, vpcID, zone, image, offerringCrn, versionCrn, planCrn string) error {
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
	// cluster changes
	if clusterNetworkAttachmentOk, ok := d.GetOk("cluster_network_attachments"); ok {
		clusterNetworkAttachmentList := clusterNetworkAttachmentOk.([]interface{})
		if len(clusterNetworkAttachmentList) > 0 {
			clusterNetworkAttachments := []vpcv1.InstanceClusterNetworkAttachmentPrototypeInstanceContext{}
			for _, clusterNetworkAttachmentsItem := range clusterNetworkAttachmentList {
				clusterNetworkAttachmentsItemModel, err := ResourceIBMIsInstanceMapToInstanceClusterNetworkAttachmentPrototypeInstanceContext(clusterNetworkAttachmentsItem.(map[string]interface{}))
				if err != nil {
					return err
				}
				clusterNetworkAttachments = append(clusterNetworkAttachments, *clusterNetworkAttachmentsItemModel)
			}
			instanceproto.ClusterNetworkAttachments = clusterNetworkAttachments
		}
	}
	// volume_prototypes
	if volumeattintf, ok := d.GetOk("volume_prototypes"); ok {
		volumeatt := []vpcv1.VolumeAttachmentPrototype{}
		for i, _ := range volumeattintf.([]interface{}) {
			volumeattItemModel := &vpcv1.VolumeAttachmentPrototype{}
			volumeattItemPrototypeModel := &vpcv1.VolumeAttachmentPrototypeVolume{}
			if attNameOk, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.name", i)); ok {
				attName := attNameOk.(string)
				if attName != "" {
					volumeattItemModel.Name = &attName
				}
			}
			if vname, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.volume_name", i)); ok {
				volName := vname.(string)
				if volName != "" {
					volumeattItemPrototypeModel.Name = &volName
				}
			}
			if volAutoDelete, ok := d.GetOkExists(fmt.Sprintf("volume_prototypes.%d.delete_volume_on_instance_delete", i)); ok {
				volumeattItemModel.DeleteVolumeOnInstanceDelete = core.BoolPtr(volAutoDelete.(bool))
			}
			if volIops, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.volume_iops", i)); ok {
				if volIops.(int) != 0 {
					volumeattItemPrototypeModel.Iops = core.Int64Ptr(int64(volIops.(int)))
				}
			}
			if volCapacity, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.volume_capacity", i)); ok {
				if volCapacity != 0 {
					volumeattItemPrototypeModel.Capacity = core.Int64Ptr(int64(volCapacity.(int)))
				}
			}
			if volEncKeyOk, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.volume_encryption_key", i)); ok {
				volEncKey := volEncKeyOk.(string)
				if volEncKey != "" {
					volumeattItemPrototypeModel.EncryptionKey = &vpcv1.EncryptionKeyIdentity{
						CRN: &volEncKey,
					}
				}
			}
			if volProfileOk, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.volume_profile", i)); ok {
				volProfile := volProfileOk.(string)
				if volProfile != "" {
					volumeattItemPrototypeModel.Profile = &vpcv1.VolumeProfileIdentity{
						Name: &volProfile,
					}
				}
			}
			if volRgOk, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.volume_resource_group", i)); ok {
				volRg := volRgOk.(string)
				if volRg != "" {
					volumeattItemPrototypeModel.ResourceGroup = &vpcv1.ResourceGroupIdentity{
						ID: &volRg,
					}
				}
			}
			if volSnapshotok, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.volume_source_snapshot", i)); ok {
				volSnapshot := volSnapshotok.(string)
				if volSnapshot != "" {
					volumeattItemPrototypeModel.SourceSnapshot = &vpcv1.SnapshotIdentity{
						ID: &volSnapshot,
					}
				}
			}
			volTags := d.Get(fmt.Sprintf("volume_prototypes.%d.volume_tags", i)).(*schema.Set)
			if volTags != nil && volTags.Len() != 0 {
				userTagsArray := make([]string, volTags.Len())
				for i, userTag := range volTags.List() {
					userTagStr := userTag.(string)
					userTagsArray[i] = userTagStr
				}
				volumeattItemPrototypeModel.UserTags = userTagsArray
			}

			volumeattItemModel.Volume = volumeattItemPrototypeModel

			volumeatt = append(volumeatt, *volumeattItemModel)
		}
		instanceproto.VolumeAttachments = volumeatt
	}
	if _, ok := d.GetOk("confidential_compute_mode"); ok {
		instanceproto.ConfidentialComputeMode = core.StringPtr(d.Get("confidential_compute_mode").(string))
	}
	if _, ok := d.GetOkExists("enable_secure_boot"); ok {
		instanceproto.EnableSecureBoot = core.BoolPtr(d.Get("enable_secure_boot").(bool))
	}
	var planOffering *vpcv1.CatalogOfferingVersionPlanIdentityCatalogOfferingVersionPlanByCRN
	planOffering = nil
	if planCrn != "" {
		planOffering = &vpcv1.CatalogOfferingVersionPlanIdentityCatalogOfferingVersionPlanByCRN{
			CRN: &planCrn,
		}
	}

	if offerringCrn != "" {
		catalogOffering := &vpcv1.CatalogOfferingIdentityCatalogOfferingByCRN{
			CRN: &offerringCrn,
		}
		offeringPrototype := &vpcv1.InstanceCatalogOfferingPrototypeCatalogOfferingByOffering{
			Offering: catalogOffering,
		}
		if planOffering != nil {
			offeringPrototype.Plan = planOffering
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
		if planOffering != nil {
			versionPrototype.Plan = planOffering
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
		iopsOk, ok := bootvol[isInstanceBootIOPS]
		iops := iopsOk.(int)
		if iops != 0 && ok {
			iopsInt64 := int64(iops)
			volTemplate.Iops = &iopsInt64
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

	if networkattachmentsintf, ok := d.GetOk("network_attachments"); ok {
		networkAttachments := []vpcv1.InstanceNetworkAttachmentPrototype{}
		for i, networkAttachmentsItem := range networkattachmentsintf.([]interface{}) {
			allowipspoofing := fmt.Sprintf("network_attachments.%d.virtual_network_interface.0.allow_ip_spoofing", i)
			autodelete := fmt.Sprintf("network_attachments.%d.virtual_network_interface.0.auto_delete", i)
			enablenat := fmt.Sprintf("network_attachments.%d.virtual_network_interface.0.enable_infrastructure_nat", i)
			// allowipspoofing := "primary_network_attachment.0.allow_ip_spoofing"
			// autodelete := "primary_network_attachment.0.autodelete"
			// enablenat := "primary_network_attachment.0.enable_infrastructure_nat"
			networkAttachmentsItemModel, err := resourceIBMIsInstanceMapToInstanceNetworkAttachmentPrototype(allowipspoofing, autodelete, enablenat, d, networkAttachmentsItem.(map[string]interface{}))
			if err != nil {
				return err
			}
			networkAttachments = append(networkAttachments, *networkAttachmentsItemModel)
		}
		instanceproto.NetworkAttachments = networkAttachments
	}
	if primnetworkattachmentintf, ok := d.GetOk("primary_network_attachment"); ok && len(primnetworkattachmentintf.([]interface{})) > 0 {
		i := 0
		allowipspoofing := fmt.Sprintf("primary_network_attachment.%d.virtual_network_interface.0.allow_ip_spoofing", i)
		autodelete := fmt.Sprintf("primary_network_attachment.%d.virtual_network_interface.0.auto_delete", i)
		enablenat := fmt.Sprintf("primary_network_attachment.%d.virtual_network_interface.0.enable_infrastructure_nat", i)
		primaryNetworkAttachmentModel, err := resourceIBMIsInstanceMapToInstanceNetworkAttachmentPrototype(allowipspoofing, autodelete, enablenat, d, primnetworkattachmentintf.([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return err
		}
		instanceproto.PrimaryNetworkAttachment = primaryNetworkAttachmentModel
	}
	if resAffinity, ok := d.GetOk(isReservationAffinity); ok {
		resAff := resAffinity.([]interface{})[0].(map[string]interface{})
		var resAffinity = &vpcv1.InstanceReservationAffinityPrototype{}
		policy, ok := resAff["policy"]
		policyStr := policy.(string)
		if policyStr != "" && ok {
			resAffinity.Policy = &policyStr
		}
		poolIntf, okPool := resAff[isReservationAffinityPool]
		if okPool && poolIntf != nil && poolIntf.([]interface{}) != nil && len(poolIntf.([]interface{})) > 0 {
			pool := poolIntf.([]interface{})[0].(map[string]interface{})
			id, okId := pool["id"]
			if okId {
				idStr, ok := id.(string)
				if idStr != "" && ok {
					var resAffPool = make([]vpcv1.ReservationIdentityIntf, 1)
					resAffPool[0] = &vpcv1.ReservationIdentity{
						ID: &idStr,
					}
					resAffinity.Pool = resAffPool
				}
			}
		}
		instanceproto.ReservationAffinity = resAffinity
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

	if keySetIntf, ok := d.GetOk(isInstanceKeys); ok {
		keySet := keySetIntf.(*schema.Set)
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
				"[ERROR] Error on create of resource instance (%s) tags: %s", d.Id(), err)
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
	// cluster changes
	if clusterNetworkAttachmentOk, ok := d.GetOk("cluster_network_attachments"); ok {
		clusterNetworkAttachmentList := clusterNetworkAttachmentOk.([]interface{})
		if len(clusterNetworkAttachmentList) > 0 {
			clusterNetworkAttachments := []vpcv1.InstanceClusterNetworkAttachmentPrototypeInstanceContext{}
			for _, clusterNetworkAttachmentsItem := range clusterNetworkAttachmentList {
				clusterNetworkAttachmentsItemModel, err := ResourceIBMIsInstanceMapToInstanceClusterNetworkAttachmentPrototypeInstanceContext(clusterNetworkAttachmentsItem.(map[string]interface{}))
				if err != nil {
					return err
				}
				clusterNetworkAttachments = append(clusterNetworkAttachments, *clusterNetworkAttachmentsItemModel)
			}
			instanceproto.ClusterNetworkAttachments = clusterNetworkAttachments
		}
	}
	// volume_prototypes
	if volumeattintf, ok := d.GetOk("volume_prototypes"); ok {
		volumeatt := []vpcv1.VolumeAttachmentPrototype{}
		for i, _ := range volumeattintf.([]interface{}) {
			volumeattItemModel := &vpcv1.VolumeAttachmentPrototype{}
			volumeattItemPrototypeModel := &vpcv1.VolumeAttachmentPrototypeVolume{}
			if attNameOk, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.name", i)); ok {
				attName := attNameOk.(string)
				if attName != "" {
					volumeattItemModel.Name = &attName
				}
			}
			if vname, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.volume_name", i)); ok {
				volName := vname.(string)
				if volName != "" {
					volumeattItemPrototypeModel.Name = &volName
				}
			}
			if volAutoDelete, ok := d.GetOkExists(fmt.Sprintf("volume_prototypes.%d.delete_volume_on_instance_delete", i)); ok {
				volumeattItemModel.DeleteVolumeOnInstanceDelete = core.BoolPtr(volAutoDelete.(bool))
			}
			if volIops, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.volume_iops", i)); ok {
				if volIops.(int) != 0 {
					volumeattItemPrototypeModel.Iops = core.Int64Ptr(int64(volIops.(int)))
				}
			}
			if volCapacity, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.volume_capacity", i)); ok {
				if volCapacity != 0 {
					volumeattItemPrototypeModel.Capacity = core.Int64Ptr(int64(volCapacity.(int)))
				}
			}
			if volEncKeyOk, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.volume_encryption_key", i)); ok {
				volEncKey := volEncKeyOk.(string)
				if volEncKey != "" {
					volumeattItemPrototypeModel.EncryptionKey = &vpcv1.EncryptionKeyIdentity{
						CRN: &volEncKey,
					}
				}
			}
			if volProfileOk, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.volume_profile", i)); ok {
				volProfile := volProfileOk.(string)
				if volProfile != "" {
					volumeattItemPrototypeModel.Profile = &vpcv1.VolumeProfileIdentity{
						Name: &volProfile,
					}
				}
			}
			if volRgOk, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.volume_resource_group", i)); ok {
				volRg := volRgOk.(string)
				if volRg != "" {
					volumeattItemPrototypeModel.ResourceGroup = &vpcv1.ResourceGroupIdentity{
						ID: &volRg,
					}
				}
			}
			if volSnapshotok, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.volume_source_snapshot", i)); ok {
				volSnapshot := volSnapshotok.(string)
				if volSnapshot != "" {
					volumeattItemPrototypeModel.SourceSnapshot = &vpcv1.SnapshotIdentity{
						ID: &volSnapshot,
					}
				}
			}
			volTags := d.Get(fmt.Sprintf("volume_prototypes.%d.volume_tags", i)).(*schema.Set)
			if volTags != nil && volTags.Len() != 0 {
				userTagsArray := make([]string, volTags.Len())
				for i, userTag := range volTags.List() {
					userTagStr := userTag.(string)
					userTagsArray[i] = userTagStr
				}
				volumeattItemPrototypeModel.UserTags = userTagsArray
			}

			volumeattItemModel.Volume = volumeattItemPrototypeModel

			volumeatt = append(volumeatt, *volumeattItemModel)
		}
		instanceproto.VolumeAttachments = volumeatt
	}
	if _, ok := d.GetOk("confidential_compute_mode"); ok {
		instanceproto.ConfidentialComputeMode = core.StringPtr(d.Get("confidential_compute_mode").(string))
	}
	if _, ok := d.GetOkExists("enable_secure_boot"); ok {
		instanceproto.EnableSecureBoot = core.BoolPtr(d.Get("enable_secure_boot").(bool))
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
		iopsOk, ok := bootvol[isInstanceBootIOPS]
		iops := iopsOk.(int)
		if iops != 0 && ok {
			iopsInt64 := int64(iops)
			volTemplate.Iops = &iopsInt64
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

	if networkattachmentsintf, ok := d.GetOk("network_attachments"); ok {
		networkAttachments := []vpcv1.InstanceNetworkAttachmentPrototype{}
		for i, networkAttachmentsItem := range networkattachmentsintf.([]interface{}) {
			allowipspoofing := fmt.Sprintf("network_attachments.%d.allow_ip_spoofing", i)
			autodelete := fmt.Sprintf("network_attachments.%d.autodelete", i)
			enablenat := fmt.Sprintf("network_attachments.%d.enable_infrastructure_nat", i)
			networkAttachmentsItemModel, err := resourceIBMIsInstanceMapToInstanceNetworkAttachmentPrototype(allowipspoofing, autodelete, enablenat, d, networkAttachmentsItem.(map[string]interface{}))
			if err != nil {
				return err
			}
			networkAttachments = append(networkAttachments, *networkAttachmentsItemModel)
		}
		instanceproto.NetworkAttachments = networkAttachments
	}
	if primnetworkattachmentintf, ok := d.GetOk("primary_network_attachment"); ok && len(primnetworkattachmentintf.([]interface{})) > 0 {
		i := 0
		allowipspoofing := fmt.Sprintf("primary_network_attachment.%d.virtual_network_interface.0.allow_ip_spoofing", i)
		autodelete := fmt.Sprintf("primary_network_attachment.%d.virtual_network_interface.0.auto_delete", i)
		enablenat := fmt.Sprintf("primary_network_attachment.%d.virtual_network_interface.0.enable_infrastructure_nat", i)
		primaryNetworkAttachmentModel, err := resourceIBMIsInstanceMapToInstanceNetworkAttachmentPrototype(allowipspoofing, autodelete, enablenat, d, primnetworkattachmentintf.([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return err
		}
		instanceproto.PrimaryNetworkAttachment = primaryNetworkAttachmentModel
	}
	if resAffinity, ok := d.GetOk(isReservationAffinity); ok {
		resAff := resAffinity.([]interface{})[0].(map[string]interface{})
		var resAffinity = &vpcv1.InstanceReservationAffinityPrototype{}
		policy, ok := resAff["policy"]
		policyStr := policy.(string)
		if policyStr != "" && ok {
			resAffinity.Policy = &policyStr
		}
		poolIntf, okPool := resAff[isReservationAffinityPool]
		if okPool && poolIntf != nil && poolIntf.([]interface{}) != nil && len(poolIntf.([]interface{})) > 0 {
			pool := poolIntf.([]interface{})[0].(map[string]interface{})
			id, okId := pool["id"]
			if okId {
				idStr, ok := id.(string)
				if idStr != "" && ok {
					var resAffPool = make([]vpcv1.ReservationIdentityIntf, 1)
					resAffPool[0] = &vpcv1.ReservationIdentity{
						ID: &idStr,
					}
					resAffinity.Pool = resAffPool
				}
			}
		}
		instanceproto.ReservationAffinity = resAffinity
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

	if keySetIntf, ok := d.GetOk(isInstanceKeys); ok {
		keySet := keySetIntf.(*schema.Set)
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
				"[ERROR] Error on create of resource instance (%s) tags: %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk(isInstanceAccessTags); ok {
		oldList, newList := d.GetChange(isInstanceAccessTags)
		err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *instance.CRN, "", isInstanceAccessTagType)
		if err != nil {
			log.Printf(
				"[ERROR] Error on create of resource instance (%s) access tags: %s", d.Id(), err)
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
	// cluster changes
	if clusterNetworkAttachmentOk, ok := d.GetOk("cluster_network_attachments"); ok {
		clusterNetworkAttachmentList := clusterNetworkAttachmentOk.([]interface{})
		if len(clusterNetworkAttachmentList) > 0 {
			clusterNetworkAttachments := []vpcv1.InstanceClusterNetworkAttachmentPrototypeInstanceContext{}
			for _, clusterNetworkAttachmentsItem := range clusterNetworkAttachmentList {
				clusterNetworkAttachmentsItemModel, err := ResourceIBMIsInstanceMapToInstanceClusterNetworkAttachmentPrototypeInstanceContext(clusterNetworkAttachmentsItem.(map[string]interface{}))
				if err != nil {
					return err
				}
				clusterNetworkAttachments = append(clusterNetworkAttachments, *clusterNetworkAttachmentsItemModel)
			}
			instanceproto.ClusterNetworkAttachments = clusterNetworkAttachments
		}
	}
	// volume_prototypes
	if volumeattintf, ok := d.GetOk("volume_prototypes"); ok {
		volumeatt := []vpcv1.VolumeAttachmentPrototype{}
		for i, _ := range volumeattintf.([]interface{}) {
			volumeattItemModel := &vpcv1.VolumeAttachmentPrototype{}
			volumeattItemPrototypeModel := &vpcv1.VolumeAttachmentPrototypeVolume{}
			if attNameOk, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.name", i)); ok {
				attName := attNameOk.(string)
				if attName != "" {
					volumeattItemModel.Name = &attName
				}
			}
			if vname, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.volume_name", i)); ok {
				volName := vname.(string)
				if volName != "" {
					volumeattItemPrototypeModel.Name = &volName
				}
			}
			if volAutoDelete, ok := d.GetOkExists(fmt.Sprintf("volume_prototypes.%d.delete_volume_on_instance_delete", i)); ok {
				volumeattItemModel.DeleteVolumeOnInstanceDelete = core.BoolPtr(volAutoDelete.(bool))
			}
			if volIops, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.volume_iops", i)); ok {
				if volIops.(int) != 0 {
					volumeattItemPrototypeModel.Iops = core.Int64Ptr(int64(volIops.(int)))
				}
			}
			if volCapacity, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.volume_capacity", i)); ok {
				if volCapacity != 0 {
					volumeattItemPrototypeModel.Capacity = core.Int64Ptr(int64(volCapacity.(int)))
				}
			}
			if volEncKeyOk, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.volume_encryption_key", i)); ok {
				volEncKey := volEncKeyOk.(string)
				if volEncKey != "" {
					volumeattItemPrototypeModel.EncryptionKey = &vpcv1.EncryptionKeyIdentity{
						CRN: &volEncKey,
					}
				}
			}
			if volProfileOk, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.volume_profile", i)); ok {
				volProfile := volProfileOk.(string)
				if volProfile != "" {
					volumeattItemPrototypeModel.Profile = &vpcv1.VolumeProfileIdentity{
						Name: &volProfile,
					}
				}
			}
			if volRgOk, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.volume_resource_group", i)); ok {
				volRg := volRgOk.(string)
				if volRg != "" {
					volumeattItemPrototypeModel.ResourceGroup = &vpcv1.ResourceGroupIdentity{
						ID: &volRg,
					}
				}
			}
			if volSnapshotok, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.volume_source_snapshot", i)); ok {
				volSnapshot := volSnapshotok.(string)
				if volSnapshot != "" {
					volumeattItemPrototypeModel.SourceSnapshot = &vpcv1.SnapshotIdentity{
						ID: &volSnapshot,
					}
				}
			}
			volTags := d.Get(fmt.Sprintf("volume_prototypes.%d.volume_tags", i)).(*schema.Set)
			if volTags != nil && volTags.Len() != 0 {
				userTagsArray := make([]string, volTags.Len())
				for i, userTag := range volTags.List() {
					userTagStr := userTag.(string)
					userTagsArray[i] = userTagStr
				}
				volumeattItemPrototypeModel.UserTags = userTagsArray
			}

			volumeattItemModel.Volume = volumeattItemPrototypeModel

			volumeatt = append(volumeatt, *volumeattItemModel)
		}
		instanceproto.VolumeAttachments = volumeatt
	}
	if _, ok := d.GetOk("confidential_compute_mode"); ok {
		instanceproto.ConfidentialComputeMode = core.StringPtr(d.Get("confidential_compute_mode").(string))
	}
	if _, ok := d.GetOkExists("enable_secure_boot"); ok {
		instanceproto.EnableSecureBoot = core.BoolPtr(d.Get("enable_secure_boot").(bool))
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
		iopsOk, ok := bootvol[isInstanceBootIOPS]
		iops := iopsOk.(int)
		if iops != 0 && ok {
			iopsInt64 := int64(iops)
			volTemplate.Iops = &iopsInt64
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
		snapshotCrn, ok := bootvol[isInstanceVolumeSnapshotCrn]
		snapshotCrnStr := snapshotCrn.(string)
		if snapshotCrnStr != "" && ok {
			volTemplate.SourceSnapshot = &vpcv1.SnapshotIdentity{
				CRN: &snapshotCrnStr,
			}
		}
		deleteboolIntf := bootvol[isInstanceVolAttVolAutoDelete]
		deletebool := deleteboolIntf.(bool)
		instanceproto.BootVolumeAttachment = &vpcv1.VolumeAttachmentPrototypeInstanceBySourceSnapshotContext{
			DeleteVolumeOnInstanceDelete: &deletebool,
			Volume:                       volTemplate,
		}
	}

	if resAffinity, ok := d.GetOk(isReservationAffinity); ok {
		resAff := resAffinity.([]interface{})[0].(map[string]interface{})
		var resAffinity = &vpcv1.InstanceReservationAffinityPrototype{}
		policy, ok := resAff["policy"]
		policyStr := policy.(string)
		if policyStr != "" && ok {
			resAffinity.Policy = &policyStr
		}
		poolIntf, okPool := resAff[isReservationAffinityPool]
		if okPool && poolIntf != nil && poolIntf.([]interface{}) != nil && len(poolIntf.([]interface{})) > 0 {
			pool := poolIntf.([]interface{})[0].(map[string]interface{})
			id, okId := pool["id"]
			if okId {
				idStr, ok := id.(string)
				if idStr != "" && ok {
					var resAffPool = make([]vpcv1.ReservationIdentityIntf, 1)
					resAffPool[0] = &vpcv1.ReservationIdentity{
						ID: &idStr,
					}
					resAffinity.Pool = resAffPool
				}
			}
		}
		instanceproto.ReservationAffinity = resAffinity
	}

	if totalVolBandwidthIntf, ok := d.GetOk(isInstanceTotalVolumeBandwidth); ok {
		totalVolBandwidthStr := int64(totalVolBandwidthIntf.(int))
		instanceproto.TotalVolumeBandwidth = &totalVolBandwidthStr
	}

	if networkattachmentsintf, ok := d.GetOk("network_attachments"); ok {
		networkAttachments := []vpcv1.InstanceNetworkAttachmentPrototype{}
		for i, networkAttachmentsItem := range networkattachmentsintf.([]interface{}) {
			allowipspoofing := fmt.Sprintf("network_attachments.%d.allow_ip_spoofing", i)
			autodelete := fmt.Sprintf("network_attachments.%d.autodelete", i)
			enablenat := fmt.Sprintf("network_attachments.%d.enable_infrastructure_nat", i)
			networkAttachmentsItemModel, err := resourceIBMIsInstanceMapToInstanceNetworkAttachmentPrototype(allowipspoofing, autodelete, enablenat, d, networkAttachmentsItem.(map[string]interface{}))
			if err != nil {
				return err
			}
			networkAttachments = append(networkAttachments, *networkAttachmentsItemModel)
		}
		instanceproto.NetworkAttachments = networkAttachments
	}
	if primnetworkattachmentintf, ok := d.GetOk("primary_network_attachment"); ok && len(primnetworkattachmentintf.([]interface{})) > 0 {
		i := 0
		allowipspoofing := fmt.Sprintf("primary_network_attachment.%d.virtual_network_interface.0.allow_ip_spoofing", i)
		autodelete := fmt.Sprintf("primary_network_attachment.%d.virtual_network_interface.0.auto_delete", i)
		enablenat := fmt.Sprintf("primary_network_attachment.%d.virtual_network_interface.0.enable_infrastructure_nat", i)
		primaryNetworkAttachmentModel, err := resourceIBMIsInstanceMapToInstanceNetworkAttachmentPrototype(allowipspoofing, autodelete, enablenat, d, primnetworkattachmentintf.([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return err
		}
		instanceproto.PrimaryNetworkAttachment = primaryNetworkAttachmentModel
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

	if keySetIntf, ok := d.GetOk(isInstanceKeys); ok {
		keySet := keySetIntf.(*schema.Set)
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
				"[ERROR] Error on create of resource instance (%s) tags: %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk(isInstanceAccessTags); ok {
		oldList, newList := d.GetChange(isInstanceAccessTags)
		err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *instance.CRN, "", isInstanceAccessTagType)
		if err != nil {
			log.Printf(
				"[ERROR] Error on create of resource instance (%s) access tags: %s", d.Id(), err)
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
	// cluster changes
	if clusterNetworkAttachmentOk, ok := d.GetOk("cluster_network_attachments"); ok {
		clusterNetworkAttachmentList := clusterNetworkAttachmentOk.([]interface{})
		if len(clusterNetworkAttachmentList) > 0 {
			clusterNetworkAttachments := []vpcv1.InstanceClusterNetworkAttachmentPrototypeInstanceContext{}
			for _, clusterNetworkAttachmentsItem := range clusterNetworkAttachmentList {
				clusterNetworkAttachmentsItemModel, err := ResourceIBMIsInstanceMapToInstanceClusterNetworkAttachmentPrototypeInstanceContext(clusterNetworkAttachmentsItem.(map[string]interface{}))
				if err != nil {
					return err
				}
				clusterNetworkAttachments = append(clusterNetworkAttachments, *clusterNetworkAttachmentsItemModel)
			}
			instanceproto.ClusterNetworkAttachments = clusterNetworkAttachments
		}
	}
	// volume_prototypes
	if volumeattintf, ok := d.GetOk("volume_prototypes"); ok {
		volumeatt := []vpcv1.VolumeAttachmentPrototype{}
		for i, _ := range volumeattintf.([]interface{}) {
			volumeattItemModel := &vpcv1.VolumeAttachmentPrototype{}
			volumeattItemPrototypeModel := &vpcv1.VolumeAttachmentPrototypeVolume{}
			if attNameOk, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.name", i)); ok {
				attName := attNameOk.(string)
				if attName != "" {
					volumeattItemModel.Name = &attName
				}
			}
			if vname, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.volume_name", i)); ok {
				volName := vname.(string)
				if volName != "" {
					volumeattItemPrototypeModel.Name = &volName
				}
			}
			if volAutoDelete, ok := d.GetOkExists(fmt.Sprintf("volume_prototypes.%d.delete_volume_on_instance_delete", i)); ok {
				volumeattItemModel.DeleteVolumeOnInstanceDelete = core.BoolPtr(volAutoDelete.(bool))
			}
			if volIops, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.volume_iops", i)); ok {
				if volIops.(int) != 0 {
					volumeattItemPrototypeModel.Iops = core.Int64Ptr(int64(volIops.(int)))
				}
			}
			if volCapacity, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.volume_capacity", i)); ok {
				if volCapacity != 0 {
					volumeattItemPrototypeModel.Capacity = core.Int64Ptr(int64(volCapacity.(int)))
				}
			}
			if volEncKeyOk, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.volume_encryption_key", i)); ok {
				volEncKey := volEncKeyOk.(string)
				if volEncKey != "" {
					volumeattItemPrototypeModel.EncryptionKey = &vpcv1.EncryptionKeyIdentity{
						CRN: &volEncKey,
					}
				}
			}
			if volProfileOk, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.volume_profile", i)); ok {
				volProfile := volProfileOk.(string)
				if volProfile != "" {
					volumeattItemPrototypeModel.Profile = &vpcv1.VolumeProfileIdentity{
						Name: &volProfile,
					}
				}
			}
			if volRgOk, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.volume_resource_group", i)); ok {
				volRg := volRgOk.(string)
				if volRg != "" {
					volumeattItemPrototypeModel.ResourceGroup = &vpcv1.ResourceGroupIdentity{
						ID: &volRg,
					}
				}
			}
			if volSnapshotok, ok := d.GetOk(fmt.Sprintf("volume_prototypes.%d.volume_source_snapshot", i)); ok {
				volSnapshot := volSnapshotok.(string)
				if volSnapshot != "" {
					volumeattItemPrototypeModel.SourceSnapshot = &vpcv1.SnapshotIdentity{
						ID: &volSnapshot,
					}
				}
			}
			volTags := d.Get(fmt.Sprintf("volume_prototypes.%d.volume_tags", i)).(*schema.Set)
			if volTags != nil && volTags.Len() != 0 {
				userTagsArray := make([]string, volTags.Len())
				for i, userTag := range volTags.List() {
					userTagStr := userTag.(string)
					userTagsArray[i] = userTagStr
				}
				volumeattItemPrototypeModel.UserTags = userTagsArray
			}

			volumeattItemModel.Volume = volumeattItemPrototypeModel

			volumeatt = append(volumeatt, *volumeattItemModel)
		}
		instanceproto.VolumeAttachments = volumeatt
	}
	if _, ok := d.GetOk("confidential_compute_mode"); ok {
		instanceproto.ConfidentialComputeMode = core.StringPtr(d.Get("confidential_compute_mode").(string))
	}
	if _, ok := d.GetOkExists("enable_secure_boot"); ok {
		instanceproto.EnableSecureBoot = core.BoolPtr(d.Get("enable_secure_boot").(bool))
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

	if resAffinity, ok := d.GetOk(isReservationAffinity); ok {
		resAff := resAffinity.([]interface{})[0].(map[string]interface{})
		var resAffinity = &vpcv1.InstanceReservationAffinityPrototype{}
		policy, ok := resAff["policy"]
		policyStr := policy.(string)
		if policyStr != "" && ok {
			resAffinity.Policy = &policyStr
		}
		poolIntf, okPool := resAff[isReservationAffinityPool]
		if okPool && poolIntf != nil && poolIntf.([]interface{}) != nil && len(poolIntf.([]interface{})) > 0 {
			pool := poolIntf.([]interface{})[0].(map[string]interface{})
			id, okId := pool["id"]
			if okId {
				idStr, ok := id.(string)
				if idStr != "" && ok {
					var resAffPool = make([]vpcv1.ReservationIdentityIntf, 1)
					resAffPool[0] = &vpcv1.ReservationIdentity{
						ID: &idStr,
					}
					resAffinity.Pool = resAffPool
				}
			}
		}
		instanceproto.ReservationAffinity = resAffinity
	}

	if totalVolBandwidthIntf, ok := d.GetOk(isInstanceTotalVolumeBandwidth); ok {
		totalVolBandwidthStr := int64(totalVolBandwidthIntf.(int))
		instanceproto.TotalVolumeBandwidth = &totalVolBandwidthStr
	}
	if networkattachmentsintf, ok := d.GetOk("network_attachments"); ok {
		networkAttachments := []vpcv1.InstanceNetworkAttachmentPrototype{}
		for i, networkAttachmentsItem := range networkattachmentsintf.([]interface{}) {
			allowipspoofing := fmt.Sprintf("network_attachments.%d.virtual_network_interface.0.allow_ip_spoofing", i)
			autodelete := fmt.Sprintf("network_attachments.%d.virtual_network_interface.0.auto_delete", i)
			enablenat := fmt.Sprintf("network_attachments.%d.virtual_network_interface.0.enable_infrastructure_nat", i)
			// allowipspoofing := "primary_network_attachment.0.allow_ip_spoofing"
			// autodelete := "primary_network_attachment.0.autodelete"
			// enablenat := "primary_network_attachment.0.enable_infrastructure_nat"
			networkAttachmentsItemModel, err := resourceIBMIsInstanceMapToInstanceNetworkAttachmentPrototype(allowipspoofing, autodelete, enablenat, d, networkAttachmentsItem.(map[string]interface{}))
			if err != nil {
				return err
			}
			networkAttachments = append(networkAttachments, *networkAttachmentsItemModel)
		}
		instanceproto.NetworkAttachments = networkAttachments
	}
	if primnetworkattachmentintf, ok := d.GetOk("primary_network_attachment"); ok && len(primnetworkattachmentintf.([]interface{})) > 0 {
		i := 0
		allowipspoofing := fmt.Sprintf("primary_network_attachment.%d.virtual_network_interface.0.allow_ip_spoofing", i)
		autodelete := fmt.Sprintf("primary_network_attachment.%d.virtual_network_interface.0.auto_delete", i)
		enablenat := fmt.Sprintf("primary_network_attachment.%d.virtual_network_interface.0.enable_infrastructure_nat", i)
		primaryNetworkAttachmentModel, err := resourceIBMIsInstanceMapToInstanceNetworkAttachmentPrototype(allowipspoofing, autodelete, enablenat, d, primnetworkattachmentintf.([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return err
		}
		instanceproto.PrimaryNetworkAttachment = primaryNetworkAttachmentModel
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

	if keySetIntf, ok := d.GetOk(isInstanceKeys); ok {
		keySet := keySetIntf.(*schema.Set)
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
				"[ERROR] Error on create of resource instance (%s) tags: %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk(isInstanceAccessTags); ok {
		oldList, newList := d.GetChange(isInstanceAccessTags)
		err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *instance.CRN, "", isInstanceAccessTagType)
		if err != nil {
			log.Printf(
				"[ERROR] Error on create of resource instance (%s) access tags: %s", d.Id(), err)
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
	bootProfile := d.Get("boot_volume.0.profile").(string)
	snapshotcrn := d.Get("boot_volume.0.snapshot_crn").(string)
	volume := d.Get("boot_volume.0.volume_id").(string)
	template := d.Get(isInstanceSourceTemplate).(string)
	if catalogOfferingOk, ok := d.GetOk(isInstanceCatalogOffering); ok {
		catalogOffering := catalogOfferingOk.([]interface{})[0].(map[string]interface{})
		offeringCrn, _ := catalogOffering[isInstanceCatalogOfferingOfferingCrn].(string)
		versionCrn, _ := catalogOffering[isInstanceCatalogOfferingVersionCrn].(string)
		planCrn, _ := catalogOffering[isInstanceCatalogOfferingPlanCrn].(string)
		err := instanceCreateByCatalogOffering(d, meta, profile, name, vpcID, zone, image, offeringCrn, versionCrn, planCrn)
		if err != nil {
			return err
		}

	} else if volume != "" {
		err := instanceCreateByVolume(d, meta, profile, name, vpcID, zone)
		if err != nil {
			return err
		}
	} else if snapshot != "" || snapshotcrn != "" {
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
		err := instanceCreateByImage(d, meta, profile, name, vpcID, zone, image, bootProfile)
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
	// cluster changes
	if !core.IsNil(instance.ClusterNetworkAttachments) {
		clusterNetworkAttachments := []map[string]interface{}{}
		for _, clusterNetworkAttachmentsItem := range instance.ClusterNetworkAttachments {
			clusterNetworkAttachmentsItemMap, err := ResourceIBMIsInstanceInstanceClusterNetworkAttachmentReferenceToMap(instanceC, &clusterNetworkAttachmentsItem, *instance.ID) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance", "read", "cluster_network_attachments-to-map")
			}
			clusterNetworkAttachments = append(clusterNetworkAttachments, clusterNetworkAttachmentsItemMap)
		}
		if err = d.Set("cluster_network_attachments", clusterNetworkAttachments); err != nil {
			err = fmt.Errorf("Error setting cluster_network_attachments: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance", "read", "set-cluster_network_attachments")
		}
	}
	clusterNetwork := make([]map[string]interface{}, 0)
	if !core.IsNil(instance.ClusterNetwork) {
		clusterNetworkMap, err := ResourceIBMIsInstanceClusterNetworkReferenceToMap(instance.ClusterNetwork)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance", "read", "cluster_network-to-map")
		}
		clusterNetwork = append(clusterNetwork, clusterNetworkMap)
	}
	if err = d.Set("cluster_network", clusterNetwork); err != nil {
		err = fmt.Errorf("Error setting cluster_network: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance", "read", "set-cluster_network")
	}
	if !core.IsNil(instance.ConfidentialComputeMode) {
		if err = d.Set("confidential_compute_mode", instance.ConfidentialComputeMode); err != nil {
			return fmt.Errorf("Error setting confidential_compute_mode: %s", err)
		}
	}
	if !core.IsNil(instance.EnableSecureBoot) {
		if err = d.Set("enable_secure_boot", instance.EnableSecureBoot); err != nil {
			return fmt.Errorf("Error setting enable_secure_boot: %s", err)
		}
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

	// volume_prototypes
	volList, _ := setVolumePrototypesInState(d, instance, instanceC)
	d.Set("volume_prototypes", volList)

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
	if instance.HealthReasons != nil {
		healthReasonsList := make([]map[string]interface{}, 0)
		for _, sr := range instance.HealthReasons {
			currentSR := map[string]interface{}{}
			if sr.Code != nil && sr.Message != nil {
				currentSR["code"] = *sr.Code
				currentSR["message"] = *sr.Message
				if sr.MoreInfo != nil {
					currentSR["more_info"] = *sr.Message
				}
				healthReasonsList = append(healthReasonsList, currentSR)
			}
		}
		d.Set("health_reasons", healthReasonsList)
	}
	if err = d.Set("health_state", instance.HealthState); err != nil {
		return err
	}
	if instance.ReservationAffinity != nil {
		reservationAffinity := []map[string]interface{}{}
		reservationAffinityMap := map[string]interface{}{}

		reservationAffinityMap[isReservationAffinityPolicyResp] = instance.ReservationAffinity.Policy
		if instance.ReservationAffinity.Pool != nil && len(instance.ReservationAffinity.Pool) > 0 {
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
		d.Set(isReservationAffinity, reservationAffinity)
	}
	resList := make([]map[string]interface{}, 0)
	if instance.Reservation != nil {
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
	}
	d.Set(isInstanceReservation, resList)

	if !core.IsNil(instance.PrimaryNetworkAttachment) {

		pnaId := *instance.PrimaryNetworkAttachment.ID
		getInstanceNetworkAttachment := &vpcv1.GetInstanceNetworkAttachmentOptions{
			InstanceID: &id,
			ID:         &pnaId,
		}
		autoDelete := true
		if autoDeleteOk, ok := d.GetOkExists("primary_network_attachment.0.virtual_network_interface.0.primary_ip.0.auto_delete"); ok {
			autoDelete = autoDeleteOk.(bool)
		}
		pna, response, err := instanceC.GetInstanceNetworkAttachment(getInstanceNetworkAttachment)
		if err != nil {
			return fmt.Errorf("[ERROR] Error on GetInstanceNetworkAttachment in instance : %s\n%s", err, response)
		}
		primaryNetworkAttachmentMap, err := resourceIBMIsInstanceInstanceNetworkAttachmentReferenceToMap(instance.PrimaryNetworkAttachment, pna, instanceC, autoDelete)
		if err != nil {
			return err
		}
		if err = d.Set("primary_network_attachment", []map[string]interface{}{primaryNetworkAttachmentMap}); err != nil {
			return fmt.Errorf("[ERROR] Error setting primary_network_attachment: %s", err)
		}
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

	if !core.IsNil(instance.NetworkAttachments) {
		networkAttachments := []map[string]interface{}{}
		for i, networkAttachmentsItem := range instance.NetworkAttachments {
			naId := *networkAttachmentsItem.ID
			if *instance.PrimaryNetworkAttachment.ID != naId {
				autoDelete := true
				if autoDeleteOk, ok := d.GetOkExists(fmt.Sprintf("network_attachments.%d.virtual_network_interface.0.primary_ip.0.auto_delete", i)); ok {
					autoDelete = autoDeleteOk.(bool)
				}
				getInstanceNetworkAttachment := &vpcv1.GetInstanceNetworkAttachmentOptions{
					InstanceID: &id,
					ID:         &naId,
				}
				na, response, err := instanceC.GetInstanceNetworkAttachment(getInstanceNetworkAttachment)
				if err != nil {
					return fmt.Errorf("[ERROR] Error on GetInstanceNetworkAttachment in instance : %s\n%s", err, response)
				}
				networkAttachmentsItemMap, err := resourceIBMIsInstanceInstanceNetworkAttachmentReferenceToMap(&networkAttachmentsItem, na, instanceC, autoDelete)
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
	if instance.Image != nil {
		d.Set(isInstanceImage, *instance.Image.ID)
	}
	if instance.NumaCount != nil {
		d.Set("numa_count", int(*instance.NumaCount))
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
				log.Printf("[ERROR] Error Getting Boot Volume (%s): %s\n%s", id, err, response)
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
					bootVol[isInstanceVolumeSnapshotCrn] = vol.SourceSnapshot.CRN
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
			"[ERROR] Error on get of resource Instance (%s) tags: %s", d.Id(), err)
	}
	d.Set(isInstanceTags, tags)
	accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *instance.CRN, "", isInstanceAccessTagType)
	if err != nil {
		log.Printf(
			"[ERROR] Error on get of resource Instance (%s) access tags: %s", d.Id(), err)
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
	// network attachments

	err = handleVolumePrototypesUpdate(d, instanceC)
	if err != nil {
		return err
	}
	err = handleClusterNetworkAttachmentUpdate(d, instanceC)
	if err != nil {
		return err
	}

	if d.HasChange("network_attachments") && !d.IsNewResource() {
		nacs := d.Get("network_attachments").([]interface{})
		ots, nts := d.GetChange("network_attachments")
		otsIntf := ots.([]interface{})
		ntsIntf := nts.([]interface{})
		out := make([]string, len(otsIntf))
		j := 0
		for _, currOtsG := range otsIntf {
			currOts := currOtsG.(map[string]interface{})
			flag := false
			for _, currNtsG := range ntsIntf {
				currNts := currNtsG.(map[string]interface{})
				if currOts["id"].(string) == currNts["id"].(string) {
					flag = true
				}
			}
			if !flag {
				log.Printf("[INFO] Nac with name (%s) will be deleted", currOts["name"].(string))
				nacId := currOts["id"]
				if nacId != nil && nacId.(string) != "" {
					nacIdStr := nacId.(string)
					if !containsNacId(out, nacIdStr) {
						out[j] = nacIdStr
						j = j + 1
						deleteInstanceNetworkAttachmentOptions := &vpcv1.DeleteInstanceNetworkAttachmentOptions{
							InstanceID: &id,
							ID:         &nacIdStr,
						}
						res, err := instanceC.DeleteInstanceNetworkAttachment(deleteInstanceNetworkAttachmentOptions)
						if err != nil {
							return fmt.Errorf("[ERROR] Error while deleting network attachment(%s) of instance(%s) \n%s: %q", nacIdStr, d.Id(), err, res)
						}
					}
				}
			}
		}

		for i, nac := range nacs {
			nacIdKey := fmt.Sprintf("network_attachments.%d.id", i)
			nacId := d.Get(nacIdKey).(string)
			// if nacId is empty, then create
			// if nacId == "" || containsNacId(out, nacId) {

			if nacId == "" {
				log.Printf("[DEBUG] nacId is empty")
				allowipspoofing := fmt.Sprintf("network_attachments.%d.virtual_network_interface.0.allow_ip_spoofing", i)
				autodelete := fmt.Sprintf("network_attachments.%d.virtual_network_interface.0.auto_delete", i)
				enablenat := fmt.Sprintf("network_attachments.%d.virtual_network_interface.0.enable_infrastructure_nat", i)
				nacMap := nac.(map[string]interface{})
				VirtualNetworkInterfaceModel, err := resourceIBMIsInstanceMapToVirtualNetworkInterfacePrototypeAttachmentContext(allowipspoofing, autodelete, enablenat, d, nacMap["virtual_network_interface"].([]interface{})[0].(map[string]interface{}))
				if err != nil {
					return err
				}
				nacNameStr := nacMap["name"].(string)
				createInstanceNetworkAttachmentOptions := &vpcv1.CreateInstanceNetworkAttachmentOptions{
					InstanceID:              &id,
					Name:                    &nacNameStr,
					VirtualNetworkInterface: VirtualNetworkInterfaceModel,
				}
				_, res, err := instanceC.CreateInstanceNetworkAttachment(createInstanceNetworkAttachmentOptions)
				if err != nil {
					return fmt.Errorf("[ERROR] Error while creating network attachment(%s) of instance(%s) \n%s: %q", nacNameStr, d.Id(), err, res)
				}
			} else {
				log.Printf("[DEBUG] nacId is not empty")
				nacName := fmt.Sprintf("network_attachments.%d.name", i)
				nacVniName := fmt.Sprintf("network_attachments.%d.virtual_network_interface", i)
				primaryipName := fmt.Sprintf("%s.%s", nacVniName, "0.primary_ip")
				sgName := fmt.Sprintf("%s.%s", nacVniName, "0.security_groups")
				if d.HasChange(nacName) {
					networkName := d.Get(nacName).(string)
					updateInstanceNetworkAttachmentOptions := &vpcv1.UpdateInstanceNetworkAttachmentOptions{
						InstanceID: &id,
						ID:         &nacId,
					}
					instanceNetworkAttachmentPatch := &vpcv1.InstanceNetworkAttachmentPatch{
						Name: &networkName,
					}
					instanceNetworkAttachmentPatchAsPatch, err := instanceNetworkAttachmentPatch.AsPatch()
					if err != nil {
						return (fmt.Errorf("[ERROR] Error encountered while apply as patch for instanceNetworkAttachmentPatchAsPatch of network attachment(%s) of instance(%s) %s", nacId, id, err))
					}
					updateInstanceNetworkAttachmentOptions.InstanceNetworkAttachmentPatch = instanceNetworkAttachmentPatchAsPatch
					_, res, err := instanceC.UpdateInstanceNetworkAttachment(updateInstanceNetworkAttachmentOptions)
					if err != nil {
						return (fmt.Errorf("[ERROR] Error encountered while updating network attachment(%s) name of instance(%s) %s/n%s", nacId, id, err, res))
					}
					// output, err := json.MarshalIndent(updateInstanceNetworkAttachmentOptions, "", "    ")
					// if err == nil {
					// 	log.Printf("%+v\n", string(output))
					// } else {
					// 	log.Printf("Error : %#v", updateInstanceNetworkAttachmentOptions)
					// }
				}
				if d.HasChange(nacVniName) {
					vniId := d.Get(fmt.Sprintf("%s.%s", nacVniName, "0.id")).(string)
					updateVirtualNetworkInterfaceOptions := &vpcv1.UpdateVirtualNetworkInterfaceOptions{
						ID: &vniId,
					}
					virtualNetworkInterfacePatch := &vpcv1.VirtualNetworkInterfacePatch{}
					autoDeleteName := fmt.Sprintf("%s.%s", nacVniName, "0.auto_delete")
					nameName := fmt.Sprintf("%s.%s", nacVniName, "0.name")
					ipsName := fmt.Sprintf("%s.%s", nacVniName, "0.ips")
					enableNatName := fmt.Sprintf("%s.%s", nacVniName, "0.enable_infrastructure_nat")
					allowIpSpoofingName := fmt.Sprintf("%s.%s", nacVniName, "0.allow_ip_spoofing")
					pStateFilteringModeSchemaName := fmt.Sprintf("%s.%s", nacVniName, "0.protocol_state_filtering_mode")
					if d.HasChange(autoDeleteName) {
						autodelete := d.Get(autoDeleteName).(bool)
						virtualNetworkInterfacePatch.AutoDelete = &autodelete
					}
					if d.HasChange(nameName) {
						name := d.Get(nameName).(string)
						virtualNetworkInterfacePatch.Name = &name
					}
					if d.HasChange(enableNatName) {
						enableNat := d.Get(enableNatName).(bool)
						virtualNetworkInterfacePatch.EnableInfrastructureNat = &enableNat
					}
					if d.HasChange(allowIpSpoofingName) {
						allIpSpoofing := d.Get(allowIpSpoofingName).(bool)
						virtualNetworkInterfacePatch.AllowIPSpoofing = &allIpSpoofing
					}
					if d.HasChange(pStateFilteringModeSchemaName) {
						pStateFilteringMode := d.Get(pStateFilteringModeSchemaName).(string)
						virtualNetworkInterfacePatch.ProtocolStateFilteringMode = &pStateFilteringMode
					}
					virtualNetworkInterfacePatchAsPatch, err := virtualNetworkInterfacePatch.AsPatch()
					if err != nil {
						return fmt.Errorf("[ERROR] Error encountered while apply as patch for virtualNetworkInterfacePatch of instance(%s) vni (%s) %s", d.Id(), vniId, err)
					}
					updateVirtualNetworkInterfaceOptions.VirtualNetworkInterfacePatch = virtualNetworkInterfacePatchAsPatch
					_, response, err := instanceC.UpdateVirtualNetworkInterface(updateVirtualNetworkInterfaceOptions)
					if err != nil {
						log.Printf("[DEBUG] UpdateVirtualNetworkInterfaceWithContext failed %s\n%s", err, response)
						return fmt.Errorf("UpdateVirtualNetworkInterfaceWithContext failed during instance(%s) network attachment patch %s\n%s", d.Id(), err, response)
					}

					if d.HasChange(ipsName) {
						oldips, newips := d.GetChange(ipsName)
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
									_, response, err := instanceC.AddVirtualNetworkInterfaceIP(addVirtualNetworkInterfaceIPOptions)
									if err != nil {
										log.Printf("[DEBUG] AddVirtualNetworkInterfaceIPWithContext failed in VirtualNetworkInterface patch during instance nac patch %s\n%s", err, response)
										return fmt.Errorf("AddVirtualNetworkInterfaceIPWithContext failed in VirtualNetworkInterface patch during instance nac patch %s\n%s", err, response)
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
									response, err := instanceC.RemoveVirtualNetworkInterfaceIP(removeVirtualNetworkInterfaceIPOptions)
									if err != nil {
										log.Printf("[DEBUG] RemoveVirtualNetworkInterfaceIPWithContext failed in VirtualNetworkInterface patch during instance nac patch %s\n%s", err, response)
										return fmt.Errorf("RemoveVirtualNetworkInterfaceIPWithContext failed in VirtualNetworkInterface patch during instance nac patch %s\n%s", err, response)
									}
								}
							}
						}
					}

					if d.HasChange(primaryipName) {
						subnetIdName := fmt.Sprintf("%s.%s", nacVniName, "0.subnet")
						ripIdName := fmt.Sprintf("%s.%s", primaryipName, "0.reserved_ip")
						subnetId := d.Get(subnetIdName).(string)
						primaryipNameName := fmt.Sprintf("%s.%s", primaryipName, "0.name")
						primaryipAutoDeleteName := fmt.Sprintf("%s.%s", primaryipName, "0.name")
						ripId := d.Get(ripIdName).(string)
						updateripoptions := &vpcv1.UpdateSubnetReservedIPOptions{
							SubnetID: &subnetId,
							ID:       &ripId,
						}
						reservedIpPath := &vpcv1.ReservedIPPatch{}
						if d.HasChange(primaryipNameName) {
							name := d.Get(primaryipNameName).(string)
							reservedIpPath.Name = &name
						}
						if d.HasChange(primaryipAutoDeleteName) {
							auto := d.Get(primaryipAutoDeleteName).(bool)
							reservedIpPath.AutoDelete = &auto
						}
						reservedIpPathAsPatch, err := reservedIpPath.AsPatch()
						if err != nil {
							return fmt.Errorf("[ERROR] Error calling reserved ip as patch on vni patch \n%s", err)
						}
						updateripoptions.ReservedIPPatch = reservedIpPathAsPatch
						_, response, err := instanceC.UpdateSubnetReservedIP(updateripoptions)
						if err != nil {
							return fmt.Errorf("[ERROR] Error updating vni reserved ip(%s): %s\n%s", ripId, err, response)
						}
					}
					if d.HasChange(sgName) {
						ovs, nvs := d.GetChange(sgName)
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
								_, response, err := instanceC.CreateSecurityGroupTargetBinding(createsgnicoptions)
								if err != nil {
									return fmt.Errorf("[ERROR] Error while creating security group %q for virtual network interface %s\n%s: %q", add[i], vniId, err, response)
								}
								_, err = isWaitForVirtualNetworkInterfaceAvailable(instanceC, vniId, d.Timeout(schema.TimeoutUpdate))
								if err != nil {
									return err
								}
							}

						}
						if len(remove) > 0 {
							for i := range remove {
								deletesgnicoptions := &vpcv1.DeleteSecurityGroupTargetBindingOptions{
									SecurityGroupID: &remove[i],
									ID:              &vniId,
								}
								response, err := instanceC.DeleteSecurityGroupTargetBinding(deletesgnicoptions)
								if err != nil {
									return fmt.Errorf("[ERROR] Error while removing security group %q for virtual network interface %s\n%s: %q", remove[i], d.Id(), err, response)
								}
								_, err = isWaitForVirtualNetworkInterfaceAvailable(instanceC, vniId, d.Timeout(schema.TimeoutUpdate))
								if err != nil {
									return err
								}
							}
						}
					}

				}

			}
			// }
		}

	}

	//primary_network_attachment
	if d.HasChange("primary_network_attachment") && !d.IsNewResource() {
		networkID := d.Get("primary_network_attachment.0.id").(string)
		networkName := "primary_network_attachment.0.name"
		nacVniName := "primary_network_attachment.0.virtual_network_interface"
		if d.HasChange(networkName) {
			networkNameString := d.Get(networkName).(string)
			updateInstanceNetworkAttachmentOptions := &vpcv1.UpdateInstanceNetworkAttachmentOptions{
				InstanceID: &id,
				ID:         &networkID,
			}
			instanceNetworkAttachmentPatch := &vpcv1.InstanceNetworkAttachmentPatch{
				Name: &networkNameString,
			}
			instanceNetworkAttachmentPatchAsPatch, err := instanceNetworkAttachmentPatch.AsPatch()
			if err != nil {
				return (fmt.Errorf("[ERROR] Error encountered while apply as patch for instanceNetworkAttachmentPatchAsPatch of pna of instance(%s) %s", id, err))
			}
			updateInstanceNetworkAttachmentOptions.InstanceNetworkAttachmentPatch = instanceNetworkAttachmentPatchAsPatch
			_, res, err := instanceC.UpdateInstanceNetworkAttachment(updateInstanceNetworkAttachmentOptions)
			if err != nil {
				return (fmt.Errorf("[ERROR] Error encountered while updating pna name of instance(%s) %s/n%s", id, err, res))
			}
		}
		if d.HasChange(nacVniName) {
			vniId := d.Get(fmt.Sprintf("%s.%s", nacVniName, "0.id")).(string)
			updateVirtualNetworkInterfaceOptions := &vpcv1.UpdateVirtualNetworkInterfaceOptions{
				ID: &vniId,
			}
			virtualNetworkInterfacePatch := &vpcv1.VirtualNetworkInterfacePatch{}
			autoDeleteName := fmt.Sprintf("%s.%s", nacVniName, "0.auto_delete")
			nameName := fmt.Sprintf("%s.%s", nacVniName, "0.name")
			ipsName := fmt.Sprintf("%s.%s", nacVniName, "0.ips")
			primaryipName := fmt.Sprintf("%s.%s", nacVniName, "0.primary_ip")
			sgName := fmt.Sprintf("%s.%s", nacVniName, "0.security_groups")
			enableNatName := fmt.Sprintf("%s.%s", nacVniName, "0.enable_infrastructure_nat")
			allowIpSpoofingName := fmt.Sprintf("%s.%s", nacVniName, "0.allow_ip_spoofing")
			pStateFilteringModeSchemaName := fmt.Sprintf("%s.%s", nacVniName, "0.protocol_state_filtering_mode")
			if d.HasChange(autoDeleteName) {
				autodelete := d.Get(autoDeleteName).(bool)
				virtualNetworkInterfacePatch.AutoDelete = &autodelete
			}
			if d.HasChange(nameName) {
				name := d.Get(nameName).(string)
				virtualNetworkInterfacePatch.Name = &name
			}
			if d.HasChange(enableNatName) {
				enableNat := d.Get(enableNatName).(bool)
				virtualNetworkInterfacePatch.EnableInfrastructureNat = &enableNat
			}
			if d.HasChange(allowIpSpoofingName) {
				allIpSpoofing := d.Get(allowIpSpoofingName).(bool)
				virtualNetworkInterfacePatch.AllowIPSpoofing = &allIpSpoofing
			}
			if d.HasChange(pStateFilteringModeSchemaName) {
				pStateFilteringMode := d.Get(pStateFilteringModeSchemaName).(string)
				virtualNetworkInterfacePatch.ProtocolStateFilteringMode = &pStateFilteringMode
			}
			virtualNetworkInterfacePatchAsPatch, err := virtualNetworkInterfacePatch.AsPatch()
			if err != nil {
				return fmt.Errorf("[ERROR] Error encountered while apply as patch for virtualNetworkInterfacePatch of instance(%s) vni (%s) %s", d.Id(), vniId, err)
			}
			updateVirtualNetworkInterfaceOptions.VirtualNetworkInterfacePatch = virtualNetworkInterfacePatchAsPatch
			_, response, err := instanceC.UpdateVirtualNetworkInterface(updateVirtualNetworkInterfaceOptions)
			if err != nil {
				log.Printf("[DEBUG] UpdateVirtualNetworkInterfaceWithContext failed %s\n%s", err, response)
				return fmt.Errorf("UpdateVirtualNetworkInterfaceWithContext failed during instance(%s) network attachment patch %s\n%s", d.Id(), err, response)
			}

			if d.HasChange(ipsName) {
				oldips, newips := d.GetChange(ipsName)
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
							_, response, err := instanceC.AddVirtualNetworkInterfaceIP(addVirtualNetworkInterfaceIPOptions)
							if err != nil {
								log.Printf("[DEBUG] AddVirtualNetworkInterfaceIPWithContext failed in VirtualNetworkInterface patch during instance nac patch %s\n%s", err, response)
								return fmt.Errorf("AddVirtualNetworkInterfaceIPWithContext failed in VirtualNetworkInterface patch during instance nac patch %s\n%s", err, response)
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
							response, err := instanceC.RemoveVirtualNetworkInterfaceIP(removeVirtualNetworkInterfaceIPOptions)
							if err != nil {
								log.Printf("[DEBUG] RemoveVirtualNetworkInterfaceIPWithContext failed in VirtualNetworkInterface patch during instance nac patch %s\n%s", err, response)
								return fmt.Errorf("RemoveVirtualNetworkInterfaceIPWithContext failed in VirtualNetworkInterface patch during instance nac patch %s\n%s", err, response)
							}
						}
					}
				}
			}

			if d.HasChange(primaryipName) {
				subnetIdName := fmt.Sprintf("%s.%s", nacVniName, "0.subnet")
				ripIdName := fmt.Sprintf("%s.%s", primaryipName, "0.reserved_ip")
				subnetId := d.Get(subnetIdName).(string)
				primaryipNameName := fmt.Sprintf("%s.%s", primaryipName, "0.name")
				primaryipAutoDeleteName := fmt.Sprintf("%s.%s", primaryipName, "0.name")
				ripId := d.Get(ripIdName).(string)
				updateripoptions := &vpcv1.UpdateSubnetReservedIPOptions{
					SubnetID: &subnetId,
					ID:       &ripId,
				}
				reservedIpPath := &vpcv1.ReservedIPPatch{}
				if d.HasChange(primaryipNameName) {
					name := d.Get(primaryipNameName).(string)
					reservedIpPath.Name = &name
				}
				if d.HasChange(primaryipAutoDeleteName) {
					auto := d.Get(primaryipAutoDeleteName).(bool)
					reservedIpPath.AutoDelete = &auto
				}
				reservedIpPathAsPatch, err := reservedIpPath.AsPatch()
				if err != nil {
					return fmt.Errorf("[ERROR] Error calling reserved ip as patch on vni patch \n%s", err)
				}
				updateripoptions.ReservedIPPatch = reservedIpPathAsPatch
				_, response, err := instanceC.UpdateSubnetReservedIP(updateripoptions)
				if err != nil {
					return fmt.Errorf("[ERROR] Error updating vni reserved ip(%s): %s\n%s", ripId, err, response)
				}
			}
			if d.HasChange(sgName) {
				ovs, nvs := d.GetChange(sgName)
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
						_, response, err := instanceC.CreateSecurityGroupTargetBinding(createsgnicoptions)
						if err != nil {
							return fmt.Errorf("[ERROR] Error while creating security group %q for virtual network interface %s\n%s: %q", add[i], vniId, err, response)
						}
						_, err = isWaitForVirtualNetworkInterfaceAvailable(instanceC, vniId, d.Timeout(schema.TimeoutUpdate))
						if err != nil {
							return err
						}
					}

				}
				if len(remove) > 0 {
					for i := range remove {
						deletesgnicoptions := &vpcv1.DeleteSecurityGroupTargetBindingOptions{
							SecurityGroupID: &remove[i],
							ID:              &vniId,
						}
						response, err := instanceC.DeleteSecurityGroupTargetBinding(deletesgnicoptions)
						if err != nil {
							return fmt.Errorf("[ERROR] Error while removing security group %q for virtual network interface %s\n%s: %q", remove[i], d.Id(), err, response)
						}
						_, err = isWaitForVirtualNetworkInterfaceAvailable(instanceC, vniId, d.Timeout(schema.TimeoutUpdate))
						if err != nil {
							return err
						}
					}
				}
			}

		}

	}

	resPol := "reservation_affinity.0.policy"
	resPool := "reservation_affinity.0.pool"

	if (d.HasChange(resPol) || d.HasChange(resPool)) && !d.IsNewResource() {
		if resAffinity, ok := d.GetOk(isReservationAffinity); ok {
			getinsOptions := &vpcv1.GetInstanceOptions{
				ID: &id,
			}
			_, response, err := instanceC.GetInstance(getinsOptions)
			if err != nil {
				return fmt.Errorf("[ERROR] Error getting instance (%s): %s\n%s", id, err, response)
			}
			eTag := response.Headers.Get("ETag")

			resAff := resAffinity.([]interface{})[0].(map[string]interface{})
			var resAffinityPatch = &vpcv1.InstanceReservationAffinityPatch{}
			policy, ok := resAff["policy"]
			policyStr := policy.(string)
			idStr := ""
			if policyStr != "" && ok {
				resAffinityPatch.Policy = &policyStr
			}
			if d.HasChange(resPool) {
				poolIntf, okPool := resAff[isReservationAffinityPool]
				if okPool && poolIntf != nil && poolIntf.([]interface{}) != nil && len(poolIntf.([]interface{})) > 0 {
					pool := poolIntf.([]interface{})[0].(map[string]interface{})
					id, okId := pool["id"]
					if okId {
						idStr, ok = id.(string)
						if idStr != "" && ok {
							var resAffPool = make([]vpcv1.ReservationIdentityIntf, 1)
							resAffPool[0] = &vpcv1.ReservationIdentity{
								ID: &idStr,
							}
							resAffinityPatch.Pool = resAffPool
						}
					}

				}
			}

			instancePatchModel := &vpcv1.InstancePatch{
				ReservationAffinity: resAffinityPatch,
			}
			mpatch, err := instancePatchModel.AsPatch()
			if err != nil {
				return fmt.Errorf("[ERROR] Error calling asPatch with reservation affinity: %s", err)
			}
			//Detaching the reservation from the reserved instance
			if policyStr == "disabled" && idStr == "" {
				resAffMap := mpatch["reservation_affinity"].(map[string]interface{})
				resAffMap["pool"] = nil
				mpatch["reservation_affinity"] = resAffMap
			}
			param := &vpcv1.UpdateInstanceOptions{
				InstancePatch: mpatch,
				ID:            &id,
			}
			param.IfMatch = &eTag
			_, _, err = instanceC.UpdateInstance(param)
			if err != nil {
				return err
			}
		}
	}

	bootVolSize := "boot_volume.0.size"
	bootIopsSize := "boot_volume.0.iops"

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
	if d.HasChange(bootIopsSize) && !d.IsNewResource() {
		_, new := d.GetChange(bootIopsSize)

		bootVolIops := int64(new.(int))
		volId := d.Get("boot_volume.0.volume_id").(string)
		updateVolumeOptions := &vpcv1.UpdateVolumeOptions{
			ID: &volId,
		}
		volPatchModel := &vpcv1.VolumePatch{
			Iops: &bootVolIops,
		}
		volPatchModelAsPatch, err := volPatchModel.AsPatch()

		if err != nil {
			return (fmt.Errorf("[ERROR] Error encountered while apply as patch for boot iops of instance %s", err))
		}

		updateVolumeOptions.VolumePatch = volPatchModelAsPatch

		vol, res, err := instanceC.UpdateVolume(updateVolumeOptions)

		if vol == nil || err != nil {
			return (fmt.Errorf("[ERROR] Error encountered while expanding boot iops of instance %s/n%s", err, res))
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
	bootVolName := "boot_volume.0.name"
	if d.HasChange(bootVolName) && !d.IsNewResource() {
		volId := d.Get("boot_volume.0.volume_id").(string)
		volName := d.Get(bootVolName).(string)
		updateVolumeOptions := &vpcv1.UpdateVolumeOptions{
			ID: &volId,
		}
		volPatchModel := &vpcv1.VolumePatch{
			Name: &volName,
		}
		volPatchModelAsPatch, err := volPatchModel.AsPatch()

		if err != nil {
			return (fmt.Errorf("[ERROR] Error encountered while apply as patch for boot volume name update of instance %s", err))
		}

		updateVolumeOptions.VolumePatch = volPatchModelAsPatch

		vol, res, err := instanceC.UpdateVolume(updateVolumeOptions)

		if vol == nil || err != nil {
			return (fmt.Errorf("[ERROR] Error encountered while updating name of boot volume of instance %s/n%s", err, res))
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

	if (d.HasChange(isInstanceName) || d.HasChange("confidential_compute_mode") || d.HasChange("enable_secure_boot")) && !d.IsNewResource() {
		restartNeeded := false
		serverstopped := false
		name := d.Get(isInstanceName).(string)
		updnetoptions := &vpcv1.UpdateInstanceOptions{
			ID: &id,
		}
		instancePatchModel := &vpcv1.InstancePatch{}
		if d.HasChange("confidential_compute_mode") {
			instancePatchModel.ConfidentialComputeMode = core.StringPtr(d.Get("confidential_compute_mode").(string))
			restartNeeded = true
		}
		if _, ok := d.GetOkExists("enable_secure_boot"); ok && d.HasChange("enable_secure_boot") {
			instancePatchModel.EnableSecureBoot = core.BoolPtr(d.Get("enable_secure_boot").(bool))
		}
		if d.HasChange("name") {
			instancePatchModel.Name = &name
		}

		instancePatch, err := instancePatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("[ERROR] Error calling asPatch for InstancePatch: %s", err)
		}
		updnetoptions.InstancePatch = instancePatch
		if restartNeeded {
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
				serverstopped = true
			}
		}
		_, _, err = instanceC.UpdateInstance(updnetoptions)
		if err != nil {
			return err
		}
		if serverstopped {
			actiontype := "start"
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
			_, err = isWaitForInstanceAvailable(instanceC, d.Id(), d.Timeout(schema.TimeoutUpdate), d)
			if err != nil {
				return err
			}
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
			return fmt.Errorf("[ERROR] Error calling asPatch for InstancePatch: %s", err)
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
			return fmt.Errorf("[ERROR] Error calling asPatch for InstancePatch: %s", err)
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
			return fmt.Errorf("[ERROR] Error calling asPatch for InstancePatch: %s", err)
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
				"[ERROR] Error on update of resource Instance (%s) tags: %s", d.Id(), err)
		}
	}
	if d.HasChange(isInstanceAccessTags) {
		oldList, newList := d.GetChange(isInstanceAccessTags)
		err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *instance.CRN, "", isInstanceAccessTagType)
		if err != nil {
			log.Printf(
				"[ERROR] Error on update of resource Instance (%s) access tags: %s", d.Id(), err)
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
			if *vol.Type == "data" && *vol.DeleteVolumeOnInstanceDelete {
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

func resourceIbmIsInstanceDedicatedHostGroupReferenceDeletedToMap(dedicatedHostGroupReferenceDeleted vpcv1.Deleted) map[string]interface{} {
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

func resourceIBMIsInstanceInstanceNetworkAttachmentReferenceToMap(model *vpcv1.InstanceNetworkAttachmentReference, pna *vpcv1.InstanceNetworkAttachment, instanceC *vpcv1.VpcV1, autoDelete bool) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Deleted != nil {
		deletedMap, err := resourceIBMIsInstanceInstanceNetworkAttachmentReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	if model.Href != nil {
		modelMap["href"] = model.Href
	}
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	vniMap := make(map[string]interface{})
	if pna.VirtualNetworkInterface != nil {
		vniMap["id"] = *pna.VirtualNetworkInterface.ID
		vniMap["name"] = pna.VirtualNetworkInterface.Name
		vniMap["resource_type"] = pna.VirtualNetworkInterface.ResourceType
	}
	if model.PrimaryIP != nil {
		primaryipmap, _ := resourceIBMIsInstancePrimaryIPReferenceToMap(model.PrimaryIP)
		modelMap["primary_ip"] = []map[string]interface{}{primaryipmap}
	}
	getVirtualNetworkInterfaceOptions := &vpcv1.GetVirtualNetworkInterfaceOptions{
		ID: pna.VirtualNetworkInterface.ID,
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
				ipsItemMap, err := resourceIBMIsVirtualNetworkInterfaceReservedIPReferenceToMap(&ipsItem, autoDelete)
				if err != nil {
					return nil, err
				}
				ips = append(ips, ipsItemMap)
			}
		}
		vniMap["ips"] = ips
	}

	if !core.IsNil(vniDetails.SecurityGroups) {
		securityGroups := make([]string, 0)
		for _, securityGroupsItem := range vniDetails.SecurityGroups {
			if securityGroupsItem.ID != nil {
				securityGroups = append(securityGroups, *securityGroupsItem.ID)
			}
		}
		vniMap["security_groups"] = securityGroups
	}
	primaryIPMap, err := resourceIBMIsInstanceReservedIPReferenceToMap(model.PrimaryIP, autoDelete)
	if err != nil {
		return modelMap, err
	}
	vniMap["primary_ip"] = []map[string]interface{}{primaryIPMap}
	if model.ResourceType != nil {
		modelMap["resource_type"] = *model.ResourceType
	}
	if model.Subnet != nil {
		vniMap["subnet"] = *model.Subnet.ID
	}
	modelMap["virtual_network_interface"] = []map[string]interface{}{vniMap}
	return modelMap, nil
}

func resourceIBMIsInstanceInstanceNetworkAttachmentReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func resourceIBMIsInstanceReservedIPReferenceToMap(model *vpcv1.ReservedIPReference, autoDelete bool) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["address"] = model.Address
	if model.Deleted != nil {
		deletedMap, err := resourceIBMIsInstanceReservedIPReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = model.Href
	modelMap["auto_delete"] = autoDelete
	modelMap["reserved_ip"] = model.ID
	modelMap["name"] = model.Name
	if model.ResourceType != nil {
		modelMap["resource_type"] = *model.ResourceType
	}
	return modelMap, nil
}
func resourceIBMIsInstancePrimaryIPReferenceToMap(model *vpcv1.ReservedIPReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["address"] = model.Address
	if model.Deleted != nil {
		deletedMap, err := resourceIBMIsInstanceReservedIPReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	if model.ResourceType != nil {
		modelMap["resource_type"] = model.ResourceType
	}
	return modelMap, nil
}

func resourceIBMIsInstanceReservedIPReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func resourceIBMIsInstanceMapToInstanceNetworkAttachmentPrototype(allowipspoofing, autodelete, enablenat string, d *schema.ResourceData, modelMap map[string]interface{}) (*vpcv1.InstanceNetworkAttachmentPrototype, error) {
	model := &vpcv1.InstanceNetworkAttachmentPrototype{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	VirtualNetworkInterfaceModel, err := resourceIBMIsInstanceMapToVirtualNetworkInterfacePrototypeAttachmentContext(allowipspoofing, autodelete, enablenat, d, modelMap["virtual_network_interface"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.VirtualNetworkInterface = VirtualNetworkInterfaceModel
	return model, nil
}
func resourceIBMIsInstanceMapToVirtualNetworkInterfacePrototypeAttachmentContext(allowipspoofing, autodelete, enablenat string, d *schema.ResourceData, modelMap map[string]interface{}) (vpcv1.InstanceNetworkAttachmentPrototypeVirtualNetworkInterfaceIntf, error) {
	model := &vpcv1.InstanceNetworkAttachmentPrototypeVirtualNetworkInterface{}
	if allowipspoofingOk, ok := d.GetOkExists(allowipspoofing); ok {
		model.AllowIPSpoofing = core.BoolPtr(allowipspoofingOk.(bool))
	}
	if autodeleteOk, ok := d.GetOkExists(autodelete); ok {
		model.AutoDelete = core.BoolPtr(autodeleteOk.(bool))
	}
	if enablenatok, ok := d.GetOkExists(enablenat); ok {
		model.EnableInfrastructureNat = core.BoolPtr(enablenatok.(bool))
	}
	if modelMap["ips"] != nil {
		ips := []vpcv1.VirtualNetworkInterfaceIPPrototypeIntf{}
		for _, ipsItem := range modelMap["ips"].(*schema.Set).List() {
			ipsItemModel, err := resourceIBMIsInstanceMapToVirtualNetworkInterfaceIPsReservedIPPrototype(ipsItem.(map[string]interface{}))
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
	if modelMap["protocol_state_filtering_mode"] != nil {
		if pStateFilteringInt, ok := modelMap["protocol_state_filtering_mode"]; ok && pStateFilteringInt.(string) != "" {
			model.ProtocolStateFilteringMode = core.StringPtr(pStateFilteringInt.(string))
		}
	}
	if modelMap["primary_ip"] != nil && len(modelMap["primary_ip"].([]interface{})) > 0 {
		PrimaryIPModel, err := resourceIBMIsInstanceMapToVirtualNetworkInterfacePrimaryIPReservedIPPrototype(modelMap["primary_ip"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.PrimaryIP = PrimaryIPModel
	}
	if modelMap["resource_group"] != nil && modelMap["resource_group"].(string) != "" {
		resourcegroupid := modelMap["resource_group"].(string)
		model.ResourceGroup = &vpcv1.ResourceGroupIdentity{
			ID: &resourcegroupid,
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
		model.Subnet = &vpcv1.SubnetIdentityByID{
			ID: &subnetId,
		}
	}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	return model, nil
}

func resourceIBMIsInstanceMapToVirtualNetworkInterfaceIPsReservedIPPrototype(modelMap map[string]interface{}) (vpcv1.VirtualNetworkInterfaceIPPrototypeIntf, error) {
	model := &vpcv1.VirtualNetworkInterfaceIPPrototype{}
	if modelMap["reserved_ip"] != nil && modelMap["reserved_ip"].(string) != "" {
		model.ID = core.StringPtr(modelMap["reserved_ip"].(string))
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
func resourceIBMIsInstanceMapToVirtualNetworkInterfacePrimaryIPReservedIPPrototype(modelMap map[string]interface{}) (vpcv1.VirtualNetworkInterfacePrimaryIPPrototypeIntf, error) {
	model := &vpcv1.VirtualNetworkInterfacePrimaryIPPrototype{}
	if modelMap["reserved_ip"] != nil && modelMap["reserved_ip"].(string) != "" {
		model.ID = core.StringPtr(modelMap["reserved_ip"].(string))
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

func containsNacId(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func ResourceIBMIsInstanceInstanceClusterNetworkAttachmentReferenceToMap(instanceC *vpcv1.VpcV1, model *vpcv1.InstanceClusterNetworkAttachmentReference, id string) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	getInstanceClusterNetworkAttachment := &vpcv1.GetInstanceClusterNetworkAttachmentOptions{
		InstanceID: &id,
		ID:         model.ID,
	}
	clusterNetworkAttachment, _, err := instanceC.GetInstanceClusterNetworkAttachment(getInstanceClusterNetworkAttachment)
	if err != nil {
		return modelMap, err
	}
	if clusterNetworkAttachment.ClusterNetworkInterface != nil {
		clusterMap, err := ResourceIBMIsInstanceClusterNetworkAttachmentToMap(clusterNetworkAttachment)
		if err != nil {
			return modelMap, err
		}
		modelMap["cluster_network_interface"] = []map[string]interface{}{clusterMap}
	}

	return modelMap, nil
}

func ResourceIBMIsInstanceClusterNetworkAttachmentToMap(cnamodel *vpcv1.InstanceClusterNetworkAttachment) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	model := cnamodel.ClusterNetworkInterface
	if model.Deleted != nil {
		deletedMap, err := ResourceIBMIsInstanceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	if model.Subnet != nil {
		subnetMap, err := ResourceIBMIsInstanceClusterNetworkInterfaceSubnetToMap(model.Subnet)
		if err != nil {
			return modelMap, err
		}
		modelMap["subnet"] = []map[string]interface{}{subnetMap}
	}
	if model.PrimaryIP != nil {
		primaryipMap, err := ResourceIBMIsInstanceClusterNetworkInterfacePrimaryIPToMap(model.PrimaryIP)
		if err != nil {
			return modelMap, err
		}
		modelMap["primary_ip"] = []map[string]interface{}{primaryipMap}
	}
	return modelMap, nil
}
func ResourceIBMIsInstanceClusterNetworkInterfacePrimaryIPToMap(model *vpcv1.ClusterNetworkSubnetReservedIPReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["address"] = *model.Address
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	return modelMap, nil
}
func ResourceIBMIsInstanceClusterNetworkInterfaceSubnetToMap(model *vpcv1.ClusterNetworkSubnetReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	return modelMap, nil
}
func ResourceIBMIsInstanceClusterNetworkReferenceToMap(model *vpcv1.ClusterNetworkReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	if model.Deleted != nil {
		deletedMap, err := ResourceIBMIsInstanceDeletedToMap(model.Deleted)
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

func ResourceIBMIsInstanceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = *model.MoreInfo
	return modelMap, nil
}

func ResourceIBMIsInstanceMapToInstanceClusterNetworkAttachmentPrototypeInstanceContext(modelMap map[string]interface{}) (*vpcv1.InstanceClusterNetworkAttachmentPrototypeInstanceContext, error) {
	model := &vpcv1.InstanceClusterNetworkAttachmentPrototypeInstanceContext{}
	ClusterNetworkInterfaceModel, err := ResourceIBMIsInstanceMapToInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterface(modelMap["cluster_network_interface"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.ClusterNetworkInterface = ClusterNetworkInterfaceModel
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	return model, nil
}

func ResourceIBMIsInstanceMapToInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterface(modelMap map[string]interface{}) (vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceIntf, error) {
	model := &vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterface{}
	if modelMap["auto_delete"] != nil {
		model.AutoDelete = core.BoolPtr(modelMap["auto_delete"].(bool))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["primary_ip"] != nil && len(modelMap["primary_ip"].([]interface{})) > 0 {
		PrimaryIPModel, err := ResourceIBMIsInstanceMapToClusterNetworkInterfacePrimaryIPPrototype(modelMap["primary_ip"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.PrimaryIP = PrimaryIPModel
	}
	if modelMap["subnet"] != nil && len(modelMap["subnet"].([]interface{})) > 0 {
		SubnetModel, err := ResourceIBMIsInstanceMapToClusterNetworkSubnetIdentity(modelMap["subnet"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Subnet = SubnetModel
	}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	return model, nil
}

func ResourceIBMIsInstanceMapToClusterNetworkInterfacePrimaryIPPrototype(modelMap map[string]interface{}) (vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeIntf, error) {
	model := &vpcv1.ClusterNetworkInterfacePrimaryIPPrototype{}
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
func ResourceIBMIsInstanceMapToClusterNetworkSubnetIdentity(modelMap map[string]interface{}) (vpcv1.ClusterNetworkSubnetIdentityIntf, error) {
	model := &vpcv1.ClusterNetworkSubnetIdentity{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	return model, nil
}
func diffSuppressVolumePrototypes(k, old, new string, d *schema.ResourceData) bool {
	if d.Id() == "" {
		return false
	}

	o, n := d.GetChange("volume_prototypes")
	oldList := o.([]interface{})
	newList := n.([]interface{})

	if len(oldList) != len(newList) {
		return false
	}

	// First, find which volume in new list corresponds to each old volume by name
	volMap := make(map[string]int) // maps attachment name to position in new list
	for i, v := range newList {
		vol := v.(map[string]interface{})
		attachmentName := vol["name"].(string)
		volMap[attachmentName] = i
	}

	// Compare each old volume with its corresponding new volume
	for _, oldVol := range oldList {
		oldVolMap := oldVol.(map[string]interface{})
		attachmentName := oldVolMap["name"].(string)

		// Find corresponding new volume
		newIndex, exists := volMap[attachmentName]
		if !exists {
			return false
		}

		newVol := newList[newIndex].(map[string]interface{})

		// Compare relevant fields
		if !volumesEqual(oldVolMap, newVol) {
			return false
		}
	}

	return true
}

func volumesEqual(oldVol, newVol map[string]interface{}) bool {
	fieldsToCompare := []string{
		"delete_volume_on_instance_delete",
		"volume_name",
		"volume_capacity",
		"volume_profile",
		"volume_source_snapshot",
		"volume_encryption_key",
		"volume_tags",
	}

	for _, field := range fieldsToCompare {
		oldVal, oldOk := oldVol[field]
		newVal, newOk := newVol[field]

		if oldOk != newOk {
			return false
		}

		if oldOk && newOk {
			if field == "volume_tags" {
				if !compareVolumeTags(oldVal, newVal) {
					return false
				}
				continue
			}

			if !reflect.DeepEqual(oldVal, newVal) {
				return false
			}
		}
	}

	// Handle IOPS specially based on profile
	profile := oldVol["volume_profile"].(string)
	if !isTieredProfile(profile) {
		oldIops, oldOk := oldVol["volume_iops"]
		newIops, newOk := newVol["volume_iops"]

		if oldOk != newOk {
			return false
		}

		if oldOk && newOk && !reflect.DeepEqual(oldIops, newIops) {
			return false
		}
	}

	return true
}

// Validation function
func ResourceValidateInstanceVolumePrototypes(diff *schema.ResourceDiff, meta interface{}) error {
	// For new resource creation
	if diff.Id() == "" {
		volProtoListIntf := diff.Get("volume_prototypes")
		if volProtoListIntf == nil {
			return nil
		}

		volProtoList := volProtoListIntf.([]interface{})
		for i, vol := range volProtoList {
			volMap := vol.(map[string]interface{})
			profile := volMap["volume_profile"].(string)

			// For tiered profiles, validate IOPS not set
			if isTieredProfile(profile) {
				if iops, ok := volMap["volume_iops"]; ok && iops.(int) != 0 {
					return fmt.Errorf("volume prototype %d (%s): iops cannot be set for tiered profile %s",
						i, volMap["volume_name"].(string), profile)
				}
			}
		}
		return nil
	}

	// For updates
	if !diff.HasChange("volume_prototypes") {
		return nil
	}

	oldVolProtoListIntf, newVolProtoListIntf := diff.GetChange("volume_prototypes")
	if oldVolProtoListIntf == nil || newVolProtoListIntf == nil {
		return nil
	}

	oldVolProtoList := oldVolProtoListIntf.([]interface{})
	newVolProtoList := newVolProtoListIntf.([]interface{})

	oldVolMap := make(map[string]map[string]interface{})
	for _, v := range oldVolProtoList {
		volMap := v.(map[string]interface{})
		oldVolMap[volMap["volume_name"].(string)] = volMap
	}

	// Validate each volume
	for _, v := range newVolProtoList {
		volMap := v.(map[string]interface{})
		volName := volMap["volume_name"].(string)
		newProfile := volMap["volume_profile"].(string)

		if oldVol, exists := oldVolMap[volName]; exists {
			oldProfile := oldVol["volume_profile"].(string)

			// Validate profile transitions
			if oldProfile != newProfile {
				if oldProfile == "custom" && newProfile != "custom" {
					return fmt.Errorf("volume %s: custom profile can only be changed to another custom profile", volName)
				}
				if isTieredProfile(oldProfile) && !isTieredProfile(newProfile) {
					return fmt.Errorf("volume %s: tiered profile can only be changed to another tiered profile", volName)
				}
			}
		}

		// Validate tiered profile constraints
		if isTieredProfile(newProfile) {
			if iops, ok := volMap["volume_iops"]; ok && iops.(int) != 0 {
				return fmt.Errorf("volume %s: iops cannot be set for tiered profile", volName)
			}
		}
	}

	return nil
}

func handleVolumePrototypesUpdate(d *schema.ResourceData, instanceC *vpcv1.VpcV1) error {
	if !d.HasChange("volume_prototypes") || d.IsNewResource() {
		return nil
	}

	instanceID := d.Id()
	o, n := d.GetChange("volume_prototypes")
	oldList := o.([]interface{})
	newList := n.([]interface{})

	// Track processed old volumes
	processedOldVolumes := make(map[string]bool)

	// First create a map of old volumes by name for easy lookup
	oldVolMap := make(map[string]map[string]interface{})
	for _, v := range oldList {
		vol := v.(map[string]interface{})
		name := vol["name"].(string)
		oldVolMap[name] = vol
	}

	// Process new list for updates and additions
	for i, newVolInterface := range newList {
		newVol := newVolInterface.(map[string]interface{})
		name := newVol["name"].(string)

		// Check if volume exists in old list
		if oldVol, exists := oldVolMap[name]; exists {
			// Mark as processed
			processedOldVolumes[name] = true

			// Check if update is needed
			if hasVolumeChanged(d, i, oldVol, newVol) {
				// Handle update
				volID := oldVol["volume_id"].(string)

				voloptions := &vpcv1.UpdateVolumeOptions{
					ID: &volID,
				}
				getvoloptions := &vpcv1.GetVolumeOptions{
					ID: &volID,
				}
				_, res, err := instanceC.GetVolume(getvoloptions)
				if err != nil {
					return fmt.Errorf("error getting volume for patch for %s: %w", name, err)
				}
				eTag := res.Headers.Get("ETag")
				volumePatchModel := &vpcv1.VolumePatch{}
				if newVol["volume_profile"].(string) != oldVol["volume_profile"].(string) && isTieredProfile(newVol["volume_profile"].(string)) {
					volumePatchModel.Profile = &vpcv1.VolumeProfileIdentity{
						Name: core.StringPtr(newVol["volume_profile"].(string)),
					}
				}
				if newVol["volume_name"].(string) != oldVol["volume_name"].(string) {
					volumePatchModel.Name = core.StringPtr(newVol["volume_name"].(string))
				}
				if newVol["volume_tags"] == nil && oldVol["volume_tags"] == nil {
					// do nothing
				} else if newVol["volume_tags"] == nil && oldVol["volume_tags"] != nil {
					volumePatchModel.UserTags = nil
				} else if (newVol["volume_tags"] != nil && oldVol["volume_tags"] == nil) || (!newVol["volume_tags"].(*schema.Set).Equal(oldVol["volume_tags"].(*schema.Set))) {
					userTags := newVol["volume_tags"].(*schema.Set)
					userTagsArray := make([]string, userTags.Len())
					for i, userTag := range userTags.List() {
						userTagStr := userTag.(string)
						userTagsArray[i] = userTagStr
					}
					volumePatchModel.UserTags = userTagsArray
				}
				volumePatch, err := volumePatchModel.AsPatch()
				if err != nil {
					return fmt.Errorf("error creating volume patch for %s: %w", name, err)
				}
				voloptions.VolumePatch = volumePatch
				voloptions.SetIfMatch(eTag)
				_, response, err := instanceC.UpdateVolume(voloptions)
				if err != nil {
					return fmt.Errorf("error updating volume %s: %s\n%s", name, err, response)
				}
				eTag = response.Headers.Get("ETag")

				// Only include IOPS for non-tiered profiles
				if !isTieredProfile(newVol["volume_profile"].(string)) && (int64(newVol["volume_iops"].(int)) != int64(oldVol["volume_iops"].(int))) {
					volumeIopsPatchModel := &vpcv1.VolumePatch{}
					iops := int64(newVol["volume_iops"].(int))
					volumeIopsPatchModel.Iops = &iops
					volumePatch, err := volumeIopsPatchModel.AsPatch()
					if err != nil {
						return fmt.Errorf("error creating volume patch for iops update %s: %w", name, err)
					}
					voloptions.VolumePatch = volumePatch
					voloptions.SetIfMatch(eTag)
					_, response, err := instanceC.UpdateVolume(voloptions)
					if err != nil {
						return fmt.Errorf("error updating volume during iops update %s: %s\n%s", name, err, response)
					}
					eTag = response.Headers.Get("ETag")
				}
				// Only include capacity update
				if int64(newVol["volume_capacity"].(int)) != int64(oldVol["volume_capacity"].(int)) {
					volumeCapacityPatchModel := &vpcv1.VolumePatch{}
					capacity := int64(newVol["volume_capacity"].(int))
					volumeCapacityPatchModel.Capacity = &capacity
					volumePatch, err := volumeCapacityPatchModel.AsPatch()
					if err != nil {
						return fmt.Errorf("error creating volume patch for capacity update %s: %w", name, err)
					}
					voloptions.SetIfMatch(eTag)
					voloptions.VolumePatch = volumePatch
					_, response, err := instanceC.UpdateVolume(voloptions)
					if err != nil {
						return fmt.Errorf("error updating volume during capacity update %s: %s\n%s", name, err, response)
					}
					eTag = response.Headers.Get("ETag")
				}
			}
		} else {
			// Handle addition
			profile := newVol["volume_profile"].(string)
			capacity := int64(newVol["volume_capacity"].(int))
			volumeName := newVol["volume_name"].(string)

			createvolattoptions := &vpcv1.CreateInstanceVolumeAttachmentOptions{
				InstanceID: &instanceID,
			}
			volAtt := &vpcv1.VolumeAttachmentPrototypeVolume{
				Name: &volumeName,
				Profile: &vpcv1.VolumeProfileIdentity{
					Name: &profile,
				},
				Capacity: &capacity,
			}
			// Handle delete_volume_on_instance_delete using GetOkExists only for new volumes
			if volAutoDelete, ok := d.GetOkExists(fmt.Sprintf("volume_prototypes.%d.delete_volume_on_instance_delete", i)); ok {
				createvolattoptions.DeleteVolumeOnInstanceDelete = core.BoolPtr(volAutoDelete.(bool))
			}

			// Only set IOPS for non-tiered profiles
			if !isTieredProfile(profile) {
				iops := int64(newVol["volume_iops"].(int))
				volAtt.Iops = &iops
			}
			createvolattoptions.Volume = volAtt
			newVolume, _, err := instanceC.CreateInstanceVolumeAttachment(createvolattoptions)
			if err != nil {
				return fmt.Errorf("error attaching volume %s: %w", name, err)
			}

			_, err = isWaitForInstanceVolumeAttached(instanceC, d, instanceID, *newVolume.ID)
			if err != nil {
				return err
			}
		}
	}

	// Handle deletions - anything in old list that wasn't processed
	for _, oldVolInterface := range oldList {
		oldVol := oldVolInterface.(map[string]interface{})
		name := oldVol["name"].(string)

		if !processedOldVolumes[name] {
			// Handle deletion
			volID := oldVol["id"].(string)

			delvolattoptions := &vpcv1.DeleteInstanceVolumeAttachmentOptions{
				InstanceID: &instanceID,
				ID:         &volID,
			}

			_, err := instanceC.DeleteInstanceVolumeAttachment(delvolattoptions)
			if err != nil {
				return fmt.Errorf("error removing volume %s: %w", name, err)
			}

			_, err = isWaitForInstanceVolumeDetached(instanceC, d, instanceID, volID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Modified to handle boolean comparison correctly between state and config
func hasVolumeChanged(d *schema.ResourceData, newIndex int, oldVol, newVol map[string]interface{}) bool {
	fieldsToCompare := []string{
		"volume_name",
		"volume_capacity",
		"volume_profile",
	}

	// Compare standard fields
	for _, field := range fieldsToCompare {
		oldVal := oldVol[field]
		newVal := newVol[field]

		if !reflect.DeepEqual(oldVal, newVal) {
			return true
		}
	}

	// Compare delete_volume_on_instance_delete
	// For old (state), direct access
	oldAutoDelete := oldVol["delete_volume_on_instance_delete"].(bool)

	// For new (config), use GetOkExists
	if newAutoDelete, ok := d.GetOkExists(fmt.Sprintf("volume_prototypes.%d.delete_volume_on_instance_delete", newIndex)); ok {
		if oldAutoDelete != newAutoDelete.(bool) {
			return true
		}
	}

	// Special handling for IOPS based on profile
	newProfile := newVol["volume_profile"].(string)

	if !isTieredProfile(newProfile) {
		oldIops := oldVol["volume_iops"].(int)
		newIops := newVol["volume_iops"].(int)
		if oldIops != newIops {
			return true
		}
	}

	return false
}

func isTieredProfile(profile string) bool {
	switch profile {
	case "general-purpose", "10iops-tier", "5iops-tier":
		return true
	default:
		return false
	}
}

// Helper function to compare volume tags
func compareVolumeTags(old, new interface{}) bool {
	if old == nil && new == nil {
		return true
	}
	if old == nil || new == nil {
		return false
	}

	oldSet := old.(*schema.Set)
	newSet := new.(*schema.Set)

	return oldSet.Len() == newSet.Len() && oldSet.Difference(newSet).Len() == 0
}

func prettifyResponse(response interface{}) string {
	output, err := json.MarshalIndent(response, "", "    ")
	if err == nil {
		return fmt.Sprintf("%+v\n", string(output))
	}
	return fmt.Sprintf("Error : %#v", response)
}

func setVolumePrototypesInState(d *schema.ResourceData, instance *vpcv1.Instance, instanceC *vpcv1.VpcV1) ([]map[string]interface{}, error) {
	if instance.VolumeAttachments == nil {
		return nil, nil
	}

	// First get the config order
	configVolumes := make(map[string]int) // maps attachment name to position
	if configList, ok := d.GetOk("volume_prototypes"); ok {
		for i, v := range configList.([]interface{}) {
			vol := v.(map[string]interface{})
			if name, ok := vol["name"].(string); ok {
				configVolumes[name] = i
			}
		}
	}

	// Create a map for current volumes
	currentVolumes := make(map[string]map[string]interface{})
	maxPosition := -1

	// Process all volumes
	for _, volume := range instance.VolumeAttachments {
		if *volume.ID != *instance.BootVolumeAttachment.ID {
			vol := map[string]interface{}{}

			if volume.Volume != nil {
				getVolOptions := &vpcv1.GetVolumeOptions{
					ID: volume.Volume.ID,
				}

				getInstanceVolumeAttachmentOptions := &vpcv1.GetInstanceVolumeAttachmentOptions{
					InstanceID: core.StringPtr(d.Id()),
					ID:         volume.ID,
				}

				volumeRef, _, err := instanceC.GetVolume(getVolOptions)
				if err != nil {
					vol["id"] = *volume.ID
					vol["volume_id"] = *volume.Volume.ID
					vol["name"] = *volume.Name
					vol["volume_name"] = *volume.Volume.Name
					vol["volume_crn"] = *volume.Volume.CRN
					vol["volume_resource_type"] = *volume.Volume.ResourceType
				} else {
					vol["id"] = *volume.ID
					vol["volume_id"] = *volume.Volume.ID
					vol["name"] = *volume.Name
					vol["volume_name"] = *volumeRef.Name
					vol["volume_profile"] = *volumeRef.Profile.Name
					vol["volume_iops"] = *volumeRef.Iops
					vol["volume_capacity"] = *volumeRef.Capacity
					vol["volume_crn"] = *volume.Volume.CRN
					vol["volume_resource_type"] = *volume.Volume.ResourceType
				}

				volumeAttRef, _, err := instanceC.GetInstanceVolumeAttachment(getInstanceVolumeAttachmentOptions)
				if err != nil {
					vol["delete_volume_on_instance_delete"] = true
				} else {
					vol["delete_volume_on_instance_delete"] = volumeAttRef.DeleteVolumeOnInstanceDelete
				}

				currentVolumes[*volume.Name] = vol

				// Track maximum position
				if pos, exists := configVolumes[*volume.Name]; exists {
					if pos > maxPosition {
						maxPosition = pos
					}
				}
			}
		}
	}

	// Create ordered list based on config positions
	orderedList := make([]map[string]interface{}, maxPosition+1)
	unorderedVolumes := make([]map[string]interface{}, 0)

	// First place volumes that exist in config
	for name, vol := range currentVolumes {
		if pos, exists := configVolumes[name]; exists {
			orderedList[pos] = vol
		} else {
			unorderedVolumes = append(unorderedVolumes, vol)
		}
	}

	// Remove nil entries and append any volumes not in config
	finalList := make([]map[string]interface{}, 0)
	for _, vol := range orderedList {
		if vol != nil {
			finalList = append(finalList, vol)
		}
	}
	finalList = append(finalList, unorderedVolumes...)

	return finalList, nil
}

// diffSuppressClusterNetworkAttachment handles comparing old and new cluster network attachments
// to determine if there are actual changes that require an update
func diffSuppressClusterNetworkAttachment(k, old, new string, d *schema.ResourceData) bool {
	// If values are equal, no changes needed

	if old == new {
		return true
	}

	// Get the lists of old and new attachments
	oldAttachments, newAttachments := []interface{}{}, []interface{}{}
	if v, ok := d.GetOk("cluster_network_attachments"); ok {
		newAttachments = v.([]interface{})
	}
	if v, ok := d.GetOk("cluster_network_attachments"); ok {
		oldAttachments = v.([]interface{})
	}

	// If lengths differ, there are definitely changes
	if len(oldAttachments) != len(newAttachments) {
		return false
	}

	// Compare each attachment
	for i := range oldAttachments {
		oldAttach := oldAttachments[i].(map[string]interface{})
		newAttach := newAttachments[i].(map[string]interface{})

		// Compare cluster_network_interface
		oldInterface := oldAttach["cluster_network_interface"].([]interface{})[0].(map[string]interface{})
		newInterface := newAttach["cluster_network_interface"].([]interface{})[0].(map[string]interface{})

		// Compare key properties
		if !compareInterfaces(oldInterface, newInterface) {
			return false
		}
	}

	return true
}

// comparePrimaryIP compares primary IP configurations
func comparePrimaryIP(old, new map[string]interface{}) bool {
	return old["id"] == new["id"] &&
		old["address"] == new["address"] &&
		old["auto_delete"] == new["auto_delete"] &&
		old["name"] == new["name"]
}

// compareSubnet compares subnet configurations
func compareSubnet(old, new map[string]interface{}) bool {
	return old["id"] == new["id"]
}
func handleClusterNetworkAttachmentUpdate(d *schema.ResourceData, instanceC *vpcv1.VpcV1) error {
	if d.HasChange("cluster_network_attachments") {
		old, new := d.GetChange("cluster_network_attachments")
		oldAttachments := old.([]interface{})
		newAttachments := new.([]interface{})

		// Build maps for both old and new attachments by name
		oldAttachMap := make(map[string]map[string]interface{})
		newAttachMap := make(map[string]map[string]interface{})

		// Map old attachments by name
		for _, attachment := range oldAttachments {
			attach := attachment.(map[string]interface{})
			name := attach["name"].(string)
			oldAttachMap[name] = attach
		}

		// Map new attachments by name and identify additions
		toAdd := []map[string]interface{}{}
		for _, attachment := range newAttachments {
			attach := attachment.(map[string]interface{})
			name := attach["name"].(string)
			newAttachMap[name] = attach

			// If name doesn't exist in old map, it's a new attachment
			if _, exists := oldAttachMap[name]; !exists {
				toAdd = append(toAdd, attach)
			}
		}

		// Identify removals by checking old names not in new map
		toRemove := []string{}
		for name, attach := range oldAttachMap {
			if _, exists := newAttachMap[name]; !exists {
				if id, ok := attach["id"].(string); ok {
					toRemove = append(toRemove, id)
				}
			}
		}

		// Process removals first
		instanceID := d.Id()
		for _, id := range toRemove {
			deleteOptions := &vpcv1.DeleteInstanceClusterNetworkAttachmentOptions{
				InstanceID: &instanceID,
				ID:         &id,
			}
			_, _, err := instanceC.DeleteInstanceClusterNetworkAttachment(deleteOptions)
			if err != nil {
				return fmt.Errorf("error removing cluster network attachment: %v", err)
			}
		}

		// Process additions
		for _, attach := range toAdd {
			createOptions := buildCreateClusterNetworkAttachmentOptions(d.Id(), attach)
			_, _, err := instanceC.CreateClusterNetworkAttachment(createOptions)
			if err != nil {
				return fmt.Errorf("error adding cluster network attachment: %v", err)
			}
		}

		// Identify and process updates for existing attachments
		for name, newAttach := range newAttachMap {
			if oldAttach, exists := oldAttachMap[name]; exists {
				// Compare the interfaces to see if an update is needed
				oldInterface := oldAttach["cluster_network_interface"].([]interface{})[0].(map[string]interface{})
				newInterface := newAttach["cluster_network_interface"].([]interface{})[0].(map[string]interface{})

				if !compareInterfaces(oldInterface, newInterface) {
					updateOptions := buildUpdateClusterNetworkAttachmentOptions(d.Id(), newAttach)
					_, _, err := instanceC.UpdateInstanceClusterNetworkAttachment(updateOptions)
					if err != nil {
						return fmt.Errorf("error updating cluster network attachment: %v", err)
					}
				}
			}
		}
	}
	return nil
}

// Helper function to compare interfaces
func compareInterfaces(old, new map[string]interface{}) bool {
	// Compare name
	if old["name"] != new["name"] {
		return false
	}

	// Compare auto_delete if present
	oldAutoDelete, oldOk := old["auto_delete"].(bool)
	newAutoDelete, newOk := new["auto_delete"].(bool)
	if oldOk != newOk || (oldOk && oldAutoDelete != newAutoDelete) {
		return false
	}

	// Compare primary_ip if present
	oldPrimaryIP, oldOk := old["primary_ip"].([]interface{})
	newPrimaryIP, newOk := new["primary_ip"].([]interface{})
	if oldOk != newOk {
		return false
	}
	if oldOk && newOk {
		if len(oldPrimaryIP) != len(newPrimaryIP) {
			return false
		}
		if len(oldPrimaryIP) > 0 && len(newPrimaryIP) > 0 {
			oldIP := oldPrimaryIP[0].(map[string]interface{})
			newIP := newPrimaryIP[0].(map[string]interface{})
			if oldIP["address"] != newIP["address"] ||
				oldIP["auto_delete"] != newIP["auto_delete"] ||
				oldIP["name"] != newIP["name"] {
				return false
			}
		}
	}

	// Compare subnet if present
	oldSubnet, oldOk := old["subnet"].([]interface{})
	newSubnet, newOk := new["subnet"].([]interface{})
	if oldOk != newOk {
		return false
	}
	if oldOk && newOk {
		if len(oldSubnet) != len(newSubnet) {
			return false
		}
		if len(oldSubnet) > 0 && len(newSubnet) > 0 {
			oldSub := oldSubnet[0].(map[string]interface{})
			newSub := newSubnet[0].(map[string]interface{})
			if oldSub["id"] != newSub["id"] {
				return false
			}
		}
	}

	return true
}

func buildCreateClusterNetworkAttachmentOptions(instanceID string, attachment map[string]interface{}) *vpcv1.CreateClusterNetworkAttachmentOptions {
	networkInterface := attachment["cluster_network_interface"].([]interface{})[0].(map[string]interface{})

	clusterNetworkInterface := &vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterface{}

	// if autoDelete, ok := networkInterface["auto_delete"].(bool); ok {
	// 	clusterNetworkInterface.AutoDelete = &autoDelete
	// }
	if autoDelete, ok := networkInterface["auto_delete"]; ok {
		// Convert interface{} to bool properly
		autoDeletBool := false
		switch v := autoDelete.(type) {
		case bool:
			autoDeletBool = v
		case string:
			autoDeletBool = v == "true"
		}
		clusterNetworkInterface.AutoDelete = &autoDeletBool
	}

	if name, ok := networkInterface["name"].(string); ok {
		clusterNetworkInterface.Name = &name
	}

	if primaryIPList, ok := networkInterface["primary_ip"].([]interface{}); ok && len(primaryIPList) > 0 {
		primaryIP := primaryIPList[0].(map[string]interface{})
		primaryIPPrototype := &vpcv1.ClusterNetworkInterfacePrimaryIPPrototype{}

		if address, ok := primaryIP["address"].(string); ok {
			primaryIPPrototype.Address = &address
		}
		if autoDelete, ok := primaryIP["auto_delete"].(bool); ok {
			primaryIPPrototype.AutoDelete = &autoDelete
		}
		if name, ok := primaryIP["name"].(string); ok {
			primaryIPPrototype.Name = &name
		}

		clusterNetworkInterface.PrimaryIP = primaryIPPrototype
	}

	// Handle subnet if present
	if subnetList, ok := networkInterface["subnet"].([]interface{}); ok && len(subnetList) > 0 {
		subnet := subnetList[0].(map[string]interface{})
		if id, ok := subnet["id"].(string); ok {
			clusterNetworkInterface.Subnet = &vpcv1.ClusterNetworkSubnetIdentity{
				ID: &id,
			}
		}
	}

	// Get attachment name
	attachmentName := attachment["name"].(string)

	// Create the options struct
	createOptions := &vpcv1.CreateClusterNetworkAttachmentOptions{
		InstanceID:              &instanceID,
		Name:                    &attachmentName,
		ClusterNetworkInterface: clusterNetworkInterface,
	}

	return createOptions
}

func buildUpdateClusterNetworkAttachmentOptions(instanceID string, attachment map[string]interface{}) *vpcv1.UpdateInstanceClusterNetworkAttachmentOptions {
	networkInterface := attachment["cluster_network_interface"].([]interface{})[0].(map[string]interface{})
	clusterNetworkInterface := &vpcv1.InstanceClusterNetworkAttachmentPatch{}

	if name, ok := networkInterface["name"].(string); ok {
		clusterNetworkInterface.Name = &name
	}
	clusterNetworkInterfaceAsPatch, _ := clusterNetworkInterface.AsPatch()
	attachmentID := attachment["id"].(string)
	updateOptions := &vpcv1.UpdateInstanceClusterNetworkAttachmentOptions{
		InstanceID:                            &instanceID,
		ID:                                    &attachmentID,
		InstanceClusterNetworkAttachmentPatch: clusterNetworkInterfaceAsPatch,
	}
	return updateOptions
}
