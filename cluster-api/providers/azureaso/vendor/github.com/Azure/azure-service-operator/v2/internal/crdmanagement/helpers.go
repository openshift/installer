// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package crdmanagement

import (
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/rotisserie/eris"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"

	"github.com/Azure/azure-service-operator/v2/internal/set"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/registration"
)

const CRDFilePrefix = "apiextensions.k8s.io_v1_customresourcedefinition_"

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
			return nil, eris.Wrapf(err, "creating GVK for obj %T", storageType.Obj)
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
	knownTypes []*registration.KnownType,
) ([]*registration.KnownType, error) {
	// include map key is by CRD name, but we need it to be by kind
	includeKinds := set.Make[schema.GroupKind]()
	for _, crd := range include {
		includeKinds.Add(schema.GroupKind{Group: crd.Spec.Group, Kind: crd.Spec.Names.Kind})
	}

	result := make([]*registration.KnownType, 0, len(knownTypes))
	for _, knownType := range knownTypes {
		// Use the provided GVK to construct a new runtime object of the desired concrete type.
		gvk, err := apiutil.GVKForObject(knownType.Obj, scheme)
		if err != nil {
			return nil, eris.Wrapf(err, "creating GVK for obj %T", knownType)
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

// groupFromFilename extracts the API group from a CRD filename.
// CRD files follow the convention: "apiextensions.k8s.io_v1_customresourcedefinition_{kindPlural}.{group}.yaml".
// For example, "apiextensions.k8s.io_v1_customresourcedefinition_virtualnetworks.network.azure.com.yaml"
// returns "network.azure.com".
func groupFromFilename(filename string) (string, error) {
	crdName, err := crdNameFromFilename(filename)
	if err != nil {
		return "", err
	}

	// CRD name is "{kindPlural}.{group}", extract group (everything after the first dot)
	idx := strings.Index(crdName, ".")
	if idx < 0 || idx == len(crdName)-1 {
		return "", eris.Errorf("CRD name %q derived from filename %q has no group component", crdName, filename)
	}

	return crdName[idx+1:], nil
}

// crdNameFromFilename derives the CRD metadata.name from a CRD filename.
// CRD files follow the convention: "apiextensions.k8s.io_v1_customresourcedefinition_{crdname}.yaml",
// where {crdname} is the CRD's metadata.name (e.g. "virtualnetworks.network.azure.com").
// For example, "apiextensions.k8s.io_v1_customresourcedefinition_virtualnetworks.network.azure.com.yaml"
// returns "virtualnetworks.network.azure.com".
func crdNameFromFilename(filename string) (string, error) {
	if !strings.HasPrefix(filename, CRDFilePrefix) {
		return "", eris.Errorf("filename %q does not have expected prefix %q", filename, CRDFilePrefix)
	}

	if !strings.HasSuffix(filename, ".yaml") {
		return "", eris.Errorf("filename %q does not have .yaml extension", filename)
	}

	// Strip prefix and .yaml suffix
	crdName := strings.TrimPrefix(filename, CRDFilePrefix)
	crdName = strings.TrimSuffix(crdName, ".yaml")

	if crdName == "" {
		return "", eris.Errorf("filename %q has no CRD name between prefix and extension", filename)
	}

	return crdName, nil
}
