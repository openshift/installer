// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"log"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPINetworkAddressGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPINetworkAddressGroupsRead,

		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_NetworkAddressGroups: {
				Computed:    true,
				Description: "list of Network Address Groups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_CRN: {
							Computed:    true,
							Description: "The Network Address Group's crn.",
							Type:        schema.TypeString,
						},
						Attr_ID: {
							Computed:    true,
							Description: "The id of the Network Address Group.",
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
				},
				Type: schema.TypeList,
			},
		},
	}
}

func dataSourceIBMPINetworkAddressGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	nagC := instance.NewIBMPINetworkAddressGroupClient(ctx, sess, cloudInstanceID)
	networkAddressGroups, err := nagC.GetAll()
	if err != nil {
		return diag.FromErr(err)
	}

	var genID, _ = uuid.GenerateUUID()
	d.SetId(genID)

	nags := []map[string]interface{}{}
	if len(networkAddressGroups.NetworkAddressGroups) > 0 {
		for _, nag := range networkAddressGroups.NetworkAddressGroups {
			modelMap := networkAddressGroupsNetworkAddressGroupToMap(nag, meta)
			nags = append(nags, modelMap)
		}
	}
	d.Set(Attr_NetworkAddressGroups, nags)

	return nil
}

func networkAddressGroupsNetworkAddressGroupToMap(networkAddressGroup *models.NetworkAddressGroup, meta interface{}) map[string]interface{} {
	nag := make(map[string]interface{})
	if networkAddressGroup.Crn != nil {
		nag[Attr_CRN] = networkAddressGroup.Crn
		userTags, err := flex.GetTagsUsingCRN(meta, string(*networkAddressGroup.Crn))
		if err != nil {
			log.Printf("Error on get of pi network address group (%s) user_tags: %s", *networkAddressGroup.ID, err)
		}
		nag[Attr_UserTags] = userTags
	}

	nag[Attr_ID] = networkAddressGroup.ID
	if len(networkAddressGroup.Members) > 0 {
		members := []map[string]interface{}{}
		for _, membersItem := range networkAddressGroup.Members {
			member := memberToMap(membersItem)
			members = append(members, member)
		}
		nag[Attr_Members] = members
	}
	nag[Attr_Name] = networkAddressGroup.Name
	return nag
}

func memberToMap(mbr *models.NetworkAddressGroupMember) map[string]interface{} {
	member := make(map[string]interface{})
	member[Attr_CIDR] = mbr.Cidr
	member[Attr_ID] = mbr.ID
	return member
}
