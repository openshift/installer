package azureprivatedns

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/services/privatedns/mgmt/2018-09-01/privatedns"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPrivateDNSSrvRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPrivateDNSSrvRecordCreateUpdate,
		Read:   resourceArmPrivateDNSSrvRecordRead,
		Update: resourceArmPrivateDNSSrvRecordCreateUpdate,
		Delete: resourceArmPrivateDNSArvRecordDelete,
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

			"record": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"priority": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"weight": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"port": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"target": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Set: resourceArmPrivateDNSSrvRecordHash,
			},

			"ttl": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmPrivateDNSSrvRecordCreateUpdate(d *schema.ResourceData, meta interface{}) error {
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
			Metadata:   expandTags(t),
			TTL:        &ttl,
			SrvRecords: expandAzureRmPrivateDNSSrvRecords(d),
		},
	}

	eTag := ""
	ifNoneMatch := "" // set to empty to allow updates to records after creation
	if _, err := client.CreateOrUpdate(ctx, resGroup, zoneName, privatedns.SRV, name, parameters, eTag, ifNoneMatch); err != nil {
		return fmt.Errorf("error creating/updating Private DNS SRV Record %q (Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
	}

	resp, err := client.Get(ctx, resGroup, zoneName, privatedns.SRV, name)
	if err != nil {
		return fmt.Errorf("error retrieving Private DNS SRV Record %q (Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("cannot read Private DNS SRV Record %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmPrivateDNSSrvRecordRead(d, meta)
}

func resourceArmPrivateDNSSrvRecordRead(d *schema.ResourceData, meta interface{}) error {
	dnsClient := meta.(*ArmClient).recordSetsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["SRV"]
	zoneName := id.Path["privateDnsZones"]

	resp, err := dnsClient.Get(ctx, resGroup, zoneName, privatedns.SRV, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error reading Private DNS SRV record %s: %+v", name, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("zone_name", zoneName)
	d.Set("ttl", resp.TTL)

	if err := d.Set("record", flattenAzureRmPrivateDNSSrvRecords(resp.SrvRecords)); err != nil {
		return err
	}
	flattenAndSetTags(d, resp.Metadata)

	return nil
}

func resourceArmPrivateDNSArvRecordDelete(d *schema.ResourceData, meta interface{}) error {
	dnsClient := meta.(*ArmClient).recordSetsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["A"]
	zoneName := id.Path["privateDnsZones"]

	resp, err := dnsClient.Delete(ctx, resGroup, zoneName, privatedns.SRV, name, "")
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error deleting Private DNS A Record %s: %+v", name, err)
	}

	return nil
}

func flattenAzureRmPrivateDNSSrvRecords(records *[]privatedns.SrvRecord) []map[string]interface{} {
	if records == nil {
		return nil
	}
	results := make([]map[string]interface{}, 0, len(*records))

	for _, record := range *records {
		results = append(results, map[string]interface{}{
			"priority": *record.Priority,
			"weight":   *record.Weight,
			"port":     *record.Port,
			"target":   *record.Target,
		})
	}

	return results
}

func expandAzureRmPrivateDNSSrvRecords(d *schema.ResourceData) *[]privatedns.SrvRecord {
	recordStrings := d.Get("record").(*schema.Set).List()
	records := make([]privatedns.SrvRecord, len(recordStrings))

	for i, v := range recordStrings {
		record := v.(map[string]interface{})
		priority := int32(record["priority"].(int))
		weight := int32(record["weight"].(int))
		port := int32(record["port"].(int))
		target := record["target"].(string)

		srvRecord := privatedns.SrvRecord{
			Priority: &priority,
			Weight:   &weight,
			Port:     &port,
			Target:   &target,
		}

		records[i] = srvRecord
	}

	return &records
}

func resourceArmPrivateDNSSrvRecordHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%d-", m["priority"].(int)))
		buf.WriteString(fmt.Sprintf("%d-", m["weight"].(int)))
		buf.WriteString(fmt.Sprintf("%d-", m["port"].(int)))
		buf.WriteString(fmt.Sprintf("%s-", m["target"].(string)))
	}

	return hashcode.String(buf.String())
}
