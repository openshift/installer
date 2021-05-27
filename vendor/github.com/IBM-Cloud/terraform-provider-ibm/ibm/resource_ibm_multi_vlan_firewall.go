// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/helpers/location"
	"github.com/softlayer/softlayer-go/helpers/product"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/sl"
)

func resourceIBMMultiVlanFirewall() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMNetworkMultiVlanCreate,
		Read:     resourceIBMMultiVlanFirewallRead,
		Delete:   resourceIBMFirewallDelete,
		Update:   resourceIBMMultiVlanFirewallUpdate,
		Exists:   resourceIBMMultiVLanFirewallExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"datacenter": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Datacenter name",
			},

			"pod": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.TrimSpace(old) == strings.TrimSpace(new)
				},
				Description: "POD name",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "name",
			},

			"public_vlan_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Public VLAN id",
			},

			"private_vlan_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Private VLAN id",
			},

			"firewall_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue([]string{"FortiGate Firewall Appliance HA Option", "FortiGate Security Appliance"}),
				Description:  "Firewall type",
			},

			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP Address",
			},

			"public_ipv6": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IPV6 IP",
			},

			"private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Private IP Address",
			},

			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User name",
			},

			"password": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Password",
			},

			"addon_configuration": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: `High Availability - [Web Filtering Add-on, NGFW Add-on, AV Add-on] or [Web Filtering Add-on, NGFW Add-on, AV Add-on]`,
			},
		},
	}
}

const (
	productPackageFilter      = `{"keyName":{"operation":"FIREWALL_APPLIANCE"}}`
	complexType               = "SoftLayer_Container_Product_Order_Network_Protection_Firewall_Dedicated"
	productPackageServiceMask = "description,prices.locationGroupId,prices.id"
	mandatoryFirewallType     = "FortiGate Security Appliance"
	multiVlansMask            = "id,customerManagedFlag,datacenter.name,bandwidthAllocation"
)

func resourceIBMNetworkMultiVlanCreate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	name := d.Get("name").(string)
	FirewallType := d.Get("firewall_type").(string)
	datacenter := d.Get("datacenter").(string)
	pod := d.Get("pod").(string)
	podName := datacenter + "." + pod
	PodService := services.GetNetworkPodService(sess)
	podMask := `frontendRouterId,name`

	// 1.Getting the router ID
	routerids, err := PodService.Filter(filter.Path("datacenterName").Eq(datacenter).Build()).Mask(podMask).GetAllObjects()
	if err != nil {
		return fmt.Errorf("Encountered problem trying to get the router ID: %s", err)
	}
	var routerid int
	for _, iterate := range routerids {
		if *iterate.Name == podName {
			routerid = *iterate.FrontendRouterId
		}
	}

	//2.Get the datacenter id
	dc, err := location.GetDatacenterByName(sess, datacenter, "id")
	if err != nil {
		return fmt.Errorf("Encountered problem trying to get the Datacenter ID: %s", err)
	}
	locationservice := services.GetLocationService(sess)

	//3. get the pricegroups that the datacenter belongs to
	priceidds, _ := locationservice.Id(*dc.Id).GetPriceGroups()
	var listofpriceids []int
	//store all the pricegroups a datacenter belongs to
	for _, priceidd := range priceidds {
		listofpriceids = append(listofpriceids, *priceidd.Id)
	}

	//4.get the addons that are specified
	var addonconfigurations []interface{}
	if _, ok := d.GetOk("addon_configuration"); ok {
		addonconfigurations, ok = d.Get("addon_configuration").([]interface{})
	}

	var actualaddons []string
	for _, addons := range addonconfigurations {
		actualaddons = append(actualaddons, addons.(string))
	}
	//appending the 20000GB Bandwidth item as it is mandatory
	actualaddons = append(actualaddons, FirewallType, "20000 GB Bandwidth Allotment")
	//appending the Fortigate Security Appliance as it is mandatory parameter for placing an order
	if FirewallType != mandatoryFirewallType {
		actualaddons = append(actualaddons, mandatoryFirewallType)
	}

	//5. Getting the priceids of items which have to be ordered
	priceItems := []datatypes.Product_Item_Price{}
	for _, addon := range actualaddons {
		actualpriceid, err := product.GetPriceIDByPackageIdandLocationGroups(sess, listofpriceids, 863, addon)
		if err != nil || actualpriceid == 0 {
			return fmt.Errorf("Encountered problem trying to get priceIds of items which have to be ordered: %s", err)
		}
		priceItem := datatypes.Product_Item_Price{
			Id: &actualpriceid,
		}
		priceItems = append(priceItems, priceItem)
	}

	//6.Get the package ID
	productpackageservice, _ := services.GetProductPackageService(sess).Filter(productPackageFilter).Mask(`id`).GetAllObjects()
	var productid int
	for _, packageid := range productpackageservice {
		productid = *packageid.Id
	}

	//7. Populate the container which needs to be sent for Verify order and Place order
	productOrderContainer := datatypes.Container_Product_Order_Network_Protection_Firewall_Dedicated{
		Container_Product_Order: datatypes.Container_Product_Order{
			PackageId:   &productid,
			Prices:      priceItems,
			Quantity:    sl.Int(1),
			Location:    &datacenter,
			ComplexType: sl.String(complexType),
		},
		Name:     sl.String(name),
		RouterId: &routerid,
	}

	//8.Calling verify order
	_, err = services.GetProductOrderService(sess.SetRetries(0)).
		VerifyOrder(&productOrderContainer)
	if err != nil {
		return fmt.Errorf("Error during Verify order for Creating: %s", err)
	}
	//9.Calling place order
	receipt, err := services.GetProductOrderService(sess.SetRetries(0)).
		PlaceOrder(&productOrderContainer, sl.Bool(false))
	if err != nil {
		return fmt.Errorf("Error during Place order for Creating: %s", err)
	}
	_, vlan, _, err := findDedicatedFirewallByOrderId(sess, *receipt.OrderId, d)
	if err != nil {
		return fmt.Errorf("Error during creation of dedicated hardware firewall: %s", err)
	}
	id := *vlan.NetworkFirewall.Id
	d.SetId(fmt.Sprintf("%d", id))
	log.Printf("[INFO] Firewall ID: %s", d.Id())
	return resourceIBMMultiVlanFirewallRead(d, meta)
}

func resourceIBMMultiVlanFirewallRead(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()

	fwID, _ := strconv.Atoi(d.Id())

	firewalls, err := services.GetAccountService(sess).
		Filter(filter.Build(
			filter.Path("networkGateways.networkFirewall.id").
				Eq(strconv.Itoa(fwID)))).
		Mask(multiVlanMask).
		GetNetworkGateways()
	if err != nil {
		return fmt.Errorf("Error retrieving firewall information: %s", err)
	}
	d.Set("datacenter", *firewalls[0].NetworkFirewall.Datacenter.Name)
	if *firewalls[0].NetworkFirewall.CustomerManagedFlag && *firewalls[0].MemberCount == 1 {
		d.Set("firewall_type", "FortiGate Security Appliance")
	} else if *firewalls[0].NetworkFirewall.CustomerManagedFlag && *firewalls[0].MemberCount > 1 {
		d.Set("firewall_type", "FortiGate Firewall Appliance HA Option")
	}
	addonConfiguration := make([]interface{}, 0, len(firewalls[0].NetworkFirewall.BillingItem.ActiveChildren))
	for _, elem := range firewalls[0].NetworkFirewall.BillingItem.ActiveChildren {
		if *elem.Description != "20000 GB Bandwidth Allotment" && *elem.Description != "FortiGate Firewall Appliance HA Option" {
			addonConfiguration = append(addonConfiguration, *elem.Description)
		}
	}
	if len(addonConfiguration) > 0 {
		d.Set("addon_configuration", addonConfiguration)
	}
	pod := *firewalls[0].NetworkFirewall.BillingItem.Notes
	pod = "pod" + strings.SplitAfter(pod, "pod")[1]
	d.Set("pod", &pod)
	d.Set("name", *firewalls[0].Name)
	d.Set("public_ip", *firewalls[0].PublicIpAddress.IpAddress)
	d.Set("public_ipv6", firewalls[0].PublicIpv6Address.IpAddress)
	d.Set("private_ip", *firewalls[0].PrivateIpAddress.IpAddress)
	d.Set("public_vlan_id", *firewalls[0].PublicVlan.Id)
	d.Set("private_vlan_id", *firewalls[0].PrivateVlan.Id)
	d.Set("username", *firewalls[0].NetworkFirewall.ManagementCredentials.Username)
	d.Set("password", *firewalls[0].NetworkFirewall.ManagementCredentials.Password)
	return nil
}

func resourceIBMMultiVlanFirewallUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("addon_configuration") {
		sess := meta.(ClientSession).SoftLayerSession()
		fwID, _ := strconv.Atoi(d.Id())
		old, new := d.GetChange("addon_configuration")
		oldaddons := old.([]interface{})
		newaddons := new.([]interface{})
		var oldaddon, newaddon, add []string
		for _, v := range oldaddons {
			oldaddon = append(oldaddon, v.(string))
		}
		for _, v := range newaddons {
			newaddon = append(newaddon, v.(string))
		}
		// 1. Remove old addons no longer appearing in the new set
		// 2. Add new addons not already provisioned
		remove := listdifference(oldaddon, newaddon)
		add = listdifference(newaddon, oldaddon)
		if len(remove) > 0 {
			firewalls, err := services.GetAccountService(sess).
				Filter(filter.Build(
					filter.Path("networkGateways.networkFirewall.id").
						Eq(strconv.Itoa(fwID)))).
				Mask(multiVlanMask).
				GetNetworkGateways()
			if err != nil {
				return fmt.Errorf("Some error occured while fetching the information of the Multi-Vlan Firewall")
			}
			for _, i := range remove {
				for _, j := range firewalls[0].NetworkFirewall.BillingItem.ActiveChildren {
					if i == *j.Description {
						cancelimmediately := true
						cancelAssociatedBillingItems := false
						reason := "No longer needed"
						customerNote := "No longer needed"
						billingitemservice, err := services.GetBillingItemService(sess).Id(*j.Id).CancelItem(&cancelimmediately, &cancelAssociatedBillingItems, &reason, &customerNote)
						if err != nil || !billingitemservice {
							return fmt.Errorf("Error while cancelling the addon")
						}
					}
				}
			}
		}
		if len(add) > 0 {
			datacentername, ok := d.GetOk("datacenter")
			if !ok {
				return fmt.Errorf("The attribute datacenter is not defined")
			}
			//2.Get the datacenter id
			dc, err := location.GetDatacenterByName(sess, datacentername.(string), "id")
			if err != nil {
				return fmt.Errorf("Datacenter not found")
			}
			locationservice := services.GetLocationService(sess)
			//3. get the pricegroups that the datacenter belongs to
			priceidds, _ := locationservice.Id(*dc.Id).GetPriceGroups()
			var listofpriceids []int
			//store all the pricegroups a datacenter belongs to
			for _, priceidd := range priceidds {
				listofpriceids = append(listofpriceids, *priceidd.Id)
			}
			priceItems := []datatypes.Product_Item_Price{}
			for _, addon := range add {
				actualpriceid, err := product.GetPriceIDByPackageIdandLocationGroups(sess, listofpriceids, 863, addon)
				if err != nil || actualpriceid == 0 {
					return fmt.Errorf("The addon or the firewall is not available for the datacenter you have selected. Please enter a different datacenter")
				}
				priceItem := datatypes.Product_Item_Price{
					Id: &actualpriceid,
				}
				priceItems = append(priceItems, priceItem)
			}
			//6.Get the package ID
			productpackageservice, _ := services.GetProductPackageService(sess).Filter(productPackageFilter).Mask(`id`).GetAllObjects()
			var productid int
			for _, packageid := range productpackageservice {
				productid = *packageid.Id
			}
			var properties []datatypes.Container_Product_Order_Property
			t := time.Now()
			upgradeproductOrderContainer := datatypes.Container_Product_Order_Network_Protection_Firewall_Dedicated_Upgrade{
				Container_Product_Order_Network_Protection_Firewall_Dedicated: datatypes.Container_Product_Order_Network_Protection_Firewall_Dedicated{
					Container_Product_Order: datatypes.Container_Product_Order{
						PackageId:   &productid,
						Prices:      priceItems,
						ComplexType: sl.String(complexType),
						Properties: append(properties, datatypes.Container_Product_Order_Property{
							Name:  sl.String("MAINTENANCE_WINDOW"),
							Value: sl.String(time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, t.Location()).UTC().String()),
						}),
					},
				},
				FirewallId: &fwID,
			}
			//8.Calling verify order
			_, err = services.GetProductOrderService(sess.SetRetries(0)).
				VerifyOrder(&upgradeproductOrderContainer)
			if err != nil {
				return fmt.Errorf("Error during Verify order for Updating: %s", err)
			}

			//9.Calling place order
			receipt, err := services.GetProductOrderService(sess.SetRetries(0)).
				PlaceOrder(&upgradeproductOrderContainer, sl.Bool(false))
			if err != nil {
				return fmt.Errorf("Error during Place order for Updating: %s", err)
			}
			_, _, _, err = findDedicatedFirewallByOrderId(sess, *receipt.OrderId, d)
			if err != nil {
				return fmt.Errorf("Error during creation of dedicated hardware firewall: %s", err)
			}
		}
	}
	return resourceIBMMultiVlanFirewallRead(d, meta)
}

func resourceIBMMultiVLanFirewallExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess := meta.(ClientSession).SoftLayerSession()

	fwID, _ := strconv.Atoi(d.Id())

	firewalls, err := services.GetAccountService(sess).
		Filter(filter.Build(
			filter.Path("networkGateways.networkFirewall.id").
				Eq(strconv.Itoa(fwID)))).
		Mask(multiVlanMask).
		GetNetworkGateways()
	if err != nil {
		return false, fmt.Errorf("Error retrieving firewall information: %s", err)
	}
	if firewalls[0].NetworkFirewall.BillingItem == nil {
		return false, nil
	}
	return true, nil
}

//This function takes two lists and returns the difference between the two lists
//listdifference([1,2] [2,3]) = [1]
func listdifference(a, b []string) []string {
	mb := map[string]bool{}
	for _, x := range b {
		mb[x] = true
	}
	ab := []string{}
	for _, x := range a {
		if _, ok := mb[x]; !ok {
			ab = append(ab, x)
		}
	}
	return ab
}
