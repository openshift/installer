package aws

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/smithy-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"

	awssession "github.com/openshift/installer/pkg/asset/installconfig/aws"
	"github.com/openshift/installer/pkg/gather"
	"github.com/openshift/installer/pkg/gather/providers"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/version"
)

// Filter holds the key/value pairs for the tags we will be matching against.
//
// A resource matches the filter if all of the key/value pairs are in its tags.
type Filter map[string]string

// Gather holds options for resources we want to gather.
type Gather struct {
	logger          logrus.FieldLogger
	filters         []Filter
	region          string
	bootstrap       string
	masters         []string
	directory       string
	serialLogBundle string
	ec2Client       *ec2.Client
}

// New returns an AWS Gather from ClusterMetadata.
func New(logger logrus.FieldLogger, serialLogBundle string, bootstrap string, masters []string, metadata *types.ClusterMetadata) (providers.Gather, error) {
	metadataAWS := metadata.ClusterPlatformMetadata.AWS

	filters := make([]Filter, 0, len(metadataAWS.Identifier))
	for _, filter := range metadataAWS.Identifier {
		filters = append(filters, filter)
	}

	ec2Client, err := awssession.NewEC2Client(context.TODO(), awssession.EndpointOptions{
		Region:    metadataAWS.Region,
		Endpoints: metadataAWS.ServiceEndpoints,
	}, ec2.WithAPIOptions(awsmiddleware.AddUserAgentKeyValue(awssession.OpenShiftInstallerGatherUserAgent, version.Raw)))
	if err != nil {
		return nil, fmt.Errorf("failed to create EC2 client: %w", err)
	}

	return &Gather{
		logger:          logger,
		region:          metadataAWS.Region,
		filters:         filters,
		serialLogBundle: serialLogBundle,
		bootstrap:       bootstrap,
		masters:         masters,
		directory:       filepath.Dir(serialLogBundle),
		ec2Client:       ec2Client,
	}, nil
}

// Run is the entrypoint to start the gather process.
func (g *Gather) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	instances, err := g.findEC2Instances(ctx)
	if err != nil {
		return err
	}

	if len(instances) == 0 {
		g.logger.Infoln("Skipping console log gathering: no instances found")
		return nil
	}

	serialLogBundleDir := strings.TrimSuffix(filepath.Base(g.serialLogBundle), ".tar.gz")
	filePathDir := filepath.Join(g.directory, serialLogBundleDir)
	err = os.MkdirAll(filePathDir, 0755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}

	var errs []error
	var files []string
	for _, instance := range instances {
		filePath, err := g.downloadConsoleOutput(ctx, instance, filePathDir)
		if err != nil {
			errs = append(errs, err)
		} else {
			files = append(files, filePath)
		}
	}

	if len(files) > 0 {
		err := gather.CreateArchive(files, g.serialLogBundle)
		if err != nil {
			errs = append(errs, errors.Wrap(err, "failed to create archive"))
		}
	}

	if err := gather.DeleteArchiveDirectory(filePathDir); err != nil {
		// Note: cleanup is best effort, it shouldn't fail the gather
		g.logger.Debugf("Failed to remove archive directory: %v", err)
	}

	return utilerrors.NewAggregate(errs)
}

// findEC2Instances returns the EC2 instances with tags that satisfy the filters.
func (g *Gather) findEC2Instances(ctx context.Context) ([]*ec2types.Instance, error) {
	region := g.ec2Client.Options().Region
	if region == "" {
		return nil, errors.New("EC2 client does not have region configured")
	}

	var instances []*ec2types.Instance
	for _, filter := range g.filters {
		g.logger.Debugf("Search for matching instances by tag in %s matching %#+v", region, filter)
		instanceFilters := make([]ec2types.Filter, 0, len(g.filters))

		for key, value := range filter {
			instanceFilters = append(instanceFilters, ec2types.Filter{
				Name:   aws.String("tag:" + key),
				Values: []string{value},
			})
		}

		input := &ec2.DescribeInstancesInput{Filters: instanceFilters}
		paginator := ec2.NewDescribeInstancesPaginator(g.ec2Client, input)

		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				return instances, fmt.Errorf("failed to get ec2 instances: %w", err)
			}

			for _, reservation := range page.Reservations {
				if reservation.OwnerId == nil {
					continue
				}

				for _, instance := range reservation.Instances {
					if instance.InstanceId != nil {
						instances = append(instances, &instance)
					}
				}
			}
		}
	}

	return instances, nil
}

// downloadConsoleOutput downloads console logs for the EC2 instance, saves it to local disk
// and returns the file name.
func (g *Gather) downloadConsoleOutput(ctx context.Context, instance *ec2types.Instance, filePathDir string) (string, error) {
	instanceName := aws.ToString(instance.InstanceId)
	for _, tags := range instance.Tags {
		if strings.EqualFold(aws.ToString(tags.Key), "Name") {
			instanceName = aws.ToString(tags.Value)
		}
	}

	logger := g.logger.WithField("Instance", instanceName)
	logger.Debugf("Attemping to download console logs for %s", instanceName)

	input := &ec2.GetConsoleOutputInput{InstanceId: instance.InstanceId}
	result, err := g.ec2Client.GetConsoleOutput(ctx, input)
	if err != nil {
		var aerr smithy.APIError
		if errors.As(err, &aerr) {
			logger.Errorf("failed to gather console logs for %s: %s", instanceName, aerr.ErrorMessage())
		}
		return "", err
	}

	filePath, err := g.saveToFile(instanceName, aws.ToString(result.Output), filePathDir)
	if err != nil {
		return "", err
	}
	logger.Debug("Download complete")

	return filePath, nil
}

func (g *Gather) saveToFile(instanceName, content, filePathDir string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode console output")
	}

	filename := filepath.Join(filePathDir, fmt.Sprintf("%s-serial.log", instanceName))

	file, err := os.Create(filename)
	if err != nil {
		return "", errors.Wrap(err, "failed to create file")
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return "", errors.Wrap(err, "failed to write to file")
	}

	return filename, nil
}
