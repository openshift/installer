// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isVPNGatewayName              = "name"
	isVPNGatewayResourceGroup     = "resource_group"
	isVPNGatewayMode              = "mode"
	isVPNGatewayCRN               = "crn"
	isVPNGatewayTags              = "tags"
	isVPNGatewaySubnet            = "subnet"
	isVPNGatewayStatus            = "status"
	isVPNGatewayDeleting          = "deleting"
	isVPNGatewayDeleted           = "done"
	isVPNGatewayProvisioning      = "provisioning"
	isVPNGatewayProvisioningDone  = "done"
	isVPNGatewayPublicIPAddress   = "public_ip_address"
	isVPNGatewayMembers           = "members"
	isVPNGatewayCreatedAt         = "created_at"
	isVPNGatewayPublicIPAddress2  = "public_ip_address2"
	isVPNGatewayPrivateIPAddress  = "private_ip_address"
	isVPNGatewayPrivateIPAddress2 = "private_ip_address2"
	isVPNGatewayAccessTags        = "access_tags"
	isVPNGatewayHealthState       = "health_state"
	isVPNGatewayHealthReasons     = "health_reasons"
	isVPNGatewayLifecycleState    = "lifecycle_state"
	isVPNGatewayLifecycleReasons  = "lifecycle_reasons"
)

func ResourceIBMISVPNGateway() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMISVPNGatewayCreate,
		Read:     resourceIBMISVPNGatewayRead,
		Update:   resourceIBMISVPNGatewayUpdate,
		Delete:   resourceIBMISVPNGatewayDelete,
		Exists:   resourceIBMISVPNGatewayExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			customdiff.Sequence(
				func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
					return flex.ResourceTagsCustomizeDiff(diff)
				},
			),
			customdiff.Sequence(
				func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
					return flex.ResourceValidateAccessTags(diff, v)
				}),
		),

		Schema: map[string]*schema.Schema{

			isVPNGatewayName: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validate.InvokeValidator("ibm_is_vpn_gateway", isVPNGatewayName),
				Description:  "VPN Gateway instance name",
			},

			isVPNGatewaySubnet: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPNGateway subnet info",
			},

			isVPNGatewayResourceGroup: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: "The resource group for this VPN gateway",
			},

			isVPNGatewayStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the VPN gateway",
			},

			isVPNGatewayHealthState: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The health of this resource.- `ok`: Healthy- `degraded`: Suffering from compromised performance, capacity, or connectivity- `faulted`: Completely unreachable, inoperative, or otherwise entirely incapacitated- `inapplicable`: The health state does not apply because of the current lifecycle state. A resource with a lifecycle state of `failed` or `deleting` will have a health state of `inapplicable`. A `pending` resource may also have this state.",
			},
			isVPNGatewayHealthReasons: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A snake case string succinctly identifying the reason for this health state.",
						},

						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An explanation of the reason for this health state.",
						},

						"more_info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Link to documentation about the reason for this health state.",
						},
					},
				},
			},

			isVPNGatewayPublicIPAddress: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The public IP address assigned to the VPN gateway member.",
			},

			isVPNGatewayPublicIPAddress2: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The second public IP address assigned to the VPN gateway member.",
			},

			isVPNGatewayPrivateIPAddress: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Private IP address assigned to the VPN gateway member.",
			},

			isVPNGatewayPrivateIPAddress2: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Second Private IP address assigned to the VPN gateway member.",
			},

			isVPNGatewayTags: {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_vpn_gateway", "tags")},
				Set:         flex.ResourceIBMVPCHash,
				Description: "VPN Gateway tags list",
			},

			isVPNGatewayAccessTags: {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_vpn_gateway", "accesstag")},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of access management tags",
			},

			flex.ResourceControllerURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the IBM Cloud dashboard that can be used to explore and view details about this instance",
			},

			flex.ResourceName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the resource",
			},

			flex.ResourceCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},

			isVPNGatewayCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},

			flex.ResourceStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the resource",
			},

			flex.ResourceGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group name in which resource is provisioned",
			},
			isVPNGatewayCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Created Time of the VPN Gateway",
			},
			isVPNGatewayLifecycleState: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the VPN route.",
			},
			isVPNGatewayLifecycleReasons: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The reasons for the current lifecycle_state (if any).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A snake case string succinctly identifying the reason for this lifecycle state.",
						},

						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An explanation of the reason for this lifecycle state.",
						},

						"more_info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Link to documentation about the reason for this lifecycle state.",
						},
					},
				},
			},
			isVPNGatewayMode: {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "route",
				ValidateFunc: validate.InvokeValidator("ibm_is_vpn_gateway", isVPNGatewayMode),
				Description:  "mode in VPN gateway(route/policy)",
			},

			isVPNGatewayMembers: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of VPN gateway members",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The public IP address assigned to the VPN gateway member",
						},

						"private_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The private IP address assigned to the VPN gateway member",
						},

						"role": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The high availability role assigned to the VPN gateway member",
						},

						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the VPN gateway member",
						},
					},
				},
			},
			"vpc": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "VPC for the VPN Gateway",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this VPC.",
						},
						"deleted": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this VPC.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this VPC.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique user-defined name for this VPC.",
						},
					},
				},
			},
		},
	}
}

func ResourceIBMISVPNGatewayValidator() *validate.ResourceValidator {

	modeCheckTypes := "route,policy"
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isVPNGatewayName,
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isVPNGatewayMode,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   false,
			AllowedValues:              modeCheckTypes})

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "tags",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9:_ .-]+$`,
			MinValueLength:             1,
			MaxValueLength:             128})

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "accesstag",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([A-Za-z0-9_.-]|[A-Za-z0-9_.-][A-Za-z0-9_ .-]*[A-Za-z0-9_.-]):([A-Za-z0-9_.-]|[A-Za-z0-9_.-][A-Za-z0-9_ .-]*[A-Za-z0-9_.-])$`,
			MinValueLength:             1,
			MaxValueLength:             128})

	ibmISVPNGatewayResourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_vpn_gateway", Schema: validateSchema}
	return &ibmISVPNGatewayResourceValidator
}

func resourceIBMISVPNGatewayCreate(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] VPNGateway create")
	name := d.Get(isVPNGatewayName).(string)
	subnetID := d.Get(isVPNGatewaySubnet).(string)
	mode := d.Get(isVPNGatewayMode).(string)

	err := vpngwCreate(d, meta, name, subnetID, mode)
	if err != nil {
		return err
	}
	return resourceIBMISVPNGatewayRead(d, meta)
}

func vpngwCreate(d *schema.ResourceData, meta interface{}, name, subnetID, mode string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	vpnGatewayPrototype := &vpcv1.VPNGatewayPrototype{
		Subnet: &vpcv1.SubnetIdentity{
			ID: &subnetID,
		},
		Name: &name,
		Mode: &mode,
	}
	options := &vpcv1.CreateVPNGatewayOptions{
		VPNGatewayPrototype: vpnGatewayPrototype,
	}

	if rgrp, ok := d.GetOk(isVPNGatewayResourceGroup); ok {
		rg := rgrp.(string)
		vpnGatewayPrototype.ResourceGroup = &vpcv1.ResourceGroupIdentity{
			ID: &rg,
		}
	}

	vpnGatewayIntf, response, err := sess.CreateVPNGateway(options)
	if err != nil {
		return fmt.Errorf("[DEBUG] Create vpc VPN Gateway %s\n%s", err, response)
	}
	vpnGateway := vpnGatewayIntf.(*vpcv1.VPNGateway)

	d.SetId(*vpnGateway.ID)
	log.Printf("[INFO] VPNGateway : %s", *vpnGateway.ID)

	_, err = isWaitForVpnGatewayAvailable(sess, *vpnGateway.ID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk(isVPNGatewayTags); ok || v != "" {
		oldList, newList := d.GetChange(isVPNGatewayTags)
		err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *vpnGateway.CRN, "", isUserTagType)
		if err != nil {
			log.Printf(
				"Error on create of resource vpc VPN Gateway (%s) tags: %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk(isVPNGatewayAccessTags); ok {
		oldList, newList := d.GetChange(isVPNGatewayAccessTags)
		err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *vpnGateway.CRN, "", isAccessTagType)
		if err != nil {
			log.Printf(
				"Error on create of resource VPN Gateway (%s) access tags: %s", d.Id(), err)
		}
	}

	return nil
}

func isWaitForVpnGatewayAvailable(vpnGateway *vpcv1.VpcV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for vpn gateway (%s) to be available.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isVPNGatewayProvisioning},
		Target:     []string{isVPNGatewayProvisioningDone, ""},
		Refresh:    isVpnGatewayRefreshFunc(vpnGateway, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isVpnGatewayRefreshFunc(vpnGateway *vpcv1.VpcV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getVpnGatewayOptions := &vpcv1.GetVPNGatewayOptions{
			ID: &id,
		}
		vpnGatewayIntf, response, err := vpnGateway.GetVPNGateway(getVpnGatewayOptions)
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] Error Getting Vpn Gateway: %s\n%s", err, response)
		}
		vpnGateway := vpnGatewayIntf.(*vpcv1.VPNGateway)

		if *vpnGateway.LifecycleState == "stable" || *vpnGateway.LifecycleState == "failed" || *vpnGateway.LifecycleState == "suspended" {
			return vpnGateway, isVPNGatewayProvisioningDone, nil
		}

		return vpnGateway, isVPNGatewayProvisioning, nil
	}
}

func resourceIBMISVPNGatewayRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()

	err := vpngwGet(d, meta, id)
	if err != nil {
		return err
	}
	return nil
}

func vpngwGet(d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	getVpnGatewayOptions := &vpcv1.GetVPNGatewayOptions{
		ID: &id,
	}
	vpnGatewayIntf, response, err := sess.GetVPNGateway(getVpnGatewayOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error Getting Vpn Gateway (%s): %s\n%s", id, err, response)
	}
	vpnGateway := vpnGatewayIntf.(*vpcv1.VPNGateway)

	d.Set(isVPNGatewayName, *vpnGateway.Name)
	d.Set(isVPNGatewaySubnet, *vpnGateway.Subnet.ID)
	if err = d.Set(isVPNGatewayHealthState, vpnGateway.HealthState); err != nil {
		return fmt.Errorf("[ERROR] Error setting health_state: %s", err)
	}
	if err := d.Set(isVPNGatewayHealthReasons, resourceVPNGatewayRouteFlattenHealthReasons(vpnGateway.HealthReasons)); err != nil {
		return fmt.Errorf("[ERROR] Error setting health_reasons: %s", err)
	}
	if err = d.Set(isVPNGatewayLifecycleState, vpnGateway.LifecycleState); err != nil {
		return fmt.Errorf("[ERROR] Error setting lifecycle_state: %s", err)
	}
	if err := d.Set(isVPNGatewayLifecycleReasons, resourceVPNGatewayFlattenLifecycleReasons(vpnGateway.LifecycleReasons)); err != nil {
		return fmt.Errorf("[ERROR] Error setting lifecycle_reasons: %s", err)
	}
	members := []vpcv1.VPNGatewayMember{}
	for _, member := range vpnGateway.Members {
		members = append(members, member)
	}
	if len(members) > 0 {
		d.Set(isVPNGatewayPublicIPAddress, *members[0].PublicIP.Address)
		if members[0].PrivateIP != nil && members[0].PrivateIP.Address != nil {
			d.Set(isVPNGatewayPrivateIPAddress, *members[0].PrivateIP.Address)
		}
	}
	if len(members) > 1 {
		d.Set(isVPNGatewayPublicIPAddress2, *members[1].PublicIP.Address)
		if members[1].PrivateIP != nil && members[1].PrivateIP.Address != nil {
			d.Set(isVPNGatewayPrivateIPAddress2, *members[1].PrivateIP.Address)
		}

	}
	tags, err := flex.GetGlobalTagsUsingCRN(meta, *vpnGateway.CRN, "", isUserTagType)
	if err != nil {
		log.Printf(
			"Error on get of resource vpc VPN Gateway (%s) tags: %s", d.Id(), err)
	}
	d.Set(isVPNGatewayTags, tags)

	accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *vpnGateway.CRN, "", isAccessTagType)
	if err != nil {
		log.Printf(
			"Error on get of resource VPC VPN Gateway (%s) access tags: %s", d.Id(), err)
	}
	d.Set(isVPNGatewayAccessTags, accesstags)

	controller, err := flex.GetBaseController(meta)
	if err != nil {
		return err
	}
	d.Set(flex.ResourceControllerURL, controller+"/vpc/network/vpngateways")
	d.Set(flex.ResourceName, *vpnGateway.Name)
	d.Set(flex.ResourceCRN, *vpnGateway.CRN)
	d.Set(isVPNGatewayCRN, *vpnGateway.CRN)
	if vpnGateway.ResourceGroup != nil {
		d.Set(flex.ResourceGroupName, vpnGateway.ResourceGroup.Name)
		d.Set(isVPNGatewayResourceGroup, vpnGateway.ResourceGroup.ID)
	}
	d.Set(isVPNGatewayMode, *vpnGateway.Mode)
	if vpnGateway.Members != nil {
		vpcMembersIpsList := make([]map[string]interface{}, 0)
		for _, memberIP := range vpnGateway.Members {
			currentMemberIP := map[string]interface{}{}
			if memberIP.PublicIP != nil {
				currentMemberIP["address"] = *memberIP.PublicIP.Address
				currentMemberIP["role"] = *memberIP.Role
				vpcMembersIpsList = append(vpcMembersIpsList, currentMemberIP)
			}
			if memberIP.PrivateIP != nil && memberIP.PrivateIP.Address != nil {
				currentMemberIP["private_address"] = *memberIP.PrivateIP.Address
			}
		}
		d.Set(isVPNGatewayMembers, vpcMembersIpsList)
	}
	if vpnGateway.CreatedAt != nil {
		d.Set(isVPNGatewayCreatedAt, (vpnGateway.CreatedAt).String())
	}
	if vpnGateway.VPC != nil {
		vpcList := []map[string]interface{}{}
		vpcList = append(vpcList, dataSourceVPNServerCollectionVPNGatewayVpcReferenceToMap(vpnGateway.VPC))
		err = d.Set("vpc", vpcList)
		if err != nil {
			return fmt.Errorf("Error setting the vpc: %s", err)
		}
	}
	return nil
}

func resourceIBMISVPNGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	name := ""
	hasChanged := false

	if d.HasChange(isVPNGatewayName) {
		name = d.Get(isVPNGatewayName).(string)
		hasChanged = true
	}

	err := vpngwUpdate(d, meta, id, name, hasChanged)
	if err != nil {
		return err
	}
	return resourceIBMISVPNGatewayRead(d, meta)
}

func vpngwUpdate(d *schema.ResourceData, meta interface{}, id, name string, hasChanged bool) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	if d.HasChange(isVPNGatewayTags) {
		getVpnGatewayOptions := &vpcv1.GetVPNGatewayOptions{
			ID: &id,
		}
		vpnGatewayIntf, response, err := sess.GetVPNGateway(getVpnGatewayOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error getting Volume : %s\n%s", err, response)
		}
		vpnGateway := vpnGatewayIntf.(*vpcv1.VPNGateway)

		oldList, newList := d.GetChange(isVPNGatewayTags)
		err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *vpnGateway.CRN, "", isUserTagType)
		if err != nil {
			log.Printf(
				"Error on update of resource vpc Vpn Gateway (%s) tags: %s", id, err)
		}
	}
	if d.HasChange(isVPNGatewayAccessTags) {
		getVpnGatewayOptions := &vpcv1.GetVPNGatewayOptions{
			ID: &id,
		}
		vpnGatewayIntf, response, err := sess.GetVPNGateway(getVpnGatewayOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error getting Volume : %s\n%s", err, response)
		}
		vpnGateway := vpnGatewayIntf.(*vpcv1.VPNGateway)

		oldList, newList := d.GetChange(isVPNGatewayAccessTags)
		err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *vpnGateway.CRN, "", isAccessTagType)
		if err != nil {
			log.Printf(
				"Error on update of resource VPC VPN Gateway  (%s) access tags: %s", d.Id(), err)
		}
	}
	if hasChanged {
		options := &vpcv1.UpdateVPNGatewayOptions{
			ID: &id,
		}
		vpnGatewayPatchModel := &vpcv1.VPNGatewayPatch{
			Name: &name,
		}
		vpnGatewayPatch, err := vpnGatewayPatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("[ERROR] Error calling asPatch for VPNGatewayPatch: %s", err)
		}
		options.VPNGatewayPatch = vpnGatewayPatch
		_, response, err := sess.UpdateVPNGateway(options)
		if err != nil {
			return fmt.Errorf("[ERROR] Error updating vpc Vpn Gateway: %s\n%s", err, response)
		}
	}
	return nil
}

func resourceIBMISVPNGatewayDelete(d *schema.ResourceData, meta interface{}) error {

	id := d.Id()

	err := vpngwDelete(d, meta, id)
	if err != nil {
		return err
	}
	return nil
}

func vpngwDelete(d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	getVpnGatewayOptions := &vpcv1.GetVPNGatewayOptions{
		ID: &id,
	}
	_, response, err := sess.GetVPNGateway(getVpnGatewayOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil
		}
		return fmt.Errorf("[ERROR] Error Getting Vpn Gateway (%s): %s\n%s", id, err, response)
	}

	options := &vpcv1.DeleteVPNGatewayOptions{
		ID: &id,
	}
	response, err = sess.DeleteVPNGateway(options)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Deleting Vpn Gateway : %s\n%s", err, response)
	}
	_, err = isWaitForVpnGatewayDeleted(sess, id, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func isWaitForVpnGatewayDeleted(vpnGateway *vpcv1.VpcV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for VPNGateway (%s) to be deleted.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isVPNGatewayDeleting},
		Target:     []string{isVPNGatewayDeleted, ""},
		Refresh:    isVpnGatewayDeleteRefreshFunc(vpnGateway, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isVpnGatewayDeleteRefreshFunc(vpnGateway *vpcv1.VpcV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getVpnGatewayOptions := &vpcv1.GetVPNGatewayOptions{
			ID: &id,
		}
		vpngw, response, err := vpnGateway.GetVPNGateway(getVpnGatewayOptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return "", isVPNGatewayDeleted, nil
			}
			return "", "", fmt.Errorf("[ERROR] Error Getting Vpn Gateway: %s\n%s", err, response)
		}
		return vpngw, isVPNGatewayDeleting, err
	}
}

func resourceIBMISVPNGatewayExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	id := d.Id()

	exists, err := vpngwExists(d, meta, id)
	return exists, err
}

func vpngwExists(d *schema.ResourceData, meta interface{}, id string) (bool, error) {
	sess, err := vpcClient(meta)
	if err != nil {
		return false, err
	}
	getVpnGatewayOptions := &vpcv1.GetVPNGatewayOptions{
		ID: &id,
	}
	_, response, err := sess.GetVPNGateway(getVpnGatewayOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("[ERROR] Error getting Vpn Gatewa: %s\n%s", err, response)
	}
	return true, nil
}

func resourceVPNGatewayRouteFlattenHealthReasons(healthReasons []vpcv1.VPNGatewayHealthReason) (healthReasonsList []map[string]interface{}) {
	healthReasonsList = make([]map[string]interface{}, 0)
	for _, lr := range healthReasons {
		currentLR := map[string]interface{}{}
		if lr.Code != nil && lr.Message != nil {
			currentLR[isInstanceLifecycleReasonsCode] = *lr.Code
			currentLR[isInstanceLifecycleReasonsMessage] = *lr.Message
			if lr.MoreInfo != nil {
				currentLR[isInstanceLifecycleReasonsMoreInfo] = *lr.MoreInfo
			}
			healthReasonsList = append(healthReasonsList, currentLR)
		}
	}
	return healthReasonsList
}

func resourceVPNGatewayFlattenLifecycleReasons(lifecycleReasons []vpcv1.VPNGatewayLifecycleReason) (lifecycleReasonsList []map[string]interface{}) {
	lifecycleReasonsList = make([]map[string]interface{}, 0)
	for _, lr := range lifecycleReasons {
		currentLR := map[string]interface{}{}
		if lr.Code != nil && lr.Message != nil {
			currentLR[isInstanceLifecycleReasonsCode] = *lr.Code
			currentLR[isInstanceLifecycleReasonsMessage] = *lr.Message
			if lr.MoreInfo != nil {
				currentLR[isInstanceLifecycleReasonsMoreInfo] = *lr.MoreInfo
			}
			lifecycleReasonsList = append(lifecycleReasonsList, currentLR)
		}
	}
	return lifecycleReasonsList
}
