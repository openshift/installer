// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"log"
	"time"

	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	cisGLBHealthCheck = "cis_healthchecks"
)

func dataSourceIBMCISHealthChecks() *schema.Resource {
	return &schema.Resource{
		Read:     dataSourceIBMCISGLBHealthCheckRead,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS Zone CRN",
			},

			cisGLBHealthCheck: {
				Type:        schema.TypeList,
				Description: "Collection of GLB Health check/monitor detail",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "GLB Monitor/Health check id",
						},
						cisID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS Zone CRN",
						},
						cisGLBHealthCheckID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "GLB Monitor/Health check id",
						},
						cisGLBHealthCheckPath: {
							Type:        schema.TypeString,
							Description: "path",
							Computed:    true,
						},
						cisGLBHealthCheckExpectedBody: {
							Type:        schema.TypeString,
							Description: "expected_body",
							Computed:    true,
						},
						cisGLBHealthCheckExpectedCodes: {
							Type:        schema.TypeString,
							Description: "expected_codes",
							Computed:    true,
						},
						cisGLBHealthCheckDesc: {
							Type:        schema.TypeString,
							Description: "description",
							Computed:    true,
						},
						cisGLBHealthCheckType: {
							Type:        schema.TypeString,
							Description: "type",
							Computed:    true,
						},
						cisGLBHealthCheckMethod: {
							Type:        schema.TypeString,
							Description: "method",
							Computed:    true,
						},
						cisGLBHealthCheckTimeout: {
							Type:        schema.TypeInt,
							Description: "timeout",
							Computed:    true,
						},
						cisGLBHealthCheckRetries: {
							Type:        schema.TypeInt,
							Description: "retries",
							Computed:    true,
						},
						cisGLBHealthCheckInterval: {
							Type:        schema.TypeInt,
							Description: "interval",
							Computed:    true,
						},
						cisGLBHealthCheckFollowRedirects: {
							Type:        schema.TypeBool,
							Description: "follow_redirects",
							Computed:    true,
						},
						cisGLBHealthCheckAllowInsecure: {
							Type:        schema.TypeBool,
							Description: "allow_insecure",
							Computed:    true,
						},
						cisGLBHealthCheckCreatedOn: {
							Type:     schema.TypeString,
							Computed: true,
						},
						cisGLBHealthCheckModifiedOn: {
							Type:     schema.TypeString,
							Computed: true,
						},
						cisGLBHealthCheckPort: {
							Type:     schema.TypeInt,
							Computed: true,
						},
						cisGLBHealthCheckHeaders: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									cisGLBHealthCheckHeadersHeader: {
										Type:     schema.TypeString,
										Computed: true,
									},
									cisGLBHealthCheckHeadersValues: {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMCISGLBHealthCheckRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).CisGLBHealthCheckClientSession()
	if err != nil {
		return err
	}

	crn := d.Get(cisID).(string)
	sess.Crn = core.StringPtr(crn)

	opt := sess.NewListAllLoadBalancerMonitorsOptions()

	result, resp, err := sess.ListAllLoadBalancerMonitors(opt)
	if err != nil {
		log.Printf("Error listing global load balancer health check detail: %s", resp)
		return err
	}

	monitors := make([]map[string]interface{}, 0)
	for _, instance := range result.Result {
		monitor := map[string]interface{}{}
		monitor["id"] = convertCisToTfTwoVar(*instance.ID, crn)
		monitor[cisID] = crn
		monitor[cisGLBHealthCheckID] = *instance.ID
		monitor[cisGLBHealthCheckDesc] = *instance.Description
		monitor[cisGLBHealthCheckPath] = *instance.Path
		monitor[cisGLBHealthCheckCreatedOn] = *instance.CreatedOn
		monitor[cisGLBHealthCheckModifiedOn] = *instance.ModifiedOn
		monitor[cisGLBHealthCheckExpectedBody] = *instance.ExpectedBody
		monitor[cisGLBHealthCheckExpectedCodes] = *instance.ExpectedCodes
		monitor[cisGLBHealthCheckType] = *instance.Type
		monitor[cisGLBHealthCheckMethod] = *instance.Method
		monitor[cisGLBHealthCheckTimeout] = *instance.Timeout
		monitor[cisGLBHealthCheckRetries] = *instance.Retries
		monitor[cisGLBHealthCheckInterval] = *instance.Interval
		monitor[cisGLBHealthCheckFollowRedirects] = *instance.FollowRedirects
		monitor[cisGLBHealthCheckAllowInsecure] = *instance.AllowInsecure
		monitor[cisGLBHealthCheckHeaders] = flattenDataSourceLoadBalancerMonitorHeader(instance.Header)
		if instance.Port != nil {
			monitor[cisGLBHealthCheckPort] = *instance.Port
		}

		monitors = append(monitors, monitor)
	}
	d.SetId(dataSourceIBMCISGLBHealthCheckID(d))
	d.Set(cisID, crn)
	d.Set(cisGLBHealthCheck, monitors)
	return nil
}

// dataSourceIBMCISDNSRecordID returns a reasonable ID for dns zones list.
func dataSourceIBMCISGLBHealthCheckID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func flattenDataSourceLoadBalancerMonitorHeader(header map[string][]string) interface{} {
	flattened := make([]interface{}, 0)
	for k, v := range header {
		cfg := map[string]interface{}{
			cisGLBHealthCheckHeadersHeader: k,
			cisGLBHealthCheckHeadersValues: v,
		}
		flattened = append(flattened, cfg)
	}
	return flattened
}
