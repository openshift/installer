package codebuild

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codebuild"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

// FindReportGroupByARN returns the Report Group corresponding to the specified Arn.
func FindReportGroupByARN(ctx context.Context, conn *codebuild.CodeBuild, arn string) (*codebuild.ReportGroup, error) {
	output, err := conn.BatchGetReportGroupsWithContext(ctx, &codebuild.BatchGetReportGroupsInput{
		ReportGroupArns: aws.StringSlice([]string{arn}),
	})
	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, nil
	}

	if len(output.ReportGroups) == 0 {
		return nil, nil
	}

	reportGroup := output.ReportGroups[0]
	if reportGroup == nil {
		return nil, nil
	}

	return reportGroup, nil
}

func FindProjectByARN(ctx context.Context, conn *codebuild.CodeBuild, arn string) (*codebuild.Project, error) {
	input := &codebuild.BatchGetProjectsInput{
		Names: []*string{aws.String(arn)},
	}

	output, err := conn.BatchGetProjectsWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	if output == nil || len(output.Projects) == 0 || output.Projects[0] == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	if count := len(output.Projects); count > 1 {
		return nil, tfresource.NewTooManyResultsError(count, input)
	}

	return output.Projects[0], nil
}

func FindResourcePolicyByARN(ctx context.Context, conn *codebuild.CodeBuild, arn string) (*codebuild.GetResourcePolicyOutput, error) {
	input := &codebuild.GetResourcePolicyInput{
		ResourceArn: aws.String(arn),
	}

	output, err := conn.GetResourcePolicyWithContext(ctx, input)
	if tfawserr.ErrMessageContains(err, codebuild.ErrCodeResourceNotFoundException, "Resource ARN does not exist") ||
		tfawserr.ErrMessageContains(err, codebuild.ErrCodeResourceNotFoundException, "Resource ARN resource policy does not exist") {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	return output, nil
}

func FindSourceCredentialByARN(ctx context.Context, conn *codebuild.CodeBuild, arn string) (*codebuild.SourceCredentialsInfo, error) {
	var result *codebuild.SourceCredentialsInfo
	input := &codebuild.ListSourceCredentialsInput{}
	output, err := conn.ListSourceCredentialsWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	for _, sourceCred := range output.SourceCredentialsInfos {
		if sourceCred == nil {
			continue
		}

		if aws.StringValue(sourceCred.Arn) == arn {
			result = sourceCred
			break
		}
	}

	if result == nil {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	return result, nil
}
