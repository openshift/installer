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

package controllers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/IBM-Cloud/power-go-client/power/models"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1" //nolint:staticcheck
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	clusterv1util "sigs.k8s.io/cluster-api/util"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"         //nolint:staticcheck
	v1beta2conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions/v1beta2" //nolint:staticcheck
	v1beta1patch "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/patch"                   //nolint:staticcheck
	"sigs.k8s.io/cluster-api/util/finalizers"

	infrav1 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/endpoints"
)

// IBMPowerVSImageReconciler reconciles a IBMPowerVSImage object.
type IBMPowerVSImageReconciler struct {
	client.Client
	Recorder        record.EventRecorder
	ServiceEndpoint []endpoints.ServiceEndpoint
	Scheme          *runtime.Scheme
}

//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsimages,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsimages/status,verbs=get;update;patch

// Reconcile implements controller runtime Reconciler interface and handles reconciliation logic for IBMPowerVSImage.
func (r *IBMPowerVSImageReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	log := ctrl.LoggerFrom(ctx)

	log.Info("Reconciling IBMPowerVSImage")
	defer log.Info("Finished reconciling IBMPowerVSImage")

	// Fetch the IBMPowerVSImage.
	ibmPowerVSImage := &infrav1.IBMPowerVSImage{}
	err := r.Client.Get(ctx, req.NamespacedName, ibmPowerVSImage)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("IBMPowerVSImage not found")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, fmt.Errorf("failed to get IBMPowerVSImage: %w", err)
	}

	// Add finalizer first if not set to avoid the race condition between init and delete.
	if finalizerAdded, err := finalizers.EnsureFinalizer(ctx, r.Client, ibmPowerVSImage, infrav1.IBMPowerVSImageFinalizer); err != nil || finalizerAdded {
		return ctrl.Result{}, err
	}

	var cluster *infrav1.IBMPowerVSCluster
	scopeParams := scope.PowerVSImageScopeParams{
		Client:          r.Client,
		IBMPowerVSImage: ibmPowerVSImage,
		ServiceEndpoint: r.ServiceEndpoint,
	}

	// Externally managed clusters might not be available during image deletion. Get the cluster only when image is still not deleted.
	if ibmPowerVSImage.DeletionTimestamp.IsZero() {
		cluster, err = scope.GetClusterByName(ctx, r.Client, ibmPowerVSImage.Namespace, ibmPowerVSImage.Spec.ClusterName)
		if err != nil {
			return ctrl.Result{}, err
		}
		scopeParams.Zone = cluster.Spec.Zone
	}

	// Initialize the patch helper
	patchHelper, err := v1beta1patch.NewHelper(ibmPowerVSImage, r.Client)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to init patch helper: %w", err)
	}

	// Always attempt to Patch the IBMPowerVSImage object and status after each reconciliation.
	defer func() {
		if err := patchIBMPowerVSImage(ctx, patchHelper, ibmPowerVSImage); err != nil {
			reterr = kerrors.NewAggregate([]error{reterr, err})
		}
	}()

	// Create the scope
	imageScope, err := scope.NewPowerVSImageScope(ctx, scopeParams)
	if err != nil {
		if errors.Is(err, scope.ErrServiceInsanceNotInActiveState) {
			v1beta2conditions.Set(imageScope.IBMPowerVSImage, metav1.Condition{
				Type:   infrav1.WorkspaceReadyV1Beta2Condition,
				Status: metav1.ConditionFalse,
				Reason: infrav1.WorkspaceNotReadyV1Beta2Reason,
			})
		}
		return ctrl.Result{}, fmt.Errorf("failed to create scope: %w", err)
	}

	// Handle deleted clusters.
	if !ibmPowerVSImage.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, imageScope)
	}

	return r.reconcile(ctx, cluster, imageScope)
}

func (r *IBMPowerVSImageReconciler) reconcile(ctx context.Context, cluster *infrav1.IBMPowerVSCluster, imageScope *scope.PowerVSImageScope) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	// Create new labels section for IBMPowerVSImage metadata if nil.
	if imageScope.IBMPowerVSImage.Labels == nil {
		imageScope.IBMPowerVSImage.Labels = make(map[string]string)
	}

	if _, ok := imageScope.IBMPowerVSImage.Labels[clusterv1.ClusterNameLabel]; !ok {
		imageScope.IBMPowerVSImage.Labels[clusterv1.ClusterNameLabel] = imageScope.IBMPowerVSImage.Spec.ClusterName
	}

	if r.shouldAdopt(*imageScope.IBMPowerVSImage) {
		log.Info("Image Controller has not yet set OwnerRef")
		imageScope.IBMPowerVSImage.OwnerReferences = clusterv1util.EnsureOwnerRef(imageScope.IBMPowerVSImage.OwnerReferences, metav1.OwnerReference{
			APIVersion: infrav1.GroupVersion.String(),
			Kind:       "IBMPowerVSCluster",
			Name:       cluster.Name,
			UID:        cluster.UID,
		})
		return ctrl.Result{}, nil
	}

	v1beta2conditions.Set(imageScope.IBMPowerVSImage, metav1.Condition{
		Type:   infrav1.WorkspaceReadyV1Beta2Condition,
		Status: metav1.ConditionTrue,
		Reason: infrav1.WorkspaceReadyV1Beta2Reason,
	})

	if jobID := imageScope.GetJobID(); jobID != "" {
		job, err := imageScope.IBMPowerVSClient.GetJob(jobID)
		if err != nil {
			log.Info("Unable to get job details", "jobID", jobID)
			return ctrl.Result{RequeueAfter: 2 * time.Minute}, err
		}

		imageScope.SetImageState(*job.Status.State)
		switch imageScope.GetImageState() {
		case infrav1.PowerVSImageStateCompleted:
			v1beta1conditions.MarkTrue(imageScope.IBMPowerVSImage, infrav1.ImageImportedCondition)
			v1beta2conditions.Set(imageScope.IBMPowerVSImage, metav1.Condition{
				Type:   infrav1.IBMPowerVSImageReadyV1Beta2Condition,
				Status: metav1.ConditionTrue,
				Reason: infrav1.IBMPowerVSImageReadyV1Beta2Reason,
			})
		case infrav1.PowerVSImageStateFailed:
			imageScope.SetNotReady()
			imageScope.SetImageState(string(infrav1.PowerVSImageStateFailed))
			v1beta1conditions.MarkFalse(imageScope.IBMPowerVSImage, infrav1.ImageImportedCondition, infrav1.ImageImportFailedReason, clusterv1beta1.ConditionSeverityError, "%s", job.Status.Message)
			v1beta2conditions.Set(imageScope.IBMPowerVSImage, metav1.Condition{
				Type:   infrav1.IBMPowerVSImageReadyV1Beta2Condition,
				Status: metav1.ConditionFalse,
				Reason: infrav1.ImageImportFailedReason,
			})
			return ctrl.Result{RequeueAfter: 2 * time.Minute}, fmt.Errorf("failed to import image, message: %s", job.Status.Message)
		case infrav1.PowerVSImageStateQueued:
			imageScope.SetNotReady()
			imageScope.SetImageState(string(infrav1.PowerVSImageStateQueued))
			v1beta1conditions.MarkFalse(imageScope.IBMPowerVSImage, infrav1.ImageImportedCondition, string(infrav1.PowerVSImageStateQueued), clusterv1beta1.ConditionSeverityInfo, "%s", job.Status.Message)
			v1beta2conditions.Set(imageScope.IBMPowerVSImage, metav1.Condition{
				Type:   infrav1.IBMPowerVSImageReadyV1Beta2Condition,
				Status: metav1.ConditionFalse,
				Reason: infrav1.ImageQueuedReason,
			})
			return ctrl.Result{RequeueAfter: 2 * time.Minute}, nil
		default:
			imageScope.SetNotReady()
			imageScope.SetImageState(string(infrav1.PowerVSImageStateImporting))
			v1beta1conditions.MarkFalse(imageScope.IBMPowerVSImage, infrav1.ImageImportedCondition, *job.Status.State, clusterv1beta1.ConditionSeverityInfo, "%s", job.Status.Message)
			v1beta2conditions.Set(imageScope.IBMPowerVSImage, metav1.Condition{
				Type:   infrav1.IBMPowerVSImageReadyV1Beta2Condition,
				Status: metav1.ConditionFalse,
				Reason: infrav1.ImageNotReadyReason,
			})
			return ctrl.Result{RequeueAfter: 2 * time.Minute}, nil
		}
	}

	img, jobRef, err := r.getOrCreate(ctx, imageScope)
	if err != nil {
		log.Error(err, "Unable to import image")
		return ctrl.Result{}, fmt.Errorf("failed to reconcile Image for IBMPowerVSImage %s/%s: %w", imageScope.IBMPowerVSImage.Namespace, imageScope.IBMPowerVSImage.Name, err)
	}

	if jobRef != nil {
		imageScope.SetJobID(*jobRef.ID)
	}
	return reconcileImage(ctx, img, imageScope)
}

func reconcileImage(ctx context.Context, img *models.ImageReference, imageScope *scope.PowerVSImageScope) (_ ctrl.Result, reterr error) {
	log := ctrl.LoggerFrom(ctx)
	if img != nil {
		image, err := imageScope.IBMPowerVSClient.GetImage(*img.ImageID)
		if err != nil {
			log.Info("Unable to get image details", "imageID", *img.ImageID)
			return ctrl.Result{}, err
		}

		imageScope.SetImageID(image.ImageID)
		log.Info("ImageID", imageScope.GetImageID())
		imageScope.SetImageState(image.State)
		log.Info("ImageState", image.State)

		switch imageScope.GetImageState() {
		case infrav1.PowerVSImageStateQueued:
			log.Info("Image is in queued state")
			imageScope.SetNotReady()
			v1beta1conditions.MarkFalse(imageScope.IBMPowerVSImage, infrav1.ImageReadyCondition, infrav1.ImageNotReadyReason, clusterv1beta1.ConditionSeverityWarning, "")
			v1beta2conditions.Set(imageScope.IBMPowerVSImage, metav1.Condition{
				Type:   infrav1.IBMPowerVSImageReadyV1Beta2Condition,
				Status: metav1.ConditionFalse,
				Reason: infrav1.IBMPowerVSImageNotReadyV1Beta2Reason,
			})
			return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
		case infrav1.PowerVSImageStateACTIVE:
			log.Info("Image is in active state")
			imageScope.SetReady()
			v1beta1conditions.MarkTrue(imageScope.IBMPowerVSImage, infrav1.ImageReadyCondition)
			v1beta2conditions.Set(imageScope.IBMPowerVSImage, metav1.Condition{
				Type:   infrav1.IBMPowerVSImageReadyV1Beta2Condition,
				Status: metav1.ConditionTrue,
				Reason: infrav1.IBMPowerVSImageReadyV1Beta2Reason,
			})

		default:
			imageScope.SetNotReady()
			log.Info("PowerVS image state is undefined", "state", image.State, "image-id", imageScope.GetImageID())
			v1beta1conditions.MarkUnknown(imageScope.IBMPowerVSImage, infrav1.ImageReadyCondition, "", "")
			v1beta2conditions.Set(imageScope.IBMPowerVSImage, metav1.Condition{
				Type:   infrav1.IBMPowerVSImageReadyV1Beta2Condition,
				Status: metav1.ConditionFalse,
				Reason: infrav1.IBMPowerVSImageReadyUnknownV1Beta2Reason,
			})
		}
	}

	// Requeue after 1 minute if image is not ready to update status of the image properly.
	if !imageScope.IsReady() {
		log.Info("Image is not yet ready, requeue", "state", imageScope.GetImageState())
		return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
	}

	return ctrl.Result{}, nil
}

func (r *IBMPowerVSImageReconciler) reconcileDelete(ctx context.Context, scope *scope.PowerVSImageScope) (_ ctrl.Result, reterr error) {
	log := ctrl.LoggerFrom(ctx)
	log.Info("Handling deleted IBMPowerVSImage")

	v1beta1conditions.MarkFalse(scope.IBMPowerVSImage, infrav1.ImageReadyCondition, clusterv1.DeletingReason, clusterv1beta1.ConditionSeverityInfo, "")
	v1beta2conditions.Set(scope.IBMPowerVSImage, metav1.Condition{
		Type:   infrav1.IBMPowerVSImageReadyV1Beta2Condition,
		Status: metav1.ConditionFalse,
		Reason: infrav1.IBMPowerVSImageDeletingV1Beta2Reason,
	})

	defer func() {
		if reterr == nil {
			// IBMPowerVSImage is deleted so remove the finalizer.
			controllerutil.RemoveFinalizer(scope.IBMPowerVSImage, infrav1.IBMPowerVSImageFinalizer)
		}
	}()

	if scope.GetImageID() == "" {
		log.Info("IBMPowerVSImage ImageID is not yet set, hence not invoking the PowerVS API to delete the image")
		if scope.GetJobID() == "" {
			log.Info("JobID is not yet set, hence not invoking the PowerVS API to delete the image import job")
			return ctrl.Result{}, nil
		}
		if err := scope.DeleteImportJob(); err != nil {
			log.Error(err, "Error deleting IBMPowerVSImage Import Job")
			return ctrl.Result{}, fmt.Errorf("error deleting IBMPowerVSImage Import Job: %w", err)
		}
		return ctrl.Result{}, nil
	}

	if scope.IBMPowerVSImage.Spec.DeletePolicy != string(infrav1.DeletePolicyRetain) {
		if err := scope.DeleteImage(); err != nil {
			v1beta1conditions.MarkFalse(scope.IBMPowerVSImage, infrav1.ImageReadyCondition, clusterv1beta1.DeletionFailedReason, clusterv1beta1.ConditionSeverityWarning, "")
			v1beta2conditions.Set(scope.IBMPowerVSImage, metav1.Condition{
				Type:    infrav1.IBMPowerVSImageReadyV1Beta2Condition,
				Status:  metav1.ConditionFalse,
				Reason:  infrav1.IBMPowerVSMachineInstanceDeletingV1Beta2Reason,
				Message: fmt.Sprintf("failed to delete IBMPowerVSImage: %v", err),
			})
			return ctrl.Result{}, fmt.Errorf("error deleting IBMPowerVSImage %v: %w", klog.KObj(scope.IBMPowerVSImage), err)
		}
	}
	return ctrl.Result{}, nil
}

func (r *IBMPowerVSImageReconciler) getOrCreate(ctx context.Context, scope *scope.PowerVSImageScope) (*models.ImageReference, *models.JobReference, error) {
	image, job, err := scope.CreateImageCOSBucket(ctx)
	return image, job, err
}

func (r *IBMPowerVSImageReconciler) shouldAdopt(i infrav1.IBMPowerVSImage) bool {
	return !clusterv1util.HasOwner(i.OwnerReferences, infrav1.GroupVersion.String(), []string{"IBMPowerVSCluster"})
}

// SetupWithManager sets up the controller with the Manager.
func (r *IBMPowerVSImageReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1.IBMPowerVSImage{}).
		Complete(r)
}

func patchIBMPowerVSImage(ctx context.Context, patchHelper *v1beta1patch.Helper, ibmPowerVSImage *infrav1.IBMPowerVSImage) error {
	// Before computing ready condition, make sure that ImageReady is always set.
	// NOTE: This is required because v1beta2 conditions comply to guideline requiring conditions to be set at the
	// first reconcile.
	if c := v1beta2conditions.Get(ibmPowerVSImage, infrav1.IBMPowerVSImageReadyV1Beta2Condition); c == nil {
		if ibmPowerVSImage.Status.Ready {
			v1beta2conditions.Set(ibmPowerVSImage, metav1.Condition{
				Type:   infrav1.IBMPowerVSImageReadyV1Beta2Condition,
				Status: metav1.ConditionTrue,
				Reason: infrav1.IBMPowerVSImageReadyV1Beta2Reason,
			})
		} else {
			v1beta2conditions.Set(ibmPowerVSImage, metav1.Condition{
				Type:   infrav1.IBMPowerVSImageReadyV1Beta2Condition,
				Status: metav1.ConditionFalse,
				Reason: infrav1.IBMPowerVSImageNotReadyV1Beta2Reason,
			})
		}
	}

	// always update the readyCondition.
	v1beta1conditions.SetSummary(ibmPowerVSImage,
		v1beta1conditions.WithConditions(
			infrav1.ImageReadyCondition,
		),
	)

	if err := v1beta2conditions.SetSummaryCondition(ibmPowerVSImage, ibmPowerVSImage, infrav1.IBMPowerVSImageReadyCondition,
		v1beta2conditions.ForConditionTypes{
			infrav1.IBMPowerVSImageReadyV1Beta2Condition,
			infrav1.WorkspaceReadyV1Beta2Condition,
		},
		// Using a custom merge strategy to override reasons applied during merge.
		v1beta2conditions.CustomMergeStrategy{
			MergeStrategy: v1beta2conditions.DefaultMergeStrategy(
				// Use custom reasons.
				v1beta2conditions.ComputeReasonFunc(v1beta2conditions.GetDefaultComputeMergeReasonFunc(
					infrav1.IBMPowerVSImageNotReadyV1Beta2Reason,
					infrav1.IBMPowerVSImageReadyUnknownV1Beta2Reason,
					infrav1.IBMPowerVSImageReadyV1Beta2Reason,
				)),
			),
		},
	); err != nil {
		return fmt.Errorf("failed to set %s condition: %w", infrav1.IBMPowerVSImageReadyCondition, err)
	}

	// Patch the IBMPowerVSImage resource.
	return patchHelper.Patch(ctx, ibmPowerVSImage, v1beta1patch.WithOwnedV1Beta2Conditions{Conditions: []string{
		infrav1.IBMPowerVSImageReadyCondition,
		infrav1.IBMPowerVSImageReadyV1Beta2Condition,
		clusterv1beta1.PausedV1Beta2Condition,
		infrav1.WorkspaceReadyV1Beta2Condition,
	}})
}
