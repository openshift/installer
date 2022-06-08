// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"fmt"

	"log"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"

	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	PIPoolName = "pi_storage_pool"
)

func DataSourceIBMPIStoragePoolCapacity() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIStoragePoolCapacityRead,
		Schema: map[string]*schema.Schema{
			helpers.PICloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			PIPoolName: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "Storage pool name",
			},
			// Computed Attributes
			MaxAllocationSize: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Maximum allocation storage size (GB)",
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
	}
}

func dataSourceIBMPIStoragePoolCapacityRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)
	storagePool := d.Get(PIPoolName).(string)

	client := st.NewIBMPIStorageCapacityClient(ctx, sess, cloudInstanceID)
	sp, err := client.GetStoragePoolCapacity(storagePool)
	if err != nil {
		log.Printf("[ERROR] get storage pool capacity failed %v", err)
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, storagePool))

	d.Set(MaxAllocationSize, sp.MaxAllocationSize)
	d.Set(StorageType, sp.StorageType)
	d.Set(TotalCapacity, sp.TotalCapacity)

	return nil
}
