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

package controllers

import (
	"context"
	"fmt"

	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	rosacontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/rosa/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/exp/utils"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	stsservice "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/sts"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/rosa"
	"sigs.k8s.io/cluster-api-provider-aws/v2/util/paused"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	expclusterv1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/predicates"
)

// ROSAClusterReconciler reconciles ROSACluster.
type ROSAClusterReconciler struct {
	client.Client
	Recorder         record.EventRecorder
	WatchFilterValue string
	NewStsClient     func(cloud.ScopeUsage, cloud.Session, logger.Wrapper, runtime.Object) stsservice.STSClient
	NewOCMClient     func(ctx context.Context, rosaScope *scope.ROSAControlPlaneScope) (rosa.OCMClient, error)
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=rosaclusters,verbs=get;list;watch;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=rosaclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=controlplane.cluster.x-k8s.io,resources=rosacontrolplanes;rosacontrolplanes/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machinepools;machinepools/status,verbs=get;list;watch;create
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=rosamachinepools;rosamachinepools/status,verbs=get;list;watch;create
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete

func (r *ROSAClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	log := logger.FromContext(ctx)
	log.Info("Reconciling ROSACluster")

	// Fetch the ROSACluster instance
	rosaCluster := &expinfrav1.ROSACluster{}
	err := r.Get(ctx, req.NamespacedName, rosaCluster)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	if !rosaCluster.DeletionTimestamp.IsZero() {
		log.Info("Deleting ROSACluster.")
		return reconcile.Result{}, nil
	}

	// Fetch the Cluster.
	cluster, err := util.GetOwnerCluster(ctx, r.Client, rosaCluster.ObjectMeta)
	if err != nil {
		return reconcile.Result{}, err
	}

	if cluster == nil {
		log.Info("Cluster Controller has not yet set OwnerRef")
		return reconcile.Result{}, nil
	}

	if isPaused, conditionChanged, err := paused.EnsurePausedCondition(ctx, r.Client, cluster, rosaCluster); err != nil || isPaused || conditionChanged {
		return ctrl.Result{}, err
	}

	log = log.WithValues("cluster", cluster.Name)

	controlPlane := &rosacontrolplanev1.ROSAControlPlane{}
	controlPlaneRef := types.NamespacedName{
		Name:      cluster.Spec.ControlPlaneRef.Name,
		Namespace: cluster.Spec.ControlPlaneRef.Namespace,
	}

	if err := r.Get(ctx, controlPlaneRef, controlPlane); err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to get control plane ref: %w", err)
	}

	log = log.WithValues("controlPlane", controlPlaneRef.Name)

	patchHelper, err := patch.NewHelper(rosaCluster, r.Client)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to init patch helper: %w", err)
	}

	// Set the values from the managed control plane
	rosaCluster.Status.Ready = true
	rosaCluster.Spec.ControlPlaneEndpoint = controlPlane.Spec.ControlPlaneEndpoint

	if err := patchHelper.Patch(ctx, rosaCluster); err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to patch ROSACluster: %w", err)
	}

	rosaScope, err := scope.NewROSAControlPlaneScope(scope.ROSAControlPlaneScopeParams{
		Client:         r.Client,
		Cluster:        cluster,
		ControlPlane:   controlPlane,
		ControllerName: "",
		Logger:         log,
		NewStsClient:   r.NewStsClient,
	})
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create rosa controlplane scope: %w", err)
	}

	err = r.syncROSAClusterNodePools(ctx, controlPlane, rosaScope)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to sync ROSA cluster nodePools: %w", err)
	}

	log.Info("Successfully reconciled ROSACluster")

	return reconcile.Result{}, nil
}

func (r *ROSAClusterReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := logger.FromContext(ctx)
	r.NewOCMClient = rosa.NewWrappedOCMClient
	r.NewStsClient = scope.NewSTSClient

	rosaCluster := &expinfrav1.ROSACluster{}

	controller, err := ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(rosaCluster).
		WithEventFilter(predicates.ResourceHasFilterLabel(mgr.GetScheme(), ctrl.LoggerFrom(ctx), r.WatchFilterValue)).
		Build(r)
	if err != nil {
		return fmt.Errorf("error creating controller: %w", err)
	}

	// Add a watch for clusterv1.Cluster unpaise
	if err = controller.Watch(
		source.Kind[client.Object](mgr.GetCache(), &clusterv1.Cluster{},
			handler.EnqueueRequestsFromMapFunc(util.ClusterToInfrastructureMapFunc(ctx, infrav1.GroupVersion.WithKind("ROSACluster"), mgr.GetClient(), &expinfrav1.ROSACluster{})),
			predicates.ClusterPausedTransitions(mgr.GetScheme(), log.GetLogger())),
	); err != nil {
		return fmt.Errorf("failed adding a watch for ready clusters: %w", err)
	}

	// Add a watch for ROSAControlPlane
	if err = controller.Watch(
		source.Kind[client.Object](mgr.GetCache(), &rosacontrolplanev1.ROSAControlPlane{},
			handler.EnqueueRequestsFromMapFunc(r.rosaControlPlaneToManagedCluster(log))),
	); err != nil {
		return fmt.Errorf("failed adding watch on ROSAControlPlane: %w", err)
	}

	return nil
}

func (r *ROSAClusterReconciler) rosaControlPlaneToManagedCluster(log *logger.Logger) handler.MapFunc {
	return func(ctx context.Context, o client.Object) []ctrl.Request {
		rosaControlPlane, ok := o.(*rosacontrolplanev1.ROSAControlPlane)
		if !ok {
			log.Error(errors.Errorf("expected a ROSAControlPlane, got %T instead", o), "failed to map ROSAControlPlane")
			return nil
		}

		log := log.WithValues("objectMapper", "rosacpTorosac", "ROSAcontrolplane", klog.KRef(rosaControlPlane.Namespace, rosaControlPlane.Name))

		if !rosaControlPlane.ObjectMeta.DeletionTimestamp.IsZero() {
			log.Info("ROSAControlPlane has a deletion timestamp, skipping mapping")
			return nil
		}

		if rosaControlPlane.Spec.ControlPlaneEndpoint.IsZero() {
			log.Debug("ROSAControlPlane has no control plane endpoint, skipping mapping")
			return nil
		}

		cluster, err := util.GetOwnerCluster(ctx, r.Client, rosaControlPlane.ObjectMeta)
		if err != nil {
			log.Error(err, "failed to get owning cluster")
			return nil
		}
		if cluster == nil {
			log.Info("no owning cluster, skipping mapping")
			return nil
		}

		rosaClusterRef := cluster.Spec.InfrastructureRef
		if rosaClusterRef == nil || rosaClusterRef.Kind != "ROSACluster" {
			log.Info("InfrastructureRef is nil or not ROSACluster, skipping mapping")
			return nil
		}

		return []ctrl.Request{
			{
				NamespacedName: types.NamespacedName{
					Name:      rosaClusterRef.Name,
					Namespace: rosaClusterRef.Namespace,
				},
			},
		}
	}
}

// getRosMachinePools returns a map of RosaMachinePool names associatd with the cluster.
func (r *ROSAClusterReconciler) getRosaMachinePoolNames(ctx context.Context, cluster *clusterv1.Cluster) (map[string]bool, error) {
	selectors := []client.ListOption{
		client.InNamespace(cluster.GetNamespace()),
		client.MatchingLabels{
			clusterv1.ClusterNameLabel: cluster.GetName(),
		},
	}

	rosaMachinePoolList := &expinfrav1.ROSAMachinePoolList{}
	err := r.Client.List(ctx, rosaMachinePoolList, selectors...)
	if err != nil {
		return nil, err
	}

	rosaMPNames := make(map[string]bool)
	for _, rosaMP := range rosaMachinePoolList.Items {
		rosaMPNames[rosaMP.Spec.NodePoolName] = true
	}

	return rosaMPNames, nil
}

// buildROSAMachinePool returns a ROSAMachinePool and its corresponding MachinePool.
func (r *ROSAClusterReconciler) buildROSAMachinePool(nodePoolName string, clusterName string, namespace string, nodePool *cmv1.NodePool) (*expinfrav1.ROSAMachinePool, *expclusterv1.MachinePool) {
	rosaMPSpec := utils.NodePoolToRosaMachinePoolSpec(nodePool)
	rosaMachinePool := &expinfrav1.ROSAMachinePool{
		TypeMeta: metav1.TypeMeta{
			APIVersion: expinfrav1.GroupVersion.String(),
			Kind:       "ROSAMachinePool",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      nodePoolName,
			Namespace: namespace,
			Labels: map[string]string{
				clusterv1.ClusterNameLabel: clusterName,
			},
		},
		Spec: rosaMPSpec,
	}
	machinePool := &expclusterv1.MachinePool{
		TypeMeta: metav1.TypeMeta{
			APIVersion: expclusterv1.GroupVersion.String(),
			Kind:       "MachinePool",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      nodePoolName,
			Namespace: namespace,
			Labels: map[string]string{
				clusterv1.ClusterNameLabel: clusterName,
			},
		},
		Spec: expclusterv1.MachinePoolSpec{
			ClusterName: clusterName,
			Replicas:    ptr.To(int32(1)),
			Template: clusterv1.MachineTemplateSpec{
				Spec: clusterv1.MachineSpec{
					ClusterName: clusterName,
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: ptr.To(string("")),
					},
					InfrastructureRef: corev1.ObjectReference{
						APIVersion: expinfrav1.GroupVersion.String(),
						Kind:       "ROSAMachinePool",
						Name:       rosaMachinePool.Name,
					},
				},
			},
		},
	}

	return rosaMachinePool, machinePool
}

// syncROSAClusterNodePools ensure every NodePool has a MachinePool and create a corresponding MachinePool if it does not exist.
func (r *ROSAClusterReconciler) syncROSAClusterNodePools(ctx context.Context, controlPlane *rosacontrolplanev1.ROSAControlPlane, rosaScope *scope.ROSAControlPlaneScope) error {
	if controlPlane.Status.Ready {
		if r.NewOCMClient == nil {
			return fmt.Errorf("failed to create OCM client: NewOCMClient is nil")
		}

		ocmClient, err := r.NewOCMClient(ctx, rosaScope)
		if err != nil || ocmClient == nil {
			return fmt.Errorf("failed to create OCM client: %w", err)
		}

		// List the ROSA-HCP nodePools and MachinePools
		nodePools, err := ocmClient.GetNodePools(rosaScope.ControlPlane.Status.ID)
		if err != nil {
			return fmt.Errorf("failed to get nodePools: %w", err)
		}

		rosaMPNames, err := r.getRosaMachinePoolNames(ctx, rosaScope.Cluster)
		if err != nil {
			return fmt.Errorf("failed to get Rosa machinePool names: %w", err)
		}

		// Ensure every NodePool has a MachinePool and create a corresponding MachinePool if it does not exist.
		var errs []error
		for _, nodePool := range nodePools {
			// continue if nodePool is not in ready state.
			if !rosa.IsNodePoolReady(nodePool) {
				continue
			}
			// continue if nodePool exist
			if rosaMPNames[nodePool.ID()] {
				continue
			}
			// create ROSAMachinePool & MachinePool
			rosaMachinePool, machinePool := r.buildROSAMachinePool(nodePool.ID(), rosaScope.Cluster.Name, rosaScope.Cluster.Namespace, nodePool)

			rosaScope.Logger.Info(fmt.Sprintf("Create ROSAMachinePool %s", rosaMachinePool.Name))
			if err = r.Client.Create(ctx, rosaMachinePool); err != nil {
				errs = append(errs, err)
			}

			rosaScope.Logger.Info(fmt.Sprintf("Create MachinePool %s", machinePool.Name))
			if err = r.Client.Create(ctx, machinePool); err != nil {
				errs = append(errs, err)
			}
		}

		if len(errs) > 0 {
			return kerrors.NewAggregate(errs)
		}
	}
	return nil
}
