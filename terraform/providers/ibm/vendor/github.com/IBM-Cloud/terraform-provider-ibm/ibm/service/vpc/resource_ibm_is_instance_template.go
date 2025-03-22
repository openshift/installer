// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isInstanceTemplateBootVolume                   = "boot_volume"
	isInstanceTemplateBootVolumeTags               = "tags"
	isInstanceTemplateVolAttTags                   = "tags"
	isInstanceTemplateCRN                          = "crn"
	isInstanceTemplateVolAttVolAutoDelete          = "auto_delete"
	isInstanceTemplateVolAttVol                    = "volume"
	isInstanceTemplateVolAttachmentName            = "name"
	isInstanceTemplateVolAttVolPrototype           = "volume_prototype"
	isInstanceTemplateVolAttVolCapacity            = "capacity"
	isInstanceTemplateVolAttVolIops                = "iops"
	isInstanceTemplateVolAttVolName                = "name"
	isInstanceTemplateVolAttVolBillingTerm         = "billing_term"
	isInstanceTemplateVolAttVolEncryptionKey       = "encryption_key"
	isInstanceTemplateVolAttVolType                = "type"
	isInstanceTemplateVolAttVolProfile             = "profile"
	isInstanceTemplateProvisioning                 = "provisioning"
	isInstanceTemplateProvisioningDone             = "done"
	isInstanceTemplateAvailable                    = "available"
	isInstanceTemplateDeleting                     = "deleting"
	isInstanceTemplateDeleteDone                   = "done"
	isInstanceTemplateFailed                       = "failed"
	isInstanceTemplateBootName                     = "name"
	isInstanceTemplateBootSize                     = "size"
	isInstanceTemplateBootIOPS                     = "iops"
	isInstanceTemplateBootEncryption               = "encryption"
	isInstanceTemplateBootProfile                  = "profile"
	isInstanceTemplateVolumeAttaching              = "attaching"
	isInstanceTemplateVolumeAttached               = "attached"
	isInstanceTemplateVolumeDetaching              = "detaching"
	isInstanceTemplatePlacementTarget              = "placement_target"
	isInstanceTemplateDedicatedHost                = "dedicated_host"
	isInstanceTemplateDedicatedHostGroup           = "dedicated_host_group"
	isInstanceTemplateResourceType                 = "resource_type"
	isInstanceTemplateVolumeDeleteOnInstanceDelete = "delete_volume_on_instance_delete"
	isInstanceTemplateMetadataServiceEnabled       = "metadata_service_enabled"
	isInstanceTemplateAvailablePolicyHostFailure   = "availability_policy_host_failure"
	isInstanceTemplateHostFailure                  = "host_failure"
	isInstanceTemplateNicPrimaryIP                 = "primary_ip"
	isInstanceTemplateNicReservedIpAddress         = "address"
	isInstanceTemplateNicReservedIpAutoDelete      = "auto_delete"
	isInstanceTemplateNicReservedIpName            = "name"
	isInstanceTemplateNicReservedIpId              = "reserved_ip"

	// catalog offering
	isInstanceTemplateCatalogOffering            = "catalog_offering"
	isInstanceTemplateCatalogOfferingOfferingCrn = "offering_crn"
	isInstanceTemplateCatalogOfferingVersionCrn  = "version_crn"
	isInstanceTemplateCatalogOfferingPlanCrn     = "plan_crn"
)

func ResourceIBMISInstanceTemplate() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMisInstanceTemplateCreate,
		Read:     resourceIBMisInstanceTemplateRead,
		Update:   resourceIBMisInstanceTemplateUpdate,
		Delete:   resourceIBMisInstanceTemplateDelete,
		Exists:   resourceIBMisInstanceTemplateExists,
		Importer: &schema.ResourceImporter{},

		CustomizeDiff: customdiff.All(
			customdiff.Sequence(
				func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
					return flex.ResourceTagsCustomizeDiff(diff)
				},
			),

			customdiff.Sequence(
				func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
					return flex.ResourceVolumeAttachmentValidate(diff)
				}),
		),

		Schema: map[string]*schema.Schema{
			isInstanceTemplateAvailablePolicyHostFailure: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The availability policy to use for this virtual server instance",
			},

			isInstanceTemplateName: {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     false,
				ValidateFunc: validate.ValidateISName,
				Description:  "Instance Template name",
			},

			// cluster changes
			"cluster_network_attachments": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "The cluster network attachments to create for this virtual server instance. A cluster network attachment represents a device that is connected to a cluster network. The number of network attachments must match one of the values from the instance profile's `cluster_network_attachment_count` before the instance can be started.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_network_interface": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Required:    true,
							ForceNew:    true,
							Description: "A cluster network interface for the instance cluster network attachment. This can bespecified using an existing cluster network interface that does not already have a `target`,or a prototype object for a new cluster network interface.This instance must reside in the same VPC as the specified cluster network interface. Thecluster network interface must reside in the same cluster network as the`cluster_network_interface` of any other `cluster_network_attachments` for this instance.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"auto_delete": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Computed:    true,
										Description: "Indicates whether this cluster network interface will be automatically deleted when `target` is deleted.",
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
													Description: "The unique identifier for this cluster network subnet reserved IP.",
												},
												"href": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The URL for this cluster network subnet reserved IP.",
												},
												"address": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The IP address to reserve, which must not already be reserved on the subnet.If unspecified, an available address on the subnet will automatically be selected.",
												},
												"auto_delete": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Default:     true,
													Description: "Indicates whether this cluster network subnet reserved IP member will be automatically deleted when either `target` is deleted, or the cluster network subnet reserved IP is unbound.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
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
													Description: "The unique identifier for this cluster network subnet.",
												},
												"href": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The URL for this cluster network subnet.",
												},
											},
										},
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The unique identifier for this cluster network interface.",
									},
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The URL for this cluster network interface.",
									},
								},
							},
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The name for this cluster network attachment. Names must be unique within the instance the cluster network attachment resides in. If unspecified, the name will be a hyphenated list of randomly-selected words. Names starting with `ibm-` are reserved for provider-owned resources, and are not allowed.",
						},
					},
				},
			},

			"confidential_compute_mode": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_instance_template", "confidential_compute_mode"),
				Description:  "The confidential compute mode to use for this virtual server instance.If unspecified, the default confidential compute mode from the profile will be used.",
			},
			"enable_secure_boot": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Indicates whether secure boot is enabled for this virtual server instance.If unspecified, the default secure boot mode from the profile will be used.",
			},

			isInstanceTemplateMetadataServiceEnabled: {
				Type:          schema.TypeBool,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				Deprecated:    "Use metadata_service instead",
				ConflictsWith: []string{isInstanceMetadataService},
				Description:   "Indicates whether the metadata service endpoint is available to the virtual server instance",
			},

			isInstanceMetadataService: {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				MinItems:      1,
				MaxItems:      1,
				ConflictsWith: []string{isInstanceTemplateMetadataServiceEnabled},
				Description:   "The metadata service configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isInstanceMetadataServiceEnabled1: {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    true,
							Computed:    true,
							Description: "Indicates whether the metadata service endpoint will be available to the virtual server instance",
						},

						isInstanceMetadataServiceProtocol: {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Computed:    true,
							Description: "The communication protocol to use for the metadata service endpoint. Applies only when the metadata service is enabled.",
						},

						isInstanceMetadataServiceRespHopLimit: {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Computed:    true,
							Description: "The hop limit (IP time to live) for IP response packets from the metadata service",
						},
					},
				},
			},
			isInstanceTemplateVPC: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "VPC id",
			},

			isInstanceTemplateZone: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Zone name",
			},

			isInstanceTemplateProfile: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
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

			isInstanceTotalVolumeBandwidth: {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_instance_template", isInstanceTotalVolumeBandwidth),
				Description:  "The amount of bandwidth (in megabits per second) allocated exclusively to instance storage volumes",
			},

			isInstanceTemplateKeys: {
				Type:             schema.TypeSet,
				Required:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				Set:              schema.HashString,
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "SSH key Ids for the instance template",
			},

			isPlacementTargetDedicatedHost: {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{isPlacementTargetDedicatedHostGroup, isPlacementTargetPlacementGroup},
				Description:   "Unique Identifier of the Dedicated Host where the instance will be placed",
			},

			isPlacementTargetDedicatedHostGroup: {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
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

			isInstanceTemplateVolumeAttachments: {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isInstanceTemplateVolumeDeleteOnInstanceDelete: {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "If set to true, when deleting the instance the volume will also be deleted.",
						},
						isInstanceTemplateVolAttachmentName: {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.InvokeValidator("ibm_is_instance_template", isInstanceTemplateVolAttachmentName),
							Description:  "The user-defined name for this volume attachment.",
						},
						isInstanceTemplateVolAttVol: {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The unique identifier for this volume.",
						},
						isInstanceTemplateVolAttVolPrototype: {
							Type:     schema.TypeList,
							MaxItems: 1,
							MinItems: 1,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isInstanceTemplateVolAttVolIops: {
										Type:        schema.TypeInt,
										Optional:    true,
										ForceNew:    true,
										Description: "The maximum I/O operations per second (IOPS) for the volume.",
									},
									isInstanceTemplateVolAttVolProfile: {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: "The  globally unique name for the volume profile to use for this volume.",
									},
									isInstanceTemplateVolAttVolCapacity: {
										Type:        schema.TypeInt,
										Required:    true,
										ForceNew:    true,
										Description: "The capacity of the volume in gigabytes. The specified minimum and maximum capacity values for creating or updating volumes may expand in the future.",
									},
									isInstanceTemplateVolAttVolEncryptionKey: {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "The CRN of the [Key Protect Root Key](https://cloud.ibm.com/docs/key-protect?topic=key-protect-getting-started-tutorial) or [Hyper Protect Crypto Service Root Key](https://cloud.ibm.com/docs/hs-crypto?topic=hs-crypto-get-started) for this resource.",
									},
									isInstanceTemplateVolAttTags: {
										Type:        schema.TypeSet,
										Optional:    true,
										ForceNew:    true,
										Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_instance_template", "tags")},
										Set:         flex.ResourceIBMVPCHash,
										Description: "UserTags for the volume instance",
									},
								},
							},
						},
					},
				},
			},

			isInstanceTemplateCatalogOffering: {
				Type:         schema.TypeList,
				MinItems:     0,
				MaxItems:     1,
				ExactlyOneOf: []string{isInstanceTemplateCatalogOffering, isInstanceTemplateImage},
				Optional:     true,
				ForceNew:     true,
				Description:  "The catalog offering or offering version to use when provisioning this virtual server instance template. If an offering is specified, the latest version of that offering will be used. The specified offering or offering version may be in a different account in the same enterprise, subject to IAM policies.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isInstanceTemplateCatalogOfferingOfferingCrn: {
							Type:          schema.TypeString,
							Optional:      true,
							ForceNew:      true,
							ConflictsWith: []string{"catalog_offering.0.version_crn"},
							Description:   "Identifies a catalog offering by a unique CRN property",
						},
						isInstanceTemplateCatalogOfferingVersionCrn: {
							Type:          schema.TypeString,
							Optional:      true,
							ForceNew:      true,
							ConflictsWith: []string{"catalog_offering.0.offering_crn"},
							Description:   "Identifies a version of a catalog offering by a unique CRN property",
						},
						isInstanceTemplateCatalogOfferingPlanCrn: {
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

			isInstanceTemplatePrimaryNetworkInterface: {
				Type:          schema.TypeList,
				MinItems:      1,
				MaxItems:      1,
				Optional:      true,
				ExactlyOneOf:  []string{"primary_network_attachment", "primary_network_interface"},
				ConflictsWith: []string{"primary_network_attachment", "network_attachments"},
				Description:   "Primary Network interface info",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isInstanceTemplateNicAllowIPSpoofing: {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						isInstanceTemplateNicName: {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						isInstanceTemplateNicPrimaryIpv4Address: {
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							ConflictsWith: []string{"primary_network_interface.0.primary_ip.0.address"},
							Deprecated:    "primary_ipv4_address is deprecated and support will be removed. Use primary_ip instead",
						},
						isInstanceTemplateNicPrimaryIP: {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: "The primary IP address to bind to the network interface. This can be specified using an existing reserved IP, or a prototype object for a new reserved IP.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isInstanceTemplateNicReservedIpAddress: {
										Type:          schema.TypeString,
										Optional:      true,
										ConflictsWith: []string{"primary_network_interface.0.primary_ipv4_address"},
										Computed:      true,
										ForceNew:      true,
										Description:   "The IP address to reserve, which must not already be reserved on the subnet.",
									},
									isInstanceTemplateNicReservedIpAutoDelete: {
										Type:        schema.TypeBool,
										Optional:    true,
										Computed:    true,
										ForceNew:    true,
										Description: "Indicates whether this reserved IP member will be automatically deleted when either target is deleted, or the reserved IP is unbound.",
									},
									isInstanceTemplateNicReservedIpName: {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										ForceNew:    true,
										Description: "The user-defined name for this reserved IP. If unspecified, the name will be a hyphenated list of randomly-selected words. Names must be unique within the subnet the reserved IP resides in. ",
									},
									isInstanceTemplateNicReservedIpId: {
										Type:          schema.TypeString,
										Optional:      true,
										Computed:      true,
										ForceNew:      true,
										ConflictsWith: []string{"primary_network_interface.0.primary_ip.0.address", "primary_network_interface.0.primary_ip.0.auto_delete", "primary_network_interface.0.primary_ip.0.name", "primary_network_interface.0.primary_ipv4_address"},
										Description:   "Identifies a reserved IP by a unique property.",
									},
								},
							},
						},
						isInstanceTemplateNicSecurityGroups: {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						isInstanceTemplateNicSubnet: {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			isInstanceTemplateNetworkInterfaces: {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"primary_network_attachment", "network_attachments"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isInstanceTemplateNicAllowIPSpoofing: {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						isInstanceTemplateNicName: {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						isInstanceTemplateNicPrimaryIpv4Address: {
							Type:       schema.TypeString,
							Optional:   true,
							Deprecated: "primary_ipv4_address is deprecated and support will be removed. Use primary_ip instead",
							Computed:   true,
						},
						isInstanceTemplateNicPrimaryIP: {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Computed:    true,
							Description: "The primary IP address to bind to the network interface. This can be specified using an existing reserved IP, or a prototype object for a new reserved IP.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isInstanceTemplateNicReservedIpAddress: {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Computed:    true,
										Description: "The IP address to reserve, which must not already be reserved on the subnet.",
									},
									isInstanceTemplateNicReservedIpAutoDelete: {
										Type:        schema.TypeBool,
										Optional:    true,
										Computed:    true,
										ForceNew:    true,
										Description: "Indicates whether this reserved IP member will be automatically deleted when either target is deleted, or the reserved IP is unbound.",
									},
									isInstanceTemplateNicReservedIpName: {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										ForceNew:    true,
										Description: "The user-defined name for this reserved IP. If unspecified, the name will be a hyphenated list of randomly-selected words. Names must be unique within the subnet the reserved IP resides in. ",
									},
									isInstanceTemplateNicReservedIpId: {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										ForceNew:    true,
										Description: "Identifies a reserved IP by a unique property.",
									},
								},
							},
						},
						isInstanceTemplateNicSecurityGroups: {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						isInstanceTemplateNicSubnet: {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			// vni

			"primary_network_attachment": &schema.Schema{
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ForceNew:      true,
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
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "The reserved IPs bound to this virtual network interface.May be empty when `lifecycle_state` is `pending`.",
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
												"auto_delete": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
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
									"name": &schema.Schema{
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validate.InvokeValidator("ibm_is_virtual_network_interface", "vni_name"),
										Description:  "The name for this virtual network interface. The name is unique across all virtual network interfaces in the VPC.",
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
									"protocol_state_filtering_mode": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validate.InvokeValidator("ibm_is_virtual_network_interface", "protocol_state_filtering_mode"),
										Description:  "The protocol state filtering mode used for this virtual network interface.",
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
										ForceNew:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Set:         schema.HashString,
										Description: "The security groups for this virtual network interface.",
									},
									"subnet": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										ForceNew:    true,
										Description: "The associated subnet id.",
									},
								},
							},
						},
					},
				},
			},

			"network_attachments": &schema.Schema{
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
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
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "The reserved IPs bound to this virtual network interface.May be empty when `lifecycle_state` is `pending`.",
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
												"auto_delete": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
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
									"name": &schema.Schema{
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validate.InvokeValidator("ibm_is_virtual_network_interface", "vni_name"),
										Description:  "The name for this virtual network interface. The name is unique across all virtual network interfaces in the VPC.",
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
									"protocol_state_filtering_mode": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validate.InvokeValidator("ibm_is_virtual_network_interface", "protocol_state_filtering_mode"),
										Description:  "The protocol state filtering mode used for this virtual network interface.",
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

			isInstanceTemplateUserData: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "User data given for the instance",
			},

			isInstanceTemplateCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for the instance",
			},

			isInstanceTemplateImage: {
				Type:         schema.TypeString,
				ForceNew:     true,
				ExactlyOneOf: []string{isInstanceTemplateCatalogOffering, isInstanceTemplateImage},
				Optional:     true,
				Description:  "image name",
			},

			isInstanceTemplateBootVolume: {
				Type:             schema.TypeList,
				DiffSuppressFunc: flex.ApplyOnce,
				Optional:         true,
				Computed:         true,
				MaxItems:         1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isInstanceTemplateBootName: {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						isInstanceTemplateBootVolumeTags: {
							Type:        schema.TypeSet,
							Optional:    true,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_instance_template", "tags")},
							Set:         flex.ResourceIBMVPCHash,
							Description: "UserTags for the volume instance",
						},
						isInstanceTemplateBootEncryption: {
							Type:     schema.TypeString,
							Optional: true,
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
						isInstanceTemplateVolumeDeleteOnInstanceDelete: {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			isInstanceTemplateResourceGroup: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: "Instance template resource group",
			},

			isInstanceTemplatePlacementTarget: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The placement restrictions for the virtual server instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this placement target.",
						},
						"crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this placement target.",
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this placement target.",
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
						isReservationAffinityPolicyResp: {
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
								},
							},
						},
					},
				},
			},
		},
	}
}

func ResourceIBMISInstanceTemplateValidator() *validate.ResourceValidator {
	host_failure := "restart, stop"
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
			Identifier:                 isInstanceTemplateVolAttachmentName,
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isInstanceTotalVolumeBandwidth,
			ValidateFunctionIdentifier: validate.IntAtLeast,
			Type:                       validate.TypeInt,
			Optional:                   true,
			MinValue:                   "500"})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isInstanceTemplateAvailablePolicyHostFailure,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              host_failure})

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "tags",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9:_ .-]+$`,
			MinValueLength:             1,
			MaxValueLength:             128})
	ibmISInstanceTemplateValidator := validate.ResourceValidator{ResourceName: "ibm_is_instance_template", Schema: validateSchema}
	return &ibmISInstanceTemplateValidator
}

func resourceIBMisInstanceTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	profile := d.Get(isInstanceTemplateProfile).(string)
	name := d.Get(isInstanceTemplateName).(string)
	vpcID := d.Get(isInstanceTemplateVPC).(string)
	zone := d.Get(isInstanceTemplateZone).(string)
	image := d.Get(isInstanceTemplateImage).(string)

	if catalogOfferingOk, ok := d.GetOk(isInstanceTemplateCatalogOffering); ok {
		catalogOffering := catalogOfferingOk.([]interface{})[0].(map[string]interface{})
		offeringCrn, _ := catalogOffering[isInstanceTemplateCatalogOfferingOfferingCrn].(string)
		versionCrn, _ := catalogOffering[isInstanceTemplateCatalogOfferingVersionCrn].(string)
		planCrn, _ := catalogOffering[isInstanceTemplateCatalogOfferingPlanCrn].(string)
		err := instanceTemplateCreateByCatalogOffering(d, meta, profile, name, vpcID, zone, offeringCrn, versionCrn, planCrn)
		if err != nil {
			return err
		}
	} else {
		err := instanceTemplateCreate(d, meta, profile, name, vpcID, zone, image)
		if err != nil {
			return err
		}
	}

	return resourceIBMisInstanceTemplateRead(d, meta)
}

func resourceIBMisInstanceTemplateRead(d *schema.ResourceData, meta interface{}) error {
	ID := d.Id()
	err := instanceTemplateGet(d, meta, ID)
	if err != nil {
		return err
	}
	return nil
}

func resourceIBMisInstanceTemplateDelete(d *schema.ResourceData, meta interface{}) error {

	ID := d.Id()

	err := instanceTemplateDelete(d, meta, ID)
	if err != nil {
		return err
	}
	return nil
}

func resourceIBMisInstanceTemplateUpdate(d *schema.ResourceData, meta interface{}) error {

	err := instanceTemplateUpdate(d, meta)
	if err != nil {
		return err
	}
	return resourceIBMisInstanceTemplateRead(d, meta)
}

func resourceIBMisInstanceTemplateExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	ID := d.Id()
	ok, err := instanceTemplateExists(d, meta, ID)
	if err != nil {
		return false, err
	}
	return ok, err
}

func instanceTemplateCreateByCatalogOffering(d *schema.ResourceData, meta interface{}, profile, name, vpcID, zone, offeringCrn, versionCrn, planCrn string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	instanceproto := &vpcv1.InstanceTemplatePrototypeInstanceTemplateByCatalogOffering{
		Zone: &vpcv1.ZoneIdentity{
			Name: &zone,
		},
		Profile: &vpcv1.InstanceProfileIdentity{
			Name: &profile,
		},
		VPC: &vpcv1.VPCIdentity{
			ID: &vpcID,
		},
	}
	// cluster changes
	if clusterNetworkAttachmentOk, ok := d.GetOk("cluster_network_attachments"); ok {
		clusterNetworkAttachmentList := clusterNetworkAttachmentOk.([]interface{})
		clusterNetworkAttachments := []vpcv1.InstanceClusterNetworkAttachmentPrototypeInstanceContext{}
		for _, clusterNetworkAttachmentsItem := range clusterNetworkAttachmentList {
			clusterNetworkAttachmentsItemModel, err := ResourceIBMIsInstanceTemplateMapToInstanceClusterNetworkAttachmentPrototypeInstanceContext(clusterNetworkAttachmentsItem.(map[string]interface{}))
			if err != nil {
				return err
			}
			clusterNetworkAttachments = append(clusterNetworkAttachments, *clusterNetworkAttachmentsItemModel)
		}
		instanceproto.ClusterNetworkAttachments = clusterNetworkAttachments
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
	if offeringCrn != "" {
		catalogOffering := &vpcv1.InstanceCatalogOfferingPrototypeCatalogOfferingByOffering{
			Offering: &vpcv1.CatalogOfferingIdentityCatalogOfferingByCRN{
				CRN: &offeringCrn,
			},
		}
		if planOffering != nil {
			catalogOffering.Plan = planOffering
		}
		instanceproto.CatalogOffering = catalogOffering
	}
	if versionCrn != "" {
		catalogOffering := &vpcv1.InstanceCatalogOfferingPrototypeCatalogOfferingByVersion{
			Version: &vpcv1.CatalogOfferingVersionIdentityCatalogOfferingVersionByCRN{
				CRN: &versionCrn,
			},
		}
		if planOffering != nil {
			catalogOffering.Plan = planOffering
		}
		instanceproto.CatalogOffering = catalogOffering
	}

	if name != "" {
		instanceproto.Name = &name
	}

	// vni
	if networkattachmentsintf, ok := d.GetOk("network_attachments"); ok {
		networkAttachments := []vpcv1.InstanceNetworkAttachmentPrototype{}
		for i, networkAttachmentsItem := range networkattachmentsintf.([]interface{}) {
			allowipspoofing := fmt.Sprintf("network_attachments.%d.allow_ip_spoofing", i)
			autodelete := fmt.Sprintf("network_attachments.%d.auto_delete", i)
			enablenat := fmt.Sprintf("network_attachments.%d.enable_infrastructure_nat", i)
			networkAttachmentsItemModel, err := resourceIBMIsInstanceTemplateMapToInstanceNetworkAttachmentPrototype(allowipspoofing, autodelete, enablenat, d, networkAttachmentsItem.(map[string]interface{}))
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
		primaryNetworkAttachmentModel, err := resourceIBMIsInstanceTemplateMapToInstanceNetworkAttachmentPrototype(allowipspoofing, autodelete, enablenat, d, primnetworkattachmentintf.([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return err
		}
		instanceproto.PrimaryNetworkAttachment = primaryNetworkAttachmentModel
	}

	metadataServiceEnabled := d.Get(isInstanceTemplateMetadataServiceEnabled).(bool)
	if metadataServiceEnabled {
		instanceproto.MetadataService = &vpcv1.InstanceMetadataServicePrototype{
			Enabled: &metadataServiceEnabled,
		}
	}

	if metadataService := GetInstanceTemplateMetadataServiceOptions(d); metadataService != nil {
		instanceproto.MetadataService = metadataService
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
	if availablePolicyHostFailureIntf, ok := d.GetOk(isInstanceTemplateAvailablePolicyHostFailure); ok {
		availablePolicyHostFailure := availablePolicyHostFailureIntf.(string)
		instanceproto.AvailabilityPolicy = &vpcv1.InstanceAvailabilityPolicyPrototype{
			HostFailure: &availablePolicyHostFailure,
		}
	}
	if dHostIdInf, ok := d.GetOk(isPlacementTargetDedicatedHost); ok {
		dHostIdStr := dHostIdInf.(string)
		dHostPlaementTarget := &vpcv1.InstancePlacementTargetPrototypeDedicatedHostIdentity{
			ID: &dHostIdStr,
		}
		instanceproto.PlacementTarget = dHostPlaementTarget
	}

	if dHostGrpIdInf, ok := d.GetOk(isPlacementTargetDedicatedHostGroup); ok {
		dHostGrpIdStr := dHostGrpIdInf.(string)
		dHostGrpPlaementTarget := &vpcv1.InstancePlacementTargetPrototypeDedicatedHostGroupIdentity{
			ID: &dHostGrpIdStr,
		}
		instanceproto.PlacementTarget = dHostGrpPlaementTarget
	}

	if placementGroupInf, ok := d.GetOk(isPlacementTargetPlacementGroup); ok {
		placementGrpStr := placementGroupInf.(string)
		placementGrp := &vpcv1.InstancePlacementTargetPrototypePlacementGroupIdentity{
			ID: &placementGrpStr,
		}
		instanceproto.PlacementTarget = placementGrp
	}

	if totalVolBandwidthIntf, ok := d.GetOk(isInstanceTotalVolumeBandwidth); ok {
		totalVolBandwidthStr := int64(totalVolBandwidthIntf.(int))
		instanceproto.TotalVolumeBandwidth = &totalVolBandwidthStr
	}

	// BOOT VOLUME ATTACHMENT for instance template
	if boot, ok := d.GetOk(isInstanceTemplateBootVolume); ok {
		bootvol := boot.([]interface{})[0].(map[string]interface{})
		var volTemplate = &vpcv1.VolumePrototypeInstanceByImageContext{}
		name, ok := bootvol[isInstanceTemplateBootName]
		namestr := name.(string)
		if ok && namestr != "" {
			volTemplate.Name = &namestr
		}

		var userTags *schema.Set
		if v, ok := bootvol[isInstanceTemplateBootVolumeTags]; ok {
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

		volcap := 100
		volcapint64 := int64(volcap)
		volprof := "general-purpose"
		volTemplate.Capacity = &volcapint64
		volTemplate.Profile = &vpcv1.VolumeProfileIdentity{
			Name: &volprof,
		}

		if encryption, ok := bootvol[isInstanceTemplateBootEncryption]; ok {
			bootEncryption := encryption.(string)
			if bootEncryption != "" {
				volTemplate.EncryptionKey = &vpcv1.EncryptionKeyIdentity{
					CRN: &bootEncryption,
				}
			}
		}

		var deleteVolumeOption bool
		if deleteVolume, ok := bootvol[isInstanceTemplateVolumeDeleteOnInstanceDelete]; ok {
			deleteVolumeOption = deleteVolume.(bool)
		}

		instanceproto.BootVolumeAttachment = &vpcv1.VolumeAttachmentPrototypeInstanceByImageContext{
			DeleteVolumeOnInstanceDelete: &deleteVolumeOption,
			Volume:                       volTemplate,
		}
	}

	// Handle volume attachments
	if volsintf, ok := d.GetOk(isInstanceTemplateVolumeAttachments); ok {
		vols := volsintf.([]interface{})
		var intfs []vpcv1.VolumeAttachmentPrototype
		for _, resource := range vols {
			vol := resource.(map[string]interface{})
			volInterface := &vpcv1.VolumeAttachmentPrototype{}
			deleteVolBool := vol[isInstanceTemplateVolumeDeleteOnInstanceDelete].(bool)
			volInterface.DeleteVolumeOnInstanceDelete = &deleteVolBool
			attachmentnamestr := vol[isInstanceTemplateVolAttachmentName].(string)
			volInterface.Name = &attachmentnamestr
			volIdStr := vol[isInstanceTemplateVolAttVol].(string)

			if volIdStr != "" {
				volInterface.Volume = &vpcv1.VolumeAttachmentPrototypeVolume{
					ID: &volIdStr,
				}
			} else {
				newvolintf := vol[isInstanceTemplateVolAttVolPrototype].([]interface{})[0]
				newvol := newvolintf.(map[string]interface{})
				profileName := newvol[isInstanceTemplateVolAttVolProfile].(string)
				capacity := int64(newvol[isInstanceTemplateVolAttVolCapacity].(int))

				volPrototype := &vpcv1.VolumeAttachmentPrototypeVolumeVolumePrototypeInstanceContext{
					Profile: &vpcv1.VolumeProfileIdentity{
						Name: &profileName,
					},
					Capacity: &capacity,
				}
				iops := int64(newvol[isInstanceTemplateVolAttVolIops].(int))
				encryptionKey := newvol[isInstanceTemplateVolAttVolEncryptionKey].(string)

				if iops != 0 {
					volPrototype.Iops = &iops
				}

				if encryptionKey != "" {
					volPrototype.EncryptionKey = &vpcv1.EncryptionKeyIdentity{
						CRN: &encryptionKey,
					}
				}
				var userTags *schema.Set
				if v, ok := newvol[isInstanceTemplateVolAttTags]; ok {
					userTags = v.(*schema.Set)
					if userTags != nil && userTags.Len() != 0 {
						userTagsArray := make([]string, userTags.Len())
						for i, userTag := range userTags.List() {
							userTagStr := userTag.(string)
							userTagsArray[i] = userTagStr
						}
						volPrototype.UserTags = userTagsArray
					}
				}
				volInterface.Volume = volPrototype
			}

			intfs = append(intfs, *volInterface)
		}
		instanceproto.VolumeAttachments = intfs
	}

	// Handle primary network interface
	if primnicintf, ok := d.GetOk(isInstanceTemplatePrimaryNetworkInterface); ok {
		primnic := primnicintf.([]interface{})[0].(map[string]interface{})
		subnetintf, _ := primnic[isInstanceTemplateNicSubnet]
		subnetintfstr := subnetintf.(string)
		var primnicobj = &vpcv1.NetworkInterfacePrototype{}
		primnicobj.Subnet = &vpcv1.SubnetIdentity{
			ID: &subnetintfstr,
		}

		if name, ok := primnic[isInstanceTemplateNicName]; ok {
			namestr := name.(string)
			if namestr != "" {
				primnicobj.Name = &namestr
			}
		}
		allowIPSpoofing, ok := primnic[isInstanceTemplateNicAllowIPSpoofing]
		allowIPSpoofingbool := allowIPSpoofing.(bool)
		if ok {
			primnicobj.AllowIPSpoofing = &allowIPSpoofingbool
		}

		secgrpintf, ok := primnic[isInstanceTemplateNicSecurityGroups]
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
		// reserved ip changes
		var PrimaryIpv4Address, reservedIp, reservedIpAddress, reservedIpName string
		var reservedIpAutoDelete, okAuto bool
		if IPAddress, ok := primnic[isInstanceTemplateNicPrimaryIpv4Address]; ok {
			PrimaryIpv4Address = IPAddress.(string)
		}
		primaryIpOk, ok := primnic[isInstanceTemplateNicPrimaryIP]
		if ok && len(primaryIpOk.([]interface{})) > 0 {
			primip := primaryIpOk.([]interface{})[0].(map[string]interface{})

			reservedipok, _ := primip[isInstanceTemplateNicReservedIpId]
			reservedIp = reservedipok.(string)
			reservedipv4Ok, _ := primip[isInstanceTemplateNicReservedIpAddress]
			reservedIpAddress = reservedipv4Ok.(string)
			reservedipnameOk, _ := primip[isInstanceTemplateNicReservedIpName]
			reservedIpName = reservedipnameOk.(string)
			reservedipautodeleteok, okAuto := primip[isInstanceTemplateNicReservedIpAutoDelete]
			if okAuto {
				reservedIpAutoDelete = reservedipautodeleteok.(bool)
			}
		}
		if PrimaryIpv4Address != "" && reservedIpAddress != "" && PrimaryIpv4Address != reservedIpAddress {
			return fmt.Errorf("[ERROR] Error creating instance template, primary_network_interface error, use either primary_ipv4_address(%s) or primary_ip.0.address(%s)", PrimaryIpv4Address, reservedIpAddress)
		}
		if reservedIp != "" {
			primnicobj.PrimaryIP = &vpcv1.NetworkInterfaceIPPrototypeReservedIPIdentity{
				ID: &reservedIp,
			}
		} else {
			if PrimaryIpv4Address != "" || reservedIpAddress != "" || reservedIpName != "" || okAuto {
				primaryipobj := &vpcv1.NetworkInterfaceIPPrototypeReservedIPPrototypeNetworkInterfaceContext{}
				if PrimaryIpv4Address != "" {
					primaryipobj.Address = &PrimaryIpv4Address
				}
				if reservedIpAddress != "" {
					primaryipobj.Address = &reservedIpAddress
				}
				if reservedIpName != "" {
					primaryipobj.Name = &reservedIpName
				}
				if okAuto {
					primaryipobj.AutoDelete = &reservedIpAutoDelete
				}
				primnicobj.PrimaryIP = primaryipobj
			}
		}
		instanceproto.PrimaryNetworkInterface = primnicobj
	}

	// Handle  additional network interface
	if nicsintf, ok := d.GetOk(isInstanceTemplateNetworkInterfaces); ok {
		nics := nicsintf.([]interface{})
		var intfs []vpcv1.NetworkInterfacePrototype
		for _, resource := range nics {
			nic := resource.(map[string]interface{})
			nwInterface := &vpcv1.NetworkInterfacePrototype{}
			subnetintf, _ := nic[isInstanceTemplateNicSubnet]
			subnetintfstr := subnetintf.(string)
			nwInterface.Subnet = &vpcv1.SubnetIdentity{
				ID: &subnetintfstr,
			}

			name, ok := nic[isInstanceTemplateNicName]
			namestr := name.(string)
			if ok && namestr != "" {
				nwInterface.Name = &namestr
			}
			allowIPSpoofing, ok := nic[isInstanceTemplateNicAllowIPSpoofing]
			allowIPSpoofingbool := allowIPSpoofing.(bool)
			if ok {
				nwInterface.AllowIPSpoofing = &allowIPSpoofingbool
			}
			secgrpintf, ok := nic[isInstanceTemplateNicSecurityGroups]
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
			// reserved ip changes
			var PrimaryIpv4Address, reservedIp, reservedIpAddress, reservedIpName string
			var reservedIpAutoDelete, okAuto bool
			if IPAddress, ok := nic[isInstanceTemplateNicPrimaryIpv4Address]; ok {
				PrimaryIpv4Address = IPAddress.(string)
			}
			primaryIpOk, ok := nic[isInstanceTemplateNicPrimaryIP]
			if ok && len(primaryIpOk.([]interface{})) > 0 {
				primip := primaryIpOk.([]interface{})[0].(map[string]interface{})

				reservedipok, _ := primip[isInstanceTemplateNicReservedIpId]
				reservedIp = reservedipok.(string)
				reservedipv4Ok, _ := primip[isInstanceTemplateNicReservedIpAddress]
				reservedIpAddress = reservedipv4Ok.(string)
				reservedipnameOk, _ := primip[isInstanceTemplateNicReservedIpName]
				reservedIpName = reservedipnameOk.(string)
				// var reservedipautodeleteok interface{}

				if v, ok := primip[isInstanceTemplateNicReservedIpAutoDelete].(bool); ok && v {
					log.Printf("[INFO] isInstanceTemplateNicReservedIpAutoDelete is v is %t and okay is %t", v, ok)
					reservedIpAutoDelete = primip[isInstanceTemplateNicReservedIpAutoDelete].(bool)
					okAuto = true
				}
				// reservedipautodeleteok, okAuto = primip[isInstanceTemplateNicReservedIpAutoDelete]
				// if okAuto {
				// 	reservedIpAutoDelete = reservedipautodeleteok.(bool)
				// }
			}
			if PrimaryIpv4Address != "" && reservedIpAddress != "" && PrimaryIpv4Address != reservedIpAddress {
				return fmt.Errorf("[ERROR] Error creating instance template, network_interfaces error, use either primary_ipv4_address(%s) or primary_ip.0.address(%s)", PrimaryIpv4Address, reservedIpAddress)
			}
			if reservedIp != "" && (PrimaryIpv4Address != "" || reservedIpAddress != "" || reservedIpName != "" || okAuto) {
				return fmt.Errorf("[ERROR] Error creating instance template, network_interfaces error, reserved_ip(%s) is mutually exclusive with other primary_ip attributes", reservedIp)
			}
			if reservedIp != "" {
				nwInterface.PrimaryIP = &vpcv1.NetworkInterfaceIPPrototypeReservedIPIdentity{
					ID: &reservedIp,
				}
			} else {
				if PrimaryIpv4Address != "" || reservedIpAddress != "" || reservedIpName != "" || okAuto {
					primaryipobj := &vpcv1.NetworkInterfaceIPPrototypeReservedIPPrototypeNetworkInterfaceContext{}
					if PrimaryIpv4Address != "" {
						primaryipobj.Address = &PrimaryIpv4Address
					}
					if reservedIpAddress != "" {
						primaryipobj.Address = &reservedIpAddress
					}
					if reservedIpName != "" {
						primaryipobj.Name = &reservedIpName
					}
					if okAuto {
						primaryipobj.AutoDelete = &reservedIpAutoDelete
					}
					nwInterface.PrimaryIP = primaryipobj
				}
			}
			intfs = append(intfs, *nwInterface)
		}
		instanceproto.NetworkInterfaces = intfs
	}

	// Handle SSH Keys
	keySet := d.Get(isInstanceTemplateKeys).(*schema.Set)
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

	// Handle user data
	if userdata, ok := d.GetOk(isInstanceTemplateUserData); ok {
		userdatastr := userdata.(string)
		instanceproto.UserData = &userdatastr
	}

	// handle resource group
	if grp, ok := d.GetOk(isInstanceTemplateResourceGroup); ok {
		grpstr := grp.(string)
		instanceproto.ResourceGroup = &vpcv1.ResourceGroupIdentity{
			ID: &grpstr,
		}

	}

	options := &vpcv1.CreateInstanceTemplateOptions{
		InstanceTemplatePrototype: instanceproto,
	}

	instanceIntf, response, err := sess.CreateInstanceTemplate(options)
	if err != nil {
		return fmt.Errorf("[ERROR] Error creating InstanceTemplate: %s\n%s", err, response)
	}
	instance := instanceIntf.(*vpcv1.InstanceTemplate)
	d.SetId(*instance.ID)
	return nil
}

func instanceTemplateCreate(d *schema.ResourceData, meta interface{}, profile, name, vpcID, zone, image string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	instanceproto := &vpcv1.InstanceTemplatePrototype{
		Image: &vpcv1.ImageIdentity{
			ID: &image,
		},
		Zone: &vpcv1.ZoneIdentity{
			Name: &zone,
		},
		Profile: &vpcv1.InstanceProfileIdentity{
			Name: &profile,
		},
		VPC: &vpcv1.VPCIdentity{
			ID: &vpcID,
		},
	}
	if name != "" {
		instanceproto.Name = &name
	}
	// cluster changes
	if clusterNetworkAttachmentOk, ok := d.GetOk("cluster_network_attachments"); ok {
		clusterNetworkAttachmentList := clusterNetworkAttachmentOk.([]interface{})
		clusterNetworkAttachments := []vpcv1.InstanceClusterNetworkAttachmentPrototypeInstanceContext{}
		for _, clusterNetworkAttachmentsItem := range clusterNetworkAttachmentList {
			clusterNetworkAttachmentsItemModel, err := ResourceIBMIsInstanceTemplateMapToInstanceClusterNetworkAttachmentPrototypeInstanceContext(clusterNetworkAttachmentsItem.(map[string]interface{}))
			if err != nil {
				return err
			}
			clusterNetworkAttachments = append(clusterNetworkAttachments, *clusterNetworkAttachmentsItemModel)
		}
		instanceproto.ClusterNetworkAttachments = clusterNetworkAttachments
	}
	if _, ok := d.GetOk("confidential_compute_mode"); ok {
		instanceproto.ConfidentialComputeMode = core.StringPtr(d.Get("confidential_compute_mode").(string))
	}
	if _, ok := d.GetOkExists("enable_secure_boot"); ok {
		instanceproto.EnableSecureBoot = core.BoolPtr(d.Get("enable_secure_boot").(bool))
	}
	metadataServiceEnabled := d.Get(isInstanceTemplateMetadataServiceEnabled).(bool)
	if metadataServiceEnabled {
		instanceproto.MetadataService = &vpcv1.InstanceMetadataServicePrototype{
			Enabled: &metadataServiceEnabled,
		}
	}

	if metadataService := GetInstanceTemplateMetadataServiceOptions(d); metadataService != nil {
		instanceproto.MetadataService = metadataService
	}

	// vni

	if networkattachmentsintf, ok := d.GetOk("network_attachments"); ok {
		networkAttachments := []vpcv1.InstanceNetworkAttachmentPrototype{}
		for i, networkAttachmentsItem := range networkattachmentsintf.([]interface{}) {
			allowipspoofing := fmt.Sprintf("network_attachments.%d.virtual_network_interface.0.allow_ip_spoofing", i)
			autodelete := fmt.Sprintf("network_attachments.%d.virtual_network_interface.0.auto_delete", i)
			enablenat := fmt.Sprintf("network_attachments.%d.virtual_network_interface.0.enable_infrastructure_nat", i)
			networkAttachmentsItemModel, err := resourceIBMIsInstanceTemplateMapToInstanceNetworkAttachmentPrototype(allowipspoofing, autodelete, enablenat, d, networkAttachmentsItem.(map[string]interface{}))
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
		primaryNetworkAttachmentModel, err := resourceIBMIsInstanceTemplateMapToInstanceNetworkAttachmentPrototype(allowipspoofing, autodelete, enablenat, d, primnetworkattachmentintf.([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return err
		}
		instanceproto.PrimaryNetworkAttachment = primaryNetworkAttachmentModel
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
	if availablePolicyHostFailureIntf, ok := d.GetOk(isInstanceTemplateAvailablePolicyHostFailure); ok {
		availablePolicyHostFailure := availablePolicyHostFailureIntf.(string)
		instanceproto.AvailabilityPolicy = &vpcv1.InstanceAvailabilityPolicyPrototype{
			HostFailure: &availablePolicyHostFailure,
		}
	}
	if dHostIdInf, ok := d.GetOk(isPlacementTargetDedicatedHost); ok {
		dHostIdStr := dHostIdInf.(string)
		dHostPlaementTarget := &vpcv1.InstancePlacementTargetPrototypeDedicatedHostIdentity{
			ID: &dHostIdStr,
		}
		instanceproto.PlacementTarget = dHostPlaementTarget
	}

	if dHostGrpIdInf, ok := d.GetOk(isPlacementTargetDedicatedHostGroup); ok {
		dHostGrpIdStr := dHostGrpIdInf.(string)
		dHostGrpPlaementTarget := &vpcv1.InstancePlacementTargetPrototypeDedicatedHostGroupIdentity{
			ID: &dHostGrpIdStr,
		}
		instanceproto.PlacementTarget = dHostGrpPlaementTarget
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

	if placementGroupInf, ok := d.GetOk(isPlacementTargetPlacementGroup); ok {
		placementGrpStr := placementGroupInf.(string)
		placementGrp := &vpcv1.InstancePlacementTargetPrototypePlacementGroupIdentity{
			ID: &placementGrpStr,
		}
		instanceproto.PlacementTarget = placementGrp
	}

	if totalVolBandwidthIntf, ok := d.GetOk(isInstanceTotalVolumeBandwidth); ok {
		totalVolBandwidthStr := int64(totalVolBandwidthIntf.(int))
		instanceproto.TotalVolumeBandwidth = &totalVolBandwidthStr
	}

	// BOOT VOLUME ATTACHMENT for instance template
	if boot, ok := d.GetOk(isInstanceTemplateBootVolume); ok {
		bootvol := boot.([]interface{})[0].(map[string]interface{})
		var volTemplate = &vpcv1.VolumePrototypeInstanceByImageContext{}
		name, ok := bootvol[isInstanceTemplateBootName]
		namestr := name.(string)
		if ok && namestr != "" {
			volTemplate.Name = &namestr
		}

		var userTags *schema.Set
		if v, ok := bootvol[isInstanceTemplateBootVolumeTags]; ok {
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

		volcap := 100
		volcapint64 := int64(volcap)
		volprof := "general-purpose"
		volTemplate.Capacity = &volcapint64
		volTemplate.Profile = &vpcv1.VolumeProfileIdentity{
			Name: &volprof,
		}

		if encryption, ok := bootvol[isInstanceTemplateBootEncryption]; ok {
			bootEncryption := encryption.(string)
			if bootEncryption != "" {
				volTemplate.EncryptionKey = &vpcv1.EncryptionKeyIdentity{
					CRN: &bootEncryption,
				}
			}
		}

		var deleteVolumeOption bool
		if deleteVolume, ok := bootvol[isInstanceTemplateVolumeDeleteOnInstanceDelete]; ok {
			deleteVolumeOption = deleteVolume.(bool)
		}

		instanceproto.BootVolumeAttachment = &vpcv1.VolumeAttachmentPrototypeInstanceByImageContext{
			DeleteVolumeOnInstanceDelete: &deleteVolumeOption,
			Volume:                       volTemplate,
		}
	}

	// Handle volume attachments
	if volsintf, ok := d.GetOk(isInstanceTemplateVolumeAttachments); ok {
		vols := volsintf.([]interface{})
		var intfs []vpcv1.VolumeAttachmentPrototype
		for _, resource := range vols {
			vol := resource.(map[string]interface{})
			volInterface := &vpcv1.VolumeAttachmentPrototype{}
			deleteVolBool := vol[isInstanceTemplateVolumeDeleteOnInstanceDelete].(bool)
			volInterface.DeleteVolumeOnInstanceDelete = &deleteVolBool
			attachmentnamestr := vol[isInstanceTemplateVolAttachmentName].(string)
			volInterface.Name = &attachmentnamestr
			volIdStr := vol[isInstanceTemplateVolAttVol].(string)

			if volIdStr != "" {
				volInterface.Volume = &vpcv1.VolumeAttachmentPrototypeVolumeVolumeIdentity{
					ID: &volIdStr,
				}
			} else {
				newvolintf := vol[isInstanceTemplateVolAttVolPrototype].([]interface{})[0]
				newvol := newvolintf.(map[string]interface{})
				profileName := newvol[isInstanceTemplateVolAttVolProfile].(string)
				capacity := int64(newvol[isInstanceTemplateVolAttVolCapacity].(int))

				volPrototype := &vpcv1.VolumeAttachmentPrototypeVolumeVolumePrototypeInstanceContext{
					Profile: &vpcv1.VolumeProfileIdentity{
						Name: &profileName,
					},
					Capacity: &capacity,
				}
				iops := int64(newvol[isInstanceTemplateVolAttVolIops].(int))
				encryptionKey := newvol[isInstanceTemplateVolAttVolEncryptionKey].(string)

				if iops != 0 {
					volPrototype.Iops = &iops
				}

				if encryptionKey != "" {
					volPrototype.EncryptionKey = &vpcv1.EncryptionKeyIdentity{
						CRN: &encryptionKey,
					}
				}
				var userTags *schema.Set
				if v, ok := newvol[isInstanceTemplateVolAttTags]; ok {
					userTags = v.(*schema.Set)
					if userTags != nil && userTags.Len() != 0 {
						userTagsArray := make([]string, userTags.Len())
						for i, userTag := range userTags.List() {
							userTagStr := userTag.(string)
							userTagsArray[i] = userTagStr
						}
						volPrototype.UserTags = userTagsArray
					}
				}
				volInterface.Volume = volPrototype
			}

			intfs = append(intfs, *volInterface)
		}
		instanceproto.VolumeAttachments = intfs
	}

	// Handle primary network interface
	if primnicintf, ok := d.GetOk(isInstanceTemplatePrimaryNetworkInterface); ok {
		primnic := primnicintf.([]interface{})[0].(map[string]interface{})
		subnetintf, _ := primnic[isInstanceTemplateNicSubnet]
		subnetintfstr := subnetintf.(string)
		var primnicobj = &vpcv1.NetworkInterfacePrototype{}
		primnicobj.Subnet = &vpcv1.SubnetIdentity{
			ID: &subnetintfstr,
		}

		if name, ok := primnic[isInstanceTemplateNicName]; ok {
			namestr := name.(string)
			if namestr != "" {
				primnicobj.Name = &namestr
			}
		}
		allowIPSpoofing, ok := primnic[isInstanceTemplateNicAllowIPSpoofing]
		allowIPSpoofingbool := allowIPSpoofing.(bool)
		if ok {
			primnicobj.AllowIPSpoofing = &allowIPSpoofingbool
		}

		secgrpintf, ok := primnic[isInstanceTemplateNicSecurityGroups]
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
		// reserved ip changes
		var PrimaryIpv4Address, reservedIp, reservedIpAddress, reservedIpName string
		var reservedIpAutoDelete, okAuto bool
		if IPAddress, ok := primnic[isInstanceTemplateNicPrimaryIpv4Address]; ok {
			PrimaryIpv4Address = IPAddress.(string)
		}
		primaryIpOk, ok := primnic[isInstanceTemplateNicPrimaryIP]
		if ok && len(primaryIpOk.([]interface{})) > 0 {
			primip := primaryIpOk.([]interface{})[0].(map[string]interface{})

			reservedipok, _ := primip[isInstanceTemplateNicReservedIpId]
			reservedIp = reservedipok.(string)
			reservedipv4Ok, _ := primip[isInstanceTemplateNicReservedIpAddress]
			reservedIpAddress = reservedipv4Ok.(string)
			reservedipnameOk, _ := primip[isInstanceTemplateNicReservedIpName]
			reservedIpName = reservedipnameOk.(string)
			reservedipautodeleteok, okAuto := primip[isInstanceTemplateNicReservedIpAutoDelete]
			if okAuto {
				reservedIpAutoDelete = reservedipautodeleteok.(bool)
			}
		}
		if PrimaryIpv4Address != "" && reservedIpAddress != "" && PrimaryIpv4Address != reservedIpAddress {
			return fmt.Errorf("[ERROR] Error creating instance template, primary_network_interface error, use either primary_ipv4_address(%s) or primary_ip.0.address(%s)", PrimaryIpv4Address, reservedIpAddress)
		}
		if reservedIp != "" {
			primnicobj.PrimaryIP = &vpcv1.NetworkInterfaceIPPrototypeReservedIPIdentity{
				ID: &reservedIp,
			}
		} else {
			if PrimaryIpv4Address != "" || reservedIpAddress != "" || reservedIpName != "" || okAuto {
				primaryipobj := &vpcv1.NetworkInterfaceIPPrototypeReservedIPPrototypeNetworkInterfaceContext{}
				if PrimaryIpv4Address != "" {
					primaryipobj.Address = &PrimaryIpv4Address
				}
				if reservedIpAddress != "" {
					primaryipobj.Address = &reservedIpAddress
				}
				if reservedIpName != "" {
					primaryipobj.Name = &reservedIpName
				}
				if okAuto {
					primaryipobj.AutoDelete = &reservedIpAutoDelete
				}
				primnicobj.PrimaryIP = primaryipobj
			}
		}
		instanceproto.PrimaryNetworkInterface = primnicobj
	}

	// Handle  additional network interface
	if nicsintf, ok := d.GetOk(isInstanceTemplateNetworkInterfaces); ok {
		nics := nicsintf.([]interface{})
		var intfs []vpcv1.NetworkInterfacePrototype
		for _, resource := range nics {
			nic := resource.(map[string]interface{})
			nwInterface := &vpcv1.NetworkInterfacePrototype{}
			subnetintf, _ := nic[isInstanceTemplateNicSubnet]
			subnetintfstr := subnetintf.(string)
			nwInterface.Subnet = &vpcv1.SubnetIdentity{
				ID: &subnetintfstr,
			}

			name, ok := nic[isInstanceTemplateNicName]
			namestr := name.(string)
			if ok && namestr != "" {
				nwInterface.Name = &namestr
			}
			allowIPSpoofing, ok := nic[isInstanceTemplateNicAllowIPSpoofing]
			allowIPSpoofingbool := allowIPSpoofing.(bool)
			if ok {
				nwInterface.AllowIPSpoofing = &allowIPSpoofingbool
			}
			secgrpintf, ok := nic[isInstanceTemplateNicSecurityGroups]
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
			// reserved ip changes
			var PrimaryIpv4Address, reservedIp, reservedIpAddress, reservedIpName string
			var reservedIpAutoDelete, okAuto bool
			if IPAddress, ok := nic[isInstanceTemplateNicPrimaryIpv4Address]; ok {
				PrimaryIpv4Address = IPAddress.(string)
			}
			primaryIpOk, ok := nic[isInstanceTemplateNicPrimaryIP]
			if ok && len(primaryIpOk.([]interface{})) > 0 {
				primip := primaryIpOk.([]interface{})[0].(map[string]interface{})

				reservedipok, _ := primip[isInstanceTemplateNicReservedIpId]
				reservedIp = reservedipok.(string)
				reservedipv4Ok, _ := primip[isInstanceTemplateNicReservedIpAddress]
				reservedIpAddress = reservedipv4Ok.(string)
				reservedipnameOk, _ := primip[isInstanceTemplateNicReservedIpName]
				reservedIpName = reservedipnameOk.(string)
				// var reservedipautodeleteok interface{}

				if v, ok := primip[isInstanceTemplateNicReservedIpAutoDelete].(bool); ok && v {
					log.Printf("[INFO] isInstanceTemplateNicReservedIpAutoDelete is v is %t and okay is %t", v, ok)
					reservedIpAutoDelete = primip[isInstanceTemplateNicReservedIpAutoDelete].(bool)
					okAuto = true
				}
				// reservedipautodeleteok, okAuto = primip[isInstanceTemplateNicReservedIpAutoDelete]
				// if okAuto {
				// 	reservedIpAutoDelete = reservedipautodeleteok.(bool)
				// }
			}
			if PrimaryIpv4Address != "" && reservedIpAddress != "" && PrimaryIpv4Address != reservedIpAddress {
				return fmt.Errorf("[ERROR] Error creating instance template, network_interfaces error, use either primary_ipv4_address(%s) or primary_ip.0.address(%s)", PrimaryIpv4Address, reservedIpAddress)
			}
			if reservedIp != "" && (PrimaryIpv4Address != "" || reservedIpAddress != "" || reservedIpName != "" || okAuto) {
				return fmt.Errorf("[ERROR] Error creating instance template, network_interfaces error, reserved_ip(%s) is mutually exclusive with other primary_ip attributes", reservedIp)
			}
			if reservedIp != "" {
				nwInterface.PrimaryIP = &vpcv1.NetworkInterfaceIPPrototypeReservedIPIdentity{
					ID: &reservedIp,
				}
			} else {
				if PrimaryIpv4Address != "" || reservedIpAddress != "" || reservedIpName != "" || okAuto {
					primaryipobj := &vpcv1.NetworkInterfaceIPPrototypeReservedIPPrototypeNetworkInterfaceContext{}
					if PrimaryIpv4Address != "" {
						primaryipobj.Address = &PrimaryIpv4Address
					}
					if reservedIpAddress != "" {
						primaryipobj.Address = &reservedIpAddress
					}
					if reservedIpName != "" {
						primaryipobj.Name = &reservedIpName
					}
					if okAuto {
						primaryipobj.AutoDelete = &reservedIpAutoDelete
					}
					nwInterface.PrimaryIP = primaryipobj
				}
			}
			intfs = append(intfs, *nwInterface)
		}
		instanceproto.NetworkInterfaces = intfs
	}

	// Handle SSH Keys
	keySet := d.Get(isInstanceTemplateKeys).(*schema.Set)
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

	// Handle user data
	if userdata, ok := d.GetOk(isInstanceTemplateUserData); ok {
		userdatastr := userdata.(string)
		instanceproto.UserData = &userdatastr
	}

	// handle resource group
	if grp, ok := d.GetOk(isInstanceTemplateResourceGroup); ok {
		grpstr := grp.(string)
		instanceproto.ResourceGroup = &vpcv1.ResourceGroupIdentity{
			ID: &grpstr,
		}

	}

	options := &vpcv1.CreateInstanceTemplateOptions{
		InstanceTemplatePrototype: instanceproto,
	}

	instanceIntf, response, err := sess.CreateInstanceTemplate(options)
	if err != nil {
		return fmt.Errorf("[ERROR] Error creating InstanceTemplate: %s\n%s", err, response)
	}
	instance := instanceIntf.(*vpcv1.InstanceTemplate)
	d.SetId(*instance.ID)
	return nil
}

func instanceTemplateGet(d *schema.ResourceData, meta interface{}, ID string) error {
	instanceC, err := vpcClient(meta)
	if err != nil {
		return err
	}
	getinsOptions := &vpcv1.GetInstanceTemplateOptions{
		ID: &ID,
	}
	instanceIntf, response, err := instanceC.GetInstanceTemplate(getinsOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Getting Instance template: %s\n%s", err, response)
	}
	instance := instanceIntf.(*vpcv1.InstanceTemplate)
	d.Set(isInstanceTemplateName, *instance.Name)
	d.Set(isInstanceTemplateCRN, *instance.CRN)
	if instance.AvailabilityPolicy != nil && instance.AvailabilityPolicy.HostFailure != nil {
		d.Set(isInstanceTemplateAvailablePolicyHostFailure, instance.AvailabilityPolicy.HostFailure)
	}
	// cluster changes
	clusterNetworkAttachments := []map[string]interface{}{}
	for _, clusterNetworkAttachmentsItem := range instance.ClusterNetworkAttachments {
		clusterNetworkAttachmentsItemMap, err := ResourceIBMIsInstanceTemplateInstanceClusterNetworkAttachmentPrototypeInstanceContextToMap(&clusterNetworkAttachmentsItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_template", "read", "cluster_network_attachments-to-map")
		}
		clusterNetworkAttachments = append(clusterNetworkAttachments, clusterNetworkAttachmentsItemMap)
	}
	if err = d.Set("cluster_network_attachments", clusterNetworkAttachments); err != nil {
		err = fmt.Errorf("Error setting cluster_network_attachments: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_template", "read", "set-cluster_network_attachments")
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

	// vni if any
	if !core.IsNil(instance.NetworkAttachments) {
		networkAttachments := []map[string]interface{}{}
		for _, networkAttachmentsItem := range instance.NetworkAttachments {
			networkAttachmentsItemMap, err := resourceIBMIsInstanceTemplateNetworkAttachmentReferenceToMap(&networkAttachmentsItem, instanceC)
			if err != nil {
				return err
			}
			networkAttachments = append(networkAttachments, networkAttachmentsItemMap)
		}
		if err = d.Set("network_attachments", networkAttachments); err != nil {
			return fmt.Errorf("[ERROR] Error  setting network_attachments: %s", err)
		}
	}

	if !core.IsNil(instance.PrimaryNetworkAttachment) {
		primaryNetworkAttachmentMap, err := resourceIBMIsInstanceTemplateNetworkAttachmentReferenceToMap(instance.PrimaryNetworkAttachment, instanceC)
		if err != nil {
			return err
		}
		if err = d.Set("primary_network_attachment", []map[string]interface{}{primaryNetworkAttachmentMap}); err != nil {
			return fmt.Errorf("[ERROR] Error  setting primary_network_attachment: %s", err)
		}
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
		d.Set(isInstanceTemplateCatalogOffering, catOfferingList)

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
		d.Set(isReservationAffinity, reservationAffinity)
	}

	if instance.Profile != nil {
		instanceProfileIntf := instance.Profile
		identity := instanceProfileIntf.(*vpcv1.InstanceProfileIdentity)
		d.Set(isInstanceTemplateProfile, *identity.Name)
	}

	if instance.DefaultTrustedProfile != nil {
		if instance.DefaultTrustedProfile.AutoLink != nil {
			d.Set(isInstanceDefaultTrustedProfileAutoLink, instance.DefaultTrustedProfile.AutoLink)
		}
		if instance.DefaultTrustedProfile.Target != nil {
			switch reflect.TypeOf(instance.DefaultTrustedProfile.Target).String() {
			case "*vpcv1.TrustedProfileIdentityTrustedProfileByID":
				{
					target := instance.DefaultTrustedProfile.Target.(*vpcv1.TrustedProfileIdentityByID)
					d.Set(isInstanceDefaultTrustedProfileTarget, target.ID)
				}
			case "*vpcv1.TrustedProfileIdentityTrustedProfileByCRN":
				{
					target := instance.DefaultTrustedProfile.Target.(*vpcv1.TrustedProfileIdentityByCRN)
					d.Set(isInstanceDefaultTrustedProfileTarget, target.CRN)
				}
			}
		}
	}

	if instance.TotalVolumeBandwidth != nil {
		d.Set(isInstanceTotalVolumeBandwidth, int(*instance.TotalVolumeBandwidth))
	}
	if instance.MetadataService != nil {
		d.Set(isInstanceTemplateMetadataServiceEnabled, instance.MetadataService.Enabled)
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

	var placementTargetMap map[string]interface{}
	if instance.PlacementTarget != nil {
		placementTargetMap = resourceIbmIsInstanceTemplateInstancePlacementTargetPrototypeToMap(*instance.PlacementTarget.(*vpcv1.InstancePlacementTargetPrototype))
	}
	if err = d.Set(isInstanceTemplatePlacementTarget, []map[string]interface{}{placementTargetMap}); err != nil {
		return fmt.Errorf("[ERROR] Error setting placement_target: %s", err)
	}

	if instance.PrimaryNetworkInterface != nil {
		primaryNicList := make([]map[string]interface{}, 0)
		currentPrimNic := map[string]interface{}{}
		currentPrimNic[isInstanceTemplateNicName] = *instance.PrimaryNetworkInterface.Name
		if instance.PrimaryNetworkInterface.PrimaryIP != nil {
			pipIntf := instance.PrimaryNetworkInterface.PrimaryIP
			// reserved ip changes
			primaryIpList := make([]map[string]interface{}, 0)
			currentPrimIp := map[string]interface{}{}
			switch reflect.TypeOf(pipIntf).String() {
			case "*vpcv1.NetworkInterfaceIPPrototype":
				{
					pip := pipIntf.(*vpcv1.NetworkInterfaceIPPrototype)
					currentPrimNic[isInstanceTemplateNicPrimaryIpv4Address] = pip.Address
					currentPrimIp[isInstanceTemplateNicReservedIpId] = pip.ID
					currentPrimIp[isInstanceTemplateNicReservedIpAddress] = pip.Address
					currentPrimIp[isInstanceTemplateNicReservedIpAutoDelete] = pip.AutoDelete
					currentPrimIp[isInstanceTemplateNicReservedIpName] = pip.Name
				}
			case "*vpcv1.NetworkInterfaceIPPrototypeReservedIPPrototypeNetworkInterfaceContext":
				{
					pip := pipIntf.(*vpcv1.NetworkInterfaceIPPrototypeReservedIPPrototypeNetworkInterfaceContext)
					currentPrimNic[isInstanceTemplateNicPrimaryIpv4Address] = pip.Address
					currentPrimIp[isInstanceTemplateNicReservedIpAddress] = pip.Address
					currentPrimIp[isInstanceTemplateNicReservedIpAutoDelete] = pip.AutoDelete
					currentPrimIp[isInstanceTemplateNicReservedIpName] = pip.Name
				}
			case "*vpcv1.NetworkInterfaceIPPrototypeReservedIPIdentity":
				{
					pip := pipIntf.(*vpcv1.NetworkInterfaceIPPrototypeReservedIPIdentity)
					currentPrimIp[isInstanceTemplateNicReservedIpId] = pip.ID
				}
			}
			primaryIpList = append(primaryIpList, currentPrimIp)
			currentPrimNic[isInstanceTemplateNicPrimaryIP] = primaryIpList
		}
		subInf := instance.PrimaryNetworkInterface.Subnet
		subnetIdentity := subInf.(*vpcv1.SubnetIdentity)
		currentPrimNic[isInstanceTemplateNicSubnet] = *subnetIdentity.ID
		if instance.PrimaryNetworkInterface.AllowIPSpoofing != nil {
			currentPrimNic[isInstanceTemplateNicAllowIPSpoofing] = *instance.PrimaryNetworkInterface.AllowIPSpoofing
		}
		if len(instance.PrimaryNetworkInterface.SecurityGroups) != 0 {
			secgrpList := []string{}
			for i := 0; i < len(instance.PrimaryNetworkInterface.SecurityGroups); i++ {
				secGrpInf := instance.PrimaryNetworkInterface.SecurityGroups[i]
				subnetIdentity := secGrpInf.(*vpcv1.SecurityGroupIdentity)
				secgrpList = append(secgrpList, string(*subnetIdentity.ID))
			}
			currentPrimNic[isInstanceTemplateNicSecurityGroups] = flex.NewStringSet(schema.HashString, secgrpList)
		}
		primaryNicList = append(primaryNicList, currentPrimNic)
		d.Set(isInstanceTemplatePrimaryNetworkInterface, primaryNicList)
	}

	if instance.NetworkInterfaces != nil {
		interfacesList := make([]map[string]interface{}, 0)
		for _, intfc := range instance.NetworkInterfaces {
			currentNic := map[string]interface{}{}
			currentNic[isInstanceTemplateNicName] = *intfc.Name
			if intfc.PrimaryIP != nil {
				// reserved ip changes
				primaryIpList := make([]map[string]interface{}, 0)
				currentPrimIp := map[string]interface{}{}
				pipIntf := intfc.PrimaryIP
				switch reflect.TypeOf(pipIntf).String() {
				case "*vpcv1.NetworkInterfaceIPPrototype":
					{
						pip := pipIntf.(*vpcv1.NetworkInterfaceIPPrototype)
						currentNic[isInstanceTemplateNicPrimaryIpv4Address] = pip.Address
						currentPrimIp[isInstanceTemplateNicReservedIpAddress] = pip.Address
						currentPrimIp[isInstanceTemplateNicReservedIpAutoDelete] = pip.AutoDelete
						currentPrimIp[isInstanceTemplateNicReservedIpName] = pip.Name
						currentPrimIp[isInstanceTemplateNicReservedIpId] = pip.ID
					}
				case "*vpcv1.NetworkInterfaceIPPrototypeReservedIPPrototypeNetworkInterfaceContext":
					{
						pip := pipIntf.(*vpcv1.NetworkInterfaceIPPrototypeReservedIPPrototypeNetworkInterfaceContext)
						currentNic[isInstanceTemplateNicPrimaryIpv4Address] = pip.Address
						currentPrimIp[isInstanceTemplateNicReservedIpAddress] = pip.Address
						currentPrimIp[isInstanceTemplateNicReservedIpAutoDelete] = pip.AutoDelete
						currentPrimIp[isInstanceTemplateNicReservedIpName] = pip.Name
					}
				case "*vpcv1.NetworkInterfaceIPPrototypeReservedIPIdentity":
					{
						pip := pipIntf.(*vpcv1.NetworkInterfaceIPPrototypeReservedIPIdentity)
						currentPrimIp[isInstanceTemplateNicReservedIpId] = pip.ID
					}
				}
				primaryIpList = append(primaryIpList, currentPrimIp)
				currentNic[isInstanceTemplateNicPrimaryIP] = primaryIpList
			}
			if intfc.AllowIPSpoofing != nil {
				currentNic[isInstanceTemplateNicAllowIPSpoofing] = *intfc.AllowIPSpoofing
			}
			subInf := intfc.Subnet
			subnetIdentity := subInf.(*vpcv1.SubnetIdentity)
			currentNic[isInstanceTemplateNicSubnet] = *subnetIdentity.ID
			if len(intfc.SecurityGroups) != 0 {
				secgrpList := []string{}
				for i := 0; i < len(intfc.SecurityGroups); i++ {
					secGrpInf := intfc.SecurityGroups[i]
					subnetIdentity := secGrpInf.(*vpcv1.SecurityGroupIdentity)
					secgrpList = append(secgrpList, string(*subnetIdentity.ID))
				}
				currentNic[isInstanceTemplateNicSecurityGroups] = flex.NewStringSet(schema.HashString, secgrpList)
			}
			interfacesList = append(interfacesList, currentNic)
		}
		d.Set(isInstanceTemplateNetworkInterfaces, interfacesList)
	}

	if instance.Image != nil {
		imageInf := instance.Image
		imageIdentity := imageInf.(*vpcv1.ImageIdentity)
		d.Set(isInstanceTemplateImage, *imageIdentity.ID)
	}
	vpcInf := instance.VPC
	vpcRef := vpcInf.(*vpcv1.VPCIdentity)
	d.Set(isInstanceTemplateVPC, vpcRef.ID)
	zoneInf := instance.Zone
	zone := zoneInf.(*vpcv1.ZoneIdentity)
	d.Set(isInstanceTemplateZone, *zone.Name)

	interfacesList := make([]map[string]interface{}, 0)
	if instance.VolumeAttachments != nil {
		for _, volume := range instance.VolumeAttachments {
			volumeAttach := map[string]interface{}{}
			volumeAttach[isInstanceTemplateVolAttName] = *volume.Name
			volumeAttach[isInstanceTemplateDeleteVolume] = *volume.DeleteVolumeOnInstanceDelete
			newVolumeArr := []map[string]interface{}{}
			newVolume := map[string]interface{}{}
			volumeIntf := volume.Volume
			volumeInst := volumeIntf.(*vpcv1.VolumeAttachmentPrototypeVolume)
			if volumeInst.ID != nil {
				volumeAttach[isInstanceTemplateVolAttVol] = *volumeInst.ID
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
				newVolume[isInstanceTemplateVolAttTags] = volumeInst.UserTags
			}
			if len(newVolume) > 0 {
				newVolumeArr = append(newVolumeArr, newVolume)
			}
			volumeAttach[isInstanceTemplateVolAttVolPrototype] = newVolumeArr
			interfacesList = append(interfacesList, volumeAttach)
		}
		d.Set(isInstanceTemplateVolumeAttachments, interfacesList)
	}
	if instance.BootVolumeAttachment != nil {
		bootVolList := make([]map[string]interface{}, 0)
		bootVol := map[string]interface{}{}
		bootVol[isInstanceTemplateDeleteVolume] = *instance.BootVolumeAttachment.DeleteVolumeOnInstanceDelete
		if instance.BootVolumeAttachment.Volume != nil {
			volumeIntf := instance.BootVolumeAttachment.Volume
			bootVol[isInstanceTemplateBootName] = volumeIntf.Name
			bootVol[isInstanceTemplateBootSize] = volumeIntf.Capacity
			if volumeIntf.Profile != nil {
				volProfIntf := volumeIntf.Profile
				volProfInst := volProfIntf.(*vpcv1.VolumeProfileIdentity)
				bootVol[isInstanceTemplateBootProfile] = volProfInst.Name
			}
			if volumeIntf.EncryptionKey != nil {
				volEncryption := volumeIntf.EncryptionKey
				volEncryptionIntf := volEncryption.(*vpcv1.EncryptionKeyIdentity)
				bootVol[isInstanceTemplateBootEncryption] = volEncryptionIntf.CRN
			}
			if volumeIntf.UserTags != nil {
				bootVol[isVolumeTags] = volumeIntf.UserTags
			}
		}

		bootVolList = append(bootVolList, bootVol)
		d.Set(isInstanceTemplateBootVolume, bootVolList)
	}

	if instance.ResourceGroup != nil {
		d.Set(isInstanceTemplateResourceGroup, instance.ResourceGroup.ID)
	}
	return nil
}

func resourceIbmIsInstanceTemplateInstancePlacementTargetPrototypeToMap(instancePlacementTargetPrototype vpcv1.InstancePlacementTargetPrototype) map[string]interface{} {
	instancePlacementTargetPrototypeMap := map[string]interface{}{}

	instancePlacementTargetPrototypeMap["id"] = instancePlacementTargetPrototype.ID
	instancePlacementTargetPrototypeMap["crn"] = instancePlacementTargetPrototype.CRN
	instancePlacementTargetPrototypeMap["href"] = instancePlacementTargetPrototype.Href

	return instancePlacementTargetPrototypeMap
}

func instanceTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	instanceC, err := vpcClient(meta)
	if err != nil {
		return err
	}
	ID := d.Id()

	if d.HasChange(isInstanceName) {
		name := d.Get(isInstanceTemplateName).(string)
		updnetoptions := &vpcv1.UpdateInstanceTemplateOptions{
			ID: &ID,
		}

		instanceTemplatePatchModel := &vpcv1.InstanceTemplatePatch{
			Name: &name,
		}
		instanceTemplatePatch, err := instanceTemplatePatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("[ERROR] Error calling asPatch for InstanceTemplatePatch: %s", err)
		}
		updnetoptions.InstanceTemplatePatch = instanceTemplatePatch

		_, _, err = instanceC.UpdateInstanceTemplate(updnetoptions)
		if err != nil {
			return err
		}
	}
	return nil
}

func instanceTemplateDelete(d *schema.ResourceData, meta interface{}, ID string) error {
	instanceC, err := vpcClient(meta)
	if err != nil {
		return err
	}

	deleteinstanceTemplateOptions := &vpcv1.DeleteInstanceTemplateOptions{
		ID: &ID,
	}
	_, err = instanceC.DeleteInstanceTemplate(deleteinstanceTemplateOptions)
	if err != nil {
		return err
	}
	return nil
}

func instanceTemplateExists(d *schema.ResourceData, meta interface{}, ID string) (bool, error) {
	instanceC, err := vpcClient(meta)
	if err != nil {
		return false, err
	}
	getinsOptions := &vpcv1.GetInstanceTemplateOptions{
		ID: &ID,
	}
	_, response, err := instanceC.GetInstanceTemplate(getinsOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("[ERROR] Error Getting InstanceTemplate: %s\n%s", err, response)
	}
	return true, nil
}

func GetInstanceTemplateMetadataServiceOptions(d *schema.ResourceData) (metadataService *vpcv1.InstanceMetadataServicePrototype) {
	if metadataServiceIntf, ok := d.GetOk(isInstanceMetadataService); ok {
		metadataServiceMap := metadataServiceIntf.([]interface{})[0].(map[string]interface{})
		enabledIntf, ok := metadataServiceMap[isInstanceMetadataServiceEnabled1]
		metadataService = &vpcv1.InstanceMetadataServicePrototype{}
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

func resourceIBMIsInstanceTemplateNetworkAttachmentReferenceToMap(model *vpcv1.InstanceNetworkAttachmentPrototype, instanceC *vpcv1.VpcV1) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	vniMap := make(map[string]interface{})
	if model.VirtualNetworkInterface != nil {
		pna := model.VirtualNetworkInterface.(*vpcv1.InstanceNetworkAttachmentPrototypeVirtualNetworkInterface)
		if !core.IsNil(pna.Name) {
			vniMap["name"] = pna.Name
		}
		if !core.IsNil(pna.ID) {
			vniMap["id"] = pna.ID
		}
		if !core.IsNil(pna.AllowIPSpoofing) {
			vniMap["allow_ip_spoofing"] = pna.AllowIPSpoofing
		}
		if !core.IsNil(pna.AutoDelete) {
			vniMap["auto_delete"] = pna.AutoDelete
		}
		if !core.IsNil(pna.EnableInfrastructureNat) {
			vniMap["enable_infrastructure_nat"] = pna.EnableInfrastructureNat
		}
		// primaryipId := *vniDetails.PrimaryIP.ID
		if !core.IsNil(pna.Ips) {
			ips := []map[string]interface{}{}
			for _, ipsItem := range pna.Ips {
				// if *ipsItem.ID != primaryipId {
				ipsItemMap, err := resourceIBMIsInstanceTemplateVirtualNetworkInterfaceReservedIPReferenceToMap(ipsItem, false)
				if err != nil {
					return nil, err
				}
				ips = append(ips, ipsItemMap)
				// }
			}
			vniMap["ips"] = ips
		}

		if !core.IsNil(pna.SecurityGroups) {
			securityGroups := make([]string, 0)
			for _, securityGroup := range pna.SecurityGroups {
				securityGroupsItem := securityGroup.(*vpcv1.SecurityGroupIdentity)
				if securityGroupsItem.ID != nil {
					securityGroups = append(securityGroups, *securityGroupsItem.ID)
				}
			}
			vniMap["security_groups"] = securityGroups
		}
		if !core.IsNil(pna.PrimaryIP) {
			primaryIPMap, err := resourceIBMIsInstanceTemplateReservedIPReferenceToMap(pna.PrimaryIP, true)
			if err != nil {
				return modelMap, err
			}
			vniMap["primary_ip"] = []map[string]interface{}{primaryIPMap}
		}
		vniMap["protocol_state_filtering_mode"] = pna.ProtocolStateFilteringMode
		if pna.Subnet != nil {
			subnet := pna.Subnet.(*vpcv1.SubnetIdentity)
			vniMap["subnet"] = subnet.ID
		}

		modelMap["virtual_network_interface"] = []map[string]interface{}{vniMap}
	}

	return modelMap, nil
}

func resourceIBMIsInstanceTemplateVirtualNetworkInterfaceReservedIPReferenceToMap(modelIntf vpcv1.VirtualNetworkInterfaceIPPrototypeIntf, autodelete bool) (map[string]interface{}, error) {
	model := modelIntf.(*vpcv1.VirtualNetworkInterfaceIPPrototype)
	modelMap := make(map[string]interface{})
	modelMap["address"] = model.Address
	modelMap["auto_delete"] = autodelete
	modelMap["href"] = model.Href
	modelMap["reserved_ip"] = model.ID
	modelMap["name"] = model.Name

	return modelMap, nil
}

func resourceIBMIsInstanceTemplateReservedIPReferenceToMap(modelIntf vpcv1.VirtualNetworkInterfacePrimaryIPPrototypeIntf, autoDelete bool) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	model := modelIntf.(*vpcv1.VirtualNetworkInterfacePrimaryIPPrototype)
	modelMap["address"] = model.Address
	modelMap["href"] = model.Href
	modelMap["auto_delete"] = autoDelete
	modelMap["reserved_ip"] = model.ID
	modelMap["name"] = model.Name
	return modelMap, nil
}

func resourceIBMIsInstanceTemplateMapToInstanceNetworkAttachmentPrototype(allowipspoofing, autodelete, enablenat string, d *schema.ResourceData, modelMap map[string]interface{}) (*vpcv1.InstanceNetworkAttachmentPrototype, error) {
	model := &vpcv1.InstanceNetworkAttachmentPrototype{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}

	VirtualNetworkInterfaceModel, err := resourceIBMIsInstanceTemplateMapToVirtualNetworkInterfacePrototypeAttachmentContext(allowipspoofing, autodelete, enablenat, d, modelMap["virtual_network_interface"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.VirtualNetworkInterface = VirtualNetworkInterfaceModel
	return model, nil
}

func resourceIBMIsInstanceTemplateMapToVirtualNetworkInterfacePrototypeAttachmentContext(allowipspoofing, autodelete, enablenat string, d *schema.ResourceData, modelMap map[string]interface{}) (vpcv1.InstanceNetworkAttachmentPrototypeVirtualNetworkInterfaceIntf, error) {
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
		for _, ipsItem := range modelMap["ips"].([]interface{}) {
			ipsItemModel, err := resourceIBMIsInstanceTemplateMapToVirtualNetworkInterfaceIPsReservedIPPrototype(ipsItem.(map[string]interface{}))
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
		PrimaryIPModel, err := resourceIBMIsInstanceTemplateMapToVirtualNetworkInterfacePrimaryIPReservedIPPrototype(modelMap["primary_ip"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.PrimaryIP = PrimaryIPModel
	}
	if modelMap["protocol_state_filtering_mode"] != nil {
		if pStateFilteringInt, ok := modelMap["protocol_state_filtering_mode"]; ok && pStateFilteringInt.(string) != "" {
			model.ProtocolStateFilteringMode = core.StringPtr(pStateFilteringInt.(string))
		}
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

func resourceIBMIsInstanceTemplateMapToVirtualNetworkInterfaceIPsReservedIPPrototype(modelMap map[string]interface{}) (vpcv1.VirtualNetworkInterfaceIPPrototypeIntf, error) {
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

func resourceIBMIsInstanceTemplateMapToVirtualNetworkInterfacePrimaryIPReservedIPPrototype(modelMap map[string]interface{}) (vpcv1.VirtualNetworkInterfacePrimaryIPPrototypeIntf, error) {
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

func hashIpsListForInstanceTemplate(v interface{}) int {
	var buf bytes.Buffer
	a := v.(map[string]interface{})
	// buf.WriteString(fmt.Sprintf("%s-", a["address"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", a["reserved_ip"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", a["address"].(string)))
	return conns.String(buf.String())
}

func ResourceIBMIsInstanceTemplateInstanceClusterNetworkAttachmentPrototypeInstanceContextToMap(model *vpcv1.InstanceClusterNetworkAttachmentPrototypeInstanceContext) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	clusterNetworkInterfaceMap, err := ResourceIBMIsInstanceTemplateInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceToMap(model.ClusterNetworkInterface)
	if err != nil {
		return modelMap, err
	}
	modelMap["cluster_network_interface"] = []map[string]interface{}{clusterNetworkInterfaceMap}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	return modelMap, nil
}

func ResourceIBMIsInstanceTemplateInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceToMap(model vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceInstanceClusterNetworkInterfacePrototypeInstanceClusterNetworkAttachment); ok {
		return ResourceIBMIsInstanceTemplateInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceInstanceClusterNetworkInterfacePrototypeInstanceClusterNetworkAttachmentToMap(model.(*vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceInstanceClusterNetworkInterfacePrototypeInstanceClusterNetworkAttachment))
	} else if _, ok := model.(*vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentity); ok {
		return ResourceIBMIsInstanceTemplateInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityToMap(model.(*vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentity))
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
			primaryIPMap, err := ResourceIBMIsInstanceTemplateClusterNetworkInterfacePrimaryIPPrototypeToMap(model.PrimaryIP)
			if err != nil {
				return modelMap, err
			}
			modelMap["primary_ip"] = []map[string]interface{}{primaryIPMap}
		}
		if model.Subnet != nil {
			subnetMap, err := ResourceIBMIsInstanceTemplateClusterNetworkSubnetIdentityToMap(model.Subnet)
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

func ResourceIBMIsInstanceTemplateClusterNetworkInterfacePrimaryIPPrototypeToMap(model vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContext); ok {
		return ResourceIBMIsInstanceTemplateClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextToMap(model.(*vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContext))
	} else if _, ok := model.(*vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPPrototypeClusterNetworkInterfacePrimaryIPContext); ok {
		return ResourceIBMIsInstanceTemplateClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPPrototypeClusterNetworkInterfacePrimaryIPContextToMap(model.(*vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPPrototypeClusterNetworkInterfacePrimaryIPContext))
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

func ResourceIBMIsInstanceTemplateClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextToMap(model vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByID); ok {
		return ResourceIBMIsInstanceTemplateClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByIDToMap(model.(*vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByID))
	} else if _, ok := model.(*vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByHref); ok {
		return ResourceIBMIsInstanceTemplateClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByHrefToMap(model.(*vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByHref))
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

func ResourceIBMIsInstanceTemplateClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByIDToMap(model *vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByID) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	return modelMap, nil
}

func ResourceIBMIsInstanceTemplateClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByHrefToMap(model *vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByHref) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	return modelMap, nil
}

func ResourceIBMIsInstanceTemplateClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPPrototypeClusterNetworkInterfacePrimaryIPContextToMap(model *vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPPrototypeClusterNetworkInterfacePrimaryIPContext) (map[string]interface{}, error) {
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

func ResourceIBMIsInstanceTemplateClusterNetworkSubnetIdentityToMap(model vpcv1.ClusterNetworkSubnetIdentityIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.ClusterNetworkSubnetIdentityByID); ok {
		return ResourceIBMIsInstanceTemplateClusterNetworkSubnetIdentityByIDToMap(model.(*vpcv1.ClusterNetworkSubnetIdentityByID))
	} else if _, ok := model.(*vpcv1.ClusterNetworkSubnetIdentityByHref); ok {
		return ResourceIBMIsInstanceTemplateClusterNetworkSubnetIdentityByHrefToMap(model.(*vpcv1.ClusterNetworkSubnetIdentityByHref))
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

func ResourceIBMIsInstanceTemplateClusterNetworkSubnetIdentityByIDToMap(model *vpcv1.ClusterNetworkSubnetIdentityByID) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	return modelMap, nil
}

func ResourceIBMIsInstanceTemplateClusterNetworkSubnetIdentityByHrefToMap(model *vpcv1.ClusterNetworkSubnetIdentityByHref) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	return modelMap, nil
}

func ResourceIBMIsInstanceTemplateInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceInstanceClusterNetworkInterfacePrototypeInstanceClusterNetworkAttachmentToMap(model *vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceInstanceClusterNetworkInterfacePrototypeInstanceClusterNetworkAttachment) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AutoDelete != nil {
		modelMap["auto_delete"] = *model.AutoDelete
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.PrimaryIP != nil {
		primaryIPMap, err := ResourceIBMIsInstanceTemplateClusterNetworkInterfacePrimaryIPPrototypeToMap(model.PrimaryIP)
		if err != nil {
			return modelMap, err
		}
		modelMap["primary_ip"] = []map[string]interface{}{primaryIPMap}
	}
	if model.Subnet != nil {
		subnetMap, err := ResourceIBMIsInstanceTemplateClusterNetworkSubnetIdentityToMap(model.Subnet)
		if err != nil {
			return modelMap, err
		}
		modelMap["subnet"] = []map[string]interface{}{subnetMap}
	}
	return modelMap, nil
}

func ResourceIBMIsInstanceTemplateInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityToMap(model vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByID); ok {
		return ResourceIBMIsInstanceTemplateInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByIDToMap(model.(*vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByID))
	} else if _, ok := model.(*vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByHref); ok {
		return ResourceIBMIsInstanceTemplateInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByHrefToMap(model.(*vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByHref))
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

func ResourceIBMIsInstanceTemplateInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByIDToMap(model *vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByID) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	return modelMap, nil
}

func ResourceIBMIsInstanceTemplateInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByHrefToMap(model *vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByHref) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	return modelMap, nil
}

func ResourceIBMIsInstanceTemplateMapToInstanceClusterNetworkAttachmentPrototypeInstanceContext(modelMap map[string]interface{}) (*vpcv1.InstanceClusterNetworkAttachmentPrototypeInstanceContext, error) {
	model := &vpcv1.InstanceClusterNetworkAttachmentPrototypeInstanceContext{}
	ClusterNetworkInterfaceModel, err := ResourceIBMIsInstanceTemplateMapToInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterface(modelMap["cluster_network_interface"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.ClusterNetworkInterface = ClusterNetworkInterfaceModel
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	return model, nil
}

func ResourceIBMIsInstanceTemplateMapToInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterface(modelMap map[string]interface{}) (vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceIntf, error) {
	model := &vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterface{}
	if modelMap["auto_delete"] != nil {
		model.AutoDelete = core.BoolPtr(modelMap["auto_delete"].(bool))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["primary_ip"] != nil && len(modelMap["primary_ip"].([]interface{})) > 0 {
		PrimaryIPModel, err := ResourceIBMIsInstanceTemplateMapToClusterNetworkInterfacePrimaryIPPrototype(modelMap["primary_ip"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.PrimaryIP = PrimaryIPModel
	}
	if modelMap["subnet"] != nil && len(modelMap["subnet"].([]interface{})) > 0 {
		SubnetModel, err := ResourceIBMIsInstanceTemplateMapToClusterNetworkSubnetIdentity(modelMap["subnet"].([]interface{})[0].(map[string]interface{}))
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

func ResourceIBMIsInstanceTemplateMapToClusterNetworkInterfacePrimaryIPPrototype(modelMap map[string]interface{}) (vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeIntf, error) {
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

func ResourceIBMIsInstanceTemplateMapToClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContext(modelMap map[string]interface{}) (vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextIntf, error) {
	model := &vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContext{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	return model, nil
}

func ResourceIBMIsInstanceTemplateMapToClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByID(modelMap map[string]interface{}) (*vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByID, error) {
	model := &vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func ResourceIBMIsInstanceTemplateMapToClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByHref(modelMap map[string]interface{}) (*vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByHref, error) {
	model := &vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func ResourceIBMIsInstanceTemplateMapToClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPPrototypeClusterNetworkInterfacePrimaryIPContext(modelMap map[string]interface{}) (*vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPPrototypeClusterNetworkInterfacePrimaryIPContext, error) {
	model := &vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPPrototypeClusterNetworkInterfacePrimaryIPContext{}
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

func ResourceIBMIsInstanceTemplateMapToClusterNetworkSubnetIdentity(modelMap map[string]interface{}) (vpcv1.ClusterNetworkSubnetIdentityIntf, error) {
	model := &vpcv1.ClusterNetworkSubnetIdentity{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	return model, nil
}

func ResourceIBMIsInstanceTemplateMapToClusterNetworkSubnetIdentityByID(modelMap map[string]interface{}) (*vpcv1.ClusterNetworkSubnetIdentityByID, error) {
	model := &vpcv1.ClusterNetworkSubnetIdentityByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func ResourceIBMIsInstanceTemplateMapToClusterNetworkSubnetIdentityByHref(modelMap map[string]interface{}) (*vpcv1.ClusterNetworkSubnetIdentityByHref, error) {
	model := &vpcv1.ClusterNetworkSubnetIdentityByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func ResourceIBMIsInstanceTemplateMapToInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceInstanceClusterNetworkInterfacePrototypeInstanceClusterNetworkAttachment(modelMap map[string]interface{}) (*vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceInstanceClusterNetworkInterfacePrototypeInstanceClusterNetworkAttachment, error) {
	model := &vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceInstanceClusterNetworkInterfacePrototypeInstanceClusterNetworkAttachment{}
	if modelMap["auto_delete"] != nil {
		model.AutoDelete = core.BoolPtr(modelMap["auto_delete"].(bool))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["primary_ip"] != nil && len(modelMap["primary_ip"].([]interface{})) > 0 {
		PrimaryIPModel, err := ResourceIBMIsInstanceTemplateMapToClusterNetworkInterfacePrimaryIPPrototype(modelMap["primary_ip"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.PrimaryIP = PrimaryIPModel
	}
	if modelMap["subnet"] != nil && len(modelMap["subnet"].([]interface{})) > 0 {
		SubnetModel, err := ResourceIBMIsInstanceTemplateMapToClusterNetworkSubnetIdentity(modelMap["subnet"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Subnet = SubnetModel
	}
	return model, nil
}

func ResourceIBMIsInstanceTemplateMapToInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentity(modelMap map[string]interface{}) (vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityIntf, error) {
	model := &vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentity{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	return model, nil
}

func ResourceIBMIsInstanceTemplateMapToInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByID(modelMap map[string]interface{}) (*vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByID, error) {
	model := &vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func ResourceIBMIsInstanceTemplateMapToInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByHref(modelMap map[string]interface{}) (*vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByHref, error) {
	model := &vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}
