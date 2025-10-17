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

// Package controllers provides a way to reconcile ROSA resources.
package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	stsv2 "github.com/aws/aws-sdk-go-v2/service/sts"
	sts "github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/google/go-cmp/cmp"
	idputils "github.com/openshift-online/ocm-common/pkg/idp/utils"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	rosaaws "github.com/openshift/rosa/pkg/aws"
	"github.com/openshift/rosa/pkg/ocm"
	"github.com/zgalor/weberr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apiserver/pkg/storage/names"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	rosacontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/rosa/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/annotations"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/rosa"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/utils"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	capiannotations "sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/kubeconfig"
	"sigs.k8s.io/cluster-api/util/predicates"
	"sigs.k8s.io/cluster-api/util/secret"
)

const (
	rosaControlPlaneKind = "ROSAControlPlane"
	// ROSAControlPlaneFinalizer allows the controller to clean up resources on delete.
	ROSAControlPlaneFinalizer = "rosacontrolplane.controlplane.cluster.x-k8s.io"

	// ROSAControlPlaneForceDeleteAnnotation annotation can be set to force the deletion of ROSAControlPlane bypassing any deletion validations/errors.
	ROSAControlPlaneForceDeleteAnnotation = "controlplane.cluster.x-k8s.io/rosacontrolplane-force-delete"

	// ExternalAuthProviderLastAppliedAnnotation annotation tracks the last applied external auth configuration to inform if an update is required.
	ExternalAuthProviderLastAppliedAnnotation = "controlplane.cluster.x-k8s.io/rosacontrolplane-last-applied-external-auth-provider"
)

// ROSAControlPlaneReconciler reconciles a ROSAControlPlane object.
type ROSAControlPlaneReconciler struct {
	client.Client
	WatchFilterValue string
	WaitInfraPeriod  time.Duration
	Endpoints        []scope.ServiceEndpoint
	NewStsClient     func(cloud.ScopeUsage, cloud.Session, logger.Wrapper, runtime.Object) stsiface.STSAPI
	NewOCMClient     func(ctx context.Context, rosaScope *scope.ROSAControlPlaneScope) (rosa.OCMClient, error)
}

// SetupWithManager is used to setup the controller.
func (r *ROSAControlPlaneReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := logger.FromContext(ctx)
	r.NewOCMClient = rosa.NewWrappedOCMClient
	r.NewStsClient = scope.NewSTSClient

	rosaControlPlane := &rosacontrolplanev1.ROSAControlPlane{}
	c, err := ctrl.NewControllerManagedBy(mgr).
		For(rosaControlPlane).
		WithOptions(options).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(mgr.GetScheme(), log.GetLogger(), r.WatchFilterValue)).
		Build(r)

	if err != nil {
		return fmt.Errorf("failed setting up the AWSManagedControlPlane controller manager: %w", err)
	}

	if err = c.Watch(
		source.Kind[client.Object](mgr.GetCache(), &clusterv1.Cluster{},
			handler.EnqueueRequestsFromMapFunc(util.ClusterToInfrastructureMapFunc(ctx, rosaControlPlane.GroupVersionKind(), mgr.GetClient(), &expinfrav1.ROSACluster{})),
			predicates.ClusterPausedTransitionsOrInfrastructureReady(mgr.GetScheme(), log.GetLogger())),
	); err != nil {
		return fmt.Errorf("failed adding a watch for ready clusters: %w", err)
	}

	if err = c.Watch(
		source.Kind[client.Object](mgr.GetCache(), &expinfrav1.ROSACluster{},
			handler.EnqueueRequestsFromMapFunc(r.rosaClusterToROSAControlPlane(log))),
	); err != nil {
		return fmt.Errorf("failed adding a watch for ROSACluster")
	}

	return nil
}

// +kubebuilder:rbac:groups=core,resources=events,verbs=get;list;watch;create;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;delete;patch
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;delete;patch
// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machinedeployments,verbs=get;list;watch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machinepools,verbs=get;list;watch
// +kubebuilder:rbac:groups=controlplane.cluster.x-k8s.io,resources=rosacontrolplanes,verbs=get;list;watch;update;patch;delete
// +kubebuilder:rbac:groups=controlplane.cluster.x-k8s.io,resources=rosacontrolplanes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=controlplane.cluster.x-k8s.io,resources=rosacontrolplanes/finalizers,verbs=update

// Reconcile will reconcile RosaControlPlane Resources.
func (r *ROSAControlPlaneReconciler) Reconcile(ctx context.Context, req ctrl.Request) (res ctrl.Result, reterr error) {
	log := logger.FromContext(ctx)

	// Get the control plane instance
	rosaControlPlane := &rosacontrolplanev1.ROSAControlPlane{}
	if err := r.Client.Get(ctx, req.NamespacedName, rosaControlPlane); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// Get the cluster
	cluster, err := util.GetOwnerCluster(ctx, r.Client, rosaControlPlane.ObjectMeta)
	if err != nil {
		log.Error(err, "Failed to retrieve owner Cluster from the API Server")
		return ctrl.Result{}, err
	}
	if cluster == nil {
		log.Info("Cluster Controller has not yet set OwnerRef")
		return ctrl.Result{}, nil
	}

	log = log.WithValues("cluster", klog.KObj(cluster))

	if capiannotations.IsPaused(cluster, rosaControlPlane) {
		log.Info("Reconciliation is paused for this object")
		return ctrl.Result{}, nil
	}

	rosaScope, err := scope.NewROSAControlPlaneScope(scope.ROSAControlPlaneScopeParams{
		Client:         r.Client,
		Cluster:        cluster,
		ControlPlane:   rosaControlPlane,
		ControllerName: strings.ToLower(rosaControlPlaneKind),
		Endpoints:      r.Endpoints,
		Logger:         log,
		NewStsClient:   r.NewStsClient,
	})
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create scope: %w", err)
	}

	// Always close the scope
	defer func() {
		if err := rosaScope.Close(); err != nil {
			reterr = errors.Join(reterr, err)
		}
	}()

	if !rosaControlPlane.ObjectMeta.DeletionTimestamp.IsZero() {
		// Handle deletion reconciliation loop.
		return r.reconcileDelete(ctx, rosaScope)
	}

	// Handle normal reconciliation loop.
	return r.reconcileNormal(ctx, rosaScope)
}

func (r *ROSAControlPlaneReconciler) reconcileNormal(ctx context.Context, rosaScope *scope.ROSAControlPlaneScope) (res ctrl.Result, reterr error) {
	rosaScope.Info("Reconciling ROSAControlPlane")

	if controllerutil.AddFinalizer(rosaScope.ControlPlane, ROSAControlPlaneFinalizer) {
		if err := rosaScope.PatchObject(); err != nil {
			return ctrl.Result{}, err
		}
	}
	if r.NewOCMClient == nil {
		return ctrl.Result{}, fmt.Errorf("failed to create OCM client: NewOCMClient is nil")
	}

	ocmClient, err := r.NewOCMClient(ctx, rosaScope)
	if err != nil || ocmClient == nil {
		// TODO: need to expose in status, as likely the credentials are invalid
		return ctrl.Result{}, fmt.Errorf("failed to create OCM client: %w", err)
	}

	creator, err := rosaaws.CreatorForCallerIdentity(convertStsV2(rosaScope.Identity))
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to transform caller identity to creator: %w", err)
	}

	validationMessage, err := validateControlPlaneSpec(ocmClient, rosaScope)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to validate ROSAControlPlane.spec: %w", err)
	}

	conditions.MarkTrue(rosaScope.ControlPlane, rosacontrolplanev1.ROSAControlPlaneValidCondition)
	if validationMessage != "" {
		conditions.MarkFalse(rosaScope.ControlPlane,
			rosacontrolplanev1.ROSAControlPlaneValidCondition,
			rosacontrolplanev1.ROSAControlPlaneInvalidConfigurationReason,
			clusterv1.ConditionSeverityError,
			"%s",
			validationMessage)
		// dont' requeue because input is invalid and manual intervention is needed.
		return ctrl.Result{}, nil
	}
	rosaScope.ControlPlane.Status.FailureMessage = nil

	cluster, err := ocmClient.GetCluster(rosaScope.ControlPlane.Spec.RosaClusterName, creator)
	if err != nil && weberr.GetType(err) != weberr.NotFound {
		return ctrl.Result{}, err
	}

	if cluster != nil {
		rosaScope.ControlPlane.Status.ID = cluster.ID()
		rosaScope.ControlPlane.Status.ConsoleURL = cluster.Console().URL()
		rosaScope.ControlPlane.Status.OIDCEndpointURL = cluster.AWS().STS().OIDCEndpointURL()
		rosaScope.ControlPlane.Status.Ready = false

		switch cluster.Status().State() {
		case cmv1.ClusterStateReady:
			conditions.MarkTrue(rosaScope.ControlPlane, rosacontrolplanev1.ROSAControlPlaneReadyCondition)
			rosaScope.ControlPlane.Status.Ready = true

			apiEndpoint, err := buildAPIEndpoint(cluster)
			if err != nil {
				return ctrl.Result{}, err
			}
			rosaScope.ControlPlane.Spec.ControlPlaneEndpoint = *apiEndpoint

			if err := r.updateOCMCluster(rosaScope, ocmClient, cluster, creator); err != nil {
				return ctrl.Result{}, fmt.Errorf("failed to update rosa control plane: %w", err)
			}
			if err := r.reconcileClusterVersion(rosaScope, ocmClient, cluster); err != nil {
				return ctrl.Result{}, err
			}

			if rosaScope.ControlPlane.Spec.EnableExternalAuthProviders {
				if err := r.reconcileExternalAuth(ctx, rosaScope, cluster); err != nil {
					return ctrl.Result{}, fmt.Errorf("failed to reconcile external auth: %w", err)
				}
			} else {
				// only reconcile a kubeconfig when external auth is not enabled.
				// The user is expected to provide the kubeconfig for CAPI.
				if err := r.reconcileKubeconfig(ctx, rosaScope, ocmClient, cluster); err != nil {
					return ctrl.Result{}, fmt.Errorf("failed to reconcile kubeconfig: %w", err)
				}
			}

			return ctrl.Result{}, nil
		case cmv1.ClusterStateError:
			errorMessage := cluster.Status().ProvisionErrorMessage()
			rosaScope.ControlPlane.Status.FailureMessage = &errorMessage

			conditions.MarkFalse(rosaScope.ControlPlane,
				rosacontrolplanev1.ROSAControlPlaneReadyCondition,
				string(cluster.Status().State()),
				clusterv1.ConditionSeverityError,
				"%s",
				cluster.Status().ProvisionErrorCode())
			// Cluster is in an unrecoverable state, returning nil error so that the request doesn't get requeued.
			return ctrl.Result{}, nil
		}

		conditions.MarkFalse(rosaScope.ControlPlane,
			rosacontrolplanev1.ROSAControlPlaneReadyCondition,
			string(cluster.Status().State()),
			clusterv1.ConditionSeverityInfo,
			"%s",
			cluster.Status().Description())

		rosaScope.Info("waiting for cluster to become ready", "state", cluster.Status().State())
		// Requeue so that status.ready is set to true when the cluster is fully created.
		return ctrl.Result{RequeueAfter: time.Second * 60}, nil
	}

	ocmClusterSpec, err := buildOCMClusterSpec(rosaScope.ControlPlane.Spec, creator)
	if err != nil {
		return ctrl.Result{}, err
	}

	cluster, err = ocmClient.CreateCluster(ocmClusterSpec)
	if err != nil {
		conditions.MarkFalse(rosaScope.ControlPlane,
			rosacontrolplanev1.ROSAControlPlaneReadyCondition,
			rosacontrolplanev1.ReconciliationFailedReason,
			clusterv1.ConditionSeverityError,
			"%s",
			err.Error())
		return ctrl.Result{}, fmt.Errorf("failed to create OCM cluster: %w", err)
	}

	rosaScope.Info("cluster created", "state", cluster.Status().State())
	rosaScope.ControlPlane.Status.ID = cluster.ID()

	return ctrl.Result{}, nil
}

func (r *ROSAControlPlaneReconciler) reconcileDelete(ctx context.Context, rosaScope *scope.ROSAControlPlaneScope) (res ctrl.Result, reterr error) {
	rosaScope.Info("Reconciling ROSAControlPlane delete")

	// Deleting MachinePools first.
	deleted, err := r.deleteMachinePools(ctx, rosaScope)
	if err != nil {
		return ctrl.Result{}, err
	}
	if !deleted {
		// Reconcile after 1 min giving time for machinePools to be deleted.
		return ctrl.Result{RequeueAfter: time.Minute}, nil
	}

	ocmClient, err := rosa.NewOCMClient(ctx, rosaScope)
	if err != nil || ocmClient == nil {
		// TODO: need to expose in status, as likely the credentials are invalid
		return ctrl.Result{}, fmt.Errorf("failed to create OCM client: %w", err)
	}

	creator, err := rosaaws.CreatorForCallerIdentity(convertStsV2(rosaScope.Identity))
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to transform caller identity to creator: %w", err)
	}

	cluster, err := ocmClient.GetCluster(rosaScope.ControlPlane.Spec.RosaClusterName, creator)
	if err != nil && weberr.GetType(err) != weberr.NotFound {
		return ctrl.Result{}, err
	}
	if cluster == nil {
		// cluster and machinepools are deleted, removing finalizer.
		controllerutil.RemoveFinalizer(rosaScope.ControlPlane, ROSAControlPlaneFinalizer)

		return ctrl.Result{}, nil
	}

	bestEffort := false
	if value, found := annotations.Get(rosaScope.ControlPlane, ROSAControlPlaneForceDeleteAnnotation); found && value != "false" {
		bestEffort = true
	}

	if cluster.Status().State() != cmv1.ClusterStateUninstalling {
		if _, err := ocmClient.DeleteCluster(cluster.ID(), bestEffort, creator); err != nil {
			conditions.MarkFalse(rosaScope.ControlPlane,
				rosacontrolplanev1.ROSAControlPlaneReadyCondition,
				rosacontrolplanev1.ROSAControlPlaneDeletionFailedReason,
				clusterv1.ConditionSeverityError,
				"failed to delete ROSAControlPlane: %s; if the error can't be resolved, set '%s' annotation to force the deletion",
				err.Error(),
				ROSAControlPlaneForceDeleteAnnotation)
			return ctrl.Result{}, err
		}
	}

	conditions.MarkFalse(rosaScope.ControlPlane,
		rosacontrolplanev1.ROSAControlPlaneReadyCondition,
		string(cluster.Status().State()),
		clusterv1.ConditionSeverityInfo,
		"deleting")
	rosaScope.ControlPlane.Status.Ready = false
	rosaScope.Info("waiting for cluster to be deleted")
	// Requeue to remove the finalizer when the cluster is fully deleted.
	return ctrl.Result{RequeueAfter: time.Second * 60}, nil
}

// deleteMachinePools check if the controlplane has related machinePools and delete them.
func (r *ROSAControlPlaneReconciler) deleteMachinePools(ctx context.Context, rosaScope *scope.ROSAControlPlaneScope) (bool, error) {
	machinePools, err := utils.GetMachinePools(ctx, rosaScope.Client, rosaScope.Cluster.Name, rosaScope.Cluster.Namespace)
	if err != nil {
		return false, err
	}

	var errs []error
	for id, mp := range machinePools {
		if !mp.DeletionTimestamp.IsZero() {
			continue
		}
		if err = rosaScope.Client.Delete(ctx, &machinePools[id]); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return false, kerrors.NewAggregate(errs)
	}

	return len(machinePools) == 0, nil
}

func (r *ROSAControlPlaneReconciler) reconcileClusterVersion(rosaScope *scope.ROSAControlPlaneScope, ocmClient rosa.OCMClient, cluster *cmv1.Cluster) error {
	version := rosaScope.ControlPlane.Spec.Version
	if version == rosa.RawVersionID(cluster.Version()) {
		conditions.MarkFalse(rosaScope.ControlPlane, rosacontrolplanev1.ROSAControlPlaneUpgradingCondition, "upgraded", clusterv1.ConditionSeverityInfo, "")

		if cluster.Version() != nil {
			rosaScope.ControlPlane.Status.AvailableUpgrades = cluster.Version().AvailableUpgrades()
		}

		// Set the version gate to WaitForAcknowledge as the previous upgrade is applied.
		if rosaScope.ControlPlane.Spec.VersionGate == rosacontrolplanev1.Acknowledge {
			rosaScope.ControlPlane.Spec.VersionGate = rosacontrolplanev1.WaitForAcknowledge
		}

		// return as there is no upgrade to schedule.
		return nil
	}

	scheduledUpgrade, err := rosa.CheckExistingScheduledUpgrade(ocmClient, cluster)
	if err != nil {
		return fmt.Errorf("failed to get existing scheduled upgrades: %w", err)
	}

	if scheduledUpgrade == nil {
		ack := (rosaScope.ControlPlane.Spec.VersionGate == rosacontrolplanev1.Acknowledge || rosaScope.ControlPlane.Spec.VersionGate == rosacontrolplanev1.AlwaysAcknowledge)
		scheduledUpgrade, err = rosa.ScheduleControlPlaneUpgrade(ocmClient, cluster, version, time.Now(), ack)
		if err != nil {
			condition := &clusterv1.Condition{
				Type:    rosacontrolplanev1.ROSAControlPlaneUpgradingCondition,
				Status:  corev1.ConditionFalse,
				Reason:  "failed",
				Message: fmt.Sprintf("failed to schedule upgrade to version %s: %v", version, err),
			}
			conditions.Set(rosaScope.ControlPlane, condition)

			return err
		}
	}

	condition := &clusterv1.Condition{
		Type:    rosacontrolplanev1.ROSAControlPlaneUpgradingCondition,
		Status:  corev1.ConditionTrue,
		Reason:  string(scheduledUpgrade.State().Value()),
		Message: fmt.Sprintf("Upgrading to version %s", scheduledUpgrade.Version()),
	}
	conditions.Set(rosaScope.ControlPlane, condition)

	// if cluster is already upgrading to another version we need to wait until the current upgrade is finished, return an error to requeue and try later.
	if scheduledUpgrade.Version() != version {
		return fmt.Errorf("there is already a %s upgrade to version %s", scheduledUpgrade.State().Value(), scheduledUpgrade.Version())
	}

	return nil
}

func (r *ROSAControlPlaneReconciler) updateOCMCluster(rosaScope *scope.ROSAControlPlaneScope, ocmClient rosa.OCMClient, cluster *cmv1.Cluster, creator *rosaaws.Creator) error {
	ocmClusterSpec, updated := r.updateOCMClusterSpec(rosaScope.ControlPlane, cluster)

	if updated {
		// Update the cluster.
		rosaScope.Info("Updating cluster")
		if err := ocmClient.UpdateCluster(cluster.ID(), creator, ocmClusterSpec); err != nil {
			conditions.MarkFalse(rosaScope.ControlPlane,
				rosacontrolplanev1.ROSAControlPlaneValidCondition,
				rosacontrolplanev1.ROSAControlPlaneInvalidConfigurationReason,
				clusterv1.ConditionSeverityError,
				"%s",
				err.Error())
			return err
		}
	}

	return nil
}

func (r *ROSAControlPlaneReconciler) updateOCMClusterSpec(rosaControlPlane *rosacontrolplanev1.ROSAControlPlane, cluster *cmv1.Cluster) (ocm.Spec, bool) {
	ocmClusterSpec := ocm.Spec{}
	updated := false

	// Check for audit role arn changes
	currentAuditLogRole := cluster.AWS().AuditLog().RoleArn()
	if currentAuditLogRole != rosaControlPlane.Spec.AuditLogRoleARN {
		ocmClusterSpec.AuditLogRoleARN = ptr.To(rosaControlPlane.Spec.AuditLogRoleARN)
		updated = true
	}

	// Check for registry config changes
	regConfig := &rosacontrolplanev1.RegistryConfig{
		RegistrySources: &rosacontrolplanev1.RegistrySources{},
	}
	if rosaControlPlane.Spec.ClusterRegistryConfig != nil {
		regConfig.AdditionalTrustedCAs = rosaControlPlane.Spec.ClusterRegistryConfig.AdditionalTrustedCAs
		regConfig.AllowedRegistriesForImport = rosaControlPlane.Spec.ClusterRegistryConfig.AllowedRegistriesForImport

		if rosaControlPlane.Spec.ClusterRegistryConfig.RegistrySources != nil {
			regConfig.RegistrySources.AllowedRegistries = rosaControlPlane.Spec.ClusterRegistryConfig.RegistrySources.AllowedRegistries
			regConfig.RegistrySources.BlockedRegistries = rosaControlPlane.Spec.ClusterRegistryConfig.RegistrySources.BlockedRegistries
			regConfig.RegistrySources.InsecureRegistries = rosaControlPlane.Spec.ClusterRegistryConfig.RegistrySources.InsecureRegistries
		}
	}
	if !reflect.DeepEqual(regConfig.AdditionalTrustedCAs, cluster.RegistryConfig().AdditionalTrustedCa()) {
		ocmClusterSpec.AdditionalTrustedCa = regConfig.AdditionalTrustedCAs
		updated = true
	}
	if !reflect.DeepEqual(regConfig.RegistrySources.AllowedRegistries, cluster.RegistryConfig().RegistrySources().AllowedRegistries()) {
		ocmClusterSpec.AllowedRegistries = regConfig.RegistrySources.AllowedRegistries
		updated = true
	}
	if !reflect.DeepEqual(regConfig.RegistrySources.BlockedRegistries, cluster.RegistryConfig().RegistrySources().BlockedRegistries()) {
		ocmClusterSpec.BlockedRegistries = regConfig.RegistrySources.BlockedRegistries
		updated = true
	}
	if !reflect.DeepEqual(regConfig.RegistrySources.InsecureRegistries, cluster.RegistryConfig().RegistrySources().InsecureRegistries()) {
		ocmClusterSpec.InsecureRegistries = regConfig.RegistrySources.InsecureRegistries
		updated = true
	}

	var newAllowedRegistries, oldAllowedRegistries []string
	if len(regConfig.AllowedRegistriesForImport) > 0 {
		for id := range regConfig.AllowedRegistriesForImport {
			newAllowedRegistries = append(newAllowedRegistries, regConfig.AllowedRegistriesForImport[id].DomainName+":"+
				strconv.FormatBool(regConfig.AllowedRegistriesForImport[id].Insecure))
		}
	}
	if len(cluster.RegistryConfig().AllowedRegistriesForImport()) > 0 {
		for id := range cluster.RegistryConfig().AllowedRegistriesForImport() {
			oldAllowedRegistries = append(oldAllowedRegistries, cluster.RegistryConfig().AllowedRegistriesForImport()[id].DomainName()+":"+
				strconv.FormatBool(cluster.RegistryConfig().AllowedRegistriesForImport()[id].Insecure()))
		}
	}
	if !reflect.DeepEqual(newAllowedRegistries, oldAllowedRegistries) {
		ocmClusterSpec.AllowedRegistriesForImport = strings.Join(newAllowedRegistries, ",")
		updated = true
	}

	// TODO: check for cluster AutoScale changes
	// rosaControlPlane.Spec.DefaultMachinePoolSpec.Autoscaling

	return ocmClusterSpec, updated
}

func (r *ROSAControlPlaneReconciler) reconcileExternalAuth(ctx context.Context, rosaScope *scope.ROSAControlPlaneScope, cluster *cmv1.Cluster) error {
	externalAuthClient, err := rosa.NewExternalAuthClient(ctx, rosaScope)
	if err != nil {
		return fmt.Errorf("failed to create external auth client: %v", err)
	}
	defer externalAuthClient.Close()

	var errs []error
	if err := r.reconcileExternalAuthProviders(ctx, externalAuthClient, rosaScope, cluster); err != nil {
		errs = append(errs, err)
		conditions.MarkFalse(rosaScope.ControlPlane,
			rosacontrolplanev1.ExternalAuthConfiguredCondition,
			rosacontrolplanev1.ReconciliationFailedReason,
			clusterv1.ConditionSeverityError,
			"%s",
			err.Error())
	} else {
		conditions.MarkTrue(rosaScope.ControlPlane, rosacontrolplanev1.ExternalAuthConfiguredCondition)
	}

	if err := r.reconcileExternalAuthBootstrapKubeconfig(ctx, externalAuthClient, rosaScope, cluster); err != nil {
		errs = append(errs, err)
	}

	return kerrors.NewAggregate(errs)
}

func (r *ROSAControlPlaneReconciler) reconcileExternalAuthProviders(ctx context.Context, externalAuthClient *rosa.ExternalAuthClient, rosaScope *scope.ROSAControlPlaneScope, cluster *cmv1.Cluster) error {
	externalAuths, err := externalAuthClient.ListExternalAuths(cluster.ID())
	if err != nil {
		return fmt.Errorf("failed to list external auths: %v", err)
	}

	if len(rosaScope.ControlPlane.Spec.ExternalAuthProviders) == 0 {
		if len(externalAuths) > 0 {
			if err := externalAuthClient.DeleteExternalAuth(cluster.ID(), externalAuths[0].ID()); err != nil {
				return fmt.Errorf("failed to delete external auth provider %s: %v", externalAuths[0].ID(), err)
			}
		}

		return nil
	}

	authProvider := rosaScope.ControlPlane.Spec.ExternalAuthProviders[0]
	shouldUpdate := false
	if len(externalAuths) > 0 {
		existingProvider := externalAuths[0]
		// name/ID can't be patched, we need to delete the old provider and create a new one.
		if existingProvider.ID() != authProvider.Name {
			if err := externalAuthClient.DeleteExternalAuth(cluster.ID(), existingProvider.ID()); err != nil {
				return fmt.Errorf("failed to delete external auth provider %s: %v", existingProvider.ID(), err)
			}
		} else {
			jsonAnnotation := rosaScope.ControlPlane.Annotations[ExternalAuthProviderLastAppliedAnnotation]
			if len(jsonAnnotation) != 0 {
				var lastAppliedAuthProvider rosacontrolplanev1.ExternalAuthProvider
				err := json.Unmarshal([]byte(jsonAnnotation), &lastAppliedAuthProvider)
				if err != nil {
					return fmt.Errorf("failed to unmarshal '%s' annotaion content: %v", ExternalAuthProviderLastAppliedAnnotation, err)
				}

				// if there were no changes, return.
				if cmp.Equal(authProvider, lastAppliedAuthProvider) {
					return nil
				}
			}

			shouldUpdate = true
		}
	}

	externalAuthBuilder := cmv1.NewExternalAuth().ID(authProvider.Name)

	// issuer builder
	audiences := make([]string, 0, len(authProvider.Issuer.Audiences))
	for _, a := range authProvider.Issuer.Audiences {
		audiences = append(audiences, string(a))
	}
	tokenIssuerBuilder := cmv1.NewTokenIssuer().URL(authProvider.Issuer.URL).
		Audiences(audiences...)

	if authProvider.Issuer.CertificateAuthority != nil {
		CertificateAuthorityConfigMap := &corev1.ConfigMap{}
		err := rosaScope.Client.Get(ctx, types.NamespacedName{Namespace: rosaScope.Namespace(), Name: authProvider.Issuer.CertificateAuthority.Name}, CertificateAuthorityConfigMap)
		if err != nil {
			return fmt.Errorf("failed to get issuer CertificateAuthority configMap %s: %v", authProvider.Issuer.CertificateAuthority.Name, err)
		}
		CertificateAuthorityValue := CertificateAuthorityConfigMap.Data["ca-bundle.crt"]

		tokenIssuerBuilder.CA(CertificateAuthorityValue)
	}
	externalAuthBuilder.Issuer(tokenIssuerBuilder)

	// oidc-clients builder
	clientsBuilders := make([]*cmv1.ExternalAuthClientConfigBuilder, 0, len(authProvider.OIDCClients))
	for _, client := range authProvider.OIDCClients {
		secretObj := &corev1.Secret{}
		err := rosaScope.Client.Get(ctx, types.NamespacedName{Namespace: rosaScope.Namespace(), Name: client.ClientSecret.Name}, secretObj)
		if err != nil {
			return fmt.Errorf("failed to get client secret %s: %v", client.ClientSecret.Name, err)
		}
		clientSecretValue := string(secretObj.Data["clientSecret"])

		clientsBuilders = append(clientsBuilders, cmv1.NewExternalAuthClientConfig().
			ID(client.ClientID).Secret(clientSecretValue).
			Component(cmv1.NewClientComponent().Name(client.ComponentName).Namespace(client.ComponentNamespace)))
	}
	externalAuthBuilder.Clients(clientsBuilders...)

	// claims builder
	if authProvider.ClaimMappings != nil {
		clainMappingsBuilder := cmv1.NewTokenClaimMappings()
		if authProvider.ClaimMappings.Groups != nil {
			clainMappingsBuilder.Groups(cmv1.NewGroupsClaim().Claim(authProvider.ClaimMappings.Groups.Claim).
				Prefix(authProvider.ClaimMappings.Groups.Prefix))
		}

		if authProvider.ClaimMappings.Username != nil {
			usernameClaimBuilder := cmv1.NewUsernameClaim().Claim(authProvider.ClaimMappings.Username.Claim).
				PrefixPolicy(string(authProvider.ClaimMappings.Username.PrefixPolicy))
			if authProvider.ClaimMappings.Username.Prefix != nil {
				usernameClaimBuilder.Prefix(*authProvider.ClaimMappings.Username.Prefix)
			}

			clainMappingsBuilder.UserName(usernameClaimBuilder)
		}

		claimBuilder := cmv1.NewExternalAuthClaim().Mappings(clainMappingsBuilder)

		validationRulesbuilders := make([]*cmv1.TokenClaimValidationRuleBuilder, 0, len(authProvider.ClaimValidationRules))
		for _, rule := range authProvider.ClaimValidationRules {
			validationRulesbuilders = append(validationRulesbuilders, cmv1.NewTokenClaimValidationRule().
				Claim(rule.RequiredClaim.Claim).RequiredValue(rule.RequiredClaim.RequiredValue))
		}
		claimBuilder.ValidationRules(validationRulesbuilders...)

		externalAuthBuilder.Claim(claimBuilder)
	}

	externalAuthConfig, err := externalAuthBuilder.Build()
	if err != nil {
		return fmt.Errorf("failed to build external auth config: %v", err)
	}

	if shouldUpdate {
		_, err = externalAuthClient.UpdateExternalAuth(cluster.ID(), externalAuthConfig)
		if err != nil {
			return fmt.Errorf("failed to update external authentication provider '%s' for cluster '%s': %v",
				externalAuthConfig.ID(), rosaScope.InfraClusterName(), err)
		}
	} else {
		_, err = externalAuthClient.CreateExternalAuth(cluster.ID(), externalAuthConfig)
		if err != nil {
			return fmt.Errorf("failed to create external authentication provider '%s' for cluster '%s': %v",
				externalAuthConfig.ID(), rosaScope.InfraClusterName(), err)
		}
	}

	lastAppliedAnnotation, err := json.Marshal(authProvider)
	if err != nil {
		return err
	}

	if rosaScope.ControlPlane.Annotations == nil {
		rosaScope.ControlPlane.Annotations = make(map[string]string)
	}
	rosaScope.ControlPlane.Annotations[ExternalAuthProviderLastAppliedAnnotation] = string(lastAppliedAnnotation)

	return nil
}

// Generates a temporarily admin kubeconfig using break-glass credentials for the user to bootstreap their environment like setting up RBAC for oidc users/groups.
// This Kubeonconfig will be created only once initially and be valid for only 24h.
// The kubeconfig secret will not be autoamticallty rotated and will be invalid after the 24h. However, users can opt to manually delete the secret to trigger the generation of a new one which will be valid for another 24h.
func (r *ROSAControlPlaneReconciler) reconcileExternalAuthBootstrapKubeconfig(ctx context.Context, externalAuthClient *rosa.ExternalAuthClient, rosaScope *scope.ROSAControlPlaneScope, cluster *cmv1.Cluster) error {
	kubeconfigSecret := rosaScope.ExternalAuthBootstrapKubeconfigSecret()
	err := r.Client.Get(ctx, client.ObjectKeyFromObject(kubeconfigSecret), kubeconfigSecret)
	if err == nil {
		// already exist.
		return nil
	} else if !apierrors.IsNotFound(err) {
		return fmt.Errorf("failed to get bootstrap kubeconfig secret: %w", err)
	}

	// kubeconfig doesn't exist, generate a new one.
	breakGlassConfig, err := cmv1.NewBreakGlassCredential().
		Username(names.SimpleNameGenerator.GenerateName("capi-admin-")). // OCM requires unique usernames
		ExpirationTimestamp(time.Now().Add(time.Hour * 24)).
		Build()
	if err != nil {
		return fmt.Errorf("failed to build break glass config: %v", err)
	}

	breakGlassCredential, err := externalAuthClient.CreateBreakGlassCredential(cluster.ID(), breakGlassConfig)
	if err != nil {
		return fmt.Errorf("failed to create break glass credential: %v", err)
	}

	kubeconfigData, err := externalAuthClient.PollKubeconfig(ctx, cluster.ID(), breakGlassCredential.ID())
	if err != nil {
		return fmt.Errorf("failed to poll break glass kubeconfig: %v", err)
	}

	kubeconfigSecret.Data = map[string][]byte{
		"value": []byte(kubeconfigData),
	}
	if err := r.Client.Create(ctx, kubeconfigSecret); err != nil {
		return fmt.Errorf("failed to create external auth bootstrap kubeconfig: %v", err)
	}

	return nil
}

func (r *ROSAControlPlaneReconciler) reconcileKubeconfig(ctx context.Context, rosaScope *scope.ROSAControlPlaneScope, ocmClient rosa.OCMClient, cluster *cmv1.Cluster) error {
	rosaScope.Debug("Reconciling ROSA kubeconfig for cluster", "cluster-name", rosaScope.RosaClusterName())

	clusterRef := client.ObjectKeyFromObject(rosaScope.Cluster)
	kubeconfigSecret, err := secret.GetFromNamespacedName(ctx, r.Client, clusterRef, secret.Kubeconfig)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return fmt.Errorf("failed to get kubeconfig secret: %w", err)
		}
	}

	// generate a new password for the cluster admin user, or retrieve an existing one.
	password, err := r.reconcileClusterAdminPassword(ctx, rosaScope)
	if err != nil {
		return fmt.Errorf("failed to reconcile cluster admin password secret: %w", err)
	}

	clusterName := rosaScope.RosaClusterName()
	userName := fmt.Sprintf("%s-capi-admin", clusterName)
	apiServerURL := cluster.API().URL()

	// create new user with admin privileges in the ROSA cluster if 'userName' doesn't already exist.
	err = rosa.CreateAdminUserIfNotExist(ocmClient, cluster.ID(), userName, password)
	if err != nil {
		return err
	}

	clientConfig := &restclient.Config{
		Host:     apiServerURL,
		Username: userName,
	}
	// request an acccess token using the credentials of the cluster admin user created earlier.
	// this token is used in the kubeconfig to authenticate with the API server.
	token, err := rosa.RequestToken(ctx, apiServerURL, userName, password, clientConfig)
	if err != nil {
		return fmt.Errorf("failed to request token: %w", err)
	}

	// create the kubeconfig spec.
	contextName := fmt.Sprintf("%s@%s", userName, clusterName)
	cfg := &api.Config{
		APIVersion: api.SchemeGroupVersion.Version,
		Clusters: map[string]*api.Cluster{
			clusterName: {
				Server: apiServerURL,
			},
		},
		Contexts: map[string]*api.Context{
			contextName: {
				Cluster:  clusterName,
				AuthInfo: userName,
			},
		},
		CurrentContext: contextName,
		AuthInfos: map[string]*api.AuthInfo{
			userName: {
				Token: token.AccessToken,
			},
		},
	}
	out, err := clientcmd.Write(*cfg)
	if err != nil {
		return fmt.Errorf("failed to serialize config to yaml: %w", err)
	}

	if kubeconfigSecret != nil {
		// update existing kubeconfig secret.
		kubeconfigSecret.Data[secret.KubeconfigDataName] = out
		if err := r.Client.Update(ctx, kubeconfigSecret); err != nil {
			return fmt.Errorf("failed to update kubeconfig secret: %w", err)
		}
	} else {
		// create new kubeconfig secret.
		controllerOwnerRef := *metav1.NewControllerRef(rosaScope.ControlPlane, rosacontrolplanev1.GroupVersion.WithKind("ROSAControlPlane"))
		kubeconfigSecret = kubeconfig.GenerateSecretWithOwner(clusterRef, out, controllerOwnerRef)
		if err := r.Client.Create(ctx, kubeconfigSecret); err != nil {
			return fmt.Errorf("failed to create kubeconfig secret: %w", err)
		}
	}

	rosaScope.ControlPlane.Status.Initialized = true
	return nil
}

// reconcileClusterAdminPassword generates and store the password of the cluster admin user in a secret which is used to request a token for kubeconfig auth.
// Since it is not possible to retrieve a user's password through the ocm API once created,
// we have to store the password in a secret as it is needed later to refresh the token.
func (r *ROSAControlPlaneReconciler) reconcileClusterAdminPassword(ctx context.Context, rosaScope *scope.ROSAControlPlaneScope) (string, error) {
	passwordSecret := rosaScope.ClusterAdminPasswordSecret()
	err := r.Client.Get(ctx, client.ObjectKeyFromObject(passwordSecret), passwordSecret)
	if err == nil {
		password := string(passwordSecret.Data["value"])
		return password, nil
	} else if !apierrors.IsNotFound(err) {
		return "", fmt.Errorf("failed to get cluster admin password secret: %w", err)
	}
	password, err := idputils.GenerateRandomPassword()
	// Generate a new password and create the secret
	if err != nil {
		return "", err
	}

	passwordSecret.Data = map[string][]byte{
		"value": []byte(password),
	}
	if err := r.Client.Create(ctx, passwordSecret); err != nil {
		return "", err
	}

	return password, nil
}

func validateControlPlaneSpec(ocmClient rosa.OCMClient, rosaScope *scope.ROSAControlPlaneScope) (string, error) {
	version := rosaScope.ControlPlane.Spec.Version
	channelGroup := string(rosaScope.ControlPlane.Spec.ChannelGroup)
	valid, err := ocmClient.ValidateHypershiftVersion(version, channelGroup)
	if err != nil {
		return "", fmt.Errorf("error validating version in this channelGroup : %w", err)
	}
	if !valid {
		return fmt.Sprintf("this version %s is not supported in this channelGroup", version), nil
	}

	// TODO: add more input validations
	return "", nil
}

func buildOCMClusterSpec(controlPlaneSpec rosacontrolplanev1.RosaControlPlaneSpec, creator *rosaaws.Creator) (ocm.Spec, error) {
	billingAccount := controlPlaneSpec.BillingAccount
	if billingAccount == "" {
		billingAccount = creator.AccountID
	}

	ocmClusterSpec := ocm.Spec{
		DryRun:                    ptr.To(false),
		Name:                      controlPlaneSpec.RosaClusterName,
		DomainPrefix:              controlPlaneSpec.DomainPrefix,
		Region:                    controlPlaneSpec.Region,
		MultiAZ:                   true,
		Version:                   ocm.CreateVersionID(controlPlaneSpec.Version, string(controlPlaneSpec.ChannelGroup)),
		ChannelGroup:              string(controlPlaneSpec.ChannelGroup),
		DisableWorkloadMonitoring: ptr.To(true),
		DefaultIngress:            ocm.NewDefaultIngressSpec(), // n.b. this is a no-op when it's set to the default value
		ComputeMachineType:        controlPlaneSpec.DefaultMachinePoolSpec.InstanceType,
		AvailabilityZones:         controlPlaneSpec.AvailabilityZones,
		Tags:                      controlPlaneSpec.AdditionalTags,
		EtcdEncryption:            controlPlaneSpec.EtcdEncryptionKMSARN != "",
		EtcdEncryptionKMSArn:      controlPlaneSpec.EtcdEncryptionKMSARN,

		SubnetIds:        controlPlaneSpec.Subnets,
		IsSTS:            true,
		RoleARN:          controlPlaneSpec.InstallerRoleARN,
		SupportRoleARN:   controlPlaneSpec.SupportRoleARN,
		WorkerRoleARN:    controlPlaneSpec.WorkerRoleARN,
		OperatorIAMRoles: operatorIAMRoles(controlPlaneSpec.RolesRef),
		OidcConfigId:     controlPlaneSpec.OIDCID,
		Mode:             "auto",
		Hypershift: ocm.Hypershift{
			Enabled: true,
		},
		BillingAccount:               billingAccount,
		AWSCreator:                   creator,
		AuditLogRoleARN:              ptr.To(controlPlaneSpec.AuditLogRoleARN),
		ExternalAuthProvidersEnabled: controlPlaneSpec.EnableExternalAuthProviders,
	}

	if controlPlaneSpec.EndpointAccess == rosacontrolplanev1.Private {
		ocmClusterSpec.Private = ptr.To(true)
		ocmClusterSpec.PrivateLink = ptr.To(true)
	}

	if networkSpec := controlPlaneSpec.Network; networkSpec != nil {
		if networkSpec.MachineCIDR != "" {
			_, machineCIDR, err := net.ParseCIDR(networkSpec.MachineCIDR)
			if err != nil {
				return ocmClusterSpec, err
			}
			ocmClusterSpec.MachineCIDR = *machineCIDR
		}

		if networkSpec.PodCIDR != "" {
			_, podCIDR, err := net.ParseCIDR(networkSpec.PodCIDR)
			if err != nil {
				return ocmClusterSpec, err
			}
			ocmClusterSpec.PodCIDR = *podCIDR
		}

		if networkSpec.ServiceCIDR != "" {
			_, serviceCIDR, err := net.ParseCIDR(networkSpec.ServiceCIDR)
			if err != nil {
				return ocmClusterSpec, err
			}
			ocmClusterSpec.ServiceCIDR = *serviceCIDR
		}

		ocmClusterSpec.HostPrefix = networkSpec.HostPrefix
		ocmClusterSpec.NetworkType = networkSpec.NetworkType
	}

	if controlPlaneSpec.DefaultMachinePoolSpec.VolumeSize >= 75 {
		ocmClusterSpec.MachinePoolRootDisk = &ocm.Volume{Size: controlPlaneSpec.DefaultMachinePoolSpec.VolumeSize}
	}

	// Set cluster compute autoscaling replicas
	// In case autoscaling is not defined and multiple zones defined, set the compute nodes equal to the zones count.
	if computeAutoscaling := controlPlaneSpec.DefaultMachinePoolSpec.Autoscaling; computeAutoscaling != nil {
		ocmClusterSpec.Autoscaling = true
		ocmClusterSpec.MaxReplicas = computeAutoscaling.MaxReplicas
		ocmClusterSpec.MinReplicas = computeAutoscaling.MinReplicas
	} else if len(controlPlaneSpec.AvailabilityZones) > 1 {
		ocmClusterSpec.ComputeNodes = len(controlPlaneSpec.AvailabilityZones)
	}

	if controlPlaneSpec.ProvisionShardID != "" {
		ocmClusterSpec.CustomProperties = map[string]string{
			"provision_shard_id": controlPlaneSpec.ProvisionShardID,
		}
	}

	// Set the cluster registry config.
	if controlPlaneSpec.ClusterRegistryConfig != nil {
		if len(controlPlaneSpec.ClusterRegistryConfig.AdditionalTrustedCAs) > 0 {
			ocmClusterSpec.AdditionalTrustedCa = controlPlaneSpec.ClusterRegistryConfig.AdditionalTrustedCAs
		}

		if len(controlPlaneSpec.ClusterRegistryConfig.AllowedRegistriesForImport) > 0 {
			registers := make([]string, 0)
			for id := range controlPlaneSpec.ClusterRegistryConfig.AllowedRegistriesForImport {
				registers = append(registers, controlPlaneSpec.ClusterRegistryConfig.AllowedRegistriesForImport[id].DomainName+":"+
					strconv.FormatBool(controlPlaneSpec.ClusterRegistryConfig.AllowedRegistriesForImport[id].Insecure))
			}
			ocmClusterSpec.AllowedRegistriesForImport = strings.Join(registers, ",")
		}

		if controlPlaneSpec.ClusterRegistryConfig.RegistrySources != nil {
			ocmClusterSpec.BlockedRegistries = controlPlaneSpec.ClusterRegistryConfig.RegistrySources.BlockedRegistries
			ocmClusterSpec.AllowedRegistries = controlPlaneSpec.ClusterRegistryConfig.RegistrySources.AllowedRegistries
			ocmClusterSpec.InsecureRegistries = controlPlaneSpec.ClusterRegistryConfig.RegistrySources.InsecureRegistries
		}
	}

	return ocmClusterSpec, nil
}

func operatorIAMRoles(rolesRef rosacontrolplanev1.AWSRolesRef) []ocm.OperatorIAMRole {
	return []ocm.OperatorIAMRole{
		{
			Name:      "cloud-credentials",
			Namespace: "openshift-ingress-operator",
			RoleARN:   rolesRef.IngressARN,
		},
		{
			Name:      "installer-cloud-credentials",
			Namespace: "openshift-image-registry",
			RoleARN:   rolesRef.ImageRegistryARN,
		},
		{
			Name:      "ebs-cloud-credentials",
			Namespace: "openshift-cluster-csi-drivers",
			RoleARN:   rolesRef.StorageARN,
		},
		{
			Name:      "cloud-credentials",
			Namespace: "openshift-cloud-network-config-controller",
			RoleARN:   rolesRef.NetworkARN,
		},
		{
			Name:      "kube-controller-manager",
			Namespace: "kube-system",
			RoleARN:   rolesRef.KubeCloudControllerARN,
		},
		{
			Name:      "kms-provider",
			Namespace: "kube-system",
			RoleARN:   rolesRef.KMSProviderARN,
		},
		{
			Name:      "control-plane-operator",
			Namespace: "kube-system",
			RoleARN:   rolesRef.ControlPlaneOperatorARN,
		},
		{
			Name:      "capa-controller-manager",
			Namespace: "kube-system",
			RoleARN:   rolesRef.NodePoolManagementARN,
		},
	}
}

func (r *ROSAControlPlaneReconciler) rosaClusterToROSAControlPlane(log *logger.Logger) handler.MapFunc {
	return func(ctx context.Context, o client.Object) []ctrl.Request {
		rosaCluster, ok := o.(*expinfrav1.ROSACluster)
		if !ok {
			log.Error(fmt.Errorf("expected a ROSACluster but got a %T", o), "Expected ROSACluster")
			return nil
		}

		if !rosaCluster.ObjectMeta.DeletionTimestamp.IsZero() {
			log.Debug("ROSACluster has a deletion timestamp, skipping mapping")
			return nil
		}

		cluster, err := util.GetOwnerCluster(ctx, r.Client, rosaCluster.ObjectMeta)
		if err != nil {
			log.Error(err, "failed to get owning cluster")
			return nil
		}
		if cluster == nil {
			log.Debug("Owning cluster not set on ROSACluster, skipping mapping")
			return nil
		}

		controlPlaneRef := cluster.Spec.ControlPlaneRef
		if controlPlaneRef == nil || controlPlaneRef.Kind != rosaControlPlaneKind {
			log.Debug("ControlPlaneRef is nil or not ROSAControlPlane, skipping mapping")
			return nil
		}

		return []ctrl.Request{
			{
				NamespacedName: types.NamespacedName{
					Name:      controlPlaneRef.Name,
					Namespace: controlPlaneRef.Namespace,
				},
			},
		}
	}
}

func buildAPIEndpoint(cluster *cmv1.Cluster) (*clusterv1.APIEndpoint, error) {
	parsedURL, err := url.ParseRequestURI(cluster.API().URL())
	if err != nil {
		return nil, err
	}
	host, portStr, err := net.SplitHostPort(parsedURL.Host)
	if err != nil {
		return nil, err
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}

	return &clusterv1.APIEndpoint{
		Host: host,
		Port: int32(port), //#nosec G109 G115
	}, nil
}

// TODO: Remove this and update the aws-sdk lib to v2.
func convertStsV2(identity *sts.GetCallerIdentityOutput) *stsv2.GetCallerIdentityOutput {
	return &stsv2.GetCallerIdentityOutput{
		Account: identity.Account,
		Arn:     identity.Arn,
		UserId:  identity.UserId,
	}
}
