// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package satellite

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	v1 "github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
	"github.com/IBM-Cloud/container-services-go-sdk/kubernetesserviceapiv1"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jinzhu/copier"
)

const (
	isLocationProvisioning = "provisioning"
	isLocationNormal       = "normal"
	isCluterDeploying      = "deploying"
	isCluterDeployFailed   = "deploy_failed"
	isClusterNormal        = "normal"
	isClusterWarning       = "warning"
	isClusterDeleting      = "deleting"
	isClusterDeleteDone    = "done"
	isWorkerDeployed       = "deployed"
)

func ResourceIBMSatelliteCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceIBMSatelliteClusterCreate,
		Read:   resourceIBMSatelliteClusterRead,
		Update: resourceIBMSatelliteClusterUpdate,
		Delete: resourceIBMSatelliteClusterDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				ID := d.Id()
				satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
				if err != nil {
					return nil, err
				}
				getSatClusterOptions := &kubernetesserviceapiv1.GetClusterOptions{
					Cluster: &ID,
				}

				cluster, response, err := satClient.GetCluster(getSatClusterOptions)
				if err != nil || cluster == nil {
					if response != nil && response.StatusCode == 404 && strings.Contains(err.Error(), "The specified cluster could not be found") {
						return nil, fmt.Errorf("Error reading satellite cluster: %s\n%s", err, response)
					}
					return nil, fmt.Errorf("Error reading satellite cluster: %s", err)
				}

				d.Set("zones", flex.FlattenSatelliteClusterZones(cluster.LocationZones))
				return []*schema.ResourceData{d}, nil
			},
		},

		CustomizeDiff: customdiff.Sequence(
			func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
				return flex.ImmutableResourceCustomizeDiff([]string{"name", "location", "resource_group_id", "crn_token"}, diff)
			},
			func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
				return flex.ResourceTagsCustomizeDiff(diff)
			},
		),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Read:   schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(120 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The unique name for the new IBM Cloud Satellite cluster",
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name or ID of the Satellite location",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the cluster.",
			},
			"infrastructure_topology": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{"single-replica", "highly-available"}),
				Description:  "String value for single node cluster option. Available options: single-replica, highly-available (default: highly-available)",
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
				Description: "The OpenShift Container Platform version",
			},
			"entitlement": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "Entitlement option reduces additional OCP Licence cost in Openshift Clusters",
			},
			"operating_system": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Operating system of the default worker pool. Options are REDHAT_7_64, REDHAT_8_64, or RHCOS.",
			},
			"wait_for_worker_update": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Wait for worker node to update during kube version update.",
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
			"master_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"master_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_config_admin": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "Grant cluster admin access to Satellite Config to manage Kubernetes resources.",
			},
			"worker_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The number of worker nodes per zone in the default worker pool. Required when '--host-label' is specified. (default: 0)",
			},
			"default_worker_pool_labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Labels on the default worker pool",
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
			"pull_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The RedHat pull secret to create the OpenShift cluster",
			},
			"pod_subnet": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "User provided value for the pod subnet",
			},
			"service_subnet": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "User provided value for service subnet",
			},
			"resource_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the resource group.",
			},
			flex.ResourceGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group name in which resource is provisioned",
			},
			"public_service_endpoint_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_service_endpoint_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_service_endpoint_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"private_service_endpoint_enabled": {
				Type:     schema.TypeBool,
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
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_satellite_cluster", "tags")},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of tags for the resources",
			},
			"host_labels": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         flex.ResourceIBMVPCHash,
				Description: "Labels that describe a Satellite host for default workerpool",
			},
			"crn_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The IBM Cloud Identity and Access Management (IAM) service CRN token for the service that creates the cluster.",
			},
			"calico_ip_autodetection": {
				Type:             schema.TypeMap,
				Optional:         true,
				Description:      "Set IP autodetection to use correct interface for Calico",
				Elem:             &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: flex.ApplyOnce,
			},
		},
	}
}

func ResourceIBMSatelliteClusterValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "tags",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9:_ .-]+$`,
			MinValueLength:             1,
			MaxValueLength:             128})

	ibmSatelliteClusteresourceValidator := validate.ResourceValidator{ResourceName: "ibm_satellite_cluster", Schema: validateSchema}
	return &ibmSatelliteClusteresourceValidator
}

func resourceIBMSatelliteClusterCreate(d *schema.ResourceData, meta interface{}) error {
	var resourceGrp, clusterId string
	pathParamsMap := make(map[string]string)
	satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
	if err != nil {
		return err
	}

	targetEnv, err := getClusterTargetHeader(d, meta)
	if err != nil {
		return err
	}

	createClusterOptions := &kubernetesserviceapiv1.CreateSatelliteClusterOptions{}
	name := d.Get("name").(string)
	createClusterOptions.Name = &name

	location := d.Get("location").(string)
	createClusterOptions.Controller = &location

	if v, ok := d.GetOk("resource_group_id"); ok {
		resourceGrp = v.(string)
		pathParamsMap = map[string]string{
			"X-Auth-Resource-Group": resourceGrp,
		}
		createClusterOptions.Headers = pathParamsMap
	}

	//Wait for location to get normal
	_, ok := d.GetOk("crn_token")
	if !ok {
		_, err = waitForLocationNormal(location, d, meta)
		if err != nil {
			return fmt.Errorf("[ERROR] Error waiting for getting location (%s) to be normal: %s", location, err)
		}
	}

	if v, ok := d.GetOk("worker_count"); ok {
		workerCount := int64(v.(int))
		createClusterOptions.WorkerCount = &workerCount
	}

	if v, ok := d.GetOk("kube_version"); ok {
		kube_version := v.(string)
		createClusterOptions.KubeVersion = &kube_version
	}

	if v, ok := d.GetOk("operating_system"); ok {
		operating_system := v.(string)
		createClusterOptions.OperatingSystem = &operating_system
	}

	if v, ok := d.GetOk("infrastructure_topology"); ok {
		infrastructure_topology := v.(string)
		createClusterOptions.InfrastructureTopology = &infrastructure_topology
	}

	if res, ok := d.GetOk("zones"); ok {
		zones := res.(*schema.Set).List()
		for _, e := range zones {
			r, _ := e.(map[string]interface{})
			if ID := r["id"]; ID != nil {
				zone := ID.(string)
				createClusterOptions.Zone = &zone
			}
			break
		}
	}

	if v, ok := d.GetOk("enable_config_admin"); ok {
		enableConfigAdmin := v.(bool)
		createClusterOptions.AdminAgentOptIn = &enableConfigAdmin
	}

	if v, ok := d.GetOk("pod_subnet"); ok {
		podSubnet := v.(string)
		createClusterOptions.PodSubnet = &podSubnet
	}

	if v, ok := d.GetOk("service_subnet"); ok {
		serviceSubnet := v.(string)
		createClusterOptions.ServiceSubnet = &serviceSubnet
	}

	if v, ok := d.GetOk("pull_secret"); ok {
		pullSecret := v.(string)
		createClusterOptions.PullSecret = &pullSecret
	}

	if v, ok := d.GetOk("host_labels"); ok {
		hostLabels := make(map[string]string)
		hl := v.(*schema.Set)
		hostLabels = flex.FlattenKeyValues(hl.List())
		createClusterOptions.Labels = hostLabels
	}

	if v, ok := d.GetOk("entitlement"); ok {
		entitlement := v.(string)
		createClusterOptions.DefaultWorkerPoolEntitlement = &entitlement
	}

	if m, ok := d.GetOk("calico_ip_autodetection"); ok {
		methods := make(map[string]string)
		for k, v := range m.(map[string]interface{}) {
			methods[k] = v.(string)
		}
		createClusterOptions.SetCalicoIPAutodetectionMethods(methods)
	}

	if v, ok := d.GetOk("crn_token"); ok {
		crnToken := v.(string)
		createRemoteClusterOptions := &kubernetesserviceapiv1.CreateSatelliteClusterRemoteOptions{}
		copier.Copy(createRemoteClusterOptions, createClusterOptions)
		createRemoteClusterOptions.XAuthSupplemental = &crnToken

		instance, response, err := satClient.CreateSatelliteClusterRemote(createRemoteClusterOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Creating Satellite Cluster for remote location: %s\n%s", err, response)
		}
		clusterId = *instance.ID
		log.Printf("[INFO] Created ROKS Satellite Cluster for remote location: %s", clusterId)
	} else {
		instance, response, err := satClient.CreateSatelliteCluster(createClusterOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Creating Satellite Cluster: %s\n%s", err, response)
		}
		clusterId = *instance.ID
		log.Printf("[INFO] Created ROKS Satellite Cluster : %s", clusterId)
	}

	d.SetId(clusterId)

	//Create zone in default workerpool
	workerPoolName := "default"
	if res, ok := d.GetOk("zones"); ok {
		zones := res.(*schema.Set).List()
		if len(zones) >= 2 {
			for i := 1; i < len(zones); i++ {
				zone := zones[i].(map[string]interface{})
				if ID := zone["id"]; ID != nil {
					zoneId := ID.(string)
					zoneOptions := &kubernetesserviceapiv1.CreateSatelliteWorkerPoolZoneOptions{
						Cluster:    &clusterId,
						Workerpool: &workerPoolName,
						ID:         &zoneId,
					}
					if pathParamsMap != nil {
						zoneOptions.Headers = pathParamsMap
					}

					response, err := satClient.CreateSatelliteWorkerPoolZone(zoneOptions)
					if err != nil {
						return fmt.Errorf("[ERROR] Error Adding Worker Pool Zone : %s\n%s", err, response)
					}
				}
			}
			_, err = WaitForSatelliteWorkerPoolAvailable(d, meta, clusterId, workerPoolName, d.Timeout(schema.TimeoutCreate), targetEnv)
			if err != nil {
				return fmt.Errorf("[ERROR] Error waiting for default workerpool (%s) to become ready: %s", d.Id(), err)
			}
		}
	}

	if l, ok := d.GetOk("default_worker_pool_labels"); ok {
		labels := make(map[string]string)
		for k, v := range l.(map[string]interface{}) {
			labels[k] = v.(string)
		}

		wpots := &kubernetesserviceapiv1.V2SetWorkerPoolLabelsOptions{
			Cluster:    &clusterId,
			Workerpool: &workerPoolName,
			Labels:     labels,
		}
		if resourceGrp != "" {
			wpots.XAuthResourceGroup = &resourceGrp
		}

		response, err := satClient.V2SetWorkerPoolLabels(wpots)
		if err != nil {
			return fmt.Errorf("[ERROR] Error updating the labels: %s\n%s", err, response)
		}
	}

	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk("tags"); ok || v != "" {
		getSatClusterOptions := &kubernetesserviceapiv1.GetClusterOptions{
			Cluster: flex.PtrToString(clusterId),
		}

		cluster, response, err := satClient.GetCluster(getSatClusterOptions)
		if err != nil || cluster == nil {
			log.Printf(
				"Error in retreiving ibm satellite cluster : %s\n%s", err, response)
		}

		oldList, newList := d.GetChange("tags")
		err = flex.UpdateTagsUsingCRN(oldList, newList, meta, *cluster.Crn)
		if err != nil {
			log.Printf(
				"Error on create of ibm satellite location (%s) tags: %s", d.Id(), err)
		}
	}

	//Wait for cluster to get warning state
	_, err = waitForClusterToReady(clusterId, d, meta)
	if err != nil {
		return fmt.Errorf("[ERROR] Error waiting for getting cluster (%s) to be warning state: %s", clusterId, err)
	}

	return resourceIBMSatelliteClusterRead(d, meta)
}

func resourceIBMSatelliteClusterRead(d *schema.ResourceData, meta interface{}) error {
	ID := d.Id()
	satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
	if err != nil {
		return err
	}

	getSatClusterOptions := &kubernetesserviceapiv1.GetClusterOptions{
		Cluster: &ID,
	}

	cluster, response, err := satClient.GetCluster(getSatClusterOptions)
	if err != nil || cluster == nil {
		if response != nil && response.StatusCode == 404 && strings.Contains(err.Error(), "The specified cluster could not be found") {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", *cluster.Name)
	d.Set("crn", *cluster.Crn)
	d.Set("kube_version", *cluster.MasterKubeVersion)
	if cluster.InfrastructureTopology != nil {
		d.Set("infrastructure_topology", *cluster.InfrastructureTopology)
	}
	d.Set("state", *cluster.State)
	if cluster.Lifecycle != nil {
		d.Set("master_status", *cluster.Lifecycle.MasterStatus)
	}
	d.Set("service_subnet", *cluster.ServiceSubnet)
	d.Set("pod_subnet", *cluster.PodSubnet)
	d.Set("master_url", *cluster.MasterURL)
	d.Set("service_subnet", *cluster.ServiceSubnet)
	d.Set("pod_subnet", *cluster.PodSubnet)
	if cluster.Ingress != nil {
		d.Set("ingress_hostname", *cluster.Ingress.Hostname)
		d.Set("ingress_secret", *cluster.Ingress.SecretName)
	}
	d.Set("resource_group_id", *cluster.ResourceGroup)
	d.Set(flex.ResourceGroupName, *cluster.ResourceGroupName)
	if cluster.ServiceEndpoints != nil {
		d.Set("public_service_endpoint_url", *cluster.ServiceEndpoints.PublicServiceEndpointURL)
		d.Set("private_service_endpoint_url", *cluster.ServiceEndpoints.PrivateServiceEndpointURL)
		d.Set("private_service_endpoint_enabled", *cluster.ServiceEndpoints.PrivateServiceEndpointEnabled)
		d.Set("public_service_endpoint_enabled", *cluster.ServiceEndpoints.PublicServiceEndpointEnabled)
	}

	if *cluster.ServiceEndpoints.PublicServiceEndpointURL != "" {
		d.Set("disable_public_service_endpoint", false)
	} else {
		d.Set("disable_public_service_endpoint", true)
	}

	workerPoolID := "default"
	getWorkerPoolOptions := &kubernetesserviceapiv1.GetWorkerPoolOptions{
		Cluster:    &ID,
		Workerpool: &workerPoolID,
	}

	workerPool, response, err := satClient.GetWorkerPool(getWorkerPoolOptions)
	if err != nil || workerPool == nil {
		log.Printf(
			"An error occured while retrieving default workerpool : %s\n%s", err, response)
	}

	tags, err := flex.GetTagsUsingCRN(meta, *cluster.Crn)
	if err != nil {
		log.Printf(
			"An error occured during reading of instance (%s) tags : %s", d.Id(), err)
	}
	d.Set("tags", tags)
	d.Set("default_worker_pool_labels", flex.IgnoreSystemLabels(workerPool.Labels))
	d.Set("host_labels", flex.FlattenWorkerPoolHostLabels(workerPool.HostLabels))
	d.Set("operating_system", workerPool.OperatingSystem)

	return nil
}

func resourceIBMSatelliteClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	clusterID := d.Id()
	workerPoolName := "default"

	satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
	if err != nil {
		return err
	}

	csClient, err := meta.(conns.ClientSession).ContainerAPI()
	if err != nil {
		return err
	}

	targetEnv, err := getClusterTargetHeader(d, meta)
	if err != nil {
		return err
	}

	wrkAPI := csClient.Workers()
	clusterAPI := csClient.Clusters()
	if (d.HasChange("kube_version") || d.HasChange("patch_version") || d.HasChange("retry_patch_version")) && !d.IsNewResource() {
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
			_, err = WaitForSatelliteClusterVersionUpdate(d, meta, targetEnv)
			if err != nil {
				return fmt.Errorf("[ERROR] Error waiting for cluster (%s) version to be updated: %s", d.Id(), err)
			}
		}
		satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
		if err != nil {
			return err
		}
		workerOpts := satClient.NewGetWorkers1Options(clusterID)
		workerFields, response, err := satClient.GetWorkers1(workerOpts)
		if err != nil {
			return fmt.Errorf("[ERROR] Error retrieving workerFields %s: %s", err, response)
		}

		getSatClusterOptions := &kubernetesserviceapiv1.GetClusterOptions{
			Cluster:            &clusterID,
			XAuthResourceGroup: &targetEnv.ResourceGroup,
		}

		cluster, response, err := satClient.GetCluster(getSatClusterOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error retrieving cluster %s: %s, %s", clusterID, err, response)
		}
		waitForWorkerUpdate := d.Get("wait_for_worker_update").(bool)
		if workerFields != nil {
			for _, w := range workerFields {
				//kubeversion update done if there is a change in Major.Minor version
				if *w.KubeVersion.Actual != *w.KubeVersion.Target {
					params := v1.WorkerUpdateParam{
						Action: "update",
					}
					err = wrkAPI.Update(clusterID, *w.ID, params, targetEnv)
					if err != nil {
						d.Set("patch_version", nil)
						return fmt.Errorf("[ERROR] Error updating worker %s: %s", *w.ID, err)
					}
				}
			}

			if waitForWorkerUpdate {
				_, err = WaitForSatelliteWorkerVersionUpdate(d, meta, *cluster.MasterKubeVersion, targetEnv)
				if err != nil {
					d.Set("patch_version", nil)
					return fmt.Errorf("[ERROR] Error waiting for workers of cluster (%s) to update kube version: %s", clusterID, err)
				}
			}
		}
	}

	if d.HasChange("worker_count") {
		workerCount := int64(d.Get("worker_count").(int))
		resizeOpts := &kubernetesserviceapiv1.V2ResizeWorkerPoolOptions{
			Cluster:            &clusterID,
			Workerpool:         &workerPoolName,
			Size:               &workerCount,
			XAuthResourceGroup: &targetEnv.ResourceGroup,
		}

		response, err := satClient.V2ResizeWorkerPool(resizeOpts)
		if err != nil {
			return fmt.Errorf("[ERROR] Error updating the worker pool size %d: %s\n%s", workerCount, err, response)
		}
	}

	if d.HasChange("default_worker_pool_labels") {
		labels := make(map[string]string)
		if l, ok := d.GetOk("default_worker_pool_labels"); ok {
			for k, v := range l.(map[string]interface{}) {
				labels[k] = v.(string)
			}
		}

		wpots := &kubernetesserviceapiv1.V2SetWorkerPoolLabelsOptions{
			Cluster:            &clusterID,
			Workerpool:         &workerPoolName,
			Labels:             labels,
			XAuthResourceGroup: &targetEnv.ResourceGroup,
		}
		response, err := satClient.V2SetWorkerPoolLabels(wpots)
		if err != nil {
			return fmt.Errorf("[ERROR] Error updating the labels: %s\n%s", err, response)
		}
	}

	v := os.Getenv("IC_ENV_TAGS")
	if d.HasChange("tags") || v != "" {
		oldList, newList := d.GetChange("tags")
		getSatClusterOptions := &kubernetesserviceapiv1.GetClusterOptions{
			Cluster:            &clusterID,
			XAuthResourceGroup: &targetEnv.ResourceGroup,
		}

		cluster, response, err := satClient.GetCluster(getSatClusterOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error retrieving cluster %s: %s\n%s", clusterID, err, response)
		}
		err = flex.UpdateTagsUsingCRN(oldList, newList, meta, *cluster.Crn)
		if err != nil {
			log.Printf(
				"An error occured during update of instance (%s) tags: %s", clusterID, err)
		}
	}

	if d.HasChange("zones") {
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
			_, err = WaitForSatelliteWorkerPoolAvailable(d, meta, clusterID, workerPoolName, d.Timeout(schema.TimeoutUpdate), targetEnv)
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

	return resourceIBMSatelliteClusterRead(d, meta)
}

func resourceIBMSatelliteClusterDelete(d *schema.ResourceData, meta interface{}) error {
	satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
	if err != nil {
		return err
	}

	removeClusterOptions := &kubernetesserviceapiv1.RemoveClusterOptions{}
	name := d.Get("name").(string)
	removeClusterOptions.IdOrName = &name

	response, err := satClient.RemoveCluster(removeClusterOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Creating Satellite Location: %s\n%s", err, response)
	}

	//Wait for cluster to get delete
	_, err = waitForClusterToDelete(name, d, meta)
	if err != nil {
		return fmt.Errorf("[ERROR] Error waiting while deleteing cluster (%s) : %s", name, err)
	}

	log.Printf("[INFO] Deleted satellite cluster : %s", name)

	d.SetId("")
	return nil
}

func waitForLocationNormal(location string, d *schema.ResourceData, meta interface{}) (interface{}, error) {
	satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
	if err != nil {
		return false, err
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{isLocationProvisioning},
		Target:  []string{isLocationNormal},
		Refresh: func() (interface{}, string, error) {
			getSatLocOptions := &kubernetesserviceapiv1.GetSatelliteLocationOptions{
				Controller: &location,
			}

			var instance *kubernetesserviceapiv1.MultishiftGetController
			var response *core.DetailedResponse
			var err error
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				instance, response, err = satClient.GetSatelliteLocation(getSatLocOptions)
				if err != nil || instance == nil {
					if response != nil && response.StatusCode == 404 {
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})

			if conns.IsResourceTimeoutError(err) {
				instance, response, err = satClient.GetSatelliteLocation(getSatLocOptions)
			}

			if instance != nil {
				if *instance.State == isLocationNormal {
					return "", isLocationNormal, nil
				}
			}

			return "", isLocationProvisioning, nil
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      60 * time.Second,
		MinTimeout: 60 * time.Second,
	}

	return stateConf.WaitForState()
}

func waitForClusterToReady(cluster string, d *schema.ResourceData, meta interface{}) (interface{}, error) {
	satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
	if err != nil {
		return false, err
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{isCluterDeploying},
		Target:  []string{isClusterNormal, isCluterDeployFailed},
		Refresh: func() (interface{}, string, error) {
			getClusterOptions := &kubernetesserviceapiv1.GetClusterOptions{
				Cluster: &cluster,
			}
			instance, response, err := satClient.GetCluster(getClusterOptions)
			if err != nil {
				return nil, "", fmt.Errorf("[ERROR] Error Getting cluster : %s\n%s", err, response)
			}
			if instance != nil && *instance.State == isCluterDeployFailed {
				return instance, isCluterDeployFailed, fmt.Errorf("[ERROR] The cluster failed to deploy : %s", cluster)
			}

			if instance != nil && (*instance.State == isClusterNormal || *instance.State == isClusterWarning) {
				return instance, isClusterNormal, nil
			}
			return instance, isCluterDeploying, nil
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      60 * time.Second,
		MinTimeout: 60 * time.Second,
	}

	return stateConf.WaitForState()
}

func waitForClusterToDelete(cluster string, d *schema.ResourceData, meta interface{}) (interface{}, error) {
	satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
	if err != nil {
		return false, err
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{isClusterDeleting},
		Target:  []string{isClusterDeleteDone},
		Refresh: func() (interface{}, string, error) {
			getClusterOptions := &kubernetesserviceapiv1.GetClusterOptions{
				Cluster: &cluster,
			}
			cluster, response, err := satClient.GetCluster(getClusterOptions)
			if err != nil {
				if response != nil && response.StatusCode == 404 {
					return cluster, isClusterDeleteDone, nil
				}
				return nil, "", err
			}
			return cluster, isClusterDeleting, nil
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      60 * time.Second,
		MinTimeout: 60 * time.Second,
	}

	return stateConf.WaitForState()
}

// WaitForSatelliteWorkerVersionUpdate Waits for worker creation
func WaitForSatelliteWorkerVersionUpdate(d *schema.ResourceData, meta interface{}, masterVersion string, target v1.ClusterTargetHeader) (interface{}, error) {
	satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
	if err != nil {
		return nil, err
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{versionUpdating},
		Target:  []string{workerNormal},
		Refresh: func() (interface{}, string, error) {
			log.Printf("Waiting for workers of the cluster (%s) to update version.", d.Id())
			id := d.Id()
			workerOpts := satClient.NewGetWorkers1Options(id)
			workerFields, _, err := satClient.GetWorkers1(workerOpts)
			if err != nil {
				return nil, "", err
			}

			// Check active updates
			count := 0
			for _, worker := range workerFields {
				if *worker.KubeVersion.Actual == *worker.KubeVersion.Target {
					count = count + 1
				}
			}
			if count == len(workerFields) {
				return workerFields, workerNormal, nil
			}
			return workerFields, versionUpdating, nil
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

// WaitForSatelliteClusterVersionUpdate Waits for cluster creation
func WaitForSatelliteClusterVersionUpdate(d *schema.ResourceData, meta interface{}, target v1.ClusterTargetHeader) (interface{}, error) {
	csClient, err := meta.(conns.ClientSession).ContainerAPI()
	if err != nil {
		return nil, err
	}
	log.Printf("Waiting for satellite cluster (%s) version to be updated.", d.Id())
	id := d.Id()

	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"retry", versionUpdating},
		Target:                    []string{clusterNormal},
		Refresh:                   satelliteClusterVersionRefreshFunc(csClient.Clusters(), id, d, target),
		Timeout:                   d.Timeout(schema.TimeoutUpdate),
		Delay:                     20 * time.Second,
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 3,
	}

	return stateConf.WaitForState()
}

func satelliteClusterVersionRefreshFunc(client v1.Clusters, instanceID string, d *schema.ResourceData, target v1.ClusterTargetHeader) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		clusterFields, err := client.FindWithOutShowResourcesCompatible(instanceID, target)
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] Error retrieving satellite cluster: %s", err)
		}

		// Check active transactions
		kubeversion := d.Get("kube_version").(string)
		log.Println("Checking cluster version", clusterFields.MasterKubeVersion, kubeversion)
		if strings.Contains(clusterFields.MasterKubeVersion, "pending") {
			return clusterFields, versionUpdating, nil
		}
		return clusterFields, clusterNormal, nil
	}
}
