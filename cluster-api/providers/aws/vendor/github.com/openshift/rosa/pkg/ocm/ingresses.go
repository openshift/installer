/*
Copyright (c) 2021 Red Hat, Inc.

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

func (c *Client) GetIngress(clusterId string, ingressKey string) (*cmv1.Ingress, error) {
	ingresses, err := c.GetIngresses(clusterId)
	if err != nil {
		return nil, err
	}

	var ingress *cmv1.Ingress
	for _, item := range ingresses {
		if ingressKey == "apps" && item.Default() {
			ingress = item
		}
		if ingressKey == "apps2" && !item.Default() {
			ingress = item
		}
		if item.ID() == ingressKey {
			ingress = item
		}
	}
	if ingress == nil {
		return nil, fmt.Errorf("Failed to get ingress '%s' for cluster '%s'", ingressKey, clusterId)
	}
	return ingress, nil
}

func (c *Client) GetIngresses(clusterID string) ([]*cmv1.Ingress, error) {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		Ingresses().
		List().Page(1).Size(-1).
		Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}

	return response.Items().Slice(), nil
}

func (c *Client) UpdateIngress(clusterID string, ingress *cmv1.Ingress) (*cmv1.Ingress, error) {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		Ingresses().Ingress(ingress.ID()).
		Update().Body(ingress).
		Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}
	return response.Body(), nil
}

func (c *Client) DeleteIngress(clusterID string, ingressID string) error {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		Ingresses().Ingress(ingressID).
		Delete().
		Send()
	if err != nil {
		return handleErr(response.Error(), err)
	}
	return nil
}
