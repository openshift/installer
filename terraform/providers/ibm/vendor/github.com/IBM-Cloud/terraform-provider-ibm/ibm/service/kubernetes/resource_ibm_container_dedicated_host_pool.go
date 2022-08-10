// Copyright IBM Corp. 2017, 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package kubernetes

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
)

const (
	DedicatedHostPoolStateCreated  = created
	DedicatedHostPoolStateCreating = creating
	DedicatedHostPoolStateDeleting = deleting
	DedicatedHostPoolStateDeleted  = deleted
)

func ResourceIBMContainerDedicatedHostPool() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceIBMContainerDedicatedHostPoolCreate,
		ReadContext:   resourceIBMContainerDedicatedHostPoolRead,
		DeleteContext: resourceIBMContainerDedicatedHostPoolDelete,
		Importer:      &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(time.Minute * 10),
			Read:   schema.DefaultTimeout(time.Minute * 10),
			Delete: schema.DefaultTimeout(time.Minute * 10),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the dedicated host pool",
			},
			"metro": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The metro to create the dedicated host pool in",
			},
			"flavor_class": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The flavor class of the dedicated host pool",
			},
			"resource_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "ID of the resource group.",
			},
			"host_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The count of the hosts under the dedicated host pool",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The state of the dedicated host pool",
			},
			"zones": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The zones of the dedicated host pool",
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
						"host_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"worker_pools": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The worker pools of the dedicated host pool",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
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

func resourceIBMContainerDedicatedHostPoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return diag.FromErr(err)
	}
	dedicatedHostPoolAPI := client.DedicatedHostPool()
	targetEnv := v2.ClusterTargetHeader{}

	if rg, ok := d.GetOk("resource_group_id"); ok {
		targetEnv.ResourceGroup = rg.(string)
	}

	params := v2.CreateDedicatedHostPoolRequest{
		FlavorClass: d.Get("flavor_class").(string),
		Metro:       d.Get("metro").(string),
		Name:        d.Get("name").(string),
	}

	res, err := dedicatedHostPoolAPI.CreateDedicatedHostPool(params, targetEnv)
	if err != nil {
		return diag.Errorf("[ERROR] Error creating host pool %v", err)
	}

	d.SetId(res.ID)

	dhp, err := waitForDedicatedHostPoolAvailable(ctx, dedicatedHostPoolAPI, res.ID, d.Timeout(schema.TimeoutCreate), targetEnv)
	if err != nil {
		return diag.Errorf("[ERROR] waitForDedicatedHostPoolAvailable failed: %v", err)
	}

	setDedicatedHostPoolFields(dhp.(v2.GetDedicatedHostPoolResponse), d)

	return nil
}

func resourceIBMContainerDedicatedHostPoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := getIBMContainerDedicatedHostPool(d.Id(), d, meta)
	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				d.SetId("")
				return nil
			}
		}
		return diag.Errorf("[ERROR] Error retrieving host pool details %v", err)
	}
	return nil
}

func getIBMContainerDedicatedHostPool(hostPoolID string, d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}
	dedicatedHostPoolAPI := client.DedicatedHostPool()
	targetEnv := v2.ClusterTargetHeader{}

	dedicatedHostPool, err := dedicatedHostPoolAPI.GetDedicatedHostPool(hostPoolID, targetEnv)
	if err != nil {
		return err
	}

	setDedicatedHostPoolFields(dedicatedHostPool, d)

	return nil
}

func setDedicatedHostPoolFields(dedicatedHostPool v2.GetDedicatedHostPoolResponse, d *schema.ResourceData) {
	d.Set("name", dedicatedHostPool.Name)
	d.Set("metro", dedicatedHostPool.Metro)
	d.Set("flavor_class", dedicatedHostPool.FlavorClass)
	d.Set("host_count", dedicatedHostPool.HostCount)
	d.Set("state", dedicatedHostPool.State)

	zones := make([]map[string]interface{}, len(dedicatedHostPool.Zones))
	for i, zone := range dedicatedHostPool.Zones {
		zones[i] = map[string]interface{}{
			"capacity": []interface{}{map[string]interface{}{
				"memory_bytes": zone.Capacity.MemoryBytes,
				"vcpu":         zone.Capacity.VCPU,
			}},
			"host_count": zone.HostCount,
			"zone":       zone.Zone,
		}
	}
	d.Set("zones", zones)

	workerpools := make([]map[string]interface{}, len(dedicatedHostPool.WorkerPools))
	for i, wpool := range dedicatedHostPool.WorkerPools {
		workerpools[i] = map[string]interface{}{
			"cluster_id":     wpool.ClusterID,
			"worker_pool_id": wpool.WorkerPoolID,
		}
	}
	d.Set("worker_pools", workerpools)
}

func resourceIBMContainerDedicatedHostPoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return diag.FromErr(err)
	}
	dedicatedHostPoolAPI := client.DedicatedHostPool()
	targetEnv := v2.ClusterTargetHeader{}

	hostPoolID := d.Id()

	params := v2.RemoveDedicatedHostPoolRequest{
		HostPoolID: hostPoolID,
	}

	if err := dedicatedHostPoolAPI.RemoveDedicatedHostPool(params, targetEnv); err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				log.Printf("[DEBUG] RemoveDedicatedHostPool couldn't find dedicated host pool with id %s", hostPoolID)
				return nil
			}
		}
		return diag.Errorf("[ERROR] Error removing host pool %v", err)
	}

	_, err = waitForDedicatedHostPoolRemove(ctx, dedicatedHostPoolAPI, hostPoolID, d.Timeout(schema.TimeoutDelete), targetEnv)
	if err != nil {
		return diag.Errorf("[ERROR] waitForDedicatedHostPoolRemove failed: %v", err)
	}

	return nil
}

func waitForDedicatedHostPoolAvailable(ctx context.Context, dedicatedHostPoolAPI v2.DedicatedHostPool, hostPoolID string, timeout time.Duration, target v2.ClusterTargetHeader) (interface{}, error) {

	log.Printf("[DEBUG] Waiting for the dedicated hostpool (%s) to be available.", hostPoolID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{DedicatedHostPoolStateCreating},
		Target:     []string{DedicatedHostPoolStateCreated},
		Refresh:    dedicatedHostPoolStateRefreshFunc(dedicatedHostPoolAPI, hostPoolID, target),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForStateContext(ctx)
}

func waitForDedicatedHostPoolRemove(ctx context.Context, dedicatedHostPoolAPI v2.DedicatedHostPool, hostPoolID string, timeout time.Duration, target v2.ClusterTargetHeader) (interface{}, error) {

	log.Printf("[DEBUG] Waiting for the dedicated hostpool (%s) to be removed.", hostPoolID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{DedicatedHostPoolStateCreated, DedicatedHostPoolStateDeleting},
		Target:     []string{DedicatedHostPoolStateDeleted},
		Refresh:    dedicatedHostPoolStateRefreshFunc(dedicatedHostPoolAPI, hostPoolID, target),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForStateContext(ctx)
}

func dedicatedHostPoolStateRefreshFunc(dedicatedHostPoolAPI v2.DedicatedHostPool, hostPoolID string, target v2.ClusterTargetHeader) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		dedicatedHostPool, err := dedicatedHostPoolAPI.GetDedicatedHostPool(hostPoolID, target)
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] Error retrieving dedicated host pool: %s", err)

		}
		return dedicatedHostPool, dedicatedHostPool.State, nil
	}
}
