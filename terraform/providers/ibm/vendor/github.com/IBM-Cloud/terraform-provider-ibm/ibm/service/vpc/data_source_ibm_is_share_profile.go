// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIbmIsShareProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmIsShareProfileRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The file share profile name.",
			},
			"capacity": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The permitted capacity range (in gigabytes) for a share with this profile.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The default capacity.",
						},
						"max": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The max capacity.",
						},
						"min": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The min capacity.",
						},
						"step": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The increment step value for this profile field.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field.",
						},
						"value": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The value for this profile field",
						},
						"values": {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Description: "The permitted values for this profile field.",
						},
					},
				},
			},
			"iops": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The permitted IOPS range for a share with this profile.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The default iops.",
						},
						"max": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The max iops.",
						},
						"min": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The min iops.",
						},
						"step": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The increment step value for this profile field.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field.",
						},
						"value": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The value for this profile field",
						},
						"values": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The permitted values for this profile field.",
						},
					},
				},
			},
			"family": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The product family this share profile belongs to.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this share profile.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
		},
	}
}

func dataSourceIbmIsShareProfileRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getShareProfileOptions := &vpcv1.GetShareProfileOptions{}

	getShareProfileOptions.SetName(d.Get("name").(string))

	shareProfile, response, err := vpcClient.GetShareProfileWithContext(context, getShareProfileOptions)
	if err != nil {
		log.Printf("[DEBUG] GetShareProfileWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	if shareProfile.Capacity != nil {
		capacityList := []map[string]interface{}{}
		capacity := shareProfile.Capacity.(*vpcv1.ShareProfileCapacity)
		capacityMap := dataSourceShareProfileCapacityToMap(*capacity)
		capacityList = append(capacityList, capacityMap)
		if err = d.Set("capacity", capacityList); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting capacity: %s", err))
		}
	}
	if shareProfile.Iops != nil {
		iopsList := []map[string]interface{}{}
		iops := shareProfile.Iops.(*vpcv1.ShareProfileIops)
		iopsMap := dataSourceShareProfileIopsToMap(*iops)
		iopsList = append(iopsList, iopsMap)
		if err = d.Set("iops", iopsList); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting iops: %s", err))
		}
	}
	d.SetId(*shareProfile.Name)
	if err = d.Set("family", shareProfile.Family); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting family: %s", err))
	}
	if err = d.Set("href", shareProfile.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("resource_type", shareProfile.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}

	return nil
}
func dataSourceShareProfileIopsToMap(iops vpcv1.ShareProfileIops) (iopsMap map[string]interface{}) {
	iopsMap = map[string]interface{}{}
	if iops.Default != nil {
		iopsMap["default"] = int(*iops.Default)
	}
	iopsMap["max"] = iops.Max
	iopsMap["min"] = iops.Min
	iopsMap["step"] = iops.Step
	iopsMap["type"] = iops.Type
	if iops.Value != nil {
		iopsMap["value"] = int(*iops.Value)
	}
	values := []int{}
	if len(iops.Values) > 0 {
		for _, value := range iops.Values {
			values = append(values, int(value))
		}
		iopsMap["values"] = values
	}
	return iopsMap
}
func dataSourceShareProfileCapacityToMap(capacity vpcv1.ShareProfileCapacity) (capacityMap map[string]interface{}) {
	capacityMap = map[string]interface{}{}
	// if capacity.Default != nil {
	// 	capacityMap["default"] = int(*capacity.Default)
	// }
	capacityMap["max"] = capacity.Max
	capacityMap["min"] = capacity.Min
	capacityMap["step"] = capacity.Step
	capacityMap["type"] = capacity.Type
	if capacity.Value != nil {
		capacityMap["value"] = int(*capacity.Value)
	}
	values := []int{}
	if len(capacity.Values) > 0 {
		for _, value := range capacity.Values {
			values = append(values, int(value))
		}
		capacityMap["values"] = values
	}
	return capacityMap
}

// dataSourceIbmIsShareProfileID returns a reasonable ID for the list.
func dataSourceIbmIsShareProfileID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
