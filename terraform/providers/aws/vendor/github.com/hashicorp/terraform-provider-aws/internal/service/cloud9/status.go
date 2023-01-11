package cloud9

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloud9"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func statusEnvironmentStatus(conn *cloud9.Cloud9, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindEnvironmentByID(conn, id)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.Lifecycle.Status), nil
	}
}
