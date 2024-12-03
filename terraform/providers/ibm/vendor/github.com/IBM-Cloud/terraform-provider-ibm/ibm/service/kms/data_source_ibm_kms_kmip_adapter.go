// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package kms

import (
	"context"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIBMKMSKmipAdapterBaseSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"adapter_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The UUID of the KMIP adapter to be fetched",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The name of the KMIP adapter to be fetched",
		},
		"description": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The description of the KMIP adapter to be fetched",
		},
		"profile": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The profile of the KMIP adapter to be fetched",
		},
		"profile_data": {
			Type:        schema.TypeMap,
			Computed:    true,
			Description: "The data specific to the KMIP Adapter profile",
		},
		"created_by": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The unique identifier that is associated with the entity that created the adapter.",
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The date when a resource was created. The date format follows RFC 3339.",
		},
		"updated_by": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The unique identifier that is associated with the entity that updated the adapter.",
		},
		"updated_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The date when a resource was updated. The date format follows RFC 3339.",
		},
	}
}

func DataSourceIBMKMSKmipAdapter() *schema.Resource {
	baseMap := dataSourceIBMKMSKmipAdapterBaseSchema()

	baseMap["endpoint_type"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Computed:     true,
		ValidateFunc: validate.ValidateAllowedStringValues([]string{"public", "private"}),
		Description:  "public or private",
	}

	baseMap["instance_id"] = &schema.Schema{
		Type:             schema.TypeString,
		Required:         true,
		Description:      "Key protect or hpcs instance GUID",
		DiffSuppressFunc: suppressKMSInstanceIDDiff,
	}
	adapterIDSchema := baseMap["adapter_id"]
	adapterIDSchema.Optional = true
	adapterIDSchema.ExactlyOneOf = []string{"adapter_id", "name"}
	baseMap["adapter_id"] = adapterIDSchema

	adapterNameSchema := baseMap["name"]
	adapterNameSchema.Optional = true
	adapterNameSchema.ExactlyOneOf = []string{"adapter_id", "name"}
	baseMap["name"] = adapterNameSchema

	return &schema.Resource{
		Read:   dataSourceIBMKMSKmipAdapterRead,
		Schema: baseMap,
	}
}

func dataSourceIBMKMSKmipAdapterRead(d *schema.ResourceData, meta interface{}) error {
	// initialize API
	instanceID := getInstanceIDFromResourceData(d, "instance_id")
	kpAPI, _, err := populateKPClient(d, meta, instanceID)
	if err != nil {
		return err
	}

	nameOrID, hasID := d.GetOk("adapter_id")
	if !hasID {
		nameOrID = d.Get("name")
	}
	adapterNameOrID := nameOrID.(string)
	// call GetKMIPAdapter api
	adapter, err := kpAPI.GetKMIPAdapter(context.Background(), adapterNameOrID)
	if err != nil {
		return err
	}

	// set computed values
	return populateKMIPAdapterSchemaDataFromStruct(d, *adapter, instanceID)
}
