/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package resolver

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"

	"github.com/Azure/azure-service-operator/v2/internal/reflecthelpers"
	"github.com/Azure/azure-service-operator/v2/internal/set"
	"github.com/Azure/azure-service-operator/v2/internal/util/kubeclient"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/core"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/registration"
)

type Resolver struct {
	client                   kubeclient.Client
	kubeSecretResolver       SecretResolver
	kubeSecretMapResolver    SecretMapResolver
	kubeConfigMapResolver    ConfigMapResolver
	reconciledResourceLookup map[schema.GroupKind]schema.GroupVersionKind
}

func NewResolver(client kubeclient.Client) *Resolver {
	return &Resolver{
		client:                   client,
		kubeSecretResolver:       NewKubeSecretResolver(client),
		kubeSecretMapResolver:    NewKubeSecretMapResolver(client),
		kubeConfigMapResolver:    NewKubeConfigMapResolver(client),
		reconciledResourceLookup: make(map[schema.GroupKind]schema.GroupVersionKind),
	}
}

func (r *Resolver) IndexStorageTypes(scheme *runtime.Scheme, objs []*registration.StorageType) error {
	for _, obj := range objs {
		gvk, err := apiutil.GVKForObject(obj.Obj, scheme)
		if err != nil {
			return errors.Wrapf(err, "creating GVK for obj %T", obj)
		}
		groupKind := schema.GroupKind{Group: gvk.Group, Kind: gvk.Kind}
		if existing, ok := r.reconciledResourceLookup[groupKind]; ok {
			if existing == gvk {
				continue
			}

			return errors.Errorf(
				"group: %q, kind: %q already has registered storage version %q, but found %q as well",
				gvk.Group,
				gvk.Kind,
				existing.Version,
				gvk.Version)
		}
		r.reconciledResourceLookup[groupKind] = gvk
	}

	return nil
}

// ResolveReferenceToARMID gets a references ARM ID. If the reference is just pointing to an ARM resource then the ARMID is returned.
// If the reference is pointing to a Kubernetes resource, that resource is looked up and its ARM ID is computed.
func (r *Resolver) ResolveReferenceToARMID(ctx context.Context, ref genruntime.NamespacedResourceReference) (string, error) {
	if ref.IsDirectARMReference() {
		return ref.ARMID, nil
	}

	obj, err := r.ResolveReference(ctx, ref)
	if err != nil {
		return "", err
	}

	// There are two ways to get the ARM ID here, we can look it up using GetResourceID, which will only work if the
	// resource has actually been successfully deployed to Azure, or we can "compute" it. Currently it's harder to compute
	// it given that a resource doesn't know what subscription it's deployed in... but we should probably change that
	// and move to computing it here.
	id, ok := genruntime.GetResourceID(obj)
	if !ok {
		// Resource doesn't have a resource ID. This probably means it's not done deploying
		return "", errors.Errorf("ref %s doesn't have an assigned ARM ID", ref)
	}

	return id, nil
}

// ResolveReferencesToARMIDs resolves all provided references to their ARM IDs.
func (r *Resolver) ResolveReferencesToARMIDs(ctx context.Context, refs map[genruntime.NamespacedResourceReference]struct{}) (genruntime.Resolved[genruntime.ResourceReference, string], error) {
	result := make(map[genruntime.ResourceReference]string, len(refs))

	for ref := range refs {
		armID, err := r.ResolveReferenceToARMID(ctx, ref)
		if err != nil {
			return genruntime.MakeResolved[genruntime.ResourceReference, string](nil), err
		}
		result[ref.ResourceReference] = armID
	}

	return genruntime.MakeResolved[genruntime.ResourceReference, string](result), nil
}

// ResolveResourceReferences resolves every reference found on the specified genruntime.ARMMetaObject to its corresponding ARM ID.
func (r *Resolver) ResolveResourceReferences(ctx context.Context, metaObject genruntime.ARMMetaObject) (genruntime.Resolved[genruntime.ResourceReference, string], error) {
	refs, err := reflecthelpers.FindResourceReferences(metaObject)
	if err != nil {
		return genruntime.Resolved[genruntime.ResourceReference, string]{}, errors.Wrapf(err, "finding references on %q", metaObject.GetName())
	}

	// Include the namespace
	namespacedRefs := make(map[genruntime.NamespacedResourceReference]struct{}, len(refs))
	for ref := range refs {
		namespacedRefs[ref.AsNamespacedRef(metaObject.GetNamespace())] = struct{}{}
	}

	// resolve them
	resolvedRefs, err := r.ResolveReferencesToARMIDs(ctx, namespacedRefs)
	if err != nil {
		return genruntime.Resolved[genruntime.ResourceReference, string]{}, errors.Wrapf(err, "failed resolving ARM IDs for references")
	}

	return resolvedRefs, nil
}

// ResolveResourceHierarchy gets the resource hierarchy for a given resource. The result is a slice of
// resources, with the uppermost parent at position 0 and the resource itself at position len(slice)-1.
// Note that there is NO GUARANTEE that this hierarchy is "complete". It may root up to a resource which uses
// the ARMID field of owner.
func (r *Resolver) ResolveResourceHierarchy(ctx context.Context, obj genruntime.ARMMetaObject) (ResourceHierarchy, error) {
	owner := obj.Owner()
	if owner == nil {
		return ResourceHierarchy{obj}, nil
	}

	ownerDetails, err := r.ResolveOwner(ctx, obj)
	if err != nil {
		return nil, err
	}

	if ownerDetails.Result == OwnerFoundARM {
		return ResourceHierarchy{obj}, nil
	}

	owners, err := r.ResolveResourceHierarchy(ctx, ownerDetails.Owner)
	if err != nil {
		return nil, errors.Wrapf(err, "getting owners for %s", ownerDetails.Owner.GetName())
	}

	return append(owners, obj), nil
}

// ResolveReference resolves a reference, or returns an error if the reference is not pointing to a KubernetesResource
func (r *Resolver) ResolveReference(ctx context.Context, ref genruntime.NamespacedResourceReference) (genruntime.ARMMetaObject, error) {
	refGVK, err := r.findGVK(ref)
	if err != nil {
		return nil, err
	}

	refNamespacedName := types.NamespacedName{
		Namespace: ref.Namespace,
		Name:      ref.Name,
	}

	refObj, err := r.client.GetObject(ctx, refNamespacedName, refGVK)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// Check if the user has mistakenly put the armID in 'name' field
			_, err1 := arm.ParseResourceID(ref.Name)
			if err1 == nil {
				return nil, errors.Errorf("couldn't resolve reference %s. 'name' looks like it might be an ARM ID; did you mean 'armID: %s'?", refNamespacedName.String(), ref.Name)
			}
			err := core.NewReferenceNotFoundError(refNamespacedName, err)
			return nil, errors.WithStack(err)
		}

		return nil, errors.Wrapf(err, "couldn't resolve reference %s", ref.String())
	}

	metaObj, ok := refObj.(genruntime.ARMMetaObject)
	if !ok {
		return nil, errors.Errorf("reference %s (%s) was not of type genruntime.ARMMetaObject", refNamespacedName, refGVK)
	}

	return metaObj, nil
}

type ResolveOwnerResult string

const (
	// OwnerFoundKubernetes indicates the owner was found in Kubernetes.
	OwnerFoundKubernetes = ResolveOwnerResult("OwnerFoundKubernetes")
	// OwnerFoundARM indicates the owner is an ARM ID. The resource the ARM ID points to may or may not exist in Azure currently.
	OwnerFoundARM = ResolveOwnerResult("OwnerFoundARM")
	// OwnerNotExpected indicates that this resource is not expected to have any owner. (Example: ResourceGroup)
	OwnerNotExpected = ResolveOwnerResult("OwnerNotExpected")
)

type OwnerDetails struct {
	Result ResolveOwnerResult
	Owner  genruntime.ARMMetaObject
	ARMID  string
}

func (det OwnerDetails) FoundKubernetesOwner() bool {
	if det.Result == OwnerNotExpected || det.Result == OwnerFoundARM {
		// If no owner is expected or the owner is only in ARM, no need to assign ownership in Kubernetes
		return false
	}

	return true
}

func OwnerDetailsFromKubernetes(owner genruntime.ARMMetaObject) OwnerDetails {
	return OwnerDetails{
		Result: OwnerFoundKubernetes,
		Owner:  owner,
	}
}

func OwnerDetailsFromARM(armID string) OwnerDetails {
	return OwnerDetails{
		Result: OwnerFoundARM,
		ARMID:  armID,
	}
}

func OwnerDetailsNotExpected() OwnerDetails {
	return OwnerDetails{
		Result: OwnerNotExpected,
	}
}

// ResolveOwner returns an OwnerDetails describing more information about the owner of the provided resource.
// If the resource is supposed to have
// an owner but doesn't, this returns an ReferenceNotFound error. If the resource is not supposed
// to have an owner (for example, ResourceGroup) or the owner points to a raw ARM ID this returns an OwnerDetails
// with the OwnerDetails.Result set appropriately.
func (r *Resolver) ResolveOwner(ctx context.Context, obj genruntime.ARMOwnedMetaObject) (OwnerDetails, error) {
	owner := obj.Owner()

	if owner == nil {
		return OwnerDetailsNotExpected(), nil
	}

	if owner.IsDirectARMReference() {
		return OwnerDetailsFromARM(owner.ARMID), nil
	}

	namespacedRef := genruntime.NamespacedResourceReference{
		ResourceReference: *owner,
		Namespace:         obj.GetNamespace(),
	}
	ownerMeta, err := r.ResolveReference(ctx, namespacedRef)
	if err != nil {
		return OwnerDetails{}, err
	}

	return OwnerDetailsFromKubernetes(ownerMeta), nil
}

// Scheme returns the current scheme from our client
func (r *Resolver) Scheme() *runtime.Scheme {
	return r.client.Scheme()
}

func (r *Resolver) findGVK(ref genruntime.NamespacedResourceReference) (schema.GroupVersionKind, error) {
	var ownerGvk schema.GroupVersionKind

	if !ref.IsKubernetesReference() {
		return ownerGvk, errors.Errorf("reference %s is not pointing to a Kubernetes resource", ref)
	}

	groupKind := schema.GroupKind{Group: ref.Group, Kind: ref.Kind}
	gvk, ok := r.reconciledResourceLookup[groupKind]
	if !ok {
		return ownerGvk, errors.Errorf("group: %q, kind: %q was not in reconciledResourceLookup", ref.Group, ref.Kind)
	}

	return gvk, nil
}

// ResolveSecretReferences resolves all provided secret references
func (r *Resolver) ResolveSecretReferences(
	ctx context.Context,
	refs set.Set[genruntime.NamespacedSecretReference],
) (genruntime.Resolved[genruntime.SecretReference, string], error) {
	return r.kubeSecretResolver.ResolveSecretReferences(ctx, refs)
}

// ResolveResourceSecretReferences resolves all of the specified genruntime.MetaObject's secret references.
func (r *Resolver) ResolveResourceSecretReferences(ctx context.Context, metaObject genruntime.MetaObject) (genruntime.Resolved[genruntime.SecretReference, string], error) {
	refs, err := reflecthelpers.FindSecretReferences(metaObject)
	if err != nil {
		return genruntime.Resolved[genruntime.SecretReference, string]{}, errors.Wrapf(err, "finding secrets on %q", metaObject.GetName())
	}

	// Include the namespace
	namespacedSecretRefs := set.Make[genruntime.NamespacedSecretReference]()
	for ref := range refs {
		namespacedSecretRefs.Add(ref.AsNamespacedRef(metaObject.GetNamespace()))
	}

	// resolve them
	resolvedSecrets, err := r.ResolveSecretReferences(ctx, namespacedSecretRefs)
	if err != nil {
		return genruntime.Resolved[genruntime.SecretReference, string]{}, errors.Wrapf(err, "failed resolving secret references")
	}

	return resolvedSecrets, nil
}

// ResolveSecretMapReferences resolves all provided secret map references
func (r *Resolver) ResolveSecretMapReferences(
	ctx context.Context,
	refs set.Set[genruntime.NamespacedSecretMapReference],
) (genruntime.Resolved[genruntime.SecretMapReference, map[string]string], error) {
	return r.kubeSecretMapResolver.ResolveSecretMapReferences(ctx, refs)
}

// ResolveResourceSecretMapReferences resolves all the specified genruntime.MetaObject's secret maps.
func (r *Resolver) ResolveResourceSecretMapReferences(
	ctx context.Context,
	metaObject genruntime.MetaObject,
) (genruntime.Resolved[genruntime.SecretMapReference, map[string]string], error) {
	refs, err := reflecthelpers.FindSecretMaps(metaObject)
	if err != nil {
		return genruntime.Resolved[genruntime.SecretMapReference, map[string]string]{}, errors.Wrapf(err, "finding secrets on %q", metaObject.GetName())
	}

	// Include the namespace
	namespacedSecretRefs := set.Make[genruntime.NamespacedSecretMapReference]()
	for ref := range refs {
		namespacedSecretRefs.Add(ref.AsNamespacedRef(metaObject.GetNamespace()))
	}

	// resolve them
	resolvedSecrets, err := r.ResolveSecretMapReferences(ctx, namespacedSecretRefs)
	if err != nil {
		return genruntime.Resolved[genruntime.SecretMapReference, map[string]string]{}, errors.Wrapf(err, "failed resolving secret references")
	}

	return resolvedSecrets, nil
}

// ResolveConfigMapReferences resolves all provided secret references
func (r *Resolver) ResolveConfigMapReferences(
	ctx context.Context,
	refs set.Set[genruntime.NamespacedConfigMapReference],
) (genruntime.Resolved[genruntime.ConfigMapReference, string], error) {
	return r.kubeConfigMapResolver.ResolveConfigMapReferences(ctx, refs)
}

// ResolveResourceConfigMapReferences resolves the specified genruntime.MetaObject's configmap references.
func (r *Resolver) ResolveResourceConfigMapReferences(ctx context.Context, metaObject genruntime.MetaObject) (genruntime.Resolved[genruntime.ConfigMapReference, string], error) {
	refs, err := reflecthelpers.FindConfigMapReferences(metaObject)
	if err != nil {
		return genruntime.Resolved[genruntime.ConfigMapReference, string]{}, errors.Wrapf(err, "finding config maps on %q", metaObject.GetName())
	}

	// Include the namespace
	namespacedConfigMapReferences := set.Make[genruntime.NamespacedConfigMapReference]()
	for ref := range refs {
		namespacedConfigMapReferences.Add(ref.AsNamespacedRef(metaObject.GetNamespace()))
	}

	// resolve them
	resolvedConfigMaps, err := r.ResolveConfigMapReferences(ctx, namespacedConfigMapReferences)
	if err != nil {
		return genruntime.Resolved[genruntime.ConfigMapReference, string]{}, errors.Wrapf(err, "failed resolving config map references")
	}

	return resolvedConfigMaps, nil
}

// ResolveAll resolves every reference on the provided genruntime.ARMMetaObject.
// This includes: owner, all resource references, and all secrets.
func (r *Resolver) ResolveAll(ctx context.Context, metaObject genruntime.ARMMetaObject) (ResourceHierarchy, genruntime.ConvertToARMResolvedDetails, error) {
	// Resolve the resource hierarchy (owner)
	resourceHierarchy, err := r.ResolveResourceHierarchy(ctx, metaObject)
	if err != nil {
		return nil, genruntime.ConvertToARMResolvedDetails{}, err
	}

	// Resolve all ARM ID references
	resolvedRefs, err := r.ResolveResourceReferences(ctx, metaObject)
	if err != nil {
		return nil, genruntime.ConvertToARMResolvedDetails{}, err
	}

	// Resolve all secrets
	resolvedSecrets, err := r.ResolveResourceSecretReferences(ctx, metaObject)
	if err != nil {
		return nil, genruntime.ConvertToARMResolvedDetails{}, err
	}

	resolvedSecretMaps, err := r.ResolveResourceSecretMapReferences(ctx, metaObject)
	if err != nil {
		return nil, genruntime.ConvertToARMResolvedDetails{}, err
	}

	// Resolve all configmaps
	resolvedConfigMaps, err := r.ResolveResourceConfigMapReferences(ctx, metaObject)
	if err != nil {
		return nil, genruntime.ConvertToARMResolvedDetails{}, err
	}

	resolvedDetails := genruntime.ConvertToARMResolvedDetails{
		Name:               resourceHierarchy.AzureName(),
		ResolvedReferences: resolvedRefs,
		ResolvedSecrets:    resolvedSecrets,
		ResolvedSecretMaps: resolvedSecretMaps,
		ResolvedConfigMaps: resolvedConfigMaps,
	}

	return resourceHierarchy, resolvedDetails, nil
}
