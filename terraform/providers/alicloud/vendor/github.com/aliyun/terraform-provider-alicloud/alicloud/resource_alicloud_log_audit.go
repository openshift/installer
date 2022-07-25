package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	slsPop "github.com/aliyun/alibaba-cloud-sdk-go/services/sls"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudLogAudit() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudLogAuditCreate,
		Read:   resourceAlicloudLogAuditRead,
		Update: resourceAlicloudLogAuditUpdate,
		Delete: resourceAlicloudLogAuditDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"aliuid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"variable_map": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"multi_account": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"resource_directory_type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"custom", "all"}, true),
				Optional:     true,
			},
		},
	}
}

func resourceAlicloudLogAuditCreate(d *schema.ResourceData, meta interface{}) error {

	d.SetId(d.Get("display_name").(string))
	return resourceAlicloudLogAuditUpdate(d, meta)
}

func resourceAlicloudLogAuditUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := slsPop.CreateAnalyzeAppLogRequest()
	request.AppType = "audit"
	request.DisplayName = d.Id()

	var variableMap = map[string]interface{}{}
	mutiAccount := expandStringList(d.Get("multi_account").(*schema.Set).List())

	if resourceDirectoryType, ok := d.GetOk("resource_directory_type"); ok {
		resourceDirectoryMap := map[string]interface{}{}
		resourceDirectoryMap["type"] = resourceDirectoryType
		resourceDirectoryMap["multi_account"] = mutiAccount
		data, err := json.Marshal(resourceDirectoryMap)
		if err != nil {
			return WrapError(err)
		}
		resultResourceDirectory := string(data)
		variableMap["resource_directory"] = resultResourceDirectory
	} else if len(mutiAccount) > 0 {
		mutiAccountList := []map[string]string{}
		for _, v := range mutiAccount {
			mutiAccountMap := map[string]string{}
			mutiAccountMap["uid"] = v
			mutiAccountList = append(mutiAccountList, mutiAccountMap)
		}
		data, err := json.Marshal(mutiAccountList)
		if err != nil {
			return WrapError(err)
		}
		resultMutiAccount := string(data)
		variableMap["multi_account"] = resultMutiAccount
	}
	variableMap["region"] = client.RegionId
	variableMap["aliuid"] = d.Get("aliuid").(string)
	variableMap["project"] = fmt.Sprintf("slsaudit-center-%s-%s", variableMap["aliuid"], variableMap["region"])
	variableMap["logstore"] = "slsaudit"

	if tempMap, ok := d.GetOk("variable_map"); ok {
		for k, v := range tempMap.(map[string]interface{}) {
			if strings.HasSuffix(k, "_policy_setting") {
				return Error("Does not support configuration %s in variable_map", k)
			}
			variableMap[k] = v
		}
	}

	b, err := json.Marshal(variableMap)
	if err != nil {
		return WrapError(err)
	} else {
		request.VariableMap = string(b[:])
	}
	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		rep, err := client.WithLogPopClient(func(client *slsPop.Client) (interface{}, error) {
			return client.AnalyzeAppLog(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), rep, request)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_audit", request.GetActionName(), AliyunLogGoSdkERROR)
	}
	return resourceAlicloudLogAuditRead(d, meta)
}

func resourceAlicloudLogAuditRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	response, err := logService.DescribeLogAudit(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	displayName, initMap, err := getInitParameter(response.GetHttpContentString())
	if err != nil {
		return WrapError(err)
	}
	d.Set("display_name", displayName)
	d.Set("aliuid", initMap["aliuid"].(string))
	if multiAccount, ok := initMap["multi_account"]; ok {
		account, err := analyzeMultiAccount(multiAccount.(string))
		if err != nil {
			return WrapError(err)
		}
		d.Set("multi_account", account)
	}
	if resourceDirectory, ok := initMap["resource_directory"]; ok {
		resourceDirectoryMap := map[string]interface{}{}
		err = json.Unmarshal([]byte(resourceDirectory.(string)), &resourceDirectoryMap)
		if err != nil {
			return WrapError(err)
		}
		if len(resourceDirectoryMap) > 0 {
			if rd_type, ok := resourceDirectoryMap["type"]; ok {
				d.Set("resource_directory_type", rd_type.(string))
			}
			if multiAccount, ok := resourceDirectoryMap["multi_account"]; ok {
				d.Set("multi_account", multiAccount)
			}
		}
	}
	for k := range initMap {
		if strings.HasSuffix(k, "_policy_setting") {
			delete(initMap, k)
		}
	}
	delete(initMap, "region")
	delete(initMap, "aliuid")
	delete(initMap, "project")
	delete(initMap, "logstore")
	delete(initMap, "multi_account")
	delete(initMap, "resource_directory")
	d.Set("variable_map", initMap)
	return nil
}

func resourceAlicloudLogAuditDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudLogAuditInstance.")
	return nil
}

func getInitParameter(rep string) (displayName string, initMap map[string]interface{}, err error) {
	m := make(map[string]interface{})
	err = json.Unmarshal([]byte(rep), &m)
	if _, ok := m["AppModel"].(map[string]interface{}); ok {
		model := make(map[string]interface{})
		err = json.Unmarshal([]byte(rep), &model)
		if d, ok := model["AppModel"].(map[string]interface{}); ok {
			displayName = d["DisplayName"].(string)
			configNew := d["Config"]
			m := make(map[string]interface{})
			err = json.Unmarshal([]byte(configNew.(string)), &m)
			initMap = m["initParam"].(map[string]interface{})
		}
	}
	return displayName, initMap, err
}

func analyzeMultiAccount(s string) ([]string, error) {
	var m []map[string]interface{}
	err := json.Unmarshal([]byte(s), &m)
	if err != nil {
		return nil, err
	}
	multiAccount := make([]string, len(m))
	for i := range m {
		if v, ok := m[i]["uid"].(string); ok {
			multiAccount[i] = v
		} else {
			multiAccount[i] = fmt.Sprintf("%.0f", m[i]["uid"].(float64))
		}
	}
	return multiAccount, nil
}
