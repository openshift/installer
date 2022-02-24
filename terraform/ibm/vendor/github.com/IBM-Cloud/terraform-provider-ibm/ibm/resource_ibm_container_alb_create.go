// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	v1 "github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
)

func resourceIBMContainerAlbCreate() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMContainerClassicAlbCreate,
		Read:     resourceIBMContainerALBRead,
		Update:   resourceIBMContainerALBUpdate,
		Delete:   resourceIBMContainerALBDelete,
		Importer: &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			//post req
			"enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If set to true, the ALB is enabled by default.",
			},
			"ingress_image": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The type of Ingress image that you want to use for your ALB deployment.",
			},
			"ip": {
				Type:     schema.TypeString,
				Optional: true,
				//ForceNew:    true,
				Description: "The IP address that you want to assign to the ALB.",
			},
			"nlb_version": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The version of the network load balancer that you want to use for the ALB.",
			},
			"alb_type": {
				Type:     schema.TypeString,
				Required: true,
				//ForceNew:    true,
				Description: "The type of ALB that you want to create.",
			},
			"vlan_id": {
				Type:     schema.TypeString,
				Required: true,
				//ForceNew:    true,
				Description: "The VLAN ID that you want to use for your ALBs.",
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
			},

			//response
			"alb_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the application load balancer (ALB).",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ALB name",
			},
			"disable_deployment": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Set to true if ALB needs to be disabled",
			},
			"user_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP assigned by the user",
			},
			"replicas": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "number of instances",
			},
			"resize": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "resize",
			},
		},
	}
}

func resourceIBMContainerClassicAlbCreate(d *schema.ResourceData, meta interface{}) error {
	log.Println("resourceIBMContainerClassicAlbCreate")
	albClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return err
	}
	albAPI := albClient.Albs()

	// "cluster":"string" //mandatory
	// "enableByDefault": true,
	// "ingressImage": "string",
	// "ip": "string",
	// "nlbVersion": "string",
	// "type": "string", //mandatory
	// "vlanID": "string", //mandatory
	// "zone": "string" //mandatory

	params := v1.CreateALB{}

	if v, ok := d.GetOkExists("alb_type"); ok {
		params.Type = v.(string)
	}

	if v, ok := d.GetOkExists("vlan_id"); ok {
		params.VlanID = v.(string)
	}

	if v, ok := d.GetOkExists("zone"); ok {
		params.Zone = v.(string)
	}

	if v, ok := d.GetOk("enable"); ok {
		params.EnableByDefault = v.(bool)
	}

	if v, ok := d.GetOk("ingress_image"); ok {
		params.IngressImage = v.(string)
	}

	if v, ok := d.GetOk("ip"); ok {
		params.IP = v.(string)
	}

	if v, ok := d.GetOk("nlb_version"); ok {
		params.NLBVersion = v.(string)
	}
	var cluster string
	if v, ok := d.GetOkExists("cluster"); ok {
		cluster = v.(string)
	}

	targetEnv, err := getAlbTargetHeader(d, meta)
	if err != nil {
		return err
	}

	//v1.AlbCreateResp
	albResp, err := albAPI.CreateALB(params, cluster, targetEnv)
	if err != nil {
		return fmt.Errorf("Error creating ALb to the cluster %s", err)
	}

	d.SetId(albResp.Alb)
	return nil
}
