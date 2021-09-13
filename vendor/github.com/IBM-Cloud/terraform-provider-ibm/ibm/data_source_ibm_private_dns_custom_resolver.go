// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMDNSCustomResolver() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMDNSCustomResolverRead,

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

func dataSourceIBMDNSCustomResolverRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}
	instanceID := d.Get(pdnsInstanceID).(string)

	opt := sess.NewListCustomResolversOptions(instanceID)
	result, resp, err := sess.ListCustomResolversWithContext(context.TODO(), opt)
	if err != nil || result == nil {
		return fmt.Errorf("Error listing the custom resolvers %s:%s", err, resp)
	}

	customResolvers := make([]interface{}, 0)
	for _, instance := range result.CustomResolvers {
		customResolver := map[string]interface{}{}
		customResolver[pdnsCRId] = *instance.ID
		customResolver[pdnsCRName] = *instance.Name
		customResolver[pdnsCRDescription] = *instance.Description
		customResolver[pdnsCRHealth] = *instance.Health
		customResolver[pdnsCREnabled] = *instance.Enabled
		customResolver[pdnsCustomResolverLocations] = flattenPdnsCRLocations(instance.Locations)

		customResolvers = append(customResolvers, customResolver)
	}
	d.SetId(dataSourceIBMPrivateDNSCustomResolverID(d))
	d.Set(pdnsInstanceID, instanceID)
	d.Set(pdnsCustomResolvers, customResolvers)
	return nil
}

// dataSourceIBMPrivateDNSCustomResolverID returns a reasonable ID for dns  custom resolver list.
func dataSourceIBMPrivateDNSCustomResolverID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
