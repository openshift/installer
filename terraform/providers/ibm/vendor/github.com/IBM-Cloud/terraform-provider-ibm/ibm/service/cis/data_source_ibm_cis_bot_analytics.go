// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	cisBotAnalyticsType          = "type"
	cisBotAnalyticsSince         = "since"
	cisBotAnalyticsUntil         = "until"
	cisBotAnalyticsResult        = "result"
	cisBotAnalyticsScoreSource   = "score_source"
	cisBotAnalyticsTimeseries    = "timeseries"
	cisBotAnalyticsTopAttributes = "top_ns"
)

func DataSourceIBMCISBotAnalytics() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMCISBotAnalyticsRead,

		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
				ValidateFunc: validate.InvokeDataSourceValidator(
					"ibm_cis_bot_analytics",
					"cis_id"),
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Description:      "Associated CIS domain",
				Required:         true,
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisBotAnalyticsType: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Bot Analytics Type",
			},
			cisBotAnalyticsSince: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Datetime for start of query",
			},
			cisBotAnalyticsUntil: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Datetime for end of query",
			},
			cisBotAnalyticsResult: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Bot Analytics result",
				Elem:        &schema.Schema{Type: schema.TypeMap},
			},
		},
	}
}

func DataSourceIBMCISBotAnalyticsValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cis_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "resource_instance",
			CloudDataRange:             []string{"service:internet-svcs"},
			Required:                   true})

	iBMCISBotAnalyticsValidator := validate.ResourceValidator{
		ResourceName: "ibm_cis_bot_analytics",
		Schema:       validateSchema}
	return &iBMCISBotAnalyticsValidator
}

func dataSourceIBMCISBotAnalyticsRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisBotAnalyticsSession()
	if err != nil {
		return err
	}

	crn := d.Get(cisID).(string)
	zoneName := d.Get(cisDomainID).(string)
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneName)

	requestType := d.Get(cisBotAnalyticsType).(string)

	layout := time.RFC3339
	parseSinceTime, err := time.Parse(layout, d.Get(cisBotAnalyticsSince).(string))
	if err != nil {
		return err
	}
	since := strfmt.DateTime(parseSinceTime)
	parseUntilTime, err := time.Parse(layout, d.Get(cisBotAnalyticsUntil).(string))
	if err != nil {
		return err
	}
	until := strfmt.DateTime(parseUntilTime)

	d.SetId(crn)
	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneName)
	d.Set(cisBotAnalyticsType, requestType)

	if requestType == cisBotAnalyticsScoreSource {
		sourceScoreOpt := cisClient.NewGetBotScoreOptions(&since, &until)
		sourceScoreResult, sourceScoreResp, sourceScoreErr := cisClient.GetBotScore(sourceScoreOpt)
		if sourceScoreErr != nil {
			log.Printf("Get Source Score Failed with response : %s\n", sourceScoreResp)
			return sourceScoreErr
		}

		jsonRes, err := json.Marshal(sourceScoreResult.Result)
		if err != nil {
			log.Printf("Response Marshal failed: %s", err)
		}
		stringRes := strings.Replace(string(jsonRes), "\"", "'", -1)

		d.Set(cisBotAnalyticsResult, stringRes)

	} else if requestType == cisBotAnalyticsTimeseries {
		timeSeriesOpt := cisClient.NewGetBotTimeseriesOptions(&since, &until)
		timeSeriesResult, timeSeriesResp, timeSeriesErr := cisClient.GetBotTimeseries(timeSeriesOpt)
		if timeSeriesErr != nil {
			log.Printf("Get TimeSeries Failed with response : %s\n", timeSeriesResp)
			return timeSeriesErr
		}

		jsonRes, err := json.Marshal(timeSeriesResult.Result)
		if err != nil {
			log.Printf("Response Marshal failed: %s", err)
		}
		stringRes := strings.Replace(string(jsonRes), "\"", "'", -1)

		d.Set(cisBotAnalyticsResult, stringRes)

	} else if requestType == cisBotAnalyticsTopAttributes {
		topAttributesOpt := cisClient.NewGetBotTopnsOptions(&since, &until)
		topAttributesResult, topAttributesResp, topAttributesErr := cisClient.GetBotTopns(topAttributesOpt)
		if topAttributesErr != nil {
			log.Printf("Get topAttributes Failed with response : %s\n", topAttributesResp)
			return topAttributesErr
		}
		jsonRes, err := json.Marshal(topAttributesResult.Result)
		if err != nil {
			log.Printf("Response Marshal failed: %s", err)
		}
		stringRes := strings.Replace(string(jsonRes), "\"", "'", -1)

		d.Set(cisBotAnalyticsResult, stringRes)

	} else {
		log.Printf("dataSourceIBMCISBotAnalyticsRead - Wrong Type provided.")
	}

	return nil
}
