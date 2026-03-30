package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/sirupsen/logrus"
)

// Host holds metadata for a dedicated host.
type Host struct {
	// ID the ID of the host.
	ID string
	// Zone the zone the host belongs to.
	Zone string
	// Tags is the map of the Host's tags.
	Tags Tags
}

// dedicatedHosts retrieves a list of dedicated hosts and returns them in a map keyed by the host ID.
func dedicatedHosts(ctx context.Context, client *ec2.Client, hosts []string) (map[string]Host, error) {
	hostsByID := map[string]Host{}

	input := &ec2.DescribeHostsInput{}
	if len(hosts) > 0 {
		input.HostIds = hosts
	}

	paginator := ec2.NewDescribeHostsPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("fetching dedicated hosts: %w", err)
		}

		for _, host := range page.Hosts {
			id := aws.ToString(host.HostId)
			if id == "" {
				// Skip entries lacking an ID (should not happen)
				continue
			}

			logrus.Debugf("Found dedicated host: %s", id)
			hostsByID[id] = Host{
				ID:   id,
				Zone: aws.ToString(host.AvailabilityZone),
				Tags: FromAWSTags(host.Tags),
			}
		}
	}

	return hostsByID, nil
}
