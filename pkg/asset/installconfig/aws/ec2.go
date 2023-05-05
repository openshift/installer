package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"k8s.io/apimachinery/pkg/util/sets"
)

// GetVpcsFromSecurityGroups will get the list of VPCs that are attached to the security groups with the ids provided.
func GetVpcsFromSecurityGroups(ctx context.Context, session *session.Session, securityGroupIDs []string, region string) ([]string, error) {
	client := ec2.New(session, aws.NewConfig().WithRegion(region))

	sgIDPtrs := []*string{}
	for _, sgid := range securityGroupIDs {
		sgidAlias := sgid
		sgIDPtrs = append(sgIDPtrs, &sgidAlias)
	}

	cctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	sgOutput, err := client.DescribeSecurityGroupsWithContext(cctx, &ec2.DescribeSecurityGroupsInput{GroupIds: sgIDPtrs})
	if err != nil {
		return nil, err
	}

	vpcIds := sets.New[string]()
	for _, sg := range sgOutput.SecurityGroups {
		vpcIds.Insert(*sg.VpcId)
	}
	return sets.List(vpcIds), nil
}
