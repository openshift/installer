package controllers

import (
	"context"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	eksbootstrapv1 "sigs.k8s.io/cluster-api-provider-aws/v2/bootstrap/eks/api/v1beta2"
)

// FileResolver provides methods to resolve files and their content from secrets.
type FileResolver struct {
	Client client.Reader
}

// ResolveFiles resolves the content of files, fetching data from referenced secrets if needed.
func (fr *FileResolver) ResolveFiles(ctx context.Context, namespace string, files []eksbootstrapv1.File) ([]eksbootstrapv1.File, error) {
	collected := make([]eksbootstrapv1.File, 0, len(files))

	for i := range files {
		in := files[i]
		if in.ContentFrom != nil {
			data, err := fr.ResolveSecretFileContent(ctx, namespace, in)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to resolve file source")
			}
			in.ContentFrom = nil
			in.Content = string(data)
		}
		collected = append(collected, in)
	}

	return collected, nil
}

// ResolveSecretFileContent fetches the content of a file from a referenced secret.
func (fr *FileResolver) ResolveSecretFileContent(ctx context.Context, ns string, source eksbootstrapv1.File) ([]byte, error) {
	secret := &corev1.Secret{}
	key := types.NamespacedName{Namespace: ns, Name: source.ContentFrom.Secret.Name}
	if err := fr.Client.Get(ctx, key, secret); err != nil {
		if apierrors.IsNotFound(err) {
			return nil, errors.Wrapf(err, "secret not found: %s", key)
		}
		return nil, errors.Wrapf(err, "failed to retrieve Secret %q", key)
	}
	data, ok := secret.Data[source.ContentFrom.Secret.Key]
	if !ok {
		return nil, errors.Errorf("secret references non-existent secret key: %q", source.ContentFrom.Secret.Key)
	}
	return data, nil
}
