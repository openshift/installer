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
	cisOriginPools = "cis_origin_pools"
)

func dataSourceIBMCISOriginPools() *schema.Resource {
	return &schema.Resource{
		Read:     dataSourceIBMCISGLBPoolsRead,
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
			cisOriginPools: {
				Type:        schema.TypeList,
				Description: "Collection of GLB pools detail",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "GLB Pools id",
						},
						cisGLBPoolID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "GLB Pool id",
						},
						cisGLBPoolName: {
							Type:        schema.TypeString,
							Description: "name",
							Computed:    true,
						},
						cisGLBPoolRegions: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "List of regions",
						},
						cisGLBPoolDesc: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the CIS Origin Pool",
						},
						cisGLBPoolEnabled: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Boolean value set to true if cis origin pool needs to be enabled",
						},
						cisGLBPoolMinimumOrigins: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Minimum number of Origins",
						},
						cisGLBPoolMonitor: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Monitor value",
						},
						cisGLBPoolNotificationEMail: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Email address configured to recieve the notifications",
						},
						cisGLBPoolOrigins: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Origins info",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									cisGLBPoolOriginsName: {
										Type:     schema.TypeString,
										Computed: true,
									},
									cisGLBPoolOriginsAddress: {
										Type:     schema.TypeString,
										Computed: true,
									},
									cisGLBPoolOriginsEnabled: {
										Type:     schema.TypeBool,
										Computed: true,
									},
									cisGLBPoolOriginsWeight: {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									cisGLBPoolOriginsHealthy: {
										Type:     schema.TypeBool,
										Computed: true,
									},
									cisGLBPoolOriginsDisabledAt: {
										Type:     schema.TypeString,
										Computed: true,
									},
									cisGLBPoolOriginsFailureReason: {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						cisGLBPoolHealthy: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Health status",
						},
						cisGLBPoolCreatedOn: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation date info",
						},
						cisGLBPoolModifiedOn: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Modified date info",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMCISGLBPoolsRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisGLBPoolClientSession()
	if err != nil {
		return err
	}

	crn := d.Get(cisID).(string)
	cisClient.Crn = core.StringPtr(crn)

	opt := cisClient.NewListAllLoadBalancerPoolsOptions()
	result, resp, err := cisClient.ListAllLoadBalancerPools(opt)
	if err != nil {
		log.Printf("Error listing global load balancer pools detail: %s", resp)
		return err
	}

	pools := make([]map[string]interface{}, 0)
	for _, instance := range result.Result {
		pool := map[string]interface{}{}
		pool["id"] = convertCisToTfTwoVar(*instance.ID, crn)
		pool[cisGLBPoolID] = *instance.ID
		pool[cisGLBPoolName] = *instance.Name
		pool[cisGLBPoolOrigins] = flattenOrigins(instance.Origins)
		pool[cisGLBPoolRegions] = instance.CheckRegions
		pool[cisGLBPoolDesc] = *instance.Description
		pool[cisGLBPoolEnabled] = *instance.Enabled
		pool[cisGLBPoolNotificationEMail] = *instance.NotificationEmail
		pool[cisGLBPoolCreatedOn] = *instance.CreatedOn
		pool[cisGLBPoolModifiedOn] = *instance.ModifiedOn
		if instance.Monitor != nil {
			pool[cisGLBPoolMonitor] = *instance.Monitor
		}
		if instance.Healthy != nil {
			pool[cisGLBPoolHealthy] = *instance.Healthy
		}

		pools = append(pools, pool)
	}
	d.SetId(dataSourceIBMCISGLBPoolsID(d))
	d.Set(cisID, crn)
	d.Set(cisOriginPools, pools)
	return nil
}

func dataSourceIBMCISGLBPoolsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
