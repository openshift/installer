// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"log"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPICloudInstance() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPICloudInstanceRead,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_Capabilities: {
				Computed:    true,
				Description: "Lists the capabilities for this cloud instance.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Type:        schema.TypeList,
			},
			Attr_Enabled: {
				Computed:    true,
				Description: "Indicates whether the tenant is enabled.",
				Type:        schema.TypeBool,
			},
			Attr_PVMInstances: {
				Computed:    true,
				Description: "PVM instances owned by the Cloud Instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_CreationDate: {
							Computed:    true,
							Description: "Date of PVM instance creation.",
							Type:        schema.TypeString,
						},
						Attr_CRN: {
							Computed:    true,
							Description: "The CRN of this resource.",
							Type:        schema.TypeString,
						},
						Attr_Href: {
							Computed:    true,
							Description: "Link to Cloud Instance resource.",
							Type:        schema.TypeString,
						},
						Attr_ID: {
							Computed:    true,
							Description: "PVM Instance ID.",
							Type:        schema.TypeString,
						},
						Attr_Name: {
							Computed:    true,
							Description: "Name of the server.",
							Type:        schema.TypeString,
						},
						Attr_Status: {
							Computed:    true,
							Description: "The status of the instance.",
							Type:        schema.TypeString,
						},
						Attr_Systype: {
							Computed:    true,
							Description: "System type used to host the instance.",
							Type:        schema.TypeString,
						},
					},
				},
				Type: schema.TypeList,
			},
			Attr_Region: {
				Computed:    true,
				Description: "The region the cloud instance lives.",
				Type:        schema.TypeString,
			},
			Attr_TenantID: {
				Computed:    true,
				Description: "The tenant ID that owns this cloud instance.",
				Type:        schema.TypeString,
			},
			Attr_TotalInstances: {
				Computed:    true,
				Description: "The count of lpars that belong to this specific cloud instance.",
				Type:        schema.TypeFloat,
			},
			Attr_TotalMemoryConsumed: {
				Computed:    true,
				Description: "The total memory consumed by this service instance.",
				Type:        schema.TypeFloat,
			},
			Attr_TotalProcessorsConsumed: {
				Computed:    true,
				Description: "The total processors consumed by this service instance.",
				Type:        schema.TypeFloat,
			},
			Attr_TotalSSDStorageConsumed: {
				Computed:    true,
				Description: "The total SSD Storage consumed by this service instance.",
				Type:        schema.TypeFloat,
			},
			Attr_TotalStandardStorageConsumed: {
				Computed:    true,
				Description: "The total Standard Storage consumed by this service instance.",
				Type:        schema.TypeFloat,
			},
		},
	}
}

func dataSourceIBMPICloudInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)

	cloud_instance := instance.NewIBMPICloudInstanceClient(ctx, sess, cloudInstanceID)
	cloud_instance_data, err := cloud_instance.Get(cloudInstanceID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*cloud_instance_data.CloudInstanceID)

	d.Set(Attr_Capabilities, cloud_instance_data.Capabilities)
	d.Set(Attr_Enabled, cloud_instance_data.Enabled)
	d.Set(Attr_PVMInstances, flattenpvminstances(cloud_instance_data.PvmInstances, meta))
	d.Set(Attr_Region, cloud_instance_data.Region)
	d.Set(Attr_TenantID, (cloud_instance_data.TenantID))
	d.Set(Attr_TotalInstances, cloud_instance_data.Usage.Instances)
	d.Set(Attr_TotalMemoryConsumed, cloud_instance_data.Usage.Memory)
	d.Set(Attr_TotalProcessorsConsumed, cloud_instance_data.Usage.Processors)
	d.Set(Attr_TotalSSDStorageConsumed, cloud_instance_data.Usage.StorageSSD)
	d.Set(Attr_TotalStandardStorageConsumed, cloud_instance_data.Usage.StorageStandard)

	return nil
}

func flattenpvminstances(list []*models.PVMInstanceReference, meta interface{}) []map[string]interface{} {
	pvms := make([]map[string]interface{}, 0)
	for _, lpars := range list {
		l := map[string]interface{}{
			Attr_CreationDate: lpars.CreationDate.String(),
			Attr_ID:           *lpars.PvmInstanceID,
			Attr_Href:         *lpars.Href,
			Attr_Name:         *lpars.ServerName,
			Attr_Status:       *lpars.Status,
			Attr_Systype:      lpars.SysType,
		}
		if lpars.Crn != "" {
			l[Attr_CRN] = lpars.Crn
			tags, err := flex.GetGlobalTagsUsingCRN(meta, string(lpars.Crn), "", UserTagType)
			if err != nil {
				log.Printf("Error on get of pi instance (%s) user_tags: %s", *lpars.PvmInstanceID, err)
			}
			l[Attr_UserTags] = tags
		}
		pvms = append(pvms, l)
	}
	return pvms
}
