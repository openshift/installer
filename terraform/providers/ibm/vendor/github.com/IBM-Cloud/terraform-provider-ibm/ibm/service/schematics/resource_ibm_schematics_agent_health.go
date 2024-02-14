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
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/schematics-go-sdk/schematicsv1"
)

func ResourceIbmSchematicsAgentHealth() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmSchematicsAgentHealthCreate,
		ReadContext:   resourceIbmSchematicsAgentHealthRead,
		UpdateContext: resourceIbmSchematicsAgentHealthUpdate,
		DeleteContext: resourceIbmSchematicsAgentHealthDelete,
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
				Description: "The agent health check job updation time.",
			},
			"updated_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Email address of user who ran the agent health check job.",
			},
			"agent_version": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Agent version.",
			},
			"status_code": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Final result of the health-check job.",
			},
			"status_message": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The outcome of the health-check job, in a formatted log string.",
			},
			"log_url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL to the full health-check job logs.",
			},
		},
	}
}

func resourceIbmSchematicsAgentHealthCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	healthCheckAgentJobOptions := &schematicsv1.HealthCheckAgentJobOptions{}
	ff := map[string]string{
		"X-Feature-Agents": "true",
		"Authorization":    iamAccessToken,
		"refresh_token":    iamRefreshToken,
	}
	healthCheckAgentJobOptions.Headers = ff

	healthCheckAgentJobOptions.SetAgentID(d.Get("agent_id").(string))
	if _, ok := d.GetOk("force"); ok {
		healthCheckAgentJobOptions.SetForce(d.Get("force").(bool))
	}

	agentHealthJob, response, err := schematicsClient.HealthCheckAgentJobWithContext(context, healthCheckAgentJobOptions)
	if err != nil {
		log.Printf("[DEBUG] HealthCheckAgentJobWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("HealthCheckAgentJobWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *healthCheckAgentJobOptions.AgentID, *agentHealthJob.JobID))

	return resourceIbmSchematicsAgentHealthRead(context, d, meta)
}

func resourceIbmSchematicsAgentHealthRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if agentData.RecentHealthJob != nil {

		if err = d.Set("agent_id", parts[0]); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting agent_id: %s", err))
		}
		if err = d.Set("job_id", agentData.RecentHealthJob.JobID); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting job_id: %s", err))
		}
		if err = d.Set("updated_at", flex.DateTimeToString(agentData.RecentHealthJob.UpdatedAt)); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
		}
		if err = d.Set("updated_by", agentData.RecentHealthJob.UpdatedBy); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting updated_by: %s", err))
		}
		if err = d.Set("agent_version", agentData.Version); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting agent_version: %s", err))
		}
		if err = d.Set("status_code", agentData.RecentHealthJob.StatusCode); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting status_code: %s", err))
		}
		if err = d.Set("status_message", agentData.RecentHealthJob.StatusMessage); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting status_message: %s", err))
		}
		if err = d.Set("log_url", agentData.RecentHealthJob.LogURL); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting log_url: %s", err))
		}

	}

	return nil
}

func resourceIbmSchematicsAgentHealthUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	healthCheckAgentJobOptions := &schematicsv1.HealthCheckAgentJobOptions{}
	ff := map[string]string{
		"X-Feature-Agents": "true",
		"Authorization":    iamAccessToken,
		"refresh_token":    iamRefreshToken,
	}
	healthCheckAgentJobOptions.Headers = ff

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	healthCheckAgentJobOptions.SetAgentID(parts[0])

	hasChange := false

	if d.HasChange("agent_id") {
		return diag.FromErr(fmt.Errorf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "agent_id"))
	}
	if d.HasChange("force") {
		healthCheckAgentJobOptions.SetForce(d.Get("force").(bool))
		hasChange = true
	}

	if hasChange {
		agentHealthJob, response, err := schematicsClient.HealthCheckAgentJobWithContext(context, healthCheckAgentJobOptions)
		if err != nil {
			log.Printf("[DEBUG] HealthCheckAgentJobWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("HealthCheckAgentJobWithContext failed %s\n%s", err, response))
		}
		d.SetId(fmt.Sprintf("%s/%s", *healthCheckAgentJobOptions.AgentID, *agentHealthJob.JobID))
	}

	return resourceIbmSchematicsAgentHealthRead(context, d, meta)
}

func resourceIbmSchematicsAgentHealthDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
