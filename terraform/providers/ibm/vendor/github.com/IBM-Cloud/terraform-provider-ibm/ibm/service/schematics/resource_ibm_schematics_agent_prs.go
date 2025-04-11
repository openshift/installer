// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package schematics

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v3/core"
	"github.com/IBM/schematics-go-sdk/schematicsv1"
)

func ResourceIbmSchematicsAgentPrs() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmSchematicsAgentPrsCreate,
		ReadContext:   resourceIbmSchematicsAgentPrsRead,
		UpdateContext: resourceIbmSchematicsAgentPrsUpdate,
		DeleteContext: resourceIbmSchematicsAgentPrsDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"agent_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Agent ID to get the details of agent.",
			},
			"force": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Equivalent to -force options in the command line, default is false.",
			},
			"job_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Job Id.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The agent prs job updation time.",
			},
			"updated_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Email address of user who ran the agent prs job.",
			},
			"agent_version": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Agent version.",
			},
			"status_code": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Final result of the pre-requisite scanner job.",
			},
			"status_message": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The outcome of the pre-requisite scanner job, in a formatted log string.",
			},
			"log_url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL to the full pre-requisite scanner job logs.",
			},
		},
	}
}

func resourceIbmSchematicsAgentPrsCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentPrsCreate schematicsClient initialization failed: %s", err.Error()), "ibm_schematics_agent_prs", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	session, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentPrsCreate bluemixClient initialization failed: %s", err.Error()), "ibm_schematics_agent_prs", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	iamAccessToken := session.Config.IAMAccessToken
	iamRefreshToken := session.Config.IAMRefreshToken

	prsAgentJobOptions := &schematicsv1.PrsAgentJobOptions{}
	ff := map[string]string{
		"Authorization": iamAccessToken,
		"refresh_token": iamRefreshToken,
	}
	prsAgentJobOptions.Headers = ff
	prsAgentJobOptions.SetAgentID(d.Get("agent_id").(string))
	if _, ok := d.GetOk("force"); ok {
		prsAgentJobOptions.SetForce(d.Get("force").(bool))
	}

	agentPrsJob, response, err := schematicsClient.PrsAgentJobWithContext(context, prsAgentJobOptions)
	if err != nil {

		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentPrsCreate PrsAgentJobWithContext failed with error: %s and response:\n%s", err, response), "ibm_schematics_agent_prs", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *prsAgentJobOptions.AgentID, *agentPrsJob.JobID))

	return resourceIbmSchematicsAgentPrsRead(context, d, meta)
}

func resourceIbmSchematicsAgentPrsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentPrsRead schematicsClient initialization failed: %s", err.Error()), "ibm_schematics_agent_prs", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentPrsRead failed: %s", err.Error()), "ibm_schematics_agent_prs", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getAgentDataOptions := &schematicsv1.GetAgentDataOptions{
		Profile: core.StringPtr("detailed"),
	}

	getAgentDataOptions.SetAgentID(parts[0])

	agentData, response, err := schematicsClient.GetAgentDataWithContext(context, getAgentDataOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentPrsRead GetAgentDataWithContext failed with error: %s and response:\n%s", err, response), "ibm_schematics_agent_prs", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if agentData.RecentPrsJob != nil {

		if err = d.Set("agent_id", parts[0]); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentPrsRead failed with error: %s", err), "ibm_schematics_agent_prs", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("job_id", agentData.RecentPrsJob.JobID); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentPrsRead failed with error: %s", err), "ibm_schematics_agent_prs", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("updated_at", flex.DateTimeToString(agentData.RecentPrsJob.UpdatedAt)); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentPrsRead failed with error: %s", err), "ibm_schematics_agent_prs", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("updated_by", agentData.RecentPrsJob.UpdatedBy); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentPrsRead failed with error: %s", err), "ibm_schematics_agent_prs", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("agent_version", agentData.Version); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentPrsRead failed with error: %s", err), "ibm_schematics_agent_prs", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("status_code", agentData.RecentPrsJob.StatusCode); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentPrsRead failed with error: %s", err), "ibm_schematics_agent_prs", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("status_message", agentData.RecentPrsJob.StatusMessage); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentPrsRead failed with error: %s", err), "ibm_schematics_agent_prs", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("log_url", agentData.RecentPrsJob.LogURL); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentPrsRead failed with error: %s", err), "ibm_schematics_agent_prs", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}

	}
	return nil
}

func resourceIbmSchematicsAgentPrsUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentPrsUpdate schematicsClient initialization failed: %s", err.Error()), "ibm_schematics_agent_prs", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	session, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentPrsUpdate bluemixClient initialization failed: %s", err.Error()), "ibm_schematics_agent_prs", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	iamAccessToken := session.Config.IAMAccessToken
	iamRefreshToken := session.Config.IAMRefreshToken

	prsAgentJobOptions := &schematicsv1.PrsAgentJobOptions{}
	ff := map[string]string{
		"Authorization": iamAccessToken,
		"refresh_token": iamRefreshToken,
	}
	prsAgentJobOptions.Headers = ff

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentPrsUpdate failed: %s", err.Error()), "ibm_schematics_agent_prs", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	prsAgentJobOptions.SetAgentID(parts[0])

	hasChange := false

	if d.HasChange("agent_id") {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentPrsUpdate failed with error: %s", err), "ibm_schematics_agent_prs", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if d.HasChange("force") {
		prsAgentJobOptions.SetForce(d.Get("force").(bool))
		hasChange = true
	}

	if hasChange {
		agentPrsJob, response, err := schematicsClient.PrsAgentJobWithContext(context, prsAgentJobOptions)
		if err != nil {

			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentPrsUpdate PrsAgentJobWithContext failed with error: %s and response:\n%s", err, response), "ibm_schematics_agent_prs", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		d.SetId(fmt.Sprintf("%s/%s", *prsAgentJobOptions.AgentID, *agentPrsJob.JobID))
	}

	return resourceIbmSchematicsAgentPrsRead(context, d, meta)
}

func resourceIbmSchematicsAgentPrsDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
