/*
Copyright 2026 The Kubernetes Authors.

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

package asomigration

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/go-logr/logr"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
)

// storedVersionMigration describes an ASO-managed CRD whose
// status.storedVersions still references a version that the bundled ASO
// release no longer lists in spec.versions. The Kubernetes API server rejects
// a CRD update that drops such a version until a storage migration ensures no
// data remains persisted in it, so an upgrade would otherwise fail with:
//
//	status.storedVersions[i]: Invalid value: "<version>": missing from
//	spec.versions; ... must remain in spec.versions until a storage migration
//	ensures no data remains persisted in <version> ...
type storedVersionMigration struct {
	// crd is the name of the affected CustomResourceDefinition.
	crd string
	// stale is the storedVersions entry the bundled ASO no longer serves.
	stale string
}

// storedVersionMigrations lists the known-affected ASO-managed CRDs. Add an
// entry here when a future ASO release removes a previously-stored version.
// ASO v2.18 removes containerservice/v1api20231001 from ManagedCluster and
// AgentPool; clusters that ran an early CAPZ-bundled ASO (v2.5-v2.8, where
// v1api20231001storage was the hub) may still list it in status.storedVersions.
var storedVersionMigrations = []storedVersionMigration{
	{crd: "fleetsmembers.containerservice.azure.com", stale: "v1api20230315previewstorage"},
	{crd: "managedclusters.containerservice.azure.com", stale: "v1api20231001storage"},
	{crd: "managedclustersagentpools.containerservice.azure.com", stale: "v1api20231001storage"},
}

// crdDeletionTimeout bounds how long we wait for a deleted CRD to disappear.
const crdDeletionTimeout = 2 * time.Minute

// MigrateStoredVersions cleans up obsolete status.storedVersions entries for
// ASO-managed CRDs before the bundled ASO controller tries to apply them. It
// runs as an init container in the ASO deployment (via the `manager
// migrate-aso-crds` subcommand) so it completes before ASO reconciles its CRDs,
// reusing the CAPZ manager image instead of depending on an external one.
func MigrateStoredVersions(ctx context.Context, cfg *rest.Config) error {
	log := ctrl.Log.WithName("migrate-aso-crds")

	crdClient, err := apiextensionsclientset.NewForConfig(cfg)
	if err != nil {
		return fmt.Errorf("creating apiextensions client: %w", err)
	}
	dynClient, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return fmt.Errorf("creating dynamic client: %w", err)
	}

	for _, m := range storedVersionMigrations {
		if err := migrateStoredVersion(ctx, log, crdClient, dynClient, m); err != nil {
			return err
		}
	}
	return nil
}

func migrateStoredVersion(ctx context.Context, log logr.Logger, crdClient apiextensionsclientset.Interface, dynClient dynamic.Interface, m storedVersionMigration) error {
	crd, err := crdClient.ApiextensionsV1().CustomResourceDefinitions().Get(ctx, m.crd, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		log.Info("CRD not present; skipping", "crd", m.crd)
		return nil
	}
	if err != nil {
		return fmt.Errorf("getting CRD %s: %w", m.crd, err)
	}

	if !slices.Contains(crd.Status.StoredVersions, m.stale) {
		log.Info("CRD has no stale stored version; skipping", "crd", m.crd, "staleVersion", m.stale)
		return nil
	}

	// Every served version exposes the same stored objects (the API server
	// converts on the fly), so take the max instance count across the served
	// versions rather than summing, which would overcount by the number of
	// versions. Unserved versions (e.g. the stale storage version) can't be
	// listed directly, but their stored objects surface through the served
	// versions we do list.
	instances := 0
	for _, v := range crd.Spec.Versions {
		if !v.Served {
			continue
		}
		gvr := schema.GroupVersionResource{
			Group:    crd.Spec.Group,
			Version:  v.Name,
			Resource: crd.Spec.Names.Plural,
		}
		n, err := countInstances(ctx, dynClient, gvr)
		if err != nil {
			// Fail closed: a real listing failure (conversion webhook down,
			// API server hiccup, RBAC denial) must never be misread as "zero
			// instances" and let us delete a CRD that still holds data.
			return fmt.Errorf("listing %s to confirm CRD %s is empty: %w", gvr, m.crd, err)
		}
		instances = max(instances, n)
	}

	if instances != 0 {
		return fmt.Errorf("CRD %s has %d instance(s) but its status.storedVersions still references %q, "+
			"which the bundled Azure Service Operator no longer serves. Run 'asoctl clean crds' against this "+
			"cluster to migrate the stored objects before retrying the CAPZ upgrade. See "+
			"https://azure.github.io/azure-service-operator/guide/crd-management/ for details",
			m.crd, instances, m.stale)
	}

	// With no instances, delete the CRD rather than patching
	// status.storedVersions. A patch is a no-op when the stale version is the
	// active storage version (the API server refuses to drop it), and the
	// bundled ASO spec.versions for these CRDs is fully disjoint from the
	// pre-existing definition, so there is no served version we could flip the
	// storage flag onto that the new bundle keeps. The new ASO controller
	// recreates the CRD from its bundled definition with status.storedVersions
	// containing only the new storage version.
	log.Info("CRD has no instances and stale storedVersions; deleting so the new ASO controller can recreate it",
		"crd", m.crd, "storedVersions", crd.Status.StoredVersions)
	if err := crdClient.ApiextensionsV1().CustomResourceDefinitions().Delete(ctx, m.crd, metav1.DeleteOptions{}); err != nil && !apierrors.IsNotFound(err) {
		return fmt.Errorf("deleting CRD %s: %w", m.crd, err)
	}
	if err := waitForCRDDeletion(ctx, crdClient, m.crd); err != nil {
		return err
	}
	log.Info("CRD deleted", "crd", m.crd)
	return nil
}

// countInstances returns the number of objects of the given served resource
// across all namespaces. Any list error is returned to the caller, which fails
// closed rather than risk deleting a CRD whose emptiness it couldn't confirm.
func countInstances(ctx context.Context, dynClient dynamic.Interface, gvr schema.GroupVersionResource) (int, error) {
	total := 0
	opts := metav1.ListOptions{Limit: 500}
	for {
		list, err := dynClient.Resource(gvr).List(ctx, opts)
		if err != nil {
			return 0, err
		}
		total += len(list.Items)
		cont := list.GetContinue()
		if cont == "" {
			return total, nil
		}
		opts.Continue = cont
	}
}

// waitForCRDDeletion blocks until the named CRD is gone or the timeout elapses.
func waitForCRDDeletion(ctx context.Context, crdClient apiextensionsclientset.Interface, name string) error {
	ctx, cancel := context.WithTimeout(ctx, crdDeletionTimeout)
	defer cancel()
	err := wait.PollUntilContextCancel(ctx, time.Second, true, func(ctx context.Context) (bool, error) {
		_, err := crdClient.ApiextensionsV1().CustomResourceDefinitions().Get(ctx, name, metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			return true, nil
		}
		return false, err
	})
	if err != nil {
		return fmt.Errorf("waiting for CRD %s to be deleted: %w", name, err)
	}
	return nil
}
