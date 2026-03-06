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

	asocontainerservicev1hub "github.com/Azure/azure-service-operator/v2/api/containerservice/v1api20240901/storage"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"k8s.io/utils/ptr"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/aso"
)

const serviceName = "agentpools"

// AgentPoolScope defines the scope interface for an agent pool.
type AgentPoolScope interface {
	aso.Scope

	Name() string
	NodeResourceGroup() string
	AgentPoolSpec() azure.ASOResourceSpecGetter[genruntime.MetaObject]
	SetAgentPoolProviderIDList([]string)
	SetAgentPoolReplicas(int32)
	SetAgentPoolReady(bool)
	SetCAPIMachinePoolReplicas(replicas *int)
	SetCAPIMachinePoolAnnotation(key, value string)
	RemoveCAPIMachinePoolAnnotation(key string)
	SetSubnetName()
	IsPreviewEnabled() bool
}

// New creates a new service.
func New(scope AgentPoolScope) *aso.Service[genruntime.MetaObject, AgentPoolScope] {
	svc := aso.NewService[genruntime.MetaObject](serviceName, scope)
	svc.Specs = []azure.ASOResourceSpecGetter[genruntime.MetaObject]{scope.AgentPoolSpec()}
	svc.ConditionType = infrav1.AgentPoolsReadyCondition
	svc.PostCreateOrUpdateResourceHook = postCreateOrUpdateResourceHook
	return svc
}

func postCreateOrUpdateResourceHook(_ context.Context, scope AgentPoolScope, obj genruntime.MetaObject, err error) error {
	if err != nil {
		return err
	}
	agentPool := &asocontainerservicev1hub.ManagedClustersAgentPool{}
	if err := obj.(conversion.Convertible).ConvertTo(agentPool); err != nil {
		return err
	}

	// When autoscaling is set, add the annotation to the machine pool and update the replica count.
	if ptr.Deref(agentPool.Status.EnableAutoScaling, false) {
		scope.SetCAPIMachinePoolAnnotation(clusterv1.ReplicasManagedByAnnotation, "true")
		scope.SetCAPIMachinePoolReplicas(agentPool.Status.Count)
	} else { // Otherwise, remove the annotation.
		scope.RemoveCAPIMachinePoolAnnotation(clusterv1.ReplicasManagedByAnnotation)
	}
	return nil
}
