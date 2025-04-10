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

package virtualmachineimages

import (
	"context"
	"regexp"
	"strings"

	"github.com/blang/semver"
	"github.com/pkg/errors"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

/* Reference images live in an Azure community gallery with this structure:
. Gallery: "ClusterAPI-f72ceb4f-5159-4c26-a0fe-2ea738f0d019"
├── Image Definition: "capi-ubun2-2404"
│   ├── Version: "1.30.4"
│   ├── Version: "1.30.5"
│   └── Version: "1.31.1"
└── Image Definition: "capi-win-2022-containerd"
    ├── Version: "1.30.4"
    ├── Version: "1.30.5"
    └── Version: "1.31.1"
*/

// Service provides operations on Azure VM Images.
type Service struct {
}

// New creates a VM Images service.
func New(_ azure.Authorizer) (*Service, error) {
	return &Service{}, nil
}

// GetDefaultLinuxImage returns the default image spec for Ubuntu.
func (s *Service) GetDefaultLinuxImage(ctx context.Context, _, k8sVersion string) (*infrav1.Image, error) {
	_, _, done := tele.StartSpanWithLogger(ctx, "azure.services.virtualmachineimages.GetDefaultLinuxImage")
	defer done()

	v, err := semver.ParseTolerant(k8sVersion)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to parse Kubernetes version \"%s\"", k8sVersion)
	}

	// Use the Azure Marketplace for specific older versions, to keep "clusterctl upgrade" from rolling new machines.
	marketplaceVersionRange := semver.MustParseRange("<=1.28.12 || >=1.29.1 <=1.29.7 || >=1.30.0 <=1.30.3")
	if marketplaceVersionRange(v) {
		if version, ok := oldUbuntu2204Versions[v.String()]; ok {
			return &infrav1.Image{
				Marketplace: &infrav1.AzureMarketplaceImage{
					ImagePlan: infrav1.ImagePlan{
						Publisher: "cncf-upstream",
						Offer:     "capi",
						SKU:       "ubuntu-2204-gen1",
					},
					Version: version,
				},
			}, nil
		}
	}
	return &infrav1.Image{
		ComputeGallery: &infrav1.AzureComputeGalleryImage{
			Gallery: azure.DefaultPublicGalleryName,
			Name:    azure.DefaultLinuxGalleryImageName,
			Version: v.String(),
		},
	}, nil
}

// GetDefaultWindowsImage returns the default image spec for Windows.
func (s *Service) GetDefaultWindowsImage(ctx context.Context, _, k8sVersion, runtime, osAndVersion string) (*infrav1.Image, error) {
	_, _, done := tele.StartSpanWithLogger(ctx, "azure.services.virtualmachineimages.GetDefaultWindowsImage")
	defer done()

	v, err := semver.ParseTolerant(k8sVersion)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to parse Kubernetes version \"%s\"", k8sVersion)
	}

	if runtime != "" && runtime != "containerd" {
		return nil, errors.Errorf("unsupported runtime %s", runtime)
	}

	// Use the Windows edition from the osAndVersion if provided.
	imageName := azure.DefaultWindowsGalleryImageName
	if osAndVersion != "" {
		match := regexp.MustCompile(`^windows-(\d{4})$`).FindStringSubmatch(osAndVersion)
		if len(match) != 2 {
			return nil, errors.Errorf("unsupported osAndVersion %s", osAndVersion)
		}
		imageName = strings.Replace(imageName, "2019", match[1], 1)
	}

	// Use the Azure Marketplace for specific older versions, to keep "clusterctl upgrade" from rolling new machines.
	marketplaceVersionRange := semver.MustParseRange("<=1.28.12 || >=1.29.1 <=1.29.7 || >=1.30.0 <=1.30.3")
	if marketplaceVersionRange(v) {
		if version, ok := oldWindows2019Versions[v.String()]; ok {
			return &infrav1.Image{
				Marketplace: &infrav1.AzureMarketplaceImage{
					ImagePlan: infrav1.ImagePlan{
						Publisher: "cncf-upstream",
						Offer:     "capi-windows",
						SKU:       "windows-2019-containerd-gen1",
					},
					Version: version,
				},
			}, nil
		}
	}
	return &infrav1.Image{
		ComputeGallery: &infrav1.AzureComputeGalleryImage{
			Gallery: azure.DefaultPublicGalleryName,
			Name:    imageName,
			Version: v.String(),
		},
	}, nil
}

// oldUbuntu2204Versions maps Kubernetes versions to Azure Marketplace image versions.
// The Marketplace offer is deprecated and won't be updated, so these values are
// hard-coded here to simplify lookup.
var oldUbuntu2204Versions = map[string]string{
	"1.27.14": "127.14.20240517",
	"1.27.15": "127.15.20240612",
	"1.27.16": "127.16.20240717",
	"1.28.1":  "128.1.20230829",
	"1.28.2":  "128.2.20230918",
	"1.28.3":  "128.3.20231023",
	"1.28.4":  "128.4.20231130",
	"1.28.5":  "128.5.20240102",
	"1.28.6":  "128.6.20240201",
	"1.28.7":  "128.7.20240223",
	"1.28.8":  "128.8.20240327",
	"1.28.9":  "128.9.20240418",
	"1.28.10": "128.10.20240517",
	"1.28.11": "128.11.20240612",
	"1.28.12": "128.12.20240717",
	"1.29.1":  "129.1.20240206",
	"1.29.2":  "129.2.20240223",
	"1.29.3":  "129.3.20240327",
	"1.29.4":  "129.4.20240418",
	"1.29.5":  "129.5.20240517",
	"1.29.6":  "129.6.20240612",
	"1.29.7":  "129.7.20240717",
	"1.30.0":  "130.0.20240506",
	"1.30.1":  "130.1.20240517",
	"1.30.2":  "130.2.20240612",
	"1.30.3":  "130.3.20240717",
}

// oldWindows2019Versions maps Kubernetes versions to Azure Marketplace image versions.
// The Marketplace offer is deprecated and won't be updated, so these values are
// hard-coded here to simplify lookup.
var oldWindows2019Versions = map[string]string{
	"1.27.14": "127.14.20240515",
	"1.28.1":  "128.1.20230830",
	"1.28.2":  "128.2.20230918",
	"1.28.3":  "128.3.20231023",
	"1.28.4":  "128.4.20231122",
	"1.28.5":  "128.5.20240102",
	"1.28.6":  "128.6.20240201",
	"1.28.7":  "128.7.20240223",
	"1.28.8":  "128.8.20240327",
	"1.28.9":  "128.9.20240418",
	"1.28.10": "128.10.20240515",
	"1.28.11": "128.11.20240612",
	"1.28.12": "128.12.20240717",
	"1.29.1":  "129.1.20240201",
	"1.29.2":  "129.2.20240223",
	"1.29.3":  "129.3.20240327",
	"1.29.4":  "129.4.20240418",
	"1.29.5":  "129.5.20240515",
	"1.29.6":  "129.6.20240612",
	"1.29.7":  "129.7.20240717",
	"1.30.0":  "130.0.20240418",
	"1.30.1":  "130.1.20240515",
	"1.30.2":  "130.2.20240612",
	"1.30.3":  "130.3.20240717",
}
