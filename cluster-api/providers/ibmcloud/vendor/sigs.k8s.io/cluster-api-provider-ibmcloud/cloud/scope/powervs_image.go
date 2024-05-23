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

package scope

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-logr/logr"

	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"

	"k8s.io/klog/v2"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"sigs.k8s.io/cluster-api/util/patch"

	infrav1beta2 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/powervs"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/resourcecontroller"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/endpoints"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/record"
)

// BucketAccess indicates if the bucket has public or private access public access.
const BucketAccess = "public"

// PowerVSImageScopeParams defines the input parameters used to create a new PowerVSImageScope.
type PowerVSImageScopeParams struct {
	Client          client.Client
	Logger          logr.Logger
	IBMPowerVSImage *infrav1beta2.IBMPowerVSImage
	ServiceEndpoint []endpoints.ServiceEndpoint
}

// PowerVSImageScope defines a scope defined around a Power VS Cluster.
type PowerVSImageScope struct {
	logr.Logger
	Client      client.Client
	patchHelper *patch.Helper

	IBMPowerVSClient powervs.PowerVS
	IBMPowerVSImage  *infrav1beta2.IBMPowerVSImage
	ServiceEndpoint  []endpoints.ServiceEndpoint
}

// NewPowerVSImageScope creates a new PowerVSImageScope from the supplied parameters.
func NewPowerVSImageScope(params PowerVSImageScopeParams) (scope *PowerVSImageScope, err error) {
	scope = &PowerVSImageScope{}

	if params.Client == nil {
		err = errors.New("failed to generate new scope from nil Client")
		return nil, err
	}
	scope.Client = params.Client

	if params.IBMPowerVSImage == nil {
		err = errors.New("failed to generate new scope from nil IBMPowerVSImage")
		return nil, err
	}
	scope.IBMPowerVSImage = params.IBMPowerVSImage

	if params.Logger == (logr.Logger{}) {
		params.Logger = klog.Background()
	}
	scope.Logger = params.Logger

	helper, err := patch.NewHelper(params.IBMPowerVSImage, params.Client)
	if err != nil {
		err = fmt.Errorf("failed to init patch helper: %w", err)
		return nil, err
	}
	scope.patchHelper = helper

	// Create Resource Controller client.
	var serviceOption resourcecontroller.ServiceOptions
	// Fetch the resource controller endpoint.
	rcEndpoint := endpoints.FetchEndpoints(string(endpoints.RC), params.ServiceEndpoint)
	if rcEndpoint != "" {
		serviceOption.URL = rcEndpoint
		params.Logger.V(3).Info("Overriding the default resource controller endpoint", "ResourceControllerEndpoint", rcEndpoint)
	}

	rc, err := resourcecontroller.NewService(serviceOption)
	if err != nil {
		return nil, err
	}

	var serviceInstanceID string
	spec := params.IBMPowerVSImage.Spec
	if spec.ServiceInstanceID != "" {
		serviceInstanceID = spec.ServiceInstanceID
	} else if params.IBMPowerVSImage.Spec.ServiceInstance != nil && params.IBMPowerVSImage.Spec.ServiceInstance.ID != nil {
		serviceInstanceID = *params.IBMPowerVSImage.Spec.ServiceInstance.ID
	} else {
		name := fmt.Sprintf("%s-%s", params.IBMPowerVSImage.Spec.ClusterName, "serviceInstance")
		if params.IBMPowerVSImage.Spec.ServiceInstance != nil && params.IBMPowerVSImage.Spec.ServiceInstance.Name != nil {
			name = *params.IBMPowerVSImage.Spec.ServiceInstance.Name
		}
		serviceInstance, err := rc.GetServiceInstance("", name)
		if err != nil {
			params.Logger.Error(err, "error failed to get service instance id from name", "name", name)
			return nil, err
		}
		if serviceInstance == nil {
			return nil, fmt.Errorf("service instance %s is not yet created", name)
		}
		if *serviceInstance.State != string(infrav1beta2.ServiceInstanceStateActive) {
			return nil, fmt.Errorf("service instance %s is not in active state", name)
		}
		serviceInstanceID = *serviceInstance.GUID
	}

	res, _, err := rc.GetResourceInstance(
		&resourcecontrollerv2.GetResourceInstanceOptions{
			ID: &serviceInstanceID,
		})
	if err != nil {
		err = fmt.Errorf("failed to get resource instance: %w", err)
		return nil, err
	}

	options := powervs.ServiceOptions{
		IBMPIOptions: &ibmpisession.IBMPIOptions{
			Debug: params.Logger.V(DEBUGLEVEL).Enabled(),
			Zone:  *res.RegionID,
		},
	}

	// Fetch the service endpoint.
	if svcEndpoint := endpoints.FetchPVSEndpoint(endpoints.ConstructRegionFromZone(*res.RegionID), params.ServiceEndpoint); svcEndpoint != "" {
		options.IBMPIOptions.URL = svcEndpoint
		scope.Logger.V(3).Info("overriding the default powervs service endpoint")
	}

	c, err := powervs.NewService(options)
	if err != nil {
		err = fmt.Errorf("failed to create NewIBMPowerVSClient error %w", err)
		return nil, err
	}

	options.CloudInstanceID = serviceInstanceID
	c.WithClients(options)
	scope.IBMPowerVSClient = c
	return scope, nil
}

func (i *PowerVSImageScope) ensureImageUnique(imageName string) (*models.ImageReference, error) {
	images, err := i.IBMPowerVSClient.GetAllImage()
	if err != nil {
		return nil, err
	}
	for _, img := range images.Images {
		if *img.Name == imageName {
			return img, nil
		}
	}
	return nil, nil
}

// CreateImageCOSBucket creates a power vs image.
func (i *PowerVSImageScope) CreateImageCOSBucket() (*models.ImageReference, *models.JobReference, error) {
	s := i.IBMPowerVSImage.Spec
	m := i.IBMPowerVSImage.ObjectMeta

	imageReply, err := i.ensureImageUnique(m.Name)
	if err != nil {
		record.Warnf(i.IBMPowerVSImage, "FailedRetrieveImage", "Failed to retrieve image %q", m.Name)
		return nil, nil, err
	} else if imageReply != nil {
		i.Info("Image already exists")
		return imageReply, nil, nil
	}

	if lastJob, _ := i.GetImportJob(); lastJob != nil {
		if *lastJob.Status.State != "completed" && *lastJob.Status.State != "failed" {
			i.Info("Previous import job not yet finished", "state", *lastJob.Status.State)
			return nil, nil, nil
		}
	}

	body := &models.CreateCosImageImportJob{
		ImageName:     &m.Name,
		BucketName:    s.Bucket,
		BucketAccess:  core.StringPtr(BucketAccess),
		Region:        s.Region,
		ImageFilename: s.Object,
		StorageType:   s.StorageType,
	}

	jobRef, err := i.IBMPowerVSClient.CreateCosImage(body)
	if err != nil {
		i.Info("Unable to create new import job request")
		record.Warnf(i.IBMPowerVSImage, "FailedCreateImageImportJob", "Failed image import job creation - %v", err)
		return nil, nil, err
	}
	i.Info("New import job request created")
	record.Eventf(i.IBMPowerVSImage, "SuccessfulCreateImageImportJob", "Created image import job %q", *jobRef.ID)
	return nil, jobRef, nil
}

// PatchObject persists the cluster configuration and status.
func (i *PowerVSImageScope) PatchObject() error {
	return i.patchHelper.Patch(context.TODO(), i.IBMPowerVSImage)
}

// Close closes the current scope persisting the cluster configuration and status.
func (i *PowerVSImageScope) Close() error {
	return i.PatchObject()
}

// DeleteImage will delete the image.
func (i *PowerVSImageScope) DeleteImage() error {
	if err := i.IBMPowerVSClient.DeleteImage(i.IBMPowerVSImage.Status.ImageID); err != nil {
		record.Warnf(i.IBMPowerVSImage, "FailedDeleteImage", "Failed image deletion - %v", err)
		return err
	}
	record.Eventf(i.IBMPowerVSImage, "SuccessfulDeleteImage", "Deleted Image %q", i.IBMPowerVSImage.Status.ImageID)
	return nil
}

// GetImportJob will get the image import job.
func (i *PowerVSImageScope) GetImportJob() (*models.Job, error) {
	return i.IBMPowerVSClient.GetCosImages(i.IBMPowerVSImage.Spec.ServiceInstanceID)
}

// DeleteImportJob will delete the image import job.
func (i *PowerVSImageScope) DeleteImportJob() error {
	if err := i.IBMPowerVSClient.DeleteJob(i.IBMPowerVSImage.Status.JobID); err != nil {
		record.Warnf(i.IBMPowerVSImage, "FailedDeleteImageImportJob", "Failed image import job deletion - %v", err)
		return err
	}
	record.Eventf(i.IBMPowerVSImage, "SuccessfulDeleteImageImportJob", "Deleted image import job %q", i.IBMPowerVSImage.Status.JobID)
	return nil
}

// SetReady will set the status as ready for the image.
func (i *PowerVSImageScope) SetReady() {
	i.IBMPowerVSImage.Status.Ready = true
}

// SetNotReady will set the status as not ready for the image.
func (i *PowerVSImageScope) SetNotReady() {
	i.IBMPowerVSImage.Status.Ready = false
}

// IsReady will return the status for the image.
func (i *PowerVSImageScope) IsReady() bool {
	return i.IBMPowerVSImage.Status.Ready
}

// SetImageID will set the id for the image.
func (i *PowerVSImageScope) SetImageID(id *string) {
	if id != nil {
		i.IBMPowerVSImage.Status.ImageID = *id
	}
}

// GetImageID will get the id for the image.
func (i *PowerVSImageScope) GetImageID() string {
	return i.IBMPowerVSImage.Status.ImageID
}

// SetImageState will set the state for the image.
func (i *PowerVSImageScope) SetImageState(status string) {
	i.IBMPowerVSImage.Status.ImageState = infrav1beta2.PowerVSImageState(status)
}

// GetImageState will get the state for the image.
func (i *PowerVSImageScope) GetImageState() infrav1beta2.PowerVSImageState {
	return i.IBMPowerVSImage.Status.ImageState
}

// SetJobID will set the id for the import image job.
func (i *PowerVSImageScope) SetJobID(id string) {
	i.IBMPowerVSImage.Status.JobID = id
}

// GetJobID will get the id for the import image job.
func (i *PowerVSImageScope) GetJobID() string {
	return i.IBMPowerVSImage.Status.JobID
}
