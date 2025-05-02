// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
)

const (
	PISharedProcessorPools = "shared_processor_pools"
)

func DataSourceIBMPISharedProcessorPools() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPISharedProcessorPoolsRead,
		Schema: map[string]*schema.Schema{
			helpers.PICloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "PI cloud instance ID",
				ValidateFunc: validation.NoZeroValues,
			},
			// Computed Attributes
			PISharedProcessorPools: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_SharedProcessorPoolID: {
							Type:     schema.TypeString,
							Computed: true,
						},
						Attr_SharedProcessorPoolAllocatedCores: {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						Attr_SharedProcessorPoolAvailableCores: {
							Type:     schema.TypeInt,
							Computed: true,
						},
						Attr_SharedProcessorPoolName: {
							Type:     schema.TypeString,
							Computed: true,
						},
						Attr_SharedProcessorPoolReservedCores: {
							Type:     schema.TypeInt,
							Computed: true,
						},
						Attr_SharedProcessorPoolHostID: {
							Type:     schema.TypeInt,
							Computed: true,
						},
						Attr_SharedProcessorPoolStatus: {
							Type:     schema.TypeString,
							Computed: true,
						},
						Attr_SharedProcessorPoolStatusDetail: {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMPISharedProcessorPoolsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)

	client := st.NewIBMPISharedProcessorPoolClient(ctx, sess, cloudInstanceID)
	pools, err := client.GetAll()
	if err != nil || pools == nil {
		return diag.Errorf("error fetching shared processor pools: %v", err)
	}

	result := make([]map[string]interface{}, 0, len(pools.SharedProcessorPools))
	for _, pool := range pools.SharedProcessorPools {
		key := map[string]interface{}{
			Attr_SharedProcessorPoolID:             *pool.ID,
			Attr_SharedProcessorPoolName:           *pool.Name,
			Attr_SharedProcessorPoolAllocatedCores: *pool.AllocatedCores,
			Attr_SharedProcessorPoolAvailableCores: *pool.AvailableCores,
			Attr_SharedProcessorPoolReservedCores:  *pool.ReservedCores,
			Attr_SharedProcessorPoolHostID:         pool.HostID,
			Attr_SharedProcessorPoolStatus:         pool.Status,
			Attr_SharedProcessorPoolStatusDetail:   pool.StatusDetail,
		}
		result = append(result, key)
	}

	var genID, _ = uuid.GenerateUUID()
	d.SetId(genID)
	d.Set(PISharedProcessorPools, result)

	return nil
}
