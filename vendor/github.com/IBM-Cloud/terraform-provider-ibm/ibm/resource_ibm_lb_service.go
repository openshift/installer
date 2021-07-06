// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
)

func resourceIBMLbService() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMLbServiceCreate,
		Read:     resourceIBMLbServiceRead,
		Update:   resourceIBMLbServiceUpdate,
		Delete:   resourceIBMLbServiceDelete,
		Exists:   resourceIBMLbServiceExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"service_group_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "service group ID",
			},
			"ip_address_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "IP Address ID",
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Port number",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Boolean value true, if enabled else false",
			},
			"health_check_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "health check type",
			},
			"weight": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Weight value",
			},
			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "Tags for the resource",
			},
		},
	}
}

func resourceIBMLbServiceCreate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()

	// SoftLayer Local LBs consist of a multi-level hierarchy of types.
	// (virtualIpAddress -> []virtualServer -> []serviceGroup -> []service)

	// Using the service group ID provided in the config, find the IDs of the
	// respective virtualServer and virtualIpAddress
	sgID := d.Get("service_group_id").(int)
	serviceGroup, err := services.GetNetworkApplicationDeliveryControllerLoadBalancerServiceGroupService(sess).
		Id(sgID).
		Mask("id,routingMethodId,routingTypeId,virtualServer[id,allocation,port,virtualIpAddress[id]]").
		GetObject()

	if err != nil {
		return fmt.Errorf("Error retrieving load balancer service group from SoftLayer, %s", err)
	}

	// Store the IDs for later use
	vsID := *serviceGroup.VirtualServer.Id
	vipID := *serviceGroup.VirtualServer.VirtualIpAddress.Id

	// Convert the health check type name to an ID
	healthCheckTypeId, err := getHealthCheckTypeId(sess, d.Get("health_check_type").(string))
	if err != nil {
		return err
	}

	// The API only exposes edit capability at the root of the tree (virtualIpAddress),
	// so need to send the full structure from the root down to the node to be added or
	// modified
	vip := datatypes.Network_Application_Delivery_Controller_LoadBalancer_VirtualIpAddress{

		VirtualServers: []datatypes.Network_Application_Delivery_Controller_LoadBalancer_VirtualServer{{
			Id:         &vsID,
			Allocation: serviceGroup.VirtualServer.Allocation,
			Port:       serviceGroup.VirtualServer.Port,

			ServiceGroups: []datatypes.Network_Application_Delivery_Controller_LoadBalancer_Service_Group{{
				Id:              &sgID,
				RoutingMethodId: serviceGroup.RoutingMethodId,
				RoutingTypeId:   serviceGroup.RoutingTypeId,

				Services: []datatypes.Network_Application_Delivery_Controller_LoadBalancer_Service{{
					Enabled:     sl.Int(1),
					Port:        sl.Int(d.Get("port").(int)),
					IpAddressId: sl.Int(d.Get("ip_address_id").(int)),

					HealthChecks: []datatypes.Network_Application_Delivery_Controller_LoadBalancer_Health_Check{{
						HealthCheckTypeId: &healthCheckTypeId,
					}},

					GroupReferences: []datatypes.Network_Application_Delivery_Controller_LoadBalancer_Service_Group_CrossReference{{
						Weight: sl.Int(d.Get("weight").(int)),
					}},
				}},
			}},
		}},
	}

	log.Println("[INFO] Creating load balancer service")

	err = updateLoadBalancerService(sess.SetRetries(0), vipID, &vip)

	if err != nil {
		return fmt.Errorf("Error creating load balancer service: %s", err)
	}

	// Retrieve the newly created object, to obtain its ID
	svcs, err := services.GetNetworkApplicationDeliveryControllerLoadBalancerServiceGroupService(sess).
		Id(sgID).
		Mask("mask[id,port,ipAddressId]").
		Filter(filter.New(
			filter.Path("services.port").Eq(d.Get("port")),
			filter.Path("services.ipAddressId").Eq(d.Get("ip_address_id"))).Build()).
		GetServices()

	if err != nil || len(svcs) == 0 {
		return fmt.Errorf("Error retrieving load balancer: %s", err)
	}

	d.SetId(strconv.Itoa(*svcs[0].Id))

	log.Printf("[INFO] Load Balancer Service ID: %s", d.Id())

	return resourceIBMLbServiceRead(d, meta)
}

func resourceIBMLbServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()

	// Using the ID stored in the config, find the IDs of the respective
	// serviceGroup, virtualServer and virtualIpAddress
	svcID, _ := strconv.Atoi(d.Id())
	svc, err := services.GetNetworkApplicationDeliveryControllerLoadBalancerServiceService(sess).
		Id(svcID).
		Mask("id,serviceGroup[id,routingTypeId,routingMethodId,virtualServer[id,allocation,port,virtualIpAddress[id]]]").
		GetObject()

	if err != nil {
		return fmt.Errorf("Error retrieving load balancer service group from SoftLayer, %s", err)
	}

	// Store the IDs for later use
	sgID := *svc.ServiceGroup.Id
	vsID := *svc.ServiceGroup.VirtualServer.Id
	vipID := *svc.ServiceGroup.VirtualServer.VirtualIpAddress.Id

	// Convert the health check type name to an ID
	healthCheckTypeId, err := getHealthCheckTypeId(sess, d.Get("health_check_type").(string))
	if err != nil {
		return err
	}

	// The API only exposes edit capability at the root of the tree (virtualIpAddress),
	// so need to send the full structure from the root down to the node to be added or
	// modified
	vip := datatypes.Network_Application_Delivery_Controller_LoadBalancer_VirtualIpAddress{

		VirtualServers: []datatypes.Network_Application_Delivery_Controller_LoadBalancer_VirtualServer{{
			Id:         &vsID,
			Allocation: svc.ServiceGroup.VirtualServer.Allocation,
			Port:       svc.ServiceGroup.VirtualServer.Port,

			ServiceGroups: []datatypes.Network_Application_Delivery_Controller_LoadBalancer_Service_Group{{
				Id:              &sgID,
				RoutingMethodId: svc.ServiceGroup.RoutingMethodId,
				RoutingTypeId:   svc.ServiceGroup.RoutingTypeId,

				Services: []datatypes.Network_Application_Delivery_Controller_LoadBalancer_Service{{
					Id:          &svcID,
					Enabled:     sl.Int(1),
					Port:        sl.Int(d.Get("port").(int)),
					IpAddressId: sl.Int(d.Get("ip_address_id").(int)),

					HealthChecks: []datatypes.Network_Application_Delivery_Controller_LoadBalancer_Health_Check{{
						HealthCheckTypeId: &healthCheckTypeId,
					}},

					GroupReferences: []datatypes.Network_Application_Delivery_Controller_LoadBalancer_Service_Group_CrossReference{{
						Weight: sl.Int(d.Get("weight").(int)),
					}},
				}},
			}},
		}},
	}

	log.Println("[INFO] Updating load balancer service")

	err = updateLoadBalancerService(sess.SetRetries(0), vipID, &vip)

	if err != nil {
		return fmt.Errorf("Error updating load balancer service: %s", err)
	}

	return resourceIBMLbServiceRead(d, meta)
}

func resourceIBMLbServiceRead(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()

	svcID, _ := strconv.Atoi(d.Id())

	svc, err := services.GetNetworkApplicationDeliveryControllerLoadBalancerServiceService(sess).
		Id(svcID).
		Mask("ipAddressId,enabled,port,healthChecks[type[keyname]],groupReferences[weight],serviceGroup[id]").
		GetObject()

	if err != nil {
		return fmt.Errorf("Error retrieving service: %s", err)
	}

	d.Set("ip_address_id", svc.IpAddressId)
	d.Set("port", svc.Port)
	d.Set("health_check_type", svc.HealthChecks[0].Type.Keyname)
	d.Set("weight", svc.GroupReferences[0].Weight)
	d.Set("enabled", (*svc.Enabled == 1))
	d.Set("service_group_id", svc.ServiceGroup.Id)

	return nil
}

func resourceIBMLbServiceDelete(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()

	svcID, _ := strconv.Atoi(d.Id())

	stateConf := &resource.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"complete"},
		Refresh: func() (interface{}, string, error) {
			err := services.GetNetworkApplicationDeliveryControllerLoadBalancerServiceService(sess).
				Id(svcID).
				DeleteObject()

			if apiErr, ok := err.(sl.Error); ok {
				switch {
				case apiErr.Exception == "SoftLayer_Exception_Network_Timeout" ||
					strings.Contains(apiErr.Message, "There was a problem saving your configuration to the load balancer.") ||
					strings.Contains(apiErr.Message, "The selected group could not be removed from the load balancer.") ||
					strings.Contains(apiErr.Message, "The resource '480' is already in use."):
					// The LB is busy with another transaction. Retry
					return false, "pending", nil
				case apiErr.StatusCode == 404 || // 404 - service was deleted on the previous attempt
					strings.Contains(apiErr.Message, "Unable to find object with id"): // xmlrpc returns 200 instead of 404
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

func resourceIBMLbServiceExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess := meta.(ClientSession).SoftLayerSession()

	svcID, _ := strconv.Atoi(d.Id())

	_, err := services.GetNetworkApplicationDeliveryControllerLoadBalancerServiceService(sess).
		Id(svcID).
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

func getHealthCheckTypeId(sess *session.Session, healthCheckTypeName string) (int, error) {
	healthCheckTypes, err := services.GetNetworkApplicationDeliveryControllerLoadBalancerHealthCheckTypeService(sess).
		Mask("id").
		Filter(filter.Build(
			filter.Path("keyname").Eq(healthCheckTypeName))).
		Limit(1).
		GetAllObjects()

	if err != nil {
		return -1, err
	}

	if len(healthCheckTypes) < 1 {
		return -1, fmt.Errorf("Invalid health check type: %s", healthCheckTypeName)
	}

	return *healthCheckTypes[0].Id, nil
}

func updateLoadBalancerService(sess *session.Session, vipID int, vip *datatypes.Network_Application_Delivery_Controller_LoadBalancer_VirtualIpAddress) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"complete"},
		Refresh: func() (interface{}, string, error) {
			_, err := services.GetNetworkApplicationDeliveryControllerLoadBalancerVirtualIpAddressService(sess).
				Id(vipID).
				EditObject(vip)

			if apiErr, ok := err.(sl.Error); ok {
				// The LB is busy with another transaction. Retry
				if apiErr.Exception == "SoftLayer_Exception_Network_Timeout" ||
					strings.Contains(apiErr.Message, "There was a problem saving your configuration to the load balancer.") ||
					strings.Contains(apiErr.Message, "The selected group could not be removed from the load balancer.") ||
					strings.Contains(apiErr.Message, "The resource '480' is already in use.") {
					return false, "pending", nil
				}

				// Any other error is unexpected. Abort
				return false, "", err
			}

			return true, "complete", nil
		},
		Timeout:    10 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()

	return err
}
