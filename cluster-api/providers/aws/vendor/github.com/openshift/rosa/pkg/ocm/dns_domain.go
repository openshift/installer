/*
Copyright (c) 2023 Red Hat, Inc.
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

package ocm

import (
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

func (c *Client) ListDNSDomains(search string) ([]*cmv1.DNSDomain, error) {
	response, err := c.ocm.ClustersMgmt().V1().
		DNSDomains().
		List().
		Parameter("search", search).
		Parameter("order", "organization.id asc").
		Page(1).Size(-1).
		Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}

	return response.Items().Slice(), nil
}

func (c *Client) DeleteDNSDomain(id string) error {
	response, err := c.ocm.ClustersMgmt().V1().
		DNSDomains().DNSDomain(id).
		Delete().
		Send()
	if err != nil {
		return handleErr(response.Error(), err)
	}

	return nil
}

func (c *Client) CreateDNSDomain(dnsDomain *cmv1.DNSDomain) (*cmv1.DNSDomain, error) {
	response, err := c.ocm.ClustersMgmt().V1().
		DNSDomains().
		Add().
		Body(dnsDomain).
		Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}

	return response.Body(), nil
}
