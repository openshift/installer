package azure

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-10-01/compute"
	"github.com/Azure/go-autorest/autorest/to"

	"github.com/openshift/installer/pkg/client/azure"
)

// AvailabilityZones retrieves a list of availability zones for the given cloud, region, and instance type.
func AvailabilityZones(session *azure.Session, region string, instanceType string) ([]string, error) {
	skusClient, err := skusClient(session)
	if err != nil {
		return nil, err
	}
	zones, err := fetchAvailabilityZones(skusClient, region, instanceType)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch availability zones: %v", err)
	}
	return zones, nil
}

func skusClient(session *azure.Session) (client *compute.ResourceSkusClient, err error) {
	skusClient := compute.NewResourceSkusClientWithBaseURI(session.Environment.ResourceManagerEndpoint, session.Credentials.SubscriptionID)
	skusClient.Authorizer = session.Authorizer
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
