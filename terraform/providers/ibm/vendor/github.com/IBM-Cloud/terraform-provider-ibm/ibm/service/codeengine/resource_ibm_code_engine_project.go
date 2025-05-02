// Copyright IBM Corp. 2023 All Rights Reserved.
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

	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
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
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_project", "name"),
				Description:  "The name of the project.",
			},
			"project_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the project.",
			},
			"resource_group_id": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_project", "resource_group_id"),
				Description:  "Optional ID of the resource group for your project deployment. If this field is not defined, the default resource group of the account will be used.",
			},
			"account_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "An alphanumeric value identifying the account ID.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the project was created.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN of the project.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When you provision a new resource, a URL is created identifying the location of the instance.",
			},
			"region": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The region for your project deployment. Possible values: 'au-syd', 'br-sao', 'ca-tor', 'eu-de', 'eu-gb', 'jp-osa', 'jp-tok', 'us-east', 'us-south'.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the project.",
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current state of the project. For example, if the project is created and ready to get used, it will return active.",
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
		return diag.FromErr(err)
	}

	createProjectOptions := &codeenginev2.CreateProjectOptions{}

	createProjectOptions.SetName(d.Get("name").(string))
	if _, ok := d.GetOk("resource_group_id"); ok {
		createProjectOptions.SetResourceGroupID(d.Get("resource_group_id").(string))
	}

	project, response, err := codeEngineClient.CreateProjectWithContext(context, createProjectOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateProjectWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateProjectWithContext failed %s\n%s", err, response))
	}

	d.SetId(*project.ID)

	_, err = waitForIbmCodeEngineProjectCreate(d, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf(
			"Error waiting for resource IbmCodeEngineProject (%s) to be created: %s", d.Id(), err))
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
				if apiErr, ok := err.(bmxerror.RequestFailure); ok && apiErr.StatusCode() == 404 {
					return nil, "", fmt.Errorf("The instance %s does not exist anymore: %s\n%s", "getProjectOptions", err, response)
				}
				return nil, "", err
			}
			failStates := map[string]bool{"creation_failed": true}
			if failStates[*stateObj.Status] {
				return stateObj, *stateObj.Status, fmt.Errorf("The instance %s failed: %s\n%s", "getProjectOptions", err, response)
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
		return diag.FromErr(err)
	}

	getProjectOptions := &codeenginev2.GetProjectOptions{}

	getProjectOptions.SetID(d.Id())

	project, response, err := codeEngineClient.GetProjectWithContext(context, getProjectOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetProjectWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetProjectWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("name", project.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("project_id", project.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting project_id: %s", err))
	}
	if !core.IsNil(project.ResourceGroupID) {
		if err = d.Set("resource_group_id", project.ResourceGroupID); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting resource_group_id: %s", err))
		}
	}
	if !core.IsNil(project.AccountID) {
		if err = d.Set("account_id", project.AccountID); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting account_id: %s", err))
		}
	}
	if !core.IsNil(project.CreatedAt) {
		if err = d.Set("created_at", project.CreatedAt); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
		}
	}
	if !core.IsNil(project.Crn) {
		if err = d.Set("crn", project.Crn); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
		}
	}
	if !core.IsNil(project.Href) {
		if err = d.Set("href", project.Href); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
		}
	}
	if !core.IsNil(project.Region) {
		if err = d.Set("region", project.Region); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting region: %s", err))
		}
	}
	if !core.IsNil(project.ResourceType) {
		if err = d.Set("resource_type", project.ResourceType); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
		}
	}
	if !core.IsNil(project.Status) {
		if err = d.Set("status", project.Status); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting status: %s", err))
		}
	}

	return nil
}

func resourceIbmCodeEngineProjectDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteProjectOptions := &codeenginev2.DeleteProjectOptions{}

	deleteProjectOptions.SetID(d.Id())

	response, err := codeEngineClient.DeleteProjectWithContext(context, deleteProjectOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteProjectWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteProjectWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}
