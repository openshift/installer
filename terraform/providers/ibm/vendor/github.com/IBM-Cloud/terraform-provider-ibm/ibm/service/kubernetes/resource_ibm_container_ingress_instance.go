// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package kubernetes

import (
	"fmt"
	"time"

	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMContainerIngressInstance() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMContainerIngressInstanceCreate,
		Read:     resourceIBMContainerIngressInstanceRead,
		Update:   resourceIBMContainerIngressInstanceUpdate,
		Delete:   resourceIBMContainerIngressInstanceDelete,
		Exists:   resourceIBMContainerIngressInstanceExists,
		Importer: &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"instance_crn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance CRN id",
				ForceNew:    true,
			},
			"cluster": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cluster ID",
				ForceNew:    true,
				ValidateFunc: validate.InvokeDataSourceValidator(
					"ibm_container_ingress_instance",
					"cluster"),
			},
			"instance_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance registration name",
			},
			"secret_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Secret group for the instance registration",
			},
			"is_default": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Designates if the instance is the default for the cluster",
			},
			"secret_group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the secret group for the instance",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance registration status",
			},
			"instance_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance type",
			},
			"user_managed": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If the instance was created by the user",
			},
		},
	}
}
func ResourceIBMContainerIngressInstanceValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cluster",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			Required:                   true,
			CloudDataType:              "cluster",
			CloudDataRange:             []string{"resolved_to:id"}})

	iBMContainerIngressInstanceValidator := validate.ResourceValidator{ResourceName: "ibm_container_ingress_instance", Schema: validateSchema}
	return &iBMContainerIngressInstanceValidator
}
func resourceIBMContainerIngressInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	ingressClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}

	instanceCRN := d.Get("instance_crn").(string)
	cluster := d.Get("cluster").(string)

	params := v2.InstanceRegisterConfig{
		CRN:     instanceCRN,
		Cluster: cluster,
	}

	if v, ok := d.GetOk("is_default"); ok {
		params.IsDefault = v.(bool)
	}

	if v, ok := d.GetOk("secret_group_id"); ok {
		params.SecretGroupID = v.(string)
	}

	ingressAPI := ingressClient.Ingresses()
	response, err := ingressAPI.RegisterIngressInstance(params)

	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%s/%s", cluster, response.Name))

	return resourceIBMContainerIngressInstanceRead(d, meta)
}

func resourceIBMContainerIngressInstanceRead(d *schema.ResourceData, meta interface{}) error {
	ingressClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}
	clusterIDOrName := parts[0]
	instanceName := parts[1]

	ingressAPI := ingressClient.Ingresses()
	ingressInstanceConfig, err := ingressAPI.GetIngressInstance(clusterIDOrName, instanceName)
	if err != nil {
		return err
	}
	d.Set("cluster", clusterIDOrName)
	d.Set("instance_name", ingressInstanceConfig.Name)
	d.Set("instance_crn", ingressInstanceConfig.CRN)
	d.Set("is_default", ingressInstanceConfig.IsDefault)
	d.Set("secret_group_id", ingressInstanceConfig.SecretGroupID)
	d.Set("secret_group_name", ingressInstanceConfig.SecretGroupName)
	d.Set("instance_type", ingressInstanceConfig.Type)
	d.Set("status", ingressInstanceConfig.Status)
	d.Set("user_managed", ingressInstanceConfig.UserManaged)

	return nil
}

func resourceIBMContainerIngressInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	ingressClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}

	ingressAPI := ingressClient.Ingresses()

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}
	clusterID := parts[0]
	instanceName := parts[1]

	params := v2.InstanceDeleteConfig{
		Cluster: clusterID,
		Name:    instanceName,
	}

	err = ingressAPI.DeleteIngressInstance(params)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceIBMContainerIngressInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	ingressClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}
	cluster := parts[0]
	instanceName := parts[1]

	params := v2.InstanceUpdateConfig{
		Cluster: cluster,
		Name:    instanceName,
	}

	hasChange := false
	if d.HasChange("is_default") {
		if v, ok := d.GetOk("is_default"); ok {
			params.IsDefault = v.(bool)
		}
		hasChange = true
	}

	if d.HasChange("secret_group_id") {
		if v, ok := d.GetOk("secret_group_id"); ok {
			params.SecretGroupID = v.(string)
		}
		hasChange = true
	}

	if hasChange {
		ingressAPI := ingressClient.Ingresses()
		err = ingressAPI.UpdateIngressInstance(params)
		if err != nil {
			return err
		}
	}

	return resourceIBMContainerIngressInstanceRead(d, meta)
}

func resourceIBMContainerIngressInstanceExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	ingressClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return false, err
	}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return false, err
	}
	clusterID := parts[0]
	instanceName := parts[1]

	ingressAPI := ingressClient.Ingresses()
	ingressInstanceConfig, err := ingressAPI.GetIngressInstance(clusterID, instanceName)

	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("[ERROR] Error getting ingress instance: %s", err)
	}

	return ingressInstanceConfig.Name == instanceName, nil
}
