// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"log"
	"strconv"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	cisLogpushJobs = "logpush_job_pack"
)

func DataSourceIBMCISLogPushJobs() *schema.Resource {
	return &schema.Resource{
		Read: ResourceIBMCISLogpushJobsRead,
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
				ValidateFunc: validate.InvokeDataSourceValidator(
					"ibm_cis_logpush_jobs",
					"cis_id"),
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Description:      "Associated CIS domain",
				Required:         true,
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisLogpushJobID: {
				Type:        schema.TypeInt,
				Description: "ID of the Job",
				Optional:    true,
			},
			cisLogpushJobs: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "logpush jobs information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						cisLogpushJobID: {
							Type:        schema.TypeInt,
							Description: "Associated CIS domain ",
							Computed:    true,
						},
						cisLogpushName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Logpush Job Name",
						},
						cisLogpushEnabled: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the logpush job enabled or not",
						},
						cisLogpullOpt: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Configuration string",
						},
						cisLogpushDestConf: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Uniquely identifies a resource (such as an s3 bucket) where data will be pushed.",
						},
						cisLogpushDataset: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Dataset to be pulled",
						},
						cisLogpushFreq: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The frequency at which CIS sends batches of logs to your destination",
						},
						cisLogpushLastComplete: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Records the last time for which logs have been successfully pushed.",
						},
						cisLogpushLastError: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Records the last time the job failed.",
						},
						cisLogpushErrorMessage: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The last failure message.",
						},
					},
				},
			},
		},
	}
}
func DataSourceIBMCISLogPushJobsValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cis_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "resource_instance",
			CloudDataRange:             []string{"service:internet-svcs"},
			Required:                   true})

	iBMCISLogPushJobsValidator := validate.ResourceValidator{
		ResourceName: "ibm_cis_logpush_jobs",
		Schema:       validateSchema}
	return &iBMCISLogPushJobsValidator
}
func ResourceIBMCISLogpushJobsRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(conns.ClientSession).CisLogpushJobsSession()
	if err != nil {
		return err
	}
	crn := d.Get(cisID).(string)
	sess.Crn = core.StringPtr(crn)
	zoneID, _, _ := flex.ConvertTftoCisTwoVar(d.Get(cisDomainID).(string))
	sess.ZoneID = core.StringPtr(zoneID)
	jobId := d.Get(cisLogpushJobID).(int)
	logPushList := make([]map[string]interface{}, 0)

	if jobId != 0 {
		opt := sess.NewGetLogpushJobV2Options(strconv.Itoa(jobId))
		result, resp, err := sess.GetLogpushJobV2(opt)
		if err != nil {
			log.Printf("[WARN] Get Logpush job failed: %v\n", resp)
			return err
		}
		logPushOpt := map[string]interface{}{}
		logPushOpt[cisLogpushJobID] = int64(*result.Result.ID)
		logPushOpt[cisLogpushName] = *result.Result.Name
		logPushOpt[cisLogpullOpt] = *result.Result.LogpullOptions
		logPushOpt[cisLogpushEnabled] = *result.Result.Enabled
		logPushOpt[cisLogpushDataset] = *result.Result.Dataset
		logPushOpt[cisLogpushFreq] = *result.Result.Frequency
		logPushOpt[cisLogpushDestConf] = *result.Result.DestinationConf
		if result.Result.LastComplete != nil {
			logPushOpt[cisLogpushLastComplete] = *result.Result.LastComplete
		}
		if result.Result.LastError != nil {
			logPushOpt[cisLogpushLastError] = *result.Result.LastError
		}
		if result.Result.ErrorMessage != nil {
			logPushOpt[cisLogpushErrorMessage] = *result.Result.ErrorMessage
		}
		logPushList = append(logPushList, logPushOpt)

	} else {

		opt := sess.NewGetLogpushJobsV2Options()
		result, resp, err := sess.GetLogpushJobsV2(opt)
		if err != nil {
			log.Printf("[WARN] List all Logpush jobs failed: %v\n", resp)
			return err
		}

		for _, logpushObj := range result.Result {
			logPushOpt := map[string]interface{}{}
			logPushOpt[cisLogpushJobID] = int64(*logpushObj.ID)
			logPushOpt[cisLogpushName] = *logpushObj.Name
			logPushOpt[cisLogpullOpt] = *logpushObj.LogpullOptions
			logPushOpt[cisLogpushEnabled] = *logpushObj.Enabled
			logPushOpt[cisLogpushDataset] = *logpushObj.Dataset
			logPushOpt[cisLogpushFreq] = *logpushObj.Frequency
			logPushOpt[cisLogpushDestConf] = *logpushObj.DestinationConf
			if logpushObj.LastComplete != nil {
				logPushOpt[cisLogpushLastComplete] = *logpushObj.LastComplete
			}
			if logpushObj.LastError != nil {
				logPushOpt[cisLogpushLastError] = *logpushObj.LastError
			}
			if logpushObj.ErrorMessage != nil {
				logPushOpt[cisLogpushErrorMessage] = *logpushObj.ErrorMessage
			}

			logPushList = append(logPushList, logPushOpt)
		}
	}
	d.SetId(dataSourceCISLogpushJobsCheckID())
	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	d.Set(cisLogpushJobs, logPushList)
	return nil
}
func dataSourceCISLogpushJobsCheckID() string {
	return time.Now().UTC().String()
}
