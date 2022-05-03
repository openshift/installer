// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIbmIsVpcAddressPrefixes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmIsVpcAddressPrefixRead,

		Schema: map[string]*schema.Schema{
			"vpc": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The VPC identifier.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The user-defined name for this address prefix. Names must be unique within the VPC the address prefix resides in.",
			},
			"address_prefixes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of address prefixes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CIDR block for this prefix.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the prefix was created.",
						},
						"has_subnets": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether subnets exist with addresses from this prefix.",
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this address prefix.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this address prefix.",
						},
						"is_default": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether this is the default prefix for this zone in this VPC. If a default prefix was automatically created when the VPC was created, the prefix is automatically named using a hyphenated list of randomly-selected words, but may be updated with a user-specified name.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this address prefix. Names must be unique within the VPC the address prefix resides in.",
						},
						"zone": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The zone this address prefix resides in.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"href": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this zone.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The globally unique name for this zone.",
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

func dataSourceIbmIsVpcAddressPrefixRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	start := ""
	allrecs := []vpcv1.AddressPrefix{}
	for {
		listVpcAddressPrefixesOptions := &vpcv1.ListVPCAddressPrefixesOptions{}

		listVpcAddressPrefixesOptions.SetVPCID(d.Get("vpc").(string))

		if start != "" {
			listVpcAddressPrefixesOptions.Start = &start
		}
		addressPrefixCollection, response, err := vpcClient.ListVPCAddressPrefixesWithContext(context, listVpcAddressPrefixesOptions)
		if err != nil {
			log.Printf("[DEBUG] ListVpcAddressPrefixesWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("ListVpcAddressPrefixesWithContext failed %s\n%s", err, response))
		}
		start = flex.GetNext(addressPrefixCollection.Next)
		allrecs = append(allrecs, addressPrefixCollection.AddressPrefixes...)
		if start == "" {
			break
		}
	}

	// Use the provided filter argument and construct a new list with only the requested resource(s)
	var matchAddressPrefixes []vpcv1.AddressPrefix
	var name string
	var suppliedFilter bool

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
		suppliedFilter = true
		for _, data := range allrecs {
			if *data.Name == name {
				matchAddressPrefixes = append(matchAddressPrefixes, data)
			}
		}
	} else {
		matchAddressPrefixes = allrecs
	}

	if suppliedFilter {
		if len(matchAddressPrefixes) == 0 {
			return diag.FromErr(fmt.Errorf("no AddressPrefixes found with name %s", name))
		}
		d.SetId(name)
	} else {
		d.SetId(dataSourceIbmIsVpcAddressPrefixID(d))
	}

	if matchAddressPrefixes != nil {
		err = d.Set("address_prefixes", dataSourceAddressPrefixCollectionFlattenAddressPrefixes(matchAddressPrefixes))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting address_prefixes %s", err))
		}
	}

	return nil
}

// dataSourceIbmIsVpcAddressPrefixID returns a reasonable ID for the list.
func dataSourceIbmIsVpcAddressPrefixID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceAddressPrefixCollectionFlattenAddressPrefixes(result []vpcv1.AddressPrefix) (addressPrefixes []map[string]interface{}) {
	for _, addressPrefixesItem := range result {
		addressPrefixes = append(addressPrefixes, dataSourceAddressPrefixCollectionAddressPrefixesToMap(addressPrefixesItem))
	}

	return addressPrefixes
}

func dataSourceAddressPrefixCollectionAddressPrefixesToMap(addressPrefixesItem vpcv1.AddressPrefix) (addressPrefixesMap map[string]interface{}) {

	addressPrefixesMap = map[string]interface{}{}

	if addressPrefixesItem.CIDR != nil {
		addressPrefixesMap["cidr"] = addressPrefixesItem.CIDR
	}
	if addressPrefixesItem.CreatedAt != nil {
		addressPrefixesMap["created_at"] = addressPrefixesItem.CreatedAt.String()
	}
	if addressPrefixesItem.HasSubnets != nil {
		addressPrefixesMap["has_subnets"] = addressPrefixesItem.HasSubnets
	}
	if addressPrefixesItem.Href != nil {
		addressPrefixesMap["href"] = addressPrefixesItem.Href
	}
	if addressPrefixesItem.ID != nil {
		addressPrefixesMap["id"] = addressPrefixesItem.ID
	}
	if addressPrefixesItem.IsDefault != nil {
		addressPrefixesMap["is_default"] = addressPrefixesItem.IsDefault
	}
	if addressPrefixesItem.Name != nil {
		addressPrefixesMap["name"] = addressPrefixesItem.Name
	}
	if addressPrefixesItem.Zone != nil {
		zoneList := []map[string]interface{}{}
		zoneMap := dataSourceAddressPrefixCollectionAddressPrefixesZoneToMap(*addressPrefixesItem.Zone)
		zoneList = append(zoneList, zoneMap)
		addressPrefixesMap["zone"] = zoneList
	}

	return addressPrefixesMap
}

func dataSourceAddressPrefixCollectionAddressPrefixesZoneToMap(zoneItem vpcv1.ZoneReference) (zoneMap map[string]interface{}) {
	zoneMap = map[string]interface{}{}

	if zoneItem.Href != nil {
		zoneMap["href"] = zoneItem.Href
	}
	if zoneItem.Name != nil {
		zoneMap["name"] = zoneItem.Name
	}

	return zoneMap
}

func dataSourceAddressPrefixCollectionFlattenFirst(result vpcv1.AddressPrefixCollectionFirst) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceAddressPrefixCollectionFirstToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceAddressPrefixCollectionFirstToMap(firstItem vpcv1.AddressPrefixCollectionFirst) (firstMap map[string]interface{}) {
	firstMap = map[string]interface{}{}

	if firstItem.Href != nil {
		firstMap["href"] = firstItem.Href
	}

	return firstMap
}

func dataSourceAddressPrefixCollectionFlattenNext(result vpcv1.AddressPrefixCollectionNext) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceAddressPrefixCollectionNextToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceAddressPrefixCollectionNextToMap(nextItem vpcv1.AddressPrefixCollectionNext) (nextMap map[string]interface{}) {
	nextMap = map[string]interface{}{}

	if nextItem.Href != nil {
		nextMap["href"] = nextItem.Href
	}

	return nextMap
}
