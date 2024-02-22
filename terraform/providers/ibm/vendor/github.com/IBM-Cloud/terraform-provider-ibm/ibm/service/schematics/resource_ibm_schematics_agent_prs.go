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
		return diag.FromErr(err)
	}
	session, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return diag.FromErr(err)
	}
	iamAccessToken := session.Config.IAMAccessToken
	iamRefreshToken := session.Config.IAMRefreshToken

	prsAgentJobOptions := &schematicsv1.PrsAgentJobOptions{}
	ff := map[string]string{
		"X-Feature-Agents": "true",
		"Authorization":    iamAccessToken,
		"refresh_token":    iamRefreshToken,
	}
	prsAgentJobOptions.Headers = ff
	prsAgentJobOptions.SetAgentID(d.Get("agent_id").(string))
	if _, ok := d.GetOk("force"); ok {
		prsAgentJobOptions.SetForce(d.Get("force").(bool))
	}

	agentPrsJob, response, err := schematicsClient.PrsAgentJobWithContext(context, prsAgentJobOptions)
	if err != nil {
		log.Printf("[DEBUG] PrsAgentJobWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("PrsAgentJobWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *prsAgentJobOptions.AgentID, *agentPrsJob.JobID))

	return resourceIbmSchematicsAgentPrsRead(context, d, meta)
}

func resourceIbmSchematicsAgentPrsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		return diag.FromErr(err)
	}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getAgentDataOptions := &schematicsv1.GetAgentDataOptions{
		// XFeatureAgents: core.BoolPtr(true),
		Profile: core.StringPtr("detailed"),
	}
	ff := map[string]string{
		"X-Feature-Agents": "true",
	}
	getAgentDataOptions.Headers = ff

	getAgentDataOptions.SetAgentID(parts[0])

	agentData, response, err := schematicsClient.GetAgentDataWithContext(context, getAgentDataOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetAgentDataWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetAgentDataWithContext failed %s\n%s", err, response))
	}
	if agentData.RecentPrsJob != nil {

		if err = d.Set("agent_id", parts[0]); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting agent_id: %s", err))
		}
		if err = d.Set("job_id", agentData.RecentPrsJob.JobID); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting job_id: %s", err))
		}
		if err = d.Set("updated_at", flex.DateTimeToString(agentData.RecentPrsJob.UpdatedAt)); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
		}
		if err = d.Set("updated_by", agentData.RecentPrsJob.UpdatedBy); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting updated_by: %s", err))
		}
		if err = d.Set("agent_version", agentData.Version); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting agent_version: %s", err))
		}
		if err = d.Set("status_code", agentData.RecentPrsJob.StatusCode); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting status_code: %s", err))
		}
		if err = d.Set("status_message", agentData.RecentPrsJob.StatusMessage); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting status_message: %s", err))
		}
		if err = d.Set("log_url", agentData.RecentPrsJob.LogURL); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting log_url: %s", err))
		}

	}
	return nil
}

func resourceIbmSchematicsAgentPrsUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		return diag.FromErr(err)
	}
	session, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return diag.FromErr(err)
	}
	iamAccessToken := session.Config.IAMAccessToken
	iamRefreshToken := session.Config.IAMRefreshToken

	prsAgentJobOptions := &schematicsv1.PrsAgentJobOptions{}
	ff := map[string]string{
		"X-Feature-Agents": "true",
		"Authorization":    iamAccessToken,
		"refresh_token":    iamRefreshToken,
	}
	prsAgentJobOptions.Headers = ff

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	prsAgentJobOptions.SetAgentID(parts[0])

	hasChange := false

	if d.HasChange("agent_id") {
		return diag.FromErr(fmt.Errorf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "agent_id"))
	}
	if d.HasChange("force") {
		prsAgentJobOptions.SetForce(d.Get("force").(bool))
		hasChange = true
	}

	if hasChange {
		agentPrsJob, response, err := schematicsClient.PrsAgentJobWithContext(context, prsAgentJobOptions)
		if err != nil {
			log.Printf("[DEBUG] PrsAgentJobWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("PrsAgentJobWithContext failed %s\n%s", err, response))
		}
		d.SetId(fmt.Sprintf("%s/%s", *prsAgentJobOptions.AgentID, *agentPrsJob.JobID))
	}

	return resourceIbmSchematicsAgentPrsRead(context, d, meta)
}

func resourceIbmSchematicsAgentPrsDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
