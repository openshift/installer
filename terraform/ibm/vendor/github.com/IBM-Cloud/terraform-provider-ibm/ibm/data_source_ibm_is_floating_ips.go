// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func dataSourceIBMIsFloatingIps() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsFloatingIpsRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique user-defined name for this floating IP.",
			},
			"floating_ips": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of floating IPs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The globally unique IP address.",
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the floating IP was created.",
						},
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this floating IP.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this floating IP.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this floating IP.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique user-defined name for this floating IP.",
						},
						"resource_group": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The resource group for this floating IP.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this resource group.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this resource group.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this resource group.",
									},
								},
							},
						},
						"status": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the floating IP.",
						},
						"target": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The target of this floating IP.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"deleted": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"more_info": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
									},
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this network interface.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this network interface.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this network interface.",
									},
									"primary_ipv4_address": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The primary IPv4 address.If the address has not yet been selected, the value will be `0.0.0.0`.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this public gateway.",
									},
								},
							},
						},
						"zone": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The zone this floating IP resides in.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this zone.",
									},
									"name": &schema.Schema{
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

func dataSourceIBMIsFloatingIpsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	start := ""
	allFloatingIPs := []vpcv1.FloatingIP{}
	for {
		floatingIPOptions := &vpcv1.ListFloatingIpsOptions{}
		if start != "" {
			floatingIPOptions.Start = &start
		}
		floatingIPs, response, err := sess.ListFloatingIps(floatingIPOptions)
		if err != nil {
			log.Printf("[DEBUG] Error Fetching floating IPs  %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("Error Fetching floating IPs %s\n%s", err, response))
		}
		start = GetNext(floatingIPs.Next)
		allFloatingIPs = append(allFloatingIPs, floatingIPs.FloatingIps...)
		if start == "" {
			break
		}
	}
	var matchFloatingIps []vpcv1.FloatingIP
	var name string
	var suppliedFilter bool

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
		suppliedFilter = true
		for _, data := range allFloatingIPs {
			if *data.Name == name {
				matchFloatingIps = append(matchFloatingIps, data)
			}
		}
	} else {
		matchFloatingIps = allFloatingIPs
	}
	if suppliedFilter {
		if len(matchFloatingIps) == 0 {
			return diag.FromErr(fmt.Errorf("no FloatingIps found with name %s", name))
		}
		d.SetId(name)
	} else {
		d.SetId(dataSourceIBMIsFloatingIpsID(d))
	}

	if matchFloatingIps != nil {
		err = d.Set("floating_ips", dataSourceFloatingIPCollectionFlattenFloatingIps(matchFloatingIps))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting floating_ips %s", err))
		}
	}
	return nil
}

// dataSourceIBMIsFloatingIpsID returns a reasonable ID for the list.
func dataSourceIBMIsFloatingIpsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceFloatingIPCollectionFlattenFirst(result vpcv1.FloatingIPCollectionFirst) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceFloatingIPCollectionFirstToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceFloatingIPCollectionFirstToMap(firstItem vpcv1.FloatingIPCollectionFirst) (firstMap map[string]interface{}) {
	firstMap = map[string]interface{}{}

	if firstItem.Href != nil {
		firstMap["href"] = firstItem.Href
	}

	return firstMap
}

func dataSourceFloatingIPCollectionFlattenFloatingIps(result []vpcv1.FloatingIP) (floatingIps []map[string]interface{}) {
	for _, floatingIpsItem := range result {
		floatingIps = append(floatingIps, dataSourceFloatingIPCollectionFloatingIpsToMap(floatingIpsItem))
	}

	return floatingIps
}

func dataSourceFloatingIPCollectionFloatingIpsToMap(floatingIpsItem vpcv1.FloatingIP) (floatingIpsMap map[string]interface{}) {
	floatingIpsMap = map[string]interface{}{}

	if floatingIpsItem.Address != nil {
		floatingIpsMap["address"] = floatingIpsItem.Address
	}
	if floatingIpsItem.CreatedAt != nil {
		floatingIpsMap["created_at"] = floatingIpsItem.CreatedAt.String()
	}
	if floatingIpsItem.CRN != nil {
		floatingIpsMap["crn"] = floatingIpsItem.CRN
	}
	if floatingIpsItem.Href != nil {
		floatingIpsMap["href"] = floatingIpsItem.Href
	}
	if floatingIpsItem.ID != nil {
		floatingIpsMap["id"] = floatingIpsItem.ID
	}
	if floatingIpsItem.Name != nil {
		floatingIpsMap["name"] = floatingIpsItem.Name
	}
	if floatingIpsItem.ResourceGroup != nil {
		resourceGroupList := []map[string]interface{}{}
		resourceGroupMap := dataSourceFloatingIPCollectionFloatingIpsResourceGroupToMap(*floatingIpsItem.ResourceGroup)
		resourceGroupList = append(resourceGroupList, resourceGroupMap)
		floatingIpsMap["resource_group"] = resourceGroupList
	}
	if floatingIpsItem.Status != nil {
		floatingIpsMap["status"] = floatingIpsItem.Status
	}
	if floatingIpsItem.Target != nil {
		targetList := []map[string]interface{}{}
		targetMap := dataSourceFloatingIPCollectionFloatingIpsTargetToMap(floatingIpsItem.Target)
		targetList = append(targetList, targetMap)
		floatingIpsMap["target"] = targetList
	}
	if floatingIpsItem.Zone != nil {
		zoneList := []map[string]interface{}{}
		zoneMap := dataSourceFloatingIPCollectionFloatingIpsZoneToMap(*floatingIpsItem.Zone)
		zoneList = append(zoneList, zoneMap)
		floatingIpsMap["zone"] = zoneList
	}

	return floatingIpsMap
}

func dataSourceFloatingIPCollectionFloatingIpsResourceGroupToMap(resourceGroupItem vpcv1.ResourceGroupReference) (resourceGroupMap map[string]interface{}) {
	resourceGroupMap = map[string]interface{}{}

	if resourceGroupItem.Href != nil {
		resourceGroupMap["href"] = resourceGroupItem.Href
	}
	if resourceGroupItem.ID != nil {
		resourceGroupMap["id"] = resourceGroupItem.ID
	}
	if resourceGroupItem.Name != nil {
		resourceGroupMap["name"] = resourceGroupItem.Name
	}

	return resourceGroupMap
}

func dataSourceFloatingIPCollectionFloatingIpsTargetToMap(targetItemIntf vpcv1.FloatingIPTargetIntf) (targetMap map[string]interface{}) {
	targetMap = map[string]interface{}{}
	targetItem := targetItemIntf.(*vpcv1.FloatingIPTarget)

	if targetItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceFloatingIPCollectionTargetDeletedToMap(*targetItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		targetMap["deleted"] = deletedList
	}
	if targetItem.Href != nil {
		targetMap["href"] = targetItem.Href
	}
	if targetItem.ID != nil {
		targetMap["id"] = targetItem.ID
	}
	if targetItem.Name != nil {
		targetMap["name"] = targetItem.Name
	}
	if targetItem.PrimaryIpv4Address != nil {
		targetMap["primary_ipv4_address"] = targetItem.PrimaryIpv4Address
	}
	if targetItem.ResourceType != nil {
		targetMap["resource_type"] = targetItem.ResourceType
	}
	if targetItem.CRN != nil {
		targetMap["crn"] = targetItem.CRN
	}

	return targetMap
}

func dataSourceFloatingIPCollectionTargetDeletedToMap(deletedItem vpcv1.NetworkInterfaceReferenceDeleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceFloatingIPCollectionFloatingIpsZoneToMap(zoneItem vpcv1.ZoneReference) (zoneMap map[string]interface{}) {
	zoneMap = map[string]interface{}{}

	if zoneItem.Href != nil {
		zoneMap["href"] = zoneItem.Href
	}
	if zoneItem.Name != nil {
		zoneMap["name"] = zoneItem.Name
	}

	return zoneMap
}

func dataSourceFloatingIPCollectionFlattenNext(result vpcv1.FloatingIPCollectionNext) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceFloatingIPCollectionNextToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceFloatingIPCollectionNextToMap(nextItem vpcv1.FloatingIPCollectionNext) (nextMap map[string]interface{}) {
	nextMap = map[string]interface{}{}

	if nextItem.Href != nil {
		nextMap["href"] = nextItem.Href
	}

	return nextMap
}
