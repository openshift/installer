// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"

	"github.com/IBM/vpc-go-sdk/vpcclassicv1"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

func resourceIBMISIKEPolicy() *schema.Resource {
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
				ValidateFunc: InvokeValidator("ibm_is_ike_policy", isIKEName),
				Description:  "IKE name",
			},

			isIKEAuthenticationAlg: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: InvokeValidator("ibm_is_ike_policy", isIKEAuthenticationAlg),
				Description:  "Authentication algorithm type",
			},

			isIKEEncryptionAlg: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: InvokeValidator("ibm_is_ike_policy", isIKEEncryptionAlg),
				Description:  "Encryption alogorithm type",
			},

			isIKEDhGroup: {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: InvokeValidator("ibm_is_ike_policy", isIKEDhGroup),
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
				ValidateFunc: validateKeyLifeTime,
				Description:  "IKE Key lifetime",
			},

			isIKEVERSION: {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: InvokeValidator("ibm_is_ike_policy", isIKEVERSION),
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

			ResourceGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group name in which resource is provisioned",
			},
		},
	}
}

func resourceIBMISIKEValidator() *ResourceValidator {

	validateSchema := make([]ValidateSchema, 1)
	authentication_algorithm := "md5, sha1, sha256, sha512"
	encryption_algorithm := "triple_des, aes128, aes256"
	dh_group := "2, 5, 14, 19"
	ike_version := "1, 2"
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isIKEName,
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Required:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isIKEAuthenticationAlg,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              authentication_algorithm})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isIKEEncryptionAlg,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              encryption_algorithm})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isIKEDhGroup,
			ValidateFunctionIdentifier: ValidateAllowedIntValue,
			Type:                       TypeInt,
			Required:                   true,
			AllowedValues:              dh_group})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isIKEVERSION,
			ValidateFunctionIdentifier: ValidateAllowedIntValue,
			Type:                       TypeInt,
			Optional:                   true,
			AllowedValues:              ike_version})

	ibmISIKEResourceValidator := ResourceValidator{ResourceName: "ibm_is_ike_policy", Schema: validateSchema}
	return &ibmISIKEResourceValidator
}

func resourceIBMISIKEPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] IKE Policy create")
	name := d.Get(isIKEName).(string)
	authenticationAlg := d.Get(isIKEAuthenticationAlg).(string)
	encryptionAlg := d.Get(isIKEEncryptionAlg).(string)
	dhGroup := int64(d.Get(isIKEDhGroup).(int))

	if userDetails.generation == 1 {
		err := classicIkepCreate(d, meta, authenticationAlg, encryptionAlg, name, dhGroup)
		if err != nil {
			return err
		}
	} else {
		err := ikepCreate(d, meta, authenticationAlg, encryptionAlg, name, dhGroup)
		if err != nil {
			return err
		}
	}
	return resourceIBMISIKEPolicyRead(d, meta)
}

func classicIkepCreate(d *schema.ResourceData, meta interface{}, authenticationAlg, encryptionAlg, name string, dhGroup int64) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	options := &vpcclassicv1.CreateIkePolicyOptions{
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
		options.ResourceGroup = &vpcclassicv1.ResourceGroupIdentity{
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
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	id := d.Id()
	if userDetails.generation == 1 {
		err := classicIkepGet(d, meta, id)
		if err != nil {
			return err
		}
	} else {
		err := ikepGet(d, meta, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func classicIkepGet(d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	getikepoptions := &vpcclassicv1.GetIkePolicyOptions{
		ID: &id,
	}
	ike, response, err := sess.GetIkePolicy(getikepoptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting IKE Policy(%s): %s\n%s", id, err, response)
	}

	d.Set(isIKEName, *ike.Name)
	d.Set(isIKEAuthenticationAlg, *ike.AuthenticationAlgorithm)
	d.Set(isIKEEncryptionAlg, *ike.EncryptionAlgorithm)
	if ike.ResourceGroup != nil {
		d.Set(isIKEResourceGroup, *ike.ResourceGroup.ID)
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
	controller, err := getBaseController(meta)
	if err != nil {
		return err
	}
	d.Set(ResourceControllerURL, controller+"/vpc/network/ikepolicies")
	d.Set(ResourceName, *ike.Name)
	if ike.ResourceGroup != nil {
		rsMangClient, err := meta.(ClientSession).ResourceManagementAPIv2()
		if err != nil {
			return err
		}
		grp, err := rsMangClient.ResourceGroup().Get(*ike.ResourceGroup.ID)
		if err != nil {
			return err
		}
		d.Set(ResourceGroupName, grp.Name)
	}
	return nil
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
		return fmt.Errorf("Error getting IKE Policy(%s): %s\n%s", id, err, response)
	}

	d.Set(isIKEName, *ike.Name)
	d.Set(isIKEAuthenticationAlg, *ike.AuthenticationAlgorithm)
	d.Set(isIKEEncryptionAlg, *ike.EncryptionAlgorithm)
	if ike.ResourceGroup != nil {
		d.Set(isIKEResourceGroup, *ike.ResourceGroup.ID)
		d.Set(ResourceGroupName, *ike.ResourceGroup.Name)
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
	controller, err := getBaseController(meta)
	if err != nil {
		return err
	}
	d.Set(ResourceControllerURL, controller+"/vpc-ext/network/ikepolicies")
	d.Set(ResourceName, *ike.Name)
	return nil
}

func resourceIBMISIKEPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	id := d.Id()

	if userDetails.generation == 1 {
		err := classicIkepUpdate(d, meta, id)
		if err != nil {
			return err
		}
	} else {
		err := ikepUpdate(d, meta, id)
		if err != nil {
			return err
		}
	}
	return resourceIBMISIKEPolicyRead(d, meta)
}

func classicIkepUpdate(d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	options := &vpcclassicv1.UpdateIkePolicyOptions{
		ID: &id,
	}
	if d.HasChange(isIKEName) || d.HasChange(isIKEAuthenticationAlg) || d.HasChange(isIKEEncryptionAlg) || d.HasChange(isIKEDhGroup) || d.HasChange(isIKEVERSION) || d.HasChange(isIKEKeyLifeTime) {
		name := d.Get(isIKEName).(string)
		authenticationAlg := d.Get(isIKEAuthenticationAlg).(string)
		encryptionAlg := d.Get(isIKEEncryptionAlg).(string)
		keyLifetime := int64(d.Get(isIKEKeyLifeTime).(int))
		dhGroup := int64(d.Get(isIKEDhGroup).(int))
		ikeVersion := int64(d.Get(isIKEVERSION).(int))

		ikePolicyPatchModel := &vpcclassicv1.IkePolicyPatch{}
		ikePolicyPatchModel.Name = &name
		ikePolicyPatchModel.AuthenticationAlgorithm = &authenticationAlg
		ikePolicyPatchModel.EncryptionAlgorithm = &encryptionAlg
		ikePolicyPatchModel.KeyLifetime = &keyLifetime
		ikePolicyPatchModel.DhGroup = &dhGroup
		ikePolicyPatchModel.IkeVersion = &ikeVersion
		ikePolicyPatch, err := ikePolicyPatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("Error calling asPatch for ikePolicyPatch: %s", err)
		}
		options.IkePolicyPatch = ikePolicyPatch

		_, response, err := sess.UpdateIkePolicy(options)
		if err != nil {
			return fmt.Errorf("Error on update of IKE Policy(%s): %s\n%s", id, err, response)
		}
	}
	return nil
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
			return fmt.Errorf("Error calling asPatch for IkePolicyPatch: %s", err)
		}
		options.IkePolicyPatch = ikePolicyPatch

		_, response, err := sess.UpdateIkePolicy(options)
		if err != nil {
			return fmt.Errorf("Error on update of IKE Policy(%s): %s\n%s", id, err, response)
		}
	}
	return nil
}

func resourceIBMISIKEPolicyDelete(d *schema.ResourceData, meta interface{}) error {

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	id := d.Id()
	if userDetails.generation == 1 {
		err := classicIkepDelete(d, meta, id)
		if err != nil {
			return err
		}
	} else {
		err := ikepDelete(d, meta, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func classicIkepDelete(d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	getikepoptions := &vpcclassicv1.GetIkePolicyOptions{
		ID: &id,
	}
	_, response, err := sess.GetIkePolicy(getikepoptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error getting IKE Policy(%s): %s\n%s", id, err, response)
	}

	deleteIkePolicyOptions := &vpcclassicv1.DeleteIkePolicyOptions{
		ID: &id,
	}
	response, err = sess.DeleteIkePolicy(deleteIkePolicyOptions)
	if err != nil {
		return fmt.Errorf("Error Deleting IKE Policy(%s): %s\n%s", id, err, response)
	}
	d.SetId("")
	return nil
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
		return fmt.Errorf("Error getting IKE Policy(%s): %s\n%s", id, err, response)
	}

	deleteIkePolicyOptions := &vpcv1.DeleteIkePolicyOptions{
		ID: &id,
	}
	response, err = sess.DeleteIkePolicy(deleteIkePolicyOptions)
	if err != nil {
		return fmt.Errorf("Error Deleting IKE Policy(%s): %s\n%s", id, err, response)
	}
	d.SetId("")
	return nil
}

func resourceIBMISIKEPolicyExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return false, err
	}
	id := d.Id()

	if userDetails.generation == 1 {
		exists, err := classicikepExists(d, meta, id)
		return exists, err
	} else {
		exists, err := ikepExists(d, meta, id)
		return exists, err
	}
}

func classicikepExists(d *schema.ResourceData, meta interface{}, id string) (bool, error) {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return false, err
	}
	options := &vpcclassicv1.GetIkePolicyOptions{
		ID: &id,
	}
	_, response, err := sess.GetIkePolicy(options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("Error getting IKE Policy(%s): %s\n%s", id, err, response)
	}

	return true, nil
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
		return false, fmt.Errorf("Error getting IKE Policy(%s): %s\n%s", id, err, response)
	}

	return true, nil
}
