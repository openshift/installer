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
	"encoding/json"
	"errors"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/tools/events"
	"k8s.io/klog/v2"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	conditions "sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	infrav1alpha1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha1"
	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/cloud/services/compute"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/cloud/services/networking"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/scope"
	controllers "sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/controllers"
)

const (
	imagePropertyForOS = "os_type"

	// annotationAllowedAddressPairs tracks the last-applied allowedAddressPairs per
	// OpenStackMachine (stored as JSON). Written as a metadata annotation so we never
	// touch the immutable OSM spec, which would trigger the spec-immutability webhook.
	annotationAllowedAddressPairs = "infrastructure.cluster.x-k8s.io/osmt-allowed-address-pairs"
)

// Set here so we can easily mock it in tests.
var (
	newComputeService    = compute.NewService
	newNetworkingService = networking.NewService
)

// OpenStackMachineTemplateReconciler reconciles a OpenStackMachineTemplate object.
// it only updates the .status field to allow auto-scaling.
type OpenStackMachineTemplateReconciler struct {
	Client           client.Client
	Recorder         events.EventRecorder
	WatchFilterValue string
	ScopeFactory     scope.Factory
	CaCertificates   []byte // PEM encoded ca certificates.
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=openstackmachinetemplates,verbs=get;list;watch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=openstackmachinetemplates/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=openstackmachines,verbs=get;list;watch;patch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=openstackservers,verbs=get;list;watch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machinesets,verbs=get;list;watch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machines,verbs=get;list;watch

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

	log = log.WithValues("OpenStackMachineTemplate", klog.KObj(openStackMachineTemplate))
	log.V(4).Info("Reconciling OpenStackMachineTemplate")

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
		log.Info("OpenStackMachineTemplate is missing owner Cluster or Cluster does not exist")
		return ctrl.Result{}, nil
	}

	log = log.WithValues("Cluster", klog.KObj(cluster))

	if annotations.IsPaused(cluster, openStackMachineTemplate) {
		log.Info("OpenStackMachineTemplate or linked Cluster is marked as paused. Won't reconcile")
		return ctrl.Result{}, nil
	}

	infraCluster, err := controllers.GetInfraCluster(ctx, r.Client, cluster)
	if err != nil {
		return ctrl.Result{}, errors.New("error getting infra provider cluster")
	}
	if infraCluster == nil {
		log.Info("OpenStackCluster is not ready", "OpenStackCluster", klog.KRef(cluster.Namespace, cluster.Spec.InfrastructureRef.Name))
		return ctrl.Result{}, nil
	}

	log = log.WithValues("OpenStackCluster", klog.KObj(infraCluster))

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

	clientScope, err := r.ScopeFactory.NewClientScopeFromObject(ctx, r.Client, r.CaCertificates, log, openStackMachineTemplate, infraCluster)
	if err != nil {
		conditions.Set(openStackMachineTemplate, metav1.Condition{
			Type:    infrav1.OpenStackAuthenticationSucceeded,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.OpenStackAuthenticationFailedReason,
			Message: fmt.Sprintf("Failed to create OpenStack client scope: %v", err),
		})
		return ctrl.Result{}, err
	}
	conditions.Set(openStackMachineTemplate, metav1.Condition{
		Type:   infrav1.OpenStackAuthenticationSucceeded,
		Status: metav1.ConditionTrue,
		Reason: infrav1.ReadyConditionReason,
	})
	scope := scope.NewWithLogger(clientScope, log)

	// Handle non-deleted OpenStackMachineTemplates
	if err := r.reconcileNormal(ctx, scope, cluster.Name, openStackMachineTemplate); err != nil {
		return ctrl.Result{}, err
	}
	log.V(4).Info("Successfully reconciled OpenStackMachineTemplate")
	return ctrl.Result{}, nil
}

func (r *OpenStackMachineTemplateReconciler) reconcileNormal(ctx context.Context, scope *scope.WithLogger, clusterName string, openStackMachineTemplate *infrav1.OpenStackMachineTemplate) (reterr error) {
	log := scope.Logger()

	computeService, err := newComputeService(scope)
	if err != nil {
		return err
	}

	flavorID, err := computeService.GetFlavorID(openStackMachineTemplate.Spec.Template.Spec.Flavor)
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

	// reconcileAllowedAddressPairs is called independently of the image/flavor logic so that
	// it is never skipped by an early return (e.g. when imageID is not yet resolvable).
	if err := r.reconcileAllowedAddressPairs(ctx, scope, clusterName, openStackMachineTemplate); err != nil {
		return err
	}

	imageID, err := computeService.GetImageID(ctx, r.Client, openStackMachineTemplate.Namespace, openStackMachineTemplate.Spec.Template.Spec.Image)
	if err != nil {
		return err
	}
	if imageID == nil {
		return nil
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

// reconcileAllowedAddressPairs updates the allowedAddressPairs on existing Neutron ports
// to match what is defined in the OpenStackMachineTemplate.
// Idempotency is tracked via an annotation on each OpenStackMachine so that only a
// metadata-only patch is needed — this avoids touching the immutable OSM spec.
func (r *OpenStackMachineTemplateReconciler) reconcileAllowedAddressPairs(ctx context.Context, scope *scope.WithLogger, clusterName string, openStackMachineTemplate *infrav1.OpenStackMachineTemplate) error {
	log := scope.Logger()

	if len(openStackMachineTemplate.Spec.Template.Spec.Ports) == 0 || clusterName == "" {
		return nil
	}

	// Build the desired state as JSON for idempotency comparison.
	type portPairs = []infrav1.AddressPair
	templatePorts := openStackMachineTemplate.Spec.Template.Spec.Ports
	desired := make([]portPairs, len(templatePorts))
	for i, p := range templatePorts {
		desired[i] = p.AllowedAddressPairs
	}
	desiredJSON, err := json.Marshal(desired)
	if err != nil {
		return err
	}
	desiredStr := string(desiredJSON)

	// List MachineSets in the namespace for this cluster.
	machineSetList := &clusterv1.MachineSetList{}
	if err := r.Client.List(ctx, machineSetList,
		client.InNamespace(openStackMachineTemplate.Namespace),
		client.MatchingLabels{clusterv1.ClusterNameLabel: clusterName},
	); err != nil {
		return err
	}

	// List Machines in the namespace for this cluster.
	machineList := &clusterv1.MachineList{}
	if err := r.Client.List(ctx, machineList,
		client.InNamespace(openStackMachineTemplate.Namespace),
		client.MatchingLabels{clusterv1.ClusterNameLabel: clusterName},
	); err != nil {
		return err
	}

	// Networking service is initialised lazily on first actual port update.
	var networkingService *networking.Service

	for i := range machineSetList.Items {
		ms := &machineSetList.Items[i]
		if ms.Spec.Template.Spec.InfrastructureRef.Name != openStackMachineTemplate.Name {
			continue
		}

		for j := range machineList.Items {
			machine := &machineList.Items[j]
			if !isOwnedByMachineSet(machine, ms) {
				continue
			}
			infraName := machine.Spec.InfrastructureRef.Name
			if infraName == "" {
				continue
			}

			osm := &infrav1.OpenStackMachine{}
			if err := r.Client.Get(ctx, client.ObjectKey{
				Namespace: openStackMachineTemplate.Namespace,
				Name:      infraName,
			}, osm); err != nil {
				if apierrors.IsNotFound(err) {
					continue
				}
				return err
			}

			// Skip if annotation already reflects the desired state.
			if osm.Annotations[annotationAllowedAddressPairs] == desiredStr {
				continue
			}

			// Port IDs are stored in the OpenStackServer status (same name as the OSM).
			openStackServer := &infrav1alpha1.OpenStackServer{}
			if err := r.Client.Get(ctx, client.ObjectKey{
				Namespace: openStackMachineTemplate.Namespace,
				Name:      infraName,
			}, openStackServer); err != nil {
				if apierrors.IsNotFound(err) {
					continue
				}
				return err
			}

			if openStackServer.Status.Resources == nil || len(openStackServer.Status.Resources.Ports) == 0 {
				continue
			}

			if networkingService == nil {
				networkingService, err = newNetworkingService(scope)
				if err != nil {
					return err
				}
			}

			for portIdx, portStatus := range openStackServer.Status.Resources.Ports {
				if portIdx >= len(templatePorts) {
					break
				}
				pairs := templatePorts[portIdx].AllowedAddressPairs
				log.Info("Updating spec.template.spec.ports[*].allowedAddressPairs on Neutron port",
					"OpenStackMachine", klog.KObj(osm), "portID", portStatus.ID, "portIndex", portIdx)
				if err := networkingService.UpdateAllowedAddressPairs(portStatus.ID, pairs); err != nil {
					log.Error(err, "Failed to update spec.template.spec.ports[*].allowedAddressPairs",
						"OpenStackMachine", klog.KObj(osm), "portID", portStatus.ID)
					return err
				}
			}

			// Record the applied state in an annotation (metadata-only patch).
			osmCopy := osm.DeepCopy()
			if osm.Annotations == nil {
				osm.Annotations = map[string]string{}
			}
			osm.Annotations[annotationAllowedAddressPairs] = desiredStr
			if err := r.Client.Patch(ctx, osm, client.MergeFrom(osmCopy)); err != nil {
				return err
			}
		}
	}

	return nil
}

// isOwnedByMachineSet returns true if the Machine has an owner reference pointing to ms.
func isOwnedByMachineSet(machine *clusterv1.Machine, ms *clusterv1.MachineSet) bool {
	for _, ref := range machine.OwnerReferences {
		if ref.Kind == "MachineSet" && ref.Name == ms.Name {
			return true
		}
	}
	return false
}

func (r *OpenStackMachineTemplateReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := ctrl.LoggerFrom(ctx)

	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&infrav1.OpenStackMachineTemplate{}).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(mgr.GetScheme(), log, r.WatchFilterValue)).
		Complete(r)
}
