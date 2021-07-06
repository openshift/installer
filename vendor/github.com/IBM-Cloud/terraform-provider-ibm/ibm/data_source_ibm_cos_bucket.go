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

var bucketTypes = []string{"single_site_location", "region_location", "cross_region_location"}

func dataSourceIBMCosBucket() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMCosBucketRead,

		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bucket_type": {
				Type:         schema.TypeString,
				ValidateFunc: validateAllowedStringValue(bucketTypes),
				Required:     true,
			},
			"bucket_region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"endpoint_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue([]string{"public", "private"}),
				Description:  "public or private",
				Default:      "public",
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CRN of resource instance",
			},
			"key_protect": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CRN of the key you want to use data at rest encryption",
			},
			"single_site_location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cross_region_location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_class": {
				Type:     schema.TypeString,
				Computed: true,
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
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of IPv4 or IPv6 addresses ",
			},
			"activity_tracking": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"read_data_events": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "If set to true, all object read events will be sent to Activity Tracker.",
						},
						"write_data_events": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "If set to true, all object write events will be sent to Activity Tracker.",
						},
						"activity_tracker_crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The instance of Activity Tracker that will receive object event data",
						},
					},
				},
			},
			"metrics_monitoring": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"usage_metrics_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Usage metrics will be sent to the monitoring service.",
						},
						"request_metrics_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Request metrics will be sent to the monitoring service.",
						},
						"metrics_monitoring_crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance of IBM Cloud Monitoring that will receive the bucket metrics.",
						},
					},
				},
			},
			"archive_rule": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Enable configuration archive_rule (glacier/accelerated) to COS Bucket after a defined period of time",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Enable or disable an archive rule for a bucket",
						},
						"days": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"expire_rule": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Enable configuration expire_rule to COS Bucket",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Enable or disable an archive rule for a bucket",
						},
						"days": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the number of days when the specific rule action takes effect.",
						},
						"prefix": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The rule applies to any objects with keys that match this prefix",
						},
					},
				},
			},
			"retention_rule": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A retention policy is enabled at the IBM Cloud Object Storage bucket level. Minimum, maximum and default retention period are defined by this policy and apply to all objects in the bucket.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "If an object is stored in the bucket without specifying a custom retention period.",
						},
						"maximum": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum duration of time an object can be kept unmodified in the bucket.",
						},
						"minimum": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Minimum duration of time an object must be kept unmodified in the bucket",
						},
						"permanent": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Enable or disable the permanent retention policy on the bucket",
						},
					},
				},
			},
			"object_versioning": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Protect objects from accidental deletion or overwrites. Versioning allows you to keep multiple versions of an object protecting from unintentional data loss.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Enable or suspend the versioning for objects in the bucket",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMCosBucketRead(d *schema.ResourceData, meta interface{}) error {
	var s3Conf *aws.Config
	rsConClient, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return err
	}
	bucketName := d.Get("bucket_name").(string)
	serviceID := d.Get("resource_instance_id").(string)
	bucketType := d.Get("bucket_type").(string)
	bucketRegion := d.Get("bucket_region").(string)
	var endpointType = d.Get("endpoint_type").(string)
	apiEndpoint, apiEndpointPrivate := selectCosApi(bucketLocationConvert(bucketType), bucketRegion)
	if endpointType == "private" {
		apiEndpoint = apiEndpointPrivate
	}
	apiEndpoint = envFallBack([]string{"IBMCLOUD_COS_ENDPOINT"}, apiEndpoint)
	if apiEndpoint == "" {
		return fmt.Errorf("The endpoint doesn't exists for given location %s and endpoint type %s", bucketRegion, endpointType)
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
	bucketLocationInput := &s3.GetBucketLocationInput{
		Bucket: aws.String(bucketName),
	}
	bucketLocationConstraint, err := s3Client.GetBucketLocation(bucketLocationInput)
	if err != nil {
		return err
	}
	bLocationConstraint := *bucketLocationConstraint.LocationConstraint

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

	head, err := s3Client.HeadBucket(headInput)
	if err != nil {
		return err
	}
	bucketID := fmt.Sprintf("%s:%s:%s:meta:%s:%s:%s", strings.Replace(serviceID, "::", "", -1), "bucket", bucketName, bucketLocationConvert(bucketType), bucketRegion, endpointType)
	d.SetId(bucketID)
	d.Set("key_protect", head.IBMSSEKPCrkId)
	bucketCRN := fmt.Sprintf("%s:%s:%s", strings.Replace(serviceID, "::", "", -1), "bucket", bucketName)
	d.Set("crn", bucketCRN)
	d.Set("resource_instance_id", serviceID)
	d.Set("s3_endpoint_public", apiEndpoint)
	d.Set("s3_endpoint_private", apiEndpointPrivate)

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
		if len(lifecycleptr.Rules) > 0 {
			archiveRules := archiveRuleGet(lifecycleptr.Rules)
			expireRules := expireRuleGet(lifecycleptr.Rules)
			if len(archiveRules) > 0 {
				d.Set("archive_rule", archiveRules)
			}
			if len(expireRules) > 0 {
				d.Set("expire_rule", expireRules)
			}
		}
	}

	// Read the retention policy
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

	// Get the object Versioning
	versionInput := &s3.GetBucketVersioningInput{
		Bucket: aws.String(bucketName),
	}
	versionPtr, err := s3Client.GetBucketVersioning(versionInput)

	if err != nil && bucketPtr != nil && bucketPtr.Firewall != nil && !strings.Contains(err.Error(), "AccessDenied: Access Denied") {
		return err
	}
	if versionPtr != nil {
		versioningData := flattenCosObejctVersioning(versionPtr)
		if len(versioningData) > 0 {
			d.Set("object_versioning", versioningData)
		}
	}

	return nil
}

func bucketLocationConvert(locationtype string) string {
	if locationtype == "cross_region_location" {
		return "crl"
	}
	if locationtype == "region_location" {
		return "rl"
	}
	if locationtype == "single_site_location" {
		return "crl"
	}
	return ""
}
