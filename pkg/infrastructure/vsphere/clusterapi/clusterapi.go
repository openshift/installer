package clusterapi

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/vmware/govmomi/object"
	"sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/session"

	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	"github.com/openshift/installer/pkg/rhcos/cache"
	"github.com/openshift/installer/pkg/types/vsphere"
)

// Provider is the vSphere implementation of the clusterapi InfraProvider.
type Provider struct {
	clusterapi.InfraProvider
}

var _ clusterapi.PreProvider = (*Provider)(nil)

// Name returns the vsphere provider name.
func (p Provider) Name() string {
	return vsphere.Name
}

// PublicGatherEndpoint indicates that machine ready checks should NOT wait for an ExternalIP
// in the status when declaring machines ready.
func (Provider) PublicGatherEndpoint() clusterapi.GatherEndpoint { return clusterapi.InternalIP }

func initializeFoldersAndTemplates(ctx context.Context, cachedImage string, failureDomain vsphere.FailureDomain, session *session.Session, diskType vsphere.DiskType, clusterID, tagID string) error {
	finder := session.Finder

	dc, err := finder.Datacenter(ctx, failureDomain.Topology.Datacenter)
	if err != nil {
		logrus.Errorf("unable to get datacenter: %v", err)
		return nil
	}

	// Upstream govmomi bug, workaround
	// https://github.com/vmware/govmomi/issues/3523
	folderPath := path.Join(dc.InventoryPath, "vm", clusterID)

	// we must set the Folder to the infraId somewhere, we will need to remove that.
	// if we are overwriting folderPath it needs to have a slash (path)
	folder := failureDomain.Topology.Folder
	if strings.Contains(folder, "/") {
		folderPath = folder
	}

	var folderObj *object.Folder

	// Only createFolder() and attach the tag if the folder does not exist prior to installing
	if folderObj, err = folderExists(ctx, folderPath, session); folderObj == nil && err == nil {
		folderObj, err = createFolder(ctx, folderPath, session)
		if err != nil {
			logrus.Errorf("unable to create folder: %v", err)
			return nil
		}
		// attach tag to folder
		err = session.TagManager.AttachTag(ctx, tagID, folderObj.Reference())
		if err != nil {
			logrus.Errorf("unable to attach tag to folder: %v", err)
			return nil
		}
	} else if err != nil {
		logrus.Errorf("unable to get folder: %v", err)
		return nil
	}

	// if the template is empty, the ova must be imported
	if len(failureDomain.Topology.Template) == 0 {
		if err = importRhcosOva(ctx, session, folderObj,
			cachedImage, clusterID, tagID, string(diskType), failureDomain); err != nil {
			logrus.Errorf("failed to import ova: %v", err)
			return nil
		}
	}
	return nil
}

// PreProvision creates the vCenter objects required prior to running capv.
func (p Provider) PreProvision(ctx context.Context, in clusterapi.PreProvisionInput) error {
	/*
	 * one locally cached image
	 * one tag and tag category per vcenter
	 * one folder per datacenter
	 * one template per region/zone aka failuredomain
	 */

	installConfig := in.InstallConfig
	clusterID := &installconfig.ClusterID{InfraID: in.InfraID}
	var tagID, cachedImage string
	var err error

	if downloadImage(installConfig) {
		cachedImage, err = cache.DownloadImageFile(in.RhcosImage.ControlPlane, cache.InstallerApplicationName)
		if err != nil {
			logrus.Errorf("failed to use cached vsphere image: %v", err)
			return nil
		}
	}

	for _, vcenter := range installConfig.Config.VSphere.VCenters {
		server := vcenter.Server
		vctrSession, err := installConfig.VSphere.Session(context.TODO(), server)

		if err != nil {
			logrus.Errorf("unable to get vCenter session: %v", err)
			return nil
		}

		tagID, err = createClusterTagID(ctx, vctrSession, clusterID.InfraID)
		if err != nil {
			logrus.Errorf("unable to create cluster tag ID: %v", err)
			return nil
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

			if failureDomain.ZoneType == vsphere.HostGroupFailureDomain {
				vmGroupAndRuleName := fmt.Sprintf("%s-%s", clusterID.InfraID, failureDomain.Name)

				err = createVMGroup(ctx, vctrSession, failureDomain.Topology.ComputeCluster, vmGroupAndRuleName)
				if err != nil {
					logrus.Errorf("unable to create VM group: %v", err)
					return nil
				}

				err = createVMHostAffinityRule(ctx, vctrSession, failureDomain.Topology.ComputeCluster, failureDomain.Topology.HostGroup, vmGroupAndRuleName, vmGroupAndRuleName)
				if err != nil {
					logrus.Errorf("unable to create VM host affinity rule: %v", err)
					return nil
				}
			}

			if err = initializeFoldersAndTemplates(ctx, cachedImage, failureDomain, vctrSession, installConfig.Config.VSphere.DiskType, clusterID.InfraID, tagID); err != nil {
				logrus.Errorf("unable to initialize folders and templates: %v", err)
				return nil
			}
		}
	}

	return nil
}

// downloadImage if any failure domains don't have a defined template, this function
// returns true.
func downloadImage(installConfig *installconfig.InstallConfig) bool {
	for _, fd := range installConfig.Config.VSphere.FailureDomains {
		if fd.Topology.Template == "" {
			return true
		}
	}
	return false
}
