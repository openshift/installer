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

package agentpools

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v4"
	"github.com/pkg/errors"
	"k8s.io/utils/ptr"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/async"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

const serviceName = "agentpools"

// AgentPoolScope defines the scope interface for an agent pool.
type AgentPoolScope interface {
	azure.ClusterDescriber
	azure.AsyncStatusUpdater

	Name() string
	NodeResourceGroup() string
	AgentPoolAnnotations() map[string]string
	AgentPoolSpec() azure.ResourceSpecGetter
	SetAgentPoolProviderIDList([]string)
	SetAgentPoolReplicas(int32)
	SetAgentPoolReady(bool)
	SetCAPIMachinePoolReplicas(replicas *int32)
	SetCAPIMachinePoolAnnotation(key, value string)
	RemoveCAPIMachinePoolAnnotation(key string)
	SetSubnetName()
}

// Service provides operations on Azure resources.
type Service struct {
	scope AgentPoolScope
	async.Reconciler
}

// New creates a new service.
func New(scope AgentPoolScope) (*Service, error) {
	client, err := newClient(scope)
	if err != nil {
		return nil, err
	}
	return &Service{
		scope: scope,
		Reconciler: async.New[armcontainerservice.AgentPoolsClientCreateOrUpdateResponse,
			armcontainerservice.AgentPoolsClientDeleteResponse](scope, client, client),
	}, nil
}

// Name returns the service name.
func (s *Service) Name() string {
	return serviceName
}

// Reconcile idempotently creates or updates an agent pool, if possible.
func (s *Service) Reconcile(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "agentpools.Service.Reconcile")
	defer done()

	var resultingErr error
	if agentPoolSpec := s.scope.AgentPoolSpec(); agentPoolSpec != nil {
		result, err := s.CreateOrUpdateResource(ctx, agentPoolSpec, serviceName)
		if err != nil {
			resultingErr = err
		} else {
			agentPool, ok := result.(armcontainerservice.AgentPool)
			if !ok {
				return errors.Errorf("%T is not an armcontainerservice.AgentPool", result)
			}
			// When autoscaling is set, add the annotation to the machine pool and update the replica count.
			if ptr.Deref(agentPool.Properties.EnableAutoScaling, false) {
				s.scope.SetCAPIMachinePoolAnnotation(clusterv1.ReplicasManagedByAnnotation, "true")
				s.scope.SetCAPIMachinePoolReplicas(agentPool.Properties.Count)
			} else { // Otherwise, remove the annotation.
				s.scope.RemoveCAPIMachinePoolAnnotation(clusterv1.ReplicasManagedByAnnotation)
			}
		}
	} else {
		return nil
	}

	s.scope.UpdatePutStatus(infrav1.AgentPoolsReadyCondition, serviceName, resultingErr)
	return resultingErr
}

// Delete deletes the virtual network with the provided name.
func (s *Service) Delete(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "agentpools.Service.Delete")
	defer done()

	var resultingErr error
	if agentPoolSpec := s.scope.AgentPoolSpec(); agentPoolSpec != nil {
		resultingErr = s.DeleteResource(ctx, agentPoolSpec, serviceName)
	} else {
		return nil
	}

	s.scope.UpdateDeleteStatus(infrav1.AgentPoolsReadyCondition, serviceName, resultingErr)
	return resultingErr
}
