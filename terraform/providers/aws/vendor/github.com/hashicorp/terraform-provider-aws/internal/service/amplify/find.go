package amplify

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/amplify"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func FindAppByID(ctx context.Context, conn *amplify.Amplify, id string) (*amplify.App, error) {
	input := &amplify.GetAppInput{
		AppId: aws.String(id),
	}

	output, err := conn.GetAppWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, amplify.ErrCodeNotFoundException) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || output.App == nil {
		return nil, &retry.NotFoundError{
			Message:     "Empty result",
			LastRequest: input,
		}
	}

	return output.App, nil
}

func FindBackendEnvironmentByAppIDAndEnvironmentName(ctx context.Context, conn *amplify.Amplify, appID, environmentName string) (*amplify.BackendEnvironment, error) {
	input := &amplify.GetBackendEnvironmentInput{
		AppId:           aws.String(appID),
		EnvironmentName: aws.String(environmentName),
	}

	output, err := conn.GetBackendEnvironmentWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, amplify.ErrCodeNotFoundException) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || output.BackendEnvironment == nil {
		return nil, &retry.NotFoundError{
			Message:     "Empty result",
			LastRequest: input,
		}
	}

	return output.BackendEnvironment, nil
}

func FindBranchByAppIDAndBranchName(ctx context.Context, conn *amplify.Amplify, appID, branchName string) (*amplify.Branch, error) {
	input := &amplify.GetBranchInput{
		AppId:      aws.String(appID),
		BranchName: aws.String(branchName),
	}

	output, err := conn.GetBranchWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, amplify.ErrCodeNotFoundException) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || output.Branch == nil {
		return nil, &retry.NotFoundError{
			Message:     "Empty result",
			LastRequest: input,
		}
	}

	return output.Branch, nil
}

func FindDomainAssociationByAppIDAndDomainName(ctx context.Context, conn *amplify.Amplify, appID, domainName string) (*amplify.DomainAssociation, error) {
	input := &amplify.GetDomainAssociationInput{
		AppId:      aws.String(appID),
		DomainName: aws.String(domainName),
	}

	output, err := conn.GetDomainAssociationWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, amplify.ErrCodeNotFoundException) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || output.DomainAssociation == nil {
		return nil, &retry.NotFoundError{
			Message:     "Empty result",
			LastRequest: input,
		}
	}

	return output.DomainAssociation, nil
}

func FindWebhookByID(ctx context.Context, conn *amplify.Amplify, id string) (*amplify.Webhook, error) {
	input := &amplify.GetWebhookInput{
		WebhookId: aws.String(id),
	}

	output, err := conn.GetWebhookWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, amplify.ErrCodeNotFoundException) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || output.Webhook == nil {
		return nil, &retry.NotFoundError{
			Message:     "Empty result",
			LastRequest: input,
		}
	}

	return output.Webhook, nil
}
