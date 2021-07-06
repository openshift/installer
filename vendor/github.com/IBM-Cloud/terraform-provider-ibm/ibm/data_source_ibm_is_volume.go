// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"

	"github.com/IBM/vpc-go-sdk/vpcclassicv1"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMISVolume() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISVolumeRead,

		Schema: map[string]*schema.Schema{

			isVolumeName: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: InvokeDataSourceValidator("ibm_is_subnet", isVolumeName),
				Description:  "Volume name",
			},

			isVolumeZone: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Zone name",
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
				Set:         resourceIBMVPCHash,
				Description: "Tags for the volume instance",
			},

			ResourceControllerURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the IBM Cloud dashboard that can be used to explore and view details about this instance",
			},

			ResourceName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the resource",
			},

			ResourceCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},

			ResourceStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the resource",
			},

			ResourceGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group name in which resource is provisioned",
			},
		},
	}
}

func dataSourceIBMISVolumeValidator() *ResourceValidator {
	validateSchema := make([]ValidateSchema, 1)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isVolumeName,
			ValidateFunctionIdentifier: ValidateNoZeroValues,
			Type:                       TypeString})

	ibmISVoulmeDataSourceValidator := ResourceValidator{ResourceName: "ibm_is_volume", Schema: validateSchema}
	return &ibmISVoulmeDataSourceValidator
}

func dataSourceIBMISVolumeRead(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	name := d.Get(isVolumeName).(string)
	if userDetails.generation == 1 {
		err := classicVolumeGet(d, meta, name)
		if err != nil {
			return err
		}
	} else {
		err := volumeGet(d, meta, name)
		if err != nil {
			return err
		}
	}
	return nil
}

func classicVolumeGet(d *schema.ResourceData, meta interface{}, name string) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	zone := ""
	if zname, ok := d.GetOk(isVolumeZone); ok {
		zone = zname.(string)
	}
	start := ""
	allrecs := []vpcclassicv1.Volume{}
	for {
		listVolumesOptions := &vpcclassicv1.ListVolumesOptions{}
		if start != "" {
			listVolumesOptions.Start = &start
		}
		if zone != "" {
			listVolumesOptions.ZoneName = &zone
		}
		listVolumesOptions.Name = &name
		vols, response, err := sess.ListVolumes(listVolumesOptions)
		if err != nil {
			return fmt.Errorf("Error Fetching volumes %s\n%s", err, response)
		}
		start = GetNext(vols.Next)
		allrecs = append(allrecs, vols.Volumes...)
		if start == "" {
			break
		}
	}
	for _, vol := range allrecs {
		d.SetId(*vol.ID)
		d.Set(isVolumeName, *vol.Name)
		d.Set(isVolumeProfileName, *vol.Profile.Name)
		d.Set(isVolumeZone, *vol.Zone.Name)
		if vol.EncryptionKey != nil {
			d.Set(isVolumeEncryptionKey, *vol.EncryptionKey.CRN)
		}
		d.Set(isVolumeIops, *vol.Iops)
		d.Set(isVolumeCapacity, *vol.Capacity)
		d.Set(isVolumeCrn, *vol.CRN)
		d.Set(isVolumeStatus, *vol.Status)
		tags, err := GetTagsUsingCRN(meta, *vol.CRN)
		if err != nil {
			log.Printf(
				"Error on get of resource vpc volume (%s) tags: %s", d.Id(), err)
		}
		d.Set(isVolumeTags, tags)
		controller, err := getBaseController(meta)
		if err != nil {
			return err
		}
		d.Set(ResourceControllerURL, controller+"/vpc/storage/storageVolumes")
		d.Set(ResourceName, *vol.Name)
		d.Set(ResourceCRN, *vol.CRN)
		d.Set(ResourceStatus, *vol.Status)
		if vol.ResourceGroup != nil {
			d.Set(ResourceGroupName, *vol.ResourceGroup.ID)
			d.Set(isVolumeResourceGroup, *vol.ResourceGroup.ID)
		}
		return nil
	}
	return fmt.Errorf("No Volume found with name %s", name)
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
			return fmt.Errorf("Error Fetching volumes %s\n%s", err, response)
		}
		start = GetNext(vols.Next)
		allrecs = append(allrecs, vols.Volumes...)
		if start == "" {
			break
		}
	}
	for _, vol := range allrecs {
		d.SetId(*vol.ID)
		d.Set(isVolumeName, *vol.Name)
		d.Set(isVolumeProfileName, *vol.Profile.Name)
		d.Set(isVolumeZone, *vol.Zone.Name)
		if vol.EncryptionKey != nil {
			d.Set(isVolumeEncryptionKey, vol.EncryptionKey.CRN)
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
		tags, err := GetTagsUsingCRN(meta, *vol.CRN)
		if err != nil {
			log.Printf(
				"Error on get of resource vpc volume (%s) tags: %s", d.Id(), err)
		}
		d.Set(isVolumeTags, tags)
		controller, err := getBaseController(meta)
		if err != nil {
			return err
		}
		d.Set(ResourceControllerURL, controller+"/vpc-ext/storage/storageVolumes")
		d.Set(ResourceName, *vol.Name)
		d.Set(ResourceCRN, *vol.CRN)
		d.Set(ResourceStatus, *vol.Status)
		if vol.ResourceGroup != nil {
			d.Set(ResourceGroupName, *vol.ResourceGroup.Name)
			d.Set(isVolumeResourceGroup, *vol.ResourceGroup.ID)
		}
		return nil
	}
	return fmt.Errorf("No Volume found with name %s", name)
}
