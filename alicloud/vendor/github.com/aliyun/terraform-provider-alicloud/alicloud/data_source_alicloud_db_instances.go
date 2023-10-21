package alicloud

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudDBInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDBInstancesRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(MySQL),
					string(SQLServer),
					string(PPAS),
					string(PostgreSQL),
				}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				// please follow the link below to see more details on available statusesplease follow the link below to see more details on available statuses
				// https://help.aliyun.com/document_detail/26315.html
			},
			"db_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Primary",
					"Readonly",
					"Guard",
					"Temp",
				}, false),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"connection_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Standard",
					"Safe",
				}, false),
			},
			"tags": tagsSchema(),
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parameters": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"force_modify": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"checking_code": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"parameter_value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"force_restart": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"parameter_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"parameter_description": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expire_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"net_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"guard_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"temp_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"readonly_instance_ids": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_string": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_storage_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_storage": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"master_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id_slave_a": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id_slave_b": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"ssl_expire_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"require_update": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"acl": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ca_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_ca_cert": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_ca_cert_expire_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_cert_revocation_list": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_modify_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modify_status_reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"replication_acl": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"require_update_item": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"require_update_reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ssl_create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ssl_enabled": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_ca_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_cert": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delete_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"encryption_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"encryption_key_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_usage": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"material_expire_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"origin": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDBInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DescribeDBInstances"
	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"SourceIp":   client.SourceIp,
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
	if v, ok := d.GetOk("engine"); ok && v.(string) != "" {
		request["Engine"] = v.(string)
	}
	if v, ok := d.GetOk("status"); ok && v.(string) != "" {
		request["DBInstanceStatus"] = v.(string)
	}
	if v, ok := d.GetOk("db_type"); ok && v.(string) != "" {
		request["DBInstanceType"] = v.(string)
	}
	if v, ok := d.GetOk("vpc_id"); ok && v.(string) != "" {
		request["VpcId"] = v.(string)
	}
	if v, ok := d.GetOk("vswitch_id"); ok && v.(string) != "" {
		request["VSwitchId"] = v.(string)
	}
	if v, ok := d.GetOk("connection_mode"); ok && v.(string) != "" {
		request["ConnectionMode"] = v.(string)
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := v.(map[string]interface{})
		bs, err := json.Marshal(tagsMap)
		if err != nil {
			return WrapError(err)
		}
		request["Tags"] = string(bs)
	}

	var response map[string]interface{}
	conn, err := client.NewRdsClient()
	if err != nil {
		return WrapError(err)
	}

	var objects []map[string]interface{}

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nameRegex = r
	}

	// ids
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_db_instances", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Items.DBInstance", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Items.DBInstance", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if nameRegex != nil {
				if !nameRegex.MatchString(fmt.Sprint(item["DBInstanceDescription"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["DBInstanceId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return rdsInstancesDescription(d, meta, objects)
}

func rdsInstancesDescription(d *schema.ResourceData, meta interface{}, objects []map[string]interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}

	var ids []string
	var names []string
	var s []map[string]interface{}

	for _, item := range objects {
		readOnlyInstanceIDs := []string{}
		for _, id := range item["ReadOnlyDBInstanceIds"].(map[string]interface{})["ReadOnlyDBInstanceId"].([]interface{}) {
			readOnlyInstanceIDs = append(readOnlyInstanceIDs, fmt.Sprint(id.(map[string]interface{})["DBInstanceId"]))
		}
		instance, err := rdsService.DescribeDBInstance(fmt.Sprint(item["DBInstanceId"]))
		if err != nil {
			return WrapError(err)
		}

		mapping := map[string]interface{}{
			"id":                       fmt.Sprint(item["DBInstanceId"]),
			"name":                     fmt.Sprint(item["DBInstanceDescription"]),
			"charge_type":              fmt.Sprint(item["PayType"]),
			"db_type":                  fmt.Sprint(item["DBInstanceType"]),
			"region_id":                fmt.Sprint(item["RegionId"]),
			"create_time":              fmt.Sprint(item["CreateTime"]),
			"expire_time":              fmt.Sprint(item["ExpireTime"]),
			"status":                   fmt.Sprint(item["DBInstanceStatus"]),
			"engine":                   fmt.Sprint(item["Engine"]),
			"engine_version":           fmt.Sprint(item["EngineVersion"]),
			"net_type":                 fmt.Sprint(item["DBInstanceNetType"]),
			"connection_mode":          fmt.Sprint(item["ConnectionMode"]),
			"instance_type":            fmt.Sprint(item["DBInstanceClass"]),
			"availability_zone":        fmt.Sprint(item["ZoneId"]),
			"master_instance_id":       fmt.Sprint(item["MasterInstanceId"]),
			"guard_instance_id":        fmt.Sprint(item["GuardDBInstanceId"]),
			"temp_instance_id":         fmt.Sprint(item["TempDBInstanceId"]),
			"readonly_instance_ids":    readOnlyInstanceIDs,
			"vpc_id":                   fmt.Sprint(item["VpcId"]),
			"vswitch_id":               fmt.Sprint(item["VSwitchId"]),
			"connection_string":        instance["ConnectionString"],
			"port":                     instance["Port"],
			"db_instance_storage_type": instance["DBInstanceStorageType"],
			"instance_storage":         instance["DBInstanceStorage"],
			"master_zone":              instance["MasterZone"],
		}
		sslResponse, sslErr := rdsService.DescribeDBInstanceSSL(fmt.Sprint(item["DBInstanceId"]))
		if sslErr == nil {
			if v, ok := sslResponse["SSLExpireTime"]; ok && v != "" {
				mapping["ssl_expire_time"] = sslResponse["SSLExpireTime"]
			}
			if v, ok := sslResponse["RequireUpdate"]; ok && v != "" {
				mapping["require_update"] = sslResponse["RequireUpdate"]
			}
			if v, ok := sslResponse["ACL"]; ok && v != "" {
				mapping["acl"] = sslResponse["ACL"]
			}
			if v, ok := sslResponse["CAType"]; ok && v != "" {
				mapping["ca_type"] = sslResponse["CAType"]
			}
			if v, ok := sslResponse["ClientCACert"]; ok && v != "" {
				mapping["client_ca_cert"] = sslResponse["ClientCACert"]
			}
			if v, ok := sslResponse["ClientCACertExpireTime"]; ok && v != "" {
				mapping["client_ca_cert_expire_time"] = sslResponse["ClientCACertExpireTime"]
			}
			if v, ok := sslResponse["ClientCertRevocationList"]; ok && v != "" {
				mapping["client_cert_revocation_list"] = sslResponse["ClientCertRevocationList"]
			}
			if v, ok := sslResponse["LastModifyStatus"]; ok && v != "" {
				mapping["last_modify_status"] = sslResponse["LastModifyStatus"]
			}
			if v, ok := sslResponse["ModifyStatusReason"]; ok && v != "" {
				mapping["modify_status_reason"] = sslResponse["ModifyStatusReason"]
			}
			if v, ok := sslResponse["ReplicationACL"]; ok && v != "" {
				mapping["replication_acl"] = sslResponse["ReplicationACL"]
			}
			if v, ok := sslResponse["RequireUpdateItem"]; ok && v != "" {
				mapping["require_update_item"] = sslResponse["RequireUpdateItem"]
			}
			if v, ok := sslResponse["RequireUpdateReason"]; ok && v != "" {
				mapping["require_update_reason"] = sslResponse["RequireUpdateReason"]
			}
			if v, ok := sslResponse["SSLCreateTime"]; ok && v != "" {
				mapping["ssl_create_time"] = sslResponse["SSLCreateTime"]
			}
			if v, ok := sslResponse["SSLEnabled"]; ok && v != "" {
				mapping["ssl_enabled"] = sslResponse["SSLEnabled"]
			}
			if v, ok := sslResponse["ServerCAUrl"]; ok && v != "" {
				mapping["server_ca_url"] = sslResponse["ServerCAUrl"]
			}
			if v, ok := sslResponse["ServerCert"]; ok && v != "" {
				mapping["server_cert"] = sslResponse["ServerCert"]
			}
			if v, ok := sslResponse["ServerKey"]; ok && v != "" {
				mapping["server_key"] = sslResponse["ServerKey"]
			}

		}

		encResponse, encError := rdsService.DescribeDBInstanceEncryptionKey(fmt.Sprint(item["DBInstanceId"]))
		if encError == nil {
			if v, ok := encResponse["Creator"]; ok && v != "" {
				mapping["creator"] = encResponse["Creator"]
			}
			if v, ok := encResponse["DeleteDate"]; ok && v != "" {
				mapping["delete_date"] = encResponse["DeleteDate"]
			}
			if v, ok := encResponse["Description"]; ok && v != "" {
				mapping["description"] = encResponse["Description"]
			}
			if v, ok := encResponse["EncryptionKey"]; ok && v != "" {
				mapping["encryption_key"] = encResponse["EncryptionKey"]
			}
			if v, ok := encResponse["EncryptionKeyStatus"]; ok && v != "" {
				mapping["encryption_key_status"] = encResponse["EncryptionKeyStatus"]
			}
			if v, ok := encResponse["KeyUsage"]; ok && v != "" {
				mapping["key_usage"] = encResponse["KeyUsage"]
			}
			if v, ok := encResponse["MaterialExpireTime"]; ok && v != "" {
				mapping["material_expire_time"] = encResponse["MaterialExpireTime"]
			}
			if v, ok := encResponse["Origin"]; ok && v != "" {
				mapping["origin"] = encResponse["Origin"]
			}
		}

		slaveZones := instance["SlaveZones"].(map[string]interface{})["SlaveZone"].([]interface{})
		if len(slaveZones) == 2 {
			mapping["zone_id_slave_a"] = slaveZones[0].(map[string]interface{})["ZoneId"]
			mapping["zone_id_slave_b"] = slaveZones[1].(map[string]interface{})["ZoneId"]
		} else if len(slaveZones) == 1 {
			mapping["zone_id_slave_a"] = slaveZones[0].(map[string]interface{})["ZoneId"]
		}

		if d.Get("enable_details").(bool) {
			paramResponse, paramError := rdsService.DescribeParameterTemplates(item["DBInstanceId"].(string), item["Engine"].(string), item["EngineVersion"].(string))
			if paramError == nil {
				parameterDetail := make([]map[string]interface{}, 0)
				for _, val := range paramResponse {
					item := val.(map[string]interface{})
					parameterDetail = append(parameterDetail, map[string]interface{}{
						"force_modify":          item["ForceModify"],
						"checking_code":         item["CheckingCode"],
						"parameter_value":       item["ParameterValue"],
						"force_restart":         item["ForceRestart"],
						"parameter_name":        item["ParameterName"],
						"parameter_description": item["ParameterDescription"],
					})
				}
				mapping["parameters"] = parameterDetail
			}
		}
		ids = append(ids, fmt.Sprint(item["DBInstanceId"]))
		names = append(names, fmt.Sprint(item["DBInstanceDescription"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("instances", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
