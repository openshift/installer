// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	pdnsResourceRecords = "dns_resource_records"
)

func dataSourceIBMPrivateDNSResourceRecords() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMPrivateDNSResourceRecordsRead,
		Schema: map[string]*schema.Schema{
			pdnsInstanceID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance ID",
			},
			pdnsZoneID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Zone Id",
			},
			pdnsResourceRecords: {
				Type:        schema.TypeList,
				Description: "Collection of dns resource records",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS record id",
						},
						pdnsRecordName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS record name",
						},
						pdnsRecordType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS record Type",
						},
						pdnsRdata: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS record Data",
						},
						pdnsRecordTTL: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "DNS record TTL",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMPrivateDNSResourceRecordsRead(d *schema.ResourceData, meta interface{}) error {

	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}
	instanceID := d.Get(pdnsInstanceID).(string)
	DnszoneID := d.Get(pdnsZoneID).(string)
	listDNSResRecOptions := sess.NewListResourceRecordsOptions(instanceID, DnszoneID)
	availableDNSResRecs, detail, err := sess.ListResourceRecords(listDNSResRecOptions)
	if err != nil {
		return fmt.Errorf("Error reading list of pdns resource records:%s\n%s", err, detail)
	}
	dnsResRecs := make([]map[string]interface{}, 0)
	for _, instance := range availableDNSResRecs.ResourceRecords {
		dnsRecord := map[string]interface{}{}
		dnsRecord["id"] = *instance.ID
		dnsRecord[pdnsRecordName] = *instance.Name
		dnsRecord[pdnsRecordType] = *instance.Type
		// Marshal the rdata map into a JSON string
		rData, err := json.Marshal(instance.Rdata)
		if err != nil {
			return fmt.Errorf("Error reading rdata map of dns resource records:%s", err)
		}
		jsonStr := string(rData)
		dnsRecord[pdnsRdata] = jsonStr
		dnsRecord[pdnsRecordTTL] = instance.TTL
		dnsResRecs = append(dnsResRecs, dnsRecord)
	}
	d.SetId(dataSourceIBMPrivateDNSResourceRecordsID(d))
	d.Set(pdnsResourceRecords, dnsResRecs)
	return nil
}

// dataSourceIBMPrivateDNSResourceRecordsID returns a reasonable ID for dns zones list.
func dataSourceIBMPrivateDNSResourceRecordsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
