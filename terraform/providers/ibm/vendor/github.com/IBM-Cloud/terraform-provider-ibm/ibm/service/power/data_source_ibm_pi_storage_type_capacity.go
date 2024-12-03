// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPIStorageTypeCapacity() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIStorageTypeCapacityRead,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_StorageType: {
				Description:  "The storage type name.",
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
			Attr_StoragePoolsCapacity: {
				Computed:    true,
				Description: "List of storage pools capacity.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_MaxAllocationSize: {
							Computed:    true,
							Description: "Maximum allocation storage size (GB).",
							Type:        schema.TypeInt,
						},
						Attr_PoolName: {
							Computed:    true,
							Description: "The pool name",
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
		},
	}
}

func dataSourceIBMPIStorageTypeCapacityRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	storageType := d.Get(Arg_StorageType).(string)

	client := instance.NewIBMPIStorageCapacityClient(ctx, sess, cloudInstanceID)
	stc, err := client.GetStorageTypeCapacity(storageType)
	if err != nil {
		log.Printf("[ERROR] get storage type capacity failed %v", err)
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, storageType))

	if stc.MaximumStorageAllocation != nil {
		msa := stc.MaximumStorageAllocation
		data := map[string]interface{}{
			Attr_MaxAllocationSize: *msa.MaxAllocationSize,
			Attr_StoragePool:       *msa.StoragePool,
			Attr_StorageType:       *msa.StorageType,
		}
		d.Set(Attr_MaximumStorageAllocation, flex.Flatten(data))
	}

	result := make([]map[string]interface{}, 0, len(stc.StoragePoolsCapacity))
	for _, sp := range stc.StoragePoolsCapacity {
		data := map[string]interface{}{
			Attr_MaxAllocationSize: *sp.MaxAllocationSize,
			Attr_PoolName:          sp.PoolName,
			Attr_StorageType:       sp.StorageType,
			Attr_TotalCapacity:     sp.TotalCapacity,
		}
		result = append(result, data)
	}
	d.Set(Attr_StoragePoolsCapacity, result)

	return nil
}
