package dms

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	dms "github.com/aws/aws-sdk-go/service/databasemigrationservice"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func statusEndpoint(ctx context.Context, conn *dms.DatabaseMigrationService, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindEndpointByID(ctx, conn, id)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.Status), nil
	}
}

func statusReplicationTask(ctx context.Context, conn *dms.DatabaseMigrationService, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindReplicationTaskByID(ctx, conn, id)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.Status), nil
	}
}
