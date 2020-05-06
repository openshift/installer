package ovirt

import (
	"fmt"
	"sync"
	"time"

	"github.com/ovirt/go-ovirt"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset/installconfig/ovirt"
	"github.com/openshift/installer/pkg/destroy/providers"
	"github.com/openshift/installer/pkg/types"
)

// ClusterUninstaller holds the various options for the cluster we want to delete.
type ClusterUninstaller struct {
	Metadata types.ClusterMetadata
	Logger   logrus.FieldLogger
}

// Run is the entrypoint to start the uninstall process.
func (uninstaller *ClusterUninstaller) Run() error {
	con, err := ovirt.NewConnection()
	if err != nil {
		return fmt.Errorf("failed to initialize connection to ovirt-engine's %s", err)
	}
	defer con.Close()

	if err := uninstaller.removeVMs(con); err != nil {
		uninstaller.Logger.Errorf("Failed to remove VMs: %s", err)
	}
	if err := uninstaller.removeTag(con); err != nil {
		uninstaller.Logger.Errorf("Failed to remove tag: %s", err)
	}
	if err := uninstaller.removeTemplate(con); err != nil {
		uninstaller.Logger.Errorf("Failed to remove template: %s", err)
	}

	return nil
}

func (uninstaller *ClusterUninstaller) removeVMs(con *ovirtsdk.Connection) error {
	// - find all vms by tag name=infraID
	vmsService := con.SystemService().VmsService()
	searchTerm := fmt.Sprintf("tag=%s", uninstaller.Metadata.InfraID)
	uninstaller.Logger.Debugf("Searching VMs by %s", searchTerm)
	vmsResponse, err := vmsService.List().Search(searchTerm).Send()
	if err != nil {
		return err
	}
	// - stop + delete VMS
	vms := vmsResponse.MustVms().Slice()
	uninstaller.Logger.Debugf("Found %d VMs", len(vms))
	wg := sync.WaitGroup{}
	wg.Add(len(vms))
	for _, vm := range vms {
		go func(vm *ovirtsdk.Vm) {
			uninstaller.stopVM(vmsService, vm)
			uninstaller.removeVM(vmsService, vm)
			wg.Done()
		}(vm)
	}
	wg.Wait()
	return nil
}

func (uninstaller *ClusterUninstaller) removeTag(con *ovirtsdk.Connection) error {
	// finally remove the tag
	tagsService := con.SystemService().TagsService()
	tagsServiceListResponse, err := tagsService.List().Send()
	if err != nil {
		return err
	}
	if tagsServiceListResponse != nil {
		for _, t := range tagsServiceListResponse.MustTags().Slice() {
			if t.MustName() == uninstaller.Metadata.InfraID {
				uninstaller.Logger.Infof("Removing tag %s", t.MustName())
				_, err := tagsService.TagService(t.MustId()).Remove().Send()
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (uninstaller *ClusterUninstaller) stopVM(vmsService *ovirtsdk.VmsService, vm *ovirtsdk.Vm) {
	vmService := vmsService.VmService(vm.MustId())
	// this is a teardown, stopping instead of shutting down.
	_, err := vmService.Stop().Send()
	if err == nil {
		uninstaller.Logger.Infof("Stopping VM %s", vm.MustName())
	} else {
		uninstaller.Logger.Errorf("Failed to stop VM %s: %s", vm.MustName(), err)
	}
	waitForDownDuration := time.Minute * 10
	err = vmService.Connection().WaitForVM(vm.MustId(), ovirtsdk.VMSTATUS_DOWN, waitForDownDuration)
	if err == nil {
		uninstaller.Logger.Infof("VM %s powered off", vm.MustName())
	} else {
		uninstaller.Logger.Warnf("Waited %d for VM %s to power off: %s", waitForDownDuration, vm.MustName(), err)
	}
}

func (uninstaller *ClusterUninstaller) removeVM(vmsService *ovirtsdk.VmsService, vm *ovirtsdk.Vm) {
	vmService := vmsService.VmService(vm.MustId())
	_, err := vmService.Remove().Send()
	if err == nil {
		uninstaller.Logger.Infof("Removing VM %s", vm.MustName())
	} else {
		uninstaller.Logger.Errorf("Failed to remove VM %s: %s", vm.MustName(), err)
	}
}

func (uninstaller *ClusterUninstaller) removeTemplate(con *ovirtsdk.Connection) error {
	if uninstaller.Metadata.Ovirt.RemoveTemplate {
		search, err := con.SystemService().TemplatesService().
			List().Search(fmt.Sprintf("name=%s", uninstaller.Metadata.InfraID)).Send()
		if err != nil {
			return fmt.Errorf("couldn't find a template with name %s", uninstaller.Metadata.InfraID)
		}
		if result, ok := search.Templates(); ok {
			// the results can potentially return a list of template
			// because the search uses wildcards
			for _, tmp := range result.Slice() {
				uninstaller.Logger.Infof("Removing Template %s", tmp.MustName())
				service := con.SystemService().TemplatesService().TemplateService(tmp.MustId())
				_, err := service.Remove().Send()
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// New returns oVirt Uninstaller from ClusterMetadata.
func New(logger logrus.FieldLogger, metadata *types.ClusterMetadata) (providers.Destroyer, error) {
	return &ClusterUninstaller{
		Metadata: *metadata,
		Logger:   logger,
	}, nil
}
