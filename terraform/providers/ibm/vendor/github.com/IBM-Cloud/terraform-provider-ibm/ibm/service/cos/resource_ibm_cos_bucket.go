// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cos

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/ibm-cos-sdk-go-config/resourceconfigurationv1"
	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/aws/credentials/ibmiam"
	token "github.com/IBM/ibm-cos-sdk-go/aws/credentials/ibmiam/token"
	"github.com/IBM/ibm-cos-sdk-go/aws/session"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var singleSiteLocation = []string{
	"ams03", "che01", "hkg02", "mel01", "mex01",
	"mil01", "mon01", "osl01", "par01", "sjc04", "sao01",
	"seo01", "sng01", "tor01",
}

var regionLocation = []string{
	"au-syd", "ca-tor", "eu-de", "eu-gb", "jp-tok", "jp-osa", "us-east", "us-south", "br-sao",
}

var crossRegionLocation = []string{
	"us", "eu", "ap",
}

var storageClass = []string{
	"standard", "vault", "cold", "smart", "flex", "onerate_active",
}

var singleSiteLocationRegex = regexp.MustCompile("^[a-z]{3}[0-9][0-9]-[a-z]{4,8}$")
var regionLocationRegex = regexp.MustCompile("^[a-z]{2}-[a-z]{2,5}[0-9]?-[a-z]{4,8}$")
var crossRegionLocationRegex = regexp.MustCompile("^[a-z]{2}-[a-z]{4,8}$")

const (
	keyAlgorithm = "AES256"
)

func caseDiffSuppress(_, old, new string, _ *schema.ResourceData) bool {
	return strings.ToUpper(old) == strings.ToUpper(new)
}

func resourceinstanceidDiffSuppress(_, old, new string, _ *schema.ResourceData) bool {
	if old != "" && strings.Contains(new, old) {
		return true
	}
	return false
}
func ResourceIBMCOSBucket() *schema.Resource {
	return &schema.Resource{
		Read:          resourceIBMCOSBucketRead,
		Create:        resourceIBMCOSBucketCreate,
		Update:        resourceIBMCOSBucketUpdate,
		Delete:        resourceIBMCOSBucketDelete,
		Exists:        resourceIBMCOSBucketExists,
		Importer:      &schema.ResourceImporter{},
		CustomizeDiff: resourceExpiryValidate,

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
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "resource instance ID",
				DiffSuppressFunc: resourceinstanceidDiffSuppress,
				ValidateFunc:     validate.InvokeValidator("ibm_cos_bucket", "resource_instance_id"),
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CRN of resource instance",
			},
			"key_protect": {
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				ConflictsWith: []string{"kms_key_crn"},
				Description:   "CRN of the key you want to use data at rest encryption",
			},
			"kms_key_crn": {
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				ConflictsWith: []string{"key_protect"},
				Description:   "CRN of the key you want to use data at rest encryption",
			},
			"satellite_location_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"cross_region_location", "single_site_location", "region_location"},
				Description:   "Provide satellite location info.",
			},
			"single_site_location": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validate.InvokeValidator("ibm_cos_bucket", "single_site_location"),
				ForceNew:      true,
				ConflictsWith: []string{"region_location", "cross_region_location", "satellite_location_id"},
				Description:   "single site location info",
			},
			"region_location": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"cross_region_location", "single_site_location", "satellite_location_id"},
				Description:   "Region Location info.",
			},
			"cross_region_location": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validate.InvokeValidator("ibm_cos_bucket", "cross_region_location"),
				ForceNew:      true,
				ConflictsWith: []string{"region_location", "single_site_location", "satellite_location_id"},
				Description:   "Cros region location info",
			},
			"storage_class": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				Description:   "Storage class info",
				ConflictsWith: []string{"satellite_location_id"},
				ValidateFunc:  validate.InvokeValidator("ibm_cos_bucket", "storage_class"),
			},
			"endpoint_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "public or private",
				ConflictsWith:    []string{"satellite_location_id"},
				DiffSuppressFunc: flex.ApplyOnce,
				Default:          "public",
				ValidateFunc:     validate.InvokeValidator("ibm_cos_bucket", "endpoint_type"),
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
			"s3_endpoint_direct": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Direct endpoint for the COS bucket",
			},
			"allowed_ip": {
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      false,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"satellite_location_id"},
				Description:   "List of IPv4 or IPv6 addresses ",
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
			"abort_incomplete_multipart_upload_days": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Enable abort incomplete multipart upload to COS Bucket after a defined period of time",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Unique identifier for the rule. Rules allow you to set a specific time frame after which objects are deleted. Set Rule ID for cos bucket",
						},
						"enable": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Enable or disable rule for a bucket",
						},
						"prefix": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The rule applies to any objects with keys that match this prefix",
						},
						"days_after_initiation": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validate.ValidateAllowedRangeInt(1, 3650),
							Description:  "Specifies the number of days when the specific rule action takes effect.",
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
							ValidateFunc: validate.ValidateAllowedRangeInt(0, 3650),
							Description:  "Specifies the number of days when the specific rule action takes effect.",
						},
						"type": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateFunc:     validate.InvokeValidator("ibm_cos_bucket", "type"),
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
						"date": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validate.ValidBucketLifecycleTimestamp,
							Description:  "Specify a rule to expire the current version of objects in bucket after a specific date.",
						},
						"days": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validate.ValidateAllowedRangeInt(1, 3650),
							Description:  "Specifies the number of days when the specific rule action takes effect.",
						},
						"expired_object_delete_marker": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Expired object delete markers can be automatically cleaned up to improve performance in bucket. This cannot be used alongside version expiration.",
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
							ValidateFunc: validate.ValidateAllowedRangeInt(0, 365243),
							Description:  "If an object is stored in the bucket without specifying a custom retention period.",
							ForceNew:     false,
						},
						"maximum": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validate.ValidateAllowedRangeInt(0, 365243),
							Description:  "Maximum duration of time an object can be kept unmodified in the bucket.",
							ForceNew:     false,
						},
						"minimum": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validate.ValidateAllowedRangeInt(0, 365243),
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
			"object_versioning": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"retention_rule"},
				Description:   "Protect objects from accidental deletion or overwrites. Versioning allows you to keep multiple versions of an object protecting from unintentional data loss.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Enable or suspend the versioning for objects in the bucket",
						},
					},
				},
			},
			"noncurrent_version_expiration": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
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
						"noncurrent_days": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validate.ValidateAllowedRangeInt(1, 3650),
							Description:  "Specifies the number of days when the specific rule action takes effect.",
						},
					},
				},
			},
			"hard_quota": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "sets a maximum amount of storage (in bytes) available for a bucket",
			},
			"force_delete": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "COS buckets need to be empty before they can be deleted. force_delete option empty the bucket and delete it.",
			},
			"object_lock": {
				Type:         schema.TypeBool,
				Optional:     true,
				RequiredWith: []string{"object_versioning"},
				Description:  "Enable objectlock for the bucket. When enabled, buckets within the container vault can have Object Lock Configuration applied to the bucket.",
			},
		},
	}
}
func ResourceIBMCOSBucketValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "resource_instance_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^crn:.+:.+:.+:.+:.+:a\/[0-9a-f]{32}:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\:\:$`,
			CloudDataType:              "resource_instance",
			CloudDataRange:             []string{"service:cloud-object-storage"}})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "single_site_location",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "ams03,che01,hkg02,mel01,mex01,mil01,mon01,osl01,par01,sjc04,sao01,seo01,sng01,tor01",
		})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cross_region_location",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "us,eu,ap",
		})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "storage_class",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "standard,vault,cold,smart,flex,onerate_active",
		})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "endpoint_type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "public,private,direct",
		})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "GLACIER,ACCELERATED,Glacier,Accelerated,glacier,accelerated",
		})

	ibmCOSBucketResourceValidator := validate.ResourceValidator{ResourceName: "ibm_cos_bucket", Schema: validateSchema}
	return &ibmCOSBucketResourceValidator
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

func nc_expRuleList(nc_expList []interface{}) []*s3.LifecycleRule {
	var nc_exp_prefix, nc_exp_status, rule_id string
	var nc_days int64
	var rules []*s3.LifecycleRule

	for _, l := range nc_expList {
		nc_expMap, _ := l.(map[string]interface{})
		//Rule ID
		if rule_idSet, exist := nc_expMap["rule_id"]; exist {
			id := rule_idSet.(string)
			rule_id = id
		}

		//Status Enable/Disable
		if nc_exp_statusSet, exist := nc_expMap["enable"]; exist {
			nc_expStatusEnabled := nc_exp_statusSet.(bool)
			if nc_expStatusEnabled == true {
				nc_exp_status = "Enabled"
			} else {
				nc_exp_status = "Disabled"
			}
		}
		//Days
		if nc_exp_daySet, exist := nc_expMap["noncurrent_days"]; exist {
			nc_exp_days := int64(nc_exp_daySet.(int))
			nc_days = nc_exp_days
		}
		//Expire Prefix
		if nc_expPrefixClassSet, exist := nc_expMap["prefix"]; exist {
			prefix_check := nc_expPrefixClassSet.(string)
			nc_exp_prefix = prefix_check
		}

		nc_exp_rule := s3.LifecycleRule{
			ID:     aws.String(rule_id),
			Status: aws.String(nc_exp_status),
			Filter: &s3.LifecycleRuleFilter{
				Prefix: aws.String(nc_exp_prefix),
			},
			NoncurrentVersionExpiration: &s3.NoncurrentVersionExpiration{
				NoncurrentDays: aws.Int64(nc_days),
			},
		}
		rules = append(rules, &nc_exp_rule)
	}
	return rules
}

func abortmpuRuleList(abortmpuList []interface{}) []*s3.LifecycleRule {
	var abort_mpu_prefix, abort_mpu_status, rule_id string
	var abort_mpu_days_init int64
	var rules []*s3.LifecycleRule

	for _, l := range abortmpuList {
		abortmpuMap, _ := l.(map[string]interface{})
		//Rule ID
		if rule_idSet, exist := abortmpuMap["rule_id"]; exist {
			id := rule_idSet.(string)
			rule_id = id
		}

		//Status Enable/Disable
		if abort_mpu_statusSet, exist := abortmpuMap["enable"]; exist {
			abort_mpuStatusEnabled := abort_mpu_statusSet.(bool)
			if abort_mpuStatusEnabled == true {
				abort_mpu_status = "Enabled"
			} else {
				abort_mpu_status = "Disabled"
			}
		}
		//Days
		if abort_mpu_daySet, exist := abortmpuMap["days_after_initiation"]; exist {
			abort_mpu_days := int64(abort_mpu_daySet.(int))
			abort_mpu_days_init = abort_mpu_days
		}
		//Expire Prefix
		if abort_mpuPrefixClassSet, exist := abortmpuMap["prefix"]; exist {
			prefix_check := abort_mpuPrefixClassSet.(string)
			abort_mpu_prefix = prefix_check
		}

		abort_mpu_rule := s3.LifecycleRule{
			ID:     aws.String(rule_id),
			Status: aws.String(abort_mpu_status),
			Filter: &s3.LifecycleRuleFilter{
				Prefix: aws.String(abort_mpu_prefix),
			},
			AbortIncompleteMultipartUpload: &s3.AbortIncompleteMultipartUpload{
				DaysAfterInitiation: aws.Int64(abort_mpu_days_init),
			},
		}
		rules = append(rules, &abort_mpu_rule)
	}
	return rules
}

func expireRuleList(expireList []interface{}) []*s3.LifecycleRule {
	var expire_prefix, expire_status, rule_id string
	var expire_date time.Time
	var days int64
	var expired_object_del_marker bool
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
			expireStatusEnabled := expire_statusSet.(bool)
			if expireStatusEnabled == true {
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
		//Date
		if dateexpireSet, exist := expireMap["date"]; exist {
			expiredatevalue := dateexpireSet.(string)
			expire_date, _ = time.Parse(time.RFC3339, fmt.Sprintf("%sT00:00:00Z", expiredatevalue))
		}
		//Expire Prefix
		if expirePrefixClassSet, exist := expireMap["prefix"]; exist {
			prefix := expirePrefixClassSet.(string)
			expire_prefix = prefix
		}
		// Expired Object Delete Marker
		if expireObjectDelMarkerSet, exist := expireMap["expired_object_delete_marker"]; exist {
			expired_object_del_marker = expireObjectDelMarkerSet.(bool)
		}
		var i *s3.LifecycleExpiration

		if expired_object_del_marker == true {
			i = &s3.LifecycleExpiration{
				ExpiredObjectDeleteMarker: aws.Bool(expired_object_del_marker),
			}
		} else if days > 0 {
			i = &s3.LifecycleExpiration{
				Days: aws.Int64(days),
			}
		} else if !expire_date.IsZero() {
			i = &s3.LifecycleExpiration{
				Date: aws.Time(expire_date),
			}
		}
		expire_rule := s3.LifecycleRule{
			ID:     aws.String(rule_id),
			Status: aws.String(expire_status),
			Filter: &s3.LifecycleRuleFilter{
				Prefix: aws.String(expire_prefix),
			},

			Expiration: i,
		}
		rules = append(rules, &expire_rule)
	}
	return rules
}

func resourceIBMCOSBucketUpdate(d *schema.ResourceData, meta interface{}) error {
	var s3Conf *aws.Config
	rsConClient, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return err
	}
	bucketName := parseBucketId(d.Id(), "bucketName")
	serviceID := parseBucketId(d.Id(), "serviceID")
	endpointType := parseBucketId(d.Id(), "endpointType")
	bLocation := parseBucketId(d.Id(), "bLocation")
	apiType := parseBucketId(d.Id(), "apiType")
	if apiType == "sl" {
		satloc_guid := strings.Split(serviceID, ":")
		bucketsatcrn := satloc_guid[0]
		serviceID = bucketsatcrn
	}

	var apiEndpoint, apiEndpointPrivate, directApiEndpoint string

	if apiType == "sl" {
		apiEndpoint = SelectSatlocCosApi(apiType, serviceID, bLocation)
	} else {
		apiEndpoint, apiEndpointPrivate, directApiEndpoint = SelectCosApi(apiType, bLocation)
		if endpointType == "private" {
			apiEndpoint = apiEndpointPrivate
		}
		if endpointType == "direct" {
			apiEndpoint = directApiEndpoint
		}

	}

	authEndpoint, err := rsConClient.Config.EndpointLocator.IAMEndpoint()

	if err != nil {
		return err
	}
	authEndpointPath := fmt.Sprintf("%s%s", authEndpoint, "/identity/token")
	apiKey := rsConClient.Config.BluemixAPIKey
	if apiKey != "" {
		s3Conf = aws.NewConfig().WithEndpoint(conns.EnvFallBack([]string{"IBMCLOUD_COS_ENDPOINT"}, apiEndpoint)).WithCredentials(ibmiam.NewStaticCredentials(aws.NewConfig(), authEndpointPath, apiKey, serviceID)).WithS3ForcePathStyle(true)
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
		s3Conf = aws.NewConfig().WithEndpoint(conns.EnvFallBack([]string{"IBMCLOUD_COS_ENDPOINT"}, apiEndpoint)).WithCredentials(ibmiam.NewCustomInitFuncCredentials(aws.NewConfig(), initFunc, authEndpointPath, serviceID)).WithS3ForcePathStyle(true)
	}
	s3Sess := session.Must(session.NewSession())
	s3Client := s3.New(s3Sess, s3Conf)

	//// Update  the lifecycle (Archive or Expire or Non Current version or Abort incomplete Multipart Upload)
	if d.HasChange("archive_rule") || d.HasChange("expire_rule") || d.HasChange("noncurrent_version_expiration") || d.HasChange("abort_incomplete_multipart_upload_days") {
		var archive, archive_ok = d.GetOk("archive_rule")
		var expire, expire_ok = d.GetOk("expire_rule")
		var noncurrentverexp, nc_exp_ok = d.GetOk("noncurrent_version_expiration")
		var abortmpu, abort_mpu_ok = d.GetOk("abort_incomplete_multipart_upload_days")
		var rules []*s3.LifecycleRule
		if archive_ok || expire_ok || nc_exp_ok || abort_mpu_ok {
			if archive_ok {
				rules = append(rules, archiveRuleList(archive.([]interface{}))...)
			}
			if nc_exp_ok {
				rules = append(rules, nc_expRuleList(noncurrentverexp.([]interface{}))...)
			}
			if abort_mpu_ok {
				rules = append(rules, abortmpuRuleList(abortmpu.([]interface{}))...)
			}
			if expire_ok {
				rules = append(rules, expireRuleList(expire.([]interface{}))...)
			}
			lInput := &s3.PutBucketLifecycleConfigurationInput{
				Bucket: aws.String(bucketName),
				LifecycleConfiguration: &s3.LifecycleConfiguration{
					Rules: rules,
				},
			}

			_, err := s3Client.PutBucketLifecycleConfiguration(lInput)
			if err != nil {
				return fmt.Errorf("failed to update the lifecyle rule on COS bucket %s, %v", bucketName, err)
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

	//update the object versioning (object versioning)
	if d.HasChange("object_versioning") {
		versioningConf := &s3.VersioningConfiguration{}
		if versioning, ok := d.GetOk("object_versioning"); ok {
			versioningList := versioning.([]interface{})
			for _, l := range versioningList {
				versioningMap, _ := l.(map[string]interface{})
				//Status Enable/Disable
				if object_versioning_statusSet, exist1 := versioningMap["enable"]; exist1 {
					versioningStatusEnabled := object_versioning_statusSet.(bool)
					if versioningStatusEnabled == true {
						versioningConf.Status = aws.String("Enabled")
					} else {
						versioningConf.Status = aws.String("Suspended")
					}
				}
			}
		} else {
			versioningConf.Status = aws.String("Suspended")
		}
		// PUT BUCKET Object Versioning
		input := &s3.PutBucketVersioningInput{
			Bucket:                  aws.String(bucketName),
			VersioningConfiguration: versioningConf,
		}

		_, err := s3Client.PutBucketVersioning(input)
		if err != nil {
			return fmt.Errorf("failed to update the object versioning on COS bucket %s, %v", bucketName, err)
		}
	}

	sess, err := meta.(conns.ClientSession).CosConfigV1API()
	if err != nil {
		return err
	}
	if endpointType == "private" {
		sess.SetServiceURL("https://config.private.cloud-object-storage.cloud.ibm.com/v1")
	}
	if endpointType == "direct" {
		sess.SetServiceURL("https://config.direct.cloud-object-storage.cloud.ibm.com/v1")
	}

	if apiType == "sl" {
		satconfig := fmt.Sprintf("https://config.%s.%s.cloud-object-storage.appdomain.cloud/v1", serviceID, bLocation)

		sess.SetServiceURL(satconfig)
	}

	hasChanged := false
	updateBucketConfigOptions := &resourceconfigurationv1.UpdateBucketConfigOptions{}

	//BucketName
	bucketName = d.Get("bucket_name").(string)
	updateBucketConfigOptions.Bucket = &bucketName

	if d.HasChange("hard_quota") {
		hasChanged = true
		updateBucketConfigOptions.HardQuota = aws.Int64(int64(d.Get("hard_quota").(int)))
	}

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
			return fmt.Errorf("[ERROR] Error Update COS Bucket: %s\n%s", err, response)
		}
	}

	return resourceIBMCOSBucketRead(d, meta)
}

func resourceIBMCOSBucketRead(d *schema.ResourceData, meta interface{}) error {
	var s3Conf *aws.Config
	var keyProtectFlag bool
	rsConClient, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return err
	}
	bucketName := parseBucketId(d.Id(), "bucketName")
	serviceID := parseBucketId(d.Id(), "serviceID")
	endpointType := parseBucketId(d.Id(), "endpointType")
	apiType := parseBucketId(d.Id(), "apiType")
	bLocation := parseBucketId(d.Id(), "bLocation")

	if _, ok := d.GetOk("key_protect"); ok {
		keyProtectFlag = true
	}

	//split satellite resource instance id to get the 1st value
	if apiType == "sl" {
		satloc_guid := strings.Split(serviceID, ":")
		bucketsatcrn := satloc_guid[0]
		serviceID = bucketsatcrn
	}

	var apiEndpoint, apiEndpointPrivate, directApiEndpoint string

	if apiType == "sl" {
		apiEndpoint = SelectSatlocCosApi(apiType, serviceID, bLocation)
	} else {
		apiEndpoint, apiEndpointPrivate, directApiEndpoint = SelectCosApi(apiType, bLocation)
		if endpointType == "private" {
			apiEndpoint = apiEndpointPrivate
		}
		if endpointType == "direct" {
			apiEndpoint = directApiEndpoint
		}

	}

	apiEndpoint = conns.EnvFallBack([]string{"IBMCLOUD_COS_ENDPOINT"}, apiEndpoint)

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

	if apiType != "sl" {
		var bLocationConstraint string
		for _, b := range bucketOutput.Buckets {
			if *b.Name == bucketName {
				bLocationConstraint = *b.LocationConstraint
			}
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
	} else {
		d.Set("satellite_location_id", bLocation)

	}

	bucketCRN := fmt.Sprintf("%s:%s:%s", strings.Replace(serviceID, "::", "", -1), "bucket", bucketName)
	d.Set("crn", bucketCRN)
	d.Set("resource_instance_id", serviceID)
	d.Set("bucket_name", bucketName)
	d.Set("s3_endpoint_public", apiEndpoint)
	d.Set("s3_endpoint_private", apiEndpointPrivate)
	d.Set("s3_endpoint_direct", directApiEndpoint)
	if endpointType != "" {
		d.Set("endpoint_type", endpointType)
	}

	getBucketConfigOptions := &resourceconfigurationv1.GetBucketConfigOptions{
		Bucket: &bucketName,
	}

	sess, err := meta.(conns.ClientSession).CosConfigV1API()
	if err != nil {
		return err
	}
	if endpointType == "private" {
		sess.SetServiceURL("https://config.private.cloud-object-storage.cloud.ibm.com/v1")
	}
	if endpointType == "direct" {
		sess.SetServiceURL("https://config.direct.cloud-object-storage.cloud.ibm.com/v1")
	}

	if apiType == "sl" {

		satconfig := fmt.Sprintf("https://config.%s.%s.cloud-object-storage.appdomain.cloud/v1", serviceID, bLocation)

		sess.SetServiceURL(satconfig)
	}

	bucketPtr, response, err := sess.GetBucketConfig(getBucketConfigOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error in getting bucket info rule: %s\n%s", err, response)
	}
	head, err := s3Client.HeadBucket(headInput)
	if err != nil {
		return err
	}
	if head.IBMSSEKPEnabled != nil {
		if *head.IBMSSEKPEnabled == true {
			if keyProtectFlag == true {
				d.Set("key_protect", head.IBMSSEKPCrkId)
			} else {
				d.Set("kms_key_crn", head.IBMSSEKPCrkId)
			}
		}
	}

	if bucketPtr != nil {

		if bucketPtr.Firewall != nil {
			d.Set("allowed_ip", flex.FlattenStringList(bucketPtr.Firewall.AllowedIp))
		} else {

			d.Set("allowed_ip", []string{})
		}
		if bucketPtr.ActivityTracking != nil {
			d.Set("activity_tracking", flex.FlattenActivityTrack(bucketPtr.ActivityTracking))
		} else {

			d.Set("activity_tracking", []interface{}{})
		}

		if bucketPtr.MetricsMonitoring != nil {
			d.Set("metrics_monitoring", flex.FlattenMetricsMonitor(bucketPtr.MetricsMonitoring))
		} else {

			d.Set("metrics_monitoring", []interface{}{})
		}
		if bucketPtr.HardQuota != nil {
			d.Set("hard_quota", bucketPtr.HardQuota)
		} else {
			d.Set("hard_quota", 0)
		}
	}
	// Read the lifecycle configuration (archive & expiration or non current version or abort incomplete multipart upload)

	gInput := &s3.GetBucketLifecycleConfigurationInput{
		Bucket: aws.String(bucketName),
	}

	lifecycleptr, err := s3Client.GetBucketLifecycleConfiguration(gInput)

	if (err != nil && !strings.Contains(err.Error(), "NoSuchLifecycleConfiguration: The lifecycle configuration does not exist")) && (err != nil && bucketPtr != nil && bucketPtr.Firewall != nil && !strings.Contains(err.Error(), "AccessDenied: Access Denied")) {
		return err
	}
	if lifecycleptr != nil {
		archiveRules := flex.ArchiveRuleGet(lifecycleptr.Rules)
		expireRules := flex.ExpireRuleGet(lifecycleptr.Rules)
		nc_expRules := flex.Nc_exp_RuleGet(lifecycleptr.Rules)
		abort_mpuRules := flex.Abort_mpu_RuleGet(lifecycleptr.Rules)
		if len(archiveRules) > 0 {
			d.Set("archive_rule", archiveRules)
		}
		if len(expireRules) > 0 {
			d.Set("expire_rule", expireRules)
		}
		if len(nc_expRules) > 0 {
			d.Set("noncurrent_version_expiration", nc_expRules)
		}
		if len(abort_mpuRules) > 0 {
			d.Set("abort_incomplete_multipart_upload_days", abort_mpuRules)
		}
	}

	// Read retention rule
	retentionInput := &s3.GetBucketProtectionConfigurationInput{
		Bucket: aws.String(bucketName),
	}
	retentionptr, err := s3Client.GetBucketProtectionConfiguration(retentionInput)

	if err != nil && bucketPtr != nil && bucketPtr.Firewall != nil && !strings.Contains(err.Error(), "AccessDenied: Access Denied") {
		return err
	}

	if retentionptr != nil {
		retentionRules := flex.RetentionRuleGet(retentionptr.ProtectionConfiguration)
		if len(retentionRules) > 0 {
			d.Set("retention_rule", retentionRules)
		}
	}

	// Read Object versioning
	versionInput := &s3.GetBucketVersioningInput{
		Bucket: aws.String(bucketName),
	}

	versionPtr, err := s3Client.GetBucketVersioning(versionInput)

	if err != nil && bucketPtr != nil && bucketPtr.Firewall != nil && !strings.Contains(err.Error(), "AccessDenied: Access Denied") {
		return err
	}
	if versionPtr != nil {
		versioningData := flex.FlattenCosObejctVersioning(versionPtr)

		if len(versioningData) > 0 {
			d.Set("object_versioning", versioningData)
		} else {
			d.Set("object_versioning", nil)
		}
	}
	// reading objectlock
	getObjectLockConfigurationInput := &s3.GetObjectLockConfigurationInput{
		Bucket: aws.String(bucketName),
	}
	output, err := s3Client.GetObjectLockConfiguration(getObjectLockConfigurationInput)
	if output.ObjectLockConfiguration != nil {
		objectLockEnabled := *output.ObjectLockConfiguration.ObjectLockEnabled
		if objectLockEnabled == "Enabled" {
			d.Set("object_lock", true)
		}
	}
	return nil
}

func resourceIBMCOSBucketCreate(d *schema.ResourceData, meta interface{}) error {
	var s3Conf *aws.Config
	rsConClient, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return err
	}
	bucketName := d.Get("bucket_name").(string)
	storageClass := d.Get("storage_class").(string)
	objectLockEnabled := d.Get("object_lock").(bool)
	var bLocation string
	var apiType string
	var satlc_id string
	serviceID := d.Get("resource_instance_id").(string)

	//satlc_id := d.GetOK("satellite_location_id")
	if satlc, ok := d.GetOk("satellite_location_id"); ok {
		satlc_id = satlc.(string)
		satloc_guid := strings.Split(serviceID, ":")
		bucketsatcrn := satloc_guid[7]
		serviceID = bucketsatcrn
	}

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
	//Add satellite location id
	if bucketLocation, ok := d.GetOk("satellite_location_id"); ok {
		bLocation = bucketLocation.(string)
		apiType = "sl"
	}
	if bLocation == "" {
		return fmt.Errorf("Provide either `cross_region_location` or `region_location` or `single_site_location` or `satellite_location_id`")
	}

	lConstraint := fmt.Sprintf("%s-%s", bLocation, storageClass)

	var endpointType = d.Get("endpoint_type").(string)

	var apiEndpoint, privateApiEndpoint, directApiEndpoint string
	if apiType == "sl" {

		apiEndpoint = SelectSatlocCosApi(apiType, serviceID, bLocation)

	} else {
		apiEndpoint, privateApiEndpoint, directApiEndpoint = SelectCosApi(apiType, bLocation)
		if endpointType == "private" {
			apiEndpoint = privateApiEndpoint
		}
		if endpointType == "direct" {
			apiEndpoint = directApiEndpoint
		}

	}

	apiEndpoint = conns.EnvFallBack([]string{"IBMCLOUD_COS_ENDPOINT"}, apiEndpoint)

	if apiEndpoint == "" {
		return fmt.Errorf("[ERROR] The endpoint doesn't exists for given location %s and endpoint type %s", bLocation, endpointType)
	}

	var create *s3.CreateBucketInput
	if satlc_id != "" || storageClass == "" {
		create = &s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
		}
	} else if objectLockEnabled == true {
		create = &s3.CreateBucketInput{
			Bucket:                     aws.String(bucketName),
			ObjectLockEnabledForBucket: aws.Bool(true),
		}
	} else {
		create = &s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
			CreateBucketConfiguration: &s3.CreateBucketConfiguration{
				LocationConstraint: aws.String(lConstraint),
			},
		}
	}

	if keyprotect, ok := d.GetOk("key_protect"); ok {
		create.IBMSSEKPCustomerRootKeyCrn = aws.String(keyprotect.(string))
		create.IBMSSEKPEncryptionAlgorithm = aws.String(keyAlgorithm)
	} else if kmsKeyCrn, ok := d.GetOk("kms_key_crn"); ok {
		create.IBMSSEKPCustomerRootKeyCrn = aws.String(kmsKeyCrn.(string))
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
	rsConClient, _ := meta.(conns.ClientSession).BluemixSession()
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

	if bucketLocation, ok := d.GetOk("satellite_location_id"); ok {
		bLocation = bucketLocation.(string)
		apiType = "sl"
	}

	endpointType := parseBucketId(d.Id(), "endpointType")

	var apiEndpoint, apiEndpointPrivate, directApiEndpoint string

	if apiType == "sl" {

		apiEndpoint = SelectSatlocCosApi(apiType, serviceID, bLocation)

	} else {
		apiEndpoint, apiEndpointPrivate, directApiEndpoint = SelectCosApi(apiType, bLocation)
		if endpointType == "private" {
			apiEndpoint = apiEndpointPrivate
		}
		if endpointType == "direct" {
			apiEndpoint = directApiEndpoint
		}

	}

	apiEndpoint = conns.EnvFallBack([]string{"IBMCLOUD_COS_ENDPOINT"}, apiEndpoint)

	if apiEndpoint == "" {
		return fmt.Errorf("[ERROR] The endpoint doesn't exists for given location %s and endpoint type %s", bLocation, endpointType)
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

	delete := &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	}
	_, err = s3Client.DeleteBucket(delete)

	if err != nil && strings.Contains(err.Error(), "BucketNotEmpty") {
		if delbucket, ok := d.GetOk("force_delete"); ok {
			if delbucket.(bool) {
				// Use a S3 service client that can handle multiple slashes in URIs.
				// While ibm_cos_bucket_object resources cannot create these object
				// keys, other AWS services and applications using the COS Bucket can.

				// bucket may have things delete them
				log.Printf("[DEBUG] COS Bucket attempting to forceDelete %+v", err)

				// Delete everything including locked objects.
				// Don't ignore any object errors or we could recurse infinitely.
				err = deleteAllCOSObjectVersions(s3Client, bucketName, "", false, false)

				if err != nil {
					return fmt.Errorf("[ERROR] Error COS Bucket force_delete: %s", err)
				}

				// this line recurses until all objects are deleted or an error is returned
				return resourceIBMCOSBucketDelete(d, meta)
			}
		}
	}
	if err != nil {
		return fmt.Errorf("[ERROR] Error deleting COS Bucket (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceIBMCOSBucketExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	var s3Conf *aws.Config
	rsConClient, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return false, err
	}
	bucket_meta := strings.Split(d.Id(), ":meta:")
	if len(bucket_meta) < 2 || len(strings.Split(bucket_meta[1], ":")) < 2 {
		return false, fmt.Errorf("[ERROR] Error parsing bucket ID. Bucket ID format must be: $CRN:meta:$buckettype:$bucketlocation")
	}

	bucketName := parseBucketId(d.Id(), "bucketName")
	serviceID := parseBucketId(d.Id(), "serviceID")

	apiType := parseBucketId(d.Id(), "apiType")
	bLocation := parseBucketId(d.Id(), "bLocation")
	endpointType := parseBucketId(d.Id(), "endpointType")

	if apiType == "sl" {
		satloc_guid := strings.Split(serviceID, ":")
		bucketsatcrn := satloc_guid[0]
		serviceID = bucketsatcrn
	}

	var apiEndpoint, apiEndpointPrivate, directApiEndpoint string

	if apiType == "sl" {

		apiEndpoint = SelectSatlocCosApi(apiType, serviceID, bLocation)

	} else {
		apiEndpoint, apiEndpointPrivate, directApiEndpoint = SelectCosApi(apiType, bLocation)
		if endpointType == "private" {
			apiEndpoint = apiEndpointPrivate
		}
		if endpointType == "direct" {
			apiEndpoint = directApiEndpoint
		}

	}

	apiEndpoint = conns.EnvFallBack([]string{"IBMCLOUD_COS_ENDPOINT"}, apiEndpoint)

	if apiEndpoint == "" {
		return false, fmt.Errorf("[ERROR] The endpoint doesn't exists for given endpoint type %s", endpointType)
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

func SelectCosApi(apiType string, bLocation string) (string, string, string) {
	if apiType == "crl" {
		return fmt.Sprintf("s3.%s.cloud-object-storage.appdomain.cloud", bLocation), fmt.Sprintf("s3.private.%s.cloud-object-storage.appdomain.cloud", bLocation), fmt.Sprintf("s3.direct.%s.cloud-object-storage.appdomain.cloud", bLocation)
	}
	if apiType == "rl" {
		return fmt.Sprintf("s3.%s.cloud-object-storage.appdomain.cloud", bLocation), fmt.Sprintf("s3.private.%s.cloud-object-storage.appdomain.cloud", bLocation), fmt.Sprintf("s3.direct.%s.cloud-object-storage.appdomain.cloud", bLocation)
	}
	if apiType == "ssl" {
		return fmt.Sprintf("s3.%s.cloud-object-storage.appdomain.cloud", bLocation), fmt.Sprintf("s3.private.%s.cloud-object-storage.appdomain.cloud", bLocation), fmt.Sprintf("s3.direct.%s.cloud-object-storage.appdomain.cloud", bLocation)
	}
	return "", "", ""
}

// /Satellite ENdpoint configuration
func SelectSatlocCosApi(apiType string, serviceID string, bLocation string) string {
	if apiType == "sl" {
		return fmt.Sprintf("s3.%s.%s.cloud-object-storage.appdomain.cloud", serviceID, bLocation)
	}
	return ""
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

func resourceExpiryValidate(_ context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	if expire, ok := diff.GetOk("expire_rule"); ok {
		expire_list := expire.([]interface{})
		for _, l := range expire_list {
			expireMap, _ := l.(map[string]interface{})
			ctr := 0
			if val, days_exist := expireMap["days"]; days_exist && val.(int) != 0 {
				ctr++
			}
			if val, date_exist := expireMap["date"]; date_exist && val != "" {
				ctr++
			}
			if val, expired_object_delete_marker_exist := expireMap["expired_object_delete_marker"]; expired_object_delete_marker_exist && val.(bool) != false {
				ctr++
			}
			if ctr > 1 {
				return fmt.Errorf("[ERROR] The expiry 3 action elements (Days, Date, ExpiredObjectDeleteMarker) are all mutually exclusive. These can not be used with each other. Please set one expiry element on the same rule of expiration.")
			}
		}
	}
	return nil
}
