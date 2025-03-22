// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
)

func ResourceIBMPIHostGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPIHostGroupCreate,
		ReadContext:   resourceIBMPIHostGroupRead,
		DeleteContext: resourceIBMPIHostGroupDelete,
		UpdateContext: resourceIBMPIHostGroupUpdate,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
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
			Arg_Hosts: {
				Description: "List of hosts to add to the group.",
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_DisplayName: {
							Description:  "Name of the host chosen by the user.",
							Required:     true,
							Type:         schema.TypeString,
							ValidateFunc: validation.NoZeroValues,
						},
						Attr_SysType: {
							Description:  "System type.",
							Required:     true,
							Type:         schema.TypeString,
							ValidateFunc: validation.NoZeroValues,
						},
						Attr_UserTags: {
							Description: "List of user tags attached to the resource.",
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Set:         schema.HashString,
							Type:        schema.TypeSet,
						},
					},
				},
				Required: true,
				Type:     schema.TypeSet,
			},
			Arg_Name: {
				Description:  "Name of the host group to create.",
				ForceNew:     true,
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_Remove: {
				Description: "A workspace ID to stop sharing the host group with.",
				Optional:    true,
				Type:        schema.TypeString,
			},
			Arg_Secondaries: {
				Description: "List of workspaces to share the host group with.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_Name: {
							Description: "Name of the host group to create in the secondary workspace.",
							Optional:    true,
							Type:        schema.TypeString,
						},
						Attr_Workspace: {
							Description: "ID of the workspace to share the host group with.",
							Required:    true,
							Type:        schema.TypeString,
						},
					},
				},
				Optional: true,
				Type:     schema.TypeSet,
			},
			// Attributes
			Attr_CreationDate: {
				Computed:    true,
				Description: "Date/Time of host group creation.",
				ForceNew:    true,
				Type:        schema.TypeString,
			},
			Attr_HostGroupID: {
				Computed:    true,
				Description: "Host group ID.",
				ForceNew:    true,
				Type:        schema.TypeString,
			},
			Attr_Hosts: {
				Computed:    true,
				Description: "List of hosts.",
				ForceNew:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Type: schema.TypeList,
			},
			Attr_Name: {
				Computed:    true,
				Description: "Name of the host group.",
				ForceNew:    true,
				Type:        schema.TypeString,
			},
			Attr_Primary: {
				Computed:    true,
				Description: "ID of the workspace owning the host group.",
				ForceNew:    true,
				Type:        schema.TypeString,
			},
			Attr_Secondaries: {
				Computed:    true,
				Description: "IDs of workspaces the host group has been shared with.",
				ForceNew:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Type: schema.TypeList,
			},
		},
	}
}

func resourceIBMPIHostGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	name := d.Get(Arg_Name).(string)
	client := instance.NewIBMPIHostGroupsClient(ctx, sess, cloudInstanceID)
	body := models.HostGroupCreate{}
	var hosts []*models.AddHost
	for _, v := range d.Get(Arg_Hosts).(*schema.Set).List() {
		hostData := v.(map[string]interface{})
		host := hostMapToAddHost(hostData)
		hosts = append(hosts, host)
	}
	body.Hosts = hosts
	body.Name = &name

	if _, ok := d.GetOk(Arg_Secondaries); ok {
		var secondaries []*models.Secondary
		for _, v := range d.Get(Arg_Secondaries).(*schema.Set).List() {
			secData := v.(map[string]interface{})
			secondary := secondaryMapToSecondary(secData)
			secondaries = append(secondaries, secondary)
		}
		body.Secondaries = secondaries
	}

	hg, err := client.CreateHostGroup(&body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, hg.ID))

	return resourceIBMPIHostGroupRead(ctx, d, meta)
}

func resourceIBMPIHostGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID, hostGroupID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	client := instance.NewIBMPIHostGroupsClient(ctx, sess, cloudInstanceID)
	hostGroup, err := client.GetHostGroup(hostGroupID)
	if err != nil {
		if strings.Contains(err.Error(), NotFound) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	d.Set(Attr_CreationDate, hostGroup.CreationDate.String())
	d.Set(Attr_HostGroupID, hostGroup.ID)
	d.Set(Attr_Hosts, hostGroup.Hosts)
	d.Set(Attr_Primary, hostGroup.Primary)
	d.Set(Attr_Secondaries, hostGroup.Secondaries)

	return nil
}

func resourceIBMPIHostGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID, hostGroupID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	client := instance.NewIBMPIHostGroupsClient(ctx, sess, cloudInstanceID)
	hostGroupUpdateBody := models.HostGroupShareOp{}
	hasChange := false
	if d.HasChange(Arg_Remove) {
		hostGroupUpdateBody.Remove = d.Get(Arg_Remove).(string)
		hasChange = true
	}

	if d.HasChange(Arg_Secondaries) {
		oldSecondaries, newSecondaries := d.GetChange(Arg_Secondaries)
		if len(oldSecondaries.([]interface{})) == len(newSecondaries.([]interface{})) {
			return diag.FromErr(fmt.Errorf("change in place not supported for: %v", Arg_Secondaries))
		}
		var add []*models.Secondary
		for _, v := range d.Get(Arg_Secondaries).([]interface{}) {
			secData := v.(map[string]interface{})
			addItem := secondaryMapToSecondary(secData)
			add = append(add, addItem)
		}
		hostGroupUpdateBody.Add = add
		hasChange = true
	}

	if hasChange {
		_, err := client.UpdateHostGroup(&hostGroupUpdateBody, hostGroupID)
		if err != nil {
			if strings.Contains(err.Error(), NotFound) {
				d.SetId("")
				return nil
			}
			return diag.FromErr(err)
		}
	}

	return resourceIBMPIHostGroupRead(ctx, d, meta)
}

func resourceIBMPIHostGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID, hostGroupID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	client := instance.NewIBMPIHostGroupsClient(ctx, sess, cloudInstanceID)
	hostGroup, err := client.GetHostGroup(hostGroupID)
	if err != nil {
		if strings.Contains(err.Error(), NotFound) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	for _, v := range hostGroup.Hosts {
		ref, err := json.Marshal(v)
		if err != nil {
			fmt.Printf("error while json Marshal: %v", err)
		}
		hostRef := string(ref)
		hostID, err := getLastPart(hostRef)
		if err != nil {
			return diag.FromErr(err)
		}
		err = client.DeleteHost(hostID)
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = isWaitForHostDeleted(ctx, client, hostID, d.Timeout(schema.TimeoutDelete))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	_, err = isWaitForHostGroupDeleted(ctx, client, hostGroupID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func isWaitForHostGroupDeleted(ctx context.Context, client *instance.IBMPIHostGroupsClient, id string, timeout time.Duration) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{State_Deleting},
		Target:     []string{NotFound},
		Refresh:    isHostGroupDeleteRefresh(client, id),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
		Timeout:    timeout,
	}
	return stateConf.WaitForStateContext(ctx)
}

func isHostGroupDeleteRefresh(client *instance.IBMPIHostGroupsClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		hg, err := client.GetHostGroup(id)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), NotFound) {

				return hg, NotFound, nil
			}
		}
		return hg, State_Deleting, nil
	}
}

func isWaitForHostDeleted(ctx context.Context, client *instance.IBMPIHostGroupsClient, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for host (%s) to be deleted.", id)
	stateConf := &retry.StateChangeConf{
		Pending:    []string{State_Deleting},
		Target:     []string{NotFound},
		Refresh:    isHostDeleteRefreshFunc(client, id),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
		Timeout:    timeout,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isHostDeleteRefreshFunc(client *instance.IBMPIHostGroupsClient, id string) retry.StateRefreshFunc {
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

func hostMapToAddHost(modelMap map[string]interface{}) *models.AddHost {
	host := &models.AddHost{}
	host.DisplayName = core.StringPtr(modelMap[Attr_DisplayName].(string))
	host.SysType = core.StringPtr(modelMap[Attr_SysType].(string))
	host.UserTags = flex.FlattenSet(modelMap[Attr_UserTags].(*schema.Set))
	return host
}

func secondaryMapToSecondary(modelMap map[string]interface{}) *models.Secondary {
	secondary := &models.Secondary{}
	if modelMap[Attr_Name].(string) != "" {
		secondary.Name = modelMap[Attr_Name].(string)
	}
	secondary.Workspace = core.StringPtr(modelMap[Attr_Workspace].(string))
	return secondary
}

// This function trims unwanted characters from each part of the id
func getLastPart(id string) (string, error) {
	parts := strings.Split(id, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid input format")
	}
	lastPart := parts[len(parts)-1]
	cleanedLastPart := strings.Trim(lastPart, `"`)
	return cleanedLastPart, nil
}
