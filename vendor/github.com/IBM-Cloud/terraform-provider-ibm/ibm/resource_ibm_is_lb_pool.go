// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcclassicv1"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	isLBPoolName                      = "name"
	isLBID                            = "lb"
	isLBPoolAlgorithm                 = "algorithm"
	isLBPoolProtocol                  = "protocol"
	isLBPoolHealthDelay               = "health_delay"
	isLBPoolHealthRetries             = "health_retries"
	isLBPoolHealthTimeout             = "health_timeout"
	isLBPoolHealthType                = "health_type"
	isLBPoolHealthMonitorURL          = "health_monitor_url"
	isLBPoolHealthMonitorPort         = "health_monitor_port"
	isLBPoolSessPersistenceType       = "session_persistence_type"
	isLBPoolSessPersistenceCookieName = "session_persistence_cookie_name"
	isLBPoolProvisioningStatus        = "provisioning_status"
	isLBPoolProxyProtocol             = "proxy_protocol"
	isLBPoolActive                    = "active"
	isLBPoolCreatePending             = "create_pending"
	isLBPoolUpdatePending             = "update_pending"
	isLBPoolDeletePending             = "delete_pending"
	isLBPoolMaintainancePending       = "maintenance_pending"
	isLBPoolFailed                    = "failed"
	isLBPoolDeleteDone                = "deleted"
	isLBPool                          = "pool_id"
)

func resourceIBMISLBPool() *schema.Resource {
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

		Schema: map[string]*schema.Schema{
			isLBPoolName: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: InvokeValidator("ibm_is_lb_pool", isLBPoolName),
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
				ValidateFunc: InvokeValidator("ibm_is_lb_pool", isLBPoolAlgorithm),
				Description:  "Load Balancer Pool algorithm",
			},

			isLBPoolProtocol: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: InvokeValidator("ibm_is_lb_pool", isLBPoolProtocol),
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
				ValidateFunc: InvokeValidator("ibm_is_lb_pool", isLBPoolHealthType),
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
				ValidateFunc: InvokeValidator("ibm_is_lb_pool", isLBPoolSessPersistenceType),
				Description:  "Load Balancer Pool session persisence type.",
			},

			isLBPoolSessPersistenceCookieName: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Load Balancer Pool session persisence cookie name",
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
				ValidateFunc: InvokeValidator("ibm_is_lb_pool", isLBPoolProxyProtocol),
				Description:  "PROXY protocol setting for this pool",
			},

			isLBPool: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The LB Pool id",
			},

			RelatedCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the LB resource",
			},
		},
	}
}

func resourceIBMISLBPoolValidator() *ResourceValidator {

	validateSchema := make([]ValidateSchema, 1)
	algorithm := "round_robin, weighted_round_robin, least_connections"
	protocol := "http, tcp, https"
	persistanceType := "source_ip"
	proxyProtocol := "disabled, v1, v2"
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isLBPoolName,
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Required:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isLBPoolAlgorithm,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              algorithm})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isLBPoolProtocol,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              protocol})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isLBPoolHealthType,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              protocol})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isLBPoolProxyProtocol,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              proxyProtocol})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isLBPoolSessPersistenceType,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              persistanceType})

	ibmISLBPoolResourceValidator := ResourceValidator{ResourceName: "ibm_is_lb_pool", Schema: validateSchema}
	return &ibmISLBPoolResourceValidator
}

func resourceIBMISLBPoolCreate(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

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

	if cn, ok := d.GetOk(isLBPoolSessPersistenceCookieName); ok {
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
	isLBKey := "load_balancer_key_" + lbID
	ibmMutexKV.Lock(isLBKey)
	defer ibmMutexKV.Unlock(isLBKey)

	if userDetails.generation == 1 {
		err := classicLBPoolCreate(d, meta, name, lbID, algorithm, protocol, healthType, spType, cName, healthMonitorURL, healthDelay, maxRetries, healthTimeOut, healthMonitorPort)
		if err != nil {
			return err
		}
	} else {
		err := lbPoolCreate(d, meta, name, lbID, algorithm, protocol, healthType, spType, cName, healthMonitorURL, pProtocol, healthDelay, maxRetries, healthTimeOut, healthMonitorPort)
		if err != nil {
			return err
		}
	}
	return resourceIBMISLBPoolRead(d, meta)
}

func classicLBPoolCreate(d *schema.ResourceData, meta interface{}, name, lbID, algorithm, protocol, healthType, spType, cName, healthMonitorURL string, healthDelay, maxRetries, healthTimeOut, healthMonitorPort int64) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}

	_, err = isWaitForClassicLBAvailable(sess, lbID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf(
			"Error checking for load balancer (%s) is active: %s", lbID, err)
	}

	options := &vpcclassicv1.CreateLoadBalancerPoolOptions{
		LoadBalancerID: &lbID,
		Algorithm:      &algorithm,
		Protocol:       &protocol,
		Name:           &name,
		HealthMonitor: &vpcclassicv1.LoadBalancerPoolHealthMonitorPrototype{
			Delay:      &healthDelay,
			MaxRetries: &maxRetries,
			Timeout:    &healthTimeOut,
			Type:       &healthType,
		},
	}
	if healthMonitorURL != "" {
		options.HealthMonitor.URLPath = &healthMonitorURL
	}
	if healthMonitorPort > int64(0) {
		options.HealthMonitor.Port = &healthMonitorPort
	}
	if spType != "" {
		options.SessionPersistence = &vpcclassicv1.LoadBalancerPoolSessionPersistencePrototype{
			Type: &spType,
		}
	}
	lbPool, response, err := sess.CreateLoadBalancerPool(options)
	if err != nil {
		return fmt.Errorf("[DEBUG] lbpool create err: %s\n%s", err, response)
	}

	d.SetId(fmt.Sprintf("%s/%s", lbID, *lbPool.ID))
	log.Printf("[INFO] lbpool : %s", *lbPool.ID)

	_, err = isWaitForClassicLBPoolActive(sess, lbID, *lbPool.ID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf(
			"Error checking for load balancer pool (%s) is active: %s", *lbPool.ID, err)
	}

	_, err = isWaitForClassicLBAvailable(sess, lbID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf(
			"Error checking for load balancer (%s) is active: %s", lbID, err)
	}

	return nil
}

func lbPoolCreate(d *schema.ResourceData, meta interface{}, name, lbID, algorithm, protocol, healthType, spType, cName, healthMonitorURL, pProtocol string, healthDelay, maxRetries, healthTimeOut, healthMonitorPort int64) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	_, err = isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf(
			"Error checking for load balancer (%s) is active: %s", lbID, err)
	}

	options := &vpcv1.CreateLoadBalancerPoolOptions{
		LoadBalancerID: &lbID,
		Algorithm:      &algorithm,
		Protocol:       &protocol,
		Name:           &name,
		HealthMonitor: &vpcv1.LoadBalancerPoolHealthMonitorPrototype{
			Delay:      &healthDelay,
			MaxRetries: &maxRetries,
			Timeout:    &healthTimeOut,
			Type:       &healthType,
		},
	}
	if healthMonitorURL != "" {
		options.HealthMonitor.URLPath = &healthMonitorURL
	}
	if healthMonitorPort > int64(0) {
		options.HealthMonitor.Port = &healthMonitorPort
	}
	if spType != "" {
		options.SessionPersistence = &vpcv1.LoadBalancerPoolSessionPersistencePrototype{
			Type: &spType,
		}
	}
	if pProtocol != "" {
		options.ProxyProtocol = &pProtocol
	}
	lbPool, response, err := sess.CreateLoadBalancerPool(options)
	if err != nil {
		return fmt.Errorf("[DEBUG] lbpool create err: %s\n%s", err, response)
	}

	d.SetId(fmt.Sprintf("%s/%s", lbID, *lbPool.ID))
	log.Printf("[INFO] lbpool : %s", *lbPool.ID)

	_, err = isWaitForLBPoolActive(sess, lbID, *lbPool.ID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf(
			"Error checking for load balancer pool (%s) is active: %s", *lbPool.ID, err)
	}

	_, err = isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf(
			"Error checking for load balancer (%s) is active: %s", lbID, err)
	}

	return nil
}

func resourceIBMISLBPoolRead(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}

	lbID := parts[0]
	lbPoolID := parts[1]

	if userDetails.generation == 1 {
		err := classicLBPoolGet(d, meta, lbID, lbPoolID)
		if err != nil {
			return err
		}
	} else {
		err := lbPoolGet(d, meta, lbID, lbPoolID)
		if err != nil {
			return err
		}
	}
	return nil
}

func classicLBPoolGet(d *schema.ResourceData, meta interface{}, lbID, lbPoolID string) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	getLoadBalancerPoolOptions := &vpcclassicv1.GetLoadBalancerPoolOptions{
		LoadBalancerID: &lbID,
		ID:             &lbPoolID,
	}
	lbPool, response, err := sess.GetLoadBalancerPool(getLoadBalancerPoolOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Getting Load Balancer Pool : %s\n%s", err, response)
	}

	d.Set(isLBPoolName, *lbPool.Name)
	d.Set(isLBPool, lbPoolID)
	d.Set(isLBID, lbID)
	d.Set(isLBPoolAlgorithm, *lbPool.Algorithm)
	d.Set(isLBPoolProtocol, *lbPool.Protocol)
	d.Set(isLBPoolHealthDelay, *lbPool.HealthMonitor.Delay)
	d.Set(isLBPoolHealthRetries, *lbPool.HealthMonitor.MaxRetries)
	d.Set(isLBPoolHealthTimeout, *lbPool.HealthMonitor.Timeout)
	if lbPool.HealthMonitor.Type != nil {
		d.Set(isLBPoolHealthType, *lbPool.HealthMonitor.Type)
	}
	if lbPool.HealthMonitor.URLPath != nil {
		d.Set(isLBPoolHealthMonitorURL, *lbPool.HealthMonitor.URLPath)
	}
	if lbPool.HealthMonitor.Port != nil {
		d.Set(isLBPoolHealthMonitorPort, *lbPool.HealthMonitor.Port)
	}
	if lbPool.SessionPersistence != nil {
		d.Set(isLBPoolSessPersistenceType, *lbPool.SessionPersistence.Type)
		// d.Set(isLBPoolSessPersistenceCookieName, *lbPool.SessionPersistence.CookieName)
	}
	d.Set(isLBPoolProvisioningStatus, *lbPool.ProvisioningStatus)
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
		return fmt.Errorf("Error Getting Load Balancer Pool : %s\n%s", err, response)
	}
	d.Set(isLBPoolName, *lbPool.Name)
	d.Set(isLBPool, lbPoolID)
	d.Set(isLBID, lbID)
	d.Set(isLBPoolAlgorithm, *lbPool.Algorithm)
	d.Set(isLBPoolProtocol, *lbPool.Protocol)
	d.Set(isLBPoolHealthDelay, *lbPool.HealthMonitor.Delay)
	d.Set(isLBPoolHealthRetries, *lbPool.HealthMonitor.MaxRetries)
	d.Set(isLBPoolHealthTimeout, *lbPool.HealthMonitor.Timeout)
	if lbPool.HealthMonitor.Type != nil {
		d.Set(isLBPoolHealthType, *lbPool.HealthMonitor.Type)
	}
	if lbPool.HealthMonitor.URLPath != nil {
		d.Set(isLBPoolHealthMonitorURL, *lbPool.HealthMonitor.URLPath)
	}
	if lbPool.HealthMonitor.Port != nil {
		d.Set(isLBPoolHealthMonitorPort, *lbPool.HealthMonitor.Port)
	}
	if lbPool.SessionPersistence != nil {
		d.Set(isLBPoolSessPersistenceType, *lbPool.SessionPersistence.Type)
		// d.Set(isLBPoolSessPersistenceCookieName, *lbPool.SessionPersistence.CookieName)
	}
	d.Set(isLBPoolProvisioningStatus, *lbPool.ProvisioningStatus)
	d.Set(isLBPoolProxyProtocol, *lbPool.ProxyProtocol)
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

func resourceIBMISLBPoolUpdate(d *schema.ResourceData, meta interface{}) error {

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}

	lbID := parts[0]
	lbPoolID := parts[1]

	if userDetails.generation == 1 {
		err := classicLBPoolUpdate(d, meta, lbID, lbPoolID)
		if err != nil {
			return err
		}
	} else {
		err := lbPoolUpdate(d, meta, lbID, lbPoolID)
		if err != nil {
			return err
		}
	}
	return resourceIBMISLBPoolRead(d, meta)
}

func classicLBPoolUpdate(d *schema.ResourceData, meta interface{}, lbID, lbPoolID string) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}

	hasChanged := false

	updateLoadBalancerPoolOptions := &vpcclassicv1.UpdateLoadBalancerPoolOptions{
		LoadBalancerID: &lbID,
		ID:             &lbPoolID,
	}

	loadBalancerPoolPatchModel := &vpcclassicv1.LoadBalancerPoolPatch{}

	if d.HasChange(isLBPoolHealthDelay) || d.HasChange(isLBPoolHealthRetries) ||
		d.HasChange(isLBPoolHealthTimeout) || d.HasChange(isLBPoolHealthType) || d.HasChange(isLBPoolHealthMonitorURL) || d.HasChange(isLBPoolHealthMonitorPort) {

		delay := int64(d.Get(isLBPoolHealthDelay).(int))
		maxretries := int64(d.Get(isLBPoolHealthRetries).(int))
		timeout := int64(d.Get(isLBPoolHealthTimeout).(int))
		healthtype := d.Get(isLBPoolHealthType).(string)
		urlpath := d.Get(isLBPoolHealthMonitorURL).(string)
		healthMonitorTemplate := &vpcclassicv1.LoadBalancerPoolHealthMonitorPatch{
			Delay:      &delay,
			MaxRetries: &maxretries,
			Timeout:    &timeout,
			Type:       &healthtype,
			URLPath:    &urlpath,
		}
		port := int64(d.Get(isLBPoolHealthMonitorPort).(int))
		if port > int64(0) {
			healthMonitorTemplate.Port = &port
		}

		loadBalancerPoolPatchModel.HealthMonitor = healthMonitorTemplate
		hasChanged = true
	}

	if d.HasChange(isLBPoolSessPersistenceType) || d.HasChange(isLBPoolSessPersistenceCookieName) {
		sesspersistancetype := d.Get(isLBPoolSessPersistenceType).(string)
		sessionPersistence := &vpcclassicv1.LoadBalancerPoolSessionPersistencePatch{
			Type: &sesspersistancetype,
		}
		loadBalancerPoolPatchModel.SessionPersistence = sessionPersistence
		hasChanged = true
	}

	if d.HasChange(isLBPoolName) || d.HasChange(isLBPoolAlgorithm) || d.HasChange(isLBPoolProtocol) || hasChanged {
		name := d.Get(isLBPoolName).(string)
		algorithm := d.Get(isLBPoolAlgorithm).(string)
		protocol := d.Get(isLBPoolProtocol).(string)

		loadBalancerPoolPatchModel.Algorithm = &algorithm
		loadBalancerPoolPatchModel.Name = &name
		loadBalancerPoolPatchModel.Protocol = &protocol

		isLBKey := "load_balancer_key_" + lbID
		ibmMutexKV.Lock(isLBKey)
		defer ibmMutexKV.Unlock(isLBKey)
		_, err := isWaitForClassicLBAvailable(sess, lbID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf(
				"Error checking for load balancer (%s) is active: %s", lbID, err)
		}

		_, err = isWaitForClassicLBPoolActive(sess, lbID, lbPoolID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf(
				"Error checking for load balancer pool (%s) is active: %s", lbPoolID, err)
		}

		LoadBalancerPoolPatch, err := loadBalancerPoolPatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("Error calling asPatch for LoadBalancerPoolPatch: %s", err)
		}
		updateLoadBalancerPoolOptions.LoadBalancerPoolPatch = LoadBalancerPoolPatch

		_, response, err := sess.UpdateLoadBalancerPool(updateLoadBalancerPoolOptions)
		if err != nil {
			return fmt.Errorf("Error Updating Load Balancer Pool : %s\n%s", err, response)
		}

		_, err = isWaitForClassicLBPoolActive(sess, lbID, lbPoolID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf(
				"Error checking for load balancer pool (%s) is active: %s", lbPoolID, err)
		}

		_, err = isWaitForClassicLBAvailable(sess, lbID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf(
				"Error checking for load balancer (%s) is active: %s", lbID, err)
		}
	}
	return nil
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

	if d.HasChange(isLBPoolHealthDelay) || d.HasChange(isLBPoolHealthRetries) ||
		d.HasChange(isLBPoolHealthTimeout) || d.HasChange(isLBPoolHealthType) || d.HasChange(isLBPoolHealthMonitorURL) || d.HasChange(isLBPoolHealthMonitorPort) {

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
		}
		loadBalancerPoolPatchModel.HealthMonitor = healthMonitorTemplate
		hasChanged = true
	}

	if d.HasChange(isLBPoolSessPersistenceType) || d.HasChange(isLBPoolSessPersistenceCookieName) {
		sesspersistancetype := d.Get(isLBPoolSessPersistenceType).(string)
		sessionPersistence := &vpcv1.LoadBalancerPoolSessionPersistencePatch{
			Type: &sesspersistancetype,
		}
		loadBalancerPoolPatchModel.SessionPersistence = sessionPersistence
		hasChanged = true
	}

	if d.HasChange(isLBPoolProxyProtocol) {
		proxyProtocol := d.Get(isLBPoolProxyProtocol).(string)
		loadBalancerPoolPatchModel.ProxyProtocol = &proxyProtocol
		hasChanged = true
	}

	if d.HasChange(isLBPoolName) || d.HasChange(isLBPoolAlgorithm) || d.HasChange(isLBPoolProtocol) || hasChanged {
		name := d.Get(isLBPoolName).(string)
		algorithm := d.Get(isLBPoolAlgorithm).(string)
		protocol := d.Get(isLBPoolProtocol).(string)

		loadBalancerPoolPatchModel.Algorithm = &algorithm
		loadBalancerPoolPatchModel.Name = &name
		loadBalancerPoolPatchModel.Protocol = &protocol

		isLBKey := "load_balancer_key_" + lbID
		ibmMutexKV.Lock(isLBKey)
		defer ibmMutexKV.Unlock(isLBKey)
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
			return fmt.Errorf("Error calling asPatch for LoadBalancerPoolPatch: %s", err)
		}
		updateLoadBalancerPoolOptions.LoadBalancerPoolPatch = LoadBalancerPoolPatch

		_, response, err := sess.UpdateLoadBalancerPool(updateLoadBalancerPoolOptions)
		if err != nil {
			return fmt.Errorf("Error Updating Load Balancer Pool : %s\n%s", err, response)
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

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}

	lbID := parts[0]
	lbPoolID := parts[1]

	isLBKey := "load_balancer_key_" + lbID
	ibmMutexKV.Lock(isLBKey)
	defer ibmMutexKV.Unlock(isLBKey)

	if userDetails.generation == 1 {
		err := classicLBPoolDelete(d, meta, lbID, lbPoolID)
		if err != nil {
			return err
		}
	} else {
		err := lbPoolDelete(d, meta, lbID, lbPoolID)
		if err != nil {
			return err
		}
	}
	return nil
}

func classicLBPoolDelete(d *schema.ResourceData, meta interface{}, lbID, lbPoolID string) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	getLoadBalancerPoolOptions := &vpcclassicv1.GetLoadBalancerPoolOptions{
		LoadBalancerID: &lbID,
		ID:             &lbPoolID,
	}
	_, response, err := sess.GetLoadBalancerPool(getLoadBalancerPoolOptions)

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Getting vpc load balancer pool(%s): %s\n%s", lbPoolID, err, response)
	}
	_, err = isWaitForClassicLBAvailable(sess, lbID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf(
			"Error checking for load balancer (%s) is active: %s", lbID, err)
	}

	_, err = isWaitForClassicLBPoolActive(sess, lbID, lbPoolID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf(
			"Error checking for load balancer pool (%s) is active: %s", lbPoolID, err)
	}

	deleteLoadBalancerPoolOptions := &vpcclassicv1.DeleteLoadBalancerPoolOptions{
		LoadBalancerID: &lbID,
		ID:             &lbPoolID,
	}
	response, err = sess.DeleteLoadBalancerPool(deleteLoadBalancerPoolOptions)
	if err != nil {
		return fmt.Errorf("Error Deleting Load Balancer Pool : %s\n%s", err, response)
	}
	_, err = isWaitForClassicLBPoolDeleted(sess, lbID, lbPoolID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf(
			"Error checking for load balancer pool (%s) is deleted: %s", lbPoolID, err)
	}

	_, err = isWaitForClassicLBAvailable(sess, lbID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf(
			"Error checking for load balancer (%s) is active: %s", lbID, err)
	}
	d.SetId("")
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
		return fmt.Errorf("Error Getting vpc load balancer pool(%s): %s\n%s", lbPoolID, err, response)
	}
	_, err = isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf(
			"Error checking for load balancer (%s) is active: %s", lbID, err)
	}

	_, err = isWaitForLBPoolActive(sess, lbID, lbPoolID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf(
			"Error checking for load balancer pool (%s) is active: %s", lbPoolID, err)
	}

	deleteLoadBalancerPoolOptions := &vpcv1.DeleteLoadBalancerPoolOptions{
		LoadBalancerID: &lbID,
		ID:             &lbPoolID,
	}
	response, err = sess.DeleteLoadBalancerPool(deleteLoadBalancerPoolOptions)
	if err != nil {
		return fmt.Errorf("Error Deleting Load Balancer Pool : %s\n%s", err, response)
	}
	_, err = isWaitForLBPoolDeleted(sess, lbID, lbPoolID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf(
			"Error checking for load balancer pool (%s) is deleted: %s", lbPoolID, err)
	}

	_, err = isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf(
			"Error checking for load balancer (%s) is active: %s", lbID, err)
	}
	d.SetId("")
	return nil
}

func resourceIBMISLBPoolExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return false, err
	}

	parts, err := idParts(d.Id())
	if err != nil {
		return false, err
	}
	if len(parts) != 2 {
		return false, fmt.Errorf("Incorrect ID %s: ID should be a combination of lbID/lbPoolID", d.Id())
	}

	lbID := parts[0]
	lbPoolID := parts[1]

	if userDetails.generation == 1 {
		exists, err := classicLBPoolExists(d, meta, lbID, lbPoolID)
		return exists, err
	} else {
		exists, err := lbPoolExists(d, meta, lbID, lbPoolID)
		return exists, err
	}
}

func classicLBPoolExists(d *schema.ResourceData, meta interface{}, lbID, lbPoolID string) (bool, error) {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return false, err
	}

	getLoadBalancerPoolOptions := &vpcclassicv1.GetLoadBalancerPoolOptions{
		LoadBalancerID: &lbID,
		ID:             &lbPoolID,
	}
	_, response, err := sess.GetLoadBalancerPool(getLoadBalancerPoolOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("Error getting Load balancer pool: %s\n%s", err, response)
	}
	return true, nil
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
		return false, fmt.Errorf("Error getting Load balancer pool: %s\n%s", err, response)
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
			return nil, "", fmt.Errorf("Error Getting Load Balancer Pool: %s\n%s", err, response)
		}

		if *lbPool.ProvisioningStatus == isLBPoolActive || *lbPool.ProvisioningStatus == isLBPoolFailed {
			return lbPool, isLBPoolActive, nil
		}

		return lbPool, *lbPool.ProvisioningStatus, nil
	}
}

func isWaitForClassicLBPoolActive(sess *vpcclassicv1.VpcClassicV1, lbId, lbPoolId string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for load balancer pool (%s) to be available.", lbPoolId)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{isLBPoolCreatePending, isLBPoolUpdatePending, isLBPoolMaintainancePending},
		Target:     []string{isLBPoolActive, ""},
		Refresh:    isClassicLBPoolRefreshFunc(sess, lbId, lbPoolId),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isClassicLBPoolRefreshFunc(sess *vpcclassicv1.VpcClassicV1, lbId, lbPoolId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		getlbpOptions := &vpcclassicv1.GetLoadBalancerPoolOptions{
			LoadBalancerID: &lbId,
			ID:             &lbPoolId,
		}
		lbPool, response, err := sess.GetLoadBalancerPool(getlbpOptions)
		if err != nil {
			return nil, "", fmt.Errorf("Error Getting Load Balancer Pool: %s\n%s", err, response)
		}

		if *lbPool.ProvisioningStatus == isLBPoolActive || *lbPool.ProvisioningStatus == isLBPoolFailed {
			return lbPool, isLBPoolActive, nil
		}

		return lbPool, *lbPool.ProvisioningStatus, nil
	}
}

func isWaitForClassicLBPoolDeleted(lbc *vpcclassicv1.VpcClassicV1, lbId, lbPoolId string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for  (%s) to be deleted.", lbPoolId)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{isLBPoolUpdatePending, isLBPoolMaintainancePending, isLBPoolDeletePending},
		Target:     []string{isLBPoolDeleteDone, ""},
		Refresh:    isClassicLBPoolDeleteRefreshFunc(lbc, lbId, lbPoolId),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isClassicLBPoolDeleteRefreshFunc(lbc *vpcclassicv1.VpcClassicV1, lbId, lbPoolId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] delete function here")
		getlbpOptions := &vpcclassicv1.GetLoadBalancerPoolOptions{
			LoadBalancerID: &lbId,
			ID:             &lbPoolId,
		}
		lbPool, response, err := lbc.GetLoadBalancerPool(getlbpOptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return lbPool, isLBPoolDeleteDone, nil
			}
			return nil, "", fmt.Errorf("The vpc load balancer pool %s failed to delete: %s\n%s", lbPoolId, err, response)
		}
		return lbPool, isLBPoolDeletePending, nil
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
		log.Printf("[DEBUG] delete function here")
		getlbpOptions := &vpcv1.GetLoadBalancerPoolOptions{
			LoadBalancerID: &lbId,
			ID:             &lbPoolId,
		}
		lbPool, response, err := lbc.GetLoadBalancerPool(getlbpOptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return lbPool, isLBPoolDeleteDone, nil
			}
			return nil, "", fmt.Errorf("The vpc load balancer pool %s failed to delete: %s\n%s", lbPoolId, err, response)
		}
		return lbPool, isLBPoolDeletePending, nil
	}
}
