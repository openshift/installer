/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package resolver

import (
	"context"

	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/core"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	"github.com/Azure/azure-service-operator/v2/internal/set"
	"github.com/Azure/azure-service-operator/v2/internal/util/kubeclient"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
)

// ConfigMapResolver is a configmap resolver
type ConfigMapResolver interface {
	ResolveConfigMapReference(ctx context.Context, ref genruntime.NamespacedConfigMapReference) (string, error)
	ResolveConfigMapReferences(ctx context.Context, refs set.Set[genruntime.NamespacedConfigMapReference]) (genruntime.Resolved[genruntime.ConfigMapReference], error)
}

// kubeConfigMapResolver resolves Kubernetes config maps
type kubeConfigMapResolver struct {
	client kubeclient.Client
}

var _ ConfigMapResolver = &kubeConfigMapResolver{}

func NewKubeConfigMapResolver(client kubeclient.Client) ConfigMapResolver {
	return &kubeConfigMapResolver{
		client: client,
	}
}

// ResolveConfigMapReference resolves the configmap reference and returns the corresponding value, or an error
// if it could not be found
func (r *kubeConfigMapResolver) ResolveConfigMapReference(ctx context.Context, ref genruntime.NamespacedConfigMapReference) (string, error) {
	refNamespacedName := types.NamespacedName{
		Namespace: ref.Namespace,
		Name:      ref.Name,
	}

	configMap := &v1.ConfigMap{}
	err := r.client.Get(ctx, refNamespacedName, configMap)
	if err != nil {
		if apierrors.IsNotFound(err) {
			err := core.NewConfigMapNotFoundError(refNamespacedName, err)
			return "", errors.WithStack(err)
		}

		return "", errors.Wrapf(err, "couldn't resolve config map reference %s", ref.String())
	}

	value, ok := configMap.Data[ref.Key] // TODO: Do we need to also check the binaryData field?
	if !ok {
		return "", core.NewConfigMapNotFoundError(refNamespacedName, errors.Errorf("ConfigMap %q does not contain key %q", refNamespacedName.String(), ref.Key))
	}

	return value, nil
}

// ResolveConfigMapReferences resolves all provided configmap references
func (r *kubeConfigMapResolver) ResolveConfigMapReferences(ctx context.Context, refs set.Set[genruntime.NamespacedConfigMapReference]) (genruntime.Resolved[genruntime.ConfigMapReference], error) {
	result := make(map[genruntime.ConfigMapReference]string, len(refs))

	for ref := range refs {
		value, err := r.ResolveConfigMapReference(ctx, ref)
		if err != nil {
			return genruntime.MakeResolved[genruntime.ConfigMapReference](nil), err
		}
		result[ref.ConfigMapReference] = value
	}

	return genruntime.MakeResolved[genruntime.ConfigMapReference](result), nil
}
