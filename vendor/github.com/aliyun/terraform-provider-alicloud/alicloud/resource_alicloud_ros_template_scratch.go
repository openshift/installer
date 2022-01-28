package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudRosTemplateScratch() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRosTemplateScratchCreate,
		Read:   resourceAlicloudRosTemplateScratchRead,
		Update: resourceAlicloudRosTemplateScratchUpdate,
		Delete: resourceAlicloudRosTemplateScratchDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"execution_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Async", "Sync"}, false),
			},
			"logical_id_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"LongTypePrefixAndIndexSuffix", "LongTypePrefixAndHashSuffix", "ShortTypePrefixAndHashSuffix"}, false),
			},
			"preference_parameters": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parameter_key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"parameter_value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"source_tag": {
				Type:         schema.TypeSet,
				Optional:     true,
				MaxItems:     1,
				ExactlyOneOf: []string{"source_resource_group", "source_tag", "source_resources"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_tags": {
							Type:     schema.TypeMap,
							Required: true,
						},
						"resource_type_filter": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							MaxItems: 20,
						},
					},
				},
			},
			"source_resource_group": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_group_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"resource_type_filter": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							MaxItems: 20,
						},
					},
				},
			},
			"source_resources": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"template_scratch_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ResourceImport", "ArchitectureReplication"}, false),
			},
		},
	}
}

func resourceAlicloudRosTemplateScratchCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateTemplateScratch"
	request := make(map[string]interface{})
	conn, err := client.NewRosClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("execution_mode"); ok {
		request["ExecutionMode"] = v
	}
	if v, ok := d.GetOk("logical_id_strategy"); ok {
		request["LogicalIdStrategy"] = v
	}
	if v, ok := d.GetOk("preference_parameters"); ok {
		preferenceParametersMaps := make([]map[string]interface{}, 0)
		for _, preferenceParameters := range v.(*schema.Set).List() {
			preferenceParametersArg := preferenceParameters.(map[string]interface{})
			preferenceParametersMap := map[string]interface{}{}
			preferenceParametersMap["ParameterKey"] = preferenceParametersArg["parameter_key"]
			preferenceParametersMap["ParameterValue"] = preferenceParametersArg["parameter_value"]
			preferenceParametersMaps = append(preferenceParametersMaps, preferenceParametersMap)
		}
		if v, err := convertListMapToJsonString(preferenceParametersMaps); err != nil {
			return WrapError(err)
		} else {
			request["PreferenceParameters"] = v
		}
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("source_resource_group"); ok {
		sourceResourceGroupMap := map[string]interface{}{}
		for _, sourceResourceGroup := range v.(*schema.Set).List() {
			sourceResourceGroupArg := sourceResourceGroup.(map[string]interface{})
			sourceResourceGroupMap["ResourceGroupId"] = sourceResourceGroupArg["resource_group_id"]
			resourceTypeFilterList := make([]string, 0)
			for _, v := range sourceResourceGroupArg["resource_type_filter"].([]interface{}) {
				resourceTypeFilterList = append(resourceTypeFilterList, v.(string))
			}
			sourceResourceGroupMap["ResourceTypeFilter"] = resourceTypeFilterList
		}
		if v, err := convertMaptoJsonString(sourceResourceGroupMap); err != nil {
			return WrapError(err)
		} else {
			request["SourceResourceGroup"] = v
		}
	}
	if v, ok := d.GetOk("source_resources"); ok {
		sourceResourcesMaps := make([]map[string]interface{}, 0)
		for _, sourceResources := range v.(*schema.Set).List() {
			sourceResourcesArg := sourceResources.(map[string]interface{})
			sourceResourcesMap := map[string]interface{}{}
			sourceResourcesMap["ResourceId"] = sourceResourcesArg["resource_id"]
			sourceResourcesMap["ResourceType"] = sourceResourcesArg["resource_type"]
			sourceResourcesMaps = append(sourceResourcesMaps, sourceResourcesMap)
		}
		if v, err := convertListMapToJsonString(sourceResourcesMaps); err != nil {
			return WrapError(err)
		} else {
			request["SourceResources"] = v
		}
	}

	if v, ok := d.GetOk("source_tag"); ok {
		sourceTagMap := map[string]interface{}{}
		for _, sourceResources := range v.(*schema.Set).List() {
			sourceTagArg := sourceResources.(map[string]interface{})
			sourceTagMap["ResourceTags"] = sourceTagArg["resource_tags"]
			resourceTypeFilterList := make([]string, 0)
			for _, v := range sourceTagArg["resource_type_filter"].([]interface{}) {
				resourceTypeFilterList = append(resourceTypeFilterList, v.(string))
			}
			sourceTagMap["ResourceTypeFilter"] = resourceTypeFilterList
		}
		if v, err := convertMaptoJsonString(sourceTagMap); err != nil {
			return WrapError(err)
		} else {
			request["SourceTag"] = v
		}
	}
	request["TemplateScratchType"] = d.Get("template_scratch_type")
	request["ClientToken"] = buildClientToken("CreateTemplateScratch")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-10"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ros_template_scratch", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["TemplateScratchId"]))
	rosService := RosService{client}
	stateConf := BuildStateConf([]string{}, []string{"GENERATE_COMPLETE"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, rosService.RosTemplateScratchStateRefreshFunc(d.Id(), []string{"GENERATE_FAILED"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudRosTemplateScratchRead(d, meta)
}
func resourceAlicloudRosTemplateScratchRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rosService := RosService{client}
	object, err := rosService.DescribeRosTemplateScratch(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ros_template_scratch rosService.DescribeRosTemplateScratch Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("description", object["Description"])
	d.Set("logical_id_strategy", object["LogicalIdStrategy"])

	if preferenceParametersList, ok := object["PreferenceParameters"]; ok && preferenceParametersList != nil {
		preferenceParametersMaps := make([]map[string]interface{}, 0)
		for _, preferenceParametersListItem := range preferenceParametersList.([]interface{}) {
			if preferenceParametersListItemMap, ok := preferenceParametersListItem.(map[string]interface{}); ok {
				if v, ok := preferenceParametersListItemMap["ParameterValue"]; ok && fmt.Sprint(v) != "" {
					preferenceParametersListItemArg := make(map[string]interface{}, 0)
					preferenceParametersListItemArg["parameter_key"] = preferenceParametersListItemMap["ParameterKey"]
					preferenceParametersListItemArg["parameter_value"] = v
					preferenceParametersMaps = append(preferenceParametersMaps, preferenceParametersListItemArg)
				}
			}
		}
		d.Set("preference_parameters", preferenceParametersMaps)
	}

	sourceTagSli := make([]map[string]interface{}, 0)
	if v, ok := object["SourceTag"]; ok {
		if sourceTag, ok := v.(map[string]interface{}); ok && len(sourceTag) > 0 {
			sourceTagMap := make(map[string]interface{})
			sourceTagMap["resource_tags"] = sourceTag["ResourceTags"]
			resourceTypeFilter := make([]interface{}, 0)
			if v, ok := sourceTag["ResourceTypeFilter"]; ok {
				if vv, ok := v.([]interface{}); ok && len(vv) > 0 {
					resourceTypeFilter = append(resourceTypeFilter, vv...)
				}
			}
			sourceTagMap["resource_type_filter"] = resourceTypeFilter
			sourceTagSli = append(sourceTagSli, sourceTagMap)
		}
	}

	d.Set("source_tag", sourceTagSli)

	sourceResourceGroupSli := make([]map[string]interface{}, 0)
	if v, ok := object["SourceResourceGroup"]; ok {
		if sourceResourceGroup, ok := v.(map[string]interface{}); ok && len(sourceResourceGroup) > 0 {
			sourceResourceGroupMap := make(map[string]interface{})
			sourceResourceGroupMap["resource_group_id"] = sourceResourceGroup["ResourceGroupId"]
			resourceTypeFilter := make([]interface{}, 0)
			if v, ok := sourceResourceGroup["ResourceTypeFilter"]; ok {
				if vv, ok := v.([]interface{}); ok && len(vv) > 0 {
					resourceTypeFilter = append(resourceTypeFilter, vv...)
				}
			}
			sourceResourceGroupMap["resource_type_filter"] = resourceTypeFilter
			sourceResourceGroupSli = append(sourceResourceGroupSli, sourceResourceGroupMap)
		}
	}

	d.Set("source_resource_group", sourceResourceGroupSli)
	if sourceResourcesList, ok := object["SourceResources"]; ok && sourceResourcesList != nil {
		sourceResourcesMaps := make([]map[string]interface{}, 0)
		for _, sourceResourcesListItem := range sourceResourcesList.([]interface{}) {
			if sourceResourcesListItemMap, ok := sourceResourcesListItem.(map[string]interface{}); ok {
				sourceResourcesListItemArg := make(map[string]interface{}, 0)
				sourceResourcesListItemArg["resource_id"] = sourceResourcesListItemMap["ResourceId"]
				sourceResourcesListItemArg["resource_type"] = sourceResourcesListItemMap["ResourceType"]
				sourceResourcesMaps = append(sourceResourcesMaps, sourceResourcesListItemArg)
			}
		}
		d.Set("source_resources", sourceResourcesMaps)
	}

	d.Set("status", object["Status"])
	d.Set("template_scratch_type", object["TemplateScratchType"])
	return nil
}
func resourceAlicloudRosTemplateScratchUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rosService := RosService{client}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"TemplateScratchId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
	}
	if d.HasChange("logical_id_strategy") {
		update = true
		if v, ok := d.GetOk("logical_id_strategy"); ok {
			request["LogicalIdStrategy"] = v
		}
	}
	if d.HasChange("source_tag") {
		update = true
		if v, ok := d.GetOkExists("source_tag"); ok {
			sourceTagMap := map[string]interface{}{}
			for _, sourceResources := range v.(*schema.Set).List() {
				sourceTagArg := sourceResources.(map[string]interface{})
				if len(sourceTagArg["resource_tags"].(map[string]interface{})) > 0 {
					sourceTagMap["ResourceTags"] = sourceTagArg["resource_tags"]
					resourceTypeFilterList := make([]string, 0)
					for _, v := range sourceTagArg["resource_type_filter"].([]interface{}) {
						resourceTypeFilterList = append(resourceTypeFilterList, v.(string))
					}
					sourceTagMap["ResourceTypeFilter"] = resourceTypeFilterList
				}
			}
			if v, err := convertMaptoJsonString(sourceTagMap); err != nil {
				return WrapError(err)
			} else {
				request["SourceTag"] = v
			}
		}
	}

	if d.HasChange("preference_parameters") {
		update = true
	}
	if v, ok := d.GetOk("preference_parameters"); ok {
		preferenceParametersMaps := make([]map[string]interface{}, 0)
		for _, preferenceParameters := range v.(*schema.Set).List() {
			preferenceParametersArg := preferenceParameters.(map[string]interface{})
			preferenceParametersMap := map[string]interface{}{}
			preferenceParametersMap["ParameterKey"] = preferenceParametersArg["parameter_key"]
			preferenceParametersMap["ParameterValue"] = preferenceParametersArg["parameter_value"]
			preferenceParametersMaps = append(preferenceParametersMaps, preferenceParametersMap)
		}
		if v, err := convertListMapToJsonString(preferenceParametersMaps); err != nil {
			return WrapError(err)
		} else {
			request["PreferenceParameters"] = v
		}
	}

	if d.HasChange("source_resource_group") {
		update = true
		if v, ok := d.GetOk("source_resource_group"); ok {
			sourceResourceGroupMap := map[string]interface{}{}
			for _, sourceResourceGroup := range v.(*schema.Set).List() {
				sourceResourceGroupArg := sourceResourceGroup.(map[string]interface{})
				sourceResourceGroupMap["ResourceGroupId"] = sourceResourceGroupArg["resource_group_id"]
				resourceTypeFilterList := make([]string, 0)
				for _, v := range sourceResourceGroupArg["resource_type_filter"].([]interface{}) {
					resourceTypeFilterList = append(resourceTypeFilterList, v.(string))
				}
				sourceResourceGroupMap["ResourceTypeFilter"] = resourceTypeFilterList
			}
			if v, err := convertMaptoJsonString(sourceResourceGroupMap); err != nil {
				return WrapError(err)
			} else {
				request["SourceResourceGroup"] = v
			}
		}
	}
	if d.HasChange("source_resources") {
		update = true
		if v, ok := d.GetOk("source_resources"); ok {
			sourceResourcesMaps := make([]map[string]interface{}, 0)
			for _, sourceResources := range v.(*schema.Set).List() {
				sourceResourcesArg := sourceResources.(map[string]interface{})
				sourceResourcesMap := map[string]interface{}{}
				sourceResourcesMap["ResourceId"] = sourceResourcesArg["resource_id"]
				sourceResourcesMap["ResourceType"] = sourceResourcesArg["resource_type"]
				sourceResourcesMaps = append(sourceResourcesMaps, sourceResourcesMap)
			}
			if v, err := convertListMapToJsonString(sourceResourcesMaps); err != nil {
				return WrapError(err)
			} else {
				request["SourceResources"] = v
			}
		}
	}
	if update {
		if v, ok := d.GetOk("execution_mode"); ok {
			request["ExecutionMode"] = v
		}
		action := "UpdateTemplateScratch"
		conn, err := client.NewRosClient()
		if err != nil {
			return WrapError(err)
		}
		request["ClientToken"] = buildClientToken("UpdateTemplateScratch")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-10"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"GENERATE_COMPLETE"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, rosService.RosTemplateScratchStateRefreshFunc(d.Id(), []string{"GENERATE_FAILED"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	return resourceAlicloudRosTemplateScratchRead(d, meta)
}
func resourceAlicloudRosTemplateScratchDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rosService := RosService{client}
	action := "DeleteTemplateScratch"
	var response map[string]interface{}
	conn, err := client.NewRosClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"TemplateScratchId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, rosService.RosTemplateScratchStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
