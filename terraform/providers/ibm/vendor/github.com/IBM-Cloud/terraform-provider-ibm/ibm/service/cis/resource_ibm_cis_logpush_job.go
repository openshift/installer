// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/networking-go-sdk/logpushjobsapiv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	cisLogpushJobID                 = "job_id"
	cisLogpushName                  = "name"
	cisLogpushEnabled               = "enabled"
	cisLogpullOpt                   = "logpull_options"
	cisLogdna                       = "logdna"
	cisLogPushCos                   = "cos"
	cisLogpushCosOwnershipChallenge = "ownership_challenge"
	cisLogPushIbmCl                 = "ibmcl"
	cisLogPushIbmClInstanceId       = "instance_id"
	cisLogPushIbmClRegion           = "region"
	cisLogPushIbmClApiKey           = "api_key"
	cisLogpushDataset               = "dataset"
	cisLogpushFreq                  = "frequency"
	cisLogpushDestConf              = "destination_conf"
	cisLogpushLastComplete          = "last_complete"
	cisLogpushLastError             = "last_error"
	cisLogpushErrorMessage          = "error_message"
)

func ResourceIBMCISLogPushJob() *schema.Resource {
	return &schema.Resource{
		Create:   ResourceIBMCISLogpushJobCreate,
		Read:     ResourceIBMCISLogpushJobRead,
		Update:   ResourceIBMCISLogpushJobUpdate,
		Delete:   ResourceIBMCISLogpushJobDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
				ValidateFunc: validate.InvokeValidator("ibm_cis_logpush_job",
					"cis_id"),
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Description:      "Associated CIS domain",
				Required:         true,
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisLogdna: {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{cisLogPushCos, cisLogPushIbmCl, cisLogpushDestConf},
				StateFunc: func(v interface{}) string {
					json, err := flex.NormalizeJSONString(v)
					if err != nil {
						return fmt.Sprintf("%q", err.Error())
					}
					return json
				},
				Description: "Information to identify the LogDNA instance the data will be pushed.",
			},
			cisLogPushCos: {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				RequiredWith:  []string{cisLogpushCosOwnershipChallenge},
				ConflictsWith: []string{cisLogdna, cisLogPushIbmCl, cisLogpushDestConf},
				StateFunc: func(v interface{}) string {
					json, err := flex.NormalizeJSONString(v)
					if err != nil {
						return fmt.Sprintf("%q", err.Error())
					}
					return json
				},
				Description: "Information to identify the COS bucket where the data will be pushed.",
			},
			cisLogpushCosOwnershipChallenge: {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Ownership challenge token to prove destination ownership.",
				Sensitive:    true,
				RequiredWith: []string{cisLogPushCos},
			},
			cisLogPushIbmCl: {
				Type:          schema.TypeList,
				Optional:      true,
				Description:   "Information to identify the IBM Cloud Log instance where the data will be pushed.",
				MaxItems:      1,
				Sensitive:     true,
				ConflictsWith: []string{cisLogdna, cisLogPushCos, cisLogpushDestConf},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						cisLogPushIbmClInstanceId: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID of the IBM Cloud Log instance where you want to send logs.",
						},
						cisLogPushIbmClRegion: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Region where the IBM Cloud Log instance is located.",
						},
						cisLogPushIbmClApiKey: {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "IBM Cloud API key used to generate a token for pushing to your IBM Cloud Log instance.",
						},
					},
				},
			},

			cisLogpushName: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Logpush Job Name",
			},
			cisLogpushEnabled: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether the logpush job enabled or not",
			},
			cisLogpullOpt: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Configuration string",
			},
			cisLogpushDataset: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Dataset to be pulled",
			},
			cisLogpushFreq: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The frequency at which CIS sends batches of logs to your destination",
			},
			cisLogpushJobID: {
				Type:        schema.TypeInt,
				Description: "Associated CIS domain",
				Computed:    true,
			},
			cisLogpushDestConf: {
				Type:          schema.TypeString,
				Computed:      true,
				Optional:      true,
				ConflictsWith: []string{cisLogdna, cisLogPushCos, cisLogPushIbmCl},
				Description:   "Uniquely identifies a resource (such as an s3 bucket) where data will be pushed.",
			},
			cisLogpushLastComplete: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Records the last time that logs have been successfully pushed.",
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
	}
}
func ResourceIBMCISLogPushJobValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cis_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "resource_instance",
			CloudDataRange:             []string{"service:internet-svcs"},
			Required:                   true})
	ibmCISLogPushValidator := validate.ResourceValidator{
		ResourceName: "ibm_cis_logpush_job",
		Schema:       validateSchema}
	return &ibmCISLogPushValidator
}

func ResourceIBMCISLogpushJobCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(conns.ClientSession).CisLogpushJobsSession()
	if err != nil {
		return err
	}
	crn := d.Get(cisID).(string)
	zoneID, _, _ := flex.ConvertTftoCisTwoVar(d.Get(cisDomainID).(string))
	sess.Crn = core.StringPtr(crn)
	sess.ZoneID = core.StringPtr(zoneID)

	logpushJob := &logpushjobsapiv1.CreateLogpushJobV2Request{}

	if a, ok := d.GetOk(cisLogpushName); ok {
		name := a.(string)
		logpushJob.Name = &name
	}
	if e, ok := d.GetOk(cisLogpushEnabled); ok {
		enabled := e.(bool)
		logpushJob.Enabled = &enabled
	}
	if lp, ok := d.GetOk(cisLogpullOpt); ok {
		logpullopt := lp.(string)
		logpushJob.LogpullOptions = &logpullopt
	}
	if log, ok := d.GetOk(cisLogdna); ok {
		var logDNA map[string]interface{}
		json.Unmarshal([]byte(log.(string)), &logDNA)
		logpushJob.Logdna = logDNA
	}
	if cos, ok := d.GetOk(cisLogPushCos); ok {
		var logCos map[string]interface{}
		json.Unmarshal([]byte(cos.(string)), &logCos)
		logpushJob.Cos = logCos
	}
	if d, ok := d.GetOk(cisLogpushDataset); ok {
		dataset := d.(string)
		logpushJob.Dataset = &dataset
	}
	if f, ok := d.GetOk(cisLogpushFreq); ok {
		freq := f.(string)
		logpushJob.Frequency = &freq
	}
	if f, ok := d.GetOk(cisLogpushCosOwnershipChallenge); ok {
		ownChallenge := f.(string)
		logpushJob.OwnershipChallenge = &ownChallenge
	}
	if f, ok := d.GetOk(cisLogpushDestConf); ok {
		destConf := f.(string)
		logpushJob.DestinationConf = &destConf
	}
	if f, ok := d.GetOk(cisLogPushIbmCl); ok {
		logpushJob.Ibmcl = extractLogPushIbmClValues(f.([]interface{}))
	}

	options := &logpushjobsapiv1.CreateLogpushJobV2Options{
		CreateLogpushJobV2Request: logpushJob,
	}

	result, response, err := sess.CreateLogpushJobV2(options)
	if err != nil {
		log.Printf("[DEBUG] Instance err %s\n%s", err, response)
		return err
	}
	JobID := strconv.Itoa(int(*result.Result.ID))

	d.SetId(flex.ConvertCisToTfThreeVar(JobID, zoneID, crn))
	return ResourceIBMCISLogpushJobRead(d, meta)
}

func ResourceIBMCISLogpushJobRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(conns.ClientSession).CisLogpushJobsSession()
	if err != nil {
		return err
	}
	logpushID, zoneID, crn, err := flex.ConvertTfToCisThreeVar(d.Id())
	if err != nil {
		return fmt.Errorf("[ERROR] Error Converting ConvertTfToCisThreeVar in Read")
	}
	sess.Crn = core.StringPtr(crn)
	sess.ZoneID = core.StringPtr(zoneID)

	opt := sess.NewGetLogpushJobV2Options(logpushID)
	result, response, err := sess.GetLogpushJobV2(opt)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error While Reading the Logpushjobs for LogDNA %s:%s", err, response)
	}
	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	d.Set(cisLogpushJobID, int64(*result.Result.ID))
	d.Set(cisLogpushName, *result.Result.Name)
	d.Set(cisLogpushEnabled, *result.Result.Enabled)
	d.Set(cisLogpushDataset, *result.Result.Dataset)
	d.Set(cisLogpushFreq, *result.Result.Frequency)
	d.Set(cisLogpullOpt, *result.Result.LogpullOptions)
	d.Set(cisLogpushDestConf, *result.Result.DestinationConf)
	if result.Result.LastComplete != nil {
		d.Set(cisLogpushLastComplete, *result.Result.LastComplete)
	}
	if result.Result.LastError != nil {
		d.Set(cisLogpushLastError, result.Result.LastError)
	}
	if result.Result.ErrorMessage != nil {
		d.Set(cisLogpushErrorMessage, result.Result.ErrorMessage)
	}
	return nil
}
func ResourceIBMCISLogpushJobUpdate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(conns.ClientSession).CisLogpushJobsSession()
	if err != nil {
		return err
	}

	crn := d.Get(cisID).(string)
	zoneID, _, _ := flex.ConvertTftoCisTwoVar(d.Get(cisDomainID).(string))
	sess.Crn = core.StringPtr(crn)
	sess.ZoneID = core.StringPtr(zoneID)

	logpushID, _, _, err := flex.ConvertTfToCisThreeVar(d.Id())

	if err != nil {
		return fmt.Errorf("[ERROR] Error Converting ConvertTfToCisThreeVar in Update")
	}
	if d.HasChange(cisLogpushEnabled) ||
		d.HasChange(cisLogpullOpt) ||
		d.HasChange(cisLogdna) ||
		d.HasChange(cisLogpushFreq) ||
		d.HasChange(cisLogPushCos) ||
		d.HasChange(cisLogPushIbmCl) ||
		d.HasChange(cisLogpushCosOwnershipChallenge) ||
		d.HasChange(cisLogpushDestConf) {

		updateLogpushJob := &logpushjobsapiv1.UpdateLogpushJobV2Request{}

		if e, ok := d.GetOk(cisLogpushEnabled); ok {
			enabled := e.(bool)
			updateLogpushJob.Enabled = &enabled
		}
		if lp, ok := d.GetOk(cisLogpullOpt); ok {
			logpullopt := lp.(string)
			updateLogpushJob.LogpullOptions = &logpullopt
		}
		if log, ok := d.GetOk(cisLogdna); ok {
			var logDNA map[string]interface{}
			json.Unmarshal([]byte(log.(string)), &logDNA)
			updateLogpushJob.Logdna = logDNA
		}
		if f, ok := d.GetOk(cisLogpushFreq); ok {
			freq := f.(string)
			updateLogpushJob.Frequency = &freq
		}
		if cos, ok := d.GetOk(cisLogPushCos); ok {
			var logCos map[string]interface{}
			json.Unmarshal([]byte(cos.(string)), &logCos)
			updateLogpushJob.Cos = logCos
		}
		if f, ok := d.GetOk(cisLogpushCosOwnershipChallenge); ok {
			ownChallenge := f.(string)
			updateLogpushJob.OwnershipChallenge = &ownChallenge
		}
		if d.HasChange(cisLogpushDestConf) {
			if f, ok := d.GetOk(cisLogpushDestConf); ok {
				destConf := f.(string)
				updateLogpushJob.DestinationConf = &destConf
			}
		}
		if f, ok := d.GetOk(cisLogPushIbmCl); ok {
			updateLogpushJob.Ibmcl = extractLogPushUpdateIbmClValues(f.([]interface{}))
		}
		options := &logpushjobsapiv1.UpdateLogpushJobV2Options{
			JobID:                     core.StringPtr(logpushID),
			UpdateLogpushJobV2Request: updateLogpushJob,
		}
		result, resp, err := sess.UpdateLogpushJobV2(options)
		if err != nil || result == nil {
			return fmt.Errorf("[ERROR] Error While Updating the Logpushjobs for LogDNA  %v, %v", err, resp)
		}
	}
	return ResourceIBMCISLogpushJobRead(d, meta)
}
func ResourceIBMCISLogpushJobDelete(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(conns.ClientSession).CisLogpushJobsSession()
	if err != nil {
		return err
	}
	crn := d.Get(cisID).(string)
	zoneID, _, _ := flex.ConvertTftoCisTwoVar(d.Get(cisDomainID).(string))
	sess.Crn = core.StringPtr(crn)
	sess.ZoneID = core.StringPtr(zoneID)

	logpushID, _, _, err := flex.ConvertTfToCisThreeVar(d.Id())
	if err != nil {
		return fmt.Errorf("[ERROR] Error Converting ConvertTfToCisThreeVar in Delete")
	}
	opt := sess.NewDeleteLogpushJobV2Options(logpushID)
	_, response, err := sess.DeleteLogpushJobV2(opt)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil
		}
		return fmt.Errorf("[ERROR] Error While Deleting the Logpushjob for LogDNA %s:%s", err, response)
	}
	d.SetId("")
	return nil
}

func extractLogPushIbmClValues(cisLogPushIbmClList []interface{}) *logpushjobsapiv1.LogpushJobIbmclReqIbmcl {
	cisLogPushIbmCl := cisLogPushIbmClList[0].(map[string]interface{})

	instanceID := cisLogPushIbmCl[cisLogPushIbmClInstanceId].(string)
	region := cisLogPushIbmCl[cisLogPushIbmClRegion].(string)
	apiKey := cisLogPushIbmCl[cisLogPushIbmClApiKey].(string)

	logPushIbmclReq := logpushjobsapiv1.LogpushJobIbmclReqIbmcl{
		InstanceID: &instanceID,
		Region:     &region,
		ApiKey:     &apiKey,
	}

	return &logPushIbmclReq
}

func extractLogPushUpdateIbmClValues(cisLogPushIbmClList []interface{}) *logpushjobsapiv1.LogpushJobsUpdateIbmclReqIbmcl {
	cisLogPushIbmCl := cisLogPushIbmClList[0].(map[string]interface{})

	instanceID := cisLogPushIbmCl[cisLogPushIbmClInstanceId].(string)
	region := cisLogPushIbmCl[cisLogPushIbmClRegion].(string)
	apiKey := cisLogPushIbmCl[cisLogPushIbmClApiKey].(string)

	logPushIbmclReq := logpushjobsapiv1.LogpushJobsUpdateIbmclReqIbmcl{
		InstanceID: &instanceID,
		Region:     &region,
		ApiKey:     &apiKey,
	}

	return &logPushIbmclReq
}
