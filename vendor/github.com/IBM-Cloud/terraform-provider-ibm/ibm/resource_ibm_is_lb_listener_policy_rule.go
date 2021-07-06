// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcclassicv1"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

const (
	isLBListenerPolicyRuleLBID             = "lb"
	isLBListenerPolicyRuleListenerID       = "listener"
	isLBListenerPolicyRulePolicyID         = "policy"
	isLBListenerPolicyRuleid               = "rule"
	isLBListenerPolicyRulecondition        = "condition"
	isLBListenerPolicyRuletype             = "type"
	isLBListenerPolicyRulevalue            = "value"
	isLBListenerPolicyRulefield            = "field"
	isLBListenerPolicyRuleStatus           = "provisioning_status"
	isLBListenerPolicyRuleAvailable        = "active"
	isLBListenerPolicyRuleFailed           = "failed"
	isLBListenerPolicyRulePending          = "pending"
	isLBListenerPolicyRuleDeleting         = "deleting"
	isLBListenerPolicyRuleDeleted          = "done"
	isLBListenerPolicyRuleRetry            = "retry"
	isLBListenerPolicyRuleProvisioning     = "provisioning"
	isLBListenerPolicyRuleProvisioningDone = "done"
)

func resourceIBMISLBListenerPolicyRule() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMISLBListenerPolicyRuleCreate,
		Read:     resourceIBMISLBListenerPolicyRuleRead,
		Update:   resourceIBMISLBListenerPolicyRuleUpdate,
		Delete:   resourceIBMISLBListenerPolicyRuleDelete,
		Exists:   resourceIBMISLBListenerPolicyRuleExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			isLBListenerPolicyRuleLBID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Loadbalancer ID",
			},

			isLBListenerPolicyRuleListenerID: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, o, n string, d *schema.ResourceData) bool {
					if o == "" {
						return false
					}
					// if state file entry and tf file entry matches
					if strings.Compare(n, o) == 0 {
						return true
					}

					if strings.Contains(n, "/") {
						//Split lbID/listenerID and fetch listenerID
						new := strings.Split(n, "/")
						if strings.Compare(new[1], o) == 0 {
							return true
						}
					}

					return false
				},
				Description: "Listener ID.",
			},

			isLBListenerPolicyRulePolicyID: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, o, n string, d *schema.ResourceData) bool {
					if o == "" {
						return false
					}
					// if state file entry and tf file entry matches
					if strings.Compare(n, o) == 0 {
						return true
					}

					if strings.Contains(n, "/") {
						//Split lbID/listenerID and fetch listenerID
						new := strings.Split(n, "/")
						if strings.Compare(new[2], o) == 0 {
							return true
						}
					}

					return false
				},
				Description: "Listener Policy ID",
			},

			isLBListenerPolicyRulecondition: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: InvokeValidator("ibm_is_lb_listener_policy_rule", isLBListenerPolicyRulecondition),
				Description:  "Condition info of the rule.",
			},

			isLBListenerPolicyRuletype: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: InvokeValidator("ibm_is_lb_listener_policy_rule", isLBListenerPolicyRuletype),
				Description:  "Policy rule type.",
			},

			isLBListenerPolicyRulevalue: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLength,
				Description:  "policy rule value info",
			},

			isLBListenerPolicyRulefield: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLength,
			},

			isLBListenerPolicyRuleid: {
				Type:     schema.TypeString,
				Computed: true,
			},

			isLBListenerPolicyStatus: {
				Type:     schema.TypeString,
				Computed: true,
			},

			RelatedCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the LB resource",
			},
		},
	}
}

func resourceIBMISLBListenerPolicyRuleValidator() *ResourceValidator {

	validateSchema := make([]ValidateSchema, 1)
	condition := "contains, equals, matches_regex"
	ruletype := "header, hostname, path, body, query"
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isLBListenerPolicyRulecondition,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              condition})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isLBListenerPolicyRuletype,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              ruletype})

	ibmISLBListenerPolicyRuleResourceValidator := ResourceValidator{ResourceName: "ibm_is_lb_listener_policy_rule", Schema: validateSchema}
	return &ibmISLBListenerPolicyRuleResourceValidator
}

func resourceIBMISLBListenerPolicyRuleCreate(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	//Read lb, listerner, policy IDs
	var field string
	lbID := d.Get(isLBListenerPolicyRuleLBID).(string)
	listenerID, err := getLbListenerID(d.Get(isLBListenerPolicyRuleListenerID).(string))
	if err != nil {
		return err
	}

	policyID, err := getLbPolicyID(d.Get(isLBListenerPolicyRulePolicyID).(string))
	if err != nil {
		return err
	}

	condition := d.Get(isLBListenerPolicyRulecondition).(string)
	ty := d.Get(isLBListenerPolicyRuletype).(string)
	value := d.Get(isLBListenerPolicyRulevalue).(string)
	if n, ok := d.GetOk(isLBListenerPolicyRulefield); ok {
		field = n.(string)
	}

	if userDetails.generation == 1 {
		err := classicLbListenerPolicyRuleCreate(d, meta, lbID, listenerID, policyID, condition, ty, value, field)
		if err != nil {
			return err
		}
	} else {
		err := lbListenerPolicyRuleCreate(d, meta, lbID, listenerID, policyID, condition, ty, value, field)
		if err != nil {
			return err
		}
	}

	return resourceIBMISLBListenerPolicyRuleRead(d, meta)
}

func getLbListenerID(id string) (string, error) {
	if strings.Contains(id, "/") {
		parts, err := idParts(id)
		if err != nil {
			return "", err
		}

		return parts[1], nil
	} else {
		return id, nil
	}
}

func getLbPolicyID(id string) (string, error) {
	if strings.Contains(id, "/") {
		parts, err := idParts(id)
		if err != nil {
			return "", err
		}

		return parts[2], nil
	} else {
		return id, nil
	}
}

func classicLbListenerPolicyRuleCreate(d *schema.ResourceData, meta interface{}, lbID, listenerID, policyID, condition, ty, value, field string) error {
	sess, err := classicVpcSdkClient(meta)
	if err != nil {
		return err
	}

	options := &vpcclassicv1.CreateLoadBalancerListenerPolicyRuleOptions{
		LoadBalancerID: &lbID,
		ListenerID:     &listenerID,
		PolicyID:       &policyID,
		Condition:      &condition,
		Type:           &ty,
		Value:          &value,
		Field:          &field,
	}

	isLBKey := "load_balancer_key_" + lbID
	ibmMutexKV.Lock(isLBKey)
	defer ibmMutexKV.Unlock(isLBKey)

	_, err = isWaitForClassicLoadbalancerAvailable(sess, lbID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf(
			"LB-LP Error checking for load balancer (%s) is active: %s", lbID, err)
	}

	rule, response, err := sess.CreateLoadBalancerListenerPolicyRule(options)
	if err != nil {
		return fmt.Errorf("Error while creating lb listener policy for LB %s: Error %v Response %v", lbID, err, *response)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s/%s", lbID, listenerID, policyID, *(rule.ID)))

	_, err = isWaitForClassicLbListenerPolicyRuleAvailable(sess, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}
	return nil
}

func isWaitForClassicLoadbalancerAvailable(vpc *vpcclassicv1.VpcClassicV1, id string, timeout time.Duration) (interface{}, error) {

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isLBListenerPolicyRuleProvisioning, "create_pending", "update_pending", "maintenance_pending"},
		Target:     []string{isLBProvisioningDone},
		Refresh:    isLoadbalancerClassicRefreshFunc(vpc, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isLoadbalancerClassicRefreshFunc(vpc *vpcclassicv1.VpcClassicV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		getLbOptions := &vpcclassicv1.GetLoadBalancerOptions{
			ID: &id,
		}

		lb, _, err := vpc.GetLoadBalancer(getLbOptions)
		if err != nil {
			return nil, "", err
		}

		if *(lb.ProvisioningStatus) == isLBListenerPolicyAvailable || *lb.ProvisioningStatus == isLBListenerPolicyFailed {
			return lb, isLBProvisioningDone, nil
		}

		return lb, isLBProvisioning, nil
	}
}

func isWaitForClassicLbListenerPolicyRuleAvailable(vpc *vpcclassicv1.VpcClassicV1, id string, timeout time.Duration) (interface{}, error) {

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isLBListenerPolicyRuleProvisioning, "create_pending", "update_pending", "maintenance_pending"},
		Target:     []string{isLBListenerPolicyProvisioningDone},
		Refresh:    isLbListenerPolicyRuleClassicRefreshFunc(vpc, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isLbListenerPolicyRuleClassicRefreshFunc(vpc *vpcclassicv1.VpcClassicV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		parts, err := idParts(id)
		if err != nil {
			return nil, "", err
		}

		lbID := parts[0]
		listenerID := parts[1]
		policyID := parts[2]
		ruleID := parts[3]

		getLbListenerPolicyRuleOptions := &vpcclassicv1.GetLoadBalancerListenerPolicyRuleOptions{
			LoadBalancerID: &lbID,
			ListenerID:     &listenerID,
			PolicyID:       &policyID,
			ID:             &ruleID,
		}

		rule, _, err := vpc.GetLoadBalancerListenerPolicyRule(getLbListenerPolicyRuleOptions)

		if err != nil {
			return rule, "", err
		}

		if *rule.ProvisioningStatus == isLBListenerPolicyAvailable || *rule.ProvisioningStatus == isLBListenerPolicyFailed {
			return rule, isLBListenerProvisioningDone, nil
		}

		return rule, *rule.ProvisioningStatus, nil
	}
}

func vpcSdkClient(meta interface{}) (*vpcv1.VpcV1, error) {
	sess, err := meta.(ClientSession).VpcV1API()
	return sess, err
}

func lbListenerPolicyRuleCreate(d *schema.ResourceData, meta interface{}, lbID, listenerID, policyID, condition, ty, value, field string) error {

	sess, err := vpcSdkClient(meta)
	if err != nil {
		return err
	}

	options := &vpcv1.CreateLoadBalancerListenerPolicyRuleOptions{
		LoadBalancerID: &lbID,
		ListenerID:     &listenerID,
		PolicyID:       &policyID,
		Condition:      &condition,
		Type:           &ty,
		Value:          &value,
		Field:          &field,
	}

	isLBKey := "load_balancer_key_" + lbID
	ibmMutexKV.Lock(isLBKey)
	defer ibmMutexKV.Unlock(isLBKey)

	_, err = isWaitForLoadbalancerAvailable(sess, lbID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf(
			"LB-LP Error checking for load balancer (%s) is active: %s", lbID, err)
	}

	rule, response, err := sess.CreateLoadBalancerListenerPolicyRule(options)
	if err != nil {
		return fmt.Errorf("Error while creating lb listener policy for LB %s: Error %v Response %v", lbID, err, *response)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s/%s", lbID, listenerID, policyID, *(rule.ID)))

	_, err = isWaitForLbListenerPolicyRuleAvailable(sess, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	return nil
}

func isWaitForLoadbalancerAvailable(vpc *vpcv1.VpcV1, id string, timeout time.Duration) (interface{}, error) {

	stateConf := &resource.StateChangeConf{
		Pending:    []string{isLBListenerPolicyRulePending},
		Target:     []string{isLBProvisioningDone},
		Refresh:    isLoadbalancerRefreshFunc(vpc, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isLoadbalancerRefreshFunc(vpc *vpcv1.VpcV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		getLbOptions := &vpcv1.GetLoadBalancerOptions{
			ID: &id,
		}

		lb, _, err := vpc.GetLoadBalancer(getLbOptions)
		if err != nil {
			return nil, "", err
		}

		if *(lb.ProvisioningStatus) == isLBListenerPolicyAvailable {
			return lb, isLBProvisioningDone, nil
		}

		return lb, isLBProvisioning, nil
	}
}

func isWaitForLbListenerPolicyRuleAvailable(vpc *vpcv1.VpcV1, id string, timeout time.Duration) (interface{}, error) {

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isLBListenerPolicyRuleProvisioning, "create_pending", "update_pending", "maintenance_pending"},
		Target:     []string{isLBListenerPolicyRuleProvisioningDone},
		Refresh:    isLbListenerPolicyRuleRefreshFunc(vpc, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isLbListenerPolicyRuleRefreshFunc(vpc *vpcv1.VpcV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		parts, err := idParts(id)
		if err != nil {
			return nil, "", err
		}

		lbID := parts[0]
		listenerID := parts[1]
		policyID := parts[2]
		ruleID := parts[3]

		getLbListenerPolicyRuleOptions := &vpcv1.GetLoadBalancerListenerPolicyRuleOptions{
			LoadBalancerID: &lbID,
			ListenerID:     &listenerID,
			PolicyID:       &policyID,
			ID:             &ruleID,
		}

		rule, _, err := vpc.GetLoadBalancerListenerPolicyRule(getLbListenerPolicyRuleOptions)

		if err != nil {
			return rule, "", err
		}

		if *rule.ProvisioningStatus == isLBListenerPolicyRuleAvailable || *rule.ProvisioningStatus == isLBListenerPolicyRuleFailed {
			return rule, isLBListenerPolicyRuleProvisioningDone, nil
		}

		return rule, *rule.ProvisioningStatus, nil
	}
}

func resourceIBMISLBListenerPolicyRuleRead(d *schema.ResourceData, meta interface{}) error {

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	ID := d.Id()
	parts, err := idParts(ID)
	if err != nil {
		return err
	}

	lbID := parts[0]
	listenerID := parts[1]
	policyID := parts[2]
	ruleID := parts[3]

	if userDetails.generation == 1 {
		err := classicLbListenerPolicyRuleGet(d, meta, lbID, listenerID, policyID, ruleID)
		if err != nil {
			return err
		}
	} else {
		err := lbListenerPolicyRuleGet(d, meta, lbID, listenerID, policyID, ruleID)
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceIBMISLBListenerPolicyRuleExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return false, err
	}
	ID := d.Id()
	if userDetails.generation == 1 {
		exists, err := classicLbListenerPolicyRuleExists(d, meta, ID)
		return exists, err
	} else {
		exists, err := lbListenerPolicyRuleExists(d, meta, ID)
		return exists, err
	}
}

func classicLbListenerPolicyRuleExists(d *schema.ResourceData, meta interface{}, ID string) (bool, error) {
	sess, err := classicVpcSdkClient(meta)
	if err != nil {
		return false, err
	}

	//Retrieve lbID, listenerID and policyID
	parts, err := idParts(d.Id())
	if err != nil {
		return false, err
	}

	lbID := parts[0]
	listenerID := parts[1]
	policyID := parts[2]
	ruleID := parts[3]

	//populate lblistenerpolicyOPtions
	getLbListenerPolicyRuleOptions := &vpcclassicv1.GetLoadBalancerListenerPolicyRuleOptions{
		LoadBalancerID: &lbID,
		ListenerID:     &listenerID,
		PolicyID:       &policyID,
		ID:             &ruleID,
	}

	//Getting lb listener policy
	_, response, err := sess.GetLoadBalancerListenerPolicyRule(getLbListenerPolicyRuleOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("Error getting policy: %s\n%s", err, response)
	}
	return true, nil
}

func lbListenerPolicyRuleExists(d *schema.ResourceData, meta interface{}, ID string) (bool, error) {
	sess, err := vpcSdkClient(meta)
	if err != nil {
		return false, err
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return false, err
	}
	if len(parts) != 4 {
		return false, fmt.Errorf("Incorrect ID %s: ID should be a combination of lbID/listenerID/policyID/ruleID", d.Id())
	}
	lbID := parts[0]
	listenerID := parts[1]
	policyID := parts[2]
	ruleID := parts[3]

	getLbListenerPolicyRuleOptions := &vpcv1.GetLoadBalancerListenerPolicyRuleOptions{
		LoadBalancerID: &lbID,
		ListenerID:     &listenerID,
		PolicyID:       &policyID,
		ID:             &ruleID,
	}

	//Getting lb listener policy
	_, response, err := sess.GetLoadBalancerListenerPolicyRule(getLbListenerPolicyRuleOptions)

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("Error getting policy: %s\n%s", err, response)
	}
	return true, nil
}
func resourceIBMISLBListenerPolicyRuleUpdate(d *schema.ResourceData, meta interface{}) error {

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}

	lbID := parts[0]
	listenerID := parts[1]
	policyID := parts[2]
	ruleID := parts[3]

	if userDetails.generation == 1 {

		err := classicLbListenerPolicyRuleUpdate(d, meta, lbID, listenerID, policyID, ruleID)
		if err != nil {
			return err
		}
	} else {

		err := lbListenerPolicyRuleUpdate(d, meta, lbID, listenerID, policyID, ruleID)
		if err != nil {
			return err
		}
	}

	return resourceIBMISLBListenerPolicyRuleRead(d, meta)
}

func classicLbListenerPolicyRuleUpdate(d *schema.ResourceData, meta interface{}, lbID, listenerID, policyID, ID string) error {
	sess, err := classicVpcSdkClient(meta)
	if err != nil {
		return err
	}

	hasChanged := false
	updatePolicyRuleOptions := vpcclassicv1.UpdateLoadBalancerListenerPolicyRuleOptions{}
	updatePolicyRuleOptions.LoadBalancerID = &lbID
	updatePolicyRuleOptions.ListenerID = &listenerID
	updatePolicyRuleOptions.PolicyID = &policyID
	updatePolicyRuleOptions.ID = &ID

	loadBalancerListenerPolicyRulePatchModel := &vpcclassicv1.LoadBalancerListenerPolicyRulePatch{}

	if d.HasChange(isLBListenerPolicyRulecondition) {
		condition := d.Get(isLBListenerPolicyRulecondition).(string)
		loadBalancerListenerPolicyRulePatchModel.Condition = &condition
		hasChanged = true
	}

	if d.HasChange(isLBListenerPolicyRuletype) {
		ty := d.Get(isLBListenerPolicyRuletype).(string)
		loadBalancerListenerPolicyRulePatchModel.Type = &ty
		hasChanged = true
	}

	if d.HasChange(isLBListenerPolicyRulevalue) {
		value := d.Get(isLBListenerPolicyRulevalue).(string)
		loadBalancerListenerPolicyRulePatchModel.Value = &value
		hasChanged = true
	}

	if d.HasChange(isLBListenerPolicyRulefield) {
		field := d.Get(isLBListenerPolicyRulefield).(string)
		loadBalancerListenerPolicyRulePatchModel.Field = &field
		hasChanged = true
	}

	isLBKey := "load_balancer_key_" + lbID
	ibmMutexKV.Lock(isLBKey)
	defer ibmMutexKV.Unlock(isLBKey)

	if hasChanged {
		loadBalancerListenerPolicyRulePatch, err := loadBalancerListenerPolicyRulePatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("Error calling asPatch for LoadBalancerListenerPolicyRulePatch: %s", err)
		}
		updatePolicyRuleOptions.LoadBalancerListenerPolicyRulePatch = loadBalancerListenerPolicyRulePatch

		_, err = isWaitForClassicLoadbalancerAvailable(sess, lbID, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return fmt.Errorf(
				"LB-LP Error checking for load balancer (%s) is active: %s", lbID, err)
		}
		_, response, err := sess.UpdateLoadBalancerListenerPolicyRule(&updatePolicyRuleOptions)
		if err != nil {
			return fmt.Errorf("Error Getting Instance: %s\n%s", err, response)
		}

		_, err = isWaitForClassicLbListenerPolicyRuleAvailable(sess, d.Id(), d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return err
		}
	}
	return nil
}

func lbListenerPolicyRuleUpdate(d *schema.ResourceData, meta interface{}, lbID, listenerID, policyID, ID string) error {
	sess, err := vpcSdkClient(meta)
	if err != nil {
		return err
	}
	hasChanged := false
	updatePolicyRuleOptions := vpcv1.UpdateLoadBalancerListenerPolicyRuleOptions{}
	updatePolicyRuleOptions.LoadBalancerID = &lbID
	updatePolicyRuleOptions.ListenerID = &listenerID
	updatePolicyRuleOptions.PolicyID = &policyID
	updatePolicyRuleOptions.ID = &ID

	loadBalancerListenerPolicyRulePatchModel := &vpcv1.LoadBalancerListenerPolicyRulePatch{}

	if d.HasChange(isLBListenerPolicyRulecondition) {
		condition := d.Get(isLBListenerPolicyRulecondition).(string)
		loadBalancerListenerPolicyRulePatchModel.Condition = &condition
		hasChanged = true
	}

	if d.HasChange(isLBListenerPolicyRuletype) {
		ty := d.Get(isLBListenerPolicyRuletype).(string)
		loadBalancerListenerPolicyRulePatchModel.Type = &ty
		hasChanged = true
	}

	if d.HasChange(isLBListenerPolicyRulevalue) {
		value := d.Get(isLBListenerPolicyRulevalue).(string)
		loadBalancerListenerPolicyRulePatchModel.Value = &value
		hasChanged = true
	}

	if d.HasChange(isLBListenerPolicyRulefield) {
		field := d.Get(isLBListenerPolicyRulefield).(string)
		loadBalancerListenerPolicyRulePatchModel.Field = &field
		hasChanged = true
	}

	if hasChanged {
		loadBalancerListenerPolicyRulePatch, err := loadBalancerListenerPolicyRulePatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("Error calling asPatch for LoadBalancerListenerPolicyRulePatch: %s", err)
		}
		updatePolicyRuleOptions.LoadBalancerListenerPolicyRulePatch = loadBalancerListenerPolicyRulePatch

		isLBKey := "load_balancer_key_" + lbID
		ibmMutexKV.Lock(isLBKey)
		defer ibmMutexKV.Unlock(isLBKey)

		_, err = isWaitForLoadbalancerAvailable(sess, lbID, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return fmt.Errorf(
				"LB-LP Error checking for load balancer (%s) is active: %s", lbID, err)
		}

		_, response, err := sess.UpdateLoadBalancerListenerPolicyRule(&updatePolicyRuleOptions)
		if err != nil {
			return fmt.Errorf("Error Updating in policy : %s\n%s", err, response)
		}

		_, err = isWaitForLbListenerPolicyRuleAvailable(sess, d.Id(), d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceIBMISLBListenerPolicyRuleDelete(d *schema.ResourceData, meta interface{}) error {

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	//Retrieve lbId, listenerId and policyID
	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}

	lbID := parts[0]
	listenerID := parts[1]
	policyID := parts[2]
	ruleID := parts[3]

	isLBKey := "load_balancer_key_" + lbID
	ibmMutexKV.Lock(isLBKey)
	defer ibmMutexKV.Unlock(isLBKey)

	if userDetails.generation == 1 {
		err := classicLbListenerPolicyRuleDelete(d, meta, lbID, listenerID, policyID, ruleID)
		if err != nil {
			return err
		}
	} else {
		err := lbListenerPolicyRuleDelete(d, meta, lbID, listenerID, policyID, ruleID)
		if err != nil {
			return err
		}
	}

	d.SetId("")
	return nil

}

func classicLbListenerPolicyRuleDelete(d *schema.ResourceData, meta interface{}, lbID, listenerID, policyID, ID string) error {
	sess, err := classicVpcSdkClient(meta)
	if err != nil {
		return err
	}

	//Getting rule optins
	getLbListenerPolicyRuleOptions := &vpcclassicv1.GetLoadBalancerListenerPolicyRuleOptions{
		LoadBalancerID: &lbID,
		ListenerID:     &listenerID,
		PolicyID:       &policyID,
		ID:             &ID,
	}

	//Getting lb listener policy
	_, response, err := sess.GetLoadBalancerListenerPolicyRule(getLbListenerPolicyRuleOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error in classicLbListenerPolicyGet : %s\n%s", err, response)
	}

	deleteLbListenerPolicyRuleOptions := &vpcclassicv1.DeleteLoadBalancerListenerPolicyRuleOptions{
		LoadBalancerID: &lbID,
		ListenerID:     &listenerID,
		PolicyID:       &policyID,
		ID:             &ID,
	}

	response, err = sess.DeleteLoadBalancerListenerPolicyRule(deleteLbListenerPolicyRuleOptions)

	if err != nil {
		return fmt.Errorf("Error in classicLbListenerPolicyRuleDelete: %s\n%s", err, response)
	}
	_, err = isWaitForLbListenerPolicyRuleClassicDeleted(sess, d.Id(), d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}
	return nil
}

func lbListenerPolicyRuleDelete(d *schema.ResourceData, meta interface{}, lbID, listenerID, policyID, ID string) error {

	sess, err := vpcSdkClient(meta)
	if err != nil {
		return err
	}

	//Getting rule optins
	getLbListenerPolicyRuleOptions := &vpcv1.GetLoadBalancerListenerPolicyRuleOptions{
		LoadBalancerID: &lbID,
		ListenerID:     &listenerID,
		PolicyID:       &policyID,
		ID:             &ID,
	}

	//Getting lb listener policy
	_, response, err := sess.GetLoadBalancerListenerPolicyRule(getLbListenerPolicyRuleOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error in LbListenerPolicyGet : %s\n%s", err, response)
	}

	deleteLbListenerPolicyRuleOptions := &vpcv1.DeleteLoadBalancerListenerPolicyRuleOptions{
		LoadBalancerID: &lbID,
		ListenerID:     &listenerID,
		PolicyID:       &policyID,
		ID:             &ID,
	}
	response, err = sess.DeleteLoadBalancerListenerPolicyRule(deleteLbListenerPolicyRuleOptions)
	if err != nil {
		return fmt.Errorf("Error in lbListenerPolicyRuleDelete: %s\n%s", err, response)
	}
	_, err = isWaitForLbListnerPolicyRuleDeleted(sess, d.Id(), d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}
	return nil
}
func isWaitForLbListnerPolicyRuleDeleted(vpc *vpcv1.VpcV1, id string, timeout time.Duration) (interface{}, error) {

	stateConf := &resource.StateChangeConf{
		Pending:    []string{isLBListenerPolicyRuleRetry, isLBListenerPolicyRuleDeleting},
		Target:     []string{isLBListenerPolicyRuleDeleted, isLBListenerPolicyRuleFailed},
		Refresh:    isLbListenerPolicyRuleDeleteRefreshFunc(vpc, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isLbListenerPolicyRuleDeleteRefreshFunc(vpc *vpcv1.VpcV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		//Retrieve lbId, listenerId and policyID
		parts, err := idParts(id)
		if err != nil {
			return nil, isLBListenerPolicyFailed, nil
		}

		lbID := parts[0]
		listenerID := parts[1]
		policyID := parts[2]
		ruleID := parts[3]

		getLbListenerPolicyRuleOptions := &vpcv1.GetLoadBalancerListenerPolicyRuleOptions{
			LoadBalancerID: &lbID,
			ListenerID:     &listenerID,
			PolicyID:       &policyID,
			ID:             &ruleID,
		}

		//Getting lb listener policy
		rule, response, err := vpc.GetLoadBalancerListenerPolicyRule(getLbListenerPolicyRuleOptions)

		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return rule, isLBListenerPolicyRuleDeleted, nil
			}
			return rule, isLBListenerPolicyRuleFailed, err
		}
		return nil, isLBListenerPolicyRuleDeleting, err
	}
}

func classicVpcSdkClient(meta interface{}) (*vpcclassicv1.VpcClassicV1, error) {
	sess, err := meta.(ClientSession).VpcClassicV1API()
	return sess, err
}

func classicLbListenerPolicyRuleGet(d *schema.ResourceData, meta interface{}, lbID, listenerID, policyID, id string) error {
	sess, err := classicVpcSdkClient(meta)
	if err != nil {
		return err
	}

	//Getting rule optins
	getLbListenerPolicyRuleOptions := &vpcclassicv1.GetLoadBalancerListenerPolicyRuleOptions{
		LoadBalancerID: &lbID,
		ListenerID:     &listenerID,
		PolicyID:       &policyID,
		ID:             &id,
	}

	//Getting lb listener policy
	rule, response, err := sess.GetLoadBalancerListenerPolicyRule(getLbListenerPolicyRuleOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error in classicLbListenerPolicyGet : %s\n%s", err, response)
	}

	d.Set(isLBListenerPolicyRuleLBID, lbID)
	d.Set(isLBListenerPolicyRuleListenerID, listenerID)
	d.Set(isLBListenerPolicyRulePolicyID, policyID)
	d.Set(isLBListenerPolicyRuleid, id)
	d.Set(isLBListenerPolicyRulecondition, rule.Condition)
	d.Set(isLBListenerPolicyRuletype, rule.Type)
	d.Set(isLBListenerPolicyRulevalue, rule.Value)
	d.Set(isLBListenerPolicyRulefield, rule.Field)
	d.Set(isLBListenerPolicyRuleStatus, rule.ProvisioningStatus)
	getLoadBalancerOptions := &vpcclassicv1.GetLoadBalancerOptions{
		ID: &lbID,
	}
	lb, response, err := sess.GetLoadBalancer(getLoadBalancerOptions)
	if err != nil {
		return fmt.Errorf("Error Getting Load Balancer : %s\n%s", err, response)
	}
	d.Set(RelatedCRN, *lb.CRN)

	return nil
}

func lbListenerPolicyRuleGet(d *schema.ResourceData, meta interface{}, lbID, listenerID, policyID, id string) error {

	sess, err := vpcSdkClient(meta)
	if err != nil {
		return err
	}

	//Getting rule optins
	getLbListenerPolicyRuleOptions := &vpcv1.GetLoadBalancerListenerPolicyRuleOptions{
		LoadBalancerID: &lbID,
		ListenerID:     &listenerID,
		PolicyID:       &policyID,
		ID:             &id,
	}

	//Getting lb listener policy
	rule, response, err := sess.GetLoadBalancerListenerPolicyRule(getLbListenerPolicyRuleOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return err
	}

	//set the argument values
	d.Set(isLBListenerPolicyRuleLBID, lbID)
	d.Set(isLBListenerPolicyRuleListenerID, listenerID)
	d.Set(isLBListenerPolicyRulePolicyID, policyID)
	d.Set(isLBListenerPolicyRuleid, id)
	d.Set(isLBListenerPolicyRulecondition, rule.Condition)
	d.Set(isLBListenerPolicyRuletype, rule.Type)
	d.Set(isLBListenerPolicyRulevalue, rule.Value)
	d.Set(isLBListenerPolicyRulefield, rule.Field)
	d.Set(isLBListenerPolicyRuleStatus, rule.ProvisioningStatus)
	getLoadBalancerOptions := &vpcv1.GetLoadBalancerOptions{
		ID: &lbID,
	}
	lb, response, err := sess.GetLoadBalancer(getLoadBalancerOptions)
	if err != nil {
		return fmt.Errorf("Error Getting Load Balancer : %s\n%s", err, response)
	}
	d.Set(RelatedCRN, *lb.CRN)

	return nil
}

func isWaitForLbListenerPolicyRuleClassicDeleted(vpc *vpcclassicv1.VpcClassicV1, id string, timeout time.Duration) (interface{}, error) {

	stateConf := &resource.StateChangeConf{
		Pending:    []string{isLBListenerPolicyRuleRetry, isLBListenerPolicyRuleDeleting, "delete_pending"},
		Target:     []string{isLBListenerPolicyRuleDeleted, isLBListenerPolicyRuleFailed},
		Refresh:    isLbListenerPolicyRuleClassicDeleteRefreshFunc(vpc, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isLbListenerPolicyRuleClassicDeleteRefreshFunc(vpc *vpcclassicv1.VpcClassicV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		//Retrieve lbId and listenerId
		parts, err := idParts(id)
		if err != nil {
			return nil, isLBListenerPolicyFailed, nil
		}

		lbID := parts[0]
		listenerID := parts[1]
		policyID := parts[2]
		ruleID := parts[3]

		getLbListenerPolicyRuleOptions := &vpcclassicv1.GetLoadBalancerListenerPolicyRuleOptions{
			LoadBalancerID: &lbID,
			ListenerID:     &listenerID,
			PolicyID:       &policyID,
			ID:             &ruleID,
		}

		//Getting lb listener policy
		rule, response, err := vpc.GetLoadBalancerListenerPolicyRule(getLbListenerPolicyRuleOptions)
		//failed := isLBListenerPolicyRuleFailed
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return rule, isLBListenerPolicyRuleDeleted, nil
			}
			return nil, isLBListenerPolicyRuleFailed, err
		}
		return rule, isLBListenerPolicyRuleDeleting, err
	}
}
