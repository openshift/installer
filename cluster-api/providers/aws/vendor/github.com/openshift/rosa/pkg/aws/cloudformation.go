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
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cloudformationtypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/smithy-go"

	"github.com/openshift/rosa/assets"
)

func readCloudFormationTemplate(path string) (string, error) {
	cfTemplate, err := assets.Asset(path)
	if err != nil {
		return "", fmt.Errorf("Unable to read cloudformation template: %s", err)
	}

	return string(cfTemplate), nil
}

// Ensure osdCcsAdmin IAM user is created
func (c *awsClient) EnsureOsdCcsAdminUser(stackName string, adminUserName string, awsRegion string) (bool, error) {
	userExists := true
	regionForInit, err := c.GetClusterRegionTagForUser(adminUserName)
	if err != nil {
		//If user doesn't exists proceed normally
		//If users exists and tag is not present then add the tag
		//if user exists and tag is present then proceed normally
		var noEntityErr *iamtypes.NoSuchEntityException
		if errors.As(err, &noEntityErr) {
			userExists = false
		} else {
			return false, err
		}
	}
	if userExists {
		if regionForInit == "" {
			err = c.TagUserRegion(adminUserName, DefaultRegion)
			if err != nil {
				return false, err
			}
		}
		return false, nil
	}
	// Check already existing cloudformation stack status
	stackReady, stackStatus, err := c.CheckStackReadyOrNotExisting(stackName)
	if err != nil {
		return false, err
	}

	// Read cloudformation template
	cfTemplatePath := "templates/cloudformation/iam_user_osdCcsAdmin.json"
	cfTemplateBody, err := readCloudFormationTemplate(cfTemplatePath)
	if err != nil {
		return false, err
	}

	// If stack CREATE_COMPLETE or UPGRADE_COMPLETE the stack is already create
	// try to update it in case the cloudformation template has changed
	if stackStatus != nil {
		if (*stackStatus == string(cloudformationtypes.StackStatusCreateComplete)) ||
			(*stackStatus == string(cloudformationtypes.StackStatusUpdateComplete)) {
			err = c.UpdateStack(cfTemplateBody, stackName)
			if err != nil {
				return false, err
			}

			return false, nil
		}
	}

	// If the Cloudformation stack isn't ready, make sure the IAM user
	// doesn't exist or the Cloudformation stack create will fail
	if !stackReady {
		err = c.CheckAdminUserNotExisting(adminUserName)
		if err != nil {
			return false, err
		}
	}

	// Create stack
	_, err = c.CreateStack(cfTemplateBody, stackName)
	if err != nil {
		return false, err
	}

	err = c.TagUserRegion(adminUserName, awsRegion)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *awsClient) CreateStack(cfTemplateBody, stackName string) (bool, error) {
	// Create cloudformation stack
	_, err := c.cfClient.CreateStack(context.Background(), buildCreateStackInput(cfTemplateBody, stackName, []cloudformationtypes.Parameter{}, []cloudformationtypes.Tag{}))
	if err != nil {
		return false, err
	}

	err = waitForStackCreateComplete(context.Background(), c.cfClient, stackName)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *awsClient) CreateStackWithParamsTags(ctx context.Context, cfTemplateBody, stackName string, stackParams, stackTags map[string]string) (*string, error) {
	// stack tags
	var cfTags []cloudformationtypes.Tag
	for k, v := range stackTags {
		cfTags = append(cfTags, cloudformationtypes.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}

	// Create a slice for CloudFormation parameters
	var cfParams []cloudformationtypes.Parameter
	for k, v := range stackParams {
		cfParams = append(cfParams, cloudformationtypes.Parameter{
			ParameterKey:   aws.String(k),
			ParameterValue: aws.String(v),
		})
	}

	// Create cloudformation stack
	stackOutput, err := c.cfClient.CreateStack(ctx, buildCreateStackInput(cfTemplateBody, stackName, cfParams, cfTags))
	if err != nil {
		return nil, err
	}

	return stackOutput.StackId, nil
}

func (c *awsClient) GetCFStack(ctx context.Context, stackName string) (*cloudformationtypes.Stack, error) {
	output, err := c.cfClient.DescribeStacks(ctx, &cloudformation.DescribeStacksInput{
		StackName: aws.String(stackName),
	})
	if err != nil {
		return nil, err
	}

	if len(output.Stacks) == 0 {
		return nil, fmt.Errorf("No CF stacks with name %s found", stackName)
	}

	return &output.Stacks[0], nil
}

func (c *awsClient) DescribeCFStackResources(ctx context.Context, stackName string) (*[]cloudformationtypes.StackResource, error) {
	output, err := c.cfClient.DescribeStackResources(ctx, &cloudformation.DescribeStackResourcesInput{
		StackName: aws.String(stackName),
	})

	if err != nil {
		return nil, err
	}

	return &output.StackResources, nil
}

func (c *awsClient) DeleteCFStack(ctx context.Context, stackName string) error {
	_, err := c.cfClient.DeleteStack(ctx, &cloudformation.DeleteStackInput{
		StackName: aws.String(stackName),
	})

	return err
}

func (c *awsClient) UpdateStack(cfTemplateBody, stackName string) error {
	_, err := c.cfClient.UpdateStack(context.TODO(), buildUpdateStackInput(cfTemplateBody, stackName))
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			if apiErr.ErrorCode() == "ValidationError" &&
				strings.Contains(apiErr.ErrorMessage(), "No updates are to be performed") {
				// No updates are to be performed
				return nil
			}
		}
		return err
	}

	// Wait for CloudFormation update to complete
	err = waitForStackUpdateComplete(context.TODO(), c.cfClient, stackName)
	if err != nil {
		return err
	}

	return nil
}

func (c *awsClient) CheckStackReadyOrNotExisting(stackName string) (stackReady bool, status *string, err error) {
	stackList, err := c.cfClient.ListStacks(context.TODO(), &cloudformation.ListStacksInput{})
	if err != nil {
		return false, nil, err
	}

	for _, summary := range stackList.StackSummaries {
		if *summary.StackName == stackName {
			if (summary.StackStatus == cloudformationtypes.StackStatusCreateComplete) ||
				(summary.StackStatus == cloudformationtypes.StackStatusUpdateComplete) {
				stackStatus := string(summary.StackStatus)
				return true, &stackStatus, nil
			}
			if summary.StackStatus != cloudformationtypes.StackStatusDeleteComplete {
				stackStatus := string(summary.StackStatus)
				return false, &stackStatus, fmt.Errorf("Error creating user: CloudFormation stack %s exists "+
					"with status %s. Expected status is %s.\n"+
					"Ensure %s CloudFormation Stack does not exist, then retry with\n"+
					"rosa init --delete-stack; rosa init",
					*summary.StackName, summary.StackStatus, cloudformationtypes.StackStatusCreateComplete, *summary.StackName)
			}
		}
	}
	return false, nil, nil
}

func (c *awsClient) DeleteOsdCcsAdminUser(stackName string) error {
	deleteStackInput := &cloudformation.DeleteStackInput{
		StackName: aws.String(stackName),
	}

	// Delete cloudformation stack
	_, err := c.cfClient.DeleteStack(context.Background(), deleteStackInput)
	if err != nil {
		var tokenExistsErr *cloudformationtypes.TokenAlreadyExistsException
		if errors.As(err, &tokenExistsErr) {
			return nil
		}
		return err
	}

	// Wait until cloudformation stack deletes
	err = waitForStackDeleteComplete(context.Background(), c.cfClient, stackName)
	if err != nil {
		return err
	}

	return nil
}

// Build cloudformation create stack input
func buildCreateStackInput(cfTemplateBody, stackName string, cfParams []cloudformationtypes.Parameter, cfTags []cloudformationtypes.Tag) *cloudformation.CreateStackInput {
	// Special cloudformation capabilities are required to create IAM resources in AWS
	cfCapabilityIAM := cloudformationtypes.CapabilityCapabilityIam
	cfCapabilityNamedIAM := cloudformationtypes.CapabilityCapabilityNamedIam
	cfCapabilityAutoExpand := cloudformationtypes.CapabilityCapabilityAutoExpand
	cfTemplateCapabilities := []cloudformationtypes.Capability{
		cfCapabilityIAM, cfCapabilityNamedIAM, cfCapabilityAutoExpand}

	return &cloudformation.CreateStackInput{
		Capabilities: cfTemplateCapabilities,
		StackName:    aws.String(stackName),
		TemplateBody: aws.String(cfTemplateBody),
		Parameters:   cfParams,
		Tags:         cfTags,
	}
}

// Build cloudformation update stack input
func buildUpdateStackInput(cfTemplateBody, stackName string) *cloudformation.UpdateStackInput {
	// Special cloudformation capabilities are required to update IAM resources in AWS
	cfCapabilityIAM := cloudformationtypes.CapabilityCapabilityIam
	cfCapabilityNamedIAM := cloudformationtypes.CapabilityCapabilityNamedIam
	cfCapabilityAutoExpand := cloudformationtypes.CapabilityCapabilityAutoExpand
	cfTemplateCapabilities := []cloudformationtypes.Capability{
		cfCapabilityIAM, cfCapabilityNamedIAM, cfCapabilityAutoExpand}

	return &cloudformation.UpdateStackInput{
		Capabilities: cfTemplateCapabilities,
		StackName:    aws.String(stackName),
		TemplateBody: aws.String(cfTemplateBody),
	}
}
