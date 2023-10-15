package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/sirupsen/logrus"
)

func createIAMResources(l *logrus.Logger, session *session.Session, infraID string, tags map[string]string) error {
	_, err := CreateComputeInstanceProfile(l, session, infraID, tags)
	return err
}

func CreateComputeInstanceProfile(l *logrus.Logger, session *session.Session, infraID string, additionalTags map[string]string) (string, error) {
	const (
		assumeRolePolicy = `{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "sts:AssumeRole",
            "Principal": {
                "Service": "ec2.amazonaws.com"
            },
            "Effect": "Allow",
            "Sid": ""
        }
    ]
}`
		policy = `{
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
	)

	return createInstanceProfile(l, session, fmt.Sprintf("%s-worker", infraID), assumeRolePolicy, policy, additionalTags)
}
