// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	pdnsGLBs = "dns_glbs"
)

func dataSourceIBMPrivateDNSGLBs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMPrivateDNSGLBsRead,

		Schema: map[string]*schema.Schema{

			pdnsInstanceID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The GUID of the private DNS. ",
			},
			pdnsZoneID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Zone GUID ",
			},
			pdnsGLBs: {
				Type:        schema.TypeList,
				Description: "Collection of GLB load balancer collectors",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						pdnsGLBID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Load balancer Id",
						},
						pdnsGLBName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the load balancer",
						},
						pdnsGLBDescription: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Descriptive text of the load balancer",
						},
						pdnsGLBEnabled: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the load balancer is enabled",
						},
						pdnsGLBTTL: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Time to live in second",
						},
						pdnsGLBHealth: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Healthy state of the load balancer.",
						},
						pdnsGLBFallbackPool: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The pool ID to use when all other pools are detected as unhealthy",
						},
						pdnsGLBDefaultPool: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A list of pool IDs ordered by their failover priority",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						pdnsGLBAZPools: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Map availability zones to pool ID's.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									pdnsGLBAvailabilityZone: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Availability zone.",
									},

									pdnsGLBAZPoolsPools: {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "List of load balancer pools",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						pdnsGLBCreatedOn: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "GLB Load Balancer creation date",
						},
						pdnsGLBModifiedOn: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "GLB Load Balancer Modification date",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMPrivateDNSGLBsRead(d *schema.ResourceData, meta interface{}) error {

	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}
	instanceID := d.Get(pdnsInstanceID).(string)
	zoneID := d.Get(pdnsZoneID).(string)
	listDNSGLBs := sess.NewListLoadBalancersOptions(instanceID, zoneID)
	availableGLBs, detail, err := sess.ListLoadBalancers(listDNSGLBs)
	if err != nil {
		return fmt.Errorf("Error reading list of pdns GLB load balancers:%s\n%s", err, detail)
	}

	dnslbs := make([]interface{}, 0)
	for _, instance := range availableGLBs.LoadBalancers {
		dnsLoadbalancer := map[string]interface{}{}
		dnsLoadbalancer[pdnsGLBID] = *instance.ID
		dnsLoadbalancer[pdnsGLBName] = *instance.Name
		dnsLoadbalancer[pdnsGLBDescription] = *instance.Description
		dnsLoadbalancer[pdnsGLBEnabled] = *instance.Enabled
		dnsLoadbalancer[pdnsGLBTTL] = *instance.TTL
		dnsLoadbalancer[pdnsGLBHealth] = *instance.Health
		dnsLoadbalancer[pdnsGLBFallbackPool] = *instance.FallbackPool
		dnsLoadbalancer[pdnsGLBCreatedOn] = *instance.CreatedOn
		dnsLoadbalancer[pdnsGLBModifiedOn] = *instance.ModifiedOn
		dnsLoadbalancer[pdnsGLBDefaultPool] = instance.DefaultPools
		dnsLoadbalancer[pdnsGLBAZPools] = flattenPDNSGlbAZpool(instance.AzPools)

		dnslbs = append(dnslbs, dnsLoadbalancer)
	}
	d.SetId(dataSourceIBMPrivateDNSGLBsID(d))
	d.Set(pdnsInstanceID, instanceID)
	d.Set(pdnsZoneID, zoneID)
	d.Set(pdnsGLBs, dnslbs)
	return nil
}

// dataSourceIBMPrivateDNSGLBMonitorsID returns a reasonable ID  list.
func dataSourceIBMPrivateDNSGLBsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
