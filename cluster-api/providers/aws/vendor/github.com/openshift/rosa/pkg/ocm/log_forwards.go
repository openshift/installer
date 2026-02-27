/*
Copyright (c) 2025 Red Hat, Inc.

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

package ocm

import (
	"reflect"

	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	errors "github.com/zgalor/weberr"

	"github.com/openshift/rosa/pkg/logforwarding"
)

func BuildLogForwarder(logForwarderConfig *cmv1.LogForwarder) *cmv1.LogForwarderBuilder {
	logForwardbldr := cmv1.NewLogForwarder()
	if logForwarderConfig != nil {
		if len(logForwarderConfig.Applications()) > 0 {
			logForwardbldr.Applications(logForwarderConfig.Applications()...)
		}
		if logForwarderConfig.Cloudwatch() != nil && (logForwarderConfig.Cloudwatch().LogGroupName() != "" ||
			logForwarderConfig.Cloudwatch().LogDistributionRoleArn() != "") {
			logForwardbldr.Cloudwatch(cmv1.NewLogForwarderCloudWatchConfig().
				LogDistributionRoleArn(logForwarderConfig.Cloudwatch().LogDistributionRoleArn()).
				LogGroupName(logForwarderConfig.Cloudwatch().LogGroupName()))
		}
		if len(logForwarderConfig.Groups()) > 0 {
			logForwarderGroupBlds := make([]*cmv1.LogForwarderGroupBuilder, 0)
			for _, group := range logForwarderConfig.Groups() {
				logForwarderGroupBlds = append(logForwarderGroupBlds, cmv1.NewLogForwarderGroup().
					Version(group.Version()).ID(group.ID()))
			}
			logForwardbldr.Groups(logForwarderGroupBlds...)
		}
		if logForwarderConfig.S3() != nil && (logForwarderConfig.S3().BucketName() != "" ||
			logForwarderConfig.S3().BucketPrefix() != "") {
			logForwardbldr.S3(cmv1.NewLogForwarderS3Config().BucketName(logForwarderConfig.S3().BucketName()).
				BucketPrefix(logForwarderConfig.S3().BucketPrefix()))
		}
	}

	return logForwardbldr
}

func (c *Client) GetLogForwarders(clusterId string) ([]*cmv1.LogForwarder, error) {
	var LogForwarderList []*cmv1.LogForwarder
	collection := c.ocm.ClustersMgmt().V1().
		Clusters().
		Cluster(clusterId).
		ControlPlane().LogForwarders().List()

	page := 1
	size := 100
	for {
		response, err := collection.
			Page(page).
			Size(size).
			Send()
		if err != nil {
			return nil, handleErr(response.Error(), err)
		}
		LogForwarderList = append(LogForwarderList, response.Items().Slice()...)
		if response.Size() < size {
			break
		}
		page++
	}

	return LogForwarderList, nil
}

func (c *Client) GetLogForwarderByID(clusterID string, logForwarderID string) (*cmv1.LogForwarder, error) {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().
		Cluster(clusterID).
		ControlPlane().
		LogForwarders().
		LogForwarder(logForwarderID).
		Get().
		Send()

	if err != nil {
		if response != nil && response.Status() == 404 {
			return nil, nil
		}
		return nil, handleErr(response.Error(), err)
	}
	return response.Body(), nil
}

func (c *Client) SetLogForwarder(clusterId string,
	logForwarder *cmv1.LogForwarder) (*cmv1.LogForwarder, error) {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterId).ControlPlane().
		LogForwarders().Add().Body(logForwarder).Send()

	if err != nil {
		return nil, handleErr(response.Error(), err)
	}
	return response.Body(), nil
}

func (c *Client) DeleteLogForwarder(clusterID string, logForwarderID string) error {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		ControlPlane().
		LogForwarders().
		LogForwarder(logForwarderID).
		Delete().
		Send()
	if err != nil {
		return handleErr(response.Error(), err)
	}
	return nil
}

func (c *Client) GetLogForwarderGroupVersions() ([]*cmv1.LogForwarderGroupVersions, error) {
	response, err := c.ocm.ClustersMgmt().V1().LogForwarding().Groups().List().Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}
	return response.Items().Slice(), nil
}

func (c *Client) FetchLogForwarder(clusterId string, logForwarderId string) (*cmv1.LogForwarder, error) {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().
		Cluster(clusterId).
		ControlPlane().
		LogForwarders().
		LogForwarder(logForwarderId).
		Get().
		Send()

	if err != nil {
		return nil, handleErr(response.Error(), err)
	}
	return response.Body(), nil
}

func (c *Client) UpdateLogForwarder(logForwarder *cmv1.LogForwarder, logForwarderId string, clusterId string) error {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().
		Cluster(clusterId).
		ControlPlane().
		LogForwarders().
		LogForwarder(logForwarderId).
		Update().
		Body(logForwarder).
		Send()

	if err != nil {
		return handleErr(response.Error(), err)
	}
	return nil
}

func (c *Client) EditLogForwarder(clusterId string, logForwarderId string,
	logForwarderYaml logforwarding.LogForwarderYaml, currentLogForwarder *cmv1.LogForwarder) error {

	if logForwarderYaml.S3 == nil && logForwarderYaml.CloudWatch == nil {
		return errors.UserErrorf("log forwarding config provided contained no valid log forwarders")
	}

	var logForwarderConfig *cmv1.LogForwarderBuilder

	if currentLogForwarder.S3() != nil && currentLogForwarder.S3().BucketName() != "" {
		logForwarderConfigBuilder := cmv1.NewLogForwarder()
		logForwarderS3ConfigBuilder := cmv1.NewLogForwarderS3Config()

		logForwarderS3ConfigBuilder.BucketPrefix(logForwarderYaml.S3.S3ConfigBucketPrefix)

		logForwarderS3ConfigBuilder.BucketName(logForwarderYaml.S3.S3ConfigBucketName)

		if !reflect.DeepEqual(currentLogForwarder.Applications(), logForwarderYaml.S3.Applications) {
			logForwarderConfigBuilder.Applications(logForwarderYaml.S3.Applications...)
		}

		if !reflect.DeepEqual(currentLogForwarder.Groups(), logForwarderYaml.S3.GroupsLogVersions) {
			var groupBuilder []*cmv1.LogForwarderGroupBuilder
			for _, group := range logForwarderYaml.S3.GroupsLogVersions {
				groupBuilder = append(groupBuilder, cmv1.NewLogForwarderGroup().ID(group))
			}
			logForwarderConfigBuilder.Groups(groupBuilder...)
		}

		logForwarderConfigBuilder.S3(logForwarderS3ConfigBuilder)
		logForwarderConfig = logForwarderConfigBuilder
	} else if currentLogForwarder.Cloudwatch() != nil && currentLogForwarder.Cloudwatch().
		LogDistributionRoleArn() != "" {

		logForwarderConfigBuilder := cmv1.NewLogForwarder()
		logForwarderCWConfigBuilder := cmv1.NewLogForwarderCloudWatchConfig()

		logForwarderCWConfigBuilder.LogGroupName(logForwarderYaml.CloudWatch.CloudWatchLogGroupName)

		logForwarderCWConfigBuilder.LogDistributionRoleArn(logForwarderYaml.CloudWatch.CloudWatchLogRoleArn)

		if !reflect.DeepEqual(currentLogForwarder.Applications(), logForwarderYaml.CloudWatch.Applications) {
			logForwarderConfigBuilder.Applications(logForwarderYaml.CloudWatch.Applications...)
		}

		if !reflect.DeepEqual(currentLogForwarder.Groups(), logForwarderYaml.CloudWatch.GroupsLogVersions) {
			var groupBuilder []*cmv1.LogForwarderGroupBuilder
			for _, group := range logForwarderYaml.CloudWatch.GroupsLogVersions {
				groupBuilder = append(groupBuilder, cmv1.NewLogForwarderGroup().ID(group))
			}
			logForwarderConfigBuilder.Groups(groupBuilder...)
		}

		logForwarderConfigBuilder.Cloudwatch(logForwarderCWConfigBuilder)
		logForwarderConfig = logForwarderConfigBuilder
	}

	if logForwarderConfig == nil {
		return errors.UserErrorf("log forwarding config changes contained no valid log forwarder changes")
	}

	body, err := logForwarderConfig.Build()
	if err != nil {
		return err
	}

	return c.UpdateLogForwarder(body, logForwarderId, clusterId)
}
