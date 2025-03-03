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
	"errors"
	"fmt"
	"strings"

	asocontainerservicev1 "github.com/Azure/azure-service-operator/v2/api/containerservice/v1api20231001"
	asocontainerservicev1hub "github.com/Azure/azure-service-operator/v2/api/containerservice/v1api20231001/storage"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	infrav1alpha "sigs.k8s.io/cluster-api-provider-azure/api/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	exputil "sigs.k8s.io/cluster-api/exp/util"
	"sigs.k8s.io/cluster-api/util/secret"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var (
	// ErrNoManagedClusterDefined describes an AzureASOManagedControlPlane without a ManagedCluster.
	ErrNoManagedClusterDefined = fmt.Errorf("no %s ManagedCluster defined in AzureASOManagedControlPlane spec.resources", asocontainerservicev1hub.GroupVersion.Group)

	// ErrNoAzureASOManagedMachinePools means no AzureASOManagedMachinePools exist for an AzureASOManagedControlPlane.
	ErrNoAzureASOManagedMachinePools = errors.New("no AzureASOManagedMachinePools found for AzureASOManagedControlPlane")
)

// SetManagedClusterDefaults propagates values defined by Cluster API to an ASO ManagedCluster.
func SetManagedClusterDefaults(ctrlClient client.Client, asoManagedControlPlane *infrav1alpha.AzureASOManagedControlPlane, cluster *clusterv1.Cluster) ResourcesMutator {
	return func(ctx context.Context, us []*unstructured.Unstructured) error {
		ctx, _, done := tele.StartSpanWithLogger(ctx, "mutators.SetManagedClusterDefaults")
		defer done()

		var managedCluster *unstructured.Unstructured
		var managedClusterPath string
		for i, u := range us {
			if u.GroupVersionKind().Group == asocontainerservicev1hub.GroupVersion.Group &&
				u.GroupVersionKind().Kind == "ManagedCluster" {
				managedCluster = u
				managedClusterPath = fmt.Sprintf("spec.resources[%d]", i)
				break
			}
		}
		if managedCluster == nil {
			return reconcile.TerminalError(ErrNoManagedClusterDefined)
		}

		if err := setManagedClusterKubernetesVersion(ctx, asoManagedControlPlane, managedClusterPath, managedCluster); err != nil {
			return err
		}

		if err := setManagedClusterServiceCIDR(ctx, cluster, managedClusterPath, managedCluster); err != nil {
			return err
		}

		if err := setManagedClusterPodCIDR(ctx, cluster, managedClusterPath, managedCluster); err != nil {
			return err
		}

		if err := setManagedClusterAgentPoolProfiles(ctx, ctrlClient, asoManagedControlPlane.Namespace, cluster, managedClusterPath, managedCluster); err != nil {
			return err
		}

		if err := setManagedClusterCredentials(ctx, cluster, managedClusterPath, managedCluster); err != nil {
			return err
		}

		return nil
	}
}

func setManagedClusterKubernetesVersion(ctx context.Context, asoManagedControlPlane *infrav1alpha.AzureASOManagedControlPlane, managedClusterPath string, managedCluster *unstructured.Unstructured) error {
	_, log, done := tele.StartSpanWithLogger(ctx, "mutators.setManagedClusterKubernetesVersion")
	defer done()

	capzK8sVersion := strings.TrimPrefix(asoManagedControlPlane.Spec.Version, "v")
	if capzK8sVersion == "" {
		// When the CAPI contract field isn't set, any value for version in the embedded ASO resource may be specified.
		return nil
	}

	k8sVersionPath := []string{"spec", "kubernetesVersion"}
	userK8sVersion, k8sVersionFound, err := unstructured.NestedString(managedCluster.UnstructuredContent(), k8sVersionPath...)
	if err != nil {
		return err
	}
	setK8sVersion := mutation{
		location: managedClusterPath + "." + strings.Join(k8sVersionPath, "."),
		val:      capzK8sVersion,
		reason:   "because spec.version is set to " + asoManagedControlPlane.Spec.Version,
	}
	if k8sVersionFound && userK8sVersion != capzK8sVersion {
		return Incompatible{
			mutation: setK8sVersion,
			userVal:  userK8sVersion,
		}
	}
	logMutation(log, setK8sVersion)
	return unstructured.SetNestedField(managedCluster.UnstructuredContent(), capzK8sVersion, k8sVersionPath...)
}

func setManagedClusterServiceCIDR(ctx context.Context, cluster *clusterv1.Cluster, managedClusterPath string, managedCluster *unstructured.Unstructured) error {
	_, log, done := tele.StartSpanWithLogger(ctx, "mutators.setManagedClusterServiceCIDR")
	defer done()

	if cluster.Spec.ClusterNetwork == nil ||
		cluster.Spec.ClusterNetwork.Services == nil ||
		len(cluster.Spec.ClusterNetwork.Services.CIDRBlocks) == 0 {
		return nil
	}

	capiCIDR := cluster.Spec.ClusterNetwork.Services.CIDRBlocks[0]

	// ManagedCluster.v1api20210501.containerservice.azure.com does not contain the plural serviceCidrs field.
	svcCIDRPath := []string{"spec", "networkProfile", "serviceCidr"}
	userSvcCIDR, found, err := unstructured.NestedString(managedCluster.UnstructuredContent(), svcCIDRPath...)
	if err != nil {
		return err
	}
	setSvcCIDR := mutation{
		location: managedClusterPath + "." + strings.Join(svcCIDRPath, "."),
		val:      capiCIDR,
		reason:   fmt.Sprintf("because spec.clusterNetwork.services.cidrBlocks[0] in Cluster %s/%s is set to %s", cluster.Namespace, cluster.Name, capiCIDR),
	}
	if found && userSvcCIDR != capiCIDR {
		return Incompatible{
			mutation: setSvcCIDR,
			userVal:  userSvcCIDR,
		}
	}
	logMutation(log, setSvcCIDR)
	return unstructured.SetNestedField(managedCluster.UnstructuredContent(), capiCIDR, svcCIDRPath...)
}

func setManagedClusterPodCIDR(ctx context.Context, cluster *clusterv1.Cluster, managedClusterPath string, managedCluster *unstructured.Unstructured) error {
	_, log, done := tele.StartSpanWithLogger(ctx, "mutators.setManagedClusterPodCIDR")
	defer done()

	if cluster.Spec.ClusterNetwork == nil ||
		cluster.Spec.ClusterNetwork.Pods == nil ||
		len(cluster.Spec.ClusterNetwork.Pods.CIDRBlocks) == 0 {
		return nil
	}

	capiCIDR := cluster.Spec.ClusterNetwork.Pods.CIDRBlocks[0]

	// ManagedCluster.v1api20210501.containerservice.azure.com does not contain the plural podCidrs field.
	podCIDRPath := []string{"spec", "networkProfile", "podCidr"}
	userPodCIDR, found, err := unstructured.NestedString(managedCluster.UnstructuredContent(), podCIDRPath...)
	if err != nil {
		return err
	}
	setPodCIDR := mutation{
		location: managedClusterPath + "." + strings.Join(podCIDRPath, "."),
		val:      capiCIDR,
		reason:   fmt.Sprintf("because spec.clusterNetwork.pods.cidrBlocks[0] in Cluster %s/%s is set to %s", cluster.Namespace, cluster.Name, capiCIDR),
	}
	if found && userPodCIDR != capiCIDR {
		return Incompatible{
			mutation: setPodCIDR,
			userVal:  userPodCIDR,
		}
	}
	logMutation(log, setPodCIDR)
	return unstructured.SetNestedField(managedCluster.UnstructuredContent(), capiCIDR, podCIDRPath...)
}

func setManagedClusterAgentPoolProfiles(ctx context.Context, ctrlClient client.Client, namespace string, cluster *clusterv1.Cluster, managedClusterPath string, managedCluster *unstructured.Unstructured) error {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "mutators.setManagedClusterAgentPoolProfiles")
	defer done()

	agentPoolProfilesPath := []string{"spec", "agentPoolProfiles"}
	userAgentPoolProfiles, agentPoolProfilesFound, err := unstructured.NestedSlice(managedCluster.UnstructuredContent(), agentPoolProfilesPath...)
	if err != nil {
		return err
	}
	setAgentPoolProfiles := mutation{
		location: managedClusterPath + "." + strings.Join(agentPoolProfilesPath, "."),
		val:      "nil",
		reason:   "because agent pool definitions must be inherited from AzureASOManagedMachinePools",
	}
	if agentPoolProfilesFound {
		return Incompatible{
			mutation: setAgentPoolProfiles,
			userVal:  fmt.Sprintf("<slice of length %d>", len(userAgentPoolProfiles)),
		}
	}

	// AKS requires ManagedClusters to be created with agent pools: https://github.com/Azure/azure-service-operator/issues/2791
	getMC := &asocontainerservicev1.ManagedCluster{}
	err = ctrlClient.Get(ctx, client.ObjectKey{Namespace: namespace, Name: managedCluster.GetName()}, getMC)
	if client.IgnoreNotFound(err) != nil {
		return err
	}
	if len(getMC.Status.AgentPoolProfiles) != 0 {
		return nil
	}

	log.V(4).Info("gathering agent pool profiles to include in ManagedCluster create")
	agentPools, err := agentPoolsFromManagedMachinePools(ctx, ctrlClient, cluster.Name, namespace)
	if err != nil {
		return err
	}
	mc, err := ctrlClient.Scheme().New(managedCluster.GroupVersionKind())
	if err != nil {
		return err
	}
	err = ctrlClient.Scheme().Convert(managedCluster, mc, nil)
	if err != nil {
		return err
	}
	setAgentPoolProfiles.val = fmt.Sprintf("<slice of length %d>", len(agentPools))
	logMutation(log, setAgentPoolProfiles)
	err = setAgentPoolProfilesFromAgentPools(mc.(conversion.Convertible), agentPools)
	if err != nil {
		return err
	}
	err = ctrlClient.Scheme().Convert(mc, managedCluster, nil)
	if err != nil {
		return err
	}

	return nil
}

func agentPoolsFromManagedMachinePools(ctx context.Context, ctrlClient client.Client, clusterName string, namespace string) ([]conversion.Convertible, error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "mutators.agentPoolsFromManagedMachinePools")
	defer done()

	asoManagedMachinePools := &infrav1alpha.AzureASOManagedMachinePoolList{}
	err := ctrlClient.List(ctx, asoManagedMachinePools,
		client.InNamespace(namespace),
		client.MatchingLabels{
			clusterv1.ClusterNameLabel: clusterName,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list AzureASOManagedMachinePools: %w", err)
	}

	var agentPools []conversion.Convertible
	for _, asoManagedMachinePool := range asoManagedMachinePools.Items {
		machinePool, err := exputil.GetOwnerMachinePool(ctx, ctrlClient, asoManagedMachinePool.ObjectMeta)
		if err != nil {
			return nil, err
		}
		if machinePool == nil {
			log.V(2).Info("Waiting for MachinePool Controller to set OwnerRef on AzureASOManagedMachinePool")
			return nil, nil
		}

		resources, err := ApplyMutators(ctx, asoManagedMachinePool.Spec.Resources,
			SetAgentPoolDefaults(ctrlClient, machinePool),
		)
		if err != nil {
			return nil, err
		}

		for _, u := range resources {
			if u.GroupVersionKind().Group != asocontainerservicev1hub.GroupVersion.Group ||
				u.GroupVersionKind().Kind != "ManagedClustersAgentPool" {
				continue
			}

			agentPool, err := ctrlClient.Scheme().New(u.GroupVersionKind())
			if err != nil {
				return nil, fmt.Errorf("error creating new %v: %w", u.GroupVersionKind(), err)
			}
			err = ctrlClient.Scheme().Convert(u, agentPool, nil)
			if err != nil {
				return nil, err
			}

			agentPools = append(agentPools, agentPool.(conversion.Convertible))
			break
		}
	}

	return agentPools, nil
}

func setAgentPoolProfilesFromAgentPools(managedCluster conversion.Convertible, agentPools []conversion.Convertible) error {
	hubMC := &asocontainerservicev1hub.ManagedCluster{}
	err := managedCluster.ConvertTo(hubMC)
	if err != nil {
		return err
	}
	hubMC.Spec.AgentPoolProfiles = nil

	for _, agentPool := range agentPools {
		hubPool := &asocontainerservicev1hub.ManagedClustersAgentPool{}
		err := agentPool.ConvertTo(hubPool)
		if err != nil {
			return err
		}

		profile := asocontainerservicev1hub.ManagedClusterAgentPoolProfile{
			AvailabilityZones:                 hubPool.Spec.AvailabilityZones,
			CapacityReservationGroupReference: hubPool.Spec.CapacityReservationGroupReference,
			Count:                             hubPool.Spec.Count,
			CreationData:                      hubPool.Spec.CreationData,
			EnableAutoScaling:                 hubPool.Spec.EnableAutoScaling,
			EnableEncryptionAtHost:            hubPool.Spec.EnableEncryptionAtHost,
			EnableFIPS:                        hubPool.Spec.EnableFIPS,
			EnableNodePublicIP:                hubPool.Spec.EnableNodePublicIP,
			EnableUltraSSD:                    hubPool.Spec.EnableUltraSSD,
			GpuInstanceProfile:                hubPool.Spec.GpuInstanceProfile,
			HostGroupReference:                hubPool.Spec.HostGroupReference,
			KubeletConfig:                     hubPool.Spec.KubeletConfig,
			KubeletDiskType:                   hubPool.Spec.KubeletDiskType,
			LinuxOSConfig:                     hubPool.Spec.LinuxOSConfig,
			MaxCount:                          hubPool.Spec.MaxCount,
			MaxPods:                           hubPool.Spec.MaxPods,
			MinCount:                          hubPool.Spec.MinCount,
			Mode:                              hubPool.Spec.Mode,
			Name:                              azure.AliasOrNil[string](&hubPool.Spec.AzureName),
			NetworkProfile:                    hubPool.Spec.NetworkProfile,
			NodeLabels:                        hubPool.Spec.NodeLabels,
			NodePublicIPPrefixReference:       hubPool.Spec.NodePublicIPPrefixReference,
			NodeTaints:                        hubPool.Spec.NodeTaints,
			OrchestratorVersion:               hubPool.Spec.OrchestratorVersion,
			OsDiskSizeGB:                      hubPool.Spec.OsDiskSizeGB,
			OsDiskType:                        hubPool.Spec.OsDiskType,
			OsSKU:                             hubPool.Spec.OsSKU,
			OsType:                            hubPool.Spec.OsType,
			PodSubnetReference:                hubPool.Spec.PodSubnetReference,
			PowerState:                        hubPool.Spec.PowerState,
			PropertyBag:                       hubPool.Spec.PropertyBag,
			ProximityPlacementGroupReference:  hubPool.Spec.ProximityPlacementGroupReference,
			ScaleDownMode:                     hubPool.Spec.ScaleDownMode,
			ScaleSetEvictionPolicy:            hubPool.Spec.ScaleSetEvictionPolicy,
			ScaleSetPriority:                  hubPool.Spec.ScaleSetPriority,
			SpotMaxPrice:                      hubPool.Spec.SpotMaxPrice,
			Tags:                              hubPool.Spec.Tags,
			Type:                              hubPool.Spec.Type,
			UpgradeSettings:                   hubPool.Spec.UpgradeSettings,
			VmSize:                            hubPool.Spec.VmSize,
			VnetSubnetReference:               hubPool.Spec.VnetSubnetReference,
			WorkloadRuntime:                   hubPool.Spec.WorkloadRuntime,
		}

		hubMC.Spec.AgentPoolProfiles = append(hubMC.Spec.AgentPoolProfiles, profile)
	}

	return managedCluster.ConvertFrom(hubMC)
}

func setManagedClusterCredentials(ctx context.Context, cluster *clusterv1.Cluster, managedClusterPath string, managedCluster *unstructured.Unstructured) error {
	_, log, done := tele.StartSpanWithLogger(ctx, "mutators.setManagedClusterCredentials")
	defer done()

	// CAPZ only cares that some set of credentials is created by ASO, but not where. CAPZ will propagate
	// whatever is defined in the ASO resource to the <cluster>-kubeconfig secret as expected by CAPI.

	_, hasUserCreds, err := unstructured.NestedMap(managedCluster.UnstructuredContent(), "spec", "operatorSpec", "secrets", "userCredentials")
	if err != nil {
		return err
	}
	if hasUserCreds {
		return nil
	}

	_, hasAdminCreds, err := unstructured.NestedMap(managedCluster.UnstructuredContent(), "spec", "operatorSpec", "secrets", "adminCredentials")
	if err != nil {
		return err
	}
	if hasAdminCreds {
		return nil
	}

	secrets := map[string]interface{}{
		"adminCredentials": map[string]interface{}{
			"name": cluster.Name + "-" + string(secret.Kubeconfig),
			"key":  secret.KubeconfigDataName,
		},
	}

	setCreds := mutation{
		location: managedClusterPath + ".spec.operatorSpec.secrets",
		val:      secrets,
		reason:   "because no userCredentials or adminCredentials are defined",
	}
	logMutation(log, setCreds)
	return unstructured.SetNestedMap(managedCluster.UnstructuredContent(), secrets, "spec", "operatorSpec", "secrets")
}
