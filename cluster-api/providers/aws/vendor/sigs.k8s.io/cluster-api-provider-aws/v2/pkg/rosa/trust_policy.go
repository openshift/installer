/*
Copyright 2026 The Kubernetes Authors.

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

package rosa

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

// ApplyTrustPolicyExternalID fetches the current trust policy for the named IAM role,
// injects an sts:ExternalId condition into every sts:AssumeRole Allow statement,
// and updates the role if a change was made.
func ApplyTrustPolicyExternalID(ctx context.Context, iamClient *iam.Client, roleName, externalID string) error {
	roleOutput, err := iamClient.GetRole(ctx, &iam.GetRoleInput{
		RoleName: aws.String(roleName),
	})
	if err != nil {
		return fmt.Errorf("failed to get role %q: %w", roleName, err)
	}

	rawPolicy := aws.ToString(roleOutput.Role.AssumeRolePolicyDocument)
	if rawPolicy == "" {
		return fmt.Errorf("role %q has no assume role policy document", roleName)
	}

	decoded, err := url.QueryUnescape(rawPolicy)
	if err != nil {
		return fmt.Errorf("failed to URL-decode trust policy for role %q: %w", roleName, err)
	}

	updated, changed, err := InjectExternalIDIntoTrustPolicy(decoded, externalID)
	if err != nil {
		return fmt.Errorf("failed to inject external ID into trust policy for role %q: %w", roleName, err)
	}
	if !changed {
		return nil
	}

	_, err = iamClient.UpdateAssumeRolePolicy(ctx, &iam.UpdateAssumeRolePolicyInput{
		RoleName:       aws.String(roleName),
		PolicyDocument: aws.String(updated),
	})
	if err != nil {
		return fmt.Errorf("failed to update assume role policy for role %q: %w", roleName, err)
	}
	return nil
}

// InjectExternalIDIntoTrustPolicy parses a trust policy JSON document and adds
// an sts:ExternalId condition to every sts:AssumeRole Allow statement.
// Returns the modified JSON and whether any change was made.
func InjectExternalIDIntoTrustPolicy(policyJSON, externalID string) (string, bool, error) {
	var doc map[string]interface{}
	if err := json.Unmarshal([]byte(policyJSON), &doc); err != nil {
		return "", false, fmt.Errorf("failed to parse trust policy: %w", err)
	}

	statements, ok := doc["Statement"].([]interface{})
	if !ok {
		return policyJSON, false, nil
	}

	changed := false
	for _, raw := range statements {
		stmt, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		if !isAssumeRoleAllowStatement(stmt) {
			continue
		}
		if setExternalIDCondition(stmt, externalID) {
			changed = true
		}
	}

	if !changed {
		return policyJSON, false, nil
	}

	out, err := json.Marshal(doc)
	if err != nil {
		return "", false, fmt.Errorf("failed to marshal updated trust policy: %w", err)
	}
	return string(out), true, nil
}

func isAssumeRoleAllowStatement(stmt map[string]interface{}) bool {
	effect, _ := stmt["Effect"].(string)
	if effect != "Allow" {
		return false
	}
	switch action := stmt["Action"].(type) {
	case string:
		return action == "sts:AssumeRole"
	case []interface{}:
		for _, a := range action {
			if s, ok := a.(string); ok && s == "sts:AssumeRole" {
				return true
			}
		}
	}
	return false
}

// setExternalIDCondition adds or updates sts:ExternalId in the StringEquals condition.
// Returns true if the statement was modified.
func setExternalIDCondition(stmt map[string]interface{}, externalID string) bool {
	condition, _ := stmt["Condition"].(map[string]interface{})
	if condition == nil {
		stmt["Condition"] = map[string]interface{}{
			"StringEquals": map[string]interface{}{
				"sts:ExternalId": externalID,
			},
		}
		return true
	}

	stringEquals, _ := condition["StringEquals"].(map[string]interface{})
	if stringEquals == nil {
		condition["StringEquals"] = map[string]interface{}{
			"sts:ExternalId": externalID,
		}
		return true
	}

	existing, _ := stringEquals["sts:ExternalId"].(string)
	if existing == externalID {
		return false
	}

	stringEquals["sts:ExternalId"] = externalID
	return true
}
