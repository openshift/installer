// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"

	"log"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"

	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	MaximumStorageAllocation = "maximum_storage_allocation"
	StoragePoolsCapacity     = "storage_pools_capacity"
	MaxAllocationSize        = "max_allocation_size"
	PoolName                 = "pool_name"
	StoragePool              = "storage_pool"
	StorageType              = "storage_type"
	TotalCapacity            = "total_capacity"
)

func DataSourceIBMPIStoragePoolsCapacity() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIStoragePoolsCapacityRead,
		Schema: map[string]*schema.Schema{
			helpers.PICloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			// Computed Attributes
			MaximumStorageAllocation: {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Maximum storage allocation",
			},
			StoragePoolsCapacity: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Storage pools capacity",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						MaxAllocationSize: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum allocation storage size (GB)",
						},
						PoolName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Pool name",
						},
						StorageType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Storage type of the storage pool",
						},
						TotalCapacity: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total pool capacity (GB)",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMPIStoragePoolsCapacityRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)

	client := st.NewIBMPIStorageCapacityClient(ctx, sess, cloudInstanceID)
	spc, err := client.GetAllStoragePoolsCapacity()
	if err != nil {
		log.Printf("[ERROR] get all storage pools capacity failed %v", err)
		return diag.FromErr(err)
	}

	var genID, _ = uuid.GenerateUUID()
	d.SetId(genID)

	if spc.MaximumStorageAllocation != nil {
		msa := spc.MaximumStorageAllocation
		data := map[string]interface{}{
			MaxAllocationSize: msa.MaxAllocationSize,
			StoragePool:       msa.StoragePool,
			StorageType:       msa.StorageType,
		}
		d.Set(MaximumStorageAllocation, flex.Flatten(data))
	}

	result := make([]map[string]string, 0, len(spc.StoragePoolsCapacity))
	for _, sp := range spc.StoragePoolsCapacity {
		data := map[string]interface{}{
			MaxAllocationSize: sp.MaxAllocationSize,
			PoolName:          sp.PoolName,
			StorageType:       sp.StorageType,
			TotalCapacity:     sp.TotalCapacity,
		}
		result = append(result, flex.Flatten(data))
	}
	d.Set(StoragePoolsCapacity, result)

	return nil
}
