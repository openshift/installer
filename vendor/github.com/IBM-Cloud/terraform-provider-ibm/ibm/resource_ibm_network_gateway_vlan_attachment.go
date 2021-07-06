// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/sl"
)

func resourceIBMNetworkGatewayVlanAttachment() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMNetworkGatewayVlanAttachmentCreate,
		Read:     resourceIBMNetworkGatewayVlanAttachmentRead,
		Update:   resourceIBMNetworkGatewayVlanAttachmentUpdate,
		Delete:   resourceIBMNetworkGatewayVlanAttachmentDelete,
		Exists:   resourceIBMNetworkGatewayVlanAttachmentExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Gateway instance ID",
			},
			"network_vlan_id": {
				Type:        schema.TypeInt,
				Description: "The Identifier of the VLAN to be associated",
				Required:    true,
				ForceNew:    true,
			},
			"bypass": {
				Type:        schema.TypeBool,
				Description: "Indicates if the VLAN should be in bypass or routed modes",
				Default:     true,
				Optional:    true,
			},
		},
	}
}

func resourceIBMNetworkGatewayVlanAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	gatewayID := d.Get("gateway_id").(int)
	networkVlanID := d.Get("network_vlan_id").(int)
	bypass := d.Get("bypass").(bool)

	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetNetworkGatewayService(sess)
	vlanService := services.GetNetworkGatewayVlanService(sess)
	result, err := service.Id(gatewayID).Mask(
		"insideVlans",
	).GetObject()
	if err == nil && len(result.InsideVlans) > 0 {
		insideVlans := result.InsideVlans
		for _, i := range insideVlans {
			if *i.NetworkVlanId == networkVlanID {
				if bypass != *i.BypassFlag {
					if !bypass {
						err = vlanService.Id(*i.Id).Unbypass()
						if err != nil {
							return err
						}
					} else {
						err = vlanService.Id(*i.Id).Bypass()
						if err != nil {
							return err
						}
					}
					_, err = waitForNetworkGatewayActiveState(*i.NetworkGatewayId, meta)
					if err != nil {
						return err
					}
				}
				vlan, err := vlanService.Id(*i.Id).GetObject()
				if err != nil {
					return fmt.Errorf("Error trying to retrieve Network Gateway Vlan: %s", err)
				}
				d.SetId(fmt.Sprintf("%d", *vlan.Id))
				d.Set("bypass", vlan.BypassFlag)
				d.Set("gateway_id", vlan.NetworkGatewayId)
				d.Set("network_vlan_id", vlan.NetworkVlanId)
				return nil
			}
		}
	}

	vlan := datatypes.Network_Gateway_Vlan{
		NetworkGatewayId: sl.Int(gatewayID),
		BypassFlag:       sl.Bool(bypass),
		NetworkVlanId:    sl.Int(networkVlanID),
	}

	resp, err := resourceIBMNetworkGatewayVlanAssociate(d, meta, vlan)
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%d", *resp.Id))
	_, err = waitForNetworkGatewayActiveState(gatewayID, meta)
	if err != nil {
		return err
	}
	return resourceIBMNetworkGatewayVlanAttachmentRead(d, meta)
}

func resourceIBMNetworkGatewayVlanAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	vlan, err := services.GetNetworkGatewayVlanService(sess).Id(id).GetObject()
	if err != nil {
		return fmt.Errorf("Error trying to retrieve Network Gateway Vlan: %s", err)
	}
	d.Set("gateway_id", vlan.NetworkGatewayId)
	d.Set("network_vlan_id", vlan.NetworkVlanId)
	d.Set("bypass", vlan.BypassFlag)
	return nil
}

func resourceIBMNetworkGatewayVlanAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetNetworkGatewayVlanService(sess)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}
	if d.HasChange("bypass") {
		bypass := d.Get("bypass").(bool)
		var err error
		if !bypass {
			err = service.Id(id).Unbypass()
			if err != nil {
				return err
			}
		} else {
			err = service.Id(id).Bypass()
			if err != nil {
				return err
			}
		}
		vlan, err := service.Id(id).GetObject()
		_, err = waitForNetworkGatewayActiveState(*vlan.NetworkGatewayId, meta)
		if err != nil {
			return err
		}
	}

	return resourceIBMNetworkGatewayVlanAttachmentRead(d, meta)
}

func resourceIBMNetworkGatewayVlanAttachmentExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	service := services.GetNetworkGatewayVlanService(meta.(ClientSession).SoftLayerSession())
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}
	result, err := service.Id(id).GetObject()
	if err != nil {
		if apiErr, ok := err.(sl.Error); !ok || apiErr.StatusCode != 404 {
			return false, fmt.Errorf("Error trying to retrieve Network Gateway Vlan: %s", err)
		}
	}
	return result.Id != nil && *result.Id == id, nil
}

func resourceIBMNetworkGatewayVlanAttachmentDelete(d *schema.ResourceData, meta interface{}) error {

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}
	vlan, err := services.GetNetworkGatewayVlanService(meta.(ClientSession).SoftLayerSession()).Id(id).GetObject()

	err = resourceIBMNetworkGatewayVlanDissociate(d, meta)
	if err != nil {
		return err
	}

	_, err = waitForNetworkGatewayActiveState(*vlan.NetworkGatewayId, meta)
	if err != nil {
		return err
	}

	return nil

}

func resourceIBMNetworkGatewayVlanAssociate(d *schema.ResourceData, meta interface{}, vlan datatypes.Network_Gateway_Vlan) (resp datatypes.Network_Gateway_Vlan, err error) {
	sess := meta.(ClientSession).SoftLayerSession()

	resp, err = services.GetNetworkGatewayVlanService(sess).CreateObject(&vlan)
	if err != nil {
		return resp, fmt.Errorf(
			"Encountered problem trying to associate the VLAN : %s", err)
	}
	return resp, nil
}

func resourceIBMNetworkGatewayVlanDissociate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}
	err = services.GetNetworkGatewayVlanService(sess).Id(id).DeleteObject()
	if err != nil {
		return fmt.Errorf(
			"Encountered problem trying to dissociate the VLAN : %s", err)
	}
	return nil
}

func waitForNetworkGatewayActiveState(id int, meta interface{}) (interface{}, error) {
	log.Printf("Waiting for Gateway (%d) to be active", id)
	service := services.GetNetworkGatewayService(meta.(ClientSession).SoftLayerSession())

	stateConf := &resource.StateChangeConf{
		Pending: []string{"updating"},
		Target:  []string{"Active"},
		Refresh: func() (interface{}, string, error) {
			gw, err := service.Id(id).GetStatus()
			if err != nil {
				return false, "updating", nil
			}

			if gw.Name != nil && *gw.Name == "Active" {
				log.Printf("Gateway (%d) is active", id)
				return gw, "Active", nil
			}
			log.Printf("Gateway (%d) status is %s", id, *gw.Name)
			return gw, "updating", nil

		},
		Timeout:        24 * time.Hour,
		Delay:          10 * time.Second,
		MinTimeout:     1 * time.Minute,
		NotFoundChecks: 24 * 60,
	}

	return stateConf.WaitForState()
}
