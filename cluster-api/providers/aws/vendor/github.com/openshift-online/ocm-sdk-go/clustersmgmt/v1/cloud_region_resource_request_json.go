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

func writeCloudRegionDeleteRequest(request *CloudRegionDeleteRequest, writer io.Writer) error {
	return nil
}
func readCloudRegionDeleteResponse(response *CloudRegionDeleteResponse, reader io.Reader) error {
	return nil
}
func writeCloudRegionGetRequest(request *CloudRegionGetRequest, writer io.Writer) error {
	return nil
}
func readCloudRegionGetResponse(response *CloudRegionGetResponse, reader io.Reader) error {
	var err error
	response.body, err = UnmarshalCloudRegion(reader)
	return err
}
func writeCloudRegionUpdateRequest(request *CloudRegionUpdateRequest, writer io.Writer) error {
	return MarshalCloudRegion(request.body, writer)
}
func readCloudRegionUpdateResponse(response *CloudRegionUpdateResponse, reader io.Reader) error {
	var err error
	response.body, err = UnmarshalCloudRegion(reader)
	return err
}
