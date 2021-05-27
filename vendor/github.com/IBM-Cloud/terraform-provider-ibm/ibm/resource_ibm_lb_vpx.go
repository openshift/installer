// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
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
	PACKAGE_ID_APPLICATION_DELIVERY_CONTROLLER = 192
	DELIMITER                                  = "_"
)

func resourceIBMLbVpx() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMLbVpxCreate,
		Read:     resourceIBMLbVpxRead,
		Update:   resourceIBMLbVpxUpdate,
		Delete:   resourceIBMLbVpxDelete,
		Exists:   resourceIBMLbVpxExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name",
			},

			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the VPX",
			},

			"datacenter": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Datacenter name",
			},

			"speed": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Speed value",
			},

			"version": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "version info",
			},

			"plan": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Plan info",
			},

			"ip_count": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "IP address count",
			},

			"public_vlan_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Piblic VLAN id",
			},

			"public_subnet": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Public subnet",
			},

			"private_vlan_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Private VLAN id",
			},

			"private_subnet": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Private subnet",
			},

			"vip_pool": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of VIP ids",
			},

			"management_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "management IP address",
			},

			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "List of the tags",
			},
		},
	}
}

func getSubnetId(subnet string, meta interface{}) (int, error) {
	service := services.GetAccountService(meta.(ClientSession).SoftLayerSession())

	subnetInfo := strings.Split(subnet, "/")
	if len(subnetInfo) != 2 {
		return 0, fmt.Errorf(
			"Unable to parse the provided subnet: %s", subnet)
	}

	networkIdentifier := subnetInfo[0]
	cidr := subnetInfo[1]

	subnets, err := service.
		Mask("id").
		Filter(
			filter.Build(
				filter.Path("subnets.cidr").Eq(cidr),
				filter.Path("subnets.networkIdentifier").Eq(networkIdentifier),
			),
		).
		GetSubnets()

	if err != nil {
		return 0, fmt.Errorf("Error looking up Subnet: %s", err)
	}

	if len(subnets) < 1 {
		return 0, fmt.Errorf(
			"Unable to locate a subnet matching the provided subnet: %s", subnet)
	}

	return *subnets[0].Id, nil
}

func getVPXVersion(id int, sess *session.Session) (string, error) {
	service := services.GetNetworkApplicationDeliveryControllerService(sess)
	getObjectResult, err := service.Id(id).Mask("description").GetObject()

	if err != nil {
		return "", fmt.Errorf("Error retrieving VPX version: %s", err)
	}

	return strings.Split(*getObjectResult.Description, " ")[3], nil
}

func getVPXPriceItemKeyName(version string, speed int, plan string) string {
	name := "NETSCALER_VPX"
	speedMeasurements := "MBPS"

	floatVersion, err := strconv.ParseFloat(version, 10)
	if err != nil {
		return ("Invalid Version :" + version)
	}

	newVersion := strconv.FormatFloat(floatVersion, 'f', -1, 64)

	versionReplaced := strings.Replace(newVersion, ".", DELIMITER, -1)

	speedString := strconv.Itoa(speed) + speedMeasurements

	return strings.Join([]string{name, versionReplaced, speedString, strings.ToUpper(plan)}, DELIMITER)
}

func getPublicIpItemKeyName(ipCount int) string {

	var name string

	if ipCount == 1 {
		name = "STATIC_PUBLIC_IP_ADDRESS"
	} else {
		name = "STATIC_PUBLIC_IP_ADDRESSES"
	}
	ipCountString := strconv.Itoa(ipCount)

	return strings.Join([]string{ipCountString, name}, DELIMITER)
}

func findVPXPriceItems(version string, speed int, plan string, ipCount int, meta interface{}) ([]datatypes.Product_Item_Price, error) {
	sess := meta.(ClientSession).SoftLayerSession()

	// Get VPX package type.
	productPackage, err := product.GetPackageByType(sess, "ADDITIONAL_SERVICES_APPLICATION_DELIVERY_APPLIANCE")
	if err != nil {
		return []datatypes.Product_Item_Price{}, err
	}

	// Get VPX product items
	items, err := product.GetPackageProducts(sess, *productPackage.Id)
	if err != nil {
		return []datatypes.Product_Item_Price{}, err
	}

	// Get VPX and static IP items
	nadcKey := getVPXPriceItemKeyName(version, speed, plan)
	ipKey := getPublicIpItemKeyName(ipCount)

	var nadcItemPrice, ipItemPrice datatypes.Product_Item_Price

	for _, item := range items {
		itemKey := item.KeyName
		if strings.Contains(*itemKey, nadcKey) {
			nadcItemPrice = item.Prices[0]
		}
		if *itemKey == ipKey {
			ipItemPrice = item.Prices[0]
		}
	}

	var errorMessages []string

	if nadcItemPrice.Id == nil {
		errorMessages = append(errorMessages, "VPX version, speed or plan have incorrect values")
	}

	if ipItemPrice.Id == nil {
		errorMessages = append(errorMessages, "IP quantity value is incorrect")
	}

	if len(errorMessages) > 0 {
		err = errors.New(strings.Join(errorMessages, "\n"))
		return []datatypes.Product_Item_Price{}, err
	}

	return []datatypes.Product_Item_Price{
		{
			Id: nadcItemPrice.Id,
		},
		{
			Id: ipItemPrice.Id,
		},
	}, nil
}

func findVPXByOrderId(orderId int, meta interface{}) (datatypes.Network_Application_Delivery_Controller, error) {
	service := services.GetAccountService(meta.(ClientSession).SoftLayerSession())

	stateConf := &resource.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"complete"},
		Refresh: func() (interface{}, string, error) {
			vpxs, err := service.
				Filter(
					filter.Build(
						filter.Path("applicationDeliveryControllers.billingItem.orderItem.order.id").Eq(orderId),
					),
				).GetApplicationDeliveryControllers()
			if err != nil {
				return datatypes.Network_Application_Delivery_Controller{}, "", err
			}

			if len(vpxs) == 1 {
				return vpxs[0], "complete", nil
			} else if len(vpxs) == 0 {
				return nil, "pending", nil
			} else {
				return nil, "", fmt.Errorf("Expected one VPX: %s", err)
			}
		},
		Timeout:    45 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	pendingResult, err := stateConf.WaitForState()

	if err != nil {
		return datatypes.Network_Application_Delivery_Controller{}, err
	}

	var result, ok = pendingResult.(datatypes.Network_Application_Delivery_Controller)

	if ok {
		return result, nil
	}

	return datatypes.Network_Application_Delivery_Controller{},
		fmt.Errorf("Cannot find Application Delivery Controller with order id '%d'", orderId)
}

func prepareHardwareOptions(d *schema.ResourceData, meta interface{}) ([]datatypes.Hardware, error) {
	hardwareOpts := make([]datatypes.Hardware, 1)
	publicVlanId := d.Get("public_vlan_id").(int)
	publicSubnet := d.Get("public_subnet").(string)

	if publicVlanId > 0 || len(publicSubnet) > 0 {
		hardwareOpts[0].PrimaryNetworkComponent = &datatypes.Network_Component{}
	}

	if publicVlanId > 0 {
		hardwareOpts[0].PrimaryNetworkComponent.NetworkVlanId = &publicVlanId
	}

	if len(publicSubnet) > 0 {
		primarySubnetId, err := getSubnetId(publicSubnet, meta)
		if err != nil {
			return nil, fmt.Errorf("Error creating network application delivery controller: %s", err)
		}
		hardwareOpts[0].PrimaryNetworkComponent.NetworkVlan = &datatypes.Network_Vlan{
			PrimarySubnetId: &primarySubnetId,
		}
	}

	privateVlanId := d.Get("private_vlan_id").(int)
	privateSubnet := d.Get("private_subnet").(string)
	if privateVlanId > 0 || len(privateSubnet) > 0 {
		hardwareOpts[0].PrimaryBackendNetworkComponent = &datatypes.Network_Component{}
	}

	if privateVlanId > 0 {
		hardwareOpts[0].PrimaryBackendNetworkComponent.NetworkVlanId = &privateVlanId
	}

	if len(privateSubnet) > 0 {
		primarySubnetId, err := getSubnetId(privateSubnet, meta)
		if err != nil {
			return nil, fmt.Errorf("Error creating network application delivery controller: %s", err)
		}
		hardwareOpts[0].PrimaryBackendNetworkComponent.NetworkVlan = &datatypes.Network_Vlan{
			PrimarySubnetId: &primarySubnetId,
		}
	}
	return hardwareOpts, nil
}

func resourceIBMLbVpxCreate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	NADCService := services.GetNetworkApplicationDeliveryControllerService(sess)
	productOrderService := services.GetProductOrderService(sess.SetRetries(0))
	var err error

	opts := datatypes.Container_Product_Order{
		PackageId: sl.Int(PACKAGE_ID_APPLICATION_DELIVERY_CONTROLLER),
		Quantity:  sl.Int(1),
	}

	opts.Prices, err = findVPXPriceItems(
		d.Get("version").(string),
		d.Get("speed").(int),
		d.Get("plan").(string),
		d.Get("ip_count").(int),
		meta)

	if err != nil {
		return fmt.Errorf("Error Cannot find Application Delivery Controller prices '%s'.", err)
	}

	datacenter := d.Get("datacenter").(string)

	if len(datacenter) > 0 {
		datacenter, err := location.GetDatacenterByName(sess, datacenter, "id")
		if err != nil {
			return fmt.Errorf("Error creating network application delivery controller: %s", err)
		}
		opts.Location = sl.String(strconv.Itoa(*datacenter.Id))
	}

	opts.Hardware, err = prepareHardwareOptions(d, meta)
	if err != nil {
		return fmt.Errorf("Error Cannot get hardware options '%s'.", err)
	}

	log.Println("[INFO] Creating network application delivery controller")

	receipt, err := productOrderService.PlaceOrder(&opts, sl.Bool(false))

	if err != nil {
		return fmt.Errorf("Error creating network application delivery controller: %s", err)
	}

	// Wait VPX provisioning
	VPX, err := findVPXByOrderId(*receipt.OrderId, meta)

	if err != nil {
		return fmt.Errorf("Error creating network application delivery controller: %s", err)
	}

	d.SetId(fmt.Sprintf("%d", *VPX.Id))

	log.Printf("[INFO] Netscaler VPX ID: %s", d.Id())

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	// Wait Virtual IP provisioning
	IsVipReady := false

	for vipWaitCount := 0; vipWaitCount < 270; vipWaitCount++ {
		getObjectResult, err := NADCService.Id(id).Mask("subnets[ipAddresses],password[password]").GetObject()
		if err != nil {
			return fmt.Errorf("Error retrieving network application delivery controller: %s", err)
		}

		ipCount := 0
		if getObjectResult.Password != nil && getObjectResult.Password.Password != nil && len(*getObjectResult.Password.Password) > 0 &&
			getObjectResult.Subnets != nil && len(getObjectResult.Subnets) > 0 && getObjectResult.Subnets[0].IpAddresses != nil {
			ipCount = len(getObjectResult.Subnets[0].IpAddresses)
		}
		if ipCount > 0 {
			IsVipReady = true
			break
		}
		log.Printf("[INFO] Wait 10 seconds for Virtual IP provisioning on Netscaler VPX ID: %d", id)
		time.Sleep(time.Second * 10)
	}

	if !IsVipReady {
		return fmt.Errorf("Failed to create VIPs for Netscaler VPX ID: %d", id)
	}

	// Wait while VPX service is initializing. GetLoadBalancers() internally calls REST API of VPX and returns
	// an error "Could not connect to host" if the REST API is not available.
	IsRESTReady := false

	for restWaitCount := 0; restWaitCount < 270; restWaitCount++ {
		_, err := NADCService.Id(id).GetLoadBalancers()
		// GetLoadBalancers returns an error "There was a problem processing the reply from the
		// application tier.  Please contact development." if the VPX version is 10.5.
		if err == nil || !strings.Contains(err.Error(), "Could not connect to host") {
			IsRESTReady = true
			break
		}
		log.Printf("[INFO] Wait 10 seconds for VPX(%d) REST Service ID", id)
		time.Sleep(time.Second * 10)
	}

	if !IsRESTReady {
		return fmt.Errorf("Failed to intialize VPX REST Service for Netscaler VPX ID: %d", id)
	}

	// Wait additional buffer time for VPX service.
	time.Sleep(time.Minute)

	return resourceIBMLbVpxRead(d, meta)
}

func resourceIBMLbVpxRead(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()

	service := services.GetNetworkApplicationDeliveryControllerService(sess)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	getObjectResult, err := service.
		Id(id).
		Mask("id,name,type[name],datacenter,networkVlans[primaryRouter],networkVlans[primarySubnets],subnets[ipAddresses],description,managementIpAddress").
		GetObject()

	if err != nil {
		return fmt.Errorf("Error retrieving network application delivery controller: %s", err)
	}

	d.Set("name", *getObjectResult.Name)
	d.Set("type", *getObjectResult.Type.Name)
	if getObjectResult.Datacenter != nil {
		d.Set("datacenter", *getObjectResult.Datacenter.Name)
	}

	for _, vlan := range getObjectResult.NetworkVlans {
		if vlan.PrimaryRouter != nil && *vlan.PrimaryRouter.Hostname != "" {
			isFcr := strings.HasPrefix(*vlan.PrimaryRouter.Hostname, "fcr")
			isBcr := strings.HasPrefix(*vlan.PrimaryRouter.Hostname, "bcr")
			if isFcr {
				d.Set("public_vlan_id", *vlan.Id)
				if vlan.PrimarySubnets != nil && len(vlan.PrimarySubnets) > 0 {
					ipAddress := *vlan.PrimarySubnets[0].NetworkIdentifier
					d.Set(
						"public_subnet",
						fmt.Sprintf("%s/%d", ipAddress, *vlan.PrimarySubnets[0].Cidr),
					)
				}
			}

			if isBcr {
				d.Set("private_vlan_id", *vlan.Id)
				if vlan.PrimarySubnets != nil && len(vlan.PrimarySubnets) > 0 {
					ipAddress := *vlan.PrimarySubnets[0].NetworkIdentifier
					d.Set(
						"private_subnet",
						fmt.Sprintf("%s/%d", ipAddress, *vlan.PrimarySubnets[0].Cidr),
					)
				}
			}
		}
	}

	vips := make([]string, 0)
	ipCount := 0
	for i, subnet := range getObjectResult.Subnets {
		for _, ipAddressObj := range subnet.IpAddresses {
			vips = append(vips, *ipAddressObj.IpAddress)
			if i == 0 {
				ipCount++
			}
		}
	}

	d.Set("vip_pool", vips)
	d.Set("ip_count", ipCount)
	d.Set("management_ip_address", *getObjectResult.ManagementIpAddress)

	description := *getObjectResult.Description
	r, _ := regexp.Compile(" [0-9]+Mbps")
	speedStr := r.FindString(description)
	r, _ = regexp.Compile("[0-9]+")
	speed, err := strconv.Atoi(r.FindString(speedStr))
	if err == nil && speed > 0 {
		d.Set("speed", speed)
	}

	r, _ = regexp.Compile(" VPX [0-9]+\\.[0-9]+ ")
	versionStr := r.FindString(description)
	r, _ = regexp.Compile("[0-9]+\\.[0-9]+")
	version := r.FindString(versionStr)
	if version != "" {
		d.Set("version", version)
	}

	r, _ = regexp.Compile(" [A-Za-z]+$")
	planStr := r.FindString(description)
	r, _ = regexp.Compile("[A-Za-z]+$")
	plan := r.FindString(planStr)
	if plan != "" {
		d.Set("plan", plan)
	}

	return nil
}

func resourceIBMLbVpxUpdate(d *schema.ResourceData, meta interface{}) error {
	//Only tags are updated and that too locally hence nothing to validate and update in terms of real API at this point
	return nil
}

func resourceIBMLbVpxDelete(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetNetworkApplicationDeliveryControllerService(sess)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	billingItem, err := service.Id(id).GetBillingItem()
	if err != nil {
		return fmt.Errorf("Error deleting network application delivery controller: %s", err)
	}

	if *billingItem.Id > 0 {
		billingItemService := services.GetBillingItemService(sess)
		deleted, err := billingItemService.Id(*billingItem.Id).CancelService()
		if err != nil {
			return fmt.Errorf("Error deleting network application delivery controller: %s", err)
		}

		if deleted {
			return nil
		}
	}

	return nil
}

func resourceIBMLbVpxExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	service := services.GetNetworkApplicationDeliveryControllerService(meta.(ClientSession).SoftLayerSession())

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}
	nadc, err := service.Mask("id").Id(id).GetObject()
	if err != nil {
		if apiErr, ok := err.(sl.Error); ok {
			if apiErr.StatusCode == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error communicating with the API: %s", err)
	}
	return nadc.Id != nil && *nadc.Id == id, nil
}
