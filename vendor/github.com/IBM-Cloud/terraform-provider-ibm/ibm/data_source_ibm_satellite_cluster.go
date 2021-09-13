// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"strings"

	"github.com/IBM-Cloud/container-services-go-sdk/kubernetesserviceapiv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMSatelliteCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMSatelliteClusterRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Name or id of the cluster",
				Type:        schema.TypeString,
				Required:    true,
			},
			"location": {
				Description: "Name or id of the location",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"state": {
				Description: "The lifecycle state of the cluster",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"status": {
				Description: "The status of the cluster",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"health": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kube_version": {
				Description: "Kubernetes version",
				Type:        schema.TypeString,
				Computed:    true,
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
						"flavour": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size_per_zone": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"isolation": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default_worker_pool_labels": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"host_labels": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
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
			"resource_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the resource group.",
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
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CRN of resource instance",
			},
			"server_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The server URL",
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
			ResourceGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group name in which resource is provisioned",
			},
			"tags": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceIBMSatelliteClusterRead(d *schema.ResourceData, meta interface{}) error {
	var name string
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}

	var resourceGrp string
	if v, ok := d.GetOk("resource_group_id"); ok {
		resourceGrp = v.(string)
	}

	satClient, err := meta.(ClientSession).SatelliteClientSession()
	if err != nil {
		return err
	}

	getSatClusterOptions := &kubernetesserviceapiv1.GetClusterOptions{}
	getSatClusterOptions.Cluster = &name
	if resourceGrp != "" {
		getSatClusterOptions.XAuthResourceGroup = &resourceGrp
	}

	clusterFields, response, err := satClient.GetCluster(getSatClusterOptions)
	if err != nil {
		return fmt.Errorf("Error retrieving cluster: %s\n%s", err, response)
	}

	getWorkersOptions := &kubernetesserviceapiv1.GetWorkers1Options{}
	getWorkersOptions.Cluster = &name
	if resourceGrp != "" {
		getWorkersOptions.XAuthResourceGroup = &resourceGrp
	}
	workerFields, response, err := satClient.GetWorkers1(getWorkersOptions)
	if err != nil {
		return fmt.Errorf("Error retrieving workers for satellite cluster: %s", err)
	}
	workers := make([]string, len(workerFields))
	for i, worker := range workerFields {
		workers[i] = *worker.ID
	}

	getSatWorkerPoolOptions := &kubernetesserviceapiv1.GetWorkerPools1Options{}
	getSatWorkerPoolOptions.Cluster = &name
	if resourceGrp != "" {
		getSatWorkerPoolOptions.XAuthResourceGroup = &resourceGrp
	}
	if v, ok := d.GetOk("region"); ok {
		wpRegion := v.(string)
		getSatWorkerPoolOptions.XRegion = &wpRegion
	}

	workerPools, response, err := satClient.GetWorkerPools1(getSatWorkerPoolOptions)
	if err != nil {
		return fmt.Errorf("Error retrieving worker pools of the cluster %s: %s\n%s", name, err, response)
	}

	d.SetId(*clusterFields.ID)
	d.Set("crn", *clusterFields.Crn)
	d.Set("location", *clusterFields.Location)
	d.Set("kube_version", *clusterFields.MasterKubeVersion)
	d.Set("worker_count", *clusterFields.WorkerCount)
	d.Set("state", *clusterFields.State)
	d.Set("status", *clusterFields.Status)
	d.Set("workers", workers)

	d.Set("worker_pools", flattenSatelliteWorkerPools(workerPools))

	if clusterFields.ServiceEndpoints != nil {
		d.Set("public_service_endpoint", *clusterFields.ServiceEndpoints.PublicServiceEndpointEnabled)
		d.Set("private_service_endpoint", *clusterFields.ServiceEndpoints.PrivateServiceEndpointEnabled)
		d.Set("public_service_endpoint_url", *clusterFields.ServiceEndpoints.PublicServiceEndpointURL)
		d.Set("private_service_endpoint_url", *clusterFields.ServiceEndpoints.PrivateServiceEndpointURL)
	}
	d.Set("server_url", *clusterFields.MasterURL)

	if clusterFields.Ingress != nil {
		d.Set("ingress_hostname", *clusterFields.Ingress.Hostname)
		d.Set("ingress_secret", *clusterFields.Ingress.SecretName)
	}
	d.Set("resource_group_id", *clusterFields.ResourceGroup)
	d.Set(ResourceGroupName, *clusterFields.ResourceGroupName)

	if clusterFields.Lifecycle != nil {
		d.Set("health", *clusterFields.Lifecycle.MasterHealth)
	}

	if strings.HasSuffix(*clusterFields.MasterKubeVersion, _OPENSHIFT) {
		d.Set("kube_version", strings.Split(*clusterFields.MasterKubeVersion, "_")[0]+_OPENSHIFT)
	} else {
		d.Set("kube_version", strings.Split(*clusterFields.MasterKubeVersion, "_")[0])
	}

	tags, err := GetTagsUsingCRN(meta, *clusterFields.Crn)
	if err != nil {
		log.Printf(
			"An error occured during reading of instance (%s) tags : %s", d.Id(), err)
	}
	d.Set("tags", tags)

	return nil
}
