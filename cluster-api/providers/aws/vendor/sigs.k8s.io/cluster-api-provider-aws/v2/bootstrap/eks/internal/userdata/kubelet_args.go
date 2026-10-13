/*
Copyright 2022 The Kubernetes Authors.

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

const argsTemplate = `{{- define "args" -}}
{{- if .KubeletExtraArgs }} --kubelet-extra-args '{{ template "kubeletArgsTemplate" .KubeletExtraArgs }}'
{{- end -}}
{{- if .ContainerRuntime }} --container-runtime {{.ContainerRuntime}}{{- end -}}
{{- if .IPFamily }} --ip-family {{.IPFamily}}{{- end -}}
{{- if .ServiceIPV6Cidr }} --service-ipv6-cidr {{.ServiceIPV6Cidr}}{{- end -}}
{{- if .UseMaxPods }} --use-max-pods {{.UseMaxPods}}{{- end -}}
{{- if .APIRetryAttempts }} --aws-api-retry-attempts {{.APIRetryAttempts}}{{- end -}}
{{- if .PauseContainerAccount }} --pause-container-account {{.PauseContainerAccount}}{{- end -}}
{{- if .PauseContainerVersion }} --pause-container-version {{.PauseContainerVersion}}{{- end -}}
{{- if .DNSClusterIP }} --dns-cluster-ip {{.DNSClusterIP}}{{- end -}}
{{- if .DockerConfigJSON }} --docker-config-json {{ .DockerConfigJSONEscaped }}{{- end -}}
{{- end -}}`

const kubeletArgsTemplate = `{{- define "kubeletArgsTemplate" -}}
{{- $first := true -}}
{{- range $k, $v := . -}}
{{- if $first -}}{{ $first = false -}}{{- else }} {{ end -}}
--{{$k}}={{$v}}
{{- end -}}
{{- end -}}
`
