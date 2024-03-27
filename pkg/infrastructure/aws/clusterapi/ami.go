package clusterapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset/installconfig"
)

// copyAMIToRegion copies the AMI to the region configured in the installConfig if needed.
func copyAMIToRegion(ctx context.Context, installConfig *installconfig.InstallConfig, infraID, rhcosImage string) (string, error) {
	osImage := strings.SplitN(rhcosImage, ",", 2)
	amiID := osImage[0]
	amiRegion := installConfig.Config.AWS.Region
	if len(osImage) > 1 {
		amiRegion = osImage[1]
	}

	// Already in target region, nothing to do
	if amiRegion == installConfig.Config.AWS.Region {
		return amiID, nil
	}

	logrus.Infof("Copying AMI to region %s", installConfig.AWS.Region)

	session, err := installConfig.AWS.Session(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get AWS session: %w", err)
	}
	client := ec2.New(session)

	res, err := client.CopyImageWithContext(ctx, &ec2.CopyImageInput{
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
	amiTags := make([]*ec2.Tag, 0, len(installConfig.Config.AWS.UserTags)+4)
	for k, v := range installConfig.Config.AWS.UserTags {
		amiTags = append(amiTags, &ec2.Tag{
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
		amiTags = append(amiTags, &ec2.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}

	_, err = client.CreateTagsWithContext(ctx, &ec2.CreateTagsInput{
		Resources: []*string{res.ImageId},
		Tags:      amiTags,
	})
	if err != nil {
		return "", fmt.Errorf("failed to tag AMI copy (%s): %w", name, err)
	}

	return aws.StringValue(res.ImageId), nil
}
