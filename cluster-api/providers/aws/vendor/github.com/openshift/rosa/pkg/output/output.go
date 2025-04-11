/*
Copyright (c) 2021 Red Hat, Inc.

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

// This file contains functions used to implement the '--output' command line option.

package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"sigs.k8s.io/yaml"

	arv1 "github.com/openshift-online/ocm-sdk-go/accesstransparency/v1"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	msv1 "github.com/openshift-online/ocm-sdk-go/servicemgmt/v1"

	"gitlab.com/c0b/go-ordered-json"

	"github.com/openshift/rosa/pkg/aws"
)

// When ocm-sdk-go encounters an empty resource list, it marshals it as a
// string that represents an empty JSON array with newline and spaces in between:
// '[', '\n', ' ', ' ', '\n', ']'. This byte-array allows us to compare that so
// that the output can be shown correctly.
var emptyBuffer = []byte{91, 10, 32, 32, 10, 93}

func Print(resource interface{}) error {
	var b bytes.Buffer

	switch reflect.TypeOf(resource).String() {
	case "[]*v1.ManagedService":
		if managedServices, ok := resource.([]*msv1.ManagedService); ok {
			msv1.MarshalManagedServiceList(managedServices, &b)
		}
	case "[]*v1.CloudRegion":
		if cloudRegions, ok := resource.([]*cmv1.CloudRegion); ok {
			cmv1.MarshalCloudRegionList(cloudRegions, &b)
		}
	case "*v1.Cluster":
		if cluster, ok := resource.(*cmv1.Cluster); ok {
			cmv1.MarshalCluster(cluster, &b)
		}
	case "[]*v1.Cluster":
		if clusters, ok := resource.([]*cmv1.Cluster); ok {
			cmv1.MarshalClusterList(clusters, &b)
		}
	case "[]*v1.DNSDomain":
		if dnsdomains, ok := resource.([]*cmv1.DNSDomain); ok {
			cmv1.MarshalDNSDomainList(dnsdomains, &b)
		}
	case "[]*v1.ExternalAuth":
		if externalAuths, ok := resource.([]*cmv1.ExternalAuth); ok {
			cmv1.MarshalExternalAuthList(externalAuths, &b)
		}
	case "*v1.ExternalAuth":
		if externalAuth, ok := resource.(*cmv1.ExternalAuth); ok {
			cmv1.MarshalExternalAuth(externalAuth, &b)
		}
	case "[]*v1.IdentityProvider":
		if idps, ok := resource.([]*cmv1.IdentityProvider); ok {
			cmv1.MarshalIdentityProviderList(idps, &b)
		}
	case "*v1.Ingress":
		if ingress, ok := resource.(*cmv1.Ingress); ok {
			cmv1.MarshalIngress(ingress, &b)
		}
	case "[]*v1.Ingress":
		if ingresses, ok := resource.([]*cmv1.Ingress); ok {
			cmv1.MarshalIngressList(ingresses, &b)
		}
	case "[]*v1.MachinePool":
		if machinePools, ok := resource.([]*cmv1.MachinePool); ok {
			cmv1.MarshalMachinePoolList(machinePools, &b)
		}
	case "*v1.MachinePool":
		if machinePool, ok := resource.(*cmv1.MachinePool); ok {
			cmv1.MarshalMachinePool(machinePool, &b)
		}
	case "[]*v1.MachineType":
		if machineTypes, ok := resource.([]*cmv1.MachineType); ok {
			cmv1.MarshalMachineTypeList(machineTypes, &b)
		}
	case "*v1.NodePool":
		if nodePool, ok := resource.(*cmv1.NodePool); ok {
			cmv1.MarshalNodePool(nodePool, &b)
		}
	case "[]*v1.NodePool":
		if nodePools, ok := resource.([]*cmv1.NodePool); ok {
			cmv1.MarshalNodePoolList(nodePools, &b)
		}
	case "[]*v1.Version":
		if versions, ok := resource.([]*cmv1.Version); ok {
			cmv1.MarshalVersionList(versions, &b)
		}
	case "[]*v1.VersionGate":
		if versionGates, ok := resource.([]*cmv1.VersionGate); ok {
			cmv1.MarshalVersionGateList(versionGates, &b)
		}
	case "[]*v1.OidcConfig":
		if oidcConfigs, ok := resource.([]*cmv1.OidcConfig); ok {
			cmv1.MarshalOidcConfigList(oidcConfigs, &b)
		}
	case "*v1.OidcConfig":
		if oidcConfig, ok := resource.(*cmv1.OidcConfig); ok {
			cmv1.MarshalOidcConfig(oidcConfig, &b)
		}
	case "[]*v1.BreakGlassCredential":
		if breakGlassCredentials, ok := resource.([]*cmv1.BreakGlassCredential); ok {
			cmv1.MarshalBreakGlassCredentialList(breakGlassCredentials, &b)
		}
	case "*v1.BreakGlassCredential":
		if breakGlassCredential, ok := resource.(*cmv1.BreakGlassCredential); ok {
			cmv1.MarshalBreakGlassCredential(breakGlassCredential, &b)
		}
	case "[]*v1.TuningConfig":
		if tuningConfigs, ok := resource.([]*cmv1.TuningConfig); ok {
			cmv1.MarshalTuningConfigList(tuningConfigs, &b)
		}
	case "*v1.TuningConfig":
		if tuningConfig, ok := resource.(*cmv1.TuningConfig); ok {
			cmv1.MarshalTuningConfig(tuningConfig, &b)
		}
	case "*v1.KubeletConfig":
		if kubeletConfig, ok := resource.(*cmv1.KubeletConfig); ok {
			cmv1.MarshalKubeletConfig(kubeletConfig, &b)
		}
	case "[]*v1.KubeletConfig":
		if kubeletConfigs, ok := resource.([]*cmv1.KubeletConfig); ok {
			cmv1.MarshalKubeletConfigList(kubeletConfigs, &b)
		}
	case "*v1.ClusterAutoscaler":
		if autoscaler, ok := resource.(*cmv1.ClusterAutoscaler); ok {
			cmv1.MarshalClusterAutoscaler(autoscaler, &b)
		}
	case "[]*v1.User":
		if users, ok := resource.([]*cmv1.User); ok {
			cmv1.MarshalUserList(users, &b)
		}
	case "*v1.SubnetNetworkVerification":
		if subnetNetworkVerification, ok := resource.(*cmv1.SubnetNetworkVerification); ok {
			cmv1.MarshalSubnetNetworkVerification(subnetNetworkVerification, &b)
		}
	case "[]aws.Role", "[]aws.OidcProviderOutput":
		{
			err := defaultEncode(resource, &b)
			if err != nil {
				return err
			}
		}
	case "map[string][]aws.Role":
		{
			for _, operatorRoles := range resource.(map[string][]aws.Role) {
				err := Print(operatorRoles)
				if err != nil {
					return err
				}
			}
		}
	case "*v1.AccessRequest":
		if accessRequest, ok := resource.(*arv1.AccessRequest); ok {
			arv1.MarshalAccessRequest(accessRequest, &b)
		}
	// default to catch non concrete types
	default:
		{
			err := defaultEncode(resource, &b)
			if err != nil {
				return err
			}
		}
	}
	// Verify if the resource is an empty string and ensure that the JSON
	// representation looks correct for STDOUT.
	if b.String() == string(emptyBuffer) {
		b = *bytes.NewBufferString("[]")
	}
	str, err := parseResource(b)
	if err != nil {
		return err
	}
	fmt.Print(str)
	return nil
}

// Provides a default encoding to JSON for types not being marshalled via the cmv1 package
func defaultEncode(resource interface{}, b *bytes.Buffer) error {
	reqBodyBytes := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBytes).Encode(resource)
	if err != nil {
		return err
	}

	err = json.Indent(b, reqBodyBytes.Bytes(), "", "  ")
	if err != nil {
		return err
	}

	return nil
}

func parseResource(body bytes.Buffer) (string, error) {
	switch o {
	case "json":
		var out bytes.Buffer
		prettifyJSON(&out, body.Bytes())
		return out.String(), nil
	case "yaml":
		out, err := yaml.JSONToYAML(body.Bytes())
		if err != nil {
			return "", err
		}
		return string(out), nil
	default:
		return "", fmt.Errorf("Unknown format '%s'. Valid formats are %s", o, formats)
	}
}

func prettifyJSON(stream io.Writer, body []byte) error {
	if len(body) == 0 {
		return nil
	}
	data := ordered.NewOrderedMap()
	err := json.Unmarshal(body, data)
	if err != nil {
		return dumpBytes(stream, body)
	}
	return dumpJSON(stream, data)
}

func dumpBytes(stream io.Writer, data []byte) error {
	_, err := stream.Write(data)
	if err != nil {
		return err
	}
	_, err = stream.Write([]byte("\n"))
	return err
}

func dumpJSON(stream io.Writer, data *ordered.OrderedMap) error {
	encoder := json.NewEncoder(stream)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}
