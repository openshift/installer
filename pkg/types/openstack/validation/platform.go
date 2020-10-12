package validation

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/validate"
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *openstack.Platform, n *types.Networking, fldPath *field.Path, c *types.InstallConfig) field.ErrorList {
	var allErrs field.ErrorList

	allErrs = append(allErrs, validateClusterName(c.ObjectMeta.Name)...)

	for _, ip := range p.ExternalDNS {
		if err := validate.IP(ip); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("externalDNS"), p.ExternalDNS, err.Error()))
		}
	}

	err := validateVIP(p.APIVIP, n)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("apiVIP"), p.APIVIP, err.Error()))
	}

	err = validateVIP(p.IngressVIP, n)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("ingressVIP"), p.IngressVIP, err.Error()))
	}

	// Ensure clusterNetwork doesn't overlap with any machine network
	for _, cn := range n.ClusterNetwork {
		if err := validateNoOverlapMachineCIDR(&cn.CIDR.IPNet, n); err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("networking").Child("clusterNetwork"), cn.CIDR.IPNet.String(), err.Error()))
		}
	}

	return allErrs
}

// validateVIP is a convenience function for validating VIP port and usage
func validateVIP(vip string, n *types.Networking) error {
	if vip != "" {
		if err := validate.IP(vip); err != nil {
			return err
		}

		if !n.MachineNetwork[0].CIDR.Contains(net.ParseIP(vip)) {
			return errors.New("IP is not in the machineNetwork")
		}
	}
	return nil
}

func validateClusterName(name string) (allErrs field.ErrorList) {
	if len(name) > 14 {
		allErrs = append(allErrs, field.Invalid(field.NewPath("metadata", "name"), name, "cluster name is too long, please restrict it to 14 characters"))
	}

	if strings.Contains(name, ".") {
		allErrs = append(allErrs, field.Invalid(field.NewPath("metadata", "name"), name, "cluster name can't contain \".\" character"))
	}

	return
}

// duplicated from baremetal
func validateNoOverlapMachineCIDR(target *net.IPNet, n *types.Networking) error {
	allIPv4 := ipnet.MustParseCIDR("0.0.0.0/0")
	allIPv6 := ipnet.MustParseCIDR("::/0")
	netIsIPv6 := target.IP.To4() == nil

	for _, machineCIDR := range n.MachineNetwork {
		machineCIDRisIPv6 := machineCIDR.CIDR.IP.To4() == nil

		// Only compare if both are the same IP version
		if netIsIPv6 == machineCIDRisIPv6 {
			var err error
			if netIsIPv6 {
				err = cidr.VerifyNoOverlap(
					[]*net.IPNet{
						target,
						&machineCIDR.CIDR.IPNet,
					},
					&allIPv6.IPNet,
				)
			} else {
				err = cidr.VerifyNoOverlap(
					[]*net.IPNet{
						target,
						&machineCIDR.CIDR.IPNet,
					},
					&allIPv4.IPNet,
				)
			}

			if err != nil {
				return fmt.Errorf("%v cannot overlap with machine network", err)
			}
		}
	}

	return nil
}
