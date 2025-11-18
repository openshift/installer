/*
Copyright 2024 The ORC Authors.

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

package image

import (
	"context"
	"slices"
	"strings"

	"github.com/gophercloud/gophercloud/v2/openstack/image/v2/imageimport"
	"github.com/gophercloud/gophercloud/v2/openstack/image/v2/images"
	ctrl "sigs.k8s.io/controller-runtime"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/internal/logging"
	orcerrors "github.com/k-orc/openstack-resource-controller/v2/internal/util/errors"
)

func requireResource(orcImage *orcv1alpha1.Image) (*orcv1alpha1.ImageResourceSpec, error) {
	resource := orcImage.Spec.Resource
	if resource == nil {
		return nil, orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "resource not provided")
	}

	return resource, nil
}

func requireResourceContent(orcImage *orcv1alpha1.Image) (*orcv1alpha1.ImageContent, error) {
	resource, err := requireResource(orcImage)
	if err != nil {
		return nil, err
	}
	if resource.Content == nil {
		return nil, orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "resource content not provided")
	}
	return resource.Content, nil
}

func (actuator imageActuator) canWebDownload(ctx context.Context, orcImage *orcv1alpha1.Image) (bool, error) {
	log := ctrl.LoggerFrom(ctx)

	debugLog := func(reason string, extra ...any) {
		log.V(logging.Verbose).Info("Image cannot use web-download", slices.Concat([]any{"reason", reason}, extra)...)
	}

	content, err := requireResourceContent(orcImage)
	if err != nil {
		return false, err
	}

	download := content.Download
	if download == nil {
		debugLog("not type URL")
		return false, nil
	}

	// web-download does not deterministically support decompression.
	// Glance can be configured to do automatic decompression of imported
	// images, but there's no way to determine if it is enabled so it can't
	// be safely used.
	if download.Decompress != nil {
		debugLog("web-download does not support decompression")
		return false, nil
	}

	// web-download can't be used with hash verification
	// Glance does publish a hash for an image whose data was imported with web-download, but:
	// * the hash it publishes is of the image contents after
	//   backend-specific processing, e.g. conversion to raw. There is no way
	//   to determine via the API what processing will be done on the image,
	//   and without downloading the image and performing that processing
	//   ourselves, no way to know what the resulting hash should be.
	// * even if the backend doesn't perform any processing, which we can't
	//   determine, if the hash we have doesn't match the algorithm configured
	//   server-side in glance, which is also not exposed via the API, we
	//   can't verify it anyway.
	if download.Hash != nil {
		debugLog("web-download does not support hash verification")
		return false, nil
	}

	// Get supported import methods from Glance
	importInfo, err := actuator.osClient.GetImportInfo(ctx)
	if err != nil {
		return false, err
	}

	if !slices.Contains(importInfo.ImportMethods.Value, string(imageimport.WebDownloadMethod)) {
		debugLog("glance is not configured with web-download", "import-methods", strings.Join(importInfo.ImportMethods.Value, ", "))
		return false, nil
	}

	return true, nil
}

func (actuator imageActuator) webDownload(ctx context.Context, orcImage *orcv1alpha1.Image, glanceImage *images.Image) error {
	log := ctrl.LoggerFrom(ctx)
	log.V(logging.Verbose).Info("Importing with web-download")

	resource := orcImage.Spec.Resource
	if resource == nil {
		// Should have been caught by validation
		return orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "resource not provided")
	}

	content, err := requireResourceContent(orcImage)
	if err != nil {
		return err
	}

	return actuator.osClient.CreateImport(ctx, glanceImage.ID, &imageimport.CreateOpts{
		Name: imageimport.WebDownloadMethod,
		URI:  content.Download.URL,
	})
}
