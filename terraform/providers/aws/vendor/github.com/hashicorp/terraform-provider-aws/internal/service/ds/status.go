package ds

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/directoryservice"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func statusDirectoryStage(conn *directoryservice.DirectoryService, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindDirectoryByID(conn, id)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.Stage), nil
	}
}

func statusDirectoryShareStatus(conn *directoryservice.DirectoryService, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindDirectoryByID(conn, id)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.ShareStatus), nil
	}
}

func statusDomainController(conn *directoryservice.DirectoryService, directoryID, domainControllerID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindDomainController(conn, directoryID, domainControllerID)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.Status), nil
	}
}

func statusRadius(_ context.Context, conn *directoryservice.DirectoryService, directoryID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindDirectoryByID(conn, directoryID)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.RadiusStatus), nil
	}
}

func statusRegion(ctx context.Context, conn *directoryservice.DirectoryService, directoryID, regionName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindRegion(ctx, conn, directoryID, regionName)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.Status), nil
	}
}

func statusSharedDirectory(ctx context.Context, conn *directoryservice.DirectoryService, ownerDirectoryID, sharedDirectoryID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindSharedDirectory(ctx, conn, ownerDirectoryID, sharedDirectoryID)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.ShareStatus), nil
	}
}
