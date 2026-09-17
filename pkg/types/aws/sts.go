package aws

import "fmt"

const stsMaxRoleNameLen = 64

// OIDCBucketName returns the S3 bucket name used to host the OIDC
// discovery documents for the given cluster infrastructure ID.
func OIDCBucketName(infraID string) string {
	return fmt.Sprintf("%s-oidc", infraID)
}

// STSRoleName computes the deterministic IAM role name for a given
// credentials request, used by both infrastructure provisioning and
// manifest generation to predict role ARNs.
func STSRoleName(infraID, namespace, name string) string {
	roleName := fmt.Sprintf("%s-%s-%s", infraID, namespace, name)
	if len(roleName) > stsMaxRoleNameLen {
		roleName = roleName[:stsMaxRoleNameLen]
	}
	return roleName
}
