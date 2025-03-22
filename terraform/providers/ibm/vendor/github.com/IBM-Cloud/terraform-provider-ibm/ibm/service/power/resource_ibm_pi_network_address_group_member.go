// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceIBMPINetworkAddressGroupMember() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPINetworkAddressGroupMemberCreate,
		ReadContext:   resourceIBMPINetworkAddressGroupMemberRead,
		DeleteContext: resourceIBMPINetworkAddressGroupMemberDelete,
		Importer:      &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_Cidr: {
				Description:  "The member to add in CIDR format.",
				ExactlyOneOf: []string{Arg_Cidr, Arg_NetworkAddressGroupMemberID},
				ForceNew:     true,
				Optional:     true,
				Type:         schema.TypeString,
			},
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				ForceNew:     true,
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_NetworkAddressGroupID: {
				Description:  "Network Address Group ID.",
				ForceNew:     true,
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_NetworkAddressGroupMemberID: {
				Description:  "The network address group member id to remove.",
				ExactlyOneOf: []string{Arg_Cidr, Arg_NetworkAddressGroupMemberID},
				ForceNew:     true,
				Optional:     true,
				Type:         schema.TypeString,
			},
			// Attributes
			Attr_CRN: {
				Computed:    true,
				Description: "The network address group's crn.",
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
func resourceIBMPINetworkAddressGroupMemberCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	nagID := d.Get(Arg_NetworkAddressGroupID).(string)
	nagC := instance.NewIBMPINetworkAddressGroupClient(ctx, sess, cloudInstanceID)
	var body = &models.NetworkAddressGroupAddMember{}
	if v, ok := d.GetOk(Arg_Cidr); ok {
		cidr := v.(string)
		body.Cidr = &cidr
		NetworkAddressGroupMember, err := nagC.AddMember(nagID, body)
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = isWaitForIBMPINetworkAddressGroupMemberAdd(ctx, nagC, nagID, *NetworkAddressGroupMember.ID, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(fmt.Sprintf("%s/%s/%s", cloudInstanceID, nagID, *NetworkAddressGroupMember.ID))
	}
	if v, ok := d.GetOk(Arg_NetworkAddressGroupMemberID); ok {
		memberID := v.(string)
		err := nagC.DeleteMember(nagID, memberID)
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = isWaitForIBMPINetworkAddressGroupMemberRemove(ctx, nagC, nagID, memberID, d.Timeout(schema.TimeoutDelete))
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, nagID))
	}

	return resourceIBMPINetworkAddressGroupMemberRead(ctx, d, meta)
}
func resourceIBMPINetworkAddressGroupMemberRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	nagC := instance.NewIBMPINetworkAddressGroupClient(ctx, sess, parts[0])
	networkAddressGroup, err := nagC.Get(parts[1])
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), NotFound) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if networkAddressGroup.Crn != nil {
		d.Set(Attr_CRN, networkAddressGroup.Crn)
		userTags, err := flex.GetTagsUsingCRN(meta, string(*networkAddressGroup.Crn))
		if err != nil {
			log.Printf("Error on get of pi network address group (%s) user_tags: %s", parts[1], err)
		}
		d.Set(Attr_UserTags, userTags)
	}
	if len(networkAddressGroup.Members) > 0 {
		members := []map[string]interface{}{}
		for _, mbr := range networkAddressGroup.Members {
			member := memberToMap(mbr)
			members = append(members, member)
		}
		d.Set(Attr_Members, members)
	} else {
		d.Set(Attr_Members, nil)
	}
	d.Set(Attr_Name, networkAddressGroup.Name)

	return nil
}
func resourceIBMPINetworkAddressGroupMemberDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if len(parts) > 2 {
		nagC := instance.NewIBMPINetworkAddressGroupClient(ctx, sess, parts[0])
		err = nagC.DeleteMember(parts[1], parts[2])
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = isWaitForIBMPINetworkAddressGroupMemberRemove(ctx, nagC, parts[1], parts[2], d.Timeout(schema.TimeoutDelete))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId("")

	return nil
}
func isWaitForIBMPINetworkAddressGroupMemberAdd(ctx context.Context, client *instance.IBMPINetworkAddressGroupClient, id, memberID string, timeout time.Duration) (interface{}, error) {

	stateConf := &retry.StateChangeConf{
		Pending:    []string{State_Pending},
		Target:     []string{State_Available},
		Refresh:    isIBMPINetworkAddressGroupMemberAddRefreshFunc(client, id, memberID),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	return stateConf.WaitForStateContext(ctx)
}
func isIBMPINetworkAddressGroupMemberAddRefreshFunc(client *instance.IBMPINetworkAddressGroupClient, id, memberID string) retry.StateRefreshFunc {

	return func() (interface{}, string, error) {
		networkAddressGroup, err := client.Get(id)
		if err != nil {
			return nil, "", err
		}

		if len(networkAddressGroup.Members) > 0 {
			var mbrIDs []string
			for _, mbr := range networkAddressGroup.Members {
				mbrIDs = append(mbrIDs, *mbr.ID)
			}
			if slices.Contains(mbrIDs, memberID) {
				return networkAddressGroup, State_Available, nil
			}
		}
		return networkAddressGroup, State_Pending, nil
	}
}
func isWaitForIBMPINetworkAddressGroupMemberRemove(ctx context.Context, client *instance.IBMPINetworkAddressGroupClient, id, memberID string, timeout time.Duration) (interface{}, error) {

	stateConf := &retry.StateChangeConf{
		Pending:    []string{State_Pending},
		Target:     []string{State_Removed},
		Refresh:    isIBMPINetworkAddressGroupMemberRemoveRefreshFunc(client, id, memberID),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	return stateConf.WaitForStateContext(ctx)
}
func isIBMPINetworkAddressGroupMemberRemoveRefreshFunc(client *instance.IBMPINetworkAddressGroupClient, id, memberID string) retry.StateRefreshFunc {

	return func() (interface{}, string, error) {
		networkAddressGroup, err := client.Get(id)
		if err != nil {
			return nil, "", err
		}
		var mbrIDs []string
		for _, mbr := range networkAddressGroup.Members {
			mbrIDs = append(mbrIDs, *mbr.ID)
		}
		if !slices.Contains(mbrIDs, memberID) {
			return networkAddressGroup, State_Removed, nil
		}
		return networkAddressGroup, State_Pending, nil
	}
}
