package nutanix

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"k8s.io/utils/ptr"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types/nutanix"
)

// GetInfrastructureNutanixPlatformSpec constructs NutanixPlatformSpec for the infrastructure spec.
func GetInfrastructureNutanixPlatformSpec(ic *installconfig.InstallConfig) (*configv1.NutanixPlatformSpec, error) {
	nutanixPlatform := ic.Config.Nutanix

	// Retrieve the prism element name
	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	var peName string
	if len(nutanixPlatform.PrismElements[0].Name) == 0 {
		nc, err := nutanix.CreateNutanixClient(ctx,
			nutanixPlatform.PrismCentral.Endpoint.Address,
			strconv.Itoa(int(nutanixPlatform.PrismCentral.Endpoint.Port)),
			nutanixPlatform.PrismCentral.Username,
			nutanixPlatform.PrismCentral.Password)
		if err != nil {
			return nil, fmt.Errorf("unable to connect to Prism Central %s: %w", nutanixPlatform.PrismCentral.Endpoint.Address, err)
		}
		pe, err := nc.V3.GetCluster(ctx, nutanixPlatform.PrismElements[0].UUID)
		if err != nil {
			return nil, fmt.Errorf("fail to find the Prism Element (cluster) with uuid %s: %w", nutanixPlatform.PrismElements[0].UUID, err)
		}
		peName = *pe.Spec.Name
	} else {
		peName = nutanixPlatform.PrismElements[0].Name
	}

	platformSpec := &configv1.NutanixPlatformSpec{
		PrismCentral: configv1.NutanixPrismEndpoint{
			Address: nutanixPlatform.PrismCentral.Endpoint.Address,
			Port:    nutanixPlatform.PrismCentral.Endpoint.Port,
		},
		PrismElements: []configv1.NutanixPrismElementEndpoint{{
			Name: peName,
			Endpoint: configv1.NutanixPrismEndpoint{
				Address: nutanixPlatform.PrismElements[0].Endpoint.Address,
				Port:    nutanixPlatform.PrismElements[0].Endpoint.Port,
			},
		}},
	}

	// failure domains configuration
	failureDomains := make([]configv1.NutanixFailureDomain, 0, len(nutanixPlatform.FailureDomains))
	for _, fd := range nutanixPlatform.FailureDomains {
		subnets := make([]configv1.NutanixResourceIdentifier, 0, len(fd.SubnetUUIDs))
		for _, subnetUUID := range fd.SubnetUUIDs {
			subnets = append(subnets, configv1.NutanixResourceIdentifier{
				Type: configv1.NutanixIdentifierUUID,
				UUID: ptr.To[string](subnetUUID),
			})
		}

		failureDomain := configv1.NutanixFailureDomain{
			Name: fd.Name,
			Cluster: configv1.NutanixResourceIdentifier{
				Type: configv1.NutanixIdentifierUUID,
				UUID: ptr.To[string](fd.PrismElement.UUID),
			},
			Subnets: subnets,
		}
		failureDomains = append(failureDomains, failureDomain)
	}
	platformSpec.FailureDomains = failureDomains

	return platformSpec, nil
}
