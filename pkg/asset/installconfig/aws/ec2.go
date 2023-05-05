package aws

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func GetVpcsFromSecurityGroups(ctx context.Context, session *session.Session, securityGroupIDs []string, region string) ([]string, error) {
	client := ec2.New(session, aws.NewConfig().WithRegion(region))

	sgIdPtrs := []*string{}
	for _, sgid := range securityGroupIDs {
		sgIdPtrs = append(sgIdPtrs, &sgid)
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Minute)
	defer cancel()

	sgOutput, err := client.DescribeSecurityGroupsWithContext(ctx, &ec2.DescribeSecurityGroupsInput{GroupIds: sgIdPtrs})
	if err != nil {
		return nil, err
	}

	vpcIds := sets.NewString()
	for _, sg := range sgOutput.SecurityGroups {
		vpcIds.Insert(*sg.VpcId)
	}
	return vpcIds.List(), nil
}
