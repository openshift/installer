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

package v1alpha1 // github.com/openshift-online/ocm-sdk-go/arohcp/v1alpha1

import "io"

func writeExternalAuthAsyncDeleteRequest(request *ExternalAuthDeleteRequest, writer io.Writer) error {
	return nil
}
func readExternalAuthAsyncDeleteResponse(response *ExternalAuthDeleteResponse, reader io.Reader) error {
	return nil
}
func writeExternalAuthAsyncUpdateRequest(request *ExternalAuthUpdateRequest, writer io.Writer) error {
	return MarshalExternalAuth(request.body, writer)
}
func readExternalAuthAsyncUpdateResponse(response *ExternalAuthUpdateResponse, reader io.Reader) error {
	var err error
	response.body, err = UnmarshalExternalAuth(reader)
	return err
}
func writeExternalAuthGetRequest(request *ExternalAuthGetRequest, writer io.Writer) error {
	return nil
}
func readExternalAuthGetResponse(response *ExternalAuthGetResponse, reader io.Reader) error {
	var err error
	response.body, err = UnmarshalExternalAuth(reader)
	return err
}
