/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ec2

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/blang/semver"
	"github.com/pkg/errors"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/api/bootstrap/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/record"
	"sigs.k8s.io/cluster-api-provider-aws/v2/util/system"
)

const (
	// DefaultArchitectureTag is the default architecture used when the architcture can't be determined from instance type.
	DefaultArchitectureTag = Amd64ArchitectureTag

	// Amd64ArchitectureTag is the reference AWS uses for amd64 architecture images.
	Amd64ArchitectureTag = "x86_64"

	// Arm64ArchitectureTag is the reference AWS uses for arm64 architecture images.
	Arm64ArchitectureTag = "arm64"

	// DefaultMachineAMIOwnerID is a CNCF owned account. This was previously a VMWare owned account
	// but the project lost access to it. See the k8s.io repo for the infra definition.
	DefaultMachineAMIOwnerID = "819546954734"

	// ubuntuOwnerID is Ubuntu owned account. Please see:
	// https://ubuntu.com/server/docs/cloud-images/amazon-ec2
	ubuntuOwnerID = "099720109477"

	ubuntuOwnerIDUsGov = "513442679011"

	// Description regex for fetching Ubuntu AMIs for bastion host.
	ubuntuImageDescription = "Canonical??Ubuntu??24.04??amd64?noble?image*"

	// defaultMachineAMILookupBaseOS is the default base operating system to use
	// when looking up machine AMIs.
	defaultMachineAMILookupBaseOS = "ubuntu-24.04"

	// DefaultAmiNameFormat is defined in the build/ directory of this project.
	// The pattern is:
	// 1. the string value `capa-ami-`
	// 2. the baseOS of the AMI, for example: ubuntu-24.04, centos-7, amazon-2
	// 3. the kubernetes version as defined by the packages produced by kubernetes/release with or without v as a prefix, for example: 1.13.0, 1.12.5-mybuild.1, v1.17.3
	// 4. a `-` followed by any additional characters.
	DefaultAmiNameFormat = "capa-ami-{{.BaseOS}}-?{{.K8sVersion}}-*"

	// Amazon's AMI timestamp format.
	createDateTimestampFormat = "2006-01-02T15:04:05.000Z"

	// EKS AMI ID SSM Parameter name.
	eksAmiSSMParameterFormat = "/aws/service/eks/optimized-ami/%s/amazon-linux-2/recommended/image_id"

	// EKS AL2023 AMI ID SSM Parameter name.
	eksAmiAl2023SSMParameterFormat = "/aws/service/eks/optimized-ami/%s/amazon-linux-2023/x86_64/standard/recommended/image_id"

	// EKS ARM64 AMI ID SSM Parameter name.
	eksARM64AmiSSMParameterFormat = "/aws/service/eks/optimized-ami/%s/amazon-linux-2-arm64/recommended/image_id"

	// EKS ARM64 AL2023 AMI ID SSM Parameter name.
	eksARM64AmiAl2023SSMParameterFormat = "/aws/service/eks/optimized-ami/%s/amazon-linux-2023/arm64/standard/recommended/image_id"

	// EKS GPU AMI ID SSM Parameter name.
	eksGPUAmiSSMParameterFormat = "/aws/service/eks/optimized-ami/%s/amazon-linux-2-gpu/recommended/image_id"

	// EKS GPU AL2023 AMI ID SSM Parameter name.
	eksGPUAmiAl2023SSMParameterFormat = "/aws/service/eks/optimized-ami/%s/amazon-linux-2023/x86_64/nvidia/recommended/image_id"
)

// AMILookup contains the parameters used to template AMI names used for lookup.
type AMILookup struct {
	BaseOS     string
	K8sVersion string
}

// GenerateAmiName will generate an AMI name.
func GenerateAmiName(amiNameFormat, baseOS, kubernetesVersion string) (string, error) {
	amiNameParameters := AMILookup{baseOS, kubernetesVersion}
	// revert to default if not specified
	if amiNameFormat == "" {
		amiNameFormat = DefaultAmiNameFormat
	}
	var templateBytes bytes.Buffer
	template, err := template.New("amiName").Parse(amiNameFormat)
	if err != nil {
		return amiNameFormat, errors.Wrapf(err, "failed create template from string: %q", amiNameFormat)
	}
	err = template.Execute(&templateBytes, amiNameParameters)
	if err != nil {
		return amiNameFormat, errors.Wrapf(err, "failed to substitute string: %q", amiNameFormat)
	}
	return templateBytes.String(), nil
}

// Determine architecture based on instance type.
func (s *Service) pickArchitectureForInstanceType(instanceType string) (string, error) {
	descInstanceTypeInput := &ec2.DescribeInstanceTypesInput{
		InstanceTypes: []*string{&instanceType},
	}
	describeInstanceTypeResult, err := s.EC2Client.DescribeInstanceTypesWithContext(context.TODO(), descInstanceTypeInput)
	if err != nil {
		// if call to DescribeInstanceTypes fails due to permissions error, log a warning and return the default architecture.
		if awserrors.IsPermissionsError(err) {
			record.Warnf(s.scope.InfraCluster(), "FailedDescribeInstanceTypes", "insufficient permissions to describe instance types for instance type %q, falling back to the default architecture of %q: %v", instanceType, DefaultArchitectureTag, err)

			return DefaultArchitectureTag, nil
		}
		return "", errors.Wrapf(err, "failed to describe instance types for instance type %q", instanceType)
	}

	if len(describeInstanceTypeResult.InstanceTypes) == 0 {
		return "", fmt.Errorf("instance type result empty for type %q", instanceType)
	}

	supportedArchs := describeInstanceTypeResult.InstanceTypes[0].ProcessorInfo.SupportedArchitectures

	logger := s.scope.GetLogger().WithValues("instance type", instanceType, "supported architectures", supportedArchs)
	logger.Info("Obtained a list of supported architectures for instance type")

	// Loop over every supported architecture for the instance type
	architecture := ""
archCheck:
	for _, a := range supportedArchs {
		switch *a {
		case Amd64ArchitectureTag:
			architecture = *a
			break archCheck
		case Arm64ArchitectureTag:
			architecture = *a
			break archCheck
		}
	}

	if architecture == "" {
		return "", fmt.Errorf("unable to find preferred architecture for instance type %q", instanceType)
	}

	logger.Info("Chosen architecture", "architecture", architecture)

	return architecture, nil
}

// DefaultAMILookup will do a default AMI lookup.
func DefaultAMILookup(ec2Client ec2iface.EC2API, ownerID, baseOS, kubernetesVersion, architecture, amiNameFormat string) (*ec2.Image, error) {
	if amiNameFormat == "" {
		amiNameFormat = DefaultAmiNameFormat
	}
	if ownerID == "" {
		ownerID = DefaultMachineAMIOwnerID
	}
	if baseOS == "" {
		baseOS = defaultMachineAMILookupBaseOS
	}

	amiName, err := GenerateAmiName(amiNameFormat, baseOS, strings.TrimPrefix(kubernetesVersion, "v"))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to process ami format: %q", amiNameFormat)
	}
	describeImageInput := &ec2.DescribeImagesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("owner-id"),
				Values: []*string{aws.String(ownerID)},
			},
			{
				Name:   aws.String("name"),
				Values: []*string{aws.String(amiName)},
			},
			{
				Name:   aws.String("architecture"),
				Values: []*string{aws.String(architecture)},
			},
			{
				Name:   aws.String("state"),
				Values: []*string{aws.String("available")},
			},
			{
				Name:   aws.String("virtualization-type"),
				Values: []*string{aws.String("hvm")},
			},
		},
	}

	out, err := ec2Client.DescribeImagesWithContext(context.TODO(), describeImageInput)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find ami: %q", amiName)
	}
	if out == nil || len(out.Images) == 0 {
		return nil, errors.Errorf("found no AMIs with the name: %q", amiName)
	}
	latestImage, err := GetLatestImage(out.Images)
	if err != nil {
		return nil, err
	}

	return latestImage, nil
}

// defaultAMIIDLookup returns the default AMI based on region.
func (s *Service) defaultAMIIDLookup(amiNameFormat, ownerID, baseOS, architecture, kubernetesVersion string) (string, error) {
	latestImage, err := DefaultAMILookup(s.EC2Client, ownerID, baseOS, kubernetesVersion, architecture, amiNameFormat)
	if err != nil {
		record.Eventf(s.scope.InfraCluster(), "FailedDescribeImages", "Failed to find ami for OS=%s, Architecture=%s and Kubernetes-version=%s: %v", baseOS, architecture, kubernetesVersion, err)
		return "", errors.Wrapf(err, "failed to find ami")
	}

	s.scope.Debug("Found and using an existing AMI", "ami-id", aws.StringValue(latestImage.ImageId))
	return aws.StringValue(latestImage.ImageId), nil
}

type images []*ec2.Image

// Len is the number of elements in the collection.
func (i images) Len() int {
	return len(i)
}

// Less reports whether the element with
// index i should sort before the element with index j.
// At this point all CreationDates have been checked for errors so ignoring the error is ok.
func (i images) Less(k, j int) bool {
	firstTime, _ := time.Parse(createDateTimestampFormat, aws.StringValue(i[k].CreationDate))
	secondTime, _ := time.Parse(createDateTimestampFormat, aws.StringValue(i[j].CreationDate))
	return firstTime.Before(secondTime)
}

// Swap swaps the elements with indexes i and j.
func (i images) Swap(k, j int) {
	i[k], i[j] = i[j], i[k]
}

// GetLatestImage assumes imgs is not empty. Responsibility of the caller to check.
func GetLatestImage(imgs []*ec2.Image) (*ec2.Image, error) {
	for _, img := range imgs {
		if _, err := time.Parse(createDateTimestampFormat, aws.StringValue(img.CreationDate)); err != nil {
			return nil, err
		}
	}
	// old to new (newest one is last)
	sort.Sort(images(imgs))
	return imgs[len(imgs)-1], nil
}

func (s *Service) defaultBastionAMILookup() (string, error) {
	describeImageInput := &ec2.DescribeImagesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("architecture"),
				Values: []*string{aws.String("x86_64")},
			},
			{
				Name:   aws.String("state"),
				Values: []*string{aws.String("available")},
			},
			{
				Name:   aws.String("virtualization-type"),
				Values: []*string{aws.String("hvm")},
			},
			{
				Name:   aws.String("description"),
				Values: aws.StringSlice([]string{ubuntuImageDescription}),
			},
		},
	}

	ownerID := ubuntuOwnerID
	partition := system.GetPartitionFromRegion(s.scope.Region())
	if strings.Contains(partition, v1beta1.PartitionNameUSGov) {
		ownerID = ubuntuOwnerIDUsGov
	}

	filter := &ec2.Filter{
		Name:   aws.String("owner-id"),
		Values: []*string{aws.String(ownerID)},
	}
	describeImageInput.Filters = append(describeImageInput.Filters, filter)

	out, err := s.EC2Client.DescribeImagesWithContext(context.TODO(), describeImageInput)
	if err != nil {
		return "", errors.Wrapf(err, "failed to describe images within region: %q", s.scope.Region())
	}
	if len(out.Images) == 0 {
		return "", errors.Errorf("found no AMIs within the region: %q", s.scope.Region())
	}
	latestImage, err := GetLatestImage(out.Images)
	if err != nil {
		return "", err
	}
	return *latestImage.ImageId, nil
}

func (s *Service) eksAMILookup(ctx context.Context, kubernetesVersion string, architecture string, amiType *infrav1.EKSAMILookupType) (string, error) {
	// format ssm parameter path properly
	formattedVersion, err := formatVersionForEKS(kubernetesVersion)
	if err != nil {
		return "", err
	}

	var paramName string

	if amiType == nil {
		amiType = new(infrav1.EKSAMILookupType)
	}

	switch *amiType {
	case infrav1.AmazonLinux2023GPU:
		paramName = fmt.Sprintf(eksGPUAmiAl2023SSMParameterFormat, formattedVersion)
	case infrav1.AmazonLinuxGPU:
		paramName = fmt.Sprintf(eksGPUAmiSSMParameterFormat, formattedVersion)
	default:
		switch architecture {
		case Arm64ArchitectureTag:
			if *amiType == infrav1.AmazonLinux2023 {
				paramName = fmt.Sprintf(eksARM64AmiAl2023SSMParameterFormat, formattedVersion)
			} else {
				paramName = fmt.Sprintf(eksARM64AmiSSMParameterFormat, formattedVersion)
			}
		case Amd64ArchitectureTag:
			if *amiType == infrav1.AmazonLinux2023 {
				paramName = fmt.Sprintf(eksAmiAl2023SSMParameterFormat, formattedVersion)
			} else {
				paramName = fmt.Sprintf(eksAmiSSMParameterFormat, formattedVersion)
			}
		default:
			return "", fmt.Errorf("cannot look up eks-optimized image for architecture %q", architecture)
		}
	}

	input := &ssm.GetParameterInput{
		Name: aws.String(paramName),
	}

	out, err := s.SSMClient.GetParameter(ctx, input)
	if err != nil {
		record.Eventf(s.scope.InfraCluster(), "FailedGetParameter", "Failed to get ami SSM parameter %q: %v", paramName, err)

		return "", errors.Wrapf(err, "failed to get ami SSM parameter: %q", paramName)
	}

	if out.Parameter == nil || out.Parameter.Value == nil {
		return "", errors.Errorf("SSM parameter returned with nil value: %q", paramName)
	}

	id := aws.StringValue(out.Parameter.Value)
	s.scope.Info("found AMI", "id", id, "version", formattedVersion)

	return id, nil
}

func formatVersionForEKS(version string) (string, error) {
	parsed, err := semver.ParseTolerant(version)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d.%d", parsed.Major, parsed.Minor), nil
}
