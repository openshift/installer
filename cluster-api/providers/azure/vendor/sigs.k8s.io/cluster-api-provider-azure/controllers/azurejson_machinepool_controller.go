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
	"fmt"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	"k8s.io/utils/ptr"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/identities"
	infrav1exp "sigs.k8s.io/cluster-api-provider-azure/exp/api/v1beta1"
	azureutil "sigs.k8s.io/cluster-api-provider-azure/util/azure"
	"sigs.k8s.io/cluster-api-provider-azure/util/reconciler"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// AzureJSONMachinePoolReconciler reconciles Azure json secrets for AzureMachinePool objects.
type AzureJSONMachinePoolReconciler struct {
	client.Client
	Recorder         record.EventRecorder
	Timeouts         reconciler.Timeouts
	WatchFilterValue string
}

// SetupWithManager initializes this controller with a manager.
func (r *AzureJSONMachinePoolReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	_, log, done := tele.StartSpanWithLogger(ctx,
		"controllers.AzureJSONMachinePoolReconciler.SetupWithManager",
	)
	defer done()

	azureMachinePoolMapper, err := util.ClusterToTypedObjectsMapper(r.Client, &infrav1exp.AzureMachinePoolList{}, mgr.GetScheme())
	if err != nil {
		return errors.Wrap(err, "failed to create mapper for Cluster to AzureMachinePools")
	}

	c, err := ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&infrav1exp.AzureMachinePool{}).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(log, r.WatchFilterValue)).
		Owns(&corev1.Secret{}).
		Build(r)

	if err != nil {
		return errors.Wrap(err, "failed to create controller")
	}

	// Add a watch on Clusters to requeue when the infraRef is set. This is needed because the infraRef is not initially
	// set in Clusters created from a ClusterClass.
	if err := c.Watch(
		source.Kind(mgr.GetCache(), &clusterv1.Cluster{}),
		handler.EnqueueRequestsFromMapFunc(azureMachinePoolMapper),
		predicates.ClusterUnpausedAndInfrastructureReady(log),
		predicates.ResourceNotPausedAndHasFilterLabel(log, r.WatchFilterValue),
	); err != nil {
		return errors.Wrap(err, "failed adding a watch for Clusters")
	}

	return nil
}

// Reconcile reconciles the Azure json for AzureMachinePool objects.
func (r *AzureJSONMachinePoolReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	ctx, cancel := context.WithTimeout(ctx, r.Timeouts.DefaultedLoopTimeout())
	defer cancel()

	ctx, log, done := tele.StartSpanWithLogger(
		ctx,
		"controllers.AzureJSONMachinePoolReconciler.Reconcile",
		tele.KVP("namespace", req.Namespace),
		tele.KVP("name", req.Name),
		tele.KVP("kind", infrav1.AzureMachinePoolKind),
	)
	defer done()

	log = log.WithValues("namespace", req.Namespace, "azureMachinePool", req.Name)

	// Fetch the AzureMachine instance
	azureMachinePool := &infrav1exp.AzureMachinePool{}
	err := r.Get(ctx, req.NamespacedName, azureMachinePool)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("object was not found")
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// Fetch the CAPI MachinePool.
	machinePool, err := GetOwnerMachinePool(ctx, r.Client, azureMachinePool.ObjectMeta)
	if err != nil {
		return reconcile.Result{}, err
	}
	if machinePool == nil {
		log.Info("MachinePool Controller has not yet set OwnerRef")
		return reconcile.Result{}, nil
	}

	log = log.WithValues("machinePool", machinePool.Name)

	// Fetch the Cluster.
	cluster, err := util.GetClusterFromMetadata(ctx, r.Client, machinePool.ObjectMeta)
	if err != nil {
		log.Info("MachinePool is missing cluster label or cluster does not exist")
		return reconcile.Result{}, nil
	}

	log = log.WithValues("cluster", cluster.Name)

	if cluster.Spec.InfrastructureRef == nil {
		log.Info("infra ref is nil")
		return ctrl.Result{}, nil
	}

	clusterScope, err := GetClusterScoper(ctx, log, r.Client, cluster, r.Timeouts)
	if err != nil {
		return reconcile.Result{}, errors.Wrapf(err, "failed to create cluster scope for cluster %s/%s", cluster.Namespace, cluster.Name)
	}

	// Construct secret for this machine
	userAssignedIdentityIfExists := ""
	if len(azureMachinePool.Spec.UserAssignedIdentities) > 0 {
		var identitiesClient identities.Client
		identitiesClient, err := getClient(clusterScope)
		if err != nil {
			return reconcile.Result{}, errors.Wrap(err, "failed to create identities client")
		}
		parsed, err := azureutil.ParseResourceID(azureMachinePool.Spec.UserAssignedIdentities[0].ProviderID)
		if err != nil {
			return reconcile.Result{}, errors.Wrapf(err, "failed to parse ProviderID %s", azureMachinePool.Spec.UserAssignedIdentities[0].ProviderID)
		}
		if parsed.SubscriptionID != clusterScope.SubscriptionID() {
			identitiesClient, err = identities.NewClientBySub(clusterScope, parsed.SubscriptionID)
			if err != nil {
				return reconcile.Result{}, errors.Wrapf(err, "failed to create identities client from subscription ID %s", parsed.SubscriptionID)
			}
		}
		userAssignedIdentityIfExists, err = identitiesClient.GetClientID(
			ctx, azureMachinePool.Spec.UserAssignedIdentities[0].ProviderID)
		if err != nil {
			return reconcile.Result{}, errors.Wrap(err, "failed to get user-assigned identity ClientID")
		}
	}

	apiVersion, kind := infrav1.GroupVersion.WithKind(infrav1.AzureMachinePoolKind).ToAPIVersionAndKind()
	owner := metav1.OwnerReference{
		APIVersion: apiVersion,
		Kind:       kind,
		Name:       azureMachinePool.GetName(),
		UID:        azureMachinePool.GetUID(),
		Controller: ptr.To(true),
	}

	if azureMachinePool.Spec.Identity == infrav1.VMIdentityNone {
		log.Info(fmt.Sprintf("WARNING, %s", spIdentityWarning))
		r.Recorder.Eventf(azureMachinePool, corev1.EventTypeWarning, "VMIdentityNone", spIdentityWarning)
	}

	newSecret, err := GetCloudProviderSecret(
		clusterScope,
		azureMachinePool.Namespace,
		azureMachinePool.Name,
		owner,
		azureMachinePool.Spec.Identity,
		userAssignedIdentityIfExists,
	)

	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "failed to create cloud provider config")
	}

	if err := reconcileAzureSecret(ctx, r.Client, owner, newSecret, clusterScope.ClusterName()); err != nil {
		r.Recorder.Eventf(azureMachinePool, corev1.EventTypeWarning, "Error reconciling cloud provider secret for AzureMachinePool", err.Error())
		return ctrl.Result{}, errors.Wrap(err, "failed to reconcile azure secret")
	}

	return ctrl.Result{}, nil
}

var getClient = identities.NewClient
