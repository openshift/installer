// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package kubernetes

import (
	"time"

	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMContainerVpcAlbCreateNew() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMContainerVpcAlbCreate,
		Read:     resourceIBMContainerVpcALBRead,
		Update:   resourceIBMContainerVpcALBUpdate,
		Delete:   resourceIBMContainerVpcALBDelete,
		Importer: &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
		},
		Schema: map[string]*schema.Schema{

			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The type of ALB that you want to create.",
			},
			"zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The zone where you want to deploy the ALB.",
			},
			"cluster": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the cluster that the ALB belongs to.",
				ValidateFunc: validate.InvokeValidator(
					"ibm_container_vpc_alb_create",
					"cluster"),
			},
			"resource_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the resource group.",
			},
			"enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable the ALB instance in the cluster",
			},
			//response
			"alb_id": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Description: "The ID of the application load balancer (ALB).",
			},

			//get
			"alb_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the ALB",
			},
			"disable_deployment": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Disable the ALB instance in the cluster",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ALB name",
			},
			"load_balancer_hostname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Load balancer host name",
			},
			"resize": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "boolean value to resize the albs",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ALB state",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the ALB",
			},
		},
	}
}

func resourceIBMContainerVpcAlbCreate(d *schema.ResourceData, meta interface{}) error {

	albClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}
	albAPI := albClient.Albs()

	params := v2.AlbCreateReq{}

	if v, ok := d.GetOkExists("cluster"); ok {
		params.Cluster = v.(string)
	}

	if v, ok := d.GetOkExists("type"); ok {
		params.Type = v.(string)
	}

	if v, ok := d.GetOkExists("zone"); ok {
		params.ZoneAlb = v.(string)
	}

	if v, ok := d.GetOk("enable"); ok {
		params.EnableByDefault = v.(bool)
	}

	targetEnv, _ := getVpcClusterTargetHeader(d)

	//v2.AlbCreateResp
	albResp, err := albAPI.CreateAlb(params, targetEnv)
	if err != nil {
		return err
	}

	d.SetId(albResp.Alb)
	return nil
}
func ResourceIBMContainerVpcAlbCreateNewValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cluster",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			Required:                   true,
			CloudDataType:              "cluster",
			CloudDataRange:             []string{"resolved_to:id"}})

	iBMContainerVpcAlbCreateNewValidator := validate.ResourceValidator{ResourceName: "ibm_container_nlb_dns", Schema: validateSchema}
	return &iBMContainerVpcAlbCreateNewValidator
}
