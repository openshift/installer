/*
Copyright 2024 The ORC Authors.

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

package server

import (
	"fmt"

	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	orcapplyconfigv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/pkg/clients/applyconfiguration/api/v1alpha1"
)

const (
	ServerStatusActive = "ACTIVE"
	ServerStatusBuild  = "BUILD"
	ServerStatusError  = "ERROR"
)

type objectApplyPT = *orcapplyconfigv1alpha1.ServerApplyConfiguration
type statusApplyPT = *orcapplyconfigv1alpha1.ServerStatusApplyConfiguration

type serverStatusWriter struct{}

var _ interfaces.ResourceStatusWriter[orcObjectPT, *osResourceT, objectApplyPT, statusApplyPT] = serverStatusWriter{}

func (serverStatusWriter) GetApplyConfig(name, namespace string) objectApplyPT {
	return orcapplyconfigv1alpha1.Server(name, namespace)
}

func (serverStatusWriter) ResourceAvailableStatus(orcObject orcObjectPT, osResource *osResourceT) (metav1.ConditionStatus, progress.ReconcileStatus) {
	if osResource == nil {
		if orcObject.Status.ID == nil {
			return metav1.ConditionFalse, nil
		} else {
			return metav1.ConditionUnknown, nil
		}
	}

	if osResource.Status == ServerStatusActive {
		return metav1.ConditionTrue, nil
	}
	// We should continue to poll if the status is not ACTIVE
	return metav1.ConditionFalse, progress.WaitingOnOpenStack(progress.WaitingOnReady, serverActivePollingPeriod)
}

func (serverStatusWriter) ApplyResourceStatus(log logr.Logger, osResource *osResourceT, statusApply statusApplyPT) {
	// TODO: Add the rest of the OpenStack data to Status
	status := orcapplyconfigv1alpha1.ServerResourceStatus().
		WithName(osResource.Name).
		WithStatus(osResource.Status).
		WithHostID(osResource.HostID).
		WithAvailabilityZone(osResource.AvailabilityZone).
		WithServerGroups(ptr.Deref(osResource.ServerGroups, []string{})...).
		WithTags(ptr.Deref(osResource.Tags, []string{})...)

	if imageID, ok := osResource.Image["id"]; ok {
		status.WithImageID(fmt.Sprintf("%s", imageID))
	}

	for i := range osResource.AttachedVolumes {
		status.WithVolumes(orcapplyconfigv1alpha1.ServerVolumeStatus().
			WithID(osResource.AttachedVolumes[i].ID))
	}

	for i := range osResource.Interfaces {
		iface := osResource.Interfaces[i]
		interfaceStatus := orcapplyconfigv1alpha1.ServerInterfaceStatus().
			WithPortID(iface.PortID).
			WithNetID(iface.NetID).
			WithMACAddr(iface.MACAddr).
			WithPortState(iface.PortState)

		for j := range iface.FixedIPs {
			interfaceStatus.WithFixedIPs(orcapplyconfigv1alpha1.ServerInterfaceFixedIP().
				WithIPAddress(iface.FixedIPs[j].IPAddress).
				WithSubnetID(iface.FixedIPs[j].SubnetID))
		}

		status.WithInterfaces(interfaceStatus)
	}

	statusApply.WithResource(status)
}
