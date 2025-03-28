// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package kms

import (
	"context"
	"fmt"
	"strings"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	kp "github.com/IBM/keyprotect-go-client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	resourceIBMKmsKMIPAdapterValidProfiles = []string{"native_1.0"}
)

func ResourceIBMKmsKMIPAdapter() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMKmsKMIPAdapterCreate,
		Read:     resourceIBMKmsKMIPAdapterRead,
		Delete:   resourceIBMKmsKMIPAdapterDelete,
		Exists:   resourceIBMKmsKMIPAdapterExists,
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
			"profile": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The profile of the KMIP adapter",
				ValidateFunc: validate.ValidateAllowedStringValues(resourceIBMKmsKMIPAdapterValidProfiles),
			},
			"profile_data": {
				Type:        schema.TypeMap,
				Required:    true,
				ForceNew:    true,
				Description: "The data specific to the KMIP Adapter profile",
			},
			"adapter_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the KMIP adapter",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The name of the KMIP adapter",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The description of the KMIP adapter",
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
		},
	}
}

func resourceIBMKmsKMIPAdapterProfileToProfileFunc(profile string, profileData map[string]string) kp.CreateKMIPAdapterProfile {
	if profile == resourceIBMKmsKMIPAdapterValidProfiles[0] {
		//native_1.0
		return kp.WithNativeProfile(profileData["crk_id"])
	}
	// Shouldn't reach here, since we check for profile validity before this.
	return nil
}

func resourceIBMKmsKMIPAdapterCreate(d *schema.ResourceData, meta interface{}) error {
	adapterToCreate, instanceID, err := ExtractAndValidateKMIPAdapterDataFromSchema(d)
	if err != nil {
		return err
	}
	kpAPI, _, err := populateKPClient(d, meta, instanceID)
	if err != nil {
		return err
	}
	if strings.Contains((kpAPI.URL).String(), "private") || strings.Contains(kpAPI.Config.BaseURL, "private") {
		d.Set("endpoint_type", "private")
	} else {
		d.Set("endpoint_type", "public")
	}
	adapter, err := kpAPI.CreateKMIPAdapter(context.Background(),
		resourceIBMKmsKMIPAdapterProfileToProfileFunc(adapterToCreate.Profile, adapterToCreate.ProfileData),
		kp.WithKMIPAdapterName(adapterToCreate.Name),
		kp.WithKMIPAdapterDescription(adapterToCreate.Description),
	)
	if err != nil {
		return flex.FmtErrorf("[ERROR] Error while creating KMIP adapter: %s", err)
	}
	return populateKMIPAdapterSchemaDataFromStruct(d, *adapter, instanceID)
}

func resourceIBMKmsKMIPAdapterRead(d *schema.ResourceData, meta interface{}) error {
	instanceID, adapterID, err := splitAdapterID(d.Id())
	if err != nil {
		return err
	}
	kpAPI, _, err := populateKPClient(d, meta, instanceID)
	if err != nil {
		return err
	}
	if strings.Contains((kpAPI.URL).String(), "private") || strings.Contains(kpAPI.Config.BaseURL, "private") {
		d.Set("endpoint_type", "private")
	} else {
		d.Set("endpoint_type", "public")
	}
	ctx := context.Background()
	adapter, err := kpAPI.GetKMIPAdapter(ctx, adapterID)
	if err != nil {
		return err
	}
	return populateKMIPAdapterSchemaDataFromStruct(d, *adapter, instanceID)
}

func resourceIBMKmsKMIPAdapterDelete(d *schema.ResourceData, meta interface{}) error {
	instanceID := d.Get("instance_id").(string)
	_, adapterID, err := splitAdapterID(d.Id())
	if err != nil {
		return err
	}
	kpAPI, _, err := populateKPClient(d, meta, instanceID)
	if err != nil {
		return err
	}
	ctx := context.Background()
	objects, err := kpAPI.GetKMIPObjects(ctx, adapterID, nil)
	if err != nil {
		return flex.FmtErrorf("[ERROR] Failed to fetch KMIP objects associated with adapter '%s' for deletion: %v", adapterID, err)
	}

	for _, object := range objects.Objects {
		err = kpAPI.DeleteKMIPObject(ctx, adapterID, object.ID, kp.WithForce(true))
		if err != nil {
			if kpError, ok := err.(*kp.Error); ok {
				if kpError.StatusCode == 404 || kpError.StatusCode == 410 {
					// if the kmip object is already deleted, do not error out
					continue
				}
			}
			return flex.FmtErrorf("[ERROR] Failed to delete KMIP object associated with adapter (%s): %s",
				adapterID,
				err,
			)
		}
	}

	err = kpAPI.DeleteKMIPAdapter(ctx, adapterID)
	return err
}

func resourceIBMKmsKMIPAdapterExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	instanceID, adapterID, err := splitAdapterID(d.Id())
	if err != nil {
		return false, err
	}
	kpAPI, _, err := populateKPClient(d, meta, instanceID)
	if err != nil {
		return false, err
	}
	ctx := context.Background()
	_, err = kpAPI.GetKMIPAdapter(ctx, adapterID)
	if err != nil {
		if kpError, ok := err.(*kp.Error); ok {
			if kpError.StatusCode == 404 {
				return false, nil
			}
		}
		return false, wrapError(err, "Error checking adapter existence")
	}
	return true, nil
}

func ExtractAndValidateKMIPAdapterDataFromSchema(d *schema.ResourceData) (adapter kp.KMIPAdapter, instanceID string, err error) {
	err = nil
	instanceID = getInstanceIDFromResourceData(d, "instance_id")
	profile, ok := d.Get("profile").(string)
	if !ok {
		err = flex.FmtErrorf("[ERROR] Error converting profile to string")
		return
	}
	adapter = kp.KMIPAdapter{
		Profile: profile,
	}
	if name, ok := d.GetOk("name"); ok {
		nameStr, ok2 := name.(string)
		if !ok2 {
			err = flex.FmtErrorf("[ERROR] Error converting name to string")
			return
		}
		adapter.Name = nameStr
	}
	if desc, ok := d.GetOk("description"); ok {
		descStr, ok2 := desc.(string)
		if !ok2 {
			err = flex.FmtErrorf("[ERROR] Error converting description to string")
			return
		}
		adapter.Description = descStr
	}
	if data, ok := d.GetOk("profile_data"); ok {
		dataMap, ok2 := data.(map[string]interface{})
		if !ok2 {
			err = flex.FmtErrorf("[ERROR] Error converting profile data to map[string]interface{}")
			return
		}
		profileData := map[string]string{}
		for key := range dataMap {
			if val, ok := dataMap[key].(string); ok {
				profileData[key] = val
			} else {
				err = flex.FmtErrorf("[ERROR] Error converting value with key {%s} into string", key)
				return
			}
		}
		adapter.ProfileData = profileData
	}
	return
}

func populateKMIPAdapterSchemaDataFromStruct(d *schema.ResourceData, adapter kp.KMIPAdapter, instanceID string) (err error) {
	d.SetId(fmt.Sprintf("%s/%s", instanceID, adapter.ID))

	if err = d.Set("name", adapter.Name); err != nil {
		return flex.FmtErrorf("[ERROR] Error setting name: %s", err)
	}
	if err = d.Set("adapter_id", adapter.ID); err != nil {
		return flex.FmtErrorf("[ERROR] Error setting adapter_id: %s", err)
	}
	if err = d.Set("instance_id", instanceID); err != nil {
		return flex.FmtErrorf("[ERROR] Error setting instance_id: %s", err)
	}
	if err = d.Set("description", adapter.Description); err != nil {
		return flex.FmtErrorf("[ERROR] Error setting description: %s", err)
	}
	if err = d.Set("profile", adapter.Profile); err != nil {
		return flex.FmtErrorf("[ERROR] Error setting profile: %s", err)
	}
	if err = d.Set("profile_data", adapter.ProfileData); err != nil {
		return flex.FmtErrorf("[ERROR] Error setting profile_data: %s", err)
	}
	if err = d.Set("created_at", adapter.CreatedAt.String()); err != nil {
		return flex.FmtErrorf("[ERROR] Error setting created_at: %s", err)
	}
	if err = d.Set("created_by", adapter.CreatedBy); err != nil {
		return flex.FmtErrorf("[ERROR] Error setting created_by: %s", err)
	}
	if err = d.Set("updated_at", adapter.UpdatedAt.String()); err != nil {
		return flex.FmtErrorf("[ERROR] Error setting updated_at: %s", err)
	}
	if err = d.Set("updated_by", adapter.UpdatedBy); err != nil {
		return flex.FmtErrorf("[ERROR] Error setting updated_by: %s", err)
	}
	return nil
}

func splitAdapterID(terraformId string) (instanceID, adapterID string, err error) {
	split, err := flex.SepIdParts(terraformId, "/")
	if err != nil {
		return "", "", err
	}
	if len(split) != 2 {
		return "", "", flex.FmtErrorf("[ERROR] The given id %s does not contain all expected sections, should be of format instance_id/adapter_id", terraformId)
	}
	for index, id := range split {
		if uuid.Validate(id) != nil {
			return "", "", flex.FmtErrorf("[ERROR] The given id %s at index %d of instance_id/adapter_id is not a valid UUID", id, index)
		}
	}
	return split[0], split[1], nil
}
