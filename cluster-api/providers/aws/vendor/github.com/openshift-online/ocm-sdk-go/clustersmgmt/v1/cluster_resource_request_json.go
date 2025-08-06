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

func writeClusterDeleteRequest(request *ClusterDeleteRequest, writer io.Writer) error {
	return nil
}
func readClusterDeleteResponse(response *ClusterDeleteResponse, reader io.Reader) error {
	return nil
}
func writeClusterGetRequest(request *ClusterGetRequest, writer io.Writer) error {
	return nil
}
func readClusterGetResponse(response *ClusterGetResponse, reader io.Reader) error {
	var err error
	response.body, err = UnmarshalCluster(reader)
	return err
}
func writeClusterHibernateRequest(request *ClusterHibernateRequest, writer io.Writer) error {
	return nil
}
func readClusterHibernateResponse(response *ClusterHibernateResponse, reader io.Reader) error {
	return nil
}
func writeClusterResumeRequest(request *ClusterResumeRequest, writer io.Writer) error {
	return nil
}
func readClusterResumeResponse(response *ClusterResumeResponse, reader io.Reader) error {
	return nil
}
func writeClusterUpdateRequest(request *ClusterUpdateRequest, writer io.Writer) error {
	return MarshalCluster(request.body, writer)
}
func readClusterUpdateResponse(response *ClusterUpdateResponse, reader io.Reader) error {
	var err error
	response.body, err = UnmarshalCluster(reader)
	return err
}
