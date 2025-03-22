// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package schematics

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/schematics-go-sdk/schematicsv1"
)

func DataSourceIBMSchematicsState() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMSchematicsStateRead,

		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the workspace for which you want to retrieve the Terraform statefile URL.  To find the workspace ID, use the GET /v1/workspaces API.",
			},
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Region of the workspace.",
			},
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the Terraform template for which you want to retrieve the Terraform statefile.  When you create a workspace, the Terraform template that your workspace points to is assigned a unique ID.  To find this ID, use the GET /v1/workspaces API and review the template_data.id value.",
			},
			"state_store": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state_store_json": {
				Type:     schema.TypeString,
				Computed: true,
			},
			flex.ResourceControllerURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the IBM Cloud dashboard that can be used to explore and view details about this workspace",
			},
		},
	}
}

func dataSourceIBMSchematicsStateRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsStateRead schematicsClient initialization failed: %s", err.Error()), "ibm_schematics_state", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if r, ok := d.GetOk("location"); ok {
		region := r.(string)
		schematicsURL, updatedURL, _ := SchematicsEndpointURL(region, meta)
		if updatedURL {
			schematicsClient.Service.Options.URL = schematicsURL
		}
	}

	getWorkspaceTemplateStateOptions := &schematicsv1.GetWorkspaceTemplateStateOptions{}

	getWorkspaceTemplateStateOptions.SetWID(d.Get("workspace_id").(string))
	getWorkspaceTemplateStateOptions.SetTID(d.Get("template_id").(string))

	_, response, _ := schematicsClient.GetWorkspaceTemplateStateWithContext(context, getWorkspaceTemplateStateOptions)
	if response.StatusCode != 200 {

		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsStateRead GetWorkspaceTemplateStateWithContext failed with error: %s and response:\n%s", err, response), "ibm_schematics_state", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIBMSchematicsStateID(d))

	var stateStore map[string]interface{}
	json.Unmarshal(response.RawResult, &stateStore)

	b := bytes.NewReader(response.RawResult)

	decoder := json.NewDecoder(b)
	decoder.UseNumber()
	decoder.Decode(&stateStore)

	statestr := fmt.Sprintf("%v", stateStore)
	d.Set("state_store", statestr)

	stateByte, err := json.MarshalIndent(stateStore, "", "")
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsStateRead failed: %s", err.Error()), "ibm_schematics_state", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	stateStoreJSON := string(stateByte[:])
	d.Set("state_store_json", stateStoreJSON)

	controller, err := flex.GetBaseController(meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsStateRead failed: %s", err.Error()), "ibm_schematics_state", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	d.Set(flex.ResourceControllerURL, controller+"/schematics")

	return nil
}

// dataSourceIBMSchematicsStateID returns a reasonable ID for the list.
func dataSourceIBMSchematicsStateID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
