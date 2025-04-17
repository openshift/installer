// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package database

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
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
				ValidateFunc: validate.InvokeDataSourceValidator(
					"ibm_database_point_in_time_recovery",
					"deployment_id"),
			},
			"earliest_point_in_time_recovery_time": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DataSourceIBMDatabasePointInTimeRecoveryValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "deployment_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			Required:                   true,
			CloudDataType:              "cloud-database",
			CloudDataRange:             []string{"resolved_to:id"}})

	iBMDatabasePointInTimeRecoveryValidator := validate.ResourceValidator{ResourceName: "ibm_database_point_in_time_recovery", Schema: validateSchema}
	return &iBMDatabasePointInTimeRecoveryValidator
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
