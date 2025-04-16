// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isVPCAddressPrefixPrefixName = "name"
	isVPCAddressPrefixZoneName   = "zone"
	isVPCAddressPrefixCIDR       = "cidr"
	isVPCAddressPrefixVPCID      = "vpc"
	isVPCAddressPrefixHasSubnets = "has_subnets"
	isVPCAddressPrefixDefault    = "is_default"
	isAddressPrefix              = "address_prefix"
)

func ResourceIBMISVpcAddressPrefix() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMISVpcAddressPrefixCreate,
		Read:     resourceIBMISVpcAddressPrefixRead,
		Update:   resourceIBMISVpcAddressPrefixUpdate,
		Delete:   resourceIBMISVpcAddressPrefixDelete,
		Exists:   resourceIBMISVpcAddressPrefixExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			isVPCAddressPrefixPrefixName: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validate.InvokeValidator("ibm_is_address_prefix", isVPCAddressPrefixPrefixName),
				Description:  "Name",
			},
			isVPCAddressPrefixZoneName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone name",
			},

			isVPCAddressPrefixCIDR: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_address_prefix", isVPCAddressPrefixCIDR),
				Description:  "CIDIR address prefix",
			},
			isVPCAddressPrefixDefault: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Is default prefix for this zone in this VPC",
			},

			isVPCAddressPrefixVPCID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPC id",
			},

			isVPCAddressPrefixHasSubnets: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Boolean value, set to true if VPC instance have subnets",
			},

			flex.RelatedCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the VPC resource",
			},

			isAddressPrefix: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the address prefix",
			},
		},
	}
}

func ResourceIBMISAddressPrefixValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isVPCAddressPrefixPrefixName,
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isVPCAddressPrefixCIDR,
			ValidateFunctionIdentifier: validate.ValidateOverlappingAddress,
			Type:                       validate.TypeString,
			ForceNew:                   true,
			Required:                   true})

	ibmISAddressPrefixResourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_address_prefix", Schema: validateSchema}
	return &ibmISAddressPrefixResourceValidator
}

func resourceIBMISVpcAddressPrefixCreate(d *schema.ResourceData, meta interface{}) error {

	isDefault := false
	prefixName := d.Get(isVPCAddressPrefixPrefixName).(string)
	zoneName := d.Get(isVPCAddressPrefixZoneName).(string)
	cidr := d.Get(isVPCAddressPrefixCIDR).(string)
	vpcID := d.Get(isVPCAddressPrefixVPCID).(string)
	if isDefaultPrefix, ok := d.GetOk(isVPCAddressPrefixDefault); ok {
		isDefault = isDefaultPrefix.(bool)
	}

	isVPCAddressPrefixKey := "vpc_address_prefix_key_" + vpcID
	conns.IbmMutexKV.Lock(isVPCAddressPrefixKey)
	defer conns.IbmMutexKV.Unlock(isVPCAddressPrefixKey)

	err := vpcAddressPrefixCreate(d, meta, prefixName, zoneName, cidr, vpcID, isDefault)
	if err != nil {
		return err
	}
	return resourceIBMISVpcAddressPrefixRead(d, meta)
}

func vpcAddressPrefixCreate(d *schema.ResourceData, meta interface{}, name, zone, cidr, vpcID string, isDefault bool) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	options := &vpcv1.CreateVPCAddressPrefixOptions{
		Name:      &name,
		VPCID:     &vpcID,
		CIDR:      &cidr,
		IsDefault: &isDefault,
		Zone: &vpcv1.ZoneIdentity{
			Name: &zone,
		},
	}
	addrPrefix, response, err := sess.CreateVPCAddressPrefix(options)
	if err != nil {
		return fmt.Errorf("[ERROR] Error while creating VPC Address Prefix %s\n%s", err, response)
	}

	addrPrefixID := *addrPrefix.ID
	d.SetId(fmt.Sprintf("%s/%s", vpcID, addrPrefixID))
	d.Set(isAddressPrefix, addrPrefixID)
	return nil
}

func resourceIBMISVpcAddressPrefixRead(d *schema.ResourceData, meta interface{}) error {
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}

	vpcID := parts[0]
	addrPrefixID := parts[1]
	error := vpcAddressPrefixGet(d, meta, vpcID, addrPrefixID)
	if error != nil {
		return error
	}

	return nil
}

func vpcAddressPrefixGet(d *schema.ResourceData, meta interface{}, vpcID, addrPrefixID string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	getvpcAddressPrefixOptions := &vpcv1.GetVPCAddressPrefixOptions{
		VPCID: &vpcID,
		ID:    &addrPrefixID,
	}
	addrPrefix, response, err := sess.GetVPCAddressPrefix(getvpcAddressPrefixOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error Getting VPC Address Prefix (%s): %s\n%s", addrPrefixID, err, response)
	}
	d.Set(isVPCAddressPrefixVPCID, vpcID)
	d.Set(isVPCAddressPrefixDefault, *addrPrefix.IsDefault)
	d.Set(isVPCAddressPrefixPrefixName, *addrPrefix.Name)
	if addrPrefix.Zone != nil {
		d.Set(isVPCAddressPrefixZoneName, *addrPrefix.Zone.Name)
	}
	d.Set(isVPCAddressPrefixCIDR, *addrPrefix.CIDR)
	d.Set(isVPCAddressPrefixHasSubnets, *addrPrefix.HasSubnets)
	d.Set(isAddressPrefix, addrPrefixID)
	getVPCOptions := &vpcv1.GetVPCOptions{
		ID: &vpcID,
	}
	vpc, response, err := sess.GetVPC(getVPCOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Getting VPC : %s\n%s", err, response)
	}
	d.Set(flex.RelatedCRN, *vpc.CRN)

	return nil
}

func resourceIBMISVpcAddressPrefixUpdate(d *schema.ResourceData, meta interface{}) error {

	name := ""
	isDefault := false
	hasNameChanged := false
	hasIsDefaultChanged := false

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}
	vpcID := parts[0]
	addrPrefixID := parts[1]

	isVPCAddressPrefixKey := "vpc_address_prefix_key_" + vpcID
	conns.IbmMutexKV.Lock(isVPCAddressPrefixKey)
	defer conns.IbmMutexKV.Unlock(isVPCAddressPrefixKey)

	if d.HasChange(isVPCAddressPrefixPrefixName) {
		name = d.Get(isVPCAddressPrefixPrefixName).(string)
		hasNameChanged = true
	}
	if d.HasChange(isVPCAddressPrefixDefault) {
		isDefault = d.Get(isVPCAddressPrefixDefault).(bool)
		hasIsDefaultChanged = true
	}
	error := vpcAddressPrefixUpdate(d, meta, vpcID, addrPrefixID, name, isDefault, hasNameChanged, hasIsDefaultChanged)
	if error != nil {
		return error
	}

	return resourceIBMISVpcAddressPrefixRead(d, meta)
}

func vpcAddressPrefixUpdate(d *schema.ResourceData, meta interface{}, vpcID, addrPrefixID, name string, isDefault, hasNameChanged, hasIsDefaultChanged bool) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	if hasNameChanged || hasIsDefaultChanged {
		updatevpcAddressPrefixoptions := &vpcv1.UpdateVPCAddressPrefixOptions{
			VPCID: &vpcID,
			ID:    &addrPrefixID,
		}

		addressPrefixPatchModel := &vpcv1.AddressPrefixPatch{}
		if hasNameChanged {
			addressPrefixPatchModel.Name = &name
		}
		if hasIsDefaultChanged {
			addressPrefixPatchModel.IsDefault = &isDefault
		}
		addressPrefixPatch, err := addressPrefixPatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("[ERROR] Error calling asPatch for AddressPrefixPatch: %s", err)
		}
		updatevpcAddressPrefixoptions.AddressPrefixPatch = addressPrefixPatch
		_, response, err := sess.UpdateVPCAddressPrefix(updatevpcAddressPrefixoptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Updating VPC Address Prefix: %s\n%s", err, response)
		}
	}
	return nil
}

func resourceIBMISVpcAddressPrefixDelete(d *schema.ResourceData, meta interface{}) error {

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}
	vpcID := parts[0]
	addrPrefixID := parts[1]

	isVPCAddressPrefixKey := "vpc_address_prefix_key_" + vpcID
	conns.IbmMutexKV.Lock(isVPCAddressPrefixKey)
	defer conns.IbmMutexKV.Unlock(isVPCAddressPrefixKey)

	error := vpcAddressPrefixDelete(d, meta, vpcID, addrPrefixID)
	if error != nil {
		return error
	}

	d.SetId("")
	return nil
}

func vpcAddressPrefixDelete(d *schema.ResourceData, meta interface{}, vpcID, addrPrefixID string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	getvpcAddressPrefixOptions := &vpcv1.GetVPCAddressPrefixOptions{
		VPCID: &vpcID,
		ID:    &addrPrefixID,
	}
	_, response, err := sess.GetVPCAddressPrefix(getvpcAddressPrefixOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil
		}
		return fmt.Errorf("[ERROR] Error Getting VPC Address Prefix (%s): %s\n%s", addrPrefixID, err, response)
	}

	deletevpcAddressPrefixOptions := &vpcv1.DeleteVPCAddressPrefixOptions{
		VPCID: &vpcID,
		ID:    &addrPrefixID,
	}
	response, err = sess.DeleteVPCAddressPrefix(deletevpcAddressPrefixOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil
		}
		return fmt.Errorf("[ERROR] Error Deleting VPC Address Prefix (%s): %s\n%s", addrPrefixID, err, response)
	}
	d.SetId("")
	return nil
}

func resourceIBMISVpcAddressPrefixExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	parts, err := flex.IdParts(d.Id())
	if len(parts) != 2 {
		return false, fmt.Errorf("[ERROR] Incorrect ID %s: ID should be a combination of vpcID/addrPrefixID", d.Id())
	}
	if err != nil {
		return false, err
	}
	vpcID := parts[0]
	addrPrefixID := parts[1]

	exists, err := vpcAddressPrefixExists(d, meta, vpcID, addrPrefixID)
	return exists, err
}

func vpcAddressPrefixExists(d *schema.ResourceData, meta interface{}, vpcID, addrPrefixID string) (bool, error) {
	sess, err := vpcClient(meta)
	if err != nil {
		return false, err
	}
	getvpcAddressPrefixOptions := &vpcv1.GetVPCAddressPrefixOptions{
		VPCID: &vpcID,
		ID:    &addrPrefixID,
	}
	_, response, err := sess.GetVPCAddressPrefix(getvpcAddressPrefixOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("[ERROR] Error getting VPC Address Prefix: %s\n%s", err, response)
	}
	return true, nil
}
