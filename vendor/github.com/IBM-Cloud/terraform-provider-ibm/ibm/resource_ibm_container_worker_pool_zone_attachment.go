// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"strings"
	"time"

	v1 "github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIBMContainerWorkerPoolZoneAttachment() *schema.Resource {

	return &schema.Resource{
		Create:   resourceIBMContainerWorkerPoolZoneAttachmentCreate,
		Read:     resourceIBMContainerWorkerPoolZoneAttachmentRead,
		Update:   resourceIBMContainerWorkerPoolZoneAttachmentUpdate,
		Delete:   resourceIBMContainerWorkerPoolZoneAttachmentDelete,
		Exists:   resourceIBMContainerWorkerPoolZoneAttachmentExists,
		Importer: &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Update: schema.DefaultTimeout(90 * time.Minute),
			Delete: schema.DefaultTimeout(90 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone name",
			},

			"cluster": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "cluster name or ID",
			},

			"worker_pool": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Workerpool name",
			},

			"private_vlan_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"public_vlan_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"resource_group_id": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "ID of the resource group.",
				ForceNew:         true,
				DiffSuppressFunc: applyOnce,
			},

			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The zone region",
				Deprecated:  "This field is deprecated",
			},

			"worker_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"wait_till_albs": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          true,
				DiffSuppressFunc: applyOnce,
				Description:      "wait_till_albs can be configured to wait for albs during the worker pool zone attachment.",
			},
		},
	}
}

func resourceIBMContainerWorkerPoolZoneAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	csClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return err
	}

	zone := d.Get("zone").(string)
	var privateVLAN, publicVLAN string
	if v, ok := d.GetOk("private_vlan_id"); ok {
		privateVLAN = v.(string)
	}

	if v, ok := d.GetOk("public_vlan_id"); ok {
		publicVLAN = v.(string)
	}

	if publicVLAN != "" && privateVLAN == "" {
		return fmt.Errorf(
			"A private_vlan_id must be specified if a public_vlan_id is specified.")
	}

	workerPoolZoneNetwork := v1.WorkerPoolZoneNetwork{
		PrivateVLAN: privateVLAN,
		PublicVLAN:  publicVLAN,
	}

	workerPoolZone := v1.WorkerPoolZone{
		ID:                    zone,
		WorkerPoolZoneNetwork: workerPoolZoneNetwork,
	}

	cluster := d.Get("cluster").(string)
	workerPool := d.Get("worker_pool").(string)

	workerPoolsAPI := csClient.WorkerPools()
	targetEnv, err := getWorkerPoolTargetHeader(d, meta)
	if err != nil {
		return err
	}

	err = workerPoolsAPI.AddZone(cluster, workerPool, workerPoolZone, targetEnv)
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%s/%s/%s", cluster, workerPool, zone))

	_, err = WaitForWorkerZoneNormal(cluster, workerPool, zone, meta, d.Timeout(schema.TimeoutUpdate), targetEnv)
	if err != nil {
		return fmt.Errorf(
			"Error waiting for workers of worker pool (%s) of cluster (%s) to become ready: %s", workerPool, cluster, err)
	}

	var waitTillALBs bool
	if v, ok := d.GetOk("wait_till_albs"); ok {
		waitTillALBs = v.(bool)
	}

	if waitTillALBs {
		_, err = waitForWorkerZoneALB(cluster, zone, meta, d.Timeout(schema.TimeoutUpdate), targetEnv)
		if err != nil {
			return fmt.Errorf(
				"Error waiting for ALBs in zone (%s) of cluster (%s) to become ready: %s", zone, cluster, err)
		}
	}

	return resourceIBMContainerWorkerPoolZoneAttachmentRead(d, meta)

}

func resourceIBMContainerWorkerPoolZoneAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	csClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return err
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	cluster := parts[0]
	workerPool := parts[1]
	zoneName := parts[2]

	workerPoolsAPI := csClient.WorkerPools()
	targetEnv, err := getWorkerPoolTargetHeader(d, meta)
	if err != nil {
		return err
	}

	workerPoolRes, err := workerPoolsAPI.GetWorkerPool(cluster, workerPool, targetEnv)
	if err != nil {
		return err
	}
	zones := workerPoolRes.Zones

	for _, zone := range zones {
		if zone.ID == zoneName {
			d.Set("public_vlan_id", zone.PublicVLAN)
			d.Set("private_vlan_id", zone.PrivateVLAN)
			d.Set("worker_count", zone.WorkerCount)
			d.Set("zone", zone.ID)
			d.Set("cluster", cluster)
			d.Set("worker_pool", workerPool)

			break
		}
	}

	return nil
}

func resourceIBMContainerWorkerPoolZoneAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	csClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return err
	}

	workerPoolsAPI := csClient.WorkerPools()

	if d.HasChange("private_vlan_id") || d.HasChange("public_vlan_id") {
		privateVLAN := d.Get("private_vlan_id").(string)
		publicVLAN := d.Get("public_vlan_id").(string)
		if publicVLAN != "" && privateVLAN == "" {
			return fmt.Errorf(
				"A private VLAN must be specified if a public VLAN is specified.")
		}
		targetEnv, err := getWorkerPoolTargetHeader(d, meta)
		if err != nil {
			return err
		}
		parts, err := idParts(d.Id())
		if err != nil {
			return err
		}
		cluster := parts[0]
		workerPool := parts[1]
		zone := parts[2]
		err = workerPoolsAPI.UpdateZoneNetwork(cluster, zone, workerPool, privateVLAN, publicVLAN, targetEnv)
		if err != nil {
			return err
		}
	}

	return resourceIBMContainerWorkerPoolZoneAttachmentRead(d, meta)
}

func resourceIBMContainerWorkerPoolZoneAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	csClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return err
	}

	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	cluster := parts[0]
	workerPool := parts[1]
	zone := parts[2]

	workerPoolsAPI := csClient.WorkerPools()
	targetEnv, err := getWorkerPoolTargetHeader(d, meta)
	if err != nil {
		return err
	}
	err = workerPoolsAPI.RemoveZone(cluster, zone, workerPool, targetEnv)
	if err != nil {
		return err
	}
	_, err = WaitForWorkerZoneDeleted(cluster, workerPool, zone, meta, d.Timeout(schema.TimeoutDelete), targetEnv)
	if err != nil {
		return fmt.Errorf(
			"Error waiting for deleting workers of worker pool (%s) of cluster (%s):  %s", workerPool, cluster, err)
	}

	return nil
}

func resourceIBMContainerWorkerPoolZoneAttachmentExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	csClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return false, err
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return false, err
	}
	cluster := parts[0]
	workerPoolID := parts[1]
	zoneID := parts[2]

	workerPoolsAPI := csClient.WorkerPools()
	targetEnv, err := getWorkerPoolTargetHeader(d, meta)
	if err != nil {
		return false, err
	}

	workerPool, err := workerPoolsAPI.GetWorkerPool(cluster, workerPoolID, targetEnv)
	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error communicating with the API: %s", err)
	}
	zones := workerPool.Zones
	var zone v1.WorkerPoolZoneResponse
	for _, z := range zones {
		if z.ID == zoneID {
			zone = z
		}
	}
	return zone.ID == zoneID, nil
}

func WaitForWorkerZoneNormal(clusterNameOrID, workerPoolNameOrID, zone string, meta interface{}, timeout time.Duration, target v1.ClusterTargetHeader) (interface{}, error) {
	csClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return nil, err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", workerProvisioning},
		Target:     []string{workerNormal},
		Refresh:    workerPoolZoneStateRefreshFunc(csClient.Workers(), clusterNameOrID, workerPoolNameOrID, zone, target),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func workerPoolZoneStateRefreshFunc(client v1.Workers, instanceID, workerPoolNameOrID, zone string, target v1.ClusterTargetHeader) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		workerFields, err := client.ListByWorkerPool(instanceID, workerPoolNameOrID, false, target)
		if err != nil {
			return nil, "", fmt.Errorf("Error retrieving workers for cluster: %s", err)
		}
		//Done worker has two fields State and Status , so check for those 2
		for _, e := range workerFields {
			if e.Location == zone {
				if strings.Contains(e.KubeVersion, "pending") || strings.Compare(e.State, workerNormal) != 0 || strings.Compare(e.Status, workerReadyState) != 0 {
					if strings.Compare(e.State, "deleted") != 0 {
						return workerFields, workerProvisioning, nil
					}
				}
			}
		}
		return workerFields, workerNormal, nil
	}
}

func WaitForWorkerZoneDeleted(clusterNameOrID, workerPoolNameOrID, zone string, meta interface{}, timeout time.Duration, target v1.ClusterTargetHeader) (interface{}, error) {
	csClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return nil, err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"deleting"},
		Target:     []string{workerDeleteState},
		Refresh:    workerPoolZoneDeleteStateRefreshFunc(csClient.Workers(), clusterNameOrID, workerPoolNameOrID, zone, target),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func workerPoolZoneDeleteStateRefreshFunc(client v1.Workers, instanceID, workerPoolNameOrID, zone string, target v1.ClusterTargetHeader) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		workerFields, err := client.ListByWorkerPool(instanceID, workerPoolNameOrID, true, target)
		if err != nil {
			return nil, "", fmt.Errorf("Error retrieving workers for cluster: %s", err)
		}
		//Done worker has two fields State and Status , so check for those 2
		for _, e := range workerFields {
			if e.Location == zone {
				if strings.Compare(e.State, "deleted") != 0 {
					return workerFields, "deleting", nil
				}
			}
		}
		return workerFields, workerDeleteState, nil
	}
}

func waitForWorkerZoneALB(clusterNameOrID, zone string, meta interface{}, timeout time.Duration, target v1.ClusterTargetHeader) (interface{}, error) {
	csClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return nil, err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"ready"},
		Refresh:    workerZoneALBStateRefreshFunc(csClient.Albs(), clusterNameOrID, zone, target),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func workerZoneALBStateRefreshFunc(client v1.Albs, instanceID, zone string, target v1.ClusterTargetHeader) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		// Get all ALBs associated with cluster
		albs, err := client.ListClusterALBs(instanceID, target)
		if err != nil {
			return nil, "", fmt.Errorf("Error retrieving ALBs for cluster: %s", err)
		}

		privateALBsByZone := []v1.ALBConfig{}
		publicALBsByZone := []v1.ALBConfig{}

		// Find ALBs by zone and type
		for _, alb := range albs {
			if alb.Zone == zone {
				if alb.ALBType == "private" {
					privateALBsByZone = append(privateALBsByZone, alb)
				}
				if alb.ALBType == "public" {
					publicALBsByZone = append(publicALBsByZone, alb)
				}
			}
		}

		// Ready if both private and public ALBs are present
		if len(privateALBsByZone) > 0 && len(publicALBsByZone) > 0 {
			return albs, "ready", nil
		}

		return albs, "pending", nil
	}
}
