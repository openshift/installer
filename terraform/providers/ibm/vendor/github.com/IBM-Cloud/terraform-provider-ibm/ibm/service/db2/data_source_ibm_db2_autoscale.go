// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.96.0-d6dec9d7-20241008-212902
 */

package db2

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/cloud-db2-go-sdk/db2saasv1"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
)

func DataSourceIbmDb2Autoscale() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmDb2AutoscaleRead,

		Schema: map[string]*schema.Schema{
			"deployment_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Encoded CRN deployment id.",
			},
			"auto_scaling_allow_plan_limit": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates the maximum number of scaling actions that are allowed within a specified time period.",
			},
			"auto_scaling_enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if automatic scaling is enabled or not.",
			},
			"auto_scaling_max_storage": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum limit for automatically increasing storage capacity to handle growing data needs.",
			},
			"auto_scaling_over_time_period": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Defines the time period over which auto-scaling adjustments are monitored and applied.",
			},
			"auto_scaling_pause_limit": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Specifies the duration to pause auto-scaling actions after a scaling event has occurred.",
			},
			"auto_scaling_threshold": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Specifies the resource utilization level that triggers an auto-scaling.",
			},
			"storage_unit": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the unit of measurement for storage capacity.",
			},
			"storage_utilization_percentage": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Represents the percentage of total storage capacity currently in use.",
			},
			"support_auto_scaling": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether a system or service can automatically adjust resources based on demand.",
			},
		},
	}
}

func dataSourceIbmDb2AutoscaleRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	db2saasClient, err := meta.(conns.ClientSession).Db2saasV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_db2_autoscale", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getDb2SaasAutoscaleOptions := &db2saasv1.GetDb2SaasAutoscaleOptions{}

	getDb2SaasAutoscaleOptions.SetXDbProfile(d.Get("deployment_id").(string))

	successAutoScaling, _, err := db2saasClient.GetDb2SaasAutoscaleWithContext(context, getDb2SaasAutoscaleOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetDb2SaasAutoscaleWithContext failed: %s", err.Error()), "(Data) ibm_db2_autoscale", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIbmDb2AutoscaleID(d))

	if err = d.Set("auto_scaling_allow_plan_limit", successAutoScaling.AutoScalingAllowPlanLimit); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting auto_scaling_allow_plan_limit: %s", err), "(Data) ibm_db2_autoscale", "read", "set-auto_scaling_allow_plan_limit").GetDiag()
	}

	if err = d.Set("auto_scaling_enabled", successAutoScaling.AutoScalingEnabled); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting auto_scaling_enabled: %s", err), "(Data) ibm_db2_autoscale", "read", "set-auto_scaling_enabled").GetDiag()
	}

	if err = d.Set("auto_scaling_max_storage", flex.IntValue(successAutoScaling.AutoScalingMaxStorage)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting auto_scaling_max_storage: %s", err), "(Data) ibm_db2_autoscale", "read", "set-auto_scaling_max_storage").GetDiag()
	}

	if err = d.Set("auto_scaling_over_time_period", flex.IntValue(successAutoScaling.AutoScalingOverTimePeriod)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting auto_scaling_over_time_period: %s", err), "(Data) ibm_db2_autoscale", "read", "set-auto_scaling_over_time_period").GetDiag()
	}

	if err = d.Set("auto_scaling_pause_limit", flex.IntValue(successAutoScaling.AutoScalingPauseLimit)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting auto_scaling_pause_limit: %s", err), "(Data) ibm_db2_autoscale", "read", "set-auto_scaling_pause_limit").GetDiag()
	}

	if err = d.Set("auto_scaling_threshold", flex.IntValue(successAutoScaling.AutoScalingThreshold)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting auto_scaling_threshold: %s", err), "(Data) ibm_db2_autoscale", "read", "set-auto_scaling_threshold").GetDiag()
	}

	if err = d.Set("storage_unit", successAutoScaling.StorageUnit); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting storage_unit: %s", err), "(Data) ibm_db2_autoscale", "read", "set-storage_unit").GetDiag()
	}

	if err = d.Set("storage_utilization_percentage", flex.IntValue(successAutoScaling.StorageUtilizationPercentage)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting storage_utilization_percentage: %s", err), "(Data) ibm_db2_autoscale", "read", "set-storage_utilization_percentage").GetDiag()
	}

	if err = d.Set("support_auto_scaling", successAutoScaling.SupportAutoScaling); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting support_auto_scaling: %s", err), "(Data) ibm_db2_autoscale", "read", "set-support_auto_scaling").GetDiag()
	}

	return nil
}

// dataSourceIbmDb2SaasAutoscaleID returns a reasonable ID for the list.
func dataSourceIbmDb2AutoscaleID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
