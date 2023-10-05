/*
Copyright 2018 The Kubernetes Authors.

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
	"encoding/base64"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"k8s.io/utils/pointer"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/userdata"
)

// GetLaunchTemplate returns the existing LaunchTemplate or nothing if it doesn't exist.
// For now by name until we need the input to be something different.
func (s *Service) GetLaunchTemplate(launchTemplateName string) (*expinfrav1.AWSLaunchTemplate, string, error) {
	if launchTemplateName == "" {
		return nil, "", nil
	}

	s.scope.V(2).Info("Looking for existing LaunchTemplates")

	input := &ec2.DescribeLaunchTemplateVersionsInput{
		LaunchTemplateName: aws.String(launchTemplateName),
		Versions:           aws.StringSlice([]string{expinfrav1.LaunchTemplateLatestVersion}),
	}

	out, err := s.EC2Client.DescribeLaunchTemplateVersions(input)
	switch {
	case awserrors.IsNotFound(err):
		return nil, "", nil
	case err != nil:
		return nil, "", err
	}

	if out == nil || out.LaunchTemplateVersions == nil || len(out.LaunchTemplateVersions) == 0 {
		return nil, "", nil
	}

	return s.SDKToLaunchTemplate(out.LaunchTemplateVersions[0])
}

// GetLaunchTemplateID returns the existing LaunchTemplateId or empty string if it doesn't exist.
func (s *Service) GetLaunchTemplateID(launchTemplateName string) (string, error) {
	if launchTemplateName == "" {
		return "", nil
	}

	input := &ec2.DescribeLaunchTemplateVersionsInput{
		LaunchTemplateName: aws.String(launchTemplateName),
		Versions:           aws.StringSlice([]string{expinfrav1.LaunchTemplateLatestVersion}),
	}

	out, err := s.EC2Client.DescribeLaunchTemplateVersions(input)
	switch {
	case awserrors.IsNotFound(err):
		return "", nil
	case err != nil:
		s.scope.Info("", "aerr", err.Error())
		return "", err
	}

	if out == nil || out.LaunchTemplateVersions == nil || len(out.LaunchTemplateVersions) == 0 {
		return "", nil
	}

	return aws.StringValue(out.LaunchTemplateVersions[0].LaunchTemplateId), nil
}

// CreateLaunchTemplate generates a launch template to be used with the autoscaling group.
func (s *Service) CreateLaunchTemplate(scope *scope.MachinePoolScope, imageID *string, userData []byte) (string, error) {
	s.scope.Info("Create a new launch template")

	launchTemplateData, err := s.createLaunchTemplateData(scope, imageID, userData)
	if err != nil {
		return "", errors.Wrapf(err, "unable to form launch template data")
	}

	input := &ec2.CreateLaunchTemplateInput{
		LaunchTemplateData: launchTemplateData,
		LaunchTemplateName: aws.String(scope.Name()),
	}

	additionalTags := scope.AdditionalTags()
	// Set the cloud provider tag
	additionalTags[infrav1.ClusterAWSCloudProviderTagKey(s.scope.Name())] = string(infrav1.ResourceLifecycleOwned)

	tags := infrav1.Build(infrav1.BuildParams{
		ClusterName: s.scope.Name(),
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		Name:        aws.String(scope.Name()),
		Role:        aws.String("node"),
		Additional:  additionalTags,
	})

	if len(tags) > 0 {
		spec := &ec2.TagSpecification{ResourceType: aws.String(ec2.ResourceTypeLaunchTemplate)}
		for key, value := range tags {
			spec.Tags = append(spec.Tags, &ec2.Tag{
				Key:   aws.String(key),
				Value: aws.String(value),
			})
		}
		input.TagSpecifications = append(input.TagSpecifications, spec)
	}

	result, err := s.EC2Client.CreateLaunchTemplate(input)
	if err != nil {
		return "", err
	}
	return aws.StringValue(result.LaunchTemplate.LaunchTemplateId), nil
}

// CreateLaunchTemplateVersion will create a launch template.
func (s *Service) CreateLaunchTemplateVersion(scope *scope.MachinePoolScope, imageID *string, userData []byte) error {
	s.scope.V(2).Info("creating new launch template version", "machine-pool", scope.Name())

	launchTemplateData, err := s.createLaunchTemplateData(scope, imageID, userData)
	if err != nil {
		return errors.Wrapf(err, "unable to form launch template data")
	}

	input := &ec2.CreateLaunchTemplateVersionInput{
		LaunchTemplateData: launchTemplateData,
		LaunchTemplateId:   aws.String(scope.AWSMachinePool.Status.LaunchTemplateID),
	}

	_, err = s.EC2Client.CreateLaunchTemplateVersion(input)
	if err != nil {
		return errors.Wrapf(err, "unable to create launch template version")
	}

	return nil
}

func (s *Service) createLaunchTemplateData(scope *scope.MachinePoolScope, imageID *string, userData []byte) (*ec2.RequestLaunchTemplateData, error) {
	lt := scope.AWSMachinePool.Spec.AWSLaunchTemplate

	// An explicit empty string for SSHKeyName means do not specify a key in the ASG launch
	var sshKeyNamePtr *string
	if lt.SSHKeyName != nil && *lt.SSHKeyName != "" {
		sshKeyNamePtr = lt.SSHKeyName
	}

	data := &ec2.RequestLaunchTemplateData{
		InstanceType: aws.String(lt.InstanceType),
		IamInstanceProfile: &ec2.LaunchTemplateIamInstanceProfileSpecificationRequest{
			Name: aws.String(lt.IamInstanceProfile),
		},
		KeyName:  sshKeyNamePtr,
		UserData: pointer.StringPtr(base64.StdEncoding.EncodeToString(userData)),
	}

	ids, err := s.GetCoreNodeSecurityGroups(scope)
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		data.SecurityGroupIds = append(data.SecurityGroupIds, aws.String(id))
	}

	// add additional security groups as well
	securityGroupIDs, err := s.GetAdditionalSecurityGroupsIDs(scope.AWSMachinePool.Spec.AWSLaunchTemplate.AdditionalSecurityGroups)
	if err != nil {
		return nil, err
	}
	data.SecurityGroupIds = append(data.SecurityGroupIds, aws.StringSlice(securityGroupIDs)...)

	// set the AMI ID
	data.ImageId = imageID

	// Set up root volume
	if lt.RootVolume != nil {
		rootDeviceName, err := s.checkRootVolume(lt.RootVolume, *data.ImageId)
		if err != nil {
			return nil, err
		}

		lt.RootVolume.DeviceName = aws.StringValue(rootDeviceName)

		req := volumeToLaunchTemplateBlockDeviceMappingRequest(lt.RootVolume)
		data.BlockDeviceMappings = []*ec2.LaunchTemplateBlockDeviceMappingRequest{
			req,
		}
	}

	data.TagSpecifications = s.buildLaunchTemplateTagSpecificationRequest(scope)

	return data, nil
}

func volumeToLaunchTemplateBlockDeviceMappingRequest(v *infrav1.Volume) *ec2.LaunchTemplateBlockDeviceMappingRequest {
	ltEbsDevice := &ec2.LaunchTemplateEbsBlockDeviceRequest{
		DeleteOnTermination: aws.Bool(true),
		VolumeSize:          aws.Int64(v.Size),
		Encrypted:           v.Encrypted,
	}

	if v.Throughput != nil {
		ltEbsDevice.Throughput = v.Throughput
	}

	if v.IOPS != 0 {
		ltEbsDevice.Iops = aws.Int64(v.IOPS)
	}

	if v.EncryptionKey != "" {
		ltEbsDevice.Encrypted = aws.Bool(true)
		ltEbsDevice.KmsKeyId = aws.String(v.EncryptionKey)
	}

	if v.Type != "" {
		ltEbsDevice.VolumeType = aws.String(string(v.Type))
	}

	return &ec2.LaunchTemplateBlockDeviceMappingRequest{
		DeviceName: &v.DeviceName,
		Ebs:        ltEbsDevice,
	}
}

// DeleteLaunchTemplate delete a launch template.
func (s *Service) DeleteLaunchTemplate(id string) error {
	s.scope.V(2).Info("Deleting launch template", "id", id)

	input := &ec2.DeleteLaunchTemplateInput{
		LaunchTemplateId: aws.String(id),
	}

	if _, err := s.EC2Client.DeleteLaunchTemplate(input); err != nil {
		return errors.Wrapf(err, "failed to delete launch template %q", id)
	}

	s.scope.V(2).Info("Deleted launch template", "id", id)
	return nil
}

// PruneLaunchTemplateVersions deletes one old launch template version.
// It does not delete the "latest" version, because that version may still be in use.
// It does not delete the "default" version, because that version cannot be deleted.
// It does not assume that versions are sequential. Versions may be deleted out of band.
func (s *Service) PruneLaunchTemplateVersions(id string) error {
	// When there is one version available, it is the default and the latest.
	// When there are two versions available, one the is the default, the other is the latest.
	// Therefore we only prune when there are at least 3 versions available.
	const minCountToAllowPrune = 3

	input := &ec2.DescribeLaunchTemplateVersionsInput{
		LaunchTemplateId: aws.String(id),
		MinVersion:       aws.String("0"),
		MaxVersion:       aws.String(expinfrav1.LaunchTemplateLatestVersion),
		MaxResults:       aws.Int64(minCountToAllowPrune),
	}

	out, err := s.EC2Client.DescribeLaunchTemplateVersions(input)
	if err != nil {
		s.scope.Info("", "aerr", err.Error())
		return err
	}

	// len(out.LaunchTemplateVersions)	|	items
	// -------------------------------- + -----------------------
	// 								1	|	[default/latest]
	// 								2	|	[default, latest]
	// 								3	| 	[default, versionToPrune, latest]
	if len(out.LaunchTemplateVersions) < minCountToAllowPrune {
		return nil
	}
	versionToPrune := out.LaunchTemplateVersions[1].VersionNumber
	return s.deleteLaunchTemplateVersion(id, versionToPrune)
}

func (s *Service) deleteLaunchTemplateVersion(id string, version *int64) error {
	s.scope.V(2).Info("Deleting launch template version", "id", id)

	if version == nil {
		return errors.New("version is a nil pointer")
	}
	versions := []string{strconv.FormatInt(*version, 10)}

	input := &ec2.DeleteLaunchTemplateVersionsInput{
		LaunchTemplateId: aws.String(id),
		Versions:         aws.StringSlice(versions),
	}

	_, err := s.EC2Client.DeleteLaunchTemplateVersions(input)
	if err != nil {
		return err
	}

	s.scope.V(2).Info("Deleted launch template", "id", id, "version", *version)
	return nil
}

// SDKToLaunchTemplate converts an AWS EC2 SDK instance to the CAPA instance type.
func (s *Service) SDKToLaunchTemplate(d *ec2.LaunchTemplateVersion) (*expinfrav1.AWSLaunchTemplate, string, error) {
	v := d.LaunchTemplateData
	i := &expinfrav1.AWSLaunchTemplate{
		Name: aws.StringValue(d.LaunchTemplateName),
		AMI: infrav1.AMIReference{
			ID: v.ImageId,
		},
		IamInstanceProfile: aws.StringValue(v.IamInstanceProfile.Name),
		InstanceType:       aws.StringValue(v.InstanceType),
		SSHKeyName:         v.KeyName,
		VersionNumber:      d.VersionNumber,
	}

	// Extract IAM Instance Profile name from ARN
	if v.IamInstanceProfile != nil && v.IamInstanceProfile.Arn != nil {
		split := strings.Split(aws.StringValue(v.IamInstanceProfile.Arn), "instance-profile/")
		if len(split) > 1 && split[1] != "" {
			i.IamInstanceProfile = split[1]
		}
	}

	for _, id := range v.SecurityGroupIds {
		// FIXME(dlipovetsky): This will include the core security groups as well, making the
		// "Additional" a bit dishonest. However, including the core groups drastically simplifies
		// comparison with the incoming security groups.
		i.AdditionalSecurityGroups = append(i.AdditionalSecurityGroups, infrav1.AWSResourceReference{ID: id})
	}

	if v.UserData == nil {
		return i, userdata.ComputeHash(nil), nil
	}
	decodedUserData, err := base64.StdEncoding.DecodeString(*v.UserData)
	if err != nil {
		return nil, "", errors.Wrap(err, "unable to decode UserData")
	}

	return i, userdata.ComputeHash(decodedUserData), nil
}

// LaunchTemplateNeedsUpdate checks if a new launch template version is needed.
//
// FIXME(dlipovetsky): This check should account for changed userdata, but does not yet do so.
// Although userdata is stored in an EC2 Launch Template, it is not a field of AWSLaunchTemplate.
func (s *Service) LaunchTemplateNeedsUpdate(scope *scope.MachinePoolScope, incoming *expinfrav1.AWSLaunchTemplate, existing *expinfrav1.AWSLaunchTemplate) (bool, error) {
	if incoming.IamInstanceProfile != existing.IamInstanceProfile {
		return true, nil
	}

	if incoming.InstanceType != existing.InstanceType {
		return true, nil
	}

	incomingIDs, err := s.GetAdditionalSecurityGroupsIDs(incoming.AdditionalSecurityGroups)
	if err != nil {
		return false, err
	}

	coreIDs, err := s.GetCoreNodeSecurityGroups(scope)
	if err != nil {
		return false, err
	}

	incomingIDs = append(incomingIDs, coreIDs...)
	existingIDs, err := s.GetAdditionalSecurityGroupsIDs(existing.AdditionalSecurityGroups)
	if err != nil {
		return false, err
	}
	sort.Strings(incomingIDs)
	sort.Strings(existingIDs)

	if !cmp.Equal(incomingIDs, existingIDs) {
		return true, nil
	}

	return false, nil
}

// DiscoverLaunchTemplateAMI will discover the AMI launch template.
func (s *Service) DiscoverLaunchTemplateAMI(scope *scope.MachinePoolScope) (*string, error) {
	lt := scope.AWSMachinePool.Spec.AWSLaunchTemplate

	if lt.AMI.ID != nil {
		return lt.AMI.ID, nil
	}

	if scope.MachinePool.Spec.Template.Spec.Version == nil {
		err := errors.New("Either AWSMachinePool's spec.awslaunchtemplate.ami.id or MachinePool's spec.template.spec.version must be defined")
		s.scope.Error(err, "")
		return nil, err
	}

	var lookupAMI string
	var err error

	imageLookupFormat := lt.ImageLookupFormat
	if imageLookupFormat == "" {
		imageLookupFormat = scope.InfraCluster.ImageLookupFormat()
	}

	imageLookupOrg := lt.ImageLookupOrg
	if imageLookupOrg == "" {
		imageLookupOrg = scope.InfraCluster.ImageLookupOrg()
	}

	imageLookupBaseOS := lt.ImageLookupBaseOS
	if imageLookupBaseOS == "" {
		imageLookupBaseOS = scope.InfraCluster.ImageLookupBaseOS()
	}

	if scope.IsEKSManaged() && imageLookupFormat == "" && imageLookupOrg == "" && imageLookupBaseOS == "" {
		lookupAMI, err = s.eksAMILookup(*scope.MachinePool.Spec.Template.Spec.Version, scope.AWSMachinePool.Spec.AWSLaunchTemplate.AMI.EKSOptimizedLookupType)
		if err != nil {
			return nil, err
		}
	} else {
		lookupAMI, err = s.defaultAMIIDLookup(imageLookupFormat, imageLookupOrg, imageLookupBaseOS, *scope.MachinePool.Spec.Template.Spec.Version)
		if err != nil {
			return nil, err
		}
	}

	return aws.String(lookupAMI), nil
}

func (s *Service) GetAdditionalSecurityGroupsIDs(securityGroups []infrav1.AWSResourceReference) ([]string, error) {
	var additionalSecurityGroupsIDs []string

	for _, sg := range securityGroups {
		if sg.ID != nil {
			additionalSecurityGroupsIDs = append(additionalSecurityGroupsIDs, *sg.ID)
		} else if sg.Filters != nil {
			id, err := s.getFilteredSecurityGroupID(sg)
			if err != nil {
				return nil, err
			}

			additionalSecurityGroupsIDs = append(additionalSecurityGroupsIDs, id)
		}
	}

	return additionalSecurityGroupsIDs, nil
}

func (s *Service) buildLaunchTemplateTagSpecificationRequest(scope *scope.MachinePoolScope) []*ec2.LaunchTemplateTagSpecificationRequest {
	tagSpecifications := make([]*ec2.LaunchTemplateTagSpecificationRequest, 0)
	additionalTags := scope.AdditionalTags()
	// Set the cloud provider tag
	additionalTags[infrav1.ClusterAWSCloudProviderTagKey(s.scope.Name())] = string(infrav1.ResourceLifecycleOwned)

	tags := infrav1.Build(infrav1.BuildParams{
		ClusterName: s.scope.Name(),
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		Name:        aws.String(scope.Name()),
		Role:        aws.String("node"),
		Additional:  additionalTags,
	})

	if len(tags) > 0 {
		// tag instances
		spec := &ec2.LaunchTemplateTagSpecificationRequest{ResourceType: aws.String(ec2.ResourceTypeInstance)}
		for key, value := range tags {
			spec.Tags = append(spec.Tags, &ec2.Tag{
				Key:   aws.String(key),
				Value: aws.String(value),
			})
		}
		tagSpecifications = append(tagSpecifications, spec)

		// tag EBS volumes
		spec = &ec2.LaunchTemplateTagSpecificationRequest{ResourceType: aws.String(ec2.ResourceTypeVolume)}
		for key, value := range tags {
			spec.Tags = append(spec.Tags, &ec2.Tag{
				Key:   aws.String(key),
				Value: aws.String(value),
			})
		}
		tagSpecifications = append(tagSpecifications, spec)
	}
	return tagSpecifications
}

// getFilteredSecurityGroupID get security group ID using filters.
func (s *Service) getFilteredSecurityGroupID(securityGroup infrav1.AWSResourceReference) (string, error) {
	if securityGroup.Filters == nil {
		return "", nil
	}

	filters := []*ec2.Filter{}
	for _, f := range securityGroup.Filters {
		filters = append(filters, &ec2.Filter{Name: aws.String(f.Name), Values: aws.StringSlice(f.Values)})
	}

	sgs, err := s.EC2Client.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{Filters: filters})
	if err != nil {
		return "", err
	}

	if len(sgs.SecurityGroups) == 0 {
		return "", fmt.Errorf("failed to find security group matching filters: %q, reason: %w", filters, err)
	}

	return *sgs.SecurityGroups[0].GroupId, nil
}
