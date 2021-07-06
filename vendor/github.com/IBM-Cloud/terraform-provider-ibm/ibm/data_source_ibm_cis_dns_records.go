// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	cisDNSRecords           = "cis_dns_records"
	cisDNSRecordsExportFile = "file"
)

func dataSourceIBMCISDNSRecords() *schema.Resource {
	return &schema.Resource{
		Read:     dataSourceIBMCISDNSRecordsRead,
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
			cisDomainID: {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Zone Id",
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisDNSRecordsExportFile: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "file to be exported",
			},

			cisDNSRecords: {
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
						cisDNSRecordID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS record id",
						},
						cisZoneName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS Record Name",
						},
						cisDNSRecordName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS Record Name",
						},
						cisDNSRecordCreatedOn: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS Record created on",
						},
						cisDNSRecordModifiedOn: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS Record modified on",
						},
						cisDNSRecordType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS Record Type",
						},
						cisDNSRecordContent: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS Record conent info",
						},
						cisDNSRecordPriority: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "DNS Record MX priority",
						},
						cisDNSRecordProxiable: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "DNS Record proxiable",
						},
						cisDNSRecordProxied: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "DNS Record proxied",
						},
						cisDNSRecordTTL: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "DNS Record Time To Live",
						},
						cisDNSRecordData: {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "DNS Record Data",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMCISDNSRecordsRead(d *schema.ResourceData, meta interface{}) error {
	var (
		crn     string
		zoneID  string
		records []map[string]interface{}
	)
	sess, err := meta.(ClientSession).CisDNSRecordClientSession()
	if err != nil {
		return err
	}

	// session options
	crn = d.Get(cisID).(string)
	zoneID, _, _ = convertTftoCisTwoVar(d.Get(cisDomainID).(string))
	sess.Crn = core.StringPtr(crn)
	sess.ZoneIdentifier = core.StringPtr(zoneID)

	if file, ok := d.GetOk(cisDNSRecordsExportFile); ok {
		sess, err := meta.(ClientSession).CisDNSRecordBulkClientSession()
		if err != nil {
			return err
		}
		sess.Crn = core.StringPtr(crn)
		sess.ZoneIdentifier = core.StringPtr(zoneID)
		opt := sess.NewGetDnsRecordsBulkOptions()
		result, response, err := sess.GetDnsRecordsBulk(opt)
		if err != nil {
			log.Printf("Error exporting dns records: %s", response)
			return err
		}
		buf, err := ioutil.ReadAll(result)
		if err != nil {
			log.Printf("Error while reading io reader")
			return err
		}

		f, err := os.Create(file.(string))
		if err != nil {
			log.Printf("Error opening file: %v", err)
			return err
		}
		defer f.Close()
		f.Write(buf)
		d.Set(cisDNSRecordsExportFile, file)
	}

	opt := sess.NewListAllDnsRecordsOptions()
	opt.SetPage(1)
	opt.SetPerPage(1000)
	result, response, err := sess.ListAllDnsRecords(opt)
	if err != nil {
		log.Printf("Error reading dns records: %s", response)
		return err
	}

	records = make([]map[string]interface{}, 0)
	for _, instance := range result.Result {
		record := map[string]interface{}{}
		record["id"] = convertCisToTfThreeVar(*instance.ID, zoneID, crn)
		record[cisDNSRecordID] = *instance.ID
		record[cisZoneName] = *instance.ZoneName
		record[cisDNSRecordCreatedOn] = *instance.CreatedOn
		record[cisDNSRecordModifiedOn] = *instance.ModifiedOn
		record[cisDNSRecordName] = *instance.Name
		record[cisDNSRecordType] = *instance.Type
		if instance.Priority != nil {
			record[cisDNSRecordPriority] = *instance.Priority
		}
		if instance.Content != nil {
			record[cisDNSRecordContent] = *instance.Content
		}
		record[cisDNSRecordProxiable] = *instance.Proxiable
		record[cisDNSRecordProxied] = *instance.Proxied
		record[cisDNSRecordTTL] = *instance.TTL
		if instance.Data != nil {
			d.Set(cisDNSRecordData, flattenData(instance.Data, *instance.ZoneName))
		}

		records = append(records, record)
	}
	d.SetId(dataSourceIBMCISDNSRecordID(d))
	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	d.Set(cisDNSRecords, records)
	return nil
}

// dataSourceIBMCISDNSRecordID returns a reasonable ID for dns zones list.
func dataSourceIBMCISDNSRecordID(d *schema.ResourceData) string {
	zoneID := d.Get(cisDomainID)
	crn := d.Get(cisID)
	return fmt.Sprintf("%s:%s", zoneID, crn)
}
