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

func DataSourceIbmSchematicsAgent() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmSchematicsAgentRead,

		Schema: map[string]*schema.Schema{
			"agent_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Agent ID to get the details of agent.",
			},
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
				Description: "The resource-group name for the agent.  By default, agent will be registered in Default Resource Group.",
			},
			"tags": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Tags for the agent.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"version": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Agent version.",
			},
			"schematics_location": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "List of locations supported by IBM Cloud Schematics service.  While creating your workspace or action, choose the right region, since it cannot be changed.  Note, this does not limit the location of the IBM Cloud resources, provisioned using Schematics.",
			},
			"agent_location": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The location where agent is deployed in the user environment.",
			},
			"agent_infrastructure": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The infrastructure parameters used by the agent.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"infra_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of target agent infrastructure.",
						},
						"cluster_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The cluster ID where agent services will be running.",
						},
						"cluster_resource_group": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource group of the cluster (is it required?).",
						},
						"cos_instance_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The COS instance name to store the agent logs.",
						},
						"cos_bucket_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The COS bucket name used to store the logs.",
						},
						"cos_bucket_region": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The COS bucket region.",
						},
					},
				},
			},
			"agent_metadata": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The metadata of an agent.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the metadata.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Value of the metadata name.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"agent_inputs": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Additional input variables for the agent.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the variable. For example, `name = \"inventory username\"`.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The value for the variable or reference to the value. For example, `value = \"<provide your ssh_key_value with \n>\"`. **Note** The SSH key should contain `\n` at the end of the key details in case of command line or API calls.",
						},
						"use_default": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "True, will ignore the data in the value attribute, instead the data in metadata.default_value will be used.",
						},
						"metadata": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "An user editable metadata for the variables.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of the variable.",
									},
									"aliases": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The list of aliases for the variable name.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"description": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The description of the meta data.",
									},
									"cloud_data_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cloud data type of the variable. eg. resource_group_id, region, vpc_id.",
									},
									"default_value": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Default value for the variable only if the override value is not specified.",
									},
									"link_status": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The status of the link.",
									},
									"secure": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Is the variable secure or sensitive ?.",
									},
									"immutable": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Is the variable readonly ?.",
									},
									"hidden": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "If **true**, the variable is not displayed on UI or Command line.",
									},
									"required": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "If the variable required?.",
									},
									"options": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The list of possible values for this variable.  If type is **integer** or **date**, then the array of string is  converted to array of integers or date during the runtime.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"min_value": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The minimum value of the variable. Applicable for the integer type.",
									},
									"max_value": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The maximum value of the variable. Applicable for the integer type.",
									},
									"min_length": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The minimum length of the variable value. Applicable for the string type.",
									},
									"max_length": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The maximum length of the variable value. Applicable for the string type.",
									},
									"matches": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The regex for the variable value.",
									},
									"position": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The relative position of this variable in a list.",
									},
									"group_by": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The display name of the group this variable belongs to.",
									},
									"source": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The source of this meta-data.",
									},
								},
							},
						},
						"link": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The reference link to the variable value By default the expression points to `$self.value`.",
						},
					},
				},
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
			"agent_crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The agent crn, obtained from the Schematics agent deployment configuration.",
			},
			"id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The agent resource id.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The agent creation date-time.",
			},
			"creation_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The email address of an user who created the agent.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The agent registration updation time.",
			},
			"updated_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Email address of user who updated the agent registration.",
			},
			"system_state": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Computed state of the agent.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status_code": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Agent Status.",
						},
						"status_message": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The agent status message.",
						},
					},
				},
			},
			"agent_kpi": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Schematics Agent key performance indicators.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_indicator": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Overall availability indicator reported by the agent.",
						},
						"lifecycle_indicator": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Overall lifecycle indicator reported by the agents.",
						},
						"percent_usage_indicator": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Percentage usage of the agent resources.",
						},
						"application_indicators": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Agent application key performance indicators.",
							Elem: &schema.Schema{
								Type: schema.TypeMap,
							},
						},
						"infra_indicators": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Agent infrastructure key performance indicators.",
							Elem: &schema.Schema{
								Type: schema.TypeMap,
							},
						},
					},
				},
			},
			"recent_prs_job": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Run a pre-requisite scanner for deploying agent.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"agent_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of the agent.",
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
				},
			},
			"recent_deploy_job": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Post-installations checks for Agent health.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"agent_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of the agent.",
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
				},
			},
			"recent_health_job": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Agent health check.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"agent_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of the agent.",
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
				},
			},
		},
	}
}

func dataSourceIbmSchematicsAgentRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentRead schematicsClient initialization failed: %s", err.Error()), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getAgentDataOptions := &schematicsv1.GetAgentDataOptions{
		Profile: core.StringPtr("detailed"),
	}

	getAgentDataOptions.SetAgentID(d.Get("agent_id").(string))

	agentData, response, err := schematicsClient.GetAgentDataWithContext(context, getAgentDataOptions)
	if err != nil {

		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentRead GetAgentDataWithContext failed with error: %s and response:\n%s", err, response), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s", *getAgentDataOptions.AgentID))

	if err = d.Set("name", agentData.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("description", agentData.Description); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("resource_group", agentData.ResourceGroup); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("version", agentData.Version); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("schematics_location", agentData.SchematicsLocation); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("agent_location", agentData.AgentLocation); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if agentData.Tags != nil {
		if err = d.Set("tags", agentData.Tags); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	agentInfrastructure := []map[string]interface{}{}
	if agentData.AgentInfrastructure != nil {
		modelMap, err := dataSourceIbmSchematicsAgentAgentInfrastructureToMap(agentData.AgentInfrastructure)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentRead failed: %s", err.Error()), "ibm_schematics_agent", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		agentInfrastructure = append(agentInfrastructure, modelMap)
	}
	if err = d.Set("agent_infrastructure", agentInfrastructure); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	agentMetadata := []map[string]interface{}{}
	if agentData.AgentMetadata != nil {
		for _, modelItem := range agentData.AgentMetadata {
			modelMap, err := dataSourceIbmSchematicsAgentAgentMetadataInfoToMap(&modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentRead failed: %s", err.Error()), "ibm_schematics_agent", "read")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			agentMetadata = append(agentMetadata, modelMap)
		}
	}
	if err = d.Set("agent_metadata", agentMetadata); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	agentInputs := []map[string]interface{}{}
	if agentData.AgentInputs != nil {
		for _, modelItem := range agentData.AgentInputs {
			modelMap, err := dataSourceIbmSchematicsAgentVariableDataToMap(&modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentRead failed: %s", err.Error()), "ibm_schematics_agent", "read")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			agentInputs = append(agentInputs, modelMap)
		}
	}
	if err = d.Set("agent_inputs", agentInputs); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	userState := []map[string]interface{}{}
	if agentData.UserState != nil {
		modelMap, err := dataSourceIbmSchematicsAgentAgentUserStateToMap(agentData.UserState)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentRead failed: %s", err.Error()), "ibm_schematics_agent", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		userState = append(userState, modelMap)
	}
	if err = d.Set("user_state", userState); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("agent_crn", agentData.AgentCrn); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("id", agentData.ID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("created_at", flex.DateTimeToString(agentData.CreatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("creation_by", agentData.CreationBy); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("updated_at", flex.DateTimeToString(agentData.UpdatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("updated_by", agentData.UpdatedBy); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	systemState := []map[string]interface{}{}
	if agentData.SystemState != nil {
		modelMap, err := dataSourceIbmSchematicsAgentAgentSystemStatusToMap(agentData.SystemState)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentRead failed: %s", err.Error()), "ibm_schematics_agent", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		systemState = append(systemState, modelMap)
	}
	if err = d.Set("system_state", systemState); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	agentKpi := []map[string]interface{}{}
	if agentData.AgentKpi != nil {
		modelMap, err := dataSourceIbmSchematicsAgentAgentKPIDataToMap(agentData.AgentKpi)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentRead failed: %s", err.Error()), "ibm_schematics_agent", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		agentKpi = append(agentKpi, modelMap)
	}
	if err = d.Set("agent_kpi", agentKpi); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	recentPrsJob := []map[string]interface{}{}
	if agentData.RecentPrsJob != nil {
		modelMap, err := dataSourceIbmSchematicsAgentAgentDataRecentPrsJobToMap(agentData.RecentPrsJob)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentRead failed: %s", err.Error()), "ibm_schematics_agent", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		recentPrsJob = append(recentPrsJob, modelMap)
	}
	if err = d.Set("recent_prs_job", recentPrsJob); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	recentDeployJob := []map[string]interface{}{}
	if agentData.RecentDeployJob != nil {
		modelMap, err := dataSourceIbmSchematicsAgentAgentDataRecentDeployJobToMap(agentData.RecentDeployJob)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentRead failed: %s", err.Error()), "ibm_schematics_agent", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		recentDeployJob = append(recentDeployJob, modelMap)
	}
	if err = d.Set("recent_deploy_job", recentDeployJob); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	recentHealthJob := []map[string]interface{}{}
	if agentData.RecentHealthJob != nil {
		modelMap, err := dataSourceIbmSchematicsAgentAgentDataRecentHealthJobToMap(agentData.RecentHealthJob)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentRead failed: %s", err.Error()), "ibm_schematics_agent", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		recentHealthJob = append(recentHealthJob, modelMap)
	}
	if err = d.Set("recent_health_job", recentHealthJob); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	return nil
}

func dataSourceIbmSchematicsAgentAgentInfrastructureToMap(model *schematicsv1.AgentInfrastructure) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.InfraType != nil {
		modelMap["infra_type"] = *model.InfraType
	}
	if model.ClusterID != nil {
		modelMap["cluster_id"] = *model.ClusterID
	}
	if model.ClusterResourceGroup != nil {
		modelMap["cluster_resource_group"] = *model.ClusterResourceGroup
	}
	if model.CosInstanceName != nil {
		modelMap["cos_instance_name"] = *model.CosInstanceName
	}
	if model.CosBucketName != nil {
		modelMap["cos_bucket_name"] = *model.CosBucketName
	}
	if model.CosBucketRegion != nil {
		modelMap["cos_bucket_region"] = *model.CosBucketRegion
	}
	return modelMap, nil
}

func dataSourceIbmSchematicsAgentAgentMetadataInfoToMap(model *schematicsv1.AgentMetadataInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Value != nil {
		modelMap["value"] = model.Value
	}
	return modelMap, nil
}

func dataSourceIbmSchematicsAgentVariableDataToMap(model *schematicsv1.VariableData) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	if model.UseDefault != nil {
		modelMap["use_default"] = *model.UseDefault
	}
	if model.Metadata != nil {
		metadataMap, err := dataSourceIbmSchematicsAgentVariableMetadataToMap(model.Metadata)
		if err != nil {
			return modelMap, err
		}
		modelMap["metadata"] = []map[string]interface{}{metadataMap}
	}
	if model.Link != nil {
		modelMap["link"] = *model.Link
	}
	return modelMap, nil
}

func dataSourceIbmSchematicsAgentVariableMetadataToMap(model *schematicsv1.VariableMetadata) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Aliases != nil {
		modelMap["aliases"] = model.Aliases
	}
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	if model.CloudDataType != nil {
		modelMap["cloud_data_type"] = *model.CloudDataType
	}
	if model.DefaultValue != nil {
		modelMap["default_value"] = *model.DefaultValue
	}
	if model.LinkStatus != nil {
		modelMap["link_status"] = *model.LinkStatus
	}
	if model.Secure != nil {
		modelMap["secure"] = *model.Secure
	}
	if model.Immutable != nil {
		modelMap["immutable"] = *model.Immutable
	}
	if model.Hidden != nil {
		modelMap["hidden"] = *model.Hidden
	}
	if model.Required != nil {
		modelMap["required"] = *model.Required
	}
	if model.Options != nil {
		modelMap["options"] = model.Options
	}
	if model.MinValue != nil {
		modelMap["min_value"] = *model.MinValue
	}
	if model.MaxValue != nil {
		modelMap["max_value"] = *model.MaxValue
	}
	if model.MinLength != nil {
		modelMap["min_length"] = *model.MinLength
	}
	if model.MaxLength != nil {
		modelMap["max_length"] = *model.MaxLength
	}
	if model.Matches != nil {
		modelMap["matches"] = *model.Matches
	}
	if model.Position != nil {
		modelMap["position"] = *model.Position
	}
	if model.GroupBy != nil {
		modelMap["group_by"] = *model.GroupBy
	}
	if model.Source != nil {
		modelMap["source"] = *model.Source
	}
	return modelMap, nil
}

func dataSourceIbmSchematicsAgentAgentUserStateToMap(model *schematicsv1.AgentUserState) (map[string]interface{}, error) {
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

func dataSourceIbmSchematicsAgentAgentSystemStatusToMap(model *schematicsv1.AgentSystemStatus) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.StatusCode != nil {
		modelMap["status_code"] = *model.StatusCode
	}
	if model.StatusMessage != nil {
		modelMap["status_message"] = *model.StatusMessage
	}
	return modelMap, nil
}

func dataSourceIbmSchematicsAgentAgentKPIDataToMap(model *schematicsv1.AgentKPIData) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AvailabilityIndicator != nil {
		modelMap["availability_indicator"] = *model.AvailabilityIndicator
	}
	if model.LifecycleIndicator != nil {
		modelMap["lifecycle_indicator"] = *model.LifecycleIndicator
	}
	if model.PercentUsageIndicator != nil {
		modelMap["percent_usage_indicator"] = *model.PercentUsageIndicator
	}
	if model.ApplicationIndicators != nil {
	}
	if model.InfraIndicators != nil {
	}
	return modelMap, nil
}

func dataSourceIbmSchematicsAgentAgentDataRecentPrsJobToMap(model *schematicsv1.AgentDataRecentPrsJob) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AgentID != nil {
		modelMap["agent_id"] = *model.AgentID
	}
	if model.JobID != nil {
		modelMap["job_id"] = *model.JobID
	}
	if model.UpdatedAt != nil {
		modelMap["updated_at"] = model.UpdatedAt.String()
	}
	if model.UpdatedBy != nil {
		modelMap["updated_by"] = *model.UpdatedBy
	}
	if model.AgentVersion != nil {
		modelMap["agent_version"] = *model.AgentVersion
	}
	if model.StatusCode != nil {
		modelMap["status_code"] = *model.StatusCode
	}
	if model.StatusMessage != nil {
		modelMap["status_message"] = *model.StatusMessage
	}
	if model.LogURL != nil {
		modelMap["log_url"] = *model.LogURL
	}
	return modelMap, nil
}

func dataSourceIbmSchematicsAgentAgentDataRecentDeployJobToMap(model *schematicsv1.AgentDataRecentDeployJob) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AgentID != nil {
		modelMap["agent_id"] = *model.AgentID
	}
	if model.JobID != nil {
		modelMap["job_id"] = *model.JobID
	}
	if model.UpdatedAt != nil {
		modelMap["updated_at"] = model.UpdatedAt.String()
	}
	if model.UpdatedBy != nil {
		modelMap["updated_by"] = *model.UpdatedBy
	}
	if model.IsRedeployed != nil {
		modelMap["is_redeployed"] = *model.IsRedeployed
	}
	if model.AgentVersion != nil {
		modelMap["agent_version"] = *model.AgentVersion
	}
	if model.StatusCode != nil {
		modelMap["status_code"] = *model.StatusCode
	}
	if model.StatusMessage != nil {
		modelMap["status_message"] = *model.StatusMessage
	}
	if model.LogURL != nil {
		modelMap["log_url"] = *model.LogURL
	}
	return modelMap, nil
}

func dataSourceIbmSchematicsAgentAgentDataRecentHealthJobToMap(model *schematicsv1.AgentDataRecentHealthJob) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AgentID != nil {
		modelMap["agent_id"] = *model.AgentID
	}
	if model.JobID != nil {
		modelMap["job_id"] = *model.JobID
	}
	if model.UpdatedAt != nil {
		modelMap["updated_at"] = model.UpdatedAt.String()
	}
	if model.UpdatedBy != nil {
		modelMap["updated_by"] = *model.UpdatedBy
	}
	if model.AgentVersion != nil {
		modelMap["agent_version"] = *model.AgentVersion
	}
	if model.StatusCode != nil {
		modelMap["status_code"] = *model.StatusCode
	}
	if model.StatusMessage != nil {
		modelMap["status_message"] = *model.StatusMessage
	}
	if model.LogURL != nil {
		modelMap["log_url"] = *model.LogURL
	}
	return modelMap, nil
}
