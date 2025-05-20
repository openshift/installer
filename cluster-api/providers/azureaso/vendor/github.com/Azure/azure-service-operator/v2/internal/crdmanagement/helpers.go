// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package crdmanagement

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"

	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"

	"github.com/Azure/azure-service-operator/v2/internal/set"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/registration"
)

func MakeCRDMap(
	crds []apiextensions.CustomResourceDefinition,
) map[string]apiextensions.CustomResourceDefinition {
	// Build a map so lookup is faster
	result := make(map[string]apiextensions.CustomResourceDefinition, len(crds))
	for _, crd := range crds {
		result[crd.Name] = crd
	}
	return result
}

func FilterStorageTypesByReadyCRDs(
	logger logr.Logger,
	scheme *runtime.Scheme,
	include map[string]apiextensions.CustomResourceDefinition,
	storageTypes []*registration.StorageType,
) ([]*registration.StorageType, error) {
	// include map key is by CRD name, but we need it to be by kind
	includeKinds := set.Make[schema.GroupKind]()
	for _, crd := range include {
		includeKinds.Add(schema.GroupKind{Group: crd.Spec.Group, Kind: crd.Spec.Names.Kind})
	}

	result := make([]*registration.StorageType, 0, len(storageTypes))

	for _, storageType := range storageTypes {
		// Use the provided GVK to construct a new runtime object of the desired concrete type.
		gvk, err := apiutil.GVKForObject(storageType.Obj, scheme)
		if err != nil {
			return nil, errors.Wrapf(err, "creating GVK for obj %T", storageType.Obj)
		}

		if !includeKinds.Contains(gvk.GroupKind()) {
			logger.V(0).Info(
				"Skipping reconciliation of resource because CRD was not installed or did not match the expected shape",
				"groupKind", gvk.GroupKind().String())
			continue
		}

		result = append(result, storageType)
	}

	return result, nil
}

func FilterKnownTypesByReadyCRDs(
	logger logr.Logger,
	scheme *runtime.Scheme,
	include map[string]apiextensions.CustomResourceDefinition,
	knownTypes []client.Object,
) ([]client.Object, error) {
	// include map key is by CRD name, but we need it to be by kind
	includeKinds := set.Make[schema.GroupKind]()
	for _, crd := range include {
		includeKinds.Add(schema.GroupKind{Group: crd.Spec.Group, Kind: crd.Spec.Names.Kind})
	}

	result := make([]client.Object, 0, len(knownTypes))
	for _, knownType := range knownTypes {
		// Use the provided GVK to construct a new runtime object of the desired concrete type.
		gvk, err := apiutil.GVKForObject(knownType, scheme)
		if err != nil {
			return nil, errors.Wrapf(err, "creating GVK for obj %T", knownType)
		}
		if !includeKinds.Contains(gvk.GroupKind()) {
			logger.V(0).Info(
				"Skipping webhooks of resource because CRD was not installed or did not match the expected shape",
				"groupKind", gvk.GroupKind().String())
			continue
		}

		result = append(result, knownType)
	}

	return result, nil
}

func makeMatchString(crd apiextensions.CustomResourceDefinition) string {
	group := crd.Spec.Group
	kind := crd.Spec.Names.Kind

	// matchString should be "group/kind"
	return fmt.Sprintf("%s/%s", group, kind)
}
