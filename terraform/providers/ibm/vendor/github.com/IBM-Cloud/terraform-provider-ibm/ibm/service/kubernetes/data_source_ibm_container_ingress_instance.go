// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package kubernetes

import (
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMContainerIngressInstance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMContainerIngressInstanceRead,
		Schema: map[string]*schema.Schema{
			"instance_crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance CRN id",
			},
			"cluster": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cluster ID",
				ValidateFunc: validate.InvokeDataSourceValidator(
					"ibm_container_ingress_instance",
					"cluster"),
			},
			"instance_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance registration name",
			},
			"secret_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Secret group for the instance registration",
			},
			"is_default": {
				Type:        schema.TypeBool,
				Computed:    true,
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

func DataSourceIBMContainerIngressInstanceValidator() *validate.ResourceValidator {
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

func dataSourceIBMContainerIngressInstanceRead(d *schema.ResourceData, meta interface{}) error {
	ingressClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}

	clusterID := d.Get("cluster").(string)
	name := d.Get("instance_name").(string)

	ingressAPI := ingressClient.Ingresses()
	ingressInstanceConfig, err := ingressAPI.GetIngressInstance(clusterID, name)
	if err != nil {
		return err
	}

	d.Set("cluster", ingressInstanceConfig.Cluster)
	d.Set("instance_name", ingressInstanceConfig.Name)
	d.Set("instance_crn", ingressInstanceConfig.CRN)
	d.Set("is_default", ingressInstanceConfig.IsDefault)
	d.Set("secret_group_id", ingressInstanceConfig.SecretGroupID)
	d.Set("secret_group_name", ingressInstanceConfig.SecretGroupName)
	d.Set("instance_type", ingressInstanceConfig.Type)
	d.Set("status", ingressInstanceConfig.Status)
	d.Set("user_managed", ingressInstanceConfig.UserManaged)
	d.SetId(fmt.Sprintf("%s/%s", clusterID, name))

	return nil
}
