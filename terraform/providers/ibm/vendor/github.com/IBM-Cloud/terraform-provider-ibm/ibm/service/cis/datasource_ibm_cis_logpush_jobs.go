// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
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
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Description:      "Associated CIS domain",
				Required:         true,
				DiffSuppressFunc: suppressDomainIDDiff,
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
					},
				},
			},
		},
	}
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
	opt := sess.NewGetLogpushJobsV2Options()
	result, resp, err := sess.GetLogpushJobsV2(opt)
	if err != nil {
		log.Printf("[WARN] List all Logpush jobs failed: %v\n", resp)
		return err
	}
	logPushList := make([]map[string]interface{}, 0)
	for _, logpushObj := range result.Result {
		logPushOpt := map[string]interface{}{}
		logPushOpt[cisLogpushJobID] = int64(*logpushObj.ID)
		logPushOpt[cisLogpushName] = *logpushObj.Name
		logPushOpt[cisLogpullOpt] = *logpushObj.LogpullOptions
		logPushOpt[cisLogpushEnabled] = *logpushObj.Enabled
		logPushOpt[cisLogpushDataset] = *logpushObj.Dataset
		logPushOpt[cisLogpushFreq] = *logpushObj.Frequency
		logPushOpt[cisLogpushDestConf] = *logpushObj.DestinationConf
		logPushList = append(logPushList, logPushOpt)
	}
	d.SetId(dataSourceCISLogpushJobsCheckID(d))
	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	d.Set(cisLogpushJobs, logPushList)
	return nil
}
func dataSourceCISLogpushJobsCheckID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
