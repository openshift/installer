// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package database

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/cloud-databases-go-sdk/clouddatabasesv5"
)

func DataSourceIBMDatabaseTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMDatabaseTasksRead,

		Schema: map[string]*schema.Schema{
			"deployment_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Deployment ID.",
				ValidateFunc: validate.InvokeDataSourceValidator(
					"ibm_database_tasks",
					"deployment_id"),
			},
			"tasks": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"task_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the task.",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Human-readable description of the task.",
						},
						"status": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the task.",
						},
						"deployment_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the deployment the task is being performed on.",
						},
						"progress_percent": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Indicator as percentage of progress of the task.",
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Date and time when the task was created.",
						},
					},
				},
			},
		},
	}
}

func DataSourceIBMDatabaseTasksValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "deployment_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			Required:                   true,
			CloudDataType:              "cloud-database",
			CloudDataRange:             []string{"resolved_to:id"}})

	iBMDatabaseTasksValidator := validate.ResourceValidator{ResourceName: "ibm_database_tasks", Schema: validateSchema}
	return &iBMDatabaseTasksValidator
}

func dataSourceIBMDatabaseTasksRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cloudDatabasesClient, err := meta.(conns.ClientSession).CloudDatabasesV5()
	if err != nil {
		return diag.FromErr(err)
	}

	listDeploymentTasksOptions := &clouddatabasesv5.ListDeploymentTasksOptions{}

	listDeploymentTasksOptions.SetID(d.Get("deployment_id").(string))

	tasks, response, err := cloudDatabasesClient.ListDeploymentTasksWithContext(context, listDeploymentTasksOptions)
	if err != nil {
		return diag.FromErr(fmt.Errorf("ListDeploymentTasksWithContext failed %s\n%s", err, response))
	}

	// Use the provided filter argument and construct a new list with only the requested resource(s)
	var matchTasks []clouddatabasesv5.Task
	var deploymentID string
	var suppliedFilter bool

	if v, ok := d.GetOk("deployment_id"); ok {
		deploymentID = v.(string)
		suppliedFilter = true
		for _, data := range tasks.Tasks {
			if *data.DeploymentID == deploymentID {
				matchTasks = append(matchTasks, data)
			}
		}
	} else {
		matchTasks = tasks.Tasks
	}
	tasks.Tasks = matchTasks

	if suppliedFilter {
		if len(tasks.Tasks) == 0 {
			return diag.FromErr(fmt.Errorf("no Tasks found with deploymentID %s", deploymentID))
		}
		d.SetId(deploymentID)
	} else {
		d.SetId(DataSourceIBMDatabaseTasksID(d))
	}

	tasks2 := []map[string]interface{}{}
	if tasks.Tasks != nil {
		for _, modelItem := range tasks.Tasks {
			modelMap, err := DataSourceIBMDatabaseTasksTaskToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			tasks2 = append(tasks2, modelMap)
		}
	}
	if err = d.Set("tasks", tasks2); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting tasks %s", err))
	}

	return nil
}

// DataSourceIBMDatabaseTasksID returns a reasonable ID for the list.
func DataSourceIBMDatabaseTasksID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIBMDatabaseTasksTaskToMap(model *clouddatabasesv5.Task) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["task_id"] = *model.ID
	}
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	if model.Status != nil {
		modelMap["status"] = *model.Status
	}
	if model.DeploymentID != nil {
		modelMap["deployment_id"] = *model.DeploymentID
	}
	if model.ProgressPercent != nil {
		modelMap["progress_percent"] = *model.ProgressPercent
	}
	if model.CreatedAt != nil {
		modelMap["created_at"] = model.CreatedAt.String()
	}
	return modelMap, nil
}
