// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	v1 "github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
)

func resourceIBMContainerBindService() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMContainerBindServiceCreate,
		Read:     resourceIBMContainerBindServiceRead,
		Update:   resourceIBMContainerBindServiceUpdate,
		Delete:   resourceIBMContainerBindServiceDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"cluster_name_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Cluster name or ID",
			},
			"service_instance_name": {
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"service_instance_id"},
				Description:   "serivice instance name",
			},
			"service_instance_id": {
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"service_instance_name"},
				Description:   "Service instance ID",
			},
			"namespace_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "namespace ID",
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
			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Key info",
			},
			"role": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Role info",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The cluster region",
				Deprecated:  "This field is deprecated",
			},
			"resource_group_id": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "ID of the resource group.",
				ForceNew:         true,
				DiffSuppressFunc: applyOnce,
			},
			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "List of tags for the resource",
			},
		},
	}
}

func getClusterTargetHeader(d *schema.ResourceData, meta interface{}) (v1.ClusterTargetHeader, error) {
	_, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return v1.ClusterTargetHeader{}, err
	}

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return v1.ClusterTargetHeader{}, err
	}
	accountID := userDetails.userAccount

	targetEnv := v1.ClusterTargetHeader{
		AccountID: accountID,
	}

	resourceGroup := ""
	if v, ok := d.GetOk("resource_group_id"); ok {
		resourceGroup = v.(string)
		targetEnv.ResourceGroup = resourceGroup
	}
	return targetEnv, nil
}

func resourceIBMContainerBindServiceCreate(d *schema.ResourceData, meta interface{}) error {
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

	bindService := v1.ServiceBindRequest{
		ClusterNameOrID:         clusterNameID,
		ServiceInstanceNameOrID: serviceInstanceNameID,
		NamespaceID:             namespaceID,
	}

	if v, ok := d.GetOk("key"); ok {
		bindService.ServiceKeyGUID = v.(string)
	}

	if v, ok := d.GetOk("role"); ok {
		bindService.Role = v.(string)
	}

	targetEnv, err := getClusterTargetHeader(d, meta)
	if err != nil {
		return err
	}
	_, err = csClient.Clusters().BindService(bindService, targetEnv)
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%s/%s/%s", clusterNameID, serviceInstanceNameID, namespaceID))

	return resourceIBMContainerBindServiceRead(d, meta)
}

func resourceIBMContainerBindServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceIBMContainerBindServiceRead(d *schema.ResourceData, meta interface{}) error {
	csClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return err
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	clusterNameID := parts[0]
	serviceInstanceNameID := parts[1]
	namespaceID := parts[2]

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
	//d.Set(key, boundService.ServiceKeyName)
	//d.Set(key, boundService.ServiceName)
	return nil
}

func resourceIBMContainerBindServiceDelete(d *schema.ResourceData, meta interface{}) error {
	csClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return err
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	clusterNameID := parts[0]
	serviceInstanceNameID := parts[1]
	namespace := parts[2]
	targetEnv, err := getClusterTargetHeader(d, meta)
	if err != nil {
		return err
	}

	err = csClient.Clusters().UnBindService(clusterNameID, namespace, serviceInstanceNameID, targetEnv)
	if err != nil {
		return fmt.Errorf("Error unbinding service: %s", err)
	}
	return nil
}

//Pure Aramda API not available, we can still find by using k8s api
/*
func resourceIBMContainerBindServiceExists(d *schema.ResourceData, meta interface{}) (bool, error) {

}*/
