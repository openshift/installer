// Package AWS provides AWS-specific tools for gathering debugging information.
package aws

import (
	"context"
	"encoding/base64"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
)

// ConsoleLogs retrieves console logs from the AWS instance with the
// given IP address.
func ConsoleLogs(ctx context.Context, session *session.Session, ip string) ([]byte, error) {
	client := ec2.New(session)
	var instanceID string
	err := client.DescribeInstancesPagesWithContext(
		ctx,
		&ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{{
				Name:   aws.String("ip-address"),
				Values: []*string{&ip},
			}},
		},
		func(results *ec2.DescribeInstancesOutput, lastPage bool) bool {
			for _, reservation := range results.Reservations {
				for _, instance := range reservation.Instances {
					if instance.InstanceId != nil {
						instanceID = *instance.InstanceId
						return false
					}
				}
			}

			return !lastPage
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "describe instances")
	}

	if instanceID == "" {
		return nil, errors.Errorf("unable to find an AWS instance ID for %q", ip)
	}

	consoleOutput, err := client.GetConsoleOutputWithContext(
		ctx,
		&ec2.GetConsoleOutputInput{
			InstanceId: &instanceID,
			Latest:     aws.Bool(true),
		},
	)
	if err != nil {
		return nil, errors.Wrapf(err, "get console output for %s", instanceID)
	}
	if consoleOutput.Output == nil {
		return nil, errors.Errorf("nil console output for %s", instanceID)
	}

	data, err := base64.StdEncoding.DecodeString(*consoleOutput.Output)
	if err != nil {
		return nil, errors.Wrapf(err, "decoding console output for %s", instanceID)
	}

	return data, nil
}
