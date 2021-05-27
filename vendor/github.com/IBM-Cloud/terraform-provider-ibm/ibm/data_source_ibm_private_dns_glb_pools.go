// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	pdnsGLBPools = "dns_glb_pools"
)

func dataSourceIBMPrivateDNSGLBPools() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMPrivateDNSGLBPoolsRead,
		Schema: map[string]*schema.Schema{
			pdnsInstanceID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance ID",
			},
			pdnsGLBPools: {
				Type:        schema.TypeList,
				Description: "Collection of dns resource records",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						pdnsGlbPoolID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS record id",
						},
						pdnsGlbPoolName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS record name",
						},
						pdnsGlbPoolDescription: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Descriptive text of the load balancer pool",
						},
						pdnsGlbPoolEnabled: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the load balancer pool is enabled",
						},
						pdnsGlbPoolHealthyOriginsThreshold: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The minimum number of origins that must be healthy for this pool to serve traffic",
						},
						pdnsGlbPoolCreatedOn: {
							Type:        schema.TypeString,
							Description: "The time when a load balancer pool is created.",
							Computed:    true,
						},
						pdnsGlbPoolModifiedOn: {
							Type:        schema.TypeString,
							Description: "The recent time when a load balancer pool is modified.",
							Computed:    true,
						},
						pdnsGlbPoolMonitor: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the load balancer monitor to be associated to this pool",
						},
						pdnsGlbPoolChannel: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The notification channel,It is a webhook url",
						},
						pdnsGlbPoolRegion: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Health check region of VSIs",
						},
						pdnsGlbPoolHealth: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether the load balancer pool is enabled",
						},
						pdnsGlbPoolOrigins: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Origins info",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									pdnsGlbPoolOriginsName: {
										Type:        schema.TypeString,
										Description: "The name of the origin server.",
										Computed:    true,
									},
									pdnsGlbPoolOriginsAddress: {
										Type:        schema.TypeString,
										Description: "The address of the origin server. It can be a hostname or an IP address.",
										Computed:    true,
									},
									pdnsGlbPoolOriginsEnabled: {
										Type:        schema.TypeBool,
										Description: "Whether the origin server is enabled.",
										Computed:    true,
									},
									pdnsGlbPoolOriginsDescription: {
										Type:        schema.TypeString,
										Description: "Description of the origin server.",
										Computed:    true,
									},
									pdnsGlbPoolOriginsHealth: {
										Type:        schema.TypeBool,
										Description: "Whether the health is `true` or `false`.",
										Computed:    true,
									},
									pdnsGlbPoolOriginsHealthFailureReason: {
										Type:        schema.TypeString,
										Description: "The Reason for health check failure",
										Computed:    true,
									},
								},
							},
						},
						pdnsGlbPoolSubnet: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Health check subnet crn of VSIs",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMPrivateDNSGLBPoolsRead(d *schema.ResourceData, meta interface{}) error {

	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}
	instanceID := d.Get(pdnsInstanceID).(string)
	listDNSGLBPooloptions := sess.NewListPoolsOptions(instanceID)
	availableGLBPools, detail, err := sess.ListPools(listDNSGLBPooloptions)
	if err != nil {
		return fmt.Errorf("Error reading list of pdns GLB pools:%s\n%s", err, detail)
	}
	d.Set(pdnsInstanceID, instanceID)
	dnsPools := make([]map[string]interface{}, 0)
	for _, instance := range availableGLBPools.Pools {
		dnsPool := map[string]interface{}{}
		dnsPool[pdnsGlbPoolID] = *instance.ID
		dnsPool[pdnsGlbPoolName] = *instance.Name
		dnsPool[pdnsGlbPoolDescription] = *instance.Description
		dnsPool[pdnsGlbPoolEnabled] = *instance.Enabled
		dnsPool[pdnsGlbPoolHealth] = *instance.Health
		dnsPool[pdnsGlbPoolHealthyOriginsThreshold] = *instance.HealthyOriginsThreshold
		dnsPool[pdnsGlbPoolCreatedOn] = *instance.CreatedOn
		dnsPool[pdnsGlbPoolModifiedOn] = *instance.ModifiedOn
		dnsPool[pdnsGlbPoolMonitor] = *instance.Monitor
		dnsPool[pdnsGlbPoolChannel] = *instance.NotificationChannel
		dnsPool[pdnsGlbPoolRegion] = *instance.HealthcheckRegion
		dnsPool[pdnsGlbPoolOrigins] = flattenPDNSGlbPoolOrigins(instance.Origins)
		dnsPool[pdnsGlbPoolSubnet] = instance.HealthcheckSubnets

		dnsPools = append(dnsPools, dnsPool)
	}
	d.SetId(dataSourceIBMPrivateDNSGLBPoolsID(d))
	d.Set(pdnsGLBPools, dnsPools)
	return nil
}

// dataSourceIBMPrivateDNSGLBMonitorsID returns a reasonable ID  list.
func dataSourceIBMPrivateDNSGLBPoolsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
