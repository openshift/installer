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
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
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
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentDeployCreate schematicsClient initialization failed: %s", err.Error()), "ibm_schematics_agent_deploy", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	session, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentDeployCreate bluemixClient initialization failed: %s", err.Error()), "ibm_schematics_agent_deploy", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	iamAccessToken := session.Config.IAMAccessToken
	iamRefreshToken := session.Config.IAMRefreshToken

	deployAgentJobOptions := &schematicsv1.DeployAgentJobOptions{}
	ff := map[string]string{
		"Authorization": iamAccessToken,
		"refresh_token": iamRefreshToken,
	}
	deployAgentJobOptions.Headers = ff
	deployAgentJobOptions.SetAgentID(d.Get("agent_id").(string))
	if _, ok := d.GetOk("force"); ok {
		deployAgentJobOptions.SetForce(d.Get("force").(bool))
	}

	agentDeployJob, response, err := schematicsClient.DeployAgentJobWithContext(context, deployAgentJobOptions)
	if err != nil {

		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentDeployCreate DeployAgentJobWithContext failed with error: %s and response:\n%s", err, response), "ibm_schematics_agent_deploy", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *deployAgentJobOptions.AgentID, *agentDeployJob.JobID))
	log.Printf("[INFO] Agent : %s", *deployAgentJobOptions.AgentID)

	_, err = isWaitForAgentAvailable(context, schematicsClient, *deployAgentJobOptions.AgentID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentDeployCreate failed with error: %s", err), "ibm_schematics_agent_deploy", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	return resourceIbmSchematicsAgentDeployRead(context, d, meta)
}

const (
	agentProvisioningStatusCodeJobCancelled      = "job_cancelled"
	agentProvisioningStatusCodeJobFailed         = "job_failed"
	agentProvisioningStatusCodeJobFinished       = "job_finished"
	agentProvisioningStatusCodeJobInProgress     = "job_in_progress"
	agentProvisioningStatusCodeJobPending        = "job_pending"
	agentProvisioningStatusCodeJobReadyToExecute = "job_ready_to_execute"
	agentProvisioningStatusCodeJobStopInProgress = "job_stop_in_progress"
	agentProvisioningStatusCodeJobStopped        = "job_stopped"
)

func isWaitForAgentAvailable(context context.Context, schematicsClient *schematicsv1.SchematicsV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for agent (%s) to be available.", id)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", agentProvisioningStatusCodeJobInProgress, agentProvisioningStatusCodeJobPending, agentProvisioningStatusCodeJobReadyToExecute, agentProvisioningStatusCodeJobStopInProgress},
		Target:     []string{agentProvisioningStatusCodeJobFinished, agentProvisioningStatusCodeJobFailed, agentProvisioningStatusCodeJobCancelled, agentProvisioningStatusCodeJobStopped, ""},
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

		agent, response, err := schematicsClient.GetAgentData(getAgentDataOptions)
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] Error Getting Agent: %s\n%s", err, response)
		}
		if agent.RecentDeployJob.StatusCode != nil {
			return agent, *agent.RecentDeployJob.StatusCode, nil
		}
		return agent, agentProvisioningStatusCodeJobPending, nil
	}
}

func resourceIbmSchematicsAgentDeployRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentDeployRead schematicsClient initialization failed: %s", err.Error()), "ibm_schematics_agent_deploy", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentDeployRead bluemixClient initialization failed: %s", err.Error()), "ibm_schematics_agent_deploy", "read")
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

		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentDeployRead GetAgentDataWithContext failed with error: %s and response:\n%s", err, response), "ibm_schematics_agent_deploy", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if agentData.RecentDeployJob != nil {

		if err = d.Set("agent_id", parts[0]); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentDeployRead failed with error: %s", err), "ibm_schematics_agent_deploy", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("job_id", agentData.RecentDeployJob.JobID); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentDeployRead failed with error: %s", err), "ibm_schematics_agent_deploy", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("updated_at", flex.DateTimeToString(agentData.RecentDeployJob.UpdatedAt)); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentDeployRead failed with error: %s", err), "ibm_schematics_agent_deploy", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("updated_by", agentData.RecentDeployJob.UpdatedBy); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentDeployRead failed with error: %s", err), "ibm_schematics_agent_deploy", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("agent_version", agentData.Version); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentDeployRead failed with error: %s", err), "ibm_schematics_agent_deploy", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("status_code", agentData.RecentDeployJob.StatusCode); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentDeployRead failed with error: %s", err), "ibm_schematics_agent_deploy", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("status_message", agentData.RecentDeployJob.StatusMessage); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentDeployRead failed with error: %s", err), "ibm_schematics_agent_deploy", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("log_url", agentData.RecentDeployJob.LogURL); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentDeployRead failed with error: %s", err), "ibm_schematics_agent_deploy", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}

	}
	return nil
}

func resourceIbmSchematicsAgentDeployUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentDeployUpdate schematicsClient initialization failed: %s", err.Error()), "ibm_schematics_agent_deploy", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	session, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentDeployUpdate bluemixClient initialization failed: %s", err.Error()), "ibm_schematics_agent_deploy", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	iamAccessToken := session.Config.IAMAccessToken
	iamRefreshToken := session.Config.IAMRefreshToken
	deployAgentJobOptions := &schematicsv1.DeployAgentJobOptions{}
	ff := map[string]string{
		"Authorization": iamAccessToken,
		"refresh_token": iamRefreshToken,
	}
	deployAgentJobOptions.Headers = ff

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentDeployUpdate failed: %s", err.Error()), "ibm_schematics_agent_deploy", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deployAgentJobOptions.SetAgentID(parts[0])

	hasChange := false

	if d.HasChange("agent_id") {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentDeployUpdate failed with error: %s", err), "ibm_schematics_agent_deploy", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if d.HasChange("force") {
		deployAgentJobOptions.SetForce(d.Get("force").(bool))
		hasChange = true
	}

	if hasChange {
		agentDeployJob, response, err := schematicsClient.DeployAgentJobWithContext(context, deployAgentJobOptions)
		if err != nil {

			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentDeployUpdate DeployAgentJobWithContext failed with error: %s and response:\n%s", err, response), "ibm_schematics_agent_deploy", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		d.SetId(fmt.Sprintf("%s/%s", *deployAgentJobOptions.AgentID, *agentDeployJob.JobID))

		_, err = isWaitForAgentAvailable(context, schematicsClient, parts[0], d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentDeployUpdate failed with error: %s", err), "ibm_schematics_agent_deploy", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIbmSchematicsAgentDeployRead(context, d, meta)
}

func resourceIbmSchematicsAgentDeployDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
