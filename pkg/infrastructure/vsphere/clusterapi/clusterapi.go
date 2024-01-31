package clusterapi

import (
	"context"
	"fmt"
	"path"
	"strings"

	"sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/session"

	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/rhcos"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	"github.com/openshift/installer/pkg/rhcos/cache"
	"github.com/openshift/installer/pkg/types/vsphere"
)

// Provider is the vSphere implementation of the clusterapi InfraProvider.
type Provider struct {
	clusterapi.InfraProvider
}

var _ clusterapi.PreProvider = Provider{}

// Name returns the vsphere provider name.
func (p Provider) Name() string {
	return vsphere.Name
}

func initializeFoldersAndTemplates(ctx context.Context, rhcosImage *rhcos.Image, installConfig *installconfig.InstallConfig, session *session.Session, clusterID, server, tagID string) error {
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
			return fmt.Errorf("unable to get datacenter folder: %w", err)
		}

		folderPath := path.Join(dcFolders.VmFolder.InventoryPath, clusterID)

		// we must set the Folder to the infraId somewhere, we will need to remove that.
		// if we are overwriting folderPath it needs to have a slash (path)
		folder := failureDomain.Topology.Folder
		if strings.Contains(folder, "/") {
			folderPath = folder
		}

		folderMo, err := createFolder(ctx, folderPath, session)
		if err != nil {
			return fmt.Errorf("unable to create folder: %w", err)
		}

		cachedImage, err := cache.DownloadImageFile(string(*rhcosImage), cache.InstallerApplicationName)
		if err != nil {
			return fmt.Errorf("failed to use cached vsphere image: %w", err)
		}

		// if the template is empty, the ova must be imported
		if len(failureDomain.Topology.Template) == 0 {
			if err = importRhcosOva(ctx, session, folderMo,
				cachedImage, clusterID, tagID, string(platform.DiskType), failureDomain); err != nil {
				return fmt.Errorf("failed to import ova: %w", err)
			}
		}
	}
	return nil
}

// PreProvision creates the vCenter objects required prior to running capv.
func (p Provider) PreProvision(ctx context.Context, in clusterapi.PreProvisionInput) error {
	installConfig := in.InstallConfig
	clusterID := &installconfig.ClusterID{InfraID: in.InfraID}
	var tagID string

	for _, vcenter := range installConfig.Config.VSphere.VCenters {
		server := vcenter.Server
		vctrSession, err := installConfig.VSphere.Session(context.TODO(), server)

		if err != nil {
			return err
		}

		tagID, err = createClusterTagID(ctx, vctrSession, clusterID.InfraID)
		if err != nil {
			return err
		}

		for i := range in.MachineManifests {
			if vm, ok := in.MachineManifests[i].(*v1beta1.VSphereMachine); ok {
				if vm.Spec.Server == server {
					vm.Spec.TagIDs = append(vm.Spec.TagIDs, tagID)
				}
			}
		}

		for _, failureDomain := range installConfig.Config.VSphere.FailureDomains {
			if failureDomain.Server != server {
				continue
			}
			if err = initializeFoldersAndTemplates(ctx, in.RhcosImage, installConfig, vctrSession, clusterID.InfraID, server, tagID); err != nil {
				return fmt.Errorf("unable to initialize folders and templates: %w", err)
			}
		}
	}

	return nil
}
