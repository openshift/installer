// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPIStoragePoolCapacity() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIStoragePoolCapacityRead,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_StoragePool: {
				Description:  "The storage pool name.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_MaxAllocationSize: {
				Computed:    true,
				Description: "Maximum allocation storage size (GB).",
				Type:        schema.TypeInt,
			},
			Attr_ReplicationEnabled: {
				Computed:    true,
				Description: "Replication status of the storage pool.",
				Type:        schema.TypeBool,
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
	}
}

func dataSourceIBMPIStoragePoolCapacityRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	storagePool := d.Get(Arg_StoragePool).(string)

	client := instance.NewIBMPIStorageCapacityClient(ctx, sess, cloudInstanceID)
	sp, err := client.GetStoragePoolCapacity(storagePool)
	if err != nil {
		log.Printf("[ERROR] get storage pool capacity failed %v", err)
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, storagePool))
	d.Set(Attr_MaxAllocationSize, *sp.MaxAllocationSize)
	d.Set(Attr_ReplicationEnabled, *sp.ReplicationEnabled)
	d.Set(Attr_StorageType, sp.StorageType)
	d.Set(Attr_TotalCapacity, sp.TotalCapacity)
	return nil
}
