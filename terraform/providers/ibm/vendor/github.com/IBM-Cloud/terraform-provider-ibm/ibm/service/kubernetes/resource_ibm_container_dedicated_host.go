// Copyright IBM Corp. 2017, 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package kubernetes

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
)

const (
	created       = "created"
	creating      = "creating"
	createFailed  = "create_failed"
	createPending = "create_pending"
	deleted       = "deleted"
	deleting      = "deleting"
	deleteFailed  = "delete_failed"

	DedicatedHostStateCreatePending = createPending
	DedicatedHostStateCreating      = creating
	DedicatedHostStateCreated       = created
	DedicatedHostStateCreateFailed  = createFailed
	DedicatedHostStateDeleting      = deleting
	DedicatedHostStateDeleted       = deleted
	DedicatedHostStateDeleteFailed  = deleteFailed
)

func ResourceIBMContainerDedicatedHost() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceIBMContainerDedicatedHostCreate,
		ReadContext:   resourceIBMContainerDedicatedHostRead,
		UpdateContext: resourceIBMContainerDedicatedHostUpdate,
		DeleteContext: resourceIBMContainerDedicatedHostDelete,
		Importer:      &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(time.Minute * 40),
			Read:   schema.DefaultTimeout(time.Minute * 10),
			Update: schema.DefaultTimeout(time.Minute * 15),
			Delete: schema.DefaultTimeout(time.Minute * 40),
		},

		Schema: map[string]*schema.Schema{
			"flavor": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The flavor of the dedicated host",
			},
			"host_pool_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the dedicated host pool the dedicated host is associated with",
			},
			"zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The zone of the dedicated host",
			},
			"placement_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enables/disables placement on the dedicated host",
			},
			"host_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the dedicated host",
			},
			"life_cycle": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"actual_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"desired_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message_details": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message_details_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"resources": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The resources of the dedicated host",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"capacity": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"memory_bytes": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"vcpu": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"consumed": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"memory_bytes": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"vcpu": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"workers": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The workers of the dedicated host",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"flavor": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"worker_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"worker_pool_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceIBMContainerDedicatedHostCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return diag.FromErr(err)
	}
	dedicatedHostAPI := client.DedicatedHost()
	targetEnv := v2.ClusterTargetHeader{}

	hostPoolID := d.Get("host_pool_id").(string)

	params := v2.CreateDedicatedHostRequest{
		Flavor:     d.Get("flavor").(string),
		HostPoolID: hostPoolID,
		Zone:       d.Get("zone").(string),
	}

	res, err := dedicatedHostAPI.CreateDedicatedHost(params, targetEnv)
	if err != nil {
		return diag.Errorf("[ERROR] CreateDedicatedHost failed: %v", err)
	}
	hostID := res.ID
	d.SetId(fmt.Sprintf("%s/%s", hostPoolID, hostID))

	dh, err := waitForDedicatedHostAvailable(ctx, dedicatedHostAPI, hostID, hostPoolID, d.Timeout(schema.TimeoutCreate), targetEnv)
	if err != nil {
		return diag.Errorf("[ERROR] waitForDedicatedHostAvailable failed: %v", err)
	}

	dedicatedHost, ok := dh.(v2.GetDedicatedHostResponse)
	if !ok {
		return diag.Errorf("[ERROR] waitForDedicatedHostAvailable response is faulty: %v", dh)
	}

	setDedicatedHostFields(d, dedicatedHost)

	placement, ok := d.GetOk("placement_enabled")
	if ok && dedicatedHost.PlacementEnabled != placement.(bool) {
		req := v2.UpdateDedicatedHostPlacementRequest{
			HostPoolID: hostPoolID,
			HostID:     hostID,
		}
		if placement.(bool) {
			if err = dedicatedHostAPI.EnableDedicatedHostPlacement(req, targetEnv); err != nil {
				return diag.Errorf("[ERROR] EnableDedicatedHostPlacement failed: %v", err)
			}
		} else {
			if err = dedicatedHostAPI.DisableDedicatedHostPlacement(req, targetEnv); err != nil {
				return diag.Errorf("[ERROR] DisableDedicatedHostPlacement failed: %v", err)
			}
		}
		_, err = waitForDedicatedHostPlacement(ctx, dedicatedHostAPI, hostID, hostPoolID, placement.(bool), d.Timeout(schema.TimeoutCreate), targetEnv)
		if err != nil {
			return diag.Errorf("[ERROR] waitForDedicatedHostPlacement failed: %v", err)
		}
		d.Set("placement_enabled", placement)
	}

	return nil
}

func resourceIBMContainerDedicatedHostRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := d.Id()
	if err := getIBMContainerDedicatedHost(id, d, meta); err != nil {
		return diag.Errorf("[ERROR] getIBMContainerDedicatedHost failed: %v", err)
	}
	return nil
}

func getIBMContainerDedicatedHost(id string, d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}
	dedicatedHostAPI := client.DedicatedHost()
	targetEnv := v2.ClusterTargetHeader{}

	// <hostpoolid>:<hostid>
	m := strings.Split(id, "/")
	if len(m) < 2 || m[0] == "" || m[1] == "" {
		return fmt.Errorf("[ERROR] unexpected format of ID (%s), the expected format is <hostpoolid>/<hostid>", id)
	}
	hostPoolID := m[0]
	hostID := m[1]

	dedicatedHost, err := dedicatedHostAPI.GetDedicatedHost(hostID, hostPoolID, targetEnv)

	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("[ERROR] Error getting container dedicatedhost: %s", err)
	}

	d.Set("host_pool_id", hostPoolID)
	setDedicatedHostFields(d, dedicatedHost)

	return nil
}

func resourceIBMContainerDedicatedHostUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return diag.FromErr(err)
	}
	dedicatedHostAPI := client.DedicatedHost()
	targetEnv := v2.ClusterTargetHeader{}

	id := d.Id()
	// <hostpoolid>/<hostid>
	m := strings.Split(id, "/")
	if len(m) < 2 || m[0] == "" || m[1] == "" {
		return diag.Errorf("[ERROR] unexpected format of ID (%s), the expected format is <hostpoolid>/<hostid>", id)
	}
	hostPoolID := m[0]
	hostID := m[1]

	if d.HasChange("placement_enabled") {
		placement := d.Get("placement_enabled").(bool)
		req := v2.UpdateDedicatedHostPlacementRequest{
			HostPoolID: hostPoolID,
			HostID:     hostID,
		}
		if placement {
			if err = dedicatedHostAPI.EnableDedicatedHostPlacement(req, targetEnv); err != nil {
				return diag.Errorf("[ERROR] EnableDedicatedHostPlacement failed: %v", err)
			}
		} else {
			if err = dedicatedHostAPI.DisableDedicatedHostPlacement(req, targetEnv); err != nil {
				return diag.Errorf("[ERROR] DisableDedicatedHostPlacement failed: %v", err)
			}
		}
		_, err = waitForDedicatedHostPlacement(ctx, dedicatedHostAPI, hostID, hostPoolID, placement, d.Timeout(schema.TimeoutUpdate), targetEnv)
		if err != nil {
			return diag.Errorf("[ERROR] waitForDedicatedHostPlacement failed: %v", err)
		}
		d.Set("placement_enabled", placement)
	}
	return nil
}

func resourceIBMContainerDedicatedHostDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return diag.FromErr(err)
	}
	dedicatedHostAPI := client.DedicatedHost()
	targetEnv := v2.ClusterTargetHeader{}

	id := d.Id()
	// <hostpoolid>/<hostid>
	m := strings.Split(id, "/")
	if len(m) < 2 || m[0] == "" || m[1] == "" {
		return diag.Errorf("[ERROR] unexpected format of ID (%s), the expected format is <hostpoolid>/<hostid>", id)
	}
	hostPoolID := m[0]
	hostID := m[1]

	placementParams := v2.UpdateDedicatedHostPlacementRequest{
		HostID:     hostID,
		HostPoolID: hostPoolID,
	}

	if err = dedicatedHostAPI.DisableDedicatedHostPlacement(placementParams, targetEnv); err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				log.Printf("[DEBUG] DedicatedHostDelete: DisableDedicatedHostPlacement couldn't find dedicated host with host pool id %s and host id %s", hostPoolID, hostID)
				return nil
			}
		}
		return diag.Errorf("[ERROR] DisableDedicatedHostPlacement failed: %v", err)
	}

	_, err = waitForDedicatedHostPlacement(ctx, dedicatedHostAPI, hostID, hostPoolID, false, d.Timeout(schema.TimeoutDelete), targetEnv)
	if err != nil {
		return diag.Errorf("[ERROR] waitForDedicatedHostPlacement failed: %v", err)
	}

	params := v2.RemoveDedicatedHostRequest{
		HostID:     hostID,
		HostPoolID: hostPoolID,
	}

	if err = dedicatedHostAPI.RemoveDedicatedHost(params, targetEnv); err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				log.Printf("[DEBUG] RemoveDedicatedHost couldn't find dedicated host with host pool id %s and host id %s", hostPoolID, hostID)
				return nil
			}
		}
		return diag.Errorf("[ERROR] RemoveDedicatedHost failed: %v", err)
	}

	_, err = waitForDedicatedHostRemove(ctx, dedicatedHostAPI, hostID, hostPoolID, d.Timeout(schema.TimeoutUpdate), targetEnv)
	if err != nil {
		return diag.Errorf("[ERROR] waitForDedicatedHostRemove failed: %v", err)
	}

	return nil
}

func waitForDedicatedHostAvailable(ctx context.Context, dedicatedHostAPI v2.DedicatedHost, hostID, hostPoolID string, timeout time.Duration, target v2.ClusterTargetHeader) (interface{}, error) {

	log.Printf("[DEBUG] Waiting for the dedicated host (%s) for hostpool (%s) to be available.", hostID, hostPoolID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{DedicatedHostStateCreatePending, DedicatedHostStateCreating},
		Target:     []string{DedicatedHostStateCreated},
		Refresh:    dedicatedHostStateRefreshFunc(dedicatedHostAPI, hostID, hostPoolID, target),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForStateContext(ctx)
}

func waitForDedicatedHostRemove(ctx context.Context, dedicatedHostAPI v2.DedicatedHost, hostID, hostPoolID string, timeout time.Duration, target v2.ClusterTargetHeader) (interface{}, error) {

	log.Printf("[DEBUG] Waiting for the dedicated host (%s) for hostpool (%s) to be removed.", hostID, hostPoolID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{DedicatedHostStateCreated, DedicatedHostStateDeleting},
		Target:     []string{DedicatedHostStateDeleted},
		Refresh:    dedicatedHostStateRefreshFunc(dedicatedHostAPI, hostID, hostPoolID, target),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForStateContext(ctx)
}

func dedicatedHostStateRefreshFunc(dedicatedHostAPI v2.DedicatedHost, hostID, hostPoolID string, target v2.ClusterTargetHeader) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		dedicatedHost, err := dedicatedHostAPI.GetDedicatedHost(hostID, hostPoolID, target)
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] Error retrieving dedicated host: %s", err)

		}
		return dedicatedHost, dedicatedHost.Lifecycle.ActualState, nil
	}
}

func waitForDedicatedHostPlacement(ctx context.Context, dedicatedHostAPI v2.DedicatedHost, hostID, hostPoolID string, placement bool, timeout time.Duration, target v2.ClusterTargetHeader) (interface{}, error) {
	placementStr := strconv.FormatBool(placement)

	log.Printf("[DEBUG] Waiting for the dedicated host (%s) for hostpool (%s) placement to be %s.", hostID, hostPoolID, placementStr)

	pendingStr := strconv.FormatBool(!placement)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{pendingStr},
		Target:     []string{placementStr},
		Refresh:    dedicatedHostPlacementRefreshFunc(dedicatedHostAPI, hostID, hostPoolID, target),
		Timeout:    timeout,
		Delay:      2 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForStateContext(ctx)
}

func dedicatedHostPlacementRefreshFunc(dedicatedHostAPI v2.DedicatedHost, hostID, hostPoolID string, target v2.ClusterTargetHeader) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		dedicatedHost, err := dedicatedHostAPI.GetDedicatedHost(hostID, hostPoolID, target)
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] Error retrieving dedicated host: %s", err)

		}
		return dedicatedHost, strconv.FormatBool(dedicatedHost.PlacementEnabled), nil
	}
}

func setDedicatedHostFields(d *schema.ResourceData, dedicatedHost v2.GetDedicatedHostResponse) {
	d.Set("flavor", dedicatedHost.Flavor)
	d.Set("zone", dedicatedHost.Zone)
	d.Set("placement_enabled", dedicatedHost.PlacementEnabled)
	d.Set("host_id", dedicatedHost.ID)

	life_cycle := []interface{}{map[string]interface{}{
		"actual_state":         dedicatedHost.Lifecycle.ActualState,
		"desired_state":        dedicatedHost.Lifecycle.DesiredState,
		"message":              dedicatedHost.Lifecycle.Message,
		"message_date":         dedicatedHost.Lifecycle.MessageDate,
		"message_details":      dedicatedHost.Lifecycle.MessageDetails,
		"message_details_date": dedicatedHost.Lifecycle.MessageDetailsDate,
	}}
	d.Set("life_cycle", life_cycle)

	capacity := []interface{}{map[string]interface{}{
		"memory_bytes": dedicatedHost.Resources.Capacity.MemoryBytes,
		"vcpu":         dedicatedHost.Resources.Capacity.VCPU,
	}}
	consumed := []interface{}{map[string]interface{}{
		"memory_bytes": dedicatedHost.Resources.Consumed.MemoryBytes,
		"vcpu":         dedicatedHost.Resources.Consumed.VCPU,
	}}
	resources := []interface{}{map[string]interface{}{
		"capacity": capacity,
		"consumed": consumed,
	}}
	d.Set("resources", resources)

	workers := make([]map[string]interface{}, len(dedicatedHost.Workers))
	for i, w := range dedicatedHost.Workers {
		workers[i] = map[string]interface{}{
			"cluster_id":     w.ClusterID,
			"flavor":         w.Flavor,
			"worker_id":      w.WorkerID,
			"worker_pool_id": w.WorkerPoolID,
		}
	}
	d.Set("workers", workers)
}
