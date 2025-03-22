// Copyright IBM Corp. 2017, 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package kubernetes

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	v1 "github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
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

const (
	DisableOutboundTrafficProtectionFlag = "disable_outbound_traffic_protection"
	EnableSecureByDefaultFlag            = "enable_secure_by_default"
)

func ResourceIBMContainerVpcCluster() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMContainerVpcClusterCreate,
		Read:     resourceIBMContainerVpcClusterRead,
		Update:   resourceIBMContainerVpcClusterUpdate,
		Delete:   resourceIBMContainerVpcClusterDelete,
		Exists:   resourceIBMContainerVpcClusterExists,
		Importer: &schema.ResourceImporter{},

		CustomizeDiff: customdiff.Sequence(
			func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
				return flex.ResourceTagsCustomizeDiff(diff)
			},
			func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
				return flex.OnlyInUpdateDiff([]string{EnableSecureByDefaultFlag}, diff)
			},
		),

		Schema: map[string]*schema.Schema{

			"flavor": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "Cluster nodes flavour",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The cluster name",
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
						"account_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Account ID of KMS instance holder - if not provided, defaults to the account in use",
						},
						"wait_for_apply": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Forces terraform to wait till the changes take effect, not marking the cluster complete till",
						},
					},
				},
			},

			"zones": {
				Type:             schema.TypeSet,
				Required:         true,
				Description:      "Zone info",
				DiffSuppressFunc: flex.ApplyOnce,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Zone for the worker pool in a multizone cluster",
							DiffSuppressFunc: flex.ApplyOnce,
						},

						"subnet_id": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "The VPC subnet to assign the cluster",
							DiffSuppressFunc: flex.ApplyOnce,
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
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          1,
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "Number of worker nodes in the default worker pool",
			},

			"worker_labels": {
				Type:             schema.TypeMap,
				Optional:         true,
				Computed:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "Labels for default worker pool",
			},

			"operating_system": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "The operating system of the workers in the default worker pool.",
			},

			"secondary_storage": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "The secondary storage option for the default worker pool.",
			},

			"taints": {
				Type:             schema.TypeSet,
				Optional:         true,
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "Taints for the default worker pool",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: flex.ApplyOnce,
							Description:      "Key for taint",
						},
						"value": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: flex.ApplyOnce,
							Description:      "Value for taint.",
						},
						"effect": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: flex.ApplyOnce,
							Description:      "Effect for taint. Accepted values are NoSchedule, PreferNoSchedule and NoExecute.",
							ValidateFunc: validate.InvokeValidator(
								"ibm_container_vpc_cluster",
								"effect"),
						},
					},
				},
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
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_container_vpc_cluster", "tags")},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of tags for the resources",
			},

			"wait_till": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          ingressReady,
				DiffSuppressFunc: flex.ApplyOnce,
				ValidateFunc:     validation.StringInSlice([]string{masterNodeReady, oneWorkerNodeReady, ingressReady, clusterNormal}, true),
				Description:      "wait_till can be configured for Master Ready, One worker Ready, Ingress Ready or Normal",
			},

			"entitlement": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "Entitlement option reduces additional OCP Licence cost in Openshift Clusters",
			},

			"cos_instance_crn": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "A standard cloud object storage instance CRN to back up the internal registry in your OpenShift on VPC Gen 2 cluster",
			},

			"force_delete_storage": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Force the removal of a cluster and its persistent storage. Deleted data cannot be recovered",
			},

			flex.ResourceControllerURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the IBM Cloud dashboard that can be used to explore and view details about this cluster",
			},

			"kms_instance_id": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "Instance ID for boot volume encryption",
				RequiredWith:     []string{"crk"},
			},

			"crk": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "Root Key ID for boot volume encryption",
				RequiredWith:     []string{"kms_instance_id"},
			},
			"kms_account_id": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "Account ID of kms instance holder - if not provided, defaults to the account in use",
				RequiredWith:     []string{"kms_instance_id", "crk"},
			},

			"security_groups": {
				Type:             schema.TypeSet,
				Optional:         true,
				Description:      "Allow user to set which security groups added to their workers",
				Elem:             &schema.Schema{Type: schema.TypeString},
				Set:              flex.ResourceIBMVPCHash,
				DiffSuppressFunc: flex.ApplyOnce,
			},

			DisableOutboundTrafficProtectionFlag: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Allow outbound connections to public destinations",
			},

			EnableSecureByDefaultFlag: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable Secure-by-default on existing clusters (note: can be used on existing clusters)",
				ValidateFunc: func(i interface{}, s string) (warnings []string, errors []error) {
					v := i.(bool)
					if !v {
						// The field can only be true
						errors = append(errors, fmt.Errorf("%s can be only true", s))
						return warnings, errors
					}

					return warnings, errors
				},
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

			"vpe_service_endpoint_url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CRN of resource instance",
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

			"image_security_enforcement": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Set true to enable image security enforcement policies",
			},

			"host_pool_id": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "The ID of the default worker pool's associated host pool",
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

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(45 * time.Minute),
		},
	}
}

func ResourceIBMContainerVpcClusterValidator() *validate.ResourceValidator {
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

	ibmContainerVpcClusteresourceValidator := validate.ResourceValidator{ResourceName: "ibm_container_vpc_cluster", Schema: validateSchema}
	return &ibmContainerVpcClusteresourceValidator
}

func resourceIBMContainerVpcClusterCreate(d *schema.ResourceData, meta interface{}) error {

	vpcProvider := "vpc-gen2"

	csClient, err := meta.(conns.ClientSession).VpcContainerAPI()
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
	imageSecurityEnabled := d.Get("image_security_enforcement").(bool)

	// timeoutStage will define the timeout stage
	var timeoutStage string
	if v, ok := d.GetOk("wait_till"); ok {
		timeoutStage = strings.ToLower(v.(string))
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
		CommonWorkerPoolConfig: v2.CommonWorkerPoolConfig{
			VpcID:       vpcID,
			Flavor:      flavor,
			WorkerCount: workerCount,
			Zones:       zonesList,
		},
	}

	if hpid, ok := d.GetOk("host_pool_id"); ok {
		workerpool.HostPoolID = hpid.(string)
	}

	if os, ok := d.GetOk("operating_system"); ok {
		workerpool.OperatingSystem = os.(string)
	}

	if secondarystorage, ok := d.GetOk("secondary_storage"); ok {
		workerpool.SecondaryStorageOption = secondarystorage.(string)
	}

	if kmsid, ok := d.GetOk("kms_instance_id"); ok {
		crk := d.Get("crk").(string)
		wve := v2.WorkerVolumeEncryption{
			KmsInstanceID:     kmsid.(string),
			WorkerVolumeCRKID: crk,
		}
		if kmsaccid, ok := d.GetOk("kms_account_id"); ok {
			wve.KMSAccountID = kmsaccid.(string)
		}
		workerpool.WorkerVolumeEncryption = &wve
	}

	if l, ok := d.GetOk("worker_labels"); ok {
		labels := make(map[string]string)
		for k, v := range l.(map[string]interface{}) {
			labels[k] = v.(string)
		}
		workerpool.Labels = labels
	}

	disableOutboundTrafficProtection := d.Get(DisableOutboundTrafficProtectionFlag).(bool)

	params := v2.ClusterCreateRequest{
		DisablePublicServiceEndpoint:     disablePublicServiceEndpoint,
		Name:                             name,
		KubeVersion:                      kubeVersion,
		PodSubnet:                        podSubnet,
		ServiceSubnet:                    serviceSubnet,
		WorkerPools:                      workerpool,
		Provider:                         vpcProvider,
		DisableOutboundTrafficProtection: disableOutboundTrafficProtection,
	}

	// Update params with Entitlement option if provided
	if v, ok := d.GetOk("entitlement"); ok {
		params.DefaultWorkerPoolEntitlement = v.(string)
	}

	// Update params with Cloud Object Store instance CRN id option if provided
	if v, ok := d.GetOk("cos_instance_crn"); ok {
		params.CosInstanceCRN = v.(string)
	}

	if v, ok := d.GetOk("security_groups"); ok {
		securityGroups := flex.FlattenSet(v.(*schema.Set))
		params.SecurityGroupIDs = securityGroups
	}

	targetEnv, err := getVpcClusterTargetHeader(d)
	if err != nil {
		return err
	}

	cls, err := csClient.Clusters().Create(params, targetEnv)
	if err != nil {
		return err
	}

	d.SetId(cls.ID)

	if imageSecurityEnabled {
		err = csClient.Clusters().EnableImageSecurityEnforcement(cls.ID, targetEnv)
		if err != nil {
			return err
		}
	}

	err = waitForVpcCluster(d, meta, timeoutStage, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	var taints []interface{}
	if taintRes, ok := d.GetOk("taints"); ok {
		taints = taintRes.(*schema.Set).List()
	}
	if err := updateWorkerpoolTaints(d, meta, cls.ID, "default", taints); err != nil {
		return err
	}

	return resourceIBMContainerVpcClusterUpdate(d, meta)
}

func resourceIBMContainerVpcClusterUpdate(d *schema.ResourceData, meta interface{}) error {

	csClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}

	targetEnv, err := getVpcClusterTargetHeader(d)
	if err != nil {
		return err
	}

	clusterID := d.Id()

	v := os.Getenv("IC_ENV_TAGS")
	if d.HasChange("tags") || v != "" {
		oldList, newList := d.GetChange("tags")
		cluster, err := csClient.Clusters().GetCluster(clusterID, targetEnv)
		if err != nil {
			return fmt.Errorf("[ERROR] Error retrieving cluster %s: %s", clusterID, err)
		}
		err = flex.UpdateTagsUsingCRN(oldList, newList, meta, cluster.CRN)
		if err != nil {
			log.Printf(
				"An error occured during update of instance (%s) tags: %s", clusterID, err)
		}
	}

	if d.HasChange("kms_config") {
		kmsConfig := v2.KmsEnableReq{}
		kmsConfig.Cluster = clusterID
		targetEnv := v2.ClusterHeader{}
		var waitForApply bool
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

				if wait, ok := kmsMap["wait_for_apply"].(bool); ok {
					waitForApply = wait
				}
			}

			err := csClient.Kms().EnableKms(kmsConfig, targetEnv)
			if err != nil {
				log.Printf(
					"An error occured during EnableKms (cluster: %s) error: %s", d.Id(), err)
				return err
			}
			if waitForApply {
				waitForVpcClusterMasterKMSApply(d, meta)
			}
		}
	}

	if d.HasChange(DisableOutboundTrafficProtectionFlag) || d.HasChange(EnableSecureByDefaultFlag) {
		ClusterClient, err := meta.(conns.ClientSession).VpcContainerAPI()
		if err != nil {
			return err
		}

		Env, err := getVpcClusterTargetHeader(d)
		if err != nil {
			return err
		}

		if d.HasChange(DisableOutboundTrafficProtectionFlag) {
			outbound_traffic_protection := !d.Get(DisableOutboundTrafficProtectionFlag).(bool)
			if err := ClusterClient.VPCs().SetOutboundTrafficProtection(clusterID, outbound_traffic_protection, Env); err != nil {
				return err
			}
		}

		if d.HasChange(EnableSecureByDefaultFlag) {
			enableSecureByDefault := d.Get(EnableSecureByDefaultFlag).(bool)
			if err := ClusterClient.VPCs().EnableSecureByDefault(clusterID, enableSecureByDefault, Env); err != nil {
				return err
			}

		}

	}

	if (d.HasChange("kube_version") || d.HasChange("update_all_workers") || d.HasChange("patch_version") || d.HasChange("retry_patch_version")) && !d.IsNewResource() {

		if d.HasChange("kube_version") {
			ClusterClient, err := meta.(conns.ClientSession).ContainerAPI()
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
			_, err = waitForVpcClusterVersionUpdate(d, meta, targetEnv)
			if err != nil {
				return fmt.Errorf("[ERROR] Error waiting for cluster (%s) version to be updated: %s", d.Id(), err)
			}
		}

		csClient, err := meta.(conns.ClientSession).VpcContainerAPI()
		if err != nil {
			return err
		}
		targetEnv, err := getVpcClusterTargetHeader(d)
		if err != nil {
			return err
		}

		clusterID := d.Id()

		// Update the worker nodes after master node kube-version is updated.
		// workers will store the existing workers info to identify the replaced node
		workersInfo := make(map[string]int)

		updateAllWorkers := d.Get("update_all_workers").(bool)
		if updateAllWorkers || d.HasChange("patch_version") || d.HasChange("retry_patch_version") {

			// patchVersion := d.Get("patch_version").(string)
			workers, err := csClient.Workers().ListWorkers(clusterID, false, targetEnv)
			if err != nil {
				d.Set("patch_version", nil)
				return fmt.Errorf("[ERROR] Error retrieving workers for cluster: %s", err)
			}

			for index, worker := range workers {
				workersInfo[worker.ID] = index
			}
			workersCount := len(workers)

			waitForWorkerUpdate := d.Get("wait_for_worker_update").(bool)

			for _, worker := range workers {
				workerPool, err := csClient.WorkerPools().GetWorkerPool(clusterID, worker.PoolID, targetEnv)
				if err != nil {
					return fmt.Errorf("[ERROR] Error retrieving worker pool: %s", err)
				}

				// check if change is present in MAJOR.MINOR version or in PATCH version
				if worker.KubeVersion.Actual != worker.KubeVersion.Target || worker.LifeCycle.ActualOperatingSystem != workerPool.OperatingSystem {
					_, err := csClient.Workers().ReplaceWokerNode(clusterID, worker.ID, targetEnv)
					// As API returns http response 204 NO CONTENT, error raised will be exempted.
					if err != nil && !strings.Contains(err.Error(), "EmptyResponseBody") {
						d.Set("patch_version", nil)
						return fmt.Errorf("[ERROR] Error replacing the worker node from the cluster: %s", err)
					}

					if waitForWorkerUpdate {
						//1. wait for worker node to delete
						_, deleteError := waitForWorkerNodetoDelete(d, meta, targetEnv, worker.ID)
						if deleteError != nil {
							d.Set("patch_version", nil)
							return fmt.Errorf("[ERROR] Worker node - %s is failed to replace", worker.ID)
						}

						//2. wait for new workerNode
						_, newWorkerError := waitForNewWorker(d, meta, targetEnv, workersCount)
						if newWorkerError != nil {
							d.Set("patch_version", nil)
							return fmt.Errorf("[ERROR] Failed to spawn new worker node")
						}

						//3. Get new worker node ID and update the map
						newWorkerID, index, newNodeError := getNewWorkerID(d, meta, targetEnv, workersInfo)
						if newNodeError != nil {
							d.Set("patch_version", nil)
							return fmt.Errorf("[ERROR] Unable to find the new worker node info")
						}

						delete(workersInfo, worker.ID)
						workersInfo[newWorkerID] = index

						//4. wait for the worker's version update and normal state
						_, Err := waitForVpcClusterWokersVersionUpdate(d, meta, targetEnv, newWorkerID)
						if Err != nil {
							d.Set("patch_version", nil)
							return fmt.Errorf(
								"[ERROR] Error waiting for cluster (%s) worker nodes kube version to be updated: %s", d.Id(), Err)
						}
					}
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

	if d.HasChange("image_security_enforcement") && !d.IsNewResource() {
		var imageSecurity bool
		if v, ok := d.GetOk("image_security_enforcement"); ok {
			imageSecurity = v.(bool)
		}
		if imageSecurity {
			csClient.Clusters().EnableImageSecurityEnforcement(clusterID, targetEnv)
		} else {
			csClient.Clusters().DisableImageSecurityEnforcement(clusterID, targetEnv)
		}
	}

	return resourceIBMContainerVpcClusterRead(d, meta)
}

func resourceIBMContainerVpcClusterRead(d *schema.ResourceData, meta interface{}) error {

	csClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}
	albsAPI := csClient.Albs()

	targetEnv, err := getVpcClusterTargetHeader(d)
	if err != nil {
		return err
	}

	clusterID := d.Id()
	cls, err := csClient.Clusters().GetCluster(clusterID, targetEnv)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving conatiner vpc cluster: %s", err)
	}

	albs, err := albsAPI.ListClusterAlbs(clusterID, targetEnv)
	if err != nil && !strings.Contains(err.Error(), "This operation is not supported for your cluster's version.") {
		return fmt.Errorf("[ERROR] Error retrieving alb's of the cluster %s: %s", clusterID, err)
	}

	d.Set("name", cls.Name)
	d.Set("crn", cls.CRN)
	d.Set("master_status", cls.Lifecycle.MasterStatus)
	if strings.HasSuffix(cls.MasterKubeVersion, "_openshift") {
		d.Set("kube_version", strings.Split(cls.MasterKubeVersion, "_")[0]+"_openshift")
	} else {
		d.Set("kube_version", strings.Split(cls.MasterKubeVersion, "_")[0])
	}
	if cls.Vpcs != nil {
		d.Set("vpc_id", cls.Vpcs[0])
	}
	d.Set("master_url", cls.MasterURL)
	d.Set("service_subnet", cls.ServiceSubnet)
	d.Set("pod_subnet", cls.PodSubnet)
	d.Set("state", cls.State)
	d.Set("ingress_hostname", cls.Ingress.HostName)
	d.Set("ingress_secret", cls.Ingress.SecretName)
	d.Set("albs", flex.FlattenVpcAlbs(albs, "all"))
	d.Set("resource_group_id", cls.ResourceGroupID)
	d.Set("public_service_endpoint_url", cls.ServiceEndpoints.PublicServiceEndpointURL)
	d.Set("private_service_endpoint_url", cls.ServiceEndpoints.PrivateServiceEndpointURL)
	d.Set("vpe_service_endpoint_url", cls.VirtualPrivateEndpointURL)
	if cls.ServiceEndpoints.PublicServiceEndpointEnabled {
		d.Set("disable_public_service_endpoint", false)
	} else {
		d.Set("disable_public_service_endpoint", true)
	}
	d.Set("image_security_enforcement", cls.ImageSecurityEnabled)

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
	d.Set(flex.ResourceControllerURL, controller+"/kubernetes/clusters")
	d.Set(flex.ResourceName, cls.Name)
	d.Set(flex.ResourceCRN, cls.CRN)
	d.Set(flex.ResourceStatus, cls.State)
	d.Set(flex.ResourceGroupName, cls.ResourceGroupName)

	return nil
}

func resourceIBMContainerVpcClusterDelete(d *schema.ResourceData, meta interface{}) error {

	targetEnv, err := getVpcClusterTargetHeader(d)
	if err != nil {
		return err
	}
	csClient, err := meta.(conns.ClientSession).VpcContainerAPI()
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

	forceDeleteStorage := d.Get("force_delete_storage").(bool)
	err = csClient.Clusters().Delete(clusterID, targetEnv, forceDeleteStorage)
	if err != nil {
		return fmt.Errorf("[ERROR] Error deleting cluster: %s", err)
	}
	_, err = waitForVpcClusterDelete(d, meta)
	if err != nil {
		return err
	}

	sess1, err := vpcClient(meta)
	if err == nil {
		listlbOptions := &vpcv1.ListLoadBalancersOptions{}
		lbs, response, err1 := sess1.ListLoadBalancers(listlbOptions)
		if err1 != nil {
			log.Printf("Error Retrieving vpc load balancers: %s\n%s", err1, response)
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
	} else {
		log.Printf("Error connecting to VPC client %s", err)
	}
	return nil
}

func resourceIBMContainerVpcClusterExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	csClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return false, err
	}
	targetEnv, err := getVpcClusterTargetHeader(d)
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
		return false, fmt.Errorf("[ERROR] Error getting container vpc cluster: %s", err)
	}
	return cls.ID == clusterID, nil
}

func vpcClient(meta interface{}) (*vpcv1.VpcV1, error) {
	sess, err := meta.(conns.ClientSession).VpcV1API()
	return sess, err
}

func isWaitForLBDeleted(lbc *vpcv1.VpcV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for  (%s) to be deleted.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", "deleting"},
		Target:     []string{"done", "failed"},
		Refresh:    isLBDeleteRefreshFunc(lbc, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isLBDeleteRefreshFunc(lbc *vpcv1.VpcV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] is lb delete function here")
		getLoadBalancerOptions := &vpcv1.GetLoadBalancerOptions{
			ID: &id,
		}
		lb, response, err := lbc.GetLoadBalancer(getLoadBalancerOptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return lb, "done", nil
			}
			return nil, "failed", fmt.Errorf("[ERROR] The vpc load balancer %s failed to delete: %s\n%s", id, err, response)
		}
		return lb, "deleting", nil
	}
}

func waitForVpcCluster(d *schema.ResourceData, meta interface{}, timeoutStage string, timeout time.Duration) error {
	var err error
	switch timeoutStage {

	case strings.ToLower(clusterNormal):
		pendingStates := []string{clusterDeploying, clusterRequested, clusterPending, clusterDeployed, clusterCritical, clusterWarning}
		_, err = waitForVpcClusterState(d, meta, clusterNormal, pendingStates, timeout)
		if err != nil {
			return err
		}

	case strings.ToLower(masterNodeReady):
		_, err = waitForVpcClusterMasterAvailable(d, meta, timeout)
		if err != nil {
			return err
		}

	case strings.ToLower(oneWorkerNodeReady):
		_, err = waitForVpcClusterOneWorkerAvailable(d, meta, timeout)
		if err != nil {
			return err
		}

	case strings.ToLower(ingressReady):
		_, err = waitForVpcClusterIngressAvailable(d, meta, timeout)
		if err != nil {
			return err
		}
	}

	return nil
}

func waitForVpcClusterDelete(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	targetEnv, err := getVpcClusterTargetHeader(d)
	if err != nil {
		return nil, err
	}
	csClient, err := meta.(conns.ClientSession).VpcContainerAPI()
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

func waitForVpcClusterOneWorkerAvailable(d *schema.ResourceData, meta interface{}, timeout time.Duration) (interface{}, error) {
	targetEnv, err := getVpcClusterTargetHeader(d)
	if err != nil {
		return nil, err
	}
	csClient, err := meta.(conns.ClientSession).VpcContainerAPI()
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
		Timeout:                   timeout,
		Delay:                     10 * time.Second,
		MinTimeout:                5 * time.Second,
		ContinuousTargetOccurence: 3,
	}
	return createStateConf.WaitForState()
}

func waitForVpcClusterState(d *schema.ResourceData, meta interface{}, waitForState string, pendingState []string, timeout time.Duration) (interface{}, error) {
	targetEnv, err := getVpcClusterTargetHeader(d)
	if err != nil {
		return nil, err
	}
	csClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return nil, err
	}
	clusterID := d.Id()
	createStateConf := &resource.StateChangeConf{
		Pending: pendingState,
		Target:  []string{waitForState},
		Refresh: func() (interface{}, string, error) {
			clusterInfo, err := csClient.Clusters().GetCluster(clusterID, targetEnv)
			if err != nil {
				return nil, "", err
			}

			if clusterInfo.State == clusterWarning {
				log.Println("[WARN] Cluster is in Warning State, this may be temporary")
			}
			if clusterInfo.State == clusterCritical {
				log.Println("[WARN] Cluster is in Critical State, this may be temporary")
			}

			return clusterInfo, clusterInfo.State, nil
		},
		Timeout:                   timeout,
		Delay:                     10 * time.Second,
		MinTimeout:                5 * time.Second,
		ContinuousTargetOccurence: 3,
	}
	return createStateConf.WaitForState()
}

func waitForVpcClusterMasterAvailable(d *schema.ResourceData, meta interface{}, timeout time.Duration) (interface{}, error) {
	targetEnv, err := getVpcClusterTargetHeader(d)
	if err != nil {
		return nil, err
	}
	csClient, err := meta.(conns.ClientSession).VpcContainerAPI()
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
				return clusterInfo, deployInProgress, clusterInfoErr
			}

			if clusterInfo.Lifecycle.MasterStatus == ready {
				return clusterInfo, ready, nil
			}
			return clusterInfo, deployInProgress, nil

		},
		Timeout:                   timeout,
		Delay:                     10 * time.Second,
		MinTimeout:                5 * time.Second,
		ContinuousTargetOccurence: 3,
	}
	return createStateConf.WaitForState()
}

func waitForVpcClusterMasterKMSApply(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	log.Printf("[DEBUG] Wait for KMS to apply to master")
	targetEnv, err := getVpcClusterTargetHeader(d)
	if err != nil {
		return nil, err
	}
	csClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return nil, err
	}
	clusterID := d.Id()
	createStateConf := &resource.StateChangeConf{
		Pending: []string{deployRequested, deployInProgress},
		Target:  []string{ready},
		Refresh: func() (interface{}, string, error) {
			log.Printf("[DEBUG] Waiting for KMS to apply to master")
			clusterInfo, clusterInfoErr := csClient.Clusters().GetCluster(clusterID, targetEnv)

			if err != nil || clusterInfoErr != nil {
				return clusterInfo, deployInProgress, clusterInfoErr
			}
			if clusterInfo.Features.KeyProtectEnabled == false {
				log.Printf("[DEBUG] KeyProtectEnabled still false")
				return clusterInfo, deployInProgress, nil
			}

			if clusterInfo.Lifecycle.MasterStatus == ready &&
				clusterInfo.Lifecycle.MasterState == masterDeployed {
				log.Printf("[DEBUG] KMS applied to master")
				return clusterInfo, ready, nil
			}
			return clusterInfo, deployInProgress, nil

		},
		Timeout:                   d.Timeout(schema.TimeoutCreate),
		Delay:                     10 * time.Second,
		MinTimeout:                5 * time.Second,
		ContinuousTargetOccurence: 1,
	}
	return createStateConf.WaitForState()
}

func waitForVpcClusterIngressAvailable(d *schema.ResourceData, meta interface{}, timeout time.Duration) (interface{}, error) {
	targetEnv, err := getVpcClusterTargetHeader(d)
	if err != nil {
		return nil, err
	}
	csClient, err := meta.(conns.ClientSession).VpcContainerAPI()
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
				return clusterInfo, deployInProgress, clusterInfoErr
			}

			if clusterInfo.Ingress.HostName != "" {
				return clusterInfo, ready, nil
			}
			return clusterInfo, deployInProgress, nil

		},
		Timeout:                   timeout,
		Delay:                     10 * time.Second,
		MinTimeout:                5 * time.Second,
		ContinuousTargetOccurence: 3,
	}
	return createStateConf.WaitForState()
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

// waitForVpcClusterVersionUpdate Waits for cluster creation
func waitForVpcClusterVersionUpdate(d *schema.ResourceData, meta interface{}, target v2.ClusterTargetHeader) (interface{}, error) {
	csClient, err := meta.(conns.ClientSession).VpcContainerAPI()
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
		ContinuousTargetOccurence: 3,
	}

	return stateConf.WaitForState()
}

func vpcClusterVersionRefreshFunc(client v2.Clusters, instanceID string, d *schema.ResourceData, target v2.ClusterTargetHeader) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		cls, err := client.GetCluster(instanceID, target)
		if err != nil {
			return nil, "retry", fmt.Errorf("[ERROR] Error retrieving conatiner vpc cluster: %s", err)
		}

		// Check active transactions
		log.Println("Checking cluster version", cls.MasterKubeVersion, d.Get("kube_version").(string))
		if strings.Contains(cls.MasterKubeVersion, "(pending)") {
			return cls, versionUpdating, nil
		}
		return cls, clusterNormal, nil
	}
}

// waitForVpcClusterWokersVersionUpdate Waits for Cluster version Update
func waitForVpcClusterWokersVersionUpdate(d *schema.ResourceData, meta interface{}, target v2.ClusterTargetHeader, workerID string) (interface{}, error) {
	csClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return nil, err
	}

	log.Printf("Waiting for worker (%s) version to be updated.", workerID)
	clusterID := d.Id()
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"retry", versionUpdating},
		Target:                    []string{workerNormal},
		Refresh:                   vpcClusterWorkersVersionRefreshFunc(csClient.Workers(), workerID, clusterID, target),
		Timeout:                   d.Timeout(schema.TimeoutUpdate),
		Delay:                     10 * time.Second,
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 3,
	}

	return stateConf.WaitForState()
}

func vpcClusterWorkersVersionRefreshFunc(client v2.Workers, workerID, clusterID string, target v2.ClusterTargetHeader) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		worker, err := client.Get(clusterID, workerID, target)
		if err != nil {
			return nil, "retry", fmt.Errorf("[ERROR] Error retrieving worker of container vpc cluster: %s", err)
		}
		// Check active updates
		if worker.Health.State == "normal" {
			return worker, workerNormal, nil
		}
		return worker, versionUpdating, nil
	}
}

func waitForWorkerNodetoDelete(d *schema.ResourceData, meta interface{}, targetEnv v2.ClusterTargetHeader, workerID string) (interface{}, error) {

	csClient, err := meta.(conns.ClientSession).VpcContainerAPI()
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
	csClient, err := meta.(conns.ClientSession).VpcContainerAPI()
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
				return workers, "", fmt.Errorf("[ERROR] Error in retriving the list of worker nodes")
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
	csClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return "", -1, err
	}

	clusterID := d.Id()

	workers, err := csClient.Workers().ListWorkers(clusterID, false, targetEnv)
	if err != nil {
		return "", -1, fmt.Errorf("[ERROR] Error in retriving the list of worker nodes")
	}

	for index, worker := range workers {
		if _, ok := workersInfo[worker.ID]; !ok {
			log.Println("found new replaced node: ", worker.ID)
			return worker.ID, index, nil
		}
	}
	return "", -1, fmt.Errorf("[ERROR] no new node found")
}
