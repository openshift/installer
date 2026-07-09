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

// Type conversion helpers for CAPO → ORC filter types.

func toOpenStackName(s string) *orcv1alpha1.OpenStackName {
	if s == "" {
		return nil
	}
	n := orcv1alpha1.OpenStackName(s)
	return &n
}

func toNeutronDescription(s string) *orcv1alpha1.NeutronDescription {
	if s == "" {
		return nil
	}
	d := orcv1alpha1.NeutronDescription(s)
	return &d
}

func convertNeutronTags(tags infrav1.FilterByNeutronTags) orcv1alpha1.FilterByNeutronTags {
	var result orcv1alpha1.FilterByNeutronTags
	for _, t := range tags.Tags {
		result.Tags = append(result.Tags, orcv1alpha1.NeutronTag(t))
	}
	for _, t := range tags.TagsAny {
		result.TagsAny = append(result.TagsAny, orcv1alpha1.NeutronTag(t))
	}
	for _, t := range tags.NotTags {
		result.NotTags = append(result.NotTags, orcv1alpha1.NeutronTag(t))
	}
	for _, t := range tags.NotTagsAny {
		result.NotTagsAny = append(result.NotTagsAny, orcv1alpha1.NeutronTag(t))
	}
	return result
}

// buildImage builds an unmanaged ORC Image for importing an existing
// Glance image. Returns nil if the user specified an imageRef (an
// existing ORC Image that should be used directly).
func buildImage(serverName, namespace string, imageParam infrav1.ImageParam, credRef orcv1alpha1.CloudCredentialsReference) *orcv1alpha1.Image {
	// If the user specified an ORC Image reference, don't create a wrapper
	if imageParam.ImageRef != nil {
		return nil
	}

	name := ImageName(serverName)
	img := &orcv1alpha1.Image{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: orcv1alpha1.ImageSpec{
			ManagementPolicy:    orcv1alpha1.ManagementPolicyUnmanaged,
			CloudCredentialsRef: credRef,
		},
	}

	switch {
	case imageParam.ID != nil:
		id := *imageParam.ID
		img.Spec.Import = &orcv1alpha1.ImageImport{ID: &id}
	case imageParam.Filter != nil:
		orcFilter := &orcv1alpha1.ImageFilter{}
		if imageParam.Filter.Name != nil {
			orcFilter.Name = toOpenStackName(*imageParam.Filter.Name)
		}
		for _, tag := range imageParam.Filter.Tags {
			orcFilter.Tags = append(orcFilter.Tags, orcv1alpha1.ImageTag(tag))
		}
		img.Spec.Import = &orcv1alpha1.ImageImport{Filter: orcFilter}
	}

	return img
}

// buildFlavor builds an unmanaged ORC Flavor for importing an existing
// Nova flavor. FlavorID takes precedence over Flavor name.
func buildFlavor(serverName, namespace string, flavor *string, flavorID *string, credRef orcv1alpha1.CloudCredentialsReference) *orcv1alpha1.Flavor {
	name := FlavorName(serverName)
	f := &orcv1alpha1.Flavor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: orcv1alpha1.FlavorSpec{
			ManagementPolicy:    orcv1alpha1.ManagementPolicyUnmanaged,
			CloudCredentialsRef: credRef,
		},
	}

	switch {
	case flavorID != nil:
		id := *flavorID
		f.Spec.Import = &orcv1alpha1.FlavorImport{ID: &id}
	case flavor != nil:
		f.Spec.Import = &orcv1alpha1.FlavorImport{
			Filter: &orcv1alpha1.FlavorFilter{
				Name: toOpenStackName(*flavor),
			},
		}
	}

	return f
}

// buildKeypair builds an unmanaged ORC KeyPair for importing an existing
// Nova keypair by name.
func buildKeypair(serverName, namespace, sshKeyName string, credRef orcv1alpha1.CloudCredentialsReference) *orcv1alpha1.KeyPair {
	return &orcv1alpha1.KeyPair{
		ObjectMeta: metav1.ObjectMeta{
			Name:      KeyPairName(serverName),
			Namespace: namespace,
		},
		Spec: orcv1alpha1.KeyPairSpec{
			ManagementPolicy: orcv1alpha1.ManagementPolicyUnmanaged,
			Import: &orcv1alpha1.KeyPairImport{
				Filter: &orcv1alpha1.KeyPairFilter{
					Name: toOpenStackName(sshKeyName),
				},
			},
			CloudCredentialsRef: credRef,
		},
	}
}

// buildServerGroup builds an unmanaged ORC ServerGroup for importing an
// existing Nova server group.
func buildServerGroup(serverName, namespace string, sg *infrav1.ServerGroupParam, credRef orcv1alpha1.CloudCredentialsReference) *orcv1alpha1.ServerGroup {
	obj := &orcv1alpha1.ServerGroup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ServerGroupORCName(serverName),
			Namespace: namespace,
		},
		Spec: orcv1alpha1.ServerGroupSpec{
			ManagementPolicy:    orcv1alpha1.ManagementPolicyUnmanaged,
			CloudCredentialsRef: credRef,
		},
	}

	switch {
	case sg.ID != nil:
		id := *sg.ID
		obj.Spec.Import = &orcv1alpha1.ServerGroupImport{ID: &id}
	case sg.Filter != nil && sg.Filter.Name != nil:
		obj.Spec.Import = &orcv1alpha1.ServerGroupImport{
			Filter: &orcv1alpha1.ServerGroupFilter{
				Name: toOpenStackName(*sg.Filter.Name),
			},
		}
	}

	return obj
}

// buildNetwork builds an unmanaged ORC Network for importing an existing
// Neutron network.
func buildNetwork(serverName, namespace string, param infrav1.NetworkParam, credRef orcv1alpha1.CloudCredentialsReference) *orcv1alpha1.Network {
	key := NetworkParamKey(param)
	obj := &orcv1alpha1.Network{
		ObjectMeta: metav1.ObjectMeta{
			Name:      NetworkORCName(serverName, key),
			Namespace: namespace,
		},
		Spec: orcv1alpha1.NetworkSpec{
			ManagementPolicy:    orcv1alpha1.ManagementPolicyUnmanaged,
			CloudCredentialsRef: credRef,
		},
	}

	switch {
	case param.ID != nil:
		id := *param.ID
		obj.Spec.Import = &orcv1alpha1.NetworkImport{ID: &id}
	case param.Filter != nil:
		orcFilter := &orcv1alpha1.NetworkFilter{
			Name:                toOpenStackName(param.Filter.Name),
			Description:         toNeutronDescription(param.Filter.Description),
			FilterByNeutronTags: convertNeutronTags(param.Filter.FilterByNeutronTags),
		}
		// Note: CAPO ProjectID cannot be mapped to ORC ProjectRef
		obj.Spec.Import = &orcv1alpha1.NetworkImport{Filter: orcFilter}
	}

	return obj
}

// buildSubnet builds an unmanaged ORC Subnet for importing an existing
// Neutron subnet.
func buildSubnet(serverName, namespace string, param infrav1.SubnetParam, credRef orcv1alpha1.CloudCredentialsReference) *orcv1alpha1.Subnet {
	key := SubnetParamKey(param)
	obj := &orcv1alpha1.Subnet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      SubnetORCName(serverName, key),
			Namespace: namespace,
		},
		Spec: orcv1alpha1.SubnetSpec{
			ManagementPolicy:    orcv1alpha1.ManagementPolicyUnmanaged,
			CloudCredentialsRef: credRef,
		},
	}

	switch {
	case param.ID != nil:
		id := *param.ID
		obj.Spec.Import = &orcv1alpha1.SubnetImport{ID: &id}
	case param.Filter != nil:
		orcFilter := &orcv1alpha1.SubnetFilter{
			Name:                toOpenStackName(param.Filter.Name),
			Description:         toNeutronDescription(param.Filter.Description),
			FilterByNeutronTags: convertNeutronTags(param.Filter.FilterByNeutronTags),
		}
		obj.Spec.Import = &orcv1alpha1.SubnetImport{Filter: orcFilter}
	}

	return obj
}

// buildSecurityGroup builds an unmanaged ORC SecurityGroup for importing
// an existing Neutron security group.
func buildSecurityGroup(serverName, namespace string, param infrav1.SecurityGroupParam, credRef orcv1alpha1.CloudCredentialsReference) *orcv1alpha1.SecurityGroup {
	key := SecurityGroupParamKey(param)
	obj := &orcv1alpha1.SecurityGroup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      SecurityGroupORCName(serverName, key),
			Namespace: namespace,
		},
		Spec: orcv1alpha1.SecurityGroupSpec{
			ManagementPolicy:    orcv1alpha1.ManagementPolicyUnmanaged,
			CloudCredentialsRef: credRef,
		},
	}

	switch {
	case param.ID != nil:
		id := *param.ID
		obj.Spec.Import = &orcv1alpha1.SecurityGroupImport{ID: &id}
	case param.Filter != nil:
		orcFilter := &orcv1alpha1.SecurityGroupFilter{
			Name:                toOpenStackName(param.Filter.Name),
			Description:         toNeutronDescription(param.Filter.Description),
			FilterByNeutronTags: convertNeutronTags(param.Filter.FilterByNeutronTags),
		}
		obj.Spec.Import = &orcv1alpha1.SecurityGroupImport{Filter: orcFilter}
	}

	return obj
}

// buildVolumeType builds an unmanaged ORC VolumeType for importing an
// existing Cinder volume type by name.
func buildVolumeType(serverName, namespace, volumeTypeName string, credRef orcv1alpha1.CloudCredentialsReference) *orcv1alpha1.VolumeType {
	return &orcv1alpha1.VolumeType{
		ObjectMeta: metav1.ObjectMeta{
			Name:      VolumeTypeORCName(serverName, volumeTypeName),
			Namespace: namespace,
		},
		Spec: orcv1alpha1.VolumeTypeSpec{
			ManagementPolicy: orcv1alpha1.ManagementPolicyUnmanaged,
			Import: &orcv1alpha1.VolumeTypeImport{
				Filter: &orcv1alpha1.VolumeTypeFilter{
					Name: toOpenStackName(volumeTypeName),
				},
			},
			CloudCredentialsRef: credRef,
		},
	}
}
