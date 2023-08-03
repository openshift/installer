// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package catalogmanagement

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/platform-services-go-sdk/catalogmanagementv1"
)

func DataSourceIBMCmPreset() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMCmPresetRead,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the preset.  Format is <catalog_id>-<object_name>@<version>",
			},
			"preset": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The map of preset values as a JSON string.",
			},
		},
	}
}

func dataSourceIBMCmPresetRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	catalogManagementClient, err := meta.(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		return diag.FromErr(err)
	}

	presetID := d.Get("id").(string)
	regex := "[A-Za-z0-9]+-[A-Za-z0-9]+-[A-Za-z0-9]+-[A-Za-z0-9]+-[A-Za-z0-9]+-([A-Za-z0-9]+(_[A-Za-z0-9]+)+)@[A-Za-z0-9]"
	match, err := regexp.MatchString(regex, presetID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error attempting regex match string %s", err))
	}
	if !match {
		return diag.FromErr(fmt.Errorf("Error: Preset ID does not match required format. Must be <catalog_id>-<object_name>@<version> %s", err))
	}
	splitID := strings.Split(presetID, "@")
	version := splitID[len(splitID)-1]
	objectID := splitID[0]
	splitID = strings.Split(presetID, "-")
	catalogID := strings.Join(splitID[:5], "-")

	getObjectOptions := &catalogmanagementv1.GetObjectOptions{}

	getObjectOptions.SetCatalogIdentifier(catalogID)
	getObjectOptions.SetObjectIdentifier(objectID)

	catalogObject, response, err := catalogManagementClient.GetObjectWithContext(context, getObjectOptions)
	if err != nil {
		log.Printf("[DEBUG] GetObjectWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetObjectWithContext failed %s\n%s", err, response))
	}

	d.SetId(presetID)

	if catalogObject.Data == nil {
		return diag.FromErr(fmt.Errorf("Error setting preset, object data is nil. %s", err))
	}
	if catalogObject.Data["versions"] == nil {
		return diag.FromErr(fmt.Errorf("Error setting preset, object data.versions is nil. %s", err))
	}
	if catalogObject.Data["versions"].(map[string]interface{})[version] == nil {
		return diag.FromErr(fmt.Errorf("Error setting preset, could not find preset with version %s. %s", version, err))
	}
	if catalogObject.Data["versions"].(map[string]interface{})[version].(map[string]interface{})["preset"] == nil {
		return diag.FromErr(fmt.Errorf("Error setting preset, preset field not found in version %s. %s", version, err))
	}

	presetMap, err := json.Marshal(catalogObject.Data["versions"].(map[string]interface{})[version].(map[string]interface{})["preset"])
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error setting preset, error with json marshal: %s", err))
	}
	if err = d.Set("preset", string(presetMap)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting preset: %s", err))
	}

	return nil
}
