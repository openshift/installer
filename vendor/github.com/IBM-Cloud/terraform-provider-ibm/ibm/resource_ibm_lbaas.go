// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/internal/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/helpers/product"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
)

var packageType = "LOAD_BALANCER_AS_A_SERVICE"
var lbMethodToId = make(map[string]string)
var lbIdToMethod = make(map[string]string)

const (
	lbActive        = "ACTIVE"
	lbPending       = "CREATE_PENDING"
	lbUpdatePening  = "UPDATE_PENDING"
	lbOnline        = "ONLINE"
	lbDeletePending = "DELETE_PENDING"
	lbDeleted       = "DELETED"
)

const NOT_FOUND = "SoftLayer_Exception_Network_LBaaS_ObjectNotFound"

const productItemMaskWithPriceLocationGroupID = "id,categories,capacity,description,units,keyName,prices[id,categories[id,name,categoryCode],locationGroupId,capacityRestrictionMaximum,capacityRestrictionMinimum,capacityRestrictionType,bareMetalReservedCapacityFlag],totalPhysicalCoreCapacity,totalProcessorCapacity"

func init() {

	lbMethodToId = map[string]string{
		"round_robin":          "ROUNDROBIN",
		"weighted_round_robin": "WEIGHTED_RR",
		"least_connection":     "LEASTCONNECTION",
	}
	for k, v := range lbMethodToId {
		lbIdToMethod[v] = k
	}
}

func resourceIBMLbaas() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMLbaasCreate,
		Read:     resourceIBMLbaasRead,
		Delete:   resourceIBMLbaasDelete,
		Exists:   resourceIBMLbaasExists,
		Update:   resourceIBMLbaasUpdate,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The load balancer's name.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of a load balancer.",
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "PUBLIC",
				ForceNew:     true,
				Description:  "Specifies if a load balancer is public or private",
				ValidateFunc: validateAllowedStringValue([]string{"PUBLIC", "PRIVATE"}),
			},
			"datacenter": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnets": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "The subnet where this Load Balancer will be provisioned.",
				Elem:        &schema.Schema{Type: schema.TypeInt},
				MinItems:    1,
				MaxItems:    1,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "The operation status 'ONLINE' or 'OFFLINE' of a load balancer.",
				Computed:    true,
			},
			"vip": {
				Type:        schema.TypeString,
				Description: "The virtual ip address of this load balancer",
				Computed:    true,
			},
			"use_system_public_ip_pool": {
				Type:             schema.TypeBool,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: applyOnce,
				Description:      `"in public loadbalancer - Public IP address allocation done by system public IP pool or public subnet."`,
			},
			"protocols": {
				Type:        schema.TypeSet,
				Description: "Protocols to be assigned to this load balancer.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"frontend_protocol": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Frontend protocol, one of 'TCP', 'HTTP', 'HTTPS'.",
							ValidateFunc: validateAllowedStringValue([]string{"HTTP", "HTTPS", "TCP"}),
						},
						"frontend_port": {
							Type:         schema.TypeInt,
							Required:     true,
							Description:  "Frontend Protocol port number. Should be in range (1, 65535)",
							ValidateFunc: validatePortRange(1, 65535),
						},
						"backend_protocol": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Backend protocol, one of 'TCP', 'HTTP', 'HTTPS'.",
							ValidateFunc: validateAllowedStringValue([]string{"HTTP", "HTTPS", "TCP"}),
						},
						"backend_port": {
							Type:         schema.TypeInt,
							Required:     true,
							Description:  "Backend Protocol port number. Should be in range (1, 65535)",
							ValidateFunc: validatePortRange(1, 65535),
						},
						"load_balancing_method": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateAllowedStringValue([]string{"round_robin", "weighted_round_robin", "least_connection"}),
							Default:      "round_robin",
							Description:  "Load balancing algorithm: 'round_robin', 'weighted_round_robin', 'least_connection'",
						},
						"session_stickiness": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Session stickness. Valid values is SOURCE_IP and HTTP_COOKIE",
							ValidateFunc: validateAllowedStringValue([]string{"SOURCE_IP", "HTTP_COOKIE"}),
						},
						"max_conn": {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "No. of connections the listener can accept. Should be between 1-64000",
							ValidateFunc: validateMaxConn,
						},
						"tls_certificate_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "This references to SSL/TLS certificate for a protocol",
						},
						"protocol_id": {
							Type:        schema.TypeString,
							Description: "The UUID of a load balancer protocol",
							Computed:    true,
						},
					},
				},
				Set: resourceIBMLBProtocolHash,
			},
			"ssl_ciphers": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				//ValidateFunc: validateAllowedStringValue([]string{"ECDHE-RSA-AES256-GCM-SHA384", "ECDHE-RSA-AES256-SHA384", "AES256-GCM-SHA384", "AES256-SHA256", "ECDHE-RSA-AES128-GCM-SHA256", "ECDHE-RSA-AES128-SHA256", "AES128-GCM-SHA256", "AES128-SHA256"}),
			},
			"wait_time_minutes": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  90,
			},
			"health_monitors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_retries": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"url_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"monitor_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			ResourceControllerURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the IBM Cloud dashboard that can be used to explore and view details about this instance",
			},
			ResourceName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the resource",
			},
			ResourceStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the resource",
			},
		},
	}
}

func resourceIBMLbaasCreate(d *schema.ResourceData, meta interface{}) error {

	sess := meta.(ClientSession).SoftLayerSession()

	// Find price items
	productOrderContainer, err := buildLbaasLBProductOrderContainer(d, sess)
	if err != nil {
		return fmt.Errorf("Error creating Load balancer: %s", err)
	}
	log.Println("[INFO] Creating Load Balancer")

	//verify order
	_, err = services.GetProductOrderService(sess).
		VerifyOrder(productOrderContainer)
	if err != nil {
		return fmt.Errorf("Error during creation of Load balancer: %s", err)
	}
	//place order
	_, err = services.GetProductOrderService(sess.SetRetries(0)).
		PlaceOrder(productOrderContainer, sl.Bool(false))
	if err != nil {
		return fmt.Errorf("Error during creation of Load balancer: %s", err)
	}

	name := d.Get("name").(string)

	lbaasLB, err := findLbaasLBByOrderId(sess, name, d)
	if err != nil {
		return fmt.Errorf("Error during creation of Load balancer: %s", err)
	}

	d.SetId(*lbaasLB.Uuid)

	return resourceIBMLbaasUpdate(d, meta)
}

func resourceIBMLbaasRead(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetNetworkLBaaSLoadBalancerService(sess)

	result, err := service.Mask("datacenter,members,listeners.defaultPool,listeners.defaultPool.sessionAffinity,listeners.defaultPool.healthMonitor,healthMonitors,sslCiphers[name],useSystemPublicIpPool,isPublic,name,description,operatingStatus,address").GetLoadBalancer(sl.String(d.Id()))
	if err != nil {
		return fmt.Errorf("Error retrieving load balancer: %s", err)
	}

	var lbType string
	if *result.IsPublic == 1 {
		lbType = "PUBLIC"
	} else {
		lbType = "PRIVATE"
	}
	//TODO THis is public subnet and we need to set the private subnet
	//subnets := [1]int{*result.IpAddress.SubnetId}
	//d.Set("subnets", subnets)
	d.Set("name", result.Name)
	d.Set("description", result.Description)
	d.Set("datacenter", result.Datacenter.Name)
	d.Set("type", lbType)
	d.Set("status", result.OperatingStatus)
	d.Set("vip", result.Address)
	d.Set("health_monitors", flattenHealthMonitors(result.Listeners))
	d.Set("protocols", flattenProtocols(result.Listeners))
	d.Set("ssl_ciphers", flattenSSLCiphers(result.SslCiphers))
	if *result.UseSystemPublicIpPool == 1 {
		d.Set("use_system_public_ip_pool", true)
	} else {
		d.Set("use_system_public_ip_pool", false)
	}
	d.Set(ResourceControllerURL, fmt.Sprintf("https://cloud.ibm.com/classic/network/loadbalancing/cloud/details/%s#Overview", d.Id()))
	d.Set(ResourceName, *result.Name)
	d.Set(ResourceStatus, *result.OperatingStatus)

	return nil
}

func resourceIBMLbaasUpdate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetNetworkLBaaSLoadBalancerService(sess.SetRetries(0))

	if d.HasChange("description") {
		_, err := service.UpdateLoadBalancer(sl.String(d.Id()), sl.String(d.Get("description").(string)))
		if err != nil {
			return err
		}
	}
	listenerService := services.GetNetworkLBaaSListenerService(sess.SetRetries(0))
	if d.HasChange("protocols") {
		o, n := d.GetChange("protocols")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		add, err := expandProtocols(ns.Difference(os).List())
		if err != nil {
			return err
		}
		rem := os.Difference(ns).List()
		removeList := make([]string, len(rem), len(rem))
		for i, remove := range rem {
			data := remove.(map[string]interface{})
			if v, ok := data["protocol_id"]; ok && v.(string) != "" {
				removeList[i] = v.(string)
			}
		}
		if len(removeList) > 0 {
			_, err := listenerService.DeleteLoadBalancerProtocols(sl.String(d.Id()), removeList)
			if err != nil {
				return fmt.Errorf("Error removing protocols: %#v", err)
			}
			_, err = waitForLbaasLBAvailable(d, meta)
			if err != nil {
				return fmt.Errorf(
					"Error waiting for load balancer (%s) to become ready: %s", d.Id(), err)
			}
		}

		if len(add) > 0 {
			_, err := listenerService.UpdateLoadBalancerProtocols(sl.String(d.Id()), add)
			if err != nil {
				return fmt.Errorf("Error adding protocols: %#v", err)
			}
			_, err = waitForLbaasLBAvailable(d, meta)
			if err != nil {
				return fmt.Errorf(
					"Error waiting for load balancer (%s) to become ready: %s", d.Id(), err)
			}
		}

	}
	if d.HasChange("ssl_ciphers") {
		if v, ok := d.GetOk("ssl_ciphers"); ok && v.(*schema.Set).Len() > 0 {
			service := services.GetNetworkLBaaSLoadBalancerService(sess.SetRetries(0))
			supportedCiphers, err := services.GetNetworkLBaaSSSLCipherService(sess).Mask("id,name").GetAllObjects()
			if err != nil {
				return fmt.Errorf("Error retreving list of ssl ciphers: %#v", err)
			}
			ciphers := make([]int, v.(*schema.Set).Len())
			for i, v := range v.(*schema.Set).List() {
				for _, c := range supportedCiphers {
					if v == *c.Name {
						ciphers[i] = *c.Id
						break
					}
				}
			}
			_, err = service.UpdateSslCiphers(sl.String(d.Id()), ciphers)
			if err != nil {
				return fmt.Errorf("Error updating ssl ciphers: %#v", err)
			}
			_, err = waitForLbaasLBAvailable(d, meta)
			if err != nil {
				return fmt.Errorf(
					"Error waiting for load balancer (%s) to become ready: %s", d.Id(), err)
			}

		}

	}

	return resourceIBMLbaasRead(d, meta)
}

func resourceIBMLbaasDelete(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetNetworkLBaaSLoadBalancerService(sess)

	_, err := service.CancelLoadBalancer(sl.String(d.Id()))
	if err != nil {
		if strings.Contains(err.Error(), "DELETE_PENDING") {
			log.Println("Deletion is already in progress, probably from previous runs")
		} else {
			return fmt.Errorf("Error deleting load balancer: %s", err)
		}
	}
	_, err = waitForLbaasLBDelete(d, meta)
	if err != nil {
		return fmt.Errorf(
			"Error waiting for load balancer (%s) to be deleted: %s", d.Id(), err)
	}
	d.SetId("")
	return nil
}

func resourceIBMLbaasExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetNetworkLBaaSLoadBalancerService(sess)

	result, err := service.GetLoadBalancer(sl.String(d.Id()))
	if err != nil {
		if apiErr, ok := err.(sl.Error); ok && (apiErr.StatusCode == 404 || apiErr.Exception == NOT_FOUND) {
			return false, nil
		}
		return false, fmt.Errorf("Error retrieving load balancer: %s", err)
	}
	return result.Uuid != nil && *result.Uuid == d.Id(), nil
}

func buildLbaasLBProductOrderContainer(d *schema.ResourceData, sess *session.Session) (*datatypes.Container_Product_Order_Network_LoadBalancer_AsAService, error) {
	// 1. Get a package
	name := d.Get("name").(string)
	subnets := d.Get("subnets").([]interface{})
	lbType := d.Get("type").(string)

	subnetsParam := []datatypes.Network_Subnet{}
	for _, subnet := range subnets {
		subnetItem := datatypes.Network_Subnet{
			Id: sl.Int(subnet.(int)),
		}
		subnetsParam = append(subnetsParam, subnetItem)
	}

	pkg, err := product.GetPackageByType(sess, packageType)
	if err != nil {
		return nil, err
	}

	// 2. Get all prices for the package
	productItems, err := product.GetPackageProducts(sess, *pkg.Id, productItemMaskWithPriceLocationGroupID)
	if err != nil {
		return &datatypes.Container_Product_Order_Network_LoadBalancer_AsAService{}, err
	}

	priceItems := []datatypes.Product_Item_Price{}
	for _, item := range productItems {
		for _, price := range item.Prices {
			if price.LocationGroupId == nil && !*price.BareMetalReservedCapacityFlag {
				priceItem := datatypes.Product_Item_Price{
					Id: price.Id,
				}
				priceItems = append(priceItems, priceItem)
				break
			}
		}
	}

	productOrderContainer := datatypes.Container_Product_Order_Network_LoadBalancer_AsAService{
		Container_Product_Order: datatypes.Container_Product_Order{
			PackageId:        pkg.Id,
			Prices:           priceItems,
			Quantity:         sl.Int(1),
			UseHourlyPricing: sl.Bool(true),
		},
		Name:    sl.String(name),
		Subnets: subnetsParam,
	}
	if d, ok := d.GetOk("description"); ok {
		productOrderContainer.Description = sl.String(d.(string))
	}

	if lbType == "PRIVATE" {
		productOrderContainer.IsPublic = sl.Bool(false)
	}
	if publicIPPool, ok := d.GetOkExists("use_system_public_ip_pool"); ok {
		productOrderContainer.UseSystemPublicIpPool = sl.Bool(publicIPPool.(bool))
	}

	return &productOrderContainer, nil
}

func findLbaasLBByOrderId(sess *session.Session, name string, d *schema.ResourceData) (*datatypes.Network_LBaaS_LoadBalancer, error) {

	isIDSet := false
	stateConf := &resource.StateChangeConf{
		Pending: []string{lbPending},
		Target:  []string{lbActive},
		Refresh: func() (interface{}, string, error) {
			/*lb, err := services.GetAccountService(sess).
			Filter(filter.Path("loadbalancer.billingItem.orderItem.order.id").
				Eq(strconv.Itoa(orderId)).Build()).
			Mask("id,activeTransaction").
			GetLoadBalancer()*/
			//TODO This is a temporary workaround to find lbass obj by name.Get the lbass obj from order id
			lb, err := services.GetNetworkLBaaSLoadBalancerService(sess).Filter(filter.Build(
				filter.Path("name").Eq(name))).GetAllObjects()
			if err != nil {
				return nil, "", err
			}
			if len(lb) == 1 {
				if *lb[0].ProvisioningStatus == lbActive && *lb[0].OperatingStatus == lbOnline {
					return lb[0], lbActive, nil
				}
				if !isIDSet && lb[0].Uuid != nil {
					d.SetId(*lb[0].Uuid)
					isIDSet = true
				}
				return lb[0], lbPending, nil
			}
			return nil, lbPending, nil
		},
		Timeout:        time.Duration(d.Get("wait_time_minutes").(int)) * time.Minute,
		Delay:          60 * time.Second,
		MinTimeout:     3 * time.Second,
		PollInterval:   60 * time.Second,
		NotFoundChecks: 40,
	}

	pendingResult, err := stateConf.WaitForState()

	if err != nil {
		return nil, err
	}

	if result, ok := pendingResult.(datatypes.Network_LBaaS_LoadBalancer); ok {
		return &result, nil
	}

	return nil,
		fmt.Errorf("Cannot find a load balancer with name '%s' ", name)
}

func waitForLbaasLBAvailable(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetNetworkLBaaSLoadBalancerService(sess)

	stateConf := &resource.StateChangeConf{
		Pending: []string{lbUpdatePening},
		Target:  []string{lbActive},
		Refresh: func() (interface{}, string, error) {
			lb, err := service.GetLoadBalancer(sl.String(d.Id()))
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
		Timeout:        time.Duration(d.Get("wait_time_minutes").(int)) * time.Minute,
		Delay:          60 * time.Second,
		MinTimeout:     3 * time.Second,
		PollInterval:   60 * time.Second,
		NotFoundChecks: 40,
	}

	return stateConf.WaitForState()
}

func waitForLbaasLBDelete(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetNetworkLBaaSLoadBalancerService(sess)

	stateConf := &resource.StateChangeConf{
		Pending: []string{lbDeletePending},
		Target:  []string{lbDeleted},
		Refresh: func() (interface{}, string, error) {
			lb, err := service.GetLoadBalancer(sl.String(d.Id()))
			if err != nil {
				if apiErr, ok := err.(sl.Error); ok && (apiErr.StatusCode == 404 || apiErr.Exception == NOT_FOUND) {
					return lb, lbDeleted, nil
				}
				return datatypes.Network_LBaaS_LoadBalancer{}, "", err
			}
			return lb, lbDeletePending, nil
		},
		Timeout:      time.Duration(d.Get("wait_time_minutes").(int)) * time.Minute,
		Delay:        60 * time.Second,
		MinTimeout:   10 * time.Second,
		PollInterval: 60 * time.Second,
	}

	return stateConf.WaitForState()
}

func resourceIBMLBProtocolHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-",
		m["frontend_protocol"].(string)))
	buf.WriteString(fmt.Sprintf("%d-", m["frontend_port"].(int)))
	buf.WriteString(fmt.Sprintf("%s-",
		m["backend_protocol"].(string)))
	buf.WriteString(fmt.Sprintf("%d-", m["backend_port"].(int)))

	if v, ok := m["tls_certificate_id"]; ok && v.(int) != 0 {
		buf.WriteString(fmt.Sprintf("%d-", v.(int)))
	}

	return hashcode.String(buf.String())
}

func resourceIBMLBMemberHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-",
		m["private_ip_address"].(string)))

	return hashcode.String(buf.String())
}
