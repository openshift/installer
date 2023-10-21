package configservice

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const (
	conformancePackCreateTimeout = 5 * time.Minute
	conformancePackDeleteTimeout = 5 * time.Minute

	conformancePackStatusNotFound = "NotFound"
	conformancePackStatusUnknown  = "Unknown"
)

func DescribeConformancePack(ctx context.Context, conn *configservice.ConfigService, name string) (*configservice.ConformancePackDetail, error) {
	input := &configservice.DescribeConformancePacksInput{
		ConformancePackNames: []*string{aws.String(name)},
	}

	for {
		output, err := conn.DescribeConformancePacksWithContext(ctx, input)

		if err != nil {
			return nil, err
		}

		for _, pack := range output.ConformancePackDetails {
			if pack == nil {
				continue
			}

			if aws.StringValue(pack.ConformancePackName) == name {
				return pack, nil
			}
		}

		if aws.StringValue(output.NextToken) == "" {
			break
		}

		input.NextToken = output.NextToken
	}

	return nil, nil
}

func describeConformancePackStatus(ctx context.Context, conn *configservice.ConfigService, name string) (*configservice.ConformancePackStatusDetail, error) {
	input := &configservice.DescribeConformancePackStatusInput{
		ConformancePackNames: []*string{aws.String(name)},
	}

	for {
		output, err := conn.DescribeConformancePackStatusWithContext(ctx, input)

		if err != nil {
			return nil, err
		}

		for _, status := range output.ConformancePackStatusDetails {
			if aws.StringValue(status.ConformancePackName) == name {
				return status, nil
			}
		}

		if aws.StringValue(output.NextToken) == "" {
			break
		}

		input.NextToken = output.NextToken
	}

	return nil, nil
}

func DescribeOrganizationConfigRule(ctx context.Context, conn *configservice.ConfigService, name string) (*configservice.OrganizationConfigRule, error) {
	input := &configservice.DescribeOrganizationConfigRulesInput{
		OrganizationConfigRuleNames: []*string{aws.String(name)},
	}

	for {
		output, err := conn.DescribeOrganizationConfigRulesWithContext(ctx, input)

		if err != nil {
			return nil, err
		}

		for _, rule := range output.OrganizationConfigRules {
			if aws.StringValue(rule.OrganizationConfigRuleName) == name {
				return rule, nil
			}
		}

		if aws.StringValue(output.NextToken) == "" {
			break
		}

		input.NextToken = output.NextToken
	}

	return nil, nil
}

func describeOrganizationConfigRuleStatus(ctx context.Context, conn *configservice.ConfigService, name string) (*configservice.OrganizationConfigRuleStatus, error) {
	input := &configservice.DescribeOrganizationConfigRuleStatusesInput{
		OrganizationConfigRuleNames: []*string{aws.String(name)},
	}

	for {
		output, err := conn.DescribeOrganizationConfigRuleStatusesWithContext(ctx, input)

		if err != nil {
			return nil, err
		}

		for _, status := range output.OrganizationConfigRuleStatuses {
			if aws.StringValue(status.OrganizationConfigRuleName) == name {
				return status, nil
			}
		}

		if aws.StringValue(output.NextToken) == "" {
			break
		}

		input.NextToken = output.NextToken
	}

	return nil, nil
}

func DescribeOrganizationConformancePack(ctx context.Context, conn *configservice.ConfigService, name string) (*configservice.OrganizationConformancePack, error) {
	input := &configservice.DescribeOrganizationConformancePacksInput{
		OrganizationConformancePackNames: []*string{aws.String(name)},
	}

	for {
		output, err := conn.DescribeOrganizationConformancePacksWithContext(ctx, input)

		if err != nil {
			return nil, err
		}

		for _, pack := range output.OrganizationConformancePacks {
			if aws.StringValue(pack.OrganizationConformancePackName) == name {
				return pack, nil
			}
		}

		if aws.StringValue(output.NextToken) == "" {
			break
		}

		input.NextToken = output.NextToken
	}

	return nil, nil
}

func describeOrganizationConformancePackStatus(ctx context.Context, conn *configservice.ConfigService, name string) (*configservice.OrganizationConformancePackStatus, error) {
	input := &configservice.DescribeOrganizationConformancePackStatusesInput{
		OrganizationConformancePackNames: []*string{aws.String(name)},
	}

	for {
		output, err := conn.DescribeOrganizationConformancePackStatusesWithContext(ctx, input)

		if err != nil {
			return nil, err
		}

		for _, status := range output.OrganizationConformancePackStatuses {
			if aws.StringValue(status.OrganizationConformancePackName) == name {
				return status, nil
			}
		}

		if aws.StringValue(output.NextToken) == "" {
			break
		}

		input.NextToken = output.NextToken
	}

	return nil, nil
}

func getOrganizationConfigRuleDetailedStatus(ctx context.Context, conn *configservice.ConfigService, ruleName, ruleStatus string) ([]*configservice.MemberAccountStatus, error) {
	input := &configservice.GetOrganizationConfigRuleDetailedStatusInput{
		Filters: &configservice.StatusDetailFilters{
			MemberAccountRuleStatus: aws.String(ruleStatus),
		},
		OrganizationConfigRuleName: aws.String(ruleName),
	}
	var statuses []*configservice.MemberAccountStatus

	for {
		output, err := conn.GetOrganizationConfigRuleDetailedStatusWithContext(ctx, input)

		if err != nil {
			return nil, err
		}

		statuses = append(statuses, output.OrganizationConfigRuleDetailedStatus...)

		if aws.StringValue(output.NextToken) == "" {
			break
		}

		input.NextToken = output.NextToken
	}

	return statuses, nil
}

func getOrganizationConformancePackDetailedStatus(ctx context.Context, conn *configservice.ConfigService, name, status string) ([]*configservice.OrganizationConformancePackDetailedStatus, error) {
	input := &configservice.GetOrganizationConformancePackDetailedStatusInput{
		Filters: &configservice.OrganizationResourceDetailedStatusFilters{
			Status: aws.String(status),
		},
		OrganizationConformancePackName: aws.String(name),
	}

	var statuses []*configservice.OrganizationConformancePackDetailedStatus

	for {
		output, err := conn.GetOrganizationConformancePackDetailedStatusWithContext(ctx, input)

		if err != nil {
			return nil, err
		}

		statuses = append(statuses, output.OrganizationConformancePackDetailedStatuses...)

		if aws.StringValue(output.NextToken) == "" {
			break
		}

		input.NextToken = output.NextToken
	}

	return statuses, nil
}

func refreshConformancePackStatus(ctx context.Context, conn *configservice.ConfigService, name string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		status, err := describeConformancePackStatus(ctx, conn, name)

		if err != nil {
			return nil, conformancePackStatusUnknown, err
		}

		if status == nil {
			return nil, conformancePackStatusNotFound, nil
		}

		if errMsg := aws.StringValue(status.ConformancePackStatusReason); errMsg != "" {
			return status, aws.StringValue(status.ConformancePackState), fmt.Errorf(errMsg)
		}

		return status, aws.StringValue(status.ConformancePackState), nil
	}
}

func refreshOrganizationConfigRuleStatus(ctx context.Context, conn *configservice.ConfigService, name string, target string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		status, err := describeOrganizationConfigRuleStatus(ctx, conn, name)

		// Transient ResourceDoesNotExist error after creation caught here
		// in cases where the StateChangeConf's delay time is not sufficient
		if target != configservice.OrganizationResourceDetailedStatusDeleteSuccessful && tfawserr.ErrCodeEquals(err, configservice.ErrCodeNoSuchOrganizationConfigRuleException) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		if status == nil {
			return nil, "", fmt.Errorf("status not found")
		}

		if status.ErrorCode != nil {
			return status, aws.StringValue(status.OrganizationRuleStatus), fmt.Errorf("%s: %s", aws.StringValue(status.ErrorCode), aws.StringValue(status.ErrorMessage))
		}

		switch aws.StringValue(status.OrganizationRuleStatus) {
		case configservice.OrganizationRuleStatusCreateFailed, configservice.OrganizationRuleStatusDeleteFailed, configservice.OrganizationRuleStatusUpdateFailed:
			// Display detailed errors for failed member accounts
			memberAccountStatuses, err := getOrganizationConfigRuleDetailedStatus(ctx, conn, name, aws.StringValue(status.OrganizationRuleStatus))

			if err != nil {
				return status, aws.StringValue(status.OrganizationRuleStatus), fmt.Errorf("unable to get Organization Config Rule detailed status for showing member account errors: %w", err)
			}

			var errBuilder strings.Builder

			for _, mas := range memberAccountStatuses {
				errBuilder.WriteString(fmt.Sprintf("Account ID (%s): %s: %s\n", aws.StringValue(mas.AccountId), aws.StringValue(mas.ErrorCode), aws.StringValue(mas.ErrorMessage)))
			}

			return status, aws.StringValue(status.OrganizationRuleStatus), fmt.Errorf("Failed in %d account(s):\n\n%s", len(memberAccountStatuses), errBuilder.String())
		}

		return status, aws.StringValue(status.OrganizationRuleStatus), nil
	}
}

func refreshOrganizationConformancePackCreationStatus(ctx context.Context, conn *configservice.ConfigService, name string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		status, err := describeOrganizationConformancePackStatus(ctx, conn, name)

		// Transient ResourceDoesNotExist error after creation caught here
		// in cases where the StateChangeConf's delay time is not sufficient
		if tfawserr.ErrCodeEquals(err, configservice.ErrCodeNoSuchOrganizationConformancePackException) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		if status == nil {
			return nil, "", nil
		}

		if status.ErrorCode != nil {
			return status, aws.StringValue(status.Status), fmt.Errorf("%s: %s", aws.StringValue(status.ErrorCode), aws.StringValue(status.ErrorMessage))
		}

		switch s := aws.StringValue(status.Status); s {
		case configservice.OrganizationResourceStatusCreateFailed, configservice.OrganizationResourceStatusDeleteFailed, configservice.OrganizationResourceStatusUpdateFailed:
			return status, s, organizationConformancePackDetailedStatusError(ctx, conn, name, s)
		}

		return status, aws.StringValue(status.Status), nil
	}
}

func refreshOrganizationConformancePackStatus(ctx context.Context, conn *configservice.ConfigService, name string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		status, err := describeOrganizationConformancePackStatus(ctx, conn, name)

		if err != nil {
			return nil, "", err
		}

		if status == nil {
			return nil, "", nil
		}

		if status.ErrorCode != nil {
			return status, aws.StringValue(status.Status), fmt.Errorf("%s: %s", aws.StringValue(status.ErrorCode), aws.StringValue(status.ErrorMessage))
		}

		switch s := aws.StringValue(status.Status); s {
		case configservice.OrganizationResourceStatusCreateFailed, configservice.OrganizationResourceStatusDeleteFailed, configservice.OrganizationResourceStatusUpdateFailed:
			return status, s, organizationConformancePackDetailedStatusError(ctx, conn, name, s)
		}

		return status, aws.StringValue(status.Status), nil
	}
}

func organizationConformancePackDetailedStatusError(ctx context.Context, conn *configservice.ConfigService, name, status string) error {
	memberAccountStatuses, err := getOrganizationConformancePackDetailedStatus(ctx, conn, name, status)

	if err != nil {
		return fmt.Errorf("unable to get Config Organization Conformance Pack detailed status for showing member account errors: %w", err)
	}

	var errBuilder strings.Builder

	for _, mas := range memberAccountStatuses {
		errBuilder.WriteString(fmt.Sprintf("Account ID (%s): %s: %s\n", aws.StringValue(mas.AccountId), aws.StringValue(mas.ErrorCode), aws.StringValue(mas.ErrorMessage)))
	}

	return fmt.Errorf("Failed in %d account(s):\n\n%s", len(memberAccountStatuses), errBuilder.String())
}

func waitForConformancePackStateCreateComplete(ctx context.Context, conn *configservice.ConfigService, name string) error {
	stateChangeConf := retry.StateChangeConf{
		Pending: []string{configservice.ConformancePackStateCreateInProgress},
		Target:  []string{configservice.ConformancePackStateCreateComplete},
		Timeout: conformancePackCreateTimeout,
		Refresh: refreshConformancePackStatus(ctx, conn, name),
	}

	_, err := stateChangeConf.WaitForStateContext(ctx)

	if tfawserr.ErrCodeEquals(err, configservice.ErrCodeNoSuchConformancePackException) {
		return nil
	}

	return err
}

func waitForConformancePackStateDeleteComplete(ctx context.Context, conn *configservice.ConfigService, name string) error {
	stateChangeConf := retry.StateChangeConf{
		Pending: []string{configservice.ConformancePackStateDeleteInProgress},
		Target:  []string{},
		Timeout: conformancePackDeleteTimeout,
		Refresh: refreshConformancePackStatus(ctx, conn, name),
	}

	_, err := stateChangeConf.WaitForStateContext(ctx)

	if tfawserr.ErrCodeEquals(err, configservice.ErrCodeNoSuchConformancePackException) {
		return nil
	}

	return err
}

func waitForOrganizationConformancePackStatusCreateSuccessful(ctx context.Context, conn *configservice.ConfigService, name string, timeout time.Duration) error {
	stateChangeConf := retry.StateChangeConf{
		Pending: []string{configservice.OrganizationResourceStatusCreateInProgress},
		Target:  []string{configservice.OrganizationResourceStatusCreateSuccessful},
		Timeout: timeout,
		Refresh: refreshOrganizationConformancePackCreationStatus(ctx, conn, name),
		// Include a delay to help avoid ResourceDoesNotExist errors
		Delay: 30 * time.Second,
	}

	_, err := stateChangeConf.WaitForStateContext(ctx)

	return err
}

func waitForOrganizationConformancePackStatusUpdateSuccessful(ctx context.Context, conn *configservice.ConfigService, name string, timeout time.Duration) error {
	stateChangeConf := retry.StateChangeConf{
		Pending: []string{configservice.OrganizationResourceStatusUpdateInProgress},
		Target:  []string{configservice.OrganizationResourceStatusUpdateSuccessful},
		Timeout: timeout,
		Refresh: refreshOrganizationConformancePackStatus(ctx, conn, name),
	}

	_, err := stateChangeConf.WaitForStateContext(ctx)

	return err
}

func waitForOrganizationConformancePackStatusDeleteSuccessful(ctx context.Context, conn *configservice.ConfigService, name string, timeout time.Duration) error {
	stateChangeConf := retry.StateChangeConf{
		Pending: []string{configservice.OrganizationResourceStatusDeleteInProgress},
		Target:  []string{configservice.OrganizationResourceStatusDeleteSuccessful},
		Timeout: timeout,
		Refresh: refreshOrganizationConformancePackStatus(ctx, conn, name),
	}

	_, err := stateChangeConf.WaitForStateContext(ctx)

	return err
}

func waitForOrganizationRuleStatusCreateSuccessful(ctx context.Context, conn *configservice.ConfigService, name string, timeout time.Duration) error {
	stateChangeConf := &retry.StateChangeConf{
		Pending:        []string{configservice.OrganizationRuleStatusCreateInProgress},
		Target:         []string{configservice.OrganizationRuleStatusCreateSuccessful},
		Refresh:        refreshOrganizationConfigRuleStatus(ctx, conn, name, configservice.OrganizationRuleStatusCreateSuccessful),
		Timeout:        timeout,
		NotFoundChecks: 10,
		Delay:          30 * time.Second,
	}

	_, err := stateChangeConf.WaitForStateContext(ctx)

	return err
}

func waitForOrganizationRuleStatusDeleteSuccessful(ctx context.Context, conn *configservice.ConfigService, name string, timeout time.Duration) error {
	stateChangeConf := &retry.StateChangeConf{
		Pending: []string{configservice.OrganizationRuleStatusDeleteInProgress},
		Target:  []string{configservice.OrganizationRuleStatusDeleteSuccessful},
		Refresh: refreshOrganizationConfigRuleStatus(ctx, conn, name, configservice.OrganizationRuleStatusDeleteSuccessful),
		Timeout: timeout,
		Delay:   10 * time.Second,
	}

	_, err := stateChangeConf.WaitForStateContext(ctx)

	if tfawserr.ErrCodeEquals(err, configservice.ErrCodeNoSuchOrganizationConfigRuleException) {
		return nil
	}

	return err
}

func waitForOrganizationRuleStatusUpdateSuccessful(ctx context.Context, conn *configservice.ConfigService, name string, timeout time.Duration) error {
	stateChangeConf := &retry.StateChangeConf{
		Pending: []string{configservice.OrganizationRuleStatusUpdateInProgress},
		Target:  []string{configservice.OrganizationRuleStatusUpdateSuccessful},
		Refresh: refreshOrganizationConfigRuleStatus(ctx, conn, name, configservice.OrganizationRuleStatusUpdateSuccessful),
		Timeout: timeout,
		Delay:   10 * time.Second,
	}

	_, err := stateChangeConf.WaitForStateContext(ctx)

	return err
}
