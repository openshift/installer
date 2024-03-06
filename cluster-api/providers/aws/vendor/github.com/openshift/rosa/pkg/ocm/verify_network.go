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

func (c *Client) GetVerifyNetworkSubnet(id string) (*cmv1.SubnetNetworkVerification, error) {
	response, err := c.ocm.ClustersMgmt().V1().NetworkVerifications().
		NetworkVerification(id).Get().Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}

	return response.Body(), nil
}

func (c *Client) VerifyNetworkSubnets(awsAccountId string, region string,
	subnets []string, tags map[string]string, platform cmv1.Platform) ([]*cmv1.SubnetNetworkVerification, error) {

	body, _ := cmv1.NewNetworkVerification().Platform(platform).CloudProviderData(cmv1.NewCloudProviderData().
		AWS(cmv1.NewAWS().STS(cmv1.NewSTS().RoleARN(awsAccountId)).Tags(tags)).
		Subnets(subnets...).
		Region(cmv1.NewCloudRegion().ID(region))).
		Build()
	response, err := c.ocm.ClustersMgmt().V1().NetworkVerifications().Add().
		Body(body).
		Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}

	return response.Body().Items(), nil
}

func (c *Client) VerifyNetworkSubnetsByCluster(clusterId string, tags map[string]string) (
	[]*cmv1.SubnetNetworkVerification, error) {
	body, _ := cmv1.NewNetworkVerification().ClusterId(clusterId).CloudProviderData(
		cmv1.NewCloudProviderData().AWS(cmv1.NewAWS().Tags(tags))).Build()
	response, err := c.ocm.ClustersMgmt().V1().NetworkVerifications().Add().
		Body(body).
		Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}

	return response.Body().Items(), nil
}
