package awsbase

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/smithy-go"
	multierror "github.com/hashicorp/go-multierror"
)

// getAccountIDAndPartition gets the account ID and associated partition.
func getAccountIDAndPartition(ctx context.Context, iamClient *iam.Client, stsClient *sts.Client, authProviderName string) (string, string, error) {
	var accountID, partition string
	var err, errors error

	if authProviderName == ec2rolecreds.ProviderName {
		accountID, partition, err = getAccountIDAndPartitionFromEC2Metadata(ctx)
	} else {
		accountID, partition, err = getAccountIDAndPartitionFromIAMGetUser(ctx, iamClient)
	}
	if accountID != "" {
		return accountID, partition, nil
	}
	errors = multierror.Append(errors, err)

	accountID, partition, err = getAccountIDAndPartitionFromSTSGetCallerIdentity(ctx, stsClient)
	if accountID != "" {
		return accountID, partition, nil
	}
	errors = multierror.Append(errors, err)

	accountID, partition, err = getAccountIDAndPartitionFromIAMListRoles(ctx, iamClient)
	if accountID != "" {
		return accountID, partition, nil
	}
	errors = multierror.Append(errors, err)

	return accountID, partition, errors
}

// getAccountIDAndPartitionFromEC2Metadata gets the account ID and associated
// partition from EC2 metadata.
func getAccountIDAndPartitionFromEC2Metadata(ctx context.Context) (string, string, error) {
	log.Println("[DEBUG] Trying to get account information via EC2 Metadata")

	cfg := aws.Config{}

	metadataClient := imds.NewFromConfig(cfg)
	info, err := metadataClient.GetIAMInfo(ctx, &imds.GetIAMInfoInput{})
	if err != nil {
		// We can end up here if there's an issue with the instance metadata service
		// or if we're getting credentials from AdRoll's Hologram (in which case IAMInfo will
		// error out).
		err = fmt.Errorf("failed getting account information via EC2 Metadata IAM information: %w", err)
		log.Printf("[DEBUG] %s", err)
		return "", "", err
	}

	return parseAccountIDAndPartitionFromARN(info.InstanceProfileArn)
}

// getAccountIDAndPartitionFromIAMGetUser gets the account ID and associated
// partition from IAM.
func getAccountIDAndPartitionFromIAMGetUser(ctx context.Context, iamClient iam.GetUserAPIClient) (string, string, error) {
	log.Println("[DEBUG] Trying to get account information via iam:GetUser")

	output, err := iamClient.GetUser(ctx, &iam.GetUserInput{})
	if err != nil {
		// AccessDenied and ValidationError can be raised
		// if credentials belong to federated profile, so we ignore these
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			switch apiErr.ErrorCode() {
			case "AccessDenied", "InvalidClientTokenId", "ValidationError":
				log.Printf("[DEBUG] Ignoring iam:GetUser error: %s", err)
				return "", "", nil
			}
		}
		err = fmt.Errorf("failed getting account information via iam:GetUser: %[1]w", err)
		log.Printf("[DEBUG] %s", err)
		return "", "", err
	}

	if output == nil || output.User == nil {
		err = errors.New("empty iam:GetUser response")
		log.Printf("[DEBUG] %s", err)
		return "", "", err
	}

	return parseAccountIDAndPartitionFromARN(aws.ToString(output.User.Arn))
}

// getAccountIDAndPartitionFromIAMListRoles gets the account ID and associated
// partition from listing IAM roles.
func getAccountIDAndPartitionFromIAMListRoles(ctx context.Context, iamClient iam.ListRolesAPIClient) (string, string, error) {
	log.Println("[DEBUG] Trying to get account information via iam:ListRoles")

	output, err := iamClient.ListRoles(ctx, &iam.ListRolesInput{
		MaxItems: aws.Int32(1),
	})
	if err != nil {
		err = fmt.Errorf("failed getting account information via iam:ListRoles: %w", err)
		log.Printf("[DEBUG] %s", err)
		return "", "", err
	}

	if output == nil || len(output.Roles) < 1 {
		err = fmt.Errorf("empty iam:ListRoles response")
		log.Printf("[DEBUG] %s", err)
		return "", "", err
	}

	return parseAccountIDAndPartitionFromARN(aws.ToString(output.Roles[0].Arn))
}

// getAccountIDAndPartitionFromSTSGetCallerIdentity gets the account ID and associated
// partition from STS caller identity.
func getAccountIDAndPartitionFromSTSGetCallerIdentity(ctx context.Context, stsClient *sts.Client) (string, string, error) {
	log.Println("[DEBUG] Trying to get account information via sts:GetCallerIdentity")

	output, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return "", "", fmt.Errorf("error calling sts:GetCallerIdentity: %[1]w", err)
	}

	if output == nil || output.Arn == nil {
		err = errors.New("empty sts:GetCallerIdentity response")
		log.Printf("[DEBUG] %s", err)
		return "", "", err
	}

	return parseAccountIDAndPartitionFromARN(aws.ToString(output.Arn))
}

func parseAccountIDAndPartitionFromARN(inputARN string) (string, string, error) {
	arn, err := arn.Parse(inputARN)
	if err != nil {
		return "", "", fmt.Errorf("error parsing ARN (%s): %s", inputARN, err)
	}
	return arn.AccountID, arn.Partition, nil
}
