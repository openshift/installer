/*
Copyright 2024 The Kubernetes Authors.

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
	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/api/v1alpha1"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
)

type orcIdentityRefProvider struct {
	orcv1alpha1.CloudCredentialsRefProvider
}

func (o orcIdentityRefProvider) GetIdentityRef() (*string, *infrav1.OpenStackIdentityReference) {
	namespace, openStackCredentialsRef := o.GetCloudCredentialsRef()
	if namespace == nil || openStackCredentialsRef == nil {
		return nil, nil
	}

	return namespace, &infrav1.OpenStackIdentityReference{
		Name:      openStackCredentialsRef.SecretName,
		CloudName: openStackCredentialsRef.CloudName,
	}
}

var _ infrav1.IdentityRefProvider = orcIdentityRefProvider{}

func IdentityRefProvider(o orcv1alpha1.CloudCredentialsRefProvider) infrav1.IdentityRefProvider {
	return orcIdentityRefProvider{o}
}
