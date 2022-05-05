// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isInstanceVolumeAttDevice                 = "device"
	isInstanceVolumeAttHref                   = "href"
	isInstanceVolumeAttStatus                 = "status"
	isInstanceVolumeAttType                   = "type"
	isInstanceVolumeAttVolumeReference        = "volume_reference"
	isInstanceVolumeAttVolumeReferenceCrn     = "volume_crn"
	isInstanceVolumeAttVolumeReferenceDeleted = "volume_deleted"
	isInstanceVolumeAttVolumeReferenceHref    = "volume_href"
	isInstanceVolumeAttVolumeReferenceId      = "volume_id"
	isInstanceVolumeAttVolumeReferenceName    = "volume_name"
)

func DataSourceIBMISInstanceVolumeAttachment() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISInstanceVolumeAttachmentRead,

		Schema: map[string]*schema.Schema{

			isInstanceId: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance id",
			},
			isInstanceVolAttName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The user-defined name for this volume attachment.",
			},

			isInstanceVolumeDeleteOnInstanceDelete: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If set to true, when deleting the instance the volume will also be deleted.",
			},

			isInstanceVolumeAttDevice: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A unique identifier for the device which is exposed to the instance operating system",
			},

			isInstanceVolumeAttHref: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this volume attachment",
			},

			isInstanceVolAttId: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for this volume attachment",
			},

			isInstanceVolumeAttStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of this volume attachment, one of [ attached, attaching, deleting, detaching ]",
			},

			isInstanceVolumeAttType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of volume attachment one of [ boot, data ]",
			},

			isInstanceVolumeAttVolumeReference: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The attached volume",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isInstanceVolumeAttVolumeReferenceId: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this volume",
						},
						isInstanceVolumeAttVolumeReferenceName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique user-defined name for this volume",
						},
						isInstanceVolumeAttVolumeReferenceCrn: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this volume",
						},
						isInstanceVolumeAttVolumeReferenceDeleted: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Link to documentation about deleted resources",
						},
						isInstanceVolumeAttVolumeReferenceHref: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this volume",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMISInstanceVolumeAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	instanceId := d.Get(isInstanceId).(string)
	name := d.Get(isInstanceName).(string)
	err := instanceVolumeAttachmentGetByName(d, meta, instanceId, name)
	if err != nil {
		return err
	}
	return nil
}

func instanceVolumeAttachmentGetByName(d *schema.ResourceData, meta interface{}, instanceId, name string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	allrecs := []vpcv1.VolumeAttachment{}
	listInstanceVolumeAttOptions := &vpcv1.ListInstanceVolumeAttachmentsOptions{
		InstanceID: &instanceId,
	}
	volumeAtts, response, err := sess.ListInstanceVolumeAttachments(listInstanceVolumeAttOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error fetching Instance volume attachments %s\n%s", err, response)
	}
	allrecs = append(allrecs, volumeAtts.VolumeAttachments...)
	for _, volumeAtt := range allrecs {
		if *volumeAtt.Name == name {
			d.SetId(makeTerraformVolAttID(instanceId, *volumeAtt.ID))
			d.Set(isInstanceVolAttName, *volumeAtt.Name)
			d.Set(isInstanceVolumeDeleteOnInstanceDelete, *volumeAtt.DeleteVolumeOnInstanceDelete)
			d.Set(isInstanceVolumeAttDevice, *volumeAtt.Device.ID)
			d.Set(isInstanceVolumeAttHref, *volumeAtt.Href)
			d.Set(isInstanceVolAttId, *volumeAtt.ID)
			d.Set(isInstanceVolumeAttStatus, *volumeAtt.Status)
			d.Set(isInstanceVolumeAttType, *volumeAtt.Type)
			volList := make([]map[string]interface{}, 0)
			if volumeAtt.Volume != nil {
				currentVol := map[string]interface{}{}
				currentVol[isInstanceVolumeAttVolumeReferenceId] = *volumeAtt.Volume.ID
				currentVol[isInstanceVolumeAttVolumeReferenceName] = *volumeAtt.Volume.Name
				currentVol[isInstanceVolumeAttVolumeReferenceCrn] = *volumeAtt.Volume.CRN
				if volumeAtt.Volume.Deleted != nil {
					currentVol[isInstanceVolumeAttVolumeReferenceDeleted] = *volumeAtt.Volume.Deleted.MoreInfo
				}
				currentVol[isInstanceVolumeAttVolumeReferenceHref] = *volumeAtt.Volume.Href
				volList = append(volList, currentVol)
			}
			d.Set(isInstanceVolumeAttVolumeReference, volList)
			return nil
		}
	}
	return fmt.Errorf("[ERROR] No Instance volume attachment found with name %s on instance %s", name, instanceId)
}
