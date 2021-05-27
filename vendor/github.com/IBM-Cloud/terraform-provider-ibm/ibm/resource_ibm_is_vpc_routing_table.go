// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	rtID                         = "routing_table"
	rtVpcID                      = "vpc"
	rtName                       = "name"
	rtRouteDirectLinkIngress     = "route_direct_link_ingress"
	rtRouteTransitGatewayIngress = "route_transit_gateway_ingress"
	rtRouteVPCZoneIngress        = "route_vpc_zone_ingress"
	rtCreateAt                   = "created_at"
	rtHref                       = "href"
	rtIsDefault                  = "is_default"
	rtResourceType               = "resource_type"
	rtLifecycleState             = "lifecycle_state"
	rtSubnets                    = "subnets"
	rtDestination                = "destination"
	rtAction                     = "action"
	rtNextHop                    = "next_hop"
	rtZone                       = "zone"
	rtOrigin                     = "origin"
)

func resourceIBMISVPCRoutingTable() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMISVPCRoutingTableCreate,
		Read:     resourceIBMISVPCRoutingTableRead,
		Update:   resourceIBMISVPCRoutingTableUpdate,
		Delete:   resourceIBMISVPCRoutingTableDelete,
		Exists:   resourceIBMISVPCRoutingTableExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			rtVpcID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The VPC identifier.",
			},
			rtRouteDirectLinkIngress: {
				Type:        schema.TypeBool,
				ForceNew:    false,
				Default:     false,
				Optional:    true,
				Description: "If set to true, this routing table will be used to route traffic that originates from Direct Link to this VPC.",
			},
			rtRouteTransitGatewayIngress: {
				Type:        schema.TypeBool,
				ForceNew:    false,
				Default:     false,
				Optional:    true,
				Description: "If set to true, this routing table will be used to route traffic that originates from Transit Gateway to this VPC.",
			},
			rtRouteVPCZoneIngress: {
				Type:        schema.TypeBool,
				ForceNew:    false,
				Default:     false,
				Optional:    true,
				Description: "If set to true, this routing table will be used to route traffic that originates from subnets in other zones in this VPC.",
			},
			rtName: {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     false,
				Computed:     true,
				ValidateFunc: InvokeValidator("ibm_is_vpc_routing_table", rtName),
				Description:  "The user-defined name for this routing table.",
			},
			rtID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The routing table identifier.",
			},
			rtHref: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Routing table Href",
			},
			rtResourceType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Routing table Resource Type",
			},
			rtCreateAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Routing table Created At",
			},
			rtLifecycleState: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Routing table Lifecycle State",
			},
			rtIsDefault: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether this is the default routing table for this VPC",
			},
			rtSubnets: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						rtName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet name",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet ID",
						},
					},
				},
			},
		},
	}
}

func resourceIBMISVPCRoutingTableValidator() *ResourceValidator {

	validateSchema := make([]ValidateSchema, 2)
	actionAllowedValues := "delegate, delegate_vpc, deliver, drop"

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 rtName,
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Required:                   false,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 rtAction,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   false,
			AllowedValues:              actionAllowedValues})

	ibmISVPCRoutingTableValidator := ResourceValidator{ResourceName: "ibm_is_vpc_routing_table", Schema: validateSchema}
	return &ibmISVPCRoutingTableValidator
}

func resourceIBMISVPCRoutingTableCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	vpcID := d.Get(rtVpcID).(string)
	rtName := d.Get(rtName).(string)

	createVpcRoutingTableOptions := sess.NewCreateVPCRoutingTableOptions(vpcID)
	createVpcRoutingTableOptions.SetName(rtName)
	if _, ok := d.GetOk(rtRouteDirectLinkIngress); ok {
		routeDirectLinkIngress := d.Get(rtRouteDirectLinkIngress).(bool)
		createVpcRoutingTableOptions.RouteDirectLinkIngress = &routeDirectLinkIngress
	}
	if _, ok := d.GetOk(rtRouteTransitGatewayIngress); ok {
		routeTransitGatewayIngress := d.Get(rtRouteTransitGatewayIngress).(bool)
		createVpcRoutingTableOptions.RouteTransitGatewayIngress = &routeTransitGatewayIngress
	}
	if _, ok := d.GetOk(rtRouteVPCZoneIngress); ok {
		routeVPCZoneIngress := d.Get(rtRouteVPCZoneIngress).(bool)
		createVpcRoutingTableOptions.RouteVPCZoneIngress = &routeVPCZoneIngress
	}
	routeTable, response, err := sess.CreateVPCRoutingTable(createVpcRoutingTableOptions)
	if err != nil {
		log.Printf("[DEBUG] Create VPC Routing table err %s\n%s", err, response)
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s", vpcID, *routeTable.ID))

	return resourceIBMISVPCRoutingTableRead(d, meta)
}

func resourceIBMISVPCRoutingTableRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	idSet := strings.Split(d.Id(), "/")
	getVpcRoutingTableOptions := sess.NewGetVPCRoutingTableOptions(idSet[0], idSet[1])
	routeTable, response, err := sess.GetVPCRoutingTable(getVpcRoutingTableOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Getting VPC Routing table: %s\n%s", err, response)
	}

	d.Set(rtID, routeTable.ID)
	d.Set(rtName, routeTable.Name)

	d.Set(rtHref, routeTable.Href)
	d.Set(rtLifecycleState, routeTable.LifecycleState)
	d.Set(rtCreateAt, routeTable.CreatedAt.String())
	d.Set(rtResourceType, routeTable.ResourceType)
	d.Set(rtRouteDirectLinkIngress, routeTable.RouteDirectLinkIngress)
	d.Set(rtRouteTransitGatewayIngress, routeTable.RouteTransitGatewayIngress)
	d.Set(rtRouteVPCZoneIngress, routeTable.RouteVPCZoneIngress)
	d.Set(rtIsDefault, routeTable.IsDefault)

	subnets := make([]map[string]interface{}, 0)

	for _, s := range routeTable.Subnets {
		subnet := make(map[string]interface{})
		subnet[ID] = *s.ID
		subnet["name"] = *s.Name
		subnets = append(subnets, subnet)
	}

	d.Set(rtSubnets, subnets)

	return nil
}

func resourceIBMISVPCRoutingTableUpdate(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	idSet := strings.Split(d.Id(), "/")
	updateVpcRoutingTableOptions := new(vpcv1.UpdateVPCRoutingTableOptions)
	updateVpcRoutingTableOptions.VPCID = &idSet[0]
	updateVpcRoutingTableOptions.ID = &idSet[1]
	// Construct an instance of the RoutingTablePatch model
	routingTablePatchModel := new(vpcv1.RoutingTablePatch)

	if d.HasChange(rtName) {
		name := d.Get(rtName).(string)
		routingTablePatchModel.Name = core.StringPtr(name)
	}
	if d.HasChange(rtRouteDirectLinkIngress) {
		routeDirectLinkIngress := d.Get(rtRouteDirectLinkIngress).(bool)
		routingTablePatchModel.RouteDirectLinkIngress = core.BoolPtr(routeDirectLinkIngress)
	}
	if d.HasChange(rtRouteTransitGatewayIngress) {
		routeTransitGatewayIngress := d.Get(rtRouteTransitGatewayIngress).(bool)
		routingTablePatchModel.RouteTransitGatewayIngress = core.BoolPtr(routeTransitGatewayIngress)
	}
	if d.HasChange(rtRouteVPCZoneIngress) {
		routeVPCZoneIngress := d.Get(rtRouteVPCZoneIngress).(bool)
		routingTablePatchModel.RouteVPCZoneIngress = core.BoolPtr(routeVPCZoneIngress)
	}
	routingTablePatchModelAsPatch, asPatchErr := routingTablePatchModel.AsPatch()
	if asPatchErr != nil {
		return fmt.Errorf("Error calling asPatch for RoutingTablePatchModel: %s", asPatchErr)
	}
	updateVpcRoutingTableOptions.RoutingTablePatch = routingTablePatchModelAsPatch
	_, response, err := sess.UpdateVPCRoutingTable(updateVpcRoutingTableOptions)
	if err != nil {
		log.Printf("[DEBUG] Update VPC Routing table err %s\n%s", err, response)
		return err
	}
	return resourceIBMISVPCRoutingTableRead(d, meta)
}

func resourceIBMISVPCRoutingTableDelete(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	idSet := strings.Split(d.Id(), "/")

	deleteTableOptions := sess.NewDeleteVPCRoutingTableOptions(idSet[0], idSet[1])
	response, err := sess.DeleteVPCRoutingTable(deleteTableOptions)
	if err != nil && response.StatusCode != 404 {
		log.Printf("Error deleting VPC Routing table : %s", response)
		return err
	}

	d.SetId("")
	return nil
}

func resourceIBMISVPCRoutingTableExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess, err := vpcClient(meta)
	if err != nil {
		return false, err
	}

	idSet := strings.Split(d.Id(), "/")
	if len(idSet) < 2 {
		return false, fmt.Errorf("Incorrect ID %s: ID should be a Combination of vpcID/routingTableID", d.Id())
	}
	getVpcRoutingTableOptions := sess.NewGetVPCRoutingTableOptions(idSet[0], idSet[1])
	_, response, err := sess.GetVPCRoutingTable(getVpcRoutingTableOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return false, nil
		}
		return false, fmt.Errorf("Error Getting VPC Routing table : %s\n%s", err, response)
	}
	return true, nil
}
