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

func writePeerDependencyDeleteRequest(request *PeerDependencyDeleteRequest, writer io.Writer) error {
	return nil
}
func readPeerDependencyDeleteResponse(response *PeerDependencyDeleteResponse, reader io.Reader) error {
	return nil
}
func writePeerDependencyGetRequest(request *PeerDependencyGetRequest, writer io.Writer) error {
	return nil
}
func readPeerDependencyGetResponse(response *PeerDependencyGetResponse, reader io.Reader) error {
	var err error
	response.body, err = UnmarshalService(reader)
	return err
}
func writePeerDependencyUpdateRequest(request *PeerDependencyUpdateRequest, writer io.Writer) error {
	return MarshalPeerDependency(request.body, writer)
}
func readPeerDependencyUpdateResponse(response *PeerDependencyUpdateResponse, reader io.Reader) error {
	var err error
	response.body, err = UnmarshalPeerDependency(reader)
	return err
}
