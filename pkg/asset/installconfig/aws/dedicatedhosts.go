package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/sirupsen/logrus"
)

// Host holds metadata for a dedicated host.
type Host struct {
	ID   string
	Zone string
}

// dedicatedHosts retrieves a list of dedicated hosts for the given region and
// returns them in a map keyed by the host ID.
func dedicatedHosts(ctx context.Context, session *session.Session, region string) (map[string]Host, error) {
	hostsByID := map[string]Host{}

	client := ec2.New(session, aws.NewConfig().WithRegion(region))
	input := &ec2.DescribeHostsInput{}

	if err := client.DescribeHostsPagesWithContext(ctx, input, func(page *ec2.DescribeHostsOutput, lastPage bool) bool {
		for _, h := range page.Hosts {
			id := aws.StringValue(h.HostId)
			if id == "" {
				// Skip entries lacking an ID (should not happen)
				continue
			}

			logrus.Debugf("Found dedicatd host: %s", id)
			hostsByID[id] = Host{
				ID:   id,
				Zone: aws.StringValue(h.AvailabilityZone),
			}
		}
		return !lastPage
	}); err != nil {
		return nil, fmt.Errorf("fetching dedicated hosts: %w", err)
	}

	return hostsByID, nil
}
