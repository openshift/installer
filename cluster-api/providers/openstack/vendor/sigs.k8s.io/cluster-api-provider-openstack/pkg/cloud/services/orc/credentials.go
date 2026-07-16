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
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"

	infrav1alpha1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha1"
	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta2"
)

// ResolveCloudCredentialsRef resolves a CAPO IdentityRef to an ORC
// CloudCredentialsReference.
//
// For type "Secret" (or empty/default), the mapping is direct: the secret
// name and cloud name are taken from the IdentityRef. The secret must be
// in the same namespace as the OpenStackServer.
//
// For type "ClusterIdentity", the controller fetches the cluster-scoped
// OpenStackClusterIdentity resource and resolves its SecretRef. The secret
// must be in the same namespace as the ORC resources (the OpenStackServer's
// namespace). Cross-namespace secrets are not yet supported.
func ResolveCloudCredentialsRef(ctx context.Context, k8sClient client.Client, namespace string, identityRef infrav1.OpenStackIdentityReference) (orcv1alpha1.CloudCredentialsReference, error) {
	switch identityRef.Type {
	case "", "Secret":
		return orcv1alpha1.CloudCredentialsReference{
			SecretName: identityRef.Name,
			CloudName:  identityRef.CloudName,
		}, nil

	case "ClusterIdentity":
		identity := &infrav1alpha1.OpenStackClusterIdentity{}
		if err := k8sClient.Get(ctx, client.ObjectKey{Name: identityRef.Name}, identity); err != nil {
			return orcv1alpha1.CloudCredentialsReference{}, fmt.Errorf("failed to get OpenStackClusterIdentity %q: %w", identityRef.Name, err)
		}

		secretRef := identity.Spec.SecretRef
		if secretRef.Namespace != namespace {
			// TODO: Support cross-namespace secret copying for ClusterIdentity.
			// ORC requires the secret to be in the same namespace as the ORC
			// resource. A future iteration could copy the secret to the target
			// namespace or contribute cross-namespace secret support to ORC.
			return orcv1alpha1.CloudCredentialsReference{}, fmt.Errorf(
				"OpenStackClusterIdentity %q references secret %s/%s, but ORC resources are in namespace %q; "+
					"cross-namespace secrets are not yet supported",
				identityRef.Name, secretRef.Namespace, secretRef.Name, namespace,
			)
		}

		return orcv1alpha1.CloudCredentialsReference{
			SecretName: secretRef.Name,
			CloudName:  identityRef.CloudName,
		}, nil

	default:
		return orcv1alpha1.CloudCredentialsReference{}, fmt.Errorf("unsupported identity reference type %q", identityRef.Type)
	}
}
