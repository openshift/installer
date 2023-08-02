package quicksight

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/quicksight"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKDataSource("aws_quicksight_data_set", name="Data Set")
func DataSourceDataSet() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceDataSetRead,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"aws_account_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: verify.ValidAccountID,
			},
			"column_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"geo_spatial_column_group": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"columns": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"country_code": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"column_level_permission_rules": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"column_names": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"principals": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"data_set_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data_set_usage_configuration": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disable_use_as_direct_query_source": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"disable_use_as_imported_source": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"field_folders": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field_folders_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"columns": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"import_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"logical_table_map": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     logicalTableMapDataSourceSchema(),
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"permissions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"actions": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"principal": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"physical_table_map": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     physicalTableMapDataSourceSchema(),
			},
			"row_level_permission_data_set": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"format_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"permission_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"row_level_permission_tag_configuration": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tag_rules": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"column_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"match_all_value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"tag_key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"tag_multi_value_delimiter": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"tags": tftags.TagsSchemaComputed(),
			"tags_all": {
				Type:       schema.TypeMap,
				Optional:   true,
				Computed:   true,
				Elem:       &schema.Schema{Type: schema.TypeString},
				Deprecated: `this attribute has been deprecated`,
			},
		},
	}
}

func logicalTableMapDataSourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"alias": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_transforms": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cast_column_type_operation": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"column_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"format": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"new_column_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"create_columns_operation": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"columns": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"column_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"column_name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"expression": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"filter_operation": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"condition_expression": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"project_operation": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"projected_columns": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"rename_column_operation": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"column_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"new_column_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"tag_column_operation": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"column_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"tags": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"column_description": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"text": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
												"column_geographic_role": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"untag_column_operation": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"column_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"tag_names": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
			"logical_table_map_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"data_set_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"join_instruction": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"left_join_key_properties": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"unique_key": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
									"left_operand": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"on_clause": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"right_join_key_properties": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"unique_key": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
									"right_operand": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"physical_table_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func physicalTableMapDataSourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"custom_sql": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"columns": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"data_source_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sql_query": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"physical_table_map_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"relational_table": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"catalog": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_source_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"input_columns": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"schema": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"s3_source": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"data_source_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"input_columns": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"upload_settings": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"contains_header": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"delimiter": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"format": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"start_from_row": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"text_qualifier": {
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

const (
	DSNameDataSet = "Data Set Data Source"
)

func dataSourceDataSetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).QuickSightConn(ctx)
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	awsAccountId := meta.(*conns.AWSClient).AccountID
	if v, ok := d.GetOk("aws_account_id"); ok {
		awsAccountId = v.(string)
	}
	dataSetId := d.Get("data_set_id").(string)

	descOpts := &quicksight.DescribeDataSetInput{
		AwsAccountId: aws.String(awsAccountId),
		DataSetId:    aws.String(dataSetId),
	}

	output, err := conn.DescribeDataSetWithContext(ctx, descOpts)
	if err != nil {
		return create.DiagError(names.QuickSight, create.ErrActionReading, DSNameDataSet, dataSetId, err)
	}

	dataSet := output.DataSet

	d.SetId(createDataSetID(awsAccountId, dataSetId))

	d.Set("arn", dataSet.Arn)
	d.Set("aws_account_id", awsAccountId)
	d.Set("data_set_id", dataSet.DataSetId)
	d.Set("name", dataSet.Name)
	d.Set("import_mode", dataSet.ImportMode)

	if err := d.Set("column_groups", flattenColumnGroups(dataSet.ColumnGroups)); err != nil {
		return diag.Errorf("setting column_groups: %s", err)
	}

	if err := d.Set("column_level_permission_rules", flattenColumnLevelPermissionRules(dataSet.ColumnLevelPermissionRules)); err != nil {
		return diag.Errorf("setting column_level_permission_rules: %s", err)
	}

	if err := d.Set("data_set_usage_configuration", flattenDataSetUsageConfiguration(dataSet.DataSetUsageConfiguration)); err != nil {
		return diag.Errorf("setting data_set_usage_configuration: %s", err)
	}

	if err := d.Set("field_folders", flattenFieldFolders(dataSet.FieldFolders)); err != nil {
		return diag.Errorf("setting field_folders: %s", err)
	}

	if err := d.Set("logical_table_map", flattenLogicalTableMap(dataSet.LogicalTableMap, logicalTableMapDataSourceSchema())); err != nil {
		return diag.Errorf("setting logical_table_map: %s", err)
	}

	if err := d.Set("physical_table_map", flattenPhysicalTableMap(dataSet.PhysicalTableMap, physicalTableMapDataSourceSchema())); err != nil {
		return diag.Errorf("setting physical_table_map: %s", err)
	}

	if err := d.Set("row_level_permission_data_set", flattenRowLevelPermissionDataSet(dataSet.RowLevelPermissionDataSet)); err != nil {
		return diag.Errorf("setting row_level_permission_data_set: %s", err)
	}

	if err := d.Set("row_level_permission_tag_configuration", flattenRowLevelPermissionTagConfiguration(dataSet.RowLevelPermissionTagConfiguration)); err != nil {
		return diag.Errorf("setting row_level_permission_tag_configuration: %s", err)
	}

	tags, err := ListTags(ctx, conn, d.Get("arn").(string))

	if err != nil {
		return diag.Errorf("listing tags for QuickSight Data Set (%s): %s", d.Id(), err)
	}

	tags = tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return diag.Errorf("setting tags: %s", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return diag.Errorf("setting tags_all: %s", err)
	}

	permsResp, err := conn.DescribeDataSetPermissionsWithContext(ctx, &quicksight.DescribeDataSetPermissionsInput{
		AwsAccountId: aws.String(awsAccountId),
		DataSetId:    aws.String(dataSetId),
	})

	if err != nil {
		return diag.Errorf("describing QuickSight Data Source (%s) Permissions: %s", d.Id(), err)
	}

	if err := d.Set("permissions", flattenPermissions(permsResp.Permissions)); err != nil {
		return diag.Errorf("setting permissions: %s", err)
	}
	return nil
}
