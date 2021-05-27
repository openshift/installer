// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"

	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/IBM/networking-go-sdk/zoneratelimitsv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	cisRLThreshold   = "threshold"
	cisRLPeriod      = "period"
	cisRLDescription = "description"
	cisRLTimeout     = "timeout"
	cisRLBody        = "body"
	cisRLURL         = "url"
)

func resourceIBMCISRateLimit() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMCISRateLimitCreate,
		Read:     resourceIBMCISRateLimitRead,
		Update:   resourceIBMCISRateLimitUpdate,
		Delete:   resourceIBMCISRateLimitDelete,
		Exists:   resourceIBMCISRateLimitExists,
		Importer: &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			"cis_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CIS Intance CRN",
			},
			"domain_id": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "CIS Domain ID",
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			"disabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether this rate limiting rule is currently disabled.",
			},
			cisRLDescription: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: InvokeValidator("ibm_cis_rate_limit", cisRLDescription),
				Description:  "A note that you can use to describe the reason for a rate limiting rule.",
			},
			"bypass": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Bypass URL",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "url",
							Description: "bypass URL name",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "bypass URL value",
						},
					},
				},
			},
			cisRLThreshold: {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: InvokeValidator("ibm_cis_rate_limit", cisRLThreshold),
				Description:  "Rate Limiting Threshold",
			},
			cisRLPeriod: {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: InvokeValidator("ibm_cis_rate_limit", cisRLPeriod),
				Description:  "Rate Limiting Period",
			},
			"correlate": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Ratelimiting Correlate",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"by": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "nat",
							ValidateFunc: InvokeValidator("ibm_cis_rate_limit", "by"),
							Description:  "Whether to enable NAT based rate limiting",
						},
					},
				},
			},
			"action": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				MaxItems:    1,
				Description: "Rate Limiting Action",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: InvokeValidator("ibm_cis_rate_limit", "mode"),
							Description:  "Type of action performed.Valid values are: 'simulate', 'ban', 'challenge', 'js_challenge'.",
						},
						"timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: InvokeValidator("ibm_cis_rate_limit", cisRLTimeout),
							Description:  "The time to perform the mitigation action. Timeout be the same or greater than the period.",
						},
						"response": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Rate Limiting Action Response",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"content_type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: InvokeValidator("ibm_cis_rate_limit", "content_type"),
										Description:  "Custom content-type and body to return. It must be one of following 'text/plain', 'text/xml', 'application/json'.",
									},
									"body": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: InvokeValidator("ibm_cis_rate_limit", cisRLBody),
										Description:  "The body to return. The content here must confirm to the 'content_type'",
									},
								},
							},
						},
					},
				},
			},
			"match": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Computed:    true,
				Description: "Rate Limiting Match",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"request": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MinItems:    1,
							MaxItems:    1,
							Description: "Rate Limiting Match Request",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"methods": {
										Type:        schema.TypeSet,
										Optional:    true,
										Computed:    true,
										Description: "HTTP Methos of matching request. It can be one or many. Example methods 'POST', 'PUT'",
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: InvokeValidator("ibm_cis_rate_limit", "methods"),
										},
									},
									"schemes": {
										Type:        schema.TypeSet,
										Optional:    true,
										Computed:    true,
										Description: "HTTP Schemes of matching request. It can be one or many. Example schemes 'HTTP', 'HTTPS'.",
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: InvokeValidator("ibm_cis_rate_limit", "schemes"),
										},
									},
									"url": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										Description:  "URL pattern of matching request",
										ValidateFunc: InvokeValidator("ibm_cis_rate_limit", cisRLURL),
									},
								},
							},
						},
						"response": {
							Type:        schema.TypeList,
							Optional:    true,
							MinItems:    1,
							MaxItems:    1,
							Description: "Rate Limiting Response",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status": {
										Type:        schema.TypeSet,
										Optional:    true,
										Description: "HTTP Status Codes of matching response. It can be one or many. Example status codes '403', '401",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
									"origin_traffic": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Origin Traffic of matching response.",
									},
									"headers": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The name of the response header to match.",
												},
												"op": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The operator when matching. Valid values are 'eq' and 'ne'.",
												},
												"value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The value of the header, which is exactly matched.",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rate Limit rule Id",
			},
		},
	}
}
func resourceIBMCISRateLimitValidator() *ResourceValidator {

	validateSchema := make([]ValidateSchema, 1)
	byValues := "nat"
	modeValues := "simulate, ban, challenge, js_challenge"
	ctypeValues := "text/plain, text/xml, application/json"
	methodValues := "GET, POST, PUT, DELETE, PATCH, HEAD, _ALL_"
	schemeValues := "HTTP, HTTPS, _ALL_"

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisRLDescription,
			ValidateFunctionIdentifier: StringLenBetween,
			Type:                       TypeString,
			Optional:                   true,
			MinValueLength:             0,
			MaxValueLength:             1024})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisRLThreshold,
			ValidateFunctionIdentifier: IntBetween,
			Type:                       TypeInt,
			Required:                   true,
			MinValue:                   "1",
			MaxValue:                   "1000000"})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisRLPeriod,
			ValidateFunctionIdentifier: IntBetween,
			Type:                       TypeInt,
			Required:                   true,
			MinValue:                   "1",
			MaxValue:                   "86400"})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "by",
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              byValues})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "mode",
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              modeValues})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "content_type",
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              ctypeValues})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "methods",
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              methodValues})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "schemes",
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              schemeValues})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisRLTimeout,
			ValidateFunctionIdentifier: IntBetween,
			Type:                       TypeInt,
			Optional:                   true,
			MinValue:                   "1",
			MaxValue:                   "86400"})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisRLBody,
			ValidateFunctionIdentifier: StringLenBetween,
			Type:                       TypeString,
			Optional:                   true,
			MinValueLength:             0,
			MaxValueLength:             10240})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisRLURL,
			ValidateFunctionIdentifier: StringLenBetween,
			Type:                       TypeString,
			Optional:                   true,
			MinValueLength:             0,
			MaxValueLength:             1024})

	ibmCISRateLimitResourceValidator := ResourceValidator{ResourceName: "ibm_cis_rate_limit", Schema: validateSchema}
	return &ibmCISRateLimitResourceValidator
}
func resourceIBMCISRateLimitCreate(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisRLClientSession()
	if err != nil {
		return err
	}

	cisID := d.Get("cis_id").(string)
	zoneID, _, err := convertTftoCisTwoVar(d.Get("domain_id").(string))
	if err != nil {
		return err
	}
	cisClient.Crn = core.StringPtr(cisID)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	//payload to create a rate limit rule
	opt := cisClient.NewCreateZoneRateLimitsOptions()
	opt.SetThreshold(int64(d.Get(cisRLThreshold).(int)))
	opt.SetPeriod(int64(d.Get(cisRLPeriod).(int)))

	if description, ok := d.GetOk(cisRLDescription); ok {
		opt.SetDescription(description.(string))
	}

	if disabled, ok := d.GetOk("disabled"); ok {
		opt.SetDisabled(disabled.(bool))
	}

	action, err := expandRateLimitAction(d)
	if err != nil {
		return fmt.Errorf("Error in getting action from expandRateLimitAction %s", err)
	}
	opt.SetAction(action)

	match, err := expandRateLimitMatch(d)
	if err != nil {
		return fmt.Errorf("Error in getting match from expandRateLimitMatch %s", err)
	}
	opt.SetMatch(match)

	correlate, err := expandRateLimitCorrelate(d)
	if err == nil {
		opt.SetCorrelate(correlate)
	}

	byPass, err := expandRateLimitBypass(d)
	if err != nil {
		return fmt.Errorf("Error in getting bypass from expandRateLimitBypass %s", err)
	}
	opt.SetBypass(byPass)

	//creating rate limit rule
	result, resp, err := cisClient.CreateZoneRateLimits(opt)
	if err != nil {
		return fmt.Errorf("Failed to create RateLimit: %v", resp)
	}
	record := result.Result
	d.SetId(convertCisToTfThreeVar(*record.ID, zoneID, cisID))
	return resourceIBMCISRateLimitRead(d, meta)
}

func resourceIBMCISRateLimitRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisRLClientSession()
	if err != nil {
		return err
	}
	recordID, zoneID, cisID, _ := convertTfToCisThreeVar(d.Id())
	if err != nil {
		return err
	}
	cisClient.Crn = core.StringPtr(cisID)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)
	opt := cisClient.NewGetRateLimitOptions(recordID)
	result, resp, err := cisClient.GetRateLimit(opt)
	if err != nil {
		return fmt.Errorf("Failed to read RateLimit: %v", resp)
	}

	rule := result.Result
	d.Set("cis_id", cisID)
	d.Set("domain_id", zoneID)
	d.Set("rule_id", recordID)
	d.Set("disabled", rule.Disabled)
	d.Set(cisRLDescription, rule.Description)
	d.Set(cisRLThreshold, rule.Threshold)
	d.Set(cisRLPeriod, rule.Period)
	d.Set("action", flattenRateLimitAction(rule.Action))
	d.Set("match", flattenRateLimitMatch(rule.Match))
	d.Set("correlate", flattenRateLimitCorrelate(rule.Correlate))
	d.Set("bypass", flattenRateLimitByPass(rule.Bypass))

	return nil
}

func resourceIBMCISRateLimitUpdate(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisRLClientSession()
	if err != nil {
		return err
	}

	recordID, zoneID, cisID, _ := convertTfToCisThreeVar(d.Id())
	cisClient.Crn = core.StringPtr(cisID)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	opt := cisClient.NewUpdateRateLimitOptions(recordID)
	if d.HasChange("disabled") ||
		d.HasChange(cisRLThreshold) ||
		d.HasChange(cisRLPeriod) ||
		d.HasChange(cisRLDescription) ||
		d.HasChange("action") ||
		d.HasChange("match") ||
		d.HasChange("correlate") ||
		d.HasChange("bypass") {

		opt.SetThreshold(int64(d.Get(cisRLThreshold).(int)))
		opt.SetPeriod(int64(d.Get(cisRLPeriod).(int)))

		if description, ok := d.GetOk(cisRLDescription); ok {
			opt.SetDescription(description.(string))
		}

		if disabled, ok := d.GetOk("disabled"); ok {
			opt.SetDisabled(disabled.(bool))
		}

		action, err := expandRateLimitAction(d)
		if err != nil {
			return fmt.Errorf("Error in getting action from expandRateLimitAction %s", err)
		}
		opt.SetAction(action)

		match, err := expandRateLimitMatch(d)
		if err != nil {
			return fmt.Errorf("Error in getting match from expandRateLimitMatch %s", err)
		}
		opt.SetMatch(match)

		correlate, err := expandRateLimitCorrelate(d)
		if err != nil {
			return fmt.Errorf("Error in getting correlate from expandRateLimitCorrelate %s", err)
		}
		opt.SetCorrelate(correlate)

		byPass, err := expandRateLimitBypass(d)
		if err != nil {
			return fmt.Errorf("Error in getting bypass from expandRateLimitBypass %s", err)
		}
		opt.SetBypass(byPass)
		_, resp, err := cisClient.UpdateRateLimit(opt)
		if err != nil {
			return fmt.Errorf("Failed to update RateLimit: %v", resp)
		}
	}
	d.SetId(convertCisToTfThreeVar(recordID, zoneID, cisID))
	return resourceIBMCISRateLimitRead(d, meta)
}

func resourceIBMCISRateLimitDelete(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisRLClientSession()
	if err != nil {
		return err
	}

	recordID, zoneID, cisID, _ := convertTfToCisThreeVar(d.Id())
	cisClient.Crn = core.StringPtr(cisID)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)
	opt := cisClient.NewDeleteZoneRateLimitOptions(recordID)
	_, resp, err := cisClient.DeleteZoneRateLimit(opt)
	if err != nil {
		return fmt.Errorf("Failed to delete RateLimit: %v", resp)
	}
	return nil
}

func resourceIBMCISRateLimitExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	cisClient, err := meta.(ClientSession).CisRLClientSession()
	if err != nil {
		return false, err
	}
	recordID, zoneID, cisID, _ := convertTfToCisThreeVar(d.Id())
	cisClient.Crn = core.StringPtr(cisID)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)
	opt := cisClient.NewGetRateLimitOptions(recordID)
	_, resp, err := cisClient.GetRateLimit(opt)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Println("ratelimit is not found")
			return false, nil
		}
		return false, fmt.Errorf("Failed to getting existing RateLimit: %v", err)
	}
	return true, nil
}

func expandRateLimitAction(d *schema.ResourceData) (
	action *zoneratelimitsv1.RatelimitInputAction, err error) {
	action = &zoneratelimitsv1.RatelimitInputAction{}
	actionRecord := d.Get("action").([]interface{})[0].(map[string]interface{})
	mode := actionRecord["mode"].(string)
	timeout := actionRecord["timeout"].(int)
	if timeout == 0 {
		if mode == "simulate" || mode == "ban" {
			return action, fmt.Errorf("For the mode 'simulate' and 'ban' timeout must be %s %s",
				"set.. valid range for timeout is 10 - 86400", err)
		}
	} else {
		if mode == "challenge" || mode == "js_challenge" {
			return action, fmt.Errorf(
				"Timeout field is only valid for 'simulate' and 'ban' modes. %s", err)
		}
	}
	action.Mode = core.StringPtr(mode)
	action.Timeout = core.Int64Ptr(int64(timeout))

	if _, ok := actionRecord["response"]; ok && len(actionRecord["response"].([]interface{})) > 0 {
		actionResponse := actionRecord["response"].([]interface{})[0].(map[string]interface{})
		action.Response = &zoneratelimitsv1.RatelimitInputActionResponse{
			ContentType: core.StringPtr(actionResponse["content_type"].(string)),
			Body:        core.StringPtr(actionResponse["body"].(string)),
		}
	}

	return action, nil
}

func expandRateLimitMatch(d *schema.ResourceData) (match *zoneratelimitsv1.RatelimitInputMatch, err error) {
	match = new(zoneratelimitsv1.RatelimitInputMatch)
	m := d.Get("match")
	if len(m.([]interface{})) == 0 {
		// Match Request is a mondatory property. So, setting default if none provided.
		match.Request = &zoneratelimitsv1.RatelimitInputMatchRequest{
			Methods: []string{"_ALL_"},
			Schemes: []string{"_ALL_"},
			URL:     core.StringPtr("*"),
		}
		return match, nil
	}
	matchRecord := m.([]interface{})[0].(map[string]interface{})

	if matchReqRecord, ok := matchRecord["request"]; ok && len(matchReqRecord.([]interface{})) > 0 {
		matchRequestRecord := matchReqRecord.([]interface{})[0].(map[string]interface{})

		url := matchRequestRecord["url"].(string)
		// If url is not provided, then set it with *
		if len(url) == 0 {
			url = "*"
		}
		matchRequest := &zoneratelimitsv1.RatelimitInputMatchRequest{
			URL: core.StringPtr(url),
		}
		if methodsRecord, ok := matchRequestRecord["methods"]; ok {
			methods := make([]string, methodsRecord.(*schema.Set).Len())
			for i, m := range methodsRecord.(*schema.Set).List() {
				methods[i] = m.(string)
			}
			matchRequest.Methods = methods
		}
		if schemesRecord, ok := matchRequestRecord["schemes"]; ok {
			schemes := make([]string, schemesRecord.(*schema.Set).Len())
			for i, s := range schemesRecord.(*schema.Set).List() {
				schemes[i] = s.(string)
			}
			matchRequest.Schemes = schemes
		}

		match.Request = matchRequest
	}
	if matchResRecord, ok := matchRecord["response"]; ok && len(matchResRecord.([]interface{})) > 0 {
		matchResponseRecord := matchResRecord.([]interface{})[0].(map[string]interface{})
		matchResponse := &zoneratelimitsv1.RatelimitInputMatchResponse{}
		if statusRecord, ok := matchResponseRecord["status"]; ok {
			status := make([]int64, statusRecord.(*schema.Set).Len())
			for i, s := range statusRecord.(*schema.Set).List() {
				status[i] = int64(s.(int))
			}
			matchResponse.Status = status
		}
		if originRecord, ok := matchResponseRecord["origin_traffic"]; ok {
			originTraffic := originRecord.(bool)
			matchResponse.OriginTraffic = &originTraffic
		}
		if headersRecord, ok := matchResponseRecord["headers"]; ok && len(headersRecord.([]interface{})) > 0 {
			matchResponseHeaders := headersRecord.([]interface{})

			responseHeaders := make([]zoneratelimitsv1.RatelimitInputMatchResponseHeadersItem, 0)

			for _, h := range matchResponseHeaders {
				header := h.(map[string]interface{})
				headerRecord := zoneratelimitsv1.RatelimitInputMatchResponseHeadersItem{}
				headerRecord.Name = core.StringPtr(header["name"].(string))
				headerRecord.Op = core.StringPtr(header["op"].(string))
				headerRecord.Value = core.StringPtr(header["value"].(string))
				responseHeaders = append(responseHeaders, headerRecord)
			}
			matchResponse.HeadersVar = responseHeaders

		}
		match.Response = matchResponse
	}

	return match, nil
}

func expandRateLimitCorrelate(d *schema.ResourceData) (
	correlate *zoneratelimitsv1.RatelimitInputCorrelate, err error) {
	correlate = &zoneratelimitsv1.RatelimitInputCorrelate{}
	c, ok := d.GetOk("correlate")
	if !ok {
		err = fmt.Errorf("correlate field is empty")
		return &zoneratelimitsv1.RatelimitInputCorrelate{}, err
	}
	correlateRecord := c.([]interface{})[0].(map[string]interface{})
	correlate.By = core.StringPtr(correlateRecord["by"].(string))

	return correlate, nil
}

func expandRateLimitBypass(d *schema.ResourceData) (
	byPass []zoneratelimitsv1.RatelimitInputBypassItem, err error) {
	b, ok := d.GetOk("bypass")
	if !ok {
		return
	}
	byPassKV := b.([]interface{})

	byPassRecord := make([]zoneratelimitsv1.RatelimitInputBypassItem, 0)

	for _, kv := range byPassKV {
		keyValue, _ := kv.(map[string]interface{})

		byPassKeyValue := zoneratelimitsv1.RatelimitInputBypassItem{}
		byPassKeyValue.Name = core.StringPtr(keyValue["name"].(string))
		byPassKeyValue.Value = core.StringPtr(keyValue["value"].(string))
		byPassRecord = append(byPassRecord, byPassKeyValue)
	}
	byPass = byPassRecord

	return byPass, nil
}

func flattenRateLimitAction(action *zoneratelimitsv1.RatelimitObjectAction) []map[string]interface{} {
	actionRecord := map[string]interface{}{
		"mode":    *action.Mode,
		"timeout": *action.Timeout,
	}

	if action.Response != nil {
		actionResponseRecord := *action.Response
		actionResponse := map[string]interface{}{
			"content_type": *actionResponseRecord.ContentType,
			"body":         *actionResponseRecord.Body,
		}
		actionRecord["response"] = []map[string]interface{}{actionResponse}
	}
	return []map[string]interface{}{actionRecord}
}

func flattenRateLimitMatch(match *zoneratelimitsv1.RatelimitObjectMatch) []map[string]interface{} {
	matchRecord := map[string]interface{}{}
	matchRecord["request"] = flattenRateLimitMatchRequest(*match.Request)
	if match.Response != nil {
		matchRecord["response"] = flattenRateLimitMatchResponse(*match.Response)
	}

	return []map[string]interface{}{matchRecord}
}

func flattenRateLimitMatchRequest(request zoneratelimitsv1.RatelimitObjectMatchRequest) []map[string]interface{} {

	requestRecord := map[string]interface{}{}
	methods := make([]string, 0)
	for _, m := range request.Methods {
		methods = append(methods, m)
	}
	requestRecord["methods"] = methods
	schemes := make([]string, 0)
	for _, s := range request.Schemes {
		schemes = append(schemes, s)
	}
	requestRecord["schemes"] = schemes

	requestRecord["url"] = *request.URL
	return []map[string]interface{}{requestRecord}
}

func flattenRateLimitMatchResponse(response zoneratelimitsv1.RatelimitObjectMatchResponse) []interface{} {
	responseRecord := map[string]interface{}{}
	flag := false
	if response.OriginTraffic != nil {
		responseRecord["origin_traffic"] = *response.OriginTraffic
		flag = true
	}

	if len(response.Status) > 0 {
		statuses := make([]int64, 0)
		for _, s := range response.Status {
			statuses = append(statuses, s)
		}
		responseRecord["status"] = statuses
		flag = true
	}

	if len(response.HeadersVar) > 0 {
		headers := make([]map[string]interface{}, 0)
		for _, h := range response.HeadersVar {
			header := map[string]interface{}{}
			header["name"] = h.Name
			header["op"] = h.Op
			header["value"] = h.Value
			headers = append(headers, header)

		}
		responseRecord["headers"] = headers
		flag = true
	}
	if flag == true {
		return []interface{}{responseRecord}
	}
	return []interface{}{}
}
func flattenRateLimitCorrelate(correlate *zoneratelimitsv1.RatelimitObjectCorrelate) []map[string]interface{} {
	if correlate == nil {
		return []map[string]interface{}{}
	}
	correlateRecord := map[string]interface{}{}
	if *correlate.By != "" {
		correlateRecord["by"] = *correlate.By
	}
	return []map[string]interface{}{correlateRecord}
}

func flattenRateLimitByPass(byPass []zoneratelimitsv1.RatelimitObjectBypassItem) []map[string]interface{} {
	byPassRecord := make([]map[string]interface{}, 0, len(byPass))
	if len(byPass) > 0 {
		for _, b := range byPass {
			byPassKV := map[string]interface{}{
				"name":  *b.Name,
				"value": *b.Value,
			}
			byPassRecord = append(byPassRecord, byPassKV)
		}
	}
	return byPassRecord
}
