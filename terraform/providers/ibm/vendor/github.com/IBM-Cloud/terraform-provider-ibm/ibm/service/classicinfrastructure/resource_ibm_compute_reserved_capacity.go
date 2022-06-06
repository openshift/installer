// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package classicinfrastructure

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/helpers/location"
	"github.com/softlayer/softlayer-go/helpers/product"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/sl"
)

func ResourceIBMComputeReservedCapacity() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMComputeReservedCapacityCreate,
		ReadContext:   resourceIBMComputeReservedCapacityRead,
		UpdateContext: resourceIBMComputeReservedCapacityUpdate,
		DeleteContext: resourceIBMComputeReservedCapacityDelete,
		Exists:        resourceIBMComputeReservedCapacityExists,
		Importer:      &schema.ResourceImporter{},
		CustomizeDiff: resourceReservedCapacityValidate,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"datacenter": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Dataceneter name",
			},

			"pod": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.TrimSpace(old) == strings.TrimSpace(new)
				},
				Description: "Pod name",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name",
			},

			"instances": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "no of the instances",
			},

			"flavor": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "flavor of the reserved capacity",
			},

			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "List of tags",
			},
			"force_create": {
				Type:             schema.TypeBool,
				Optional:         true,
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "Force the creation of reserved capacity with same name",
			},
		},
	}
}

func resourceIBMComputeReservedCapacityCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess := meta.(conns.ClientSession).SoftLayerSession()
	name := d.Get("name").(string)
	datacenter := d.Get("datacenter").(string)
	pod := d.Get("pod").(string)
	podName := datacenter + "." + pod
	PodService := services.GetNetworkPodService(sess)
	podMask := `backendRouterId,name`
	instances := d.Get("instances").(int)
	flavor := d.Get("flavor").(string)

	// 1.Getting the router ID
	routerids, err := PodService.Filter(filter.Path("datacenterName").Eq(datacenter).Build()).Mask(podMask).GetAllObjects()
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Encountered problem trying to get the router ID: %s", err))
	}
	var routerid int
	for _, iterate := range routerids {
		if *iterate.Name == podName {
			routerid = *iterate.BackendRouterId
		}
	}

	pkg, err := product.GetPackageByType(sess, "RESERVED_CAPACITY")
	if err != nil {
		return diag.FromErr(err)
	}

	// 2. Get all prices for the package
	productItems, err := product.GetPackageProducts(sess, *pkg.Id, productItemMaskWithPriceLocationGroupID)
	if err != nil {
		return diag.FromErr(err)
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

	// Lookup the data center ID
	dc, err := location.GetDatacenterByName(sess, datacenter)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] No data centers matching %s could be found", datacenter))
	}

	productOrderContainer := datatypes.Container_Product_Order_Virtual_ReservedCapacity{
		Container_Product_Order: datatypes.Container_Product_Order{
			PackageId:   pkg.Id,
			Location:    sl.String(strconv.Itoa(*dc.Id)),
			Prices:      priceItems,
			ComplexType: sl.String("SoftLayer_Container_Product_Order_Virtual_ReservedCapacity"),
			Quantity:    &instances,
		},
		Name:            sl.String(name),
		BackendRouterId: &routerid,
	}
	log.Println("[INFO] Creating reserved capacity")

	//verify order
	_, err = services.GetProductOrderService(sess.SetRetries(0)).
		VerifyOrder(&productOrderContainer)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] during creation of reserved capacity: %s", err))
	}
	//place order
	_, err = services.GetProductOrderService(sess.SetRetries(0)).
		PlaceOrder(&productOrderContainer, sl.Bool(false))
	if err != nil {
		return diag.FromErr(fmt.Errorf("[Error] during creation of reserved capacity: %s", err))
	}

	// wait for machine availability
	reservedCapacity, err := findReservedCapacityByOrderID(name, d, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf(
			"[Error] waiting for reserved capacity (%s) to become ready: %s", d.Id(), err))
	}

	id := *reservedCapacity.(datatypes.Virtual_ReservedCapacityGroup).Id
	d.SetId(fmt.Sprintf("%d", id))
	return resourceIBMComputeReservedCapacityRead(context, d, meta)
}
func findReservedCapacityByOrderID(name string, r *schema.ResourceData, meta interface{}) (interface{}, error) {

	log.Printf("Waiting for reserved capacity  (%s) to have to be provisioned", name)

	stateConf := &resource.StateChangeConf{
		Pending: []string{"retry", "pending"},
		Target:  []string{"provisioned"},
		Refresh: func() (interface{}, string, error) {
			service := services.GetAccountService(meta.(conns.ClientSession).SoftLayerSession())
			reservedCapacitys, err := service.Filter(
				filter.Build(
					filter.Path("reservedCapacityGroups.name").Eq(name),
				),
			).Mask("id,createDate").GetReservedCapacityGroups()
			if err != nil {
				return false, "retry", nil
			}

			if len(reservedCapacitys) == 0 || reservedCapacitys[0].CreateDate == nil {
				return datatypes.Virtual_ReservedCapacityGroup{}, "pending", nil
			}
			return reservedCapacitys[0], "provisioned", nil

		},
		Timeout:    r.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 1 * time.Minute,
	}

	return stateConf.WaitForState()
}

func resourceIBMComputeReservedCapacityRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess := meta.(conns.ClientSession).SoftLayerSession()
	service := services.GetVirtualReservedCapacityGroupService(sess)

	rgrpID, _ := strconv.Atoi(d.Id())

	rgrp, err := service.Id(rgrpID).Mask("id,name,instancesCount,backendRouter[hostname,datacenter[name]],instances[billingItem[item[keyName]]]").GetObject()
	if err != nil {
		if err, ok := err.(sl.Error); ok {
			if err.StatusCode == http.StatusNotFound {
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(fmt.Errorf("[Error] retrieving reserved capacity: %s", err))
	}

	d.Set("name", rgrp.Name)
	d.Set("datacenter", rgrp.BackendRouter.Datacenter.Name)
	pod := strings.SplitAfter(*rgrp.BackendRouter.Hostname, ".")[0]
	r, _ := regexp.Compile("[0-9]{2}")
	pod = "pod" + r.FindString(pod)
	d.Set("pod", pod)
	d.Set("instances", rgrp.InstancesCount)

	keyName, ok := sl.GrabOk(rgrp, "Instances.0.BillingItem.Item.KeyName")
	if ok {
		d.Set("flavor", keyName)
	}

	return nil
}

func resourceIBMComputeReservedCapacityUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess := meta.(conns.ClientSession).SoftLayerSession()
	service := services.GetVirtualReservedCapacityGroupService(sess)

	rgrpID, _ := strconv.Atoi(d.Id())

	opts := datatypes.Virtual_ReservedCapacityGroup{}

	if d.HasChange("name") {
		opts.Name = sl.String(d.Get("name").(string))
		_, err := service.Id(rgrpID).EditObject(&opts)

		if err != nil {
			return diag.FromErr(fmt.Errorf("[Error] editing reserved capacity: %s", err))
		}
	}

	return nil
}

func resourceIBMComputeReservedCapacityExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	sess := meta.(conns.ClientSession).SoftLayerSession()
	service := services.GetVirtualReservedCapacityGroupService(sess)

	rgrpID, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, fmt.Errorf("[Error] Not a valid ID, must be an integer: %s", err)
	}

	result, err := service.Id(rgrpID).GetObject()
	if err != nil {
		if apiErr, ok := err.(sl.Error); ok {
			if apiErr.StatusCode == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("[Error] communicating with the API: %s", err)
	}
	return result.Id != nil && *result.Id == rgrpID, nil
}

func resourceIBMComputeReservedCapacityDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[WARN]: `terraform destroy` does not remove the reserved capacity but only clears the state file. We cannot cancel reserved capacity")
	d.SetId("")
	return nil
}

func resourceReservedCapacityValidate(_ context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	forceCreate := diff.Get("force_create").(bool)
	if diff.Id() == "" && !forceCreate {
		name := diff.Get("name").(string)
		service := services.GetAccountService(meta.(conns.ClientSession).SoftLayerSession())
		reservedCapacities, _ := service.Filter(
			filter.Build(
				filter.Path("reservedCapacityGroups.name").Eq(name),
			),
		).Mask("id,createDate").GetReservedCapacityGroups()
		if len(reservedCapacities) > 0 {
			return fmt.Errorf("reserved capacity exists with same name [%s] if you still want to provision with same name set force_create argument to true", name)
		}
	}

	return nil

}
