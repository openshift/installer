/*
Copyright 2023 The Kubernetes Authors.

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
	"fmt"

	"github.com/aws/aws-sdk-go/aws/arn"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
)

const (
	fakePartition     = "aws"
	fakeRegion        = "fake-region"
	fakeAccount       = "fake-account"
	elbService        = "elasticloadbalancing"
	elbResourcePrefix = "loadbalancer/"
	sgService         = "ec2"
	sgResourcePrefix  = "security-group/"

	// maxDescribeTagsRequest is the maximum number of resources for the DescribeTags API call
	// see: https://docs.aws.amazon.com/elasticloadbalancing/latest/APIReference/API_DescribeTags.html.
	maxDescribeTagsRequest = 20
)

// composeFakeArn composes a resource arn with correct service and resource, but fake partition, region and account.
// This fake arn is used to compose an *AWSResource object that can be consumed by existing cleanupFuncs of gc service.
func composeFakeArn(service, resource string) string {
	return "arn:" + fakePartition + ":" + service + ":" + fakeRegion + ":" + fakeAccount + ":" + resource
}

// composeAWSResource composes *AWSResource object for an aws resource.
func composeAWSResource(resourceARN string, resourceTags infrav1.Tags) (*AWSResource, error) {
	parsedArn, err := arn.Parse(resourceARN)
	if err != nil {
		return nil, fmt.Errorf("parsing resource arn %s: %w", resourceARN, err)
	}

	resource := &AWSResource{
		ARN:  &parsedArn,
		Tags: resourceTags,
	}

	return resource, nil
}
