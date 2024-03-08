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

// Package instancestate provides a way to interact with the EC2 instance state.
package instancestate

import (
	"github.com/aws/aws-sdk-go/service/eventbridge/eventbridgeiface"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
)

// Service defines the specs for a service.
type Service struct {
	scope             scope.EC2Scope
	EventBridgeClient eventbridgeiface.EventBridgeAPI
	SQSClient         sqsiface.SQSAPI
}

// NewService returns a new service given the ec2 api client.
func NewService(clusterScope scope.EC2Scope) *Service {
	return &Service{
		scope:             clusterScope,
		EventBridgeClient: scope.NewEventBridgeClient(clusterScope, clusterScope, clusterScope.InfraCluster()),
		SQSClient:         scope.NewSQSClient(clusterScope, clusterScope, clusterScope.InfraCluster()),
	}
}
