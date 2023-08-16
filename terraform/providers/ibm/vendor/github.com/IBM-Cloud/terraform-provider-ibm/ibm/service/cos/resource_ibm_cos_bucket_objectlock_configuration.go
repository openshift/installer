package cos

import (
	"fmt"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMCOSBucketObjectlock() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMCOSBucketObjectlockCreate,
		Read:     resourceIBMCOSBucketObjectlockRead,
		Update:   resourceIBMCOSBucketObjectlockUpdate,
		Delete:   resourceIBMCOSBucketObjectlockDelete,
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
			"object_lock_configuration": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Bucket level object lock settings includes Days, Years, Mode.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"object_lock_enabled": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enable object lock on a COS bucket. This can be used to enable objectlock on an existing bucket",
						},
						"object_lock_rule": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"default_retention": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "An object lock configuration on the object at a bucket level, in the form of a days , years and mode that establishes a point in time after which the object can be deleted. This is applied at bucket level hence it is by default applied to all the object in the bucket unless a seperate retention period is set on the object.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"mode": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Retention modes apply different levels of protection to the objects.",
												},
												"years": {
													Type:          schema.TypeInt,
													Optional:      true,
													ConflictsWith: []string{"object_lock_configuration.0.object_lock_rule.0.default_retention.0.days"},
													Description:   "Retention period in terms of years after which the object can be deleted.",
												},
												"days": {
													Type:          schema.TypeInt,
													Optional:      true,
													ConflictsWith: []string{"object_lock_configuration.0.object_lock_rule.0.default_retention.0.years"},
													Description:   "Retention period in terms of days after which the object can be deleted.",
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
		},
	}
}

func resourceIBMCOSBucketObjectlockCreate(d *schema.ResourceData, meta interface{}) error {

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
	var objectLockConfiguration *s3.ObjectLockConfiguration
	configuration, ok := d.GetOk("object_lock_configuration")
	if ok {
		objectLockConfiguration = objectLockConfigurationSet(configuration.([]interface{}))

	}
	putObjectLockConfigurationInput := &s3.PutObjectLockConfigurationInput{
		Bucket:                  aws.String(bucketName),
		ObjectLockConfiguration: objectLockConfiguration,
	}
	_, err = s3Client.PutObjectLockConfiguration(putObjectLockConfigurationInput)

	if err != nil {
		return fmt.Errorf("failed to put objectlock configuration on the COS bucket %s, %v", bucketName, err)
	}
	bktID := fmt.Sprintf("%s:%s:%s:meta:%s:%s", strings.Replace(instanceCRN, "::", "", -1), "bucket", bucketName, bucketLocation, endpointType)
	d.SetId(bktID)
	return resourceIBMCOSBucketObjectlockUpdate(d, meta)
}

func resourceIBMCOSBucketObjectlockUpdate(d *schema.ResourceData, meta interface{}) error {
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
	if d.HasChange("object_lock_configuration") {
		var objectLockConfiguration *s3.ObjectLockConfiguration
		configuration, ok := d.GetOk("object_lock_configuration")
		if ok {
			objectLockConfiguration = objectLockConfigurationSet(configuration.([]interface{}))
		}
		putObjectLockConfigurationInput := &s3.PutObjectLockConfigurationInput{
			Bucket:                  aws.String(bucketName),
			ObjectLockConfiguration: objectLockConfiguration,
		}
		_, err = s3Client.PutObjectLockConfiguration(putObjectLockConfigurationInput)

		if err != nil {
			return fmt.Errorf("failed to update objectlock configuration on the COS bucket %s, %v", bucketName, err)
		}
	}
	return resourceIBMCOSBucketObjectlockRead(d, meta)
}

func resourceIBMCOSBucketObjectlockRead(d *schema.ResourceData, meta interface{}) error {
	bucketCRN := parseObjectLockId(d.Id(), "bucketCRN")
	bucketName := parseObjectLockId(d.Id(), "bucketName")
	bucketLocation := parseObjectLockId(d.Id(), "bucketLocation")
	instanceCRN := parseObjectLockId(d.Id(), "instanceCRN")
	endpointType := parseObjectLockId(d.Id(), "endpointType")

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
	//object lock configuration
	getObjectLockConfigurationInput := &s3.GetObjectLockConfigurationInput{
		Bucket: aws.String(bucketName),
	}

	output, err := s3Client.GetObjectLockConfiguration(getObjectLockConfigurationInput)

	if err != nil && !strings.Contains(err.Error(), "AccessDenied: Access Denied") {
		return err
	}
	if output.ObjectLockConfiguration != nil {
		objectLockConfigurationptr := output.ObjectLockConfiguration

		objectLockConfiguration := flex.ObjectLockConfigurationGet(objectLockConfigurationptr)
		if len(objectLockConfiguration) > 0 {
			d.Set("object_lock_configuration", objectLockConfiguration)
		}

	}
	return nil
}

func resourceIBMCOSBucketObjectlockDelete(d *schema.ResourceData, meta interface{}) error {
	bucketName := parseObjectLockId(d.Id(), "bucketName")
	bucketLocation := parseObjectLockId(d.Id(), "bucketLocation")
	instanceCRN := parseObjectLockId(d.Id(), "instanceCRN")
	endpointType := parseObjectLockId(d.Id(), "endpointType")

	bxSession, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return err
	}

	s3Client, err := getS3ClientSession(bxSession, bucketLocation, endpointType, instanceCRN)
	if err != nil {
		return err
	}
	putObjectLockConfigurationInput := &s3.PutObjectLockConfigurationInput{
		Bucket: aws.String(bucketName),
		ObjectLockConfiguration: &s3.ObjectLockConfiguration{
			ObjectLockEnabled: aws.String("Enabled"),
		},
	}
	_, err = s3Client.PutObjectLockConfiguration(putObjectLockConfigurationInput)
	if err != nil {
		return fmt.Errorf("failed to put objectlock configuration on the COS bucket %s, %v", bucketName, err)
	}
	return nil
}

func parseObjectLockId(id string, info string) string {
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

func objectLockDefaultRetentionSetFunction(objectLockDefaultRetentionList []interface{}) *s3.DefaultRetention {
	var defaultRetention *s3.DefaultRetention
	var days, years int64
	var mode string

	for _, l := range objectLockDefaultRetentionList {
		object_lock_default_retention := s3.DefaultRetention{}
		defaultRetentionMap, _ := l.(map[string]interface{})

		if modeSet, exist := defaultRetentionMap["mode"]; exist {
			mode = modeSet.(string)
			object_lock_default_retention.Mode = aws.String(mode)

		}
		if daysSet, exist := defaultRetentionMap["days"]; exist {
			objectLockdays := int64(daysSet.(int))
			if objectLockdays > 0 {
				days = objectLockdays
				object_lock_default_retention.Days = aws.Int64(days)
			}
		}
		if yearsSet, exist := defaultRetentionMap["years"]; exist {
			objectLockyears := int64(yearsSet.(int))
			if objectLockyears > 0 {
				years = objectLockyears
				object_lock_default_retention.Years = aws.Int64(years)
			}

		}
		defaultRetention = &object_lock_default_retention
	}

	return defaultRetention
}
func objectLockRuleSetFunction(objectLockRuleList []interface{}) *s3.ObjectLockRule {
	var rules *s3.ObjectLockRule

	for _, l := range objectLockRuleList {
		object_lock_rule := s3.ObjectLockRule{}
		ruleMap, _ := l.(map[string]interface{})

		if defaultRetentionSet, exist := ruleMap["default_retention"]; exist {
			object_lock_rule.DefaultRetention = objectLockDefaultRetentionSetFunction(defaultRetentionSet.([]interface{}))

		}

		rules = &object_lock_rule

	}
	return rules
}
func objectLockConfigurationSet(objectLockConfigurationList []interface{}) *s3.ObjectLockConfiguration {
	var objectLockConfig *s3.ObjectLockConfiguration

	for _, l := range objectLockConfigurationList {
		object_lock_configuration := s3.ObjectLockConfiguration{}
		configurationMap, _ := l.(map[string]interface{})
		//objectlock enabled
		if objectLockEnabledSet, exist := configurationMap["object_lock_enabled"]; exist {
			objectLockEnabledValue := objectLockEnabledSet.(string)
			if objectLockEnabledValue == "enabled" || objectLockEnabledValue == "Enabled" {
				object_lock_configuration.ObjectLockEnabled = aws.String(s3.ObjectLockEnabledEnabled)
			} else {
				object_lock_configuration.ObjectLockEnabled = aws.String(objectLockEnabledValue)
			}

		}
		//ObjectLock configuration
		if objectLockRuleSet, exist := configurationMap["object_lock_rule"]; exist {
			object_lock_configuration.Rule = objectLockRuleSetFunction(objectLockRuleSet.([]interface{}))

		}
		objectLockConfig = &object_lock_configuration

	}
	return objectLockConfig
}
