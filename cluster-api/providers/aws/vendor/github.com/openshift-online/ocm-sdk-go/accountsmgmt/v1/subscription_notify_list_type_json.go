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

// MarshalSubscriptionNotifyList writes a list of values of the 'subscription_notify' type to
// the given writer.
func MarshalSubscriptionNotifyList(list []*SubscriptionNotify, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeSubscriptionNotifyList(list, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeSubscriptionNotifyList writes a list of value of the 'subscription_notify' type to
// the given stream.
func writeSubscriptionNotifyList(list []*SubscriptionNotify, stream *jsoniter.Stream) {
	stream.WriteArrayStart()
	for i, value := range list {
		if i > 0 {
			stream.WriteMore()
		}
		writeSubscriptionNotify(value, stream)
	}
	stream.WriteArrayEnd()
}

// UnmarshalSubscriptionNotifyList reads a list of values of the 'subscription_notify' type
// from the given source, which can be a slice of bytes, a string or a reader.
func UnmarshalSubscriptionNotifyList(source interface{}) (items []*SubscriptionNotify, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	items = readSubscriptionNotifyList(iterator)
	err = iterator.Error
	return
}

// readSubscriptionNotifyList reads list of values of the ”subscription_notify' type from
// the given iterator.
func readSubscriptionNotifyList(iterator *jsoniter.Iterator) []*SubscriptionNotify {
	list := []*SubscriptionNotify{}
	for iterator.ReadArray() {
		item := readSubscriptionNotify(iterator)
		list = append(list, item)
	}
	return list
}
