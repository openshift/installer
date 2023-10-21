package opensearchserverless

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/opensearchserverless"
	"github.com/aws/aws-sdk-go-v2/service/opensearchserverless/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func FindAccessPolicyByNameAndType(ctx context.Context, conn *opensearchserverless.Client, id, policyType string) (*types.AccessPolicyDetail, error) {
	in := &opensearchserverless.GetAccessPolicyInput{
		Name: aws.String(id),
		Type: types.AccessPolicyType(policyType),
	}
	out, err := conn.GetAccessPolicy(ctx, in)
	if err != nil {
		var nfe *types.ResourceNotFoundException
		if errors.As(err, &nfe) {
			return nil, &retry.NotFoundError{
				LastError:   err,
				LastRequest: in,
			}
		}

		return nil, err
	}

	if out == nil || out.AccessPolicyDetail == nil {
		return nil, tfresource.NewEmptyResultError(in)
	}

	return out.AccessPolicyDetail, nil
}

func FindCollectionByID(ctx context.Context, conn *opensearchserverless.Client, id string) (*types.CollectionDetail, error) {
	in := &opensearchserverless.BatchGetCollectionInput{
		Ids: []string{id},
	}
	out, err := conn.BatchGetCollection(ctx, in)
	if err != nil {
		var nfe *types.ResourceNotFoundException
		if errors.As(err, &nfe) {
			return nil, &retry.NotFoundError{
				LastError:   err,
				LastRequest: in,
			}
		}

		return nil, err
	}

	if out == nil || out.CollectionDetails == nil || len(out.CollectionDetails) == 0 {
		return nil, tfresource.NewEmptyResultError(in)
	}

	return &out.CollectionDetails[0], nil
}

func FindSecurityConfigByID(ctx context.Context, conn *opensearchserverless.Client, id string) (*types.SecurityConfigDetail, error) {
	in := &opensearchserverless.GetSecurityConfigInput{
		Id: aws.String(id),
	}
	out, err := conn.GetSecurityConfig(ctx, in)
	if err != nil {
		var nfe *types.ResourceNotFoundException
		if errors.As(err, &nfe) {
			return nil, &retry.NotFoundError{
				LastError:   err,
				LastRequest: in,
			}
		}

		return nil, err
	}

	if out == nil || out.SecurityConfigDetail == nil {
		return nil, tfresource.NewEmptyResultError(in)
	}

	return out.SecurityConfigDetail, nil
}

func FindSecurityPolicyByNameAndType(ctx context.Context, conn *opensearchserverless.Client, name, policyType string) (*types.SecurityPolicyDetail, error) {
	in := &opensearchserverless.GetSecurityPolicyInput{
		Name: aws.String(name),
		Type: types.SecurityPolicyType(policyType),
	}
	out, err := conn.GetSecurityPolicy(ctx, in)
	if err != nil {
		var nfe *types.ResourceNotFoundException
		if errors.As(err, &nfe) {
			return nil, &retry.NotFoundError{
				LastError:   err,
				LastRequest: in,
			}
		}

		return nil, err
	}

	if out == nil || out.SecurityPolicyDetail == nil {
		return nil, tfresource.NewEmptyResultError(in)
	}

	return out.SecurityPolicyDetail, nil
}

func FindVPCEndpointByID(ctx context.Context, conn *opensearchserverless.Client, id string) (*types.VpcEndpointDetail, error) {
	in := &opensearchserverless.BatchGetVpcEndpointInput{
		Ids: []string{id},
	}
	out, err := conn.BatchGetVpcEndpoint(ctx, in)

	if err != nil {
		var nfe *types.ResourceNotFoundException
		if errors.As(err, &nfe) {
			return nil, &retry.NotFoundError{
				LastError:   err,
				LastRequest: in,
			}
		}

		return nil, err
	}

	if out == nil || out.VpcEndpointDetails == nil || len(out.VpcEndpointDetails) == 0 {
		return nil, tfresource.NewEmptyResultError(in)
	}

	return &out.VpcEndpointDetails[0], nil
}
