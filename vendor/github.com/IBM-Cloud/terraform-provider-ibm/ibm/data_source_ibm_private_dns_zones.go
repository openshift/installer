// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	pdnsZones = "dns_zones"
)

func dataSourceIBMPrivateDNSZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMPrivateDNSZonesRead,
		Schema: map[string]*schema.Schema{
			pdnsInstanceID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance ID",
			},
			pdnsZones: {
				Type:        schema.TypeList,
				Description: "Collection of dns zones",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						pdnsInstanceID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance ID",
						},
						pdnsZoneID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Zone ID",
						},
						pdnsZoneName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Zone name",
						},
						pdnsZoneDescription: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Zone description",
						},
						pdnsZoneState: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Zone state",
						},
						pdnsZoneLabel: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Label",
						},
						pdnsZoneCreatedOn: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation date",
						},
						pdnsZoneModifiedOn: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Modification date",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMPrivateDNSZonesRead(d *schema.ResourceData, meta interface{}) error {

	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}
	instanceID := d.Get(pdnsInstanceID).(string)
	listDNSZonesOptions := sess.NewListDnszonesOptions(instanceID)
	availableDNSZones, detail, err := sess.ListDnszones(listDNSZonesOptions)
	if err != nil {
		return fmt.Errorf("Error reading list of dns zones:%s\n%s", err, detail)
	}
	dnsZones := make([]map[string]interface{}, 0)
	for _, instance := range availableDNSZones.Dnszones {
		dnsZone := map[string]interface{}{}
		dnsZone[pdnsInstanceID] = instance.InstanceID
		dnsZone[pdnsZoneID] = instance.ID
		dnsZone[pdnsZoneName] = instance.Name
		dnsZone[pdnsZoneDescription] = instance.Description
		dnsZone[pdnsZoneLabel] = instance.Label
		dnsZone[pdnsZoneCreatedOn] = instance.CreatedOn
		dnsZone[pdnsZoneModifiedOn] = instance.ModifiedOn
		dnsZone[pdnsZoneState] = instance.State
		dnsZones = append(dnsZones, dnsZone)
	}
	d.SetId(dataSourceIBMPrivateDNSZonesID(d))
	d.Set(pdnsZones, dnsZones)
	return nil
}

// dataSourceIBMPrivateDnsZonesID returns a reasonable ID for dns zones list.
func dataSourceIBMPrivateDNSZonesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
