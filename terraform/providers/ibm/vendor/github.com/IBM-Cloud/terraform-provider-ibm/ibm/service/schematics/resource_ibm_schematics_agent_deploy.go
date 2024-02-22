// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package schematics

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
	"github.com/IBM/go-sdk-core/v3/core"
	"github.com/IBM/schematics-go-sdk/schematicsv1"
)

func ResourceIbmSchematicsAgentDeploy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmSchematicsAgentDeployCreate,
		ReadContext:   resourceIbmSchematicsAgentDeployRead,
		UpdateContext: resourceIbmSchematicsAgentDeployUpdate,
		DeleteContext: resourceIbmSchematicsAgentDeployDelete,
		Importer:      &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

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
				Description: "The agent deploy job updation time.",
			},
			"updated_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Email address of user who ran the agent deploy job.",
			},
			"is_redeployed": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "True, when the same version of the agent was redeployed.",
			},
			"agent_version": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Agent version.",
			},
			"status_code": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Final result of the agent deployment job.",
			},
			"status_message": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The outcome of the agent deployment job, in a formatted log string.",
			},
			"log_url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL to the full agent deployment job logs.",
			},
		},
	}
}

func resourceIbmSchematicsAgentDeployCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	deployAgentJobOptions := &schematicsv1.DeployAgentJobOptions{}
	ff := map[string]string{
		"X-Feature-Agents": "true",
		"Authorization":    iamAccessToken,
		"refresh_token":    iamRefreshToken,
	}
	deployAgentJobOptions.Headers = ff
	deployAgentJobOptions.SetAgentID(d.Get("agent_id").(string))
	if _, ok := d.GetOk("force"); ok {
		deployAgentJobOptions.SetForce(d.Get("force").(bool))
	}

	agentDeployJob, response, err := schematicsClient.DeployAgentJobWithContext(context, deployAgentJobOptions)
	if err != nil {
		log.Printf("[DEBUG] DeployAgentJobWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeployAgentJobWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *deployAgentJobOptions.AgentID, *agentDeployJob.JobID))
	log.Printf("[INFO] Agent : %s", *deployAgentJobOptions.AgentID)

	_, err = isWaitForAgentAvailable(context, schematicsClient, *deployAgentJobOptions.AgentID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(fmt.Errorf("Waiting for agent to be available failed %s", err))
	}

	return resourceIbmSchematicsAgentDeployRead(context, d, meta)
}

const (
	agentProvisioningTriggered = "Triggered deployment"
	agentProvisioningDone      = "success"
	agentProvisioningPending   = "PENDING"
	agentProvisioninFailed     = "Job Failed"
)

func isWaitForAgentAvailable(context context.Context, schematicsClient *schematicsv1.SchematicsV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for agent (%s) to be available.", id)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", agentProvisioningPending, agentProvisioningTriggered},
		Target:     []string{agentProvisioningDone, agentProvisioninFailed, ""},
		Refresh:    agentRefreshFunc(schematicsClient, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	return stateConf.WaitForStateContext(context)
}
func agentRefreshFunc(schematicsClient *schematicsv1.SchematicsV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getAgentDataOptions := &schematicsv1.GetAgentDataOptions{
			AgentID: core.StringPtr(id),
			Profile: core.StringPtr("detailed"),
		}
		ff := map[string]string{
			"X-Feature-Agents": "true",
		}
		getAgentDataOptions.Headers = ff

		agent, response, err := schematicsClient.GetAgentData(getAgentDataOptions)
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] Error Getting Agent: %s\n%s", err, response)
		}
		if *agent.RecentDeployJob.StatusMessage == agentProvisioninFailed || *agent.RecentDeployJob.StatusMessage == agentProvisioningDone {
			return agent, agentProvisioningDone, nil
		}
		return agent, agentProvisioningPending, nil
	}
}

func resourceIbmSchematicsAgentDeployRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		return diag.FromErr(err)
	}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getAgentDataOptions := &schematicsv1.GetAgentDataOptions{
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
	if agentData.RecentDeployJob != nil {

		if err = d.Set("agent_id", parts[0]); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting agent_id: %s", err))
		}
		if err = d.Set("job_id", agentData.RecentDeployJob.JobID); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting job_id: %s", err))
		}
		if err = d.Set("updated_at", flex.DateTimeToString(agentData.RecentDeployJob.UpdatedAt)); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
		}
		if err = d.Set("updated_by", agentData.RecentDeployJob.UpdatedBy); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting updated_by: %s", err))
		}
		if err = d.Set("agent_version", agentData.Version); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting agent_version: %s", err))
		}
		if err = d.Set("status_code", agentData.RecentDeployJob.StatusCode); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting status_code: %s", err))
		}
		if err = d.Set("status_message", agentData.RecentDeployJob.StatusMessage); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting status_message: %s", err))
		}
		if err = d.Set("log_url", agentData.RecentDeployJob.LogURL); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting log_url: %s", err))
		}

	}
	return nil
}

func resourceIbmSchematicsAgentDeployUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	deployAgentJobOptions := &schematicsv1.DeployAgentJobOptions{}
	ff := map[string]string{
		"X-Feature-Agents": "true",
		"Authorization":    iamAccessToken,
		"refresh_token":    iamRefreshToken,
	}
	deployAgentJobOptions.Headers = ff

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	deployAgentJobOptions.SetAgentID(parts[0])

	hasChange := false

	if d.HasChange("agent_id") {
		return diag.FromErr(fmt.Errorf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "agent_id"))
	}
	if d.HasChange("force") {
		deployAgentJobOptions.SetForce(d.Get("force").(bool))
		hasChange = true
	}

	if hasChange {
		agentDeployJob, response, err := schematicsClient.DeployAgentJobWithContext(context, deployAgentJobOptions)
		if err != nil {
			log.Printf("[DEBUG] DeployAgentJobWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("DeployAgentJobWithContext failed %s\n%s", err, response))
		}
		d.SetId(fmt.Sprintf("%s/%s", *deployAgentJobOptions.AgentID, *agentDeployJob.JobID))

		_, err = isWaitForAgentAvailable(context, schematicsClient, parts[0], d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Waiting for agent to be available failed %s", err))
		}
	}

	return resourceIbmSchematicsAgentDeployRead(context, d, meta)
}

func resourceIbmSchematicsAgentDeployDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
