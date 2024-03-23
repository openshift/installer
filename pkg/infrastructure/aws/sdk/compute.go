package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/sirupsen/logrus"
)

type computeInputOptions struct {
	infraID            string
	partitionDNSSuffix string
	tags               map[string]string
}

func createComputeResources(ctx context.Context, logger logrus.FieldLogger, iamClient iamiface.IAMAPI, input *computeInputOptions) error {
	profileName := fmt.Sprintf("%s-worker", input.infraID)
	_, err := createComputeInstanceProfile(ctx, logger, iamClient, profileName, input.partitionDNSSuffix, input.tags)
	if err != nil {
		return fmt.Errorf("failed to create compute instance profile: %w", err)
	}
	logger.Infoln("Created compute instance profile")
	return nil
}

func createComputeInstanceProfile(ctx context.Context, logger logrus.FieldLogger, client iamiface.IAMAPI, name string, partitionDNSSuffix string, tags map[string]string) (*iam.InstanceProfile, error) {
	assumeRolePolicy := fmt.Sprintf(`{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "sts:AssumeRole",
            "Principal": {
                "Service": "ec2.%s"
            },
            "Effect": "Allow",
            "Sid": ""
        }
    ]
}`, partitionDNSSuffix)

	const policy = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
		"ec2:DescribeInstances",
		"ec2:DescribeRegions"
      ],
      "Resource": "*",
      "Effect": "Allow"
    }
  ]
}`

	input := &instanceProfileOptions{
		namePrefix:       name,
		assumeRolePolicy: assumeRolePolicy,
		policyDocument:   policy,
		tags:             tags,
	}
	return createInstanceProfile(ctx, logger, client, input)
}
