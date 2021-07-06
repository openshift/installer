// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"log"

	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	ibmCISHealthCheck                = "ibm_cis_healthcheck"
	cisGLBHealthCheckID              = "monitor_id"
	cisGLBHealthCheckPath            = "path"
	cisGLBHealthCheckPort            = "port"
	cisGLBHealthCheckExpectedBody    = "expected_body"
	cisGLBHealthCheckExpectedCodes   = "expected_codes"
	cisGLBHealthCheckDesc            = "description"
	cisGLBHealthCheckType            = "type"
	cisGLBHealthCheckMethod          = "method"
	cisGLBHealthCheckTimeout         = "timeout"
	cisGLBHealthCheckRetries         = "retries"
	cisGLBHealthCheckInterval        = "interval"
	cisGLBHealthCheckFollowRedirects = "follow_redirects"
	cisGLBHealthCheckAllowInsecure   = "allow_insecure"
	cisGLBHealthCheckCreatedOn       = "create_on"
	cisGLBHealthCheckModifiedOn      = "modified_on"
	cisGLBHealthCheckHeaders         = "headers"
	cisGLBHealthCheckHeadersHeader   = "header"
	cisGLBHealthCheckHeadersValues   = "values"
)

func resourceIBMCISHealthCheck() *schema.Resource {
	return &schema.Resource{

		Create:   resourceCISHealthCheckCreate,
		Read:     resourceCISHealthCheckRead,
		Update:   resourceCISHealthCheckUpdate,
		Delete:   resourceCISHealthCheckDelete,
		Exists:   resourceCISHealthCheckExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
			},
			cisGLBHealthCheckID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "GLB Monitor/Health check id",
			},
			cisGLBHealthCheckPath: {
				Type:         schema.TypeString,
				Description:  "path",
				Optional:     true,
				Default:      "/",
				ValidateFunc: validateURLPath,
			},
			cisGLBHealthCheckExpectedBody: {
				Type:        schema.TypeString,
				Description: "expected_body",
				Optional:    true,
			},
			cisGLBHealthCheckExpectedCodes: {
				Type:        schema.TypeString,
				Description: "expected_codes",
				Optional:    true,
			},
			cisGLBHealthCheckDesc: {
				Type:        schema.TypeString,
				Description: "description",
				Default:     " ",
				Optional:    true,
			},
			cisGLBHealthCheckType: {
				Type:         schema.TypeString,
				Description:  "type",
				Optional:     true,
				Default:      "http",
				ValidateFunc: InvokeValidator(ibmCISHealthCheck, cisGLBHealthCheckType),
			},
			cisGLBHealthCheckMethod: {
				Type:         schema.TypeString,
				Description:  "method",
				Optional:     true,
				Default:      "GET",
				ValidateFunc: InvokeValidator(ibmCISHealthCheck, cisGLBHealthCheckMethod),
			},
			cisGLBHealthCheckTimeout: {
				Type:         schema.TypeInt,
				Description:  "timeout",
				Optional:     true,
				Default:      5,
				ValidateFunc: InvokeValidator(ibmCISHealthCheck, cisGLBHealthCheckTimeout),
			},
			cisGLBHealthCheckRetries: {
				Type:         schema.TypeInt,
				Description:  "retries",
				Optional:     true,
				Default:      2,
				ValidateFunc: InvokeValidator(ibmCISHealthCheck, cisGLBHealthCheckRetries),
			},
			cisGLBHealthCheckInterval: {
				Type:         schema.TypeInt,
				Description:  "interval",
				Optional:     true,
				Default:      60,
				ValidateFunc: InvokeValidator(ibmCISHealthCheck, cisGLBHealthCheckInterval),
			},
			cisGLBHealthCheckFollowRedirects: {
				Type:        schema.TypeBool,
				Description: "follow_redirects",
				Default:     false,
				Optional:    true,
			},
			cisGLBHealthCheckAllowInsecure: {
				Type:        schema.TypeBool,
				Description: "allow_insecure",
				Optional:    true,
				Default:     false,
			},
			cisGLBHealthCheckCreatedOn: {
				Type:     schema.TypeString,
				Computed: true,
			},
			cisGLBHealthCheckModifiedOn: {
				Type:     schema.TypeString,
				Computed: true,
			},
			cisGLBHealthCheckPort: {
				Type:         schema.TypeInt,
				Description:  "port number",
				Computed:     true,
				Optional:     true,
				ValidateFunc: InvokeValidator(ibmCISHealthCheck, cisGLBHealthCheckPort),
			},
			cisGLBHealthCheckHeaders: {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						cisGLBHealthCheckHeadersHeader: {
							Type:     schema.TypeString,
							Required: true,
						},

						cisGLBHealthCheckHeadersValues: {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
				Set: hashByMapKey(cisGLBHealthCheckHeadersHeader),
			},
		},
	}
}

func resourceIBMCISHealthCheckValidator() *ResourceValidator {
	healthCheckTypes := "http, https, tcp"
	methods := "GET, HEAD"

	validateSchema := make([]ValidateSchema, 1)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisGLBHealthCheckType,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              healthCheckTypes})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisGLBHealthCheckMethod,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              methods})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisGLBHealthCheckTimeout,
			ValidateFunctionIdentifier: IntBetween,
			Type:                       TypeInt,
			Required:                   true,
			MinValue:                   "1",
			MaxValue:                   "10"})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisGLBHealthCheckRetries,
			ValidateFunctionIdentifier: IntBetween,
			Type:                       TypeInt,
			Required:                   true,
			MinValue:                   "1",
			MaxValue:                   "3"})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisGLBHealthCheckInterval,
			ValidateFunctionIdentifier: IntBetween,
			Type:                       TypeInt,
			Required:                   true,
			MinValue:                   "5",
			MaxValue:                   "3600"})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisGLBHealthCheckPort,
			ValidateFunctionIdentifier: IntBetween,
			Type:                       TypeInt,
			Required:                   true,
			MinValue:                   "1",
			MaxValue:                   "65535"})
	cisHealthCheckValidator := ResourceValidator{ResourceName: ibmCISHealthCheck, Schema: validateSchema}
	return &cisHealthCheckValidator
}

func resourceCISHealthCheckCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).CisGLBHealthCheckClientSession()
	if err != nil {
		return err
	}

	crn := d.Get(cisID).(string)
	sess.Crn = core.StringPtr(crn)

	opt := sess.NewCreateLoadBalancerMonitorOptions()

	if monType, ok := d.GetOk(cisGLBHealthCheckType); ok {
		opt.SetType(monType.(string))
	}
	if expCodes, ok := d.GetOk(cisGLBHealthCheckExpectedCodes); ok {
		opt.SetExpectedCodes(expCodes.(string))
	}
	if expBody, ok := d.GetOk(cisGLBHealthCheckExpectedBody); ok {
		opt.SetExpectedBody(expBody.(string))
	}
	if monPath, ok := d.GetOk(cisGLBHealthCheckPath); ok {
		opt.SetPath(monPath.(string))
	}
	if description, ok := d.GetOk(cisGLBHealthCheckDesc); ok {
		opt.SetDescription(description.(string))
	}
	if method, ok := d.GetOk(cisGLBHealthCheckMethod); ok {
		opt.SetMethod(method.(string))
	}
	if timeout, ok := d.GetOk(cisGLBHealthCheckTimeout); ok {
		opt.SetTimeout(int64(timeout.(int)))
	}
	if retries, ok := d.GetOk(cisGLBHealthCheckRetries); ok {
		opt.SetRetries(int64(retries.(int)))
	}
	if interval, ok := d.GetOk(cisGLBHealthCheckInterval); ok {
		opt.SetInterval(int64(interval.(int)))
	}
	if followRedirects, ok := d.GetOk(cisGLBHealthCheckFollowRedirects); ok {
		opt.SetFollowRedirects(followRedirects.(bool))
	}
	if allowInsecure, ok := d.GetOk(cisGLBHealthCheckAllowInsecure); ok {
		opt.SetAllowInsecure(allowInsecure.(bool))
	}
	if port, ok := d.GetOk(cisGLBHealthCheckPort); ok {
		opt.SetPort(int64(port.(int)))
	}
	if header, ok := d.GetOk(cisGLBHealthCheckHeaders); ok {
		opt.SetHeader(expandLoadBalancerMonitorHeader(header))
	}

	result, resp, err := sess.CreateLoadBalancerMonitor(opt)
	if err != nil {
		log.Printf("create global load balancer health check failed %s", resp)
		return err
	}
	log.Printf("global load balancer created successfully : %s", *result.Result.ID)
	d.SetId(convertCisToTfTwoVar(*result.Result.ID, crn))
	return resourceCISHealthCheckRead(d, meta)
}

func resourceCISHealthCheckRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).CisGLBHealthCheckClientSession()
	if err != nil {
		return err
	}

	monitorID, crn, err := convertTftoCisTwoVar(d.Id())
	if err != nil {
		return err
	}
	sess.Crn = core.StringPtr(crn)

	opt := sess.NewGetLoadBalancerMonitorOptions(monitorID)

	result, resp, err := sess.GetLoadBalancerMonitor(opt)
	if err != nil {
		log.Printf("Error reading global load balancer health check detail: %s", resp)
		return err
	}
	d.Set(cisGLBHealthCheckID, result.Result.ID)
	d.Set(cisID, crn)
	d.Set(cisGLBHealthCheckDesc, result.Result.Description)
	d.Set(cisGLBHealthCheckPath, result.Result.Path)
	d.Set(cisGLBHealthCheckExpectedBody, result.Result.ExpectedBody)
	d.Set(cisGLBHealthCheckExpectedCodes, result.Result.ExpectedCodes)
	d.Set(cisGLBHealthCheckType, result.Result.Type)
	d.Set(cisGLBHealthCheckMethod, result.Result.Method)
	d.Set(cisGLBHealthCheckTimeout, result.Result.Timeout)
	d.Set(cisGLBHealthCheckRetries, result.Result.Retries)
	d.Set(cisGLBHealthCheckInterval, result.Result.Interval)
	d.Set(cisGLBHealthCheckFollowRedirects, result.Result.FollowRedirects)
	d.Set(cisGLBHealthCheckAllowInsecure, result.Result.AllowInsecure)
	d.Set(cisGLBHealthCheckPort, result.Result.Port)
	d.Set(cisGLBHealthCheckCreatedOn, result.Result.CreatedOn)
	d.Set(cisGLBHealthCheckModifiedOn, result.Result.ModifiedOn)
	if err := d.Set(cisGLBHealthCheckHeaders, flattenLoadBalancerMonitorHeader(result.Result.Header)); err != nil {
		log.Printf("[WARN] Error setting header for load balancer monitor %q: %s", d.Id(), err)
	}

	return nil
}

func resourceCISHealthCheckUpdate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).CisGLBHealthCheckClientSession()
	if err != nil {
		return err
	}

	monitorID, crn, err := convertTftoCisTwoVar(d.Id())
	if err != nil {
		return err
	}
	sess.Crn = core.StringPtr(crn)

	opt := sess.NewEditLoadBalancerMonitorOptions(monitorID)
	if d.HasChange(cisGLBHealthCheckType) ||
		d.HasChange(cisGLBHealthCheckDesc) ||
		d.HasChange(cisGLBHealthCheckPort) ||
		d.HasChange(cisGLBHealthCheckExpectedCodes) ||
		d.HasChange(cisGLBHealthCheckExpectedBody) ||
		d.HasChange(cisGLBHealthCheckMethod) ||
		d.HasChange(cisGLBHealthCheckTimeout) ||
		d.HasChange(cisGLBHealthCheckRetries) ||
		d.HasChange(cisGLBHealthCheckInterval) ||
		d.HasChange(cisGLBHealthCheckFollowRedirects) ||
		d.HasChange(cisGLBHealthCheckAllowInsecure) ||
		d.HasChange(cisGLBHealthCheckPort) ||
		d.HasChange(cisGLBHealthCheckHeaders) {
		if monType, ok := d.GetOk(cisGLBHealthCheckType); ok {
			opt.SetType(monType.(string))
		}
		if expCodes, ok := d.GetOk(cisGLBHealthCheckExpectedCodes); ok {
			opt.SetExpectedCodes(expCodes.(string))
		}
		if expBody, ok := d.GetOk(cisGLBHealthCheckExpectedBody); ok {
			opt.SetExpectedBody(expBody.(string))
		}
		if monPath, ok := d.GetOk(cisGLBHealthCheckPath); ok {
			opt.SetPath(monPath.(string))
		}
		if description, ok := d.GetOk(cisGLBHealthCheckDesc); ok {
			opt.SetDescription(description.(string))
		}
		if method, ok := d.GetOk(cisGLBHealthCheckMethod); ok {
			opt.SetMethod(method.(string))
		}
		if timeout, ok := d.GetOk(cisGLBHealthCheckTimeout); ok {
			opt.SetTimeout(int64(timeout.(int)))
		}
		if retries, ok := d.GetOk(cisGLBHealthCheckRetries); ok {
			opt.SetRetries(int64(retries.(int)))
		}
		if interval, ok := d.GetOk(cisGLBHealthCheckInterval); ok {
			opt.SetInterval(int64(interval.(int)))
		}
		if followRedirects, ok := d.GetOk(cisGLBHealthCheckFollowRedirects); ok {
			opt.SetFollowRedirects(followRedirects.(bool))
		}
		if allowInsecure, ok := d.GetOk(cisGLBHealthCheckAllowInsecure); ok {
			opt.SetAllowInsecure(allowInsecure.(bool))
		}
		if port, ok := d.GetOk(cisGLBHealthCheckPort); ok {
			opt.SetPort(int64(port.(int)))
		}
		if header, ok := d.GetOk(cisGLBHealthCheckHeaders); ok {
			opt.SetHeader(expandLoadBalancerMonitorHeader(header))
		}
		result, resp, err := sess.EditLoadBalancerMonitor(opt)
		if err != nil {
			log.Printf("Error updating global load balancer health check detail: %s", resp)
			return err
		}
		log.Printf("Monitor update succesful : %s", *result.Result.ID)
	}

	return resourceCISHealthCheckRead(d, meta)
}

func resourceCISHealthCheckDelete(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).CisGLBHealthCheckClientSession()
	if err != nil {
		return err
	}

	monitorID, crn, err := convertTftoCisTwoVar(d.Id())
	if err != nil {
		return err
	}
	sess.Crn = core.StringPtr(crn)

	opt := sess.NewDeleteLoadBalancerMonitorOptions(monitorID)

	result, resp, err := sess.DeleteLoadBalancerMonitor(opt)
	if err != nil {
		log.Printf("Error deleting global load balancer health check detail: %s", resp)
		return err
	}
	log.Printf("Monitor ID: %s", *result.Result.ID)
	return nil
}

func resourceCISHealthCheckExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess, err := meta.(ClientSession).CisGLBHealthCheckClientSession()
	if err != nil {
		return false, err
	}

	monitorID, crn, err := convertTftoCisTwoVar(d.Id())
	if err != nil {
		return false, err
	}
	sess.Crn = core.StringPtr(crn)

	opt := sess.NewGetLoadBalancerMonitorOptions(monitorID)

	result, response, err := sess.GetLoadBalancerMonitor(opt)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			log.Printf("global load balancer health check does not exist.")
			return false, nil
		}
		log.Printf("Error : %s", response)
		return false, err
	}
	log.Printf("global load balancer health check exists: %s", *result.Result.ID)
	return true, nil
}

func hashByMapKey(key string) func(v interface{}) int {
	return func(v interface{}) int {
		m := v.(map[string]interface{})
		return schema.HashString(m[key])
	}
}

func expandLoadBalancerMonitorHeader(cfgSet interface{}) map[string][]string {
	header := make(map[string][]string)
	cfgList := cfgSet.(*schema.Set).List()
	for _, item := range cfgList {
		cfg := item.(map[string]interface{})
		header[cfg[cisGLBHealthCheckHeadersHeader].(string)] =
			expandStringList(cfg[cisGLBHealthCheckHeadersValues].(*schema.Set).List())
	}
	return header
}

func flattenLoadBalancerMonitorHeader(header map[string][]string) *schema.Set {
	flattened := make([]interface{}, 0)
	for k, v := range header {
		cfg := map[string]interface{}{
			cisGLBHealthCheckHeadersHeader: k,
			cisGLBHealthCheckHeadersValues: schema.NewSet(schema.HashString, flattenStringList(v)),
		}
		flattened = append(flattened, cfg)
	}
	return schema.NewSet(hashByMapKey(cisGLBHealthCheckHeadersHeader), flattened)
}
