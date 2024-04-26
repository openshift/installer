/*
Copyright 2021 The Kubernetes Authors.

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

package networking

import (
	"context"
	"fmt"
	"time"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/trunks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"

	"sigs.k8s.io/cluster-api-provider-openstack/pkg/record"
	capoerrors "sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/errors"
)

const (
	timeoutTrunkDelete       = 3 * time.Minute
	retryIntervalTrunkDelete = 5 * time.Second
)

func (s *Service) GetTrunkSupport() (bool, error) {
	allExts, err := s.client.ListExtensions()
	if err != nil {
		return false, err
	}

	for _, ext := range allExts {
		if ext.Alias == "trunk" {
			return true, nil
		}
	}
	return false, nil
}

func (s *Service) getOrCreateTrunkForPort(eventObject runtime.Object, port *ports.Port) (*trunks.Trunk, error) {
	trunkList, err := s.client.ListTrunk(trunks.ListOpts{
		Name:   port.Name,
		PortID: port.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("searching for existing trunk for server: %v", err)
	}

	if len(trunkList) != 0 {
		return &trunkList[0], nil
	}

	trunkCreateOpts := trunks.CreateOpts{
		Name:        port.Name,
		PortID:      port.ID,
		Description: port.Description,
	}

	trunk, err := s.client.CreateTrunk(trunkCreateOpts)
	if err != nil {
		return nil, err
	}

	record.Eventf(eventObject, "SuccessfulCreateTrunk", "Created trunk %s with id %s", trunk.Name, trunk.ID)
	return trunk, nil
}

func (s *Service) DeleteTrunk(eventObject runtime.Object, portID string) error {
	listOpts := trunks.ListOpts{
		PortID: portID,
	}
	trunkInfo, err := s.client.ListTrunk(listOpts)
	if err != nil {
		return err
	}
	if len(trunkInfo) != 1 {
		return nil
	}

	err = wait.PollUntilContextTimeout(context.TODO(), retryIntervalTrunkDelete, timeoutTrunkDelete, true, func(_ context.Context) (bool, error) {
		if err := s.client.DeleteTrunk(trunkInfo[0].ID); err != nil {
			if capoerrors.IsNotFound(err) {
				record.Eventf(eventObject, "SuccessfulDeleteTrunk", "Trunk %s with id %s did not exist", trunkInfo[0].Name, trunkInfo[0].ID)
				return true, nil
			}
			if capoerrors.IsConflict(err) {
				return false, nil
			}
			if capoerrors.IsRetryable(err) {
				return false, nil
			}
			return false, err
		}
		return true, nil
	})
	if err != nil {
		record.Warnf(eventObject, "FailedDeleteTrunk", "Failed to delete trunk %s with id %s: %v", trunkInfo[0].Name, trunkInfo[0].ID, err)
		return err
	}

	record.Eventf(eventObject, "SuccessfulDeleteTrunk", "Deleted trunk %s with id %s", trunkInfo[0].Name, trunkInfo[0].ID)
	return nil
}
