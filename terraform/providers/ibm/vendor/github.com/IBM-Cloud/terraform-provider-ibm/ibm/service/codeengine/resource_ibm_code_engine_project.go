// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package codeengine

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/code-engine-go-sdk/codeenginev2"
	"github.com/IBM/go-sdk-core/v5/core"
)

func ResourceIbmCodeEngineProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmCodeEngineProjectCreate,
		ReadContext:   resourceIbmCodeEngineProjectRead,
		DeleteContext: resourceIbmCodeEngineProjectDelete,
		Importer:      &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_project", "name"),
				Description:  "The name of the project.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the project.",
			},
			"resource_group_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_project", "resource_group_id"),
				Description:  "The ID of the resource group.",
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
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The region for your project deployment. Possible values: 'au-syd', 'br-sao', 'ca-tor', 'eu-de', 'eu-gb', 'jp-osa', 'jp-tok', 'us-east', 'us-south'.",
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

func ResourceIbmCodeEngineProjectValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^([^\x00-\x7F]|[a-zA-Z0-9\-._: ])+$`,
			MinValueLength:             1,
			MaxValueLength:             128,
		},
		validate.ValidateSchema{
			Identifier:                 "resource_group_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[a-z0-9]*$`,
			MinValueLength:             32,
			MaxValueLength:             32,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_code_engine_project", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmCodeEngineProjectCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_code_engine_project", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createProjectOptions := &codeenginev2.CreateProjectOptions{}

	createProjectOptions.SetName(d.Get("name").(string))
	if _, ok := d.GetOk("resource_group_id"); ok {
		createProjectOptions.SetResourceGroupID(d.Get("resource_group_id").(string))
	}
	if _, ok := d.GetOk("tags"); ok {
		var tags []string
		for _, v := range d.Get("tags").([]interface{}) {
			tagsItem := v.(string)
			tags = append(tags, tagsItem)
		}
		createProjectOptions.SetTags(tags)
	}

	project, _, err := codeEngineClient.CreateProjectWithContext(context, createProjectOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateProjectWithContext failed: %s", err.Error()), "ibm_code_engine_project", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*project.ID)

	_, err = waitForIbmCodeEngineProjectCreate(d, meta)
	if err != nil {
		errMsg := fmt.Sprintf("Error waiting for resource IbmCodeEngineProject (%s) to be created: %s", d.Id(), err)
		tfErr := flex.TerraformErrorf(err, errMsg, "ibm_code_engine_project", "create")
		return tfErr.GetDiag()
	}

	return resourceIbmCodeEngineProjectRead(context, d, meta)
}

func waitForIbmCodeEngineProjectCreate(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		return false, err
	}
	getProjectOptions := &codeenginev2.GetProjectOptions{}

	getProjectOptions.SetID(d.Id())

	stateConf := &resource.StateChangeConf{
		Pending: []string{"creating", "preparing"},
		Target:  []string{"active"},
		Refresh: func() (interface{}, string, error) {
			stateObj, response, err := codeEngineClient.GetProject(getProjectOptions)
			if err != nil {
				if sdkErr, ok := err.(*core.SDKProblem); ok && response.GetStatusCode() == 404 {
					sdkErr.Summary = fmt.Sprintf("The instance %s does not exist anymore: %s", "getProjectOptions", err)
					return nil, "", sdkErr
				}
				return nil, "", err
			}
			failStates := map[string]bool{"creation_failed": true}
			if failStates[*stateObj.Status] {
				return stateObj, *stateObj.Status, fmt.Errorf("the instance %s failed: %s", "getProjectOptions", err)
			}
			return stateObj, *stateObj.Status, nil
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      20 * time.Second,
		MinTimeout: 20 * time.Second,
	}

	return stateConf.WaitForState()
}

func resourceIbmCodeEngineProjectRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_code_engine_project", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getProjectOptions := &codeenginev2.GetProjectOptions{}

	getProjectOptions.SetID(d.Id())

	project, response, err := codeEngineClient.GetProjectWithContext(context, getProjectOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetProjectWithContext failed: %s", err.Error()), "ibm_code_engine_project", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("name", project.Name); err != nil {
		return diag.FromErr(fmt.Errorf("error setting name: %s", err))
	}
	if err = d.Set("project_id", project.ID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting project_id: %s", err))
	}
	if !core.IsNil(project.ResourceGroupID) {
		if err = d.Set("resource_group_id", project.ResourceGroupID); err != nil {
			return diag.FromErr(fmt.Errorf("error setting resource_group_id: %s", err))
		}
	}
	if !core.IsNil(project.AccountID) {
		if err = d.Set("account_id", project.AccountID); err != nil {
			return diag.FromErr(fmt.Errorf("error setting account_id: %s", err))
		}
	}
	if !core.IsNil(project.CreatedAt) {
		if err = d.Set("created_at", project.CreatedAt); err != nil {
			return diag.FromErr(fmt.Errorf("error setting created_at: %s", err))
		}
	}
	if !core.IsNil(project.Crn) {
		if err = d.Set("crn", project.Crn); err != nil {
			return diag.FromErr(fmt.Errorf("error setting crn: %s", err))
		}
	}
	if !core.IsNil(project.Href) {
		if err = d.Set("href", project.Href); err != nil {
			return diag.FromErr(fmt.Errorf("error setting href: %s", err))
		}
	}
	if !core.IsNil(project.Region) {
		if err = d.Set("region", project.Region); err != nil {
			return diag.FromErr(fmt.Errorf("error setting region: %s", err))
		}
	}
	if !core.IsNil(project.ResourceType) {
		if err = d.Set("resource_type", project.ResourceType); err != nil {
			return diag.FromErr(fmt.Errorf("error setting resource_type: %s", err))
		}
	}
	if !core.IsNil(project.Status) {
		if err = d.Set("status", project.Status); err != nil {
			return diag.FromErr(fmt.Errorf("error setting status: %s", err))
		}
	}

	return nil
}

func resourceIbmCodeEngineProjectDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_code_engine_project", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteProjectOptions := &codeenginev2.DeleteProjectOptions{}

	deleteProjectOptions.SetID(d.Id())

	_, err = codeEngineClient.DeleteProjectWithContext(context, deleteProjectOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteProjectWithContext failed: %s", err.Error()), "ibm_code_engine_project", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}
