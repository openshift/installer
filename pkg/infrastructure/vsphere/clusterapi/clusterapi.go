package clusterapi

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"crypto/sha256"

	"github.com/openshift/installer/pkg/asset"
	vcentercontexts "github.com/openshift/installer/pkg/asset/cluster/vsphere"
	"github.com/openshift/installer/pkg/asset/installconfig"
	icasset "github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/rhcos"
	"github.com/openshift/installer/pkg/infrastructure"
	"github.com/openshift/installer/pkg/rhcos/cache"
	ictypes "github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
	"github.com/pkg/errors"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/govc/importx"
	"github.com/vmware/govmomi/nfc"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/ovf"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/session"
)

// InfraProvider is the AWS SDK infra provider.
type InfraProvider struct{}

// InitializeProvider initializes the AWS SDK provider.
func InitializeProvider() infrastructure.Provider {
	return InfraProvider{}
}

func attachTag(ctx context.Context, session *session.Session, vmMoRefValue, tagId string) error {
	tagManager := session.TagManager

	moRef := types.ManagedObjectReference{
		Value: vmMoRefValue,
		Type:  "VirtualMachine",
	}

	err := tagManager.AttachTag(ctx, tagId, moRef)

	if err != nil {
		return fmt.Errorf("unable to attach tag: %s", err)
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

func importRhcosOva(ctx context.Context, session *session.Session, folder *object.Folder, cachedImage, clusterId, tagId, diskProvisioningType string, failureDomain vsphere.FailureDomain) error {
	name := fmt.Sprintf("%s-rhcos-%s-%s", clusterId, failureDomain.Region, failureDomain.Zone)

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

		return fmt.Errorf("ova %s has a sha256 of %x and a size of %d bytes, failed to read the ovf descriptor %s", cachedImage, h.Sum(nil), written, err)
	}

	ovfEnvelope, err := archive.ReadEnvelope(ovfDescriptor)
	if err != nil {
		return fmt.Errorf("failed to parse ovf: %s", err)
	}

	// The RHCOS OVA only has one network defined by default
	// The OVF envelope defines this.  We need a 1:1 mapping
	// between networks with the OVF and the host
	if len(ovfEnvelope.Network.Networks) != 1 {
		return fmt.Errorf("expected the OVA to only have a single network adapter")
	}

	cluster, err := session.Finder.ClusterComputeResource(ctx, failureDomain.Topology.ComputeCluster)
	if err != nil {
		return fmt.Errorf("failed to find compute cluster: %s", err)
	}

	clusterHostSystems, err := cluster.Hosts(ctx)

	if err != nil {
		return fmt.Errorf("failed to get cluster hosts: %s", err)
	}
	resourcePool, err := session.Finder.ResourcePool(ctx, failureDomain.Topology.ResourcePool)
	if err != nil {
		return fmt.Errorf("failed to find resource pool: %s", err)
	}

	networkPath := path.Join(cluster.InventoryPath, failureDomain.Topology.Networks[0])

	networkRef, err := session.Finder.Network(ctx, networkPath)
	if err != nil {
		return fmt.Errorf("failed to find network: %s", err)

	}
	datastore, err := session.Finder.Datastore(ctx, failureDomain.Topology.Datastore)
	if err != nil {
		return fmt.Errorf("failed to find datastore: %s", err)
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

	m := ovf.NewManager(session.Client.Client)
	spec, err := m.CreateImportSpec(ctx,
		string(ovfDescriptor),
		resourcePool.Reference(),
		datastore.Reference(),
		cisp)

	if err != nil {
		return fmt.Errorf("failed to create import spec: %s", err)
	}
	if spec.Error != nil {
		return errors.New(spec.Error[0].LocalizedMessage)
	}

	hostSystem, err := findAvailableHostSystems(ctx, session, clusterHostSystems)
	if err != nil {
		return fmt.Errorf("failed to find available host system: %s", err)
	}

	lease, err := resourcePool.ImportVApp(ctx, spec.ImportSpec, folder, hostSystem)

	if err != nil {
		return fmt.Errorf("failed to import vapp: %s", err)
	}

	info, err := lease.Wait(ctx, spec.FileItem)
	if err != nil {
		return fmt.Errorf("failed to lease wait: %s", err)
	}

	u := lease.StartUpdater(ctx, info)
	defer u.Done()

	for _, i := range info.Items {
		// upload the vmdk to which ever host that was first
		// available with the required network and datastore.
		err = upload(ctx, archive, lease, i)
		if err != nil {
			return fmt.Errorf("failed to upload: %s", err)
		}
	}

	err = lease.Complete(ctx)
	if err != nil {
		return fmt.Errorf("failed to lease complete: %s", err)
	}

	vm := object.NewVirtualMachine(session.Client.Client, info.Entity)
	if vm == nil {
		return fmt.Errorf("error VirtualMachine not found, managed object id: %s", info.Entity.Value)
	}

	err = vm.MarkAsTemplate(ctx)
	if err != nil {
		return fmt.Errorf("failed to mark vm as template: %s", err)
	}
	err = attachTag(ctx, session, vm.Reference().Value, tagId)
	if err != nil {
		return fmt.Errorf("failed to attach tag: %s", err)
	}

	return nil
}

func createFolder(ctx context.Context, fullpath string, session *session.Session) (*object.Folder, error) {
	dir := path.Dir(fullpath)
	base := path.Base(fullpath)
	finder := session.Finder

	folder, err := finder.Folder(ctx, fullpath)

	if folder == nil {
		folder, err = finder.Folder(ctx, dir)

		var notFoundError *find.NotFoundError
		if errors.As(err, &notFoundError) {
			folder, err = createFolder(ctx, dir, session)
			if err != nil {
				return folder, err
			}
		}

		if folder != nil {
			folder, err = folder.CreateFolder(ctx, base)
			if err != nil {
				return folder, err
			}
		}
	}
	return folder, err
}

func initializeFoldersAndTemplates(ctx context.Context, rhcosImage *rhcos.Image, installConfig *installconfig.InstallConfig, session *session.Session, clusterId, server string, vcenterContexts *vcentercontexts.VCenterContexts) error {
	finder := session.Finder

	platform := installConfig.Config.VSphere
	failureDomains := platform.FailureDomains

	for _, failureDomain := range failureDomains {
		dc, err := finder.Datacenter(ctx, failureDomain.Topology.Datacenter)
		if err != nil {
			return err
		}
		dcFolders, err := dc.Folders(ctx)
		if err != nil {
			return fmt.Errorf("unable to get datacenter folder: %v", err)
		}

		folderPath := path.Join(dcFolders.VmFolder.InventoryPath, clusterId)

		// we must set the Folder to the infraId somewhere, we will need to remove that.
		// if we are overwriting folderPath it needs to have a slash (path)
		folder := failureDomain.Topology.Folder
		if strings.Contains(folder, "/") {
			folderPath = folder
		}

		folderMo, err := createFolder(ctx, folderPath, session)
		if err != nil {
			return fmt.Errorf("unable to create folder: %v", err)
		}

		cachedImage, err := cache.DownloadImageFile(string(*rhcosImage), cache.InstallerApplicationName)
		if err != nil {
			return fmt.Errorf("failed to use cached vsphere image: %v", err)
		}

		// if the template is empty, the ova must be imported
		if len(failureDomain.Topology.Template) == 0 {
			if err = importRhcosOva(ctx, session, folderMo,
				cachedImage, clusterId, vcenterContexts.VCenters[server].TagID, string(platform.DiskType), failureDomain); err != nil {
				return fmt.Errorf("failed to import ova: %v", err)
			}
		}
	}
	return nil
}

func (a InfraProvider) Name() string {
	return vsphere.Name
}

func (a InfraProvider) Provision(dir string, parents asset.Parents) ([]*asset.File, error) {
	ctx := context.TODO()

	installConfig := &icasset.InstallConfig{}
	vcenterContexts := &vcentercontexts.VCenterContexts{}
	rhcosImage := new(rhcos.Image)
	clusterID := &installconfig.ClusterID{}

	parents.Get(
		vcenterContexts,
		installConfig,
		rhcosImage,
		clusterID)

	for _, vcenter := range installConfig.Config.VSphere.VCenters {
		server := vcenter.Server
		params := session.NewParams().WithServer(server).WithUserInfo(vcenter.Username, vcenter.Password)
		tempConnection, err := session.GetOrCreate(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("unable to create session: %v", err)
		}

		defer tempConnection.CloseIdleConnections()

		for _, failureDomain := range installConfig.Config.VSphere.FailureDomains {
			if failureDomain.Server != server {
				continue
			}
			if err = initializeFoldersAndTemplates(ctx, rhcosImage, installConfig, tempConnection, clusterID.InfraID, server, vcenterContexts); err != nil {
				return nil, fmt.Errorf("unable to initialize folders and templates: %v", err)
			}
		}
	}

	return nil, nil
}

func (a InfraProvider) DestroyBootstrap(dir string) error {
	return nil
}

func (a InfraProvider) ExtractHostAddresses(dir string, ic *ictypes.InstallConfig, ha *infrastructure.HostAddresses) error {
	return nil
}
