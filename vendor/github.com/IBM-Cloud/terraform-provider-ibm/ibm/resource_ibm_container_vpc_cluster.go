// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/IBM/go-sdk-core/v3/core"
	"github.com/IBM/vpc-go-sdk/vpcclassicv1"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	v1 "github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
	"github.com/IBM-Cloud/bluemix-go/bmxerror"
)

const (
	deployRequested    = "Deploy requested"
	deployInProgress   = "Deploy in progress"
	ready              = "Ready"
	normal             = "normal"
	masterNodeReady    = "MasterNodeReady"
	oneWorkerNodeReady = "OneWorkerNodeReady"
	ingressReady       = "IngressReady"
)

func resourceIBMContainerVpcCluster() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMContainerVpcClusterCreate,
		Read:     resourceIBMContainerVpcClusterRead,
		Update:   resourceIBMContainerVpcClusterUpdate,
		Delete:   resourceIBMContainerVpcClusterDelete,
		Exists:   resourceIBMContainerVpcClusterExists,
		Importer: &schema.ResourceImporter{},

		CustomizeDiff: customdiff.Sequence(
			func(diff *schema.ResourceDiff, v interface{}) error {
				return resourceTagsCustomizeDiff(diff)
			},
		),

		Schema: map[string]*schema.Schema{

			"flavor": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cluster nodes flavour",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The cluster name",
			},

			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The vpc id where the cluster is",
			},

			"kms_config": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Enables KMS on a given cluster ",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID of the KMS instance to use to encrypt the cluster.",
						},
						"crk_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID of the customer root key.",
						},
						"private_endpoint": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Specify this option to use the KMS public service endpoint.",
						},
					},
				},
			},

			"zones": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Zone info",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Zone for the worker pool in a multizone cluster",
						},

						"subnet_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The VPC subnet to assign the cluster",
						},
					},
				},
			},
			//Optionals in cluster creation

			"kube_version": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				DiffSuppressFunc: func(k, o, n string, d *schema.ResourceData) bool {
					if o == "" {
						return false
					}
					new := strings.Split(n, ".")
					old := strings.Split(o, ".")

					if strings.Compare(new[0]+"."+strings.Split(new[1], "_")[0], old[0]+"."+strings.Split(old[1], "_")[0]) == 0 {
						return true
					}
					return false
				},
				Description: "Kubernetes version",
			},

			"update_all_workers": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Updates all the woker nodes if sets to true",
			},

			"patch_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Kubernetes patch version",
			},

			"retry_patch_version": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Argument which helps to retry the patch version updates on worker nodes. Increment the value to retry the patch updates if the previous apply fails",
			},

			"wait_for_worker_update": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Wait for worker node to update during kube version update.",
			},

			"service_subnet": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Custom subnet CIDR to provide private IP addresses for services",
				Computed:    true,
			},

			"pod_subnet": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Custom subnet CIDR to provide private IP addresses for pods",
				Computed:    true,
			},

			"worker_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Number of worker nodes in the cluster",
			},

			"worker_labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Labels for default worker pool",
			},

			"disable_public_service_endpoint": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Boolean value true if Public service endpoint to be disabled",
			},

			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: InvokeValidator("ibm_container_vpc_cluster", "tag")},
				Set:         resourceIBMVPCHash,
				Description: "List of tags for the resources",
			},

			"wait_till": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          ingressReady,
				DiffSuppressFunc: applyOnce,
				ValidateFunc:     validation.StringInSlice([]string{masterNodeReady, oneWorkerNodeReady, ingressReady}, true),
				Description:      "wait_till can be configured for Master Ready, One worker Ready or Ingress Ready",
			},

			"entitlement": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: applyOnce,
				Description:      "Entitlement option reduces additional OCP Licence cost in Openshift Clusters",
			},

			"cos_instance_crn": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: applyOnce,
				Description:      "A standard cloud object storage instance CRN to back up the internal registry in your OpenShift on VPC Gen 2 cluster",
			},

			"force_delete_storage": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Force the removal of a cluster and its persistent storage. Deleted data cannot be recovered",
			},

			ResourceControllerURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the IBM Cloud dashboard that can be used to explore and view details about this cluster",
			},

			//Get Cluster info Request
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"master_status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"resource_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "ID of the resource group.",
			},

			"master_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"albs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alb_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"load_balancer_hostname": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resize": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"disable_deployment": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"public_service_endpoint_url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"private_service_endpoint_url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CRN of resource instance",
			},

			"ingress_hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ingress_secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			ResourceName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the resource",
			},

			ResourceCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},

			ResourceStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the resource",
			},

			ResourceGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group name in which resource is provisioned",
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(45 * time.Minute),
		},
	}
}

func resourceIBMContainerVpcClusterValidator() *ResourceValidator {
	validateSchema := make([]ValidateSchema, 1)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "tag",
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9:_ .-]+$`,
			MinValueLength:             1,
			MaxValueLength:             128})

	ibmContainerVpcClusteresourceValidator := ResourceValidator{ResourceName: "ibm_container_vpc_cluster", Schema: validateSchema}
	return &ibmContainerVpcClusteresourceValidator
}

func resourceIBMContainerVpcClusterCreate(d *schema.ResourceData, meta interface{}) error {

	var vpcProvider string
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	if userDetails.generation == 1 {
		vpcProvider = "vpc-classic"
	} else {
		vpcProvider = "vpc-gen2"
	}

	csClient, err := meta.(ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}

	disablePublicServiceEndpoint := d.Get("disable_public_service_endpoint").(bool)
	name := d.Get("name").(string)
	var kubeVersion string
	if v, ok := d.GetOk("kube_version"); ok {
		kubeVersion = v.(string)
	}
	podSubnet := d.Get("pod_subnet").(string)
	serviceSubnet := d.Get("service_subnet").(string)
	vpcID := d.Get("vpc_id").(string)
	flavor := d.Get("flavor").(string)
	workerCount := d.Get("worker_count").(int)

	// timeoutStage will define the timeout stage
	var timeoutStage string
	if v, ok := d.GetOk("wait_till"); ok {
		timeoutStage = v.(string)
	}

	var zonesList = make([]v2.Zone, 0)

	if res, ok := d.GetOk("zones"); ok {
		zones := res.(*schema.Set).List()
		for _, e := range zones {
			r, _ := e.(map[string]interface{})
			if ID, subnetID := r["name"], r["subnet_id"]; ID != nil && subnetID != nil {
				zoneParam := v2.Zone{}
				zoneParam.ID, zoneParam.SubnetID = ID.(string), subnetID.(string)
				zonesList = append(zonesList, zoneParam)
			}

		}
	}

	workerpool := v2.WorkerPoolConfig{
		VpcID:       vpcID,
		Flavor:      flavor,
		WorkerCount: workerCount,
		Zones:       zonesList,
	}

	if l, ok := d.GetOk("worker_labels"); ok {
		labels := make(map[string]string)
		for k, v := range l.(map[string]interface{}) {
			labels[k] = v.(string)
		}
		workerpool.Labels = labels
	}

	params := v2.ClusterCreateRequest{
		DisablePublicServiceEndpoint: disablePublicServiceEndpoint,
		Name:                         name,
		KubeVersion:                  kubeVersion,
		PodSubnet:                    podSubnet,
		ServiceSubnet:                serviceSubnet,
		WorkerPools:                  workerpool,
		Provider:                     vpcProvider,
	}

	// Update params with Entitlement option if provided
	if v, ok := d.GetOk("entitlement"); ok {
		params.DefaultWorkerPoolEntitlement = v.(string)
	}

	// Update params with Cloud Object Store instance CRN id option if provided
	if v, ok := d.GetOk("cos_instance_crn"); ok {
		params.CosInstanceCRN = v.(string)
	}

	targetEnv, err := getVpcClusterTargetHeader(d, meta)
	if err != nil {
		return err
	}

	cls, err := csClient.Clusters().Create(params, targetEnv)

	if err != nil {
		return err
	}
	d.SetId(cls.ID)
	switch strings.ToLower(timeoutStage) {

	case strings.ToLower(masterNodeReady):
		_, err = waitForVpcClusterMasterAvailable(d, meta)
		if err != nil {
			return err
		}

	case strings.ToLower(oneWorkerNodeReady):
		_, err = waitForVpcClusterOneWorkerAvailable(d, meta)
		if err != nil {
			return err
		}

	case strings.ToLower(ingressReady):
		_, err = waitForVpcClusterIngressAvailable(d, meta)
		if err != nil {
			return err
		}

	}
	return resourceIBMContainerVpcClusterUpdate(d, meta)

}

func resourceIBMContainerVpcClusterUpdate(d *schema.ResourceData, meta interface{}) error {

	csClient, err := meta.(ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}

	targetEnv, err := getVpcClusterTargetHeader(d, meta)
	if err != nil {
		return err
	}

	clusterID := d.Id()

	v := os.Getenv("IC_ENV_TAGS")
	if d.HasChange("tags") || v != "" {
		oldList, newList := d.GetChange("tags")
		cluster, err := csClient.Clusters().GetCluster(clusterID, targetEnv)
		if err != nil {
			return fmt.Errorf("Error retrieving cluster %s: %s", clusterID, err)
		}
		err = UpdateTagsUsingCRN(oldList, newList, meta, cluster.CRN)
		if err != nil {
			log.Printf(
				"An error occured during update of instance (%s) tags: %s", clusterID, err)
		}
	}

	if d.HasChange("kms_config") {
		kmsConfig := v2.KmsEnableReq{}
		kmsConfig.Cluster = clusterID
		targetEnv := v2.ClusterHeader{}
		if kms, ok := d.GetOk("kms_config"); ok {

			kmsConfiglist := kms.([]interface{})

			for _, l := range kmsConfiglist {
				kmsMap, _ := l.(map[string]interface{})

				//instance_id - Required field
				instanceID := kmsMap["instance_id"].(string)
				kmsConfig.Kms = instanceID

				//crk_id - Required field
				crk := kmsMap["crk_id"].(string)
				kmsConfig.Crk = crk

				//Read event - as its optional check for existence
				if privateEndpoint := kmsMap["private_endpoint"]; privateEndpoint != nil {
					endpoint := privateEndpoint.(bool)
					kmsConfig.PrivateEndpoint = endpoint
				}
			}
		}

		err := csClient.Kms().EnableKms(kmsConfig, targetEnv)
		if err != nil {
			log.Printf(
				"An error occured during EnableKms (cluster: %s) error: %s", d.Id(), err)
			return err
		}

	}

	if (d.HasChange("kube_version") || d.HasChange("update_all_workers") || d.HasChange("patch_version") || d.HasChange("retry_patch_version")) && !d.IsNewResource() {

		if d.HasChange("kube_version") {
			ClusterClient, err := meta.(ClientSession).ContainerAPI()
			if err != nil {
				return err
			}
			var masterVersion string
			if v, ok := d.GetOk("kube_version"); ok {
				masterVersion = v.(string)
			}
			params := v1.ClusterUpdateParam{
				Action:  "update",
				Force:   true,
				Version: masterVersion,
			}

			Env, err := getClusterTargetHeader(d, meta)

			if err != nil {
				return err
			}
			Error := ClusterClient.Clusters().Update(clusterID, params, Env)
			if Error != nil {
				return Error
			}
			_, err = WaitForVpcClusterVersionUpdate(d, meta, targetEnv)
			if err != nil {
				return fmt.Errorf(
					"Error waiting for cluster (%s) version to be updated: %s", d.Id(), err)
			}
		}

		csClient, err := meta.(ClientSession).VpcContainerAPI()
		if err != nil {
			return err
		}
		targetEnv, err := getVpcClusterTargetHeader(d, meta)
		if err != nil {
			return err
		}

		clusterID := d.Id()
		cls, err := csClient.Clusters().GetCluster(clusterID, targetEnv)
		if err != nil {
			return fmt.Errorf("Error retrieving conatiner vpc cluster: %s", err)
		}

		// Update the worker nodes after master node kube-version is updated.
		// workers will store the existing workers info to identify the replaced node
		workersInfo := make(map[string]int, 0)

		updateAllWorkers := d.Get("update_all_workers").(bool)
		if updateAllWorkers || d.HasChange("patch_version") || d.HasChange("retry_patch_version") {

			patchVersion := d.Get("patch_version").(string)
			workers, err := csClient.Workers().ListWorkers(clusterID, false, targetEnv)
			if err != nil {
				d.Set("patch_version", nil)
				return fmt.Errorf("Error retrieving workers for cluster: %s", err)
			}

			for index, worker := range workers {
				workersInfo[worker.ID] = index
			}
			workersCount := len(workers)

			waitForWorkerUpdate := d.Get("wait_for_worker_update").(bool)

			for _, worker := range workers {
				// check if change is present in MAJOR.MINOR version or in PATCH version
				if strings.Split(worker.KubeVersion.Actual, "_")[0] != strings.Split(cls.MasterKubeVersion, "_")[0] || (strings.Split(worker.KubeVersion.Actual, ".")[2] != patchVersion && patchVersion == strings.Split(worker.KubeVersion.Target, ".")[2]) {
					_, err := csClient.Workers().ReplaceWokerNode(clusterID, worker.ID, targetEnv)
					// As API returns http response 204 NO CONTENT, error raised will be exempted.
					if err != nil && !strings.Contains(err.Error(), "EmptyResponseBody") {
						d.Set("patch_version", nil)
						return fmt.Errorf("Error replacing the worker node from the cluster: %s", err)
					}

					if waitForWorkerUpdate {
						//1. wait for worker node to delete
						_, deleteError := waitForWorkerNodetoDelete(d, meta, targetEnv, worker.ID)
						if deleteError != nil {
							d.Set("patch_version", nil)
							return fmt.Errorf("Worker node - %s is failed to replace", worker.ID)
						}

						//2. wait for new workerNode
						_, newWorkerError := waitForNewWorker(d, meta, targetEnv, workersCount)
						if newWorkerError != nil {
							d.Set("patch_version", nil)
							return fmt.Errorf("Failed to spawn new worker node")
						}

						//3. Get new worker node ID and update the map
						newWorkerID, index, newNodeError := getNewWorkerID(d, meta, targetEnv, workersInfo)
						if newNodeError != nil {
							d.Set("patch_version", nil)
							return fmt.Errorf("Unable to find the new worker node info")
						}

						delete(workersInfo, worker.ID)
						workersInfo[newWorkerID] = index

						//4. wait for the worker's version update and normal state
						_, Err := WaitForVpcClusterWokersVersionUpdate(d, meta, targetEnv, cls.MasterKubeVersion, newWorkerID)
						if Err != nil {
							d.Set("patch_version", nil)
							return fmt.Errorf(
								"Error waiting for cluster (%s) worker nodes kube version to be updated: %s", d.Id(), Err)
						}
					}
				}
			}
		}
	}

	if d.HasChange("worker_labels") && !d.IsNewResource() {
		labels := make(map[string]string)
		if l, ok := d.GetOk("worker_labels"); ok {
			for k, v := range l.(map[string]interface{}) {
				labels[k] = v.(string)
			}
		}

		ClusterClient, err := meta.(ClientSession).ContainerAPI()
		if err != nil {
			return err
		}
		Env := v1.ClusterTargetHeader{ResourceGroup: targetEnv.ResourceGroup}

		err = ClusterClient.WorkerPools().UpdateLabelsWorkerPool(clusterID, "default", labels, Env)
		if err != nil {
			return fmt.Errorf(
				"Error updating the labels: %s", err)
		}
	}

	if d.HasChange("worker_count") && !d.IsNewResource() {
		count := d.Get("worker_count").(int)
		ClusterClient, err := meta.(ClientSession).ContainerAPI()
		if err != nil {
			return err
		}
		Env := v1.ClusterTargetHeader{ResourceGroup: targetEnv.ResourceGroup}

		err = ClusterClient.WorkerPools().ResizeWorkerPool(clusterID, "default", count, Env)
		if err != nil {
			return fmt.Errorf(
				"Error updating the worker_count %d: %s", count, err)
		}
	}
	if d.HasChange("zones") && !d.IsNewResource() {
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
				zoneParam := v2.WorkerPoolZone{
					Cluster:      clusterID,
					Id:           newZone["name"].(string),
					SubnetID:     newZone["subnet_id"].(string),
					WorkerPoolID: "default",
				}
				err = csClient.WorkerPools().CreateWorkerPoolZone(zoneParam, targetEnv)
				if err != nil {
					return fmt.Errorf("Error adding zone to conatiner vpc cluster: %s", err)
				}
				_, err = WaitForWorkerPoolAvailable(d, meta, clusterID, "default", d.Timeout(schema.TimeoutCreate), targetEnv)
				if err != nil {
					return fmt.Errorf(
						"Error waiting for workerpool (%s) to become ready: %s", d.Id(), err)
				}

			}
		}
		if len(remove) > 0 {
			for _, zone := range remove {
				oldZone := zone.(map[string]interface{})
				ClusterClient, err := meta.(ClientSession).ContainerAPI()
				if err != nil {
					return err
				}
				Env := v1.ClusterTargetHeader{ResourceGroup: targetEnv.ResourceGroup}
				err = ClusterClient.WorkerPools().RemoveZone(clusterID, oldZone["name"].(string), "default", Env)
				if err != nil {
					return fmt.Errorf("Error deleting zone to conatiner vpc cluster: %s", err)
				}
				_, err = WaitForV2WorkerZoneDeleted(clusterID, "default", oldZone["name"].(string), meta, d.Timeout(schema.TimeoutDelete), targetEnv)
				if err != nil {
					return fmt.Errorf(
						"Error waiting for deleting workers of worker pool (%s) of cluster (%s):  %s", "default", clusterID, err)
				}
			}
		}
	}

	if d.HasChange("force_delete_storage") {
		var forceDeleteStorage bool
		if v, ok := d.GetOk("force_delete_storage"); ok {
			forceDeleteStorage = v.(bool)
		}
		d.Set("force_delete_storage", forceDeleteStorage)
	}

	return resourceIBMContainerVpcClusterRead(d, meta)
}
func WaitForV2WorkerZoneDeleted(clusterNameOrID, workerPoolNameOrID, zone string, meta interface{}, timeout time.Duration, target v2.ClusterTargetHeader) (interface{}, error) {
	csClient, err := meta.(ClientSession).VpcContainerAPI()
	if err != nil {
		return nil, err
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"deleting"},
		Target:     []string{workerDeleteState},
		Refresh:    workerPoolV2ZoneDeleteStateRefreshFunc(csClient.Workers(), clusterNameOrID, workerPoolNameOrID, zone, target),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}
func workerPoolV2ZoneDeleteStateRefreshFunc(client v2.Workers, instanceID, workerPoolNameOrID, zone string, target v2.ClusterTargetHeader) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		workerFields, err := client.ListByWorkerPool(instanceID, workerPoolNameOrID, true, target)
		if err != nil {
			return nil, "", fmt.Errorf("Error retrieving workers for cluster: %s", err)
		}
		//Done worker has two fields State and Status , so check for those 2
		for _, e := range workerFields {
			if e.Location == zone {
				if strings.Compare(e.LifeCycle.ActualState, "deleted") != 0 {
					return workerFields, "deleting", nil
				}
			}
		}
		return workerFields, workerDeleteState, nil
	}
}
func resourceIBMContainerVpcClusterRead(d *schema.ResourceData, meta interface{}) error {

	csClient, err := meta.(ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}
	albsAPI := csClient.Albs()

	targetEnv, err := getVpcClusterTargetHeader(d, meta)
	if err != nil {
		return err
	}

	clusterID := d.Id()
	cls, err := csClient.Clusters().GetCluster(clusterID, targetEnv)
	if err != nil {
		return fmt.Errorf("Error retrieving conatiner vpc cluster: %s", err)
	}

	workerPool, err := csClient.WorkerPools().GetWorkerPool(clusterID, "default", targetEnv)

	var zones = make([]map[string]interface{}, 0)
	for _, zone := range workerPool.Zones {
		for _, subnet := range zone.Subnets {
			if subnet.Primary == true {
				zoneInfo := map[string]interface{}{
					"name":      zone.ID,
					"subnet_id": subnet.ID,
				}
				zones = append(zones, zoneInfo)
			}
		}
	}

	albs, err := albsAPI.ListClusterAlbs(clusterID, targetEnv)
	if err != nil && !strings.Contains(err.Error(), "This operation is not supported for your cluster's version.") {
		return fmt.Errorf("Error retrieving alb's of the cluster %s: %s", clusterID, err)
	}

	d.Set("name", cls.Name)
	d.Set("crn", cls.CRN)
	d.Set("master_status", cls.Lifecycle.MasterStatus)
	d.Set("zones", zones)
	if strings.HasSuffix(cls.MasterKubeVersion, "_openshift") {
		d.Set("kube_version", strings.Split(cls.MasterKubeVersion, "_")[0]+"_openshift")
	} else {
		d.Set("kube_version", strings.Split(cls.MasterKubeVersion, "_")[0])
	}
	d.Set("worker_count", workerPool.WorkerCount)
	d.Set("worker_labels", IgnoreSystemLabels(workerPool.Labels))
	if cls.Vpcs != nil {
		d.Set("vpc_id", cls.Vpcs[0])
	}
	d.Set("master_url", cls.MasterURL)
	d.Set("flavor", workerPool.Flavor)
	d.Set("service_subnet", cls.ServiceSubnet)
	d.Set("pod_subnet", cls.PodSubnet)
	d.Set("state", cls.State)
	d.Set("ingress_hostname", cls.Ingress.HostName)
	d.Set("ingress_secret", cls.Ingress.SecretName)
	d.Set("albs", flattenVpcAlbs(albs, "all"))
	d.Set("resource_group_id", cls.ResourceGroupID)
	d.Set("public_service_endpoint_url", cls.ServiceEndpoints.PublicServiceEndpointURL)
	d.Set("private_service_endpoint_url", cls.ServiceEndpoints.PrivateServiceEndpointURL)
	if cls.ServiceEndpoints.PublicServiceEndpointURL != "" {
		d.Set("disable_public_service_endpoint", false)
	} else {
		d.Set("disable_public_service_endpoint", true)
	}

	tags, err := GetTagsUsingCRN(meta, cls.CRN)
	if err != nil {
		log.Printf(
			"An error occured during reading of instance (%s) tags : %s", d.Id(), err)
	}
	d.Set("tags", tags)
	controller, err := getBaseController(meta)
	if err != nil {
		return err
	}
	d.Set(ResourceControllerURL, controller+"/kubernetes/clusters")
	d.Set(ResourceName, cls.Name)
	d.Set(ResourceCRN, cls.CRN)
	d.Set(ResourceStatus, cls.State)
	d.Set(ResourceGroupName, cls.ResourceGroupName)

	return nil
}

func resourceIBMContainerVpcClusterDelete(d *schema.ResourceData, meta interface{}) error {

	targetEnv, err := getVpcClusterTargetHeader(d, meta)
	if err != nil {
		return err
	}
	csClient, err := meta.(ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}
	clusterID := d.Id()

	var zonesList = make([]v2.Zone, 0)

	if res, ok := d.GetOk("zones"); ok {
		zones := res.(*schema.Set).List()
		for _, e := range zones {
			r, _ := e.(map[string]interface{})
			if ID, subnetID := r["name"], r["subnet_id"]; ID != nil && subnetID != nil {
				zoneParam := v2.Zone{}
				zoneParam.ID, zoneParam.SubnetID = ID.(string), subnetID.(string)
				zonesList = append(zonesList, zoneParam)
			}

		}
	}
	var region = ""
	if len(zonesList) > 0 {
		splitZone := strings.Split(zonesList[0].ID, "-")
		region = splitZone[0] + "-" + splitZone[1]
	}

	bxsession, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return err
	}

	var authenticator *core.BearerTokenAuthenticator
	if strings.HasPrefix(bxsession.Config.IAMAccessToken, "Bearer") {
		authenticator = &core.BearerTokenAuthenticator{
			BearerToken: bxsession.Config.IAMAccessToken[7:],
		}
	} else {
		authenticator = &core.BearerTokenAuthenticator{
			BearerToken: bxsession.Config.IAMAccessToken,
		}
	}

	forceDeleteStorage := d.Get("force_delete_storage").(bool)
	err = csClient.Clusters().Delete(clusterID, targetEnv, forceDeleteStorage)
	if err != nil {
		return fmt.Errorf("Error deleting cluster: %s", err)
	}
	_, err = waitForVpcClusterDelete(d, meta)
	if err != nil {
		return err
	}

	if region != "" {
		userDetails, err := meta.(ClientSession).BluemixUserDetails()
		if err != nil {
			return err
		}
		if userDetails.generation == 1 {
			vpcclassicurl := fmt.Sprintf("https://%s.iaas.cloud.ibm.com/v1", region)
			vpcclassicoptions := &vpcclassicv1.VpcClassicV1Options{
				URL:           envFallBack([]string{"IBMCLOUD_IS_API_ENDPOINT"}, vpcclassicurl),
				Authenticator: authenticator,
			}
			sess1, err := vpcclassicv1.NewVpcClassicV1(vpcclassicoptions)
			if err != nil {
				log.Println("error creating vpcclassic session", err)
			}
			listlbOptions := &vpcclassicv1.ListLoadBalancersOptions{}
			lbs, response, err1 := sess1.ListLoadBalancers(listlbOptions)
			if err1 != nil {
				log.Printf("Error Retrieving vpc load balancers: %s\n%s", err, response)
			}
			if lbs != nil && lbs.LoadBalancers != nil && len(lbs.LoadBalancers) > 0 {
				for _, lb := range lbs.LoadBalancers {
					if strings.Contains(*lb.Name, clusterID) {
						log.Println("Deleting Load Balancer", *lb.Name)
						id := *lb.ID
						_, err = isWaitForClassicLBDeleted(sess1, id, d.Timeout(schema.TimeoutDelete))
						if err != nil {
							log.Printf("Error waiting for vpc load balancer to be deleted: %s\n", err)

						}
					}
				}
			}
		} else {
			vpcurl := fmt.Sprintf("https://%s.iaas.cloud.ibm.com/v1", region)
			vpcoptions := &vpcv1.VpcV1Options{
				URL:           envFallBack([]string{"IBMCLOUD_IS_NG_API_ENDPOINT"}, vpcurl),
				Authenticator: authenticator,
			}
			sess1, err := vpcv1.NewVpcV1(vpcoptions)
			if err != nil {
				log.Println("error creating vpc session", err)
			}
			listlbOptions := &vpcv1.ListLoadBalancersOptions{}
			lbs, response, err1 := sess1.ListLoadBalancers(listlbOptions)
			if err1 != nil {
				log.Printf("Error Retrieving vpc load balancers: %s\n%s", err, response)
			}
			if lbs != nil && lbs.LoadBalancers != nil && len(lbs.LoadBalancers) > 0 {
				for _, lb := range lbs.LoadBalancers {
					if strings.Contains(*lb.Name, clusterID) {
						log.Println("Deleting Load Balancer", *lb.Name)
						id := *lb.ID
						_, err = isWaitForLBDeleted(sess1, id, d.Timeout(schema.TimeoutDelete))
						if err != nil {
							log.Printf("Error waiting for vpc load balancer to be deleted: %s\n", err)

						}
					}
				}
			}
		}
	}
	return nil
}

func waitForVpcClusterDelete(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	targetEnv, err := getVpcClusterTargetHeader(d, meta)
	if err != nil {
		return nil, err
	}
	csClient, err := meta.(ClientSession).VpcContainerAPI()
	if err != nil {
		return nil, err
	}
	clusterID := d.Id()
	deleteStateConf := &resource.StateChangeConf{
		Pending: []string{clusterDeletePending},
		Target:  []string{clusterDeleted},
		Refresh: func() (interface{}, string, error) {
			cluster, err := csClient.Clusters().GetCluster(clusterID, targetEnv)
			if err != nil {
				if apiErr, ok := err.(bmxerror.RequestFailure); ok && (apiErr.StatusCode() == 404) {
					return cluster, clusterDeleted, nil
				}
				return nil, "", err
			}
			return cluster, clusterDeletePending, nil
		},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		MinTimeout:   5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	return deleteStateConf.WaitForState()
}

func waitForVpcClusterOneWorkerAvailable(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	targetEnv, err := getVpcClusterTargetHeader(d, meta)
	if err != nil {
		return nil, err
	}
	csClient, err := meta.(ClientSession).VpcContainerAPI()
	if err != nil {
		return nil, err
	}
	clusterID := d.Id()
	createStateConf := &resource.StateChangeConf{
		Pending: []string{deployRequested, deployInProgress},
		Target:  []string{normal},
		Refresh: func() (interface{}, string, error) {
			workers, err := csClient.Workers().ListByWorkerPool(clusterID, "default", false, targetEnv)
			if err != nil {
				return workers, deployInProgress, err
			}
			if len(workers) == 0 {
				return workers, deployInProgress, nil
			}

			for _, worker := range workers {
				log.Println("worker: ", worker.ID)
				log.Println("worker health state:  ", worker.Health.State)

				if worker.Health.State == normal {
					return workers, normal, nil
				}
			}
			return workers, deployInProgress, nil

		},
		Timeout:                   d.Timeout(schema.TimeoutCreate),
		Delay:                     10 * time.Second,
		MinTimeout:                5 * time.Second,
		ContinuousTargetOccurence: 5,
	}
	return createStateConf.WaitForState()
}

func waitForVpcClusterMasterAvailable(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	targetEnv, err := getVpcClusterTargetHeader(d, meta)
	if err != nil {
		return nil, err
	}
	csClient, err := meta.(ClientSession).VpcContainerAPI()
	if err != nil {
		return nil, err
	}
	clusterID := d.Id()
	createStateConf := &resource.StateChangeConf{
		Pending: []string{deployRequested, deployInProgress},
		Target:  []string{ready},
		Refresh: func() (interface{}, string, error) {
			clusterInfo, clusterInfoErr := csClient.Clusters().GetCluster(clusterID, targetEnv)

			if err != nil || clusterInfoErr != nil {
				return clusterInfo, deployInProgress, err
			}

			if clusterInfo.Lifecycle.MasterStatus == ready {
				return clusterInfo, ready, nil
			}
			return clusterInfo, deployInProgress, nil

		},
		Timeout:                   d.Timeout(schema.TimeoutCreate),
		Delay:                     10 * time.Second,
		MinTimeout:                5 * time.Second,
		ContinuousTargetOccurence: 5,
	}
	return createStateConf.WaitForState()
}

func waitForVpcClusterIngressAvailable(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	targetEnv, err := getVpcClusterTargetHeader(d, meta)
	if err != nil {
		return nil, err
	}
	csClient, err := meta.(ClientSession).VpcContainerAPI()
	if err != nil {
		return nil, err
	}
	clusterID := d.Id()
	createStateConf := &resource.StateChangeConf{
		Pending: []string{deployRequested, deployInProgress},
		Target:  []string{ready},
		Refresh: func() (interface{}, string, error) {
			clusterInfo, clusterInfoErr := csClient.Clusters().GetCluster(clusterID, targetEnv)

			if err != nil || clusterInfoErr != nil {
				return clusterInfo, deployInProgress, err
			}

			if clusterInfo.Ingress.HostName != "" {
				return clusterInfo, ready, nil
			}
			return clusterInfo, deployInProgress, nil

		},
		Timeout:                   d.Timeout(schema.TimeoutCreate),
		Delay:                     10 * time.Second,
		MinTimeout:                5 * time.Second,
		ContinuousTargetOccurence: 5,
	}
	return createStateConf.WaitForState()
}

func getVpcClusterTargetHeader(d *schema.ResourceData, meta interface{}) (v2.ClusterTargetHeader, error) {
	targetEnv := v2.ClusterTargetHeader{}
	var resourceGroup string
	if rg, ok := d.GetOk("resource_group_id"); ok {
		resourceGroup = rg.(string)
		targetEnv.ResourceGroup = resourceGroup
	}

	return targetEnv, nil
}

func resourceIBMContainerVpcClusterExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	csClient, err := meta.(ClientSession).VpcContainerAPI()
	if err != nil {
		return false, err
	}
	targetEnv, err := getVpcClusterTargetHeader(d, meta)
	if err != nil {
		return false, err
	}
	clusterID := d.Id()
	cls, err := csClient.Clusters().GetCluster(clusterID, targetEnv)
	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 && strings.Contains(apiErr.Description(), "The specified cluster could not be found") {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error communicating with the API: %s", err)
	}
	return cls.ID == clusterID, nil
}

// WaitForVpcClusterVersionUpdate Waits for cluster creation
func WaitForVpcClusterVersionUpdate(d *schema.ResourceData, meta interface{}, target v2.ClusterTargetHeader) (interface{}, error) {
	csClient, err := meta.(ClientSession).VpcContainerAPI()
	if err != nil {
		return nil, err
	}
	log.Printf("Waiting for cluster (%s) version to be updated.", d.Id())
	id := d.Id()

	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"retry", versionUpdating},
		Target:                    []string{clusterNormal},
		Refresh:                   vpcClusterVersionRefreshFunc(csClient.Clusters(), id, d, target),
		Timeout:                   d.Timeout(schema.TimeoutUpdate),
		Delay:                     10 * time.Second,
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 5,
	}

	return stateConf.WaitForState()
}

func vpcClusterVersionRefreshFunc(client v2.Clusters, instanceID string, d *schema.ResourceData, target v2.ClusterTargetHeader) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		cls, err := client.GetCluster(instanceID, target)
		if err != nil {
			return nil, "retry", fmt.Errorf("Error retrieving conatiner vpc cluster: %s", err)
		}

		// Check active transactions
		log.Println("Checking cluster version", cls.MasterKubeVersion, d.Get("kube_version").(string))
		if strings.Contains(cls.MasterKubeVersion, "(pending)") {
			return cls, versionUpdating, nil
		}
		return cls, clusterNormal, nil
	}
}

// WaitForVpcClusterWokersVersionUpdate Waits for Cluster version Update
func WaitForVpcClusterWokersVersionUpdate(d *schema.ResourceData, meta interface{}, target v2.ClusterTargetHeader, masterVersion, workerID string) (interface{}, error) {
	csClient, err := meta.(ClientSession).VpcContainerAPI()
	if err != nil {
		return nil, err
	}

	log.Printf("Waiting for worker (%s) version to be updated.", workerID)
	clusterID := d.Id()
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"retry", versionUpdating},
		Target:                    []string{workerNormal},
		Refresh:                   vpcClusterWorkersVersionRefreshFunc(csClient.Workers(), workerID, clusterID, d, target, masterVersion),
		Timeout:                   d.Timeout(schema.TimeoutUpdate),
		Delay:                     10 * time.Second,
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 5,
	}

	return stateConf.WaitForState()
}

func vpcClusterWorkersVersionRefreshFunc(client v2.Workers, workerID, clusterID string, d *schema.ResourceData, target v2.ClusterTargetHeader, masterVersion string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		worker, err := client.Get(clusterID, workerID, target)
		if err != nil {
			return nil, "retry", fmt.Errorf("Error retrieving worker of container vpc cluster: %s", err)
		}
		// Check active updates
		if worker.Health.State == "normal" && strings.Split(worker.KubeVersion.Actual, "_")[0] == strings.Split(masterVersion, "_")[0] {
			return worker, workerNormal, nil
		}
		return worker, versionUpdating, nil
	}
}

func waitForWorkerNodetoDelete(d *schema.ResourceData, meta interface{}, targetEnv v2.ClusterTargetHeader, workerID string) (interface{}, error) {

	csClient, err := meta.(ClientSession).VpcContainerAPI()
	if err != nil {
		return nil, err
	}

	clusterID := d.Id()
	deleteStateConf := &resource.StateChangeConf{
		Pending: []string{workerDeletePending},
		Target:  []string{workerDeleteState},
		Refresh: func() (interface{}, string, error) {
			worker, err := csClient.Workers().Get(clusterID, workerID, targetEnv)
			if err != nil {
				return worker, workerDeletePending, nil
			}
			if worker.LifeCycle.ActualState == "deleted" {
				return worker, workerDeleteState, nil
			}
			return worker, workerDeletePending, nil
		},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		MinTimeout:   5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	return deleteStateConf.WaitForState()
}

func waitForNewWorker(d *schema.ResourceData, meta interface{}, targetEnv v2.ClusterTargetHeader, workersCount int) (interface{}, error) {
	csClient, err := meta.(ClientSession).VpcContainerAPI()
	if err != nil {
		return nil, err
	}

	clusterID := d.Id()
	stateConf := &resource.StateChangeConf{
		Pending: []string{"creating"},
		Target:  []string{"created"},
		Refresh: func() (interface{}, string, error) {
			workers, err := csClient.Workers().ListWorkers(clusterID, false, targetEnv)
			if err != nil {
				return workers, "", fmt.Errorf("Error in retriving the list of worker nodes")
			}
			if len(workers) == workersCount {
				return workers, "created", nil
			}
			return workers, "creating", nil
		},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		MinTimeout:   5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	return stateConf.WaitForState()
}

func getNewWorkerID(d *schema.ResourceData, meta interface{}, targetEnv v2.ClusterTargetHeader, workersInfo map[string]int) (string, int, error) {
	csClient, err := meta.(ClientSession).VpcContainerAPI()
	if err != nil {
		return "", -1, err
	}

	clusterID := d.Id()

	workers, err := csClient.Workers().ListWorkers(clusterID, false, targetEnv)
	if err != nil {
		return "", -1, fmt.Errorf("Error in retriving the list of worker nodes")
	}

	for index, worker := range workers {
		if _, ok := workersInfo[worker.ID]; !ok {
			log.Println("found new replaced node: ", worker.ID)
			return worker.ID, index, nil
		}
	}
	return "", -1, fmt.Errorf("no new node found")
}
