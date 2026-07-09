/*
Copyright 2018 The Kubernetes Authors.

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

package compute

import (
	"context"
	"errors"
	"fmt"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/v2/openstack/image/v2/images"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta2"
	capoerrors "sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/errors"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/filterconvert"
)

// Helper function for getting image ID from name, ID, or tags.
func (s *Service) GetImageID(ctx context.Context, k8sClient client.Client, namespace string, image infrav1.ImageParam) (*string, error) {
	switch {
	case image.ID != nil:
		return image.ID, nil
	case image.Filter != nil:
		return s.getImageIDByFilter(image.Filter)
	case image.ImageRef != nil:
		return s.getImageIDByReference(ctx, k8sClient, namespace, image.ImageRef)
	default:
		// Should have been caught by validation
		return nil, errors.New("image id, filter, and reference are all nil")
	}
}

func (s *Service) getImageIDByFilter(filter *infrav1.ImageFilter) (*string, error) {
	listOpts := filterconvert.ImageFilterToListOpts(filter)
	allImages, err := s.getImageClient().ListImages(listOpts)
	if err != nil {
		return nil, err
	}

	switch len(allImages) {
	case 0:
		var name string
		if filter.Name != nil {
			name = *filter.Name
		}
		return nil, fmt.Errorf("no images were found with the given image filter: name=%v, tags=%v", name, filter.Tags)
	case 1:
		return &allImages[0].ID, nil
	default:
		// this should never happen
		var name string
		if filter.Name != nil {
			name = *filter.Name
		}
		return nil, fmt.Errorf("too many images were found with the given image filter: name=%v, tags=%v", name, filter.Tags)
	}
}

func (s *Service) getImageIDByReference(ctx context.Context, k8sClient client.Client, namespace string, ref *infrav1.ResourceReference) (*string, error) {
	orcImage := &orcv1alpha1.Image{}
	err := k8sClient.Get(ctx, client.ObjectKey{
		Namespace: namespace,
		Name:      ref.Name,
	}, orcImage)
	if err != nil {
		// Not an error if it doesn't exist yet
		if apierrors.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	if orcv1alpha1.IsAvailable(orcImage) {
		return orcImage.Status.ID, nil
	}

	if !orcv1alpha1.IsReconciliationComplete(orcImage) {
		return nil, nil
	}

	err = orcv1alpha1.GetTerminalError(orcImage)
	if err != nil {
		return nil, capoerrors.Terminal(infrav1.DependencyFailedReason, orcImage.Kind+" "+orcImage.GetNamespace()+"/"+orcImage.GetName()+" failed: "+err.Error())
	}

	return nil, nil
}

func (s *Service) GetFlavorID(flavorParam infrav1.FlavorParam) (string, error) {
	if flavorParam.ID != nil {
		return *flavorParam.ID, nil
	}

	if flavorParam.Filter == nil || flavorParam.Filter.Name == nil {
		return "", fmt.Errorf("no flavors were found: no name set")
	}

	allFlavors, err := s.getComputeClient().ListFlavors()
	if err != nil {
		return "", err
	}

	for _, flavor := range allFlavors {
		if flavor.Name == *flavorParam.Filter.Name {
			return flavor.ID, nil
		}
	}

	return "", fmt.Errorf("no flavors were found: name=%v", *flavorParam.Filter.Name)
}

func (s *Service) GetFlavor(flavorID string) (*flavors.Flavor, error) {
	return s.getComputeClient().GetFlavor(flavorID)
}

func (s *Service) GetImageDetails(imageID string) (*images.Image, error) {
	return s.getImageClient().GetImage(imageID)
}

// GetManagementPort returns the port which is used for management and external
// traffic. Cluster floating IPs must be associated with this port.
func (s *Service) GetManagementPort(openStackCluster *infrav1.OpenStackCluster, instanceStatus *InstanceStatus) (*ports.Port, error) {
	ns, err := instanceStatus.NetworkStatus()
	if err != nil {
		return nil, err
	}

	networkingService, err := s.getNetworkingService()
	if err != nil {
		return nil, err
	}

	allPorts, err := networkingService.GetPortFromInstanceIP(instanceStatus.ID(), ns.IP(openStackCluster.Status.Network.Name))
	if err != nil {
		return nil, fmt.Errorf("lookup management port for server %s: %w", instanceStatus.ID(), err)
	}
	if len(allPorts) < 1 {
		return nil, fmt.Errorf("did not find management port for server %s", instanceStatus.ID())
	}
	return &allPorts[0], nil
}

func (s *Service) GetInstanceStatus(resourceID string) (instance *InstanceStatus, err error) {
	if resourceID == "" {
		return nil, fmt.Errorf("resourceId should be specified to get detail")
	}

	server, err := s.getComputeClient().GetServer(resourceID)
	if err != nil {
		if capoerrors.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("get server %q detail failed: %v", resourceID, err)
	}

	return &InstanceStatus{server, s.scope.Logger()}, nil
}
