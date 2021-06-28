// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	isVirtualEndpointGatewayID                   = "gateway"
	isVirtualEndpointGatewayIPID                 = "reserved_ip"
	isVirtualEndpointGatewayIPName               = "name"
	isVirtualEndpointGatewayIPAddress            = "address"
	isVirtualEndpointGatewayIPResourceType       = "resource_type"
	isVirtualEndpointGatewayIPAutoDelete         = "auto_delete"
	isVirtualEndpointGatewayIPCreatedAt          = "created_at"
	isVirtualEndpointGatewayIPTarget             = "target"
	isVirtualEndpointGatewayIPTargetID           = "id"
	isVirtualEndpointGatewayIPTargetName         = "name"
	isVirtualEndpointGatewayIPTargetResourceType = "resource_type"
)

func resourceIBMISEndpointGatewayIP() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMisVirtualEndpointGatewayIPCreate,
		Read:     resourceIBMisVirtualEndpointGatewayIPRead,
		Delete:   resourceIBMisVirtualEndpointGatewayIPDelete,
		Exists:   resourceIBMisVirtualEndpointGatewayIPExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			isVirtualEndpointGatewayID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Endpoint gateway ID",
			},
			isVirtualEndpointGatewayIPID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Endpoint gateway IP id",
			},
			isVirtualEndpointGatewayIPName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Endpoint gateway IP name",
			},
			isVirtualEndpointGatewayIPResourceType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Endpoint gateway IP resource type",
			},
			isVirtualEndpointGatewayIPCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Endpoint gateway IP created date and time",
			},
			isVirtualEndpointGatewayIPAutoDelete: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Endpoint gateway IP auto delete",
			},
			isVirtualEndpointGatewayIPAddress: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Endpoint gateway IP address",
			},
			isVirtualEndpointGatewayIPTarget: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Endpoint gateway detail",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isVirtualEndpointGatewayIPTargetID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IPs target id",
						},
						isVirtualEndpointGatewayIPTargetName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IPs target name",
						},
						isVirtualEndpointGatewayIPTargetResourceType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Endpoint gateway resource type",
						},
					},
				},
			},
		},
	}
}

func resourceIBMisVirtualEndpointGatewayIPCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	gatewayID := d.Get(isVirtualEndpointGatewayID).(string)
	ipID := d.Get(isVirtualEndpointGatewayIPID).(string)
	opt := sess.NewAddEndpointGatewayIPOptions(gatewayID, ipID)
	_, response, err := sess.AddEndpointGatewayIP(opt)
	if err != nil {
		log.Printf("Add Endpoint Gateway failed: %v", response)
		return err
	}
	d.SetId(fmt.Sprintf("%s/%s", gatewayID, ipID))
	return resourceIBMisVirtualEndpointGatewayIPRead(d, meta)
}

func resourceIBMisVirtualEndpointGatewayIPRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	gatewayID := parts[0]
	ipID := parts[1]
	opt := sess.NewGetEndpointGatewayIPOptions(gatewayID, ipID)
	result, response, err := sess.GetEndpointGatewayIP(opt)
	if err != nil {
		log.Printf("Get Endpoint Gateway IP failed: %v", response)
		return err
	}
	d.Set(isVirtualEndpointGatewayIPID, result.ID)
	d.Set(isVirtualEndpointGatewayIPName, result.Name)
	d.Set(isVirtualEndpointGatewayIPAddress, result.Address)
	d.Set(isVirtualEndpointGatewayIPCreatedAt, (result.CreatedAt).String())
	d.Set(isVirtualEndpointGatewayIPResourceType, result.ResourceType)
	d.Set(isVirtualEndpointGatewayIPAutoDelete, result.AutoDelete)
	d.Set(isVirtualEndpointGatewayIPTarget,
		flattenEndpointGatewayIPTarget(result.Target.(*vpcv1.ReservedIPTarget)))
	return nil
}

func resourceIBMisVirtualEndpointGatewayIPDelete(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	gatewayID := parts[0]
	ipID := parts[1]
	opt := sess.NewRemoveEndpointGatewayIPOptions(gatewayID, ipID)
	response, err := sess.RemoveEndpointGatewayIP(opt)
	if err != nil && response.StatusCode != 404 {
		log.Printf("Remove Endpoint Gateway IP failed: %v", response)
		return err
	}
	d.SetId("")
	return nil
}

func resourceIBMisVirtualEndpointGatewayIPExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess, err := vpcClient(meta)
	if err != nil {
		return false, err
	}

	parts, err := idParts(d.Id())
	if err != nil {
		return false, err
	}
	if len(parts) != 3 {
		return false, fmt.Errorf("Incorrect ID %s: ID should be a combination of gatewayID/ipID", d.Id())
	}
	gatewayID := parts[0]
	ipID := parts[1]
	opt := sess.NewGetEndpointGatewayIPOptions(gatewayID, ipID)
	_, response, err := sess.GetEndpointGatewayIP(opt)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			log.Printf("Endpoint Gateway IP does not exist.")
			return false, nil
		}
		log.Printf("Error : %s", response)
		return false, err
	}
	return true, nil
}

func flattenEndpointGatewayIPTarget(target *vpcv1.ReservedIPTarget) interface{} {
	targetSlice := []interface{}{}
	targetOutput := map[string]string{}
	if target == nil {
		return targetOutput
	}
	if target.ID != nil {
		targetOutput[isVirtualEndpointGatewayIPTargetID] = *target.ID
	}
	if target.Name != nil {
		targetOutput[isVirtualEndpointGatewayIPTargetName] = *target.Name
	}
	if target.ResourceType != nil {
		targetOutput[isVirtualEndpointGatewayIPTargetResourceType] = *target.ResourceType
	}
	targetSlice = append(targetSlice, targetOutput)
	return targetSlice
}
