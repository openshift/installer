// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package kubernetes

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	v1 "github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
)

func ResourceIBMContainerBindService() *schema.Resource {
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
				ValidateFunc: validate.InvokeValidator(
					"ibm_container_bind_service",
					"cluster_name_id"),
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
				DiffSuppressFunc: flex.ApplyOnce,
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
func ResourceIBMContainerBindServiceValidator() *validate.ResourceValidator {
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

func getClusterTargetHeader(d *schema.ResourceData, meta interface{}) (v1.ClusterTargetHeader, error) {
	_, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return v1.ClusterTargetHeader{}, err
	}

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return v1.ClusterTargetHeader{}, err
	}
	accountID := userDetails.UserAccount

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
	csClient, err := meta.(conns.ClientSession).ContainerAPI()
	if err != nil {
		return err
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}
	if len(parts) < 3 {
		return fmt.Errorf("[ERROR] Incorrect ID %s: Id should be a combination of clusterNameID/serviceInstanceNameID/namespaceID", d.Id())
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
	d.Set("cluster_name_id", clusterNameID)
	//d.Set(key, boundService.ServiceKeyName)
	//d.Set(key, boundService.ServiceName)
	return nil
}

func resourceIBMContainerBindServiceDelete(d *schema.ResourceData, meta interface{}) error {
	csClient, err := meta.(conns.ClientSession).ContainerAPI()
	if err != nil {
		return err
	}
	parts, err := flex.IdParts(d.Id())
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
		return fmt.Errorf("[ERROR] Error unbinding service: %s", err)
	}
	return nil
}

//Pure Aramda API not available, we can still find by using k8s api
/*
func resourceIBMContainerBindServiceExists(d *schema.ResourceData, meta interface{}) (bool, error) {

}*/
