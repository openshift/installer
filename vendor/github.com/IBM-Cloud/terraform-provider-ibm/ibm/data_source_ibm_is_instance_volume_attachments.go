// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMISInstanceVolumeAttachments() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISInstanceVolumeAttachmentsRead,

		Schema: map[string]*schema.Schema{
			isInstanceId: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance id",
			},

			isInstanceVolumeAttachments: {
				Type:        schema.TypeList,
				Description: "List of volume attachments on an instance",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isInstanceVolumeDeleteOnInstanceDelete: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "If set to true, when deleting the instance the volume will also be deleted.",
						},
						isInstanceVolAttName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this volume attachment.",
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
						isInstanceVolumeAttVolumeReference: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The attached volume",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMISInstanceVolumeAttachmentsRead(d *schema.ResourceData, meta interface{}) error {
	instanceId := d.Get(isInstanceId).(string)

	err := instanceGetVolumeAttachments(d, meta, instanceId)
	if err != nil {
		return err
	}

	return nil
}

func instanceGetVolumeAttachments(d *schema.ResourceData, meta interface{}, instanceId string) error {
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
		return fmt.Errorf("Error Fetching Instance volume attachments %s\n%s", err, response)
	}
	allrecs = append(allrecs, volumeAtts.VolumeAttachments...)
	volAttList := make([]map[string]interface{}, 0)
	for _, volumeAtt := range allrecs {
		currentVolAtt := map[string]interface{}{}
		currentVolAtt[isInstanceVolAttName] = *volumeAtt.Name
		currentVolAtt[isInstanceVolumeDeleteOnInstanceDelete] = *volumeAtt.DeleteVolumeOnInstanceDelete
		currentVolAtt[isInstanceVolumeAttDevice] = *volumeAtt.Device.ID
		currentVolAtt[isInstanceVolumeAttHref] = *volumeAtt.Href
		currentVolAtt[isInstanceVolAttId] = *volumeAtt.ID
		currentVolAtt[isInstanceVolumeAttStatus] = *volumeAtt.Status
		currentVolAtt[isInstanceVolumeAttType] = *volumeAtt.Type

		if volumeAtt.Volume != nil {
			currentVolAtt[isInstanceVolumeAttVolumeReferenceId] = *volumeAtt.Volume.ID
			currentVolAtt[isInstanceVolumeAttVolumeReferenceName] = *volumeAtt.Volume.Name
			currentVolAtt[isInstanceVolumeAttVolumeReferenceCrn] = *volumeAtt.Volume.CRN
			if volumeAtt.Volume.Deleted != nil {
				currentVolAtt[isInstanceVolumeAttVolumeReferenceDeleted] = *volumeAtt.Volume.Deleted.MoreInfo
			}
			currentVolAtt[isInstanceVolumeAttVolumeReferenceHref] = *volumeAtt.Volume.Href
		}

		volAttList = append(volAttList, currentVolAtt)
	}
	d.SetId(dataSourceIBMISInstanceVolumeAttachmentsID(d))
	d.Set(isInstanceVolumeAttachments, volAttList)
	return nil
}

// dataSourceIBMISInstanceVolumeAttachmentsID returns a reasonable ID for a Instance volume attachments list.
func dataSourceIBMISInstanceVolumeAttachmentsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
