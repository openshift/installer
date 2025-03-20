// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isLBPoolName                          = "name"
	isLBID                                = "lb"
	isLBPoolAlgorithm                     = "algorithm"
	isLBPoolProtocol                      = "protocol"
	isLBPoolHealthDelay                   = "health_delay"
	isLBPoolHealthRetries                 = "health_retries"
	isLBPoolHealthTimeout                 = "health_timeout"
	isLBPoolHealthType                    = "health_type"
	isLBPoolHealthMonitorURL              = "health_monitor_url"
	isLBPoolHealthMonitorPort             = "health_monitor_port"
	isLBPoolSessPersistenceType           = "session_persistence_type"
	isLBPoolSessPersistenceAppCookieName  = "session_persistence_app_cookie_name"
	isLBPoolSessPersistenceHttpCookieName = "session_persistence_http_cookie_name"
	isLBPoolProvisioningStatus            = "provisioning_status"
	isLBPoolProxyProtocol                 = "proxy_protocol"
	isLBPoolActive                        = "active"
	isLBPoolCreatePending                 = "create_pending"
	isLBPoolUpdatePending                 = "update_pending"
	isLBPoolDeletePending                 = "delete_pending"
	isLBPoolMaintainancePending           = "maintenance_pending"
	isLBPoolFailed                        = "failed"
	isLBPoolDeleteDone                    = "deleted"
	isLBPool                              = "pool_id"
)

func ResourceIBMISLBPool() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMISLBPoolCreate,
		Read:     resourceIBMISLBPoolRead,
		Update:   resourceIBMISLBPoolUpdate,
		Delete:   resourceIBMISLBPoolDelete,
		Exists:   resourceIBMISLBPoolExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: customdiff.Sequence(
			func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
				return flex.ResourceIBMISLBPoolCookieValidate(diff)
			},
		),

		Schema: map[string]*schema.Schema{
			isLBPoolName: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_lb_pool", isLBPoolName),
				Description:  "Load Balancer Pool name",
			},

			isLBID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Load Balancer ID",
			},

			isLBPoolAlgorithm: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_lb_pool", isLBPoolAlgorithm),
				Description:  "Load Balancer Pool algorithm",
			},

			"failsafe_policy": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "The failsafe policy to use for this pool.If unspecified, the default failsafe policy action from the profile will be used.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "A load balancer failsafe policy action:- `forward`: Forwards requests to the `target` pool.- `fail`: Rejects requests with an HTTP `503` status code.The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
						},
						"healthy_member_threshold_count": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The healthy member count at which the failsafe policy action will be triggered. At present, this is always `0`, but may be modifiable in the future.",
						},
						"target": &schema.Schema{
							Type:             schema.TypeList,
							MaxItems:         1,
							Optional:         true,
							DiffSuppressFunc: suppressNullTarget,
							Description:      "If `action` is `forward`, the target pool to forward to.If `action` is `fail`, this property will be absent.The targets supported by this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
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
										Optional:    true,
										Computed:    true,
										Description: "The URL for this load balancer pool.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The unique identifier for this load balancer pool.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this load balancer pool. The name is unique across all pools for the load balancer.",
									},
								},
							},
						},
					},
				},
			},
			isLBPoolProtocol: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_lb_pool", isLBPoolProtocol),
				Description:  "Load Balancer Protocol",
			},

			isLBPoolHealthDelay: {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Load Blancer health delay time period",
			},

			isLBPoolHealthRetries: {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Load Balancer health retry count",
			},

			isLBPoolHealthTimeout: {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Load Balancer health timeout interval",
			},

			isLBPoolHealthType: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_lb_pool", isLBPoolHealthType),
				Description:  "Load Balancer health type",
			},

			isLBPoolHealthMonitorURL: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Health monitor URL of LB Pool",
			},

			isLBPoolHealthMonitorPort: {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Health monitor Port the LB Pool",
			},

			isLBPoolSessPersistenceType: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_lb_pool", isLBPoolSessPersistenceType),
				Description:  "Load Balancer Pool session persisence type.",
			},

			isLBPoolSessPersistenceAppCookieName: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_lb_pool", isLBPoolSessPersistenceAppCookieName),
				Description:  "Load Balancer Pool session persisence app cookie name.",
			},

			isLBPoolSessPersistenceHttpCookieName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Load Balancer Pool session persisence http cookie name.",
			},

			isLBPoolProvisioningStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the LB Pool",
			},

			isLBPoolProxyProtocol: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_lb_pool", isLBPoolProxyProtocol),
				Description:  "PROXY protocol setting for this pool",
			},

			isLBPool: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The LB Pool id",
			},

			flex.RelatedCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the LB resource",
			},
		},
	}
}

func ResourceIBMISLBPoolValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)
	algorithm := "round_robin, weighted_round_robin, least_connections"
	protocol := "http, tcp, https, udp"
	persistanceType := "source_ip, app_cookie, http_cookie"
	proxyProtocol := "disabled, v1, v2"
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isLBPoolName,
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isLBPoolAlgorithm,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              algorithm})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isLBPoolProtocol,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              protocol})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isLBPoolHealthType,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              protocol})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isLBPoolProxyProtocol,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              proxyProtocol})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isLBPoolSessPersistenceAppCookieName,
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     "^[-A-Za-z0-9!#$%&'*+.^_`~|]+$",
			MinValueLength:             1,
			MaxValueLength:             63})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isLBPoolSessPersistenceType,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              persistanceType})

	ibmISLBPoolResourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_lb_pool", Schema: validateSchema}
	return &ibmISLBPoolResourceValidator
}

func resourceIBMISLBPoolCreate(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] LB Pool create")
	name := d.Get(isLBPoolName).(string)
	lbID := d.Get(isLBID).(string)
	algorithm := d.Get(isLBPoolAlgorithm).(string)
	protocol := d.Get(isLBPoolProtocol).(string)
	healthDelay := int64(d.Get(isLBPoolHealthDelay).(int))
	maxRetries := int64(d.Get(isLBPoolHealthRetries).(int))
	healthTimeOut := int64(d.Get(isLBPoolHealthTimeout).(int))
	healthType := d.Get(isLBPoolHealthType).(string)

	var spType, cName, healthMonitorURL, pProtocol string
	var healthMonitorPort int64
	if pt, ok := d.GetOk(isLBPoolSessPersistenceType); ok {
		spType = pt.(string)
	}

	if cn, ok := d.GetOk(isLBPoolSessPersistenceAppCookieName); ok {
		cName = cn.(string)
	}
	if pp, ok := d.GetOk(isLBPoolProxyProtocol); ok {
		pProtocol = pp.(string)
	}

	if hmu, ok := d.GetOk(isLBPoolHealthMonitorURL); ok {
		healthMonitorURL = hmu.(string)
	}

	if hmp, ok := d.GetOk(isLBPoolHealthMonitorPort); ok {
		healthMonitorPort = int64(hmp.(int))
	}

	err := lbPoolCreate(d, meta, name, lbID, algorithm, protocol, healthType, spType, cName, healthMonitorURL, pProtocol, healthDelay, maxRetries, healthTimeOut, healthMonitorPort)
	if err != nil {
		return err
	}

	return resourceIBMISLBPoolRead(d, meta)
}

func lbPoolCreate(d *schema.ResourceData, meta interface{}, name, lbID, algorithm, protocol, healthType, spType, cName, healthMonitorURL, pProtocol string, healthDelay, maxRetries, healthTimeOut, healthMonitorPort int64) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	_, err = isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("[ERROR] Error checking for load balancer (%s) is active: %s", lbID, err)
	}
	healthMonitor := &vpcv1.LoadBalancerPoolHealthMonitorPrototype{
		Delay:      &healthDelay,
		MaxRetries: &maxRetries,
		Timeout:    &healthTimeOut,
		Type:       &healthType,
	}
	options := &vpcv1.CreateLoadBalancerPoolOptions{
		LoadBalancerID: &lbID,
		Algorithm:      &algorithm,
		Protocol:       &protocol,
		Name:           &name,
	}
	if healthMonitorURL != "" {
		healthMonitor.URLPath = &healthMonitorURL
	}
	if healthMonitorPort > int64(0) {
		healthMonitor.Port = &healthMonitorPort
	}
	if spType != "" {
		options.SessionPersistence = &vpcv1.LoadBalancerPoolSessionPersistencePrototype{
			Type: &spType,
		}
		if cName != "" {
			options.SessionPersistence.CookieName = &cName
		}
	}
	if pProtocol != "" {
		options.ProxyProtocol = &pProtocol
	}
	if _, ok := d.GetOk("failsafe_policy"); ok {
		failsafePolicyModel, err := resourceIBMIsLbPoolMapToLoadBalancerPoolFailsafePolicyPrototype(d.Get("failsafe_policy.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_lb_pool", "create", "parse-failsafe_policy")
		}
		options.SetFailsafePolicy(failsafePolicyModel)
	}
	options.HealthMonitor = healthMonitor
	lbPool, response, err := sess.CreateLoadBalancerPool(options)
	if err != nil {
		return fmt.Errorf("[DEBUG] lbpool create err: %s\n%s", err, response)
	}

	d.SetId(fmt.Sprintf("%s/%s", lbID, *lbPool.ID))
	log.Printf("[INFO] lbpool : %s", *lbPool.ID)

	_, err = isWaitForLBPoolActive(sess, lbID, *lbPool.ID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("[ERROR] Error checking for load balancer pool (%s) is active: %s", *lbPool.ID, err)
	}

	_, err = isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("[ERROR] Error checking for load balancer (%s) is active: %s", lbID, err)
	}

	return nil
}

func resourceIBMISLBPoolRead(d *schema.ResourceData, meta interface{}) error {

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}

	lbID := parts[0]
	lbPoolID := parts[1]

	err = lbPoolGet(d, meta, lbID, lbPoolID)
	if err != nil {
		return err
	}

	return nil
}

func lbPoolGet(d *schema.ResourceData, meta interface{}, lbID, lbPoolID string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	getLoadBalancerPoolOptions := &vpcv1.GetLoadBalancerPoolOptions{
		LoadBalancerID: &lbID,
		ID:             &lbPoolID,
	}

	lbPool, response, err := sess.GetLoadBalancerPool(getLoadBalancerPoolOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error Getting Load Balancer Pool : %s\n%s", err, response)
	}
	d.Set(isLBPoolName, *lbPool.Name)
	d.Set(isLBPool, lbPoolID)
	d.Set(isLBID, lbID)
	d.Set(isLBPoolAlgorithm, *lbPool.Algorithm)
	d.Set(isLBPoolProtocol, *lbPool.Protocol)
	if lbPool.HealthMonitor != nil {
		poolHealthMonitor := lbPool.HealthMonitor.(*vpcv1.LoadBalancerPoolHealthMonitor)
		d.Set(isLBPoolHealthDelay, *poolHealthMonitor.Delay)
		d.Set(isLBPoolHealthRetries, *poolHealthMonitor.MaxRetries)
		d.Set(isLBPoolHealthTimeout, *poolHealthMonitor.Timeout)
		if poolHealthMonitor.Type != nil {
			d.Set(isLBPoolHealthType, *poolHealthMonitor.Type)
		}
		if poolHealthMonitor.URLPath != nil {
			d.Set(isLBPoolHealthMonitorURL, *poolHealthMonitor.URLPath)
		}
		if poolHealthMonitor.Port != nil {
			d.Set(isLBPoolHealthMonitorPort, *poolHealthMonitor.Port)
		}
	}
	if !core.IsNil(lbPool.FailsafePolicy) {
		failsafePolicyMap, err := resourceIBMIsLbPoolLoadBalancerPoolFailsafePolicyToMap(lbPool.FailsafePolicy)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_lb_pool", "read", "failsafe_policy-to-map")
		}
		if err = d.Set("failsafe_policy", []map[string]interface{}{failsafePolicyMap}); err != nil {
			err = fmt.Errorf("Error setting failsafe_policy: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_lb_pool", "read", "set-failsafe_policy")
		}
	}
	if lbPool.SessionPersistence != nil {
		d.Set(isLBPoolSessPersistenceType, *lbPool.SessionPersistence.Type)
		if lbPool.SessionPersistence.CookieName != nil {
			if *lbPool.SessionPersistence.Type == "app_cookie" {
				d.Set(isLBPoolSessPersistenceAppCookieName, *lbPool.SessionPersistence.CookieName)
				d.Set(isLBPoolSessPersistenceHttpCookieName, "")
			} else if *lbPool.SessionPersistence.Type == "http_cookie" {
				d.Set(isLBPoolSessPersistenceHttpCookieName, *lbPool.SessionPersistence.CookieName)
				d.Set(isLBPoolSessPersistenceAppCookieName, "")
			}

		}
	} else {
		d.Set(isLBPoolSessPersistenceType, "")
		d.Set(isLBPoolSessPersistenceHttpCookieName, "")
		d.Set(isLBPoolSessPersistenceAppCookieName, "")
	}

	d.Set(isLBPoolProvisioningStatus, *lbPool.ProvisioningStatus)
	d.Set(isLBPoolProxyProtocol, *lbPool.ProxyProtocol)
	getLoadBalancerOptions := &vpcv1.GetLoadBalancerOptions{
		ID: &lbID,
	}
	lb, response, err := sess.GetLoadBalancer(getLoadBalancerOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Getting Load Balancer : %s\n%s", err, response)
	}
	d.Set(flex.RelatedCRN, *lb.CRN)
	return nil
}

func resourceIBMISLBPoolUpdate(d *schema.ResourceData, meta interface{}) error {

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}

	lbID := parts[0]
	lbPoolID := parts[1]

	err = lbPoolUpdate(d, meta, lbID, lbPoolID)
	if err != nil {
		return err
	}

	return resourceIBMISLBPoolRead(d, meta)
}

func lbPoolUpdate(d *schema.ResourceData, meta interface{}, lbID, lbPoolID string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	hasChanged := false

	updateLoadBalancerPoolOptions := &vpcv1.UpdateLoadBalancerPoolOptions{
		LoadBalancerID: &lbID,
		ID:             &lbPoolID,
	}

	loadBalancerPoolPatchModel := &vpcv1.LoadBalancerPoolPatch{}

	lBPoolHealthMonitorPortRemoved := false
	isFailSafePolicyTargetNull := false
	hasFailSafeChanged := false
	if d.HasChange("failsafe_policy") {
		failsafePolicy := &vpcv1.LoadBalancerPoolFailsafePolicyPatch{}
		hasFailSafeChanged = true
		if d.HasChange("failsafe_policy.0.action") {
			failsafepolicyAction := d.Get("failsafe_policy.0.action").(string)
			failsafePolicy.Action = &failsafepolicyAction
		}
		if d.HasChange("failsafe_policy.0.target") {
			failsafepolicyAction := d.Get("failsafe_policy.0.action").(string)
			failsafePolicy.Action = &failsafepolicyAction
			targetPatch := &vpcv1.LoadBalancerPoolFailsafePolicyTargetPatch{}
			failsafepolicyTargetId := d.Get("failsafe_policy.0.target.0.id").(string)
			if failsafepolicyTargetId == "null" {
				isFailSafePolicyTargetNull = true
				var nullStringPtr *string
				targetPatch.ID = nullStringPtr
				failsafePolicy.Target = targetPatch
			} else {
				targetPatch.ID = &failsafepolicyTargetId
				failsafePolicy.Target = targetPatch

			}
			failsafepolicyTargetHref := d.Get("failsafe_policy.0.target.0.href").(string)
			if failsafepolicyTargetHref == "null" {
				isFailSafePolicyTargetNull = true
				var nullStringPtr *string
				targetPatch.Href = nullStringPtr
				failsafePolicy.Target = targetPatch
			} else {
				targetPatch.Href = &failsafepolicyTargetHref
				failsafePolicy.Target = targetPatch
			}
		}
		loadBalancerPoolPatchModel.FailsafePolicy = failsafePolicy
	}

	if d.HasChange(isLBPoolHealthDelay) || d.HasChange(isLBPoolHealthRetries) ||
		d.HasChange(isLBPoolHealthTimeout) || d.HasChange(isLBPoolHealthType) || d.HasChange(isLBPoolHealthMonitorURL) || d.HasChange(isLBPoolHealthMonitorPort) || d.HasChange("failsafe_policy") {

		delay := int64(d.Get(isLBPoolHealthDelay).(int))
		maxretries := int64(d.Get(isLBPoolHealthRetries).(int))
		timeout := int64(d.Get(isLBPoolHealthTimeout).(int))
		healthtype := d.Get(isLBPoolHealthType).(string)
		urlpath := d.Get(isLBPoolHealthMonitorURL).(string)
		healthMonitorTemplate := &vpcv1.LoadBalancerPoolHealthMonitorPatch{
			Delay:      &delay,
			MaxRetries: &maxretries,
			Timeout:    &timeout,
			Type:       &healthtype,
			URLPath:    &urlpath,
		}
		port := int64(d.Get(isLBPoolHealthMonitorPort).(int))
		if port > int64(0) {
			healthMonitorTemplate.Port = &port
		} else {
			lBPoolHealthMonitorPortRemoved = true
		}
		loadBalancerPoolPatchModel.HealthMonitor = healthMonitorTemplate
		hasChanged = true
	}

	sessionPersistenceRemoved := false
	if d.HasChange(isLBPoolSessPersistenceType) || d.HasChange(isLBPoolSessPersistenceAppCookieName) {
		sesspersistancetype := d.Get(isLBPoolSessPersistenceType).(string)
		sessPersistanceCookieName := d.Get(isLBPoolSessPersistenceAppCookieName).(string)
		sessionPersistence := &vpcv1.LoadBalancerPoolSessionPersistencePatch{}
		if sesspersistancetype != "" {
			sessionPersistence.Type = &sesspersistancetype
			if sessPersistanceCookieName != "" {
				sessionPersistence.CookieName = &sessPersistanceCookieName
			}
		} else {
			sessionPersistenceRemoved = true
		}

		loadBalancerPoolPatchModel.SessionPersistence = sessionPersistence

		hasChanged = true
	}

	if d.HasChange(isLBPoolProxyProtocol) {
		proxyProtocol := d.Get(isLBPoolProxyProtocol).(string)
		loadBalancerPoolPatchModel.ProxyProtocol = &proxyProtocol
		hasChanged = true
	}

	if d.HasChange(isLBPoolName) || d.HasChange(isLBPoolAlgorithm) || d.HasChange(isLBPoolProtocol) || hasChanged || hasFailSafeChanged {
		name := d.Get(isLBPoolName).(string)
		algorithm := d.Get(isLBPoolAlgorithm).(string)
		protocol := d.Get(isLBPoolProtocol).(string)

		loadBalancerPoolPatchModel.Algorithm = &algorithm
		loadBalancerPoolPatchModel.Name = &name
		loadBalancerPoolPatchModel.Protocol = &protocol

		isLBKey := "load_balancer_key_" + lbID
		conns.IbmMutexKV.Lock(isLBKey)
		defer conns.IbmMutexKV.Unlock(isLBKey)
		_, err := isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf(
				"Error checking for load balancer (%s) is active: %s", lbID, err)
		}

		_, err = isWaitForLBPoolActive(sess, lbID, lbPoolID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf(
				"Error checking for load balancer pool (%s) is active: %s", lbPoolID, err)
		}

		LoadBalancerPoolPatch, err := loadBalancerPoolPatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("[ERROR] Error calling asPatch for LoadBalancerPoolPatch: %s", err)
		}
		if sessionPersistenceRemoved {
			LoadBalancerPoolPatch["session_persistence"] = nil
		}
		if lBPoolHealthMonitorPortRemoved {
			LoadBalancerPoolPatch["health_monitor"].(map[string]interface{})["port"] = nil
		}
		if isFailSafePolicyTargetNull {
			LoadBalancerPoolPatch["failsafe_policy"].(map[string]interface{})["target"] = nil
		}

		updateLoadBalancerPoolOptions.LoadBalancerPoolPatch = LoadBalancerPoolPatch

		_, response, err := sess.UpdateLoadBalancerPool(updateLoadBalancerPoolOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Updating Load Balancer Pool : %s\n%s", err, response)
		}

		_, err = isWaitForLBPoolActive(sess, lbID, lbPoolID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf(
				"Error checking for load balancer pool (%s) is active: %s", lbPoolID, err)
		}

		_, err = isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf(
				"Error checking for load balancer (%s) is active: %s", lbID, err)
		}
	}
	return nil
}

func resourceIBMISLBPoolDelete(d *schema.ResourceData, meta interface{}) error {

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}

	lbID := parts[0]
	lbPoolID := parts[1]

	isLBKey := "load_balancer_key_" + lbID
	conns.IbmMutexKV.Lock(isLBKey)
	defer conns.IbmMutexKV.Unlock(isLBKey)

	err = lbPoolDelete(d, meta, lbID, lbPoolID)
	if err != nil {
		return err
	}

	return nil
}

func lbPoolDelete(d *schema.ResourceData, meta interface{}, lbID, lbPoolID string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	getLoadBalancerPoolOptions := &vpcv1.GetLoadBalancerPoolOptions{
		LoadBalancerID: &lbID,
		ID:             &lbPoolID,
	}
	_, response, err := sess.GetLoadBalancerPool(getLoadBalancerPoolOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error Getting vpc load balancer pool(%s): %s\n%s", lbPoolID, err, response)
	}
	_, err = isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf("[ERROR] Error checking for load balancer (%s) is active: %s", lbID, err)
	}

	_, err = isWaitForLBPoolActive(sess, lbID, lbPoolID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf("[ERROR] Error checking for load balancer pool (%s) is active: %s", lbPoolID, err)
	}

	deleteLoadBalancerPoolOptions := &vpcv1.DeleteLoadBalancerPoolOptions{
		LoadBalancerID: &lbID,
		ID:             &lbPoolID,
	}
	response, err = sess.DeleteLoadBalancerPool(deleteLoadBalancerPoolOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Deleting Load Balancer Pool : %s\n%s", err, response)
	}
	_, err = isWaitForLBPoolDeleted(sess, lbID, lbPoolID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf("[ERROR] Error checking for load balancer pool (%s) is deleted: %s", lbPoolID, err)
	}

	_, err = isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf("[ERROR] Error checking for load balancer (%s) is active: %s", lbID, err)
	}
	d.SetId("")
	return nil
}

func resourceIBMISLBPoolExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return false, err
	}
	if len(parts) != 2 {
		return false, fmt.Errorf("[ERROR] Incorrect ID %s: ID should be a combination of lbID/lbPoolID", d.Id())
	}

	lbID := parts[0]
	lbPoolID := parts[1]

	exists, err := lbPoolExists(d, meta, lbID, lbPoolID)
	return exists, err

}

func lbPoolExists(d *schema.ResourceData, meta interface{}, lbID, lbPoolID string) (bool, error) {
	sess, err := vpcClient(meta)
	if err != nil {
		return false, err
	}

	getLoadBalancerPoolOptions := &vpcv1.GetLoadBalancerPoolOptions{
		LoadBalancerID: &lbID,
		ID:             &lbPoolID,
	}
	_, response, err := sess.GetLoadBalancerPool(getLoadBalancerPoolOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("[ERROR] Error getting Load balancer pool: %s\n%s", err, response)
	}
	return true, nil
}

func isWaitForLBPoolActive(sess *vpcv1.VpcV1, lbId, lbPoolId string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for load balancer pool (%s) to be available.", lbPoolId)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{isLBPoolCreatePending, isLBPoolUpdatePending, isLBPoolMaintainancePending},
		Target:     []string{isLBPoolActive, ""},
		Refresh:    isLBPoolRefreshFunc(sess, lbId, lbPoolId),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isLBPoolRefreshFunc(sess *vpcv1.VpcV1, lbId, lbPoolId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		getlbpOptions := &vpcv1.GetLoadBalancerPoolOptions{
			LoadBalancerID: &lbId,
			ID:             &lbPoolId,
		}
		lbPool, response, err := sess.GetLoadBalancerPool(getlbpOptions)
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] Error Getting Load Balancer Pool: %s\n%s", err, response)
		}

		if *lbPool.ProvisioningStatus == isLBPoolActive || *lbPool.ProvisioningStatus == isLBPoolFailed {
			return lbPool, isLBPoolActive, nil
		}

		return lbPool, *lbPool.ProvisioningStatus, nil
	}
}

func isWaitForLBPoolDeleted(lbc *vpcv1.VpcV1, lbId, lbPoolId string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for  (%s) to be deleted.", lbPoolId)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{isLBPoolUpdatePending, isLBPoolMaintainancePending, isLBPoolDeletePending},
		Target:     []string{isLBPoolDeleteDone, ""},
		Refresh:    isLBPoolDeleteRefreshFunc(lbc, lbId, lbPoolId),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isLBPoolDeleteRefreshFunc(lbc *vpcv1.VpcV1, lbId, lbPoolId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] is lb pool delete function here")
		getlbpOptions := &vpcv1.GetLoadBalancerPoolOptions{
			LoadBalancerID: &lbId,
			ID:             &lbPoolId,
		}
		lbPool, response, err := lbc.GetLoadBalancerPool(getlbpOptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return lbPool, isLBPoolDeleteDone, nil
			}
			return nil, "", fmt.Errorf("[ERROR] The vpc load balancer pool %s failed to delete: %s\n%s", lbPoolId, err, response)
		}
		return lbPool, isLBPoolDeletePending, nil
	}
}

func resourceIBMIsLbPoolMapToLoadBalancerPoolFailsafePolicyPrototype(modelMap map[string]interface{}) (*vpcv1.LoadBalancerPoolFailsafePolicyPrototype, error) {
	model := &vpcv1.LoadBalancerPoolFailsafePolicyPrototype{}
	if modelMap["action"] != nil && modelMap["action"].(string) != "" {
		model.Action = core.StringPtr(modelMap["action"].(string))
	}
	if modelMap["target"] != nil && len(modelMap["target"].([]interface{})) > 0 {
		TargetModel, err := resourceIBMIsLbPoolMapToLoadBalancerPoolIdentity(modelMap["target"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Target = TargetModel
	}
	return model, nil
}

func resourceIBMIsLbPoolMapToLoadBalancerPoolIdentity(modelMap map[string]interface{}) (vpcv1.LoadBalancerPoolIdentityIntf, error) {
	model := &vpcv1.LoadBalancerPoolIdentity{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	return model, nil
}

func resourceIBMIsLbPoolLoadBalancerPoolFailsafePolicyToMap(model *vpcv1.LoadBalancerPoolFailsafePolicy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["action"] = *model.Action
	modelMap["healthy_member_threshold_count"] = flex.IntValue(model.HealthyMemberThresholdCount)
	if model.Target != nil {
		targetMap, err := resourceIBMIsLbPoolLoadBalancerPoolReferenceToMap(model.Target)
		if err != nil {
			return modelMap, err
		}
		modelMap["target"] = []map[string]interface{}{targetMap}
	}
	return modelMap, nil
}

func resourceIBMIsLbPoolLoadBalancerPoolReferenceToMap(model *vpcv1.LoadBalancerPoolReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Deleted != nil {
		deletedMap, err := resourceIBMIsLbPoolDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	return modelMap, nil
}

func resourceIBMIsLbPoolDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = *model.MoreInfo
	return modelMap, nil
}

func suppressNullTarget(k, old, new string, d *schema.ResourceData) bool {
	// If resource already exists (has an ID) and new value is "null" while old is empty
	// then suppress the diff
	if new != old && new == "null" && old == "" && d.Id() != "" {
		return true
	}

	oldId, newId := d.GetChange("failsafe_policy.0.target.0.id")
	oldHref, newHref := d.GetChange("failsafe_policy.0.target.0.href")

	// Check id field
	if newId.(string) == "null" && oldId.(string) == "" && d.Id() != "" {
		return true
	}

	// Check href field
	if newHref.(string) == "null" && oldHref.(string) == "" && d.Id() != "" {
		return true
	}

	return false
}
