// Copyright IBM Corp. 2022 All Rights Reserved.
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

func DataSourceIBMPIStorageTypesCapacity() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIStorageTypesCapacityRead,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_MaximumStorageAllocation: {
				Computed:    true,
				Description: "Maximum storage allocation.",
				Type:        schema.TypeMap,
			},
			Attr_StorageTypesCapacity: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of storage types capacity.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_MaximumStorageAllocation: {
							Computed:    true,
							Description: "Maximum storage allocation.",
							Type:        schema.TypeMap,
						},
						Attr_StoragePoolsCapacity: {
							Computed:    true,
							Description: "List of storage types capacity.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									Attr_MaxAllocationSize: {
										Computed:    true,
										Description: "Maximum allocation storage size (GB).",
										Type:        schema.TypeInt,
									},
									Attr_PoolName: {
										Computed:    true,
										Description: "The pool name.",
										Type:        schema.TypeString,
									},
									Attr_StorageType: {
										Computed:    true,
										Description: "Storage type of the storage pool.",
										Type:        schema.TypeString,
									},
									Attr_TotalCapacity: {
										Computed:    true,
										Description: "Total pool capacity (GB).",
										Type:        schema.TypeInt,
									},
								},
							},
							Type: schema.TypeList,
						},
						Attr_StorageType: {
							Computed:    true,
							Description: "The storage type.",
							Type:        schema.TypeString,
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

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)

	client := instance.NewIBMPIStorageCapacityClient(ctx, sess, cloudInstanceID)
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
			Attr_MaxAllocationSize: *msa.MaxAllocationSize,
			Attr_StoragePool:       *msa.StoragePool,
			Attr_StorageType:       *msa.StorageType,
		}
		d.Set(Attr_MaximumStorageAllocation, flex.Flatten(data))
	}
	stcResult := make([]map[string]interface{}, 0, len(stc.StorageTypesCapacity))
	for _, st := range stc.StorageTypesCapacity {
		stResult := map[string]interface{}{}
		if st.MaximumStorageAllocation != nil {
			msa := st.MaximumStorageAllocation
			data := map[string]interface{}{
				Attr_MaxAllocationSize: *msa.MaxAllocationSize,
				Attr_StoragePool:       *msa.StoragePool,
				Attr_StorageType:       *msa.StorageType,
			}
			stResult[Attr_MaximumStorageAllocation] = flex.Flatten(data)
		}
		spc := make([]map[string]interface{}, 0, len(st.StoragePoolsCapacity))
		for _, sp := range st.StoragePoolsCapacity {
			data := map[string]interface{}{
				Attr_MaxAllocationSize: *sp.MaxAllocationSize,
				Attr_PoolName:          sp.PoolName,
				Attr_StorageType:       sp.StorageType,
				Attr_TotalCapacity:     sp.TotalCapacity,
			}
			spc = append(spc, data)
		}
		stResult[Attr_StoragePoolsCapacity] = spc
		stResult[Attr_StorageType] = st.StorageType
		stcResult = append(stcResult, stResult)
	}

	d.Set(Attr_StorageTypesCapacity, stcResult)

	return nil
}
