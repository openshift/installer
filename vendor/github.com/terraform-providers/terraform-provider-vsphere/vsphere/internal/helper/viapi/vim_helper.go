package viapi

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/provider"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/task"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
)

// ManagedObject allows for working generically with managed objects.
type ManagedObject interface {
	Reference() types.ManagedObjectReference
}

// ErrVirtualCenterOnly is the error message that validateVirtualCenter returns.
const ErrVirtualCenterOnly = "this operation is only supported on vCenter"

// soapFault extracts the SOAP fault from an error fault, if it exists. Check
// the returned boolean value to see if you have a SoapFault.
func soapFault(err error) (*soap.Fault, bool) {
	if soap.IsSoapFault(err) {
		return soap.ToSoapFault(err), true
	}
	return nil, false
}

// vimSoapFault extracts the VIM fault. Check the returned boolean value to see
// if you have a fault, which will need to be further asserted into the error
// that you are looking for.
func vimSoapFault(err error) (types.AnyType, bool) {
	if sf, ok := soapFault(err); ok {
		return sf.VimFault(), true
	}
	return nil, false
}

// taskFault extracts the task fault from a supplied task.Error. Check the
// returned boolean value to see if the fault was extracted correctly, after
// which you will need to do further checking.
func taskFault(err error) (types.BaseMethodFault, bool) {
	if te, ok := err.(task.Error); ok {
		return te.Fault(), true
	}
	return nil, false
}

// IsManagedObjectNotFoundError checks an error to see if it's of the
// ManagedObjectNotFound type.
func IsManagedObjectNotFoundError(err error) bool {
	if f, ok := vimSoapFault(err); ok {
		if _, ok := f.(types.ManagedObjectNotFound); ok {
			return true
		}
	}
	return false
}

// IsInvalidStateError checks an error to see if it's of the
// InvalidState type.
func IsInvalidStateError(err error) bool {
	if f, ok := vimSoapFault(err); ok {
		if _, ok := f.(types.InvalidState); ok {
			return true
		}
	}
	return false
}

// IsInvalidPowerStateError checks an error to see if it's of the
// InvalidState type.
func IsInvalidPowerStateError(err error) bool {
	if f, ok := vimSoapFault(err); ok {
		if _, ok := f.(types.InvalidPowerState); ok {
			return true
		}
	}
	return false
}

// isNotFoundError checks an error to see if it's of the NotFoundError type.
//
// Note this is different from the other "not found" faults and is an error
// type in its own right. Use IsAnyNotFoundError to check for any "not found"
// type.
func isNotFoundError(err error) bool {
	if f, ok := vimSoapFault(err); ok {
		if _, ok := f.(types.NotFound); ok {
			return true
		}
	}
	return false
}

// IsAnyNotFoundError checks to see if the fault is of any not found error type
// that we track.
func IsAnyNotFoundError(err error) bool {
	switch {
	case IsManagedObjectNotFoundError(err):
		fallthrough
	case isNotFoundError(err):
		return true
	}
	return false
}

// IsResourceInUseError checks an error to see if it's of the
// ResourceInUse type.
func IsResourceInUseError(err error) bool {
	if f, ok := vimSoapFault(err); ok {
		if _, ok := f.(types.ResourceInUse); ok {
			return true
		}
	}
	return false
}

// isConcurrentAccessError checks an error to see if it's of the
// ConcurrentAccess type.
func isConcurrentAccessError(err error) bool {
	// ConcurrentAccess comes from a task more than it usually does from a direct
	// SOAP call, so we need to handle both here.
	var f types.AnyType
	var ok bool
	f, ok = vimSoapFault(err)
	if !ok {
		f, ok = taskFault(err)
	}
	if ok {
		switch f.(type) {
		case types.ConcurrentAccess, *types.ConcurrentAccess:
			return true
		}
	}
	return false
}

// RenameObject renames a MO and tracks the task to make sure it completes.
func RenameObject(client *govmomi.Client, ref types.ManagedObjectReference, new string) error {
	req := types.Rename_Task{
		This:    ref,
		NewName: new,
	}

	rctx, rcancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer rcancel()
	res, err := methods.Rename_Task(rctx, client.Client, &req)
	if err != nil {
		return err
	}

	t := object.NewTask(client.Client, res.Returnval)
	tctx, tcancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer tcancel()
	return t.Wait(tctx)
}

// ValidateVirtualCenter ensures that the client is connected to vCenter.
func ValidateVirtualCenter(c *govmomi.Client) error {
	return VimValidateVirtualCenter(c.Client)
}

// VimValidateVirtualCenter ensures that the client is connected to vCenter.
//
// This is a lower-level method that does not take the wrapped client from the
// higher-level govmomi object, and can be used to facilitate validation when
// it's not available.
func VimValidateVirtualCenter(c *vim25.Client) error {
	if c.ServiceContent.About.ApiType != "VirtualCenter" {
		return errors.New(ErrVirtualCenterOnly)
	}
	return nil
}

// VSphereVersion represents a version number of a ESXi/vCenter server
// instance.
type VSphereVersion struct {
	// The product name. Example: "VMware vCenter Server", or "VMware ESXi".
	Product string

	// The major version. Example: If "6.5.1" is the full version, the major
	// version is "6".
	Major int

	// The minor version. Example: If "6.5.1" is the full version, the minor
	// version is "5".
	Minor int

	// The patch version. Example: If "6.5.1" is the full version, the patch
	// version is "1".
	Patch int

	// The build number. This is usually a lengthy integer. This number should
	// not be used to compare versions on its own.
	Build int
}

// parseVersion creates a new VSphereVersion from a parsed version string and
// build number.
func parseVersion(name, version, build string) (VSphereVersion, error) {
	v := VSphereVersion{
		Product: name,
	}
	s := strings.Split(version, ".")
	if len(s) > 3 {
		return v, fmt.Errorf("version string %q has more than 3 components", version)
	}
	var err error
	v.Major, err = strconv.Atoi(s[0])
	if err != nil {
		return v, fmt.Errorf("could not parse major version %q from version string %q", s[0], version)
	}
	v.Minor, err = strconv.Atoi(s[1])
	if err != nil {
		return v, fmt.Errorf("could not parse minor version %q from version string %q", s[1], version)
	}
	v.Patch, err = strconv.Atoi(s[2])
	if err != nil {
		return v, fmt.Errorf("could not parse patch version %q from version string %q", s[2], version)
	}
	v.Build, err = strconv.Atoi(build)
	if err != nil {
		return v, fmt.Errorf("could not parse build version string %q", build)
	}

	return v, nil
}

// parseVersionFromAboutInfo returns a populated VSphereVersion from an
// AboutInfo data object.
//
// This function panics if it cannot parse the version correctly, as given our
// source of truth is a valid AboutInfo object, such an error is indicative of
// a major issue with our version parsing logic.
func parseVersionFromAboutInfo(info types.AboutInfo) VSphereVersion {
	v, err := parseVersion(info.Name, info.Version, info.Build)
	if err != nil {
		panic(err)
	}
	return v
}

// ParseVersionFromClient returns a populated VSphereVersion from a client
// connection.
func ParseVersionFromClient(client *govmomi.Client) VSphereVersion {
	return parseVersionFromAboutInfo(client.Client.ServiceContent.About)
}

// String implements stringer for VSphereVersion.
func (v VSphereVersion) String() string {
	return fmt.Sprintf("%s %d.%d.%d build-%d", v.Product, v.Major, v.Minor, v.Patch, v.Build)
}

// ProductEqual returns true if this version's product name is the same as the
// supplied version's name.
func (v VSphereVersion) ProductEqual(other VSphereVersion) bool {
	return v.Product == other.Product
}

// newerVersion checks the major/minor/patch part of the version to see it's
// higher than the version supplied in other. This is broken off from the main
// test so that it can be checked in Older before the build number is compared.
func (v VSphereVersion) newerVersion(other VSphereVersion) bool {
	// Assuming here that VMware is a loooong way away from having a
	// major/minor/patch that's bigger than 254.
	vc := v.Major<<16 + v.Minor<<8 + v.Patch
	vo := other.Major<<16 + other.Minor<<8 + other.Patch
	return vc > vo
}

// Newer returns true if this version's product is the same, and composite of
// the version and build numbers, are newer than the supplied version's
// information.
func (v VSphereVersion) Newer(other VSphereVersion) bool {
	if !v.ProductEqual(other) {
		return false
	}
	if v.newerVersion(other) {
		return true
	}

	// Double check this version is not actually older by version number before
	// moving on to the build number
	if v.olderVersion(other) {
		return false
	}

	if v.Build > other.Build {
		return true
	}
	return false
}

// olderVersion checks the major/minor/patch part of the version to see it's
// older than the version supplied in other. This is broken off from the main
// test so that it can be checked in Newer before the build number is compared.
func (v VSphereVersion) olderVersion(other VSphereVersion) bool {
	// Assuming here that VMware is a loooong way away from having a
	// major/minor/patch that's bigger than 254.
	vc := v.Major<<16 + v.Minor<<8 + v.Patch
	vo := other.Major<<16 + other.Minor<<8 + other.Patch
	return vc < vo
}

// Older returns true if this version's product is the same, and composite of
// the version and build numbers, are older than the supplied version's
// information.
func (v VSphereVersion) Older(other VSphereVersion) bool {
	if !v.ProductEqual(other) {
		return false
	}
	if v.olderVersion(other) {
		return true
	}

	// Double check this version is not actually newer by version number before
	// moving on to the build number
	if v.newerVersion(other) {
		return false
	}

	if v.Build < other.Build {
		return true
	}
	return false
}

// Equal returns true if the version is equal to the supplied version.
func (v VSphereVersion) Equal(other VSphereVersion) bool {
	return v.ProductEqual(other) && !v.Older(other) && !v.Newer(other)
}
