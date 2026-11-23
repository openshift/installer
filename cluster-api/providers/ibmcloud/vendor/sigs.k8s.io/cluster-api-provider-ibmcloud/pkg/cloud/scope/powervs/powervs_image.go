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

package powervs

import (
	"context"
	"errors"
	"fmt"

	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"

	corev1 "k8s.io/api/core/v1"
	cgrecord "k8s.io/client-go/tools/record"
	"k8s.io/utils/ptr"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/powervs/v1beta3"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/endpoints"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/powervs"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/resourcecontroller"
)

// BucketAccess indicates if the bucket has public or private access public access.
const BucketAccess = "public"

var (
	// ErrWorkspaceNotInActiveState indicates error if workspace is inactive.
	ErrWorkspaceNotInActiveState = errors.New("workspace is not in active state")
)

// ImageScopeParams defines the input parameters used to create a new ImageScope.
type ImageScopeParams struct {
	Client          client.Client
	IBMPowerVSImage *infrav1.IBMPowerVSImage
	Zone            string
	ServiceEndpoint []endpoints.ServiceEndpoint
	ClientBuilder   ClientBuilder
	Recorder        cgrecord.EventRecorder
}

// ImageScope defines a scope defined around a Power VS Cluster.
type ImageScope struct {
	Client client.Client

	IBMPowerVSClient powervs.PowerVS
	ResourceClient   resourcecontroller.ResourceController

	// Recorder is the event recorder for emitting Kubernetes events.
	Recorder cgrecord.EventRecorder

	IBMPowerVSImage *infrav1.IBMPowerVSImage
	ServiceEndpoint []endpoints.ServiceEndpoint
	workspaceID     string
}

// NewPowerVSImageScope creates a new ImageScope from the supplied parameters.
func NewPowerVSImageScope(ctx context.Context, params ImageScopeParams) (*ImageScope, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	scope := &ImageScope{
		Client:          params.Client,
		IBMPowerVSImage: params.IBMPowerVSImage,
		Recorder:        params.Recorder,
	}

	if err := scope.initClients(ctx, &params); err != nil {
		return nil, fmt.Errorf("failed to initialize IBM Cloud clients for image: %w", err)
	}

	return scope, nil
}

// validate ensures all required fields are present before scope creation.
func (p *ImageScopeParams) validate() error {
	if p.Client == nil {
		return errors.New("failed to generate new scope: client is nil")
	}
	if p.IBMPowerVSImage == nil {
		return errors.New("failed to generate new scope: IBMPowerVSImage is nil")
	}
	if p.ClientBuilder == nil {
		return errors.New("failed to generate new scope: ClientBuilder is nil")
	}
	return nil
}

// initClients bootstraps the required IBM Cloud clients based on the current image state.
func (s *ImageScope) initClients(ctx context.Context, params *ImageScopeParams) error {
	log := ctrl.LoggerFrom(ctx)

	auth, err := params.ClientBuilder.GetAuthenticator(ctx)
	if err != nil {
		return fmt.Errorf("failed to create authenticator: %w", err)
	}

	opts := ClientOptions{
		Authenticator:   auth,
		ServiceEndpoint: params.ServiceEndpoint,
		Debug:           log.V(DEBUGLEVEL).Enabled(),
	}

	// 1. Build ResourceController
	s.ResourceClient, err = params.ClientBuilder.GetResourceControllerClient(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to create Resource Controller client: %w", err)
	}

	// 2. Resolve Workspace ID and Zone
	workspaceID, workspaceZone, err := s.resolveWorkspace(ctx, params.Zone)
	if err != nil {
		return err
	}

	// 3. Build PowerVS Client with the resolved Workspace ID
	opts.WorkspaceID = workspaceID
	opts.Zone = workspaceZone

	s.IBMPowerVSClient, err = params.ClientBuilder.GetPowerVSClient(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to create PowerVS client: %w", err)
	}

	s.workspaceID = workspaceID
	return nil
}

// resolveWorkspace figures out which workspace this image should belong to, validates it, and returns the ID and Zone.
func (s *ImageScope) resolveWorkspace(_ context.Context, specZone string) (string, string, error) {
	var workspaceID string

	// 1. Check if we were given an explicit ID or need to resolve a Name
	if s.IBMPowerVSImage.Spec.Workspace.ID != "" {
		workspaceID = s.IBMPowerVSImage.Spec.Workspace.ID
	} else {
		workspaceName := s.IBMPowerVSImage.Spec.Workspace.Name
		if workspaceName == "" {
			workspaceName = fmt.Sprintf("%s-workspace", s.IBMPowerVSImage.Spec.ClusterName)
		}

		filter := resourcecontroller.InstanceFilter{
			Name:           workspaceName,
			Zone:           &specZone,
			ResourceID:     resourcecontroller.PowerVSResourceID,
			ResourcePlanID: resourcecontroller.PowerVSResourcePlanID,
		}

		workspace, err := s.ResourceClient.GetResourceInstanceByFilter(filter)
		if err != nil {
			return "", "", fmt.Errorf("failed to look up workspace name %q: %w", workspaceName, err)
		}
		if workspace == nil || workspace.GUID == nil {
			return "", "", fmt.Errorf("PowerVS workspace %q is not yet created or GUID is nil", workspaceName)
		}
		if workspace.State == nil || *workspace.State != string(infrav1.WorkspaceStateActive) {
			return "", "", fmt.Errorf("%w: workspace %q", ErrWorkspaceNotInActiveState, workspaceName)
		}

		workspaceID = *workspace.GUID
	}

	// 2. Fetch the concrete resource to get the RegionID (which is the PowerVS Zone)
	res, _, err := s.ResourceClient.GetResourceInstance(&resourcecontrollerv2.GetResourceInstanceOptions{
		ID: &workspaceID,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to get resource instance by ID %q: %w", workspaceID, err)
	}
	if res.RegionID == nil {
		return "", "", fmt.Errorf("resource instance %q has no RegionID", workspaceID)
	}

	return workspaceID, *res.RegionID, nil
}

// GetOrImportImage verifies if the image exists, and if not, triggers a COS import job.
func (s *ImageScope) GetOrImportImage(ctx context.Context) (*models.ImageReference, *models.JobReference, error) {
	log := ctrl.LoggerFrom(ctx)
	imageSpec := s.IBMPowerVSImage.Spec
	imageName := s.IBMPowerVSImage.Name

	// 1. Idempotency Check
	imageReply, err := s.ensureImageUnique(ctx, imageName)
	if err != nil {
		s.Recorder.Eventf(s.IBMPowerVSImage, corev1.EventTypeWarning, "FailedRetrieveImage", "Failed to retrieve image %q", imageName)
		return nil, nil, fmt.Errorf("failed to verify image uniqueness: %w", err)
	} else if imageReply != nil {
		log.Info("Image already exists", "imageName", imageName)
		return imageReply, nil, nil
	}

	// 2. In-Progress Job Check
	lastJob, err := s.getImportJob(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to check for existing import job: %w", err)
	}

	if lastJob != nil && lastJob.Status != nil && lastJob.Status.State != nil {
		state := *lastJob.Status.State
		if state != string(infrav1.PowerVSImageStateCompleted) && state != string(infrav1.PowerVSImageStateFailed) {
			log.Info("Previous import job not yet finished, requeuing", "state", state)
			return nil, nil, nil
		}
		log.Info("Previous import job ended, triggering new import", "state", state)
	} else {
		log.Info("No prior import job found, triggering first import")
	}

	// 3. Trigger New Import Job
	body := &models.CreateCosImageImportJob{
		ImageName:     ptr.To(imageName),
		BucketName:    ptr.To(imageSpec.Bucket),
		BucketAccess:  ptr.To(BucketAccess),
		Region:        ptr.To(imageSpec.Region),
		ImageFilename: ptr.To(imageSpec.Object),
	}

	if imageSpec.StorageType != "" {
		body.StorageType = string(imageSpec.StorageType)
	}

	jobRef, err := s.IBMPowerVSClient.CreateCosImage(ctx, body)
	if err != nil {
		s.Recorder.Eventf(s.IBMPowerVSImage, corev1.EventTypeWarning, "FailedCreateImageImportJob", "Failed image import job creation: %v", err)
		return nil, nil, fmt.Errorf("failed to create COS image import job: %w", err)
	}

	log.Info("New import job request created", "jobID", *jobRef.ID)
	s.Recorder.Eventf(s.IBMPowerVSImage, corev1.EventTypeNormal, "SuccessfulCreateImageImportJob", "Created image import job %q", *jobRef.ID)
	return nil, jobRef, nil
}

// DeleteImage will delete the image.
func (s *ImageScope) DeleteImage(ctx context.Context) error {
	imageID := s.GetImageID()
	if imageID == "" {
		return nil
	}

	if err := s.IBMPowerVSClient.DeleteImage(ctx, imageID); err != nil {
		s.Recorder.Eventf(s.IBMPowerVSImage, corev1.EventTypeWarning, "FailedDeleteImage", "Failed image deletion: %v", err)
		return fmt.Errorf("failed to delete PowerVS image %s: %w", imageID, err)
	}

	s.Recorder.Eventf(s.IBMPowerVSImage, corev1.EventTypeNormal, "SuccessfulDeleteImage", "Deleted Image %q", imageID)
	return nil
}

// DeleteImportJob will delete the image import job.
func (s *ImageScope) DeleteImportJob(ctx context.Context) error {
	jobID := s.GetJobID()
	if jobID == "" {
		return nil
	}

	if err := s.IBMPowerVSClient.DeleteJob(ctx, jobID); err != nil {
		s.Recorder.Eventf(s.IBMPowerVSImage, corev1.EventTypeWarning, "FailedDeleteImageImportJob", "Failed image import job deletion: %v", err)
		return fmt.Errorf("failed to delete COS image import job %s: %w", jobID, err)
	}

	s.Recorder.Eventf(s.IBMPowerVSImage, corev1.EventTypeNormal, "SuccessfulDeleteImageImportJob", "Deleted image import job %q", jobID)
	return nil
}

// SetImageActive sets the image state to active.
func (s *ImageScope) SetImageActive() {
	s.IBMPowerVSImage.Status.ImageState = infrav1.PowerVSImageStateACTIVE
}

// IsImageActive returns true when the image state is active.
func (s *ImageScope) IsImageActive() bool {
	return s.IBMPowerVSImage.Status.ImageState == infrav1.PowerVSImageStateACTIVE
}

// SetImageID will set the id for the image.
func (s *ImageScope) SetImageID(id *string) {
	if id != nil {
		s.IBMPowerVSImage.Status.ImageID = *id
	}
}

// GetImageID will get the id for the image.
func (s *ImageScope) GetImageID() string {
	return s.IBMPowerVSImage.Status.ImageID
}

// SetImageState will set the state for the image.
func (s *ImageScope) SetImageState(status string) {
	s.IBMPowerVSImage.Status.ImageState = infrav1.PowerVSImageState(status)
}

// GetImageState will get the state for the image.
func (s *ImageScope) GetImageState() infrav1.PowerVSImageState {
	return s.IBMPowerVSImage.Status.ImageState
}

// SetJobID will set the id for the import image job.
func (s *ImageScope) SetJobID(id string) {
	s.IBMPowerVSImage.Status.JobID = id
}

// GetJobID will get the id for the import image job.
func (s *ImageScope) GetJobID() string {
	return s.IBMPowerVSImage.Status.JobID
}

// ensureImageUnique checks whether an image with the given name already exists
// in the workspace. Returns the existing reference if found, nil otherwise.
func (s *ImageScope) ensureImageUnique(ctx context.Context, imageName string) (*models.ImageReference, error) {
	images, err := s.IBMPowerVSClient.ListImages(ctx)
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

// getImportJob returns the current COS image import job for the workspace.
func (s *ImageScope) getImportJob(ctx context.Context) (*models.Job, error) {
	return s.IBMPowerVSClient.GetCosImages(ctx, s.workspaceID)
}
