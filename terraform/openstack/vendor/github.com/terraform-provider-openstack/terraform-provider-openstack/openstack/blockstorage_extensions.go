package openstack

import (
	"bytes"
	"fmt"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/schedulerhints"
	"github.com/gophercloud/utils/terraform/hashcode"
)

func expandBlockStorageExtensionsSchedulerHints(v schedulerhints.SchedulerHints) map[string]interface{} {
	schedulerHints := make(map[string]interface{})

	differentHost := make([]interface{}, len(v.DifferentHost))
	for i, dh := range v.DifferentHost {
		differentHost[i] = dh
	}

	sameHost := make([]interface{}, len(v.SameHost))
	for i, sh := range v.SameHost {
		sameHost[i] = sh
	}

	schedulerHints["different_host"] = differentHost
	schedulerHints["same_host"] = sameHost
	schedulerHints["local_to_instance"] = v.LocalToInstance
	schedulerHints["query"] = v.Query
	schedulerHints["additional_properties"] = v.AdditionalProperties
	return schedulerHints
}

func blockStorageExtensionsSchedulerHintsHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if m["query"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["query"].(string)))
	}

	if m["local_to_instance"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["local_to_instance"].(string)))
	}

	if m["additional_properties"] != nil {
		for _, v := range m["additional_properties"].(map[string]interface{}) {
			buf.WriteString(fmt.Sprintf("%s-", v))
		}
	}

	buf.WriteString(fmt.Sprintf("%s-", m["different_host"].([]interface{})))
	buf.WriteString(fmt.Sprintf("%s-", m["same_host"].([]interface{})))

	return hashcode.String(buf.String())
}

func resourceBlockStorageSchedulerHints(schedulerHintsRaw map[string]interface{}) schedulerhints.SchedulerHints {
	schedulerHints := schedulerhints.SchedulerHints{
		Query:                schedulerHintsRaw["query"].(string),
		LocalToInstance:      schedulerHintsRaw["local_to_instance"].(string),
		AdditionalProperties: schedulerHintsRaw["additional_properties"].(map[string]interface{}),
	}

	if v, ok := schedulerHintsRaw["different_host"].([]interface{}); ok {
		differentHost := make([]string, len(v))

		for i, dh := range v {
			differentHost[i] = dh.(string)
		}

		schedulerHints.DifferentHost = differentHost
	}

	if v, ok := schedulerHintsRaw["same_host"].([]interface{}); ok {
		sameHost := make([]string, len(v))

		for i, sh := range v {
			sameHost[i] = sh.(string)
		}

		schedulerHints.SameHost = sameHost
	}

	return schedulerHints
}
