package cos

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	validation "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceIBMCOSBucketWebsiteConfiguration() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMCOSBucketWebsiteConfigurationCreate,
		Read:     resourceIBMCOSBucketWebsiteConfigurationRead,
		Update:   resourceIBMCOSBucketWebsiteConfigurationUpdate,
		Delete:   resourceIBMCOSBucketWebsiteConfigurationDelete,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bucket_crn": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "COS bucket CRN",
			},
			"bucket_location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "COS bucket location",
			},
			"endpoint_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{"public", "private", "direct"}),
				Description:  "COS endpoint type: public, private, direct",
				Default:      "public",
			},
			"website_configuration": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Configuration for Hosting a static website on COS with public access.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"error_document": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "This is returned when an error occurs.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"index_document": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Home or the default page of the website.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"suffix": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"redirect_all_requests_to": {
							Type:          schema.TypeList,
							Optional:      true,
							MaxItems:      1,
							Description:   "Redirect requests can be set to specific page documents, individual routing rules, or redirect all requests globally to one bucket or domain.",
							ConflictsWith: []string{"website_configuration.0.error_document", "website_configuration.0.index_document", "website_configuration.0.routing_rule", "website_configuration.0.routing_rules"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"protocol": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice(s3.Protocol_Values(), false),
									},
								},
							},
						},
						"routing_rule": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Rules that define when a redirect is applied and the redirect behavior.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"condition": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "A condition that must be met for the specified redirect to be applie.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"http_error_code_returned_equals": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The HTTP error code when the redirect is applied. Valid codes are 4XX or 5XX..",
												},
												"key_prefix_equals": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The object key name prefix when the redirect is applied..",
												},
											},
										},
									},
									"redirect": {
										Type:        schema.TypeList,
										Required:    true,
										MaxItems:    1,
										Description: ".",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"host_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The host name the request should be redirected to.",
												},
												"http_redirect_code": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The HTTP redirect code to use on the response. Valid codes are 3XX except 300..",
												},
												"protocol": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringInSlice(s3.Protocol_Values(), false),
													Description:  "Protocol to be used in the Location header that is returned in the response.",
												},
												"replace_key_prefix_with": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The prefix of the object key name that replaces the value of KeyPrefixEquals in the redirect request.",
												},
												"replace_key_with": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The object key to be used in the Location header that is returned in the response.",
												},
											},
										},
									},
								},
							},
						},
						"routing_rules": {
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							Description:   "Rules that define when a redirect is applied and the redirect behavior.",
							ConflictsWith: []string{"website_configuration.0.routing_rule"},
							ValidateFunc:  validation.StringIsJSON,
							StateFunc: func(v interface{}) string {
								json, _ := structure.NormalizeJsonString(v)
								return json
							},
						},
					},
				},
			},
			"website_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func errorDocumentSetFunction(errorDocument []interface{}) *s3.ErrorDocument {
	var error_document *s3.ErrorDocument
	errorDocumentValue := s3.ErrorDocument{}
	if len(errorDocument) != 0 {
		errorDocumentMap := errorDocument[0].(map[string]interface{})
		if keyValue, ok := errorDocumentMap["key"].(string); ok && keyValue != "" {
			errorDocumentValue.Key = aws.String(keyValue)
		}
		error_document = &errorDocumentValue
		return error_document
	} else {
		return nil
	}
}

func indexDocumentSetFunction(indexDocument []interface{}) *s3.IndexDocument {
	var index_document *s3.IndexDocument
	indexDocumentValue := s3.IndexDocument{}
	if len(indexDocument) != 0 {
		indexDocumentMap, _ := indexDocument[0].(map[string]interface{})

		if suffixValue, ok := indexDocumentMap["suffix"].(string); ok && suffixValue != "" {
			indexDocumentValue.Suffix = aws.String(suffixValue)
		}
		index_document = &indexDocumentValue
		return index_document
	} else {
		return nil
	}
}

func redirectAllRequestsSetFunction(redirectAllRequestsSet []interface{}) *s3.RedirectAllRequestsTo {
	var redirect_all_requests *s3.RedirectAllRequestsTo
	redirectAllRequestsValue := s3.RedirectAllRequestsTo{}
	if len(redirectAllRequestsSet) != 0 {
		redirectAllRequestsMap, _ := redirectAllRequestsSet[0].(map[string]interface{})
		if hostName, ok := redirectAllRequestsMap["host_name"].(string); ok && hostName != "" {
			redirectAllRequestsValue.HostName = aws.String(hostName)
		}
		if protocol, ok := redirectAllRequestsMap["protocol"].(string); ok && protocol != "" {
			redirectAllRequestsValue.Protocol = aws.String(protocol)
		}
		redirect_all_requests = &redirectAllRequestsValue
		return redirect_all_requests
	} else {
		return nil
	}
}

func routingRouteConditionSet(routingRuleConditionSet []interface{}) *s3.Condition {
	var routingRuleCondition *s3.Condition
	conditionValue := s3.Condition{}
	if len(routingRuleConditionSet) != 0 {
		routingRuleConditionMap, _ := routingRuleConditionSet[0].(map[string]interface{})
		if httpErrorCodeReturnedEquals, ok := routingRuleConditionMap["http_error_code_returned_equals"].(string); ok && httpErrorCodeReturnedEquals != "" {
			conditionValue.HttpErrorCodeReturnedEquals = aws.String(httpErrorCodeReturnedEquals)
		}
		if keyPrefixEquals, ok := routingRuleConditionMap["key_prefix_equals"].(string); ok && keyPrefixEquals != "" {
			conditionValue.KeyPrefixEquals = aws.String(keyPrefixEquals)
		}
		routingRuleCondition = &conditionValue
		return routingRuleCondition
	} else {
		return nil
	}
}

func routingRouteRedirectSet(routingRuleRedirectSet []interface{}) *s3.Redirect {
	var routingRuleRedirect *s3.Redirect
	redirectValue := s3.Redirect{}
	if len(routingRuleRedirectSet) != 0 {
		routingRuleRedirectMap, _ := routingRuleRedirectSet[0].(map[string]interface{})
		if hostName, ok := routingRuleRedirectMap["host_name"].(string); ok && hostName != "" {
			redirectValue.HostName = aws.String(hostName)
		}
		if httpRedirectCode, ok := routingRuleRedirectMap["http_redirect_code"].(string); ok && httpRedirectCode != "" {
			redirectValue.HttpRedirectCode = aws.String(httpRedirectCode)
		}
		if protocol, ok := routingRuleRedirectMap["protocol"].(string); ok && protocol != "" {
			redirectValue.Protocol = aws.String(protocol)
		}
		if replaceKeyPrefixWith, ok := routingRuleRedirectMap["replace_key_prefix_with"].(string); ok && replaceKeyPrefixWith != "" {
			redirectValue.ReplaceKeyPrefixWith = aws.String(replaceKeyPrefixWith)
		}
		if replaceKeyWith, ok := routingRuleRedirectMap["replace_key_with"].(string); ok && replaceKeyWith != "" {
			redirectValue.ReplaceKeyWith = aws.String(replaceKeyWith)
		}
		routingRuleRedirect = &redirectValue
		return routingRuleRedirect
	} else {
		return nil
	}
}

func routingRuleSetFunction(routingRuleSet []interface{}) []*s3.RoutingRule {
	var rules []*s3.RoutingRule
	for _, l := range routingRuleSet {
		ruleMap, ok := l.(map[string]interface{})
		if !ok {
			continue
		}
		routing_rule := s3.RoutingRule{}
		if condition, ok := ruleMap["condition"].([]interface{}); ok && len(condition) > 0 && condition[0] != nil {
			routing_rule.Condition = routingRouteConditionSet(condition)
		}
		if redirect, ok := ruleMap["redirect"].([]interface{}); ok && len(redirect) > 0 && redirect[0] != nil {
			routing_rule.Redirect = routingRouteRedirectSet(redirect)
		}
		rules = append(rules, &routing_rule)
	}
	return rules
}

func websiteConfigurationSet(websiteConfigurationList []interface{}) (*s3.WebsiteConfiguration, error) {
	var websiteConfig *s3.WebsiteConfiguration
	website_configuration := s3.WebsiteConfiguration{}
	configurationMap, _ := websiteConfigurationList[0].(map[string]interface{})
	// error document is provided
	if errorDocumentSet, exist := configurationMap["error_document"]; exist {
		website_configuration.ErrorDocument = errorDocumentSetFunction(errorDocumentSet.([]interface{}))
	}
	//index document is provided
	if indexDocumentSet, exist := configurationMap["index_document"]; exist {
		website_configuration.IndexDocument = indexDocumentSetFunction(indexDocumentSet.([]interface{}))
	}
	//redirect_all_requests_to is provided
	if redirectAllRequestsSet, exist := configurationMap["redirect_all_requests_to"]; exist {

		website_configuration.RedirectAllRequestsTo = redirectAllRequestsSetFunction(redirectAllRequestsSet.([]interface{}))
	}
	// routing_rules provided
	if routingRulesSet, exist := configurationMap["routing_rule"]; exist {
		website_configuration.RoutingRules = routingRuleSetFunction(routingRulesSet.([]interface{}))
	}
	// if json routing routes are provided
	if routingRulesJsonSet, exist := configurationMap["routing_rules"]; exist {
		jsonValue := fmt.Sprint(configurationMap["routing_rules"])
		if jsonValue != "" {
			var unmarshalledRules []*s3.RoutingRule
			if err := json.Unmarshal([]byte(routingRulesJsonSet.(string)), &unmarshalledRules); err != nil {
				return nil, fmt.Errorf("failed to update the json routing rules in the website configuration : %v", err)
			}
			fmt.Println("unmarshal rules : ", unmarshalledRules)
			website_configuration.RoutingRules = unmarshalledRules
		}
	}
	websiteConfig = &website_configuration
	return websiteConfig, nil
}

func resourceIBMCOSBucketWebsiteConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	bucketCRN := d.Get("bucket_crn").(string)
	bucketName := strings.Split(bucketCRN, ":bucket:")[1]
	instanceCRN := fmt.Sprintf("%s::", strings.Split(bucketCRN, ":bucket:")[0])
	bucketLocation := d.Get("bucket_location").(string)
	endpointType := d.Get("endpoint_type").(string)
	bxSession, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return err
	}
	s3Client, err := getS3ClientSession(bxSession, bucketLocation, endpointType, instanceCRN)
	var websiteConfiguration *s3.WebsiteConfiguration
	configuration, ok := d.GetOk("website_configuration")
	if ok {
		websiteConfiguration, err = websiteConfigurationSet(configuration.([]interface{}))
	}
	if err != nil {
		return fmt.Errorf("failed to read website configuration for COS bucket %s, %v", bucketName, err)
	}
	putBucketWebsiteConfigurationInput := s3.PutBucketWebsiteInput{
		Bucket:               aws.String(bucketName),
		WebsiteConfiguration: websiteConfiguration,
	}
	_, err = s3Client.PutBucketWebsite(&putBucketWebsiteConfigurationInput)

	if err != nil {
		return fmt.Errorf("failed to put website configuration on the COS bucket %s, %v", bucketName, err)
	}
	bktID := fmt.Sprintf("%s:%s:%s:meta:%s:%s", strings.Replace(instanceCRN, "::", "", -1), "bucket", bucketName, bucketLocation, endpointType)
	d.SetId(bktID)
	return resourceIBMCOSBucketWebsiteConfigurationUpdate(d, meta)
}

func resourceIBMCOSBucketWebsiteConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	bucketCRN := d.Get("bucket_crn").(string)
	bucketName := strings.Split(bucketCRN, ":bucket:")[1]
	instanceCRN := fmt.Sprintf("%s::", strings.Split(bucketCRN, ":bucket:")[0])
	bucketLocation := d.Get("bucket_location").(string)
	endpointType := d.Get("endpoint_type").(string)
	bxSession, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return err
	}
	s3Client, err := getS3ClientSession(bxSession, bucketLocation, endpointType, instanceCRN)
	if err != nil {
		return err
	}
	if d.HasChange("website_configuration") {
		var websiteConfiguration *s3.WebsiteConfiguration
		configuration, ok := d.GetOk("website_configuration")
		if ok {
			websiteConfiguration, err = websiteConfigurationSet(configuration.([]interface{}))

		}
		if err != nil {
			return fmt.Errorf("failed to read website configuration for COS bucket %s, %v", bucketName, err)
		}
		putBucketWebsiteConfigurationInput := s3.PutBucketWebsiteInput{
			Bucket:               aws.String(bucketName),
			WebsiteConfiguration: websiteConfiguration,
		}
		_, err = s3Client.PutBucketWebsite(&putBucketWebsiteConfigurationInput)

		if err != nil {
			return fmt.Errorf("failed to update website configuration on the COS bucket %s, %v", bucketName, err)
		}
	}
	return resourceIBMCOSBucketWebsiteConfigurationRead(d, meta)
}

func resourceIBMCOSBucketWebsiteConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	bucketCRN := parseWebsiteId(d.Id(), "bucketCRN")
	bucketName := parseWebsiteId(d.Id(), "bucketName")
	bucketLocation := parseWebsiteId(d.Id(), "bucketLocation")
	instanceCRN := parseWebsiteId(d.Id(), "instanceCRN")
	endpointType := parseWebsiteId(d.Id(), "endpointType")
	d.Set("bucket_crn", bucketCRN)
	d.Set("bucket_location", bucketLocation)
	if endpointType != "" {
		d.Set("endpoint_type", endpointType)
	}
	bxSession, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return err
	}
	s3Client, err := getS3ClientSession(bxSession, bucketLocation, endpointType, instanceCRN)
	if err != nil {
		return err
	}
	//getBucketConfiguration
	getBucketWebsiteConfigurationInput := &s3.GetBucketWebsiteInput{
		Bucket: aws.String(bucketName),
	}
	output, err := s3Client.GetBucketWebsite(getBucketWebsiteConfigurationInput)
	var outputptr *s3.WebsiteConfiguration
	outputptr = (*s3.WebsiteConfiguration)(output)
	if err != nil && !strings.Contains(err.Error(), "AccessDenied: Access Denied") {
		return err
	}
	if output != nil {
		websiteConfiguration := flex.WebsiteConfigurationGet(outputptr)
		if len(websiteConfiguration) > 0 {
			d.Set("website_configuration", websiteConfiguration)
		}
		websiteEndpoint := getWebsiteEndpoint(bucketName, bucketLocation)
		if websiteEndpoint != "" {
			d.Set("website_endpoint", websiteEndpoint)
		}
	}
	return nil
}

func resourceIBMCOSBucketWebsiteConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	bucketName := parseWebsiteId(d.Id(), "bucketName")
	bucketLocation := parseWebsiteId(d.Id(), "bucketLocation")
	instanceCRN := parseWebsiteId(d.Id(), "instanceCRN")
	endpointType := parseWebsiteId(d.Id(), "endpointType")
	bxSession, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return err
	}
	s3Client, err := getS3ClientSession(bxSession, bucketLocation, endpointType, instanceCRN)
	if err != nil {
		return err
	}
	deleteBucketWebsiteInput := &s3.DeleteBucketWebsiteInput{
		Bucket: aws.String(bucketName),
	}
	_, err = s3Client.DeleteBucketWebsite(deleteBucketWebsiteInput)
	if err != nil {
		return fmt.Errorf("failed to delete the objectlock configuration on the COS bucket %s, %v", bucketName, err)
	}
	return nil
}

func parseWebsiteId(id string, info string) string {
	bucketCRN := strings.Split(id, ":meta:")[0]
	meta := strings.Split(id, ":meta:")[1]
	if info == "bucketName" {
		return strings.Split(bucketCRN, ":bucket:")[1]
	}
	if info == "instanceCRN" {
		return fmt.Sprintf("%s::", strings.Split(bucketCRN, ":bucket:")[0])
	}
	if info == "bucketCRN" {
		return bucketCRN
	}
	if info == "bucketLocation" {
		return strings.Split(meta, ":")[0]
	}
	if info == "endpointType" {
		return strings.Split(meta, ":")[1]
	}
	if info == "keyName" {
		return strings.Split(meta, ":key:")[1]
	}
	return parseBucketId(bucketCRN, info)
}

func getWebsiteEndpoint(bucketName string, bucketLocation string) string {
	return fmt.Sprintf("https://%s.s3-web.%s.cloud-object-storage.appdomain.cloud", bucketName, bucketLocation)
}
