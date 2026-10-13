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
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalSubscription writes a value of the 'subscription' type to the given writer.
func MarshalSubscription(object *Subscription, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteSubscription(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteSubscription writes a value of the 'subscription' type to the given stream.
func WriteSubscription(object *Subscription, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(SubscriptionLinkKind)
	} else {
		stream.WriteString(SubscriptionKind)
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
		stream.WriteObjectField("billing_marketplace_account")
		stream.WriteString(object.billingMarketplaceAccount)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4] && object.capabilities != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("capabilities")
		WriteCapabilityList(object.capabilities, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cloud_account_id")
		stream.WriteString(object.cloudAccountID)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cloud_provider_id")
		stream.WriteString(object.cloudProviderID)
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cluster_id")
		stream.WriteString(object.clusterID)
		count++
	}
	present_ = len(object.fieldSet_) > 8 && object.fieldSet_[8]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cluster_billing_model")
		stream.WriteString(string(object.clusterBillingModel))
		count++
	}
	present_ = len(object.fieldSet_) > 9 && object.fieldSet_[9]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("console_url")
		stream.WriteString(object.consoleURL)
		count++
	}
	present_ = len(object.fieldSet_) > 10 && object.fieldSet_[10]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("consumer_uuid")
		stream.WriteString(object.consumerUUID)
		count++
	}
	present_ = len(object.fieldSet_) > 11 && object.fieldSet_[11]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cpu_total")
		stream.WriteInt(object.cpuTotal)
		count++
	}
	present_ = len(object.fieldSet_) > 12 && object.fieldSet_[12]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("created_at")
		stream.WriteString((object.createdAt).Format(time.RFC3339))
		count++
	}
	present_ = len(object.fieldSet_) > 13 && object.fieldSet_[13] && object.creator != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("creator")
		WriteAccount(object.creator, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 14 && object.fieldSet_[14]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("display_name")
		stream.WriteString(object.displayName)
		count++
	}
	present_ = len(object.fieldSet_) > 15 && object.fieldSet_[15]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("external_cluster_id")
		stream.WriteString(object.externalClusterID)
		count++
	}
	present_ = len(object.fieldSet_) > 16 && object.fieldSet_[16] && object.labels != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("labels")
		WriteLabelList(object.labels, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 17 && object.fieldSet_[17]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("last_reconcile_date")
		stream.WriteString((object.lastReconcileDate).Format(time.RFC3339))
		count++
	}
	present_ = len(object.fieldSet_) > 18 && object.fieldSet_[18]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("last_released_at")
		stream.WriteString((object.lastReleasedAt).Format(time.RFC3339))
		count++
	}
	present_ = len(object.fieldSet_) > 19 && object.fieldSet_[19]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("last_telemetry_date")
		stream.WriteString((object.lastTelemetryDate).Format(time.RFC3339))
		count++
	}
	present_ = len(object.fieldSet_) > 20 && object.fieldSet_[20]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("managed")
		stream.WriteBool(object.managed)
		count++
	}
	present_ = len(object.fieldSet_) > 21 && object.fieldSet_[21] && object.metrics != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("metrics")
		WriteSubscriptionMetricsList(object.metrics, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 22 && object.fieldSet_[22] && object.notificationContacts != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("notification_contacts")
		WriteAccountList(object.notificationContacts, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 23 && object.fieldSet_[23]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("organization_id")
		stream.WriteString(object.organizationID)
		count++
	}
	present_ = len(object.fieldSet_) > 24 && object.fieldSet_[24] && object.plan != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("plan")
		WritePlan(object.plan, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 25 && object.fieldSet_[25]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("product_bundle")
		stream.WriteString(object.productBundle)
		count++
	}
	present_ = len(object.fieldSet_) > 26 && object.fieldSet_[26]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("provenance")
		stream.WriteString(object.provenance)
		count++
	}
	present_ = len(object.fieldSet_) > 27 && object.fieldSet_[27]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("region_id")
		stream.WriteString(object.regionID)
		count++
	}
	present_ = len(object.fieldSet_) > 28 && object.fieldSet_[28]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("released")
		stream.WriteBool(object.released)
		count++
	}
	present_ = len(object.fieldSet_) > 29 && object.fieldSet_[29]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("rh_region_id")
		stream.WriteString(object.rhRegionID)
		count++
	}
	present_ = len(object.fieldSet_) > 30 && object.fieldSet_[30]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("service_level")
		stream.WriteString(object.serviceLevel)
		count++
	}
	present_ = len(object.fieldSet_) > 31 && object.fieldSet_[31]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("socket_total")
		stream.WriteInt(object.socketTotal)
		count++
	}
	present_ = len(object.fieldSet_) > 32 && object.fieldSet_[32]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("status")
		stream.WriteString(object.status)
		count++
	}
	present_ = len(object.fieldSet_) > 33 && object.fieldSet_[33]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("support_level")
		stream.WriteString(object.supportLevel)
		count++
	}
	present_ = len(object.fieldSet_) > 34 && object.fieldSet_[34]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("system_units")
		stream.WriteString(object.systemUnits)
		count++
	}
	present_ = len(object.fieldSet_) > 35 && object.fieldSet_[35]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("trial_end_date")
		stream.WriteString((object.trialEndDate).Format(time.RFC3339))
		count++
	}
	present_ = len(object.fieldSet_) > 36 && object.fieldSet_[36]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("updated_at")
		stream.WriteString((object.updatedAt).Format(time.RFC3339))
		count++
	}
	present_ = len(object.fieldSet_) > 37 && object.fieldSet_[37]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("usage")
		stream.WriteString(object.usage)
	}
	stream.WriteObjectEnd()
}

// UnmarshalSubscription reads a value of the 'subscription' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalSubscription(source interface{}) (object *Subscription, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadSubscription(iterator)
	err = iterator.Error
	return
}

// ReadSubscription reads a value of the 'subscription' type from the given iterator.
func ReadSubscription(iterator *jsoniter.Iterator) *Subscription {
	object := &Subscription{
		fieldSet_: make([]bool, 38),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == SubscriptionLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "billing_marketplace_account":
			value := iterator.ReadString()
			object.billingMarketplaceAccount = value
			object.fieldSet_[3] = true
		case "capabilities":
			value := ReadCapabilityList(iterator)
			object.capabilities = value
			object.fieldSet_[4] = true
		case "cloud_account_id":
			value := iterator.ReadString()
			object.cloudAccountID = value
			object.fieldSet_[5] = true
		case "cloud_provider_id":
			value := iterator.ReadString()
			object.cloudProviderID = value
			object.fieldSet_[6] = true
		case "cluster_id":
			value := iterator.ReadString()
			object.clusterID = value
			object.fieldSet_[7] = true
		case "cluster_billing_model":
			text := iterator.ReadString()
			value := BillingModel(text)
			object.clusterBillingModel = value
			object.fieldSet_[8] = true
		case "console_url":
			value := iterator.ReadString()
			object.consoleURL = value
			object.fieldSet_[9] = true
		case "consumer_uuid":
			value := iterator.ReadString()
			object.consumerUUID = value
			object.fieldSet_[10] = true
		case "cpu_total":
			value := iterator.ReadInt()
			object.cpuTotal = value
			object.fieldSet_[11] = true
		case "created_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.createdAt = value
			object.fieldSet_[12] = true
		case "creator":
			value := ReadAccount(iterator)
			object.creator = value
			object.fieldSet_[13] = true
		case "display_name":
			value := iterator.ReadString()
			object.displayName = value
			object.fieldSet_[14] = true
		case "external_cluster_id":
			value := iterator.ReadString()
			object.externalClusterID = value
			object.fieldSet_[15] = true
		case "labels":
			value := ReadLabelList(iterator)
			object.labels = value
			object.fieldSet_[16] = true
		case "last_reconcile_date":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.lastReconcileDate = value
			object.fieldSet_[17] = true
		case "last_released_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.lastReleasedAt = value
			object.fieldSet_[18] = true
		case "last_telemetry_date":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.lastTelemetryDate = value
			object.fieldSet_[19] = true
		case "managed":
			value := iterator.ReadBool()
			object.managed = value
			object.fieldSet_[20] = true
		case "metrics":
			value := ReadSubscriptionMetricsList(iterator)
			object.metrics = value
			object.fieldSet_[21] = true
		case "notification_contacts":
			value := ReadAccountList(iterator)
			object.notificationContacts = value
			object.fieldSet_[22] = true
		case "organization_id":
			value := iterator.ReadString()
			object.organizationID = value
			object.fieldSet_[23] = true
		case "plan":
			value := ReadPlan(iterator)
			object.plan = value
			object.fieldSet_[24] = true
		case "product_bundle":
			value := iterator.ReadString()
			object.productBundle = value
			object.fieldSet_[25] = true
		case "provenance":
			value := iterator.ReadString()
			object.provenance = value
			object.fieldSet_[26] = true
		case "region_id":
			value := iterator.ReadString()
			object.regionID = value
			object.fieldSet_[27] = true
		case "released":
			value := iterator.ReadBool()
			object.released = value
			object.fieldSet_[28] = true
		case "rh_region_id":
			value := iterator.ReadString()
			object.rhRegionID = value
			object.fieldSet_[29] = true
		case "service_level":
			value := iterator.ReadString()
			object.serviceLevel = value
			object.fieldSet_[30] = true
		case "socket_total":
			value := iterator.ReadInt()
			object.socketTotal = value
			object.fieldSet_[31] = true
		case "status":
			value := iterator.ReadString()
			object.status = value
			object.fieldSet_[32] = true
		case "support_level":
			value := iterator.ReadString()
			object.supportLevel = value
			object.fieldSet_[33] = true
		case "system_units":
			value := iterator.ReadString()
			object.systemUnits = value
			object.fieldSet_[34] = true
		case "trial_end_date":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.trialEndDate = value
			object.fieldSet_[35] = true
		case "updated_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.updatedAt = value
			object.fieldSet_[36] = true
		case "usage":
			value := iterator.ReadString()
			object.usage = value
			object.fieldSet_[37] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
