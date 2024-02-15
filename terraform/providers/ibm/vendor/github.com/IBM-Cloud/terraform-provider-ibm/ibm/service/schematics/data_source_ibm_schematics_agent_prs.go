// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package schematics

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v3/core"
	"github.com/IBM/schematics-go-sdk/schematicsv1"
)

func DataSourceIbmSchematicsAgentPrs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmSchematicsAgentPrsRead,

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

func dataSourceIbmSchematicsAgentPrsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
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
	getAgentDataOptions.SetAgentID(d.Get("agent_id").(string))

	agentData, response, err := schematicsClient.GetAgentDataWithContext(context, getAgentDataOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId(DataSourceIBMSchematicsAgentID(d))
			return nil
		}
		log.Printf("[DEBUG] GetAgentDataWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetAgentDataWithContext failed %s\n%s", err, response))
	}
	d.SetId(DataSourceIBMSchematicsAgentID(d))

	if agentData.RecentPrsJob != nil {

		if err = d.Set("agent_id", getAgentDataOptions.AgentID); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting agent_id: %s", err))
		}
		if err = d.Set("job_id", agentData.RecentPrsJob.JobID); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting job_id: %s", err))
		} else if agentData.RecentPrsJob.JobID != nil {
			d.SetId(fmt.Sprintf("%s", *agentData.RecentPrsJob.JobID))
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

// DataSourceIBMSchematicsAgentID returns a reasonable ID
func DataSourceIBMSchematicsAgentID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
