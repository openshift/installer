// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceIBMPINetworkAddressGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPINetworkAddressGroupCreate,
		ReadContext:   resourceIBMPINetworkAddressGroupRead,
		UpdateContext: resourceIBMPINetworkAddressGroupUpdate,
		DeleteContext: resourceIBMPINetworkAddressGroupDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		CustomizeDiff: customdiff.Sequence(
			func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
				return flex.ResourcePowerUserTagsCustomizeDiff(diff)
			},
		),
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				ForceNew:     true,
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_Name: {
				Description:  "The name of the Network Address Group.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_UserTags: {
				Computed:    true,
				Description: "The user tags associated with this resource.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Set:         schema.HashString,
				Type:        schema.TypeSet,
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
				Optional: true,
				Type:     schema.TypeList,
			},
			Attr_NetworkAddressGroupID: {
				Computed:    true,
				Description: "The unique identifier of the network address group.",
				Type:        schema.TypeString,
			},
		},
	}
}

func resourceIBMPINetworkAddressGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	name := d.Get(Arg_Name).(string)
	nagC := instance.NewIBMPINetworkAddressGroupClient(ctx, sess, cloudInstanceID)
	var body = &models.NetworkAddressGroupCreate{
		Name: &name,
	}

	if v, ok := d.GetOk(Arg_UserTags); ok {
		body.UserTags = flex.ExpandStringList(v.([]interface{}))
	}

	networkAddressGroup, err := nagC.Create(body)
	if err != nil {
		return diag.FromErr(err)
	}
	if _, ok := d.GetOk(Arg_UserTags); ok {
		if networkAddressGroup.Crn != nil {
			oldList, newList := d.GetChange(Arg_UserTags)
			err := flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, string(*networkAddressGroup.Crn), "", UserTagType)
			if err != nil {
				log.Printf("Error on update of pi network address group (%s) pi_user_tags during creation: %s", *networkAddressGroup.ID, err)
			}
		}
	}
	d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, *networkAddressGroup.ID))

	return resourceIBMPINetworkAddressGroupRead(ctx, d, meta)
}

func resourceIBMPINetworkAddressGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID, nagID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	nagC := instance.NewIBMPINetworkAddressGroupClient(ctx, sess, cloudInstanceID)
	networkAddressGroup, err := nagC.Get(nagID)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), NotFound) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	d.Set(Arg_Name, networkAddressGroup.Name)
	if networkAddressGroup.Crn != nil {
		d.Set(Attr_CRN, networkAddressGroup.Crn)
		userTags, err := flex.GetTagsUsingCRN(meta, string(*networkAddressGroup.Crn))
		if err != nil {
			log.Printf("Error on get of network address group (%s) pi_user_tags: %s", nagID, err)
		}
		d.Set(Arg_UserTags, userTags)
	}

	d.Set(Attr_NetworkAddressGroupID, networkAddressGroup.ID)
	members := []map[string]interface{}{}
	if len(networkAddressGroup.Members) > 0 {
		for _, mbr := range networkAddressGroup.Members {
			member := memberToMap(mbr)
			members = append(members, member)
		}
		d.Set(Attr_Members, members)
	} else {
		d.Set(Attr_Members, nil)
	}

	return nil
}

func resourceIBMPINetworkAddressGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	hasChange := false
	body := &models.NetworkAddressGroupUpdate{}
	if d.HasChange(Arg_UserTags) {
		if crn, ok := d.GetOk(Attr_CRN); ok {
			oldList, newList := d.GetChange(Arg_UserTags)
			err := flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, crn.(string), "", UserTagType)
			if err != nil {
				log.Printf("Error on update of pi network address group (%s) pi_user_tags: %s", parts[1], err)
			}
		}
	}
	if d.HasChange(Arg_Name) {
		body.Name = d.Get(Arg_Name).(string)
		hasChange = true
	}
	nagC := instance.NewIBMPINetworkAddressGroupClient(ctx, sess, parts[0])
	if hasChange {
		_, err := nagC.Update(parts[1], body)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceIBMPINetworkAddressGroupRead(ctx, d, meta)
}

func resourceIBMPINetworkAddressGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	nagC := instance.NewIBMPINetworkAddressGroupClient(ctx, sess, parts[0])
	err = nagC.Delete(parts[1])
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = isWaitForIBMPINetworkAddressGroupDeleted(ctx, nagC, parts[1], d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")

	return nil
}

func isWaitForIBMPINetworkAddressGroupDeleted(ctx context.Context, client *instance.IBMPINetworkAddressGroupClient, nagID string, timeout time.Duration) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{State_Deleting},
		Target:     []string{State_NotFound},
		Refresh:    isIBMPINetworkAddressGroupDeleteRefreshFunc(client, nagID),
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
		Timeout:    timeout,
	}

	return stateConf.WaitForStateContext(ctx)
}
func isIBMPINetworkAddressGroupDeleteRefreshFunc(client *instance.IBMPINetworkAddressGroupClient, nagID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		nag, err := client.Get(nagID)
		if err != nil {
			return nag, State_NotFound, nil
		}
		return nag, State_Deleting, nil
	}
}
