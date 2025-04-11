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
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentHealthCreate schematicsClient initialization failed: %s", err.Error()), "ibm_schematics_agent_health", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	session, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentHealthCreate bluemixClient initialization failed: %s", err.Error()), "ibm_schematics_agent_health", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	iamAccessToken := session.Config.IAMAccessToken
	iamRefreshToken := session.Config.IAMRefreshToken

	healthCheckAgentJobOptions := &schematicsv1.HealthCheckAgentJobOptions{}
	ff := map[string]string{
		"Authorization": iamAccessToken,
		"refresh_token": iamRefreshToken,
	}
	healthCheckAgentJobOptions.Headers = ff

	healthCheckAgentJobOptions.SetAgentID(d.Get("agent_id").(string))
	if _, ok := d.GetOk("force"); ok {
		healthCheckAgentJobOptions.SetForce(d.Get("force").(bool))
	}

	agentHealthJob, response, err := schematicsClient.HealthCheckAgentJobWithContext(context, healthCheckAgentJobOptions)
	if err != nil {

		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentHealthCreate HealthCheckAgentJobWithContext failed with error: %s and response:\n%s", err, response), "ibm_schematics_agent_health", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *healthCheckAgentJobOptions.AgentID, *agentHealthJob.JobID))

	return resourceIbmSchematicsAgentHealthRead(context, d, meta)
}

func resourceIbmSchematicsAgentHealthRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentHealthRead schematicsClient initialization failed: %s", err.Error()), "ibm_schematics_agent_health", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentHealthRead failed: %s", err.Error()), "ibm_schematics_agent_health", "read")
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

		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentHealthRead GetAgentDataWithContext failed with error: %s and response:\n%s", err, response), "ibm_schematics_agent_health", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if agentData.RecentHealthJob != nil {

		if err = d.Set("agent_id", parts[0]); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentHealthRead failed with error: %s", err), "ibm_schematics_agent_health", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("job_id", agentData.RecentHealthJob.JobID); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentHealthRead failed with error: %s", err), "ibm_schematics_agent_health", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("updated_at", flex.DateTimeToString(agentData.RecentHealthJob.UpdatedAt)); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentHealthRead failed with error: %s", err), "ibm_schematics_agent_health", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("updated_by", agentData.RecentHealthJob.UpdatedBy); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentHealthRead failed with error: %s", err), "ibm_schematics_agent_health", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("agent_version", agentData.Version); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentHealthRead failed with error: %s", err), "ibm_schematics_agent_health", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("status_code", agentData.RecentHealthJob.StatusCode); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentHealthRead failed with error: %s", err), "ibm_schematics_agent_health", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("status_message", agentData.RecentHealthJob.StatusMessage); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentHealthRead failed with error: %s", err), "ibm_schematics_agent_health", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("log_url", agentData.RecentHealthJob.LogURL); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentHealthRead failed with error: %s", err), "ibm_schematics_agent_health", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}

	}

	return nil
}

func resourceIbmSchematicsAgentHealthUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentHealthUpdate schematicsClient initialization failed: %s", err.Error()), "ibm_schematics_agent_health", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	session, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentHealthUpdate bluemixClient initialization failed: %s", err.Error()), "ibm_schematics_agent_health", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	iamAccessToken := session.Config.IAMAccessToken
	iamRefreshToken := session.Config.IAMRefreshToken

	healthCheckAgentJobOptions := &schematicsv1.HealthCheckAgentJobOptions{}
	ff := map[string]string{
		"Authorization": iamAccessToken,
		"refresh_token": iamRefreshToken,
	}
	healthCheckAgentJobOptions.Headers = ff

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentHealthUpdate failed: %s", err.Error()), "ibm_schematics_agent_health", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	healthCheckAgentJobOptions.SetAgentID(parts[0])

	hasChange := false

	if d.HasChange("agent_id") {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentHealthUpdate failed with error: %s", err), "ibm_schematics_agent_health", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if d.HasChange("force") {
		healthCheckAgentJobOptions.SetForce(d.Get("force").(bool))
		hasChange = true
	}

	if hasChange {
		agentHealthJob, response, err := schematicsClient.HealthCheckAgentJobWithContext(context, healthCheckAgentJobOptions)
		if err != nil {

			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentHealthUpdate HealthCheckAgentJobWithContext failed with error: %s and response:\n%s", err, response), "ibm_schematics_agent_health", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		d.SetId(fmt.Sprintf("%s/%s", *healthCheckAgentJobOptions.AgentID, *agentHealthJob.JobID))
	}

	return resourceIbmSchematicsAgentHealthRead(context, d, meta)
}

func resourceIbmSchematicsAgentHealthDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
