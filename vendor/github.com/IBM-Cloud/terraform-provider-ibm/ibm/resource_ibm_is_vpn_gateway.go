// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcclassicv1"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	isVPNGatewayName              = "name"
	isVPNGatewayResourceGroup     = "resource_group"
	isVPNGatewayMode              = "mode"
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
)

func resourceIBMISVPNGateway() *schema.Resource {
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

		CustomizeDiff: customdiff.Sequence(
			func(diff *schema.ResourceDiff, v interface{}) error {
				return resourceTagsCustomizeDiff(diff)
			},
		),

		Schema: map[string]*schema.Schema{

			isVPNGatewayName: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: InvokeValidator("ibm_is_route", isVPNGatewayName),
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
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: InvokeValidator("ibm_is_vpn_gateway", "tag")},
				Set:         resourceIBMVPCHash,
				Description: "VPN Gateway tags list",
			},

			ResourceControllerURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the IBM Cloud dashboard that can be used to explore and view details about this instance",
			},

			ResourceName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the resource",
			},

			ResourceCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},

			ResourceStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the resource",
			},

			ResourceGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group name in which resource is provisioned",
			},
			isVPNGatewayCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Created Time of the VPN Gateway",
			},
			isVPNGatewayMode: {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "route",
				ValidateFunc: InvokeValidator("ibm_is_vpn_gateway", isVPNGatewayMode),
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
		},
	}
}

func resourceIBMISVPNGatewayValidator() *ResourceValidator {

	modeCheckTypes := "route,policy"
	validateSchema := make([]ValidateSchema, 2)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isVPNGatewayName,
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Required:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isVPNGatewayMode,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   false,
			AllowedValues:              modeCheckTypes})

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "tag",
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9:_ .-]+$`,
			MinValueLength:             1,
			MaxValueLength:             128})

	ibmISVPNGatewayResourceValidator := ResourceValidator{ResourceName: "ibm_is_vpn_gateway", Schema: validateSchema}
	return &ibmISVPNGatewayResourceValidator
}

func resourceIBMISVPNGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] VPNGateway create")
	name := d.Get(isVPNGatewayName).(string)
	subnetID := d.Get(isVPNGatewaySubnet).(string)
	mode := d.Get(isVPNGatewayMode).(string)

	if userDetails.generation == 1 {
		err := classicVpngwCreate(d, meta, name, subnetID)
		if err != nil {
			return err
		}
	} else {
		err := vpngwCreate(d, meta, name, subnetID, mode)
		if err != nil {
			return err
		}
	}
	return resourceIBMISVPNGatewayRead(d, meta)
}

func classicVpngwCreate(d *schema.ResourceData, meta interface{}, name, subnetID string) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	vpnGatewayPrototype := &vpcclassicv1.VPNGatewayPrototype{
		Subnet: &vpcclassicv1.SubnetIdentity{
			ID: &subnetID,
		},
		Name: &name,
	}
	options := &vpcclassicv1.CreateVPNGatewayOptions{
		VPNGatewayPrototype: vpnGatewayPrototype,
	}

	if rgrp, ok := d.GetOk(isVPNGatewayResourceGroup); ok {
		rg := rgrp.(string)
		vpnGatewayPrototype.ResourceGroup = &vpcclassicv1.ResourceGroupIdentity{
			ID: &rg,
		}
	}

	vpnGatewayIntf, response, err := sess.CreateVPNGateway(options)
	if err != nil {
		return fmt.Errorf("[DEBUG] Create vpc VPN Gateway %s\n%s", err, response)
	}
	vpnGateway := vpnGatewayIntf.(*vpcclassicv1.VPNGateway)
	_, err = isWaitForClassicVpnGatewayAvailable(sess, *vpnGateway.ID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	d.SetId(*vpnGateway.ID)
	log.Printf("[INFO] VPNGateway : %s", *vpnGateway.ID)

	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk(isVPNGatewayTags); ok || v != "" {
		oldList, newList := d.GetChange(isVPNGatewayTags)
		err = UpdateTagsUsingCRN(oldList, newList, meta, *vpnGateway.CRN)
		if err != nil {
			log.Printf(
				"Error on create of resource vpc VPN Gateway (%s) tags: %s", d.Id(), err)
		}
	}
	return nil
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

	_, err = isWaitForVpnGatewayAvailable(sess, *vpnGateway.ID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	d.SetId(*vpnGateway.ID)
	log.Printf("[INFO] VPNGateway : %s", *vpnGateway.ID)

	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk(isVPNGatewayTags); ok || v != "" {
		oldList, newList := d.GetChange(isVPNGatewayTags)
		err = UpdateTagsUsingCRN(oldList, newList, meta, *vpnGateway.CRN)
		if err != nil {
			log.Printf(
				"Error on create of resource vpc VPN Gateway (%s) tags: %s", d.Id(), err)
		}
	}
	return nil
}

func isWaitForClassicVpnGatewayAvailable(vpnGateway *vpcclassicv1.VpcClassicV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for vpn gateway (%s) to be available.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isVPNGatewayProvisioning},
		Target:     []string{isVPNGatewayProvisioningDone, ""},
		Refresh:    isClassicVpnGatewayRefreshFunc(vpnGateway, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isClassicVpnGatewayRefreshFunc(vpnGateway *vpcclassicv1.VpcClassicV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getVpnGatewayOptions := &vpcclassicv1.GetVPNGatewayOptions{
			ID: &id,
		}
		vpnGatewayIntf, response, err := vpnGateway.GetVPNGateway(getVpnGatewayOptions)
		if err != nil {
			return nil, "", fmt.Errorf("Error Getting Vpn Gateway: %s\n%s", err, response)
		}
		vpnGateway := vpnGatewayIntf.(*vpcclassicv1.VPNGateway)

		if *vpnGateway.Status == "available" || *vpnGateway.Status == "failed" || *vpnGateway.Status == "running" {
			return vpnGateway, isVPNGatewayProvisioningDone, nil
		}

		return vpnGateway, isVPNGatewayProvisioning, nil
	}
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
			return nil, "", fmt.Errorf("Error Getting Vpn Gateway: %s\n%s", err, response)
		}
		vpnGateway := vpnGatewayIntf.(*vpcv1.VPNGateway)

		if *vpnGateway.Status == "available" || *vpnGateway.Status == "failed" || *vpnGateway.Status == "running" {
			return vpnGateway, isVPNGatewayProvisioningDone, nil
		}

		return vpnGateway, isVPNGatewayProvisioning, nil
	}
}

func resourceIBMISVPNGatewayRead(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	id := d.Id()
	if userDetails.generation == 1 {
		err := classicVpngwGet(d, meta, id)
		if err != nil {
			return err
		}
	} else {
		err := vpngwGet(d, meta, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func classicVpngwGet(d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	getVpnGatewayOptions := &vpcclassicv1.GetVPNGatewayOptions{
		ID: &id,
	}
	vpnGatewayIntf, response, err := sess.GetVPNGateway(getVpnGatewayOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Getting Vpn Gateway (%s): %s\n%s", id, err, response)
	}
	vpnGateway := vpnGatewayIntf.(*vpcclassicv1.VPNGateway)
	d.Set(isVPNGatewayName, *vpnGateway.Name)
	d.Set(isVPNGatewaySubnet, *vpnGateway.Subnet.ID)
	d.Set(isVPNGatewayStatus, *vpnGateway.Status)
	members := []vpcclassicv1.VPNGatewayMember{}
	for _, member := range vpnGateway.Members {
		members = append(members, member)
	}
	if len(members) > 0 {
		d.Set(isVPNGatewayPublicIPAddress, *members[0].PublicIP.Address)
	}
	if len(members) > 1 {
		d.Set(isVPNGatewayPublicIPAddress2, *members[1].PublicIP.Address)
	}
	tags, err := GetTagsUsingCRN(meta, *vpnGateway.CRN)
	if err != nil {
		log.Printf(
			"Error on get of resource vpc VPN Gateway (%s) tags: %s", d.Id(), err)
	}
	d.Set(isVPNGatewayTags, tags)
	controller, err := getBaseController(meta)
	if err != nil {
		return err
	}
	d.Set(ResourceControllerURL, controller+"/vpc/network/vpngateways")
	d.Set(ResourceName, *vpnGateway.Name)
	d.Set(ResourceCRN, *vpnGateway.CRN)
	d.Set(ResourceStatus, *vpnGateway.Status)
	if vpnGateway.ResourceGroup != nil {
		d.Set(ResourceGroupName, *vpnGateway.ResourceGroup.ID)
		d.Set(isVPNGatewayResourceGroup, *vpnGateway.ResourceGroup.ID)
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
		return fmt.Errorf("Error Getting Vpn Gateway (%s): %s\n%s", id, err, response)
	}
	vpnGateway := vpnGatewayIntf.(*vpcv1.VPNGateway)

	d.Set(isVPNGatewayName, *vpnGateway.Name)
	d.Set(isVPNGatewaySubnet, *vpnGateway.Subnet.ID)
	d.Set(isVPNGatewayStatus, *vpnGateway.Status)
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
	tags, err := GetTagsUsingCRN(meta, *vpnGateway.CRN)
	if err != nil {
		log.Printf(
			"Error on get of resource vpc VPN Gateway (%s) tags: %s", d.Id(), err)
	}
	d.Set(isVPNGatewayTags, tags)
	controller, err := getBaseController(meta)
	if err != nil {
		return err
	}
	d.Set(ResourceControllerURL, controller+"/vpc/network/vpngateways")
	d.Set(ResourceName, *vpnGateway.Name)
	d.Set(ResourceCRN, *vpnGateway.CRN)
	d.Set(ResourceStatus, *vpnGateway.Status)
	if vpnGateway.ResourceGroup != nil {
		d.Set(ResourceGroupName, *vpnGateway.ResourceGroup.Name)
		d.Set(isVPNGatewayResourceGroup, *vpnGateway.ResourceGroup.ID)
	}
	d.Set(isVPNGatewayMode, *vpnGateway.Mode)
	if vpnGateway.Members != nil {
		vpcMembersIpsList := make([]map[string]interface{}, 0)
		for _, memberIP := range vpnGateway.Members {
			currentMemberIP := map[string]interface{}{}
			if memberIP.PublicIP != nil {
				currentMemberIP["address"] = *memberIP.PublicIP.Address
				currentMemberIP["role"] = *memberIP.Role
				currentMemberIP["status"] = *memberIP.Status
				vpcMembersIpsList = append(vpcMembersIpsList, currentMemberIP)
			}
			if memberIP.PrivateIP != nil {
				currentMemberIP["private_address"] = *memberIP.PrivateIP.Address
			}
		}
		d.Set(isVPNGatewayMembers, vpcMembersIpsList)
	}
	if vpnGateway.CreatedAt != nil {
		d.Set(isVPNGatewayCreatedAt, (vpnGateway.CreatedAt).String())
	}
	return nil
}

func resourceIBMISVPNGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	id := d.Id()
	name := ""
	hasChanged := false

	if d.HasChange(isVPNGatewayName) {
		name = d.Get(isVPNGatewayName).(string)
		hasChanged = true
	}

	if userDetails.generation == 1 {
		err := classicVpngwUpdate(d, meta, id, name, hasChanged)
		if err != nil {
			return err
		}
	} else {
		err := vpngwUpdate(d, meta, id, name, hasChanged)
		if err != nil {
			return err
		}
	}
	return resourceIBMISVPNGatewayRead(d, meta)
}

func classicVpngwUpdate(d *schema.ResourceData, meta interface{}, id, name string, hasChanged bool) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	if d.HasChange(isVPNGatewayTags) {
		getVpnGatewayOptions := &vpcclassicv1.GetVPNGatewayOptions{
			ID: &id,
		}
		vpnGatewayIntf, response, err := sess.GetVPNGateway(getVpnGatewayOptions)
		if err != nil {
			return fmt.Errorf("Error getting Volume : %s\n%s", err, response)
		}
		vpnGateway := vpnGatewayIntf.(*vpcclassicv1.VPNGateway)

		oldList, newList := d.GetChange(isVPNGatewayTags)
		err = UpdateTagsUsingCRN(oldList, newList, meta, *vpnGateway.CRN)
		if err != nil {
			log.Printf(
				"Error on update of resource vpc Vpn Gateway (%s) tags: %s", id, err)
		}
	}
	if hasChanged {
		options := &vpcclassicv1.UpdateVPNGatewayOptions{
			ID: &id,
		}
		vpnGatewayPatchModel := &vpcclassicv1.VPNGatewayPatch{
			Name: &name,
		}
		vpnGatewayPatch, err := vpnGatewayPatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("Error calling asPatch for VPNGatewayPatch: %s", err)
		}
		options.VPNGatewayPatch = vpnGatewayPatch
		_, response, err := sess.UpdateVPNGateway(options)
		if err != nil {
			return fmt.Errorf("Error updating vpc Vpn Gateway: %s\n%s", err, response)
		}
	}
	return nil
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
			return fmt.Errorf("Error getting Volume : %s\n%s", err, response)
		}
		vpnGateway := vpnGatewayIntf.(*vpcv1.VPNGateway)

		oldList, newList := d.GetChange(isVPNGatewayTags)
		err = UpdateTagsUsingCRN(oldList, newList, meta, *vpnGateway.CRN)
		if err != nil {
			log.Printf(
				"Error on update of resource vpc Vpn Gateway (%s) tags: %s", id, err)
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
			return fmt.Errorf("Error calling asPatch for VPNGatewayPatch: %s", err)
		}
		options.VPNGatewayPatch = vpnGatewayPatch
		_, response, err := sess.UpdateVPNGateway(options)
		if err != nil {
			return fmt.Errorf("Error updating vpc Vpn Gateway: %s\n%s", err, response)
		}
	}
	return nil
}

func resourceIBMISVPNGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	id := d.Id()
	if userDetails.generation == 1 {
		err := classicVpngwDelete(d, meta, id)
		if err != nil {
			return err
		}
	} else {
		err := vpngwDelete(d, meta, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func classicVpngwDelete(d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}

	getVpnGatewayOptions := &vpcclassicv1.GetVPNGatewayOptions{
		ID: &id,
	}
	_, response, err := sess.GetVPNGateway(getVpnGatewayOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil
		}
		return fmt.Errorf("Error Getting Vpn Gateway (%s): %s\n%s", id, err, response)
	}

	options := &vpcclassicv1.DeleteVPNGatewayOptions{
		ID: &id,
	}
	response, err = sess.DeleteVPNGateway(options)
	if err != nil {
		return fmt.Errorf("Error Deleting Vpn Gateway : %s\n%s", err, response)
	}
	_, err = isWaitForClassicVpnGatewayDeleted(sess, id, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}
	d.SetId("")
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
		return fmt.Errorf("Error Getting Vpn Gateway (%s): %s\n%s", id, err, response)
	}

	options := &vpcv1.DeleteVPNGatewayOptions{
		ID: &id,
	}
	response, err = sess.DeleteVPNGateway(options)
	if err != nil {
		return fmt.Errorf("Error Deleting Vpn Gateway : %s\n%s", err, response)
	}
	_, err = isWaitForVpnGatewayDeleted(sess, id, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func isWaitForClassicVpnGatewayDeleted(vpnGateway *vpcclassicv1.VpcClassicV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for VPNGateway (%s) to be deleted.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isVPNGatewayDeleting},
		Target:     []string{isVPNGatewayDeleted, ""},
		Refresh:    isClassicVpnGatewayDeleteRefreshFunc(vpnGateway, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isClassicVpnGatewayDeleteRefreshFunc(vpnGateway *vpcclassicv1.VpcClassicV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getVpnGatewayOptions := &vpcclassicv1.GetVPNGatewayOptions{
			ID: &id,
		}
		vpngw, response, err := vpnGateway.GetVPNGateway(getVpnGatewayOptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return "", isVPNGatewayDeleted, nil
			}
			return "", "", fmt.Errorf("Error Getting Vpn Gateway: %s\n%s", err, response)
		}
		return vpngw, isVPNGatewayDeleting, err
	}
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
			return "", "", fmt.Errorf("Error Getting Vpn Gateway: %s\n%s", err, response)
		}
		return vpngw, isVPNGatewayDeleting, err
	}
}

func resourceIBMISVPNGatewayExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return false, err
	}
	id := d.Id()
	if userDetails.generation == 1 {
		exists, err := classicVpngwExists(d, meta, id)
		return exists, err
	} else {
		exists, err := vpngwExists(d, meta, id)
		return exists, err
	}
}

func classicVpngwExists(d *schema.ResourceData, meta interface{}, id string) (bool, error) {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return false, err
	}
	getVpnGatewayOptions := &vpcclassicv1.GetVPNGatewayOptions{
		ID: &id,
	}
	_, response, err := sess.GetVPNGateway(getVpnGatewayOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("Error getting Vpn Gatewa: %s\n%s", err, response)
	}
	return true, nil
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
		return false, fmt.Errorf("Error getting Vpn Gatewa: %s\n%s", err, response)
	}
	return true, nil
}
