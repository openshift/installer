// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isLBListenerPolicyLBID                    = "lb"
	isLBListenerPolicyListenerID              = "listener"
	isLBListenerPolicyAction                  = "action"
	isLBListenerPolicyPriority                = "priority"
	isLBListenerPolicyName                    = "name"
	isLBListenerPolicyID                      = "policy_id"
	isLBListenerPolicyRules                   = "rules"
	isLBListenerPolicyRulesInfo               = "rule_info"
	isLBListenerPolicyTargetID                = "target_id"
	isLBListenerPolicyTargetHTTPStatusCode    = "target_http_status_code"
	isLBListenerPolicyTargetURL               = "target_url"
	isLBListenerPolicyStatus                  = "provisioning_status"
	isLBListenerPolicyRuleID                  = "rule_id"
	isLBListenerPolicyAvailable               = "active"
	isLBListenerPolicyFailed                  = "failed"
	isLBListenerPolicyPending                 = "pending"
	isLBListenerPolicyDeleting                = "deleting"
	isLBListenerPolicyDeleted                 = "done"
	isLBListenerPolicyRetry                   = "retry"
	isLBListenerPolicyRuleCondition           = "condition"
	isLBListenerPolicyRuleType                = "type"
	isLBListenerPolicyRuleValue               = "value"
	isLBListenerPolicyRuleField               = "field"
	isLBListenerPolicyProvisioning            = "provisioning"
	isLBListenerPolicyProvisioningDone        = "done"
	isLBListenerPolicyHTTPSRedirectStatusCode = "target_https_redirect_status_code"
	isLBListenerPolicyHTTPSRedirectURI        = "target_https_redirect_uri"
	isLBListenerPolicyHTTPSRedirectListener   = "target_https_redirect_listener"
)

func ResourceIBMISLBListenerPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMISLBListenerPolicyCreate,
		ReadContext:   resourceIBMISLBListenerPolicyRead,
		UpdateContext: resourceIBMISLBListenerPolicyUpdate,
		DeleteContext: resourceIBMISLBListenerPolicyDelete,
		Exists:        resourceIBMISLBListenerPolicyExists,
		Importer:      &schema.ResourceImporter{},

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

			isLBListenerPolicyHTTPSRedirectStatusCode: {
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{isLBListenerPolicyHTTPSRedirectListener},
				Deprecated:   "Please use the argument 'target'",
				Description:  "The HTTP status code to be returned in the redirect response",
			},

			isLBListenerPolicyHTTPSRedirectURI: {
				Type:         schema.TypeString,
				Optional:     true,
				Deprecated:   "Please use the argument 'target'",
				RequiredWith: []string{isLBListenerPolicyHTTPSRedirectListener, isLBListenerPolicyHTTPSRedirectStatusCode},
				Description:  "Target URI where traffic will be redirected",
			},

			isLBListenerPolicyHTTPSRedirectListener: {
				Type:         schema.TypeString,
				Optional:     true,
				Deprecated:   "Please use the argument 'target'",
				RequiredWith: []string{isLBListenerPolicyHTTPSRedirectStatusCode},
				Description:  "ID of the listener that will be set as http redirect target",
			},

			isLBListenerPolicyAction: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_lb_listener_policy", isLBListenerPolicyAction),
				Description:  "Policy Action",
			},

			isLBListenerPolicyPriority: {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validate.ValidateLBListenerPolicyPriority,
				Description:  "Listener Policy Priority",
			},

			isLBListenerPolicyName: {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     false,
				Computed:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_lb_listener_policy", isLBListenerPolicyName),
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
							ValidateFunc: validate.InvokeValidator("ibm_is_lb_listener_policy_rule", isLBListenerPolicyRulecondition),
							Description:  "Condition of the rule",
						},

						isLBListenerPolicyRuleType: {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.InvokeValidator("ibm_is_lb_listener_policy_rule", isLBListenerPolicyRuleType),
							Description:  "Type of the rule",
						},

						isLBListenerPolicyRuleValue: {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.ValidateStringLength,
							Description:  "Value to be matched for rule condition",
						},

						isLBListenerPolicyRuleField: {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validate.ValidateStringLength,
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
				Type:       schema.TypeString,
				ForceNew:   false,
				Optional:   true,
				Deprecated: "Please use the argument 'target'",
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
				Deprecated:  "Please use the argument 'target'",
				Description: "Listener Policy target HTTPS Status code.",
			},

			isLBListenerPolicyTargetURL: {
				Type:        schema.TypeString,
				ForceNew:    false,
				Optional:    true,
				Deprecated:  "Please use the argument 'target'",
				Description: "Policy Target URL",
			},
			"target": &schema.Schema{
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: []string{"target_url", "target_http_status_code", "target_id", "target_https_redirect_listener", "target_https_redirect_uri", "target_https_redirect_status_code"},
				Description:   "- If `action` is `forward`, the response is a `LoadBalancerPoolReference`- If `action` is `redirect`, the response is a `LoadBalancerListenerPolicyRedirectURL`- If `action` is `https_redirect`, the response is a `LoadBalancerListenerHTTPSRedirect`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"deleted": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The pool's canonical URL.",
						},
						"id": &schema.Schema{
							Type:          schema.TypeString,
							Optional:      true,
							AtLeastOneOf:  []string{"target.0.id", "target.0.http_status_code", "target.0.url", "target.0.listener"},
							ConflictsWith: []string{"target.0.http_status_code", "target.0.url", "target.0.listener", "target.0.uri"},
							Description:   "The unique identifier for this load balancer pool.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this load balancer pool. The name is unique across all pools for the load balancer.",
						},
						"http_status_code": &schema.Schema{
							Type:          schema.TypeInt,
							Optional:      true,
							AtLeastOneOf:  []string{"target.0.id", "target.0.http_status_code", "target.0.url", "target.0.listener"},
							ConflictsWith: []string{"target.0.id"},
							Description:   "The HTTP status code for this redirect.",
						},
						"url": &schema.Schema{
							Type:          schema.TypeString,
							Optional:      true,
							AtLeastOneOf:  []string{"target.0.id", "target.0.http_status_code", "target.0.url", "target.0.listener"},
							ConflictsWith: []string{"target.0.id", "target.0.listener", "target.0.uri"},
							Description:   "The redirect target URL.",
						},
						"listener": &schema.Schema{
							Type:          schema.TypeList,
							MaxItems:      1,
							Optional:      true,
							AtLeastOneOf:  []string{"target.0.id", "target.0.http_status_code", "target.0.url", "target.0.listener"},
							ConflictsWith: []string{"target.0.id", "target.0.url"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"deleted": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"more_info": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
									},
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The listener's canonical URL.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "The unique identifier for this load balancer listener.",
									},
								},
							},
						},
						"uri": &schema.Schema{
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"target.0.id", "target.0.url"},
							Description:   "The redirect relative target URI.",
						},
					},
				},
			},
			isLBListenerPolicyStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Listner Policy status",
			},

			flex.RelatedCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the LB resource",
			},
		},
	}
}

func ResourceIBMISLBListenerPolicyValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)
	action := "forward, redirect, reject, https_redirect"
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isLBListenerPolicyName,
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isLBListenerPolicyAction,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              action})

	ibmISLBListenerPolicyResourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_lb_listener_policy", Schema: validateSchema}
	return &ibmISLBListenerPolicyResourceValidator
}

func resourceIBMISLBListenerPolicyCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	//Get the Load balancer ID
	lbID := d.Get(isLBListenerPolicyLBID).(string)

	//User can set listener id as combination of lbID/listenerID, parse and get the listenerID
	listenerID, err := getListenerID(d.Get(isLBListenerPolicyListenerID).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	action := d.Get(isLBListenerPolicyAction).(string)
	priority := int64(d.Get(isLBListenerPolicyPriority).(int))

	//user-defined name for this policy.
	var name string
	if n, ok := d.GetOk(isLBListenerPolicyName); ok {
		name = n.(string)
	}

	errDiag := lbListenerPolicyCreate(d, meta, lbID, listenerID, action, name, priority)
	if errDiag != nil {
		return errDiag
	}

	return resourceIBMISLBListenerPolicyRead(context, d, meta)
}

func getListenerID(id string) (string, error) {
	if strings.Contains(id, "/") {
		parts, err := flex.IdParts(id)
		if err != nil {
			return "", err
		}

		return parts[1], nil
	} else {
		return id, nil
	}
}

func getPoolID(id string) (string, error) {
	if strings.Contains(id, "/") {
		parts, err := flex.IdParts(id)
		if err != nil {
			return "", err
		}

		return parts[1], nil
	}
	return id, nil

}

func lbListenerPolicyCreate(d *schema.ResourceData, meta interface{}, lbID, listenerID, action, name string, priority int64) diag.Diagnostics {

	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
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

	listener, listenerSet := d.GetOk(isLBListenerPolicyHTTPSRedirectListener)
	httpsStatusCode, httpsStatusSet := d.GetOk(isLBListenerPolicyHTTPSRedirectStatusCode)
	uri, uriSet := d.GetOk(isLBListenerPolicyHTTPSRedirectURI)
	if _, ok := d.GetOk("target"); ok {
		target, err = resourceIBMIsLbListenerPolicyMapToLoadBalancerListenerPolicyTargetPrototype(d.Get("target.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}

	} else {
		if actionChk.(string) == "forward" {
			if targetIDSet {

				//User can set the poolId as combination of lbID/poolID, if so parse the string & get the poolID
				id, err := getPoolID(tID.(string))
				if err != nil {
					return diag.FromErr(err)
				}

				target = &vpcv1.LoadBalancerListenerPolicyTargetPrototypeLoadBalancerPoolIdentity{
					ID: &id,
				}
			} else {
				return diag.FromErr(fmt.Errorf("When action is forward please specify target_id"))
			}
		} else if actionChk.(string) == "redirect" {

			urlPrototype := vpcv1.LoadBalancerListenerPolicyTargetPrototypeLoadBalancerListenerPolicyRedirectURLPrototype{}

			if statusSet {
				sc := int64(statusCode.(int))
				urlPrototype.HTTPStatusCode = &sc
			} else {
				return diag.FromErr(fmt.Errorf("When action is redirect please specify target_http_status_code"))
			}

			if urlSet {
				link := url.(string)
				urlPrototype.URL = &link
			} else {
				return diag.FromErr(fmt.Errorf("When action is redirect please specify target_url"))
			}

			target = &urlPrototype
		} else if actionChk.(string) == "https_redirect" {

			urlPrototype := vpcv1.LoadBalancerListenerPolicyTargetPrototypeLoadBalancerListenerPolicyHTTPSRedirectPrototype{}

			if listenerSet {
				listener := listener.(string)
				urlPrototype.Listener = &vpcv1.LoadBalancerListenerIdentity{
					ID: &listener,
				}
			} else {
				return diag.FromErr(fmt.Errorf("When action is https_redirect please specify target_https_redirect_listener"))
			}

			if httpsStatusSet {
				sc := int64(httpsStatusCode.(int))
				urlPrototype.HTTPStatusCode = &sc
			} else {
				return diag.FromErr(fmt.Errorf("When action is https_redirect please specify target_https_redirect_status_code"))
			}

			if uriSet {
				link := uri.(string)
				urlPrototype.URI = &link
			}

			target = &urlPrototype
		}
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
	conns.IbmMutexKV.Lock(isLBKey)
	defer conns.IbmMutexKV.Unlock(isLBKey)

	_, err = isWaitForLbAvailable(sess, lbID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(fmt.Errorf(
			"LB-LP Error checking for load balancer (%s) is active: %s", lbID, err))
	}

	policy, response, err := sess.CreateLoadBalancerListenerPolicy(options)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error while creating lb listener policy for LB %s: Error %v Response %v", lbID, err, *response))
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", lbID, listenerID, *(policy.ID)))

	_, err = isWaitForLbListenerPolicyAvailable(sess, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
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

		parts, err := flex.IdParts(id)
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

func resourceIBMISLBListenerPolicyRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	ID := d.Id()
	parts, err := flex.IdParts(ID)
	if err != nil {
		return diag.FromErr(err)
	}

	lbID := parts[0]
	listenerID := parts[1]
	policyID := parts[2]

	diag := lbListenerPolicyGet(d, meta, lbID, listenerID, policyID)
	if diag != nil {
		return diag
	}

	return nil
}

func resourceIBMISLBListenerPolicyExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	ID := d.Id()

	exists, err := lbListenerPolicyExists(d, meta, ID)
	return exists, err

}

func lbListenerPolicyExists(d *schema.ResourceData, meta interface{}, ID string) (bool, error) {
	sess, err := vpcClient(meta)
	if err != nil {
		return false, err
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return false, err
	}
	if len(parts) != 3 {
		return false, fmt.Errorf("[ERROR] Incorrect ID %s: ID should be a combination of lbID/listenerID/policyID", d.Id())
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
		return false, fmt.Errorf("[ERROR] Error getting Load balancer policy: %s\n%s", err, response)
	}
	return true, nil
}
func resourceIBMISLBListenerPolicyUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	lbID := parts[0]
	listenerID := parts[1]
	policyID := parts[2]

	diagErr := lbListenerPolicyUpdate(d, meta, lbID, listenerID, policyID)
	if diagErr != nil {
		return diagErr
	}

	return resourceIBMISLBListenerPolicyRead(context, d, meta)
}

func lbListenerPolicyUpdate(d *schema.ResourceData, meta interface{}, lbID, listenerID, ID string) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
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
	httpsURIRemoved := false
	var target vpcv1.LoadBalancerListenerPolicyTargetPatchIntf
	if d.HasChange("target") {
		target, err := resourceIBMIsLbListenerPolicyMapToLoadBalancerListenerPolicyTargetPatch(d, d.Get("target.0").(map[string]interface{}), &httpsURIRemoved)
		if err != nil {
			return diag.FromErr(err)
		}
		loadBalancerListenerPolicyPatchModel.Target = target
		hasChanged = true

	} else {

		//If Action is forward and TargetID is changed, set the target to pool ID
		if d.Get(isLBListenerPolicyAction).(string) == "forward" && d.HasChange(isLBListenerPolicyTargetID) {

			//User can set the poolId as combination of lbID/poolID, if so parse the string & get the poolID
			id, err := getPoolID(d.Get(isLBListenerPolicyTargetID).(string))
			if err != nil {
				return diag.FromErr(err)
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
		} else if d.Get(isLBListenerPolicyAction).(string) == "https_redirect" {

			httpsRedirectPatch := vpcv1.LoadBalancerListenerPolicyTargetPatchLoadBalancerListenerPolicyHTTPSRedirectPatch{}

			targetChange := false
			if d.HasChange(isLBListenerPolicyHTTPSRedirectListener) {
				listener := d.Get(isLBListenerPolicyHTTPSRedirectListener).(string)
				httpsRedirectPatch.Listener = &vpcv1.LoadBalancerListenerIdentity{
					ID: &listener,
				}
				hasChanged = true
				targetChange = true
			}

			if d.HasChange(isLBListenerPolicyHTTPSRedirectStatusCode) {
				status := d.Get(isLBListenerPolicyHTTPSRedirectStatusCode).(int)
				sc := int64(status)
				httpsRedirectPatch.HTTPStatusCode = &sc
				hasChanged = true
				targetChange = true
			}

			if d.HasChange(isLBListenerPolicyHTTPSRedirectURI) {
				uri := d.Get(isLBListenerPolicyHTTPSRedirectURI).(string)
				httpsRedirectPatch.URI = &uri
				hasChanged = true
				targetChange = true
				if uri == "" {
					httpsURIRemoved = true
				}
			}

			//Update the target only if there is a change in either listener, statusCode or URI
			if targetChange {
				target = &httpsRedirectPatch
				loadBalancerListenerPolicyPatchModel.Target = target
			}
		}
	}
	if hasChanged {
		loadBalancerListenerPolicyPatch, err := loadBalancerListenerPolicyPatchModel.AsPatch()
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error calling asPatch for LoadBalancerListenerPolicyPatch: %s", err))
		}
		if httpsURIRemoved {
			loadBalancerListenerPolicyPatch["target"].(map[string]interface{})["uri"] = nil
		}
		updatePolicyOptions.LoadBalancerListenerPolicyPatch = loadBalancerListenerPolicyPatch
		isLBKey := "load_balancer_key_" + lbID
		conns.IbmMutexKV.Lock(isLBKey)
		defer conns.IbmMutexKV.Unlock(isLBKey)

		_, err = isWaitForLbAvailable(sess, lbID, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(fmt.Errorf(
				"LB-LP Error checking for load balancer (%s) is active: %s", lbID, err))
		}
		_, response, err := sess.UpdateLoadBalancerListenerPolicy(&updatePolicyOptions)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error Updating in policy : %s\n%s", err, response))
		}

		_, err = isWaitForLbListenerPolicyAvailable(sess, d.Id(), d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

func resourceIBMISLBListenerPolicyDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	//Retrieve lbId, listenerId and policyID
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	lbID := parts[0]
	listenerID := parts[1]
	policyID := parts[2]

	isLBKey := "load_balancer_key_" + lbID
	conns.IbmMutexKV.Lock(isLBKey)
	defer conns.IbmMutexKV.Unlock(isLBKey)

	err = lbListenerPolicyDelete(d, meta, lbID, listenerID, policyID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
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
		return fmt.Errorf("[ERROR] Error in lbListenerPolicyDelete: %s\n%s", err, response)
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
		parts, err := flex.IdParts(id)
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

func lbListenerPolicyGet(d *schema.ResourceData, meta interface{}, lbID, listenerID, id string) diag.Diagnostics {

	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
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
		return diag.FromErr(err)
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
		for _, policyRuleItem := range policyRules {
			ruleId := *policyRuleItem.ID
			getLbListenerPolicyRuleOptions := &vpcv1.GetLoadBalancerListenerPolicyRuleOptions{
				LoadBalancerID: &lbID,
				ListenerID:     &listenerID,
				PolicyID:       &id,
				ID:             &ruleId,
			}

			//Getting lb listener policy rule
			rule, response, err := sess.GetLoadBalancerListenerPolicyRule(getLbListenerPolicyRuleOptions)
			if err != nil {
				if response != nil && response.StatusCode == 404 {
					d.SetId("")
					return nil
				}
				return diag.FromErr(err)
			}

			l := map[string]interface{}{
				isLBListenerPolicyRuleCondition: rule.Condition,
				isLBListenerPolicyRuleType:      rule.Type,
				isLBListenerPolicyRuleField:     rule.Field,
				isLBListenerPolicyRuleValue:     rule.Value,
				isLBListenerPolicyRuleID:        rule.ID,
			}
			rulesInfo = append(rulesInfo, l)
		}
		d.Set(isLBListenerPolicyRules, rulesInfo)
	}

	// `LoadBalancerPoolReference` is in the response if `action` is `forward`.
	// `LoadBalancerListenerPolicyRedirectURL` is in the response if `action` is `redirect`.

	if !core.IsNil(policy.Target) {
		if _, ok := d.GetOk("target"); ok {

			targetMap, err := resourceIBMIsLbListenerPolicyLoadBalancerListenerPolicyTargetToMap(policy.Target)
			if err != nil {
				return diag.FromErr(err)
			}
			if err = d.Set("target", []map[string]interface{}{targetMap}); err != nil {
				return diag.FromErr(fmt.Errorf("Error setting target: %s", err))
			}
		} else {
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
			} else if *(policy.Action) == "https_redirect" {
				if reflect.TypeOf(policy.Target).String() == "*vpcv1.LoadBalancerListenerPolicyTargetLoadBalancerListenerPolicyHTTPSRedirect" {
					target, ok := policy.Target.(*vpcv1.LoadBalancerListenerPolicyTargetLoadBalancerListenerPolicyHTTPSRedirect)
					if ok {
						d.Set(isLBListenerPolicyHTTPSRedirectListener, target.Listener.ID)
						d.Set(isLBListenerPolicyHTTPSRedirectStatusCode, target.HTTPStatusCode)
						d.Set(isLBListenerPolicyHTTPSRedirectURI, target.URI)
					}
				}
			}
		}
	}
	getLoadBalancerOptions := &vpcv1.GetLoadBalancerOptions{
		ID: &lbID,
	}
	lb, response, err := sess.GetLoadBalancer(getLoadBalancerOptions)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error Getting Load Balancer : %s\n%s", err, response))
	}
	d.Set(flex.RelatedCRN, *lb.CRN)

	return nil
}

func resourceIBMIsLbListenerPolicyMapToLoadBalancerListenerPolicyTargetPrototype(modelMap map[string]interface{}) (vpcv1.LoadBalancerListenerPolicyTargetPrototypeIntf, error) {
	model := &vpcv1.LoadBalancerListenerPolicyTargetPrototype{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["http_status_code"] != nil && modelMap["http_status_code"].(int) != 0 {
		model.HTTPStatusCode = core.Int64Ptr(int64(modelMap["http_status_code"].(int)))
	}
	if modelMap["url"] != nil && modelMap["url"].(string) != "" {
		model.URL = core.StringPtr(modelMap["url"].(string))
	}
	if modelMap["listener"] != nil && len(modelMap["listener"].([]interface{})) > 0 {
		ListenerModel, err := resourceIBMIsLbListenerPolicyMapToLoadBalancerListenerIdentity(modelMap["listener"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Listener = ListenerModel
	}
	if modelMap["uri"] != nil && modelMap["uri"].(string) != "" {
		model.URI = core.StringPtr(modelMap["uri"].(string))
	}
	return model, nil
}
func resourceIBMIsLbListenerPolicyMapToLoadBalancerListenerIdentity(modelMap map[string]interface{}) (vpcv1.LoadBalancerListenerIdentityIntf, error) {
	model := &vpcv1.LoadBalancerListenerIdentity{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	return model, nil
}

func resourceIBMIsLbListenerPolicyLoadBalancerListenerPolicyTargetToMap(model vpcv1.LoadBalancerListenerPolicyTargetIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.LoadBalancerListenerPolicyTargetLoadBalancerPoolReference); ok {
		return resourceIBMIsLbListenerPolicyLoadBalancerListenerPolicyTargetLoadBalancerPoolReferenceToMap(model.(*vpcv1.LoadBalancerListenerPolicyTargetLoadBalancerPoolReference))
	} else if _, ok := model.(*vpcv1.LoadBalancerListenerPolicyTargetLoadBalancerListenerPolicyRedirectURL); ok {
		return resourceIBMIsLbListenerPolicyLoadBalancerListenerPolicyTargetLoadBalancerListenerPolicyRedirectURLToMap(model.(*vpcv1.LoadBalancerListenerPolicyTargetLoadBalancerListenerPolicyRedirectURL))
	} else if _, ok := model.(*vpcv1.LoadBalancerListenerPolicyTargetLoadBalancerListenerPolicyHTTPSRedirect); ok {
		return resourceIBMIsLbListenerPolicyLoadBalancerListenerPolicyTargetLoadBalancerListenerHTTPSRedirectToMap(model.(*vpcv1.LoadBalancerListenerPolicyTargetLoadBalancerListenerPolicyHTTPSRedirect))
	} else if _, ok := model.(*vpcv1.LoadBalancerListenerPolicyTarget); ok {
		modelMap := make(map[string]interface{})
		model := model.(*vpcv1.LoadBalancerListenerPolicyTarget)
		if model.Deleted != nil {
			deletedMap, err := resourceIBMIsLbListenerPolicyLoadBalancerPoolReferenceDeletedToMap(model.Deleted)
			if err != nil {
				return modelMap, err
			}
			modelMap["deleted"] = []map[string]interface{}{deletedMap}
		}
		if model.Href != nil {
			modelMap["href"] = model.Href
		}
		if model.ID != nil {
			modelMap["id"] = model.ID
		}
		if model.Name != nil {
			modelMap["name"] = model.Name
		}
		if model.HTTPStatusCode != nil {
			modelMap["http_status_code"] = flex.IntValue(model.HTTPStatusCode)
		}
		if model.URL != nil {
			modelMap["url"] = model.URL
		}
		if model.Listener != nil {
			listenerMap, err := resourceIBMIsLbListenerPolicyLoadBalancerListenerReferenceToMap(model.Listener)
			if err != nil {
				return modelMap, err
			}
			modelMap["listener"] = []map[string]interface{}{listenerMap}
		}
		if model.URI != nil {
			modelMap["uri"] = model.URI
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized vpcv1.LoadBalancerListenerPolicyTargetIntf subtype encountered")
	}
}

func resourceIBMIsLbListenerPolicyLoadBalancerPoolReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func resourceIBMIsLbListenerPolicyLoadBalancerListenerReferenceToMap(model *vpcv1.LoadBalancerListenerReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Deleted != nil {
		deletedMap, err := resourceIBMIsLbListenerPolicyLoadBalancerListenerReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	return modelMap, nil
}

func resourceIBMIsLbListenerPolicyLoadBalancerListenerReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func resourceIBMIsLbListenerPolicyLoadBalancerListenerPolicyTargetLoadBalancerPoolReferenceToMap(model *vpcv1.LoadBalancerListenerPolicyTargetLoadBalancerPoolReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Deleted != nil {
		deletedMap, err := resourceIBMIsLbListenerPolicyLoadBalancerPoolReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	return modelMap, nil
}

func resourceIBMIsLbListenerPolicyLoadBalancerListenerPolicyTargetLoadBalancerListenerPolicyRedirectURLToMap(model *vpcv1.LoadBalancerListenerPolicyTargetLoadBalancerListenerPolicyRedirectURL) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["http_status_code"] = flex.IntValue(model.HTTPStatusCode)
	modelMap["url"] = model.URL
	return modelMap, nil
}

func resourceIBMIsLbListenerPolicyLoadBalancerListenerPolicyTargetLoadBalancerListenerHTTPSRedirectToMap(model *vpcv1.LoadBalancerListenerPolicyTargetLoadBalancerListenerPolicyHTTPSRedirect) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["http_status_code"] = flex.IntValue(model.HTTPStatusCode)
	listenerMap, err := resourceIBMIsLbListenerPolicyLoadBalancerListenerReferenceToMap(model.Listener)
	if err != nil {
		return modelMap, err
	}
	modelMap["listener"] = []map[string]interface{}{listenerMap}
	if model.URI != nil {
		modelMap["uri"] = model.URI
	}
	return modelMap, nil
}

func resourceIBMIsLbListenerPolicyMapToLoadBalancerListenerPolicyTargetPatch(d *schema.ResourceData, modelMap map[string]interface{}, httpsURIRemoved *bool) (vpcv1.LoadBalancerListenerPolicyTargetPatchIntf, error) {
	model := &vpcv1.LoadBalancerListenerPolicyTargetPatch{}
	if d.HasChange("target.0.id") && modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if d.HasChange("target.0.http_status_code") && modelMap["http_status_code"] != nil {
		model.HTTPStatusCode = core.Int64Ptr(int64(modelMap["http_status_code"].(int)))
	}
	if d.HasChange("target.0.url") && modelMap["url"] != nil && modelMap["url"].(string) != "" {
		model.URL = core.StringPtr(modelMap["url"].(string))
	}
	if d.HasChange("target.0.listener") && modelMap["listener"] != nil && len(modelMap["listener"].([]interface{})) > 0 {
		ListenerModel, err := resourceIBMIsLbListenerPolicyMapToLoadBalancerListenerIdentity(modelMap["listener"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Listener = ListenerModel
	}
	if d.HasChange("target.0.uri") {
		if modelMap["uri"] != nil && modelMap["uri"].(string) != "" {
			model.URI = core.StringPtr(modelMap["uri"].(string))
		} else {
			*httpsURIRemoved = true
		}
	}
	return model, nil
}
