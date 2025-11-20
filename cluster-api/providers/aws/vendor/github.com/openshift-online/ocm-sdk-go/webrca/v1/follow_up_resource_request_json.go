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

package v1 // github.com/openshift-online/ocm-sdk-go/webrca/v1

import "io"

func writeFollowUpDeleteRequest(request *FollowUpDeleteRequest, writer io.Writer) error {
	return nil
}
func readFollowUpDeleteResponse(response *FollowUpDeleteResponse, reader io.Reader) error {
	return nil
}
func writeFollowUpGetRequest(request *FollowUpGetRequest, writer io.Writer) error {
	return nil
}
func readFollowUpGetResponse(response *FollowUpGetResponse, reader io.Reader) error {
	var err error
	response.body, err = UnmarshalFollowUp(reader)
	return err
}
func writeFollowUpUpdateRequest(request *FollowUpUpdateRequest, writer io.Writer) error {
	return MarshalFollowUp(request.body, writer)
}
func readFollowUpUpdateResponse(response *FollowUpUpdateResponse, reader io.Reader) error {
	var err error
	response.body, err = UnmarshalFollowUp(reader)
	return err
}
