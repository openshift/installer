// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isIpSecName              = "name"
	isIpSecAuthenticationAlg = "authentication_algorithm"
	isIpSecEncryptionAlg     = "encryption_algorithm"
	isIpSecPFS               = "pfs"
	isIpSecKeyLifeTime       = "key_lifetime"
	isIPSecResourceGroup     = "resource_group"
	isIPSecEncapsulationMode = "encapsulation_mode"
	isIPSecTransformProtocol = "transform_protocol"
	isIPSecVPNConnections    = "vpn_connections"
	isIPSecVPNConnectionName = "name"
	isIPSecVPNConnectionId   = "id"
	isIPSecVPNConnectionHref = "href"
)

func ResourceIBMISIPSecPolicy() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMISIPSecPolicyCreate,
		Read:     resourceIBMISIPSecPolicyRead,
		Update:   resourceIBMISIPSecPolicyUpdate,
		Delete:   resourceIBMISIPSecPolicyDelete,
		Exists:   resourceIBMISIPSecPolicyExists,
		Importer: &schema.ResourceImporter{},

		CustomizeDiff: customdiff.All(
			customdiff.Sequence(
				func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
					return flex.ResourceIPSecPolicyValidate(diff)
				}),
		),

		Schema: map[string]*schema.Schema{
			isIpSecName: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_ipsec_policy", isIpSecName),
				Description:  "IPSEC name",
			},

			isIpSecAuthenticationAlg: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_ipsec_policy", isIpSecAuthenticationAlg),
				Description:  "Authentication alorothm",
			},

			isIpSecEncryptionAlg: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_ipsec_policy", isIpSecEncryptionAlg),
				Description:  "Encryption algorithm",
			},

			isIpSecPFS: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_ipsec_policy", isIpSecPFS),
				Description:  "PFS info",
			},

			isIPSecResourceGroup: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: "Resource group info",
			},

			isIpSecKeyLifeTime: {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3600,
				ValidateFunc: validate.ValidateKeyLifeTime,
				Description:  "IPSEC key lifetime",
			},

			isIPSecEncapsulationMode: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IPSEC encapsulation mode",
			},

			isIPSecTransformProtocol: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IPSEC transform protocol",
			},

			isIPSecVPNConnections: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isIPSecVPNConnectionName: {
							Type:     schema.TypeString,
							Computed: true,
						},
						isIPSecVPNConnectionId: {
							Type:     schema.TypeString,
							Computed: true,
						},
						isIPSecVPNConnectionHref: {
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

			flex.ResourceCRN: {
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

func ResourceIBMISIPSECValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)
	authentication_algorithm := "md5, sha1, sha256, sha512, sha384, disabled"
	encryption_algorithm := "triple_des, aes128, aes256, aes128gcm16, aes192gcm16, aes256gcm16"
	pfs := "disabled, group_2, group_5, group_14, group_19, group_15, group_16, group_17, group_18, group_20, group_21, group_22, group_23, group_24, group_31"
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isIpSecName,
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isIpSecAuthenticationAlg,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              authentication_algorithm})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isIpSecEncryptionAlg,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              encryption_algorithm})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isIpSecPFS,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              pfs})

	ibmISIPSECResourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_ipsec_policy", Schema: validateSchema}
	return &ibmISIPSECResourceValidator
}

func resourceIBMISIPSecPolicyCreate(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] Ip Sec create")
	name := d.Get(isIpSecName).(string)
	authenticationAlg := d.Get(isIpSecAuthenticationAlg).(string)
	encryptionAlg := d.Get(isIpSecEncryptionAlg).(string)
	pfs := d.Get(isIpSecPFS).(string)

	err := ipsecpCreate(d, meta, authenticationAlg, encryptionAlg, name, pfs)
	if err != nil {
		return err
	}
	return resourceIBMISIPSecPolicyRead(d, meta)
}

func ipsecpCreate(d *schema.ResourceData, meta interface{}, authenticationAlg, encryptionAlg, name, pfs string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	options := &vpcv1.CreateIpsecPolicyOptions{
		AuthenticationAlgorithm: &authenticationAlg,
		EncryptionAlgorithm:     &encryptionAlg,
		Pfs:                     &pfs,
		Name:                    &name,
	}

	if keylt, ok := d.GetOk(isIpSecKeyLifeTime); ok {
		keyLifetime := int64(keylt.(int))
		options.KeyLifetime = &keyLifetime
	} else {
		keyLifetime := int64(3600)
		options.KeyLifetime = &keyLifetime
	}

	if rgrp, ok := d.GetOk(isIPSecResourceGroup); ok {
		rg := rgrp.(string)
		options.ResourceGroup = &vpcv1.ResourceGroupIdentity{
			ID: &rg,
		}
	}
	ipSec, response, err := sess.CreateIpsecPolicy(options)
	if err != nil {
		return fmt.Errorf("[DEBUG] ipSec policy err %s\n%s", err, response)
	}
	d.SetId(*ipSec.ID)
	log.Printf("[INFO] ipSec policy : %s", *ipSec.ID)
	return nil
}

func resourceIBMISIPSecPolicyRead(d *schema.ResourceData, meta interface{}) error {

	id := d.Id()
	return ipsecpGet(d, meta, id)
}

func ipsecpGet(d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	getIpsecPolicyOptions := &vpcv1.GetIpsecPolicyOptions{
		ID: &id,
	}
	ipSec, response, err := sess.GetIpsecPolicy(getIpsecPolicyOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error getting IPSEC Policy(%s): %s\n%s", id, err, response)
	}
	d.Set(isIpSecName, *ipSec.Name)
	d.Set(isIpSecAuthenticationAlg, *ipSec.AuthenticationAlgorithm)
	d.Set(isIpSecEncryptionAlg, *ipSec.EncryptionAlgorithm)
	if ipSec.ResourceGroup != nil {
		d.Set(isIPSecResourceGroup, *ipSec.ResourceGroup.ID)
		d.Set(flex.ResourceGroupName, *ipSec.ResourceGroup.Name)
	} else {
		d.Set(isIPSecResourceGroup, nil)
	}
	d.Set(isIpSecPFS, *ipSec.Pfs)
	if ipSec.KeyLifetime != nil {
		d.Set(isIpSecKeyLifeTime, *ipSec.KeyLifetime)
	}
	d.Set(isIPSecEncapsulationMode, *ipSec.EncapsulationMode)
	d.Set(isIPSecTransformProtocol, *ipSec.TransformProtocol)

	connList := make([]map[string]interface{}, 0)
	if ipSec.Connections != nil && len(ipSec.Connections) > 0 {
		for _, connection := range ipSec.Connections {
			conn := map[string]interface{}{}
			conn[isIPSecVPNConnectionName] = *connection.Name
			conn[isIPSecVPNConnectionId] = *connection.ID
			conn[isIPSecVPNConnectionHref] = *connection.Href
			connList = append(connList, conn)
		}
	}
	d.Set(isIPSecVPNConnections, connList)
	controller, err := flex.GetBaseController(meta)
	if err != nil {
		return err
	}
	d.Set(flex.ResourceControllerURL, controller+"/vpc-ext/network/ipsecpolicies")
	d.Set(flex.ResourceName, *ipSec.Name)
	// d.Set(flex.ResourceCRN, *ipSec.Crn)
	return nil
}

func resourceIBMISIPSecPolicyUpdate(d *schema.ResourceData, meta interface{}) error {

	id := d.Id()
	err := ipsecpUpdate(d, meta, id)
	if err != nil {
		return err
	}

	return resourceIBMISIPSecPolicyRead(d, meta)
}

func ipsecpUpdate(d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	options := &vpcv1.UpdateIpsecPolicyOptions{
		ID: &id,
	}
	if d.HasChange(isIpSecName) || d.HasChange(isIpSecAuthenticationAlg) || d.HasChange(isIpSecEncryptionAlg) || d.HasChange(isIpSecPFS) || d.HasChange(isIpSecKeyLifeTime) {
		name := d.Get(isIpSecName).(string)
		authenticationAlg := d.Get(isIpSecAuthenticationAlg).(string)
		encryptionAlg := d.Get(isIpSecEncryptionAlg).(string)
		pfs := d.Get(isIpSecPFS).(string)
		keyLifetime := int64(d.Get(isIpSecKeyLifeTime).(int))

		ipsecPolicyPatchModel := &vpcv1.IPsecPolicyPatch{
			Name:                    &name,
			AuthenticationAlgorithm: &authenticationAlg,
			EncryptionAlgorithm:     &encryptionAlg,
			Pfs:                     &pfs,
			KeyLifetime:             &keyLifetime,
		}
		ipsecPolicyPatch, err := ipsecPolicyPatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("[ERROR] Error calling asPatch for IPsecPolicyPatch: %s", err)
		}
		options.IPsecPolicyPatch = ipsecPolicyPatch

		_, response, err := sess.UpdateIpsecPolicy(options)
		if err != nil {
			return fmt.Errorf("[ERROR] Error on update of IPSEC Policy(%s): %s\n%s", id, err, response)
		}
	}
	return nil
}

func resourceIBMISIPSecPolicyDelete(d *schema.ResourceData, meta interface{}) error {

	id := d.Id()
	return ipsecpDelete(d, meta, id)
}

func ipsecpDelete(d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	getIpsecPolicyOptions := &vpcv1.GetIpsecPolicyOptions{
		ID: &id,
	}
	_, response, err := sess.GetIpsecPolicy(getIpsecPolicyOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error getting IPSEC Policy(%s): %s\n%s", id, err, response)
	}
	deleteIpsecPolicyOptions := &vpcv1.DeleteIpsecPolicyOptions{
		ID: &id,
	}
	response, err = sess.DeleteIpsecPolicy(deleteIpsecPolicyOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Deleting IPSEC Policy(%s): %s\n%s", id, err, response)
	}
	d.SetId("")
	return nil
}

func resourceIBMISIPSecPolicyExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	id := d.Id()
	return ipsecpExists(d, meta, id)
}

func ipsecpExists(d *schema.ResourceData, meta interface{}, id string) (bool, error) {
	sess, err := vpcClient(meta)
	if err != nil {
		return false, err
	}
	getIpsecPolicyOptions := &vpcv1.GetIpsecPolicyOptions{
		ID: &id,
	}
	_, response, err := sess.GetIpsecPolicy(getIpsecPolicyOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("[ERROR] Error getting IPSEC Policy(%s): %s\n%s", id, err, response)
	}
	return true, nil
}
