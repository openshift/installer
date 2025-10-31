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

package v1 // github.com/openshift-online/ocm-sdk-go/statusboard/v1

import "io"

func writeApplicationDependencyDeleteRequest(request *ApplicationDependencyDeleteRequest, writer io.Writer) error {
	return nil
}
func readApplicationDependencyDeleteResponse(response *ApplicationDependencyDeleteResponse, reader io.Reader) error {
	return nil
}
func writeApplicationDependencyGetRequest(request *ApplicationDependencyGetRequest, writer io.Writer) error {
	return nil
}
func readApplicationDependencyGetResponse(response *ApplicationDependencyGetResponse, reader io.Reader) error {
	var err error
	response.body, err = UnmarshalApplicationDependency(reader)
	return err
}
func writeApplicationDependencyUpdateRequest(request *ApplicationDependencyUpdateRequest, writer io.Writer) error {
	return MarshalApplicationDependency(request.body, writer)
}
func readApplicationDependencyUpdateResponse(response *ApplicationDependencyUpdateResponse, reader io.Reader) error {
	var err error
	response.body, err = UnmarshalApplicationDependency(reader)
	return err
}
