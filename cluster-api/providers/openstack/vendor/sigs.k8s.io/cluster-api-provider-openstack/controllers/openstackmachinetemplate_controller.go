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
	"errors"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/cloud/services/compute"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/scope"
	controllers "sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/controllers"
)

const imagePropertyForOS = "os_type"

// Set here so we can easily mock it in tests.
var newComputeService = compute.NewService

// OpenStackMachineTemplateReconciler reconciles a OpenStackMachineTemplate object.
// it only updates the .status field to allow auto-scaling.
type OpenStackMachineTemplateReconciler struct {
	Client           client.Client
	Recorder         record.EventRecorder
	WatchFilterValue string
	ScopeFactory     scope.Factory
	CaCertificates   []byte // PEM encoded ca certificates.
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=openstackmachinetemplates,verbs=get;list;watch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=openstackmachinetemplates/status,verbs=get;update;patch

func (r *OpenStackMachineTemplateReconciler) Reconcile(ctx context.Context, req ctrl.Request) (result ctrl.Result, reterr error) {
	log := ctrl.LoggerFrom(ctx)

	// Fetch the OpenStackMachine instance.
	openStackMachineTemplate := &infrav1.OpenStackMachineTemplate{}
	err := r.Client.Get(ctx, req.NamespacedName, openStackMachineTemplate)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	log = log.WithValues("openStackMachineTemplate", openStackMachineTemplate.Name)
	log.V(4).Info("Reconciling openStackMachineTemplate")

	// If OSMT is set for deletion, do nothing
	if !openStackMachineTemplate.DeletionTimestamp.IsZero() {
		log.Info("OpenStackMachineTemplate marked for deletion, skipping reconciliation")
		return ctrl.Result{}, nil
	}

	// Fetch the Cluster.
	// OSMT can be a valid OSMT owned a running cluster OR a OSMT owned by a ClusterClass.
	// We skip reconciliation on the latter as in this case OSMT.spec might have values that
	// are patched by the CC and thus not valid.
	cluster, err := util.GetOwnerCluster(ctx, r.Client, openStackMachineTemplate.ObjectMeta)
	if err != nil || cluster == nil {
		log.Info("openStackMachineTemplate is missing owner cluster or cluster does not exist")
		return ctrl.Result{}, nil
	}

	log = log.WithValues("cluster", cluster.Name)

	if annotations.IsPaused(cluster, openStackMachineTemplate) {
		log.Info("OpenStackMachineTemplate or linked Cluster is marked as paused. Won't reconcile")
		return ctrl.Result{}, nil
	}

	infraCluster, err := controllers.GetInfraCluster(ctx, r.Client, cluster)
	if err != nil {
		return ctrl.Result{}, errors.New("error getting infra provider cluster")
	}
	if infraCluster == nil {
		log.Info("OpenStackCluster not ready", "name", cluster.Spec.InfrastructureRef.Name)
		return ctrl.Result{}, nil
	}

	log = log.WithValues("openStackCluster", infraCluster.Name)

	clientScope, err := r.ScopeFactory.NewClientScopeFromObject(ctx, r.Client, r.CaCertificates, log, openStackMachineTemplate, infraCluster)
	if err != nil {
		return ctrl.Result{}, err
	}
	scope := scope.NewWithLogger(clientScope, log)

	// Initialize the patch helper
	patchHelper, err := patch.NewHelper(openStackMachineTemplate, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Always patch the openStackMachine when exiting this function so we can persist any OpenStackMachineTemplate changes.
	defer func() {
		if err := patchHelper.Patch(ctx, openStackMachineTemplate); err != nil {
			log.Error(err, "Failed to patch OpenStackMachineTemplate after reconciliation")
			result = ctrl.Result{}
			reterr = kerrors.NewAggregate([]error{reterr, err})
		}
	}()

	// Handle non-deleted OpenStackMachineTemplates
	if err := r.reconcileNormal(ctx, scope, openStackMachineTemplate); err != nil {
		return ctrl.Result{}, err
	}
	log.V(4).Info("Successfully reconciled OpenStackMachineTemplate")
	return ctrl.Result{}, nil
}

func (r *OpenStackMachineTemplateReconciler) reconcileNormal(ctx context.Context, scope *scope.WithLogger, openStackMachineTemplate *infrav1.OpenStackMachineTemplate) (reterr error) {
	log := scope.Logger()

	computeService, err := newComputeService(scope)
	if err != nil {
		return err
	}

	flavorID, err := computeService.GetFlavorID(openStackMachineTemplate.Spec.Template.Spec.FlavorID, openStackMachineTemplate.Spec.Template.Spec.Flavor)
	if err != nil {
		return err
	}

	flavor, err := computeService.GetFlavor(flavorID)
	if err != nil {
		return err
	}

	log.V(4).Info("Retrieved flavor details", "flavorID", flavorID)

	if openStackMachineTemplate.Status.Capacity == nil {
		log.V(4).Info("Initializing status capacity map")
		openStackMachineTemplate.Status.Capacity = corev1.ResourceList{}
	}

	if flavor.VCPUs > 0 {
		openStackMachineTemplate.Status.Capacity[corev1.ResourceCPU] = *resource.NewQuantity(int64(flavor.VCPUs), resource.DecimalSI)
	}

	if flavor.RAM > 0 {
		// flavor.RAM is in MiB -> convert to bytes
		ramBytes := int64(flavor.RAM) * 1024 * 1024
		openStackMachineTemplate.Status.Capacity[corev1.ResourceMemory] = *resource.NewQuantity(ramBytes, resource.BinarySI)
	}

	if flavor.Ephemeral > 0 {
		// flavor.Ephemeral is in GiB -> convert to bytes
		ephemeralBytes := int64(flavor.Ephemeral) * 1024 * 1024 * 1024
		openStackMachineTemplate.Status.Capacity[corev1.ResourceEphemeralStorage] = *resource.NewQuantity(ephemeralBytes, resource.BinarySI)
	}

	// storage depends on whether user boots-from-volume or not
	if openStackMachineTemplate.Spec.Template.Spec.RootVolume != nil && openStackMachineTemplate.Spec.Template.Spec.RootVolume.SizeGiB > 0 {
		// RootVolume.SizeGib is in GiB -> convert to bytes
		storageBytes := int64(openStackMachineTemplate.Spec.Template.Spec.RootVolume.SizeGiB) * 1024 * 1024 * 1024
		openStackMachineTemplate.Status.Capacity[corev1.ResourceStorage] = *resource.NewQuantity(storageBytes, resource.BinarySI)
	} else if flavor.Disk > 0 {
		// flavor.Disk is in GiB -> convert to bytes
		storageBytes := int64(flavor.Disk) * 1024 * 1024 * 1024
		openStackMachineTemplate.Status.Capacity[corev1.ResourceStorage] = *resource.NewQuantity(storageBytes, resource.BinarySI)
	}

	imageID, err := computeService.GetImageID(ctx, r.Client, openStackMachineTemplate.Namespace, openStackMachineTemplate.Spec.Template.Spec.Image)
	if err != nil {
		return err
	}

	image, err := computeService.GetImageDetails(*imageID)
	if err != nil {
		return err
	}

	log.V(4).Info("Retrieved image details", "imageID", imageID)

	if image.Properties != nil {
		if v, ok := image.Properties[imagePropertyForOS]; ok {
			if osType, ok := v.(string); ok {
				openStackMachineTemplate.Status.NodeInfo.OperatingSystem = osType
			}
		}
	}

	return nil
}

func (r *OpenStackMachineTemplateReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := ctrl.LoggerFrom(ctx)

	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&infrav1.OpenStackMachineTemplate{}).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(mgr.GetScheme(), log, r.WatchFilterValue)).
		Complete(r)
}
