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

	pkgerrors "github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/klog/v2"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	ipamv1 "sigs.k8s.io/cluster-api/api/ipam/v1beta2"
	clusterutilv1 "sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/conditions"
	deprecatedv1beta1conditions "sigs.k8s.io/cluster-api/util/conditions/deprecated/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/api/govmomi/v1beta2"
	capvcontext "sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/util"
)

// +kubebuilder:rbac:groups=ipam.cluster.x-k8s.io,resources=ipaddressclaims,verbs=get;create;patch;watch;list;update
// +kubebuilder:rbac:groups=ipam.cluster.x-k8s.io,resources=ipaddresses,verbs=get;list;watch

// reconcileIPAddressClaims ensures that VSphereVMs that are configured with .spec.network.devices.addressFromPools
// have corresponding IPAddressClaims.
func (r vmReconciler) reconcileIPAddressClaims(ctx context.Context, vmCtx *capvcontext.VMContext) error {
	totalClaims, claimsCreated := 0, 0
	claimsFulfilled := 0
	log := ctrl.LoggerFrom(ctx)

	var (
		deprecatedV1beta1Claims []deprecatedv1beta1conditions.Getter
		claims                  []conditions.Getter
		errList                 []error
	)

	for devIdx, device := range vmCtx.VSphereVM.Spec.Network.Devices {
		for poolRefIdx, poolRef := range device.AddressesFromPools {
			totalClaims++
			ipAddrClaimName := util.IPAddressClaimName(vmCtx.VSphereVM.Name, devIdx, poolRefIdx)
			ipAddrClaim := &ipamv1.IPAddressClaim{}
			ipAddrClaimKey := client.ObjectKey{
				Namespace: vmCtx.VSphereVM.Namespace,
				Name:      ipAddrClaimName,
			}

			// Note: We have to use := here to create a new variable and not overwrite log & ctx outside the for loop.
			log := log.WithValues("IPAddressClaim", klog.KRef(ipAddrClaimKey.Namespace, ipAddrClaimKey.Name))
			ctx := ctrl.LoggerInto(ctx, log)

			err := vmCtx.Client.Get(ctx, ipAddrClaimKey, ipAddrClaim)
			if err != nil && !apierrors.IsNotFound(err) {
				return pkgerrors.Wrapf(err, "failed to get IPAddressClaim %s", klog.KRef(ipAddrClaimKey.Namespace, ipAddrClaimKey.Name))
			}
			ipAddrClaim, created, err := createOrPatchIPAddressClaim(ctx, vmCtx, ipAddrClaimName, poolRef)
			if err != nil {
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
				deprecatedV1beta1Claims = append(deprecatedV1beta1Claims, ipAddrClaim)
				claims = append(claims, ipAddrClaim)
			}
		}
	}

	if len(errList) > 0 {
		aggregatedErr := kerrors.NewAggregate(errList)
		deprecatedv1beta1conditions.MarkFalse(vmCtx.VSphereVM,
			infrav1.IPAddressClaimedV1Beta1Condition,
			infrav1.IPAddressClaimNotFoundV1Beta1Reason,
			clusterv1.ConditionSeverityError,
			"%v",
			aggregatedErr)
		conditions.Set(vmCtx.VSphereVM, metav1.Condition{
			Type:    infrav1.VSphereVMIPAddressClaimsFulfilledCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.VSphereVMIPAddressClaimsNotFulfilledReason,
			Message: aggregatedErr.Error(),
		})
		return aggregatedErr
	}

	// Calculating the VSphereVMIPAddressClaimsFulfilledCondition from the Ready Condition of the individual IPAddressClaims.
	// This will not work if the IPAM provider does not set the Ready condition on the IPAddressClaim.
	// To correctly calculate the status of the condition, we would want all the IPAddressClaim objects
	// to report the Ready Condition.
	if len(claims) == totalClaims {
		deprecatedv1beta1conditions.SetAggregate(vmCtx.VSphereVM,
			infrav1.IPAddressClaimedV1Beta1Condition,
			deprecatedV1beta1Claims,
			deprecatedv1beta1conditions.AddSourceRef(),
			deprecatedv1beta1conditions.WithStepCounter())

		if len(claims) > 0 {
			if err := conditions.SetAggregateCondition(claims, vmCtx.VSphereVM, clusterv1.ReadyCondition, conditions.TargetConditionType(infrav1.VSphereVMIPAddressClaimsFulfilledCondition)); err != nil {
				return pkgerrors.Wrap(err, "failed to aggregate Ready condition from IPAddressClaims")
			}
		} else {
			conditions.Set(vmCtx.VSphereVM, metav1.Condition{
				Type:   infrav1.VSphereVMIPAddressClaimsFulfilledCondition,
				Status: metav1.ConditionTrue,
				Reason: infrav1.VSphereVMIPAddressClaimsNotFulfilledReason,
			})
		}
		return nil
	}

	// Fallback logic to calculate the state of the IPAddressClaimed condition
	switch {
	case totalClaims == claimsFulfilled:
		deprecatedv1beta1conditions.MarkTrue(vmCtx.VSphereVM, infrav1.IPAddressClaimedV1Beta1Condition)
		conditions.Set(vmCtx.VSphereVM, metav1.Condition{
			Type:   infrav1.VSphereVMIPAddressClaimsFulfilledCondition,
			Status: metav1.ConditionTrue,
			Reason: infrav1.VSphereVMIPAddressClaimsFulfilledReason,
		})
	case claimsFulfilled < totalClaims && claimsCreated > 0:
		deprecatedv1beta1conditions.MarkFalse(vmCtx.VSphereVM, infrav1.IPAddressClaimedV1Beta1Condition,
			infrav1.IPAddressClaimsBeingCreatedV1Beta1Reason, clusterv1.ConditionSeverityInfo,
			"%d/%d claims being created", claimsCreated, totalClaims)
		conditions.Set(vmCtx.VSphereVM, metav1.Condition{
			Type:    infrav1.VSphereVMIPAddressClaimsFulfilledCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.VSphereVMIPAddressClaimsBeingCreatedReason,
			Message: fmt.Sprintf("%d/%d claims being created", claimsCreated, totalClaims),
		})
	case claimsFulfilled < totalClaims && claimsCreated == 0:
		deprecatedv1beta1conditions.MarkFalse(vmCtx.VSphereVM, infrav1.IPAddressClaimedV1Beta1Condition,
			infrav1.WaitingForIPAddressV1Beta1Reason, clusterv1.ConditionSeverityInfo,
			"%d/%d claims being processed", totalClaims-claimsFulfilled, totalClaims)
		conditions.Set(vmCtx.VSphereVM, metav1.Condition{
			Type:    infrav1.VSphereVMIPAddressClaimsFulfilledCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.VSphereVMIPAddressClaimsWaitingForIPAddressReason,
			Message: fmt.Sprintf("%d/%d claims being processed", totalClaims-claimsFulfilled, totalClaims),
		})
	}
	return nil
}

// createOrPatchIPAddressClaim creates/patches an IPAddressClaim object for a device requesting an address
// from an externally managed IPPool. Ensures that the claim has a reference to the cluster of the VM to
// support pausing reconciliation.
// The responsibility of the IP address resolution is handled by an external IPAM provider.
func createOrPatchIPAddressClaim(ctx context.Context, vmCtx *capvcontext.VMContext, name string, poolRef infrav1.IPPoolReference) (*ipamv1.IPAddressClaim, bool, error) {
	claim := &ipamv1.IPAddressClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: vmCtx.VSphereVM.Namespace,
		},
	}
	mutateFn := func() (err error) {
		claim.SetOwnerReferences(clusterutilv1.EnsureOwnerRef(
			claim.OwnerReferences,
			metav1.OwnerReference{
				APIVersion: infrav1.GroupVersion.String(),
				Kind:       "VSphereVM",
				Name:       vmCtx.VSphereVM.Name,
				UID:        vmCtx.VSphereVM.UID,
			}))

		ctrlutil.AddFinalizer(claim, infrav1.IPAddressClaimFinalizer)

		if claim.Labels == nil {
			claim.Labels = make(map[string]string)
		}
		claim.Labels[clusterv1.ClusterNameLabel] = vmCtx.VSphereVM.Labels[clusterv1.ClusterNameLabel]

		claim.Spec.PoolRef.APIGroup = poolRef.APIGroup
		claim.Spec.PoolRef.Kind = poolRef.Kind
		claim.Spec.PoolRef.Name = poolRef.Name
		return nil
	}
	log := ctrl.LoggerFrom(ctx)

	result, err := ctrlutil.CreateOrPatch(ctx, vmCtx.Client, claim, mutateFn)
	if err != nil {
		return nil, false, pkgerrors.Wrap(err, "failed to CreateOrPatch IPAddressClaim")
	}
	switch result {
	case ctrlutil.OperationResultCreated:
		log.Info("Created IPAddressClaim")
		return claim, true, nil
	case ctrlutil.OperationResultUpdated:
		log.Info("Updated IPAddressClaim")
	case ctrlutil.OperationResultNone, ctrlutil.OperationResultUpdatedStatus, ctrlutil.OperationResultUpdatedStatusOnly:
		log.V(3).Info("No change required for IPAddressClaim", "operationResult", result)
	}
	return claim, false, nil
}

// deleteIPAddressClaims removes the finalizers from the IPAddressClaim objects
// thus freeing them up for garbage collection.
func (r vmReconciler) deleteIPAddressClaims(ctx context.Context, vmCtx *capvcontext.VMContext) error {
	log := ctrl.LoggerFrom(ctx)
	for devIdx, device := range vmCtx.VSphereVM.Spec.Network.Devices {
		for poolRefIdx := range device.AddressesFromPools {
			// check if claim exists
			ipAddrClaim := &ipamv1.IPAddressClaim{}
			ipAddrClaimName := util.IPAddressClaimName(vmCtx.VSphereVM.Name, devIdx, poolRefIdx)
			ipAddrClaimKey := client.ObjectKey{
				Namespace: vmCtx.VSphereVM.Namespace,
				Name:      ipAddrClaimName,
			}
			if err := vmCtx.Client.Get(ctx, ipAddrClaimKey, ipAddrClaim); err != nil {
				if apierrors.IsNotFound(err) {
					continue
				}
				return pkgerrors.Wrapf(err, "failed to get IPAddressClaim %q to remove the finalizer", ipAddrClaimName)
			}

			if ctrlutil.RemoveFinalizer(ipAddrClaim, infrav1.IPAddressClaimFinalizer) {
				log.Info(fmt.Sprintf("Removing finalizer %s", infrav1.IPAddressClaimFinalizer), "IPAddressClaim", klog.KObj(ipAddrClaim))
				if err := vmCtx.Client.Update(ctx, ipAddrClaim); err != nil {
					return pkgerrors.Wrapf(err, "failed to update IPAddressClaim %s", klog.KObj(ipAddrClaim))
				}
			}
		}
	}
	return nil
}
