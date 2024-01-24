package alicloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/gpdb"

	"regexp"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cdn"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/elasticsearch"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func String(v string) *string {
	return &v
}

func tagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
	}
}

func tagsSchemaComputed() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Computed: true,
	}
}

func tagsSchemaWithIgnore() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validation.StringDoesNotMatch(regexp.MustCompile(`(^acs:.*)|(^aliyun.*)|(/.*http://.*\.\w+/gm)|(/.*https://.*\.\w+/gm)`), "It cannot begin with \"aliyun\", \"acs:\"; without \"http://\", and \"https://\"."),
		},
	}
}

func parsingTags(d *schema.ResourceData) (map[string]interface{}, []string) {
	oraw, nraw := d.GetChange("tags")
	removedTags := oraw.(map[string]interface{})
	addedTags := nraw.(map[string]interface{})
	// Build the list of what to remove
	removed := make([]string, 0)
	for key, value := range removedTags {
		old, ok := addedTags[key]
		if !ok || old != value {
			// Delete it!
			removed = append(removed, key)
		}
	}

	return addedTags, removed
}

func tagsToMap(tags interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	if tags == nil {
		return result
	}
	switch v := tags.(type) {
	case map[string]interface{}:
		for key, value := range tags.(map[string]interface{}) {
			if !tagIgnored(key, value) {
				result[key] = value
			}
		}
	case []interface{}:
		if len(tags.([]interface{})) < 1 {
			return result
		}
		for _, tag := range tags.([]interface{}) {
			t := tag.(map[string]interface{})
			var tagKey string
			var tagValue interface{}
			if v, ok := t["TagKey"]; ok {
				tagKey = v.(string)
				tagValue = t["TagValue"]
			} else if v, ok := t["Key"]; ok {
				tagKey = v.(string)
				tagValue = t["Value"]
			}
			if !tagIgnored(tagKey, tagValue) {
				result[tagKey] = tagValue
			}
		}
	default:
		log.Printf("\u001B[31m[ERROR]\u001B[0m Unknown tags type %s. The tags value is: %v.", v, tags)
	}
	return result
}

func tagIgnored(tagKey string, tagValue interface{}) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, tagKey)
		ok, _ := regexp.MatchString(v, tagKey)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific tag %s (val: %s), ignoring.\n", tagKey, tagValue)
			return true
		}
	}
	return false
}

// setTags is a helper to set the tags for a resource. It expects the
// tags field to be named "tags"
func setTags(client *connectivity.AliyunClient, resourceType TagResourceType, d *schema.ResourceData) error {
	if d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		return updateTags(client, []string{d.Id()}, resourceType, oraw, nraw)
	}

	return nil
}

func setCdnTags(client *connectivity.AliyunClient, resourceType TagResourceType, d *schema.ResourceData) error {
	if d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		return updateCdnTags(client, []string{d.Id()}, resourceType, oraw, nraw)
	}

	return nil
}

func setVolumeTags(client *connectivity.AliyunClient, resourceType TagResourceType, d *schema.ResourceData) error {
	ecsService := EcsService{client}
	if d.HasChange("volume_tags") {
		request := ecs.CreateDescribeDisksRequest()
		request.InstanceId = d.Id()
		var response *ecs.DescribeDisksResponse
		wait := incrementalWait(1*time.Second, 1*time.Second)
		err := resource.Retry(10*time.Minute, func() *resource.RetryError {
			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.DescribeDisks(request)
			})
			if err != nil {
				if IsThrottling(err) {
					wait()
					return resource.RetryableError(err)

				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			response, _ = raw.(*ecs.DescribeDisksResponse)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		if len(response.Disks.Disk) == 0 {
			return WrapError(Error("no specified system disk"))
		}

		var ids []string
		systemDiskTag := make(map[string]interface{})
		for _, disk := range response.Disks.Disk {
			ids = append(ids, disk.DiskId)
			if disk.Type == "system" {
				for _, t := range disk.Tags.Tag {
					if !ecsService.ecsTagIgnored(t) {
						systemDiskTag[t.TagKey] = t.TagValue
					}
				}
			}
		}

		oraw, nraw := d.GetChange("volume_tags")
		if d.IsNewResource() {
			oraw = systemDiskTag
		}
		return updateTags(client, ids, resourceType, oraw, nraw)
	}

	return nil
}

func updateTags(client *connectivity.AliyunClient, ids []string, resourceType TagResourceType, oraw, nraw interface{}) error {
	o := oraw.(map[string]interface{})
	n := nraw.(map[string]interface{})
	create, remove := diffTags(tagsFromMap(o), tagsFromMap(n))

	// Set tags
	if len(remove) > 0 {
		request := ecs.CreateUntagResourcesRequest()
		request.ResourceType = string(resourceType)
		request.ResourceId = &ids

		var tagsKey []string
		for _, t := range remove {
			tagsKey = append(tagsKey, t.Key)
		}
		request.TagKey = &tagsKey

		wait := incrementalWait(1*time.Second, 1*time.Second)
		err := resource.Retry(10*time.Minute, func() *resource.RetryError {
			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.UntagResources(request)
			})
			if err != nil {
				if IsThrottling(err) {
					wait()
					return resource.RetryableError(err)

				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, ids, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}

	if len(create) > 0 {
		request := ecs.CreateTagResourcesRequest()
		request.ResourceType = string(resourceType)
		request.ResourceId = &ids

		var tags []ecs.TagResourcesTag
		for _, t := range create {
			tags = append(tags, ecs.TagResourcesTag{
				Key:   t.Key,
				Value: t.Value,
			})
		}
		request.Tag = &tags

		wait := incrementalWait(1*time.Second, 1*time.Second)
		err := resource.Retry(10*time.Minute, func() *resource.RetryError {
			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.TagResources(request)
			})
			if err != nil {
				if IsThrottling(err) {
					wait()
					return resource.RetryableError(err)

				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, ids, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}

	return nil
}

func updateCdnTags(client *connectivity.AliyunClient, ids []string, resourceType TagResourceType, oraw, nraw interface{}) error {
	o := oraw.(map[string]interface{})
	n := nraw.(map[string]interface{})
	create, remove := diffTags(tagsFromMap(o), tagsFromMap(n))

	// Set tags
	if len(remove) > 0 {
		request := cdn.CreateUntagResourcesRequest()
		request.ResourceType = string(resourceType)
		request.ResourceId = &ids

		var tagsKey []string
		for _, t := range remove {
			tagsKey = append(tagsKey, t.Key)
		}
		request.TagKey = &tagsKey

		raw, err := client.WithCdnClient_new(func(cdnClient *cdn.Client) (interface{}, error) {
			return cdnClient.UntagResources(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, ids, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	if len(create) > 0 {
		request := cdn.CreateTagResourcesRequest()
		request.ResourceType = string(resourceType)
		request.ResourceId = &ids

		var tags []cdn.TagResourcesTag
		for _, t := range create {
			tags = append(tags, cdn.TagResourcesTag{
				Key:   t.Key,
				Value: t.Value,
			})
		}
		request.Tag = &tags

		raw, err := client.WithCdnClient_new(func(cdnClient *cdn.Client) (interface{}, error) {
			return cdnClient.TagResources(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, ids, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return nil
}

// diffTags takes our tags locally and the ones remotely and returns
// the set of tags that must be created, and the set of tags that must
// be destroyed.
func diffTags(oldTags, newTags []Tag) ([]Tag, []Tag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []Tag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			// Delete it!
			remove = append(remove, t)
		}
	}

	return tagsFromMap(create), remove
}

func diffGpdbTags(oldTags, newTags []gpdb.TagResourcesTag) ([]gpdb.TagResourcesTag, []gpdb.TagResourcesTag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []gpdb.TagResourcesTag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			// Delete it!
			remove = append(remove, t)
		}
	}

	return gpdbTagsFromMap(create), remove
}

// tagsFromMap returns the tags for the given map of data.
func tagsFromMap(m map[string]interface{}) []Tag {
	result := make([]Tag, 0, len(m))
	for k, v := range m {
		result = append(result, Tag{
			Key:   k,
			Value: v.(string),
		})
	}

	return result
}

func gpdbTagsFromMap(m map[string]interface{}) []gpdb.TagResourcesTag {
	result := make([]gpdb.TagResourcesTag, 0, len(m))
	for k, v := range m {
		result = append(result, gpdb.TagResourcesTag{
			Key:   k,
			Value: v.(string),
		})
	}

	return result
}

func ecsTagsToMap(tags []ecs.Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !ecsTagIgnored(t) {
			result[t.TagKey] = t.TagValue
		}
	}

	return result
}

func elasticsearchTagsToMap(tags []elasticsearch.Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !elasticsearchTagIgnored(t.TagKey, t.TagValue) {
			result[t.TagKey] = t.TagValue
		}
	}

	return result
}

func vpcTagsToMap(tags []vpc.Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !vpcTagIgnored(t) {
			result[t.Key] = t.Value
		}
	}
	return result
}

func cdnTagsToMap(tags []cdn.TagItem) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !cdnTagIgnored(t) {
			result[t.Key] = t.Value
		}
	}

	return result
}

func slbTagsToMap(tags []slb.TagResource) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !slbTagIgnored(t) {
			result[t.TagKey] = t.TagValue
		}
	}

	return result
}

func essTagsToMap(tags []ess.Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !essTagIgnored(t) {
			result[t.Key] = t.Value
		}
	}

	return result
}

func otsTagsToMap(tags []ots.TagInfo) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		result[t.TagKey] = t.TagValue
	}

	return result
}

func tagsMapEqual(expectMap map[string]interface{}, compareMap map[string]string) bool {
	if len(expectMap) != len(compareMap) {
		return false
	} else {
		for key, eVal := range expectMap {
			if eStr, ok := eVal.(string); !ok {
				// type is mismatch.
				return false
			} else {
				if cStr, ok := compareMap[key]; ok {
					if eStr != cStr {
						return false
					}
				} else {
					return false
				}
			}
		}
	}
	return true
}

// tagIgnored compares a tag against a list of strings and checks if it should be ignored or not
func ecsTagIgnored(t ecs.Tag) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.TagKey)
		ok, _ := regexp.MatchString(v, t.TagKey)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific tag %s (val: %s), ignoring.\n", t.TagKey, t.TagValue)
			return true
		}
	}
	return false
}

func vpcTagIgnored(t vpc.Tag) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.Key)
		ok, _ := regexp.MatchString(v, t.Key)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific tag %s (val: %s), ignoring.\n", t.Key, t.Value)
			return true
		}
	}
	return false
}

// tagIgnored compares a tag against a list of strings and checks if it should be ignored or not
func essTagIgnored(t ess.Tag) bool {
	filter := []string{"^aliyun", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.Key)
		ok, _ := regexp.MatchString(v, t.Key)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific tag %s (val: %s), ignoring.\n", t.Key, t.Value)
			return true
		}
	}
	return false
}

func cdnTagIgnored(t cdn.TagItem) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.Key)
		ok, _ := regexp.MatchString(v, t.Key)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific tag %s (val: %s), ignoring.\n", t.Key, t.Value)
			return true
		}
	}
	return false
}

func slbTagIgnored(t slb.TagResource) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.TagKey)
		ok, _ := regexp.MatchString(v, t.TagKey)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific tag %s (val: %s), ignoring.\n", t.TagKey, t.TagValue)
			return true
		}
	}
	return false
}

func elasticsearchTagIgnored(tagKey, tagValue string) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, tagKey)
		ok, _ := regexp.MatchString(v, tagKey)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific tag %s (val: %s), ignoring.\n", tagKey, tagValue)
			return true
		}
	}
	return false
}

func ignoredTags(tagKey, tagValue string) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, tagKey)
		ok, _ := regexp.MatchString(v, tagKey)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific tag %s (val: %s), ignoring.\n", tagKey, tagValue)
			return true
		}
	}
	return false
}
