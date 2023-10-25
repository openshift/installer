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
	"fmt"
	"sort"
	"strings"

	"github.com/blang/semver"
	"github.com/pkg/errors"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

// Service provides operations on Azure VM Images.
type Service struct {
	Client
	azure.Authorizer
}

// New creates a VM Images service.
func New(auth azure.Authorizer) (*Service, error) {
	client, err := NewClient(auth)
	if err != nil {
		return nil, err
	}
	return &Service{
		Client:     client,
		Authorizer: auth,
	}, nil
}

// GetDefaultUbuntuImage returns the default image spec for Ubuntu.
func (s *Service) GetDefaultUbuntuImage(ctx context.Context, location, k8sVersion string) (*infrav1.Image, error) {
	v, err := semver.ParseTolerant(k8sVersion)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to parse Kubernetes version \"%s\"", k8sVersion)
	}

	osVersion := getUbuntuOSVersion(v.Major, v.Minor, v.Patch)
	publisher, offer := azure.DefaultImagePublisherID, azure.DefaultImageOfferID
	skuID, version, err := s.getSKUAndVersion(
		ctx, location, publisher, offer, k8sVersion, fmt.Sprintf("ubuntu-%s", osVersion))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get default image")
	}

	defaultImage := &infrav1.Image{
		Marketplace: &infrav1.AzureMarketplaceImage{
			ImagePlan: infrav1.ImagePlan{
				Publisher: publisher,
				Offer:     offer,
				SKU:       skuID,
			},
			Version: version,
		},
	}

	return defaultImage, nil
}

// GetDefaultWindowsImage returns the default image spec for Windows.
func (s *Service) GetDefaultWindowsImage(ctx context.Context, location, k8sVersion, runtime, osAndVersion string) (*infrav1.Image, error) {
	v122 := semver.MustParse("1.22.0")
	v, err := semver.ParseTolerant(k8sVersion)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to parse Kubernetes version \"%s\"", k8sVersion)
	}

	// If containerd is specified we don't currently support less than 1.22
	if v.LE(v122) && runtime == "containerd" {
		return nil, errors.New("containerd image only supported in 1.22+")
	}

	if osAndVersion == "" {
		osAndVersion = azure.DefaultWindowsOsAndVersion
	}

	// Starting with 1.22 we default to containerd for Windows unless the runtime flag is set.
	if v.GE(v122) && runtime != "dockershim" && !strings.HasSuffix(osAndVersion, "-containerd") {
		osAndVersion += "-containerd"
	}

	publisher, offer := azure.DefaultImagePublisherID, azure.DefaultWindowsImageOfferID
	skuID, version, err := s.getSKUAndVersion(
		ctx, location, publisher, offer, k8sVersion, osAndVersion)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get default image")
	}

	defaultImage := &infrav1.Image{
		Marketplace: &infrav1.AzureMarketplaceImage{
			ImagePlan: infrav1.ImagePlan{
				Publisher: publisher,
				Offer:     offer,
				SKU:       skuID,
			},
			Version: version,
		},
	}

	return defaultImage, nil
}

// getSKUAndVersion gets the SKU ID and version of the image to use for the provided version of Kubernetes.
// note: osAndVersion is expected to be in the format of {os}-{version} (ex: ubuntu-2004 or windows-2022)
func (s *Service) getSKUAndVersion(ctx context.Context, location, publisher, offer, k8sVersion, osAndVersion string) (skuID string, imageVersion string, err error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "virtualmachineimages.Service.getSKUAndVersion")
	defer done()

	log.V(4).Info("Getting VM image SKU and version", "location", location, "publisher", publisher, "offer", offer, "k8sVersion", k8sVersion, "osAndVersion", osAndVersion)

	v, err := semver.ParseTolerant(k8sVersion)
	if err != nil {
		return "", "", errors.Wrapf(err, "unable to parse Kubernetes version \"%s\" in spec, expected valid SemVer string", k8sVersion)
	}

	// Old SKUs before 1.21.12, 1.22.9, or 1.23.6 are named like "k8s-1dot21dot2-ubuntu-2004".
	if k8sVersionInSKUName(v.Major, v.Minor, v.Patch) {
		return fmt.Sprintf("k8s-%ddot%ddot%d-%s", v.Major, v.Minor, v.Patch, osAndVersion), azure.LatestVersion, nil
	}

	// New SKUs don't contain the Kubernetes version and are named like "ubuntu-2004-gen1".
	sku := fmt.Sprintf("%s-gen1", osAndVersion)

	imageCache, err := GetCache(s.Authorizer)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to get image cache")
	}
	imageCache.client = s.Client

	listImagesResponse, err := imageCache.Get(ctx, location, publisher, offer, sku)
	if err != nil {
		return "", "", errors.Wrapf(err, "unable to list VM images for publisher \"%s\" offer \"%s\" sku \"%s\"", publisher, offer, sku)
	}

	vmImages := listImagesResponse.VirtualMachineImageResourceArray
	if len(vmImages) == 0 {
		return "", "", errors.Errorf("no VM images found for publisher \"%s\" offer \"%s\" sku \"%s\"", publisher, offer, sku)
	}

	// Sort the VM image names descending, so more recent dates sort first.
	// (The date is encoded into the end of the name, for example "124.0.20220512").
	names := []string{}
	for _, vmImage := range vmImages {
		names = append(names, *vmImage.Name)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(names)))

	// Pick the first (most recent) one whose k8s version matches.
	var version string
	id := fmt.Sprintf("%d%d.%d", v.Major, v.Minor, v.Patch)
	for _, name := range names {
		if strings.HasPrefix(name, id) {
			version = name
			break
		}
	}
	if version == "" {
		return "", "", errors.Errorf("no VM image found for publisher \"%s\" offer \"%s\" sku \"%s\" with Kubernetes version \"%s\"", publisher, offer, sku, k8sVersion)
	}

	log.V(4).Info("Found VM image SKU and version", "location", location, "publisher", publisher, "offer", offer, "sku", sku, "version", version)

	return sku, version, nil
}

// getUbuntuOSVersion returns the default Ubuntu OS version for the given Kubernetes version.
func getUbuntuOSVersion(major, minor, patch uint64) string {
	// Default to Ubuntu 22.04 LTS for Kubernetes v1.25.3 and later.
	osVersion := "2204"
	if major == 1 && minor == 21 && patch < 2 ||
		major == 1 && minor == 20 && patch < 8 ||
		major == 1 && minor == 19 && patch < 12 ||
		major == 1 && minor == 18 && patch < 20 ||
		major == 1 && minor < 18 {
		osVersion = "1804"
	} else if major == 1 && minor == 25 && patch < 3 ||
		major == 1 && minor < 25 {
		osVersion = "2004"
	}
	return osVersion
}

// k8sVersionInSKUName returns true if the k8s version is in the SKU name (the older style of naming).
func k8sVersionInSKUName(major, minor, patch uint64) bool {
	return (major == 1 && minor < 21) ||
		(major == 1 && minor == 21 && patch <= 12) ||
		(major == 1 && minor == 22 && patch <= 9) ||
		(major == 1 && minor == 23 && patch <= 6)
}
