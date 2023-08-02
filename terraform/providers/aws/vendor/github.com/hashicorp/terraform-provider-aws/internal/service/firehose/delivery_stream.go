package firehose

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/firehose"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

const (
	destinationTypeExtendedS3    = "extended_s3"
	destinationTypeElasticsearch = "elasticsearch"
	destinationTypeOpensearch    = "opensearch"
	destinationTypeRedshift      = "redshift"
	destinationTypeSplunk        = "splunk"
	destinationTypeHTTPEndpoint  = "http_endpoint"
)

func cloudWatchLoggingOptionsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"enabled": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},

				"log_group_name": {
					Type:     schema.TypeString,
					Optional: true,
				},

				"log_stream_name": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func dynamicPartitioningConfigurationSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Optional: true,
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"enabled": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
					ForceNew: true,
				},
				"retry_duration": {
					Type:         schema.TypeInt,
					Optional:     true,
					Default:      300,
					ValidateFunc: validation.IntBetween(0, 7200),
				},
			},
		},
	}
}

func requestConfigurationSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"content_encoding": {
					Type:         schema.TypeString,
					Optional:     true,
					Default:      firehose.ContentEncodingNone,
					ValidateFunc: validation.StringInSlice(firehose.ContentEncoding_Values(), false),
				},

				"common_attributes": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"name": {
								Type:     schema.TypeString,
								Required: true,
							},
							"value": {
								Type:     schema.TypeString,
								Required: true,
							},
						},
					},
				},
			},
		},
	}
}

func s3BackupConfigurationSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Optional: true,
		Elem:     s3ConfigurationElem(),
	}
}

func s3ConfigurationSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Required: true,
		Elem:     s3ConfigurationElem(),
	}
}

func s3ConfigurationElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"bucket_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: verify.ValidARN,
			},

			"buffering_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      5,
				ValidateFunc: validation.IntAtLeast(1),
			},

			"buffering_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      300,
				ValidateFunc: validation.IntAtLeast(60),
			},

			"compression_format": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      firehose.CompressionFormatUncompressed,
				ValidateFunc: validation.StringInSlice(firehose.CompressionFormat_Values(), false),
			},

			"error_output_prefix": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 1024),
			},

			"kms_key_arn": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: verify.ValidARN,
			},

			"role_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: verify.ValidARN,
			},

			"prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"cloudwatch_logging_options": cloudWatchLoggingOptionsSchema(),
		},
	}
}

func processingConfigurationSchema() *schema.Schema {
	return &schema.Schema{
		Type:             schema.TypeList,
		Optional:         true,
		MaxItems:         1,
		DiffSuppressFunc: verify.SuppressMissingOptionalConfigurationBlock,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"enabled": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"processors": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"parameters": {
								Type:     schema.TypeList,
								Optional: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"parameter_name": {
											Type:         schema.TypeString,
											Required:     true,
											ValidateFunc: validation.StringInSlice(firehose.ProcessorParameterName_Values(), false),
										},
										"parameter_value": {
											Type:         schema.TypeString,
											Required:     true,
											ValidateFunc: validation.StringLenBetween(1, 5120),
										},
									},
								},
							},
							"type": {
								Type:         schema.TypeString,
								Required:     true,
								ValidateFunc: validation.StringInSlice(firehose.ProcessorType_Values(), false),
							},
						},
					},
				},
			},
		},
	}
}

func flattenCloudWatchLoggingOptions(clo *firehose.CloudWatchLoggingOptions) []interface{} {
	if clo == nil {
		return []interface{}{}
	}

	cloudwatchLoggingOptions := map[string]interface{}{
		"enabled": aws.BoolValue(clo.Enabled),
	}
	if aws.BoolValue(clo.Enabled) {
		cloudwatchLoggingOptions["log_group_name"] = aws.StringValue(clo.LogGroupName)
		cloudwatchLoggingOptions["log_stream_name"] = aws.StringValue(clo.LogStreamName)
	}
	return []interface{}{cloudwatchLoggingOptions}
}

func flattenElasticsearchConfiguration(description *firehose.ElasticsearchDestinationDescription) []map[string]interface{} {
	if description == nil {
		return []map[string]interface{}{}
	}

	m := map[string]interface{}{
		"cloudwatch_logging_options": flattenCloudWatchLoggingOptions(description.CloudWatchLoggingOptions),
		"role_arn":                   aws.StringValue(description.RoleARN),
		"type_name":                  aws.StringValue(description.TypeName),
		"index_name":                 aws.StringValue(description.IndexName),
		"s3_backup_mode":             aws.StringValue(description.S3BackupMode),
		"s3_configuration":           flattenS3Configuration(description.S3DestinationDescription),
		"index_rotation_period":      aws.StringValue(description.IndexRotationPeriod),
		"vpc_config":                 flattenVPCConfiguration(description.VpcConfigurationDescription),
		"processing_configuration":   flattenProcessingConfiguration(description.ProcessingConfiguration, aws.StringValue(description.RoleARN)),
	}

	if description.DomainARN != nil {
		m["domain_arn"] = aws.StringValue(description.DomainARN)
	}

	if description.ClusterEndpoint != nil {
		m["cluster_endpoint"] = aws.StringValue(description.ClusterEndpoint)
	}

	if description.BufferingHints != nil {
		m["buffering_interval"] = int(aws.Int64Value(description.BufferingHints.IntervalInSeconds))
		m["buffering_size"] = int(aws.Int64Value(description.BufferingHints.SizeInMBs))
	}

	if description.RetryOptions != nil {
		m["retry_duration"] = int(aws.Int64Value(description.RetryOptions.DurationInSeconds))
	}

	return []map[string]interface{}{m}
}

func flattenOpensearchConfiguration(description *firehose.AmazonopensearchserviceDestinationDescription) []map[string]interface{} {
	if description == nil {
		return []map[string]interface{}{}
	}

	m := map[string]interface{}{
		"cloudwatch_logging_options": flattenCloudWatchLoggingOptions(description.CloudWatchLoggingOptions),
		"role_arn":                   aws.StringValue(description.RoleARN),
		"type_name":                  aws.StringValue(description.TypeName),
		"index_name":                 aws.StringValue(description.IndexName),
		"s3_backup_mode":             aws.StringValue(description.S3BackupMode),
		"s3_configuration":           flattenS3Configuration(description.S3DestinationDescription),
		"index_rotation_period":      aws.StringValue(description.IndexRotationPeriod),
		"vpc_config":                 flattenVPCConfiguration(description.VpcConfigurationDescription),
		"processing_configuration":   flattenProcessingConfiguration(description.ProcessingConfiguration, aws.StringValue(description.RoleARN)),
	}

	if description.DomainARN != nil {
		m["domain_arn"] = aws.StringValue(description.DomainARN)
	}

	if description.ClusterEndpoint != nil {
		m["cluster_endpoint"] = aws.StringValue(description.ClusterEndpoint)
	}

	if description.BufferingHints != nil {
		m["buffering_interval"] = int(aws.Int64Value(description.BufferingHints.IntervalInSeconds))
		m["buffering_size"] = int(aws.Int64Value(description.BufferingHints.SizeInMBs))
	}

	if description.RetryOptions != nil {
		m["retry_duration"] = int(aws.Int64Value(description.RetryOptions.DurationInSeconds))
	}

	return []map[string]interface{}{m}
}

func flattenVPCConfiguration(description *firehose.VpcConfigurationDescription) []map[string]interface{} {
	if description == nil {
		return []map[string]interface{}{}
	}

	m := map[string]interface{}{
		"vpc_id":             aws.StringValue(description.VpcId),
		"subnet_ids":         flex.FlattenStringSet(description.SubnetIds),
		"security_group_ids": flex.FlattenStringSet(description.SecurityGroupIds),
		"role_arn":           aws.StringValue(description.RoleARN),
	}

	return []map[string]interface{}{m}
}

func flattenExtendedS3Configuration(description *firehose.ExtendedS3DestinationDescription) []map[string]interface{} {
	if description == nil {
		return []map[string]interface{}{}
	}

	m := map[string]interface{}{
		"bucket_arn":                           aws.StringValue(description.BucketARN),
		"cloudwatch_logging_options":           flattenCloudWatchLoggingOptions(description.CloudWatchLoggingOptions),
		"compression_format":                   aws.StringValue(description.CompressionFormat),
		"data_format_conversion_configuration": flattenDataFormatConversionConfiguration(description.DataFormatConversionConfiguration),
		"error_output_prefix":                  aws.StringValue(description.ErrorOutputPrefix),
		"prefix":                               aws.StringValue(description.Prefix),
		"processing_configuration":             flattenProcessingConfiguration(description.ProcessingConfiguration, aws.StringValue(description.RoleARN)),
		"dynamic_partitioning_configuration":   flattenDynamicPartitioningConfiguration(description.DynamicPartitioningConfiguration),
		"role_arn":                             aws.StringValue(description.RoleARN),
		"s3_backup_configuration":              flattenS3Configuration(description.S3BackupDescription),
		"s3_backup_mode":                       aws.StringValue(description.S3BackupMode),
	}

	if description.BufferingHints != nil {
		m["buffering_interval"] = int(aws.Int64Value(description.BufferingHints.IntervalInSeconds))
		m["buffering_size"] = int(aws.Int64Value(description.BufferingHints.SizeInMBs))
	}

	if description.EncryptionConfiguration != nil && description.EncryptionConfiguration.KMSEncryptionConfig != nil {
		m["kms_key_arn"] = aws.StringValue(description.EncryptionConfiguration.KMSEncryptionConfig.AWSKMSKeyARN)
	}

	return []map[string]interface{}{m}
}

func flattenRedshiftConfiguration(description *firehose.RedshiftDestinationDescription, configuredPassword string) []map[string]interface{} {
	if description == nil {
		return []map[string]interface{}{}
	}

	m := map[string]interface{}{
		"cloudwatch_logging_options": flattenCloudWatchLoggingOptions(description.CloudWatchLoggingOptions),
		"cluster_jdbcurl":            aws.StringValue(description.ClusterJDBCURL),
		"password":                   configuredPassword,
		"processing_configuration":   flattenProcessingConfiguration(description.ProcessingConfiguration, aws.StringValue(description.RoleARN)),
		"role_arn":                   aws.StringValue(description.RoleARN),
		"s3_backup_configuration":    flattenS3Configuration(description.S3BackupDescription),
		"s3_backup_mode":             aws.StringValue(description.S3BackupMode),
		"s3_configuration":           flattenS3Configuration(description.S3DestinationDescription),
		"username":                   aws.StringValue(description.Username),
	}

	if description.CopyCommand != nil {
		m["copy_options"] = aws.StringValue(description.CopyCommand.CopyOptions)
		m["data_table_columns"] = aws.StringValue(description.CopyCommand.DataTableColumns)
		m["data_table_name"] = aws.StringValue(description.CopyCommand.DataTableName)
	}

	if description.RetryOptions != nil {
		m["retry_duration"] = int(aws.Int64Value(description.RetryOptions.DurationInSeconds))
	}

	return []map[string]interface{}{m}
}

func flattenSplunkConfiguration(description *firehose.SplunkDestinationDescription) []map[string]interface{} {
	if description == nil {
		return []map[string]interface{}{}
	}
	m := map[string]interface{}{
		"cloudwatch_logging_options": flattenCloudWatchLoggingOptions(description.CloudWatchLoggingOptions),
		"hec_acknowledgment_timeout": int(aws.Int64Value(description.HECAcknowledgmentTimeoutInSeconds)),
		"hec_endpoint_type":          aws.StringValue(description.HECEndpointType),
		"hec_endpoint":               aws.StringValue(description.HECEndpoint),
		"hec_token":                  aws.StringValue(description.HECToken),
		"processing_configuration":   flattenProcessingConfiguration(description.ProcessingConfiguration, ""),
		"s3_backup_mode":             aws.StringValue(description.S3BackupMode),
		"s3_configuration":           flattenS3Configuration(description.S3DestinationDescription),
	}

	if description.RetryOptions != nil {
		m["retry_duration"] = int(aws.Int64Value(description.RetryOptions.DurationInSeconds))
	}

	return []map[string]interface{}{m}
}

func flattenS3Configuration(description *firehose.S3DestinationDescription) []map[string]interface{} {
	if description == nil {
		return []map[string]interface{}{}
	}

	m := map[string]interface{}{
		"bucket_arn":                 aws.StringValue(description.BucketARN),
		"cloudwatch_logging_options": flattenCloudWatchLoggingOptions(description.CloudWatchLoggingOptions),
		"compression_format":         aws.StringValue(description.CompressionFormat),
		"error_output_prefix":        aws.StringValue(description.ErrorOutputPrefix),
		"prefix":                     aws.StringValue(description.Prefix),
		"role_arn":                   aws.StringValue(description.RoleARN),
	}

	if description.BufferingHints != nil {
		m["buffering_interval"] = int(aws.Int64Value(description.BufferingHints.IntervalInSeconds))
		m["buffering_size"] = int(aws.Int64Value(description.BufferingHints.SizeInMBs))
	}

	if description.EncryptionConfiguration != nil && description.EncryptionConfiguration.KMSEncryptionConfig != nil {
		m["kms_key_arn"] = aws.StringValue(description.EncryptionConfiguration.KMSEncryptionConfig.AWSKMSKeyARN)
	}

	return []map[string]interface{}{m}
}

func flattenDataFormatConversionConfiguration(dfcc *firehose.DataFormatConversionConfiguration) []map[string]interface{} {
	if dfcc == nil {
		return []map[string]interface{}{}
	}

	enabled := aws.BoolValue(dfcc.Enabled)
	ifc := flattenInputFormatConfiguration(dfcc.InputFormatConfiguration)
	ofc := flattenOutputFormatConfiguration(dfcc.OutputFormatConfiguration)
	sc := flattenSchemaConfiguration(dfcc.SchemaConfiguration)

	// The AWS SDK can represent "no data format conversion configuration" in two ways:
	// 1. With a nil value
	// 2. With enabled set to false and nil for ALL the config sections.
	// We normalize this with an empty configuration in the state due
	// to the existing Default: true on the enabled attribute.
	if !enabled && len(ifc) == 0 && len(ofc) == 0 && len(sc) == 0 {
		return []map[string]interface{}{}
	}

	m := map[string]interface{}{
		"enabled":                     enabled,
		"input_format_configuration":  ifc,
		"output_format_configuration": ofc,
		"schema_configuration":        sc,
	}

	return []map[string]interface{}{m}
}

func flattenInputFormatConfiguration(ifc *firehose.InputFormatConfiguration) []map[string]interface{} {
	if ifc == nil {
		return []map[string]interface{}{}
	}

	m := map[string]interface{}{
		"deserializer": flattenDeserializer(ifc.Deserializer),
	}

	return []map[string]interface{}{m}
}

func flattenDeserializer(deserializer *firehose.Deserializer) []map[string]interface{} {
	if deserializer == nil {
		return []map[string]interface{}{}
	}

	m := map[string]interface{}{
		"hive_json_ser_de":   flattenHiveJSONSerDe(deserializer.HiveJsonSerDe),
		"open_x_json_ser_de": flattenOpenXJSONSerDe(deserializer.OpenXJsonSerDe),
	}

	return []map[string]interface{}{m}
}

func flattenHiveJSONSerDe(hjsd *firehose.HiveJsonSerDe) []map[string]interface{} {
	if hjsd == nil {
		return []map[string]interface{}{}
	}

	m := map[string]interface{}{
		"timestamp_formats": flex.FlattenStringList(hjsd.TimestampFormats),
	}

	return []map[string]interface{}{m}
}

func flattenOpenXJSONSerDe(oxjsd *firehose.OpenXJsonSerDe) []map[string]interface{} {
	if oxjsd == nil {
		return []map[string]interface{}{}
	}

	m := map[string]interface{}{
		"column_to_json_key_mappings":              aws.StringValueMap(oxjsd.ColumnToJsonKeyMappings),
		"convert_dots_in_json_keys_to_underscores": aws.BoolValue(oxjsd.ConvertDotsInJsonKeysToUnderscores),
	}

	// API omits default values
	// Return defaults that are not type zero values to prevent extraneous difference

	m["case_insensitive"] = true
	if oxjsd.CaseInsensitive != nil {
		m["case_insensitive"] = aws.BoolValue(oxjsd.CaseInsensitive)
	}

	return []map[string]interface{}{m}
}

func flattenOutputFormatConfiguration(ofc *firehose.OutputFormatConfiguration) []map[string]interface{} {
	if ofc == nil {
		return []map[string]interface{}{}
	}

	m := map[string]interface{}{
		"serializer": flattenSerializer(ofc.Serializer),
	}

	return []map[string]interface{}{m}
}

func flattenSerializer(serializer *firehose.Serializer) []map[string]interface{} {
	if serializer == nil {
		return []map[string]interface{}{}
	}

	m := map[string]interface{}{
		"orc_ser_de":     flattenOrcSerDe(serializer.OrcSerDe),
		"parquet_ser_de": flattenParquetSerDe(serializer.ParquetSerDe),
	}

	return []map[string]interface{}{m}
}

func flattenOrcSerDe(osd *firehose.OrcSerDe) []map[string]interface{} {
	if osd == nil {
		return []map[string]interface{}{}
	}

	m := map[string]interface{}{
		"bloom_filter_columns":     aws.StringValueSlice(osd.BloomFilterColumns),
		"dictionary_key_threshold": aws.Float64Value(osd.DictionaryKeyThreshold),
		"enable_padding":           aws.BoolValue(osd.EnablePadding),
	}

	// API omits default values
	// Return defaults that are not type zero values to prevent extraneous difference

	m["block_size_bytes"] = 268435456
	if osd.BlockSizeBytes != nil {
		m["block_size_bytes"] = int(aws.Int64Value(osd.BlockSizeBytes))
	}

	m["bloom_filter_false_positive_probability"] = 0.05
	if osd.BloomFilterFalsePositiveProbability != nil {
		m["bloom_filter_false_positive_probability"] = aws.Float64Value(osd.BloomFilterFalsePositiveProbability)
	}

	m["compression"] = firehose.OrcCompressionSnappy
	if osd.Compression != nil {
		m["compression"] = aws.StringValue(osd.Compression)
	}

	m["format_version"] = firehose.OrcFormatVersionV012
	if osd.FormatVersion != nil {
		m["format_version"] = aws.StringValue(osd.FormatVersion)
	}

	m["padding_tolerance"] = 0.05
	if osd.PaddingTolerance != nil {
		m["padding_tolerance"] = aws.Float64Value(osd.PaddingTolerance)
	}

	m["row_index_stride"] = 10000
	if osd.RowIndexStride != nil {
		m["row_index_stride"] = int(aws.Int64Value(osd.RowIndexStride))
	}

	m["stripe_size_bytes"] = 67108864
	if osd.StripeSizeBytes != nil {
		m["stripe_size_bytes"] = int(aws.Int64Value(osd.StripeSizeBytes))
	}

	return []map[string]interface{}{m}
}

func flattenParquetSerDe(psd *firehose.ParquetSerDe) []map[string]interface{} {
	if psd == nil {
		return []map[string]interface{}{}
	}

	m := map[string]interface{}{
		"enable_dictionary_compression": aws.BoolValue(psd.EnableDictionaryCompression),
		"max_padding_bytes":             int(aws.Int64Value(psd.MaxPaddingBytes)),
	}

	// API omits default values
	// Return defaults that are not type zero values to prevent extraneous difference

	m["block_size_bytes"] = 268435456
	if psd.BlockSizeBytes != nil {
		m["block_size_bytes"] = int(aws.Int64Value(psd.BlockSizeBytes))
	}

	m["compression"] = firehose.ParquetCompressionSnappy
	if psd.Compression != nil {
		m["compression"] = aws.StringValue(psd.Compression)
	}

	m["page_size_bytes"] = 1048576
	if psd.PageSizeBytes != nil {
		m["page_size_bytes"] = int(aws.Int64Value(psd.PageSizeBytes))
	}

	m["writer_version"] = firehose.ParquetWriterVersionV1
	if psd.WriterVersion != nil {
		m["writer_version"] = aws.StringValue(psd.WriterVersion)
	}

	return []map[string]interface{}{m}
}

func flattenSchemaConfiguration(sc *firehose.SchemaConfiguration) []map[string]interface{} {
	if sc == nil {
		return []map[string]interface{}{}
	}

	m := map[string]interface{}{
		"catalog_id":    aws.StringValue(sc.CatalogId),
		"database_name": aws.StringValue(sc.DatabaseName),
		"region":        aws.StringValue(sc.Region),
		"role_arn":      aws.StringValue(sc.RoleARN),
		"table_name":    aws.StringValue(sc.TableName),
		"version_id":    aws.StringValue(sc.VersionId),
	}

	return []map[string]interface{}{m}
}

func flattenRequestConfiguration(rc *firehose.HttpEndpointRequestConfiguration) []map[string]interface{} {
	if rc == nil {
		return []map[string]interface{}{}
	}

	requestConfiguration := make([]map[string]interface{}, 1)

	commonAttributes := make([]interface{}, 0)
	for _, params := range rc.CommonAttributes {
		name := aws.StringValue(params.AttributeName)
		value := aws.StringValue(params.AttributeValue)

		commonAttributes = append(commonAttributes, map[string]interface{}{
			"name":  name,
			"value": value,
		})
	}

	requestConfiguration[0] = map[string]interface{}{
		"common_attributes": commonAttributes,
		"content_encoding":  aws.StringValue(rc.ContentEncoding),
	}

	return requestConfiguration
}

func flattenProcessingConfiguration(pc *firehose.ProcessingConfiguration, roleArn string) []map[string]interface{} {
	if pc == nil {
		return []map[string]interface{}{}
	}

	processingConfiguration := make([]map[string]interface{}, 1)

	// It is necessary to explicitly filter this out
	// to prevent diffs during routine use and retain the ability
	// to show diffs if any field has drifted
	defaultLambdaParams := map[string]string{
		"NumberOfRetries":         "3",
		"RoleArn":                 roleArn,
		"BufferSizeInMBs":         "3",
		"BufferIntervalInSeconds": "60",
	}

	processors := make([]interface{}, len(pc.Processors))
	for i, p := range pc.Processors {
		t := aws.StringValue(p.Type)
		parameters := make([]interface{}, 0)

		for _, params := range p.Parameters {
			name := aws.StringValue(params.ParameterName)
			value := aws.StringValue(params.ParameterValue)

			if t == firehose.ProcessorTypeLambda {
				// Ignore defaults
				if v, ok := defaultLambdaParams[name]; ok && v == value {
					continue
				}
			}

			parameters = append(parameters, map[string]interface{}{
				"parameter_name":  name,
				"parameter_value": value,
			})
		}

		processors[i] = map[string]interface{}{
			"type":       t,
			"parameters": parameters,
		}
	}
	processingConfiguration[0] = map[string]interface{}{
		"enabled":    aws.BoolValue(pc.Enabled),
		"processors": processors,
	}
	return processingConfiguration
}

func flattenDynamicPartitioningConfiguration(dpc *firehose.DynamicPartitioningConfiguration) []map[string]interface{} {
	if dpc == nil {
		return []map[string]interface{}{}
	}

	dynamicPartitioningConfiguration := make([]map[string]interface{}, 1)

	dynamicPartitioningConfiguration[0] = map[string]interface{}{
		"enabled": aws.BoolValue(dpc.Enabled),
	}

	if dpc.RetryOptions != nil && dpc.RetryOptions.DurationInSeconds != nil {
		dynamicPartitioningConfiguration[0]["retry_duration"] = int(aws.Int64Value(dpc.RetryOptions.DurationInSeconds))
	}

	return dynamicPartitioningConfiguration
}

func flattenSourceConfiguration(desc *firehose.KinesisStreamSourceDescription) []interface{} {
	if desc == nil {
		return []interface{}{}
	}

	mDesc := map[string]interface{}{
		"kinesis_stream_arn": aws.StringValue(desc.KinesisStreamARN),
		"role_arn":           aws.StringValue(desc.RoleARN),
	}

	return []interface{}{mDesc}
}

func flattenDeliveryStream(d *schema.ResourceData, s *firehose.DeliveryStreamDescription) error {
	d.Set("version_id", s.VersionId)
	d.Set("arn", s.DeliveryStreamARN)
	d.Set("name", s.DeliveryStreamName)

	sseOptions := map[string]interface{}{
		"enabled":  false,
		"key_type": firehose.KeyTypeAwsOwnedCmk,
	}
	if s.DeliveryStreamEncryptionConfiguration != nil &&
		aws.StringValue(s.DeliveryStreamEncryptionConfiguration.Status) == firehose.DeliveryStreamEncryptionStatusEnabled {
		sseOptions["enabled"] = true

		if v := s.DeliveryStreamEncryptionConfiguration.KeyARN; v != nil {
			sseOptions["key_arn"] = aws.StringValue(v)
		}
		if v := s.DeliveryStreamEncryptionConfiguration.KeyType; v != nil {
			sseOptions["key_type"] = aws.StringValue(v)
		}
	}

	if err := d.Set("server_side_encryption", []map[string]interface{}{sseOptions}); err != nil {
		return fmt.Errorf("setting server_side_encryption: %s", err)
	}

	if s.Source != nil {
		if err := d.Set("kinesis_source_configuration", flattenSourceConfiguration(s.Source.KinesisStreamSourceDescription)); err != nil {
			return fmt.Errorf("setting kinesis_source_configuration: %s", err)
		}
	}

	if len(s.Destinations) > 0 {
		destination := s.Destinations[0]
		if destination.RedshiftDestinationDescription != nil {
			d.Set("destination", destinationTypeRedshift)
			configuredPassword := d.Get("redshift_configuration.0.password").(string)
			if err := d.Set("redshift_configuration", flattenRedshiftConfiguration(destination.RedshiftDestinationDescription, configuredPassword)); err != nil {
				return fmt.Errorf("setting redshift_configuration: %s", err)
			}
		} else if destination.ElasticsearchDestinationDescription != nil {
			d.Set("destination", destinationTypeElasticsearch)
			if err := d.Set("elasticsearch_configuration", flattenElasticsearchConfiguration(destination.ElasticsearchDestinationDescription)); err != nil {
				return fmt.Errorf("setting elasticsearch_configuration: %s", err)
			}
		} else if destination.AmazonopensearchserviceDestinationDescription != nil {
			d.Set("destination", destinationTypeOpensearch)
			if err := d.Set("opensearch_configuration", flattenOpensearchConfiguration(destination.AmazonopensearchserviceDestinationDescription)); err != nil {
				return fmt.Errorf("setting opensearch_configuration: %s", err)
			}
		} else if destination.SplunkDestinationDescription != nil {
			d.Set("destination", destinationTypeSplunk)
			if err := d.Set("splunk_configuration", flattenSplunkConfiguration(destination.SplunkDestinationDescription)); err != nil {
				return fmt.Errorf("setting splunk_configuration: %s", err)
			}
		} else if destination.HttpEndpointDestinationDescription != nil {
			d.Set("destination", destinationTypeHTTPEndpoint)
			configuredAccessKey := d.Get("http_endpoint_configuration.0.access_key").(string)
			if err := d.Set("http_endpoint_configuration", flattenHTTPEndpointConfiguration(destination.HttpEndpointDestinationDescription, configuredAccessKey)); err != nil {
				return fmt.Errorf("setting http_endpoint_configuration: %s", err)
			}
		} else {
			d.Set("destination", destinationTypeExtendedS3)
			if err := d.Set("extended_s3_configuration", flattenExtendedS3Configuration(destination.ExtendedS3DestinationDescription)); err != nil {
				return fmt.Errorf("setting extended_s3_configuration: %s", err)
			}
		}
		d.Set("destination_id", destination.DestinationId)
	}

	return nil
}

func flattenHTTPEndpointConfiguration(description *firehose.HttpEndpointDestinationDescription, configuredAccessKey string) []map[string]interface{} {
	if description == nil {
		return []map[string]interface{}{}
	}
	m := map[string]interface{}{
		"access_key":                 configuredAccessKey,
		"url":                        aws.StringValue(description.EndpointConfiguration.Url),
		"name":                       aws.StringValue(description.EndpointConfiguration.Name),
		"role_arn":                   aws.StringValue(description.RoleARN),
		"s3_backup_mode":             aws.StringValue(description.S3BackupMode),
		"s3_configuration":           flattenS3Configuration(description.S3DestinationDescription),
		"request_configuration":      flattenRequestConfiguration(description.RequestConfiguration),
		"cloudwatch_logging_options": flattenCloudWatchLoggingOptions(description.CloudWatchLoggingOptions),
		"processing_configuration":   flattenProcessingConfiguration(description.ProcessingConfiguration, aws.StringValue(description.RoleARN)),
	}

	if description.RetryOptions != nil {
		m["retry_duration"] = int(aws.Int64Value(description.RetryOptions.DurationInSeconds))
	}

	if description.BufferingHints != nil {
		m["buffering_interval"] = int(aws.Int64Value(description.BufferingHints.IntervalInSeconds))
		m["buffering_size"] = int(aws.Int64Value(description.BufferingHints.SizeInMBs))
	}

	return []map[string]interface{}{m}
}

// @SDKResource("aws_kinesis_firehose_delivery_stream", name="Delivery Stream")
// @Tags(identifierAttribute="name")
func ResourceDeliveryStream() *schema.Resource {
	//lintignore:R011
	return &schema.Resource{
		CreateWithoutTimeout: resourceDeliveryStreamCreate,
		ReadWithoutTimeout:   resourceDeliveryStreamRead,
		UpdateWithoutTimeout: resourceDeliveryStreamUpdate,
		DeleteWithoutTimeout: resourceDeliveryStreamDelete,

		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idErr := fmt.Errorf("Expected ID in format of arn:PARTITION:firehose:REGION:ACCOUNTID:deliverystream/NAME and provided: %s", d.Id())
				resARN, err := arn.Parse(d.Id())
				if err != nil {
					return nil, idErr
				}
				resourceParts := strings.Split(resARN.Resource, "/")
				if len(resourceParts) != 2 {
					return nil, idErr
				}
				d.Set("name", resourceParts[1])
				return []*schema.ResourceData{d}, nil
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: verify.SetTagsDiff,

		SchemaVersion: 1,
		MigrateState:  MigrateState,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},

			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),

			"server_side_encryption": {
				Type:             schema.TypeList,
				Optional:         true,
				MaxItems:         1,
				DiffSuppressFunc: verify.SuppressMissingOptionalConfigurationBlock,
				ConflictsWith:    []string{"kinesis_source_configuration"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},

						"key_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      firehose.KeyTypeAwsOwnedCmk,
							ValidateFunc: validation.StringInSlice(firehose.KeyType_Values(), false),
							RequiredWith: []string{"server_side_encryption.0.enabled"},
						},

						"key_arn": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: verify.ValidARN,
							RequiredWith: []string{"server_side_encryption.0.enabled", "server_side_encryption.0.key_type"},
						},
					},
				},
			},

			"kinesis_source_configuration": {
				Type:          schema.TypeList,
				ForceNew:      true,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"server_side_encryption"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kinesis_stream_arn": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: verify.ValidARN,
						},

						"role_arn": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: verify.ValidARN,
						},
					},
				},
			},

			"destination": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				StateFunc: func(v interface{}) string {
					value := v.(string)
					return strings.ToLower(value)
				},
				ValidateFunc: validation.StringInSlice([]string{
					destinationTypeExtendedS3,
					destinationTypeRedshift,
					destinationTypeElasticsearch,
					destinationTypeOpensearch,
					destinationTypeSplunk,
					destinationTypeHTTPEndpoint,
				}, false),
			},

			"extended_s3_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket_arn": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: verify.ValidARN,
						},

						"buffering_size": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  5,
						},

						"buffering_interval": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  300,
						},

						"compression_format": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      firehose.CompressionFormatUncompressed,
							ValidateFunc: validation.StringInSlice(firehose.CompressionFormat_Values(), false),
						},

						"data_format_conversion_configuration": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"input_format_configuration": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"deserializer": {
													Type:     schema.TypeList,
													Required: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"hive_json_ser_de": {
																Type:          schema.TypeList,
																Optional:      true,
																MaxItems:      1,
																ConflictsWith: []string{"extended_s3_configuration.0.data_format_conversion_configuration.0.input_format_configuration.0.deserializer.0.open_x_json_ser_de"},
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"timestamp_formats": {
																			Type:     schema.TypeList,
																			Optional: true,
																			Elem:     &schema.Schema{Type: schema.TypeString},
																		},
																	},
																},
															},
															"open_x_json_ser_de": {
																Type:          schema.TypeList,
																Optional:      true,
																MaxItems:      1,
																ConflictsWith: []string{"extended_s3_configuration.0.data_format_conversion_configuration.0.input_format_configuration.0.deserializer.0.hive_json_ser_de"},
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"case_insensitive": {
																			Type:     schema.TypeBool,
																			Optional: true,
																			Default:  true,
																		},
																		"column_to_json_key_mappings": {
																			Type:     schema.TypeMap,
																			Optional: true,
																			Elem:     &schema.Schema{Type: schema.TypeString},
																		},
																		"convert_dots_in_json_keys_to_underscores": {
																			Type:     schema.TypeBool,
																			Optional: true,
																			Default:  false,
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
									"output_format_configuration": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"serializer": {
													Type:     schema.TypeList,
													Required: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"orc_ser_de": {
																Type:          schema.TypeList,
																Optional:      true,
																MaxItems:      1,
																ConflictsWith: []string{"extended_s3_configuration.0.data_format_conversion_configuration.0.output_format_configuration.0.serializer.0.parquet_ser_de"},
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"block_size_bytes": {
																			Type:     schema.TypeInt,
																			Optional: true,
																			// 256 MiB
																			Default: 268435456,
																			// 64 MiB
																			ValidateFunc: validation.IntAtLeast(67108864),
																		},
																		"bloom_filter_columns": {
																			Type:     schema.TypeList,
																			Optional: true,
																			Elem:     &schema.Schema{Type: schema.TypeString},
																		},
																		"bloom_filter_false_positive_probability": {
																			Type:     schema.TypeFloat,
																			Optional: true,
																			Default:  0.05,
																		},
																		"compression": {
																			Type:         schema.TypeString,
																			Optional:     true,
																			Default:      firehose.OrcCompressionSnappy,
																			ValidateFunc: validation.StringInSlice(firehose.OrcCompression_Values(), false),
																		},
																		"dictionary_key_threshold": {
																			Type:     schema.TypeFloat,
																			Optional: true,
																			Default:  0.0,
																		},
																		"enable_padding": {
																			Type:     schema.TypeBool,
																			Optional: true,
																			Default:  false,
																		},
																		"format_version": {
																			Type:         schema.TypeString,
																			Optional:     true,
																			Default:      firehose.OrcFormatVersionV012,
																			ValidateFunc: validation.StringInSlice(firehose.OrcFormatVersion_Values(), false),
																		},
																		"padding_tolerance": {
																			Type:     schema.TypeFloat,
																			Optional: true,
																			Default:  0.05,
																		},
																		"row_index_stride": {
																			Type:         schema.TypeInt,
																			Optional:     true,
																			Default:      10000,
																			ValidateFunc: validation.IntAtLeast(1000),
																		},
																		"stripe_size_bytes": {
																			Type:     schema.TypeInt,
																			Optional: true,
																			// 64 MiB
																			Default: 67108864,
																			// 8 MiB
																			ValidateFunc: validation.IntAtLeast(8388608),
																		},
																	},
																},
															},
															"parquet_ser_de": {
																Type:          schema.TypeList,
																Optional:      true,
																MaxItems:      1,
																ConflictsWith: []string{"extended_s3_configuration.0.data_format_conversion_configuration.0.output_format_configuration.0.serializer.0.orc_ser_de"},
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"block_size_bytes": {
																			Type:     schema.TypeInt,
																			Optional: true,
																			// 256 MiB
																			Default: 268435456,
																			// 64 MiB
																			ValidateFunc: validation.IntAtLeast(67108864),
																		},
																		"compression": {
																			Type:         schema.TypeString,
																			Optional:     true,
																			Default:      firehose.ParquetCompressionSnappy,
																			ValidateFunc: validation.StringInSlice(firehose.ParquetCompression_Values(), false),
																		},
																		"enable_dictionary_compression": {
																			Type:     schema.TypeBool,
																			Optional: true,
																			Default:  false,
																		},
																		"max_padding_bytes": {
																			Type:     schema.TypeInt,
																			Optional: true,
																			Default:  0,
																		},
																		"page_size_bytes": {
																			Type:     schema.TypeInt,
																			Optional: true,
																			// 1 MiB
																			Default: 1048576,
																			// 64 KiB
																			ValidateFunc: validation.IntAtLeast(65536),
																		},
																		"writer_version": {
																			Type:         schema.TypeString,
																			Optional:     true,
																			Default:      firehose.ParquetWriterVersionV1,
																			ValidateFunc: validation.StringInSlice(firehose.ParquetWriterVersion_Values(), false),
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
									"schema_configuration": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"catalog_id": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"database_name": {
													Type:     schema.TypeString,
													Required: true,
												},
												"region": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"role_arn": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: verify.ValidARN,
												},
												"table_name": {
													Type:     schema.TypeString,
													Required: true,
												},
												"version_id": {
													Type:     schema.TypeString,
													Optional: true,
													Default:  "LATEST",
												},
											},
										},
									},
								},
							},
						},

						"error_output_prefix": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(0, 1024),
						},

						"kms_key_arn": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: verify.ValidARN,
						},

						"role_arn": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: verify.ValidARN,
						},

						"prefix": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"s3_backup_mode": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      firehose.S3BackupModeDisabled,
							ValidateFunc: validation.StringInSlice(firehose.S3BackupMode_Values(), false),
						},

						"s3_backup_configuration": s3BackupConfigurationSchema(),

						"cloudwatch_logging_options": cloudWatchLoggingOptionsSchema(),

						"dynamic_partitioning_configuration": dynamicPartitioningConfigurationSchema(),

						"processing_configuration": processingConfigurationSchema(),
					},
				},
			},

			"redshift_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloudwatch_logging_options": cloudWatchLoggingOptionsSchema(),

						"cluster_jdbcurl": {
							Type:     schema.TypeString,
							Required: true,
						},

						"copy_options": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"data_table_columns": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"data_table_name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"password": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},

						"processing_configuration": processingConfigurationSchema(),

						"retry_duration": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      3600,
							ValidateFunc: validation.IntBetween(0, 7200),
						},

						"role_arn": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: verify.ValidARN,
						},

						"s3_backup_mode": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      firehose.S3BackupModeDisabled,
							ValidateFunc: validation.StringInSlice(firehose.S3BackupMode_Values(), false),
						},

						"s3_backup_configuration": s3BackupConfigurationSchema(),

						"s3_configuration": s3ConfigurationSchema(),

						"username": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"elasticsearch_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"buffering_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      300,
							ValidateFunc: validation.IntBetween(60, 900),
						},

						"buffering_size": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      5,
							ValidateFunc: validation.IntBetween(1, 100),
						},

						"cloudwatch_logging_options": cloudWatchLoggingOptionsSchema(),

						"cluster_endpoint": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"elasticsearch_configuration.0.domain_arn"},
						},

						"domain_arn": {
							Type:          schema.TypeString,
							Optional:      true,
							ValidateFunc:  verify.ValidARN,
							ConflictsWith: []string{"elasticsearch_configuration.0.cluster_endpoint"},
						},

						"index_name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"index_rotation_period": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      firehose.ElasticsearchIndexRotationPeriodOneDay,
							ValidateFunc: validation.StringInSlice(firehose.ElasticsearchIndexRotationPeriod_Values(), false),
						},

						"processing_configuration": processingConfigurationSchema(),

						"retry_duration": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      300,
							ValidateFunc: validation.IntBetween(0, 7200),
						},

						"role_arn": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: verify.ValidARN,
						},

						"s3_configuration": s3ConfigurationSchema(),

						"s3_backup_mode": {
							Type:         schema.TypeString,
							ForceNew:     true,
							Optional:     true,
							Default:      firehose.ElasticsearchS3BackupModeFailedDocumentsOnly,
							ValidateFunc: validation.StringInSlice(firehose.ElasticsearchS3BackupMode_Values(), false),
						},

						"type_name": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(0, 100),
						},

						"vpc_config": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"subnet_ids": {
										Type:     schema.TypeSet,
										Required: true,
										ForceNew: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"security_group_ids": {
										Type:     schema.TypeSet,
										Required: true,
										ForceNew: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"role_arn": {
										Type:         schema.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: verify.ValidARN,
									},
								},
							},
						},
					},
				},
			},

			"opensearch_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"buffering_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      300,
							ValidateFunc: validation.IntBetween(60, 900),
						},

						"buffering_size": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      5,
							ValidateFunc: validation.IntBetween(1, 100),
						},

						"domain_arn": {
							Type:          schema.TypeString,
							Optional:      true,
							ValidateFunc:  verify.ValidARN,
							ConflictsWith: []string{"opensearch_configuration.0.cluster_endpoint"},
						},

						"index_name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"index_rotation_period": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      firehose.AmazonopensearchserviceIndexRotationPeriodOneDay,
							ValidateFunc: validation.StringInSlice(firehose.AmazonopensearchserviceIndexRotationPeriod_Values(), false),
						},

						"retry_duration": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      300,
							ValidateFunc: validation.IntBetween(0, 7200),
						},

						"role_arn": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: verify.ValidARN,
						},

						"s3_configuration": s3ConfigurationSchema(),

						"s3_backup_mode": {
							Type:         schema.TypeString,
							ForceNew:     true,
							Optional:     true,
							Default:      firehose.AmazonopensearchserviceS3BackupModeFailedDocumentsOnly,
							ValidateFunc: validation.StringInSlice(firehose.AmazonopensearchserviceS3BackupMode_Values(), false),
						},

						"type_name": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(0, 100),
						},

						"vpc_config": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"subnet_ids": {
										Type:     schema.TypeSet,
										Required: true,
										ForceNew: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"security_group_ids": {
										Type:     schema.TypeSet,
										Required: true,
										ForceNew: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"role_arn": {
										Type:         schema.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: verify.ValidARN,
									},
								},
							},
						},

						"cloudwatch_logging_options": cloudWatchLoggingOptionsSchema(),

						"processing_configuration": processingConfigurationSchema(),
						"cluster_endpoint": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"opensearch_configuration.0.domain_arn"},
						},
					},
				},
			},

			"splunk_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hec_acknowledgment_timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      180,
							ValidateFunc: validation.IntBetween(180, 600),
						},

						"hec_endpoint": {
							Type:     schema.TypeString,
							Required: true,
						},

						"hec_endpoint_type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  firehose.HECEndpointTypeRaw,
							ValidateFunc: validation.StringInSlice([]string{
								firehose.HECEndpointTypeRaw,
								firehose.HECEndpointTypeEvent,
							}, false),
						},

						"hec_token": {
							Type:     schema.TypeString,
							Required: true,
						},

						"retry_duration": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      3600,
							ValidateFunc: validation.IntBetween(0, 7200),
						},

						"s3_backup_mode": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  firehose.SplunkS3BackupModeFailedEventsOnly,
							ValidateFunc: validation.StringInSlice([]string{
								firehose.SplunkS3BackupModeFailedEventsOnly,
								firehose.SplunkS3BackupModeAllEvents,
							}, false),
						},

						"s3_configuration": s3ConfigurationSchema(),

						"cloudwatch_logging_options": cloudWatchLoggingOptionsSchema(),

						"processing_configuration": processingConfigurationSchema(),
					},
				},
			},

			"http_endpoint_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.All(
								validation.StringLenBetween(1, 1000),
								validation.StringMatch(regexp.MustCompile(`^https://.*$`), ""),
							),
						},

						"name": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.All(
								validation.StringLenBetween(1, 256),
							),
						},

						"access_key": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(0, 4096),
							Sensitive:    true,
						},

						"retry_duration": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      300,
							ValidateFunc: validation.IntBetween(0, 7200),
						},

						"role_arn": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: verify.ValidARN,
						},

						"s3_backup_mode": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      firehose.HttpEndpointS3BackupModeFailedDataOnly,
							ValidateFunc: validation.StringInSlice(firehose.HttpEndpointS3BackupMode_Values(), false),
						},

						"s3_configuration": s3ConfigurationSchema(),

						"buffering_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      300,
							ValidateFunc: validation.IntBetween(60, 900),
						},

						"buffering_size": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      5,
							ValidateFunc: validation.IntBetween(1, 100),
						},

						"request_configuration": requestConfigurationSchema(),

						"cloudwatch_logging_options": cloudWatchLoggingOptionsSchema(),

						"processing_configuration": processingConfigurationSchema(),
					},
				},
			},

			"arn": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"version_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"destination_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func createSourceConfig(source map[string]interface{}) *firehose.KinesisStreamSourceConfiguration {
	configuration := &firehose.KinesisStreamSourceConfiguration{
		KinesisStreamARN: aws.String(source["kinesis_stream_arn"].(string)),
		RoleARN:          aws.String(source["role_arn"].(string)),
	}

	return configuration
}

func createS3Config(tfList []interface{}) *firehose.S3DestinationConfiguration {
	s3 := tfList[0].(map[string]interface{})

	configuration := &firehose.S3DestinationConfiguration{
		BucketARN: aws.String(s3["bucket_arn"].(string)),
		RoleARN:   aws.String(s3["role_arn"].(string)),
		BufferingHints: &firehose.BufferingHints{
			IntervalInSeconds: aws.Int64(int64(s3["buffering_interval"].(int))),
			SizeInMBs:         aws.Int64(int64(s3["buffering_size"].(int))),
		},
		Prefix:                  extractPrefixConfiguration(s3),
		CompressionFormat:       aws.String(s3["compression_format"].(string)),
		EncryptionConfiguration: extractEncryptionConfiguration(s3),
	}

	if v, ok := s3["error_output_prefix"].(string); ok && v != "" {
		configuration.ErrorOutputPrefix = aws.String(v)
	}

	if _, ok := s3["cloudwatch_logging_options"]; ok {
		configuration.CloudWatchLoggingOptions = extractCloudWatchLoggingConfiguration(s3)
	}

	return configuration
}

func expandS3BackupConfig(d map[string]interface{}) *firehose.S3DestinationConfiguration {
	config := d["s3_backup_configuration"].([]interface{})
	if len(config) == 0 {
		return nil
	}

	s3 := config[0].(map[string]interface{})

	configuration := &firehose.S3DestinationConfiguration{
		BucketARN: aws.String(s3["bucket_arn"].(string)),
		RoleARN:   aws.String(s3["role_arn"].(string)),
		BufferingHints: &firehose.BufferingHints{
			IntervalInSeconds: aws.Int64(int64(s3["buffering_interval"].(int))),
			SizeInMBs:         aws.Int64(int64(s3["buffering_size"].(int))),
		},
		Prefix:                  extractPrefixConfiguration(s3),
		CompressionFormat:       aws.String(s3["compression_format"].(string)),
		EncryptionConfiguration: extractEncryptionConfiguration(s3),
	}

	if _, ok := s3["cloudwatch_logging_options"]; ok {
		configuration.CloudWatchLoggingOptions = extractCloudWatchLoggingConfiguration(s3)
	}

	if v, ok := s3["error_output_prefix"].(string); ok && v != "" {
		configuration.ErrorOutputPrefix = aws.String(v)
	}

	return configuration
}

func createExtendedS3Config(d *schema.ResourceData) *firehose.ExtendedS3DestinationConfiguration {
	s3 := d.Get("extended_s3_configuration").([]interface{})[0].(map[string]interface{})

	configuration := &firehose.ExtendedS3DestinationConfiguration{
		BucketARN: aws.String(s3["bucket_arn"].(string)),
		RoleARN:   aws.String(s3["role_arn"].(string)),
		BufferingHints: &firehose.BufferingHints{
			IntervalInSeconds: aws.Int64(int64(s3["buffering_interval"].(int))),
			SizeInMBs:         aws.Int64(int64(s3["buffering_size"].(int))),
		},
		Prefix:                            extractPrefixConfiguration(s3),
		CompressionFormat:                 aws.String(s3["compression_format"].(string)),
		DataFormatConversionConfiguration: expandDataFormatConversionConfiguration(s3["data_format_conversion_configuration"].([]interface{})),
		EncryptionConfiguration:           extractEncryptionConfiguration(s3),
	}

	if _, ok := s3["processing_configuration"]; ok {
		configuration.ProcessingConfiguration = extractProcessingConfiguration(s3)
	}

	if _, ok := s3["dynamic_partitioning_configuration"]; ok {
		configuration.DynamicPartitioningConfiguration = extractDynamicPartitioningConfiguration(s3)
	}

	if _, ok := s3["cloudwatch_logging_options"]; ok {
		configuration.CloudWatchLoggingOptions = extractCloudWatchLoggingConfiguration(s3)
	}

	if v, ok := s3["error_output_prefix"].(string); ok && v != "" {
		configuration.ErrorOutputPrefix = aws.String(v)
	}

	if s3BackupMode, ok := s3["s3_backup_mode"]; ok {
		configuration.S3BackupMode = aws.String(s3BackupMode.(string))
		configuration.S3BackupConfiguration = expandS3BackupConfig(d.Get("extended_s3_configuration").([]interface{})[0].(map[string]interface{}))
	}

	return configuration
}

func updateS3Config(tfList []interface{}) *firehose.S3DestinationUpdate {
	s3 := tfList[0].(map[string]interface{})

	configuration := &firehose.S3DestinationUpdate{
		BucketARN: aws.String(s3["bucket_arn"].(string)),
		RoleARN:   aws.String(s3["role_arn"].(string)),
		BufferingHints: &firehose.BufferingHints{
			IntervalInSeconds: aws.Int64((int64)(s3["buffering_interval"].(int))),
			SizeInMBs:         aws.Int64((int64)(s3["buffering_size"].(int))),
		},
		ErrorOutputPrefix:        aws.String(s3["error_output_prefix"].(string)),
		Prefix:                   extractPrefixConfiguration(s3),
		CompressionFormat:        aws.String(s3["compression_format"].(string)),
		EncryptionConfiguration:  extractEncryptionConfiguration(s3),
		CloudWatchLoggingOptions: extractCloudWatchLoggingConfiguration(s3),
	}

	if _, ok := s3["cloudwatch_logging_options"]; ok {
		configuration.CloudWatchLoggingOptions = extractCloudWatchLoggingConfiguration(s3)
	}

	return configuration
}

func updateS3BackupConfig(d map[string]interface{}) *firehose.S3DestinationUpdate {
	config := d["s3_backup_configuration"].([]interface{})
	if len(config) == 0 {
		return nil
	}

	s3 := config[0].(map[string]interface{})

	configuration := &firehose.S3DestinationUpdate{
		BucketARN: aws.String(s3["bucket_arn"].(string)),
		RoleARN:   aws.String(s3["role_arn"].(string)),
		BufferingHints: &firehose.BufferingHints{
			IntervalInSeconds: aws.Int64((int64)(s3["buffering_interval"].(int))),
			SizeInMBs:         aws.Int64((int64)(s3["buffering_size"].(int))),
		},
		ErrorOutputPrefix:        aws.String(s3["error_output_prefix"].(string)),
		Prefix:                   extractPrefixConfiguration(s3),
		CompressionFormat:        aws.String(s3["compression_format"].(string)),
		EncryptionConfiguration:  extractEncryptionConfiguration(s3),
		CloudWatchLoggingOptions: extractCloudWatchLoggingConfiguration(s3),
	}

	if _, ok := s3["cloudwatch_logging_options"]; ok {
		configuration.CloudWatchLoggingOptions = extractCloudWatchLoggingConfiguration(s3)
	}

	return configuration
}

func updateExtendedS3Config(d *schema.ResourceData) *firehose.ExtendedS3DestinationUpdate {
	s3 := d.Get("extended_s3_configuration").([]interface{})[0].(map[string]interface{})

	configuration := &firehose.ExtendedS3DestinationUpdate{
		BucketARN: aws.String(s3["bucket_arn"].(string)),
		RoleARN:   aws.String(s3["role_arn"].(string)),
		BufferingHints: &firehose.BufferingHints{
			IntervalInSeconds: aws.Int64((int64)(s3["buffering_interval"].(int))),
			SizeInMBs:         aws.Int64((int64)(s3["buffering_size"].(int))),
		},
		ErrorOutputPrefix:                 aws.String(s3["error_output_prefix"].(string)),
		Prefix:                            extractPrefixConfiguration(s3),
		CompressionFormat:                 aws.String(s3["compression_format"].(string)),
		EncryptionConfiguration:           extractEncryptionConfiguration(s3),
		DataFormatConversionConfiguration: expandDataFormatConversionConfiguration(s3["data_format_conversion_configuration"].([]interface{})),
		CloudWatchLoggingOptions:          extractCloudWatchLoggingConfiguration(s3),
		ProcessingConfiguration:           extractProcessingConfiguration(s3),
	}

	if _, ok := s3["cloudwatch_logging_options"]; ok {
		configuration.CloudWatchLoggingOptions = extractCloudWatchLoggingConfiguration(s3)
	}

	if _, ok := s3["dynamic_partitioning_configuration"]; ok {
		configuration.DynamicPartitioningConfiguration = extractDynamicPartitioningConfiguration(s3)
	}

	if s3BackupMode, ok := s3["s3_backup_mode"]; ok {
		configuration.S3BackupMode = aws.String(s3BackupMode.(string))
		configuration.S3BackupUpdate = updateS3BackupConfig(d.Get("extended_s3_configuration").([]interface{})[0].(map[string]interface{}))
	}

	return configuration
}

func expandDataFormatConversionConfiguration(l []interface{}) *firehose.DataFormatConversionConfiguration {
	if len(l) == 0 || l[0] == nil {
		// It is possible to just pass nil here, but this seems to be the
		// canonical form that AWS uses, and is less likely to produce diffs.
		return &firehose.DataFormatConversionConfiguration{
			Enabled: aws.Bool(false),
		}
	}

	m := l[0].(map[string]interface{})

	return &firehose.DataFormatConversionConfiguration{
		Enabled:                   aws.Bool(m["enabled"].(bool)),
		InputFormatConfiguration:  expandInputFormatConfiguration(m["input_format_configuration"].([]interface{})),
		OutputFormatConfiguration: expandOutputFormatConfiguration(m["output_format_configuration"].([]interface{})),
		SchemaConfiguration:       expandSchemaConfiguration(m["schema_configuration"].([]interface{})),
	}
}

func expandInputFormatConfiguration(l []interface{}) *firehose.InputFormatConfiguration {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})

	return &firehose.InputFormatConfiguration{
		Deserializer: expandDeserializer(m["deserializer"].([]interface{})),
	}
}

func expandDeserializer(l []interface{}) *firehose.Deserializer {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})

	return &firehose.Deserializer{
		HiveJsonSerDe:  expandHiveJSONSerDe(m["hive_json_ser_de"].([]interface{})),
		OpenXJsonSerDe: expandOpenXJSONSerDe(m["open_x_json_ser_de"].([]interface{})),
	}
}

func expandHiveJSONSerDe(l []interface{}) *firehose.HiveJsonSerDe {
	if len(l) == 0 {
		return nil
	}

	if l[0] == nil {
		return &firehose.HiveJsonSerDe{}
	}

	m := l[0].(map[string]interface{})

	return &firehose.HiveJsonSerDe{
		TimestampFormats: flex.ExpandStringList(m["timestamp_formats"].([]interface{})),
	}
}

func expandOpenXJSONSerDe(l []interface{}) *firehose.OpenXJsonSerDe {
	if len(l) == 0 {
		return nil
	}

	if l[0] == nil {
		return &firehose.OpenXJsonSerDe{}
	}

	m := l[0].(map[string]interface{})

	return &firehose.OpenXJsonSerDe{
		CaseInsensitive:                    aws.Bool(m["case_insensitive"].(bool)),
		ColumnToJsonKeyMappings:            flex.ExpandStringMap(m["column_to_json_key_mappings"].(map[string]interface{})),
		ConvertDotsInJsonKeysToUnderscores: aws.Bool(m["convert_dots_in_json_keys_to_underscores"].(bool)),
	}
}

func expandOutputFormatConfiguration(l []interface{}) *firehose.OutputFormatConfiguration {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})

	return &firehose.OutputFormatConfiguration{
		Serializer: expandSerializer(m["serializer"].([]interface{})),
	}
}

func expandSerializer(l []interface{}) *firehose.Serializer {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})

	return &firehose.Serializer{
		OrcSerDe:     expandOrcSerDe(m["orc_ser_de"].([]interface{})),
		ParquetSerDe: expandParquetSerDe(m["parquet_ser_de"].([]interface{})),
	}
}

func expandOrcSerDe(l []interface{}) *firehose.OrcSerDe {
	if len(l) == 0 {
		return nil
	}

	if l[0] == nil {
		return &firehose.OrcSerDe{}
	}

	m := l[0].(map[string]interface{})

	orcSerDe := &firehose.OrcSerDe{
		BlockSizeBytes:                      aws.Int64(int64(m["block_size_bytes"].(int))),
		BloomFilterFalsePositiveProbability: aws.Float64(m["bloom_filter_false_positive_probability"].(float64)),
		Compression:                         aws.String(m["compression"].(string)),
		DictionaryKeyThreshold:              aws.Float64(m["dictionary_key_threshold"].(float64)),
		EnablePadding:                       aws.Bool(m["enable_padding"].(bool)),
		FormatVersion:                       aws.String(m["format_version"].(string)),
		PaddingTolerance:                    aws.Float64(m["padding_tolerance"].(float64)),
		RowIndexStride:                      aws.Int64(int64(m["row_index_stride"].(int))),
		StripeSizeBytes:                     aws.Int64(int64(m["stripe_size_bytes"].(int))),
	}

	if v, ok := m["bloom_filter_columns"].([]interface{}); ok && len(v) > 0 {
		orcSerDe.BloomFilterColumns = flex.ExpandStringList(v)
	}

	return orcSerDe
}

func expandParquetSerDe(l []interface{}) *firehose.ParquetSerDe {
	if len(l) == 0 {
		return nil
	}

	if l[0] == nil {
		return &firehose.ParquetSerDe{}
	}

	m := l[0].(map[string]interface{})

	return &firehose.ParquetSerDe{
		BlockSizeBytes:              aws.Int64(int64(m["block_size_bytes"].(int))),
		Compression:                 aws.String(m["compression"].(string)),
		EnableDictionaryCompression: aws.Bool(m["enable_dictionary_compression"].(bool)),
		MaxPaddingBytes:             aws.Int64(int64(m["max_padding_bytes"].(int))),
		PageSizeBytes:               aws.Int64(int64(m["page_size_bytes"].(int))),
		WriterVersion:               aws.String(m["writer_version"].(string)),
	}
}

func expandSchemaConfiguration(l []interface{}) *firehose.SchemaConfiguration {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})

	config := &firehose.SchemaConfiguration{
		DatabaseName: aws.String(m["database_name"].(string)),
		RoleARN:      aws.String(m["role_arn"].(string)),
		TableName:    aws.String(m["table_name"].(string)),
		VersionId:    aws.String(m["version_id"].(string)),
	}

	if v, ok := m["catalog_id"].(string); ok && v != "" {
		config.CatalogId = aws.String(v)
	}
	if v, ok := m["region"].(string); ok && v != "" {
		config.Region = aws.String(v)
	}

	return config
}

func extractDynamicPartitioningConfiguration(s3 map[string]interface{}) *firehose.DynamicPartitioningConfiguration {
	config := s3["dynamic_partitioning_configuration"].([]interface{})
	if len(config) == 0 {
		return nil
	}

	dynamicPartitioningConfig := config[0].(map[string]interface{})
	DynamicPartitioningConfiguration := &firehose.DynamicPartitioningConfiguration{
		Enabled: aws.Bool(dynamicPartitioningConfig["enabled"].(bool)),
	}

	if retryDuration, ok := dynamicPartitioningConfig["retry_duration"]; ok {
		DynamicPartitioningConfiguration.RetryOptions = &firehose.RetryOptions{
			DurationInSeconds: aws.Int64(int64(retryDuration.(int))),
		}
	}

	return DynamicPartitioningConfiguration
}

func extractProcessingConfiguration(s3 map[string]interface{}) *firehose.ProcessingConfiguration {
	config := s3["processing_configuration"].([]interface{})
	if len(config) == 0 || config[0] == nil {
		// It is possible to just pass nil here, but this seems to be the
		// canonical form that AWS uses, and is less likely to produce diffs.
		return &firehose.ProcessingConfiguration{
			Enabled:    aws.Bool(false),
			Processors: []*firehose.Processor{},
		}
	}

	processingConfiguration := config[0].(map[string]interface{})

	return &firehose.ProcessingConfiguration{
		Enabled:    aws.Bool(processingConfiguration["enabled"].(bool)),
		Processors: extractProcessors(processingConfiguration["processors"].([]interface{})),
	}
}

func extractProcessors(processingConfigurationProcessors []interface{}) []*firehose.Processor {
	processors := []*firehose.Processor{}

	for _, processor := range processingConfigurationProcessors {
		extractedProcessor := extractProcessor(processor.(map[string]interface{}))
		if extractedProcessor != nil {
			processors = append(processors, extractedProcessor)
		}
	}

	return processors
}

func extractProcessor(processingConfigurationProcessor map[string]interface{}) *firehose.Processor {
	var processor *firehose.Processor
	processorType := processingConfigurationProcessor["type"].(string)
	if processorType != "" {
		processor = &firehose.Processor{
			Type:       aws.String(processorType),
			Parameters: extractProcessorParameters(processingConfigurationProcessor["parameters"].([]interface{})),
		}
	}
	return processor
}

func extractProcessorParameters(processorParameters []interface{}) []*firehose.ProcessorParameter {
	parameters := []*firehose.ProcessorParameter{}

	for _, attr := range processorParameters {
		parameters = append(parameters, extractProcessorParameter(attr.(map[string]interface{})))
	}

	return parameters
}

func extractProcessorParameter(processorParameter map[string]interface{}) *firehose.ProcessorParameter {
	parameter := &firehose.ProcessorParameter{
		ParameterName:  aws.String(processorParameter["parameter_name"].(string)),
		ParameterValue: aws.String(processorParameter["parameter_value"].(string)),
	}

	return parameter
}

func extractEncryptionConfiguration(s3 map[string]interface{}) *firehose.EncryptionConfiguration {
	if key, ok := s3["kms_key_arn"]; ok && len(key.(string)) > 0 {
		return &firehose.EncryptionConfiguration{
			KMSEncryptionConfig: &firehose.KMSEncryptionConfig{
				AWSKMSKeyARN: aws.String(key.(string)),
			},
		}
	}

	return &firehose.EncryptionConfiguration{
		NoEncryptionConfig: aws.String(firehose.NoEncryptionConfigNoEncryption),
	}
}

func extractCloudWatchLoggingConfiguration(s3 map[string]interface{}) *firehose.CloudWatchLoggingOptions {
	config := s3["cloudwatch_logging_options"].([]interface{})
	if len(config) == 0 {
		return nil
	}

	loggingConfig := config[0].(map[string]interface{})
	loggingOptions := &firehose.CloudWatchLoggingOptions{
		Enabled: aws.Bool(loggingConfig["enabled"].(bool)),
	}

	if v, ok := loggingConfig["log_group_name"]; ok {
		loggingOptions.LogGroupName = aws.String(v.(string))
	}

	if v, ok := loggingConfig["log_stream_name"]; ok {
		loggingOptions.LogStreamName = aws.String(v.(string))
	}

	return loggingOptions
}

func extractVPCConfiguration(es map[string]interface{}) *firehose.VpcConfiguration {
	config := es["vpc_config"].([]interface{})
	if len(config) == 0 {
		return nil
	}

	vpcConfig := config[0].(map[string]interface{})

	return &firehose.VpcConfiguration{
		RoleARN:          aws.String(vpcConfig["role_arn"].(string)),
		SubnetIds:        flex.ExpandStringSet(vpcConfig["subnet_ids"].(*schema.Set)),
		SecurityGroupIds: flex.ExpandStringSet(vpcConfig["security_group_ids"].(*schema.Set)),
	}
}

func extractPrefixConfiguration(s3 map[string]interface{}) *string {
	if v, ok := s3["prefix"]; ok {
		return aws.String(v.(string))
	}

	return nil
}

func createRedshiftConfig(d *schema.ResourceData) (*firehose.RedshiftDestinationConfiguration, error) {
	redshiftRaw, ok := d.GetOk("redshift_configuration")
	if !ok {
		return nil, elasticsearchDestinationRequiredParamErr("redshift_configuration", destinationTypeRedshift)
	}
	rl := redshiftRaw.([]interface{})

	redshift := rl[0].(map[string]interface{})

	configuration := &firehose.RedshiftDestinationConfiguration{
		ClusterJDBCURL:  aws.String(redshift["cluster_jdbcurl"].(string)),
		RetryOptions:    extractRedshiftRetryOptions(redshift),
		Password:        aws.String(redshift["password"].(string)),
		Username:        aws.String(redshift["username"].(string)),
		RoleARN:         aws.String(redshift["role_arn"].(string)),
		CopyCommand:     extractCopyCommandConfiguration(redshift),
		S3Configuration: createS3Config(redshift["s3_configuration"].([]interface{})),
	}

	if _, ok := redshift["cloudwatch_logging_options"]; ok {
		configuration.CloudWatchLoggingOptions = extractCloudWatchLoggingConfiguration(redshift)
	}
	if _, ok := redshift["processing_configuration"]; ok {
		configuration.ProcessingConfiguration = extractProcessingConfiguration(redshift)
	}
	if s3BackupMode, ok := redshift["s3_backup_mode"]; ok {
		configuration.S3BackupMode = aws.String(s3BackupMode.(string))
		configuration.S3BackupConfiguration = expandS3BackupConfig(d.Get("redshift_configuration").([]interface{})[0].(map[string]interface{}))
	}

	return configuration, nil
}

func updateRedshiftConfig(d *schema.ResourceData) (*firehose.RedshiftDestinationUpdate, error) {
	redshiftRaw, ok := d.GetOk("redshift_configuration")
	if !ok {
		return nil, elasticsearchDestinationRequiredParamErr("redshift_configuration", destinationTypeRedshift)
	}
	rl := redshiftRaw.([]interface{})

	redshift := rl[0].(map[string]interface{})

	configuration := &firehose.RedshiftDestinationUpdate{
		ClusterJDBCURL: aws.String(redshift["cluster_jdbcurl"].(string)),
		RetryOptions:   extractRedshiftRetryOptions(redshift),
		Password:       aws.String(redshift["password"].(string)),
		Username:       aws.String(redshift["username"].(string)),
		RoleARN:        aws.String(redshift["role_arn"].(string)),
		CopyCommand:    extractCopyCommandConfiguration(redshift),
	}

	s3Config := updateS3Config(redshift["s3_configuration"].([]interface{}))
	// Redshift does not currently support ErrorOutputPrefix,
	// which is set to the empty string within "updateS3Config",
	// thus we must remove it here to avoid an InvalidArgumentException.
	s3Config.ErrorOutputPrefix = nil
	configuration.S3Update = s3Config

	if _, ok := redshift["cloudwatch_logging_options"]; ok {
		configuration.CloudWatchLoggingOptions = extractCloudWatchLoggingConfiguration(redshift)
	}
	if _, ok := redshift["processing_configuration"]; ok {
		configuration.ProcessingConfiguration = extractProcessingConfiguration(redshift)
	}
	if s3BackupMode, ok := redshift["s3_backup_mode"]; ok {
		configuration.S3BackupMode = aws.String(s3BackupMode.(string))
		configuration.S3BackupUpdate = updateS3BackupConfig(d.Get("redshift_configuration").([]interface{})[0].(map[string]interface{}))
		if configuration.S3BackupUpdate != nil {
			// Redshift does not currently support ErrorOutputPrefix,
			// which is set to the empty string within "updateS3BackupConfig",
			// thus we must remove it here to avoid an InvalidArgumentException.
			configuration.S3BackupUpdate.ErrorOutputPrefix = nil
		}
	}

	return configuration, nil
}

func elasticsearchDestinationRequiredParamErr(param, destination string) error {
	return fmt.Errorf(`%q is required when "destination" is %q`, param, destination)
}

func createElasticsearchConfig(d *schema.ResourceData) (*firehose.ElasticsearchDestinationConfiguration, error) {
	esConfig, ok := d.GetOk("elasticsearch_configuration")
	if !ok {
		return nil, elasticsearchDestinationRequiredParamErr("elasticsearch_configuration", destinationTypeElasticsearch)
	}
	esList := esConfig.([]interface{})

	es := esList[0].(map[string]interface{})

	config := &firehose.ElasticsearchDestinationConfiguration{
		BufferingHints:  extractElasticsearchBufferingHints(es),
		IndexName:       aws.String(es["index_name"].(string)),
		RetryOptions:    extractElasticsearchRetryOptions(es),
		RoleARN:         aws.String(es["role_arn"].(string)),
		TypeName:        aws.String(es["type_name"].(string)),
		S3Configuration: createS3Config(es["s3_configuration"].([]interface{})),
	}

	if v, ok := es["domain_arn"]; ok && v.(string) != "" {
		config.DomainARN = aws.String(v.(string))
	}

	if v, ok := es["cluster_endpoint"]; ok && v.(string) != "" {
		config.ClusterEndpoint = aws.String(v.(string))
	}

	if _, ok := es["cloudwatch_logging_options"]; ok {
		config.CloudWatchLoggingOptions = extractCloudWatchLoggingConfiguration(es)
	}

	if _, ok := es["processing_configuration"]; ok {
		config.ProcessingConfiguration = extractProcessingConfiguration(es)
	}

	if indexRotationPeriod, ok := es["index_rotation_period"]; ok {
		config.IndexRotationPeriod = aws.String(indexRotationPeriod.(string))
	}
	if s3BackupMode, ok := es["s3_backup_mode"]; ok {
		config.S3BackupMode = aws.String(s3BackupMode.(string))
	}

	if _, ok := es["vpc_config"]; ok {
		config.VpcConfiguration = extractVPCConfiguration(es)
	}

	return config, nil
}

func updateElasticsearchConfig(d *schema.ResourceData) (*firehose.ElasticsearchDestinationUpdate, error) {
	esConfig, ok := d.GetOk("elasticsearch_configuration")
	if !ok {
		return nil, elasticsearchDestinationRequiredParamErr("elasticsearch_configuration", destinationTypeElasticsearch)
	}
	esList := esConfig.([]interface{})

	es := esList[0].(map[string]interface{})

	update := &firehose.ElasticsearchDestinationUpdate{
		BufferingHints: extractElasticsearchBufferingHints(es),
		IndexName:      aws.String(es["index_name"].(string)),
		RetryOptions:   extractElasticsearchRetryOptions(es),
		RoleARN:        aws.String(es["role_arn"].(string)),
		TypeName:       aws.String(es["type_name"].(string)),
		S3Update:       updateS3Config(es["s3_configuration"].([]interface{})),
	}

	if v, ok := es["domain_arn"]; ok && v.(string) != "" {
		update.DomainARN = aws.String(v.(string))
	}

	if v, ok := es["cluster_endpoint"]; ok && v.(string) != "" {
		update.ClusterEndpoint = aws.String(v.(string))
	}

	if _, ok := es["cloudwatch_logging_options"]; ok {
		update.CloudWatchLoggingOptions = extractCloudWatchLoggingConfiguration(es)
	}

	if _, ok := es["processing_configuration"]; ok {
		update.ProcessingConfiguration = extractProcessingConfiguration(es)
	}

	if indexRotationPeriod, ok := es["index_rotation_period"]; ok {
		update.IndexRotationPeriod = aws.String(indexRotationPeriod.(string))
	}

	return update, nil
}

func createOpensearchConfig(d *schema.ResourceData) (*firehose.AmazonopensearchserviceDestinationConfiguration, error) {
	esConfig, ok := d.GetOk("opensearch_configuration")
	if !ok {
		return nil, elasticsearchDestinationRequiredParamErr("opensearch_configuration", destinationTypeOpensearch)
	}
	esList := esConfig.([]interface{})

	es := esList[0].(map[string]interface{})

	config := &firehose.AmazonopensearchserviceDestinationConfiguration{
		BufferingHints:  extractOpensearchBufferingHints(es),
		IndexName:       aws.String(es["index_name"].(string)),
		RetryOptions:    extractOpensearchRetryOptions(es),
		RoleARN:         aws.String(es["role_arn"].(string)),
		TypeName:        aws.String(es["type_name"].(string)),
		S3Configuration: createS3Config(es["s3_configuration"].([]interface{})),
	}

	if v, ok := es["domain_arn"]; ok && v.(string) != "" {
		config.DomainARN = aws.String(v.(string))
	}

	if v, ok := es["cluster_endpoint"]; ok && v.(string) != "" {
		config.ClusterEndpoint = aws.String(v.(string))
	}

	if _, ok := es["cloudwatch_logging_options"]; ok {
		config.CloudWatchLoggingOptions = extractCloudWatchLoggingConfiguration(es)
	}

	if _, ok := es["processing_configuration"]; ok {
		config.ProcessingConfiguration = extractProcessingConfiguration(es)
	}

	if indexRotationPeriod, ok := es["index_rotation_period"]; ok {
		config.IndexRotationPeriod = aws.String(indexRotationPeriod.(string))
	}
	if s3BackupMode, ok := es["s3_backup_mode"]; ok {
		config.S3BackupMode = aws.String(s3BackupMode.(string))
	}

	if _, ok := es["vpc_config"]; ok {
		config.VpcConfiguration = extractVPCConfiguration(es)
	}

	return config, nil
}

func updateOpensearchConfig(d *schema.ResourceData) (*firehose.AmazonopensearchserviceDestinationUpdate, error) {
	esConfig, ok := d.GetOk("opensearch_configuration")
	if !ok {
		return nil, elasticsearchDestinationRequiredParamErr("elasticsearch_configuration", destinationTypeOpensearch)
	}
	esList := esConfig.([]interface{})

	es := esList[0].(map[string]interface{})

	update := &firehose.AmazonopensearchserviceDestinationUpdate{
		BufferingHints: extractOpensearchBufferingHints(es),
		IndexName:      aws.String(es["index_name"].(string)),
		RetryOptions:   extractOpensearchRetryOptions(es),
		RoleARN:        aws.String(es["role_arn"].(string)),
		TypeName:       aws.String(es["type_name"].(string)),
		S3Update:       updateS3Config(es["s3_configuration"].([]interface{})),
	}

	if v, ok := es["domain_arn"]; ok && v.(string) != "" {
		update.DomainARN = aws.String(v.(string))
	}

	if v, ok := es["cluster_endpoint"]; ok && v.(string) != "" {
		update.ClusterEndpoint = aws.String(v.(string))
	}

	if _, ok := es["cloudwatch_logging_options"]; ok {
		update.CloudWatchLoggingOptions = extractCloudWatchLoggingConfiguration(es)
	}

	if _, ok := es["processing_configuration"]; ok {
		update.ProcessingConfiguration = extractProcessingConfiguration(es)
	}

	if indexRotationPeriod, ok := es["index_rotation_period"]; ok {
		update.IndexRotationPeriod = aws.String(indexRotationPeriod.(string))
	}

	return update, nil
}

func createSplunkConfig(d *schema.ResourceData) (*firehose.SplunkDestinationConfiguration, error) {
	splunkRaw, ok := d.GetOk("splunk_configuration")
	if !ok {
		return nil, elasticsearchDestinationRequiredParamErr("splunk_configuration", destinationTypeSplunk)
	}
	sl := splunkRaw.([]interface{})

	splunk := sl[0].(map[string]interface{})

	configuration := &firehose.SplunkDestinationConfiguration{
		HECToken:                          aws.String(splunk["hec_token"].(string)),
		HECEndpointType:                   aws.String(splunk["hec_endpoint_type"].(string)),
		HECEndpoint:                       aws.String(splunk["hec_endpoint"].(string)),
		HECAcknowledgmentTimeoutInSeconds: aws.Int64(int64(splunk["hec_acknowledgment_timeout"].(int))),
		RetryOptions:                      extractSplunkRetryOptions(splunk),
		S3Configuration:                   createS3Config(splunk["s3_configuration"].([]interface{})),
	}

	if _, ok := splunk["processing_configuration"]; ok {
		configuration.ProcessingConfiguration = extractProcessingConfiguration(splunk)
	}

	if _, ok := splunk["cloudwatch_logging_options"]; ok {
		configuration.CloudWatchLoggingOptions = extractCloudWatchLoggingConfiguration(splunk)
	}
	if s3BackupMode, ok := splunk["s3_backup_mode"]; ok {
		configuration.S3BackupMode = aws.String(s3BackupMode.(string))
	}

	return configuration, nil
}

func updateSplunkConfig(d *schema.ResourceData) (*firehose.SplunkDestinationUpdate, error) {
	splunkRaw, ok := d.GetOk("splunk_configuration")
	if !ok {
		return nil, elasticsearchDestinationRequiredParamErr("splunk_configuration", destinationTypeSplunk)
	}
	sl := splunkRaw.([]interface{})

	splunk := sl[0].(map[string]interface{})

	configuration := &firehose.SplunkDestinationUpdate{
		HECToken:                          aws.String(splunk["hec_token"].(string)),
		HECEndpointType:                   aws.String(splunk["hec_endpoint_type"].(string)),
		HECEndpoint:                       aws.String(splunk["hec_endpoint"].(string)),
		HECAcknowledgmentTimeoutInSeconds: aws.Int64(int64(splunk["hec_acknowledgment_timeout"].(int))),
		RetryOptions:                      extractSplunkRetryOptions(splunk),
		S3Update:                          updateS3Config(splunk["s3_configuration"].([]interface{})),
	}

	if _, ok := splunk["processing_configuration"]; ok {
		configuration.ProcessingConfiguration = extractProcessingConfiguration(splunk)
	}

	if _, ok := splunk["cloudwatch_logging_options"]; ok {
		configuration.CloudWatchLoggingOptions = extractCloudWatchLoggingConfiguration(splunk)
	}
	if s3BackupMode, ok := splunk["s3_backup_mode"]; ok {
		configuration.S3BackupMode = aws.String(s3BackupMode.(string))
	}

	return configuration, nil
}

func createHTTPEndpointConfig(d *schema.ResourceData) (*firehose.HttpEndpointDestinationConfiguration, error) {
	HttpEndpointRaw, ok := d.GetOk("http_endpoint_configuration")
	if !ok {
		return nil, elasticsearchDestinationRequiredParamErr("http_endpoint_configuration", destinationTypeHTTPEndpoint)
	}
	sl := HttpEndpointRaw.([]interface{})

	HttpEndpoint := sl[0].(map[string]interface{})

	configuration := &firehose.HttpEndpointDestinationConfiguration{
		RetryOptions:    extractHTTPEndpointRetryOptions(HttpEndpoint),
		RoleARN:         aws.String(HttpEndpoint["role_arn"].(string)),
		S3Configuration: createS3Config(HttpEndpoint["s3_configuration"].([]interface{})),
	}

	configuration.EndpointConfiguration = extractHTTPEndpointConfiguration(HttpEndpoint)

	bufferingHints := &firehose.HttpEndpointBufferingHints{}

	if bufferingInterval, ok := HttpEndpoint["buffering_interval"].(int); ok {
		bufferingHints.IntervalInSeconds = aws.Int64(int64(bufferingInterval))
	}
	if bufferingSize, ok := HttpEndpoint["buffering_size"].(int); ok {
		bufferingHints.SizeInMBs = aws.Int64(int64(bufferingSize))
	}
	configuration.BufferingHints = bufferingHints

	if _, ok := HttpEndpoint["processing_configuration"]; ok {
		configuration.ProcessingConfiguration = extractProcessingConfiguration(HttpEndpoint)
	}

	if _, ok := HttpEndpoint["request_configuration"]; ok {
		configuration.RequestConfiguration = extractRequestConfiguration(HttpEndpoint)
	}

	if _, ok := HttpEndpoint["cloudwatch_logging_options"]; ok {
		configuration.CloudWatchLoggingOptions = extractCloudWatchLoggingConfiguration(HttpEndpoint)
	}
	if s3BackupMode, ok := HttpEndpoint["s3_backup_mode"]; ok {
		configuration.S3BackupMode = aws.String(s3BackupMode.(string))
	}

	return configuration, nil
}

func updateHTTPEndpointConfig(d *schema.ResourceData) (*firehose.HttpEndpointDestinationUpdate, error) {
	HttpEndpointRaw, ok := d.GetOk("http_endpoint_configuration")
	if !ok {
		return nil, elasticsearchDestinationRequiredParamErr("http_endpoint_configuration", destinationTypeHTTPEndpoint)
	}
	sl := HttpEndpointRaw.([]interface{})

	HttpEndpoint := sl[0].(map[string]interface{})

	configuration := &firehose.HttpEndpointDestinationUpdate{
		RetryOptions: extractHTTPEndpointRetryOptions(HttpEndpoint),
		RoleARN:      aws.String(HttpEndpoint["role_arn"].(string)),
		S3Update:     updateS3Config(HttpEndpoint["s3_configuration"].([]interface{})),
	}

	configuration.EndpointConfiguration = extractHTTPEndpointConfiguration(HttpEndpoint)

	bufferingHints := &firehose.HttpEndpointBufferingHints{}

	if bufferingInterval, ok := HttpEndpoint["buffering_interval"].(int); ok {
		bufferingHints.IntervalInSeconds = aws.Int64(int64(bufferingInterval))
	}
	if bufferingSize, ok := HttpEndpoint["buffering_size"].(int); ok {
		bufferingHints.SizeInMBs = aws.Int64(int64(bufferingSize))
	}
	configuration.BufferingHints = bufferingHints

	if _, ok := HttpEndpoint["processing_configuration"]; ok {
		configuration.ProcessingConfiguration = extractProcessingConfiguration(HttpEndpoint)
	}

	if _, ok := HttpEndpoint["request_configuration"]; ok {
		configuration.RequestConfiguration = extractRequestConfiguration(HttpEndpoint)
	}

	if _, ok := HttpEndpoint["cloudwatch_logging_options"]; ok {
		configuration.CloudWatchLoggingOptions = extractCloudWatchLoggingConfiguration(HttpEndpoint)
	}

	if s3BackupMode, ok := HttpEndpoint["s3_backup_mode"]; ok {
		configuration.S3BackupMode = aws.String(s3BackupMode.(string))
	}

	return configuration, nil
}

func extractCommonAttributes(ca []interface{}) []*firehose.HttpEndpointCommonAttribute {
	CommonAttributes := make([]*firehose.HttpEndpointCommonAttribute, 0, len(ca))

	for _, raw := range ca {
		data := raw.(map[string]interface{})

		a := &firehose.HttpEndpointCommonAttribute{
			AttributeName:  aws.String(data["name"].(string)),
			AttributeValue: aws.String(data["value"].(string)),
		}
		CommonAttributes = append(CommonAttributes, a)
	}

	return CommonAttributes
}

func extractRequestConfiguration(rc map[string]interface{}) *firehose.HttpEndpointRequestConfiguration {
	config := rc["request_configuration"].([]interface{})
	if len(config) == 0 {
		return nil
	}

	requestConfig := config[0].(map[string]interface{})
	RequestConfiguration := &firehose.HttpEndpointRequestConfiguration{}

	if contentEncoding, ok := requestConfig["content_encoding"]; ok {
		RequestConfiguration.ContentEncoding = aws.String(contentEncoding.(string))
	}

	if commonAttributes, ok := requestConfig["common_attributes"]; ok {
		RequestConfiguration.CommonAttributes = extractCommonAttributes(commonAttributes.([]interface{}))
	}

	return RequestConfiguration
}

func extractHTTPEndpointConfiguration(ep map[string]interface{}) *firehose.HttpEndpointConfiguration {
	endpointConfiguration := &firehose.HttpEndpointConfiguration{
		Url: aws.String(ep["url"].(string)),
	}

	if Name, ok := ep["name"]; ok {
		endpointConfiguration.Name = aws.String(Name.(string))
	}

	if AccessKey, ok := ep["access_key"]; ok {
		endpointConfiguration.AccessKey = aws.String(AccessKey.(string))
	}

	return endpointConfiguration
}

func extractElasticsearchBufferingHints(es map[string]interface{}) *firehose.ElasticsearchBufferingHints {
	bufferingHints := &firehose.ElasticsearchBufferingHints{}

	if bufferingInterval, ok := es["buffering_interval"].(int); ok {
		bufferingHints.IntervalInSeconds = aws.Int64(int64(bufferingInterval))
	}
	if bufferingSize, ok := es["buffering_size"].(int); ok {
		bufferingHints.SizeInMBs = aws.Int64(int64(bufferingSize))
	}

	return bufferingHints
}

func extractOpensearchBufferingHints(es map[string]interface{}) *firehose.AmazonopensearchserviceBufferingHints {
	bufferingHints := &firehose.AmazonopensearchserviceBufferingHints{}

	if bufferingInterval, ok := es["buffering_interval"].(int); ok {
		bufferingHints.IntervalInSeconds = aws.Int64(int64(bufferingInterval))
	}
	if bufferingSize, ok := es["buffering_size"].(int); ok {
		bufferingHints.SizeInMBs = aws.Int64(int64(bufferingSize))
	}

	return bufferingHints
}

func extractElasticsearchRetryOptions(es map[string]interface{}) *firehose.ElasticsearchRetryOptions {
	retryOptions := &firehose.ElasticsearchRetryOptions{}

	if retryDuration, ok := es["retry_duration"].(int); ok {
		retryOptions.DurationInSeconds = aws.Int64(int64(retryDuration))
	}

	return retryOptions
}

func extractOpensearchRetryOptions(es map[string]interface{}) *firehose.AmazonopensearchserviceRetryOptions {
	retryOptions := &firehose.AmazonopensearchserviceRetryOptions{}

	if retryDuration, ok := es["retry_duration"].(int); ok {
		retryOptions.DurationInSeconds = aws.Int64(int64(retryDuration))
	}

	return retryOptions
}

func extractHTTPEndpointRetryOptions(tfMap map[string]interface{}) *firehose.HttpEndpointRetryOptions {
	retryOptions := &firehose.HttpEndpointRetryOptions{}

	if retryDuration, ok := tfMap["retry_duration"].(int); ok {
		retryOptions.DurationInSeconds = aws.Int64(int64(retryDuration))
	}

	return retryOptions
}

func extractRedshiftRetryOptions(redshift map[string]interface{}) *firehose.RedshiftRetryOptions {
	retryOptions := &firehose.RedshiftRetryOptions{}

	if retryDuration, ok := redshift["retry_duration"].(int); ok {
		retryOptions.DurationInSeconds = aws.Int64(int64(retryDuration))
	}

	return retryOptions
}

func extractSplunkRetryOptions(splunk map[string]interface{}) *firehose.SplunkRetryOptions {
	retryOptions := &firehose.SplunkRetryOptions{}

	if retryDuration, ok := splunk["retry_duration"].(int); ok {
		retryOptions.DurationInSeconds = aws.Int64(int64(retryDuration))
	}

	return retryOptions
}

func extractCopyCommandConfiguration(redshift map[string]interface{}) *firehose.CopyCommand {
	cmd := &firehose.CopyCommand{
		DataTableName: aws.String(redshift["data_table_name"].(string)),
	}
	if copyOptions, ok := redshift["copy_options"]; ok {
		cmd.CopyOptions = aws.String(copyOptions.(string))
	}
	if columns, ok := redshift["data_table_columns"]; ok {
		cmd.DataTableColumns = aws.String(columns.(string))
	}

	return cmd
}

func resourceDeliveryStreamCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	sn := d.Get("name").(string)

	if err := validateSchema(d); err != nil {
		return sdkdiag.AppendErrorf(diags, "creating Kinesis Firehose Delivery Stream (%s): %s", sn, err)
	}

	conn := meta.(*conns.AWSClient).FirehoseConn(ctx)
	input := &firehose.CreateDeliveryStreamInput{
		DeliveryStreamName: aws.String(sn),
		Tags:               GetTagsIn(ctx),
	}

	if v, ok := d.GetOk("kinesis_source_configuration"); ok {
		sourceConfig := createSourceConfig(v.([]interface{})[0].(map[string]interface{}))
		input.KinesisStreamSourceConfiguration = sourceConfig
		input.DeliveryStreamType = aws.String(firehose.DeliveryStreamTypeKinesisStreamAsSource)
	} else {
		input.DeliveryStreamType = aws.String(firehose.DeliveryStreamTypeDirectPut)
	}

	if d.Get("destination").(string) == destinationTypeExtendedS3 {
		extendedS3Config := createExtendedS3Config(d)
		input.ExtendedS3DestinationConfiguration = extendedS3Config
	} else {
		if d.Get("destination").(string) == destinationTypeElasticsearch {
			esConfig, err := createElasticsearchConfig(d)
			if err != nil {
				return sdkdiag.AppendErrorf(diags, "creating Kinesis Firehose Delivery Stream (%s): %s", sn, err)
			}
			input.ElasticsearchDestinationConfiguration = esConfig
		} else if d.Get("destination").(string) == destinationTypeOpensearch {
			osConfig, err := createOpensearchConfig(d)
			if err != nil {
				return sdkdiag.AppendErrorf(diags, "creating Kinesis Firehose Delivery Stream (%s): %s", sn, err)
			}
			input.AmazonopensearchserviceDestinationConfiguration = osConfig
		} else if d.Get("destination").(string) == destinationTypeRedshift {
			rc, err := createRedshiftConfig(d)
			if err != nil {
				return sdkdiag.AppendErrorf(diags, "creating Kinesis Firehose Delivery Stream (%s): %s", sn, err)
			}
			input.RedshiftDestinationConfiguration = rc
		} else if d.Get("destination").(string) == destinationTypeSplunk {
			spConfig, err := createSplunkConfig(d)
			if err != nil {
				return sdkdiag.AppendErrorf(diags, "creating Kinesis Firehose Delivery Stream (%s): %s", sn, err)
			}
			input.SplunkDestinationConfiguration = spConfig
		} else if d.Get("destination").(string) == destinationTypeHTTPEndpoint {
			httpConfig, err := createHTTPEndpointConfig(d)
			if err != nil {
				return sdkdiag.AppendErrorf(diags, "creating Kinesis Firehose Delivery Stream (%s): %s", sn, err)
			}
			input.HttpEndpointDestinationConfiguration = httpConfig
		}
	}

	err := retry.RetryContext(ctx, propagationTimeout, func() *retry.RetryError {
		_, err := conn.CreateDeliveryStreamWithContext(ctx, input)
		if err != nil {
			// Access was denied when calling Glue. Please ensure that the role specified in the data format conversion configuration has the necessary permissions.
			if tfawserr.ErrMessageContains(err, firehose.ErrCodeInvalidArgumentException, "Access was denied") {
				return retry.RetryableError(err)
			}

			if tfawserr.ErrMessageContains(err, firehose.ErrCodeInvalidArgumentException, "is not authorized to") {
				return retry.RetryableError(err)
			}

			if tfawserr.ErrMessageContains(err, firehose.ErrCodeInvalidArgumentException, "Please make sure the role specified in VpcConfiguration has permissions") {
				return retry.RetryableError(err)
			}

			// InvalidArgumentException: Verify that the IAM role has access to the Elasticsearch domain.
			if tfawserr.ErrMessageContains(err, firehose.ErrCodeInvalidArgumentException, "Verify that the IAM role has access") {
				return retry.RetryableError(err)
			}

			if tfawserr.ErrMessageContains(err, firehose.ErrCodeInvalidArgumentException, "Firehose is unable to assume role") {
				return retry.RetryableError(err)
			}

			return retry.NonRetryableError(err)
		}

		return nil
	})
	if tfresource.TimedOut(err) {
		_, err = conn.CreateDeliveryStreamWithContext(ctx, input)
	}
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating Kinesis Firehose Delivery Stream (%s): %s", sn, err)
	}

	s, err := waitDeliveryStreamCreated(ctx, conn, sn, d.Timeout(schema.TimeoutCreate))

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating Kinesis Firehose Delivery Stream (%s): waiting for completion: %s", sn, err)
	}

	d.SetId(aws.StringValue(s.DeliveryStreamARN))
	d.Set("arn", s.DeliveryStreamARN)

	if v, ok := d.GetOk("server_side_encryption"); ok && !isDeliveryStreamOptionDisabled(v) {
		startInput := &firehose.StartDeliveryStreamEncryptionInput{
			DeliveryStreamName:                         aws.String(sn),
			DeliveryStreamEncryptionConfigurationInput: expandDeliveryStreamEncryptionConfigurationInput(v.([]interface{})),
		}

		_, err := conn.StartDeliveryStreamEncryptionWithContext(ctx, startInput)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "starting Kinesis Firehose Delivery Stream (%s) encryption: %s", sn, err)
		}

		if _, err := waitDeliveryStreamEncryptionEnabled(ctx, conn, sn, d.Timeout(schema.TimeoutCreate)); err != nil {
			return sdkdiag.AppendErrorf(diags, "waiting for Kinesis Firehose Delivery Stream (%s) encryption enable: %s", sn, err)
		}
	}

	return append(diags, resourceDeliveryStreamRead(ctx, d, meta)...)
}

func validateSchema(d *schema.ResourceData) error {
	_, extendedS3Exists := d.GetOk("extended_s3_configuration")

	if d.Get("destination").(string) == destinationTypeExtendedS3 {
		if !extendedS3Exists {
			return fmt.Errorf(
				"when destination is 'extended_s3', extended_s3_configuration is required",
			)
		}
	}

	return nil
}

func resourceDeliveryStreamUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	sn := d.Get("name").(string)

	if err := validateSchema(d); err != nil {
		return sdkdiag.AppendErrorf(diags, "updating Kinesis Firehose Delivery Stream (%s): %s", sn, err)
	}

	conn := meta.(*conns.AWSClient).FirehoseConn(ctx)

	if d.HasChangesExcept("tags", "tags_all") {
		updateInput := &firehose.UpdateDestinationInput{
			DeliveryStreamName:             aws.String(sn),
			CurrentDeliveryStreamVersionId: aws.String(d.Get("version_id").(string)),
			DestinationId:                  aws.String(d.Get("destination_id").(string)),
		}

		if d.Get("destination").(string) == destinationTypeExtendedS3 {
			extendedS3Config := updateExtendedS3Config(d)
			updateInput.ExtendedS3DestinationUpdate = extendedS3Config
		} else {
			if d.Get("destination").(string) == destinationTypeElasticsearch {
				esUpdate, err := updateElasticsearchConfig(d)
				if err != nil {
					return sdkdiag.AppendErrorf(diags, "updating Kinesis Firehose Delivery Stream (%s): %s", sn, err)
				}
				updateInput.ElasticsearchDestinationUpdate = esUpdate
			} else if d.Get("destination").(string) == destinationTypeOpensearch {
				esUpdate, err := updateOpensearchConfig(d)
				if err != nil {
					return sdkdiag.AppendErrorf(diags, "updating Kinesis Firehose Delivery Stream (%s): %s", sn, err)
				}
				updateInput.AmazonopensearchserviceDestinationUpdate = esUpdate
			} else if d.Get("destination").(string) == destinationTypeRedshift {
				rc, err := updateRedshiftConfig(d)
				if err != nil {
					return sdkdiag.AppendErrorf(diags, "updating Kinesis Firehose Delivery Stream (%s): %s", sn, err)
				}
				updateInput.RedshiftDestinationUpdate = rc
			} else if d.Get("destination").(string) == destinationTypeSplunk {
				rc, err := updateSplunkConfig(d)
				if err != nil {
					return sdkdiag.AppendErrorf(diags, "updating Kinesis Firehose Delivery Stream (%s): %s", sn, err)
				}
				updateInput.SplunkDestinationUpdate = rc
			} else if d.Get("destination").(string) == destinationTypeHTTPEndpoint {
				rc, err := updateHTTPEndpointConfig(d)
				if err != nil {
					return sdkdiag.AppendErrorf(diags, "updating Kinesis Firehose Delivery Stream (%s): %s", sn, err)
				}
				updateInput.HttpEndpointDestinationUpdate = rc
			}
		}

		err := retry.RetryContext(ctx, propagationTimeout, func() *retry.RetryError {
			_, err := conn.UpdateDestinationWithContext(ctx, updateInput)
			if err != nil {
				// Access was denied when calling Glue. Please ensure that the role specified in the data format conversion configuration has the necessary permissions.
				if tfawserr.ErrMessageContains(err, firehose.ErrCodeInvalidArgumentException, "Access was denied") {
					return retry.RetryableError(err)
				}

				if tfawserr.ErrMessageContains(err, firehose.ErrCodeInvalidArgumentException, "is not authorized to") {
					return retry.RetryableError(err)
				}

				if tfawserr.ErrMessageContains(err, firehose.ErrCodeInvalidArgumentException, "Please make sure the role specified in VpcConfiguration has permissions") {
					return retry.RetryableError(err)
				}

				// InvalidArgumentException: Verify that the IAM role has access to the Elasticsearch domain.
				if tfawserr.ErrMessageContains(err, firehose.ErrCodeInvalidArgumentException, "Verify that the IAM role has access") {
					return retry.RetryableError(err)
				}

				if tfawserr.ErrMessageContains(err, firehose.ErrCodeInvalidArgumentException, "Firehose is unable to assume role") {
					return retry.RetryableError(err)
				}

				return retry.NonRetryableError(err)
			}

			return nil
		})

		if tfresource.TimedOut(err) {
			_, err = conn.UpdateDestinationWithContext(ctx, updateInput)
		}

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "updating Kinesis Firehose Delivery Stream (%s): %s", sn, err)
		}
	}

	if d.HasChange("server_side_encryption") {
		_, n := d.GetChange("server_side_encryption")
		if isDeliveryStreamOptionDisabled(n) {
			_, err := conn.StopDeliveryStreamEncryptionWithContext(ctx, &firehose.StopDeliveryStreamEncryptionInput{
				DeliveryStreamName: aws.String(sn),
			})

			if err != nil {
				return sdkdiag.AppendErrorf(diags, "stopping Kinesis Firehose Delivery Stream (%s) encryption: %s", sn, err)
			}

			if _, err := waitDeliveryStreamEncryptionDisabled(ctx, conn, sn, d.Timeout(schema.TimeoutUpdate)); err != nil {
				return sdkdiag.AppendErrorf(diags, "waiting for Kinesis Firehose Delivery Stream (%s) encryption disable: %s", sn, err)
			}
		} else {
			startInput := &firehose.StartDeliveryStreamEncryptionInput{
				DeliveryStreamName:                         aws.String(sn),
				DeliveryStreamEncryptionConfigurationInput: expandDeliveryStreamEncryptionConfigurationInput(n.([]interface{})),
			}

			_, err := conn.StartDeliveryStreamEncryptionWithContext(ctx, startInput)

			if err != nil {
				return sdkdiag.AppendErrorf(diags, "starting Kinesis Firehose Delivery Stream (%s) encryption: %s", sn, err)
			}

			if _, err := waitDeliveryStreamEncryptionEnabled(ctx, conn, sn, d.Timeout(schema.TimeoutUpdate)); err != nil {
				return sdkdiag.AppendErrorf(diags, "waiting for Kinesis Firehose Delivery Stream (%s) encryption enable: %s", sn, err)
			}
		}
	}

	return append(diags, resourceDeliveryStreamRead(ctx, d, meta)...)
}

func resourceDeliveryStreamRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).FirehoseConn(ctx)

	sn := d.Get("name").(string)
	s, err := FindDeliveryStreamByName(ctx, conn, sn)

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Kinesis Firehose Delivery Stream (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Kinesis Firehose Delivery Stream (%s): %s", sn, err)
	}

	if err := flattenDeliveryStream(d, s); err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Kinesis Firehose Delivery Stream (%s): %s", sn, err)
	}

	return diags
}

func resourceDeliveryStreamDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).FirehoseConn(ctx)

	sn := d.Get("name").(string)
	log.Printf("[DEBUG] Deleting Kinesis Firehose Delivery Stream: (%s)", sn)
	_, err := conn.DeleteDeliveryStreamWithContext(ctx, &firehose.DeleteDeliveryStreamInput{
		DeliveryStreamName: aws.String(sn),
	})

	if tfawserr.ErrCodeEquals(err, firehose.ErrCodeResourceNotFoundException) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting Kinesis Firehose Delivery Stream (%s): %s", sn, err)
	}

	_, err = waitDeliveryStreamDeleted(ctx, conn, sn, d.Timeout(schema.TimeoutDelete))

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for Kinesis Firehose Delivery Stream (%s) delete: %s", sn, err)
	}

	return diags
}

func isDeliveryStreamOptionDisabled(v interface{}) bool {
	options := v.([]interface{})
	if len(options) == 0 || options[0] == nil {
		return true
	}
	optionMap := options[0].(map[string]interface{})

	var enabled bool

	if v, ok := optionMap["enabled"]; ok {
		enabled = v.(bool)
	}

	return !enabled
}

func expandDeliveryStreamEncryptionConfigurationInput(tfList []interface{}) *firehose.DeliveryStreamEncryptionConfigurationInput {
	if len(tfList) == 0 {
		return nil
	}

	tfMap, ok := tfList[0].(map[string]interface{})

	if !ok {
		return nil
	}

	apiObject := &firehose.DeliveryStreamEncryptionConfigurationInput{}

	if v, ok := tfMap["key_arn"].(string); ok && v != "" {
		apiObject.KeyARN = aws.String(v)
	}

	if v, ok := tfMap["key_type"].(string); ok && v != "" {
		apiObject.KeyType = aws.String(v)
	}

	return apiObject
}
