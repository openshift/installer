package azure

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest/to"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-10-01/compute"
	azureutil "github.com/openshift/installer/pkg/asset/installconfig/azure"
	"github.com/openshift/installer/pkg/types/azure"
)

// AvailabilityZones retrieves a list of availability zones for the given cloud, region, and instance type.
func AvailabilityZones(cloud azure.CloudEnvironment, region string, instanceType string) ([]string, error) {
	skusClient, err := skusClient(cloud)
	if err != nil {
		return nil, err
	}
	zones, err := fetchAvailabilityZones(skusClient, region, instanceType)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch availability zones: %v", err)
	}
	return zones, nil
}

func skusClient(cloudName azure.CloudEnvironment) (client *compute.ResourceSkusClient, err error) {
	ssn, err := azureutil.GetSession(cloudName)
	if err != nil {
		return nil, err
	}

	skusClient := compute.NewResourceSkusClient(ssn.Credentials.SubscriptionID)
	skusClient.Authorizer = ssn.Authorizer
	return &skusClient, nil
}

func fetchAvailabilityZones(client *compute.ResourceSkusClient, region string, instanceType string) ([]string, error) {
	var zones []string
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	for res, err := client.List(ctx); res.NotDone(); err = res.NextWithContext(ctx) {
		if err != nil {
			return zones, err
		}

		for _, resSku := range res.Values() {
			if strings.EqualFold(to.String(resSku.Name), instanceType) {
				for _, locationInfo := range *resSku.LocationInfo {
					if strings.EqualFold(to.String(locationInfo.Location), region) {
						zones = *locationInfo.Zones
					}
				}
			}
		}
	}

	return zones, nil
}
