/*
Copyright 2020 The Kubernetes Authors.

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

package userdata

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/alessio/shellescape"
)

const (
	nodeUserData = `#!/bin/bash
/etc/eks/bootstrap.sh {{.ClusterName}} {{- template "args" . }}
`
)

// NodeInput defines the context to generate a node user data.
type NodeInput struct {
	ClusterName           string
	KubeletExtraArgs      map[string]string
	ContainerRuntime      *string
	DNSClusterIP          *string
	DockerConfigJSON      *string
	APIRetryAttempts      *int
	PauseContainerAccount *string
	PauseContainerVersion *string
	UseMaxPods            *bool
	// NOTE: currently the IPFamily/ServiceIPV6Cidr isn't exposed to the user.
	// TODO (richardcase): remove the above comment when IPV6 / dual stack is implemented.
	IPFamily        *string
	ServiceIPV6Cidr *string
}

func (ni *NodeInput) DockerConfigJSONEscaped() string {
	if ni.DockerConfigJSON == nil || len(*ni.DockerConfigJSON) == 0 {
		return "''"
	}

	return shellescape.Quote(*ni.DockerConfigJSON)
}

// NewNode returns the user data string to be used on a node instance.
func NewNode(input *NodeInput) ([]byte, error) {
	tm := template.New("Node")

	if _, err := tm.Parse(argsTemplate); err != nil {
		return nil, fmt.Errorf("failed to parse args template: %w", err)
	}

	if _, err := tm.Parse(kubeletArgsTemplate); err != nil {
		return nil, fmt.Errorf("failed to parse kubeletExtraArgs template: %w", err)
	}

	t, err := tm.Parse(nodeUserData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Node template: %w", err)
	}

	var out bytes.Buffer
	if err := t.Execute(&out, input); err != nil {
		return nil, fmt.Errorf("failed to generate Node template: %w", err)
	}

	return out.Bytes(), nil
}
