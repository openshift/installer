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
	"github.com/softlayer/softlayer-go/helpers/hardware"
	"github.com/softlayer/softlayer-go/helpers/location"
	"github.com/softlayer/softlayer-go/helpers/product"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/sl"
)

var dedicatedHostPackageType = "DEDICATED_HOST"

func resourceIBMComputeDedicatedHost() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMComputeDedicatedHostCreate,
		Read:     resourceIBMComputeDedicatedHostRead,
		Delete:   resourceIBMComputeDedicatedHostDelete,
		Exists:   resourceIBMComputeDedicatedHostExists,
		Update:   resourceIBMComputeDedicatedHostUpdate,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The host name of dedicatated host.",
			},
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The domain of dedicatated host.",
			},
			"datacenter": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The data center in which the dedicatated host is to be provisioned.",
			},
			"flavor": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "56_CORES_X_242_RAM_X_1_4_TB",
				ForceNew:    true,
				Description: "The flavor of the dedicatated host.",
			},
			"hourly_billing": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				ForceNew:    true,
				Description: "The billing type for the dedicatated host.",
			},
			"router_hostname": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The hostname of the primary router that the dedicated host is associated with.",
			},
			"cpu_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The capacity that the dedicated host's CPU allocation is restricted to.",
			},
			"disk_capacity": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The capacity that the dedicated host's disk allocation is restricted to.",
			},
			"memory_capacity": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The capacity that the dedicated host's memory allocation is restricted to.",
			},
			"wait_time_minutes": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  90,
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func resourceIBMComputeDedicatedHostCreate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()

	pkg, err := product.GetPackageByType(sess, dedicatedHostPackageType)
	if err != nil {
		return err
	}

	datacenter := d.Get("datacenter").(string)
	router := d.Get("router_hostname").(string)
	flavor := d.Get("flavor").(string)

	// Lookup the data center ID
	dc, err := location.GetDatacenterByName(sess, datacenter)
	if err != nil {
		return fmt.Errorf("No data centers matching %s could be found", datacenter)
	}

	rt, err := hardware.GetRouterByName(sess, router, "id")
	if err != nil {
		return fmt.Errorf("Error creating dedicated host: %s", err)
	}

	primaryBackendNetworkComponent := datatypes.Network_Component{
		Router: &datatypes.Hardware{
			Id: rt.Id,
		},
	}

	hardware := datatypes.Hardware{
		Hostname:                       sl.String(d.Get("hostname").(string)),
		Domain:                         sl.String(d.Get("domain").(string)),
		HourlyBillingFlag:              sl.Bool(d.Get("hourly_billing").(bool)),
		PrimaryBackendNetworkComponent: &primaryBackendNetworkComponent,
	}

	// 2. Get all prices for the package
	productItems, err := product.GetPackageProducts(sess, *pkg.Id, productItemMaskWithPriceLocationGroupID)
	if err != nil {
		return err
	}

	priceItems := []datatypes.Product_Item_Price{}
	for _, item := range productItems {
		if *item.KeyName == flavor {
			for _, price := range item.Prices {
				if price.LocationGroupId == nil {
					priceItem := datatypes.Product_Item_Price{
						Id: price.Id,
					}
					priceItems = append(priceItems, priceItem)
					break
				}
			}

		}

	}

	productOrderContainer := datatypes.Container_Product_Order_Virtual_DedicatedHost{
		Container_Product_Order: datatypes.Container_Product_Order{
			PackageId:        pkg.Id,
			Location:         sl.String(strconv.Itoa(*dc.Id)),
			Prices:           priceItems,
			Quantity:         sl.Int(1),
			UseHourlyPricing: sl.Bool(true),
		},
	}

	productOrderContainer.Hardware = make([]datatypes.Hardware, 0, 1)
	productOrderContainer.Hardware = append(
		productOrderContainer.Hardware,
		hardware,
	)

	log.Println("[INFO] Creating dedicated host")

	//verify order
	_, err = services.GetProductOrderService(sess.SetRetries(0)).
		VerifyOrder(&productOrderContainer)
	if err != nil {
		return fmt.Errorf("Error during creation of dedicated host: %s", err)
	}
	//place order
	_, err = services.GetProductOrderService(sess.SetRetries(0)).
		PlaceOrder(&productOrderContainer, sl.Bool(false))
	if err != nil {
		return fmt.Errorf("Error during creation of dedicated host: %s", err)
	}

	// wait for machine availability
	dedicated, err := findDedicatedHostByOrderID(&hardware, d, meta)
	if err != nil {
		return fmt.Errorf(
			"Error waiting for dedicated host (%s) to become ready: %s", d.Id(), err)
	}

	id := *dedicated.(datatypes.Virtual_DedicatedHost).Id
	d.SetId(fmt.Sprintf("%d", id))
	return resourceIBMComputeDedicatedHostRead(d, meta)
}

func resourceIBMComputeDedicatedHostRead(d *schema.ResourceData, meta interface{}) error {
	service := services.GetVirtualDedicatedHostService(meta.(ClientSession).SoftLayerSession())

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	result, err := service.Id(id).Mask(
		"name,cpuCount,datacenter,memoryCapacity,diskCapacity,backendRouter[hostname]",
	).GetObject()

	if err != nil {
		return fmt.Errorf("Error retrieving dedicated host: %s", err)
	}

	d.Set("hostname", result.Name)
	d.Set("datacenter", result.Datacenter.Name)
	d.Set("cpu_count", result.CpuCount)
	d.Set("disk_capacity", result.DiskCapacity)
	d.Set("memory_capacity", result.MemoryCapacity)
	d.Set("router_hostname", result.BackendRouter.Hostname)
	return nil
}

func resourceIBMComputeDedicatedHostUpdate(d *schema.ResourceData, meta interface{}) error {

	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetVirtualDedicatedHostService(sess.SetRetries(0))

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	result, err := service.Id(id).GetObject()
	if err != nil {
		return fmt.Errorf("Error retrieving dedicated host: %s", err)
	}

	if d.HasChange("hostname") {
		result.Name = sl.String(d.Get("hostname").(string))
		_, err = service.Id(id).EditObject(&result)
		if err != nil {
			return fmt.Errorf("Couldn't update dedicated host: %s", err)
		}

	}
	return resourceIBMComputeDedicatedHostRead(d, meta)
}

func resourceIBMComputeDedicatedHostDelete(d *schema.ResourceData, meta interface{}) error {
	service := services.GetVirtualDedicatedHostService(meta.(ClientSession).SoftLayerSession())

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	ok, err := service.Id(id).DeleteObject()

	if err != nil {
		return fmt.Errorf("Error deleting dedicated host: %s", err)
	}

	if !ok {
		return fmt.Errorf(
			"API reported it was unsuccessful in removing the dedicated host '%d'", id)
	}

	d.SetId("")
	return nil
}

func resourceIBMComputeDedicatedHostExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	service := services.GetVirtualDedicatedHostService(meta.(ClientSession).SoftLayerSession())
	dedicatedID, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	result, err := service.Id(dedicatedID).GetObject()
	if err != nil {
		if apiErr, ok := err.(sl.Error); ok {
			if apiErr.StatusCode == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error communicating with the API: %s", err)
	}

	return result.Id != nil && *result.Id == dedicatedID, nil
}

func findDedicatedHostByOrderID(d *datatypes.Hardware, r *schema.ResourceData, meta interface{}) (interface{}, error) {
	hostname := *d.Hostname

	log.Printf("Waiting for dedicated host (%s) to have to be provisioned", hostname)

	stateConf := &resource.StateChangeConf{
		Pending: []string{"retry", "pending"},
		Target:  []string{"provisioned"},
		Refresh: func() (interface{}, string, error) {
			service := services.GetAccountService(meta.(ClientSession).SoftLayerSession())
			dedicatedHosts, err := service.Filter(
				filter.Build(
					filter.Path("dedicatedHosts.name").Eq(hostname),
				),
			).Mask("id,createDate").GetDedicatedHosts()
			if err != nil {
				return false, "retry", nil
			}

			if len(dedicatedHosts) == 0 || dedicatedHosts[0].CreateDate == nil {
				return datatypes.Hardware{}, "pending", nil
			}
			return dedicatedHosts[0], "provisioned", nil

		},
		Timeout:    time.Duration(r.Get("wait_time_minutes").(int)) * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 1 * time.Minute,
	}

	return stateConf.WaitForState()
}
