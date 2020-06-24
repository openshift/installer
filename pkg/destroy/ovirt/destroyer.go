package ovirt

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ovirt/go-ovirt"
	"github.com/sirupsen/logrus"

	"k8s.io/apimachinery/pkg/util/errors"

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
func (uninstaller *ClusterUninstaller) Run(ctx context.Context) error {
	con, err := ovirt.NewConnection()
	if err != nil {
		return fmt.Errorf("failed to initialize connection to ovirt-engine's %s", err)
	}
	defer con.Close()

	// Tags
	tagVMs := uninstaller.Metadata.InfraID
	tagVMbootstrap := uninstaller.Metadata.InfraID + "-bootstrap"
	tags := [2]string{tagVMs, tagVMbootstrap}

	for _, tag := range tags {
		if err := uninstaller.removeVMs(ctx, con, tag); err != nil {
			uninstaller.Logger.Errorf("failed to remove VMs: %s", err)
		}
		if err := uninstaller.removeTag(ctx, con, tag); err != nil {
			uninstaller.Logger.Errorf("failed to remove tag: %s", err)
		}
	}
	if err := uninstaller.removeTemplate(ctx, con); err != nil {
		uninstaller.Logger.Errorf("Failed to remove template: %s", err)
	}

	return nil
}

func (uninstaller *ClusterUninstaller) removeVMs(ctx context.Context, con *ovirtsdk.Connection, tag string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	// - find all vms by tag name=infraID
	vmsService := con.SystemService().VmsService()
	searchTerm := fmt.Sprintf("tag=%s", tag)
	uninstaller.Logger.Debugf("Searching VMs by %s", searchTerm)
	vmsResponse, err := vmsService.List().Search(searchTerm).Send()
	if err != nil {
		return err
	}
	// - stop + delete VMS
	vms := vmsResponse.MustVms().Slice()
	uninstaller.Logger.Debugf("Found %d VMs", len(vms))
	errChan := make(chan (error), len(vms))
	wg := sync.WaitGroup{}
	wg.Add(len(vms))
	for _, vm := range vms {
		go func(vm *ovirtsdk.Vm) {
			defer wg.Done()
			if err := uninstaller.stopVM(ctx, vmsService, vm); err != nil {
				errChan <- err
				return
			}
			if err := uninstaller.removeVM(vmsService, vm); err != nil {
				errChan <- err
			}
		}(vm)
	}
	wg.Wait()
	close(errChan)
	var errorList []error
	for err := range errChan {
		errorList = append(errorList, err)
	}
	return errors.NewAggregate(errorList)
}

func (uninstaller *ClusterUninstaller) removeTag(ctx context.Context, con *ovirtsdk.Connection, tag string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	// finally remove the tag
	tagsService := con.SystemService().TagsService()
	tagsServiceListResponse, err := tagsService.List().Send()
	if err != nil {
		return err
	}
	if tagsServiceListResponse != nil {
		for _, t := range tagsServiceListResponse.MustTags().Slice() {
			if t.MustName() == tag {
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

func (uninstaller *ClusterUninstaller) stopVM(ctx context.Context, vmsService *ovirtsdk.VmsService, vm *ovirtsdk.Vm) error {
	vmService := vmsService.VmService(vm.MustId())
	// this is a teardown, stopping instead of shutting down.
	if _, err := vmService.Stop().Send(); err != nil {
		uninstaller.Logger.Errorf("Failed to stop VM %s: %s", vm.MustName(), err)
		return err
	}
	uninstaller.Logger.Infof("Stopping VM %s", vm.MustName())
	waitForDownDuration := timeoutWithContext(ctx, time.Minute*10)
	if err := vmService.Connection().WaitForVM(vm.MustId(), ovirtsdk.VMSTATUS_DOWN, waitForDownDuration); err != nil {
		uninstaller.Logger.Warnf("Waited %d for VM %s to power off: %s", waitForDownDuration, vm.MustName(), err)
		return err
	}
	uninstaller.Logger.Infof("VM %s powered off", vm.MustName())
	return nil
}

func (uninstaller *ClusterUninstaller) removeVM(vmsService *ovirtsdk.VmsService, vm *ovirtsdk.Vm) error {
	vmService := vmsService.VmService(vm.MustId())
	if _, err := vmService.Remove().Send(); err != nil {
		uninstaller.Logger.Errorf("Failed to remove VM %s: %s", vm.MustName(), err)
		return err
	}
	uninstaller.Logger.Infof("Removing VM %s", vm.MustName())
	return nil
}

func (uninstaller *ClusterUninstaller) removeTemplate(ctx context.Context, con *ovirtsdk.Connection) error {
	if uninstaller.Metadata.Ovirt.RemoveTemplate {
		if err := ctx.Err(); err != nil {
			return err
		}
		search, err := con.SystemService().TemplatesService().
			List().Search(fmt.Sprintf("name=%s-rhcos", uninstaller.Metadata.InfraID)).Send()
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

// timeoutWithContext returns the smaller timeout between the context timeout and the specified timeout.
func timeoutWithContext(ctx context.Context, timeout time.Duration) time.Duration {
	deadline, ok := ctx.Deadline()
	if !ok {
		return timeout
	}
	timeUntilDeadline := time.Until(deadline)
	if timeUntilDeadline < timeout {
		return timeUntilDeadline
	}
	return timeout
}
