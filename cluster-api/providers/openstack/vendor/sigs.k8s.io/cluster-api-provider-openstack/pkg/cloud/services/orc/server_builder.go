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
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"

	infrav1alpha1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha1"
	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta2"
)

// VolumeAttachment links an ORC Volume name to the device name used
// when attaching to the server.
type VolumeAttachment struct {
	ORCName string
	Device  string
}

// buildServer builds a managed ORC Server referencing all previously
// created ORC sub-resources (Image/Flavor/KeyPair/Ports/Volumes/etc.).
func buildServer(
	serverName, namespace string,
	spec *infrav1alpha1.OpenStackServerSpec,
	imageORCName, flavorORCName string,
	keypairORCName, serverGroupORCName, rootVolumeORCName string,
	portORCNames []string,
	additionalVolumes []VolumeAttachment,
	credRef orcv1alpha1.CloudCredentialsReference,
) *orcv1alpha1.Server {
	serverSpec := &orcv1alpha1.ServerResourceSpec{
		FlavorRef: orcv1alpha1.KubernetesNameRef(flavorORCName),
	}

	// Boot mode: boot-from-volume or boot-from-image
	if rootVolumeORCName != "" {
		serverSpec.BootVolume = &orcv1alpha1.ServerBootVolumeSpec{
			VolumeRef: orcv1alpha1.KubernetesNameRef(rootVolumeORCName),
		}
	} else {
		ref := orcv1alpha1.KubernetesNameRef(imageORCName)
		serverSpec.ImageRef = &ref
	}

	// Ports
	for _, portName := range portORCNames {
		ref := orcv1alpha1.KubernetesNameRef(portName)
		serverSpec.Ports = append(serverSpec.Ports, orcv1alpha1.ServerPortSpec{
			PortRef: &ref,
		})
	}

	// Additional volumes
	for _, vol := range additionalVolumes {
		device := vol.Device
		serverSpec.Volumes = append(serverSpec.Volumes, orcv1alpha1.ServerVolumeSpec{
			VolumeRef: orcv1alpha1.KubernetesNameRef(vol.ORCName),
			Device:    &device,
		})
	}

	// KeyPair
	if keypairORCName != "" {
		ref := orcv1alpha1.KubernetesNameRef(keypairORCName)
		serverSpec.KeypairRef = &ref
	}

	// Availability zone
	if spec.AvailabilityZone != nil {
		serverSpec.AvailabilityZone = *spec.AvailabilityZone
	}

	// Tags
	for _, t := range spec.Tags {
		serverSpec.Tags = append(serverSpec.Tags, orcv1alpha1.ServerTag(t))
	}

	// Config drive
	if spec.ConfigDrive != nil {
		serverSpec.ConfigDrive = spec.ConfigDrive
	}

	// Metadata
	for _, m := range spec.ServerMetadata {
		serverSpec.Metadata = append(serverSpec.Metadata, orcv1alpha1.ServerMetadata{
			Key:   m.Key,
			Value: m.Value,
		})
	}

	// User data
	if spec.UserDataRef != nil {
		ref := orcv1alpha1.KubernetesNameRef(spec.UserDataRef.Name)
		serverSpec.UserData = &orcv1alpha1.UserDataSpec{
			SecretRef: &ref,
		}
	}

	// Scheduler hints
	hints := convertSchedulerHints(serverGroupORCName, spec.SchedulerHintAdditionalProperties)
	if hints != nil {
		serverSpec.SchedulerHints = hints
	}

	return &orcv1alpha1.Server{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ServerName(serverName),
			Namespace: namespace,
		},
		Spec: orcv1alpha1.ServerSpec{
			ManagementPolicy:    orcv1alpha1.ManagementPolicyManaged,
			Resource:            serverSpec,
			CloudCredentialsRef: credRef,
		},
	}
}

// convertSchedulerHints converts CAPO scheduler hints to ORC format.
func convertSchedulerHints(serverGroupORCName string, props []infrav1.SchedulerHintAdditionalProperty) *orcv1alpha1.ServerSchedulerHints {
	if serverGroupORCName == "" && len(props) == 0 {
		return nil
	}

	hints := &orcv1alpha1.ServerSchedulerHints{}

	if serverGroupORCName != "" {
		ref := orcv1alpha1.KubernetesNameRef(serverGroupORCName)
		hints.ServerGroupRef = &ref
	}

	if len(props) > 0 {
		hints.AdditionalProperties = make(map[string]string, len(props))
		for _, prop := range props {
			var value string
			switch prop.Value.Type {
			case infrav1.SchedulerHintTypeBool:
				value = strconv.FormatBool(ptr.Deref(prop.Value.Bool, false))
			case infrav1.SchedulerHintTypeString:
				value = ptr.Deref(prop.Value.String, "")
			case infrav1.SchedulerHintTypeNumber:
				value = strconv.Itoa(int(ptr.Deref(prop.Value.Number, 0)))
			}
			hints.AdditionalProperties[prop.Name] = value
		}
	}

	return hints
}
