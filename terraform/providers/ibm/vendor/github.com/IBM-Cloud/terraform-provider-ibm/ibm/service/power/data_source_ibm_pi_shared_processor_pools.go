// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"log"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPISharedProcessorPools() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPISharedProcessorPoolsRead,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_SharedProcessorPools: {
				Computed:    true,
				Description: "List of all the shared processor pools.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_AllocatedCores: {
							Computed:    true,
							Description: "The allocated cores in the shared processor pool.",
							Type:        schema.TypeFloat,
						},
						Attr_AvailableCores: {
							Computed:    true,
							Description: "The available cores in the shared processor pool.",
							Type:        schema.TypeInt,
						},
						Attr_CRN: {
							Computed:    true,
							Description: "The CRN of this resource.",
							Type:        schema.TypeString,
						},
						Attr_DedicatedHostID: {
							Computed:    true,
							Description: "The dedicated host ID where the shared processor pool resides.",
							Type:        schema.TypeString,
						},
						Attr_HostID: {
							Computed:    true,
							Description: "The host ID where the shared processor pool resides.",
							Type:        schema.TypeInt,
						},
						Attr_Name: {
							Computed:    true,
							Description: "The name of the shared processor pool.",
							Type:        schema.TypeString,
						},
						Attr_ReservedCores: {
							Computed:    true,
							Description: "The amount of reserved cores for the shared processor pool.",
							Type:        schema.TypeInt,
						},
						Attr_SharedProcessorPoolID: {
							Computed:    true,
							Description: "The shared processor pool's unique ID.",
							Type:        schema.TypeString,
						},
						Attr_Status: {
							Computed:    true,
							Description: "The status of the shared processor pool.",
							Type:        schema.TypeString,
						},
						Attr_StatusDetail: {
							Computed:    true,
							Description: "The status details of the shared processor pool.",
							Type:        schema.TypeString,
						},
						Attr_UserTags: {
							Computed:    true,
							Description: "List of user tags attached to the resource.",
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         schema.HashString,
							Type:        schema.TypeSet,
						},
					},
				},
				Type: schema.TypeList,
			},
		},
	}
}

func dataSourceIBMPISharedProcessorPoolsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)

	client := instance.NewIBMPISharedProcessorPoolClient(ctx, sess, cloudInstanceID)
	pools, err := client.GetAll()
	if err != nil || pools == nil {
		return diag.Errorf("error fetching shared processor pools: %v", err)
	}

	result := make([]map[string]interface{}, 0, len(pools.SharedProcessorPools))
	for _, pool := range pools.SharedProcessorPools {
		key := map[string]interface{}{
			Attr_AllocatedCores:        *pool.AllocatedCores,
			Attr_AvailableCores:        *pool.AvailableCores,
			Attr_DedicatedHostID:       pool.DedicatedHostID,
			Attr_HostID:                pool.HostID,
			Attr_Name:                  *pool.Name,
			Attr_ReservedCores:         *pool.ReservedCores,
			Attr_SharedProcessorPoolID: *pool.ID,
			Attr_Status:                pool.Status,
			Attr_StatusDetail:          pool.StatusDetail,
		}
		if pool.Crn != "" {
			key[Attr_CRN] = pool.Crn
			tags, err := flex.GetGlobalTagsUsingCRN(meta, string(pool.Crn), "", UserTagType)
			if err != nil {
				log.Printf("Error on get of pi shared_processor_pool (%s) user_tags: %s", *pool.ID, err)
			}
			key[Attr_UserTags] = tags
		}
		result = append(result, key)
	}

	var genID, _ = uuid.GenerateUUID()
	d.SetId(genID)
	d.Set(Attr_SharedProcessorPools, result)

	return nil
}
