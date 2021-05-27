// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	_OPENSHIFT = "_openshift"
)

func dataSourceIBMContainerVPCCluster() *schema.Resource {
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
						"labels": {
							Type:     schema.TypeMap,
							Computed: true,
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
				ValidateFunc: validateAllowedStringValue([]string{"private", "public", "all"}),
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

func dataSourceIBMContainerClusterVPCRead(d *schema.ResourceData, meta interface{}) error {
	csClient, err := meta.(ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}

	targetEnv, err := getVpcClusterTargetHeader(d, meta)
	if err != nil {
		return err
	}

	var clusterID string

	if v, ok := d.GetOk("cluster_name_id"); ok {
		clusterID = v.(string)
	}
	if v, ok := d.GetOk("name"); ok {
		clusterID = v.(string)
	}

	cls, err := csClient.Clusters().GetCluster(clusterID, targetEnv)
	if err != nil {
		return fmt.Errorf("Error retrieving container vpc cluster: %s", err)
	}

	d.SetId(cls.ID)
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
	d.Set("public_service_endpoint", cls.ServiceEndpoints.PublicServiceEndpointEnabled)
	d.Set("private_service_endpoint", cls.ServiceEndpoints.PrivateServiceEndpointEnabled)
	d.Set("ingress_hostname", cls.Ingress.HostName)
	d.Set("ingress_secret", cls.Ingress.SecretName)

	workerFields, err := csClient.Workers().ListWorkers(clusterID, false, targetEnv)
	if err != nil {
		return fmt.Errorf("Error retrieving workers for cluster: %s", err)
	}
	workers := make([]string, len(workerFields))
	for i, worker := range workerFields {
		workers[i] = worker.ID
	}

	d.Set("workers", workers)

	//Get worker pools
	pools, err := csClient.WorkerPools().ListWorkerPools(clusterID, targetEnv)
	if err != nil {
		return fmt.Errorf("Error retrieving worker pools for container vpc cluster: %s", err)
	}

	d.Set("worker_pools", flattenVpcWorkerPools(pools))

	if !strings.HasSuffix(cls.MasterKubeVersion, _OPENSHIFT) {
		albs, err := csClient.Albs().ListClusterAlbs(clusterID, targetEnv)
		if err != nil && !strings.Contains(err.Error(), "The specified cluster is a lite cluster.") {
			return fmt.Errorf("Error retrieving alb's of the cluster %s: %s", clusterID, err)
		}

		filterType := d.Get("alb_type").(string)
		filteredAlbs := flattenVpcAlbs(albs, filterType)

		d.Set("albs", filteredAlbs)
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
	csClientv1, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return err
	}
	apikeyAPI := csClientv1.Apikeys()
	v1targetEnv, err := getClusterTargetHeader(d, meta)
	if err != nil {
		return err
	}
	apikeyConfig, err := apikeyAPI.GetApiKeyInfo(clusterID, v1targetEnv)
	if err != nil {
		return err
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
	d.Set(ResourceControllerURL, controller+"/kubernetes/clusters")
	d.Set(ResourceName, cls.Name)
	d.Set(ResourceCRN, cls.CRN)
	d.Set(ResourceStatus, cls.State)
	d.Set(ResourceGroupName, cls.ResourceGroupName)

	return nil
}
