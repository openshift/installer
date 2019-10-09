package gcp

import (
	"context"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
)

// Validator provides validation which requires access to the GCP API.
type Validator struct {
	Client        API
	InstallConfig *types.InstallConfig
}

// Validate executes platform-specific validation.
func (v *Validator) Validate() field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, v.validateVPC()...)

	return allErrs
}

// validateVPC checks that the user-provided VPC is in the project and the provided subnets are valid.
func (v *Validator) validateVPC() field.ErrorList {
	allErrs := field.ErrorList{}

	network, err := v.Client.GetNetwork(context.TODO(), v.InstallConfig.GCP.Network, v.InstallConfig.GCP.ProjectID)
	if err != nil {
		return append(allErrs, field.Invalid(field.NewPath("network"), v.InstallConfig.GCP.Network, err.Error()))
	}

	if ok, errMsg := validateSubnet(network.Subnetworks, v.InstallConfig.GCP.ControlPlaneSubnet, v.InstallConfig.GCP.Region); !ok {
		allErrs = append(allErrs, field.Invalid(field.NewPath("controlPlaneSubnet"), v.InstallConfig.GCP.ControlPlaneSubnet, errMsg))
	}

	if ok, errMsg := validateSubnet(network.Subnetworks, v.InstallConfig.GCP.ComputeSubnet, v.InstallConfig.GCP.Region); !ok {
		allErrs = append(allErrs, field.Invalid(field.NewPath("computeSubnet"), v.InstallConfig.GCP.ComputeSubnet, errMsg))
	}

	return allErrs
}

// validateSubnets checks that the subnets are in the provided VPC and in the cluster region.
func validateSubnet(subnets []string, userSubnet, region string) (bool, string) {
	for _, vpcSubnet := range subnets {
		if userSubnet == name(vpcSubnet) {
			if !strings.Contains(vpcSubnet, region) {
				return false, fmt.Sprintf("subnet %s not in region %s", vpcSubnet, region)
			}
			return true, ""
		}
	}
	return false, fmt.Sprintf("could not find subnet %s", userSubnet)
}

// name takes the fully-qualified URL of the subnetwork, which is provided in the API call and returns the name.
// Example url: https://www.googleapis.com/compute/v1/projects/openshift-dev-installer/regions/us-east4/subnetworks/byovpc-master-subnet
// Returns: byovpc-master-subnet
func name(url string) string {
	splitSubnet := strings.Split(url, "/")
	return splitSubnet[len(splitSubnet)-1]
}
