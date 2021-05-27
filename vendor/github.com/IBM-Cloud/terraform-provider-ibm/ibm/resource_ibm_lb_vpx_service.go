// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	dt "github.com/minsikl/netscaler-nitro-go/datatypes"
	"github.com/minsikl/netscaler-nitro-go/op"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/helpers/network"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
)

var (
	// Healthcheck mapping tables
	healthCheckMapFromSLtoVPX105 = map[string]string{
		"HTTP": "http",
		"TCP":  "tcp",
		"ICMP": "ping",
		"icmp": "ping",
		"DNS":  "dns",
	}

	healthCheckMapFromVPX105toSL = map[string]string{
		"http": "HTTP",
		"tcp":  "TCP",
		"ping": "ICMP",
		"dns":  "DNS",
	}
)

func resourceIBMLbVpxService() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMLbVpxServiceCreate,
		Read:     resourceIBMLbVpxServiceRead,
		Update:   resourceIBMLbVpxServiceUpdate,
		Delete:   resourceIBMLbVpxServiceDelete,
		Exists:   resourceIBMLbVpxServiceExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{

			"vip_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VIP id",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "name",
			},

			"destination_ip_address": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Destination IP Address",
			},

			"destination_port": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Destination Port number",
			},

			"weight": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Weight value",
			},

			"connection_limit": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Number of connections limit",
			},

			"health_check": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, o, n string, d *schema.ResourceData) bool {
					if strings.ToUpper(o) == strings.ToUpper(n) {
						return true
					}
					return false
				},
				Description: "Health check info",
			},

			"usip": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "NO",
				Description: "usip info",
			},

			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "list of tags associated with the resource",
			},
		},
	}
}

func parseServiceId(id string) (string, int, string, error) {
	parts := strings.Split(id, ":")
	vipId := parts[1]
	nacdId, err := strconv.Atoi(parts[0])
	if err != nil {
		return "", -1, "", fmt.Errorf("Error parsing vip id: %s", err)
	}

	serviceName := ""
	if len(parts) > 2 {
		serviceName = parts[2]
	}

	return vipId, nacdId, serviceName, nil
}

func updateVpxService(sess *session.Session, nadcId int, lbVip *datatypes.Network_LoadBalancer_VirtualIpAddress) (bool, error) {
	service := services.GetNetworkApplicationDeliveryControllerService(sess)
	serviceName := *lbVip.Services[0].Name
	successFlag := true
	var err error
	for count := 0; count < 10; count++ {
		successFlag, err = service.Id(nadcId).UpdateLiveLoadBalancer(lbVip)
		log.Printf("[INFO] Updating LoadBalancer Service %s successFlag : %t", serviceName, successFlag)

		if err != nil && strings.Contains(err.Error(), "Operation already in progress") {
			log.Printf("[INFO] Updating LoadBalancer Service %s Error : %s. Retry in 10 secs", serviceName, err.Error())
			time.Sleep(time.Second * 10)
			continue
		}

		break
	}
	return successFlag, err
}

func resourceIBMLbVpxServiceCreate(d *schema.ResourceData, meta interface{}) error {
	vipId := d.Get("vip_id").(string)
	_, nadcId, _, err := parseServiceId(vipId)

	if err != nil {
		return fmt.Errorf("Error parsing vip id: %s", err)
	}

	version, err := getVPXVersion(nadcId, meta.(ClientSession).SoftLayerSession())
	if err != nil {
		return fmt.Errorf("Error creating Virtual Ip Address: %s", err)
	}

	if version == VPX_VERSION_10_1 {
		return resourceIBMLbVpxServiceCreate101(d, meta)
	}

	return resourceIBMLbVpxServiceCreate105(d, meta)
}

func resourceIBMLbVpxServiceRead(d *schema.ResourceData, meta interface{}) error {
	_, nadcId, _, err := parseServiceId(d.Id())
	if err != nil {
		return fmt.Errorf("Error parsing vip id: %s", err)
	}

	version, err := getVPXVersion(nadcId, meta.(ClientSession).SoftLayerSession())
	if err != nil {
		return fmt.Errorf("Error Reading Virtual Ip Address: %s", err)
	}

	if version == VPX_VERSION_10_1 {
		return resourceIBMLbVpxServiceRead101(d, meta)
	}

	return resourceIBMLbVpxServiceRead105(d, meta)
}

func resourceIBMLbVpxServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	_, nadcId, _, err := parseServiceId(d.Id())
	if err != nil {
		return fmt.Errorf("Error updating Virtual IP Address: %s", err)
	}

	version, err := getVPXVersion(nadcId, meta.(ClientSession).SoftLayerSession())
	if err != nil {
		return fmt.Errorf("Error updating Virtual Ip Address: %s", err)
	}

	if version == VPX_VERSION_10_1 {
		return resourceIBMLbVpxServiceUpdate101(d, meta)
	}

	return resourceIBMLbVpxServiceUpdate105(d, meta)
}

func resourceIBMLbVpxServiceDelete(d *schema.ResourceData, meta interface{}) error {
	_, nadcId, _, err := parseServiceId(d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting Virtual Ip Address: %s", err)
	}

	version, err := getVPXVersion(nadcId, meta.(ClientSession).SoftLayerSession())
	if err != nil {
		return fmt.Errorf("Error deleting Virtual Ip Address: %s", err)
	}

	if version == VPX_VERSION_10_1 {
		return resourceIBMLbVpxServiceDelete101(d, meta)
	}

	return resourceIBMLbVpxServiceDelete105(d, meta)
}

func resourceIBMLbVpxServiceExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	_, nadcId, _, err := parseServiceId(d.Id())
	if err != nil {
		return false, fmt.Errorf("Error in exists: %s", err)
	}

	version, err := getVPXVersion(nadcId, meta.(ClientSession).SoftLayerSession())
	if err != nil {
		return false, fmt.Errorf("Error in exists: %s", err)
	}

	if version == VPX_VERSION_10_1 {
		return resourceIBMLbVpxServiceExists101(d, meta)
	}

	return resourceIBMLbVpxServiceExists105(d, meta)
}

func resourceIBMLbVpxServiceCreate101(d *schema.ResourceData, meta interface{}) error {

	sess := meta.(ClientSession).SoftLayerSession()

	vipId := d.Get("vip_id").(string)
	vipName, nadcId, _, err := parseServiceId(vipId)
	serviceName := d.Get("name").(string)

	if err != nil {
		return fmt.Errorf("Error parsing vip id: %s", err)
	}

	lb_services := []datatypes.Network_LoadBalancer_Service{
		{
			Name:                 sl.String(d.Get("name").(string)),
			DestinationIpAddress: sl.String(d.Get("destination_ip_address").(string)),
			DestinationPort:      sl.Int(d.Get("destination_port").(int)),
			Weight:               sl.Int(d.Get("weight").(int)),
			HealthCheck:          sl.String(d.Get("health_check").(string)),
			ConnectionLimit:      sl.Int(d.Get("connection_limit").(int)),
		},
	}

	lbVip := &datatypes.Network_LoadBalancer_VirtualIpAddress{
		Name:     sl.String(vipName),
		Services: lb_services,
	}

	// Check if there is an existed loadbalancer service which has same name.
	log.Printf("[INFO] Creating LoadBalancer Service Name %s validation", serviceName)

	_, err = network.GetNadcLbVipServiceByName(sess, nadcId, vipName, serviceName)
	if err == nil {
		return fmt.Errorf("Error creating LoadBalancer Service: The service name '%s' is already used.",
			serviceName)
	}

	log.Printf("[INFO] Creating LoadBalancer Service %s", serviceName)

	successFlag, err := updateVpxService(sess.SetRetries(0), nadcId, lbVip)

	if err != nil {
		return fmt.Errorf("Error creating LoadBalancer Service: %s", err)
	}

	if !successFlag {
		return errors.New("Error creating LoadBalancer Service")
	}

	d.SetId(fmt.Sprintf("%s:%s", vipId, serviceName))

	return resourceIBMLbVpxServiceRead(d, meta)
}

func resourceIBMLbVpxServiceCreate105(d *schema.ResourceData, meta interface{}) error {

	sess := meta.(ClientSession).SoftLayerSession()

	vipId := d.Get("vip_id").(string)
	vipName, nadcId, _, err := parseServiceId(vipId)
	serviceName := d.Get("name").(string)

	if err != nil {
		return fmt.Errorf("Error parsing vip id: %s", err)
	}

	nClient, err := getNitroClient(sess, nadcId)
	if err != nil {
		return fmt.Errorf("Error getting netscaler information ID: %d", nadcId)
	}

	// Create a service
	svcReq := dt.ServiceReq{
		Service: &dt.Service{
			Name:      op.String(d.Get("name").(string)),
			Ip:        op.String(d.Get("destination_ip_address").(string)),
			Port:      op.Int(d.Get("destination_port").(int)),
			Maxclient: op.String(strconv.Itoa(d.Get("connection_limit").(int))),
			Usip:      op.String(d.Get("usip").(string)),
		},
	}

	// Get serviceType of a virtual server
	vip := dt.LbvserverRes{}
	err = nClient.Get(&vip, vipName)
	if err != nil {
		return fmt.Errorf("Error creating LoadBalancer Service : %s", err)
	}

	if vip.Lbvserver[0].ServiceType != nil {
		svcReq.Service.ServiceType = vip.Lbvserver[0].ServiceType
	} else {
		return fmt.Errorf("Error creating LoadBalancer : type of VIP '%s' is null.", vipName)
	}

	// SSL offload
	if *svcReq.Service.ServiceType == "SSL" {
		*svcReq.Service.ServiceType = "HTTP"
	}

	log.Printf("[INFO] Creating LoadBalancer Service %s", serviceName)

	// Add the service
	err = nClient.Add(&svcReq)
	if err != nil {
		return fmt.Errorf("Error creating LoadBalancer Service: %s", err)
	}

	// Bind the virtual server and the service
	lbvserverServiceBindingReq := dt.LbvserverServiceBindingReq{
		LbvserverServiceBinding: &dt.LbvserverServiceBinding{
			Name:        op.String(vipName),
			ServiceName: op.String(serviceName),
		},
	}

	err = nClient.Add(&lbvserverServiceBindingReq)
	if err != nil {
		return fmt.Errorf("Error creating LoadBalancer Service: %s", err)
	}

	// Bind Health_check monitor
	healthCheck := d.Get("health_check").(string)
	if len(healthCheckMapFromSLtoVPX105[healthCheck]) > 0 {
		healthCheck = healthCheckMapFromSLtoVPX105[healthCheck]
	}

	serviceLbmonitorBindingReq := dt.ServiceLbmonitorBindingReq{
		ServiceLbmonitorBinding: &dt.ServiceLbmonitorBinding{
			Name:        op.String(serviceName),
			MonitorName: op.String(healthCheck),
		},
	}

	err = nClient.Add(&serviceLbmonitorBindingReq)
	if err != nil {
		return fmt.Errorf("Error creating LoadBalancer Service: %s", err)
	}

	d.SetId(fmt.Sprintf("%s:%s", vipId, serviceName))

	return resourceIBMLbVpxServiceRead(d, meta)
}

func resourceIBMLbVpxServiceRead101(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()

	vipName, nadcId, serviceName, err := parseServiceId(d.Id())
	if err != nil {
		return fmt.Errorf("Error parsing vip id: %s", err)
	}

	lbService, err := network.GetNadcLbVipServiceByName(sess, nadcId, vipName, serviceName)
	if err != nil {
		return fmt.Errorf("Unable to get load balancer service %s: %s", serviceName, err)
	}

	d.Set("vip_id", strconv.Itoa(nadcId)+":"+vipName)
	d.Set("name", *lbService.Name)
	d.Set("destination_ip_address", *lbService.DestinationIpAddress)
	d.Set("destination_port", *lbService.DestinationPort)
	d.Set("weight", *lbService.Weight)
	d.Set("health_check", *lbService.HealthCheck)
	d.Set("connection_limit", *lbService.ConnectionLimit)

	return nil
}

func resourceIBMLbVpxServiceRead105(d *schema.ResourceData, meta interface{}) error {
	vipName, nadcId, serviceName, err := parseServiceId(d.Id())
	if err != nil {
		return fmt.Errorf("Error parsing vip id: %s", err)
	}

	nClient, err := getNitroClient(meta.(ClientSession).SoftLayerSession(), nadcId)
	if err != nil {
		return fmt.Errorf("Error getting netscaler information ID: %d", nadcId)
	}

	// Read a service

	svc := dt.ServiceRes{}
	err = nClient.Get(&svc, serviceName)
	if err != nil {
		fmt.Printf("Error getting service information : %s", err.Error())
	}
	d.Set("vip_id", strconv.Itoa(nadcId)+":"+vipName)
	d.Set("name", *svc.Service[0].Name)
	d.Set("destination_ip_address", *svc.Service[0].Ipaddress)
	d.Set("destination_port", *svc.Service[0].Port)
	d.Set("usip", *svc.Service[0].Usip)

	maxClientStr, err := strconv.Atoi(*svc.Service[0].Maxclient)
	if err == nil {
		d.Set("connection_limit", maxClientStr)
	}

	// Read a monitor information
	healthCheck := dt.ServiceLbmonitorBindingRes{}
	err = nClient.Get(&healthCheck, serviceName)
	if err != nil {
		fmt.Printf("Error getting service information : %s", err.Error())
	}
	if healthCheck.ServiceLbmonitorBinding[0].MonitorName != nil {
		healthCheck := *healthCheck.ServiceLbmonitorBinding[0].MonitorName
		if len(healthCheckMapFromVPX105toSL[healthCheck]) > 0 {
			healthCheck = healthCheckMapFromVPX105toSL[healthCheck]
		}
		d.Set("health_check", healthCheck)
	}

	return nil
}

func resourceIBMLbVpxServiceUpdate101(d *schema.ResourceData, meta interface{}) error {

	sess := meta.(ClientSession).SoftLayerSession()

	vipName, nadcId, serviceName, err := parseServiceId(d.Id())
	if err != nil {
		return fmt.Errorf("Error parsing vip id: %s", err)
	}

	lbService, err := network.GetNadcLbVipServiceByName(sess, nadcId, vipName, serviceName)
	if err != nil {
		return fmt.Errorf("Unable to get load balancer service: %s", err)
	}

	// copy current service
	template := datatypes.Network_LoadBalancer_Service(*lbService)

	if data, ok := d.GetOk("name"); ok {
		template.Name = sl.String(data.(string))
	}
	if data, ok := d.GetOk("destination_ip_address"); ok {
		template.DestinationIpAddress = sl.String(data.(string))
	}
	if data, ok := d.GetOk("destination_port"); ok {
		template.DestinationPort = sl.Int(data.(int))
	}
	if data, ok := d.GetOk("weight"); ok {
		template.Weight = sl.Int(data.(int))
	}
	if data, ok := d.GetOk("health_check"); ok {
		template.HealthCheck = sl.String(data.(string))
	}
	if data, ok := d.GetOk("connection_limit"); ok {
		template.ConnectionLimit = sl.Int(data.(int))
	}

	lbVip := &datatypes.Network_LoadBalancer_VirtualIpAddress{
		Name: sl.String(vipName),
		Services: []datatypes.Network_LoadBalancer_Service{
			template},
	}

	successFlag, err := updateVpxService(sess.SetRetries(0), nadcId, lbVip)

	if err != nil {
		return fmt.Errorf("Error updating LoadBalancer Service: %s", err)
	}

	if !successFlag {
		return errors.New("Error updating LoadBalancer Service")
	}

	return nil
}

func resourceIBMLbVpxServiceUpdate105(d *schema.ResourceData, meta interface{}) error {
	_, nadcId, serviceName, err := parseServiceId(d.Id())
	if err != nil {
		return fmt.Errorf("Error parsing vip id: %s", err)
	}

	nClient, err := getNitroClient(meta.(ClientSession).SoftLayerSession(), nadcId)
	if err != nil {
		return fmt.Errorf("Error getting netscaler information ID: %d", nadcId)
	}

	// Update a service
	svcReq := dt.ServiceReq{
		Service: &dt.Service{
			Name: op.String(d.Get("name").(string)),
		},
	}

	updateFlag := false

	if d.HasChange("health_check") {
		healthCheck := dt.ServiceLbmonitorBindingRes{}
		err = nClient.Get(&healthCheck, serviceName)
		if err != nil {
			fmt.Printf("Error getting service information : %s", err.Error())
		}
		monitorName := healthCheck.ServiceLbmonitorBinding[0].MonitorName
		if monitorName != nil && *monitorName != "tcp-default" {
			// Delete previous health_check
			err = nClient.Delete(&dt.ServiceLbmonitorBindingReq{}, serviceName, "args=monitor_name:"+*monitorName)
			if err != nil {
				return fmt.Errorf("Error deleting monitor %s: %s", *monitorName, err)
			}
		}

		// Add a new health_check
		monitor := d.Get("health_check").(string)
		if len(healthCheckMapFromSLtoVPX105[monitor]) > 0 {
			monitor = healthCheckMapFromSLtoVPX105[monitor]
		}

		serviceLbmonitorBindingReq := dt.ServiceLbmonitorBindingReq{
			ServiceLbmonitorBinding: &dt.ServiceLbmonitorBinding{
				Name:        op.String(serviceName),
				MonitorName: op.String(monitor),
			},
		}

		err = nClient.Add(&serviceLbmonitorBindingReq)
		if err != nil {
			return fmt.Errorf("Error adding a monitor: %s", err)
		}
	}

	if d.HasChange("connection_limit") {
		svcReq.Service.Maxclient = op.String(strconv.Itoa(d.Get("connection_limit").(int)))
		updateFlag = true
	}

	if d.HasChange("usip") {
		svcReq.Service.Usip = op.String(d.Get("usip").(string))
		updateFlag = true
	}

	log.Printf("[INFO] Updating LoadBalancer Service %s", serviceName)

	if updateFlag {
		err = nClient.Update(&svcReq)
	}

	if err != nil {
		return fmt.Errorf("Error updating LoadBalancer Service: %s", err)
	}

	return nil
}

func resourceIBMLbVpxServiceDelete101(d *schema.ResourceData, meta interface{}) error {
	vipName, nadcId, serviceName, err := parseServiceId(d.Id())
	if err != nil {
		return fmt.Errorf("Error parsing vip id: %s", err)
	}

	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetNetworkApplicationDeliveryControllerService(sess)

	lbSvc := datatypes.Network_LoadBalancer_Service{
		Name: sl.String(serviceName),
		Vip: &datatypes.Network_LoadBalancer_VirtualIpAddress{
			Name: sl.String(vipName),
		},
	}

	for count := 0; count < 10; count++ {
		err = service.Id(nadcId).DeleteLiveLoadBalancerService(&lbSvc)
		log.Printf("[INFO] Deleting Loadbalancer service %s", serviceName)

		if err != nil &&
			(strings.Contains(err.Error(), "Operation already in progress") ||
				strings.Contains(err.Error(), "Internal Error")) {
			log.Printf("[INFO] Deleting Loadbalancer service Error : %s. Retry in 10 secs", err.Error())
			time.Sleep(time.Second * 10)
			continue
		}

		if err != nil &&
			(strings.Contains(err.Error(), "No Service") ||
				strings.Contains(err.Error(), "Unable to find object with unknown identifier of")) {
			log.Printf("[INFO] Deleting Loadbalancer service %s Error : %s ", serviceName, err.Error())
			err = nil
		}

		break
	}

	if err != nil {
		return fmt.Errorf("Error deleting LoadBalancer Service %s: %s", serviceName, err)
	}

	return nil
}

func resourceIBMLbVpxServiceDelete105(d *schema.ResourceData, meta interface{}) error {
	_, nadcId, serviceName, err := parseServiceId(d.Id())
	if err != nil {
		return fmt.Errorf("Error parsing vip id: %s", err)
	}

	nClient, err := getNitroClient(meta.(ClientSession).SoftLayerSession(), nadcId)
	if err != nil {
		return fmt.Errorf("Error getting netscaler information ID: %d", nadcId)
	}

	// Delete a service
	err = nClient.Delete(&dt.ServiceReq{}, serviceName)
	if err != nil {
		return fmt.Errorf("Error deleting service %s: %s", serviceName, err)
	}

	return nil
}

func resourceIBMLbVpxServiceExists101(d *schema.ResourceData, meta interface{}) (bool, error) {
	vipName, nadcId, serviceName, err := parseServiceId(d.Id())
	if err != nil {
		return false, fmt.Errorf("Error parsing vip id: %s", err)
	}
	lbService, err := network.GetNadcLbVipServiceByName(meta.(ClientSession).SoftLayerSession(), nadcId, vipName, serviceName)
	if err != nil {
		if apiErr, ok := err.(sl.Error); ok {
			if apiErr.StatusCode == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error communicating with the API: %s", err)
	}
	return *lbService.Name == serviceName, nil
}

func resourceIBMLbVpxServiceExists105(d *schema.ResourceData, meta interface{}) (bool, error) {
	_, nadcId, serviceName, err := parseServiceId(d.Id())
	if err != nil {
		return false, fmt.Errorf("Error parsing vip id: %s", err)
	}

	nClient, err := getNitroClient(meta.(ClientSession).SoftLayerSession(), nadcId)
	if err != nil {
		return false, fmt.Errorf("Error getting netscaler information ID: %d", nadcId)
	}

	svc := dt.ServiceRes{}
	err = nClient.Get(&svc, serviceName)
	if err != nil && strings.Contains(err.Error(), "No Service") {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("Unable to get load balancer service %s: %s", serviceName, err)
	}

	return *svc.Service[0].Name == serviceName, nil
}
