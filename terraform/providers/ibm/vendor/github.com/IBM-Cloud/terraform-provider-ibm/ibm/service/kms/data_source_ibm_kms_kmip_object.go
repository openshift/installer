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

func dataSourceIBMKMSKMIPObjectBaseSchema(isForList bool) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"object_id": {
			Type:        schema.TypeString,
			Required:    !isForList,
			Computed:    isForList,
			Description: "The id of the KMIP Object to be fetched",
		},
		"object_type": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "The type of KMIP object",
		},
		"object_state": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "The state of the KMIP object",
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
		"created_by_cert_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The ID of the certificate that created the object",
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
		"updated_by_cert_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The ID of the certificate that updated the object",
		},
		"destroyed_by": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The unique identifier that is associated with the entity that destroyed the adapter.",
		},
		"destroyed_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The date when a resource was destroyed. The date format follows RFC 3339.",
		},
		"destroyed_by_cert_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The ID of the certificate that destroyed the object",
		},
	}
}

func DataSourceIBMKMSKMIPObject() *schema.Resource {
	baseMap := dataSourceIBMKMSKMIPObjectBaseSchema(false)

	baseMap["instance_id"] = &schema.Schema{
		Type:             schema.TypeString,
		Required:         true,
		Description:      "Key protect Instance GUID",
		ForceNew:         true,
		DiffSuppressFunc: suppressKMSInstanceIDDiff,
	}
	baseMap["adapter_id"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Computed:     true,
		Description:  "The id of the KMIP adapter that contains the kmip object",
		ForceNew:     true,
		ExactlyOneOf: []string{"adapter_id", "adapter_name"},
	}
	baseMap["adapter_name"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Computed:     true,
		Description:  "The name of the KMIP adapter that contains the kmip object",
		ForceNew:     true,
		ExactlyOneOf: []string{"adapter_id", "adapter_name"},
	}

	baseMap["endpoint_type"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Computed:     true,
		ValidateFunc: validate.ValidateAllowedStringValues([]string{"public", "private"}),
		Description:  "public or private",
	}

	return &schema.Resource{
		Read:     dataSourceIBMKmsKMIPObjectRead,
		Importer: &schema.ResourceImporter{},
		Schema:   baseMap,
	}
}

func dataSourceIBMKmsKMIPObjectRead(d *schema.ResourceData, meta interface{}) error {
	// initialize API
	instanceID := getInstanceIDFromResourceData(d, "instance_id")
	kpAPI, _, err := populateKPClient(d, meta, instanceID)
	if err != nil {
		return err
	}
	// get adapterID and certID
	nameOrID, hasID := d.GetOk("adapter_id")
	if !hasID {
		nameOrID = d.Get("adapter_name")
	}
	adapterNameOrID := nameOrID.(string)

	objectID := d.Get("object_id").(string)

	ctx := context.Background()
	adapter, err := kpAPI.GetKMIPAdapter(ctx, adapterNameOrID)
	if err != nil {
		return flex.FmtErrorf("[ERROR] Error while retriving KMIP adapter to get KMIP object: %s", err)
	}
	if err = d.Set("adapter_id", adapter.ID); err != nil {
		return flex.FmtErrorf("[ERROR] Error setting adapter_id: %s", err)
	}
	if err = d.Set("adapter_name", adapter.Name); err != nil {
		return flex.FmtErrorf("[ERROR] Error setting adapter_name: %s", err)
	}

	object, err := kpAPI.GetKMIPObject(ctx, adapterNameOrID, objectID)
	if err != nil {
		return err
	}
	err = populateKMIPObjectSchemaDataFromStruct(d, *object)
	if err != nil {
		return err
	}
	return nil
}

func populateKMIPObjectSchemaDataFromStruct(d *schema.ResourceData, object kp.KMIPObject) (err error) {
	if err = d.Set("object_id", object.ID); err != nil {
		return flex.FmtErrorf("[ERROR] Error setting object_id: %s", err)
	}
	if err = d.Set("object_type", object.KMIPObjectType); err != nil {
		return flex.FmtErrorf("[ERROR] Error setting object_type: %s", err)
	}
	if err = d.Set("object_state", object.ObjectState); err != nil {
		return flex.FmtErrorf("[ERROR] Error setting object_state: %s", err)
	}
	if object.CreatedAt != nil {
		if err = d.Set("created_at", object.CreatedAt.String()); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting created_at: %s", err)
		}
		if err = d.Set("created_by", object.CreatedBy); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting created_by: %s", err)
		}
		if err = d.Set("created_by_cert_id", object.CreatedByCertID); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting created_by_cert_id: %s", err)
		}
	}
	if object.UpdatedAt != nil {
		if err = d.Set("updated_at", object.UpdatedAt.String()); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting updated_at: %s", err)
		}
		if err = d.Set("updated_by", object.UpdatedBy); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting created_by: %s", err)
		}
		if err = d.Set("updated_by_cert_id", object.UpdatedByCertID); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting updated_by_cert_id: %s", err)
		}
	}
	if object.DestroyedAt != nil {
		if err = d.Set("destroyed_at", object.DestroyedAt.String()); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting destroyed_at: %s", err)
		}
		if err = d.Set("destroyed_by", object.DestroyedBy); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting destroyed_by: %s", err)
		}
		if err = d.Set("destroyed_by_cert_id", object.DestroyedByCertID); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting destroyed_by_cert_id: %s", err)
		}
	}
	d.SetId(object.ID)
	return nil
}
