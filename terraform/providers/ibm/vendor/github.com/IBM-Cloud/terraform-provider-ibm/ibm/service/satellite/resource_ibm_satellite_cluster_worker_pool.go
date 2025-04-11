// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package satellite

import (
	"fmt"
	"log"
	"strings"
	"time"

	v1 "github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
	"github.com/IBM-Cloud/container-services-go-sdk/kubernetesserviceapiv1"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	workerPoolDesired    = "deployed"
	clusterNormal        = "normal"
	clusterDeletePending = "deleting"
	clusterDeleted       = "deleted"
	workerNormal         = "normal"
	subnetNormal         = "normal"
	workerReadyState     = "Ready"
	workerDeleteState    = "deleted"
	workerDeletePending  = "deleting"

	versionUpdating     = "updating"
	clusterProvisioning = "provisioning"
	workerProvisioning  = "provisioning"
	subnetProvisioning  = "provisioning"
)

func ResourceIBMSatelliteClusterWorkerPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceIBMSatelliteClusterWorkerPoolCreate,
		Read:   resourceIBMSatelliteClusterWorkerPoolRead,
		Update: resourceIBMSatelliteClusterWorkerPoolUpdate,
		Delete: resourceIBMSatelliteClusterWorkerPoolDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				parts, err := flex.IdParts(d.Id())
				if err != nil {
					return nil, err
				}
				clusterID := parts[0]
				workerPoolID := parts[1]

				satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
				if err != nil {
					return nil, err
				}
				getWorkerPoolOptions := &kubernetesserviceapiv1.GetWorkerPoolOptions{
					Cluster:    &clusterID,
					Workerpool: &workerPoolID,
				}
				workerPool, response, err := satClient.GetWorkerPool(getWorkerPoolOptions)
				if err != nil {
					return nil, fmt.Errorf("Error reading satellite worker pool: %s\n%s", response, err)
				}

				var zones = make([]map[string]interface{}, 0)
				for _, zone := range workerPool.Zones {
					zoneInfo := map[string]interface{}{
						"id": *zone.ID,
					}
					zones = append(zones, zoneInfo)
				}

				d.Set("zones", zones)
				return []*schema.ResourceData{d}, nil
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Delete: schema.DefaultTimeout(90 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name for the worker pool",
			},
			"cluster": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The unique name for the new IBM Cloud Satellite cluster",
			},
			"flavor": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The flavor defines the amount of virtual CPU, memory, and disk space that is set up in each worker node",
			},
			"disk_encryption": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Disk encryption for worker node",
			},
			"isolation": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"entitlement": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "Entitlement option reduces additional OCP Licence cost in Openshift Clusters",
			},
			"operating_system": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Operating system of the worker pool. Options are REDHAT_7_64, REDHAT_8_64, or RHCOS.",
			},
			"worker_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Specify the desired number of workers per zone in this worker pool",
			},
			"zones": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "Zone info for worker pool",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Zone for the worker pool in a multizone cluster",
						},
					},
				},
			},
			"worker_pool_labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Labels on all the workers in the worker pool",
			},
			"host_labels": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         flex.ResourceIBMVPCHash,
				Description: "Labels that describe a Satellite host",
			},
			"resource_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the resource group.",
			},
		},
	}
}
func getClusterTargetHeader(d *schema.ResourceData, meta interface{}) (v1.ClusterTargetHeader, error) {
	_, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return v1.ClusterTargetHeader{}, err
	}

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return v1.ClusterTargetHeader{}, err
	}
	accountID := userDetails.UserAccount

	targetEnv := v1.ClusterTargetHeader{
		AccountID: accountID,
	}

	resourceGroup := ""
	if v, ok := d.GetOk("resource_group_id"); ok {
		resourceGroup = v.(string)
		targetEnv.ResourceGroup = resourceGroup
	}
	return targetEnv, nil
}
func resourceIBMSatelliteClusterWorkerPoolCreate(d *schema.ResourceData, meta interface{}) error {
	satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
	if err != nil {
		return err
	}

	targetEnv, err := getClusterTargetHeader(d, meta)
	if err != nil {
		return err
	}

	createWorkerPoolOptions := &kubernetesserviceapiv1.CreateSatelliteWorkerPoolOptions{}
	name := d.Get("name").(string)
	createWorkerPoolOptions.Name = &name

	cluster := d.Get("cluster").(string)
	createWorkerPoolOptions.Cluster = &cluster

	if v, ok := d.GetOk("resource_group_id"); ok {
		pathParamsMap := map[string]string{
			"X-Auth-Resource-Group": v.(string),
		}
		createWorkerPoolOptions.Headers = pathParamsMap
	}

	if v, ok := d.GetOk("operating_system"); ok {
		operating_system := v.(string)
		createWorkerPoolOptions.OperatingSystem = &operating_system
	}

	if v, ok := d.GetOk("worker_count"); ok {
		workerCount := int64(v.(int))
		createWorkerPoolOptions.WorkerCount = &workerCount
	}

	if v, ok := d.GetOk("zones"); ok {
		z := v.(*schema.Set)
		createWorkerPoolOptions.Zones = flex.FlattenSatelliteWorkerPoolZones(z)
	}

	hostLabels := make(map[string]string)
	if v, ok := d.GetOk("host_labels"); ok {
		hl := v.(*schema.Set)
		hostLabels = flex.FlattenKeyValues(hl.List())
		createWorkerPoolOptions.HostLabels = hostLabels
	} else {
		createWorkerPoolOptions.HostLabels = hostLabels
	}

	labels := make(map[string]string)
	if l, ok := d.GetOk("worker_pool_labels"); ok {
		for k, v := range l.(map[string]interface{}) {
			labels[k] = v.(string)
		}
		createWorkerPoolOptions.Labels = labels
	} else {
		createWorkerPoolOptions.Labels = labels
	}

	if v, ok := d.GetOk("flavor"); ok {
		flavor := v.(string)
		createWorkerPoolOptions.Flavor = &flavor
	}

	if v, ok := d.GetOk("disk_encryption"); ok {
		diskEncryption := v.(bool)
		createWorkerPoolOptions.DiskEncryption = &diskEncryption
	}

	if v, ok := d.GetOk("isolation"); ok {
		isolation := v.(string)
		createWorkerPoolOptions.Isolation = &isolation
	}

	if v, ok := d.GetOk("entitlement"); ok {
		entitlement := v.(string)
		createWorkerPoolOptions.Entitlement = &entitlement
	}

	instance, response, err := satClient.CreateSatelliteWorkerPool(createWorkerPoolOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Creating Satellite cluster worker pool: %s\n%s", err, response)
	}

	d.SetId(fmt.Sprintf("%s/%s", cluster, *instance.WorkerPoolID))
	log.Printf("[INFO] Created satellite cluster worker pool: %s", *instance.WorkerPoolID)

	_, err = WaitForSatelliteWorkerPoolAvailable(d, meta, cluster, *instance.WorkerPoolID, d.Timeout(schema.TimeoutCreate), targetEnv)
	if err != nil {
		return fmt.Errorf("[ERROR] Error waiting for workerpool (%s) to become ready: %s", d.Id(), err)
	}

	return resourceIBMSatelliteClusterWorkerPoolRead(d, meta)
}

func resourceIBMSatelliteClusterWorkerPoolRead(d *schema.ResourceData, meta interface{}) error {
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}
	clusterID := parts[0]
	workerPoolID := parts[1]

	satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
	if err != nil {
		return err
	}

	getWorkerPoolOptions := &kubernetesserviceapiv1.GetWorkerPoolOptions{
		Cluster:    &clusterID,
		Workerpool: &workerPoolID,
	}

	workerPool, response, err := satClient.GetWorkerPool(getWorkerPoolOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", workerPool.PoolName)
	d.Set("cluster", clusterID)
	d.Set("flavor", workerPool.Flavor)
	d.Set("isolation", workerPool.Isolation)
	d.Set("operating_system", workerPool.OperatingSystem)
	d.Set("worker_count", workerPool.WorkerCount)
	d.Set("worker_pool_labels", flex.IgnoreSystemLabels(workerPool.Labels))
	d.Set("host_labels", flex.FlattenWorkerPoolHostLabels(workerPool.HostLabels))

	return nil
}

func resourceIBMSatelliteClusterWorkerPoolUpdate(d *schema.ResourceData, meta interface{}) error {
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}
	clusterNameOrID := parts[0]
	workerPoolName := parts[1]

	satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
	if err != nil {
		return err
	}

	targetEnv, err := getClusterTargetHeader(d, meta)
	if err != nil {
		return err
	}

	if d.HasChange("worker_pool_labels") {
		labels := make(map[string]string)
		if l, ok := d.GetOk("worker_pool_labels"); ok {
			for k, v := range l.(map[string]interface{}) {
				labels[k] = v.(string)
			}
		}

		wpots := &kubernetesserviceapiv1.V2SetWorkerPoolLabelsOptions{
			Cluster:    &clusterNameOrID,
			Workerpool: &workerPoolName,
			Labels:     labels,
		}
		response, err := satClient.V2SetWorkerPoolLabels(wpots)
		if err != nil {
			return fmt.Errorf("[ERROR] Error updating the labels: %s\n%s", err, response)
		}
	}

	if d.HasChange("worker_count") {
		clusterNameOrID := d.Get("cluster").(string)
		workerPoolName := d.Get("name").(string)
		count := d.Get("worker_count").(int)
		targetEnv, err := getVpcClusterTargetHeader(d)
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

	if d.HasChange("zones") {
		clusterID := d.Get("cluster").(string)
		workerPoolName := d.Get("name").(string)

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
			for _, zone := range add {
				newZone := zone.(map[string]interface{})
				zID := newZone["id"].(string)
				zoneOptions := &kubernetesserviceapiv1.CreateSatelliteWorkerPoolZoneOptions{
					Cluster:    &clusterID,
					Workerpool: &workerPoolName,
					ID:         &zID,
				}
				response, err := satClient.CreateSatelliteWorkerPoolZone(zoneOptions)
				if err != nil {
					return fmt.Errorf("[ERROR] Error Adding Worker Pool Zone : %s\n%s", err, response)
				}
			}
			_, err = WaitForSatelliteWorkerPoolAvailable(d, meta, clusterID, workerPoolName, d.Timeout(schema.TimeoutCreate), targetEnv)
			if err != nil {
				return fmt.Errorf("[ERROR] Error waiting for workerpool (%s) to become ready: %s", d.Id(), err)
			}
		}
		if len(remove) > 0 {
			for _, zone := range remove {
				oldZone := zone.(map[string]interface{})
				zID := oldZone["id"].(string)
				zoneOptions := &kubernetesserviceapiv1.RemoveWorkerPoolZoneOptions{
					IdOrName:     &clusterID,
					PoolidOrName: &workerPoolName,
					Zoneid:       &zID,
				}
				response, err := satClient.RemoveWorkerPoolZone(zoneOptions)
				if err != nil {
					return fmt.Errorf("[ERROR] Error deleting Worker Pool Zone : %s\n%s", err, response)
				}
			}
		}
	}
	return resourceIBMSatelliteClusterWorkerPoolRead(d, meta)
}

func resourceIBMSatelliteClusterWorkerPoolDelete(d *schema.ResourceData, meta interface{}) error {
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}
	clusterID := parts[0]
	workerPoolID := parts[1]

	satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
	if err != nil {
		return err
	}

	targetEnv, err := getVpcClusterTargetHeader(d)
	if err != nil {
		return err
	}

	wpOptions := &kubernetesserviceapiv1.RemoveWorkerPoolOptions{
		IdOrName:     &clusterID,
		PoolidOrName: &workerPoolID,
	}

	response, err := satClient.RemoveWorkerPool(wpOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil
		}
		return fmt.Errorf("[ERROR] Error Deleting Satellite Cluster WorkerPool: %s\n%s", err, response)
	}

	_, err = WaitForSatelliteWorkerDelete(clusterID, workerPoolID, meta, d.Timeout(schema.TimeoutDelete), targetEnv)
	if err != nil {
		return fmt.Errorf("[ERROR] Error waiting for removing workers of worker pool (%s) of cluster (%s): %s", workerPoolID, clusterID, err)
	}

	d.SetId("")
	return nil
}
func getVpcClusterTargetHeader(d *schema.ResourceData) (v2.ClusterTargetHeader, error) {
	targetEnv := v2.ClusterTargetHeader{}
	var resourceGroup string
	if rg, ok := d.GetOk("resource_group_id"); ok {
		resourceGroup = rg.(string)
		targetEnv.ResourceGroup = resourceGroup
	}

	return targetEnv, nil
}

// WaitForSatelliteWorkerPoolAvailable Waits for workerpool deployed
func WaitForSatelliteWorkerPoolAvailable(d *schema.ResourceData, meta interface{}, clusterNameOrID, workerPoolNameOrID string, timeout time.Duration, target v1.ClusterTargetHeader) (interface{}, error) {
	clusterID := clusterNameOrID
	workerPoolID := workerPoolNameOrID

	satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
	if err != nil {
		return nil, err
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"provision_pending"},
		Target:  []string{workerPoolDesired},
		Refresh: func() (interface{}, string, error) {
			getWorkersOptions := &kubernetesserviceapiv1.GetWorkers1Options{
				Cluster:            &clusterID,
				XAuthResourceGroup: &target.ResourceGroup,
			}

			workers, response, err := satClient.GetWorkers1(getWorkersOptions)
			if err != nil {
				return nil, "", fmt.Errorf("[ERROR] Error retrieving workers for cluster: %s\n%s", err, response)
			}

			//Check active transactions
			//Check for workerpool state to be deployed
			//Done workerpool has two fields desiredState and actualState , so check for those 2
			for _, e := range workers {
				if *e.PoolName == workerPoolID || *e.PoolID == workerPoolID {
					if strings.Compare(*e.Lifecycle.ActualState, workerPoolDesired) != 0 {
						log.Printf("worker: %s state: %s", *e.ID, *e.Lifecycle.ActualState)
						return workers, "provision_pending", nil
					}
				}
			}
			return workers, workerPoolDesired, nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	return stateConf.WaitForState()
}

func WaitForSatelliteWorkerDelete(clusterNameOrID, workerPoolNameOrID string, meta interface{}, timeout time.Duration, target v2.ClusterTargetHeader) (interface{}, error) {
	satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
	if err != nil {
		return nil, err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"deleting"},
		Target:     []string{workerDeleteState},
		Refresh:    satelliteWorkerPoolDeleteStateRefreshFunc(satClient, clusterNameOrID, workerPoolNameOrID, target),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func satelliteWorkerPoolDeleteStateRefreshFunc(satClient *kubernetesserviceapiv1.KubernetesServiceApiV1, clusterID, workerPoolNameOrID string, target v2.ClusterTargetHeader) resource.StateRefreshFunc {

	return func() (interface{}, string, error) {

		getWorkersOptions := &kubernetesserviceapiv1.GetWorkers1Options{
			Cluster:            &clusterID,
			XAuthResourceGroup: &target.ResourceGroup,
		}

		workerFields, response, err := satClient.GetWorkers1(getWorkersOptions)
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] Error retrieving workers for cluster: %s\n%s", err, response)
		}
		//Done worker has two fields desiredState and actualState , so check for those 2
		for _, e := range workerFields {
			if *e.PoolName == workerPoolNameOrID || *e.PoolID == workerPoolNameOrID {
				if strings.Compare(*e.Lifecycle.ActualState, "deleted") != 0 {
					log.Printf("Deleting worker %s", *e.ID)
					return workerFields, "deleting", nil
				}
			}
		}
		return workerFields, workerDeleteState, nil
	}
}
