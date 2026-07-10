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

package orc

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"

	infrav1alpha1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha1"
	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta2"
	capoerrors "sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/errors"
)

// Reconciler creates and monitors ORC resources on behalf of an
// OpenStackServer. It replaces all direct Gophercloud API calls for
// server, port, and volume management.
type Reconciler struct {
	Client client.Client
	Scheme *runtime.Scheme
}

// ReconcileResult contains the outcome of an ORC reconciliation cycle.
type ReconcileResult struct {
	// Done is true when the ORC Server is Available and its status has
	// been extracted.
	Done bool
	// ServerID is the OpenStack Nova server UUID (set when Done is true).
	ServerID string
	// ServerState is the Nova instance state (set when Done is true).
	ServerState infrav1.InstanceState
	// Addresses are the server's network addresses (set when Done is true).
	Addresses []corev1.NodeAddress
}

// Reconcile creates all ORC sub-resources for the given OpenStackServer
// and monitors their status. It returns a ReconcileResult indicating
// whether the server is fully ready.
//
// On each call it:
//  1. Resolves credentials from the IdentityRef
//  2. Builds all ORC objects (Image, Flavor, KeyPair, Networks, Subnets,
//     SecurityGroups, Ports, Trunks, Volumes, Server)
//  3. Ensures each object exists (creates if missing; does not update)
//  4. Checks all objects for terminal errors
//  5. Returns the ORC Server status when Available
func (r *Reconciler) Reconcile(ctx context.Context, openStackServer *infrav1alpha1.OpenStackServer) (*ReconcileResult, error) {
	log := ctrl.LoggerFrom(ctx)
	serverName := openStackServer.Name
	namespace := openStackServer.Namespace
	spec := &openStackServer.Spec

	// ── 1. Resolve credentials ──────────────────────────────────────
	credRef, err := ResolveCloudCredentialsRef(ctx, r.Client, namespace, spec.IdentityRef)
	if err != nil {
		return nil, fmt.Errorf("resolving cloud credentials: %w", err)
	}

	// ── 2. Resolve ports → deduplicated Network/Subnet/SG objects ───
	portRes := ResolvePortsToORC(serverName, namespace, spec.Ports, spec.SecurityGroups, credRef)

	// ── 3. Determine image ORC name ─────────────────────────────────
	imageORCName := ImageName(serverName)
	if spec.Image.ImageRef != nil {
		imageORCName = spec.Image.ImageRef.Name
	}

	// ── 4. Build all ORC objects ────────────────────────────────────
	var orcObjects []client.Object

	// Image (nil when user provides an ImageRef to an existing ORC Image)
	if img := buildImage(serverName, namespace, spec.Image, credRef); img != nil {
		orcObjects = append(orcObjects, img)
	}

	// Flavor
	orcObjects = append(orcObjects, buildFlavor(serverName, namespace, spec.Flavor, spec.FlavorID, credRef))

	// KeyPair
	var keypairORCName string
	if spec.SSHKeyName != "" {
		keypairORCName = KeyPairName(serverName)
		orcObjects = append(orcObjects, buildKeypair(serverName, namespace, spec.SSHKeyName, credRef))
	}

	// ServerGroup
	var serverGroupORCName string
	if spec.ServerGroup != nil {
		serverGroupORCName = ServerGroupORCName(serverName)
		orcObjects = append(orcObjects, buildServerGroup(serverName, namespace, spec.ServerGroup, credRef))
	}

	// Networks, Subnets, SecurityGroups from port resolution
	for _, obj := range portRes.Networks {
		orcObjects = append(orcObjects, obj)
	}
	for _, obj := range portRes.Subnets {
		orcObjects = append(orcObjects, obj)
	}
	for _, obj := range portRes.SecurityGroups {
		orcObjects = append(orcObjects, obj)
	}

	// Volume types (deduplicated)
	volumeTypeNameMap := make(map[string]string)
	collectVolumeType := func(typeName string) {
		if typeName == "" {
			return
		}
		if _, exists := volumeTypeNameMap[typeName]; exists {
			return
		}
		vtName := VolumeTypeORCName(serverName, typeName)
		volumeTypeNameMap[typeName] = vtName
		orcObjects = append(orcObjects, buildVolumeType(serverName, namespace, typeName, credRef))
	}
	if spec.RootVolume != nil {
		collectVolumeType(spec.RootVolume.Type)
	}
	for i := range spec.AdditionalBlockDevices {
		bd := &spec.AdditionalBlockDevices[i]
		if bd.Storage.Type == infrav1.VolumeBlockDevice && bd.Storage.Volume != nil {
			collectVolumeType(bd.Storage.Volume.Type)
		}
	}

	// Root volume
	serverAZ := ptr.Deref(spec.AvailabilityZone, "")
	var rootVolumeORCName string
	if spec.RootVolume != nil && spec.RootVolume.SizeGiB > 0 {
		rootVolumeORCName = RootVolumeName(serverName)
		orcObjects = append(orcObjects, buildRootVolume(serverName, namespace, spec.RootVolume, imageORCName, volumeTypeNameMap, serverAZ, credRef))
	}

	// Additional volumes
	var additionalVolumes []VolumeAttachment
	for i := range spec.AdditionalBlockDevices {
		bd := &spec.AdditionalBlockDevices[i]
		switch bd.Storage.Type {
		case infrav1.VolumeBlockDevice:
			volName := AdditionalVolumeName(serverName, bd.Name)
			additionalVolumes = append(additionalVolumes, VolumeAttachment{ORCName: volName, Device: bd.Name})
			orcObjects = append(orcObjects, buildAdditionalVolume(serverName, namespace, *bd, volumeTypeNameMap, serverAZ, credRef))
		case infrav1.LocalBlockDevice:
			log.Info("WARNING: Local block devices are not supported by ORC, skipping", "blockDevice", bd.Name)
		}
	}

	// Ports and trunks
	var portORCNames []string
	for i, portOpts := range spec.Ports {
		portName := PortORCName(serverName, i)
		portORCNames = append(portORCNames, portName)

		portObj := buildPort(serverName, namespace, i, portOpts, spec.SecurityGroups,
			ptr.Deref(spec.Trunk, false), spec.Tags,
			portRes.NetworkNameMap, portRes.SubnetNameMap, portRes.SGNameMap, credRef)
		orcObjects = append(orcObjects, portObj)

		trunkEnabled := ptr.Deref(portOpts.Trunk, ptr.Deref(spec.Trunk, false))
		if trunkEnabled {
			orcObjects = append(orcObjects, buildTrunk(serverName, namespace, i, portName, spec.Tags, credRef))
		}
	}

	// ORC Server (depends on everything above)
	serverObj := buildServer(serverName, namespace, spec, imageORCName, FlavorName(serverName),
		keypairORCName, serverGroupORCName, rootVolumeORCName,
		portORCNames, additionalVolumes, credRef)
	orcObjects = append(orcObjects, serverObj)

	// ── 5. Set owner references and ensure all objects exist ────────
	for _, obj := range orcObjects {
		if err := controllerutil.SetControllerReference(openStackServer, obj, r.Scheme); err != nil {
			return nil, fmt.Errorf("setting owner reference on %s: %w", obj.GetName(), err)
		}
		if err := r.ensureORCResource(ctx, obj); err != nil {
			return nil, fmt.Errorf("ensuring ORC resource %s: %w", obj.GetName(), err)
		}
	}

	// ── 6. Check for terminal errors on all ORC objects ─────────────
	for _, obj := range orcObjects {
		if err := r.Client.Get(ctx, client.ObjectKeyFromObject(obj), obj); err != nil {
			return nil, fmt.Errorf("fetching ORC resource %s: %w", obj.GetName(), err)
		}
		if condObj, ok := obj.(orcv1alpha1.ObjectWithConditions); ok {
			if termErr := orcv1alpha1.GetTerminalError(condObj); termErr != nil {
				return nil, capoerrors.Terminal(infrav1.DependencyFailedReason,
					fmt.Sprintf("ORC resource %s/%s failed: %v", obj.GetNamespace(), obj.GetName(), termErr))
			}
		}
	}

	// ── 7. Check if the ORC Server is Available ─────────────────────
	orcServer := &orcv1alpha1.Server{}
	serverKey := client.ObjectKey{Namespace: namespace, Name: ServerName(serverName)}
	if err := r.Client.Get(ctx, serverKey, orcServer); err != nil {
		return nil, fmt.Errorf("fetching ORC Server: %w", err)
	}

	result := &ReconcileResult{}

	if !orcv1alpha1.IsAvailable(orcServer) {
		return result, nil
	}

	// Server is Available → extract status
	result.Done = true
	if orcServer.Status.ID != nil {
		result.ServerID = *orcServer.Status.ID
	}
	if orcServer.Status.Resource != nil {
		result.ServerState = MapInstanceState(orcServer.Status.Resource.Status)
		result.Addresses = MapAddresses(orcServer.Status.Resource.Interfaces)
	}

	return result, nil
}

// DeleteORCResources deletes the ORC Server and waits for its deletion.
// All other ORC sub-resources are garbage-collected via OwnerReferences
// when the OpenStackServer is deleted.
//
// Returns (true, nil) when the ORC Server no longer exists.
// Returns (false, nil) when the deletion is in progress.
func (r *Reconciler) DeleteORCResources(ctx context.Context, openStackServer *infrav1alpha1.OpenStackServer) (bool, error) {
	orcServer := &orcv1alpha1.Server{}
	serverKey := client.ObjectKey{
		Namespace: openStackServer.Namespace,
		Name:      ServerName(openStackServer.Name),
	}

	err := r.Client.Get(ctx, serverKey, orcServer)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return true, nil // already gone
		}
		return false, fmt.Errorf("getting ORC Server for deletion: %w", err)
	}

	// Delete if not already being deleted
	if orcServer.DeletionTimestamp.IsZero() {
		if err := r.Client.Delete(ctx, orcServer); err != nil {
			if apierrors.IsNotFound(err) {
				return true, nil
			}
			return false, fmt.Errorf("deleting ORC Server: %w", err)
		}
	}

	// Server exists but is being deleted — the Owns() watch will
	// re-trigger the controller when the Server is fully removed.
	return false, nil
}

// ensureORCResource creates an ORC resource if it does not already exist.
// It does not update existing resources (ORC specs are immutable).
func (r *Reconciler) ensureORCResource(ctx context.Context, obj client.Object) error {
	existing := obj.DeepCopyObject().(client.Object)
	err := r.Client.Get(ctx, client.ObjectKeyFromObject(obj), existing)
	if err != nil {
		if apierrors.IsNotFound(err) {
			if createErr := r.Client.Create(ctx, obj); createErr != nil {
				if apierrors.IsAlreadyExists(createErr) {
					return nil // race condition, object was created concurrently
				}
				return createErr
			}
			return nil
		}
		return err
	}
	// Object already exists — nothing to do.
	return nil
}
