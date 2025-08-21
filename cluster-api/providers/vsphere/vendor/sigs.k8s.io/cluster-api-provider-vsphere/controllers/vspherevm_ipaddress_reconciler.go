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

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/klog/v2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	ipamv1 "sigs.k8s.io/cluster-api/exp/ipam/api/v1beta1"
	clusterutilv1 "sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/conditions"
	v1beta2conditions "sigs.k8s.io/cluster-api/util/conditions/v1beta2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
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
		claims        []conditions.Getter
		v1beta2Claims []v1beta2conditions.Getter
		errList       []error
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
				return errors.Wrapf(err, "failed to get IPAddressClaim %s", klog.KRef(ipAddrClaimKey.Namespace, ipAddrClaimKey.Name))
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
				claims = append(claims, ipAddrClaim)
				v1beta2Claims = append(v1beta2Claims, ipAddrClaim)
			}
		}
	}

	if len(errList) > 0 {
		aggregatedErr := kerrors.NewAggregate(errList)
		conditions.MarkFalse(vmCtx.VSphereVM,
			infrav1.IPAddressClaimedCondition,
			infrav1.IPAddressClaimNotFoundReason,
			clusterv1.ConditionSeverityError,
			aggregatedErr.Error())
		v1beta2conditions.Set(vmCtx.VSphereVM, metav1.Condition{
			Type:    infrav1.VSphereVMIPAddressClaimsFulfilledV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.VSphereVMIPAddressClaimsNotFulfilledV1Beta2Reason,
			Message: aggregatedErr.Error(),
		})
		return aggregatedErr
	}

	// Calculating the IPAddressClaimedCondition from the Ready Condition of the individual IPAddressClaims.
	// This will not work if the IPAM provider does not set the Ready condition on the IPAddressClaim.
	// To correctly calculate the status of the condition, we would want all the IPAddressClaim objects
	// to report the Ready Condition.
	if len(claims) == totalClaims {
		conditions.SetAggregate(vmCtx.VSphereVM,
			infrav1.IPAddressClaimedCondition,
			claims,
			conditions.AddSourceRef(),
			conditions.WithStepCounter())

		if len(v1beta2Claims) > 0 {
			if err := v1beta2conditions.SetAggregateCondition(v1beta2Claims, vmCtx.VSphereVM, clusterv1.ReadyV1Beta2Condition, v1beta2conditions.TargetConditionType(infrav1.VSphereVMIPAddressClaimsFulfilledV1Beta2Condition)); err != nil {
				return errors.Wrap(err, "failed to aggregate Ready condition from IPAddressClaims")
			}
		} else {
			v1beta2conditions.Set(vmCtx.VSphereVM, metav1.Condition{
				Type:   infrav1.VSphereVMIPAddressClaimsFulfilledV1Beta2Condition,
				Status: metav1.ConditionTrue,
				Reason: infrav1.VSphereVMIPAddressClaimsNotFulfilledV1Beta2Reason,
			})
		}
		return nil
	}

	// Fallback logic to calculate the state of the IPAddressClaimed condition
	switch {
	case totalClaims == claimsFulfilled:
		conditions.MarkTrue(vmCtx.VSphereVM, infrav1.IPAddressClaimedCondition)
		v1beta2conditions.Set(vmCtx.VSphereVM, metav1.Condition{
			Type:   infrav1.VSphereVMIPAddressClaimsFulfilledV1Beta2Condition,
			Status: metav1.ConditionTrue,
			Reason: infrav1.VSphereVMIPAddressClaimsFulfilledV1Beta2Reason,
		})
	case claimsFulfilled < totalClaims && claimsCreated > 0:
		conditions.MarkFalse(vmCtx.VSphereVM, infrav1.IPAddressClaimedCondition,
			infrav1.IPAddressClaimsBeingCreatedReason, clusterv1.ConditionSeverityInfo,
			"%d/%d claims being created", claimsCreated, totalClaims)
		v1beta2conditions.Set(vmCtx.VSphereVM, metav1.Condition{
			Type:    infrav1.VSphereVMIPAddressClaimsFulfilledV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.VSphereVMIPAddressClaimsBeingCreatedV1Beta2Reason,
			Message: fmt.Sprintf("%d/%d claims being created", claimsCreated, totalClaims),
		})
	case claimsFulfilled < totalClaims && claimsCreated == 0:
		conditions.MarkFalse(vmCtx.VSphereVM, infrav1.IPAddressClaimedCondition,
			infrav1.WaitingForIPAddressReason, clusterv1.ConditionSeverityInfo,
			"%d/%d claims being processed", totalClaims-claimsFulfilled, totalClaims)
		v1beta2conditions.Set(vmCtx.VSphereVM, metav1.Condition{
			Type:    infrav1.VSphereVMIPAddressClaimsFulfilledV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.VSphereVMIPAddressClaimsWaitingForIPAddressV1Beta2Reason,
			Message: fmt.Sprintf("%d/%d claims being processed", totalClaims-claimsFulfilled, totalClaims),
		})
	}
	return nil
}

// createOrPatchIPAddressClaim creates/patches an IPAddressClaim object for a device requesting an address
// from an externally managed IPPool. Ensures that the claim has a reference to the cluster of the VM to
// support pausing reconciliation.
// The responsibility of the IP address resolution is handled by an external IPAM provider.
func createOrPatchIPAddressClaim(ctx context.Context, vmCtx *capvcontext.VMContext, name string, poolRef corev1.TypedLocalObjectReference) (*ipamv1.IPAddressClaim, bool, error) {
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
		return nil, false, errors.Wrap(err, "failed to CreateOrPatch IPAddressClaim")
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
				return errors.Wrapf(err, "failed to get IPAddressClaim %q to remove the finalizer", ipAddrClaimName)
			}

			if ctrlutil.RemoveFinalizer(ipAddrClaim, infrav1.IPAddressClaimFinalizer) {
				log.Info(fmt.Sprintf("Removing finalizer %s", infrav1.IPAddressClaimFinalizer), "IPAddressClaim", klog.KObj(ipAddrClaim))
				if err := vmCtx.Client.Update(ctx, ipAddrClaim); err != nil {
					return errors.Wrapf(err, "failed to update IPAddressClaim %s", klog.KObj(ipAddrClaim))
				}
			}
		}
	}
	return nil
}
