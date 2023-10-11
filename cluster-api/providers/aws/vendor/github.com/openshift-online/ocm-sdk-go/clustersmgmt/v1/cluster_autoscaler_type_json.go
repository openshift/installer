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

package v1 // github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalClusterAutoscaler writes a value of the 'cluster_autoscaler' type to the given writer.
func MarshalClusterAutoscaler(object *ClusterAutoscaler, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeClusterAutoscaler(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeClusterAutoscaler writes a value of the 'cluster_autoscaler' type to the given stream.
func writeClusterAutoscaler(object *ClusterAutoscaler, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(ClusterAutoscalerLinkKind)
	} else {
		stream.WriteString(ClusterAutoscalerKind)
	}
	count++
	if object.bitmap_&2 != 0 {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
		count++
	}
	if object.bitmap_&4 != 0 {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("href")
		stream.WriteString(object.href)
		count++
	}
	var present_ bool
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("balance_similar_node_groups")
		stream.WriteBool(object.balanceSimilarNodeGroups)
		count++
	}
	present_ = object.bitmap_&16 != 0 && object.balancingIgnoredLabels != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("balancing_ignored_labels")
		writeStringList(object.balancingIgnoredLabels, stream)
		count++
	}
	present_ = object.bitmap_&32 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("ignore_daemonsets_utilization")
		stream.WriteBool(object.ignoreDaemonsetsUtilization)
		count++
	}
	present_ = object.bitmap_&64 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("log_verbosity")
		stream.WriteInt(object.logVerbosity)
		count++
	}
	present_ = object.bitmap_&128 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("max_node_provision_time")
		stream.WriteString(object.maxNodeProvisionTime)
		count++
	}
	present_ = object.bitmap_&256 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("max_pod_grace_period")
		stream.WriteInt(object.maxPodGracePeriod)
		count++
	}
	present_ = object.bitmap_&512 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("pod_priority_threshold")
		stream.WriteInt(object.podPriorityThreshold)
		count++
	}
	present_ = object.bitmap_&1024 != 0 && object.resourceLimits != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("resource_limits")
		writeAutoscalerResourceLimits(object.resourceLimits, stream)
		count++
	}
	present_ = object.bitmap_&2048 != 0 && object.scaleDown != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("scale_down")
		writeAutoscalerScaleDownConfig(object.scaleDown, stream)
		count++
	}
	present_ = object.bitmap_&4096 != 0
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
	object = readClusterAutoscaler(iterator)
	err = iterator.Error
	return
}

// readClusterAutoscaler reads a value of the 'cluster_autoscaler' type from the given iterator.
func readClusterAutoscaler(iterator *jsoniter.Iterator) *ClusterAutoscaler {
	object := &ClusterAutoscaler{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == ClusterAutoscalerLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "balance_similar_node_groups":
			value := iterator.ReadBool()
			object.balanceSimilarNodeGroups = value
			object.bitmap_ |= 8
		case "balancing_ignored_labels":
			value := readStringList(iterator)
			object.balancingIgnoredLabels = value
			object.bitmap_ |= 16
		case "ignore_daemonsets_utilization":
			value := iterator.ReadBool()
			object.ignoreDaemonsetsUtilization = value
			object.bitmap_ |= 32
		case "log_verbosity":
			value := iterator.ReadInt()
			object.logVerbosity = value
			object.bitmap_ |= 64
		case "max_node_provision_time":
			value := iterator.ReadString()
			object.maxNodeProvisionTime = value
			object.bitmap_ |= 128
		case "max_pod_grace_period":
			value := iterator.ReadInt()
			object.maxPodGracePeriod = value
			object.bitmap_ |= 256
		case "pod_priority_threshold":
			value := iterator.ReadInt()
			object.podPriorityThreshold = value
			object.bitmap_ |= 512
		case "resource_limits":
			value := readAutoscalerResourceLimits(iterator)
			object.resourceLimits = value
			object.bitmap_ |= 1024
		case "scale_down":
			value := readAutoscalerScaleDownConfig(iterator)
			object.scaleDown = value
			object.bitmap_ |= 2048
		case "skip_nodes_with_local_storage":
			value := iterator.ReadBool()
			object.skipNodesWithLocalStorage = value
			object.bitmap_ |= 4096
		default:
			iterator.ReadAny()
		}
	}
	return object
}
