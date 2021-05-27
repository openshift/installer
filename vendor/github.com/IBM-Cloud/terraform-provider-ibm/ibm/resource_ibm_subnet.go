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

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/helpers/product"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
)

const (
	SubnetMask = "id,addressSpace,subnetType,version,ipAddressCount," +
		"networkIdentifier,cidr,note,endPointIpAddress[ipAddress],networkVlan[id],totalIpAddresses"
)

var (
	// Map subnet types to product package keyname in SoftLayer_Product_Item
	subnetPackageTypeMap = map[string]string{
		"Static":   "ADDITIONAL_SERVICES_STATIC_IP_ADDRESSES",
		"Portable": "ADDITIONAL_SERVICES_PORTABLE_IP_ADDRESSES",
	}
)

func resourceIBMSubnet() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMSubnetCreate,
		Read:     resourceIBMSubnetRead,
		Update:   resourceIBMSubnetUpdate,
		Delete:   resourceIBMSubnetDelete,
		Exists:   resourceIBMSubnetExists,
		Importer: &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			"private": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "private subnet",
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errs []error) {
					typeStr := v.(string)
					if typeStr != "Portable" && typeStr != "Static" {
						errs = append(errs, errors.New(
							"type should be either Portable or Static"))
					}
					return
				},
				Description: "subnet type",
			},

			// IP version 4 or IP version 6
			"ip_version": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  4,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errs []error) {
					ipVersion := v.(int)
					if ipVersion != 4 && ipVersion != 6 {
						errs = append(errs, errors.New(
							"ip version should be either 4 or 6"))
					}
					return
				},
				Description: "ip version",
			},

			"capacity": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "number of ip addresses in the subnet",
			},

			// vlan_id should be configured when type is "Portable"
			"vlan_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"endpoint_ip"},
				Description:   "VLAN ID for the subnet",
			},

			// endpoint_ip should be configured when type is "Static"
			"endpoint_ip": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"vlan_id"},
				Description:   "endpoint IP",
			},

			// Provides IP address/cidr format (ex. 10.10.10.10/28)
			"subnet_cidr": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CIDR notation for the subnet",
			},

			"notes": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Notes",
			},

			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "tags set for the resource",
			},
		},
	}
}

func resourceIBMSubnetCreate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()

	// Find price items with AdditionalServicesSubnetAddresses
	productOrderContainer, err := buildSubnetProductOrderContainer(d, sess)
	if err != nil {
		return fmt.Errorf("Error creating subnet: %s", err)
	}

	log.Println("[INFO] Creating subnet")

	receipt, err := services.GetProductOrderService(sess.SetRetries(0)).
		PlaceOrder(productOrderContainer, sl.Bool(false))
	if err != nil {
		return fmt.Errorf("Error during creation of subnet: %s", err)
	}

	Subnet, err := findSubnetByOrderID(sess, *receipt.OrderId, d)
	if err != nil {
		return fmt.Errorf("Error during creation of subnet: %s", err)
	}

	d.SetId(fmt.Sprintf("%d", *Subnet.Id))

	return resourceIBMSubnetUpdate(d, meta)
}

func resourceIBMSubnetRead(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetNetworkSubnetService(sess)

	subnetID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Not a valid subnet ID, must be an integer: %s", err)
	}

	subnet, err := service.Id(subnetID).Mask(SubnetMask).GetObject()
	if err != nil {
		return fmt.Errorf("Error retrieving a subnet: %s", err)
	}

	if *subnet.AddressSpace == "PRIVATE" {
		d.Set("private", true)
	} else if *subnet.AddressSpace == "PUBLIC" {
		d.Set("private", false)
	}

	if subnet.SubnetType == nil {
		return fmt.Errorf("Invalid vlan type: the subnet type is null")
	}
	if strings.Contains(*subnet.SubnetType, "STATIC") {
		d.Set("type", "Static")
	} else if strings.Contains(*subnet.SubnetType, "VLAN") {
		d.Set("type", "Portable")
	} else {
		return fmt.Errorf("Invalid vlan type: %s", *subnet.SubnetType)
	}
	d.Set("ip_version", *subnet.Version)
	d.Set("capacity", *subnet.TotalIpAddresses)
	if *subnet.Version == 6 {
		d.Set("capacity", 64)
	}
	d.Set("subnet_cidr", *subnet.NetworkIdentifier+"/"+strconv.Itoa(*subnet.Cidr))
	if subnet.Note != nil {
		d.Set("notes", *subnet.Note)
	}
	if subnet.EndPointIpAddress != nil {
		d.Set("endpoint_ip", *subnet.EndPointIpAddress.IpAddress)
	}
	if subnet.NetworkVlan != nil {
		d.Set("vlan_id", subnet.NetworkVlan.Id)
	}
	d.Set("notes", sl.Get(subnet.Note, nil))

	return nil
}

func resourceIBMSubnetUpdate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetNetworkSubnetService(sess)

	subnetID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Not a valid subnet ID, must be an integer: %s", err)
	}

	if d.HasChange("notes") {
		_, err = service.Id(subnetID).EditNote(sl.String(d.Get("notes").(string)))
		if err != nil {
			return fmt.Errorf("Error updating subnet: %s", err)
		}
	}
	return resourceIBMSubnetRead(d, meta)
}

func resourceIBMSubnetDelete(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetNetworkSubnetService(sess)

	subnetID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Not a valid subnet ID, must be an integer: %s", err)
	}

	billingItem, err := service.Id(subnetID).GetBillingItem()
	if err != nil {
		return fmt.Errorf("Error deleting subnet: %s", err)
	}

	if billingItem.Id == nil {
		return nil
	}
	_, err = services.GetBillingItemService(sess).Id(*billingItem.Id).CancelService()
	if err != nil {
		return fmt.Errorf("Error deleting subnet: %s", err)
	}

	return err
}

func resourceIBMSubnetExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetNetworkSubnetService(sess)

	subnetID, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	result, err := service.Id(subnetID).GetObject()
	if err != nil {
		if apiErr, ok := err.(sl.Error); ok && apiErr.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("Error retrieving subnet: %s", err)
	}
	return result.Id != nil && *result.Id == subnetID, nil
}

func findSubnetByOrderID(sess *session.Session, orderID int, d *schema.ResourceData) (datatypes.Network_Subnet, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"complete"},
		Refresh: func() (interface{}, string, error) {
			subnets, err := services.GetAccountService(sess).
				Filter(filter.Path("subnets.billingItem.orderItem.order.id").
					Eq(strconv.Itoa(orderID)).Build()).
				Mask("id,activeTransaction").
				GetSubnets()
			if err != nil {
				return datatypes.Network_Subnet{}, "", err
			}

			if len(subnets) == 1 && subnets[0].ActiveTransaction == nil {
				return subnets[0], "complete", nil
			}
			return nil, "pending", nil
		},
		Timeout:        d.Timeout(schema.TimeoutCreate),
		Delay:          5 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 1440,
	}

	pendingResult, err := stateConf.WaitForState()

	if err != nil {
		return datatypes.Network_Subnet{}, err
	}

	if result, ok := pendingResult.(datatypes.Network_Subnet); ok {
		return result, nil
	}

	return datatypes.Network_Subnet{},
		fmt.Errorf("Cannot find a subnet with order id '%d'", orderID)
}

func buildSubnetProductOrderContainer(d *schema.ResourceData, sess *session.Session) (
	*datatypes.Container_Product_Order_Network_Subnet, error) {

	// 1. Get a package
	typeStr := d.Get("type").(string)
	vlanID := d.Get("vlan_id").(int)
	private := d.Get("private").(bool)
	network := "PUBLIC"
	if private {
		network = "PRIVATE"
	}

	pkg, err := product.GetPackageByType(sess, subnetPackageTypeMap[typeStr])
	if err != nil {
		return &datatypes.Container_Product_Order_Network_Subnet{}, err
	}

	// 2. Get all prices for the package
	productItems, err := product.GetPackageProducts(sess, *pkg.Id)
	if err != nil {
		return &datatypes.Container_Product_Order_Network_Subnet{}, err
	}

	// 3. Select items which have a matching capacity, network, and IP version.
	capacity := d.Get("capacity").(int)
	ipVersionStr := "_IP_"
	if d.Get("ip_version").(int) == 6 {
		ipVersionStr = "_IPV6_"
	}
	SubnetItems := []datatypes.Product_Item{}
	for _, item := range productItems {
		if int(*item.Capacity) == d.Get("capacity").(int) &&
			strings.Contains(*item.KeyName, network) &&
			strings.Contains(*item.KeyName, ipVersionStr) {
			SubnetItems = append(SubnetItems, item)
		}
	}

	if len(SubnetItems) == 0 {
		return &datatypes.Container_Product_Order_Network_Subnet{},
			fmt.Errorf("No product items matching with capacity %d could be found", capacity)
	}

	productOrderContainer := datatypes.Container_Product_Order_Network_Subnet{
		Container_Product_Order: datatypes.Container_Product_Order{
			PackageId: pkg.Id,
			Prices: []datatypes.Product_Item_Price{
				{
					Id: SubnetItems[0].Prices[0].Id,
				},
			},
			Quantity: sl.Int(1),
		},
		EndPointVlanId: sl.Int(vlanID),
	}

	if endpointIP, ok := d.GetOk("endpoint_ip"); ok {
		if typeStr != "Static" {
			return &datatypes.Container_Product_Order_Network_Subnet{},
				fmt.Errorf("endpoint_ip is only available when type is Static")
		}
		endpointIPStr := endpointIP.(string)
		subnet, err := services.GetNetworkSubnetService(sess).Mask("ipAddresses").GetSubnetForIpAddress(sl.String(endpointIPStr))
		if err != nil {
			return &datatypes.Container_Product_Order_Network_Subnet{}, err
		}
		for _, ipSubnet := range subnet.IpAddresses {
			if *ipSubnet.IpAddress == endpointIPStr {
				productOrderContainer.EndPointIpAddressId = ipSubnet.Id
			}
		}
		if productOrderContainer.EndPointIpAddressId == nil {
			return &datatypes.Container_Product_Order_Network_Subnet{},
				fmt.Errorf("Unable to find an ID of ipAddress: %s", endpointIPStr)
		}
	}
	return &productOrderContainer, nil
}

func getVlanType(sess *session.Session, vlanID int) (string, error) {
	vlan, err := services.GetNetworkVlanService(sess).Id(vlanID).Mask(VlanMask).GetObject()

	if err != nil {
		return "", fmt.Errorf("Error retrieving vlan: %s", err)
	}

	if vlan.PrimaryRouter != nil {
		if strings.HasPrefix(*vlan.PrimaryRouter.Hostname, "fcr") {
			return "PUBLIC", nil
		} else {
			return "PRIVATE", nil
		}
	}
	return "", fmt.Errorf("Unable to determine network")
}
