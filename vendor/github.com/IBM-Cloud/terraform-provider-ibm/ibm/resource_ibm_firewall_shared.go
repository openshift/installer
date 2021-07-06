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
	"github.com/softlayer/softlayer-go/helpers/product"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/sl"
)

const (
	FwHardwarePackageType = "ADDITIONAL_SERVICES_FIREWALL"
)

func resourceIBMFirewallShared() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMFirewallSharedCreate,
		Read:     resourceIBMFirewallSharedRead,
		Delete:   resourceIBMFirewallSharedDelete,
		Exists:   resourceIBMFirewallSharedExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"billing_item_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Billing Item ID",
			},
			"firewall_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue([]string{"10MBPS_HARDWARE_FIREWALL", "20MBPS_HARDWARE_FIREWALL", "100MBPS_HARDWARE_FIREWALL", "1000MBPS_HARDWARE_FIREWALL", "200MBPS_HARDWARE_FIREWALL", "2000MBPS_HARDWARE_FIREWALL"}),
				Description:  "Firewall type",
			},
			"virtual_instance_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"hardware_instance_id"},
				Description:   "Virtual instance ID",
			},
			"hardware_instance_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"virtual_instance_id"},
				Description:   "Hardware instance ID",
			},
		},
	}
}

// keyName is in between:[10MBPS_HARDWARE_FIREWALL, 20MBPS_HARDWARE_FIREWALL,
//                         100MBPS_HARDWARE_FIREWALL, 1000MBPS_HARDWARE_FIREWALL]
func resourceIBMFirewallSharedCreate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()

	keyName := d.Get("firewall_type").(string)

	var virtualId, hardwareId int
	if vID, ok := d.GetOk("virtual_instance_id"); ok {
		virtualId = vID.(int)
	}

	if hID, ok := d.GetOk("hardware_instance_id"); ok {
		hardwareId = hID.(int)
	}

	if virtualId == 0 && hardwareId == 0 {
		return fmt.Errorf("Provide either `virtual_instance_id` or `hardware_instance_id`")
	}

	//var productOrderContainer *string
	pkg, err := product.GetPackageByType(sess, FwHardwarePackageType)
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

	masked := "id,firewallServiceComponent[id,status]"
	if virtualId > 0 {
		productOrderContainer := datatypes.Container_Product_Order_Network_Protection_Firewall{
			Container_Product_Order: datatypes.Container_Product_Order{
				PackageId: pkg.Id,
				Prices: []datatypes.Product_Item_Price{
					{
						Id: targetItems[0].Prices[0].Id,
					},
				},
				Quantity: sl.Int(1),
				VirtualGuests: []datatypes.Virtual_Guest{{
					Id: sl.Int(virtualId),
				},
				},
			},
		}
		_, err := services.GetProductOrderService(sess.SetRetries(0)).PlaceOrder(&productOrderContainer, sl.Bool(false))
		if err != nil {
			return nil
		}

		log.Printf("[INFO] Wait one minute before fetching the firewall/device.")
		time.Sleep(time.Second * 30)
		service := services.GetVirtualGuestService(sess)
		stateConf := &resource.StateChangeConf{
			Target:  []string{"completed"},
			Pending: []string{"pending"},
			Refresh: func() (interface{}, string, error) {
				result, err := service.Id(virtualId).Mask(masked).GetObject()
				if err != nil {
					return nil, "", err
				}
				status, ok := sl.GrabOk(result, "FirewallServiceComponent.Status")
				if ok && status == "bypass" {
					return result, "completed", nil
				}
				return result, "pending", nil
			},
			Timeout:        d.Timeout(schema.TimeoutCreate),
			Delay:          10 * time.Second,
			MinTimeout:     10 * time.Second,
			NotFoundChecks: 24 * 60,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return err
		}

		result, err := service.Id(virtualId).Mask(masked).GetObject()
		idd := *result.FirewallServiceComponent.Id
		log.Print(idd)
		d.SetId(fmt.Sprintf("%d", idd))

		if err != nil {
			return fmt.Errorf("Error during creation of hardware firewall: %s", err)
		}

	}
	if hardwareId > 0 {
		productOrderContainer := datatypes.Container_Product_Order_Network_Protection_Firewall{
			Container_Product_Order: datatypes.Container_Product_Order{
				PackageId: pkg.Id,
				Prices: []datatypes.Product_Item_Price{
					{
						Id: targetItems[0].Prices[0].Id,
					},
				},
				Quantity: sl.Int(1),
				Hardware: []datatypes.Hardware{{
					Id: sl.Int(hardwareId),
				},
				},
			},
		}
		_, err := services.GetProductOrderService(sess.SetRetries(0)).PlaceOrder(&productOrderContainer, sl.Bool(false))
		if err != nil {
			return nil
		}

		log.Printf("[INFO] Wait one minute before fetching the firewall/device.")
		time.Sleep(time.Second * 30)

		service := services.GetHardwareService(sess)
		stateConf := &resource.StateChangeConf{
			Target:  []string{"completed"},
			Pending: []string{"pending"},
			Refresh: func() (interface{}, string, error) {
				result, err := service.Id(hardwareId).Mask(masked).GetObject()
				if err != nil {
					return nil, "", err
				}
				status, ok := sl.GrabOk(result, "FirewallServiceComponent.Status")
				if ok && status == "bypass" {
					return result, "completed", nil
				}
				return result, "pending", nil
			},
			Timeout:        d.Timeout(schema.TimeoutCreate),
			Delay:          10 * time.Second,
			MinTimeout:     10 * time.Second,
			NotFoundChecks: 24 * 60,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return err
		}

		resultNew, err := service.Id(hardwareId).Mask(masked).GetObject()
		idd2 := *resultNew.FirewallServiceComponent.Id

		d.SetId(fmt.Sprintf("%d", idd2))
		log.Print(idd2)
		if err != nil {
			return fmt.Errorf("Error during creation of hardware firewall: %s", err)
		}

	}
	log.Println("[INFO] Creating hardware firewall shared")

	return resourceIBMFirewallSharedRead(d, meta)
}

func resourceIBMFirewallSharedRead(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()

	firewall_type := (d.Get("firewall_type").(string))
	d.Set("firewall_type", firewall_type)

	fservice := services.GetNetworkComponentFirewallService(sess)

	fwID, _ := strconv.Atoi(d.Id())

	data, err := fservice.Id(fwID).Mask("billingItem.id").GetObject()
	d.Set("billing_item_id", *data.BillingItem.Id)
	if err != nil {
		return fmt.Errorf("Error during creation of hardware firewall: %s", err)
	}

	return nil
}

//detach hardware firewall from particular machine
func resourceIBMFirewallSharedDelete(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	idd2 := (d.Get("billing_item_id")).(int)

	success, err := services.GetBillingItemService(sess).Id(idd2).CancelService()
	log.Print(success)
	if err != nil {
		return err
	}

	if !success {
		return fmt.Errorf("SoftLayer reported an unsuccessful cancellation")
	}
	return nil
}

//exists method
func resourceIBMFirewallSharedExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess := meta.(ClientSession).SoftLayerSession()
	fservice := services.GetNetworkComponentFirewallService(sess)
	id, err := strconv.Atoi(d.Id())
	response, err := fservice.Id(id).GetObject()

	if err != nil {
		log.Printf("error fetching the firewall resource: %s", err)
		return false, err
	}
	log.Print(response)
	return true, nil
}
