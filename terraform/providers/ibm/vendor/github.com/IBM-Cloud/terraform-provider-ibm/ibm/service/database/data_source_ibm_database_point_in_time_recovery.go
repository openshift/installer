// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package database

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/cloud-databases-go-sdk/clouddatabasesv5"
)

func DataSourceIBMDatabasePointInTimeRecovery() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceIBMDatabasePointInTimeRecoveryRead,

		Schema: map[string]*schema.Schema{
			"deployment_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Deployment ID.",
			},
			"earliest_point_in_time_recovery_time": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DataSourceIBMDatabasePointInTimeRecoveryRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cloudDatabasesClient, err := meta.(conns.ClientSession).CloudDatabasesV5()
	if err != nil {
		return diag.FromErr(err)
	}

	getPitrDataOptions := &clouddatabasesv5.GetPitrDataOptions{}

	getPitrDataOptions.SetID(d.Get("deployment_id").(string))

	pointInTimeRecoveryData, response, err := cloudDatabasesClient.GetPitrDataWithContext(context, getPitrDataOptions)
	if err != nil {
		log.Printf("[DEBUG] GetPitrDataWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetPitrDataWithContext failed %s\n%s", err, response))
	}

	d.SetId(d.Get("deployment_id").(string))

	if pointInTimeRecoveryData.PointInTimeRecoveryData.EarliestPointInTimeRecoveryTime != nil {
		pitr := pointInTimeRecoveryData.PointInTimeRecoveryData.EarliestPointInTimeRecoveryTime
		if err = d.Set("earliest_point_in_time_recovery_time", pitr); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting earliest_point_in_time_recovery_time: %s", err))
		}
	}

	return nil
}
