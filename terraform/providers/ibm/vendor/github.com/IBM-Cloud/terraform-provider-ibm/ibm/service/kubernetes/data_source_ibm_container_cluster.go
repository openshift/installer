// Copyright IBM Corp. 2017, 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package kubernetes

import (
	"fmt"
	"log"
	"strings"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMContainerCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMContainerClusterRead,

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
			"is_trusted": {
				Type:     schema.TypeBool,
				Computed: true,
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
			"bounded_services": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_key_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespace": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"vlans": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
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
									"ips": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"is_public": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"is_byoip": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"cidr": {
										Type:     schema.TypeString,
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
			"ingress_hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ingress_secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
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
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The cluster region",
				Deprecated:  "This field is deprecated",
			},
			"resource_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the resource group.",
				Computed:    true,
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

			"server_url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"list_bounded_services": {
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
				Description: "If set to false bounded services won't be listed.",
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

func dataSourceIBMContainerClusterRead(d *schema.ResourceData, meta interface{}) error {
	csClient, err := meta.(conns.ClientSession).ContainerAPI()
	if err != nil {
		return err
	}
	csAPI := csClient.Clusters()
	wrkAPI := csClient.Workers()
	workerPoolsAPI := csClient.WorkerPools()
	albsAPI := csClient.Albs()

	targetEnv, err := getClusterTargetHeader(d, meta)
	if err != nil {
		return err
	}

	var name string

	if v, ok := d.GetOk("cluster_name_id"); ok {
		name = v.(string)
	}
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}
	clusterFields, err := csAPI.Find(name, targetEnv)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving cluster: %s", err)
	}
	workerFields, err := wrkAPI.List(name, targetEnv)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving workers for cluster: %s", err)
	}
	workers := make([]string, len(workerFields))
	for i, worker := range workerFields {
		workers[i] = worker.ID
	}

	listBoundedServices := d.Get("list_bounded_services").(bool)
	boundedServices := make([]map[string]interface{}, 0)
	if listBoundedServices {
		servicesBoundToCluster, err := csAPI.ListServicesBoundToCluster(name, "", targetEnv)
		if err != nil {
			return fmt.Errorf("[ERROR] Error retrieving services bound to cluster: %s", err)
		}
		for _, service := range servicesBoundToCluster {
			boundedService := make(map[string]interface{})
			boundedService["service_name"] = service.ServiceName
			boundedService["service_id"] = service.ServiceID
			boundedService["service_key_name"] = service.ServiceKeyName
			boundedService["namespace"] = service.Namespace
			boundedServices = append(boundedServices, boundedService)
		}
	}

	workerPools, err := workerPoolsAPI.ListWorkerPools(name, targetEnv)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving worker pools of the cluster %s: %s", name, err)
	}

	albs, err := albsAPI.ListClusterALBs(name, targetEnv)
	if err != nil && !strings.Contains(err.Error(), "The specified cluster is a lite cluster.") && !strings.Contains(err.Error(), "This operation is not supported for your cluster's version.") && !strings.Contains(err.Error(), "The specified cluster is a free cluster.") {
		return fmt.Errorf("[ERROR] Error retrieving alb's of the cluster %s: %s", name, err)
	}

	filterType := d.Get("alb_type").(string)
	filteredAlbs := flex.FlattenAlbs(albs, filterType)

	d.SetId(clusterFields.ID)
	d.Set("worker_count", clusterFields.WorkerCount)
	d.Set("workers", workers)
	d.Set("region", clusterFields.Region)
	d.Set("bounded_services", boundedServices)
	d.Set("vlans", flex.FlattenVlans(clusterFields.Vlans))
	d.Set("is_trusted", clusterFields.IsTrusted)
	d.Set("worker_pools", flex.FlattenWorkerPools(workerPools))
	d.Set("albs", filteredAlbs)
	d.Set("resource_group_id", clusterFields.ResourceGroupID)
	d.Set("public_service_endpoint", clusterFields.PublicServiceEndpointEnabled)
	d.Set("private_service_endpoint", clusterFields.PrivateServiceEndpointEnabled)
	d.Set("public_service_endpoint_url", clusterFields.PublicServiceEndpointURL)
	d.Set("private_service_endpoint_url", clusterFields.PrivateServiceEndpointURL)
	d.Set("crn", clusterFields.CRN)
	d.Set("server_url", clusterFields.ServerURL)
	d.Set("ingress_hostname", clusterFields.IngressHostname)
	d.Set("ingress_secret", clusterFields.IngressSecretName)

	controller, err := flex.GetBaseController(meta)
	if err != nil {
		return err
	}
	d.Set(flex.ResourceControllerURL, controller+"/kubernetes/clusters")
	apikeyAPI := csClient.Apikeys()
	apikeyConfig, err := apikeyAPI.GetApiKeyInfo(name, targetEnv)
	if err != nil {
		log.Printf("[ERROR] Error in GetApiKeyInfo, %s", err)
	}
	d.Set("api_key_id", apikeyConfig.ID)
	d.Set("api_key_owner_name", apikeyConfig.Name)
	d.Set("api_key_owner_email", apikeyConfig.Email)
	d.Set("image_security_enforcement", clusterFields.ImageSecurityEnabled)
	d.Set(flex.ResourceName, clusterFields.Name)
	d.Set(flex.ResourceCRN, clusterFields.CRN)
	d.Set(flex.ResourceStatus, clusterFields.State)
	d.Set(flex.ResourceGroupName, clusterFields.ResourceGroupName)

	return nil
}
