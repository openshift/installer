// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
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
	isLBListenerIdleConnectionTimeout   = "idle_connection_timeout"
)

func ResourceIBMISLBListener() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMISLBListenerCreate,
		ReadContext:   resourceIBMISLBListenerRead,
		UpdateContext: resourceIBMISLBListenerUpdate,
		DeleteContext: resourceIBMISLBListenerDelete,
		Exists:        resourceIBMISLBListenerExists,
		Importer:      &schema.ResourceImporter{},

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

			isLBListenerIdleConnectionTimeout: {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				Description:  "idle connection timeout of listener",
				ValidateFunc: validate.InvokeValidator("ibm_is_lb_listener", isLBListenerIdleConnectionTimeout),
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
				Type:          schema.TypeInt,
				Optional:      true,
				RequiredWith:  []string{isLBListenerHTTPSRedirectListener},
				ConflictsWith: []string{"https_redirect"},
				Deprecated:    "Please use the argument 'https_redirect'",
				Description:   "The HTTP status code to be returned in the redirect response",
			},

			isLBListenerHTTPSRedirectURI: {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"https_redirect"},
				Deprecated:    "Please use the argument 'https_redirect'",
				RequiredWith:  []string{isLBListenerHTTPSRedirectStatusCode, isLBListenerHTTPSRedirectListener},
				Description:   "Target URI where traffic will be redirected",
			},

			isLBListenerHTTPSRedirectListener: {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"https_redirect"},
				Deprecated:    "Please use the argument 'https_redirect'",
				RequiredWith:  []string{isLBListenerHTTPSRedirectStatusCode},
				Description:   "ID of the listener that will be set as http redirect target",
			},
			"https_redirect": &schema.Schema{
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: []string{"https_redirect_status_code", "https_redirect_uri", "https_redirect_listener"},
				Description:   "If present, the target listener that requests are redirected to.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"http_status_code": &schema.Schema{
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The HTTP status code for this redirect.",
						},
						"listener": &schema.Schema{
							Type:     schema.TypeList,
							MinItems: 1,
							MaxItems: 1,
							Required: true,
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
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The redirect relative target URI.",
						},
					},
				},
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
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isLBListenerIdleConnectionTimeout,
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			MinValue:                   "50",
			MaxValue:                   "7200"})
	ibmISLBListenerResourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_lb_listener", Schema: validateSchema}
	return &ibmISLBListenerResourceValidator
}

func resourceIBMISLBListenerCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

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
			diag.FromErr(err)
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

	return resourceIBMISLBListenerRead(context, d, meta)
}

func lbListenerCreate(d *schema.ResourceData, meta interface{}, lbID, protocol, defPool, certificateCRN, listener, uri string, port, portMin, portMax, connLimit, httpStatusCode int64) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
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
		return diag.FromErr(fmt.Errorf("[ERROR] Error getting Load Balancer : %s\n%s", err, response))
	}
	if lb != nil && *lb.RouteMode && lb.Profile != nil && *lb.Profile.Name == "network-fixed" {
		if portMin > 0 && portMax > 0 && portMin != 1 && portMax != 65535 {
			return diag.FromErr(fmt.Errorf("[ERROR] Only acceptable value for port_min is 1 and port_max is 65535 for route_mode enabled private network load balancer"))
		}
		pmin := int64(1)
		pmax := int64(65535)

		options.PortMin = &pmin
		options.PortMax = &pmax
	} else if lb != nil && lb.Profile != nil {
		if strings.EqualFold(*lb.Profile.Family, "network") && *lb.IsPublic {
			if port == 0 && (portMin == 0 || portMax == 0) {
				return diag.FromErr(fmt.Errorf(
					"[ERROR] Error port_min(%d)/port_max(%d) for public network load balancer(%s) needs to be in between 1-65335", portMin, portMax, lbID))
			} else {
				if port != 0 {
					options.Port = &port
				} else {
					options.PortMin = &portMin
					options.PortMax = &portMax
				}
			}
		} else if portMin != portMax {
			return diag.FromErr(fmt.Errorf("[ERROR] Listener port_min and port_max values have to be equal for ALB and private NLB (excluding route mode)"))
		} else {
			if port != 0 && (portMin == 0 || port == portMin) {
				options.Port = &port
			} else {
				options.PortMin = &portMin
				options.PortMax = &portMax
			}
		}
	}
	if idleconnectiontimeout, ok := d.GetOk(isLBListenerIdleConnectionTimeout); ok {
		idleConnectionTimeout := int64(idleconnectiontimeout.(int))
		options.IdleConnectionTimeout = &idleConnectionTimeout
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
	if _, ok := d.GetOk("https_redirect"); ok {
		httpsRedirectModel, err := resourceIBMIsLbListenerMapToLoadBalancerListenerHTTPSRedirectPrototype(d.Get("https_redirect.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		options.SetHTTPSRedirect(httpsRedirectModel)
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
		return diag.FromErr(fmt.Errorf("[ERROR] Error checking for load balancer (%s) is active: %s", lbID, err))
	}

	lbListener, response, err := sess.CreateLoadBalancerListener(options)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error while creating Load Balanacer Listener err %s\n%s", err, response))
	}
	d.SetId(fmt.Sprintf("%s/%s", lbID, *lbListener.ID))
	_, err = isWaitForLBListenerAvailable(sess, lbID, *lbListener.ID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error waiting for load balancer listener(%s) to become ready: %s", d.Id(), err))
	}
	_, err = isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error waiting for load balancer (%s) to become ready: %s", lbID, err))
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

func resourceIBMISLBListenerRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		diag.FromErr(err)
	}

	lbID := parts[0]
	lbListenerID := parts[1]

	diagEerr := lbListenerGet(d, meta, lbID, lbListenerID)
	if diagEerr != nil {
		return diagEerr
	}

	return nil
}

func lbListenerGet(d *schema.ResourceData, meta interface{}, lbID, lbListenerID string) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		diag.FromErr(err)
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
		return diag.FromErr(fmt.Errorf("[ERROR] Error Getting Load Balancer Listener : %s\n%s", err, response))
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
		if _, ok := d.GetOk("https_redirect"); ok {
			httpsRedirectMap, err := resourceIBMIsLbListenerLoadBalancerListenerHTTPSRedirectToMap(lbListener.HTTPSRedirect)
			if err != nil {
				return diag.FromErr(err)
			}
			if err = d.Set("https_redirect", []map[string]interface{}{httpsRedirectMap}); err != nil {
				return diag.FromErr(fmt.Errorf("Error setting https_redirect: %s", err))
			}

		} else {
			d.Set(isLBListenerHTTPSRedirectStatusCode, *lbListener.HTTPSRedirect.HTTPStatusCode)
			d.Set(isLBListenerHTTPSRedirectListener, *lbListener.HTTPSRedirect.Listener.ID)
			if lbListener.HTTPSRedirect.URI != nil {
				d.Set(isLBListenerHTTPSRedirectURI, *lbListener.HTTPSRedirect.URI)
			}
		}
	}
	if lbListener.CertificateInstance != nil {
		d.Set(isLBListenerCertificateInstance, *lbListener.CertificateInstance.CRN)
	}
	if lbListener.ConnectionLimit != nil {
		d.Set(isLBListenerConnectionLimit, *lbListener.ConnectionLimit)
	}
	d.Set(isLBListenerStatus, *lbListener.ProvisioningStatus)

	if lbListener.IdleConnectionTimeout != nil {
		d.Set(isLBListenerIdleConnectionTimeout, *lbListener.IdleConnectionTimeout)
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

func resourceIBMISLBListenerUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		diag.FromErr(err)
	}

	lbID := parts[0]
	lbListenerID := parts[1]

	diagEerr := lbListenerUpdate(d, meta, lbID, lbListenerID)
	if diagEerr != nil {
		return diagEerr
	}

	return resourceIBMISLBListenerRead(context, d, meta)
}

func lbListenerUpdate(d *schema.ResourceData, meta interface{}, lbID, lbListenerID string) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		diag.FromErr(err)
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
			diag.FromErr(err)
		}
		defPool = lbpool
		loadBalancerListenerPatchModel.DefaultPool = &vpcv1.LoadBalancerListenerDefaultPoolPatch{
			ID: &defPool,
		}
		hasChanged = true
	}
	httpsRedirectRemoved := false
	httpsURIRemoved := false
	if d.HasChange("https_redirect") {
		httpsRedirect := &vpcv1.LoadBalancerListenerHTTPSRedirectPatch{}
		if _, ok := d.GetOk("https_redirect"); !ok {
			httpsRedirectRemoved = true
		} else {
			if d.HasChange("https_redirect.0.http_status_code") {
				httpStatusCode := int64(d.Get("https_redirect.0.http_status_code").(int))
				httpsRedirect.HTTPStatusCode = &httpStatusCode
			}
			if d.HasChange("https_redirect.0.listener.0.id") {
				listenerId := d.Get("https_redirect.0.listener.0.id").(string)
				httpsRedirect.Listener = &vpcv1.LoadBalancerListenerIdentityByID{ID: &listenerId}
			}
			if d.HasChange("https_redirect.0.uri") {
				uri := d.Get("https_redirect.0.uri").(string)
				if uri == "" {
					httpsURIRemoved = true
				} else {
					httpsRedirect.URI = &uri
				}
			}
		}
		loadBalancerListenerPatchModel.HTTPSRedirect = httpsRedirect
		hasChanged = true
	} else {
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
	if d.HasChange(isLBListenerIdleConnectionTimeout) {
		idleConnectionTimeout := int64(d.Get(isLBListenerIdleConnectionTimeout).(int))
		loadBalancerListenerPatchModel.IdleConnectionTimeout = &idleConnectionTimeout
		hasChanged = true
	}
	if hasChanged {
		loadBalancerListenerPatch, err := loadBalancerListenerPatchModel.AsPatch()
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error calling asPatch for LoadBalancerListenerPatch: %s", err))
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
			return diag.FromErr(fmt.Errorf(
				"Error checking for load balancer (%s) is active: %s", lbID, err))
		}
		_, response, err := sess.UpdateLoadBalancerListener(updateLoadBalancerListenerOptions)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error Updating Load Balancer Listener : %s\n%s", err, response))
		}

		_, err = isWaitForLBListenerAvailable(sess, lbID, lbListenerID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(fmt.Errorf(
				"Error waiting for load balancer listener(%s) to become ready: %s", d.Id(), err))
		}

		_, err = isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(fmt.Errorf(
				"Error waiting for load balancer (%s) to become ready: %s", lbID, err))
		}
	}
	return nil
}

func resourceIBMISLBListenerDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		diag.FromErr(err)
	}

	lbID := parts[0]
	lbListenerID := parts[1]

	isLBKey := "load_balancer_key_" + lbID
	conns.IbmMutexKV.Lock(isLBKey)
	defer conns.IbmMutexKV.Unlock(isLBKey)

	diagEerr := lbListenerDelete(d, meta, lbID, lbListenerID)
	if diagEerr != nil {
		return diagEerr
	}

	return nil
}

func lbListenerDelete(d *schema.ResourceData, meta interface{}, lbID, lbListenerID string) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		diag.FromErr(err)
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
		return diag.FromErr(fmt.Errorf("[ERROR] Error Getting vpc load balancer listener(%s): %s\n%s", lbListenerID, err, response))
	}
	_, err = isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error checking for load balancer (%s) is active: %s", lbID, err))
	}
	deleteLoadBalancerListenerOptions := &vpcv1.DeleteLoadBalancerListenerOptions{
		LoadBalancerID: &lbID,
		ID:             &lbListenerID,
	}
	response, err = sess.DeleteLoadBalancerListener(deleteLoadBalancerListenerOptions)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error Deleting Load Balancer Pool : %s\n%s", err, response))
	}
	_, err = isWaitForLBListenerDeleted(sess, lbID, lbListenerID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		diag.FromErr(err)
	}
	_, err = isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error waiting for load balancer (%s) to be active: %s", lbID, err))
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

func resourceIBMIsLbListenerMapToLoadBalancerListenerHTTPSRedirectPrototype(modelMap map[string]interface{}) (*vpcv1.LoadBalancerListenerHTTPSRedirectPrototype, error) {
	model := &vpcv1.LoadBalancerListenerHTTPSRedirectPrototype{}
	model.HTTPStatusCode = core.Int64Ptr(int64(modelMap["http_status_code"].(int)))
	ListenerModel, err := resourceIBMIsLbListenerMapToLoadBalancerListenerIdentity(modelMap["listener"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Listener = ListenerModel
	if modelMap["uri"] != nil && modelMap["uri"].(string) != "" {
		model.URI = core.StringPtr(modelMap["uri"].(string))
	}
	return model, nil
}

func resourceIBMIsLbListenerMapToLoadBalancerListenerIdentity(modelMap map[string]interface{}) (vpcv1.LoadBalancerListenerIdentityIntf, error) {
	model := &vpcv1.LoadBalancerListenerIdentity{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	return model, nil
}

func resourceIBMIsLbListenerLoadBalancerListenerHTTPSRedirectToMap(model *vpcv1.LoadBalancerListenerHTTPSRedirect) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["http_status_code"] = flex.IntValue(model.HTTPStatusCode)
	listenerMap, err := resourceIBMIsLbListenerLoadBalancerListenerReferenceToMap(model.Listener)
	if err != nil {
		return modelMap, err
	}
	modelMap["listener"] = []map[string]interface{}{listenerMap}
	if model.URI != nil {
		modelMap["uri"] = model.URI
	}
	return modelMap, nil
}

func resourceIBMIsLbListenerLoadBalancerListenerReferenceToMap(model *vpcv1.LoadBalancerListenerReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Deleted != nil {
		deletedMap, err := resourceIBMIsLbListenerLoadBalancerListenerReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	return modelMap, nil
}
func resourceIBMIsLbListenerLoadBalancerListenerReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}
