// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isInstanceTemplateBootVolume                   = "boot_volume"
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

			isInstanceTemplateMetadataServiceEnabled: {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Indicates whether the metadata service endpoint is available to the virtual server instance",
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
				Computed: true,
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
								},
							},
						},
					},
				},
			},

			isInstanceTemplatePrimaryNetworkInterface: {
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Description: "Primary Network interface info",
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
				Type:     schema.TypeList,
				Optional: true,
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
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "image name",
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
		},
	}
}

func ResourceIBMISInstanceTemplateValidator() *validate.ResourceValidator {
	host_failure := "restart, stop"
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
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

	ibmISInstanceTemplateValidator := validate.ResourceValidator{ResourceName: "ibm_is_instance_template", Schema: validateSchema}
	return &ibmISInstanceTemplateValidator
}

func resourceIBMisInstanceTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	profile := d.Get(isInstanceTemplateProfile).(string)
	name := d.Get(isInstanceTemplateName).(string)
	vpcID := d.Get(isInstanceTemplateVPC).(string)
	zone := d.Get(isInstanceTemplateZone).(string)
	image := d.Get(isInstanceTemplateImage).(string)

	err := instanceTemplateCreate(d, meta, profile, name, vpcID, zone, image)
	if err != nil {
		return err
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

	metadataServiceEnabled := d.Get(isInstanceTemplateMetadataServiceEnabled).(bool)
	if metadataServiceEnabled {
		instanceproto.MetadataService = &vpcv1.InstanceMetadataServicePrototype{
			Enabled: &metadataServiceEnabled,
		}
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
		instanceproto.AvailabilityPolicy = &vpcv1.InstanceAvailabilityPrototype{
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
		var intfs []vpcv1.VolumeAttachmentPrototypeInstanceContext
		for _, resource := range vols {
			vol := resource.(map[string]interface{})
			volInterface := &vpcv1.VolumeAttachmentPrototypeInstanceContext{}
			deleteVolBool := vol[isInstanceTemplateVolumeDeleteOnInstanceDelete].(bool)
			volInterface.DeleteVolumeOnInstanceDelete = &deleteVolBool
			attachmentnamestr := vol[isInstanceTemplateVolAttachmentName].(string)
			volInterface.Name = &attachmentnamestr
			volIdStr := vol[isInstanceTemplateVolAttVol].(string)

			if volIdStr != "" {
				volInterface.Volume = &vpcv1.VolumeAttachmentVolumePrototypeInstanceContextVolumeIdentity{
					ID: &volIdStr,
				}
			} else {
				newvolintf := vol[isInstanceTemplateVolAttVolPrototype].([]interface{})[0]
				newvol := newvolintf.(map[string]interface{})
				profileName := newvol[isInstanceTemplateVolAttVolProfile].(string)
				capacity := int64(newvol[isInstanceTemplateVolAttVolCapacity].(int))

				volPrototype := &vpcv1.VolumeAttachmentVolumePrototypeInstanceContextVolumePrototypeInstanceContext{
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
					log.Printf("[INFO] UJJK isInstanceTemplateNicReservedIpAutoDelete is v is %t and okay is %t", v, ok)
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
					target := instance.DefaultTrustedProfile.Target.(*vpcv1.TrustedProfileIdentityTrustedProfileByID)
					d.Set(isInstanceDefaultTrustedProfileTarget, target.ID)
				}
			case "*vpcv1.TrustedProfileIdentityTrustedProfileByCRN":
				{
					target := instance.DefaultTrustedProfile.Target.(*vpcv1.TrustedProfileIdentityTrustedProfileByCRN)
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
			volumeInst := volumeIntf.(*vpcv1.VolumeAttachmentVolumePrototypeInstanceContext)
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
