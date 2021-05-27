// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/IBM/ibm-cos-sdk-go-config/resourceconfigurationv1"
	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/aws/credentials/ibmiam"
	token "github.com/IBM/ibm-cos-sdk-go/aws/credentials/ibmiam/token"
	"github.com/IBM/ibm-cos-sdk-go/aws/session"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var singleSiteLocation = []string{
	"ams03", "che01", "hkg02", "mel01", "mex01",
	"mil01", "mon01", "osl01", "par01", "sjc04", "sao01",
	"seo01", "sng01", "tor01",
}

var regionLocation = []string{
	"au-syd", "ca-tor", "eu-de", "eu-gb", "jp-tok", "jp-osa", "us-east", "us-south",
}

var crossRegionLocation = []string{
	"us", "eu", "ap",
}

var storageClass = []string{
	"standard", "vault", "cold", "flex", "smart",
}

const (
	keyAlgorithm = "AES256"
)

func resourceIBMCOSBucket() *schema.Resource {
	return &schema.Resource{
		Read:     resourceIBMCOSBucketRead,
		Create:   resourceIBMCOSBucketCreate,
		Update:   resourceIBMCOSBucketUpdate,
		Delete:   resourceIBMCOSBucketDelete,
		Exists:   resourceIBMCOSBucketExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "COS Bucket name",
			},
			"resource_instance_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "resource instance ID",
				ValidateFunc: validateRegexp(`^crn:.+:.+:.+:.+:.+:a\/[0-9a-f]{32}:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\:\:$`),
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CRN of resource instance",
			},
			"key_protect": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "CRN of the key you want to use data at rest encryption",
			},
			"single_site_location": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validateAllowedStringValue(singleSiteLocation),
				ForceNew:      true,
				ConflictsWith: []string{"region_location", "cross_region_location"},
				Description:   "single site location info",
			},
			"region_location": {
				Type:     schema.TypeString,
				Optional: true,
				//ValidateFunc:  validateAllowedStringValue(regionLocation),
				ForceNew:      true,
				ConflictsWith: []string{"cross_region_location", "single_site_location"},
				Description:   "Region Location info.",
			},
			"cross_region_location": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validateAllowedStringValue(crossRegionLocation),
				ForceNew:      true,
				ConflictsWith: []string{"region_location", "single_site_location"},
				Description:   "Cros region location info",
			},
			"storage_class": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue(storageClass),
				ForceNew:     true,
				Description:  "Storage class info",
			},
			"endpoint_type": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validateAllowedStringValue([]string{"public", "private"}),
				Description:      "public or private",
				DiffSuppressFunc: applyOnce,
				Default:          "public",
			},
			"s3_endpoint_public": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public endpoint for the COS bucket",
			},
			"s3_endpoint_private": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Private endpoint for the COS bucket",
			},
			"allowed_ip": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    false,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of IPv4 or IPv6 addresses ",
			},
			"activity_tracking": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Enables sending log data to Activity Tracker and LogDNA to provide visibility into object read and write events",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"read_data_events": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "If set to true, all object read events will be sent to Activity Tracker.",
						},
						"write_data_events": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "If set to true, all object write events will be sent to Activity Tracker.",
						},
						"activity_tracker_crn": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The instance of Activity Tracker that will receive object event data",
						},
					},
				},
			},
			"metrics_monitoring": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Enables sending metrics to IBM Cloud Monitoring.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"usage_metrics_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Usage metrics will be sent to the monitoring service.",
						},
						"request_metrics_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Request metrics will be sent to the monitoring service.",
						},
						"metrics_monitoring_crn": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Instance of IBM Cloud Monitoring that will receive the bucket metrics.",
						},
					},
				},
			},
			"archive_rule": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Enable configuration archive_rule (glacier/accelerated) to COS Bucket after a defined period of time",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Unique identifier for the rule.Archive rules allow you to set a specific time frame after which objects transition to the archive. Set Rule ID for cos bucket",
						},
						"enable": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Enable or disable an archive rule for a bucket",
						},
						"days": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validateAllowedRangeInt(0, 3650),
							Description:  "Specifies the number of days when the specific rule action takes effect.",
						},
						"type": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateFunc:     validateAllowedStringValue([]string{"GLACIER", "ACCELERATED", "Glacier", "Accelerated", "glacier", "accelerated"}),
							DiffSuppressFunc: caseDiffSuppress,
							Description:      "Specifies the storage class/archive type to which you want the object to transition. It can be Glacier or Accelerated",
						},
					},
				},
			},
			"expire_rule": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1000,
				Description: "Enable configuration expire_rule to COS Bucket after a defined period of time",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Unique identifier for the rule.Expire rules allow you to set a specific time frame after which objects are deleted. Set Rule ID for cos bucket",
						},
						"enable": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Enable or disable an expire rule for a bucket",
						},
						"prefix": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The rule applies to any objects with keys that match this prefix",
						},
						"days": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validateAllowedRangeInt(0, 3650),
							Description:  "Specifies the number of days when the specific rule action takes effect.",
						},
					},
				},
			},
			"retention_rule": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "A retention policy is enabled at the IBM Cloud Object Storage bucket level. Minimum, maximum and default retention period are defined by this policy and apply to all objects in the bucket.",
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validateAllowedRangeInt(0, 365243),
							Description:  "If an object is stored in the bucket without specifying a custom retention period.",
							ForceNew:     false,
						},
						"maximum": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validateAllowedRangeInt(0, 365243),
							Description:  "Maximum duration of time an object can be kept unmodified in the bucket.",
							ForceNew:     false,
						},
						"minimum": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validateAllowedRangeInt(0, 365243),
							Description:  "Minimum duration of time an object must be kept unmodified in the bucket",
							ForceNew:     false,
						},
						"permanent": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Enable or disable the permanent retention policy on the bucket",
						},
					},
				},
			},
			"force_delete": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "COS buckets need to be empty before they can be deleted. force_delete option empty the bucket and delete it.",
			},
		},
	}
}

func archiveRuleList(archiveList []interface{}) []*s3.LifecycleRule {
	var archive_status, archiveStorageClass, rule_id string
	var days int64
	var rules []*s3.LifecycleRule

	for _, l := range archiveList {
		archiveMap, _ := l.(map[string]interface{})
		//Rule ID
		if rule_idSet, exist := archiveMap["rule_id"]; exist {
			id := rule_idSet.(string)
			rule_id = id
		}

		//Status Enable/Disable
		if archive_statusSet, exist := archiveMap["enable"]; exist {
			archiveStatusEnabled := archive_statusSet.(bool)
			if archiveStatusEnabled == true {
				archive_status = "Enabled"
			} else {
				archive_status = "Disabled"
			}
		}
		//Days
		if daysarchiveSet, exist := archiveMap["days"]; exist {
			daysarchive := int64(daysarchiveSet.(int))
			days = daysarchive
		}
		//Archive Type
		if archiveStorgaeClassSet, exist := archiveMap["type"]; exist {
			archiveType := archiveStorgaeClassSet.(string)
			archiveStorageClass = archiveType
		}

		archive_rule := s3.LifecycleRule{
			ID:     aws.String(rule_id),
			Status: aws.String(archive_status),
			Filter: &s3.LifecycleRuleFilter{},
			Transitions: []*s3.Transition{
				{
					Days:         aws.Int64(days),
					StorageClass: aws.String(archiveStorageClass),
				},
			},
		}

		rules = append(rules, &archive_rule)
	}
	return rules
}

func expireRuleList(expireList []interface{}) []*s3.LifecycleRule {
	var expire_prefix, expire_status, rule_id string
	var days int64
	var rules []*s3.LifecycleRule

	for _, l := range expireList {
		expireMap, _ := l.(map[string]interface{})
		//Rule ID
		if rule_idSet, exist := expireMap["rule_id"]; exist {
			id := rule_idSet.(string)
			rule_id = id
		}

		//Status Enable/Disable
		if expire_statusSet, exist := expireMap["enable"]; exist {
			archiveStatusEnabled := expire_statusSet.(bool)
			if archiveStatusEnabled == true {
				expire_status = "Enabled"
			} else {
				expire_status = "Disabled"
			}
		}
		//Days
		if daysexpireSet, exist := expireMap["days"]; exist {
			daysexpire := int64(daysexpireSet.(int))
			days = daysexpire
		}
		//Expire Prefix
		if expirePrefixClassSet, exist := expireMap["prefix"]; exist {
			expire_prefix = expirePrefixClassSet.(string)
		}

		expire_rule := s3.LifecycleRule{
			ID:     aws.String(rule_id),
			Status: aws.String(expire_status),
			Filter: &s3.LifecycleRuleFilter{
				Prefix: aws.String(expire_prefix),
			},
			Expiration: &s3.LifecycleExpiration{
				Days: aws.Int64(days),
			},
		}

		rules = append(rules, &expire_rule)
	}
	return rules
}

func resourceIBMCOSBucketUpdate(d *schema.ResourceData, meta interface{}) error {
	var s3Conf *aws.Config
	rsConClient, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return err
	}
	bucketName := parseBucketId(d.Id(), "bucketName")
	serviceID := parseBucketId(d.Id(), "serviceID")
	endpointType := parseBucketId(d.Id(), "endpointType")
	apiEndpoint, apiEndpointPrivate := selectCosApi(parseBucketId(d.Id(), "apiType"), parseBucketId(d.Id(), "bLocation"))
	if endpointType == "private" {
		apiEndpoint = apiEndpointPrivate
	}
	authEndpoint, err := rsConClient.Config.EndpointLocator.IAMEndpoint()
	if err != nil {
		return err
	}
	authEndpointPath := fmt.Sprintf("%s%s", authEndpoint, "/identity/token")
	apiKey := rsConClient.Config.BluemixAPIKey
	if apiKey != "" {
		s3Conf = aws.NewConfig().WithEndpoint(envFallBack([]string{"IBMCLOUD_COS_ENDPOINT"}, apiEndpoint)).WithCredentials(ibmiam.NewStaticCredentials(aws.NewConfig(), authEndpointPath, apiKey, serviceID)).WithS3ForcePathStyle(true)
	}
	iamAccessToken := rsConClient.Config.IAMAccessToken
	if iamAccessToken != "" {
		initFunc := func() (*token.Token, error) {
			return &token.Token{
				AccessToken:  rsConClient.Config.IAMAccessToken,
				RefreshToken: rsConClient.Config.IAMRefreshToken,
				TokenType:    "Bearer",
				ExpiresIn:    int64((time.Hour * 248).Seconds()) * -1,
				Expiration:   time.Now().Add(-1 * time.Hour).Unix(),
			}, nil
		}
		s3Conf = aws.NewConfig().WithEndpoint(envFallBack([]string{"IBMCLOUD_COS_ENDPOINT"}, apiEndpoint)).WithCredentials(ibmiam.NewCustomInitFuncCredentials(aws.NewConfig(), initFunc, authEndpointPath, serviceID)).WithS3ForcePathStyle(true)
	}
	s3Sess := session.Must(session.NewSession())
	s3Client := s3.New(s3Sess, s3Conf)

	//// Update  the lifecycle (Archive or Expire)
	if d.HasChange("archive_rule") || d.HasChange("expire_rule") {
		var archive, archive_ok = d.GetOk("archive_rule")
		var expire, expire_ok = d.GetOk("expire_rule")
		var rules []*s3.LifecycleRule
		if archive_ok || expire_ok {
			if expire_ok {
				rules = append(rules, expireRuleList(expire.([]interface{}))...)
			}
			if archive_ok {
				rules = append(rules, archiveRuleList(archive.([]interface{}))...)
			}

			lInput := &s3.PutBucketLifecycleConfigurationInput{
				Bucket: aws.String(bucketName),
				LifecycleConfiguration: &s3.LifecycleConfiguration{
					Rules: rules,
				},
			}
			_, err := s3Client.PutBucketLifecycleConfiguration(lInput)
			if err != nil {
				return fmt.Errorf("failed to update the archive rule on COS bucket %s, %v", bucketName, err)
			}

		} else {
			DelInput := &s3.DeleteBucketLifecycleInput{
				Bucket: aws.String(bucketName),
			}

			delarchive, _ := s3Client.DeleteBucketLifecycleRequest(DelInput)
			err := delarchive.Send()
			if err != nil {
				return err
			}
		}
	}

	//// Update  the Retention policy
	if d.HasChange("retention_rule") {
		var defaultretention, minretention, maxretention int64
		var permanentretention bool
		if retention, ok := d.GetOk("retention_rule"); ok {
			retentionList := retention.([]interface{})
			if len(retentionList) > 1 {
				return fmt.Errorf("Can not more than 1 retention policy")
			}
			for _, l := range retentionList {
				retentionMap, _ := l.(map[string]interface{})
				//Default Days
				if defaultretentionSet, exist := retentionMap["default"]; exist {
					defaultdays := int64(defaultretentionSet.(int))
					defaultretention = defaultdays
				}
				//Maximum Days
				if maxretentionSet, exist := retentionMap["maximum"]; exist {
					maxdays := int64(maxretentionSet.(int))
					maxretention = maxdays
				}
				//Minimum Days
				if minretentionSet, exist := retentionMap["minimum"]; exist {
					mindays := int64(minretentionSet.(int))
					minretention = mindays
				}
				//Permanent Retention Enable/Disable
				if permanentretentionSet, exist := retentionMap["permanent"]; exist {
					permanentretention = permanentretentionSet.(bool)
				}
			}
			// PUT BUCKET PROTECTION CONFIGURATION
			pInput := &s3.PutBucketProtectionConfigurationInput{
				Bucket: aws.String(bucketName),
				ProtectionConfiguration: &s3.ProtectionConfiguration{
					DefaultRetention: &s3.BucketProtectionDefaultRetention{
						Days: aws.Int64(defaultretention),
					},
					MaximumRetention: &s3.BucketProtectionMaximumRetention{
						Days: aws.Int64(maxretention),
					},
					MinimumRetention: &s3.BucketProtectionMinimumRetention{
						Days: aws.Int64(minretention),
					},
					Status:                   aws.String("Retention"),
					EnablePermanentRetention: aws.Bool(permanentretention),
				},
			}
			_, err := s3Client.PutBucketProtectionConfiguration(pInput)
			if err != nil {
				return fmt.Errorf("failed to update the retention rule on COS bucket %s, %v", bucketName, err)
			}
		}
	}

	sess, err := meta.(ClientSession).CosConfigV1API()
	if err != nil {
		return err
	}
	if endpointType == "private" {
		sess.SetServiceURL("https://config.private.cloud-object-storage.cloud.ibm.com/v1")
	}

	hasChanged := false
	updateBucketConfigOptions := &resourceconfigurationv1.UpdateBucketConfigOptions{}

	//BucketName
	bucketName = d.Get("bucket_name").(string)
	updateBucketConfigOptions.Bucket = &bucketName

	if d.HasChange("allowed_ip") {
		firewall := &resourceconfigurationv1.Firewall{}
		var ips = make([]string, 0)
		if ip, ok := d.GetOk("allowed_ip"); ok && ip != nil {
			for _, i := range ip.([]interface{}) {
				ips = append(ips, i.(string))
			}
			firewall.AllowedIp = ips
		} else {
			firewall.AllowedIp = []string{}
		}
		hasChanged = true
		updateBucketConfigOptions.Firewall = firewall
	}

	if d.HasChange("activity_tracking") {
		activityTracker := &resourceconfigurationv1.ActivityTracking{}
		if activity, ok := d.GetOk("activity_tracking"); ok {
			activitylist := activity.([]interface{})
			for _, l := range activitylist {
				activityMap, _ := l.(map[string]interface{})

				//Read event - as its optional check for existence
				if readEvent := activityMap["read_data_events"]; readEvent != nil {
					readSet := readEvent.(bool)
					activityTracker.ReadDataEvents = &readSet
				}

				//Write Event - as its optional check for existence
				if writeEvent := activityMap["write_data_events"]; writeEvent != nil {
					writeSet := writeEvent.(bool)
					activityTracker.WriteDataEvents = &writeSet
				}

				//crn - Required field
				crn := activityMap["activity_tracker_crn"].(string)
				activityTracker.ActivityTrackerCrn = &crn
			}
		}
		hasChanged = true
		updateBucketConfigOptions.ActivityTracking = activityTracker
	}

	if d.HasChange("metrics_monitoring") {
		metricsMonitor := &resourceconfigurationv1.MetricsMonitoring{}
		if metrics, ok := d.GetOk("metrics_monitoring"); ok {
			metricslist := metrics.([]interface{})
			for _, l := range metricslist {
				metricsMap, _ := l.(map[string]interface{})

				//metrics enabled - as its optional check for existence
				if metricsSet := metricsMap["usage_metrics_enabled"]; metricsSet != nil {
					metrics := metricsSet.(bool)
					metricsMonitor.UsageMetricsEnabled = &metrics
				}
				// request metrics enabled - as its optional check for existence
				if metricsSet := metricsMap["request_metrics_enabled"]; metricsSet != nil {
					metrics := metricsSet.(bool)
					metricsMonitor.RequestMetricsEnabled = &metrics
				}
				//crn - Required field
				crn := metricsMap["metrics_monitoring_crn"].(string)
				metricsMonitor.MetricsMonitoringCrn = &crn
			}
		}
		hasChanged = true
		updateBucketConfigOptions.MetricsMonitoring = metricsMonitor
	}

	if hasChanged {
		response, err := sess.UpdateBucketConfig(updateBucketConfigOptions)
		if err != nil {
			return fmt.Errorf("Error Update COS Bucket: %s\n%s", err, response)
		}
	}

	return resourceIBMCOSBucketRead(d, meta)
}

func resourceIBMCOSBucketRead(d *schema.ResourceData, meta interface{}) error {
	var s3Conf *aws.Config
	rsConClient, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return err
	}
	bucketName := parseBucketId(d.Id(), "bucketName")
	serviceID := parseBucketId(d.Id(), "serviceID")
	endpointType := parseBucketId(d.Id(), "endpointType")
	apiEndpoint, apiEndpointPrivate := selectCosApi(parseBucketId(d.Id(), "apiType"), parseBucketId(d.Id(), "bLocation"))
	if endpointType == "private" {
		apiEndpoint = apiEndpointPrivate
	}
	apiEndpoint = envFallBack([]string{"IBMCLOUD_COS_ENDPOINT"}, apiEndpoint)
	authEndpoint, err := rsConClient.Config.EndpointLocator.IAMEndpoint()
	if err != nil {
		return err
	}
	authEndpointPath := fmt.Sprintf("%s%s", authEndpoint, "/identity/token")
	apiKey := rsConClient.Config.BluemixAPIKey
	if apiKey != "" {
		s3Conf = aws.NewConfig().WithEndpoint(apiEndpoint).WithCredentials(ibmiam.NewStaticCredentials(aws.NewConfig(), authEndpointPath, apiKey, serviceID)).WithS3ForcePathStyle(true)
	}
	iamAccessToken := rsConClient.Config.IAMAccessToken
	if iamAccessToken != "" && apiKey == "" {
		initFunc := func() (*token.Token, error) {
			return &token.Token{
				AccessToken:  rsConClient.Config.IAMAccessToken,
				RefreshToken: rsConClient.Config.IAMRefreshToken,
				TokenType:    "Bearer",
				ExpiresIn:    int64((time.Hour * 248).Seconds()) * -1,
				Expiration:   time.Now().Add(-1 * time.Hour).Unix(),
			}, nil
		}
		s3Conf = aws.NewConfig().WithEndpoint(apiEndpoint).WithCredentials(ibmiam.NewCustomInitFuncCredentials(aws.NewConfig(), initFunc, authEndpointPath, serviceID)).WithS3ForcePathStyle(true)
	}
	s3Sess := session.Must(session.NewSession())
	s3Client := s3.New(s3Sess, s3Conf)

	headInput := &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	}
	err = s3Client.WaitUntilBucketExists(headInput)
	if err != nil {
		return fmt.Errorf("failed waiting for bucket %s to be created, %v",
			bucketName, err)
	}

	bucketOutput, err := s3Client.ListBucketsExtended(&s3.ListBucketsExtendedInput{})

	if err != nil {
		return err
	}
	var bLocationConstraint string
	for _, b := range bucketOutput.Buckets {
		if *b.Name == bucketName {
			bLocationConstraint = *b.LocationConstraint
		}
	}

	singleSiteLocationRegex, err := regexp.Compile("^[a-z]{3}[0-9][0-9]-[a-z]{4,8}$")
	if err != nil {
		return err
	}
	regionLocationRegex, err := regexp.Compile("^[a-z]{2}-[a-z]{2,5}-[a-z]{4,8}$")
	if err != nil {
		return err
	}
	crossRegionLocationRegex, err := regexp.Compile("^[a-z]{2}-[a-z]{4,8}$")
	if err != nil {
		return err
	}

	if singleSiteLocationRegex.MatchString(bLocationConstraint) {
		d.Set("single_site_location", strings.Split(bLocationConstraint, "-")[0])
		d.Set("storage_class", strings.Split(bLocationConstraint, "-")[1])
	}
	if regionLocationRegex.MatchString(bLocationConstraint) {
		d.Set("region_location", fmt.Sprintf("%s-%s", strings.Split(bLocationConstraint, "-")[0], strings.Split(bLocationConstraint, "-")[1]))
		d.Set("storage_class", strings.Split(bLocationConstraint, "-")[2])
	}
	if crossRegionLocationRegex.MatchString(bLocationConstraint) {
		d.Set("cross_region_location", strings.Split(bLocationConstraint, "-")[0])
		d.Set("storage_class", strings.Split(bLocationConstraint, "-")[1])
	}

	bucketCRN := fmt.Sprintf("%s:%s:%s", strings.Replace(serviceID, "::", "", -1), "bucket", bucketName)
	d.Set("crn", bucketCRN)
	d.Set("resource_instance_id", serviceID)
	d.Set("bucket_name", bucketName)
	d.Set("s3_endpoint_public", apiEndpoint)
	d.Set("s3_endpoint_private", apiEndpointPrivate)
	if endpointType != "" {
		d.Set("endpoint_type", endpointType)
	}

	getBucketConfigOptions := &resourceconfigurationv1.GetBucketConfigOptions{
		Bucket: &bucketName,
	}

	sess, err := meta.(ClientSession).CosConfigV1API()
	if err != nil {
		return err
	}
	if endpointType == "private" {
		sess.SetServiceURL("https://config.private.cloud-object-storage.cloud.ibm.com/v1")
	}

	bucketPtr, response, err := sess.GetBucketConfig(getBucketConfigOptions)
	if err != nil {
		return fmt.Errorf("Error in getting bucket info rule: %s\n%s", err, response)
	}

	if bucketPtr != nil {

		if bucketPtr.Firewall != nil {
			d.Set("allowed_ip", flattenStringList(bucketPtr.Firewall.AllowedIp))
		}
		if bucketPtr.ActivityTracking != nil {
			d.Set("activity_tracking", flattenActivityTrack(bucketPtr.ActivityTracking))
		}
		if bucketPtr.MetricsMonitoring != nil {
			d.Set("metrics_monitoring", flattenMetricsMonitor(bucketPtr.MetricsMonitoring))
		}
	}
	// Read the lifecycle configuration (archive)

	gInput := &s3.GetBucketLifecycleConfigurationInput{
		Bucket: aws.String(bucketName),
	}

	lifecycleptr, err := s3Client.GetBucketLifecycleConfiguration(gInput)

	if (err != nil && !strings.Contains(err.Error(), "NoSuchLifecycleConfiguration: The lifecycle configuration does not exist")) && (err != nil && bucketPtr != nil && bucketPtr.Firewall != nil && !strings.Contains(err.Error(), "AccessDenied: Access Denied")) {
		return err
	}

	if lifecycleptr != nil {
		archiveRules := archiveRuleGet(lifecycleptr.Rules)
		expireRules := expireRuleGet(lifecycleptr.Rules)
		if len(archiveRules) > 0 {
			d.Set("archive_rule", archiveRules)
		}
		if len(expireRules) > 0 {
			d.Set("expire_rule", expireRules)
		}
	}

	retentionInput := &s3.GetBucketProtectionConfigurationInput{
		Bucket: aws.String(bucketName),
	}
	retentionptr, err := s3Client.GetBucketProtectionConfiguration(retentionInput)

	if err != nil && bucketPtr != nil && bucketPtr.Firewall != nil && !strings.Contains(err.Error(), "AccessDenied: Access Denied") {
		return err
	}

	if retentionptr != nil {
		retentionRules := retentionRuleGet(retentionptr.ProtectionConfiguration)
		if len(retentionRules) > 0 {
			d.Set("retention_rule", retentionRules)
		}
	}

	return nil
}

func resourceIBMCOSBucketCreate(d *schema.ResourceData, meta interface{}) error {
	var s3Conf *aws.Config
	rsConClient, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return err
	}
	bucketName := d.Get("bucket_name").(string)
	storageClass := d.Get("storage_class").(string)
	var bLocation string
	var apiType string
	serviceID := d.Get("resource_instance_id").(string)

	if bucketLocation, ok := d.GetOk("cross_region_location"); ok {
		bLocation = bucketLocation.(string)
		apiType = "crl"
	}
	if bucketLocation, ok := d.GetOk("region_location"); ok {
		bLocation = bucketLocation.(string)
		apiType = "rl"
	}
	if bucketLocation, ok := d.GetOk("single_site_location"); ok {
		bLocation = bucketLocation.(string)
		apiType = "ssl"
	}
	if bLocation == "" {
		return fmt.Errorf("Provide either `cross_region_location` or `region_location` or `single_site_location`")
	}
	lConstraint := fmt.Sprintf("%s-%s", bLocation, storageClass)
	var endpointType = d.Get("endpoint_type").(string)
	apiEndpoint, privateApiEndpoint := selectCosApi(apiType, bLocation)
	if endpointType == "private" {
		apiEndpoint = privateApiEndpoint
	}
	apiEndpoint = envFallBack([]string{"IBMCLOUD_COS_ENDPOINT"}, apiEndpoint)
	if apiEndpoint == "" {
		return fmt.Errorf("The endpoint doesn't exists for given location %s and endpoint type %s", bLocation, endpointType)
	}
	create := &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String(lConstraint),
		},
	}

	if keyprotect, ok := d.GetOk("key_protect"); ok {
		create.IBMSSEKPCustomerRootKeyCrn = aws.String(keyprotect.(string))
		create.IBMSSEKPEncryptionAlgorithm = aws.String(keyAlgorithm)
	}

	authEndpoint, err := rsConClient.Config.EndpointLocator.IAMEndpoint()
	if err != nil {
		return err
	}
	authEndpointPath := fmt.Sprintf("%s%s", authEndpoint, "/identity/token")
	apiKey := rsConClient.Config.BluemixAPIKey
	if apiKey != "" {
		s3Conf = aws.NewConfig().WithEndpoint(apiEndpoint).WithCredentials(ibmiam.NewStaticCredentials(aws.NewConfig(), authEndpointPath, apiKey, serviceID)).WithS3ForcePathStyle(true)
	}
	iamAccessToken := rsConClient.Config.IAMAccessToken
	if iamAccessToken != "" && apiKey == "" {
		initFunc := func() (*token.Token, error) {
			return &token.Token{
				AccessToken:  rsConClient.Config.IAMAccessToken,
				RefreshToken: rsConClient.Config.IAMRefreshToken,
				TokenType:    "Bearer",
				ExpiresIn:    int64((time.Hour * 248).Seconds()) * -1,
				Expiration:   time.Now().Add(-1 * time.Hour).Unix(),
			}, nil
		}
		s3Conf = aws.NewConfig().WithEndpoint(apiEndpoint).WithCredentials(ibmiam.NewCustomInitFuncCredentials(aws.NewConfig(), initFunc, authEndpointPath, serviceID)).WithS3ForcePathStyle(true)
	}

	s3Sess := session.Must(session.NewSession())
	s3Client := s3.New(s3Sess, s3Conf)

	_, err = s3Client.CreateBucket(create)
	if err != nil {
		return err
	}
	// Generating a fake id which contains every information about to get the bucket via s3 api
	bucketID := fmt.Sprintf("%s:%s:%s:meta:%s:%s:%s", strings.Replace(serviceID, "::", "", -1), "bucket", bucketName, apiType, bLocation, endpointType)
	d.SetId(bucketID)

	return resourceIBMCOSBucketUpdate(d, meta)

}

func resourceIBMCOSBucketDelete(d *schema.ResourceData, meta interface{}) error {
	var s3Conf *aws.Config
	rsConClient, _ := meta.(ClientSession).BluemixSession()
	bucketName := parseBucketId(d.Id(), "bucketName")
	serviceID := d.Get("resource_instance_id").(string)
	var bLocation string
	var apiType string
	if bucketLocation, ok := d.GetOk("cross_region_location"); ok {
		bLocation = bucketLocation.(string)
		apiType = "crl"
	}
	if bucketLocation, ok := d.GetOk("region_location"); ok {
		bLocation = bucketLocation.(string)
		apiType = "rl"
	}
	if bucketLocation, ok := d.GetOk("single_site_location"); ok {
		bLocation = bucketLocation.(string)
		apiType = "ssl"
	}
	endpointType := parseBucketId(d.Id(), "endpointType")
	apiEndpoint, apiEndpointPrivate := selectCosApi(apiType, bLocation)
	if endpointType == "private" {
		apiEndpoint = apiEndpointPrivate
	}
	apiEndpoint = envFallBack([]string{"IBMCLOUD_COS_ENDPOINT"}, apiEndpoint)
	if apiEndpoint == "" {
		return fmt.Errorf("The endpoint doesn't exists for given location %s and endpoint type %s", bLocation, endpointType)
	}
	authEndpoint, err := rsConClient.Config.EndpointLocator.IAMEndpoint()
	if err != nil {
		return err
	}
	authEndpointPath := fmt.Sprintf("%s%s", authEndpoint, "/identity/token")

	apiKey := rsConClient.Config.BluemixAPIKey
	if apiKey != "" {
		s3Conf = aws.NewConfig().WithEndpoint(apiEndpoint).WithCredentials(ibmiam.NewStaticCredentials(aws.NewConfig(), authEndpointPath, apiKey, serviceID)).WithS3ForcePathStyle(true)
	}
	iamAccessToken := rsConClient.Config.IAMAccessToken
	if iamAccessToken != "" && apiKey == "" {
		initFunc := func() (*token.Token, error) {
			return &token.Token{
				AccessToken:  rsConClient.Config.IAMAccessToken,
				RefreshToken: rsConClient.Config.IAMRefreshToken,
				TokenType:    "Bearer",
				ExpiresIn:    int64((time.Hour * 248).Seconds()) * -1,
				Expiration:   time.Now().Add(-1 * time.Hour).Unix(),
			}, nil
		}
		s3Conf = aws.NewConfig().WithEndpoint(apiEndpoint).WithCredentials(ibmiam.NewCustomInitFuncCredentials(aws.NewConfig(), initFunc, authEndpointPath, serviceID)).WithS3ForcePathStyle(true)
	}

	s3Sess := session.Must(session.NewSession())
	s3Client := s3.New(s3Sess, s3Conf)

	if delbucket, ok := d.GetOk("force_delete"); ok {
		if delbucket.(bool) {

			// List objects within a bucket
			resp, err := s3Client.ListObjects(&s3.ListObjectsInput{Bucket: aws.String(bucketName)})
			if err != nil {
				return fmt.Errorf("Unable to list items in bucket %s, %v", bucketName, err)
			}
			for _, item := range resp.Contents {
				// Delete object within the bucket
				_, err = s3Client.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(bucketName), Key: aws.String(*item.Key)})

				if err != nil {
					return fmt.Errorf("Unable to delete object %s from bucket %s, %v", *item.Key, bucketName, err)
				}

				err = s3Client.WaitUntilObjectNotExists(&s3.HeadObjectInput{
					Bucket: aws.String(bucketName),
					Key:    aws.String(*item.Key),
				})
				if err != nil {
					return fmt.Errorf("Error occurred while waiting for object %s to be deleted %v", *item.Key, err)
				}
			}
		}
	}

	delete := &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	}
	_, err = s3Client.DeleteBucket(delete)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceIBMCOSBucketExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	var s3Conf *aws.Config
	rsConClient, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return false, err
	}

	bucketName := parseBucketId(d.Id(), "bucketName")
	serviceID := parseBucketId(d.Id(), "serviceID")
	endpointType := parseBucketId(d.Id(), "endpointType")
	apiEndpoint, apiEndpointPrivate := selectCosApi(parseBucketId(d.Id(), "apiType"), parseBucketId(d.Id(), "bLocation"))
	if endpointType == "private" {
		apiEndpoint = apiEndpointPrivate
	}
	apiEndpoint = envFallBack([]string{"IBMCLOUD_COS_ENDPOINT"}, apiEndpoint)
	if apiEndpoint == "" {
		return false, fmt.Errorf("The endpoint doesn't exists for given endpoint type %s", endpointType)
	}
	authEndpoint, err := rsConClient.Config.EndpointLocator.IAMEndpoint()
	if err != nil {
		return false, err
	}
	authEndpointPath := fmt.Sprintf("%s%s", authEndpoint, "/identity/token")

	apiKey := rsConClient.Config.BluemixAPIKey
	if apiKey != "" {
		s3Conf = aws.NewConfig().WithEndpoint(apiEndpoint).WithCredentials(ibmiam.NewStaticCredentials(aws.NewConfig(), authEndpointPath, apiKey, serviceID)).WithS3ForcePathStyle(true)
	}
	iamAccessToken := rsConClient.Config.IAMAccessToken
	if iamAccessToken != "" && apiKey == "" {
		initFunc := func() (*token.Token, error) {
			return &token.Token{
				AccessToken:  rsConClient.Config.IAMAccessToken,
				RefreshToken: rsConClient.Config.IAMRefreshToken,
				TokenType:    "Bearer",
				ExpiresIn:    int64((time.Hour * 248).Seconds()) * -1,
				Expiration:   time.Now().Add(-1 * time.Hour).Unix(),
			}, nil
		}
		s3Conf = aws.NewConfig().WithEndpoint(apiEndpoint).WithCredentials(ibmiam.NewCustomInitFuncCredentials(aws.NewConfig(), initFunc, authEndpointPath, serviceID)).WithS3ForcePathStyle(true)
	}

	s3Sess := session.Must(session.NewSession())
	s3Client := s3.New(s3Sess, s3Conf)

	bucketList, err := s3Client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return false, err
	}
	for _, bucket := range bucketList.Buckets {
		if *bucket.Name == bucketName {
			return true, nil
		}
	}
	return false, nil
}

func selectCosApi(apiType string, bLocation string) (string, string) {
	if apiType == "crl" {
		return fmt.Sprintf("s3.%s.cloud-object-storage.appdomain.cloud", bLocation), fmt.Sprintf("s3.private.%s.cloud-object-storage.appdomain.cloud", bLocation)
	}
	if apiType == "rl" {
		return fmt.Sprintf("s3.%s.cloud-object-storage.appdomain.cloud", bLocation), fmt.Sprintf("s3.private.%s.cloud-object-storage.appdomain.cloud", bLocation)
	}
	if apiType == "ssl" {
		return fmt.Sprintf("s3.%s.cloud-object-storage.appdomain.cloud", bLocation), fmt.Sprintf("s3.private.%s.cloud-object-storage.appdomain.cloud", bLocation)
	}
	return "", ""
}

func parseBucketId(id string, info string) string {
	crn := strings.Split(id, ":meta:")[0]
	meta := strings.Split(id, ":meta:")[1]

	if info == "bucketName" {
		return strings.Split(crn, ":bucket:")[1]
	}
	if info == "serviceID" {
		return fmt.Sprintf("%s::", strings.Split(crn, ":bucket:")[0])
	}
	if info == "apiType" {
		return strings.Split(meta, ":")[0]
	}
	if info == "bLocation" {
		return strings.Split(meta, ":")[1]
	}
	if info == "endpointType" {
		s := strings.Split(meta, ":")
		if len(s) > 2 {
			return strings.Split(meta, ":")[2]
		}
		return ""

	}
	return ""
}
