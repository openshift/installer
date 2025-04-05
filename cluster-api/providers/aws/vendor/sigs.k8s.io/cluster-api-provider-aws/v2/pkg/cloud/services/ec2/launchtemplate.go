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
	"context"
	"encoding/base64"
	"encoding/json"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apimachinerytypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/userdata"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/record"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
)

const (
	// TagsLastAppliedAnnotation is the key for the AWSMachinePool object annotation
	// which tracks the tags that the AWSMachinePool actuator is responsible
	// for. These are the tags that have been handled by the
	// AdditionalTags in the AWSMachinePool Provider Config.
	// See https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/
	// for annotation formatting rules.
	TagsLastAppliedAnnotation = "sigs.k8s.io/cluster-api-provider-aws-last-applied-tags"
)

// ReconcileLaunchTemplate reconciles a launch template and triggers instance refresh conditionally, depending on
// changes.
//
//nolint:gocyclo
func (s *Service) ReconcileLaunchTemplate(
	scope scope.LaunchTemplateScope,
	ec2svc services.EC2Interface,
	canUpdateLaunchTemplate func() (bool, error),
	runPostLaunchTemplateUpdateOperation func() error,
) error {
	bootstrapData, bootstrapDataSecretKey, err := scope.GetRawBootstrapData()
	if err != nil {
		record.Eventf(scope.GetMachinePool(), corev1.EventTypeWarning, "FailedGetBootstrapData", err.Error())
		return err
	}
	bootstrapDataHash := userdata.ComputeHash(bootstrapData)

	scope.Info("checking for existing launch template")
	launchTemplate, launchTemplateUserDataHash, launchTemplateUserDataSecretKey, err := ec2svc.GetLaunchTemplate(scope.LaunchTemplateName())
	if err != nil {
		conditions.MarkUnknown(scope.GetSetter(), expinfrav1.LaunchTemplateReadyCondition, expinfrav1.LaunchTemplateNotFoundReason, "%s", err.Error())
		return err
	}

	imageID, err := ec2svc.DiscoverLaunchTemplateAMI(scope)
	if err != nil {
		conditions.MarkFalse(scope.GetSetter(), expinfrav1.LaunchTemplateReadyCondition, expinfrav1.LaunchTemplateCreateFailedReason, clusterv1.ConditionSeverityError, "%s", err.Error())
		return err
	}

	if launchTemplate == nil {
		scope.Info("no existing launch template found, creating")
		launchTemplateID, err := ec2svc.CreateLaunchTemplate(scope, imageID, *bootstrapDataSecretKey, bootstrapData)
		if err != nil {
			conditions.MarkFalse(scope.GetSetter(), expinfrav1.LaunchTemplateReadyCondition, expinfrav1.LaunchTemplateCreateFailedReason, clusterv1.ConditionSeverityError, "%s", err.Error())
			return err
		}

		scope.SetLaunchTemplateIDStatus(launchTemplateID)
		return scope.PatchObject()
	}

	// LaunchTemplateID is set during LaunchTemplate creation, but for a scenario such as `clusterctl move`, status fields become blank.
	// If launchTemplate already exists but LaunchTemplateID field in the status is empty, get the ID and update the status.
	if scope.GetLaunchTemplateIDStatus() == "" {
		launchTemplateID, err := ec2svc.GetLaunchTemplateID(scope.LaunchTemplateName())
		if err != nil {
			conditions.MarkUnknown(scope.GetSetter(), expinfrav1.LaunchTemplateReadyCondition, expinfrav1.LaunchTemplateNotFoundReason, "%s", err.Error())
			return err
		}
		scope.SetLaunchTemplateIDStatus(launchTemplateID)
		return scope.PatchObject()
	}

	if scope.GetLaunchTemplateLatestVersionStatus() == "" {
		launchTemplateVersion, err := ec2svc.GetLaunchTemplateLatestVersion(scope.GetLaunchTemplateIDStatus())
		if err != nil {
			conditions.MarkUnknown(scope.GetSetter(), expinfrav1.LaunchTemplateReadyCondition, expinfrav1.LaunchTemplateNotFoundReason, "%s", err.Error())
			return err
		}
		scope.SetLaunchTemplateLatestVersionStatus(launchTemplateVersion)
		if err := scope.PatchObject(); err != nil {
			return err
		}
	}

	annotation, err := MachinePoolAnnotationJSON(scope, TagsLastAppliedAnnotation)
	if err != nil {
		return err
	}

	// Check if the instance tags were changed. If they were, create a new LaunchTemplate.
	tagsChanged, _, _, _ := tagsChanged(annotation, scope.AdditionalTags()) //nolint:dogsled

	needsUpdate, err := ec2svc.LaunchTemplateNeedsUpdate(scope, scope.GetLaunchTemplate(), launchTemplate)
	if err != nil {
		return err
	}

	amiChanged := *imageID != *launchTemplate.AMI.ID

	// `launchTemplateUserDataSecretKey` can be nil since it comes from a tag on the launch template
	// which may not exist in older launch templates created by older CAPA versions.
	// On change, we trigger instance refresh (rollout of new nodes). Therefore, do not consider it a change if the
	// launch template does not have the respective tag yet, as it could be surprising to users. Instead, ensure the
	// tag is stored on the newly-generated launch template version, without rolling out nodes.
	userDataSecretKeyChanged := launchTemplateUserDataSecretKey != nil && bootstrapDataSecretKey.String() != launchTemplateUserDataSecretKey.String()
	launchTemplateNeedsUserDataSecretKeyTag := launchTemplateUserDataSecretKey == nil

	if needsUpdate || tagsChanged || amiChanged || userDataSecretKeyChanged {
		canUpdate, err := canUpdateLaunchTemplate()
		if err != nil {
			return err
		}
		if !canUpdate {
			conditions.MarkFalse(scope.GetSetter(), expinfrav1.PreLaunchTemplateUpdateCheckCondition, expinfrav1.PreLaunchTemplateUpdateCheckFailedReason, clusterv1.ConditionSeverityWarning, "")
			return errors.New("Cannot update the launch template, prerequisite not met")
		}
	}

	userDataHashChanged := launchTemplateUserDataHash != bootstrapDataHash

	// Create a new launch template version if there's a difference in configuration, tags,
	// userdata, OR we've discovered a new AMI ID.
	if needsUpdate || tagsChanged || amiChanged || userDataHashChanged || userDataSecretKeyChanged || launchTemplateNeedsUserDataSecretKeyTag {
		scope.Info("creating new version for launch template", "existing", launchTemplate, "incoming", scope.GetLaunchTemplate(), "needsUpdate", needsUpdate, "tagsChanged", tagsChanged, "amiChanged", amiChanged, "userDataHashChanged", userDataHashChanged, "userDataSecretKeyChanged", userDataSecretKeyChanged)
		// There is a limit to the number of Launch Template Versions.
		// We ensure that the number of versions does not grow without bound by following a simple rule: Before we create a new version, we delete one old version, if there is at least one old version that is not in use.
		if err := ec2svc.PruneLaunchTemplateVersions(scope.GetLaunchTemplateIDStatus()); err != nil {
			return err
		}
		if err := ec2svc.CreateLaunchTemplateVersion(scope.GetLaunchTemplateIDStatus(), scope, imageID, *bootstrapDataSecretKey, bootstrapData); err != nil {
			return err
		}
		version, err := ec2svc.GetLaunchTemplateLatestVersion(scope.GetLaunchTemplateIDStatus())
		if err != nil {
			return err
		}

		scope.SetLaunchTemplateLatestVersionStatus(version)
		if err := scope.PatchObject(); err != nil {
			return err
		}
	}

	if needsUpdate || tagsChanged || amiChanged || userDataSecretKeyChanged {
		if err := runPostLaunchTemplateUpdateOperation(); err != nil {
			conditions.MarkFalse(scope.GetSetter(), expinfrav1.PostLaunchTemplateUpdateOperationCondition, expinfrav1.PostLaunchTemplateUpdateOperationFailedReason, clusterv1.ConditionSeverityError, "%s", err.Error())
			return err
		}
		conditions.MarkTrue(scope.GetSetter(), expinfrav1.PostLaunchTemplateUpdateOperationCondition)
	}

	return nil
}

// ReconcileTags reconciles the tags for the AWSMachinePool instances.
func (s *Service) ReconcileTags(scope scope.LaunchTemplateScope, resourceServicesToUpdate []scope.ResourceServiceToUpdate) error {
	additionalTags := scope.AdditionalTags()

	_, err := s.ensureTags(scope, resourceServicesToUpdate, additionalTags)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) ensureTags(scope scope.LaunchTemplateScope, resourceServicesToUpdate []scope.ResourceServiceToUpdate, additionalTags map[string]string) (bool, error) {
	annotation, err := MachinePoolAnnotationJSON(scope, TagsLastAppliedAnnotation)
	if err != nil {
		return false, err
	}

	// Check if the instance tags were changed. If they were, update them.
	// It would be possible here to only send new/updated tags, but for the
	// moment we send everything, even if only a single tag was created or
	// upated.
	changed, created, deleted, newAnnotation := tagsChanged(annotation, additionalTags)
	if changed {
		for _, resourceServiceToUpdate := range resourceServicesToUpdate {
			err := resourceServiceToUpdate.ResourceService.UpdateResourceTags(resourceServiceToUpdate.ResourceID, created, deleted)
			if err != nil {
				return false, err
			}
		}

		// We also need to update the annotation if anything changed.
		err = UpdateMachinePoolAnnotationJSON(scope, TagsLastAppliedAnnotation, newAnnotation)
		if err != nil {
			return false, err
		}
	}

	return changed, nil
}

// MachinePoolAnnotationJSON returns the annotation's json value as a map.
func MachinePoolAnnotationJSON(lts scope.LaunchTemplateScope, annotation string) (map[string]interface{}, error) {
	out := map[string]interface{}{}

	jsonAnnotation := machinePoolAnnotation(lts, annotation)
	if len(jsonAnnotation) == 0 {
		return out, nil
	}

	err := json.Unmarshal([]byte(jsonAnnotation), &out)
	if err != nil {
		return out, err
	}

	return out, nil
}

func machinePoolAnnotation(lts scope.LaunchTemplateScope, annotation string) string {
	return lts.GetObjectMeta().GetAnnotations()[annotation]
}

// UpdateMachinePoolAnnotationJSON updates the annotation with the given content.
func UpdateMachinePoolAnnotationJSON(lts scope.LaunchTemplateScope, annotation string, content map[string]interface{}) error {
	b, err := json.Marshal(content)
	if err != nil {
		return err
	}

	updateMachinePoolAnnotation(lts, annotation, string(b))
	return nil
}

func updateMachinePoolAnnotation(lts scope.LaunchTemplateScope, annotation, content string) {
	// Get the annotations
	annotations := lts.GetObjectMeta().GetAnnotations()

	if annotations == nil {
		annotations = make(map[string]string)
	}

	// Set our annotation to the given content.
	annotations[annotation] = content

	// Update the machine object with these annotations
	lts.GetObjectMeta().SetAnnotations(annotations)
}

// tagsChanged determines which tags to delete and which to add.
func tagsChanged(annotation map[string]interface{}, src map[string]string) (bool, map[string]string, map[string]string, map[string]interface{}) {
	// Bool tracking if we found any changed state.
	changed := false

	// Tracking for created/updated
	created := map[string]string{}

	// Tracking for tags that were deleted.
	deleted := map[string]string{}

	// The new annotation that we need to set if anything is created/updated.
	newAnnotation := map[string]interface{}{}

	// Loop over annotation, checking if entries are in src.
	// If an entry is present in annotation but not src, it has been deleted
	// since last time. We flag this in the deleted map.
	for t, v := range annotation {
		_, ok := src[t]

		// Entry isn't in src, it has been deleted.
		if !ok {
			// Cast v to a string here. This should be fine, tags are always
			// strings.
			deleted[t] = v.(string)
			changed = true
		}
	}

	// Loop over src, checking for entries in annotation.
	//
	// If an entry is in src, but not annotation, it has been created since
	// last time.
	//
	// If an entry is in both src and annotation, we compare their values, if
	// the value in src differs from that in annotation, the tag has been
	// updated since last time.
	for t, v := range src {
		av, ok := annotation[t]

		// Entries in the src always need to be noted in the newAnnotation. We
		// know they're going to be created or updated.
		newAnnotation[t] = v

		// Entry isn't in annotation, it's new.
		if !ok {
			created[t] = v
			newAnnotation[t] = v
			changed = true
			continue
		}

		// Entry is in annotation, has the value changed?
		if v != av {
			created[t] = v
			changed = true
		}

		// Entry existed in both src and annotation, and their values were
		// equal. Nothing to do.
	}

	// We made it through the loop, and everything that was in src, was also
	// in dst. Nothing changed.
	return changed, created, deleted, newAnnotation
}

// GetLaunchTemplate returns the existing LaunchTemplate or nothing if it doesn't exist.
// For now by name until we need the input to be something different.
func (s *Service) GetLaunchTemplate(launchTemplateName string) (*expinfrav1.AWSLaunchTemplate, string, *apimachinerytypes.NamespacedName, error) {
	if launchTemplateName == "" {
		return nil, "", nil, nil
	}

	s.scope.Debug("Looking for existing LaunchTemplates")

	input := &ec2.DescribeLaunchTemplateVersionsInput{
		LaunchTemplateName: aws.String(launchTemplateName),
		Versions:           aws.StringSlice([]string{expinfrav1.LaunchTemplateLatestVersion}),
	}

	out, err := s.EC2Client.DescribeLaunchTemplateVersionsWithContext(context.TODO(), input)
	switch {
	case awserrors.IsNotFound(err):
		return nil, "", nil, nil
	case err != nil:
		return nil, "", nil, err
	}

	if out == nil || out.LaunchTemplateVersions == nil || len(out.LaunchTemplateVersions) == 0 {
		return nil, "", nil, nil
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

	out, err := s.EC2Client.DescribeLaunchTemplateVersionsWithContext(context.TODO(), input)
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
func (s *Service) CreateLaunchTemplate(scope scope.LaunchTemplateScope, imageID *string, userDataSecretKey apimachinerytypes.NamespacedName, userData []byte) (string, error) {
	s.scope.Info("Create a new launch template")

	launchTemplateData, err := s.createLaunchTemplateData(scope, imageID, userDataSecretKey, userData)
	if err != nil {
		return "", errors.Wrapf(err, "unable to form launch template data")
	}

	input := &ec2.CreateLaunchTemplateInput{
		LaunchTemplateData: launchTemplateData,
		LaunchTemplateName: aws.String(scope.LaunchTemplateName()),
	}

	additionalTags := scope.AdditionalTags()
	// Set the cloud provider tag
	additionalTags[infrav1.ClusterAWSCloudProviderTagKey(s.scope.KubernetesClusterName())] = string(infrav1.ResourceLifecycleOwned)

	tags := infrav1.Build(infrav1.BuildParams{
		ClusterName: s.scope.KubernetesClusterName(),
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		Name:        aws.String(scope.LaunchTemplateName()),
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

	result, err := s.EC2Client.CreateLaunchTemplateWithContext(context.TODO(), input)
	if err != nil {
		return "", err
	}
	return aws.StringValue(result.LaunchTemplate.LaunchTemplateId), nil
}

// CreateLaunchTemplateVersion will create a launch template.
func (s *Service) CreateLaunchTemplateVersion(id string, scope scope.LaunchTemplateScope, imageID *string, userDataSecretKey apimachinerytypes.NamespacedName, userData []byte) error {
	s.scope.Debug("creating new launch template version", "machine-pool", scope.LaunchTemplateName())

	launchTemplateData, err := s.createLaunchTemplateData(scope, imageID, userDataSecretKey, userData)
	if err != nil {
		return errors.Wrapf(err, "unable to form launch template data")
	}

	input := &ec2.CreateLaunchTemplateVersionInput{
		LaunchTemplateData: launchTemplateData,
		LaunchTemplateId:   &id,
	}

	_, err = s.EC2Client.CreateLaunchTemplateVersionWithContext(context.TODO(), input)
	if err != nil {
		return errors.Wrapf(err, "unable to create launch template version")
	}

	return nil
}

func (s *Service) createLaunchTemplateData(scope scope.LaunchTemplateScope, imageID *string, userDataSecretKey apimachinerytypes.NamespacedName, userData []byte) (*ec2.RequestLaunchTemplateData, error) {
	lt := scope.GetLaunchTemplate()

	// An explicit empty string for SSHKeyName means do not specify a key in the ASG launch
	var sshKeyNamePtr *string
	if lt.SSHKeyName != nil && *lt.SSHKeyName != "" {
		sshKeyNamePtr = lt.SSHKeyName
	}

	data := &ec2.RequestLaunchTemplateData{
		InstanceType: aws.String(lt.InstanceType),
		KeyName:      sshKeyNamePtr,
		UserData:     ptr.To[string](base64.StdEncoding.EncodeToString(userData)),
	}

	if lt.InstanceMetadataOptions != nil {
		data.MetadataOptions = &ec2.LaunchTemplateInstanceMetadataOptionsRequest{
			HttpEndpoint:         aws.String(string(lt.InstanceMetadataOptions.HTTPEndpoint)),
			InstanceMetadataTags: aws.String(string(lt.InstanceMetadataOptions.InstanceMetadataTags)),
		}

		if lt.InstanceMetadataOptions.HTTPTokens != "" {
			data.MetadataOptions.HttpTokens = aws.String(string(lt.InstanceMetadataOptions.HTTPTokens))
		}
		if lt.InstanceMetadataOptions.HTTPPutResponseHopLimit != 0 {
			data.MetadataOptions.HttpPutResponseHopLimit = aws.Int64(lt.InstanceMetadataOptions.HTTPPutResponseHopLimit)
		}
	}

	if len(lt.IamInstanceProfile) > 0 {
		data.IamInstanceProfile = &ec2.LaunchTemplateIamInstanceProfileSpecificationRequest{
			Name: aws.String(lt.IamInstanceProfile),
		}
	}

	ids, err := s.GetCoreNodeSecurityGroups(scope)
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		data.SecurityGroupIds = append(data.SecurityGroupIds, aws.String(id))
	}

	// add additional security groups as well
	securityGroupIDs, err := s.GetAdditionalSecurityGroupsIDs(scope.GetLaunchTemplate().AdditionalSecurityGroups)
	if err != nil {
		return nil, err
	}
	data.SecurityGroupIds = append(data.SecurityGroupIds, aws.StringSlice(securityGroupIDs)...)

	// set the AMI ID
	data.ImageId = imageID

	instanceMarketOptions, err := getLaunchTemplateInstanceMarketOptionsRequest(scope.GetLaunchTemplate())
	if err != nil {
		return nil, err
	}
	data.InstanceMarketOptions = instanceMarketOptions
	data.PrivateDnsNameOptions = getLaunchTemplatePrivateDNSNameOptionsRequest(scope.GetLaunchTemplate().PrivateDNSName)

	blockDeviceMappings := []*ec2.LaunchTemplateBlockDeviceMappingRequest{}

	// Set up root volume
	if lt.RootVolume != nil {
		rootDeviceName, err := s.checkRootVolume(lt.RootVolume, *data.ImageId)
		if err != nil {
			return nil, err
		}

		lt.RootVolume.DeviceName = aws.StringValue(rootDeviceName)

		req := volumeToLaunchTemplateBlockDeviceMappingRequest(lt.RootVolume)
		blockDeviceMappings = append(blockDeviceMappings, req)
	}

	for vi := range lt.NonRootVolumes {
		nonRootVolume := lt.NonRootVolumes[vi]

		blockDeviceMapping := volumeToLaunchTemplateBlockDeviceMappingRequest(&nonRootVolume)
		blockDeviceMappings = append(blockDeviceMappings, blockDeviceMapping)
	}

	if len(blockDeviceMappings) > 0 {
		data.BlockDeviceMappings = blockDeviceMappings
	}

	data.TagSpecifications = s.buildLaunchTemplateTagSpecificationRequest(scope, userDataSecretKey)

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
	s.scope.Debug("Deleting launch template", "id", id)

	input := &ec2.DeleteLaunchTemplateInput{
		LaunchTemplateId: aws.String(id),
	}

	if _, err := s.EC2Client.DeleteLaunchTemplateWithContext(context.TODO(), input); err != nil {
		return errors.Wrapf(err, "failed to delete launch template %q", id)
	}

	s.scope.Debug("Deleted launch template", "id", id)
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

	out, err := s.EC2Client.DescribeLaunchTemplateVersionsWithContext(context.TODO(), input)
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

// GetLaunchTemplateLatestVersion returns the latest version of a launch template.
func (s *Service) GetLaunchTemplateLatestVersion(id string) (string, error) {
	input := &ec2.DescribeLaunchTemplateVersionsInput{
		LaunchTemplateId: aws.String(id),
		Versions:         aws.StringSlice([]string{expinfrav1.LaunchTemplateLatestVersion}),
	}

	out, err := s.EC2Client.DescribeLaunchTemplateVersionsWithContext(context.TODO(), input)
	if err != nil {
		s.scope.Info("", "aerr", err.Error())
		return "", err
	}

	if len(out.LaunchTemplateVersions) == 0 {
		return "", errors.Wrapf(err, "failed to get latest launch template version %q", id)
	}

	return strconv.Itoa(int(*out.LaunchTemplateVersions[0].VersionNumber)), nil
}

func (s *Service) deleteLaunchTemplateVersion(id string, version *int64) error {
	s.scope.Debug("Deleting launch template version", "id", id)

	if version == nil {
		return errors.New("version is a nil pointer")
	}
	versions := []string{strconv.FormatInt(*version, 10)}

	input := &ec2.DeleteLaunchTemplateVersionsInput{
		LaunchTemplateId: aws.String(id),
		Versions:         aws.StringSlice(versions),
	}

	_, err := s.EC2Client.DeleteLaunchTemplateVersionsWithContext(context.TODO(), input)
	if err != nil {
		return err
	}

	s.scope.Debug("Deleted launch template", "id", id, "version", *version)
	return nil
}

// SDKToLaunchTemplate converts an AWS EC2 SDK instance to the CAPA instance type.
func (s *Service) SDKToLaunchTemplate(d *ec2.LaunchTemplateVersion) (*expinfrav1.AWSLaunchTemplate, string, *apimachinerytypes.NamespacedName, error) {
	v := d.LaunchTemplateData
	i := &expinfrav1.AWSLaunchTemplate{
		Name: aws.StringValue(d.LaunchTemplateName),
		AMI: infrav1.AMIReference{
			ID: v.ImageId,
		},
		InstanceType:  aws.StringValue(v.InstanceType),
		SSHKeyName:    v.KeyName,
		VersionNumber: d.VersionNumber,
	}

	if v.MetadataOptions != nil {
		i.InstanceMetadataOptions = &infrav1.InstanceMetadataOptions{
			HTTPPutResponseHopLimit: aws.Int64Value(v.MetadataOptions.HttpPutResponseHopLimit),
			HTTPTokens:              infrav1.HTTPTokensState(aws.StringValue(v.MetadataOptions.HttpTokens)),
			HTTPEndpoint:            infrav1.InstanceMetadataEndpointStateEnabled,
			InstanceMetadataTags:    infrav1.InstanceMetadataEndpointStateDisabled,
		}
		if v.MetadataOptions.HttpEndpoint != nil && aws.StringValue(v.MetadataOptions.HttpEndpoint) == "disabled" {
			i.InstanceMetadataOptions.HTTPEndpoint = infrav1.InstanceMetadataEndpointStateDisabled
		}
		if v.MetadataOptions.InstanceMetadataTags != nil && aws.StringValue(v.MetadataOptions.InstanceMetadataTags) == "enabled" {
			i.InstanceMetadataOptions.InstanceMetadataTags = infrav1.InstanceMetadataEndpointStateEnabled
		}
	}

	if v.PrivateDnsNameOptions != nil {
		i.PrivateDNSName = &infrav1.PrivateDNSName{
			EnableResourceNameDNSAAAARecord: v.PrivateDnsNameOptions.EnableResourceNameDnsAAAARecord,
			EnableResourceNameDNSARecord:    v.PrivateDnsNameOptions.EnableResourceNameDnsARecord,
			HostnameType:                    v.PrivateDnsNameOptions.HostnameType,
		}
	}

	if v.IamInstanceProfile != nil {
		i.IamInstanceProfile = aws.StringValue(v.IamInstanceProfile.Name)
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
		return i, userdata.ComputeHash(nil), nil, nil
	}
	decodedUserData, err := base64.StdEncoding.DecodeString(*v.UserData)
	if err != nil {
		return nil, "", nil, errors.Wrap(err, "unable to decode UserData")
	}
	decodedUserDataHash := userdata.ComputeHash(decodedUserData)

	for _, tagSpecification := range v.TagSpecifications {
		if tagSpecification.ResourceType != nil && *tagSpecification.ResourceType == ec2.ResourceTypeInstance {
			for _, tag := range tagSpecification.Tags {
				if tag.Key != nil && *tag.Key == infrav1.LaunchTemplateBootstrapDataSecret && tag.Value != nil && strings.Contains(*tag.Value, "/") {
					parts := strings.SplitN(*tag.Value, "/", 2)
					launchTemplateUserDataSecretKey := &apimachinerytypes.NamespacedName{
						Namespace: parts[0],
						Name:      parts[1],
					}
					return i, decodedUserDataHash, launchTemplateUserDataSecretKey, nil
				}
			}
		}
	}

	return i, decodedUserDataHash, nil, nil
}

// LaunchTemplateNeedsUpdate checks if a new launch template version is needed.
//
// FIXME(dlipovetsky): This check should account for changed userdata, but does not yet do so.
// Although userdata is stored in an EC2 Launch Template, it is not a field of AWSLaunchTemplate.
func (s *Service) LaunchTemplateNeedsUpdate(scope scope.LaunchTemplateScope, incoming *expinfrav1.AWSLaunchTemplate, existing *expinfrav1.AWSLaunchTemplate) (bool, error) {
	if incoming.IamInstanceProfile != existing.IamInstanceProfile {
		return true, nil
	}

	if incoming.InstanceType != existing.InstanceType {
		return true, nil
	}
	if !cmp.Equal(incoming.InstanceMetadataOptions, existing.InstanceMetadataOptions) {
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
func (s *Service) DiscoverLaunchTemplateAMI(scope scope.LaunchTemplateScope) (*string, error) {
	lt := scope.GetLaunchTemplate()

	if lt.AMI.ID != nil {
		return lt.AMI.ID, nil
	}

	templateVersion := scope.GetMachinePool().Spec.Template.Spec.Version
	if templateVersion == nil {
		err := errors.New("Either AWSMachinePool's spec.awslaunchtemplate.ami.id or MachinePool's spec.template.spec.version must be defined")
		s.scope.Error(err, "")
		return nil, err
	}

	var lookupAMI string
	var err error

	imageLookupFormat := lt.ImageLookupFormat
	if imageLookupFormat == "" {
		imageLookupFormat = scope.GetEC2Scope().ImageLookupFormat()
	}

	imageLookupOrg := lt.ImageLookupOrg
	if imageLookupOrg == "" {
		imageLookupOrg = scope.GetEC2Scope().ImageLookupOrg()
	}

	imageLookupBaseOS := lt.ImageLookupBaseOS
	if imageLookupBaseOS == "" {
		imageLookupBaseOS = scope.GetEC2Scope().ImageLookupBaseOS()
	}

	instanceType := lt.InstanceType

	// If instance type is not specified on a launch template, we can safely assume the instance type will be a `t3.medium`.
	// As specified in the AWS docs https://docs.aws.amazon.com/eks/latest/userguide/launch-templates.html.
	// We will set the default architecture to `x86_64` as a result.
	imageArchitecture := Amd64ArchitectureTag

	if instanceType != "" {
		imageArchitecture, err = s.pickArchitectureForInstanceType(instanceType)
		if err != nil {
			return nil, err
		}
	}

	if scope.IsEKSManaged() && imageLookupFormat == "" && imageLookupOrg == "" && imageLookupBaseOS == "" {
		lookupAMI, err = s.eksAMILookup(
			*templateVersion,
			imageArchitecture,
			scope.GetLaunchTemplate().AMI.EKSOptimizedLookupType,
		)
		if err != nil {
			return nil, err
		}
	} else {
		lookupAMI, err = s.defaultAMIIDLookup(
			imageLookupFormat,
			imageLookupOrg,
			imageLookupBaseOS,
			imageArchitecture,
			*templateVersion,
		)
		if err != nil {
			return nil, err
		}
	}

	return aws.String(lookupAMI), nil
}

// GetAdditionalSecurityGroupsIDs returns the security group IDs for the additional security groups.
func (s *Service) GetAdditionalSecurityGroupsIDs(securityGroups []infrav1.AWSResourceReference) ([]string, error) {
	var additionalSecurityGroupsIDs []string

	for _, sg := range securityGroups {
		if sg.ID != nil {
			additionalSecurityGroupsIDs = append(additionalSecurityGroupsIDs, *sg.ID)
		} else if sg.Filters != nil {
			ids, err := s.getFilteredSecurityGroupIDs(sg)
			if err != nil {
				return nil, err
			}

			additionalSecurityGroupsIDs = append(additionalSecurityGroupsIDs, ids...)
		}
	}

	return additionalSecurityGroupsIDs, nil
}

func (s *Service) buildLaunchTemplateTagSpecificationRequest(scope scope.LaunchTemplateScope, userDataSecretKey apimachinerytypes.NamespacedName) []*ec2.LaunchTemplateTagSpecificationRequest {
	tagSpecifications := make([]*ec2.LaunchTemplateTagSpecificationRequest, 0)
	additionalTags := scope.AdditionalTags()
	// Set the cloud provider tag
	additionalTags[infrav1.ClusterAWSCloudProviderTagKey(s.scope.KubernetesClusterName())] = string(infrav1.ResourceLifecycleOwned)

	tags := infrav1.Build(infrav1.BuildParams{
		ClusterName: s.scope.KubernetesClusterName(),
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		Name:        aws.String(scope.LaunchTemplateName()),
		Role:        aws.String("node"),
		Additional:  additionalTags,
	})

	// tag instances
	{
		instanceTags := tags.DeepCopy()
		instanceTags[infrav1.LaunchTemplateBootstrapDataSecret] = userDataSecretKey.String()

		spec := &ec2.LaunchTemplateTagSpecificationRequest{ResourceType: aws.String(ec2.ResourceTypeInstance)}
		for key, value := range instanceTags {
			spec.Tags = append(spec.Tags, &ec2.Tag{
				Key:   aws.String(key),
				Value: aws.String(value),
			})
		}
		// Sort so that unit tests can expect a stable order
		sort.Slice(spec.Tags, func(i, j int) bool { return *spec.Tags[i].Key < *spec.Tags[j].Key })
		tagSpecifications = append(tagSpecifications, spec)
	}

	// tag EBS volumes
	if len(tags) > 0 {
		spec := &ec2.LaunchTemplateTagSpecificationRequest{ResourceType: aws.String(ec2.ResourceTypeVolume)}
		for key, value := range tags {
			spec.Tags = append(spec.Tags, &ec2.Tag{
				Key:   aws.String(key),
				Value: aws.String(value),
			})
		}
		// Sort so that unit tests can expect a stable order
		sort.Slice(spec.Tags, func(i, j int) bool { return *spec.Tags[i].Key < *spec.Tags[j].Key })
		tagSpecifications = append(tagSpecifications, spec)
	}

	return tagSpecifications
}

// getFilteredSecurityGroupIDs get security group IDs using filters.
func (s *Service) getFilteredSecurityGroupIDs(securityGroup infrav1.AWSResourceReference) ([]string, error) {
	if securityGroup.Filters == nil {
		return nil, nil
	}

	filters := []*ec2.Filter{}
	for _, f := range securityGroup.Filters {
		filters = append(filters, &ec2.Filter{Name: aws.String(f.Name), Values: aws.StringSlice(f.Values)})
	}

	sgs, err := s.EC2Client.DescribeSecurityGroupsWithContext(context.TODO(), &ec2.DescribeSecurityGroupsInput{Filters: filters})
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(sgs.SecurityGroups))
	for _, sg := range sgs.SecurityGroups {
		ids = append(ids, *sg.GroupId)
	}

	return ids, nil
}

func getLaunchTemplateInstanceMarketOptionsRequest(i *expinfrav1.AWSLaunchTemplate) (*ec2.LaunchTemplateInstanceMarketOptionsRequest, error) {
	if i.MarketType != "" && i.MarketType == infrav1.MarketTypeCapacityBlock && i.SpotMarketOptions != nil {
		return nil, errors.New("can't create spot capacity-blocks, remove spot market request")
	}

	// Infer MarketType if not explicitly set and SpotMarketOptions specified
	if i.SpotMarketOptions != nil && i.MarketType == "" {
		i.MarketType = infrav1.MarketTypeSpot
	}

	// Infer MarketType if not explicitly set
	if i.MarketType == "" {
		i.MarketType = infrav1.MarketTypeOnDemand
	}

	switch i.MarketType {
	case infrav1.MarketTypeCapacityBlock:
		// Handle Capacity Block case.
		if i.CapacityReservationID == nil {
			return nil, errors.Errorf("capacityReservationID is required when CapacityBlock is enabled")
		}
		return &ec2.LaunchTemplateInstanceMarketOptionsRequest{
			MarketType: aws.String(ec2.MarketTypeCapacityBlock),
		}, nil

	case infrav1.MarketTypeSpot:
		// Set required values for Spot instances
		spotOptions := &ec2.LaunchTemplateSpotMarketOptionsRequest{}

		// Persistent option is not available for EC2 autoscaling, EC2 makes a one-time request by default and setting request type should not be allowed.
		// For one-time requests, only terminate option is available as interruption behavior, and default for spotOptions.SetInstanceInterruptionBehavior() is terminate, so it is not set here explicitly.

		if maxPrice := aws.StringValue(i.SpotMarketOptions.MaxPrice); maxPrice != "" {
			spotOptions.SetMaxPrice(maxPrice)
		}

		launchTemplateInstanceMarketOptionsRequest := &ec2.LaunchTemplateInstanceMarketOptionsRequest{}
		launchTemplateInstanceMarketOptionsRequest.SetMarketType(ec2.MarketTypeSpot)
		launchTemplateInstanceMarketOptionsRequest.SetSpotOptions(spotOptions)

		return launchTemplateInstanceMarketOptionsRequest, nil
	case infrav1.MarketTypeOnDemand:
		// Instance is on-demand
		return nil, nil
	default:
		// Invalid MarketType provided
		return nil, errors.Errorf("invalid MarketType %s, must be spot/capacity-block or empty", i.MarketType)
	}
}

func getLaunchTemplatePrivateDNSNameOptionsRequest(privateDNSName *infrav1.PrivateDNSName) *ec2.LaunchTemplatePrivateDnsNameOptionsRequest {
	if privateDNSName == nil {
		return nil
	}

	return &ec2.LaunchTemplatePrivateDnsNameOptionsRequest{
		EnableResourceNameDnsAAAARecord: privateDNSName.EnableResourceNameDNSAAAARecord,
		EnableResourceNameDnsARecord:    privateDNSName.EnableResourceNameDNSARecord,
		HostnameType:                    privateDNSName.HostnameType,
	}
}
