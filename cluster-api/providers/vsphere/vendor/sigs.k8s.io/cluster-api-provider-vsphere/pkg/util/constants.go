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

// Package util contains utils.
package util

const metadataFormat = `
instance-id: "{{ .Hostname }}"
local-hostname: "{{ .Hostname }}"
wait-on-network:
  ipv4: {{ .WaitForIPv4 }}
  ipv6: {{ .WaitForIPv6 }}
network:
  version: 2
  ethernets:
    {{- range $i, $net := .Devices }}
    id{{ $i }}:
      match:
        macaddress: "{{ $net.MACAddr }}"
      {{- if $net.DeviceName }}
      set-name: "{{ $net.DeviceName }}"
      {{- else }}
      set-name: "eth{{ $i }}"
      {{- end }}
      wakeonlan: true
      {{- if or $net.DHCP4 $net.DHCP6 }}
      dhcp4: {{ $net.DHCP4 }}
	  {{- if $net.DHCP4Overrides }}
      dhcp4-overrides:
	    {{- if $net.DHCP4Overrides.Hostname }}
        hostname: "{{ $net.DHCP4Overrides.Hostname }}"
	    {{- end }}
	    {{- if $net.DHCP4Overrides.RouteMetric }}
        route-metric: {{ $net.DHCP4Overrides.RouteMetric }}
	    {{- end }}
	    {{- if $net.DHCP4Overrides.SendHostname }}
        send-hostname: {{ $net.DHCP4Overrides.SendHostname }}
	    {{- end }}
	    {{- if $net.DHCP4Overrides.UseDNS }}
        use-dns: {{ $net.DHCP4Overrides.UseDNS }}
	    {{- end }}
	    {{- if $net.DHCP4Overrides.UseDomains }}
        use-domains: {{ $net.DHCP4Overrides.UseDomains }}
	    {{- end }}
	    {{- if $net.DHCP4Overrides.UseHostname }}
        use-hostname: {{ $net.DHCP4Overrides.UseHostname }}
	    {{- end }}
	    {{- if $net.DHCP4Overrides.UseMTU }}
        use-mtu: {{ $net.DHCP4Overrides.UseMTU }}
	    {{- end }}
	    {{- if $net.DHCP4Overrides.UseNTP }}
        use-ntp: {{ $net.DHCP4Overrides.UseNTP }}
	    {{- end }}
	    {{- if $net.DHCP4Overrides.UseRoutes }}
        use-routes: "{{ $net.DHCP4Overrides.UseRoutes }}"
	    {{- end }}
	  {{- end }}
      dhcp6: {{ $net.DHCP6 }}
	  {{- if $net.DHCP6Overrides }}
      dhcp6-overrides:
	    {{- if $net.DHCP6Overrides.Hostname }}
        hostname: "{{ $net.DHCP6Overrides.Hostname }}"
	    {{- end }}
	    {{- if $net.DHCP6Overrides.RouteMetric }}
        route-metric: {{ $net.DHCP6Overrides.RouteMetric }}
	    {{- end }}
	    {{- if $net.DHCP6Overrides.SendHostname }}
        send-hostname: {{ $net.DHCP6Overrides.SendHostname }}
	    {{- end }}
	    {{- if $net.DHCP6Overrides.UseDNS }}
        use-dns: {{ $net.DHCP6Overrides.UseDNS }}
	    {{- end }}
	    {{- if $net.DHCP6Overrides.UseDomains }}
        use-domains: {{ $net.DHCP6Overrides.UseDomains }}
	    {{- end }}
	    {{- if $net.DHCP6Overrides.UseHostname }}
        use-hostname: {{ $net.DHCP6Overrides.UseHostname }}
	    {{- end }}
	    {{- if $net.DHCP6Overrides.UseMTU }}
        use-mtu: {{ $net.DHCP6Overrides.UseMTU }}
	    {{- end }}
	    {{- if $net.DHCP6Overrides.UseNTP }}
        use-ntp: {{ $net.DHCP6Overrides.UseNTP }}
	    {{- end }}
	    {{- if $net.DHCP6Overrides.UseRoutes }}
        use-routes: "{{ $net.DHCP6Overrides.UseRoutes }}"
	    {{- end }}
	  {{- end }}
      {{- end }}
      {{- if $net.IPAddrs }}
      addresses:
      {{- range $net.IPAddrs }}
      - "{{ . }}"
      {{- end }}
      {{- end }}
      {{- if $net.Gateway4 }}
      gateway4: "{{ $net.Gateway4 }}"
      {{- end }}
      {{- if $net.Gateway6 }}
      gateway6: "{{ $net.Gateway6 }}"
      {{- end }}
      {{- if .MTU }}
      mtu: {{ .MTU }}
      {{- end }}
      {{- if .Routes }}
      routes:
      {{- range .Routes }}
      - to: "{{ .To }}"
        via: "{{ .Via }}"
        metric: {{ .Metric }}
      {{- end }}
      {{- end }}
      {{- if nameservers $net }}
      nameservers:
        {{- if $net.Nameservers }}
        addresses:
        {{- range $net.Nameservers }}
        - "{{ . }}"
        {{- end }}
        {{- end }}
        {{- if $net.SearchDomains }}
        search:
        {{- range $net.SearchDomains }}
        - "{{ . }}"
        {{- end }}
        {{- end }}
      {{- end }}
    {{- end }}
  {{- if .Routes }}
  routes:
  {{- range .Routes }}
  - to: "{{ .To }}"
    via: "{{ .Via }}"
    metric: {{ .Metric }}
  {{- end }}
  {{- end }}
`
