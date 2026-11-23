package clusterapi

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/rhcos"
)

// copyAMIToRegion copies the AMI to the region configured in the installConfig if needed.
func copyAMIToRegion(ctx context.Context, installConfig *installconfig.InstallConfig, infraID string, rhcosImage *rhcos.Image) (string, error) {
	osImage := strings.SplitN(rhcosImage.ControlPlane, ",", 2)
	amiID, amiRegion := osImage[0], osImage[1]

	logrus.Infof("Copying AMI %s to region %s", amiID, installConfig.AWS.Region)

	client, err := installConfig.AWS.EC2Client(ctx)
	if err != nil {
		return "", err
	}

	res, err := client.CopyImage(ctx, &ec2.CopyImageInput{
		Name:          aws.String(fmt.Sprintf("%s-master", infraID)),
		ClientToken:   aws.String(infraID),
		Description:   aws.String("Created by Openshift Installer"),
		SourceImageId: aws.String(amiID),
		SourceRegion:  aws.String(amiRegion),
		Encrypted:     aws.Bool(true),
	})
	if err != nil {
		return "", err
	}

	name := fmt.Sprintf("%s-ami-%s", infraID, installConfig.AWS.Region)
	amiTags := make([]ec2types.Tag, 0, len(installConfig.Config.AWS.UserTags)+4)
	for k, v := range installConfig.Config.AWS.UserTags {
		amiTags = append(amiTags, ec2types.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}
	for k, v := range map[string]string{
		"Name":         name,
		"sourceAMI":    amiID,
		"sourceRegion": amiRegion,
		fmt.Sprintf("kubernetes.io/cluster/%s", infraID): "owned",
	} {
		amiTags = append(amiTags, ec2types.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}

	_, err = client.CreateTags(ctx, &ec2.CreateTagsInput{
		Resources: []string{aws.ToString(res.ImageId)},
		Tags:      amiTags,
	})
	if err != nil {
		return "", fmt.Errorf("failed to tag AMI copy (%s): %w", name, err)
	}

	describeImagesInput := &ec2.DescribeImagesInput{
		ImageIds: []string{aws.ToString(res.ImageId)},
	}

	logrus.Infof("Waiting for AMI [%s] to exist in region [%s]",
		aws.ToString(res.ImageId), installConfig.AWS.Region)

	existsWaiter := ec2.NewImageExistsWaiter(client)
	if err := existsWaiter.Wait(ctx, describeImagesInput, 5*time.Minute); err != nil {
		return "", fmt.Errorf("failed to wait for AMI [%s] to exist: %w",
			aws.ToString(res.ImageId), err)
	}

	logrus.Infof("AMI [%s] exists, waiting for it to become available", aws.ToString(res.ImageId))

	// Wait for the AMI to transition from pending to available
	availableWaiter := ec2.NewImageAvailableWaiter(client)
	if err := availableWaiter.Wait(ctx, describeImagesInput, 10*time.Minute); err != nil {
		return "", fmt.Errorf("failed to wait for AMI [%s] to become available: %w",
			aws.ToString(res.ImageId), err)
	}

	logrus.Infof("AMI [%s] is now available", aws.ToString(res.ImageId))

	return aws.ToString(res.ImageId), nil
}
