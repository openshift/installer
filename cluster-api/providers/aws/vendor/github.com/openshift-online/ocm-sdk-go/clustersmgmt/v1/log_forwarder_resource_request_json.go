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

import "io"

func writeLogForwarderDeleteRequest(request *LogForwarderDeleteRequest, writer io.Writer) error {
	return nil
}
func readLogForwarderDeleteResponse(response *LogForwarderDeleteResponse, reader io.Reader) error {
	return nil
}
func writeLogForwarderGetRequest(request *LogForwarderGetRequest, writer io.Writer) error {
	return nil
}
func readLogForwarderGetResponse(response *LogForwarderGetResponse, reader io.Reader) error {
	var err error
	response.body, err = UnmarshalLogForwarder(reader)
	return err
}
func writeLogForwarderUpdateRequest(request *LogForwarderUpdateRequest, writer io.Writer) error {
	return MarshalLogForwarder(request.body, writer)
}
func readLogForwarderUpdateResponse(response *LogForwarderUpdateResponse, reader io.Reader) error {
	var err error
	response.body, err = UnmarshalLogForwarder(reader)
	return err
}
