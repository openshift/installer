/*
Copyright 2025 The Kubernetes Authors.

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

package controllers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	v1beta1patch "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/patch"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/scope"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/resourceskus"
	"sigs.k8s.io/cluster-api-provider-azure/util/reconciler"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

// AzureMachineTemplateReconciler reconciles AzureMachineTemplateReconciler objects.
type AzureMachineTemplateReconciler struct {
	client.Client
	Recorder         record.EventRecorder
	Timeouts         reconciler.Timeouts
	WatchFilterValue string
	CredentialCache  azure.CredentialCache
}

// SetupWithManager initializes this controller with a manager.
func (r *AzureMachineTemplateReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	_, log, done := tele.StartSpanWithLogger(ctx,
		"controllers.AzureMachineTemplateReconciler.SetupWithManager",
	)
	defer done()

	azureMachineTemplateMapper, err := util.ClusterToTypedObjectsMapper(r.Client, &infrav1.AzureMachineTemplateList{}, mgr.GetScheme())
	if err != nil {
		return errors.Wrap(err, "failed to create mapper for Cluster to AzureMachineTemplates")
	}

	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&infrav1.AzureMachineTemplate{}).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(mgr.GetScheme(), log, r.WatchFilterValue)).
		// Add a watch on Clusters to requeue when the infraRef is set. This is needed because the infraRef is not initially
		// set in Clusters created from a ClusterClass.
		Watches(
			&clusterv1.Cluster{},
			handler.EnqueueRequestsFromMapFunc(azureMachineTemplateMapper),
			builder.WithPredicates(
				predicates.ClusterPausedTransitionsOrInfrastructureProvisioned(mgr.GetScheme(), log),
				predicates.ResourceNotPausedAndHasFilterLabel(mgr.GetScheme(), log, r.WatchFilterValue),
			),
		).
		Complete(r)
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=azuremachinetemplates,verbs=get;list;watch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=azuremachinetemplates/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters,verbs=get;list;watch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=azureclusters,verbs=get;list;watch

// Reconcile reconciles the AzureMachineTemplate status to populate capacity and nodeInfo for autoscaling-from-zero support.
func (r *AzureMachineTemplateReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	ctx, cancel := context.WithTimeout(ctx, r.Timeouts.DefaultedLoopTimeout())
	defer cancel()

	ctx, log, done := tele.StartSpanWithLogger(ctx, "controllers.AzureMachineTemplateReconciler.Reconcile",
		tele.KVP("namespace", req.Namespace),
		tele.KVP("name", req.Name),
		tele.KVP("kind", "AzureMachineTemplate"),
	)
	defer done()

	// Fetch the AzureMachineTemplate instance
	azureMachineTemplate := &infrav1.AzureMachineTemplate{}
	err := r.Get(ctx, req.NamespacedName, azureMachineTemplate)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// Fetch the owner Cluster
	cluster, err := util.GetOwnerCluster(ctx, r.Client, azureMachineTemplate.ObjectMeta)
	if err != nil {
		return reconcile.Result{}, err
	}
	if cluster == nil {
		log.Info("Cluster Controller has not yet set OwnerRef")
		return reconcile.Result{}, nil
	}

	log = log.WithValues("cluster", cluster.Name)

	// Return early if the object or Cluster is paused.
	if annotations.IsPaused(cluster, azureMachineTemplate) {
		log.Info("AzureMachineTemplate or linked Cluster is marked as paused. Won't reconcile")
		return ctrl.Result{}, nil
	}

	// Only look at Azure clusters
	if !cluster.Spec.InfrastructureRef.IsDefined() {
		log.Info("infra ref is not defined")
		return ctrl.Result{}, nil
	}
	if cluster.Spec.InfrastructureRef.Kind != infrav1.AzureClusterKind {
		log.WithValues("kind", cluster.Spec.InfrastructureRef.Kind).Info("infra ref was not an AzureCluster")
		return ctrl.Result{}, nil
	}

	// Fetch the corresponding AzureCluster
	azureCluster := &infrav1.AzureCluster{}
	azureClusterName := types.NamespacedName{
		Namespace: req.Namespace,
		Name:      cluster.Spec.InfrastructureRef.Name,
	}

	if err := r.Get(ctx, azureClusterName, azureCluster); err != nil {
		log.Error(err, "failed to fetch AzureCluster")
		return reconcile.Result{}, err
	}

	return r.reconcileNormal(ctx, cluster, azureCluster, azureMachineTemplate)
}

func (r *AzureMachineTemplateReconciler) reconcileNormal(ctx context.Context, cluster *clusterv1.Cluster, azureCluster *infrav1.AzureCluster, azureMachineTemplate *infrav1.AzureMachineTemplate) (ctrl.Result, error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "controllers.AzureMachineTemplateReconciler.reconcileNormal")
	defer done()

	// Create patch helper
	patchHelper, err := v1beta1patch.NewHelper(azureMachineTemplate, r.Client)
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to init patch helper")
	}

	// Create the cluster scope
	clusterScope, err := scope.NewClusterScope(ctx, scope.ClusterScopeParams{
		Client:          r.Client,
		Cluster:         cluster,
		AzureCluster:    azureCluster,
		Timeouts:        r.Timeouts,
		CredentialCache: r.CredentialCache,
	})
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to create scope")
	}

	vmSize := azureMachineTemplate.Spec.Template.Spec.VMSize

	// Get SKU cache
	skuCache, err := resourceskus.GetCache(clusterScope, clusterScope.Location())
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to get SKU cache")
	}

	// Get capacity from VM size
	capacity, err := r.getVMSizeCapacity(ctx, skuCache, vmSize)
	if err != nil {
		log.Error(err, "failed to get VM size capacity")
		r.Recorder.Eventf(azureMachineTemplate, corev1.EventTypeWarning, "SKUNotFound",
			"Failed to get capacity for VM size %s: %v", vmSize, err)
		return reconcile.Result{}, err
	}

	// Get node info from VM size
	nodeInfo, err := r.getVMSizeNodeInfo(ctx, skuCache, vmSize, azureMachineTemplate.Spec.Template.Spec.OSDisk.OSType)
	if err != nil {
		log.Error(err, "failed to get VM size node info")
		r.Recorder.Eventf(azureMachineTemplate, corev1.EventTypeWarning, "NodeInfoError",
			"Failed to get node info for VM size %s: %v", vmSize, err)
		return reconcile.Result{}, errors.Wrap(err, "failed to get node info")
	}

	// Update status
	azureMachineTemplate.Status.Capacity = capacity
	azureMachineTemplate.Status.NodeInfo = nodeInfo

	// Patch the object
	if err := patchHelper.Patch(ctx, azureMachineTemplate); err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to patch AzureMachineTemplate")
	}

	cpuQty := capacity[corev1.ResourceCPU]
	memQty := capacity[corev1.ResourceMemory]
	log.V(2).Info("Successfully updated AzureMachineTemplate status",
		"cpu", cpuQty.String(),
		"memory", memQty.String(),
		"architecture", nodeInfo.Architecture,
		"os", nodeInfo.OperatingSystem,
	)

	return ctrl.Result{}, nil
}

// getVMSizeCapacity retrieves the resource capacity for a given Azure VM size.
func (r *AzureMachineTemplateReconciler) getVMSizeCapacity(ctx context.Context, skuCache *resourceskus.Cache, vmSize string) (corev1.ResourceList, error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "controllers.AzureMachineTemplateReconciler.getVMSizeCapacity")
	defer done()

	// Query SKU for the VM size
	sku, err := skuCache.Get(ctx, vmSize, resourceskus.VirtualMachines)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get SKU for VM size %s", vmSize)
	}

	// Extract capacity from SKU
	return extractCapacityFromSKU(sku)
}

// extractCapacityFromSKU extracts CPU and memory resources from an Azure SKU.
func extractCapacityFromSKU(sku resourceskus.SKU) (corev1.ResourceList, error) {
	capacity := corev1.ResourceList{}

	// Get vCPUs
	vcpuStr, ok := sku.GetCapability(resourceskus.VCPUs)
	if !ok {
		return nil, errors.New("SKU does not have vCPUs capability")
	}
	capacity[corev1.ResourceCPU] = resource.MustParse(vcpuStr)

	// Get memory - Azure's MemoryGB capability returns GiB values (despite the name)
	memoryStr, ok := sku.GetCapability(resourceskus.MemoryGB)
	if !ok {
		return nil, errors.New("SKU does not have MemoryGB capability")
	}
	memoryGB, err := strconv.ParseFloat(memoryStr, 64)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse MemoryGB value %q", memoryStr)
	}

	// Format as Gi for whole numbers, Mi for fractions to ensure integer quantities
	intMemoryGB := int64(memoryGB)
	if memoryGB == float64(intMemoryGB) {
		capacity[corev1.ResourceMemory] = resource.MustParse(fmt.Sprintf("%dGi", intMemoryGB))
	} else {
		memoryMi := int64(memoryGB * 1024)
		capacity[corev1.ResourceMemory] = resource.MustParse(fmt.Sprintf("%dMi", memoryMi))
	}

	return capacity, nil
}

// getVMSizeNodeInfo retrieves node architecture and OS information for a given VM size.
func (r *AzureMachineTemplateReconciler) getVMSizeNodeInfo(ctx context.Context, skuCache *resourceskus.Cache, vmSize string, osType string) (*infrav1.NodeInfo, error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "controllers.AzureMachineTemplateReconciler.getVMSizeNodeInfo")
	defer done()

	// Query SKU for the VM size
	sku, err := skuCache.Get(ctx, vmSize, resourceskus.VirtualMachines)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get SKU for VM size %s", vmSize)
	}

	// Extract architecture from SKU
	azureArch, ok := sku.GetCapability(resourceskus.CPUArchitectureType)
	if !ok {
		return nil, errors.New("SKU does not have CPUArchitectureType capability")
	}
	var architecture infrav1.Architecture
	switch azureArch {
	case string(armcompute.ArchitectureX64):
		architecture = infrav1.ArchitectureAmd64
	case string(armcompute.ArchitectureArm64):
		architecture = infrav1.ArchitectureArm64
	default:
		return nil, errors.Errorf("unsupported architecture: %v", azureArch)
	}

	var operatingSystem infrav1.OperatingSystem
	switch osType {
	case azure.LinuxOS:
		operatingSystem = infrav1.OperatingSystemLinux
	case azure.WindowsOS:
		operatingSystem = infrav1.OperatingSystemWindows
	default:
		operatingSystem = infrav1.OperatingSystemLinux
	}

	return &infrav1.NodeInfo{
		Architecture:    architecture,
		OperatingSystem: operatingSystem,
	}, nil
}
