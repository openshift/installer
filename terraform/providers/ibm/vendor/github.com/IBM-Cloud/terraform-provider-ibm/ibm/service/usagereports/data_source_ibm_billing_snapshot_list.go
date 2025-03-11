// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package usagereports

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/platform-services-go-sdk/usagereportsv4"
)

func DataSourceIBMBillingSnapshotList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMBillingSnapshotListRead,

		Schema: map[string]*schema.Schema{
			"month": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The month for which billing report snapshot is requested.  Format is yyyy-mm.",
			},
			"date_from": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Timestamp in milliseconds for which billing report snapshot is requested.",
			},
			"date_to": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Timestamp in milliseconds for which billing report snapshot is requested.",
			},
			"limit": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     30,
				Description: "Number of usage records returned. The default value is 30. Maximum value is 200.",
			},
			"snapshotcount": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of total snapshots.",
			},
			"snapshots": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Account ID for which billing report snapshot is configured.",
						},
						"month": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Month of captured snapshot.",
						},
						"account_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of account. Possible values are [enterprise, account].",
						},
						"expected_processed_at": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Timestamp of snapshot processed.",
						},
						"state": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of the billing snapshot configuration. Possible values are [enabled, disabled].",
						},
						"billing_period": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Period of billing in snapshot.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Date and time of start of billing in the respective snapshot.",
									},
									"end": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Date and time of end of billing in the respective snapshot.",
									},
								},
							},
						},
						"snapshot_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of the snapshot captured.",
						},
						"charset": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Character encoding used.",
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
						"bucket": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the COS bucket to store the snapshot of the billing reports.",
						},
						"version": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Version of the snapshot.",
						},
						"created_on": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Date and time of creation of snapshot.",
						},
						"report_types": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of report types configured for the snapshot.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of billing report of the snapshot. Possible values are [account_summary, enterprise_summary, account_resource_instance_usage].",
									},
									"version": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Version of the snapshot.",
									},
								},
							},
						},
						"files": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of location of reports.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"report_types": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of billing report stored. Possible values are [account_summary, enterprise_summary, account_resource_instance_usage].",
									},
									"location": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Absolute path of the billing report in the COS instance.",
									},
									"account_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Account ID for which billing report is captured.",
									},
								},
							},
						},
						"processed_at": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Timestamp at which snapshot is captured.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMBillingSnapshotListRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	usageReportsClient, err := meta.(conns.ClientSession).UsageReportsV4()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_billing_snapshot_list", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_billing_snapshot_list", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getReportsSnapshotOptions := &usagereportsv4.GetReportsSnapshotOptions{}

	getReportsSnapshotOptions.SetAccountID(userDetails.UserAccount)
	getReportsSnapshotOptions.SetMonth(d.Get("month").(string))
	if _, ok := d.GetOk("date_from"); ok {
		getReportsSnapshotOptions.SetDateFrom(int64(d.Get("date_from").(int)))
	}
	if _, ok := d.GetOk("date_to"); ok {
		getReportsSnapshotOptions.SetDateTo(int64(d.Get("date_to").(int)))
	}
	if _, ok := d.GetOk("limit"); ok {
		getReportsSnapshotOptions.SetLimit(int64(d.Get("limit").(int)))
	}

	var pager *usagereportsv4.GetReportsSnapshotPager
	pager, err = usageReportsClient.NewGetReportsSnapshotPager(getReportsSnapshotOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_billing_snapshot_list", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	snapshotList, err := pager.GetAll()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetReportsSnapshotPager.GetAll() failed %s", err), "(Data) ibm_billing_snapshot_list", "read")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIBMBillingSnapshotListID(d))

	snapshots := []map[string]interface{}{}
	for _, modelItem := range snapshotList {
		modelMap, err := dataSourceIBMBillingSnapshotListSnapshotListSnapshotsItemToMap(&modelItem)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_billing_snapshot_list", "read")
			return tfErr.GetDiag()
		}
		snapshots = append(snapshots, modelMap)
	}

	if err = d.Set("snapshots", snapshots); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting snapshots %s", err), "(Data) ibm_billing_snapshot_list", "read")
		return tfErr.GetDiag()
	}

	return nil
}

// dataSourceIBMBillingSnapshotListID returns a reasonable ID for the list.
func dataSourceIBMBillingSnapshotListID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIBMBillingSnapshotListSnapshotListSnapshotsItemToMap(model *usagereportsv4.SnapshotListSnapshotsItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AccountID != nil {
		modelMap["account_id"] = model.AccountID
	}
	if model.Month != nil {
		modelMap["month"] = model.Month
	}
	if model.AccountType != nil {
		modelMap["account_type"] = model.AccountType
	}
	if model.ExpectedProcessedAt != nil {
		modelMap["expected_processed_at"] = flex.IntValue(model.ExpectedProcessedAt)
	}
	if model.State != nil {
		modelMap["state"] = model.State
	}
	if model.BillingPeriod != nil {
		billingPeriodMap, err := dataSourceIBMBillingSnapshotListSnapshotListSnapshotsItemBillingPeriodToMap(model.BillingPeriod)
		if err != nil {
			return modelMap, err
		}
		modelMap["billing_period"] = []map[string]interface{}{billingPeriodMap}
	}
	if model.SnapshotID != nil {
		modelMap["snapshot_id"] = model.SnapshotID
	}
	if model.Charset != nil {
		modelMap["charset"] = model.Charset
	}
	if model.Compression != nil {
		modelMap["compression"] = model.Compression
	}
	if model.ContentType != nil {
		modelMap["content_type"] = model.ContentType
	}
	if model.Bucket != nil {
		modelMap["bucket"] = model.Bucket
	}
	if model.Version != nil {
		modelMap["version"] = model.Version
	}
	if model.CreatedOn != nil {
		modelMap["created_on"] = model.CreatedOn
	}
	if model.ReportTypes != nil {
		reportTypes := []map[string]interface{}{}
		for _, reportTypesItem := range model.ReportTypes {
			reportTypesItemMap, err := dataSourceIBMBillingSnapshotListSnapshotListSnapshotsItemReportTypesItemToMap(&reportTypesItem)
			if err != nil {
				return modelMap, err
			}
			reportTypes = append(reportTypes, reportTypesItemMap)
		}
		modelMap["report_types"] = reportTypes
	}
	if model.Files != nil {
		files := []map[string]interface{}{}
		for _, filesItem := range model.Files {
			filesItemMap, err := dataSourceIBMBillingSnapshotListSnapshotListSnapshotsItemFilesItemToMap(&filesItem)
			if err != nil {
				return modelMap, err
			}
			files = append(files, filesItemMap)
		}
		modelMap["files"] = files
	}
	if model.ProcessedAt != nil {
		modelMap["processed_at"] = flex.IntValue(model.ProcessedAt)
	}
	return modelMap, nil
}

func dataSourceIBMBillingSnapshotListSnapshotListSnapshotsItemBillingPeriodToMap(model *usagereportsv4.SnapshotListSnapshotsItemBillingPeriod) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Start != nil {
		modelMap["start"] = model.Start
	}
	if model.End != nil {
		modelMap["end"] = model.End
	}
	return modelMap, nil
}

func dataSourceIBMBillingSnapshotListSnapshotListSnapshotsItemReportTypesItemToMap(model *usagereportsv4.SnapshotListSnapshotsItemReportTypesItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	if model.Version != nil {
		modelMap["version"] = model.Version
	}
	return modelMap, nil
}

func dataSourceIBMBillingSnapshotListSnapshotListSnapshotsItemFilesItemToMap(model *usagereportsv4.SnapshotListSnapshotsItemFilesItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ReportTypes != nil {
		modelMap["report_types"] = model.ReportTypes
	}
	if model.Location != nil {
		modelMap["location"] = model.Location
	}
	if model.AccountID != nil {
		modelMap["account_id"] = model.AccountID
	}
	return modelMap, nil
}
