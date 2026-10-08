/*
Copyright 2026 The Kubernetes Authors.

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
	"fmt"
	"sync"

	rosaaws "github.com/openshift/rosa/pkg/aws"
	rosalogging "github.com/openshift/rosa/pkg/logging"
	"github.com/openshift/rosa/pkg/reporter"
	rosacli "github.com/openshift/rosa/pkg/rosa"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	caparosa "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/rosa"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/rosa/ocmrole"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"
	"sigs.k8s.io/cluster-api/util/predicates"
)

// ROSAOCMRoleConfigReconciler reconciles an ROSAOCMRoleConfig object.
type ROSAOCMRoleConfigReconciler struct {
	client.Client
	Recorder         record.EventRecorder
	WatchFilterValue string
	NewOCMClient     func(ctx context.Context, scope caparosa.OCMSecretsRetriever) (caparosa.OCMClient, error)
	// runtimeFactory overrides runtime creation per reconciliation. Used in tests to inject mock clients.
	runtimeFactory func(ctx context.Context, scope *scope.ROSAOCMRoleConfigScope) (*rosacli.Runtime, error)

	// Per-organization linking mutex
	orgLinkMutexes sync.Map // map[orgID]*sync.Mutex
}

func (r *ROSAOCMRoleConfigReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := logger.FromContext(ctx)
	r.NewOCMClient = caparosa.NewWrappedOCMClientWithoutControlPlane

	return ctrl.NewControllerManagedBy(mgr).
		For(&expinfrav1.ROSAOCMRoleConfig{}).
		WithOptions(options).
		WithEventFilter(predicates.ResourceHasFilterLabel(mgr.GetScheme(), log.GetLogger(), r.WatchFilterValue)).
		Complete(r)
}

// +kubebuilder:rbac:groups=core,resources=events,verbs=get;list;watch;create;patch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=rosaocmroleconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=rosaocmroleconfigs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=rosaocmroleconfigs/finalizers,verbs=update

// Reconcile reconciles ROSAOCMRoleConfig.
func (r *ROSAOCMRoleConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (res ctrl.Result, reterr error) {
	log := logger.FromContext(ctx)

	ocmRoleConfig := &expinfrav1.ROSAOCMRoleConfig{}
	if err := r.Get(ctx, req.NamespacedName, ocmRoleConfig); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get ROSAOCMRoleConfig")
		return ctrl.Result{Requeue: true}, nil
	}

	log = log.WithValues("ocmRoleConfig", klog.KObj(ocmRoleConfig))
	scope, err := scope.NewROSAOCMRoleConfigScope(scope.ROSAOCMRoleConfigScopeParams{
		Client:            r.Client,
		ROSAOCMRoleConfig: ocmRoleConfig,
		ControllerName:    "rosaocmroleconfig",
		Logger:            log,
	})
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create rosaocmroleconfig scope: %w", err)
	}

	// Always close scope and set summary condition
	defer func() {
		v1beta1conditions.SetSummary(scope.ROSAOCMRoleConfig,
			v1beta1conditions.WithConditions(expinfrav1.ROSAOCMRoleConfigReadyCondition),
			v1beta1conditions.WithStepCounter())
		if err := scope.PatchObject(); err != nil {
			reterr = errors.Join(reterr, err)
		}
	}()

	rt, err := r.setUpRuntime(ctx, scope)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to set up runtime: %w", err)
	}

	if !scope.ROSAOCMRoleConfig.DeletionTimestamp.IsZero() {
		scope.Info("Deleting ROSAOCMRoleConfig.")
		v1beta1conditions.MarkFalse(scope.ROSAOCMRoleConfig, expinfrav1.ROSAOCMRoleConfigReadyCondition, expinfrav1.ROSAOCMRoleConfigDeletionStarted, clusterv1beta1.ConditionSeverityInfo, "Deletion of ROSAOCMRoleConfig started")
		err = r.reconcileDelete(ctx, scope, rt)
		if err != nil {
			v1beta1conditions.MarkFalse(scope.ROSAOCMRoleConfig, expinfrav1.ROSAOCMRoleConfigReadyCondition,
				expinfrav1.ROSAOCMRoleConfigDeletionFailedReason, clusterv1beta1.ConditionSeverityError,
				"Failed to delete ROSAOCMRoleConfig: %s", err.Error())
			return ctrl.Result{}, err
		}
		controllerutil.RemoveFinalizer(scope.ROSAOCMRoleConfig, expinfrav1.ROSAOCMRoleConfigFinalizer)
		return ctrl.Result{}, nil
	}

	if controllerutil.AddFinalizer(scope.ROSAOCMRoleConfig, expinfrav1.ROSAOCMRoleConfigFinalizer) {
		return ctrl.Result{}, nil
	}

	return r.reconcileOCMRole(ctx, scope, rt)
}

func (r *ROSAOCMRoleConfigReconciler) reconcileOCMRole(_ context.Context, scope *scope.ROSAOCMRoleConfigScope, rt *rosacli.Runtime) (ctrl.Result, error) {
	ocmRoleConfig := scope.ROSAOCMRoleConfig

	orgID, externalID, err := rt.OCMClient.GetCurrentOrganization()
	if err != nil {
		v1beta1conditions.MarkFalse(ocmRoleConfig, expinfrav1.ROSAOCMRoleConfigReadyCondition,
			expinfrav1.ROSAOCMRoleConfigReconciliationFailedReason, clusterv1beta1.ConditionSeverityError,
			"Failed to get organization: %s", err.Error())
		return ctrl.Result{}, fmt.Errorf("failed to get organization: %w", err)
	}

	roleName := rosaaws.GetOCMRoleName(ocmRoleConfig.Spec.RolePrefix, rosaaws.OCMRole, externalID)
	roleExists, existingRoleARN, err := rt.AWSClient.CheckRoleExists(roleName)
	if err != nil {
		v1beta1conditions.MarkFalse(ocmRoleConfig, expinfrav1.ROSAOCMRoleConfigReadyCondition,
			expinfrav1.ROSAOCMRoleConfigReconciliationFailedReason, clusterv1beta1.ConditionSeverityError,
			"Failed to check role existence: %s", err.Error())
		return ctrl.Result{}, fmt.Errorf("failed to check role existence: %w", err)
	}

	// If role exists, check link status
	if roleExists {
		ocmRoleConfig.Status.RoleARN = existingRoleARN
		ocmRoleConfig.Status.OrganizationID = orgID

		// Check if role is already linked to organization
		existsOnOCM, _, linkedARN, err := rt.OCMClient.CheckIfAWSAccountExists(orgID, rt.Creator.AccountID)
		if err != nil {
			v1beta1conditions.MarkFalse(ocmRoleConfig, expinfrav1.ROSAOCMRoleConfigReadyCondition,
				expinfrav1.ROSAOCMRoleConfigReconciliationFailedReason, clusterv1beta1.ConditionSeverityError,
				"Failed to check OCM link status: %s", err.Error())
			return ctrl.Result{}, fmt.Errorf("failed to check OCM link status: %w", err)
		}

		isLinked := existsOnOCM && linkedARN == existingRoleARN
		if isLinked {
			// Ensure the existing linked role's profile matches the requested profile
			existingProfile, err := r.detectRoleProfile(rt.AWSClient, roleName)
			if err != nil {
				v1beta1conditions.MarkFalse(ocmRoleConfig, expinfrav1.ROSAOCMRoleConfigReadyCondition,
					expinfrav1.ROSAOCMRoleConfigReconciliationFailedReason, clusterv1beta1.ConditionSeverityError,
					"Failed to detect role profile: %s", err.Error())
				return ctrl.Result{}, fmt.Errorf("failed to detect role profile: %w", err)
			}

			requestedProfile := string(ocmRoleConfig.Spec.Profile)
			if existingProfile != requestedProfile {
				v1beta1conditions.MarkFalse(ocmRoleConfig, expinfrav1.ROSAOCMRoleConfigReadyCondition,
					expinfrav1.ROSAOCMRoleConfigReconciliationFailedReason, clusterv1beta1.ConditionSeverityError,
					"Profile mismatch: role has %s profile but %s was requested", existingProfile, requestedProfile)
				return ctrl.Result{}, fmt.Errorf("profile mismatch: role %s has %s profile but %s was requested",
					roleName, existingProfile, requestedProfile)
			}

			v1beta1conditions.MarkTrue(ocmRoleConfig, expinfrav1.ROSAOCMRoleConfigReadyCondition)
			return ctrl.Result{}, nil
		}

		// Role exists but not linked - will proceed to GetOrCreateOCMRole for self-healing/linking
		scope.Logger.Info("OCM role exists but not linked, will complete creation and link to organization", "roleARN", existingRoleARN)
	}

	// Convert profile from CRD enum to extracted code enum
	extractedProfile, err := r.convertProfile(ocmRoleConfig.Spec.Profile)
	if err != nil {
		v1beta1conditions.MarkFalse(ocmRoleConfig, expinfrav1.ROSAOCMRoleConfigReadyCondition,
			expinfrav1.ROSAOCMRoleConfigReconciliationFailedReason, clusterv1beta1.ConditionSeverityError,
			"Invalid profile: %s", err.Error())
		return ctrl.Result{}, err
	}

	roleARN, orgID, _, err := ocmrole.GetOrCreateOCMRole(
		rt,
		&scope.Logger,
		ocmRoleConfig.Spec.RolePrefix,
		extractedProfile,
		ocmRoleConfig.Spec.PermissionsBoundaryARN,
		ocmRoleConfig.Spec.Path,
	)
	if err != nil {
		v1beta1conditions.MarkFalse(ocmRoleConfig, expinfrav1.ROSAOCMRoleConfigReadyCondition,
			expinfrav1.ROSAOCMRoleConfigReconciliationFailedReason, clusterv1beta1.ConditionSeverityError,
			"Failed to get or create OCM role: %s", err.Error())
		return ctrl.Result{}, fmt.Errorf("failed to get or create OCM role: %w", err)
	}

	ocmRoleConfig.Status.RoleARN = roleARN
	ocmRoleConfig.Status.OrganizationID = orgID

	// Acquire per-organization mutex to prevent concurrent LinkOrgToRole races.
	// Multiple ROSAOCMRoleConfigs with different AWS accounts can link to the same org,
	// causing read-modify-write races on the org label. Serialize linking per org.
	linkMu := r.getLinkMutex(orgID)
	linkMu.Lock()
	defer linkMu.Unlock()

	// LinkOrgToRole is idempotent - it checks if already linked and handles all cases:
	// - Not linked: links and returns (true, nil)
	// - Same role already linked: returns (false, nil)
	// - Different role already linked: returns (false, error)
	_, err = rt.OCMClient.LinkOrgToRole(orgID, roleARN)
	if err != nil {
		v1beta1conditions.MarkFalse(ocmRoleConfig, expinfrav1.ROSAOCMRoleConfigReadyCondition,
			expinfrav1.ROSAOCMRoleConfigReconciliationFailedReason, clusterv1beta1.ConditionSeverityError,
			"Failed to link OCM role: %s", err.Error())
		return ctrl.Result{}, fmt.Errorf("failed to link OCM role to organization: %w", err)
	}

	// Set ready condition
	v1beta1conditions.MarkTrue(ocmRoleConfig, expinfrav1.ROSAOCMRoleConfigReadyCondition)

	return ctrl.Result{}, nil
}

func (r *ROSAOCMRoleConfigReconciler) reconcileDelete(_ context.Context, scope *scope.ROSAOCMRoleConfigScope, rt *rosacli.Runtime) error {
	ocmRoleConfig := scope.ROSAOCMRoleConfig

	// Check deletion policy
	deletionPolicy := ocmRoleConfig.Spec.DeletionPolicy
	if deletionPolicy == "" {
		deletionPolicy = expinfrav1.ROSAOCMRoleDeletionPolicyDelete // default
	}

	if deletionPolicy == expinfrav1.ROSAOCMRoleDeletionPolicyRetain {
		scope.Info("DeletionPolicy is Retain, skipping OCM role deletion")
		return nil
	}

	orgID := ocmRoleConfig.Status.OrganizationID
	if orgID == "" {
		scope.Info("No organization ID in status, skipping cleanup")
		return nil
	}

	roleARN := ocmRoleConfig.Status.RoleARN
	if roleARN == "" {
		scope.Info("No role ARN in status, skipping cleanup")
		return nil
	}

	roleName, err := rosaaws.GetResourceIdFromARN(roleARN)
	if err != nil {
		return fmt.Errorf("failed to parse role ARN: %w", err)
	}

	if err := func() error {
		// Acquire per-organization mutex to prevent concurrent unlink/link races.
		// Multiple ROSAOCMRoleConfigs with different AWS accounts can link to the same org,
		// and we need to prevent a race where one CR is unlinking while another is linking.
		linkMu := r.getLinkMutex(orgID)
		linkMu.Lock()
		defer linkMu.Unlock()

		// CheckIfAWSAccountExists returns the actual linked ARN for this AWS account in this org
		existsOnOCM, _, linkedARN, err := rt.OCMClient.CheckIfAWSAccountExists(orgID, rt.Creator.AccountID)
		if err != nil {
			return fmt.Errorf("failed to check AWS account link status: %w", err)
		}

		if existsOnOCM && linkedARN == roleARN {
			if err := rt.OCMClient.UnlinkOCMRoleFromOrg(orgID, roleARN); err != nil {
				return fmt.Errorf("failed to unlink the role: %w", err)
			}
			scope.Info("Successfully unlinked role", "roleARN", roleARN, "orgID", orgID)
		} else if existsOnOCM && linkedARN != roleARN {
			scope.Info("Different role linked in OCM, skipping unlink", "expectedARN", roleARN, "linkedARN", linkedARN, "orgID", orgID)
		} else {
			scope.Info("No role linked in OCM for this AWS account, skipping unlink", "roleARN", roleARN, "orgID", orgID)
		}

		return nil
	}(); err != nil {
		return err
	}

	// Delete IAM role
	roleExists, _, err := rt.AWSClient.CheckRoleExists(roleName)
	if err != nil {
		return fmt.Errorf("failed to check role existence: %w", err)
	}

	if roleExists {
		err = rt.AWSClient.DeleteOCMRole(roleName, false)
		if err != nil {
			return fmt.Errorf("failed to delete IAM role: %w", err)
		}
		scope.Info("Successfully deleted OCM role", "roleName", roleName)
	} else {
		scope.Info("IAM role does not exist, skipping deletion", "roleName", roleName)
	}

	return nil
}

func (r *ROSAOCMRoleConfigReconciler) convertProfile(profile expinfrav1.ROSAOCMRoleProfile) (ocmrole.RoleProfile, error) {
	switch profile {
	case expinfrav1.ROSAOCMRoleProfileStandard:
		return ocmrole.ProfileStandard, nil
	case expinfrav1.ROSAOCMRoleProfileAdmin:
		return ocmrole.ProfileAdmin, nil
	case expinfrav1.ROSAOCMRoleProfileNoConsole:
		return ocmrole.ProfileNoConsole, nil
	default:
		return "", fmt.Errorf("unknown profile: %s", profile)
	}
}

func (r *ROSAOCMRoleConfigReconciler) getLinkMutex(orgID string) *sync.Mutex {
	mu, _ := r.orgLinkMutexes.LoadOrStore(orgID, &sync.Mutex{})
	return mu.(*sync.Mutex)
}

// setUpRuntime creates a ROSA runtime for the current reconciliation using the scope's AWS session.
// A new runtime is created per reconciliation so that the AWSClient always uses the correct
// identity (identityRef) for the ROSAOCMRoleConfig being reconciled.
func (r *ROSAOCMRoleConfigReconciler) setUpRuntime(ctx context.Context, scope *scope.ROSAOCMRoleConfigScope) (*rosacli.Runtime, error) {
	if r.runtimeFactory != nil {
		return r.runtimeFactory(ctx, scope)
	}

	// Create OCM client
	ocmClient, err := r.NewOCMClient(ctx, scope)
	if err != nil {
		return nil, fmt.Errorf("failed to create OCM client: %w", err)
	}

	ocm, err := caparosa.ConvertToRosaOcmClient(ocmClient)
	if err != nil || ocm == nil {
		return nil, fmt.Errorf("failed to convert OCM client: %w", err)
	}

	rt := rosacli.NewRuntime()
	rt.OCMClient = ocm
	rt.Reporter = reporter.CreateReporter()
	rt.Logger = rosalogging.NewLogger()

	session := scope.Session()
	awsClient, err := rosaaws.NewClient().
		Logger(rt.Logger).
		ExternalConfig(&session).
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS client: %w", err)
	}
	rt.AWSClient = awsClient

	creator, err := awsClient.GetCreator()
	if err != nil {
		return nil, fmt.Errorf("failed to get creator: %w", err)
	}
	rt.Creator = creator

	return rt, nil
}

// detectRoleProfile determines the role profile by checking role tags.
func (r *ROSAOCMRoleConfigReconciler) detectRoleProfile(awsClient rosaaws.Client, roleName string) (string, error) {
	isAdmin, err := awsClient.IsAdminRole(roleName)
	if err != nil {
		return "", fmt.Errorf("failed to check if role is admin: %w", err)
	}
	if isAdmin {
		return string(expinfrav1.ROSAOCMRoleProfileAdmin), nil
	}

	isNoConsole, err := ocmrole.IsNoConsoleRole(awsClient, roleName)
	if err != nil {
		return "", fmt.Errorf("failed to check if role is no-console: %w", err)
	}
	if isNoConsole {
		return string(expinfrav1.ROSAOCMRoleProfileNoConsole), nil
	}

	return string(expinfrav1.ROSAOCMRoleProfileStandard), nil
}
