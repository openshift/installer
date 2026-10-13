/*
Copyright 2022 The Kubernetes Authors.

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

package endpoints

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
)

// ServiceEndpointFormat is used to identify the custom endpoint of IBM Cloud service.
var ServiceEndpointFormat string

const (
	// VPC used to identify VPC service.
	VPC serviceID = "vpc"
	// PowerVS used to identify PowerVS service.
	PowerVS serviceID = "powervs"
	// RC used to identify Resource-Controller service.
	RC serviceID = "rc"
	// TransitGateway  service.
	TransitGateway serviceID = "transitgateway"
	// COS service.
	COS serviceID = "cos"
	// RM used to identify Resource-Manager service.
	RM serviceID = "rm"
	// GlobalTagging used to identify the Global Tagging service.
	GlobalTagging serviceID = "globaltagging"
)

type serviceID string

var serviceIDs = []serviceID{VPC, PowerVS, RC, TransitGateway, COS, RM, GlobalTagging}

// ServiceEndpoint holds the Service endpoint specific information.
type ServiceEndpoint struct {
	ID     string
	URL    string
	Region string
}

var (
	errServiceEndpointFormat      = errors.New("must be formatted as ${ServiceID}=${URL}")
	errServiceEndpointRegion      = errors.New("must be formatted as ${ServiceRegion}:${ServiceID1}=${URL1},${ServiceID2}=${URL2...}")
	errServiceEndpointURL         = errors.New("must use a valid URL as a service-endpoint")
	errServiceEndpointID          = errors.New("invalid service ID: %s must use a valid")
	errServiceEndpointDuplicateID = errors.New("same ID defined twice for service region")
)

// ParseServiceEndpointFlag parses the command line flag of service endpoint in the format ${ServiceRegion}:${ServiceID1}=${URL1},${ServiceID2}=${URL2...}
// returning a list of ServiceEndpoint.
func ParseServiceEndpointFlag(serviceEndpoints string) ([]ServiceEndpoint, error) {
	if serviceEndpoints == "" || serviceEndpoints == "none" {
		return nil, nil
	}

	serviceRegionConfigs := strings.Split(serviceEndpoints, ";")
	endpoints := []ServiceEndpoint{}
	for _, regionConfig := range serviceRegionConfigs {
		components := strings.SplitN(regionConfig, ":", 2)
		if len(components) != 2 {
			return nil, errServiceEndpointRegion
		}
		serviceRegion := components[0]
		servicePairs := strings.Split(components[1], ",")
		seenServices := []string{}
		for _, servicePair := range servicePairs {
			kv := strings.Split(servicePair, "=")
			if len(kv) != 2 {
				return nil, errServiceEndpointFormat
			}
			var serviceID = ""
			for _, id := range serviceIDs {
				if kv[0] == string(id) {
					serviceID = kv[0]
					break
				}
			}
			if serviceID == "" {
				return nil, errServiceEndpointID
			}
			if containsString(seenServices, serviceID) {
				return nil, errServiceEndpointDuplicateID
			}
			seenServices = append(seenServices, serviceID)
			URL, err := url.ParseRequestURI(kv[1])
			if err != nil {
				return nil, errServiceEndpointURL
			}
			endpoints = append(endpoints, ServiceEndpoint{
				ID:     serviceID,
				URL:    URL.String(),
				Region: serviceRegion,
			})
		}
	}
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

// FetchVPCEndpoint will return VPC service endpoint.
func FetchVPCEndpoint(region string, serviceEndpoint []ServiceEndpoint) string {
	svcEndpoint := "https://" + region + ".iaas.cloud.ibm.com/v1"
	for _, vpcEndpoint := range serviceEndpoint {
		if vpcEndpoint.Region == region && vpcEndpoint.ID == string(VPC) {
			return vpcEndpoint.URL
		}
	}
	return svcEndpoint
}

// FetchPVSEndpoint will return PowerVS service endpoint.
// Deprecated: User FetchEndpoints instead.
func FetchPVSEndpoint(region string, serviceEndpoint []ServiceEndpoint) string {
	for _, powervsEndpoint := range serviceEndpoint {
		if powervsEndpoint.Region == region && powervsEndpoint.ID == string(PowerVS) {
			return powervsEndpoint.URL
		}
	}
	return ""
}

// FetchRCEndpoint will return resource controller endpoint.
// Deprecated: User FetchEndpoints instead.
func FetchRCEndpoint(serviceEndpoint []ServiceEndpoint) string {
	for _, rcEndpoint := range serviceEndpoint {
		if rcEndpoint.ID == string(RC) {
			return rcEndpoint.URL
		}
	}
	return ""
}

// FetchEndpoints returns the endpoint associated with serviceID otherwise empty string.
func FetchEndpoints(serviceID string, serviceEndpoint []ServiceEndpoint) string {
	for _, endpoint := range serviceEndpoint {
		if endpoint.ID == serviceID {
			return endpoint.URL
		}
	}
	return ""
}

// ConstructRegionFromZone Calculate region based on location/zone.
func ConstructRegionFromZone(zone string) string {
	var regex string
	if strings.Contains(zone, "-") {
		// it's a region or AZ
		regex = "-[0-9]+$"
	} else {
		// it's a datacenter
		regex = "[0-9]+$"
	}

	reg, _ := regexp.Compile(regex)
	return reg.ReplaceAllString(zone, "")
}
