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
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/networking-go-sdk/logpushjobsapiv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	cisLogpushJobID    = "job_id"
	cisLogpushName     = "name"
	cisLogpushEnabled  = "enabled"
	cisLogpullOpt      = "logpull_options"
	cisLogdna          = "logdna"
	cisLogpushDataset  = "dataset"
	cisLogpushFreq     = "frequency"
	cisLogpushDestConf = "destination_conf"
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
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Description:      "Associated CIS domain",
				Required:         true,
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisLogdna: {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				StateFunc: func(v interface{}) string {
					json, err := flex.NormalizeJSONString(v)
					if err != nil {
						return fmt.Sprintf("%q", err.Error())
					}
					return json
				},
				Description: "Information to identify the LogDNA instance the data will be pushed.",
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
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Uniquely identifies a resource (such as an s3 bucket) where data will be pushed.",
			},
		},
	}
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

	logpushJob := &logpushjobsapiv1.CreateLogpushJobV2RequestLogpushJobLogdnaReq{}

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
		var logDNA interface{}
		json.Unmarshal([]byte(log.(string)), &logDNA)
		logpushJob.Logdna = logDNA
	}
	if d, ok := d.GetOk(cisLogpushDataset); ok {
		dataset := d.(string)
		logpushJob.Dataset = &dataset
	}
	if f, ok := d.GetOk(cisLogpushFreq); ok {
		freq := f.(string)
		logpushJob.Frequency = &freq
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
	logpushID, zoneID, crn, _ := flex.ConvertTfToCisThreeVar(d.Id())
	if err != nil {
		return fmt.Errorf("[ERROR] Error Converting ConvertTfToCisThreeVar in Read")
	}
	JobID, _ := strconv.Atoi(logpushID)
	sess.Crn = core.StringPtr(crn)
	sess.ZoneID = core.StringPtr(zoneID)

	opt := sess.NewGetLogpushJobV2Options(int64(JobID))
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

	logpushID, zoneID, crn, _ := flex.ConvertTfToCisThreeVar(d.Id())
	JobId, _ := strconv.Atoi(logpushID)

	if err != nil {
		return fmt.Errorf("[ERROR] Error Converting ConvertTfToCisThreeVar in Update")
	}
	if d.HasChange(cisLogpushEnabled) ||
		d.HasChange(cisLogpullOpt) ||
		d.HasChange(cisLogdna) ||
		d.HasChange(cisLogpushFreq) {

		updateLogpushJob := &logpushjobsapiv1.UpdateLogpushJobV2RequestLogpushJobsUpdateLogdnaReq{}

		if e, ok := d.GetOk(cisLogpushEnabled); ok {
			enabled := e.(bool)
			updateLogpushJob.Enabled = &enabled
		}
		if lp, ok := d.GetOk(cisLogpullOpt); ok {
			logpullopt := lp.(string)
			updateLogpushJob.LogpullOptions = &logpullopt
		}
		if log, ok := d.GetOk(cisLogdna); ok {
			var logDNA interface{}
			json.Unmarshal([]byte(log.(string)), &logDNA)
			updateLogpushJob.Logdna = logDNA
		}
		if f, ok := d.GetOk(cisLogpushFreq); ok {
			freq := f.(string)
			updateLogpushJob.Frequency = &freq
		}
		options := &logpushjobsapiv1.UpdateLogpushJobV2Options{
			JobID:                     core.Int64Ptr(int64(JobId)),
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

	logpushID, zoneID, crn, _ := flex.ConvertTfToCisThreeVar(d.Id())
	JobID, _ := strconv.Atoi(logpushID)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Converting ConvertTfToCisThreeVar in Delete")
	}
	opt := sess.NewDeleteLogpushJobV2Options(int64(JobID))
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
