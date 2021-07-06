// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcclassicv1"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	isLBListenerPolicyLBID                 = "lb"
	isLBListenerPolicyListenerID           = "listener"
	isLBListenerPolicyAction               = "action"
	isLBListenerPolicyPriority             = "priority"
	isLBListenerPolicyName                 = "name"
	isLBListenerPolicyID                   = "policy_id"
	isLBListenerPolicyRules                = "rules"
	isLBListenerPolicyRulesInfo            = "rule_info"
	isLBListenerPolicyTargetID             = "target_id"
	isLBListenerPolicyTargetHTTPStatusCode = "target_http_status_code"
	isLBListenerPolicyTargetURL            = "target_url"
	isLBListenerPolicyStatus               = "provisioning_status"
	isLBListenerPolicyRuleID               = "rule_id"
	isLBListenerPolicyAvailable            = "active"
	isLBListenerPolicyFailed               = "failed"
	isLBListenerPolicyPending              = "pending"
	isLBListenerPolicyDeleting             = "deleting"
	isLBListenerPolicyDeleted              = "done"
	isLBListenerPolicyRetry                = "retry"
	isLBListenerPolicyRuleCondition        = "condition"
	isLBListenerPolicyRuleType             = "type"
	isLBListenerPolicyRuleValue            = "value"
	isLBListenerPolicyRuleField            = "field"
	isLBListenerPolicyProvisioning         = "provisioning"
	isLBListenerPolicyProvisioningDone     = "done"
)

func resourceIBMISLBListenerPolicy() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMISLBListenerPolicyCreate,
		Read:     resourceIBMISLBListenerPolicyRead,
		Update:   resourceIBMISLBListenerPolicyUpdate,
		Delete:   resourceIBMISLBListenerPolicyDelete,
		Exists:   resourceIBMISLBListenerPolicyExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			isLBListenerPolicyLBID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Load Balancer Listener Policy",
			},

			isLBListenerPolicyListenerID: {
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
				Description: "Listener ID",
			},

			isLBListenerPolicyAction: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: InvokeValidator("ibm_is_lb_listener_policy", isLBListenerPolicyAction),
				Description:  "Policy Action",
			},

			isLBListenerPolicyPriority: {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validateLBListenerPolicyPriority,
				Description:  "Listener Policy Priority",
			},

			isLBListenerPolicyName: {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     false,
				Computed:     true,
				ValidateFunc: InvokeValidator("ibm_is_lb_listener_policy", isLBListenerPolicyName),
				Description:  "Policy name",
			},

			isLBListenerPolicyID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Listener Policy ID",
			},

			isLBListenerPolicyRules: {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Policy Rules",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isLBListenerPolicyRuleCondition: {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: InvokeValidator("ibm_is_lb_listener_policy_rule", isLBListenerPolicyRulecondition),
							Description:  "Condition of the rule",
						},

						isLBListenerPolicyRuleType: {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: InvokeValidator("ibm_is_lb_listener_policy_rule", isLBListenerPolicyRuleType),
							Description:  "Type of the rule",
						},

						isLBListenerPolicyRuleValue: {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateStringLength,
							Description:  "Value to be matched for rule condition",
						},

						isLBListenerPolicyRuleField: {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateStringLength,
							Description:  "HTTP header field. This is only applicable to rule type.",
						},

						isLBListenerPolicyRuleID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule ID",
						},
					},
				},
			},

			isLBListenerPolicyTargetID: {
				Type:     schema.TypeString,
				ForceNew: false,
				Optional: true,
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
				Description: "Listener Policy Target ID",
			},

			isLBListenerPolicyTargetHTTPStatusCode: {
				Type:        schema.TypeInt,
				ForceNew:    false,
				Optional:    true,
				Description: "Listener Policy target HTTPS Status code.",
			},

			isLBListenerPolicyTargetURL: {
				Type:        schema.TypeString,
				ForceNew:    false,
				Optional:    true,
				Description: "Policy Target URL",
			},

			isLBListenerPolicyStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Listner Policy status",
			},

			RelatedCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the LB resource",
			},
		},
	}
}

func resourceIBMISLBListenerPolicyValidator() *ResourceValidator {

	validateSchema := make([]ValidateSchema, 1)
	action := "forward, redirect, reject"
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isLBListenerPolicyName,
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Required:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isLBListenerPolicyAction,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              action})

	ibmISLBListenerPolicyResourceValidator := ResourceValidator{ResourceName: "ibm_is_lb_listener_policy", Schema: validateSchema}
	return &ibmISLBListenerPolicyResourceValidator
}

func resourceIBMISLBListenerPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	//Get the Load balancer ID
	lbID := d.Get(isLBListenerPolicyLBID).(string)

	//User can set listener id as combination of lbID/listenerID, parse and get the listenerID
	listenerID, err := getListenerID(d.Get(isLBListenerPolicyListenerID).(string))
	if err != nil {
		return err
	}

	action := d.Get(isLBListenerPolicyAction).(string)
	priority := int64(d.Get(isLBListenerPolicyPriority).(int))

	//user-defined name for this policy.
	var name string
	if n, ok := d.GetOk(isLBListenerPolicyName); ok {
		name = n.(string)
	}

	if userDetails.generation == 1 {
		err := classicLbListenerPolicyCreate(d, meta, lbID, listenerID, action, name, priority)
		if err != nil {
			return err
		}
	} else {
		err := lbListenerPolicyCreate(d, meta, lbID, listenerID, action, name, priority)
		if err != nil {
			return err
		}
	}

	return resourceIBMISLBListenerPolicyRead(d, meta)
}

func getListenerID(id string) (string, error) {
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

func classicLbListenerPolicyCreate(d *schema.ResourceData, meta interface{}, lbID, listenerID, action, name string, priority int64) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}

	// When `action` is `forward`, `LoadBalancerPoolIdentity` is required to specify which
	// pool the load balancer forwards the traffic to. When `action` is `redirect`,
	// `LoadBalancerListenerPolicyRedirectURLPrototype` is required to specify the url and
	// http status code used in the redirect response.

	actionChk := d.Get(isLBListenerPolicyAction)
	tID, targetIDSet := d.GetOk(isLBListenerPolicyTargetID)
	statusCode, statusSet := d.GetOk(isLBListenerPolicyTargetHTTPStatusCode)
	url, urlSet := d.GetOk(isLBListenerPolicyTargetURL)

	var target vpcclassicv1.LoadBalancerListenerPolicyTargetPrototypeIntf

	if actionChk.(string) == "forward" {
		if targetIDSet {

			//User can set the poolId as combination of lbID/poolID, if so parse the string & get the poolID
			id, err := getPoolID(tID.(string))
			if err != nil {
				return err
			}

			//id := lbPoolID.(string)
			target = &vpcclassicv1.LoadBalancerListenerPolicyTargetPrototypeLoadBalancerPoolIdentity{
				ID: &id,
			}
		} else {
			return fmt.Errorf("When action is forward please specify target_id")
		}
	} else if actionChk.(string) == "redirect" {

		urlPrototype := vpcclassicv1.LoadBalancerListenerPolicyTargetPrototypeLoadBalancerListenerPolicyRedirectURLPrototype{}

		if statusSet {
			sc := int64(statusCode.(int))
			urlPrototype.HTTPStatusCode = &sc
		} else {
			return fmt.Errorf("When action is redirect please specify target_http_status_code")
		}

		if urlSet {
			link := url.(string)
			urlPrototype.URL = &link
		} else {
			return fmt.Errorf("When action is redirect please specify target_url")
		}

		target = &urlPrototype
	}

	rulesInfo := make([]vpcclassicv1.LoadBalancerListenerPolicyRulePrototype, 0)
	if rules, rulesSet := d.GetOk(isLBListenerPolicyRules); rulesSet {
		policyRules := rules.([]interface{})
		for _, rule := range policyRules {
			rulex := rule.(map[string]interface{})

			//condition, type and value are mandatory params
			var condition string
			if rulex[isLBListenerPolicyRuleCondition] != nil {
				condition = rulex[isLBListenerPolicyRuleCondition].(string)
			}

			var ty string
			if rulex[isLBListenerPolicyRuleType] != nil {
				ty = rulex[isLBListenerPolicyRuleType].(string)
			}

			var value string
			if rulex[isLBListenerPolicyRuleValue] != nil {
				value = rulex[isLBListenerPolicyRuleValue].(string)
			}

			field := rulex[isLBListenerPolicyRuleField].(string)

			r := vpcclassicv1.LoadBalancerListenerPolicyRulePrototype{
				Condition: &condition,
				Field:     &field,
				Type:      &ty,
				Value:     &value,
			}

			rulesInfo = append(rulesInfo, r)
		}
	}

	options := &vpcclassicv1.CreateLoadBalancerListenerPolicyOptions{
		LoadBalancerID: &lbID,
		ListenerID:     &listenerID,
		Action:         &action,
		Priority:       &priority,
		Name:           &name,
		Target:         target,
		Rules:          rulesInfo,
	}

	isLBKey := "load_balancer_key_" + lbID
	ibmMutexKV.Lock(isLBKey)
	defer ibmMutexKV.Unlock(isLBKey)

	_, err = isWaitForClassicLbAvailable(sess, lbID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf(
			"LB-LP Error checking for load balancer (%s) is active: %s", lbID, err)
	}

	policy, response, err := sess.CreateLoadBalancerListenerPolicy(options)
	if err != nil {
		return fmt.Errorf("Error while creating lb listener policy for LB %s: Error %v Response %v", lbID, err, *response)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", lbID, listenerID, *(policy.ID)))

	_, err = isWaitForClassicLbListenerPolicyAvailable(sess, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}
	return nil
}

func getPoolID(id string) (string, error) {
	if strings.Contains(id, "/") {
		parts, err := idParts(id)
		if err != nil {
			return "", err
		}

		return parts[1], nil
	}
	return id, nil

}

func isWaitForClassicLbAvailable(vpc *vpcclassicv1.VpcClassicV1, id string, timeout time.Duration) (interface{}, error) {

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isLBListenerProvisioning, "create_pending", "update_pending", "maintenance_pending"},
		Target:     []string{isLBProvisioningDone},
		Refresh:    isLbClassicRefreshFunc(vpc, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isLbClassicRefreshFunc(vpc *vpcclassicv1.VpcClassicV1, id string) resource.StateRefreshFunc {
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

func isWaitForClassicLbListenerPolicyAvailable(vpc *vpcclassicv1.VpcClassicV1, id string, timeout time.Duration) (interface{}, error) {

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isLBListenerProvisioning, "create_pending", "update_pending", "maintenance_pending"},
		Target:     []string{isLBListenerPolicyProvisioningDone},
		Refresh:    isLbListenerPolicyClassicRefreshFunc(vpc, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isLbListenerPolicyClassicRefreshFunc(vpc *vpcclassicv1.VpcClassicV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		parts, err := idParts(id)
		if err != nil {
			return nil, "", err
		}

		lbID := parts[0]
		listenerID := parts[1]
		policyID := parts[2]

		getLbListenerPolicyOptions := &vpcclassicv1.GetLoadBalancerListenerPolicyOptions{
			LoadBalancerID: &lbID,
			ListenerID:     &listenerID,
			ID:             &policyID,
		}

		policy, _, err := vpc.GetLoadBalancerListenerPolicy(getLbListenerPolicyOptions)

		if err != nil {
			return policy, "", err
		}

		if *policy.ProvisioningStatus == isLBListenerPolicyAvailable || *policy.ProvisioningStatus == isLBListenerPolicyFailed {
			return policy, isLBListenerProvisioningDone, nil
		}

		return policy, *policy.ProvisioningStatus, nil
	}
}

func lbListenerPolicyCreate(d *schema.ResourceData, meta interface{}, lbID, listenerID, action, name string, priority int64) error {

	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	// When `action` is `forward`, `LoadBalancerPoolIdentity` is required to specify which
	// pool the load balancer forwards the traffic to. When `action` is `redirect`,
	// `LoadBalancerListenerPolicyRedirectURLPrototype` is required to specify the url and
	// http status code used in the redirect response.
	actionChk := d.Get(isLBListenerPolicyAction)
	tID, targetIDSet := d.GetOk(isLBListenerPolicyTargetID)
	statusCode, statusSet := d.GetOk(isLBListenerPolicyTargetHTTPStatusCode)
	url, urlSet := d.GetOk(isLBListenerPolicyTargetURL)

	var target vpcv1.LoadBalancerListenerPolicyTargetPrototypeIntf

	if actionChk.(string) == "forward" {
		if targetIDSet {

			//User can set the poolId as combination of lbID/poolID, if so parse the string & get the poolID
			id, err := getPoolID(tID.(string))
			if err != nil {
				return err
			}

			target = &vpcv1.LoadBalancerListenerPolicyTargetPrototypeLoadBalancerPoolIdentity{
				ID: &id,
			}
		} else {
			return fmt.Errorf("When action is forward please specify target_id")
		}
	} else if actionChk.(string) == "redirect" {

		urlPrototype := vpcv1.LoadBalancerListenerPolicyTargetPrototypeLoadBalancerListenerPolicyRedirectURLPrototype{}

		if statusSet {
			sc := int64(statusCode.(int))
			urlPrototype.HTTPStatusCode = &sc
		} else {
			return fmt.Errorf("When action is redirect please specify target_http_status_code")
		}

		if urlSet {
			link := url.(string)
			urlPrototype.URL = &link
		} else {
			return fmt.Errorf("When action is redirect please specify target_url")
		}

		target = &urlPrototype
	}

	//Read Rules
	rulesInfo := make([]vpcv1.LoadBalancerListenerPolicyRulePrototype, 0)
	if rules, rulesSet := d.GetOk(isLBListenerPolicyRules); rulesSet {
		policyRules := rules.([]interface{})
		for _, rule := range policyRules {
			rulex := rule.(map[string]interface{})

			//condition, type and value are mandatory params
			var condition string
			if rulex[isLBListenerPolicyRuleCondition] != nil {
				condition = rulex[isLBListenerPolicyRuleCondition].(string)
			}

			var ty string
			if rulex[isLBListenerPolicyRuleType] != nil {
				ty = rulex[isLBListenerPolicyRuleType].(string)
			}

			var value string
			if rulex[isLBListenerPolicyRuleValue] != nil {
				value = rulex[isLBListenerPolicyRuleValue].(string)
			}

			field := rulex[isLBListenerPolicyRuleField].(string)

			r := vpcv1.LoadBalancerListenerPolicyRulePrototype{
				Condition: &condition,
				Field:     &field,
				Type:      &ty,
				Value:     &value,
			}

			rulesInfo = append(rulesInfo, r)
		}
	}

	options := &vpcv1.CreateLoadBalancerListenerPolicyOptions{
		LoadBalancerID: &lbID,
		ListenerID:     &listenerID,
		Action:         &action,
		Priority:       &priority,
		Target:         target,
		Name:           &name,
		Rules:          rulesInfo,
	}

	isLBKey := "load_balancer_key_" + lbID
	ibmMutexKV.Lock(isLBKey)
	defer ibmMutexKV.Unlock(isLBKey)

	_, err = isWaitForLbAvailable(sess, lbID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf(
			"LB-LP Error checking for load balancer (%s) is active: %s", lbID, err)
	}

	policy, response, err := sess.CreateLoadBalancerListenerPolicy(options)
	if err != nil {
		return fmt.Errorf("Error while creating lb listener policy for LB %s: Error %v Response %v", lbID, err, *response)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", lbID, listenerID, *(policy.ID)))

	_, err = isWaitForLbListenerPolicyAvailable(sess, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}
	return nil
}

func isWaitForLbAvailable(vpc *vpcv1.VpcV1, id string, timeout time.Duration) (interface{}, error) {

	stateConf := &resource.StateChangeConf{
		Pending:    []string{isLBListenerPolicyPending},
		Target:     []string{isLBProvisioningDone},
		Refresh:    isLbRefreshFunc(vpc, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isLbRefreshFunc(vpc *vpcv1.VpcV1, id string) resource.StateRefreshFunc {
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

func isWaitForLbListenerPolicyAvailable(vpc *vpcv1.VpcV1, id string, timeout time.Duration) (interface{}, error) {

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isLBListenerProvisioning, "create_pending", "update_pending", "maintenance_pending"},
		Target:     []string{isLBListenerProvisioningDone},
		Refresh:    isLbListenerPolicyRefreshFunc(vpc, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isLbListenerPolicyRefreshFunc(vpc *vpcv1.VpcV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		parts, err := idParts(id)
		if err != nil {
			return nil, "", err
		}

		lbID := parts[0]
		listenerID := parts[1]
		policyID := parts[2]

		getLbListenerPolicyOptions := &vpcv1.GetLoadBalancerListenerPolicyOptions{
			LoadBalancerID: &lbID,
			ListenerID:     &listenerID,
			ID:             &policyID,
		}

		policy, _, err := vpc.GetLoadBalancerListenerPolicy(getLbListenerPolicyOptions)

		if err != nil {
			return policy, "", err
		}

		if *policy.ProvisioningStatus == isLBListenerPolicyAvailable || *policy.ProvisioningStatus == isLBListenerPolicyFailed {
			return policy, isLBListenerProvisioningDone, nil
		}

		return policy, *policy.ProvisioningStatus, nil
	}
}

func resourceIBMISLBListenerPolicyRead(d *schema.ResourceData, meta interface{}) error {

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

	if userDetails.generation == 1 {
		err := classicLbListenerPolicyGet(d, meta, lbID, listenerID, policyID)
		if err != nil {
			return err
		}
	} else {
		err := lbListenerPolicyGet(d, meta, lbID, listenerID, policyID)
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceIBMISLBListenerPolicyExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return false, err
	}
	ID := d.Id()
	if userDetails.generation == 1 {
		exists, err := classicLbListenerPolicyExists(d, meta, ID)
		return exists, err
	} else {
		exists, err := lbListenerPolicyExists(d, meta, ID)
		return exists, err
	}
}

func classicLbListenerPolicyExists(d *schema.ResourceData, meta interface{}, ID string) (bool, error) {
	sess, err := classicVpcClient(meta)
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

	//populate lblistenerpolicyOPtions
	getLbListenerPolicyOptions := &vpcclassicv1.GetLoadBalancerListenerPolicyOptions{
		LoadBalancerID: &lbID,
		ListenerID:     &listenerID,
		ID:             &policyID,
	}

	//Getting lb listener policy
	_, response, err := sess.GetLoadBalancerListenerPolicy(getLbListenerPolicyOptions)

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("Error getting Load balancer policy: %s\n%s", err, response)
	}

	return true, nil
}

func lbListenerPolicyExists(d *schema.ResourceData, meta interface{}, ID string) (bool, error) {
	sess, err := vpcClient(meta)
	if err != nil {
		return false, err
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return false, err
	}
	if len(parts) != 3 {
		return false, fmt.Errorf("Incorrect ID %s: ID should be a combination of lbID/listenerID/policyID", d.Id())
	}

	lbID := parts[0]
	listenerID := parts[1]
	policyID := parts[2]

	getLbListenerPolicyOptions := &vpcv1.GetLoadBalancerListenerPolicyOptions{
		LoadBalancerID: &lbID,
		ListenerID:     &listenerID,
		ID:             &policyID,
	}

	//Getting lb listener policy
	_, response, err := sess.GetLoadBalancerListenerPolicy(getLbListenerPolicyOptions)

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("Error getting Load balancer policy: %s\n%s", err, response)
	}
	return true, nil
}
func resourceIBMISLBListenerPolicyUpdate(d *schema.ResourceData, meta interface{}) error {

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

	if userDetails.generation == 1 {

		err := classicLbListenerPolicyUpdate(d, meta, lbID, listenerID, policyID)
		if err != nil {
			return err
		}
	} else {

		err := lbListenerPolicyUpdate(d, meta, lbID, listenerID, policyID)
		if err != nil {
			return err
		}
	}

	return resourceIBMISLBListenerPolicyRead(d, meta)
}

func classicLbListenerPolicyUpdate(d *schema.ResourceData, meta interface{}, lbID, listenerID, ID string) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}

	hasChanged := false
	updatePolicyOptions := vpcclassicv1.UpdateLoadBalancerListenerPolicyOptions{}
	updatePolicyOptions.LoadBalancerID = &lbID
	updatePolicyOptions.ListenerID = &listenerID
	updatePolicyOptions.ID = &ID

	loadBalancerListenerPolicyPatchModel := &vpcclassicv1.LoadBalancerListenerPolicyPatch{}

	if d.HasChange(isLBListenerPolicyName) {
		policy := d.Get(isLBListenerPolicyName).(string)
		loadBalancerListenerPolicyPatchModel.Name = &policy
		hasChanged = true
	}

	if d.HasChange(isLBListenerPolicyPriority) {
		prio := d.Get(isLBListenerPolicyPriority).(int)
		priority := int64(prio)
		loadBalancerListenerPolicyPatchModel.Priority = &priority
		hasChanged = true
	}

	var target vpcclassicv1.LoadBalancerListenerPolicyTargetPatchIntf

	//If Action is forward and TargetID is changed, set the target to pool ID
	if d.Get(isLBListenerPolicyAction).(string) == "forward" && d.HasChange(isLBListenerPolicyTargetID) {

		//User can set the poolId as combination of lbID/poolID, if so parse the string & get the poolID
		id, err := getPoolID(d.Get(isLBListenerPolicyTargetID).(string))
		if err != nil {
			return err
		}

		target = &vpcclassicv1.LoadBalancerListenerPolicyTargetPatch{
			ID: &id,
		}

		loadBalancerListenerPolicyPatchModel.Target = target
		hasChanged = true
	} else if d.Get(isLBListenerPolicyAction).(string) == "redirect" {
		//if Action is redirect and either status code or URL chnaged, set accordingly
		//LoadBalancerListenerPolicyPatchTargetLoadBalancerListenerPolicyRedirectURLPatch

		redirectPatch := vpcclassicv1.LoadBalancerListenerPolicyTargetPatchLoadBalancerListenerPolicyRedirectURLPatch{}

		targetChange := false
		if d.HasChange(isLBListenerPolicyTargetHTTPStatusCode) {
			status := d.Get(isLBListenerPolicyTargetHTTPStatusCode).(int)
			sc := int64(status)
			redirectPatch.HTTPStatusCode = &sc
			hasChanged = true
			targetChange = true
		}

		if d.HasChange(isLBListenerPolicyTargetURL) {
			url := d.Get(isLBListenerPolicyTargetURL).(string)
			redirectPatch.URL = &url
			hasChanged = true
			targetChange = true
		}

		//Update the target only if there is a change in either statusCode or URL
		if targetChange {
			target = &redirectPatch
			loadBalancerListenerPolicyPatchModel.Target = target
		}
	}

	isLBKey := "load_balancer_key_" + lbID
	ibmMutexKV.Lock(isLBKey)
	defer ibmMutexKV.Unlock(isLBKey)

	if hasChanged {
		loadBalancerListenerPolicyPatch, err := loadBalancerListenerPolicyPatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("Error calling asPatch for LoadBalancerListenerPolicyPatch: %s", err)
		}
		updatePolicyOptions.LoadBalancerListenerPolicyPatch = loadBalancerListenerPolicyPatch
		_, err = isWaitForClassicLbAvailable(sess, lbID, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return fmt.Errorf(
				"LB-LP Error checking for load balancer (%s) is active: %s", lbID, err)
		}
		_, response, err := sess.UpdateLoadBalancerListenerPolicy(&updatePolicyOptions)
		if err != nil {
			return fmt.Errorf("Error Getting Instance: %s\n%s", err, response)
		}

		_, err = isWaitForClassicLbListenerPolicyAvailable(sess, d.Id(), d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return err
		}
	}
	return nil
}

func lbListenerPolicyUpdate(d *schema.ResourceData, meta interface{}, lbID, listenerID, ID string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	hasChanged := false
	updatePolicyOptions := vpcv1.UpdateLoadBalancerListenerPolicyOptions{}
	updatePolicyOptions.LoadBalancerID = &lbID
	updatePolicyOptions.ListenerID = &listenerID
	updatePolicyOptions.ID = &ID

	loadBalancerListenerPolicyPatchModel := &vpcv1.LoadBalancerListenerPolicyPatch{}

	if d.HasChange(isLBListenerPolicyName) {
		policy := d.Get(isLBListenerPolicyName).(string)
		loadBalancerListenerPolicyPatchModel.Name = &policy
		hasChanged = true
	}

	if d.HasChange(isLBListenerPolicyPriority) {
		prio := d.Get(isLBListenerPolicyPriority).(int)
		priority := int64(prio)
		loadBalancerListenerPolicyPatchModel.Priority = &priority
		hasChanged = true
	}

	var target vpcv1.LoadBalancerListenerPolicyTargetPatchIntf
	//If Action is forward and TargetID is changed, set the target to pool ID
	if d.Get(isLBListenerPolicyAction).(string) == "forward" && d.HasChange(isLBListenerPolicyTargetID) {

		//User can set the poolId as combination of lbID/poolID, if so parse the string & get the poolID
		id, err := getPoolID(d.Get(isLBListenerPolicyTargetID).(string))
		if err != nil {
			return err
		}
		target = &vpcv1.LoadBalancerListenerPolicyTargetPatchLoadBalancerPoolIdentity{
			ID: &id,
		}

		loadBalancerListenerPolicyPatchModel.Target = target
		hasChanged = true
	} else if d.Get(isLBListenerPolicyAction).(string) == "redirect" {
		//if Action is redirect and either status code or URL chnaged, set accordingly
		//LoadBalancerListenerPolicyPatchTargetLoadBalancerListenerPolicyRedirectURLPatch

		redirectPatch := vpcv1.LoadBalancerListenerPolicyTargetPatchLoadBalancerListenerPolicyRedirectURLPatch{}

		targetChange := false
		if d.HasChange(isLBListenerPolicyTargetHTTPStatusCode) {
			status := d.Get(isLBListenerPolicyTargetHTTPStatusCode).(int)
			sc := int64(status)
			redirectPatch.HTTPStatusCode = &sc
			hasChanged = true
			targetChange = true
		}

		if d.HasChange(isLBListenerPolicyTargetURL) {
			url := d.Get(isLBListenerPolicyTargetURL).(string)
			redirectPatch.URL = &url
			hasChanged = true
			targetChange = true
		}

		//Update the target only if there is a change in either statusCode or URL
		if targetChange {
			target = &redirectPatch
			loadBalancerListenerPolicyPatchModel.Target = target
		}
	}

	if hasChanged {
		loadBalancerListenerPolicyPatch, err := loadBalancerListenerPolicyPatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("Error calling asPatch for LoadBalancerListenerPolicyPatch: %s", err)
		}
		updatePolicyOptions.LoadBalancerListenerPolicyPatch = loadBalancerListenerPolicyPatch
		isLBKey := "load_balancer_key_" + lbID
		ibmMutexKV.Lock(isLBKey)
		defer ibmMutexKV.Unlock(isLBKey)

		_, err = isWaitForLbAvailable(sess, lbID, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return fmt.Errorf(
				"LB-LP Error checking for load balancer (%s) is active: %s", lbID, err)
		}
		_, response, err := sess.UpdateLoadBalancerListenerPolicy(&updatePolicyOptions)
		if err != nil {
			return fmt.Errorf("Error Updating in policy : %s\n%s", err, response)
		}

		_, err = isWaitForLbListenerPolicyAvailable(sess, d.Id(), d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceIBMISLBListenerPolicyDelete(d *schema.ResourceData, meta interface{}) error {

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

	isLBKey := "load_balancer_key_" + lbID
	ibmMutexKV.Lock(isLBKey)
	defer ibmMutexKV.Unlock(isLBKey)

	if userDetails.generation == 1 {
		err := classicLbListenerPolicycDelete(d, meta, lbID, listenerID, policyID)
		if err != nil {
			return err
		}
	} else {
		err := lbListenerPolicyDelete(d, meta, lbID, listenerID, policyID)
		if err != nil {
			return err
		}
	}
	d.SetId("")
	return nil

}

func classicLbListenerPolicycDelete(d *schema.ResourceData, meta interface{}, lbID, listenerID, ID string) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}

	//Getting policy optins
	getLbListenerPolicyOptions := &vpcclassicv1.GetLoadBalancerListenerPolicyOptions{
		LoadBalancerID: &lbID,
		ListenerID:     &listenerID,
		ID:             &ID,
	}

	//Getting lb listener policy
	_, response, err := sess.GetLoadBalancerListenerPolicy(getLbListenerPolicyOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error in classicLbListenerPolicyGet : %s\n%s", err, response)
	}

	deleteLbListenerPolicyOptions := &vpcclassicv1.DeleteLoadBalancerListenerPolicyOptions{
		LoadBalancerID: &lbID,
		ListenerID:     &listenerID,
		ID:             &ID,
	}

	_, err = isWaitForClassicLbAvailable(sess, lbID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf(
			"LB-LP Error checking for load balancer (%s) is active: %s", lbID, err)
	}

	response, err = sess.DeleteLoadBalancerListenerPolicy(deleteLbListenerPolicyOptions)
	if err != nil {
		return fmt.Errorf("Error in classicLbListenerPolicycDelete: %s\n%s", err, response)
	}
	_, err = isWaitForLbListenerPolicyClassicDeleted(sess, d.Id(), d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}
	return nil
}

func lbListenerPolicyDelete(d *schema.ResourceData, meta interface{}, lbID, listenerID, ID string) error {

	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	//Getting policy optins
	getLbListenerPolicyOptions := &vpcv1.GetLoadBalancerListenerPolicyOptions{
		LoadBalancerID: &lbID,
		ListenerID:     &listenerID,
		ID:             &ID,
	}

	//Getting lb listener policy
	_, response, err := sess.GetLoadBalancerListenerPolicy(getLbListenerPolicyOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
	}

	deleteLbListenerPolicyOptions := &vpcv1.DeleteLoadBalancerListenerPolicyOptions{
		LoadBalancerID: &lbID,
		ListenerID:     &listenerID,
		ID:             &ID,
	}

	_, err = isWaitForLbAvailable(sess, lbID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf(
			"LB-LP Error checking for load balancer (%s) is active: %s", lbID, err)
	}

	response, err = sess.DeleteLoadBalancerListenerPolicy(deleteLbListenerPolicyOptions)
	if err != nil {
		return fmt.Errorf("Error in lbListenerPolicyDelete: %s\n%s", err, response)
	}
	_, err = isWaitForLbListnerPolicyDeleted(sess, d.Id(), d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}
	return nil
}
func isWaitForLbListnerPolicyDeleted(vpc *vpcv1.VpcV1, id string, timeout time.Duration) (interface{}, error) {

	stateConf := &resource.StateChangeConf{
		Pending:    []string{isLBListenerPolicyRetry, isLBListenerPolicyDeleting},
		Target:     []string{isLBListenerPolicyFailed, isLBListenerPolicyDeleted},
		Refresh:    isLbListenerPolicyDeleteRefreshFunc(vpc, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isLbListenerPolicyDeleteRefreshFunc(vpc *vpcv1.VpcV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		//Retrieve lbId, listenerId and policyID
		parts, err := idParts(id)
		if err != nil {
			return nil, isLBListenerPolicyFailed, nil
		}

		lbID := parts[0]
		listenerID := parts[1]
		policyID := parts[2]

		getLbListenerPolicyOptions := &vpcv1.GetLoadBalancerListenerPolicyOptions{
			LoadBalancerID: &lbID,
			ListenerID:     &listenerID,
			ID:             &policyID,
		}

		//Getting lb listener policy
		policy, response, err := vpc.GetLoadBalancerListenerPolicy(getLbListenerPolicyOptions)

		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return policy, isLBListenerPolicyDeleted, nil
			}
			return nil, isLBListenerPolicyFailed, err
		}
		return policy, isLBListenerPolicyDeleting, err
	}
}

func classicLbListenerPolicyGet(d *schema.ResourceData, meta interface{}, lbID, listenerID, id string) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}

	//Getting policy optins
	getLbListenerPolicyOptions := &vpcclassicv1.GetLoadBalancerListenerPolicyOptions{
		LoadBalancerID: &lbID,
		ListenerID:     &listenerID,
		ID:             &id,
	}

	//Getting lb listener policy
	policy, response, err := sess.GetLoadBalancerListenerPolicy(getLbListenerPolicyOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error in classicLbListenerPolicyGet : %s\n%s", err, response)
	}

	d.Set(isLBListenerPolicyLBID, lbID)
	d.Set(isLBListenerPolicyListenerID, listenerID)
	d.Set(isLBListenerPolicyAction, policy.Action)
	d.Set(isLBListenerPolicyID, id)
	d.Set(isLBListenerPolicyPriority, policy.Priority)
	d.Set(isLBListenerPolicyName, policy.Name)
	d.Set(isLBListenerPolicyStatus, policy.ProvisioningStatus)

	if policy.Rules != nil {
		rulesSet := make([]interface{}, 0)
		for _, rule := range policy.Rules {
			getLbListenerPolicyRulesOptions := &vpcclassicv1.GetLoadBalancerListenerPolicyRuleOptions{
				LoadBalancerID: &lbID,
				ListenerID:     &listenerID,
				ID:             rule.ID,
				PolicyID:       &id,
			}
			ruleInfo, response, err := sess.GetLoadBalancerListenerPolicyRule(getLbListenerPolicyRulesOptions)
			if err != nil {
				return fmt.Errorf("Error in classicLbListenerPolicyGet rule: %s\n%s", err, response)
			}

			r := map[string]interface{}{
				isLBListenerPolicyRuleID:        *ruleInfo.ID,
				isLBListenerPolicyRuleCondition: *ruleInfo.Condition,
				isLBListenerPolicyRuleType:      *ruleInfo.Type,
				isLBListenerPolicyRuleField:     *ruleInfo.Field,
				isLBListenerPolicyRuleValue:     *ruleInfo.Value,
			}
			rulesSet = append(rulesSet, r)
		}
		d.Set(isLBListenerPolicyRulesInfo, rulesSet)
	}

	// `LoadBalancerPoolReference` is in the response if `action` is `forward`.
	// `LoadBalancerListenerPolicyRedirectURL` is in the response if `action` is `redirect`.

	if *(policy.Action) == "forward" {
		if reflect.TypeOf(policy.Target).String() == "*vpcclassicv1.LoadBalancerListenerPolicyTargetLoadBalancerPoolReference" {
			target, ok := policy.Target.(*vpcclassicv1.LoadBalancerListenerPolicyTargetLoadBalancerPoolReference)
			if ok {
				d.Set(isLBListenerPolicyTargetID, target.ID)
			}
		}

	} else if *(policy.Action) == "redirect" {
		if reflect.TypeOf(policy.Target).String() == "*vpcclassicv1.LoadBalancerListenerPolicyTargetLoadBalancerListenerPolicyRedirectURL" {
			target, ok := policy.Target.(*vpcclassicv1.LoadBalancerListenerPolicyTargetLoadBalancerListenerPolicyRedirectURL)
			if ok {
				d.Set(isLBListenerPolicyTargetURL, target.URL)
				d.Set(isLBListenerPolicyTargetHTTPStatusCode, target.HTTPStatusCode)
			}
		}
	}

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

func lbListenerPolicyGet(d *schema.ResourceData, meta interface{}, lbID, listenerID, id string) error {

	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	//Getting policy optins
	getLbListenerPolicyOptions := &vpcv1.GetLoadBalancerListenerPolicyOptions{
		LoadBalancerID: &lbID,
		ListenerID:     &listenerID,
		ID:             &id,
	}

	//Getting lb listener policy
	policy, response, err := sess.GetLoadBalancerListenerPolicy(getLbListenerPolicyOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return err
	}

	//set the argument values
	d.Set(isLBListenerPolicyLBID, lbID)
	d.Set(isLBListenerPolicyListenerID, listenerID)
	d.Set(isLBListenerPolicyAction, policy.Action)
	d.Set(isLBListenerPolicyID, id)
	d.Set(isLBListenerPolicyPriority, policy.Priority)
	d.Set(isLBListenerPolicyName, policy.Name)
	d.Set(isLBListenerPolicyStatus, policy.ProvisioningStatus)

	//set rules - Doubt (Rules has condition, type, value, field and id where as SDK has only Href and id, so setting only id)
	if policy.Rules != nil {
		policyRules := policy.Rules
		rulesInfo := make([]map[string]interface{}, 0)
		for _, index := range policyRules {

			l := map[string]interface{}{
				isLBListenerPolicyRuleID: index.ID,
			}
			rulesInfo = append(rulesInfo, l)
		}
		d.Set(isLBListenerPolicyRules, rulesInfo)
	}

	// `LoadBalancerPoolReference` is in the response if `action` is `forward`.
	// `LoadBalancerListenerPolicyRedirectURL` is in the response if `action` is `redirect`.

	if *(policy.Action) == "forward" {
		if reflect.TypeOf(policy.Target).String() == "*vpcv1.LoadBalancerListenerPolicyTargetLoadBalancerPoolReference" {
			target, ok := policy.Target.(*vpcv1.LoadBalancerListenerPolicyTargetLoadBalancerPoolReference)
			if ok {
				d.Set(isLBListenerPolicyTargetID, target.ID)
			}
		}

	} else if *(policy.Action) == "redirect" {
		if reflect.TypeOf(policy.Target).String() == "*vpcv1.LoadBalancerListenerPolicyTargetLoadBalancerListenerPolicyRedirectURL" {
			target, ok := policy.Target.(*vpcv1.LoadBalancerListenerPolicyTargetLoadBalancerListenerPolicyRedirectURL)
			if ok {
				d.Set(isLBListenerPolicyTargetURL, target.URL)
				d.Set(isLBListenerPolicyTargetHTTPStatusCode, target.HTTPStatusCode)
			}
		}
	}

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

func isWaitForLbListenerPolicyClassicDeleted(vpc *vpcclassicv1.VpcClassicV1, id string, timeout time.Duration) (interface{}, error) {

	stateConf := &resource.StateChangeConf{
		Pending:    []string{isLBListenerPolicyRetry, isLBListenerPolicyDeleting, "delete_pending"},
		Target:     []string{isLBListenerPolicyFailed, isLBListenerPolicyDeleted},
		Refresh:    isLbListenerPolicyClassicDeleteRefreshFunc(vpc, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isLbListenerPolicyClassicDeleteRefreshFunc(vpc *vpcclassicv1.VpcClassicV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		//Retrieve lbId and listenerId
		parts, err := idParts(id)
		if err != nil {
			return nil, isLBListenerPolicyFailed, nil
		}

		lbID := parts[0]
		listenerID := parts[1]
		policyID := parts[2]

		getLbListenerPolicyOptions := &vpcclassicv1.GetLoadBalancerListenerPolicyOptions{
			LoadBalancerID: &lbID,
			ListenerID:     &listenerID,
			ID:             &policyID,
		}

		//Getting lb listener policy
		policy, response, err := vpc.GetLoadBalancerListenerPolicy(getLbListenerPolicyOptions)

		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return policy, isLBListenerPolicyDeleted, nil
			}

			return nil, isLBListenerPolicyFailed, err
		}

		return policy, isLBListenerPolicyDeleting, err
	}
}
