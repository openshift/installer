/*
Copyright 2021 The Kubernetes Authors.

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

// Package identity contains tools for managing secrets used to access the VCenter API.
package identity

import (
	"context"
	"errors"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
)

const (
	// UsernameKey is the key used for the username.
	UsernameKey = "username"
	// PasswordKey is the key used for the password.
	PasswordKey = "password"
)

// Credentials are the user credentials used with the VSphere API.
type Credentials struct {
	Username string
	Password string
}

// GetCredentials returns the VCenter credentials for the VSphereCluster.
func GetCredentials(ctx context.Context, c client.Client, cluster *infrav1.VSphereCluster, controllerNamespace string) (*Credentials, error) {
	if err := validateInputs(c, cluster); err != nil {
		return nil, err
	}

	ref := cluster.Spec.IdentityRef
	secret := &corev1.Secret{}
	var secretKey client.ObjectKey

	switch ref.Kind {
	case infrav1.SecretKind:
		secretKey = client.ObjectKey{
			Namespace: cluster.Namespace,
			Name:      ref.Name,
		}
	case infrav1.VSphereClusterIdentityKind:
		identity := &infrav1.VSphereClusterIdentity{}
		key := client.ObjectKey{
			Name: ref.Name,
		}
		if err := c.Get(ctx, key, identity); err != nil {
			return nil, err
		}

		if !identity.Status.Ready {
			return nil, errors.New("identity isn't ready to be used yet")
		}

		if identity.Spec.AllowedNamespaces == nil {
			return nil, errors.New("allowedNamespaces set to nil, no namespaces are allowed to use this identity")
		}

		selector, err := metav1.LabelSelectorAsSelector(&identity.Spec.AllowedNamespaces.Selector)
		if err != nil {
			return nil, errors.New("failed to build selector")
		}

		ns := &corev1.Namespace{}
		nsKey := client.ObjectKey{
			Name: cluster.Namespace,
		}
		if err := c.Get(ctx, nsKey, ns); err != nil {
			return nil, err
		}
		if !selector.Matches(labels.Set(ns.GetLabels())) {
			return nil, fmt.Errorf("namespace %s is not allowed to use specifified identity", cluster.Namespace)
		}

		secretKey = client.ObjectKey{
			Name:      identity.Spec.SecretName,
			Namespace: controllerNamespace,
		}
	default:
		return nil, fmt.Errorf("unknown type %s used for Identity", ref.Kind)
	}

	if err := c.Get(ctx, secretKey, secret); err != nil {
		return nil, err
	}

	credentials := &Credentials{
		Username: getData(secret, UsernameKey),
		Password: getData(secret, PasswordKey),
	}

	return credentials, nil
}

func validateInputs(c client.Client, cluster *infrav1.VSphereCluster) error {
	if c == nil {
		return errors.New("kubernetes client is required")
	}
	if cluster == nil {
		return errors.New("vsphere cluster is required")
	}
	ref := cluster.Spec.IdentityRef
	if ref == nil {
		return errors.New("IdentityRef is required")
	}
	return nil
}

// IsSecretIdentity returns true if the VSphereCluster identity is a Secret.
func IsSecretIdentity(cluster *infrav1.VSphereCluster) bool {
	if cluster == nil || cluster.Spec.IdentityRef == nil {
		return false
	}
	return cluster.Spec.IdentityRef.Kind == infrav1.SecretKind
}

// IsOwnedByIdentityOrCluster discovers if a secret is owned by a VSphereCluster or VSphereClusterIdentity.
func IsOwnedByIdentityOrCluster(ownerReferences []metav1.OwnerReference) bool {
	if len(ownerReferences) > 0 {
		for _, ownerReference := range ownerReferences {
			if !strings.Contains(ownerReference.APIVersion, infrav1.GroupName+"/") {
				continue
			}
			if ownerReference.Kind == "VSphereCluster" || ownerReference.Kind == "VSphereClusterIdentity" {
				return true
			}
		}
	}
	return false
}

func getData(secret *corev1.Secret, key string) string {
	if secret.Data == nil {
		return ""
	}
	if val, ok := secret.Data[key]; ok {
		return string(val)
	}
	return ""
}
