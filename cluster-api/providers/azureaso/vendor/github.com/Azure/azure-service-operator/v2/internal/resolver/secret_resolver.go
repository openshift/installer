/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package resolver

import (
	"context"

	"github.com/rotisserie/eris"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	"github.com/Azure/azure-service-operator/v2/internal/set"
	"github.com/Azure/azure-service-operator/v2/internal/util/kubeclient"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/core"
)

// SecretResolver is a secret resolver
type SecretResolver interface {
	ResolveSecretReference(ctx context.Context, ref genruntime.NamespacedSecretReference) (string, error)
	ResolveSecretReferences(ctx context.Context, refs set.Set[genruntime.NamespacedSecretReference]) (genruntime.Resolved[genruntime.SecretReference, string], error)
}

// kubeSecretResolver resolves Kubernetes secrets
type kubeSecretResolver struct {
	client kubeclient.Client
}

var _ SecretResolver = &kubeSecretResolver{}

func NewKubeSecretResolver(client kubeclient.Client) SecretResolver {
	return &kubeSecretResolver{
		client: client,
	}
}

// ResolveSecretReference resolves the secret reference and returns the corresponding secret value, or an error
// if it could not be found
func (r *kubeSecretResolver) ResolveSecretReference(ctx context.Context, ref genruntime.NamespacedSecretReference) (string, error) {
	refNamespacedName := types.NamespacedName{
		Namespace: ref.Namespace,
		Name:      ref.Name,
	}

	secret := &v1.Secret{}
	err := r.client.Get(ctx, refNamespacedName, secret)
	if err != nil {
		if apierrors.IsNotFound(err) {
			err := core.NewSecretNotFoundError(refNamespacedName, err)
			return "", eris.Wrapf(
				err,
				"couldn't resolve secret reference %s/%s.%s",
				ref.Namespace,
				ref.Name,
				ref.Key)
		}
	}

	// TODO: Do we want to confirm that the type is Opaque?

	valueBytes, ok := secret.Data[ref.Key]
	if !ok {
		return "", core.NewSecretNotFoundError(refNamespacedName, eris.Errorf("Secret %q does not contain key %q", refNamespacedName.String(), ref.Key))
	}

	return string(valueBytes), nil
}

// ResolveSecretReferences resolves all provided secret references
func (r *kubeSecretResolver) ResolveSecretReferences(ctx context.Context, refs set.Set[genruntime.NamespacedSecretReference]) (genruntime.Resolved[genruntime.SecretReference, string], error) {
	result := make(map[genruntime.SecretReference]string, len(refs))

	for ref := range refs {
		value, err := r.ResolveSecretReference(ctx, ref)
		if err != nil {
			return genruntime.MakeResolved[genruntime.SecretReference, string](nil), err
		}
		result[ref.SecretReference] = value
	}

	return genruntime.MakeResolved[genruntime.SecretReference, string](result), nil
}
