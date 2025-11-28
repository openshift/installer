/*
Copyright 2025 The ORC Authors.

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

package volume

import (
	"strconv"

	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	orcapplyconfigv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/pkg/clients/applyconfiguration/api/v1alpha1"
)

// Ideally, these constants are defined in gophercloud.
const (
	VolumeStatusAvailable = "available"
	VolumeStatusInUse     = "in-use"
	VolumeStatusDeleting  = "deleting"
)

type volumeStatusWriter struct{}

type objectApplyT = orcapplyconfigv1alpha1.VolumeApplyConfiguration
type statusApplyT = orcapplyconfigv1alpha1.VolumeStatusApplyConfiguration

var _ interfaces.ResourceStatusWriter[*orcv1alpha1.Volume, *osResourceT, *objectApplyT, *statusApplyT] = volumeStatusWriter{}

func (volumeStatusWriter) GetApplyConfig(name, namespace string) *objectApplyT {
	return orcapplyconfigv1alpha1.Volume(name, namespace)
}

func (volumeStatusWriter) ResourceAvailableStatus(orcObject *orcv1alpha1.Volume, osResource *osResourceT) (metav1.ConditionStatus, progress.ReconcileStatus) {
	if osResource == nil {
		if orcObject.Status.ID == nil {
			return metav1.ConditionFalse, nil
		} else {
			return metav1.ConditionUnknown, nil
		}
	}

	if osResource.Status == VolumeStatusAvailable || osResource.Status == VolumeStatusInUse {
		return metav1.ConditionTrue, nil
	}

	// Otherwise we should continue to poll
	return metav1.ConditionFalse, progress.WaitingOnOpenStack(progress.WaitingOnReady, volumeAvailablePollingPeriod)
}

func (volumeStatusWriter) ApplyResourceStatus(log logr.Logger, osResource *osResourceT, statusApply *statusApplyT) {
	resourceStatus := orcapplyconfigv1alpha1.VolumeResourceStatus().
		WithName(osResource.Name).
		WithVolumeType(osResource.VolumeType).
		WithSize(int32(osResource.Size)).
		WithStatus(osResource.Status).
		WithUserID(osResource.UserID).
		WithEncrypted(osResource.Encrypted).
		WithMultiattach(osResource.Multiattach).
		WithCreatedAt(metav1.NewTime(osResource.CreatedAt))

	if !osResource.UpdatedAt.IsZero() {
		resourceStatus.WithUpdatedAt(metav1.NewTime(osResource.UpdatedAt))
	}

	if osResource.Description != "" {
		resourceStatus.WithDescription(osResource.Description)
	}

	if osResource.Bootable != "" {
		boolValue, err := strconv.ParseBool(osResource.Bootable)
		if err != nil {
			log.Info("Failed to parse boolean value", err)
		} else {
			resourceStatus.WithBootable(boolValue)
		}
	}

	for k, v := range osResource.Metadata {
		resourceStatus.WithMetadata(orcapplyconfigv1alpha1.VolumeMetadataStatus().
			WithName(k).
			WithValue(v))
	}

	if osResource.AvailabilityZone != "" {
		resourceStatus.WithAvailabilityZone(osResource.AvailabilityZone)
	}

	if osResource.SnapshotID != "" {
		resourceStatus.WithSnapshotID(osResource.SnapshotID)
	}

	if osResource.SourceVolID != "" {
		resourceStatus.WithSourceVolID(osResource.SourceVolID)
	}

	if osResource.BackupID != nil {
		resourceStatus.WithBackupID(*osResource.BackupID)
	}

	if osResource.ReplicationStatus != "" {
		resourceStatus.WithReplicationStatus(osResource.ReplicationStatus)
	}

	if osResource.ConsistencyGroupID != "" {
		resourceStatus.WithConsistencyGroupID(osResource.ConsistencyGroupID)
	}

	if osResource.Host != "" {
		resourceStatus.WithHost(osResource.Host)
	}

	if osResource.TenantID != "" {
		resourceStatus.WithTenantID(osResource.TenantID)
	}

	for i := range osResource.Attachments {
		resourceStatus.WithAttachments(orcapplyconfigv1alpha1.VolumeAttachmentStatus().
			WithAttachmentID(osResource.Attachments[i].AttachmentID).
			WithServerID(osResource.Attachments[i].ServerID).
			WithDevice(osResource.Attachments[i].Device).
			WithAttachedAt(metav1.NewTime(osResource.Attachments[i].AttachedAt)))
	}

	statusApply.WithResource(resourceStatus)
}
