package gamelift

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/gamelift"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func FindBuildByID(conn *gamelift.GameLift, id string) (*gamelift.Build, error) {
	input := &gamelift.DescribeBuildInput{
		BuildId: aws.String(id),
	}

	output, err := conn.DescribeBuild(input)

	if tfawserr.ErrCodeEquals(err, gamelift.ErrCodeNotFoundException) {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || output.Build == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.Build, nil
}

func FindFleetByID(conn *gamelift.GameLift, id string) (*gamelift.FleetAttributes, error) {
	input := &gamelift.DescribeFleetAttributesInput{
		FleetIds: aws.StringSlice([]string{id}),
	}

	output, err := conn.DescribeFleetAttributes(input)

	if tfawserr.ErrCodeEquals(err, gamelift.ErrCodeNotFoundException) {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if len(output.FleetAttributes) == 0 || output.FleetAttributes[0] == nil {
		return nil, tfresource.NewEmptyResultError(output)
	}

	if count := len(output.FleetAttributes); count > 1 {
		return nil, tfresource.NewTooManyResultsError(count, output)
	}

	fleet := output.FleetAttributes[0]

	if aws.StringValue(fleet.FleetId) != id {
		return nil, tfresource.NewEmptyResultError(id)
	}

	return fleet, nil
}

func FindGameServerGroupByName(conn *gamelift.GameLift, name string) (*gamelift.GameServerGroup, error) {
	input := &gamelift.DescribeGameServerGroupInput{
		GameServerGroupName: aws.String(name),
	}

	output, err := conn.DescribeGameServerGroup(input)

	if tfawserr.ErrCodeEquals(err, gamelift.ErrCodeNotFoundException) {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || output.GameServerGroup == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.GameServerGroup, nil
}

func FindScriptByID(conn *gamelift.GameLift, id string) (*gamelift.Script, error) {
	input := &gamelift.DescribeScriptInput{
		ScriptId: aws.String(id),
	}

	output, err := conn.DescribeScript(input)

	if tfawserr.ErrCodeEquals(err, gamelift.ErrCodeNotFoundException) {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || output.Script == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.Script, nil
}
