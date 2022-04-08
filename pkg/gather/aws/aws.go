package aws

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
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

	// Session is the AWS session to be used for deletion.  If nil, a
	// new session will be created based on the usual credential
	// configuration (AWS_PROFILE, AWS_ACCESS_KEY_ID, etc.).
	session *session.Session
}

// New returns a AWS Gather from ClusterMetadata.
func New(logger logrus.FieldLogger, serialLogBundle string, bootstrap string, masters []string, metadata *types.ClusterMetadata) (providers.Gather, error) {
	filters := make([]Filter, 0, len(metadata.ClusterPlatformMetadata.AWS.Identifier))
	for _, filter := range metadata.ClusterPlatformMetadata.AWS.Identifier {
		filters = append(filters, filter)
	}
	region := metadata.ClusterPlatformMetadata.AWS.Region
	session, err := awssession.GetSessionWithOptions(
		awssession.WithRegion(region),
		awssession.WithServiceEndpoints(region, metadata.ClusterPlatformMetadata.AWS.ServiceEndpoints),
	)
	if err != nil {
		return nil, err
	}

	return &Gather{
		logger:          logger,
		region:          region,
		filters:         filters,
		session:         session,
		serialLogBundle: serialLogBundle,
		bootstrap:       bootstrap,
		masters:         masters,
		directory:       filepath.Dir(serialLogBundle),
	}, nil
}

// Run is the entrypoint to start the gather process.
func (g *Gather) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	var err error
	awsSession := g.session
	if awsSession == nil {
		// Relying on appropriate AWS ENV vars (eg AWS_PROFILE, AWS_ACCESS_KEY_ID, etc)
		awsSession, err = session.NewSession(aws.NewConfig().WithRegion(g.region))
		if err != nil {
			return err
		}
	}
	awsSession.Handlers.Build.PushBackNamed(request.NamedHandler{
		Name: "openshiftInstaller.OpenshiftInstallerUserAgentHandler",
		Fn:   request.MakeAddToUserAgentHandler("OpenShift/4.x Destroyer", version.Raw),
	})

	ec2Client := ec2.New(awsSession)

	instances, err := g.getInstanceIDs(ctx, ec2Client)
	if err != nil {
		return err
	}

	if len(instances) == 0 {
		g.logger.Debugln("No instances found")
		return nil
	}

	serialLogBundleDir := strings.TrimSuffix(filepath.Base(g.serialLogBundle), ".tar.gz")
	filePathDir := filepath.Join(g.directory, serialLogBundleDir)
	err = os.MkdirAll(filePathDir, 0755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}

	g.logger.Debugf("found %d instances\n", len(instances))

	var errs []error
	var files []string
	for _, instance := range instances {
		filePath, err := g.downloadConsoleOutput(ctx, ec2Client, instance, filePathDir)
		if err != nil {
			g.logger.Debug(err)
			errs = append(errs, err)
		} else {
			files = append(files, filePath)
		}
	}

	if len(files) > 0 {
		err := gather.CreateArchive(files, g.serialLogBundle)
		if err != nil {
			g.logger.Debugf("failed to create archive: %s", err.Error())
			errs = append(errs, err)
		}
	}

	// clean up the mess we've made.
	_, err = os.Stat(filePathDir)
	if err == nil && !strings.HasPrefix(filePathDir, ".") {
		err := os.RemoveAll(filePathDir)
		if err != nil {
			errs = append(errs, errors.Wrap(err, "failed to removed archive directory"))
		}
	}

	return utilerrors.NewAggregate(errs)
}

// findEC2Instances returns the EC2 instances with tags that satisfy the filters.
// FIXME: use more appropriate name for this func
func (g *Gather) getInstanceIDs(ctx context.Context, ec2Client *ec2.EC2) ([]*ec2.Instance, error) {
	if ec2Client.Config.Region == nil {
		return nil, errors.New("EC2 client does not have region configured")
	}

	var instances []*ec2.Instance
	for _, filter := range g.filters {
		g.logger.Debugf("search for matching instances by tag in %s matching %#+v", *ec2Client.Config.Region, filter)
		instanceFilters := make([]*ec2.Filter, 0, len(g.filters))
		for key, value := range filter {
			instanceFilters = append(instanceFilters, &ec2.Filter{
				Name:   aws.String("tag:" + key),
				Values: []*string{aws.String(value)},
			})
		}

		err := ec2Client.DescribeInstancesPagesWithContext(
			ctx,
			&ec2.DescribeInstancesInput{Filters: instanceFilters},
			func(results *ec2.DescribeInstancesOutput, lastPage bool) bool {
				for _, reservation := range results.Reservations {
					if reservation.OwnerId == nil {
						continue
					}

					for _, instance := range reservation.Instances {
						// if instance.InstanceId == nil || instance.State == nil {
						// 	continue
						// }
						if instance.InstanceId != nil {
							instances = append(instances, instance)
						}
					}
				}
				return !lastPage
			},
		)
		if err != nil {
			err = errors.Wrap(err, "get ec2 instances")
			g.logger.Info(err)
			return instances, err
		}
	}

	return instances, nil
}

func (g *Gather) downloadConsoleOutput(ctx context.Context, ec2Client *ec2.EC2, instance *ec2.Instance, filePathDir string) (string, error) {
	logger := g.logger.WithField("instance", *instance.InstanceId)

	input := &ec2.GetConsoleOutputInput{
		InstanceId: aws.String(*instance.InstanceId),
	}
	result, err := ec2Client.GetConsoleOutput(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				logger.Errorln(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			logger.Errorln(err.Error())
		}
		return "", err
	}

	instanceName := *result.InstanceId
	for _, tags := range instance.Tags {
		if strings.EqualFold(*tags.Key, "Name") {
			instanceName = *tags.Value
		}
	}

	logger.Debugf("attemping to download console logs for %s", instanceName)
	filePath, err := g.saveToFile(instanceName, *result.Output, filePathDir)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func (g *Gather) saveToFile(instanceName, content, filePathDir string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode console output")
	}

	filePath := filepath.Join(filePathDir, fmt.Sprintf("%s-console.log", instanceName))

	file, err := os.Create(filePath)
	if err != nil {
		return "", errors.Wrap(err, "failed to create file")
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return "", errors.Wrap(err, "failed to wrtie to file")
	}

	return filePath, nil
}
