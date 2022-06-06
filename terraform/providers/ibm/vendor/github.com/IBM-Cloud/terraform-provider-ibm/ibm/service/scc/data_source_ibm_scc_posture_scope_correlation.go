// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/scc-go-sdk/v3/posturemanagementv2"
)

func DataSourceIBMSccPostureScopeCorrelation() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMSccPostureScopeCorrelationRead,

		Schema: map[string]*schema.Schema{
			"correlation_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A correlation_Id is created when a scope is created and discovery task is triggered or when a validation is triggered on a Scope. This is used to get the status of the task(discovery or validation).",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Returns the current status of a task.",
			},
			"start_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Returns the time that task started.",
			},
			"last_heartbeat": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Returns the time that the scope was last updated. This value exists when collector is installed and running.",
			},
		},
	}
}

func dataSourceIBMSccPostureScopeCorrelationRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	postureManagementClient, err := meta.(conns.ClientSession).PostureManagementV2()
	if err != nil {
		return diag.FromErr(err)
	}

	getCorrelationIDOptions := &posturemanagementv2.GetCorrelationIDOptions{}
	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error getting userDetails %s", err))
	}

	accountID := userDetails.UserAccount
	getCorrelationIDOptions.SetAccountID(accountID)

	getCorrelationIDOptions.SetCorrelationID(d.Get("correlation_id").(string))

	scopeTaskStatus, response, err := postureManagementClient.GetCorrelationIDWithContext(context, getCorrelationIDOptions)
	if err != nil {
		log.Printf("[DEBUG] GetCorrelationIDWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetCorrelationIDWithContext failed %s\n%s", err, response))
	}

	d.SetId(dataSourceIBMSccPostureScopeCorrelationID(d))
	if err = d.Set("status", scopeTaskStatus.Status); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting status: %s", err))
	}
	if err = d.Set("start_time", scopeTaskStatus.StartTime); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting start_time: %s", err))
	}
	if err = d.Set("last_heartbeat", flex.DateTimeToString(scopeTaskStatus.LastHeartbeat)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting last_heartbeat: %s", err))
	}

	return nil
}

// dataSourceIBMScopeCorrelationID returns a reasonable ID for the list.
func dataSourceIBMSccPostureScopeCorrelationID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
