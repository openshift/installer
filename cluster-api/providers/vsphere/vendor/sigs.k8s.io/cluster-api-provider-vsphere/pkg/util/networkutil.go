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

package util

import (
	"context"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apitypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// NCPSNATKey is the key used for the NCPSNAT annotation.
	NCPSNATKey = "ncp/snat_ip"
	// NCPVersionKey is the key used for version information in the NCP configmap.
	NCPVersionKey = "version"
	// NCPNamespace is the namespace of the NCP configmap.
	NCPNamespace = "vmware-system-nsx"
	// NCPVersionConfigMap is a name of the NCP config map.
	NCPVersionConfigMap = "nsx-ncp-version-config"
	// NCPVersionSupportFW 3.0.1 is where NCP starts to support "whitelist_source_ranges" specification in VNET and enforce FW rules on GC T1.
	NCPVersionSupportFW = "3.0.1"
	// NCPVersionSupportFWEnded 3.1.0 is where NCP stopped to support "whitelist_source_ranges" specification in VNET.
	NCPVersionSupportFWEnded = "3.1.0"

	// EmptyAnnotationErrorMsg is an error message returned when no annotations are found.
	EmptyAnnotationErrorMsg = "annotation not found"
	// EmptyNCPSNATKeyMsg is an error message returned when the annotation can not be found.
	EmptyNCPSNATKeyMsg = NCPSNATKey + " key not found"
)

var (
	// NCPVersionSupportFWSemver is the SemVer representation of the minimum NCPVersion for enforcing FW rules.
	NCPVersionSupportFWSemver = semver.MustParse(NCPVersionSupportFW)
	// NCPVersionSupportFWEndedSemver is the SemVer representation of the maximum NCPVersion for enforcing FW rules.
	NCPVersionSupportFWEndedSemver = semver.MustParse(NCPVersionSupportFWEnded)
)

// GetNamespaceNetSnatIP finds out the namespace's corresponding network's SNAT IP.
func GetNamespaceNetSnatIP(ctx context.Context, controllerClient client.Client, namespace string) (string, error) {
	namespaceObj := &corev1.Namespace{}
	namespacedName := apitypes.NamespacedName{
		Name: namespace,
	}

	if err := controllerClient.Get(ctx, namespacedName, namespaceObj); err != nil {
		return "", err
	}

	annotations := namespaceObj.GetAnnotations()
	if annotations == nil {
		return "", errors.New(EmptyAnnotationErrorMsg)
	}

	snatIP := annotations[NCPSNATKey]
	if snatIP == "" {
		return "", errors.New(EmptyNCPSNATKeyMsg)
	}

	return snatIP, nil
}

// GetNCPVersion finds out the running ncp's version from its configmap.
// If the version contains more than 3 segments, it will get trimmed down to 3.
func GetNCPVersion(ctx context.Context, controllerClient client.Client) (string, error) {
	configmapObj := &corev1.ConfigMap{}
	namespacedName := apitypes.NamespacedName{
		Name:      NCPVersionConfigMap,
		Namespace: NCPNamespace,
	}

	if err := controllerClient.Get(ctx, namespacedName, configmapObj); err != nil {
		return "", err
	}

	version := configmapObj.Data[NCPVersionKey]

	// NSX doesn't stritcly follow SemVer and there are versions like 4.0.1.3
	// This will cause an error if directly used in semver.Parse() and prevent the cluster from reconciling.
	// Since GetNCPVersion is only used to check >= 3.0.1 and < 3.1.0, it's safe to trim the last segment
	if segments := strings.Split(version, "."); len(segments) > 3 {
		return strings.Join(segments[:3], "."), nil
	}
	return version, nil
}

// NCPSupportFW checks the version of running NCP and return true if it supports FW rule enforcement on GC T1 Router.
func NCPSupportFW(ctx context.Context, controllerClient client.Client) (bool, error) {
	ncpVersion, err := GetNCPVersion(ctx, controllerClient)
	if err != nil {
		return false, err
	}
	currVersion, err := semver.Parse(ncpVersion)
	if err != nil {
		return false, err
	}
	supported := currVersion.GTE(NCPVersionSupportFWSemver) && currVersion.LT(NCPVersionSupportFWEndedSemver)
	return supported, nil
}
