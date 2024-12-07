package cos

import (
	// "encoding/json"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validation "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceIBMCOSBucketLifecycleConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMCOSBucketLifecycleConfigurationCreate,
		ReadContext:   resourceIBMCOSBucketLifecycleConfigurationRead,
		UpdateContext: resourceIBMCOSBucketLifecycleConfigurationUpdate,
		DeleteContext: resourceIBMCOSBucketLifecycleConfigurationDelete,
		Importer:      &schema.ResourceImporter{},
		CustomizeDiff: customdiff.Sequence(
			func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
				return resourceValidateLifecycleRule(diff)
			},
		),

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
			"lifecycle_rule": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"abort_incomplete_multipart_upload": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"days_after_initiation": {
										Type:         schema.TypeInt,
										ValidateFunc: validate.ValidateAllowedRangeInt(1, 3650),
										Optional:     true,
									},
								},
							},
						},
						"expiration": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"date": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validate.ValidateUTCFormat,
									},
									"days": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validate.ValidateAllowedRangeInt(1, 3650),
										Default:      0, // API returns 0
									},
									"expired_object_delete_marker": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"filter": {
							Type:             schema.TypeList,
							Required:         true,
							DiffSuppressFunc: suppressMissingFilterConfigurationBlock,
							// IBM has filter as required parameter
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"prefix": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"rule_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 255),
						},
						"noncurrent_version_expiration": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"noncurrent_days": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntAtLeast(1),
									},
								},
							},
						},
						"status": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.ValidateAllowedStringValues([]string{"enable", "disable"}),
						},
						"transition": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"date": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validate.ValidateUTCFormat,
									},
									"days": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validate.ValidateAllowedRangeInt(0, 3650),
									},
									"storage_class": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validate.ValidateAllowedStringValues([]string{"GLACIER",
											"ACCELERATED"}),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func abortIncompleteMultipartUploadSet(abortIncompleteMultipartUpload map[string]interface{}) *s3.AbortIncompleteMultipartUpload {
	var abort_incomplete_multipart_upload *s3.AbortIncompleteMultipartUpload
	abortIncompleteMultipartUploadValue := s3.AbortIncompleteMultipartUpload{}
	if len(abortIncompleteMultipartUpload) != 0 {
		if daysSet, exists := abortIncompleteMultipartUpload["days_after_initiation"]; exists {
			days := int64(daysSet.(int))
			abortIncompleteMultipartUploadValue.DaysAfterInitiation = aws.Int64(days)
		}
		abort_incomplete_multipart_upload = &abortIncompleteMultipartUploadValue
		return abort_incomplete_multipart_upload
	} else {
		return nil
	}
}

func lifecycleExpirationSet(expiration []interface{}) *s3.LifecycleExpiration {
	if len(expiration) == 0 {
		return nil
	}

	var result *s3.LifecycleExpiration

	resultValue := s3.LifecycleExpiration{}

	if expiration[0] == nil {
		return result
	}

	expirationMap := expiration[0].(map[string]interface{})
	if v, ok := expirationMap["date"].(string); ok && v != "" {
		t, _ := time.Parse(time.RFC3339, v)
		resultValue.Date = aws.Time(t)
	}

	if v, ok := expirationMap["days"]; ok {
		daysSet := int64(v.(int))
		if daysSet > 0 {
			days := daysSet
			resultValue.Days = aws.Int64(days)
		}
	}

	// This cannot be specified with Days or Date
	if v, exist := expirationMap["expired_object_delete_marker"]; exist && resultValue.Date == nil && resultValue.Days == nil {
		resultValue.ExpiredObjectDeleteMarker = aws.Bool(v.(bool))
	}

	result = &resultValue

	return result
}

func lifecycleRuleFilterSet(filter []interface{}) *s3.LifecycleRuleFilter {
	var result *s3.LifecycleRuleFilter
	resultValue := s3.LifecycleRuleFilter{}
	if filter[0] == nil {
		resultValue.Prefix = aws.String("")
		result = &resultValue
		return result
	}
	filterMap := filter[0].(map[string]interface{})
	if v, ok := filterMap["prefix"].(string); ok {
		resultValue.Prefix = aws.String(v)
	}
	result = &resultValue
	return result
}

func noncurrentVersionExpirationSet(noncurrentVersionExpiration map[string]interface{}) *s3.NoncurrentVersionExpiration {
	var result *s3.NoncurrentVersionExpiration
	resultValue := s3.NoncurrentVersionExpiration{}
	if v, ok := noncurrentVersionExpiration["noncurrent_days"]; ok {
		resultValue.NoncurrentDays = aws.Int64(int64(v.(int)))
	}
	result = &resultValue
	return result
}

func transitionsSet(transitions []interface{}) []*s3.Transition {
	if len(transitions) == 0 {
		return nil
	}
	var results []*s3.Transition
	var transition = s3.Transition{}
	if transitions[0] == nil {
		return results
	}
	transitionMap := transitions[0].(map[string]interface{})
	if v, ok := transitionMap["date"].(string); ok && v != "" {
		t, _ := time.Parse(time.RFC3339, v)
		transition.Date = aws.Time(t)
	}

	// Only one of "date" and "days" can be configured
	// so only set the transition.Days value when transition.Date is nil
	// By default, transitionMap["days"] = 0 if not explicitly configured in terraform.
	if v, ok := transitionMap["days"]; ok && v.(int) >= 0 && transition.Date == nil {
		transition.Days = aws.Int64(int64(v.(int)))
	}
	if v, ok := transitionMap["storage_class"].(string); ok && v != "" {
		transition.StorageClass = aws.String(v)
	}
	results = append(results, &transition)

	return results
}

func lifecycleConfigurationSet(lifecycleConfigurationList []interface{}) []*s3.LifecycleRule {
	var lifecycleRules []*s3.LifecycleRule
	for _, lifecycleRuleMapRaw := range lifecycleConfigurationList {
		lifecycleRuleMap, ok := lifecycleRuleMapRaw.(map[string]interface{})
		if !ok {
			continue
		}
		lifecycleRule := s3.LifecycleRule{} // single rule to be appended to the list of lifecycle rules
		// check for abort incomplete multipart
		if v, ok := lifecycleRuleMap["abort_incomplete_multipart_upload"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
			lifecycleRule.AbortIncompleteMultipartUpload = abortIncompleteMultipartUploadSet(v[0].(map[string]interface{}))
		}
		// check for expiration
		if v, ok := lifecycleRuleMap["expiration"].([]interface{}); ok && len(v) > 0 {
			lifecycleRule.Expiration = lifecycleExpirationSet(v)
		}
		// check for filter (required)
		if v, ok := lifecycleRuleMap["filter"].([]interface{}); ok && len(v) > 0 {
			lifecycleRule.Filter = lifecycleRuleFilterSet(v)
		}
		// check for rule id (required)
		if v, ok := lifecycleRuleMap["rule_id"].(string); ok {
			lifecycleRule.ID = aws.String(v)
		}
		// check for noncurrent version expiration
		if v, ok := lifecycleRuleMap["noncurrent_version_expiration"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
			lifecycleRule.NoncurrentVersionExpiration = noncurrentVersionExpirationSet(v[0].(map[string]interface{}))
		}
		// check for status
		if v, ok := lifecycleRuleMap["status"].(string); ok && v != "" {
			if v == "enable" {
				lifecycleRule.Status = aws.String("Enabled")
			} else if v == "disable" {
				lifecycleRule.Status = aws.String("Disabled")
			}

		}

		// check for transition
		if v, ok := lifecycleRuleMap["transition"].([]interface{}); ok && len(v) > 0 {
			lifecycleRule.Transitions = transitionsSet(v)
		}

		lifecycleRules = append(lifecycleRules, &lifecycleRule)
	}

	return lifecycleRules
}

func resourceIBMCOSBucketLifecycleConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	bucketCRN := d.Get("bucket_crn").(string)
	bucketName := strings.Split(bucketCRN, ":bucket:")[1]
	instanceCRN := fmt.Sprintf("%s::", strings.Split(bucketCRN, ":bucket:")[0])
	bucketLocation := d.Get("bucket_location").(string)
	endpointType := d.Get("endpoint_type").(string)
	bxSession, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return diag.Errorf("%v", err)
	}
	s3Client, err := getS3ClientSession(bxSession, bucketLocation, endpointType, instanceCRN)
	lifecycleRule := d.Get("lifecycle_rule")
	rules := lifecycleConfigurationSet(lifecycleRule.([]interface{})) // setting each lifecycle rule
	if err != nil {
		return diag.Errorf("Failed to read lifecycle configuration for COS bucket %s, %v", bucketName, err)
	}
	putBucketLifecycleConfigurationInput := s3.PutBucketLifecycleConfigurationInput{
		Bucket: aws.String(bucketName),
		LifecycleConfiguration: &s3.LifecycleConfiguration{
			Rules: rules,
		},
	}
	_, err = s3Client.PutBucketLifecycleConfiguration(&putBucketLifecycleConfigurationInput)
	if err != nil {
		return diag.Errorf("Failed to put Lifecycle configuration on the COS bucket %s, %v", bucketName, err)
	}
	bktID := fmt.Sprintf("%s:%s:%s:meta:%s:%s", strings.Replace(instanceCRN, "::", "", -1), "bucket", bucketName, bucketLocation, endpointType)
	d.SetId(bktID)
	return resourceIBMCOSBucketLifecycleConfigurationUpdate(ctx, d, meta)
}

func resourceIBMCOSBucketLifecycleConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	bucketCRN := d.Get("bucket_crn").(string)
	bucketName := strings.Split(bucketCRN, ":bucket:")[1]
	instanceCRN := fmt.Sprintf("%s::", strings.Split(bucketCRN, ":bucket:")[0])
	bucketLocation := d.Get("bucket_location").(string)
	endpointType := d.Get("endpoint_type").(string)
	bxSession, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return diag.Errorf("%v", err)
	}
	s3Client, err := getS3ClientSession(bxSession, bucketLocation, endpointType, instanceCRN)
	if d.HasChange("lifecycle_rule") {
		lifecycleRule := d.Get("lifecycle_rule")
		rules := lifecycleConfigurationSet(lifecycleRule.([]interface{}))
		if err != nil {
			return diag.Errorf("Failed to read lifecycle configuration for COS bucket %s, %v", bucketName, err)
		}
		putBucketLifecycleConfigurationInput := s3.PutBucketLifecycleConfigurationInput{
			Bucket: aws.String(bucketName),
			LifecycleConfiguration: &s3.LifecycleConfiguration{
				Rules: rules,
			},
		}
		_, err = s3Client.PutBucketLifecycleConfiguration(&putBucketLifecycleConfigurationInput)
		if err != nil {
			return diag.Errorf("Failed to put Lifecycle configuration on the COS bucket %s, %v", bucketName, err)
		}
	}
	return resourceIBMCOSBucketLifecycleConfigurationRead(ctx, d, meta)
}

func resourceIBMCOSBucketLifecycleConfigurationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return diag.Errorf("%v", err)
	}
	s3Client, err := getS3ClientSession(bxSession, bucketLocation, endpointType, instanceCRN)
	if err != nil {
		return diag.Errorf("%v", err)
	}
	//getBucketConfiguration
	const (
		lifecycleConfigurationRulesSteadyTimeout = 2 * time.Minute
	)
	getLifecycleConfigurationInput := &s3.GetBucketLifecycleConfigurationInput{
		Bucket: aws.String(bucketName),
	}
	var output *s3.GetBucketLifecycleConfigurationOutput

	// Adding a retry to overcome the NoSuchLifecycleConfiguration error as it takes time for the lifecycle rules to  get applied.
	err = resource.Retry(lifecycleConfigurationRulesSteadyTimeout, func() *resource.RetryError {
		var err error
		output, err = s3Client.GetBucketLifecycleConfiguration(getLifecycleConfigurationInput)

		if d.IsNewResource() && err != nil && strings.Contains(err.Error(), "NoSuchLifecycleConfiguration: The lifecycle configuration does not exist") {

			return resource.RetryableError(err)
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if conns.IsResourceTimeoutError(err) {
		output, err = s3Client.GetBucketLifecycleConfiguration(getLifecycleConfigurationInput)
	}
	if err != nil {
		return diag.Errorf("[ERROR] Error getting Lifecycle Configuration for the bucket %s", bucketName)
	}
	var outputptr *s3.LifecycleConfiguration
	outputptr = (*s3.LifecycleConfiguration)(output)
	if err != nil && !strings.Contains(err.Error(), "AccessDenied: Access Denied") {
		return diag.Errorf("%v", err)
	}
	if outputptr.Rules != nil {
		lifecycleConfiguration := flex.LifecylceRuleGet(outputptr.Rules)
		if len(lifecycleConfiguration) > 0 {
			d.Set("lifecycle_rule", lifecycleConfiguration)
		}
	} else {
		d.SetId("")
	}
	return nil
}

func resourceIBMCOSBucketLifecycleConfigurationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	bucketName := parseWebsiteId(d.Id(), "bucketName")
	bucketLocation := parseWebsiteId(d.Id(), "bucketLocation")
	instanceCRN := parseWebsiteId(d.Id(), "instanceCRN")
	endpointType := parseWebsiteId(d.Id(), "endpointType")
	bxSession, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return diag.Errorf("%v", err)
	}
	s3Client, err := getS3ClientSession(bxSession, bucketLocation, endpointType, instanceCRN)
	if err != nil {
		return diag.Errorf("%v", err)
	}
	deleteBucketLifecycleInput := &s3.DeleteBucketLifecycleInput{
		Bucket: aws.String(bucketName),
	}
	_, err = s3Client.DeleteBucketLifecycle(deleteBucketLifecycleInput)
	if err != nil {
		return diag.Errorf("failed to delete the Lifecycle configuration on the COS bucket %s, %v", bucketName, err)
	}
	return nil
}

func parseLifecycleId(id string, info string) string {
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

func resourceValidateLifecycleRule(diff *schema.ResourceDiff) error {
	if lifecycle, ok := diff.GetOk("lifecycle_rule"); ok {
		lifecycle_list := lifecycle.([]interface{})
		for _, l := range lifecycle_list {
			lifecycleMap, _ := l.(map[string]interface{})
			expiration := lifecycleMap["expiration"]
			expirationParameterCounter := 0
			if expiration != nil {
				expirationList, _ := expiration.([]interface{})
				for _, e := range expirationList {
					expirationMap, _ := e.(map[string]interface{})
					if val, days_exist := expirationMap["days"]; days_exist && val.(int) != 0 {
						expirationParameterCounter++
					}
					if val, date_exist := expirationMap["date"]; date_exist && val != "" {
						expirationParameterCounter++
					}
					if val, expired_object_delete_marker_exist := expirationMap["expired_object_delete_marker"]; expired_object_delete_marker_exist && val.(bool) != false {
						expirationParameterCounter++
					}
					if expirationParameterCounter > 1 {
						return fmt.Errorf("[ERROR] The expiry 3 action elements (Days, Date, ExpiredObjectDeleteMarker) are all mutually exclusive. These can not be used with each other. Please set one expiry element on the same rule of expiration.")
					}
				}
			}

		}
	}
	return nil
}

func suppressMissingFilterConfigurationBlock(k, old, new string, d *schema.ResourceData) bool {
	if strings.HasSuffix(k, "filter.#") {
		oraw, nraw := d.GetChange(k)
		o, n := oraw.(int), nraw.(int)
		if o == 1 && n == 0 {
			return true
		}
		if o == 1 && n == 1 {
			return old == "1" && new == "0"
		}
		return false
	}
	return false
}
