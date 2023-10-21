package appautoscaling

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/applicationautoscaling"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func FindScheduledAction(ctx context.Context, conn *applicationautoscaling.ApplicationAutoScaling, name, serviceNamespace, resourceId string) (*applicationautoscaling.ScheduledAction, error) {
	var result *applicationautoscaling.ScheduledAction

	input := &applicationautoscaling.DescribeScheduledActionsInput{
		ScheduledActionNames: []*string{aws.String(name)},
		ServiceNamespace:     aws.String(serviceNamespace),
		ResourceId:           aws.String(resourceId),
	}
	err := conn.DescribeScheduledActionsPagesWithContext(ctx, input, func(page *applicationautoscaling.DescribeScheduledActionsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, item := range page.ScheduledActions {
			if item == nil {
				continue
			}

			if name == aws.StringValue(item.ScheduledActionName) {
				result = item
				return false
			}
		}

		return !lastPage
	})
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, &retry.NotFoundError{
			LastRequest: input,
		}
	}

	return result, nil
}
