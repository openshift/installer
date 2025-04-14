// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isIKEName              = "name"
	isIKEAuthenticationAlg = "authentication_algorithm"
	isIKEEncryptionAlg     = "encryption_algorithm"
	isIKEDhGroup           = "dh_group"
	isIKEVERSION           = "ike_version"
	isIKEKeyLifeTime       = "key_lifetime"
	isIKEResourceGroup     = "resource_group"
	isIKENegotiationMode   = "negotiation_mode"
	isIKEVPNConnections    = "vpn_connections"
	isIKEVPNConnectionName = "name"
	isIKEVPNConnectionId   = "id"
	isIKEVPNConnectionHref = "href"
	isIKEHref              = "href"
)

func ResourceIBMISIKEPolicy() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMISIKEPolicyCreate,
		Read:     resourceIBMISIKEPolicyRead,
		Update:   resourceIBMISIKEPolicyUpdate,
		Delete:   resourceIBMISIKEPolicyDelete,
		Exists:   resourceIBMISIKEPolicyExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			isIKEName: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_ike_policy", isIKEName),
				Description:  "IKE name",
			},

			isIKEAuthenticationAlg: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_ike_policy", isIKEAuthenticationAlg),
				Description:  "Authentication algorithm type",
			},

			isIKEEncryptionAlg: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_ike_policy", isIKEEncryptionAlg),
				Description:  "Encryption alogorithm type",
			},

			isIKEDhGroup: {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_ike_policy", isIKEDhGroup),
				Description:  "IKE DH group",
			},

			isIKEResourceGroup: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: "IKE resource group ID",
			},

			isIKEKeyLifeTime: {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      28800,
				ValidateFunc: validate.ValidateKeyLifeTime,
				Description:  "IKE Key lifetime",
			},

			isIKEVERSION: {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_ike_policy", isIKEVERSION),
				Description:  "IKE version",
			},

			isIKENegotiationMode: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IKE negotiation mode",
			},

			isIKEHref: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IKE href value",
			},

			isIKEVPNConnections: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isIKEVPNConnectionName: {
							Type:     schema.TypeString,
							Computed: true,
						},
						isIKEVPNConnectionId: {
							Type:     schema.TypeString,
							Computed: true,
						},
						isIKEVPNConnectionHref: {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
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

			flex.ResourceGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group name in which resource is provisioned",
			},
		},
	}
}

func ResourceIBMISIKEValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)
	authentication_algorithm := "md5, sha1, sha256, sha512, sha384"
	encryption_algorithm := "triple_des, aes128, aes192, aes256"
	dh_group := "2, 5, 14, 19, 15, 16, 17, 18, 20, 21, 22, 23, 24, 31"
	ike_version := "1, 2"
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isIKEName,
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isIKEAuthenticationAlg,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              authentication_algorithm})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isIKEEncryptionAlg,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              encryption_algorithm})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isIKEDhGroup,
			ValidateFunctionIdentifier: validate.ValidateAllowedIntValue,
			Type:                       validate.TypeInt,
			Required:                   true,
			AllowedValues:              dh_group})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isIKEVERSION,
			ValidateFunctionIdentifier: validate.ValidateAllowedIntValue,
			Type:                       validate.TypeInt,
			Optional:                   true,
			AllowedValues:              ike_version})

	ibmISIKEResourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_ike_policy", Schema: validateSchema}
	return &ibmISIKEResourceValidator
}

func resourceIBMISIKEPolicyCreate(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] IKE Policy create")
	name := d.Get(isIKEName).(string)
	authenticationAlg := d.Get(isIKEAuthenticationAlg).(string)
	encryptionAlg := d.Get(isIKEEncryptionAlg).(string)
	dhGroup := int64(d.Get(isIKEDhGroup).(int))

	err := ikepCreate(d, meta, authenticationAlg, encryptionAlg, name, dhGroup)
	if err != nil {
		return err
	}
	return resourceIBMISIKEPolicyRead(d, meta)
}

func ikepCreate(d *schema.ResourceData, meta interface{}, authenticationAlg, encryptionAlg, name string, dhGroup int64) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	options := &vpcv1.CreateIkePolicyOptions{
		AuthenticationAlgorithm: &authenticationAlg,
		EncryptionAlgorithm:     &encryptionAlg,
		DhGroup:                 &dhGroup,
		Name:                    &name,
	}

	if keylt, ok := d.GetOk(isIKEKeyLifeTime); ok {
		keyLifetime := int64(keylt.(int))
		options.KeyLifetime = &keyLifetime
	} else {
		keyLifetime := int64(28800)
		options.KeyLifetime = &keyLifetime
	}

	if ikev, ok := d.GetOk(isIKEVERSION); ok {
		ikeVersion := int64(ikev.(int))
		options.IkeVersion = &ikeVersion
	}

	if rgrp, ok := d.GetOk(isIKEResourceGroup); ok {
		rg := rgrp.(string)
		options.ResourceGroup = &vpcv1.ResourceGroupIdentity{
			ID: &rg,
		}
	}
	ike, response, err := sess.CreateIkePolicy(options)
	if err != nil {
		return fmt.Errorf("[DEBUG] ike policy err %s\n%s", err, response)
	}
	d.SetId(*ike.ID)
	log.Printf("[INFO] ike policy : %s", *ike.ID)
	return nil
}

func resourceIBMISIKEPolicyRead(d *schema.ResourceData, meta interface{}) error {

	id := d.Id()
	return ikepGet(d, meta, id)
}

func ikepGet(d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	getikepoptions := &vpcv1.GetIkePolicyOptions{
		ID: &id,
	}
	ike, response, err := sess.GetIkePolicy(getikepoptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error getting IKE Policy(%s): %s\n%s", id, err, response)
	}

	d.Set(isIKEName, *ike.Name)
	d.Set(isIKEAuthenticationAlg, *ike.AuthenticationAlgorithm)
	d.Set(isIKEEncryptionAlg, *ike.EncryptionAlgorithm)
	if ike.ResourceGroup != nil {
		d.Set(isIKEResourceGroup, *ike.ResourceGroup.ID)
		d.Set(flex.ResourceGroupName, *ike.ResourceGroup.Name)
	} else {
		d.Set(isIKEResourceGroup, nil)
	}
	if ike.KeyLifetime != nil {
		d.Set(isIKEKeyLifeTime, *ike.KeyLifetime)
	}
	d.Set(isIKEHref, *ike.Href)
	d.Set(isIKENegotiationMode, *ike.NegotiationMode)
	d.Set(isIKEVERSION, *ike.IkeVersion)
	d.Set(isIKEDhGroup, *ike.DhGroup)
	connList := make([]map[string]interface{}, 0)
	if ike.Connections != nil && len(ike.Connections) > 0 {
		for _, connection := range ike.Connections {
			conn := map[string]interface{}{}
			conn[isIKEVPNConnectionName] = *connection.Name
			conn[isIKEVPNConnectionId] = *connection.ID
			conn[isIKEVPNConnectionHref] = *connection.Href
			connList = append(connList, conn)
		}
	}
	d.Set(isIKEVPNConnections, connList)
	controller, err := flex.GetBaseController(meta)
	if err != nil {
		return err
	}
	d.Set(flex.ResourceControllerURL, controller+"/vpc-ext/network/ikepolicies")
	d.Set(flex.ResourceName, *ike.Name)
	return nil
}

func resourceIBMISIKEPolicyUpdate(d *schema.ResourceData, meta interface{}) error {

	id := d.Id()
	err := ikepUpdate(d, meta, id)
	if err != nil {
		return err
	}
	return resourceIBMISIKEPolicyRead(d, meta)
}

func ikepUpdate(d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	options := &vpcv1.UpdateIkePolicyOptions{
		ID: &id,
	}
	if d.HasChange(isIKEName) || d.HasChange(isIKEAuthenticationAlg) || d.HasChange(isIKEEncryptionAlg) || d.HasChange(isIKEDhGroup) || d.HasChange(isIKEVERSION) || d.HasChange(isIKEKeyLifeTime) {
		name := d.Get(isIKEName).(string)
		authenticationAlg := d.Get(isIKEAuthenticationAlg).(string)
		encryptionAlg := d.Get(isIKEEncryptionAlg).(string)
		keyLifetime := int64(d.Get(isIKEKeyLifeTime).(int))
		dhGroup := int64(d.Get(isIKEDhGroup).(int))
		ikeVersion := int64(d.Get(isIKEVERSION).(int))

		ikePolicyPatchModel := &vpcv1.IkePolicyPatch{}
		ikePolicyPatchModel.Name = &name
		ikePolicyPatchModel.AuthenticationAlgorithm = &authenticationAlg
		ikePolicyPatchModel.EncryptionAlgorithm = &encryptionAlg
		ikePolicyPatchModel.KeyLifetime = &keyLifetime
		ikePolicyPatchModel.DhGroup = &dhGroup
		ikePolicyPatchModel.IkeVersion = &ikeVersion
		ikePolicyPatch, err := ikePolicyPatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("[ERROR] Error calling asPatch for IkePolicyPatch: %s", err)
		}
		options.IkePolicyPatch = ikePolicyPatch

		_, response, err := sess.UpdateIkePolicy(options)
		if err != nil {
			return fmt.Errorf("[ERROR] Error on update of IKE Policy(%s): %s\n%s", id, err, response)
		}
	}
	return nil
}

func resourceIBMISIKEPolicyDelete(d *schema.ResourceData, meta interface{}) error {

	id := d.Id()
	return ikepDelete(d, meta, id)
}

func ikepDelete(d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	getikepoptions := &vpcv1.GetIkePolicyOptions{
		ID: &id,
	}
	_, response, err := sess.GetIkePolicy(getikepoptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error getting IKE Policy(%s): %s\n%s", id, err, response)
	}

	deleteIkePolicyOptions := &vpcv1.DeleteIkePolicyOptions{
		ID: &id,
	}
	response, err = sess.DeleteIkePolicy(deleteIkePolicyOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Deleting IKE Policy(%s): %s\n%s", id, err, response)
	}
	d.SetId("")
	return nil
}

func resourceIBMISIKEPolicyExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	id := d.Id()
	return ikepExists(d, meta, id)
}

func ikepExists(d *schema.ResourceData, meta interface{}, id string) (bool, error) {
	sess, err := vpcClient(meta)
	if err != nil {
		return false, err
	}
	options := &vpcv1.GetIkePolicyOptions{
		ID: &id,
	}
	_, response, err := sess.GetIkePolicy(options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("[ERROR] Error getting IKE Policy(%s): %s\n%s", id, err, response)
	}

	return true, nil
}
