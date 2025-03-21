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
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/eks"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/predicates"
)

// AWSFargateProfileReconciler reconciles a AWSFargateProfile object.
type AWSFargateProfileReconciler struct {
	client.Client
	Recorder         record.EventRecorder
	Endpoints        []scope.ServiceEndpoint
	EnableIAM        bool
	WatchFilterValue string
}

// SetupWithManager is used to setup the controller.
func (r *AWSFargateProfileReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	managedControlPlaneToFargateProfileMap := managedControlPlaneToFargateProfileMapFunc(r.Client, logger.FromContext(ctx))
	return ctrl.NewControllerManagedBy(mgr).
		For(&expinfrav1.AWSFargateProfile{}).
		WithOptions(options).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(mgr.GetScheme(), logger.FromContext(ctx).GetLogger(), r.WatchFilterValue)).
		Watches(
			&ekscontrolplanev1.AWSManagedControlPlane{},
			handler.EnqueueRequestsFromMapFunc(managedControlPlaneToFargateProfileMap),
		).
		Complete(r)
}

// +kubebuilder:rbac:groups=core,resources=events,verbs=get;list;watch;create;patch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=controlplane.cluster.x-k8s.io,resources=awsmanagedcontrolplanes;awsmanagedcontrolplanes/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=awsfargateprofiles,verbs=get;list;watch;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=awsfargateprofiles/status,verbs=get;update;patch

// Reconcile reconciles AWSFargateProfiles.
func (r *AWSFargateProfileReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	log := logger.FromContext(ctx)

	fargateProfile := &expinfrav1.AWSFargateProfile{}
	if err := r.Get(ctx, req.NamespacedName, fargateProfile); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{Requeue: true}, nil
	}

	cluster, err := util.GetClusterByName(ctx, r.Client, fargateProfile.Namespace, fargateProfile.Spec.ClusterName)
	if err != nil {
		log.Info("Failed to retrieve Cluster from AWSFargateProfile")
		return reconcile.Result{}, nil
	}

	log = log.WithValues("cluster", klog.KObj(cluster))

	controlPlaneKey := client.ObjectKey{
		Namespace: fargateProfile.Namespace,
		Name:      cluster.Spec.ControlPlaneRef.Name,
	}
	controlPlane := &ekscontrolplanev1.AWSManagedControlPlane{}
	if err := r.Client.Get(ctx, controlPlaneKey, controlPlane); err != nil {
		log.Info("Failed to retrieve ControlPlane from AWSFargateProfile")
		return reconcile.Result{}, nil
	}

	fargateProfileScope, err := scope.NewFargateProfileScope(scope.FargateProfileScopeParams{
		Client:         r.Client,
		ControllerName: "awsfargateprofile",
		Cluster:        cluster,
		ControlPlane:   controlPlane,
		FargateProfile: fargateProfile,
		EnableIAM:      r.EnableIAM,
		Endpoints:      r.Endpoints,
	})
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "failed to create scope")
	}

	defer func() {
		applicableConditions := []clusterv1.ConditionType{
			expinfrav1.IAMFargateRolesReadyCondition,
			expinfrav1.EKSFargateProfileReadyCondition,
		}

		conditions.SetSummary(fargateProfileScope.FargateProfile, conditions.WithConditions(applicableConditions...), conditions.WithStepCounter())

		if err := fargateProfileScope.Close(); err != nil && reterr == nil {
			reterr = err
		}
	}()

	if !controlPlane.Status.Ready {
		log.Info("Control plane is not ready yet")
		conditions.MarkFalse(fargateProfile, clusterv1.ReadyCondition, expinfrav1.WaitingForEKSControlPlaneReason, clusterv1.ConditionSeverityInfo, "")
		return ctrl.Result{}, nil
	}

	if !fargateProfile.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, fargateProfileScope)
	}

	return r.reconcileNormal(ctx, fargateProfileScope)
}

func (r *AWSFargateProfileReconciler) reconcileNormal(
	_ context.Context,
	fargateProfileScope *scope.FargateProfileScope,
) (ctrl.Result, error) {
	fargateProfileScope.Info("Reconciling AWSFargateProfile")

	if controllerutil.AddFinalizer(fargateProfileScope.FargateProfile, expinfrav1.FargateProfileFinalizer) {
		if err := fargateProfileScope.PatchObject(); err != nil {
			return ctrl.Result{}, err
		}
	}

	ekssvc := eks.NewFargateService(fargateProfileScope)

	res, err := ekssvc.Reconcile()
	if err != nil {
		return res, errors.Wrapf(err, "failed to reconcile fargate profile for AWSFargateProfile %s/%s", fargateProfileScope.FargateProfile.Namespace, fargateProfileScope.FargateProfile.Name)
	}

	return res, nil
}

func (r *AWSFargateProfileReconciler) reconcileDelete(
	_ context.Context,
	fargateProfileScope *scope.FargateProfileScope,
) (ctrl.Result, error) {
	fargateProfileScope.Info("Reconciling deletion of AWSFargateProfile")

	ekssvc := eks.NewFargateService(fargateProfileScope)

	res, err := ekssvc.ReconcileDelete()
	if err != nil {
		return res, errors.Wrapf(err, "failed to reconcile fargate profile deletion for AWSFargateProfile %s/%s", fargateProfileScope.FargateProfile.Namespace, fargateProfileScope.FargateProfile.Name)
	}

	if res.IsZero() {
		controllerutil.RemoveFinalizer(fargateProfileScope.FargateProfile, expinfrav1.FargateProfileFinalizer)
	}

	return res, nil
}

func managedControlPlaneToFargateProfileMapFunc(c client.Client, log logger.Wrapper) handler.MapFunc {
	return func(ctx context.Context, o client.Object) []ctrl.Request {
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

		fargateProfileForClusterList := expinfrav1.AWSFargateProfileList{}
		if err := c.List(
			ctx, &fargateProfileForClusterList, client.InNamespace(clusterKey.Namespace), client.MatchingLabels{clusterv1.ClusterNameLabel: clusterKey.Name},
		); err != nil {
			log.Error(err, "couldn't list fargate profiles for cluster")
			return nil
		}

		var results []ctrl.Request
		for i := range fargateProfileForClusterList.Items {
			fp := fargateProfileForClusterList.Items[i]
			results = append(results, reconcile.Request{
				NamespacedName: client.ObjectKey{
					Namespace: fp.Namespace,
					Name:      fp.Name,
				},
			})
		}

		return results
	}
}
