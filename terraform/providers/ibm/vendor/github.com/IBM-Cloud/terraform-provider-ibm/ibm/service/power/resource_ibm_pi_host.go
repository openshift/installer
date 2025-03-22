// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
)

func ResourceIBMPIHost() *schema.Resource {
	return &schema.Resource{
		CustomizeDiff: customdiff.Sequence(
			func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
				return customizeUserTagsPIHostDiff(diff)
			},
		),
		CreateContext: resourceIBMPIHostCreate,
		ReadContext:   resourceIBMPIHostRead,
		UpdateContext: resourceIBMPIHostUpdate,
		DeleteContext: resourceIBMPIHostDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				ForceNew:     true,
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_HostGroupID: {
				Description:  "ID of the host group to which the host should be added.",
				ForceNew:     true,
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_Host: {
				Description: "Host to add to a host group.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_DisplayName: {
							Description: "Name of the host chosen by the user.",
							Required:    true,
							Type:        schema.TypeString,
						},
						Attr_SysType: {
							Description: "System type.",
							ForceNew:    true,
							Required:    true,
							Type:        schema.TypeString,
						},
						Attr_UserTags: {
							Computed:    true,
							Description: "List of user tags attached to the resource.",
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Set:         schema.HashString,
							Type:        schema.TypeSet,
						},
					},
				},
				MaxItems: 1,
				Required: true,
				Type:     schema.TypeList,
			},
			// Attributes
			Attr_Capacity: {
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_AvailableCores: {
							Computed:    true,
							Description: "Number of cores currently available.",
							Type:        schema.TypeFloat,
						},
						Attr_AvailableMemory: {
							Computed:    true,
							Description: "Amount of memory currently available (in GB).",
							Type:        schema.TypeFloat,
						},
						Attr_ReservedCore: {
							Computed:    true,
							Description: "Number of cores reserved for system use.",
							Type:        schema.TypeFloat,
						},
						Attr_ReservedMemory: {
							Computed:    true,
							Description: "Amount of memory reserved for system use (in GB).",
							Type:        schema.TypeFloat,
						},
						Attr_TotalCore: {
							Computed:    true,
							Description: "Total number of cores of the host.",
							Type:        schema.TypeFloat,
						},
						Attr_TotalMemory: {
							Computed:    true,
							Description: "Total amount of memory of the host (in GB).",
							Type:        schema.TypeFloat,
						},
						Attr_UsedCore: {
							Computed:    true,
							Description: "Number of cores in use on the host.",
							Type:        schema.TypeFloat,
						},
						Attr_UsedMemory: {
							Computed:    true,
							Description: "Amount of memory used on the host (in GB).",
							Type:        schema.TypeFloat,
						},
					},
				},
				Type: schema.TypeList,
			},
			Attr_CRN: {
				Computed:    true,
				Description: "The CRN of this resource.",
				Type:        schema.TypeString,
			},
			Attr_DisplayName: {
				Computed:    true,
				Description: "Name of the host (chosen by the user).",
				Type:        schema.TypeString,
			},
			Attr_HostID: {
				Computed:    true,
				Description: "ID of the host.",
				Type:        schema.TypeString,
			},
			Attr_HostGroup: {
				Computed:    true,
				Description: "Link to host group resource.",
				Type:        schema.TypeMap,
			},
			Attr_State: {
				Computed:    true,
				Description: "State of the host (up/down).",
				Type:        schema.TypeString,
			},
			Attr_Status: {
				Computed:    true,
				Description: "Status of the host (enabled/disabled).",
				Type:        schema.TypeString,
			},
			Attr_SysType: {
				Computed:    true,
				Description: "System type.",
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

func resourceIBMPIHostCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	client := instance.NewIBMPIHostGroupsClient(ctx, sess, cloudInstanceID)
	hosts := d.Get(Arg_Host).([]interface{})
	hostGroupID := d.Get(Arg_HostGroupID).(string)
	body := models.HostCreate{}
	hostBody := make([]*models.AddHost, 0, len(hosts))
	for _, v := range hosts {
		host := v.(map[string]interface{})
		hs := models.AddHost{
			DisplayName: core.StringPtr(host[Attr_DisplayName].(string)),
			SysType:     core.StringPtr(host[Attr_SysType].(string)),
			UserTags:    flex.FlattenSet(host[Attr_UserTags].(*schema.Set)),
		}
		hostBody = append(hostBody, &hs)
	}

	body.Hosts = hostBody
	body.HostGroupID = &hostGroupID
	hostResponse, err := client.CreateHost(&body)
	if err != nil {
		return diag.FromErr(err)
	}

	hostID := hostResponse[0].ID
	d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, hostID))
	_, err = isWaitForIBMPIHostAvailable(ctx, client, hostID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	host := hosts[0].(map[string]interface{})
	tags := flex.FlattenSet(host[Attr_UserTags].(*schema.Set))
	if hostResponse[0].Crn != "" && len(tags) > 0 {
		oldList, newList := d.GetChange(Arg_Host + ".0." + Attr_UserTags)
		err := flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, string(hostResponse[0].Crn), "", UserTagType)
		if err != nil {
			log.Printf("Error on update of pi host (%s) user_tags during creation: %s", hostResponse[0].ID, err)
		}
	}

	return resourceIBMPIHostRead(ctx, d, meta)
}

func resourceIBMPIHostRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID, hostID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	client := instance.NewIBMPIHostGroupsClient(ctx, sess, cloudInstanceID)
	host, err := client.GetHost(hostID)
	if err != nil {
		if strings.Contains(err.Error(), NotFound) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	d.Set(Arg_CloudInstanceID, cloudInstanceID)
	hostGroupID, err := getLastPart(host.HostGroup.Href)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set(Arg_HostGroupID, hostGroupID)
	d.Set(Attr_HostID, host.ID)

	if host.Capacity != nil {
		d.Set(Attr_Capacity, hostCapacityToMap(host.Capacity))
	}
	if host.Crn != "" {
		d.Set(Attr_CRN, host.Crn)
		tags, err := flex.GetGlobalTagsUsingCRN(meta, string(host.Crn), "", UserTagType)
		if err != nil {
			log.Printf("Error on get of pi host (%s) user_tags: %s", host.ID, err)
		}
		d.Set(Attr_UserTags, tags)
	}
	if host.DisplayName != "" {
		d.Set(Attr_DisplayName, host.DisplayName)
	}
	if host.HostGroup != nil {
		d.Set(Attr_HostGroup, hostGroupToMap(host.HostGroup))
	}
	if host.State != "" {
		d.Set(Attr_State, host.State)
	}
	if host.Status != "" {
		d.Set(Attr_Status, host.Status)
	}
	if host.SysType != "" {
		d.Set(Attr_SysType, host.SysType)
	}
	d.Set(Arg_Host, flattenHostArgumentToList(d, meta))

	return nil
}

func resourceIBMPIHostUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID, hostID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	client := instance.NewIBMPIHostGroupsClient(ctx, sess, cloudInstanceID)
	if d.HasChange(Arg_Host) {
		oldHost, newHost := d.GetChange(Arg_Host + ".0")

		displayNameOld := oldHost.(map[string]interface{})[Attr_DisplayName].(string)
		displayNameNew := newHost.(map[string]interface{})[Attr_DisplayName].(string)

		if displayNameNew != displayNameOld {
			hostBody := models.HostPut{
				DisplayName: &displayNameNew,
			}
			_, err := client.UpdateHost(&hostBody, hostID)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if crn, ok := d.GetOk(Attr_CRN); ok {
			userTagsOld := oldHost.(map[string]interface{})[Attr_UserTags].(*schema.Set)
			userTagsNew := newHost.(map[string]interface{})[Attr_UserTags].(*schema.Set)
			if !userTagsNew.Equal(userTagsOld) {
				err = flex.UpdateGlobalTagsUsingCRN(userTagsOld, userTagsNew, meta, crn.(string), "", UserTagType)
				if err != nil {
					log.Printf("Error on update of pi host (%s) pi_host user_tags: %s", d.Get(Attr_HostID), err)
				}
			}
		}

	}

	return resourceIBMPIHostRead(ctx, d, meta)
}

func resourceIBMPIHostDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID, hostID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	client := instance.NewIBMPIHostGroupsClient(ctx, sess, cloudInstanceID)
	err = client.DeleteHost(hostID)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = isWaitForPIHostDeleted(ctx, client, hostID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")

	return nil
}
func isWaitForPIHostDeleted(ctx context.Context, client *instance.IBMPIHostGroupsClient, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for host (%s) to be deleted.", id)
	stateConf := &retry.StateChangeConf{
		Pending:    []string{State_Deleting},
		Target:     []string{NotFound},
		Refresh:    isPIHostDeleteRefreshFunc(client, id),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
		Timeout:    timeout,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isPIHostDeleteRefreshFunc(client *instance.IBMPIHostGroupsClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		host, err := client.GetHost(id)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), NotFound) {

				return host, NotFound, nil
			}
		}
		return host, State_Deleting, nil
	}
}
func hostCapacityToMap(capicity *models.HostCapacity) (hostCapacity []map[string]interface{}) {
	hostCapacityMap := make(map[string]interface{})
	if capicity.Cores.Available != nil {
		hostCapacityMap[Attr_AvailableCores] = capicity.Cores.Available
	}
	if capicity.Memory.Available != nil {
		hostCapacityMap[Attr_AvailableMemory] = capicity.Memory.Available
	}
	if capicity.Cores.Reserved != nil {
		hostCapacityMap[Attr_ReservedCore] = capicity.Cores.Reserved
	}
	if capicity.Memory.Reserved != nil {
		hostCapacityMap[Attr_ReservedMemory] = capicity.Memory.Reserved
	}
	if capicity.Cores.Total != nil {
		hostCapacityMap[Attr_TotalCore] = capicity.Cores.Total
	}
	if capicity.Memory.Total != nil {
		hostCapacityMap[Attr_TotalMemory] = capicity.Memory.Total
	}
	if capicity.Cores.Used != nil {
		hostCapacityMap[Attr_UsedCore] = capicity.Cores.Used
	}
	if capicity.Memory.Used != nil {
		hostCapacityMap[Attr_UsedMemory] = capicity.Memory.Used
	}
	hostCapacity = append(hostCapacity, hostCapacityMap)
	return hostCapacity
}
func hostGroupToMap(hostgroup *models.HostGroupSummary) map[string]interface{} {
	hostGroupMap := make(map[string]interface{})
	if hostgroup.Access != "" {
		hostGroupMap[Attr_Access] = hostgroup.Access
	}
	if hostgroup.Href != "" {
		hostGroupMap[Attr_Href] = hostgroup.Href
	}
	if hostgroup.Name != "" {
		hostGroupMap[Attr_Name] = hostgroup.Name
	}
	return hostGroupMap
}

func isWaitForIBMPIHostAvailable(ctx context.Context, client *instance.IBMPIHostGroupsClient, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for  host (%s) to be available.", id)
	stateConf := &retry.StateChangeConf{
		Pending:    []string{State_Down},
		Target:     []string{State_Up},
		Refresh:    isIBMPIHostRefreshFunc(client, id),
		Timeout:    timeout,
		Delay:      20 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isIBMPIHostRefreshFunc(client *instance.IBMPIHostGroupsClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		host, err := client.GetHost(id)
		if err != nil {
			return nil, "", err
		}
		if host.State == State_Up {
			return host, State_Up, nil
		}
		return host, State_Down, nil
	}
}

func flattenHostArgumentToList(d *schema.ResourceData, meta interface{}) []map[string]interface{} {
	hostListType := make([]map[string]interface{}, 0)
	h := map[string]interface{}{}
	if v, ok := d.GetOk(Attr_DisplayName); ok {
		displayName := v.(string)
		h[Attr_DisplayName] = displayName
	}
	if v, ok := d.GetOk(Attr_SysType); ok {
		sysType := v.(string)
		h[Attr_SysType] = sysType
	}
	if v, ok := d.GetOk(Attr_UserTags); ok {
		tags := v
		h[Attr_UserTags] = tags
	}
	hostListType = append(hostListType, h)
	return hostListType
}

func customizeUserTagsPIHostDiff(diff *schema.ResourceDiff) error {
	if diff.Id() != "" && diff.HasChange(Arg_Host+".0."+Attr_UserTags) {
		o, n := diff.GetChange(Arg_Host + ".0." + Attr_UserTags)
		oldSet := o.(*schema.Set)
		newSet := n.(*schema.Set)
		removeInt := oldSet.Difference(newSet).List()
		addInt := newSet.Difference(oldSet).List()
		if v := os.Getenv("IC_ENV_TAGS"); v != "" {
			s := strings.Split(v, ",")
			if len(removeInt) == len(s) && len(addInt) == 0 {
				fmt.Println("Suppresing the TAG diff ")
				return diff.Clear(Arg_Host + ".0." + Attr_UserTags)
			}
		}
		diff.SetNewComputed(Attr_UserTags)
	}
	return nil
}
