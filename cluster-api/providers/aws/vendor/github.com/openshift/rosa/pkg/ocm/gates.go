/*
Copyright (c) 2022 Red Hat, Inc.

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
	"fmt"

	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

func (c *Client) ListStsGates(version string) (stsVersionGates []*cmv1.VersionGate, err error) {
	versionGates, err := c.ListAllOcpGates(version)
	if err != nil {
		return nil, err
	}

	for _, gate := range versionGates {
		if gate.STSOnly() {
			stsVersionGates = append(stsVersionGates, gate)
		}
	}

	return
}

func (c *Client) ListOcpGates(version string) (stsVersionGates []*cmv1.VersionGate, err error) {
	versionGates, err := c.ListAllOcpGates(version)
	if err != nil {
		return nil, err
	}

	for _, gate := range versionGates {
		if !gate.STSOnly() {
			stsVersionGates = append(stsVersionGates, gate)
		}
	}

	return
}

func (c *Client) ListAllOcpGates(version string) (versionGates []*cmv1.VersionGate, err error) {
	versionGatesRequest := c.ocm.ClustersMgmt().V1().VersionGates()

	page := 1
	size := 100

	query := fmt.Sprintf("version_raw_id_prefix = '%s'", version)

	for {
		response, err := versionGatesRequest.List().
			Page(page).
			Size(size).
			Search(query).
			Send()

		if err != nil {
			return nil, handleErr(response.Error(), err)
		}

		versionGates = append(versionGates, response.Items().Slice()...)

		if response.Size() < size {
			break
		}
		page++
	}

	return
}

func (c *Client) AcknowledgeGate(versionGates []*cmv1.VersionGate) (err error) {
	return
}
