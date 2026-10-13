package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

// SimulateParams captures any additional details that should be used
// when simulating permissions.
type SimulateParams struct {
	Region string
}

// ValidateSCP attempts to validate SCP policies by ensuring we have the correct permissions
func (c *awsClient) ValidateSCP(target *string, policies map[string]*cmv1.AWSSTSPolicy) (bool, error) {
	policyDetails := GetPolicyDetails(policies, "osd_scp_policy")

	sParams := &SimulateParams{
		Region: c.GetRegion(),
	}
	// Read installer permissions and OSD SCP Policy permissions
	osdPolicyDocument, err := ParsePolicyDocument(policyDetails)
	if err != nil {
		return false, err
	}

	// Get Creator details
	creator, err := c.GetCreator()
	if err != nil {
		return false, err
	}

	// Find target user
	var targetUserARN arn.ARN
	if target == nil {
		var err error
		callerIdentity, _, err := getClientDetails(c)
		if err != nil {
			return false, fmt.Errorf("getClientDetails: %v\n"+
				"Run 'rosa init' and try again", err)
		}
		targetUserARN, err = arn.Parse(*callerIdentity.Arn)
		if err != nil {
			return false, fmt.Errorf("unable to parse caller ARN %v", err)
		}
		// If the client is using STS credentials want to validate the role
		// the user has assumed. GetCreator() resolves that for us and updates
		// the ARN
		if creator.IsSTS {
			targetUserARN, err = arn.Parse(creator.ARN)
			if err != nil {
				return false, err
			}
		}
	} else {
		targetIAMOutput, err := c.iamClient.GetUser(context.Background(), &iam.GetUserInput{UserName: target})
		if err != nil {
			return false, fmt.Errorf("iamClient.GetUser: %v\n"+
				"To reset the '%s' account, run 'rosa init --delete-stack' and try again", *target, err)
		}
		targetUserARN, err = arn.Parse(*targetIAMOutput.User.Arn)
		if err != nil {
			return false, fmt.Errorf("unable to parse caller ARN %v", err)
		}
	}

	// Validate permissions
	hasPermissions, err := osdPolicyDocument.checkPermissionsUsingQueryClient(c, targetUserARN.String(), sParams)
	if err != nil {
		return false, err
	}
	if !hasPermissions {
		return false, err
	}

	return true, nil
}
