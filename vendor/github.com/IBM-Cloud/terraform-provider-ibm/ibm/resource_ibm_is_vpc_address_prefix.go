// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	"github.com/IBM/vpc-go-sdk/vpcclassicv1"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	isVPCAddressPrefixPrefixName = "name"
	isVPCAddressPrefixZoneName   = "zone"
	isVPCAddressPrefixCIDR       = "cidr"
	isVPCAddressPrefixVPCID      = "vpc"
	isVPCAddressPrefixHasSubnets = "has_subnets"
	isVPCAddressPrefixDefault    = "is_default"
)

func resourceIBMISVpcAddressPrefix() *schema.Resource {
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
				ValidateFunc: InvokeValidator("ibm_is_address_prefix", isVPCAddressPrefixPrefixName),
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
				ValidateFunc: InvokeValidator("ibm_is_address_prefix", isVPCAddressPrefixCIDR),
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

			RelatedCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the VPC resource",
			},
		},
	}
}

func resourceIBMISAddressPrefixValidator() *ResourceValidator {

	validateSchema := make([]ValidateSchema, 1)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isVPCAddressPrefixPrefixName,
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
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isVPCAddressPrefixCIDR,
			ValidateFunctionIdentifier: ValidateOverlappingAddress,
			Type:                       TypeString,
			ForceNew:                   true,
			Required:                   true})

	ibmISAddressPrefixResourceValidator := ResourceValidator{ResourceName: "ibm_is_address_prefix", Schema: validateSchema}
	return &ibmISAddressPrefixResourceValidator
}

func resourceIBMISVpcAddressPrefixCreate(d *schema.ResourceData, meta interface{}) error {

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	isDefault := false
	prefixName := d.Get(isVPCAddressPrefixPrefixName).(string)
	zoneName := d.Get(isVPCAddressPrefixZoneName).(string)
	cidr := d.Get(isVPCAddressPrefixCIDR).(string)
	vpcID := d.Get(isVPCAddressPrefixVPCID).(string)
	if isDefaultPrefix, ok := d.GetOk(isVPCAddressPrefixDefault); ok {
		isDefault = isDefaultPrefix.(bool)
	}

	isVPCAddressPrefixKey := "vpc_address_prefix_key_" + vpcID
	ibmMutexKV.Lock(isVPCAddressPrefixKey)
	defer ibmMutexKV.Unlock(isVPCAddressPrefixKey)

	if userDetails.generation == 1 {
		err := classicVpcAddressPrefixCreate(d, meta, prefixName, zoneName, cidr, vpcID)
		if err != nil {
			return err
		}
	} else {
		err := vpcAddressPrefixCreate(d, meta, prefixName, zoneName, cidr, vpcID, isDefault)
		if err != nil {
			return err
		}
	}
	return resourceIBMISVpcAddressPrefixRead(d, meta)
}

func classicVpcAddressPrefixCreate(d *schema.ResourceData, meta interface{}, name, zone, cidr, vpcID string) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	options := &vpcclassicv1.CreateVPCAddressPrefixOptions{
		Name:  &name,
		VPCID: &vpcID,
		CIDR:  &cidr,
		Zone: &vpcclassicv1.ZoneIdentity{
			Name: &zone,
		},
	}
	addrPrefix, response, err := sess.CreateVPCAddressPrefix(options)
	if err != nil {
		return fmt.Errorf("Error while creating VPC Address Prefix %s\n%s", err, response)
	}

	addrPrefixID := *addrPrefix.ID

	d.SetId(fmt.Sprintf("%s/%s", vpcID, addrPrefixID))
	return nil
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
		return fmt.Errorf("Error while creating VPC Address Prefix %s\n%s", err, response)
	}

	addrPrefixID := *addrPrefix.ID
	d.SetId(fmt.Sprintf("%s/%s", vpcID, addrPrefixID))
	return nil
}

func resourceIBMISVpcAddressPrefixRead(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}

	vpcID := parts[0]
	addrPrefixID := parts[1]
	if userDetails.generation == 1 {
		err := classicVpcAddressPrefixGet(d, meta, vpcID, addrPrefixID)
		if err != nil {
			return err
		}
	} else {
		err := vpcAddressPrefixGet(d, meta, vpcID, addrPrefixID)
		if err != nil {
			return err
		}
	}

	return nil
}

func classicVpcAddressPrefixGet(d *schema.ResourceData, meta interface{}, vpcID, addrPrefixID string) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	getvpcAddressPrefixOptions := &vpcclassicv1.GetVPCAddressPrefixOptions{
		VPCID: &vpcID,
		ID:    &addrPrefixID,
	}
	addrPrefix, response, err := sess.GetVPCAddressPrefix(getvpcAddressPrefixOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Getting VPC Address Prefix (%s): %s\n%s", addrPrefixID, err, response)
	}
	d.Set(isVPCAddressPrefixVPCID, vpcID)
	d.Set(isVPCAddressPrefixPrefixName, *addrPrefix.Name)
	if addrPrefix.Zone != nil {
		d.Set(isVPCAddressPrefixZoneName, *addrPrefix.Zone.Name)
	}
	d.Set(isVPCAddressPrefixCIDR, *addrPrefix.CIDR)
	d.Set(isVPCAddressPrefixHasSubnets, *addrPrefix.HasSubnets)
	getVPCOptions := &vpcclassicv1.GetVPCOptions{
		ID: &vpcID,
	}
	vpc, response, err := sess.GetVPC(getVPCOptions)
	if err != nil {
		return fmt.Errorf("Error Getting VPC : %s\n%s", err, response)
	}
	d.Set(RelatedCRN, *vpc.CRN)

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
		return fmt.Errorf("Error Getting VPC Address Prefix (%s): %s\n%s", addrPrefixID, err, response)
	}
	d.Set(isVPCAddressPrefixVPCID, vpcID)
	d.Set(isVPCAddressPrefixDefault, *addrPrefix.IsDefault)
	d.Set(isVPCAddressPrefixPrefixName, *addrPrefix.Name)
	if addrPrefix.Zone != nil {
		d.Set(isVPCAddressPrefixZoneName, *addrPrefix.Zone.Name)
	}
	d.Set(isVPCAddressPrefixCIDR, *addrPrefix.CIDR)
	d.Set(isVPCAddressPrefixHasSubnets, *addrPrefix.HasSubnets)
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

func resourceIBMISVpcAddressPrefixUpdate(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	name := ""
	isDefault := false
	hasNameChanged := false
	hasIsDefaultChanged := false

	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	vpcID := parts[0]
	addrPrefixID := parts[1]

	isVPCAddressPrefixKey := "vpc_address_prefix_key_" + vpcID
	ibmMutexKV.Lock(isVPCAddressPrefixKey)
	defer ibmMutexKV.Unlock(isVPCAddressPrefixKey)

	if d.HasChange(isVPCAddressPrefixPrefixName) {
		name = d.Get(isVPCAddressPrefixPrefixName).(string)
		hasNameChanged = true
	}
	if d.HasChange(isVPCAddressPrefixDefault) {
		isDefault = d.Get(isVPCAddressPrefixDefault).(bool)
		hasIsDefaultChanged = true
	}
	if userDetails.generation == 1 {
		err := classicVpcAddressPrefixUpdate(d, meta, vpcID, addrPrefixID, name, hasNameChanged)
		if err != nil {
			return err
		}
	} else {
		err := vpcAddressPrefixUpdate(d, meta, vpcID, addrPrefixID, name, isDefault, hasNameChanged, hasIsDefaultChanged)
		if err != nil {
			return err
		}
	}

	return resourceIBMISVpcAddressPrefixRead(d, meta)
}

func classicVpcAddressPrefixUpdate(d *schema.ResourceData, meta interface{}, vpcID, addrPrefixID, name string, hasChanged bool) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	if hasChanged {
		updatevpcAddressPrefixoptions := &vpcclassicv1.UpdateVPCAddressPrefixOptions{
			VPCID: &vpcID,
			ID:    &addrPrefixID,
		}
		addressPrefixPatchModel := &vpcclassicv1.AddressPrefixPatch{
			Name: &name,
		}
		addressPrefixPatch, err := addressPrefixPatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("Error calling asPatch for AddressPrefixPatch: %s", err)
		}
		updatevpcAddressPrefixoptions.AddressPrefixPatch = addressPrefixPatch
		_, response, err := sess.UpdateVPCAddressPrefix(updatevpcAddressPrefixoptions)
		if err != nil {
			return fmt.Errorf("Error Updating VPC Address Prefix: %s\n%s", err, response)
		}
	}
	return nil
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
			return fmt.Errorf("Error calling asPatch for AddressPrefixPatch: %s", err)
		}
		updatevpcAddressPrefixoptions.AddressPrefixPatch = addressPrefixPatch
		_, response, err := sess.UpdateVPCAddressPrefix(updatevpcAddressPrefixoptions)
		if err != nil {
			return fmt.Errorf("Error Updating VPC Address Prefix: %s\n%s", err, response)
		}
	}
	return nil
}

func resourceIBMISVpcAddressPrefixDelete(d *schema.ResourceData, meta interface{}) error {

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	vpcID := parts[0]
	addrPrefixID := parts[1]

	isVPCAddressPrefixKey := "vpc_address_prefix_key_" + vpcID
	ibmMutexKV.Lock(isVPCAddressPrefixKey)
	defer ibmMutexKV.Unlock(isVPCAddressPrefixKey)

	if userDetails.generation == 1 {
		err := classicVpcAddressPrefixDelete(d, meta, vpcID, addrPrefixID)
		if err != nil {
			return err
		}
	} else {
		err := vpcAddressPrefixDelete(d, meta, vpcID, addrPrefixID)
		if err != nil {
			return err
		}
	}

	d.SetId("")
	return nil
}

func classicVpcAddressPrefixDelete(d *schema.ResourceData, meta interface{}, vpcID, addrPrefixID string) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}

	getvpcAddressPrefixOptions := &vpcclassicv1.GetVPCAddressPrefixOptions{
		VPCID: &vpcID,
		ID:    &addrPrefixID,
	}
	_, response, err := sess.GetVPCAddressPrefix(getvpcAddressPrefixOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil
		}
		return fmt.Errorf("Error Getting VPC Address Prefix (%s): %s\n%s", addrPrefixID, err, response)
	}
	deletevpcAddressPrefixOptions := &vpcclassicv1.DeleteVPCAddressPrefixOptions{
		VPCID: &vpcID,
		ID:    &addrPrefixID,
	}
	response, err = sess.DeleteVPCAddressPrefix(deletevpcAddressPrefixOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil
		}
		return fmt.Errorf("Error Deleting VPC Address Prefix (%s): %s\n%s", addrPrefixID, err, response)
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
		return fmt.Errorf("Error Getting VPC Address Prefix (%s): %s\n%s", addrPrefixID, err, response)
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
		return fmt.Errorf("Error Deleting VPC Address Prefix (%s): %s\n%s", addrPrefixID, err, response)
	}
	d.SetId("")
	return nil
}

func resourceIBMISVpcAddressPrefixExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return false, err
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return false, err
	}
	vpcID := parts[0]
	addrPrefixID := parts[1]

	if userDetails.generation == 1 {
		exists, err := classicVpcAddressPrefixExists(d, meta, vpcID, addrPrefixID)
		return exists, err
	} else {
		exists, err := vpcAddressPrefixExists(d, meta, vpcID, addrPrefixID)
		return exists, err
	}
}

func classicVpcAddressPrefixExists(d *schema.ResourceData, meta interface{}, vpcID, addrPrefixID string) (bool, error) {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return false, err
	}
	getvpcAddressPrefixOptions := &vpcclassicv1.GetVPCAddressPrefixOptions{
		VPCID: &vpcID,
		ID:    &addrPrefixID,
	}
	_, response, err := sess.GetVPCAddressPrefix(getvpcAddressPrefixOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("Error getting VPC Address Prefix: %s\n%s", err, response)
	}
	return true, nil
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
		return false, fmt.Errorf("Error getting VPC Address Prefix: %s\n%s", err, response)
	}
	return true, nil
}
