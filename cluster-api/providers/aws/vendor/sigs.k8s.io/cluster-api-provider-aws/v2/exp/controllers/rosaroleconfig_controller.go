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
	"k8s.io/apimachinery/pkg/runtime"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/rosa/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	stsiface "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/sts"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/rosa"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/predicates"
)

// ROSARoleConfigReconciler reconciles a ROSARoleConfig object.
type ROSARoleConfigReconciler struct {
	client.Client
	Recorder         record.EventRecorder
	WatchFilterValue string
	NewStsClient     func(cloud.ScopeUsage, cloud.Session, logger.Wrapper, runtime.Object) stsiface.STSClient
	NewOCMClient     func(ctx context.Context, scope rosa.OCMSecretsRetriever) (rosa.OCMClient, error)
	Runtime          *rosacli.Runtime
}

func (r *ROSARoleConfigReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := logger.FromContext(ctx)
	r.NewOCMClient = rosa.NewWrappedOCMClientWithoutControlPlane
	r.NewStsClient = scope.NewSTSClient

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
		conditions.SetSummary(scope.RosaRoleConfig, conditions.WithConditions(expinfrav1.RosaRoleConfigReadyCondition), conditions.WithStepCounter())
		if err := scope.PatchObject(); err != nil {
			reterr = errors.Join(reterr, err)
		}
	}()

	if err := r.setUpRuntime(ctx, scope); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to set up runtime: %w", err)
	}

	if !roleConfig.DeletionTimestamp.IsZero() {
		scope.Info("Deleting ROSARoleConfig.")
		conditions.MarkFalse(scope.RosaRoleConfig, expinfrav1.RosaRoleConfigReadyCondition, expinfrav1.RosaRoleConfigDeletionStarted, clusterv1.ConditionSeverityInfo, "Deletion of RosaRolesConfig started")
		err = r.reconcileDelete(scope)
		if err == nil {
			controllerutil.RemoveFinalizer(scope.RosaRoleConfig, expinfrav1.RosaRoleConfigFinalizer)
		}

		return ctrl.Result{}, err
	}

	if controllerutil.AddFinalizer(scope.RosaRoleConfig, expinfrav1.RosaRoleConfigFinalizer) {
		return ctrl.Result{}, err
	}

	if err := r.reconcileAccountRoles(scope); err != nil {
		conditions.MarkFalse(scope.RosaRoleConfig, expinfrav1.RosaRoleConfigReadyCondition, expinfrav1.RosaRoleConfigReconciliationFailedReason, clusterv1.ConditionSeverityError, "Account Roles failure: %v", err)
		return ctrl.Result{}, fmt.Errorf("account Roles: %w", err)
	}

	if err := r.reconcileOIDC(scope); err != nil {
		conditions.MarkFalse(scope.RosaRoleConfig, expinfrav1.RosaRoleConfigReadyCondition, expinfrav1.RosaRoleConfigReconciliationFailedReason, clusterv1.ConditionSeverityError, "OIDC Config/provider failure: %v", err)
		return ctrl.Result{}, fmt.Errorf("oicd Config: %w", err)
	}

	if err := r.reconcileOperatorRoles(scope); err != nil {
		conditions.MarkFalse(scope.RosaRoleConfig, expinfrav1.RosaRoleConfigReadyCondition, expinfrav1.RosaRoleConfigReconciliationFailedReason, clusterv1.ConditionSeverityError, "Operator Roles failure: %v", err)
		return ctrl.Result{}, fmt.Errorf("operator Roles: %w", err)
	}

	if r.rosaRolesConfigReady(scope.RosaRoleConfig) {
		conditions.Set(scope.RosaRoleConfig,
			&clusterv1.Condition{
				Type:     expinfrav1.RosaRoleConfigReadyCondition,
				Status:   corev1.ConditionTrue,
				Reason:   expinfrav1.RosaRoleConfigCreatedReason,
				Severity: clusterv1.ConditionSeverityInfo,
				Message:  "RosaRoleConfig is ready",
			})
	} else {
		conditions.Set(scope.RosaRoleConfig,
			&clusterv1.Condition{
				Type:     expinfrav1.RosaRoleConfigReadyCondition,
				Status:   corev1.ConditionFalse,
				Reason:   expinfrav1.RosaRoleConfigCreatedReason,
				Severity: clusterv1.ConditionSeverityInfo,
				Message:  "RosaRoleConfig not ready",
			})
	}

	return ctrl.Result{}, nil
}

func (r *ROSARoleConfigReconciler) reconcileDelete(scope *scope.RosaRoleConfigScope) error {
	if err := r.deleteOperatorRoles(scope); err != nil {
		conditions.MarkFalse(scope.RosaRoleConfig, expinfrav1.RosaRoleConfigReadyCondition, expinfrav1.RosaRoleConfigDeletionFailedReason, clusterv1.ConditionSeverityError, "Failed to delete operator roles: %v", err)
		return err
	}

	if err := r.deleteOIDC(scope); err != nil {
		conditions.MarkFalse(scope.RosaRoleConfig, expinfrav1.RosaRoleConfigReadyCondition, expinfrav1.RosaRoleConfigDeletionFailedReason, clusterv1.ConditionSeverityError, "Failed to delete OIDC provider: %v", err)
		return err
	}

	if err := r.deleteAccountRoles(scope); err != nil {
		conditions.MarkFalse(scope.RosaRoleConfig, expinfrav1.RosaRoleConfigReadyCondition, expinfrav1.RosaRoleConfigDeletionFailedReason, clusterv1.ConditionSeverityError, "Failed to delete account roles: %v", err)
		return err
	}

	return nil
}

func (r *ROSARoleConfigReconciler) reconcileOperatorRoles(scope *scope.RosaRoleConfigScope) error {
	operatorRoles, err := r.Runtime.AWSClient.ListOperatorRoles("", "", scope.RosaRoleConfig.Spec.OperatorRoleConfig.Prefix)
	if err != nil {
		return err
	}

	operatorRolesRef := v1beta2.AWSRolesRef{}
	for _, role := range operatorRoles[scope.RosaRoleConfig.Spec.OperatorRoleConfig.Prefix] {
		if strings.Contains(role.RoleName, expinfrav1.IngressOperatorARNSuffix) {
			operatorRolesRef.IngressARN = role.RoleARN
		} else if strings.Contains(role.RoleName, expinfrav1.ImageRegistryARNSuffix) {
			operatorRolesRef.ImageRegistryARN = role.RoleARN
		} else if strings.Contains(role.RoleName, expinfrav1.StorageARNSuffix) {
			operatorRolesRef.StorageARN = role.RoleARN
		} else if strings.Contains(role.RoleName, expinfrav1.NetworkARNSuffix) {
			operatorRolesRef.NetworkARN = role.RoleARN
		} else if strings.Contains(role.RoleName, expinfrav1.KubeCloudControllerARNSuffix) {
			operatorRolesRef.KubeCloudControllerARN = role.RoleARN
		} else if strings.Contains(role.RoleName, expinfrav1.NodePoolManagementARNSuffix) {
			operatorRolesRef.NodePoolManagementARN = role.RoleARN
		} else if strings.Contains(role.RoleName, expinfrav1.ControlPlaneOperatorARNSuffix) {
			operatorRolesRef.ControlPlaneOperatorARN = role.RoleARN
		} else if strings.Contains(role.RoleName, expinfrav1.KMSProviderARNSuffix) {
			operatorRolesRef.KMSProviderARN = role.RoleARN
		}
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

	policies, err := r.Runtime.OCMClient.GetPolicies("OperatorRole")
	if err != nil {
		return err
	}

	// create operator roles
	config := scope.RosaRoleConfig.Spec.OperatorRoleConfig
	return operatorroles.CreateOperatorRoles(r.Runtime, rosa.GetOCMClientEnv(r.Runtime.OCMClient), config.PermissionsBoundaryARN,
		interactive.ModeAuto, policies, "", config.SharedVPCConfig.IsSharedVPC(), config.Prefix, true, installerRoleArn,
		true, oidcConfigID, config.SharedVPCConfig.RouteRoleARN, ocm.DefaultChannelGroup,
		config.SharedVPCConfig.VPCEndpointRoleARN)
}

func (r *ROSARoleConfigReconciler) reconcileOIDC(scope *scope.RosaRoleConfigScope) error {
	oidcID := ""
	switch scope.RosaRoleConfig.Spec.OidcProviderType {
	case expinfrav1.Managed:
		// Create oidcConfig if not exist
		if scope.RosaRoleConfig.Status.OIDCID == "" {
			oidcID, createErr := oidcconfig.CreateOIDCConfig(r.Runtime, true, "", "")
			if createErr != nil {
				return fmt.Errorf("failed to Create OIDC config: %w", createErr)
			}
			scope.RosaRoleConfig.Status.OIDCID = oidcID
		}
		oidcID = scope.RosaRoleConfig.Status.OIDCID
	case expinfrav1.Unmanaged:
		oidcID = scope.RosaRoleConfig.Spec.OperatorRoleConfig.OIDCID
	}

	// Check if oidc Config exist
	oidcConfig, err := r.Runtime.OCMClient.GetOidcConfig(oidcID)
	if err != nil || oidcConfig == nil {
		return fmt.Errorf("failed to get OIDC config: %w", err)
	}

	scope.RosaRoleConfig.Status.OIDCID = oidcConfig.ID()

	// check oidc providers
	providers, err := r.Runtime.AWSClient.ListOidcProviders("", oidcConfig)
	if err != nil {
		return err
	}

	// set oidc Provider Arn
	for _, provider := range providers {
		if strings.Contains(provider.Arn, oidcID) {
			scope.RosaRoleConfig.Status.OIDCProviderARN = provider.Arn
			return nil
		}
	}

	// create oidc provider if not exist.
	if scope.RosaRoleConfig.Status.OIDCProviderARN == "" {
		if err := oidcprovider.CreateOIDCProvider(r.Runtime, oidcID, "", true); err != nil {
			return err
		}
		providerArn, err := r.Runtime.AWSClient.GetOpenIDConnectProviderByOidcEndpointUrl(oidcConfig.IssuerUrl())
		if err != nil {
			return err
		}
		scope.RosaRoleConfig.Status.OIDCProviderARN = providerArn
	}

	return nil
}

func (r *ROSARoleConfigReconciler) reconcileAccountRoles(scope *scope.RosaRoleConfigScope) error {
	accountRoles, err := r.Runtime.AWSClient.ListAccountRoles(scope.RosaRoleConfig.Spec.AccountRoleConfig.Version)
	if err != nil {
		// ListAccountRoles return error if roles does not exist. return for any other error
		if !strings.Contains(err.Error(), "no account roles found") {
			return err
		}
	}

	accountRolesRef := expinfrav1.AccountRolesRef{}
	for _, role := range accountRoles {
		if role.RoleName == fmt.Sprintf("%s%s", scope.RosaRoleConfig.Spec.AccountRoleConfig.Prefix, expinfrav1.HCPROSAInstallerRole) {
			accountRolesRef.InstallerRoleARN = role.RoleARN
		} else if role.RoleName == fmt.Sprintf("%s%s", scope.RosaRoleConfig.Spec.AccountRoleConfig.Prefix, expinfrav1.HCPROSASupportRole) {
			accountRolesRef.SupportRoleARN = role.RoleARN
		} else if role.RoleName == fmt.Sprintf("%s%s", scope.RosaRoleConfig.Spec.AccountRoleConfig.Prefix, expinfrav1.HCPROSAWorkerRole) {
			accountRolesRef.WorkerRoleARN = role.RoleARN
		}
	}

	// Set account role ref if ready
	if r.accountRolesReady(accountRolesRef) {
		scope.RosaRoleConfig.Status.AccountRolesRef = accountRolesRef
		return nil
	}

	policies, err := r.Runtime.OCMClient.GetPolicies("AccountRole")
	if err != nil {
		return err
	}

	return accountroles.CreateHCPRoles(r.Runtime, scope.RosaRoleConfig.Spec.AccountRoleConfig.Prefix, true, scope.RosaRoleConfig.Spec.AccountRoleConfig.PermissionsBoundaryARN,
		rosa.GetOCMClientEnv(r.Runtime.OCMClient), policies, scope.RosaRoleConfig.Spec.AccountRoleConfig.Version, scope.RosaRoleConfig.Spec.AccountRoleConfig.Path,
		scope.RosaRoleConfig.Spec.AccountRoleConfig.SharedVPCConfig.IsSharedVPC(), scope.RosaRoleConfig.Spec.AccountRoleConfig.SharedVPCConfig.RouteRoleARN,
		scope.RosaRoleConfig.Spec.AccountRoleConfig.SharedVPCConfig.VPCEndpointRoleARN)
}

func (r *ROSARoleConfigReconciler) deleteAccountRoles(scope *scope.RosaRoleConfigScope) error {
	// list all account role names.
	prefix := scope.RosaRoleConfig.Spec.AccountRoleConfig.Prefix
	hasSharedVpcPolicies := scope.RosaRoleConfig.Spec.AccountRoleConfig.SharedVPCConfig.IsSharedVPC()
	roleNames := []string{fmt.Sprintf("%s%s", prefix, expinfrav1.HCPROSAInstallerRole),
		fmt.Sprintf("%s%s", prefix, expinfrav1.HCPROSASupportRole),
		fmt.Sprintf("%s%s", prefix, expinfrav1.HCPROSAWorkerRole)}

	var errs []error
	for _, roleName := range roleNames {
		if err := r.Runtime.AWSClient.DeleteAccountRole(roleName, prefix, true, hasSharedVpcPolicies); err != nil {
			errs = append(errs, err)
		}
	}

	return kerrors.NewAggregate(errs)
}

func (r *ROSARoleConfigReconciler) deleteOIDC(scope *scope.RosaRoleConfigScope) error {
	// Delete only managed oidc
	if scope.RosaRoleConfig.Spec.OidcProviderType == expinfrav1.Managed && scope.RosaRoleConfig.Status.OIDCID != "" {
		oidcConfig, err := r.Runtime.OCMClient.GetOidcConfig(scope.RosaRoleConfig.Status.OIDCID)
		if err != nil {
			return err
		}

		oidcEndpointURL := oidcConfig.IssuerUrl()
		if usedOidcProvider, err := r.Runtime.OCMClient.HasAClusterUsingOidcProvider(oidcEndpointURL, r.Runtime.Creator.AccountID); err != nil {
			return err
		} else if usedOidcProvider {
			return fmt.Errorf("clusters using OIDC provider '%s', cannot be deleted", oidcEndpointURL)
		}

		if err = r.Runtime.AWSClient.DeleteOpenIDConnectProvider(scope.RosaRoleConfig.Status.OIDCProviderARN); err != nil {
			return err
		}

		return r.Runtime.OCMClient.DeleteOidcConfig(oidcConfig.ID())
	}

	return nil
}

func (r *ROSARoleConfigReconciler) deleteOperatorRoles(scope *scope.RosaRoleConfigScope) error {
	prefix := scope.RosaRoleConfig.Spec.OperatorRoleConfig.Prefix
	if usedOperatorRoles, err := r.Runtime.OCMClient.HasAClusterUsingOperatorRolesPrefix(prefix); err != nil {
		return err
	} else if usedOperatorRoles {
		return fmt.Errorf("operator Roles with Prefix '%s' are in use cannot be deleted", prefix)
	}

	// list all operator role names.
	roleNames := []string{fmt.Sprintf("%s%s", prefix, expinfrav1.ControlPlaneOperatorARNSuffix),
		fmt.Sprintf("%s%s", prefix, expinfrav1.ImageRegistryARNSuffix),
		fmt.Sprintf("%s%s", prefix, expinfrav1.IngressOperatorARNSuffix),
		fmt.Sprintf("%s%s", prefix, expinfrav1.KMSProviderARNSuffix),
		fmt.Sprintf("%s%s", prefix, expinfrav1.KubeCloudControllerARNSuffix),
		fmt.Sprintf("%s%s", prefix, expinfrav1.NetworkARNSuffix),
		fmt.Sprintf("%s%s", prefix, expinfrav1.NodePoolManagementARNSuffix),
		fmt.Sprintf("%s%s", prefix, expinfrav1.StorageARNSuffix)}

	allSharedVpcPoliciesNotDeleted := make(map[string]bool)
	var errs []error
	for _, roleName := range roleNames {
		policiesNotDeleted, err := r.Runtime.AWSClient.DeleteOperatorRole(roleName, true, true)
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

// setUpRuntime sets up the ROSA runtime if it doesn't exist.
func (r *ROSARoleConfigReconciler) setUpRuntime(ctx context.Context, scope *scope.RosaRoleConfigScope) error {
	if r.Runtime != nil {
		return nil
	}

	// Create OCM client
	ocm, err := r.NewOCMClient(ctx, scope)
	if err != nil {
		return fmt.Errorf("failed to create OCM client: %w", err)
	}

	ocmClient, err := rosa.ConvertToRosaOcmClient(ocm)
	if err != nil || ocmClient == nil {
		return fmt.Errorf("failed to create OCM client: %w", err)
	}

	r.Runtime = rosacli.NewRuntime()
	r.Runtime.OCMClient = ocmClient
	r.Runtime.Reporter = reporter.CreateReporter() // &rosa.Reporter{}
	r.Runtime.Logger = rosalogging.NewLogger()

	r.Runtime.AWSClient, err = aws.NewClient().Logger(r.Runtime.Logger).Build()
	if err != nil {
		return fmt.Errorf("failed to create aws client: %w", err)
	}

	r.Runtime.Creator, err = r.Runtime.AWSClient.GetCreator()
	if err != nil {
		return fmt.Errorf("failed to get creator: %w", err)
	}

	return nil
}
