package kms

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func FindAliasByName(conn *kms.KMS, name string) (*kms.AliasListEntry, error) {
	input := &kms.ListAliasesInput{}
	var output *kms.AliasListEntry

	err := conn.ListAliasesPages(input, func(page *kms.ListAliasesOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, alias := range page.Aliases {
			if aws.StringValue(alias.AliasName) == name {
				output = alias

				return false
			}
		}

		return !lastPage
	})

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, &resource.NotFoundError{}
	}

	return output, nil
}

func FindCustomKeyStoreByID(ctx context.Context, conn *kms.KMS, in *kms.DescribeCustomKeyStoresInput) (*kms.CustomKeyStoresListEntry, error) {
	out, err := conn.DescribeCustomKeyStoresWithContext(ctx, in)

	if tfawserr.ErrCodeEquals(err, kms.ErrCodeCustomKeyStoreNotFoundException) {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: in,
		}
	}
	if err != nil {
		return nil, err
	}

	if out == nil || out.CustomKeyStores[0] == nil {
		return nil, tfresource.NewEmptyResultError(in)
	}

	return out.CustomKeyStores[0], nil
}

func FindKeyByID(conn *kms.KMS, id string) (*kms.KeyMetadata, error) {
	input := &kms.DescribeKeyInput{
		KeyId: aws.String(id),
	}

	output, err := conn.DescribeKey(input)

	if tfawserr.ErrCodeEquals(err, kms.ErrCodeNotFoundException) {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || output.KeyMetadata == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	keyMetadata := output.KeyMetadata

	// Once the CMK is in the pending (replica) deletion state Terraform considers it logically deleted.
	if state := aws.StringValue(keyMetadata.KeyState); state == kms.KeyStatePendingDeletion || state == kms.KeyStatePendingReplicaDeletion {
		return nil, &resource.NotFoundError{
			Message:     state,
			LastRequest: input,
		}
	}

	return keyMetadata, nil
}

func FindKeyPolicyByKeyIDAndPolicyName(conn *kms.KMS, keyID, policyName string) (*string, error) {
	input := &kms.GetKeyPolicyInput{
		KeyId:      aws.String(keyID),
		PolicyName: aws.String(policyName),
	}

	output, err := conn.GetKeyPolicy(input)

	if tfawserr.ErrCodeEquals(err, kms.ErrCodeNotFoundException) {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.Policy, nil
}

func FindKeyRotationEnabledByKeyID(conn *kms.KMS, keyID string) (*bool, error) {
	input := &kms.GetKeyRotationStatusInput{
		KeyId: aws.String(keyID),
	}

	output, err := conn.GetKeyRotationStatus(input)

	if tfawserr.ErrCodeEquals(err, kms.ErrCodeNotFoundException) {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.KeyRotationEnabled, nil
}
