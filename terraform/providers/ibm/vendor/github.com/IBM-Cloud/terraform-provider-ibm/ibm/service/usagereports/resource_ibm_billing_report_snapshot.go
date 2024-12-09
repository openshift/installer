// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package usagereports

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/usagereportsv4"
)

func ResourceIBMBillingReportSnapshot() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMBillingReportSnapshotCreate,
		ReadContext:   resourceIBMBillingReportSnapshotRead,
		DeleteContext: resourceIBMBillingReportSnapshotDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"interval": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_billing_report_snapshot", "interval"),
				Description:  "Frequency of taking the snapshot of the billing reports.",
			},
			"versioning": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "new",
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_billing_report_snapshot", "versioning"),
				Description:  "A new version of report is created or the existing report version is overwritten with every update.",
			},
			"report_types": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "The type of billing reports to take snapshot of. Possible values are [account_summary, enterprise_summary, account_resource_instance_usage].",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"cos_reports_folder": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "IBMCloud-Billing-Reports",
				ForceNew:    true,
				Description: "The billing reports root folder to store the billing reports snapshots. Defaults to \"IBMCloud-Billing-Reports\".",
			},
			"cos_bucket": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the COS bucket to store the snapshot of the billing reports.",
			},
			"cos_location": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Region of the COS instance.",
			},
			"state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the billing snapshot configuration. Possible values are [enabled, disabled].",
			},
			"account_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of account. Possible values are [enterprise, account].",
			},
			"compression": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Compression format of the snapshot report.",
			},
			"content_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of content stored in snapshot report.",
			},
			"cos_endpoint": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The endpoint of the COS instance.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Timestamp in milliseconds when the snapshot configuration was created.",
			},
			"last_updated_at": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Timestamp in milliseconds when the snapshot configuration was last updated.",
			},
			"history": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of previous versions of the snapshot configurations.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Timestamp in milliseconds when the snapshot configuration was created.",
						},
						"end_time": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Timestamp in milliseconds when the snapshot configuration ends.",
						},
						"updated_by": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Account that updated the billing snapshot configuration.",
						},
						"account_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Account ID for which billing report snapshot is configured.",
						},
						"state": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of the billing snapshot configuration. Possible values are [enabled, disabled].",
						},
						"account_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of account. Possible values [enterprise, account].",
						},
						"interval": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Frequency of taking the snapshot of the billing reports.",
						},
						"versioning": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A new version of report is created or the existing report version is overwritten with every update.",
						},
						"report_types": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The type of billing reports to take snapshot of. Possible values are [account_summary, enterprise_summary, account_resource_instance_usage].",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"compression": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Compression format of the snapshot report.",
						},
						"content_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of content stored in snapshot report.",
						},
						"cos_reports_folder": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The billing reports root folder to store the billing reports snapshots. Defaults to \"IBMCloud-Billing-Reports\".",
						},
						"cos_bucket": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the COS bucket to store the snapshot of the billing reports.",
						},
						"cos_location": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region of the COS instance.",
						},
						"cos_endpoint": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The endpoint of the COS instance.",
						},
					},
				},
			},
		},
	}
}

func ResourceIBMBillingReportSnapshotValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "interval",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "daily",
		},
		validate.ValidateSchema{
			Identifier:                 "versioning",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "new, overwrite",
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_billing_report_snapshot", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMBillingReportSnapshotCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	usageReportsClient, err := meta.(conns.ClientSession).UsageReportsV4()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_billing_report_snapshot", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createReportsSnapshotConfigOptions := &usagereportsv4.CreateReportsSnapshotConfigOptions{}

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_billing_report_snapshot", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createReportsSnapshotConfigOptions.SetAccountID(userDetails.UserAccount)
	createReportsSnapshotConfigOptions.SetInterval(d.Get("interval").(string))
	createReportsSnapshotConfigOptions.SetCosBucket(d.Get("cos_bucket").(string))
	createReportsSnapshotConfigOptions.SetCosLocation(d.Get("cos_location").(string))
	if _, ok := d.GetOk("cos_reports_folder"); ok {
		createReportsSnapshotConfigOptions.SetCosReportsFolder(d.Get("cos_reports_folder").(string))
	}
	if _, ok := d.GetOk("report_types"); ok {
		var reportTypes []string
		for _, v := range d.Get("report_types").([]interface{}) {
			reportTypesItem := v.(string)
			reportTypes = append(reportTypes, reportTypesItem)
		}
		createReportsSnapshotConfigOptions.SetReportTypes(reportTypes)
	}
	if _, ok := d.GetOk("versioning"); ok {
		createReportsSnapshotConfigOptions.SetVersioning(d.Get("versioning").(string))
	}

	snapshotConfig, _, err := usageReportsClient.CreateReportsSnapshotConfigWithContext(context, createReportsSnapshotConfigOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateReportsSnapshotConfigWithContext failed: %s", err.Error()), "ibm_billing_report_snapshot", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*snapshotConfig.AccountID)

	return resourceIBMBillingReportSnapshotRead(context, d, meta)
}

func resourceIBMBillingReportSnapshotRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	usageReportsClient, err := meta.(conns.ClientSession).UsageReportsV4()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_billing_report_snapshot", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getReportsSnapshotConfigOptions := &usagereportsv4.GetReportsSnapshotConfigOptions{}

	getReportsSnapshotConfigOptions.SetAccountID(d.Id())

	snapshotConfig, response, err := usageReportsClient.GetReportsSnapshotConfigWithContext(context, getReportsSnapshotConfigOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetReportsSnapshotConfigWithContext failed: %s", err.Error()), "ibm_billing_report_snapshot", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("interval", snapshotConfig.Interval); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting interval: %s", err))
	}
	if !core.IsNil(snapshotConfig.Versioning) {
		if err = d.Set("versioning", snapshotConfig.Versioning); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting versioning: %s", err))
		}
	}
	if !core.IsNil(snapshotConfig.ReportTypes) {
		if err = d.Set("report_types", snapshotConfig.ReportTypes); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting report_types: %s", err))
		}
	}
	if !core.IsNil(snapshotConfig.CosReportsFolder) {
		if err = d.Set("cos_reports_folder", snapshotConfig.CosReportsFolder); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting cos_reports_folder: %s", err))
		}
	}
	if err = d.Set("cos_bucket", snapshotConfig.CosBucket); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting cos_bucket: %s", err))
	}
	if err = d.Set("cos_location", snapshotConfig.CosLocation); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting cos_location: %s", err))
	}
	if !core.IsNil(snapshotConfig.State) {
		if err = d.Set("state", snapshotConfig.State); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting state: %s", err))
		}
	}
	if !core.IsNil(snapshotConfig.AccountType) {
		if err = d.Set("account_type", snapshotConfig.AccountType); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting account_type: %s", err))
		}
	}
	if !core.IsNil(snapshotConfig.Compression) {
		if err = d.Set("compression", snapshotConfig.Compression); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting compression: %s", err))
		}
	}
	if !core.IsNil(snapshotConfig.ContentType) {
		if err = d.Set("content_type", snapshotConfig.ContentType); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting content_type: %s", err))
		}
	}
	if !core.IsNil(snapshotConfig.CosEndpoint) {
		if err = d.Set("cos_endpoint", snapshotConfig.CosEndpoint); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting cos_endpoint: %s", err))
		}
	}
	if !core.IsNil(snapshotConfig.CreatedAt) {
		if err = d.Set("created_at", flex.IntValue(snapshotConfig.CreatedAt)); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
		}
	}
	if !core.IsNil(snapshotConfig.LastUpdatedAt) {
		if err = d.Set("last_updated_at", flex.IntValue(snapshotConfig.LastUpdatedAt)); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting last_updated_at: %s", err))
		}
	}
	if !core.IsNil(snapshotConfig.History) {
		history := []map[string]interface{}{}
		for _, historyItem := range snapshotConfig.History {
			historyItemMap, err := resourceIBMBillingReportSnapshotSnapshotConfigHistoryItemToMap(&historyItem)
			if err != nil {
				return diag.FromErr(err)
			}
			history = append(history, historyItemMap)
		}
		if err = d.Set("history", history); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting history: %s", err))
		}
	}

	return nil
}

func resourceIBMBillingReportSnapshotDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	usageReportsClient, err := meta.(conns.ClientSession).UsageReportsV4()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_billing_report_snapshot", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteReportsSnapshotConfigOptions := &usagereportsv4.DeleteReportsSnapshotConfigOptions{}

	deleteReportsSnapshotConfigOptions.SetAccountID(d.Id())

	_, err = usageReportsClient.DeleteReportsSnapshotConfigWithContext(context, deleteReportsSnapshotConfigOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteReportsSnapshotConfigWithContext failed: %s", err.Error()), "ibm_billing_report_snapshot", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func resourceIBMBillingReportSnapshotSnapshotConfigHistoryItemToMap(model *usagereportsv4.SnapshotConfigHistoryItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.StartTime != nil {
		modelMap["start_time"] = flex.IntValue(model.StartTime)
	}
	if model.EndTime != nil {
		modelMap["end_time"] = flex.IntValue(model.EndTime)
	}
	if model.UpdatedBy != nil {
		modelMap["updated_by"] = model.UpdatedBy
	}
	if model.AccountID != nil {
		modelMap["account_id"] = model.AccountID
	}
	if model.State != nil {
		modelMap["state"] = model.State
	}
	if model.AccountType != nil {
		modelMap["account_type"] = model.AccountType
	}
	if model.Interval != nil {
		modelMap["interval"] = model.Interval
	}
	if model.Versioning != nil {
		modelMap["versioning"] = model.Versioning
	}
	if model.ReportTypes != nil {
		modelMap["report_types"] = model.ReportTypes
	}
	if model.Compression != nil {
		modelMap["compression"] = model.Compression
	}
	if model.ContentType != nil {
		modelMap["content_type"] = model.ContentType
	}
	if model.CosReportsFolder != nil {
		modelMap["cos_reports_folder"] = model.CosReportsFolder
	}
	if model.CosBucket != nil {
		modelMap["cos_bucket"] = model.CosBucket
	}
	if model.CosLocation != nil {
		modelMap["cos_location"] = model.CosLocation
	}
	if model.CosEndpoint != nil {
		modelMap["cos_endpoint"] = model.CosEndpoint
	}
	return modelMap, nil
}
