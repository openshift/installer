/*
Copyright (c) 2020 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// IMPORTANT: This file has been generated automatically, refrain from modifying it manually as all
// your changes will be lost when the file is generated again.

package v1alpha1 // github.com/openshift-online/ocm-api-model/clientapi/arohcp/v1alpha1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalClusterAutoscaler writes a value of the 'cluster_autoscaler' type to the given writer.
func MarshalClusterAutoscaler(object *ClusterAutoscaler, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteClusterAutoscaler(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteClusterAutoscaler writes a value of the 'cluster_autoscaler' type to the given stream.
func WriteClusterAutoscaler(object *ClusterAutoscaler, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(ClusterAutoscalerLinkKind)
	} else {
		stream.WriteString(ClusterAutoscalerKind)
	}
	count++
	if len(object.fieldSet_) > 1 && object.fieldSet_[1] {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
		count++
	}
	if len(object.fieldSet_) > 2 && object.fieldSet_[2] {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("href")
		stream.WriteString(object.href)
		count++
	}
	var present_ bool
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("balance_similar_node_groups")
		stream.WriteBool(object.balanceSimilarNodeGroups)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4] && object.balancingIgnoredLabels != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("balancing_ignored_labels")
		WriteStringList(object.balancingIgnoredLabels, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("ignore_daemonsets_utilization")
		stream.WriteBool(object.ignoreDaemonsetsUtilization)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("log_verbosity")
		stream.WriteInt(object.logVerbosity)
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("max_node_provision_time")
		stream.WriteString(object.maxNodeProvisionTime)
		count++
	}
	present_ = len(object.fieldSet_) > 8 && object.fieldSet_[8]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("max_pod_grace_period")
		stream.WriteInt(object.maxPodGracePeriod)
		count++
	}
	present_ = len(object.fieldSet_) > 9 && object.fieldSet_[9]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("pod_priority_threshold")
		stream.WriteInt(object.podPriorityThreshold)
		count++
	}
	present_ = len(object.fieldSet_) > 10 && object.fieldSet_[10] && object.resourceLimits != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("resource_limits")
		WriteAutoscalerResourceLimits(object.resourceLimits, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 11 && object.fieldSet_[11] && object.scaleDown != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("scale_down")
		WriteAutoscalerScaleDownConfig(object.scaleDown, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 12 && object.fieldSet_[12]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("skip_nodes_with_local_storage")
		stream.WriteBool(object.skipNodesWithLocalStorage)
	}
	stream.WriteObjectEnd()
}

// UnmarshalClusterAutoscaler reads a value of the 'cluster_autoscaler' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalClusterAutoscaler(source interface{}) (object *ClusterAutoscaler, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadClusterAutoscaler(iterator)
	err = iterator.Error
	return
}

// ReadClusterAutoscaler reads a value of the 'cluster_autoscaler' type from the given iterator.
func ReadClusterAutoscaler(iterator *jsoniter.Iterator) *ClusterAutoscaler {
	object := &ClusterAutoscaler{
		fieldSet_: make([]bool, 13),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == ClusterAutoscalerLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "balance_similar_node_groups":
			value := iterator.ReadBool()
			object.balanceSimilarNodeGroups = value
			object.fieldSet_[3] = true
		case "balancing_ignored_labels":
			value := ReadStringList(iterator)
			object.balancingIgnoredLabels = value
			object.fieldSet_[4] = true
		case "ignore_daemonsets_utilization":
			value := iterator.ReadBool()
			object.ignoreDaemonsetsUtilization = value
			object.fieldSet_[5] = true
		case "log_verbosity":
			value := iterator.ReadInt()
			object.logVerbosity = value
			object.fieldSet_[6] = true
		case "max_node_provision_time":
			value := iterator.ReadString()
			object.maxNodeProvisionTime = value
			object.fieldSet_[7] = true
		case "max_pod_grace_period":
			value := iterator.ReadInt()
			object.maxPodGracePeriod = value
			object.fieldSet_[8] = true
		case "pod_priority_threshold":
			value := iterator.ReadInt()
			object.podPriorityThreshold = value
			object.fieldSet_[9] = true
		case "resource_limits":
			value := ReadAutoscalerResourceLimits(iterator)
			object.resourceLimits = value
			object.fieldSet_[10] = true
		case "scale_down":
			value := ReadAutoscalerScaleDownConfig(iterator)
			object.scaleDown = value
			object.fieldSet_[11] = true
		case "skip_nodes_with_local_storage":
			value := iterator.ReadBool()
			object.skipNodesWithLocalStorage = value
			object.fieldSet_[12] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
