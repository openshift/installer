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

package v1 // github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalSubscriptionMetrics writes a value of the 'subscription_metrics' type to the given writer.
func MarshalSubscriptionMetrics(object *SubscriptionMetrics, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeSubscriptionMetrics(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeSubscriptionMetrics writes a value of the 'subscription_metrics' type to the given stream.
func writeSubscriptionMetrics(object *SubscriptionMetrics, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cloud_provider")
		stream.WriteString(object.cloudProvider)
		count++
	}
	present_ = object.bitmap_&2 != 0 && object.computeNodesCpu != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("compute_nodes_cpu")
		writeClusterResource(object.computeNodesCpu, stream)
		count++
	}
	present_ = object.bitmap_&4 != 0 && object.computeNodesMemory != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("compute_nodes_memory")
		writeClusterResource(object.computeNodesMemory, stream)
		count++
	}
	present_ = object.bitmap_&8 != 0 && object.computeNodesSockets != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("compute_nodes_sockets")
		writeClusterResource(object.computeNodesSockets, stream)
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("console_url")
		stream.WriteString(object.consoleUrl)
		count++
	}
	present_ = object.bitmap_&32 != 0 && object.cpu != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cpu")
		writeClusterResource(object.cpu, stream)
		count++
	}
	present_ = object.bitmap_&64 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("critical_alerts_firing")
		stream.WriteFloat64(object.criticalAlertsFiring)
		count++
	}
	present_ = object.bitmap_&128 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("health_state")
		stream.WriteString(object.healthState)
		count++
	}
	present_ = object.bitmap_&256 != 0 && object.memory != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("memory")
		writeClusterResource(object.memory, stream)
		count++
	}
	present_ = object.bitmap_&512 != 0 && object.nodes != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("nodes")
		writeClusterMetricsNodes(object.nodes, stream)
		count++
	}
	present_ = object.bitmap_&1024 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("openshift_version")
		stream.WriteString(object.openshiftVersion)
		count++
	}
	present_ = object.bitmap_&2048 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("operating_system")
		stream.WriteString(object.operatingSystem)
		count++
	}
	present_ = object.bitmap_&4096 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("operators_condition_failing")
		stream.WriteFloat64(object.operatorsConditionFailing)
		count++
	}
	present_ = object.bitmap_&8192 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("region")
		stream.WriteString(object.region)
		count++
	}
	present_ = object.bitmap_&16384 != 0 && object.sockets != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("sockets")
		writeClusterResource(object.sockets, stream)
		count++
	}
	present_ = object.bitmap_&32768 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("state")
		stream.WriteString(object.state)
		count++
	}
	present_ = object.bitmap_&65536 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("state_description")
		stream.WriteString(object.stateDescription)
		count++
	}
	present_ = object.bitmap_&131072 != 0 && object.storage != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("storage")
		writeClusterResource(object.storage, stream)
		count++
	}
	present_ = object.bitmap_&262144 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subscription_cpu_total")
		stream.WriteFloat64(object.subscriptionCpuTotal)
		count++
	}
	present_ = object.bitmap_&524288 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subscription_obligation_exists")
		stream.WriteFloat64(object.subscriptionObligationExists)
		count++
	}
	present_ = object.bitmap_&1048576 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subscription_socket_total")
		stream.WriteFloat64(object.subscriptionSocketTotal)
		count++
	}
	present_ = object.bitmap_&2097152 != 0 && object.upgrade != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("upgrade")
		writeClusterUpgrade(object.upgrade, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalSubscriptionMetrics reads a value of the 'subscription_metrics' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalSubscriptionMetrics(source interface{}) (object *SubscriptionMetrics, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readSubscriptionMetrics(iterator)
	err = iterator.Error
	return
}

// readSubscriptionMetrics reads a value of the 'subscription_metrics' type from the given iterator.
func readSubscriptionMetrics(iterator *jsoniter.Iterator) *SubscriptionMetrics {
	object := &SubscriptionMetrics{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "cloud_provider":
			value := iterator.ReadString()
			object.cloudProvider = value
			object.bitmap_ |= 1
		case "compute_nodes_cpu":
			value := readClusterResource(iterator)
			object.computeNodesCpu = value
			object.bitmap_ |= 2
		case "compute_nodes_memory":
			value := readClusterResource(iterator)
			object.computeNodesMemory = value
			object.bitmap_ |= 4
		case "compute_nodes_sockets":
			value := readClusterResource(iterator)
			object.computeNodesSockets = value
			object.bitmap_ |= 8
		case "console_url":
			value := iterator.ReadString()
			object.consoleUrl = value
			object.bitmap_ |= 16
		case "cpu":
			value := readClusterResource(iterator)
			object.cpu = value
			object.bitmap_ |= 32
		case "critical_alerts_firing":
			value := iterator.ReadFloat64()
			object.criticalAlertsFiring = value
			object.bitmap_ |= 64
		case "health_state":
			value := iterator.ReadString()
			object.healthState = value
			object.bitmap_ |= 128
		case "memory":
			value := readClusterResource(iterator)
			object.memory = value
			object.bitmap_ |= 256
		case "nodes":
			value := readClusterMetricsNodes(iterator)
			object.nodes = value
			object.bitmap_ |= 512
		case "openshift_version":
			value := iterator.ReadString()
			object.openshiftVersion = value
			object.bitmap_ |= 1024
		case "operating_system":
			value := iterator.ReadString()
			object.operatingSystem = value
			object.bitmap_ |= 2048
		case "operators_condition_failing":
			value := iterator.ReadFloat64()
			object.operatorsConditionFailing = value
			object.bitmap_ |= 4096
		case "region":
			value := iterator.ReadString()
			object.region = value
			object.bitmap_ |= 8192
		case "sockets":
			value := readClusterResource(iterator)
			object.sockets = value
			object.bitmap_ |= 16384
		case "state":
			value := iterator.ReadString()
			object.state = value
			object.bitmap_ |= 32768
		case "state_description":
			value := iterator.ReadString()
			object.stateDescription = value
			object.bitmap_ |= 65536
		case "storage":
			value := readClusterResource(iterator)
			object.storage = value
			object.bitmap_ |= 131072
		case "subscription_cpu_total":
			value := iterator.ReadFloat64()
			object.subscriptionCpuTotal = value
			object.bitmap_ |= 262144
		case "subscription_obligation_exists":
			value := iterator.ReadFloat64()
			object.subscriptionObligationExists = value
			object.bitmap_ |= 524288
		case "subscription_socket_total":
			value := iterator.ReadFloat64()
			object.subscriptionSocketTotal = value
			object.bitmap_ |= 1048576
		case "upgrade":
			value := readClusterUpgrade(iterator)
			object.upgrade = value
			object.bitmap_ |= 2097152
		default:
			iterator.ReadAny()
		}
	}
	return object
}
