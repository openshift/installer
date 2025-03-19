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

func DataSourceIBMKMSKMIPObjects() *schema.Resource {
	return &schema.Resource{
		Read:     dataSourceIBMKmsKMIPObjectList,
		Importer: &schema.ResourceImporter{},
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
				Description:      "Key protect Instance GUID",
				ForceNew:         true,
				DiffSuppressFunc: suppressKMSInstanceIDDiff,
			},
			"adapter_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The id of the KMIP adapter that contains the cert",
				ForceNew:     true,
				ExactlyOneOf: []string{"adapter_id", "adapter_name"},
			},
			"adapter_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The name of the KMIP adapter that contains the cert",
				ForceNew:     true,
				ExactlyOneOf: []string{"adapter_id", "adapter_name"},
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Limit of how many objects to be fetched",
			},
			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Offset of objects to be fetched",
			},
			"show_total_count": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Flag to return the count of how many objects there are in total after the filter",
			},
			"object_state_filter": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of integers representing Object States to filter for",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "If show_total_count is true, this will contain the total number of KMIP objects after pagination.",
			},
			"objects": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The objects contained in the specified adapter",
				Elem: &schema.Resource{
					Schema: dataSourceIBMKMSKMIPObjectBaseSchema(true),
				},
			},
		},
	}
}

func dataSourceIBMKmsKMIPObjectList(d *schema.ResourceData, meta interface{}) error {
	// initialize API
	instanceID := getInstanceIDFromResourceData(d, "instance_id")
	kpAPI, _, err := populateKPClient(d, meta, instanceID)
	if err != nil {
		return err
	}
	// get adapterID
	nameOrID, hasID := d.GetOk("adapter_id")
	if !hasID {
		nameOrID = d.Get("adapter_name")
	}
	adapterNameOrID := nameOrID.(string)

	opts := &kp.ListKmipObjectsOptions{}
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
	if stateFilter, ok := d.GetOk("object_state_filter"); ok {
		arrayVal, ok2 := stateFilter.([]any)
		if !ok2 {
			return flex.FmtErrorf("[ERROR] Error converting object_state_filter into []any")
		}
		int32Arr := make([]int32, 0, len(arrayVal))
		for _, myint := range arrayVal {
			int32Arr = append(int32Arr, int32(myint.(int)))
		}
		opts.ObjectStateFilter = &int32Arr
	}

	ctx := context.Background()
	adapter, err := kpAPI.GetKMIPAdapter(ctx, adapterNameOrID)
	if err != nil {
		return flex.FmtErrorf("[ERROR] Error while retriving KMIP adapter to list KMIP objects: %s", err)
	}
	if err = d.Set("adapter_id", adapter.ID); err != nil {
		return flex.FmtErrorf("[ERROR] Error setting adapter_id: %s", err)
	}
	if err = d.Set("adapter_name", adapter.Name); err != nil {
		return flex.FmtErrorf("[ERROR] Error setting adapter_name: %s", err)
	}
	objs, err := kpAPI.GetKMIPObjects(ctx, adapterNameOrID, opts)
	if err != nil {
		return flex.FmtErrorf("[ERROR] Error while retriving KMIP objects associated with adapter ID '%s': %v", adapter.ID, err)
	}
	objsList := objs.Objects
	// set computed values
	mySlice := make([]map[string]interface{}, 0, len(objsList))
	for _, obj := range objsList {
		objMap := dataSourceIBMKMSKmipObjectToMap(obj)
		mySlice = append(mySlice, objMap)
	}
	d.Set("objects", mySlice)
	d.SetId(adapter.ID)

	if showTotalCountEnabled {
		d.Set("total_count", objs.Metadata.TotalCount)
	}

	return nil
}

func dataSourceIBMKMSKmipObjectToMap(model kp.KMIPObject) map[string]interface{} {
	modelMap := make(map[string]interface{})
	modelMap["object_id"] = model.ID
	modelMap["object_type"] = model.KMIPObjectType
	modelMap["object_state"] = model.ObjectState
	if model.CreatedAt != nil {
		modelMap["created_at"] = model.CreatedAt.String()
		modelMap["created_by"] = model.CreatedBy
		modelMap["created_by_cert_id"] = model.CreatedByCertID
	}
	if model.UpdatedAt != nil {
		modelMap["updated_at"] = model.UpdatedAt.String()
		modelMap["updated_by"] = model.UpdatedBy
		modelMap["updated_by_cert_id"] = model.UpdatedByCertID
	}
	if model.DestroyedAt != nil {
		modelMap["destroyed_at"] = model.DestroyedAt.String()
		modelMap["destroyed_by"] = model.DestroyedBy
		modelMap["destroyed_by_cert_id"] = model.DestroyedByCertID
	}

	return modelMap
}
