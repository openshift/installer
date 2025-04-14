// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	cisDNSRecordsImportFile               = "file"
	cisDNSRecordsImportTotalRecordsParsed = "total_records_parsed"
	cisDNSRecordsImportRecordsAdded       = "records_added"
)

func ResourceIBMCISDNSRecordsImport() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
				ValidateFunc: validate.InvokeValidator("ibm_cis_dns_records_import",
					"cis_id"),
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
func ResourceIBMCISDnsRecordsImportValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cis_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "resource_instance",
			CloudDataRange:             []string{"service:internet-svcs"},
			Required:                   true})
	ibmCISDNSRecordsImportValidator := validate.ResourceValidator{
		ResourceName: "ibm_cis_dns_records_import",
		Schema:       validateSchema}
	return &ibmCISDNSRecordsImportValidator
}
func resourceCISDNSRecordsImportUpdate(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisDNSRecordBulkClientSession()
	if err != nil {
		return err
	}

	crn := d.Get(cisID).(string)
	zoneID, _, _ := flex.ConvertTftoCisTwoVar(d.Get(cisDomainID).(string))
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
