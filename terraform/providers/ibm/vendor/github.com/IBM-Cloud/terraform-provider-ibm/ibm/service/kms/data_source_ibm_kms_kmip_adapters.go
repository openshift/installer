// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package kms

import (
	"context"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	kp "github.com/IBM/keyprotect-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMKMSKmipAdapters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMKMSKmipAdaptersList,
		Schema: map[string]*schema.Schema{
			"endpoint_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{"public", "private"}),
				Description:  "public or private",
			},
			"instance_id": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Key protect or hpcs instance GUID",
				DiffSuppressFunc: suppressKMSInstanceIDDiff,
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Limit of how many adapters to be fetched",
			},
			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Offset of adapters to be fetched",
			},
			"show_total_count": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Flag to return the count of how many adapters there are in total",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "If show_total_count is true, this will contain the total number of adapters after pagination",
			},
			"adapters": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A collection of KMIP adapters",
				Elem: &schema.Resource{
					Schema: dataSourceIBMKMSKmipAdapterBaseSchema(),
				},
			},
		},
	}
}

func dataSourceIBMKMSKmipAdaptersList(d *schema.ResourceData, meta interface{}) error {
	// initialize API
	instanceID := getInstanceIDFromResourceData(d, "instance_id")
	api, _, err := populateKPClient(d, meta, instanceID)
	if err != nil {
		return err
	}

	// call GetKMIPAdapters api
	opts := &kp.ListKmipAdaptersOptions{}
	if limit, ok := d.GetOk("limit"); ok {
		limitVal := uint32(limit.(int))
		opts.Limit = &limitVal
	}
	if offset, ok := d.GetOk("offset"); ok {
		offsetVal := uint32(offset.(int))
		opts.Offset = &offsetVal
	}
	showTotalCountEnabled := false
	if showTotalCount, ok := d.GetOk("show_total_count"); ok {
		showTotalCountEnabled = showTotalCount.(bool)
		opts.TotalCount = &showTotalCountEnabled
	}

	adapters, err := api.GetKMIPAdapters(context.Background(), opts)
	if err != nil {
		return flex.FmtErrorf("[ERROR] Error listing KMIP adapters: %s", err)
	}

	adaptersList := adapters.Adapters

	// set computed values
	mySlice := make([]map[string]interface{}, 0, len(adaptersList))
	for _, adapter := range adaptersList {
		adapterMap := dataSourceIBMKMSKmipAdapterToMap(adapter)
		mySlice = append(mySlice, adapterMap)
	}
	d.Set("adapters", mySlice)
	d.Set("instance_id", instanceID)
	d.SetId(instanceID)
	if showTotalCountEnabled {
		d.Set("total_count", adapters.Metadata.TotalCount)
	}
	return nil
}

func dataSourceIBMKMSKmipAdapterToMap(model kp.KMIPAdapter) map[string]interface{} {
	modelMap := make(map[string]interface{})
	modelMap["adapter_id"] = model.ID
	modelMap["name"] = model.Name
	modelMap["description"] = model.Description
	modelMap["profile"] = model.Profile
	modelMap["profile_data"] = model.ProfileData
	modelMap["created_at"] = model.CreatedAt.String()
	modelMap["created_by"] = model.CreatedBy
	modelMap["updated_at"] = model.UpdatedAt.String()
	modelMap["updated_by"] = model.CreatedBy
	return modelMap
}
