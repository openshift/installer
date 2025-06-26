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

// Package endpoints contains aws endpoint related utilities.
package endpoints

import (
	"errors"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws/endpoints"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/endpointsv2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
)

var (
	errServiceEndpointFormat             = errors.New("must be formatted as ${ServiceID}=${URL}")
	errServiceEndpointSigningRegion      = errors.New("must be formatted as ${SigningRegion}:${ServiceID1}=${URL1},${ServiceID2}=${URL2...}")
	errServiceEndpointURL                = errors.New("must use a valid URL as a service-endpoint")
	errServiceEndpointServiceID          = errors.New("must use a valid serviceID from the AWS GO SDK")
	errServiceEndpointDuplicateServiceID = errors.New("same serviceID defined twice for signing region")
)

func serviceEnum() []string {
	var serviceIDs = []string{}
	resolver := endpoints.DefaultResolver()
	partitions := resolver.(endpoints.EnumPartitions).Partitions()
	for _, p := range partitions {
		for id := range p.Services() {
			var add = true
			for _, s := range serviceIDs {
				if id == s {
					add = false
				}
			}
			if add {
				serviceIDs = append(serviceIDs, id)
			}
		}
	}

	return serviceIDs
}

// ParseFlag parses the command line flag of service endponts in the format ${SigningRegion1}:${ServiceID1}=${URL1},${ServiceID2}=${URL2}...;${SigningRegion2}...
// returning a set of ServiceEndpoints.
func ParseFlag(serviceEndpoints string) ([]scope.ServiceEndpoint, error) {
	if serviceEndpoints == "" {
		return nil, nil
	}
	serviceIDs := serviceEnum()
	signingRegionConfigs := strings.Split(serviceEndpoints, ";")
	endpoints := []scope.ServiceEndpoint{}
	for _, regionConfig := range signingRegionConfigs {
		components := strings.SplitN(regionConfig, ":", 2)
		if len(components) != 2 {
			return nil, errServiceEndpointSigningRegion
		}
		signingRegion := components[0]
		servicePairs := strings.Split(components[1], ",")
		seenServices := []string{}
		for _, servicePair := range servicePairs {
			kv := strings.Split(servicePair, "=")
			if len(kv) != 2 {
				return nil, errServiceEndpointFormat
			}
			var serviceID = ""
			for _, id := range serviceIDs {
				if kv[0] == id {
					serviceID = kv[0]

					break
				}
			}
			if serviceID == "" {
				return nil, errServiceEndpointServiceID
			}
			if containsString(seenServices, serviceID) {
				return nil, errServiceEndpointDuplicateServiceID
			}
			seenServices = append(seenServices, serviceID)
			URL, err := url.ParseRequestURI(kv[1])
			if err != nil {
				return nil, errServiceEndpointURL
			}
			endpoints = append(endpoints, scope.ServiceEndpoint{
				ServiceID:     serviceID,
				URL:           URL.String(),
				SigningRegion: signingRegion,
			})
		}
	}
	// For Go SDK V2 migration
	saveToServiceEndpointV2Map(endpoints)

	return endpoints, nil
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}

	return false
}

// TODO: punkwalker - remove this after Go SDK V2 migration.
func saveToServiceEndpointV2Map(src []scope.ServiceEndpoint) {
	for _, svc := range src {
		// convert service ID to UpperCase as service IDs in AWS SDK GO V2 are UpperCase & Go map is Case Sensitve
		// This is for backward compabitibility
		// Ref: SDK V2 https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/ec2#pkg-constants
		// Ref: SDK V1 https://pkg.go.dev/github.com/aws/aws-sdk-go/aws/endpoints#pkg-constants
		serviceID := strings.ToUpper(svc.ServiceID)
		endpoint := endpointsv2.ServiceEndpoint{
			ServiceID:     serviceID,
			URL:           svc.URL,
			SigningRegion: svc.SigningRegion,
		}
		endpointsv2.ServiceEndpointsMap[svc.ServiceID] = endpoint
	}
}
