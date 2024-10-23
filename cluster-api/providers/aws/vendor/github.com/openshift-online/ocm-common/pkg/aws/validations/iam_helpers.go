package validations

import (
	"fmt"
	"maps"

	"github.com/aws/aws-sdk-go-v2/aws"

	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	semver "github.com/hashicorp/go-version"
	. "github.com/openshift-online/ocm-common/pkg/aws/consts"
	. "github.com/openshift-online/ocm-common/pkg/rosa/accountroles"
	. "github.com/openshift-online/ocm-common/pkg/rosa/operatorroles"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

const (
	duplicateIamRoleArnErrorMsg = "ROSA IAM roles must have unique ARNs " +
		"and should not be shared with other IAM roles within the same cluster. " +
		"Duplicated role arn: %s"
)

func GetRoleName(prefix string, role string) string {
	name := fmt.Sprintf("%s-%s-Role", prefix, role)
	if len(name) > MaxAwsRoleLength {
		name = name[0:MaxAwsRoleLength]
	}
	return name
}

func IsManagedRole(roleTags []iamtypes.Tag) bool {
	for _, tag := range roleTags {
		if aws.ToString(tag.Key) == ManagedPolicies && aws.ToString(tag.Value) == "true" {
			return true
		}
	}

	return false
}

func HasCompatibleVersionTags(iamTags []iamtypes.Tag, version string) (bool, error) {
	if len(iamTags) == 0 {
		return false, nil
	}

	wantedVersion, err := semver.NewVersion(version)
	if err != nil {
		return false, err
	}

	for _, tag := range iamTags {
		if aws.ToString(tag.Key) == OpenShiftVersion {
			if version == aws.ToString(tag.Value) {
				return true, nil
			}

			currentVersion, err := semver.NewVersion(aws.ToString(tag.Value))
			if err != nil {
				return false, err
			}
			return currentVersion.GreaterThanOrEqual(wantedVersion), nil
		}
	}
	return false, nil
}

func IamResourceHasTag(iamTags []iamtypes.Tag, tagKey string, tagValue string) bool {
	for _, tag := range iamTags {
		if aws.ToString(tag.Key) == tagKey && aws.ToString(tag.Value) == tagValue {
			return true
		}
	}

	return false
}

func IamRoleArnsValidator(cluster *cmv1.Cluster) error {
	validatingMap := map[string]struct{}{}
	clusterIamRoles := GetAccountRolesArnsMap(cluster)
	maps.Copy(clusterIamRoles, GetOperatorRolesArnsMap(cluster))

	for _, arn := range clusterIamRoles {
		if _, exist := validatingMap[arn]; exist {
			return fmt.Errorf(duplicateIamRoleArnErrorMsg, arn)
		}
		validatingMap[arn] = struct{}{}
	}
	return nil
}
