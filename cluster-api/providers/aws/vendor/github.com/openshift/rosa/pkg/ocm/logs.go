/*
Copyright (c) 2020 Red Hat, Inc.

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
	"context"
	"fmt"
	"net/http"
	"time"

	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	errors "github.com/zgalor/weberr"
)

const interval = 15 * time.Second

func (c *Client) GetInstallLogs(clusterID string, tail int) (logs *cmv1.Log, err error) {
	logsClient := c.ocm.ClustersMgmt().V1().Clusters().
		Cluster(clusterID).
		Logs().
		Install()
	response, err := logsClient.Get().
		Parameter("tail", tail).
		Send()
	if err != nil {
		err = handleErr(response.Error(), err)
		if response.Status() == http.StatusNotFound {
			err = errors.NotFound.UserErrorf("Failed to get logs for cluster '%s'", clusterID)
		}
		return
	}

	return response.Body(), nil
}

func (c *Client) GetUninstallLogs(clusterID string, tail int) (logs *cmv1.Log, err error) {
	logsClient := c.ocm.ClustersMgmt().V1().Clusters().
		Cluster(clusterID).
		Logs().
		Uninstall()
	response, err := logsClient.Get().
		Parameter("tail", tail).
		Send()
	if err != nil {
		err = handleErr(response.Error(), err)
		if response.Status() == http.StatusNotFound {
			err = errors.NotFound.UserErrorf("Failed to get logs for cluster '%s'", clusterID)
		}
		return
	}

	return response.Body(), nil
}

func (c *Client) PollInstallLogs(clusterID string, cb func(*cmv1.LogGetResponse) bool) (logs *cmv1.Log, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer func() {
		cancel()
	}()

	logsClient := c.ocm.ClustersMgmt().V1().Clusters().
		Cluster(clusterID).
		Logs().
		Install()
	response, err := logsClient.Poll().
		Parameter("tail", 100).
		Interval(interval).
		Predicate(cb).
		StartContext(ctx)
	if err != nil {
		err = fmt.Errorf("Failed to poll logs for cluster '%s': %v", clusterID, err)
		if response.Status() == http.StatusNotFound {
			err = errors.NotFound.UserErrorf("Failed to poll logs for cluster '%s'", clusterID)
		}
		return
	}

	return response.Body(), nil
}

func (c *Client) PollUninstallLogs(clusterID string,
	cb func(*cmv1.LogGetResponse) bool) (logs *cmv1.Log, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer func() {
		cancel()
	}()

	logsClient := c.ocm.ClustersMgmt().V1().Clusters().
		Cluster(clusterID).
		Logs().
		Uninstall()
	response, err := logsClient.Poll().
		Parameter("tail", 100).
		Interval(interval).
		Predicate(cb).
		StartContext(ctx)
	if err != nil {
		err = fmt.Errorf("Failed to poll logs for cluster '%s': %v", clusterID, err)
		if response.Status() == http.StatusNotFound {
			err = errors.NotFound.UserErrorf("Failed to poll logs for cluster '%s'", clusterID)
		}
		return
	}

	return response.Body(), nil
}
