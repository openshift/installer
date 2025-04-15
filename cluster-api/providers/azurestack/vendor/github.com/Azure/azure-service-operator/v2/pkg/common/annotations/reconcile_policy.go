// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package annotations

// ReconcilePolicy describes the reconcile policy for the resource in question.
// A reconcile policy describes what action (if any) the operator is allowed to take when
// reconciling the resource.
// If no reconcile policy is specified, the default is "run"
const ReconcilePolicy = "serviceoperator.azure.com/reconcile-policy"

type ReconcilePolicyValue string

const (
	// ReconcilePolicyManage instructs the operator to manage the resource in question.
	// This includes issuing PUTs to update it and DELETE's to delete it from Azure if deleted in Kubernetes.
	// This is the default policy when no policy is specified.
	ReconcilePolicyManage = ReconcilePolicyValue("manage")

	// ReconcilePolicySkip instructs the operator to skip all reconciliation actions. This includes creating
	// the resource.
	ReconcilePolicySkip = ReconcilePolicyValue("skip")

	// ReconcilePolicyDetachOnDelete instructs the operator to skip deletion of resources in Azure. This allows
	// deletion of the resource in Kubernetes to go through but does not delete the underlying Azure resource.
	ReconcilePolicyDetachOnDelete = ReconcilePolicyValue("detach-on-delete")
)

// AllowsDelete determines if the policy allows deletion of the backing Azure resource
func (r ReconcilePolicyValue) AllowsDelete() bool {
	return r == ReconcilePolicyManage
}

// AllowsModify determines if the policy allows modification of the backing Azure resource
func (r ReconcilePolicyValue) AllowsModify() bool {
	return r == ReconcilePolicyManage || r == ReconcilePolicyDetachOnDelete
}
