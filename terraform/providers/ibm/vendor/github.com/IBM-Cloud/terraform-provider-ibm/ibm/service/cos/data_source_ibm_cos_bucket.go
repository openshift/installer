// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cos

import (
	"fmt"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	rcsdk "github.com/IBM/ibm-cos-sdk-go-config/v2/resourceconfigurationv1"
	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/aws/credentials/ibmiam"
	token "github.com/IBM/ibm-cos-sdk-go/aws/credentials/ibmiam/token"
	"github.com/IBM/ibm-cos-sdk-go/aws/session"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
)

var bucketTypes = []string{"single_site_location", "region_location", "cross_region_location"}

var cosConfigUrls = map[string]string{
	"private": "https://config.private.cloud-object-storage.cloud.ibm.com/v1",
	"direct":  "https://config.direct.cloud-object-storage.cloud.ibm.com/v1",
}

func DataSourceIBMCosBucket() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMCosBucketRead,

		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bucket_type": {
				Type: schema.TypeString,
				// ValidateFunc:  validate.ValidateAllowedStringValues(bucketTypes),
				ValidateFunc:  validate.InvokeDataSourceValidator("ibm_cos_bucket", "bucket_type"),
				Optional:      true,
				RequiredWith:  []string{"bucket_region"},
				ConflictsWith: []string{"satellite_location_id"},
			},
			"bucket_region": {
				Type:          schema.TypeString,
				Optional:      true,
				RequiredWith:  []string{"bucket_type"},
				ConflictsWith: []string{"satellite_location_id"},
			},
			"resource_instance_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeDataSourceValidator("ibm_cos_bucket", "resource_instance_id"),
			},
			"satellite_location_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"bucket_type", "bucket_region"},
				ExactlyOneOf:  []string{"satellite_location_id", "bucket_region"},
			},
			"endpoint_type": {
				Type:     schema.TypeString,
				Optional: true,
				// ValidateFunc:  validate.ValidateAllowedStringValues([]string{"public", "private", "direct"}),
				ValidateFunc:  validate.InvokeDataSourceValidator("ibm_cos_bucket", "endpoint_type"),
				Description:   "COS endpoint type: public, private, direct",
				ConflictsWith: []string{"satellite_location_id"},
				Default:       "public",
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
			"kms_key_crn": {
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
			"s3_endpoint_direct": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Direct endpoint for the COS bucket",
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
							Description: "If set to true, all object read events (i.e. downloads) will be sent to Activity Tracker.",
						},
						"write_data_events": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "If set to true, all object write events (i.e. uploads) will be sent to Activity Tracker.",
						},
						"management_events": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "This field only applies if `activity_tracker_crn` is not populated. If set to true, all bucket management events will be sent to Activity Tracker.",
						},
						"activity_tracker_crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "When the `activity_tracker_crn` is not populated, then enabled events are sent to the Activity Tracker instance associated to the container's location unless otherwise specified in the Activity Tracker Event Routing service configuration.If `activity_tracker_crn` is populated, then enabled events are sent to the Activity Tracker instance specified and bucket management events are always enabled.",
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
							Description: "If set to `true`, all usage metrics (i.e. `bytes_used`) will be sent to the monitoring service.",
						},
						"request_metrics_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "If set to true, all request metrics (i.e. `rest.object.head`) will be sent to the monitoring service",
						},
						"metrics_monitoring_crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "When the `metrics_monitoring_crn` is not populated, then enabled metrics are sent to the monitoring instance associated to the container's location unless otherwise specified in the Metrics Router service configuration.If `metrics_monitoring_crn` is populated, then enabled events are sent to the Metrics Monitoring instance specified.",
						},
					},
				},
			},
			"abort_incomplete_multipart_upload_days": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique identifier for the rule. Rules allow you to set a specific time frame after which objects are deleted. Set Rule ID for cos bucket",
						},
						"enable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Enable or disable rule for a bucket",
						},
						"prefix": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The rule applies to any objects with keys that match this prefix",
						},
						"days_after_initiation": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the number of days when the specific rule action takes effect.",
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
						"date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the date when the specific rule action takes effect.",
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
						"expired_object_delete_marker": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Expired object delete markers can be automatically cleaned up to improve performance in bucket. This cannot be used alongside version expiration.",
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
			"noncurrent_version_expiration": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Enable configuration expire_rule to COS Bucket after a defined period of time",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique identifier for the rule.Expire rules allow you to set a specific time frame after which objects are deleted. Set Rule ID for cos bucket",
						},
						"enable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Enable or disable an expire rule for a bucket",
						},
						"prefix": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The rule applies to any objects with keys that match this prefix",
						},
						"noncurrent_days": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the number of days when the specific rule action takes effect.",
						},
					},
				},
			},
			"replication_rule": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Replicate objects between buckets, replicate across source and destination. A container for replication rules can add up to 1,000 rules. The maximum size of a replication configuration is 2 MB.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A unique identifier for the rule. The maximum value is 255 characters.",
						},
						"priority": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"enable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Enable or disable an replication rule for a bucket",
						},
						"prefix": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The rule applies to any objects with keys that match this prefix",
						},
						"deletemarker_replication_status": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether to replicate delete markers. It should be either Enable or Disable",
						},
						"destination_bucket_crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Cloud Resource Name (CRN) of the bucket where you want COS to store the results",
						},
					},
				},
			},
			"hard_quota": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "sets a maximum amount of storage (in bytes) available for a bucket",
			},
			"object_lock": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Description",
			},
			"object_lock_configuration": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Bucket level object lock settings includes Days, Years, Mode.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"object_lock_enabled": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Enable object lock on a COS bucket. This can be used to enable objectlock on an existing bucket",
						},
						"object_lock_rule": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"default_retention": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "An object lock configuration on the object at a bucket level, in the form of a days , years and mode that establishes a point in time after which the object can be deleted. This is applied at bucket level hence it is by default applied to all the object in the bucket unless a seperate retention period is set on the object.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Retention modes apply different levels of protection to the objects.",
												},
												"years": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Retention period in terms of years after which the object can be deleted.",
												},
												"days": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Retention period in terms of days after which the object can be deleted.",
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
			"website_configuration": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"error_document": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"index_document": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"suffix": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"redirect_all_requests_to": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"routing_rule": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Rules that define when a redirect is applied and the redirect behavior.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"condition": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "A condition that must be met for the specified redirect to be applie.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"http_error_code_returned_equals": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The HTTP error code when the redirect is applied. Valid codes are 4XX or 5XX..",
												},
												"key_prefix_equals": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The object key name prefix when the redirect is applied..",
												},
											},
										},
									},
									"redirect": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: ".",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"host_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The host name the request should be redirected to.",
												},
												"http_redirect_code": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The HTTP redirect code to use on the response. Valid codes are 3XX except 300..",
												},
												"protocol": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Protocol to be used in the Location header that is returned in the response.",
												},
												"replace_key_prefix_with": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The prefix of the object key name that replaces the value of KeyPrefixEquals in the redirect request.",
												},
												"replace_key_with": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The object key to be used in the Location header that is returned in the response.",
												},
											},
										},
									},
								},
							},
						},
						"routing_rules": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Rules that define when a redirect is applied and the redirect behavior.",
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
			"lifecycle_rule": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"abort_incomplete_multipart_upload": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"days_after_initiation": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"expiration": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"date": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"days": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"expired_object_delete_marker": {
										Type:     schema.TypeBool,
										Computed: true, // API returns false; conflicts with date and days
									},
								},
							},
						},
						"filter": {
							Type:     schema.TypeList,
							Computed: true,
							// IBM has filter a required parameter
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"prefix": {
										Type:     schema.TypeString, // check if prefix empty is eccepted and if filter empty is accepted
										Computed: true,
									},
									"object_size_greater_than": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"object_size_less_than": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"tag": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"value": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"and": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"object_size_greater_than": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"object_size_less_than": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"prefix": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"tags": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"key": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"value": {
																Type:     schema.TypeString,
																Computed: true,
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
							Type:     schema.TypeString,
							Computed: true,
						},
						"noncurrent_version_expiration": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"noncurrent_days": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transition": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"date": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"days": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"storag_class": {
										Type:     schema.TypeString,
										Computed: true,
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

func DataSourceIBMCosBucketValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "resource_instance_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			CloudDataType:              "resource_instance",
			CloudDataRange:             []string{"service:cloud-object-storage"}})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "bucket_type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "single_site_location,region_location,cross_region_location",
		})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "endpoint_type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "public,private,direct",
		})

	ibmCOSBucketDataSourceValidator := validate.ResourceValidator{ResourceName: "ibm_cos_bucket", Schema: validateSchema}
	return &ibmCOSBucketDataSourceValidator
}
func dataSourceIBMCosBucketRead(d *schema.ResourceData, meta interface{}) error {
	var s3Conf *aws.Config
	var keyProtectFlag bool
	rsConClient, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return err
	}
	bucketName := d.Get("bucket_name").(string)
	serviceID := d.Get("resource_instance_id").(string)
	bucketType := d.Get("bucket_type").(string)
	bucketRegion := d.Get("bucket_region").(string)
	endpointType := d.Get("endpoint_type").(string)
	if _, ok := d.GetOk("key_protect"); ok {
		keyProtectFlag = true
	}

	var satlc_id, apiEndpoint, apiEndpointPublic, apiEndpointPrivate, directApiEndpoint, visibility string

	if satlc, ok := d.GetOk("satellite_location_id"); ok {
		satlc_id = satlc.(string)
		satloc_guid := strings.Split(serviceID, ":")
		bucketsatcrn := satloc_guid[7]
		serviceID = bucketsatcrn
		bucketType = "sl"
	}

	if bucketType == "sl" {
		apiEndpoint = SelectSatlocCosApi(bucketType, serviceID, satlc_id)

	} else {
		apiEndpoint, apiEndpointPrivate, directApiEndpoint = SelectCosApi(bucketLocationConvert(bucketType), bucketRegion)
		visibility = endpointType
		if endpointType == "private" {
			apiEndpoint = apiEndpointPrivate
		}
		if endpointType == "direct" {
			// visibility type "direct" is not supported in endpoints file.
			visibility = "private"
			apiEndpoint = directApiEndpoint
		}

	}

	apiEndpoint = conns.FileFallBack(rsConClient.Config.EndpointsFile, visibility, "IBMCLOUD_COS_ENDPOINT", bucketRegion, apiEndpoint)
	apiEndpoint = conns.EnvFallBack([]string{"IBMCLOUD_COS_ENDPOINT"}, apiEndpoint)
	if apiEndpoint == "" {
		return fmt.Errorf("[ERROR] The endpoint doesn't exists for given location %s and endpoint type %s", bucketRegion, endpointType)
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

	if bucketType != "sl" {
		bucketLocationInput := &s3.GetBucketLocationInput{
			Bucket: aws.String(bucketName),
		}
		bucketLocationConstraint, err := s3Client.GetBucketLocation(bucketLocationInput)
		if err != nil {
			return err
		}
		bLocationConstraint := *bucketLocationConstraint.LocationConstraint

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
		d.Set("satellite_location_id", satlc_id)
	}

	head, err := s3Client.HeadBucket(headInput)
	if err != nil {
		return err
	}
	bucketID := fmt.Sprintf("%s:%s:%s:meta:%s:%s:%s", strings.Replace(serviceID, "::", "", -1), "bucket", bucketName, bucketLocationConvert(bucketType), bucketRegion, endpointType)
	d.SetId(bucketID)
	if head.IBMSSEKPEnabled != nil {
		if *head.IBMSSEKPEnabled == true {
			if keyProtectFlag == true {
				d.Set("key_protect", head.IBMSSEKPCrkId)
			} else {
				d.Set("kms_key_crn", head.IBMSSEKPCrkId)
			}
		}
	}

	bucketCRN := fmt.Sprintf("%s:%s:%s", strings.Replace(serviceID, "::", "", -1), "bucket", bucketName)
	d.Set("crn", bucketCRN)
	d.Set("resource_instance_id", serviceID)

	testEnv := strings.Contains(apiEndpoint, ".test.")
	apiEndpointPublic, apiEndpointPrivate, directApiEndpoint = SelectCosApi(bucketLocationConvert(bucketType), bucketRegion)
	if testEnv {
		d.Set(fmt.Sprintf("s3_endpoint_%s", endpointType), apiEndpoint)
	} else {
		d.Set("s3_endpoint_public", apiEndpointPublic)
		d.Set("s3_endpoint_private", apiEndpointPrivate)
		d.Set("s3_endpoint_direct", directApiEndpoint)
	}
	sess, err := meta.(conns.ClientSession).CosConfigV1API()
	if err != nil {
		return err
	}
	if endpointType != "public" {
		// User is expected to define both private and direct url type under "private" in endpoints file since visibility type "direct" is not supported.
		cosConfigURL := conns.FileFallBack(rsConClient.Config.EndpointsFile, "private", "IBMCLOUD_COS_CONFIG_ENDPOINT", bucketRegion, cosConfigUrls[endpointType])
		cosConfigURL = conns.EnvFallBack([]string{"IBMCLOUD_COS_CONFIG_ENDPOINT"}, cosConfigURL)
		if cosConfigURL != "" {
			sess.SetServiceURL(cosConfigURL)
		}
	}

	if bucketType == "sl" {
		satconfig := fmt.Sprintf("https://config.%s.%s.cloud-object-storage.appdomain.cloud/v1", serviceID, satlc_id)

		sess.SetServiceURL(satconfig)
	}

	getOptions := new(rcsdk.GetBucketConfigOptions)
	getOptions.SetBucket(bucketName)
	bucketPtr, response, err := sess.GetBucketConfig(getOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error in getting bucket info rule: %s\n%s", err, response)
	}

	if bucketPtr != nil {
		if bucketPtr.Firewall != nil {
			d.Set("allowed_ip", flex.FlattenStringList(bucketPtr.Firewall.AllowedIp))
		}
		if bucketPtr.ActivityTracking != nil {
			d.Set("activity_tracking", flex.FlattenActivityTrack(bucketPtr.ActivityTracking))
		}
		if bucketPtr.MetricsMonitoring != nil {
			d.Set("metrics_monitoring", flex.FlattenMetricsMonitor(bucketPtr.MetricsMonitoring))
		}
		if bucketPtr.HardQuota != nil {
			d.Set("hard_quota", bucketPtr.HardQuota)
		}

	}

	// Read the lifecycle configuration

	gInput := &s3.GetBucketLifecycleConfigurationInput{
		Bucket: aws.String(bucketName),
	}

	lifecycleptr, err := s3Client.GetBucketLifecycleConfiguration(gInput)

	if (err != nil && !strings.Contains(err.Error(), "NoSuchLifecycleConfiguration: The lifecycle configuration does not exist")) && (err != nil && bucketPtr != nil && bucketPtr.Firewall != nil && !strings.Contains(err.Error(), "AccessDenied: Access Denied")) {
		return err
	}

	if lifecycleptr != nil {
		if len(lifecycleptr.Rules) > 0 {
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
	}

	//lifecycle configuration new resource read
	getLifecycleConfigurationInput := &s3.GetBucketLifecycleConfigurationInput{
		Bucket: aws.String(bucketName),
	}
	var outputLifecycleConfig *s3.GetBucketLifecycleConfigurationOutput
	outputLifecycleConfig, err = s3Client.GetBucketLifecycleConfiguration(getLifecycleConfigurationInput)
	var outputLifecycleConfigptr *s3.LifecycleConfiguration
	outputLifecycleConfigptr = (*s3.LifecycleConfiguration)(outputLifecycleConfig)
	if (err != nil && !strings.Contains(err.Error(), "NoSuchLifecycleConfiguration: The lifecycle configuration does not exist")) && (err != nil && bucketPtr != nil && bucketPtr.Firewall != nil && !strings.Contains(err.Error(), "AccessDenied: Access Denied")) {
		return err
	}
	if outputLifecycleConfigptr.Rules != nil {
		lifecycleConfiguration := flex.LifecylceRuleGet(outputLifecycleConfigptr.Rules)
		if len(lifecycleConfiguration) > 0 {
			d.Set("lifecycle_rule", lifecycleConfiguration)
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
		retentionRules := flex.RetentionRuleGet(retentionptr.ProtectionConfiguration)
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
		versioningData := flex.FlattenCosObejctVersioning(versionPtr)
		if len(versioningData) > 0 {
			d.Set("object_versioning", versioningData)
		}
	}

	// Get the replication rules
	getBucketReplicationInput := &s3.GetBucketReplicationInput{
		Bucket: aws.String(bucketName),
	}

	replicationptr, err := s3Client.GetBucketReplication(getBucketReplicationInput)

	if err != nil && !strings.Contains(err.Error(), "AccessDenied: Access Denied") && !strings.Contains(err.Error(), "The replication configuration was not found") {
		return err
	}

	if replicationptr != nil {
		replicationRules := flex.ReplicationRuleGet(replicationptr.ReplicationConfiguration)
		if len(replicationRules) > 0 {
			d.Set("replication_rule", replicationRules)
		}
	}
	// reading objectlock for bucket
	getObjectLockConfigurationInput := &s3.GetObjectLockConfigurationInput{
		Bucket: aws.String(bucketName),
	}
	output, err := s3Client.GetObjectLockConfiguration(getObjectLockConfigurationInput)
	if output.ObjectLockConfiguration != nil {
		objectLockEnabled := *output.ObjectLockConfiguration.ObjectLockEnabled
		if objectLockEnabled == "Enabled" {
			d.Set("object_lock", true)
		}
		objectLockConfigurationptr := output.ObjectLockConfiguration
		objectLockConfiguration := flex.ObjectLockConfigurationGet(objectLockConfigurationptr)
		if len(objectLockConfiguration) > 0 {
			d.Set("object_lock_configuration", objectLockConfiguration)
		}
	} //getBucketConfiguration
	getBucketWebsiteConfigurationInput := &s3.GetBucketWebsiteInput{
		Bucket: aws.String(bucketName),
	}
	outputwebsite, err := s3Client.GetBucketWebsite(getBucketWebsiteConfigurationInput)
	var outputptr *s3.WebsiteConfiguration
	outputptr = (*s3.WebsiteConfiguration)(outputwebsite)
	if err != nil && !strings.Contains(err.Error(), "AccessDenied: Access Denied") && !strings.Contains(err.Error(), "The specified bucket does not have a website configuration") {
		return err
	}
	if outputwebsite.IndexDocument != nil || outputwebsite.RedirectAllRequestsTo != nil {
		websiteConfiguration := flex.WebsiteConfigurationGet(outputptr)
		if len(websiteConfiguration) > 0 {
			d.Set("website_configuration", websiteConfiguration)
		}
		websiteEndpoint := getWebsiteEndpoint(bucketName, bucketRegion)
		if websiteEndpoint != "" {
			d.Set("website_endpoint", websiteEndpoint)
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
		return "ssl"
	}
	return ""
}
