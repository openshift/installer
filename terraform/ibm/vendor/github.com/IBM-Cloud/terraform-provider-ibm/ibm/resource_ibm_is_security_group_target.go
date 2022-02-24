// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isSecurityGroupTargetID     = "target"
	isSecurityGroupResourceType = "resource_type"
)

func resourceIBMISSecurityGroupTarget() *schema.Resource {

	return &schema.Resource{
		Create:   resourceIBMISSecurityGroupTargetCreate,
		Read:     resourceIBMISSecurityGroupTargetRead,
		Delete:   resourceIBMISSecurityGroupTargetDelete,
		Exists:   resourceIBMISSecurityGroupTargetExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			"security_group": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Security group id",
			},

			isSecurityGroupTargetID: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "security group target identifier",
				ValidateFunc: InvokeValidator("ibm_is_security_group_target", isSecurityGroupTargetID),
			},

			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Security group target name",
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this Security group target",
			},

			isSecurityGroupResourceType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource Type",
			},
		},
	}
}

func resourceIBMISSecurityGroupTargetValidator() *ResourceValidator {

	validateSchema := make([]ValidateSchema, 0)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isSecurityGroupTargetID,
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Required:                   true,
			Regexp:                     `^[-0-9a-z_]+$`,
			MinValueLength:             1,
			MaxValueLength:             64},
		ValidateSchema{
			Identifier:                 "security_group",
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Required:                   true,
			Regexp:                     `^[-0-9a-z_]+$`,
			MinValueLength:             1,
			MaxValueLength:             64})

	ibmISSecurityGroupResourceValidator := ResourceValidator{ResourceName: "ibm_is_security_group_target", Schema: validateSchema}
	return &ibmISSecurityGroupResourceValidator
}

func resourceIBMISSecurityGroupTargetCreate(d *schema.ResourceData, meta interface{}) error {

	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	securityGroupID := d.Get("security_group").(string)
	targetID := d.Get(isSecurityGroupTargetID).(string)

	createSecurityGroupTargetBindingOptions := &vpcv1.CreateSecurityGroupTargetBindingOptions{}
	createSecurityGroupTargetBindingOptions.SecurityGroupID = &securityGroupID
	createSecurityGroupTargetBindingOptions.ID = &targetID

	sg, response, err := sess.CreateSecurityGroupTargetBinding(createSecurityGroupTargetBindingOptions)
	if err != nil || sg == nil {
		return fmt.Errorf("error while creating Security Group Target Binding %s\n%s", err, response)
	}
	sgtarget := sg.(*vpcv1.SecurityGroupTargetReference)
	d.SetId(fmt.Sprintf("%s/%s", securityGroupID, *sgtarget.ID))
	crn := sgtarget.CRN
	if crn != nil && *crn != "" && strings.Contains(*crn, "load-balancer") {
		lbid := sgtarget.ID
		_, errsgt := isWaitForLbSgTargetCreateAvailable(sess, *lbid, d.Timeout(schema.TimeoutCreate))
		if errsgt != nil {
			return errsgt
		}
	}

	return resourceIBMISSecurityGroupTargetRead(d, meta)
}

func resourceIBMISSecurityGroupTargetRead(d *schema.ResourceData, meta interface{}) error {

	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	securityGroupID := parts[0]
	securityGroupTargetID := parts[1]

	getSecurityGroupTargetOptions := &vpcv1.GetSecurityGroupTargetOptions{
		SecurityGroupID: &securityGroupID,
		ID:              &securityGroupTargetID,
	}

	data, response, err := sess.GetSecurityGroupTarget(getSecurityGroupTargetOptions)
	if err != nil || data == nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error getting Security Group Target : %s\n%s", err, response)
	}

	target := data.(*vpcv1.SecurityGroupTargetReference)
	d.Set("name", *target.Name)
	d.Set("crn", target.CRN)
	if target.ResourceType != nil && *target.ResourceType != "" {
		d.Set(isSecurityGroupResourceType, *target.ResourceType)
	}

	return nil
}

func resourceIBMISSecurityGroupTargetDelete(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	securityGroupID := parts[0]
	securityGroupTargetID := parts[1]

	getSecurityGroupTargetOptions := &vpcv1.GetSecurityGroupTargetOptions{
		SecurityGroupID: &securityGroupID,
		ID:              &securityGroupTargetID,
	}
	sgt, response, err := sess.GetSecurityGroupTarget(getSecurityGroupTargetOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error Getting Security Group Targets (%s): %s\n%s", securityGroupID, err, response)
	}

	deleteSecurityGroupTargetBindingOptions := sess.NewDeleteSecurityGroupTargetBindingOptions(securityGroupID, securityGroupTargetID)
	response, err = sess.DeleteSecurityGroupTargetBinding(deleteSecurityGroupTargetBindingOptions)
	if err != nil {
		return fmt.Errorf("error Deleting Security Group Targets : %s\n%s", err, response)
	}
	securityGroupTargetReference := sgt.(*vpcv1.SecurityGroupTargetReference)
	crn := securityGroupTargetReference.CRN
	if crn != nil && *crn != "" && strings.Contains(*crn, "load-balancer") {
		lbid := securityGroupTargetReference.ID
		_, errsgt := isWaitForLBRemoveAvailable(sess, sgt, *lbid, securityGroupID, securityGroupTargetID, d.Timeout(schema.TimeoutDelete))
		if errsgt != nil {
			return errsgt
		}
	}
	d.SetId("")
	return nil
}

func resourceIBMISSecurityGroupTargetExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	sess, err := vpcClient(meta)
	if err != nil {
		return false, err
	}

	parts, err := idParts(d.Id())
	if err != nil {
		return false, err
	}
	securityGroupID := parts[0]
	securityGroupTargetID := parts[1]

	getSecurityGroupTargetOptions := &vpcv1.GetSecurityGroupTargetOptions{
		SecurityGroupID: &securityGroupID,
		ID:              &securityGroupTargetID,
	}

	_, response, err := sess.GetSecurityGroupTarget(getSecurityGroupTargetOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return false, nil
		}
		return false, fmt.Errorf("error getting Security Group Target : %s\n%s", err, response)
	}
	return true, nil

}

func isWaitForLBRemoveAvailable(sess *vpcv1.VpcV1, sgt vpcv1.SecurityGroupTargetReferenceIntf, lbId, securityGroupID, securityGroupTargetID string, timeout time.Duration) (interface{}, error) {
	log.Printf("[INFO] Waiting for load balancer binding (%s) to be removed.", lbId)

	stateConf := &resource.StateChangeConf{
		Pending:        []string{isLBProvisioning},
		Target:         []string{isLBProvisioningDone},
		Refresh:        isLBRemoveRefreshFunc(sess, sgt, lbId, securityGroupID, securityGroupTargetID),
		Timeout:        timeout,
		Delay:          10 * time.Second,
		MinTimeout:     10 * time.Second,
		NotFoundChecks: 1,
	}

	return stateConf.WaitForState()
}

func isLBRemoveRefreshFunc(sess *vpcv1.VpcV1, sgt vpcv1.SecurityGroupTargetReferenceIntf, lbId, securityGroupID, securityGroupTargetID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		getSecurityGroupTargetOptions := &vpcv1.GetSecurityGroupTargetOptions{
			SecurityGroupID: &securityGroupID,
			ID:              &securityGroupTargetID,
		}
		_, response, err := sess.GetSecurityGroupTarget(getSecurityGroupTargetOptions)

		if err != nil {
			if response != nil && response.StatusCode == 404 {
				getlboptions := &vpcv1.GetLoadBalancerOptions{
					ID: &lbId,
				}
				lb, response, err := sess.GetLoadBalancer(getlboptions)
				if err != nil {
					return nil, "", fmt.Errorf("Error Getting Load Balancer : %s\n%s", err, response)
				}

				if *lb.ProvisioningStatus == "active" || *lb.ProvisioningStatus == "failed" {
					return sgt, isLBProvisioningDone, nil
				} else {
					return sgt, isLBProvisioning, nil
				}
			}
			return nil, isLBProvisioningDone, fmt.Errorf("error getting Security Group Target : %s\n%s", err, response)
		}
		return sgt, isLBProvisioning, nil
	}
}

func isWaitForLbSgTargetCreateAvailable(sess *vpcv1.VpcV1, lbId string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for load balancer (%s) to be available.", lbId)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isLBProvisioning, "update_pending"},
		Target:     []string{isLBProvisioningDone, ""},
		Refresh:    isLBSgTargetRefreshFunc(sess, lbId),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isLBSgTargetRefreshFunc(sess *vpcv1.VpcV1, lbId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		getlboptions := &vpcv1.GetLoadBalancerOptions{
			ID: &lbId,
		}
		lb, response, err := sess.GetLoadBalancer(getlboptions)
		if err != nil {
			return nil, "", fmt.Errorf("Error Getting Load Balancer : %s\n%s", err, response)
		}

		if *lb.ProvisioningStatus == "active" || *lb.ProvisioningStatus == "failed" {
			return lb, isLBProvisioningDone, nil
		}

		return lb, isLBProvisioning, nil
	}
}
