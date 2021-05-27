// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMContainerBindService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMContainerBindServiceRead,

		Schema: map[string]*schema.Schema{
			"cluster_name_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cluster name or ID",
			},
			"service_instance_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"service_instance_name"},
				Description:   "Service instance ID",
			},
			"service_instance_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"service_instance_id"},
				Description:   "serivice instance name",
			},
			"namespace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "namespace ID",
			},
			"service_key_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Key info",
			},
		},
	}
}

func dataSourceIBMContainerBindServiceRead(d *schema.ResourceData, meta interface{}) error {
	csClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return err
	}

	clusterNameID := d.Get("cluster_name_id").(string)
	namespaceID := d.Get("namespace_id").(string)
	var serviceInstanceNameID string
	if serviceInstanceName, ok := d.GetOk("service_instance_name"); ok {
		serviceInstanceNameID = serviceInstanceName.(string)
	} else if serviceInstanceID, ok := d.GetOk("service_instance_id"); ok {
		serviceInstanceNameID = serviceInstanceID.(string)
	} else {
		return fmt.Errorf("Please set either service_instance_name or service_instance_id")
	}

	targetEnv, err := getClusterTargetHeader(d, meta)
	if err != nil {
		return err
	}

	boundService, err := csClient.Clusters().FindServiceBoundToCluster(clusterNameID, serviceInstanceNameID, namespaceID, targetEnv)
	if err != nil {
		return err
	}
	d.Set("namespace_id", boundService.Namespace)
	d.Set("service_instance_name", boundService.ServiceName)
	d.Set("service_instance_id", boundService.ServiceID)
	d.Set("service_key_name", boundService.ServiceKeyName)
	d.SetId(fmt.Sprintf("%s/%s/%s", clusterNameID, serviceInstanceNameID, namespaceID))
	return nil
}
