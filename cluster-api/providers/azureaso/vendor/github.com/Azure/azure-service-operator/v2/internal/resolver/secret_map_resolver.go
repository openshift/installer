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

// SecretMapResolver is a secret collection resolver
type SecretMapResolver interface {
	// ResolveSecretMapReference takes a secret reference and returns the actual secret values.
	ResolveSecretMapReference(
		ctx context.Context,
		ref genruntime.NamespacedSecretMapReference,
	) (map[string]string, error)

	// ResolveSecretMapReferences takes a collection of secret references and returns a map of ref -> value
	// for each secret reference.
	ResolveSecretMapReferences(
		ctx context.Context,
		refs set.Set[genruntime.NamespacedSecretMapReference],
	) (genruntime.Resolved[genruntime.SecretMapReference, map[string]string], error)
}

// kubeSecretResolver resolves Kubernetes secrets
type kubeSecretMapResolver struct {
	client kubeclient.Client
}

var _ SecretMapResolver = &kubeSecretMapResolver{}

func NewKubeSecretMapResolver(client kubeclient.Client) SecretMapResolver {
	return &kubeSecretMapResolver{
		client: client,
	}
}

// ResolveSecretMapReference resolves a secret map reference and returns the corresponding secret values, or an error
// if it could not be found
func (r *kubeSecretMapResolver) ResolveSecretMapReference(
	ctx context.Context,
	ref genruntime.NamespacedSecretMapReference,
) (map[string]string, error) {
	refNamespacedName := types.NamespacedName{
		Namespace: ref.Namespace,
		Name:      ref.Name,
	}

	secret := &v1.Secret{}
	err := r.client.Get(ctx, refNamespacedName, secret)
	if err != nil {
		if apierrors.IsNotFound(err) {
			err := core.NewSecretNotFoundError(refNamespacedName, err)
			return nil, eris.Wrapf(
				err,
				"couldn't resolve secret collection %s/%s",
				ref.Namespace,
				ref.Name)
		}

		return nil, eris.Wrapf(err, "couldn't resolve secret collection %s", ref.String())
	}

	// TODO: Do we want to confirm that the type is Opaque?

	result := make(map[string]string, len(secret.Data))
	for k, v := range secret.Data {
		result[k] = string(v)
	}

	return result, nil
}

// ResolveSecretMapReferences resolves all provided secret map references
func (r *kubeSecretMapResolver) ResolveSecretMapReferences(
	ctx context.Context,
	refs set.Set[genruntime.NamespacedSecretMapReference],
) (genruntime.Resolved[genruntime.SecretMapReference, map[string]string], error) {
	result := make(map[genruntime.SecretMapReference]map[string]string, len(refs))

	for ref := range refs {
		value, err := r.ResolveSecretMapReference(ctx, ref)
		if err != nil {
			return genruntime.MakeResolved[genruntime.SecretMapReference, map[string]string](nil), err
		}

		result[ref.SecretMapReference] = value
	}

	return genruntime.MakeResolved(result), nil
}
