/*
Copyright 2023 The Kubernetes Authors.

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

package aso

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	asoannotations "github.com/Azure/azure-service-operator/v2/pkg/common/annotations"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/conditions"
	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/yaml"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/util/aso"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

const (
	// prePauseReconcilePolicyAnnotation is the annotation key for the value of
	// asoannotations.ReconcilePolicy that was set before pausing.
	prePauseReconcilePolicyAnnotation = "sigs.k8s.io/cluster-api-provider-azure-pre-pause-reconcile-policy"

	requeueInterval = 20 * time.Second

	createOrUpdateFutureType = "ASOCreateOrUpdate"
	deleteFutureType         = "ASODelete"
)

// reconciler is an implementation of the Reconciler interface. It handles creation
// and deletion of resources using ASO.
type reconciler[T genruntime.MetaObject] struct {
	client.Client

	clusterName string
	owner       client.Object
}

// New creates a new ASO reconciler.
func New[T genruntime.MetaObject](ctrlClient client.Client, clusterName string, owner client.Object) Reconciler[T] {
	return &reconciler[T]{
		Client:      ctrlClient,
		clusterName: clusterName,
		owner:       owner,
	}
}

// CreateOrUpdateResource implements the logic for creating a new or updating an
// existing resource with ASO.
func (r *reconciler[T]) CreateOrUpdateResource(ctx context.Context, spec azure.ASOResourceSpecGetter[T], serviceName string) (T, error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "services.aso.CreateOrUpdateResource")
	defer done()

	resource := spec.ResourceRef()
	resource.SetNamespace(r.owner.GetNamespace())
	resourceName := resource.GetName()
	resourceNamespace := resource.GetNamespace()

	log = log.WithValues("service", serviceName, "resource", resourceName, "namespace", resourceNamespace)

	var readyErr error
	var adopt bool
	var existing T
	var zero T // holds the zero value, to be returned with non-nil errors.
	resourceExists := false
	if err := r.Client.Get(ctx, client.ObjectKeyFromObject(resource), resource); err != nil {
		if !apierrors.IsNotFound(err) {
			return zero, errors.Wrapf(err, "failed to get existing resource %s/%s (service: %s)", resourceNamespace, resourceName, serviceName)
		}
		log.V(2).Info("existing resource not found, will create a new one")
	} else {
		existing = resource
		resourceExists = true
		log.V(2).Info("successfully got existing resource")

		if isOwned, err := isOwnedBy(resource, r.owner, r.Scheme()); err != nil {
			return zero, err
		} else if !isOwned && !hasLegacyOwnedByLabel(resource.GetLabels(), r.clusterName) {
			log.V(4).Info("skipping reconcile for unmanaged resource")
			return existing, nil
		}

		// Check if there is an ongoing long running operation.
		conds := existing.GetConditions()
		i, readyExists := conds.FindIndexByType(conditions.ConditionTypeReady)
		if !readyExists {
			return zero, azure.WithTransientError(errors.New("ready status unknown"), requeueInterval)
		}
		if cond := conds[i]; cond.Status != metav1.ConditionTrue {
			switch {
			case cond.Reason == conditions.ReasonAzureResourceNotFound.Name &&
				existing.GetAnnotations()[asoannotations.ReconcilePolicy] == string(asoannotations.ReconcilePolicySkip):
				// This resource was originally created by CAPZ and a
				// corresponding Azure resource has been found not to exist, so
				// CAPZ will tell ASO to adopt the resource by setting its
				// reconcile policy to "manage". This extra step is necessary to
				// handle user-managed resources that already exist in Azure and
				// should not be reconciled by ASO while ensuring they're still
				// represented in ASO.
				log.V(2).Info("resource not found in Azure and \"skip\" reconcile-policy set, adopting")
				// Don't set readyErr so the resource can be adopted with an
				// update instead of returning early.
				adopt = true
			case cond.Reason == conditions.ReasonReconciling.Name:
				readyErr = azure.NewOperationNotDoneError(&infrav1.Future{
					Type:          createOrUpdateFutureType,
					ResourceGroup: existing.GetNamespace(),
					Name:          existing.GetName(),
				})
			default:
				readyErr = fmt.Errorf("resource is not Ready: %s", conds[i].Message)
			}

			if readyErr != nil {
				if conds[i].Severity == conditions.ConditionSeverityError {
					readyErr = azure.WithTerminalError(readyErr)
				} else {
					readyErr = azure.WithTransientError(readyErr, requeueInterval)
				}
			}
		}
	}

	// Construct parameters using the resource spec and information from the existing resource, if there is one.
	var existingCopy T
	if resourceExists {
		existingCopy = existing.DeepCopyObject().(T)
	}
	parameters, err := PatchedParameters(ctx, r.Scheme(), spec, existingCopy)
	if err != nil {
		return zero, errors.Wrapf(err, "failed to get desired parameters for resource %s/%s (service: %s)", resourceNamespace, resourceName, serviceName)
	}

	parameters.SetName(resourceName)
	parameters.SetNamespace(resourceNamespace)

	if err := controllerutil.SetControllerReference(r.owner, parameters, r.Client.Scheme()); err != nil {
		return zero, errors.Wrap(err, "failed to set owner ref")
	}

	if t, ok := spec.(TagsGetterSetter[T]); ok {
		if err := reconcileTags(t, existing, resourceExists, parameters); err != nil {
			return zero, errors.Wrap(err, "failed to reconcile tags")
		}
	}

	labels := parameters.GetLabels()
	if labels == nil {
		labels = make(map[string]string)
	}
	labels[clusterv1.ClusterNameLabel] = r.clusterName

	annotations := parameters.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}

	if prevReconcilePolicy, ok := annotations[prePauseReconcilePolicyAnnotation]; ok {
		annotations[asoannotations.ReconcilePolicy] = prevReconcilePolicy
		delete(annotations, prePauseReconcilePolicyAnnotation)
	}
	if !resourceExists {
		// Create the ASO resource with "skip" in case a matching resource
		// already exists in Azure, in which case CAPZ will assume it is managed
		// by the user and ASO should not actively reconcile changes to the ASO
		// resource. In the canonical "entirely managed by CAPZ" case, the next
		// reconciliation will reveal the resource does not already exist in
		// Azure and the ASO resource will be adopted by changing this
		// annotation to "manage".
		annotations[asoannotations.ReconcilePolicy] = string(asoannotations.ReconcilePolicySkip)
	} else {
		adopt = adopt || spec.WasManaged(existing)
	}
	if adopt {
		annotations[asoannotations.ReconcilePolicy] = string(asoannotations.ReconcilePolicyManage)
	}

	// Set the secret name annotation in order to leverage the ASO resource credential scope as defined in
	// https://azure.github.io/azure-service-operator/guide/authentication/credential-scope/#resource-scope.
	annotations[asoannotations.PerResourceSecret] = aso.GetASOSecretName(r.clusterName)

	if len(labels) == 0 {
		labels = nil
	}
	parameters.SetLabels(labels)
	if len(annotations) == 0 {
		annotations = nil
	}
	parameters.SetAnnotations(annotations)

	diff := cmp.Diff(existing, parameters)
	if diff == "" {
		if readyErr != nil {
			// Only return this error when the resource is up to date in order to permit updates from
			// Parameters which may fix the resource's current state.
			return zero, readyErr
		}
		log.V(2).Info("resource up to date")
		return existing, nil
	}
	log.V(2).Info("creating or updating resource", "diff", diff)
	return r.createOrUpdateResource(ctx, existing, parameters, resourceExists, serviceName)
}

// PatchedParameters returns the Parameters of spec with patches applied.
func PatchedParameters[T genruntime.MetaObject](ctx context.Context, scheme *runtime.Scheme, spec azure.ASOResourceSpecGetter[T], existing T) (T, error) {
	var zero T // to be returned with non-nil errors
	parameters, err := spec.Parameters(ctx, existing)
	if err != nil {
		return zero, err
	}
	return applyPatches(scheme, spec, parameters)
}

func applyPatches[T genruntime.MetaObject](scheme *runtime.Scheme, spec azure.ASOResourceSpecGetter[T], parameters T) (T, error) {
	p, ok := spec.(Patcher)
	if !ok {
		return parameters, nil
	}

	var zero T // to be returned with non-nil errors

	gvk, err := apiutil.GVKForObject(parameters, scheme)
	if err != nil {
		return zero, errors.Wrap(err, "failed to get GroupVersionKind for object")
	}

	parameters.GetObjectKind().SetGroupVersionKind(gvk)
	paramData, err := json.Marshal(parameters)
	if err != nil {
		return zero, errors.Wrap(err, "failed to marshal JSON for patch")
	}

	for i, extraPatch := range p.ExtraPatches() {
		jsonPatch, err := yaml.ToJSON([]byte(extraPatch))
		if err != nil {
			return zero, errors.Wrapf(err, "failed to convert patch at index %d to JSON", i)
		}
		paramData, err = jsonpatch.MergePatch(paramData, jsonPatch)
		if err != nil {
			return zero, errors.Wrapf(err, "failed to apply patch at index %d", i)
		}
	}

	decoder := serializer.NewCodecFactory(scheme).UniversalDeserializer()
	obj, _, err := decoder.Decode(paramData, nil, nil)
	if err != nil {
		return zero, errors.Wrap(err, "failed to decode object")
	}

	t, ok := obj.(T)
	if !ok {
		return zero, fmt.Errorf("decoded patched object is %T, not %T", obj, parameters)
	}

	return t, nil
}

func (r *reconciler[T]) createOrUpdateResource(ctx context.Context, existing T, parameters client.Object, resourceExists bool, serviceName string) (T, error) {
	var zero T
	var err error
	var logMessageVerbPrefix string
	if resourceExists {
		logMessageVerbPrefix = "updat"
		err = r.Client.Patch(ctx, parameters, client.MergeFrom(existing))
	} else {
		logMessageVerbPrefix = "creat"
		err = r.Client.Create(ctx, parameters)
	}
	if err == nil {
		// Resources need to be requeued to wait for the create or update to finish.
		return zero, azure.WithTransientError(azure.NewOperationNotDoneError(&infrav1.Future{
			Type:          createOrUpdateFutureType,
			ResourceGroup: parameters.GetNamespace(),
			Name:          parameters.GetName(),
		}), requeueInterval)
	}
	return zero, errors.Wrapf(err, "failed to %se resource %s/%s (service: %s)", logMessageVerbPrefix, parameters.GetNamespace(), parameters.GetName(), serviceName)
}

// DeleteResource implements the logic for deleting a resource Asynchronously.
func (r *reconciler[T]) DeleteResource(ctx context.Context, resource T, serviceName string) (err error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "services.aso.DeleteResource")
	defer done()

	resource.SetNamespace(r.owner.GetNamespace())
	resourceName := resource.GetName()
	resourceNamespace := resource.GetNamespace()

	log = log.WithValues("service", serviceName, "resource", resourceName, "namespace", resourceNamespace)

	managed, err := IsManaged(ctx, r.Client, resource, r.owner)
	if apierrors.IsNotFound(err) {
		// already deleted
		log.V(2).Info("successfully deleted resource")
		return nil
	}
	if err != nil {
		return errors.Wrap(err, "failed to determine if resource is managed")
	}
	if !managed {
		log.V(4).Info("skipping delete for unmanaged resource")
		return nil
	}

	log.V(2).Info("deleting resource")
	err = r.Client.Delete(ctx, resource)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// already deleted
			log.V(2).Info("successfully deleted resource")
			return nil
		}
		return errors.Wrapf(err, "failed to delete resource %s/%s (service: %s)", resourceNamespace, resourceName, serviceName)
	}

	return azure.WithTransientError(azure.NewOperationNotDoneError(&infrav1.Future{
		Type:          deleteFutureType,
		ResourceGroup: resourceNamespace,
		Name:          resourceName,
	}), requeueInterval)
}

// IsManaged returns whether the ASO resource referred to by spec was created by
// CAPZ and therefore whether CAPZ should manage its lifecycle.
func IsManaged[T genruntime.MetaObject](ctx context.Context, ctrlClient client.Client, resource T, owner client.Object) (bool, error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "services.aso.IsManaged")
	defer done()

	resource.SetNamespace(owner.GetNamespace())

	err := ctrlClient.Get(ctx, client.ObjectKeyFromObject(resource), resource)
	if err != nil {
		return false, errors.Wrap(err, "error getting resource")
	}

	return isOwnedBy(resource, owner, ctrlClient.Scheme())
}

func isOwnedBy(resource client.Object, owner client.Object, scheme *runtime.Scheme) (bool, error) {
	ownerGVK, err := apiutil.GVKForObject(owner, scheme)
	if err != nil {
		return false, err
	}
	existingOwner := metav1.GetControllerOf(resource)
	return existingOwner != nil &&
		existingOwner.APIVersion == ownerGVK.GroupVersion().String() &&
		existingOwner.Kind == ownerGVK.Kind &&
		existingOwner.Name == owner.GetName(), nil
}

func hasLegacyOwnedByLabel(labels map[string]string, clusterName string) bool {
	return labels[infrav1.OwnedByClusterLabelKey] == clusterName //nolint:staticcheck // Referencing this deprecated value is required for backwards compatibility.
}

// PauseResource pauses an ASO resource by updating its `reconcile-policy` to `skip`.
func (r *reconciler[T]) PauseResource(ctx context.Context, resource T, serviceName string) error {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "services.aso.PauseResource")
	defer done()

	resource.SetNamespace(r.owner.GetNamespace())

	log = log.WithValues("service", serviceName, "resource", resource.GetName(), "namespace", resource.GetNamespace())

	if err := r.Client.Get(ctx, client.ObjectKeyFromObject(resource), resource); err != nil {
		return err
	}
	if isOwned, err := isOwnedBy(resource, r.owner, r.Scheme()); err != nil {
		return err
	} else if !isOwned {
		log.V(4).Info("Skipping pause of unmanaged resource")
		return nil
	}

	annotations := resource.GetAnnotations()
	if _, exists := annotations[prePauseReconcilePolicyAnnotation]; exists {
		log.V(4).Info("resource is already paused")
		return nil
	}

	log.V(4).Info("Pausing resource")
	before := resource.DeepCopyObject().(genruntime.MetaObject)

	if annotations == nil {
		annotations = make(map[string]string, 2)
	}
	annotations[prePauseReconcilePolicyAnnotation] = annotations[asoannotations.ReconcilePolicy]
	annotations[asoannotations.ReconcilePolicy] = string(asoannotations.ReconcilePolicySkip)
	resource.SetAnnotations(annotations)

	return r.Client.Patch(ctx, resource, client.MergeFrom(before))
}
