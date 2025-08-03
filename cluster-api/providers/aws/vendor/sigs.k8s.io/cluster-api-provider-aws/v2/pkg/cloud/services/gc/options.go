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

package gc

import (
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/elb"
)

// ServiceOption is an option for creating the service.
type ServiceOption func(*Service)

// withELBClient is an option for specifying a AWS ELB Client.
func withELBClient(client elb.ELBAPI) ServiceOption {
	return func(s *Service) {
		s.elbClient = client
	}
}

// withELBv2Client is an option for specifying a AWS ELBv2 Client.
func withELBv2Client(client elb.ELBV2API) ServiceOption {
	return func(s *Service) {
		s.elbv2Client = client
	}
}

// withResourceTaggingClient is an option for specifying a AWS Resource Tagging Client.
func withResourceTaggingClient(client elb.ResourceGroupsTaggingAPIAPI) ServiceOption {
	return func(s *Service) {
		s.resourceTaggingClient = client
	}
}

// withEC2Client is an option for specifying a AWS EC2 Client.
func withEC2Client(client ec2iface.EC2API) ServiceOption {
	return func(s *Service) {
		s.ec2Client = client
	}
}

// WithGCStrategy is an option for specifying using the alternative GC strategy.
func WithGCStrategy(alternativeGCStrategy bool) ServiceOption {
	if alternativeGCStrategy {
		return func(s *Service) {
			addAlternativeCollectFuncs(s)
		}
	}
	return func(s *Service) {
		addDefaultCollectFuncs(s)
	}
}
