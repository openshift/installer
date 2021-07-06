// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcclassicv1"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	isLBListenerLBID                = "lb"
	isLBListenerPort                = "port"
	isLBListenerProtocol            = "protocol"
	isLBListenerCertificateInstance = "certificate_instance"
	isLBListenerConnectionLimit     = "connection_limit"
	isLBListenerDefaultPool         = "default_pool"
	isLBListenerStatus              = "status"
	isLBListenerDeleting            = "deleting"
	isLBListenerDeleted             = "done"
	isLBListenerProvisioning        = "provisioning"
	isLBListenerAcceptProxyProtocol = "accept_proxy_protocol"
	isLBListenerProvisioningDone    = "done"
	isLBListenerID                  = "listener_id"
)

func resourceIBMISLBListener() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMISLBListenerCreate,
		Read:     resourceIBMISLBListenerRead,
		Update:   resourceIBMISLBListenerUpdate,
		Delete:   resourceIBMISLBListenerDelete,
		Exists:   resourceIBMISLBListenerExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			isLBListenerLBID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Loadbalancer listener ID",
			},

			isLBListenerPort: {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateLBListenerPort,
				Description:  "Loadbalancer listener port",
			},

			isLBListenerProtocol: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: InvokeValidator("ibm_is_lb_listener", isLBListenerProtocol),
				Description:  "Loadbalancer protocol",
			},

			isLBListenerCertificateInstance: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "certificate instance for the Loadbalancer",
			},

			isLBListenerAcceptProxyProtocol: {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Listener will forward proxy protocol",
			},

			isLBListenerConnectionLimit: {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateLBListenerConnectionLimit,
				Description:  "Connection limit for Loadbalancer",
			},

			isLBListenerDefaultPool: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, o, n string, d *schema.ResourceData) bool {
					if o == "" {
						return false
					}
					// if state file entry and tf file entry matches
					if strings.Compare(n, o) == 0 {
						return true
					}

					if strings.Contains(n, "/") {
						new := strings.Split(n, "/")
						if strings.Compare(new[1], o) == 0 {
							return true
						}
					}

					return false
				},
				Description: "Loadbalancer default pool info",
			},

			isLBListenerStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Loadbalancer listener status",
			},

			isLBListenerID: {
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

func resourceIBMISLBListenerValidator() *ResourceValidator {

	validateSchema := make([]ValidateSchema, 1)
	protocol := "https, http, tcp"
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isLBListenerProtocol,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              protocol})

	ibmISLBListenerResourceValidator := ResourceValidator{ResourceName: "ibm_is_lb_listener", Schema: validateSchema}
	return &ibmISLBListenerResourceValidator
}

func resourceIBMISLBListenerCreate(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] LB Listener create")
	lbID := d.Get(isLBListenerLBID).(string)
	port := int64(d.Get(isLBListenerPort).(int))
	protocol := d.Get(isLBListenerProtocol).(string)
	acceptProxyProtocol := d.Get(isLBListenerAcceptProxyProtocol).(bool)
	var defPool, certificateCRN string
	if pool, ok := d.GetOk(isLBListenerDefaultPool); ok {
		lbPool, err := getPoolId(pool.(string))
		if err != nil {
			return err
		}
		defPool = lbPool
	}

	if crn, ok := d.GetOk(isLBListenerCertificateInstance); ok {
		certificateCRN = crn.(string)
	}

	var connLimit int64

	if limit, ok := d.GetOk(isLBListenerConnectionLimit); ok {
		connLimit = int64(limit.(int))
	}

	isLBKey := "load_balancer_key_" + lbID
	ibmMutexKV.Lock(isLBKey)
	defer ibmMutexKV.Unlock(isLBKey)

	if userDetails.generation == 1 {
		err := classicLBListenerCreate(d, meta, lbID, protocol, defPool, certificateCRN, port, connLimit)
		if err != nil {
			return err
		}
	} else {
		err := lbListenerCreate(d, meta, lbID, protocol, defPool, certificateCRN, port, connLimit, acceptProxyProtocol)
		if err != nil {
			return err
		}
	}
	return resourceIBMISLBListenerRead(d, meta)
}

func classicLBListenerCreate(d *schema.ResourceData, meta interface{}, lbID, protocol, defPool, certificateCRN string, port, connLimit int64) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	options := &vpcclassicv1.CreateLoadBalancerListenerOptions{
		LoadBalancerID: &lbID,
		Port:           &port,
		Protocol:       &protocol,
	}
	if defPool != "" {
		options.DefaultPool = &vpcclassicv1.LoadBalancerPoolIdentity{
			ID: &defPool,
		}
	}
	if certificateCRN != "" {
		options.CertificateInstance = &vpcclassicv1.CertificateInstanceIdentity{
			CRN: &certificateCRN,
		}
	}
	if connLimit > int64(0) {
		options.ConnectionLimit = &connLimit
	}
	_, err = isWaitForClassicLBAvailable(sess, lbID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf(
			"Error checking for load balancer (%s) is active: %s", lbID, err)
	}

	lbListener, response, err := sess.CreateLoadBalancerListener(options)
	if err != nil {
		return fmt.Errorf("Error while creating Load Balanacer Listener err %s\n%s", err, response)
	}
	d.SetId(fmt.Sprintf("%s/%s", lbID, *lbListener.ID))
	_, err = isWaitForClassicLBListenerAvailable(sess, lbID, *lbListener.ID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf(
			"Error waiting for load balancer listener(%s) to become ready: %s", d.Id(), err)
	}
	_, err = isWaitForClassicLBAvailable(sess, lbID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf(
			"Error waiting for load balancer (%s) to become ready: %s", lbID, err)
	}

	log.Printf("[INFO] Load balancer Listener : %s", *lbListener.ID)
	return nil
}

func lbListenerCreate(d *schema.ResourceData, meta interface{}, lbID, protocol, defPool, certificateCRN string, port, connLimit int64, acceptProxyProtocol bool) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	options := &vpcv1.CreateLoadBalancerListenerOptions{
		LoadBalancerID:      &lbID,
		Port:                &port,
		Protocol:            &protocol,
		AcceptProxyProtocol: &acceptProxyProtocol,
	}
	if defPool != "" {
		options.DefaultPool = &vpcv1.LoadBalancerPoolIdentity{
			ID: &defPool,
		}
	}
	if certificateCRN != "" {
		options.CertificateInstance = &vpcv1.CertificateInstanceIdentity{
			CRN: &certificateCRN,
		}
	}
	if connLimit > int64(0) {
		options.ConnectionLimit = &connLimit
	}
	_, err = isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf(
			"Error checking for load balancer (%s) is active: %s", lbID, err)
	}

	lbListener, response, err := sess.CreateLoadBalancerListener(options)
	if err != nil {
		return fmt.Errorf("Error while creating Load Balanacer Listener err %s\n%s", err, response)
	}
	d.SetId(fmt.Sprintf("%s/%s", lbID, *lbListener.ID))
	_, err = isWaitForLBListenerAvailable(sess, lbID, *lbListener.ID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf(
			"Error waiting for load balancer listener(%s) to become ready: %s", d.Id(), err)
	}
	_, err = isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf(
			"Error waiting for load balancer (%s) to become ready: %s", lbID, err)
	}

	log.Printf("[INFO] Load balancer Listener : %s", *lbListener.ID)
	return nil
}

func isWaitForClassicLBListenerAvailable(sess *vpcclassicv1.VpcClassicV1, lbID, lbListenerID string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for load balancer Listener(%s) to be available.", lbListenerID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isLBListenerProvisioning, "create_pending", "update_pending", "maintenance_pending"},
		Target:     []string{isLBListenerProvisioningDone, ""},
		Refresh:    isClassicLBListenerRefreshFunc(sess, lbID, lbListenerID),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isClassicLBListenerRefreshFunc(sess *vpcclassicv1.VpcClassicV1, lbID, lbListenerID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		getLoadBalancerListenerOptions := &vpcclassicv1.GetLoadBalancerListenerOptions{
			LoadBalancerID: &lbID,
			ID:             &lbListenerID,
		}
		lblis, response, err := sess.GetLoadBalancerListener(getLoadBalancerListenerOptions)
		if err != nil {
			return nil, "", fmt.Errorf("Error Getting Load Balancer Listener: %s\n%s", err, response)
		}

		if *lblis.ProvisioningStatus == "active" || *lblis.ProvisioningStatus == "failed" {
			return lblis, isLBListenerProvisioningDone, nil
		}

		return lblis, *lblis.ProvisioningStatus, nil
	}
}

func isWaitForLBListenerAvailable(sess *vpcv1.VpcV1, lbID, lbListenerID string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for load balancer Listener(%s) to be available.", lbListenerID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isLBListenerProvisioning, "create_pending", "update_pending", "maintenance_pending"},
		Target:     []string{isLBListenerProvisioningDone, ""},
		Refresh:    isLBListenerRefreshFunc(sess, lbID, lbListenerID),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isLBListenerRefreshFunc(sess *vpcv1.VpcV1, lbID, lbListenerID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		getLoadBalancerListenerOptions := &vpcv1.GetLoadBalancerListenerOptions{
			LoadBalancerID: &lbID,
			ID:             &lbListenerID,
		}
		lblis, response, err := sess.GetLoadBalancerListener(getLoadBalancerListenerOptions)
		if err != nil {
			return nil, "", fmt.Errorf("Error Getting Load Balancer Listener: %s\n%s", err, response)
		}

		if *lblis.ProvisioningStatus == "active" || *lblis.ProvisioningStatus == "failed" {
			return lblis, isLBListenerProvisioningDone, nil
		}

		return lblis, *lblis.ProvisioningStatus, nil
	}
}

func resourceIBMISLBListenerRead(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}

	lbID := parts[0]
	lbListenerID := parts[1]

	if userDetails.generation == 1 {
		err := classicLBListenerGet(d, meta, lbID, lbListenerID)
		if err != nil {
			return err
		}
	} else {
		err := lbListenerGet(d, meta, lbID, lbListenerID)
		if err != nil {
			return err
		}
	}

	return nil
}

func classicLBListenerGet(d *schema.ResourceData, meta interface{}, lbID, lbListenerID string) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	getLoadBalancerListenerOptions := &vpcclassicv1.GetLoadBalancerListenerOptions{
		LoadBalancerID: &lbID,
		ID:             &lbListenerID,
	}
	lbListener, response, err := sess.GetLoadBalancerListener(getLoadBalancerListenerOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Getting Load Balancer Listener : %s\n%s", err, response)
	}
	d.Set(isLBListenerLBID, lbID)
	d.Set(isLBListenerPort, *lbListener.Port)
	d.Set(isLBListenerProtocol, *lbListener.Protocol)
	d.Set(isLBListenerID, lbListenerID)
	if lbListener.DefaultPool != nil {
		d.Set(isLBListenerDefaultPool, *lbListener.DefaultPool.ID)
	}
	if lbListener.CertificateInstance != nil {
		d.Set(isLBListenerCertificateInstance, *lbListener.CertificateInstance.CRN)
	}
	if lbListener.ConnectionLimit != nil {
		d.Set(isLBListenerConnectionLimit, *lbListener.ConnectionLimit)
	}
	d.Set(isLBListenerStatus, *lbListener.ProvisioningStatus)
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

func lbListenerGet(d *schema.ResourceData, meta interface{}, lbID, lbListenerID string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	getLoadBalancerListenerOptions := &vpcv1.GetLoadBalancerListenerOptions{
		LoadBalancerID: &lbID,
		ID:             &lbListenerID,
	}
	lbListener, response, err := sess.GetLoadBalancerListener(getLoadBalancerListenerOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Getting Load Balancer Listener : %s\n%s", err, response)
	}
	d.Set(isLBListenerLBID, lbID)
	d.Set(isLBListenerPort, *lbListener.Port)
	d.Set(isLBListenerProtocol, *lbListener.Protocol)
	d.Set(isLBListenerAcceptProxyProtocol, *lbListener.AcceptProxyProtocol)
	d.Set(isLBListenerID, lbListenerID)
	if lbListener.DefaultPool != nil {
		d.Set(isLBListenerDefaultPool, *lbListener.DefaultPool.ID)
	}
	if lbListener.CertificateInstance != nil {
		d.Set(isLBListenerCertificateInstance, *lbListener.CertificateInstance.CRN)
	}
	if lbListener.ConnectionLimit != nil {
		d.Set(isLBListenerConnectionLimit, *lbListener.ConnectionLimit)
	}
	d.Set(isLBListenerStatus, *lbListener.ProvisioningStatus)
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

func resourceIBMISLBListenerUpdate(d *schema.ResourceData, meta interface{}) error {

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}

	lbID := parts[0]
	lbListenerID := parts[1]

	if userDetails.generation == 1 {
		err := classicLBListenerUpdate(d, meta, lbID, lbListenerID)
		if err != nil {
			return err
		}
	} else {
		err := lbListenerUpdate(d, meta, lbID, lbListenerID)
		if err != nil {
			return err
		}
	}

	return resourceIBMISLBListenerRead(d, meta)
}

func classicLBListenerUpdate(d *schema.ResourceData, meta interface{}, lbID, lbListenerID string) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	hasChanged := false
	var certificateInstance, defPool, protocol string
	var connLimit, port int64
	updateLoadBalancerListenerOptions := &vpcclassicv1.UpdateLoadBalancerListenerOptions{
		LoadBalancerID: &lbID,
		ID:             &lbListenerID,
	}
	loadBalancerListenerPatchModel := &vpcclassicv1.LoadBalancerListenerPatch{}

	if d.HasChange(isLBListenerCertificateInstance) {
		certificateInstance = d.Get(isLBListenerCertificateInstance).(string)
		loadBalancerListenerPatchModel.CertificateInstance = &vpcclassicv1.CertificateInstanceIdentity{
			CRN: &certificateInstance,
		}
		hasChanged = true
	}

	if d.HasChange(isLBListenerDefaultPool) {
		lbpool, err := getPoolId(d.Get(isLBListenerDefaultPool).(string))
		if err != nil {
			return err
		}
		defPool = lbpool
		loadBalancerListenerPatchModel.DefaultPool = &vpcclassicv1.LoadBalancerPoolIdentity{
			ID: &defPool,
		}
		hasChanged = true
	}
	if d.HasChange(isLBListenerPort) {
		port = int64(d.Get(isLBListenerPort).(int))
		loadBalancerListenerPatchModel.Port = &port
		hasChanged = true
	}

	if d.HasChange(isLBListenerProtocol) {
		protocol = d.Get(isLBListenerProtocol).(string)
		loadBalancerListenerPatchModel.Protocol = &protocol
		hasChanged = true
	}

	if d.HasChange(isLBListenerConnectionLimit) {
		connLimit = int64(d.Get(isLBListenerConnectionLimit).(int))
		loadBalancerListenerPatchModel.ConnectionLimit = &connLimit
		hasChanged = true
	}

	if hasChanged {
		loadBalancerListenerPatch, err := loadBalancerListenerPatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("Error calling asPatch for LoadBalancerListenerPatch: %s", err)
		}
		updateLoadBalancerListenerOptions.LoadBalancerListenerPatch = loadBalancerListenerPatch

		isLBKey := "load_balancer_key_" + lbID
		ibmMutexKV.Lock(isLBKey)
		defer ibmMutexKV.Unlock(isLBKey)

		_, err = isWaitForClassicLBAvailable(sess, lbID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf(
				"Error checking for load balancer (%s) is active: %s", lbID, err)
		}
		_, response, err := sess.UpdateLoadBalancerListener(updateLoadBalancerListenerOptions)
		if err != nil {
			return fmt.Errorf("Error Updating Load Balancer Listener : %s\n%s", err, response)
		}

		_, err = isWaitForClassicLBListenerAvailable(sess, lbID, lbListenerID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf(
				"Error waiting for load balancer listener(%s) to become ready: %s", d.Id(), err)
		}

		_, err = isWaitForClassicLBAvailable(sess, lbID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf(
				"Error waiting for load balancer (%s) to become ready: %s", lbID, err)
		}
	}
	return nil
}

func lbListenerUpdate(d *schema.ResourceData, meta interface{}, lbID, lbListenerID string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	hasChanged := false
	var certificateInstance, defPool, protocol string
	var connLimit, port int64
	updateLoadBalancerListenerOptions := &vpcv1.UpdateLoadBalancerListenerOptions{
		LoadBalancerID: &lbID,
		ID:             &lbListenerID,
	}

	loadBalancerListenerPatchModel := &vpcv1.LoadBalancerListenerPatch{}

	if d.HasChange(isLBListenerCertificateInstance) {
		certificateInstance = d.Get(isLBListenerCertificateInstance).(string)
		loadBalancerListenerPatchModel.CertificateInstance = &vpcv1.CertificateInstanceIdentity{
			CRN: &certificateInstance,
		}
		hasChanged = true
	}

	if d.HasChange(isLBListenerDefaultPool) {
		lbpool, err := getPoolId(d.Get(isLBListenerDefaultPool).(string))
		if err != nil {
			return err
		}
		defPool = lbpool
		loadBalancerListenerPatchModel.DefaultPool = &vpcv1.LoadBalancerPoolIdentity{
			ID: &defPool,
		}
		hasChanged = true
	}
	if d.HasChange(isLBListenerPort) {
		port = int64(d.Get(isLBListenerPort).(int))
		loadBalancerListenerPatchModel.Port = &port
		hasChanged = true
	}

	if d.HasChange(isLBListenerProtocol) {
		protocol = d.Get(isLBListenerProtocol).(string)
		loadBalancerListenerPatchModel.Protocol = &protocol
		hasChanged = true
	}

	if d.HasChange(isLBListenerAcceptProxyProtocol) {
		acceptProxyProtocol := d.Get(isLBListenerAcceptProxyProtocol).(bool)
		loadBalancerListenerPatchModel.AcceptProxyProtocol = &acceptProxyProtocol
		hasChanged = true
	}

	if d.HasChange(isLBListenerConnectionLimit) {
		connLimit = int64(d.Get(isLBListenerConnectionLimit).(int))
		loadBalancerListenerPatchModel.ConnectionLimit = &connLimit
		hasChanged = true
	}

	if hasChanged {
		loadBalancerListenerPatch, err := loadBalancerListenerPatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("Error calling asPatch for LoadBalancerListenerPatch: %s", err)
		}
		updateLoadBalancerListenerOptions.LoadBalancerListenerPatch = loadBalancerListenerPatch

		isLBKey := "load_balancer_key_" + lbID
		ibmMutexKV.Lock(isLBKey)
		defer ibmMutexKV.Unlock(isLBKey)

		_, err = isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf(
				"Error checking for load balancer (%s) is active: %s", lbID, err)
		}
		_, response, err := sess.UpdateLoadBalancerListener(updateLoadBalancerListenerOptions)
		if err != nil {
			return fmt.Errorf("Error Updating Load Balancer Listener : %s\n%s", err, response)
		}

		_, err = isWaitForLBListenerAvailable(sess, lbID, lbListenerID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf(
				"Error waiting for load balancer listener(%s) to become ready: %s", d.Id(), err)
		}

		_, err = isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf(
				"Error waiting for load balancer (%s) to become ready: %s", lbID, err)
		}
	}
	return nil
}

func resourceIBMISLBListenerDelete(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}

	lbID := parts[0]
	lbListenerID := parts[1]

	isLBKey := "load_balancer_key_" + lbID
	ibmMutexKV.Lock(isLBKey)
	defer ibmMutexKV.Unlock(isLBKey)

	if userDetails.generation == 1 {
		err := classicLBListenerDelete(d, meta, lbID, lbListenerID)
		if err != nil {
			return err
		}
	} else {
		err := lbListenerDelete(d, meta, lbID, lbListenerID)
		if err != nil {
			return err
		}
	}
	return nil
}

func classicLBListenerDelete(d *schema.ResourceData, meta interface{}, lbID, lbListenerID string) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	getLoadBalancerListenerOptions := &vpcclassicv1.GetLoadBalancerListenerOptions{
		LoadBalancerID: &lbID,
		ID:             &lbListenerID,
	}
	_, response, err := sess.GetLoadBalancerListener(getLoadBalancerListenerOptions)

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Getting vpc load balancer listener(%s): %s\n%s", lbListenerID, err, response)
	}
	_, err = isWaitForClassicLBAvailable(sess, lbID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf(
			"Error checking for load balancer (%s) is active: %s", lbID, err)
	}
	deleteLoadBalancerListenerOptions := &vpcclassicv1.DeleteLoadBalancerListenerOptions{
		LoadBalancerID: &lbID,
		ID:             &lbListenerID,
	}
	response, err = sess.DeleteLoadBalancerListener(deleteLoadBalancerListenerOptions)
	if err != nil {
		return fmt.Errorf("Error Deleting Load Balancer Pool : %s\n%s", err, response)
	}
	_, err = isWaitForClassicLBListenerDeleted(sess, lbID, lbListenerID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}
	_, err = isWaitForClassicLBAvailable(sess, lbID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf(
			"Error waiting for load balancer (%s) to be active: %s", lbID, err)
	}

	d.SetId("")
	return nil
}

func lbListenerDelete(d *schema.ResourceData, meta interface{}, lbID, lbListenerID string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	getLoadBalancerListenerOptions := &vpcv1.GetLoadBalancerListenerOptions{
		LoadBalancerID: &lbID,
		ID:             &lbListenerID,
	}
	_, response, err := sess.GetLoadBalancerListener(getLoadBalancerListenerOptions)

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Getting vpc load balancer listener(%s): %s\n%s", lbListenerID, err, response)
	}
	_, err = isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf(
			"Error checking for load balancer (%s) is active: %s", lbID, err)
	}
	deleteLoadBalancerListenerOptions := &vpcv1.DeleteLoadBalancerListenerOptions{
		LoadBalancerID: &lbID,
		ID:             &lbListenerID,
	}
	response, err = sess.DeleteLoadBalancerListener(deleteLoadBalancerListenerOptions)
	if err != nil {
		return fmt.Errorf("Error Deleting Load Balancer Pool : %s\n%s", err, response)
	}
	_, err = isWaitForLBListenerDeleted(sess, lbID, lbListenerID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}
	_, err = isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf(
			"Error waiting for load balancer (%s) to be active: %s", lbID, err)
	}

	d.SetId("")
	return nil
}

func isWaitForClassicLBListenerDeleted(lbc *vpcclassicv1.VpcClassicV1, lbID, lbListenerID string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for  (%s) to be deleted.", lbListenerID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isLBListenerDeleting, "delete_pending"},
		Target:     []string{isLBListenerDeleted, ""},
		Refresh:    isClassicLBListenerDeleteRefreshFunc(lbc, lbID, lbListenerID),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isClassicLBListenerDeleteRefreshFunc(lbc *vpcclassicv1.VpcClassicV1, lbID, lbListenerID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getLoadBalancerListenerOptions := &vpcclassicv1.GetLoadBalancerListenerOptions{
			LoadBalancerID: &lbID,
			ID:             &lbListenerID,
		}
		lbLis, response, err := lbc.GetLoadBalancerListener(getLoadBalancerListenerOptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return lbLis, isLBListenerDeleted, nil
			}
			return nil, "", fmt.Errorf("The vpc load balancer listener %s failed to delete: %s\n%s", lbListenerID, err, response)
		}
		return lbLis, isLBListenerDeleting, nil
	}
}

func isWaitForLBListenerDeleted(lbc *vpcv1.VpcV1, lbID, lbListenerID string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for  (%s) to be deleted.", lbListenerID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isLBListenerDeleting, "delete_pending"},
		Target:     []string{isLBListenerDeleted, ""},
		Refresh:    isLBListenerDeleteRefreshFunc(lbc, lbID, lbListenerID),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isLBListenerDeleteRefreshFunc(lbc *vpcv1.VpcV1, lbID, lbListenerID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getLoadBalancerListenerOptions := &vpcv1.GetLoadBalancerListenerOptions{
			LoadBalancerID: &lbID,
			ID:             &lbListenerID,
		}
		lbLis, response, err := lbc.GetLoadBalancerListener(getLoadBalancerListenerOptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return lbLis, isLBListenerDeleted, nil
			}
			return nil, "", fmt.Errorf("The vpc load balancer listener %s failed to delete: %s\n%s", lbListenerID, err, response)
		}
		return lbLis, isLBListenerDeleting, nil
	}
}

func resourceIBMISLBListenerExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return false, err
	}

	parts, err := idParts(d.Id())
	if err != nil {
		return false, err
	}
	if len(parts) != 2 {
		return false, fmt.Errorf("Incorrect ID %s: ID should be a combination of lbID/lbListenerID", d.Id())
	}
	lbID := parts[0]
	lbListenerID := parts[1]

	if userDetails.generation == 1 {
		exists, err := classicLBListenerExists(d, meta, lbID, lbListenerID)
		return exists, err
	} else {
		exists, err := lbListenerExists(d, meta, lbID, lbListenerID)
		return exists, err
	}
}

func classicLBListenerExists(d *schema.ResourceData, meta interface{}, lbID, lbListenerID string) (bool, error) {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return false, err
	}

	getLoadBalancerListenerOptions := &vpcclassicv1.GetLoadBalancerListenerOptions{
		LoadBalancerID: &lbID,
		ID:             &lbListenerID,
	}
	_, response, err := sess.GetLoadBalancerListener(getLoadBalancerListenerOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("Error getting Load balancer Listener: %s\n%s", err, response)
	}
	return true, nil
}

func lbListenerExists(d *schema.ResourceData, meta interface{}, lbID, lbListenerID string) (bool, error) {
	sess, err := vpcClient(meta)
	if err != nil {
		return false, err
	}

	getLoadBalancerListenerOptions := &vpcv1.GetLoadBalancerListenerOptions{
		LoadBalancerID: &lbID,
		ID:             &lbListenerID,
	}
	_, response, err := sess.GetLoadBalancerListener(getLoadBalancerListenerOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("Error getting Load balancer Listener: %s\n%s", err, response)
	}
	return true, nil
}
