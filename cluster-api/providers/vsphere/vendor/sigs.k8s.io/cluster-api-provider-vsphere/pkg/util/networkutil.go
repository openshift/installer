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

	"github.com/hashicorp/go-version"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	apitypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	NCPSNATKey          = "ncp/snat_ip"
	NCPVersionKey       = "version"
	NCPNamespace        = "vmware-system-nsx"
	NCPVersionConfigMap = "nsx-ncp-version-config"
	// 3.0.1 is where NCP starts to support "whitelist_source_ranges" specification in VNET and enforce FW rules on GC T1.
	NCPVersionSupportFW = "3.0.1"
	// 3.1.0 is where NCP stopped to support "whitelist_source_ranges" specification in VNET.
	NCPVersionSupportFWEnded = "3.1.0"

	EmptyAnnotationErrorMsg = "annotation not found"
	EmptyNCPSNATKeyMsg      = NCPSNATKey + " key not found"
)

// GetNamespaceNetSnatIP finds out the namespace's corresponding network's SNAT IP.
func GetNamespaceNetSnatIP(ctx context.Context, controllerClient client.Client, namespace string) (string, error) {
	namespaceObj := &v1.Namespace{}
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
func GetNCPVersion(ctx context.Context, controllerClient client.Client) (string, error) {
	configmapObj := &v1.ConfigMap{}
	namespacedName := apitypes.NamespacedName{
		Name:      NCPVersionConfigMap,
		Namespace: NCPNamespace,
	}

	if err := controllerClient.Get(ctx, namespacedName, configmapObj); err != nil {
		return "", err
	}

	version := configmapObj.Data[NCPVersionKey]
	return version, nil
}

// NCPSupportFW checks the version of running NCP and return true if it supports FW rule enforcement on GC T1 Router.
func NCPSupportFW(ctx context.Context, controllerClient client.Client) (bool, error) {
	ncpVersion, err := GetNCPVersion(ctx, controllerClient)
	if err != nil {
		return false, err
	}
	currVersion, err := version.NewVersion(ncpVersion)
	if err != nil {
		return false, err
	}
	supportStartedVersion, err := version.NewVersion(NCPVersionSupportFW)
	if err != nil {
		return false, err
	}
	supportEndedVersion, err := version.NewVersion(NCPVersionSupportFWEnded)
	if err != nil {
		return false, err
	}
	supported := currVersion.GreaterThanOrEqual(supportStartedVersion) && currVersion.LessThan(supportEndedVersion)
	return supported, nil
}
