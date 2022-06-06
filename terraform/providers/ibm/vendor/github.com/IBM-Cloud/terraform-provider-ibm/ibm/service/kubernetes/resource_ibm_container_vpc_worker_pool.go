// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package kubernetes

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	v1 "github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
)

const (
	workerDesired = "deployed"
)

func ResourceIBMContainerVpcWorkerPool() *schema.Resource {

	return &schema.Resource{
		Create:   resourceIBMContainerVpcWorkerPoolCreate,
		Update:   resourceIBMContainerVpcWorkerPoolUpdate,
		Read:     resourceIBMContainerVpcWorkerPoolRead,
		Delete:   resourceIBMContainerVpcWorkerPoolDelete,
		Exists:   resourceIBMContainerVpcWorkerPoolExists,
		Importer: &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Delete: schema.DefaultTimeout(90 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"cluster": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cluster name",
			},

			"flavor": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "cluster node falvor",
			},

			"worker_pool_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "worker pool name",
			},

			"zones": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Zones info",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "zone name",
						},

						"subnet_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "subnet ID",
						},
					},
				},
			},

			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Labels",
			},

			"worker_pool_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"taints": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "WorkerPool Taints",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Key for taint",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Value for taint.",
						},
						"effect": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Effect for taint. Accepted values are NoSchedule, PreferNoSchedule and NoExecute.",
							ValidateFunc: validate.InvokeValidator(
								"ibm_container_vpc_worker_pool",
								"worker_taints"),
						},
					},
				},
			},

			"resource_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the resource group.",
				ForceNew:    true,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The vpc id where the cluster is",
				ForceNew:    true,
			},
			"worker_count": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The number of workers",
			},
			"entitlement": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "Entitlement option reduces additional OCP Licence cost in Openshift Clusters",
			},
			flex.ResourceControllerURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource Controller URL",
			},
		},
	}
}
func ResourceIBMContainerVPCWorkerPoolValidator() *validate.ResourceValidator {
	tainteffects := "NoSchedule,PreferNoSchedule,NoExecute"
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "worker_taints",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              tainteffects})

	containerVPCWorkerPoolTaintsValidator := validate.ResourceValidator{ResourceName: "ibm_container_vpc_worker_pool", Schema: validateSchema}
	return &containerVPCWorkerPoolTaintsValidator
}

func resourceIBMContainerVpcWorkerPoolCreate(d *schema.ResourceData, meta interface{}) error {

	wpClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}

	clusterNameorID := d.Get("cluster").(string)
	var zonei []interface{}

	zone := []v2.Zone{}

	if res, ok := d.GetOk("zones"); ok {
		zonei = res.(*schema.Set).List()
		for _, e := range zonei {
			r, _ := e.(map[string]interface{})
			zoneParam := v2.Zone{
				ID:       r["name"].(string),
				SubnetID: r["subnet_id"].(string),
			}
			zone = append(zone, zoneParam)
		}

	}

	// for _, e := range d.Get("zones").(*schema.Set).List() {
	// 	value := e.(map[string]interface{})
	// 	id := value["id"].(string)
	// 	subnetid := value["subnet_id"].(string)

	// }

	workerPoolConfig := v2.WorkerPoolConfig{
		Name:        d.Get("worker_pool_name").(string),
		VpcID:       d.Get("vpc_id").(string),
		Flavor:      d.Get("flavor").(string),
		WorkerCount: d.Get("worker_count").(int),
		Zones:       zone,
	}

	if l, ok := d.GetOk("labels"); ok {
		labels := make(map[string]string)
		for k, v := range l.(map[string]interface{}) {
			labels[k] = v.(string)
		}
		workerPoolConfig.Labels = labels
	}

	// Update workerpoolConfig with Entitlement option if provided
	if v, ok := d.GetOk("entitlement"); ok {
		workerPoolConfig.Entitlement = v.(string)
	}

	params := v2.WorkerPoolRequest{
		WorkerPoolConfig: workerPoolConfig,
		Cluster:          clusterNameorID,
	}

	workerPoolsAPI := wpClient.WorkerPools()
	targetEnv, err := getVpcClusterTargetHeader(d, meta)
	if err != nil {
		return err
	}

	res, err := workerPoolsAPI.CreateWorkerPool(params, targetEnv)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s", clusterNameorID, res.ID))

	//wait for workerpool availability
	_, err = WaitForWorkerPoolAvailable(d, meta, clusterNameorID, res.ID, d.Timeout(schema.TimeoutCreate), targetEnv)
	if err != nil {
		return fmt.Errorf("[ERROR] Error waiting for workerpool (%s) to become ready: %s", d.Id(), err)
	}

	return resourceIBMContainerVpcWorkerPoolUpdate(d, meta)
}

func resourceIBMContainerVpcWorkerPoolUpdate(d *schema.ResourceData, meta interface{}) error {

	if d.HasChange("labels") && !d.IsNewResource() {
		clusterNameOrID := d.Get("cluster").(string)
		workerPoolName := d.Get("worker_pool_name").(string)

		labels := make(map[string]string)
		if l, ok := d.GetOk("labels"); ok {
			for k, v := range l.(map[string]interface{}) {
				labels[k] = v.(string)
			}
		}

		targetEnv, err := getVpcClusterTargetHeader(d, meta)
		if err != nil {
			return err
		}
		ClusterClient, err := meta.(conns.ClientSession).ContainerAPI()
		if err != nil {
			return err
		}
		Env := v1.ClusterTargetHeader{ResourceGroup: targetEnv.ResourceGroup}

		err = ClusterClient.WorkerPools().UpdateLabelsWorkerPool(clusterNameOrID, workerPoolName, labels, Env)
		if err != nil {
			return fmt.Errorf("[ERROR] Error updating the labels: %s", err)
		}
	}
	if d.HasChange("taints") {
		clusterNameOrID := d.Get("cluster").(string)
		workerPoolName := d.Get("worker_pool_name").(string)
		taintParam := expandWorkerPoolTaints(d, meta, clusterNameOrID, workerPoolName)

		targetEnv, err := getVpcClusterTargetHeader(d, meta)
		if err != nil {
			return err
		}
		ClusterClient, err := meta.(conns.ClientSession).VpcContainerAPI()
		if err != nil {
			return err
		}
		err = ClusterClient.WorkerPools().UpdateWorkerPoolTaints(taintParam, targetEnv)
		if err != nil {
			return fmt.Errorf("[ERROR] Error updating the taints: %s", err)
		}
	}

	if d.HasChange("worker_count") {
		clusterNameOrID := d.Get("cluster").(string)
		workerPoolName := d.Get("worker_pool_name").(string)
		count := d.Get("worker_count").(int)
		targetEnv, err := getVpcClusterTargetHeader(d, meta)
		if err != nil {
			return err
		}
		ClusterClient, err := meta.(conns.ClientSession).ContainerAPI()
		if err != nil {
			return err
		}
		Env := v1.ClusterTargetHeader{ResourceGroup: targetEnv.ResourceGroup}

		err = ClusterClient.WorkerPools().ResizeWorkerPool(clusterNameOrID, workerPoolName, count, Env)
		if err != nil {
			return fmt.Errorf("[ERROR] Error updating the worker_count %d: %s", count, err)
		}
	}

	if d.HasChange("zones") && !d.IsNewResource() {
		clusterID := d.Get("cluster").(string)
		workerPoolName := d.Get("worker_pool_name").(string)
		targetEnv, err := getVpcClusterTargetHeader(d, meta)
		if err != nil {
			return err
		}
		oldList, newList := d.GetChange("zones")
		if oldList == nil {
			oldList = new(schema.Set)
		}
		if newList == nil {
			newList = new(schema.Set)
		}
		os := oldList.(*schema.Set)
		ns := newList.(*schema.Set)
		remove := os.Difference(ns).List()
		add := ns.Difference(os).List()
		if len(add) > 0 {
			csClient, err := meta.(conns.ClientSession).VpcContainerAPI()
			if err != nil {
				return err
			}
			for _, zone := range add {
				newZone := zone.(map[string]interface{})
				zoneParam := v2.WorkerPoolZone{
					Cluster:      clusterID,
					Id:           newZone["name"].(string),
					SubnetID:     newZone["subnet_id"].(string),
					WorkerPoolID: workerPoolName,
				}
				err = csClient.WorkerPools().CreateWorkerPoolZone(zoneParam, targetEnv)
				if err != nil {
					return fmt.Errorf("[ERROR] Error adding zone to conatiner vpc cluster: %s", err)
				}
				_, err = WaitForWorkerPoolAvailable(d, meta, clusterID, workerPoolName, d.Timeout(schema.TimeoutCreate), targetEnv)
				if err != nil {
					return fmt.Errorf("[ERROR] Error waiting for workerpool (%s) to become ready: %s", d.Id(), err)
				}

			}
		}
		if len(remove) > 0 {
			for _, zone := range remove {
				oldZone := zone.(map[string]interface{})
				ClusterClient, err := meta.(conns.ClientSession).ContainerAPI()
				if err != nil {
					return err
				}
				Env := v1.ClusterTargetHeader{ResourceGroup: targetEnv.ResourceGroup}
				err = ClusterClient.WorkerPools().RemoveZone(clusterID, oldZone["name"].(string), workerPoolName, Env)
				if err != nil {
					return fmt.Errorf("[ERROR] Error deleting zone to conatiner vpc cluster: %s", err)
				}
				_, err = WaitForV2WorkerZoneDeleted(clusterID, workerPoolName, oldZone["name"].(string), meta, d.Timeout(schema.TimeoutDelete), targetEnv)
				if err != nil {
					return fmt.Errorf("[ERROR] Error waiting for deleting workers of worker pool (%s) of cluster (%s):  %s", workerPoolName, clusterID, err)
				}
			}
		}
	}
	return resourceIBMContainerVpcWorkerPoolRead(d, meta)
}

func expandWorkerPoolTaints(d *schema.ResourceData, meta interface{}, clusterNameOrID, workerPoolName string) v2.WorkerPoolTaintRequest {
	taintBody := make(map[string]string)
	if res, ok := d.GetOk("taints"); ok {
		taints := res.(*schema.Set).List()
		for _, t := range taints {
			r, _ := t.(map[string]interface{})
			key := r["key"].(string)
			value := r["value"].(string)
			effect := r["effect"].(string)
			taintBody[key] = fmt.Sprintf("%s:%s", value, effect)
		}
	}
	taintParam := v2.WorkerPoolTaintRequest{
		Cluster:    clusterNameOrID,
		WorkerPool: workerPoolName,
		Taints:     taintBody,
	}
	return taintParam
}
func flattenWorkerPoolTaints(taints v2.GetWorkerPoolResponse) []map[string]interface{} {
	taintslist := make([]map[string]interface{}, 0)
	for k, v := range taints.Taints {
		taint := make(map[string]interface{})
		taint["key"] = k
		ve := strings.Split(v, ":")
		taint["value"] = ve[0]
		taint["effect"] = ve[1]
		taintslist = append(taintslist, taint)
	}
	return taintslist
}
func resourceIBMContainerVpcWorkerPoolRead(d *schema.ResourceData, meta interface{}) error {
	wpClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}
	cluster := parts[0]
	workerPoolID := parts[1]

	workerPoolsAPI := wpClient.WorkerPools()
	targetEnv, err := getVpcClusterTargetHeader(d, meta)
	if err != nil {
		return err
	}

	workerPool, err := workerPoolsAPI.GetWorkerPool(cluster, workerPoolID, targetEnv)
	if err != nil {
		return err
	}

	var zones = make([]map[string]interface{}, 0)
	for _, zone := range workerPool.Zones {
		for _, subnet := range zone.Subnets {
			zoneInfo := map[string]interface{}{
				"name":      zone.ID,
				"subnet_id": subnet.ID,
			}
			zones = append(zones, zoneInfo)
		}
	}

	cls, err := wpClient.Clusters().GetCluster(cluster, targetEnv)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving conatiner vpc cluster: %s", err)
	}

	d.Set("worker_pool_name", workerPool.PoolName)
	d.Set("flavor", workerPool.Flavor)
	d.Set("worker_count", workerPool.WorkerCount)
	d.Set("worker_pool_id", workerPoolID)
	// d.Set("provider", workerPool.Provider)
	d.Set("labels", flex.IgnoreSystemLabels(workerPool.Labels))
	d.Set("zones", zones)
	d.Set("resource_group_id", cls.ResourceGroupID)
	d.Set("cluster", cluster)
	d.Set("vpc_id", workerPool.VpcID)
	if workerPool.Taints != nil {
		d.Set("taints", flattenWorkerPoolTaints(workerPool))
	}
	controller, err := flex.GetBaseController(meta)
	if err != nil {
		return err
	}
	d.Set(flex.ResourceControllerURL, controller+"/kubernetes/clusters")
	return nil
}

func resourceIBMContainerVpcWorkerPoolDelete(d *schema.ResourceData, meta interface{}) error {
	wpClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}
	clusterNameorID := parts[0]
	workerPoolNameorID := parts[1]

	workerPoolsAPI := wpClient.WorkerPools()
	targetEnv, err := getVpcClusterTargetHeader(d, meta)
	if err != nil {
		return err
	}

	err = workerPoolsAPI.DeleteWorkerPool(clusterNameorID, workerPoolNameorID, targetEnv)
	if err != nil {
		return err
	}
	_, err = WaitForVpcWorkerDelete(clusterNameorID, workerPoolNameorID, meta, d.Timeout(schema.TimeoutDelete), targetEnv)
	if err != nil {
		return fmt.Errorf("[ERROR] Error waiting for removing workers of worker pool (%s) of cluster (%s): %s", workerPoolNameorID, clusterNameorID, err)
	}
	d.SetId("")
	return nil
}

func resourceIBMContainerVpcWorkerPoolExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	wpClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return false, err
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return false, err
	}
	if len(parts) < 2 {
		return false, fmt.Errorf("[ERROR] Incorrect ID %s: Id should be a combination of clusterID/WorkerPoolID", d.Id())
	}
	cluster := parts[0]
	workerPoolID := parts[1]

	workerPoolsAPI := wpClient.WorkerPools()
	targetEnv, err := getVpcClusterTargetHeader(d, meta)
	if err != nil {
		return false, err
	}

	workerPool, err := workerPoolsAPI.GetWorkerPool(cluster, workerPoolID, targetEnv)
	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 && strings.Contains(apiErr.Description(), "The specified worker pool could not be found") {
				return false, nil
			}
		}
		return false, fmt.Errorf("[ERROR] Error getting container vpc workerpool: %s", err)
	}
	if strings.Compare(workerPool.Lifecycle.ActualState, "deleted") == 0 {
		return false, nil
	}

	return workerPool.ID == workerPoolID, nil
}

// WaitForWorkerPoolAvailable Waits for worker creation
func WaitForWorkerPoolAvailable(d *schema.ResourceData, meta interface{}, clusterNameOrID, workerPoolNameOrID string, timeout time.Duration, target v2.ClusterTargetHeader) (interface{}, error) {
	wpClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return nil, err
	}
	log.Printf("Waiting for workerpool (%s) to be available.", d.Id())
	// id := d.Id()

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"provision_pending"},
		Target:     []string{workerDesired},
		Refresh:    vpcWorkerPoolStateRefreshFunc(wpClient.Workers(), clusterNameOrID, workerPoolNameOrID, target),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func vpcWorkerPoolStateRefreshFunc(client v2.Workers, instanceID string, workerPoolNameOrID string, target v2.ClusterTargetHeader) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		workerFields, err := client.ListByWorkerPool(instanceID, "", false, target)
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] Error retrieving workers for cluster: %s", err)
		}
		// Check active transactions
		//Check for worker state to be deployed
		//Done worker has two fields desiredState and actualState , so check for those 2
		for _, e := range workerFields {
			if e.PoolName == workerPoolNameOrID || e.PoolID == workerPoolNameOrID {
				if strings.Compare(e.LifeCycle.ActualState, "deployed") != 0 {
					log.Printf("worker: %s state: %s", e.ID, e.LifeCycle.ActualState)
					return workerFields, "provision_pending", nil
				}
			}
		}
		return workerFields, workerDesired, nil
	}
}

func WaitForVpcWorkerDelete(clusterNameOrID, workerPoolNameOrID string, meta interface{}, timeout time.Duration, target v2.ClusterTargetHeader) (interface{}, error) {
	wpClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return nil, err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"deleting"},
		Target:     []string{workerDeleteState},
		Refresh:    vpcworkerPoolDeleteStateRefreshFunc(wpClient.Workers(), clusterNameOrID, workerPoolNameOrID, target),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func vpcworkerPoolDeleteStateRefreshFunc(client v2.Workers, instanceID, workerPoolNameOrID string, target v2.ClusterTargetHeader) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		workerFields, err := client.ListByWorkerPool(instanceID, "", true, target)
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] Error retrieving workers for cluster: %s", err)
		}
		//Done worker has two fields desiredState and actualState , so check for those 2
		for _, e := range workerFields {
			if e.PoolName == workerPoolNameOrID || e.PoolID == workerPoolNameOrID {
				if strings.Compare(e.LifeCycle.ActualState, "deleted") != 0 {
					log.Printf("Deleting worker %s", e.ID)
					return workerFields, "deleting", nil
				}
			}
		}
		return workerFields, workerDeleteState, nil
	}
}
