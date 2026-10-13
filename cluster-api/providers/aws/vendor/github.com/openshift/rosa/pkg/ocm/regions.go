/*
Copyright (c) 2020 Red Hat, Inc.

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
	"errors"
	"fmt"

	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	"github.com/zgalor/weberr"

	"github.com/openshift/rosa/pkg/arguments"
	"github.com/openshift/rosa/pkg/aws"
	"github.com/openshift/rosa/pkg/helper"
	"github.com/openshift/rosa/pkg/logging"
)

// GetFilteredRegionsByVersion fetches a list of regions. The 'version' argument is optional for filtering.
func (c *Client) GetFilteredRegionsByVersion(roleARN string, version string,
	awsClient aws.Client, externalID string) (regions []*cmv1.CloudRegion, err error) {
	cloudProviderDataBuilder, err := c.createCloudProviderDataBuilder(roleARN, awsClient, externalID)
	if err != nil {
		return []*cmv1.CloudRegion{}, err
	}
	if version != "" {
		cloudProviderDataBuilder = cloudProviderDataBuilder.Version(cmv1.NewVersion().ID(version))
	}

	cloudProviderData, err := cloudProviderDataBuilder.Build()
	if err != nil {
		return []*cmv1.CloudRegion{}, err
	}

	return c.getFilteredRegions(cloudProviderData)
}

func (c *Client) getFilteredRegions(cloudProviderData *cmv1.CloudProviderData) ([]*cmv1.CloudRegion, error) {
	collection := c.ocm.ClustersMgmt().V1().AWSInquiries().Regions()
	page := 1
	size := 100

	var cloudRegions []*cmv1.CloudRegion
	for {
		response, err := collection.Search().
			Body(cloudProviderData).
			Page(page).
			Size(size).
			Send()
		if err != nil {
			return []*cmv1.CloudRegion{}, err
		}

		cloudRegions = append(cloudRegions, response.Items().Slice()...)
		if response.Size() < size {
			break
		}
		page++
	}

	return cloudRegions, nil
}

func (c *Client) GetRegions(roleARN string, externalID string) (regions []*cmv1.CloudRegion, err error) {
	// Retrieve AWS credentials from the local AWS user
	// pass these to OCM to validate what regions are available
	// in this AWS account

	// Build AWS client and retrieve credentials
	// This ensures we use the profile flag if passed to rosa
	// Create the AWS client:
	logger := logging.NewLogger()

	awsBuilder := cmv1.NewAWS()
	if roleARN != "" {
		stsBuilder := cmv1.NewSTS().RoleARN(roleARN)
		if externalID != "" {
			stsBuilder = stsBuilder.ExternalID(externalID)
		}
		awsBuilder = awsBuilder.STS(stsBuilder)
	} else {
		awsClient, err := aws.NewClient().
			Logger(logger).
			Build()
		if err != nil {
			return nil, fmt.Errorf("Error creating AWS client: %v", err)
		}

		// Get AWS region
		currentAWSCreds, err := awsClient.GetIAMCredentials()

		if err != nil {
			return nil, fmt.Errorf("Failed to get local AWS credentials: %v", err)
		}

		awsBuilder = awsBuilder.
			AccessKeyID(currentAWSCreds.AccessKeyID).
			SecretAccessKey(currentAWSCreds.SecretAccessKey)
	}

	awsCredentials, err := awsBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("Failed to build AWS credentials for user '%s': %v", aws.AdminUserName, err)
	}

	collection := c.ocm.ClustersMgmt().V1().
		CloudProviders().
		CloudProvider("aws").
		AvailableRegions()
	page := 1
	size := 100
	for {
		var response *cmv1.AvailableRegionsSearchResponse
		response, err = collection.Search().
			Page(page).
			Size(size).
			Body(awsCredentials).
			Send()
		if err != nil {
			errMsg := response.Error().Reason()
			if errMsg == "" {
				errMsg = err.Error()
			}
			return nil, errors.New(errMsg)
		}
		regions = append(regions, response.Items().Slice()...)
		if response.Size() < size {
			break
		}
		page++
	}
	return
}

func (c *Client) GetRegionList(multiAZ bool, roleARN string,
	externalID string, version string, awsClient aws.Client, isHostedCP bool,
	shardPinningEnabled bool) (regionList []string,
	regionAZ map[string]bool, err error) {
	regions, err := c.GetFilteredRegionsByVersion(roleARN, version, awsClient, externalID)
	if err != nil {
		err = fmt.Errorf("Failed to retrieve AWS regions: %s", err)
		return
	}

	regionAZ = make(map[string]bool, len(regions))

	for _, v := range regions {
		if !v.Enabled() {
			continue
		}

		if isHostedCP && !shardPinningEnabled && !v.SupportsHypershift() {
			continue
		}

		if !multiAZ || v.SupportsMultiAZ() {
			regionList = append(regionList, v.ID())
		}
		regionAZ[v.ID()] = v.SupportsMultiAZ()
	}

	return
}

func (c *Client) GetDatabaseRegionList() ([]string, error) {
	response, err := c.ocm.ClustersMgmt().V1().CloudProviders().CloudProvider("aws").Regions().List().Send()
	if err != nil {
		return []string{}, weberr.Errorf("Failed to get regions listing: %v", err)
	}
	supportedRegions := []string{}
	response.Items().Range(func(index int, item *cmv1.CloudRegion) bool {
		supportedRegions = append(supportedRegions, item.ID())
		return true
	})
	return supportedRegions, nil
}

func (c *Client) ValidateAwsClientRegion() error {
	awsRegionInUserConfig, err := aws.GetRegion(arguments.GetRegion())
	if err != nil {
		return err
	}
	if awsRegionInUserConfig == "" {
		return fmt.Errorf("AWS region not set")
	}

	supportedRegions, err := c.GetDatabaseRegionList()
	if err != nil {
		return err
	}
	if !helper.Contains(supportedRegions, awsRegionInUserConfig) {
		return fmt.Errorf("Unsupported region '%s', available regions: %s",
			awsRegionInUserConfig, helper.SliceToSortedString(supportedRegions))
	}
	return nil
}
