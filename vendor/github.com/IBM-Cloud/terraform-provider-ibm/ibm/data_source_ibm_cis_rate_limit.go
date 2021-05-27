// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"time"

	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMCISRateLimit() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMCISRateLimitRead,
		Schema: map[string]*schema.Schema{
			"cis_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rate_limit": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bypass": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"threshold": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"period": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"correlate": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"by": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"action": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mode": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"timeout": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"response": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"content_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"body": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"match": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"request": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"methods": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"schemes": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"url": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"response": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"status": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeInt},
												},
												"origin_traffic": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"headers": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"op": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"value": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMCISRateLimitRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisRLClientSession()
	if err != nil {
		return err
	}
	cisID := d.Get("cis_id").(string)
	zoneID, _, err := convertTftoCisTwoVar(d.Get("domain_id").(string))
	if err != nil {
		return err
	}
	cisClient.Crn = core.StringPtr(cisID)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)
	opt := cisClient.NewListAllZoneRateLimitsOptions()
	rateLimitRecord, resp, err := cisClient.ListAllZoneRateLimits(opt)
	if err != nil {
		return fmt.Errorf("Failed to read RateLimit: %v", resp)
	}
	rules := make([]map[string]interface{}, 0)
	for _, r := range rateLimitRecord.Result {
		rule := map[string]interface{}{}
		rule["rule_id"] = *r.ID
		rule["disabled"] = *r.Disabled
		rule["description"] = *r.Description
		rule["threshold"] = *r.Threshold
		rule["period"] = *r.Period
		rule["action"] = flattenRateLimitAction(r.Action)
		rule["match"] = flattenRateLimitMatch(r.Match)
		rule["correlate"] = flattenRateLimitCorrelate(r.Correlate)
		rule["bypass"] = flattenRateLimitByPass(r.Bypass)
		rules = append(rules, rule)

	}
	d.SetId(dataSourceIBMCISRateLimitID(d))
	d.Set("rate_limit", rules)
	d.Set("cis_id", cisID)
	d.Set("domain_id", zoneID)
	return nil
}

func dataSourceIBMCISRateLimitID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
