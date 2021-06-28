// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"strconv"

	"log"
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
	FwHardwareDedicatedPackageType = "ADDITIONAL_SERVICES_FIREWALL"

	vlanMask = "firewallNetworkComponents,networkVlanFirewall.billingItem.orderItem.order.id,dedicatedFirewallFlag" +
		",firewallGuestNetworkComponents,firewallInterfaces,firewallRules,highAvailabilityFirewallFlag"
	fwMask        = "id,datacenter,primaryIpAddress,networkVlan.highAvailabilityFirewallFlag,managementCredentials,tagReferences[id,tag[name]]"
	multiVlanMask = "id,name,networkFirewall[id,customerManagedFlag,datacenter.name,billingItem[orderItem.order.id,activeChildren[categoryCode, description,id]],managementCredentials,firewallType],publicIpAddress.ipAddress,publicIpv6Address.ipAddress,publicVlan[id,primaryRouter.hostname],privateVlan[id,primaryRouter.hostname],privateIpAddress.ipAddress,insideVlans[id],memberCount,status.keyName"
)

func resourceIBMFirewall() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMFirewallCreate,
		Read:     resourceIBMFirewallRead,
		Update:   resourceIBMFirewallUpdate,
		Delete:   resourceIBMFirewallDelete,
		Exists:   resourceIBMFirewallExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"firewall_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "HARDWARE_FIREWALL_DEDICATED",
				ValidateFunc: validateAllowedStringValue([]string{
					"HARDWARE_FIREWALL_DEDICATED",
					"FORTIGATE_SECURITY_APPLIANCE",
				}),
				Description: "Firewall type",
			},

			"ha_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "set to true if High availability is enabled",
			},
			"public_vlan_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Public VLAN ID",
			},
			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "List of tags for the firewall",
			},
			"location": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Location info",
			},
			"primary_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Primary IP address",
			},
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User name",
			},
			"password": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Password for the given User",
			},
		},
	}
}

func resourceIBMFirewallCreate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()

	keyName := "HARDWARE_FIREWALL_DEDICATED"
	firewallType := d.Get("firewall_type").(string)
	haEnabled := d.Get("ha_enabled").(bool)
	if haEnabled {
		if firewallType == "HARDWARE_FIREWALL_DEDICATED" {
			keyName = "HARDWARE_FIREWALL_HIGH_AVAILABILITY"
		} else {
			keyName = "FORTIGATE_SECURITY_APPLIANCE_HIGH_AVAILABILITY"
		}
	} else {
		keyName = firewallType
	}

	publicVlanId := d.Get("public_vlan_id").(int)

	pkg, err := product.GetPackageByType(sess, FwHardwareDedicatedPackageType)
	if err != nil {
		return err
	}

	// Get all prices for ADDITIONAL_SERVICES_FIREWALL with the given capacity
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

	productOrderContainer := datatypes.Container_Product_Order_Network_Protection_Firewall_Dedicated{
		Container_Product_Order: datatypes.Container_Product_Order{
			PackageId: pkg.Id,
			Prices: []datatypes.Product_Item_Price{
				{
					Id: targetItems[0].Prices[0].Id,
				},
			},
			Quantity: sl.Int(1),
		},
		VlanId: sl.Int(publicVlanId),
	}

	log.Println("[INFO] Creating dedicated hardware firewall")

	receipt, err := services.GetProductOrderService(sess.SetRetries(0)).
		PlaceOrder(&productOrderContainer, sl.Bool(false))
	if err != nil {
		return fmt.Errorf("Error during creation of dedicated hardware firewall: %s", err)
	}
	vlan, _, _, err := findDedicatedFirewallByOrderId(sess, *receipt.OrderId, d)
	if err != nil {
		return fmt.Errorf("Error during creation of dedicated hardware firewall: %s", err)
	}

	id := *vlan.NetworkVlanFirewall.Id
	d.SetId(fmt.Sprintf("%d", id))
	d.Set("ha_enabled", *vlan.HighAvailabilityFirewallFlag)
	d.Set("public_vlan_id", *vlan.Id)

	log.Printf("[INFO] Firewall ID: %s", d.Id())

	// Set tags
	tags := getTags(d)
	if tags != "" {
		//Try setting only when it is non empty as we are creating Firewall
		err = setFirewallTags(id, tags, meta)
		if err != nil {
			return err
		}
	}

	return resourceIBMFirewallRead(d, meta)
}

func resourceIBMFirewallRead(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()

	fwID, _ := strconv.Atoi(d.Id())

	fw, err := services.GetNetworkVlanFirewallService(sess).
		Id(fwID).
		Mask(fwMask).
		GetObject()

	if err != nil {
		return fmt.Errorf("Error retrieving firewall information: %s", err)
	}

	d.Set("public_vlan_id", *fw.NetworkVlan.Id)
	d.Set("ha_enabled", *fw.NetworkVlan.HighAvailabilityFirewallFlag)
	d.Set("location", *fw.Datacenter.Name)
	d.Set("primary_ip", *fw.PrimaryIpAddress)
	if fw.ManagementCredentials != nil {
		d.Set("username", *fw.ManagementCredentials.Username)
		d.Set("password", *fw.ManagementCredentials.Password)
	}

	tagRefs := fw.TagReferences
	tagRefsLen := len(tagRefs)
	if tagRefsLen > 0 {
		tags := make([]string, tagRefsLen, tagRefsLen)
		for i, tagRef := range tagRefs {
			tags[i] = *tagRef.Tag.Name
		}
		d.Set("tags", tags)
	}

	return nil
}

func resourceIBMFirewallUpdate(d *schema.ResourceData, meta interface{}) error {

	fwID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Not a valid firewall ID, must be an integer: %s", err)
	}

	// Update tags
	if d.HasChange("tags") {
		tags := getTags(d)
		err := setFirewallTags(fwID, tags, meta)
		if err != nil {
			return err
		}
	}
	return resourceIBMFirewallRead(d, meta)
}

func resourceIBMFirewallDelete(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	fwService := services.GetNetworkVlanFirewallService(sess)

	fwID, _ := strconv.Atoi(d.Id())

	// Get billing item associated with the firewall
	billingItem, err := fwService.Id(fwID).GetBillingItem()

	if err != nil {
		return fmt.Errorf("Error while looking up billing item associated with the firewall: %s", err)
	}

	if billingItem.Id == nil {
		return fmt.Errorf("Error while looking up billing item associated with the firewall: No billing item for ID:%d", fwID)
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

func resourceIBMFirewallExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess := meta.(ClientSession).SoftLayerSession()

	fwID, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	_, err = services.GetNetworkVlanFirewallService(sess).
		Id(fwID).
		GetObject()

	if err != nil {
		if apiErr, ok := err.(sl.Error); ok && apiErr.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("Error retrieving firewall information: %s", err)
	}

	return true, nil
}

func findDedicatedFirewallByOrderId(sess *session.Session, orderId int, d *schema.ResourceData) (datatypes.Network_Vlan, datatypes.Network_Gateway, datatypes.Product_Upgrade_Request, error) {
	filterPath := "networkVlans.networkVlanFirewall.billingItem.orderItem.order.id"
	multivlanfilterpath := "networkGateways.networkFirewall.billingItem.orderItem.order.id"
	var vlans []datatypes.Network_Vlan
	var err error
	var firewalls []datatypes.Network_Gateway
	var upgraderequest datatypes.Product_Upgrade_Request
	stateConf := &resource.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"complete"},
		Refresh: func() (interface{}, string, error) {

			if ok := d.HasChange("addon_configuration") && !d.IsNewResource(); ok {
				fwID, _ := strconv.Atoi(d.Id())
				upgraderequest, err = services.GetNetworkVlanFirewallService(sess).
					Id(fwID).
					Mask("status").
					GetUpgradeRequest()
				if err != nil {
					return datatypes.Product_Upgrade_Request{}, "", err
				}
			} else if _, ok := d.GetOk("pod"); ok {
				firewalls, err = services.GetAccountService(sess).
					Filter(filter.Build(
						filter.Path(multivlanfilterpath).
							Eq(strconv.Itoa(orderId)))).
					Mask(multiVlanMask).
					GetNetworkGateways()
				if err != nil {
					return datatypes.Network_Gateway{}, "", err
				}
			} else {
				vlans, err = services.GetAccountService(sess).
					Filter(filter.Build(
						filter.Path(filterPath).
							Eq(strconv.Itoa(orderId)))).
					Mask(vlanMask).
					GetNetworkVlans()
				if err != nil {
					return datatypes.Network_Vlan{}, "", err
				}
			}
			status, ok := sl.GrabOk(upgraderequest, "Status.Name")
			if ok && status == "Complete" {
				return upgraderequest, "complete", nil
			} else if len(vlans) == 1 {
				return vlans[0], "complete", nil
			} else if len(firewalls) == 1 {
				return firewalls[0], "complete", nil
			} else if len(vlans) == 0 || len(firewalls) == 0 || *upgraderequest.Status.Name != "Complete" {
				return datatypes.Network_Vlan{}, "pending", nil
			}
			return nil, "", fmt.Errorf("Expected one dedicated firewall: %s", err)
		},
		Timeout:        2 * time.Hour,
		Delay:          10 * time.Second,
		MinTimeout:     10 * time.Second,
		NotFoundChecks: 24 * 60,
	}

	pendingResult, err := stateConf.WaitForState()

	if err != nil {
		return datatypes.Network_Vlan{}, datatypes.Network_Gateway{}, datatypes.Product_Upgrade_Request{}, err
	}
	if ok := d.HasChange("addon_configuration") && !d.IsNewResource(); ok {
		if result, ok := pendingResult.(datatypes.Product_Upgrade_Request); ok {
			return datatypes.Network_Vlan{}, datatypes.Network_Gateway{}, result, nil
		}
		return datatypes.Network_Vlan{}, datatypes.Network_Gateway{}, datatypes.Product_Upgrade_Request{},
			fmt.Errorf("Something went wrong while upgrading '%d'", orderId)
	} else if _, ok := d.GetOk("pod"); ok {
		if result, ok := pendingResult.(datatypes.Network_Gateway); ok {
			return datatypes.Network_Vlan{}, result, datatypes.Product_Upgrade_Request{}, nil
		}
		return datatypes.Network_Vlan{}, datatypes.Network_Gateway{}, datatypes.Product_Upgrade_Request{},
			fmt.Errorf("Cannot find Dedicated Firewall with order id '%d'", orderId)
	}
	var result, ok = pendingResult.(datatypes.Network_Vlan)

	if ok {
		return result, datatypes.Network_Gateway{}, datatypes.Product_Upgrade_Request{}, nil
	}

	return datatypes.Network_Vlan{}, datatypes.Network_Gateway{}, datatypes.Product_Upgrade_Request{},
		fmt.Errorf("Cannot find Dedicated Firewall with order id '%d'", orderId)
}

func setFirewallTags(id int, tags string, meta interface{}) error {
	service := services.GetNetworkVlanFirewallService(meta.(ClientSession).SoftLayerSession())
	_, err := service.Id(id).SetTags(sl.String(tags))
	if err != nil {
		return fmt.Errorf("Could not set tags on firewall %d", id)
	}
	return nil
}
