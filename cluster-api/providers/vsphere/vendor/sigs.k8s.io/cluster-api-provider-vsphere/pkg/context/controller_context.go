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

package context

import (
	"fmt"

	"github.com/go-logr/logr"

	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/record"
)

// ControllerContext is the context of a controller.
type ControllerContext struct {
	*ControllerManagerContext

	// Name is the name of the controller.
	Name string

	// Logger is the controller's logger.
	Logger logr.Logger

	// Recorder is used to record events.
	Recorder record.Recorder
}

// String returns ControllerManagerName/ControllerName.
func (c *ControllerContext) String() string {
	return fmt.Sprintf("%s/%s", c.ControllerManagerContext.String(), c.Name)
}
