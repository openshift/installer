package servicediscovery

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicediscovery"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

// StatusOperation fetches the Operation and its Status
func StatusOperation(ctx context.Context, conn *servicediscovery.ServiceDiscovery, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindOperationByID(ctx, conn, id)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.Status), nil
	}
}
