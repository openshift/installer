// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package database

import (
	"context"
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/cloud-databases-go-sdk/clouddatabasesv5"
)

func DataSourceIBMDatabaseTask() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMDatabaseTaskRead,

		Schema: map[string]*schema.Schema{
			"task_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Task ID.",
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
	}
}

func dataSourceIBMDatabaseTaskRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cloudDatabasesClient, err := meta.(conns.ClientSession).CloudDatabasesV5()
	if err != nil {
		return diag.FromErr(err)
	}

	getTaskOptions := &clouddatabasesv5.GetTaskOptions{}

	getTaskOptions.SetID(d.Get("task_id").(string))

	task, response, err := cloudDatabasesClient.GetTaskWithContext(context, getTaskOptions)

	if err != nil {
		return diag.FromErr(fmt.Errorf("GetTaskWithContext failed %s\n%s", err, response))
	}

	d.SetId(*task.Task.ID)

	if err = d.Set("task_id", task.Task.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting id: %s", err))
	}

	if task.Task.Description != nil {
		if err = d.Set("description", task.Task.Description); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
		}
	}

	if task.Task.Status != nil {
		if err = d.Set("status", task.Task.Status); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting status: %s", err))
		}
	}

	if task.Task.DeploymentID != nil {
		if err = d.Set("deployment_id", task.Task.DeploymentID); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting deployment_id: %s", err))
		}
	}

	if task.Task.ProgressPercent != nil {
		if err = d.Set("progress_percent", task.Task.ProgressPercent); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting progress_percent: %s", err))
		}
	}

	if task.Task.CreatedAt != nil {
		if err = d.Set("created_at", flex.DateTimeToString(task.Task.CreatedAt)); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
		}
	}

	return nil
}
