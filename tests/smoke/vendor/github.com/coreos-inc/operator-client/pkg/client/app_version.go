package client

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/coreos-inc/operator-client/pkg/types"

	"github.com/golang/glog"

	"k8s.io/client-go/pkg/api/errors"
	"k8s.io/client-go/pkg/runtime"
	"k8s.io/client-go/pkg/util/wait"
)

// GetAppVersion will return the AppVersion resource for the given name.
func (c *Client) GetAppVersion(name string) (*types.AppVersion, error) {
	glog.V(4).Infof("[GET AppVersion]: %s", name)
	obj, err := c.GetThirdPartyResource(
		types.TectonicAPIGroup,
		types.AppVersionTPRVersion,
		types.TectonicNamespace,
		types.AppVersionTPRKind,
		name,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting AppVersion %s: %v", name, err)
	}
	return appVersionFromUnstructured(obj)
}

// UpdateAppVersion will update the given AppVersion resource.
func (c *Client) UpdateAppVersion(av *types.AppVersion) (*types.AppVersion, error) {
	glog.V(4).Infof("[UPDATE AppVersion]: %s", av.GetName())
	obj, err := unstructuredFromAppVersion(av)
	if err != nil {
		return nil, fmt.Errorf("error getting AppVersion %s: %v", av.GetName(), err)
	}
	err = c.UpdateThirdPartyResource(obj)
	if err != nil {
		return nil, fmt.Errorf("error updating AppVersion %s: %v", av.GetName(), err)
	}
	return c.GetAppVersion(obj.GetName())
}

// AtomicUpdateAppVersion takes an update function which is executed before attempting
// to update the AppVersion resource. Upon conflict, the update function is run
// again, until the update is successful or a non-conflict error is returned.
func (c *Client) AtomicUpdateAppVersion(name string, f types.AppVersionModifier) (*types.AppVersion, error) {
	glog.V(4).Infof("[ATOMIC UPDATE AppVersion]: %s", name)
	var nav *types.AppVersion
	wait.ExponentialBackoff(wait.Backoff{
		Duration: time.Second,
		Factor:   2.0,
		Jitter:   0.1,
		Steps:    5,
	}, func() (bool, error) {
		av, err := c.GetAppVersion(name)
		if err != nil {
			return false, fmt.Errorf("error getting current AppVersion resource: %v", err)
		}
		if err = f(av); err != nil {
			return false, err
		}
		nav, err = c.UpdateAppVersion(av)
		if err != nil {
			if errors.IsConflict(err) {
				glog.Info("conflict updating AppVersion resource, will try again")
				return false, nil
			}
			return false, err
		}
		return true, nil
	})
	return nav, nil
}

func unstructuredFromAppVersion(av *types.AppVersion) (*runtime.Unstructured, error) {
	avb, err := json.Marshal(av)
	if err != nil {
		return nil, fmt.Errorf("error marshaling AppVersion resource: %v", err)
	}
	var r runtime.Unstructured
	if err := json.Unmarshal(avb, &r.Object); err != nil {
		return nil, fmt.Errorf("error unmarshaling marshaled resource: %v", err)
	}
	return &r, nil
}

func appVersionFromUnstructured(r *runtime.Unstructured) (*types.AppVersion, error) {
	avb, err := json.Marshal(r.Object)
	if err != nil {
		return nil, fmt.Errorf("error marshaling unstructured resource: %v", err)
	}
	var av types.AppVersion
	if err := json.Unmarshal(avb, &av); err != nil {
		return nil, fmt.Errorf("error unmarshmaling marshaled resource to TectonicAppVersion: %v", err)
	}
	return &av, nil
}

// SetFailureStatus sets the failure status in the AppVersion.Status.
// If nil is passed, then the failure status is cleared.
func (c *Client) SetFailureStatus(name string, failureStatus *types.FailureStatus) error {
	_, err := c.AtomicUpdateAppVersion(name, func(av *types.AppVersion) error {
		av.Status.FailureStatus = failureStatus
		return nil
	})
	return err
}

// SetTaskStatuses sets the task status list in the AppVersion.Status.
// If nil is passed, then the task status list is cleared.
func (c *Client) SetTaskStatuses(name string, ts []types.TaskStatus) error {
	_, err := c.AtomicUpdateAppVersion(name, func(av *types.AppVersion) error {
		av.Status.TaskStatuses = ts
		return nil
	})
	return err
}

// UpdateTaskStatus updates the task status in the AppVersion.Status.Taskstatues list.
// It will return the error if the name of the task is not found in the list.
func (c *Client) UpdateTaskStatus(name string, ts types.TaskStatus) error {
	_, err := c.AtomicUpdateAppVersion(name, func(av *types.AppVersion) error {
		var found bool
		for i, v := range av.Status.TaskStatuses {
			if v.Name == ts.Name {
				av.Status.TaskStatuses[i] = ts
				found = true
			}
		}
		if !found {
			return fmt.Errorf("%q is not found in TaskStatus", ts.Name)
		}
		return nil
	})
	return err
}
