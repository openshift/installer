package v1alpha1

import (
	"fmt"
	"net"
	"regexp"

	"github.com/google/uuid"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/metal3-io/baremetal-operator/pkg/hardwareutils/bmc"
)

// log is for logging in this package.
var log = logf.Log.WithName("baremetalhost-validation")

// validateHost validates BareMetalHost resource for creation
func (host *BareMetalHost) validateHost() []error {
	log.Info("validate create", "name", host.Name)
	var errs []error
	var bmcAccess bmc.AccessDetails

	if host.Spec.BMC.Address != "" {
		var err error
		bmcAccess, err = bmc.NewAccessDetails(host.Spec.BMC.Address, host.Spec.BMC.DisableCertificateVerification)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if raid_errors := validateRAID(host.Spec.RAID); raid_errors != nil {
		errs = append(errs, raid_errors...)
	}

	errs = append(errs, validateBMCAccess(host.Spec, bmcAccess)...)

	if err := validateBMHName(host.Name); err != nil {
		errs = append(errs, err)
	}

	return errs
}

// validateChanges validates BareMetalHost resource on changes
// but also covers the validations of creation
func (host *BareMetalHost) validateChanges(old *BareMetalHost) []error {
	log.Info("validate update", "name", host.Name)
	var errs []error

	if err := host.validateHost(); err != nil {
		errs = append(errs, err...)
	}

	if old.Spec.BMC.Address != "" && host.Spec.BMC.Address != old.Spec.BMC.Address {
		errs = append(errs, fmt.Errorf("BMC address can not be changed once it is set"))
	}

	if old.Spec.BootMACAddress != "" && host.Spec.BootMACAddress != old.Spec.BootMACAddress {
		errs = append(errs, fmt.Errorf("bootMACAddress can not be changed once it is set"))
	}

	return errs
}

func validateBMCAccess(s BareMetalHostSpec, bmcAccess bmc.AccessDetails) []error {
	var errs []error

	if bmcAccess == nil {
		return errs
	}

	if s.RAID != nil && len(s.RAID.HardwareRAIDVolumes) > 0 {
		if bmcAccess.RAIDInterface() == "no-raid" {
			errs = append(errs, fmt.Errorf("BMC driver %s does not support configuring RAID", bmcAccess.Type()))
		}
	}

	if s.Firmware != nil {
		if _, err := bmcAccess.BuildBIOSSettings((*bmc.FirmwareConfig)(s.Firmware)); err != nil {
			errs = append(errs, err)
		}
	}

	if bmcAccess.NeedsMAC() && s.BootMACAddress == "" {
		errs = append(errs, fmt.Errorf("BMC driver %s requires a BootMACAddress value", bmcAccess.Type()))
	}

	if s.BootMACAddress != "" {
		_, err := net.ParseMAC(s.BootMACAddress)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if s.BootMode == UEFISecureBoot && !bmcAccess.SupportsSecureBoot() {
		errs = append(errs, fmt.Errorf("BMC driver %s does not support secure boot", bmcAccess.Type()))
	}

	return errs
}

func validateRAID(r *RAIDConfig) []error {
	var errors []error

	if r == nil {
		return nil
	}

	// check if both hardware and software RAID are specified
	if len(r.HardwareRAIDVolumes) > 0 && len(r.SoftwareRAIDVolumes) > 0 {
		errors = append(errors, fmt.Errorf("hardwareRAIDVolumes and softwareRAIDVolumes can not be set at the same time"))
	}

	for index, volume := range r.HardwareRAIDVolumes {
		// check if physicalDisks are specified without a controller
		if len(volume.PhysicalDisks) != 0 {
			if volume.Controller == "" {
				errors = append(errors, fmt.Errorf("'physicalDisks' specified without 'controller' in hardware RAID volume %d", index))
			}
		}
		// check if numberOfPhysicalDisks is not same as len(physicalDisks)
		if volume.NumberOfPhysicalDisks != nil && len(volume.PhysicalDisks) != 0 {
			if *volume.NumberOfPhysicalDisks != len(volume.PhysicalDisks) {
				errors = append(errors, fmt.Errorf("the 'numberOfPhysicalDisks'[%d] and number of 'physicalDisks'[%d] is not same for volume %d", *volume.NumberOfPhysicalDisks, len(volume.PhysicalDisks), index))
			}
		}
	}

	return errors
}

func validateBMHName(bmhname string) error {

	invalidname, _ := regexp.MatchString(`[^A-Za-z0-9\.\-\_]`, bmhname)
	if invalidname {
		return fmt.Errorf("BareMetalHost resource name cannot contain characters other than [A-Za-z0-9._-]")
	}

	_, err := uuid.Parse(bmhname)
	if err == nil {
		return fmt.Errorf("BareMetalHost resource name cannot be a UUID")
	}

	return nil
}
