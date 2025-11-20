package network

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cfTypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/sirupsen/logrus"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
)

// logStackEvents fetches and logs stack events
func logStackEvents(cfClient *cloudformation.Client, stackName string, logger *logrus.Logger) {
	events, err := cfClient.DescribeStackEvents(context.TODO(), &cloudformation.DescribeStackEventsInput{
		StackName: aws.String(stackName),
	})
	if err != nil {
		logger.Errorf("Failed to describe stack events: %v", err)
		return
	}

	// Group events by resource and keep the latest event for each resource
	latestEvents := make(map[string]cfTypes.StackEvent)
	for _, event := range events.StackEvents {
		resource := aws.ToString(event.LogicalResourceId)
		if existingEvent, exists := latestEvents[resource]; !exists || event.Timestamp.After(*existingEvent.Timestamp) {
			latestEvents[resource] = event
		}
	}

	logger.Info("---------------------------------------------")
	// Log the latest event for each resource with color
	for resource, event := range latestEvents {
		statusColor := getStatusColor(event.ResourceStatus)
		reason := aws.ToString(event.ResourceStatusReason)

		readableReason := strings.ReplaceAll(reason, ". ", ".\n    ")

		statusStr := fmt.Sprintf("%s%s%s", statusColor, event.ResourceStatus, ColorReset)

		if strings.Contains(reason, "AccessDenied") {
			logger.Warnf("Resource: %s, Status: %s, Reason: %s (Access Denied)",
				resource, statusStr, readableReason)
		} else {
			logger.Infof("Resource: %s, Status: %s, Reason: %s",
				resource, statusStr, readableReason)
		}
	}
}

func getStatusColor(status cfTypes.ResourceStatus) string {
	switch status {
	case cfTypes.ResourceStatusCreateComplete, cfTypes.ResourceStatusUpdateComplete:
		return ColorGreen
	case cfTypes.ResourceStatusCreateFailed, cfTypes.ResourceStatusDeleteFailed, cfTypes.ResourceStatusUpdateFailed:
		return ColorRed
	default:
		return ColorYellow
	}
}

// ParseParams converts the list of parameter strings into a map and sets default values
func ParseParams(params []string) (map[string]string, map[string]string, error) {
	result := make(map[string]string)
	userTags := map[string]string{}

	for _, param := range params {
		parts := strings.SplitN(param, "=", 2)
		if len(parts) != 2 {
			return nil, nil, errors.New("invalid parameter format")
		}
		if parts[0] == "Tags" {
			tagEntries := strings.Split(parts[1], ",")
			for _, entry := range tagEntries {
				tagParts := strings.SplitN(strings.TrimSpace(entry), "=", 2)
				if len(tagParts) != 2 {
					return nil, nil, errors.New("invalid tag format")
				}
				if _, ok := userTags[tagParts[0]]; ok {
					return nil, nil, errors.New("duplicate tag key " + tagParts[0])
				}
				userTags[tagParts[0]] = tagParts[1]
			}
		} else {
			if _, ok := result[parts[0]]; ok {
				return nil, nil, errors.New("duplicate parameter key")
			}
			result[parts[0]] = parts[1]
		}
	}

	return result, userTags, nil
}

// SelectTemplate selects the appropriate template file based on the template name
func SelectTemplate(templateDir, command string) string {
	return fmt.Sprintf("%s/%s/cloudformation.yaml", templateDir, command)
}

func formatParams(params map[string]string) string {
	var paramStr string
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		paramStr += fmt.Sprintf("ParameterKey=%s,ParameterValue=%s ", k, params[k])
	}
	return paramStr
}

func formatTags(tags map[string]string) string {
	var tagStr string
	keys := make([]string, 0, len(tags))
	for k := range tags {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		tagStr += fmt.Sprintf("Key=%s,Value=%s ", k, tags[k])
	}
	return tagStr
}

func deleteHelperMessage(logger *logrus.Logger, params map[string]string, err error) {
	logger.Errorf("Failed to create CloudFormation stack: %v", err)
	logger.Infof("To delete all created resource stacks, run "+
		"`aws cloudformation delete-stack --stack-name %s --region %s`",
		params["Name"], params["Region"])
}

func ManualModeHelperMessage(params map[string]string, tags map[string]string) string {
	return fmt.Sprintf("Run the following command to create the stack manually:\n"+
		"aws cloudformation create-stack --stack-name %s --template-body file://<template-file-path>"+
		" --param %s --tags %s --region %s",
		params["Name"], formatParams(params), formatTags(tags), params["Region"])
}
