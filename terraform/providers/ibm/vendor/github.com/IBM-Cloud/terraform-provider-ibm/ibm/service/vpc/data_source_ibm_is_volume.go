// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMISVolume() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMISVolumeRead,

		Schema: map[string]*schema.Schema{

			isVolumeName: {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{isVolumeName, "identifier"},
				ValidateFunc: validate.InvokeDataSourceValidator("ibm_is_volume", isVolumeName),
				Description:  "Volume name",
			},
			"identifier": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{isVolumeName, "identifier"},
				ValidateFunc: validate.InvokeDataSourceValidator("ibm_is_volume", "identifier"),
				Description:  "Volume name",
			},

			isVolumeZone: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Zone name",
			},
			isVolumesActive: &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether a running virtual server instance has an attachment to this volume.",
			},
			// defined_performance changes
			"adjustable_capacity_states": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The attachment states that support adjustable capacity for this volume.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"adjustable_iops_states": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The attachment states that support adjustable IOPS for this volume.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			isVolumeAttachmentState: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The attachment state of the volume.",
			},
			isVolumeBandwidth: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum bandwidth (in megabits per second) for the volume",
			},
			isVolumesBusy: &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether this volume is performing an operation that must be serialized. If an operation specifies that it requires serialization, the operation will fail unless this property is `false`.",
			},
			isVolumesCreatedAt: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the volume was created.",
			},
			isVolumeResourceGroup: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource group name",
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
			isVolumeProfileName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Volume profile name",
			},

			isVolumeEncryptionKey: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Volume encryption key info",
			},

			isVolumeEncryptionType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Volume encryption type info",
			},

			isVolumeCapacity: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Vloume capacity value",
			},

			isVolumeIops: {
				Type:        schema.TypeInt,
				Computed:    true,
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

			isVolumeHealthState: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The health of this resource.",
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

			isVolumeTags: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         flex.ResourceIBMVPCHash,
				Description: "Tags for the volume instance",
			},

			isVolumeAccessTags: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         flex.ResourceIBMVPCHash,
				Description: "Access management tags for the volume instance",
			},

			isVolumeSourceSnapshot: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identifier of the snapshot from which this volume was cloned",
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
		},
	}
}

func DataSourceIBMISVolumeValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "identifier",
			ValidateFunctionIdentifier: validate.ValidateNoZeroValues,
			Type:                       validate.TypeString})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isVolumeName,
			ValidateFunctionIdentifier: validate.ValidateNoZeroValues,
			Type:                       validate.TypeString})

	ibmISVoulmeDataSourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_volume", Schema: validateSchema}
	return &ibmISVoulmeDataSourceValidator
}

func dataSourceIBMISVolumeRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	err := volumeGet(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func volumeGet(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	var vol vpcv1.Volume
	if volName, ok := d.GetOk(isVolumeName); ok {
		name := volName.(string)
		zone := ""
		if zname, ok := d.GetOk(isVolumeZone); ok {
			zone = zname.(string)
		}
		listVolumesOptions := &vpcv1.ListVolumesOptions{
			Name: &name,
		}

		if zone != "" {
			listVolumesOptions.ZoneName = &zone
		}
		listVolumesOptions.Name = &name
		vols, response, err := sess.ListVolumes(listVolumesOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Fetching volumes %s\n%s", err, response)
		}
		allrecs := vols.Volumes

		if len(allrecs) == 0 {
			return fmt.Errorf("[ERROR] No Volume found with name %s", name)
		}
		vol = allrecs[0]
	} else {
		identifier := d.Get("identifier").(string)
		getVolumeOptions := &vpcv1.GetVolumeOptions{
			ID: &identifier,
		}

		volPtr, response, err := sess.GetVolume(getVolumeOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error on get volume %s\n%s", err, response)
		}
		vol = *volPtr
	}
	d.SetId(*vol.ID)
	if vol.Active != nil {
		d.Set(isVolumesActive, vol.Active)
	}
	if vol.AttachmentState != nil {
		d.Set(isVolumeAttachmentState, vol.AttachmentState)
	}
	d.Set(isVolumeBandwidth, int(*vol.Bandwidth))
	if vol.Busy != nil {
		d.Set(isVolumesBusy, vol.Busy)
	}
	if vol.Capacity != nil {
		d.Set(isVolumesCapacity, vol.Capacity)
	}
	if vol.CreatedAt != nil {
		d.Set(isVolumesCreatedAt, flex.DateTimeToString(vol.CreatedAt))
	}
	d.Set(isVolumeName, *vol.Name)
	d.Set("identifier", *vol.ID)
	if vol.OperatingSystem != nil {
		operatingSystemList := []map[string]interface{}{}
		operatingSystemMap := dataSourceVolumeCollectionVolumesOperatingSystemToMap(*vol.OperatingSystem)
		operatingSystemList = append(operatingSystemList, operatingSystemMap)
		d.Set(isVolumesOperatingSystem, operatingSystemList)
	}
	d.Set(isVolumeProfileName, *vol.Profile.Name)
	d.Set(isVolumeZone, *vol.Zone.Name)
	if vol.EncryptionKey != nil {
		d.Set(isVolumeEncryptionKey, vol.EncryptionKey.CRN)
	}
	if vol.Encryption != nil {
		d.Set(isVolumeEncryptionType, vol.Encryption)
	}
	if vol.SourceSnapshot != nil {
		d.Set(isVolumeSourceSnapshot, *vol.SourceSnapshot.ID)
	}
	d.Set(isVolumeIops, *vol.Iops)
	d.Set(isVolumeCapacity, *vol.Capacity)
	d.Set(isVolumeCrn, *vol.CRN)
	d.Set(isVolumeStatus, *vol.Status)
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
			d.Set(isVolumeStatusReasons, statusReasonsList)
		}
	}
	d.Set(isVolumeTags, vol.UserTags)
	accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *vol.CRN, "", isVolumeAccessTagType)
	if err != nil {
		log.Printf(
			"Error on get of resource vpc volume (%s) access tags: %s", d.Id(), err)
	}
	d.Set(isVolumeAccessTags, accesstags)
	controller, err := flex.GetBaseController(meta)
	if err != nil {
		return err
	}
	d.Set(flex.ResourceControllerURL, controller+"/vpc-ext/storage/storageVolumes")
	d.Set(flex.ResourceName, *vol.Name)
	d.Set(flex.ResourceCRN, *vol.CRN)
	d.Set(flex.ResourceStatus, *vol.Status)
	if vol.ResourceGroup != nil {
		d.Set(flex.ResourceGroupName, vol.ResourceGroup.Name)
		d.Set(isVolumeResourceGroup, *vol.ResourceGroup.ID)
	}

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
	if vol.CatalogOffering != nil {
		versionCrn := ""
		if vol.CatalogOffering.Version != nil && vol.CatalogOffering.Version.CRN != nil {
			versionCrn = *vol.CatalogOffering.Version.CRN
		}
		catalogList := make([]map[string]interface{}, 0)
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
		d.Set(isVolumeCatalogOffering, catalogList)
	}
	if vol.HealthState != nil {
		d.Set(isVolumeHealthState, *vol.HealthState)
	}

	if err = d.Set("adjustable_capacity_states", vol.AdjustableCapacityStates); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting adjustable_capacity_states: %s", err), "(Data) ibm_is_volume", "read", "set-adjustable_capacity_states")
	}

	if err = d.Set("adjustable_iops_states", vol.AdjustableIopsStates); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting adjustable_iops_states: %s", err), "(Data) ibm_is_volume", "read", "set-adjustable_iops_states")
	}

	return nil
}

func resourceIbmIsVolumeCatalogOfferingVersionPlanReferenceDeletedToMap(catalogOfferingVersionPlanReferenceDeleted vpcv1.Deleted) map[string]interface{} {
	catalogOfferingVersionPlanReferenceDeletedMap := map[string]interface{}{}

	catalogOfferingVersionPlanReferenceDeletedMap["more_info"] = catalogOfferingVersionPlanReferenceDeleted.MoreInfo

	return catalogOfferingVersionPlanReferenceDeletedMap
}
