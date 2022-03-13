package alicloud

import (
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

func dataSourceAlicloudHbrVaults() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudHbrVaultsRead,
		Schema: map[string]*schema.Schema{
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"CREATED", "ERROR", "INITIALIZING", "UNKNOWN"}, false),
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"vault_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "STANDARD",
				ValidateFunc: validation.StringInSlice([]string{"STANDARD"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vaults": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bytes_done": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dedup": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"index_available": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"index_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"index_update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"latest_replication_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"replication": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"replication_source_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"replication_source_vault_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"retention": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"search_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"source_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_size": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vault_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vault_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vault_status_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vault_storage_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vault_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudHbrVaultsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeVaults"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}
	if v, ok := d.GetOk("vault_type"); ok {
		request["VaultType"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var vaultNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		vaultNameRegex = r
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	var response map[string]interface{}
	conn, err := client.NewHbrClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_hbr_vaults", action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		resp, err := jsonpath.Get("$.Vaults.Vault", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Vaults.Vault", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if vaultNameRegex != nil && !vaultNameRegex.MatchString(fmt.Sprint(item["VaultName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["VaultId"])]; !ok {
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
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"bucket_name":                  object["BucketName"],
			"bytes_done":                   fmt.Sprint(object["BytesDone"]),
			"created_time":                 fmt.Sprint(object["CreatedTime"]),
			"dedup":                        object["Dedup"],
			"description":                  object["Description"],
			"index_available":              object["IndexAvailable"],
			"index_level":                  object["IndexLevel"],
			"index_update_time":            fmt.Sprint(object["IndexUpdateTime"]),
			"latest_replication_time":      fmt.Sprint(object["LatestReplicationTime"]),
			"payment_type":                 object["ChargeType"],
			"replication":                  object["Replication"],
			"replication_source_region_id": object["ReplicationSourceRegionId"],
			"replication_source_vault_id":  object["ReplicationSourceVaultId"],
			"retention":                    fmt.Sprint(object["Retention"]),
			"search_enabled":               object["SearchEnabled"],
			"source_types":                 object["SourceTypes"].(map[string]interface{})["SourceType"],
			"status":                       object["Status"],
			"storage_size":                 fmt.Sprint(object["StorageSize"]),
			"updated_time":                 fmt.Sprint(object["UpdatedTime"]),
			"id":                           fmt.Sprint(object["VaultId"]),
			"vault_id":                     fmt.Sprint(object["VaultId"]),
			"vault_name":                   object["VaultName"],
			"vault_status_message":         object["VaultStatusMessage"],
			"vault_storage_class":          object["VaultStorageClass"],
			"vault_type":                   object["VaultType"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["VaultName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("vaults", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
