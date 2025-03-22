// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package dnsservices

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMPrivateDNSCustomResolver() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMDNSCustomResolverRead,

		Schema: map[string]*schema.Schema{
			pdnsInstanceID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance ID",
			},
			pdnsCustomResolvers: {
				Type:        schema.TypeList,
				Description: "Custom resolver details",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						pdnsCRId: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identifier of the custom resolver",
						},
						pdnsCRName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the custom resolver",
						},
						pdnsCRDescription: {
							Type:     schema.TypeString,
							Computed: true,
						},
						pdnsCREnabled: {
							Type:     schema.TypeBool,
							Computed: true,
						},
						pdnsCRHealth: {
							Type:     schema.TypeString,
							Computed: true,
						},
						pdnsCRProfile: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The profile name of the custom resolver.",
						},
						pdnsCRAllowDisruptiveUpdates: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether a disruptive update is allowed for the custom resolver",
						},
						pdnsCustomResolverLocations: {
							Type:        schema.TypeList,
							Description: "Locations on which the custom resolver will be running",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									pdnsCRLocationId: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Identifier of the custom resolver",
									},
									pdnsCRLocationSubnetCrn: {
										Type:     schema.TypeString,
										Computed: true,
									},
									pdnsCRLocationEnabled: {
										Type:     schema.TypeBool,
										Computed: true,
									},
									pdnsCRLocationHealthy: {
										Type:     schema.TypeBool,
										Computed: true,
									},
									pdnsCRLocationDnsServerIp: {
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
	}
}

func dataSourceIBMDNSCustomResolverRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).PrivateDNSClientSession()
	if err != nil {
		return diag.FromErr(err)
	}
	instanceID := d.Get(pdnsInstanceID).(string)

	opt := sess.NewListCustomResolversOptions(instanceID)
	result, resp, err := sess.ListCustomResolversWithContext(context, opt)
	if err != nil || result == nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error listing the custom resolvers %s:%s", err, resp))
	}

	customResolvers := make([]interface{}, 0)
	for _, instance := range result.CustomResolvers {
		customResolver := map[string]interface{}{}
		customResolver[pdnsCRId] = *instance.ID
		customResolver[pdnsCRName] = *instance.Name
		customResolver[pdnsCRDescription] = *instance.Description
		customResolver[pdnsCRHealth] = *instance.Health
		customResolver[pdnsCREnabled] = *instance.Enabled
		customResolver[pdnsCRProfile] = *instance.Profile
		customResolver[pdnsCRAllowDisruptiveUpdates] = *instance.AllowDisruptiveUpdates
		customResolver[pdnsCustomResolverLocations] = flattenPdnsCRLocations(instance.Locations)

		customResolvers = append(customResolvers, customResolver)
	}
	d.SetId(dataSourceIBMPrivateDNSCustomResolverID())
	d.Set(pdnsInstanceID, instanceID)
	d.Set(pdnsCustomResolvers, customResolvers)
	return nil
}

// dataSourceIBMPrivateDNSCustomResolverID returns a reasonable ID for dns  custom resolver list.
func dataSourceIBMPrivateDNSCustomResolverID() string {
	return time.Now().UTC().String()
}
