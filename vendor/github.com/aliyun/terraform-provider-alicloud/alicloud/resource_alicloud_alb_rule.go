package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudAlbRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAlbRuleCreate,
		Read:   resourceAlicloudAlbRuleRead,
		Update: resourceAlicloudAlbRuleUpdate,
		Delete: resourceAlicloudAlbRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"rule_actions": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"fixed_response_config": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"content": {
										Type:     schema.TypeString,
										Required: true,
									},
									"content_type": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"text/plain", "text/css", "text/html", "application/javascript", "application/json"}, false),
									},
									"http_code": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[2-5][0-9]{2}$`), "The http code must be an HTTP_2xx,HTTP_4xx or HTTP_5xx.x is a digit."),
									},
								},
							},
						},
						"forward_group_config": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"server_group_tuples": {
										Type:     schema.TypeSet,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"server_group_id": {
													Type:     schema.TypeString,
													Computed: true,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
						"insert_header_config": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[A-Za-z0-9_-]{1,40}$`), "The name of the header. The name must be 1 to 40 characters in length and can contain letters, digits, underscores (_), and hyphens (-)."),
									},
									"value": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"ClientSrcPort", "ClientSrcIp", "Protocol", "SLBId", "SLBPort", "UserDefined"}, false),
									},
									"value_type": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"UserDefined", "ReferenceHeader", "SystemDefined"}, false),
									},
								},
							},
						},
						"order": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 50000),
						},
						"redirect_config": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-z0-9\-\.\*\?]{3,128}$`), "The host name must be 3 to128 characters in length, and can contain lowercase letters, digits, hyphens (-), periods (.), asterisks (*), and question marks (?)."),
									},
									"http_code": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"301", "302", "303", "307", "308"}, false),
									},
									"path": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringMatch(regexp.MustCompile(`^\/[A-Za-z0-9\$\-_\.\+\/\&\~\@\:]{1,127}$`), "The value must be 1 to 128 characters in length and must start with a forward slash (/). The value can contain letters, digits, and the following special characters: $ - _ .+ / & ~ @ :. It cannot contain the following special characters: \" % # ; ! ( ) [ ]^ , \". The value is case-sensitive and can contain asterisks (*) and question marks (?)."),
									},
									"port": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: intBetween(1, 63335),
									},
									"protocol": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"HTTP", "HTTPS"}, false),
									},
									"query": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringLenBetween(1, 128),
									},
								},
							},
						},
						"rewrite_config": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-z0-9\-\.\*\?]{3,128}$`), "The host name must be 3 to128 characters in length, and can contain lowercase letters, digits, hyphens (-), periods (.), asterisks (*), and question marks (?)."),
									},
									"path": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringMatch(regexp.MustCompile(`^\/[A-Za-z0-9\$\-_\.\+\/\&\~\@\:]{1,127}$`), "The value must be 1 to 128 characters in length and must start with a forward slash (/). The value can contain letters, digits, and the following special characters: $ - _ .+ / & ~ @ :. It cannot contain the following special characters: \" % # ; ! ( ) [ ]^ , \". The value is case-sensitive and can contain asterisks (*) and question marks (?)."),
									},
									"query": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringLenBetween(1, 128),
									},
								},
							},
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"ForwardGroup", "Redirect", "FixedResponse", "Rewrite", "InsertHeader"}, false),
						},
					},
				},
			},
			"rule_conditions": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cookie_config": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"values": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringLenBetween(1, 100),
												},
												"value": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringLenBetween(1, 128),
												},
											},
										},
									},
								},
							},
						},
						"header_config": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[A-Za-z0-9_-]{1,40}$`), "The name of the header. The name must be 1 to 40 characters in length and can contain letters, digits, underscores (_), and hyphens (-)."),
									},
									"values": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"host_config": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"values": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"method_config": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"values": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"path_config": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"values": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"query_string_config": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"values": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringLenBetween(1, 100),
												},
												"value": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringLenBetween(1, 128),
												},
											},
										},
									},
								},
							},
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"Host", "Path", "Header", "HTTP", "QueryString", "Method", "Cookie"}, false),
						},
					},
				},
			},
			"rule_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudAlbRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateRule"
	request := make(map[string]interface{})
	conn, err := client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["ListenerId"] = d.Get("listener_id")
	request["Priority"] = d.Get("priority")
	ruleActionsMaps := make([]map[string]interface{}, 0)
	for _, ruleActions := range d.Get("rule_actions").(*schema.Set).List() {
		ruleActionsArg := ruleActions.(map[string]interface{})
		ruleActionsMap := map[string]interface{}{}
		ruleActionsMap["Order"] = ruleActionsArg["order"]
		ruleActionsMap["Type"] = ruleActionsArg["type"]

		if ruleActionsMap["Type"] == "" {
			continue
		}

		fixedResponseConfigMap := map[string]interface{}{}
		for _, fixedResponseConfig := range ruleActionsArg["fixed_response_config"].(*schema.Set).List() {
			fixedResponseConfigArg := fixedResponseConfig.(map[string]interface{})
			fixedResponseConfigMap["Content"] = fixedResponseConfigArg["content"]
			fixedResponseConfigMap["ContentType"] = fixedResponseConfigArg["content_type"]
			fixedResponseConfigMap["HttpCode"] = fixedResponseConfigArg["http_code"]
			ruleActionsMap["FixedResponseConfig"] = fixedResponseConfigMap
		}

		forwardGroupConfigMap := map[string]interface{}{}
		for _, forwardGroupConfig := range ruleActionsArg["forward_group_config"].(*schema.Set).List() {
			forwardGroupConfigArg := forwardGroupConfig.(map[string]interface{})
			serverGroupTuplesMaps := make([]map[string]interface{}, 0)
			for _, serverGroupTuples := range forwardGroupConfigArg["server_group_tuples"].(*schema.Set).List() {
				serverGroupTuplesArg := serverGroupTuples.(map[string]interface{})
				serverGroupTuplesMap := map[string]interface{}{}
				serverGroupTuplesMap["ServerGroupId"] = serverGroupTuplesArg["server_group_id"]
				serverGroupTuplesMaps = append(serverGroupTuplesMaps, serverGroupTuplesMap)
			}
			forwardGroupConfigMap["ServerGroupTuples"] = serverGroupTuplesMaps
			ruleActionsMap["ForwardGroupConfig"] = forwardGroupConfigMap
		}

		insertHeaderConfigMap := map[string]interface{}{}
		for _, insertHeaderConfig := range ruleActionsArg["insert_header_config"].(*schema.Set).List() {
			insertHeaderConfigArg := insertHeaderConfig.(map[string]interface{})
			insertHeaderConfigMap["Key"] = insertHeaderConfigArg["key"]
			insertHeaderConfigMap["Value"] = insertHeaderConfigArg["value"]
			insertHeaderConfigMap["ValueType"] = insertHeaderConfigArg["value_type"]
			ruleActionsMap["InsertHeaderConfig"] = insertHeaderConfigMap
		}

		redirectConfigMap := map[string]interface{}{}
		for _, redirectConfig := range ruleActionsArg["redirect_config"].(*schema.Set).List() {
			redirectConfigArg := redirectConfig.(map[string]interface{})
			redirectConfigMap["Host"] = redirectConfigArg["host"]
			redirectConfigMap["HttpCode"] = redirectConfigArg["http_code"]
			redirectConfigMap["Path"] = redirectConfigArg["path"]
			redirectConfigMap["Port"] = redirectConfigArg["port"]
			redirectConfigMap["Protocol"] = redirectConfigArg["protocol"]
			redirectConfigMap["Query"] = redirectConfigArg["query"]
			ruleActionsMap["RedirectConfig"] = redirectConfigMap
		}

		rewriteConfigMap := map[string]interface{}{}
		for _, rewriteConfig := range ruleActionsArg["rewrite_config"].(*schema.Set).List() {
			rewriteConfigArg := rewriteConfig.(map[string]interface{})
			rewriteConfigMap["Host"] = rewriteConfigArg["host"]
			rewriteConfigMap["Path"] = rewriteConfigArg["path"]
			rewriteConfigMap["Query"] = rewriteConfigArg["query"]
			ruleActionsMap["RewriteConfig"] = rewriteConfigMap
		}

		ruleActionsMaps = append(ruleActionsMaps, ruleActionsMap)
	}
	request["RuleActions"] = ruleActionsMaps

	ruleConditionsMaps := make([]map[string]interface{}, 0)
	for _, ruleConditions := range d.Get("rule_conditions").(*schema.Set).List() {
		ruleConditionsArg := ruleConditions.(map[string]interface{})
		ruleConditionsMap := map[string]interface{}{}
		ruleConditionsMap["Type"] = ruleConditionsArg["type"]

		cookieConfigMap := map[string]interface{}{}
		for _, cookieConfig := range ruleConditionsArg["cookie_config"].(*schema.Set).List() {
			cookieConfigArg := cookieConfig.(map[string]interface{})
			valuesMaps := make([]map[string]interface{}, 0)
			for _, values := range cookieConfigArg["values"].(*schema.Set).List() {
				valuesArg := values.(map[string]interface{})
				valuesMap := map[string]interface{}{}
				valuesMap["Key"] = valuesArg["key"]
				valuesMap["Value"] = valuesArg["value"]
				valuesMaps = append(valuesMaps, valuesMap)
			}
			cookieConfigMap["Values"] = valuesMaps
			ruleConditionsMap["CookieConfig"] = cookieConfigMap
		}

		headerConfigMap := map[string]interface{}{}
		for _, headerConfig := range ruleConditionsArg["header_config"].(*schema.Set).List() {
			headerConfigArg := headerConfig.(map[string]interface{})
			headerConfigMap["Key"] = headerConfigArg["key"]
			headerConfigMap["Values"] = headerConfigArg["values"].(*schema.Set).List()
			ruleConditionsMap["HeaderConfig"] = headerConfigMap
		}

		hostConfigMap := map[string]interface{}{}
		for _, hostConfig := range ruleConditionsArg["host_config"].(*schema.Set).List() {
			hostConfigArg := hostConfig.(map[string]interface{})
			hostConfigMap["Values"] = hostConfigArg["values"].(*schema.Set).List()
			ruleConditionsMap["HostConfig"] = hostConfigMap
		}

		methodConfigMap := map[string]interface{}{}
		for _, methodConfig := range ruleConditionsArg["method_config"].(*schema.Set).List() {
			methodConfigArg := methodConfig.(map[string]interface{})
			methodConfigMap["Values"] = methodConfigArg["values"].(*schema.Set).List()
			ruleConditionsMap["MethodConfig"] = methodConfigMap
		}

		pathConfigMap := map[string]interface{}{}
		for _, pathConfig := range ruleConditionsArg["path_config"].(*schema.Set).List() {
			pathConfigArg := pathConfig.(map[string]interface{})
			pathConfigMap["Values"] = pathConfigArg["values"].(*schema.Set).List()
			ruleConditionsMap["PathConfig"] = pathConfigMap
		}

		queryStringConfigMap := map[string]interface{}{}
		for _, queryStringConfig := range ruleConditionsArg["query_string_config"].(*schema.Set).List() {
			queryStringConfigArg := queryStringConfig.(map[string]interface{})
			valuesMaps := make([]map[string]interface{}, 0)
			for _, values := range queryStringConfigArg["values"].(*schema.Set).List() {
				valuesArg := values.(map[string]interface{})
				valuesMap := map[string]interface{}{}
				valuesMap["Key"] = valuesArg["key"]
				valuesMap["Value"] = valuesArg["value"]
				valuesMaps = append(valuesMaps, valuesMap)
			}
			queryStringConfigMap["Values"] = valuesMaps
			ruleConditionsMap["QueryStringConfig"] = queryStringConfigMap
		}

		ruleConditionsMaps = append(ruleConditionsMaps, ruleConditionsMap)
	}
	request["RuleConditions"] = ruleConditionsMaps

	request["RuleName"] = d.Get("rule_name")
	request["ClientToken"] = buildClientToken("CreateRule")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"IdempotenceProcessing", "IncorrectStatus.Listener", "SystemBusy", "Throttling"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alb_rule", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["RuleId"]))
	albService := AlbService{client}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, albService.AlbRuleStateRefreshFunc(d.Id(), []string{"CreateFailed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudAlbRuleRead(d, meta)
}
func resourceAlicloudAlbRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albService := AlbService{client}
	object, err := albService.DescribeAlbRule(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alb_rule albService.DescribeAlbRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("listener_id", object["ListenerId"])
	if v, ok := object["Priority"]; ok && fmt.Sprint(v) != "0" {
		d.Set("priority", formatInt(v))
	}

	if ruleActionsList, ok := object["RuleActions"]; ok {
		ruleActionsMaps := make([]map[string]interface{}, 0)
		for _, ruleActions := range ruleActionsList.([]interface{}) {
			ruleActionsArg := ruleActions.(map[string]interface{})
			ruleActionsMap := map[string]interface{}{}
			ruleActionsMap["type"] = ruleActionsArg["Type"]
			ruleActionsMap["order"] = formatInt(ruleActionsArg["Order"])

			if forwardGroupConfig, ok := ruleActionsArg["ForwardGroupConfig"]; ok {
				forwardGroupConfigArg := forwardGroupConfig.(map[string]interface{})
				if len(forwardGroupConfigArg) > 0 {
					serverGroupTuplesMaps := make([]map[string]interface{}, 0)
					if forwardGroupConfigArgs, ok := forwardGroupConfigArg["ServerGroupTuples"].([]interface{}); ok {
						for _, serverGroupTuples := range forwardGroupConfigArgs {
							serverGroupTuplesArg := serverGroupTuples.(map[string]interface{})
							serverGroupTuplesMap := map[string]interface{}{}
							serverGroupTuplesMap["server_group_id"] = serverGroupTuplesArg["ServerGroupId"]
							serverGroupTuplesMaps = append(serverGroupTuplesMaps, serverGroupTuplesMap)
						}
					}
					if len(serverGroupTuplesMaps) > 0 {
						forwardGroupConfigMaps := make([]map[string]interface{}, 0)
						forwardGroupConfigMap := map[string]interface{}{}
						forwardGroupConfigMap["server_group_tuples"] = serverGroupTuplesMaps
						forwardGroupConfigMaps = append(forwardGroupConfigMaps, forwardGroupConfigMap)
						ruleActionsMap["forward_group_config"] = forwardGroupConfigMaps
						ruleActionsMaps = append(ruleActionsMaps, ruleActionsMap)
					}
				}
			}

			if fixedResponseConfig, ok := ruleActionsArg["FixedResponseConfig"]; ok {
				fixedResponseConfigArg := fixedResponseConfig.(map[string]interface{})
				if len(fixedResponseConfigArg) > 0 {
					fixedResponseConfigMaps := make([]map[string]interface{}, 0)
					fixedResponseConfigMap := make(map[string]interface{}, 0)
					fixedResponseConfigMap["content"] = fixedResponseConfigArg["Content"]
					fixedResponseConfigMap["content_type"] = fixedResponseConfigArg["ContentType"]
					fixedResponseConfigMap["http_code"] = fixedResponseConfigArg["HttpCode"]
					fixedResponseConfigMaps = append(fixedResponseConfigMaps, fixedResponseConfigMap)
					ruleActionsMap["fixed_response_config"] = fixedResponseConfigMaps
					ruleActionsMaps = append(ruleActionsMaps, ruleActionsMap)
				}
			}

			if insertHeaderConfig, ok := ruleActionsArg["InsertHeaderConfig"]; ok {
				insertHeaderConfigArg := insertHeaderConfig.(map[string]interface{})
				if len(insertHeaderConfigArg) > 0 {
					insertHeaderConfigMaps := make([]map[string]interface{}, 0)
					insertHeaderConfigMap := make(map[string]interface{}, 0)
					insertHeaderConfigMap["key"] = insertHeaderConfigArg["Key"]
					insertHeaderConfigMap["value"] = insertHeaderConfigArg["Value"]
					insertHeaderConfigMap["value_type"] = insertHeaderConfigArg["ValueType"]
					insertHeaderConfigMaps = append(insertHeaderConfigMaps, insertHeaderConfigMap)
					ruleActionsMap["insert_header_config"] = insertHeaderConfigMaps
					ruleActionsMaps = append(ruleActionsMaps, ruleActionsMap)
				}
			}

			if redirectConfig, ok := ruleActionsArg["RedirectConfig"]; ok {
				redirectConfigArg := redirectConfig.(map[string]interface{})
				if len(redirectConfigArg) > 0 {
					redirectConfigMaps := make([]map[string]interface{}, 0)
					redirectConfigMap := make(map[string]interface{}, 0)
					redirectConfigMap["host"] = redirectConfigArg["Host"]
					redirectConfigMap["http_code"] = redirectConfigArg["HttpCode"]
					redirectConfigMap["path"] = redirectConfigArg["Path"]
					redirectConfigMap["port"] = formatInt(redirectConfigArg["Port"])
					redirectConfigMap["protocol"] = redirectConfigArg["Protocol"]
					redirectConfigMap["query"] = redirectConfigArg["Query"]
					redirectConfigMaps = append(redirectConfigMaps, redirectConfigMap)
					ruleActionsMap["redirect_config"] = redirectConfigMaps
					ruleActionsMaps = append(ruleActionsMaps, ruleActionsMap)
				}
			}

			if rewriteConfig, ok := ruleActionsArg["RewriteConfig"]; ok {
				rewriteConfigArg := rewriteConfig.(map[string]interface{})
				if len(rewriteConfigArg) > 0 {
					rewriteConfigMaps := make([]map[string]interface{}, 0)
					rewriteConfigMap := make(map[string]interface{}, 0)
					rewriteConfigMap["host"] = rewriteConfigArg["Host"]
					rewriteConfigMap["path"] = rewriteConfigArg["Path"]
					rewriteConfigMap["query"] = rewriteConfigArg["Query"]
					rewriteConfigMaps = append(rewriteConfigMaps, rewriteConfigMap)
					ruleActionsMap["rewrite_config"] = rewriteConfigMaps
					ruleActionsMaps = append(ruleActionsMaps, ruleActionsMap)
				}
			}
		}
		d.Set("rule_actions", ruleActionsMaps)
	}

	if ruleConditionsList, ok := object["RuleConditions"]; ok {
		ruleConditionsMaps := make([]map[string]interface{}, 0)
		for _, ruleConditions := range ruleConditionsList.([]interface{}) {
			ruleConditionsArg := ruleConditions.(map[string]interface{})
			ruleConditionsMap := map[string]interface{}{}
			ruleConditionsMap["type"] = ruleConditionsArg["Type"]

			if cookieConfig, ok := ruleConditionsArg["CookieConfig"]; ok {
				cookieConfigArg := cookieConfig.(map[string]interface{})
				if len(cookieConfigArg) > 0 {
					cookieConfigMaps := make([]map[string]interface{}, 0)
					valuesMaps := make([]map[string]interface{}, 0)
					for _, values := range cookieConfigArg["Values"].([]interface{}) {
						valuesArg := values.(map[string]interface{})
						valuesMap := map[string]interface{}{}
						valuesMap["key"] = valuesArg["Key"]
						valuesMap["value"] = valuesArg["Value"]
						valuesMaps = append(valuesMaps, valuesMap)
					}
					cookieConfigMap := map[string]interface{}{}
					cookieConfigMap["values"] = valuesMaps
					cookieConfigMaps = append(cookieConfigMaps, cookieConfigMap)
					ruleConditionsMap["cookie_config"] = cookieConfigMaps
					ruleConditionsMaps = append(ruleConditionsMaps, ruleConditionsMap)
				}
			}

			if headerConfig, ok := ruleConditionsArg["HeaderConfig"]; ok {
				headerConfigArg := headerConfig.(map[string]interface{})
				if len(headerConfigArg) > 0 {
					headerConfigMaps := make([]map[string]interface{}, 0)
					headerConfigMap := map[string]interface{}{}
					headerConfigMap["values"] = headerConfigArg["Values"].([]interface{})
					headerConfigMap["key"] = headerConfigArg["Key"]
					headerConfigMaps = append(headerConfigMaps, headerConfigMap)
					ruleConditionsMap["header_config"] = headerConfigMaps
					ruleConditionsMaps = append(ruleConditionsMaps, ruleConditionsMap)
				}
			}

			if queryStringConfig, ok := ruleConditionsArg["QueryStringConfig"]; ok {
				queryStringConfigArg := queryStringConfig.(map[string]interface{})
				if len(queryStringConfigArg) > 0 {
					queryStringConfigMaps := make([]map[string]interface{}, 0)
					queryStringValuesMaps := make([]map[string]interface{}, 0)
					for _, values := range queryStringConfigArg["Values"].([]interface{}) {
						valuesArg := values.(map[string]interface{})
						valuesMap := map[string]interface{}{}
						valuesMap["key"] = valuesArg["Key"]
						valuesMap["value"] = valuesArg["Value"]
						queryStringValuesMaps = append(queryStringValuesMaps, valuesMap)
					}
					queryStringConfigMap := map[string]interface{}{}
					queryStringConfigMap["values"] = queryStringValuesMaps
					queryStringConfigMaps = append(queryStringConfigMaps, queryStringConfigMap)
					ruleConditionsMap["query_string_config"] = queryStringConfigMaps
					ruleConditionsMaps = append(ruleConditionsMaps, ruleConditionsMap)
				}
			}

			if hostConfig, ok := ruleConditionsArg["HostConfig"]; ok {
				hostConfigArg := hostConfig.(map[string]interface{})
				if len(hostConfigArg) > 0 {
					hostConfigMaps := make([]map[string]interface{}, 0)
					hostConfigMap := map[string]interface{}{}
					hostConfigMap["values"] = hostConfigArg["Values"].([]interface{})
					hostConfigMaps = append(hostConfigMaps, hostConfigMap)
					ruleConditionsMap["host_config"] = hostConfigMaps
					ruleConditionsMaps = append(ruleConditionsMaps, ruleConditionsMap)
				}
			}

			if methodConfig, ok := ruleConditionsArg["MethodConfig"]; ok {
				methodConfigArg := methodConfig.(map[string]interface{})
				if len(methodConfigArg) > 0 {
					methodConfigMaps := make([]map[string]interface{}, 0)
					methodConfigMap := map[string]interface{}{}
					methodConfigMap["values"] = methodConfigArg["Values"].([]interface{})
					methodConfigMaps = append(methodConfigMaps, methodConfigMap)
					ruleConditionsMap["method_config"] = methodConfigMaps
					ruleConditionsMaps = append(ruleConditionsMaps, ruleConditionsMap)
				}
			}

			if pathConfig, ok := ruleConditionsArg["PathConfig"]; ok {
				pathConfigArg := pathConfig.(map[string]interface{})
				if len(pathConfigArg) > 0 {
					pathConfigMaps := make([]map[string]interface{}, 0)
					pathConfigMap := map[string]interface{}{}
					pathConfigMap["values"] = pathConfigArg["Values"].([]interface{})
					pathConfigMaps = append(pathConfigMaps, pathConfigMap)
					ruleConditionsMap["path_config"] = pathConfigMaps
					ruleConditionsMaps = append(ruleConditionsMaps, ruleConditionsMap)
				}
			}
		}

		d.Set("rule_conditions", ruleConditionsMaps)
	}

	d.Set("rule_name", object["RuleName"])
	d.Set("status", object["RuleStatus"])
	return nil
}
func resourceAlicloudAlbRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"RuleId": d.Id(),
	}
	if d.HasChange("priority") {
		update = true
		request["Priority"] = d.Get("priority")
	}
	if d.HasChange("rule_name") {
		update = true
		request["RuleName"] = d.Get("rule_name")
	}

	if d.HasChange("rule_actions") {
		update = true
		ruleActionsMaps := make([]map[string]interface{}, 0)
		for _, ruleActions := range d.Get("rule_actions").(*schema.Set).List() {
			ruleActionsArg := ruleActions.(map[string]interface{})
			ruleActionsMap := map[string]interface{}{}
			ruleActionsMap["Order"] = ruleActionsArg["order"]
			ruleActionsMap["Type"] = ruleActionsArg["type"]

			if ruleActionsMap["Type"] == "" {
				continue
			}

			fixedResponseConfigMap := map[string]interface{}{}
			for _, fixedResponseConfig := range ruleActionsArg["fixed_response_config"].(*schema.Set).List() {
				fixedResponseConfigArg := fixedResponseConfig.(map[string]interface{})
				fixedResponseConfigMap["Content"] = fixedResponseConfigArg["content"]
				fixedResponseConfigMap["ContentType"] = fixedResponseConfigArg["content_type"]
				fixedResponseConfigMap["HttpCode"] = fixedResponseConfigArg["http_code"]
				ruleActionsMap["FixedResponseConfig"] = fixedResponseConfigMap
			}

			forwardGroupConfigMap := map[string]interface{}{}
			for _, forwardGroupConfig := range ruleActionsArg["forward_group_config"].(*schema.Set).List() {
				forwardGroupConfigArg := forwardGroupConfig.(map[string]interface{})
				serverGroupTuplesMaps := make([]map[string]interface{}, 0)
				for _, serverGroupTuples := range forwardGroupConfigArg["server_group_tuples"].(*schema.Set).List() {
					serverGroupTuplesArg := serverGroupTuples.(map[string]interface{})
					serverGroupTuplesMap := map[string]interface{}{}
					serverGroupTuplesMap["ServerGroupId"] = serverGroupTuplesArg["server_group_id"]
					serverGroupTuplesMaps = append(serverGroupTuplesMaps, serverGroupTuplesMap)
				}
				forwardGroupConfigMap["ServerGroupTuples"] = serverGroupTuplesMaps
				ruleActionsMap["ForwardGroupConfig"] = forwardGroupConfigMap
			}

			insertHeaderConfigMap := map[string]interface{}{}
			for _, insertHeaderConfig := range ruleActionsArg["insert_header_config"].(*schema.Set).List() {
				insertHeaderConfigArg := insertHeaderConfig.(map[string]interface{})
				insertHeaderConfigMap["Key"] = insertHeaderConfigArg["key"]
				insertHeaderConfigMap["Value"] = insertHeaderConfigArg["value"]
				insertHeaderConfigMap["ValueType"] = insertHeaderConfigArg["value_type"]
				ruleActionsMap["InsertHeaderConfig"] = insertHeaderConfigMap
			}

			redirectConfigMap := map[string]interface{}{}
			for _, redirectConfig := range ruleActionsArg["redirect_config"].(*schema.Set).List() {
				redirectConfigArg := redirectConfig.(map[string]interface{})
				redirectConfigMap["Host"] = redirectConfigArg["host"]
				redirectConfigMap["HttpCode"] = redirectConfigArg["http_code"]
				redirectConfigMap["Path"] = redirectConfigArg["path"]
				redirectConfigMap["Port"] = redirectConfigArg["port"]
				redirectConfigMap["Protocol"] = redirectConfigArg["protocol"]
				redirectConfigMap["Query"] = redirectConfigArg["query"]
				ruleActionsMap["RedirectConfig"] = redirectConfigMap
			}

			rewriteConfigMap := map[string]interface{}{}
			for _, rewriteConfig := range ruleActionsArg["rewrite_config"].(*schema.Set).List() {
				rewriteConfigArg := rewriteConfig.(map[string]interface{})
				rewriteConfigMap["Host"] = rewriteConfigArg["host"]
				rewriteConfigMap["Path"] = rewriteConfigArg["path"]
				rewriteConfigMap["Query"] = rewriteConfigArg["query"]
				ruleActionsMap["RewriteConfig"] = rewriteConfigMap
			}

			ruleActionsMaps = append(ruleActionsMaps, ruleActionsMap)
		}
		request["RuleActions"] = ruleActionsMaps
	}

	if d.HasChange("rule_conditions") {
		update = true
		ruleConditionsMaps := make([]map[string]interface{}, 0)
		for _, ruleConditions := range d.Get("rule_conditions").(*schema.Set).List() {
			ruleConditionsArg := ruleConditions.(map[string]interface{})
			ruleConditionsMap := map[string]interface{}{}
			ruleConditionsMap["Type"] = ruleConditionsArg["type"]

			cookieConfigMap := map[string]interface{}{}
			for _, cookieConfig := range ruleConditionsArg["cookie_config"].(*schema.Set).List() {
				cookieConfigArg := cookieConfig.(map[string]interface{})
				valuesMaps := make([]map[string]interface{}, 0)
				for _, values := range cookieConfigArg["values"].(*schema.Set).List() {
					valuesArg := values.(map[string]interface{})
					valuesMap := map[string]interface{}{}
					valuesMap["Key"] = valuesArg["key"]
					valuesMap["Value"] = valuesArg["value"]
					valuesMaps = append(valuesMaps, valuesMap)
				}
				cookieConfigMap["Values"] = valuesMaps
				ruleConditionsMap["CookieConfig"] = cookieConfigMap
			}

			headerConfigMap := map[string]interface{}{}
			for _, headerConfig := range ruleConditionsArg["header_config"].(*schema.Set).List() {
				headerConfigArg := headerConfig.(map[string]interface{})
				headerConfigMap["Key"] = headerConfigArg["key"]
				headerConfigMap["Values"] = headerConfigArg["values"].(*schema.Set).List()
				ruleConditionsMap["HeaderConfig"] = headerConfigMap
			}

			hostConfigMap := map[string]interface{}{}
			for _, hostConfig := range ruleConditionsArg["host_config"].(*schema.Set).List() {
				hostConfigArg := hostConfig.(map[string]interface{})
				hostConfigMap["Values"] = hostConfigArg["values"].(*schema.Set).List()
				ruleConditionsMap["HostConfig"] = hostConfigMap
			}

			methodConfigMap := map[string]interface{}{}
			for _, methodConfig := range ruleConditionsArg["method_config"].(*schema.Set).List() {
				methodConfigArg := methodConfig.(map[string]interface{})
				methodConfigMap["Values"] = methodConfigArg["values"].(*schema.Set).List()
				ruleConditionsMap["MethodConfig"] = methodConfigMap
			}

			pathConfigMap := map[string]interface{}{}
			for _, pathConfig := range ruleConditionsArg["path_config"].(*schema.Set).List() {
				pathConfigArg := pathConfig.(map[string]interface{})
				pathConfigMap["Values"] = pathConfigArg["values"].(*schema.Set).List()
				ruleConditionsMap["PathConfig"] = pathConfigMap
			}

			queryStringConfigMap := map[string]interface{}{}
			for _, queryStringConfig := range ruleConditionsArg["query_string_config"].(*schema.Set).List() {
				queryStringConfigArg := queryStringConfig.(map[string]interface{})
				valuesMaps := make([]map[string]interface{}, 0)
				for _, values := range queryStringConfigArg["values"].(*schema.Set).List() {
					valuesArg := values.(map[string]interface{})
					valuesMap := map[string]interface{}{}
					valuesMap["Key"] = valuesArg["key"]
					valuesMap["Value"] = valuesArg["value"]
					valuesMaps = append(valuesMaps, valuesMap)
				}
				queryStringConfigMap["Values"] = valuesMaps
				ruleConditionsMap["QueryStringConfig"] = queryStringConfigMap
			}

			ruleConditionsMaps = append(ruleConditionsMaps, ruleConditionsMap)
		}
		request["RuleConditions"] = ruleConditionsMaps
	}

	if update {
		if v, ok := d.GetOkExists("dry_run"); ok {
			request["DryRun"] = v
		}
		action := "UpdateRuleAttribute"
		conn, err := client.NewAlbClient()
		if err != nil {
			return WrapError(err)
		}
		request["ClientToken"] = buildClientToken("UpdateRuleAttribute")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		albService := AlbService{client}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, albService.AlbRuleStateRefreshFunc(d.Id(), []string{"CreateFailed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	return resourceAlicloudAlbRuleRead(d, meta)
}
func resourceAlicloudAlbRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteRule"
	var response map[string]interface{}
	conn, err := client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"RuleId": d.Id(),
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["ClientToken"] = buildClientToken("DeleteRule")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"IdempotenceProcessing", "IncorrectStatus.Rule", "SystemBusy", "Throttling"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"ResourceNotFound.Rule"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
