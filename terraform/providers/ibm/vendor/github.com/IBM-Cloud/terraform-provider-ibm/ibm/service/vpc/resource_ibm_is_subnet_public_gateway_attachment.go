// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isPublicGatewayID                  = "public_gateway"
	IsPublicGatewayResourceType        = "resource_type"
	IsPublicGatewayAttachmentAvailable = "available"
	IsPublicGatewayAttachmentDeleting  = "deleting"
	IsPublicGatewayAttachmentFailed    = "failed"
	IsPublicGatewayAttachmentPending   = "pending"
)

func ResourceIBMISSubnetPublicGatewayAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMISSubnetPublicGatewayAttachmentCreate,
		ReadContext:   resourceIBMISSubnetPublicGatewayAttachmentRead,
		UpdateContext: resourceIBMISSubnetPublicGatewayAttachmentUpdate,
		DeleteContext: resourceIBMISSubnetPublicGatewayAttachmentDelete,
		Importer:      &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			isSubnetID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The subnet identifier",
			},

			isPublicGatewayID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of public gateway",
			},

			isPublicGatewayName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the Public gateway instance",
			},

			isPublicGatewayFloatingIP: {
				Type:     schema.TypeMap,
				Computed: true,
			},

			isPublicGatewayStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public gateway instance status",
			},

			isPublicGatewayResourceGroup: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public gateway resource group info",
			},

			isPublicGatewayVPC: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public gateway VPC info",
			},

			isPublicGatewayZone: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public gateway zone info",
			},

			IsPublicGatewayResourceType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the resource",
			},

			isPublicGatewayCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},

			flex.ResourceGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group name in which resource is provisioned",
			},
		},
	}
}

func resourceIBMISSubnetPublicGatewayAttachmentCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	subnet := d.Get(isSubnetID).(string)
	publicGateway := d.Get(isPublicGatewayID).(string)

	publicGatewayIdentity := &vpcv1.PublicGatewayIdentity{
		ID: &publicGateway,
	}

	// Construct an instance of the SetSubnetPublicGatewayOptions
	setSubnetPublicGatewayOptions := &vpcv1.SetSubnetPublicGatewayOptions{
		ID:                    &subnet,
		PublicGatewayIdentity: publicGatewayIdentity,
	}

	pg, response, err := sess.SetSubnetPublicGatewayWithContext(context, setSubnetPublicGatewayOptions)

	if err != nil {
		log.Printf("[DEBUG] Error while attaching public gateway(%s) to subnet(%s) %s\n%s", publicGateway, subnet, err, response)
		return diag.FromErr(fmt.Errorf("[ERROR] Error while attaching public gateway(%s) to subnet(%s) %s\n%s", publicGateway, subnet, err, response))
	}
	d.SetId(subnet)
	_, err = isWaitForSubnetPublicGatewayAvailable(sess, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Public Gateway : %s", *pg.ID)
	log.Printf("[INFO] Subnet ID : %s", subnet)

	return resourceIBMISSubnetPublicGatewayAttachmentRead(context, d, meta)
}

func resourceIBMISSubnetPublicGatewayAttachmentRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := d.Id()
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	getSubnetPublicGatewayOptionsModel := &vpcv1.GetSubnetPublicGatewayOptions{
		ID: &id,
	}
	pg, response, err := sess.GetSubnetPublicGatewayWithContext(context, getSubnetPublicGatewayOptionsModel)

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("[ERROR] Error getting subnet's (%s) attached public gateway: %s\n%s", id, err, response))
	}
	d.Set(isPublicGatewayName, pg.Name)
	d.Set(isSubnetID, id)
	d.Set(isPublicGatewayID, pg.ID)
	if pg.FloatingIP != nil {
		floatIP := map[string]interface{}{
			"id":                             *pg.FloatingIP.ID,
			isPublicGatewayFloatingIPAddress: *pg.FloatingIP.Address,
		}
		d.Set(isPublicGatewayFloatingIP, floatIP)
	}
	d.Set(isPublicGatewayStatus, pg.Status)
	if pg.ResourceGroup != nil {
		d.Set(isPublicGatewayResourceGroup, *pg.ResourceGroup.ID)
		d.Set(flex.ResourceGroupName, *pg.ResourceGroup.Name)
	}
	d.Set(isPublicGatewayVPC, *pg.VPC.ID)
	d.Set(isPublicGatewayZone, *pg.Zone.Name)
	d.Set(IsPublicGatewayResourceType, pg.ResourceType)
	d.Set(isPublicGatewayCRN, pg.CRN)

	return nil
}

func resourceIBMISSubnetPublicGatewayAttachmentUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	if d.HasChange(isPublicGatewayID) {
		subnet := d.Get(isSubnetID).(string)
		publicGateway := d.Get(isPublicGatewayID).(string)

		publicGatewayIdentity := &vpcv1.PublicGatewayIdentity{
			ID: &publicGateway,
		}

		// Construct an instance of the SetSubnetPublicGatewayOptions
		setSubnetPublicGatewayOptions := &vpcv1.SetSubnetPublicGatewayOptions{
			ID:                    &subnet,
			PublicGatewayIdentity: publicGatewayIdentity,
		}

		pg, response, err := sess.SetSubnetPublicGatewayWithContext(context, setSubnetPublicGatewayOptions)

		if err != nil || pg == nil {
			log.Printf("[DEBUG] Error while attaching public gateway(%s) to subnet(%s) %s\n%s", publicGateway, subnet, err, response)
			return diag.FromErr(fmt.Errorf("[ERROR] Error while attaching public gateway(%s) to subnet(%s) %s\n%s", publicGateway, subnet, err, response))
		}
		log.Printf("[INFO] Updated subnet %s with public gateway(%s)", subnet, publicGateway)

		d.SetId(subnet)
		return resourceIBMISSubnetPublicGatewayAttachmentRead(context, d, meta)
	}

	return resourceIBMISSubnetPublicGatewayAttachmentRead(context, d, meta)
}

func resourceIBMISSubnetPublicGatewayAttachmentDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := d.Id()
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	// Get subnet details
	getSubnetOptions := &vpcv1.GetSubnetOptions{
		ID: &id,
	}
	_, response, err := sess.GetSubnetWithContext(context, getSubnetOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("[ERROR] Error Getting Subnet (%s): %s\n%s", id, err, response))
	}

	// Construct an instance of the UnsetSubnetPublicGatewayOptions model
	unsetSubnetPublicGatewayOptions := &vpcv1.UnsetSubnetPublicGatewayOptions{
		ID: &id,
	}
	res, err := sess.UnsetSubnetPublicGatewayWithContext(context, unsetSubnetPublicGatewayOptions)

	if err != nil {
		log.Printf("[DEBUG] Error while detaching public gateway to subnet %s\n%s", err, res)
		return diag.FromErr(fmt.Errorf("[ERROR] Error while detaching public gateway to subnet %s\n%s", err, res))
	}
	_, err = isWaitForSubnetPublicGatewayDelete(sess, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func isWaitForSubnetPublicGatewayAvailable(subnetC *vpcv1.VpcV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for subnet (%s) public gateway attachment to be available.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{IsPublicGatewayAttachmentPending, IsPublicGatewayAttachmentDeleting},
		Target:     []string{IsPublicGatewayAttachmentAvailable, IsPublicGatewayAttachmentFailed, ""},
		Refresh:    isSubnetPublicGatewayRefreshFunc(subnetC, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isSubnetPublicGatewayRefreshFunc(subnetC *vpcv1.VpcV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getSubnetPublicGatewayOptionsModel := &vpcv1.GetSubnetPublicGatewayOptions{
			ID: &id,
		}
		pg, response, err := subnetC.GetSubnetPublicGateway(getSubnetPublicGatewayOptionsModel)

		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return pg, "", fmt.Errorf("[ERROR] Error getting subnet's (%s) attached public gateway: %s\n%s", id, err, response)
			}
			return pg, "", fmt.Errorf("[ERROR] Error getting subnet's (%s) attached public gateway: %s\n%s", id, err, response)
		}

		if *pg.Status == "failed" {
			return pg, IsPublicGatewayAttachmentFailed, fmt.Errorf("[ERROR] Error subnet (%s) public gateway attachment failed: %s\n%s", id, err, response)
		}

		return pg, *pg.Status, nil
	}
}

func isWaitForSubnetPublicGatewayDelete(subnetC *vpcv1.VpcV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for subnet (%s) public gateway attachment to be detached.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{IsPublicGatewayAttachmentPending, IsPublicGatewayAttachmentDeleting},
		Target:     []string{IsPublicGatewayAttachmentAvailable, IsPublicGatewayAttachmentFailed, ""},
		Refresh:    isSubnetPublicGatewayDeleteRefreshFunc(subnetC, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isSubnetPublicGatewayDeleteRefreshFunc(subnetC *vpcv1.VpcV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getSubnetPublicGatewayOptionsModel := &vpcv1.GetSubnetPublicGatewayOptions{
			ID: &id,
		}
		pg, response, err := subnetC.GetSubnetPublicGateway(getSubnetPublicGatewayOptionsModel)

		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return pg, "", nil
			}
			return pg, "", fmt.Errorf("[ERROR] Error getting subnet's (%s) attached public gateway: %s\n%s", id, err, response)
		}

		if *pg.Status == "failed" {
			return pg, IsPublicGatewayAttachmentFailed, fmt.Errorf("[ERROR] Error subnet (%s) public gateway attachment failed: %s\n%s", id, err, response)
		}

		return pg, *pg.Status, nil
	}
}
