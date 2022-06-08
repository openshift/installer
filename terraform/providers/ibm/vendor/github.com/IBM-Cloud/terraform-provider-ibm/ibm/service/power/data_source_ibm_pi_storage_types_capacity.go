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
	StorageTypesCapacity = "storage_types_capacity"
)

func DataSourceIBMPIStorageTypesCapacity() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIStorageTypesCapacityRead,
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
			StorageTypesCapacity: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Storage types capacity",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						StorageType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The storage type",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMPIStorageTypesCapacityRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)

	client := st.NewIBMPIStorageCapacityClient(ctx, sess, cloudInstanceID)
	stc, err := client.GetAllStorageTypesCapacity()
	if err != nil {
		log.Printf("[ERROR] get all storage types capacity failed %v", err)
		return diag.FromErr(err)
	}

	var genID, _ = uuid.GenerateUUID()
	d.SetId(genID)

	if stc.MaximumStorageAllocation != nil {
		msa := stc.MaximumStorageAllocation
		data := map[string]interface{}{
			MaxAllocationSize: msa.MaxAllocationSize,
			StoragePool:       msa.StoragePool,
			StorageType:       msa.StorageType,
		}
		d.Set(MaximumStorageAllocation, flex.Flatten(data))
	}
	stcResult := make([]map[string]interface{}, 0, len(stc.StorageTypesCapacity))
	for _, st := range stc.StorageTypesCapacity {
		stResult := map[string]interface{}{}
		if st.MaximumStorageAllocation != nil {
			msa := st.MaximumStorageAllocation
			data := map[string]interface{}{
				MaxAllocationSize: msa.MaxAllocationSize,
				StoragePool:       msa.StoragePool,
				StorageType:       msa.StorageType,
			}
			stResult[MaximumStorageAllocation] = flex.Flatten(data)
		}
		spc := make([]map[string]string, 0, len(st.StoragePoolsCapacity))
		for _, sp := range st.StoragePoolsCapacity {
			data := map[string]interface{}{
				MaxAllocationSize: sp.MaxAllocationSize,
				PoolName:          sp.PoolName,
				StorageType:       sp.StorageType,
				TotalCapacity:     sp.TotalCapacity,
			}
			spc = append(spc, flex.Flatten(data))
		}
		stResult[StoragePoolsCapacity] = spc
		stResult[StorageType] = st.StorageType
		stcResult = append(stcResult, stResult)
	}

	d.Set(StorageTypesCapacity, stcResult)

	return nil
}
