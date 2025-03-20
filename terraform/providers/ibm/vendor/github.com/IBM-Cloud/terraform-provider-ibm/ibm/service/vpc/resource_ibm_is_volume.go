// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
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
	isVolumeName                  = "name"
	isVolumeProfileName           = "profile"
	isVolumeZone                  = "zone"
	isVolumeEncryptionKey         = "encryption_key"
	isVolumeEncryptionType        = "encryption_type"
	isVolumeCapacity              = "capacity"
	isVolumeIops                  = "iops"
	isVolumeCrn                   = "crn"
	isVolumeTags                  = "tags"
	isVolumeStatus                = "status"
	isVolumeStatusReasons         = "status_reasons"
	isVolumeStatusReasonsCode     = "code"
	isVolumeStatusReasonsMessage  = "message"
	isVolumeStatusReasonsMoreInfo = "more_info"
	isVolumeDeleting              = "deleting"
	isVolumeDeleted               = "done"
	isVolumeProvisioning          = "provisioning"
	isVolumeProvisioningDone      = "done"
	isVolumeResourceGroup         = "resource_group"
	isVolumeSourceSnapshot        = "source_snapshot"
	isVolumeSourceSnapshotCrn     = "source_snapshot_crn"
	isVolumeDeleteAllSnapshots    = "delete_all_snapshots"
	isVolumeBandwidth             = "bandwidth"
	isVolumeAccessTags            = "access_tags"
	isVolumeUserTagType           = "user"
	isVolumeAccessTagType         = "access"
	isVolumeHealthReasons         = "health_reasons"
	isVolumeHealthReasonsCode     = "code"
	isVolumeHealthReasonsMessage  = "message"
	isVolumeHealthReasonsMoreInfo = "more_info"
	isVolumeHealthState           = "health_state"

	isVolumeCatalogOffering           = "catalog_offering"
	isVolumeCatalogOfferingPlanCrn    = "plan_crn"
	isVolumeCatalogOfferingVersionCrn = "version_crn"
)

func ResourceIBMISVolume() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMISVolumeCreate,
		Read:     resourceIBMISVolumeRead,
		Update:   resourceIBMISVolumeUpdate,
		Delete:   resourceIBMISVolumeDelete,
		Exists:   resourceIBMISVolumeExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			customdiff.Sequence(
				func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
					return flex.ResourceTagsCustomizeDiff(diff)
				},
			),
			customdiff.Sequence(
				func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
					return flex.ResourceVolumeValidate(diff)
				}),
			customdiff.Sequence(
				func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
					return flex.ResourceValidateAccessTags(diff, v)
				}),
		),

		Schema: map[string]*schema.Schema{

			isVolumeName: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_volume", isVolumeName),
				Description:  "Volume name",
			},

			isVolumeProfileName: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_volume", isVolumeProfileName),
				Description:  "Volume profile name",
			},

			isVolumeZone: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone name",
			},

			isVolumeEncryptionKey: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Volume encryption key info",
			},

			isVolumeEncryptionType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Volume encryption type info",
			},

			isVolumeCapacity: {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Computed: true,
				// ValidateFunc: validate.InvokeValidator("ibm_is_volume", isVolumeCapacity),
				Description: "Volume capacity value",
			},
			isVolumeSourceSnapshot: {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{isVolumeSourceSnapshotCrn},
				ValidateFunc:  validate.InvokeValidator("ibm_is_volume", isVolumeSourceSnapshot),
				Description:   "The unique identifier for this snapshot",
			},
			isVolumeSourceSnapshotCrn: {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{isVolumeSourceSnapshot},
				Description:   "The crn for this snapshot",
			},
			isVolumeResourceGroup: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Resource group name",
			},
			isVolumeIops: {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				// ValidateFunc: validate.InvokeValidator("ibm_is_volume", isVolumeIops),
				Description: "IOPS value for the Volume",
			},
			isVolumeCrn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CRN value for the volume instance",
			},
			isVolumeStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Volume status",
			},

			isVolumeStatusReasons: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isVolumeStatusReasonsCode: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A snake case string succinctly identifying the status reason",
						},

						isVolumeStatusReasonsMessage: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An explanation of the status reason",
						},

						isVolumeStatusReasonsMoreInfo: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Link to documentation about this status reason",
						},
					},
				},
			},
			isVolumeCatalogOffering: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The catalog offering this volume was created from. If a virtual server instance is provisioned with a boot_volume_attachment specifying this volume, the virtual server instance will use this volume's catalog offering, including its pricing plan.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isVolumeCatalogOfferingPlanCrn: {
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
						isVolumeCatalogOfferingVersionCrn: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this version of a catalog offering",
						},
					},
				},
			},
			// defined_performance changes
			"adjustable_capacity_states": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The attachment states that support adjustable capacity for this volume.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"adjustable_iops_states": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The attachment states that support adjustable IOPS for this volume.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			isVolumeHealthReasons: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isVolumeHealthReasonsCode: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A snake case string succinctly identifying the reason for this health state.",
						},

						isVolumeHealthReasonsMessage: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An explanation of the reason for this health state.",
						},

						isVolumeHealthReasonsMoreInfo: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Link to documentation about the reason for this health state.",
						},
					},
				},
			},

			isVolumeHealthState: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The health of this resource.",
			},
			isVolumeDeleteAllSnapshots: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Deletes all snapshots created from this volume",
			},
			isVolumeTags: {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_volume", "tags")},
				Set:         flex.ResourceIBMVPCHash,
				Description: "UserTags for the volume instance",
			},
			isVolumeAccessTags: {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_volume", "accesstag")},
				Set:         flex.ResourceIBMVPCHash,
				Description: "Access management tags for the volume instance",
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

			isVolumeBandwidth: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum bandwidth (in megabits per second) for the volume",
			},

			isVolumesOperatingSystem: &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The operating system associated with this volume. If absent, this volume was notcreated from an image, or the image did not include an operating system.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isVolumeArchitecture: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The operating system architecture.",
						},
						isVolumeDHOnly: &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Images with this operating system can only be used on dedicated hosts or dedicated host groups.",
						},
						isVolumeDisplayName: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A unique, display-friendly name for the operating system.",
						},
						isVolumeOSFamily: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The software family for this operating system.",
						},

						isVolumesOperatingSystemHref: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this operating system.",
						},
						isVolumesOperatingSystemName: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The globally unique name for this operating system.",
						},
						isVolumeOSVendor: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The vendor of the operating system.",
						},
						isVolumeOSVersion: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The major release version of this operating system.",
						},
					},
				},
			},
		},
	}
}

func ResourceIBMISVolumeValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isVolumeName,
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isVolumeSourceSnapshot,
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[-0-9a-z_]+$`,
			MinValueLength:             1,
			MaxValueLength:             64})
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
			Identifier:                 isVolumeProfileName,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "general-purpose, 5iops-tier, 10iops-tier, custom, sdp",
		})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isVolumeCapacity,
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			MinValue:                   "10"})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isVolumeIops,
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			MinValue:                   "100"})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "accesstag",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([A-Za-z0-9_.-]|[A-Za-z0-9_.-][A-Za-z0-9_ .-]*[A-Za-z0-9_.-]):([A-Za-z0-9_.-]|[A-Za-z0-9_.-][A-Za-z0-9_ .-]*[A-Za-z0-9_.-])$`,
			MinValueLength:             1,
			MaxValueLength:             128})

	ibmISVolumeResourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_volume", Schema: validateSchema}
	return &ibmISVolumeResourceValidator
}

func resourceIBMISVolumeCreate(d *schema.ResourceData, meta interface{}) error {

	volName := d.Get(isVolumeName).(string)
	profile := d.Get(isVolumeProfileName).(string)
	zone := d.Get(isVolumeZone).(string)

	err := volCreate(d, meta, volName, profile, zone)
	if err != nil {
		return err
	}

	return resourceIBMISVolumeRead(d, meta)
}

func volCreate(d *schema.ResourceData, meta interface{}, volName, profile, zone string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	options := &vpcv1.CreateVolumeOptions{}
	volTemplate := &vpcv1.VolumePrototype{
		Name: &volName,
		Zone: &vpcv1.ZoneIdentity{
			Name: &zone,
		},
		Profile: &vpcv1.VolumeProfileIdentity{
			Name: &profile,
		},
	}

	if sourceSnapsht, ok := d.GetOk(isVolumeSourceSnapshot); ok {
		sourceSnapshot := sourceSnapsht.(string)
		snapshotIdentity := &vpcv1.SnapshotIdentity{
			ID: &sourceSnapshot,
		}
		volTemplate.SourceSnapshot = snapshotIdentity
		if capacity, ok := d.GetOk(isVolumeCapacity); ok {
			if int64(capacity.(int)) > 0 {
				volCapacity := int64(capacity.(int))
				volTemplate.Capacity = &volCapacity
			}
		}
	} else if sourceSnapshtCrn, ok := d.GetOk(isVolumeSourceSnapshotCrn); ok {
		sourceSnapshot := sourceSnapshtCrn.(string)

		snapshotIdentity := &vpcv1.SnapshotIdentity{
			CRN: &sourceSnapshot,
		}
		volTemplate.SourceSnapshot = snapshotIdentity
		if capacity, ok := d.GetOk(isVolumeCapacity); ok {
			if int64(capacity.(int)) > 0 {
				volCapacity := int64(capacity.(int))
				volTemplate.Capacity = &volCapacity
			}
		}
	} else if capacity, ok := d.GetOk(isVolumeCapacity); ok {
		if int64(capacity.(int)) > 0 {
			volCapacity := int64(capacity.(int))
			volTemplate.Capacity = &volCapacity
		}
	} else {
		volCapacity := int64(100)
		volTemplate.Capacity = &volCapacity
	}

	if key, ok := d.GetOk(isVolumeEncryptionKey); ok {
		encryptionKey := key.(string)
		volTemplate.EncryptionKey = &vpcv1.EncryptionKeyIdentity{
			CRN: &encryptionKey,
		}
	}

	if rgrp, ok := d.GetOk(isVolumeResourceGroup); ok {
		rg := rgrp.(string)
		volTemplate.ResourceGroup = &vpcv1.ResourceGroupIdentity{
			ID: &rg,
		}
	}

	if i, ok := d.GetOk(isVolumeIops); ok {
		iops := int64(i.(int))
		volTemplate.Iops = &iops
	}

	var userTags *schema.Set
	if v, ok := d.GetOk(isVolumeTags); ok {
		userTags = v.(*schema.Set)
		if userTags != nil && userTags.Len() != 0 {
			userTagsArray := make([]string, userTags.Len())
			for i, userTag := range userTags.List() {
				userTagStr := userTag.(string)
				userTagsArray[i] = userTagStr
			}
			schematicTags := os.Getenv("IC_ENV_TAGS")
			var envTags []string
			if schematicTags != "" {
				envTags = strings.Split(schematicTags, ",")
				userTagsArray = append(userTagsArray, envTags...)
			}
			volTemplate.UserTags = userTagsArray
		}
	}
	options.VolumePrototype = volTemplate
	vol, response, err := sess.CreateVolume(options)
	if err != nil {
		return fmt.Errorf("[DEBUG] Create volume err %s\n%s", err, response)
	}
	d.SetId(*vol.ID)
	log.Printf("[INFO] Volume : %s", *vol.ID)
	_, err = isWaitForVolumeAvailable(sess, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	if _, ok := d.GetOk(isVolumeAccessTags); ok {
		oldList, newList := d.GetChange(isVolumeAccessTags)
		err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *vol.CRN, "", isVolumeAccessTagType)
		if err != nil {
			log.Printf(
				"Error on create of resource vpc volume (%s) access tags: %s", d.Id(), err)
		}
	}
	return nil
}

func resourceIBMISVolumeRead(d *schema.ResourceData, meta interface{}) error {

	id := d.Id()
	err := volGet(d, meta, id)
	if err != nil {
		return err
	}
	return nil
}

func volGet(d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	options := &vpcv1.GetVolumeOptions{
		ID: &id,
	}
	vol, response, err := sess.GetVolume(options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error getting Volume (%s): %s\n%s", id, err, response)
	}
	d.SetId(*vol.ID)
	d.Set(isVolumeName, *vol.Name)
	d.Set(isVolumeProfileName, *vol.Profile.Name)
	d.Set(isVolumeZone, *vol.Zone.Name)
	if vol.EncryptionKey != nil {
		d.Set(isVolumeEncryptionKey, vol.EncryptionKey.CRN)
	}
	if vol.Encryption != nil {
		d.Set(isVolumeEncryptionType, vol.Encryption)
	}
	d.Set(isVolumeIops, *vol.Iops)
	d.Set(isVolumeCapacity, *vol.Capacity)
	d.Set(isVolumeCrn, *vol.CRN)
	if vol.SourceSnapshot != nil {
		d.Set(isVolumeSourceSnapshot, *vol.SourceSnapshot.ID)
		d.Set(isVolumeSourceSnapshotCrn, *vol.SourceSnapshot.CRN)
	}
	d.Set(isVolumeStatus, *vol.Status)
	if vol.HealthState != nil {
		d.Set(isVolumeHealthState, *vol.HealthState)
	}
	d.Set(isVolumeBandwidth, int(*vol.Bandwidth))
	//set the status reasons
	if vol.StatusReasons != nil {
		statusReasonsList := make([]map[string]interface{}, 0)
		for _, sr := range vol.StatusReasons {
			currentSR := map[string]interface{}{}
			if sr.Code != nil && sr.Message != nil {
				currentSR[isVolumeStatusReasonsCode] = *sr.Code
				currentSR[isVolumeStatusReasonsMessage] = *sr.Message
				if sr.MoreInfo != nil {
					currentSR[isVolumeStatusReasonsMoreInfo] = *sr.Message
				}
				statusReasonsList = append(statusReasonsList, currentSR)
			}
		}
		d.Set(isVolumeStatusReasons, statusReasonsList)
	}
	if vol.UserTags != nil {
		if err = d.Set(isVolumeTags, vol.UserTags); err != nil {
			return fmt.Errorf("Error setting user tags: %s", err)
		}
	}
	accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *vol.CRN, "", isVolumeAccessTagType)
	if err != nil {
		log.Printf(
			"Error on get of resource volume (%s) access tags: %s", d.Id(), err)
	}
	d.Set(isVolumeAccessTags, accesstags)
	if vol.HealthReasons != nil {
		healthReasonsList := make([]map[string]interface{}, 0)
		for _, sr := range vol.HealthReasons {
			currentSR := map[string]interface{}{}
			if sr.Code != nil && sr.Message != nil {
				currentSR[isVolumeHealthReasonsCode] = *sr.Code
				currentSR[isVolumeHealthReasonsMessage] = *sr.Message
				if sr.MoreInfo != nil {
					currentSR[isVolumeHealthReasonsMoreInfo] = *sr.Message
				}
				healthReasonsList = append(healthReasonsList, currentSR)
			}
		}
		d.Set(isVolumeHealthReasons, healthReasonsList)
	}
	// catalog
	catalogList := make([]map[string]interface{}, 0)
	if vol.CatalogOffering != nil {
		versionCrn := ""
		if vol.CatalogOffering.Version != nil && vol.CatalogOffering.Version.CRN != nil {
			versionCrn = *vol.CatalogOffering.Version.CRN
		}
		catalogMap := map[string]interface{}{}
		if versionCrn != "" {
			catalogMap[isVolumeCatalogOfferingVersionCrn] = versionCrn
		}
		if vol.CatalogOffering.Plan != nil {
			planCrn := ""
			if vol.CatalogOffering.Plan.CRN != nil {
				planCrn = *vol.CatalogOffering.Plan.CRN
			}
			if planCrn != "" {
				catalogMap[isVolumeCatalogOfferingPlanCrn] = *vol.CatalogOffering.Plan.CRN
			}
			if vol.CatalogOffering.Plan.Deleted != nil {
				deletedMap := resourceIbmIsVolumeCatalogOfferingVersionPlanReferenceDeletedToMap(*vol.CatalogOffering.Plan.Deleted)
				catalogMap["deleted"] = []map[string]interface{}{deletedMap}
			}
		}
		catalogList = append(catalogList, catalogMap)
	}
	d.Set(isVolumeCatalogOffering, catalogList)
	controller, err := flex.GetBaseController(meta)
	if err != nil {
		return err
	}
	// defined_performance changes

	if err = d.Set("adjustable_capacity_states", vol.AdjustableCapacityStates); err != nil {
		err = fmt.Errorf("Error setting adjustable_capacity_states: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_volume", "read", "set-adjustable_capacity_states")
	}
	if err = d.Set("adjustable_iops_states", vol.AdjustableIopsStates); err != nil {
		err = fmt.Errorf("Error setting adjustable_iops_states: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_volume", "read", "set-adjustable_iops_states")
	}
	d.Set(flex.ResourceControllerURL, controller+"/vpc-ext/storage/storageVolumes")
	d.Set(flex.ResourceName, *vol.Name)
	d.Set(flex.ResourceCRN, *vol.CRN)
	d.Set(flex.ResourceStatus, *vol.Status)
	if vol.ResourceGroup != nil {
		d.Set(flex.ResourceGroupName, vol.ResourceGroup.Name)
		d.Set(isVolumeResourceGroup, *vol.ResourceGroup.ID)
	}
	operatingSystemList := []map[string]interface{}{}
	if vol.OperatingSystem != nil {
		operatingSystemMap := dataSourceVolumeCollectionVolumesOperatingSystemToMap(*vol.OperatingSystem)
		operatingSystemList = append(operatingSystemList, operatingSystemMap)
	}
	d.Set(isVolumesOperatingSystem, operatingSystemList)
	return nil
}

func resourceIBMISVolumeUpdate(d *schema.ResourceData, meta interface{}) error {

	id := d.Id()
	name := ""
	hasNameChanged := false
	delete := false

	if delete_all_snapshots, ok := d.GetOk(isVolumeDeleteAllSnapshots); ok && delete_all_snapshots.(bool) {
		delete = true
	}

	if d.HasChange(isVolumeName) {
		name = d.Get(isVolumeName).(string)
		hasNameChanged = true
	}

	err := volUpdate(d, meta, id, name, hasNameChanged, delete)
	if err != nil {
		return err
	}
	return resourceIBMISVolumeRead(d, meta)
}

func volUpdate(d *schema.ResourceData, meta interface{}, id, name string, hasNameChanged, delete bool) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	var capacity int64
	if delete {
		deleteAllSnapshots(sess, id)
	}

	if d.HasChange(isVolumeAccessTags) {
		options := &vpcv1.GetVolumeOptions{
			ID: &id,
		}
		vol, response, err := sess.GetVolume(options)
		if err != nil {
			return fmt.Errorf("Error getting Volume : %s\n%s", err, response)
		}
		oldList, newList := d.GetChange(isVolumeAccessTags)

		err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *vol.CRN, "", isVolumeAccessTagType)
		if err != nil {
			log.Printf(
				"Error on update of resource vpc volume (%s) access tags: %s", id, err)
		}
	}

	optionsget := &vpcv1.GetVolumeOptions{
		ID: &id,
	}
	oldVol, response, err := sess.GetVolume(optionsget)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting Volume (%s): %s\n%s", id, err, response)
	}
	eTag := response.Headers.Get("ETag")
	options := &vpcv1.UpdateVolumeOptions{
		ID: &id,
	}
	options.IfMatch = &eTag

	//name update
	volumeNamePatchModel := &vpcv1.VolumePatch{}
	if hasNameChanged {
		volumeNamePatchModel.Name = &name
		volumeNamePatch, err := volumeNamePatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("[ERROR] Error calling asPatch for volumeNamePatch: %s", err)
		}
		options.VolumePatch = volumeNamePatch
		_, response, err = sess.UpdateVolume(options)
		if err != nil {
			return err
		}
		_, err = isWaitForVolumeAvailable(sess, d.Id(), d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return err
		}
		eTag = response.Headers.Get("ETag")
		options.IfMatch = &eTag
	}

	// profile/ iops update
	if !d.HasChange(isVolumeProfileName) && *oldVol.Profile.Name == "sdp" && d.HasChange(isVolumeIops) {
		volumeProfilePatchModel := &vpcv1.VolumePatch{}
		iops := int64(d.Get(isVolumeIops).(int))
		volumeProfilePatchModel.Iops = &iops
		volumeProfilePatch, err := volumeProfilePatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("[ERROR] Error calling asPatch for VolumeProfilePatch for sdp profiles : %s", err)
		}
		options.VolumePatch = volumeProfilePatch
		_, response, err = sess.UpdateVolume(options)
		if err != nil {
			return err
		}
		_, err = isWaitForVolumeAvailable(sess, d.Id(), d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return err
		}
		eTag = response.Headers.Get("ETag")
		options.IfMatch = &eTag
	} else if d.HasChange(isVolumeProfileName) || d.HasChange(isVolumeIops) {
		volumeProfilePatchModel := &vpcv1.VolumePatch{}
		volId := d.Id()
		getvoloptions := &vpcv1.GetVolumeOptions{
			ID: &volId,
		}
		vol, response, err := sess.GetVolume(getvoloptions)
		if err != nil || vol == nil {
			return fmt.Errorf("[ERROR] Error retrieving Volume (%s) details: %s\n%s", volId, err, response)
		}
		if vol.VolumeAttachments == nil || len(vol.VolumeAttachments) < 1 {
			return fmt.Errorf("[ERROR] Error updating Volume profile/iops because the specified volume %s is not attached to a virtual server instance ", volId)
		}
		volAtt := &vol.VolumeAttachments[0]
		insId := *volAtt.Instance.ID
		getinsOptions := &vpcv1.GetInstanceOptions{
			ID: &insId,
		}
		instance, response, err := sess.GetInstance(getinsOptions)
		if err != nil || instance == nil {
			return fmt.Errorf("[ERROR] Error retrieving Instance (%s) to which the volume (%s) is attached : %s\n%s", insId, volId, err, response)
		}
		if instance != nil && *instance.Status != "running" {
			actiontype := "start"
			createinsactoptions := &vpcv1.CreateInstanceActionOptions{
				InstanceID: &insId,
				Type:       &actiontype,
			}
			_, response, err = sess.CreateInstanceAction(createinsactoptions)
			if err != nil {
				return fmt.Errorf("[ERROR] Error starting Instance (%s) to which the volume (%s) is attached  : %s\n%s", insId, volId, err, response)
			}
			_, err = isWaitForInstanceAvailable(sess, insId, d.Timeout(schema.TimeoutCreate), d)
			if err != nil {
				return err
			}
		}
		if d.HasChange(isVolumeProfileName) {
			profile := d.Get(isVolumeProfileName).(string)
			volumeProfilePatchModel.Profile = &vpcv1.VolumeProfileIdentity{
				Name: &profile,
			}
		} else if d.HasChange(isVolumeIops) {
			profile := d.Get(isVolumeProfileName).(string)
			volumeProfilePatchModel.Profile = &vpcv1.VolumeProfileIdentity{
				Name: &profile,
			}
			iops := int64(d.Get(isVolumeIops).(int))
			volumeProfilePatchModel.Iops = &iops
		}

		volumeProfilePatch, err := volumeProfilePatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("[ERROR] Error calling asPatch for VolumeProfilePatch: %s", err)
		}
		options.VolumePatch = volumeProfilePatch
		_, response, err = sess.UpdateVolume(options)
		if err != nil {
			return err
		}
		eTag = response.Headers.Get("ETag")
		options.IfMatch = &eTag
		_, err = isWaitForVolumeAvailable(sess, d.Id(), d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return err
		}
	}

	// capacity update
	if d.HasChange(isVolumeCapacity) {
		id := d.Id()
		getvolumeoptions := &vpcv1.GetVolumeOptions{
			ID: &id,
		}
		vol, response, err := sess.GetVolume(getvolumeoptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("[ERROR] Error Getting Volume (%s): %s\n%s", id, err, response)
		}
		eTag = response.Headers.Get("ETag")
		options.IfMatch = &eTag
		if *vol.Profile.Name != "sdp" {
			if vol.VolumeAttachments == nil || len(vol.VolumeAttachments) == 0 || *vol.VolumeAttachments[0].ID == "" {
				return fmt.Errorf("[ERROR] Error volume capacity can't be updated since volume %s is not attached to any instance for VolumePatch", id)
			}
			insId := vol.VolumeAttachments[0].Instance.ID
			getinsOptions := &vpcv1.GetInstanceOptions{
				ID: insId,
			}
			instance, response, err := sess.GetInstance(getinsOptions)
			if err != nil || instance == nil {
				return fmt.Errorf("[ERROR] Error retrieving Instance (%s) : %s\n%s", *insId, err, response)
			}
			if instance != nil && *instance.Status != "running" {
				actiontype := "start"
				createinsactoptions := &vpcv1.CreateInstanceActionOptions{
					InstanceID: insId,
					Type:       &actiontype,
				}
				_, response, err = sess.CreateInstanceAction(createinsactoptions)
				if err != nil {
					return fmt.Errorf("[ERROR] Error starting Instance (%s) : %s\n%s", *insId, err, response)
				}
				_, err = isWaitForInstanceAvailable(sess, *insId, d.Timeout(schema.TimeoutCreate), d)
				if err != nil {
					return err
				}
			}
		}

		capacity = int64(d.Get(isVolumeCapacity).(int))
		volumeCapacityPatchModel := &vpcv1.VolumePatch{}
		volumeCapacityPatchModel.Capacity = &capacity

		volumeCapacityPatch, err := volumeCapacityPatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("[ERROR] Error calling asPatch for volumeCapacityPatch: %s", err)
		}
		options.VolumePatch = volumeCapacityPatch
		_, response, err = sess.UpdateVolume(options)
		if err != nil {
			return fmt.Errorf("[ERROR] Error updating vpc volume: %s\n%s", err, response)
		}
		_, err = isWaitForVolumeAvailable(sess, d.Id(), d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return err
		}
	}

	// user tags update
	if d.HasChange(isVolumeTags) {
		var userTags *schema.Set
		if v, ok := d.GetOk(isVolumeTags); ok {
			userTags = v.(*schema.Set)
			if userTags != nil && userTags.Len() != 0 {
				userTagsArray := make([]string, userTags.Len())
				for i, userTag := range userTags.List() {
					userTagStr := userTag.(string)
					userTagsArray[i] = userTagStr
				}
				schematicTags := os.Getenv("IC_ENV_TAGS")
				var envTags []string
				if schematicTags != "" {
					envTags = strings.Split(schematicTags, ",")
					userTagsArray = append(userTagsArray, envTags...)
				}
				volumeNamePatchModel := &vpcv1.VolumePatch{}
				volumeNamePatchModel.UserTags = userTagsArray
				volumeNamePatch, err := volumeNamePatchModel.AsPatch()
				if err != nil {
					return fmt.Errorf("Error calling asPatch for volumeNamePatch: %s", err)
				}
				options.IfMatch = &eTag
				options.VolumePatch = volumeNamePatch
				_, response, err := sess.UpdateVolume(options)
				if err != nil {
					return fmt.Errorf("Error updating volume : %s\n%s", err, response)
				}
				_, err = isWaitForVolumeAvailable(sess, d.Id(), d.Timeout(schema.TimeoutCreate))
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func resourceIBMISVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()

	err := volDelete(d, meta, id)
	if err != nil {
		return err
	}
	return nil
}

func volDelete(d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	getvoloptions := &vpcv1.GetVolumeOptions{
		ID: &id,
	}
	volDetails, response, err := sess.GetVolume(getvoloptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil
		}
		return fmt.Errorf("[ERROR] Error getting Volume (%s): %s\n%s", id, err, response)
	}

	if volDetails.VolumeAttachments != nil {
		for _, volAtt := range volDetails.VolumeAttachments {
			deleteVolumeAttachment := &vpcv1.DeleteInstanceVolumeAttachmentOptions{
				InstanceID: volAtt.Instance.ID,
				ID:         volAtt.ID,
			}
			_, err := sess.DeleteInstanceVolumeAttachment(deleteVolumeAttachment)
			if err != nil {
				return fmt.Errorf("[ERROR] Error while removing volume attachment %q for instance %s: %q", *volAtt.ID, *volAtt.Instance.ID, err)
			}
			_, err = isWaitForInstanceVolumeDetached(sess, d, d.Id(), *volAtt.ID)
			if err != nil {
				return err
			}

		}
	}

	options := &vpcv1.DeleteVolumeOptions{
		ID: &id,
	}
	response, err = sess.DeleteVolume(options)
	if err != nil {
		return fmt.Errorf("[ERROR] Error deleting Volume : %s\n%s", err, response)
	}
	_, err = isWaitForVolumeDeleted(sess, id, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func isWaitForVolumeDeleted(vol *vpcv1.VpcV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for  (%s) to be deleted.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isVolumeDeleting},
		Target:     []string{"done", ""},
		Refresh:    isVolumeDeleteRefreshFunc(vol, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isVolumeDeleteRefreshFunc(vol *vpcv1.VpcV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		volgetoptions := &vpcv1.GetVolumeOptions{
			ID: &id,
		}
		vol, response, err := vol.GetVolume(volgetoptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return vol, isVolumeDeleted, nil
			}
			return vol, "", fmt.Errorf("[ERROR] Error getting Volume: %s\n%s", err, response)
		}
		return vol, isVolumeDeleting, err
	}
}

func resourceIBMISVolumeExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	id := d.Id()

	exists, err := volExists(d, meta, id)
	return exists, err
}

func volExists(d *schema.ResourceData, meta interface{}, id string) (bool, error) {
	sess, err := vpcClient(meta)
	if err != nil {
		return false, err
	}
	options := &vpcv1.GetVolumeOptions{
		ID: &id,
	}
	_, response, err := sess.GetVolume(options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("[ERROR] Error getting Volume: %s\n%s", err, response)
	}
	return true, nil
}

func isWaitForVolumeAvailable(client *vpcv1.VpcV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for Volume (%s) to be available.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isVolumeProvisioning},
		Target:     []string{isVolumeProvisioningDone, ""},
		Refresh:    isVolumeRefreshFunc(client, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isVolumeRefreshFunc(client *vpcv1.VpcV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		volgetoptions := &vpcv1.GetVolumeOptions{
			ID: &id,
		}
		vol, response, err := client.GetVolume(volgetoptions)
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] Error getting volume: %s\n%s", err, response)
		}

		if *vol.Status == "available" {
			return vol, isVolumeProvisioningDone, nil
		}

		return vol, isVolumeProvisioning, nil
	}
}

func deleteAllSnapshots(sess *vpcv1.VpcV1, id string) error {
	delete_all_snapshots := new(vpcv1.DeleteSnapshotsOptions)
	delete_all_snapshots.SourceVolumeID = &id
	response, err := sess.DeleteSnapshots(delete_all_snapshots)
	if err != nil {
		return fmt.Errorf("[ERROR] Error deleting snapshots from volume %s\n%s", err, response)
	}
	return nil
}
