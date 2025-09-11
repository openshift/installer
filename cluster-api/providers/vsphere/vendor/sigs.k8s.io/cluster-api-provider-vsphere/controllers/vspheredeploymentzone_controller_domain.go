/*
Copyright 2021 The Kubernetes Authors.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	clusterutilv1 "sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/conditions"
	v1beta2conditions "sigs.k8s.io/cluster-api/util/conditions/v1beta2"
	ctrl "sigs.k8s.io/controller-runtime"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	capvcontext "sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services/govmomi/cluster"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services/govmomi/metadata"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/taggable"
)

func (r vsphereDeploymentZoneReconciler) reconcileFailureDomain(ctx context.Context, deploymentZoneCtx *capvcontext.VSphereDeploymentZoneContext, vsphereFailureDomain *infrav1.VSphereFailureDomain) error {
	// verify the failure domain for the region
	if err := r.reconcileInfraFailureDomain(ctx, deploymentZoneCtx, vsphereFailureDomain, vsphereFailureDomain.Spec.Region); err != nil {
		conditions.MarkFalse(deploymentZoneCtx.VSphereDeploymentZone, infrav1.VSphereFailureDomainValidatedCondition, infrav1.RegionMisconfiguredReason, clusterv1.ConditionSeverityError, err.Error())
		v1beta2conditions.Set(deploymentZoneCtx.VSphereDeploymentZone, metav1.Condition{
			Type:    infrav1.VSphereDeploymentZoneFailureDomainValidatedV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.VSphereDeploymentZoneFailureDomainRegionMisconfiguredV1Beta2Reason,
			Message: err.Error(),
		})
		return errors.Wrapf(err, "failed to reconcile failure domain: region is not configured correctly")
	}

	// verify the failure domain for the zone
	if err := r.reconcileInfraFailureDomain(ctx, deploymentZoneCtx, vsphereFailureDomain, vsphereFailureDomain.Spec.Zone); err != nil {
		conditions.MarkFalse(deploymentZoneCtx.VSphereDeploymentZone, infrav1.VSphereFailureDomainValidatedCondition, infrav1.ZoneMisconfiguredReason, clusterv1.ConditionSeverityError, err.Error())
		v1beta2conditions.Set(deploymentZoneCtx.VSphereDeploymentZone, metav1.Condition{
			Type:    infrav1.VSphereDeploymentZoneFailureDomainValidatedV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.VSphereDeploymentZoneFailureDomainZoneMisconfiguredV1Beta2Reason,
			Message: err.Error(),
		})
		return errors.Wrapf(err, "failed to reconcile failure domain: zone is not configured correctly")
	}

	if computeCluster := vsphereFailureDomain.Spec.Topology.ComputeCluster; computeCluster != nil {
		if err := r.reconcileComputeCluster(ctx, deploymentZoneCtx, vsphereFailureDomain); err != nil {
			return errors.Wrapf(err, "failed to reconcile failure domain: compute cluster %s is not configured correctly", *computeCluster)
		}
	}

	if err := r.reconcileTopology(ctx, deploymentZoneCtx, vsphereFailureDomain); err != nil {
		return errors.Wrap(err, "failed to reconcile failure domain: topology is not configured correctly")
	}

	// Ensure the VSphereDeploymentZone is marked as an owner of the VSphereFailureDomain.
	if err := updateOwnerReferences(ctx, vsphereFailureDomain, r.Client,
		func() []metav1.OwnerReference {
			return clusterutilv1.EnsureOwnerRef(
				vsphereFailureDomain.OwnerReferences,
				metav1.OwnerReference{
					APIVersion: infrav1.GroupVersion.String(),
					Kind:       "VSphereDeploymentZone",
					Name:       deploymentZoneCtx.VSphereDeploymentZone.Name,
					UID:        deploymentZoneCtx.VSphereDeploymentZone.UID,
				})
		}); err != nil {
		v1beta2conditions.Set(deploymentZoneCtx.VSphereDeploymentZone, metav1.Condition{
			Type:    infrav1.VSphereDeploymentZoneFailureDomainValidatedV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.VSphereDeploymentZoneFailureDomainValidationFailedV1Beta2Reason,
			Message: "failed to update owner references of failure domain",
		})
		return err
	}

	// Mark the VSphereDeploymentZone as having a valid VSphereFailureDomain.
	conditions.MarkTrue(deploymentZoneCtx.VSphereDeploymentZone, infrav1.VSphereFailureDomainValidatedCondition)
	v1beta2conditions.Set(deploymentZoneCtx.VSphereDeploymentZone, metav1.Condition{
		Type:   infrav1.VSphereDeploymentZoneFailureDomainValidatedV1Beta2Condition,
		Status: metav1.ConditionTrue,
		Reason: infrav1.VSphereDeploymentZoneFailureDomainValidatedV1Beta2Reason,
	})
	return nil
}

func (r vsphereDeploymentZoneReconciler) reconcileInfraFailureDomain(ctx context.Context, deploymentZoneCtx *capvcontext.VSphereDeploymentZoneContext, vsphereFailureDomain *infrav1.VSphereFailureDomain, failureDomain infrav1.FailureDomain) error {
	if *failureDomain.AutoConfigure {
		return r.createAndAttachMetadata(ctx, deploymentZoneCtx, vsphereFailureDomain, failureDomain)
	}
	return r.verifyFailureDomain(ctx, deploymentZoneCtx, vsphereFailureDomain, failureDomain)
}

func (r vsphereDeploymentZoneReconciler) reconcileTopology(ctx context.Context, deploymentZoneCtx *capvcontext.VSphereDeploymentZoneContext, vsphereFailureDomain *infrav1.VSphereFailureDomain) error {
	topology := vsphereFailureDomain.Spec.Topology
	if datastore := topology.Datastore; datastore != "" {
		if _, err := deploymentZoneCtx.AuthSession.Finder.Datastore(ctx, datastore); err != nil {
			conditions.MarkFalse(deploymentZoneCtx.VSphereDeploymentZone, infrav1.VSphereFailureDomainValidatedCondition, infrav1.DatastoreNotFoundReason, clusterv1.ConditionSeverityError, "datastore %s is misconfigured", datastore)
			v1beta2conditions.Set(deploymentZoneCtx.VSphereDeploymentZone, metav1.Condition{
				Type:    infrav1.VSphereDeploymentZoneFailureDomainValidatedV1Beta2Condition,
				Status:  metav1.ConditionFalse,
				Reason:  infrav1.VSphereDeploymentZoneFailureDomainDatastoreNotFoundV1Beta2Reason,
				Message: fmt.Sprintf("datastore %s is misconfigured", datastore),
			})
			return errors.Wrapf(err, "unable to find datastore %s", datastore)
		}
	}

	for _, network := range topology.Networks {
		if _, err := deploymentZoneCtx.AuthSession.Finder.Network(ctx, network); err != nil {
			conditions.MarkFalse(deploymentZoneCtx.VSphereDeploymentZone, infrav1.VSphereFailureDomainValidatedCondition, infrav1.NetworkNotFoundReason, clusterv1.ConditionSeverityError, "network %s is not found", network)
			v1beta2conditions.Set(deploymentZoneCtx.VSphereDeploymentZone, metav1.Condition{
				Type:    infrav1.VSphereDeploymentZoneFailureDomainValidatedV1Beta2Condition,
				Status:  metav1.ConditionFalse,
				Reason:  infrav1.VSphereDeploymentZoneFailureDomainNetworkNotFoundV1Beta2Reason,
				Message: fmt.Sprintf("network %s is not found", network),
			})
			return errors.Wrapf(err, "unable to find network %s", network)
		}
	}

	for _, networkConfig := range topology.NetworkConfigurations {
		if _, err := deploymentZoneCtx.AuthSession.Finder.Network(ctx, networkConfig.NetworkName); err != nil {
			conditions.MarkFalse(deploymentZoneCtx.VSphereDeploymentZone, infrav1.VSphereFailureDomainValidatedCondition, infrav1.NetworkNotFoundReason, clusterv1.ConditionSeverityError, "network %s is not found", networkConfig.NetworkName)
			v1beta2conditions.Set(deploymentZoneCtx.VSphereDeploymentZone, metav1.Condition{
				Type:    infrav1.VSphereDeploymentZoneFailureDomainValidatedV1Beta2Condition,
				Status:  metav1.ConditionFalse,
				Reason:  infrav1.VSphereDeploymentZoneFailureDomainNetworkNotFoundV1Beta2Reason,
				Message: fmt.Sprintf("network %s is not found", networkConfig.NetworkName),
			})
			return errors.Wrapf(err, "unable to find network %s", networkConfig.NetworkName)
		}
	}

	if hostPlacementInfo := topology.Hosts; hostPlacementInfo != nil {
		rule, err := cluster.VerifyAffinityRule(ctx, deploymentZoneCtx, *topology.ComputeCluster, hostPlacementInfo.HostGroupName, hostPlacementInfo.VMGroupName)
		if err != nil {
			conditions.MarkFalse(deploymentZoneCtx.VSphereDeploymentZone, infrav1.VSphereFailureDomainValidatedCondition, infrav1.HostsMisconfiguredReason, clusterv1.ConditionSeverityError, "vm host affinity does not exist")
			v1beta2conditions.Set(deploymentZoneCtx.VSphereDeploymentZone, metav1.Condition{
				Type:    infrav1.VSphereDeploymentZoneFailureDomainValidatedV1Beta2Condition,
				Status:  metav1.ConditionFalse,
				Reason:  infrav1.VSphereDeploymentZoneFailureDomainHostsMisconfiguredV1Beta2Reason,
				Message: "vm host affinity rule does not exist",
			})
			return err
		}

		if rule.Disabled() {
			ctrl.LoggerFrom(ctx).V(4).Info("WARNING: vm-host rule for the failure domain is disabled", "hostGroup", hostPlacementInfo.HostGroupName, "vmGroup", hostPlacementInfo.VMGroupName)
		}
	}

	return nil
}

func (r vsphereDeploymentZoneReconciler) reconcileComputeCluster(ctx context.Context, deploymentZoneCtx *capvcontext.VSphereDeploymentZoneContext, vsphereFailureDomain *infrav1.VSphereFailureDomain) error {
	computeCluster := vsphereFailureDomain.Spec.Topology.ComputeCluster
	if computeCluster == nil {
		return nil
	}

	ccr, err := deploymentZoneCtx.AuthSession.Finder.ClusterComputeResource(ctx, *computeCluster)
	if err != nil {
		conditions.MarkFalse(deploymentZoneCtx.VSphereDeploymentZone, infrav1.VSphereFailureDomainValidatedCondition, infrav1.ComputeClusterNotFoundReason, clusterv1.ConditionSeverityError, "compute cluster %s not found", *computeCluster)
		v1beta2conditions.Set(deploymentZoneCtx.VSphereDeploymentZone, metav1.Condition{
			Type:    infrav1.VSphereDeploymentZoneFailureDomainValidatedV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.VSphereDeploymentZoneFailureDomainComputeClusterNotFoundV1Beta2Reason,
			Message: fmt.Sprintf("compute cluster %s not found", *computeCluster),
		})
		return errors.Wrap(err, "compute cluster not found")
	}

	if resourcePool := deploymentZoneCtx.VSphereDeploymentZone.Spec.PlacementConstraint.ResourcePool; resourcePool != "" {
		rp, err := deploymentZoneCtx.AuthSession.Finder.ResourcePool(ctx, resourcePool)
		if err != nil {
			return errors.Wrapf(err, "unable to find resource pool")
		}

		ref, err := rp.Owner(ctx)
		if err != nil {
			conditions.MarkFalse(deploymentZoneCtx.VSphereDeploymentZone, infrav1.VSphereFailureDomainValidatedCondition, infrav1.ComputeClusterNotFoundReason, clusterv1.ConditionSeverityError, "resource pool owner not found")
			v1beta2conditions.Set(deploymentZoneCtx.VSphereDeploymentZone, metav1.Condition{
				Type:    infrav1.VSphereDeploymentZoneFailureDomainValidatedV1Beta2Condition,
				Status:  metav1.ConditionFalse,
				Reason:  infrav1.VSphereDeploymentZoneFailureDomainComputeClusterNotFoundV1Beta2Reason,
				Message: "resource pool owner not found",
			})
			return errors.Wrap(err, "unable to find owner compute resource")
		}
		if ref.Reference() != ccr.Reference() {
			conditions.MarkFalse(deploymentZoneCtx.VSphereDeploymentZone, infrav1.VSphereFailureDomainValidatedCondition, infrav1.ResourcePoolNotFoundReason, clusterv1.ConditionSeverityError, "resource pool is not owned by compute cluster")
			v1beta2conditions.Set(deploymentZoneCtx.VSphereDeploymentZone, metav1.Condition{
				Type:    infrav1.VSphereDeploymentZoneFailureDomainValidatedV1Beta2Condition,
				Status:  metav1.ConditionFalse,
				Reason:  infrav1.VSphereDeploymentZoneFailureDomainResourcePoolNotFoundV1Beta2Reason,
				Message: "resource pool is not owned by compute cluster",
			})
			return errors.Errorf("compute cluster %s does not own resource pool %s", *computeCluster, resourcePool)
		}
	}
	return nil
}

// verifyFailureDomain verifies the Failure Domain. It verifies the existence of tag and category specified and
// checks whether the specified tags exist on the DataCenter or Compute Cluster or Hosts (in a HostGroup).
func (r vsphereDeploymentZoneReconciler) verifyFailureDomain(ctx context.Context, deploymentZoneCtx *capvcontext.VSphereDeploymentZoneContext, vsphereFailureDomain *infrav1.VSphereFailureDomain, failureDomain infrav1.FailureDomain) error {
	if _, err := deploymentZoneCtx.AuthSession.TagManager.GetTagForCategory(ctx, failureDomain.Name, failureDomain.TagCategory); err != nil {
		return errors.Wrapf(err, "failed to verify tag %s and category %s", failureDomain.Name, failureDomain.TagCategory)
	}

	objects, err := taggable.GetObjects(ctx, deploymentZoneCtx, vsphereFailureDomain, failureDomain.Type)
	if err != nil {
		return errors.Wrapf(err, "failed to get objects of type %s", failureDomain.Type)
	}

	// All the objects should be associated to the tag
	for _, obj := range objects {
		hasTag, err := obj.HasTag(ctx, failureDomain.Name)
		if err != nil {
			return errors.Wrapf(err, "failed to verify if object %s has tag %s", obj, failureDomain.Name)
		}
		if !hasTag {
			return errors.Errorf("object %s does not have tag %s", obj, failureDomain.Name)
		}
	}
	return nil
}

func (r vsphereDeploymentZoneReconciler) createAndAttachMetadata(ctx context.Context, deploymentZoneCtx *capvcontext.VSphereDeploymentZoneContext, vsphereFailureDomain *infrav1.VSphereFailureDomain, failureDomain infrav1.FailureDomain) error {
	log := ctrl.LoggerFrom(ctx, "tagName", failureDomain.Name, "tagCategory", failureDomain.TagCategory, "failureDomainType", failureDomain.Type)
	categoryID, err := metadata.CreateCategory(ctx, deploymentZoneCtx, failureDomain.TagCategory, failureDomain.Type)
	if err != nil {
		return errors.Wrapf(err, "failed to create category %s", failureDomain.TagCategory)
	}
	err = metadata.CreateTag(ctx, deploymentZoneCtx, failureDomain.Name, categoryID)
	if err != nil {
		return errors.Wrapf(err, "failed to create tag %s", failureDomain.Name)
	}

	objects, err := taggable.GetObjects(ctx, deploymentZoneCtx, vsphereFailureDomain, failureDomain.Type)
	if err != nil {
		return errors.Wrapf(err, "failed to get objects of type %s", failureDomain.Type)
	}

	var errList []error
	for _, obj := range objects {
		log.V(4).Info("Attaching tag to object")
		err := obj.AttachTag(ctx, failureDomain.Name)
		if err != nil {
			errList = append(errList, errors.Wrapf(err, "failed to attach tag %s to object %s", failureDomain.Name, obj))
		}
	}
	return kerrors.NewAggregate(errList)
}
