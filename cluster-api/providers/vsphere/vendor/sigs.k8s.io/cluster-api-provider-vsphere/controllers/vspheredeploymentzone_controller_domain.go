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
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apierrors "k8s.io/apimachinery/pkg/util/errors"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	clusterutilv1 "sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/conditions"
	ctrl "sigs.k8s.io/controller-runtime"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services/govmomi/cluster"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services/govmomi/metadata"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/taggable"
)

func (r vsphereDeploymentZoneReconciler) reconcileFailureDomain(ctx *context.VSphereDeploymentZoneContext) error {
	logger := ctrl.LoggerFrom(ctx).WithValues("failure domain", ctx.VSphereFailureDomain.Name)

	// verify the failure domain for the region
	if err := r.reconcileInfraFailureDomain(ctx, ctx.VSphereFailureDomain.Spec.Region); err != nil {
		conditions.MarkFalse(ctx.VSphereDeploymentZone, infrav1.VSphereFailureDomainValidatedCondition, infrav1.RegionMisconfiguredReason, clusterv1.ConditionSeverityError, err.Error())
		logger.Error(err, "region is not configured correctly")
		return errors.Wrapf(err, "region is not configured correctly")
	}

	// verify the failure domain for the zone
	if err := r.reconcileInfraFailureDomain(ctx, ctx.VSphereFailureDomain.Spec.Zone); err != nil {
		conditions.MarkFalse(ctx.VSphereDeploymentZone, infrav1.VSphereFailureDomainValidatedCondition, infrav1.ZoneMisconfiguredReason, clusterv1.ConditionSeverityError, err.Error())
		logger.Error(err, "zone is not configured correctly")
		return errors.Wrapf(err, "zone is not configured correctly")
	}

	if computeCluster := ctx.VSphereFailureDomain.Spec.Topology.ComputeCluster; computeCluster != nil {
		if err := r.reconcileComputeCluster(ctx); err != nil {
			logger.Error(err, "compute cluster is not configured correctly", "name", *computeCluster)
			return errors.Wrap(err, "compute cluster is not configured correctly")
		}
	}

	if err := r.reconcileTopology(ctx); err != nil {
		logger.Error(err, "topology is not configured correctly")
		return errors.Wrap(err, "topology is not configured correctly")
	}

	// Ensure the VSphereDeploymentZone is marked as an owner of the VSphereFailureDomain.
	if err := updateOwnerReferences(ctx, ctx.VSphereFailureDomain, r.Client,
		func() []metav1.OwnerReference {
			return clusterutilv1.EnsureOwnerRef(
				ctx.VSphereFailureDomain.OwnerReferences,
				metav1.OwnerReference{
					APIVersion: infrav1.GroupVersion.String(),
					Kind:       ctx.VSphereDeploymentZone.Kind,
					Name:       ctx.VSphereDeploymentZone.Name,
					UID:        ctx.VSphereDeploymentZone.UID,
				})
		}); err != nil {
		return err
	}

	// Mark the VSphereDeploymentZone as having a valid VSphereFailureDomain.
	conditions.MarkTrue(ctx.VSphereDeploymentZone, infrav1.VSphereFailureDomainValidatedCondition)
	return nil
}

func (r vsphereDeploymentZoneReconciler) reconcileInfraFailureDomain(ctx *context.VSphereDeploymentZoneContext, failureDomain infrav1.FailureDomain) error {
	if *failureDomain.AutoConfigure {
		return r.createAndAttachMetadata(ctx, failureDomain)
	}
	return r.verifyFailureDomain(ctx, failureDomain)
}

func (r vsphereDeploymentZoneReconciler) reconcileTopology(ctx *context.VSphereDeploymentZoneContext) error {
	topology := ctx.VSphereFailureDomain.Spec.Topology
	if datastore := topology.Datastore; datastore != "" {
		if _, err := ctx.AuthSession.Finder.Datastore(ctx, datastore); err != nil {
			conditions.MarkFalse(ctx.VSphereDeploymentZone, infrav1.VSphereFailureDomainValidatedCondition, infrav1.DatastoreNotFoundReason, clusterv1.ConditionSeverityError, "datastore %s is misconfigured", datastore)
			return errors.Wrapf(err, "unable to find datastore %s", datastore)
		}
	}

	for _, network := range topology.Networks {
		if _, err := ctx.AuthSession.Finder.Network(ctx, network); err != nil {
			conditions.MarkFalse(ctx.VSphereDeploymentZone, infrav1.VSphereFailureDomainValidatedCondition, infrav1.NetworkNotFoundReason, clusterv1.ConditionSeverityError, "network %s is misconfigured", network)
			return errors.Wrapf(err, "unable to find network %s", network)
		}
	}

	if hostPlacementInfo := topology.Hosts; hostPlacementInfo != nil {
		rule, err := cluster.VerifyAffinityRule(ctx, *topology.ComputeCluster, hostPlacementInfo.HostGroupName, hostPlacementInfo.VMGroupName)
		switch {
		case err != nil:
			conditions.MarkFalse(ctx.VSphereDeploymentZone, infrav1.VSphereFailureDomainValidatedCondition, infrav1.HostsMisconfiguredReason, clusterv1.ConditionSeverityError, "vm host affinity does not exist")
			return err
		case rule.Disabled():
			ctrl.LoggerFrom(ctx).V(4).Info("warning: vm-host rule for the failure domain is disabled", "hostgroup", hostPlacementInfo.HostGroupName, "vmGroup", hostPlacementInfo.VMGroupName)
			conditions.MarkFalse(ctx.VSphereDeploymentZone, infrav1.VSphereFailureDomainValidatedCondition, infrav1.HostsAffinityMisconfiguredReason, clusterv1.ConditionSeverityWarning, "vm host affinity is disabled")
		default:
			conditions.MarkTrue(ctx.VSphereDeploymentZone, infrav1.VSphereFailureDomainValidatedCondition)
		}
	}
	return nil
}

func (r vsphereDeploymentZoneReconciler) reconcileComputeCluster(ctx *context.VSphereDeploymentZoneContext) error {
	computeCluster := ctx.VSphereFailureDomain.Spec.Topology.ComputeCluster
	if computeCluster == nil {
		return nil
	}

	ccr, err := ctx.AuthSession.Finder.ClusterComputeResource(ctx, *computeCluster)
	if err != nil {
		conditions.MarkFalse(ctx.VSphereDeploymentZone, infrav1.VSphereFailureDomainValidatedCondition, infrav1.ComputeClusterNotFoundReason, clusterv1.ConditionSeverityError, "compute cluster %s not found", *computeCluster)
		return errors.Wrap(err, "compute cluster not found")
	}

	if resourcePool := ctx.VSphereDeploymentZone.Spec.PlacementConstraint.ResourcePool; resourcePool != "" {
		rp, err := ctx.AuthSession.Finder.ResourcePool(ctx, resourcePool)
		if err != nil {
			return errors.Wrapf(err, "unable to find resource pool")
		}

		ref, err := rp.Owner(ctx)
		if err != nil {
			conditions.MarkFalse(ctx.VSphereDeploymentZone, infrav1.VSphereFailureDomainValidatedCondition, infrav1.ComputeClusterNotFoundReason, clusterv1.ConditionSeverityError, "resource pool owner not found")
			return errors.Wrap(err, "unable to find owner compute resource")
		}
		if ref.Reference() != ccr.Reference() {
			conditions.MarkFalse(ctx.VSphereDeploymentZone, infrav1.VSphereFailureDomainValidatedCondition, infrav1.ResourcePoolNotFoundReason, clusterv1.ConditionSeverityError, "resource pool is not owned by compute cluster")
			return errors.Errorf("compute cluster %s does not own resource pool %s", *computeCluster, resourcePool)
		}
	}
	return nil
}

// verifyFailureDomain verifies the Failure Domain. It verifies the existence of tag and category specified and
// checks whether the specified tags exist on the DataCenter or Compute Cluster or Hosts (in a HostGroup).
func (r vsphereDeploymentZoneReconciler) verifyFailureDomain(ctx *context.VSphereDeploymentZoneContext, failureDomain infrav1.FailureDomain) error {
	if _, err := ctx.AuthSession.TagManager.GetTagForCategory(ctx, failureDomain.Name, failureDomain.TagCategory); err != nil {
		return errors.Wrapf(err, "failed to verify tag %s and category %s", failureDomain.Name, failureDomain.TagCategory)
	}

	objects, err := taggable.GetObjects(ctx, failureDomain.Type)
	if err != nil {
		return errors.Wrapf(err, "failed to find object")
	}

	// All the objects should be associated to the tag
	for _, obj := range objects {
		hasTag, err := obj.HasTag(ctx, failureDomain.Name)
		if err != nil {
			return errors.Wrapf(err, "failed to verify tag association")
		}
		if !hasTag {
			return errors.Errorf("tag %s is not associated to object %s", failureDomain.Name, obj)
		}
	}
	return nil
}

func (r vsphereDeploymentZoneReconciler) createAndAttachMetadata(ctx *context.VSphereDeploymentZoneContext, failureDomain infrav1.FailureDomain) error {
	logger := ctrl.LoggerFrom(ctx, "tag", failureDomain.Name, "category", failureDomain.TagCategory)
	categoryID, err := metadata.CreateCategory(ctx, failureDomain.TagCategory, failureDomain.Type)
	if err != nil {
		logger.V(4).Error(err, "category creation failed")
		return errors.Wrapf(err, "failed to create category %s", failureDomain.TagCategory)
	}
	err = metadata.CreateTag(ctx, failureDomain.Name, categoryID)
	if err != nil {
		logger.V(4).Error(err, "tag creation failed")
		return errors.Wrapf(err, "failed to create tag %s", failureDomain.Name)
	}

	logger = logger.WithValues("type", failureDomain.Type)
	objects, err := taggable.GetObjects(ctx, failureDomain.Type)
	if err != nil {
		logger.V(4).Error(err, "failed to find object")
		return err
	}

	var errList []error
	for _, obj := range objects {
		logger.V(4).Info("attaching tag to object")
		err := obj.AttachTag(ctx, failureDomain.Name)
		if err != nil {
			logger.V(4).Error(err, "failed to find object")
			errList = append(errList, errors.Wrapf(err, "failed to attach tag"))
		}
	}
	return apierrors.NewAggregate(errList)
}
