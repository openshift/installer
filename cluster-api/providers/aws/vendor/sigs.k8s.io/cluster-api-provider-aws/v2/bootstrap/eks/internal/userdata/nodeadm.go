/*
Copyright 2026 The Kubernetes Authors.

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

	"k8s.io/apimachinery/pkg/runtime"

	eksbootstrapv1 "sigs.k8s.io/cluster-api-provider-aws/v2/bootstrap/eks/api/v1beta2"
)

const (
	// boundary is the MIME multipart boundary used to separate sections in cloud-init userdata.
	boundary = "//"

	// this does not start with a boundary because it is the last item that is processed.
	cloudInitUserData = `
Content-Type: text/cloud-config
MIME-Version: 1.0
Content-Transfer-Encoding: 7bit
Content-Disposition: attachment; filename="cloud-config.yaml"

#cloud-config
{{- if .Files }}
{{template "files" .Files}}
{{- end }}
{{- if .NTP }}
{{- template "ntp" .NTP }}
{{- end }}
{{- if .Users }}
{{- template "users" .Users }}
{{- end }}
{{- if .DiskSetup }}
{{- template "disk_setup" .DiskSetup }}
{{- template "fs_setup" .DiskSetup }}
{{- end }}
{{- if .Mounts }}
{{- template "mounts" .Mounts }}
{{- end }}
--{{.Boundary}}`

	// Shell script part template for nodeadm.
	shellScriptPartTemplate = `
--{{.Boundary}}
Content-Type: text/x-shellscript; charset="us-ascii"
MIME-Version: 1.0
Content-Transfer-Encoding: 7bit
Content-Disposition: attachment; filename="commands.sh"

#!/bin/bash
set -o errexit
set -o pipefail
set -o nounset
{{- range .PreNodeadmCommands}}
{{.}}
{{- end}}
--{{ .Boundary }}`

	// Node config part template for nodeadm.
	nodeConfigPartTemplate = `
--{{.Boundary}}
Content-Type: application/node.eks.aws

---
apiVersion: node.eks.aws/v1alpha1
kind: NodeConfig
spec:
  cluster:
    name: {{.ClusterName}}
    apiServerEndpoint: {{.APIServerEndpoint}}
    certificateAuthority: {{.CACert}}
    cidr: {{if .ServiceCIDR}}{{.ServiceCIDR}}{{else}}172.20.0.0/16{{end}}
  {{- if .FeatureGates }}
  featureGates:
  {{- range $k, $v := .FeatureGates }}
    {{$k}}: {{$v}}
  {{- end }}
  {{- end }}
{{- if or .KubeletConfig .KubeletFlags }}
  kubelet:
    {{- if .KubeletConfig }}
    config:
{{ Indent 6 (toYaml .KubeletConfig) }}
    {{- end }}
    {{- if .KubeletFlags }}
    flags:
      {{- range $flag := .KubeletFlags }}
       - "{{$flag}}"
      {{- end }}
    {{- end }}
{{- end }}
  {{- if or .ContainerdConfig .ContainerdBaseRuntimeSpec }}
  containerd:
    {{- if .ContainerdConfig }}
    config:
{{ Indent 6 .ContainerdConfig }}
    {{- end }}
    {{- if .ContainerdBaseRuntimeSpec }}
    baseRuntimeSpec:
{{ Indent 6 (toYaml .ContainerdBaseRuntimeSpec) }}
    {{- end }}
  {{- end }}

--{{.Boundary}}`
)

// NodeadmInput contains all the information required to generate user data for a node.
type NodeadmInput struct {
	ClusterName               string
	KubeletFlags              []string
	KubeletConfig             *runtime.RawExtension
	ContainerdConfig          string
	ContainerdBaseRuntimeSpec *runtime.RawExtension
	FeatureGates              map[eksbootstrapv1.Feature]bool

	PreNodeadmCommands []string
	Files              []eksbootstrapv1.File
	DiskSetup          *eksbootstrapv1.DiskSetup
	Mounts             []eksbootstrapv1.MountPoints
	Users              []eksbootstrapv1.User
	NTP                *eksbootstrapv1.NTP

	AMIImageID        string
	APIServerEndpoint string
	Boundary          string
	CACert            string
	ServiceCIDR       string // Service CIDR range for the cluster
	ClusterDNS        string
}

// validateNodeInput validates the input for nodeadm user data generation.
func validateNodeadmInput(input *NodeadmInput) error {
	if input.APIServerEndpoint == "" {
		return fmt.Errorf("API server endpoint is required for nodeadm")
	}
	if input.CACert == "" {
		return fmt.Errorf("CA certificate is required for nodeadm")
	}
	if input.ClusterName == "" {
		return fmt.Errorf("cluster name is required for nodeadm")
	}
	if input.Boundary == "" {
		input.Boundary = boundary
	}

	return nil
}

// NewNodeadmUserdata returns the user data string to be used on a node instance.
func NewNodeadmUserdata(input *NodeadmInput) ([]byte, error) {
	if err := validateNodeadmInput(input); err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	// Write MIME header
	if _, err := buf.WriteString(fmt.Sprintf("MIME-Version: 1.0\nContent-Type: multipart/mixed; boundary=%q\n\n", input.Boundary)); err != nil {
		return nil, fmt.Errorf("failed to write MIME header: %v", err)
	}

	// Write shell script part if needed
	if len(input.PreNodeadmCommands) > 0 {
		shellScriptTemplate := template.Must(template.New("shell").Parse(shellScriptPartTemplate))
		if err := shellScriptTemplate.Execute(&buf, input); err != nil {
			return nil, fmt.Errorf("failed to execute shell script template: %v", err)
		}
		if _, err := buf.WriteString("\n"); err != nil {
			return nil, fmt.Errorf("failed to write newline: %v", err)
		}
	}

	// Write node config part
	nodeConfigTemplate := template.Must(
		template.New("node").
			Funcs(defaultTemplateFuncMap).
			Parse(nodeConfigPartTemplate),
	)
	if err := nodeConfigTemplate.Execute(&buf, input); err != nil {
		return nil, fmt.Errorf("failed to execute node config template: %v", err)
	}

	// Write cloud-config part
	tm := template.New("Node").Funcs(defaultTemplateFuncMap)
	// if any of the input fields are set, we need to write the cloud-config part.
	// Each sub-template is parsed individually so that failures produce clear error messages
	// indicating which specific section (files, ntp, users, etc.) failed to parse.
	if input.NTP != nil || input.DiskSetup != nil || input.Mounts != nil || input.Users != nil || input.Files != nil {
		if _, err := tm.Parse(filesTemplate); err != nil {
			return nil, fmt.Errorf("failed to parse args template: %w", err)
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

		t, err := tm.Parse(cloudInitUserData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Node template: %w", err)
		}

		if err := t.Execute(&buf, input); err != nil {
			return nil, fmt.Errorf("failed to execute node user data template: %w", err)
		}
	}
	// write the final boundary closing, all of the ones in the script use intermediate boundaries
	buf.Write([]byte("--"))
	return buf.Bytes(), nil
}
