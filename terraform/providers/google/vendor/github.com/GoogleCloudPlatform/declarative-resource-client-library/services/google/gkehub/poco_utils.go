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
package gkehub

func alsoExpandEmptyBundlesInMap(c *Client, f map[string]FeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContentBundles, res *FeatureMembership) (map[string]any, error) {
	if len(f) == 0 {
		return nil, nil
	}

	items := make(map[string]any)
	for k, v := range f {
		i, err := alsoExpandEmptyBundles(c, &v, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}
	return items, nil
}

func alsoExpandEmptyBundles(c *Client, f *FeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContentBundles, res *FeatureMembership) (map[string]any, error) {
	m := make(map[string]any)
	if v := f.ExemptedNamespaces; v != nil {
		m["exemptedNamespaces"] = v
	}
	return m, nil
}
