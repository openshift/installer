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
	"fmt"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	ipamv1 "sigs.k8s.io/cluster-api/exp/ipam/api/v1alpha1"
	clusterutilv1 "sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/util"
)

// +kubebuilder:rbac:groups=ipam.cluster.x-k8s.io,resources=ipaddressclaims,verbs=get;create;patch;watch;list;update
// +kubebuilder:rbac:groups=ipam.cluster.x-k8s.io,resources=ipaddresses,verbs=get;list;watch

// reconcileIPAddressClaims ensures that VSphereVMs that are configured with .spec.network.devices.addressFromPools
// have corresponding IPAddressClaims.
func (r vmReconciler) reconcileIPAddressClaims(ctx *context.VMContext) error {
	totalClaims, claimsCreated := 0, 0
	claimsFulfilled := 0

	var (
		claims  []conditions.Getter
		errList []error
	)

	for devIdx, device := range ctx.VSphereVM.Spec.Network.Devices {
		for poolRefIdx, poolRef := range device.AddressesFromPools {
			totalClaims++
			ipAddrClaimName := util.IPAddressClaimName(ctx.VSphereVM.Name, devIdx, poolRefIdx)
			ipAddrClaim := &ipamv1.IPAddressClaim{}
			ipAddrClaimKey := client.ObjectKey{
				Namespace: ctx.VSphereVM.Namespace,
				Name:      ipAddrClaimName,
			}
			err := ctx.Client.Get(ctx, ipAddrClaimKey, ipAddrClaim)
			if err != nil && !apierrors.IsNotFound(err) {
				ctx.Logger.Error(err, "fetching IPAddressClaim failed", "name", ipAddrClaimName)
				return err
			}
			ipAddrClaim, created, err := createOrPatchIPAddressClaim(ctx, ipAddrClaimName, poolRef)
			if err != nil {
				ctx.Logger.Error(err, "createOrPatchIPAddressClaim failed", "name", ipAddrClaimName)
				errList = append(errList, err)
				continue
			}
			if created {
				claimsCreated++
			}
			if ipAddrClaim.Status.AddressRef.Name != "" {
				claimsFulfilled++
			}

			// Since this is eventually used to calculate the status of the
			// IPAddressClaimed condition for the VSphereVM object.
			if conditions.Has(ipAddrClaim, clusterv1.ReadyCondition) {
				claims = append(claims, ipAddrClaim)
			}
		}
	}

	if len(errList) > 0 {
		aggregatedErr := kerrors.NewAggregate(errList)
		conditions.MarkFalse(ctx.VSphereVM,
			infrav1.IPAddressClaimedCondition,
			infrav1.IPAddressClaimNotFoundReason,
			clusterv1.ConditionSeverityError,
			aggregatedErr.Error())
		return aggregatedErr
	}

	// Calculating the IPAddressClaimedCondition from the Ready Condition of the individual IPAddressClaims.
	// This will not work if the IPAM provider does not set the Ready condition on the IPAddressClaim.
	// To correctly calculate the status of the condition, we would want all the IPAddressClaim objects
	// to report the Ready Condition.
	if len(claims) == totalClaims {
		conditions.SetAggregate(ctx.VSphereVM,
			infrav1.IPAddressClaimedCondition,
			claims,
			conditions.AddSourceRef(),
			conditions.WithStepCounter())
		return nil
	}

	// Fallback logic to calculate the state of the IPAddressClaimed condition
	switch {
	case totalClaims == claimsFulfilled:
		conditions.MarkTrue(ctx.VSphereVM, infrav1.IPAddressClaimedCondition)
	case claimsFulfilled < totalClaims && claimsCreated > 0:
		conditions.MarkFalse(ctx.VSphereVM, infrav1.IPAddressClaimedCondition,
			infrav1.IPAddressClaimsBeingCreatedReason, clusterv1.ConditionSeverityInfo,
			"%d/%d claims being created", claimsCreated, totalClaims)
	case claimsFulfilled < totalClaims && claimsCreated == 0:
		conditions.MarkFalse(ctx.VSphereVM, infrav1.IPAddressClaimedCondition,
			infrav1.WaitingForIPAddressReason, clusterv1.ConditionSeverityInfo,
			"%d/%d claims being processed", totalClaims-claimsFulfilled, totalClaims)
	}
	return nil
}

// createOrPatchIPAddressClaim creates/patches an IPAddressClaim object for a device requesting an address
// from an externally managed IPPool. Ensures that the claim has a reference to the cluster of the VM to
// support pausing reconciliation.
// The responsibility of the IP address resolution is handled by an external IPAM provider.
func createOrPatchIPAddressClaim(ctx *context.VMContext, name string, poolRef corev1.TypedLocalObjectReference) (*ipamv1.IPAddressClaim, bool, error) {
	claim := &ipamv1.IPAddressClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ctx.VSphereVM.Namespace,
		},
	}
	mutateFn := func() (err error) {
		claim.SetOwnerReferences(clusterutilv1.EnsureOwnerRef(
			claim.OwnerReferences,
			metav1.OwnerReference{
				APIVersion: ctx.VSphereVM.APIVersion,
				Kind:       ctx.VSphereVM.Kind,
				Name:       ctx.VSphereVM.Name,
				UID:        ctx.VSphereVM.UID,
			}))

		ctrlutil.AddFinalizer(claim, infrav1.IPAddressClaimFinalizer)

		if claim.Labels == nil {
			claim.Labels = make(map[string]string)
		}
		claim.Labels[clusterv1.ClusterNameLabel] = ctx.VSphereVM.Labels[clusterv1.ClusterNameLabel]

		claim.Spec.PoolRef.APIGroup = poolRef.APIGroup
		claim.Spec.PoolRef.Kind = poolRef.Kind
		claim.Spec.PoolRef.Name = poolRef.Name
		return nil
	}

	result, err := ctrlutil.CreateOrPatch(ctx, ctx.Client, claim, mutateFn)
	if err != nil {
		ctx.Logger.Error(
			err,
			"failed to CreateOrPatch IPAddressClaim",
			"namespace",
			claim.Namespace,
			"name",
			claim.Name,
		)
		return nil, false, err
	}
	key := types.NamespacedName{
		Namespace: claim.Namespace,
		Name:      claim.Name,
	}
	switch result {
	case ctrlutil.OperationResultCreated:
		ctx.Logger.Info(
			"created claim",
			"claim",
			key,
		)
		return claim, true, nil
	case ctrlutil.OperationResultUpdated:
		ctx.Logger.Info(
			"updated claim",
			"claim",
			key,
		)
	case ctrlutil.OperationResultNone, ctrlutil.OperationResultUpdatedStatus, ctrlutil.OperationResultUpdatedStatusOnly:
		ctx.Logger.V(5).Info(
			"no change required for claim",
			"claim", key,
			"operation", result,
		)
	}
	return claim, false, nil
}

// deleteIPAddressClaims removes the finalizers from the IPAddressClaim objects
// thus freeing them up for garbage collection.
func (r vmReconciler) deleteIPAddressClaims(ctx *context.VMContext) error {
	for devIdx, device := range ctx.VSphereVM.Spec.Network.Devices {
		for poolRefIdx := range device.AddressesFromPools {
			// check if claim exists
			ipAddrClaim := &ipamv1.IPAddressClaim{}
			ipAddrClaimName := util.IPAddressClaimName(ctx.VSphereVM.Name, devIdx, poolRefIdx)
			ctx.Logger.Info("removing finalizer", "IPAddressClaim", ipAddrClaimName)
			ipAddrClaimKey := client.ObjectKey{
				Namespace: ctx.VSphereVM.Namespace,
				Name:      ipAddrClaimName,
			}
			if err := ctx.Client.Get(ctx, ipAddrClaimKey, ipAddrClaim); err != nil {
				if apierrors.IsNotFound(err) {
					continue
				}
				return errors.Wrapf(err, fmt.Sprintf("failed to find IPAddressClaim %q to remove the finalizer", ipAddrClaimName))
			}
			if ctrlutil.RemoveFinalizer(ipAddrClaim, infrav1.IPAddressClaimFinalizer) {
				if err := ctx.Client.Update(ctx, ipAddrClaim); err != nil {
					return errors.Wrapf(err, fmt.Sprintf("failed to update IPAddressClaim %q", ipAddrClaimName))
				}
			}
		}
	}
	return nil
}
