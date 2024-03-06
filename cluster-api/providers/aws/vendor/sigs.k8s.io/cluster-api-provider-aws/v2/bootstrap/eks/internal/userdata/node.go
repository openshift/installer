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

	eksbootstrapv1 "sigs.k8s.io/cluster-api-provider-aws/v2/bootstrap/eks/api/v1beta2"
)

const (
	defaultBootstrapCommand = "/etc/eks/bootstrap.sh"

	nodeUserData = `#cloud-config
{{template "files" .Files}}
runcmd:
{{- template "commands" .PreBootstrapCommands }}
  - {{ .BootstrapCommand }} {{.ClusterName}} {{- template "args" . }}
{{- template "commands" .PostBootstrapCommands }}
{{- template "ntp" .NTP }}
{{- template "users" .Users }}
{{- template "disk_setup" .DiskSetup}}
{{- template "fs_setup" .DiskSetup}}
{{- template "mounts" .Mounts}}
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
	IPFamily                 *string
	ServiceIPV6Cidr          *string
	PreBootstrapCommands     []string
	PostBootstrapCommands    []string
	BootstrapCommandOverride *string
	Files                    []eksbootstrapv1.File
	DiskSetup                *eksbootstrapv1.DiskSetup
	Mounts                   []eksbootstrapv1.MountPoints
	Users                    []eksbootstrapv1.User
	NTP                      *eksbootstrapv1.NTP
}

// DockerConfigJSONEscaped returns the DockerConfigJSON escaped for use in cloud-init.
func (ni *NodeInput) DockerConfigJSONEscaped() string {
	if ni.DockerConfigJSON == nil || len(*ni.DockerConfigJSON) == 0 {
		return "''"
	}

	return shellescape.Quote(*ni.DockerConfigJSON)
}

// BootstrapCommand returns the bootstrap command to be used on a node instance.
func (ni *NodeInput) BootstrapCommand() string {
	if ni.BootstrapCommandOverride != nil && *ni.BootstrapCommandOverride != "" {
		return *ni.BootstrapCommandOverride
	}

	return defaultBootstrapCommand
}

// NewNode returns the user data string to be used on a node instance.
func NewNode(input *NodeInput) ([]byte, error) {
	tm := template.New("Node").Funcs(defaultTemplateFuncMap)

	if _, err := tm.Parse(filesTemplate); err != nil {
		return nil, fmt.Errorf("failed to parse args template: %w", err)
	}

	if _, err := tm.Parse(argsTemplate); err != nil {
		return nil, fmt.Errorf("failed to parse args template: %w", err)
	}

	if _, err := tm.Parse(kubeletArgsTemplate); err != nil {
		return nil, fmt.Errorf("failed to parse kubeletExtraArgs template: %w", err)
	}

	if _, err := tm.Parse(commandsTemplate); err != nil {
		return nil, fmt.Errorf("failed to parse commandsTemplate template: %w", err)
	}

	if _, err := tm.Parse(ntpTemplate); err != nil {
		return nil, fmt.Errorf("failed to parse ntp template: %w", err)
	}

	if _, err := tm.Parse(usersTemplate); err != nil {
		return nil, fmt.Errorf("failed to parse users template: %w", err)
	}

	if _, err := tm.Parse(diskSetupTemplate); err != nil {
		return nil, fmt.Errorf("failed to parse disk setup template: %w", err)
	}

	if _, err := tm.Parse(fsSetupTemplate); err != nil {
		return nil, fmt.Errorf("failed to parse fs setup template: %w", err)
	}

	if _, err := tm.Parse(mountsTemplate); err != nil {
		return nil, fmt.Errorf("failed to parse mounts template: %w", err)
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
