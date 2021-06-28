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
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/helpers/location"
	"github.com/softlayer/softlayer-go/helpers/product"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
)

const (
	LB_LARGE_150000_CONNECTIONS = 150000
	LB_SMALL_15000_CONNECTIONS  = 15000

	LbLocalPackageType = "ADDITIONAL_SERVICES_LOAD_BALANCER"

	lbMask = "id,dedicatedFlag,connectionLimit,ipAddressId,securityCertificateId,highAvailabilityFlag," +
		"sslEnabledFlag,sslActiveFlag,loadBalancerHardware[datacenter[name]],ipAddress[ipAddress,subnetId],billingItem[upgradeItems[capacity]]"
)

func resourceIBMLb() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMLbCreate,
		Read:     resourceIBMLbRead,
		Update:   resourceIBMLbUpdate,
		Delete:   resourceIBMLbDelete,
		Exists:   resourceIBMLbExists,
		Importer: &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"connections": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Connections value",
			},
			"datacenter": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Datacenter name info",
			},
			"ha_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "true if High availability is enabled",
			},
			"security_certificate_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Security certificate ID",
			},
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"dedicated": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Boolena value true if Load balncer is dedicated type",
			},
			"ssl_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ssl_offload": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "boolean value true if SSL offload is enabled",
			},
			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "Tags associated with resource",
			},

			"hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIBMLbCreate(d *schema.ResourceData, meta interface{}) error {

	sess := meta.(ClientSession).SoftLayerSession()

	connections := d.Get("connections").(int)
	haEnabled := d.Get("ha_enabled").(bool)
	dedicated := d.Get("dedicated").(bool)

	var categoryCode string

	// SoftLayer capacities don't match the published capacities as seen in the local lb
	// ordering screen in the customer portal. Terraform exposes the published capacities.
	// Create a translation map for those cases where the published capacity does not
	// equal the actual actual capacity on the product_item.
	capacities := map[int]float64{
		15000:  65000.0,
		150000: 130000.0,
	}

	var capacity float64
	if c, ok := capacities[connections]; !ok {
		capacity = float64(connections)
	} else {
		capacity = c
	}

	var keyFormatter string
	if dedicated {
		// Dedicated local LB always comes with SSL support
		d.Set("ssl_enabled", true)
		categoryCode = product.DedicatedLoadBalancerCategoryCode
		if haEnabled {
			keyFormatter = "DEDICATED_LOAD_BALANCER_WITH_HIGH_AVAILABILITY_AND_SSL_%d_CONNECTIONS"
		} else {
			keyFormatter = "LOAD_BALANCER_DEDICATED_WITH_SSL_OFFLOAD_%d_CONNECTIONS"
		}
	} else {
		if d.Get("ha_enabled").(bool) {
			return fmt.Errorf("High Availability is not supported for shared local load balancers")
		}
		categoryCode = product.ProxyLoadBalancerCategoryCode
		if _, ok := d.GetOk("security_certificate_id"); ok {
			d.Set("ssl_enabled", true)
			keyFormatter = "LOAD_BALANCER_%d_VIP_CONNECTIONS_WITH_SSL_OFFLOAD"
		} else {
			d.Set("ssl_enabled", false)
			keyFormatter = "LOAD_BALANCER_%d_VIP_CONNECTIONS"
		}
	}

	keyName := fmt.Sprintf(keyFormatter, connections)

	pkg, err := product.GetPackageByType(sess, LbLocalPackageType)
	if err != nil {
		return err
	}

	// Get all prices for ADDITIONAL_SERVICE_LOAD_BALANCER with the given capacity
	productItems, err := product.GetPackageProducts(sess, *pkg.Id)
	if err != nil {
		return err
	}

	// Select only those product items with a matching keyname
	targetItems := []datatypes.Product_Item{}
	for _, item := range productItems {
		if *item.KeyName == keyName {
			targetItems = append(targetItems, item)
		}
	}

	if len(targetItems) == 0 {
		return fmt.Errorf("No product items matching %s could be found", keyName)
	}

	//select prices with the required capacity
	prices := product.SelectProductPricesByCategory(
		targetItems,
		map[string]float64{
			categoryCode: capacity,
		},
	)

	// Lookup the datacenter ID
	dc, err := location.GetDatacenterByName(sess, d.Get("datacenter").(string))

	productOrderContainer := datatypes.Container_Product_Order_Network_LoadBalancer{
		Container_Product_Order: datatypes.Container_Product_Order{
			PackageId: pkg.Id,
			Location:  sl.String(strconv.Itoa(*dc.Id)),
			Prices:    prices[:1],
			Quantity:  sl.Int(1),
		},
	}

	log.Println("[INFO] Creating load balancer")

	receipt, err := services.GetProductOrderService(sess.SetRetries(0)).
		PlaceOrder(&productOrderContainer, sl.Bool(false))
	if err != nil {
		return fmt.Errorf("Error during creation of load balancer: %s", err)
	}

	loadBalancer, err := findLoadBalancerByOrderId(sess, *receipt.OrderId, dedicated, d)
	if err != nil {
		return fmt.Errorf("Error during creation of load balancer: %s", err)
	}

	d.SetId(fmt.Sprintf("%d", *loadBalancer.Id))
	d.Set("connections", getConnectionLimit(*loadBalancer.ConnectionLimit))
	d.Set("datacenter", loadBalancer.LoadBalancerHardware[0].Datacenter.Name)
	d.Set("ip_address", loadBalancer.IpAddress.IpAddress)
	d.Set("subnet_id", loadBalancer.IpAddress.SubnetId)
	d.Set("ha_enabled", loadBalancer.HighAvailabilityFlag)

	log.Printf("[INFO] Load Balancer ID: %s", d.Id())

	return resourceIBMLbUpdate(d, meta)
}

func resourceIBMLbUpdate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()

	vipID, _ := strconv.Atoi(d.Id())

	certID := d.Get("security_certificate_id").(int)

	err := setLocalLBSecurityCert(sess, vipID, certID)
	if err != nil {
		return fmt.Errorf("Update load balancer failed: %s", err)
	}

	if d.HasChange("connections") {
		vip, err := services.GetNetworkApplicationDeliveryControllerLoadBalancerVirtualIpAddressService(sess).
			Id(vipID).
			Mask(lbMask).
			GetObject()
		if err != nil {
			return err
		}
		ors, nrs := d.GetChange("connections")
		oldValue := ors.(int)
		newValue := nrs.(int)

		if oldValue > 0 {
			if *vip.DedicatedFlag {
				return fmt.Errorf("Error Updating load balancer connection limit: Upgrade for dedicated loadbalancer is not supported")
			}
			if vip.BillingItem.UpgradeItems[0].Capacity != nil {
				validUpgradeValue := vip.BillingItem.UpgradeItems[0].Capacity
				if newValue == int(*validUpgradeValue) {
					_, err := services.GetNetworkApplicationDeliveryControllerLoadBalancerVirtualIpAddressService(sess).
						Id(vipID).UpgradeConnectionLimit()
					if err != nil {
						return fmt.Errorf("Error Updating load balancer connection limit: %s", err)
					}
				} else {

					return fmt.Errorf("Error Updating load balancer connection limit : Valid value to which connection limit can be upgraded is : %d ", int(*validUpgradeValue))

				}

			} else {
				return fmt.Errorf("Error Updating load balancer connection limit: No upgrade available, already it has maximum connection limit")
			}
		}

	}

	if d.HasChange("ssl_offload") && !d.IsNewResource() {

		if d.Get("ssl_offload").(bool) {

			_, err := services.GetNetworkApplicationDeliveryControllerLoadBalancerVirtualIpAddressService(sess).
				Id(vipID).StartSsl()
			if err != nil {
				return fmt.Errorf("Error starting ssl acceleration for load balancer : %s", err)
			}

		} else {

			_, err := services.GetNetworkApplicationDeliveryControllerLoadBalancerVirtualIpAddressService(sess).
				Id(vipID).StopSsl()
			if err != nil {
				return fmt.Errorf("Error stopping ssl acceleration for load balancer : %s", err)
			}

		}
	}

	return resourceIBMLbRead(d, meta)
}

func resourceIBMLbRead(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	vipID, _ := strconv.Atoi(d.Id())

	vip, err := services.GetNetworkApplicationDeliveryControllerLoadBalancerVirtualIpAddressService(sess).
		Id(vipID).
		Mask(lbMask).
		GetObject()

	if err != nil {
		return fmt.Errorf("Error retrieving load balancer: %s", err)
	}

	d.Set("connections", getConnectionLimit(*vip.ConnectionLimit))
	d.Set("datacenter", vip.LoadBalancerHardware[0].Datacenter.Name)
	d.Set("ip_address", vip.IpAddress.IpAddress)
	d.Set("subnet_id", vip.IpAddress.SubnetId)
	d.Set("ha_enabled", vip.HighAvailabilityFlag)
	d.Set("dedicated", vip.DedicatedFlag)
	d.Set("ssl_enabled", vip.SslEnabledFlag)
	d.Set("ssl_offload", vip.SslActiveFlag)
	// Optional fields.  Guard against nil pointer dereferences
	d.Set("security_certificate_id", sl.Get(vip.SecurityCertificateId, nil))
	d.Set("hostname", vip.LoadBalancerHardware[0].Hostname)
	return nil
}

func resourceIBMLbDelete(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	vipService := services.GetNetworkApplicationDeliveryControllerLoadBalancerVirtualIpAddressService(sess)
	vipID, _ := strconv.Atoi(d.Id())

	certID := d.Get("security_certificate_id").(int)

	if certID > 0 {
		err := setLocalLBSecurityCert(sess, vipID, 0)
		if err != nil {
			return fmt.Errorf("Remove certificate before deleting load balancer failed: %s", err)
		}

	}

	var billingItem datatypes.Billing_Item_Network_LoadBalancer
	var err error

	// Get billing item associated with the load balancer
	if d.Get("dedicated").(bool) {
		billingItem, err = vipService.
			Id(vipID).
			GetDedicatedBillingItem()
	} else {
		billingItem.Billing_Item, err = vipService.
			Id(vipID).
			GetBillingItem()
	}

	if err != nil {
		return fmt.Errorf("Error while looking up billing item associated with the load balancer: %s", err)
	}

	if billingItem.Id == nil {
		return fmt.Errorf("Error while looking up billing item associated with the load balancer: No billing item for ID:%d", vipID)
	}
	success, err := services.GetBillingItemService(sess).Id(*billingItem.Id).CancelService()
	if err != nil {
		return err
	}

	if !success {
		return fmt.Errorf("SoftLayer reported an unsuccessful cancellation")
	}

	return nil
}

func resourceIBMLbExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess := meta.(ClientSession).SoftLayerSession()

	vipID, _ := strconv.Atoi(d.Id())

	_, err := services.GetNetworkApplicationDeliveryControllerLoadBalancerVirtualIpAddressService(sess).
		Id(vipID).
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

/* When requesting 15000 SL creates between 15000 and 150000. When requesting 150000 SL creates >= 150000 */
func getConnectionLimit(connectionLimit int) int {
	if connectionLimit >= LB_LARGE_150000_CONNECTIONS {
		return LB_LARGE_150000_CONNECTIONS
	} else if connectionLimit >= LB_SMALL_15000_CONNECTIONS &&
		connectionLimit < LB_LARGE_150000_CONNECTIONS {
		return LB_SMALL_15000_CONNECTIONS
	} else {
		return connectionLimit
	}
}

func findLoadBalancerByOrderId(sess *session.Session, orderId int, dedicated bool, d *schema.ResourceData) (datatypes.Network_Application_Delivery_Controller_LoadBalancer_VirtualIpAddress, error) {
	var filterPath string
	if dedicated {
		filterPath = "adcLoadBalancers.dedicatedBillingItem.orderItem.order.id"
	} else {
		filterPath = "adcLoadBalancers.billingItem.orderItem.order.id"
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"complete"},
		Refresh: func() (interface{}, string, error) {
			lbs, err := services.GetAccountService(sess).
				Filter(filter.Build(
					filter.Path(filterPath).
						Eq(strconv.Itoa(orderId)))).
				Mask(lbMask).
				GetAdcLoadBalancers()
			if err != nil {
				return datatypes.Network_Application_Delivery_Controller_LoadBalancer_VirtualIpAddress{}, "", err
			}

			if len(lbs) == 1 {
				return lbs[0], "complete", nil
			} else if len(lbs) == 0 {
				return datatypes.Network_Application_Delivery_Controller_LoadBalancer_VirtualIpAddress{}, "pending", nil
			} else {
				return nil, "", fmt.Errorf("Expected one load balancer: %s", err)
			}
		},
		Timeout:        d.Timeout(schema.TimeoutCreate),
		Delay:          5 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 24 * 60,
	}

	pendingResult, err := stateConf.WaitForState()

	if err != nil {
		return datatypes.Network_Application_Delivery_Controller_LoadBalancer_VirtualIpAddress{}, err
	}

	var result, ok = pendingResult.(datatypes.Network_Application_Delivery_Controller_LoadBalancer_VirtualIpAddress)

	if ok {
		return result, nil
	}

	return datatypes.Network_Application_Delivery_Controller_LoadBalancer_VirtualIpAddress{},
		fmt.Errorf("Cannot find Application Delivery Controller Load Balancer with order id '%d'", orderId)
}

func setLocalLBSecurityCert(sess *session.Session, vipID int, certID int) error {
	var vip struct {
		SecurityCertificateId *int `json:"securityCertificateId"`
	}

	var success bool

	if certID == 0 {
		vip.SecurityCertificateId = nil
	} else {
		vip.SecurityCertificateId = &certID
	}

	// In order to send a null value, need to invoke DoRequest directly with a custom struct
	err := sess.DoRequest(
		"SoftLayer_Network_Application_Delivery_Controller_LoadBalancer_VirtualIpAddress",
		"editObject",
		[]interface{}{&vip},
		&sl.Options{Id: &vipID},
		&success,
	)

	if !success && err == nil {
		return fmt.Errorf("Unable to remove ssl security certificate from load balancer")
	}

	return err
}
