// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isLBListenerLBID                    = "lb"
	isLBListenerPort                    = "port"
	isLBListenerPortMin                 = "port_min"
	isLBListenerPortMax                 = "port_max"
	isLBListenerProtocol                = "protocol"
	isLBListenerCertificateInstance     = "certificate_instance"
	isLBListenerConnectionLimit         = "connection_limit"
	isLBListenerDefaultPool             = "default_pool"
	isLBListenerStatus                  = "status"
	isLBListenerDeleting                = "deleting"
	isLBListenerDeleted                 = "done"
	isLBListenerProvisioning            = "provisioning"
	isLBListenerAcceptProxyProtocol     = "accept_proxy_protocol"
	isLBListenerProvisioningDone        = "done"
	isLBListenerID                      = "listener_id"
	isLBListenerHTTPSRedirectListener   = "https_redirect_listener"
	isLBListenerHTTPSRedirectStatusCode = "https_redirect_status_code"
	isLBListenerHTTPSRedirectURI        = "https_redirect_uri"
)

func ResourceIBMISLBListener() *schema.Resource {
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
				Optional:     true,
				ValidateFunc: validate.ValidateLBListenerPort,
				Computed:     true,
				Description:  "Loadbalancer listener port",
			},
			isLBListenerPortMin: {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validate.ValidateLBListenerPort,
				Computed:     true,
				Description:  "The inclusive lower bound of the range of ports used by this listener. Only load balancers in the `network` family support more than one port per listener.",
			},

			isLBListenerPortMax: {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validate.ValidateLBListenerPort,
				Computed:     true,
				Description:  "The inclusive upper bound of the range of ports used by this listener. Only load balancers in the `network` family support more than one port per listener",
			},

			isLBListenerProtocol: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_lb_listener", isLBListenerProtocol),
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

			isLBListenerHTTPSRedirectStatusCode: {
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{isLBListenerHTTPSRedirectListener},
				Description:  "The HTTP status code to be returned in the redirect response",
			},

			isLBListenerHTTPSRedirectURI: {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{isLBListenerHTTPSRedirectStatusCode, isLBListenerHTTPSRedirectListener},
				Description:  "Target URI where traffic will be redirected",
			},

			isLBListenerHTTPSRedirectListener: {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{isLBListenerHTTPSRedirectStatusCode},
				Description:  "ID of the listener that will be set as http redirect target",
			},

			isLBListenerConnectionLimit: {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validate.ValidateLBListenerConnectionLimit,
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

			flex.RelatedCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the LB resource",
			},
		},
	}
}

func ResourceIBMISLBListenerValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)
	protocol := "https, http, tcp, udp"
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isLBListenerProtocol,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              protocol})

	ibmISLBListenerResourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_lb_listener", Schema: validateSchema}
	return &ibmISLBListenerResourceValidator
}

func resourceIBMISLBListenerCreate(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] LB Listener create")
	lbID := d.Get(isLBListenerLBID).(string)
	port := int64(d.Get(isLBListenerPort).(int))
	portMin := int64(d.Get(isLBListenerPortMin).(int))
	portMax := int64(d.Get(isLBListenerPortMax).(int))
	protocol := d.Get(isLBListenerProtocol).(string)
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

	var httpStatusCode int64

	if statusCode, ok := d.GetOk(isLBListenerHTTPSRedirectStatusCode); ok {
		httpStatusCode = int64(statusCode.(int))
	}

	var uri string

	if redirecturi, ok := d.GetOk(isLBListenerHTTPSRedirectURI); ok {
		uri = redirecturi.(string)
	}

	var listener string

	if redirectListener, ok := d.GetOk(isLBListenerHTTPSRedirectListener); ok {
		listener = redirectListener.(string)
	}

	isLBKey := "load_balancer_key_" + lbID
	conns.IbmMutexKV.Lock(isLBKey)
	defer conns.IbmMutexKV.Unlock(isLBKey)

	err := lbListenerCreate(d, meta, lbID, protocol, defPool, certificateCRN, listener, uri, port, portMin, portMax, connLimit, httpStatusCode)
	if err != nil {
		return err
	}

	return resourceIBMISLBListenerRead(d, meta)
}

func lbListenerCreate(d *schema.ResourceData, meta interface{}, lbID, protocol, defPool, certificateCRN, listener, uri string, port, portMin, portMax, connLimit, httpStatusCode int64) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	options := &vpcv1.CreateLoadBalancerListenerOptions{
		LoadBalancerID: &lbID,
		Protocol:       &protocol,
	}

	getlboptions := &vpcv1.GetLoadBalancerOptions{
		ID: &lbID,
	}
	lb, response, err := sess.GetLoadBalancer(getlboptions)

	if err != nil || lb == nil {
		return fmt.Errorf("[ERROR] Error getting Load Balancer : %s\n%s", err, response)
	}
	if lb != nil && *lb.RouteMode && lb.Profile != nil && *lb.Profile.Name == "network-fixed" {
		if portMin > 0 && portMax > 0 && portMin != 1 && portMax != 65535 {
			return fmt.Errorf("[ERROR] Only acceptable value for port_min is 1 and port_max is 65535 for route_mode enabled private network load balancer")
		}
		pmin := int64(1)
		pmax := int64(65535)

		options.PortMin = &pmin
		options.PortMax = &pmax
	} else if lb != nil && lb.Profile != nil {
		if strings.EqualFold(*lb.Profile.Family, "network") && *lb.IsPublic {
			if port == 0 && (portMin == 0 || portMax == 0) {
				return fmt.Errorf(
					"[ERROR] Error port_min(%d)/port_max(%d) for public network load balancer(%s) needs to be in between 1-65335", portMin, portMax, lbID)
			} else {
				if port != 0 {
					options.Port = &port
				} else {
					options.PortMin = &portMin
					options.PortMax = &portMax
				}
			}
		} else if portMin != portMax {
			return fmt.Errorf("[ERROR] Listener port_min and port_max values have to be equal for ALB and private NLB (excluding route mode)")
		} else {
			if port != 0 && (portMin == 0 || port == portMin) {
				options.Port = &port
			} else {
				options.PortMin = &portMin
				options.PortMax = &portMax
			}
		}
	}

	if app, ok := d.GetOk(isLBListenerAcceptProxyProtocol); ok {
		acceptProxyProtocol := app.(bool)
		options.AcceptProxyProtocol = &acceptProxyProtocol
	}

	if defPool != "" {
		options.DefaultPool = &vpcv1.LoadBalancerPoolIdentity{
			ID: &defPool,
		}
	}

	if listener != "" {
		httpsRedirect := &vpcv1.LoadBalancerListenerHTTPSRedirectPrototype{
			HTTPStatusCode: &httpStatusCode,
			Listener: &vpcv1.LoadBalancerListenerIdentity{
				ID: &listener,
			},
		}
		if uri != "" {
			httpsRedirect.URI = &uri
		}
		options.HTTPSRedirect = httpsRedirect
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
		return fmt.Errorf("[ERROR] Error checking for load balancer (%s) is active: %s", lbID, err)
	}

	lbListener, response, err := sess.CreateLoadBalancerListener(options)
	if err != nil {
		return fmt.Errorf("[ERROR] Error while creating Load Balanacer Listener err %s\n%s", err, response)
	}
	d.SetId(fmt.Sprintf("%s/%s", lbID, *lbListener.ID))
	_, err = isWaitForLBListenerAvailable(sess, lbID, *lbListener.ID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("[ERROR] Error waiting for load balancer listener(%s) to become ready: %s", d.Id(), err)
	}
	_, err = isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("[ERROR] Error waiting for load balancer (%s) to become ready: %s", lbID, err)
	}

	log.Printf("[INFO] Load balancer Listener : %s", *lbListener.ID)
	return nil
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
			return nil, "", fmt.Errorf("[ERROR] Error Getting Load Balancer Listener: %s\n%s", err, response)
		}

		if *lblis.ProvisioningStatus == "active" || *lblis.ProvisioningStatus == "failed" {
			return lblis, isLBListenerProvisioningDone, nil
		}

		return lblis, *lblis.ProvisioningStatus, nil
	}
}

func resourceIBMISLBListenerRead(d *schema.ResourceData, meta interface{}) error {

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}

	lbID := parts[0]
	lbListenerID := parts[1]

	err = lbListenerGet(d, meta, lbID, lbListenerID)
	if err != nil {
		return err
	}

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
		return fmt.Errorf("[ERROR] Error Getting Load Balancer Listener : %s\n%s", err, response)
	}
	d.Set(isLBListenerLBID, lbID)
	if lbListener.Port != nil {
		d.Set(isLBListenerPort, *lbListener.Port)
	}
	if lbListener.PortMin != nil {
		d.Set(isLBListenerPortMin, *lbListener.PortMin)
	}
	if lbListener.PortMax != nil {
		d.Set(isLBListenerPortMax, *lbListener.PortMax)
	}
	d.Set(isLBListenerProtocol, *lbListener.Protocol)
	d.Set(isLBListenerAcceptProxyProtocol, *lbListener.AcceptProxyProtocol)
	d.Set(isLBListenerID, lbListenerID)
	if lbListener.DefaultPool != nil {
		d.Set(isLBListenerDefaultPool, *lbListener.DefaultPool.ID)
	}
	if lbListener.HTTPSRedirect != nil {
		d.Set(isLBListenerHTTPSRedirectStatusCode, *lbListener.HTTPSRedirect.HTTPStatusCode)
		d.Set(isLBListenerHTTPSRedirectListener, *lbListener.HTTPSRedirect.Listener.ID)
		if lbListener.HTTPSRedirect.URI != nil {
			d.Set(isLBListenerHTTPSRedirectURI, *lbListener.HTTPSRedirect.URI)
		}
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
		return fmt.Errorf("[ERROR] Error Getting Load Balancer : %s\n%s", err, response)
	}
	d.Set(flex.RelatedCRN, *lb.CRN)
	return nil
}

func resourceIBMISLBListenerUpdate(d *schema.ResourceData, meta interface{}) error {

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}

	lbID := parts[0]
	lbListenerID := parts[1]

	err = lbListenerUpdate(d, meta, lbID, lbListenerID)
	if err != nil {
		return err
	}

	return resourceIBMISLBListenerRead(d, meta)
}

func lbListenerUpdate(d *schema.ResourceData, meta interface{}, lbID, lbListenerID string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	hasChanged := false
	var certificateInstance, defPool, protocol, listener, uri string
	var connLimit, port, httpStatusCode int64
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
	httpsRedirectRemoved := false
	httpsURIRemoved := false
	if d.HasChange(isLBListenerHTTPSRedirectListener) || d.HasChange(isLBListenerHTTPSRedirectURI) || d.HasChange(isLBListenerHTTPSRedirectStatusCode) {
		hasChanged = true
		listener = d.Get(isLBListenerHTTPSRedirectListener).(string)
		httpStatusCode = int64(d.Get(isLBListenerHTTPSRedirectStatusCode).(int))
		uri = d.Get(isLBListenerHTTPSRedirectURI).(string)
		if listener == "" {
			httpsRedirectRemoved = true
		} else {
			HTTPSRedirect := &vpcv1.LoadBalancerListenerHTTPSRedirectPatch{
				HTTPStatusCode: &httpStatusCode,
				Listener:       &vpcv1.LoadBalancerListenerIdentityByID{ID: &listener},
			}
			if d.HasChange(isLBListenerHTTPSRedirectURI) {
				if uri == "" {
					HTTPSRedirect.URI = nil
					httpsURIRemoved = true
				} else {
					HTTPSRedirect.URI = &uri
				}
			}

			loadBalancerListenerPatchModel.HTTPSRedirect = HTTPSRedirect
		}
	}
	if _, ok := d.GetOk(isLBListenerPort); ok && d.HasChange(isLBListenerPort) {
		port = int64(d.Get(isLBListenerPort).(int))
		loadBalancerListenerPatchModel.Port = &port
		hasChanged = true
	}
	if d.HasChange(isLBListenerPortMin) {
		portMin := int64(d.Get(isLBListenerPortMin).(int))
		loadBalancerListenerPatchModel.PortMin = &portMin
		hasChanged = true
	}
	if d.HasChange(isLBListenerPortMax) {
		portMax := int64(d.Get(isLBListenerPortMax).(int))
		loadBalancerListenerPatchModel.PortMax = &portMax
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
			return fmt.Errorf("[ERROR] Error calling asPatch for LoadBalancerListenerPatch: %s", err)
		}
		if httpsRedirectRemoved {
			loadBalancerListenerPatch["https_redirect"] = nil
		}
		if httpsURIRemoved {
			loadBalancerListenerPatch["https_redirect"].(map[string]interface{})["uri"] = nil
		}
		updateLoadBalancerListenerOptions.LoadBalancerListenerPatch = loadBalancerListenerPatch

		isLBKey := "load_balancer_key_" + lbID
		conns.IbmMutexKV.Lock(isLBKey)
		defer conns.IbmMutexKV.Unlock(isLBKey)

		_, err = isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf(
				"Error checking for load balancer (%s) is active: %s", lbID, err)
		}
		_, response, err := sess.UpdateLoadBalancerListener(updateLoadBalancerListenerOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Updating Load Balancer Listener : %s\n%s", err, response)
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

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}

	lbID := parts[0]
	lbListenerID := parts[1]

	isLBKey := "load_balancer_key_" + lbID
	conns.IbmMutexKV.Lock(isLBKey)
	defer conns.IbmMutexKV.Unlock(isLBKey)

	err = lbListenerDelete(d, meta, lbID, lbListenerID)
	if err != nil {
		return err
	}

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
		return fmt.Errorf("[ERROR] Error Getting vpc load balancer listener(%s): %s\n%s", lbListenerID, err, response)
	}
	_, err = isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf("[ERROR] Error checking for load balancer (%s) is active: %s", lbID, err)
	}
	deleteLoadBalancerListenerOptions := &vpcv1.DeleteLoadBalancerListenerOptions{
		LoadBalancerID: &lbID,
		ID:             &lbListenerID,
	}
	response, err = sess.DeleteLoadBalancerListener(deleteLoadBalancerListenerOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Deleting Load Balancer Pool : %s\n%s", err, response)
	}
	_, err = isWaitForLBListenerDeleted(sess, lbID, lbListenerID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}
	_, err = isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf("[ERROR] Error waiting for load balancer (%s) to be active: %s", lbID, err)
	}

	d.SetId("")
	return nil
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
			return nil, "", fmt.Errorf("[ERROR] The vpc load balancer listener %s failed to delete: %s\n%s", lbListenerID, err, response)
		}
		return lbLis, isLBListenerDeleting, nil
	}
}

func resourceIBMISLBListenerExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return false, err
	}
	if len(parts) != 2 {
		return false, fmt.Errorf("[ERROR] Incorrect ID %s: ID should be a combination of lbID/lbListenerID", d.Id())
	}
	lbID := parts[0]
	lbListenerID := parts[1]

	exists, err := lbListenerExists(d, meta, lbID, lbListenerID)
	return exists, err

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
		return false, fmt.Errorf("[ERROR] Error getting Load balancer Listener: %s\n%s", err, response)
	}
	return true, nil
}
