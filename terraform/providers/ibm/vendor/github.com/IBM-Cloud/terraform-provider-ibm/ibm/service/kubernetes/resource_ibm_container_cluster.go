// Copyright IBM Corp. 2017, 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package kubernetes

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	v1 "github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	clusterNormal        = "normal"
	clusterDeletePending = "deleting"
	clusterDeleted       = "deleted"
	clusterDeployed      = "deployed"
	clusterDeploying     = "deploying"
	clusterPending       = "pending"
	clusterRequested     = "requested"
	clusterCritical      = "critical"
	clusterWarning       = "warning"

	workerNormal        = "normal"
	subnetNormal        = "normal"
	workerReadyState    = "Ready"
	workerDeleteState   = "deleted"
	workerDeletePending = "deleting"

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

	masterDeployed = "deployed"
)

const PUBLIC_SUBNET_TYPE = "public"

func ResourceIBMContainerCluster() *schema.Resource {
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
			func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
				return flex.ResourceTagsCustomizeDiff(diff)
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
						"account_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Account ID of KMS instance holder - if not provided, defaults to the account in use",
						},
					},
				},
			},

			"default_pool_size": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          1,
				Description:      "The size of the default worker pool",
				DiffSuppressFunc: flex.ApplyOnce,
				ValidateFunc:     validate.ValidateWorkerNum,
			},

			"labels": {
				Type:             schema.TypeMap,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: flex.ApplyOnce,
				Elem:             &schema.Schema{Type: schema.TypeString},
				Description:      "list of labels to the default worker pool",
			},
			"taints": {
				Type:             schema.TypeSet,
				Optional:         true,
				Description:      "WorkerPool Taints",
				DiffSuppressFunc: flex.ApplyOnce,
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
								"ibm_container_cluster",
								"effect"),
						},
					},
				},
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
				Type:             schema.TypeBool,
				Optional:         true,
				DiffSuppressFunc: flex.ApplyOnce,
				Default:          true,
				Description:      "disc encryption done, if set to true.",
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
				Type:             schema.TypeString,
				DiffSuppressFunc: flex.ApplyOnce,
				Optional:         true,
				Description:      "Machine type",
			},

			"hardware": {
				Type:             schema.TypeString,
				DiffSuppressFunc: flex.ApplyOnce,
				Required:         true,
				ValidateFunc:     validate.ValidateAllowedStringValues([]string{hardwareShared, hardwareDedicated}),
				Description:      "Hardware type",
			},

			"public_vlan_id": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          nil,
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "Public VLAN ID",
			},

			"private_vlan_id": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          nil,
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "Private VLAN ID",
			},

			"entitlement": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "Entitlement option reduces additional OCP Licence cost in Openshift Clusters",
			},

			"operating_system": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: flex.ApplyOnce,
				Computed:         true,
				Description:      "The operating system of the workers in the default worker pool.",
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
				DiffSuppressFunc: flex.ApplyOnce,
				ValidateFunc:     validation.StringInSlice([]string{masterNodeReady, oneWorkerNodeReady, ingressReady, clusterNormal}, true),
				Description:      "wait_till can be configured for Master Ready, One worker Ready, Ingress Ready or Normal",
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
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "Boolean value set to true when subnet creation is not required.",
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
							ValidateFunc: validate.ValidateAllowedStringValues([]string{"slack"}),
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
				DiffSuppressFunc: flex.ApplyOnce,
			},

			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_container_cluster", "tags")},
				Set:         flex.ResourceIBMVPCHash,
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
				DiffSuppressFunc: flex.ApplyOnce,
				Default:          false,
				Description:      "Set true for gateway enabled clusters",
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CRN of resource instance",
			},
			"image_security_enforcement": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Set true to enable image security enforcement policies",
			},
			flex.ResourceControllerURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the IBM Cloud dashboard that can be used to explore and view details about this cluster",
			},

			flex.ResourceName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the resource",
			},

			flex.ResourceCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},

			flex.ResourceStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the resource",
			},

			flex.ResourceGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group name in which resource is provisioned",
			},
		},
	}
}

func ResourceIBMContainerClusterValidator() *validate.ResourceValidator {
	tainteffects := "NoSchedule,PreferNoSchedule,NoExecute"
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "tags",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9:_ .-]+$`,
			MinValueLength:             1,
			MaxValueLength:             128},
		validate.ValidateSchema{
			Identifier:                 "effect",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              tainteffects})

	ibmContainerClusterResourceValidator := validate.ResourceValidator{ResourceName: "ibm_container_cluster", Schema: validateSchema}
	return &ibmContainerClusterResourceValidator
}

func resourceIBMContainerClusterCreate(d *schema.ResourceData, meta interface{}) error {

	csClient, err := meta.(conns.ClientSession).ContainerAPI()
	if err != nil {
		return err
	}
	csClientV2, err := meta.(conns.ClientSession).VpcContainerAPI()
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
	imageSecurityEnabled := d.Get("image_security_enforcement").(bool)

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
				return fmt.Errorf("[ERROR] set private_service_endpoint to true for gateway_enabled clusters")
			}
		} else {
			return fmt.Errorf("[ERROR] set private_service_endpoint to true for gateway_enabled clusters")
		}
	}
	if v, ok := d.GetOk("kube_version"); ok {
		params.MasterVersion = v.(string)
	}
	if v, ok := d.GetOk("operating_system"); ok {
		params.OperatingSystem = v.(string)
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

	targetEnvV2, err := getVpcClusterTargetHeader(d)
	if err != nil {
		return err
	}
	if imageSecurityEnabled {
		err = csClientV2.Clusters().EnableImageSecurityEnforcement(cls.ID, targetEnvV2)
		if err != nil {
			return err
		}
	}

	_, err = waitForClusterMasterAvailable(d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	timeoutStage := strings.ToLower(d.Get("wait_till").(string))
	err = waitForCluster(d, timeoutStage, d.Timeout(schema.TimeoutCreate), meta)
	if err != nil {
		return err
	}

	d.Set("force_delete_storage", d.Get("force_delete_storage").(bool))

	//labels
	workerPoolsAPI := csClient.WorkerPools()
	workerPools, err := workerPoolsAPI.ListWorkerPools(cls.ID, targetEnv)
	if err != nil {
		return err
	}

	if len(workerPools) == 0 || !workerPoolContains(workerPools, defaultWorkerPool) {
		return fmt.Errorf("[ERROR] The default worker pool does not exist. Use ibm_container_worker_pool and ibm_container_worker_pool_zone attachment resources to make changes to your cluster, such as adding zones, adding worker nodes, or updating worker nodes")
	}

	labels := make(map[string]string)
	if l, ok := d.GetOk("labels"); ok {
		for k, v := range l.(map[string]interface{}) {
			labels[k] = v.(string)
		}
	}
	err = workerPoolsAPI.UpdateLabelsWorkerPool(cls.ID, defaultWorkerPool, labels, targetEnv)
	if err != nil {
		return fmt.Errorf("[ERROR] Error updating the labels %s", err)
	}

	//taints
	var taints []interface{}
	if taintRes, ok := d.GetOk("taints"); ok {
		taints = taintRes.(*schema.Set).List()
	}
	if err := updateWorkerpoolTaints(d, meta, cls.ID, defaultWorkerPool, taints); err != nil {
		return err
	}

	return resourceIBMContainerClusterUpdate(d, meta)
}

func waitForCluster(d *schema.ResourceData, timeoutStage string, timeout time.Duration, meta interface{}) error {
	switch timeoutStage {
	case strings.ToLower(masterNodeReady):
		_, err := waitForClusterMasterAvailable(d, meta, timeout)
		if err != nil {
			return err
		}

	case strings.ToLower(oneWorkerNodeReady):
		_, err := waitForClusterOneWorkerAvailable(d, meta, timeout)
		if err != nil {
			return err
		}

	case clusterNormal:
		pendingStates := []string{clusterDeploying, clusterRequested, clusterPending, clusterDeployed, clusterCritical, clusterWarning}
		_, err := waitForClusterState(d, meta, clusterNormal, pendingStates, timeout)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceIBMContainerClusterRead(d *schema.ResourceData, meta interface{}) error {
	csClient, err := meta.(conns.ClientSession).ContainerAPI()
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
		return fmt.Errorf("[ERROR] Error retrieving armada cluster: %s", err)
	}

	workerFields, err := wrkAPI.List(clusterID, targetEnv)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving workers for cluster: %s", err)
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
			return fmt.Errorf("[ERROR] Error retrieving workers of default worker pool for cluster: %s", err)
		}

		// to get the private and public vlan IDs of the gateway enabled cluster.
		if poolName == computeWorkerPool {
			gatewayWorkersByPool, err := wrkAPI.ListByWorkerPool(clusterID, gatewayWorkerpool, false, targetEnv)
			if err != nil {
				return fmt.Errorf("[ERROR] Error retrieving workers of default worker pool for cluster: %s", err)
			}
			d.Set("public_vlan_id", gatewayWorkersByPool[0].PublicVlan)
			d.Set("private_vlan_id", gatewayWorkersByPool[0].PrivateVlan)
		} else {
			d.Set("public_vlan_id", workersByPool[0].PublicVlan)
			d.Set("private_vlan_id", workersByPool[0].PrivateVlan)
		}
		d.Set("machine_type", strings.Split(workersByPool[0].MachineType, ".encrypted")[0])
		if strings.HasSuffix(workersByPool[0].MachineType, ".encrypted") {
			d.Set("disk_encryption", true)
		} else {
			d.Set("disk_encryption", false)
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
		d.Set("labels", flex.IgnoreSystemLabels(defaultWorkerPool.Labels))
		d.Set("operating_system", defaultWorkerPool.OperatingSystem)
		zones := defaultWorkerPool.Zones
		for _, zone := range zones {
			if zone.ID == cls.DataCenter {
				d.Set("default_pool_size", zone.WorkerCount)
				break
			}
		}
		d.Set("worker_pools", flex.FlattenWorkerPools(workerPools))
	}

	albs, err := albsAPI.ListClusterALBs(clusterID, targetEnv)
	if err != nil && !strings.Contains(err.Error(), "This operation is not supported for your cluster's version.") {

		return fmt.Errorf("[ERROR] Error retrieving alb's of the cluster %s: %s", clusterID, err)
	}

	d.Set("name", cls.Name)
	d.Set("server_url", cls.ServerURL)
	d.Set("ingress_hostname", cls.IngressHostname)
	d.Set("ingress_secret", cls.IngressSecretName)
	d.Set("region", cls.Region)
	d.Set("datacenter", cls.DataCenter)
	d.Set("service_subnet", cls.ServiceSubnet)
	d.Set("pod_subnet", cls.PodSubnet)
	d.Set("subnet_id", d.Get("subnet_id").(*schema.Set))
	d.Set("workers_info", workers)
	if strings.HasSuffix(cls.MasterKubeVersion, "_openshift") {
		d.Set("kube_version", strings.Split(cls.MasterKubeVersion, "_")[0]+"_openshift")
	} else {
		d.Set("kube_version", strings.Split(cls.MasterKubeVersion, "_")[0])
	}
	d.Set("albs", flex.FlattenAlbs(albs, "all"))
	d.Set("resource_group_id", cls.ResourceGroupID)
	d.Set("public_service_endpoint", cls.PublicServiceEndpointEnabled)
	d.Set("private_service_endpoint", cls.PrivateServiceEndpointEnabled)
	d.Set("public_service_endpoint_url", cls.PublicServiceEndpointURL)
	d.Set("private_service_endpoint_url", cls.PrivateServiceEndpointURL)
	d.Set("crn", cls.CRN)
	tags, err := flex.GetTagsUsingCRN(meta, cls.CRN)
	if err != nil {
		log.Printf(
			"An error occured during reading of instance (%s) tags : %s", d.Id(), err)
	}
	d.Set("tags", tags)
	controller, err := flex.GetBaseController(meta)
	if err != nil {
		return err
	}
	d.Set("image_security_enforcement", cls.ImageSecurityEnabled)
	d.Set(flex.ResourceControllerURL, controller+"/kubernetes/clusters")
	d.Set(flex.ResourceName, cls.Name)
	d.Set(flex.ResourceCRN, cls.CRN)
	d.Set(flex.ResourceStatus, cls.State)
	d.Set(flex.ResourceGroupName, cls.ResourceGroupName)
	return nil
}

func resourceIBMContainerClusterUpdate(d *schema.ResourceData, meta interface{}) error {

	csClient, err := meta.(conns.ClientSession).ContainerAPI()
	if err != nil {
		return err
	}

	csClientV2, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}

	targetEnv, err := getClusterTargetHeader(d, meta)
	if err != nil {
		return err
	}

	targetEnvV2, err := getVpcClusterTargetHeader(d)
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
				return fmt.Errorf("[ERROR] Error waiting for cluster (%s) version to be updated: %s", d.Id(), err)
			}
		}
		// "update_all_workers" deafult is false, enable to true when all worker nodes to be updated
		// with major and minor updates.
		updateAllWorkers := d.Get("update_all_workers").(bool)
		if updateAllWorkers || d.HasChange("patch_version") || d.HasChange("retry_patch_version") {
			workerFields, err := csClientV2.Workers().ListAllWorkers(clusterID, false, targetEnvV2)
			if err != nil {
				return fmt.Errorf("[ERROR] Error retrieving workers for cluster: %s", err)
			}

			waitForWorkerUpdate := d.Get("wait_for_worker_update").(bool)

			for _, w := range workerFields {
				workerPool, err := csClient.WorkerPools().GetWorkerPool(clusterID, w.PoolID, targetEnv)
				if err != nil {
					return fmt.Errorf("[ERROR] Error retrieving worker pool: %s", err)
				}

				/*kubeversion update done if
				1. There is a change in Major.Minor version
				2. Therese is a change in patch_version & Traget kube patch version and patch_version are same
				*/
				if w.KubeVersion.Actual != w.KubeVersion.Target || w.LifeCycle.ActualOperatingSystem != workerPool.OperatingSystem {
					params := v1.WorkerUpdateParam{
						Action: "update",
					}
					err = wrkAPI.Update(clusterID, w.ID, params, targetEnv)
					if err != nil {
						d.Set("patch_version", nil)
						return fmt.Errorf("[ERROR] Error updating worker %s: %s", w.ID, err)
					}
					if waitForWorkerUpdate {
						_, err = WaitForWorkerAvailable(d, meta, targetEnv)
						if err != nil {
							d.Set("patch_version", nil)
							return fmt.Errorf("[ERROR] Error waiting for workers of cluster (%s) to become ready: %s", d.Id(), err)
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

				//Read optional account id
				if accountid := kmsMap["account_id"]; accountid != nil {
					accountid_string := accountid.(string)
					kmsConfig.AccountID = accountid_string
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
						return fmt.Errorf("[ERROR] workers_info Error retrieving cluster %s: %s", clusterID, err)
					}
					if newPack["version"].(string) != strings.Split(cluster.MasterKubeVersion, "_")[0] {
						return fmt.Errorf("[ERROR] Worker version %s should match the master kube version %s", newPack["version"].(string), strings.Split(cluster.MasterKubeVersion, "_")[0])
					}
					params := v1.WorkerUpdateParam{
						Action: "update",
					}
					err = wrkAPI.Update(clusterID, oldPack["id"].(string), params, targetEnv)
					if err != nil {
						return fmt.Errorf("[ERROR] Error updating worker %s: %s", oldPack["id"].(string), err)
					}

					_, err = WaitForWorkerAvailable(d, meta, targetEnv)
					if err != nil {
						return fmt.Errorf("[ERROR] Error waiting for workers of cluster (%s) to become ready: %s", d.Id(), err)
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
	if !noSubnet && publicVlanID != "" {
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
			return fmt.Errorf("[ERROR] Error waiting for initializing ingress hostname and secret: %s", err)
		}
	}

	v := os.Getenv("IC_ENV_TAGS")
	if d.HasChange("tags") || v != "" {
		oldList, newList := d.GetChange("tags")
		cluster, err := clusterAPI.Find(clusterID, targetEnv)
		if err != nil {
			return fmt.Errorf("[ERROR] tags Error retrieving cluster %s: %s", clusterID, err)
		}
		err = flex.UpdateTagsUsingCRN(oldList, newList, meta, cluster.CRN)
		if err != nil {
			log.Printf(
				"An error occured during update of instance (%s) tags: %s", clusterID, err)
		}

	}

	if d.HasChange("image_security_enforcement") && !d.IsNewResource() {
		var imageSecurity bool
		if v, ok := d.GetOk("image_security_enforcement"); ok {
			imageSecurity = v.(bool)
		}
		if imageSecurity {
			csClientV2.Clusters().EnableImageSecurityEnforcement(clusterID, targetEnvV2)
		} else {
			csClientV2.Clusters().DisableImageSecurityEnforcement(clusterID, targetEnvV2)
		}
	}

	return resourceIBMContainerClusterRead(d, meta)
}

func getID(d *schema.ResourceData, meta interface{}, clusterID string, oldWorkers []interface{}, workerInfo []map[string]string) (string, error) {
	targetEnv, err := getClusterTargetHeader(d, meta)
	if err != nil {
		return "", err
	}
	csClient, err := meta.(conns.ClientSession).ContainerAPI()
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

	return "", fmt.Errorf("[ERROR] Unable to get ID of worker")
}

func resourceIBMContainerClusterDelete(d *schema.ResourceData, meta interface{}) error {

	targetEnv, err := getClusterTargetHeader(d, meta)
	if err != nil {
		return err
	}
	csClient, err := meta.(conns.ClientSession).ContainerAPI()
	if err != nil {
		return err
	}
	clusterID := d.Id()
	forceDeleteStorage := d.Get("force_delete_storage").(bool)
	err = csClient.Clusters().Delete(clusterID, targetEnv, forceDeleteStorage)
	if err != nil {
		return fmt.Errorf("[ERROR] Error deleting cluster: %s", err)
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
	csClient, err := meta.(conns.ClientSession).ContainerAPI()
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

// waitForClusterMasterAvailable Waits for cluster creation
func waitForClusterMasterAvailable(d *schema.ResourceData, meta interface{}, timeout time.Duration) (interface{}, error) {
	targetEnv, err := getClusterTargetHeader(d, meta)
	if err != nil {
		return nil, err
	}
	csClient, err := meta.(conns.ClientSession).ContainerAPI()
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
				return nil, "", fmt.Errorf("[ERROR] waitForClusterMasterAvailable Error retrieving cluster: %s", err)
			}

			if clusterFields.MasterStatus == ready {
				return clusterFields, ready, nil
			}
			return clusterFields, deployInProgress, nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func waitForClusterState(d *schema.ResourceData, meta interface{}, waitForState string, pendingState []string, timeout time.Duration) (interface{}, error) {
	targetEnv, err := getClusterTargetHeader(d, meta)
	if err != nil {
		return nil, err
	}
	csClient, err := meta.(conns.ClientSession).ContainerAPI()
	if err != nil {
		return nil, err
	}
	clusterID := d.Id()

	stateConf := &resource.StateChangeConf{
		Pending: pendingState,
		Target:  []string{waitForState},
		Refresh: func() (interface{}, string, error) {
			cls, err := csClient.Clusters().FindWithOutShowResourcesCompatible(clusterID, targetEnv)
			if err != nil {
				return nil, "", fmt.Errorf("[ERROR] waitForClusterState Error retrieving cluster: %s", err)
			}

			if cls.State == clusterWarning {
				log.Println("[WARN] Cluster is in Warning State, this may be temporary")
			}
			if cls.State == clusterCritical {
				log.Println("[WARN] Cluster is in Critical State, this may be temporary")
			}

			return cls, cls.State, nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

// waitForClusterOneWorkerAvailable Waits for cluster creation
func waitForClusterOneWorkerAvailable(d *schema.ResourceData, meta interface{}, timeout time.Duration) (interface{}, error) {
	targetEnv, err := getClusterTargetHeader(d, meta)
	if err != nil {
		return nil, err
	}
	csClient, err := meta.(conns.ClientSession).ContainerAPI()
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
					return nil, "", fmt.Errorf("[ERROR] Error retrieving workers of default worker pool for cluster: %s", err)
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
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

// WaitForWorkerAvailable Waits for worker creation
func WaitForWorkerAvailable(d *schema.ResourceData, meta interface{}, target v1.ClusterTargetHeader) (interface{}, error) {
	csClient, err := meta.(conns.ClientSession).ContainerAPI()
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
			return nil, "", fmt.Errorf("[ERROR] Error retrieving workers for cluster: %s", err)
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

func WaitForSubnetAvailable(d *schema.ResourceData, meta interface{}, target v1.ClusterTargetHeader) (interface{}, error) {
	csClient, err := meta.(conns.ClientSession).ContainerAPI()
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
			return nil, "", fmt.Errorf("[ERROR]subnetStateRefresh Error retrieving cluster: %s", err)
		}
		if cluster.IngressHostname == "" || cluster.IngressSecretName == "" {
			return cluster, subnetProvisioning, nil
		}
		return cluster, subnetNormal, nil
	}
}

// WaitForClusterVersionUpdate Waits for cluster creation
func WaitForClusterVersionUpdate(d *schema.ResourceData, meta interface{}, target v1.ClusterTargetHeader) (interface{}, error) {
	csClient, err := meta.(conns.ClientSession).ContainerAPI()
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
		ContinuousTargetOccurence: 3,
	}

	return stateConf.WaitForState()
}

func clusterVersionRefreshFunc(client v1.Clusters, instanceID string, d *schema.ResourceData, target v1.ClusterTargetHeader) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		clusterFields, err := client.FindWithOutShowResourcesCompatible(instanceID, target)
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] clusterVersionRefresh Error retrieving cluster: %s", err)
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

	csClient, err := meta.(conns.ClientSession).ContainerAPI()
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
		return false, fmt.Errorf("[ERROR] Error getting container cluster: %s", err)
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
