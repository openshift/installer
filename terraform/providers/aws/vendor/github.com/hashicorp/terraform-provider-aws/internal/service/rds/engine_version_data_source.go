package rds

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/generate/namevaluesfilters"
)

// @SDKDataSource("aws_rds_engine_version")
func DataSourceEngineVersion() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceEngineVersionRead,
		Schema: map[string]*schema.Schema{
			"default_character_set": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"default_only": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"engine": {
				Type:     schema.TypeString,
				Required: true,
			},

			"engine_description": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"exportable_log_types": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Set:      schema.HashString,
			},

			"filter": namevaluesfilters.Schema(),

			"include_all": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"parameter_group_family": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"preferred_versions": {
				Type:          schema.TypeList,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"version"},
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"supported_character_sets": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Set:      schema.HashString,
			},

			"supported_feature_names": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Set:      schema.HashString,
			},

			"supported_modes": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Set:      schema.HashString,
			},

			"supported_timezones": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Set:      schema.HashString,
			},

			"supports_global_databases": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"supports_log_exports_to_cloudwatch": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"supports_parallel_query": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"supports_read_replica": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"valid_upgrade_targets": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Set:      schema.HashString,
			},

			"version": {
				Type:          schema.TypeString,
				Computed:      true,
				Optional:      true,
				ConflictsWith: []string{"preferred_versions"},
			},

			"version_description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceEngineVersionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).RDSConn(ctx)

	input := &rds.DescribeDBEngineVersionsInput{
		ListSupportedCharacterSets: aws.Bool(true),
		ListSupportedTimezones:     aws.Bool(true),
	}

	if v, ok := d.GetOk("engine"); ok {
		input.Engine = aws.String(v.(string))
	}

	if v, ok := d.GetOk("filter"); ok {
		input.Filters = namevaluesfilters.New(v.(*schema.Set)).RDSFilters()
	}

	if v, ok := d.GetOk("include_all"); ok {
		input.IncludeAll = aws.Bool(v.(bool))
	}

	if v, ok := d.GetOk("parameter_group_family"); ok {
		input.DBParameterGroupFamily = aws.String(v.(string))
	}

	if v, ok := d.GetOk("version"); ok {
		input.EngineVersion = aws.String(v.(string))
	}

	if v, ok := d.GetOk("default_only"); ok {
		input.DefaultOnly = aws.Bool(v.(bool))
	} else if _, ok := d.GetOk("version"); !ok {
		if _, ok := d.GetOk("preferred_versions"); !ok {
			input.DefaultOnly = aws.Bool(true)
		}
	}

	log.Printf("[DEBUG] Reading RDS engine versions: %v", input)
	var engineVersions []*rds.DBEngineVersion

	err := conn.DescribeDBEngineVersionsPagesWithContext(ctx, input, func(resp *rds.DescribeDBEngineVersionsOutput, lastPage bool) bool {
		for _, engineVersion := range resp.DBEngineVersions {
			if engineVersion == nil {
				continue
			}

			engineVersions = append(engineVersions, engineVersion)
		}
		return !lastPage
	})
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading RDS engine versions: %s", err)
	}

	if len(engineVersions) == 0 {
		return sdkdiag.AppendErrorf(diags, "no RDS engine versions found")
	}

	// preferred versions
	var found *rds.DBEngineVersion
	if l := d.Get("preferred_versions").([]interface{}); len(l) > 0 {
		for _, elem := range l {
			preferredVersion, ok := elem.(string)

			if !ok {
				continue
			}

			for _, engineVersion := range engineVersions {
				if preferredVersion == aws.StringValue(engineVersion.EngineVersion) {
					found = engineVersion
					break
				}
			}

			if found != nil {
				break
			}
		}
	}

	if found == nil && len(engineVersions) > 1 {
		return sdkdiag.AppendErrorf(diags, "multiple RDS engine versions (%v) match the criteria", engineVersions)
	}

	if found == nil && len(engineVersions) == 1 {
		found = engineVersions[0]
	}

	if found == nil {
		return sdkdiag.AppendErrorf(diags, "no RDS engine versions match the criteria")
	}

	d.SetId(aws.StringValue(found.EngineVersion))

	if found.DefaultCharacterSet != nil {
		d.Set("default_character_set", found.DefaultCharacterSet.CharacterSetName)
	}

	d.Set("engine", found.Engine)
	d.Set("engine_description", found.DBEngineDescription)
	d.Set("exportable_log_types", found.ExportableLogTypes)
	d.Set("parameter_group_family", found.DBParameterGroupFamily)
	d.Set("status", found.Status)

	var characterSets []string
	for _, cs := range found.SupportedCharacterSets {
		characterSets = append(characterSets, aws.StringValue(cs.CharacterSetName))
	}
	d.Set("supported_character_sets", characterSets)

	d.Set("supported_feature_names", found.SupportedFeatureNames)
	d.Set("supported_modes", found.SupportedEngineModes)

	var timezones []string
	for _, tz := range found.SupportedTimezones {
		timezones = append(timezones, aws.StringValue(tz.TimezoneName))
	}
	d.Set("supported_timezones", timezones)

	d.Set("supports_global_databases", found.SupportsGlobalDatabases)
	d.Set("supports_log_exports_to_cloudwatch", found.SupportsLogExportsToCloudwatchLogs)
	d.Set("supports_parallel_query", found.SupportsParallelQuery)
	d.Set("supports_read_replica", found.SupportsReadReplica)

	var upgradeTargets []string
	for _, ut := range found.ValidUpgradeTarget {
		upgradeTargets = append(upgradeTargets, aws.StringValue(ut.EngineVersion))
	}
	d.Set("valid_upgrade_targets", upgradeTargets)

	d.Set("version", found.EngineVersion)
	d.Set("version_description", found.DBEngineVersionDescription)

	return diags
}
