package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// DescribeSecurityGroups returns the list of ec2 Security Groups that contain the group id and vpc id.
func DescribeSecurityGroups(ctx context.Context, session *session.Session, securityGroupIDs []string, region string) ([]*ec2.SecurityGroup, error) {
	client := ec2.New(session, aws.NewConfig().WithRegion(region))

	sgIDPtrs := []*string{}
	for _, sgid := range securityGroupIDs {
		sgid := sgid
		sgIDPtrs = append(sgIDPtrs, &sgid)
	}

	cctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	sgOutput, err := client.DescribeSecurityGroupsWithContext(cctx, &ec2.DescribeSecurityGroupsInput{GroupIds: sgIDPtrs})
	if err != nil {
		return nil, err
	}
	return sgOutput.SecurityGroups, nil
}
