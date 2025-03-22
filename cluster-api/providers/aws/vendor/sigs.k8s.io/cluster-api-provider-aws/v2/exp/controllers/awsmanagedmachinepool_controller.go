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

package controllers

import (
	"context"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/ec2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/eks"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	expclusterv1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/predicates"
)

// AWSManagedMachinePoolReconciler reconciles a AWSManagedMachinePool object.
type AWSManagedMachinePoolReconciler struct {
	client.Client
	Recorder                     record.EventRecorder
	Endpoints                    []scope.ServiceEndpoint
	EnableIAM                    bool
	AllowAdditionalRoles         bool
	WatchFilterValue             string
	TagUnmanagedNetworkResources bool
}

// SetupWithManager is used to setup the controller.
func (r *AWSManagedMachinePoolReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := logger.FromContext(ctx)

	gvk, err := apiutil.GVKForObject(new(expinfrav1.AWSManagedMachinePool), mgr.GetScheme())
	if err != nil {
		return errors.Wrapf(err, "failed to find GVK for AWSManagedMachinePool")
	}
	managedControlPlaneToManagedMachinePoolMap := managedControlPlaneToManagedMachinePoolMapFunc(r.Client, gvk, log)
	return ctrl.NewControllerManagedBy(mgr).
		For(&expinfrav1.AWSManagedMachinePool{}).
		WithOptions(options).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(mgr.GetScheme(), log.GetLogger(), r.WatchFilterValue)).
		Watches(
			&expclusterv1.MachinePool{},
			handler.EnqueueRequestsFromMapFunc(machinePoolToInfrastructureMapFunc(gvk)),
		).
		Watches(
			&ekscontrolplanev1.AWSManagedControlPlane{},
			handler.EnqueueRequestsFromMapFunc(managedControlPlaneToManagedMachinePoolMap),
		).
		Complete(r)
}

// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machinepools;machinepools/status,verbs=get;list;watch;patch
// +kubebuilder:rbac:groups=core,resources=events,verbs=get;list;watch;create;patch
// +kubebuilder:rbac:groups=controlplane.cluster.x-k8s.io,resources=awsmanagedcontrolplanes;awsmanagedcontrolplanes/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=awsmanagedmachinepools,verbs=get;list;watch;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=awsmanagedmachinepools/status,verbs=get;update;patch

// Reconcile reconciles AWSManagedMachinePools.
func (r *AWSManagedMachinePoolReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	log := logger.FromContext(ctx)

	awsPool := &expinfrav1.AWSManagedMachinePool{}
	if err := r.Get(ctx, req.NamespacedName, awsPool); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{Requeue: true}, nil
	}

	machinePool, err := getOwnerMachinePool(ctx, r.Client, awsPool.ObjectMeta)
	if err != nil {
		log.Error(err, "Failed to retrieve owner MachinePool from the API Server")
		return ctrl.Result{}, err
	}
	if machinePool == nil {
		log.Info("MachinePool Controller has not yet set OwnerRef")
		return ctrl.Result{}, nil
	}

	log = log.WithValues("MachinePool", klog.KObj(machinePool))

	cluster, err := util.GetClusterFromMetadata(ctx, r.Client, machinePool.ObjectMeta)
	if err != nil {
		log.Info("Failed to retrieve Cluster from MachinePool")
		return reconcile.Result{}, nil
	}

	if annotations.IsPaused(cluster, awsPool) {
		log.Info("Reconciliation is paused for this object")
		return ctrl.Result{}, nil
	}

	log = log.WithValues("cluster", klog.KObj(cluster))

	controlPlaneKey := client.ObjectKey{
		Namespace: awsPool.Namespace,
		Name:      cluster.Spec.ControlPlaneRef.Name,
	}
	controlPlane := &ekscontrolplanev1.AWSManagedControlPlane{}
	if err := r.Client.Get(ctx, controlPlaneKey, controlPlane); err != nil {
		log.Info("Failed to retrieve ControlPlane from MachinePool")
		return reconcile.Result{}, nil
	}

	managedControlPlaneScope, err := scope.NewManagedControlPlaneScope(scope.ManagedControlPlaneScopeParams{
		Client:                       r.Client,
		Logger:                       log,
		Cluster:                      cluster,
		ControlPlane:                 controlPlane,
		ControllerName:               "awsManagedControlPlane",
		TagUnmanagedNetworkResources: r.TagUnmanagedNetworkResources,
	})
	if err != nil {
		return ctrl.Result{}, errors.New("error getting managed control plane scope")
	}

	if !controlPlane.Status.Ready {
		log.Info("Control plane is not ready yet")
		conditions.MarkFalse(awsPool, expinfrav1.EKSNodegroupReadyCondition, expinfrav1.WaitingForEKSControlPlaneReason, clusterv1.ConditionSeverityInfo, "")
		return ctrl.Result{}, nil
	}

	machinePoolScope, err := scope.NewManagedMachinePoolScope(scope.ManagedMachinePoolScopeParams{
		Client:               r.Client,
		ControllerName:       "awsmanagedmachinepool",
		Cluster:              cluster,
		ControlPlane:         controlPlane,
		MachinePool:          machinePool,
		ManagedMachinePool:   awsPool,
		EnableIAM:            r.EnableIAM,
		AllowAdditionalRoles: r.AllowAdditionalRoles,
		Endpoints:            r.Endpoints,
		InfraCluster:         managedControlPlaneScope,
	})
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "failed to create scope")
	}

	defer func() {
		applicableConditions := []clusterv1.ConditionType{
			expinfrav1.EKSNodegroupReadyCondition,
			expinfrav1.IAMNodegroupRolesReadyCondition,
			expinfrav1.LaunchTemplateReadyCondition,
		}

		conditions.SetSummary(machinePoolScope.ManagedMachinePool, conditions.WithConditions(applicableConditions...), conditions.WithStepCounter())

		if err := machinePoolScope.Close(); err != nil && reterr == nil {
			reterr = err
		}
	}()

	if !awsPool.ObjectMeta.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, r.reconcileDelete(ctx, machinePoolScope, managedControlPlaneScope)
	}

	return ctrl.Result{}, r.reconcileNormal(ctx, machinePoolScope, managedControlPlaneScope)
}

func (r *AWSManagedMachinePoolReconciler) reconcileNormal(
	ctx context.Context,
	machinePoolScope *scope.ManagedMachinePoolScope,
	ec2Scope scope.EC2Scope,
) error {
	machinePoolScope.Info("Reconciling AWSManagedMachinePool")

	if controllerutil.AddFinalizer(machinePoolScope.ManagedMachinePool, expinfrav1.ManagedMachinePoolFinalizer) {
		if err := machinePoolScope.PatchObject(); err != nil {
			return err
		}
	}

	ekssvc := eks.NewNodegroupService(machinePoolScope)
	ec2svc := r.getEC2Service(ec2Scope)
	reconSvc := r.getReconcileService(ec2Scope)

	if machinePoolScope.ManagedMachinePool.Spec.AWSLaunchTemplate != nil {
		canUpdateLaunchTemplate := func() (bool, error) {
			return true, nil
		}
		runPostLaunchTemplateUpdateOperation := func() error {
			return nil
		}
		if err := reconSvc.ReconcileLaunchTemplate(machinePoolScope, ec2svc, canUpdateLaunchTemplate, runPostLaunchTemplateUpdateOperation); err != nil {
			r.Recorder.Eventf(machinePoolScope.ManagedMachinePool, corev1.EventTypeWarning, "FailedLaunchTemplateReconcile", "Failed to reconcile launch template: %v", err)
			machinePoolScope.Error(err, "failed to reconcile launch template")
			conditions.MarkFalse(machinePoolScope.ManagedMachinePool, expinfrav1.LaunchTemplateReadyCondition, expinfrav1.LaunchTemplateReconcileFailedReason, clusterv1.ConditionSeverityError, "")
			return err
		}

		launchTemplateID := machinePoolScope.GetLaunchTemplateIDStatus()
		resourceServiceToUpdate := []scope.ResourceServiceToUpdate{{
			ResourceID:      &launchTemplateID,
			ResourceService: ec2svc,
		}}
		if err := reconSvc.ReconcileTags(machinePoolScope, resourceServiceToUpdate); err != nil {
			return errors.Wrap(err, "error updating tags")
		}

		// set the LaunchTemplateReady condition
		conditions.MarkTrue(machinePoolScope.ManagedMachinePool, expinfrav1.LaunchTemplateReadyCondition)
	}

	if err := ekssvc.ReconcilePool(ctx); err != nil {
		return errors.Wrapf(err, "failed to reconcile machine pool for AWSManagedMachinePool %s/%s", machinePoolScope.ManagedMachinePool.Namespace, machinePoolScope.ManagedMachinePool.Name)
	}

	return nil
}

func (r *AWSManagedMachinePoolReconciler) reconcileDelete(
	_ context.Context,
	machinePoolScope *scope.ManagedMachinePoolScope,
	ec2Scope scope.EC2Scope,
) error {
	machinePoolScope.Info("Reconciling deletion of AWSManagedMachinePool")

	ekssvc := eks.NewNodegroupService(machinePoolScope)
	ec2Svc := ec2.NewService(ec2Scope)

	if err := ekssvc.ReconcilePoolDelete(); err != nil {
		return errors.Wrapf(err, "failed to reconcile machine pool deletion for AWSManagedMachinePool %s/%s", machinePoolScope.ManagedMachinePool.Namespace, machinePoolScope.ManagedMachinePool.Name)
	}

	if machinePoolScope.ManagedMachinePool.Spec.AWSLaunchTemplate != nil {
		launchTemplateID := machinePoolScope.ManagedMachinePool.Status.LaunchTemplateID
		launchTemplate, _, _, err := ec2Svc.GetLaunchTemplate(machinePoolScope.LaunchTemplateName())
		if err != nil {
			return err
		}

		if launchTemplate == nil {
			machinePoolScope.Debug("Unable to find matching launch template")
			r.Recorder.Eventf(machinePoolScope.ManagedMachinePool, corev1.EventTypeNormal, "NoLaunchTemplateFound", "Unable to find matching launch template")
			controllerutil.RemoveFinalizer(machinePoolScope.ManagedMachinePool, expinfrav1.ManagedMachinePoolFinalizer)
			return nil
		}

		machinePoolScope.Info("deleting launch template", "name", launchTemplate.Name)
		if err := ec2Svc.DeleteLaunchTemplate(*launchTemplateID); err != nil {
			r.Recorder.Eventf(machinePoolScope.ManagedMachinePool, corev1.EventTypeWarning, "FailedDelete", "Failed to delete launch template %q: %v", launchTemplate.Name, err)
			return errors.Wrap(err, "failed to delete launch template")
		}

		machinePoolScope.Info("successfully deleted launch template")
	}

	controllerutil.RemoveFinalizer(machinePoolScope.ManagedMachinePool, expinfrav1.ManagedMachinePoolFinalizer)

	return nil
}

// GetOwnerClusterKey returns only the Cluster name and namespace.
func GetOwnerClusterKey(obj metav1.ObjectMeta) (*client.ObjectKey, error) {
	for _, ref := range obj.OwnerReferences {
		if ref.Kind != "Cluster" {
			continue
		}
		gv, err := schema.ParseGroupVersion(ref.APIVersion)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if gv.Group == clusterv1.GroupVersion.Group {
			return &client.ObjectKey{
				Namespace: obj.Namespace,
				Name:      ref.Name,
			}, nil
		}
	}
	return nil, nil
}

func managedControlPlaneToManagedMachinePoolMapFunc(c client.Client, gvk schema.GroupVersionKind, log logger.Wrapper) handler.MapFunc {
	return func(ctx context.Context, o client.Object) []reconcile.Request {
		awsControlPlane, ok := o.(*ekscontrolplanev1.AWSManagedControlPlane)
		if !ok {
			klog.Errorf("Expected a AWSManagedControlPlane but got a %T", o)
		}

		if !awsControlPlane.ObjectMeta.DeletionTimestamp.IsZero() {
			return nil
		}

		clusterKey, err := GetOwnerClusterKey(awsControlPlane.ObjectMeta)
		if err != nil {
			log.Error(err, "couldn't get AWS control plane owner ObjectKey")
			return nil
		}
		if clusterKey == nil {
			return nil
		}

		managedPoolForClusterList := expclusterv1.MachinePoolList{}
		if err := c.List(
			ctx, &managedPoolForClusterList, client.InNamespace(clusterKey.Namespace), client.MatchingLabels{clusterv1.ClusterNameLabel: clusterKey.Name},
		); err != nil {
			log.Error(err, "couldn't list pools for cluster")
			return nil
		}

		mapFunc := machinePoolToInfrastructureMapFunc(gvk)

		var results []ctrl.Request
		for i := range managedPoolForClusterList.Items {
			managedPool := mapFunc(ctx, &managedPoolForClusterList.Items[i])
			results = append(results, managedPool...)
		}

		return results
	}
}

func (r *AWSManagedMachinePoolReconciler) getEC2Service(scope scope.EC2Scope) services.EC2Interface {
	return ec2.NewService(scope)
}

func (r *AWSManagedMachinePoolReconciler) getReconcileService(scope scope.EC2Scope) services.MachinePoolReconcileInterface {
	return ec2.NewService(scope)
}
