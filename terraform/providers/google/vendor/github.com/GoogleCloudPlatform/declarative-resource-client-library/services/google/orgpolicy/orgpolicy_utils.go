// Copyright 2024 Google LLC. All Rights Reserved.
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
//     http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// Package orgpolicy defines types and methods for working with orgpolicy GCP resources.
package orgpolicy

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func expandPolicyName(_ *Client, name *string, res *Policy) (*string, error) {
	nameParts := strings.Split(dcl.ValueOrEmptyString(name), "/")
	if len(nameParts) == 4 {
		fullName := strings.Join(nameParts, "/")
		return &fullName, nil
	}
	shortName := nameParts[len(nameParts)-1]
	fullName := fmt.Sprintf("%s/policies/%s", dcl.ValueOrEmptyString(res.Parent), shortName)
	return &fullName, nil
}

func equalsPolicyName(m, n *string) bool {
	if m == nil && n == nil {
		return true
	}
	if m == nil || n == nil {
		return false
	}
	return *dcl.SelfLinkToName(m) == *dcl.SelfLinkToName(n)
}

// Compares two values of policy name. Custom diff function required because API returns project numbers.
func canonicalizePolicyName(m, n interface{}) bool {
	mString, ok := m.(*string)
	if !ok {
		return false
	}
	nString, ok := n.(*string)
	if !ok {
		return false
	}
	return equalsPolicyName(mString, nString)
}

func equalsPolicyRulesConditionExpression(m, n *string) bool {
	if m == nil && n == nil {
		return true
	}
	if m == nil || n == nil {
		return false
	}
	mReplaced := strings.ReplaceAll(strings.ReplaceAll(*m, "Labels", "TagId"), "label", "tag")
	nReplaced := strings.ReplaceAll(strings.ReplaceAll(*n, "Labels", "TagId"), "label", "tag")
	return mReplaced == nReplaced
}

// Compares two values of policy rules condition expression. Custom diff function required due to API substitutions.
func canonicalizePolicyRulesConditionExpression(m, n interface{}) bool {
	mString, ok := m.(*string)
	if !ok {
		return false
	}
	nString, ok := n.(*string)
	if !ok {
		return false
	}
	return equalsPolicyRulesConditionExpression(mString, nString)
}

// Update has a custom method because the update mask needs to be in the request body.
func (op *updatePolicyUpdatePolicyOperation) do(ctx context.Context, r *Policy, c *Client) error {
	_, err := c.GetPolicy(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "UpdatePolicy")
	if err != nil {
		return err
	}

	updateMask := "*"
	if (r.Spec != nil && !r.Spec.Empty()) && (r.DryRunSpec == nil || r.DryRunSpec.Empty()) {
		updateMask = "policy.spec"
	}
	if (r.DryRunSpec != nil && !r.DryRunSpec.Empty()) && (r.Spec == nil || r.Spec.Empty()) {
		updateMask = "policy.dryRunSpec"
	}
	u, err = dcl.AddQueryParams(u, map[string]string{"updateMask": updateMask})
	if err != nil {
		return err
	}

	req, err := newUpdatePolicyUpdatePolicyRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdatePolicyUpdatePolicyRequest(c, req)
	if err != nil {
		return err
	}
	_, err = dcl.SendRequest(ctx, c.Config, "PATCH", u, bytes.NewBuffer(body), c.Config.RetryProvider)
	if err != nil {
		return err
	}

	return nil
}
