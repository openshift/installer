package azureprivatedns

import (
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/services/privatedns/mgmt/2018-09-01/privatedns"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPrivateDNSAAAARecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPrivateDNSAAAARecordCreateUpdate,
		Read:   resourceArmPrivateDNSAAAARecordRead,
		Update: resourceArmPrivateDNSAAAARecordCreateUpdate,
		Delete: resourceArmPrivateDNSAAAARecordDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"zone_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"records": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"ttl": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmPrivateDNSAAAARecordCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).recordSetsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	zoneName := d.Get("zone_name").(string)

	ttl := int64(d.Get("ttl").(int))
	t := d.Get("tags").(map[string]interface{})

	parameters := privatedns.RecordSet{
		Name: &name,
		RecordSetProperties: &privatedns.RecordSetProperties{
			Metadata:    expandTags(t),
			TTL:         &ttl,
			AaaaRecords: expandAzureRmPrivateDNSAAAARecords(d),
		},
	}

	eTag := ""
	ifNoneMatch := "" // set to empty to allow updates to records after creation
	if _, err := client.CreateOrUpdate(ctx, resGroup, zoneName, privatedns.AAAA, name, parameters, eTag, ifNoneMatch); err != nil {
		return fmt.Errorf("error creating/updating Private DNS AAAA Record %q (Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
	}

	resp, err := client.Get(ctx, resGroup, zoneName, privatedns.AAAA, name)
	if err != nil {
		return fmt.Errorf("error retrieving Private DNS AAAA Record %q (Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("cannot read Private DNS AAAA Record %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmPrivateDNSAAAARecordRead(d, meta)
}

func resourceArmPrivateDNSAAAARecordRead(d *schema.ResourceData, meta interface{}) error {
	dnsClient := meta.(*ArmClient).recordSetsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["AAAA"]
	zoneName := id.Path["privateDnsZones"]

	resp, err := dnsClient.Get(ctx, resGroup, zoneName, privatedns.AAAA, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error reading Private DNS A record %s: %+v", name, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("zone_name", zoneName)
	d.Set("ttl", resp.TTL)

	if err := d.Set("records", flattenAzureRmPrivateDNSAAAARecords(resp.AaaaRecords)); err != nil {
		return err
	}
	flattenAndSetTags(d, resp.Metadata)

	return nil
}

func resourceArmPrivateDNSAAAARecordDelete(d *schema.ResourceData, meta interface{}) error {
	dnsClient := meta.(*ArmClient).recordSetsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["AAAA"]
	zoneName := id.Path["privateDnsZones"]

	resp, err := dnsClient.Delete(ctx, resGroup, zoneName, privatedns.AAAA, name, "")
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error deleting Private DNS AAAA Record %s: %+v", name, err)
	}

	return nil
}

func flattenAzureRmPrivateDNSAAAARecords(records *[]privatedns.AaaaRecord) []string {
	results := make([]string, 0)
	if records == nil {
		return results
	}

	for _, record := range *records {
		if record.Ipv6Address == nil {
			continue
		}

		results = append(results, *record.Ipv6Address)
	}

	return results
}

func expandAzureRmPrivateDNSAAAARecords(d *schema.ResourceData) *[]privatedns.AaaaRecord {
	recordStrings := d.Get("records").(*schema.Set).List()
	records := make([]privatedns.AaaaRecord, len(recordStrings))

	for i, v := range recordStrings {
		ipv6 := v.(string)
		records[i] = privatedns.AaaaRecord{
			Ipv6Address: &ipv6,
		}
	}

	return &records
}
