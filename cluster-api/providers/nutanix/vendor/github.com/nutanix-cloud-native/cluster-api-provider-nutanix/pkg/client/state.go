/*
Copyright 2022 Nutanix

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

package client

import (
	"context"
	"fmt"
	"math"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/nutanix-cloud-native/prism-go-client/utils"
	nutanixClientV3 "github.com/nutanix-cloud-native/prism-go-client/v3"
)

type stateRefreshFunc func() (string, error)

func WaitForTaskCompletion(ctx context.Context, conn *nutanixClientV3.Client, uuid string) error {
	errCh := make(chan error, 1)
	go waitForState(
		errCh,
		"SUCCEEDED",
		waitUntilTaskStateFunc(ctx, conn, uuid))

	err := <-errCh
	return err
}

func waitForState(errCh chan<- error, target string, refresh stateRefreshFunc) {
	err := Retry(2, 2, 0, func(_ uint) (bool, error) {
		state, err := refresh()
		if err != nil {
			return false, err
		} else if state == target {
			return true, nil
		}
		return false, nil
	})
	errCh <- err
}

func waitUntilTaskStateFunc(ctx context.Context, conn *nutanixClientV3.Client, uuid string) stateRefreshFunc {
	return func() (string, error) {
		return GetTaskState(ctx, conn, uuid)
	}
}

func GetTaskState(ctx context.Context, client *nutanixClientV3.Client, taskUUID string) (string, error) {
	log := ctrl.LoggerFrom(ctx)
	log.V(1).Info(fmt.Sprintf("Getting task with UUID %s", taskUUID))
	v, err := client.V3.GetTask(ctx, taskUUID)
	if err != nil {
		log.Error(err, fmt.Sprintf("error occurred while waiting for task with UUID %s", taskUUID))
		return "", err
	}

	if *v.Status == "INVALID_UUID" || *v.Status == "FAILED" {
		return *v.Status,
			fmt.Errorf("error_detail: %s, progress_message: %s", utils.StringValue(v.ErrorDetail), utils.StringValue(v.ProgressMessage))
	}
	taskStatus := *v.Status
	log.V(1).Info(fmt.Sprintf("Status for task with UUID %s: %s", taskUUID, taskStatus))
	return taskStatus, nil
}

// RetryableFunc performs an action and returns a bool indicating whether the
// function is done, or if it should keep retrying, and an error which will
// abort the retry and be returned by the Retry function. The 0-indexed attempt
// is passed with each call.
type RetryableFunc func(uint) (bool, error)

/*
Retry retries a function up to numTries times with exponential backoff.
If numTries == 0, retry indefinitely.
If interval == 0, Retry will not delay retrying and there will be no
exponential backoff.
If maxInterval == 0, maxInterval is set to +Infinity.
Intervals are in seconds.
Returns an error if initial > max intervals, if retries are exhausted, or if the passed function returns
an error.
*/
func Retry(initialInterval float64, maxInterval float64, numTries uint, function RetryableFunc) error {
	if maxInterval == 0 {
		maxInterval = math.Inf(1)
	} else if initialInterval < 0 || initialInterval > maxInterval {
		return fmt.Errorf("invalid retry intervals (negative or initial < max). Initial: %f, Max: %f", initialInterval, maxInterval)
	}

	var err error
	done := false
	interval := initialInterval
	for i := uint(0); !done && (numTries == 0 || i < numTries); i++ {
		done, err = function(i)
		if err != nil {
			return err
		}

		if !done {
			// Retry after delay. Calculate next delay.
			time.Sleep(time.Duration(interval) * time.Second)
			interval = math.Min(interval*2, maxInterval)
		}
	}

	if !done {
		return fmt.Errorf("function never succeeded in Retry")
	}
	return nil
}
