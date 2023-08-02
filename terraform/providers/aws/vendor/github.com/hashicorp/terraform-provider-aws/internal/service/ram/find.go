package ram

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ram"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

const (
	FindInvitationTimeout    = 2 * time.Minute
	FindResourceShareTimeout = 1 * time.Minute
)

// FindResourceShareOwnerOtherAccountsByARN returns the resource share owned by other accounts corresponding to the specified ARN.
// Returns nil if no configuration is found.
func FindResourceShareOwnerOtherAccountsByARN(ctx context.Context, conn *ram.RAM, arn string) (*ram.ResourceShare, error) {
	listResourceSharesInput := &ram.GetResourceSharesInput{
		ResourceOwner:     aws.String(ram.ResourceOwnerOtherAccounts),
		ResourceShareArns: aws.StringSlice([]string{arn}),
	}

	return resourceShare(ctx, conn, listResourceSharesInput)
}

// FindResourceShareOwnerSelfByARN returns the resource share owned by own account corresponding to the specified ARN.
// Returns nil if no configuration is found.
func FindResourceShareOwnerSelfByARN(ctx context.Context, conn *ram.RAM, arn string) (*ram.ResourceShare, error) {
	listResourceSharesInput := &ram.GetResourceSharesInput{
		ResourceOwner:     aws.String(ram.ResourceOwnerSelf),
		ResourceShareArns: aws.StringSlice([]string{arn}),
	}

	return resourceShare(ctx, conn, listResourceSharesInput)
}

// FindResourceShareInvitationByResourceShareARNAndStatus returns the resource share invitation corresponding to the specified resource share ARN.
// Returns nil if no configuration is found.
func FindResourceShareInvitationByResourceShareARNAndStatus(ctx context.Context, conn *ram.RAM, resourceShareArn, status string) (*ram.ResourceShareInvitation, error) {
	var invitation *ram.ResourceShareInvitation

	// Retry for Ram resource share invitation eventual consistency
	err := retry.RetryContext(ctx, FindInvitationTimeout, func() *retry.RetryError {
		i, err := resourceShareInvitationByResourceShareARNAndStatus(ctx, conn, resourceShareArn, status)
		invitation = i

		if err != nil {
			return retry.NonRetryableError(err)
		}

		if invitation == nil {
			return retry.RetryableError(&retry.NotFoundError{})
		}

		return nil
	})

	if tfresource.TimedOut(err) {
		invitation, err = resourceShareInvitationByResourceShareARNAndStatus(ctx, conn, resourceShareArn, status)
	}

	if invitation == nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return invitation, nil
}

// FindResourceShareInvitationByARN returns the resource share invitation corresponding to the specified ARN.
// Returns nil if no configuration is found.
func FindResourceShareInvitationByARN(ctx context.Context, conn *ram.RAM, arn string) (*ram.ResourceShareInvitation, error) {
	var invitation *ram.ResourceShareInvitation

	// Retry for Ram resource share invitation eventual consistency
	err := retry.RetryContext(ctx, FindInvitationTimeout, func() *retry.RetryError {
		i, err := resourceShareInvitationByARN(ctx, conn, arn)
		invitation = i

		if err != nil {
			return retry.NonRetryableError(err)
		}

		if invitation == nil {
			retry.RetryableError(&retry.NotFoundError{})
		}

		return nil
	})

	if tfresource.TimedOut(err) {
		invitation, err = resourceShareInvitationByARN(ctx, conn, arn)
	}

	if invitation == nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return invitation, nil
}

func resourceShare(ctx context.Context, conn *ram.RAM, input *ram.GetResourceSharesInput) (*ram.ResourceShare, error) {
	var shares *ram.GetResourceSharesOutput

	// Retry for Ram resource share eventual consistency
	err := retry.RetryContext(ctx, FindResourceShareTimeout, func() *retry.RetryError {
		ss, err := conn.GetResourceSharesWithContext(ctx, input)
		shares = ss

		if tfawserr.ErrCodeEquals(err, ram.ErrCodeUnknownResourceException) {
			return retry.RetryableError(err)
		}

		if err != nil {
			return retry.NonRetryableError(err)
		}

		if len(shares.ResourceShares) == 0 {
			return retry.RetryableError(&retry.NotFoundError{})
		}

		return nil
	})

	if tfresource.TimedOut(err) {
		shares, err = conn.GetResourceSharesWithContext(ctx, input)
	}

	if err != nil {
		return nil, err
	}

	if shares == nil || len(shares.ResourceShares) == 0 {
		return nil, nil
	}

	return shares.ResourceShares[0], nil
}

func resourceShareInvitationByResourceShareARNAndStatus(ctx context.Context, conn *ram.RAM, resourceShareArn, status string) (*ram.ResourceShareInvitation, error) {
	var invitation *ram.ResourceShareInvitation

	input := &ram.GetResourceShareInvitationsInput{
		ResourceShareArns: []*string{aws.String(resourceShareArn)},
	}

	err := conn.GetResourceShareInvitationsPagesWithContext(ctx, input, func(page *ram.GetResourceShareInvitationsOutput, lastPage bool) bool {
		for _, rsi := range page.ResourceShareInvitations {
			if aws.StringValue(rsi.Status) == status {
				invitation = rsi
				return false
			}
		}

		return !lastPage
	})

	if err != nil {
		return nil, err
	}

	return invitation, nil
}

func resourceShareInvitationByARN(ctx context.Context, conn *ram.RAM, arn string) (*ram.ResourceShareInvitation, error) {
	input := &ram.GetResourceShareInvitationsInput{
		ResourceShareInvitationArns: []*string{aws.String(arn)},
	}

	output, err := conn.GetResourceShareInvitationsWithContext(ctx, input)

	if err != nil {
		return nil, err
	}

	if len(output.ResourceShareInvitations) == 0 {
		return nil, nil
	}

	return output.ResourceShareInvitations[0], nil
}

func FindResourceSharePrincipalAssociationByShareARNPrincipal(ctx context.Context, conn *ram.RAM, resourceShareARN, principal string) (*ram.ResourceShareAssociation, error) {
	input := &ram.GetResourceShareAssociationsInput{
		AssociationType:   aws.String(ram.ResourceShareAssociationTypePrincipal),
		Principal:         aws.String(principal),
		ResourceShareArns: aws.StringSlice([]string{resourceShareARN}),
	}

	output, err := conn.GetResourceShareAssociationsWithContext(ctx, input)

	if err != nil {
		return nil, err
	}

	if output == nil || len(output.ResourceShareAssociations) == 0 || output.ResourceShareAssociations[0] == nil {
		return nil, nil
	}

	return output.ResourceShareAssociations[0], nil
}
