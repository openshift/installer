// Copyright IBM Corp. 2017, 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package kubernetes

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	_OPENSHIFT = "_openshift"
)

func DataSourceIBMContainerVPCCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMContainerClusterVPCRead,

		Schema: map[string]*schema.Schema{
			"cluster_name_id": {
				Description:  "Name or id of the cluster",
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"cluster_name_id", "name"},
				Deprecated:   "use name instead",
			},
			"name": {
				Description:  "Name or id of the cluster",
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"cluster_name_id", "name"},
				ValidateFunc: validate.InvokeDataSourceValidator(
					"ibm_container_vpc_cluster",
					"name"),
			},
			"wait_till": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{masterNodeReady, oneWorkerNodeReady, ingressReady, clusterNormal}, true),
				Description:  "wait_till can be configured for Master Ready, One worker Ready, Ingress Ready or Normal",
			},
			"wait_till_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      "20",
				Description:  "timeout for wait_till in minutes",
				RequiredWith: []string{"wait_till"},
			},
			"worker_count": {
				Description: "Number of workers",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"workers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
						"flavor": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"worker_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"isolation": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_pool_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"labels": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"operating_system": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The operating system of the workers in the worker pool",
						},
						"secondary_storage": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The optional secondary storage configuration of the workers in the worker pool.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"device_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"raid_configuration": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"profile": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"state": {
							Type:     schema.TypeString,
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
									"subnets": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"primary": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
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
			"alb_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "all",
				ValidateFunc: validate.ValidateAllowedStringValues([]string{"private", "public", "all"}),
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
			"service_subnet": {
				Type:        schema.TypeString,
				Description: "Custom subnet CIDR to provide private IP addresses for services",
				Computed:    true,
			},
			"pod_subnet": {
				Type:        schema.TypeString,
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
			"resource_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the resource group.",
				Computed:    true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_service_endpoint": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"private_service_endpoint": {
				Type:     schema.TypeBool,
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
			"vpe_service_endpoint_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CRN of resource instance",
			},

			"master_url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the cluster master",
			},

			"health": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"kube_version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"api_key_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of APIkey",
			},
			"api_key_owner_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the key owner",
			},
			"api_key_owner_email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "email id of the key owner",
			},
			"image_security_enforcement": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "True if image security enforcement is enabled",
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

func DataSourceIBMContainerVPCClusterValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			Optional:                   true,
			CloudDataType:              "cluster",
			CloudDataRange:             []string{"resolved_to:id"}})

	iBMContainerVPCClusterValidator := validate.ResourceValidator{ResourceName: "ibm_container_vpc_cluster", Schema: validateSchema}
	return &iBMContainerVPCClusterValidator
}
func dataSourceIBMContainerClusterVPCRead(d *schema.ResourceData, meta interface{}) error {
	csClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}

	targetEnv, err := getVpcClusterTargetHeader(d)
	if err != nil {
		return err
	}
	var clusterNameOrID string

	if v, ok := d.GetOk("cluster_name_id"); ok {
		clusterNameOrID = v.(string)
	}
	if v, ok := d.GetOk("name"); ok {
		clusterNameOrID = v.(string)
	}

	// timeoutStage will define the timeout stage
	var timeoutStage string
	var timeout time.Duration = 20 * time.Minute
	if v, ok := d.GetOk("wait_till"); ok {
		timeoutStage = strings.ToLower(v.(string))
		timeoutInt := d.Get("wait_till_timeout").(int)
		timeout = time.Duration(timeoutInt) * time.Minute
	}

	cls, err := csClient.Clusters().GetCluster(clusterNameOrID, targetEnv)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving container vpc cluster: %s", err)
	}

	d.SetId(cls.ID)

	if timeoutStage != "" {
		err = waitForVpcCluster(d, meta, timeoutStage, timeout)
		if err != nil {
			return err
		}

		cls, err = csClient.Clusters().GetCluster(clusterNameOrID, targetEnv)
		if err != nil {
			return fmt.Errorf("[ERROR] Error retrieving container vpc cluster: %s", err)
		}
	}

	d.Set("crn", cls.CRN)
	d.Set("status", cls.Lifecycle.MasterStatus)
	d.Set("health", cls.Lifecycle.MasterHealth)
	if strings.HasSuffix(cls.MasterKubeVersion, _OPENSHIFT) {
		d.Set("kube_version", strings.Split(cls.MasterKubeVersion, "_")[0]+_OPENSHIFT)
	} else {
		d.Set("kube_version", strings.Split(cls.MasterKubeVersion, "_")[0])

	}
	d.Set("master_url", cls.MasterURL)
	d.Set("worker_count", cls.WorkerCount)
	d.Set("service_subnet", cls.ServiceSubnet)
	d.Set("pod_subnet", cls.PodSubnet)
	d.Set("state", cls.State)
	d.Set("resource_group_id", cls.ResourceGroupID)
	d.Set("public_service_endpoint_url", cls.ServiceEndpoints.PublicServiceEndpointURL)
	d.Set("private_service_endpoint_url", cls.ServiceEndpoints.PrivateServiceEndpointURL)
	d.Set("vpe_service_endpoint_url", cls.VirtualPrivateEndpointURL)
	d.Set("public_service_endpoint", cls.ServiceEndpoints.PublicServiceEndpointEnabled)
	d.Set("private_service_endpoint", cls.ServiceEndpoints.PrivateServiceEndpointEnabled)
	d.Set("ingress_hostname", cls.Ingress.HostName)
	d.Set("ingress_secret", cls.Ingress.SecretName)

	workerFields, err := csClient.Workers().ListWorkers(clusterNameOrID, false, targetEnv)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving workers for cluster: %s", err)
	}
	workers := make([]string, len(workerFields))
	for i, worker := range workerFields {
		workers[i] = worker.ID
	}

	d.Set("workers", workers)

	//Get worker pools
	pools, err := csClient.WorkerPools().ListWorkerPools(clusterNameOrID, targetEnv)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving worker pools for container vpc cluster: %s", err)
	}

	d.Set("worker_pools", flex.FlattenVpcWorkerPools(pools))

	if !strings.HasSuffix(cls.MasterKubeVersion, _OPENSHIFT) {
		albs, err := csClient.Albs().ListClusterAlbs(clusterNameOrID, targetEnv)
		if err != nil {
			return fmt.Errorf("[ERROR] Error retrieving alb's of the cluster %s: %s", clusterNameOrID, err)
		}

		filterType := d.Get("alb_type").(string)
		filteredAlbs := flex.FlattenVpcAlbs(albs, filterType)

		d.Set("albs", filteredAlbs)
	}
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
	csClientv1, err := meta.(conns.ClientSession).ContainerAPI()
	if err != nil {
		return err
	}
	apikeyAPI := csClientv1.Apikeys()
	v1targetEnv, err := getClusterTargetHeader(d, meta)
	if err != nil {
		return err
	}
	apikeyConfig, err := apikeyAPI.GetApiKeyInfo(clusterNameOrID, v1targetEnv)
	if err != nil {
		log.Printf("Error in GetApiKeyInfo, %s", err)
		//return err
	}
	if &apikeyConfig != nil {
		if &apikeyConfig.Name != nil {
			d.Set("api_key_id", apikeyConfig.ID)
		}
		if &apikeyConfig.ID != nil {
			d.Set("api_key_owner_name", apikeyConfig.Name)
		}
		if &apikeyConfig.Email != nil {
			d.Set("api_key_owner_email", apikeyConfig.Email)
		}
	}
	d.Set("image_security_enforcement", cls.ImageSecurityEnabled)
	d.Set(flex.ResourceControllerURL, controller+"/kubernetes/clusters")
	d.Set(flex.ResourceName, cls.Name)
	d.Set(flex.ResourceCRN, cls.CRN)
	d.Set(flex.ResourceStatus, cls.State)
	d.Set(flex.ResourceGroupName, cls.ResourceGroupName)

	return nil
}
