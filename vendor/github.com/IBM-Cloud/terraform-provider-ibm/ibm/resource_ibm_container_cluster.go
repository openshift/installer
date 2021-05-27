// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	v1 "github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
	"github.com/IBM-Cloud/bluemix-go/bmxerror"
)

const (
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

	hardwareShared    = "shared"
	hardwareDedicated = "dedicated"
	isolationPublic   = "public"
	isolationPrivate  = "private"

	defaultWorkerPool = "default"
	computeWorkerPool = "compute"
	gatewayWorkerpool = "gateway"
)

const PUBLIC_SUBNET_TYPE = "public"

func resourceIBMContainerCluster() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMContainerClusterCreate,
		Read:     resourceIBMContainerClusterRead,
		Update:   resourceIBMContainerClusterUpdate,
		Delete:   resourceIBMContainerClusterDelete,
		Exists:   resourceIBMContainerClusterExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Delete: schema.DefaultTimeout(45 * time.Minute),
			Update: schema.DefaultTimeout(45 * time.Minute),
		},

		CustomizeDiff: customdiff.Sequence(
			func(diff *schema.ResourceDiff, v interface{}) error {
				return resourceTagsCustomizeDiff(diff)
			},
		),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The cluster name",
			},
			"datacenter": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The datacenter where this cluster will be deployed",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Deprecated:  "This field is deprecated",
				Computed:    true,
				Description: "The cluster region",
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

			"worker_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				Description:  "Number of worker nodes",
				ValidateFunc: validateWorkerNum,
				Deprecated:   "This field is deprecated",
			},

			"default_pool_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				Description:  "The size of the default worker pool",
				ValidateFunc: validateWorkerNum,
			},

			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "list of labels to the default worker pool",
			},

			"workers_info": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"pool_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Description: "The IDs of the worker node",
			},

			"disk_encryption": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     true,
				Description: "disc encryption done, if set to true.",
			},

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
				Description: "Kubernetes version info",
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

			"update_all_workers": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Updates all the woker nodes if sets to true",
			},

			"machine_type": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "Machine type",
			},

			"hardware": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{hardwareShared, hardwareDedicated}),
				Description:  "Hardware type",
			},

			"billing": {
				Type:             schema.TypeString,
				Optional:         true,
				Deprecated:       "This field is deprecated",
				DiffSuppressFunc: applyOnce,
			},
			"public_vlan_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  nil,
				DiffSuppressFunc: func(k, o, n string, d *schema.ResourceData) bool {
					if o == "" {
						return false
					}
					if o != "" && n == "" {
						return true
					}
					return false
				},
				Description: "Public VLAN ID",
			},

			"private_vlan_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  nil,
				DiffSuppressFunc: func(k, o, n string, d *schema.ResourceData) bool {
					if o == "" {
						return false
					}
					if o != "" && n == "" {
						return true
					}
					return false
				},
				Description: "Private VLAN ID",
			},
			"entitlement": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: applyOnce,
				Description:      "Entitlement option reduces additional OCP Licence cost in Openshift Clusters",
			},

			"wait_for_worker_update": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Wait for worker node to update during kube version update.",
			},
			"wait_till": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          ingressReady,
				DiffSuppressFunc: applyOnce,
				ValidateFunc:     validation.StringInSlice([]string{masterNodeReady, oneWorkerNodeReady, ingressReady}, true),
				Description:      "wait_till can be configured for Master Ready, One worker Ready or Ingress Ready",
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

			"ingress_hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ingress_secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"no_subnet": {
				Type:             schema.TypeBool,
				Optional:         true,
				ForceNew:         true,
				Default:          false,
				DiffSuppressFunc: applyOnce,
				Description:      "Boolean value set to true when subnet creation is not required.",
			},
			"is_trusted": {
				Type:             schema.TypeBool,
				Optional:         true,
				Deprecated:       "This field is deprecated",
				DiffSuppressFunc: applyOnce,
			},
			"server_url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"subnet_id": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "List of subnet IDs",
			},
			"webhook": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"level": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateAllowedStringValue([]string{"slack"}),
						},
						"url": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"force_delete_storage": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Force the removal of a cluster and its persistent storage. Deleted data cannot be recovered",
			},

			"resource_group_id": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "ID of the resource group.",
				Computed:         true,
				DiffSuppressFunc: applyOnce,
			},

			"org_guid": {
				Description: "The bluemix organization guid this cluster belongs to",
				Type:        schema.TypeString,
				Optional:    true,
				Deprecated:  "This field is deprecated",
			},
			"space_guid": {
				Description: "The bluemix space guid this cluster belongs to",
				Type:        schema.TypeString,
				Optional:    true,
				Deprecated:  "This field is deprecated",
			},
			"account_guid": {
				Description: "The bluemix account guid this cluster belongs to",
				Type:        schema.TypeString,
				Optional:    true,
				Deprecated:  "This field is deprecated",
			},
			"wait_time_minutes": {
				Type:       schema.TypeInt,
				Optional:   true,
				Deprecated: "This field is deprecated",
			},
			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: InvokeValidator("ibm_container_cluster", "tag")},
				Set:         resourceIBMVPCHash,
				Description: "Tags for the resource",
			},

			"worker_pools": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"machine_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size_per_zone": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"hardware": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"labels": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"zones": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zone": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"private_vlan": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"public_vlan": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"worker_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
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
						"num_of_instances": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alb_ip": {
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
			"public_service_endpoint": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"private_service_endpoint": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"public_service_endpoint_url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"private_service_endpoint_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gateway_enabled": {
				Type:             schema.TypeBool,
				Optional:         true,
				DiffSuppressFunc: applyOnce,
				Default:          false,
				Description:      "Set true for gateway enabled clusters",
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CRN of resource instance",
			},

			ResourceControllerURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the IBM Cloud dashboard that can be used to explore and view details about this cluster",
			},

			ResourceName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the resource",
			},

			ResourceCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},

			ResourceStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the resource",
			},

			ResourceGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group name in which resource is provisioned",
			},
		},
	}
}

func resourceIBMContainerClusterValidator() *ResourceValidator {
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

	ibmContainerClusterResourceValidator := ResourceValidator{ResourceName: "ibm_container_cluster", Schema: validateSchema}
	return &ibmContainerClusterResourceValidator
}

func resourceIBMContainerClusterCreate(d *schema.ResourceData, meta interface{}) error {

	csClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	datacenter := d.Get("datacenter").(string)
	machineType := d.Get("machine_type").(string)
	publicVlanID := d.Get("public_vlan_id").(string)
	privateVlanID := d.Get("private_vlan_id").(string)
	noSubnet := d.Get("no_subnet").(bool)
	diskEncryption := d.Get("disk_encryption").(bool)
	defaultPoolSize := d.Get("default_pool_size").(int)
	gatewayEnabled := d.Get("gateway_enabled").(bool)
	hardware := d.Get("hardware").(string)
	switch strings.ToLower(hardware) {
	case hardwareDedicated:
		hardware = isolationPrivate
	case hardwareShared:
		hardware = isolationPublic
	}

	params := v1.ClusterCreateRequest{
		Name:           name,
		Datacenter:     datacenter,
		WorkerNum:      defaultPoolSize,
		MachineType:    machineType,
		PublicVlan:     publicVlanID,
		PrivateVlan:    privateVlanID,
		NoSubnet:       noSubnet,
		Isolation:      hardware,
		DiskEncryption: diskEncryption,
	}

	// Update params with Entitlement option if provided
	if v, ok := d.GetOk("entitlement"); ok {
		params.DefaultWorkerPoolEntitlement = v.(string)
	}
	if v, ok := d.GetOk("pod_subnet"); ok {
		params.PodSubnet = v.(string)
	}
	if v, ok := d.GetOk("service_subnet"); ok {
		params.ServiceSubnet = v.(string)
	}

	if gatewayEnabled {
		if v, ok := d.GetOkExists("private_service_endpoint"); ok {
			if v.(bool) {
				params.PrivateEndpointEnabled = v.(bool)
				params.GatewayEnabled = gatewayEnabled
			} else {
				return fmt.Errorf("set private_service_endpoint to true for gateway_enabled clusters")
			}
		} else {
			return fmt.Errorf("set private_service_endpoint to true for gateway_enabled clusters")
		}
	}
	if v, ok := d.GetOk("kube_version"); ok {
		params.MasterVersion = v.(string)
	}
	if v, ok := d.GetOkExists("private_service_endpoint"); ok {
		params.PrivateEndpointEnabled = v.(bool)
	}
	if v, ok := d.GetOkExists("public_service_endpoint"); ok {
		params.PublicEndpointEnabled = v.(bool)
	}

	targetEnv, err := getClusterTargetHeader(d, meta)
	if err != nil {
		return err
	}

	cls, err := csClient.Clusters().Create(params, targetEnv)
	if err != nil {
		return err
	}
	d.SetId(cls.ID)

	_, err = waitForClusterMasterAvailable(d, meta)
	if err != nil {
		return err
	}
	if d.Get("wait_till").(string) == oneWorkerNodeReady {
		_, err = waitForClusterOneWorkerAvailable(d, meta)
		if err != nil {
			return err
		}
	}
	d.Set("force_delete_storage", d.Get("force_delete_storage").(bool))

	return resourceIBMContainerClusterUpdate(d, meta)
}

func resourceIBMContainerClusterRead(d *schema.ResourceData, meta interface{}) error {
	csClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return err
	}
	wrkAPI := csClient.Workers()
	workerPoolsAPI := csClient.WorkerPools()
	albsAPI := csClient.Albs()

	targetEnv, err := getClusterTargetHeader(d, meta)
	if err != nil {
		return err
	}

	clusterID := d.Id()
	cls, err := csClient.Clusters().Find(clusterID, targetEnv)
	if err != nil {
		return fmt.Errorf("Error retrieving armada cluster: %s", err)
	}

	workerFields, err := wrkAPI.List(clusterID, targetEnv)
	if err != nil {
		return fmt.Errorf("Error retrieving workers for cluster: %s", err)
	}
	workerCount := 0
	workers := []map[string]string{}
	for _, w := range workerFields {
		var worker = map[string]string{
			"id":        w.ID,
			"version":   strings.Split(w.KubeVersion, "_")[0],
			"pool_name": w.PoolName,
		}
		workers = append(workers, worker)
		if w.PoolID == "" && w.PoolName == "" {
			workerCount = workerCount + 1
		}
	}

	d.Set("worker_num", workerCount)

	workerPools, err := workerPoolsAPI.ListWorkerPools(clusterID, targetEnv)
	if err != nil {
		return err
	}
	var poolName string
	var poolContains bool

	if len(workerPools) > 0 && workerPoolContains(workerPools, defaultWorkerPool) {
		poolName = defaultWorkerPool
		poolContains = true
	} else if len(workerPools) > 0 && workerPoolContains(workerPools, computeWorkerPool) && workerPoolContains(workerPools, gatewayWorkerpool) {
		poolName = computeWorkerPool
		poolContains = true
	}
	if poolContains {
		workersByPool, err := wrkAPI.ListByWorkerPool(clusterID, poolName, false, targetEnv)
		if err != nil {
			return fmt.Errorf("Error retrieving workers of default worker pool for cluster: %s", err)
		}

		// to get the private and public vlan IDs of the gateway enabled cluster.
		if poolName == computeWorkerPool {
			gatewayWorkersByPool, err := wrkAPI.ListByWorkerPool(clusterID, gatewayWorkerpool, false, targetEnv)
			if err != nil {
				return fmt.Errorf("Error retrieving workers of default worker pool for cluster: %s", err)
			}
			d.Set("public_vlan_id", gatewayWorkersByPool[0].PublicVlan)
			d.Set("private_vlan_id", gatewayWorkersByPool[0].PrivateVlan)
		} else {
			d.Set("public_vlan_id", workersByPool[0].PublicVlan)
			d.Set("private_vlan_id", workersByPool[0].PrivateVlan)
		}
		d.Set("machine_type", strings.Split(workersByPool[0].MachineType, ".encrypted")[0])
		d.Set("datacenter", cls.DataCenter)
		if workersByPool[0].MachineType != "free" {
			if strings.HasSuffix(workersByPool[0].MachineType, ".encrypted") {
				d.Set("disk_encryption", true)
			} else {
				d.Set("disk_encryption", false)
			}
		}

		if len(workersByPool) > 0 {
			hardware := workersByPool[0].Isolation
			switch strings.ToLower(hardware) {
			case "":
				hardware = hardwareShared
			case isolationPrivate:
				hardware = hardwareDedicated
			case isolationPublic:
				hardware = hardwareShared
			}
			d.Set("hardware", hardware)
		}

		defaultWorkerPool, err := workerPoolsAPI.GetWorkerPool(clusterID, poolName, targetEnv)
		if err != nil {
			return err
		}
		d.Set("labels", IgnoreSystemLabels(defaultWorkerPool.Labels))
		zones := defaultWorkerPool.Zones
		for _, zone := range zones {
			if zone.ID == cls.DataCenter {
				d.Set("default_pool_size", zone.WorkerCount)
				break
			}
		}
		d.Set("worker_pools", flattenWorkerPools(workerPools))
	}

	albs, err := albsAPI.ListClusterALBs(clusterID, targetEnv)
	if err != nil && !strings.Contains(err.Error(), "The specified cluster is a lite cluster.") && !strings.Contains(err.Error(), "This operation is not supported for your cluster's version.") && !strings.Contains(err.Error(), "The specified cluster is a free cluster.") {

		return fmt.Errorf("Error retrieving alb's of the cluster %s: %s", clusterID, err)
	}

	d.Set("name", cls.Name)
	d.Set("server_url", cls.ServerURL)
	d.Set("ingress_hostname", cls.IngressHostname)
	d.Set("ingress_secret", cls.IngressSecretName)
	d.Set("region", cls.Region)
	d.Set("service_subnet", cls.ServiceSubnet)
	d.Set("pod_subnet", cls.PodSubnet)
	d.Set("subnet_id", d.Get("subnet_id").(*schema.Set))
	d.Set("workers_info", workers)
	if strings.HasSuffix(cls.MasterKubeVersion, "_openshift") {
		d.Set("kube_version", strings.Split(cls.MasterKubeVersion, "_")[0]+"_openshift")
	} else {
		d.Set("kube_version", strings.Split(cls.MasterKubeVersion, "_")[0])
	}
	d.Set("albs", flattenAlbs(albs, "all"))
	d.Set("resource_group_id", cls.ResourceGroupID)
	d.Set("public_service_endpoint", cls.PublicServiceEndpointEnabled)
	d.Set("private_service_endpoint", cls.PrivateServiceEndpointEnabled)
	d.Set("public_service_endpoint_url", cls.PublicServiceEndpointURL)
	d.Set("private_service_endpoint_url", cls.PrivateServiceEndpointURL)
	d.Set("crn", cls.CRN)
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

func resourceIBMContainerClusterUpdate(d *schema.ResourceData, meta interface{}) error {

	csClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return err
	}

	targetEnv, err := getClusterTargetHeader(d, meta)
	if err != nil {
		return err
	}

	subnetAPI := csClient.Subnets()
	whkAPI := csClient.WebHooks()
	wrkAPI := csClient.Workers()
	clusterAPI := csClient.Clusters()
	kmsAPI := csClient.Kms()

	clusterID := d.Id()

	if (d.HasChange("kube_version") || d.HasChange("update_all_workers") || d.HasChange("patch_version") || d.HasChange("retry_patch_version")) && !d.IsNewResource() {
		if d.HasChange("kube_version") {
			var masterVersion string
			if v, ok := d.GetOk("kube_version"); ok {
				masterVersion = v.(string)
			}
			params := v1.ClusterUpdateParam{
				Action:  "update",
				Force:   true,
				Version: masterVersion,
			}
			err := clusterAPI.Update(clusterID, params, targetEnv)
			if err != nil {
				return err
			}
			_, err = WaitForClusterVersionUpdate(d, meta, targetEnv)
			if err != nil {
				return fmt.Errorf(
					"Error waiting for cluster (%s) version to be updated: %s", d.Id(), err)
			}
		}
		// "update_all_workers" deafult is false, enable to true when all worker nodes to be updated
		// with major and minor updates.
		updateAllWorkers := d.Get("update_all_workers").(bool)
		if updateAllWorkers || d.HasChange("patch_version") || d.HasChange("retry_patch_version") {
			patchVersion := d.Get("patch_version").(string)
			workerFields, err := wrkAPI.List(clusterID, targetEnv)
			if err != nil {
				return fmt.Errorf("Error retrieving workers for cluster: %s", err)
			}
			cluster, err := clusterAPI.Find(clusterID, targetEnv)
			if err != nil {
				return fmt.Errorf("Error retrieving cluster %s: %s", clusterID, err)
			}

			waitForWorkerUpdate := d.Get("wait_for_worker_update").(bool)

			for _, w := range workerFields {
				/*kubeversion update done if
				1. There is a change in Major.Minor version
				2. Therese is a change in patch_version & Traget kube patch version and patch_version are same
				*/
				if strings.Split(w.KubeVersion, "_")[0] != strings.Split(cluster.MasterKubeVersion, "_")[0] || (strings.Split(w.KubeVersion, ".")[2] != patchVersion && strings.Split(w.TargetVersion, ".")[2] == patchVersion) {
					params := v1.WorkerUpdateParam{
						Action: "update",
					}
					err = wrkAPI.Update(clusterID, w.ID, params, targetEnv)
					if err != nil {
						d.Set("patch_version", nil)
						return fmt.Errorf("Error updating worker %s: %s", w.ID, err)
					}
					if waitForWorkerUpdate {
						_, err = WaitForWorkerAvailable(d, meta, targetEnv)
						if err != nil {
							d.Set("patch_version", nil)
							return fmt.Errorf(
								"Error waiting for workers of cluster (%s) to become ready: %s", d.Id(), err)
						}
					}
				}
			}
		}
	}

	if d.HasChange("kms_config") {
		kmsConfig := v1.KmsEnableReq{}
		kmsConfig.Cluster = clusterID
		targetEnv := v1.ClusterHeader{}
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

		err := kmsAPI.EnableKms(kmsConfig, targetEnv)
		if err != nil {
			log.Printf(
				"An error occured during EnableKms (cluster: %s) error: %s", d.Id(), err)
			return err
		}
	}

	if d.HasChange("force_delete_storage") {
		var forceDeleteStorage bool
		if v, ok := d.GetOk("force_delete_storage"); ok {
			forceDeleteStorage = v.(bool)
		}
		d.Set("force_delete_storage", forceDeleteStorage)
	}

	if d.HasChange("default_pool_size") && !d.IsNewResource() {
		workerPoolsAPI := csClient.WorkerPools()
		workerPools, err := workerPoolsAPI.ListWorkerPools(clusterID, targetEnv)
		if err != nil {
			return err
		}
		var poolName string
		var poolContains bool

		if len(workerPools) > 0 && workerPoolContains(workerPools, defaultWorkerPool) {
			poolName = defaultWorkerPool

			poolContains = true
		} else if len(workerPools) > 0 && workerPoolContains(workerPools, computeWorkerPool) && workerPoolContains(workerPools, gatewayWorkerpool) {
			poolName = computeWorkerPool
			poolContains = true
		}
		if poolContains {
			poolSize := d.Get("default_pool_size").(int)
			err = workerPoolsAPI.ResizeWorkerPool(clusterID, poolName, poolSize, targetEnv)
			if err != nil {
				return fmt.Errorf(
					"Error updating the default_pool_size %d: %s", poolSize, err)
			}

			_, err = WaitForWorkerAvailable(d, meta, targetEnv)
			if err != nil {
				return fmt.Errorf(
					"Error waiting for workers of cluster (%s) to become ready: %s", d.Id(), err)
			}
		} else {
			return fmt.Errorf(
				"The default worker pool does not exist. Use ibm_container_worker_pool and ibm_container_worker_pool_zone attachment resources to make changes to your cluster, such as adding zones, adding worker nodes, or updating worker nodes..")
		}
	}

	if d.HasChange("labels") {
		workerPoolsAPI := csClient.WorkerPools()
		workerPools, err := workerPoolsAPI.ListWorkerPools(clusterID, targetEnv)
		if err != nil {
			return err
		}
		var poolName string
		var poolContains bool

		if len(workerPools) > 0 && workerPoolContains(workerPools, defaultWorkerPool) {
			poolName = defaultWorkerPool
			poolContains = true
		} else if len(workerPools) > 0 && workerPoolContains(workerPools, computeWorkerPool) && workerPoolContains(workerPools, gatewayWorkerpool) {
			poolName = computeWorkerPool
			poolContains = true
		}
		if poolContains {
			labels := make(map[string]string)
			if l, ok := d.GetOk("labels"); ok {
				for k, v := range l.(map[string]interface{}) {
					labels[k] = v.(string)
				}
			}
			err = workerPoolsAPI.UpdateLabelsWorkerPool(clusterID, poolName, labels, targetEnv)
			if err != nil {
				return fmt.Errorf(
					"Error updating the labels %s", err)
			}

			_, err = WaitForWorkerAvailable(d, meta, targetEnv)
			if err != nil {
				return fmt.Errorf(
					"Error waiting for workers of cluster (%s) to become ready: %s", d.Id(), err)
			}
		} else {
			return fmt.Errorf(
				"The default worker pool does not exist. Use ibm_container_worker_pool and ibm_container_worker_pool_zone attachment resources to make changes to your cluster, such as adding zones, adding worker nodes, or updating worker nodes..")
		}
	}

	if d.HasChange("worker_num") {
		old, new := d.GetChange("worker_num")
		oldCount := old.(int)
		newCount := new.(int)
		if newCount > oldCount {
			count := newCount - oldCount
			machineType := d.Get("machine_type").(string)
			publicVlanID := d.Get("public_vlan_id").(string)
			privateVlanID := d.Get("private_vlan_id").(string)
			hardware := d.Get("hardware").(string)
			switch strings.ToLower(hardware) {
			case hardwareDedicated:
				hardware = isolationPrivate
			case hardwareShared:
				hardware = isolationPublic
			}
			params := v1.WorkerParam{
				WorkerNum:   count,
				MachineType: machineType,
				PublicVlan:  publicVlanID,
				PrivateVlan: privateVlanID,
				Isolation:   hardware,
			}
			wrkAPI.Add(clusterID, params, targetEnv)
		} else if oldCount > newCount {
			count := oldCount - newCount
			workerFields, err := wrkAPI.List(clusterID, targetEnv)
			if err != nil {
				return fmt.Errorf("Error retrieving workers for cluster: %s", err)
			}
			for i := 0; i < count; i++ {
				err := wrkAPI.Delete(clusterID, workerFields[i].ID, targetEnv)
				if err != nil {
					return fmt.Errorf(
						"Error deleting workers of cluster (%s): %s", d.Id(), err)
				}
			}
		}

		_, err = WaitForWorkerAvailable(d, meta, targetEnv)
		if err != nil {
			return fmt.Errorf(
				"Error waiting for workers of cluster (%s) to become ready: %s", d.Id(), err)
		}
	}

	if d.HasChange("workers_info") {
		oldWorkers, newWorkers := d.GetChange("workers_info")
		oldWorker := oldWorkers.([]interface{})
		newWorker := newWorkers.([]interface{})
		for _, nW := range newWorker {
			newPack := nW.(map[string]interface{})
			for _, oW := range oldWorker {
				oldPack := oW.(map[string]interface{})
				if strings.Compare(newPack["version"].(string), oldPack["version"].(string)) != 0 {
					cluster, err := clusterAPI.Find(clusterID, targetEnv)
					if err != nil {
						return fmt.Errorf("Error retrieving cluster %s: %s", clusterID, err)
					}
					if newPack["version"].(string) != strings.Split(cluster.MasterKubeVersion, "_")[0] {
						return fmt.Errorf("Worker version %s should match the master kube version %s", newPack["version"].(string), strings.Split(cluster.MasterKubeVersion, "_")[0])
					}
					params := v1.WorkerUpdateParam{
						Action: "update",
					}
					err = wrkAPI.Update(clusterID, oldPack["id"].(string), params, targetEnv)
					if err != nil {
						return fmt.Errorf("Error updating worker %s: %s", oldPack["id"].(string), err)
					}

					_, err = WaitForWorkerAvailable(d, meta, targetEnv)
					if err != nil {
						return fmt.Errorf(
							"Error waiting for workers of cluster (%s) to become ready: %s", d.Id(), err)
					}
				}
			}
		}

	}

	//TODO put webhooks can't deleted in the error message if such case is observed in the chnages
	if d.HasChange("webhook") {
		oldHooks, newHooks := d.GetChange("webhook")
		oldHook := oldHooks.([]interface{})
		newHook := newHooks.([]interface{})
		for _, nH := range newHook {
			newPack := nH.(map[string]interface{})
			exists := false
			for _, oH := range oldHook {
				oldPack := oH.(map[string]interface{})
				if (strings.Compare(newPack["level"].(string), oldPack["level"].(string)) == 0) && (strings.Compare(newPack["type"].(string), oldPack["type"].(string)) == 0) && (strings.Compare(newPack["url"].(string), oldPack["url"].(string)) == 0) {
					exists = true
				}
			}
			if !exists {
				webhook := v1.WebHook{
					Level: newPack["level"].(string),
					Type:  newPack["type"].(string),
					URL:   newPack["url"].(string),
				}

				whkAPI.Add(clusterID, webhook, targetEnv)
			}
		}
	}
	//TODO put subnet can't deleted in the error message if such case is observed in the chnages
	var publicSubnetAdded bool
	noSubnet := d.Get("no_subnet").(bool)
	publicVlanID := d.Get("public_vlan_id").(string)
	if noSubnet == false && publicVlanID != "" {
		publicSubnetAdded = true
	}
	if d.HasChange("subnet_id") {
		oldSubnets, newSubnets := d.GetChange("subnet_id")
		oldSubnet := oldSubnets.(*schema.Set)
		newSubnet := newSubnets.(*schema.Set)
		rem := oldSubnet.Difference(newSubnet).List()
		if len(rem) > 0 {
			return fmt.Errorf("Subnet(s) %v cannot be deleted", rem)
		}
		metro := d.Get("datacenter").(string)
		//from datacenter retrive the metro for filtering the subnets
		metro = metro[0:3]
		subnets, err := subnetAPI.List(targetEnv, metro)
		if err != nil {
			return err
		}
		for _, nS := range newSubnet.List() {
			exists := false
			for _, oS := range oldSubnet.List() {
				if strings.Compare(nS.(string), oS.(string)) == 0 {
					exists = true
				}
			}
			if !exists {
				err := subnetAPI.AddSubnet(clusterID, nS.(string), targetEnv)
				if err != nil {
					return err
				}
				subnet := getSubnet(subnets, nS.(string))
				if subnet.Type == PUBLIC_SUBNET_TYPE {
					publicSubnetAdded = true
				}
			}
		}
	}
	if publicSubnetAdded && d.Get("wait_till").(string) == ingressReady {
		_, err = WaitForSubnetAvailable(d, meta, targetEnv)
		if err != nil {
			return fmt.Errorf(
				"Error waiting for initializing ingress hostname and secret: %s", err)
		}
	}

	v := os.Getenv("IC_ENV_TAGS")
	if d.HasChange("tags") || v != "" {
		oldList, newList := d.GetChange("tags")
		cluster, err := clusterAPI.Find(clusterID, targetEnv)
		if err != nil {
			return fmt.Errorf("Error retrieving cluster %s: %s", clusterID, err)
		}
		err = UpdateTagsUsingCRN(oldList, newList, meta, cluster.CRN)
		if err != nil {
			log.Printf(
				"An error occured during update of instance (%s) tags: %s", clusterID, err)
		}

	}

	return resourceIBMContainerClusterRead(d, meta)
}

func getID(d *schema.ResourceData, meta interface{}, clusterID string, oldWorkers []interface{}, workerInfo []map[string]string) (string, error) {
	targetEnv, err := getClusterTargetHeader(d, meta)
	if err != nil {
		return "", err
	}
	csClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return "", err
	}
	workerFields, err := csClient.Workers().List(clusterID, targetEnv)
	if err != nil {
		return "", err
	}
	for _, wF := range workerFields {
		exists := false
		for _, oW := range oldWorkers {
			oldPack := oW.(map[string]interface{})
			if strings.Compare(wF.ID, oldPack["id"].(string)) == 0 || strings.Compare(wF.State, "deleted") == 0 {
				exists = true
			}
		}
		if !exists {
			for i := 0; i < len(workerInfo); i++ {
				pack := workerInfo[i]
				exists = exists || (strings.Compare(wF.ID, pack["id"]) == 0)
			}
			if !exists {
				return wF.ID, nil
			}
		}
	}

	return "", fmt.Errorf("Unable to get ID of worker")
}

func resourceIBMContainerClusterDelete(d *schema.ResourceData, meta interface{}) error {

	targetEnv, err := getClusterTargetHeader(d, meta)
	if err != nil {
		return err
	}
	csClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return err
	}
	clusterID := d.Id()
	forceDeleteStorage := d.Get("force_delete_storage").(bool)
	err = csClient.Clusters().Delete(clusterID, targetEnv, forceDeleteStorage)
	if err != nil {
		return fmt.Errorf("Error deleting cluster: %s", err)
	}
	_, err = waitForClusterDelete(d, meta)
	if err != nil {
		return err
	}
	return nil
}

func waitForClusterDelete(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	targetEnv, err := getClusterTargetHeader(d, meta)
	if err != nil {
		return nil, err
	}
	csClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return nil, err
	}
	clusterID := d.Id()
	stateConf := &resource.StateChangeConf{
		Pending: []string{clusterDeletePending},
		Target:  []string{clusterDeleted},
		Refresh: func() (interface{}, string, error) {
			cluster, err := csClient.Clusters().Find(clusterID, targetEnv)
			if err != nil {
				if apiErr, ok := err.(bmxerror.RequestFailure); ok && (apiErr.StatusCode() == 404) {
					return cluster, clusterDeleted, nil
				}
				return nil, "", err
			}
			return cluster, clusterDeletePending, nil
		},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        60 * time.Second,
		MinTimeout:   10 * time.Second,
		PollInterval: 60 * time.Second,
	}

	return stateConf.WaitForState()
}

// WaitForClusterAvailable Waits for cluster creation
func WaitForClusterAvailable(d *schema.ResourceData, meta interface{}, target v1.ClusterTargetHeader) (interface{}, error) {
	csClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return nil, err
	}
	log.Printf("Waiting for cluster (%s) to be available.", d.Id())
	id := d.Id()

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", clusterProvisioning},
		Target:     []string{clusterNormal},
		Refresh:    clusterStateRefreshFunc(csClient.Clusters(), id, target),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func clusterStateRefreshFunc(client v1.Clusters, instanceID string, target v1.ClusterTargetHeader) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		clusterFields, err := client.FindWithOutShowResourcesCompatible(instanceID, target)
		if err != nil {
			return nil, "", fmt.Errorf("Error retrieving cluster: %s", err)
		}
		// Check active transactions
		log.Println("Checking cluster")
		//Check for cluster state to be normal
		log.Println("Checking cluster state", strings.Compare(clusterFields.State, clusterNormal))
		if strings.Compare(clusterFields.State, clusterNormal) != 0 {
			return clusterFields, clusterProvisioning, nil
		}
		return clusterFields, clusterNormal, nil
	}
}

// waitForClusterMasterAvailable Waits for cluster creation
func waitForClusterMasterAvailable(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	targetEnv, err := getClusterTargetHeader(d, meta)
	if err != nil {
		return nil, err
	}
	csClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return nil, err
	}
	clusterID := d.Id()

	stateConf := &resource.StateChangeConf{
		Pending: []string{deployRequested, deployInProgress},
		Target:  []string{ready},
		Refresh: func() (interface{}, string, error) {
			clusterFields, err := csClient.Clusters().FindWithOutShowResourcesCompatible(clusterID, targetEnv)
			if err != nil {
				return nil, "", fmt.Errorf("Error retrieving cluster: %s", err)
			}

			if clusterFields.MasterStatus == ready {
				return clusterFields, ready, nil
			}
			return clusterFields, deployInProgress, nil
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

// waitForClusterOneWorkerAvailable Waits for cluster creation
func waitForClusterOneWorkerAvailable(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	targetEnv, err := getClusterTargetHeader(d, meta)
	if err != nil {
		return nil, err
	}
	csClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return nil, err
	}
	clusterID := d.Id()

	stateConf := &resource.StateChangeConf{
		Pending: []string{"retry", "deploying", "provisioning"},
		Target:  []string{normal},
		Refresh: func() (interface{}, string, error) {

			workerPoolsAPI := csClient.WorkerPools()
			workerPools, err := workerPoolsAPI.ListWorkerPools(clusterID, targetEnv)
			if err != nil {
				return nil, "", err
			}
			var poolName string
			var poolContains bool

			if len(workerPools) > 0 && workerPoolContains(workerPools, defaultWorkerPool) {
				poolName = defaultWorkerPool
				poolContains = true
			} else if len(workerPools) > 0 && workerPoolContains(workerPools, computeWorkerPool) && workerPoolContains(workerPools, gatewayWorkerpool) {
				poolName = computeWorkerPool
				poolContains = true
			}
			if poolContains {
				wrkAPI := csClient.Workers()
				workersByPool, err := wrkAPI.ListByWorkerPool(clusterID, poolName, false, targetEnv)
				if err != nil {
					return nil, "", fmt.Errorf("Error retrieving workers of default worker pool for cluster: %s", err)
				}
				if len(workersByPool) == 0 {
					return workersByPool, "provisioning", nil
				}
				for _, worker := range workersByPool {

					if worker.State == normal {
						return workersByPool, normal, nil
					}
				}
				return workersByPool, "deploying", nil
			}
			return nil, normal, nil
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

// WaitForWorkerAvailable Waits for worker creation
func WaitForWorkerAvailable(d *schema.ResourceData, meta interface{}, target v1.ClusterTargetHeader) (interface{}, error) {
	csClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return nil, err
	}
	log.Printf("Waiting for worker of the cluster (%s) to be available.", d.Id())
	id := d.Id()

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", workerProvisioning},
		Target:     []string{workerNormal},
		Refresh:    workerStateRefreshFunc(csClient.Workers(), id, target),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func workerStateRefreshFunc(client v1.Workers, instanceID string, target v1.ClusterTargetHeader) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		workerFields, err := client.List(instanceID, target)
		if err != nil {
			return nil, "", fmt.Errorf("Error retrieving workers for cluster: %s", err)
		}
		log.Println("Checking workers...")
		//Done worker has two fields State and Status , so check for those 2
		for _, e := range workerFields {
			if strings.Contains(e.KubeVersion, "pending") || strings.Compare(e.State, workerNormal) != 0 || strings.Compare(e.Status, workerReadyState) != 0 {
				if strings.Compare(e.State, "deleted") != 0 {
					return workerFields, workerProvisioning, nil
				}
			}
		}
		return workerFields, workerNormal, nil
	}
}

func WaitForClusterCreation(d *schema.ResourceData, meta interface{}, target v1.ClusterTargetHeader) (interface{}, error) {
	csClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return nil, err
	}
	log.Printf("Waiting for cluster (%s) to be available.", d.Id())
	ClusterID := d.Id()

	stateConf := &resource.StateChangeConf{
		Pending: []string{"retry", clusterProvisioning},
		Target:  []string{clusterNormal},
		Refresh: func() (interface{}, string, error) {
			workerFields, err := csClient.Workers().List(ClusterID, target)
			log.Println("Total workers: ", len(workerFields))
			if err != nil {
				return nil, "", fmt.Errorf("Error retrieving workers for cluster: %s", err)
			}
			log.Println("Checking workers...")
			//verifying for atleast sing node to be in normal state
			for _, e := range workerFields {
				log.Println("Worker node status: ", e.State)
				if e.State == workerNormal {
					return workerFields, workerNormal, nil
				}
			}
			return workerFields, workerProvisioning, nil
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func WaitForSubnetAvailable(d *schema.ResourceData, meta interface{}, target v1.ClusterTargetHeader) (interface{}, error) {
	csClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return nil, err
	}
	log.Printf("Waiting for Ingress Subdomain and secret being assigned.")
	id := d.Id()

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", workerProvisioning},
		Target:     []string{workerNormal},
		Refresh:    subnetStateRefreshFunc(csClient.Clusters(), id, d, target),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func subnetStateRefreshFunc(client v1.Clusters, instanceID string, d *schema.ResourceData, target v1.ClusterTargetHeader) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		cluster, err := client.FindWithOutShowResourcesCompatible(instanceID, target)
		if err != nil {
			return nil, "", fmt.Errorf("Error retrieving cluster: %s", err)
		}
		if cluster.IngressHostname == "" || cluster.IngressSecretName == "" {
			return cluster, subnetProvisioning, nil
		}
		return cluster, subnetNormal, nil
	}
}

// WaitForClusterVersionUpdate Waits for cluster creation
func WaitForClusterVersionUpdate(d *schema.ResourceData, meta interface{}, target v1.ClusterTargetHeader) (interface{}, error) {
	csClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return nil, err
	}
	log.Printf("Waiting for cluster (%s) version to be updated.", d.Id())
	id := d.Id()

	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"retry", versionUpdating},
		Target:                    []string{clusterNormal},
		Refresh:                   clusterVersionRefreshFunc(csClient.Clusters(), id, d, target),
		Timeout:                   d.Timeout(schema.TimeoutUpdate),
		Delay:                     20 * time.Second,
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 5,
	}

	return stateConf.WaitForState()
}

func clusterVersionRefreshFunc(client v1.Clusters, instanceID string, d *schema.ResourceData, target v1.ClusterTargetHeader) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		clusterFields, err := client.FindWithOutShowResourcesCompatible(instanceID, target)
		if err != nil {
			return nil, "", fmt.Errorf("Error retrieving cluster: %s", err)
		}
		// Check active transactions
		kubeversion := d.Get("kube_version").(string)
		log.Println("Checking cluster version", clusterFields.MasterKubeVersion, d.Get("kube_version").(string))
		if strings.Contains(clusterFields.MasterKubeVersion, "pending") {
			return clusterFields, versionUpdating, nil
		} else if !strings.Contains(clusterFields.MasterKubeVersion, kubeversion) {
			return clusterFields, versionUpdating, nil
		}
		return clusterFields, clusterNormal, nil
	}
}

func resourceIBMContainerClusterExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	csClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return false, err
	}
	targetEnv, err := getClusterTargetHeader(d, meta)
	if err != nil {
		return false, err
	}
	clusterID := d.Id()
	cls, err := csClient.Clusters().FindWithOutShowResourcesCompatible(clusterID, targetEnv)
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

func getSubnet(subnets []v1.Subnet, subnetId string) v1.Subnet {
	for _, subnet := range subnets {
		if subnet.ID == subnetId {
			return subnet
		}
	}
	return v1.Subnet{}
}

func workerPoolContains(workerPools []v1.WorkerPoolResponse, pool string) bool {
	for _, workerPool := range workerPools {
		if workerPool.Name == pool {
			return true
		}
	}
	return false
}
