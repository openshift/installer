package v1alpha1

import (
	"fmt"
	"regexp"

	"github.com/google/uuid"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	_ "github.com/metal3-io/baremetal-operator/pkg/hardwareutils/bmc"
)

// log is for logging in this package.
var log = logf.Log.WithName("baremetalhost-validation")

// validateHost validates BareMetalHost resource for creation
func (host *BareMetalHost) validateHost() []error {
	log.Info("validate create", "name", host.Name)
	var errs []error

	if err := validateRAID(host.Spec.RAID); err != nil {
		errs = append(errs, err)
	}

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

func validateRAID(r *RAIDConfig) error {
	if r == nil {
		return nil
	}

	if len(r.HardwareRAIDVolumes) > 0 && len(r.SoftwareRAIDVolumes) > 0 {
		return fmt.Errorf("hardwareRAIDVolumes and softwareRAIDVolumes can not be set at the same time")
	}

	return nil
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
