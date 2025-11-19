package network

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cfTypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/sirupsen/logrus"
)

type NetworkService interface {
	CreateStack(templateFile *string, templateBody *[]byte, params map[string]string, tags map[string]string) error
}

type network struct {
}

var _ NetworkService = &network{}

func NewNetworkService() NetworkService {
	return &network{}
}

// CreateStack creates a CloudFormation stack
func (s *network) CreateStack(templateFile *string, templateBody *[]byte,
	params map[string]string, tags map[string]string) error {
	// Load the AWS configuration
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(params["Region"]))
	if err != nil {
		return fmt.Errorf("unable to load SDK config, %v", err)
	}

	var cfTags []cfTypes.Tag
	for k, v := range tags {
		cfTags = append(cfTags, cfTypes.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}

	// Create a CloudFormation client
	logger.Info("Creating CloudFormation client")
	cfClient := cloudformation.NewFromConfig(cfg)

	// Create a slice for CloudFormation parameters
	var cfParams []cfTypes.Parameter
	for k, v := range params {
		cfParams = append(cfParams, cfTypes.Parameter{
			ParameterKey:   aws.String(k),
			ParameterValue: aws.String(v),
		})
	}

	var template string
	template = *templateFile
	if templateBody != nil && len(*templateBody) > 0 {
		template = string(*templateBody)
	}

	// Create the stack
	logger.Info("Creating CloudFormation stack")
	_, err = cfClient.CreateStack(context.TODO(), &cloudformation.CreateStackInput{
		StackName:    aws.String(params["Name"]),
		TemplateBody: aws.String(template),
		Parameters:   cfParams,
		Tags:         cfTags,
		Capabilities: []cfTypes.Capability{
			cfTypes.CapabilityCapabilityIam,
			cfTypes.CapabilityCapabilityNamedIam,
			cfTypes.CapabilityCapabilityAutoExpand,
		},
	})
	if err != nil {
		deleteHelperMessage(logger, params, err)
		return fmt.Errorf("failed to create stack, %v", err)
	}

	// Fetch and log stack events periodically
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			<-ticker.C
			logStackEvents(cfClient, params["Name"], logger)
		}
	}()

	// Wait until the stack is created
	waiter := cloudformation.NewStackCreateCompleteWaiter(cfClient)
	err = waiter.Wait(context.TODO(), &cloudformation.DescribeStacksInput{
		StackName: aws.String(params["Name"]),
	}, 10*time.Minute, func(o *cloudformation.StackCreateCompleteWaiterOptions) {
		o.MinDelay = 30 * time.Second
		o.MaxDelay = 60 * time.Second
	})
	if err != nil {
		deleteHelperMessage(logger, params, err)
		helperMsg := ManualModeHelperMessage(params, tags)
		logger.Infof(helperMsg)
		return fmt.Errorf("failed to wait for stack creation, %v", err)
	}

	// Describe the stack resources
	describeStackResourcesOutput, err := cfClient.DescribeStackResources(context.TODO(),
		&cloudformation.DescribeStackResourcesInput{
			StackName: aws.String(params["Name"]),
		})
	if err != nil {
		return fmt.Errorf("failed to describe stack resources, %v", err)
	}

	logger.Info("--------------------------------")
	logger.Info("Resources created in stack:")
	for _, resource := range describeStackResourcesOutput.StackResources {
		logger.Infof("Resource: %s%s%s, Type: %s, ID: %s%s%s", ColorBlue,
			aws.ToString(resource.LogicalResourceId), ColorReset,
			aws.ToString(resource.ResourceType), ColorGreen,
			aws.ToString(resource.PhysicalResourceId), ColorReset)
	}

	logger.Infof("Stack %s created", params["Name"])
	return nil
}
