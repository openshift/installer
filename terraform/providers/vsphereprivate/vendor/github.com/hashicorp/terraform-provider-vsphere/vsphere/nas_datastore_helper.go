package vsphere

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/hostsystem"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/viapi"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
)

// nasDatastoreMountProcessor is an object that wraps the "complex" mounting
// and unmounting workflows in the NAS datastore resource. We are using an
// object as the process itself is a bit too complex for a pure functional
// approach.
type nasDatastoreMountProcessor struct {
	// The client connection.
	client *govmomi.Client

	// A list of old (current) hosts mounted to the datastore.
	oldHSIDs []string

	// The list of hosts that should be mounted to the datastore.
	newHSIDs []string

	// The NAS datastore volume spec, used for mounting new hosts to a datastore.
	volSpec *types.HostNasVolumeSpec

	// The datastore. If this is not populated by the time the first host is
	// mounted, it's assumed that the datastore is new and we populate this field
	// with that newly created datastore. If this is missing, unmount operations
	// will also be skipped.
	ds *object.Datastore
}

// diffOldNew returns any elements of old that were missing in new.
func (p *nasDatastoreMountProcessor) diffOldNew() []string {
	return p.diff(p.oldHSIDs, p.newHSIDs)
}

// diffNewOld returns any elements of new that were missing in old.
func (p *nasDatastoreMountProcessor) diffNewOld() []string {
	return p.diff(p.newHSIDs, p.oldHSIDs)
}

// diff is what diffOldNew and diffNewOld hand off to.
func (p *nasDatastoreMountProcessor) diff(a, b []string) []string {
	var found bool
	c := make([]string, 0)
	for _, v1 := range a {
		for _, v2 := range b {
			if v1 == v2 {
				found = true
			}
		}
		if !found {
			c = append(c, v1)
		}
	}
	return c
}

// processMountOperations processes all pending mount operations by diffing old
// and new and adding any hosts that were not found in old. The datastore is
// returned, along with any error.
func (p *nasDatastoreMountProcessor) processMountOperations() (*object.Datastore, error) {
	hosts := p.diffNewOld()
	if len(hosts) < 1 {
		// Nothing to do
		return p.ds, nil
	}
	// Validate we are vCenter if we are working with multiple hosts
	if len(hosts) > 1 {
		if err := viapi.ValidateVirtualCenter(p.client); err != nil {
			return p.ds, fmt.Errorf("cannot mount on multiple hosts: %s", err)
		}
	}
	for _, hsID := range hosts {
		dss, err := hostDatastoreSystemFromHostSystemID(p.client, hsID)
		if err != nil {
			return p.ds, fmt.Errorf("host %q: %s", hostsystem.NameOrID(p.client, hsID), err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
		defer cancel()
		ds, err := dss.CreateNasDatastore(ctx, *p.volSpec)
		if err != nil {
			return p.ds, fmt.Errorf("host %q: %s", hostsystem.NameOrID(p.client, hsID), err)
		}
		if err := p.validateDatastore(ds); err != nil {
			return p.ds, fmt.Errorf("datastore validation error on host %q: %s", hostsystem.NameOrID(p.client, hsID), err)
		}
	}
	return p.ds, nil
}

// processUnmountOperations processes all pending unmount operations by diffing old
// and new and removing any hosts that were not found in new. This operation
// only proceeds if the datastore field in the processor is populated.
func (p *nasDatastoreMountProcessor) processUnmountOperations() error {
	hosts := p.diffOldNew()
	if len(hosts) < 1 || p.ds == nil {
		// Nothing to do
		return nil
	}
	for _, hsID := range hosts {
		dss, err := hostDatastoreSystemFromHostSystemID(p.client, hsID)
		if err != nil {
			return fmt.Errorf("host %q: %s", hostsystem.NameOrID(p.client, hsID), err)
		}
		if err := removeDatastore(dss, p.ds); err != nil {
			return fmt.Errorf("host %q: %s", hostsystem.NameOrID(p.client, hsID), err)
		}
	}
	return nil
}

// validateDatastore does one of two things: either stores the current
// datastore in the processor, if it's missing, or validates the supplied
// datastore with the one currently in the processor by checking if their IDs
// match.
func (p *nasDatastoreMountProcessor) validateDatastore(ds *object.Datastore) error {
	if p.ds == nil {
		p.ds = ds
		return nil
	}
	expected := p.ds.Reference().Value
	actual := ds.Reference().Value
	if expected != actual {
		return fmt.Errorf("expected datastore ID to be %q, got %q", expected, actual)
	}
	return nil
}
