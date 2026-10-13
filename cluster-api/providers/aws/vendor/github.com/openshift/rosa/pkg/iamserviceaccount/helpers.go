/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package iamserviceaccount

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	// ServiceAccountRoleType is the role type tag for service account roles
	ServiceAccountRoleType = "ServiceAccountRole"

	// RoleTypeTagKey is the tag key for the role type
	RoleTypeTagKey = "rosa_role_type"

	// ServiceAccountTagKey is the tag key for the service account name
	ServiceAccountTagKey = "rosa.openshift.io/service-account"

	// NamespaceTagKey is the tag key for the namespace
	NamespaceTagKey = "rosa.openshift.io/namespace"

	// ClusterTagKey is the tag key for the cluster name
	ClusterTagKey = "rosa.openshift.io/cluster"
)

var (
	// ServiceAccountNameRE validates Kubernetes service account names
	ServiceAccountNameRE = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`)

	// NamespaceNameRE validates Kubernetes namespace names
	NamespaceNameRE = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)
)

// GenerateRoleName creates a standardized role name for a service account
func GenerateRoleName(clusterName, namespace, serviceAccountName string) string {
	// Create a role name following the pattern: {cluster-name}-{namespace}-{service-account}-role
	// Ensure it meets AWS IAM role name requirements (64 chars max, alphanumeric + _+=,.@-)
	roleName := fmt.Sprintf("%s-%s-%s-role", clusterName, namespace, serviceAccountName)

	// Truncate if too long
	if len(roleName) > 64 {
		// Keep the suffix and truncate the beginning
		suffix := fmt.Sprintf("-%s-%s-role", namespace, serviceAccountName)
		maxClusterLen := 64 - len(suffix)
		if maxClusterLen > 0 {
			roleName = clusterName[:maxClusterLen] + suffix
		} else {
			// If still too long, create a hash-based name
			roleName = fmt.Sprintf("rosa-%s-%s-role", namespace, serviceAccountName)
			if len(roleName) > 64 {
				roleName = roleName[:64]
			}
		}
	}

	return roleName
}

// GenerateTrustPolicy creates an OIDC trust policy for the service account
func GenerateTrustPolicy(oidcProviderARN, namespace, serviceAccountName string) string {
	// Extract the OIDC provider URL from the ARN
	// ARN format: arn:aws:iam::123456789012:oidc-provider/rh-oidc.s3.us-east-1.amazonaws.com/1234567890abcdef
	parts := strings.Split(oidcProviderARN, "/")
	if len(parts) < 2 {
		return ""
	}
	oidcProviderURL := strings.Join(parts[1:], "/")

	trustPolicy := fmt.Sprintf(`{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Federated": "%s"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
          "%s:sub": "system:serviceaccount:%s:%s"
        }
      }
    }
  ]
}`, oidcProviderARN, oidcProviderURL, namespace, serviceAccountName)

	return trustPolicy
}

// GenerateTrustPolicyMultiple creates an OIDC trust policy for multiple service accounts
func GenerateTrustPolicyMultiple(oidcProviderARN string, serviceAccounts []ServiceAccountIdentifier) string {
	// Extract the OIDC provider URL from the ARN
	parts := strings.Split(oidcProviderARN, "/")
	if len(parts) < 2 {
		return ""
	}
	oidcProviderURL := strings.Join(parts[1:], "/")

	// Build the list of service account subjects
	subjects := make([]string, 0, len(serviceAccounts))
	for _, sa := range serviceAccounts {
		subject := fmt.Sprintf("system:serviceaccount:%s:%s", sa.Namespace, sa.Name)
		subjects = append(subjects, subject)
	}

	// Generate the trust policy based on the number of subjects
	var trustPolicy string
	if len(subjects) == 1 {
		// Single subject - use string format for backwards compatibility
		trustPolicy = fmt.Sprintf(`{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Federated": "%s"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
          "%s:sub": "%s"
        }
      }
    }
  ]
}`, oidcProviderARN, oidcProviderURL, subjects[0])
	} else {
		// Multiple subjects - use array format
		subjectsJSON := `[`
		for i, subject := range subjects {
			if i > 0 {
				subjectsJSON += ", "
			}
			subjectsJSON += fmt.Sprintf(`"%s"`, subject)
		}
		subjectsJSON += `]`

		trustPolicy = fmt.Sprintf(`{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Federated": "%s"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
          "%s:sub": %s
        }
      }
    }
  ]
}`, oidcProviderARN, oidcProviderURL, subjectsJSON)
	}

	return trustPolicy
}

// ValidateServiceAccountName validates a Kubernetes service account name
func ValidateServiceAccountName(name string) error {
	if name == "" {
		return fmt.Errorf("service account name cannot be empty")
	}

	if len(name) > 253 {
		return fmt.Errorf("service account name cannot be longer than 253 characters")
	}

	if !ServiceAccountNameRE.MatchString(name) {
		return fmt.Errorf("service account name must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character")
	}

	return nil
}

// ValidateNamespaceName validates a Kubernetes namespace name
func ValidateNamespaceName(name string) error {
	if name == "" {
		return fmt.Errorf("namespace name cannot be empty")
	}

	if len(name) > 63 {
		return fmt.Errorf("namespace name cannot be longer than 63 characters")
	}

	if !NamespaceNameRE.MatchString(name) {
		return fmt.Errorf("namespace name must consist of lower case alphanumeric characters or '-', and must start and end with an alphanumeric character")
	}

	// Reserved namespaces (only system namespaces)
	reserved := []string{"kube-system", "kube-public", "kube-node-lease"}
	for _, r := range reserved {
		if name == r {
			return fmt.Errorf("namespace '%s' is reserved and cannot be used", name)
		}
	}

	return nil
}

// GenerateDefaultTags creates default tags for service account roles
func GenerateDefaultTags(clusterName, namespace, serviceAccountName string) map[string]string {
	return map[string]string{
		RoleTypeTagKey:       ServiceAccountRoleType,
		ClusterTagKey:        clusterName,
		NamespaceTagKey:      namespace,
		ServiceAccountTagKey: serviceAccountName,
		"red-hat-managed":    "true",
	}
}

// GetRoleARN constructs the ARN for a role given the account ID, role name, path, and partition
func GetRoleARN(accountID, roleName, path, partition string) string {
	if path == "" {
		path = "/"
	}
	return fmt.Sprintf("arn:%s:iam::%s:role%s%s", partition, accountID, path, roleName)
}

// ServiceAccountNameValidator is an interactive validator for service account names
func ServiceAccountNameValidator(val interface{}) error {
	return ValidateServiceAccountName(val.(string))
}

// NamespaceNameValidator is an interactive validator for namespace names
func NamespaceNameValidator(val interface{}) error {
	return ValidateNamespaceName(val.(string))
}
