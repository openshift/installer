// Copyright IBM Corp. 2017, 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	rID          = "route_id"
	rDestination = "destination"
	rAction      = "action"
	rNextHop     = "next_hop"
	rName        = "name"
	rZone        = "zone"
)

func ResourceIBMISVPCRoutingTableRoute() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMISVPCRoutingTableRouteCreate,
		Read:     resourceIBMISVPCRoutingTableRouteRead,
		Update:   resourceIBMISVPCRoutingTableRouteUpdate,
		Delete:   resourceIBMISVPCRoutingTableRouteDelete,
		Exists:   resourceIBMISVPCRoutingTableRouteExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			rtID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The routing table identifier.",
			},
			rtVpcID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The VPC identifier.",
			},
			rDestination: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The destination of the route.",
			},
			rZone: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The zone to apply the route to. Traffic from subnets in this zone will be subject to this route.",
			},
			rNextHop: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "If action is deliver, the next hop that packets will be delivered to. For other action values, its address will be 0.0.0.0.",
			},
			rAction: {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "deliver",
				Description:  "The action to perform with a packet matching the route.",
				ValidateFunc: validate.InvokeValidator("ibm_is_vpc_routing_table_route", rAction),
			},
			"advertise": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates whether this route will be advertised to the ingress sources specified by the `advertise_routes_to` routing table property.",
			},
			rName: {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     false,
				Computed:     true,
				Description:  "The user-defined name for this route.",
				ValidateFunc: validate.InvokeValidator("ibm_is_vpc_routing_table_route", rName),
			},
			rID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The routing table route identifier.",
			},
			rtHref: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Routing table route Href",
			},
			rtCreateAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Routing table route Created At",
			},
			"creator": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "If present, the resource that created the route. Routes with this property present cannot bedirectly deleted. All routes with an `origin` of `learned` or `service` will have thisproperty set, and future `origin` values may also have this property set.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VPN gateway's CRN.",
						},
						"deleted": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VPN gateway's canonical URL.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this VPN gateway.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this VPN gateway.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			rtLifecycleState: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Routing table route Lifecycle State",
			},
			"priority": {
				Type:         schema.TypeInt,
				Computed:     true,
				Optional:     true,
				Description:  "The route's priority. Smaller values have higher priority.",
				ValidateFunc: validate.InvokeValidator("ibm_is_vpc_routing_table_route", "priority"),
			},
			rtOrigin: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The origin of this route.",
			},
		},
	}
}

func ResourceIBMISVPCRoutingTableRouteValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)
	actionAllowedValues := "delegate, delegate_vpc, deliver, drop"

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 rtName,
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   false,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 rAction,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   false,
			AllowedValues:              actionAllowedValues})

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "priority",
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			MinValue:                   "0",
			MaxValue:                   "4"})
	ibmVPCRoutingTableRouteValidator := validate.ResourceValidator{ResourceName: "ibm_is_vpc_routing_table_route", Schema: validateSchema}
	return &ibmVPCRoutingTableRouteValidator
}

func resourceIBMISVPCRoutingTableRouteCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	vpcID := d.Get(rtVpcID).(string)
	tableID := d.Get(rtID).(string)
	destination := d.Get(rDestination).(string)
	zone := d.Get(rZone).(string)
	z := &vpcv1.ZoneIdentityByName{
		Name: core.StringPtr(zone),
	}

	createVpcRoutingTableRouteOptions := sess.NewCreateVPCRoutingTableRouteOptions(vpcID, tableID, destination, z)
	createVpcRoutingTableRouteOptions.SetZone(z)
	createVpcRoutingTableRouteOptions.SetDestination(destination)

	if add, ok := d.GetOk(rNextHop); ok {
		item := add.(string)
		if net.ParseIP(item) == nil {
			nhConnectionID := &vpcv1.RouteNextHopPrototype{
				ID: core.StringPtr(item),
			}
			createVpcRoutingTableRouteOptions.SetNextHop(nhConnectionID)
		} else {
			nh := &vpcv1.RouteNextHopPrototype{
				Address: core.StringPtr(item),
			}
			createVpcRoutingTableRouteOptions.SetNextHop(nh)
		}
	}

	if action, ok := d.GetOk(rAction); ok {
		routeAction := action.(string)
		createVpcRoutingTableRouteOptions.SetAction(routeAction)
	}

	if advertiseVal, ok := d.GetOk("advertise"); ok {
		advertise := advertiseVal.(bool)
		createVpcRoutingTableRouteOptions.SetAdvertise(advertise)
	}

	if name, ok := d.GetOk(rName); ok {
		routeName := name.(string)
		createVpcRoutingTableRouteOptions.SetName(routeName)
	}

	// Using GetOkExists to detet 0 as the possible values.
	if priority, ok := d.GetOkExists("priority"); ok {
		routePriority := priority.(int)
		createVpcRoutingTableRouteOptions.SetPriority(int64(routePriority))
	}

	route, response, err := sess.CreateVPCRoutingTableRoute(createVpcRoutingTableRouteOptions)
	if err != nil {
		log.Printf("[DEBUG] Create VPC Routing table route err %s\n%s", err, response)
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", vpcID, tableID, *route.ID))
	d.Set(rID, *route.ID)
	return resourceIBMISVPCRoutingTableRouteRead(d, meta)
}

func resourceIBMISVPCRoutingTableRouteRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	idSet := strings.Split(d.Id(), "/")
	getVpcRoutingTableRouteOptions := sess.NewGetVPCRoutingTableRouteOptions(idSet[0], idSet[1], idSet[2])
	route, response, err := sess.GetVPCRoutingTableRoute(getVpcRoutingTableRouteOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error Getting VPC Routing table route: %s\n%s", err, response)
	}

	d.Set(rID, *route.ID)
	if route.Advertise != nil {
		d.Set("Advertise", route.Advertise)
	}
	d.Set(rName, *route.Name)
	d.Set(rDestination, *route.Destination)
	if route.NextHop != nil {
		nexthop := route.NextHop.(*vpcv1.RouteNextHop)
		if nexthop.Address != nil {
			d.Set(rNextHop, *nexthop.Address)
		}
		if nexthop.ID != nil {
			d.Set(rNextHop, *nexthop.ID)
		}
	}
	if err = d.Set("origin", route.Origin); err != nil {
		return fmt.Errorf("[ERROR] Error setting origin %s", err)
	}
	if route.Zone != nil {
		d.Set(rZone, *route.Zone.Name)
	}
	d.Set(rtHref, route.Href)
	d.Set(rtLifecycleState, route.LifecycleState)
	d.Set(rtCreateAt, route.CreatedAt.String())
	creator := []map[string]interface{}{}
	if route.Creator != nil {
		mm, err := dataSourceIBMIsRouteCreatorToMap(route.Creator)
		if err != nil {
			log.Printf("Error reading VPC Routing Table Routes' creator:%s", err)
			return err
		}
		creator = append(creator, mm)
	}
	d.Set("creator", creator)
	d.Set("priority", route.Priority)
	return nil
}

func resourceIBMISVPCRoutingTableRouteUpdate(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	idSet := strings.Split(d.Id(), "/")
	hasChange := false
	routePatch := make(map[string]interface{})
	updateVpcRoutingTableRouteOptions := sess.NewUpdateVPCRoutingTableRouteOptions(idSet[0], idSet[1], idSet[2], routePatch)

	// Construct an instance of the RoutePatch model
	routePatchModel := new(vpcv1.RoutePatch)
	if d.HasChange("advertise") {
		advertiseVal := d.Get("advertise").(bool)
		routePatchModel.Advertise = &advertiseVal
		hasChange = true

	}
	if d.HasChange(rName) {
		name := d.Get(rName).(string)
		routePatchModel.Name = &name
		hasChange = true
	}
	if d.HasChange("priority") {
		rp := d.Get("priority").(int)
		routePriority := int64(rp)
		routePatchModel.Priority = &routePriority
		hasChange = true
	}

	if d.HasChange(rNextHop) {
		if add, ok := d.GetOk(rNextHop); ok {
			item := add.(string)
			if net.ParseIP(item) == nil {
				routePatchModel.NextHop = &vpcv1.RouteNextHopPatch{
					ID: core.StringPtr(item),
				}
				hasChange = true
			} else {
				routePatchModel.NextHop = &vpcv1.RouteNextHopPatch{
					Address: core.StringPtr(item),
				}
				hasChange = true
			}
		}
	}
	if hasChange {
		routePatchModelAsPatch, patchErr := routePatchModel.AsPatch()
		if patchErr != nil {
			return fmt.Errorf("[ERROR] Error calling asPatch for VPC Routing Table Route Patch: %s", patchErr)
		}
		updateVpcRoutingTableRouteOptions.RoutePatch = routePatchModelAsPatch
		_, response, err := sess.UpdateVPCRoutingTableRoute(updateVpcRoutingTableRouteOptions)
		if err != nil {
			log.Printf("[DEBUG] Update VPC Routing table route err %s\n%s", err, response)
			return err
		}
	}
	return resourceIBMISVPCRoutingTableRouteRead(d, meta)
}

func resourceIBMISVPCRoutingTableRouteDelete(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	idSet := strings.Split(d.Id(), "/")
	deleteVpcRoutingTableRouteOptions := sess.NewDeleteVPCRoutingTableRouteOptions(idSet[0], idSet[1], idSet[2])
	response, err := sess.DeleteVPCRoutingTableRoute(deleteVpcRoutingTableRouteOptions)
	if err != nil && response.StatusCode != 404 {
		log.Printf("Error deleting VPC Routing table route : %s", response)
		return err
	}

	d.SetId("")
	return nil
}

func resourceIBMISVPCRoutingTableRouteExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess, err := vpcClient(meta)
	if err != nil {
		return false, err
	}

	idSet := strings.Split(d.Id(), "/")
	if len(idSet) != 3 {
		return false, fmt.Errorf("[ERROR] Incorrect ID %s: ID should be a combination of vpcID/routingTableID/routeID", d.Id())
	}
	getVpcRoutingTableRouteOptions := sess.NewGetVPCRoutingTableRouteOptions(idSet[0], idSet[1], idSet[2])
	_, response, err := sess.GetVPCRoutingTableRoute(getVpcRoutingTableRouteOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return false, nil
		}
		return false, fmt.Errorf("[ERROR] Error Getting VPC Routing table route : %s\n%s", err, response)
	}
	return true, nil
}
