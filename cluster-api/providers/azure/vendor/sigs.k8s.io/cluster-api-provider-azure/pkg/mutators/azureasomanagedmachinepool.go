/*
Copyright 2024 The Kubernetes Authors.

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

package mutators

import (
	"context"
	"fmt"
	"strings"

	asocontainerservicev1 "github.com/Azure/azure-service-operator/v2/api/containerservice/v1api20231001"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	infrav1alpha "sigs.k8s.io/cluster-api-provider-azure/api/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	expv1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ErrNoManagedClustersAgentPoolDefined describes an AzureASOManagedMachinePool without a ManagedClustersAgentPool.
var ErrNoManagedClustersAgentPoolDefined = fmt.Errorf("no %s ManagedClustersAgentPools defined in AzureASOManagedMachinePool spec.resources", asocontainerservicev1.GroupVersion.Group)

// SetAgentPoolDefaults propagates config from a MachinePool to an AzureASOManagedMachinePool's defined ManagedClustersAgentPool.
func SetAgentPoolDefaults(ctrlClient client.Client, machinePool *expv1.MachinePool) ResourcesMutator {
	return func(ctx context.Context, us []*unstructured.Unstructured) error {
		ctx, _, done := tele.StartSpanWithLogger(ctx, "mutators.SetAgentPoolDefaults")
		defer done()

		var agentPool *unstructured.Unstructured
		var agentPoolPath string
		for i, u := range us {
			if u.GroupVersionKind().Group == asocontainerservicev1.GroupVersion.Group &&
				u.GroupVersionKind().Kind == "ManagedClustersAgentPool" {
				agentPool = u
				agentPoolPath = fmt.Sprintf("spec.resources[%d]", i)
				break
			}
		}
		if agentPool == nil {
			return reconcile.TerminalError(ErrNoManagedClustersAgentPoolDefined)
		}

		if err := setAgentPoolOrchestratorVersion(ctx, machinePool, agentPoolPath, agentPool); err != nil {
			return err
		}

		if err := reconcileAutoscaling(agentPool, machinePool); err != nil {
			return err
		}

		if err := setAgentPoolCount(ctx, ctrlClient, machinePool, agentPoolPath, agentPool); err != nil {
			return err
		}

		return nil
	}
}

func setAgentPoolOrchestratorVersion(ctx context.Context, machinePool *expv1.MachinePool, agentPoolPath string, agentPool *unstructured.Unstructured) error {
	_, log, done := tele.StartSpanWithLogger(ctx, "mutators.setAgentPoolOrchestratorVersion")
	defer done()

	if machinePool.Spec.Template.Spec.Version == nil {
		return nil
	}

	k8sVersionPath := []string{"spec", "orchestratorVersion"}
	capiK8sVersion := strings.TrimPrefix(*machinePool.Spec.Template.Spec.Version, "v")
	userK8sVersion, k8sVersionFound, err := unstructured.NestedString(agentPool.UnstructuredContent(), k8sVersionPath...)
	if err != nil {
		return err
	}
	setK8sVersion := mutation{
		location: agentPoolPath + "." + strings.Join(k8sVersionPath, "."),
		val:      capiK8sVersion,
		reason:   fmt.Sprintf("because MachinePool %s's spec.template.spec.version is %s", machinePool.Name, *machinePool.Spec.Template.Spec.Version),
	}
	if k8sVersionFound && userK8sVersion != capiK8sVersion {
		return Incompatible{
			mutation: setK8sVersion,
			userVal:  userK8sVersion,
		}
	}
	logMutation(log, setK8sVersion)
	return unstructured.SetNestedField(agentPool.UnstructuredContent(), capiK8sVersion, k8sVersionPath...)
}

func reconcileAutoscaling(agentPool *unstructured.Unstructured, machinePool *expv1.MachinePool) error {
	autoscaling, _, err := unstructured.NestedBool(agentPool.UnstructuredContent(), "spec", "enableAutoScaling")
	if err != nil {
		return err
	}

	// Update the MachinePool replica manager annotation. This isn't wrapped in a mutation object because
	// it's not modifying an ASO resource and users are not expected to set this manually. This behavior
	// is documented by CAPI as expected of a provider.
	replicaManager, ok := machinePool.Annotations[clusterv1.ReplicasManagedByAnnotation]
	if autoscaling {
		if !ok {
			if machinePool.Annotations == nil {
				machinePool.Annotations = make(map[string]string)
			}
			machinePool.Annotations[clusterv1.ReplicasManagedByAnnotation] = infrav1alpha.ReplicasManagedByAKS
		} else if replicaManager != infrav1alpha.ReplicasManagedByAKS {
			return fmt.Errorf("failed to enable autoscaling, replicas are already being managed by %s according to MachinePool %s's %s annotation", replicaManager, machinePool.Name, clusterv1.ReplicasManagedByAnnotation)
		}
	} else if !autoscaling && replicaManager == infrav1alpha.ReplicasManagedByAKS {
		// Removing this annotation informs the MachinePool controller that this MachinePool is no longer
		// being autoscaled.
		delete(machinePool.Annotations, clusterv1.ReplicasManagedByAnnotation)
	}

	return nil
}

func setAgentPoolCount(ctx context.Context, ctrlClient client.Client, machinePool *expv1.MachinePool, agentPoolPath string, agentPool *unstructured.Unstructured) error {
	_, log, done := tele.StartSpanWithLogger(ctx, "mutators.setAgentPoolCount")
	defer done()

	if machinePool.Spec.Replicas == nil {
		return nil
	}

	// When managed by any autoscaler, CAPZ should not provide any spec.count to the ManagedClustersAgentPool
	// to prevent ASO from overwriting the autoscaler's opinion of the replica count.
	// The MachinePool's spec.replicas is used to seed an initial value as required by AKS.
	if _, autoscaling := machinePool.Annotations[clusterv1.ReplicasManagedByAnnotation]; autoscaling {
		existingAgentPool := &asocontainerservicev1.ManagedClustersAgentPool{}
		err := ctrlClient.Get(ctx, client.ObjectKey{Namespace: machinePool.GetNamespace(), Name: agentPool.GetName()}, existingAgentPool)
		if client.IgnoreNotFound(err) != nil {
			return err
		}
		if err == nil && existingAgentPool.Status.Count != nil {
			return nil
		}
	}

	countPath := []string{"spec", "count"}
	capiCount := int64(*machinePool.Spec.Replicas)
	userCount, countFound, err := unstructured.NestedInt64(agentPool.UnstructuredContent(), countPath...)
	if err != nil {
		return err
	}
	setCount := mutation{
		location: agentPoolPath + "." + strings.Join(countPath, "."),
		val:      capiCount,
		reason:   fmt.Sprintf("because MachinePool %s's spec.replicas is %d", machinePool.Name, capiCount),
	}
	if countFound && userCount != capiCount {
		return Incompatible{
			mutation: setCount,
			userVal:  userCount,
		}
	}
	logMutation(log, setCount)
	return unstructured.SetNestedField(agentPool.UnstructuredContent(), capiCount, countPath...)
}
