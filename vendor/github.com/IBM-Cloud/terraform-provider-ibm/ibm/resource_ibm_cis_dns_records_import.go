// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	cisDNSRecordsImportFile               = "file"
	cisDNSRecordsImportTotalRecordsParsed = "total_records_parsed"
	cisDNSRecordsImportRecordsAdded       = "records_added"
)

func resourceIBMCISDNSRecordsImport() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Description:      "Associated CIS domain",
				Required:         true,
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisDNSRecordsImportFile: {
				Type:        schema.TypeString,
				Description: "File to import",
				Required:    true,
				ForceNew:    true,
			},
			cisDNSRecordsImportTotalRecordsParsed: {
				Type:        schema.TypeInt,
				Description: "total records parsed",
				Computed:    true,
			},
			cisDNSRecordsImportRecordsAdded: {
				Type:        schema.TypeInt,
				Description: "added records count",
				Computed:    true,
			},
		},

		Create:   resourceCISDNSRecordsImportUpdate,
		Read:     resourceCISDNSRecordsImportRead,
		Update:   resourceCISDNSRecordsImportRead,
		Delete:   resourceCISDNSRecordsImportDelete,
		Importer: &schema.ResourceImporter{},
	}
}
func resourceCISDNSRecordsImportUpdate(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisDNSRecordBulkClientSession()
	if err != nil {
		return err
	}

	crn := d.Get(cisID).(string)
	zoneID, _, err := convertTftoCisTwoVar(d.Get(cisDomainID).(string))
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)
	file := d.Get(cisDNSRecordsImportFile).(string)

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	opt := cisClient.NewPostDnsRecordsBulkOptions()
	opt.SetFile(f)
	result, response, err := cisClient.PostDnsRecordsBulk(opt)
	if err != nil {
		log.Printf("Error importing dns records: %v", response)
		return err
	}
	id := fmt.Sprintf("%v:%v:%s:%s:%s", *result.Result.TotalRecordsParsed,
		*result.Result.RecsAdded, file, zoneID, crn)
	d.SetId(id)

	return nil

}

func resourceCISDNSRecordsImportRead(d *schema.ResourceData, meta interface{}) error {
	idSplitStr := strings.SplitN(d.Id(), ":", 5)
	parsed, _ := strconv.Atoi(idSplitStr[0])
	added, _ := strconv.Atoi(idSplitStr[1])
	file := idSplitStr[2]
	zoneID := idSplitStr[3]
	crn := idSplitStr[4]
	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	d.Set(cisDNSRecordsImportFile, file)
	d.Set(cisDNSRecordsImportTotalRecordsParsed, parsed)
	d.Set(cisDNSRecordsImportRecordsAdded, added)
	return nil
}

func resourceCISDNSRecordsImportDelete(d *schema.ResourceData, meta interface{}) error {
	// Nothing to delete on CIS DNS Record import resource
	d.SetId("")
	return nil
}
