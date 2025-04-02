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

func DataSourceIbmSchematicsAgentHealth() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmSchematicsAgentHealthRead,

		Schema: map[string]*schema.Schema{
			"agent_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Agent ID to get the details of agent.",
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

func dataSourceIbmSchematicsAgentHealthRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentHealthRead schematicsClient initialization failed: %s", err.Error()), "ibm_schematics_agent_health", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getAgentDataOptions := &schematicsv1.GetAgentDataOptions{
		Profile: core.StringPtr("detailed"),
	}
	getAgentDataOptions.SetAgentID(d.Get("agent_id").(string))

	agentData, response, err := schematicsClient.GetAgentDataWithContext(context, getAgentDataOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId(DataSourceIBMSchematicsAgentID(d))
			return nil
		}

		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentHealthRead GetAgentDataWithContext failed with error: %s and response:\n%s", err, response), "ibm_schematics_agent_health", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	d.SetId(DataSourceIBMSchematicsAgentID(d))

	if agentData.RecentHealthJob != nil {

		if err = d.Set("agent_id", getAgentDataOptions.AgentID); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentHealthRead failed with error: %s", err), "ibm_schematics_agent_health", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("job_id", agentData.RecentHealthJob.JobID); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentHealthRead failed with error: %s", err), "ibm_schematics_agent_health", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		} else if agentData.RecentHealthJob.JobID != nil {
			d.SetId(fmt.Sprintf("%s", *agentData.RecentHealthJob.JobID))
		}
		if err = d.Set("updated_at", flex.DateTimeToString(agentData.RecentHealthJob.UpdatedAt)); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentHealthRead failed with error: %s", err), "ibm_schematics_agent_health", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("updated_by", agentData.RecentHealthJob.UpdatedBy); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentHealthRead failed with error: %s", err), "ibm_schematics_agent_health", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("agent_version", agentData.Version); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentHealthRead failed with error: %s", err), "ibm_schematics_agent_health", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("status_code", agentData.RecentHealthJob.StatusCode); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentHealthRead failed with error: %s", err), "ibm_schematics_agent_health", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("status_message", agentData.RecentHealthJob.StatusMessage); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentHealthRead failed with error: %s", err), "ibm_schematics_agent_health", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("log_url", agentData.RecentHealthJob.LogURL); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentHealthRead failed with error: %s", err), "ibm_schematics_agent_health", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}

	}
	return nil
}
