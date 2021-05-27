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
	isPublicGatewayName              = "name"
	isPublicGatewayFloatingIP        = "floating_ip"
	isPublicGatewayStatus            = "status"
	isPublicGatewayVPC               = "vpc"
	isPublicGatewayZone              = "zone"
	isPublicGatewayFloatingIPAddress = "address"
	isPublicGatewayTags              = "tags"

	isPublicGatewayProvisioning     = "provisioning"
	isPublicGatewayProvisioningDone = "available"
	isPublicGatewayDeleting         = "deleting"
	isPublicGatewayDeleted          = "done"

	isPublicGatewayResourceGroup = "resource_group"
)

func resourceIBMISPublicGateway() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMISPublicGatewayCreate,
		Read:     resourceIBMISPublicGatewayRead,
		Update:   resourceIBMISPublicGatewayUpdate,
		Delete:   resourceIBMISPublicGatewayDelete,
		Exists:   resourceIBMISPublicGatewayExists,
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
			isPublicGatewayName: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: InvokeValidator("ibm_is_public_gateway", isPublicGatewayName),
				Description:  "Name of the Public gateway instance",
			},

			isPublicGatewayFloatingIP: {
				Type:             schema.TypeMap,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: applyOnce,
			},

			isPublicGatewayStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public gateway instance status",
			},

			isPublicGatewayResourceGroup: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: "Public gateway resource group info",
			},

			isPublicGatewayVPC: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Public gateway VPC info",
			},

			isPublicGatewayZone: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Public gateway zone info",
			},

			isPublicGatewayTags: {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: InvokeValidator("ibm_is_public_gateway", "tag")},
				Set:         resourceIBMVPCHash,
				Description: "Service tags for the public gateway instance",
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
		},
	}
}

func resourceIBMISPublicGatewayValidator() *ResourceValidator {

	validateSchema := make([]ValidateSchema, 1)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isPublicGatewayName,
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Required:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "tag",
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9:_ .-]+$`,
			MinValueLength:             1,
			MaxValueLength:             128})

	ibmISPublicGatewayResourceValidator := ResourceValidator{ResourceName: "ibm_is_public_gateway", Schema: validateSchema}
	return &ibmISPublicGatewayResourceValidator
}

func resourceIBMISPublicGatewayCreate(d *schema.ResourceData, meta interface{}) error {

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	name := d.Get(isPublicGatewayName).(string)
	vpc := d.Get(isPublicGatewayVPC).(string)
	zone := d.Get(isPublicGatewayZone).(string)

	if userDetails.generation == 1 {
		err := classicPgwCreate(d, meta, name, vpc, zone)
		if err != nil {
			return err
		}
	} else {
		err := pgwCreate(d, meta, name, vpc, zone)
		if err != nil {
			return err
		}
	}
	return resourceIBMISPublicGatewayRead(d, meta)
}

func classicPgwCreate(d *schema.ResourceData, meta interface{}, name, vpc, zone string) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}

	options := &vpcclassicv1.CreatePublicGatewayOptions{
		Name: &name,
		VPC: &vpcclassicv1.VPCIdentity{
			ID: &vpc,
		},
		Zone: &vpcclassicv1.ZoneIdentity{
			Name: &zone,
		},
	}
	floatingipID := ""
	floatingipadd := ""
	if floatingipdataIntf, ok := d.GetOk(isPublicGatewayFloatingIP); ok && floatingipdataIntf != nil {
		fip := &vpcclassicv1.PublicGatewayFloatingIPPrototype{}
		floatingipdata := floatingipdataIntf.(map[string]interface{})
		if floatingipidintf, ok := floatingipdata["id"]; ok && floatingipidintf != nil {
			floatingipID = floatingipidintf.(string)
			fip.ID = &floatingipID
		}
		if floatingipaddintf, ok := floatingipdata[isPublicGatewayFloatingIPAddress]; ok && floatingipaddintf != nil {
			floatingipadd = floatingipaddintf.(string)
			fip.Address = &floatingipadd
		}
		options.FloatingIP = fip
	}

	publicgw, response, err := sess.CreatePublicGateway(options)
	if err != nil {
		return fmt.Errorf("Error while creating Public Gateway %s\n%s", err, response)
	}
	d.SetId(*publicgw.ID)
	log.Printf("[INFO] PublicGateway : %s", *publicgw.ID)

	_, err = isWaitForClassicPublicGatewayAvailable(sess, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk(isPublicGatewayTags); ok || v != "" {
		oldList, newList := d.GetChange(isPublicGatewayTags)
		err = UpdateTagsUsingCRN(oldList, newList, meta, *publicgw.CRN)
		if err != nil {
			log.Printf(
				"Error on create of vpc public gateway (%s) tags: %s", d.Id(), err)
		}
	}
	return nil
}

func pgwCreate(d *schema.ResourceData, meta interface{}, name, vpc, zone string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	options := &vpcv1.CreatePublicGatewayOptions{
		Name: &name,
		VPC: &vpcv1.VPCIdentity{
			ID: &vpc,
		},
		Zone: &vpcv1.ZoneIdentity{
			Name: &zone,
		},
	}
	floatingipID := ""
	floatingipadd := ""
	if floatingipdataIntf, ok := d.GetOk(isPublicGatewayFloatingIP); ok && floatingipdataIntf != nil {
		fip := &vpcv1.PublicGatewayFloatingIPPrototype{}
		floatingipdata := floatingipdataIntf.(map[string]interface{})
		if floatingipidintf, ok := floatingipdata["id"]; ok && floatingipidintf != nil {
			floatingipID = floatingipidintf.(string)
			fip.ID = &floatingipID
		}
		if floatingipaddintf, ok := floatingipdata[isPublicGatewayFloatingIPAddress]; ok && floatingipaddintf != nil {
			floatingipadd = floatingipaddintf.(string)
			fip.Address = &floatingipadd
		}
		options.FloatingIP = fip
	}
	if grp, ok := d.GetOk(isPublicGatewayResourceGroup); ok {
		rg := grp.(string)
		options.ResourceGroup = &vpcv1.ResourceGroupIdentity{
			ID: &rg,
		}
	}

	publicgw, response, err := sess.CreatePublicGateway(options)
	if err != nil {
		return fmt.Errorf("Error while creating Public Gateway %s\n%s", err, response)
	}
	d.SetId(*publicgw.ID)
	log.Printf("[INFO] PublicGateway : %s", *publicgw.ID)

	_, err = isWaitForPublicGatewayAvailable(sess, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk(isPublicGatewayTags); ok || v != "" {
		oldList, newList := d.GetChange(isPublicGatewayTags)
		err = UpdateTagsUsingCRN(oldList, newList, meta, *publicgw.CRN)
		if err != nil {
			log.Printf(
				"Error on create of vpc public gateway (%s) tags: %s", d.Id(), err)
		}
	}
	return nil
}

func isWaitForClassicPublicGatewayAvailable(publicgwC *vpcclassicv1.VpcClassicV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for public gateway (%s) to be available.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isPublicGatewayProvisioning},
		Target:     []string{isPublicGatewayProvisioningDone, ""},
		Refresh:    isClassicPublicGatewayRefreshFunc(publicgwC, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isClassicPublicGatewayRefreshFunc(publicgwC *vpcclassicv1.VpcClassicV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getPublicGatewayOptions := &vpcclassicv1.GetPublicGatewayOptions{
			ID: &id,
		}
		publicgw, response, err := publicgwC.GetPublicGateway(getPublicGatewayOptions)
		if err != nil {
			return nil, "", fmt.Errorf("Error getting Public Gateway : %s\n%s", err, response)
		}

		if *publicgw.Status == isPublicGatewayProvisioningDone {
			return publicgw, isPublicGatewayProvisioningDone, nil
		}

		return publicgw, isPublicGatewayProvisioning, nil
	}
}

func isWaitForPublicGatewayAvailable(publicgwC *vpcv1.VpcV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for public gateway (%s) to be available.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isPublicGatewayProvisioning},
		Target:     []string{isPublicGatewayProvisioningDone, ""},
		Refresh:    isPublicGatewayRefreshFunc(publicgwC, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isPublicGatewayRefreshFunc(publicgwC *vpcv1.VpcV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getPublicGatewayOptions := &vpcv1.GetPublicGatewayOptions{
			ID: &id,
		}
		publicgw, response, err := publicgwC.GetPublicGateway(getPublicGatewayOptions)
		if err != nil {
			return nil, "", fmt.Errorf("Error getting Public Gateway : %s\n%s", err, response)
		}

		if *publicgw.Status == isPublicGatewayProvisioningDone {
			return publicgw, isPublicGatewayProvisioningDone, nil
		}

		return publicgw, isPublicGatewayProvisioning, nil
	}
}

func resourceIBMISPublicGatewayRead(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	id := d.Id()
	if userDetails.generation == 1 {
		err := classicPgwGet(d, meta, id)
		if err != nil {
			return err
		}
	} else {
		err := pgwGet(d, meta, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func classicPgwGet(d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	getPublicGatewayOptions := &vpcclassicv1.GetPublicGatewayOptions{
		ID: &id,
	}
	publicgw, response, err := sess.GetPublicGateway(getPublicGatewayOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting Public Gateway : %s\n%s", err, response)
	}
	d.Set(isPublicGatewayName, *publicgw.Name)
	if publicgw.FloatingIP != nil {
		floatIP := map[string]interface{}{
			"id":                             *publicgw.FloatingIP.ID,
			isPublicGatewayFloatingIPAddress: *publicgw.FloatingIP.Address,
		}
		d.Set(isPublicGatewayFloatingIP, floatIP)

	}
	d.Set(isPublicGatewayStatus, *publicgw.Status)
	d.Set(isPublicGatewayZone, *publicgw.Zone.Name)
	d.Set(isPublicGatewayVPC, *publicgw.VPC.ID)
	tags, err := GetTagsUsingCRN(meta, *publicgw.CRN)
	if err != nil {
		log.Printf(
			"Error on get of vpc public gateway (%s) tags: %s", id, err)
	}
	d.Set(isPublicGatewayTags, tags)
	controller, err := getBaseController(meta)
	if err != nil {
		return err
	}
	d.Set(ResourceControllerURL, controller+"/vpc/network/publicGateways")
	d.Set(ResourceName, *publicgw.Name)
	d.Set(ResourceCRN, *publicgw.CRN)
	d.Set(ResourceStatus, *publicgw.Status)
	return nil
}

func pgwGet(d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	getPublicGatewayOptions := &vpcv1.GetPublicGatewayOptions{
		ID: &id,
	}
	publicgw, response, err := sess.GetPublicGateway(getPublicGatewayOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting Public Gateway : %s\n%s", err, response)
	}
	d.Set(isPublicGatewayName, *publicgw.Name)
	if publicgw.FloatingIP != nil {
		floatIP := map[string]interface{}{
			"id":                             *publicgw.FloatingIP.ID,
			isPublicGatewayFloatingIPAddress: *publicgw.FloatingIP.Address,
		}
		d.Set(isPublicGatewayFloatingIP, floatIP)

	}
	d.Set(isPublicGatewayStatus, *publicgw.Status)
	d.Set(isPublicGatewayZone, *publicgw.Zone.Name)
	d.Set(isPublicGatewayVPC, *publicgw.VPC.ID)
	tags, err := GetTagsUsingCRN(meta, *publicgw.CRN)
	if err != nil {
		log.Printf(
			"Error on get of vpc public gateway (%s) tags: %s", id, err)
	}
	d.Set(isPublicGatewayTags, tags)
	controller, err := getBaseController(meta)
	if err != nil {
		return err
	}
	d.Set(ResourceControllerURL, controller+"/vpc-ext/network/publicGateways")
	d.Set(ResourceName, *publicgw.Name)
	d.Set(ResourceCRN, *publicgw.CRN)
	d.Set(ResourceStatus, *publicgw.Status)
	if publicgw.ResourceGroup != nil {
		d.Set(isPublicGatewayResourceGroup, *publicgw.ResourceGroup.ID)
		d.Set(ResourceGroupName, *publicgw.ResourceGroup.Name)
	}
	return nil
}

func resourceIBMISPublicGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	id := d.Id()

	name := ""
	hasChanged := false

	if d.HasChange(isPublicGatewayName) {
		name = d.Get(isPublicGatewayName).(string)
		hasChanged = true
	}
	if userDetails.generation == 1 {
		err := classicPgwUpdate(d, meta, id, name, hasChanged)
		if err != nil {
			return err
		}
	} else {
		err := pgwUpdate(d, meta, id, name, hasChanged)
		if err != nil {
			return err
		}
	}
	return resourceIBMISPublicGatewayRead(d, meta)
}

func classicPgwUpdate(d *schema.ResourceData, meta interface{}, id, name string, hasChanged bool) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	if d.HasChange(isPublicGatewayTags) {
		getPublicGatewayOptions := &vpcclassicv1.GetPublicGatewayOptions{
			ID: &id,
		}
		publicgw, response, err := sess.GetPublicGateway(getPublicGatewayOptions)
		if err != nil {
			return fmt.Errorf("Error getting Public Gateway : %s\n%s", err, response)
		}
		oldList, newList := d.GetChange(isPublicGatewayTags)
		err = UpdateTagsUsingCRN(oldList, newList, meta, *publicgw.CRN)
		if err != nil {
			log.Printf(
				"Error on update of resource Public Gateway (%s) tags: %s", id, err)
		}
	}
	if hasChanged {
		updatePublicGatewayOptions := &vpcclassicv1.UpdatePublicGatewayOptions{
			ID: &id,
		}

		PublicGatewayPatchModel := &vpcclassicv1.PublicGatewayPatch{
			Name: &name,
		}
		PublicGatewayPatch, err := PublicGatewayPatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("Error calling asPatch for PublicGatewayPatch: %s", err)
		}
		updatePublicGatewayOptions.PublicGatewayPatch = PublicGatewayPatch

		_, response, err := sess.UpdatePublicGateway(updatePublicGatewayOptions)
		if err != nil {
			return fmt.Errorf("Error Updating Public Gateway  : %s\n%s", err, response)
		}
	}
	return nil
}

func pgwUpdate(d *schema.ResourceData, meta interface{}, id, name string, hasChanged bool) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	if d.HasChange(isPublicGatewayTags) {
		getPublicGatewayOptions := &vpcv1.GetPublicGatewayOptions{
			ID: &id,
		}
		publicgw, response, err := sess.GetPublicGateway(getPublicGatewayOptions)
		if err != nil {
			return fmt.Errorf("Error getting Public Gateway : %s\n%s", err, response)
		}
		oldList, newList := d.GetChange(isPublicGatewayTags)
		err = UpdateTagsUsingCRN(oldList, newList, meta, *publicgw.CRN)
		if err != nil {
			log.Printf(
				"Error on update of resource Public Gateway (%s) tags: %s", id, err)
		}
	}
	if hasChanged {
		updatePublicGatewayOptions := &vpcv1.UpdatePublicGatewayOptions{
			ID: &id,
		}
		PublicGatewayPatchModel := &vpcv1.PublicGatewayPatch{
			Name: &name,
		}
		PublicGatewayPatch, err := PublicGatewayPatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("Error calling asPatch for PublicGatewayPatch: %s", err)
		}
		updatePublicGatewayOptions.PublicGatewayPatch = PublicGatewayPatch
		_, response, err := sess.UpdatePublicGateway(updatePublicGatewayOptions)
		if err != nil {
			return fmt.Errorf("Error Updating Public Gateway  : %s\n%s", err, response)
		}
	}
	return nil
}

func resourceIBMISPublicGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	id := d.Id()
	if userDetails.generation == 1 {
		err := classicPgwDelete(d, meta, id)
		if err != nil {
			return err
		}
	} else {
		err := pgwDelete(d, meta, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func classicPgwDelete(d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}

	getPublicGatewayOptions := &vpcclassicv1.GetPublicGatewayOptions{
		ID: &id,
	}
	_, response, err := sess.GetPublicGateway(getPublicGatewayOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Getting Public Gateway (%s): %s\n%s", id, err, response)
	}

	deletePublicGatewayOptions := &vpcclassicv1.DeletePublicGatewayOptions{
		ID: &id,
	}
	response, err = sess.DeletePublicGateway(deletePublicGatewayOptions)
	if err != nil {
		return fmt.Errorf("Error Deleting Public Gateway : %s\n%s", err, response)
	}
	_, err = isWaitForClassicPublicGatewayDeleted(sess, id, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func pgwDelete(d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	getPublicGatewayOptions := &vpcv1.GetPublicGatewayOptions{
		ID: &id,
	}
	_, response, err := sess.GetPublicGateway(getPublicGatewayOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Getting Public Gateway (%s): %s\n%s", id, err, response)
	}

	deletePublicGatewayOptions := &vpcv1.DeletePublicGatewayOptions{
		ID: &id,
	}
	response, err = sess.DeletePublicGateway(deletePublicGatewayOptions)
	if err != nil {
		return fmt.Errorf("Error Deleting Public Gateway : %s\n%s", err, response)
	}
	_, err = isWaitForPublicGatewayDeleted(sess, id, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func isWaitForClassicPublicGatewayDeleted(pg *vpcclassicv1.VpcClassicV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for public gateway (%s) to be deleted.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isPublicGatewayDeleting},
		Target:     []string{isPublicGatewayDeleted, ""},
		Refresh:    isClassicPublicGatewayDeleteRefreshFunc(pg, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isClassicPublicGatewayDeleteRefreshFunc(pg *vpcclassicv1.VpcClassicV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] delete function here")
		getPublicGatewayOptions := &vpcclassicv1.GetPublicGatewayOptions{
			ID: &id,
		}
		pgw, response, err := pg.GetPublicGateway(getPublicGatewayOptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return pgw, isPublicGatewayDeleted, nil
			}
			return nil, "", fmt.Errorf("The Public Gateway %s failed to delete: %s\n%s", id, err, response)
		}
		return pgw, isPublicGatewayDeleting, nil
	}
}

func isWaitForPublicGatewayDeleted(pg *vpcv1.VpcV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for public gateway (%s) to be deleted.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isPublicGatewayDeleting},
		Target:     []string{isPublicGatewayDeleted, ""},
		Refresh:    isPublicGatewayDeleteRefreshFunc(pg, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isPublicGatewayDeleteRefreshFunc(pg *vpcv1.VpcV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] delete function here")
		getPublicGatewayOptions := &vpcv1.GetPublicGatewayOptions{
			ID: &id,
		}
		pgw, response, err := pg.GetPublicGateway(getPublicGatewayOptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return pgw, isPublicGatewayDeleted, nil
			}
			return nil, "", fmt.Errorf("The Public Gateway %s failed to delete: %s\n%s", id, err, response)
		}
		return pgw, isPublicGatewayDeleting, nil
	}
}

func resourceIBMISPublicGatewayExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return false, err
	}
	id := d.Id()
	if userDetails.generation == 1 {
		exists, err := classicPgwExists(d, meta, id)
		return exists, err
	} else {
		exists, err := pgwExists(d, meta, id)
		return exists, err
	}
}

func classicPgwExists(d *schema.ResourceData, meta interface{}, id string) (bool, error) {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return false, err
	}
	getPublicGatewayOptions := &vpcclassicv1.GetPublicGatewayOptions{
		ID: &id,
	}
	_, response, err := sess.GetPublicGateway(getPublicGatewayOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("Error getting Public Gateway: %s\n%s", err, response)
	}
	return true, nil
}

func pgwExists(d *schema.ResourceData, meta interface{}, id string) (bool, error) {
	sess, err := vpcClient(meta)
	if err != nil {
		return false, err
	}
	getPublicGatewayOptions := &vpcv1.GetPublicGatewayOptions{
		ID: &id,
	}
	_, response, err := sess.GetPublicGateway(getPublicGatewayOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("Error getting Public Gateway: %s\n%s", err, response)
	}
	return true, nil
}
