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

func resourceIBMStorageEvault() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMStorageEvaultCreate,
		Read:     resourceIBMStorageEvaultRead,
		Update:   resourceIBMStorageEvaultUpdate,
		Delete:   resourceIBMStorageEvaultDelete,
		Exists:   resourceIBMStorageEvaultExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"datacenter": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Datacenter name",
			},
			"capacity": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Capacity",
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
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "user name",
			},
			"password": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "password",
			},
			"service_resource_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "service resource name",
			},
			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "Tags set for the resource",
			},
		},
	}
}

const (
	evaultPackageFilter = `{"keyName":{"operation":"ADDITIONAL_PRODUCTS"}}`
	evaultStorageMask   = "id,billingItem.orderItem.order.id"
)

func resourceIBMStorageEvaultCreate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()

	// Find price items
	productOrderContainer, err := buildEvaultProductOrderContainer(d, sess)
	if err != nil {
		return fmt.Errorf("Error creating evault: %s", err)
	}

	log.Println("[INFO] Creating Evault")

	receipt, err := services.GetProductOrderService(sess.SetRetries(0)).
		PlaceOrder(productOrderContainer, sl.Bool(false))
	if err != nil {
		return fmt.Errorf("Error during creation of evault: %s", err)
	}
	evaultStorage, err := findEvaultStorageByOrderID(d, meta, *receipt.OrderId)

	if err != nil {
		return fmt.Errorf("Error during creation of storage: %s", err)
	}

	d.SetId(fmt.Sprintf("%d", *evaultStorage.Id))

	// Wait for storage availability
	_, err = WaitForEvaultAvailable(d, meta, schema.TimeoutCreate)

	if err != nil {
		return fmt.Errorf(
			"Error waiting for evault (%s) to become ready: %s", d.Id(), err)
	}

	// SoftLayer changes the device ID after completion of provisioning. It is necessary to refresh device ID.
	evaultStorage, err = findEvaultStorageByOrderID(d, meta, *receipt.OrderId)

	if err != nil {
		return fmt.Errorf("Error during creation of storage: %s", err)
	}
	d.SetId(fmt.Sprintf("%d", *evaultStorage.Id))

	log.Printf("[INFO] Storage ID: %s", d.Id())

	return resourceIBMStorageEvaultRead(d, meta)
}

func buildEvaultProductOrderContainer(d *schema.ResourceData, sess *session.Session) (
	*datatypes.Container_Product_Order, error) {
	datacenter := d.Get("datacenter").(string)
	capacity := d.Get("capacity").(int)

	var virtualID, hardwareID int
	if vID, ok := d.GetOk("virtual_instance_id"); ok {
		virtualID = vID.(int)
	}

	if hID, ok := d.GetOk("hardware_instance_id"); ok {
		hardwareID = hID.(int)
	}

	if virtualID == 0 && hardwareID == 0 {
		return &datatypes.Container_Product_Order{}, fmt.Errorf("Provide either `virtual_instance_id` or `hardware_instance_id`")
	}

	/*pkg, err := product.GetPackageByType(sess, "ADDITIONAL_PRODUCTS")
	if err != nil {
		return nil, err
	}*/
	productpackageservice, _ := services.GetProductPackageService(sess).Filter(evaultPackageFilter).Mask(`id`).GetAllObjects()
	var productid int
	for _, packageid := range productpackageservice {
		productid = *packageid.Id
	}

	// Lookup the data center ID
	dc, err := location.GetDatacenterByName(sess, datacenter)
	if err != nil {
		return &datatypes.Container_Product_Order{},
			fmt.Errorf("No data centers matching %s could be found", datacenter)
	}

	locationservice := services.GetLocationService(sess)

	//3. get the pricegroups that the datacenter belongs to
	priceidds, _ := locationservice.Id(*dc.Id).GetPriceGroups()

	var listofpriceids []int

	//store all the pricegroups a datacenter belongs to
	for _, priceidd := range priceidds {
		listofpriceids = append(listofpriceids, *priceidd.Id)
	}

	description := strconv.Itoa(capacity) + "GB IBM Cloud Backup"

	priceItems := []datatypes.Product_Item_Price{}
	actualpriceid, err := product.GetPriceIDByPackageIdandLocationGroups(sess, listofpriceids, 0, description)
	if err != nil || actualpriceid == 0 {
		return &datatypes.Container_Product_Order{}, fmt.Errorf("The evault with the given capacity is not available for the datacenter you have selected. Please enter a different capacity : %s", err)
	}
	priceItem := datatypes.Product_Item_Price{
		Id: &actualpriceid,
	}
	priceItems = append(priceItems, priceItem)

	order := datatypes.Container_Product_Order{
		ComplexType: sl.String("SoftLayer_Container_Product_Order_Network_Storage_Backup_Evault_Vault"),
		PackageId:   &productid,
		Prices:      priceItems,
		Location:    sl.String(strconv.Itoa(*dc.Id)),
	}

	if virtualID > 0 {
		var guest datatypes.Virtual_Guest
		guest.Id = sl.Int(virtualID)
		order.VirtualGuests = []datatypes.Virtual_Guest{
			guest,
		}

	} else {
		var hardware datatypes.Hardware
		hardware.Id = sl.Int(hardwareID)
		order.Hardware = []datatypes.Hardware{
			hardware,
		}

	}

	return &order, nil
}

func resourceIBMStorageEvaultRead(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()

	evaultID, _ := strconv.Atoi(d.Id())

	evault, err := services.GetNetworkStorageBackupEvaultService(sess).
		Id(evaultID).Mask("billingItem[location[name]]").
		GetObject()

	if err != nil {
		return fmt.Errorf("Error retrieving evault information: %s", err)
	}

	d.Set("capacity", evault.CapacityGb)
	d.Set("datacenter", evault.BillingItem.Location.Name)

	if evault.GuestId != nil {
		d.Set("virtual_instance_id", evault.GuestId)
	}

	if evault.HardwareId != nil {
		d.Set("hardware_instance_id", evault.HardwareId)
	}

	d.Set("username", evault.Username)
	d.Set("password", evault.Password)
	d.Set("service_resource_name", evault.ServiceResourceName)

	return nil
}

func resourceIBMStorageEvaultUpdate(d *schema.ResourceData, meta interface{}) error {

	if d.HasChange("capacity") && !d.IsNewResource() {
		sess := meta.(ClientSession).SoftLayerSession()

		evaultID, err := strconv.Atoi(d.Id())
		if err != nil {
			return fmt.Errorf("Not a valid ID, must be an integer: %s", err)
		}

		priceID, err := getEvaultUpgradePriceItem(d, sess)
		if err != nil {
			return err
		}

		_, err = services.GetNetworkStorageBackupEvaultService(sess).
			Id(evaultID).UpgradeVolumeCapacity(sl.Int(priceID))

		if err != nil {
			return err
		}

		// Wait for storage availability
		_, err = WaitForEvaultAvailable(d, meta, schema.TimeoutUpdate)

		if err != nil {
			return fmt.Errorf(
				"Error waiting for evault upgrade (%s) to become ready: %s", d.Id(), err)
		}

		return resourceIBMStorageEvaultRead(d, meta)
	}
	return nil

}

func resourceIBMStorageEvaultDelete(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	evaultService := services.GetNetworkStorageBackupEvaultService(sess)
	evaultID, _ := strconv.Atoi(d.Id())

	// Get billing item associated with the storage
	billingItem, err := evaultService.Id(evaultID).GetBillingItem()

	if err != nil {
		return fmt.Errorf("Error while looking up billing item associated with the evault: %s", err)
	}

	if billingItem.Id == nil {
		return fmt.Errorf("Error while looking up billing item associated with the evault: No billing item for ID:%d", evaultID)
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

func findEvaultStorageByOrderID(d *schema.ResourceData, meta interface{}, orderId int) (datatypes.Network_Storage, error) {
	filterPath := "evaultNetworkStorage.billingItem.orderItem.order.id"
	sess := meta.(ClientSession).SoftLayerSession()

	stateConf := &resource.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"complete"},
		Refresh: func() (interface{}, string, error) {
			storages, err := services.GetAccountService(sess).
				Filter(filter.Build(
					filter.Path(filterPath).
						Eq(strconv.Itoa(orderId)))).
				Mask(storageMask).
				GetEvaultNetworkStorage()
			if err != nil {
				return datatypes.Network_Storage{}, "", err
			}

			if len(storages) == 1 {
				return storages[0], "complete", nil
			} else if len(storages) == 0 {
				return nil, "pending", nil
			} else {
				return nil, "", fmt.Errorf("Expected one evault: %s", err)
			}

		},
		Timeout:        d.Timeout(schema.TimeoutCreate),
		Delay:          10 * time.Second,
		MinTimeout:     10 * time.Second,
		NotFoundChecks: 300,
	}

	pendingResult, err := stateConf.WaitForState()

	if err != nil {
		return datatypes.Network_Storage{}, err
	}

	var result, ok = pendingResult.(datatypes.Network_Storage)

	if ok {
		return result, nil
	}

	return datatypes.Network_Storage{},
		fmt.Errorf("Cannot find evault with order id '%d'", orderId)
}

// Waits for storage provisioning
func WaitForEvaultAvailable(d *schema.ResourceData, meta interface{}, timeout string) (interface{}, error) {
	log.Printf("Waiting for evault (%s) to be available.", d.Id())
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return nil, fmt.Errorf("The evault ID %s must be numeric", d.Id())
	}
	sess := meta.(ClientSession).SoftLayerSession()
	stateConf := &resource.StateChangeConf{
		Pending: []string{"retry", "provisioning"},
		Target:  []string{"available"},
		Refresh: func() (interface{}, string, error) {
			// Check active transactions
			service := services.GetNetworkStorageBackupEvaultService(sess)
			result, err := service.Id(id).Mask("activeTransactionCount").GetObject()
			if err != nil {
				if apiErr, ok := err.(sl.Error); ok && apiErr.StatusCode == 404 {
					return nil, "", fmt.Errorf("Error retrieving evault: %s", err)
				}
				return false, "retry", nil
			}

			log.Println("Checking active transactions.")
			if *result.ActiveTransactionCount > 0 {
				return result, "provisioning", nil
			}

			return result, "available", nil
		},
		Timeout:    d.Timeout(timeout),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func resourceIBMStorageEvaultExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess := meta.(ClientSession).SoftLayerSession()

	evaultID, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	_, err = services.GetNetworkStorageBackupEvaultService(sess).
		Id(evaultID).
		GetObject()

	if err != nil {
		if apiErr, ok := err.(sl.Error); ok && apiErr.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("Error retrieving evault information: %s", err)
	}
	return true, nil
}

func getEvaultUpgradePriceItem(d *schema.ResourceData, sess *session.Session) (int, error) {
	evaultID, _ := strconv.Atoi(d.Id())

	evault, err := services.GetNetworkStorageBackupEvaultService(sess).
		Id(evaultID).Mask("id, billingItem[id,upgradeItems[prices]]").
		GetObject()

	if err != nil {
		return 0, fmt.Errorf("Error retrieving evault information: %s", err)
	}

	capacity := d.Get("capacity")

	len := len(evault.BillingItem.UpgradeItems)
	validCapacities := make([]int, len)

	for i, item := range evault.BillingItem.UpgradeItems {
		if int(*item.Capacity) == capacity.(int) {
			return *item.Id, nil
		}

		validCapacities[i] = int(*item.Capacity)
	}

	return 0, fmt.Errorf("The given capacity is not a valid upgrade value. Valid capacity upgrades are: %d", validCapacities)

}
