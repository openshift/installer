/*
Copyright 2026 The Kubernetes Authors.

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

package orc

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta2"
)

// buildRootVolume builds a managed ORC Volume for the server's root disk.
// The volume is bootable (created from an image via imageRef).
func buildRootVolume(
	serverName, namespace string,
	rootVolume *infrav1.RootVolume,
	imageORCName string,
	volumeTypeNameMap map[string]string,
	serverAZ string,
	credRef orcv1alpha1.CloudCredentialsReference,
) *orcv1alpha1.Volume {
	imageRef := orcv1alpha1.KubernetesNameRef(imageORCName)
	volSpec := &orcv1alpha1.VolumeResourceSpec{
		Size:     rootVolume.SizeGiB,
		ImageRef: &imageRef,
	}

	// Volume type
	if rootVolume.Type != "" {
		if vtName, ok := volumeTypeNameMap[rootVolume.Type]; ok {
			vtRef := orcv1alpha1.KubernetesNameRef(vtName)
			volSpec.VolumeTypeRef = &vtRef
		}
	}

	// Availability zone
	volSpec.AvailabilityZone = resolveVolumeAZ(&rootVolume.BlockDeviceVolume, serverAZ)

	return &orcv1alpha1.Volume{
		ObjectMeta: metav1.ObjectMeta{
			Name:      RootVolumeName(serverName),
			Namespace: namespace,
		},
		Spec: orcv1alpha1.VolumeSpec{
			ManagementPolicy:    orcv1alpha1.ManagementPolicyManaged,
			Resource:            volSpec,
			CloudCredentialsRef: credRef,
		},
	}
}

// buildAdditionalVolume builds a managed ORC Volume for an additional
// block device. Only supports VolumeBlockDevice type; LocalBlockDevice
// is not supported by ORC and should be skipped by the caller.
func buildAdditionalVolume(
	serverName, namespace string,
	bd infrav1.AdditionalBlockDevice,
	volumeTypeNameMap map[string]string,
	serverAZ string,
	credRef orcv1alpha1.CloudCredentialsReference,
) *orcv1alpha1.Volume {
	volSpec := &orcv1alpha1.VolumeResourceSpec{
		Size: bd.SizeGiB,
	}

	if bd.Storage.Volume != nil {
		if bd.Storage.Volume.Type != "" {
			if vtName, ok := volumeTypeNameMap[bd.Storage.Volume.Type]; ok {
				vtRef := orcv1alpha1.KubernetesNameRef(vtName)
				volSpec.VolumeTypeRef = &vtRef
			}
		}
		volSpec.AvailabilityZone = resolveVolumeAZ(bd.Storage.Volume, serverAZ)
	}

	return &orcv1alpha1.Volume{
		ObjectMeta: metav1.ObjectMeta{
			Name:      AdditionalVolumeName(serverName, bd.Name),
			Namespace: namespace,
		},
		Spec: orcv1alpha1.VolumeSpec{
			ManagementPolicy:    orcv1alpha1.ManagementPolicyManaged,
			Resource:            volSpec,
			CloudCredentialsRef: credRef,
		},
	}
}

// resolveVolumeAZ determines the availability zone for a volume based
// on the CAPO VolumeAvailabilityZone configuration.
func resolveVolumeAZ(vol *infrav1.BlockDeviceVolume, serverAZ string) string {
	if vol == nil || vol.AvailabilityZone == nil {
		return ""
	}
	switch vol.AvailabilityZone.From {
	case "", infrav1.VolumeAZFromName:
		if vol.AvailabilityZone.Name != nil {
			return string(*vol.AvailabilityZone.Name)
		}
	case infrav1.VolumeAZFromMachine:
		return serverAZ
	}
	return ""
}
