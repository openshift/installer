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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/accountsmgmt/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalSubscriptionMetrics writes a value of the 'subscription_metrics' type to the given writer.
func MarshalSubscriptionMetrics(object *SubscriptionMetrics, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteSubscriptionMetrics(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteSubscriptionMetrics writes a value of the 'subscription_metrics' type to the given stream.
func WriteSubscriptionMetrics(object *SubscriptionMetrics, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cloud_provider")
		stream.WriteString(object.cloudProvider)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1] && object.computeNodesCpu != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("compute_nodes_cpu")
		WriteClusterResource(object.computeNodesCpu, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2] && object.computeNodesMemory != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("compute_nodes_memory")
		WriteClusterResource(object.computeNodesMemory, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3] && object.computeNodesSockets != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("compute_nodes_sockets")
		WriteClusterResource(object.computeNodesSockets, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("console_url")
		stream.WriteString(object.consoleUrl)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5] && object.cpu != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cpu")
		WriteClusterResource(object.cpu, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("critical_alerts_firing")
		stream.WriteFloat64(object.criticalAlertsFiring)
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("health_state")
		stream.WriteString(object.healthState)
		count++
	}
	present_ = len(object.fieldSet_) > 8 && object.fieldSet_[8] && object.memory != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("memory")
		WriteClusterResource(object.memory, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 9 && object.fieldSet_[9] && object.nodes != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("nodes")
		WriteClusterMetricsNodes(object.nodes, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 10 && object.fieldSet_[10]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("openshift_version")
		stream.WriteString(object.openshiftVersion)
		count++
	}
	present_ = len(object.fieldSet_) > 11 && object.fieldSet_[11]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("operating_system")
		stream.WriteString(object.operatingSystem)
		count++
	}
	present_ = len(object.fieldSet_) > 12 && object.fieldSet_[12]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("operators_condition_failing")
		stream.WriteFloat64(object.operatorsConditionFailing)
		count++
	}
	present_ = len(object.fieldSet_) > 13 && object.fieldSet_[13]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("region")
		stream.WriteString(object.region)
		count++
	}
	present_ = len(object.fieldSet_) > 14 && object.fieldSet_[14] && object.sockets != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("sockets")
		WriteClusterResource(object.sockets, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 15 && object.fieldSet_[15]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("state")
		stream.WriteString(object.state)
		count++
	}
	present_ = len(object.fieldSet_) > 16 && object.fieldSet_[16]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("state_description")
		stream.WriteString(object.stateDescription)
		count++
	}
	present_ = len(object.fieldSet_) > 17 && object.fieldSet_[17] && object.storage != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("storage")
		WriteClusterResource(object.storage, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 18 && object.fieldSet_[18]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subscription_cpu_total")
		stream.WriteFloat64(object.subscriptionCpuTotal)
		count++
	}
	present_ = len(object.fieldSet_) > 19 && object.fieldSet_[19]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subscription_obligation_exists")
		stream.WriteFloat64(object.subscriptionObligationExists)
		count++
	}
	present_ = len(object.fieldSet_) > 20 && object.fieldSet_[20]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subscription_socket_total")
		stream.WriteFloat64(object.subscriptionSocketTotal)
		count++
	}
	present_ = len(object.fieldSet_) > 21 && object.fieldSet_[21] && object.upgrade != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("upgrade")
		WriteClusterUpgrade(object.upgrade, stream)
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
	object = ReadSubscriptionMetrics(iterator)
	err = iterator.Error
	return
}

// ReadSubscriptionMetrics reads a value of the 'subscription_metrics' type from the given iterator.
func ReadSubscriptionMetrics(iterator *jsoniter.Iterator) *SubscriptionMetrics {
	object := &SubscriptionMetrics{
		fieldSet_: make([]bool, 22),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "cloud_provider":
			value := iterator.ReadString()
			object.cloudProvider = value
			object.fieldSet_[0] = true
		case "compute_nodes_cpu":
			value := ReadClusterResource(iterator)
			object.computeNodesCpu = value
			object.fieldSet_[1] = true
		case "compute_nodes_memory":
			value := ReadClusterResource(iterator)
			object.computeNodesMemory = value
			object.fieldSet_[2] = true
		case "compute_nodes_sockets":
			value := ReadClusterResource(iterator)
			object.computeNodesSockets = value
			object.fieldSet_[3] = true
		case "console_url":
			value := iterator.ReadString()
			object.consoleUrl = value
			object.fieldSet_[4] = true
		case "cpu":
			value := ReadClusterResource(iterator)
			object.cpu = value
			object.fieldSet_[5] = true
		case "critical_alerts_firing":
			value := iterator.ReadFloat64()
			object.criticalAlertsFiring = value
			object.fieldSet_[6] = true
		case "health_state":
			value := iterator.ReadString()
			object.healthState = value
			object.fieldSet_[7] = true
		case "memory":
			value := ReadClusterResource(iterator)
			object.memory = value
			object.fieldSet_[8] = true
		case "nodes":
			value := ReadClusterMetricsNodes(iterator)
			object.nodes = value
			object.fieldSet_[9] = true
		case "openshift_version":
			value := iterator.ReadString()
			object.openshiftVersion = value
			object.fieldSet_[10] = true
		case "operating_system":
			value := iterator.ReadString()
			object.operatingSystem = value
			object.fieldSet_[11] = true
		case "operators_condition_failing":
			value := iterator.ReadFloat64()
			object.operatorsConditionFailing = value
			object.fieldSet_[12] = true
		case "region":
			value := iterator.ReadString()
			object.region = value
			object.fieldSet_[13] = true
		case "sockets":
			value := ReadClusterResource(iterator)
			object.sockets = value
			object.fieldSet_[14] = true
		case "state":
			value := iterator.ReadString()
			object.state = value
			object.fieldSet_[15] = true
		case "state_description":
			value := iterator.ReadString()
			object.stateDescription = value
			object.fieldSet_[16] = true
		case "storage":
			value := ReadClusterResource(iterator)
			object.storage = value
			object.fieldSet_[17] = true
		case "subscription_cpu_total":
			value := iterator.ReadFloat64()
			object.subscriptionCpuTotal = value
			object.fieldSet_[18] = true
		case "subscription_obligation_exists":
			value := iterator.ReadFloat64()
			object.subscriptionObligationExists = value
			object.fieldSet_[19] = true
		case "subscription_socket_total":
			value := iterator.ReadFloat64()
			object.subscriptionSocketTotal = value
			object.fieldSet_[20] = true
		case "upgrade":
			value := ReadClusterUpgrade(iterator)
			object.upgrade = value
			object.fieldSet_[21] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
