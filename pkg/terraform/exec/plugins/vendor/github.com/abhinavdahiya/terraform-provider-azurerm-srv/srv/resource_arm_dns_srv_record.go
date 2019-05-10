package srv

import (
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/services/preview/dns/mgmt/2018-03-01-preview/dns"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArmDnsSrvRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDnsSrvRecordCreateUpdate,
		Read:   resourceArmDnsSrvRecordRead,
		Update: resourceArmDnsSrvRecordCreateUpdate,
		Delete: resourceArmDnsSrvRecordDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": schemaResourceGroupName(),

			"zone_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"record": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"priority": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},

						"weight": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},

						"port": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},

						"target": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},

			"ttl": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmDnsSrvRecordCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*armClient).dnsClient
	ctx := meta.(*armClient).StopContext

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	zoneName := d.Get("zone_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, zoneName, name, dns.SRV)
		if err != nil {
			if !responseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing DNS SRV Record %q (Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return fmt.Errorf("resource azurerm_dns_srv_record %q already exists, cannot create new", *existing.ID)
		}
	}

	ttl := int64(d.Get("ttl").(int))
	tags := d.Get("tags").(map[string]interface{})

	parameters := dns.RecordSet{
		Name: &name,
		RecordSetProperties: &dns.RecordSetProperties{
			Metadata:   expandTags(tags),
			TTL:        &ttl,
			SrvRecords: expandAzureRmDnsSrvRecords(d),
		},
	}

	eTag := ""
	ifNoneMatch := "" // set to empty to allow updates to records after creation
	if _, err := client.CreateOrUpdate(ctx, resGroup, zoneName, name, dns.SRV, parameters, eTag, ifNoneMatch); err != nil {
		return fmt.Errorf("Error creating/updating DNS SRV Record %q (Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
	}

	resp, err := client.Get(ctx, resGroup, zoneName, name, dns.SRV)
	if err != nil {
		return fmt.Errorf("Error retrieving DNS SRV Record %q (Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read DNS SRV Record %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmDnsSrvRecordRead(d, meta)
}

func resourceArmDnsSrvRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*armClient).dnsClient
	ctx := meta.(*armClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["SRV"]
	zoneName := id.Path["dnszones"]

	resp, err := client.Get(ctx, resGroup, zoneName, name, dns.SRV)
	if err != nil {
		if responseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading DNS SRV record %s: %v", name, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("zone_name", zoneName)
	d.Set("ttl", resp.TTL)

	if err := d.Set("record", flattenAzureRmDnsSrvRecords(resp.SrvRecords)); err != nil {
		return err
	}
	flattenAndSetTags(d, resp.Metadata)

	return nil
}

func resourceArmDnsSrvRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*armClient).dnsClient
	ctx := meta.(*armClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["SRV"]
	zoneName := id.Path["dnszones"]

	resp, err := client.Delete(ctx, resGroup, zoneName, name, dns.SRV, "")
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error deleting DNS SRV Record %s: %+v", name, err)
	}

	return nil
}

func flattenAzureRmDnsSrvRecords(records *[]dns.SrvRecord) []map[string]interface{} {
	results := make([]map[string]interface{}, 0, len(*records))

	if records != nil {
		for _, record := range *records {
			results = append(results, map[string]interface{}{
				"priority": *record.Priority,
				"weight":   *record.Weight,
				"port":     *record.Port,
				"target":   *record.Target,
			})
		}
	}

	return results
}

func expandAzureRmDnsSrvRecords(d *schema.ResourceData) *[]dns.SrvRecord {
	recordStrings := d.Get("record").([]interface{})
	records := make([]dns.SrvRecord, len(recordStrings))

	for i, v := range recordStrings {
		record := v.(map[string]interface{})
		priority := int32(record["priority"].(int))
		weight := int32(record["weight"].(int))
		port := int32(record["port"].(int))
		target := record["target"].(string)

		srvRecord := dns.SrvRecord{
			Priority: &priority,
			Weight:   &weight,
			Port:     &port,
			Target:   &target,
		}

		records[i] = srvRecord
	}

	return &records
}
