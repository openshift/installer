package nutanix

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	nutanixtypes "github.com/openshift/installer/pkg/types/nutanix"
)

// Validate executes platform-specific validation.
func Validate(ic *types.InstallConfig) error {
	if ic.Platform.Nutanix == nil {
		return field.Required(field.NewPath("platform", "nutanix"), "nutanix validation requires a nutanix platform configuration")
	}

	return nil
}

// ValidateForProvisioning performs platform validation specifically for installer-
// provisioned infrastructure. In this case, self-hosted networking is a requirement
// when the installer creates infrastructure for nutanix clusters.
func ValidateForProvisioning(ic *types.InstallConfig) error {
	errList := field.ErrorList{}
	parentPath := field.NewPath("platform", "nutanix")

	if ic.Platform.Nutanix == nil {
		return field.Required(parentPath, "nutanix validation requires a nutanix platform configuration")
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()
	p := ic.Platform.Nutanix
	nc, err := nutanixtypes.CreateNutanixClient(ctx,
		p.PrismCentral.Endpoint.Address,
		strconv.Itoa(int(p.PrismCentral.Endpoint.Port)),
		p.PrismCentral.Username,
		p.PrismCentral.Password)
	if err != nil {
		return field.Invalid(parentPath.Child("prismCentral"), p.PrismCentral.Endpoint.Address, fmt.Sprintf("failed to connect to the prism-central with the configured credentials: %v", err))
	}

	// validate whether a prism element with the UUID actually exists
	for _, pe := range p.PrismElements {
		_, err = nc.V3.GetCluster(ctx, pe.UUID)
		if err != nil {
			errList = append(errList, field.Invalid(parentPath.Child("prismElements"), pe.UUID, fmt.Sprintf("the prism element %s's UUID does not correspond to a valid prism element in Prism: %v", pe.Name, err)))
		}
	}

	// validate whether a subnet with the UUID actually exists
	for _, subnetUUID := range p.SubnetUUIDs {
		_, err = nc.V3.GetSubnet(ctx, subnetUUID)
		if err != nil {
			errList = append(errList, field.Invalid(parentPath.Child("subnetUUIDs"), subnetUUID, fmt.Sprintf("the subnet UUID does not correspond to a valid subnet in Prism: %v", err)))
		}
	}

	// validate each FailureDomain configuration
	for _, fd := range p.FailureDomains {
		// validate whether the prism element with the UUID exists
		_, err = nc.V3.GetCluster(ctx, fd.PrismElement.UUID)
		if err != nil {
			errList = append(errList, field.Invalid(parentPath.Child("failureDomains", "prismElements"), fd.PrismElement.UUID,
				fmt.Sprintf("the failure domain %s configured prism element UUID does not correspond to a valid prism element in Prism: %v", fd.Name, err)))
		}

		// validate whether a subnet with the UUID actually exists
		for _, subnetUUID := range fd.SubnetUUIDs {
			_, err = nc.V3.GetSubnet(ctx, subnetUUID)
			if err != nil {
				errList = append(errList, field.Invalid(parentPath.Child("failureDomains", "subnetUUIDs"), subnetUUID,
					fmt.Sprintf("the failure domain %s configured subnet UUID does not correspond to a valid subnet in Prism: %v", fd.Name, err)))
			}
		}
	}

	return errList.ToAggregate()
}
