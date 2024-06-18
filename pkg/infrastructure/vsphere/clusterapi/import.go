package clusterapi

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vmware/govmomi/govc/importx"
	"github.com/vmware/govmomi/nfc"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/ovf"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/session"

	"github.com/openshift/installer/pkg/types/vsphere"
)

func importRhcosOva(ctx context.Context, session *session.Session, folder *object.Folder, cachedImage, clusterID, tagID, diskProvisioningType string, failureDomain vsphere.FailureDomain) error {
	name := fmt.Sprintf("%s-rhcos-%s-%s", clusterID, failureDomain.Region, failureDomain.Zone)
	logrus.Infof("Importing OVA %v into failure domain %v.", name, failureDomain.Name)
	archive := &importx.ArchiveFlag{Archive: &importx.TapeArchive{Path: cachedImage}}

	ovfDescriptor, err := archive.ReadOvf("*.ovf")
	if err != nil {
		// Open the corrupt OVA file
		f, ferr := os.Open(cachedImage)
		if ferr != nil {
			err = fmt.Errorf("%s, %w", err.Error(), ferr)
		}
		defer f.Close()

		// Get a sha256 on the corrupt OVA file
		// and the size of the file
		h := sha256.New()
		written, cerr := io.Copy(h, f)
		if cerr != nil {
			err = fmt.Errorf("%s, %w", err.Error(), cerr)
		}

		return fmt.Errorf("ova %s has a sha256 of %x and a size of %d bytes, failed to read the ovf descriptor %w", cachedImage, h.Sum(nil), written, err)
	}

	ovfEnvelope, err := archive.ReadEnvelope(ovfDescriptor)
	if err != nil {
		return fmt.Errorf("failed to parse ovf: %w", err)
	}

	// The RHCOS OVA only has one network defined by default
	// The OVF envelope defines this.  We need a 1:1 mapping
	// between networks with the OVF and the host
	if len(ovfEnvelope.Network.Networks) != 1 {
		return fmt.Errorf("expected the OVA to only have a single network adapter")
	}

	cluster, err := session.Finder.ClusterComputeResource(ctx, failureDomain.Topology.ComputeCluster)
	if err != nil {
		return fmt.Errorf("failed to find compute cluster: %w", err)
	}

	clusterHostSystems, err := cluster.Hosts(ctx)

	if err != nil {
		return fmt.Errorf("failed to get cluster hosts: %w", err)
	}
	resourcePool, err := session.Finder.ResourcePool(ctx, failureDomain.Topology.ResourcePool)
	if err != nil {
		return fmt.Errorf("failed to find resource pool: %w", err)
	}

	networkPath := path.Join(cluster.InventoryPath, failureDomain.Topology.Networks[0])

	networkRef, err := session.Finder.Network(ctx, networkPath)
	if err != nil {
		return fmt.Errorf("failed to find network: %w", err)
	}
	datastore, err := session.Finder.Datastore(ctx, failureDomain.Topology.Datastore)
	if err != nil {
		return fmt.Errorf("failed to find datastore: %w", err)
	}

	// Create mapping between OVF and the network object
	// found by Name
	networkMappings := []types.OvfNetworkMapping{{
		Name:    ovfEnvelope.Network.Networks[0].Name,
		Network: networkRef.Reference(),
	}}

	// This is a very minimal spec for importing an OVF.
	cisp := types.OvfCreateImportSpecParams{
		EntityName:     name,
		NetworkMapping: networkMappings,
	}

	switch diskProvisioningType {
	case "":
		// Disk provisioning type will be set according to the default storage policy of vsphere.
	case "thin":
		cisp.DiskProvisioning = string(types.OvfCreateImportSpecParamsDiskProvisioningTypeThin)
	case "thick":
		cisp.DiskProvisioning = string(types.OvfCreateImportSpecParamsDiskProvisioningTypeThick)
	case "eagerZeroedThick":
		cisp.DiskProvisioning = string(types.OvfCreateImportSpecParamsDiskProvisioningTypeEagerZeroedThick)
	default:
		return errors.Errorf("disk provisioning type %q is not supported", diskProvisioningType)
	}

	m := ovf.NewManager(session.Client.Client)
	spec, err := m.CreateImportSpec(ctx,
		string(ovfDescriptor),
		resourcePool.Reference(),
		datastore.Reference(),
		cisp)

	if err != nil {
		return fmt.Errorf("failed to create import spec: %w", err)
	}
	if spec.Error != nil {
		return errors.New(spec.Error[0].LocalizedMessage)
	}

	hostSystem, err := findAvailableHostSystems(ctx, session, clusterHostSystems)
	if err != nil {
		return fmt.Errorf("failed to find available host system: %w", err)
	}

	lease, err := resourcePool.ImportVApp(ctx, spec.ImportSpec, folder, hostSystem)

	if err != nil {
		return fmt.Errorf("failed to import vapp: %w", err)
	}

	info, err := lease.Wait(ctx, spec.FileItem)
	if err != nil {
		return fmt.Errorf("failed to lease wait: %w", err)
	}

	u := lease.StartUpdater(ctx, info)
	defer u.Done()

	for _, i := range info.Items {
		// upload the vmdk to which ever host that was first
		// available with the required network and datastore.
		err = upload(ctx, archive, lease, i)
		if err != nil {
			return fmt.Errorf("failed to upload: %w", err)
		}
	}

	err = lease.Complete(ctx)
	if err != nil {
		return fmt.Errorf("failed to lease complete: %w", err)
	}

	vm := object.NewVirtualMachine(session.Client.Client, info.Entity)
	if vm == nil {
		return fmt.Errorf("error VirtualMachine not found, managed object id: %s", info.Entity.Value)
	}

	err = vm.MarkAsTemplate(ctx)
	if err != nil {
		return fmt.Errorf("failed to mark vm as template: %w", err)
	}
	err = attachTag(ctx, session, vm.Reference().Value, tagID)
	if err != nil {
		return fmt.Errorf("failed to attach tag: %w", err)
	}

	return nil
}

func findAvailableHostSystems(ctx context.Context, session *session.Session, clusterHostSystems []*object.HostSystem) (*object.HostSystem, error) {
	var hostSystemManagedObject mo.HostSystem
	for _, hostObj := range clusterHostSystems {
		err := hostObj.Properties(ctx, hostObj.Reference(), []string{"config.product", "network", "datastore", "runtime"}, &hostSystemManagedObject)
		if err != nil {
			return nil, err
		}
		if hostSystemManagedObject.Runtime.InMaintenanceMode {
			continue
		}
		return hostObj, nil
	}
	return nil, errors.New("all hosts unavailable")
}

// Used govc/importx/ovf.go as an example to implement
// resourceVspherePrivateImportOvaCreate and upload functions
// See: https://github.com/vmware/govmomi/blob/cc10a0758d5b4d4873388bcea417251d1ad03e42/govc/importx/ovf.go#L196-L324
func upload(ctx context.Context, archive *importx.ArchiveFlag, lease *nfc.Lease, item nfc.FileItem) error {
	file := item.Path

	f, size, err := archive.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	opts := soap.Upload{
		ContentLength: size,
	}

	return lease.Upload(ctx, item, f, opts)
}
