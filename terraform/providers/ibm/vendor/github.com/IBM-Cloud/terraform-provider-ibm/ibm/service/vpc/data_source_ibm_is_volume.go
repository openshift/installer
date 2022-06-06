// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMISVolume() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISVolumeRead,

		Schema: map[string]*schema.Schema{

			isVolumeName: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeDataSourceValidator("ibm_is_subnet", isVolumeName),
				Description:  "Volume name",
			},

			isVolumeZone: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Zone name",
			},

			isVolumeBandwidth: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum bandwidth (in megabits per second) for the volume",
			},

			isVolumeResourceGroup: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource group name",
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
			Identifier:                 isVolumeName,
			ValidateFunctionIdentifier: validate.ValidateNoZeroValues,
			Type:                       validate.TypeString})

	ibmISVoulmeDataSourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_volume", Schema: validateSchema}
	return &ibmISVoulmeDataSourceValidator
}

func dataSourceIBMISVolumeRead(d *schema.ResourceData, meta interface{}) error {

	name := d.Get(isVolumeName).(string)

	err := volumeGet(d, meta, name)
	if err != nil {
		return err
	}
	return nil
}

func volumeGet(d *schema.ResourceData, meta interface{}, name string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	zone := ""
	if zname, ok := d.GetOk(isVolumeZone); ok {
		zone = zname.(string)
	}
	start := ""
	allrecs := []vpcv1.Volume{}
	for {
		listVolumesOptions := &vpcv1.ListVolumesOptions{}
		if start != "" {
			listVolumesOptions.Start = &start
		}
		if zone != "" {
			listVolumesOptions.ZoneName = &zone
		}
		listVolumesOptions.Name = &name
		vols, response, err := sess.ListVolumes(listVolumesOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Fetching volumes %s\n%s", err, response)
		}
		start = flex.GetNext(vols.Next)
		allrecs = append(allrecs, vols.Volumes...)
		if start == "" {
			break
		}
	}
	for _, vol := range allrecs {
		d.SetId(*vol.ID)
		d.Set(isVolumeBandwidth, int(*vol.Bandwidth))
		d.Set(isVolumeName, *vol.Name)
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
					statusReasonsList = append(statusReasonsList, currentSR)
				}
			}
			d.Set(isVolumeStatusReasons, statusReasonsList)
		}
		tags, err := flex.GetTagsUsingCRN(meta, *vol.CRN)
		if err != nil {
			log.Printf(
				"Error on get of resource vpc volume (%s) tags: %s", d.Id(), err)
		}
		d.Set(isVolumeTags, tags)
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
		return nil
	}
	return fmt.Errorf("[ERROR] No Volume found with name %s", name)
}
