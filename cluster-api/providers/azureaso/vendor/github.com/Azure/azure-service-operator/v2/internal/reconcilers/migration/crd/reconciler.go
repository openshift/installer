/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package crd

import (
	"context"
	"time"

	. "github.com/Azure/azure-service-operator/v2/internal/logging"

	"github.com/rotisserie/eris"
	"github.com/samber/lo"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/selection"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/Azure/azure-service-operator/v2/internal/reflecthelpers"
	"github.com/Azure/azure-service-operator/v2/internal/set"
	"github.com/Azure/azure-service-operator/v2/internal/version"
	"github.com/Azure/azure-service-operator/v2/pkg/common/controllerutils/singleton"
	asolabels "github.com/Azure/azure-service-operator/v2/pkg/common/labels"
)

var _ singleton.Reconciler = &Reconciler{}

type Options struct {
	VersionProvider func() string
}

type Reconciler struct {
	kubeClient      client.Client
	cache           cache.Cache
	crdsToMigrate   map[string][]string // map of CRD names to versions that need to be migrated
	crdsMigrated    set.Set[string]
	versionProvider func() string
}

// NewReconciler creates a new Reconciler for migrating resources to the latest version.
func NewReconciler(
	kubeClient client.Client, // should pass a client that is direct from server, at least for CRDs themselves.
	cache cache.Cache,
	crdsToMigrate map[string][]string, // map of CRD names (ex: trustedaccessrolebindings.containerservice.azure.com) to versions that need to be migrated
	opts Options,
) *Reconciler {
	if opts.VersionProvider == nil {
		opts.VersionProvider = func() string {
			return version.BuildVersion
		}
	}

	return &Reconciler{
		kubeClient:      kubeClient,
		cache:           cache,
		crdsToMigrate:   crdsToMigrate,
		crdsMigrated:    set.Make[string](),
		versionProvider: opts.VersionProvider,
	}
}

func (r *Reconciler) Reconcile(ctx context.Context) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Get CRDs (only the ones we're interested in, one by one). This results in lower memory usage as it avoids
	// a huge CRD list. We're in no hurry to migrate all CRDs at once so this is fine.
	for crdName, deprecatedVersions := range r.crdsToMigrate {
		log := log.WithValues("crd", crdName)

		// Check if we have already migrated this CRD
		if r.crdsMigrated.Contains(crdName) {
			log.V(Debug).Info("CRD already migrated, skipping")
			continue
		}

		// Get the CRD
		obj := &apiextensions.CustomResourceDefinition{}
		if err := r.kubeClient.Get(ctx, client.ObjectKey{Name: crdName}, obj); err != nil {
			if client.IgnoreNotFound(err) != nil {
				return ctrl.Result{}, eris.Wrapf(err, "failed to get CRD %s", crdName)
			}
			log.V(Debug).Info("CRD not found, skipping")
			continue
		}

		// Check if the CRD still has the deprecated version(s). If it does, we need to migrate it.
		// If it doesn't, we can add it to the list of migrated CRDs and skip it.
		if !lo.Some(obj.Status.StoredVersions, deprecatedVersions) {
			log.V(Debug).Info("CRD does not contain deprecated versions, skipping", "deprecatedVersions", deprecatedVersions)
			r.crdsMigrated.Add(crdName)
			continue
		}

		// Check if all instances of the CRD are migrated
		allMigrated, err := r.areAllInstancesMigrated(ctx, obj)
		if err != nil {
			return ctrl.Result{}, eris.Wrapf(err, "failed to check if all instances of CRD %s are migrated", crdName)
		}
		if !allMigrated {
			log.V(Debug).Info("Not all instances of CRD are migrated, skipping")
			continue
		}

		// Actually remove the deprecated versions from the status.storedVersions
		old := obj.DeepCopy()
		obj.Status.StoredVersions, _ = lo.Difference(obj.Status.StoredVersions, deprecatedVersions)
		log.V(Status).Info("Removing deprecated versions from CRD storedVersions", "deprecatedVersions", deprecatedVersions, "storedVersions", obj.Status.StoredVersions)

		if err := r.kubeClient.Status().Patch(ctx, obj, client.MergeFrom(old)); err != nil {
			return ctrl.Result{}, eris.Wrapf(err, "failed to patch CRD %s", crdName)
		}

		r.crdsMigrated.Add(crdName)
	}

	// All done, requeue to check again in an hour
	return ctrl.Result{
		RequeueAfter: 1 * time.Hour,
	}, nil
}

func (r *Reconciler) areAllInstancesMigrated(ctx context.Context, crd *apiextensions.CustomResourceDefinition) (bool, error) {
	log := log.FromContext(ctx).WithValues("crd", crd.Name)

	// Pick the latest supported version in the CRD.
	gvk, err := r.getStorageVersion(crd)
	if err != nil {
		return false, eris.Wrapf(err, "failed to get storage version for CRD %s", crd.Name)
	}

	if r.cache != nil {
		// This is a bit of a hack, but we don't want to run this controller until the cache is synced.
		// Otherwise, there's a race where we can list resources and hit an empty cache, which will cause us to
		// think it's safe to migrate a resource when it is not.
		log.V(Debug).Info("Waiting for cache to sync")
		_, err = r.cache.GetInformerForKind(ctx, gvk, cache.BlockUntilSynced(true))
		if err != nil {
			return false, eris.Wrapf(err, "failed to get informer for GVK %s", gvk)
		}
		log.V(Debug).Info("Done waiting for cache to sync")
	}

	gvk.Kind = gvk.Kind + "List" // We want the list kind, not the singular kind

	// Use the full type, not partialObjectMetadata, as controller-runtime caches per-type
	// (where the types are Unstructured, PartialObjectMetadata, and Structured). We do NOT want to
	// have a separate cache for these, so we use the full type to re-use the cache we've already got.
	// Worth noting that we've disabled the cache for CRDs because holding them in memory is expensive,
	// and this controller should eventually migrate everything at which point we don't need to make any calls
	// anymore.
	scheme := r.kubeClient.Scheme()
	// Create an instance of the type to get the type name
	obj, err := scheme.New(gvk)
	if err != nil {
		return false, eris.Wrapf(err, "failed to create instance of GVK %s", gvk)
	}

	// We expect that obj is of type client.ObjectList
	list, ok := obj.(client.ObjectList)
	if !ok {
		return false, eris.Errorf("expected object of type client.ObjectList, got %T", obj)
	}

	selector := labels.NewSelector()
	req, err := labels.NewRequirement(
		asolabels.LastReconciledVersionLabel,
		selection.NotEquals,
		[]string{r.versionProvider()},
	)
	if err != nil {
		return false, eris.Wrapf(err, "failed to create label requirement for %s", crd.Name)
	}

	selector = selector.Add(*req)
	err = r.kubeClient.List(
		ctx,
		list,
		&client.ListOptions{
			LabelSelector: selector,
			Namespace:     metav1.NamespaceAll, // We want to check all namespaces
		},
	)
	if err != nil {
		return false, eris.Wrapf(err, "failed to list resources for CRD %s", crd.Name)
	}

	items, err := reflecthelpers.GetObjectListItems(list)
	if err != nil {
		return false, eris.Wrapf(err, "failed to get items from list for CRD %s", crd.Name)
	}
	log.V(Status).Info("Checking if CRD has all instances migrated", "countNotMigrated", len(items))
	return len(items) == 0, nil
}

func (r *Reconciler) getStorageVersion(crd *apiextensions.CustomResourceDefinition) (schema.GroupVersionKind, error) {
	for _, version := range crd.Spec.Versions {
		if version.Served && version.Storage {
			return schema.GroupVersionKind{
				Group:   crd.Spec.Group,
				Version: version.Name,
				Kind:    crd.Spec.Names.Kind,
			}, nil
		}
	}
	return schema.GroupVersionKind{}, eris.Errorf("no storage version found for CRD %s", crd.Name)
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: 1, // We never allow more than 1 reconcile at a time for this resource
		}).
		Named("migration_crd").
		WatchesRawSource(singleton.Source()).
		Complete(singleton.AsReconciler(r))
}
