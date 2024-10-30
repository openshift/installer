// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package codeengine

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/code-engine-go-sdk/codeenginev2"
)

func DataSourceIbmCodeEngineProject() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmCodeEngineProjectRead,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the project.",
			},
			"account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "An alphanumeric value identifying the account ID.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the project was created.",
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN of the project.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When you provision a new resource, a URL is created identifying the location of the instance.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the project.",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The region for your project deployment. Possible values: 'au-syd', 'br-sao', 'ca-tor', 'eu-de', 'eu-gb', 'jp-osa', 'jp-tok', 'us-east', 'us-south'.",
			},
			"resource_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the resource group.",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the project.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current state of the project. For example, when the project is created and is ready for use, the status of the project is active.",
			},
		},
	}
}

func dataSourceIbmCodeEngineProjectRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_code_engine_project", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getProjectOptions := &codeenginev2.GetProjectOptions{}

	getProjectOptions.SetID(d.Get("project_id").(string))
	log.Printf("[DEBUG]\n%+v", getProjectOptions)

	project, _, err := codeEngineClient.GetProjectWithContext(context, getProjectOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetProjectWithContext failed: %s", err.Error()), "(Data) ibm_code_engine_project", "read")

		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*getProjectOptions.ID)

	if err = d.Set("account_id", project.AccountID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting account_id: %s", err), "(Data) ibm_code_engine_project", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("created_at", project.CreatedAt); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_at: %s", err), "(Data) ibm_code_engine_project", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("crn", project.Crn); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting crn: %s", err), "(Data) ibm_code_engine_project", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("href", project.Href); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting href: %s", err), "(Data) ibm_code_engine_project", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("name", project.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_code_engine_project", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("region", project.Region); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting region: %s", err), "(Data) ibm_code_engine_project", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("resource_group_id", project.ResourceGroupID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting resource_group_id: %s", err), "(Data) ibm_code_engine_project", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("resource_type", project.ResourceType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting resource_type: %s", err), "(Data) ibm_code_engine_project", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("status", project.Status); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting status: %s", err), "(Data) ibm_code_engine_project", "read")
		return tfErr.GetDiag()
	}

	return nil
}
