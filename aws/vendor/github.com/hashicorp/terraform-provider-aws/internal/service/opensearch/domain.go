package opensearch

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/opensearchservice"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	awspolicy "github.com/hashicorp/awspolicyequivalence"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_opensearch_domain", name="Domain")
// @Tags(identifierAttribute="id")
func ResourceDomain() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceDomainCreate,
		ReadWithoutTimeout:   resourceDomainRead,
		UpdateWithoutTimeout: resourceDomainUpdate,
		DeleteWithoutTimeout: resourceDomainDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDomainImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(180 * time.Minute),
			Delete: schema.DefaultTimeout(90 * time.Minute),
		},

		CustomizeDiff: customdiff.Sequence(
			customdiff.ForceNewIf("engine_version", func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool {
				newVersion := d.Get("engine_version").(string)
				domainName := d.Get("domain_name").(string)

				conn := meta.(*conns.AWSClient).OpenSearchConn(ctx)
				resp, err := conn.GetCompatibleVersionsWithContext(ctx, &opensearchservice.GetCompatibleVersionsInput{
					DomainName: aws.String(domainName),
				})
				if err != nil {
					log.Printf("[ERROR] Failed to get compatible OpenSearch versions %s", domainName)
					return false
				}
				if len(resp.CompatibleVersions) != 1 {
					return true
				}
				for _, targetVersion := range resp.CompatibleVersions[0].TargetVersions {
					if aws.StringValue(targetVersion) == newVersion {
						return false
					}
				}
				return true
			}),
			customdiff.ForceNewIf("encrypt_at_rest.0.enabled", func(_ context.Context, d *schema.ResourceDiff, meta interface{}) bool {
				o, n := d.GetChange("encrypt_at_rest.0.enabled")
				if o.(bool) && !n.(bool) {
					return true
				}

				return !inPlaceEncryptionEnableVersion(d.Get("engine_version").(string))
			}),
			customdiff.ForceNewIf("node_to_node_encryption.0.enabled", func(_ context.Context, d *schema.ResourceDiff, meta interface{}) bool {
				o, n := d.GetChange("node_to_node_encryption.0.enabled")
				if o.(bool) && !n.(bool) {
					return true
				}

				return !inPlaceEncryptionEnableVersion(d.Get("engine_version").(string))
			}),
			customdiff.ForceNewIf("advanced_security_options.0.enabled", func(_ context.Context, d *schema.ResourceDiff, meta interface{}) bool {
				o, n := d.GetChange("advanced_security_options.0.enabled")
				if o.(bool) && !n.(bool) {
					return true
				}

				return false
			}),
			verify.SetTagsDiff,
		),

		Schema: map[string]*schema.Schema{
			"access_policies": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: verify.SuppressEquivalentPolicyDiffs,
				StateFunc: func(v interface{}) string {
					json, _ := structure.NormalizeJsonString(v)
					return json
				},
			},
			"advanced_options": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"advanced_security_options": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"anonymous_auth_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"internal_user_database_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"master_user_options": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"master_user_arn": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: verify.ValidARN,
									},
									"master_user_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"master_user_password": {
										Type:      schema.TypeString,
										Optional:  true,
										Sensitive: true,
									},
								},
							},
						},
					},
				},
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_tune_options": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"desired_state": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(opensearchservice.AutoTuneDesiredState_Values(), false),
						},
						"maintenance_schedule": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cron_expression_for_recurrence": {
										Type:     schema.TypeString,
										Required: true,
									},
									"duration": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"unit": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validation.StringInSlice(opensearchservice.TimeUnit_Values(), false),
												},
												"value": {
													Type:     schema.TypeInt,
													Required: true,
												},
											},
										},
									},
									"start_at": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.IsRFC3339Time,
									},
								},
							},
						},
						"rollback_on_disable": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice(opensearchservice.RollbackOnDisable_Values(), false),
						},
					},
				},
			},
			"cluster_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cold_storage_options": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"dedicated_master_count": {
							Type:             schema.TypeInt,
							Optional:         true,
							DiffSuppressFunc: isDedicatedMasterDisabled,
						},
						"dedicated_master_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"dedicated_master_type": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: isDedicatedMasterDisabled,
						},
						"instance_count": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  opensearchservice.OpenSearchPartitionInstanceTypeM3MediumSearch,
						},
						"warm_count": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(2, 150),
						},
						"warm_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"warm_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(opensearchservice.OpenSearchWarmPartitionInstanceType_Values(), false),
						},
						"zone_awareness_config": {
							Type:             schema.TypeList,
							Optional:         true,
							MaxItems:         1,
							DiffSuppressFunc: verify.SuppressMissingOptionalConfigurationBlock,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"availability_zone_count": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      2,
										ValidateFunc: validation.IntInSlice([]int{2, 3}),
									},
								},
							},
						},
						"zone_awareness_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"cognito_options": {
				Type:             schema.TypeList,
				Optional:         true,
				MaxItems:         1,
				DiffSuppressFunc: verify.SuppressMissingOptionalConfigurationBlock,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"identity_pool_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"role_arn": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: verify.ValidARN,
						},
						"user_pool_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"dashboard_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_endpoint_options": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"custom_endpoint": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: isCustomEndpointDisabled,
						},
						"custom_endpoint_certificate_arn": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateFunc:     verify.ValidARN,
							DiffSuppressFunc: isCustomEndpointDisabled,
						},
						"custom_endpoint_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"enforce_https": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"tls_security_policy": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice(opensearchservice.TLSSecurityPolicy_Values(), false),
						},
					},
				},
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-z][0-9a-z\-]{2,27}$`),
					"must start with a lowercase alphabet and be at least 3 and no more than 28 characters long."+
						" Valid characters are a-z (lowercase letters), 0-9, and - (hyphen)."),
			},
			"ebs_options": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ebs_enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"iops": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"throughput": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntAtLeast(125),
						},
						"volume_size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"volume_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice(opensearchservice.VolumeType_Values(), false),
						},
					},
				},
			},
			"encrypt_at_rest": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"kms_key_id": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							ForceNew:         true,
							DiffSuppressFunc: suppressEquivalentKMSKeyIDs,
						},
					},
				},
			},
			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"kibana_endpoint": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "use 'dashboard_endpoint' attribute instead",
			},
			"log_publishing_options": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloudwatch_log_group_arn": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: verify.ValidARN,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"log_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(opensearchservice.LogType_Values(), false),
						},
					},
				},
			},
			"node_to_node_encryption": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},
			"off_peak_window_options": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"off_peak_window": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"window_start_time": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"hours": {
													Type:     schema.TypeInt,
													Optional: true,
													Computed: true,
												},
												"minutes": {
													Type:     schema.TypeInt,
													Optional: true,
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
			"snapshot_options": {
				Type:             schema.TypeList,
				Optional:         true,
				MaxItems:         1,
				DiffSuppressFunc: verify.SuppressMissingOptionalConfigurationBlock,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"automated_snapshot_start_hour": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
			"vpc_options": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zones": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						"security_group_ids": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						"subnet_ids": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceDomainImport(ctx context.Context,
	d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	d.Set("domain_name", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceDomainCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).OpenSearchConn(ctx)

	// The API doesn't check for duplicate names
	// so w/out this check Create would act as upsert
	// and might cause duplicate domain to appear in state
	resp, err := FindDomainByName(ctx, conn, d.Get("domain_name").(string))
	if err == nil {
		return sdkdiag.AppendErrorf(diags, "OpenSearch Domain %q already exists", aws.StringValue(resp.DomainName))
	}

	input := &opensearchservice.CreateDomainInput{
		DomainName: aws.String(d.Get("domain_name").(string)),
		TagList:    GetTagsIn(ctx),
	}

	if v, ok := d.GetOk("engine_version"); ok {
		input.EngineVersion = aws.String(v.(string))
	}

	if v, ok := d.GetOk("access_policies"); ok {
		policy, err := structure.NormalizeJsonString(v.(string))

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "policy (%s) is invalid JSON: %s", policy, err)
		}

		input.AccessPolicies = aws.String(policy)
	}

	if v, ok := d.GetOk("advanced_options"); ok {
		input.AdvancedOptions = flex.ExpandStringMap(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("advanced_security_options"); ok {
		input.AdvancedSecurityOptions = expandAdvancedSecurityOptions(v.([]interface{}))
	}

	if v, ok := d.GetOk("auto_tune_options"); ok && len(v.([]interface{})) > 0 {
		input.AutoTuneOptions = expandAutoTuneOptionsInput(v.([]interface{})[0].(map[string]interface{}))
	}

	if v, ok := d.GetOk("ebs_options"); ok {
		options := v.([]interface{})

		if len(options) == 1 {
			if options[0] == nil {
				return sdkdiag.AppendErrorf(diags, "At least one field is expected inside ebs_options")
			}

			s := options[0].(map[string]interface{})
			input.EBSOptions = expandEBSOptions(s)
		}
	}

	if v, ok := d.GetOk("encrypt_at_rest"); ok {
		options := v.([]interface{})
		if options[0] == nil {
			return sdkdiag.AppendErrorf(diags, "At least one field is expected inside encrypt_at_rest")
		}

		s := options[0].(map[string]interface{})
		input.EncryptionAtRestOptions = expandEncryptAtRestOptions(s)
	}

	if v, ok := d.GetOk("cluster_config"); ok {
		config := v.([]interface{})

		if len(config) == 1 {
			if config[0] == nil {
				return sdkdiag.AppendErrorf(diags, "At least one field is expected inside cluster_config")
			}
			m := config[0].(map[string]interface{})
			input.ClusterConfig = expandClusterConfig(m)
		}
	}

	if v, ok := d.GetOk("node_to_node_encryption"); ok {
		options := v.([]interface{})

		s := options[0].(map[string]interface{})
		input.NodeToNodeEncryptionOptions = expandNodeToNodeEncryptionOptions(s)
	}

	if v, ok := d.GetOk("snapshot_options"); ok {
		options := v.([]interface{})

		if len(options) == 1 {
			if options[0] == nil {
				return sdkdiag.AppendErrorf(diags, "At least one field is expected inside snapshot_options")
			}

			o := options[0].(map[string]interface{})

			snapshotOptions := opensearchservice.SnapshotOptions{
				AutomatedSnapshotStartHour: aws.Int64(int64(o["automated_snapshot_start_hour"].(int))),
			}

			input.SnapshotOptions = &snapshotOptions
		}
	}

	if v, ok := d.GetOk("vpc_options"); ok {
		options := v.([]interface{})
		if options[0] == nil {
			return sdkdiag.AppendErrorf(diags, "At least one field is expected inside vpc_options")
		}

		s := options[0].(map[string]interface{})
		input.VPCOptions = expandVPCOptions(s)
	}

	if v, ok := d.GetOk("log_publishing_options"); ok {
		input.LogPublishingOptions = expandLogPublishingOptions(v.(*schema.Set))
	}

	if v, ok := d.GetOk("domain_endpoint_options"); ok {
		input.DomainEndpointOptions = expandDomainEndpointOptions(v.([]interface{}))
	}

	if v, ok := d.GetOk("cognito_options"); ok {
		input.CognitoOptions = expandCognitoOptions(v.([]interface{}))
	}

	if v, ok := d.GetOk("off_peak_window_options"); ok && len(v.([]interface{})) > 0 {
		input.OffPeakWindowOptions = expandOffPeakWindowOptions(v.([]interface{})[0].(map[string]interface{}))

		// This option is only available when modifying a domain created prior to February 16, 2023, not when creating a new domain.
		// An off-peak window is required for a domain and cannot be disabled.
		if input.OffPeakWindowOptions != nil {
			input.OffPeakWindowOptions.Enabled = aws.Bool(true)
		}
	}

	// IAM Roles can take some time to propagate if set in AccessPolicies and created in the same terraform
	var out *opensearchservice.CreateDomainOutput
	err = retry.RetryContext(ctx, propagationTimeout, func() *retry.RetryError {
		var err error
		out, err = conn.CreateDomainWithContext(ctx, input)
		if err != nil {
			if tfawserr.ErrMessageContains(err, "InvalidTypeException", "Error setting policy") {
				return retry.RetryableError(err)
			}
			if tfawserr.ErrMessageContains(err, "ValidationException", "enable a service-linked role to give Amazon ES permissions") {
				return retry.RetryableError(err)
			}
			if tfawserr.ErrMessageContains(err, "ValidationException", "Domain is still being deleted") {
				return retry.RetryableError(err)
			}
			if tfawserr.ErrMessageContains(err, "ValidationException", "Amazon OpenSearch Service must be allowed to use the passed role") {
				return retry.RetryableError(err)
			}
			if tfawserr.ErrMessageContains(err, "ValidationException", "The passed role has not propagated yet") {
				return retry.RetryableError(err)
			}
			if tfawserr.ErrMessageContains(err, "ValidationException", "Authentication error") {
				return retry.RetryableError(err)
			}
			if tfawserr.ErrMessageContains(err, "ValidationException", "Unauthorized Operation: OpenSearch Service must be authorised to describe") {
				return retry.RetryableError(err)
			}
			if tfawserr.ErrMessageContains(err, "ValidationException", "The passed role must authorize Amazon OpenSearch Service to describe") {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if tfresource.TimedOut(err) {
		out, err = conn.CreateDomainWithContext(ctx, input)
	}
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating OpenSearch Domain: %s", err)
	}

	d.SetId(aws.StringValue(out.DomainStatus.ARN))

	log.Printf("[DEBUG] Waiting for OpenSearch Domain %q to be created", d.Id())
	if err := WaitForDomainCreation(ctx, conn, d.Get("domain_name").(string), d.Timeout(schema.TimeoutCreate)); err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for OpenSearch Domain (%s) to be created: %s", d.Id(), err)
	}

	log.Printf("[DEBUG] OpenSearch Domain %q created", d.Id())

	if v, ok := d.GetOk("auto_tune_options"); ok && len(v.([]interface{})) > 0 {
		log.Printf("[DEBUG] Modifying config for OpenSearch Domain %q", d.Id())

		input := &opensearchservice.UpdateDomainConfigInput{
			DomainName: aws.String(d.Get("domain_name").(string)),
		}

		input.AutoTuneOptions = expandAutoTuneOptions(v.([]interface{})[0].(map[string]interface{}))

		_, err = conn.UpdateDomainConfigWithContext(ctx, input)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "modifying config for OpenSearch Domain: %s", err)
		}

		log.Printf("[DEBUG] Config for OpenSearch Domain %q modified", d.Id())
	}

	return append(diags, resourceDomainRead(ctx, d, meta)...)
}

func resourceDomainRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).OpenSearchConn(ctx)

	ds, err := FindDomainByName(ctx, conn, d.Get("domain_name").(string))

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] OpenSearch Domain (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading OpenSearch Domain (%s): %s", d.Id(), err)
	}

	outDescribeDomainConfig, err := conn.DescribeDomainConfigWithContext(ctx, &opensearchservice.DescribeDomainConfigInput{
		DomainName: aws.String(d.Get("domain_name").(string)),
	})

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading OpenSearch Domain (%s): %s", d.Id(), err)
	}

	dc := outDescribeDomainConfig.DomainConfig

	if ds.AccessPolicies != nil && aws.StringValue(ds.AccessPolicies) != "" {
		policies, err := verify.PolicyToSet(d.Get("access_policies").(string), aws.StringValue(ds.AccessPolicies))

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "reading OpenSearch Domain (%s): %s", d.Id(), err)
		}

		d.Set("access_policies", policies)
	}

	options := advancedOptionsIgnoreDefault(d.Get("advanced_options").(map[string]interface{}), flex.PointersMapToStringList(ds.AdvancedOptions))
	if err = d.Set("advanced_options", options); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting advanced_options %v: %s", options, err)
	}

	d.SetId(aws.StringValue(ds.ARN))
	d.Set("arn", ds.ARN)
	d.Set("domain_id", ds.DomainId)
	d.Set("domain_name", ds.DomainName)
	d.Set("engine_version", ds.EngineVersion)

	if err := d.Set("ebs_options", flattenEBSOptions(ds.EBSOptions)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting ebs_options: %s", err)
	}

	if err := d.Set("encrypt_at_rest", flattenEncryptAtRestOptions(ds.EncryptionAtRestOptions)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting encrypt_at_rest: %s", err)
	}

	if err := d.Set("cluster_config", flattenClusterConfig(ds.ClusterConfig)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting cluster_config: %s", err)
	}

	if err := d.Set("cognito_options", flattenCognitoOptions(ds.CognitoOptions)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting cognito_options: %s", err)
	}

	if err := d.Set("node_to_node_encryption", flattenNodeToNodeEncryptionOptions(ds.NodeToNodeEncryptionOptions)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting node_to_node_encryption: %s", err)
	}

	// Populate AdvancedSecurityOptions with values returned from
	// DescribeDomainConfig, if enabled, else use
	// values from resource; additionally, append MasterUserOptions
	// from resource as they are not returned from the API
	if ds.AdvancedSecurityOptions != nil && aws.BoolValue(ds.AdvancedSecurityOptions.Enabled) {
		advSecOpts := flattenAdvancedSecurityOptions(ds.AdvancedSecurityOptions)
		advSecOpts[0]["master_user_options"] = getMasterUserOptions(d)
		if err := d.Set("advanced_security_options", advSecOpts); err != nil {
			return sdkdiag.AppendErrorf(diags, "setting advanced_security_options: %s", err)
		}
	}

	if v := dc.AutoTuneOptions; v != nil {
		err = d.Set("auto_tune_options", []interface{}{flattenAutoTuneOptions(v.Options)})
		if err != nil {
			return sdkdiag.AppendErrorf(diags, "reading OpenSearch Domain (%s): setting auto_tune_options: %s", d.Id(), err)
		}
	}

	if err := d.Set("snapshot_options", flattenSnapshotOptions(ds.SnapshotOptions)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting snapshot_options: %s", err)
	}

	if ds.VPCOptions != nil {
		if err := d.Set("vpc_options", flattenVPCDerivedInfo(ds.VPCOptions)); err != nil {
			return sdkdiag.AppendErrorf(diags, "setting vpc_options: %s", err)
		}

		endpoints := flex.PointersMapToStringList(ds.Endpoints)
		d.Set("endpoint", endpoints["vpc"])
		d.Set("dashboard_endpoint", getDashboardEndpoint(d))
		d.Set("kibana_endpoint", getKibanaEndpoint(d))
		if ds.Endpoint != nil {
			return sdkdiag.AppendErrorf(diags, "%q: OpenSearch Domain in VPC expected to have null Endpoint value", d.Id())
		}
	} else {
		if ds.Endpoint != nil {
			d.Set("endpoint", ds.Endpoint)
			d.Set("dashboard_endpoint", getDashboardEndpoint(d))
			d.Set("kibana_endpoint", getKibanaEndpoint(d))
		}
		if ds.Endpoints != nil {
			return sdkdiag.AppendErrorf(diags, "%q: OpenSearch Domain not in VPC expected to have null Endpoints value", d.Id())
		}
	}

	if err := d.Set("log_publishing_options", flattenLogPublishingOptions(ds.LogPublishingOptions)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting log_publishing_options: %s", err)
	}

	if err := d.Set("domain_endpoint_options", flattenDomainEndpointOptions(ds.DomainEndpointOptions)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting domain_endpoint_options: %s", err)
	}

	if ds.OffPeakWindowOptions != nil {
		if err := d.Set("off_peak_window_options", []interface{}{flattenOffPeakWindowOptions(ds.OffPeakWindowOptions)}); err != nil {
			return sdkdiag.AppendErrorf(diags, "setting off_peak_window_options: %s", err)
		}
	} else {
		d.Set("off_peak_window_options", nil)
	}

	return diags
}

func resourceDomainUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).OpenSearchConn(ctx)

	if d.HasChangesExcept("tags", "tags_all") {
		input := opensearchservice.UpdateDomainConfigInput{
			DomainName: aws.String(d.Get("domain_name").(string)),
		}

		if d.HasChange("access_policies") {
			o, n := d.GetChange("access_policies")

			if equivalent, err := awspolicy.PoliciesAreEquivalent(o.(string), n.(string)); err != nil || !equivalent {
				input.AccessPolicies = aws.String(d.Get("access_policies").(string))
			}
		}

		if d.HasChange("advanced_options") {
			input.AdvancedOptions = flex.ExpandStringMap(d.Get("advanced_options").(map[string]interface{}))
		}

		if d.HasChange("advanced_security_options") {
			input.AdvancedSecurityOptions = expandAdvancedSecurityOptions(d.Get("advanced_security_options").([]interface{}))
		}

		if d.HasChange("auto_tune_options") {
			input.AutoTuneOptions = expandAutoTuneOptions(d.Get("auto_tune_options").([]interface{})[0].(map[string]interface{}))
		}

		if d.HasChange("cognito_options") {
			input.CognitoOptions = expandCognitoOptions(d.Get("cognito_options").([]interface{}))
		}

		if d.HasChange("domain_endpoint_options") {
			input.DomainEndpointOptions = expandDomainEndpointOptions(d.Get("domain_endpoint_options").([]interface{}))
		}

		if d.HasChanges("ebs_options", "cluster_config") {
			options := d.Get("ebs_options").([]interface{})

			if len(options) == 1 {
				s := options[0].(map[string]interface{})
				input.EBSOptions = expandEBSOptions(s)
			}

			if d.HasChange("cluster_config") {
				config := d.Get("cluster_config").([]interface{})

				if len(config) == 1 {
					m := config[0].(map[string]interface{})
					input.ClusterConfig = expandClusterConfig(m)

					// Work around "ValidationException: Your domain's Elasticsearch version does not support cold storage options. Upgrade to Elasticsearch 7.9 or later.".
					if engineType, version, err := ParseEngineVersion(d.Get("engine_version").(string)); err == nil {
						switch engineType {
						case opensearchservice.EngineTypeElasticsearch:
							if verify.SemVerLessThan(version, "7.9") {
								input.ClusterConfig.ColdStorageOptions = nil
							}
						case opensearchservice.EngineTypeOpenSearch:
							// All OpenSearch versions support cold storage options.
						default:
							log.Printf("[WARN] unknown engine type: %s", engineType)
						}
					} else {
						log.Printf("[WARN] %s", err)
					}
				}
			}
		}

		if d.HasChange("encrypt_at_rest") {
			input.EncryptionAtRestOptions = nil
			if v, ok := d.GetOk("encrypt_at_rest"); ok {
				options := v.([]interface{})
				if options[0] == nil {
					return sdkdiag.AppendErrorf(diags, "at least one field is expected inside encrypt_at_rest")
				}

				s := options[0].(map[string]interface{})
				input.EncryptionAtRestOptions = expandEncryptAtRestOptions(s)
			}
		}

		if d.HasChange("log_publishing_options") {
			input.LogPublishingOptions = expandLogPublishingOptions(d.Get("log_publishing_options").(*schema.Set))
		}

		if d.HasChange("node_to_node_encryption") {
			input.NodeToNodeEncryptionOptions = nil
			if v, ok := d.GetOk("node_to_node_encryption"); ok {
				options := v.([]interface{})

				s := options[0].(map[string]interface{})
				input.NodeToNodeEncryptionOptions = expandNodeToNodeEncryptionOptions(s)
			}
		}

		if d.HasChange("off_peak_window_options") {
			input.OffPeakWindowOptions = expandOffPeakWindowOptions(d.Get("off_peak_window_options").([]interface{})[0].(map[string]interface{}))
		}

		if d.HasChange("snapshot_options") {
			options := d.Get("snapshot_options").([]interface{})

			if len(options) == 1 {
				o := options[0].(map[string]interface{})

				snapshotOptions := opensearchservice.SnapshotOptions{
					AutomatedSnapshotStartHour: aws.Int64(int64(o["automated_snapshot_start_hour"].(int))),
				}

				input.SnapshotOptions = &snapshotOptions
			}
		}

		if d.HasChange("vpc_options") {
			options := d.Get("vpc_options").([]interface{})
			s := options[0].(map[string]interface{})
			input.VPCOptions = expandVPCOptions(s)
		}

		_, err := conn.UpdateDomainConfigWithContext(ctx, &input)
		if err != nil {
			return sdkdiag.AppendErrorf(diags, "updating OpenSearch Domain (%s): %s", d.Id(), err)
		}

		if err := waitForDomainUpdate(ctx, conn, d.Get("domain_name").(string), d.Timeout(schema.TimeoutUpdate)); err != nil {
			return sdkdiag.AppendErrorf(diags, "updating OpenSearch Domain (%s): waiting for completion: %s", d.Id(), err)
		}

		if d.HasChange("engine_version") {
			upgradeInput := opensearchservice.UpgradeDomainInput{
				DomainName:    aws.String(d.Get("domain_name").(string)),
				TargetVersion: aws.String(d.Get("engine_version").(string)),
			}

			_, err := conn.UpgradeDomainWithContext(ctx, &upgradeInput)
			if err != nil {
				return sdkdiag.AppendErrorf(diags, "updating OpenSearch Domain (%s): upgrading: %s", d.Id(), err)
			}

			if _, err := waitUpgradeSucceeded(ctx, conn, d.Get("domain_name").(string), d.Timeout(schema.TimeoutUpdate)); err != nil {
				return sdkdiag.AppendErrorf(diags, "updating OpenSearch Domain (%s): upgrading: waiting for completion: %s", d.Id(), err)
			}
		}
	}

	return append(diags, resourceDomainRead(ctx, d, meta)...)
}

func resourceDomainDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).OpenSearchConn(ctx)
	domainName := d.Get("domain_name").(string)

	log.Printf("[DEBUG] Deleting OpenSearch Domain: %q", domainName)
	_, err := conn.DeleteDomainWithContext(ctx, &opensearchservice.DeleteDomainInput{
		DomainName: aws.String(domainName),
	})
	if err != nil {
		if tfawserr.ErrCodeEquals(err, opensearchservice.ErrCodeResourceNotFoundException) {
			return diags
		}
		return sdkdiag.AppendErrorf(diags, "deleting OpenSearch Domain (%s): %s", d.Id(), err)
	}

	log.Printf("[DEBUG] Waiting for OpenSearch Domain %q to be deleted", domainName)
	if err := waitForDomainDelete(ctx, conn, d.Get("domain_name").(string), d.Timeout(schema.TimeoutDelete)); err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting OpenSearch Domain (%s): waiting for completion: %s", d.Id(), err)
	}

	return diags
}

// inPlaceEncryptionEnableVersion returns true if, based on version, encryption
// can be enabled in place (without ForceNew)
func inPlaceEncryptionEnableVersion(version string) bool {
	if engineType, version, err := ParseEngineVersion(version); err == nil {
		switch engineType {
		case opensearchservice.EngineTypeElasticsearch:
			if verify.SemVerGreaterThanOrEqual(version, "6.7") {
				return true
			}
		case opensearchservice.EngineTypeOpenSearch:
			// All OpenSearch versions support enabling encryption in-place.
			return true
		}
	}

	return false
}

func suppressEquivalentKMSKeyIDs(k, old, new string, d *schema.ResourceData) bool {
	// The OpenSearch API accepts a short KMS key id but always returns the ARN of the key.
	// The ARN is of the format 'arn:aws:kms:REGION:ACCOUNT_ID:key/KMS_KEY_ID'.
	// These should be treated as equivalent.
	return strings.Contains(old, new)
}

func getDashboardEndpoint(d *schema.ResourceData) string {
	return d.Get("endpoint").(string) + "/_dashboards"
}

func getKibanaEndpoint(d *schema.ResourceData) string {
	return d.Get("endpoint").(string) + "/_plugin/kibana/"
}

func isDedicatedMasterDisabled(k, old, new string, d *schema.ResourceData) bool {
	v, ok := d.GetOk("cluster_config")
	if ok {
		clusterConfig := v.([]interface{})[0].(map[string]interface{})
		return !clusterConfig["dedicated_master_enabled"].(bool)
	}
	return false
}

func isCustomEndpointDisabled(k, old, new string, d *schema.ResourceData) bool {
	v, ok := d.GetOk("domain_endpoint_options")
	if ok {
		domainEndpointOptions := v.([]interface{})[0].(map[string]interface{})
		return !domainEndpointOptions["custom_endpoint_enabled"].(bool)
	}
	return false
}

func expandNodeToNodeEncryptionOptions(s map[string]interface{}) *opensearchservice.NodeToNodeEncryptionOptions {
	options := opensearchservice.NodeToNodeEncryptionOptions{}

	if v, ok := s["enabled"]; ok {
		options.Enabled = aws.Bool(v.(bool))
	}
	return &options
}

func flattenNodeToNodeEncryptionOptions(o *opensearchservice.NodeToNodeEncryptionOptions) []map[string]interface{} {
	if o == nil {
		return []map[string]interface{}{}
	}

	m := map[string]interface{}{}
	if o.Enabled != nil {
		m["enabled"] = aws.BoolValue(o.Enabled)
	}

	return []map[string]interface{}{m}
}

func expandClusterConfig(m map[string]interface{}) *opensearchservice.ClusterConfig {
	config := opensearchservice.ClusterConfig{}

	if v, ok := m["dedicated_master_enabled"]; ok {
		isEnabled := v.(bool)
		config.DedicatedMasterEnabled = aws.Bool(isEnabled)

		if isEnabled {
			if v, ok := m["dedicated_master_count"]; ok && v.(int) > 0 {
				config.DedicatedMasterCount = aws.Int64(int64(v.(int)))
			}
			if v, ok := m["dedicated_master_type"]; ok && v.(string) != "" {
				config.DedicatedMasterType = aws.String(v.(string))
			}
		}
	}

	if v, ok := m["instance_count"]; ok {
		config.InstanceCount = aws.Int64(int64(v.(int)))
	}
	if v, ok := m["instance_type"]; ok {
		config.InstanceType = aws.String(v.(string))
	}

	if v, ok := m["zone_awareness_enabled"]; ok {
		isEnabled := v.(bool)
		config.ZoneAwarenessEnabled = aws.Bool(isEnabled)

		if isEnabled {
			if v, ok := m["zone_awareness_config"]; ok {
				config.ZoneAwarenessConfig = expandZoneAwarenessConfig(v.([]interface{}))
			}
		}
	}

	if v, ok := m["cold_storage_options"]; ok {
		config.ColdStorageOptions = expandColdStorageOptions(v.([]interface{}))
	}

	if v, ok := m["warm_enabled"]; ok {
		isEnabled := v.(bool)
		config.WarmEnabled = aws.Bool(isEnabled)

		if isEnabled {
			if v, ok := m["warm_count"]; ok {
				config.WarmCount = aws.Int64(int64(v.(int)))
			}

			if v, ok := m["warm_type"]; ok {
				config.WarmType = aws.String(v.(string))
			}
		}
	}

	return &config
}

func expandZoneAwarenessConfig(l []interface{}) *opensearchservice.ZoneAwarenessConfig {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})

	zoneAwarenessConfig := &opensearchservice.ZoneAwarenessConfig{}

	if v, ok := m["availability_zone_count"]; ok && v.(int) > 0 {
		zoneAwarenessConfig.AvailabilityZoneCount = aws.Int64(int64(v.(int)))
	}

	return zoneAwarenessConfig
}

func expandColdStorageOptions(l []interface{}) *opensearchservice.ColdStorageOptions {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})

	ColdStorageOptions := &opensearchservice.ColdStorageOptions{}

	if v, ok := m["enabled"]; ok {
		ColdStorageOptions.Enabled = aws.Bool(v.(bool))
	}

	return ColdStorageOptions
}

func flattenClusterConfig(c *opensearchservice.ClusterConfig) []map[string]interface{} {
	m := map[string]interface{}{
		"zone_awareness_config":  flattenZoneAwarenessConfig(c.ZoneAwarenessConfig),
		"zone_awareness_enabled": aws.BoolValue(c.ZoneAwarenessEnabled),
	}

	if c.ColdStorageOptions != nil {
		m["cold_storage_options"] = flattenColdStorageOptions(c.ColdStorageOptions)
	}
	if c.DedicatedMasterCount != nil {
		m["dedicated_master_count"] = aws.Int64Value(c.DedicatedMasterCount)
	}
	if c.DedicatedMasterEnabled != nil {
		m["dedicated_master_enabled"] = aws.BoolValue(c.DedicatedMasterEnabled)
	}
	if c.DedicatedMasterType != nil {
		m["dedicated_master_type"] = aws.StringValue(c.DedicatedMasterType)
	}
	if c.InstanceCount != nil {
		m["instance_count"] = aws.Int64Value(c.InstanceCount)
	}
	if c.InstanceType != nil {
		m["instance_type"] = aws.StringValue(c.InstanceType)
	}
	if c.WarmEnabled != nil {
		m["warm_enabled"] = aws.BoolValue(c.WarmEnabled)
	}
	if c.WarmCount != nil {
		m["warm_count"] = aws.Int64Value(c.WarmCount)
	}
	if c.WarmType != nil {
		m["warm_type"] = aws.StringValue(c.WarmType)
	}

	return []map[string]interface{}{m}
}

func flattenZoneAwarenessConfig(zoneAwarenessConfig *opensearchservice.ZoneAwarenessConfig) []interface{} {
	if zoneAwarenessConfig == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"availability_zone_count": aws.Int64Value(zoneAwarenessConfig.AvailabilityZoneCount),
	}

	return []interface{}{m}
}

func flattenColdStorageOptions(coldStorageOptions *opensearchservice.ColdStorageOptions) []interface{} {
	if coldStorageOptions == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"enabled": aws.BoolValue(coldStorageOptions.Enabled),
	}

	return []interface{}{m}
}

// advancedOptionsIgnoreDefault checks for defaults in the n map and, if
// they don't exist in the o map, it deletes them. AWS returns default advanced
// options that cause perpetual diffs.
func advancedOptionsIgnoreDefault(o map[string]interface{}, n map[string]interface{}) map[string]interface{} {
	for k, v := range n {
		switch fmt.Sprintf("%s=%s", k, v) {
		case "override_main_response_version=false":
			if _, ok := o[k]; !ok {
				delete(n, "override_main_response_version")
			}
		case "rest.action.multi.allow_explicit_index=true":
			if _, ok := o[k]; !ok {
				delete(n, "rest.action.multi.allow_explicit_index")
			}
		}
	}

	return n
}

// ParseEngineVersion parses a domain's engine version string into engine type and semver string.
// engine_version is a string of format Elasticsearch_X.Y or OpenSearch_X.Y.
func ParseEngineVersion(engineVersion string) (string, string, error) {
	parts := strings.Split(engineVersion, "_")

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format for engine version (%s)", engineVersion)
	}

	return parts[0], parts[1], nil
}

// EBSVolumeTypePermitsIopsInput returns true if the volume type supports the Iops input
//
// This check prevents a ValidationException when updating EBS volume types from a value
// that supports IOPS (ex. gp3) to one that doesn't (ex. gp2).
func EBSVolumeTypePermitsIopsInput(volumeType string) bool {
	permittedTypes := []string{opensearchservice.VolumeTypeGp3, opensearchservice.VolumeTypeIo1}
	for _, t := range permittedTypes {
		if volumeType == t {
			return true
		}
	}
	return false
}

// EBSVolumeTypePermitsThroughputInput returns true if the volume type supports the Throughput input
//
// This check prevents a ValidationException when updating EBS volume types from a value
// that supports Throughput (ex. gp3) to one that doesn't (ex. gp2).
func EBSVolumeTypePermitsThroughputInput(volumeType string) bool {
	permittedTypes := []string{opensearchservice.VolumeTypeGp3}
	for _, t := range permittedTypes {
		if volumeType == t {
			return true
		}
	}
	return false
}
