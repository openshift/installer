// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/sl"
)

func resourceIBMLbaasServerInstanceAttachment() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMLbaasServerInstanceAttachmentCreate,
		Read:     resourceIBMLbaasServerInstanceAttachmentRead,
		Delete:   resourceIBMLbaasServerInstanceAttachmentDelete,
		Exists:   resourceIBMLbaasServerInstanceAttachmentExists,
		Update:   resourceIBMLbaasServerInstanceAttachmentUpdate,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"private_ip_address": {
				Type:         schema.TypeString,
				Description:  "The Private IP address of a load balancer member.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateIP,
			},
			"weight": {
				Type:         schema.TypeInt,
				Description:  "The weight of a load balancer member.",
				Computed:     true,
				Optional:     true,
				ValidateFunc: validateWeight,
			},
			"lbaas_id": {
				Type:        schema.TypeString,
				Description: "The UUID of a load balancer",
				ForceNew:    true,
				Required:    true,
			},
			"uuid": {
				Type:        schema.TypeString,
				Description: "The UUID of a load balancer member",
				Computed:    true,
			},
		},
	}
}

func resourceIBMLbaasServerInstanceAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetNetworkLBaaSLoadBalancerService(sess)
	memberService := services.GetNetworkLBaaSMemberService(sess)
	privateIPAddress := d.Get("private_ip_address").(string)
	weight := d.Get("weight").(int)
	lbaasId := d.Get("lbaas_id").(string)
	p := &datatypes.Network_LBaaS_LoadBalancerServerInstanceInfo{}
	p.PrivateIpAddress = sl.String(privateIPAddress)
	p.Weight = sl.Int(weight)
	members := make([]datatypes.Network_LBaaS_LoadBalancerServerInstanceInfo, 0, 1)
	members = append(members, *p)
	_, err := waitForLbaasLBActive(d, meta)
	if err != nil {
		return fmt.Errorf(
			"Error checking for load balancer (%s) is active: %s", d.Id(), err)
	}
	_, err = memberService.AddLoadBalancerMembers(sl.String(lbaasId), members)
	if err != nil {
		return fmt.Errorf("Error adding server instances: %#v", err)
	}
	_, err = waitForLbaasLBActive(d, meta)
	if err != nil {
		return fmt.Errorf(
			"Error waiting for load balancer (%s) to become ready: %s", d.Id(), err)
	}
	result, err := service.Mask("members").GetLoadBalancer(sl.String(lbaasId))
	lbaasMembers := result.Members

	for _, member := range lbaasMembers {
		if *member.Address == privateIPAddress {
			d.SetId(strconv.Itoa(*member.Id))
		}
	}

	return resourceIBMLbaasServerInstanceAttachmentRead(d, meta)
}

func resourceIBMLbaasServerInstanceAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	memberService := services.GetNetworkLBaaSMemberService(sess)
	id := d.Id()
	memId, _ := strconv.Atoi(d.Id())
	member, err := memberService.Id(memId).GetObject()
	if err != nil {
		return fmt.Errorf(
			"Error retrieving load balancer member(%s) : %s", id, err)
	}
	d.Set("private_ip_address", member.Address)
	d.Set("weight", member.Weight)
	d.Set("uuid", member.Uuid)

	return nil
}

func resourceIBMLbaasServerInstanceAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	memberService := services.GetNetworkLBaaSMemberService(sess)
	if d.HasChange("weight") {
		weight := d.Get("weight").(int)
		lbaasId := d.Get("lbaas_id").(string)
		uuid := d.Get("uuid").(string)
		privateIpAddress := d.Get("private_ip_address").(string)

		updateParam := &datatypes.Network_LBaaS_Member{}
		updateParam.Weight = sl.Int(weight)
		updateParam.Uuid = sl.String(uuid)
		updateParam.Address = sl.String(privateIpAddress)
		members := make([]datatypes.Network_LBaaS_Member, 0, 1)
		members = append(members, *updateParam)
		_, err := waitForLbaasLBActive(d, meta)
		if err != nil {
			return fmt.Errorf(
				"Error checking for load balancer (%s) is active: %s", d.Id(), err)
		}
		_, err = memberService.UpdateLoadBalancerMembers(sl.String(lbaasId), members)
		if err != nil {
			return fmt.Errorf("Error updating loadbalnacer: %#v", err)
		}
		_, err = waitForLbaasLBActive(d, meta)
		if err != nil {
			return fmt.Errorf(
				"Error waiting for load balancer (%s) to become ready: %s", lbaasId, err)
		}

	}

	return resourceIBMLbaasServerInstanceAttachmentRead(d, meta)
}

func resourceIBMLbaasServerInstanceAttachmentExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess := meta.(ClientSession).SoftLayerSession()
	memberService := services.GetNetworkLBaaSMemberService(sess)
	memId, _ := strconv.Atoi(d.Id())
	result, err := memberService.Id(memId).GetObject()
	if err != nil {
		if apiErr, ok := err.(sl.Error); ok && (apiErr.StatusCode == 404 || apiErr.Exception == NOT_FOUND) {
			return false, nil
		}
		return false, fmt.Errorf("Error retrieving load balancer member: %s", err)
	}
	return result.Id != nil && *result.Id == memId, nil
}

func resourceIBMLbaasServerInstanceAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	memberService := services.GetNetworkLBaaSMemberService(sess)
	lbaasId := d.Get("lbaas_id").(string)
	removeList := make([]string, 0, 1)
	removeList = append(removeList, d.Get("uuid").(string))
	_, err := waitForLbaasLBActive(d, meta)
	if err != nil {
		return fmt.Errorf(
			"Error checking for load balancer (%s) is active: %s", d.Id(), err)
	}
	_, err = memberService.DeleteLoadBalancerMembers(sl.String(lbaasId), removeList)
	if err != nil {
		return fmt.Errorf("Error removing server instances: %#v", err)
	}
	_, err = waitForLbaasLBActive(d, meta)
	if err != nil {
		return fmt.Errorf(
			"Error waiting for load balancer (%s) to become ready: %s", d.Id(), err)
	}
	return nil
}

func waitForLbaasLBActive(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetNetworkLBaaSLoadBalancerService(sess)
	lbaasId := d.Get("lbaas_id").(string)

	stateConf := &resource.StateChangeConf{
		Pending: []string{lbUpdatePening},
		Target:  []string{lbActive},
		Refresh: func() (interface{}, string, error) {
			lb, err := service.GetLoadBalancer(sl.String(lbaasId))
			if err != nil {
				if apiErr, ok := err.(sl.Error); ok && (apiErr.StatusCode == 404 || apiErr.Exception == NOT_FOUND) {
					return nil, "", fmt.Errorf("The load balancer %s does not exist anymore: %v", d.Id(), err)
				}
				return nil, "", err
			}
			if *lb.ProvisioningStatus == lbActive && *lb.OperatingStatus == lbOnline {
				return lb, lbActive, nil
			}
			return lb, lbUpdatePening, nil
		},
		Timeout:        10 * time.Minute,
		Delay:          60 * time.Second,
		MinTimeout:     3 * time.Second,
		PollInterval:   60 * time.Second,
		NotFoundChecks: 40,
	}

	return stateConf.WaitForState()
}
