package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudLogOssShipper() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudLogOssShipperCreate,
		Read:   resourceAlicloudLogOssShipperRead,
		Update: resourceAlicloudLogOssShipperUpdate,
		Delete: resourceAlicloudLogOssShipperDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"logstore_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"shipper_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"oss_bucket": {
				Type:     schema.TypeString,
				Required: true,
			},
			"oss_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"buffer_interval": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"buffer_size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"role_arn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"compress_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"none", "snappy"}, false),
				Default:      "none",
			},
			"path_format": {
				Type:     schema.TypeString,
				Required: true,
			},

			"format": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"json", "parquet", "csv"}, false),
			},

			"json_enable_tag": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"csv_config_delimiter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"csv_config_header": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"csv_config_linefeed": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"csv_config_columns": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"csv_config_nullidentifier": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"csv_config_quote": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parquet_config": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudLogOssShipperCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var requestInfo *sls.Client
	projectName := d.Get("project_name").(string)
	logstoreName := d.Get("logstore_name").(string)
	shipperName := d.Get("shipper_name").(string)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	if err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			project, _ := sls.NewLogProject(projectName, slsClient.Endpoint, slsClient.AccessKeyID, slsClient.AccessKeySecret)
			project, _ = project.WithToken(slsClient.SecurityToken)
			logstore, _ := sls.NewLogStore(logstoreName, project)
			return nil, logstore.CreateShipper(buildConfig(d))
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("CreateLogShipper", raw, requestInfo, map[string]string{
			"project_name":  projectName,
			"logstore_name": logstoreName,
			"shipper_name":  shipperName,
		})
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_oss_shipper", "CreateLogOssShipper", AliyunLogGoSdkERROR)
	}
	d.SetId(fmt.Sprintf("%s%s%s%s%s", projectName, COLON_SEPARATED, logstoreName, COLON_SEPARATED, shipperName))
	return resourceAlicloudLogOssShipperRead(d, meta)
}

func resourceAlicloudLogOssShipperRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	shipper, err := logService.DescribeLogOssShipper(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_log_oss_shipper LogService.DescribeLogOssShipper Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	ossShipperConfig := shipper.TargetConfiguration.(*sls.OSSShipperConfig)
	d.Set("project_name", parts[0])
	d.Set("logstore_name", parts[1])
	d.Set("shipper_name", parts[2])
	d.Set("oss_bucket", ossShipperConfig.OssBucket)
	d.Set("oss_prefix", ossShipperConfig.OssPrefix)
	d.Set("buffer_interval", ossShipperConfig.BufferInterval)
	d.Set("buffer_size", ossShipperConfig.BufferSize)
	d.Set("role_arn", ossShipperConfig.RoleArn)
	d.Set("compress_type", ossShipperConfig.CompressType)
	d.Set("path_format", ossShipperConfig.PathFormat)
	d.Set("format", ossShipperConfig.Storage.Format)

	if ossShipperConfig.Storage.Format == "json" {
		detail := ossShipperConfig.Storage.Detail.(map[string]interface{})
		d.Set("json_enable_tag", detail["enableTag"])
	} else if ossShipperConfig.Storage.Format == "csv" {
		detail := ossShipperConfig.Storage.Detail.(map[string]interface{})
		d.Set("csv_config_delimiter", detail["delimiter"])
		d.Set("csv_config_header", detail["header"])
		d.Set("csv_config_linefeed", detail["linefeed"])
		d.Set("csv_config_columns", detail["columns"])
		d.Set("csv_config_nullidentifier", detail["nullidentifier"])
		d.Set("csv_config_quote", detail["quote"])
	} else if ossShipperConfig.Storage.Format == "parquet" {
		var config []map[string]interface{}
		detail := ossShipperConfig.Storage.Detail.(map[string]interface{})
		tempConfigs := detail["columns"].([]interface{})
		for _, temp := range tempConfigs {
			parquetConfig := temp.(map[string]interface{})
			tempMap := map[string]interface{}{
				"name": parquetConfig["name"],
				"type": parquetConfig["type"],
			}
			config = append(config, tempMap)
		}
		d.Set("parquet_config", config)
	}
	return nil
}

func resourceAlicloudLogOssShipperUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	if err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			project, _ := sls.NewLogProject(parts[0], slsClient.Endpoint, slsClient.AccessKeyID, slsClient.AccessKeySecret)
			project, _ = project.WithToken(slsClient.SecurityToken)
			logstore, _ := sls.NewLogStore(parts[1], project)
			return nil, logstore.UpdateShipper(buildConfig(d))
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateLogOssShipper", AliyunLogGoSdkERROR)
	}
	return resourceAlicloudLogOssShipperRead(d, meta)

}

func resourceAlicloudLogOssShipperDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	var requestInfo *sls.Client
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			project, _ := sls.NewLogProject(parts[0], slsClient.Endpoint, slsClient.AccessKeyID, slsClient.AccessKeySecret)
			project, _ = project.WithToken(slsClient.SecurityToken)
			logstore, _ := sls.NewLogStore(parts[1], project)
			return nil, logstore.DeleteShipper(parts[2])
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if debugOn() {
			addDebug("DeleteLogOssShipper", raw, requestInfo, map[string]interface{}{
				"project_name":  parts[0],
				"logstore_name": parts[1],
				"shipper_name":  parts[2],
			})
		}
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ShipperNotExist"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_oss_shipper", "DeleteLogOssShipper", AliyunLogGoSdkERROR)
	}
	return WrapError(logService.WaitForLogOssShipper(d.Id(), Deleted, DefaultTimeout))

}

func buildConfig(d *schema.ResourceData) *sls.Shipper {
	format := d.Get("format").(string)
	var storage sls.ShipperStorage
	if format == "json" {
		enableTag := d.Get("json_enable_tag").(bool)
		storage = sls.ShipperStorage{
			Format: "json",
			Detail: sls.OssStorageJsonDetail{EnableTag: enableTag},
		}
	} else if format == "parquet" {
		detail := sls.OssStoreageParquet{}
		for _, f := range d.Get("parquet_config").(*schema.Set).List() {
			v := f.(map[string]interface{})
			config := sls.ParquetConfig{
				Name: v["name"].(string),
				Type: v["type"].(string),
			}
			detail.Columns = append(detail.Columns, config)
		}
		storage = sls.ShipperStorage{
			Format: "parquet",
			Detail: detail,
		}

	} else if format == "csv" {
		detail := sls.OssStoreageCsvDetail{
			Delimiter:      d.Get("csv_config_delimiter").(string),
			Header:         d.Get("csv_config_header").(bool),
			LineFeed:       d.Get("csv_config_linefeed").(string),
			NullIdentifier: d.Get("csv_config_nullidentifier").(string),
			Quote:          d.Get("csv_config_quote").(string),
		}
		columns := []string{}
		for _, v := range d.Get("csv_config_columns").([]interface{}) {
			columns = append(columns, v.(string))
		}
		detail.Columns = columns
		storage = sls.ShipperStorage{
			Format: "csv",
			Detail: detail,
		}
	}

	ossShipperConfig := &sls.OSSShipperConfig{
		OssBucket:      d.Get("oss_bucket").(string),
		OssPrefix:      d.Get("oss_prefix").(string),
		RoleArn:        d.Get("role_arn").(string),
		BufferInterval: d.Get("buffer_interval").(int),
		BufferSize:     d.Get("buffer_size").(int),
		CompressType:   d.Get("compress_type").(string),
		PathFormat:     d.Get("path_format").(string),
		Storage:        storage,
	}
	ossShipper := &sls.Shipper{
		ShipperName:         d.Get("shipper_name").(string),
		TargetType:          "oss",
		TargetConfiguration: ossShipperConfig,
	}
	return ossShipper

}
