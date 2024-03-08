/*
Copyright 2020 The Kubernetes Authors.

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

// Package ssm provides a service to generate userdata for AWS Systems Manager.
package ssm

import (
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/internal/mime"
)

const (
	serviceID = "ssm"
)

// UserData creates a multi-part MIME document including a script boothook to
// download userdata from AWS Systems Manager and then restart cloud-init, and an include part
// specifying the on disk location of the new userdata.
func (s *Service) UserData(secretPrefix string, chunks int32, region string, endpoints []scope.ServiceEndpoint) ([]byte, error) {
	var serviceEndpoint = ""
	for _, v := range endpoints {
		if v.ServiceID == serviceID {
			serviceEndpoint = v.URL
		}
	}
	var userData, err = mime.GenerateInitDocument(secretPrefix, chunks, region, serviceEndpoint, secretFetchScript)
	if err != nil {
		return []byte{}, err
	}
	return userData, nil
}
