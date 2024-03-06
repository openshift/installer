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

package aws

import (
	"fmt"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"

	"github.com/openshift/rosa/pkg/aws/tags"
)

const (
	OIDCClientIDOpenShift = "openshift"
	OIDCClientIDSTSAWS    = "sts.amazonaws.com"
)

func (c *awsClient) CreateOpenIDConnectProvider(providerURL string, thumbprint string, clusterID string) (
	string, error) {
	iamTags := []*iam.Tag{
		{
			Key:   aws.String(tags.RedHatManaged),
			Value: aws.String(tags.True),
		},
	}
	if clusterID != "" {
		iamTags = append(iamTags, &iam.Tag{
			Key:   aws.String(tags.ClusterID),
			Value: aws.String(clusterID),
		})
	}
	output, err := c.iamClient.CreateOpenIDConnectProvider(&iam.CreateOpenIDConnectProviderInput{
		ClientIDList: []*string{
			aws.String(OIDCClientIDOpenShift),
			aws.String(OIDCClientIDSTSAWS),
		},
		ThumbprintList: []*string{aws.String(thumbprint)},
		Url:            aws.String(providerURL),
		Tags:           iamTags,
	})
	if err != nil {
		return "", err
	}

	return aws.StringValue(output.OpenIDConnectProviderArn), nil
}

func (c *awsClient) HasOpenIDConnectProvider(issuerURL string, accountID string) (bool, error) {
	parsedIssuerURL, err := url.ParseRequestURI(issuerURL)
	if err != nil {
		return false, err
	}
	providerURL := fmt.Sprintf("%s%s", parsedIssuerURL.Host, parsedIssuerURL.Path)

	oidcProviderARN := GetOIDCProviderARN(accountID, providerURL)
	output, err := c.iamClient.GetOpenIDConnectProvider(&iam.GetOpenIDConnectProviderInput{
		OpenIDConnectProviderArn: aws.String(oidcProviderARN),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				return false, nil
			default:
				return false, err
			}
		}
	}
	if aws.StringValue(output.Url) != providerURL {
		return false, fmt.Errorf("The OIDC provider exists but is misconfigured")
	}
	return true, nil
}

func (c *awsClient) DeleteOpenIDConnectProvider(oidcProviderARN string) error {
	_, err := c.iamClient.DeleteOpenIDConnectProvider(&iam.DeleteOpenIDConnectProviderInput{
		OpenIDConnectProviderArn: aws.String(oidcProviderARN),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				return fmt.Errorf("OIDC provider '%s' does not exists",
					oidcProviderARN)
			}
		}
		return err
	}
	return nil
}
