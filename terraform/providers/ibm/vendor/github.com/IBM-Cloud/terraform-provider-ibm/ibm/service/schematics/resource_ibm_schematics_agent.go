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
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/schematics-go-sdk/schematicsv1"
)

func ResourceIbmSchematicsAgent() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmSchematicsAgentCreate,
		ReadContext:   resourceIbmSchematicsAgentRead,
		UpdateContext: resourceIbmSchematicsAgentUpdate,
		DeleteContext: resourceIbmSchematicsAgentDelete,
		Importer:      &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the agent (must be unique, for an account).",
			},
			"resource_group": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The resource-group name for the agent.  By default, agent will be registered in Default Resource Group.",
			},
			"version": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Agent version.",
			},
			"schematics_location": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "List of locations supported by IBM Cloud Schematics service.  While creating your workspace or action, choose the right region, since it cannot be changed.  Note, this does not limit the location of the IBM Cloud resources, provisioned using Schematics.",
			},
			"agent_location": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The location where agent is deployed in the user environment.",
			},
			"agent_infrastructure": &schema.Schema{
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Description: "The infrastructure parameters used by the agent.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"infra_type": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Type of target agent infrastructure.",
						},
						"cluster_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The cluster ID where agent services will be running.",
						},
						"cluster_resource_group": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The resource group of the cluster (is it required?).",
						},
						"cos_instance_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The COS instance name to store the agent logs.",
						},
						"cos_bucket_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The COS bucket name used to store the logs.",
						},
						"cos_bucket_region": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The COS bucket region.",
						},
					},
				},
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Agent description.",
			},
			"tags": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Tags for the agent.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"agent_metadata": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The metadata of an agent.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the metadata.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Value of the metadata name.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"agent_inputs": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Additional input variables for the agent.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the variable. For example, `name = \"inventory username\"`.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The value for the variable or reference to the value. For example, `value = \"<provide your ssh_key_value with \n>\"`. **Note** The SSH key should contain `\n` at the end of the key details in case of command line or API calls.",
						},
						"use_default": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "True, will ignore the data in the value attribute, instead the data in metadata.default_value will be used.",
						},
						"metadata": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "An user editable metadata for the variables.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Type of the variable.",
									},
									"aliases": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "The list of aliases for the variable name.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"description": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The description of the meta data.",
									},
									"cloud_data_type": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Cloud data type of the variable. eg. resource_group_id, region, vpc_id.",
									},
									"default_value": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Default value for the variable only if the override value is not specified.",
									},
									"link_status": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The status of the link.",
									},
									"secure": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Is the variable secure or sensitive ?.",
									},
									"immutable": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Is the variable readonly ?.",
									},
									"hidden": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "If **true**, the variable is not displayed on UI or Command line.",
									},
									"required": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "If the variable required?.",
									},
									"options": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "The list of possible values for this variable.  If type is **integer** or **date**, then the array of string is  converted to array of integers or date during the runtime.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"min_value": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The minimum value of the variable. Applicable for the integer type.",
									},
									"max_value": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The maximum value of the variable. Applicable for the integer type.",
									},
									"min_length": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The minimum length of the variable value. Applicable for the string type.",
									},
									"max_length": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The maximum length of the variable value. Applicable for the string type.",
									},
									"matches": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The regex for the variable value.",
									},
									"position": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The relative position of this variable in a list.",
									},
									"group_by": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The display name of the group this variable belongs to.",
									},
									"source": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
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
			"run_destroy_resources": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Argument which helps to run destroy resources job. Increment the value to destroy resources associated with agent deployment.",
			},
			"user_state": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "User defined status of the agent.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
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
							Elem:        &schema.Schema{Type: schema.TypeMap},
						},
						"infra_indicators": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Agent infrastructure key performance indicators.",
							Elem:        &schema.Schema{Type: schema.TypeMap},
						},
					},
				},
			},
			"agent_crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The agent crn, obtained from the Schematics agent deployment configuration.",
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
							Optional:    true,
							Computed:    true,
							Description: "Agent Status.",
						},
						"status_message": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The agent status message.",
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
							Optional:    true,
							Computed:    true,
							Description: "Id of the agent.",
						},
						"job_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
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
							Optional:    true,
							Computed:    true,
							Description: "Final result of the pre-requisite scanner job.",
						},
						"status_message": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The outcome of the pre-requisite scanner job, in a formatted log string.",
						},
						"log_url": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
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
							Optional:    true,
							Computed:    true,
							Description: "Id of the agent.",
						},
						"job_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
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
							Optional:    true,
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
							Optional:    true,
							Computed:    true,
							Description: "Final result of the agent deployment job.",
						},
						"status_message": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The outcome of the agent deployment job, in a formatted log string.",
						},
						"log_url": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
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
							Optional:    true,
							Computed:    true,
							Description: "Id of the agent.",
						},
						"job_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
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
							Optional:    true,
							Computed:    true,
							Description: "Final result of the health-check job.",
						},
						"status_message": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The outcome of the health-check job, in a formatted log string.",
						},
						"log_url": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "URL to the full health-check job logs.",
						},
					},
				},
			},
		},
	}
}

func resourceIbmSchematicsAgentCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentCreate schematicsClient initialization failed: %s", err.Error()), "ibm_schematics_agent", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createAgentDataOptions := &schematicsv1.CreateAgentDataOptions{}

	createAgentDataOptions.SetName(d.Get("name").(string))
	createAgentDataOptions.SetResourceGroup(d.Get("resource_group").(string))
	createAgentDataOptions.SetVersion(d.Get("version").(string))
	createAgentDataOptions.SetSchematicsLocation(d.Get("schematics_location").(string))
	createAgentDataOptions.SetAgentLocation(d.Get("agent_location").(string))
	agentInfrastructureModel, err := resourceIbmSchematicsAgentMapToAgentInfrastructure(d.Get("agent_infrastructure.0").(map[string]interface{}))
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentCreate failed: %s", err.Error()), "ibm_schematics_agent", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	createAgentDataOptions.SetAgentInfrastructure(agentInfrastructureModel)
	if _, ok := d.GetOk("description"); ok {
		createAgentDataOptions.SetDescription(d.Get("description").(string))
	}
	if _, ok := d.GetOk("tags"); ok {
		createAgentDataOptions.SetTags(flex.ExpandStringList(d.Get("tags").([]interface{})))
	}
	if _, ok := d.GetOk("agent_metadata"); ok {
		var agentMetadata []schematicsv1.AgentMetadataInfo
		for _, e := range d.Get("agent_metadata").([]interface{}) {
			value := e.(map[string]interface{})
			agentMetadataItem, err := resourceIbmSchematicsAgentMapToAgentMetadataInfo(value)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentCreate failed: %s", err.Error()), "ibm_schematics_agent", "create")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			agentMetadata = append(agentMetadata, *agentMetadataItem)
		}
		createAgentDataOptions.SetAgentMetadata(agentMetadata)
	}
	if _, ok := d.GetOk("agent_inputs"); ok {
		var agentInputs []schematicsv1.VariableData
		for _, e := range d.Get("agent_inputs").([]interface{}) {
			value := e.(map[string]interface{})
			agentInputsItem, err := resourceIbmSchematicsAgentMapToVariableData(value)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentCreate failed: %s", err.Error()), "ibm_schematics_agent", "create")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			agentInputs = append(agentInputs, *agentInputsItem)
		}
		createAgentDataOptions.SetAgentInputs(agentInputs)
	}
	if _, ok := d.GetOk("user_state"); ok {
		userStateModel, err := resourceIbmSchematicsAgentMapToAgentUserState(d.Get("user_state.0").(map[string]interface{}))
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentCreate failed: %s", err.Error()), "ibm_schematics_agent", "create")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		createAgentDataOptions.SetUserState(userStateModel)
	}

	agentData, response, err := schematicsClient.CreateAgentDataWithContext(context, createAgentDataOptions)
	if err != nil {

		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentCreate CreateAgentDataWithContext failed with error: %s and response:\n%s", err, response), "ibm_schematics_agent", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*agentData.ID)

	return resourceIbmSchematicsAgentRead(context, d, meta)
}

func resourceIbmSchematicsAgentRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentRead schematicsClient initialization failed: %s", err.Error()), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getAgentDataOptions := &schematicsv1.GetAgentDataOptions{
		Profile: core.StringPtr("detailed"),
	}

	getAgentDataOptions.SetAgentID(d.Id())

	agentData, response, err := schematicsClient.GetAgentDataWithContext(context, getAgentDataOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentRead GetAgentDataWithContext failed with error: %s and response:\n%s", err, response), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("name", agentData.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("resource_group", agentData.ResourceGroup); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("version", agentData.Version); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("schematics_location", agentData.SchematicsLocation); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("agent_location", agentData.AgentLocation); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	agentInfrastructureMap, err := resourceIbmSchematicsAgentAgentInfrastructureToMap(agentData.AgentInfrastructure)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentRead failed: %s", err.Error()), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("agent_infrastructure", []map[string]interface{}{agentInfrastructureMap}); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("description", agentData.Description); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if agentData.Tags != nil {
		if err = d.Set("tags", agentData.Tags); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	agentMetadata := []map[string]interface{}{}
	if agentData.AgentMetadata != nil {
		for _, agentMetadataItem := range agentData.AgentMetadata {
			agentMetadataItemMap, err := resourceIbmSchematicsAgentAgentMetadataInfoToMap(&agentMetadataItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentRead failed: %s", err.Error()), "ibm_schematics_agent", "read")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			agentMetadata = append(agentMetadata, agentMetadataItemMap)
		}
	}
	if err = d.Set("agent_metadata", agentMetadata); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	agentInputs := []map[string]interface{}{}
	if agentData.AgentInputs != nil {
		for _, agentInputsItem := range agentData.AgentInputs {
			agentInputsItemMap, err := resourceIbmSchematicsAgentVariableDataToMap(&agentInputsItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentRead failed: %s", err.Error()), "ibm_schematics_agent", "read")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			agentInputs = append(agentInputs, agentInputsItemMap)
		}
	}
	if err = d.Set("agent_inputs", agentInputs); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if agentData.UserState != nil {
		userStateMap, err := resourceIbmSchematicsAgentAgentUserStateToMap(agentData.UserState)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentRead failed: %s", err.Error()), "ibm_schematics_agent", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("user_state", []map[string]interface{}{userStateMap}); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	if agentData.AgentKpi != nil {
		agentKpiMap, err := resourceIbmSchematicsAgentAgentKPIDataToMap(agentData.AgentKpi)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentRead failed: %s", err.Error()), "ibm_schematics_agent", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("agent_kpi", []map[string]interface{}{agentKpiMap}); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("agent_crn", agentData.AgentCrn); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("created_at", flex.DateTimeToString(agentData.CreatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("creation_by", agentData.CreationBy); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("updated_at", flex.DateTimeToString(agentData.UpdatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("updated_by", agentData.UpdatedBy); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if agentData.SystemState != nil {
		systemStateMap, err := resourceIbmSchematicsAgentAgentSystemStatusToMap(agentData.SystemState)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentRead failed: %s", err.Error()), "ibm_schematics_agent", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("system_state", []map[string]interface{}{systemStateMap}); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	if agentData.RecentPrsJob != nil {
		recentPrsJobMap, err := resourceIbmSchematicsAgentAgentDataRecentPrsJobToMap(agentData.RecentPrsJob)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentRead failed: %s", err.Error()), "ibm_schematics_agent", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("recent_prs_job", []map[string]interface{}{recentPrsJobMap}); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	if agentData.RecentDeployJob != nil {
		recentDeployJobMap, err := resourceIbmSchematicsAgentAgentDataRecentDeployJobToMap(agentData.RecentDeployJob)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentRead failed: %s", err.Error()), "ibm_schematics_agent", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("recent_deploy_job", []map[string]interface{}{recentDeployJobMap}); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	if agentData.RecentHealthJob != nil {
		recentHealthJobMap, err := resourceIbmSchematicsAgentAgentDataRecentHealthJobToMap(agentData.RecentHealthJob)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentRead failed: %s", err.Error()), "ibm_schematics_agent", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("recent_health_job", []map[string]interface{}{recentHealthJobMap}); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentRead failed with error: %s", err), "ibm_schematics_agent", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return nil
}
func isWaitForAgentDestroyResources(context context.Context, schematicsClient *schematicsv1.SchematicsV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for agent (%s) resources to be destroyed.", id)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", agentProvisioningStatusCodeJobInProgress, agentProvisioningStatusCodeJobPending, agentProvisioningStatusCodeJobReadyToExecute, agentProvisioningStatusCodeJobStopInProgress},
		Target:     []string{agentProvisioningStatusCodeJobFinished, agentProvisioningStatusCodeJobFailed, agentProvisioningStatusCodeJobCancelled, agentProvisioningStatusCodeJobStopped, ""},
		Refresh:    agentDestroyRefreshFunc(schematicsClient, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	return stateConf.WaitForStateContext(context)
}
func agentDestroyRefreshFunc(schematicsClient *schematicsv1.SchematicsV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getAgentDataOptions := &schematicsv1.GetAgentDataOptions{
			AgentID: core.StringPtr(id),
			Profile: core.StringPtr("detailed"),
		}

		agent, response, err := schematicsClient.GetAgentData(getAgentDataOptions)
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] Error Getting Agent: %s\n%s", err, response)
		}
		if agent.RecentDestroyJob.StatusCode != nil {
			return agent, *agent.RecentDestroyJob.StatusCode, nil
		}
		return agent, agentProvisioningStatusCodeJobPending, nil
	}
}
func resourceIbmSchematicsAgentUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentUpdate schematicsClient initialization failed: %s", err.Error()), "ibm_schematics_agent", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateAgentDataOptions := &schematicsv1.UpdateAgentDataOptions{}
	session, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentUpdate bluemixClient initialization failed: %s", err.Error()), "ibm_schematics_agent", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	iamAccessToken := session.Config.IAMAccessToken
	iamRefreshToken := session.Config.IAMRefreshToken
	ff := map[string]string{
		"Authorization": iamAccessToken,
		"refresh_token": iamRefreshToken,
	}
	updateAgentDataOptions.Headers = ff

	updateAgentDataOptions.SetAgentID(d.Id())

	hasChange := false

	getAgentDataOptions := &schematicsv1.GetAgentDataOptions{
		Profile: core.StringPtr("detailed"),
	}

	getAgentDataOptions.SetAgentID(d.Id())

	agentData, response, err := schematicsClient.GetAgentDataWithContext(context, getAgentDataOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentUpdate GetAgentDataWithContext failed with error: %s and response:\n%s", err, response), "ibm_schematics_agent", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if agentData.Name != nil {
		updateAgentDataOptions.Name = agentData.Name
	}
	if agentData.ResourceGroup != nil {
		updateAgentDataOptions.ResourceGroup = agentData.ResourceGroup
	}
	if agentData.Version != nil {
		updateAgentDataOptions.Version = agentData.Version
	}
	if agentData.SchematicsLocation != nil {
		updateAgentDataOptions.SchematicsLocation = agentData.SchematicsLocation
	}
	if agentData.AgentLocation != nil {
		updateAgentDataOptions.AgentLocation = agentData.AgentLocation
	}
	if agentData.AgentInfrastructure != nil {
		updateAgentDataOptions.AgentInfrastructure = agentData.AgentInfrastructure
	}
	if agentData.Description != nil {
		updateAgentDataOptions.Description = agentData.Description
	}
	if agentData.Tags != nil {
		updateAgentDataOptions.Tags = agentData.Tags
	}
	if agentData.AgentMetadata != nil {
		updateAgentDataOptions.AgentMetadata = agentData.AgentMetadata
	}
	if agentData.AgentInputs != nil {
		updateAgentDataOptions.AgentInputs = agentData.AgentInputs
	}
	if agentData.UserState != nil {
		updateAgentDataOptions.UserState = agentData.UserState
	}
	if agentData.AgentKpi != nil {
		updateAgentDataOptions.AgentKpi = agentData.AgentKpi
	}

	if d.HasChange("schematics_location") || d.HasChange("agent_location") {
		updateAgentDataOptions.SetSchematicsLocation(d.Get("schematics_location").(string))
		updateAgentDataOptions.SetAgentLocation(d.Get("agent_location").(string))
		hasChange = true
	}
	if d.HasChange("name") {
		updateAgentDataOptions.SetName(d.Get("name").(string))
		hasChange = true
	}
	if d.HasChange("version") {
		updateAgentDataOptions.SetVersion(d.Get("version").(string))
		hasChange = true
	}
	if d.HasChange("agent_infrastructure") {
		agentInfrastructure, err := resourceIbmSchematicsAgentMapToAgentInfrastructure(d.Get("agent_infrastructure.0").(map[string]interface{}))
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentUpdate failed: %s", err.Error()), "ibm_schematics_agent", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		updateAgentDataOptions.SetAgentInfrastructure(agentInfrastructure)
		hasChange = true
	}
	if d.HasChange("description") {
		updateAgentDataOptions.SetDescription(d.Get("description").(string))
		hasChange = true
	}
	if d.HasChange("tags") {
		updateAgentDataOptions.SetTags(flex.ExpandStringList(d.Get("tags").([]interface{})))
		hasChange = true
	}
	if d.HasChange("agent_metadata") {
		var agentMetadata []schematicsv1.AgentMetadataInfo
		for _, e := range d.Get("agent_metadata").([]interface{}) {
			value := e.(map[string]interface{})
			agentMetadataItem, err := resourceIbmSchematicsAgentMapToAgentMetadataInfo(value)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentUpdate failed: %s", err.Error()), "ibm_schematics_agent", "update")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			agentMetadata = append(agentMetadata, *agentMetadataItem)
		}
		updateAgentDataOptions.SetAgentMetadata(agentMetadata)
		hasChange = true
	}
	if d.HasChange("run_destroy_resources") {
		deleteAgentResourcesOptions := &schematicsv1.DeleteAgentResourcesOptions{}
		deleteAgentResourcesOptions.Headers = ff

		deleteAgentResourcesOptions.SetAgentID(d.Id())
		deleteAgentResourcesOptions.SetRefreshToken(iamRefreshToken)

		response, err := schematicsClient.DeleteAgentResourcesWithContext(context, deleteAgentResourcesOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentUpdate DeleteAgentResourcesWithContext failed with error: %s and response:\n%s", err, response), "ibm_schematics_agent", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		} else {
			_, err = isWaitForAgentDestroyResources(context, schematicsClient, *deleteAgentResourcesOptions.AgentID, d.Timeout(schema.TimeoutUpdate))
			if err != nil {

			}
		}
	}
	if hasChange {
		_, response, err := schematicsClient.UpdateAgentDataWithContext(context, updateAgentDataOptions)
		if err != nil {

			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentUpdate failed with error: %s and response:\n%s", err, response), "ibm_schematics_agent", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIbmSchematicsAgentRead(context, d, meta)
}

func resourceIbmSchematicsAgentDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentDelete schematicsClient initialization failed: %s", err.Error()), "ibm_schematics_agent", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteAgentDataOptions := &schematicsv1.DeleteAgentDataOptions{
		Force: core.BoolPtr(true),
	}
	session, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentDelete bluemixClient initialization failed: %s", err.Error()), "ibm_schematics_agent", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	iamAccessToken := session.Config.IAMAccessToken
	iamRefreshToken := session.Config.IAMRefreshToken
	ff := map[string]string{
		"Authorization": iamAccessToken,
		"refresh_token": iamRefreshToken,
	}
	deleteAgentDataOptions.Headers = ff

	deleteAgentDataOptions.SetAgentID(d.Id())

	// first try destroying resources associated with agent deploy and then delete the agent

	deleteAgentResourcesOptions := &schematicsv1.DeleteAgentResourcesOptions{}
	deleteAgentResourcesOptions.Headers = ff

	deleteAgentResourcesOptions.SetAgentID(d.Id())
	deleteAgentResourcesOptions.SetRefreshToken(iamRefreshToken)

	response, err := schematicsClient.DeleteAgentResourcesWithContext(context, deleteAgentResourcesOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentDelete DeleteAgentResourcesWithContext failed with error: %s and response:\n%s", err, response), "ibm_schematics_agent", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
	} else {
		_, err = isWaitForAgentDestroyResources(context, schematicsClient, *deleteAgentResourcesOptions.AgentID, d.Timeout(schema.TimeoutDelete))
		if err != nil {

		}
	}

	// After deploy associated resources are destroyed, now attempt to delete the agent

	deleteresponse, err := schematicsClient.DeleteAgentDataWithContext(context, deleteAgentDataOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsAgentDelete DeleteAgentDataWithContext failed with error: %s and response:\n%s", err, deleteresponse), "ibm_schematics_agent", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func resourceIbmSchematicsAgentMapToAgentInfrastructure(modelMap map[string]interface{}) (*schematicsv1.AgentInfrastructure, error) {
	model := &schematicsv1.AgentInfrastructure{}
	if modelMap["infra_type"] != nil && modelMap["infra_type"].(string) != "" {
		model.InfraType = core.StringPtr(modelMap["infra_type"].(string))
	}
	if modelMap["cluster_id"] != nil && modelMap["cluster_id"].(string) != "" {
		model.ClusterID = core.StringPtr(modelMap["cluster_id"].(string))
	}
	if modelMap["cluster_resource_group"] != nil && modelMap["cluster_resource_group"].(string) != "" {
		model.ClusterResourceGroup = core.StringPtr(modelMap["cluster_resource_group"].(string))
	}
	if modelMap["cos_instance_name"] != nil && modelMap["cos_instance_name"].(string) != "" {
		model.CosInstanceName = core.StringPtr(modelMap["cos_instance_name"].(string))
	}
	if modelMap["cos_bucket_name"] != nil && modelMap["cos_bucket_name"].(string) != "" {
		model.CosBucketName = core.StringPtr(modelMap["cos_bucket_name"].(string))
	}
	if modelMap["cos_bucket_region"] != nil && modelMap["cos_bucket_region"].(string) != "" {
		model.CosBucketRegion = core.StringPtr(modelMap["cos_bucket_region"].(string))
	}
	return model, nil
}

func resourceIbmSchematicsAgentMapToAgentMetadataInfo(modelMap map[string]interface{}) (*schematicsv1.AgentMetadataInfo, error) {
	model := &schematicsv1.AgentMetadataInfo{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["value"] != nil {
		value := []string{}
		for _, valueItem := range modelMap["value"].([]interface{}) {
			value = append(value, valueItem.(string))
		}
		model.Value = value
	}
	return model, nil
}

func resourceIbmSchematicsAgentMapToVariableData(modelMap map[string]interface{}) (*schematicsv1.VariableData, error) {
	model := &schematicsv1.VariableData{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["value"] != nil && modelMap["value"].(string) != "" {
		model.Value = core.StringPtr(modelMap["value"].(string))
	}
	if modelMap["use_default"] != nil {
		model.UseDefault = core.BoolPtr(modelMap["use_default"].(bool))
	}
	if modelMap["metadata"] != nil && len(modelMap["metadata"].([]interface{})) > 0 {
		MetadataModel, err := resourceIbmSchematicsAgentMapToVariableMetadata(modelMap["metadata"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Metadata = MetadataModel
	}
	if modelMap["link"] != nil && modelMap["link"].(string) != "" {
		model.Link = core.StringPtr(modelMap["link"].(string))
	}
	return model, nil
}

func resourceIbmSchematicsAgentMapToVariableMetadata(modelMap map[string]interface{}) (*schematicsv1.VariableMetadata, error) {
	model := &schematicsv1.VariableMetadata{}
	if modelMap["type"] != nil && modelMap["type"].(string) != "" {
		model.Type = core.StringPtr(modelMap["type"].(string))
	}
	if modelMap["aliases"] != nil {
		aliases := []string{}
		for _, aliasesItem := range modelMap["aliases"].([]interface{}) {
			aliases = append(aliases, aliasesItem.(string))
		}
		model.Aliases = aliases
	}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["cloud_data_type"] != nil && modelMap["cloud_data_type"].(string) != "" {
		model.CloudDataType = core.StringPtr(modelMap["cloud_data_type"].(string))
	}
	if modelMap["default_value"] != nil && modelMap["default_value"].(string) != "" {
		model.DefaultValue = core.StringPtr(modelMap["default_value"].(string))
	}
	if modelMap["link_status"] != nil && modelMap["link_status"].(string) != "" {
		model.LinkStatus = core.StringPtr(modelMap["link_status"].(string))
	}
	if modelMap["secure"] != nil {
		model.Secure = core.BoolPtr(modelMap["secure"].(bool))
	}
	if modelMap["immutable"] != nil {
		model.Immutable = core.BoolPtr(modelMap["immutable"].(bool))
	}
	if modelMap["hidden"] != nil {
		model.Hidden = core.BoolPtr(modelMap["hidden"].(bool))
	}
	if modelMap["required"] != nil {
		model.Required = core.BoolPtr(modelMap["required"].(bool))
	}
	if modelMap["options"] != nil {
		options := []string{}
		for _, optionsItem := range modelMap["options"].([]interface{}) {
			options = append(options, optionsItem.(string))
		}
		model.Options = options
	}
	if modelMap["min_value"] != nil {
		model.MinValue = core.Int64Ptr(int64(modelMap["min_value"].(int)))
	}
	if modelMap["max_value"] != nil {
		model.MaxValue = core.Int64Ptr(int64(modelMap["max_value"].(int)))
	}
	if modelMap["min_length"] != nil {
		model.MinLength = core.Int64Ptr(int64(modelMap["min_length"].(int)))
	}
	if modelMap["max_length"] != nil {
		model.MaxLength = core.Int64Ptr(int64(modelMap["max_length"].(int)))
	}
	if modelMap["matches"] != nil && modelMap["matches"].(string) != "" {
		model.Matches = core.StringPtr(modelMap["matches"].(string))
	}
	if modelMap["position"] != nil {
		model.Position = core.Int64Ptr(int64(modelMap["position"].(int)))
	}
	if modelMap["group_by"] != nil && modelMap["group_by"].(string) != "" {
		model.GroupBy = core.StringPtr(modelMap["group_by"].(string))
	}
	if modelMap["source"] != nil && modelMap["source"].(string) != "" {
		model.Source = core.StringPtr(modelMap["source"].(string))
	}
	return model, nil
}

func resourceIbmSchematicsAgentMapToAgentUserState(modelMap map[string]interface{}) (*schematicsv1.AgentUserState, error) {
	model := &schematicsv1.AgentUserState{}
	if modelMap["state"] != nil && modelMap["state"].(string) != "" {
		model.State = core.StringPtr(modelMap["state"].(string))
	}
	if modelMap["set_by"] != nil && modelMap["set_by"].(string) != "" {
		model.SetBy = core.StringPtr(modelMap["set_by"].(string))
	}
	if modelMap["set_at"] != nil {

	}
	return model, nil
}

func resourceIbmSchematicsAgentAgentInfrastructureToMap(model *schematicsv1.AgentInfrastructure) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.InfraType != nil {
		modelMap["infra_type"] = model.InfraType
	}
	if model.ClusterID != nil {
		modelMap["cluster_id"] = model.ClusterID
	}
	if model.ClusterResourceGroup != nil {
		modelMap["cluster_resource_group"] = model.ClusterResourceGroup
	}
	if model.CosInstanceName != nil {
		modelMap["cos_instance_name"] = model.CosInstanceName
	}
	if model.CosBucketName != nil {
		modelMap["cos_bucket_name"] = model.CosBucketName
	}
	if model.CosBucketRegion != nil {
		modelMap["cos_bucket_region"] = model.CosBucketRegion
	}
	return modelMap, nil
}

func resourceIbmSchematicsAgentAgentMetadataInfoToMap(model *schematicsv1.AgentMetadataInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Value != nil {
		modelMap["value"] = model.Value
	}
	return modelMap, nil
}

func resourceIbmSchematicsAgentVariableDataToMap(model *schematicsv1.VariableData) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Value != nil {
		modelMap["value"] = model.Value
	}
	if model.UseDefault != nil {
		modelMap["use_default"] = model.UseDefault
	}
	if model.Metadata != nil {
		metadataMap, err := resourceIbmSchematicsAgentVariableMetadataToMap(model.Metadata)
		if err != nil {
			return modelMap, err
		}
		modelMap["metadata"] = []map[string]interface{}{metadataMap}
	}
	if model.Link != nil {
		modelMap["link"] = model.Link
	}
	return modelMap, nil
}

func resourceIbmSchematicsAgentVariableMetadataToMap(model *schematicsv1.VariableMetadata) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	if model.Aliases != nil {
		modelMap["aliases"] = model.Aliases
	}
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.CloudDataType != nil {
		modelMap["cloud_data_type"] = model.CloudDataType
	}
	if model.DefaultValue != nil {
		modelMap["default_value"] = model.DefaultValue
	}
	if model.LinkStatus != nil {
		modelMap["link_status"] = model.LinkStatus
	}
	if model.Secure != nil {
		modelMap["secure"] = model.Secure
	}
	if model.Immutable != nil {
		modelMap["immutable"] = model.Immutable
	}
	if model.Hidden != nil {
		modelMap["hidden"] = model.Hidden
	}
	if model.Required != nil {
		modelMap["required"] = model.Required
	}
	if model.Options != nil {
		modelMap["options"] = model.Options
	}
	if model.MinValue != nil {
		modelMap["min_value"] = flex.IntValue(model.MinValue)
	}
	if model.MaxValue != nil {
		modelMap["max_value"] = flex.IntValue(model.MaxValue)
	}
	if model.MinLength != nil {
		modelMap["min_length"] = flex.IntValue(model.MinLength)
	}
	if model.MaxLength != nil {
		modelMap["max_length"] = flex.IntValue(model.MaxLength)
	}
	if model.Matches != nil {
		modelMap["matches"] = model.Matches
	}
	if model.Position != nil {
		modelMap["position"] = flex.IntValue(model.Position)
	}
	if model.GroupBy != nil {
		modelMap["group_by"] = model.GroupBy
	}
	if model.Source != nil {
		modelMap["source"] = model.Source
	}
	return modelMap, nil
}

func resourceIbmSchematicsAgentAgentUserStateToMap(model *schematicsv1.AgentUserState) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.State != nil {
		modelMap["state"] = model.State
	}
	if model.SetBy != nil {
		modelMap["set_by"] = model.SetBy
	}
	if model.SetAt != nil {
		modelMap["set_at"] = model.SetAt.String()
	}
	return modelMap, nil
}

func resourceIbmSchematicsAgentAgentKPIDataToMap(model *schematicsv1.AgentKPIData) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AvailabilityIndicator != nil {
		modelMap["availability_indicator"] = model.AvailabilityIndicator
	}
	if model.LifecycleIndicator != nil {
		modelMap["lifecycle_indicator"] = model.LifecycleIndicator
	}
	if model.PercentUsageIndicator != nil {
		modelMap["percent_usage_indicator"] = model.PercentUsageIndicator
	}
	if model.ApplicationIndicators != nil {
		applicationIndicators := []map[string]interface{}{}
		for _, applicationIndicatorsItem := range model.ApplicationIndicators {
			applicationIndicators = append(applicationIndicators, applicationIndicatorsItem.(map[string]interface{}))
		}
		modelMap["application_indicators"] = applicationIndicators
	}
	if model.InfraIndicators != nil {
		infraIndicators := []map[string]interface{}{}
		for _, infraIndicatorsItem := range model.InfraIndicators {
			infraIndicators = append(infraIndicators, infraIndicatorsItem.(map[string]interface{}))
		}
		modelMap["infra_indicators"] = infraIndicators
	}
	return modelMap, nil
}

func resourceIbmSchematicsAgentAgentSystemStatusToMap(model *schematicsv1.AgentSystemStatus) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.StatusCode != nil {
		modelMap["status_code"] = model.StatusCode
	}
	if model.StatusMessage != nil {
		modelMap["status_message"] = model.StatusMessage
	}
	return modelMap, nil
}

func resourceIbmSchematicsAgentAgentDataRecentPrsJobToMap(model *schematicsv1.AgentDataRecentPrsJob) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AgentID != nil {
		modelMap["agent_id"] = model.AgentID
	}
	if model.JobID != nil {
		modelMap["job_id"] = model.JobID
	}
	if model.UpdatedAt != nil {
		modelMap["updated_at"] = model.UpdatedAt.String()
	}
	if model.UpdatedBy != nil {
		modelMap["updated_by"] = model.UpdatedBy
	}
	if model.AgentVersion != nil {
		modelMap["agent_version"] = model.AgentVersion
	}
	if model.StatusCode != nil {
		modelMap["status_code"] = model.StatusCode
	}
	if model.StatusMessage != nil {
		modelMap["status_message"] = model.StatusMessage
	}
	if model.LogURL != nil {
		modelMap["log_url"] = model.LogURL
	}
	return modelMap, nil
}

func resourceIbmSchematicsAgentAgentDataRecentDeployJobToMap(model *schematicsv1.AgentDataRecentDeployJob) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AgentID != nil {
		modelMap["agent_id"] = model.AgentID
	}
	if model.JobID != nil {
		modelMap["job_id"] = model.JobID
	}
	if model.UpdatedAt != nil {
		modelMap["updated_at"] = model.UpdatedAt.String()
	}
	if model.UpdatedBy != nil {
		modelMap["updated_by"] = model.UpdatedBy
	}
	if model.IsRedeployed != nil {
		modelMap["is_redeployed"] = model.IsRedeployed
	}
	if model.AgentVersion != nil {
		modelMap["agent_version"] = model.AgentVersion
	}
	if model.StatusCode != nil {
		modelMap["status_code"] = model.StatusCode
	}
	if model.StatusMessage != nil {
		modelMap["status_message"] = model.StatusMessage
	}
	if model.LogURL != nil {
		modelMap["log_url"] = model.LogURL
	}
	return modelMap, nil
}

func resourceIbmSchematicsAgentAgentDataRecentHealthJobToMap(model *schematicsv1.AgentDataRecentHealthJob) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AgentID != nil {
		modelMap["agent_id"] = model.AgentID
	}
	if model.JobID != nil {
		modelMap["job_id"] = model.JobID
	}
	if model.UpdatedAt != nil {
		modelMap["updated_at"] = model.UpdatedAt.String()
	}
	if model.UpdatedBy != nil {
		modelMap["updated_by"] = model.UpdatedBy
	}
	if model.AgentVersion != nil {
		modelMap["agent_version"] = model.AgentVersion
	}
	if model.StatusCode != nil {
		modelMap["status_code"] = model.StatusCode
	}
	if model.StatusMessage != nil {
		modelMap["status_message"] = model.StatusMessage
	}
	if model.LogURL != nil {
		modelMap["log_url"] = model.LogURL
	}
	return modelMap, nil
}
