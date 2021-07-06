// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"

	"strconv"

	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
)

func resourceIBMLbServiceGroup() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMLbServiceGroupCreate,
		Read:     resourceIBMLbServiceGroupRead,
		Update:   resourceIBMLbServiceGroupUpdate,
		Delete:   resourceIBMLbServiceGroupDelete,
		Exists:   resourceIBMLbServiceGroupExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"virtual_server_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Virtual server ID",
			},
			"service_group_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Service group ID",
			},
			"load_balancer_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Loadbalancer ID",
			},
			"allocation": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Allocation type",
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Port number",
			},
			"routing_method": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Routing method",
			},
			"routing_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Routing type",
			},
			"timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateLBTimeout,
				Description:  "Timeout value",
			},
			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "List of tags",
			},
		},
	}
}

func resourceIBMLbServiceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()

	vipID := d.Get("load_balancer_id").(int)

	routingMethodID, err := getRoutingMethodId(sess, d.Get("routing_method").(string))
	if err != nil {
		return err
	}

	routingTypeID, err := getRoutingTypeId(sess, d.Get("routing_type").(string))
	if err != nil {
		return err
	}

	timeout := d.Get("timeout").(int)

	vip := datatypes.Network_Application_Delivery_Controller_LoadBalancer_VirtualIpAddress{

		VirtualServers: []datatypes.Network_Application_Delivery_Controller_LoadBalancer_VirtualServer{{
			Allocation: sl.Int(d.Get("allocation").(int)),
			Port:       sl.Int(d.Get("port").(int)),
			ServiceGroups: []datatypes.Network_Application_Delivery_Controller_LoadBalancer_Service_Group{{
				RoutingMethodId: &routingMethodID,
				RoutingTypeId:   &routingTypeID,
			}},
		}},
	}

	if timeout > 0 {
		vip.VirtualServers[0].ServiceGroups[0].Timeout = sl.Int(timeout)
	}

	log.Println("[INFO] Creating load balancer service group")

	err = updateLoadBalancerService(sess.SetRetries(0), vipID, &vip)

	if err != nil {
		return fmt.Errorf("Error creating load balancer service group: %s", err)
	}

	// Retrieve the newly created object, to obtain its ID
	vs, err := services.GetNetworkApplicationDeliveryControllerLoadBalancerVirtualIpAddressService(sess).
		Id(vipID).
		Filter(filter.New(filter.Path("virtualServers.port").Eq(d.Get("port"))).Build()).
		Mask("id,serviceGroups[id]").
		GetVirtualServers()

	if err != nil {
		return fmt.Errorf("Error retrieving load balancer: %s", err)
	}

	d.SetId(strconv.Itoa(*vs[0].Id))
	d.Set("service_group_id", vs[0].ServiceGroups[0].Id)

	log.Printf("[INFO] Load Balancer Service Group ID: %s", d.Id())

	return resourceIBMLbServiceGroupRead(d, meta)
}
func resourceIBMLbServiceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()

	vipID := d.Get("load_balancer_id").(int)
	vsID, _ := strconv.Atoi(d.Id())
	sgID := d.Get("service_group_id").(int)

	routingMethodId, err := getRoutingMethodId(sess, d.Get("routing_method").(string))
	if err != nil {
		return err
	}

	routingTypeId, err := getRoutingTypeId(sess, d.Get("routing_type").(string))
	if err != nil {
		return err
	}

	vip := datatypes.Network_Application_Delivery_Controller_LoadBalancer_VirtualIpAddress{

		VirtualServers: []datatypes.Network_Application_Delivery_Controller_LoadBalancer_VirtualServer{{
			Id:         &vsID,
			Allocation: sl.Int(d.Get("allocation").(int)),
			Port:       sl.Int(d.Get("port").(int)),

			ServiceGroups: []datatypes.Network_Application_Delivery_Controller_LoadBalancer_Service_Group{{
				Id:              &sgID,
				RoutingMethodId: &routingMethodId,
				RoutingTypeId:   &routingTypeId,
			}},
		}},
	}

	if d.HasChange("timeout") {
		timeout := d.Get("timeout").(int)
		if timeout > 0 {
			vip.VirtualServers[0].ServiceGroups[0].Timeout = sl.Int(timeout)
		}

	}

	log.Println("[INFO] Updating load balancer service group")

	err = updateLoadBalancerService(sess.SetRetries(0), vipID, &vip)

	if err != nil {
		return fmt.Errorf("Error creating load balancer service group: %s", err)
	}

	return resourceIBMLbServiceGroupRead(d, meta)
}

func resourceIBMLbServiceGroupRead(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()

	vsID, _ := strconv.Atoi(d.Id())

	vs, err := services.GetNetworkApplicationDeliveryControllerLoadBalancerVirtualServerService(sess).
		Id(vsID).
		Mask("allocation,port,serviceGroups[id,routingMethod[keyname],routingType[keyname], timeout],virtualIpAddressId").
		GetObject()

	if err != nil {
		return fmt.Errorf("Error retrieving load balancer: %s", err)
	}

	d.Set("allocation", vs.Allocation)
	d.Set("port", vs.Port)
	d.Set("routing_method", vs.ServiceGroups[0].RoutingMethod.Keyname)
	d.Set("routing_type", vs.ServiceGroups[0].RoutingType.Keyname)
	d.Set("load_balancer_id", vs.VirtualIpAddressId)
	d.Set("service_group_id", vs.ServiceGroups[0].Id)
	d.Set("timeout", vs.ServiceGroups[0].Timeout)

	return nil
}

func resourceIBMLbServiceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()

	vsID, _ := strconv.Atoi(d.Id())

	stateConf := &resource.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"complete"},
		Refresh: func() (interface{}, string, error) {
			err := services.GetNetworkApplicationDeliveryControllerLoadBalancerVirtualServerService(sess).
				Id(vsID).
				DeleteObject()

			if apiErr, ok := err.(sl.Error); ok {
				switch {
				case apiErr.Exception == "SoftLayer_Exception_Network_Timeout" ||
					strings.Contains(apiErr.Message, "There was a problem saving your configuration to the load balancer.") ||
					strings.Contains(apiErr.Message, "The selected group could not be removed from the load balancer.") ||
					strings.Contains(apiErr.Message, "An error has occurred while processing your request.") ||
					strings.Contains(apiErr.Message, "The resource '480' is already in use."):
					// The LB is busy with another transaction. Retry
					return false, "pending", nil
				case apiErr.StatusCode == 404:
					// 404 - service was deleted on the previous attempt
					return true, "complete", nil
				default:
					// Any other error is unexpected. Abort
					return false, "", err
				}
			}

			return true, "complete", nil
		},
		Timeout:    10 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()

	if err != nil {
		return fmt.Errorf("Error deleting service: %s", err)
	}

	return nil
}

func resourceIBMLbServiceGroupExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess := meta.(ClientSession).SoftLayerSession()

	vsID, _ := strconv.Atoi(d.Id())

	_, err := services.GetNetworkApplicationDeliveryControllerLoadBalancerVirtualServerService(sess).
		Id(vsID).
		Mask("id").
		GetObject()

	if err != nil {
		if apiErr, ok := err.(sl.Error); ok {
			if apiErr.StatusCode == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error communicating with the API: %s", err)
	}

	return true, nil
}

func getRoutingTypeId(sess *session.Session, routingTypeName string) (int, error) {
	routingTypes, err := services.GetNetworkApplicationDeliveryControllerLoadBalancerRoutingTypeService(sess).
		Mask("id").
		Filter(filter.Build(
			filter.Path("keyname").Eq(routingTypeName))).
		Limit(1).
		GetAllObjects()

	if err != nil {
		return -1, err
	}

	if len(routingTypes) < 1 {
		return -1, fmt.Errorf("Invalid routing type: %s", routingTypeName)
	}

	return *routingTypes[0].Id, nil
}

func getRoutingMethodId(sess *session.Session, routingMethodName string) (int, error) {
	routingMethods, err := services.GetNetworkApplicationDeliveryControllerLoadBalancerRoutingMethodService(sess).
		Mask("id").
		Filter(filter.Build(
			filter.Path("keyname").Eq(routingMethodName))).
		Limit(1).
		GetAllObjects()

	if err != nil {
		return -1, err
	}

	if len(routingMethods) < 1 {
		return -1, fmt.Errorf("Invalid routing method: %s", routingMethodName)
	}

	return *routingMethods[0].Id, nil
}
