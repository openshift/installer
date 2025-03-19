// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
)

func DataSourceIBMPINetworkAddressGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPINetworkAddressGroupRead,

		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_NetworkAddressGroupID: {
				Description: "Network Address Group ID.",
				Required:    true,
				Type:        schema.TypeString,
			},
			// Attributes
			Attr_CRN: {
				Computed:    true,
				Description: "The Network Address Group's crn.",
				Type:        schema.TypeString,
			},
			Attr_Members: {
				Computed:    true,
				Description: "The list of IP addresses in CIDR notation (for example 192.168.66.2/32) in the Network Address Group.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_CIDR: {
							Computed:    true,
							Description: "The IP addresses in CIDR notation for example 192.168.1.5/32.",
							Type:        schema.TypeString,
						},
						Attr_ID: {
							Computed:    true,
							Description: "The id of the Network Address Group member IP addresses.",
							Type:        schema.TypeString,
						},
					},
				},
				Type: schema.TypeList,
			},
			Attr_Name: {
				Computed:    true,
				Description: "The name of the Network Address Group.",
				Type:        schema.TypeString,
			},
			Attr_UserTags: {
				Computed:    true,
				Description: "List of user tags attached to the resource.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Type:        schema.TypeSet,
			},
		},
	}
}

func dataSourceIBMPINetworkAddressGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	nagID := d.Get(Arg_NetworkAddressGroupID).(string)
	nagC := instance.NewIBMPINetworkAddressGroupClient(ctx, sess, cloudInstanceID)
	networkAddressGroup, err := nagC.Get(nagID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*networkAddressGroup.ID)
	if networkAddressGroup.Crn != nil {
		d.Set(Attr_CRN, networkAddressGroup.Crn)
		userTags, err := flex.GetTagsUsingCRN(meta, string(*networkAddressGroup.Crn))
		if err != nil {
			log.Printf("Error on get of pi network address group (%s) user_tags: %s", nagID, err)
		}
		d.Set(Attr_UserTags, userTags)
	}

	members := []map[string]interface{}{}
	if len(networkAddressGroup.Members) > 0 {
		for _, mbr := range networkAddressGroup.Members {
			member := memberToMap(mbr)
			members = append(members, member)
		}
		d.Set(Attr_Members, members)
	}
	d.Set(Attr_Name, networkAddressGroup.Name)

	return nil
}
