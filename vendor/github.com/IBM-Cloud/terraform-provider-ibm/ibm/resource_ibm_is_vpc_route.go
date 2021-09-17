// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	isVPCRouteName            = "name"
	isVPCRouteState           = "status"
	isVPCRouteNextHop         = "next_hop"
	isVPCRouteDestinationCIDR = "destination"
	isVPCRouteLocation        = "zone"
	isVPCRouteVPCID           = "vpc"

	isRouteStatusPending  = "pending"
	isRouteStatusUpdating = "updating"
	isRouteStatusStable   = "stable"
	isRouteStatusFailed   = "failed"

	isRouteStatusDeleting = "deleting"
	isRouteStatusDeleted  = "deleted"
)

func resourceIBMISVpcRoute() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMISVpcRouteCreate,
		Read:     resourceIBMISVpcRouteRead,
		Update:   resourceIBMISVpcRouteUpdate,
		Delete:   resourceIBMISVpcRouteDelete,
		Exists:   resourceIBMISVpcRouteExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			isVPCRouteName: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: InvokeValidator("ibm_is_route", isVPCRouteName),
				Description:  "VPC route name",
			},
			isVPCRouteLocation: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPC route location",
			},

			isVPCRouteDestinationCIDR: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: InvokeValidator("ibm_is_route", isVPCRouteDestinationCIDR),
				Description:  "VPC route destination CIDR value",
			},

			isVPCRouteState: {
				Type:     schema.TypeString,
				Computed: true,
			},

			isVPCRouteVPCID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPC instance ID",
			},

			isVPCRouteNextHop: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPC route next hop value",
			},

			RelatedCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the VPC resource",
			},
		},
	}
}

func resourceIBMISRouteValidator() *ResourceValidator {

	validateSchema := make([]ValidateSchema, 1)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isVPCRouteName,
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Required:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isVPCRouteDestinationCIDR,
			ValidateFunctionIdentifier: ValidateCIDRAddress,
			Type:                       TypeString,
			ForceNew:                   true,
			Required:                   true})

	ibmISRouteResourceValidator := ResourceValidator{ResourceName: "ibm_is_route", Schema: validateSchema}
	return &ibmISRouteResourceValidator
}

func resourceIBMISVpcRouteCreate(d *schema.ResourceData, meta interface{}) error {

	routeName := d.Get(isVPCRouteName).(string)
	zoneName := d.Get(isVPCRouteLocation).(string)
	cidr := d.Get(isVPCRouteDestinationCIDR).(string)
	vpcID := d.Get(isVPCRouteVPCID).(string)
	nextHop := d.Get(isVPCRouteNextHop).(string)
	err := vpcRouteCreate(d, meta, routeName, zoneName, cidr, vpcID, nextHop)
	if err != nil {
		return err
	}
	return resourceIBMISVpcRouteRead(d, meta)
}

func vpcRouteCreate(d *schema.ResourceData, meta interface{}, routeName, zoneName, cidr, vpcID, nextHop string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	createRouteOptions := &vpcv1.CreateVPCRouteOptions{
		VPCID:       &vpcID,
		Destination: &cidr,
		Name:        &routeName,
		NextHop: &vpcv1.RouteNextHopPrototype{
			Address: &nextHop,
		},
		Zone: &vpcv1.ZoneIdentity{
			Name: &zoneName,
		},
	}
	route, response, err := sess.CreateVPCRoute(createRouteOptions)
	if err != nil {
		return fmt.Errorf("Error while creating VPC Route err %s\n%s", err, response)
	}
	routeID := *route.ID

	d.SetId(fmt.Sprintf("%s/%s", vpcID, routeID))

	_, err = isWaitForRouteStable(sess, d, vpcID, routeID)
	if err != nil {
		return err
	}
	return nil
}

func isWaitForRouteStable(sess *vpcv1.VpcV1, d *schema.ResourceData, vpcID, routeID string) (interface{}, error) {

	stateConf := &resource.StateChangeConf{
		Pending: []string{isRouteStatusPending, isRouteStatusUpdating},
		Target:  []string{isRouteStatusStable, isRouteStatusFailed},
		Refresh: func() (interface{}, string, error) {
			getVpcRouteOptions := &vpcv1.GetVPCRouteOptions{
				VPCID: &vpcID,
				ID:    &routeID,
			}
			route, response, err := sess.GetVPCRoute(getVpcRouteOptions)
			if err != nil {
				return route, "", fmt.Errorf("Error Getting VPC Route: %s\n%s", err, response)
			}

			if *route.LifecycleState == "stable" || *route.LifecycleState == "failed" {
				return route, *route.LifecycleState, nil
			}
			return route, *route.LifecycleState, nil
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func resourceIBMISVpcRouteRead(d *schema.ResourceData, meta interface{}) error {

	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}

	vpcID := parts[0]
	routeID := parts[1]
	error := vpcRouteGet(d, meta, vpcID, routeID)
	if error != nil {
		return error
	}
	return nil
}

func vpcRouteGet(d *schema.ResourceData, meta interface{}, vpcID, routeID string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	getVpcRouteOptions := &vpcv1.GetVPCRouteOptions{
		VPCID: &vpcID,
		ID:    &routeID,
	}
	route, response, err := sess.GetVPCRoute(getVpcRouteOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Getting VPC Route (%s): %s\n%s", routeID, err, response)
	}
	d.Set(isVPCRouteVPCID, vpcID)
	d.Set(isVPCRouteName, route.Name)
	if route.Zone != nil {
		d.Set(isVPCRouteLocation, *route.Zone.Name)
	}
	d.Set(isVPCRouteDestinationCIDR, *route.Destination)
	nexthop := route.NextHop.(*vpcv1.RouteNextHop)
	d.Set(isVPCRouteNextHop, *nexthop.Address)
	d.Set(isVPCRouteState, *route.LifecycleState)
	getVPCOptions := &vpcv1.GetVPCOptions{
		ID: &vpcID,
	}
	vpc, response, err := sess.GetVPC(getVPCOptions)
	if err != nil {
		return fmt.Errorf("Error Getting VPC : %s\n%s", err, response)
	}
	d.Set(RelatedCRN, *vpc.CRN)

	return nil
}

func resourceIBMISVpcRouteUpdate(d *schema.ResourceData, meta interface{}) error {

	name := ""
	hasChanged := false

	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}

	vpcID := parts[0]
	routeID := parts[1]
	if d.HasChange(isVPCRouteName) {
		name = d.Get(isVPCRouteName).(string)
		hasChanged = true
	}

	error := vpcRouteUpdate(d, meta, vpcID, routeID, name, hasChanged)
	if error != nil {
		return error
	}

	return resourceIBMISVpcRouteRead(d, meta)
}

func vpcRouteUpdate(d *schema.ResourceData, meta interface{}, vpcID, routeID, name string, hasChanged bool) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	if hasChanged {
		updateVpcRouteOptions := &vpcv1.UpdateVPCRouteOptions{
			VPCID: &vpcID,
			ID:    &routeID,
		}
		routePatchModel := &vpcv1.RoutePatch{
			Name: &name,
		}
		routePatch, err := routePatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("Error calling asPatch for RoutePatch: %s", err)
		}
		updateVpcRouteOptions.RoutePatch = routePatch
		_, response, err := sess.UpdateVPCRoute(updateVpcRouteOptions)
		if err != nil {
			return fmt.Errorf("Error Updating VPC Route: %s\n%s", err, response)
		}
	}
	return nil
}

func resourceIBMISVpcRouteDelete(d *schema.ResourceData, meta interface{}) error {
	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	vpcID := parts[0]
	routeID := parts[1]
	error := vpcRouteDelete(d, meta, vpcID, routeID)
	if error != nil {
		return error
	}

	d.SetId("")
	return nil
}

func vpcRouteDelete(d *schema.ResourceData, meta interface{}, vpcID, routeID string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	getVpcRouteOptions := &vpcv1.GetVPCRouteOptions{
		VPCID: &vpcID,
		ID:    &routeID,
	}
	_, response, err := sess.GetVPCRoute(getVpcRouteOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil
		}
		return fmt.Errorf("Error Getting VPC Route (%s): %s\n%s", routeID, err, response)
	}
	deleteRouteOptions := &vpcv1.DeleteVPCRouteOptions{
		VPCID: &vpcID,
		ID:    &routeID,
	}
	response, err = sess.DeleteVPCRoute(deleteRouteOptions)
	if err != nil {
		return fmt.Errorf("Error Deleting VPC Route: %s\n%s", err, response)
	}
	_, err = isWaitForVPCRouteDeleted(sess, vpcID, routeID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func isWaitForVPCRouteDeleted(sess *vpcv1.VpcV1, vpcID, routeID string, timeout time.Duration) (interface{}, error) {

	log.Printf("Waiting for VPC Route (%s) to be deleted.", routeID)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"retry", isRouteStatusDeleting},
		Target:  []string{isRouteStatusDeleted, isRouteStatusFailed},
		Refresh: func() (interface{}, string, error) {
			getVpcRouteOptions := &vpcv1.GetVPCRouteOptions{
				VPCID: &vpcID,
				ID:    &routeID,
			}
			route, response, err := sess.GetVPCRoute(getVpcRouteOptions)
			if err != nil {
				if response != nil && response.StatusCode == 404 {
					return route, isRouteStatusDeleted, nil
				}
				return route, isRouteStatusDeleting, fmt.Errorf("The VPC route %s failed to delete: %s\n%s", routeID, err, response)
			}
			return route, isRouteStatusDeleting, nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func resourceIBMISVpcRouteExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	parts, err := idParts(d.Id())
	if err != nil {
		return false, err
	}
	if len(parts) != 2 {
		return false, fmt.Errorf("Incorrect ID %s: ID should be a combination of vpcID/routeID", d.Id())
	}
	vpcID := parts[0]
	routeID := parts[1]
	exists, err := vpcRouteExists(d, meta, vpcID, routeID)
	return exists, err
}

func vpcRouteExists(d *schema.ResourceData, meta interface{}, vpcID, routeID string) (bool, error) {
	sess, err := vpcClient(meta)
	if err != nil {
		return false, err
	}
	getVpcRouteOptions := &vpcv1.GetVPCRouteOptions{
		VPCID: &vpcID,
		ID:    &routeID,
	}
	_, response, err := sess.GetVPCRoute(getVpcRouteOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("Error getting VPC Route: %s\n%s", err, response)
	}
	return true, nil
}
