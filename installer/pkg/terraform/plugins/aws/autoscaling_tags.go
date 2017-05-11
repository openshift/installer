package aws

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

// autoscalingTagSchema returns the schema to use for the tag element.
func autoscalingTagSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},

				"value": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},

				"propagate_at_launch": &schema.Schema{
					Type:     schema.TypeBool,
					Required: true,
				},
			},
		},
		Set: autoscalingTagToHash,
	}
}

func autoscalingTagToHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["key"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["value"].(string)))
	buf.WriteString(fmt.Sprintf("%t-", m["propagate_at_launch"].(bool)))

	return hashcode.String(buf.String())
}

// setTags is a helper to set the tags for a resource. It expects the
// tags field to be named "tag"
func setAutoscalingTags(conn *autoscaling.AutoScaling, d *schema.ResourceData) error {
	resourceID := d.Get("name").(string)
	var createTags, removeTags []*autoscaling.Tag

	if d.HasChange("tag") || d.HasChange("tags") {
		oraw, nraw := d.GetChange("tag")
		o := setToMapByKey(oraw.(*schema.Set), "key")
		n := setToMapByKey(nraw.(*schema.Set), "key")

		c, r := diffAutoscalingTags(
			autoscalingTagsFromMap(o, resourceID),
			autoscalingTagsFromMap(n, resourceID),
			resourceID)

		createTags = append(createTags, c...)
		removeTags = append(removeTags, r...)

		oraw, nraw = d.GetChange("tags")

		c, r = diffAutoscalingTags(
			autoscalingTagsFromList(oraw.([]interface{}), resourceID),
			autoscalingTagsFromList(nraw.([]interface{}), resourceID),
			resourceID)

		createTags = append(createTags, c...)
		removeTags = append(removeTags, r...)
	}

	// Set tags
	if len(removeTags) > 0 {
		log.Printf("[DEBUG] Removing autoscaling tags: %#v", removeTags)

		remove := autoscaling.DeleteTagsInput{
			Tags: removeTags,
		}

		if _, err := conn.DeleteTags(&remove); err != nil {
			return err
		}
	}

	if len(createTags) > 0 {
		log.Printf("[DEBUG] Creating autoscaling tags: %#v", createTags)

		create := autoscaling.CreateOrUpdateTagsInput{
			Tags: createTags,
		}

		if _, err := conn.CreateOrUpdateTags(&create); err != nil {
			return err
		}
	}

	return nil
}

// diffTags takes our tags locally and the ones remotely and returns
// the set of tags that must be created, and the set of tags that must
// be destroyed.
func diffAutoscalingTags(oldTags, newTags []*autoscaling.Tag, resourceID string) ([]*autoscaling.Tag, []*autoscaling.Tag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		tag := map[string]interface{}{
			"key":                 *t.Key,
			"value":               *t.Value,
			"propagate_at_launch": *t.PropagateAtLaunch,
		}
		create[*t.Key] = tag
	}

	// Build the list of what to remove
	var remove []*autoscaling.Tag
	for _, t := range oldTags {
		old, ok := create[*t.Key].(map[string]interface{})

		if !ok || old["value"] != *t.Value || old["propagate_at_launch"] != *t.PropagateAtLaunch {
			// Delete it!
			remove = append(remove, t)
		}
	}

	return autoscalingTagsFromMap(create, resourceID), remove
}

func autoscalingTagsFromList(vs []interface{}, resourceID string) []*autoscaling.Tag {
	result := make([]*autoscaling.Tag, 0, len(vs))
	for _, tag := range vs {
		attr, ok := tag.(map[string]interface{})
		if !ok {
			continue
		}

		if t := autoscalingTagFromMap(attr, resourceID); t != nil {
			result = append(result, t)
		}
	}
	return result
}

// tagsFromMap returns the tags for the given map of data.
func autoscalingTagsFromMap(m map[string]interface{}, resourceID string) []*autoscaling.Tag {
	result := make([]*autoscaling.Tag, 0, len(m))
	for _, v := range m {
		attr, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		t := autoscalingTagFromMap(attr, resourceID)
		if t != nil {
			result = append(result, t)
		}
	}

	return result
}

func autoscalingTagFromMap(attr map[string]interface{}, resourceID string) *autoscaling.Tag {
	if _, ok := attr["key"]; !ok {
		return nil
	}

	if _, ok := attr["value"]; !ok {
		return nil
	}

	if _, ok := attr["propagate_at_launch"]; !ok {
		return nil
	}

	var propagate_at_launch bool

	if v, ok := attr["propagate_at_launch"].(bool); ok {
		propagate_at_launch = v
	}

	if v, ok := attr["propagate_at_launch"].(string); ok {
		propagate_at_launch, _ = strconv.ParseBool(v)
	}

	t := &autoscaling.Tag{
		Key:               aws.String(attr["key"].(string)),
		Value:             aws.String(attr["value"].(string)),
		PropagateAtLaunch: aws.Bool(propagate_at_launch),
		ResourceId:        aws.String(resourceID),
		ResourceType:      aws.String("auto-scaling-group"),
	}

	if tagIgnoredAutoscaling(t) {
		return nil
	}

	return t
}

// autoscalingTagsToMap turns the list of tags into a map.
func autoscalingTagsToMap(ts []*autoscaling.Tag) map[string]interface{} {
	tags := make(map[string]interface{})
	for _, t := range ts {
		tag := map[string]interface{}{
			"key":                 *t.Key,
			"value":               *t.Value,
			"propagate_at_launch": *t.PropagateAtLaunch,
		}
		tags[*t.Key] = tag
	}

	return tags
}

// autoscalingTagDescriptionsToMap turns the list of tags into a map.
func autoscalingTagDescriptionsToMap(ts *[]*autoscaling.TagDescription) map[string]map[string]interface{} {
	tags := make(map[string]map[string]interface{})
	for _, t := range *ts {
		tag := map[string]interface{}{
			"key":                 *t.Key,
			"value":               *t.Value,
			"propagate_at_launch": *t.PropagateAtLaunch,
		}
		tags[*t.Key] = tag
	}

	return tags
}

// autoscalingTagDescriptionsToSlice turns the list of tags into a slice.
func autoscalingTagDescriptionsToSlice(ts []*autoscaling.TagDescription) []map[string]interface{} {
	tags := make([]map[string]interface{}, 0, len(ts))
	for _, t := range ts {
		tags = append(tags, map[string]interface{}{
			"key":                 *t.Key,
			"value":               *t.Value,
			"propagate_at_launch": *t.PropagateAtLaunch,
		})
	}

	return tags
}

func setToMapByKey(s *schema.Set, key string) map[string]interface{} {
	result := make(map[string]interface{})
	for _, rawData := range s.List() {
		data := rawData.(map[string]interface{})
		result[data[key].(string)] = data
	}

	return result
}

// compare a tag against a list of strings and checks if it should
// be ignored or not
func tagIgnoredAutoscaling(t *autoscaling.Tag) bool {
	filter := []string{"^aws:*"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching %v with %v\n", v, *t.Key)
		if r, _ := regexp.MatchString(v, *t.Key); r == true {
			log.Printf("[DEBUG] Found AWS specific tag %s (val: %s), ignoring.\n", *t.Key, *t.Value)
			return true
		}
	}
	return false
}
