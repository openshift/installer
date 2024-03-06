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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/iam"

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
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				userExists = false
			default:
				return false, err
			}
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
		if (*stackStatus == cloudformation.StackStatusCreateComplete) ||
			(*stackStatus == cloudformation.StackStatusUpdateComplete) {
			_, err = c.UpdateStack(cfTemplateBody, stackName)
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
	_, err := c.cfClient.CreateStack(buildCreateStackInput(cfTemplateBody, stackName))
	if err != nil {
		return false, err
	}

	// Wait until cloudformation stack creates
	err = c.cfClient.WaitUntilStackCreateComplete(&cloudformation.DescribeStacksInput{
		StackName: aws.String(stackName),
	})
	if err != nil {
		switch typed := err.(type) {
		case awserr.Error:
			// Waiter reached maximum attempts waiting for the resource to be ready
			if typed.Code() == request.WaiterResourceNotReadyErrorCode {
				c.logger.Errorf("Max retries reached waiting for stack to create")
				return false, err
			}
		}
		return false, err
	}

	return true, nil
}

func (c *awsClient) UpdateStack(cfTemplateBody, stackName string) (bool, error) {
	_, err := c.cfClient.UpdateStack(buildUpdateStackInput(cfTemplateBody, stackName))
	if err != nil {
		switch typed := err.(type) {
		case awserr.Error:
			// Exit true if there is no update to be performed on the cloudformation stack
			if typed.Code() == "ValidationError" {
				if typed.Message() == "No updates are to be performed." {
					return true, nil
				}
			}
		}
		return false, err
	}

	// Wait for CloudFormation update to complete
	err = c.cfClient.WaitUntilStackUpdateComplete(&cloudformation.DescribeStacksInput{
		StackName: aws.String(stackName),
	})

	if err != nil {
		switch typed := err.(type) {
		case awserr.Error:
			// Waiter reached maximum attempts waiting for the resource to be ready
			if typed.Code() == request.WaiterResourceNotReadyErrorCode {
				c.logger.Errorf("Max retries reached waiting for stack to create")
				return false, err
			}
		}
		return false, err
	}

	return true, err
}

func (c *awsClient) CheckStackReadyOrNotExisting(stackName string) (stackReady bool, status *string, err error) {
	stackList, err := c.cfClient.ListStacks(&cloudformation.ListStacksInput{})
	if err != nil {
		return false, nil, err
	}

	for _, summary := range stackList.StackSummaries {
		if *summary.StackName == stackName {
			if (*summary.StackStatus == cloudformation.StackStatusCreateComplete) ||
				(*summary.StackStatus == cloudformation.StackStatusUpdateComplete) {
				return true, summary.StackStatus, nil
			}
			if *summary.StackStatus != cloudformation.StackStatusDeleteComplete {
				return false, summary.StackStatus, fmt.Errorf("Error creating user: Cloudformation stack %s exists "+
					"with status %s. Expected status is %s.\n"+
					"Ensure %s CloudFormation Stack does not exist, then retry with\n"+
					"rosa init --delete-stack; rosa init",
					*summary.StackName, *summary.StackStatus, cloudformation.StackStatusCreateComplete, *summary.StackName)
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
	_, err := c.cfClient.DeleteStack(deleteStackInput)
	if err != nil {
		switch typed := err.(type) {
		case awserr.Error:
			if typed.Code() == cloudformation.ErrCodeTokenAlreadyExistsException {
				return nil
			}
		}
		return err
	}

	// Wait until cloudformation stack deletes
	err = c.cfClient.WaitUntilStackDeleteComplete(&cloudformation.DescribeStacksInput{
		StackName: aws.String(stackName),
	})
	if err != nil {
		switch typed := err.(type) {
		case awserr.Error:
			// Waiter reached maximum attempts waiting for the resource to be ready
			if typed.Code() == request.WaiterResourceNotReadyErrorCode {
				c.logger.Errorf("Max retries reached waiting for stack to delete")
				return err
			}
		}
		return err
	}

	return nil
}

// Build cloudformation create stack input
func buildCreateStackInput(cfTemplateBody, stackName string) *cloudformation.CreateStackInput {
	// Special cloudformation capabilities are required to create IAM resources in AWS
	cfCapabilityIAM := "CAPABILITY_IAM"
	cfCapabilityNamedIAM := "CAPABILITY_NAMED_IAM"
	cfTemplateCapabilities := []*string{&cfCapabilityIAM, &cfCapabilityNamedIAM}

	return &cloudformation.CreateStackInput{
		Capabilities: cfTemplateCapabilities,
		StackName:    aws.String(stackName),
		TemplateBody: aws.String(cfTemplateBody),
	}
}

// Build cloudformation update stack input
func buildUpdateStackInput(cfTemplateBody, stackName string) *cloudformation.UpdateStackInput {
	// Special cloudformation capabilities are required to update IAM resources in AWS
	cfCapabilityIAM := "CAPABILITY_IAM"
	cfCapabilityNamedIAM := "CAPABILITY_NAMED_IAM"
	cfTemplateCapabilities := []*string{&cfCapabilityIAM, &cfCapabilityNamedIAM}

	return &cloudformation.UpdateStackInput{
		Capabilities: cfTemplateCapabilities,
		StackName:    aws.String(stackName),
		TemplateBody: aws.String(cfTemplateBody),
	}
}
