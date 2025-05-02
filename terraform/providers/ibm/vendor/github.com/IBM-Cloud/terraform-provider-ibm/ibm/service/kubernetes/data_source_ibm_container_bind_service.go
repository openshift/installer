// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package kubernetes

import (
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMContainerBindService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMContainerBindServiceRead,

		Schema: map[string]*schema.Schema{
			"cluster_name_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cluster name or ID",
				ValidateFunc: validate.InvokeDataSourceValidator(
					"ibm_container_bind_service",
					"cluster_name_id"),
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
func DataSourceIBMContainerBindServiceValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cluster_name_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			Required:                   true,
			CloudDataType:              "cluster",
			CloudDataRange:             []string{"resolved_to:id"}})

	iBMContainerBindServiceValidator := validate.ResourceValidator{ResourceName: "ibm_container_bind_service", Schema: validateSchema}
	return &iBMContainerBindServiceValidator
}

func dataSourceIBMContainerBindServiceRead(d *schema.ResourceData, meta interface{}) error {
	csClient, err := meta.(conns.ClientSession).ContainerAPI()
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
		return fmt.Errorf("[ERROR] Please set either service_instance_name or service_instance_id")
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
