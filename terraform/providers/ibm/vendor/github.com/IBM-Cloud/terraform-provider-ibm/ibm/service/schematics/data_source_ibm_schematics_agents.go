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
	"github.com/IBM/schematics-go-sdk/schematicsv1"
)

func DataSourceIbmSchematicsAgents() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmSchematicsAgentsRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the agent (must be unique, for an account).",
			},
			"total_count": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of records.",
			},
			"limit": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of records returned.",
			},
			"offset": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The skipped number of records.",
			},
			"agents": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of agents in the account.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the agent (must be unique, for an account).",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Agent description.",
						},
						"resource_group": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource-group name for the agent.  By default, Agent will be registered in Default Resource Group.",
						},
						"tags": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Tags for the agent.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"agent_location": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The location where agent is deployed in the user environment.",
						},
						"location": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "List of locations supported by IBM Cloud Schematics service.  While creating your workspace or action, choose the right region, since it cannot be changed.  Note, this does not limit the location of the IBM Cloud resources, provisioned using Schematics.",
						},
						"profile_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IAM trusted profile id, used by the Agent instance.",
						},
						"agent_crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Agent crn, obtained from the Schematics Agent deployment configuration.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Agent registration id.",
						},
						"registered_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Agent registration date-time.",
						},
						"registered_by": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The email address of an user who registered the Agent.",
						},
						"updated_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Agent registration updation time.",
						},
						"updated_by": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Email address of user who updated the Agent registration.",
						},
						"user_state": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "User defined status of the agent.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"state": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "User-defined states  * `enable`  Agent is enabled by the user.  * `disable` Agent is disbaled by the user.",
									},
									"set_by": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of the User who set the state of the Object.",
									},
									"set_at": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "When the User who set the state of the Object.",
									},
								},
							},
						},
						"connection_state": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Connection status of the agent.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"state": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Agent Connection Status  * `Connected` When Schematics is able to connect to the agent.  * `Disconnected` When Schematics is able not connect to the agent.",
									},
									"checked_at": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "When the connection state is modified.",
									},
								},
							},
						},
						"system_state": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Computed state of the agent.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"state": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Agent Status.",
									},
									"message": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The Agent status message.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmSchematicsAgentsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentsRead schematicsClient initialization failed: %s", err.Error()), "ibm_schematics_agents", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	listAgentOptions := &schematicsv1.ListAgentOptions{}

	agentList, response, err := schematicsClient.ListAgentWithContext(context, listAgentOptions)
	if err != nil {

		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentsRead ListAgentWithContext failed with error: %s and response:\n%s", err, response), "ibm_schematics_agents", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	// Use the provided filter argument and construct a new list with only the requested resource(s)
	var matchAgents []schematicsv1.Agent
	var name string
	var suppliedFilter bool

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
		suppliedFilter = true
		for _, data := range agentList.Agents {
			if *data.Name == name {
				matchAgents = append(matchAgents, data)
			}
		}
	} else {
		matchAgents = agentList.Agents
	}
	agentList.Agents = matchAgents

	if suppliedFilter {
		if len(agentList.Agents) == 0 {

			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentsRead failed with error: %s", err), "ibm_schematics_agents", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		d.SetId(name)
	} else {
		d.SetId(dataSourceIbmSchematicsAgentsID(d))
	}

	if err = d.Set("total_count", flex.IntValue(agentList.TotalCount)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentsRead failed with error: %s", err), "ibm_schematics_agents", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("limit", flex.IntValue(agentList.Limit)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentsRead failed with error: %s", err), "ibm_schematics_agents", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("offset", flex.IntValue(agentList.Offset)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentsRead failed with error: %s", err), "ibm_schematics_agents", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	agents := []map[string]interface{}{}
	if agentList.Agents != nil {
		for _, modelItem := range agentList.Agents {
			modelMap, err := dataSourceIbmSchematicsAgentsAgentToMap(&modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentsRead failed: %s", err.Error()), "ibm_schematics_agents", "read")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			agents = append(agents, modelMap)
		}
	}
	if err = d.Set("agents", agents); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentsRead failed with error: %s", err), "ibm_schematics_agents", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	return nil
}

// dataSourceIbmSchematicsAgentsID returns a reasonable ID for the list.
func dataSourceIbmSchematicsAgentsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIbmSchematicsAgentsAgentToMap(model *schematicsv1.Agent) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	if model.ResourceGroup != nil {
		modelMap["resource_group"] = *model.ResourceGroup
	}
	if model.Tags != nil {
		modelMap["tags"] = model.Tags
	}
	if model.AgentLocation != nil {
		modelMap["agent_location"] = *model.AgentLocation
	}
	if model.Location != nil {
		modelMap["location"] = *model.Location
	}
	if model.ProfileID != nil {
		modelMap["profile_id"] = *model.ProfileID
	}
	if model.AgentCrn != nil {
		modelMap["agent_crn"] = *model.AgentCrn
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.RegisteredAt != nil {
		modelMap["registered_at"] = model.RegisteredAt.String()
	}
	if model.RegisteredBy != nil {
		modelMap["registered_by"] = *model.RegisteredBy
	}
	if model.UpdatedAt != nil {
		modelMap["updated_at"] = model.UpdatedAt.String()
	}
	if model.UpdatedBy != nil {
		modelMap["updated_by"] = *model.UpdatedBy
	}
	if model.UserState != nil {
		userStateMap, err := dataSourceIbmSchematicsAgentsAgentUserStateToMap(model.UserState)
		if err != nil {
			return modelMap, err
		}
		modelMap["user_state"] = []map[string]interface{}{userStateMap}
	}
	if model.ConnectionState != nil {
		connectionStateMap, err := dataSourceIbmSchematicsAgentsConnectionStateToMap(model.ConnectionState)
		if err != nil {
			return modelMap, err
		}
		modelMap["connection_state"] = []map[string]interface{}{connectionStateMap}
	}
	if model.SystemState != nil {
		systemStateMap, err := dataSourceIbmSchematicsAgentsAgentSystemStateToMap(model.SystemState)
		if err != nil {
			return modelMap, err
		}
		modelMap["system_state"] = []map[string]interface{}{systemStateMap}
	}
	return modelMap, nil
}

func dataSourceIbmSchematicsAgentsAgentUserStateToMap(model *schematicsv1.AgentUserState) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.State != nil {
		modelMap["state"] = *model.State
	}
	if model.SetBy != nil {
		modelMap["set_by"] = *model.SetBy
	}
	if model.SetAt != nil {
		modelMap["set_at"] = model.SetAt.String()
	}
	return modelMap, nil
}

func dataSourceIbmSchematicsAgentsConnectionStateToMap(model *schematicsv1.ConnectionState) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.State != nil {
		modelMap["state"] = *model.State
	}
	if model.CheckedAt != nil {
		modelMap["checked_at"] = model.CheckedAt.String()
	}
	return modelMap, nil
}

func dataSourceIbmSchematicsAgentsAgentSystemStateToMap(model *schematicsv1.AgentSystemState) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.State != nil {
		modelMap["state"] = *model.State
	}
	if model.Message != nil {
		modelMap["message"] = *model.Message
	}
	return modelMap, nil
}
