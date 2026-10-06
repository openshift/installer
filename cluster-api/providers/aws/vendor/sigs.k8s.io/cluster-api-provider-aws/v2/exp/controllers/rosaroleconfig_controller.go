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
	"fmt"
	"maps"
	"strings"

	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	accountroles "github.com/openshift/rosa/cmd/create/accountroles"
	oidcconfig "github.com/openshift/rosa/cmd/create/oidcconfig"
	oidcprovider "github.com/openshift/rosa/cmd/create/oidcprovider"
	operatorroles "github.com/openshift/rosa/cmd/create/operatorroles"
	"github.com/openshift/rosa/pkg/aws"
	interactive "github.com/openshift/rosa/pkg/interactive"
	rosalogging "github.com/openshift/rosa/pkg/logging"
	"github.com/openshift/rosa/pkg/ocm"
	"github.com/openshift/rosa/pkg/reporter"
	rosacli "github.com/openshift/rosa/pkg/rosa"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/rosa/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/rosa"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"
	"sigs.k8s.io/cluster-api/util/predicates"
)

// ROSARoleConfigReconciler reconciles a ROSARoleConfig object.
type ROSARoleConfigReconciler struct {
	client.Client
	Recorder         record.EventRecorder
	WatchFilterValue string
	NewOCMClient     func(ctx context.Context, scope rosa.OCMSecretsRetriever) (rosa.OCMClient, error)
	// runtimeFactory overrides runtime creation per reconciliation. Used in tests to inject mock clients.
	runtimeFactory func(ctx context.Context, scope *scope.RosaRoleConfigScope) (*rosacli.Runtime, error)
}

// roleNameLookup pairs an IAM role name with the string field that should receive its ARN.
type roleNameLookup struct {
	name string
	dest *string
}

func (r *ROSARoleConfigReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := logger.FromContext(ctx)
	r.NewOCMClient = rosa.NewWrappedOCMClientWithoutControlPlane

	return ctrl.NewControllerManagedBy(mgr).
		For(&expinfrav1.ROSARoleConfig{}).
		WithOptions(options).
		WithEventFilter(predicates.ResourceHasFilterLabel(mgr.GetScheme(), log.GetLogger(), r.WatchFilterValue)).
		Complete(r)
}

// +kubebuilder:rbac:groups=core,resources=events,verbs=get;list;watch;create;patch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=rosaroleconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=rosaroleconfigs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=rosaroleconfigs/finalizers,verbs=update

// Reconcile reconciles ROSARoleConfig.
func (r *ROSARoleConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (res ctrl.Result, reterr error) {
	log := logger.FromContext(ctx)

	roleConfig := &expinfrav1.ROSARoleConfig{}
	if err := r.Get(ctx, req.NamespacedName, roleConfig); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get ROSARoleConfig")
		return ctrl.Result{Requeue: true}, nil
	}

	log = log.WithValues("roleConfig", klog.KObj(roleConfig))
	scope, err := scope.NewRosaRoleConfigScope(scope.RosaRoleConfigScopeParams{
		Client:         r.Client,
		RosaRoleConfig: roleConfig,
		ControllerName: "rosaroleconfig",
		Logger:         log,
	})
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create rosaroleconfig scope: %w", err)
	}

	// Always close the scope and set summary condition
	defer func() {
		v1beta1conditions.SetSummary(scope.RosaRoleConfig, v1beta1conditions.WithConditions(expinfrav1.RosaRoleConfigReadyCondition), v1beta1conditions.WithStepCounter())
		if err := scope.PatchObject(); err != nil {
			reterr = errors.Join(reterr, err)
		}
	}()

	rt, err := r.setUpRuntime(ctx, scope)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to set up runtime: %w", err)
	}

	if !roleConfig.DeletionTimestamp.IsZero() {
		scope.Info("Deleting ROSARoleConfig.")
		v1beta1conditions.MarkFalse(scope.RosaRoleConfig, expinfrav1.RosaRoleConfigReadyCondition, expinfrav1.RosaRoleConfigDeletionStarted, clusterv1beta1.ConditionSeverityInfo, "Deletion of RosaRolesConfig started")
		err = r.reconcileDelete(scope, rt)
		if err == nil {
			controllerutil.RemoveFinalizer(scope.RosaRoleConfig, expinfrav1.RosaRoleConfigFinalizer)
		}

		return ctrl.Result{}, err
	}

	if controllerutil.AddFinalizer(scope.RosaRoleConfig, expinfrav1.RosaRoleConfigFinalizer) {
		return ctrl.Result{}, err
	}

	if err := r.reconcileAccountRoles(scope, rt); err != nil {
		v1beta1conditions.MarkFalse(scope.RosaRoleConfig, expinfrav1.RosaRoleConfigReadyCondition, expinfrav1.RosaRoleConfigReconciliationFailedReason, clusterv1beta1.ConditionSeverityError, "Account Roles failure: %v", err)
		return ctrl.Result{}, fmt.Errorf("account Roles: %w", err)
	}

	if externalID := scope.RosaRoleConfig.Spec.AccountRoleConfig.TrustPolicyExternalID; externalID != "" {
		prefix := scope.RosaRoleConfig.Spec.AccountRoleConfig.Prefix
		for _, suffix := range []string{expinfrav1.HCPROSAInstallerRole, expinfrav1.HCPROSASupportRole} {
			roleName := prefix + suffix
			if err := rosa.ApplyTrustPolicyExternalID(ctx, scope.IAMClient(), roleName, externalID); err != nil {
				v1beta1conditions.MarkFalse(scope.RosaRoleConfig, expinfrav1.RosaRoleConfigReadyCondition, expinfrav1.RosaRoleConfigReconciliationFailedReason, clusterv1beta1.ConditionSeverityError, "Trust policy external ID failure: %v", err)
				return ctrl.Result{}, fmt.Errorf("trust policy external ID: %w", err)
			}
		}
	}

	if err := r.reconcileOIDC(scope, rt); err != nil {
		v1beta1conditions.MarkFalse(scope.RosaRoleConfig, expinfrav1.RosaRoleConfigReadyCondition, expinfrav1.RosaRoleConfigReconciliationFailedReason, clusterv1beta1.ConditionSeverityError, "OIDC Config/provider failure: %v", err)
		return ctrl.Result{}, fmt.Errorf("oicd Config: %w", err)
	}

	if err := r.reconcileOperatorRoles(scope, rt); err != nil {
		v1beta1conditions.MarkFalse(scope.RosaRoleConfig, expinfrav1.RosaRoleConfigReadyCondition, expinfrav1.RosaRoleConfigReconciliationFailedReason, clusterv1beta1.ConditionSeverityError, "Operator Roles failure: %v", err)
		return ctrl.Result{}, fmt.Errorf("operator Roles: %w", err)
	}

	if r.rosaRolesConfigReady(scope.RosaRoleConfig) {
		v1beta1conditions.Set(scope.RosaRoleConfig,
			&clusterv1beta1.Condition{
				Type:     expinfrav1.RosaRoleConfigReadyCondition,
				Status:   corev1.ConditionTrue,
				Reason:   expinfrav1.RosaRoleConfigCreatedReason,
				Severity: clusterv1beta1.ConditionSeverityInfo,
				Message:  "RosaRoleConfig is ready",
			})
	} else {
		v1beta1conditions.Set(scope.RosaRoleConfig,
			&clusterv1beta1.Condition{
				Type:     expinfrav1.RosaRoleConfigReadyCondition,
				Status:   corev1.ConditionFalse,
				Reason:   expinfrav1.RosaRoleConfigCreatedReason,
				Severity: clusterv1beta1.ConditionSeverityInfo,
				Message:  "RosaRoleConfig not ready",
			})
	}

	return ctrl.Result{}, nil
}

func (r *ROSARoleConfigReconciler) reconcileDelete(scope *scope.RosaRoleConfigScope, rt *rosacli.Runtime) error {
	if err := r.deleteOperatorRoles(scope, rt); err != nil {
		v1beta1conditions.MarkFalse(scope.RosaRoleConfig, expinfrav1.RosaRoleConfigReadyCondition, expinfrav1.RosaRoleConfigDeletionFailedReason, clusterv1beta1.ConditionSeverityError, "Failed to delete operator roles: %v", err)
		return err
	}

	if err := r.deleteOIDC(scope, rt); err != nil {
		v1beta1conditions.MarkFalse(scope.RosaRoleConfig, expinfrav1.RosaRoleConfigReadyCondition, expinfrav1.RosaRoleConfigDeletionFailedReason, clusterv1beta1.ConditionSeverityError, "Failed to delete OIDC provider: %v", err)
		return err
	}

	if err := r.deleteAccountRoles(scope, rt); err != nil {
		v1beta1conditions.MarkFalse(scope.RosaRoleConfig, expinfrav1.RosaRoleConfigReadyCondition, expinfrav1.RosaRoleConfigDeletionFailedReason, clusterv1beta1.ConditionSeverityError, "Failed to delete account roles: %v", err)
		return err
	}

	return nil
}

func (r *ROSARoleConfigReconciler) reconcileOperatorRoles(scope *scope.RosaRoleConfigScope, rt *rosacli.Runtime) error {
	prefix := scope.RosaRoleConfig.Spec.OperatorRoleConfig.Prefix

	// Use targeted GetRoleByName lookups instead of ListOperatorRoles. ListOperatorRoles
	// calls ListRoleTags for every IAM role in the account, which causes throttling.
	operatorRolesRef, err := r.lookupOperatorRolesRef(rt, prefix)
	if err != nil {
		return err
	}

	if r.operatorRolesReady(operatorRolesRef) {
		scope.RosaRoleConfig.Status.OperatorRolesRef = operatorRolesRef
		return nil
	}

	installerRoleArn := scope.RosaRoleConfig.Status.AccountRolesRef.InstallerRoleARN
	if installerRoleArn == "" {
		scope.Logger.Info("installerRoleARN is empty, waiting for installer role to be created.")
		return nil
	}
	oidcConfigID := scope.RosaRoleConfig.Status.OIDCID
	if oidcConfigID == "" {
		scope.Logger.Info("oidcID is empty, waiting for oidcConfig to be created.")
		return nil
	}

	policies, err := rt.OCMClient.GetPolicies("OperatorRole")
	if err != nil {
		return err
	}

	config := scope.RosaRoleConfig.Spec.OperatorRoleConfig
	return operatorroles.CreateOperatorRoles(rt, rosa.GetOCMClientEnv(rt.OCMClient), config.PermissionsBoundaryARN,
		interactive.ModeAuto, policies, "", config.SharedVPCConfig.IsSharedVPC(), config.Prefix, true, installerRoleArn,
		true, oidcConfigID, config.SharedVPCConfig.RouteRoleARN, ocm.DefaultChannelGroup,
		config.SharedVPCConfig.VPCEndpointRoleARN)
}

// lookupOperatorRolesRef fetches each operator role ARN by its exact name using GetRoleByName.
// This avoids ListOperatorRoles, which calls ListRoleTags for every IAM role in the account.
func (r *ROSARoleConfigReconciler) lookupOperatorRolesRef(rt *rosacli.Runtime, prefix string) (v1beta2.AWSRolesRef, error) {
	ref := v1beta2.AWSRolesRef{}
	err := lookupRoleARNs(rt.AWSClient, []roleNameLookup{
		{fmt.Sprintf("%s%s", prefix, expinfrav1.IngressOperatorARNSuffix), &ref.IngressARN},
		{fmt.Sprintf("%s%s", prefix, expinfrav1.ImageRegistryARNSuffix), &ref.ImageRegistryARN},
		{fmt.Sprintf("%s%s", prefix, expinfrav1.StorageARNSuffix), &ref.StorageARN},
		{fmt.Sprintf("%s%s", prefix, expinfrav1.NetworkARNSuffix), &ref.NetworkARN},
		{fmt.Sprintf("%s%s", prefix, expinfrav1.KubeCloudControllerARNSuffix), &ref.KubeCloudControllerARN},
		{fmt.Sprintf("%s%s", prefix, expinfrav1.NodePoolManagementARNSuffix), &ref.NodePoolManagementARN},
		{fmt.Sprintf("%s%s", prefix, expinfrav1.ControlPlaneOperatorARNSuffix), &ref.ControlPlaneOperatorARN},
		{fmt.Sprintf("%s%s", prefix, expinfrav1.KMSProviderARNSuffix), &ref.KMSProviderARN},
	})
	return ref, err
}

func (r *ROSARoleConfigReconciler) reconcileOIDC(scope *scope.RosaRoleConfigScope, rt *rosacli.Runtime) error {
	oidcID := ""
	switch scope.RosaRoleConfig.Spec.OidcProviderType {
	case expinfrav1.Managed:
		// Create oidcConfig if not exist
		if scope.RosaRoleConfig.Status.OIDCID == "" {
			var createErr error
			oidcID, createErr = oidcconfig.CreateOIDCConfig(rt, true, "", "")
			if createErr != nil {
				return fmt.Errorf("failed to Create OIDC config: %w", createErr)
			}
			scope.RosaRoleConfig.Status.OIDCID = oidcID
			// Persist the OIDC config ID immediately so a subsequent reconcile does not
			// create a second config if anything after this point fails.
			if err := scope.PatchObject(); err != nil {
				return fmt.Errorf("failed to persist OIDC config ID: %w", err)
			}
		}
		oidcID = scope.RosaRoleConfig.Status.OIDCID
	case expinfrav1.Unmanaged:
		oidcID = scope.RosaRoleConfig.Spec.OperatorRoleConfig.OIDCID
	}

	// Check if oidc Config exist
	oidcConfig, err := rt.OCMClient.GetOidcConfig(oidcID)
	if err != nil || oidcConfig == nil {
		return fmt.Errorf("failed to get OIDC config: %w", err)
	}

	scope.RosaRoleConfig.Status.OIDCID = oidcConfig.ID()

	// Look up the provider by issuer URL. GetOpenIDConnectProviderByOidcEndpointUrl only
	// calls ListOpenIDConnectProviders once and matches by ARN string; unlike
	// ListOidcProviders it does not call ListOpenIDConnectProviderTags per provider,
	// which is the IAM API that triggers throttling under load.
	providerArn, err := rt.AWSClient.GetOpenIDConnectProviderByOidcEndpointUrl(oidcConfig.IssuerUrl())
	if err != nil {
		return err
	}
	if providerArn != "" {
		scope.RosaRoleConfig.Status.OIDCProviderARN = providerArn
		return nil
	}

	if err := oidcprovider.CreateOIDCProvider(rt, oidcID, "", true); err != nil {
		return err
	}
	providerArn, err = rt.AWSClient.GetOpenIDConnectProviderByOidcEndpointUrl(oidcConfig.IssuerUrl())
	if err != nil {
		return err
	}
	scope.RosaRoleConfig.Status.OIDCProviderARN = providerArn
	// Persist the provider ARN immediately so a subsequent reconcile does not
	// create a second provider if anything after this point fails.
	if err := scope.PatchObject(); err != nil {
		return fmt.Errorf("failed to persist OIDC provider ARN: %w", err)
	}

	return nil
}

func (r *ROSARoleConfigReconciler) reconcileAccountRoles(scope *scope.RosaRoleConfigScope, rt *rosacli.Runtime) error {
	prefix := scope.RosaRoleConfig.Spec.AccountRoleConfig.Prefix

	// Use targeted GetRoleByName lookups instead of ListAccountRoles. ListAccountRoles
	// calls ListRoleTags for every IAM role in the account, which causes throttling.
	accountRolesRef, err := r.lookupAccountRolesRef(rt, prefix)
	if err != nil {
		return err
	}

	if r.accountRolesReady(accountRolesRef) {
		scope.RosaRoleConfig.Status.AccountRolesRef = accountRolesRef
		return nil
	}

	policies, err := rt.OCMClient.GetPolicies("AccountRole")
	if err != nil {
		return err
	}

	return accountroles.CreateHCPRoles(rt, prefix, true, scope.RosaRoleConfig.Spec.AccountRoleConfig.PermissionsBoundaryARN,
		rosa.GetOCMClientEnv(rt.OCMClient), policies, scope.RosaRoleConfig.Spec.AccountRoleConfig.Version, scope.RosaRoleConfig.Spec.AccountRoleConfig.Path,
		scope.RosaRoleConfig.Spec.AccountRoleConfig.SharedVPCConfig.IsSharedVPC(), scope.RosaRoleConfig.Spec.AccountRoleConfig.SharedVPCConfig.RouteRoleARN,
		scope.RosaRoleConfig.Spec.AccountRoleConfig.SharedVPCConfig.VPCEndpointRoleARN)
}

// lookupAccountRolesRef fetches each account role ARN by its exact name using GetRoleByName.
// This avoids ListAccountRoles, which calls ListRoleTags for every IAM role in the account.
func (r *ROSARoleConfigReconciler) lookupAccountRolesRef(rt *rosacli.Runtime, prefix string) (expinfrav1.AccountRolesRef, error) {
	ref := expinfrav1.AccountRolesRef{}
	err := lookupRoleARNs(rt.AWSClient, []roleNameLookup{
		{fmt.Sprintf("%s%s", prefix, expinfrav1.HCPROSAInstallerRole), &ref.InstallerRoleARN},
		{fmt.Sprintf("%s%s", prefix, expinfrav1.HCPROSASupportRole), &ref.SupportRoleARN},
		{fmt.Sprintf("%s%s", prefix, expinfrav1.HCPROSAWorkerRole), &ref.WorkerRoleARN},
	})
	return ref, err
}

// lookupRoleARNs calls GetRoleByName for each entry and writes the ARN into dest.
// Roles that do not yet exist are silently skipped; any other error is returned immediately.
// This avoids List*Roles calls that enumerate every IAM role and call ListRoleTags per role.
func lookupRoleARNs(awsClient aws.Client, lookups []roleNameLookup) error {
	for _, lk := range lookups {
		role, err := awsClient.GetRoleByName(lk.name)
		if err != nil {
			var notFound *iamtypes.NoSuchEntityException
			if errors.As(err, &notFound) {
				continue
			}
			return err
		}
		if role.Arn != nil {
			*lk.dest = *role.Arn
		}
	}
	return nil
}

func (r *ROSARoleConfigReconciler) deleteAccountRoles(scope *scope.RosaRoleConfigScope, rt *rosacli.Runtime) error {
	// list all account role names.
	prefix := scope.RosaRoleConfig.Spec.AccountRoleConfig.Prefix
	hasSharedVpcPolicies := scope.RosaRoleConfig.Spec.AccountRoleConfig.SharedVPCConfig.IsSharedVPC()
	roleNames := []string{
		fmt.Sprintf("%s%s", prefix, expinfrav1.HCPROSAInstallerRole),
		fmt.Sprintf("%s%s", prefix, expinfrav1.HCPROSASupportRole),
		fmt.Sprintf("%s%s", prefix, expinfrav1.HCPROSAWorkerRole),
	}

	var errs []error
	for _, roleName := range roleNames {
		if err := rt.AWSClient.DeleteAccountRole(roleName, prefix, true, hasSharedVpcPolicies); err != nil {
			errs = append(errs, err)
		}
	}

	return kerrors.NewAggregate(errs)
}

func (r *ROSARoleConfigReconciler) deleteOIDC(scope *scope.RosaRoleConfigScope, rt *rosacli.Runtime) error {
	// Delete only managed oidc
	if scope.RosaRoleConfig.Spec.OidcProviderType == expinfrav1.Managed && scope.RosaRoleConfig.Status.OIDCID != "" {
		oidcConfig, err := rt.OCMClient.GetOidcConfig(scope.RosaRoleConfig.Status.OIDCID)
		if err != nil {
			return err
		}

		oidcEndpointURL := oidcConfig.IssuerUrl()
		if usedOidcProvider, err := rt.OCMClient.HasAClusterUsingOidcProvider(oidcEndpointURL, rt.Creator.AccountID); err != nil {
			return err
		} else if usedOidcProvider {
			return fmt.Errorf("clusters using OIDC provider '%s', cannot be deleted", oidcEndpointURL)
		}

		if err = rt.AWSClient.DeleteOpenIDConnectProvider(scope.RosaRoleConfig.Status.OIDCProviderARN); err != nil {
			return err
		}

		return rt.OCMClient.DeleteOidcConfig(oidcConfig.ID())
	}

	return nil
}

func (r *ROSARoleConfigReconciler) deleteOperatorRoles(scope *scope.RosaRoleConfigScope, rt *rosacli.Runtime) error {
	prefix := scope.RosaRoleConfig.Spec.OperatorRoleConfig.Prefix
	if usedOperatorRoles, err := rt.OCMClient.HasAClusterUsingOperatorRolesPrefix(prefix); err != nil {
		return err
	} else if usedOperatorRoles {
		return fmt.Errorf("operator Roles with Prefix '%s' are in use cannot be deleted", prefix)
	}

	// list all operator role names.
	roleNames := []string{
		fmt.Sprintf("%s%s", prefix, expinfrav1.ControlPlaneOperatorARNSuffix),
		fmt.Sprintf("%s%s", prefix, expinfrav1.ImageRegistryARNSuffix),
		fmt.Sprintf("%s%s", prefix, expinfrav1.IngressOperatorARNSuffix),
		fmt.Sprintf("%s%s", prefix, expinfrav1.KMSProviderARNSuffix),
		fmt.Sprintf("%s%s", prefix, expinfrav1.KubeCloudControllerARNSuffix),
		fmt.Sprintf("%s%s", prefix, expinfrav1.NetworkARNSuffix),
		fmt.Sprintf("%s%s", prefix, expinfrav1.NodePoolManagementARNSuffix),
		fmt.Sprintf("%s%s", prefix, expinfrav1.StorageARNSuffix),
	}

	allSharedVpcPoliciesNotDeleted := make(map[string]bool)
	var errs []error
	for _, roleName := range roleNames {
		policiesNotDeleted, err := rt.AWSClient.DeleteOperatorRole(roleName, true, true)
		if err != nil && (!strings.Contains(err.Error(), "does not exists") && !strings.Contains(err.Error(), "NoSuchEntity")) {
			errs = append(errs, err)
		}

		maps.Copy(allSharedVpcPoliciesNotDeleted, policiesNotDeleted)
	}

	for policyOutput, notDeleted := range allSharedVpcPoliciesNotDeleted {
		if notDeleted {
			scope.Logger.Info("unable to delete policy %s: Policy still attached to other resources", policyOutput)
		}
	}

	return kerrors.NewAggregate(errs)
}

func (r ROSARoleConfigReconciler) rosaRolesConfigReady(rosaRoleConfig *expinfrav1.ROSARoleConfig) bool {
	return rosaRoleConfig.Status.OIDCID != "" &&
		r.operatorRolesReady(rosaRoleConfig.Status.OperatorRolesRef) &&
		r.accountRolesReady(rosaRoleConfig.Status.AccountRolesRef)
}

func (r ROSARoleConfigReconciler) accountRolesReady(accountRolesRef expinfrav1.AccountRolesRef) bool {
	return accountRolesRef.InstallerRoleARN != "" &&
		accountRolesRef.SupportRoleARN != "" &&
		accountRolesRef.WorkerRoleARN != ""
}

func (r ROSARoleConfigReconciler) operatorRolesReady(operatorRolesRef v1beta2.AWSRolesRef) bool {
	return operatorRolesRef.ControlPlaneOperatorARN != "" &&
		operatorRolesRef.ImageRegistryARN != "" &&
		operatorRolesRef.IngressARN != "" &&
		operatorRolesRef.KMSProviderARN != "" &&
		operatorRolesRef.KubeCloudControllerARN != "" &&
		operatorRolesRef.NetworkARN != "" &&
		operatorRolesRef.NodePoolManagementARN != "" &&
		operatorRolesRef.StorageARN != ""
}

// setUpRuntime creates a ROSA runtime for the current reconciliation using the scope's AWS session.
// A new runtime is created per reconciliation so that the AWSClient always uses the correct
// identity (identityRef) for the ROSARoleConfig being reconciled.
func (r *ROSARoleConfigReconciler) setUpRuntime(ctx context.Context, scope *scope.RosaRoleConfigScope) (*rosacli.Runtime, error) {
	if r.runtimeFactory != nil {
		return r.runtimeFactory(ctx, scope)
	}
	// Create OCM client
	ocm, err := r.NewOCMClient(ctx, scope)
	if err != nil {
		return nil, fmt.Errorf("failed to create OCM client: %w", err)
	}

	ocmClient, err := rosa.ConvertToRosaOcmClient(ocm)
	if err != nil || ocmClient == nil {
		return nil, fmt.Errorf("failed to create OCM client: %w", err)
	}

	rt := rosacli.NewRuntime()
	rt.OCMClient = ocmClient
	rt.Reporter = reporter.CreateReporter()
	rt.Logger = rosalogging.NewLogger()

	session := scope.Session()
	awsClient, err := aws.NewClient().
		Logger(rt.Logger).
		ExternalConfig(&session).
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create aws client: %w", err)
	}
	rt.AWSClient = awsClient

	creator, err := awsClient.GetCreator()
	if err != nil {
		return nil, fmt.Errorf("failed to get creator: %w", err)
	}
	rt.Creator = creator

	return rt, nil
}
